package controller

import (
	. "eaciit/wfdemo-git/library/core"
	. "eaciit/wfdemo-git/library/models"
	"eaciit/wfdemo-git/web/helper"
	"strings"
	// "time"
	// "fmt"
	"sort"

	"github.com/eaciit/crowd"
	"github.com/eaciit/knot/knot.v1"
	tk "github.com/eaciit/toolkit"
)

type AnalyticWindDistributionController struct {
	App
}

func CreateAnalyticWindDistributionController() *AnalyticWindDistributionController {
	var controller = new(AnalyticWindDistributionController)
	return controller
}

var windCats = [...]float64{1, 1.5, 2, 2.5, 3, 3.5, 4, 4.5, 5, 5.5, 6, 6.5, 7, 7.5, 8, 8.5, 9, 9.5, 10, 10.5, 11, 11.5, 12, 12.5, 13, 13.5, 14, 14.5, 15}

//var windCats = [...]float64{0,0.25,0.5,0.75,1,1.25,1.5,1.75, 2,2.25,2.5,2.75,	3,3.25,3.5,3.75,	4,4.25,4.5,4.75,	5,5.25,5.5,5.75,	6,6.25,6.5,6.75,	7,7.25,7.5,7.75,	8,8.25,8.5,8.75,	9,9.25,9.5,9.75,	10,10.25,10.5,10.75,	11,11.25,11.5,11.75,	12,12.25,12.5,12.75,	13,13.25,13.5,13.75,	14,14.25,14.5,14.75,	15}

func getWindDistrCategory(windValue float64) float64 {
	var datas float64

	for _, val := range windCats {
		if val >= windValue {
			datas = val
			return datas
		}
	}

	return datas
}

type ScadaAnalyticsWDData struct {
	Turbine  string
	Category float64
	Minutes  float64
}

func (m *AnalyticWindDistributionController) GetList(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	var dataSeries []tk.M

	p := new(PayloadAnalytic)
	e := k.GetPayload(&p)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	tStart, tEnd, e := helper.GetStartEndDate(k, p.Period, p.DateStart, p.DateEnd)
	query := []tk.M{}
	pipes := []tk.M{}
	query = append(query, tk.M{"_id": tk.M{"$ne": ""}})
	query = append(query, tk.M{"dateinfo.dateid": tk.M{"$gte": tStart}})
	query = append(query, tk.M{"dateinfo.dateid": tk.M{"$lte": tEnd}})
	query = append(query, tk.M{"avgwindspeed": tk.M{"$gte": 0.5}})
	if p.Project != "" {
		anProject := strings.Split(p.Project, "(")
		query = append(query, tk.M{"projectname": strings.TrimRight(anProject[0], " ")})
	}

	turbine := []string{}
	if len(p.Turbine) == 0 {
		pipes = append(pipes, tk.M{"$match": tk.M{"$and": query}})
		pipes = append(pipes, tk.M{"$group": tk.M{"_id": "$turbine"}})

		csr, _ := DB().Connection.NewQuery().From(new(ScadaData).TableName()).
			Command("pipe", pipes).Cursor(nil)
		_turbine := map[string]string{}
		for {
			e = csr.Fetch(&_turbine, 1, false)
			if e != nil {
				break
			}
			turbine = append(turbine, _turbine["_id"])
		}
		csr.Close()
	} else {
		bufferTurbine := []string{}
		for _, val := range p.Turbine {
			bufferTurbine = append(bufferTurbine, val.(string))
		}
		turbine = append(turbine, bufferTurbine...)
	}
	sort.Strings(turbine)

	type ScadaAnalyticsWDDataGroup struct {
		Turbine  string
		Category float64
	}

	type MiniScada struct {
		NacelDirection float64
		AvgWindSpeed   float64
		Turbine        string
	}
	tmpResult := []MiniScada{}
	_data := MiniScada{}
	for _, turbineX := range turbine {
		pipes = []tk.M{}
		tmpResult = []MiniScada{}
		queryT := query
		queryT = append(queryT, tk.M{"turbine": turbineX})
		pipes = append(pipes, tk.M{"$match": tk.M{"$and": queryT}})
		pipes = append(pipes, tk.M{"$project": tk.M{"turbine": 1, "avgwindspeed": 1}})
		csr, _ := DB().Connection.NewQuery().From(new(ScadaData).TableName()).
			Command("pipe", pipes).Cursor(nil)

		for {
			e = csr.Fetch(&_data, 1, false)
			if e != nil {
				break
			}
			tmpResult = append(tmpResult, _data)
		}
		csr.Close()

		if len(tmpResult) > 0 {
			datas := crowd.From(&tmpResult).Apply(func(x interface{}) interface{} {
				dt := x.(MiniScada)

				var di ScadaAnalyticsWDData
				di.Turbine = dt.Turbine
				di.Category = getWindDistrCategory(dt.AvgWindSpeed)
				di.Minutes = 1

				return di
			}).Exec().Group(func(x interface{}) interface{} {
				dt := x.(ScadaAnalyticsWDData)

				var dig ScadaAnalyticsWDDataGroup
				dig.Turbine = dt.Turbine
				dig.Category = dt.Category

				return dig
			}, nil).Exec()

			dts := datas.Apply(func(x interface{}) interface{} {
				kv := x.(crowd.KV)
				keys := kv.Key.(ScadaAnalyticsWDDataGroup)
				vs := kv.Value.([]ScadaAnalyticsWDData)
				total := len(vs)

				var di ScadaAnalyticsWDData
				di.Turbine = keys.Turbine
				di.Category = keys.Category
				di.Minutes = float64(total)

				return di
			}).Exec().Result.Data().([]ScadaAnalyticsWDData)

			totalMinutes := 0.0

			if len(dts) > 0 {
				totalMinutes = crowd.From(&dts).Sum(func(x interface{}) interface{} {
					dt := x.(ScadaAnalyticsWDData)
					return dt.Minutes
				}).Exec().Result.Sum
			}

			for _, wc := range windCats {
				exist := crowd.From(&dts).Where(func(x interface{}) interface{} {
					y := x.(ScadaAnalyticsWDData)
					Turbine := y.Turbine == turbineX
					Category := y.Category == wc
					return Turbine && Category
				}).Exec().Result.Data().([]ScadaAnalyticsWDData)

				distHelper := tk.M{}

				if len(exist) > 0 {
					distHelper.Set("Turbine", turbineX)
					distHelper.Set("Category", wc)

					Minute := crowd.From(&exist).Sum(func(x interface{}) interface{} {
						dt := x.(ScadaAnalyticsWDData)
						return dt.Minutes
					}).Exec().Result.Sum

					distHelper.Set("Contribute", Minute/totalMinutes)
				} else {
					distHelper.Set("Turbine", turbineX)
					distHelper.Set("Category", wc)
					distHelper.Set("Contribute", -0.0)
				}

				dataSeries = append(dataSeries, distHelper)
			}
		}
	}

	data := struct {
		Data []tk.M
	}{
		Data: dataSeries,
	}

	return helper.CreateResult(true, data, "success")
}

// maxWind := crowd.From(&resultScada).Max(func(x interface{}) interface{} {
// 			dt := x.(ScadaAnalyticsWDData)
// 			return dt.Category
// 		}).Exec().Result.Max

// var windCats = [...]float64{}

// for  i := 0 ; i <= 10 ;  i++ { //maxWind.(int)
// 	for  j := 0 ; j < 4 ;  j++ {
// 		windCats[i] = float64(i) + (float64(j)*0.25)
// 	}
// }
