package controller

import (
	. "eaciit/wfdemo-git/library/core"
	. "eaciit/wfdemo-git/library/models"
	"eaciit/wfdemo-git/web/helper"
	"sort"
	"strings"
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
	Text  string
	Min   float64
	Max   float64
	Order int
}

type PayloadXyAnalysis struct {
	Period    string
	Project   []string
	Engine    string
	Turbine   []interface{}
	DateStart time.Time
	DateEnd   time.Time
	XAxis     FieldAnalysis
	Y1Axis    FieldAnalysis
	Y2Axis    FieldAnalysis
}

type byOrder []FieldAnalysis

func (s byOrder) Len() int {
	return len(s)
}
func (s byOrder) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s byOrder) Less(i, j int) bool {
	return s[i].Order < s[j].Order
}

// func (m *XyAnalysis) GetXYFieldList(k *knot.WebContext) interface{} {
// 	k.Config.OutputType = knot.OutputJson
// 	result := []FieldAnalysis{}

// 	//manual first
// 	result = append(result, FieldAnalysis{Id: "ActivePower_kW", Name: "ActivePower", Aggr: "$sum", Order: 1, Text: "ActivePower (kWh)"})
// 	result = append(result, FieldAnalysis{Id: "WindSpeed_ms", Name: "WindSpeed", Aggr: "$avg", Order: 2, Text: "WindSpeed (m/s)"})
// 	result = append(result, FieldAnalysis{Id: "NacellePos", Name: "NacellePos", Aggr: "$avg", Order: 3, Text: "NacellePos (-)"})
// 	result = append(result, FieldAnalysis{Id: "WindDirection", Name: "WindDirection", Aggr: "$avg", Order: 4, Text: "WindDirection (-)"})
// 	result = append(result, FieldAnalysis{Id: "PitchAngle", Name: "PitchAngle", Aggr: "$avg", Order: 5, Text: "PitchAngle (-)"})
// 	result = append(result, FieldAnalysis{Id: "TempOutdoor", Name: "TempOutdoor", Aggr: "$avg", Order: 6, Text: "TempOutdoor (&deg;C)"})

// 	return helper.CreateResult(true, result, "success")
// }

func (m *XyAnalysis) GetXYFieldList(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson
	result := []FieldAnalysis{}

	p := new(PayloadXyAnalysis)
	e := k.GetPayload(&p)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	pipe := []tk.M{}
	matches := tk.M{}.Set("enable", true)
	if len(p.Project) >= 1 {
		matches.Set("projectname", p.Project[0])
		if len(p.Project) > 1 {
			matches.Set("projectname", tk.M{}.Set("$in", p.Project))
		}
	}

	pipe = append(pipe, tk.M{"$match": matches})
	pipe = append(pipe, tk.M{"$group": tk.M{
		"_id":   tk.M{"realtime": "$realtimefield", "name": "$label", "units": "$units"},
		"order": tk.M{"$max": "$order"},
	}})

	csr, e := DB().Connection.NewQuery().
		From("ref_databrowsertag").
		Command("pipe", pipe).
		Cursor(nil)

	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}
	defer csr.Close()

	for {
		tkm := tk.M{}
		e = csr.Fetch(&tkm, 1, false)

		if e != nil {
			break
		}

		_id := tkm.Get("_id", tk.M{}).(tk.M)

		_fa := FieldAnalysis{
			Id:    _id.GetString("realtime"),
			Name:  _id.GetString("name"),
			Order: tkm.GetInt("order"),
			Text:  tk.Sprintf("%s (%s)", _id.GetString("name"), _id.GetString("units")),
		}

		result = append(result, _fa)
	}

	sort.Sort(byOrder(result))

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
	if len(p.Project) == 1 {
		turbineList, _ = helper.GetTurbineList([]interface{}{p.Project[0]})
	} else {
		turbineList, _ = helper.GetTurbineList(nil)
	}

	dataSeries := []tk.M{}

	matches := tk.M{"dateinfo.dateid": tk.M{"$gte": tStart, "$lte": tEnd}}
	if len(p.Project) == 1 {
		matches.Set("projectname", p.Project[0])
	} else {
		matches.Set("projectname", tk.M{}.Set("$in", p.Project))
	}
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
		{"$project": tk.M{"projectname": 1, "turbine": 1, strings.ToLower(p.Y1Axis.Id): 1, strings.ToLower(p.Y2Axis.Id): 1, strings.ToLower(p.XAxis.Id): 1}}}

	csr, e := DB().Connection.NewQuery().
		From("Scada10MinHFD").
		Command("pipe", pipe).
		Cursor(nil)

	// tk.Println(pipe)

	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}
	defer csr.Close()

	alltkm, allturb := tk.M{}, []string{}
	p.XAxis.Min, p.XAxis.Max = 0, 0
	p.Y1Axis.Min, p.Y1Axis.Max = 0, 0
	p.Y2Axis.Min, p.Y2Axis.Max = 0, 0
	for {
		tkm, ds := tk.M{}, tk.M{}
		e = csr.Fetch(&tkm, 1, false)

		if e != nil {
			break
		}

		// tk.Println(tkm)

		_key := ""
		for _, v := range turbineList {
			if tkm.GetString("turbine") == v.Value {
				// capacity = v.Capacity
				_key = v.Turbine
				// ds.Set("turbine", v.Turbine)
			}
		}

		ds = alltkm.Get(_key, tk.M{}).(tk.M)

		data1y := ds.Get("data1y", map[float64]interface{}{}).(map[float64]interface{})
		data2y := ds.Get("data2y", map[float64]interface{}{}).(map[float64]interface{})

		xaxis := tk.ToFloat64(tkm.Get(strings.ToLower(p.XAxis.Id)), 2, tk.RoundingAuto)

		data1y[xaxis] = tkm.Get(strings.ToLower(p.Y1Axis.Id), nil)
		data2y[xaxis] = tkm.Get(strings.ToLower(p.Y2Axis.Id), nil)

		if data1y[xaxis] != nil {
			_val := tk.ToFloat64(data1y[xaxis], 2, tk.RoundingAuto)
			data1y[xaxis] = _val
			if p.Y1Axis.Min == 0 || p.Y1Axis.Min > _val {
				p.Y1Axis.Min = _val
			}

			if p.Y1Axis.Max == 0 || p.Y1Axis.Max < _val {
				p.Y1Axis.Max = _val
			}
		}

		if data2y[xaxis] != nil {
			_val := tk.ToFloat64(data2y[xaxis], 2, tk.RoundingAuto)
			data2y[xaxis] = _val

			if p.Y2Axis.Min == 0 || p.Y2Axis.Min > _val {
				p.Y2Axis.Min = _val
			}

			if p.Y2Axis.Max == 0 || p.Y2Axis.Max < _val {
				p.Y2Axis.Max = _val
			}
		}

		if p.XAxis.Min == 0 || p.XAxis.Min > xaxis {
			p.XAxis.Min = xaxis
		}

		if p.XAxis.Max == 0 || p.XAxis.Max < xaxis {
			p.XAxis.Max = xaxis
		}

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
		data1y := ds.Get("data1y", map[float64]interface{}{}).(map[float64]interface{})
		data2y := ds.Get("data2y", map[float64]interface{}{}).(map[float64]interface{})

		_data2y := make([][]interface{}, 0)

		// _data1y, _data2y := make([][]float64, 0), make([][]float64, 0)

		for _x, _y := range data1y {
			_y2 := data2y[_x]
			_data2y = append(_data2y, []interface{}{_x, _y, _y2})
		}

		// for _x, _y := range data2y {
		// 	_data2y = append(_data2y, []float64{_x, _y})
		// }

		// ds.Set("data1y", _data1y)
		// ds.Set("data2y", _data2y)

		ds.Set("detail", _data2y)

		ds.Unset("data1y")
		ds.Unset("data2y")

		alltkm.Set(key, ds)
	}

	sort.Strings(allturb)
	for _, turb := range allturb {
		dataSeries = append(dataSeries, alltkm.Get(turb, tk.M{}).(tk.M))
	}

	result := tk.M{}.
		Set("totalturbine", len(p.Turbine)).
		Set("data", dataSeries).
		Set("turbine", allturb).
		Set("axisinfo", []FieldAnalysis{p.XAxis, p.Y1Axis, p.Y2Axis})
	return helper.CreateResult(true, result, "success")
}
