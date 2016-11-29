package controller

import (
	. "eaciit/wfdemo-git-dev/library/core"
	. "eaciit/wfdemo-git-dev/library/models"
	"eaciit/wfdemo-git-dev/web/helper"

	// "time"
	// "fmt"

	"github.com/eaciit/crowd"
	"github.com/eaciit/dbox"
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

	var (
		filter     []*dbox.Filter
		dataSeries []tk.M
	)

	p := new(PayloadAnalytic)
	e := k.GetPayload(&p)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	// tStart, _ := time.Parse("2006-01-02", p.DateStart.UTC().Format("2006-01-02"))
	// tEnd, _ := time.Parse("2006-01-02 15:04:05", p.DateEnd.UTC().Format("2006-01-02")+" 23:59:59")
	tStart, tEnd, e := helper.GetStartEndDate(k, p.Period, p.DateStart, p.DateEnd)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}
	turbine := p.Turbine
	project := p.Project

	filter = append(filter, dbox.Ne("_id", ""))
	filter = append(filter, dbox.Gte("dateinfo.dateid", tStart))
	filter = append(filter, dbox.Lte("dateinfo.dateid", tEnd))
	if len(project) != 0 {
		filter = append(filter, dbox.Eq("projectname", project))
	}
	filter = append(filter, dbox.Gte("avgwindspeed", 0.5)) //Only >= 1
	if len(turbine) != 0 {
		filter = append(filter, dbox.In("turbine", turbine...))
	}

	csr, e := DB().Connection.NewQuery().
		From(new(ScadaData).TableName()).
		//Command("pipe", pipes).
		Where(dbox.And(filter...)).
		Cursor(nil)

	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}
	tmpResult := make([]ScadaData, 0)
	e = csr.Fetch(&tmpResult, 0, false)

	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	defer csr.Close()

	if len(p.Turbine) == 0 {
		for _, scadaVal := range tmpResult {
			exist := false
			for _, val := range turbine {
				if scadaVal.Turbine == val {
					exist = true
				}
			}
			if exist == false {
				turbine = append(turbine, scadaVal.Turbine)
			}
		}
	}

	type ScadaAnalyticsWDDataGroup struct {
		Turbine  string
		Category float64
	}

	if len(tmpResult) > 0 {
		datas := crowd.From(&tmpResult).Apply(func(x interface{}) interface{} {
			dt := x.(ScadaData)

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
			//minutes := crowd.From(&vs).Sum(func(x interface{}) interface{} {
			// 	dt := x.(ScadaAnalyticsWDData)
			// 	return dt.Minutes
			// }).Exec().Result.Sum

			var di ScadaAnalyticsWDData
			di.Turbine = keys.Turbine
			di.Category = keys.Category
			di.Minutes = float64(total)

			return di
		}).Exec().Result.Data().([]ScadaAnalyticsWDData)

		/*totalMinutes := crowd.From(&dts).Sum(func(x interface{}) interface{} {
			dt := x.(ScadaAnalyticsWDData)
			return dt.Minutes
		}).Exec().Result.Sum*/
		totalMinutes := 0.0

		for _, turbineX := range turbine {
			onotah := crowd.From(&dts).Where(func(x interface{}) interface{} {
				y := x.(ScadaAnalyticsWDData)
				Turbine := y.Turbine == turbineX
				return Turbine
			}).Exec().Result.Data().([]ScadaAnalyticsWDData)
			if len(onotah) > 0 {
				totalMinutes = crowd.From(&onotah).Sum(func(x interface{}) interface{} {
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

				//tk.Printf("dt %v\ntb %v\nct %v\n", exist, turbineX, wc)

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
