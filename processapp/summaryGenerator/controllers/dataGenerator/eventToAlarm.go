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

type OemPCValue struct {
	PCValue float64
	IsCalc  bool
}

var (
	mtx         = &sync.Mutex{}
	counterRow  = 0
	projectName = "Tejuva"
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
	// countx := 0
	// xTurbines := []string{}
	_ = wg
	for turbine, _ := range ev.BaseController.RefTurbines {
		// xTurbines = append(xTurbines, turbine)
		// wg.Add(1)
		// go func(t string) {
		t := turbine
		filterX := []*dbox.Filter{}
		filterX = append(filterX, dbox.Eq("projectname", projectName))
		filterX = append(filterX, dbox.Eq("turbine", t))

		latestDate := ev.BaseController.GetLatest("Alarm", projectName, t)

		// log.Printf(">>> db.EventDown.find({turbine: \"%v\",timeend: {$gt: ISODate(%v+0000)}}).count()\n", t, latestDate.UTC().Format("2006-01-02T15:04:05.000"))

		if latestDate.Format("2006") != "0001" {
			filterX = append(filterX, dbox.Gt("timeend", latestDate.UTC()))
		}

		// ev.BaseController.Ctx.DeleteMany(new(Alarm), dbox.Gt("startdate", latestDate))

		csr, e := ctx.NewQuery().From(new(EventDown).TableName()).
			Where(dbox.And(filterX...)).Cursor(nil)

		defer csr.Close()

		countData := csr.Count()
		events := []*EventDown{}

		// do process here
		e = csr.Fetch(&events, 0, false)
		ErrorHandler(e, funcName)

		tk.Printf("Event to Alarm for %v | %v \n", t, countData)
		for _, d := range events {

			mtx.Lock()
			dataInput := d
			//tk.Printf("%s ", idx)

			ev.doConversion(dataInput)
			mtx.Unlock()
		}
		tk.Printf("end process for %v \n", t)

		csr.Close()
		// wg.Done()
		// }(turbine)

		// countx++

		// if countx%5 == 0 || (len(ev.BaseController.RefTurbines) == countx) {
		// 	for _, tx := range xTurbines {
		// 		log.Printf(">> %v \n", tx)
		// 	}
		// 	xTurbines = []string{}
		// 	wg.Wait()
		// }
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
		alarm.ReduceAvailability = event.ReduceAvailability

		timeStartWhr := event.TimeStart
		timeEndWhr := event.TimeEnd.Add(10 * time.Minute)

		allscadapcval := make(map[string]OemPCValue)

		csr2, e := ctx.Connection.NewQuery().
			Select("timestamp", "pcvalue", "denpower", "ai_intern_activpower").
			From(new(ScadaDataOEM).TableName()).
			Where(dbox.And(
				dbox.Eq("projectname", projectName),
				dbox.Eq("turbine", event.Turbine),
				dbox.Gte("timestamp", timeStartWhr),
				dbox.Lte("timestamp", timeEndWhr))).
			Order("timestamp").
			Cursor(nil)

		ErrorHandler(e, "Convert Event to Alarm")
		for {
			scadasoem := ScadaDataOEM{}
			e = csr2.Fetch(&scadasoem, 1, false)
			if e != nil {
				break
			}

			_key := scadasoem.TimeStamp.Format("20060102150405")
			_ipcval := OemPCValue{}
			_ipcval.PCValue = scadasoem.PCValue
			if scadasoem.DenPower > scadasoem.AI_intern_ActivPower {
				_ipcval.IsCalc = true
			}

			allscadapcval[_key] = _ipcval
		}
		csr2.Close()

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

		next10min := GetNext10Min(alarm.StartDate)
		for {
			detail.StartDate = next10min.Add(-10 * time.Minute)
			detail.EndDate = next10min

			if detail.StartDate.Before(alarm.StartDate) {
				detail.StartDate = alarm.StartDate
			}

			if detail.EndDate.After(alarm.EndDate) {
				detail.EndDate = alarm.EndDate
			}

			detail.DetailDateInfo = GetDateInfo(detail.StartDate)
			detail.Duration = detail.EndDate.Sub(detail.StartDate).Hours()

			_key := next10min.Format("20060102150405")
			scadapcval, _has := allscadapcval[_key]
			if !_has {
				scadapcval = ev.getPCValue(projectName, event.Turbine, next10min)
			}

			detail.Power = scadapcval.PCValue

			if scadapcval.IsCalc {
				detail.PowerLost = detail.Power * detail.Duration
				alarm.PowerLost += detail.PowerLost
			}

			alarm.Detail = append(alarm.Detail, detail)

			if next10min.After(alarm.EndDate) {
				break
			}

			next10min = next10min.Add(10 * time.Minute)
		}

		ctx.Insert(alarm)
	}
}

func (e *EventToAlarm) getPCValue(projectname, turbine string, timestamp time.Time) (pcvalue OemPCValue) {
	pcvalue = OemPCValue{}

	turbineinfo := e.BaseController.RefTurbines.Get(turbine, tk.M{}).(tk.M)
	topcorrel := turbineinfo.Get("topcorrelation", []string{}).([]string)

	avgws := getAvgWsForLostEnergy(projectname, turbine, topcorrel, 0, timestamp, e.Ctx.Connection)

	pcvalue.PCValue, _ = GetPowerCurveCubicInterpolation(e.Ctx.Connection, projectname, avgws)
	pcvalue.IsCalc = true

	return
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

func GetNext10Min(current time.Time) time.Time {
	current = current.UTC()
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
