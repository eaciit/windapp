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

// ConvScadaData
type GenScadaDataExceptionDurationTime struct {
	*BaseController
}

// Generate
func (d *GenScadaDataExceptionDurationTime) Generate(base *BaseController) {
	funcName := "GenScadaDataExceptionDurationTime Data"
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

		csr, e := ctx.NewQuery().From(new(ScadaData).TableName()).Cursor(nil)

		e = csr.Fetch(&scadas, 0, false)
		ErrorHandler(e, funcName)
		csr.Close()
		tk.Println("GenScadaDataExceptionDurationTime Data")
		for _, data := range scadas {
			// totalTimeDuration := data.AlarmUnknownTime + data.AlarmWeatherStop + data.ExternalStopTime + data.GridDownTime + data.InternalLineDown + data.MachineDownTime + data.OkTime
			totalTimeDuration := data.ExternalStopTime + data.GridDownTime + data.InternalLineDown + data.MachineDownTime + data.OkTime
			if totalTimeDuration >= 600.02 || totalTimeDuration <= 599.98 {
				e = ctx.NewQuery().Update().From(new(ScadaData).TableName()).Where(dbox.Eq("_id", data.ID)).Exec(tk.M{}.Set("data", tk.M{}.Set("isvalidtimeduration", false)))
				if e != nil {
					tk.Printf("Update fail: %s", e.Error())
				}

				count++
				total++
			} else {
				e = ctx.NewQuery().Update().From(new(ScadaData).TableName()).Where(dbox.Eq("_id", data.ID)).Exec(tk.M{}.Set("data", tk.M{}.Set("isvalidtimeduration", true)))
				if e != nil {
					tk.Printf("Update fail: %s", e.Error())
				}
			}

			if count == 1000 {
				tk.Printf("count: %v \n", total)
				count = 0
			}
		}

		tk.Printf("totaldata: %v \n", total)
		// tk.Printf("totaldata: %v \n", len(result))
	}
}
