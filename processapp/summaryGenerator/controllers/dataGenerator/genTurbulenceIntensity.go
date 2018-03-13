package generatorControllers

import (
	. "eaciit/wfdemo-git/library/models"
	. "eaciit/wfdemo-git/processapp/summaryGenerator/controllers"
	"github.com/eaciit/dbox"
	_ "github.com/eaciit/dbox/dbc/mongo"
	tk "github.com/eaciit/toolkit"
	"sync"
	"time"
)

type TurbulenceIntensitySummary struct {
	*BaseController
}

type TurbulenceIntensity struct {
	ID                string ` bson:"_id" , json:"_id" `
	Projectname       string
	Turbine           string
	Timestamp         time.Time
	WindspeedBin      float64
	WindSpeedTotal    float64
	WindSpeedStdTotal float64
	WindSpeedCount    float64
	WindSpeedStdCount float64
	Type              string
}

func (m *TurbulenceIntensity) TableName() string {
	return "rpt_turbulenceintensity"
}

type LatestTurbulence struct {
	ID          string ` bson:"_id" , json:"_id" `
	Projectname string
	Turbine     string
	LastUpdate  time.Time
	Type        string
}

func (m *LatestTurbulence) TableName() string {
	return "log_latestturbulence"
}

func (ev *TurbulenceIntensitySummary) CreateTurbulenceIntensitySummary(base *BaseController) {
	ev.BaseController = base

	ev.Log.AddLog("===================== Start processing Turbulence Intensity Summary...", sInfo)

	var wg sync.WaitGroup
	wg.Add(2)

	go ev.processDataScada(&wg)
	go ev.processDataMet(&wg)

	wg.Wait()

	ev.Log.AddLog("===================== End processing Turbulence Intensity Summary...", sInfo)
}

func (ev *TurbulenceIntensitySummary) getTurbinePerProject() (result map[string][]string) {
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

func (ev *TurbulenceIntensitySummary) getLatestData(tipe string) (result map[string]time.Time) {
	ev.Log.AddLog("Get latest data for each turbine", sInfo)

	latestData := []LatestTurbulence{}
	csrt, e := ev.Ctx.Connection.NewQuery().
		From(new(LatestTurbulence).TableName()).
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

func (ev *TurbulenceIntensitySummary) updateLastData(projectname, tipe string, turbineList []string) {
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
		From(new(TurbulenceIntensity).TableName()).
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
		From(new(LatestTurbulence).TableName()).Save()
	defer csrSave.Close()

	for _, _turbine := range turbineList {
		data := LatestTurbulence{}
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

func (ev *TurbulenceIntensitySummary) processDataScada(wgScada *sync.WaitGroup) {
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

func (ev *TurbulenceIntensitySummary) projectWorker(projectname string, turbineList []string, lastUpdate map[string]time.Time, wgProject *sync.WaitGroup) {
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

func (ev *TurbulenceIntensitySummary) turbineWorker(projectname, turbine string, lastupdate time.Time, wgTurbine *sync.WaitGroup) {
	defer wgTurbine.Done()

	ev.Log.AddLog(tk.Sprintf("Processing %s (%s) from %s", turbine, projectname, lastupdate.String()), sInfo)
	countWS := tk.M{"$cond": tk.M{}.
		Set("if", tk.M{
			"$and": []tk.M{
				tk.M{"$ifNull": []interface{}{"$windspeed_ms", false}},
				tk.M{"$gte": []interface{}{"$windspeed_ms", -200}},
			},
		}).
		Set("then", 1).
		Set("else", 0)}
	countWSStd := tk.M{"$cond": tk.M{}.
		Set("if", tk.M{
			"$and": []tk.M{
				tk.M{"$ifNull": []interface{}{"$windspeed_ms_stddev", false}},
				tk.M{"$gte": []interface{}{"$windspeed_ms_stddev", -200}},
			},
		}).
		Set("then", 1).
		Set("else", 0)}
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
		tk.M{"$group": tk.M{
			"_id": tk.M{
				"projectname":  "$projectname",
				"turbine":      "$turbine",
				"windspeedbin": "$windspeed_ms_bin",
				"timestamp":    "$dateinfo.dateid",
			},
			"windspeedtotal":    tk.M{"$sum": "$windspeed_ms"},
			"windspeedstdtotal": tk.M{"$sum": "$windspeed_ms_stddev"},
			"windspeedcount":    tk.M{"$sum": countWS},
			"windspeedstdcount": tk.M{"$sum": countWSStd},
		}},
	}

	turbulenceData := []tk.M{}
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

	data := TurbulenceIntensity{}

	csrSave := ev.Ctx.Connection.NewQuery().SetConfig("multiexec", true).
		From(new(TurbulenceIntensity).TableName()).Save()
	defer csrSave.Close()

	for _, val := range turbulenceData {
		data = TurbulenceIntensity{}
		ids := val.Get("_id", tk.M{}).(tk.M)
		data.Projectname = ids.GetString("projectname")
		data.Turbine = ids.GetString("turbine")
		data.Timestamp = ids.Get("timestamp", time.Time{}).(time.Time).UTC()
		data.WindspeedBin = ids.GetFloat64("windspeedbin")
		data.ID = tk.Sprintf("%s_%s_%s_%s", data.Projectname, data.Turbine, tk.Sprintf("%.1f", data.WindspeedBin), data.Timestamp.Format("20060102"))

		data.WindSpeedTotal = val.GetFloat64("windspeedtotal")
		data.WindSpeedStdTotal = val.GetFloat64("windspeedstdtotal")
		data.WindSpeedCount = val.GetFloat64("windspeedcount")
		data.WindSpeedStdCount = val.GetFloat64("windspeedstdcount")
		data.Type = "SCADA"

		e = csrSave.Exec(tk.M{"data": data})
		if e != nil {
			ev.Log.AddLog(tk.Sprintf("Error on Save : %s", e.Error()), sError)
		}
	}
}

func (ev *TurbulenceIntensitySummary) projectWorkerMet(projectname string, lastupdate time.Time, wgProject *sync.WaitGroup) {
	defer wgProject.Done()

	countWS := tk.M{"$cond": tk.M{}.
		Set("if", tk.M{
			"$and": []tk.M{
				tk.M{"$ifNull": []interface{}{"$vhubws90mavg", false}},
				tk.M{"$gte": []interface{}{"$vhubws90mavg", -200}},
			},
		}).
		Set("then", 1).
		Set("else", 0)}
	countWSStd := tk.M{"$cond": tk.M{}.
		Set("if", tk.M{
			"$and": []tk.M{
				tk.M{"$ifNull": []interface{}{"$vhubws90mstddev", false}},
				tk.M{"$gte": []interface{}{"$vhubws90mstddev", -200}},
			},
		}).
		Set("then", 1).
		Set("else", 0)}
	pipe := []tk.M{
		tk.M{"$match": tk.M{
			"$and": []tk.M{
				tk.M{"timestamp": tk.M{"$gte": lastupdate}},
				tk.M{"projectname": projectname},
				tk.M{"windspeedbin": tk.M{"$gte": 0}},
				tk.M{"windspeedbin": tk.M{"$lte": 25}},
			},
		}},
		tk.M{"$group": tk.M{
			"_id": tk.M{
				"projectname":  "$projectname",
				"windspeedbin": "$windspeedbin",
				"timestamp":    "$dateinfo.dateid",
			},
			"windspeedtotal":    tk.M{"$sum": "$vhubws90mavg"},
			"windspeedstdtotal": tk.M{"$sum": "$vhubws90mstddev"},
			"windspeedcount":    tk.M{"$sum": countWS},
			"windspeedstdcount": tk.M{"$sum": countWSStd},
		}},
	}

	turbulenceData := []tk.M{}
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

	data := TurbulenceIntensity{}

	csrSave := ev.Ctx.Connection.NewQuery().SetConfig("multiexec", true).
		From(new(TurbulenceIntensity).TableName()).Save()
	defer csrSave.Close()

	for _, val := range turbulenceData {
		data = TurbulenceIntensity{}
		ids := val.Get("_id", tk.M{}).(tk.M)
		data.Projectname = ids.GetString("projectname")
		data.Timestamp = ids.Get("timestamp", time.Time{}).(time.Time).UTC()
		data.WindspeedBin = ids.GetFloat64("windspeedbin")
		data.ID = tk.Sprintf("%s_%s_%s", data.Projectname, tk.Sprintf("%.1f", data.WindspeedBin), data.Timestamp.Format("20060102"))

		data.WindSpeedTotal = val.GetFloat64("windspeedtotal")
		data.WindSpeedStdTotal = val.GetFloat64("windspeedstdtotal")
		data.WindSpeedCount = val.GetFloat64("windspeedcount")
		data.WindSpeedStdCount = val.GetFloat64("windspeedstdcount")
		data.Type = "MET"

		e = csrSave.Exec(tk.M{"data": data})
		if e != nil {
			ev.Log.AddLog(tk.Sprintf("Error on Save : %s", e.Error()), sError)
		}
	}

	ev.updateLastData(projectname, "MET", []string{})
}

func (ev *TurbulenceIntensitySummary) processDataMet(wgMet *sync.WaitGroup) {
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
