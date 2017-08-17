package generatorControllers

import (
	. "eaciit/wfdemo-git/library/helper"
	. "eaciit/wfdemo-git/library/models"
	. "eaciit/wfdemo-git/processapp/summaryGenerator/controllers"
	"math"
	"os"
	// "sort"
	"strings"
	"sync"
	"time"

	"github.com/eaciit/dbox"
	_ "github.com/eaciit/dbox/dbc/mongo"
	tk "github.com/eaciit/toolkit"
)

// UpdateScadaMinutes
type UpdateScadaOemMinutes struct {
	*BaseController
}

type minEventDown struct {
	TimeStart          time.Time
	TimeEnd            time.Time
	Turbine            string
	ReduceAvailability bool
	DownGrid           bool
	DownEnvironment    bool
	DownMachine        bool
}

var (
	mtxOem = &sync.Mutex{}
)

func (d *UpdateScadaOemMinutes) GenerateDensity(base *BaseController) {
	sProjectName := "Tejuva"
	funcName := "UpdateScadaOemDensity Data"
	if base != nil {
		d.BaseController = base

		conn, e := PrepareConnection()
		if e != nil {
			ErrorHandler(e, funcName)
			os.Exit(0)
		}
		defer conn.Close()

		tk.Println("UpdateScadaOemDensity Data")
		var wg sync.WaitGroup

		// #faisal
		// get latest scadadata from scadadata and put condition to get the scadadataoem based on latest scadadata
		// put some match condition here
		// tk.Println(d.BaseController.RefTurbines)
		for turbine, _ := range d.BaseController.RefTurbines {
			filter := []*dbox.Filter{}
			filter = append(filter, dbox.Eq("projectname", sProjectName))
			filter = append(filter, dbox.Eq("turbine", turbine))

			latestDate := d.BaseController.GetLatest("ScadaData", sProjectName, turbine)
			if latestDate.Format("2006") != "0001" {
				filter = append(filter, dbox.Gt("timestamp", latestDate))
			}

			// filter = append(filter, dbox.Gt("timestamp", d.BaseController.LatestData.MapScadaData["Tejuva#"+turbine]))

			csr, e := conn.NewQuery().From(new(ScadaDataOEM).TableName()).
				Where(filter...).
				Order("timestamp").
				Cursor(nil)
			ErrorHandler(e, funcName)

			defer csr.Close()

			counter := 0
			isDone := false
			countPerProcess := 1000
			countData := csr.Count()

			tk.Printf("\nDensity for %v | %v \n", turbine, countData)

			for !isDone && countData > 0 {
				scadas := []*ScadaDataOEM{}
				e = csr.Fetch(&scadas, countPerProcess, false)

				if len(scadas) < countPerProcess {
					isDone = true
				}

				wg.Add(1)
				go func(datas []*ScadaDataOEM, endIndex int) {
					tk.Printf("Starting process %v data\n", endIndex)
					defer wg.Done()

					lenDatas := len(datas)
					if lenDatas == 0 {
						tk.Println("Data is 0 return process")
						return
					}

					// == Get Data Downtime ==
					_filter := []*dbox.Filter{}
					_filter = append(_filter, dbox.Eq("projectname", sProjectName))
					_filter = append(_filter, dbox.Eq("turbine", turbine))
					_filter = append(_filter, dbox.Or(dbox.And(dbox.Gte("timestart", datas[0].TimeStamp.Add(time.Minute*-10)), dbox.Lte("timestart", datas[lenDatas-1].TimeStamp)),
						dbox.And(dbox.Gte("timeend", datas[0].TimeStamp.Add(time.Minute*-10)), dbox.Lte("timeend", datas[lenDatas-1].TimeStamp))))

					_csr, _err := conn.NewQuery().
						Select("timestart", "timeend", "turbine", "reduceavailability", "downgrid", "downenvironment", "downmachine").
						From(new(EventDown).TableName()).
						Where(_filter...).
						Order("timestamp").
						Cursor(nil)

					_arrMED := []minEventDown{}
					if _err == nil {
						_ = _csr.Fetch(&_arrMED, 0, false)
						_csr.Close()
					}
					// == Get Data Downtime ==

					mtxOem.Lock()
					logStart := time.Now()

					for _, data := range datas {
						d.updateScadaOEM(data, _arrMED)
					}

					logDuration := time.Now().Sub(logStart)
					mtxOem.Unlock()

					tk.Printf("End processing for %v data about %v sec(s)\n", endIndex, logDuration.Seconds())
				}(scadas, ((counter + 1) * countPerProcess))

				counter++
				if counter%10 == 0 || isDone {
					wg.Wait()
				}
			}
		}
	}
}

func (u *UpdateScadaOemMinutes) updateScadaOEM(data *ScadaDataOEM, arrMED []minEventDown) {
	ctx := u.Ctx
	turbine := GetExactTurbineId(strings.TrimSpace(data.Turbine))

	turbines := u.BaseController.RefTurbines.Get(turbine)

	turbineinfo := u.BaseController.RefTurbines.Get(turbine, tk.M{}).(tk.M)
	topcorrel := turbineinfo.Get("topcorrelation", []string{}).([]string)

	power := data.AI_intern_ActivPower
	energy := tk.Div(power, 6)

	energyLost := 0.0 // tk.Div(data.PowerLost, 6)

	pH := 0.0
	elevation := 0.0
	temperature := data.Temp_Outdoor
	if turbines != nil {
		if turbines.(tk.M).Has("turbineelevation") {
			elevation = turbines.(tk.M).GetFloat64("turbineelevation")
			exponen := (-(9.80665) * 28.9644 * elevation) / (8314.32 * 288.15)
			pH = 101325 * math.Exp(exponen)
		}
	}

	denWs := 0.0
	denPower := 0.0
	adjDenWs := 0.0

	density := pH / (287.05 * (273.15 + temperature))

	// avgWs := data.AI_intern_WindSpeed
	avgWs := getAvgWsForLostEnergy("Tejuva", turbine, topcorrel, data.AI_intern_WindSpeed, data.TimeStamp, ctx.Connection)
	adjWs := tk.RoundingDown64(avgWs, 1)
	// tk.Println(data.AI_intern_WindSpeed, " - ", avgWs, " - ", adjWs)
	pcValue, _ := GetPowerCurveCubicInterpolation(ctx.Connection, "Tejuva", avgWs)
	pcValueAdj, _ := GetPowerCurveCubicInterpolation(ctx.Connection, "Tejuva", adjWs)
	pcDeviation := pcValue - power

	denWs = avgWs * math.Pow((density/1.225), (1.0/3.0))
	adjDenWs = tk.RoundingDown64(denWs, 0)
	if denWs < 3.75 && denWs >= 3.5 {
		adjDenWs = 3.5
	}
	denPower, _ = GetPowerCurveCubicInterpolation(ctx.Connection, "Tejuva", denWs)

	denPcValue, _ := GetPowerCurveCubicInterpolation(ctx.Connection, "Tejuva", denWs)
	denPcDeviation := denPcValue - denPower

	deviationPct := 0.0
	denDeviationPct := 0.0
	if pcDeviation > 0 {
		deviationPct = tk.Div(pcDeviation, pcValue)
	}
	if denPcDeviation > 0 {
		denDeviationPct = tk.Div(denPcDeviation, denPcValue)
	}

	oktime := 600.0
	// machinedown := 0.0
	// griddown := 0.0

	timestamp := data.TimeStamp
	timestamp0 := timestamp.Add(-10 * time.Minute)

	totalDurationMttf := 0.0
	totalDowntime := 0
	aDuration := 0.0

	gridDowntime := 0.0
	machineDowntime := 0.0
	unknownDowntime := 0.0

	gridDowntimeAll := 0.0
	machineDowntimeAll := 0.0
	unknownDowntimeAll := 0.0

	// getting alarms from min EventDown
	selArrMed := []minEventDown{}
	if len(arrMED) > 0 {
		for _, a := range arrMED {
			if a.TimeStart.Sub(timestamp0) >= 0 && a.TimeEnd.Sub(timestamp) <= 0 {
				selArrMed = append(selArrMed, a)
			} else if a.TimeStart.Sub(timestamp0) >= 0 && a.TimeStart.Sub(timestamp) <= 0 {
				selArrMed = append(selArrMed, a)
			} else if a.TimeEnd.Sub(timestamp0) >= 0 && a.TimeEnd.Sub(timestamp) <= 0 {
				selArrMed = append(selArrMed, a)
			} else if a.TimeStart.Sub(timestamp0) <= 0 && a.TimeEnd.Sub(timestamp) >= 0 {
				selArrMed = append(selArrMed, a)
			}
		}
	}

	// log.Printf("%v -> scada: %v - %v \n", turbine, timestamp0.UTC().String(), timestamp.UTC().String())
	if len(selArrMed) > 0 {
		for _, a := range selArrMed {
			// log.Printf("alarm: %v - %v \n", a.TimeStart.UTC().String(), a.TimeEnd.UTC().String())

			startTime := a.TimeStart
			endTime := a.TimeEnd

			if timestamp0.Sub(startTime) > 0 {
				startTime = timestamp0
			}
			if timestamp.Sub(endTime) < 0 {
				endTime = timestamp
			}

			/*
				startTime := timestamp0
				endTime := timestamp
				if startTime.Sub(a.TimeStart) > 0 {
					startTime = a.TimeStart
				}
				if endTime.Sub(a.TimeEnd) > 0 {
					endTime = a.TimeEnd
				}
				aDuration += endTime.Sub(startTime).Seconds()
			*/

			//@ASP 29-07-2017 : Reduce availability only for machine down
			if a.ReduceAvailability || !a.DownMachine {
				aDuration += endTime.Sub(startTime).Seconds()

				if a.DownGrid {
					gridDowntime += endTime.Sub(startTime).Seconds()
				} else if a.DownMachine {
					machineDowntime += endTime.Sub(startTime).Seconds()
				} else if a.DownEnvironment {
					unknownDowntime += endTime.Sub(startTime).Seconds()
				}

				totalDowntime++
			}

			if a.DownGrid {
				gridDowntimeAll += endTime.Sub(startTime).Seconds()
			} else if a.DownMachine {
				machineDowntimeAll += endTime.Sub(startTime).Seconds()
			} else if a.DownEnvironment {
				unknownDowntimeAll += endTime.Sub(startTime).Seconds()
			}
			/*log.Printf("endTime: %v | startTime: %v \n", endTime.UTC().String(), startTime.UTC().String())
			log.Printf("aDuration: %v | machineDowntime: %v | gridDowntime: %v | unknownDowntime: %v \n", aDuration, machineDowntime, gridDowntime, unknownDowntime)*/
		}
		if aDuration > 600 {
			aDuration = 600
		}

		if machineDowntime > 600 {
			machineDowntime = 600
		}

		if gridDowntime > 600 {
			gridDowntime = 600
		}

		if unknownDowntime > 600 {
			unknownDowntime = 600
		}

		if machineDowntimeAll > 600 {
			machineDowntimeAll = 600
		}

		if gridDowntimeAll > 600 {
			gridDowntimeAll = 600
		}

		if unknownDowntimeAll > 600 {
			unknownDowntimeAll = 600
		}

		totalDurationMttf = tk.Div(aDuration, float64(totalDowntime))
	}

	// set mttr & mttf
	mttr := totalDurationMttf
	mttf := oktime - aDuration
	mtbf := tk.Div(totalDurationMttf, float64(totalDowntime))

	if denPower > data.AI_intern_ActivPower {
		energyLost = tk.Div(aDuration, 3600) * pcValue
	}

	totalavail := tk.Div(oktime, 600.0)
	machineavail := tk.Div((600.0 - machineDowntime), 600.0)
	gridavail := tk.Div((600.0 - gridDowntime), 600.0)

	totalavailall := tk.Div(oktime, 600.0)
	machineavailall := tk.Div((600.0 - machineDowntimeAll), 600.0)
	gridavailall := tk.Div((600.0 - gridDowntimeAll), 600.0)

	powerLost := denPower - data.AI_intern_ActivPower

	perfIndex := 0.0
	if denPower > 0 && data.AI_intern_ActivPower > 0 {
		perfIndex = tk.Div(data.AI_intern_ActivPower, denPower)
	}

	retadjws := tk.RoundingAuto64(data.AI_intern_WindSpeed, 0)
	retavgws := tk.RoundingAuto64(data.AI_intern_WindSpeed, 0)
	//for PC
	// retadjws := tk.RoundingAuto64(data.AI_intern_WindSpeed, 0)
	// retavgws := data.AI_intern_WindSpeed

	e := ctx.Connection.NewQuery().Update().From(new(ScadaDataOEM).TableName()).
		Where(dbox.Eq("_id", data.ID)).
		Exec(tk.M{}.Set("data", tk.M{}.
			Set("turbine", turbine).
			Set("totalavail", totalavail).
			Set("machineavail", machineavail).
			Set("gridavail", gridavail).
			Set("totalavailall", totalavailall).
			Set("machineavailall", machineavailall).
			Set("gridavailall", gridavailall).
			Set("turbineelevation", elevation).
			Set("wsadjforpc", retadjws).
			Set("wsavgforpc", retavgws).
			Set("pcdeviation", pcDeviation).
			Set("pcvalue", pcValue).
			Set("pcvalueadj", pcValueAdj).
			Set("powerlost", powerLost).
			Set("energylost", energyLost).
			Set("energy", energy).
			Set("denvalue", density).
			Set("denph", pH).
			Set("denadjwindspeed", adjDenWs).
			Set("denwindspeed", denWs).
			Set("denpower", denPower).
			Set("denpcdeviation", denPcDeviation).
			Set("dendeviationpct", denDeviationPct).
			Set("denpcvalue", denPcValue).
			Set("deviationpct", deviationPct).
			Set("mttr", mttr).
			Set("mttf", mttf).
			Set("mtbf", mtbf).
			Set("performanceindex", perfIndex).
			Set("griddowntime", gridDowntime).
			Set("machinedowntime", machineDowntime).
			Set("unknowndowntime", unknownDowntime).
			Set("griddowntimeall", gridDowntimeAll).
			Set("machinedowntimeall", machineDowntimeAll).
			Set("unknowndowntimeall", unknownDowntimeAll)))

	if e != nil {
		tk.Printf("Update fail: %s", e.Error())
	}
}

//================================================
//=== GET AvgWindSpeed for Lost Energy Calculation

func getAvgWsForLostEnergy(project, turbine string, topcorrel []string, oemavgws float64, itime time.Time, ctx dbox.IConnection) (iavgws float64) {
	iavgws = 0

	if oemavgws >= 0 && oemavgws <= 50 {
		iavgws = oemavgws
		return
	}

	itopcorrel := []interface{}{}
	for _, str := range topcorrel {
		itopcorrel = append(itopcorrel, str)
	}

	csroem, err := ctx.NewQuery().
		Select("timestamp", "turbine", "ai_intern_windspeed").
		From(new(ScadaDataOEM).TableName()).
		Where(dbox.And(dbox.Eq("projectname", project),
			dbox.Eq("timestamp", itime),
			dbox.In("turbine", itopcorrel...),
			dbox.Lte("ai_intern_windspeed", 50),
			dbox.Gte("ai_intern_windspeed", 0))).
		Cursor(nil)
	if err != nil {
		return
	}
	defer csroem.Close()

	allres := tk.M{}
	for {
		trx := new(ScadaDataOEM)
		e := csroem.Fetch(trx, 1, false)
		if e != nil {
			break
		}

		if trx.AI_intern_WindSpeed > 0 && trx.AI_intern_WindSpeed < 50 {
			allres.Set(trx.Turbine, trx.AI_intern_WindSpeed)
		}
	}

	for _, key := range topcorrel {
		if allres.Has(key) {
			iavgws = allres.GetFloat64(key)
			break
		}
	}

	if iavgws > 0 && iavgws <= 50 {
		return
	}

	csr, err := ctx.NewQuery().
		Select("vhubws90mavg").
		From("MetTower").
		Where(dbox.Eq("timestamp", itime)).
		Cursor(nil)
	if err == nil && csr.Count() > 0 {
		_tkm := tk.M{}
		err = csr.Fetch(&_tkm, 1, false)
		if err == nil {
			iavgws = _tkm.GetFloat64("vhubws90mavg")
		}
	}

	if err == nil {
		csr.Close()
	}

	return
}

func (d *UpdateScadaOemMinutes) UpdateDeviation(base *BaseController) {
	d.BaseController = base
	sProjectName := "Tejuva"
	funcName := "Update Deviation Data"
	if base != nil {
		d.BaseController = base

		conn, e := PrepareConnection()
		if e != nil {
			ErrorHandler(e, funcName)
			os.Exit(0)
		}
		defer conn.Close()

		tk.Println("Update Deviation Data")
		var wg sync.WaitGroup

		// #faisal
		// get latest scadadata from scadadata and put condition to get the scadadataoem based on latest scadadata
		// put some match condition here
		// tk.Println(d.BaseController.RefTurbines)
		for turbine, _ := range d.BaseController.RefTurbines {
			filter := []*dbox.Filter{}
			filter = append(filter, dbox.Eq("projectname", sProjectName))
			filter = append(filter, dbox.Eq("turbine", turbine))

			csr, e := conn.NewQuery().From(new(ScadaData).TableName()).
				Where(filter...).
				Order("timestamp").
				Cursor(nil)
			ErrorHandler(e, funcName)

			defer csr.Close()

			counter := 0
			isDone := false
			countPerProcess := 1000
			countData := csr.Count()

			tk.Printf("\nDeviation for %v | %v \n", turbine, countData)

			for !isDone && countData > 0 {
				scadas := []*ScadaData{}
				e = csr.Fetch(&scadas, countPerProcess, false)

				if len(scadas) < countPerProcess {
					isDone = true
				}

				wg.Add(1)
				go func(datas []*ScadaData, endIndex int) {
					tk.Printf("Starting process %v data\n", endIndex)
					defer wg.Done()

					lenDatas := len(datas)
					if lenDatas == 0 {
						tk.Println("Data is 0 return process")
						return
					}

					logStart := time.Now()

					for _, data := range datas {

						// mtxOem.Lock()
						// defer mtxOem.Unlock()

						denpower := data.DenPower
						power := data.Power
						deviation := denpower - power
						deviationpct := math.Abs(tk.Div(deviation, denpower))

						// tk.Printf("%v %v %v %v %v\n", data.ID, denpower, power, deviation, deviationpct)

						e := d.Ctx.Connection.NewQuery().Update().From(new(ScadaData).TableName()).
							Where(dbox.Eq("_id", data.ID)).
							Exec(tk.M{}.Set("data", tk.M{}.
								Set("denpcdeviation", deviation).
								Set("dendeviationpct", deviationpct)))

						if e != nil {
							tk.Printf("Update fail: %s", e.Error())
						}
					}

					logDuration := time.Now().Sub(logStart)

					tk.Printf("End processing for %v data about %v sec(s)\n", endIndex, logDuration.Seconds())
				}(scadas, ((counter + 1) * countPerProcess))

				counter++
				if counter%10 == 0 || isDone {
					wg.Wait()
				}
			}
		}
	}
}

/* First Logic
func (u *UpdateScadaOemMinutes) getAvgWsForLostEnergy(project, turbine string, oemavgws float64, itime time.Time) (iavgws float64) {
	iavgws = 0

	if oemavgws >= 0 && oemavgws <= 50 {
		iavgws = oemavgws
		return
	}

	csr, err := u.Ctx.Connection.NewQuery().
		Select("vhubws90mavg").
		From("MetTower").
		Where(dbox.Eq("timestamp", itime)).
		Cursor(nil)
	if err == nil && csr.Count() > 0 {
		_tkm := tk.M{}
		err = csr.Fetch(&_tkm, 1, false)
		if err == nil {
			iavgws = _tkm.GetFloat64("vhubws90mavg")
		}
	}

	if err == nil {
		csr.Close()
	}

	if iavgws > 0 && iavgws <= 50 {
		return
	}

	//GetCorrelWs
	csroem, err := u.Ctx.Connection.NewQuery().
		Select("timestamp", "turbine", "ai_intern_windspeed").
		From(new(ScadaDataOEM).TableName()).
		Where(dbox.And(dbox.Eq("projectname", project), dbox.Lte("timestamp", itime), dbox.Gte("timestamp", itime.AddDate(0, 0, -1)))).
		Cursor(nil)
	if err != nil {
		return
	}
	defer csroem.Close()

	allres, _tturbine, bcorrelval := tk.M{}, tk.M{}, float64(-2)
	ikey := itime.Format("20060102150405")

	for {
		trx := new(ScadaDataOEM)
		e := csroem.Fetch(trx, 1, false)
		if e != nil {
			break
		}

		dkey := trx.TimeStamp.Format("20060102150405")

		_tkm := allres.Get(trx.Turbine, tk.M{}).(tk.M)
		if trx.AI_intern_WindSpeed != -99999.0 {
			_tkm.Set(dkey, trx.AI_intern_WindSpeed)
		}

		allres.Set(trx.Turbine, _tkm)
		_tturbine.Set(trx.Turbine, 1)
	}

	for _turbine, _ := range _tturbine {
		if turbine == _turbine {
			continue
		}

		_dt01 := allres.Get(turbine, tk.M{}).(tk.M)
		_dt02 := allres.Get(_turbine, tk.M{}).(tk.M)

		if len(_dt01) > 0 && len(_dt02) > 0 && _dt02.Has(ikey) {
			_icorrel := GetCorrelation(_dt01, _dt02)
			if _icorrel > bcorrelval {
				iavgws = _dt02.GetFloat64(ikey)
			}
		}
	}

	return
}
*/

//================================================
//================================================
