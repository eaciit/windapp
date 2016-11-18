package main

import (
	"bufio"
	"log"
	"math"
	"os"
	// "strconv"
	"strings"

	"github.com/eaciit/dbox"
	tk "github.com/eaciit/toolkit"

	"github.com/eaciit/windapp/library/helper"
	. "github.com/eaciit/windapp/library/models"

	_ "github.com/eaciit/dbox/dbc/mongo"

	"flag"
	"reflect"
	"time"
	// "github.com/eaciit/orm"
	// dc "github.com/eaciit/windapp/processapp/threeextractor/dataconversion"
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

	intstartdate = int(20160821)
	intenddate   = int(20160821)
)

func main() {

	flag.IntVar(&intstartdate, "sdate", 20160821, "Start date for processing data")
	flag.IntVar(&intenddate, "edate", 20160821, "End date for processing data")
	flag.Parse()

	startdate := tk.String2Date(tk.Sprintf("%d", intstartdate), "YYYYMMdd").UTC()
	enddate := tk.String2Date(tk.Sprintf("%d", intenddate), "YYYYMMdd").UTC()

	conn, err := PrepareConnection()
	if err != nil {
		tk.Println(err)
	}
	defer conn.Close()
	// ctx := orm.New(conn)

	start := time.Now()
	log.Println(tk.Sprintf("Convert Data from %v to %v", startdate, enddate))
	//=============================
	cdate := startdate
	for !cdate.After(enddate) {
		t0 := time.Now()
		log.Println(tk.Sprintf("Preparing Process Data for %v ", cdate))

		arrtimeinterval := getinterval(new(ScadaThreeSecs).TableName(), cdate)
		count := len(arrtimeinterval)

		log.Println("Found Interval Data : ", count)

		step := getstep(count)
		sresult := make(chan int, count)
		sdata := make(chan time.Time, count)
		for i := 0; i < 10; i++ {
			go calcdata(i, sdata, sresult)
		}

		for _, _v := range arrtimeinterval {
			sdata <- _v
		}

		close(sdata)

		_countdata := int(0)
		for i := 0; i < count; i++ {
			_countdata += <-sresult

			if i%step == 0 {
				log.Println(tk.Sprintf("Saved %d of %d (%d pct) in %s",
					i, count, i*100/count, time.Since(t0).String()))
			}
		}

		cdate = cdate.AddDate(0, 0, 1)
		log.Println(tk.Sprintf("Done Process Data for %v, in %s total %d rows saved",
			cdate, time.Since(t0).String(), _countdata))
	}

	//=============================

	log.Printf("All data conversion done in %s \n",
		time.Since(start).String())
}

func calcdata(wi int, jobs <-chan time.Time, result chan<- int) {
	workerconn, _ := PrepareConnection()
	defer workerconn.Close()

	sresult := make(chan int, 100)
	sdata := make(chan ScadaConvTenMin, 100)
	for i := 0; i < 5; i++ {
		go workersave(i, sdata, sresult)
	}

	_tinterval := time.Time{}
	for _tinterval = range jobs {
		csr, e := workerconn.NewQuery().
			Select().From(new(ScadaThreeSecs).TableName()).
			Where(dbox.Eq("timestampconverted", _tinterval)).
			Cursor(nil)

		if e != nil {
			log.Printf("ERRROR: %v | for data [%v] \n", e.Error(), _tinterval)
			continue
		}
		defer csr.Close()

		mapscadatenmin := make(map[string]*ScadaConvTenMin, 0)
		for {
			sts := new(ScadaThreeSecs)
			e = csr.Fetch(sts, 1, false)
			if e != nil {
				break
			}

			//timeStampStr := m.TimeStamp.UTC().Format("060102_1504")
			//m.ID = timeStampStr + "#" + m.ProjectName + "#" + m.Turbine
			key := sts.TimeStampConverted.UTC().Format("060102_1504") + "#" + sts.ProjectName + "#" + sts.Turbine
			sctm := new(ScadaConvTenMin)
			if _, exist := mapscadatenmin[key]; exist {
				sctm = mapscadatenmin[key]
			}

			fillto(sts, sctm, key)

			mapscadatenmin[key] = sctm
		}

		for _, sctm := range mapscadatenmin {
			sdata <- *sctm
		}

		for i := 0; i < len(mapscadatenmin); i++ {
			<-sresult
		}

		result <- len(mapscadatenmin)
		// tk.Printfn("[DONE] Interval %v with %d data", _tinterval, len(mapscadatenmin))
	}

	close(sdata)
	return
}

func workersave(wi int, jobs <-chan ScadaConvTenMin, result chan<- int) {
	workerconn, _ := PrepareConnection()
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

				istddev2 = tk.Div(istddev2, (icount * (icount - 1)))
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
	config := ReadConfig()

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
		defer file.Close()

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

	return ret
}

func getstep(count int) int {
	v := count / 5
	if v == 0 {
		return 1
	}
	return v
}

func getinterval(tablename string, cdate time.Time) (arrval []time.Time) {
	conn, err := PrepareConnection()
	if err != nil {
		tk.Println(err)
		return
	}
	defer conn.Close()

	pipes := []tk.M{}
	//ISODate("2016-08-21T19:10:00.000+0000")
	match := tk.M{}.Set("timestampconverted", tk.M{"$gt": cdate}).
		Set("timestampconverted", tk.M{"$lte": cdate.AddDate(0, 0, 1)})

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
			/*
				if strings.Contains(e.Error(), "Not found") {
					log.Printf("EOF")
				} else {
					log.Printf("ERRROR: %v \n", e.Error())
				}
			*/
			break
		}
		arrval = append(arrval, tkm.Get("_id", time.Time{}).(time.Time).UTC())
	}

	return
}

func fillto(_sts *ScadaThreeSecs, _sctm *ScadaConvTenMin, key string) {

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

	/*
		for _, _str := range arrvar {
			if _sctm.ID == "" {
				_tkmsctm.Set(_str, emptyValueBig)
				_tkmsctm.Set(tk.Sprintf("%s_Min", _str), emptyValueBig)
				_tkmsctm.Set(tk.Sprintf("%s_Max", _str), emptyValueBig)
				_tkmsctm.Set(tk.Sprintf("%s_StdDev", _str), emptyValueBig)
			}

			if _tkmsts.GetFloat64(_str) != emptyValueBig {
				if _tkmsctm.GetFloat64(_str) == emptyValueBig {
					_tkmsctm.Set(_str, 0)
				}

				ival := _tkmsts.GetFloat64(_str)
				ifloat64 := _tkmsctm.GetFloat64(_str) + ival

				_tkmsctm.Set(_str, ifloat64)

				_strmax := tk.Sprintf("%s_Max", _str)
				_strmin := tk.Sprintf("%s_Min", _str)

				_strdev := tk.Sprintf("%s_StdDev", _str)
				_strcount := tk.Sprintf("%s_Count", _str)

				if _tkmsctm.GetFloat64(_strmax) < ival {
					_tkmsctm.Set(_strmax, ival)
				}

				if _tkmsctm.GetFloat64(_strmin) == emptyValueBig || _tkmsctm.GetFloat64(_strmin) > ival {
					_tkmsctm.Set(_strmin, ival)
				}

				_tkmsctm.Set(_strdev, 0)

				iint := _tkmsctm.GetInt(_strcount) + 1
				_tkmsctm.Set(_strcount, iint)
			}
		}
		err := tk.Serde(_tkmsctm, _sctm, "json")
		if err != nil {
			tk.Println(err.Error())
		}
	*/

	// if _sctm.ID == "" {
	// 	_sctm.Fast_ActivePower_kW = emptyValueBig
	// 	_sctm.Fast_ActivePower_kW_StdDev = emptyValueBig
	// 	_sctm.Fast_ActivePower_kW_Min = emptyValueBig
	// 	_sctm.Fast_ActivePower_kW_Max = emptyValueBig
	// }

	// if _sts.Fast_ActivePower_kW != emptyValueBig {
	// 	if _sctm.Fast_ActivePower_kW == emptyValueBig {
	// 		_sctm.Fast_ActivePower_kW = 0
	// 	}

	// 	_sctm.Fast_ActivePower_kW += _sts.Fast_ActivePower_kW

	// 	if _sctm.Fast_ActivePower_kW_Max < _sts.Fast_ActivePower_kW {
	// 		_sctm.Fast_ActivePower_kW_Max = _sts.Fast_ActivePower_kW
	// 	}

	// 	if _sctm.Fast_ActivePower_kW_Min == emptyValueBig || _sctm.Fast_ActivePower_kW_Min > _sts.Fast_ActivePower_kW {
	// 		_sctm.Fast_ActivePower_kW_Min = _sts.Fast_ActivePower_kW
	// 	}

	// 	_sctm.Fast_ActivePower_kW_StdDev = 0
	// 	_sctm.Fast_ActivePower_kW_Count += 1
	// }

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
