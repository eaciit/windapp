package generatorControllers

import (
	. "eaciit/wfdemo-git/library/helper"
	. "eaciit/wfdemo-git/library/models"
	. "eaciit/wfdemo-git/processapp/summaryGenerator/controllers"
	"github.com/eaciit/dbox"
	_ "github.com/eaciit/dbox/dbc/mongo"
	tk "github.com/eaciit/toolkit"
	"sync"
	"time"
)

type TurbulenceIntensityGenerator struct {
	*BaseController
}

type TurbulenceIntensityRaw struct {
	ID              string ` bson:"_id" , json:"_id" `
	Projectname     string
	Turbine         string
	Timestamp       time.Time
	DateInfo        DateInfo
	WindspeedBin    float64
	WindSpeed       float64
	WindSpeedStdDev float64
	Type            string
}

type FetchScada struct {
	ID                  string ` bson:"_id" , json:"_id" `
	Projectname         string
	Turbine             string
	Timestamp           time.Time
	DateInfo            DateInfo
	Windspeed_ms_bin    float64
	Windspeed_ms        float64
	Windspeed_ms_stddev float64
}

type FetchMet struct {
	ID              string ` bson:"_id" , json:"_id" `
	Projectname     string
	Timestamp       time.Time
	DateInfo        DateInfo
	Windspeedbin    float64
	Vhubws90mavg    float64
	Vhubws90mstddev float64
}

func (m *TurbulenceIntensityRaw) TableName() string {
	return "TurbulenceIntensity10Min"
}

type LatestTurbulenceRaw struct {
	ID          string ` bson:"_id" , json:"_id" `
	Projectname string
	LastUpdate  time.Time
	Type        string
}

func (m *LatestTurbulenceRaw) TableName() string {
	return "log_latestturbulence10min"
}

func (ev *TurbulenceIntensityGenerator) CreateTurbulenceIntensity10Min(base *BaseController) {
	ev.BaseController = base

	ev.Log.AddLog("===================== Start processing Turbulence 10 Min...", sInfo)

	var wg sync.WaitGroup
	wg.Add(2)

	go ev.processDataScada(&wg)
	// go ev.processInitialDataScada(&wg)
	go ev.processDataMet(&wg)

	wg.Wait()

	ev.Log.AddLog("===================== End processing Turbulence 10 Min...", sInfo)
}

func getstep(count int) int {
	v := count / 5
	if v == 0 {
		return 1
	}
	return v
}

func (ev *TurbulenceIntensityGenerator) processInitialDataScada(wgScada *sync.WaitGroup) {
	defer wgScada.Done()

	t0 := time.Now()
	projectList := ev.getProjectList()
	lastUpdatePerProject := ev.getLatestData("SCADA")

	tStart := time.Date(2015, 1, 1, 0, 0, 0, 0, time.UTC)
	tNow := time.Now()
	for {
		if tStart.After(tNow) {
			break
		}
		for _project := range lastUpdatePerProject {
			lastUpdatePerProject[_project] = tStart.UTC()
		}
		var wg sync.WaitGroup
		wg.Add(len(projectList))
		ev.Log.AddLog(tk.Sprintf("Updating data for %s", tStart.UTC().String()), sInfo)
		for _, _project := range projectList {
			go ev.projectInitialWorker(_project, lastUpdatePerProject[_project], &wg)
		}
		wg.Wait()

		tStart = tStart.AddDate(0, 0, 1)
	}

	ev.Log.AddLog(tk.Sprintf("Duration processing scada data %f minutes", time.Since(t0).Minutes()), sInfo)
}

func (ev *TurbulenceIntensityGenerator) projectInitialWorker(projectname string, lastUpdate time.Time, wgProject *sync.WaitGroup) {
	defer wgProject.Done()

	csr, e := ev.Ctx.Connection.NewQuery().
		From("Scada10MinHFD").
		Select("projectname", "turbine", "timestamp", "dateinfo", "windspeed_ms", "windspeed_ms_bin", "windspeed_ms_stddev").
		Where(dbox.And(dbox.Eq("dateinfo.dateid", lastUpdate),
			dbox.Eq("projectname", projectname),
			dbox.Eq("isnull", false),
			dbox.Gte("windspeed_ms_bin", 0),
			dbox.Lte("windspeed_ms_bin", 25))).
		Cursor(nil)
	if e != nil {
		ev.Log.AddLog(tk.Sprintf("Error on cursor : %s", e.Error()), sError)
	}
	defer csr.Close()

	var wg sync.WaitGroup
	totalData := csr.Count()
	totalWorker := 4
	dataChan := make(chan TurbulenceIntensityRaw, totalData)

	wg.Add(totalWorker)
	for i := 0; i < totalWorker; i++ {
		go func() {
			defer wg.Done()
			ctxWorker, e := PrepareConnection()
			if e != nil {
				ev.Log.AddLog(e.Error(), sError)
			}
			defer ctxWorker.Close()
			csrSave := ctxWorker.NewQuery().SetConfig("multiexec", true).
				From(new(TurbulenceIntensityRaw).TableName()).Save()
			defer csrSave.Close()
			for data := range dataChan {
				e = csrSave.Exec(tk.M{"data": data})
				if e != nil {
					ev.Log.AddLog(tk.Sprintf("Error on Save : %s", e.Error()), sError)
				}
			}
		}()
	}

	data := TurbulenceIntensityRaw{}
	_data := FetchScada{}
	maxTimestamp := time.Time{}

loopFetchScada:
	for {
		_data = FetchScada{}
		e = csr.Fetch(&_data, 1, false)
		if e != nil {
			break loopFetchScada
		}
		data = TurbulenceIntensityRaw{}
		data.Projectname = _data.Projectname
		data.Turbine = _data.Turbine
		data.Timestamp = _data.Timestamp.UTC()
		data.DateInfo = _data.DateInfo
		data.WindspeedBin = _data.Windspeed_ms_bin
		data.ID = tk.Sprintf("%s_%s_%s", data.Projectname, data.Turbine, data.Timestamp.Format("20060102150405"))

		if data.Timestamp.After(maxTimestamp) { /* get max timestamp for each project */
			maxTimestamp = data.Timestamp
		}

		data.WindSpeed = _data.Windspeed_ms
		data.WindSpeedStdDev = _data.Windspeed_ms_stddev
		data.Type = "SCADA"

		dataChan <- data
	}

	close(dataChan)
	wg.Wait()

	ev.updateLastData(projectname, "SCADA", maxTimestamp)
}

func (ev *TurbulenceIntensityGenerator) getProjectList() (result []string) {
	ev.Log.AddLog("Get Turbine Per Project", sInfo)

	projectData := []tk.M{}
	csrt, e := ev.Ctx.Connection.NewQuery().
		From("ref_project").Where(dbox.Eq("active", true)).Cursor(nil)
	if e != nil {
		ev.Log.AddLog(tk.Sprintf("Error on cursor at getProjectList due to : %s", e.Error()), sError)
		return
	}
	defer csrt.Close()
	e = csrt.Fetch(&projectData, 0, false)
	if e != nil {
		ev.Log.AddLog(tk.Sprintf("Error on fetch at getProjectList due to : %s", e.Error()), sError)
		return
	}
	result = []string{}
	for _, val := range projectData {
		result = append(result, val.GetString("projectid"))
	}
	ev.Log.AddLog("Finish getting Turbine Per Project", sInfo)

	return
}

func (ev *TurbulenceIntensityGenerator) getLatestData(tipe string) (result map[string]time.Time) {
	ev.Log.AddLog("Get latest data for each turbine", sInfo)

	latestData := []LatestTurbulenceRaw{}
	csrt, e := ev.Ctx.Connection.NewQuery().
		From(new(LatestTurbulenceRaw).TableName()).
		Where(dbox.Eq("type", tipe)).Cursor(nil)
	if e != nil {
		ev.Log.AddLog(tk.Sprintf("Error on cursor at getLatestData due to : %s", e.Error()), sError)
		return
	}
	defer csrt.Close()
	e = csrt.Fetch(&latestData, 0, false)
	if e != nil {
		ev.Log.AddLog(tk.Sprintf("Error on fetch at getLatestData due to : %s", e.Error()), sError)
		return
	}
	result = map[string]time.Time{}
	for _, val := range latestData {
		result[val.Projectname] = val.LastUpdate
	}
	ev.Log.AddLog("Finish getting latest data for each turbine", sInfo)

	return
}

func (ev *TurbulenceIntensityGenerator) updateLastData(projectname, tipe string, maxTimeStamp time.Time) {
	data := LatestTurbulenceRaw{}
	data.Projectname = projectname
	data.ID = tk.Sprintf("%s_%s", data.Projectname, tipe)
	data.LastUpdate = maxTimeStamp
	data.Type = tipe

	e := ev.Ctx.Connection.NewQuery().SetConfig("multiexec", true).
		From(new(LatestTurbulenceRaw).TableName()).Save().Exec(tk.M{"data": data})

	if e != nil {
		ev.Log.AddLog(tk.Sprintf("Error on Save at updateLastData due to : %s", e.Error()), sError)
	}

	ev.Log.AddLog(tk.Sprintf("Finish updating last data for %s on %s", projectname, tipe), sInfo)
}

func (ev *TurbulenceIntensityGenerator) processDataScada(wgScada *sync.WaitGroup) {
	defer wgScada.Done()

	t0 := time.Now()
	projectList := ev.getProjectList()
	lastUpdatePerProject := ev.getLatestData("SCADA")

	var wg sync.WaitGroup
	wg.Add(len(projectList))
	for _, _project := range projectList {
		go ev.projectWorker(_project, lastUpdatePerProject[_project], &wg)
	}
	wg.Wait()

	ev.Log.AddLog(tk.Sprintf("Duration processing scada data %f minutes", time.Since(t0).Minutes()), sInfo)
}

func (ev *TurbulenceIntensityGenerator) projectWorker(projectname string, lastUpdate time.Time, wgProject *sync.WaitGroup) {
	defer wgProject.Done()

	csr, e := ev.Ctx.Connection.NewQuery().
		From("Scada10MinHFD").
		Select("projectname", "turbine", "timestamp", "dateinfo", "windspeed_ms", "windspeed_ms_bin", "windspeed_ms_stddev").
		Where(dbox.And(dbox.Gte("timestamp", lastUpdate),
			dbox.Eq("projectname", projectname),
			dbox.Eq("isnull", false),
			dbox.Gte("windspeed_ms_bin", 0),
			dbox.Lte("windspeed_ms_bin", 25))).
		Cursor(nil)
	if e != nil {
		ev.Log.AddLog(tk.Sprintf("Error on cursor : %s", e.Error()), sError)
	}
	defer csr.Close()

	var wg sync.WaitGroup
	totalData := csr.Count()
	totalWorker := 4
	dataChan := make(chan TurbulenceIntensityRaw, totalData)

	wg.Add(totalWorker)
	for i := 0; i < totalWorker; i++ {
		go func() {
			defer wg.Done()
			ctxWorker, e := PrepareConnection()
			if e != nil {
				ev.Log.AddLog(e.Error(), sError)
			}
			defer ctxWorker.Close()
			csrSave := ctxWorker.NewQuery().SetConfig("multiexec", true).
				From(new(TurbulenceIntensityRaw).TableName()).Save()
			defer csrSave.Close()
			for data := range dataChan {
				e = csrSave.Exec(tk.M{"data": data})
				if e != nil {
					ev.Log.AddLog(tk.Sprintf("Error on Save : %s", e.Error()), sError)
				}
			}
		}()
	}

	data := TurbulenceIntensityRaw{}
	_data := FetchScada{}
	maxTimestamp := time.Time{}

loopFetchScada:
	for {
		_data = FetchScada{}
		e = csr.Fetch(&_data, 1, false)
		if e != nil {
			break loopFetchScada
		}
		data = TurbulenceIntensityRaw{}
		data.Projectname = _data.Projectname
		data.Turbine = _data.Turbine
		data.Timestamp = _data.Timestamp.UTC()
		data.DateInfo = _data.DateInfo
		data.WindspeedBin = _data.Windspeed_ms_bin
		data.ID = tk.Sprintf("%s_%s_%s", data.Projectname, data.Turbine, data.Timestamp.Format("20060102150405"))

		if data.Timestamp.After(maxTimestamp) { /* get max timestamp for each project */
			maxTimestamp = data.Timestamp
		}

		data.WindSpeed = _data.Windspeed_ms
		data.WindSpeedStdDev = _data.Windspeed_ms_stddev
		data.Type = "SCADA"

		dataChan <- data
	}

	close(dataChan)
	wg.Wait()

	ev.updateLastData(projectname, "SCADA", maxTimestamp)
}

func (ev *TurbulenceIntensityGenerator) projectWorkerMet(projectname string, lastUpdate time.Time, wgProject *sync.WaitGroup) {
	defer wgProject.Done()

	csr, e := ev.Ctx.Connection.NewQuery().
		From(new(MetTower).TableName()).
		Select("projectname", "timestamp", "dateinfo", "windspeedbin", "vhubws90mavg", "vhubws90mstddev").
		Where(dbox.And(dbox.Gte("timestamp", lastUpdate),
			dbox.Eq("projectname", projectname),
			dbox.Gte("windspeedbin", 0),
			dbox.Lte("windspeedbin", 25))).
		Cursor(nil)
	if e != nil {
		ev.Log.AddLog(tk.Sprintf("Error on cursor : %s", e.Error()), sError)
	}
	defer csr.Close()

	var wg sync.WaitGroup
	totalData := csr.Count()
	totalWorker := 4
	dataChan := make(chan TurbulenceIntensityRaw, totalData)

	wg.Add(totalWorker)
	for i := 0; i < totalWorker; i++ {
		go func() {
			defer wg.Done()
			ctxWorker, e := PrepareConnection()
			if e != nil {
				ev.Log.AddLog(e.Error(), sError)
			}
			defer ctxWorker.Close()
			csrSave := ctxWorker.NewQuery().SetConfig("multiexec", true).
				From(new(TurbulenceIntensityRaw).TableName()).Save()
			defer csrSave.Close()
			for data := range dataChan {
				e = csrSave.Exec(tk.M{"data": data})
				if e != nil {
					ev.Log.AddLog(tk.Sprintf("Error on Save : %s", e.Error()), sError)
				}
			}
		}()
	}

	data := TurbulenceIntensityRaw{}
	_data := FetchMet{}
	maxTimestamp := time.Time{}

loopFetchMet:
	for {
		_data = FetchMet{}
		e = csr.Fetch(&_data, 1, false)
		if e != nil {
			break loopFetchMet
		}
		data = TurbulenceIntensityRaw{}
		data.Projectname = _data.Projectname
		data.Timestamp = _data.Timestamp.UTC()
		data.DateInfo = _data.DateInfo
		data.WindspeedBin = _data.Windspeedbin
		data.ID = tk.Sprintf("%s_%s", data.Projectname, data.Timestamp.Format("20060102"))

		if data.Timestamp.After(maxTimestamp) { /* get max timestamp for each project */
			maxTimestamp = data.Timestamp
		}

		data.WindSpeed = _data.Vhubws90mavg
		data.WindSpeedStdDev = _data.Vhubws90mstddev
		data.Type = "MET"

		dataChan <- data
	}

	close(dataChan)
	wg.Wait()

	ev.updateLastData(projectname, "MET", maxTimestamp)
}

func (ev *TurbulenceIntensityGenerator) processDataMet(wgMet *sync.WaitGroup) {
	defer wgMet.Done()
	t0 := time.Now()

	projectList := ev.getProjectList()
	lastUpdatePerProject := ev.getLatestData("MET")

	var wg sync.WaitGroup
	wg.Add(len(projectList))
	for _, _project := range projectList {
		go ev.projectWorkerMet(_project, lastUpdatePerProject[_project], &wg)
	}
	wg.Wait()

	ev.Log.AddLog(tk.Sprintf("Duration process met tower data %f minutes", time.Since(t0).Minutes()), sInfo)
}
