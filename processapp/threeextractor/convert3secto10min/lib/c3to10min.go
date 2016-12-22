package c3to10min

import (
	"bufio"
	"log"
	"math"
	"os"
	"strings"

	"github.com/eaciit/dbox"
	tk "github.com/eaciit/toolkit"

	"eaciit/wfdemo-git/library/helper"
	. "eaciit/wfdemo-git/library/models"

	_ "github.com/eaciit/dbox/dbc/mongo"

	"errors"
	"reflect"
	"time"
	// "github.com/eaciit/orm"
	// dc "eaciit/wfdemo-git/processapp/threeextractor/dataconversion"
)

var (
	emptyValueSmall = -0.000001
	emptyValueBig   = -9999999.0

	wd = func() string {
		d, _ := os.Getwd()
		return d + separator
	}()
	separator = string(os.PathSeparator)

	structlist = []string{"Fast_ActivePower_kW", "Fast_WindSpeed_ms", "Slow_NacellePos", "Slow_WindDirection",
		"Fast_CurrentL3", "Fast_CurrentL1", "Fast_ActivePowerSetpoint_kW", "Fast_CurrentL2", "Fast_DrTrVibValue",
		"Fast_GenSpeed_RPM", "Fast_PitchAccuV1", "Fast_PitchAngle", "Fast_PitchAngle3", "Fast_PitchAngle2", "Fast_PitchConvCurrent1",
		"Fast_PitchConvCurrent3", "Fast_PitchConvCurrent2", "Fast_PowerFactor", "Fast_ReactivePowerSetpointPPC_kVA", "Fast_ReactivePower_kVAr",
		"Fast_RotorSpeed_RPM", "Fast_VoltageL1", "Fast_VoltageL2", "Slow_CapableCapacitiveReactPwr_kVAr", "Slow_CapableInductiveReactPwr_kVAr", "Slow_DateTime_Sec",
		"Fast_PitchAngle1", "Fast_VoltageL3", "Slow_CapableCapacitivePwrFactor", "Fast_Total_Production_kWh", "Fast_Total_Prod_Day_kWh", "Fast_Total_Prod_Month_kWh",
		"Fast_ActivePowerOutPWCSell_kW", "Fast_Frequency_Hz", "Slow_TempG1L2", "Slow_TempG1L3", "Slow_TempGearBoxHSSDE", "Slow_TempGearBoxIMSNDE",
		"Slow_TempOutdoor", "Fast_PitchAccuV3", "Slow_TotalTurbineActiveHours", "Slow_TotalTurbineOKHours", "Slow_TotalTurbineTimeAllHours",
		"Slow_TempG1L1", "Slow_TempGearBoxOilSump", "Fast_PitchAccuV2", "Slow_TotalGridOkHours", "Slow_TotalActPowerOut_kWh",
		"Fast_YawService", "Fast_YawAngle", "Slow_CapableInductivePwrFactor", "Slow_TempGearBoxHSSNDE", "Slow_TempHubBearing",
		"Slow_TotalG1ActiveHours", "Slow_TotalActPowerOutG1_kWh", "Slow_TotalReactPowerInG1_kVArh", "Slow_NacelleDrill",
		"Slow_TempGearBoxIMSDE", "Fast_Total_Operating_hrs", "Slow_TempNacelle", "Fast_Total_Grid_OK_hrs", "Fast_Total_WTG_OK_hrs",
		"Slow_TempCabinetTopBox", "Slow_TempGeneratorBearingNDE", "Fast_Total_Access_hrs", "Slow_TempBottomPowerSection", "Slow_TempGeneratorBearingDE",
		"Slow_TotalReactPowerIn_kVArh", "Slow_TempBottomControlSection", "Slow_TempConv1", "Fast_ActivePowerRated_kW", "Fast_NodeIP",
		"Fast_PitchSpeed1", "Slow_CFCardSize", "Slow_CPU_Number", "Slow_CFCardSpaceLeft", "Slow_TempBottomCapSection",
		"Slow_RatedPower", "Slow_TempConv3", "Slow_TempConv2", "Slow_TotalActPowerIn_kWh", "Slow_TotalActPowerInG1_kWh", "Slow_TotalActPowerInG2_kWh",
		"Slow_TotalActPowerOutG2_kWh", "Slow_TotalG2ActiveHours", "Slow_TotalReactPowerInG2_kVArh", "Slow_TotalReactPowerOut_kVArh", "Slow_UTCoffset_int"}

	// intstartdate = int(20160801)
	// intenddate   = int(20160831)

	startdate, enddate, cdate = time.Time{}, time.Time{}, time.Time{} //date, year
	arridate                  = []string{}                            //adate
	strfilename               = string("")                            //file

	config = map[string]string{}
)

//>> Param >>
// selector = date | file | adate | year
// >> selector = date
// sdate / edate exp. "20160821" / "20160831"
// exp tk.M{}.Set("selector", "date").Set("sdate", "20160821").Set("edate","20160831")
// >> selector = file
// file exp. "DataFile20160821-01.csv"
// exp tk.M{}.Set("selector", "file").Set("file", "DataFile20160821-01.csv")
// >> selector = adate
// adate exp []string{"20160821", "20160831"}
// exp tk.M{}.Set("selector", "adate").Set("adate", []string{"20160821", "20160831"})
// >> selector = year
// year exp 2016
// exp tk.M{}.Set("selector", "year").Set("year", 2016)
//>>>>>>>>>>>

func Generate(param tk.M) (ferr error) {

	linit := string("")
	ferr, linit = checkparam(param)
	if ferr != nil {
		log.Println(ferr.Error())
		return
	}

	sselector := param.GetString("selector")

	log.Println(tk.Sprintf("Convert Data for %s", linit))
	config = ReadConfig()

	log.Println(tk.Sprintf("Connect to %s, %s", config["host"], config["database"]))
	conn, err := PrepareConnection()
	if err != nil {
		tk.Println(err)
	}
	defer conn.Close()

	start := time.Now()

	//=============================
	// Prepare go routine to calculate data
	//=============================

	sresult := make(chan int, 200)
	sdata := make(chan time.Time, 200)
	for i := 0; i < 20; i++ {
		go calcdata(i, sdata, sresult)
	}

	//=============================

	iarr := int(0)
	isbreak := false
	for {
		t0 := time.Now()

		match := tk.M{}
		scond := ""

		//date | file | adate | year
		switch sselector {
		case "date", "year":
			if cdate.After(enddate) {
				isbreak = true
			} else {
				match.Set("$and", []tk.M{tk.M{}.Set("timestampconverted", tk.M{"$gt": cdate}),
					tk.M{}.Set("timestampconverted", tk.M{"$lte": cdate.AddDate(0, 0, 1)})})
				scond = cdate.Format("2006-01-02")
			}
			cdate = cdate.AddDate(0, 0, 1)
		case "file":
			if iarr > 0 {
				isbreak = true
			} else {
				match.Set("file", strfilename)
				scond = strfilename
			}
			iarr++
		case "adate":
			if iarr >= len(arridate) {
				isbreak = true
			} else {
				idate := tk.String2Date(arridate[iarr], "YYYYMMdd").UTC()
				match.Set("$and", []tk.M{tk.M{}.Set("timestampconverted", tk.M{"$gt": idate}),
					tk.M{}.Set("timestampconverted", tk.M{"$lte": idate.AddDate(0, 0, 1)})})
				scond = idate.Format("2006-01-02")
			}
			iarr++
		}

		if isbreak {
			break
		}

		log.Println(tk.Sprintf("Preparing Process Data for %s ", scond))
		arrtimeinterval := getinterval(new(ScadaThreeSecs).TableName(), match)
		count := len(arrtimeinterval)

		log.Println(tk.Sprintf("Found Interval Data : %d in %s", count, time.Since(t0).String()))
		step := getstep(count)
		for _, _v := range arrtimeinterval {
			sdata <- _v
		}

		_countdata := int(0)
		for i := 0; i < count; i++ {
			_countdata += <-sresult

			if i%step == 0 {
				log.Println(tk.Sprintf("Saved %d of %d (%d pct) in %s",
					i, count, i*100/count, time.Since(t0).String()))
			}
		}

		log.Println(tk.Sprintf("Done Process Data for %v, in %s total %d rows saved",
			scond, time.Since(t0).String(), _countdata))

	}

	close(sdata)
	close(sresult)

	//=============================

	log.Printf("All data conversion done in %s \n",
		time.Since(start).String())

	return
}

func checkparam(_param tk.M) (_ferr error, _linit string) {
	_ferr = nil
	_linit = string("")

	_select := _param.GetString("selector")
	switch _select {
	case "date":
		if !_param.Has("sdate") || !_param.Has("edate") {
			_ferr = errors.New(tk.Sprintf("%s value is not found", _select))
			return
		}
		_linit = tk.Sprintf("%s to %s", _param.GetString("sdate"), _param.GetString("edate"))
		startdate = tk.String2Date(_param.GetString("sdate"), "YYYYMMdd").UTC()
		cdate = startdate
		enddate = tk.String2Date(_param.GetString("edate"), "YYYYMMdd").UTC()
	case "adate":
		if !_param.Has(_select) {
			_ferr = errors.New(tk.Sprintf("%s value is not found", _select))
			return
		}
		_linit = strings.Join(_param[_select].([]string), ",")
		arridate = _param[_select].([]string)
	case "file":
		if !_param.Has(_select) {
			_ferr = errors.New(tk.Sprintf("%s value is not found", _select))
			return
		}
		_linit = _param.GetString(_select)
		strfilename = _param.GetString(_select)
	case "year":
		if !_param.Has(_select) {
			_ferr = errors.New(tk.Sprintf("%s value is not found", _select))
			return
		}
		_linit = _param.GetString(_select)

		tdate := tk.Sprintf("%d0101", _param.GetInt(_select))
		startdate = tk.String2Date(tdate, "YYYYMMdd").UTC()
		cdate = startdate

		tdate = tk.Sprintf("%d0101", _param.GetInt(_select)+1)
		enddate = tk.String2Date(tdate, "YYYYMMdd").UTC()
	}

	return
}

func calcdata(wi int, jobs <-chan time.Time, result chan<- int) {
	// workerconn, _ := PrepareConnection()
	// defer workerconn.Close()

	dtablename := tk.Sprintf("%s", new(ScadaThreeSecs).TableName())

	sresult := make(chan int, 100)
	sdata := make(chan ScadaConvTenMin, 100)
	for i := 0; i < 5; i++ {
		go workersave(i, sdata, sresult)
	}

	_tinterval := time.Time{}
	for _tinterval = range jobs {
		workerconn, _ := PrepareConnection()
		csr, e := workerconn.NewQuery().
			Select().From(dtablename).
			Where(dbox.Eq("timestampconverted", _tinterval)).
			Cursor(nil)

		if e != nil {
			log.Printf("ERRROR: %v | for data [%v] \n", e.Error(), _tinterval)
			continue
		}

		mapscada3avg := make(map[string]*ScadaThreeSecs, 0)
		mapavgcount := make(map[string]tk.M, 0)
		for {
			sts := new(ScadaThreeSecs)
			e = csr.Fetch(sts, 1, false)
			if e != nil {
				break
			}

			//Round time Second
			timeStamp := sts.TimeStamp1.UTC()
			seconds := tk.Div(tk.ToFloat64(timeStamp.Nanosecond(), 1, tk.RoundingAuto), 1000000000)
			secondsInt := tk.ToInt(seconds, tk.RoundingAuto)
			newTimeTmp := timeStamp.Add(time.Duration(secondsInt) * time.Second)
			TimeStampSecondGroup, _ := time.Parse("20060102 15:04:05", newTimeTmp.Format("20060102 15:04:05"))

			//timeStampStr := m.TimeStamp.UTC().Format("060102_1504") //m.ID = timeStampStr + "#" + m.ProjectName + "#" + m.Turbine
			key := TimeStampSecondGroup.UTC().Format("20060102_150405") + "#" + sts.ProjectName + "#" + sts.Turbine
			stsavg := new(ScadaThreeSecs)
			if _, exist := mapscada3avg[key]; exist {
				stsavg = mapscada3avg[key]
			}

			tkmcount := tk.M{}
			if _, exist := mapavgcount[key]; exist {
				tkmcount = mapavgcount[key]
			}

			fillto3secaggr(sts, stsavg, tkmcount, key)
			mapscada3avg[key] = stsavg
			mapavgcount[key] = tkmcount
		}

		csr.Close()
		workerconn.Close()

		// tk.Println(">>>>>>>>>>> length mapscada3avg : ", len(mapscada3avg))
		mapscadatenmin := make(map[string]*ScadaConvTenMin, 0)
		for _key, sts := range mapscada3avg {

			// ===================================================================
			// ===Count and set aggr=====================================================

			tkmcount := tk.M{}
			if _, exist := mapavgcount[_key]; exist {
				tkmcount = mapavgcount[_key]
			}

			for _, _str := range structlist {
				rval := reflect.ValueOf(sts).Elem().FieldByName(_str)

				if !rval.IsValid() {
					continue
				}

				ival := rval.Float()
				if ival == emptyValueBig {
					continue
				}

				if _str != "Fast_YawService" {
					ival = tk.Div(ival, tkmcount.GetFloat64(_str))
				}

				rval.SetFloat(ival)
			}

			// ===================================================================
			// ===================================================================

			//timeStampStr := m.TimeStamp.UTC().Format("060102_1504")
			//m.ID = timeStampStr + "#" + m.ProjectName + "#" + m.Turbine
			key := sts.TimeStampConverted.UTC().Format("20060102_1504") + "#" + sts.ProjectName + "#" + sts.Turbine
			sctm := new(ScadaConvTenMin)
			if _, exist := mapscadatenmin[key]; exist {
				sctm = mapscadatenmin[key]
			}

			fillto10min(sts, sctm, key)

			mapscadatenmin[key] = sctm
		}

		// sresult := make(chan int, 100)
		// sdata := make(chan ScadaConvTenMin, 100)
		// for i := 0; i < 5; i++ {
		// 	go workersave(i, sdata, sresult)
		// }

		for _, sctm := range mapscadatenmin {
			sdata <- *sctm
		}

		// close(sdata)

		for i := 0; i < len(mapscadatenmin); i++ {
			<-sresult
		}

		// close(sresult)

		result <- len(mapscadatenmin)
		// tk.Printfn("[DONE] Interval %v with %d data", _tinterval, len(mapscadatenmin))
	}

	close(sdata)
	close(sresult)

	return
}

func workersave(wi int, jobs <-chan ScadaConvTenMin, result chan<- int) {
	var workerconn dbox.IConnection
	for {
		var err error
		workerconn, err = PrepareConnection()
		if err == nil {
			break
		} else {
			tk.Printfn("==#DB-ERRCONN==\n %s \n", err.Error())
		}
	}
	defer workerconn.Close()

	dtablename := tk.Sprintf("%s", new(ScadaConvTenMin).TableName())

	qSave := workerconn.NewQuery().
		From(dtablename).
		SetConfig("multiexec", true).
		Save()

	trx := ScadaConvTenMin{}
	for trx = range jobs {
		//Average
		for _, _str := range structlist {
			rfloat := reflect.ValueOf(&trx).Elem().FieldByName(_str)
			ifloat := emptyValueBig
			if rfloat.IsValid() {
				ifloat = rfloat.Float()
			}

			_strdev := tk.Sprintf("%s_StdDev", _str)
			if ifloat != emptyValueBig {
				icount := float64(reflect.ValueOf(&trx).Elem().FieldByName(tk.Sprintf("%s_Count", _str)).Int())
				tVal := tk.Div(ifloat, icount)
				if _str == "Fast_YawService" {
					tVal = reflect.ValueOf(&trx).Elem().FieldByName(tk.Sprintf("%s_Min", _str)).Float()
				}
				reflect.ValueOf(&trx).Elem().FieldByName(_str).SetFloat(tVal)

				istddev2 := reflect.ValueOf(&trx).Elem().FieldByName(_strdev).Float()
				istddev2 = (icount * istddev2) - (ifloat * ifloat)
				// tk.Printfn(" >>>>>>> istddev2 : %#v ", istddev2)

				istddev2 = tk.Div(istddev2, (icount * icount))
				// tk.Printfn(" div >>>>>>> istddev2 : %#v ", istddev2)

				istddev := math.Sqrt(istddev2)
				// tk.Printfn(" sqrt >>>>>>> istddev : %#v ", istddev)
				reflect.ValueOf(&trx).Elem().FieldByName(_strdev).SetFloat(istddev)
			}
		}

		err := qSave.Exec(tk.M{}.Set("data", trx))
		if err != nil {
			tk.Println(err)
		}

		result <- 1
	}

	return
}

func PrepareConnection() (dbox.IConnection, error) {
	ci := &dbox.ConnectionInfo{config["host"], config["database"], config["username"], config["password"], tk.M{}.Set("timeout", 3000)}
	c, e := dbox.NewConnection("mongo", ci)

	if e != nil {
		return nil, e
	}

	e = c.Connect()
	if e != nil {
		return nil, e
	}

	return c, nil
}

func ReadConfig() map[string]string {
	ret := make(map[string]string)
	file, err := os.Open("../conf" + separator + "app.conf")
	if err == nil {
		reader := bufio.NewReader(file)
		for {
			line, _, e := reader.ReadLine()
			if e != nil {
				break
			}

			sval := strings.Split(string(line), "=")
			ret[sval[0]] = sval[1]
		}
	} else {
		tk.Println(err.Error())
	}

	file.Close()
	return ret
}

func getstep(count int) int {
	v := count / 5
	if v == 0 {
		return 1
	}
	return v
}

func getinterval(tablename string, match tk.M) (arrval []time.Time) {

	// arrval = make([]time.Time, 0)
	// _enddate := cdate.AddDate(0, 0, 1)
	// for cdate.Before(_enddate) {
	// 	cdate = cdate.Add(time.Minute * 10)
	// 	arrval = append(arrval, cdate)
	// }
	conn, err := PrepareConnection()
	if err != nil {
		tk.Println(err)
		return
	}
	defer conn.Close()

	pipes := []tk.M{}
	//ISODate("2016-08-21T19:10:00.000+0000")
	// match := tk.M{}.Set("timestampconverted", tk.M{"$gt": time.Date(2016, 8, 21, 22, 0, 0, 0, time.UTC)})

	group := tk.M{}.Set("_id", "$timestampconverted")
	sort := tk.M{}.Set("_id", 1)

	pipes = append(pipes, tk.M{"$match": match})
	pipes = append(pipes, tk.M{"$group": group})
	pipes = append(pipes, tk.M{"$sort": sort})

	// tk.Printfn(">>>> %v", match)

	csr, e := conn.NewQuery().
		From(tablename).
		Command("pipe", pipes).
		Cursor(nil)

	if e != nil {
		log.Printf("ERRROR: %v \n", e.Error())
		os.Exit(1)
	}
	defer csr.Close()

	arrval = []time.Time{}
	for {
		tkm := tk.M{}
		e = csr.Fetch(&tkm, 1, false)
		if e != nil {
			break
		}
		arrval = append(arrval, tkm.Get("_id", time.Time{}).(time.Time).UTC())
	}

	return
}

func fillto10min(_sts *ScadaThreeSecs, _sctm *ScadaConvTenMin, key string) {

	for _, _str := range structlist {
		_strmax := tk.Sprintf("%s_Max", _str)
		_strmin := tk.Sprintf("%s_Min", _str)

		_strdev := tk.Sprintf("%s_StdDev", _str)
		_strcount := tk.Sprintf("%s_Count", _str)

		if _sctm.ID == "" {
			reflect.ValueOf(_sctm).Elem().FieldByName(_str).SetFloat(emptyValueBig)
			reflect.ValueOf(_sctm).Elem().FieldByName(_strmin).SetFloat(emptyValueBig)
			reflect.ValueOf(_sctm).Elem().FieldByName(_strmax).SetFloat(emptyValueBig)
			reflect.ValueOf(_sctm).Elem().FieldByName(_strdev).SetFloat(emptyValueBig)
		}

		rval := reflect.ValueOf(_sts).Elem().FieldByName(_str)
		ival := emptyValueBig

		if rval.IsValid() {
			ival = rval.Float()
		}
		if ival != emptyValueBig {
			tval := reflect.ValueOf(_sctm).Elem().FieldByName(_str).Float()
			if tval == emptyValueBig {
				tval = 0
			}

			tval += ival
			reflect.ValueOf(_sctm).Elem().FieldByName(_str).SetFloat(tval)

			tval = reflect.ValueOf(_sctm).Elem().FieldByName(_strmax).Float()
			if tval < ival {
				reflect.ValueOf(_sctm).Elem().FieldByName(_strmax).SetFloat(ival)
			}

			tval = reflect.ValueOf(_sctm).Elem().FieldByName(_strmin).Float()
			if tval == emptyValueBig || tval > ival {
				reflect.ValueOf(_sctm).Elem().FieldByName(_strmin).SetFloat(ival)
			}

			tval = reflect.ValueOf(_sctm).Elem().FieldByName(_strdev).Float()
			if tval == emptyValueBig {
				tval = 0
			}
			tval += (ival * ival)
			reflect.ValueOf(_sctm).Elem().FieldByName(_strdev).SetFloat(tval)

			iint := reflect.ValueOf(_sctm).Elem().FieldByName(_strcount).Int() + 1
			reflect.ValueOf(_sctm).Elem().FieldByName(_strcount).SetInt(iint)
		}
	}

	_sctm.ID = key
	_sctm.TimeStamp = _sts.TimeStampConverted
	_sctm.TimeStampInt = int64(tk.ToInt(_sts.TimeStampConverted.Format("20060102150405"), tk.RoundingAuto))
	_sctm.DateInfo = helper.GetDateInfo(_sctm.TimeStamp)

	_sctm.ProjectName = _sts.ProjectName
	_sctm.Turbine = _sts.Turbine

	_sctm.File = _sts.File
	// No    int
	_sctm.Count += 1

	return
}

func fillto3secaggr(_sts *ScadaThreeSecs, _stsavg *ScadaThreeSecs, tkm tk.M, key string) {

	for _, _str := range structlist {
		_rstavg := reflect.ValueOf(_stsavg).Elem().FieldByName(_str)
		if !_rstavg.IsValid() {
			continue
		}

		if _stsavg.ID == "" {
			_rstavg.SetFloat(emptyValueBig)
		}

		rval := reflect.ValueOf(_sts).Elem().FieldByName(_str)
		ival := emptyValueBig

		if rval.IsValid() {
			ival = rval.Float()
		}

		if ival != emptyValueBig {
			tval := _rstavg.Float()
			if tval == emptyValueBig {
				tval = 0
			}

			if _str == "Fast_YawService" {
				if ival < tval {
					tval = ival
				}
			} else {
				tval += ival
			}

			_rstavg.SetFloat(tval)

			iint := tkm.GetInt(_str) + 1
			tkm.Set(_str, iint)
		}
	}

	_stsavg.ID = key

	_stsavg.TimeStampConverted = _sts.TimeStampConverted
	_stsavg.TimeStampSecondGroup = _sts.TimeStampSecondGroup

	_stsavg.ProjectName = _sts.ProjectName
	_stsavg.Turbine = _sts.Turbine

	_stsavg.File = _sts.File
	// _stsavg.Count += 1

	return
}
