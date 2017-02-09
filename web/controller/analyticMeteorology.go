package controller

import (
	. "eaciit/wfdemo-git/library/core"
	. "eaciit/wfdemo-git/library/helper"
	. "eaciit/wfdemo-git/library/models"
	"eaciit/wfdemo-git/web/helper"
	"sort"
	"strings"
	"time"

	"github.com/eaciit/dbox"
	"github.com/eaciit/knot/knot.v1"
	tk "github.com/eaciit/toolkit"
)

type AnalyticMeteorologyController struct {
	App
}

func CreateAnalyticMeteorologyController() *AnalyticMeteorologyController {
	var controller = new(AnalyticMeteorologyController)
	return controller
}

func (m *AnalyticMeteorologyController) GetWindCorrelation(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	var dataSeries []tk.M
	p := new(PayloadAnalytic)
	e := k.GetPayload(&p)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}
	tStart, tEnd, e := helper.GetStartEndDate(k, p.Period, p.DateStart, p.DateEnd)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	query := []tk.M{}
	pipes := []tk.M{}
	pipesmet := []tk.M{}
	query = append(query, tk.M{"_id": tk.M{"$ne": ""}})
	query = append(query, tk.M{"timestamp": tk.M{"$gte": tStart}})
	query = append(query, tk.M{"timestamp": tk.M{"$lte": tEnd}})
	// query = append(query, tk.M{"avgwindspeed": tk.M{"$gte": 0.5}})

	// tk.Println("Date select : ", tStart, " ~ ", tEnd)

	pipesmet = append(pipesmet, tk.M{"$match": tk.M{"$and": query}})
	pipesmet = append(pipesmet, tk.M{"$project": tk.M{"vhubws90mavg": 1, "timestamp": 1}})

	if p.Project != "" {
		anProject := strings.Split(p.Project, "(")
		query = append(query, tk.M{"projectname": strings.TrimRight(anProject[0], " ")})
	}

	pipes = append(pipes, tk.M{"$match": tk.M{"$and": query}})
	pipes = append(pipes, tk.M{"$project": tk.M{"turbine": 1, "avgwindspeed": 1, "timestamp": 1}})

	csr, err := DB().Connection.NewQuery().From(new(ScadaData).TableName()).
		Command("pipe", pipes).Cursor(nil)

	if err != nil {
		return helper.CreateResult(false, nil, err.Error())
	}
	defer csr.Close()

	allres := tk.M{}
	arrturbine := []string{}
	_tturbine := tk.M{}

	for {
		trx := new(ScadaData)
		e := csr.Fetch(trx, 1, false)
		if e != nil {
			break
		}

		dkey := trx.TimeStamp.Format("20060102150405")

		_tkm := allres.Get(trx.Turbine, tk.M{}).(tk.M)
		if trx.AvgWindSpeed != -99999.0 {
			_tkm.Set(dkey, trx.AvgWindSpeed)
		}

		allres.Set(trx.Turbine, _tkm)
		_tturbine.Set(trx.Turbine, 1)
	}

	csrx, err := DB().Connection.NewQuery().From(new(MetTower).TableName()).
		Command("pipe", pipesmet).Cursor(nil)

	if err != nil {
		return helper.CreateResult(false, nil, err.Error())
	}
	defer csrx.Close()

	for {
		trx := new(MetTower)
		e := csrx.Fetch(trx, 1, false)
		if e != nil {
			break
		}

		dkey := trx.TimeStamp.Format("20060102150405")

		_tkm := allres.Get("MetTower", tk.M{}).(tk.M)
		_tkm.Set(dkey, trx.VHubWS90mAvg)

		allres.Set("MetTower", _tkm)
	}

	for key, _ := range _tturbine {
		arrturbine = append(arrturbine, key)
	}

	sort.Strings(arrturbine)
	pturbine := append([]string{"MetTower"}, arrturbine...)
	arrturbine = append([]string{"Turbine", "MetTower"}, arrturbine...)

	if len(p.Turbine) > 0 {
		pturbine = []string{}
		for _, _v := range p.Turbine {
			pturbine = append(pturbine, tk.ToString(_v))
		}
	}

	for _, _turbine := range pturbine {
		_tkm := tk.M{}.Set("Turbine", _turbine)
		for i := 1; i < len(arrturbine); i++ {
			_dt01 := allres.Get(_turbine, tk.M{}).(tk.M)
			_dt02 := allres.Get(arrturbine[i], tk.M{}).(tk.M)
			if arrturbine[i] != _turbine && len(_dt01) > 0 && len(_dt02) > 0 {
				_tkm.Set(arrturbine[i],
					GetCorrelation(_dt01, _dt02))
			} else {
				_tkm.Set(arrturbine[i], "-")
			}
		}
		dataSeries = append(dataSeries, _tkm)
	}

	data := struct {
		Column []string
		Data   []tk.M
	}{
		Column: arrturbine,
		Data:   dataSeries,
	}

	return helper.CreateResult(true, data, "success")
}

func (c *AnalyticMeteorologyController) AverageWindSpeed(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	type PayloadAvgWindSpeed struct {
		Period          string
		Project         string
		Turbine         []interface{}
		DateStart       time.Time
		DateEnd         time.Time
		SeriesBreakDown string
		TimeBreakDown   string
	}

	var (
		pipes []tk.M
		// metTower []tk.M
		turbines []tk.M
		list     []tk.M
	)

	p := new(PayloadAvgWindSpeed)
	e := k.GetPayload(&p)

	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	now := time.Now()
	last := time.Now().AddDate(0, -12, 0)

	tStart, _ := time.Parse("20060102", last.Format("200601")+"01")
	tEnd, _ := time.Parse("20060102", now.Format("200601")+"01")

	// tStart, tEnd, e := helper.GetStartEndDate(k, p.Period, startDate, endDate)
	// tStart, tEnd, e := helper.GetStartEndDate(k, p.Period, p.DateStart, p.DateEnd)

	// log.Printf("X. %#v | %#v", startDate.String(), endDate.String())

	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	match := tk.M{}

	match.Set("dateinfo.dateid", tk.M{"$gte": tStart, "$lt": tEnd})

	if p.Project != "" {
		anProject := strings.Split(p.Project, "(")
		match.Set("projectname", strings.TrimRight(anProject[0], " "))
	}

	if len(p.Turbine) > 0 {
		match.Set("turbine", tk.M{"$in": p.Turbine})
	}

	group := tk.M{
		"windspeed": tk.M{"$avg": "$avgwindspeed"},
	}

	groupID := tk.M{}

	/*if strings.ToLower(p.TimeBreakDown) == "date" {
		groupID.Set("dateid", "$dateinfo.dateid")
	} else if strings.ToLower(p.TimeBreakDown) == "monthly" {*/
	groupID.Set("monthid", "$dateinfo.monthid")
	groupID.Set("monthdesc", "$dateinfo.monthdesc")
	// }

	// if strings.ToLower(p.SeriesBreakDown) == "byturbine" {
	groupID.Set("turbine", "$turbine")
	// }

	group.Set("_id", groupID)

	pipes = append(pipes, tk.M{"$match": match})
	pipes = append(pipes, tk.M{"$group": group})
	pipes = append(pipes, tk.M{"$sort": tk.M{"_id.monthid": 1, "_id.turbine": 1}})

	csr, e := DB().Connection.NewQuery().
		From(new(ScadaData).TableName()).
		Command("pipe", pipes).
		Cursor(nil)

	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	e = csr.Fetch(&list, 0, false)

	csr.Close()

	// log.Printf("scada: %#v \n", list)
	// wra

	wraList := []tk.M{}
	wra := tk.M{}
	filter := []*dbox.Filter{}
	filter = append(filter, dbox.Eq("time", "Avg"))

	csr, e = DB().Connection.NewQuery().
		From(new(WRA).TableName()).
		Where(dbox.And(filter...)).
		Cursor(nil)

	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	e = csr.Fetch(&wraList, 0, false)

	csr.Close()

	if len(wraList) > 0 {
		wra = wraList[0]
	}

	// combine

	tmpRes := tk.M{}

	for _, val := range list {
		id := val.Get("_id").(tk.M)
		monthDesc := id.GetString("monthdesc")
		split := strings.Split(monthDesc, " ")
		trim := split[0][:3]
		reswra := wra.GetFloat64(strings.ToLower(trim))
		wind := val.GetFloat64("windspeed")
		turbine := id.GetString("turbine")

		details := []tk.M{}

		if tmpRes.Get(turbine) != nil {
			details = tmpRes.Get(turbine).([]tk.M)
		}

		time := tk.M{}
		time.Set("time", strings.ToUpper(trim)+" "+split[1])

		col := tk.M{}
		col.Set("WRA", reswra)
		col.Set("Onsite", wind)

		time.Set("col", col)

		details = append(details, time)

		tmpRes.Set(strings.Trim(turbine, " "), details)
	}

	turbineList := []string{}
	for key := range tmpRes {
		turbineList = append(turbineList, key)
	}

	sort.Strings(turbineList)
	for _, val := range turbineList {
		/*data = append(data, tk.M{
			"hours":   val,
			"details": tmpRes[val],
		})*/

		turbines = append(turbines, tk.M{}.Set("turbine", val).Set("details", tmpRes.Get(val)))
	}

	// met tower

	/*list = []tk.M{}

	match = tk.M{}

	match.Set("dateinfo.dateid", tk.M{"$gte": tStart, "$lte": tEnd})

	group = tk.M{
		"windspeed": tk.M{"$avg": "$vhubws90mavg"},
	}

	groupID = tk.M{}

	if strings.ToLower(p.TimeBreakDown) == "date" {
		groupID.Set("dateid", "$dateinfo.dateid")
	} else if strings.ToLower(p.TimeBreakDown) == "monthly" {
		groupID.Set("monthdesc", "$dateinfo.monthdesc")
	}

	group.Set("_id", groupID)

	pipes = []tk.M{}
	pipes = append(pipes, tk.M{"$match": match})
	pipes = append(pipes, tk.M{"$group": group})
	pipes = append(pipes, tk.M{"$sort": tk.M{"_id": 1}})

	csr, e = DB().Connection.NewQuery().
		From(new(MetTower).TableName()).
		Command("pipe", pipes).
		Cursor(nil)

	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	e = csr.Fetch(&list, 0, false)

	csr.Close()

	// log.Printf("met: %#v \n", list)

	for _, val := range list {
		id := val.Get("_id").(tk.M)
		turVal := tk.M{}

		if id.GetString("dateid") == "" {
			turVal.Set("time", id.Get("dateid").(time.Time).Format("01 Feb 2006"))
		} else {
			turVal.Set("time", id.GetString("monthdesc"))
		}

		wind := val.GetFloat64("windspeed")
		turVal.Set("name", "Met Tower")
		turVal.Set("value", wind)

		metTower = append(metTower, turVal)
	}*/

	data := struct {
		Data tk.M
	}{
		Data: tk.M{"turbine": turbines},
	}

	return helper.CreateResult(true, data, "success")
}

func (c *AnalyticMeteorologyController) Table1224(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson
	type Payload1224 struct {
		Project string
		Turbine []interface{}
	}

	p := new(Payload1224)
	e := k.GetPayload(&p)

	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	now := time.Now()
	last := time.Now().AddDate(0, -12, 0)

	tStart, _ := time.Parse("20060102", last.Format("200601")+"01")
	tEnd, _ := time.Parse("20060102", now.Format("200601")+"01")

	matchTurbine := tk.M{"dateinfo.dateid": tk.M{"$gte": tStart, "$lt": tEnd}}
	matchMet := tk.M{"dateinfo.dateid": tk.M{"$gte": tStart, "$lt": tEnd}}

	if p.Project != "" {
		anProject := strings.Split(p.Project, "(")
		matchTurbine.Set("projectname", strings.TrimRight(anProject[0], " "))
	}

	if len(p.Turbine) > 0 {
		matchTurbine.Set("turbine", tk.M{"$in": p.Turbine})
	}

	groupTurbine := tk.M{
		"windspeed": tk.M{"$avg": "$avgwindspeed"},
		"temp":      tk.M{"$avg": "$nacelletemperature"},
		"power":     tk.M{"$avg": "$power"},
	}
	tablenameTurbine := new(ScadaData).TableName()
	groupMet := tk.M{
		"windspeed": tk.M{"$avg": "$vhubws90mavg"},
		"temp":      tk.M{"$avg": "$thubhhubtemp855mavg"},
	}
	tablenameMet := new(MetTower).TableName()

	dataTurbine := processTableData(groupTurbine, matchTurbine, tablenameTurbine, "turbine")
	dataMet := processTableData(groupMet, matchMet, tablenameMet, "met")

	result := tk.M{"DataTurbine": dataTurbine, "DataMet": dataMet}

	return helper.CreateResult(true, result, "success")
}

func (c *AnalyticMeteorologyController) GetListMtbf(k *knot.WebContext) interface{}{

	k.Config.OutputType = knot.OutputJson

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

	var (
		query    []tk.M
		pipes    []tk.M
		datas    []tk.M
	)
	scadaOem := make([]tk.M, 0)

	query = append(query, tk.M{"_id": tk.M{"$ne": ""}})
	query = append(query, tk.M{"timestamp": tk.M{"$gte": tStart}})
	query = append(query, tk.M{"timestamp": tk.M{"$lte": tEnd}})
	if len(turbine) > 0 {
		query = append(query, tk.M{"turbine": tk.M{"$in": turbine}})
	}


	pipes = append(pipes, tk.M{"$match": tk.M{"$and": query}})
	pipes = append(pipes, tk.M{"$group": tk.M{"_id": "$turbine",
		"avgmttr": tk.M{"$avg": "$mttr"}, 
		"avgmttf": tk.M{"$avg": "$mttf"},
		"totmttf": tk.M{"$sum": "$mttf"},
		"totmachinedowntime": tk.M{"$sum": "$machinedowntime"},
		"totunknowndowntime": tk.M{"$sum": "$unknowndowntime"},
		"totgriddowntime": tk.M{"$sum": "$griddowntime"},
		},
	})
	pipes = append(pipes, tk.M{"$sort": tk.M{"_id": 1}})


	csr, e := DB().Connection.NewQuery().
		From(new(ScadaDataOEM).TableName()).
		Command("pipe", pipes).
		Cursor(nil)

	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	e = csr.Fetch(&scadaOem, 0, false)

	csr.Close()


	for _, m := range scadaOem {
		id := m.GetString("_id")
		avgmttr := m.GetFloat64("avgmttr")
		avgmttf := m.GetFloat64("avgmttf")
		totmttf := m.GetFloat64("totmttf")
		totmachinedowntime := m.GetFloat64("totmachinedowntime")
		totunknowndowntime := m.GetFloat64("totunknowndowntime")
		totgriddowntime := m.GetFloat64("totgriddowntime")
		totDowntime := totmachinedowntime + totunknowndowntime + totgriddowntime

		datas = append(datas, tk.M{
			"id": id,
			"avgmttr": avgmttr,
			"avgmttf": avgmttf,
			"totmttf": totmttf,
			"totmachinedowntime": totmachinedowntime,
			"totunknowndowntime": totunknowndowntime,
			"totgriddowntime": totgriddowntime,
			"totDowntime": totDowntime,
			"avgmtbf":     tk.Div(totmttf, totDowntime),
		})
	}

	if datas == nil {
		datas = make([]tk.M, 0)
	}

	return helper.CreateResult(true, datas, "success")
}

func processTableData(group, match tk.M, tablename, dataType string) (data []tk.M) {
	var (
		pipes []tk.M
		list  []tk.M
	)
	groupID := tk.M{}                           /*inside _id group*/
	groupID.Set("monthid", "$dateinfo.monthid") /*for sorting purpose*/
	groupID.Set("monthdesc", "$dateinfo.monthdesc")
	groupID.Set("hours", tk.M{"$dateToString": tk.M{"format": "%H:00", "date": "$timestamp"}}) /*to format HH:MM*/
	group.Set("_id", groupID)

	pipes = append(pipes, tk.M{"$match": match})
	pipes = append(pipes, tk.M{"$group": group})
	pipes = append(pipes, tk.M{"$sort": tk.M{"_id.monthid": 1, "_id.hours": 1}})

	csr, e := DB().Connection.NewQuery().
		From(tablename).
		Command("pipe", pipes).
		Cursor(nil)

	if e != nil {
		return
	}

	e = csr.Fetch(&list, 0, false)
	if e != nil {
		return
	}

	csr.Close()

	// combine

	tmpRes := tk.M{}

	for _, val := range list {
		id := val.Get("_id").(tk.M)
		monthDesc := id.GetString("monthdesc")
		split := strings.Split(monthDesc, " ")
		trim := split[0][:3]
		wind := 0.0
		if val.GetString("windspeed") != "NaN" {
			wind = val.GetFloat64("windspeed")
		}
		temp := 0.0
		if val.GetString("temp") != "NaN" {
			temp = val.GetFloat64("temp")
		}
		power := 0.0
		if dataType == "turbine" {
			power = val.GetFloat64("power")
		}

		hours := id.GetString("hours")

		details := []tk.M{}

		if tmpRes.Get(hours) != nil {
			details = tmpRes.Get(hours).([]tk.M)
		}

		time := tk.M{}
		time.Set("time", strings.ToUpper(trim)+" "+split[1])

		col := tk.M{}
		col.Set("WS", wind)
		col.Set("Temp", temp)
		if dataType == "turbine" {
			col.Set("Power", power)
		}

		time.Set("col", col)

		details = append(details, time)

		tmpRes.Set(hours, details)
	}

	hoursList := []string{}
	for key := range tmpRes {
		hoursList = append(hoursList, key)
	}
	sort.Strings(hoursList)
	for _, val := range hoursList {
		data = append(data, tk.M{
			"hours":   val,
			"details": tmpRes[val],
		})
	}

	return
}
