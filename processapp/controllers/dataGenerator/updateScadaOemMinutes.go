package generatorControllers

import (
	. "eaciit/wfdemo-git/library/helper"
	. "eaciit/wfdemo-git/library/models"
	. "eaciit/wfdemo-git/processapp/controllers"
	"github.com/eaciit/dbox"
	_ "github.com/eaciit/dbox/dbc/mongo"
	tk "github.com/eaciit/toolkit"
	"math"
	"os"
	"strings"
	"sync"
	"time"
)

// UpdateScadaMinutes
type UpdateScadaOemMinutes struct {
	*BaseController
}

var (
	mtxOem      = &sync.Mutex{}
	refTurbines = tk.M{}
	refAlarms   = tk.M{}
)

func (u *UpdateScadaOemMinutes) prepareDataReff() {
	tk.Println("Getting data refference")
	logStart := time.Now()

	turbines := []TurbineMaster{}
	csrt, e := u.Ctx.Connection.NewQuery().From(new(TurbineMaster).TableName()).Cursor(nil)

	e = csrt.Fetch(&turbines, 0, false)
	ErrorHandler(e, funcName)
	csrt.Close()

	refTurbines = tk.M{}
	for _, t := range turbines {
		refTurbines.Set(t.TurbineId, tk.M{}.
			Set("turbinename", t.TurbineName).
			Set("turbineelevation", t.Elevation))
	}

	alarms := []DowntimeEvent{}
	csr2, e := u.Ctx.Connection.NewQuery().From(new(DowntimeEvent).TableName()).
		Where(dbox.Eq("projectname", "Tejuva")).Cursor(nil)

	e = csr2.Fetch(&alarms, 0, false)
	ErrorHandler(e, "get alarm data")
	csr2.Close()

	refAlarms = tk.M{}
	details := []DowntimeEvent{}
	for _, a := range alarms {
		if refAlarms.Has(a.Turbine) {
			details = refAlarms.Get(a.Turbine).([]DowntimeEvent)
		} else {
			details = []DowntimeEvent{}
		}

		details = append(details, a)
		refAlarms.Set(a.Turbine, details)
	}

	logDurationg := time.Now().Sub(logStart).Seconds()
	tk.Printf("Getting refference data about %v secs\n", logDurationg)
}

func (d *UpdateScadaOemMinutes) GenerateDensity(base *BaseController) {
	funcName := "UpdateScadaOemDensity Data"
	if base != nil {
		d.BaseController = base

		conn, e := PrepareConnection()
		if e != nil {
			ErrorHandler(e, funcName)
			os.Exit(0)
		}

		tk.Println("UpdateScadaOemDensity Data")

		d.prepareDataReff()

		var wg sync.WaitGroup

		// tglProses, _ := time.Parse("2006-01-02 15:04:05", "2016-10-28 00:00:00")

		// csr, e := conn.NewQuery().From(new(ScadaDataOEM).TableName()).
		// 	Where(dbox.And(dbox.Eq("projectname", "Tejuva"), dbox.Eq("dateinfo.dateid", tglProses))).Cursor(nil)
		// ErrorHandler(e, funcName)

		csr, e := conn.NewQuery().From(new(ScadaDataOEM).TableName()).
			Where(dbox.Eq("projectname", "Tejuva")).Cursor(nil)
		ErrorHandler(e, funcName)

		defer csr.Close()

		counter := 0
		isDone := false
		countPerProcess := 1000
		countData := csr.Count()

		for !isDone && countData > 0 {
			scadas := []*ScadaDataOEM{}
			e = csr.Fetch(&scadas, countPerProcess, false)

			if len(scadas) < countPerProcess {
				isDone = true
			}

			wg.Add(1)
			go func(datas []*ScadaDataOEM, endIndex int) {
				tk.Printf("Starting process %v data\n", endIndex)

				mtxOem.Lock()
				logStart := time.Now()

				for _, data := range datas {
					d.updateScadaOEM(data)
				}

				logDurationg := time.Now().Sub(logStart)
				mtxOem.Unlock()

				tk.Printf("End processing for %v data about %v sec(s)\n", endIndex, logDurationg.Seconds())
				wg.Done()
			}(scadas, ((counter + 1) * countPerProcess))

			counter++
			if counter%10 == 0 || isDone {
				wg.Wait()
			}
		}
	}
}

func (u *UpdateScadaOemMinutes) updateScadaOEM(data *ScadaDataOEM) {
	ctx := u.Ctx
	turbine := GetExactTurbineId(strings.TrimSpace(data.Turbine))

	turbines := refTurbines.Get(turbine)

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

	avgWs := data.AI_intern_WindSpeed
	adjWs := tk.RoundingDown64(avgWs, 1)
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
	machinedown := 0.0
	griddown := 0.0

	timestamp := data.TimeStamp
	timestamp0 := timestamp.Add(-10 * time.Minute)

	totalDurationMttf := 0.0
	aDuration := 0.0

	// getting mttf in seconds
	refEvents := refAlarms.Get(turbine)
	alarms := []DowntimeEvent{}
	if refEvents != nil {
		dataAlarms := refEvents.([]DowntimeEvent)
		for _, a := range dataAlarms {
			if a.TimeStart.Sub(timestamp0) >= 0 && a.TimeEnd.Sub(timestamp) <= 0 {
				alarms = append(alarms, a)
			}
			if a.TimeStart.Sub(timestamp0) >= 0 && a.TimeStart.Sub(timestamp) <= 0 {
				alarms = append(alarms, a)
			}
			if a.TimeEnd.Sub(timestamp0) >= 0 && a.TimeEnd.Sub(timestamp) <= 0 {
				alarms = append(alarms, a)
			}
			if a.TimeStart.Sub(timestamp0) <= 0 && a.TimeStart.Sub(timestamp) >= 0 {
				alarms = append(alarms, a)
			}
		}
	}

	gridDowntime := 0.0
	machineDowntime := 0.0
	unknownDowntime := 0.0
	if len(alarms) > 0 {
		for _, a := range alarms {
			startTime := timestamp0
			endTime := timestamp
			if startTime.Sub(a.TimeStart) > 0 {
				startTime = a.TimeStart
			}
			if endTime.Sub(a.TimeEnd) > 0 {
				endTime = a.TimeEnd
			}
			aDuration += endTime.Sub(startTime).Seconds()
			if a.DownGrid {
				gridDowntime += endTime.Sub(startTime).Seconds()
			}
			if a.DownMachine {
				machineDowntime += endTime.Sub(startTime).Seconds()
			}
			if a.DownEnvironment {
				unknownDowntime += endTime.Sub(startTime).Seconds()
			}
		}
		totalDurationMttf = tk.Div(aDuration, float64(len(alarms)))
	}
	mttf := totalDurationMttf
	mttr := oktime - aDuration

	if denPower > data.AI_intern_ActivPower {
		energyLost = tk.Div(aDuration, 3600) * pcValue
	}

	totalavail := tk.Div(oktime, 600.0)
	machineavail := tk.Div((600.0 - machinedown), 600.0)
	gridavail := tk.Div((600.0 - griddown), 600.0)
	powerLost := denPower - data.AI_intern_ActivPower

	perfIndex := 0.0
	if denPower > 0 && data.AI_intern_ActivPower > 0 {
		perfIndex = tk.Div(data.AI_intern_ActivPower, denPower)
	}

	retadjws := tk.RoundingAuto64(data.AI_intern_WindSpeed, 0)
	retavgws := tk.RoundingAuto64(data.AI_intern_WindSpeed, 0)

	e := ctx.Connection.NewQuery().Update().From(new(ScadaDataOEM).TableName()).
		Where(dbox.Eq("_id", data.ID)).
		Exec(tk.M{}.Set("data", tk.M{}.
		Set("turbine", turbine).
		Set("totalavail", totalavail).
		Set("machineavail", machineavail).
		Set("gridavail", gridavail).
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
		Set("performanceindex", perfIndex).
		Set("griddowntime", gridDowntime).
		Set("machinedowntime", machineDowntime).
		Set("unknowndowntime", unknownDowntime)))

	if e != nil {
		tk.Printf("Update fail: %s", e.Error())
	}
}
