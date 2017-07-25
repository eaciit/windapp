package generatorControllers

import (
	. "eaciit/wfdemo-git/library/helper"
	. "eaciit/wfdemo-git/library/models"
	. "eaciit/wfdemo-git/processapp/summaryGenerator/controllers"
	_ "math"
	"os"
	_ "strings"
	"sync"
	"time"

	"github.com/eaciit/dbox"
	_ "github.com/eaciit/dbox/dbc/mongo"
	tk "github.com/eaciit/toolkit"
)

type UpdateOEMToScada struct {
	*BaseController
}

var (
	minValueFloat = -99999.00
	minValueInt   = -99999
)

func (u *UpdateOEMToScada) RunMapping(base *BaseController) {
	u.BaseController = base
	funcName := "Mapping Scada OEM to Scada Data"

	conn, e := PrepareConnection()
	if e != nil {
		ErrorHandler(e, funcName)
		os.Exit(0)
	}

	tk.Println(funcName)

	var wg sync.WaitGroup

	for turbine, _ := range u.BaseController.RefTurbines {
		filter := []*dbox.Filter{}
		filter = append(filter, dbox.Eq("projectname", "Tejuva"))
		filter = append(filter, dbox.Eq("turbine", turbine))

		latestDate := u.BaseController.GetLatest("ScadaData", "Tejuva", turbine)
		if latestDate.Format("2006") != "0001" {
			filter = append(filter, dbox.Gt("timestamp", latestDate))
		}
		// filter = append(filter, dbox.Gt("timestamp", u.BaseController.LatestData.MapScadaData["Tejuva#"+turbine]))

		csr, e := conn.NewQuery().From(new(ScadaDataOEM).TableName()).
			Where(filter...).Cursor(nil)
		ErrorHandler(e, funcName)

		defer csr.Close()

		counter := 0
		isDone := false
		countPerProcess := 1000
		countData := csr.Count()

		tk.Printf("\nOEM to Scada for %v | %v \n", turbine, countData)

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
					u.mapOEMToScada(data)
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

func (u *UpdateOEMToScada) mapOEMToScada(data *ScadaDataOEM) {
	scada := new(ScadaData).New()

	scada.TimeStamp = data.TimeStamp
	scada.DateInfo = data.DateInfo
	scada.ProjectName = data.ProjectName
	scada.Turbine = data.Turbine
	scada.Minutes = 10
	scada.GridFrequency = data.AI_intern_Frequency_Grid
	scada.ReactivePower = data.AI_intern_ReactivPower
	scada.AlarmExtStopTime = minValueFloat
	scada.AlarmGridDownTime = minValueFloat
	scada.AlarmInterLineDown = minValueFloat
	scada.AlarmMachDownTime = minValueFloat
	scada.AlarmOkTime = minValueFloat
	scada.AlarmUnknownTime = minValueFloat
	scada.AlarmWeatherStop = minValueFloat
	scada.ExternalStopTime = minValueFloat
	scada.GridDownTime = data.GridDowntime
	scada.GridDownTimeAll = data.GridDowntimeAll
	scada.GridOkSecs = 600.0 - data.GridDowntime
	scada.InternalLineDown = minValueFloat
	scada.MachineDownTime = data.MachineDowntime
	scada.MachineDownTimeAll = data.MachineDowntimeAll
	// scada.OkSecs = data.MTTR
	// scada.OkTime = data.MTTR

	scada.OkTime = (600 - (data.GridDowntime + data.MachineDowntime + data.UnknownDowntime))
	scada.OkSecs = scada.OkTime

	scada.UnknownTime = data.UnknownDowntime
	scada.UnknownTimeAll = data.UnknownDowntimeAll
	scada.WeatherStopTime = minValueFloat
	scada.GeneratorRPM = data.C_intern_SpeedGenerator
	scada.NacelleYawPositionUntwist = data.AI_intern_NacelleDrill_at_NorthPosSensor
	scada.NacelleTemperature = data.Temp_Nacelle
	scada.AdjWindSpeed = tk.RoundingAuto64(data.AI_intern_WindSpeed, 1)
	scada.AmbientTemperature = data.Temp_Outdoor
	scada.AvgBladeAngle = minValueFloat
	scada.AvgWindSpeed = data.AI_intern_WindSpeed
	scada.UnitsGenerated = minValueFloat
	scada.EstimatedPower = data.DenPower
	scada.EstimatedEnergy = data.DenEnergy
	scada.NacelDirection = data.AI_intern_NacellePos
	scada.Power = data.AI_intern_ActivPower
	scada.PowerLost = data.DenPower - data.AI_intern_ActivPower
	scada.Energy = data.Energy
	scada.EnergyLost = data.EnergyLost
	scada.RotorRPM = data.C_intern_SpeedRotor
	scada.WindDirection = data.AI_intern_WindDirection
	scada.Line = data.Line
	scada.IsValidTimeDuration = true
	scada.TotalTime = 600.0
	scada.Available = 1
	scada.DenValue = data.DenValue
	scada.DenPh = data.DenPh
	scada.DenWindSpeed = data.DenWindSpeed
	scada.DenAdjWindSpeed = data.DenAdjWindSpeed
	scada.DenPower = data.DenPower
	scada.DenEnergy = data.DenEnergy
	scada.PCValue = data.PCValue
	scada.PCValueAdj = data.PCValueAdj
	scada.PCDeviation = data.PCDeviation
	scada.WSAdjForPC = data.WSAdjForPC
	scada.WSAvgForPC = data.WSAvgForPC
	scada.TotalAvail = tk.Div(scada.OkTime, 600.0) //data.MTTR / 600.0
	scada.MachineAvail = (600.0 - data.MachineDowntime) / 600
	scada.GridAvail = (600.0 - data.GridDowntime) / 600.0
	scada.MachineAvailAll = (600.0 - data.MachineDowntimeAll) / 600
	scada.GridAvailAll = (600.0 - data.GridDowntimeAll) / 600.0
	scada.DenPcDeviation = data.DenPcDeviation
	scada.DenDeviationPct = data.DenDeviationPct
	scada.DenPcValue = data.DenPcValue
	scada.DeviationPct = data.DeviationPct
	scada.MTTR = data.MTTR
	scada.MTTF = data.MTTF
	scada.PerformanceIndex = data.PerformanceIndex

	u.Ctx.Insert(scada)
}
