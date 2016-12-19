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
	"os"
	x "github.com/tealeg/xlsx"
	"strconv"
	// f "path/filepath"
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

	tStart, _ := time.Parse("2006-01-02", p.DateStart.UTC().Format("2006-01-02"))
	tEnd, _ := time.Parse("2006-01-02 15:04:05", p.DateEnd.UTC().Format("2006-01-02")+" 23:59:59")
	turbine := p.Turbine


	filter = append(filter, dbox.Ne("_id", ""))
	// filter = append(filter, dbox.Ne("powerlost", ""))
	// filter = append(filter, dbox.Ne("ai_intern_activpower", ""))
	// filter = append(filter, dbox.Ne("ai_intern_windspeed", ""))
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

	// tmpResult := make([]EventRaw, 0)
	results := make([]EventRaw, 0)
	e = csr.Fetch(&results, 0, false)

	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	// for _, val := range tmpResult {
	// 	val.TimeStamp = val.TimeStamp.UTC()
	// 	results = append(results, val)
	// }

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

// Generate excel

func (m *DataBrowserController) GenExcelScadaOem(k *knot.WebContext) interface{} {

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

	var pathDownload string
	typeExcel := "ScadaOem"
	TimeCreate := time.Now().Format("2006-01-02_150405")
	CreateDateTime := typeExcel + TimeCreate

	if err := os.RemoveAll("web/assets/Excel/" + typeExcel + "/"); err != nil {
		tk.Println(err)
	}

	filter = append(filter, dbox.Ne("_id", ""))
	// filter = append(filter, dbox.Ne("powerlost", ""))
	// filter = append(filter, dbox.Ne("ai_intern_activpower", ""))
	// filter = append(filter, dbox.Ne("ai_intern_windspeed", ""))
	filter = append(filter, dbox.Gte("timestamp", tStart))
	filter = append(filter, dbox.Lte("timestamp", tEnd))
	if len(turbine) != 0 {
		filter = append(filter, dbox.In("turbine", turbine...))
	}

	query := DB().Connection.NewQuery().From(new(ScadaDataOEM).TableName()).Where(dbox.And(filter...))

	csr, e := query.Cursor(nil)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}
	defer csr.Close()

	tmpResult := make([]ScadaDataOEM, 0)
	results := make([]ScadaDataOEM, 0)
	e = csr.Fetch(&tmpResult, 0, false)

	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	for _, val := range tmpResult {
		val.TimeStamp = val.TimeStamp.UTC()
		results = append(results, val)
	}
	//web/assets/Excel/

	if _, err := os.Stat("web/assets/Excel/" + typeExcel + "/"); os.IsNotExist(err) {
		os.MkdirAll("web/assets/Excel/"+typeExcel+"/", 0777)
	}

	DeserializeScadaOem(results, 0, typeExcel, CreateDateTime)
	pathDownload = "res/Excel/" + typeExcel + "/" + CreateDateTime + ".xlsx"
	// tk.Println(pathDownload)

	return helper.CreateResult(true, pathDownload, "success")
}

func (m *DataBrowserController) GenExcelDowntimeEvent(k *knot.WebContext) interface{} {

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

	var pathDownload string
	typeExcel := "DowntimeEvent"
	TimeCreate := time.Now().Format("2006-01-02_150405")
	CreateDateTime := typeExcel + TimeCreate

	if err := os.RemoveAll("web/assets/Excel/" + typeExcel + "/"); err != nil {
		tk.Println(err)
	}

	filter = append(filter, dbox.Ne("_id", ""))
	filter = append(filter, dbox.Gte("timestart", tStart))
	filter = append(filter, dbox.Lte("timestart", tEnd))
	if len(turbine) != 0 {
		filter = append(filter, dbox.In("turbine", turbine...))
	}

	query := DB().Connection.NewQuery().From(new(EventDown).TableName()).Where(dbox.And(filter...))

	csr, e := query.Cursor(nil)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}
	defer csr.Close()

	tmpResult := make([]EventDown, 0)
	// results := make([]EventDown, 0)
	e = csr.Fetch(&tmpResult, 0, false)

	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	// for _, val := range tmpResult {
	// 	val.TimeStart = val.TimeStart.UTC()
	// 	results = append(results, val)
	// }
	//web/assets/Excel/

	if _, err := os.Stat("web/assets/Excel/" + typeExcel + "/"); os.IsNotExist(err) {
		os.MkdirAll("web/assets/Excel/"+typeExcel+"/", 0777)
	}

	DeserializeEventDown(tmpResult, 0, typeExcel, CreateDateTime)
	pathDownload = "res/Excel/" + typeExcel + "/" + CreateDateTime + ".xlsx"
	// tk.Println(pathDownload)

	return helper.CreateResult(true, pathDownload, "success")
}

func (m *DataBrowserController) GenExcelEventRaw(k *knot.WebContext) interface{} {

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

	var pathDownload string
	typeExcel := "EventRaw"
	TimeCreate := time.Now().Format("2006-01-02_150405")
	CreateDateTime := typeExcel + TimeCreate

	if err := os.RemoveAll("web/assets/Excel/" + typeExcel + "/"); err != nil {
		tk.Println(err)
	}

	filter = append(filter, dbox.Ne("_id", ""))
	filter = append(filter, dbox.Gte("timestamp", tStart))
	filter = append(filter, dbox.Lte("timestamp", tEnd))
	if len(turbine) != 0 {
		filter = append(filter, dbox.In("turbine", turbine...))
	}

	query := DB().Connection.NewQuery().From(new(EventRaw).TableName()).Where(dbox.And(filter...))

	csr, e := query.Cursor(nil)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}
	defer csr.Close()

	tmpResult := make([]EventRaw, 0)
	e = csr.Fetch(&tmpResult, 0, false)

	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	if _, err := os.Stat("web/assets/Excel/" + typeExcel + "/"); os.IsNotExist(err) {
		os.MkdirAll("web/assets/Excel/"+typeExcel+"/", 0777)
	}

	DeserializeEventRaw(tmpResult, 0, typeExcel, CreateDateTime)
	pathDownload = "res/Excel/" + typeExcel + "/" + CreateDateTime + ".xlsx"

	return helper.CreateResult(true, pathDownload, "success")
}

func (m *DataBrowserController) GenExcelMet(k *knot.WebContext) interface{} {

	k.Config.OutputType = knot.OutputJson

	var filter []*dbox.Filter

	p := new(helper.PayloadsDB)
	e := k.GetPayload(&p)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	tStart, _ := time.Parse("2006-01-02", p.DateStart.UTC().Format("2006-01-02"))
	tEnd, _ := time.Parse("2006-01-02 15:04:05", p.DateEnd.UTC().Format("2006-01-02")+" 23:59:59")
	// turbine := p.Turbine

	var pathDownload string
	typeExcel := "MetTower"
	TimeCreate := time.Now().Format("2006-01-02_150405")
	CreateDateTime := typeExcel + TimeCreate

	if err := os.RemoveAll("web/assets/Excel/" + typeExcel + "/"); err != nil {
		tk.Println(err)
	}

	filter = append(filter, dbox.Ne("_id", ""))
	filter = append(filter, dbox.Gte("timestamp", tStart))
	filter = append(filter, dbox.Lte("timestamp", tEnd))

	query := DB().Connection.NewQuery().From(new(MetTower).TableName()).Where(dbox.And(filter...))

	csr, e := query.Cursor(nil)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}
	defer csr.Close()

	tmpResult := make([]MetTower, 0)
	e = csr.Fetch(&tmpResult, 0, false)

	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	if _, err := os.Stat("web/assets/Excel/" + typeExcel + "/"); os.IsNotExist(err) {
		os.MkdirAll("web/assets/Excel/"+typeExcel+"/", 0777)
	}

	DeserializeMetTower(tmpResult, 0, typeExcel, CreateDateTime)
	pathDownload = "res/Excel/" + typeExcel + "/" + CreateDateTime + ".xlsx"

	return helper.CreateResult(true, pathDownload, "success")
}

// // Deserialize

func DeserializeScadaOem(data []ScadaDataOEM, j int, typeExcel string, CreateDateTime string) error {
	//savecipo += 1
	filename := ""
	filename = "web/assets/Excel/" + typeExcel + "/" + CreateDateTime + ".xlsx"

	file := x.NewFile()
	sheet, _ := file.AddSheet("Sheet1")
	header := []string{"TimeStamp","Turbine","AI intern R PidAngleOut","AI intern ActivPower ","AI intern I1 ","AI intern I2","AI intern I3","AI intern NacelleDrill ","AI intern NacellePos ","AI intern PitchAkku V1 ","AI intern PitchAkku V2 ","AI intern PitchAkku V3 ","AI intern PitchAngle1 ","AI intern PitchAngle2 ","AI intern PitchAngle3 ","AI intern PitchConv Current1 ","AI intern PitchConv Current2 ","AI intern PitchConv Current3 ","AI intern PitchAngleSP Diff1 ","AI intern PitchAngleSP Diff2 ","AI intern PitchAngleSP Diff3 ","AI intern ReactivPower ","AI intern RpmDiff ","AI intern U1 ","AI intern U2 ","AI intern U3 ","AI intern WindDirection ","AI intern WindSpeed ","AI Intern WindSpeedDif ","AI speed RotFR ","AI WindSpeed1 ","AI WindSpeed2 ","AI WindVane1 ","AI WindVane2 ","AI internCurrentAsym ","Temp GearBox IMS NDE ","AI intern WindVaneDiff ","C intern SpeedGenerator ","C intern SpeedRotor ","AI intern Speed RPMDiff FR1 RotCNT ","AI intern Frequency Grid ","Temp GearBox HSS NDE ","AI DrTrVibValue ","AI intern InLastErrorConv1 ","AI intern InLastErrorConv2 ","AI intern InLastErrorConv3 ","AI intern TempConv1 ","AI intern TempConv2 ","AI intern TempConv3 ","AI intern PitchSpeed2","Temp YawBrake 1 ","Temp YawBrake 2 ","Temp G1L1 ","Temp G1L2 ","Temp G1L3 ","Temp YawBrake 4","AI HydrSystemPressure ","Temp BottomControlSection Low ","Temp GearBox HSS DE ","Temp GearOilSump ","Temp GeneratorBearing DE ","Temp GeneratorBearing NDE ","Temp MainBearing ","Temp GearBox IMS DE ","Temp Nacelle ","Temp Outdoor ","AI TowerVibValueAxial ","AI intern DiffGenSpeedSPToAct ","Temp YawBrake 5","AI intern SpeedGenerator Proximity ","AI intern SpeedDiff Encoder Proximity ","AI GearOilPressure ","Temp CabinetTopBox Low ","Temp CabinetTopBox ","Temp BottomControlSection ","Temp BottomPowerSection ","Temp BottomPowerSection Low ","AI intern Pitch1 Status High ","AI intern Pitch2 Status High ","AI intern Pitch3 Status High ","AI intern InPosition1 ch3","AI intern InPosition2 ch3","AI intern InPosition3 ch3","AI intern Temp Brake Blade1 ","AI intern Temp Brake Blade2 ","AI intern Temp Brake Blade3 ","AI intern Temp PitchMotor Blade1 ","AI intern Temp PitchMotor Blade2 ","AI intern Temp PitchMotor Blade3 ","AI intern Temp Hub Additional1 ","AI intern Temp Hub Additional2 ","AI intern Temp Hub Additional3 ","AI intern Pitch1 Status Low ","AI intern Pitch2 Status Low ","AI intern Pitch3 Status Low ","AI intern Battery VoltageBlade1 center ","AI intern Battery VoltageBlade2 center ","AI intern Battery VoltageBlade3 center ","AI intern Battery ChargingCur Blade1 ","AI intern Battery ChargingCur Blade2 ","AI intern Battery ChargingCur Blade3 ","AI intern Battery DischargingCur Blade1 ","AI intern Battery DischargingCur Blade2 ","AI intern Battery DischargingCur Blade3 ","AI intern PitchMotor BrakeVoltage Blade1 ","AI intern PitchMotor BrakeVoltage Blade2 ","AI intern PitchMotor BrakeVoltage Blade3 ","AI intern PitchMotor BrakeCurrent Blade1 ","AI intern PitchMotor BrakeCurrent Blade2 ","AI intern PitchMotor BrakeCurrent Blade3 ","AI intern Temp HubBox Blade1 ","AI intern Temp HubBox Blade2 ","AI intern Temp HubBox Blade3 ","AI intern Temp Pitch1 HeatSink ","AI intern Temp Pitch2 HeatSink ","AI intern Temp Pitch3 HeatSink ","AI intern ErrorStackBlade1 ","AI intern ErrorStackBlade2 ","AI intern ErrorStackBlade3 ","AI intern Temp BatteryBox Blade1 ","AI intern Temp BatteryBox Blade2 ","AI intern Temp BatteryBox Blade3 ","AI intern DC LinkVoltage1 ","AI intern DC LinkVoltage2 ","AI intern DC LinkVoltage3 ","Temp Yaw Motor1 ","Temp Yaw Motor2 ","Temp Yaw Motor3 ","Temp Yaw Motor4 ","AO DFIG Power Setpiont ","AO DFIG Q Setpoint ","AI DFIG Torque actual ","AI DFIG SpeedGenerator Encoder ","AI intern DFIG DC Link Voltage actual ","AI intern DFIG MSC current ","AI intern DFIG Main voltage ","AI intern DFIG Main current ","AI intern DFIG active power actual ","AI intern DFIG reactive power actual ","AI intern DFIG active power actual LSC ","AI intern DFIG LSC current ","AI intern DFIG Data log number ","AI intern Damper OscMagnitude ","AI intern Damper PassbandFullLoad ","AI YawBrake TempRise1 ","AI YawBrake TempRise2 ","AI YawBrake TempRise3 ","AI YawBrake TempRise4 ","AI intern NacelleDrill at NorthPosSensor ",
	}

	for i, each := range data {
		if i == 0 {
			rowHeader := sheet.AddRow()
			for _, hdr := range header {

				cell := rowHeader.AddCell()
				cell.Value = hdr
			}
		}

		rowContent := sheet.AddRow()

		cell := rowContent.AddCell()
		cell.Value = each.TimeStamp.Format("2006-01-02 15:04:05")

		cell = rowContent.AddCell()
		cell.Value = each.Turbine

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.AI_intern_R_PidAngleOut , 'f', -1, 64) 
		 
		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.AI_intern_ActivPower  , 'f', -1, 64) 
		 
		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.AI_intern_I1  , 'f', -1, 64) 
		 
		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.AI_intern_I2 , 'f', -1, 64) 
		 
		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.AI_intern_I3 , 'f', -1, 64) 
		 
		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.AI_intern_NacelleDrill  , 'f', -1, 64) 
		 
		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.AI_intern_NacellePos  , 'f', -1, 64) 
		 
		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.AI_intern_PitchAkku_V1  , 'f', -1, 64) 
		 
		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.AI_intern_PitchAkku_V2  , 'f', -1, 64) 
		 
		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.AI_intern_PitchAkku_V3  , 'f', -1, 64) 
		 
		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.AI_intern_PitchAngle1  , 'f', -1, 64) 
		 
		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.AI_intern_PitchAngle2  , 'f', -1, 64) 
		 
		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.AI_intern_PitchAngle3  , 'f', -1, 64) 
		 
		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.AI_intern_PitchConv_Current1  , 'f', -1, 64) 
		 
		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.AI_intern_PitchConv_Current2  , 'f', -1, 64) 
		 
		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.AI_intern_PitchConv_Current3  , 'f', -1, 64) 
		 
		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.AI_intern_PitchAngleSP_Diff1  , 'f', -1, 64) 
		 
		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.AI_intern_PitchAngleSP_Diff2  , 'f', -1, 64) 
		 
		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.AI_intern_PitchAngleSP_Diff3  , 'f', -1, 64) 
		 
		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.AI_intern_ReactivPower  , 'f', -1, 64) 
		 
		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.AI_intern_RpmDiff  , 'f', -1, 64) 
		 
		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.AI_intern_U1  , 'f', -1, 64) 
		 
		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.AI_intern_U2  , 'f', -1, 64) 
		 
		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.AI_intern_U3  , 'f', -1, 64) 
		 
		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.AI_intern_WindDirection  , 'f', -1, 64) 
		 
		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.AI_intern_WindSpeed  , 'f', -1, 64) 
		 
		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.AI_Intern_WindSpeedDif  , 'f', -1, 64) 
		 
		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.AI_speed_RotFR  , 'f', -1, 64) 
		 
		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.AI_WindSpeed1  , 'f', -1, 64) 
		 
		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.AI_WindSpeed2  , 'f', -1, 64) 
		 
		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.AI_WindVane1  , 'f', -1, 64) 
		 
		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.AI_WindVane2  , 'f', -1, 64) 
		 
		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.AI_internCurrentAsym  , 'f', -1, 64) 
		 
		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Temp_GearBox_IMS_NDE  , 'f', -1, 64) 
		 
		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.AI_intern_WindVaneDiff  , 'f', -1, 64) 
		 
		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.C_intern_SpeedGenerator  , 'f', -1, 64) 
		 
		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.C_intern_SpeedRotor  , 'f', -1, 64) 
		 
		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.AI_intern_Speed_RPMDiff_FR1_RotCNT  , 'f', -1, 64) 
		 
		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.AI_intern_Frequency_Grid  , 'f', -1, 64) 
		 
		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Temp_GearBox_HSS_NDE  , 'f', -1, 64) 
		 
		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.AI_DrTrVibValue  , 'f', -1, 64) 
		 
		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.AI_intern_InLastErrorConv1  , 'f', -1, 64) 
		 
		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.AI_intern_InLastErrorConv2  , 'f', -1, 64) 
		 
		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.AI_intern_InLastErrorConv3  , 'f', -1, 64) 
		 
		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.AI_intern_TempConv1  , 'f', -1, 64) 
		 
		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.AI_intern_TempConv2  , 'f', -1, 64) 
		 
		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.AI_intern_TempConv3  , 'f', -1, 64) 
		 
		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.AI_intern_PitchSpeed1  , 'f', -1, 64) 
		 
		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Temp_YawBrake_1  , 'f', -1, 64) 
		 
		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Temp_YawBrake_2  , 'f', -1, 64) 
		 
		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Temp_G1L1  , 'f', -1, 64) 
		 
		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Temp_G1L2  , 'f', -1, 64) 
		 
		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Temp_G1L3  , 'f', -1, 64) 
		 
		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Temp_YawBrake_3  , 'f', -1, 64) 
		 
		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.AI_HydrSystemPressure  , 'f', -1, 64) 
		 
		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Temp_BottomControlSection_Low  , 'f', -1, 64) 
		 
		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Temp_GearBox_HSS_DE  , 'f', -1, 64) 
		 
		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Temp_GearOilSump  , 'f', -1, 64) 
		 
		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Temp_GeneratorBearing_DE  , 'f', -1, 64) 
		 
		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Temp_GeneratorBearing_NDE  , 'f', -1, 64) 
		 
		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Temp_MainBearing  , 'f', -1, 64) 
		 
		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Temp_GearBox_IMS_DE  , 'f', -1, 64) 
		 
		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Temp_Nacelle  , 'f', -1, 64) 
		 
		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Temp_Outdoor  , 'f', -1, 64) 
		 
		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.AI_TowerVibValueAxial  , 'f', -1, 64) 
		 
		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.AI_intern_DiffGenSpeedSPToAct  , 'f', -1, 64) 
		 
		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Temp_YawBrake_4  , 'f', -1, 64) 
		 
		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.AI_intern_SpeedGenerator_Proximity  , 'f', -1, 64) 
		 
		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.AI_intern_SpeedDiff_Encoder_Proximity  , 'f', -1, 64) 
		 
		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.AI_GearOilPressure  , 'f', -1, 64) 
		 
		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Temp_CabinetTopBox_Low  , 'f', -1, 64) 
		 
		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Temp_CabinetTopBox  , 'f', -1, 64) 
		 
		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Temp_BottomControlSection  , 'f', -1, 64) 
		 
		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Temp_BottomPowerSection  , 'f', -1, 64) 
		 
		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Temp_BottomPowerSection_Low  , 'f', -1, 64) 
		 
		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.AI_intern_Pitch1_Status_High  , 'f', -1, 64) 
		 
		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.AI_intern_Pitch2_Status_High  , 'f', -1, 64) 
		 
		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.AI_intern_Pitch3_Status_High  , 'f', -1, 64) 
		 
		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.AI_intern_InPosition1_ch2  , 'f', -1, 64) 
		 
		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.AI_intern_InPosition2_ch2  , 'f', -1, 64) 
		 
		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.AI_intern_InPosition3_ch2  , 'f', -1, 64) 
		 
		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.AI_intern_Temp_Brake_Blade1  , 'f', -1, 64) 
		 
		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.AI_intern_Temp_Brake_Blade2  , 'f', -1, 64) 
		 
		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.AI_intern_Temp_Brake_Blade3  , 'f', -1, 64) 
		 
		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.AI_intern_Temp_PitchMotor_Blade1  , 'f', -1, 64) 
		 
		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.AI_intern_Temp_PitchMotor_Blade2  , 'f', -1, 64) 
		 
		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.AI_intern_Temp_PitchMotor_Blade3  , 'f', -1, 64) 
		 
		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.AI_intern_Temp_Hub_Additional1  , 'f', -1, 64) 
		 
		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.AI_intern_Temp_Hub_Additional2  , 'f', -1, 64) 
		 
		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.AI_intern_Temp_Hub_Additional3  , 'f', -1, 64) 
		 
		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.AI_intern_Pitch1_Status_Low  , 'f', -1, 64) 
		 
		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.AI_intern_Pitch2_Status_Low  , 'f', -1, 64) 
		 
		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.AI_intern_Pitch3_Status_Low  , 'f', -1, 64) 
		 
		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.AI_intern_Battery_VoltageBlade1_center  , 'f', -1, 64) 
		 
		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.AI_intern_Battery_VoltageBlade2_center  , 'f', -1, 64) 
		 
		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.AI_intern_Battery_VoltageBlade3_center  , 'f', -1, 64) 
		 
		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.AI_intern_Battery_ChargingCur_Blade1  , 'f', -1, 64) 
		 
		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.AI_intern_Battery_ChargingCur_Blade2  , 'f', -1, 64) 
		 
		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.AI_intern_Battery_ChargingCur_Blade3  , 'f', -1, 64) 
		 
		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.AI_intern_Battery_DischargingCur_Blade1  , 'f', -1, 64) 
		 
		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.AI_intern_Battery_DischargingCur_Blade2  , 'f', -1, 64) 
		 
		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.AI_intern_Battery_DischargingCur_Blade3  , 'f', -1, 64) 
		 
		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.AI_intern_PitchMotor_BrakeVoltage_Blade1  , 'f', -1, 64) 
		 
		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.AI_intern_PitchMotor_BrakeVoltage_Blade2  , 'f', -1, 64) 
		 
		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.AI_intern_PitchMotor_BrakeVoltage_Blade3  , 'f', -1, 64) 
		 
		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.AI_intern_PitchMotor_BrakeCurrent_Blade1  , 'f', -1, 64) 
		 
		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.AI_intern_PitchMotor_BrakeCurrent_Blade2  , 'f', -1, 64) 
		 
		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.AI_intern_PitchMotor_BrakeCurrent_Blade3  , 'f', -1, 64) 
		 
		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.AI_intern_Temp_HubBox_Blade1  , 'f', -1, 64) 
		 
		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.AI_intern_Temp_HubBox_Blade2  , 'f', -1, 64) 
		 
		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.AI_intern_Temp_HubBox_Blade3  , 'f', -1, 64) 
		 
		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.AI_intern_Temp_Pitch1_HeatSink  , 'f', -1, 64) 
		 
		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.AI_intern_Temp_Pitch2_HeatSink  , 'f', -1, 64) 
		 
		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.AI_intern_Temp_Pitch3_HeatSink  , 'f', -1, 64) 
		 
		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.AI_intern_ErrorStackBlade1  , 'f', -1, 64) 
		 
		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.AI_intern_ErrorStackBlade2  , 'f', -1, 64) 
		 
		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.AI_intern_ErrorStackBlade3  , 'f', -1, 64) 
		 
		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.AI_intern_Temp_BatteryBox_Blade1  , 'f', -1, 64) 
		 
		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.AI_intern_Temp_BatteryBox_Blade2  , 'f', -1, 64) 
		 
		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.AI_intern_Temp_BatteryBox_Blade3  , 'f', -1, 64) 
		 
		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.AI_intern_DC_LinkVoltage1  , 'f', -1, 64) 
		 
		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.AI_intern_DC_LinkVoltage2  , 'f', -1, 64) 
		 
		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.AI_intern_DC_LinkVoltage3  , 'f', -1, 64) 
		 
		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Temp_Yaw_Motor1  , 'f', -1, 64) 
		 
		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Temp_Yaw_Motor2  , 'f', -1, 64) 
		 
		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Temp_Yaw_Motor3  , 'f', -1, 64) 
		 
		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Temp_Yaw_Motor4  , 'f', -1, 64) 
		 
		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.AO_DFIG_Power_Setpiont  , 'f', -1, 64) 
		 
		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.AO_DFIG_Q_Setpoint  , 'f', -1, 64) 
		 
		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.AI_DFIG_Torque_actual  , 'f', -1, 64) 
		 
		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.AI_DFIG_SpeedGenerator_Encoder  , 'f', -1, 64) 
		 
		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.AI_intern_DFIG_DC_Link_Voltage_actual  , 'f', -1, 64) 
		 
		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.AI_intern_DFIG_MSC_current  , 'f', -1, 64) 
		 
		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.AI_intern_DFIG_Main_voltage  , 'f', -1, 64) 
		 
		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.AI_intern_DFIG_Main_current  , 'f', -1, 64) 
		 
		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.AI_intern_DFIG_active_power_actual  , 'f', -1, 64) 
		 
		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.AI_intern_DFIG_reactive_power_actual  , 'f', -1, 64) 
		 
		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.AI_intern_DFIG_active_power_actual_LSC  , 'f', -1, 64) 
		 
		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.AI_intern_DFIG_LSC_current  , 'f', -1, 64) 
		 
		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.AI_intern_DFIG_Data_log_number  , 'f', -1, 64) 
		 
		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.AI_intern_Damper_OscMagnitude  , 'f', -1, 64) 
		 
		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.AI_intern_Damper_PassbandFullLoad  , 'f', -1, 64) 
		 
		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.AI_YawBrake_TempRise1  , 'f', -1, 64) 
		 
		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.AI_YawBrake_TempRise2  , 'f', -1, 64) 
		 
		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.AI_YawBrake_TempRise3  , 'f', -1, 64) 
		 
		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.AI_YawBrake_TempRise4  , 'f', -1, 64) 
		 
		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.AI_intern_NacelleDrill_at_NorthPosSensor  , 'f', -1, 64) 
		 
	}

	tk.Println(filename)

	err := file.Save(filename)
	if err != nil {
		return err
	}

	return nil
}

func DeserializeEventDown(data []EventDown, j int, typeExcel string, CreateDateTime string) error {
	//savecipo += 1
	filename := ""
	filename = "web/assets/Excel/" + typeExcel + "/" + CreateDateTime + ".xlsx"

	file := x.NewFile()
	sheet, _ := file.AddSheet("Sheet1")
	header := []string{"Turbine","TimeStart","TimeEnd","Down Grid","Down Environment","Down Machine","Alarm Description","Duration (Second)"}

	for i, each := range data {
		if i == 0 {
			rowHeader := sheet.AddRow()
			for _, hdr := range header {

				cell := rowHeader.AddCell()
				cell.Value = hdr
			}
		}

		rowContent := sheet.AddRow()

		cell := rowContent.AddCell()
		cell.Value = each.Turbine

		cell = rowContent.AddCell()
		cell.Value = each.TimeStart.Format("2006-01-02 15:04:05")

		cell = rowContent.AddCell()
		cell.Value = each.TimeEnd.Format("2006-01-02 15:04:05")

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatBool(each.DownGrid)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatBool(each.DownEnvironment)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatBool(each.DownMachine)

		cell = rowContent.AddCell()
		cell.Value = each.AlarmDescription//strconv.FormatFloat(each.AI_intern_R_PidAngleOut , 'f', -1, 64) 
		 
		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Duration  , 'f', -1, 64) 
		 
	}

	err := file.Save(filename)
	if err != nil {
		return err
	}

	return nil
}

func DeserializeEventRaw(data []EventRaw, j int, typeExcel string, CreateDateTime string) error {
	//savecipo += 1
	filename := ""
	filename = "web/assets/Excel/" + typeExcel + "/" + CreateDateTime + ".xlsx"

	file := x.NewFile()
	sheet, _ := file.AddSheet("Sheet1")
	header := []string{"TimeStamp","Project Name","Turbine", "Event Type","Alarm Description", "Turbine Status", "Brake Type", "Brake Program", "Alarm Id", "Alarm Toggle"}

	for i, each := range data {
		if i == 0 {
			rowHeader := sheet.AddRow()
			for _, hdr := range header {

				cell := rowHeader.AddCell()
				cell.Value = hdr
			}
		}

		rowContent := sheet.AddRow()

		cell := rowContent.AddCell()
		cell.Value = each.TimeStamp.Format("2006-01-02 15:04:05")

		cell = rowContent.AddCell()
		cell.Value = each.ProjectName

		cell = rowContent.AddCell()
		cell.Value = each.Turbine

		cell = rowContent.AddCell()
		cell.Value = each.EventType

		cell = rowContent.AddCell()
		cell.Value = each.AlarmDescription

		cell = rowContent.AddCell()
		cell.Value = each.TurbineStatus

		cell = rowContent.AddCell()
		cell.Value = each.BrakeType

		cell = rowContent.AddCell()
		cell.Value = strconv.Itoa(each.BrakeProgram)//strconv.Formatint(each.BrakeProgram , 'f', -1, 64) 

		cell = rowContent.AddCell()
		cell.Value = strconv.Itoa(each.AlarmId)//strconv.Formatint(each.AlarmId , 'f', -1, 64) 

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatBool(each.AlarmToggle)

	}

	err := file.Save(filename)
	if err != nil {
		return err
	}

	return nil
}

func DeserializeMetTower(data []MetTower, j int, typeExcel string, CreateDateTime string) error {
	//savecipo += 1
	filename := ""
	filename = "web/assets/Excel/" + typeExcel + "/" + CreateDateTime + ".xlsx"

	file := x.NewFile()
	sheet, _ := file.AddSheet("Sheet1")
	header := []string{"TimeStamp", "WindDirNo", "WindDirDesc", "WSCategoryNo", "WSCategoryDesc", "VHubWS90mAvg", "VHubWS90mMax", "VHubWS90mMin", "VHubWS90mStdDev", "VHubWS90mCount", "VRefWS88mAvg", "VRefWS88mMax", "VRefWS88mMin", "VRefWS88mStdDev", "VRefWS88mCount", "VTipWS42mAvg", "VTipWS42mMax", "VTipWS42mMin", "VTipWS42mStdDev", "VTipWS42mCount", "DHubWD88mAvg", "DHubWD88mMax", "DHubWD88mMin", "DHubWD88mStdDev", "DHubWD88mCount", "DRefWD86mAvg", "DRefWD86mMax", "DRefWD86mMin", "DRefWD86mStdDev", "DRefWD86mCount", "THubHHubHumid855mAvg", "THubHHubHumid855mMax", "THubHHubHumid855mMin", "THubHHubHumid855mStdDev", "THubHHubHumid855mCount", "TRefHRefHumid855mAvg", "TRefHRefHumid855mMax", "TRefHRefHumid855mMin", "TRefHRefHumid855mStdDev", "TRefHRefHumid855mCount", "THubHHubTemp855mAvg", "THubHHubTemp855mMax", "THubHHubTemp855mMin", "THubHHubTemp855mStdDev", "THubHHubTemp855mCount", "TRefHRefTemp855mAvg", "TRefHRefTemp855mMax", "TRefHRefTemp855mMin", "TRefHRefTemp855mStdDev", "TRefHRefTemp855mCount", "BaroAirPress855mAvg", "BaroAirPress855mMax", "BaroAirPress855mMin", "BaroAirPress855mStdDev", "BaroAirPress855mCount", "YawAngleVoltageAvg", "YawAngleVoltageMax", "YawAngleVoltageMin", "YawAngleVoltageStdDev", "YawAngleVoltageCount", "OtherSensorVoltageAI1Avg", "OtherSensorVoltageAI1Max", "OtherSensorVoltageAI1Min", "OtherSensorVoltageAI1StdDev", "OtherSensorVoltageAI1Count", "OtherSensorVoltageAI2Avg", "OtherSensorVoltageAI2Max", "OtherSensorVoltageAI2Min", "OtherSensorVoltageAI2StdDev", "OtherSensorVoltageAI2Count", "OtherSensorVoltageAI3Avg", "OtherSensorVoltageAI3Max", "OtherSensorVoltageAI3Min", "OtherSensorVoltageAI3StdDev", "OtherSensorVoltageAI3Count", "OtherSensorVoltageAI4Avg", "OtherSensorVoltageAI4Max", "OtherSensorVoltageAI4Min", "OtherSensorVoltageAI4StdDev", "OtherSensorVoltageAI4Count", "GenRPMCurrentAvg", "GenRPMCurrentMax", "GenRPMCurrentMin", "GenRPMCurrentStdDev", "GenRPMCurrentCount", "WS_SCSCurrentAvg", "WS_SCSCurrentMax", "WS_SCSCurrentMin", "WS_SCSCurrentStdDev", "WS_SCSCurrentCount", "RainStatusCount", "RainStatusSum", "OtherSensor2StatusIO1Avg", "OtherSensor2StatusIO1Max", "OtherSensor2StatusIO1Min", "OtherSensor2StatusIO1StdDev", "OtherSensor2StatusIO1Count", "OtherSensor2StatusIO2Avg", "OtherSensor2StatusIO2Max", "OtherSensor2StatusIO2Min", "OtherSensor2StatusIO2StdDev", "OtherSensor2StatusIO2Count", "OtherSensor2StatusIO3Avg", "OtherSensor2StatusIO3Max", "OtherSensor2StatusIO3Min", "OtherSensor2StatusIO3StdDev", "OtherSensor2StatusIO3Count", "OtherSensor2StatusIO4Avg", "OtherSensor2StatusIO4Max", "OtherSensor2StatusIO4Min", "OtherSensor2StatusIO4StdDev", "OtherSensor2StatusIO4Count", "OtherSensor2StatusIO5Avg", "OtherSensor2StatusIO5Max", "OtherSensor2StatusIO5Min", "OtherSensor2StatusIO5StdDev", "OtherSensor2StatusIO5Count", "A1Avg", "A1Max", "A1Min", "A1StdDev", "A1Count", "A2Avg", "A2Max", "A2Min", "A2StdDev", "A2Count", "A3Avg", "A3Max", "A3Min", "A3StdDev", "A3Count", "A4Avg", "A4Max", "A4Min", "A4StdDev", "A4Count", "A5Avg", "A5Max", "A5Min", "A5StdDev", "A5Count", "A6Avg", "A6Max", "A6Min", "A6StdDev", "A6Count", "A7Avg", "A7Max", "A7Min", "A7StdDev", "A7Count", "A8Avg", "A8Max", "A8Min", "A8StdDev", "A8Count", "A9Avg", "A9Max", "A9Min", "A9StdDev", "A9Count", "A10Avg", "A10Max", "A10Min", "A10StdDev", "A10Count", "AC1Avg", "AC1Max", "AC1Min", "AC1StdDev", "AC1Count", "AC2Avg", "AC2Max", "AC2Min", "AC2StdDev", "AC2Count", "C1Avg", "C1Max", "C1Min", "C1StdDev", "C1Count", "C2Avg", "C2Max", "C2Min", "C2StdDev", "C2Count", "C3Avg", "C3Max", "C3Min", "C3StdDev", "C3Count", "D1Avg", "D1Max", "D1Min", "D1StdDev", "M1_1Avg", "M1_1Max", "M1_1Min", "M1_1StdDev", "M1_1Count", "M1_2Avg", "M1_2Max", "M1_2Min", "M1_2StdDev", "M1_2Count", "M1_3Avg", "M1_3Max", "M1_3Min", "M1_3StdDev", "M1_3Count", "M1_4Avg", "M1_4Max", "M1_4Min", "M1_4StdDev", "M1_4Count", "M1_5Avg", "M1_5Max", "M1_5Min", "M1_5StdDev", "M1_5Count", "M2_1Avg", "M2_1Max", "M2_1Min", "M2_1StdDev", "M2_1Count", "M2_2Avg", "M2_2Max", "M2_2Min", "M2_2StdDev", "M2_2Count", "M2_3Avg", "M2_3Max", "M2_3Min", "M2_3StdDev", "M2_3Count", "M2_4Avg", "M2_4Max", "M2_4Min", "M2_4StdDev", "M2_4Count", "M2_5Avg", "M2_5Max", "M2_5Min", "M2_5StdDev", "M2_5Count", "M2_6Avg", "M2_6Max", "M2_6Min", "M2_6StdDev", "M2_6Count", "M2_7Avg", "M2_7Max", "M2_7Min", "M2_7StdDev", "M2_7Count", "M2_8Avg", "M2_8Max", "M2_8Min", "M2_8StdDev", "M2_8Count", "VAvg", "VMax", "VMin", "IAvg", "IMax", "IMin", "T", "Addr"}

	for i, each := range data {
		if i == 0 {
			rowHeader := sheet.AddRow()
			for _, hdr := range header {

				cell := rowHeader.AddCell()
				cell.Value = hdr
			}
		}

		rowContent := sheet.AddRow()

		cell := rowContent.AddCell()
		cell.Value = each.TimeStamp.Format("2006-01-02 15:04:05")

		cell = rowContent.AddCell()
		cell.Value = each.WindDirDesc

		cell = rowContent.AddCell()
		cell.Value = each.WSCategoryDesc

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.VHubWS90mAvg , 'f', -1, 64) 

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.VHubWS90mMax , 'f', -1, 64) 

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.VHubWS90mMin , 'f', -1, 64) 

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.VHubWS90mStdDev , 'f', -1, 64) 

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.VHubWS90mCount , 'f', -1, 64) 

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.VRefWS88mAvg , 'f', -1, 64) 

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.VRefWS88mMax , 'f', -1, 64) 

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.VRefWS88mMin , 'f', -1, 64) 

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.VRefWS88mStdDev , 'f', -1, 64) 

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.VRefWS88mCount , 'f', -1, 64) 

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.VTipWS42mAvg , 'f', -1, 64) 

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.VTipWS42mMax , 'f', -1, 64) 

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.VTipWS42mMin , 'f', -1, 64) 

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.VTipWS42mStdDev , 'f', -1, 64) 

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.VTipWS42mCount , 'f', -1, 64) 

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.DHubWD88mAvg , 'f', -1, 64) 

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.DHubWD88mMax , 'f', -1, 64) 

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.DHubWD88mMin , 'f', -1, 64) 

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.DHubWD88mStdDev , 'f', -1, 64) 

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.DHubWD88mCount , 'f', -1, 64) 

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.DRefWD86mAvg , 'f', -1, 64) 

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.DRefWD86mMax , 'f', -1, 64) 

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.DRefWD86mMin , 'f', -1, 64) 

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.DRefWD86mStdDev , 'f', -1, 64) 

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.DRefWD86mCount , 'f', -1, 64) 

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.THubHHubHumid855mAvg , 'f', -1, 64) 

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.THubHHubHumid855mMax , 'f', -1, 64) 

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.THubHHubHumid855mMin , 'f', -1, 64) 

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.THubHHubHumid855mStdDev , 'f', -1, 64) 

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.THubHHubHumid855mCount , 'f', -1, 64) 

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.TRefHRefHumid855mAvg , 'f', -1, 64) 

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.TRefHRefHumid855mMax , 'f', -1, 64) 

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.TRefHRefHumid855mMin , 'f', -1, 64) 

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.TRefHRefHumid855mStdDev , 'f', -1, 64) 

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.TRefHRefHumid855mCount , 'f', -1, 64) 

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.THubHHubTemp855mAvg , 'f', -1, 64) 

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.THubHHubTemp855mMax , 'f', -1, 64) 

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.THubHHubTemp855mMin , 'f', -1, 64) 

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.THubHHubTemp855mStdDev , 'f', -1, 64) 

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.THubHHubTemp855mCount , 'f', -1, 64) 

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.TRefHRefTemp855mAvg , 'f', -1, 64) 

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.TRefHRefTemp855mMax , 'f', -1, 64) 

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.TRefHRefTemp855mMin , 'f', -1, 64) 

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.TRefHRefTemp855mStdDev , 'f', -1, 64) 

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.TRefHRefTemp855mCount , 'f', -1, 64) 

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.BaroAirPress855mAvg , 'f', -1, 64) 

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.BaroAirPress855mMax , 'f', -1, 64) 

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.BaroAirPress855mMin , 'f', -1, 64) 

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.BaroAirPress855mStdDev , 'f', -1, 64) 

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.BaroAirPress855mCount , 'f', -1, 64) 

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.YawAngleVoltageAvg , 'f', -1, 64) 

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.YawAngleVoltageMax , 'f', -1, 64) 

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.YawAngleVoltageMin , 'f', -1, 64) 

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.YawAngleVoltageStdDev , 'f', -1, 64) 

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.YawAngleVoltageCount , 'f', -1, 64) 

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.OtherSensorVoltageAI1Avg , 'f', -1, 64) 

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.OtherSensorVoltageAI1Max , 'f', -1, 64) 

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.OtherSensorVoltageAI1Min , 'f', -1, 64) 

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.OtherSensorVoltageAI1StdDev , 'f', -1, 64) 

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.OtherSensorVoltageAI1Count , 'f', -1, 64) 

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.OtherSensorVoltageAI2Avg , 'f', -1, 64) 

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.OtherSensorVoltageAI2Max , 'f', -1, 64) 

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.OtherSensorVoltageAI2Min , 'f', -1, 64) 

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.OtherSensorVoltageAI2StdDev , 'f', -1, 64) 

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.OtherSensorVoltageAI2Count , 'f', -1, 64) 

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.OtherSensorVoltageAI3Avg , 'f', -1, 64) 

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.OtherSensorVoltageAI3Max , 'f', -1, 64) 

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.OtherSensorVoltageAI3Min , 'f', -1, 64) 

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.OtherSensorVoltageAI3StdDev , 'f', -1, 64) 

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.OtherSensorVoltageAI3Count , 'f', -1, 64) 

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.OtherSensorVoltageAI4Avg , 'f', -1, 64) 

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.OtherSensorVoltageAI4Max , 'f', -1, 64) 

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.OtherSensorVoltageAI4Min , 'f', -1, 64) 

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.OtherSensorVoltageAI4StdDev , 'f', -1, 64) 

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.OtherSensorVoltageAI4Count , 'f', -1, 64) 

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.GenRPMCurrentAvg , 'f', -1, 64) 

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.GenRPMCurrentMax , 'f', -1, 64) 

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.GenRPMCurrentMin , 'f', -1, 64) 

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.GenRPMCurrentStdDev , 'f', -1, 64) 

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.GenRPMCurrentCount , 'f', -1, 64) 

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.WS_SCSCurrentAvg , 'f', -1, 64) 

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.WS_SCSCurrentMax , 'f', -1, 64) 

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.WS_SCSCurrentMin , 'f', -1, 64) 

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.WS_SCSCurrentStdDev , 'f', -1, 64) 

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.WS_SCSCurrentCount , 'f', -1, 64) 

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.RainStatusCount , 'f', -1, 64) 

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.RainStatusSum , 'f', -1, 64) 

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.OtherSensor2StatusIO1Avg , 'f', -1, 64) 

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.OtherSensor2StatusIO1Max , 'f', -1, 64) 

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.OtherSensor2StatusIO1Min , 'f', -1, 64) 

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.OtherSensor2StatusIO1StdDev , 'f', -1, 64) 

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.OtherSensor2StatusIO1Count , 'f', -1, 64) 

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.OtherSensor2StatusIO2Avg , 'f', -1, 64) 

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.OtherSensor2StatusIO2Max , 'f', -1, 64) 

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.OtherSensor2StatusIO2Min , 'f', -1, 64) 

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.OtherSensor2StatusIO2StdDev , 'f', -1, 64) 

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.OtherSensor2StatusIO2Count , 'f', -1, 64) 

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.OtherSensor2StatusIO3Avg , 'f', -1, 64) 

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.OtherSensor2StatusIO3Max , 'f', -1, 64) 

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.OtherSensor2StatusIO3Min , 'f', -1, 64) 

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.OtherSensor2StatusIO3StdDev , 'f', -1, 64) 

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.OtherSensor2StatusIO3Count , 'f', -1, 64) 

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.OtherSensor2StatusIO4Avg , 'f', -1, 64) 

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.OtherSensor2StatusIO4Max , 'f', -1, 64) 

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.OtherSensor2StatusIO4Min , 'f', -1, 64) 

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.OtherSensor2StatusIO4StdDev , 'f', -1, 64) 

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.OtherSensor2StatusIO4Count , 'f', -1, 64) 

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.OtherSensor2StatusIO5Avg , 'f', -1, 64) 

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.OtherSensor2StatusIO5Max , 'f', -1, 64) 

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.OtherSensor2StatusIO5Min , 'f', -1, 64) 

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.OtherSensor2StatusIO5StdDev , 'f', -1, 64) 

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.OtherSensor2StatusIO5Count , 'f', -1, 64) 

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.A1Avg , 'f', -1, 64) 

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.A1Max , 'f', -1, 64) 

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.A1Min , 'f', -1, 64) 

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.A1StdDev , 'f', -1, 64) 

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.A1Count , 'f', -1, 64) 

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.A2Avg , 'f', -1, 64) 

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.A2Max , 'f', -1, 64) 

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.A2Min , 'f', -1, 64) 

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.A2StdDev , 'f', -1, 64) 

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.A2Count , 'f', -1, 64) 

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.A3Avg , 'f', -1, 64) 

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.A3Max , 'f', -1, 64) 

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.A3Min , 'f', -1, 64) 

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.A3StdDev , 'f', -1, 64) 

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.A3Count , 'f', -1, 64) 

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.A4Avg , 'f', -1, 64) 

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.A4Max , 'f', -1, 64) 

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.A4Min , 'f', -1, 64) 

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.A4StdDev , 'f', -1, 64) 

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.A4Count , 'f', -1, 64) 

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.A5Avg , 'f', -1, 64) 

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.A5Max , 'f', -1, 64) 

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.A5Min , 'f', -1, 64) 

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.A5StdDev , 'f', -1, 64) 

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.A5Count , 'f', -1, 64) 

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.A6Avg , 'f', -1, 64) 

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.A6Max , 'f', -1, 64) 

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.A6Min , 'f', -1, 64) 

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.A6StdDev , 'f', -1, 64) 

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.A6Count , 'f', -1, 64) 

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.A7Avg , 'f', -1, 64) 

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.A7Max , 'f', -1, 64) 

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.A7Min , 'f', -1, 64) 

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.A7StdDev , 'f', -1, 64) 

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.A7Count , 'f', -1, 64) 

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.A8Avg , 'f', -1, 64) 

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.A8Max , 'f', -1, 64) 

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.A8Min , 'f', -1, 64) 

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.A8StdDev , 'f', -1, 64) 

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.A8Count , 'f', -1, 64) 

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.A9Avg , 'f', -1, 64) 

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.A9Max , 'f', -1, 64) 

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.A9Min , 'f', -1, 64) 

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.A9StdDev , 'f', -1, 64) 

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.A9Count , 'f', -1, 64) 

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.A10Avg , 'f', -1, 64) 

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.A10Max , 'f', -1, 64) 

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.A10Min , 'f', -1, 64) 

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.A10StdDev , 'f', -1, 64) 

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.A10Count , 'f', -1, 64) 

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.AC1Avg , 'f', -1, 64) 

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.AC1Max , 'f', -1, 64) 

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.AC1Min , 'f', -1, 64) 

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.AC1StdDev , 'f', -1, 64) 

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.AC1Count , 'f', -1, 64) 

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.AC2Avg , 'f', -1, 64) 

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.AC2Max , 'f', -1, 64) 

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.AC2Min , 'f', -1, 64) 

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.AC2StdDev , 'f', -1, 64) 

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.AC2Count , 'f', -1, 64) 

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.C1Avg , 'f', -1, 64) 

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.C1Max , 'f', -1, 64) 

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.C1Min , 'f', -1, 64) 

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.C1StdDev , 'f', -1, 64) 

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.C1Count , 'f', -1, 64) 

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.C2Avg , 'f', -1, 64) 

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.C2Max , 'f', -1, 64) 

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.C2Min , 'f', -1, 64) 

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.C2StdDev , 'f', -1, 64) 

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.C2Count , 'f', -1, 64) 

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.C3Avg , 'f', -1, 64) 

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.C3Max , 'f', -1, 64) 

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.C3Min , 'f', -1, 64) 

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.C3StdDev , 'f', -1, 64) 

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.C3Count , 'f', -1, 64) 

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.D1Avg , 'f', -1, 64) 

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.D1Max , 'f', -1, 64) 

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.D1Min , 'f', -1, 64) 

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.D1StdDev , 'f', -1, 64) 

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.M1_1Avg , 'f', -1, 64) 

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.M1_1Max , 'f', -1, 64) 

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.M1_1Min , 'f', -1, 64) 

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.M1_1StdDev , 'f', -1, 64) 

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.M1_1Count , 'f', -1, 64) 

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.M1_2Avg , 'f', -1, 64) 

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.M1_2Max , 'f', -1, 64) 

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.M1_2Min , 'f', -1, 64) 

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.M1_2StdDev , 'f', -1, 64) 

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.M1_2Count , 'f', -1, 64) 

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.M1_3Avg , 'f', -1, 64) 

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.M1_3Max , 'f', -1, 64) 

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.M1_3Min , 'f', -1, 64) 

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.M1_3StdDev , 'f', -1, 64) 

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.M1_3Count , 'f', -1, 64) 

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.M1_4Avg , 'f', -1, 64) 

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.M1_4Max , 'f', -1, 64) 

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.M1_4Min , 'f', -1, 64) 

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.M1_4StdDev , 'f', -1, 64) 

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.M1_4Count , 'f', -1, 64) 

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.M1_5Avg , 'f', -1, 64) 

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.M1_5Max , 'f', -1, 64) 

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.M1_5Min , 'f', -1, 64) 

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.M1_5StdDev , 'f', -1, 64) 

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.M1_5Count , 'f', -1, 64) 

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.M2_1Avg , 'f', -1, 64) 

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.M2_1Max , 'f', -1, 64) 

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.M2_1Min , 'f', -1, 64) 

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.M2_1StdDev , 'f', -1, 64) 

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.M2_1Count , 'f', -1, 64) 

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.M2_2Avg , 'f', -1, 64) 

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.M2_2Max , 'f', -1, 64) 

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.M2_2Min , 'f', -1, 64) 

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.M2_2StdDev , 'f', -1, 64) 

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.M2_2Count , 'f', -1, 64) 

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.M2_3Avg , 'f', -1, 64) 

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.M2_3Max , 'f', -1, 64) 

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.M2_3Min , 'f', -1, 64) 

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.M2_3StdDev , 'f', -1, 64) 

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.M2_3Count , 'f', -1, 64) 

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.M2_4Avg , 'f', -1, 64) 

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.M2_4Max , 'f', -1, 64) 

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.M2_4Min , 'f', -1, 64) 

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.M2_4StdDev , 'f', -1, 64) 

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.M2_4Count , 'f', -1, 64) 

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.M2_5Avg , 'f', -1, 64) 

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.M2_5Max , 'f', -1, 64) 

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.M2_5Min , 'f', -1, 64) 

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.M2_5StdDev , 'f', -1, 64) 

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.M2_5Count , 'f', -1, 64) 

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.M2_6Avg , 'f', -1, 64) 

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.M2_6Max , 'f', -1, 64) 

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.M2_6Min , 'f', -1, 64) 

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.M2_6StdDev , 'f', -1, 64) 

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.M2_6Count , 'f', -1, 64) 

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.M2_7Avg , 'f', -1, 64) 

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.M2_7Max , 'f', -1, 64) 

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.M2_7Min , 'f', -1, 64) 

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.M2_7StdDev , 'f', -1, 64) 

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.M2_7Count , 'f', -1, 64) 

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.M2_8Avg , 'f', -1, 64) 

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.M2_8Max , 'f', -1, 64) 

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.M2_8Min , 'f', -1, 64) 

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.M2_8StdDev , 'f', -1, 64) 

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.M2_8Count , 'f', -1, 64) 

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.VAvg , 'f', -1, 64) 

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.VMax , 'f', -1, 64) 

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.VMin , 'f', -1, 64) 

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.IAvg , 'f', -1, 64) 

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.IMax , 'f', -1, 64) 

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.IMin , 'f', -1, 64) 

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.T , 'f', -1, 64) 

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Addr , 'f', -1, 64) 


	}

	err := file.Save(filename)
	if err != nil {
		return err
	}

	return nil
}

