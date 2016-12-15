package controller

import (
	. "eaciit/wfdemo-git/library/core"
	. "eaciit/wfdemo-git/library/models"
	"eaciit/wfdemo-git/web/helper"
	"strings"

	// "fmt"
	"github.com/eaciit/dbox"
	"github.com/eaciit/knot/knot.v1"
	tk "github.com/eaciit/toolkit"
	"time"
)

type DataBrowserController struct {
	App
}

func CreateDataBrowserController() *DataBrowserController {
	var controller = new(DataBrowserController)
	return controller
}

func (m *DataBrowserController) GetScadaList(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson
	var filter []*dbox.Filter

	p := new(helper.PayloadsDB)
	e := k.GetPayload(&p)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}
	// filter, _ := p.ParseFilter()

	tStart, _ := time.Parse("2006-01-02", p.DateStart.UTC().Format("2006-01-02"))
	tEnd, _ := time.Parse("2006-01-02 15:04:05", p.DateEnd.UTC().Format("2006-01-02")+" 23:59:59")
	turbine := p.Turbine


	filter = append(filter, dbox.Ne("_id", ""))
	filter = append(filter, dbox.Gte("timestamp", tStart))
	filter = append(filter, dbox.Lte("timestamp", tEnd))
	if len(turbine) != 0 {
		filter = append(filter, dbox.In("turbine", turbine...))
	}


	query := DB().Connection.NewQuery().From(new(ScadaData).TableName()).Skip(p.Skip).Take(p.Take)
	query.Where(dbox.And(filter...))

	if len(p.Sort) > 0 {
		var arrsort []string
		for _, val := range p.Sort {
			if val.Dir == "desc" {
				arrsort = append(arrsort, strings.ToLower("-"+val.Field))
			} else {
				arrsort = append(arrsort, strings.ToLower(val.Field))
			}
		}
		query = query.Order(arrsort...)
	}
	csr, e := query.Cursor(nil)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}
	defer csr.Close()

	tmpResult := make([]ScadaData, 0)
	results := make([]ScadaData, 0)
	e = csr.Fetch(&tmpResult, 0, false)

	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	for _, val := range tmpResult {
		val.TimeStamp = val.TimeStamp.UTC()
		results = append(results, val)
	}

	// queryC := DB().Connection.NewQuery().From(new(ScadaData).TableName()).Where(dbox.And(filter...))
	// ccount, e := queryC.Cursor(nil)
	// if e != nil {
	// 	return helper.CreateResult(false, nil, e.Error())
	// }
	// defer ccount.Close()

	totalPower := 0.0
	totalPowerLost := 0.0
	totalTurbine := 0
	totalProduction := 0.0
	sumWindSpeed := 0.0
	countData := 0.0

	aggrData := []tk.M{}

	queryAggr := DB().Connection.NewQuery().From(new(ScadaData).TableName()).
		Aggr(dbox.AggrSum, "$power", "TotalPower").
		Aggr(dbox.AggrSum, "$powerlost", "TotalPowerLost").
		Aggr(dbox.AggrSum, "$power", "totalProduction").
		Aggr(dbox.AggrSum, "$avgwindspeed", "sumWindSpeed").
		Aggr(dbox.AggrSum, 1, "countData").
		Group("turbine").Where(dbox.And(filter...))

	caggr, e := queryAggr.Cursor(nil)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}
	defer caggr.Close()
	e = caggr.Fetch(&aggrData, 0, false)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	for _, val := range aggrData {
		totalPower += val.GetFloat64("TotalPower")
		totalPowerLost += val.GetFloat64("TotalPowerLost")
		totalProduction += val.GetFloat64("totalProduction")
		sumWindSpeed += val.GetFloat64("sumWindSpeed")
		countData += val.GetFloat64("countData")
	}
	totalTurbine = tk.SliceLen(aggrData)

	data := struct {
		Data            []ScadaData
		Total           float64
		TotalPower      float64
		TotalPowerLost  float64
		TotalProduction float64
		AvgWindSpeed    float64
		TotalTurbine    int
	}{
		Data:            results,
		Total:           countData,
		TotalPower:      totalPower,
		TotalPowerLost:  totalPowerLost,
		TotalProduction: totalProduction,// / 6,
		AvgWindSpeed:    sumWindSpeed / countData,
		TotalTurbine:    totalTurbine,
	}

	return helper.CreateResult(true, data, "success")
}

func (m *DataBrowserController) GetScadaAnomalyList(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	p := new(helper.Payloads)
	e := k.GetPayload(&p)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	filter, _ := p.ParseFilter()

	query := DB().Connection.NewQuery().From(new(ScadaAlarmAnomaly).TableName()).Skip(p.Skip).Take(p.Take)
	query.Where(dbox.And(filter...))

	if len(p.Sort) > 0 {
		var arrsort []string
		for _, val := range p.Sort {
			if val.Dir == "desc" {
				arrsort = append(arrsort, strings.ToLower("-"+val.Field))
			} else {
				arrsort = append(arrsort, strings.ToLower(val.Field))
			}
		}
		query = query.Order(arrsort...)
	}
	csr, e := query.Cursor(nil)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}
	defer csr.Close()

	tmpResult := make([]ScadaAlarmAnomaly, 0)
	results := make([]ScadaAlarmAnomaly, 0)
	e = csr.Fetch(&tmpResult, 0, false)

	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	for _, val := range tmpResult {
		val.TimeStamp = val.TimeStamp.UTC()
		results = append(results, val)
	}

	queryC := DB().Connection.NewQuery().From(new(ScadaAlarmAnomaly).TableName()).Where(dbox.And(filter...))
	ccount, e := queryC.Cursor(nil)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}
	defer ccount.Close()

	totalPower := 0.0
	totalPowerLost := 0.0
	// totalProduction := 0.0
	sumWindSpeed := 0.0
	totalTurbine := 0

	aggrData := []tk.M{}

	queryAggr := DB().Connection.NewQuery().From(new(ScadaAlarmAnomaly).TableName()).
		Aggr(dbox.AggrSum, "$power", "TotalPower").
		Aggr(dbox.AggrSum, "$powerlost", "TotalPowerLost").
		Aggr(dbox.AggrSum, "$avgwindspeed", "sumWindSpeed").
		Group("turbine").Where(dbox.And(filter...))

	caggr, e := queryAggr.Cursor(nil)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}
	defer caggr.Close()
	e = caggr.Fetch(&aggrData, 0, false)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	for _, val := range aggrData {
		totalPower += val.GetFloat64("TotalPower")
		totalPowerLost += val.GetFloat64("TotalPowerLost")
		sumWindSpeed += val.GetFloat64("sumWindSpeed")
	}
	totalTurbine = tk.SliceLen(aggrData)

	data := struct {
		Data            []ScadaAlarmAnomaly
		Total           int
		TotalPower      float64
		TotalPowerLost  float64
		TotalProduction float64
		AvgWindSpeed    float64
		TotalTurbine    int
	}{
		Data:            results,
		Total:           ccount.Count(),
		TotalPower:      totalPower,
		TotalPowerLost:  totalPowerLost,
		TotalProduction: totalPower / 6,
		AvgWindSpeed:    sumWindSpeed / float64(ccount.Count()),
		TotalTurbine:    totalTurbine,
	}

	return helper.CreateResult(true, data, "success")
}

func (m *DataBrowserController) GetAlarmList(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	p := new(helper.Payloads)

	e := k.GetPayload(&p)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	filter, _ := p.ParseFilter()
	query := DB().Connection.NewQuery().From(new(Alarm).TableName()).
		Skip(p.Skip).Take(p.Take).Where(dbox.And(filter...))

	if len(p.Sort) > 0 {
		var arrsort []string
		for _, val := range p.Sort {
			if val.Dir == "desc" {
				arrsort = append(arrsort, strings.ToLower("-"+val.Field))
			} else {
				arrsort = append(arrsort, strings.ToLower(val.Field))
			}
		}
		query = query.Order(arrsort...)
	}

	csr, e := query.Cursor(nil)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}
	defer csr.Close()

	tmpResult := make([]Alarm, 0)
	results := make([]Alarm, 0)
	e = csr.Fetch(&tmpResult, 0, false)

	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	for _, val := range tmpResult {
		val.StartDate = val.StartDate.UTC()
		val.EndDate = val.EndDate.UTC()
		results = append(results, val)
	}

	queryC := DB().Connection.NewQuery().From(new(Alarm).TableName()).Where(dbox.And(filter...))

	ccount, e := queryC.Cursor(nil)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}
	defer ccount.Close()

	totalTurbine := 0
	aggrData := []tk.M{}

	queryAggr := DB().Connection.NewQuery().From(new(Alarm).TableName()).
		Group("turbine").Where(dbox.And(filter...))

	caggr, e := queryAggr.Cursor(nil)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}
	defer caggr.Close()
	e = caggr.Fetch(&aggrData, 0, false)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	totalTurbine = tk.SliceLen(aggrData)

	data := struct {
		Data         []Alarm
		Total        int
		TotalTurbine int
	}{
		Data:         results,
		Total:        ccount.Count(),
		TotalTurbine: totalTurbine,
	}

	return helper.CreateResult(true, data, "success")
}

func (m *DataBrowserController) GetAlarmScadaAnomalyList(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	p := new(helper.Payloads)

	e := k.GetPayload(&p)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	filter, _ := p.ParseFilter()
	query := DB().Connection.NewQuery().From(new(AlarmScadaAnomaly).TableName()).
		Skip(p.Skip).Take(p.Take).Where(dbox.And(filter...))

	if len(p.Sort) > 0 {
		var arrsort []string
		for _, val := range p.Sort {
			if val.Dir == "desc" {
				arrsort = append(arrsort, strings.ToLower("-"+val.Field))
			} else {
				arrsort = append(arrsort, strings.ToLower(val.Field))
			}
		}
		query = query.Order(arrsort...)
	}

	csr, e := query.Cursor(nil)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}
	defer csr.Close()

	tmpResult := make([]AlarmScadaAnomaly, 0)
	results := make([]AlarmScadaAnomaly, 0)
	e = csr.Fetch(&tmpResult, 0, false)

	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	for _, val := range tmpResult {
		val.StartDate = val.StartDate.UTC()
		val.EndDate = val.EndDate.UTC()
		results = append(results, val)
	}

	queryC := DB().Connection.NewQuery().From(new(AlarmScadaAnomaly).TableName()).
		Where(dbox.And(filter...))

	ccount, e := queryC.Cursor(nil)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}
	defer ccount.Close()

	totalTurbine := 0
	aggrData := []tk.M{}

	queryAggr := DB().Connection.NewQuery().From(new(Alarm).TableName()).
		Group("turbine").Where(dbox.And(filter...))

	caggr, e := queryAggr.Cursor(nil)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}
	defer caggr.Close()
	e = caggr.Fetch(&aggrData, 0, false)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	totalTurbine = tk.SliceLen(aggrData)

	data := struct {
		Data         []AlarmScadaAnomaly
		Total        int
		TotalTurbine int
	}{
		Data:         results,
		Total:        ccount.Count(),
		TotalTurbine: totalTurbine,
	}

	return helper.CreateResult(true, data, "success")
}

func (m *DataBrowserController) GetAlarmOverlappingList(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	p := new(helper.Payloads)

	e := k.GetPayload(&p)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	filter, _ := p.ParseFilter()
	query := DB().Connection.NewQuery().From(new(AlarmOverlapping).TableName()).
		Skip(p.Skip).Take(p.Take).Where(dbox.And(filter...))

	if len(p.Sort) > 0 {
		var arrsort []string
		for _, val := range p.Sort {
			if val.Dir == "desc" {
				arrsort = append(arrsort, strings.ToLower("-"+val.Field))
			} else {
				arrsort = append(arrsort, strings.ToLower(val.Field))
			}
		}
		query = query.Order(arrsort...)
	}

	csr, e := query.Cursor(nil)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}
	defer csr.Close()

	tmpResult := make([]AlarmOverlapping, 0)
	results := make([]AlarmOverlapping, 0)
	e = csr.Fetch(&tmpResult, 0, false)

	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	for _, val := range tmpResult {
		val.StartDate = val.StartDate.UTC()
		val.EndDate = val.EndDate.UTC()
		results = append(results, val)
	}

	queryC := DB().Connection.NewQuery().From(new(AlarmOverlapping).TableName()).
		Where(dbox.And(filter...))

	ccount, e := queryC.Cursor(nil)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}
	defer ccount.Close()

	totalTurbine := 0
	aggrData := []tk.M{}

	queryAggr := DB().Connection.NewQuery().From(new(AlarmOverlapping).TableName()).
		Group("turbine").Where(dbox.And(filter...))

	caggr, e := queryAggr.Cursor(nil)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}
	defer caggr.Close()
	e = caggr.Fetch(&aggrData, 0, false)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	totalTurbine = tk.SliceLen(aggrData)

	data := struct {
		Data         []AlarmOverlapping
		Total        int
		TotalTurbine int
	}{
		Data:         results,
		Total:        ccount.Count(),
		TotalTurbine: totalTurbine,
	}

	return helper.CreateResult(true, data, "success")
}

func (m *DataBrowserController) GetAlarmOverlappingDetails(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	p := new(helper.Payloads)

	e := k.GetPayload(&p)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	filter, _ := p.ParseFilter()

	var (
		pipes []tk.M
	)

	pipes = append(pipes, tk.M{"$unwind": "$alarms"})
	pipes = append(pipes, tk.M{"$project": tk.M{
		"farm":             "$alarms.farm",
		"startdate":        "$alarms.startdate",
		"enddate":          "$alarms.enddate",
		"turbine":          "$alarms.turbine",
		"alertdescription": "$alarms.alertdescription",
		"externalstop":     "$alarms.externalstop",
		"griddown":         "$alarms.griddown",
		"internalgrid":     "$alarms.internalgrid",
		"machinedown":      "$alarms.machinedown",
		"aebok":            "$alarms.aebok",
		"unknown":          "$alarms.unknown",
		"weatherstop":      "$alarms.weatherstop",
		"line":             "$alarms.line",
	}})

	query := DB().Connection.NewQuery().From(new(AlarmOverlapping).TableName()).
		Command("pipe", pipes).Where(dbox.And(filter...))

	if len(p.Sort) > 0 {
		var arrsort []string
		for _, val := range p.Sort {
			if val.Dir == "desc" {
				arrsort = append(arrsort, strings.ToLower("-"+val.Field))
			} else {
				arrsort = append(arrsort, strings.ToLower(val.Field))
			}
		}
		query = query.Order(arrsort...)
	}
	queryCount := query

	query.Skip(p.Skip).Take(p.Take)
	csr, e := query.Cursor(nil)
	if e != nil {
		return e.Error()
	}
	defer csr.Close()

	tmpResult := make([]Alarm, 0)
	results := make([]Alarm, 0)
	e = csr.Fetch(&tmpResult, 0, false)

	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	for _, val := range tmpResult {
		val.StartDate = val.StartDate.UTC()
		val.EndDate = val.EndDate.UTC()
		results = append(results, val)
	}

	ccount, e := queryCount.Cursor(nil)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}
	defer ccount.Close()

	data := struct {
		Data  []Alarm
		Total int
	}{
		Data:  results,
		Total: ccount.Count(),
	}

	return helper.CreateResult(true, data, "success")
}

func (m *DataBrowserController) GetJMRList(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	p := new(helper.Payloads)
	e := k.GetPayload(&p)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}
	filter, _ := p.ParseFilter()

	query := DB().Connection.NewQuery().From(new(JMR).TableName()).Skip(p.Skip).Take(p.Take)
	query.Where(dbox.And(filter...))

	if len(p.Sort) > 0 {
		var arrsort []string
		for _, val := range p.Sort {
			if val.Dir == "desc" {
				arrsort = append(arrsort, strings.ToLower("-"+val.Field))
			} else {
				arrsort = append(arrsort, strings.ToLower(val.Field))
			}
		}
		query = query.Order(arrsort...)
	}
	csr, e := query.Cursor(nil)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}
	defer csr.Close()

	tmpResult := make([]JMR, 0)
	results := make([]JMR, 0)
	e = csr.Fetch(&tmpResult, 0, false)

	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	for _, val := range tmpResult {
		val.Sections = nil
		results = append(results, val)
	}

	queryC := DB().Connection.NewQuery().From(new(JMR).TableName()).Where(dbox.And(filter...))
	ccount, e := queryC.Cursor(nil)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}
	defer ccount.Close()

	data := struct {
		Data  []JMR
		Total int
	}{
		Data:  results,
		Total: ccount.Count(),
	}

	return helper.CreateResult(true, data, "success")
}

func (m *DataBrowserController) GetJMRDetails(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	p := new(helper.Payloads)

	e := k.GetPayload(&p)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	filter, _ := p.ParseFilter()

	query := DB().Connection.NewQuery().From(new(JMR).TableName())
	query.Where(dbox.And(filter...))
	csr, e := query.Cursor(nil)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	jmrResult := make([]JMR, 0)
	e = csr.Fetch(&jmrResult, 0, false)

	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}
	defer csr.Close()

	result := make([]JMRSection, 0)

	turbines := make(map[string]bool, 0)

	for _, fil := range filter {
		if fil.Field == "sections.turbine" {
			for _, str := range fil.Value.([]interface{}) {
				turbines[tk.ToString(str)] = true
			}
		}
	}

	clean := tk.M{}

	for _, jmr := range jmrResult {
		sectionsClean := []JMRSection{}
		for _, section := range jmr.Sections {
			if turbines[section.Turbine] {
				sectionsClean = append(sectionsClean, section)
			}
		}

		if len(sectionsClean) > 0 {
			clean.Set(jmr.DateInfo.MonthDesc, sectionsClean)
		}
	}

	for _, jmr := range jmrResult {
		for _, total := range jmr.TotalDetails {

			var contrGenTotal float64
			var boEExportTotal float64
			var boEImportTotal float64
			var boENetTotal float64

			sectionsClean := clean.Get(jmr.DateInfo.MonthDesc).([]JMRSection)

			for _, section := range sectionsClean {
				if total.Section == section.Description {
					result = append(result, section)
					contrGenTotal += section.ContrGen
					boEExportTotal += section.BoEExport
					boEImportTotal += section.BoEImport
					boENetTotal += section.BoENet
				}
			}

			if contrGenTotal != 0 {
				tmpSection := JMRSection{}
				tmpSection.Company = "Total"
				tmpSection.ContrGen = contrGenTotal
				tmpSection.BoEExport = boEExportTotal
				tmpSection.BoEImport = boEImportTotal
				tmpSection.BoENet = boENetTotal

				result = append(result, tmpSection)
			}

		}
	}

	return helper.CreateResult(true, result, "success")
}

func (m *DataBrowserController) GetMETList(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	p := new(helper.Payloads)
	e := k.GetPayload(&p)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}
	filter, _ := p.ParseFilter()

	query := DB().Connection.NewQuery().From(new(MetTower).TableName()).Skip(p.Skip).Take(p.Take)
	query.Where(dbox.And(filter...))

	if len(p.Sort) > 0 {
		var arrsort []string
		for _, val := range p.Sort {
			if val.Dir == "desc" {
				arrsort = append(arrsort, strings.ToLower("-"+val.Field))
			} else {
				arrsort = append(arrsort, strings.ToLower(val.Field))
			}
		}
		query = query.Order(arrsort...)
	}
	csr, e := query.Cursor(nil)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}
	defer csr.Close()

	tmpResult := make([]MetTower, 0)
	results := make([]MetTower, 0)
	e = csr.Fetch(&tmpResult, 0, false)

	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	for _, val := range tmpResult {
		val.TimeStamp = val.TimeStamp.UTC()
		results = append(results, val)
	}

	queryC := DB().Connection.NewQuery().From(new(MetTower).TableName()).Where(dbox.And(filter...))
	ccount, e := queryC.Cursor(nil)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}
	defer ccount.Close()

	data := struct {
		Data  []MetTower
		Total int
	}{
		Data:  results,
		Total: ccount.Count(),
	}

	return helper.CreateResult(true, data, "success")
}


func (m *DataBrowserController) GetEventList(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson
	var filter []*dbox.Filter

	p := new(helper.PayloadsDB)
	e := k.GetPayload(&p)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	tStart, _ := time.Parse("2006-01-02", p.DateStart.UTC().Format("2006-01-02"))
	tEnd, _ := time.Parse("2006-01-02 15:04:05", p.DateEnd.UTC().Format("2006-01-02")+" 23:59:59")
	turbine := p.Turbine

	filter = append(filter, dbox.Ne("_id", ""))
	filter = append(filter, dbox.Gte("timestamp", tStart))
	filter = append(filter, dbox.Lte("timestamp", tEnd))
	if len(turbine) != 0 {
		filter = append(filter, dbox.In("turbine", turbine...))
	}

	query := DB().Connection.NewQuery().From(new(EventRaw).TableName()).Skip(p.Skip).Take(p.Take)
	query.Where(dbox.And(filter...))

	if len(p.Sort) > 0 {
		var arrsort []string
		for _, val := range p.Sort {
			if val.Dir == "desc" {
				arrsort = append(arrsort, strings.ToLower("-"+val.Field))
			} else {
				arrsort = append(arrsort, strings.ToLower(val.Field))
			}
		}
		query = query.Order(arrsort...)
	}
	csr, e := query.Cursor(nil)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}
	defer csr.Close()

	tmpResult := make([]EventRaw, 0)
	results := make([]EventRaw, 0)
	e = csr.Fetch(&tmpResult, 0, false)

	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	for _, val := range tmpResult {
		val.TimeStamp = val.TimeStamp.UTC()
		results = append(results, val)
	}

	totalTurbine := 0
	countData := 0.0

	aggrData := []tk.M{}

	queryAggr := DB().Connection.NewQuery().From(new(EventRaw).TableName()).
		Aggr(dbox.AggrSum, 1, "countData").
		Group("turbine").Where(dbox.And(filter...))

	caggr, e := queryAggr.Cursor(nil)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}
	defer caggr.Close()
	e = caggr.Fetch(&aggrData, 0, false)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	for _, val := range aggrData {
		countData += val.GetFloat64("countData")
	}
	totalTurbine = tk.SliceLen(aggrData)

	data := struct {
		Data            []EventRaw
		Total           float64
		TotalTurbine    int
	}{
		Data:            results,
		Total:           countData,
		TotalTurbine:    totalTurbine,
	}

	return helper.CreateResult(true, data, "success")
}


// Get date info each tab

func (m *DataBrowserController) GetAvailDate(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	Scadaresults := make([]time.Time, 0)
	Alarmresults := make([]time.Time, 0)
	JMRresults := make([]time.Time, 0)
	METresults := make([]time.Time, 0)
	Durationresults := make([]time.Time, 0)
	ScadaAnomalyresults := make([]time.Time, 0)
	AlarmOverlappingresults := make([]time.Time, 0)
	AlarmScadaAnomalyresults := make([]time.Time, 0)

	// Scada Data
	for i := 0; i < 2; i++ {
		var arrsort []string
		if i == 0 {
			arrsort = append(arrsort, "timestamp")
		} else {
			arrsort = append(arrsort, "-timestamp")
		}

		query := DB().Connection.NewQuery().From(new(ScadaData).TableName()).Skip(0).Take(1)
		query = query.Order(arrsort...)

		csr, e := query.Cursor(nil)
		if e != nil {
			return helper.CreateResult(false, nil, e.Error())
		}
		defer csr.Close()

		Result := make([]ScadaData, 0)
		e = csr.Fetch(&Result, 0, false)

		if e != nil {
			return helper.CreateResult(false, nil, e.Error())
		}

		for _, val := range Result {
			Scadaresults = append(Scadaresults, val.TimeStamp.UTC())
		}
	}

	// Alarm Data
	for i := 0; i < 2; i++ {
		var arrsort []string
		if i == 0 {
			arrsort = append(arrsort, "startdate")
		} else {
			arrsort = append(arrsort, "-startdate")
		}

		query := DB().Connection.NewQuery().From(new(Alarm).TableName()).Skip(0).Take(1)
		query = query.Order(arrsort...)

		csr, e := query.Cursor(nil)
		if e != nil {
			return helper.CreateResult(false, nil, e.Error())
		}
		defer csr.Close()

		Result := make([]Alarm, 0)
		e = csr.Fetch(&Result, 0, false)

		if e != nil {
			return helper.CreateResult(false, nil, e.Error())
		}

		for _, val := range Result {
			Alarmresults = append(Alarmresults, val.StartDate.UTC())
		}
	}

	// JMR Data
	for i := 0; i < 2; i++ {
		var arrsort []string
		if i == 0 {
			arrsort = append(arrsort, "dateinfo.dateid")
		} else {
			arrsort = append(arrsort, "-dateinfo.dateid")
		}

		query := DB().Connection.NewQuery().From(new(ScadaData).TableName()).Skip(0).Take(1)
		query = query.Order(arrsort...)

		csr, e := query.Cursor(nil)
		if e != nil {
			return helper.CreateResult(false, nil, e.Error())
		}
		defer csr.Close()

		Result := make([]ScadaData, 0)
		e = csr.Fetch(&Result, 0, false)

		if e != nil {
			return helper.CreateResult(false, nil, e.Error())
		}

		for _, val := range Result {
			JMRresults = append(JMRresults, val.DateInfo.DateId.UTC())
		}
	}

	// MET Data
	for i := 0; i < 2; i++ {
		var arrsort []string
		if i == 0 {
			arrsort = append(arrsort, "timestamp")
		} else {
			arrsort = append(arrsort, "-timestamp")
		}

		query := DB().Connection.NewQuery().From(new(MetTower).TableName()).Skip(0).Take(1)
		query = query.Order(arrsort...)

		csr, e := query.Cursor(nil)
		if e != nil {
			return helper.CreateResult(false, nil, e.Error())
		}
		defer csr.Close()

		Result := make([]MetTower, 0)
		e = csr.Fetch(&Result, 0, false)

		if e != nil {
			return helper.CreateResult(false, nil, e.Error())
		}

		for _, val := range Result {
			METresults = append(METresults, val.TimeStamp.UTC())
		}
	}

	// Duration Data
	for i := 0; i < 2; i++ {
		var arrsort []string
		if i == 0 {
			arrsort = append(arrsort, "timestamp")
		} else {
			arrsort = append(arrsort, "-timestamp")
		}

		query := DB().Connection.NewQuery().From(new(ScadaData).TableName()).Where(dbox.And(dbox.Eq("isvalidtimeduration", false))).Skip(0).Take(1)
		query = query.Order(arrsort...)

		csr, e := query.Cursor(nil)
		if e != nil {
			return helper.CreateResult(false, nil, e.Error())
		}
		defer csr.Close()

		Result := make([]ScadaData, 0)
		e = csr.Fetch(&Result, 0, false)

		if e != nil {
			return helper.CreateResult(false, nil, e.Error())
		}

		for _, val := range Result {
			Durationresults = append(Durationresults, val.TimeStamp.UTC())
		}
	}

	// Anomaly Data
	for i := 0; i < 2; i++ {
		var arrsort []string
		if i == 0 {
			arrsort = append(arrsort, "timestamp")
		} else {
			arrsort = append(arrsort, "-timestamp")
		}

		query := DB().Connection.NewQuery().From(new(ScadaData).TableName()).Where(dbox.And(dbox.Eq("isvalidtimeduration", true))).Skip(0).Take(1)
		query = query.Order(arrsort...)

		csr, e := query.Cursor(nil)
		if e != nil {
			return helper.CreateResult(false, nil, e.Error())
		}
		defer csr.Close()

		Result := make([]ScadaData, 0)
		e = csr.Fetch(&Result, 0, false)

		if e != nil {
			return helper.CreateResult(false, nil, e.Error())
		}

		for _, val := range Result {
			ScadaAnomalyresults = append(ScadaAnomalyresults, val.TimeStamp.UTC())
		}
	}

	// AlarmOverlapping Data
	for i := 0; i < 2; i++ {
		var arrsort []string
		if i == 0 {
			arrsort = append(arrsort, "startdate")
		} else {
			arrsort = append(arrsort, "-startdate")
		}

		query := DB().Connection.NewQuery().From(new(AlarmOverlapping).TableName()).Skip(0).Take(1)
		query = query.Order(arrsort...)

		csr, e := query.Cursor(nil)
		if e != nil {
			return helper.CreateResult(false, nil, e.Error())
		}
		defer csr.Close()

		Result := make([]AlarmOverlapping, 0)
		e = csr.Fetch(&Result, 0, false)

		if e != nil {
			return helper.CreateResult(false, nil, e.Error())
		}

		for _, val := range Result {
			AlarmOverlappingresults = append(AlarmOverlappingresults, val.StartDate.UTC())
		}
	}

	// AlarmScadaAnomaly Data
	for i := 0; i < 2; i++ {
		var arrsort []string
		if i == 0 {
			arrsort = append(arrsort, "startdate")
		} else {
			arrsort = append(arrsort, "-startdate")
		}

		query := DB().Connection.NewQuery().From(new(AlarmScadaAnomaly).TableName()).Skip(0).Take(1)
		query = query.Order(arrsort...)

		csr, e := query.Cursor(nil)
		if e != nil {
			return helper.CreateResult(false, nil, e.Error())
		}
		defer csr.Close()

		Result := make([]AlarmScadaAnomaly, 0)
		e = csr.Fetch(&Result, 0, false)

		if e != nil {
			return helper.CreateResult(false, nil, e.Error())
		}

		for _, val := range Result {
			AlarmScadaAnomalyresults = append(AlarmScadaAnomalyresults, val.StartDate.UTC())
		}
	}

	data := struct {
		ScadaData         []time.Time
		Alarm             []time.Time
		JMR               []time.Time
		MET               []time.Time
		Duration          []time.Time
		ScadaAnomaly      []time.Time
		AlarmOverlapping  []time.Time
		AlarmScadaAnomaly []time.Time
	}{
		ScadaData:         Scadaresults,
		Alarm:             Alarmresults,
		JMR:               JMRresults,
		MET:               METresults,
		Duration:          Durationresults,
		ScadaAnomaly:      ScadaAnomalyresults,
		AlarmOverlapping:  AlarmOverlappingresults,
		AlarmScadaAnomaly: AlarmScadaAnomalyresults,
	}

	return helper.CreateResult(true, data, "success")
}

func (m *DataBrowserController) GetScadaAvailDate(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	Scadaresults := make([]time.Time, 0)

	// Scada Data
	for i := 0; i < 2; i++ {
		var arrsort []string
		if i == 0 {
			arrsort = append(arrsort, "timestamp")
		} else {
			arrsort = append(arrsort, "-timestamp")
		}

		query := DB().Connection.NewQuery().From(new(ScadaData).TableName()).Skip(0).Take(1)
		query = query.Order(arrsort...)

		csr, e := query.Cursor(nil)
		if e != nil {
			return helper.CreateResult(false, nil, e.Error())
		}
		defer csr.Close()

		Result := make([]ScadaData, 0)
		e = csr.Fetch(&Result, 0, false)

		if e != nil {
			return helper.CreateResult(false, nil, e.Error())
		}

		for _, val := range Result {
			Scadaresults = append(Scadaresults, val.TimeStamp.UTC())
		}
	}

	data := struct {
		ScadaData []time.Time
	}{
		ScadaData: Scadaresults,
	}

	return helper.CreateResult(true, data, "success")
}

func (m *DataBrowserController) GetAlarmAvailDate(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	Alarmresults := make([]time.Time, 0)

	// Alarm Data
	for i := 0; i < 2; i++ {
		var arrsort []string
		if i == 0 {
			arrsort = append(arrsort, "startdate")
		} else {
			arrsort = append(arrsort, "-startdate")
		}

		query := DB().Connection.NewQuery().From(new(Alarm).TableName()).Skip(0).Take(1)
		query = query.Order(arrsort...)

		csr, e := query.Cursor(nil)
		if e != nil {
			return helper.CreateResult(false, nil, e.Error())
		}
		defer csr.Close()

		Result := make([]Alarm, 0)
		e = csr.Fetch(&Result, 0, false)

		if e != nil {
			return helper.CreateResult(false, nil, e.Error())
		}

		for _, val := range Result {
			Alarmresults = append(Alarmresults, val.StartDate.UTC())
		}
	}

	data := struct {
		Alarm []time.Time
	}{
		Alarm: Alarmresults,
	}

	return helper.CreateResult(true, data, "success")
}

func (m *DataBrowserController) GetJMRAvailDate(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	JMRresults := make([]time.Time, 0)

	// JMR Data
	for i := 0; i < 2; i++ {
		var arrsort []string
		if i == 0 {
			arrsort = append(arrsort, "dateinfo.dateid")
		} else {
			arrsort = append(arrsort, "-dateinfo.dateid")
		}

		query := DB().Connection.NewQuery().From(new(JMR).TableName()).Skip(0).Take(1)
		query = query.Order(arrsort...)

		csr, e := query.Cursor(nil)
		if e != nil {
			return helper.CreateResult(false, nil, e.Error())
		}
		defer csr.Close()

		Result := make([]ScadaData, 0)
		e = csr.Fetch(&Result, 0, false)

		if e != nil {
			return helper.CreateResult(false, nil, e.Error())
		}

		for _, val := range Result {
			JMRresults = append(JMRresults, val.DateInfo.DateId.UTC())
		}
	}

	data := struct {
		JMR []time.Time
	}{
		JMR: JMRresults,
	}

	return helper.CreateResult(true, data, "success")
}

func (m *DataBrowserController) GetMETAvailDate(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	METresults := make([]time.Time, 0)

	// MET Data
	for i := 0; i < 2; i++ {
		var arrsort []string
		if i == 0 {
			arrsort = append(arrsort, "timestamp")
		} else {
			arrsort = append(arrsort, "-timestamp")
		}

		query := DB().Connection.NewQuery().From(new(MetTower).TableName()).Skip(0).Take(1)
		query = query.Order(arrsort...)

		csr, e := query.Cursor(nil)
		if e != nil {
			return helper.CreateResult(false, nil, e.Error())
		}
		defer csr.Close()

		Result := make([]MetTower, 0)
		e = csr.Fetch(&Result, 0, false)

		if e != nil {
			return helper.CreateResult(false, nil, e.Error())
		}

		for _, val := range Result {
			METresults = append(METresults, val.TimeStamp.UTC())
		}
	}

	data := struct {
		MET []time.Time
	}{
		MET: METresults,
	}

	return helper.CreateResult(true, data, "success")
}

func (m *DataBrowserController) GetDurationAvailDate(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	Durationresults := make([]time.Time, 0)

	// Duration Data
	for i := 0; i < 2; i++ {
		var arrsort []string
		if i == 0 {
			arrsort = append(arrsort, "timestamp")
		} else {
			arrsort = append(arrsort, "-timestamp")
		}

		// query := DB().Connection.NewQuery().From(new(ScadaData).TableName()).Where(dbox.And(dbox.Eq("isvalidtimeduration", false))).Skip(0).Take(1)
		query := DB().Connection.NewQuery().From(new(ScadaData).TableName()).Skip(0).Take(1)
		query = query.Order(arrsort...)

		csr, e := query.Cursor(nil)
		if e != nil {
			return helper.CreateResult(false, nil, e.Error())
		}
		defer csr.Close()

		Result := make([]ScadaData, 0)
		e = csr.Fetch(&Result, 0, false)

		if e != nil {
			return helper.CreateResult(false, nil, e.Error())
		}

		for _, val := range Result {
			Durationresults = append(Durationresults, val.TimeStamp.UTC())
		}
	}

	data := struct {
		Duration []time.Time
	}{
		Duration: Durationresults,
	}

	return helper.CreateResult(true, data, "success")
}

func (m *DataBrowserController) GetScadaAnomalyAvailDate(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	ScadaAnomalyresults := make([]time.Time, 0)

	// Anomaly Data
	for i := 0; i < 2; i++ {
		var arrsort []string
		if i == 0 {
			arrsort = append(arrsort, "timestamp")
		} else {
			arrsort = append(arrsort, "-timestamp")
		}

		query := DB().Connection.NewQuery().From(new(ScadaData).TableName()).Where(dbox.And(dbox.Eq("isvalidtimeduration", true))).Skip(0).Take(1)
		query = query.Order(arrsort...)

		csr, e := query.Cursor(nil)
		if e != nil {
			return helper.CreateResult(false, nil, e.Error())
		}
		defer csr.Close()

		Result := make([]ScadaData, 0)
		e = csr.Fetch(&Result, 0, false)

		if e != nil {
			return helper.CreateResult(false, nil, e.Error())
		}

		for _, val := range Result {
			ScadaAnomalyresults = append(ScadaAnomalyresults, val.TimeStamp.UTC())
		}
	}

	data := struct {
		ScadaAnomaly []time.Time
	}{
		ScadaAnomaly: ScadaAnomalyresults,
	}

	return helper.CreateResult(true, data, "success")
}

func (m *DataBrowserController) GetAlarmOverlappingAvailDate(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	AlarmOverlappingresults := make([]time.Time, 0)

	// AlarmOverlapping Data
	for i := 0; i < 2; i++ {
		var arrsort []string
		if i == 0 {
			arrsort = append(arrsort, "startdate")
		} else {
			arrsort = append(arrsort, "-startdate")
		}

		query := DB().Connection.NewQuery().From(new(AlarmOverlapping).TableName()).Skip(0).Take(1)
		query = query.Order(arrsort...)

		csr, e := query.Cursor(nil)
		if e != nil {
			return helper.CreateResult(false, nil, e.Error())
		}
		defer csr.Close()

		Result := make([]AlarmOverlapping, 0)
		e = csr.Fetch(&Result, 0, false)

		if e != nil {
			return helper.CreateResult(false, nil, e.Error())
		}

		for _, val := range Result {
			AlarmOverlappingresults = append(AlarmOverlappingresults, val.StartDate.UTC())
		}
	}

	data := struct {
		AlarmOverlapping []time.Time
	}{
		AlarmOverlapping: AlarmOverlappingresults,
	}

	return helper.CreateResult(true, data, "success")
}

func (m *DataBrowserController) GetAlarmScadaAnomalyAvailDate(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	AlarmScadaAnomalyresults := make([]time.Time, 0)

	// AlarmScadaAnomaly Data
	for i := 0; i < 2; i++ {
		var arrsort []string
		if i == 0 {
			arrsort = append(arrsort, "startdate")
		} else {
			arrsort = append(arrsort, "-startdate")
		}

		query := DB().Connection.NewQuery().From(new(AlarmScadaAnomaly).TableName()).Skip(0).Take(1)
		query = query.Order(arrsort...)

		csr, e := query.Cursor(nil)
		if e != nil {
			return helper.CreateResult(false, nil, e.Error())
		}
		defer csr.Close()

		Result := make([]AlarmScadaAnomaly, 0)
		e = csr.Fetch(&Result, 0, false)

		if e != nil {
			return helper.CreateResult(false, nil, e.Error())
		}

		for _, val := range Result {
			AlarmScadaAnomalyresults = append(AlarmScadaAnomalyresults, val.StartDate.UTC())
		}
	}

	data := struct {
		AlarmScadaAnomaly []time.Time
	}{
		AlarmScadaAnomaly: AlarmScadaAnomalyresults,
	}

	return helper.CreateResult(true, data, "success")
}

func (m *DataBrowserController) GetEventAvailDate(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	EventDateresults := make([]time.Time, 0)

	// AlarmScadaAnomaly Data
	for i := 0; i < 2; i++ {
		var arrsort []string
		if i == 0 {
			arrsort = append(arrsort, "timestamp")
		} else {
			arrsort = append(arrsort, "-timestamp")
		}

		query := DB().Connection.NewQuery().From(new(EventRaw).TableName()).Skip(0).Take(1)
		query = query.Order(arrsort...)

		csr, e := query.Cursor(nil)
		if e != nil {
			return helper.CreateResult(false, nil, e.Error())
		}
		defer csr.Close()

		Result := make([]EventRaw, 0)
		e = csr.Fetch(&Result, 0, false)

		if e != nil {
			return helper.CreateResult(false, nil, e.Error())
		}

		for _, val := range Result {
			EventDateresults = append(EventDateresults, val.TimeStamp.UTC())
		}
	}

	data := struct {
		EventDate []time.Time
	}{
		EventDate: EventDateresults,
	}

	return helper.CreateResult(true, data, "success")
}