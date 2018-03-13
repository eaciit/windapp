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
	Turbine     string
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

	// go ev.processDataScada(&wg)
	go ev.processInitialDataScada(&wg)
	go ev.processDataMet(&wg)

	wg.Wait()

	ev.Log.AddLog("===================== End processing Turbulence 10 Min...", sInfo)
}

func (ev *TurbulenceIntensityGenerator) processInitialDataScada(wgScada *sync.WaitGroup) {
	defer wgScada.Done()

	t0 := time.Now()
	turbinePerProject := ev.getTurbinePerProject()
	lastUpdateTurbine := ev.getLatestData("SCADA")
	lastUpdatePerDay := map[string]time.Time{}

	tStart := time.Date(2018, 3, 1, 0, 0, 0, 0, time.UTC)
	tNow := time.Now()
	for {
		if tStart.After(tNow) {
			break
		}
		for keys := range lastUpdateTurbine {
			lastUpdatePerDay[keys] = tStart.UTC()
		}
		var wg sync.WaitGroup
		wg.Add(len(turbinePerProject))
		ev.Log.AddLog(tk.Sprintf("Updating data for %s", tStart.UTC().String()), sInfo)
		for _project, _turbine := range turbinePerProject {
			go ev.projectInitialWorker(_project, _turbine, lastUpdatePerDay, &wg)
		}
		wg.Wait()

		tStart = tStart.AddDate(0, 0, 1)
	}

	ev.Log.AddLog(tk.Sprintf("Duration processing scada data %f minutes", time.Since(t0).Minutes()), sInfo)
}

func (ev *TurbulenceIntensityGenerator) projectInitialWorker(projectname string, turbineList []string, lastUpdate map[string]time.Time, wgProject *sync.WaitGroup) {
	defer wgProject.Done()
	var wg sync.WaitGroup
	wg.Add(len(turbineList))
	for _, _turbine := range turbineList {
		keys := tk.Sprintf("%s_%s", projectname, _turbine)
		go ev.turbineInitialWorker(projectname, _turbine, lastUpdate[keys], &wg)
	}
	wg.Wait()
	ev.updateLastData(projectname, "SCADA", turbineList)
}

func (ev *TurbulenceIntensityGenerator) turbineInitialWorker(projectname, turbine string, lastupdate time.Time, wgTurbine *sync.WaitGroup) {
	defer wgTurbine.Done()
	pipe := []tk.M{
		tk.M{"$match": tk.M{
			"$and": []tk.M{
				tk.M{"dateinfo.dateid": lastupdate},
				tk.M{"projectname": projectname},
				tk.M{"turbine": turbine},
				tk.M{"isnull": false},
				tk.M{"windspeed_ms_bin": tk.M{"$gte": 0}},
				tk.M{"windspeed_ms_bin": tk.M{"$lte": 25}},
			},
		}},
		tk.M{"$project": tk.M{
			"projectname":         1,
			"turbine":             1,
			"timestamp":           1,
			"dateinfo":            1,
			"windspeed_ms":        1,
			"windspeed_ms_bin":    1,
			"windspeed_ms_stddev": 1,
		}},
	}

	// turbulenceData := []FetchScada{}
	csr, e := ev.Ctx.Connection.NewQuery().
		From("Scada10MinHFD").
		Command("pipe", pipe).Cursor(nil)
	if e != nil {
		ev.Log.AddLog(tk.Sprintf("Error on cursor : %s", e.Error()), sError)
	}
	defer csr.Close()

	data := TurbulenceIntensityRaw{}

	csrSave := ev.Ctx.Connection.NewQuery().SetConfig("multiexec", true).
		From(new(TurbulenceIntensityRaw).TableName()).Save()
	defer csrSave.Close()
	_data := FetchScada{}

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

		data.WindSpeed = _data.Windspeed_ms
		data.WindSpeedStdDev = _data.Windspeed_ms_stddev
		data.Type = "SCADA"

		e = csrSave.Exec(tk.M{"data": data})
		if e != nil {
			ev.Log.AddLog(tk.Sprintf("Error on Save : %s", e.Error()), sError)
		}
	}

	/*e = csr.Fetch(&turbulenceData, 0, false)
	if e != nil {
		ev.Log.AddLog(tk.Sprintf("Error on Fetch : %s", e.Error()), sError)
	}*/

	/*for _, val := range turbulenceData {

	}*/
}

func (ev *TurbulenceIntensityGenerator) getTurbinePerProject() (result map[string][]string) {
	ev.Log.AddLog("Get Turbine Per Project", sInfo)

	turbineData := []tk.M{}
	csrt, e := ev.Ctx.Connection.NewQuery().
		From("ref_turbine").Cursor(nil)
	if e != nil {
		ev.Log.AddLog(tk.Sprintf("Error on cursor at getTurbinePerProjectFunc due to : %s", e.Error()), sError)
		return
	}
	defer csrt.Close()
	e = csrt.Fetch(&turbineData, 0, false)
	if e != nil {
		ev.Log.AddLog(tk.Sprintf("Error on fetch at getTurbinePerProjectFunc due to : %s", e.Error()), sError)
		return
	}
	result = map[string][]string{}
	for _, val := range turbineData {
		result[val.GetString("project")] = append(result[val.GetString("project")], val.GetString("turbineid"))
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
		result[val.ID] = val.LastUpdate
	}
	ev.Log.AddLog("Finish getting latest data for each turbine", sInfo)

	return
}

func (ev *TurbulenceIntensityGenerator) updateLastData(projectname, tipe string, turbineList []string) {
	pipes := []tk.M{} /* aggregate for rpt_turbulenceintensity to get max date */
	if tipe == "SCADA" {
		pipes = []tk.M{
			tk.M{"$match": tk.M{"$and": []tk.M{
				tk.M{"projectname": projectname},
				tk.M{"turbine": tk.M{"$in": turbineList}},
				tk.M{"type": tipe},
			}}},
			tk.M{"$group": tk.M{
				"_id":     "$turbine",
				"maxDate": tk.M{"$max": "$timestamp"},
			}},
		}
	} else {
		pipes = []tk.M{
			tk.M{"$match": tk.M{"$and": []tk.M{
				tk.M{"projectname": projectname},
				tk.M{"type": tipe},
			}}},
			tk.M{"$group": tk.M{
				"_id":     "$projectname",
				"maxDate": tk.M{"$max": "$timestamp"},
			}},
		}
		turbineList = []string{projectname}
	}
	csrt, e := ev.Ctx.Connection.NewQuery().
		From(new(TurbulenceIntensityRaw).TableName()).
		Command("pipe", pipes).Cursor(nil)
	if e != nil {
		ev.Log.AddLog(tk.Sprintf("Error on cursor at updateLastData due to : %s", e.Error()), sError)
		return
	}
	defer csrt.Close()

	latestData := []tk.M{}
	e = csrt.Fetch(&latestData, 0, false)
	if e != nil {
		ev.Log.AddLog(tk.Sprintf("Error on fetch at updateLastData due to : %s", e.Error()), sError)
		return
	}
	timestampPerTurbine := map[string]time.Time{}
	for _, val := range latestData {
		timestampPerTurbine[val.GetString("_id")] = val.Get("maxDate", time.Time{}).(time.Time).UTC()
	}

	csrSave := ev.Ctx.Connection.NewQuery().SetConfig("multiexec", true).
		From(new(LatestTurbulenceRaw).TableName()).Save()
	defer csrSave.Close()

	for _, _turbine := range turbineList {
		data := LatestTurbulenceRaw{}
		data.Projectname = projectname
		if tipe == "SCADA" {
			data.Turbine = _turbine
			data.ID = tk.Sprintf("%s_%s", data.Projectname, data.Turbine)
		} else {
			data.ID = tk.Sprintf("%s_MET", data.Projectname)
		}
		data.LastUpdate = timestampPerTurbine[_turbine]
		data.Type = tipe
		e = csrSave.Exec(tk.M{"data": data})
		if e != nil {
			ev.Log.AddLog(tk.Sprintf("Error on Save at updateLastData due to : %s", e.Error()), sError)
		}
	}
	ev.Log.AddLog(tk.Sprintf("Finish updating last data for %s on %s", projectname, tipe), sInfo)
}

func (ev *TurbulenceIntensityGenerator) processDataScada(wgScada *sync.WaitGroup) {
	defer wgScada.Done()
	t0 := time.Now()
	turbinePerProject := ev.getTurbinePerProject()
	lastUpdateTurbine := ev.getLatestData("SCADA")

	var wg sync.WaitGroup
	wg.Add(len(turbinePerProject))
	for _project, _turbine := range turbinePerProject {
		go ev.projectWorker(_project, _turbine, lastUpdateTurbine, &wg)
	}
	wg.Wait()

	ev.Log.AddLog(tk.Sprintf("Duration processing scada data %f minutes", time.Since(t0).Minutes()), sInfo)
}

func (ev *TurbulenceIntensityGenerator) projectWorker(projectname string, turbineList []string, lastUpdate map[string]time.Time, wgProject *sync.WaitGroup) {
	defer wgProject.Done()
	var wg sync.WaitGroup
	wg.Add(len(turbineList))
	for _, _turbine := range turbineList {
		keys := tk.Sprintf("%s_%s", projectname, _turbine)
		go ev.turbineWorker(projectname, _turbine, lastUpdate[keys], &wg)
	}
	wg.Wait()
	ev.updateLastData(projectname, "SCADA", turbineList)
}

func (ev *TurbulenceIntensityGenerator) turbineWorker(projectname, turbine string, lastupdate time.Time, wgTurbine *sync.WaitGroup) {
	defer wgTurbine.Done()
	pipe := []tk.M{
		tk.M{"$match": tk.M{
			"$and": []tk.M{
				tk.M{"timestamp": tk.M{"$gte": lastupdate}},
				tk.M{"projectname": projectname},
				tk.M{"turbine": turbine},
				tk.M{"isnull": false},
				tk.M{"windspeed_ms_bin": tk.M{"$gte": 0}},
				tk.M{"windspeed_ms_bin": tk.M{"$lte": 25}},
			},
		}},
		tk.M{"$project": tk.M{
			"projectname":         1,
			"turbine":             1,
			"timestamp":           1,
			"dateinfo":            1,
			"windspeed_ms":        1,
			"windspeed_ms_bin":    1,
			"windspeed_ms_stddev": 1,
		}},
	}

	turbulenceData := []FetchScada{}
	csr, e := ev.Ctx.Connection.NewQuery().
		From("Scada10MinHFD").
		Command("pipe", pipe).Cursor(nil)
	if e != nil {
		ev.Log.AddLog(tk.Sprintf("Error on cursor : %s", e.Error()), sError)
	}
	defer csr.Close()

	e = csr.Fetch(&turbulenceData, 0, false)
	if e != nil {
		ev.Log.AddLog(tk.Sprintf("Error on Fetch : %s", e.Error()), sError)
	}

	data := TurbulenceIntensityRaw{}

	csrSave := ev.Ctx.Connection.NewQuery().SetConfig("multiexec", true).
		From(new(TurbulenceIntensityRaw).TableName()).Save()
	defer csrSave.Close()

	for _, val := range turbulenceData {
		data = TurbulenceIntensityRaw{}
		data.Projectname = val.Projectname
		data.Turbine = val.Turbine
		data.Timestamp = val.Timestamp.UTC()
		data.DateInfo = val.DateInfo
		data.WindspeedBin = val.Windspeed_ms_bin
		data.ID = tk.Sprintf("%s_%s_%s", data.Projectname, data.Turbine, data.Timestamp.Format("20060102150405"))

		data.WindSpeed = val.Windspeed_ms
		data.WindSpeedStdDev = val.Windspeed_ms_stddev
		data.Type = "SCADA"

		e = csrSave.Exec(tk.M{"data": data})
		if e != nil {
			ev.Log.AddLog(tk.Sprintf("Error on Save : %s", e.Error()), sError)
		}
	}
}

func (ev *TurbulenceIntensityGenerator) projectWorkerMet(projectname string, lastupdate time.Time, wgProject *sync.WaitGroup) {
	defer wgProject.Done()

	pipe := []tk.M{
		tk.M{"$match": tk.M{
			"$and": []tk.M{
				tk.M{"timestamp": tk.M{"$gte": lastupdate}},
				tk.M{"projectname": projectname},
				tk.M{"windspeedbin": tk.M{"$gte": 0}},
				tk.M{"windspeedbin": tk.M{"$lte": 25}},
			},
		}},
		tk.M{"$project": tk.M{
			"projectname":     1,
			"timestamp":       1,
			"dateinfo":        1,
			"windspeedbin":    1,
			"vhubws90mavg":    1,
			"vhubws90mstddev": 1,
		}},
	}

	turbulenceData := []FetchMet{}
	csr, e := ev.Ctx.Connection.NewQuery().
		From(new(MetTower).TableName()).
		Command("pipe", pipe).Cursor(nil)
	if e != nil {
		ev.Log.AddLog(tk.Sprintf("Error on cursor : %s", e.Error()), sError)
	}
	defer csr.Close()

	e = csr.Fetch(&turbulenceData, 0, false)
	if e != nil {
		ev.Log.AddLog(tk.Sprintf("Error on Fetch : %s", e.Error()), sError)
	}

	data := TurbulenceIntensityRaw{}

	csrSave := ev.Ctx.Connection.NewQuery().SetConfig("multiexec", true).
		From(new(TurbulenceIntensityRaw).TableName()).Save()
	defer csrSave.Close()

	for _, val := range turbulenceData {
		data = TurbulenceIntensityRaw{}

		data.Projectname = val.Projectname
		data.Timestamp = val.Timestamp.UTC()
		data.DateInfo = val.DateInfo
		data.WindspeedBin = val.Windspeedbin
		data.ID = tk.Sprintf("%s_%s", data.Projectname, data.Timestamp.Format("20060102"))

		data.WindSpeed = val.Vhubws90mavg
		data.WindSpeedStdDev = val.Vhubws90mstddev
		data.Type = "MET"

		e = csrSave.Exec(tk.M{"data": data})
		if e != nil {
			ev.Log.AddLog(tk.Sprintf("Error on Save : %s", e.Error()), sError)
		}
	}

	ev.updateLastData(projectname, "MET", []string{})
}

func (ev *TurbulenceIntensityGenerator) processDataMet(wgMet *sync.WaitGroup) {
	defer wgMet.Done()
	t0 := time.Now()

	turbinePerProject := ev.getTurbinePerProject()
	lastUpdateTurbine := ev.getLatestData("MET")

	var wg sync.WaitGroup
	wg.Add(len(turbinePerProject))
	for _project := range turbinePerProject {
		keys := _project + "_MET"
		go ev.projectWorkerMet(_project, lastUpdateTurbine[keys], &wg)
	}
	wg.Wait()

	ev.Log.AddLog(tk.Sprintf("Duration process met tower data %f minutes", time.Since(t0).Minutes()), sInfo)
}
