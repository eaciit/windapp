package controller

import (
	. "eaciit/wfdemo-git/library/core"
	. "eaciit/wfdemo-git/library/models"
	"eaciit/wfdemo-git/web/helper"
	"time"
	// "time"

	tk "github.com/eaciit/toolkit"
	"strings"
	"github.com/eaciit/dbox"
	"github.com/eaciit/knot/knot.v1"
)

type AnalyticDgrScadaController struct {
	App
}

func CreateAnalyticDgrScadaController() *AnalyticDgrScadaController {
	var controller = new(AnalyticDgrScadaController)
	return controller
}

func (m *AnalyticDgrScadaController) GetData(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson
	type DataItem struct {
		power        float64
		energy       float64
		windspeed    float64
		downtime     float64
		plf          float64
		gridavail    float64
		machineavail float64
		trueavail    float64
	}

	type DataReturn struct {
		dgr      DataItem
		scada    DataItem
		variance DataItem
	}

	var totalTurbine, plfDivider float64
	var data DataReturn

	p := new(PayloadAnalytic)
	e := k.GetPayload(&p)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	tStart, tEnd, e := helper.GetStartEndDate(k, p.Period, p.DateStart, p.DateEnd)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	duration := tk.ToFloat64(tEnd.Sub(tStart).Hours()/24, 0, tk.RoundingAuto) // duration in days
	turbine := p.Turbine
	project := p.Project

	var (
		pipes  []tk.M
		filter []*dbox.Filter
	)

	filter = append(filter, dbox.Ne("_id", ""))
	filter = append(filter, dbox.Gte("dateinfo.dateid", tStart))
	filter = append(filter, dbox.Lte("dateinfo.dateid", tEnd))

	if project != "" {
		filter = append(filter, dbox.Eq("projectname", project))
	}

	var turbineList []TurbineOut
	if project != "" {
		turbineList, _ = helper.GetTurbineList([]interface{}{project})
	} else {
		turbineList, _ = helper.GetTurbineList(nil)
	}

	if len(turbine) > 0 {
		for _, vt := range turbine {
			for _, v := range turbineList {
				if vt == v.Value {
					plfDivider += v.Capacity
					totalTurbine += 1
				}
			}
		}

		filter = append(filter, dbox.In("turbine", turbine...))
	} else {
		if project != "" {
			for _, v := range turbineList {
				if project == v.Project {
					plfDivider += v.Capacity
					totalTurbine += 1
				}
			}
		} else {
			for _, v := range turbineList {
				plfDivider += v.Capacity
				totalTurbine += 1
			}
		}

	}

	pipes = append(pipes, tk.M{"$group": tk.M{
		"_id":              "$projectname",
		"PowerKW":          tk.M{"$sum": "$powerkw"},
		"Production":       tk.M{"$sum": "$production"},
		"WS":               tk.M{"$avg": "$avgwindspeed"},
		"OKTime":           tk.M{"$sum": "$oktime"},
		"MachineDownLoss":  tk.M{"$sum": "$machinedownloss"},
		"GridDownLoss":     tk.M{"$sum": "$griddownloss"},
		"PCDeviation":      tk.M{"$sum": "$pcdeviation"},
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
	if e != nil {
		helper.CreateResult(false, nil, e.Error())
	}

	defer csr.Close()

	var scadaItem DataItem

	sPower := 0.0
	sEnergy := 0.0
	sDowntime := 0.0
	sOktime := 0.0
	sGriddowntime := 0.0
	_ = sGriddowntime
	sMachinedowntime := 0.0
	_ = sMachinedowntime
	sWindspeed := 0.0
	sPlf := 0.0
	sGridavail := 0.0
	sMachineavail := 0.0
	sTrueavail := 0.0
	// totalTimeStamp := 0
	scadaDataAvailable := true

	if len(resultScada) > 0 {
		scada := resultScada[0]
		sPower = scada.GetFloat64("PowerKW") / 1000
		sEnergy = scada.GetFloat64("Production") / 1000
		// sWindspeed = scada.GetFloat64("WS")
		sDowntime = scada.GetFloat64("DownTimeDuration")
		sOktime = scada.GetFloat64("OKTime")
	} else {
		scadaDataAvailable = false
	}

	// get scadadata

	pipes = []tk.M{}
	pipes = append(pipes, tk.M{"$match": tk.M{"available": 1}})
	pipes = append(pipes,
		tk.M{"$group": tk.M{"_id": "$projectname",
			"power":           tk.M{"$sum": "$power"},
			"windspeed":       tk.M{"$avg": "$avgwindspeed"},
			"oktime":          tk.M{"$sum": "$oktime"},
			"griddowntime":    tk.M{"$sum": "$griddowntime"},
			"powerlost":       tk.M{"$sum": "$powerlost"},
			"totaltimestamp":  tk.M{"$sum": 1},
			"machinedowntime": tk.M{"$sum": "$machinedowntime"},
			"available":       tk.M{"$sum": "$available"},
			"minutes":         tk.M{"$sum": "$minutes"},
			"maxdate":         tk.M{"$max": "$dateinfo.dateid"},
			"mindate":         tk.M{"$min": "$dateinfo.dateid"},
			"unknowntime":     tk.M{"$sum": "$unknowntime"},
		}})
	pipes = append(pipes, tk.M{"$sort": tk.M{"_id": 1}})

	csr, e = DB().Connection.NewQuery().
		From(new(ScadaData).TableName()).
		Command("pipe", pipes).
		Where(dbox.And(filter...)).
		Cursor(nil)

	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	list := []tk.M{}
	e = csr.Fetch(&list, 0, false)
	defer csr.Close()

	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	if len(list) > 0 {
		scada := list[0]
		sWindspeed = scada.GetFloat64("windspeed")
		minDate := scada.Get("mindate").(time.Time)
		maxDate := scada.Get("maxdate").(time.Time)
		hourValue := helper.GetHourValue(tStart.UTC(), tEnd.UTC(), minDate.UTC(), maxDate.UTC())

		// log.Printf(">> %v | %v - %v | %v >> %v \n", tStart.UTC(), tEnd.UTC(), minDate.UTC(), maxDate.UTC(), hourValue)

		sPlf = sEnergy / (plfDivider * 1000 * hourValue) * 100 * 1000 //sEnergy / (totalTurbine * hourValue * 2100) * 100 * 1000
		sTrueavail = (sOktime / 3600) / (totalTurbine * hourValue) * 100

		minutes := scada.GetFloat64("minutes") / 60
		// totalTimeStamp = scada.GetInt("totaltimestamp")
		sMachineavail = (minutes - (scada.GetFloat64("machinedowntime"))/3600) / (totalTurbine * hourValue) * 100
		sGridavail = (minutes - (scada.GetFloat64("griddowntime"))/3600) / (totalTurbine * hourValue) * 100
	} else {
		scadaDataAvailable = false
	}

	/*maxCount10min := int(totalTurbine) * 144 * tk.ToInt(tk.Div(tEnd.Sub(tStart).Hours(), 24), tk.RoundingUp)

	if totalTimeStamp < maxCount10min {
		sDowntime += tk.Div(tk.ToFloat64((maxCount10min-totalTimeStamp), 0, tk.RoundingAuto)*600.0, 3600.0)
	}*/

	scadaItem.power = sPower
	scadaItem.energy = sEnergy
	scadaItem.windspeed = sWindspeed
	scadaItem.downtime = sDowntime
	scadaItem.plf = sPlf
	scadaItem.gridavail = sGridavail
	scadaItem.machineavail = sMachineavail
	scadaItem.trueavail = sTrueavail

	// ========================================================= DGR

	var filterd []*dbox.Filter
	filterd = append(filterd, dbox.Gte("dateinfo.dateid", tStart))
	filterd = append(filterd, dbox.Lte("dateinfo.dateid", tEnd))
	if project != "" {
		filterd = append(filterd, dbox.Eq("site", project))
	}

	if len(turbine) != 0 {
		filterd = append(filterd, dbox.In("turbine", turbine...))
	}

	pipes = nil
	pipes = append(pipes, tk.M{"$group": tk.M{"_id": "$site", "genkwhday": tk.M{"$sum": "$genkwhday"}, "lostenergy": tk.M{"$sum": "$lostenergy"}, "gridavailability": tk.M{"$avg": "$gridavailability"}, "downtimehours": tk.M{"$sum": "$downtimehours"}, "machineavailability": tk.M{"$avg": "$machineavailability"}, "plfday": tk.M{"$avg": "$plfday"}, "operationalhours": tk.M{"$sum": "$operationalhours"}}})
	//pipes = append(pipes, tk.M{"$sort": tk.M{"_id": 1}})

	csr, e = DB().Connection.NewQuery().
		From(new(DGRModel).TableName()).
		Command("pipe", pipes).
		Where(filterd...).
		Cursor(nil)
	defer csr.Close()

	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	list = nil
	e = csr.Fetch(&list, 0, false)

	// tk.Printf("%#v \n", list)

	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	var dgrItem DataItem

	sPower = 0.0
	sEnergy = 0.0
	sDowntime = 0.0
	sOktime = 0.0
	sGriddowntime = 0.0
	sMachinedowntime = 0.0
	sWindspeed = -0.0
	sPlf = 0.0
	sGridavail = 0.0
	sMachineavail = 0.0
	sTrueavail = 0.0

	dgrDataAvailable := true

	if len(list) > 0 {
		dgr := list[0]
		sEnergy = tk.ToFloat64(dgr["genkwhday"], 2, tk.RoundingAuto) / 1000
		sMachineavail = tk.ToFloat64(dgr["machineavailability"], 2, tk.RoundingAuto)
		sDowntime = tk.ToFloat64(dgr["downtimehours"], 2, tk.RoundingAuto)
		sPlf = tk.ToFloat64(dgr["plfday"], 2, tk.RoundingAuto)
		sGridavail = tk.ToFloat64(dgr["gridavailability"], 2, tk.RoundingAuto)
		sOktime = tk.ToFloat64(dgr["operationalhours"], 2, tk.RoundingAuto)
		sPower = sEnergy * 6
		// sTrueavail = (sOktime) / (duration * 24 * 24)
		sTrueavail = (sOktime) / (duration * 24 * 24) * 100
	} else {
		dgrDataAvailable = false
	}

	dgrItem.power = sPower
	dgrItem.energy = sEnergy
	dgrItem.downtime = sDowntime
	dgrItem.windspeed = sWindspeed
	dgrItem.plf = sPlf
	dgrItem.gridavail = sGridavail
	dgrItem.machineavail = sMachineavail
	dgrItem.trueavail = sTrueavail

	var varItem DataItem
	varItem.power = dgrItem.power - scadaItem.power
	varItem.energy = dgrItem.energy - scadaItem.energy
	varItem.windspeed = dgrItem.windspeed - scadaItem.windspeed
	varItem.plf = dgrItem.plf - scadaItem.plf
	varItem.gridavail = dgrItem.gridavail - scadaItem.gridavail
	varItem.machineavail = dgrItem.machineavail - scadaItem.machineavail
	varItem.trueavail = dgrItem.trueavail - scadaItem.trueavail
	varItem.downtime = dgrItem.downtime - scadaItem.downtime

	/*varItem.power = math.Abs(dgrItem.power - scadaItem.power)
	varItem.energy = math.Abs(dgrItem.energy - scadaItem.energy)
	varItem.windspeed = math.Abs(dgrItem.windspeed - scadaItem.windspeed)
	varItem.plf = math.Abs(dgrItem.plf - scadaItem.plf)
	varItem.gridavail = math.Abs(dgrItem.gridavail - scadaItem.gridavail)
	varItem.machineavail = math.Abs(dgrItem.machineavail - scadaItem.machineavail)
	varItem.trueavail = math.Abs(dgrItem.trueavail - scadaItem.trueavail)*/

	// ========================================================= HFD
	var scadaHfdItem DataItem
	scadaDataHfdAvailable := false
	// _midate := time.Time{}
	// _madate := time.Time{}

	// pipes = []tk.M{}
	// pipes = append(pipes,
	// 	tk.M{"$group": tk.M{"_id": "$projectname",
	// 		"power":          tk.M{"$sum": "$fast_activepower_kw"},
	// 		"windspeed":      tk.M{"$avg": "$fast_windspeed_ms"},
	// 		"totaltimestamp": tk.M{"$sum": 1},
	// 		"maxdate":        tk.M{"$max": "$dateinfo.dateid"},
	// 		"mindate":        tk.M{"$min": "$dateinfo.dateid"},
	// 	}})
	// pipes = append(pipes, tk.M{"$sort": tk.M{"_id": 1}})
	// //fast_activepower_kw
	// //fast_windspeed_ms
	// //
	csrhfd, e := DB().Connection.NewQuery().
		Select("activepower_kw", "windspeed_ms", "dateinfo.dateid").
		From("Scada10MinHFD").
		Where(dbox.And(filter...)).
		Cursor(nil)

	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	ltkm := tk.M{}
	if csrhfd.Count() > 0 {
		scadaDataHfdAvailable = true
	}

	for {
		_tkm := tk.M{}
		_ex := csrhfd.Fetch(&_tkm, 1, false)
		if _ex != nil {
			break
		}

		_dtime := _tkm.Get("dateinfo", tk.M{}).(tk.M).Get("dateid", time.Time{}).(time.Time)
		_midate := ltkm.Get("midate", time.Time{}).(time.Time)
		_madate := ltkm.Get("madate", time.Time{}).(time.Time)
		if !_dtime.IsZero() {
			if _midate.IsZero() || _midate.After(_dtime) {
				_midate = _dtime
			}

			if _madate.IsZero() || _madate.Before(_dtime) {
				_madate = _dtime
			}
		}
		ltkm.Set("midate", _midate)
		ltkm.Set("madate", _madate)

		_dws := _tkm.GetFloat64("windspeed_ms")
		_dap := _tkm.GetFloat64("activepower_kw")
		_cws := float64(1)
		if _dws == -9999999.0 {
			_dws = 0
			_cws = 0
		}

		if _dap == -9999999.0 {
			_dap = 0
		}

		_dws += ltkm.GetFloat64("windspeed")
		_dap += ltkm.GetFloat64("power")
		_cws += ltkm.GetFloat64("cws")

		ltkm.Set("windspeed", _dws)
		ltkm.Set("power", _dap)
		ltkm.Set("cws", _cws)
	}

	csrhfd.Close()

	if scadaDataHfdAvailable {
		scadaHfdItem.power = ltkm.GetFloat64("power") / 1000
		scadaHfdItem.energy = scadaHfdItem.power / 6
		scadaHfdItem.windspeed = tk.Div(ltkm.GetFloat64("windspeed"), ltkm.GetFloat64("cws"))

		minDate := ltkm.Get("midate", time.Time{}).(time.Time)
		maxDate := ltkm.Get("madate", time.Time{}).(time.Time)
		//(totalTurbine * hourValue * 2100) * 100 * 1000
		// tk.Println(" >>> ltkm >>> ", ltkm)
		hourValue := helper.GetHourValue(tStart.UTC(), tEnd.UTC(), minDate.UTC(), maxDate.UTC())
		// tk.Println(" >>> hour >>> ", hourValue)
		scadaHfdItem.plf = tk.Div(scadaHfdItem.energy, (plfDivider*hourValue*1000)) * 100 * 1000
	}

	var diffDgrHfd DataItem
	diffDgrHfd.power = dgrItem.power - scadaHfdItem.power
	diffDgrHfd.energy = dgrItem.energy - scadaHfdItem.energy
	diffDgrHfd.windspeed = dgrItem.windspeed - scadaHfdItem.windspeed
	diffDgrHfd.plf = dgrItem.plf - scadaHfdItem.plf
	diffDgrHfd.gridavail = dgrItem.gridavail - scadaHfdItem.gridavail
	diffDgrHfd.machineavail = dgrItem.machineavail - scadaHfdItem.machineavail
	diffDgrHfd.trueavail = dgrItem.trueavail - scadaHfdItem.trueavail
	diffDgrHfd.downtime = dgrItem.downtime - scadaHfdItem.downtime

	// _ = scadaHfdItem
	// _ = scadaDataHfdAvailable
	/*=====================*/

	data.scada = scadaItem
	data.dgr = dgrItem
	data.variance = varItem

	result := []tk.M{}
	result = append(result, tk.M{"desc": "Power (MW)", "dgr": dgrItem.power, "scada": scadaItem.power, "difference": varItem.power, "ScadaHFD": scadaHfdItem.power, "diffdgrhfd": diffDgrHfd.power})
	result = append(result, tk.M{"desc": "Energy (MWh)", "dgr": dgrItem.energy, "scada": scadaItem.energy, "difference": varItem.energy, "ScadaHFD": scadaHfdItem.energy, "diffdgrhfd": diffDgrHfd.energy})
	result = append(result, tk.M{"desc": "Avg. Wind Speed (m/s)", "dgr": "N/A", "scada": scadaItem.windspeed, "difference": varItem.windspeed, "ScadaHFD": scadaHfdItem.windspeed, "diffdgrhfd": diffDgrHfd.windspeed})
	result = append(result, tk.M{"desc": "Downtime (Hours)", "dgr": dgrItem.downtime, "scada": scadaItem.downtime, "difference": varItem.downtime, "ScadaHFD": "N/A", "diffdgrhfd": diffDgrHfd.downtime})
	result = append(result, tk.M{"desc": "PLF", "dgr": dgrItem.plf, "scada": scadaItem.plf, "difference": varItem.plf, "ScadaHFD": scadaHfdItem.plf, "diffdgrhfd": diffDgrHfd.plf})
	result = append(result, tk.M{"desc": "Grid Availability", "dgr": dgrItem.gridavail, "scada": scadaItem.gridavail, "difference": varItem.gridavail, "ScadaHFD": "N/A", "diffdgrhfd": diffDgrHfd.gridavail})
	result = append(result, tk.M{"desc": "Machine Availability", "dgr": dgrItem.machineavail, "scada": scadaItem.machineavail, "difference": varItem.machineavail, "ScadaHFD": "N/A", "diffdgrhfd": diffDgrHfd.machineavail})
	result = append(result, tk.M{"desc": "True Availability", "dgr": dgrItem.trueavail, "scada": scadaItem.trueavail, "difference": varItem.trueavail, "ScadaHFD": "N/A", "diffdgrhfd": diffDgrHfd.trueavail})

	if scadaDataAvailable == false {
		for _, val := range result {
			val["scada"] = "N/A"
		}
	}

	if dgrDataAvailable == false {
		for _, val := range result {
			val["dgr"] = "N/A"
		}
	}

	if !scadaDataHfdAvailable {
		for _, val := range result {
			val["ScadaHFD"] = "N/A"
		}
	}

	return helper.CreateResult(true, result, "success")
}

func (m *AnalyticDgrScadaController) GetDataRev(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson
	type DataItem struct {
		power        float64
		energy       float64
		windspeed    float64
		downtime     float64
		plf          float64
		gridavail    float64
		machineavail float64
		trueavail    float64
	}

	type DataReturn struct {
		dgr      DataItem
		scada    DataItem
		variance DataItem
	}

	var totalTurbine, turbinecapacity float64
	availdate := tk.M{}
	// var data DataReturn

	p := new(PayloadAnalytic)
	e := k.GetPayload(&p)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	tStart, tEnd, e := helper.GetStartEndDate(k, p.Period, p.DateStart, p.DateEnd)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	duration := tk.ToFloat64(tEnd.Sub(tStart).Hours()/24, 0, tk.RoundingAuto) // duration in days
	turbine := p.Turbine
	project := p.Project

	var (
		pipes      []tk.M
		filter     []*dbox.Filter
		prodfilter []*dbox.Filter
	)

	filter = append(filter, dbox.Ne("_id", ""))

	if project != "" {
		filter = append(filter, dbox.Eq("projectname", project))
	}

	var turbineList []TurbineOut
	if project != "" {
		turbineList, _ = helper.GetTurbineList([]interface{}{project})
	} else {
		turbineList, _ = helper.GetTurbineList(nil)
	}

	if len(turbine) > 0 {
		for _, vt := range turbine {
			for _, v := range turbineList {
				if vt == v.Value {
					turbinecapacity += v.Capacity
					totalTurbine += 1
				}
			}
		}

		filter = append(filter, dbox.In("turbine", turbine...))
	} else {
		if project != "" {
			for _, v := range turbineList {
				if project == v.Project {
					turbinecapacity += v.Capacity
					totalTurbine += 1
				}
			}
		} else {
			for _, v := range turbineList {
				turbinecapacity += v.Capacity
				totalTurbine += 1
			}
		}

	}

	prodfilter = append(prodfilter, filter...)

	filter = append(filter, dbox.Gte("dateinfo.dateid", tStart))
	filter = append(filter, dbox.Lte("dateinfo.dateid", tEnd))

	prodfilter = append(prodfilter, dbox.Gte("timestamp", tStart))
	prodfilter = append(prodfilter, dbox.Lte("timestamp", tEnd))

	pipes = append(pipes, tk.M{"$group": tk.M{
		"_id":              "$projectname",
		"PowerKW":          tk.M{"$sum": "$powerkw"},
		"Production":       tk.M{"$sum": "$production"},
		"WS":               tk.M{"$avg": "$avgwindspeed"},
		"OKTime":           tk.M{"$sum": "$oktime"},
		"DownTimeDuration": tk.M{"$sum": "$downtimehours"}, //totaldowntimehours
		"machinedownhours": tk.M{"$sum": "$machinedownhours"},
		"griddownhours":    tk.M{"$sum": "$griddownhours"},
		"LossEnergy":       tk.M{"$sum": "$lostenergy"},
		"SumWindSpeed":     tk.M{"$sum": "$detwindspeed.sumwindspeed"},
		"CountWindSpeed":   tk.M{"$sum": "$detwindspeed.countwindspeed"},
		"plf":              tk.M{"$avg": "$plf"},
		"machineavail":     tk.M{"$avg": "$machineavail"},
		"gridavail":        tk.M{"$avg": "$gridavail"},
		"mindate":          tk.M{"$min": "$dateinfo.dateid"},
		"maxdate":          tk.M{"$max": "$dateinfo.dateid"},
	}})

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
	if e != nil {
		helper.CreateResult(false, nil, e.Error())
	}

	defer csr.Close()

	var scadaItem DataItem

	sPower := 0.0
	sEnergy := 0.0
	sDowntime := 0.0
	sOktime := 0.0
	sWindspeed := 0.0
	sPlf := 0.0
	sGridavail := 0.0
	sMachineavail := 0.0
	sTrueavail := 0.0
	scadaDataAvailable := true

	if len(resultScada) > 0 {
		scada := resultScada[0]

		minDate := scada.Get("mindate").(time.Time).UTC()
		maxDate := scada.Get("maxdate").(time.Time).UTC()
		duration = maxDate.AddDate(0, 0, 1).UTC().Sub(minDate.UTC()).Hours()
		totalhours := duration * totalTurbine

		availdate.Set("scada", []time.Time{minDate, maxDate})

		sPower = scada.GetFloat64("PowerKW") / 1000
		sWindspeed = tk.Div(scada.GetFloat64("SumWindSpeed"), scada.GetFloat64("CountWindSpeed"))
		sEnergy = scada.GetFloat64("Production") / 1000
		sDowntime = scada.GetFloat64("machinedownhours") + scada.GetFloat64("griddownhours")
		// sDowntime = scada.GetFloat64("DownTimeDuration")
		sOktime = scada.GetFloat64("OKTime")
		// tk.Println(sOktime, " / (", duration, " * ", totalTurbine, " * ", 3600)
		sTrueavail = (sOktime / (duration * totalTurbine * 3600)) * 100
		sPlf = (sEnergy / (turbinecapacity * duration)) * 100
		sMachineavail = tk.Div(totalhours-scada.GetFloat64("machinedownhours"), totalhours) * 100
		sGridavail = tk.Div(totalhours-scada.GetFloat64("griddownhours"), totalhours) * 100
	} else {
		scadaDataAvailable = false
	}

	scadaItem.power = sPower
	scadaItem.energy = sEnergy
	scadaItem.windspeed = sWindspeed
	scadaItem.downtime = sDowntime
	scadaItem.plf = sPlf
	scadaItem.gridavail = sGridavail
	scadaItem.machineavail = sMachineavail
	scadaItem.trueavail = sTrueavail

	// ========================================================= DGR

	query := []tk.M{}
	query = append(query, tk.M{"dateinfo.dateid": tk.M{"$gte": tStart}})
	query = append(query, tk.M{"dateinfo.dateid": tk.M{"$lte": tEnd}})

	if project != "" {
		query = append(query, tk.M{"chosensite": project})
	}

	reffdgrturb := getturbinedgr(project)
	if len(turbine) != 0 {
		intturbine := []string{}
		for _, _turb := range turbine {
			intturbine = append(intturbine, reffdgrturb.GetString(tk.ToString(_turb)))
		}
		query = append(query, tk.M{"turbine": tk.M{"$in": intturbine}})
	}
	// tk.Println(">>> ", query)
	pipes = nil
	pipes = append(pipes, tk.M{"$match": tk.M{"$and": query}})
	pipes = append(pipes, tk.M{"$group": tk.M{"_id": "$chosensite",
		"genkwhday":           tk.M{"$sum": "$genkwhday"},
		"lostenergy":          tk.M{"$sum": "$lostenergy"},
		"gridavailability":    tk.M{"$avg": "$intga"}, //extga
		"downtimehours":       tk.M{"$sum": "$totaldowntimehours"},
		"machineavailability": tk.M{"$avg": "$ma"},
		"plfday":              tk.M{"$avg": "$plfday"},
		"mindate":             tk.M{"$min": "$dateinfo.dateid"},
		"maxdate":             tk.M{"$max": "$dateinfo.dateid"},
		"generationhours":     tk.M{"$sum": "$generationhours"},
		//for calculate "A" value
		"egscheduled":     tk.M{"$sum": "$egscheduled"},
		"egunscheduled":   tk.M{"$sum": "$egunscheduled"},
		"fmenvironmental": tk.M{"$sum": "$fmenvironmental"},
		"fmtheft":         tk.M{"$sum": "$fmtheft"},
		"rowignonoem":     tk.M{"$sum": "$rowignonoem"},
		"rowmcnonoem":     tk.M{"$sum": "$rowmcnonoem"},
		//for calculate grid avail value
		"igscheduled":    tk.M{"$sum": "$igscheduled"},
		"igunscheduled":  tk.M{"$sum": "$igunscheduled"},
		"pssscheduled":   tk.M{"$sum": "$pssscheduled"},
		"pssunscheduled": tk.M{"$sum": "$pssunscheduled"},
		"rowigoem":       tk.M{"$sum": "$rowigoem"},
		//for calculate machine avail value
		"wtgscheduled":   tk.M{"$sum": "$wtgscheduled"},
		"wtgunscheduled": tk.M{"$sum": "$wtgunscheduled"},
		"rowmcoem":       tk.M{"$sum": "$rowmcoem"},
		"aor":            tk.M{"$sum": "$aor"},
	}})

	csr, e = DB().Connection.NewQuery().
		From(new(DGRModel).TableName()).
		Command("pipe", pipes).
		Cursor(nil)
	defer csr.Close()

	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	list := []tk.M{}
	e = csr.Fetch(&list, 0, false)

	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	var dgrItem DataItem

	sPower = 0.0
	sEnergy = 0.0
	sDowntime = 0.0
	sOktime = 0.0
	sWindspeed = -0.0
	sPlf = 0.0
	sGridavail = 0.0
	sMachineavail = 0.0
	sTrueavail = 0.0

	dgrDataAvailable := true
	if len(list) > 0 {
		dgr := list[0]

		minDate := dgr.Get("mindate").(time.Time).UTC()
		maxDate := dgr.Get("maxdate").(time.Time).UTC()
		duration = maxDate.AddDate(0, 0, 1).UTC().Sub(minDate.UTC()).Hours()

		availdate.Set("dgr", []time.Time{minDate, maxDate})

		// aduration := (duration * float64(len(turbine))) - dgr.GetFloat64("egscheduled") - dgr.GetFloat64("egunscheduled") - dgr.GetFloat64("fmenvironmental") - dgr.GetFloat64("fmtheft") - dgr.GetFloat64("rowignonoem") - dgr.GetFloat64("rowmcnonoem")
		// gdown := dgr.GetFloat64("igscheduled") + dgr.GetFloat64("igunscheduled") + dgr.GetFloat64("pssscheduled") + dgr.GetFloat64("pssunscheduled") + dgr.GetFloat64("rowigoem")
		// mdown := dgr.GetFloat64("wtgscheduled") + dgr.GetFloat64("wtgunscheduled") + dgr.GetFloat64("rowmcoem") + dgr.GetFloat64("aor")

		sEnergy = dgr.GetFloat64("genkwhday") / 1000
		// sDowntime = dgr.GetFloat64("downtimehours")

		// sGridavail = tk.Div((aduration-gdown), aduration) * 100
		// sPlf = (sEnergy / (turbinecapacity * duration)) * 100
		// sMachineavail = tk.Div((aduration-mdown), aduration) * 100

		sOktime = dgr.GetFloat64("generationhours")
		sPower = sEnergy * 6
		// tk.Println(sOktime, " / ", duration, " * ", len(turbine))
		sTrueavail = (sOktime) / (duration * float64(len(turbine))) * 100
	} else {
		dgrDataAvailable = false
	}

	query = append(query, tk.M{"category": tk.M{"$ne": ""}})
	downGroupClause := tk.M{
		"_id" : "$category",
		"sumbreakdownhours" : tk.M{"$sum": "$breakdownhours"},
	}
	pipes = nil
	pipes = append(pipes, tk.M{"$match": tk.M{"$and": query}})
	pipes = append(pipes, tk.M{"$group": downGroupClause})
	tk.Println("match", query)
	csr, e = DB().Connection.NewQuery().
		From("rpt_downtime").
		Command("pipe", pipes).
		Cursor(nil)
	defer csr.Close()

	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	downList := []tk.M{}
	e = csr.Fetch(&downList, 0, false)

	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}
	totalhours := float64( len(turbine) * 24 * daysDiff(tEnd, tStart) )

	tmpA := 0.0
	for _, down := range downList{
		category := strings.ToUpper(down.GetString("_id"))
		if !strings.Contains(category, "LOAD"){
			sDowntime += down.GetFloat64("sumbreakdownhours")
		}

		switch category{
			case "S-M/C", "U-M/C", "ROW(M/C)-OEM", "AOR":
				sMachineavail += down.GetFloat64("sumbreakdownhours")
				break
			case "S-EG", "U-EG":
				tmpA += down.GetFloat64("sumbreakdownhours")
				sGridavail += down.GetFloat64("sumbreakdownhours")
				break
			case "ROW(IG)-NONOEM", "ROW(M/C)-NONOEM", "FM-ENV", "FM-THEFT" :
				tmpA += down.GetFloat64("sumbreakdownhours")
				break
		}
		
	}

	sPlf = tk.Div(sEnergy, totalhours * getturbinemw(project)) * 100
	A := totalhours - tmpA
	dgrItem.power = sPower
	dgrItem.energy = sEnergy
	dgrItem.downtime = sDowntime
	dgrItem.windspeed = sWindspeed
	dgrItem.plf = sPlf
	dgrItem.gridavail = tk.Div(totalhours - sGridavail, totalhours) * 100
	dgrItem.machineavail = tk.Div(A - sMachineavail, A) * 100
	dgrItem.trueavail = sTrueavail

	var varItem DataItem
	varItem.power = dgrItem.power - scadaItem.power
	varItem.energy = dgrItem.energy - scadaItem.energy
	varItem.windspeed = dgrItem.windspeed - scadaItem.windspeed
	varItem.plf = dgrItem.plf - scadaItem.plf
	varItem.gridavail = dgrItem.gridavail - scadaItem.gridavail
	varItem.machineavail = dgrItem.machineavail - scadaItem.machineavail
	varItem.trueavail = dgrItem.trueavail - scadaItem.trueavail
	varItem.downtime = dgrItem.downtime - scadaItem.downtime

	//===== GET Data from latest production day
	pipes = nil
	pipes = append(pipes, tk.M{"$group": tk.M{
		"_id":        "$projectname",
		"production": tk.M{"$sum": "$value"},
		"min":        tk.M{"$min": "$timestamp"},
		"max":        tk.M{"$max": "$timestamp"},
	}})

	csr, e = DBRealtime().NewQuery().
		From("log_latestdataproduction").
		Command("pipe", pipes).
		Where(dbox.And(prodfilter...)).
		Cursor(nil)

	if e != nil {
		helper.CreateResult(false, nil, e.Error())
	}
	defer csr.Close()

	resultprod := []tk.M{}
	e = csr.Fetch(&resultprod, 0, false)
	if e != nil {
		helper.CreateResult(false, nil, e.Error())
	}

	defer csr.Close()

	tk.Println(">>> ", resultprod)
	rprod := float64(0)
	if len(resultprod) > 0 {
		rprod = resultprod[0].GetFloat64("production")
		if resultprod[0].GetString("_id") != "Lahori" {
			rprod = rprod / 1000
		}
	}
	//===== GET Data from latest production day

	result := []tk.M{}
	result = append(result, tk.M{"desc": "Power (MW)", "dgr": dgrItem.power, "scada": scadaItem.power, "difference": varItem.power, "ScadaHFD": rprod * 6, "diffdgrhfd": "N/A"})
	result = append(result, tk.M{"desc": "Energy (MWh)", "dgr": dgrItem.energy, "scada": scadaItem.energy, "difference": varItem.energy, "ScadaHFD": rprod, "diffdgrhfd": "N/A"})
	result = append(result, tk.M{"desc": "Avg. Wind Speed (m/s)", "dgr": "N/A", "scada": scadaItem.windspeed, "difference": varItem.windspeed, "ScadaHFD": "N/A", "diffdgrhfd": "N/A"})
	result = append(result, tk.M{"desc": "Downtime (Hours)", "dgr": dgrItem.downtime, "scada": scadaItem.downtime, "difference": varItem.downtime, "ScadaHFD": "N/A", "diffdgrhfd": "N/A"})
	result = append(result, tk.M{"desc": "PLF", "dgr": dgrItem.plf, "scada": scadaItem.plf, "difference": varItem.plf, "ScadaHFD": "N/A", "diffdgrhfd": "N/A"})
	result = append(result, tk.M{"desc": "Grid Availability", "dgr": dgrItem.gridavail, "scada": scadaItem.gridavail, "difference": varItem.gridavail, "ScadaHFD": "N/A", "diffdgrhfd": "N/A"})
	result = append(result, tk.M{"desc": "Machine Availability", "dgr": dgrItem.machineavail, "scada": scadaItem.machineavail, "difference": varItem.machineavail, "ScadaHFD": "N/A", "diffdgrhfd": "N/A"})
	result = append(result, tk.M{"desc": "True Availability", "dgr": dgrItem.trueavail, "scada": scadaItem.trueavail, "difference": varItem.trueavail, "ScadaHFD": "N/A", "diffdgrhfd": "N/A"})

	if scadaDataAvailable == false {
		for _, val := range result {
			val["scada"] = "N/A"
			availdate.Set("scada", "N/A")
		}
	}

	if dgrDataAvailable == false {
		for _, val := range result {
			val["dgr"] = "N/A"
			availdate.Set("dgr", "N/A")
		}
	}

	lastres := tk.M{}.Set("data", result).Set("availdate", availdate)
	return helper.CreateResult(true, lastres, "success")
}

func getturbinedgr(iproject string) (tkm tk.M) {
	tkm = tk.M{}

	csr, e := DB().Connection.NewQuery().
		Select("turbineid", "turbinedgr",).
		From("ref_turbine").
		Where(dbox.Eq("project", iproject)).
		Cursor(nil)
	defer csr.Close()

	res := []tk.M{}
	e = csr.Fetch(&res, 0, false)
	if e != nil {
		return
	}

	for _, val := range res {
		tkm.Set(val.GetString("turbineid"), val.GetString("turbinedgr"))
	}

	return
}

func getturbinemw(iproject string) (mw float64) {
	mw = 0.0

	csr, e := DB().Connection.NewQuery().
		Select("capacitymw",).
		From("ref_turbine").
		Where(dbox.Eq("project", iproject)).
		Cursor(nil)
	defer csr.Close()

	res := []tk.M{}
	e = csr.Fetch(&res, 0, false)
	if e != nil {
		return
	}

	if len(res) > 0{
		mw  = res[0].GetFloat64("capacitymw")
	}

	return
}

func lastDayOfYear(t time.Time) time.Time {
	return time.Date(t.Year(), 12, 31, 0, 0, 0, 0, t.Location())
}

func firstDayOfNextYear(t time.Time) time.Time {
	return time.Date(t.Year()+1, 1, 1, 0, 0, 0, 0, t.Location())
}

// a - b in days
func daysDiff(a, b time.Time) (days int) {
	cur := b
	for cur.Year() < a.Year() {
		// add 1 to count the last day of the year too.
		days += lastDayOfYear(cur).YearDay() - cur.YearDay() + 1
		cur = firstDayOfNextYear(cur)	
	}
	days += a.YearDay() - cur.YearDay()
	if b.AddDate(0, 0, days).After(a) {
		days -= 1
	}
	days += 1
	if days < 0 {
		return -days
	}

	return days
}
