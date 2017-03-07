package generatorControllers

import (
	. "eaciit/wfdemo-git/library/helper"
	. "eaciit/wfdemo-git/library/models"
	. "eaciit/wfdemo-git/processapp/summaryGenerator/controllers"
	"log"
	"os"
	"sync"

	"time"

	_ "github.com/eaciit/dbox/dbc/mongo"
	tk "github.com/eaciit/toolkit"
)

type DataAvailabilitySummary struct {
	*BaseController
}

func (ev *DataAvailabilitySummary) ConvertDataAvailabilitySummary(base *BaseController) {
	ev.BaseController = base
	tk.Println("===================== Start process Data Availability Summary...")

	// mtx.Lock()
	// OEM
	availOEM := ev.scadaOEMSummary()
	e := ev.Ctx.Insert(availOEM)
	if e != nil {
		log.Printf("e: ", e.Error())
	}

	// HFD
	availHFD := ev.scadaHFDSummary()
	e = ev.Ctx.Insert(availHFD)
	if e != nil {
		log.Printf("e: ", e.Error())
	}

	// mtx.Unlock()

	tk.Println("===================== End process Data Availability Summary...")
}

func (ev *DataAvailabilitySummary) scadaOEMSummary() *DataAvailability {
	availability := new(DataAvailability)
	availability.Name = "Scada Data OEM"
	availability.Type = "SCADA_DATA_OEM"

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
	availability.PeriodFrom = periodFrom
	availability.PeriodTo = periodTo
	projectName := "Tejuva"

	for turbine, _ := range ev.BaseController.RefTurbines {
		wg.Add(1)

		go func(t string) {
			detail := []DataAvailabilityDetail{}
			start := time.Now()

			match := tk.M{}
			match.Set("projectname", projectName)
			match.Set("turbine", t)
			match.Set("timestamp", tk.M{"$gte": periodFrom})

			pipes := []tk.M{}
			pipes = append(pipes, tk.M{"$match": match})
			pipes = append(pipes, tk.M{"$project": tk.M{"projectname": 1, "turbine": 1, "timestamp": 1}})
			pipes = append(pipes, tk.M{"$sort": tk.M{"timestamp": 1}})

			csr, e := ctx.NewQuery().From(new(ScadaDataOEM).TableName()).
				Command("pipe", pipes).Cursor(nil)

			countError := 0

			for {
				countError++
				if e != nil {
					csr, e = ctx.NewQuery().From(new(ScadaDataOEM).TableName()).
						Command("pipe", pipes).Cursor(nil)
					log.Printf("e: %v \n", e.Error())
				} else {
					break
				}

				if countError == 5 {
					break
				}
			}

			defer csr.Close()

			list := []ScadaDataOEM{}

			for {
				countError++
				e = csr.Fetch(&list, 0, false)
				if e != nil {
					log.Printf("e: %v \n", e.Error())
				} else {
					break
				}

				if countError == 5 {
					break
				}
			}

			before := ScadaDataOEM{}
			from := ScadaDataOEM{}
			latestData := ScadaDataOEM{}
			hoursGap := 0.0
			duration := 0.0
			countID := 0

			if len(list) > 0 {
				for idx, oem := range list {
					if idx > 0 {
						before = list[idx-1]
						hoursGap = oem.TimeStamp.UTC().Sub(before.TimeStamp.UTC()).Hours()
						// log.Printf("xxx: %v - %v = %v \n", oem.TimeStamp.UTC().String(), from.TimeStamp.UTC().String(), hoursGap/24)

						if hoursGap > 24 {
							countID++
							// log.Printf("hrs gap: %v \n", hoursGap)
							// set duration for available datas
							duration = tk.ToFloat64(before.TimeStamp.UTC().Sub(from.TimeStamp.UTC()).Hours()/24, 2, tk.RoundingAuto)
							details = append(details, setDataAvailDetail(from.TimeStamp, before.TimeStamp, projectName, t, duration, true, countID))
							// set duration for unavailable datas
							countID++
							duration = tk.ToFloat64(hoursGap/24, 2, tk.RoundingAuto)
							details = append(details, setDataAvailDetail(before.TimeStamp, oem.TimeStamp, projectName, t, duration, false, countID))
							from = oem
						}
					} else {
						from = oem

						// set gap from stardate until first data in db
						hoursGap = from.TimeStamp.UTC().Sub(periodFrom.UTC()).Hours()
						// log.Printf("idx=0 hrs gap: %v | %v | %v \n", hoursGap, from.TimeStamp.UTC().String(), periodFrom.UTC().String())
						if hoursGap > 24 {
							countID++
							duration = tk.ToFloat64(hoursGap/24, 2, tk.RoundingAuto)
							detail = append(detail, setDataAvailDetail(periodFrom, from.TimeStamp, projectName, t, duration, false, countID))
						}
					}

					latestData = oem
				}

				hoursGap = latestData.TimeStamp.UTC().Sub(from.TimeStamp.UTC()).Hours()

				if hoursGap > 24 {
					countID++
					duration = tk.ToFloat64(hoursGap/24, 2, tk.RoundingAuto)
					details = append(details, setDataAvailDetail(from.TimeStamp, latestData.TimeStamp, projectName, t, duration, true, countID))
				}

				// set gap from last data until periodTo
				hoursGap = periodTo.UTC().Sub(latestData.TimeStamp.UTC()).Hours()
				if hoursGap > 24 {
					countID++
					duration = tk.ToFloat64(hoursGap/24, 2, tk.RoundingAuto)
					detail = append(detail, setDataAvailDetail(latestData.TimeStamp, periodTo, projectName, t, duration, false, countID))
				}
			}
			if len(detail) == 0 {
				countID++
				duration = tk.ToFloat64(periodTo.Sub(periodFrom).Hours()/24, 2, tk.RoundingAuto)
				detail = append(detail, setDataAvailDetail(periodFrom, periodTo, projectName, t, duration, true, countID))
			}

			// log.Printf("xxx: %v - %v = %v \n", latestData.TimeStamp.UTC().String(), periodTo.UTC().String(), hoursGap)
			mtx.Lock()
			// for _, filt := range filter {
			// 	if filt.Field == "timestamp" {
			// 		log.Printf("timestamp: %#v \n", filt.Value.(time.Time).String())
			// 	} else {
			// 		log.Printf("%#v \n", filt)
			// 	}
			// }

			details = append(details, detail...)
			log.Printf(">> DONE: %v | %v | %v secs \n", t, len(list), time.Now().Sub(start).Seconds())
			mtx.Unlock()
			// defer wg.Done()

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

func (ev *DataAvailabilitySummary) scadaHFDSummary() *DataAvailability {
	availability := new(DataAvailability)
	availability.Name = "Scada Data HFD"
	availability.Type = "SCADA_DATA_HFD"

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
	id := now.Format("20060102_150405_SCADAHFD")

	// latest 6 month
	periodFrom := GetNormalAddDateMonth(periodTo.UTC(), -6)

	availability.ID = id
	availability.PeriodFrom = periodFrom
	availability.PeriodTo = periodTo
	projectName := "Tejuva"

	for turbine, _ := range ev.BaseController.RefTurbines {
		wg.Add(1)

		go func(t string) {
			detail := []DataAvailabilityDetail{}
			start := time.Now()

			match := tk.M{}
			match.Set("projectname", projectName)
			match.Set("turbine", t)
			match.Set("timestamp", tk.M{"$gte": periodFrom})
			match.Set("isnull", false)

			pipes := []tk.M{}
			pipes = append(pipes, tk.M{"$match": match})
			pipes = append(pipes, tk.M{"$project": tk.M{"projectname": 1, "turbine": 1, "timestamp": 1}})
			pipes = append(pipes, tk.M{"$sort": tk.M{"timestamp": 1}})

			csr, e := ctx.NewQuery().From(new(ScadaDataHFD).TableName()).
				Command("pipe", pipes).Cursor(nil)

			countError := 0

			for {
				countError++
				if e != nil {
					csr, e = ctx.NewQuery().From(new(ScadaDataOEM).TableName()).
						Command("pipe", pipes).Cursor(nil)
					log.Printf("e: %v \n", e.Error())
				} else {
					break
				}

				if countError == 5 {
					break
				}
			}

			defer csr.Close()

			list := []ScadaDataOEM{}

			for {
				countError++
				e = csr.Fetch(&list, 0, false)
				if e != nil {
					log.Printf("e: %v \n", e.Error())
				} else {
					break
				}

				if countError == 5 {
					break
				}
			}

			before := ScadaDataOEM{}
			from := ScadaDataOEM{}
			latestData := ScadaDataOEM{}
			hoursGap := 0.0
			duration := 0.0
			countID := 0

			if len(list) > 0 {
				for idx, oem := range list {
					if idx > 0 {
						before = list[idx-1]
						hoursGap = oem.TimeStamp.UTC().Sub(before.TimeStamp.UTC()).Hours()
						// log.Printf("xxx: %v - %v = %v \n", oem.TimeStamp.UTC().String(), from.TimeStamp.UTC().String(), hoursGap/24)

						if hoursGap > 24 {
							countID++
							// log.Printf("hrs gap: %v \n", hoursGap)
							// set duration for available datas
							duration = tk.ToFloat64(before.TimeStamp.UTC().Sub(from.TimeStamp.UTC()).Hours()/24, 2, tk.RoundingAuto)
							details = append(details, setDataAvailDetail(from.TimeStamp, before.TimeStamp, projectName, t, duration, true, countID))
							// set duration for unavailable datas
							countID++
							duration = tk.ToFloat64(hoursGap/24, 2, tk.RoundingAuto)
							details = append(details, setDataAvailDetail(before.TimeStamp, oem.TimeStamp, projectName, t, duration, false, countID))
							from = oem
						}
					} else {
						from = oem

						// set gap from stardate until first data in db
						hoursGap = from.TimeStamp.UTC().Sub(periodFrom.UTC()).Hours()
						// log.Printf("idx=0 hrs gap: %v | %v | %v \n", hoursGap, from.TimeStamp.UTC().String(), periodFrom.UTC().String())
						if hoursGap > 24 {
							countID++
							duration = tk.ToFloat64(hoursGap/24, 2, tk.RoundingAuto)
							detail = append(detail, setDataAvailDetail(periodFrom, from.TimeStamp, projectName, t, duration, false, countID))
						}
					}

					latestData = oem
				}

				hoursGap = latestData.TimeStamp.UTC().Sub(from.TimeStamp.UTC()).Hours()

				if hoursGap > 24 {
					countID++
					duration = tk.ToFloat64(hoursGap/24, 2, tk.RoundingAuto)
					details = append(details, setDataAvailDetail(from.TimeStamp, latestData.TimeStamp, projectName, t, duration, true, countID))
				}

				// set gap from last data until periodTo
				hoursGap = periodTo.UTC().Sub(latestData.TimeStamp.UTC()).Hours()
				if hoursGap > 24 {
					countID++
					duration = tk.ToFloat64(hoursGap/24, 2, tk.RoundingAuto)
					detail = append(detail, setDataAvailDetail(latestData.TimeStamp, periodTo, projectName, t, duration, false, countID))
				}
			}
			if len(detail) == 0 {
				countID++
				duration = tk.ToFloat64(periodTo.Sub(periodFrom).Hours()/24, 2, tk.RoundingAuto)
				detail = append(detail, setDataAvailDetail(periodFrom, periodTo, projectName, t, duration, true, countID))
			}

			// log.Printf("xxx: %v - %v = %v \n", latestData.TimeStamp.UTC().String(), periodTo.UTC().String(), hoursGap)
			mtx.Lock()
			// for _, filt := range filter {
			// 	if filt.Field == "timestamp" {
			// 		log.Printf("timestamp: %#v \n", filt.Value.(time.Time).String())
			// 	} else {
			// 		log.Printf("%#v \n", filt)
			// 	}
			// }

			details = append(details, detail...)
			log.Printf(">> DONE: %v | %v | %v secs \n", t, len(list), time.Now().Sub(start).Seconds())
			mtx.Unlock()
			// defer wg.Done()

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

func setDataAvailDetail(from time.Time, to time.Time, project string, turbine string, duration float64, isAvail bool, id int) DataAvailabilityDetail {

	res := DataAvailabilityDetail{
		ID:          id,
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
