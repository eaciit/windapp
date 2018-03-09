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
	turbineName map[string]string
)

type DataAvailabilitySummary struct {
	*BaseController
}

func (ev *DataAvailabilitySummary) ConvertDataAvailabilitySummary(base *BaseController) {
	ev.BaseController = base
	ev.Log.AddLog("===================== Start process Data Availability Summary...", sInfo)

	turbineName = map[string]string{}
	var e error

	turbineName, e = GetTurbineNameListAll("")
	if e != nil {
		ev.Log.AddLog(tk.Sprintf("Found : %s"+e.Error()), sError)
	}

	var wg sync.WaitGroup
	wg.Add(5)

	availOEM := new(DataAvailability)
	availHFD := new(DataAvailability)
	availMet := new(DataAvailability)
	availOEMProject := new(DataAvailability)
	availHFDProject := new(DataAvailability)
	// availHFDDaily := new(DataAvailability)

	// OEM
	go func() {
		defer wg.Done()
		availOEM = ev.scadaOEMSummary()
	}()
	// HFD
	go func() {
		defer wg.Done()
		availHFD = ev.scadaHFDSummary()
	}()
	// Met Tower
	go func() {
		defer wg.Done()
		availMet = ev.metTowerSummary()
	}()
	// OEM PROJECT
	go func() {
		defer wg.Done()
		availOEMProject = ev.scadaOEMSummaryProject()
	}()
	// HFD PROJECT
	go func() {
		defer wg.Done()
		availHFDProject = ev.scadaHFDSummaryProject()
	}()
	// HFD Daily
	// go func() {
	// 	availHFDDaily = ev.scadaHFDSummaryDailyTurbine()
	// 	wg.Done()
	// }()

	wg.Wait()

	ev.Ctx.DeleteMany(new(DataAvailability), nil)

	e = ev.Ctx.Insert(availOEM)
	if e != nil {
		ev.Log.AddLog(tk.Sprintf("Found : %s"+e.Error()), sError)
	}
	e = ev.Ctx.Insert(availHFD)
	if e != nil {
		ev.Log.AddLog(tk.Sprintf("Found : %s"+e.Error()), sError)
	}
	e = ev.Ctx.Insert(availMet)
	if e != nil {
		ev.Log.AddLog(tk.Sprintf("Found : %s"+e.Error()), sError)
	}
	// e = ev.Ctx.Insert(availHFDDaily)
	// if e != nil {
	// 	ev.Log.AddLog(tk.Sprintf("Found : %s"+e.Error()), sError)
	// }
	/* ============== PROJECT LEVEL ============== */
	availMet.ID = availMet.ID + "_PROJECT"
	availMet.Type = availMet.Type + "_PROJECT"
	availMet.Name = availMet.Name + " PROJECT"
	e = ev.Ctx.Insert(availMet)
	if e != nil {
		ev.Log.AddLog(tk.Sprintf("Found : %s"+e.Error()), sError)
	}
	e = ev.Ctx.Insert(availOEMProject)
	if e != nil {
		ev.Log.AddLog(tk.Sprintf("Found : %s"+e.Error()), sError)
	}
	e = ev.Ctx.Insert(availHFDProject)
	if e != nil {
		ev.Log.AddLog(tk.Sprintf("Found : %s"+e.Error()), sError)
	}

	ev.Log.AddLog("===================== End process Data Availability Summary...", sInfo)
}

func (ev *DataAvailabilitySummary) scadaOEMSummary() *DataAvailability {
	ev.Log.AddLog("===================== SCADA DATA OEM...", sInfo)
	availability := new(DataAvailability)
	availability.Name = "Scada Data OEM"
	availability.Type = "SCADA_DATA_OEM"

	var wg sync.WaitGroup

	ctx, e := PrepareConnection()
	defer ctx.Close()
	if e != nil {
		ev.Log.AddLog(tk.Sprintf("Found : %s"+e.Error()), sWarning)
		os.Exit(0)
	}

	details := []DataAvailabilityDetail{}

	now := getTimeNow()

	periodTo, _ := time.Parse("20060102_150405", now.Format("20060102_")+"000000")
	id := now.Format("20060102_150405_SCADAOEM")

	// latest 6 month
	periodFrom := GetNormalAddDateMonth(periodTo.UTC(), monthBefore)
	periodTo = now

	availability.ID = id
	availability.PeriodFrom = periodFrom
	availability.PeriodTo = periodTo

	wg.Add(len(ev.BaseController.RefTurbines))
	for turbine, turbineVal := range ev.BaseController.RefTurbines {
		value, _ := tk.ToM(turbineVal)
		match := tk.M{}
		match.Set("projectname", value.GetString("project"))
		match.Set("turbine", turbine)
		match.Set("timestamp", tk.M{"$gte": periodFrom})
		go workerTurbine(turbine, value.GetString("project"), new(ScadaDataOEM).TableName(), match, &details, ev, ctx, &wg)
	}
	wg.Wait()

	availability.Details = details

	return availability
}

func (ev *DataAvailabilitySummary) scadaHFDSummary() *DataAvailability {
	ev.Log.AddLog("===================== SCADA DATA HFD...", sInfo)
	availability := new(DataAvailability)
	availability.Name = "Scada Data HFD"
	availability.Type = "SCADA_DATA_HFD"

	var wg sync.WaitGroup

	ctx, e := PrepareConnection()
	defer ctx.Close()
	if e != nil {
		ev.Log.AddLog(tk.Sprintf("Found : %s"+e.Error()), sWarning)
		os.Exit(0)
	}

	details := []DataAvailabilityDetail{}

	now := getTimeNow()

	periodTo, _ := time.Parse("20060102_150405", now.Format("20060102_")+"000000")
	id := now.Format("20060102_150405_SCADAHFD")

	// latest 6 month
	periodFrom := GetNormalAddDateMonth(periodTo.UTC(), monthBefore)
	periodTo = now

	availability.ID = id
	availability.PeriodFrom = periodFrom
	availability.PeriodTo = periodTo
	match := tk.M{}

	wg.Add(len(ev.BaseController.RefTurbines))
	for turbine, turbineVal := range ev.BaseController.RefTurbines {
		value, _ := tk.ToM(turbineVal)

		match.Set("projectname", value.GetString("project"))
		match.Set("turbine", turbine)
		match.Set("timestamp", tk.M{"$gte": periodFrom})
		match.Set("isnull", false)
		go workerTurbine(turbine, value.GetString("project"), "Scada10MinHFD", match, &details, ev, ctx, &wg)
	}
	wg.Wait()

	availability.Details = details

	return availability
}

func (ev *DataAvailabilitySummary) metTowerSummary() *DataAvailability {
	ev.Log.AddLog("===================== MET TOWER...", sInfo)
	availability := new(DataAvailability)
	availability.Name = "Met Tower"
	availability.Type = "MET_TOWER"

	var wg sync.WaitGroup

	ctx, e := PrepareConnection()
	defer ctx.Close()
	if e != nil {
		ev.Log.AddLog(tk.Sprintf("Found : %s"+e.Error()), sWarning)
		os.Exit(0)
	}

	details := []DataAvailabilityDetail{}

	now := getTimeNow()

	periodTo, _ := time.Parse("20060102_150405", now.Format("20060102_")+"000000")
	id := now.Format("20060102_150405_METTOWER")

	// latest 6 month
	periodFrom := GetNormalAddDateMonth(periodTo.UTC(), monthBefore)
	periodTo = now

	availability.ID = id
	availability.PeriodFrom = periodFrom
	availability.PeriodTo = periodTo

	wg.Add(len(ev.BaseController.RefTurbines))
	for turbine, turbineVal := range ev.BaseController.RefTurbines {
		value, _ := tk.ToM(turbineVal)
		match := tk.M{}
		match.Set("timestamp", tk.M{"$gte": periodFrom})
		go workerTurbine(turbine, value.GetString("project"), new(MetTower).TableName(), match, &details, ev, ctx, &wg)
	}
	wg.Wait()

	availability.Details = details

	return availability
}

func (ev *DataAvailabilitySummary) scadaOEMSummaryProject() *DataAvailability {
	ev.Log.AddLog("===================== SCADA DATA OEM PROJECT LEVEL . . .", sInfo)
	availability := new(DataAvailability)
	availability.Name = "Scada Data OEM PROJECT"
	availability.Type = "SCADA_DATA_OEM_PROJECT"

	var wg sync.WaitGroup

	ctx, e := PrepareConnection()
	defer ctx.Close()
	if e != nil {
		ev.Log.AddLog(tk.Sprintf("Found : %s"+e.Error()), sWarning)
		os.Exit(0)
	}

	details := []DataAvailabilityDetail{}

	now := getTimeNow()

	periodTo, _ := time.Parse("20060102_150405", now.Format("20060102_")+"000000")
	id := now.Format("20060102_150405_SCADAOEM_PROJECT")

	// latest 6 month
	periodFrom := GetNormalAddDateMonth(periodTo.UTC(), monthBefore)
	periodTo = now

	availability.ID = id
	availability.PeriodFrom = periodFrom
	availability.PeriodTo = periodTo

	wg.Add(len(ev.BaseController.ProjectList))
	for _, projectData := range ev.BaseController.ProjectList {
		projectID := projectData.ProjectId
		if projectID != "" {
			match := tk.M{}
			match.Set("projectname", projectID)
			match.Set("timestamp", tk.M{"$gte": periodFrom})
			go workerProject(projectID, new(ScadaDataOEM).TableName(), match, &details, ev, ctx, &wg)
		}
	}
	wg.Wait()

	availability.Details = details

	return availability
}

func (ev *DataAvailabilitySummary) scadaHFDSummaryProject() *DataAvailability {
	ev.Log.AddLog("===================== SCADA DATA HFD PROJECT LEVEL . . .", sInfo)
	availability := new(DataAvailability)
	availability.Name = "Scada Data HFD PROJECT"
	availability.Type = "SCADA_DATA_HFD_PROJECT"

	var wg sync.WaitGroup

	ctx, e := PrepareConnection()
	defer ctx.Close()
	if e != nil {
		ev.Log.AddLog(tk.Sprintf("Found : %s"+e.Error()), sWarning)
		os.Exit(0)
	}

	details := []DataAvailabilityDetail{}

	now := getTimeNow()

	periodTo, _ := time.Parse("20060102_150405", now.Format("20060102_")+"000000")
	id := now.Format("20060102_150405_SCADAHFD_PROJECT")

	// latest 6 month
	periodFrom := GetNormalAddDateMonth(periodTo.UTC(), monthBefore)
	periodTo = now

	availability.ID = id
	availability.PeriodFrom = periodFrom
	availability.PeriodTo = periodTo

	wg.Add(len(ev.BaseController.ProjectList))
	for _, projectData := range ev.BaseController.ProjectList {
		projectID := projectData.ProjectId
		if projectID != "" {
			match := tk.M{}
			match.Set("projectname", projectID)
			match.Set("timestamp", tk.M{"$gte": periodFrom})
			match.Set("isnull", false)
			go workerProject(projectID, "Scada10MinHFD", match, &details, ev, ctx, &wg)
		}
	}
	wg.Wait()

	availability.Details = details

	return availability
}

func workerTurbine(t, projectName, tablename string, match tk.M, details *[]DataAvailabilityDetail,
	ev *DataAvailabilitySummary, ctx dbox.IConnection, wg *sync.WaitGroup) {
	defer wg.Done()
	now := getTimeNow() /* bulan ini */
	periodTo, _ := time.Parse("20060102_150405", now.Format("20060102_")+"000000")
	periodFrom := GetNormalAddDateMonth(periodTo.UTC(), monthBefore) /* sampai 6 bulan ke belakang */
	periodTo = now

	start := time.Now()

	pipes := []tk.M{}
	pipes = append(pipes, tk.M{"$match": match})
	pipes = append(pipes, tk.M{"$project": tk.M{"timestamp": 1}})
	pipes = append(pipes, tk.M{"$sort": tk.M{"timestamp": 1}})

	csr, e := ctx.NewQuery().From(tablename).
		Command("pipe", pipes).Cursor(nil)
	if e != nil {
		ev.Log.AddLog(tk.Sprintf("Found : %s"+e.Error()), sWarning)
	}
	defer csr.Close()

	type CustomObject struct {
		TimeStamp time.Time
	}
	list := []CustomObject{}
	e = csr.Fetch(&list, 0, false)
	if e != nil {
		ev.Log.AddLog(tk.Sprintf("Found : %s"+e.Error()), sWarning)
	}

	before := CustomObject{} /* untuk penentuan awal durasi unavailable */
	from := CustomObject{}   /* untuk penentuan awal durasi available */
	latestData := CustomObject{}
	hoursGap := 0.0
	duration := 0.0
	countID := 0 /* buat index ordering saat data udah disimpan di DB */
	/* detail adalah additional detail untuk unavailable data jika ada gap di :
	1. mulai dari awal timestamp defined sampai awal timestamp dari DB
	2. mulai dari latest timestamp dari DB sampai akhir timestamp defined
	*/
	detail := []DataAvailabilityDetail{}

	/* ============================ LOGIC ========================
		unavailable => jika oem - before > 24 jam
		available => adalah before - from
	===============================================================*/

	if len(list) > 0 {
		/* ============= START looping data dari DB ============== */
		for idx, oem := range list {
			if idx > 0 {
				before = list[idx-1]
				hoursGap = oem.TimeStamp.UTC().Sub(before.TimeStamp.UTC()).Hours()

				if hoursGap > 24 { /* ketika ketemu yang unavailable lagi */
					countID++
					// set duration for available datas
					duration = tk.ToFloat64(before.TimeStamp.UTC().Sub(from.TimeStamp.UTC()).Hours()/24, 6, tk.RoundingAuto)
					/* durasi available mulai dari data from yang terakhir tersimpan sampai data index-1 */
					mtx.Lock()
					*details = append(*details, setDataAvailDetail(from.TimeStamp, before.TimeStamp, projectName, t, duration, true, countID))
					mtx.Unlock()
					// set duration for unavailable datas
					countID++
					duration = tk.ToFloat64(hoursGap/24, 6, tk.RoundingAuto)
					/* durasi unavailable mulai dari data index-1 sampai data index saat ini */
					mtx.Lock()
					*details = append(*details, setDataAvailDetail(before.TimeStamp, oem.TimeStamp, projectName, t, duration, false, countID))
					mtx.Unlock()
					from = oem
				}
			} else {
				from = oem

				// set gap from stardate defined by us until actual first data in db
				hoursGap = from.TimeStamp.UTC().Sub(periodFrom.UTC()).Hours()
				if hoursGap > 24 {
					countID++
					duration = tk.ToFloat64(hoursGap/24, 6, tk.RoundingAuto) /* dibuat per hari */
					/* dianggap false (not avail) karena gap startdate defined dengan startdate actual DB > 24 jam */
					detail = append(detail, setDataAvailDetail(periodFrom, from.TimeStamp, projectName, t, duration, false, countID))
				}
			}

			latestData = oem
		}
		/* =============  END OF looping data dari DB ============== */

		hoursGap = latestData.TimeStamp.UTC().Sub(from.TimeStamp.UTC()).Hours()
		/* jika data terakhir dikurangi data from yang terakhir tersimpan > 24 jam
		maka dianggap AVAILABLE karena data ini bisa jadi belum ter plot
		karena jika oem-before < 24 jam akan dilewati pada logic di dalam looping di atas */
		if hoursGap > 24 {
			countID++
			duration = tk.ToFloat64(hoursGap/24, 6, tk.RoundingAuto)
			mtx.Lock()
			*details = append(*details, setDataAvailDetail(from.TimeStamp, latestData.TimeStamp, projectName, t, duration, true, countID))
			mtx.Unlock()
		}

		// set gap from last data until periodTo
		hoursGap = periodTo.UTC().Sub(latestData.TimeStamp.UTC()).Hours()
		/* jika gap dari time.Now - latestData > 24 jam maka set sebagai UNAVAILABLE
		karena bisa jadi latestData yang tersimpan di DB tidak sampai time.Now */
		if hoursGap > 24 {
			countID++
			duration = tk.ToFloat64(hoursGap/24, 6, tk.RoundingAuto)
			detail = append(detail, setDataAvailDetail(latestData.TimeStamp, periodTo, projectName, t, duration, false, countID))
		}
		/* jika tidak ada additional detail sama sekali maka data dianggap FULL AVAILABLE selam 6 bulan */
		if len(detail) == 0 {
			countID++
			duration = tk.ToFloat64(periodTo.Sub(periodFrom).Hours()/24, 6, tk.RoundingAuto)
			detail = append(detail, setDataAvailDetail(periodFrom, periodTo, projectName, t, duration, true, countID))
		}
	} else {
		/* jika list data di DB tidak ada maka di set FULL UNAVAILABLE selama 6 bulan*/
		countID++
		duration = tk.ToFloat64(periodTo.Sub(periodFrom).Hours()/24, 6, tk.RoundingAuto)
		detail = append(detail, setDataAvailDetail(periodFrom, periodTo, projectName, t, duration, false, countID))
	}
	mtx.Lock()
	*details = append(*details, detail...)
	ev.Log.AddLog(tk.Sprintf(">> DONE: %v | %v | %v secs \n", t, len(list), time.Now().Sub(start).Seconds()), sInfo)
	mtx.Unlock()
}

func workerProject(projectName, tablename string, match tk.M, details *[]DataAvailabilityDetail,
	ev *DataAvailabilitySummary, ctx dbox.IConnection, wg *sync.WaitGroup) {
	defer wg.Done()
	now := getTimeNow()
	periodTo, _ := time.Parse("20060102_150405", now.Format("20060102_")+"000000")
	periodFrom := GetNormalAddDateMonth(periodTo.UTC(), monthBefore)
	periodTo = now

	detail := []DataAvailabilityDetail{}
	start := time.Now()

	pipes := []tk.M{}
	pipes = append(pipes, tk.M{"$match": match})
	pipes = append(pipes, tk.M{"$project": tk.M{"projectname": 1, "timestamp": 1}})
	pipes = append(pipes, tk.M{"$sort": tk.M{"timestamp": 1}})

	csr, e := ctx.NewQuery().From(tablename).
		Command("pipe", pipes).Cursor(nil)
	if e != nil {
		ev.Log.AddLog(tk.Sprintf("Found : %s"+e.Error()), sWarning)
	}

	defer csr.Close()

	type Scada10MinHFDCustom struct {
		TimeStamp time.Time
	}
	list := []Scada10MinHFDCustom{}
	e = csr.Fetch(&list, 0, false)
	if e != nil {
		ev.Log.AddLog(tk.Sprintf("Found : %s"+e.Error()), sWarning)
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
}

func (ev *DataAvailabilitySummary) scadaHFDSummaryDailyProject() *DataAvailability {
	ev.Log.AddLog("===================== SCADA DATA HFD DAILY PROJECT LEVEL...", sInfo)
	var wg sync.WaitGroup

	ctx, e := PrepareConnection()
	defer ctx.Close()
	if e != nil {
		ev.Log.AddLog(tk.Sprintf("Found : %s"+e.Error()), sWarning)
		os.Exit(0)
	}
	now := getTimeNow()

	periodTo, _ := time.Parse("20060102_150405", now.Format("20060102_")+"000000")
	id := now.Format("20060102_150405_SCADAHFD_DAILY_PROJECT")

	// latest 6 month
	periodFrom := GetNormalAddDateMonth(periodTo.UTC(), monthBefore)
	periodTo = now

	availabilityDaily := new(DataAvailability)
	availabilityDaily.Name = "Scada Data HFD DAILY PROJECT"
	availabilityDaily.Type = "SCADA_DATA_HFD_DAILY_PROJECT"
	availabilityDaily.ID = id
	availabilityDaily.PeriodFrom = periodFrom
	availabilityDaily.PeriodTo = periodTo

	matches := []tk.M{
		tk.M{"dateinfo.dateid": tk.M{"$gte": periodFrom}},
		tk.M{"isnull": false},
	}
	groups := tk.M{
		"_id": tk.M{
			"projectname": "$projectname",
			"tanggal":     "$dateinfo.dateid",
		},
		"totaldata": tk.M{"$sum": 1},
	}
	projection := tk.M{
		"projectname": "$_id.projectname",
		"tanggal":     "$_id.tanggal",
		"totaldata":   1,
	}

	pipes := []tk.M{}
	pipes = append(pipes, tk.M{"$match": tk.M{"$and": matches}})
	pipes = append(pipes, tk.M{"$group": groups})
	pipes = append(pipes, tk.M{"$project": projection})
	pipes = append(pipes, tk.M{"$sort": tk.M{"_id.tanggal": 1}})

	csr, e := ctx.NewQuery().From("Scada10MinHFD").
		Command("pipe", pipes).Cursor(nil)
	if e != nil {
		ev.Log.AddLog(tk.Sprintf("Found : %s"+e.Error()), sWarning)
	}

	_data := tk.M{}
	dataPerProject := map[string][]tk.M{} /* for appending data per project */
	currProject := ""
	for {
		_data = tk.M{}
		e = csr.Fetch(&_data, 1, false)
		if e != nil {
			break
		}
		currProject = _data.GetString("projectname")
		dataPerProject[currProject] = append(dataPerProject[currProject], _data)
	}
	csr.Close()

	detailsDaily := []DataAvailabilityDetail{}

	wg.Add(len(ev.BaseController.ProjectList))
	// for turbine, turbineVal := range ev.BaseController.ProjectList {
	// 	value, _ := tk.ToM(turbineVal)
	// 	currProject = value.GetString("project")
	// 	go workerDaily(dataPerProject[currProject], 1.0, &detailsDaily, &wg)
	// }
	wg.Wait()

	availabilityDaily.Details = detailsDaily
	ev.Log.AddLog(tk.Sprintf(">> DONE SCADA HFD DAILY PROJECT LEVEL"), sInfo)

	return availabilityDaily
}

func (ev *DataAvailabilitySummary) scadaHFDSummaryDailyTurbine() *DataAvailability {
	ev.Log.AddLog("===================== SCADA DATA HFD DAILY TURBINE LEVEL...", sInfo)
	var wg sync.WaitGroup

	now := getTimeNow()

	periodTo, _ := time.Parse("20060102_150405", now.Format("20060102_")+"000000")
	id := now.Format("20060102_150405_SCADAHFD_DAILY")

	// latest 6 month
	periodFrom := GetNormalAddDateMonth(periodTo.UTC(), monthBefore)
	periodTo = now

	availabilityDaily := new(DataAvailability)
	availabilityDaily.Name = "Scada Data HFD DAILY"
	availabilityDaily.Type = "SCADA_DATA_HFD_DAILY"
	availabilityDaily.ID = id
	availabilityDaily.PeriodFrom = periodFrom
	availabilityDaily.PeriodTo = periodTo

	countx := 0
	maxPar := 5

	detailsDaily := []DataAvailabilityDetail{}

	for currTurbine, turbineVal := range ev.BaseController.RefTurbines {
		wg.Add(1)
		value, _ := tk.ToM(turbineVal)
		currProject := value.GetString("project")
		go func(project, turbine string) {
			defer wg.Done()
			ev.workerDailyTurbine(project, turbine, &detailsDaily)
		}(currProject, currTurbine)

		countx++

		if countx%maxPar == 0 || (len(ev.BaseController.RefTurbines) == countx) {
			// wg.Wait()
		}
	}
	wg.Wait()

	availabilityDaily.Details = detailsDaily
	ev.Log.AddLog(tk.Sprintf(">> DONE SCADA HFD DAILY"), sInfo)

	return availabilityDaily
}

func (ev *DataAvailabilitySummary) workerRangeTurbine(project, turbine string, details *[]DataAvailabilityDetail) {
	detail := []DataAvailabilityDetail{}
	start := time.Now()

	totalTurbine := 1.0
	maxDataPerDay := 6.0 * 24.0 * totalTurbine /* dalam 1 jam ada 6 data karena per 10 menit dikalikan 24 karena 1 hari ada 24 jam*/
	duration := 0.0
	countID := 1

	now := getTimeNow() /* time.Now India */
	periodTo, _ := time.Parse("20060102_150405", now.Format("20060102_")+"000000")
	periodFrom := GetNormalAddDateMonth(periodTo.UTC(), monthBefore) // latest 6 month
	periodTo = now

	ctx, e := PrepareConnection()
	defer ctx.Close()
	if e != nil {
		ev.Log.AddLog(tk.Sprintf("Found : %s"+e.Error()), sWarning)
		os.Exit(0)
	}

	matches := []tk.M{
		tk.M{"projectname": project},
		tk.M{"turbine": turbine},
		tk.M{"dateinfo.dateid": tk.M{"$gte": periodFrom}},
		tk.M{"isnull": false},
	}
	groups := tk.M{
		"_id": tk.M{
			"tanggal": "$dateinfo.dateid",
		},
		"totaldata": tk.M{"$sum": 1},
	}
	projection := tk.M{
		"tanggal":   "$_id.tanggal",
		"totaldata": 1,
	}

	pipes := []tk.M{}
	pipes = append(pipes, tk.M{"$match": tk.M{"$and": matches}})
	pipes = append(pipes, tk.M{"$group": groups})
	pipes = append(pipes, tk.M{"$project": projection})
	pipes = append(pipes, tk.M{"$sort": tk.M{"_id.tanggal": 1}})

	csr, e := ctx.NewQuery().From("Scada10MinHFD").
		Command("pipe", pipes).Cursor(nil)
	if e != nil {
		ev.Log.AddLog(tk.Sprintf("Found : %s"+e.Error()), sWarning)
	}

	type ObjectFetch struct {
		Tanggal   time.Time
		TotalData float64
	}
	data := []ObjectFetch{}
	_data := ObjectFetch{}
	for {
		_data = ObjectFetch{}
		e = csr.Fetch(&_data, 1, false)
		if e != nil {
			break
		}
		data = append(data, _data)
	}
	csr.Close()

	before := time.Time{}
	latestTimeStamp := time.Time{}
	currTimeStamp := time.Time{}
	detail = []DataAvailabilityDetail{}
	hoursGap := 0.0
	var mux sync.Mutex

	if len(data) > 0 {
		/* ============= START looping data dari DB ============== */
		for idx, hfd := range data {
			currTimeStamp = hfd.Tanggal.UTC()
			if idx > 0 {
				before = data[idx-1].Tanggal.UTC()
				hoursGap = currTimeStamp.UTC().Sub(before.UTC()).Hours()

				if hoursGap > 24 { /* jika timestamp saat ini dibanding timestamp sebelumnya lebih dari 1 hari maka otomatis unavailable */
					// set duration for unavailable datas
					countID++
					duration = tk.ToFloat64(hoursGap, 2, tk.RoundingAuto)
					/* durasi unavailable mulai dari data index-1 sampai data index saat ini */
					mtx.Lock()
					*details = append(*details, setDataAvailDetail(before, currTimeStamp, project, turbine, duration, false, countID))
					mtx.Unlock()
				} else {
					duration = tk.Div(hfd.TotalData, maxDataPerDay)
					if duration >= 0.5 { /* dianggap available karena durasi >= setengah hari */
						countID++
						mux.Lock()
						*details = append(*details, setDataAvailDetail(currTimeStamp, currTimeStamp, project,
							turbine, duration, true, countID))
						mux.Unlock()
					} else { /* dianggap not available karena durasi < setengah hari */
						countID++
						mux.Lock()
						*details = append(*details, setDataAvailDetail(currTimeStamp, currTimeStamp, project,
							turbine, duration, false, countID))
						mux.Unlock()
					}
				}
			} else {
				// set gap from stardate defined by us until actual first data in db
				hoursGap = currTimeStamp.UTC().Sub(periodFrom.UTC()).Hours()

				if hoursGap > 24 {
					countID++
					duration = tk.ToFloat64(hoursGap, 2, tk.RoundingAuto)
					detail = append(detail, setDataAvailDetail(periodFrom, currTimeStamp, project, turbine, duration, false, countID))
				} else {
					duration = tk.Div(hfd.TotalData, maxDataPerDay)
					if duration >= 0.5 { /* dianggap available karena durasi >= setengah hari */
						countID++
						detail = append(detail, setDataAvailDetail(currTimeStamp, currTimeStamp, project,
							turbine, duration, true, countID))
					} else { /* dianggap not available karena durasi < setengah hari */
						countID++
						detail = append(detail, setDataAvailDetail(currTimeStamp, currTimeStamp, project,
							turbine, duration, false, countID))
					}
				}
			}

			latestTimeStamp = currTimeStamp
		}
		/* =============  END OF looping data dari DB ============== */

		// set gap from last data until periodTo
		hoursGap = periodTo.Sub(latestTimeStamp).Hours()
		/* jika gap dari time.Now - latestTimeStamp > 24 jam maka set sebagai UNAVAILABLE
		karena bisa jadi latestTimeStamp yang tersimpan di DB tidak sampai time.Now */
		if hoursGap > 24 {
			countID++
			duration = 24 // duration langsung diset 24 jam
			mulai := latestTimeStamp
		perDayLoopLast:
			for {
				if mulai.After(periodTo) {
					break perDayLoopLast
				}
				detail = append(detail, setDataAvailDetail(mulai, mulai, project,
					turbine, duration, false, countID))
				mulai = mulai.AddDate(0, 0, 1)

			}
		}
		/* jika tidak ada additional detail sama sekali maka data dianggap FULL AVAILABLE selam 6 bulan */
		if len(detail) == 0 {
			countID++
			duration = 24 // duration langsung diset 24 jam
			mulai := periodFrom
		perDayLoopFullAvail:
			for {
				if mulai.After(periodTo) {
					break perDayLoopFullAvail
				}
				detail = append(detail, setDataAvailDetail(mulai, mulai, project,
					turbine, duration, true, countID))
				mulai = mulai.AddDate(0, 0, 1)

			}
		}
	} else {
		/* jika list data di DB tidak ada maka di set FULL UNAVAILABLE selama 6 bulan*/
		countID++
		duration = 24 // duration langsung diset 24 jam
		mulai := periodFrom
	perDayLoopFullNotAvail:
		for {
			if mulai.After(periodTo) {
				break perDayLoopFullNotAvail
			}
			detail = append(detail, setDataAvailDetail(mulai, mulai, project,
				turbine, duration, false, countID))
			mulai = mulai.AddDate(0, 0, 1)

		}
	}
	mux.Lock()
	*details = append(*details, detail...)
	ev.Log.AddLog(tk.Sprintf(">> DONE: %v | %v | %v secs \n", turbine, len(data), time.Now().Sub(start).Seconds()), sInfo)
	mux.Unlock()
}

func (ev *DataAvailabilitySummary) workerDailyTurbine(project, turbine string, details *[]DataAvailabilityDetail) {
	detail := []DataAvailabilityDetail{}
	start := time.Now()

	totalTurbine := 1.0
	maxDataPerDay := 6.0 * 24.0 * totalTurbine /* dalam 1 jam ada 6 data karena per 10 menit dikalikan 24 karena 1 hari ada 24 jam*/
	duration := 0.0
	countID := 1

	now := getTimeNow() /* time.Now India */
	periodTo, _ := time.Parse("20060102_150405", now.Format("20060102_")+"000000")
	periodFrom := GetNormalAddDateMonth(periodTo.UTC(), monthBefore) // latest 6 month
	periodTo = now

	ctx, e := PrepareConnection()
	defer ctx.Close()
	if e != nil {
		ev.Log.AddLog(tk.Sprintf("Found : %s"+e.Error()), sWarning)
		os.Exit(0)
	}

	matches := []tk.M{
		tk.M{"projectname": project},
		tk.M{"turbine": turbine},
		tk.M{"dateinfo.dateid": tk.M{"$gte": periodFrom}},
		tk.M{"isnull": false},
	}
	groups := tk.M{
		"_id": tk.M{
			"tanggal": "$dateinfo.dateid",
		},
		"totaldata": tk.M{"$sum": 1},
	}
	projection := tk.M{
		"tanggal":   "$_id.tanggal",
		"totaldata": 1,
	}

	pipes := []tk.M{}
	pipes = append(pipes, tk.M{"$match": tk.M{"$and": matches}})
	pipes = append(pipes, tk.M{"$group": groups})
	pipes = append(pipes, tk.M{"$project": projection})
	pipes = append(pipes, tk.M{"$sort": tk.M{"_id.tanggal": 1}})

	csr, e := ctx.NewQuery().From("Scada10MinHFD").
		Command("pipe", pipes).Cursor(nil)
	if e != nil {
		ev.Log.AddLog(tk.Sprintf("Found : %s"+e.Error()), sWarning)
	}

	type ObjectFetch struct {
		Tanggal   time.Time
		TotalData float64
	}
	data := []ObjectFetch{}
	_data := ObjectFetch{}
	for {
		_data = ObjectFetch{}
		e = csr.Fetch(&_data, 1, false)
		if e != nil {
			break
		}
		data = append(data, _data)
	}
	csr.Close()

	before := time.Time{}
	latestTimeStamp := time.Time{}
	currTimeStamp := time.Time{}
	detail = []DataAvailabilityDetail{}
	hoursGap := 0.0
	var mux sync.Mutex

	if len(data) > 0 {
		/* ============= START looping data dari DB ============== */
		for idx, hfd := range data {
			currTimeStamp = hfd.Tanggal.UTC()
			if idx > 0 {
				before = data[idx-1].Tanggal.UTC()
				hoursGap = currTimeStamp.UTC().Sub(before.UTC()).Hours()

				if hoursGap > 24 { /* jika timestamp saat ini dibanding timestamp sebelumnya lebih dari 1 hari maka otomatis unavailable */
					// set duration for unavailable datas
					countID++
					duration = 24 // duration langsung diset 24 jam
					mulai := before
				perDayLoop:
					for {
						/* dianggap false (not avail) karena gap startdate defined dengan startdate actual DB > 24 jam */
						if mulai.After(currTimeStamp.AddDate(0, 0, -1)) {
							break perDayLoop
						}
						mux.Lock()
						*details = append(*details, setDataAvailDetail(mulai, mulai, project,
							turbine, duration, false, countID))
						mux.Unlock()
						mulai = mulai.AddDate(0, 0, 1)

					}
				} else {
					duration = tk.Div(hfd.TotalData, maxDataPerDay)
					if duration >= 0.5 { /* dianggap available karena durasi >= setengah hari */
						countID++
						mux.Lock()
						*details = append(*details, setDataAvailDetail(currTimeStamp, currTimeStamp, project,
							turbine, duration, true, countID))
						mux.Unlock()
					} else { /* dianggap not available karena durasi < setengah hari */
						countID++
						mux.Lock()
						*details = append(*details, setDataAvailDetail(currTimeStamp, currTimeStamp, project,
							turbine, duration, false, countID))
						mux.Unlock()
					}
				}
			} else {
				// set gap from stardate defined by us until actual first data in db
				hoursGap = currTimeStamp.UTC().Sub(periodFrom.UTC()).Hours()
				if hoursGap > 24 {
					countID++
					duration = 24 // duration langsung diset 24 jam
					mulai := periodFrom
				perDayLoopFirst:
					for {
						/* dianggap false (not avail) karena gap startdate defined dengan startdate actual DB > 24 jam */
						if mulai.After(currTimeStamp.AddDate(0, 0, -1)) {
							break perDayLoopFirst
						}
						detail = append(detail, setDataAvailDetail(mulai, mulai, project,
							turbine, duration, false, countID))
						mulai = mulai.AddDate(0, 0, 1)

					}
				} else {
					duration = tk.Div(hfd.TotalData, maxDataPerDay)
					if duration >= 0.5 { /* dianggap available karena durasi >= setengah hari */
						countID++
						detail = append(detail, setDataAvailDetail(currTimeStamp, currTimeStamp, project,
							turbine, duration, true, countID))
					} else { /* dianggap not available karena durasi < setengah hari */
						countID++
						detail = append(detail, setDataAvailDetail(currTimeStamp, currTimeStamp, project,
							turbine, duration, false, countID))
					}
				}
			}

			latestTimeStamp = currTimeStamp
		}
		/* =============  END OF looping data dari DB ============== */

		// set gap from last data until periodTo
		hoursGap = periodTo.Sub(latestTimeStamp).Hours()
		/* jika gap dari time.Now - latestTimeStamp > 24 jam maka set sebagai UNAVAILABLE
		karena bisa jadi latestTimeStamp yang tersimpan di DB tidak sampai time.Now */
		if hoursGap > 24 {
			countID++
			duration = 24 // duration langsung diset 24 jam
			mulai := latestTimeStamp
		perDayLoopLast:
			for {
				if mulai.After(periodTo) {
					break perDayLoopLast
				}
				detail = append(detail, setDataAvailDetail(mulai, mulai, project,
					turbine, duration, false, countID))
				mulai = mulai.AddDate(0, 0, 1)

			}
		}
		/* jika tidak ada additional detail sama sekali maka data dianggap FULL AVAILABLE selam 6 bulan */
		if len(detail) == 0 {
			countID++
			duration = 24 // duration langsung diset 24 jam
			mulai := periodFrom
		perDayLoopFullAvail:
			for {
				if mulai.After(periodTo) {
					break perDayLoopFullAvail
				}
				detail = append(detail, setDataAvailDetail(mulai, mulai, project,
					turbine, duration, true, countID))
				mulai = mulai.AddDate(0, 0, 1)

			}
		}
	} else {
		/* jika list data di DB tidak ada maka di set FULL UNAVAILABLE selama 6 bulan*/
		countID++
		duration = 24 // duration langsung diset 24 jam
		mulai := periodFrom
	perDayLoopFullNotAvail:
		for {
			if mulai.After(periodTo) {
				break perDayLoopFullNotAvail
			}
			detail = append(detail, setDataAvailDetail(mulai, mulai, project,
				turbine, duration, false, countID))
			mulai = mulai.AddDate(0, 0, 1)

		}
	}
	mux.Lock()
	*details = append(*details, detail...)
	ev.Log.AddLog(tk.Sprintf(">> DONE: %v | %v | %v secs \n", turbine, len(data), time.Now().Sub(start).Seconds()), sInfo)
	mux.Unlock()
}

func setDataAvailDetail(from time.Time, to time.Time, project string, turbine string, duration float64, isAvail bool, id int) DataAvailabilityDetail {

	res := DataAvailabilityDetail{
		ID:          id,
		ProjectName: project,
		Turbine:     turbine,
		TurbineName: turbineName[tk.Sprintf("%s_%s", project, turbine)],
		Start:       from.UTC(),
		StartInfo:   GetDateInfo(from.UTC()),
		End:         to.UTC(),
		EndInfo:     GetDateInfo(to.UTC()),
		Duration:    duration,
		IsAvail:     isAvail,
	}

	return res
}

func GetTurbineNameListAll(project string) (turbineNameData map[string]string, err error) {
	ctx, err := PrepareConnection()
	defer ctx.Close()
	if err != nil {
		return
	}
	query := ctx.NewQuery().From("ref_turbine")
	if project != "" && project != "Fleet" {
		pipes := []tk.M{
			tk.M{"$match": tk.M{"project": project}},
		}
		query = query.Command("pipe", pipes)
	}
	csrTurbine, err := query.Cursor(nil)
	if err != nil {
		return
	}
	defer csrTurbine.Close()
	turbineList := []tk.M{}
	err = csrTurbine.Fetch(&turbineList, 0, false)
	if err != nil {
		return
	}
	turbineNameData = map[string]string{}
	for _, val := range turbineList {
		if project != "" {
			turbineNameData[val.GetString("turbineid")] = val.GetString("turbinename")
		} else {
			turbineNameData[tk.Sprintf("%s_%s", val.GetString("project"), val.GetString("turbineid"))] = val.GetString("turbinename")
		}
	}
	return
}

func getTimeNow() (tNow time.Time) {
	config := ReadConfig()
	loc, err := time.LoadLocation(config["ReadTimeLoc"])
	_Now := time.Now().UTC().Add(-time.Minute * 330)
	if err != nil {
		tk.Printfn("Get time in %s found %s", config["ReadTimeLoc"], err.Error())
	} else {
		_Now = time.Now().In(loc)
	}

	tNow = time.Date(_Now.Year(), _Now.Month(), _Now.Day(), _Now.Hour(), _Now.Minute(), _Now.Second(), _Now.Nanosecond(), time.UTC)
	return
}
