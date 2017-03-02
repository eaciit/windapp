package generatorControllers

import (
	// . "eaciit/wfdemo-git/library/helper"
	. "eaciit/wfdemo-git/library/models"
	. "eaciit/wfdemo-git/processapp/summaryGenerator/controllers"
	"log"
	"os"
	"sync"

	"time"

	"github.com/eaciit/dbox"
	_ "github.com/eaciit/dbox/dbc/mongo"
	tk "github.com/eaciit/toolkit"
)

type DataAvailabilitySummary struct {
	*BaseController
}

func (ev *DataAvailabilitySummary) ConvertDataAvailabilitySummary(base *BaseController) {
	ev.BaseController = base
	tk.Println("Start process Data Availability Summary...")

	ev.scadaOEMSummary()

	tk.Println("End process Data Availability Summary...")
}

func (ev *DataAvailabilitySummary) scadaOEMSummary() {
	var wg sync.WaitGroup

	ctx, e := PrepareConnection()
	if e != nil {
		log.Printf("e: %v \n", e.Error())
		os.Exit(0)
	}

	countx := 0
	maxPar := 10

	for turbine, _ := range ev.BaseController.RefTurbines {
		wg.Add(1)

		go func(t string, wgs *sync.WaitGroup) {
			start := time.Now()

			filter := []*dbox.Filter{}
			filter = append(filter, dbox.Eq("projectname", "Tejuva"))
			filter = append(filter, dbox.Eq("turbine", t))
			log.Printf("turbine: %v", t)
			// latestDate := ev.BaseController.GetLatest("ScadaDataOEM", "Tejuva", turbine)

			latestDate := time.Now()
			id := latestDate.Format("20060102_150405_SCADAOEM")
			_ = id

			if latestDate.Format("2006") != "0001" {
				filter = append(filter, dbox.Gt("timestamp", latestDate.AddDate(0, -6, 0)))
			}

			csr, e := ctx.NewQuery().From(new(ScadaDataOEM).TableName()).
				Where(filter...).Cursor(nil)

			defer csr.Close()

			list := []ScadaDataOEM{}

			e = csr.Fetch(&list, 0, false)
			if e != nil {
				log.Printf("e: %v \n", e.Error())
			}

			// tData := csr.Count()

			for _, d := range list {
				mtx.Lock()
				dataInput := d
				_ = dataInput
				//tk.Printf("%s ", idx)

				// log.Printf("%v ", tc)

				// ev.doConversion(dataInput)
				// LogProcess("scadaOEMSummary."+turbine, float64(totalData), float64(idx))
				mtx.Unlock()
			}
			log.Printf(">> DONE: %v | %v secs", t, time.Now().Sub(start).Seconds())
			wgs.Done()
		}(turbine, &wg)

		countx++

		if countx%maxPar == 0 || (len(ev.BaseController.RefTurbines) == countx) {
			// log.Println("wait....")
			wg.Wait()
		}
	}
}
