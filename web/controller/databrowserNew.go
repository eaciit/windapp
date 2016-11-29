package controller

import (
	. "eaciit/wfdemo-git-dev/library/core"
	. "eaciit/wfdemo-git-dev/library/models"
	"eaciit/wfdemo-git-dev/web/helper"
	"strings"

	"time"

	"github.com/eaciit/dbox"
	"github.com/eaciit/knot/knot.v1"
	tk "github.com/eaciit/toolkit"
)

type DataBrowserNewController struct {
	App
}

func CreateDataBrowserNewController() *DataBrowserNewController {
	var controller = new(DataBrowserNewController)
	return controller
}

func (m *DataBrowserNewController) GetScadaList(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	p := new(helper.Payloads)
	e := k.GetPayload(&p)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}
	filter, _ := p.ParseFilter()

	query := DB().Connection.NewQuery().From(new(ScadaDataOEM).TableName()).Skip(p.Skip).Take(p.Take)
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

	queryC := DB().Connection.NewQuery().From(new(ScadaDataOEM).TableName()).Where(dbox.And(filter...))
	ccount, e := queryC.Cursor(nil)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}
	defer ccount.Close()

	totalPower := 0.0
	totalPowerLost := 0.0
	totalProduction := 0.0
	avgWindSpeed := 0.0
	totalTurbine := 0

	aggrData := []tk.M{}

	queryAggr := DB().Connection.NewQuery().From(new(ScadaDataOEM).TableName()).
		Aggr(dbox.AggrSum, "$power", "TotalPower").
		Aggr(dbox.AggrSum, "$powerlost", "TotalPowerLost").
		Aggr(dbox.AggrSum, "$ai_intern_activpower", "TotalProduction").
		Aggr(dbox.AggrAvr, "$ai_intern_windspeed", "AvgWindSpeed").
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
		totalProduction += val.GetFloat64("TotalProduction")
		avgWindSpeed += val.GetFloat64("AvgWindSpeed")
	}
	totalTurbine = tk.SliceLen(aggrData)

	data := struct {
		Data            []ScadaDataOEM
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
		TotalProduction: totalProduction,
		AvgWindSpeed:    avgWindSpeed,
		TotalTurbine:    totalTurbine,
	}

	return helper.CreateResult(true, data, "success")
}

func (m *DataBrowserNewController) GetScadaDataOEMAvailDate(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	ScadaDataOEMresults := make([]time.Time, 0)

	// ScadaDataOEM Data
	for i := 0; i < 2; i++ {
		var arrsort []string
		if i == 0 {
			arrsort = append(arrsort, "timestamp")
		} else {
			arrsort = append(arrsort, "-timestamp")
		}

		query := DB().Connection.NewQuery().From(new(ScadaDataOEM).TableName()).Skip(0).Take(1)
		query = query.Order(arrsort...)

		csr, e := query.Cursor(nil)
		if e != nil {
			return helper.CreateResult(false, nil, e.Error())
		}
		defer csr.Close()

		Result := make([]ScadaDataOEM, 0)
		e = csr.Fetch(&Result, 0, false)

		if e != nil {
			return helper.CreateResult(false, nil, e.Error())
		}

		for _, val := range Result {
			ScadaDataOEMresults = append(ScadaDataOEMresults, val.TimeStamp.UTC())
		}
	}

	data := struct {
		ScadaDataOEM []time.Time
	}{
		ScadaDataOEM: ScadaDataOEMresults,
	}

	return helper.CreateResult(true, data, "success")
}

func (m *DataBrowserNewController) GetDowntimeEventList(k *knot.WebContext) interface{} {
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
	filter = append(filter, dbox.Gte("timestart", tStart))
	filter = append(filter, dbox.Lte("timestart", tEnd))
	if len(turbine) != 0 {
		filter = append(filter, dbox.In("turbine", turbine...))
	}

	query := DB().Connection.NewQuery().From(new(DowntimeEvent).TableName()).Skip(p.Skip).Take(p.Take)
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

	tmpResult := make([]DowntimeEvent, 0)
	results := make([]DowntimeEvent, 0)
	e = csr.Fetch(&tmpResult, 0, false)
	// tk.Printf("FILTER : %s \n", filter)

	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	for _, val := range tmpResult {
		val.TimeStart = val.TimeStart.UTC()
		results = append(results, val)
	}

	queryC := DB().Connection.NewQuery().From(new(DowntimeEvent).TableName()).Where(dbox.And(filter...))
	ccount, e := queryC.Cursor(nil)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}
	defer ccount.Close()

	totalDuration := 0.0
	totalTurbine := 0

	aggrData := []tk.M{}

	queryAggr := DB().Connection.NewQuery().From(new(DowntimeEvent).TableName()).
		Aggr(dbox.AggrSum, "$duration", "duration").
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
		totalDuration += val.GetFloat64("duration")
	}
	totalTurbine = tk.SliceLen(aggrData)

	data := struct {
		Data          []DowntimeEvent
		Total         int
		TotalTurbine  int
		TotalDuration float64
	}{
		Data:          results,
		Total:         ccount.Count(),
		TotalTurbine:  totalTurbine,
		TotalDuration: totalDuration,
	}

	return helper.CreateResult(true, data, "success")
}

func (m *DataBrowserNewController) GetDowntimeEventvailDate(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	DowntimeEventresults := make([]time.Time, 0)

	// Downtime Event Data
	for i := 0; i < 2; i++ {
		var arrsort []string
		if i == 0 {
			arrsort = append(arrsort, "timestart")
		} else {
			arrsort = append(arrsort, "-timestart")
		}

		query := DB().Connection.NewQuery().From(new(DowntimeEvent).TableName()).Skip(0).Take(1)
		query = query.Order(arrsort...)

		csr, e := query.Cursor(nil)
		if e != nil {
			return helper.CreateResult(false, nil, e.Error())
		}
		defer csr.Close()

		Result := make([]DowntimeEvent, 0)
		e = csr.Fetch(&Result, 0, false)

		if e != nil {
			return helper.CreateResult(false, nil, e.Error())
		}

		for _, val := range Result {
			DowntimeEventresults = append(DowntimeEventresults, val.TimeStart.UTC())
		}
	}

	data := struct {
		DowntimeEvent []time.Time
	}{
		DowntimeEvent: DowntimeEventresults,
	}

	return helper.CreateResult(true, data, "success")
}

func GetCustomFieldList() []tk.M {
	atkm := []tk.M{}

	_ascadaoem_label := []string{"Ai Dfig Torque Actual", "Ai Dr Tr Vib Value", "Ai Gear Oil Pressure", "Ai Hydr System Pressure", "Ai Intern Active Power",
		"Ai Intern Dfig Active Power Actual", "Ai Intern Nacelle Drill", "Ai Intern Nacelle Drill At North Pos Sensor", "Ai Intern Nacelle Pos", "Ai Intern Pitch Angle1",
		"Ai Intern Pitch Angle2", "Ai Intern Pitch Angle3", "Ai Intern Pitch Speed1", "Ai Intern Reactive Power", "Ai Intern Wind Direction",
		"Ai Intern Wind Speed", "Ai Intern Wind Speed Dif", "Ai Tower Vib Value Axial", "Ai Wind Speed1", "Ai Wind Speed2",
		"Ai Wind Vane1", "Ai Wind Vane2", "C Intern Speed Generator", "C Intern Speed Rotor", "Temp Bottom Control Section",
		"Temp Bottom Control Section Low", "Temp Bottom Power Section", "Temp Cabinet Top Box", "Temp Gearbox Hss De", "Temp Gear Box Hss Nde",
		"Temp Gear Box Ims De", "Temp Gear Box Ims Nde", "Temp Gear Oil Sump", "Temp Generator Bearing De", "Temp Generator Bearing Nde",
		"Temp Main Bearing", "Temp Nacelle", "Temp Outdoor", "Time Stamp", "Turbine",
	}

	_ascadaoem_field := []string{"ai_dfig_torque_actual", "ai_drtrvibvalue", "ai_gearoilpressure", "ai_hydrsystempressure", "ai_intern_activpower",
		"ai_intern_dfig_active_power_actual", "ai_intern_nacelledrill", "ai_intern_nacelledrill_at_northpossensor", "ai_intern_nacellepos", "ai_intern_pitchangle1",
		"ai_intern_pitchangle2", "ai_intern_pitchangle3", "ai_intern_pitchspeed1", "ai_intern_reactivpower", "ai_intern_winddirection",
		"ai_intern_windspeed", "ai_intern_windspeeddif", "ai_towervibvalueaxial", "ai_windspeed1", "ai_windspeed2",
		"ai_windvane1", "ai_windvane2", "c_intern_speedgenerator", "c_intern_speedrotor", "temp_bottomcontrolsection",
		"temp_bottomcontrolsection_low", "temp_bottompowersection", "temp_cabinettopbox", "temp_gearbox_hss_de", "temp_gearbox_hss_nde",
		"temp_gearbox_ims_de", "temp_gearbox_ims_nde", "temp_gearoilsump", "temp_generatorbearing_de", "temp_generatorbearing_nde",
		"temp_mainbearing", "temp_nacelle", "temp_outdoor", "timestamp", "turbine",
	}

	_amettower_label := []string{"Time Stamp", "V Hub WS 90m Avg", "V Hub WS 90m Std Dev", "V Ref WS 88m Avg", "V Ref WS 88m Std Dev",
		"V Tip WS 42m Avg", "V Tip WS 42m Std Dev", "D Hub WD 88m Avg", "D Hub WD 88m Std Dev", "D Ref WD 86m Avg",
		"D Ref WD 86m Std Dev", "T Hub & H Hub Humid 85m Avg", "T Hub & H Hub Humid 85m Std Dev", "T Ref & H Ref Humid 85.5m Avg", "T Ref & H Ref Humid 85.5m Std Dev",
		"T Hub & H Hub Temp 85.5m Avg", "T Hub & H Hub Temp 85.5m Std Dev", "T Ref & H Ref Temp 85.5 Avg", "T Ref & H Ref Temp 85.5 Std Dev", "Baro Air Pressure 85.5m Avg", "Baro Air Pressure 85.5m Std Dev",

	}

	_amettower_field := []string{"vhubws90mavg", "vhubws90mstddev", "vrefws88mavg", "vrefws88mstddev", "vtipws42mavg",
		"vtipws42mstddev", "dhubwd88mavg", "dhubwd88mstddev", "drefwd86mavg", "drefwd86mstddev",
		"thubhhubhumid855mavg", "thubhhubhumid855mstddev", "trefhrefhumid855mavg", "trefhrefhumid855mstddev", "thubhhubtemp855mavg",
		"thubhhubtemp855mstddev", "trefhreftemp855mavg", "trefhreftemp855mstddev", "baroairpress855mavg", "baroairpress855mstddev",
	}

	for i, str := range _ascadaoem_field {
		tkm := tk.M{}.
			Set("_id", str).
			Set("label", _ascadaoem_label[i]).
			Set("source", "ScadaDataOEM")

		atkm = append(atkm, tkm)
	}

	for i, str := range _amettower_field {
		tkm := tk.M{}.
			Set("_id", str).
			Set("label", _amettower_label[i]).
			Set("source", "MetTower")

		atkm = append(atkm, tkm)
	}

	return atkm
}

func (m *DataBrowserNewController) GetCustomList(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	p := new(helper.Payloads)
	e := k.GetPayload(&p)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}
	filter, _ := p.ParseFilter()

	istimestamp := false
	arrscadaoem := []string{"_id", "timestamputc"}
	arrmettower := []string{}
	if p.Custom.Has("ColumnList") {
		for _, _val := range p.Custom["ColumnList"].([]interface{}) {
			_tkm, _ := tk.ToM(_val)
			if _tkm.GetString("source") == "ScadaDataOEM" {
				arrscadaoem = append(arrscadaoem, _tkm.GetString("_id"))
				if _tkm.GetString("_id") == "timestamp" {
					istimestamp = true
				}
			} else if _tkm.GetString("source") == "MetTower" {
				arrmettower = append(arrmettower, _tkm.GetString("_id"))
			}
		}
	}

	// tk.Printfn(">>>>>>>>>>>>>Field : %#v", arrscadaoem)

	query := DB().Connection.NewQuery().
		Select(arrscadaoem...).
		From(new(ScadaDataOEM).TableName()).
		Skip(p.Skip).
		Take(p.Take)
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

	// tmpResult := make([]ScadaDataOEM, 0)
	results := make([]tk.M, 0)
	e = csr.Fetch(&results, 0, false)

	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	arrmettowercond := []interface{}{}

	for i, val := range results {
		if val.Has("timestamputc") {
			itime := val.Get("timestamputc", time.Time{}).(time.Time).UTC()
			arrmettowercond = append(arrmettowercond, itime)
			val.Set("timestamputc", itime)
			results[i] = val
		}
		if istimestamp {
			itime := val.Get("timestamp", time.Time{}).(time.Time)
			val.Set("timestamp", itime.UTC())
			results[i] = val
		}
	}

	tkmmet := tk.M{}
	if len(arrmettower) > 0 && len(arrmettowercond) > 0 {
		arrmettower = append(arrmettower, "timestamp")
		_csr, _e := DB().Connection.NewQuery().
			Select(arrmettower...).
			From("MetTower").
			Where(dbox.In("timestamp", arrmettowercond...)).Cursor(nil)
		if _e != nil {
			return helper.CreateResult(false, nil, _e.Error())
		}
		defer _csr.Close()

		_resmet := make([]tk.M, 0)
		_e = _csr.Fetch(&_resmet, 0, false)

		if _e != nil {
			return helper.CreateResult(false, nil, _e.Error())
		}

		for _, val := range _resmet {
			itime := val.Get("timestamp", time.Time{}).(time.Time).UTC().String()
			tkmmet.Set(itime, val)
		}
	}

	if len(tkmmet) > 0 {
		for i, val := range results {
			itime := val.Get("timestamputc", time.Time{}).(time.Time).UTC().String()
			if tkmmet.Has(itime) {
				for _key, _val := range tkmmet[itime].(tk.M) {
					if _key != "timestamp" {
						val.Set(_key, _val)
					}
				}
			}
			val.Unset("timestamputc")
			results[i] = val
		}
	}

	queryC := DB().Connection.NewQuery().From(new(ScadaDataOEM).TableName()).Where(dbox.And(filter...))
	ccount, e := queryC.Cursor(nil)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}
	defer ccount.Close()

	totalPower := 0.0
	totalPowerLost := 0.0
	totalProduction := 0.0
	avgWindSpeed := 0.0
	totalTurbine := 0

	aggrData := []tk.M{}

	queryAggr := DB().Connection.NewQuery().From(new(ScadaDataOEM).TableName()).
		Aggr(dbox.AggrSum, "$power", "TotalPower").
		Aggr(dbox.AggrSum, "$powerlost", "TotalPowerLost").
		Aggr(dbox.AggrSum, "$ai_intern_activpower", "TotalProduction").
		Aggr(dbox.AggrAvr, "$ai_intern_windspeed", "AvgWindSpeed").
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
		totalProduction += val.GetFloat64("TotalProduction")
		avgWindSpeed += val.GetFloat64("AvgWindSpeed")
	}
	totalTurbine = tk.SliceLen(aggrData)

	data := struct {
		Data            []tk.M
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
		TotalProduction: totalProduction,
		AvgWindSpeed:    avgWindSpeed,
		TotalTurbine:    totalTurbine,
	}

	return helper.CreateResult(true, data, "success")
}

func (m *DataBrowserNewController) GetCustomAvailDate(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	Dateresults := make([]time.Time, 0)

	// ScadaDataOEM
	for i := 0; i < 2; i++ {
		var arrsort []string
		if i == 0 {
			arrsort = append(arrsort, "timestamp")
		} else {
			arrsort = append(arrsort, "-timestamp")
		}

		query := DB().Connection.NewQuery().From(new(ScadaDataOEM).TableName()).Skip(0).Take(1)
		query = query.Order(arrsort...)

		queryMetTower := DB().Connection.NewQuery().From(new(MetTower).TableName()).Skip(0).Take(1)
		queryMetTower = queryMetTower.Order(arrsort...)

		csr, e := query.Cursor(nil)
		csrM, eM := queryMetTower.Cursor(nil)
		if e != nil {
			return helper.CreateResult(false, nil, e.Error())
		}
		defer csr.Close()

		if eM != nil {
			return helper.CreateResult(false, nil, e.Error())
		}
		defer csrM.Close()

		Result := make([]ScadaDataOEM, 0)
		e = csr.Fetch(&Result, 0, false)

		ResultMetTower := make([]MetTower, 0)
		eM = csrM.Fetch(&ResultMetTower, 0, false)

		tk.Printf("Result : %s \n", Result)
		tk.Printf("ResultMetTower : %s \n", ResultMetTower)

		for _, val := range Result {
			Dateresults = append(Dateresults, val.TimeStamp.UTC())
		}
		for _, val := range ResultMetTower {
			Dateresults = append(Dateresults, val.TimeStamp.UTC())
		}
	}


	data := struct {
		CustomDate []time.Time
	}{
		CustomDate: Dateresults,
	}

	return helper.CreateResult(true, data, "success")
}
