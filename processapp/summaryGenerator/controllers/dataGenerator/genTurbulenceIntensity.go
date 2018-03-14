package generatorControllers

import (
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

func (ev *TurbulenceIntensitySummary) getProjectList() (result []string) {
	ev.Log.AddLog("Get Project List", sInfo)

	projectData := []tk.M{}
	csrt, e := ev.Ctx.Connection.NewQuery().
		From("ref_project").Cursor(nil)
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
	ev.Log.AddLog("Finish getting Project List", sInfo)

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
		result[val.Projectname] = val.LastUpdate
	}
	ev.Log.AddLog("Finish getting latest data for each turbine", sInfo)

	return
}

func (ev *TurbulenceIntensitySummary) updateLastData(projectname, tipe string, maxTimeStamp time.Time) {
	if !maxTimeStamp.IsZero() {
		data := LatestTurbulence{}
		data.Projectname = projectname
		data.ID = tk.Sprintf("%s_%s", data.Projectname, tipe)
		data.LastUpdate = maxTimeStamp
		data.Type = tipe

		e := ev.Ctx.Connection.NewQuery().SetConfig("multiexec", true).
			From(new(LatestTurbulence).TableName()).Save().Exec(tk.M{"data": data})

		if e != nil {
			ev.Log.AddLog(tk.Sprintf("Error on Save at updateLastData due to : %s", e.Error()), sError)
		}
	}
	ev.Log.AddLog(tk.Sprintf("Finish updating last data for %s on %s at %s", projectname, tipe, maxTimeStamp.String()), sInfo)
}

func (ev *TurbulenceIntensitySummary) processDataScada(wgScada *sync.WaitGroup) {
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

func (ev *TurbulenceIntensitySummary) projectWorker(projectname string, lastUpdate time.Time, wgProject *sync.WaitGroup) {
	defer wgProject.Done()

	countWS := tk.M{"$cond": tk.M{}.
		Set("if", tk.M{
			"$and": []tk.M{
				tk.M{"$ifNull": []interface{}{"$windspeed", false}},
				tk.M{"$gte": []interface{}{"$windspeed", -200}},
			},
		}).
		Set("then", 1).
		Set("else", 0)}
	countWSStd := tk.M{"$cond": tk.M{}.
		Set("if", tk.M{
			"$and": []tk.M{
				tk.M{"$ifNull": []interface{}{"$windspeedstddev", false}},
				tk.M{"$gte": []interface{}{"$windspeedstddev", -200}},
			},
		}).
		Set("then", 1).
		Set("else", 0)}

	ev.Log.AddLog(tk.Sprintf("Update data %s from %s", projectname, lastUpdate.String()), sInfo)
	pipe := []tk.M{
		tk.M{"$match": tk.M{
			"$and": []tk.M{
				tk.M{"dateinfo.dateid": tk.M{"$gte": lastUpdate}},
				tk.M{"projectname": projectname},
				tk.M{"type": "SCADA"},
			},
		}},
		tk.M{"$group": tk.M{
			"_id": tk.M{
				"projectname":  "$projectname",
				"turbine":      "$turbine",
				"windspeedbin": "$windspeedbin",
				"timestamp":    "$dateinfo.dateid",
			},
			"windspeedtotal":    tk.M{"$sum": "$windspeed"},
			"windspeedstdtotal": tk.M{"$sum": "$windspeedstddev"},
			"windspeedcount":    tk.M{"$sum": countWS},
			"windspeedstdcount": tk.M{"$sum": countWSStd},
		}},
	}

	csr, e := ev.Ctx.Connection.NewQuery().
		From("TurbulenceIntensity10Min").
		Command("pipe", pipe).Cursor(nil)
	if e != nil {
		ev.Log.AddLog(tk.Sprintf("Error on cursor : %s", e.Error()), sError)
	}
	defer csr.Close()

	turbulenceData := []tk.M{}
	e = csr.Fetch(&turbulenceData, 0, false)
	if e != nil {
		ev.Log.AddLog(tk.Sprintf("Error on fetch : %s", e.Error()), sError)
	}

	var wg sync.WaitGroup
	totalData := len(turbulenceData)
	totalWorker := 4
	dataChan := make(chan TurbulenceIntensity, totalData)

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
				From(new(TurbulenceIntensity).TableName()).Save()
			defer csrSave.Close()
			for data := range dataChan {
				e = csrSave.Exec(tk.M{"data": data})
				if e != nil {
					ev.Log.AddLog(tk.Sprintf("Error on Save : %s", e.Error()), sError)
				}
			}
		}()
	}

	data := TurbulenceIntensity{}
	maxTimeStamp := time.Time{}

	for _, _data := range turbulenceData {
		data = TurbulenceIntensity{}
		ids := _data.Get("_id", tk.M{}).(tk.M)
		data.Projectname = ids.GetString("projectname")
		data.Turbine = ids.GetString("turbine")
		data.Timestamp = ids.Get("timestamp", time.Time{}).(time.Time).UTC()
		data.WindspeedBin = ids.GetFloat64("windspeedbin")
		data.ID = tk.Sprintf("%s_%s_%s_%s", data.Projectname, data.Turbine, tk.Sprintf("%.1f", data.WindspeedBin), data.Timestamp.Format("20060102"))

		data.WindSpeedTotal = _data.GetFloat64("windspeedtotal")
		data.WindSpeedStdTotal = _data.GetFloat64("windspeedstdtotal")
		data.WindSpeedCount = _data.GetFloat64("windspeedcount")
		data.WindSpeedStdCount = _data.GetFloat64("windspeedstdcount")
		data.Type = "SCADA"

		if data.Timestamp.After(maxTimeStamp) {
			maxTimeStamp = data.Timestamp
		}

		dataChan <- data
	}

	close(dataChan)
	wg.Wait()

	ev.updateLastData(projectname, "SCADA", maxTimeStamp)
}

func (ev *TurbulenceIntensitySummary) projectWorkerMet(projectname string, lastupdate time.Time, wgProject *sync.WaitGroup) {
	defer wgProject.Done()

	countWS := tk.M{"$cond": tk.M{}.
		Set("if", tk.M{
			"$and": []tk.M{
				tk.M{"$ifNull": []interface{}{"$windspeed", false}},
				tk.M{"$gte": []interface{}{"$windspeed", -200}},
			},
		}).
		Set("then", 1).
		Set("else", 0)}
	countWSStd := tk.M{"$cond": tk.M{}.
		Set("if", tk.M{
			"$and": []tk.M{
				tk.M{"$ifNull": []interface{}{"$windspeedstddev", false}},
				tk.M{"$gte": []interface{}{"$windspeedstddev", -200}},
			},
		}).
		Set("then", 1).
		Set("else", 0)}
	pipe := []tk.M{
		tk.M{"$match": tk.M{
			"$and": []tk.M{
				tk.M{"dateinfo.dateid": tk.M{"$gte": lastupdate}},
				tk.M{"projectname": projectname},
				tk.M{"type": "MET"},
			},
		}},
		tk.M{"$group": tk.M{
			"_id": tk.M{
				"projectname":  "$projectname",
				"windspeedbin": "$windspeedbin",
				"timestamp":    "$dateinfo.dateid",
			},
			"windspeedtotal":    tk.M{"$sum": "$windspeed"},
			"windspeedstdtotal": tk.M{"$sum": "$windspeedstddev"},
			"windspeedcount":    tk.M{"$sum": countWS},
			"windspeedstdcount": tk.M{"$sum": countWSStd},
		}},
	}

	turbulenceData := []tk.M{}
	csr, e := ev.Ctx.Connection.NewQuery().
		From("TurbulenceIntensity10Min").
		Command("pipe", pipe).Cursor(nil)
	if e != nil {
		ev.Log.AddLog(tk.Sprintf("Error on cursor : %s", e.Error()), sError)
	}
	defer csr.Close()

	e = csr.Fetch(&turbulenceData, 0, false)
	if e != nil {
		ev.Log.AddLog(tk.Sprintf("Error on Fetch : %s", e.Error()), sError)
	}

	var wg sync.WaitGroup
	totalData := len(turbulenceData)
	totalWorker := 4
	dataChan := make(chan TurbulenceIntensity, totalData)

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
				From(new(TurbulenceIntensity).TableName()).Save()
			defer csrSave.Close()
			for data := range dataChan {
				e = csrSave.Exec(tk.M{"data": data})
				if e != nil {
					ev.Log.AddLog(tk.Sprintf("Error on Save : %s", e.Error()), sError)
				}
			}
		}()
	}

	data := TurbulenceIntensity{}
	maxTimeStamp := time.Time{}

	for _, _data := range turbulenceData {
		data = TurbulenceIntensity{}
		ids := _data.Get("_id", tk.M{}).(tk.M)
		data.Projectname = ids.GetString("projectname")
		data.Timestamp = ids.Get("timestamp", time.Time{}).(time.Time).UTC()
		data.WindspeedBin = ids.GetFloat64("windspeedbin")
		data.ID = tk.Sprintf("%s_%s_%s", data.Projectname, tk.Sprintf("%.1f", data.WindspeedBin), data.Timestamp.Format("20060102"))

		data.WindSpeedTotal = _data.GetFloat64("windspeedtotal")
		data.WindSpeedStdTotal = _data.GetFloat64("windspeedstdtotal")
		data.WindSpeedCount = _data.GetFloat64("windspeedcount")
		data.WindSpeedStdCount = _data.GetFloat64("windspeedstdcount")
		data.Type = "MET"

		if data.Timestamp.After(maxTimeStamp) {
			maxTimeStamp = data.Timestamp
		}

		dataChan <- data
	}

	close(dataChan)
	wg.Wait()

	ev.updateLastData(projectname, "MET", maxTimeStamp)
}

func (ev *TurbulenceIntensitySummary) processDataMet(wgMet *sync.WaitGroup) {
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
