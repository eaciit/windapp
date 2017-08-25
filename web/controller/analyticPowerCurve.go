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
	colorField = [...]string{"#ff9933", "#21c4af", "#ff7663", "#ffb74f", "#a2df53", "#1c9ec4", "#ff63a5", "#f44336", "#69d2e7", "#8877A9", "#9A12B3", "#26C281", "#E7505A", "#C49F47", "#ff5597", "#c3260c", "#d4735e", "#ff2ad7", "#34ac8b", "#11b2eb", "#004c79", "#ff0037", "#507ca3", "#ff6565", "#ffd664", "#72aaff", "#795548", "#383271", "#6a4795", "#bec554", "#ab5919", "#f5b1e1", "#7b3416", "#002fef", "#8d731b", "#1f8805", "#ff9900", "#9C27B0", "#6c7d8a", "#d73c1c", "#5be7a0", "#da02d4", "#afa56e", "#7e32cb", "#a2eaf7", "#9cb8f4", "#9E9E9E", "#065806", "#044082", "#18937d", "#2c787a", "#a57c0c", "#234341", "#1aae7a", "#7ac610", "#736f5f", "#4e741e", "#68349d", "#1df3b6", "#e02b09", "#d9cfab", "#6e4e52", "#f31880", "#7978ec", "#f5ace8", "#3db6ae", "#5e06b0", "#16d0b9", "#a25a5b", "#1e603a", "#4b0981", "#62975f", "#1c8f2f", "#b0c80c", "#642794", "#e2060d", "#2125f0"}
	// colorField            = [...]string{"#cc2a35","#87c5da","#115b74","#e18876","#95204d","#c5a5ca","#7d277e","#ffd145","#145b9b","#dab5cb","#dab5cb","#007ca7", "#26C281", "#E7505A", "#C49F47", "#ff5597", "#c3260c", "#d4735e", "#ff2ad7", "#34ac8b", "#11b2eb", "#f35838", "#ff0037", "#507ca3", "#ff6565", "#ffd664", "#72aaff", "#795548"}
	// colorField  		  = [...]string{"#87c5da","#cc2a35", "#d66b76", "#5d1b62", "#f1c175","#95204c","#8f4bc5","#7d287d","#00818e","#c8c8c8","#546698","#66c99a","#f3d752","#20adb8","#333d6b","#d077b1","#aab664","#01a278","#c1d41a","#807063","#ff5975","#01a3d4","#ca9d08","#026e51","#4c653f"}
	// colorFieldDegradation = [...]string{"#ffcf9e", "#a6e7df", "#ffc8c0", "#ffe2b8", "#d9f2ba", "#a4d8e7", "#ffc0db", "#fab3ae", "#efa5a2", "#cfc8dc", "#d6a0e0", "#a8e6cc", "#f5b9bd", "#e7d8b5", "#ffbbd5", "#e7a89d", "#edc7be", "#ffa9ef", "#adddd0", "#9fe0f7", "#fabcaf", "#ff99af", "#b9cada", "#ffc1c1", "#ffeec1", "#c6ddff", "#c9bbb5"}
	colorFieldDegradation = [...]string{"#FFD6AD", "#A6E7DF", "#FFC8C0", "#FFE2B8", "#D9F2BA", "#A4D8E7", "#FFC0DB", "#FAB3AE", "#C3EDF5", "#CFC8DC", "#D6A0E0", "#A8E6CC", "#F5B9BD", "#E7D8B5", "#FFBBD5", "#E7A89D", "#EDC7BE", "#FFA9EF", "#ADDDD0", "#9FE0F7", "#99B7C9", "#FF99AF", "#B9CADA", "#FFC1C1", "#FFEEC1", "#C6DDFF", "#C9BBB5", "#AFADC6", "#C3B5D4", "#E5E7BA", "#DDBCA3", "#FBDFF3", "#CAADA1", "#99ABF8", "#D1C7A3", "#A5CF9B", "#FFD699", "#D7A8DF", "#C4CBD0", "#EFB1A4", "#BDF5D9", "#F099ED", "#DFDBC5", "#CBADEA", "#D9F6FB", "#D7E2FA", "#D8D8D8", "#9BBC9B", "#9AB2CD", "#A2D3CB", "#AAC9C9", "#DBCA9D", "#A7B3B3", "#A3DEC9", "#C9E89F", "#C7C5BF", "#B8C7A5", "#C2ADD7", "#A4FAE1", "#F2AA9C", "#EFEBDD", "#C5B8B9", "#FAA2CC", "#C9C9F7", "#FBDDF5", "#B1E1DE", "#BE9BDF", "#A1ECE3", "#D9BDBD", "#A5BFB0", "#B79CCC", "#C0D5BF", "#A4D2AB", "#DFE99D", "#C1A8D4", "#F39B9E", "#A6A7F9"}
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
	filter = nil
	project := p.Project
	if project != "" {
		filter = append(filter, dbox.Eq("projectname", project))
	}
	IsDeviation := p.IsDeviation
	DeviationVal := p.DeviationVal
	DeviationOpr := tk.ToInt(p.DeviationOpr, tk.RoundingAuto)
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

	filter = append(filter, dbox.Ne("_id", ""))
	filter = append(filter, dbox.Gte("dateinfo.dateid", tStart))
	filter = append(filter, dbox.Lte("dateinfo.dateid", tEnd))
	filter = append(filter, dbox.Ne("turbine", ""))
	filter = append(filter, dbox.Gt("power", 0))
	filter = append(filter, dbox.Eq("available", 1))

	// modify by ams, 2017-08-11
	if IsDeviation {
		if DeviationOpr > 0 {
			filter = append(filter, dbox.Gte(colDeviation, dVal))
		} else {
			filter = append(filter, dbox.Lte(colDeviation, dVal))
		}
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

	selArr := 1
	turbineName := map[string]string{}
	if p.Project != "" {
		turbineName, _ = helper.GetTurbineNameList(p.Project)
	}
	for idx, turbineX := range sortTurbines {
		exist := crowd.From(&list).Where(func(x interface{}) interface{} {
			y := x.(tk.M)
			id := y.Get("_id").(tk.M)

			return id.GetString("Turbine") == turbineX
		}).Exec().Result.Data().([]tk.M)

		var datas [][]float64
		turbineData := tk.M{}
		turbineData.Set("name", turbineName[turbineX])
		turbineData.Set("turbineid", turbineX)
		turbineData.Set("type", "scatterLine")
		turbineData.Set("style", "smooth")
		turbineData.Set("dashType", "solid")
		turbineData.Set("markers", tk.M{"visible": false})
		turbineData.Set("width", 2)
		turbineData.Set("color", colorField[selArr])
		turbineData.Set("idxseries", idx+1)

		for _, val := range exist {
			idD := val.Get("_id").(tk.M)

			datas = append(datas, []float64{idD.GetFloat64("colId"), val.GetFloat64("production")}) //tk.Div(val.GetFloat64("production"), val.GetFloat64("totaldata"))
		}

		if len(datas) > 0 {
			turbineData.Set("data", datas)
		}

		dataSeries = append(dataSeries, turbineData)

		if selArr == len(colorField)-1 {
			selArr = 1
		} else {
			selArr++
		}

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

	project := p.Project

	now := time.Now()
	last := time.Now().AddDate(0, -12, 0)

	tStart, _ := time.Parse("20060102", last.Format("200601")+"01")
	tEnd, _ := time.Parse("20060102", now.Format("200601")+"01")

	//tk.Printf("Start : #%v\n", tStart)
	//tk.Printf("End : #%v\n", tEnd)

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
	match = append(match, tk.M{"available": 1})

	if project != "" {
		match = append(match, tk.M{"projectname": project})
	}
	if len(p.Turbine) > 0 {
		match = append(match, tk.M{"turbine": tk.M{"$in": p.Turbine}})
	}

	//tk.Printf("%#v\n", match)

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

	//tk.Printf("Sort Turbines : %#v\n", sortTurbines)

	selArr := 0
	dataSeries = append(dataSeries, pcData)

	turbineName := map[string]string{}
	if p.Project != "" {
		turbineName, _ = helper.GetTurbineNameList(p.Project)
	}
	//tk.Printf("Turbines : %#v\n", turbineName)
	turbineXid := ""
	for _, turbineX := range sortTurbines {
		selArr = 0
		turbineXid = turbineName[turbineX]
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
			monthData.Set("name", turbineXid)
			monthData.Set("turbineid", turbineX)
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
			"Name": turbineXid, /*for chart name*/
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
	catList := tk.M{"category": "Power Curve", "color": "#ff9933"}
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
	PC1project := p.PC1Project
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
	filter = append(filter, dbox.Eq("available", 1))

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
	filter = append(filter, dbox.Eq("available", 1))

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

func getScatterValue(list []tk.M, tipe, field string) (resWSvsPower []tk.M, resWSvsTipe []tk.M) {
	resWSvsPower = []tk.M{}
	resWSvsTipe = []tk.M{}
	dataWSvsPower := tk.M{}
	dataWSvsTipe := tk.M{}
	for _, val := range list {
		dataWSvsPower = tk.M{}
		dataWSvsTipe = tk.M{}

		dataWSvsPower.Set("WindSpeed", val.GetFloat64("avgwindspeed"))
		dataWSvsPower.Set("Power", val.GetFloat64("power"))
		dataWSvsPower.Set("valueColor", colorField[1])

		resWSvsPower = append(resWSvsPower, dataWSvsPower)

		dataWSvsTipe.Set("WindSpeed", val.GetFloat64("avgwindspeed"))
		dataWSvsTipe.Set(tipe, val.GetFloat64(field))
		dataWSvsTipe.Set("valueColor", colorField[2])

		resWSvsTipe = append(resWSvsTipe, dataWSvsTipe)
	}
	return
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
		Power, AvgWindSpeed, AvgBladeAngle   float64
		NacelleTemperature, NacelleDeviation float64
	}
	list := []tk.M{}
	var dataSeries []tk.M

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
	project := p.Project
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
	// filter = append(filter, dbox.Eq("oktime", 600))
	filter = append(filter, dbox.Gt("avgwindspeed", 0))
	filter = append(filter, dbox.Gt("power", 0))
	filter = append(filter, dbox.Eq("available", 1))

	csr, e := DB().Connection.NewQuery().
		Select("power", "avgwindspeed", "avgbladeangle", "nacelletemperature", "nacelledeviation", "ambienttemperature").
		From(new(ScadaData).TableName()).
		Where(dbox.And(filter...)).
		Take(10000).
		Cursor(nil)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}
	_list := tk.M{}
	for {
		_list = tk.M{}
		e = csr.Fetch(&_list, 1, false)
		if e != nil {
			break
		}
		list = append(list, _list)
	}

	defer csr.Close()

	turbineData := tk.M{}
	seriesData := tk.M{}
	resWSvsPower := []tk.M{}
	resWSvsTipe := []tk.M{}

	switch p.ScatterType {
	case "temp":
		resWSvsPower, resWSvsTipe = getScatterValue(list, "Temperature", "nacelletemperature")
		seriesData = setScatterData("Temperature", "WindSpeed", "Temperature", colorField[2], "tempAxis", tk.M{"size": 2}, resWSvsTipe)
	case "deviation":
		resWSvsPower, resWSvsTipe = getScatterValue(list, "Deviation", "nacelledeviation")
		seriesData = setScatterData("Deviation", "WindSpeed", "Deviation", colorField[2], "deviationAxis", tk.M{"size": 2}, resWSvsTipe)
	case "pitch":
		resWSvsPower, resWSvsTipe = getScatterValue(list, "Pitch", "avgbladeangle")
		seriesData = setScatterData("Pitch", "WindSpeed", "Pitch", colorField[2], "pitchAxis", tk.M{"size": 2}, resWSvsTipe)
	case "ambient":
		resWSvsPower, resWSvsTipe = getScatterValue(list, "Ambient", "ambienttemperature")
		seriesData = setScatterData("Ambient", "WindSpeed", "Ambient", colorField[2], "ambientAxis", tk.M{"size": 2}, resWSvsTipe)
	}
	turbineData = setScatterData("Power", "WindSpeed", "Power", colorField[1], "powerAxis", tk.M{"size": 2}, resWSvsPower)
	dataSeries = append(dataSeries, turbineData)
	dataSeries = append(dataSeries, seriesData)

	data := struct {
		Data []tk.M
	}{
		Data: dataSeries,
	}

	return helper.CreateResult(true, data, "success")
}

func (m *AnalyticPowerCurveController) GetPCScatterOperational(k *knot.WebContext) interface{} {
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
		Power         float64
		RotorRPM      float64
		AvgBladeAngle float64
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
	project := p.Project
	minAxisX := 0.0
	maxAxisX := 0.0
	minAxisY := 0.0
	maxAxisY := 0.0

	var filter []*dbox.Filter
	filter = []*dbox.Filter{}
	filter = append(filter, dbox.Ne("_id", ""))
	filter = append(filter, dbox.Gte("timestamp", tStart))
	filter = append(filter, dbox.Lte("timestamp", tEnd))
	filter = append(filter, dbox.Eq("turbine", turbine))
	filter = append(filter, dbox.Eq("projectname", project))
	// filter = append(filter, dbox.Eq("oktime", 600))
	filter = append(filter, dbox.Gt("power", 0))
	filter = append(filter, dbox.Gt("avgwindspeed", 0))
	filter = append(filter, dbox.Eq("available", 1))

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

	data := tk.M{}
	datas := []tk.M{}
	seriesData := tk.M{}
	switch p.ScatterType {
	case "rotor":
		for _, val := range list {
			data = tk.M{}
			if val.RotorRPM < minAxisX {
				minAxisX = val.RotorRPM
			}
			if val.RotorRPM > maxAxisX {
				maxAxisX = val.RotorRPM
			}
			if val.Power < minAxisY {
				minAxisY = val.Power
			}
			if val.Power > maxAxisY {
				maxAxisY = val.Power
			}

			data.Set("Rotor", val.RotorRPM)
			data.Set("Power", val.Power)
			data.Set("valueColor", colorField[1])

			datas = append(datas, data)
		}
		seriesData = setScatterData("Rotor RPM", "Rotor", "Power", colorField[1], "powerAxis", tk.M{"size": 2}, datas)
	case "pitch":
		for _, val := range list {
			data = tk.M{}
			if val.AvgBladeAngle >= -10.0 && val.AvgBladeAngle <= 120 {
				if val.AvgBladeAngle < minAxisX {
					minAxisX = val.AvgBladeAngle
				}
				if val.AvgBladeAngle > maxAxisX {
					maxAxisX = val.AvgBladeAngle
				}
				if val.Power < minAxisY {
					minAxisY = val.Power
				}
				if val.Power > maxAxisY {
					maxAxisY = val.Power
				}
				data.Set("Power", val.Power)
				data.Set("Pitch", val.AvgBladeAngle)
				data.Set("valueColor", colorField[1])
				datas = append(datas, data)
			}
		}
		seriesData = setScatterData("Pitch Angle", "Pitch", "Power", colorField[1], "powerAxis", tk.M{"size": 2}, datas)
	}
	seriesData.Unset("name")
	dataSeries = append(dataSeries, seriesData)

	result := struct {
		Data     []tk.M
		MinAxisX float64
		MaxAxisX float64
		MinAxisY float64
		MaxAxisY float64
	}{
		Data:     dataSeries,
		MinAxisX: minAxisX,
		MaxAxisX: maxAxisX,
		MinAxisY: minAxisY,
		MaxAxisY: maxAxisY,
	}

	return helper.CreateResult(true, result, "success")
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
		Power, AvgWindSpeed             float64
		NacelleDeviation, AvgBladeAngle float64
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
	project := p.Project
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
	filter = append(filter, dbox.Gt("avgwindspeed", 0))
	filter = append(filter, dbox.Eq("available", 1))

	// filter = append(filter, dbox.Eq("oktime", 600))

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
	lessDev := tk.ToFloat64(p.LessValue, 2, tk.RoundingAuto)
	greatDev := tk.ToFloat64(p.GreaterValue, 2, tk.RoundingAuto)

	if p.ScatterType != "pitch" { /*processing data non pitch*/
		for _, val := range list {
			scatterData = tk.M{}
			scatterData.Set("WindSpeed", val.AvgWindSpeed)
			scatterData.Set("Power", val.Power)

			if val.NacelleDeviation < lessDev {
				scatterDatas1 = append(scatterDatas1, scatterData)
			}
			if val.NacelleDeviation > greatDev {
				scatterDatas2 = append(scatterDatas2, scatterData)
			}
		}
	} else {
		for _, val := range list { /*processing pitch data*/
			scatterData = tk.M{}
			scatterData.Set("WindSpeed", val.AvgWindSpeed)
			scatterData.Set("Power", val.Power)

			if val.AvgBladeAngle >= -10.0 && val.AvgBladeAngle <= 120.0 {
				if val.AvgBladeAngle < lessDev {
					scatterDatas1 = append(scatterDatas1, scatterData)
				}
				if val.AvgBladeAngle > greatDev {
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
	project := p.Project
	colors := p.Color
	// colordeg := p.ColorDeg
	colorIndex := map[string]int{}
	for key, val := range colorField {
		colorIndex[val] = key
	}
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

	turbineName, e := helper.GetTurbineNameList(p.Project)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	selArr := 0
	for _, turbineX := range turbine {
		var filter []*dbox.Filter
		filter = append(filter, dbox.Gte("dateinfo.dateid", tStart))
		filter = append(filter, dbox.Lte("dateinfo.dateid", tEnd))
		filter = append(filter, dbox.Eq("turbine", turbineX))
		filter = append(filter, dbox.Eq("projectname", project))
		filter = append(filter, dbox.Eq("available", 1))

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
		turbineData.Set("name", "Scatter-"+turbineName[turbineX.(string)])
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
						// datas.Set("valueColor", colordeg[selArr])
						datas.Set("valueColor", colorFieldDegradation[colorIndex[tk.ToString(colors[selArr])]])
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
						// datas.Set("valueColor", colordeg[selArr])
						datas.Set("valueColor", colorFieldDegradation[colorIndex[tk.ToString(colors[selArr])]])
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
		filter = append(filter, dbox.Eq("available", 1))

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
		"color":     "#ff9933",
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
