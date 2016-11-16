package converterControllers

import (
	"bufio"
	. "github.com/eaciit/windapp/library/helper"
	. "github.com/eaciit/windapp/library/models"
	. "github.com/eaciit/windapp/processapp/controllers"
	"encoding/csv"
	_ "github.com/eaciit/dbox/dbc/mongo"
	tk "github.com/eaciit/toolkit"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

type ConvScadaDataCsv struct {
	*BaseController
}

func (d *ConvScadaDataCsv) Generate(base *BaseController) {
	if base != nil {
		project := "Tejuva"

		d.BaseController = base
		ctx := d.BaseController.Ctx

		sourceDir := "D:\\Works\\2016\\Ostro Wind Farm\\Documents\\Ostro Files\\10 min data tejua\\Merged"
		files, _ := ioutil.ReadDir(sourceDir)
		for _, f := range files {
			if strings.Index(f.Name(), "DataFile-T") > -1 {
				fr, _ := os.Open(sourceDir + "\\" + f.Name())
				tk.Println(f.Name())
				read := csv.NewReader(bufio.NewReader(fr))
				count := 0
				for {
					record, err := read.Read()
					if err == io.EOF {
						break
					}
					ts, e := time.Parse("02-Jan-2006 15:04", string(record[0]))
					if e == nil {
						turbines := strings.Split(record[2], ".")
						turbine := turbines[2]
						dateInfo := GetDateInfo(ts)

						// scadaToDelete := new(ScadaDataNew)
						// _ = ctx.Get(scadaToDelete, tk.M{}.Set("Where", tk.M{}.Set("timestamp", ts).Set("turbine", turbine)))

						// if scadaToDelete.ID != "" {
						// ctx.Delete(scadaToDelete)
						// tk.Printf("Delete scada for %v - %v", ts, turbine)
						// }

						power, _ := tk.StringToFloat(record[3])
						windSpeed, _ := tk.StringToFloat(record[97])
						windDirection, _ := tk.StringToFloat(record[212])
						nacelPos, e := tk.StringToFloat(record[127])
						if !e {
							nacelPos = 0.0
						}

						gridFreq := 0.0
						reactivePower, _ := tk.StringToFloat(record[78])
						alarmExtStopTime := 0.0
						alarmGridDownTime := 0.0
						AlarmerLineDown := 0.0
						AlarmMachDownTime := 0.0
						AlarmOkTime := 0.0
						AlarmUnknownTime := 0.0
						AlarmWeatherStop := 0.0
						ExternalStopTime := 0.0
						GridDownTime := 0.0
						GridOkSecs := 0.0
						ernalLineDown := 0.0
						MachineDownTime := 0.0
						OkSecs := 0.0
						OkTime := 0.0
						UnknownTime := 0.0
						WeatherStopTime := 0.0
						GeneratorRPM := 0.0
						NacelleYawPositionUntwist := 0.0
						NacelleTemperature := 0.0
						AdjWindSpeed := tk.RoundingAuto64(windSpeed, 1)
						AmbientTemperature := 0.0
						AvgBladeAngle := 0.0
						AvgWindSpeed := windSpeed
						UnitsGenerated := 0.0
						EstimatedPower := 0.0
						EstimatedEnergy := 0.0
						NacelDirection := nacelPos
						Power := power
						PowerLost := 0.0
						Energy := tk.Div(power, 6.0)
						EnergyLost := 0.0
						RotorRPM, _ := tk.StringToFloat(record[84])
						WindDirection := windDirection
						Line := count
						IsValidTimeDuration := false
						TotalTime := 600.0
						Minutes := 10
						Available := 1
						DenWindSpeed := 0.0
						DenAdjWindSpeed := 0.0
						DenPower := 0.0
						DenEnergy := 0.0

						if windSpeed >= 3.5 && power <= 0 {
							Available = 0
						} else {
							OkSecs = 600.0
							IsValidTimeDuration = true
							TotalTime = 600.0
						}

						scada := new(ScadaDataNew).New()

						// sid := tk.Sprintf("%v_%v_%v", project, turbine, strings.Replace(string(record[0]), " ", "_", -1))
						// scada.ID = sid
						scada.TimeStamp = ts
						scada.DateInfo = dateInfo
						scada.Turbine = turbine
						scada.ProjectName = project
						scada.GridFrequency = gridFreq
						scada.ReactivePower = reactivePower
						scada.AlarmExtStopTime = alarmExtStopTime
						scada.AlarmGridDownTime = alarmGridDownTime
						scada.AlarmInterLineDown = AlarmerLineDown
						scada.AlarmMachDownTime = AlarmMachDownTime
						scada.AlarmOkTime = AlarmOkTime
						scada.AlarmUnknownTime = AlarmUnknownTime
						scada.AlarmWeatherStop = AlarmWeatherStop
						scada.ExternalStopTime = ExternalStopTime
						scada.GridDownTime = GridDownTime
						scada.GridOkSecs = GridOkSecs
						scada.InternalLineDown = ernalLineDown
						scada.MachineDownTime = MachineDownTime
						scada.OkSecs = OkSecs
						scada.OkTime = OkTime
						scada.UnknownTime = UnknownTime
						scada.WeatherStopTime = WeatherStopTime
						scada.GeneratorRPM = GeneratorRPM
						scada.NacelleYawPositionUntwist = NacelleYawPositionUntwist
						scada.NacelleTemperature = NacelleTemperature
						scada.AdjWindSpeed = AdjWindSpeed
						scada.AmbientTemperature = AmbientTemperature
						scada.AvgBladeAngle = AvgBladeAngle
						scada.AvgWindSpeed = AvgWindSpeed
						scada.UnitsGenerated = UnitsGenerated
						scada.EstimatedPower = EstimatedPower
						scada.EstimatedEnergy = EstimatedEnergy
						scada.NacelDirection = NacelDirection
						scada.Power = Power
						scada.PowerLost = PowerLost
						scada.Energy = Energy
						scada.EnergyLost = EnergyLost
						scada.RotorRPM = RotorRPM
						scada.WindDirection = WindDirection
						scada.Line = Line
						scada.IsValidTimeDuration = IsValidTimeDuration
						scada.TotalTime = TotalTime
						scada.Minutes = Minutes
						scada.Available = Available
						scada.DenWindSpeed = DenWindSpeed
						scada.DenAdjWindSpeed = DenAdjWindSpeed
						scada.DenPower = DenPower
						scada.DenEnergy = DenEnergy

						ctx.Insert(scada)

						count++
					}
				}
			}
		}
	}
}
