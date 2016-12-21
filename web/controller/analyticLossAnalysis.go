package controller

import (
	. "eaciit/wfdemo-git/library/core"
	. "eaciit/wfdemo-git/library/models"
	"eaciit/wfdemo-git/web/helper"
	"fmt"
	"strconv"
	"strings"
	"time"

	c "github.com/eaciit/crowd"
	"github.com/eaciit/dbox"
	"github.com/eaciit/knot/knot.v1"
	tk "github.com/eaciit/toolkit"
)

type AnalyticLossAnalysisController struct {
	App
}

var colorFields = [...]string{"#ff880e", "#21c4af", "#ff7663", "#ffb74f", "#a2df53", "#1c9ec4", "#ff63a5", "#f44336", "#D91E18", "#8877A9", "#9A12B3", "#26C281", "#E7505A", "#C49F47", "#ff5597", "#c3260c", "#d4735e", "#ff2ad7", "#34ac8b", "#11b2eb", "#f35838", "#ff0037", "#507ca3", "#ff6565", "#ffd664", "#72aaff", "#795548"}

func CreateAnalyticLossAnalysisController() *AnalyticLossAnalysisController {
	var controller = new(AnalyticLossAnalysisController)
	return controller
}

func (m *AnalyticLossAnalysisController) GetScadaSummaryList(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	var (
		filter []*dbox.Filter
		pipes  []tk.M
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

	if project != "" {
		filter = append(filter, dbox.Eq("projectname", project))
	}
	if len(turbine) != 0 {
		filter = append(filter, dbox.In("turbine", turbine...))
	}

	ids := "$turbine"
	if project == "" {
		ids = "$projectname"
	}

	pipes = append(pipes, tk.M{"$group": tk.M{"_id": ids,
		"Production":       tk.M{"$sum": "$production"},
		"MachineDownLoss":  tk.M{"$sum": "$machinedownloss"},
		"GridDownLoss":     tk.M{"$sum": "$griddownloss"},
		"PCDeviation":      tk.M{"$sum": "$pcdeviation"},
		"ElectricalLosses": tk.M{"$sum": "$electricallosses"},
		"OtherDownLoss":    tk.M{"$sum": "$otherdownloss"},
		"DownTimeDuration": tk.M{"$sum": "$downtimehours"},
		"MachineDownHours": tk.M{"$sum": "$machinedownhours"},
		"GridDownHours":    tk.M{"$sum": "$griddownhours"},
		"LossEnergy":       tk.M{"$sum": "$lostenergy"}}})

	csr, e := DB().Connection.NewQuery().
		From(new(ScadaSummaryDaily).TableName()).
		Command("pipe", pipes).
		Where(dbox.And(filter...)).
		Cursor(nil)

	if e != nil {
		helper.CreateResult(false, nil, e.Error())
	}
	defer csr.Close()

	resultScada := []tk.M{}
	e = csr.Fetch(&resultScada, 0, false)

	LossAnalysisResult := []tk.M{}

	for _, val := range resultScada {
		// dummyRes := []tk.M{}
		LossAnalysisResult = append(LossAnalysisResult, tk.M{
			"Id":               val.GetString("_id"),
			"Production":       val.GetFloat64("Production") / 1000,
			"LossEnergy":       val.GetFloat64("LossEnergy") / 1000,
			"MachineDownHours": val.GetFloat64("MachineDownHours"),
			"GridDownHours":    val.GetFloat64("GridDownHours"),
			"EnergyyMD":        val.GetFloat64("MachineDownLoss") / 1000,
			"EnergyyGD":        val.GetFloat64("GridDownLoss") / 1000,
			"ElectricLoss":     val.GetFloat64("ElectricalLosses") / 1000,
			"PCDeviation":      val.GetFloat64("PCDeviation") / 1000,
			"Others":           val.GetFloat64("OtherDownLoss") / 1000,
			"DownTimeDuration": val.GetFloat64("DownTimeDuration"),
		})
	}

	data := struct {
		Data []tk.M
	}{
		Data: LossAnalysisResult,
	}

	return helper.CreateResult(true, data, "success")
}

func (m *AnalyticLossAnalysisController) GetScadaSummaryChart(k *knot.WebContext) interface{} {
	keys := []string{"Energy Lost Due to Machine Down", "Energy Lost Due to Grid Down", "Others", "Electrical Losses", "PC Deviation"}
	categories := []string{}

	k.Config.OutputType = knot.OutputJson

	var (
		filter     []*dbox.Filter
		pipes      []tk.M
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
	breakDown := p.BreakDown

	filter = append(filter, dbox.Ne("_id", ""))
	filter = append(filter, dbox.Gte("dateinfo.dateid", tStart))
	filter = append(filter, dbox.Lte("dateinfo.dateid", tEnd))

	if project != "" {
		filter = append(filter, dbox.Eq("projectname", project))
	}
	if len(turbine) != 0 {
		filter = append(filter, dbox.In("turbine", turbine...))
	}

	ids := "$" + breakDown

	if breakDown == "dateinfo.monthdesc" {
		pipes = append(pipes, tk.M{"$group": tk.M{"_id": tk.M{}.Set("monthid", "$dateinfo.monthid").Set("monthdesc", "$dateinfo.monthdesc"),
			"LostEnergy":       tk.M{"$sum": "$lostenergy"},
			"MachineDownLoss":  tk.M{"$sum": "$machinedownloss"},
			"GridDownLoss":     tk.M{"$sum": "$griddownloss"},
			"PCDeviation":      tk.M{"$sum": "$pcdeviation"},
			"ElectricalLosses": tk.M{"$sum": "$electricallosses"},
			"OtherDownLoss":    tk.M{"$sum": "$otherdownloss"},
			"Production":       tk.M{"$sum": "$production"}}})
		pipes = append(pipes, tk.M{"$sort": tk.M{"_id.monthid": 1}})
	} else {
		pipes = append(pipes, tk.M{"$group": tk.M{"_id": ids,
			"LostEnergy":       tk.M{"$sum": "$lostenergy"},
			"MachineDownLoss":  tk.M{"$sum": "$machinedownloss"},
			"GridDownLoss":     tk.M{"$sum": "$griddownloss"},
			"PCDeviation":      tk.M{"$sum": "$pcdeviation"},
			"ElectricalLosses": tk.M{"$sum": "$electricallosses"},
			"OtherDownLoss":    tk.M{"$sum": "$otherdownloss"},
			"Production":       tk.M{"$sum": "$production"}}})
		pipes = append(pipes, tk.M{"$sort": tk.M{"_id": 1}})
	}
	csr, e := DB().Connection.NewQuery().
		From(new(ScadaSummaryDaily).TableName()).
		Command("pipe", pipes).
		Where(dbox.And(filter...)).
		Cursor(nil)

	if e != nil {
		helper.CreateResult(false, nil, e.Error())
	}
	defer csr.Close()

	resultScada := []tk.M{}
	e = csr.Fetch(&resultScada, 0, false)
	// add by ams, 2016-10-07
	csr.Close()

	selArr := 1
	for _, key := range keys {
		// series := tk.M{}
		// series.Set("name", key)
		// series.Set("type", "bar")
		// series.Set("axis", "lossBar")
		// series.Set("opacity", 0.6)

		// if key == "Energy Lost Due to Machine Down" ||  key == "Energy Lost Due to Grid Down" {
		// 	series.Set("type", "line")
		// 	series.Set("style", "smooth")
		// 	series.Set("axis", "lossLine")
		// 	series.Set("dashType", "solid")
		// 	series.Set("markers", tk.M{"visible": false})
		// 	series.Set("width", 3)
		// }

		unitMeter := " (MWh)"
		if key == "Energy Lost Due to Machine Down" || key == "Energy Lost Due to Grid Down" {
			unitMeter = " (%)"
		} else if key == "PC Deviation" {
			unitMeter = " (MW)"
		}

		series := tk.M{}
		series.Set("name", key+unitMeter)
		series.Set("type", "line")
		series.Set("style", "smooth")
		series.Set("axis", "lossLine")
		series.Set("dashType", "solid")
		series.Set("markers", tk.M{"visible": false})
		series.Set("width", 3)
		series.Set("color", colorFields[selArr])

		if key == "PC Deviation" || key == "Electrical Losses" {
			series.Set("type", "bar")
			series.Set("axis", "lossBar")
			series.Set("opacity", 0.6)
		}

		var datas []float64
		for _, val := range resultScada {
			LostEnergy := val.GetFloat64("LostEnergy")
			if LostEnergy == 0 {
				LostEnergy = 1
			}
			// tk.Printf("dt %v\n", LostEnergy)

			if key == "Energy Lost Due to Machine Down" {
				datas = append(datas, tk.ToFloat64((val.GetFloat64("MachineDownLoss")/LostEnergy)*100, 2, tk.RoundingAuto))
			} else if key == "Energy Lost Due to Grid Down" {
				datas = append(datas, tk.ToFloat64((val.GetFloat64("GridDownLoss")/LostEnergy)*100, 2, tk.RoundingAuto))
			} else if key == "PC Deviation" {
				datas = append(datas, tk.ToFloat64((val.GetFloat64("PCDeviation")/1000), 2, tk.RoundingAuto))
			} else if key == "Electrical Losses" {
				datas = append(datas, tk.ToFloat64((val.GetFloat64("ElectricalLosses")/1000), 2, tk.RoundingAuto))
			} else if key == "Others" {
				datas = append(datas, tk.ToFloat64((val.GetFloat64("OtherDownLoss")/LostEnergy)*100, 2, tk.RoundingAuto))
			}
		}

		if len(datas) > 0 {
			series.Set("data", datas)
		}
		dataSeries = append(dataSeries, series)
		selArr++
	}

	for _, val := range resultScada {
		id := val.Get("_id")
		monthdescKey := ""
		if breakDown == "dateinfo.monthdesc" {
			_id := val.Get("_id").(tk.M)
			// tk.Printf("_id => %#v\n", _id)
			monthdescKey = tk.ToString(_id["monthdesc"])
		}

		if breakDown == "dateinfo.dateid" {
			dt := id.(time.Time)
			// if index == 0 || dt.Day() == 1 {
			// 	categories = append(categories, tk.ToString(dt.Day())+" "+dt.Month().String()[:3])
			// } else {
			categories = append(categories, tk.ToString(dt.Day()))
			// }
			// categories = append(categories, tk.ToString(dt.Day())+"/"+dt.Month().String()[:3])
		} else if breakDown == "dateinfo.monthdesc" {
			if id != "" {
				categories = append(categories, monthdescKey)
			}
		} else if breakDown == "dateinfo.year" {
			categories = append(categories, tk.ToString(id))
		} else if breakDown == "projectname" {
			categories = append(categories, tk.ToString(id))
		} else if breakDown == "turbine" {
			categories = append(categories, tk.ToString(id))
		}
	}

	result := struct {
		Series     []tk.M
		Categories []string
	}{
		Series:     dataSeries,
		Categories: categories,
	}

	return helper.CreateResult(true, result, "success")
}

// func (m *AnalyticLossAnalysisController) GetScadaList(k *knot.WebContext) interface{} {
// 	k.Config.OutputType = knot.OutputJson

// 	var (
// 		filter     		[]*dbox.Filter
// 		pipes 			[]tk.M
// 	)

// 	p := new(PayloadAnalytic)
// 	e := k.GetPayload(&p)
// 	if e != nil {
// 		return helper.CreateResult(false, nil, e.Error())
// 	}

// 	tStart, _ := time.Parse("2006-01-02", p.DateStart.UTC().Format("2006-01-02"))
// 	tEnd, _ := time.Parse("2006-01-02 15:04:05", p.DateEnd.UTC().Format("2006-01-02")+" 23:59:59")
// 	turbine := p.Turbine
// 	project := p.Project

// 	filter = append(filter, dbox.Ne("_id", ""))
// 	filter = append(filter, dbox.Gte("dateinfo.dateid", tStart))
// 	filter = append(filter, dbox.Lte("dateinfo.dateid", tEnd))

// 	if project != "" {
// 		filter = append(filter, dbox.Eq("projectname", project))
// 	}
// 	if len(turbine) != 0 {
// 		filter = append(filter, dbox.In("turbine", turbine...))
// 	}

// 	ids := "$turbine"
// 	if project == "" {
// 		ids = "$projectname"
// 	}

// 	pipes = append(pipes, tk.M{"$group": tk.M{"_id": ids, "Power": tk.M{"$sum": "$power"}, "Powerlost": tk.M{"$sum": "$powerlost"}}})

// 	csr, e := DB().Connection.NewQuery().
// 	From(new(ScadaData).TableName()).
// 	Command("pipe", pipes).
// 	Where(dbox.And(filter...)).
// 	Cursor(nil)

// 	if e != nil {
// 		helper.CreateResult(false, nil, e.Error())
// 	}
// 	defer csr.Close()

// 	resultScada := []tk.M{}
// 	e = csr.Fetch(&resultScada, 0, false)

// 	LossAnalysisResult := []tk.M{}
// 	for _, val := range resultScada {
// 		// dummyRes := []tk.M{}
// 		val.Set("Id", val.GetString("_id"))
// 		val.Set("Production", (val.GetFloat64("Power")/6)/1000)
// 		val.Set("EnergyyMD", 0)
// 		val.Set("EnergyyGD", 0)
// 		val.Set("PCDeviation", 0)
// 		val.Set("ElectricLoss", 0)
// 		val.Set("Others", 0)
// 		LossAnalysisResult = append(LossAnalysisResult, val)
// 	}

// 	data := struct {
// 		Data []tk.M
// 	}{
// 		Data: LossAnalysisResult,
// 	}

// 	return helper.CreateResult(true, data, "success")
// }

// func (m *AnalyticLossAnalysisController) GetScadaList(k *knot.WebContext) interface{} {
// 	k.Config.OutputType = knot.OutputJson

// 	p := new(helper.Payloads)
// 	e := k.GetPayload(&p)
// 	if e != nil {
// 		return helper.CreateResult(false, nil, e.Error())
// 	}

// 	filter := p.ParseFilter()

// 	tk.Printf("dt %v\n", filter)

// 	query := DB().Connection.NewQuery().From(new(ScadaData).TableName()).Skip(p.Skip).Take(p.Take)
// 	query.Where(dbox.And(filter...))

// 	if len(p.Sort) > 0 {
// 		var arrsort []string
// 		for _, val := range p.Sort {
// 			if val.Dir == "desc" {
// 				arrsort = append(arrsort, strings.ToLower("-"+val.Field))
// 			} else {
// 				arrsort = append(arrsort, strings.ToLower(val.Field))
// 			}
// 		}
// 		query = query.Order(arrsort...)
// 	}
// 	csr, e := query.Cursor(nil)
// 	if e != nil {
// 		return helper.CreateResult(false, nil, e.Error())
// 	}
// 	defer csr.Close()

// 	tmpResult := make([]ScadaData, 0)
// 	results := make([]ScadaData, 0)
// 	e = csr.Fetch(&tmpResult, 0, false)

// 	if e != nil {
// 		return helper.CreateResult(false, nil, e.Error())
// 	}

// 	for _, val := range tmpResult {
// 		val.TimeStamp = val.TimeStamp.UTC()
// 		results = append(results, val)
// 	}

// 	data := struct { Data	[]ScadaData }{ Data:	results }

// 	return helper.CreateResult(true, data, "success")
// }

func (m *AnalyticLossAnalysisController) GetTop10(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	p := new(PayloadAnalytic)
	e := k.GetPayload(&p)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	result := tk.M{}
	duration, e := getDownTimeTopFiltered("duration", p, k)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}
	result.Set("duration", duration)
	frequency, e := getDownTimeTopFiltered("frequency", p, k)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}
	result.Set("frequency", frequency)
	loss, e := getDownTimeTopFiltered("loss", p, k)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}
	result.Set("loss", loss)
	catloss, e := getCatLossTopFiltered("loss", p, k)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}
	result.Set("catloss", catloss)
	catlossduration, e := getCatLossTopFiltered("duration", p, k)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}
	result.Set("catlossduration", catlossduration)
	catlossfreq, e := getCatLossTopFiltered("frequency", p, k)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}
	result.Set("catlossfreq", catlossfreq)

	return helper.CreateResult(true, result, "success")
}

func getCatLossTopFiltered(topType string, p *PayloadAnalytic, k *knot.WebContext) ([]tk.M, error) {
	var result []tk.M
	var e error
	var pipes []tk.M
	match := tk.M{}

	if p != nil {
		// tStart, _ := time.Parse("2006-01-02", p.DateStart.UTC().Format("2006-01-02"))
		// tEnd, _ := time.Parse("2006-01-02 15:04:05", p.DateEnd.UTC().Format("2006-01-02")+" 23:59:59")
		tStart, tEnd, e := helper.GetStartEndDate(k, p.Period, p.DateStart, p.DateEnd)
		if e != nil {
			return result, e
		}

		match.Set("detail.startdate", tk.M{"$gte": tStart, "$lte": tEnd})

		if p.Project != "" {
			match.Set("projectname", p.Project)
		}

		if len(p.Turbine) != 0 {
			match.Set("turbine", tk.M{"$in": p.Turbine})
		}

		pipes = append(pipes, tk.M{"$match": match})

		downCause := tk.M{}
		downCause.Set("aebok", "AEBOK")
		downCause.Set("externalstop", "External Stop")
		downCause.Set("griddown", "Grid Down")
		downCause.Set("internalgrid", "Internal Grid")
		downCause.Set("machinedown", "Machine Down")
		downCause.Set("unknown", "Unknown")
		downCause.Set("weatherstop", "Weather Stop")

		tmpResult := []tk.M{}
		downDone := []string{}

		for f, t := range downCause {
			pipes = []tk.M{}
			loopMatch := match
			field := tk.ToString(f)
			title := tk.ToString(t)

			downDone = append(downDone, field)

			for _, done := range downDone {
				match.Unset("detail." + done)
			}

			loopMatch.Set("detail."+field, true)

			pipes = append(pipes, tk.M{"$unwind": "$detail"})
			pipes = append(pipes, tk.M{"$match": loopMatch})
			if topType == "loss" {
				pipes = append(pipes,
					tk.M{
						"$group": tk.M{"_id": tk.M{"id1": field, "id2": title}, "result": tk.M{"$sum": "$detail.powerlost"}},
					},
				)
			} else if topType == "duration" {
				pipes = append(pipes,
					tk.M{
						"$group": tk.M{"_id": tk.M{"id1": field, "id2": title}, "result": tk.M{"$sum": "$detail.duration"}},
					},
				)
			} else if topType == "frequency" {
				pipes = append(pipes,
					tk.M{
						"$group": tk.M{"_id": tk.M{"id1": field, "id2": title}, "result": tk.M{"$sum": 1}},
					},
				)
			}

			/*for _, v := range pipes {
				log.Printf("pipes: %#v \n", v)
			}*/

			csr, e := DB().Connection.NewQuery().
				From(new(Alarm).TableName()).
				Command("pipe", pipes).
				Cursor(nil)

			if e != nil {
				return result, e
			}

			resLoop := []tk.M{}
			e = csr.Fetch(&resLoop, 0, false)

			csr.Close()

			for _, res := range resLoop {
				tmpResult = append(tmpResult, res)
			}
		}

		size := len(tmpResult)

		if size > 1 {
			for i := 0; i < size; i++ {
				for j := size - 1; j >= i+1; j-- {
					a := tmpResult[j].GetFloat64("result")
					b := tmpResult[j-1].GetFloat64("result")

					if a > b {
						tmpResult[j], tmpResult[j-1] = tmpResult[j-1], tmpResult[j]
					}
				}
			}
		}

		result = tmpResult
	}

	return result, e
}

func getDownTimeTopFiltered(topType string, p *PayloadAnalytic, k *knot.WebContext) ([]tk.M, error) {
	var result []tk.M
	var e error
	var pipes []tk.M
	match := tk.M{}

	if p != nil {
		// tStart, _ := time.Parse("2006-01-02", p.DateStart.UTC().Format("2006-01-02"))
		// tEnd, _ := time.Parse("2006-01-02 15:04:05", p.DateEnd.UTC().Format("2006-01-02")+" 23:59:59")
		tStart, tEnd, e := helper.GetStartEndDate(k, p.Period, p.DateStart, p.DateEnd)
		if e != nil {
			return result, e
		}

		match.Set("detail.startdate", tk.M{"$gte": tStart, "$lte": tEnd})

		if p.Project != "" {
			match.Set("projectname", p.Project)
		}

		if len(p.Turbine) != 0 {
			match.Set("turbine", tk.M{"$in": p.Turbine})
		}

		pipes = append(pipes, tk.M{"$unwind": "$detail"})
		pipes = append(pipes, tk.M{"$match": match})
		if topType == "duration" {
			pipes = append(pipes, tk.M{"$group": tk.M{"_id": "$turbine", "result": tk.M{"$sum": "$detail.duration"}}})
		} else if topType == "frequency" {
			pipes = append(pipes, tk.M{"$group": tk.M{"_id": "$turbine", "result": tk.M{"$sum": 1}}})
		} else if topType == "loss" {
			pipes = append(pipes, tk.M{"$group": tk.M{"_id": "$turbine", "result": tk.M{"$sum": "$detail.powerlost"}}})
		}

		pipes = append(pipes, tk.M{"$sort": tk.M{"result": -1}})
		pipes = append(pipes, tk.M{"$limit": 10})

		/*log.Printf("date: %v | %v \n", tStart, tEnd)

		for _, v := range pipes {
			log.Printf("pipes: %#v \n", v)
		}*/

		// get the top 10
		csr, e := DB().Connection.NewQuery().
			//Select("_id").
			From(new(Alarm).TableName()).
			Command("pipe", pipes).
			Cursor(nil)

		if e != nil {
			return result, e
		}

		top10Turbines := []tk.M{}
		e = csr.Fetch(&top10Turbines, 0, false)

		// add by ams, 2016-10-07
		csr.Close()

		if e != nil {
			return result, e
		}

		// get the downtime
		turbines := []string{}
		turbinesVal := tk.M{}

		for _, turbine := range top10Turbines {
			turbines = append(turbines, turbine.Get("_id").(string))
			turbinesVal.Set(turbine.Get("_id").(string), turbine.GetFloat64("result"))
		}

		match.Set("turbine", tk.M{"$in": turbines})

		downCause := tk.M{}
		downCause.Set("aebok", "AEBOK")
		downCause.Set("externalstop", "External Stop")
		downCause.Set("griddown", "Grid Down")
		downCause.Set("internalgrid", "Internal Grid")
		downCause.Set("machinedown", "Machine Down")
		downCause.Set("unknown", "Unknown")
		downCause.Set("weatherstop", "Weather Stop")

		tmpResult := []tk.M{}
		downDone := []string{}

		for f, t := range downCause {
			pipes = []tk.M{}
			loopMatch := match
			field := tk.ToString(f)
			title := tk.ToString(t)

			downDone = append(downDone, field)

			for _, done := range downDone {
				match.Unset("detail." + done)
			}

			loopMatch.Set("detail."+field, true)

			pipes = append(pipes, tk.M{"$unwind": "$detail"})
			pipes = append(pipes, tk.M{"$match": loopMatch})
			if topType == "duration" {
				pipes = append(pipes,
					tk.M{
						"$group": tk.M{"_id": tk.M{"id1": "$detail.detaildateinfo.monthid", "id2": "$detail.detaildateinfo.monthdesc", "id3": "$turbine", "id4": title},
							"result": tk.M{"$sum": "$detail.duration"},
						},
					},
				)
			} else if topType == "frequency" {
				pipes = append(pipes,
					tk.M{
						"$group": tk.M{"_id": tk.M{"id1": "$detail.detaildateinfo.monthid", "id2": "$detail.detaildateinfo.monthdesc", "id3": "$turbine", "id4": title},
							"result": tk.M{"$sum": 1},
						},
					},
				)
			} else if topType == "loss" {
				pipes = append(pipes,
					tk.M{
						"$group": tk.M{"_id": tk.M{"id1": "$detail.detaildateinfo.monthid", "id2": "$detail.detaildateinfo.monthdesc", "id3": "$turbine", "id4": title},
							"result": tk.M{"$sum": "$detail.powerlost"},
						},
					},
				)
			}

			pipes = append(pipes, tk.M{"$sort": tk.M{"result": -1}})

			csr, e := DB().Connection.NewQuery().
				From(new(Alarm).TableName()).
				Command("pipe", pipes).
				Cursor(nil)

			if e != nil {
				return result, e
			}

			resLoop := []tk.M{}
			e = csr.Fetch(&resLoop, 0, false)

			// add by ams, 2016-10-07
			csr.Close()

			for _, res := range resLoop {
				tmpResult = append(tmpResult, res)
			}
		}

		resY := []tk.M{}

		for _, t := range downCause {
			title := tk.ToString(t)

			for _, turbine := range turbines {
				resX := tk.M{}
				resX.Set("_id", tk.M{"id3": turbine, "id4": title})
				resX.Set("result", 0)

			out:
				for _, res := range tmpResult {
					id3 := res.Get("_id").(tk.M).GetString("id3")
					id4 := res.Get("_id").(tk.M).GetString("id4")

					if id3 == turbine && id4 == title {
						resX = res
						break out
					}
				}

				// if title == "External Stop" {
				resY = append(resY, resX)
				// }
			}
		}

		for _, turbine := range turbines {
			resVal := tk.M{}
			resVal.Set("_id", turbine)

			for _, val := range resY {
				valTurbine := val.Get("_id").(tk.M).GetString("id3")
				valResult := val.GetFloat64("result")
				valTitle := ""

				splitTitle := strings.Split(val.Get("_id").(tk.M).GetString("id4"), " ")

				if len(splitTitle) > 1 {
					valTitle = splitTitle[0] + "" + splitTitle[1]
				} else {
					valTitle = splitTitle[0]
				}

				if turbine == valTurbine && valResult != 0 {
					resVal.Set(valTitle, valResult)
				} else if resVal.Get(valTitle) == nil {
					resVal.Set(valTitle, 0)
				}
			}

			resVal.Set("Total", turbinesVal.GetFloat64(turbine))
			result = append(result, resVal)
		}
	}

	return result, e
}

// func (m *AnalyticLossAnalysisController) GetHistogramData(k *knot.WebContext) interface{} {
// 	k.Config.OutputType = knot.OutputJson

// 	p := new(PayloadAnalytic)
// 	e := k.GetPayload(&p)

// 	if e != nil {
// 		return helper.CreateResult(false, nil, e.Error())
// 	}

// 	tStart, _ := time.Parse("2006-01-02", p.DateStart.UTC().Format("2006-01-02"))
// 	tEnd, _ := time.Parse("2006-01-02 15:04:05", p.DateEnd.UTC().Format("2006-01-02")+" 23:59:59")
// 	turbine := p.Turbine
// 	project := p.Project

// 	match := tk.M{}
// 	match.Set("dateinfo.dateid", tk.M{}.Set("$lte", tEnd).Set("$gte", tStart))
// 	match.Set("projectname", project)
// 	match.Set("avgwindspeed", tk.M{}.Set("$gte", 3).Set("$lt", 25))

// 	if len(turbine) > 0 {
// 		match.Set("turbine", tk.M{}.Set("$in", turbine))
// 	}

// 	group := tk.M{
// 		"_id":   "$wsavgforpc",
// 		"total": tk.M{}.Set("$sum", 1),
// 	}

// 	sort := tk.M{
// 		"_id": 1,
// 	}

// 	var pipes []tk.M
// 	pipes = append(pipes, tk.M{}.Set("$match", match))
// 	pipes = append(pipes, tk.M{}.Set("$group", group))
// 	pipes = append(pipes, tk.M{}.Set("$sort", sort))

// 	csr, e := DB().Connection.NewQuery().
// 		From(new(ScadaData).TableName()).
// 		Command("pipe", pipes).
// 		Cursor(nil)

// 	defer csr.Close()

// 	if e != nil {
// 		return helper.CreateResult(false, nil, "Error query : "+e.Error())
// 	}

// 	results := make([]tk.M, 0)
// 	e = csr.Fetch(&results, 0, false)

// 	if e != nil {
// 		return helper.CreateResult(false, nil, "Error facing results : "+e.Error())
// 	}

// 	totalData := c.From(&results).Sum(func(x interface{}) interface{} {
// 		dt := x.(tk.M)
// 		return dt["total"].(int)
// 	}).Exec().Result.Sum

// 	valuewindspeed := tk.M{"3.0": 0}
// 	valuewindspeed.Set("3.5", 0)

// 	categorywindspeed := []string{}
// 	categorywindspeed = append(categorywindspeed, "3 - 3.5")
// 	categorywindspeed = append(categorywindspeed, "3.5 - 4")
// 	for i := 4; i <= 24; i++ {
// 		nextPhase := i + 1
// 		categorywindspeed = append(categorywindspeed, strconv.Itoa(i)+" - "+strconv.Itoa(nextPhase))
// 		valuewindspeed.Set(strconv.Itoa(i)+".0", 0)
// 	}

// 	for _, x := range results {
// 		id := tk.RoundingAuto64(x["_id"].(float64), 1)
// 		total := x["total"].(int)
// 		value := tk.Div(float64(total), totalData)

// 		sId := strconv.FormatFloat(id, 'f', 1, 64)

// 		valuewindspeed.Set(sId, value)
// 	}

// 	retvaluews := []float64{}
// 	retvaluews = append(retvaluews, valuewindspeed.GetFloat64("3.0"))
// 	retvaluews = append(retvaluews, valuewindspeed.GetFloat64("3.5"))
// 	for i := 4; i <= 24; i++ {
// 		retvaluews = append(retvaluews, valuewindspeed.GetFloat64(strconv.Itoa(i)+".0"))
// 	}

// 	data := tk.M{
// 		"categorywindspeed": categorywindspeed,
// 		"valuewindspeed":    retvaluews,
// 		"totaldata":         totalData,
// 	}

// 	return helper.CreateResult(true, data, "success")
// }

func (m *AnalyticLossAnalysisController) GetHistogramProduction(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	p := new(PayloadAnalytic)
	e := k.GetPayload(&p)

	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	tStart, _ := time.Parse("2006-01-02", p.DateStart.UTC().Format("2006-01-02"))
	tEnd, _ := time.Parse("2006-01-02 15:04:05", p.DateEnd.UTC().Format("2006-01-02")+" 23:59:59")
	turbine := p.Turbine
	project := p.Project

	match := tk.M{}
	match.Set("dateinfo.dateid", tk.M{}.Set("$lte", tEnd).Set("$gte", tStart))
	match.Set("projectname", project)
	match.Set("avgwindspeed", tk.M{}.Set("$gte", 3).Set("$lt", 25))

	if len(turbine) > 0 {
		match.Set("turbine", tk.M{}.Set("$in", turbine))
	}

	group := tk.M{
		"_id":   "$wsavgforpc",
		"total": tk.M{}.Set("$sum", 1),
	}

	sort := tk.M{
		"_id": 1,
	}

	var pipes []tk.M
	pipes = append(pipes, tk.M{}.Set("$match", match))
	pipes = append(pipes, tk.M{}.Set("$group", group))
	pipes = append(pipes, tk.M{}.Set("$sort", sort))

	csr, e := DB().Connection.NewQuery().
		From(new(ScadaData).TableName()).
		Command("pipe", pipes).
		Cursor(nil)

	defer csr.Close()

	if e != nil {
		return helper.CreateResult(false, nil, "Error query : "+e.Error())
	}

	results := make([]tk.M, 0)
	e = csr.Fetch(&results, 0, false)

	if e != nil {
		return helper.CreateResult(false, nil, "Error facing results : "+e.Error())
	}

	totalData := c.From(&results).Sum(func(x interface{}) interface{} {
		dt := x.(tk.M)
		return dt["total"].(int)
	}).Exec().Result.Sum

	valuewindspeed := tk.M{"3.0": 0}
	valuewindspeed.Set("3.5", 0)

	categorywindspeed := []string{}
	categorywindspeed = append(categorywindspeed, "3 - 3.5")
	categorywindspeed = append(categorywindspeed, "3.5 - 4")
	for i := 4; i <= 24; i++ {
		nextPhase := i + 1
		categorywindspeed = append(categorywindspeed, strconv.Itoa(i)+" - "+strconv.Itoa(nextPhase))
		valuewindspeed.Set(strconv.Itoa(i)+".0", 0)
	}

	for _, x := range results {
		id := tk.RoundingAuto64(x["_id"].(float64), 1)
		total := x["total"].(int)
		value := tk.Div(float64(total), totalData)

		sId := strconv.FormatFloat(id, 'f', 1, 64)

		valuewindspeed.Set(sId, value)
	}

	retvaluews := []float64{}
	retvaluews = append(retvaluews, valuewindspeed.GetFloat64("3.0"))
	retvaluews = append(retvaluews, valuewindspeed.GetFloat64("3.5"))
	for i := 4; i <= 24; i++ {
		retvaluews = append(retvaluews, valuewindspeed.GetFloat64(strconv.Itoa(i)+".0"))
	}

	data := tk.M{
		"categorywindspeed": categorywindspeed,
		"valuewindspeed":    retvaluews,
		"totaldata":         totalData,
	}

	return helper.CreateResult(true, data, "success")
}

type PayloadHistogram struct {
	MaxValue float64
	MinValue float64
	BinValue int
	Filter   PayloadAnalytic
}

func (m *AnalyticLossAnalysisController) GetHistogramData(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	p := new(PayloadHistogram)
	e := k.GetPayload(&p)

	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	/*tStart, _ := time.Parse("2006-01-02", p.Filter.DateStart.UTC().Format("2006-01-02"))
	tEnd, _ := time.Parse("2006-01-02 15:04:05", p.Filter.DateEnd.UTC().Format("2006-01-02")+" 23:59:59")*/
	tStart, tEnd, e := helper.GetStartEndDate(k, p.Filter.Period, p.Filter.DateStart, p.Filter.DateEnd)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}
	turbine := p.Filter.Turbine
	project := p.Filter.Project

	categorywindspeed := []string{}
	valuewindspeed := []float64{}
	interval := (p.MaxValue - p.MinValue) / float64(p.BinValue)
	startcategory := p.MinValue
	totalData := 0.0

	for i := 0; i < (p.BinValue); i++ {
		categorywindspeed = append(categorywindspeed, fmt.Sprintf("%.0f", startcategory)+" ~ "+fmt.Sprintf("%.0f", (startcategory+interval)))

		match := tk.M{}
		match.Set("avgwindspeed", tk.M{}.Set("$lt", (startcategory+interval)).Set("$gte", startcategory))
		match.Set("dateinfo.dateid", tk.M{}.Set("$lte", tEnd).Set("$gte", tStart))
		if len(project) > 0 {
			match.Set("projectname", project)
		}
		if len(turbine) > 0 {
			match.Set("turbine", tk.M{}.Set("$in", turbine))
		}

		group := tk.M{
			"_id":   "",
			"total": tk.M{}.Set("$sum", 1),
		}

		var pipes []tk.M
		pipes = append(pipes, tk.M{}.Set("$match", match))
		pipes = append(pipes, tk.M{}.Set("$group", group))

		csr, e := DB().Connection.NewQuery().
			From(new(ScadaData).TableName()).
			Command("pipe", pipes).
			Cursor(nil)

		defer csr.Close()

		if e != nil {
			return helper.CreateResult(false, nil, "Error query : "+e.Error())
		}

		resultCategory := []tk.M{}
		e = csr.Fetch(&resultCategory, 0, false)

		if len(resultCategory) > 0 {
			valuewindspeed = append(valuewindspeed, float64(resultCategory[0]["total"].(int)))
			totalData = totalData + float64(resultCategory[0]["total"].(int))
		} else {
			valuewindspeed = append(valuewindspeed, 0.00)
		}

		startcategory = startcategory + interval
	}

	for i := 0; i < len(valuewindspeed); i++ {
		valuewindspeed[i] = float64(int((valuewindspeed[i]/totalData*100)*100)) / 100
	}

	data := tk.M{
		"categorywindspeed": categorywindspeed,
		"valuewindspeed":    valuewindspeed,
		"totaldata":         totalData,
	}

	return helper.CreateResult(true, data, "success")
}

func (m *AnalyticLossAnalysisController) GetProductionHistogramData(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	p := new(PayloadHistogram)
	e := k.GetPayload(&p)

	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	/*tStart, _ := time.Parse("2006-01-02", p.Filter.DateStart.UTC().Format("2006-01-02"))
	tEnd, _ := time.Parse("2006-01-02 15:04:05", p.Filter.DateEnd.UTC().Format("2006-01-02")+" 23:59:59")*/
	tStart, tEnd, e := helper.GetStartEndDate(k, p.Filter.Period, p.Filter.DateStart, p.Filter.DateEnd)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}
	turbine := p.Filter.Turbine
	project := p.Filter.Project

	categoryproduction := []string{}
	valueproduction := []float64{}
	interval := (p.MaxValue - p.MinValue) / float64(p.BinValue)
	startcategory := p.MinValue
	totalData := 0.0

	for i := 0; i < (p.BinValue); i++ {
		categoryproduction = append(categoryproduction, fmt.Sprintf("%.0f", startcategory)+" ~ "+fmt.Sprintf("%.0f", (startcategory+interval)))

		match := tk.M{}
		match.Set("power", tk.M{}.Set("$lt", (startcategory+interval)).Set("$gte", startcategory))
		match.Set("dateinfo.dateid", tk.M{}.Set("$lte", tEnd).Set("$gte", tStart))
		if len(project) > 0 {
			match.Set("projectname", project)
		}
		if len(turbine) > 0 {
			match.Set("turbine", tk.M{}.Set("$in", turbine))
		}

		group := tk.M{
			"_id":   "",
			"total": tk.M{}.Set("$sum", 1),
		}

		var pipes []tk.M
		pipes = append(pipes, tk.M{}.Set("$match", match))
		pipes = append(pipes, tk.M{}.Set("$group", group))

		csr, e := DB().Connection.NewQuery().
			From(new(ScadaData).TableName()).
			Command("pipe", pipes).
			Cursor(nil)

		defer csr.Close()

		if e != nil {
			return helper.CreateResult(false, nil, "Error query : "+e.Error())
		}

		resultCategory := []tk.M{}
		e = csr.Fetch(&resultCategory, 0, false)

		if len(resultCategory) > 0 {
			valueproduction = append(valueproduction, float64(resultCategory[0]["total"].(int)))
			totalData = totalData + float64(resultCategory[0]["total"].(int))
		} else {
			valueproduction = append(valueproduction, 0.00)
		}

		startcategory = startcategory + interval
	}

	for i := 0; i < len(valueproduction); i++ {
		valueproduction[i] = float64(int((valueproduction[i]/totalData*100)*100)) / 100
	}

	data := tk.M{
		"categoryproduction": categoryproduction,
		"valueproduction":    valueproduction,
		"totaldata":          totalData,
	}

	return helper.CreateResult(true, data, "success")
}

func (m *AnalyticLossAnalysisController) GetAvailDate(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	return k.Session("availdate", "")
}
