package generatorControllers

import (
	. "eaciit/wfdemo-git/library/helper"
	. "eaciit/wfdemo-git/library/models"
	. "eaciit/wfdemo-git/processapp/summaryGenerator/controllers"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/eaciit/dbox"
	_ "github.com/eaciit/dbox/dbc/mongo"
	tk "github.com/eaciit/toolkit"
)

type EventToAlarm struct {
	*BaseController
}

var (
	mtx        = &sync.Mutex{}
	counterRow = 0
)

func (ev *EventToAlarm) ConvertEventToAlarm(base *BaseController) {
	ev.BaseController = base
	tk.Println("Start process converting Event to Alarm...")

	funcName := "EventToAlarmConversion"
	var wg sync.WaitGroup

	ctx, e := PrepareConnection()
	if e != nil {
		ErrorHandler(e, funcName)
		os.Exit(0)
	}

	// #faisal
	// add condition to get the eventdown started from the latest data that already in alarm, so no need to generate the alarm data from begining
	// remove delete function
	countx := 0
	for turbine, _ := range ev.BaseController.RefTurbines {
		filter := []*dbox.Filter{}
		filter = append(filter, dbox.Eq("projectname", "Tejuva"))
		filter = append(filter, dbox.Eq("turbine", turbine))
		filter = append(filter, dbox.Gt("timeend", ev.BaseController.LatestData.MapAlarm["Tejuva#"+turbine]))

		// ev.BaseController.Ctx.DeleteMany(new(Alarm), dbox.Gt("startdate", ev.BaseController.LatestData.MapAlarm["Tejuva#"+turbine]))

		csr, e := ctx.NewQuery().From(new(EventDown).TableName()).
			Where(filter...).Cursor(nil)

		defer csr.Close()

		// counter := 0
		countData := csr.Count()
		isDone := false
		// countPerProcess := 5

		// for !isDone && countData > 0 {
		events := []*EventDown{}

		// do process here
		e = csr.Fetch(&events, 0, false)
		ErrorHandler(e, funcName)

		wg.Add(1)
		go func(datas []*EventDown, count int, t string) {
			tk.Printf("Event to Alarm for %v | %v \n", t, count)
			for _, d := range datas {

				mtx.Lock()
				dataInput := d
				//tk.Printf("%s ", idx)

				ev.doConversion(dataInput)
				mtx.Unlock()
			}
			tk.Printf("end process for %v \n", t)

			wg.Done()
		}(events, countData, turbine)

		countx++

		if len(ev.BaseController.RefTurbines) == countx {
			isDone = true
		}

		if countx%5 == 0 || isDone {
			wg.Wait()
		}
		// }
	}

	tk.Println("End process converting Event to Alarm...")
}

func (ev *EventToAlarm) doConversion(event *EventDown) {
	if event.Turbine != "" {
		ctx := ev.Ctx
		counterRow++

		turbine := GetExactTurbineId(strings.TrimSpace(event.Turbine))

		alarm := new(Alarm).New()
		alarm.StartDate = event.TimeStart
		alarm.EndDate = event.TimeEnd
		alarm.Duration = tk.RoundingAuto64(tk.Div(event.Duration, 3600.0), 2)
		alarm.StartDateInfo = event.DateInfoStart
		alarm.Turbine = turbine
		alarm.ProjectName = event.ProjectName
		alarm.Farm = event.ProjectName
		alarm.AlertDescription = event.AlarmDescription
		alarm.BrakeType = event.BrakeType // add by ams, regarding to add new req | 20170130
		alarm.ExternalStop = false
		alarm.GridDown = event.DownGrid
		alarm.InternalGrid = false
		alarm.MachineDown = event.DownMachine
		alarm.Unknown = event.DownEnvironment
		alarm.AEbOK = false
		alarm.WeatherStop = false
		alarm.Line = counterRow

		timeStartWhr := event.TimeStart.Add(-10 * time.Minute)
		timeEndWhr := event.TimeEnd.Add(10 * time.Minute)

		scadas := make([]ScadaDataOEM, 0)
		csr2, e := ctx.Connection.NewQuery().From(new(ScadaDataOEM).TableName()).
			Where(dbox.And(
				dbox.Eq("projectname", "Tejuva"),
				dbox.Eq("turbine", event.Turbine),
				dbox.Gte("timestamputc", timeStartWhr),
				dbox.Lte("timestamputc", timeEndWhr))).
			Order("timestamputc").
			Cursor(nil)

		e = csr2.Fetch(&scadas, 0, false)
		ErrorHandler(e, "Convert Event to Alarm")
		csr2.Close()

		// details := make([]*AlarmDetail, 0)

		// currMonthId := 0
		// detail := new(AlarmDetail)

		// if len(scadas) > 0 {
		// 	totalPower := 0.0
		// 	durationTS := 0.0
		// 	for _, scada := range scadas {
		// 		//if (event.TimeStart.Sub(scada.TimeStamp) <= 0 && event.TimeEnd.Sub(scada.TimeStamp) <= 0) || (event.TimeStart.Sub(scada.TimeStamp) <= 0 && event.TimeEnd.Sub(scada.TimeStamp) >= 0) {
		// 		if scada.DenPower > scada.AI_intern_ActivPower {
		// 			power, err := GetPowerCurveCubicInterpolation(ctx.Connection, "Tejuva", scada.AI_intern_WindSpeed)
		// 			if err != nil {
		// 				power = 0.0
		// 			}

		// 			if currMonthId != scada.DateInfo.MonthId {
		// 				if detail.AlertDescription != "" {
		// 					details = append(details, detail)
		// 				}

		// 				detail = new(AlarmDetail)

		// 				startTime := scada.TimeStamp
		// 				strDate0 := tk.Sprintf("%v-%v-%v %v:%v:%v", scada.DateInfo.Year, int(scada.DateInfo.DateId.Month()), "01", "00", "00", "00")
		// 				//tk.Println(strDate0)
		// 				startDate, _ := time.Parse("2006-1-02 15:04:05", strDate0)
		// 				//tk.Println(startDate)

		// 				if currMonthId > 0 {
		// 					startTime = startDate
		// 				}

		// 				detail.StartDate = startTime

		// 				detail.AlertDescription = event.AlarmDescription
		// 				detail.BrakeType = event.BrakeType // add by ams, regarding to add new req | 20170130
		// 				detail.AEbOK = false
		// 				detail.ExternalStop = false
		// 				detail.InternalGrid = false
		// 				detail.WeatherStop = false
		// 				detail.GridDown = event.DownGrid
		// 				detail.MachineDown = event.DownMachine
		// 				detail.Unknown = event.DownEnvironment
		// 				detail.Power = 0.0
		// 				detail.Duration = 0.0
		// 				detail.PowerLost = 0.0

		// 				currMonthId = scada.DateInfo.MonthId
		// 			}

		// 			detail.EndDate = scada.TimeStamp
		// 			// tk.Println(detail.EndDate)

		// 			// lastDateNo := daysIn(detail.EndDate.UTC().Month(), detail.EndDate.UTC().Year())
		// 			// strDate := tk.Sprintf("%v-%v-%v %v:%v:%v", detail.EndDate.UTC().Year(), int(detail.EndDate.UTC().Month()), lastDateNo, 23, 59, 59)
		// 			// lastDate, _ := time.Parse("2006-1-2 15:04:05", strDate)

		// 			durationTS = tk.Div(10.0, 60.0)

		// 			if scada.TimeStamp.Sub(event.TimeStart) > 0 && scada.TimeStamp.Add(-10*time.Minute).Sub(event.TimeStart) <= 0 {
		// 				durationTS = scada.TimeStamp.Sub(event.TimeStart).Hours()
		// 				detail.StartDate = event.TimeStart
		// 				// tk.Println("masuk kondisi 2")
		// 			} else if scada.TimeStamp.Sub(event.TimeEnd) > 0 && scada.TimeStamp.Add(-10*time.Minute).Sub(event.TimeEnd) <= 0 {
		// 				durationTS = event.TimeEnd.Sub(scada.TimeStamp.Add(-10 * time.Minute)).Hours()
		// 				detail.EndDate = event.TimeEnd
		// 				// tk.Println("masuk kondisi 3")
		// 			} else if scada.TimeStamp.Sub(event.TimeStart) > 0 && scada.TimeStamp.Add(-10*time.Minute).Sub(event.TimeStart) <= 0 && scada.TimeStamp.Sub(event.TimeEnd) > 0 && scada.TimeStamp.Add(-10*time.Minute).Sub(event.TimeEnd) <= 0 {
		// 				durationTS = event.TimeEnd.Sub(event.TimeStart).Hours()
		// 				// tk.Println("masuk kondisi 4")
		// 			} else if detail.EndDate.Sub(event.TimeEnd) >= 0 && detail.EndDate.Add(-10*time.Minute).Sub(event.TimeEnd) <= 0 {
		// 				detail.EndDate = event.TimeEnd
		// 				durationTS = event.TimeEnd.Sub(scada.TimeStamp.Add(-10 * time.Minute)).Hours()
		// 				// tk.Println("masuk kondisi 1")
		// 			}

		// 			/*if durationTS < 0 {
		// 				tk.Println(durationTS)
		// 				tk.Println(event.TimeStart)
		// 				tk.Println(event.TimeEnd)
		// 				tk.Println(event.Turbine)
		// 			}
		// 			*/
		// 			detail.Duration += durationTS
		// 			detail.PowerLost += (power * durationTS)
		// 			detail.Power += power

		// 			//tk.Println(idx, detail)

		// 			totalPower += power
		// 		}
		// 		//}
		// 	}
		// 	details = append(details, detail)

		// 	powerLost := 0.0
		// 	newDuration := 0.0

		// 	detailResults := make([]AlarmDetail, 0)
		// 	for _, dt := range details {
		// 		var dtl AlarmDetail
		// 		dtl.AEbOK = dt.AEbOK
		// 		dtl.AlertDescription = dt.AlertDescription
		// 		dtl.BrakeType = dt.BrakeType // add by ams, regarding to add new req | 20170130
		// 		dtl.DetailDateInfo = GetDateInfo(dt.StartDate)
		// 		dtl.StartDate = dt.StartDate
		// 		dtl.EndDate = dt.EndDate
		// 		dtl.Duration = dt.Duration
		// 		dtl.Power = dt.Power
		// 		dtl.PowerLost = dt.PowerLost
		// 		dtl.ExternalStop = dt.ExternalStop
		// 		dtl.GridDown = dt.GridDown
		// 		dtl.InternalGrid = dt.InternalGrid
		// 		dtl.MachineDown = dt.MachineDown
		// 		dtl.Unknown = dt.Unknown
		// 		dtl.WeatherStop = dt.WeatherStop

		// 		powerLost += dt.PowerLost
		// 		newDuration += dt.Duration
		// 		// tk.Println(dtl)

		// 		detailResults = append(detailResults, dtl)
		// 	}

		// 	//alarm.Duration = newDuration
		// 	alarm.PowerLost = powerLost
		// 	alarm.Detail = detailResults
		// } else {
		// 	detail := AlarmDetail{}
		// 	detail.StartDate = alarm.StartDate
		// 	detail.DetailDateInfo = GetDateInfo(alarm.StartDate)
		// 	detail.EndDate = alarm.EndDate
		// 	detail.Duration = alarm.Duration
		// 	detail.AlertDescription = alarm.AlertDescription
		// 	detail.BrakeType = alarm.BrakeType // add by ams, regarding to add new req | 20170130
		// 	detail.ExternalStop = alarm.ExternalStop
		// 	// detail.Power = alarm.Power
		// 	detail.PowerLost = alarm.PowerLost
		// 	detail.GridDown = alarm.GridDown
		// 	detail.InternalGrid = alarm.InternalGrid
		// 	detail.MachineDown = alarm.MachineDown
		// 	detail.AEbOK = alarm.AEbOK
		// 	detail.Unknown = alarm.Unknown
		// 	detail.WeatherStop = alarm.WeatherStop

		// 	alarm.Detail = append(alarm.Detail, detail)
		// }

		detail := AlarmDetail{}
		detail.StartDate = alarm.StartDate
		detail.DetailDateInfo = GetDateInfo(alarm.StartDate)
		detail.EndDate = alarm.EndDate
		detail.Duration = alarm.Duration
		detail.AlertDescription = alarm.AlertDescription
		detail.BrakeType = alarm.BrakeType // add by ams, regarding to add new req | 20170130
		detail.ExternalStop = alarm.ExternalStop
		// detail.Power = alarm.Power
		detail.PowerLost = alarm.PowerLost
		detail.GridDown = alarm.GridDown
		detail.InternalGrid = alarm.InternalGrid
		detail.MachineDown = alarm.MachineDown
		detail.AEbOK = alarm.AEbOK
		detail.Unknown = alarm.Unknown
		detail.WeatherStop = alarm.WeatherStop

		if len(scadas) > 0 {
			for _, scada := range scadas {
				if scada.DenPower > scada.AI_intern_ActivPower {
					power, err := GetPowerCurveCubicInterpolation(ctx.Connection, "Tejuva", scada.AI_intern_WindSpeed)
					if err != nil {
						power = 0.0
					}
					detail.Power = power
					alarm.PowerLost = power * alarm.Duration
					detail.PowerLost = alarm.PowerLost
				}

				if alarm.PowerLost != 0.0 {
					break
				}
			}

		}

		alarm.Detail = append(alarm.Detail, detail)

		ctx.Insert(alarm)
	}
}

func GetExactTurbineId(tId string) string {
	turbine := tId
	numTurbine := 0
	if strings.Contains(turbine, "HBR") && len(turbine) < 6 {
		numTurbine = tk.ToInt(strings.Replace(turbine, "HBR", "", 1), "0")
		nolnya := ""
		for i := 0; i < (3 - len(tk.ToString(numTurbine))); i++ {
			nolnya += "0"
		}
		turbine = "HBR" + nolnya + tk.ToString(numTurbine)
	} else if strings.Contains(turbine, "SSE") && len(turbine) < 6 {
		numTurbine = tk.ToInt(strings.Replace(turbine, "SSE", "", 1), "0")
		nolnya := ""
		for i := 0; i < (3 - len(tk.ToString(numTurbine))); i++ {
			nolnya += "0"
		}
		turbine = "SSE" + nolnya + tk.ToString(numTurbine)
	} else if strings.Contains(turbine, "TJW") && len(turbine) < 6 {
		numTurbine = tk.ToInt(strings.Replace(turbine, "TJW", "", 1), "0")
		nolnya := ""
		for i := 0; i < (3 - len(tk.ToString(numTurbine))); i++ {
			nolnya += "0"
		}
		turbine = "TJW" + nolnya + tk.ToString(numTurbine)
	} else if strings.Contains(turbine, "TJ") && len(turbine) < 5 {
		numTurbine = tk.ToInt(strings.Replace(turbine, "TJ", "", 1), "0")
		nolnya := ""
		for i := 0; i < (3 - len(tk.ToString(numTurbine))); i++ {
			nolnya += "0"
		}
		turbine = "TJ" + nolnya + tk.ToString(numTurbine)
	}

	return turbine
}