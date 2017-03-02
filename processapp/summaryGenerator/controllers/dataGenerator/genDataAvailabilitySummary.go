package generatorControllers

import (
	. "eaciit/wfdemo-git/library/helper"
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
	tk.Println("===================== Start process Data Availability Summary...")

	availability := new(DataAvailability)
	availability.ID = "SCADA_DATA_OEM"
	availability.Name = "Scada Data OEM"

	availability = ev.scadaOEMSummary(availability)
	// mtx.Lock()
	e := ev.Ctx.Insert(availability)
	// mtx.Unlock()
	if e != nil {
		log.Printf("e: ", e.Error())
	}

	tk.Println("===================== End process Data Availability Summary...")
}

func (ev *DataAvailabilitySummary) scadaOEMSummary(availability *DataAvailability) *DataAvailability {
	var wg sync.WaitGroup

	ctx, e := PrepareConnection()
	if e != nil {
		log.Printf("e: %v \n", e.Error())
		os.Exit(0)
	}

	countx := 0
	maxPar := 1

	details := []DataAvailabilityDetail{}

	latestDate := time.Now().UTC()
	id := latestDate.Format("20060102_150405_SCADAOEM")

	availability.ID = id
	availability.Timestamp = latestDate
	projectName := "Tejuva"

	for turbine, _ := range ev.BaseController.RefTurbines {
		wg.Add(1)

		go func(t string, wgs *sync.WaitGroup, detail *[]DataAvailabilityDetail) {
			start := time.Now()

			filter := []*dbox.Filter{}
			filter = append(filter, dbox.Eq("projectname", projectName))
			filter = append(filter, dbox.Eq("turbine", t))
			// latestDate := ev.BaseController.GetLatest("ScadaDataOEM", "Tejuva", turbine)

			// latest 6 month
			startDate := latestDate.AddDate(0, -6, 0).UTC()
			if latestDate.Format("2006") != "0001" {
				filter = append(filter, dbox.Gt("timestamp", startDate))
			}

			csr, e := ctx.NewQuery().From(new(ScadaDataOEM).TableName()).
				Where(filter...).Order("timestamp").Cursor(nil)

			defer csr.Close()

			list := []ScadaDataOEM{}

			e = csr.Fetch(&list, 0, false)
			if e != nil {
				log.Printf("e: %v \n", e.Error())
			}

			// tData := csr.Count()

			// log.Printf("turbine: %v | %v \n", t, len(list))

			before := ScadaDataOEM{}
			from := ScadaDataOEM{}
			latestData := ScadaDataOEM{}
			hoursGap := 0.0

			for idx, oem := range list {
				if idx != 0 {
					before = list[idx-1]
					hoursGap = oem.TimeStamp.UTC().Sub(before.TimeStamp.UTC()).Hours()

					// log.Printf("xxx: %v - %v = %v \n", oem.TimeStamp.UTC().String(), before.TimeStamp.UTC().String(), hoursGap/24)

					if hoursGap > 24 {
						// log.Printf("hrs gap: %v \n", hoursGap)

						duration := before.TimeStamp.UTC().Sub(from.TimeStamp.UTC()).Hours() / 24
						details = append(details, setDataAvailDetail(from.TimeStamp, before.TimeStamp, projectName, turbine, duration, true))

						duration = oem.TimeStamp.UTC().Sub(before.TimeStamp.UTC()).Hours() / 24
						details = append(details, setDataAvailDetail(before.TimeStamp, oem.TimeStamp, projectName, turbine, duration, false))

						from = oem
					}
				} else {
					from = oem

					// set gap from stardate until first data in db
					hoursGap = from.TimeStamp.UTC().Sub(startDate.UTC()).Hours()
					if hoursGap > 24 {
						details = append(details, setDataAvailDetail(startDate, from.TimeStamp, projectName, turbine, hoursGap/24, false))
					}
				}

				latestData = oem
			}

			hoursGap = latestData.TimeStamp.UTC().Sub(from.TimeStamp.UTC()).Hours()

			if hoursGap > 24 {
				details = append(details, setDataAvailDetail(from.TimeStamp, latestData.TimeStamp, projectName, turbine, hoursGap/24, true))
			}

			hoursGap = latestDate.UTC().Sub(latestData.TimeStamp.UTC()).Hours()

			if hoursGap > 24 {
				details = append(details, setDataAvailDetail(latestData.TimeStamp, latestDate, projectName, turbine, hoursGap/24, false))
			}

			// log.Printf("xxx: %v - %v = %v \n", latestData.TimeStamp.UTC().String(), latestDate.UTC().String(), hoursGap)

			log.Printf(">> DONE: %v | %v secs", t, time.Now().Sub(start).Seconds())
			wgs.Done()
		}(turbine, &wg, &details)

		countx++

		if countx%maxPar == 0 || (len(ev.BaseController.RefTurbines) == countx) {
			// log.Println("wait....")
			wg.Wait()
		}
	}

	availability.Details = details
	// log.Printf("%#v \n", details)

	return availability
}

func setDataAvailDetail(from time.Time, to time.Time, turbine string, project string, duration float64, isAvail bool) DataAvailabilityDetail {

	res := DataAvailabilityDetail{
		ProjectName: project,
		Turbine:     turbine,
		Start:       from.UTC(),
		StartInfo:   GetDateInfo(from.UTC()),
		End:         to.UTC(),
		EndInfo:     GetDateInfo(to.UTC()),
		Duration:    duration,
		IsAvail:     isAvail,
	}

	return res
}
