package generatorControllers

import (
	. "eaciit/ostrowfm/library/helper"
	. "eaciit/ostrowfm/library/models"
	. "eaciit/ostrowfm/processapp/controllers"
	"os"

	"github.com/eaciit/dbox"
	_ "github.com/eaciit/dbox/dbc/mongo"
	tk "github.com/eaciit/toolkit"
)

// ConvScadaClean
type UpdateProjectScadaAndAlarm struct {
	*BaseController
}

// Generate
func (d *UpdateProjectScadaAndAlarm) Generate(base *BaseController) {
	funcName := "UpdateProjectScadaAndAlarm Data"
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

		scadas := []ScadaClean{}
		scadasAlarm := []ScadaAlarmAnomaly{}

		csr, e := ctx.NewQuery().From(new(ScadaClean).TableName()).Cursor(nil)

		e = csr.Fetch(&scadas, 0, false)
		ErrorHandler(e, funcName)
		csr.Close()
		tk.Println("UpdateProjectScadaAndAlarm Data")
		for _, data := range scadas {
			// totalTimeDuration := data.AlarmUnknownTime + data.AlarmWeatherStop + data.ExternalStopTime + data.GridDownTime + data.InternalLineDown + data.MachineDownTime + data.OkTime

			e = ctx.NewQuery().Update().From(new(ScadaClean).TableName()).Where(dbox.Eq("_id", data.ID)).Exec(tk.M{}.Set("data", tk.M{}.Set("projectname", "Tejuva")))
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
		tk.Println("UpdateProjectScadaAndAlarm ScadaAlarm Data")
		for _, data := range scadasAlarm {
			e = ctx.NewQuery().Update().From(new(ScadaAlarmAnomaly).TableName()).Where(dbox.Eq("_id", data.ID)).Exec(tk.M{}.Set("data", tk.M{}.Set("projectname", "Tejuva")))
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

		csr, e = ctx.NewQuery().From(new(AlarmClean).TableName()).Cursor(nil)
		alarms := []AlarmClean{}

		e = csr.Fetch(&alarms, 0, false)
		ErrorHandler(e, funcName)
		csr.Close()
		tk.Println("UpdateProjectAlarms Data")
		for _, data := range alarms {
			e = ctx.NewQuery().Update().From(new(AlarmClean).TableName()).Where(dbox.Eq("_id", data.ID)).Exec(tk.M{}.Set("data", tk.M{}.Set("projectname", "Tejuva")))
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

		csr, e = ctx.NewQuery().From(new(AlarmOverlapping).TableName()).Cursor(nil)
		alarmsOverlapping := []AlarmOverlapping{}

		e = csr.Fetch(&alarmsOverlapping, 0, false)
		ErrorHandler(e, funcName)
		csr.Close()
		tk.Println("UpdateProjectAlarms Data")
		for _, data := range alarmsOverlapping {
			e = ctx.NewQuery().Update().From(new(AlarmOverlapping).TableName()).Where(dbox.Eq("_id", data.ID)).Exec(tk.M{}.Set("data", tk.M{}.Set("projectname", "Tejuva")))
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

		csr, e = ctx.NewQuery().From(new(AlarmScadaAnomaly).TableName()).Cursor(nil)
		alarmScadaAnomaly := []AlarmScadaAnomaly{}

		e = csr.Fetch(&alarmScadaAnomaly, 0, false)
		ErrorHandler(e, funcName)
		csr.Close()
		tk.Println("UpdateProjectAlarms Data")
		for _, data := range alarmScadaAnomaly {
			e = ctx.NewQuery().Update().From(new(AlarmScadaAnomaly).TableName()).Where(dbox.Eq("_id", data.ID)).Exec(tk.M{}.Set("data", tk.M{}.Set("projectname", "Tejuva")))
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
	}
}
