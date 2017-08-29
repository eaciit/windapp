package controller

import (
	. "eaciit/wfdemo-git/library/core"
	. "eaciit/wfdemo-git/library/models"
	"eaciit/wfdemo-git/web/helper"
	"math"
	// "time"
	// "fmt"
	"github.com/eaciit/dbox"
	"sort"

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
	Turbine    string
	Category   float64
	Contribute float64
}

func setContribution(turbine string, dataCatCount map[string]float64, countPerWSCat float64) (results []ScadaAnalyticsWDData) {
	results = []ScadaAnalyticsWDData{}
	for _, val := range windCats {
		results = append(results, ScadaAnalyticsWDData{
			Turbine:    turbine,
			Category:   val,
			Contribute: tk.Div(dataCatCount[tk.ToString(val)], countPerWSCat),
		})
	}
	return
}

func (m *AnalyticWindDistributionController) GetList(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	dataSeries := []ScadaAnalyticsWDData{}

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
	query = append(query, tk.M{"avgwindspeed": tk.M{"$gt": 0.5}})
	query = append(query, tk.M{"available": 1})
	if p.Project != "" {
		query = append(query, tk.M{"projectname": p.Project})
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
	turbineInt := []interface{}{}
	for _, val := range turbine {
		turbineInt = append(turbineInt, val)
	}

	type ScadaAnalyticsWDDataGroup struct {
		Turbine  string
		Category float64
	}

	type MiniScada struct {
		AvgWindSpeed float64
		Turbine      string
	}
	_data := MiniScada{}
	turbineName, e := helper.GetTurbineNameList(p.Project)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}
	queryT := []*dbox.Filter{}
	queryT = append(queryT, dbox.Gte("dateinfo.dateid", tStart))
	queryT = append(queryT, dbox.Lte("dateinfo.dateid", tEnd))
	queryT = append(queryT, dbox.Gt("avgwindspeed", 0.5))
	queryT = append(queryT, dbox.Eq("available", 1))
	if p.Project != "" {
		queryT = append(queryT, dbox.Eq("projectname", p.Project))
	}
	queryT = append(queryT, dbox.In("turbine", turbineInt...))

	csrData, _ := DB().Connection.NewQuery().
		Select("turbine", "avgwindspeed").
		From(new(ScadaData).TableName()).
		Where(dbox.And(queryT...)).
		Order("turbine").
		Cursor(nil)

	lastTurbine := ""
	_turbine := ""
	groupKey := ""
	countPerWSCat := 0.0
	dataCatCount := map[string]float64{}
	category := 0.0
	diff := windCats[1] - windCats[0]
	modus := 0.0
	pengurang := 1
	maxWS := windCats[len(windCats)-1]
	for {
		e = csrData.Fetch(&_data, 1, false)
		if e != nil {
			break
		}
		_turbine = turbineName[_data.Turbine]
		if lastTurbine != _turbine {
			if lastTurbine != "" {
				dataSeries = append(dataSeries, setContribution(lastTurbine, dataCatCount, countPerWSCat)...)
			}
			dataCatCount = map[string]float64{}
			lastTurbine = _turbine
			countPerWSCat = 0.0
		}
		countPerWSCat++
		if _data.AvgWindSpeed > maxWS {
			_data.AvgWindSpeed = maxWS
		}
		modus = math.Mod(_data.AvgWindSpeed, diff)
		if modus == 0 {
			pengurang = 2
		} else {
			pengurang = 1
		}
		category = windCats[int(tk.Div(_data.AvgWindSpeed, diff))-pengurang]
		groupKey = tk.ToString(category)
		dataCatCount[groupKey] = dataCatCount[groupKey] + 1
	}
	if lastTurbine != "" {
		dataSeries = append(dataSeries, setContribution(lastTurbine, dataCatCount, countPerWSCat)...)
	}
	csrData.Close()

	data := struct {
		Data []ScadaAnalyticsWDData
	}{
		Data: dataSeries,
	}

	return helper.CreateResult(true, data, "success")
}
