package controller

import (
	. "eaciit/wfdemo-git/library/core"
	. "eaciit/wfdemo-git/library/models"
	"eaciit/wfdemo-git/web/helper"
	"sort"
	"time"

	"github.com/eaciit/knot/knot.v1"
	tk "github.com/eaciit/toolkit"
)

type XyAnalysis struct {
	App
}

func CreateXyAnalysisController() *XyAnalysis {
	var controller = new(XyAnalysis)
	return controller
}

type FieldAnalysis struct {
	Id    string
	Name  string
	Aggr  string
	Order int
}

type PayloadXyAnalysis struct {
	Period    string
	Project   string
	Engine    string
	Turbine   []interface{}
	DateStart time.Time
	DateEnd   time.Time
	XAxis     FieldAnalysis
	Y1Axis    FieldAnalysis
	Y2Axis    FieldAnalysis
}

func (m *XyAnalysis) GetXYFieldList(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson
	result := []FieldAnalysis{}

	//manual first
	result = append(result, FieldAnalysis{Id: "ActivePower_kW", Name: "ActivePower", Aggr: "$sum", Order: 1})
	result = append(result, FieldAnalysis{Id: "WindSpeed_ms", Name: "WindSpeed", Aggr: "$avg", Order: 2})
	result = append(result, FieldAnalysis{Id: "NacellePos", Name: "NacellePos", Aggr: "$avg", Order: 3})
	result = append(result, FieldAnalysis{Id: "WindDirection", Name: "WindDirection", Aggr: "$avg", Order: 4})
	result = append(result, FieldAnalysis{Id: "PitchAngle", Name: "PitchAngle", Aggr: "$avg", Order: 5})
	result = append(result, FieldAnalysis{Id: "TempOutdoor", Name: "TempOutdoor", Aggr: "$avg", Order: 6})

	return helper.CreateResult(true, result, "success")
}

func (m *XyAnalysis) GetData(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	p := new(PayloadXyAnalysis)
	e := k.GetPayload(&p)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	tStart, tEnd, e := helper.GetStartEndDate(k, p.Period, p.DateStart, p.DateEnd)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	var turbineList []TurbineOut
	if p.Project != "" {
		turbineList, _ = helper.GetTurbineList([]interface{}{p.Project})
	} else {
		turbineList, _ = helper.GetTurbineList(nil)
	}

	dataSeries := []tk.M{}

	matches := tk.M{"dateinfo.dateid": tk.M{"$gte": tStart, "$lte": tEnd}}
	matches.Set("projectname", p.Project)
	matches.Set("turbine", tk.M{}.Set("$in", p.Turbine))

	// ids := tk.M{"project": "$projectname", "turbine": "$turbine", "xaxis": "$" + p.XAxis.Id}

	// county1 := tk.M{"$cond": tk.M{}.
	// 	Set("if", tk.M{"$ifNull": []interface{}{"$" + p.Y1Axis.Id, false}}).
	// 	Set("then", 1).
	// 	Set("else", 0)}

	// county2 := tk.M{"$cond": tk.M{}.
	// 	Set("if", tk.M{"$ifNull": []interface{}{"$" + p.Y2Axis.Id, false}}).
	// 	Set("then", 1).
	// 	Set("else", 0)}

	// pipe := []tk.M{{"$match": matches},
	// 	{"$group": tk.M{
	// 		"_id":          ids,
	// 		"y1axis_sum":   tk.M{"$sum": "$" + p.Y1Axis.Id},
	// 		"y1axis_count": tk.M{"$sum": county1},
	// 		"y2axis_sum":   tk.M{"$sum": "$" + p.Y2Axis.Id},
	// 		"y2axis_count": tk.M{"$sum": county2},
	// 	}}}

	pipe := []tk.M{{"$match": matches},
		{"$project": tk.M{"projectname": 1, "turbine": 1, p.Y1Axis.Id: 1, p.Y2Axis.Id: 1, p.XAxis.Id: 1}}}

	csr, e := DB().Connection.NewQuery().
		From("Scada10MinHFD").
		Command("pipe", pipe).
		Cursor(nil)

	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}
	defer csr.Close()

	alltkm, allturb := tk.M{}, []string{}
	for {
		tkm, ds := tk.M{}, tk.M{}
		e = csr.Fetch(&tkm, 1, false)

		if e != nil {
			break
		}

		_key := ""
		for _, v := range turbineList {
			if tkm.GetString("turbine") == v.Value {
				// capacity = v.Capacity
				_key = v.Turbine
				// ds.Set("turbine", v.Turbine)
			}
		}

		ds = alltkm.Get(_key, tk.M{}).(tk.M)

		data1y := ds.Get("data1y", map[float64]float64{}).(map[float64]float64)
		data2y := ds.Get("data2y", map[float64]float64{}).(map[float64]float64)

		xaxis := tk.ToFloat64(tkm.Get(p.XAxis.Id), 2, tk.RoundingAuto)

		data1y[xaxis] = tkm.GetFloat64(p.Y1Axis.Id)
		data2y[xaxis] = tkm.GetFloat64(p.Y2Axis.Id)

		ds.Set("id", tkm.GetString("turbine"))
		ds.Set("turbine", _key)
		ds.Set("dashType", "solid")
		ds.Set("style", "smooth")
		ds.Set("type", "scatter")
		ds.Set("data1y", data1y)
		ds.Set("data2y", data2y)

		alltkm.Set(_key, ds)
	}

	for key, _ := range alltkm {
		allturb = append(allturb, key)
		ds := alltkm.Get(key, tk.M{}).(tk.M)
		data1y := ds.Get("data1y", map[float64]float64{}).(map[float64]float64)
		data2y := ds.Get("data2y", map[float64]float64{}).(map[float64]float64)

		_data1y, _data2y := make([][]float64, 0), make([][]float64, 0)

		for _x, _y := range data1y {
			_data1y = append(_data1y, []float64{_x, _y})
		}

		for _x, _y := range data2y {
			_data2y = append(_data2y, []float64{_x, _y})
		}

		ds.Set("data1y", _data1y)
		ds.Set("data2y", _data2y)

		alltkm.Set(key, ds)
	}

	sort.Strings(allturb)
	for _, turb := range allturb {
		dataSeries = append(dataSeries, alltkm.Get(turb, tk.M{}).(tk.M))
	}

	result := tk.M{}.Set("totalturbine", len(p.Turbine)).Set("data", dataSeries).Set("turbine", allturb)
	return helper.CreateResult(true, result, "success")
}
