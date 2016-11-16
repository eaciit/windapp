package generatorControllers

import (
	. "github.com/eaciit/windapp/library/helper"
	. "github.com/eaciit/windapp/library/models"
	. "github.com/eaciit/windapp/processapp/controllers"
	"os"
	"sync"
	"time"

	"github.com/eaciit/crowd"
	"github.com/eaciit/dbox"
	_ "github.com/eaciit/dbox/dbc/mongo"
	tk "github.com/eaciit/toolkit"
)

var (
	maxDataEachProcessScada = 1000
	muinsertScada           = &sync.Mutex{}
	funcNameScada           = "GenScadaAlarmAnomaly Data"
)

// GenScadaAlarmAnomaly
type GenScadaAlarmAnomaly struct {
	*BaseController
}

// Generate
func (d *GenScadaAlarmAnomaly) Generate(base *BaseController) {
	startProcessTime := time.Now()

	if base != nil {
		d.BaseController = base

		ctx, e := PrepareConnection()
		if e != nil {
			ErrorHandler(e, funcNameScada)
			os.Exit(0)
		}

		scadaDatas := []ScadaData{}
		csr, e := ctx.NewQuery().From(new(ScadaData).TableName()).Where(dbox.Lt("oktime", 600)).Cursor(nil)
		e = csr.Fetch(&scadaDatas, 0, false)
		ErrorHandler(e, funcNameScada)
		csr.Close()

		tk.Println("GenScadaAlarmAnomaly Data")

		totalData := len(scadaDatas)
		worker := tk.ToInt(tk.ToFloat64(totalData/maxDataEachProcessScada, 0, tk.RoundingUp), tk.RoundingUp)

		wg := &sync.WaitGroup{}
		wg.Add(worker + 1)

		tk.Printf("totalData: %v worker: %v \n", totalData, worker)

		alarms := []Alarm{}
		csr, e = ctx.NewQuery().From(new(Alarm).TableName()).Cursor(nil)
		e = csr.Fetch(&alarms, 0, false)
		ErrorHandler(e, funcNameScada)
		csr.Close()

		for work := 0; work <= worker; work++ {
			from := work * maxDataEachProcessScada
			until := maxDataEachProcessScada + from

			if work == worker {
				until = totalData
			}

			// tk.Printf("%v | %v | %v \n", work, from, until)
			go d.doGen(scadaDatas[from:until], alarms, wg)
		}

		otherScada := []ScadaData{}
		csr, e = ctx.NewQuery().From(new(ScadaData).TableName()).Where(dbox.Gte("oktime", 600)).Cursor(nil)
		e = csr.Fetch(&otherScada, 0, false)
		ErrorHandler(e, funcNameScada)
		csr.Close()

		for _, scada := range otherScada {
			clean := new(ScadaClean).New()

			clean.TimeStamp = scada.TimeStamp
			clean.DateInfo = scada.DateInfo
			clean.Turbine = scada.Turbine
			clean.GridFrequency = scada.GridFrequency
			clean.ReactivePower = scada.ReactivePower
			clean.AlarmExtStopTime = scada.AlarmExtStopTime
			clean.AlarmGridDownTime = scada.AlarmGridDownTime
			clean.AlarmInterLineDown = scada.AlarmInterLineDown
			clean.AlarmMachDownTime = scada.AlarmMachDownTime
			clean.AlarmOkTime = scada.AlarmOkTime
			clean.AlarmUnknownTime = scada.AlarmUnknownTime
			clean.AlarmWeatherStop = scada.AlarmWeatherStop
			clean.ExternalStopTime = scada.ExternalStopTime
			clean.GridDownTime = scada.GridDownTime
			clean.GridOkSecs = scada.GridOkSecs
			clean.InternalLineDown = scada.InternalLineDown
			clean.MachineDownTime = scada.MachineDownTime
			clean.OkSecs = scada.OkSecs
			clean.OkTime = scada.OkTime
			clean.UnknownTime = scada.UnknownTime
			clean.WeatherStopTime = scada.WeatherStopTime
			clean.GeneratorRPM = scada.GeneratorRPM
			clean.NacelleYawPositionUntwist = scada.NacelleYawPositionUntwist
			clean.NacelleTemperature = scada.NacelleTemperature
			clean.AdjWindSpeed = scada.AdjWindSpeed
			clean.AmbientTemperature = scada.AmbientTemperature
			clean.AvgBladeAngle = scada.AvgBladeAngle
			clean.AvgWindSpeed = scada.AvgWindSpeed
			clean.UnitsGenerated = scada.UnitsGenerated
			clean.EstimatedPower = scada.EstimatedPower
			clean.NacelDirection = scada.NacelDirection
			clean.Power = scada.Power
			clean.PowerLost = scada.PowerLost
			clean.RotorRPM = scada.RotorRPM
			clean.WindDirection = scada.WindDirection
			clean.IsValidTimeDuration = scada.IsValidTimeDuration
			clean.Line = scada.Line
			clean.TotalTime = clean.ExternalStopTime + clean.GridDownTime + clean.InternalLineDown + clean.MachineDownTime + clean.OkTime
			clean.Minutes = scada.Minutes

			// muinsert.Lock()
			e := d.BaseController.Ctx.Insert(clean)
			ErrorHandler(e, funcName)
		}

		wg.Wait()

		// tk.Printf("totaldata: %v \n", total)
		// tk.Printf("totaldata: %v \n", len(result))
	}
	duration := time.Since(startProcessTime)
	tk.Printf("finish in: %v \n", duration.String())
}

func (d *GenScadaAlarmAnomaly) doGen(scadaDatas []ScadaData, alarms []Alarm, wg *sync.WaitGroup) {
	// ctx := d.BaseController.Ctx.Connection
	count := 0
	total := 0
	countClean := 0
	totalClean := 0

	for _, scada := range scadaDatas {

		// tk.Printf("%v - %v | %v - %v -> %v len: %v line: %v \n", startDate.String(), anomaly.StartDate.UTC().String(), anomaly.EndDate.UTC().String(), endDate.String(), anomaly.Turbine, len(scadaDatas), anomaly.Line)

		/*var filter []*dbox.Filter

		filter = append(filter, dbox.Lte("stardate", scada.TimeStamp.UTC()))
		filter = append(filter, dbox.Gte("enddate", scada.TimeStamp.UTC()))
		filter = append(filter, dbox.Eq("turbine", scada.Turbine))

		csr, e := ctx.NewQuery().From(new(Alarm).TableName()).Where(dbox.And(filter...)).Cursor(nil)
		ErrorHandler(e, funcNameScada)
		csr.Close()

		exist := csr.Count()*/

		exist := len(crowd.From(&alarms).Where(func(x interface{}) interface{} {
			y := x.(Alarm)

			isBefore := y.StartDate.UTC().Before(scada.TimeStamp.UTC()) || y.StartDate.UTC().Equal(scada.TimeStamp.UTC())
			isAfter := y.EndDate.UTC().After(scada.TimeStamp.UTC()) || y.EndDate.UTC().Equal(scada.TimeStamp.UTC())

			return isBefore && isAfter && y.Turbine == scada.Turbine
		}).Exec().Result.Data().([]Alarm))

		// tk.Printf("len: %v \n", exist)

		if exist == 0 {
			anomaly := new(ScadaAlarmAnomaly).New()

			anomaly.TimeStamp = scada.TimeStamp
			anomaly.DateInfo = scada.DateInfo
			anomaly.Turbine = scada.Turbine
			anomaly.GridFrequency = scada.GridFrequency
			anomaly.ReactivePower = scada.ReactivePower
			anomaly.AlarmExtStopTime = scada.AlarmExtStopTime
			anomaly.AlarmGridDownTime = scada.AlarmGridDownTime
			anomaly.AlarmInterLineDown = scada.AlarmInterLineDown
			anomaly.AlarmMachDownTime = scada.AlarmMachDownTime
			anomaly.AlarmOkTime = scada.AlarmOkTime
			anomaly.AlarmUnknownTime = scada.AlarmUnknownTime
			anomaly.AlarmWeatherStop = scada.AlarmWeatherStop
			anomaly.ExternalStopTime = scada.ExternalStopTime
			anomaly.GridDownTime = scada.GridDownTime
			anomaly.GridOkSecs = scada.GridOkSecs
			anomaly.InternalLineDown = scada.InternalLineDown
			anomaly.MachineDownTime = scada.MachineDownTime
			anomaly.OkSecs = scada.OkSecs
			anomaly.OkTime = scada.OkTime
			anomaly.UnknownTime = scada.UnknownTime
			anomaly.WeatherStopTime = scada.WeatherStopTime
			anomaly.GeneratorRPM = scada.GeneratorRPM
			anomaly.NacelleYawPositionUntwist = scada.NacelleYawPositionUntwist
			anomaly.NacelleTemperature = scada.NacelleTemperature
			anomaly.AdjWindSpeed = scada.AdjWindSpeed
			anomaly.AmbientTemperature = scada.AmbientTemperature
			anomaly.AvgBladeAngle = scada.AvgBladeAngle
			anomaly.AvgWindSpeed = scada.AvgWindSpeed
			anomaly.UnitsGenerated = scada.UnitsGenerated
			anomaly.EstimatedPower = scada.EstimatedPower
			anomaly.NacelDirection = scada.NacelDirection
			anomaly.Power = scada.Power
			anomaly.PowerLost = scada.PowerLost
			anomaly.RotorRPM = scada.RotorRPM
			anomaly.WindDirection = scada.WindDirection
			anomaly.IsValidTimeDuration = scada.IsValidTimeDuration
			anomaly.Line = scada.Line
			anomaly.TotalTime = anomaly.ExternalStopTime + anomaly.GridDownTime + anomaly.InternalLineDown + anomaly.MachineDownTime + anomaly.OkTime
			anomaly.Minutes = scada.Minutes

			// muinsert.Lock()
			e := d.BaseController.Ctx.Insert(anomaly)
			ErrorHandler(e, funcName)
			// muinsert.Unlock()

			count++
			total++
		} else {
			clean := new(ScadaClean).New()

			clean.TimeStamp = scada.TimeStamp
			clean.DateInfo = scada.DateInfo
			clean.Turbine = scada.Turbine
			clean.GridFrequency = scada.GridFrequency
			clean.ReactivePower = scada.ReactivePower
			clean.AlarmExtStopTime = scada.AlarmExtStopTime
			clean.AlarmGridDownTime = scada.AlarmGridDownTime
			clean.AlarmInterLineDown = scada.AlarmInterLineDown
			clean.AlarmMachDownTime = scada.AlarmMachDownTime
			clean.AlarmOkTime = scada.AlarmOkTime
			clean.AlarmUnknownTime = scada.AlarmUnknownTime
			clean.AlarmWeatherStop = scada.AlarmWeatherStop
			clean.ExternalStopTime = scada.ExternalStopTime
			clean.GridDownTime = scada.GridDownTime
			clean.GridOkSecs = scada.GridOkSecs
			clean.InternalLineDown = scada.InternalLineDown
			clean.MachineDownTime = scada.MachineDownTime
			clean.OkSecs = scada.OkSecs
			clean.OkTime = scada.OkTime
			clean.UnknownTime = scada.UnknownTime
			clean.WeatherStopTime = scada.WeatherStopTime
			clean.GeneratorRPM = scada.GeneratorRPM
			clean.NacelleYawPositionUntwist = scada.NacelleYawPositionUntwist
			clean.NacelleTemperature = scada.NacelleTemperature
			clean.AdjWindSpeed = scada.AdjWindSpeed
			clean.AmbientTemperature = scada.AmbientTemperature
			clean.AvgBladeAngle = scada.AvgBladeAngle
			clean.AvgWindSpeed = scada.AvgWindSpeed
			clean.UnitsGenerated = scada.UnitsGenerated
			clean.EstimatedPower = scada.EstimatedPower
			clean.NacelDirection = scada.NacelDirection
			clean.Power = scada.Power
			clean.PowerLost = scada.PowerLost
			clean.RotorRPM = scada.RotorRPM
			clean.WindDirection = scada.WindDirection
			clean.IsValidTimeDuration = scada.IsValidTimeDuration
			clean.Line = scada.Line
			clean.TotalTime = clean.ExternalStopTime + clean.GridDownTime + clean.InternalLineDown + clean.MachineDownTime + clean.OkTime
			clean.Minutes = scada.Minutes

			// muinsert.Lock()
			e := d.BaseController.Ctx.Insert(clean)
			ErrorHandler(e, funcName)
			// muinsert.Unlock()

			countClean++
			totalClean++
		}

		/*tk.Printf("idx: %v \n", idx+1)

		if count == 10 {
			tk.Printf("count: %v \n", total)
			count = 0
		}*/
	}
	tk.Printf("count: %v \n", total)
	tk.Printf("count clean: %v \n", totalClean)
	wg.Done()
}
