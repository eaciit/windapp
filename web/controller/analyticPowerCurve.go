package controller

import (
	. "github.com/eaciit/windapp/library/core"
	. "github.com/eaciit/windapp/library/helper"
	. "github.com/eaciit/windapp/library/models"
	"github.com/eaciit/windapp/web/helper"

	// "fmt"
	"time"

	"github.com/eaciit/crowd"
	"github.com/eaciit/dbox"
	"github.com/eaciit/knot/knot.v1"
	tk "github.com/eaciit/toolkit"
)

type AnalyticPowerCurveController struct {
	App
}

var (
	colorField            = [...]string{"#ff880e", "#21c4af", "#ff7663", "#ffb74f", "#a2df53", "#1c9ec4", "#ff63a5", "#f44336", "#D91E18", "#8877A9", "#9A12B3", "#26C281", "#E7505A", "#C49F47", "#ff5597", "#c3260c", "#d4735e", "#ff2ad7", "#34ac8b", "#11b2eb", "#f35838", "#ff0037", "#507ca3", "#ff6565", "#ffd664", "#72aaff", "#795548"}
	colorFieldDegradation = [...]string{"#ffcf9e", "#a6e7df", "#ffc8c0", "#ffe2b8", "#d9f2ba", "#a4d8e7", "#ffc0db", "#fab3ae", "#efa5a2", "#cfc8dc", "#d6a0e0", "#a8e6cc", "#f5b9bd", "#e7d8b5", "#ffbbd5", "#e7a89d", "#edc7be", "#ffa9ef", "#adddd0", "#9fe0f7", "#fabcaf", "#ff99af", "#b9cada", "#ffc1c1", "#ffeec1", "#c6ddff", "#c9bbb5"}
	downColor             = [...]string{"#000", "#444", "#666", "#888", "#aaa", "#ccc", "#eee"}
	// downIcon   = [...]string{"triangle", "square", "triangle", "cross", "square", "triangle", "cross"}
)

func CreateAnalyticPowerCurveController() *AnalyticPowerCurveController {
	var controller = new(AnalyticPowerCurveController)
	return controller
}

func (m *AnalyticPowerCurveController) GetList(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	var (
		pipes      []tk.M
		filter     []*dbox.Filter
		list       []tk.M
		dataSeries []tk.M
	)

	p := new(PayloadAnalytic)
	e := k.GetPayload(&p)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	tStart, _ := time.Parse("2006-01-02", p.DateStart.UTC().Format("2006-01-02"))
	tEnd, _ := time.Parse("2006-01-02 15:04:05", p.DateEnd.UTC().Format("2006-01-02")+" 23:59:59")
	turbine := p.Turbine
	project := p.Project
	isClean := p.IsClean
	isAverage := p.IsAverage

	pcData, e := getPCData(project)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	dataSeries = append(dataSeries, pcData)

	pipes = append(pipes, tk.M{"$unwind": "$dataitems"})
	pipes = append(pipes, tk.M{"$project": tk.M{
		"dateid":      "$dateinfo.dateid",
		"projectname": "$projectname",
		"turbineid":   "$turbineid",
		"wsclass":     "$dataitems.wsclass",
		"production":  "$dataitems.production",
		"totaldata":   "$dataitems.totaldata",
	}})
	pipes = append(pipes, tk.M{"$group": tk.M{"_id": "$wsclass", "production": tk.M{"$sum": "$production"}, "totaldata": tk.M{"$sum": "$totaldata"}}})
	pipes = append(pipes, tk.M{"$sort": tk.M{"_id": 1}})

	var collName string

	if isClean && isAverage {
		collName = "rpt_scadapowercurve_genonly_avg"
	} else if !isClean && isAverage {
		collName = "rpt_scadapowercurve_avg"
	} else if isClean && !isAverage {
		collName = "rpt_scadapowercurve_genonly"
	} else if !isClean && !isAverage {
		collName = "rpt_scadapowercurve"
	}

	selArr := 1
	for _, turbineX := range turbine {
		filter = nil
		filter = append(filter, dbox.Gte("dateinfo.dateid", tStart))
		filter = append(filter, dbox.Lte("dateinfo.dateid", tEnd))
		filter = append(filter, dbox.Lte("projectname", project))
		filter = append(filter, dbox.Eq("turbineid", turbineX))

		csr, e := DB().Connection.NewQuery().
			From(collName).
			Command("pipe", pipes).
			Where(dbox.And(filter...)).
			Cursor(nil)

		if e != nil {
			return helper.CreateResult(false, nil, e.Error())
		}

		e = csr.Fetch(&list, 0, false)

		if e != nil {
			return helper.CreateResult(false, nil, e.Error())
		}

		defer csr.Close()

		var datas [][]float64
		turbineData := tk.M{}
		turbineData.Set("name", turbineX)
		turbineData.Set("type", "scatterLine")
		turbineData.Set("style", "smooth")
		turbineData.Set("dashType", "solid")
		turbineData.Set("markers", tk.M{"visible": false})
		turbineData.Set("width", 2)
		turbineData.Set("color", colorField[selArr])

		for _, val := range list {
			// tk.Printf("%#v\n", val)
			datas = append(datas, []float64{val.GetFloat64("_id"), tk.Div(val.GetFloat64("production"), val.GetFloat64("totaldata"))})
		}

		if len(datas) > 0 {
			turbineData.Set("data", datas)
		}

		dataSeries = append(dataSeries, turbineData)
		selArr++
	}

	// tk.Printf("%#v\n", pipes)

	data := struct {
		Data []tk.M
	}{
		Data: dataSeries,
	}

	return helper.CreateResult(true, data, "success")
}

func (m *AnalyticPowerCurveController) GetListDensity(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	var (
		pipes      []tk.M
		filter     []*dbox.Filter
		list       []tk.M
		dataSeries []tk.M
	)

	p := new(PayloadAnalyticPC)
	e := k.GetPayload(&p)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	tStart, _ := time.Parse("2006-01-02", p.DateStart.UTC().Format("2006-01-02"))
	tEnd, _ := time.Parse("2006-01-02 15:04:05", p.DateEnd.UTC().Format("2006-01-02")+" 23:59:59")
	turbine := p.Turbine
	project := p.Project
	IsDeviation := p.IsDeviation
	DeviationVal := p.DeviationVal
	// isClean := p.IsClean
	// isAverage := p.IsAverage

	pcData, e := getPCData(project)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	dataSeries = append(dataSeries, pcData)

	// pipes = append(pipes, tk.M{"$group": tk.M{"_id": "$wsclass", "production": tk.M{"$sum": "$production"}, "totaldata": tk.M{"$sum": "$totaldata"}}})
	pipes = append(pipes, tk.M{"$group": tk.M{"_id": "$denadjwindspeed", "production": tk.M{"$avg": "$denpower"}, "totaldata": tk.M{"$sum": 1}}})
	pipes = append(pipes, tk.M{"$sort": tk.M{"_id": 1}})

	var collName string

	collName = "ScadaData"

	dVal := (tk.ToFloat64(tk.ToInt(DeviationVal, tk.RoundingAuto), 2, tk.RoundingUp) / 100.0)

	selArr := 1
	for _, turbineX := range turbine {

		filter = nil
		filter = append(filter, dbox.Gte("dateinfo.dateid", tStart))
		filter = append(filter, dbox.Lte("dateinfo.dateid", tEnd))
		filter = append(filter, dbox.Lte("projectname", project))
		filter = append(filter, dbox.Eq("turbine", turbineX))
		if !IsDeviation {
			filter = append(filter, dbox.Gte("dendeviationpct", dVal))
		}

		csr, e := DB().Connection.NewQuery().
			From(collName).
			Command("pipe", pipes).
			Where(dbox.And(filter...)).
			Cursor(nil)

		if e != nil {
			return helper.CreateResult(false, nil, e.Error())
		}
		e = csr.Fetch(&list, 0, false)
		defer csr.Close()

		var datas [][]float64
		turbineData := tk.M{}
		turbineData.Set("name", turbineX)
		turbineData.Set("type", "scatterLine")
		turbineData.Set("style", "smooth")
		turbineData.Set("dashType", "solid")
		turbineData.Set("markers", tk.M{"visible": false})
		turbineData.Set("width", 2)
		turbineData.Set("color", colorField[selArr])

		for _, val := range list {
			// tk.Printf("%#v\n", val)
			datas = append(datas, []float64{val.GetFloat64("_id"), tk.Div(val.GetFloat64("production"), val.GetFloat64("totaldata"))})
		}

		if len(datas) > 0 {
			turbineData.Set("data", datas)
		}

		dataSeries = append(dataSeries, turbineData)
		selArr++
	}

	// tk.Printf("%#v\n", pipes)

	data := struct {
		Data []tk.M
	}{
		Data: dataSeries,
	}

	return helper.CreateResult(true, data, "success")
}

func (m *AnalyticPowerCurveController) GetListPowerCurveScada(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	var (
		pipes      []tk.M
		filter     []*dbox.Filter
		list       []tk.M
		dataSeries []tk.M
	)

	p := new(PayloadAnalyticPC)
	e := k.GetPayload(&p)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	tStart, tEnd, e := helper.GetStartEndDate(k, p.Period, p.DateStart, p.DateEnd)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}
	turbine := p.Turbine
	project := p.Project
	IsDeviation := p.IsDeviation
	DeviationVal := p.DeviationVal
	viewSession := p.ViewSession
	isClean := p.IsClean

	colId := "$wsavgforpc"
	colValue := "$power"
	colDeviation := "deviationpct"
	switch viewSession {
	case "density":
		colId = "$denadjwindspeed"
		colValue = "$denpower"
		colDeviation = "dendeviationpct"
	case "adj":
		colId = "$wsadjforpc"
		colValue = "$power"
	default:
		colId = "$wsavgforpc"
		colDeviation = "deviationpct"
	}

	pcData, e := getPCData(project)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	dataSeries = append(dataSeries, pcData)

	pipes = append(pipes, tk.M{"$group": tk.M{"_id": tk.M{"colId": colId, "Turbine": "$turbine"}, "production": tk.M{"$avg": colValue}, "totaldata": tk.M{"$sum": 1}}})
	pipes = append(pipes, tk.M{"$sort": tk.M{"_id": 1}})

	var collName string
	collName = "ScadaData"
	dVal := (tk.ToFloat64(tk.ToInt(DeviationVal, tk.RoundingAuto), 2, tk.RoundingUp) / 100.0)
	selArr := 1

	filter = nil
	filter = append(filter, dbox.Ne("_id", ""))
	filter = append(filter, dbox.Gte("dateinfo.dateid", tStart))
	filter = append(filter, dbox.Lte("dateinfo.dateid", tEnd))
	filter = append(filter, dbox.Ne("turbine", ""))
	if !IsDeviation {
		filter = append(filter, dbox.Gte(colDeviation, dVal))
	}
	if isClean {
		filter = append(filter, dbox.Eq("oktime", 600))
	}

	csr, e := DB().Connection.NewQuery().
		From(collName).
		Command("pipe", pipes).
		Where(dbox.And(filter...)).
		Cursor(nil)

	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}
	e = csr.Fetch(&list, 0, false)
	defer csr.Close()

	if len(p.Turbine) == 0 {
		for _, listVal := range list {
			exist := false
			for _, val := range turbine {
				if listVal["_id"].(tk.M)["Turbine"] == val {
					exist = true
				}
			}
			if exist == false {
				turbine = append(turbine, listVal["_id"].(tk.M)["Turbine"])
			}
		}
	}

	for _, turbineX := range turbine {

		exist := crowd.From(&list).Where(func(x interface{}) interface{} {
			y := x.(tk.M)
			id := y.Get("_id").(tk.M)

			return id.GetString("Turbine") == turbineX
		}).Exec().Result.Data().([]tk.M)

		var datas [][]float64
		turbineData := tk.M{}
		turbineData.Set("name", turbineX)
		turbineData.Set("type", "scatterLine")
		turbineData.Set("style", "smooth")
		turbineData.Set("dashType", "solid")
		turbineData.Set("markers", tk.M{"visible": false})
		turbineData.Set("width", 2)
		turbineData.Set("color", colorField[selArr])

		for _, val := range exist {
			idD := val.Get("_id").(tk.M)
			datas = append(datas, []float64{idD.GetFloat64("colId"), tk.Div(val.GetFloat64("production"), val.GetFloat64("totaldata"))})
		}

		if len(datas) > 0 {
			turbineData.Set("data", datas)
		}

		dataSeries = append(dataSeries, turbineData)
		selArr++
	}

	data := struct {
		Data []tk.M
	}{
		Data: dataSeries,
	}

	return helper.CreateResult(true, data, "success")
}

func (m *AnalyticPowerCurveController) GetPowerCurve(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	var (
		list       []ScadaData
		listAlarm  []Alarm
		dataSeries []tk.M
	)

	p := new(PayloadAnalyticPC)
	e := k.GetPayload(&p)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	tStart, tEnd, e := helper.GetStartEndDate(k, p.Period, p.DateStart, p.DateEnd)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}
	turbine := p.Turbine
	project := p.Project
	colors := p.Color
	colordeg := p.ColorDeg
	IsDeviation := p.IsDeviation
	DeviationVal := p.DeviationVal
	viewSession := p.ViewSession
	isClean := p.IsClean
	dVal := (tk.ToFloat64(tk.ToInt(DeviationVal, tk.RoundingAuto), 2, tk.RoundingUp) / 100.0)

	colDeviation := "deviationpct"
	switch viewSession {
	case "density":
		colDeviation = "dendeviationpct"
	case "adj":
		colDeviation = "deviationpct"
	default:
		colDeviation = "deviationpct"
	}

	selArr := 0
	for _, turbineX := range turbine {
		var filter []*dbox.Filter
		filter = append(filter, dbox.Gte("dateinfo.dateid", tStart))
		filter = append(filter, dbox.Lte("dateinfo.dateid", tEnd))
		filter = append(filter, dbox.Eq("turbine", turbineX))
		filter = append(filter, dbox.Eq("projectname", project))

		if !IsDeviation {
			filter = append(filter, dbox.Gte(colDeviation, dVal))
		}
		if isClean {
			filter = append(filter, dbox.Eq("oktime", 600))
		}
		filter = append(filter, dbox.Ne("_id", ""))

		csr, e := DB().Connection.NewQuery().From(new(ScadaData).TableName()).Where(dbox.And(filter...)).Take(10000).Cursor(nil)
		if e != nil {
			return helper.CreateResult(false, nil, e.Error())
		}

		e = csr.Fetch(&list, 0, false)

		if e != nil {
			return helper.CreateResult(false, nil, e.Error())
		}

		defer csr.Close()

		turbineData := tk.M{}
		turbineData.Set("name", "Scatter-"+turbineX.(string))
		turbineData.Set("xField", "WindSpeed")
		turbineData.Set("yField", "Power")
		turbineData.Set("colorField", "valueColor")
		turbineData.Set("type", "scatter")
		turbineData.Set("markers", tk.M{"size": 2})

		datas := tk.M{}
		arrDatas := []tk.M{}
		for _, val := range list {
			datas = tk.M{}

			switch viewSession {
			case "density":
				if val.DenWindSpeed > 0 && val.DenPower > 0 {
					datas.Set("WindSpeed", val.DenWindSpeed)
					datas.Set("Power", val.DenPower)

					if val.DenDeviationPct <= dVal {
						datas.Set("valueColor", colordeg[selArr])
					} else {
						datas.Set("valueColor", colors[selArr])
					}

					arrDatas = append(arrDatas, datas)
				}
			default:
				if val.AvgWindSpeed > 0 && val.Power > 0 {

					datas.Set("WindSpeed", val.AvgWindSpeed)
					datas.Set("Power", val.Power)
					if val.DeviationPct <= dVal {
						datas.Set("valueColor", colordeg[selArr])
					} else {
						datas.Set("valueColor", colors[selArr])
					}

					arrDatas = append(arrDatas, datas)
				}
			}
		}

		turbineData.Set("data", arrDatas)
		dataSeries = append(dataSeries, turbineData)
		selArr++

		if p.IsDownTime {
			for idx, dw := range helper.DownTypes {
				down := dw.GetString("down")
				var filterAlarm []*dbox.Filter
				filterAlarm = append(filterAlarm, dbox.Gte("startdateinfo.dateid", tStart))
				filterAlarm = append(filterAlarm, dbox.Lte("startdateinfo.dateid", tEnd))
				filterAlarm = append(filterAlarm, dbox.Eq("turbine", turbineX))
				filterAlarm = append(filterAlarm, dbox.Eq("projectname", project))
				filterAlarm = append(filterAlarm, dbox.Eq(down, true))

				csr, e := DB().Connection.NewQuery().From(new(Alarm).TableName()).Where(dbox.And(filterAlarm...)).Cursor(nil)

				if e != nil {
					return helper.CreateResult(false, nil, e.Error())
				}

				e = csr.Fetch(&listAlarm, 0, false)

				if e != nil {
					return helper.CreateResult(false, nil, e.Error())
				}

				defer csr.Close()

				turbineDataAlarm := tk.M{}
				turbineDataAlarm.Set("name", down)
				turbineDataAlarm.Set("type", "scatter")
				turbineDataAlarm.Set("color", downColor[idx])
				turbineDataAlarm.Set("markers", tk.M{
					"size":       2,
					"background": downColor[idx],
				})

				var datasDown [][]float64
				if len(list) > 0 {
					for _, alarm := range listAlarm {
						startDate := GetDateRange(alarm.StartDate.UTC(), true)
						endDate := GetDateRange(alarm.EndDate.UTC(), false)

						exist := crowd.From(&list).Where(func(x interface{}) interface{} {
							y := x.(ScadaData)
							isBefore := y.TimeStamp.UTC().Before(endDate.UTC()) || y.TimeStamp.UTC().Equal(endDate.UTC())
							isAfter := y.TimeStamp.UTC().After(startDate.UTC()) || y.TimeStamp.UTC().Equal(startDate.UTC())
							return isBefore && isAfter
						}).Exec().Result.Data().([]ScadaData)

						if len(exist) > 0 {
							for _, ex := range exist {
								datasDown = append(datasDown, []float64{ex.AvgWindSpeed, ex.Power})
							}
						}
					}
					if len(datasDown) > 0 {
						turbineDataAlarm.Set("data", datasDown)
					}
				}

				found := false

			out:
				for _, ds := range dataSeries {
					if ds.GetString("name") == down {
						var tmp [][]float64
						if ds.Get("data") != nil && turbineDataAlarm.Get("data") != nil {
							tmp = ds.Get("data").([][]float64)
							tmp = append(tmp, turbineDataAlarm.Get("data").([][]float64)...)
						} else if turbineDataAlarm.Get("data") != nil {
							tmp = turbineDataAlarm.Get("data").([][]float64)
						}

						if tmp != nil {
							ds.Set("data", tmp)
						}

						found = true
						break out
					}
				}

				if !found {
					dataSeries = append(dataSeries, turbineDataAlarm)
				}
				idx++
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

func (m *AnalyticPowerCurveController) GetDetails(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	var (
		list       []ScadaData
		dataSeries []tk.M
	)

	p := new(PayloadAnalytic)
	e := k.GetPayload(&p)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	tStart, tEnd, e := helper.GetStartEndDate(k, p.Period, p.DateStart, p.DateEnd)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}
	turbine := p.Turbine
	project := p.Project
	colors := p.Color

	pcData, e := getPCData(project)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	dataSeries = append(dataSeries, pcData)

	if len(turbine) == 1 {
		var filter []*dbox.Filter
		filter = append(filter, dbox.Gte("dateinfo.dateid", tStart))
		filter = append(filter, dbox.Lte("dateinfo.dateid", tEnd))

		filter = append(filter, dbox.Eq("turbine", turbine[0]))

		filter = append(filter, dbox.Eq("projectname", project))

		csr, e := DB().Connection.NewQuery().From(new(ScadaData).TableName()).Where(dbox.And(filter...)).Cursor(nil)
		if e != nil {
			return helper.CreateResult(false, nil, e.Error())
		}

		e = csr.Fetch(&list, 0, false)

		if e != nil {
			return helper.CreateResult(false, nil, e.Error())
		}

		defer csr.Close()

		turbineData := tk.M{}
		turbineData.Set("name", turbine[0])
		turbineData.Set("type", "scatter")
		turbineData.Set("markers", tk.M{"size": 1})
		turbineData.Set("color", colors[0])

		var datas [][]float64

		for _, val := range list {
			datas = append(datas, []float64{val.AvgWindSpeed, val.Power})
		}

		turbineData.Set("data", datas)
		dataSeries = append(dataSeries, turbineData)
	}

	data := struct {
		Data []tk.M
	}{
		Data: dataSeries,
	}

	return helper.CreateResult(true, data, "success")
}

func getPCData(project string) (pcData tk.M, e error) {
	powerCurve := []PowerCurveModel{}

	csr, e := DB().Connection.NewQuery().From(new(PowerCurveModel).TableName()).Where(dbox.Eq("model", project)).Order("windspeed").Cursor(nil)
	if e != nil {
		return
	}

	e = csr.Fetch(&powerCurve, 0, false)
	defer csr.Close()

	if e != nil {
		return
	}

	var datas [][]float64

	for _, val := range powerCurve {
		datas = append(datas, []float64{val.WindSpeed, val.Power1})
	}

	pcData = tk.M{
		"name":     "Power Curve",
		"type":     "scatterLine",
		"dashType": "longDash",
		"style":    "smooth",
		"color":    "#ff880e",
		"markers":  tk.M{"visible": false},
		"width":    3,
	}

	if len(datas) > 0 {
		pcData.Set("data", datas)
	}

	return
}

func (m *AnalyticPowerCurveController) GetDownList(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	result := []tk.M{}
	for idx, dw := range helper.DownTypes {
		down := dw.GetString("down")
		label := dw.GetString("label")
		res := tk.M{
			"down":  down,
			"label": label,
			"color": downColor[idx],
		}
		result = append(result, res)
		idx++
	}

	return helper.CreateResult(true, result, "success")
}
