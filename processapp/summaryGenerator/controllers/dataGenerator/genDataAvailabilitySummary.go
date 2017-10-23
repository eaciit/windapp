package generatorControllers

import (
	. "eaciit/wfdemo-git/library/helper"
	. "eaciit/wfdemo-git/library/models"
	. "eaciit/wfdemo-git/processapp/summaryGenerator/controllers"
	"os"
	"sync"

	"time"

	"github.com/eaciit/dbox"
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
	// OEM PROJECT
	availOEMProject := ev.scadaOEMSummaryProject()
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
	/* ============== PROJECT LEVEL ============== */
	availMet.ID = availMet.ID + "_PROJECT"
	availMet.Type = availMet.Type + "_PROJECT"
	availMet.Name = availMet.Name + " PROJECT"
	e = ev.Ctx.Insert(availMet)
	if e != nil {
		ev.Log.AddLog(tk.Sprintf("Found : %s"+e.Error()), sWarning)
	}
	e = ev.Ctx.Insert(availOEMProject)
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
		match := tk.M{}
		match.Set("projectname", value.GetString("project"))
		match.Set("turbine", turbine)
		match.Set("timestamp", tk.M{"$gte": periodFrom})
		go workerTurbine(turbine, value.GetString("project"), new(ScadaDataOEM).TableName(), match, &details, ev, ctx, &wg)

		countx++

		if countx%maxPar == 0 || (len(ev.BaseController.RefTurbines) == countx) {
			wg.Wait()
		}
	}

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
	match := tk.M{}

	for turbine, turbineVal := range ev.BaseController.RefTurbines {
		wg.Add(1)
		value, _ := tk.ToM(turbineVal)

		match.Set("projectname", value.GetString("project"))
		match.Set("turbine", turbine)
		match.Set("timestamp", tk.M{"$gte": periodFrom})
		match.Set("isnull", false)
		go workerTurbine(turbine, value.GetString("project"), "Scada10MinHFD", match, &details, ev, ctx, &wg)

		countx++

		if countx%maxPar == 0 || (len(ev.BaseController.RefTurbines) == countx) {
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
		match := tk.M{}
		match.Set("timestamp", tk.M{"$gte": periodFrom})
		go workerTurbine(turbine, value.GetString("project"), new(MetTower).TableName(), match, &details, ev, ctx, &wg)

		countx++

		if countx%maxPar == 0 || (len(ev.BaseController.RefTurbines) == countx) {
			wg.Wait()
		}
	}

	availability.Details = details

	return availability
}

func (ev *DataAvailabilitySummary) scadaOEMSummaryProject() *DataAvailability {
	tk.Println("===================== SCADA DATA OEM PROJECT LEVEL . . .")
	availability := new(DataAvailability)
	availability.Name = "Scada Data OEM PROJECT"
	availability.Type = "SCADA_DATA_OEM_PROJECT"

	var wg sync.WaitGroup

	ctx, e := PrepareConnection()
	if e != nil {
		ev.Log.AddLog(tk.Sprintf("Found : %s"+e.Error()), sWarning)
		os.Exit(0)
	}

	countx := 0

	details := []DataAvailabilityDetail{}

	now := time.Now().UTC()

	periodTo, _ := time.Parse("20060102_150405", now.Format("20060102_")+"000000")
	id := now.Format("20060102_150405_SCADAOEM_PROJECT")

	// latest 6 month
	periodFrom := GetNormalAddDateMonth(periodTo.UTC(), monthBefore)

	availability.ID = id
	availability.PeriodFrom = periodFrom
	availability.PeriodTo = periodTo

	for _, projectData := range ev.BaseController.ProjectList {
		projectID := projectData.ProjectId
		if projectID != "" {
			wg.Add(1)
			match := tk.M{}
			match.Set("projectname", projectID)
			match.Set("timestamp", tk.M{"$gte": periodFrom})
			workerProject(projectID, new(ScadaDataOEM).TableName(), match, &details, ev, ctx, &wg)
		}
		countx++
		if len(ev.BaseController.ProjectList) == countx {
			wg.Wait()
		}
	}

	availability.Details = details

	return availability
}

func (ev *DataAvailabilitySummary) scadaHFDSummaryProject() *DataAvailability {
	tk.Println("===================== SCADA DATA HFD PROJECT LEVEL . . .")
	availability := new(DataAvailability)
	availability.Name = "Scada Data HFD PROJECT"
	availability.Type = "SCADA_DATA_HFD_PROJECT"

	var wg sync.WaitGroup

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
		projectID := projectData.ProjectId
		if projectID != "" {
			wg.Add(1)
			match := tk.M{}
			match.Set("projectname", projectID)
			match.Set("timestamp", tk.M{"$gte": periodFrom})
			match.Set("isnull", false)
			workerProject(projectID, "Scada10MinHFD", match, &details, ev, ctx, &wg)
		}
		countx++
		if len(ev.BaseController.ProjectList) == countx {
			wg.Wait()
		}
	}

	availability.Details = details

	return availability
}

func workerTurbine(t, projectName, tablename string, match tk.M, details *[]DataAvailabilityDetail,
	ev *DataAvailabilitySummary, ctx dbox.IConnection, wg *sync.WaitGroup) {
	now := time.Now().UTC()
	periodTo, _ := time.Parse("20060102_150405", now.Format("20060102_")+"000000")
	periodFrom := GetNormalAddDateMonth(periodTo.UTC(), monthBefore)

	detail := []DataAvailabilityDetail{}
	start := time.Now()

	pipes := []tk.M{}
	pipes = append(pipes, tk.M{"$match": match})
	pipes = append(pipes, tk.M{"$project": tk.M{"timestamp": 1}})
	pipes = append(pipes, tk.M{"$sort": tk.M{"timestamp": 1}})

	csr, e := ctx.NewQuery().From(tablename).
		Command("pipe", pipes).Cursor(nil)

	countError := 0

	for {
		countError++
		if e != nil {
			csr, e = ctx.NewQuery().From(tablename).
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

	type CustomObject struct {
		TimeStamp time.Time
	}
	list := []CustomObject{}

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

	before := CustomObject{}
	from := CustomObject{}
	latestData := CustomObject{}
	hoursGap := 0.0
	duration := 0.0
	countID := 0

	if len(list) > 0 {
		for idx, oem := range list {
			if idx > 0 {
				before = list[idx-1]
				hoursGap = oem.TimeStamp.UTC().Sub(before.TimeStamp.UTC()).Hours()

				if hoursGap > 24 {
					countID++
					// set duration for available datas
					duration = tk.ToFloat64(before.TimeStamp.UTC().Sub(from.TimeStamp.UTC()).Hours()/24, 2, tk.RoundingAuto)
					mtx.Lock()
					*details = append(*details, setDataAvailDetail(from.TimeStamp, before.TimeStamp, projectName, t, duration, true, countID))
					mtx.Unlock()
					// set duration for unavailable datas
					countID++
					duration = tk.ToFloat64(hoursGap/24, 2, tk.RoundingAuto)
					mtx.Lock()
					*details = append(*details, setDataAvailDetail(before.TimeStamp, oem.TimeStamp, projectName, t, duration, false, countID))
					mtx.Unlock()
					from = oem
				}
			} else {
				from = oem

				// set gap from stardate until first data in db
				hoursGap = from.TimeStamp.UTC().Sub(periodFrom.UTC()).Hours()
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
			mtx.Lock()
			*details = append(*details, setDataAvailDetail(from.TimeStamp, latestData.TimeStamp, projectName, t, duration, true, countID))
			mtx.Unlock()
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
	mtx.Lock()
	*details = append(*details, detail...)
	ev.Log.AddLog(tk.Sprintf(">> DONE: %v | %v | %v secs \n", t, len(list), time.Now().Sub(start).Seconds()), sInfo)
	mtx.Unlock()
	// defer wg.Done()

	csr.Close()
	wg.Done()
}

func workerProject(projectName, tablename string, match tk.M, details *[]DataAvailabilityDetail,
	ev *DataAvailabilitySummary, ctx dbox.IConnection, wg *sync.WaitGroup) {
	now := time.Now().UTC()
	periodTo, _ := time.Parse("20060102_150405", now.Format("20060102_")+"000000")
	periodFrom := GetNormalAddDateMonth(periodTo.UTC(), monthBefore)

	detail := []DataAvailabilityDetail{}
	start := time.Now()

	pipes := []tk.M{}
	pipes = append(pipes, tk.M{"$match": match})
	pipes = append(pipes, tk.M{"$project": tk.M{"projectname": 1, "timestamp": 1}})
	pipes = append(pipes, tk.M{"$sort": tk.M{"timestamp": 1}})

	csr, e := ctx.NewQuery().From(tablename).
		Command("pipe", pipes).Cursor(nil)

	countError := 0

	for {
		countError++
		if e != nil {
			csr, e = ctx.NewQuery().From(tablename).
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
					mtx.Lock()
					*details = append(*details, setDataAvailDetail(from.TimeStamp, before.TimeStamp, projectName, "", duration, true, countID))
					mtx.Unlock()
					// set duration for unavailable datas
					countID++
					duration = tk.ToFloat64(hoursGap/24, 2, tk.RoundingAuto)
					mtx.Lock()
					*details = append(*details, setDataAvailDetail(before.TimeStamp, oem.TimeStamp, projectName, "", duration, false, countID))
					mtx.Unlock()
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
			mtx.Lock()
			*details = append(*details, setDataAvailDetail(from.TimeStamp, latestData.TimeStamp, projectName, "", duration, true, countID))
			mtx.Unlock()
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
	*details = append(*details, detail...)
	ev.Log.AddLog(tk.Sprintf(">> DONE: %v | %v | %v secs \n", projectName, len(list), time.Now().Sub(start).Seconds()), sInfo)
	mtx.Unlock()

	csr.Close()
	wg.Done()
}

func (ev *DataAvailabilitySummary) scadaHFDSummaryDaily() *DataAvailability {
	tk.Println("===================== SCADA DATA HFD DAILY...")
	var wg sync.WaitGroup

	ctx, e := PrepareConnection()
	if e != nil {
		ev.Log.AddLog(tk.Sprintf("Found : %s"+e.Error()), sWarning)
		os.Exit(0)
	}
	now := time.Now().UTC()

	periodTo, _ := time.Parse("20060102_150405", now.Format("20060102_")+"000000")
	id := now.Format("20060102_150405_SCADAHFD_DAILY")

	// latest 6 month
	periodFrom := GetNormalAddDateMonth(periodTo.UTC(), monthBefore)

	availabilityDaily := new(DataAvailability)
	availabilityDaily.Name = "Scada Data HFD_DAILY"
	availabilityDaily.Type = "SCADA_DATA_HFD_DAILY"
	availabilityDaily.ID = id
	availabilityDaily.PeriodFrom = periodFrom
	availabilityDaily.PeriodTo = periodTo

	matches := []tk.M{
		tk.M{"timestamp": tk.M{"$gte": periodFrom}},
		tk.M{"isnull": false},
	}
	groups := tk.M{
		"_id": tk.M{
			"projectname": "$projectname",
			"turbine":     "$turbine",
			"tanggal":     "$dateinfo.dateid",
		},
		"totaldata": tk.M{"$sum": 1},
	}
	pipes := []tk.M{}
	pipes = append(pipes, tk.M{"$match": tk.M{"$and": matches}})
	pipes = append(pipes, tk.M{"$group": groups})
	pipes = append(pipes, tk.M{"$project": tk.M{"projectname": 1, "timestamp": 1}})
	pipes = append(pipes, tk.M{"$sort": tk.M{"_id.tanggal": 1}})

	csr, e := ctx.NewQuery().From("Scada10MinHFD").
		Command("pipe", pipes).Cursor(nil)
	if e != nil {
		ev.Log.AddLog(tk.Sprintf("Found : %s"+e.Error()), sWarning)
	}

	_data := tk.M{}
	dataPerturbine := map[string][]tk.M{} /* for appending data per project per turbine */
	currProject := ""
	currTurbine := ""
	keys := ""
	for {
		_data = tk.M{}
		e = csr.Fetch(&_data, 1, false)
		if e != nil {
			break
		}
		currProject = _data.GetString("projectname")
		currTurbine = _data.GetString("turbine")
		keys = tk.Sprintf("%s_%s", currProject, currTurbine)
		dataPerturbine[keys] = append(dataPerturbine[keys], _data)
	}
	csr.Close()

	countx := 0
	maxPar := 5

	detailsDaily := []DataAvailabilityDetail{}

	for turbine, turbineVal := range ev.BaseController.RefTurbines {
		wg.Add(1)
		value, _ := tk.ToM(turbineVal)
		currProject = value.GetString("project")
		keys = tk.Sprintf("%s_%s", currProject, turbine)
		go workerDaily(dataPerturbine[keys], &detailsDaily, &wg)

		countx++

		if countx%maxPar == 0 || (len(ev.BaseController.RefTurbines) == countx) {
			wg.Wait()
		}
	}

	availabilityDaily.Details = detailsDaily
	ev.Log.AddLog(tk.Sprintf(">> DONE SCADA HFD DAILY"), sInfo)

	return availabilityDaily
}

func workerDaily(data []tk.M, details *[]DataAvailabilityDetail, wg *sync.WaitGroup) {
	detail := []DataAvailabilityDetail{}
	/*groups := tk.M{
		"_id": tk.M{
			"projectname": "$projectname",
			"turbine":     "$turbine",
			"tanggal":     "$dateinfo.dateid",
		},
		"totaldata": tk.M{"$sum": 1},
	}*/
	ids := tk.M{}
	maxDataPerDay := 6.0 * 24.0 /* dalam 1 jam ada 6 data karena per 10 menit*/
	duration := 0.0
	countID := 1
	periodFrom := time.Time{}
	periodTo := time.Time{}
	if len(data) > 0 {
		for _, datum := range data {
			ids = datum.Get("_id", tk.M{}).(tk.M)
			periodFrom = ids.Get("tanggal", time.Time{}).(time.Time)
			periodTo = periodFrom
			duration = tk.Div(datum.GetFloat64("totaldata"), maxDataPerDay)
			if duration >= 0.5 {
				detail = append(detail, setDataAvailDetail(periodFrom, periodTo, datum.GetString("projectname"),
					datum.GetString("turbine"), duration, true, countID))
			} else {
				detail = append(detail, setDataAvailDetail(periodFrom, periodTo, datum.GetString("projectname"),
					datum.GetString("turbine"), duration, false, countID))
			}
			countID++
		}
	}
	mtx.Lock()
	*details = append(*details, detail...)
	mtx.Unlock()

	wg.Done()
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
