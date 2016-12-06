package controller

import (
	. "eaciit/wfdemo-git/library/core"
	. "eaciit/wfdemo-git/library/models"
	"eaciit/wfdemo-git/web/helper"
	"math"
	// "time"

	tk "github.com/eaciit/toolkit"

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

	duration := tEnd.Sub(tStart).Hours() / 24 // duration in days
	turbine := p.Turbine
	project := p.Project

	var (
		pipes  []tk.M
		filter []*dbox.Filter
	)

	var data DataReturn

	filter = append(filter, dbox.Gte("dateinfo.dateid", tStart))
	filter = append(filter, dbox.Lte("dateinfo.dateid", tEnd))

	if project != "" {
		filter = append(filter, dbox.Eq("projectname", project))
	}

	if len(turbine) != 0 {
		filter = append(filter, dbox.In("turbine", turbine...))
	}

	pipes = append(pipes, tk.M{"$group": tk.M{"_id": "$projectname", "power": tk.M{"$sum": "$power"}, "windspeed": tk.M{"$avg": "$avgwindspeed"}, "oktime": tk.M{"$sum": "$oktime"}, "griddowntime": tk.M{"$sum": "$griddowntime"}, "powerlost": tk.M{"$sum": "$powerlost"}, "totaltimestamp": tk.M{"$sum": 1}, "machinedowntime": tk.M{"$sum": "$machinedowntime"}, "available": tk.M{"$sum": "$available"}, "minutes": tk.M{"$sum": "$minutes"}}})
	pipes = append(pipes, tk.M{"$sort": tk.M{"_id": 1}})

	csr, e := DB().Connection.NewQuery().
		From(new(ScadaData).TableName()).
		Command("pipe", pipes).
		Where(dbox.And(filter...)).
		Cursor(nil)

	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	list := []tk.M{}
	e = csr.Fetch(&list, 0, false)
	csr.Close()

	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

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

	scadaDataAvailable := true

	if len(list) > 0 {
		scada := list[0]
		sPower = scada.GetFloat64("power") / 1000 // in KWh
		sMinutesInHour := scada.GetFloat64("minutes") / 60.0
		sEnergy = sPower / 6
		sOktime = scada.GetFloat64("oktime")
		sDowntime = ((24 * duration * 144 * 600) - sOktime) / 3600
		sWindspeed = scada.GetFloat64("windspeed")
		sPlf = sEnergy / (24 * duration * 24 * 2100) * 100 * 1000
		sTrueavail = (sOktime / 3600) / (duration * 24 * 24) * 100

		sMachineavail = (sMinutesInHour - (scada.GetFloat64("machinedowntime") / 3600)) / (24 * 24 * duration) * 100
		sGridavail = (sMinutesInHour - (scada.GetFloat64("griddowntime") / 3600)) / (24 * 24 * duration) * 100
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
	varItem.power = math.Abs(scadaItem.power - dgrItem.power)
	varItem.energy = math.Abs(scadaItem.energy - dgrItem.energy)
	varItem.windspeed = math.Abs(scadaItem.windspeed - dgrItem.windspeed)
	varItem.plf = math.Abs(scadaItem.plf - dgrItem.plf)
	varItem.gridavail = math.Abs(scadaItem.gridavail - dgrItem.gridavail)
	varItem.machineavail = math.Abs(scadaItem.machineavail - dgrItem.machineavail)
	varItem.trueavail = math.Abs(scadaItem.trueavail - dgrItem.trueavail)
	varItem.downtime = math.Abs(scadaItem.downtime - dgrItem.downtime)

	/*varItem.power = math.Abs(dgrItem.power - scadaItem.power)
	varItem.energy = math.Abs(dgrItem.energy - scadaItem.energy)
	varItem.windspeed = math.Abs(dgrItem.windspeed - scadaItem.windspeed)
	varItem.plf = math.Abs(dgrItem.plf - scadaItem.plf)
	varItem.gridavail = math.Abs(dgrItem.gridavail - scadaItem.gridavail)
	varItem.machineavail = math.Abs(dgrItem.machineavail - scadaItem.machineavail)
	varItem.trueavail = math.Abs(dgrItem.trueavail - scadaItem.trueavail)*/

	data.scada = scadaItem
	data.dgr = dgrItem
	data.variance = varItem

	result := []tk.M{}
	result = append(result, tk.M{"desc": "Power (MW)", "dgr": dgrItem.power, "scada": scadaItem.power, "difference": varItem.power})
	result = append(result, tk.M{"desc": "Energy (MWh)", "dgr": dgrItem.energy, "scada": scadaItem.energy, "difference": varItem.energy})
	result = append(result, tk.M{"desc": "Avg. Wind Speed (m/s)", "dgr": "N/A", "scada": scadaItem.windspeed, "difference": varItem.windspeed})
	result = append(result, tk.M{"desc": "Downtime (Hours)", "dgr": dgrItem.downtime, "scada": scadaItem.downtime, "difference": varItem.downtime})
	result = append(result, tk.M{"desc": "PLF", "dgr": dgrItem.plf, "scada": scadaItem.plf, "difference": varItem.plf})
	result = append(result, tk.M{"desc": "Grid Availability", "dgr": dgrItem.gridavail, "scada": scadaItem.gridavail, "difference": varItem.gridavail})
	result = append(result, tk.M{"desc": "Machine Availability", "dgr": dgrItem.machineavail, "scada": scadaItem.machineavail, "difference": varItem.machineavail})
	result = append(result, tk.M{"desc": "True Availability", "dgr": dgrItem.trueavail, "scada": scadaItem.trueavail, "difference": varItem.trueavail})

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

	return helper.CreateResult(true, result, "success")
}

/*func (m *AnalyticDgrScadaController) GetData(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	p := new(helper.Payloads)
	e := k.GetPayload(&p)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	filter := p.ParseFilter()

	if e != nil {
		helper.CreateResult(false, nil, e.Error())
	}

	return helper.CreateResult(true, data, "success")
}
*/
