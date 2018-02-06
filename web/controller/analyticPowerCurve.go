package controller

import (
	. "eaciit/wfdemo-git/library/core"
	. "eaciit/wfdemo-git/library/helper"
	. "eaciit/wfdemo-git/library/models"
	"eaciit/wfdemo-git/web/helper"
	"strings"

	// "fmt"
	"math"
	"sort"
	"sync"
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

	pcData, e := getPCData(project, p.Engine, true)
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

	pcData, e := getPCData(project, p.Engine, true)
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
		pipesx       []tk.M
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
	tNow := time.Now()
	if tEnd.Sub(tNow).Hours() > 0.0 {
		tEnd, _ = time.Parse("20060102", tNow.Format("20060102"))
	}
	// tk.Printf("TEnd : %#v\n", tEnd.Format("2006-01-02 15:04:05"))
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

	issitespecific := false
	if p.IsSpecific && p.ViewSession != "density" {
		issitespecific = true
	}

	pcData, e := getPCData(project, p.Engine, issitespecific)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	pipes = append(pipes, tk.M{"$group": tk.M{"_id": tk.M{"colId": colId, "Turbine": "$turbine"}, "production": tk.M{"$avg": colValue}, "totaldata": tk.M{"$sum": 1}}})
	pipes = append(pipes, tk.M{"$sort": tk.M{"_id": 1}})

	dVal := (tk.ToFloat64(tk.ToInt(DeviationVal, tk.RoundingAuto), 2, tk.RoundingUp) / 100.0)

	filter = append(filter, dbox.Ne("_id", ""))
	filter = append(filter, dbox.Gte("dateinfo.dateid", tStart))
	filter = append(filter, dbox.Lte("dateinfo.dateid", tEnd))
	filter = append(filter, dbox.Ne("turbine", ""))
	filter = append(filter, dbox.Eq("projectname", project))
	filter = append(filter, dbox.In("turbine", turbine...))

	// temporary
	filter = append(filter, dbox.Ne("power", 0.0))

	filter = append(filter, dbox.Ne("power", nil))

	filter = append(filter, dbox.Ne("avgwindspeed", nil))

	//// as per Neeraj Request on Oct 23, 2017
	// if !p.IsPower0 {
	// 	filter = append(filter, dbox.Gt("power", 0))
	// }

	// if viewSession == "density" {
	// 	filter = append(filter, dbox.Gt("denadjwindspeed", 3.0))
	// } else {
	// 	filter = append(filter, dbox.Gte("avgwindspeed", 3.0))
	// }

	filter = append(filter, dbox.Eq("available", 1))

	// modify by ams, 2017-08-11
	if IsDeviation {
		if DeviationOpr > 0 {
			filter = append(filter, dbox.Or(dbox.Gte(colDeviation, dVal), dbox.Lte(colDeviation, (-1.0*dVal))))
		} else {
			filter = append(filter, dbox.And(dbox.Lte(colDeviation, dVal), dbox.Gte(colDeviation, (-1.0*dVal))))
		}
	}
	if isClean {
		filter = append(filter, dbox.Eq("isvalidstate", true))
		// filter = append(filter, dbox.Eq("oktime", 600))
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

	pipesx = append(pipesx, tk.M{"$group": tk.M{"_id": tk.M{"Turbine": "$turbine"}, "totaldata": tk.M{"$sum": 1}}})
	pipesx = append(pipesx, tk.M{"$sort": tk.M{"_id": 1}})

	var filterx []*dbox.Filter
	filterx = append(filterx, dbox.Gte("dateinfo.dateid", tStart))
	filterx = append(filterx, dbox.Lte("dateinfo.dateid", tEnd))
	filterx = append(filterx, dbox.Eq("projectname", project))
	filterx = append(filterx, dbox.Ne("power", 0.0))
	filterx = append(filterx, dbox.Eq("available", 1))
	filterx = append(filterx, dbox.Ne("turbine", ""))
	filterx = append(filterx, dbox.Ne("_id", ""))
	filterx = append(filterx, dbox.In("turbine", turbine...))

	// if IsDeviation {
	// 	if DeviationOpr > 0 {
	// 		filterx = append(filterx, dbox.Or(dbox.Gte(colDeviation, dVal), dbox.Lte(colDeviation, (-1.0*dVal))))
	// 	} else {
	// 		filterx = append(filterx, dbox.Or(dbox.Lte(colDeviation, dVal), dbox.Gte(colDeviation, (-1.0*dVal))))
	// 	}
	// }
	// if isClean {
	// 	filterx = append(filterx, dbox.Eq("isvalidstate", true))
	// 	// filter = append(filter, dbox.Eq("oktime", 600))
	// }

	var listAll []tk.M
	csrAll, e := DB().Connection.NewQuery().
		From(new(ScadaData).TableName()).
		Command("pipe", pipesx).
		Where(dbox.And(filterx...)).
		Cursor(nil)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	e = csrAll.Fetch(&listAll, 0, false)
	defer csrAll.Close()

	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	totalAllPerTurbines := map[string]tk.M{}
	totalDays := tk.ToInt(tk.Div(tEnd.Sub(tStart).Hours(), 24.0), "0") + 1
	totalDataShouldBe := totalDays * 144
	totalDataAll := 0
	if len(listAll) > 0 {
		for _, dt := range listAll {
			id := dt.Get("_id").(tk.M)
			turbine := id.GetString("Turbine")
			totalDataPerTurbine := dt.GetInt("totaldata")
			totalDataAll += totalDataPerTurbine

			perTurbine := tk.M{
				"totaldata":     totalDataPerTurbine,
				"totalshouldbe": totalDataShouldBe,
				"totaldays":     totalDays,
				"avail":         tk.Div(float64(totalDataPerTurbine), float64(totalDataShouldBe)),
			}

			totalAllPerTurbines[turbine] = perTurbine
		}
	}
	totalDataAvail := tk.Div(float64(totalDataAll), (float64(totalDataShouldBe) * float64(len(p.Turbine))))

	if len(p.Turbine) == 0 {
		for _, listVal := range list {
			exist := false
			for _, val := range p.Turbine {
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

		totalDataPerTurbine := 0
		for _, val := range exist {
			idD := val.Get("_id").(tk.M)
			totalData := val.GetInt("totaldata")
			totalDataPerTurbine += totalData
			datas = append(datas, []float64{idD.GetFloat64("colId"), val.GetFloat64("production")}) //tk.Div(val.GetFloat64("production"), val.GetFloat64("totaldata"))
		}

		turbineData.Set("totaldays", totalDays)
		turbineData.Set("totaldatashouldbe", totalDataShouldBe)
		turbineData.Set("totaldata", totalDataPerTurbine)
		turbineData.Set("dataavailpct", tk.Div(float64(totalDataPerTurbine), float64(totalDataShouldBe)))
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

	dataSeries = append(dataSeries, pcData)

	data := struct {
		Data              []tk.M
		TotalData         int
		TotalDataAvail    float64
		TotalDataShouldBe float64
		TotalPerTurbine   map[string]tk.M
	}{
		Data:              dataSeries,
		TotalData:         totalDataAll,
		TotalDataAvail:    totalDataAvail,
		TotalPerTurbine:   totalAllPerTurbines,
		TotalDataShouldBe: (float64(totalDataShouldBe) * float64(len(p.Turbine))),
	}

	return helper.CreateResult(true, data, "success")
}

func (m *AnalyticPowerCurveController) GetListPowerCurveMonthly(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	/* ================= EXPECTED RESULT ==================
	Power against Wind Speed for each Month, for each Turbine
	*/

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
	//tStart, tEnd, e := helper.GetStartEndDate(k, p.Period, p.DateStart, p.DateEnd)

	//tk.Printf("TEnd : %#v\n", tEnd.Format("2006-01-02T15:04:05Z"))

	colId := "$wsavgforpc"
	colValue := "$power"

	match := []tk.M{}
	match = append(match, tk.M{"_id": tk.M{"$ne": ""}})
	match = append(match, tk.M{"dateinfo.dateid": tk.M{"$gte": tStart}})
	match = append(match, tk.M{"dateinfo.dateid": tk.M{"$lt": tEnd}})
	match = append(match, tk.M{"turbine": tk.M{"$ne": ""}})
	match = append(match, tk.M{"power": tk.M{"$ne": 0.0}})
	match = append(match, tk.M{"power": tk.M{"$ne": nil}})
	match = append(match, tk.M{"avgwindspeed": tk.M{"$ne": nil}})
	//match = append(match, tk.M{"power": tk.M{"$gte": 0}})
	//match = append(match, tk.M{"$or": []tk.M{
	//		tk.M{"$and": []tk.M{tk.M{"power": tk.M{"$lt": 10}}, tk.M{"avgwindspeed": tk.M{"$lt": 3}}}},
	//		tk.M{"$and": []tk.M{tk.M{"power": tk.M{"$gte": 10}}, tk.M{"avgwindspeed": tk.M{"$gte": 3}}}}}})

	//match = append(match, tk.M{"oktime": 600})
	match = append(match, tk.M{"isvalidstate": true})
	match = append(match, tk.M{"available": 1})

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

	for _, listVal := range list { /* only to get list of turbine and list of month for sorting purpose*/
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
	pcData, e := getPCData(project, p.Engine, true)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}
	dataSeries = append(dataSeries, pcData)

	turbineName := map[string]string{}
	if p.Project != "" {
		turbineName, _ = helper.GetTurbineNameList(p.Project)
	}

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

			monthIndex.Set(tk.ToString(selArr), simpleMonth) /*{"1": "jan 16", "2": "feb 16", ...}*/

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

func (m *AnalyticPowerCurveController) GetListPowerCurveMonthlyScatter(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson
	var (
		pipes      []tk.M
		dataSeries []tk.M
	)

	p := new(PayloadAnalyticPC)
	e := k.GetPayload(&p)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	project := p.Project
	tStart, _ := time.Parse("20060102", p.DateStart.Format("200601")+"01")
	tEnd := tStart.AddDate(0, 1, 0)

	tNow := time.Now()
	if tEnd.Sub(tNow).Hours() > 0.0 {
		tEnd, _ = time.Parse("20060102", tNow.Format("20060102"))
		tEnd = tEnd.AddDate(0, 0, 1)
	}

	pcData, e := getPCData(project, p.Engine, true)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	match := []tk.M{}
	match = append(match, tk.M{"dateinfo.dateid": tk.M{"$gte": tStart}})
	match = append(match, tk.M{"dateinfo.dateid": tk.M{"$lt": tEnd}})
	//match = append(match, tk.M{"power": tk.M{"$gt": 0}})
	match = append(match, tk.M{"power": tk.M{"$ne": 0.0}})
	//match = append(match, tk.M{"oktime": 600})
	match = append(match, tk.M{"_id": tk.M{"$ne": ""}})
	//match = append(match, tk.M{"turbine": tk.M{"$ne": ""}})
	//match = append(match, tk.M{"power": tk.M{"$ne": 0.0}})
	match = append(match, tk.M{"projectname": project})
	match = append(match, tk.M{"isvalidstate": true})
	match = append(match, tk.M{"available": 1})
	match = append(match, tk.M{"power": tk.M{"$ne": nil}})
	match = append(match, tk.M{"avgwindspeed": tk.M{"$ne": nil}})

	if len(p.Turbine) > 0 {
		match = append(match, tk.M{"turbine": tk.M{"$in": p.Turbine}})
	}

	pipes = append(pipes, tk.M{"$match": tk.M{"$and": match}})
	pipes = append(pipes, tk.M{"$project": tk.M{"turbine": 1, "power": 1, "avgwindspeed": 1}})
	pipes = append(pipes, tk.M{"$group": tk.M{
		"_id": tk.M{
			"turbine":      "$turbine",
			"avgwindspeed": "$avgwindspeed",
		},
		"power":     tk.M{"$avg": "$power"},
		"totaldata": tk.M{"$sum": 1},
	}})

	csr, e := DB().Connection.NewQuery().
		From(new(ScadaData).TableName()).
		Command("pipe", pipes).
		Cursor(nil)

	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}
	alltkm := []tk.M{}
	e = csr.Fetch(&alltkm, 0, false)
	defer csr.Close()

	// getting all turbines total availability
	match = []tk.M{}
	match = append(match, tk.M{"dateinfo.dateid": tk.M{"$gte": tStart}})
	match = append(match, tk.M{"dateinfo.dateid": tk.M{"$lt": tEnd}})
	match = append(match, tk.M{"power": tk.M{"$ne": 0.0}})
	match = append(match, tk.M{"_id": tk.M{"$ne": ""}})
	match = append(match, tk.M{"projectname": project})
	match = append(match, tk.M{"available": 1})
	// match = append(match, tk.M{"power": tk.M{"$ne": nil}})
	// match = append(match, tk.M{"avgwindspeed": tk.M{"$ne": nil}})

	if len(p.Turbine) > 0 {
		match = append(match, tk.M{"turbine": tk.M{"$in": p.Turbine}})
	}

	//tk.Printf("%#v\n", match)

	pipes = []tk.M{}
	pipes = append(pipes, tk.M{"$match": tk.M{"$and": match}})
	pipes = append(pipes, tk.M{"$group": tk.M{
		"_id":       "$turbine",
		"totaldata": tk.M{"$sum": 1},
	}})

	//tk.Printf("%#v\n", pipes)

	csrta, e := DB().Connection.NewQuery().
		From(new(ScadaData).TableName()).
		Command("pipe", pipes).
		Cursor(nil)

	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}
	allta := []tk.M{}
	e = csrta.Fetch(&allta, 0, false)
	defer csrta.Close()

	var datas [][]float64
	results := []tk.M{}
	sortTurbines := []string{}

	totalDays := tk.Div(tEnd.Sub(tStart).Hours(), 24.0)
	totalDataShouldBe := totalDays * 144

	resData := tk.M{}
	resTotal := tk.M{}
	for _, tkm := range alltkm {
		ids, _ := tk.ToM(tkm["_id"])
		sturbine := ids.GetString("turbine")

		lfloat64 := resData.Get(sturbine, map[float64]float64{}).(map[float64]float64)
		ws, pwr := tk.ToFloat64(ids.Get("avgwindspeed"), 3, tk.RoundingAuto), tk.ToFloat64(tkm.Get("power"), 3, tk.RoundingAuto)
		totalPerRec := tkm.GetInt("totaldata")
		totalPerTurbine := resTotal.Get(sturbine, 0).(int)
		totalAllTurbine := totalPerTurbine + totalPerRec
		resTotal.Set(sturbine, totalAllTurbine)

		// if ws < 3 && pwr > 10 {
		// 	continue
		// }

		lfloat64[ws] = pwr
		resData.Set(sturbine, lfloat64)
	}
	//tk.Printf("%#v\n", allta)
	resAllTa := tk.M{}
	for _, ta := range allta {
		resAllTa.Set(ta.GetString("_id"), ta.GetInt("totaldata"))
	}
	//tk.Printf("%#v\n", resAllTa)

	for _, turX := range p.Turbine {
		sortTurbines = append(sortTurbines, tk.ToString(turX))
	}
	sort.Strings(sortTurbines)

	turbineName := map[string]string{}
	if p.Project != "" {
		turbineName, _ = helper.GetTurbineNameList(p.Project)
	}

	for _, turbineX := range sortTurbines {
		dataSeries = []tk.M{}                   /*clear variable for next data*/
		dataSeries = append(dataSeries, pcData) /*always append expected value at beginning*/

		datas = [][]float64{}
		for ws, power := range resData.Get(turbineX, map[float64]float64{}).(map[float64]float64) {
			datas = append(datas, []float64{ws, power})
		}

		monthData := tk.M{}
		monthData.Set("name", turbineName[turbineX])
		monthData.Set("turbineid", turbineX)
		monthData.Set("type", "scatter")
		monthData.Set("style", "smooth")
		monthData.Set("dashType", "solid")
		monthData.Set("markers", tk.M{"visible": true, "size": 1})
		monthData.Set("width", 2)
		monthData.Set("color", "#21c4af")
		monthData.Set("idxseries", "Data")
		monthData.Set("data", datas)
		dataSeries = append(dataSeries, monthData)

		totalPerTurbine := resTotal.GetInt(turbineX)

		totalDataTurbine := 0
		totalDataAvail := 0.0
		if resAllTa.Has(turbineX) {
			totalDataTurbine = resAllTa.GetInt(turbineX)
			if totalDataTurbine > 0 {
				totalDataAvail = tk.Div(float64(totalDataTurbine), totalDataShouldBe)
			}
		}

		turbineData := tk.M{
			"NameID":                turbineX,
			"Name":                  turbineName[turbineX],
			"Data":                  dataSeries,
			"DataTotalAvailability": totalDataAvail,
			"DataTotalAll":          totalDataTurbine,
			"DataAvailability":      tk.Div(float64(totalPerTurbine), totalDataShouldBe),
			"DataTotal":             totalPerTurbine,
			"TotalDays":             totalDays,
			"TotalShouldBe":         totalDataShouldBe,
		}
		results = append(results, turbineData)
	}

	categoryList := []tk.M{}
	categoryList = append(categoryList, tk.M{"category": "Power Curve", "color": "#ff9933"})
	categoryList = append(categoryList, tk.M{"category": "Data", "color": "#21c4af"})

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
	PC2project := p.PC2Project

	colId := "$wsavgforpc"
	colValue := "$power"

	// Hardcode first
	engine := ""
	if PC1project == "Dewas" {
		engine = "S-97"
	}

	PC1Data, e := getPCData(PC1project, engine, true)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}
	dataSeries = append(dataSeries, PC1Data)

	turbineName1, e := helper.GetTurbineNameList(PC1project)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}
	turbineName2 := turbineName1
	legendName1 := turbineName1[p.PC1Turbine] + " (" + PC1tStart.Format("02-Jan-2006") + "  to " + PC1tEnd.Format("02-Jan-2006") + ")"
	legendName2 := turbineName2[p.PC2Turbine] + " (" + PC2tStart.Format("02-Jan-2006") + "  to " + PC2tEnd.Format("02-Jan-2006") + ")"

	if PC1project != PC2project {
		engine = ""
		if PC2project == "Dewas" {
			engine = "S-97"
		}
		PC2Data, e := getPCData(PC2project, engine, true)
		PC2Data.Set("name", "Power Curve "+PC2project)
		PC2Data.Set("idxseries", 1)
		PC2Data.Set("color", colorField[1])
		dataSeries = append(dataSeries, PC2Data)
		dataSeries[0].Set("name", "Power Curve "+PC1project)

		turbineName2, e = helper.GetTurbineNameList(PC2project)
		if e != nil {
			return helper.CreateResult(false, nil, e.Error())
		}
		legendName2 = turbineName2[p.PC2Turbine] + " (" + PC2tStart.Format("02-Jan-2006") + "  to " + PC2tEnd.Format("02-Jan-2006") + ")"
	}

	pipes = append(pipes, tk.M{"$group": tk.M{"_id": colId, "production": tk.M{"$avg": colValue}, "totaldata": tk.M{"$sum": 1}}})
	pipes = append(pipes, tk.M{"$sort": tk.M{"_id": 1}})

	filter = nil
	filter = append(filter, dbox.Ne("_id", ""))
	filter = append(filter, dbox.Eq("projectname", PC1project))
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

	defer csr.Close()

	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}
	e = csr.Fetch(&list, 0, false)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	var datas [][]float64
	turbineData := tk.M{}
	turbineData.Set("name", legendName1)
	turbineData.Set("type", "scatterLine")
	turbineData.Set("style", "smooth")
	turbineData.Set("dashType", "solid")
	turbineData.Set("markers", tk.M{"visible": false})
	turbineData.Set("width", 2)
	turbineData.Set("color", colorField[4])

	for _, val := range list {
		datas = append(datas, []float64{val.GetFloat64("_id"), val.GetFloat64("production")})
	}

	if len(datas) > 0 {
		turbineData.Set("data", datas)
	}

	dataSeries = append(dataSeries, turbineData)

	filter = nil
	filter = append(filter, dbox.Ne("_id", ""))
	filter = append(filter, dbox.Eq("projectname", PC2project))
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

	defer csr.Close()
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}
	list = []tk.M{}
	e = csr.Fetch(&list, 0, false)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	var datasC2 [][]float64
	turbineData = tk.M{}
	turbineData.Set("name", legendName2)
	turbineData.Set("type", "scatterLine")
	turbineData.Set("style", "smooth")
	turbineData.Set("dashType", "solid")
	turbineData.Set("markers", tk.M{"visible": false})
	turbineData.Set("width", 2)
	turbineData.Set("color", colorField[7])

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

		if val.Has(field) {
			dataWSvsTipe.Set("WindSpeed", val.GetFloat64("avgwindspeed"))
			dataWSvsTipe.Set(tipe, val.GetFloat64(field))
			dataWSvsTipe.Set("valueColor", colorField[2])

			resWSvsTipe = append(resWSvsTipe, dataWSvsTipe)
		}
	}
	return
}

func getScatterValue10Min(list []tk.M, tipe, field string, isTI bool) (resWSvsPower []tk.M, resWSvsTipe []tk.M) {
	resWSvsPower = []tk.M{}
	resWSvsTipe = []tk.M{}
	dataWSvsPower := tk.M{}
	dataWSvsTipe := tk.M{}
	for _, val := range list {
		dataWSvsPower = tk.M{}
		dataWSvsTipe = tk.M{}

		dataWSvsPower.Set("WindSpeed", val.GetFloat64("windspeed_ms"))
		dataWSvsPower.Set("Power", val.GetFloat64("activepower_kw"))
		dataWSvsPower.Set("valueColor", colorField[1])

		resWSvsPower = append(resWSvsPower, dataWSvsPower)

		dataWSvsTipe.Set("WindSpeed", val.GetFloat64("windspeed_ms"))
		if isTI {
			dataWSvsTipe.Set(tipe, tk.Div(val.GetFloat64(field), val.GetFloat64("windspeed_ms")))
		} else {
			dataWSvsTipe.Set(tipe, val.GetFloat64(field))
		}
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
		Engine      string
		ScatterType string
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
	pcData, e := getPCData(project, p.Engine, true)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}
	dataSeries = append(dataSeries, pcData)

	var filter []*dbox.Filter

	scatterType := p.ScatterType
	isScada10Min := true
	if scatterType == "temp" || scatterType == "deviation" || scatterType == "pitch" || scatterType == "ambient" {
		isScada10Min = false
	}
	_list := tk.M{}
	var csr dbox.ICursor
	if !isScada10Min {
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

		csr, e = DB().Connection.NewQuery().
			Select("power", "avgwindspeed", "avgbladeangle", "nacelletemperature", "winddirection", "ambienttemperature").
			From(new(ScadaData).TableName()).
			Where(dbox.And(filter...)).
			Take(10000).
			Cursor(nil)
		if e != nil {
			return helper.CreateResult(false, nil, e.Error())
		}

	} else {
		filter = []*dbox.Filter{}
		filter = append(filter, dbox.Ne("_id", ""))
		filter = append(filter, dbox.Gte("timestamp", tStart))
		filter = append(filter, dbox.Lte("timestamp", tEnd))
		filter = append(filter, dbox.Eq("turbine", turbine))
		filter = append(filter, dbox.Eq("projectname", project))
		filter = append(filter, dbox.Gt("windspeed_ms", 0))
		filter = append(filter, dbox.Gt("activepower_kw", 0))

		csr, e = DB().Connection.NewQuery().
			Select("activepower_kw", "windspeed_ms_stddev", "windspeed_ms").
			From("Scada10MinHFD").
			Where(dbox.And(filter...)).
			Take(10000).
			Cursor(nil)
		if e != nil {
			return helper.CreateResult(false, nil, e.Error())
		}
	}

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

	switch scatterType {
	case "temp":
		resWSvsPower, resWSvsTipe = getScatterValue(list, "NacelleTemperature", "nacelletemperature")
		seriesData = setScatterData("Nacelle Temperature", "WindSpeed", "NacelleTemperature", colorField[2], "tempAxis", tk.M{"size": 2}, resWSvsTipe)
	case "deviation":
		resWSvsPower, resWSvsTipe = getScatterValue(list, "NacelleDeviation", "winddirection")
		seriesData = setScatterData("Nacelle Deviation", "WindSpeed", "NacelleDeviation", colorField[2], "deviationAxis", tk.M{"size": 2}, resWSvsTipe)
	case "pitch":
		resWSvsPower, resWSvsTipe = getScatterValue(list, "PitchAngle", "avgbladeangle")
		seriesData = setScatterData("Pitch Angle", "WindSpeed", "PitchAngle", colorField[2], "pitchAxis", tk.M{"size": 2}, resWSvsTipe)
	case "ambient":
		resWSvsPower, resWSvsTipe = getScatterValue(list, "AmbientTemperature", "ambienttemperature")
		seriesData = setScatterData("Ambient Temperature", "WindSpeed", "AmbientTemperature", colorField[2], "ambientAxis", tk.M{"size": 2}, resWSvsTipe)
	case "windspeed_dev":
		resWSvsPower, resWSvsTipe = getScatterValue10Min(list, "WindSpeedStdDev", "windspeed_ms_stddev", false)
		seriesData = setScatterData("Wind Speed Std. Dev.", "WindSpeed", "WindSpeedStdDev", colorField[2], "windspeed_dev", tk.M{"size": 2}, resWSvsTipe)
	case "windspeed_ti":
		resWSvsPower, resWSvsTipe = getScatterValue10Min(list, "WindSpeedTI", "windspeed_ms_stddev", true)
		seriesData = setScatterData("Wind Speed TI", "WindSpeed", "WindSpeedTI", colorField[2], "windspeed_ti", tk.M{"size": 2}, resWSvsTipe)
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

	type PayloadOperational struct {
		Period      string
		Project     string
		Turbine     string
		DateStart   time.Time
		DateEnd     time.Time
		ScatterType string
	}

	payload := []*PayloadOperational{}
	e := k.GetPayload(&payload)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}
	minAxisX := 0.0
	maxAxisX := 0.0
	minAxisY := 0.0
	maxAxisY := 0.0
	dataSeries := []tk.M{}

	var mux sync.Mutex
	var wg sync.WaitGroup
	wg.Add(2)
	turbineNameOrder := []string{} /* for sort the result */
	for idx, p := range payload {
		idx++
		turbineName, e := helper.GetTurbineNameList(p.Project)
		turbineNameOrder = append(turbineNameOrder, turbineName[p.Turbine])
		if e != nil {
			return helper.CreateResult(false, nil, e.Error())
		}
		tStart, tEnd, e := helper.GetStartEndDate(k, p.Period, p.DateStart, p.DateEnd)
		if e != nil {
			return helper.CreateResult(false, nil, e.Error())
		}
		go func(p *PayloadOperational, dataSeries *[]tk.M, minAxisX, maxAxisX, minAxisY, maxAxisY *float64, index int,
			turbineName map[string]string, tStart, tEnd time.Time, wg *sync.WaitGroup) {
			defer wg.Done()
			list := []tk.M{}
			if e != nil {
				return
			}
			turbine := p.Turbine
			project := p.Project

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
				return
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

			data := tk.M{}
			datas := []tk.M{}
			seriesData := tk.M{}
			typeDetail := map[string]tk.M{
				"rotor": tk.M{
					"field":      "rotorrpm",
					"seriesname": "Rotor RPM",
					"xfieldname": "Rotor",
				},
				"pitch": tk.M{
					"field":      "avgbladeangle",
					"seriesname": "Pitch Angle",
					"xfieldname": "Pitch",
				},
				"generatorrpm": tk.M{
					"field":      "generatorrpm",
					"seriesname": "Generator RPM",
					"xfieldname": "Generator",
				},
				"windspeed": tk.M{
					"field":      "avgwindspeed",
					"seriesname": "Wind Speed",
					"xfieldname": "WindSpeed",
				},
			}
			typeSelected := typeDetail[p.ScatterType]
			fieldName := typeSelected.GetString("field")
			seriesName := typeSelected.GetString("seriesname")
			xFieldName := typeSelected.GetString("xfieldname")
			for _, val := range list {
				data = tk.M{}
				if val.GetFloat64(fieldName) < *minAxisX {
					mux.Lock()
					*minAxisX = val.GetFloat64(fieldName)
					mux.Unlock()
				}
				if val.GetFloat64(fieldName) > *maxAxisX {
					mux.Lock()
					*maxAxisX = val.GetFloat64(fieldName)
					mux.Unlock()
				}
				if val.GetFloat64("power") < *minAxisY {
					mux.Lock()
					*minAxisY = val.GetFloat64("power")
					mux.Unlock()
				}
				if val.GetFloat64("power") > *maxAxisY {
					mux.Lock()
					*maxAxisY = val.GetFloat64("power")
					mux.Unlock()
				}

				data.Set(xFieldName, val.GetFloat64(fieldName))
				data.Set("Power", val.GetFloat64("power"))
				data.Set("valueColor", colorField[index])

				datas = append(datas, data)
			}
			seriesData = setScatterData(seriesName, xFieldName, "Power", colorField[index], "powerAxis", tk.M{"size": 2}, datas)
			seriesData.Set("name", turbineName[turbine]+" ("+p.DateStart.Format("02-Jan-2006")+" to "+p.DateEnd.Format("02-Jan-2006")+")")
			mux.Lock()
			*dataSeries = append(*dataSeries, seriesData)
			mux.Unlock()

		}(p, &dataSeries, &minAxisX, &maxAxisX, &minAxisY, &maxAxisY, idx, turbineName, tStart, tEnd, &wg)
	}
	wg.Wait()

	/* sort the dataseries after processing on go routine */
	tempResult := []tk.M{}
	for _, val := range turbineNameOrder {
		for _, dtSeries := range dataSeries {
			split := strings.Split(dtSeries.GetString("name"), " (")
			if split[0] == val {
				tempResult = append(tempResult, dtSeries)
			}
		}
	}
	dataSeries = tempResult

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
		Engine        string
		ScatterType   string
		LessValue     float64
		GreaterValue  float64
		LessColor     string
		GreaterColor  string
		GreaterMarker string
		LessMarker    string
	}

	type ScadaMini struct {
		Power, AvgWindSpeed               float64
		AvgBladeAngle                     float64
		WindDirection, NacelleTemperature float64
		AmbientTemperature                float64
	}

	var (
		list       []ScadaMini
		dataSeries []tk.M
		list10Min  []tk.M
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
	pcData, e := getPCData(project, p.Engine, true)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}
	dataSeries = append(dataSeries, pcData)

	isScada10Min := true
	if p.ScatterType == "temp" || p.ScatterType == "deviation" || p.ScatterType == "pitch" || p.ScatterType == "ambient" {
		isScada10Min = false
	}

	powerData := []tk.M{}

	pipes := []tk.M{}
	if !isScada10Min {
		/*=======POWER LINE QUERY =========*/
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
		// filter = append(filter, dbox.Gt("power", 0))
		// filter = append(filter, dbox.Gt("avgwindspeed", 0))
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

		e = csrPower.Fetch(&powerData, 0, false)
		defer csrPower.Close()

		/*===== END OF POWER LINE =======*/
		// filter is same with power filter

		csr, e := DB().Connection.NewQuery().
			From(new(ScadaData).TableName()).
			Where(dbox.And(filter...)).
			//Take(10000).
			Cursor(nil)

		if e != nil {
			return helper.CreateResult(false, nil, e.Error())
		}
		var _list ScadaMini
		for {
			_list = ScadaMini{}
			e = csr.Fetch(&_list, 1, false)
			if e != nil {
				break
			}
			list = append(list, _list)
		}

		defer csr.Close()
	} else {
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
		// filter = append(filter, dbox.Gt("power", 0))
		// filter = append(filter, dbox.Gt("avgwindspeed", 0))
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

		e = csrPower.Fetch(&powerData, 0, false)
		defer csrPower.Close()

		/*===== END OF POWER LINE =======*/
		// filter is same with power filter
		var filter2 []*dbox.Filter
		filter2 = []*dbox.Filter{}
		filter2 = append(filter2, dbox.Ne("_id", ""))
		filter2 = append(filter2, dbox.Gte("timestamp", tStart))
		filter2 = append(filter2, dbox.Lte("timestamp", tEnd))
		filter2 = append(filter2, dbox.Eq("turbine", turbine))
		filter2 = append(filter2, dbox.Eq("projectname", project))
		filter2 = append(filter2, dbox.Gt("windspeed_ms", 0))
		filter2 = append(filter2, dbox.Gt("activepower_kw", 0))

		csr, e := DB().Connection.NewQuery().
			Select("activepower_kw", "windspeed_ms_stddev", "windspeed_ms").
			From("Scada10MinHFD").Where(dbox.And(filter2...)).
			//Take(10000).
			Cursor(nil)
		if e != nil {
			return helper.CreateResult(false, nil, e.Error())
		}
		var _list10Min tk.M
		list10Min = []tk.M{}
		// countList := 0
		for {
			_list10Min = tk.M{}
			e = csr.Fetch(&_list10Min, 1, false)
			if e != nil {
				break
			}
			list10Min = append(list10Min, _list10Min)
			// countList++
		}

		defer csr.Close()
	}

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

	scatterData := tk.M{}
	scatterDatas1 := []tk.M{}
	scatterDatas2 := []tk.M{}
	lessDev := tk.ToFloat64(p.LessValue, 2, tk.RoundingAuto)
	greatDev := tk.ToFloat64(p.GreaterValue, 2, tk.RoundingAuto)

	switch p.ScatterType {
	case "pitch":
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
	case "deviation":
		for _, val := range list {
			scatterData = tk.M{}
			scatterData.Set("WindSpeed", val.AvgWindSpeed)
			scatterData.Set("Power", val.Power)
			if val.WindDirection < lessDev {
				scatterDatas1 = append(scatterDatas1, scatterData)
			}
			if val.WindDirection > greatDev {
				scatterDatas2 = append(scatterDatas2, scatterData)
			}
		}
	case "temp":
		for _, val := range list { /*processing pitch data*/
			scatterData = tk.M{}
			scatterData.Set("WindSpeed", val.AvgWindSpeed)
			scatterData.Set("Power", val.Power)

			if val.NacelleTemperature < lessDev {
				scatterDatas1 = append(scatterDatas1, scatterData)
			}
			if val.NacelleTemperature > greatDev {
				scatterDatas2 = append(scatterDatas2, scatterData)
			}
		}
	case "ambient":
		for _, val := range list { /*processing pitch data*/
			scatterData = tk.M{}
			scatterData.Set("WindSpeed", val.AvgWindSpeed)
			scatterData.Set("Power", val.Power)

			if val.AmbientTemperature >= -10.0 && val.AmbientTemperature <= 120.0 {
				if val.AmbientTemperature < lessDev {
					scatterDatas1 = append(scatterDatas1, scatterData)
				}
				if val.AmbientTemperature > greatDev {
					scatterDatas2 = append(scatterDatas2, scatterData)
				}
			}
		}
	case "windspeed_dev":
		for _, val := range list10Min {
			scatterData = tk.M{}
			activepower_kw := val.GetFloat64("activepower_kw")
			windspeed_ms := val.GetFloat64("windspeed_ms")
			windspeed_std_dev := val.GetFloat64("windspeed_ms_stddev")
			scatterData.Set("WindSpeed", windspeed_ms)
			scatterData.Set("Power", activepower_kw)
			if windspeed_std_dev < lessDev {
				scatterDatas1 = append(scatterDatas1, scatterData)
			}
			if windspeed_std_dev > greatDev {
				scatterDatas2 = append(scatterDatas2, scatterData)
			}
		}
	case "windspeed_ti":
		// countdata := 0
		for _, val := range list10Min {
			// if countdata == 10 {
			// 	break
			// }
			scatterData = tk.M{}
			activepower_kw := val.GetFloat64("activepower_kw")
			windspeed_ms := val.GetFloat64("windspeed_ms")
			windspeed_std_dev := val.GetFloat64("windspeed_ms_stddev")
			windspeed_ti := tk.Div(windspeed_std_dev, windspeed_ms)
			scatterData.Set("WindSpeed", windspeed_ms)
			scatterData.Set("Power", activepower_kw)
			if windspeed_ti < lessDev {
				scatterDatas1 = append(scatterDatas1, scatterData)
			}
			if windspeed_ti > greatDev {
				scatterDatas2 = append(scatterDatas2, scatterData)
			}
			// countdata++
		}
	}

	// if p.ScatterType != "pitch" { /*processing data non pitch*/
	// 	for _, val := range list {
	// 		scatterData = tk.M{}
	// 		scatterData.Set("WindSpeed", val.AvgWindSpeed)
	// 		scatterData.Set("Power", val.Power)
	// 		if val.WindDirection < lessDev {
	// 			scatterDatas1 = append(scatterDatas1, scatterData)
	// 		}
	// 		if val.WindDirection > greatDev {
	// 			scatterDatas2 = append(scatterDatas2, scatterData)
	// 		}
	// 	}
	// } else {
	// 	for _, val := range list { /*processing pitch data*/
	// 		scatterData = tk.M{}
	// 		scatterData.Set("WindSpeed", val.AvgWindSpeed)
	// 		scatterData.Set("Power", val.Power)

	// 		if val.AvgBladeAngle >= -10.0 && val.AvgBladeAngle <= 120.0 {
	// 			if val.AvgBladeAngle < lessDev {
	// 				scatterDatas1 = append(scatterDatas1, scatterData)
	// 			}
	// 			if val.AvgBladeAngle > greatDev {
	// 				scatterDatas2 = append(scatterDatas2, scatterData)
	// 			}
	// 		}
	// 	}
	// }
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
	case "ambient":
		seriesData1 := setScatterData("Ambient Temp. < "+tk.ToString(p.LessValue), "WindSpeed", "Power", p.LessColor, "powerAxis", tk.M{"size": 2, "type": p.LessMarker, "background": p.LessColor}, scatterDatas1)
		seriesData1.Unset("colorField")
		dataSeries = append(dataSeries, seriesData1)
		seriesData2 := setScatterData("Ambient Temp. > "+tk.ToString(p.GreaterValue), "WindSpeed", "Power", p.GreaterColor, "powerAxis", tk.M{"size": 2, "type": p.GreaterMarker, "background": p.GreaterColor}, scatterDatas2)
		seriesData2.Unset("colorField")
		dataSeries = append(dataSeries, seriesData2)
	case "temp":
		seriesData1 := setScatterData("Nacelle Temp. < "+tk.ToString(p.LessValue), "WindSpeed", "Power", p.LessColor, "powerAxis", tk.M{"size": 2, "type": p.LessMarker, "background": p.LessColor}, scatterDatas1)
		seriesData1.Unset("colorField")
		dataSeries = append(dataSeries, seriesData1)
		seriesData2 := setScatterData("Nacelle Temp. > "+tk.ToString(p.GreaterValue), "WindSpeed", "Power", p.GreaterColor, "powerAxis", tk.M{"size": 2, "type": p.GreaterMarker, "background": p.GreaterColor}, scatterDatas2)
		seriesData2.Unset("colorField")
		dataSeries = append(dataSeries, seriesData2)
	case "windspeed_dev":
		seriesData1 := setScatterData("Wind Speed Std. Dev. < "+tk.ToString(tk.ToFloat64(p.LessValue, 2, tk.RoundingAuto)), "WindSpeed", "Power", p.LessColor, "powerAxis", tk.M{"size": 2, "type": p.LessMarker, "background": p.LessColor}, scatterDatas1)
		seriesData1.Unset("colorField")
		dataSeries = append(dataSeries, seriesData1)
		seriesData2 := setScatterData("Wind Speed Std. Dev. > "+tk.ToString(tk.ToFloat64(p.GreaterValue, 2, tk.RoundingAuto)), "WindSpeed", "Power", p.GreaterColor, "powerAxis", tk.M{"size": 2, "type": p.GreaterMarker, "background": p.GreaterColor}, scatterDatas2)
		seriesData2.Unset("colorField")
		dataSeries = append(dataSeries, seriesData2)
	case "windspeed_ti":
		seriesData1 := setScatterData("TI Wind Speed < "+tk.ToString(tk.ToFloat64(p.LessValue, 2, tk.RoundingAuto)), "WindSpeed", "Power", p.LessColor, "powerAxis", tk.M{"size": 2, "type": p.LessMarker, "background": p.LessColor}, scatterDatas1)
		seriesData1.Unset("colorField")
		dataSeries = append(dataSeries, seriesData1)
		seriesData2 := setScatterData("TI Wind Speed > "+tk.ToString(tk.ToFloat64(p.GreaterValue, 2, tk.RoundingAuto)), "WindSpeed", "Power", p.GreaterColor, "powerAxis", tk.M{"size": 2, "type": p.GreaterMarker, "background": p.GreaterColor}, scatterDatas2)
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

	tNow := time.Now()
	if tEnd.Sub(tNow).Hours() > 0.0 {
		tEnd, _ = time.Parse("20060102", tNow.Format("20060102"))
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
	DeviationOpr := tk.ToInt(p.DeviationOpr, tk.RoundingAuto)
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

	totalDays := tk.Div(tEnd.Sub(tStart).Hours(), 24.0)
	totalDataShouldBe := totalDays * 144

	selArr := 0
	for _, turbineX := range turbine {
		var filter []*dbox.Filter
		filter = append(filter, dbox.Gte("dateinfo.dateid", tStart))
		filter = append(filter, dbox.Lte("dateinfo.dateid", tEnd))
		filter = append(filter, dbox.Eq("turbine", turbineX))
		filter = append(filter, dbox.Eq("projectname", project))
		filter = append(filter, dbox.Eq("available", 1))

		// temporary
		filter = append(filter, dbox.Ne("power", 0.0))

		//// as per Neeraj Request Oct 23, 2017
		// if !p.IsPower0 {
		// filter = append(filter, dbox.Gt("power", 0.0))
		// }
		// filter = append(filter, dbox.Gte("avgwindspeed", 3))

		// if !IsDeviation {
		// 	filter = append(filter, dbox.Gte(colDeviation, dVal))
		// }
		if IsDeviation {
			if DeviationOpr > 0 {
				filter = append(filter, dbox.Or(dbox.Gte(colDeviation, dVal), dbox.Lte(colDeviation, (-1.0*dVal))))
			} else {
				filter = append(filter, dbox.Or(dbox.Lte(colDeviation, dVal), dbox.Gte(colDeviation, (-1.0*dVal))))
			}
		}
		if isClean {
			filter = append(filter, dbox.Eq("isvalidstate", true))
			// filter = append(filter, dbox.Eq("oktime", 600))
		}
		filter = append(filter, dbox.Ne("_id", ""))

		// csr, e := DB().Connection.NewQuery().From(new(ScadaData).TableName()).Where(dbox.And(filter...)).Take(10000).Cursor(nil)
		csr, e := DB().Connection.NewQuery().From(new(ScadaData).TableName()).Where(dbox.And(filter...)).Cursor(nil)
		if e != nil {
			return helper.CreateResult(false, nil, e.Error())
		}

		e = csr.Fetch(&list, 0, false)

		if e != nil {
			return helper.CreateResult(false, nil, e.Error())
		}

		defer csr.Close()

		totalData := len(list)

		turbineData := tk.M{}
		turbineData.Set("name", "Scatter-"+turbineName[turbineX.(string)])
		turbineData.Set("xField", "WindSpeed")
		turbineData.Set("yField", "Power")
		turbineData.Set("colorField", "valueColor")
		turbineData.Set("type", "scatter")
		turbineData.Set("totaldatashouldbe", totalDataShouldBe)
		turbineData.Set("totaldays", totalDays)
		turbineData.Set("totaldata", totalData)
		turbineData.Set("dataavailpct", tk.Div(float64(totalData), totalDataShouldBe))
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

					if IsDeviation {
						if math.Abs(val.DenDeviationPct) <= dVal {
							// datas.Set("valueColor", colordeg[selArr])
							datas.Set("valueColor", colors[selArr])
						} else {
							datas.Set("valueColor", colorFieldDegradation[colorIndex[tk.ToString(colors[selArr])]])
						}
					} else {
						datas.Set("valueColor", colors[selArr])
					}

					arrDatas = append(arrDatas, datas)
				}
			default:
				isShow := true
				if !p.IsPower0 {
					if val.AvgWindSpeed > 0 && val.Power > 0 {
						isShow = true
					} else {
						isShow = false
					}
				}
				if isShow {

					datas.Set("WindSpeed", val.AvgWindSpeed)
					datas.Set("Power", val.Power)
					if IsDeviation {
						if math.Abs(val.DeviationPct) <= dVal {
							// datas.Set("valueColor", colordeg[selArr])
							datas.Set("valueColor", colors[selArr])
						} else {
							datas.Set("valueColor", colorFieldDegradation[colorIndex[tk.ToString(colors[selArr])]])
						}
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

	pcData, e := getPCData(project, p.Engine, true)
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

func getPCData(project string, engine string, issitespecific bool) (pcData tk.M, e error) {
	powerCurve := []PowerCurveModel{}

	filter := dbox.Eq("model", project)
	if engine != "" {
		filter = dbox.And(filter, dbox.Eq("engine", engine))
	}

	csr, e := DB().Connection.NewQuery().
		From(new(PowerCurveModel).
			TableName()).
		Where(filter).
		Order("windspeed").
		Cursor(nil)
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
		if issitespecific {
			datas = append(datas, []float64{val.WindSpeed, val.Power1})
		} else {
			datas = append(datas, []float64{val.WindSpeed, val.Standard})
		}
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

func getPCFilter(project string, engine string, turbine []interface{}, dateStart time.Time, dateEnd time.Time, isValid bool, isDeviation bool, deviationOpr string, deviationValue string, colDeviation string) []*dbox.Filter {
	var filter []*dbox.Filter

	if project != "" {
		filter = append(filter, dbox.Eq("projectname", project))
	}

	dOpr := tk.ToInt(deviationOpr, tk.RoundingAuto)
	dVal := (tk.ToFloat64(tk.ToInt(deviationValue, tk.RoundingAuto), 2, tk.RoundingUp) / 100.0)

	filter = append(filter, dbox.Ne("_id", ""))
	filter = append(filter, dbox.Gte("dateinfo.dateid", dateStart))
	filter = append(filter, dbox.Lte("dateinfo.dateid", dateEnd))
	filter = append(filter, dbox.Ne("turbine", ""))
	filter = append(filter, dbox.Eq("projectname", project))
	filter = append(filter, dbox.Eq("available", 1))

	if isValid {
		filter = append(filter, dbox.Eq("isvalidstate", true))
	}

	if isDeviation {
		if dOpr > 0 {
			filter = append(filter, dbox.Or(dbox.Gte(colDeviation, dVal), dbox.Lte(colDeviation, (-1.0*dVal))))
		} else {
			filter = append(filter, dbox.And(dbox.Lte(colDeviation, dVal), dbox.Gte(colDeviation, (-1.0*dVal))))
		}
	}

	// temporary
	filter = append(filter, dbox.Ne("power", 0.0))
	filter = append(filter, dbox.Ne("power", nil))
	filter = append(filter, dbox.Ne("avgwindspeed", nil))

	return filter
}
