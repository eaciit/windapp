package generatorControllers

import (
	. "eaciit/wfdemo/library/helper"
	. "eaciit/wfdemo/library/models"
	. "eaciit/wfdemo/processapp/controllers"
	"os"

	_ "github.com/eaciit/dbox/dbc/mongo"
	tk "github.com/eaciit/toolkit"
)

// ConvAlarm
type GenAlarmOverlapping struct {
	*BaseController
}

// Generate
func (d *GenAlarmOverlapping) Generate(base *BaseController) {
	funcName := "Generating Alarm Overlapping Data"
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

		result := []*AlarmOverlapping{}
		alarms := []Alarm{}
		checked := []Alarm{}

		csr, e := ctx.NewQuery().From(new(Alarm).TableName()).Order("startdate").Cursor(nil)

		e = csr.Fetch(&alarms, 0, false)
		ErrorHandler(e, funcName)
		csr.Close()
		tk.Println("Generate Alarm Overlapping Data")
		for idx, alarm := range alarms {
			isChecked := false

		bChecked:
			for _, val := range checked {
				if alarm.StartDate.UnixNano() == val.StartDate.UnixNano() &&
					alarm.EndDate.UnixNano() == val.EndDate.UnixNano() &&
					alarm.Turbine == val.Turbine &&
					alarm.AlertDescription == val.AlertDescription &&
					alarm.ExternalStop == val.ExternalStop &&
					alarm.GridDown == val.GridDown &&
					alarm.InternalGrid == val.InternalGrid &&
					alarm.MachineDown == val.MachineDown &&
					alarm.AEbOK == val.AEbOK &&
					alarm.Unknown == val.Unknown &&
					alarm.WeatherStop == val.WeatherStop {

					isChecked = true
					break bChecked
				}
			}

			if isChecked == false {
				alarmsTmp := []Alarm{}
				alarmOverlapping := new(AlarmOverlapping).New()

				alarmOverlapping.Farm = alarm.Farm
				alarmOverlapping.StartDate = alarm.StartDate
				alarmOverlapping.StartDateInfo = alarm.StartDateInfo
				alarmOverlapping.Turbine = alarm.Turbine
				alarmOverlapping.AlertDescription = alarm.AlertDescription
				alarmOverlapping.ExternalStop = alarm.ExternalStop
				alarmOverlapping.GridDown = alarm.GridDown
				alarmOverlapping.InternalGrid = alarm.InternalGrid
				alarmOverlapping.MachineDown = alarm.MachineDown
				alarmOverlapping.AEbOK = alarm.AEbOK
				alarmOverlapping.Unknown = alarm.Unknown
				alarmOverlapping.WeatherStop = alarm.WeatherStop
				alarmOverlapping.EndDate = alarm.EndDate
				alarmOverlapping.Duration = alarm.Duration

				for _, alarmSub := range alarms[idx+1:] {

					if ((alarm.StartDate.UnixNano() <= alarmSub.StartDate.UnixNano() && alarm.EndDate.UnixNano() >= alarmSub.StartDate.UnixNano()) || (alarm.StartDate.UnixNano() <= alarmSub.EndDate.UnixNano() && alarm.EndDate.UnixNano() >= alarmSub.EndDate.UnixNano())) &&
						alarm.Turbine == alarmSub.Turbine &&
						alarm.AlertDescription != alarmSub.AlertDescription &&
						(alarm.ExternalStop != alarmSub.ExternalStop ||
							alarm.GridDown != alarmSub.GridDown ||
							alarm.InternalGrid != alarmSub.InternalGrid ||
							alarm.MachineDown != alarmSub.MachineDown ||
							alarm.AEbOK != alarmSub.AEbOK ||
							alarm.Unknown != alarmSub.Unknown ||
							alarm.WeatherStop != alarmSub.WeatherStop) {

						if alarmOverlapping.StartDate.UnixNano() > alarmSub.StartDate.UnixNano() {
							alarmOverlapping.StartDate = alarmSub.StartDate
							alarmOverlapping.StartDateInfo = alarmSub.StartDateInfo
						}

						if alarmOverlapping.EndDate.UnixNano() < alarmSub.EndDate.UnixNano() {
							alarmOverlapping.EndDate = alarmSub.EndDate
						}

						if len(alarmsTmp) == 0 {
							alarmsTmp = append(alarmsTmp, alarm)
						}

						alarmsTmp = append(alarmsTmp, alarmSub)
					}

				}

				if len(alarmsTmp) > 0 {
					checked = append(checked, alarmsTmp...)

					alarmOverlapping.Alarms = alarmsTmp
					result = append(result, alarmOverlapping)

					d.BaseController.Ctx.Insert(alarmOverlapping)
					count++
					total++
				}

				if count == 1000 {
					tk.Printf("count: %v \n", total)
					count = 0
				}
			}

		}

		tk.Printf("totaldata: %v \n", total)
		// tk.Printf("totaldata: %v \n", len(result))
	}
}
