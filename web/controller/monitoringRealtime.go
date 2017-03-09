package controller

import (
	. "eaciit/wfdemo-git/library/core"
	// . "eaciit/wfdemo-git/library/helper"
	. "eaciit/wfdemo-git/library/models"

	"eaciit/wfdemo-git/web/helper"

	cr "github.com/eaciit/crowd"
	"github.com/eaciit/dbox"

	"github.com/eaciit/knot/knot.v1"
	tk "github.com/eaciit/toolkit"

	"strings"
	"time"
)

type MonitoringRealtimeController struct {
	App
}

func CreateMonitoringRealtimeController() *MonitoringRealtimeController {
	var controller = new(MonitoringRealtimeController)
	return controller
}

var (
	defaultValue = -999999.0
)

type MiniScada struct {
	NacellePosition float64
	WindSpeed       float64
	Turbine         string
}

func (c *MonitoringRealtimeController) GetDataProject(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson
	k.Config.NoLog = true

	results := c.GetMonitoringByProject("Tejuva")

	return helper.CreateResult(true, results, "success")
}

func (c *MonitoringRealtimeController) getValue() float64 {
	retVal := 0.0

	return retVal
}

func (c *MonitoringRealtimeController) GetMonitoringByProject(project string) (rtkm tk.M) {

	rtkm = tk.M{}

	csrt, err := DB().Connection.NewQuery().Select("turbineid", "feeder").
		From("ref_turbine").
		Where(dbox.Eq("project", project)).Cursor(nil)

	if err != nil {
		tk.Println(err.Error())
	}

	_result := []tk.M{}
	err = csrt.Fetch(&_result, 0, false)
	if err != nil {
		tk.Println(err.Error())
	}
	csrt.Close()

	alldata := tk.M{}
	arrfield := []string{"ActivePower", "WindSpeed", "WindDirection", "NacellePosition", "Temperature",
		"PitchAngle", "RotorRPM"}
	lastUpdate := time.Time{}
	PowerGen, AvgWindSpeed, CountWS := float64(0), float64(0), float64(0)
	turbinedown := 0

	arrturbinestatus := getTurbineStatus(project)

	for _, _tkm := range _result {
		aturbine := tk.M{}
		strturbine := _tkm.GetString("turbineid")
		aturbine.Set("Turbine", strturbine)
		aturbine.Set("DataComing", 0)

		for _, afield := range arrfield {
			aturbine.Set(afield, defaultValue)

			_tlafield := strings.ToLower(afield)
			icsrt, err := DB().Connection.NewQuery().Select("timestamp", _tlafield).From(new(ScadaRealTime).TableName()).
				Where(dbox.And(dbox.Eq("turbine", strturbine), dbox.Ne(_tlafield, defaultValue), dbox.Eq("projectname", project))).
				Order("-timestamp").Cursor(nil)
			if err != nil {
				tk.Println(err.Error())
			}

			_tdata := tk.M{}
			if icsrt.Count() > 0 {
				err = icsrt.Fetch(&_tdata, 1, false)
			}
			if err != nil {
				tk.Println(err.Error())
			}
			icsrt.Close()

			ifloat := _tdata.GetFloat64(_tlafield)
			if len(_tdata) > 0 && ifloat != defaultValue {
				tstamp := _tdata.Get("timestamp", time.Time{}).(time.Time)
				utime := aturbine.Get("TimeUpdate", time.Time{}).(time.Time)
				aturbine.Set(afield, ifloat)
				aturbine.Set("DataComing", 1)

				if tstamp.After(utime) {
					aturbine.Set("TimeUpdate", tstamp)
				}

				if tstamp.After(lastUpdate) {
					lastUpdate = tstamp
				}

				switch afield {
				case "ActivePower":
					PowerGen += ifloat
				case "WindSpeed":
					AvgWindSpeed += ifloat
					CountWS += 1
				}
			}
		}

		aturbine.Set("AlarmCode", arrturbinestatus[strturbine].AlarmCode).
			Set("AlarmDesc", arrturbinestatus[strturbine].AlarmDesc).
			Set("Status", arrturbinestatus[strturbine].Status).
			Set("AlarmUpdate", arrturbinestatus[strturbine].TimeUpdate)
		if arrturbinestatus[strturbine].Status == 0 {
			turbinedown += 1
		}

		arrturbine := alldata.Get(_tkm.GetString("feeder"), []tk.M{}).([]tk.M)
		arrturbine = append(arrturbine, aturbine)
		alldata.Set(_tkm.GetString("feeder"), arrturbine)
	}

	rtkm.Set("Data", alldata)
	rtkm.Set("TimeStamp", lastUpdate)
	rtkm.Set("PowerGeneration", PowerGen)
	rtkm.Set("AvgWindSpeed", tk.Div(AvgWindSpeed, CountWS))
	rtkm.Set("PLF", tk.Div(PowerGen, (50400*100)))
	rtkm.Set("TurbineActive", len(_result)-turbinedown)
	rtkm.Set("TurbineDown", turbinedown)

	return
}

func (c *MonitoringRealtimeController) GetDataAlarm(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson
	k.Config.NoLog = true

	p := new(helper.Payloads)
	err := k.GetPayload(&p)
	if err != nil {
		return helper.CreateResult(false, nil, err.Error())
	}

	query := tk.M{}.Set("where", dbox.Eq("projectname", "Tejuva"))
	csr, err := DB().Find(new(AlarmRealtime), query)
	if err != nil {
		return helper.CreateResult(false, nil, err.Error())
	}
	totalData := csr.Count()
	csr.Close()

	query.Set("take", p.Take).Set("skip", p.Skip)
	csr, err = DB().Connection.NewQuery().From(new(AlarmRealtime).TableName()).
		Skip(p.Skip).Take(p.Take).Cursor(nil)
	if err != nil {
		return helper.CreateResult(false, nil, err.Error())
	}

	results := make([]AlarmRealtime, 0)
	err = csr.Fetch(&results, 0, false)
	if err != nil {
		return helper.CreateResult(false, nil, err.Error())
	}
	csr.Close()

	return helper.CreateResult(true, tk.M{}.Set("Data", results).Set("Total", totalData), "success")
}

func (c *MonitoringRealtimeController) GetDataTurbine(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson
	k.Config.NoLog = true
	sessid := k.Session("sessionid", "")
	accs := "GetDataTurbine"

	p := struct {
		Turbine string
	}{}
	e := k.GetPayload(&p)
	if e != nil {
		WriteLog(sessid, accs, e.Error())
	}

	power := 0.0
	windSpeed := 0.0
	cWindSpeed := 0
	windDir := 0.0
	cWindDir := 0
	nacellePos := 0.0
	cNacellePos := 0
	temperature := 0.0
	cTemperature := 0
	pitch := 0.0
	cPitch := 0
	rotor := 0.0
	cRotor := 0

	t := p.Turbine

	var detail ScadaMonitoringItem
	detail.Turbine = t

	csrt, err := DB().Connection.NewQuery().From(new(ScadaRealTime).TableName()).
		Where(dbox.And(dbox.Eq("turbine", t), dbox.Ne("activepower", defaultValue))).
		Order("-timestamp").Cursor(nil)
	if err != nil {
		tk.Println(err.Error())
	}
	results := make([]ScadaRealTime, 0)
	err = csrt.Fetch(&results, 1, false)
	if err != nil {
		tk.Println(err.Error())
	}
	csrt.Close()

	if len(results) > 0 {
		result := results[0]
		detail.ActivePower = result.ActivePower
	} else {
		detail.ActivePower = defaultValue
	}

	if detail.ActivePower > defaultValue {
		power += detail.ActivePower
	}

	csrt, err = DB().Connection.NewQuery().From(new(ScadaRealTime).TableName()).
		Where(dbox.And(dbox.Eq("turbine", t), dbox.Ne("windspeed", defaultValue))).
		Order("-timestamp").Cursor(nil)
	if err != nil {
		tk.Println(err.Error())
	}
	results = make([]ScadaRealTime, 0)
	err = csrt.Fetch(&results, 1, false)
	if err != nil {
		tk.Println(err.Error())
	}
	csrt.Close()

	if len(results) > 0 {
		result := results[0]
		detail.WindSpeed = result.WindSpeed
	} else {
		detail.WindSpeed = defaultValue
	}

	if detail.WindSpeed != defaultValue {
		windSpeed += detail.WindSpeed
		cWindSpeed++
	}

	csrt, err = DB().Connection.NewQuery().From(new(ScadaRealTime).TableName()).
		Where(dbox.And(dbox.Eq("turbine", t), dbox.Ne("winddirection", defaultValue))).
		Order("-timestamp").Cursor(nil)
	if err != nil {
		tk.Println(err.Error())
	}
	results = make([]ScadaRealTime, 0)
	err = csrt.Fetch(&results, 1, false)
	if err != nil {
		tk.Println(err.Error())
	}
	csrt.Close()

	if len(results) > 0 {
		result := results[0]
		detail.WindDirection = result.WindDirection
	} else {
		detail.WindDirection = defaultValue
	}

	if detail.WindDirection != defaultValue {
		windDir += detail.WindDirection
		cWindDir++
	}

	csrt, err = DB().Connection.NewQuery().From(new(ScadaRealTime).TableName()).
		Where(dbox.And(dbox.Eq("turbine", t), dbox.Ne("nacelleposition", defaultValue))).
		Order("-timestamp").Cursor(nil)
	if err != nil {
		tk.Println(err.Error())
	}
	results = make([]ScadaRealTime, 0)
	err = csrt.Fetch(&results, 1, false)
	if err != nil {
		tk.Println(err.Error())
	}
	csrt.Close()

	if len(results) > 0 {
		result := results[0]
		detail.NacellePosition = result.NacellePosition
	} else {
		detail.NacellePosition = defaultValue
	}

	if detail.NacellePosition != defaultValue {
		nacellePos += detail.NacellePosition
		cNacellePos++
	}

	csrt, err = DB().Connection.NewQuery().From(new(ScadaRealTime).TableName()).
		Where(dbox.And(dbox.Eq("turbine", t), dbox.Ne("temperature", defaultValue))).
		Order("-timestamp").Cursor(nil)
	if err != nil {
		tk.Println(err.Error())
	}
	results = make([]ScadaRealTime, 0)
	err = csrt.Fetch(&results, 1, false)
	if err != nil {
		tk.Println(err.Error())
	}
	csrt.Close()

	if len(results) > 0 {
		result := results[0]
		detail.Temperature = result.Temperature
	} else {
		detail.Temperature = defaultValue
	}

	if detail.Temperature != defaultValue {
		temperature += detail.Temperature
		cTemperature++
	}

	csrt, err = DB().Connection.NewQuery().From(new(ScadaRealTime).TableName()).
		Where(dbox.And(dbox.Eq("turbine", t), dbox.Ne("pitchangle", defaultValue))).
		Order("-timestamp").Cursor(nil)
	if err != nil {
		tk.Println(err.Error())
	}
	results = make([]ScadaRealTime, 0)
	err = csrt.Fetch(&results, 1, false)
	if err != nil {
		tk.Println(err.Error())
	}
	csrt.Close()

	if len(results) > 0 {
		result := results[0]
		detail.PitchAngle = result.PitchAngle
	} else {
		detail.PitchAngle = defaultValue
	}

	if detail.PitchAngle != defaultValue {
		pitch += detail.PitchAngle
		cPitch++
	}

	csrt, err = DB().Connection.NewQuery().From(new(ScadaRealTime).TableName()).
		Where(dbox.And(dbox.Eq("turbine", t), dbox.Ne("rotorrpm", defaultValue))).
		Order("-timestamp").Cursor(nil)
	if err != nil {
		tk.Println(err.Error())
	}
	results = make([]ScadaRealTime, 0)
	err = csrt.Fetch(&results, 1, false)
	if err != nil {
		tk.Println(err.Error())
	}
	csrt.Close()

	if len(results) > 0 {
		result := results[0]
		detail.RotorRPM = result.RotorRPM
	} else {
		detail.RotorRPM = defaultValue
	}

	if detail.RotorRPM != defaultValue {
		rotor += detail.RotorRPM
		cRotor++
	}

	return detail
}

func (c *MonitoringRealtimeController) GetWindRoseData(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson
	k.Config.NoLog = true
	sessid := k.Session("sessionid", "")
	accs := "GetWindRoseData"

	// WindRoseResult = []tk.M{}

	p := struct {
		Turbine string
	}{}
	e := k.GetPayload(&p)
	if e != nil {
		WriteLog(sessid, accs, e.Error())
	}

	query := []tk.M{}
	pipes := []tk.M{}
	section = 12
	getFullWSCategory()

	data := []MiniScada{}
	_data := MiniScada{}

	lastDateData, e := time.Parse(time.RFC3339, "2017-01-22T00:00:00+00:00")
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	turbines := p.Turbine
	defaultValue := -999999.00

	groupdata := tk.M{}
	groupdata.Set("Name", turbines)

	query = append(query, tk.M{"_id": tk.M{"$ne": nil}})
	query = append(query, tk.M{"nacelleposition": tk.M{"$ne": defaultValue}})
	query = append(query, tk.M{"dateinfo.dateid": lastDateData})
	query = append(query, tk.M{"turbine": turbines})
	pipes = append(pipes, tk.M{"$match": tk.M{"$and": query}})
	pipes = append(pipes, tk.M{"$project": tk.M{"nacelleposition": 1, "windspeed": 1}})
	csr, e := DB().Connection.NewQuery().From(new(ScadaRealTime).TableName()).
		Command("pipe", pipes).Cursor(nil)

	for {
		e = csr.Fetch(&_data, 1, false)
		if e != nil {
			break
		}
		data = append(data, _data)
	}
	csr.Close()

	if tk.SliceLen(data) > 0 {
		totalDuration := float64(len(data)) /* Tot data * 2 for get total minutes*/
		datas := cr.From(&data).Apply(func(x interface{}) interface{} {
			dt := x.(MiniScada)
			var di DataItems

			dirNo, dirDesc := getDirection(dt.NacellePosition, section)
			wsNo, wsDesc := getWsCategory(dt.WindSpeed)

			di.DirectionNo = dirNo
			di.DirectionDesc = dirDesc
			di.WsCategoryNo = wsNo
			di.WsCategoryDesc = wsDesc
			di.Frequency = 1

			return di
		}).Exec().Group(func(x interface{}) interface{} {
			dt := x.(DataItems)

			var dig DataItemsGroup
			dig.DirectionNo = dt.DirectionNo
			dig.DirectionDesc = dt.DirectionDesc
			dig.WsCategoryNo = dt.WsCategoryNo
			dig.WsCategoryDesc = dt.WsCategoryDesc

			return dig
		}, nil).Exec()

		dts := datas.Apply(func(x interface{}) interface{} {
			kv := x.(cr.KV)
			vv := kv.Key.(DataItemsGroup)
			vs := kv.Value.([]DataItems)

			sumFreq := cr.From(&vs).Sum(func(x interface{}) interface{} {
				dt := x.(DataItems)
				return dt.Frequency
			}).Exec().Result.Sum

			var di DataItemsResult

			di.DirectionNo = vv.DirectionNo
			di.DirectionDesc = vv.DirectionDesc
			di.WsCategoryNo = vv.WsCategoryNo
			di.WsCategoryDesc = vv.WsCategoryDesc
			di.Hours = tk.Div(sumFreq, 6.0)
			di.Contribution = tk.RoundingAuto64(tk.Div(sumFreq, totalDuration)*100.0, 2)

			// key := turbines + "_" + tk.ToString(di.DirectionNo)

			// if !tkMaxVal.Has(key) {
			// 	tkMaxVal.Set(key, di.Contribution)
			// } else {
			// 	tkMaxVal.Set(key, tkMaxVal.GetFloat64(key)+di.Contribution)
			// }

			di.Frequency = int(sumFreq)

			return di
		}).Exec()

		results := dts.Result.Data().([]DataItemsResult)
		wsCategoryList := []string{}
		for _, dataRes := range results {
			wsCategoryList = append(wsCategoryList, tk.ToString(dataRes.DirectionNo)+
				"_"+tk.ToString(dataRes.WsCategoryNo)+"_"+dataRes.WsCategoryDesc)
		}
		splitCatList := []string{}
		for _, wsCat := range fullWSCatList {
			if !tk.HasMember(wsCategoryList, wsCat) {
				splitCatList = strings.Split(wsCat, "_")
				emptyRes := DataItemsResult{}
				emptyRes.DirectionNo = tk.ToInt(splitCatList[0], tk.RoundingAuto)
				divider := section

				emptyRes.DirectionDesc = (360 / divider) * emptyRes.DirectionNo
				emptyRes.WsCategoryNo = tk.ToInt(splitCatList[1], tk.RoundingAuto)
				emptyRes.WsCategoryDesc = splitCatList[2]
				results = append(results, emptyRes)
			}
		}
		groupdata.Set("Data", results)

		// tk.Printf("results : %s \n", tk.SliceLen(results))
		// tk.Printf("fullWSCatList : %s \n", fullWSCatList)

		// WindRoseResult = append(WindRoseResult, groupdata)
	} else {
		splitCatList := []string{}
		results := []DataItemsResult{}
		for _, wsCat := range fullWSCatList {
			splitCatList = strings.Split(wsCat, "_")
			emptyRes := DataItemsResult{}
			emptyRes.DirectionNo = tk.ToInt(splitCatList[0], tk.RoundingAuto)
			divider := section

			emptyRes.DirectionDesc = (360 / divider) * emptyRes.DirectionNo
			emptyRes.WsCategoryNo = tk.ToInt(splitCatList[1], tk.RoundingAuto)
			emptyRes.WsCategoryDesc = splitCatList[2]
			results = append(results, emptyRes)
		}
		groupdata.Set("Data", results)
		// WindRoseResult = append(WindRoseResult, groupdata)
	}

	// tk.Printf("groupdata : %s \n", tk.SliceLen(groupdata))

	dataresult := struct {
		WindRose tk.M
	}{
		WindRose: groupdata,
	}

	return helper.CreateResult(true, dataresult, "success")
}

func (c *MonitoringRealtimeController) GetDataLine(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson
	k.Config.NoLog = true
	sessid := k.Session("sessionid", "")
	accs := "GetDataLine"

	var (
		pipes      []tk.M
		filter     []*dbox.Filter
		list       []tk.M
		dataSeries []tk.M
	)

	p := struct {
		Turbine string
	}{}
	e := k.GetPayload(&p)
	if e != nil {
		WriteLog(sessid, accs, e.Error())
	}

	lastDateData, e := time.Parse(time.RFC3339, "2017-01-22T00:00:00+00:00")
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	turbines := p.Turbine
	defaultValue := -999999.00

	pipes = append(pipes, tk.M{"$group": tk.M{"_id": tk.M{"colId": "$timestamp", "Turbine": "$turbine"},
		"avgwindspeed": tk.M{"$avg": "$windspeed"},
		"sumwindspeed": tk.M{"$sum": "$windspeed"},
		"activepower":  tk.M{"$sum": "$activepower"},
		"rotorrpm":     tk.M{"$sum": "$rotorrpm"},
		"totaldata":    tk.M{"$sum": 1}}})
	pipes = append(pipes, tk.M{"$sort": tk.M{"_id": 1}})

	filter = nil
	filter = append(filter, dbox.Ne("_id", ""))
	filter = append(filter, dbox.Eq("dateinfo.dateid", lastDateData))
	filter = append(filter, dbox.Eq("turbine", turbines))
	filter = append(filter, dbox.Ne("activepower", defaultValue))
	filter = append(filter, dbox.Ne("windspeed", defaultValue))

	csr, e := DB().Connection.NewQuery().
		From(new(ScadaRealTime).TableName()).
		Command("pipe", pipes).
		Where(dbox.And(filter...)).
		Cursor(nil)

	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}
	e = csr.Fetch(&list, 0, false)
	defer csr.Close()

	totactivepower := 0.0
	totwindspeed := 0.0
	totrotorrpm := 0.0
	totData := 0.0
	dataMonitoring := tk.M{}
	for _, val := range list {

		seriesData := tk.M{}
		avgwindspeed := val.GetFloat64("avgwindspeed")
		sumwindspeed := val.GetFloat64("sumwindspeed")
		activepower := val.GetFloat64("activepower")
		rotorrpm := val.GetFloat64("rotorrpm")
		totaldata := val.GetFloat64("totaldata")
		idD := val.Get("_id").(tk.M)
		Turbine := idD.Get("Turbine")
		timestamp := idD.Get("colId").(time.Time).UTC().Format("2006-01-02 15:04:05")

		seriesData.Set("turbine", Turbine)
		seriesData.Set("timestamp", timestamp)
		seriesData.Set("activepower", tk.Div(activepower, 1000.0))
		seriesData.Set("avgwindspeed", avgwindspeed)

		dataSeries = append(dataSeries, seriesData)

		totactivepower = totactivepower + activepower
		totwindspeed = totwindspeed + sumwindspeed
		totrotorrpm = totrotorrpm + rotorrpm
		totData = totData + totaldata

	}

	dataMonitoring.Set("Power", tk.Div(totactivepower, 1000.0))
	dataMonitoring.Set("WindSpeed", tk.Div(totwindspeed, totData))
	dataMonitoring.Set("RotorRpm", totrotorrpm)

	data := struct {
		Data       []tk.M
		Monitoring tk.M
	}{
		Data:       dataSeries,
		Monitoring: dataMonitoring,
	}

	return helper.CreateResult(true, data, "success")
}

func getTurbineStatus(project string) (res map[string]TurbineStatus) {
	res = map[string]TurbineStatus{}

	csr, err := DB().Connection.NewQuery().From(new(TurbineStatus).TableName()).
		Cursor(nil)

	if err != nil {
		return
	}

	results := make([]TurbineStatus, 0)
	err = csr.Fetch(&results, 0, false)
	if err != nil {
		return
	}
	csr.Close()

	for _, result := range results {
		res[result.ID] = result
	}

	return
}

/*
func (c *MonitoringRealtimeController) GetMonitoring() tk.M {
	turbines := []string{
		"SSE017", "SSE018", "SSE019", "SSE020", "TJ013", "TJ016", "HBR038", "TJ021", "TJ022", "TJ023", "TJ024",
		"TJ025", "HBR004", "HBR005", "HBR006", "TJW024", "HBR007", "SSE001", "SSE002", "SSE007", "SSE006", "SSE011",
		"SSE015", "SSE012",
	}
	defaultValue := -999999.00
	defaultProject := "Tejuva"

	mdl := new(ScadaMonitoring).New()

	mdl.TimeStamp = time.Now()
	mdl.DateInfo = GetDateInfo(mdl.TimeStamp)
	mdl.ActivePower = defaultValue
	mdl.Production = defaultValue
	mdl.OprHours = 0.0
	mdl.WtgOkHours = 0.0
	mdl.WindSpeed = defaultValue
	mdl.WindDirection = defaultValue
	mdl.NacellePosition = defaultValue
	mdl.Temperature = defaultValue
	mdl.PitchAngle = defaultValue
	mdl.RotorRPM = defaultValue
	mdl.ProjectName = defaultProject

	mdl.WindSpeedCount = 0
	mdl.WindDirectionCount = 0
	mdl.NacellePositionCount = 0
	mdl.TemperatureCount = 0
	mdl.PitchAngleCount = 0
	mdl.RotorRPMCount = 0

	power := 0.0
	windSpeed := 0.0
	cWindSpeed := 0
	windDir := 0.0
	cWindDir := 0
	nacellePos := 0.0
	cNacellePos := 0
	temperature := 0.0
	cTemperature := 0
	pitch := 0.0
	cPitch := 0
	rotor := 0.0
	cRotor := 0

	timeUpdate := time.Now().UTC().Add(-720 * time.Hour)
	details := make([]ScadaMonitoringItem, 0)
	for _, t := range turbines {
		var detail ScadaMonitoringItem
		detail.Turbine = t
		detail.DataComing = 0

		csrt, err := DB().Connection.NewQuery().From(new(ScadaRealTime).TableName()).
			Where(dbox.And(dbox.Eq("turbine", t), dbox.Ne("activepower", defaultValue))).
			Order("-timestamp").Cursor(nil)
		if err != nil {
			tk.Println(err.Error())
		}
		results := make([]ScadaRealTime, 0)
		err = csrt.Fetch(&results, 1, false)
		if err != nil {
			tk.Println(err.Error())
		}
		csrt.Close()

		if len(results) > 0 {
			result := results[0]
			detail.ActivePower = result.ActivePower
			detail.TimeUpdate = result.LastUpdate
			timeNow := time.Now() //.UTC().Add(5.5 * time.Hour)
			if timeNow.Sub(result.LastUpdate).Minutes() <= 3 {
				detail.DataComing = 1
			}
		} else {
			detail.ActivePower = defaultValue
		}

		if detail.ActivePower > defaultValue {
			power += detail.ActivePower
		}

		csrt, err = DB().Connection.NewQuery().From(new(ScadaRealTime).TableName()).
			Where(dbox.And(dbox.Eq("turbine", t), dbox.Ne("windspeed", defaultValue))).
			Order("-timestamp").Cursor(nil)
		if err != nil {
			tk.Println(err.Error())
		}
		results = make([]ScadaRealTime, 0)
		err = csrt.Fetch(&results, 1, false)
		if err != nil {
			tk.Println(err.Error())
		}
		csrt.Close()

		if len(results) > 0 {
			result := results[0]
			detail.WindSpeed = result.WindSpeed
			if result.LastUpdate.Sub(detail.TimeUpdate).Seconds() > 0 {
				detail.TimeUpdate = result.LastUpdate
			}
			timeNow := time.Now() //.UTC().Add(5.5 * time.Hour)
			if timeNow.Sub(result.LastUpdate).Minutes() <= 3 {
				detail.DataComing = 1
			}
		} else {
			detail.WindSpeed = defaultValue
		}

		if detail.WindSpeed != defaultValue {
			windSpeed += detail.WindSpeed
			cWindSpeed++
		}

		csrt, err = DB().Connection.NewQuery().From(new(ScadaRealTime).TableName()).
			Where(dbox.And(dbox.Eq("turbine", t), dbox.Ne("winddirection", defaultValue))).
			Order("-timestamp").Cursor(nil)
		if err != nil {
			tk.Println(err.Error())
		}
		results = make([]ScadaRealTime, 0)
		err = csrt.Fetch(&results, 1, false)
		if err != nil {
			tk.Println(err.Error())
		}
		csrt.Close()

		if len(results) > 0 {
			result := results[0]
			detail.WindDirection = result.WindDirection
			if result.LastUpdate.Sub(detail.TimeUpdate).Seconds() > 0 {
				detail.TimeUpdate = result.LastUpdate
			}
			timeNow := time.Now() //.UTC().Add(5.5 * time.Hour)
			if timeNow.Sub(result.LastUpdate).Minutes() <= 3 {
				detail.DataComing = 1
			}
		} else {
			detail.WindDirection = defaultValue
		}

		if detail.WindDirection != defaultValue {
			windDir += detail.WindDirection
			cWindDir++
		}

		csrt, err = DB().Connection.NewQuery().From(new(ScadaRealTime).TableName()).
			Where(dbox.And(dbox.Eq("turbine", t), dbox.Ne("nacelleposition", defaultValue))).
			Order("-timestamp").Cursor(nil)
		if err != nil {
			tk.Println(err.Error())
		}
		results = make([]ScadaRealTime, 0)
		err = csrt.Fetch(&results, 1, false)
		if err != nil {
			tk.Println(err.Error())
		}
		csrt.Close()

		if len(results) > 0 {
			result := results[0]
			detail.NacellePosition = result.NacellePosition
			if result.LastUpdate.Sub(detail.TimeUpdate).Seconds() > 0 {
				detail.TimeUpdate = result.LastUpdate
			}
			timeNow := time.Now() //.UTC().Add(5.5 * time.Hour)
			if timeNow.Sub(result.LastUpdate).Minutes() <= 3 {
				detail.DataComing = 1
			}
		} else {
			detail.NacellePosition = defaultValue
		}

		if detail.NacellePosition != defaultValue {
			nacellePos += detail.NacellePosition
			cNacellePos++
		}

		csrt, err = DB().Connection.NewQuery().From(new(ScadaRealTime).TableName()).
			Where(dbox.And(dbox.Eq("turbine", t), dbox.Ne("temperature", defaultValue))).
			Order("-timestamp").Cursor(nil)
		if err != nil {
			tk.Println(err.Error())
		}
		results = make([]ScadaRealTime, 0)
		err = csrt.Fetch(&results, 1, false)
		if err != nil {
			tk.Println(err.Error())
		}
		csrt.Close()

		if len(results) > 0 {
			result := results[0]
			detail.Temperature = result.Temperature
			if result.LastUpdate.Sub(detail.TimeUpdate).Seconds() > 0 {
				detail.TimeUpdate = result.LastUpdate
			}
			timeNow := time.Now() //.UTC().Add(5.5 * time.Hour)
			if timeNow.Sub(result.LastUpdate).Minutes() <= 3 {
				detail.DataComing = 1
			}
		} else {
			detail.Temperature = defaultValue
		}

		if detail.Temperature != defaultValue {
			temperature += detail.Temperature
			cTemperature++
		}

		csrt, err = DB().Connection.NewQuery().From(new(ScadaRealTime).TableName()).
			Where(dbox.And(dbox.Eq("turbine", t), dbox.Ne("pitchangle", defaultValue))).
			Order("-timestamp").Cursor(nil)
		if err != nil {
			tk.Println(err.Error())
		}
		results = make([]ScadaRealTime, 0)
		err = csrt.Fetch(&results, 1, false)
		if err != nil {
			tk.Println(err.Error())
		}
		csrt.Close()

		if len(results) > 0 {
			result := results[0]
			detail.PitchAngle = result.PitchAngle
			if result.LastUpdate.Sub(detail.TimeUpdate).Seconds() > 0 {
				detail.TimeUpdate = result.LastUpdate
			}
			timeNow := time.Now() //.UTC().Add(5.5 * time.Hour)
			if timeNow.Sub(result.LastUpdate).Minutes() <= 3 {
				detail.DataComing = 1
			}
		} else {
			detail.PitchAngle = defaultValue
		}

		if detail.PitchAngle != defaultValue {
			pitch += detail.PitchAngle
			cPitch++
		}

		csrt, err = DB().Connection.NewQuery().From(new(ScadaRealTime).TableName()).
			Where(dbox.And(dbox.Eq("turbine", t), dbox.Ne("rotorrpm", defaultValue))).
			Order("-timestamp").Cursor(nil)
		if err != nil {
			tk.Println(err.Error())
		}
		results = make([]ScadaRealTime, 0)
		err = csrt.Fetch(&results, 1, false)
		if err != nil {
			tk.Println(err.Error())
		}
		csrt.Close()

		if len(results) > 0 {
			result := results[0]
			detail.RotorRPM = result.RotorRPM
			if result.LastUpdate.Sub(detail.TimeUpdate).Seconds() > 0 {
				detail.TimeUpdate = result.LastUpdate
			}
			timeNow := time.Now() //.UTC().Add(5.5 * time.Hour)
			if timeNow.Sub(result.LastUpdate).Minutes() <= 3 {
				detail.DataComing = 1
			}
		} else {
			detail.RotorRPM = defaultValue
		}

		if detail.RotorRPM != defaultValue {
			rotor += detail.RotorRPM
			cRotor++
		}

		details = append(details, detail)
		if detail.TimeUpdate.Sub(timeUpdate) >= 0 {
			timeUpdate = detail.TimeUpdate
		}
	}

	mdl.TimeStamp = timeUpdate
	mdl.Detail = details

	// getting turbine status
	csra, err := DB().Connection.NewQuery().From(new(TurbineStatus).TableName()).
		Cursor(nil)
	if err != nil {
		tk.Println(err.Error())
	}
	rests := make([]TurbineStatus, 0)
	err = csra.Fetch(&rests, 0, false)
	if err != nil {
		tk.Println(err.Error())
	}
	csra.Close()

	ret := tk.M{}.
		Set("Data", mdl).
		Set("TurbineStatus", rests)

	return ret
}
*/
