package controller

import (
	. "eaciit/wfdemo-git/library/core"
	. "eaciit/wfdemo-git/library/helper"
	. "eaciit/wfdemo-git/library/models"
	"eaciit/wfdemo-git/web/helper"
	"strings"

	// "fmt"
	"sort"
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
	project := ""
	if p.Project != "" {
		anProject := strings.Split(p.Project, "(")
		project = strings.TrimRight(anProject[0], " ")
	}
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
	project := ""
	if p.Project != "" {
		anProject := strings.Split(p.Project, "(")
		project = strings.TrimRight(anProject[0], " ")
	}
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
	pipes = append(pipes, tk.M{"$group": tk.M{"_id": "$denadjwindspeed", "production": tk.M{"$avg": "$power"}, "totaldata": tk.M{"$sum": 1}}})
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
		pipes        []tk.M
		filter       []*dbox.Filter
		list         []tk.M
		dataSeries   []tk.M
		sortTurbines []string
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
	project := ""
	if p.Project != "" {
		anProject := strings.Split(p.Project, "(")
		project = strings.TrimRight(anProject[0], " ")
	}
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
		colValue = "$power"
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

	dVal := (tk.ToFloat64(tk.ToInt(DeviationVal, tk.RoundingAuto), 2, tk.RoundingUp) / 100.0)
	selArr := 1

	filter = nil
	filter = append(filter, dbox.Ne("_id", ""))
	filter = append(filter, dbox.Gte("dateinfo.dateid", tStart))
	filter = append(filter, dbox.Lte("dateinfo.dateid", tEnd))
	filter = append(filter, dbox.Ne("turbine", ""))
	filter = append(filter, dbox.Gt("power", 0))
	if !IsDeviation {
		filter = append(filter, dbox.Gte(colDeviation, dVal))
	}
	if isClean {
		filter = append(filter, dbox.Eq("oktime", 600))
	}

	csr, e := DB().Connection.NewQuery().
		From(new(ScadaData).TableName()).
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

	for _, turX := range turbine {
		sortTurbines = append(sortTurbines, turX.(string))
	}
	sort.Strings(sortTurbines)

	for _, turbineX := range sortTurbines {

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
		turbineData.Set("idxseries", selArr)

		for _, val := range exist {
			idD := val.Get("_id").(tk.M)

			datas = append(datas, []float64{idD.GetFloat64("colId"), val.GetFloat64("production")}) //tk.Div(val.GetFloat64("production"), val.GetFloat64("totaldata"))
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

func (m *AnalyticPowerCurveController) GetListPowerCurveMonthly(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	var (
		pipes      []tk.M
		list       []tk.M
		dataSeries []tk.M
	)

	p := new(PayloadAnalyticPC)
	e := k.GetPayload(&p)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	project := ""
	if p.Project != "" {
		anProject := strings.Split(p.Project, "(")
		project = strings.TrimRight(anProject[0], " ")
	}

	now := time.Now()
	last := time.Now().AddDate(0, -12, 0)

	tStart, _ := time.Parse("20060102", last.Format("200601")+"01")
	tEnd, _ := time.Parse("20060102", now.Format("200601")+"01")

	colId := "$wsavgforpc"
	colValue := "$power"

	pcData, e := getPCData(project)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}
	match := []tk.M{}
	match = append(match, tk.M{"_id": tk.M{"$ne": ""}})
	match = append(match, tk.M{"dateinfo.dateid": tk.M{"$gte": tStart}})
	match = append(match, tk.M{"dateinfo.dateid": tk.M{"$lt": tEnd}})
	match = append(match, tk.M{"turbine": tk.M{"$ne": ""}})
	match = append(match, tk.M{"power": tk.M{"$gt": 0}})
	match = append(match, tk.M{"oktime": 600})
	if project != "" {
		match = append(match, tk.M{"projectname": project})
	}
	if len(p.Turbine) > 0 {
		match = append(match, tk.M{"turbine": tk.M{"$in": p.Turbine}})
	}

	pipes = append(pipes, tk.M{"$match": tk.M{"$and": match}})

	pipes = append(pipes, tk.M{"$group": tk.M{
		"_id": tk.M{
			"Turbine":   "$turbine",
			"monthid":   "$dateinfo.monthid",
			"monthdesc": "$dateinfo.monthdesc",
			"colId":     colId,
		},
		"production": tk.M{"$avg": colValue},
	}})
	pipes = append(pipes, tk.M{"$sort": tk.M{
		"_id.Turbine": 1,
		"_id.monthid": 1,
		"_id.colId":   1,
	}})
	csr, e := DB().Connection.NewQuery().
		From(new(ScadaData).TableName()).
		Command("pipe", pipes).
		Cursor(nil)

	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}
	e = csr.Fetch(&list, 0, false)
	defer csr.Close()
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	ids := tk.M{}
	var datas [][]float64
	results := []tk.M{}
	monthIndex := tk.M{}
	monthList := tk.M{}
	splitMonth := []string{}
	simpleMonth := ""
	sortMonth := []int{}
	turbine := p.Turbine
	sortTurbines := []string{}

	for _, listVal := range list {
		ids, _ = tk.ToM(listVal["_id"])
		if len(p.Turbine) == 0 {
			exist := false
			for _, val := range turbine {
				if ids["Turbine"] == val {
					exist = true
				}
			}
			if exist == false {
				turbine = append(turbine, ids["Turbine"])
			}
		}
		monthList.Set(tk.ToString(ids.GetInt("monthid")), ids.GetString("monthdesc"))
	}
	for key := range monthList {
		sortMonth = append(sortMonth, tk.ToInt(key, tk.RoundingAuto))
	}
	sort.Ints(sortMonth)
	for _, turX := range turbine {
		sortTurbines = append(sortTurbines, turX.(string))
	}
	sort.Strings(sortTurbines)
	selArr := 0
	dataSeries = append(dataSeries, pcData)

	for _, turbineX := range sortTurbines {
		selArr = 0
		for _, monthX := range sortMonth {
			monthExist := crowd.From(&list).Where(func(x interface{}) interface{} {
				y := x.(tk.M)
				id := y.Get("_id").(tk.M)

				return id.GetInt("monthid") == monthX && id.GetString("Turbine") == turbineX
			}).Exec().Result.Data().([]tk.M)

			datas = [][]float64{}
			selArr++
			splitMonth = strings.Split(monthList.GetString(tk.ToString(monthX)), " ")
			simpleMonth = splitMonth[0][0:3] + " " + splitMonth[1][2:4] /*it will be jan 16, feb 16, and so on*/

			monthData := tk.M{}
			monthData.Set("name", turbineX)
			monthData.Set("type", "scatterLine")
			monthData.Set("style", "smooth")
			monthData.Set("dashType", "solid")
			monthData.Set("markers", tk.M{"visible": false})
			monthData.Set("width", 2)
			monthData.Set("color", colorField[selArr])
			monthData.Set("idxseries", selArr)

			monthIndex.Set(tk.ToString(selArr), simpleMonth)

			for _, val := range monthExist {
				idD := val.Get("_id").(tk.M)

				datas = append(datas, []float64{idD.GetFloat64("colId"), val.GetFloat64("production")})
			}

			if len(datas) > 0 {
				monthData.Set("data", datas)
			}

			dataSeries = append(dataSeries, monthData)
		}
		turbineData := tk.M{
			"Name": turbineX, /*for chart name*/
			"Data": dataSeries,
		}
		results = append(results, turbineData)
		dataSeries = []tk.M{}                   /*clear variable for next data*/
		dataSeries = append(dataSeries, pcData) /*always append expected value at beginning*/
	}

	sortedIndex := []int{}
	for key := range monthIndex {
		sortedIndex = append(sortedIndex, tk.ToInt(key, tk.RoundingAuto))
	}
	sort.Ints(sortedIndex)

	categoryList := []tk.M{}
	catList := tk.M{"category": "Power Curve", "color": "#ea5b19"}
	categoryList = append(categoryList, catList)

	for _, idx := range sortedIndex {
		catList = tk.M{"category": monthIndex.GetString(tk.ToString(idx)), "color": colorField[idx]}
		categoryList = append(categoryList, catList)
	}

	data := struct {
		Data     []tk.M
		Category []tk.M
	}{
		Data:     results,
		Category: categoryList,
	}

	return helper.CreateResult(true, data, "success")
}

func (m *AnalyticPowerCurveController) GetListPowerCurveComparison(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	var (
		pipes      []tk.M
		filter     []*dbox.Filter
		list       []tk.M
		dataSeries []tk.M
		// sortTurbines []string
	)

	p := new(PayloadPCComparison)
	e := k.GetPayload(&p)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	PC1tStart, PC1tEnd, e := helper.GetStartEndDate(k, p.PC1Period, p.PC1DateStart, p.PC1DateEnd)
	PC2tStart, PC2tEnd, e := helper.GetStartEndDate(k, p.PC2Period, p.PC2DateStart, p.PC2DateEnd)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}
	PC1turbine := p.PC1Turbine
	PC1project := ""
	if p.PC1Project != "" {
		anProject := strings.Split(p.PC1Project, "(")
		PC1project = strings.TrimRight(anProject[0], " ")
	}

	PC2turbine := p.PC2Turbine
	// PC2project := ""
	// if p.PC2Project != "" {
	// 	anProject := strings.Split(p.PC2Project, "(")
	// 	PC2project = strings.TrimRight(anProject[0], " ")
	// }

	colId := "$wsavgforpc"
	colValue := "$power"

	PC1Data, e := getPCData(PC1project)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	// PC2Data, e := getPCData(PC2project)
	// if e != nil {
	// 	return helper.CreateResult(false, nil, e.Error())
	// }

	dataSeries = append(dataSeries, PC1Data)
	// dataSeries = append(dataSeries, PC2Data)

	pipes = append(pipes, tk.M{"$group": tk.M{"_id": colId, "production": tk.M{"$avg": colValue}, "totaldata": tk.M{"$sum": 1}}})
	pipes = append(pipes, tk.M{"$sort": tk.M{"_id": 1}})

	filter = nil
	filter = append(filter, dbox.Ne("_id", ""))
	filter = append(filter, dbox.Gte("dateinfo.dateid", PC1tStart))
	filter = append(filter, dbox.Lte("dateinfo.dateid", PC1tEnd))
	filter = append(filter, dbox.Eq("turbine", PC1turbine))
	filter = append(filter, dbox.Gt("power", 0))

	csr, e := DB().Connection.NewQuery().
		From(new(ScadaData).TableName()).
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
	// turbineData.Set("name", PC1turbine)
	turbineData.Set("name", "Turbine 1")
	turbineData.Set("type", "scatterLine")
	turbineData.Set("style", "smooth")
	turbineData.Set("dashType", "solid")
	turbineData.Set("markers", tk.M{"visible": false})
	turbineData.Set("width", 2)
	turbineData.Set("color", colorField[1])

	for _, val := range list {

		datas = append(datas, []float64{val.GetFloat64("_id"), val.GetFloat64("production")})
	}

	if len(datas) > 0 {
		turbineData.Set("data", datas)
	}

	dataSeries = append(dataSeries, turbineData)

	filter = nil
	filter = append(filter, dbox.Ne("_id", ""))
	filter = append(filter, dbox.Gte("dateinfo.dateid", PC2tStart))
	filter = append(filter, dbox.Lte("dateinfo.dateid", PC2tEnd))
	filter = append(filter, dbox.Eq("turbine", PC2turbine))
	filter = append(filter, dbox.Gt("power", 0))

	csr, e = DB().Connection.NewQuery().
		From(new(ScadaData).TableName()).
		Command("pipe", pipes).
		Where(dbox.And(filter...)).
		Cursor(nil)

	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}
	e = csr.Fetch(&list, 0, false)
	defer csr.Close()

	var datasC2 [][]float64
	turbineData = tk.M{}
	// turbineData.Set("name", PC2turbine)
	turbineData.Set("name", "Turbine 2")
	turbineData.Set("type", "scatterLine")
	turbineData.Set("style", "smooth")
	turbineData.Set("dashType", "solid")
	turbineData.Set("markers", tk.M{"visible": false})
	turbineData.Set("width", 2)
	turbineData.Set("color", colorField[6])

	for _, val := range list {

		datasC2 = append(datasC2, []float64{val.GetFloat64("_id"), val.GetFloat64("production")})
	}

	if len(datasC2) > 0 {
		turbineData.Set("data", datasC2)
	}

	dataSeries = append(dataSeries, turbineData)

	data := struct {
		Data []tk.M
	}{
		Data: dataSeries,
	}

	return helper.CreateResult(true, data, "success")
}

func setScatterData(name, xField, yField, color, yAxis string, marker tk.M, data []tk.M) tk.M {
	return tk.M{
		"name":       name,
		"xField":     xField,
		"yField":     yField,
		"colorField": "valueColor",
		"color":      color,
		"type":       "scatter",
		"markers":    marker,
		"yAxis":      yAxis,
		"data":       data,
	}
}

func (m *AnalyticPowerCurveController) GetPowerCurveScatter(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	type PayloadScatter struct {
		Period      string
		DateStart   time.Time
		DateEnd     time.Time
		Turbine     string
		Project     string
		ScatterType string
	}
	type ScadaMini struct {
		Power, AvgWindSpeed               float64
		NacelleTemperature, WindDirection float64
	}

	type ScadaOEMMini struct {
		AI_intern_WindSpeed, AI_intern_Pitchangle1   float64
		AI_intern_Pitchangle2, AI_intern_Pitchangle3 float64
	}
	var (
		list       []ScadaMini
		dataSeries []tk.M
	)

	p := new(PayloadScatter)
	e := k.GetPayload(&p)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	tStart, tEnd, e := helper.GetStartEndDate(k, p.Period, p.DateStart, p.DateEnd)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}
	turbine := p.Turbine
	project := ""
	if p.Project != "" {
		anProject := strings.Split(p.Project, "(")
		project = strings.TrimRight(anProject[0], " ")
	}
	pcData, e := getPCData(project)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}
	dataSeries = append(dataSeries, pcData)

	var filter []*dbox.Filter
	filter = []*dbox.Filter{}
	filter = append(filter, dbox.Ne("_id", ""))
	filter = append(filter, dbox.Gte("timestamp", tStart))
	filter = append(filter, dbox.Lte("timestamp", tEnd))
	filter = append(filter, dbox.Eq("turbine", turbine))
	filter = append(filter, dbox.Eq("projectname", project))
	filter = append(filter, dbox.Eq("oktime", 600))
	filter = append(filter, dbox.Gt("avgwindspeed", 0))

	csr, e := DB().Connection.NewQuery().
		From(new(ScadaData).TableName()).
		Where(dbox.And(filter...)).
		Take(10000).
		Cursor(nil)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}
	var _list ScadaMini
	for {
		e = csr.Fetch(&_list, 1, false)
		if e != nil {
			break
		}
		list = append(list, _list)
	}

	defer csr.Close()

	datas := tk.M{}
	arrDatas := []tk.M{}
	tempData := tk.M{}
	tempDatas := []tk.M{}
	deviationData := tk.M{}
	deviationDatas := []tk.M{}
	for _, val := range list {
		datas = tk.M{}
		tempData = tk.M{}
		deviationData = tk.M{}

		if val.Power > 0 {
			datas.Set("WindSpeed", val.AvgWindSpeed)
			datas.Set("Power", val.Power)
			datas.Set("valueColor", colorField[1])

			arrDatas = append(arrDatas, datas)
		}

		if p.ScatterType != "pitch" { /*processing NON pitch data*/
			switch p.ScatterType {
			case "temp":
				if val.NacelleTemperature > 0 {
					tempData.Set("WindSpeed", val.AvgWindSpeed)
					tempData.Set("Temperature", val.NacelleTemperature)
					tempData.Set("valueColor", colorField[2])

					tempDatas = append(tempDatas, tempData)
				}
			case "deviation":
				deviationData.Set("WindSpeed", val.AvgWindSpeed)
				deviationData.Set("Deviation", val.WindDirection)
				deviationData.Set("valueColor", colorField[2])

				deviationDatas = append(deviationDatas, deviationData)
			}
		}
	}
	turbineData := setScatterData("Power", "WindSpeed", "Power", colorField[1], "powerAxis", tk.M{"size": 2}, arrDatas)
	dataSeries = append(dataSeries, turbineData)

	/*================== SCADA OEM PART ==================*/
	if p.ScatterType == "pitch" {
		var filterOEM []*dbox.Filter
		filterOEM = []*dbox.Filter{}
		filterOEM = append(filterOEM, dbox.Ne("_id", ""))
		filterOEM = append(filterOEM, dbox.Gte("timestamp", tStart))
		filterOEM = append(filterOEM, dbox.Lte("timestamp", tEnd))
		filterOEM = append(filterOEM, dbox.Eq("turbine", turbine))
		filterOEM = append(filterOEM, dbox.Eq("projectname", project))
		filterOEM = append(filterOEM, dbox.Eq("mttr", 600.0))
		filterOEM = append(filterOEM, dbox.Gt("ai_intern_windspeed", 0.0))

		csrOEM, e := DB().Connection.NewQuery().
			From(new(ScadaDataOEM).TableName()).
			Where(dbox.And(filterOEM...)).
			Take(10000).
			Cursor(nil)
		if e != nil {
			return helper.CreateResult(false, nil, e.Error())
		}
		var _listOEM ScadaOEMMini
		var listScadaOEM []ScadaOEMMini
		for {
			e = csrOEM.Fetch(&_listOEM, 1, false)
			if e != nil {
				break
			}
			listScadaOEM = append(listScadaOEM, _listOEM)
		}
		defer csrOEM.Close()

		pitchData := tk.M{}
		pitchDatas := []tk.M{}
		count := 0.0
		pitchAngle := 0.0
		for _, val := range listScadaOEM { /*processing pitch data*/
			pitchData = tk.M{}
			pitchAngle = 0.0
			if val.AI_intern_Pitchangle1 >= -10.0 && val.AI_intern_Pitchangle1 <= 120.0 {
				pitchAngle = val.AI_intern_Pitchangle1
				count++
			} else if val.AI_intern_Pitchangle2 >= -10.0 && val.AI_intern_Pitchangle2 <= 120.0 {
				pitchAngle += val.AI_intern_Pitchangle2
				count++
			} else if val.AI_intern_Pitchangle3 >= -10.0 && val.AI_intern_Pitchangle3 <= 120.0 {
				pitchAngle += val.AI_intern_Pitchangle3
				count++
			}
			pitchAngle /= count /*average pitch angle*/

			if pitchAngle != 0.0 || count > 0 {
				pitchData.Set("WindSpeed", val.AI_intern_WindSpeed)
				pitchData.Set("Pitch", pitchAngle)
				pitchData.Set("valueColor", colorField[2])
				pitchDatas = append(pitchDatas, pitchData)
			}
		}
		seriesData := setScatterData("Pitch", "WindSpeed", "Pitch", colorField[2], "pitchAxis", tk.M{"size": 2}, pitchDatas)
		dataSeries = append(dataSeries, seriesData)
	}
	/*================== END OF SCADA OEM PART ==================*/
	switch p.ScatterType {
	case "temp":
		/*set data series*/
		seriesData := setScatterData("Temperature", "WindSpeed", "Temperature", colorField[2], "tempAxis", tk.M{"size": 2}, tempDatas)
		dataSeries = append(dataSeries, seriesData)
	case "deviation":
		seriesData := setScatterData("Deviation", "WindSpeed", "Deviation", colorField[2], "deviationAxis", tk.M{"size": 2}, deviationDatas)
		dataSeries = append(dataSeries, seriesData)
	}

	data := struct {
		Data []tk.M
	}{
		Data: dataSeries,
	}

	return helper.CreateResult(true, data, "success")
}

func (m *AnalyticPowerCurveController) GetPCScatterAnalysis(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	type PayloadScatter struct {
		Period        string
		DateStart     time.Time
		DateEnd       time.Time
		Turbine       string
		Project       string
		ScatterType   string
		LessValue     int
		GreaterValue  int
		LessColor     string
		GreaterColor  string
		GreaterMarker string
		LessMarker    string
	}

	type ScadaMini struct {
		Power, AvgWindSpeed float64
		WindDirection       float64
	}

	type ScadaOEMMini struct {
		AI_intern_WindSpeed, AI_intern_Pitchangle1   float64
		AI_intern_Pitchangle2, AI_intern_Pitchangle3 float64
		AI_intern_ActivPower                         float64
	}

	var (
		list       []ScadaMini
		dataSeries []tk.M
	)

	p := new(PayloadScatter)
	e := k.GetPayload(&p)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	tStart, tEnd, e := helper.GetStartEndDate(k, p.Period, p.DateStart, p.DateEnd)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}
	turbine := p.Turbine
	project := ""
	if p.Project != "" {
		anProject := strings.Split(p.Project, "(")
		project = strings.TrimRight(anProject[0], " ")
	}
	pcData, e := getPCData(project)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}
	dataSeries = append(dataSeries, pcData)

	/*=======POWER LINE QUERY =========*/
	pipes := []tk.M{}
	pipes = append(pipes, tk.M{
		"$group": tk.M{
			"_id":        "$wsavgforpc",
			"production": tk.M{"$avg": "$power"},
		},
	})
	pipes = append(pipes, tk.M{"$sort": tk.M{"_id": 1}})

	var filter []*dbox.Filter
	filter = []*dbox.Filter{}
	filter = append(filter, dbox.Ne("_id", ""))
	filter = append(filter, dbox.Gte("timestamp", tStart))
	filter = append(filter, dbox.Lte("timestamp", tEnd))
	filter = append(filter, dbox.Eq("turbine", turbine))
	filter = append(filter, dbox.Eq("projectname", project))
	filter = append(filter, dbox.Gt("power", 0))
	filter = append(filter, dbox.Eq("oktime", 600))

	csrPower, e := DB().Connection.NewQuery().
		From(new(ScadaData).TableName()).
		Command("pipe", pipes).
		Where(dbox.And(filter...)).
		Cursor(nil)

	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}
	powerData := []tk.M{}
	e = csrPower.Fetch(&powerData, 0, false)
	defer csrPower.Close()

	var datas [][]float64
	turbineData := tk.M{}
	turbineData.Set("name", turbine)
	turbineData.Set("type", "scatterLine")
	turbineData.Set("style", "smooth")
	turbineData.Set("dashType", "solid")
	turbineData.Set("markers", tk.M{"visible": false})
	turbineData.Set("width", 2)
	turbineData.Set("color", colorField[1])
	turbineData.Set("idxseries", 1)

	for _, val := range powerData {
		datas = append(datas, []float64{val.GetFloat64("_id"), val.GetFloat64("production")}) //tk.Div(val.GetFloat64("production"), val.GetFloat64("totaldata"))
	}

	if len(datas) > 0 {
		turbineData.Set("data", datas)
	}

	dataSeries = append(dataSeries, turbineData)

	/*===== END OF POWER LINE =======*/
	// filter is same with power filter

	csr, e := DB().Connection.NewQuery().
		From(new(ScadaData).TableName()).
		Where(dbox.And(filter...)).
		Take(10000).
		Cursor(nil)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}
	var _list ScadaMini
	for {
		e = csr.Fetch(&_list, 1, false)
		if e != nil {
			break
		}
		list = append(list, _list)
	}

	defer csr.Close()

	scatterData := tk.M{}
	scatterDatas1 := []tk.M{}
	scatterDatas2 := []tk.M{}
	lessDev := tk.ToFloat64(p.LessValue, 2, tk.RoundingAuto) / 100.0
	greatDev := tk.ToFloat64(p.GreaterValue, 2, tk.RoundingAuto) / 100.0

	if p.ScatterType != "pitch" { /*processing data non pitch*/
		for _, val := range list {
			scatterData = tk.M{}
			scatterData.Set("WindSpeed", val.AvgWindSpeed)
			scatterData.Set("Power", val.Power)
			switch p.ScatterType {
			case "deviation":
				if val.WindDirection < lessDev {
					scatterDatas1 = append(scatterDatas1, scatterData)
				}
				if val.WindDirection > greatDev {
					scatterDatas2 = append(scatterDatas2, scatterData)
				}
			}
		}
	} else {
		/*================== SCADA OEM PART ==================*/
		var filterOEM []*dbox.Filter
		filterOEM = []*dbox.Filter{}
		filterOEM = append(filterOEM, dbox.Ne("_id", ""))
		filterOEM = append(filterOEM, dbox.Gte("timestamp", tStart))
		filterOEM = append(filterOEM, dbox.Lte("timestamp", tEnd))
		filterOEM = append(filterOEM, dbox.Eq("turbine", turbine))
		filterOEM = append(filterOEM, dbox.Eq("projectname", project))
		filterOEM = append(filterOEM, dbox.Eq("mttr", 600.0))
		filterOEM = append(filterOEM, dbox.Gt("ai_intern_activpower", 0.0))

		csrOEM, e := DB().Connection.NewQuery().
			From(new(ScadaDataOEM).TableName()).
			Where(dbox.And(filterOEM...)).
			Take(10000).
			Cursor(nil)
		if e != nil {
			return helper.CreateResult(false, nil, e.Error())
		}
		var _listOEM ScadaOEMMini
		var listScadaOEM []ScadaOEMMini
		for {
			e = csrOEM.Fetch(&_listOEM, 1, false)
			if e != nil {
				break
			}
			listScadaOEM = append(listScadaOEM, _listOEM)
		}
		defer csrOEM.Close()

		count := 0.0
		pitchAngle := 0.0

		for _, val := range listScadaOEM { /*processing pitch data*/
			pitchAngle = 0.0
			scatterData = tk.M{}
			scatterData.Set("WindSpeed", val.AI_intern_WindSpeed)
			scatterData.Set("Power", val.AI_intern_ActivPower)

			if val.AI_intern_Pitchangle1 >= -10.0 && val.AI_intern_Pitchangle1 <= 120.0 {
				pitchAngle = val.AI_intern_Pitchangle1
				count++
			} else if val.AI_intern_Pitchangle2 >= -10.0 && val.AI_intern_Pitchangle2 <= 120.0 {
				pitchAngle += val.AI_intern_Pitchangle2
				count++
			} else if val.AI_intern_Pitchangle3 >= -10.0 && val.AI_intern_Pitchangle3 <= 120.0 {
				pitchAngle += val.AI_intern_Pitchangle3
				count++
			}
			pitchAngle /= count /*average pitch angle*/

			if pitchAngle != 0.0 || count > 0 {
				if pitchAngle < lessDev {
					scatterDatas1 = append(scatterDatas1, scatterData)
				}
				if pitchAngle > greatDev {
					scatterDatas2 = append(scatterDatas2, scatterData)
				}
			}
		}
	}
	/*================== END OF SCADA OEM PART ==================*/

	switch p.ScatterType {
	case "deviation":
		seriesData1 := setScatterData("Nacelle Deviation < "+tk.ToString(p.LessValue), "WindSpeed", "Power", p.LessColor, "powerAxis", tk.M{"size": 2, "type": p.LessMarker, "background": p.LessColor}, scatterDatas1)
		seriesData1.Unset("colorField")
		dataSeries = append(dataSeries, seriesData1)
		seriesData2 := setScatterData("Nacelle Deviation > "+tk.ToString(p.GreaterValue), "WindSpeed", "Power", p.GreaterColor, "powerAxis", tk.M{"size": 2, "type": p.GreaterMarker, "background": p.GreaterColor}, scatterDatas2)
		seriesData2.Unset("colorField")
		dataSeries = append(dataSeries, seriesData2)
	case "pitch":
		seriesData1 := setScatterData("Pitch Angle < "+tk.ToString(p.LessValue), "WindSpeed", "Power", p.LessColor, "powerAxis", tk.M{"size": 2, "type": p.LessMarker, "background": p.LessColor}, scatterDatas1)
		seriesData1.Unset("colorField")
		dataSeries = append(dataSeries, seriesData1)
		seriesData2 := setScatterData("Pitch Angle > "+tk.ToString(p.GreaterValue), "WindSpeed", "Power", p.GreaterColor, "powerAxis", tk.M{"size": 2, "type": p.GreaterMarker, "background": p.GreaterColor}, scatterDatas2)
		seriesData2.Unset("colorField")
		dataSeries = append(dataSeries, seriesData2)
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
	project := ""
	if p.Project != "" {
		anProject := strings.Split(p.Project, "(")
		project = strings.TrimRight(anProject[0], " ")
	}
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
		// turbineData.Set("markers", tk.M{
		// 			"size":       10,
		// 			"type":       "triangle",
		// 			"border" : tk.M{
		// 				"width" : 2,
		// 			"color" : "red",
		// 			},
		// 		})
		turbineData.Set("markers", tk.M{"size": 2})

		datas := tk.M{}
		arrDatas := []tk.M{}
		for _, val := range list {
			datas = tk.M{}

			switch viewSession {
			case "density":
				if val.DenWindSpeed > 0 && val.Power > 0 {
					datas.Set("WindSpeed", val.DenWindSpeed)
					datas.Set("Power", val.Power)

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
					"type":       "triangle",
					"background": downColor[idx],
					// "border" : tk.M{
					// 	"width" : 2,
					// 	"color" : "red",
					// },
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
	project := ""
	if p.Project != "" {
		anProject := strings.Split(p.Project, "(")
		project = strings.TrimRight(anProject[0], " ")
	}
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
		"name":      "Power Curve",
		"idxseries": 0,
		"type":      "scatterLine",
		"dashType":  "longDash",
		"style":     "smooth",
		"color":     "#ea5b19",
		"markers":   tk.M{"visible": false},
		"width":     3,
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
