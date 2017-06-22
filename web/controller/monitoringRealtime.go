package controller

import (
	. "eaciit/wfdemo-git/library/core"
	lh "eaciit/wfdemo-git/library/helper"
	. "eaciit/wfdemo-git/library/models"
	"log"

	"eaciit/wfdemo-git/web/helper"

	"github.com/eaciit/dbox"
	_ "github.com/eaciit/dbox/dbc/mongo"

	"github.com/eaciit/knot/knot.v1"
	tk "github.com/eaciit/toolkit"

	"bufio"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"

	c "github.com/eaciit/crowd"
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
	arrlabel     = map[string]string{"Wind speed Avg": "WindSpeed_ms", "Wind speed 1": "", "Wind speed 2": "",
		"Wind Direction": "WindDirection", "Vane 1 wind direction": "",
		"Vane 2 wind direction": "", "Nacelle Direction": "NacellePos",
		"Rotor RPM": "RotorSpeed_RPM", "Generator RPM": "GenSpeed_RPM",
		"DFIG speed generator encoder": "", "Blade Angle 1": "PitchAngle1",
		"Blade Angle 2": "PitchAngle2", "Blade Angle 3": "PitchAngle3",
		"Volt. Battery - blade 1": "PitchAccuV1", "Volt. Battery - blade 2": "PitchAccuV2",
		"Volt. Battery - blade 3": "PitchAccuV3", "Current 1 Pitch Motor": "PitchConvCurrent1",
		"Current 2 Pitch Motor": "PitchConvCurrent2", "Current 3 Pitch Motor": "PitchConvCurrent3",
		"Pitch motor temperature - Blade 1": "TempConv1", "Pitch motor temperature - Blade 2": "TempConv2",
		"Pitch motor temperature - Blade 3": "TempConv3", "Phase 1 voltage": "VoltageL1",
		"Phase 2 voltage": "VoltageL2", "Phase 3 voltage": "VoltageL3", "Phase 1 current": "CurrentL1",
		"Phase 2 current": "CurrentL2", "Phase 3 current": "CurrentL3", "Power": "ActivePower_kW",
		"Power Reactive": "ReactivePower_kVAr", "Freq. Grid": "Frequency_Hz", "Production": "Total_Prod_Day_kWh",
		"Cos Phi": "PowerFactor", "DFIG active power": "", "DFIG reactive power": "", "DFIG mains Frequency": "",
		"DFIG main voltage": "", "DFIG main current": "", "DFIG DC link voltage": "",
		"Rotor R current": "", "Roter Y current": "", "Roter B current": "",
		"Temp. generator 1 phase 1 coil": "TempG1L1", "Temp. generator 1 phase 2 coil": "TempG1L2", "Temp. generator 1 phase 3 coil": "TempG1L3",
		"Temp. generator bearing driven End": "TempGeneratorBearingDE", "Temp. generator bearing non-driven End": "TempGeneratorBearingNDE",
		"Temp. Gearbox driven end": "TempGearBoxHSSDE", "Temp. Gearbox non-driven end": "TempGearBoxHSSNDE", "Temp. Gearbox inter. driven end": "TempGearBoxIMSDE",
		"Temp. Gearbox inter. non-driven end": "TempGearBoxIMSNDE", "Pressure Gear box oil": "",
		"Temp. Gear box oil": "TempGearBoxOilSump", "Temp. Nacelle": "TempNacelle", "Temp. Ambient": "TempOutdoor",
		"Temp. Main bearing": "TempHubBearing", "Damper Oscillation mag.": "", "Drive train vibration": "DrTrVibValue",
		"Tower vibration": "",
	}
)

type MiniScadaHFD struct {
	Slow_Nacellepos    float64
	Fast_Windspeed_Ms  float64
	Slow_Winddirection float64
	Turbine            string
}

func (m *MonitoringRealtimeController) GetWindRoseMonitoring(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	type PayloadWSMonitoring struct {
		Project   string
		Turbine   string
		BreakDown int
	}

	p := new(PayloadWSMonitoring)
	e := k.GetPayload(&p)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	var tStart, tEnd time.Time
	now := time.Now().UTC()
	// // now := time.Date(2017, 3, 8, 9, 20, 0, 0, time.UTC)
	// last := now.AddDate(0, 0, -24)

	indiaLoc, _ := time.LoadLocation("Asia/Kolkata")
	indiaTime := now.In(indiaLoc)
	indiaNow := time.Date(indiaTime.Year(), indiaTime.Month(), indiaTime.Day(), indiaTime.Hour(), indiaTime.Minute(), indiaTime.Second(), indiaTime.Nanosecond(), time.UTC)

	last := indiaNow.Add(time.Duration(-24) * time.Hour)

	// tStart, _ = time.Parse("20060102", last.Format("200601")+"01")
	// tEnd, _ = time.Parse("20060102", indiaNow.Format("200601")+"01")

	tStart = last
	tEnd = indiaNow

	// log.Printf(">> %v | %v \n", last.String(), indiaNow.String())
	// log.Printf(">> %v | %v \n", tStart.String(), tEnd.String())

	section = p.BreakDown
	getFullWSCategory()

	query := []tk.M{}
	pipes := []tk.M{}
	query = append(query, tk.M{"_id": tk.M{"$ne": nil}})
	query = append(query, tk.M{"timestamp": tk.M{"$gte": tStart}})
	query = append(query, tk.M{"timestamp": tk.M{"$lt": tEnd}})
	query = append(query, tk.M{"fast_windspeed_ms": tk.M{"$gt": -999999.0}})
	query = append(query, tk.M{"slow_winddirection": tk.M{"$gt": -999999.0}})
	query = append(query, tk.M{"slow_nacellepos": tk.M{"$gt": -999999.0}})

	if p.Project != "" {
		query = append(query, tk.M{"projectname": p.Project})
	}

	data := []MiniScadaHFD{}
	_data := MiniScadaHFD{}

	turbineVal := p.Turbine

	pipes = []tk.M{}
	data = []MiniScadaHFD{}
	queryT := query
	queryT = append(queryT, tk.M{"turbine": turbineVal})
	pipes = append(pipes, tk.M{"$match": tk.M{"$and": queryT}})
	pipes = append(pipes, tk.M{"$project": tk.M{"slow_nacellepos": 1, "fast_windspeed_ms": 1, "slow_winddirection": 1}})
	csr, _ := DB().Connection.NewQuery().From(new(ScadaDataHFD).TableName()).
		Command("pipe", pipes).Cursor(nil)

	for {
		e = csr.Fetch(&_data, 1, false)
		if e != nil {
			break
		}
		data = append(data, _data)
	}
	csr.Close()
	// dataNacelle := GenerateWindRose(data, "nacelle", turbineVal)
	// dataWindDir := GenerateWindRose(data, "winddir", turbineVal)
	datas := GenerateWindRose(data, "NacellePlusWind", turbineVal)

	return helper.CreateResult(true, datas, "success")

}

func GenerateWindRose(data []MiniScadaHFD, tipe, turbineVal string) tk.M {
	WsMonitoringRes := []tk.M{}
	maxValue := 0.0
	tkMaxVal := tk.M{}
	groupdata := tk.M{}
	if tk.SliceLen(data) > 0 {
		totalDuration := float64(len(data)) /* Tot data * 2 for get total minutes*/
		datas := c.From(&data).Apply(func(x interface{}) interface{} {
			dt := x.(MiniScadaHFD)
			var di DataItems
			var dirNo, dirDesc int

			if tipe == "nacelle" {
				dirNo, dirDesc = getDirection(dt.Slow_Nacellepos, section)
			} else if tipe == "winddir" {
				dirNo, dirDesc = getDirection(dt.Slow_Winddirection+300, section)
			} else {
				dirNo, dirDesc = getDirection(dt.Slow_Nacellepos+dt.Slow_Winddirection+300, section)
			}

			wsNo, wsDesc := getWsCategory(dt.Fast_Windspeed_Ms)

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
			kv := x.(c.KV)
			vv := kv.Key.(DataItemsGroup)
			vs := kv.Value.([]DataItems)

			sumFreq := c.From(&vs).Sum(func(x interface{}) interface{} {
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

			key := turbineVal + "_" + tk.ToString(di.DirectionNo)

			if !tkMaxVal.Has(key) {
				tkMaxVal.Set(key, di.Contribution)
			} else {
				tkMaxVal.Set(key, tkMaxVal.GetFloat64(key)+di.Contribution)
			}

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

		WsMonitoringRes = append(WsMonitoringRes, groupdata)
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
		WsMonitoringRes = append(WsMonitoringRes, groupdata)
	}

	for _, val := range tkMaxVal {
		if val.(float64) > maxValue {
			maxValue = val.(float64)
		}
	}

	switch {
	case maxValue >= 90 && maxValue <= 100:
		maxValue = 100
	case maxValue >= 80 && maxValue < 90:
		maxValue = 90
	case maxValue >= 70 && maxValue < 80:
		maxValue = 80
	case maxValue >= 60 && maxValue < 70:
		maxValue = 70
	case maxValue >= 50 && maxValue < 60:
		maxValue = 60
	case maxValue >= 40 && maxValue < 50:
		maxValue = 50
	case maxValue >= 30 && maxValue < 40:
		maxValue = 40
	case maxValue >= 20 && maxValue < 30:
		maxValue = 30
	case maxValue >= 10 && maxValue < 20:
		maxValue = 20
	case maxValue >= 0 && maxValue < 10:
		maxValue = 10
	}

	result := tk.M{
		"WindRose": WsMonitoringRes,
		"MaxValue": maxValue,
	}

	return result
}

func (c *MonitoringRealtimeController) GetDataProject(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson
	k.Config.NoLog = true

	p := struct {
		Project string
	}{}

	err := k.GetPayload(&p)
	if err != nil {
		return helper.CreateResult(false, nil, err.Error())
	}

	results := GetMonitoringByProjectV2(p.Project, "monitoring")

	return helper.CreateResult(true, results, "success")
}

func (c *MonitoringRealtimeController) getValue() float64 {
	retVal := 0.0

	return retVal
}

func GetMonitoringByProjectV2(project string, pageType string) (rtkm tk.M) {
	rtkm = tk.M{}
	alldata, allturbine := []tk.M{}, tk.M{}
	turbineMap := map[string]tk.M{}
	totalCapacity := 0.0

	csrt, err := DB().Connection.NewQuery().Select("turbineid", "feeder", "turbinename", "latitude", "longitude", "capacitymw").
		From("ref_turbine").
		Where(dbox.Eq("project", project)).
		Order("turbinename").
		Cursor(nil)

	if err != nil {
		tk.Println(err.Error())
	}

	_result := []tk.M{}
	err = csrt.Fetch(&_result, 0, false)
	if err != nil {
		tk.Println(err.Error())
	}
	csrt.Close()
	for _, _tkm := range _result {
		turbine := _tkm.GetString("turbineid")
		lturbine := allturbine.Get(_tkm.GetString("feeder"), []string{}).([]string)
		lturbine = append(lturbine, turbine)
		sort.Strings(lturbine)
		allturbine.Set(_tkm.GetString("feeder"), lturbine)
		turbineMap[turbine] = tk.M{"coords": []float64{_tkm.GetFloat64("latitude"), _tkm.GetFloat64("longitude")}, "name": _tkm.GetString("turbinename"), "capacity": _tkm.GetFloat64("capacitymw") * 1000.0}
		totalCapacity += _tkm.GetFloat64("capacitymw")
	}

	// arrfield := []string{"ActivePower", "WindSpeed", "WindDirection", "NacellePosition", "Temperature",
	// 	"PitchAngle", "RotorRPM"}
	arrfield := map[string]string{"ActivePower_kW": "ActivePower", "WindSpeed_ms": "WindSpeed",
		"WindDirection": "WindDirection", "NacellePos": "NacellePosition", "TempOutdoor": "Temperature",
		"PitchAngle": "PitchAngle", "RotorSpeed_RPM": "RotorRPM"}

	lastUpdate := time.Time{}
	PowerGen, AvgWindSpeed, CountWS := float64(0), float64(0), float64(0)
	turbinedown, turbnotavail := 0, 0
	t0 := getTimeNow()

	arrturbinestatus := GetTurbineStatus(project, "")
	timemax := getMaxRealTime(project, "")
	timecond := time.Date(timemax.Year(), timemax.Month(), timemax.Day(), 0, 0, 0, 0, timemax.Location())

	// rconn := lh.GetConnRealtime()
	// defer rconn.Close()
	rconn := DBRealtime()

	csr, err := rconn.NewQuery().From(new(ScadaRealTimeNew).TableName()).
		Where(dbox.And(dbox.Gte("timestamp", timecond), dbox.Eq("projectname", project))).
		Order("turbine", "-timestamp").Cursor(nil)
	if err != nil {
		tk.Println(err.Error())
	}

	_iTurbine, _iContinue, _itkm := "", false, tk.M{}
	for {
		_tdata := tk.M{}
		err = csr.Fetch(&_tdata, 1, false)
		if err != nil {
			break
		}

		tags := _tdata.GetString("tags")

		// isfound := false
		// for k, _ := range arrfield {
		// 	if arrfield.GetString(k) == tags {
		// 		isfound = true
		// 	}
		// }

		// if !isfound {
		// 	continue
		// }

		_tTurbine := _tdata.GetString("turbine")
		if _iContinue && _iTurbine == _tTurbine {
			continue
		}
		/*tstamp := _tdata.Get("timestamp", time.Time{}).(time.Time)

		if tstamp.After(lastUpdate) {
			lastUpdate = tstamp.UTC()
		}*/

		tstamp := _tdata.Get("timestamp", time.Time{}).(time.Time)

		if tstamp.After(lastUpdate) {
			lastUpdate = tstamp
		}

		if _iTurbine != _tTurbine {
			if _iTurbine != "" {
				alldata = append(alldata, _itkm)
			}
			_iContinue = false
			_iTurbine = _tTurbine
			turbineMp := turbineMap[_tTurbine]

			if pageType == "monitoring" {
				_itkm = tk.M{}.
					Set("Turbine", _tTurbine).
					Set("Name", turbineMp.GetString("name")).
					Set("DataComing", 0).
					Set("AlarmCode", 0).
					Set("AlarmDesc", "").
					Set("Status", 1).
					Set("IsWarning", false).
					Set("AlarmUpdate", time.Time{}).
					Set("Capacity", turbineMp.GetFloat64("capacity"))

				for _, afield := range arrfield {
					_itkm.Set(afield, defaultValue)
				}

				turbnotavail++
				if t0.Sub(tstamp.UTC()).Minutes() <= 3 {
					_itkm.Set("DataComing", 1)
					turbnotavail--
				}

				if _idt, _cond := arrturbinestatus[_tTurbine]; _cond {
					_itkm.Set("AlarmCode", _idt.AlarmCode).
						Set("AlarmDesc", _idt.AlarmDesc).
						Set("Status", _idt.Status).
						Set("IsWarning", _idt.IsWarning).
						Set("AlarmUpdate", _idt.TimeUpdate.UTC())
					if _idt.Status == 0 && _itkm.GetInt("DataComing") == 1 {
						turbinedown += 1
					}
				}
			} else if pageType == "dashboard" {
				_itkm = tk.M{}.
					Set("DataComing", 0).
					Set("Status", 1).
					Set("IsWarning", false)

				if t0.Sub(tstamp.UTC()).Minutes() <= 3 {
					_itkm.Set("DataComing", 1)
				}

				if _idt, _cond := arrturbinestatus[_tTurbine]; _cond {
					_itkm.
						Set("Status", _idt.Status).
						Set("IsWarning", _idt.IsWarning)
					if _idt.Status == 0 {
						turbinedown += 1
					}
				}

				_itkm.
					Set("coords", turbineMp.Get("coords")).
					Set("name", turbineMp.GetString("name")).
					Set("value", _tTurbine)
			}

		}

		// _iContinue = true

		afield, isexist := arrfield[tags]
		_ifloat := _tdata.GetFloat64("value")
		if _ifloat != defaultValue && isexist {
			switch afield {
			case "ActivePower":
				PowerGen += _ifloat
			case "WindSpeed":
				AvgWindSpeed += _ifloat
				CountWS += 1
			}

			_itkm.Set(afield, _ifloat)
		}

		// for _, afield := range arrfield {
		// 	_lafield := strings.ToLower(afield)
		// 	if _ifloat := _tdata.GetFloat64(_lafield); _ifloat != defaultValue && (_itkm.GetFloat64(afield) == 0 || _itkm.GetFloat64(afield) == defaultValue) {
		// 		_itkm.Set(afield, _ifloat)

		// 		switch afield {
		// 		case "ActivePower":
		// 			PowerGen += _ifloat
		// 		case "WindSpeed":
		// 			AvgWindSpeed += _ifloat
		// 			CountWS += 1
		// 		}

		// 	} else {
		// 		_iContinue = false
		// 	}
		// }
		_itkm.Set("IconStatus", "fa fa-circle fa-project-info fa-green")
		if _itkm.GetInt("Status") == 0 {
			_itkm.Set("IconStatus", "fa fa-circle fa-project-info fa-red")
		} else if _itkm.GetInt("Status") == 1 && _itkm["IsWarning"].(bool) {
			_itkm.Set("IconStatus", "fa fa-circle fa-project-info fa-orange")
		}
		if _itkm.GetInt("DataComing") == 0 {
			_itkm.Set("IconStatus", "fa fa-circle fa-project-info fa-grey")
		}
		if _itkm.GetFloat64("ActivePower") < 0 {
			_itkm.Set("ActivePowerColor", "redvalue")
		} else {
			_itkm.Set("ActivePowerColor", "defaultcolor")
		}
		if _itkm.GetFloat64("WindSpeed") < 3.5 {
			_itkm.Set("WindSpeedColor", "orangevalue")
		} else {
			_itkm.Set("WindSpeedColor", "defaultcolor")
		}
		if _itkm.GetFloat64("Temperature") > 38 {
			_itkm.Set("TemperatureColor", "txt-red")
		} else if _itkm.GetFloat64("Temperature") >= 30 {
			_itkm.Set("TemperatureColor", "txt-orange")
		} else {
			_itkm.Set("TemperatureColor", "txt-grey")
		}
	}
	csr.Close()
	if _iTurbine != "" {
		alldata = append(alldata, _itkm)
	}

	isInDetail := func(_turbine string) bool {
		for _, _tkm := range alldata {
			if _turbine == _tkm.GetString("Turbine") {
				return true
			}
		}
		return false
	}

	for _, _tkm := range _result {
		_turbine := _tkm.GetString("turbineid")
		if isInDetail(_turbine) {
			continue
		}

		turbineMp := turbineMap[_turbine]
		turbnotavail++

		_itkm = tk.M{}.
			Set("Turbine", _turbine).
			Set("Name", turbineMp.GetString("name")).
			Set("DataComing", 0).
			Set("AlarmCode", 0).
			Set("AlarmDesc", "").
			Set("Status", 0).
			Set("IsWarning", false).
			Set("AlarmUpdate", time.Time{}).
			Set("DataComing", 0).
			Set("IconStatus", "fa fa-circle fa-project-info fa-grey").
			Set("ActivePowerColor", "defaultcolor").
			Set("TemperatureColor", "defaultcolor").
			Set("WindSpeedColor", "defaultcolor").
			Set("Capacity", turbineMp.GetFloat64("capacity"))

		for _, afield := range arrfield {
			_itkm.Set(afield, defaultValue)
		}

		alldata = append(alldata, _itkm)
	}

	indiaLoc, _ := time.LoadLocation("Asia/Kolkata")
	indiaTime := lastUpdate.In(indiaLoc)
	lastUpdateIndia := time.Date(indiaTime.Year(), indiaTime.Month(), indiaTime.Day(), indiaTime.Hour(), indiaTime.Minute(), indiaTime.Second(), indiaTime.Nanosecond(), time.UTC)

	if pageType == "monitoring" {
		turbineactive := len(_result) - turbinedown - turbnotavail
		rtkm.Set("ListOfTurbine", allturbine)
		rtkm.Set("Detail", alldata)
		rtkm.Set("TimeNow", t0)
		rtkm.Set("TimeStamp", lastUpdateIndia)
		rtkm.Set("TimeMax", timemax)
		rtkm.Set("PowerGeneration", PowerGen)
		rtkm.Set("AvgWindSpeed", tk.Div(AvgWindSpeed, CountWS))
		rtkm.Set("PLF", tk.Div(PowerGen, (totalCapacity*1000))*100)
		rtkm.Set("TurbineActive", turbineactive)
		rtkm.Set("TurbineDown", turbinedown)
		rtkm.Set("TurbineNotAvail", turbnotavail)
	} else if pageType == "dashboard" {
		rtkm.Set("Detail", alldata)
		rtkm.Set("TurbineDown", turbinedown)
		rtkm.Set("TurbineActive", len(_result)-turbinedown)
	}

	return
}

func (c *MonitoringRealtimeController) GetDataAlarm(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson
	k.Config.NoLog = true

	type MyPayloads struct {
		Turbine   []interface{}
		DateStart time.Time
		DateEnd   time.Time
		Skip      int
		Take      int
		Sort      []Sorting
		Project   string
		Period    string
		Tipe      string
	}

	type Sorting struct {
		Field string
		Dir   string
	}

	p := new(MyPayloads)
	err := k.GetPayload(&p)
	if err != nil {
		return helper.CreateResult(false, nil, err.Error())
	}

	tStart, tEnd, e := helper.GetStartEndDate(k, p.Period, p.DateStart, p.DateEnd)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	project := p.Project
	tablename := "AlarmHFD"
	if p.Tipe == "warning" {
		tablename = "AlarmWarning"
	}

	dfilter := []*dbox.Filter{dbox.Eq("isdeleted", false)}
	dfilter = append(dfilter, dbox.Eq("projectname", project))
	dfilter = append(dfilter, dbox.Gte("timestart", tStart), dbox.Lte("timestart", tEnd))
	if len(p.Turbine) > 0 {
		dfilter = append(dfilter, dbox.In("turbine", p.Turbine...))
	}

	// rconn := lh.GetConnRealtime()
	// defer rconn.Close()
	rconn := DBRealtime()

	csr, err := rconn.NewQuery().From(tablename).
		Aggr(dbox.AggrSum, "$duration", "duration").
		Aggr(dbox.AggrSum, 1, "countdata").
		Group("projectname").
		Where(dbox.And(dfilter...)).Cursor(nil)

	if err != nil {
		return helper.CreateResult(false, nil, err.Error())
	}
	defer csr.Close()

	tkmgroup := tk.M{}
	_ = csr.Fetch(&tkmgroup, 1, false)

	totalData := tkmgroup.GetInt("countdata")
	totalDuration := tkmgroup.GetInt("duration")
	query := rconn.NewQuery().From(tablename).
		Where(dbox.And(dfilter...)).
		// Order("-timestart").
		Skip(p.Skip).Take(p.Take)

	if len(p.Sort) > 0 {
		var arrsort []string
		for _, val := range p.Sort {
			if val.Dir == "desc" {
				arrsort = append(arrsort, strings.ToLower("-"+strings.ToLower(val.Field)))
			} else {
				arrsort = append(arrsort, strings.ToLower(strings.ToLower(val.Field)))
			}
		}
		query = query.Order(arrsort...)
	}
	csr, err = query.Cursor(nil)
	if err != nil {
		return helper.CreateResult(false, nil, err.Error())
	}
	defer csr.Close()

	results := make([]AlarmHFD, 0)
	err = csr.Fetch(&results, 0, false)
	if err != nil {
		return helper.CreateResult(false, nil, err.Error())
	}

	// tStart, tEnd = tStart.UTC(), tEnd.UTC()

	// rStart := time.Date(tStart.Y, month, day, hour, min, sec, nsec, loc)

	csrTurbine, err := rconn.NewQuery().From("ref_turbine").
		Where(dbox.Eq("project", project)).Cursor(nil)
	if err != nil {
		return helper.CreateResult(false, nil, err.Error())
	}
	defer csrTurbine.Close()
	turbineList := []tk.M{}
	err = csrTurbine.Fetch(&turbineList, 0, false)
	if err != nil {
		return helper.CreateResult(false, nil, err.Error())
	}
	turbineName := map[string]string{}
	for _, val := range turbineList {
		turbineName[val.GetString("turbineid")] = val.GetString("turbinename")
	}
	for idx, val := range results {
		results[idx].Turbine = turbineName[val.Turbine]
	}

	retData := tk.M{}.Set("Data", results).
		Set("Total", totalData).
		Set("Duration", totalDuration).
		Set("mindate", tStart.UTC()).
		Set("maxdate", tEnd.UTC())

	return helper.CreateResult(true, retData, "success")
}

func (c *MonitoringRealtimeController) GetDataAlarmAvailDate(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson
	k.Config.NoLog = true

	type MyPayloads struct {
		Tipe    string
		Project string
	}

	p := new(MyPayloads)
	err := k.GetPayload(&p)
	if err != nil {
		return helper.CreateResult(false, nil, err.Error())
	}

	project := p.Project

	dfilter := []*dbox.Filter{}
	dfilter = append(dfilter, dbox.Eq("projectname", project))
	dfilter = append(dfilter, dbox.Ne("timestart", time.Time{}))

	// rconn := lh.GetConnRealtime()
	// defer rconn.Close()
	rconn := DBRealtime()
	tablename := "AlarmHFD"
	if p.Tipe == "warning" {
		tablename = "AlarmWarning"
	}
	csr, err := rconn.NewQuery().From(tablename).
		Aggr(dbox.AggrMin, "$timestart", "minstart").
		Aggr(dbox.AggrMax, "$timestart", "maxstart").
		Group("projectname").
		Where(dbox.And(dfilter...)).Cursor(nil)

	if err != nil {
		return helper.CreateResult(false, nil, err.Error())
	}

	tkmgroup := tk.M{}
	_ = csr.Fetch(&tkmgroup, 1, false)
	csr.Close()

	minDate := tkmgroup.Get("minstart", time.Time{}).(time.Time).UTC()
	maxDate := tkmgroup.Get("maxstart", time.Time{}).(time.Time).UTC()

	return helper.CreateResult(true, tk.M{}.Set("Data", []time.Time{minDate, maxDate}), "success")
}

func (c *MonitoringRealtimeController) GetDataTurbine(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson
	k.Config.NoLog = true

	p := struct {
		Project string
		Turbine string
	}{}

	alldata := tk.M{}
	err := k.GetPayload(&p)
	if err != nil {
		return helper.CreateResult(false, nil, err.Error())
	}

	project := p.Project

	timemax := getMaxRealTime(project, p.Turbine).UTC()
	// alltkmdata := getLastValueFromRaw(timemax, project, p.Turbine)
	// ============== get realtime data =================
	csr, err := DBRealtime().NewQuery().From(new(ScadaRealTimeNew).TableName()).
		Where(dbox.And(dbox.Eq("turbine", p.Turbine), dbox.Eq("projectname", project))).Cursor(nil)
	if err != nil {
		return helper.CreateResult(false, nil, err.Error())
	}
	scadaRealtimeData := []ScadaRealTimeNew{}
	err = csr.Fetch(&scadaRealtimeData, 0, false)
	if err != nil {
		return helper.CreateResult(false, nil, err.Error())
	}
	csr.Close()
	alltkmdata := tk.M{}
	for _, val := range scadaRealtimeData {
		alltkmdata.Set(val.Tags, val.Value)
	}
	// ============== end of get realtime data =================

	arrturbinestatus := GetTurbineStatus(project, p.Turbine)

	alldata.Set("turbine", p.Turbine).Set("lastupdate", timemax.UTC()).Set("projectname", project)
	for key, str := range arrlabel {
		if !alldata.Has(key) {
			alldata.Set(key, defaultValue)
		}

		if str == "" {
			continue
		}

		if alltkmdata.Has(str) {
			if _ival := alltkmdata.GetFloat64(str); _ival != defaultValue && alldata.GetFloat64(key) == defaultValue {
				alldata.Set(key, _ival)
			}
		}
	}
	if _idt, _cond := arrturbinestatus[p.Turbine]; _cond {
		alldata.Set("Turbine Status", _idt.Status)
	} else {
		alldata.Set("Turbine Status", -999)
	}

	t0 := getTimeNow()
	if t0.Sub(timemax.UTC()).Minutes() > 3 {
		alldata.Set("Turbine Status", -999)
	}

	return helper.CreateResult(true, alldata, "success")
}

func (c *MonitoringRealtimeController) GetDataTurbineOld(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson
	k.Config.NoLog = true

	p := struct {
		Project string
		Turbine string
	}{}

	alldata := tk.M{}
	err := k.GetPayload(&p)
	if err != nil {
		return helper.CreateResult(false, nil, err.Error())
	}

	project := p.Project

	timemax := getMaxRealTime(project, p.Turbine).UTC()
	alltkmdata := getLastValueFromRaw(timemax, project, p.Turbine)
	arrturbinestatus := GetTurbineStatus(project, p.Turbine)

	alldata.Set("turbine", p.Turbine).Set("lastupdate", timemax.UTC()).Set("projectname", project)
	for key, str := range arrlabel {
		if !alldata.Has(key) {
			alldata.Set(key, defaultValue)
		}

		if str == "" {
			continue
		}

		if str == "WindSpeed_ms" || str == "ActivePower_kW" {
			// log.Printf(">> %v | %v | %v \n", key, str, alltkmdata.GetFloat64(str))
		}

		if alltkmdata.Has(str) {
			if _ival := alltkmdata.GetFloat64(str); _ival != defaultValue && alldata.GetFloat64(key) == defaultValue {
				alldata.Set(key, _ival)
				if str == "WindSpeed_ms" || str == "ActivePower_kW" {
					log.Printf(">> ival: %v \n", _ival)
				}
			}
		}
	}

	if _idt, _cond := arrturbinestatus[p.Turbine]; _cond {
		alldata.Set("Turbine Status", _idt.Status)
	} else {
		alldata.Set("Turbine Status", -999)
	}

	t0 := getTimeNow()
	if t0.Sub(timemax.UTC()).Minutes() > 3 {
		alldata.Set("Turbine Status", -999)
	}

	return helper.CreateResult(true, alldata, "success")
}

func GetTurbineStatus(project string, turbine string) (res map[string]TurbineStatus) {
	res = map[string]TurbineStatus{}

	filtercond := []*dbox.Filter{}
	if project != "Fleet" || project != "" {
		filtercond = append(filtercond, dbox.Eq("projectname", project))
	}

	if turbine != "" {
		// if project == "Lahori" {
		// 	turbine = project+"_"+turbine
		// }
		filtercond = append(filtercond, dbox.Eq("_id", project+"_"+turbine))
	}

	// rconn := lh.GetConnRealtime()
	// defer rconn.Close()
	rconn := DBRealtime()

	csr, err := rconn.NewQuery().From(new(TurbineStatus).TableName()).
		Where(dbox.And(filtercond...)).
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
		res[result.Turbine] = result
	}

	return
}

func getMaxRealTime(project, turbine string) (timemax time.Time) {
	timemax = time.Time{}

	// rconn := lh.GetConnRealtime()
	// defer rconn.Close()
	rconn := DBRealtime()

	_Query := rconn.NewQuery().From(new(ScadaRealTimeNew).TableName()).
		Aggr(dbox.AggrMax, "$timestamp", "timestamp")

	if turbine != "" {
		_Query = _Query.Group("turbine").
			Where(dbox.And(dbox.Eq("turbine", turbine), dbox.Eq("projectname", project)))
	} else {
		_Query = _Query.Group("projectname").
			Where(dbox.Eq("projectname", project))
	}

	csr, err := _Query.Cursor(nil)

	if err != nil {
		return
	}

	tkmgroup := tk.M{}
	err = csr.Fetch(&tkmgroup, 1, false)
	if err != nil {
		return
	}
	csr.Close()

	timemax = tkmgroup.Get("timestamp", time.Time{}).(time.Time)

	return
}

func getNext10Min(current time.Time) time.Time {
	date1, _ := time.Parse("2006-01-02", current.Format("2006-01-02"))

	thour := current.Hour()
	tminute := current.Minute()
	tsecond := current.Second()
	tminutevalue := float64(tminute) + tk.Div(float64(tsecond), 60.0)
	tminutecategory := tk.ToInt(tk.RoundingUp64(tk.Div(tminutevalue, 10), 0)*10, "0")
	if tminutecategory == 60 {
		tminutecategory = 0
		thour = thour + 1
	}
	newTimeStamp := date1.Add(time.Duration(thour) * time.Hour).Add(time.Duration(tminutecategory) * time.Minute)
	timestampconverted := newTimeStamp.UTC()

	return timestampconverted
}

func getLastValueFromRaw(timemax time.Time, project string, turbine string) (tkm tk.M) {
	tkm = tk.M{}
	timeFolder := getNext10Min(timemax).UTC()
	aTimeFolder := []time.Time{timeFolder.Add(time.Minute * -10), timeFolder}

	for _, _tFolder := range aTimeFolder {
		fullpath := filepath.Join(helper.GetHFDFolder(),
			strings.ToLower(project),
			_tFolder.Format("20060102"), // "20170210",
			_tFolder.Format("15"),       // "11",
			_tFolder.Format("1504"),     // "1120",
		)

		// log.Printf(">> %v \n", fullpath)

		afile := getListFile(fullpath)
		for _, _file := range afile {
			ffile := filepath.Join(fullpath, _file)
			loadFileByTurbine(turbine, ffile, tkm)
		}
	}

	return
}

func getListFile(dir string) (_arrfile []string) {
	_arrfile = []string{}
	_pattern := "^(.*)(\\.[Cc][Ss][Vv])$"

	files, e := ioutil.ReadDir(dir)
	if e != nil {
		tk.Printfn("Get list file found %s", e.Error())
		return
	}

	icount := 0
	for _, file := range files {
		icount++
		filename := file.Name()
		if cond, _ := regexp.MatchString(_pattern, filename); cond {
			_arrfile = append(_arrfile, filename)
		}
	}

	return
}

func loadFileByTurbine(turbine, _fpath string, tkm tk.M) {
	_file, err := os.Open(_fpath)
	if err != nil {
		tk.Printfn("Open %s found %s", _fpath, err.Error())
		return
	}

	scanner := bufio.NewScanner(_file)
	for scanner.Scan() {

		_tData := strings.Split(scanner.Text(), ",")
		if len(_tData) < 4 || _tData[1] != turbine {
			continue
		}

		_val := tk.ToFloat64(_tData[3], 6, tk.RoundingAuto)
		if _val == defaultValue {
			continue
		}

		tkm.Set(_tData[2], _val)

	}

	if err := scanner.Err(); err != nil {
		tk.Printfn("Fetch %s found %s", _fpath, err.Error())
	}

	_file.Close()
	return
}

func getTimeNow() (tNow time.Time) {
	config := lh.ReadConfig()

	loc, err := time.LoadLocation(config["ReadTimeLoc"])
	_Now := time.Now().UTC().Add(-time.Minute * 330)
	if err != nil {
		tk.Printfn("Get time in %s found %s", config["ReadTimeLoc"], err.Error())
	} else {
		_Now = time.Now().In(loc)
	}

	tNow = time.Date(_Now.Year(), _Now.Month(), _Now.Day(), _Now.Hour(), _Now.Minute(), _Now.Second(), _Now.Nanosecond(), time.UTC)
	// tNow = tNow.Add(-10 * time.Minute)
	return
}

/*
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

	alldata, allturbine := tk.M{}, tk.M{}
	arrfield := []string{"ActivePower", "WindSpeed", "WindDirection", "NacellePosition", "Temperature",
		"PitchAngle", "RotorRPM"}
	lastUpdate := time.Time{}
	PowerGen, AvgWindSpeed, CountWS := float64(0), float64(0), float64(0)
	turbinedown := 0
	t0 := time.Now().UTC()

	arrturbinestatus := GetTurbineStatus(project)

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

				if t0.Sub(tstamp.UTC()).Minutes() <= 3 {
					aturbine.Set("DataComing", 1)
				}

				if tstamp.After(utime) {
					aturbine.Set("TimeUpdate", tstamp.UTC())
				}

				if tstamp.After(lastUpdate) {
					lastUpdate = tstamp.UTC()
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
			Set("AlarmUpdate", arrturbinestatus[strturbine].TimeUpdate.UTC())
		if arrturbinestatus[strturbine].Status == 0 {
			turbinedown += 1
		}

		arrturbine := alldata.Get(_tkm.GetString("feeder"), []tk.M{}).([]tk.M)
		arrturbine = append(arrturbine, aturbine)
		alldata.Set(_tkm.GetString("feeder"), arrturbine)

		lturbine := allturbine.Get(_tkm.GetString("feeder"), []string{}).([]string)
		lturbine = append(lturbine, strturbine)
		sort.Strings(lturbine)
		allturbine.Set(_tkm.GetString("feeder"), lturbine)
	}

	rtkm.Set("ListOfTurbine", allturbine)
	rtkm.Set("Data", alldata)
	rtkm.Set("TimeStamp", lastUpdate)
	rtkm.Set("PowerGeneration", PowerGen)
	rtkm.Set("AvgWindSpeed", tk.Div(AvgWindSpeed, CountWS))
	rtkm.Set("PLF", tk.Div(PowerGen, (50400*100)))
	rtkm.Set("TurbineActive", len(_result)-turbinedown)
	rtkm.Set("TurbineDown", turbinedown)

	return
}
*/

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

/*
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
*/

// func (c *MonitoringRealtimeController) GetWindRoseData(k *knot.WebContext) interface{} {
// 	k.Config.OutputType = knot.OutputJson
// 	k.Config.NoLog = true
// 	sessid := k.Session("sessionid", "")
// 	accs := "GetWindRoseData"

// 	// WindRoseResult = []tk.M{}

// 	p := struct {
// 		Turbine string
// 	}{}
// 	e := k.GetPayload(&p)
// 	if e != nil {
// 		WriteLog(sessid, accs, e.Error())
// 	}

// 	query := []tk.M{}
// 	pipes := []tk.M{}
// 	section = 12
// 	getFullWSCategory()

// 	data := []MiniScada{}
// 	_data := MiniScada{}

// 	lastDateData, e := time.Parse(time.RFC3339, "2017-01-22T00:00:00+00:00")
// 	if e != nil {
// 		return helper.CreateResult(false, nil, e.Error())
// 	}

// 	turbines := p.Turbine
// 	defaultValue := -999999.00

// 	groupdata := tk.M{}
// 	groupdata.Set("Name", turbines)

// 	query = append(query, tk.M{"_id": tk.M{"$ne": nil}})
// 	query = append(query, tk.M{"nacelleposition": tk.M{"$ne": defaultValue}})
// 	query = append(query, tk.M{"dateinfo.dateid": lastDateData})
// 	query = append(query, tk.M{"turbine": turbines})
// 	pipes = append(pipes, tk.M{"$match": tk.M{"$and": query}})
// 	pipes = append(pipes, tk.M{"$project": tk.M{"nacelleposition": 1, "windspeed": 1}})
// 	csr, e := DB().Connection.NewQuery().From(new(ScadaRealTime).TableName()).
// 		Command("pipe", pipes).Cursor(nil)

// 	for {
// 		e = csr.Fetch(&_data, 1, false)
// 		if e != nil {
// 			break
// 		}
// 		data = append(data, _data)
// 	}
// 	csr.Close()

// 	if tk.SliceLen(data) > 0 {
// 		totalDuration := float64(len(data)) /* Tot data * 2 for get total minutes*/
// 		datas := cr.From(&data).Apply(func(x interface{}) interface{} {
// 			dt := x.(MiniScada)
// 			var di DataItems

// 			dirNo, dirDesc := getDirection(dt.NacellePosition, section)
// 			wsNo, wsDesc := getWsCategory(dt.WindSpeed)

// 			di.DirectionNo = dirNo
// 			di.DirectionDesc = dirDesc
// 			di.WsCategoryNo = wsNo
// 			di.WsCategoryDesc = wsDesc
// 			di.Frequency = 1

// 			return di
// 		}).Exec().Group(func(x interface{}) interface{} {
// 			dt := x.(DataItems)

// 			var dig DataItemsGroup
// 			dig.DirectionNo = dt.DirectionNo
// 			dig.DirectionDesc = dt.DirectionDesc
// 			dig.WsCategoryNo = dt.WsCategoryNo
// 			dig.WsCategoryDesc = dt.WsCategoryDesc

// 			return dig
// 		}, nil).Exec()

// 		dts := datas.Apply(func(x interface{}) interface{} {
// 			kv := x.(cr.KV)
// 			vv := kv.Key.(DataItemsGroup)
// 			vs := kv.Value.([]DataItems)

// 			sumFreq := cr.From(&vs).Sum(func(x interface{}) interface{} {
// 				dt := x.(DataItems)
// 				return dt.Frequency
// 			}).Exec().Result.Sum

// 			var di DataItemsResult

// 			di.DirectionNo = vv.DirectionNo
// 			di.DirectionDesc = vv.DirectionDesc
// 			di.WsCategoryNo = vv.WsCategoryNo
// 			di.WsCategoryDesc = vv.WsCategoryDesc
// 			di.Hours = tk.Div(sumFreq, 6.0)
// 			di.Contribution = tk.RoundingAuto64(tk.Div(sumFreq, totalDuration)*100.0, 2)

// 			// key := turbines + "_" + tk.ToString(di.DirectionNo)

// 			// if !tkMaxVal.Has(key) {
// 			// 	tkMaxVal.Set(key, di.Contribution)
// 			// } else {
// 			// 	tkMaxVal.Set(key, tkMaxVal.GetFloat64(key)+di.Contribution)
// 			// }

// 			di.Frequency = int(sumFreq)

// 			return di
// 		}).Exec()

// 		results := dts.Result.Data().([]DataItemsResult)
// 		wsCategoryList := []string{}
// 		for _, dataRes := range results {
// 			wsCategoryList = append(wsCategoryList, tk.ToString(dataRes.DirectionNo)+
// 				"_"+tk.ToString(dataRes.WsCategoryNo)+"_"+dataRes.WsCategoryDesc)
// 		}
// 		splitCatList := []string{}
// 		for _, wsCat := range fullWSCatList {
// 			if !tk.HasMember(wsCategoryList, wsCat) {
// 				splitCatList = strings.Split(wsCat, "_")
// 				emptyRes := DataItemsResult{}
// 				emptyRes.DirectionNo = tk.ToInt(splitCatList[0], tk.RoundingAuto)
// 				divider := section

// 				emptyRes.DirectionDesc = (360 / divider) * emptyRes.DirectionNo
// 				emptyRes.WsCategoryNo = tk.ToInt(splitCatList[1], tk.RoundingAuto)
// 				emptyRes.WsCategoryDesc = splitCatList[2]
// 				results = append(results, emptyRes)
// 			}
// 		}
// 		groupdata.Set("Data", results)

// 		// tk.Printf("results : %s \n", tk.SliceLen(results))
// 		// tk.Printf("fullWSCatList : %s \n", fullWSCatList)

// 		// WindRoseResult = append(WindRoseResult, groupdata)
// 	} else {
// 		splitCatList := []string{}
// 		results := []DataItemsResult{}
// 		for _, wsCat := range fullWSCatList {
// 			splitCatList = strings.Split(wsCat, "_")
// 			emptyRes := DataItemsResult{}
// 			emptyRes.DirectionNo = tk.ToInt(splitCatList[0], tk.RoundingAuto)
// 			divider := section

// 			emptyRes.DirectionDesc = (360 / divider) * emptyRes.DirectionNo
// 			emptyRes.WsCategoryNo = tk.ToInt(splitCatList[1], tk.RoundingAuto)
// 			emptyRes.WsCategoryDesc = splitCatList[2]
// 			results = append(results, emptyRes)
// 		}
// 		groupdata.Set("Data", results)
// 		// WindRoseResult = append(WindRoseResult, groupdata)
// 	}

// 	// tk.Printf("groupdata : %s \n", tk.SliceLen(groupdata))

// 	dataresult := struct {
// 		WindRose tk.M
// 	}{
// 		WindRose: groupdata,
// 	}

// 	return helper.CreateResult(true, dataresult, "success")
// }

// func (c *MonitoringRealtimeController) GetDataLine(k *knot.WebContext) interface{} {
// 	k.Config.OutputType = knot.OutputJson
// 	k.Config.NoLog = true
// 	sessid := k.Session("sessionid", "")
// 	accs := "GetDataLine"

// 	var (
// 		pipes      []tk.M
// 		filter     []*dbox.Filter
// 		list       []tk.M
// 		dataSeries []tk.M
// 	)

// 	p := struct {
// 		Turbine string
// 	}{}
// 	e := k.GetPayload(&p)
// 	if e != nil {
// 		WriteLog(sessid, accs, e.Error())
// 	}

// 	lastDateData, e := time.Parse(time.RFC3339, "2017-01-22T00:00:00+00:00")
// 	if e != nil {
// 		return helper.CreateResult(false, nil, e.Error())
// 	}

// 	turbines := p.Turbine
// 	defaultValue := -999999.00

// 	pipes = append(pipes, tk.M{"$group": tk.M{"_id": tk.M{"colId": "$timestamp", "Turbine": "$turbine"},
// 		"avgwindspeed": tk.M{"$avg": "$windspeed"},
// 		"sumwindspeed": tk.M{"$sum": "$windspeed"},
// 		"activepower":  tk.M{"$sum": "$activepower"},
// 		"rotorrpm":     tk.M{"$sum": "$rotorrpm"},
// 		"totaldata":    tk.M{"$sum": 1}}})
// 	pipes = append(pipes, tk.M{"$sort": tk.M{"_id": 1}})

// 	filter = nil
// 	filter = append(filter, dbox.Ne("_id", ""))
// 	filter = append(filter, dbox.Eq("dateinfo.dateid", lastDateData))
// 	filter = append(filter, dbox.Eq("turbine", turbines))
// 	filter = append(filter, dbox.Ne("activepower", defaultValue))
// 	filter = append(filter, dbox.Ne("windspeed", defaultValue))

// 	csr, e := DB().Connection.NewQuery().
// 		From(new(ScadaRealTime).TableName()).
// 		Command("pipe", pipes).
// 		Where(dbox.And(filter...)).
// 		Cursor(nil)

// 	if e != nil {
// 		return helper.CreateResult(false, nil, e.Error())
// 	}
// 	e = csr.Fetch(&list, 0, false)
// 	defer csr.Close()

// 	totactivepower := 0.0
// 	totwindspeed := 0.0
// 	totrotorrpm := 0.0
// 	totData := 0.0
// 	dataMonitoring := tk.M{}
// 	for _, val := range list {

// 		seriesData := tk.M{}
// 		avgwindspeed := val.GetFloat64("avgwindspeed")
// 		sumwindspeed := val.GetFloat64("sumwindspeed")
// 		activepower := val.GetFloat64("activepower")
// 		rotorrpm := val.GetFloat64("rotorrpm")
// 		totaldata := val.GetFloat64("totaldata")
// 		idD := val.Get("_id").(tk.M)
// 		Turbine := idD.Get("Turbine")
// 		timestamp := idD.Get("colId").(time.Time).UTC().Format("2006-01-02 15:04:05")

// 		seriesData.Set("turbine", Turbine)
// 		seriesData.Set("timestamp", timestamp)
// 		seriesData.Set("activepower", tk.Div(activepower, 1000.0))
// 		seriesData.Set("avgwindspeed", avgwindspeed)

// 		dataSeries = append(dataSeries, seriesData)

// 		totactivepower = totactivepower + activepower
// 		totwindspeed = totwindspeed + sumwindspeed
// 		totrotorrpm = totrotorrpm + rotorrpm
// 		totData = totData + totaldata

// 	}

// 	dataMonitoring.Set("Power", tk.Div(totactivepower, 1000.0))
// 	dataMonitoring.Set("WindSpeed", tk.Div(totwindspeed, totData))
// 	dataMonitoring.Set("RotorRpm", totrotorrpm)

// 	data := struct {
// 		Data       []tk.M
// 		Monitoring tk.M
// 	}{
// 		Data:       dataSeries,
// 		Monitoring: dataMonitoring,
// 	}

// 	return helper.CreateResult(true, data, "success")
// }
