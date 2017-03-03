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
	maxPar := 5

	details := []DataAvailabilityDetail{}

	now := time.Now().UTC()

	periodTo, _ := time.Parse("20060102_150405", now.Format("20060102_")+"000000")
	id := now.Format("20060102_150405_SCADAOEM")

	periodFrom := GetNormalAddDateMonth(periodTo.UTC(), -6)

	availability.ID = id
	availability.PeriodTo = periodTo
	availability.PeriodFrom = periodFrom
	projectName := "Tejuva"

	for turbine, _ := range ev.BaseController.RefTurbines {
		wg.Add(1)

		go func(t string, wgs *sync.WaitGroup, detail *[]DataAvailabilityDetail) {
			start := time.Now()

			filter := []*dbox.Filter{}
			filter = append(filter, dbox.Eq("projectname", projectName))
			filter = append(filter, dbox.Eq("turbine", t))
			// latest 6 month
			// periodFrom := periodTo.AddDate(0, -6, 0).UTC()

			// if periodTo.Format("2006") != "0001" {
			filter = append(filter, dbox.Gte("timestamp", periodFrom))
			// }

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
					hoursGap = from.TimeStamp.UTC().Sub(periodFrom.UTC()).Hours()
					if hoursGap > 24 {
						details = append(details, setDataAvailDetail(periodFrom, from.TimeStamp, projectName, turbine, hoursGap/24, false))
					}
				}

				latestData = oem
			}

			hoursGap = latestData.TimeStamp.UTC().Sub(from.TimeStamp.UTC()).Hours()

			if hoursGap > 24 {
				details = append(details, setDataAvailDetail(from.TimeStamp, latestData.TimeStamp, projectName, turbine, hoursGap/24, true))
			}

			hoursGap = periodTo.UTC().Sub(latestData.TimeStamp.UTC()).Hours()

			if hoursGap > 24 {
				details = append(details, setDataAvailDetail(latestData.TimeStamp, periodTo, projectName, turbine, hoursGap/24, false))
			}

			// log.Printf("xxx: %v - %v = %v \n", latestData.TimeStamp.UTC().String(), periodTo.UTC().String(), hoursGap)

			log.Printf(">> DONE: %v | %v | %v secs", t, len(list), time.Now().Sub(start).Seconds())
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

func setDataAvailDetail(from time.Time, to time.Time, project string, turbine string, duration float64, isAvail bool) DataAvailabilityDetail {

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
