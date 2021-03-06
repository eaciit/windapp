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

	for k, _ := range turbineName {
		isappend := true
		for _, val := range arrturbine {
			if val == k {
				isappend = false
				break
			}
		}

		if isappend {
			arrturbine = append(arrturbine, k)
		}
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
		pipes    []tk.M
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

	tStart := last.Year()*100 + int(last.Month())
	tEnd := now.Year()*100 + int(now.Month())

	matches := []tk.M{
		tk.M{"monthid": tk.M{"$gte": tStart}},
		tk.M{"monthid": tk.M{"$lt": tEnd}},
		tk.M{"type": "SCADA"},
	}

	if p.Project != "" {
		matches = append(matches, tk.M{"projectname": p.Project})
	}

	if len(p.Turbine) > 0 {
		matches = append(matches, tk.M{"turbine": tk.M{"$in": p.Turbine}})
	}

	group := tk.M{
		"_id": tk.M{
			"monthid":   "$monthid",
			"monthdesc": "$monthdesc",
			"turbine":   "$turbine",
		},
		"windspeedtotal": tk.M{"$sum": "$windspeedtotal"},
		"windspeedcount": tk.M{"$sum": "$windspeedcount"},
	}

	pipes = append(pipes, tk.M{"$match": tk.M{"$and": matches}})
	pipes = append(pipes, tk.M{"$group": group})
	pipes = append(pipes, tk.M{"$sort": tk.M{"_id.monthid": 1, "_id.turbine": 1}})

	csr, e := DB().Connection.NewQuery().
		From("rpt_dineuralprofile").
		Command("pipe", pipes).
		Cursor(nil)

	defer csr.Close()

	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	for {
		item := tk.M{}
		e = csr.Fetch(&item, 1, false)
		if e != nil {
			break
		}
		list = append(list, item)
	}
	// wra

	wraList := []tk.M{}
	wra := tk.M{}
	filter := []*dbox.Filter{}
	filter = append(filter, dbox.Eq("time", "Avg"))

	csr, e = DB().Connection.NewQuery().
		From(new(WRA).TableName()).
		Where(dbox.And(filter...)).
		Cursor(nil)

	defer csr.Close()

	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	for {
		item := tk.M{}
		e = csr.Fetch(&item, 1, false)
		if e != nil {
			break
		}
		wraList = append(wraList, item)
	}

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
		wind := tk.Div(val.GetFloat64("windspeedtotal"), val.GetFloat64("windspeedcount"))
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
		turbines = append(turbines, tk.M{}.Set("turbine", turbineName[val]).Set("details", tmpRes.Get(val)))
	}

	data := struct {
		Data tk.M
	}{
		Data: tk.M{"turbine": turbines},
	}

	return helper.CreateResult(true, data, "success")
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

	tStart := last.Year()*100 + int(last.Month())
	tEnd := now.Year()*100 + int(now.Month())

	matches := []tk.M{
		tk.M{"monthid": tk.M{"$gte": tStart}},
		tk.M{"monthid": tk.M{"$lt": tEnd}},
	}

	totalIndex := 1

	if p.Project != "" {
		totalIndex++
		matches = append(matches, tk.M{"projectname": p.Project})
	}

	if len(p.Turbine) > 0 {
		totalIndex++
		matches = append(matches, tk.M{"turbine": tk.M{"$in": p.Turbine}})
	}

	group := tk.M{
		"_id": tk.M{
			"monthid":   "$monthid",
			"monthdesc": "$monthdesc",
			"hours":     "$hours",
		},
		"windspeedtotal":   tk.M{"$sum": "$windspeedtotal"},
		"temperaturetotal": tk.M{"$sum": "$temperaturetotal"},
		"powertotal":       tk.M{"$sum": "$powertotal"},
		"windspeedcount":   tk.M{"$sum": "$windspeedcount"},
		"temperaturecount": tk.M{"$sum": "$temperaturecount"},
	}

	dataTurbine, err := processGraphData(group, matches, "rpt_dineuralprofile", "turbine")
	dataMet, err := processGraphData(group, matches[0:totalIndex], "rpt_dineuralprofile", "met")

	if err != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	result := tk.M{"DataTurbine": dataTurbine, "DataMet": dataMet}

	return helper.CreateResult(true, result, "success")
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

	tStart := last.Year()*100 + int(last.Month())
	tEnd := now.Year()*100 + int(now.Month())

	matches := []tk.M{
		tk.M{"monthid": tk.M{"$gte": tStart}},
		tk.M{"monthid": tk.M{"$lt": tEnd}},
	}
	totalIndex := 1

	if p.Project != "" {
		totalIndex++
		matches = append(matches, tk.M{"projectname": p.Project})
	}

	if len(p.Turbine) > 0 {
		totalIndex++
		matches = append(matches, tk.M{"turbine": tk.M{"$in": p.Turbine}})
	}

	group := tk.M{
		"_id": tk.M{
			"monthid":   "$monthid",
			"monthdesc": "$monthdesc",
			"hours":     "$hours",
		},
		"windspeedtotal":   tk.M{"$sum": "$windspeedtotal"},
		"temperaturetotal": tk.M{"$sum": "$temperaturetotal"},
		"powertotal":       tk.M{"$sum": "$powertotal"},
		"windspeedcount":   tk.M{"$sum": "$windspeedcount"},
		"temperaturecount": tk.M{"$sum": "$temperaturecount"},
	}

	dataTurbine, totalTurbine, e := processTableData(group, matches, "rpt_dineuralprofile", "turbine")
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}
	dataMet, totalMet, e := processTableData(group, matches[0:totalIndex], "rpt_dineuralprofile", "met")
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	result := tk.M{"DataTurbine": dataTurbine, "DataMet": dataMet, "TotalTurbine": totalTurbine, "TotalMet": totalMet}

	return helper.CreateResult(true, result, "success")
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

func processTableData(group tk.M, matches []tk.M, tablename, dataType string) (data []tk.M, totalData []tk.M, e error) {
	if dataType == "turbine" {
		matches = append(matches, tk.M{"type": "SCADA"})
	} else {
		matches = append(matches, tk.M{"type": "MET"})
	}
	var pipes []tk.M

	pipes = append(pipes, tk.M{"$match": tk.M{"$and": matches}})
	pipes = append(pipes, tk.M{"$group": group})
	pipes = append(pipes, tk.M{"$sort": tk.M{"_id.monthid": 1}})

	csr, e := DB().Connection.NewQuery().
		From(tablename).
		Command("pipe", pipes).
		Cursor(nil)
	defer csr.Close()
	if e != nil {
		return
	}
	// combine

	tmpRes := tk.M{}
	windTotal := 0.0
	tempTotal := 0.0
	powerTotal := 0.0
	windCount := 0.0
	tempCount := 0.0
	monthDesc := ""
	totalData = []tk.M{}
	split := []string{}
	trim := ""
	var err error
	lastMonthID := 0
	currMonthID := 0

	_list := tk.M{}
	for {
		_list = tk.M{}
		err = csr.Fetch(&_list, 1, false)
		if err != nil {
			break
		}

		id := _list.Get("_id").(tk.M)
		monthDesc = id.GetString("monthdesc")
		currMonthID = id.GetInt("monthid")
		split = strings.Split(monthDesc, " ")
		trim = split[0][:3]

		if lastMonthID != currMonthID {
			if lastMonthID != 0 {
				windAvg := tk.Div(windTotal, windCount)
				tempAvg := tk.Div(tempTotal, tempCount)
				totalData = append(totalData, tk.M{
					"_id":       lastMonthID,
					"power":     powerTotal / 1000,
					"temp":      tempAvg,
					"windspeed": windAvg,
				})
			}
			windTotal = 0.0
			tempTotal = 0.0
			powerTotal = 0.0
			windCount = 0.0
			tempCount = 0.0
			lastMonthID = currMonthID
		}
		windTotal += _list.GetFloat64("windspeedtotal")
		windCount += _list.GetFloat64("windspeedcount")
		tempTotal += _list.GetFloat64("temperaturetotal")
		tempCount += _list.GetFloat64("temperaturecount")
		powerTotal += _list.GetFloat64("powertotal")

		hours := id.GetString("hours")
		details := []tk.M{}

		if tmpRes.Get(hours) != nil {
			details = tmpRes.Get(hours).([]tk.M)
		}

		time := tk.M{}
		time.Set("time", strings.ToUpper(trim)+" "+split[1])

		col := tk.M{}
		col.Set("WS", tk.Div(_list.GetFloat64("windspeedtotal"), _list.GetFloat64("windspeedcount")))
		col.Set("Temp", tk.Div(_list.GetFloat64("temperaturetotal"), _list.GetFloat64("temperaturecount")))
		if dataType == "turbine" {
			col.Set("Power", _list.GetFloat64("powertotal")/1000.0)
		}

		time.Set("col", col)

		details = append(details, time)

		tmpRes.Set(hours, details)
	}
	if lastMonthID != 0 {
		windAvg := tk.Div(windTotal, windCount)
		tempAvg := tk.Div(tempTotal, tempCount)
		totalData = append(totalData, tk.M{
			"_id":       lastMonthID,
			"power":     powerTotal / 1000,
			"temp":      tempAvg,
			"windspeed": windAvg,
		})
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

func processGraphData(group tk.M, matches []tk.M, tablename, dataType string) (data []tk.M, e error) {
	if dataType == "turbine" {
		matches = append(matches, tk.M{"type": "SCADA"})
	} else {
		matches = append(matches, tk.M{"type": "MET"})
	}

	pipes := []tk.M{}
	pipes = append(pipes, tk.M{"$match": tk.M{"$and": matches}})
	pipes = append(pipes, tk.M{"$group": group})
	pipes = append(pipes, tk.M{"$sort": tk.M{"_id.monthid": 1}})

	csr, e := DB().Connection.NewQuery().
		From(tablename).
		Command("pipe", pipes).
		Cursor(nil)
	defer csr.Close()
	if e != nil {
		return
	}

	data = []tk.M{}
	listmonth, listhour, listall := tk.M{}, tk.M{}, tk.M{}
	for {
		list := tk.M{}
		err := csr.Fetch(&list, 1, false)
		if err != nil {
			break
		}

		ids := list.Get("_id").(tk.M)
		_dt := tk.M{
			"time":    ids.GetString("monthdesc"),
			"timeint": tk.Sprintf("%d01", ids.GetInt("monthid")),
			"hours":   ids.GetString("hours"),
			"ws":      tk.Div(list.GetFloat64("windspeedtotal"), list.GetFloat64("windspeedcount")),
			"temp":    tk.Div(list.GetFloat64("temperaturetotal"), list.GetFloat64("temperaturecount")),
			"power":   0,
		}

		listhour.Set(ids.GetString("hours"), 1)
		listmonth.Set(_dt.GetString("time"), _dt.GetString("timeint"))
		listall.Set(tk.Sprintf("%s_%s", _dt.GetString("time"), ids.GetString("hours")), 1)

		if dataType == "turbine" {
			_dt.Set("power", list.GetFloat64("powertotal")/1000)
		}

		data = append(data, _dt)
	}

	for month, _ := range listmonth {
		for hour, _ := range listhour {
			if !listall.Has(tk.Sprintf("%s_%s", month, hour)) {
				_dt := tk.M{
					"time":    month,
					"timeint": listmonth.GetString(month),
					"hours":   hour,
					"ws":      nil,
					"temp":    nil,
					"power":   nil,
				}
				data = append(data, _dt)
			}
		}
	}

	return
}
