package generatorControllers

import (
	. "eaciit/wfdemo-git/library/helper"
	. "eaciit/wfdemo-git/library/models"
	. "eaciit/wfdemo-git/processapp/controllers"
	"os"
	"sync"
	"time"

	"github.com/eaciit/dbox"
	_ "github.com/eaciit/dbox/dbc/mongo"
	tk "github.com/eaciit/toolkit"
)

type UpdateScadaHFD struct {
	*BaseController
}

var (
	hmtx = &sync.Mutex{}
)

func (c *UpdateScadaHFD) DoUpdateWsBin(base *BaseController) {
	funcName := "Update ws bins for ScadaHFD"
	c.BaseController = base

	var wg sync.WaitGroup

	if base != nil {
		ctx, e := PrepareConnection()
		if e != nil {
			ErrorHandler(e, "Scada Summary")
			os.Exit(0)
		}

		csr, e := ctx.NewQuery().From(new(ScadaDataHFD).TableName()).Cursor(nil)

		defer csr.Close()

		counter := 0
		countData := csr.Count()
		isDone := false
		countPerProcess := 1000

		for !isDone && countData > 0 {
			scadas := []*ScadaDataHFD{}

			e = csr.Fetch(&scadas, countPerProcess, false)
			ErrorHandler(e, funcName)

			if len(scadas) < countPerProcess {
				isDone = true
			}

			wg.Add(1)
			go func(datas []*ScadaDataHFD, counter int) {
				tk.Println("start process ", countPerProcess*(counter+1))
				for _, d := range datas {
					hmtx.Lock()

					dId := d.ID
					wsBin := tk.RoundingAuto64(d.Fast_WindSpeed_ms, 0)
					tk.Println("Updating data for ID = ", dId, wsBin)
					e = ctx.NewQuery().Update().From(new(ScadaDataHFD).TableName()).
						Where(dbox.Eq("_id", dId)).
						Exec(tk.M{}.Set("data", tk.M{}.Set("fast_windspeed_bin", wsBin)))
					ErrorHandler(e, funcName)

					hmtx.Unlock()
				}
				tk.Println("end process ", countPerProcess*(counter+1))
				wg.Done()
			}(scadas, counter)

			counter++
			if counter%10 == 0 || isDone {
				wg.Wait()
			}
		}
	}

	tk.Println("End process updating wind speed bin for ScadaData HFD...")
}

func (c *UpdateScadaHFD) DoUpdateTurbineStateForScadaData(base *BaseController) {
	funcName := "Update turbine state for Scada"
	c.BaseController = base

	startTime := time.Now()

	// var wg sync.WaitGroup

	if base != nil {
		ctx, e := PrepareConnection()
		if e != nil {
			ErrorHandler(e, "Update Turbine State Scada Data")
			os.Exit(0)
		}

		csr, e := ctx.NewQuery().From("Scada10MinHFD").Cursor(nil)

		defer csr.Close()

		// counter := 0
		countData := csr.Count()
		isDone := false
		countPerProcess := 1000

		for !isDone && countData > 0 {
			scadas := []*ScadaDataHFD{}

			startTimeFetch := time.Now()

			e = csr.Fetch(&scadas, countPerProcess, false)
			ErrorHandler(e, funcName)

			if len(scadas) < countPerProcess {
				isDone = true
			}

			durationFetch := time.Now().Sub(startTimeFetch).Seconds()
			tk.Printf("Fetch %v data about %v sec(s)", len(scadas), durationFetch)

			// wg.Add(1)
			// go func(datas []*ScadaDataHFD, counter int) {
			// 	startTimeGo := time.Now()
			// 	tk.Println("start process ", countPerProcess*(counter+1))
			// 	for _, d := range datas {
			// 		hmtx.Lock()

			// 		dId := d.ID
			// 		turbineState := d.TurbineState
			// 		tk.Println("Updating data for ID = ", dId, turbineState)
			// 		// e = ctx.NewQuery().Update().From(new(ScadaDataHFD).TableName()).
			// 		// 	Where(dbox.Eq("_id", dId)).
			// 		// 	Exec(tk.M{}.Set("data", tk.M{}.Set("fast_windspeed_bin", wsBin)))
			// 		// ErrorHandler(e, funcName)

			// 		hmtx.Unlock()
			// 	}
			// 	durationGo := time.Now().Sub(startTimeGo).Seconds()
			// 	tk.Printf("end process %v data about %v sec(s) \n", countPerProcess*(counter+1), durationGo)
			// 	wg.Done()
			// }(scadas, counter)

			// counter++
			// if counter%10 == 0 || isDone {
			// 	wg.Wait()
			// }
		}
	}

	duration := time.Now().Sub(startTime).Seconds()

	tk.Printf("End process updating wind speed bin for ScadaData HFD about %v sec(s)\n", duration)
}
