package generatorControllers

import (
	. "eaciit/wfdemo-git/library/helper"
	. "eaciit/wfdemo-git/library/models"
	. "eaciit/wfdemo-git/processapp/summaryGenerator/controllers"
	"os"
	"sync"

	"time"

	_ "github.com/eaciit/dbox/dbc/mongo"
	tk "github.com/eaciit/toolkit"
)

const (
	monthBefore = -5
)

var (
// projectName = "Tejuva"
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
	// HFD
	availHFD := ev.scadaHFDSummary()
	// Met Tower
	availMet := ev.metTowerSummary()
	// HFD PROJECT
	availHFDProject := ev.scadaHFDSummaryProject()

	ev.Ctx.DeleteMany(new(DataAvailability), nil)

	e := ev.Ctx.Insert(availOEM)
	if e != nil {
		ev.Log.AddLog(tk.Sprintf("Found : %s"+e.Error()), sWarning)
	}
	e = ev.Ctx.Insert(availHFD)
	if e != nil {
		ev.Log.AddLog(tk.Sprintf("Found : %s"+e.Error()), sWarning)
	}
	e = ev.Ctx.Insert(availMet)
	if e != nil {
		ev.Log.AddLog(tk.Sprintf("Found : %s"+e.Error()), sWarning)
	}
	e = ev.Ctx.Insert(availHFDProject)
	if e != nil {
		ev.Log.AddLog(tk.Sprintf("Found : %s"+e.Error()), sWarning)
	}

	// mtx.Unlock()

	tk.Println("===================== End process Data Availability Summary...")
}

func (ev *DataAvailabilitySummary) scadaOEMSummary() *DataAvailability {
	tk.Println("===================== SCADA DATA OEM...")
	availability := new(DataAvailability)
	availability.Name = "Scada Data OEM"
	availability.Type = "SCADA_DATA_OEM"

	var wg sync.WaitGroup

	ctx, e := PrepareConnection()
	if e != nil {
		ev.Log.AddLog(tk.Sprintf("Found : %s"+e.Error()), sWarning)
		os.Exit(0)
	}

	countx := 0
	maxPar := 5

	details := []DataAvailabilityDetail{}

	now := time.Now().UTC()

	periodTo, _ := time.Parse("20060102_150405", now.Format("20060102_")+"000000")
	id := now.Format("20060102_150405_SCADAOEM")

	// latest 6 month
	periodFrom := GetNormalAddDateMonth(periodTo.UTC(), monthBefore)

	availability.ID = id
	availability.PeriodFrom = periodFrom
	availability.PeriodTo = periodTo

	for turbine, turbineVal := range ev.BaseController.RefTurbines {
		wg.Add(1)
		value, _ := tk.ToM(turbineVal)

		go func(t string, projectName string) {
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
					ev.Log.AddLog(tk.Sprintf("Found : %s"+e.Error()), sWarning)
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
					ev.Log.AddLog(tk.Sprintf("Found : %s"+e.Error()), sWarning)
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
				if len(detail) == 0 {
					countID++
					duration = tk.ToFloat64(periodTo.Sub(periodFrom).Hours()/24, 2, tk.RoundingAuto)
					detail = append(detail, setDataAvailDetail(periodFrom, periodTo, projectName, t, duration, true, countID))
				}
			} else {
				countID++
				duration = tk.ToFloat64(periodTo.Sub(periodFrom).Hours()/24, 2, tk.RoundingAuto)
				detail = append(detail, setDataAvailDetail(periodFrom, periodTo, projectName, t, duration, false, countID))
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
			ev.Log.AddLog(tk.Sprintf(">> DONE: %v | %v | %v secs \n", t, len(list), time.Now().Sub(start).Seconds()), sInfo)
			mtx.Unlock()
			// defer wg.Done()

			csr.Close()
			wg.Done()
		}(turbine, value.GetString("project"))

		countx++

		if countx%maxPar == 0 || (len(ev.BaseController.RefTurbines) == countx) {
			// log.Println("wait....")
			wg.Wait()
		}
	}

	availability.Details = details

	return availability
}

func (ev *DataAvailabilitySummary) scadaHFDSummaryProject() *DataAvailability {
	tk.Println("===================== SCADA DATA HFD...")
	availability := new(DataAvailability)
	availability.Name = "Scada Data HFD_PROJECT"
	availability.Type = "SCADA_DATA_HFD_PROJECT"

	ctx, e := PrepareConnection()
	if e != nil {
		ev.Log.AddLog(tk.Sprintf("Found : %s"+e.Error()), sWarning)
		os.Exit(0)
	}

	countx := 0

	details := []DataAvailabilityDetail{}

	now := time.Now().UTC()

	periodTo, _ := time.Parse("20060102_150405", now.Format("20060102_")+"000000")
	id := now.Format("20060102_150405_SCADAHFD_PROJECT")

	// latest 6 month
	periodFrom := GetNormalAddDateMonth(periodTo.UTC(), monthBefore)

	availability.ID = id
	availability.PeriodFrom = periodFrom
	availability.PeriodTo = periodTo

	for _, projectData := range ev.BaseController.ProjectList {
		projectName := projectData.ProjectId
		detail := []DataAvailabilityDetail{}
		start := time.Now()

		match := tk.M{}
		match.Set("projectname", projectName)
		match.Set("timestamp", tk.M{"$gte": periodFrom})
		match.Set("isnull", false)

		pipes := []tk.M{}
		pipes = append(pipes, tk.M{"$match": match})
		pipes = append(pipes, tk.M{"$project": tk.M{"projectname": 1, "timestamp": 1}})
		pipes = append(pipes, tk.M{"$sort": tk.M{"timestamp": 1}})

		csr, e := ctx.NewQuery().From("Scada10MinHFD").
			Command("pipe", pipes).Cursor(nil)

		countError := 0

		for {
			countError++
			if e != nil {
				csr, e = ctx.NewQuery().From("Scada10MinHFD").
					Command("pipe", pipes).Cursor(nil)
				ev.Log.AddLog(tk.Sprintf("Found : %s"+e.Error()), sWarning)
			} else {
				break
			}

			if countError == 5 {
				break
			}
		}

		defer csr.Close()

		type Scada10MinHFDCustom struct {
			TimeStamp time.Time
		}
		list := []Scada10MinHFDCustom{}

		for {
			countError++
			e = csr.Fetch(&list, 0, false)
			if e != nil {
				ev.Log.AddLog(tk.Sprintf("Found : %s"+e.Error()), sWarning)
			} else {
				break
			}

			if countError == 5 {
				break
			}
		}

		before := Scada10MinHFDCustom{}
		from := Scada10MinHFDCustom{}
		latestData := Scada10MinHFDCustom{}
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
						details = append(details, setDataAvailDetail(from.TimeStamp, before.TimeStamp, projectName, "", duration, true, countID))
						// set duration for unavailable datas
						countID++
						duration = tk.ToFloat64(hoursGap/24, 2, tk.RoundingAuto)
						details = append(details, setDataAvailDetail(before.TimeStamp, oem.TimeStamp, projectName, "", duration, false, countID))
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
						detail = append(detail, setDataAvailDetail(periodFrom, from.TimeStamp, projectName, "", duration, false, countID))
					}
				}

				latestData = oem
			}

			hoursGap = latestData.TimeStamp.UTC().Sub(from.TimeStamp.UTC()).Hours()

			if hoursGap > 24 {
				countID++
				duration = tk.ToFloat64(hoursGap/24, 2, tk.RoundingAuto)
				details = append(details, setDataAvailDetail(from.TimeStamp, latestData.TimeStamp, projectName, "", duration, true, countID))
			}

			// set gap from last data until periodTo
			hoursGap = periodTo.UTC().Sub(latestData.TimeStamp.UTC()).Hours()
			if hoursGap > 24 {
				countID++
				duration = tk.ToFloat64(hoursGap/24, 2, tk.RoundingAuto)
				detail = append(detail, setDataAvailDetail(latestData.TimeStamp, periodTo, projectName, "", duration, false, countID))
			}

			if len(detail) == 0 {
				countID++
				duration = tk.ToFloat64(periodTo.Sub(periodFrom).Hours()/24, 2, tk.RoundingAuto)
				detail = append(detail, setDataAvailDetail(periodFrom, periodTo, projectName, "", duration, true, countID))
			}
		} else {
			countID++
			duration = tk.ToFloat64(periodTo.Sub(periodFrom).Hours()/24, 2, tk.RoundingAuto)
			detail = append(detail, setDataAvailDetail(periodFrom, periodTo, projectName, "", duration, false, countID))
		}

		mtx.Lock()
		details = append(details, detail...)
		ev.Log.AddLog(tk.Sprintf(">> DONE: %v | %v | %v secs \n", projectName, len(list), time.Now().Sub(start).Seconds()), sInfo)
		mtx.Unlock()

		csr.Close()
	}

	countx++

	availability.Details = details

	return availability
}

func (ev *DataAvailabilitySummary) scadaHFDSummary() *DataAvailability {
	tk.Println("===================== SCADA DATA HFD...")
	availability := new(DataAvailability)
	availability.Name = "Scada Data HFD"
	availability.Type = "SCADA_DATA_HFD"

	var wg sync.WaitGroup

	ctx, e := PrepareConnection()
	if e != nil {
		ev.Log.AddLog(tk.Sprintf("Found : %s"+e.Error()), sWarning)
		os.Exit(0)
	}

	countx := 0
	maxPar := 5

	details := []DataAvailabilityDetail{}

	now := time.Now().UTC()

	periodTo, _ := time.Parse("20060102_150405", now.Format("20060102_")+"000000")
	id := now.Format("20060102_150405_SCADAHFD")

	// latest 6 month
	periodFrom := GetNormalAddDateMonth(periodTo.UTC(), monthBefore)

	availability.ID = id
	availability.PeriodFrom = periodFrom
	availability.PeriodTo = periodTo

	for turbine, turbineVal := range ev.BaseController.RefTurbines {
		wg.Add(1)
		value, _ := tk.ToM(turbineVal)

		go func(t string, projectName string) {
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

			csr, e := ctx.NewQuery().From("Scada10MinHFD").
				Command("pipe", pipes).Cursor(nil)

			countError := 0

			for {
				countError++
				if e != nil {
					csr, e = ctx.NewQuery().From("Scada10MinHFD").
						Command("pipe", pipes).Cursor(nil)
					ev.Log.AddLog(tk.Sprintf("Found : %s"+e.Error()), sWarning)
				} else {
					break
				}

				if countError == 5 {
					break
				}
			}

			defer csr.Close()

			type Scada10MinHFDCustom struct {
				TimeStamp time.Time
			}
			list := []Scada10MinHFDCustom{}

			for {
				countError++
				e = csr.Fetch(&list, 0, false)
				if e != nil {
					ev.Log.AddLog(tk.Sprintf("Found : %s"+e.Error()), sWarning)
				} else {
					break
				}

				if countError == 5 {
					break
				}
			}

			before := Scada10MinHFDCustom{}
			from := Scada10MinHFDCustom{}
			latestData := Scada10MinHFDCustom{}
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

				if len(detail) == 0 {
					countID++
					duration = tk.ToFloat64(periodTo.Sub(periodFrom).Hours()/24, 2, tk.RoundingAuto)
					detail = append(detail, setDataAvailDetail(periodFrom, periodTo, projectName, t, duration, true, countID))
				}
			} else {
				countID++
				duration = tk.ToFloat64(periodTo.Sub(periodFrom).Hours()/24, 2, tk.RoundingAuto)
				detail = append(detail, setDataAvailDetail(periodFrom, periodTo, projectName, t, duration, false, countID))
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
			ev.Log.AddLog(tk.Sprintf(">> DONE: %v | %v | %v secs \n", t, len(list), time.Now().Sub(start).Seconds()), sInfo)
			mtx.Unlock()
			// defer wg.Done()

			csr.Close()
			wg.Done()
		}(turbine, value.GetString("project"))

		countx++

		if countx%maxPar == 0 || (len(ev.BaseController.RefTurbines) == countx) {
			// log.Println("wait....")
			wg.Wait()
		}
	}

	availability.Details = details

	return availability
}

func (ev *DataAvailabilitySummary) metTowerSummary() *DataAvailability {
	tk.Println("===================== MET TOWER...")
	availability := new(DataAvailability)
	availability.Name = "Met Tower"
	availability.Type = "MET_TOWER"

	var wg sync.WaitGroup

	ctx, e := PrepareConnection()
	if e != nil {
		ev.Log.AddLog(tk.Sprintf("Found : %s"+e.Error()), sWarning)
		os.Exit(0)
	}

	countx := 0
	maxPar := 5

	details := []DataAvailabilityDetail{}

	now := time.Now().UTC()

	periodTo, _ := time.Parse("20060102_150405", now.Format("20060102_")+"000000")
	id := now.Format("20060102_150405_METTOWER")

	// latest 6 month
	periodFrom := GetNormalAddDateMonth(periodTo.UTC(), monthBefore)

	availability.ID = id
	availability.PeriodFrom = periodFrom
	availability.PeriodTo = periodTo

	for turbine, turbineVal := range ev.BaseController.RefTurbines {
		wg.Add(1)
		value, _ := tk.ToM(turbineVal)

		go func(t string, projectName string) {
			detail := []DataAvailabilityDetail{}
			start := time.Now()

			match := tk.M{}
			// match.Set("projectname", projectName)
			// match.Set("turbine", t)
			match.Set("timestamp", tk.M{"$gte": periodFrom})

			pipes := []tk.M{}
			pipes = append(pipes, tk.M{"$match": match})
			pipes = append(pipes, tk.M{"$project": tk.M{"timestamp": 1}})
			pipes = append(pipes, tk.M{"$sort": tk.M{"timestamp": 1}})

			csr, e := ctx.NewQuery().From(new(MetTower).TableName()).
				Command("pipe", pipes).Cursor(nil)

			countError := 0

			for {
				countError++
				if e != nil {
					csr, e = ctx.NewQuery().From(new(MetTower).TableName()).
						Command("pipe", pipes).Cursor(nil)
					ev.Log.AddLog(tk.Sprintf("Found : %s"+e.Error()), sWarning)
				} else {
					break
				}

				if countError == 5 {
					break
				}
			}

			defer csr.Close()

			list := []MetTower{}

			for {
				countError++
				e = csr.Fetch(&list, 0, false)
				if e != nil {
					ev.Log.AddLog(tk.Sprintf("Found : %s"+e.Error()), sWarning)
				} else {
					break
				}

				if countError == 5 {
					break
				}
			}

			before := MetTower{}
			from := MetTower{}
			latestData := MetTower{}
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
				if len(detail) == 0 {
					countID++
					duration = tk.ToFloat64(periodTo.Sub(periodFrom).Hours()/24, 2, tk.RoundingAuto)
					detail = append(detail, setDataAvailDetail(periodFrom, periodTo, projectName, t, duration, true, countID))
				}
			} else {
				countID++
				duration = tk.ToFloat64(periodTo.Sub(periodFrom).Hours()/24, 2, tk.RoundingAuto)
				detail = append(detail, setDataAvailDetail(periodFrom, periodTo, projectName, t, duration, false, countID))
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
			ev.Log.AddLog(tk.Sprintf(">> DONE: %v | %v | %v secs \n", t, len(list), time.Now().Sub(start).Seconds()), sInfo)
			mtx.Unlock()
			// defer wg.Done()

			csr.Close()
			wg.Done()
		}(turbine, value.GetString("project"))

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
