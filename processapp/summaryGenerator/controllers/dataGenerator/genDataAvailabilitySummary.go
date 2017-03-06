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
	availability.Name = "Scada Data OEM"
	availability.Type = "SCADA_DATA_OEM"

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

	// latest 6 month
	periodFrom := GetNormalAddDateMonth(periodTo.UTC(), -6)

	availability.ID = id
	availability.PeriodTo = periodTo
	availability.PeriodFrom = periodFrom
	projectName := "Tejuva"

	for turbine, _ := range ev.BaseController.RefTurbines {
		wg.Add(1)

		go func(t string) {
			detail := []DataAvailabilityDetail{}
			start := time.Now()

			filter := []*dbox.Filter{}
			filter = append(filter, dbox.Eq("projectname", projectName))
			filter = append(filter, dbox.Eq("turbine", t))
			filter = append(filter, dbox.Gte("timestamp", periodFrom))

			csr, e := ctx.NewQuery().From(new(ScadaDataOEM).TableName()).
				Where(filter...).Order("timestamp").Cursor(nil)

			defer csr.Close()

			list := []ScadaDataOEM{}

			for {
				e = csr.Fetch(&list, 0, false)
				if e != nil {
					log.Printf("e: %v \n", e.Error())
				} else {
					break
				}
			}

			before := ScadaDataOEM{}
			from := ScadaDataOEM{}
			latestData := ScadaDataOEM{}
			hoursGap := 0.0
			duration := 0.0

			for idx, oem := range list {
				if idx > 0 {
					before = list[idx-1]
					hoursGap = oem.TimeStamp.UTC().Sub(before.TimeStamp.UTC()).Hours()
					// log.Printf("xxx: %v - %v = %v \n", oem.TimeStamp.UTC().String(), from.TimeStamp.UTC().String(), hoursGap/24)

					if hoursGap > 24 {
						// log.Printf("hrs gap: %v \n", hoursGap)
						// set duration for available datas
						duration = before.TimeStamp.UTC().Sub(from.TimeStamp.UTC()).Hours() / 24
						details = append(details, setDataAvailDetail(from.TimeStamp, before.TimeStamp, projectName, t, duration, true))
						// set duration for unavailable datas
						duration = hoursGap / 24
						details = append(details, setDataAvailDetail(before.TimeStamp, oem.TimeStamp, projectName, t, duration, false))
						from = oem
					}
				} else {
					from = oem

					// set gap from stardate until first data in db
					hoursGap = from.TimeStamp.UTC().Sub(periodFrom.UTC()).Hours()
					// log.Printf("idx=0 hrs gap: %v \n", hoursGap)
					if hoursGap > 24 {
						duration = hoursGap / 24
						detail = append(detail, setDataAvailDetail(periodFrom, from.TimeStamp, projectName, t, duration, false))
					}
				}

				latestData = oem
			}

			hoursGap = latestData.TimeStamp.UTC().Sub(from.TimeStamp.UTC()).Hours()

			if hoursGap > 24 {
				details = append(details, setDataAvailDetail(from.TimeStamp, latestData.TimeStamp, projectName, t, hoursGap/24, true))
			}

			// set gap from last data until periodTo
			hoursGap = periodTo.UTC().Sub(latestData.TimeStamp.UTC()).Hours()
			if hoursGap > 24 {
				duration = hoursGap / 24
				detail = append(detail, setDataAvailDetail(latestData.TimeStamp, periodTo, projectName, t, duration, false))
			}

			if len(detail) == 0 {
				duration = periodTo.Sub(periodFrom).Hours() / 24
				detail = append(detail, setDataAvailDetail(periodFrom, periodTo, projectName, t, duration, true))
			}

			// log.Printf("xxx: %v - %v = %v \n", latestData.TimeStamp.UTC().String(), periodTo.UTC().String(), hoursGap)
			mtx.Lock()
			for _, filt := range filter {
				if filt.Field == "timestamp" {
					log.Printf("timestamp: %#v \n", filt.Value.(time.Time).String())
				} else {
					log.Printf("%#v \n", filt)
				}
			}

			details = append(details, detail...)
			log.Printf(">> DONE: %v | %v | %v secs >> %v \n", t, len(list), time.Now().Sub(start).Seconds(), len(detail))
			mtx.Unlock()
			defer wg.Done()

			csr.Close()
			wg.Done()
		}(turbine)

		countx++

		if countx%maxPar == 0 || (len(ev.BaseController.RefTurbines) == countx) {
			// log.Println("wait....")
			wg.Wait()
		}
	}

	availability.Details = details

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
