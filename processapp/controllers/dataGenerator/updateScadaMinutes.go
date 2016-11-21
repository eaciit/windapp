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
	"os/exec"
	"time"
)

// UpdateScadaMinutes
type UpdateScadaMinutes struct {
	*BaseController
}

// Generate
func (d *UpdateScadaMinutes) Generate(base *BaseController) {
	funcName := "UpdateScadaMinutes Data"
	count := 0
	total := 0

	_ = count
	_ = total
	if base != nil {
		d.BaseController = base

		ctx, e := PrepareConnection()
		if e != nil {
			ErrorHandler(e, funcName)
			os.Exit(0)
		}

		scadas := []ScadaData{}
		scadasAlarm := []ScadaAlarmAnomaly{}

		csr, e := ctx.NewQuery().From(new(ScadaData).TableName()).Cursor(nil)

		e = csr.Fetch(&scadas, 0, false)
		ErrorHandler(e, funcName)
		csr.Close()
		tk.Println("UpdateScadaMinutes Data")
		for _, data := range scadas {
			// totalTimeDuration := data.AlarmUnknownTime + data.AlarmWeatherStop + data.ExternalStopTime + data.GridDownTime + data.InternalLineDown + data.MachineDownTime + data.OkTime

			var available int

			if data.AvgWindSpeed < 4 || (data.AvgWindSpeed >= 4 && data.Power > 0) {
				available = 1
			} else if data.AvgWindSpeed >= 4 && data.Power <= 0 {
				available = 0
			}

			e = ctx.NewQuery().Update().From(new(ScadaData).TableName()).Where(dbox.Eq("_id", data.ID)).Exec(tk.M{}.Set("data", tk.M{}.Set("minutes", 10).Set("available", available)))
			if e != nil {
				tk.Printf("Update fail: %s", e.Error())
			}

			count++
			total++

			if count == 1000 {
				tk.Printf("count: %v \n", total)
				count = 0
			}
		}

		tk.Printf("totaldata: %v \n", total)

		csr, e = ctx.NewQuery().From(new(ScadaAlarmAnomaly).TableName()).Cursor(nil)

		e = csr.Fetch(&scadasAlarm, 0, false)
		ErrorHandler(e, funcName)
		csr.Close()
		tk.Println("UpdateScadaMinutes ScadaAlarm Data")
		for _, data := range scadasAlarm {

			var available int

			if data.AvgWindSpeed < 4 || (data.AvgWindSpeed >= 4 && data.Power > 0) {
				available = 1
			} else if data.AvgWindSpeed >= 4 && data.Power <= 0 {
				available = 0
			}

			e = ctx.NewQuery().Update().From(new(ScadaAlarmAnomaly).TableName()).Where(dbox.Eq("_id", data.ID)).Exec(tk.M{}.Set("data", tk.M{}.Set("minutes", 10).Set("available", available)))
			if e != nil {
				tk.Printf("Update fail: %s", e.Error())
			}

			count++
			total++

			if count == 1000 {
				tk.Printf("count: %v \n", total)
				count = 0
			}
		}

		tk.Printf("totaldata: %v \n", total)
		// tk.Printf("totaldata: %v \n", len(result))
	}
}

func (d *UpdateScadaMinutes) GenerateDensity(base *BaseController) {
	exec.Command("name", "aksjdkasd").Output()

	funcName := "UpdateScadaDensity Data"
	count := 0
	total := 0

	//_ = count
	_ = total
	if base != nil {
		d.BaseController = base

		ctx, e := PrepareConnection()
		if e != nil {
			ErrorHandler(e, funcName)
			os.Exit(0)
		}

		scadas := []ScadaData{}
		csr, e := ctx.NewQuery().From(new(ScadaData).TableName()).
			Where(dbox.Eq("projectname", "Tejuva")).Cursor(nil)

		e = csr.Fetch(&scadas, 0, false)
		ErrorHandler(e, funcName)
		csr.Close()
		tk.Println("UpdateScadaDensity Data")
		for _, data := range scadas {

			// tk.Println("Processing data for " + data.Turbine + " on " + data.ID.String())

			turbines := []TurbineMaster{}
			csrt, e := ctx.NewQuery().From(new(TurbineMaster).TableName()).
				Where(dbox.Eq("turbineid", data.Turbine)).Cursor(nil)

			e = csrt.Fetch(&turbines, 0, false)
			ErrorHandler(e, funcName)
			csrt.Close()

			energy := tk.Div(data.Power, 6)
			estimatedEnergy := tk.Div(data.EstimatedPower, 6)
			energyLost := tk.Div(data.PowerLost, 6)

			pH := 0.0
			elevation := 0.0
			temperature := data.NacelleTemperature
			//tk.Printf("Temp : #%v\n", temperature)
			if turbines != nil {
				elevation = turbines[0].Elevation
				//tk.Printf("Elev : #%v\n", elevation)
				exponen := (-(9.80665) * 28.9644 * elevation) / (8314.32 * 288.15)
				//tk.Printf("Exp : #%v\n", exponen)
				pH = 101325 * math.Exp(exponen)
				//tk.Printf("PH : #%v\n", pH)
			}

			denWs := 0.0
			denPower := 0.0
			adjDenWs := 0.0

			density := pH / (287.05 * (273.15 + temperature))

			// pcValue := 0.0
			avgWs := data.AvgWindSpeed
			adjWs := data.AdjWindSpeed
			// tk.Printf("Ws scd = %v\n", avgWs)
			pcValue, retavgws, _ := GetPowerCurve(ctx, avgWs)
			pcValueAdj, retadjws, _ := GetPowerCurve(ctx, adjWs)
			pcDeviation := pcValue - data.Power

			denWs = avgWs * math.Pow((density/1.225), (1.0/3.0))
			adjDenWs = tk.RoundingDown64(denWs, 0)
			if denWs < 3.75 && denWs >= 3.5 {
				adjDenWs = 3.5
			}
			denPower, _, _ = GetPowerCurve(ctx, denWs)

			denPcValue, _, _ := GetPowerCurve(ctx, denWs)
			denPcDeviation := denPcValue - denPower

			deviationPct := 0.0
			denDeviationPct := 0.0
			if pcDeviation > 0 {
				deviationPct = tk.Div(pcDeviation, pcValue)
			}
			if denPcDeviation > 0 {
				denDeviationPct = tk.Div(denPcDeviation, denPcValue)
			}

			oktime := data.OkTime
			machinedown := data.MachineDownTime
			griddown := data.GridDownTime

			// getting mttr in seconds
			mttr := oktime

			timestamp := data.TimeStamp
			timestamp0 := timestamp.Add(10 * time.Minute)

			// getting mttf in seconds
			alarms := []Alarm{}
			csr, e := ctx.NewQuery().From(new(Alarm).TableName()).
				Where(dbox.And(dbox.Eq("projectname", "Tejuva"), dbox.Eq("turbine", data.Turbine), dbox.Lte("startdate", timestamp0), dbox.Gte("enddate", timestamp))).Cursor(nil)

			e = csr.Fetch(&alarms, 0, false)
			ErrorHandler(e, "get alarm data")
			csr.Close()

			totalDurationMttf := 0.0
			if len(alarms) > 0 {
				aDuration := 0.0
				for _, a := range alarms {
					startTime := timestamp0
					endTime := timestamp
					if a.StartDate.Sub(timestamp0) > 0 {
						startTime = a.StartDate
					}
					if timestamp.Sub(a.EndDate) > 0 {
						endTime = a.EndDate
					}
					aDuration += endTime.Sub(startTime).Seconds()
				}
				totalDurationMttf = tk.Div(aDuration, float64(len(alarms)))
			}
			mttf := totalDurationMttf

			totalavail := tk.Div(oktime, 600.0)
			machineavail := tk.Div((600.0 - machinedown), 600.0)
			gridavail := tk.Div((600.0 - griddown), 600.0)

			// powerCurve := []PowerCurveModel{}
			// csrp, e := ctx.NewQuery().From(new(PowerCurveModel).TableName()).
			// 	Where(dbox.Eq("windspeed", adjDenWs)).Cursor(nil)

			// e = csrp.Fetch(&powerCurve, 0, false)
			// ErrorHandler(e, funcName)
			// csrp.Close()

			// if len(powerCurve) > 0 {
			// 	denPower = powerCurve[0].Power1
			// 	//tk.Printf("Power Curve : #%v\n", denPower)
			// }

			e = ctx.NewQuery().Update().From(new(ScadaData).TableName()).
				Where(dbox.Eq("_id", data.ID)).
				Exec(tk.M{}.Set("data", tk.M{}.
				Set("totalavail", totalavail).
				Set("machineavail", machineavail).
				Set("gridavail", gridavail).
				Set("wsadjforpc", retadjws).
				Set("wsavgforpc", retavgws).
				Set("pcdeviation", pcDeviation).
				Set("pcvalue", pcValue).
				Set("pcvalueadj", pcValueAdj).
				Set("estimatedenergy", estimatedEnergy).
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
				Set("mttf", mttf)))
			if e != nil {
				tk.Printf("Update fail: %s", e.Error())
			}

			count++
			total++

			if count == 1000 {
				tk.Printf("count: %v \n", total)
				count = 0
			}
			// break
		}

		tk.Printf("totaldata: %v \n", total)
	}
}

func GetPowerCurve(ctx dbox.IConnection, avgWs float64) (float64, float64, float64) {
	totalPower := 0.0

	wsf0 := 0.0
	wsf1 := 0.0
	wsret := 0.0
	wsavgret := tk.RoundingAuto64(avgWs, 1)

	if avgWs >= 3.5 && avgWs < 3.75 {
		wsf0 = 3.5
		wsf1 = 4
		wsret = 3.5
	} else if avgWs >= 3 && avgWs < 3.5 {
		wsf0 = 3
		wsf1 = 3.5
		wsret = 3
	} else {
		wsf0 = avgWs
		wsf1 = avgWs
		wsret = tk.RoundingAuto64(avgWs, 0)
	}
	// tk.Printf("%v-%v-%v\n", avgWs, wsf0, wsf1)

	pcs0 := []PowerCurveModel{}
	pcs1 := []PowerCurveModel{}
	csrps0, e := ctx.NewQuery().From(new(PowerCurveModel).TableName()).
		Where(dbox.Lte("windspeed", wsf0)).
		Order("-windspeed").Take(1).Skip(0).Cursor(nil)
	e = csrps0.Fetch(&pcs0, 0, false)
	ErrorHandler(e, funcName)
	csrps0.Close()
	// tk.Printf("%v\n", pcs0)
	csrps1, e := ctx.NewQuery().From(new(PowerCurveModel).TableName()).
		Where(dbox.Gte("windspeed", wsf1)).
		Order("windspeed").Take(1).Skip(0).Cursor(nil)
	e = csrps1.Fetch(&pcs1, 0, false)
	ErrorHandler(e, funcName)
	csrps1.Close()
	// tk.Printf("%v\n", pcs1)

	ws0 := 0.0
	power0 := 0.0
	if len(pcs0) > 0 {
		power0 = pcs0[0].Power1
		ws0 = pcs0[0].WindSpeed
	}

	ws1 := 0.0
	power1 := 0.0
	if len(pcs1) > 0 {
		power1 = pcs1[0].Power1
		ws1 = pcs1[0].WindSpeed
	}

	if ws1-ws0 > 0 {
		totalPower = power0 + (avgWs-ws0)*(power1-power0)/(ws1-ws0)
	} else {
		totalPower = power0
	}

	return totalPower, wsret, wsavgret
}
