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

	type HeatMap struct {
		Color   string
		Opacity float64
	}

	var dataSeries []tk.M
	var dataHeat []tk.M

	p := new(PayloadAnalytic)
	e := k.GetPayload(&p)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}
	tStart, tEnd, e := helper.GetStartEndDate(k, p.Period, p.DateStart, p.DateEnd)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	query, queryMet := []tk.M{}, []tk.M{}
	pipes := []tk.M{}
	pipesmet := []tk.M{}
	query = append(query, tk.M{"_id": tk.M{"$ne": ""}})
	query = append(query, tk.M{"timestamp": tk.M{"$gte": tStart}})
	query = append(query, tk.M{"timestamp": tk.M{"$lte": tEnd}})

	if p.Project != "" {
		query = append(query, tk.M{"projectname": p.Project})
	}
	queryMet = append(queryMet, query...)
	queryMet = append(queryMet, tk.M{"vhubws90mavg": tk.M{"$gte": 0}})
	queryMet = append(queryMet, tk.M{"vhubws90mavg": tk.M{"$lte": 30}})

	query = append(query, tk.M{"avgwindspeed": tk.M{"$gte": 0}})
	query = append(query, tk.M{"avgwindspeed": tk.M{"$lte": 30}})

	pipes = append(pipes, tk.M{"$match": tk.M{"$and": query}})
	pipes = append(pipes, tk.M{"$project": tk.M{"turbine": 1, "avgwindspeed": 1, "timestamp": 1}})

	pipesmet = append(pipesmet, tk.M{"$match": tk.M{"$and": queryMet}})
	pipesmet = append(pipesmet, tk.M{"$project": tk.M{"vhubws90mavg": 1, "timestamp": 1}})

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

	turbineName, e := helper.GetTurbineNameList(p.Project)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	for _, _turbine := range pturbine {
		_tkm := tk.M{}.Set(tk.Sprintf("%#v", "Turbine"), turbineName[_turbine])
		for i := 1; i < len(arrturbine); i++ {
			_dt01 := allres.Get(_turbine, tk.M{}).(tk.M)
			_dt02 := allres.Get(arrturbine[i], tk.M{}).(tk.M)
			if arrturbine[i] != _turbine && len(_dt01) > 0 && len(_dt02) > 0 {
				_tkm.Set(tk.Sprintf("%#v", arrturbine[i]),
					GetCorrelation(_dt01, _dt02))
			} else {
				_tkm.Set(tk.Sprintf("%#v", arrturbine[i]), "-")
			}
		}
		dataSeries = append(dataSeries, _tkm)
	}

	for _, _tkm := range dataSeries {
		_heattkm := tk.M{}

		_aint := []float64{}

		for _key, _val := range _tkm {
			if tk.ToString(_val) == "-" || _key == "Turbine" {
				continue
			}
			_num := tk.ToFloat64(_val, 2, tk.RoundingAuto)
			if !tk.HasMember(_aint, _num) {
				_aint = append(_aint, _num)
			}
		}

		sort.Float64s(_aint)
		_mapunique := map[float64]float64{}
		_median := float64((len(_aint) + 1)) / 2
		for _i, _val := range _aint {
			_mapunique[_val] = float64(_i) + 1
		}

		// tk.Println("MAP : ", _mapunique, " MEDIAN : ", _median)

		for _key, _val := range _tkm {
			_dt := HeatMap{}
			_dt.Color = "white"
			_dt.Opacity = 1

			if tk.ToString(_val) != "-" && _key != "Turbine" {
				_fval := tk.ToFloat64(_val, 2, tk.RoundingAuto)
				if _median != _mapunique[_fval] {
					if _median > _mapunique[_fval] {
						_dt.Color = "red"
						_dt.Opacity = tk.Div((_median - _mapunique[_fval]), _median)
					} else {
						_dt.Color = "green"
						_dt.Opacity = tk.Div((_mapunique[_fval] - _median), _median)
					}
				}
			}

			_heattkm.Set(_key, _dt)
		}

		dataHeat = append(dataHeat, _heattkm)
	}

	data := struct {
		Column      []string
		Data        []tk.M
		Heat        []tk.M
		TurbineName map[string]string
	}{
		Column:      arrturbine,
		Data:        dataSeries,
		Heat:        dataHeat,
		TurbineName: turbineName,
	}

	return helper.CreateResult(true, data, "success")
}

func (m *AnalyticMeteorologyController) GetEnergyCorrelation(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	type HeatMap struct {
		Color   string
		Opacity float64
	}

	var dataSeries []tk.M
	var dataHeat []tk.M

	p := new(PayloadAnalytic)
	e := k.GetPayload(&p)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}
	tStart, tEnd, e := helper.GetStartEndDate(k, p.Period, p.DateStart, p.DateEnd)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	query, pipes := []tk.M{}, []tk.M{}
	query = append(query, tk.M{"timestamp": tk.M{"$gte": tStart}})
	query = append(query, tk.M{"timestamp": tk.M{"$lte": tEnd}})

	if p.Project != "" {
		query = append(query, tk.M{"projectname": p.Project})
	}

	query = append(query, tk.M{"isvalidstate": true})

	pipes = append(pipes, tk.M{"$match": tk.M{"$and": query}})
	pipes = append(pipes, tk.M{"$project": tk.M{"turbine": 1, "power": 1, "timestamp": 1, "statedescription": 1}})

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

		// lstatedesc := strings.ToLower(trx.StateDescription)
		// if strings.Contains(lstatedesc, "ready") || strings.Contains(lstatedesc, "wind") {
		// 	continue
		// }

		dkey := trx.TimeStamp.Format("20060102150405")

		_tkm := allres.Get(trx.Turbine, tk.M{}).(tk.M)
		if trx.Power != -99999.0 {
			_tkm.Set(dkey, trx.Power)
		}

		allres.Set(trx.Turbine, _tkm)
		_tturbine.Set(trx.Turbine, 1)
	}

	for key, _ := range _tturbine {
		arrturbine = append(arrturbine, key)
	}

	sort.Strings(arrturbine)
	pturbine := arrturbine
	arrturbine = append([]string{"Turbine"}, arrturbine...)

	if len(p.Turbine) > 0 {
		pturbine = []string{}
		for _, _v := range p.Turbine {
			pturbine = append(pturbine, tk.ToString(_v))
		}
	}

	turbineName, e := helper.GetTurbineNameList(p.Project)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	for _, _turbine := range pturbine {
		_tkm := tk.M{}.Set(tk.Sprintf("%#v", "Turbine"), turbineName[_turbine])
		for i := 1; i < len(arrturbine); i++ {
			_dt01 := allres.Get(_turbine, tk.M{}).(tk.M)
			_dt02 := allres.Get(arrturbine[i], tk.M{}).(tk.M)
			if arrturbine[i] != _turbine && len(_dt01) > 0 && len(_dt02) > 0 {
				_tkm.Set(tk.Sprintf("%#v", arrturbine[i]),
					GetCorrelation(_dt01, _dt02))
			} else {
				_tkm.Set(tk.Sprintf("%#v", arrturbine[i]), "-")
			}
		}
		dataSeries = append(dataSeries, _tkm)
	}

	for _, _tkm := range dataSeries {
		_heattkm := tk.M{}

		_aint := []float64{}

		for _key, _val := range _tkm {
			if tk.ToString(_val) == "-" || _key == "Turbine" {
				continue
			}
			_num := tk.ToFloat64(_val, 2, tk.RoundingAuto)
			if !tk.HasMember(_aint, _num) {
				_aint = append(_aint, _num)
			}
		}

		sort.Float64s(_aint)
		_mapunique := map[float64]float64{}
		_median := float64((len(_aint) + 1)) / 2
		for _i, _val := range _aint {
			_mapunique[_val] = float64(_i) + 1
		}

		// tk.Println("MAP : ", _mapunique, " MEDIAN : ", _median)

		for _key, _val := range _tkm {
			_dt := HeatMap{}
			_dt.Color = "white"
			_dt.Opacity = 1

			if tk.ToString(_val) != "-" && _key != "Turbine" {
				_fval := tk.ToFloat64(_val, 2, tk.RoundingAuto)
				if _median != _mapunique[_fval] {
					if _median > _mapunique[_fval] {
						_dt.Color = "red"
						_dt.Opacity = tk.Div((_median - _mapunique[_fval]), _median)
					} else {
						_dt.Color = "green"
						_dt.Opacity = tk.Div((_mapunique[_fval] - _median), _median)
					}
				}
			}

			_heattkm.Set(_key, _dt)
		}

		dataHeat = append(dataHeat, _heattkm)
	}

	sort.Strings(arrturbine)

	data := struct {
		Column      []string
		Data        []tk.M
		Heat        []tk.M
		TurbineName map[string]string
	}{
		Column:      arrturbine,
		Data:        dataSeries,
		Heat:        dataHeat,
		TurbineName: turbineName,
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
	match.Set("available", 1)

	if p.Project != "" {
		match.Set("projectname", p.Project)
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
	monthMap := map[string]int{}

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
		monthMap[strings.ToUpper(trim)+" "+split[1]] = 1

		col := tk.M{}
		col.Set("WRA", reswra)
		col.Set("Onsite", wind)

		time.Set("col", col)

		details = append(details, time)

		tmpRes.Set(strings.Trim(turbine, " "), details)
	}
	monthList := []string{}
	for key := range monthMap {
		monthList = append(monthList, key)
	}

	turbineList := []string{}
	for namaTurbine, val := range tmpRes {
		details := val.([]tk.M)
		monthExist := map[string]bool{}
		for _, detail := range details {
			monthExist[detail.GetString("time")] = true
		}
		for _, month := range monthList {
			_, hasMonth := monthExist[month]
			if !hasMonth {
				dummyItem := tk.M{
					"col":  tk.M{"WRA": 0.0, "Onsite": 0.0},
					"time": month,
				}
				details = append(details, dummyItem)
			}
		}
		tmpRes.Set(namaTurbine, details)
		turbineList = append(turbineList, namaTurbine)
	}

	sort.Strings(turbineList)
	turbineName, e := helper.GetTurbineNameList(p.Project)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	for _, val := range turbineList {
		/*data = append(data, tk.M{
			"hours":   val,
			"details": tmpRes[val],
		})*/

		turbines = append(turbines, tk.M{}.Set("turbine", turbineName[val]).Set("details", tmpRes.Get(val)))
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

	matchTurbine := []tk.M{
		tk.M{"dateinfo.dateid": tk.M{"$gte": tStart}},
		tk.M{"dateinfo.dateid": tk.M{"$lt": tEnd}},
		tk.M{"available": 1},
	}
	matchMet := tk.M{"dateinfo.dateid": tk.M{"$gte": tStart, "$lt": tEnd}}

	if p.Project != "" {
		matchTurbine = append(matchTurbine, tk.M{"projectname": p.Project})
		matchMet.Set("projectname", p.Project)
	}

	if len(p.Turbine) > 0 {
		matchTurbine = append(matchTurbine, tk.M{"turbine": tk.M{"$in": p.Turbine}})
	}
	matchScada := tk.M{"$and": matchTurbine}

	groupTurbine := tk.M{
		"windspeed": tk.M{"$avg": "$avgwindspeed"},
		"temp":      tk.M{"$avg": "$nacelletemperature"},
		"power":     tk.M{"$sum": "$power"},
	}
	tablenameTurbine := new(ScadaData).TableName()
	groupMet := tk.M{
		"windspeed": tk.M{"$avg": "$vhubws90mavg"},
		"temp":      tk.M{"$avg": "$thubhhubtemp855mavg"},
	}
	tablenameMet := new(MetTower).TableName()

	dataTurbine, totalTurbine, e := processTableData(groupTurbine, matchScada, tablenameTurbine, "turbine")
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}
	dataMet, totalMet, e := processTableData(groupMet, matchMet, tablenameMet, "met")
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	result := tk.M{"DataTurbine": dataTurbine, "DataMet": dataMet, "TotalTurbine": totalTurbine, "TotalMet": totalMet}

	return helper.CreateResult(true, result, "success")
}

func (c *AnalyticMeteorologyController) Graph1224(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson
	type PayloadGraph struct {
		Project string
		Turbine []interface{}
		Data    string
	}

	p := new(PayloadGraph)
	e := k.GetPayload(&p)

	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	now := time.Now()
	last := time.Now().AddDate(0, -12, 0)

	tStart, _ := time.Parse("20060102", last.Format("200601")+"01")
	tEnd, _ := time.Parse("20060102", now.Format("200601")+"01")

	var (
		data []tk.M
		err  error
	)

	match := []tk.M{
		tk.M{"dateinfo.dateid": tk.M{"$gte": tStart}},
		tk.M{"dateinfo.dateid": tk.M{"$lt": tEnd}},
	}

	if p.Project != "" {
		match = append(match, tk.M{"projectname": p.Project})
	}

	dtype := strings.ToLower(p.Data)

	if dtype == "turbine" || dtype == "" {
		match = append(match, tk.M{"available": 1})
		if len(p.Turbine) > 0 {
			match = append(match, tk.M{"turbine": tk.M{"$in": p.Turbine}})
		}
		//
		group := tk.M{
			"windspeed": tk.M{"$avg": "$avgwindspeed"},
			"temp":      tk.M{"$avg": "$nacelletemperature"},
			"power":     tk.M{"$sum": "$power"},
		}

		data, err = processGraphData(group, tk.M{"$and": match}, new(ScadaData).TableName(), "turbine")

	} else {
		group := tk.M{
			"windspeed": tk.M{"$avg": "$vhubws90mavg"},
			"temp":      tk.M{"$avg": "$thubhhubtemp855mavg"},
		}

		data, err = processGraphData(group, tk.M{"$and": match}, new(MetTower).TableName(), "met")
	}

	if err != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	return helper.CreateResult(true, data, "success")
}

func (c *AnalyticMeteorologyController) GetListMtbf(k *knot.WebContext) interface{} {

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
		query []tk.M
		pipes []tk.M
		datas []tk.M
	)
	scadaOem := make([]tk.M, 0)

	query = append(query, tk.M{"_id": tk.M{"$ne": ""}})
	query = append(query, tk.M{"dateinfo.dateid": tk.M{"$gte": tStart}})
	query = append(query, tk.M{"dateinfo.dateid": tk.M{"$lte": tEnd}})
	if len(turbine) > 0 {
		query = append(query, tk.M{"turbine": tk.M{"$in": turbine}})
	}

	pipes = append(pipes, tk.M{"$match": tk.M{"$and": query}})
	pipes = append(pipes, tk.M{"$group": tk.M{"_id": "$turbine",
		"downtimehours": tk.M{"$sum": "$downtimehours"},
		"oktime":        tk.M{"$sum": "$oktime"},
		"nooffailures":  tk.M{"$sum": "$nooffailures"},
	},
	})
	pipes = append(pipes, tk.M{"$sort": tk.M{"_id": 1}})

	// timenow := time.Now()
	csr, e := DB().Connection.NewQuery().
		From(new(ScadaSummaryDaily).TableName()).
		Command("pipe", pipes).
		Cursor(nil)
	//duration := time.Now().Sub(timenow).Seconds()
	// tk.Println("Duration 1: ", duration)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	//timenow = time.Now()
	e = csr.Fetch(&scadaOem, 0, false)
	// duration = time.Now().Sub(timenow).Seconds()
	// tk.Println("Duration 2: ", duration)

	csr.Close()

	turbineName, e := helper.GetTurbineNameList(p.Project)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	// timenow = time.Now()
	for _, m := range scadaOem {
		id := turbineName[m.GetString("_id")]

		oktime := m.GetFloat64("oktime") / 3600
		nooffailures := m.GetFloat64("nooffailures")
		downtimehours := m.GetFloat64("downtimehours")

		if nooffailures == 0 && downtimehours > 0 {
			nooffailures = 1
		}

		mtbf := tk.Div(oktime, nooffailures)
		mttr := tk.Div(downtimehours, nooffailures)
		mttf := mtbf - mttr

		datas = append(datas, tk.M{
			"id":             id,
			"mtbf":           mtbf,
			"mttr":           mttr,
			"mttf":           mttf,
			"totoptime":      oktime,
			"totdowntime":    downtimehours,
			"totnooffailure": nooffailures,
		})
	}
	// duration = time.Now().Sub(timenow).Seconds()
	// tk.Println("Duration 3: ", duration)

	if datas == nil {
		datas = make([]tk.M, 0)
	}

	return helper.CreateResult(true, datas, "success")
}

func processTableData(group, match tk.M, tablename, dataType string) (data []tk.M, totalData []tk.M, e error) {
	var (
		pipes []tk.M
		// list  []tk.M
		pipesTotal []tk.M
	)
	groupID := tk.M{}                           /*inside _id group*/
	groupID.Set("monthid", "$dateinfo.monthid") /*for sorting purpose*/
	groupID.Set("monthdesc", "$dateinfo.monthdesc")
	groupID.Set("hours", tk.M{"$dateToString": tk.M{"format": "%H:00", "date": "$timestamp"}}) /*to format HH:MM*/
	group.Set("_id", groupID)

	pipes = append(pipes, tk.M{"$match": match})
	pipes = append(pipes, tk.M{"$group": group})
	pipes = append(pipes, tk.M{"$sort": tk.M{"_id.monthid": 1}})

	// timenow := time.Now()
	csr, e := DB().Connection.NewQuery().
		From(tablename).
		Command("pipe", pipes).
		Cursor(nil)
	// duration := time.Now().Sub(timenow).Seconds()
	// tk.Println("Kondisi 1: ", duration)
	defer csr.Close()
	if e != nil {
		return
	}
	// combine

	tmpRes := tk.M{}
	wind := 0.0
	temp := 0.0
	power := 0.0
	monthDesc := ""
	totalData = []tk.M{}
	split := []string{}
	trim := ""
	var err error

	// timenow = time.Now()
	_list := tk.M{}
	for {
		_list = tk.M{}
		err = csr.Fetch(&_list, 1, false)
		if err != nil {
			break
		}

		id := _list.Get("_id").(tk.M)
		monthDesc = id.GetString("monthdesc")
		split = strings.Split(monthDesc, " ")
		trim = split[0][:3]

		wind = _list.GetFloat64("windspeed")
		temp = _list.GetFloat64("temp")
		if dataType == "turbine" {
			power = _list.GetFloat64("power") / 1000
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
	// duration = time.Now().Sub(timenow).Seconds()
	// tk.Println("Kondisi 2: ", duration)

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

	group.Set("_id", "$dateinfo.monthid")
	pipesTotal = append(pipesTotal, tk.M{"$match": match})
	pipesTotal = append(pipesTotal, tk.M{"$group": group})
	pipesTotal = append(pipesTotal, tk.M{"$sort": tk.M{"_id": 1}})

	// timenow = time.Now()
	csrTotal, e := DB().Connection.NewQuery().
		From(tablename).
		Command("pipe", pipesTotal).
		Cursor(nil)
	tk.Printf("%#v\n", tablename)
	tk.Printf("%#v\n", pipesTotal)
	// duration = time.Now().Sub(timenow).Seconds()
	// tk.Println("Kondisi 3: ", duration)
	defer csrTotal.Close()
	if e != nil {
		return
	}
	// timenow = time.Now()
	totalDataItem := tk.M{}
	for {
		totalDataItem = tk.M{}
		e = csrTotal.Fetch(&totalDataItem, 1, false)
		if e != nil {
			e = nil
			break
		}
		totalData = append(totalData, totalDataItem)
	}
	//e = csrTotal.Fetch(&totalData, 0, false)
	// duration = time.Now().Sub(timenow).Seconds()
	// tk.Println("Kondisi 4: ", duration)
	if e != nil {
		return
	}

	// timenow = time.Now()
	for idx := range totalData {
		if dataType == "turbine" {
			totalData[idx].Set("power", totalData[idx].GetFloat64("power")/1000)
		}
	}
	// duration = time.Now().Sub(timenow).Seconds()
	// tk.Println("Kondisi 5: ", duration)

	return
}

func processGraphData(group, match tk.M, tablename, dataType string) (data []tk.M, e error) {
	pipes := []tk.M{}

	groupID := tk.M{}                           /*inside _id group*/
	groupID.Set("monthid", "$dateinfo.monthid") /*for sorting purpose*/
	groupID.Set("monthdesc", "$dateinfo.monthdesc")
	groupID.Set("hours", tk.M{"$dateToString": tk.M{"format": "%H:00", "date": "$timestamp"}}) /*to format HH:MM*/
	group.Set("_id", groupID)

	pipes = append(pipes, tk.M{"$match": match})
	pipes = append(pipes, tk.M{"$group": group})
	pipes = append(pipes, tk.M{"$sort": tk.M{"_id.monthid": 1}})

	// timenow := time.Now()
	csr, e := DB().Connection.NewQuery().
		From(tablename).
		Command("pipe", pipes).
		Cursor(nil)
	// duration := time.Now().Sub(timenow).Seconds()
	// tk.Println("Kondisiong 1: ", duration)
	defer csr.Close()
	if e != nil {
		return
	}

	// timenow = time.Now()
	listmonth, listhour, listall := tk.M{}, tk.M{}, tk.M{}
	for {
		list := tk.M{}
		err := csr.Fetch(&list, 1, false)
		if err != nil {
			break
		}

		_dt := tk.M{}

		id := list.Get("_id").(tk.M)
		sMonthDesc := strings.Split(id.GetString("monthdesc"), " ")
		if len(sMonthDesc) > 1 {
			_dt.Set("time", strings.ToUpper(sMonthDesc[0][:3])+" "+sMonthDesc[1])
			_dt.Set("timeint", tk.Sprintf("%d01", id.GetInt("monthid")))
		}

		_dt.Set("hours", id.GetString("hours"))
		_dt.Set("ws", list.GetFloat64("windspeed"))
		_dt.Set("temp", list.GetFloat64("temp"))
		_dt.Set("power", 0)

		listhour.Set(id.GetString("hours"), 1)
		listmonth.Set(_dt.GetString("time"), _dt.GetString("timeint"))
		listall.Set(tk.Sprintf("%s_%s", _dt.GetString("time"), id.GetString("hours")), 1)

		if dataType == "turbine" {
			_dt.Set("power", list.GetFloat64("power")/1000)
		}

		data = append(data, _dt)
	}
	// duration = time.Now().Sub(timenow).Seconds()
	// tk.Println("Kondisiong 2: ", duration)

	// timenow = time.Now()
	for month, _ := range listmonth {
		for hour, _ := range listhour {
			if !listall.Has(tk.Sprintf("%s_%s", month, hour)) {
				_dt := tk.M{}
				_dt.Set("time", month)
				_dt.Set("timeint", listmonth.GetString(month))
				_dt.Set("hours", hour)
				_dt.Set("ws", nil)
				_dt.Set("temp", nil)
				_dt.Set("power", nil)
				data = append(data, _dt)
			}
		}
	}
	// duration = time.Now().Sub(timenow).Seconds()
	// tk.Println("Kondisiong 3: ", duration)

	return
}
