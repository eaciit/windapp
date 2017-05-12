package generatorControllers

import (
	. "eaciit/wfdemo-git/library/helper"
	. "eaciit/wfdemo-git/library/models"
	. "eaciit/wfdemo-git/processapp/summaryGenerator/controllers"
	"os"

	"github.com/eaciit/dbox"
	_ "github.com/eaciit/dbox/dbc/mongo"
	tk "github.com/eaciit/toolkit"
)

type EventReduceAvailability struct {
	*BaseController
}

func (ev *EventReduceAvailability) ConvertEventReduceAvailability(base *BaseController) {
	ev.BaseController = base
	tk.Println("Start process ConvertEventReduceAvailability...")

	funcName := "EventReduceAvailabilityConversion"

	ctx, e := PrepareConnection()
	if e != nil {
		ErrorHandler(e, funcName)
		os.Exit(0)
	}

	brakeReducesAvailability, e := PopulateReducesAvailability(ev.Ctx)

	for turbine, _ := range ev.BaseController.RefTurbines {
		// xTurbines = append(xTurbines, turbine)
		// wg.Add(1)
		// go func(t string) {
		t := turbine
		filterX := []*dbox.Filter{}
		filterX = append(filterX, dbox.Eq("projectname", projectName))
		filterX = append(filterX, dbox.Eq("turbine", t))

		csr, e := ctx.NewQuery().From(new(EventDown).TableName()).
			Where(dbox.And(filterX...)).Cursor(nil)

		defer csr.Close()

		countData := csr.Count()
		events := []*EventDown{}

		// do process here
		e = csr.Fetch(&events, 0, false)
		ErrorHandler(e, funcName)

		tk.Printf("ConvertEventReduceAvailability for %v | %v \n", t, countData)
		for _, d := range events {

			mtx.Lock()

			e = ctx.NewQuery().Update().From(new(EventDown).TableName()).
				Where(dbox.Eq("_id", d.ID)).
				Exec(tk.M{}.Set("data", tk.M{}.
					Set("reduceavailability", brakeReducesAvailability[d.AlarmDescription])))

			mtx.Unlock()
		}
		tk.Printf("end process for %v \n", t)

		csr.Close()
	}

	tk.Println("End process ConvertEventReduceAvailability...")
}
