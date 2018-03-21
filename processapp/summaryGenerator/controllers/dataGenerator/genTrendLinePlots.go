package generatorControllers

import (
	. "eaciit/wfdemo-git/processapp/summaryGenerator/controllers"
	"github.com/eaciit/dbox"
	_ "github.com/eaciit/dbox/dbc/mongo"
	tk "github.com/eaciit/toolkit"
	"strings"
	"sync"
	"time"
)

type TrendLinePlotSummary struct {
	*BaseController
}

type LatestTrendLine struct {
	ID          string ` bson:"_id" , json:"_id" `
	Projectname string
	LastUpdate  time.Time
	Type        string
}

func (m *LatestTrendLine) TableName() string {
	return "log_latesttrendline"
}

func (ev *TrendLinePlotSummary) CreateTrendLinePlotSummary(base *BaseController) {
	ev.BaseController = base

	ev.Log.AddLog("===================== Start processing Trend Line Plots Summary...", sInfo)

	var wg sync.WaitGroup
	wg.Add(2)

	go ev.processDataScada(&wg)
	go ev.processDataMet(&wg)

	wg.Wait()

	ev.Log.AddLog("===================== End processing Trend Line Plots Summary...", sInfo)
}

func (ev *TrendLinePlotSummary) getProjectList() (result []string) {
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

func (ev *TrendLinePlotSummary) getTemperatureField() (result map[string][]string) {
	csr, e := ev.Ctx.Connection.NewQuery().
		From("ref_databrowsertag").
		Order("projectname", "label").
		Cursor(nil)
	defer csr.Close()

	if e != nil {
		ev.Log.AddLog(tk.Sprintf("Error on cursor at getTemperatureField due to : %s", e.Error()), sError)
		return
	}

	_data := tk.M{}
	lastProject := ""
	currProject := ""
	tempList := []string{}
	result = map[string][]string{}
	for {
		_data = tk.M{}
		e = csr.Fetch(&_data, 1, false)
		if e != nil {
			break
		}
		currProject = _data.GetString("projectname")
		if lastProject != currProject {
			if lastProject != "" {
				result[lastProject] = tempList
				tempList = []string{}
			}
			lastProject = currProject
		}
		if strings.Contains(strings.ToLower(_data.GetString("realtimefield")), "temp") {
			tempList = append(tempList, strings.ToLower(_data.GetString("realtimefield")))
		}
	}
	if lastProject != "" {
		result[lastProject] = tempList
	}

	return
}

func (ev *TrendLinePlotSummary) getLatestData(tipe string) (result map[string]time.Time) {
	ev.Log.AddLog("Get latest data for each turbine", sInfo)

	latestData := []LatestTrendLine{}
	csrt, e := ev.Ctx.Connection.NewQuery().
		From(new(LatestTrendLine).TableName()).
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

func (ev *TrendLinePlotSummary) updateLastData(projectname, tipe string, maxTimeStamp time.Time) {
	if !maxTimeStamp.IsZero() {
		data := LatestTrendLine{}
		data.Projectname = projectname
		data.ID = tk.Sprintf("%s_%s", data.Projectname, tipe)
		data.LastUpdate = maxTimeStamp
		data.Type = tipe

		e := ev.Ctx.Connection.NewQuery().SetConfig("multiexec", true).
			From(new(LatestTrendLine).TableName()).Save().Exec(tk.M{"data": data})

		if e != nil {
			ev.Log.AddLog(tk.Sprintf("Error on Save at updateLastData due to : %s", e.Error()), sError)
		}
	}
	ev.Log.AddLog(tk.Sprintf("Finish updating last data for %s on %s at %s", projectname, tipe, maxTimeStamp.String()), sInfo)
}

func (ev *TrendLinePlotSummary) processDataScada(wgScada *sync.WaitGroup) {
	defer wgScada.Done()

	t0 := time.Now()
	projectList := ev.getProjectList()
	lastUpdatePerProject := ev.getLatestData("SCADA")
	temperatureList := ev.getTemperatureField()

	var wg sync.WaitGroup
	wg.Add(len(projectList))
	for _, _project := range projectList {
		go ev.projectWorker(_project, lastUpdatePerProject[_project], temperatureList[_project], &wg)
	}
	wg.Wait()

	ev.Log.AddLog(tk.Sprintf("Duration processing scada data %f minutes", time.Since(t0).Minutes()), sInfo)
}

func (ev *TrendLinePlotSummary) projectWorker(projectname string, lastUpdate time.Time, tempList []string, wgProject *sync.WaitGroup) {
	defer wgProject.Done()

	groups := tk.M{"_id": tk.M{
		"turbine":   "$turbine",
		"timestamp": "$dateinfo.dateid",
	}}
	for _, field := range tempList {
		totalName := field + "total"
		countName := field + "count"
		fieldName := "$" + field

		countCondition := tk.M{"$cond": tk.M{}.
			Set("if", tk.M{
				"$and": []tk.M{
					tk.M{"$ifNull": []interface{}{fieldName, false}},
					tk.M{"$lte": []interface{}{fieldName, 200}},
				},
			}).
			Set("then", 1).
			Set("else", 0)}
		groups.Set(totalName, tk.M{"$sum": fieldName})
		groups.Set(countName, tk.M{"$sum": countCondition})
	}

	ev.Log.AddLog(tk.Sprintf("Update data %s from %s", projectname, lastUpdate.String()), sInfo)
	pipe := []tk.M{
		tk.M{"$match": tk.M{
			"$and": []tk.M{
				tk.M{"dateinfo.dateid": tk.M{"$gte": lastUpdate}},
				tk.M{"projectname": projectname},
				tk.M{"isnull": false},
			},
		}},
		tk.M{"$group": groups},
	}

	csr, e := ev.Ctx.Connection.NewQuery().
		From("Scada10MinHFD").
		Command("pipe", pipe).Cursor(nil)
	if e != nil {
		ev.Log.AddLog(tk.Sprintf("Error on cursor : %s", e.Error()), sError)
	}
	defer csr.Close()

	trendLineData := []tk.M{}
	e = csr.Fetch(&trendLineData, 0, false)
	if e != nil {
		ev.Log.AddLog(tk.Sprintf("Error on fetch : %s", e.Error()), sError)
	}

	var wg sync.WaitGroup
	totalData := len(trendLineData)
	totalWorker := 4
	dataChan := make(chan tk.M, totalData)

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
				From("rpt_trendlineplot").Save()
			defer csrSave.Close()
			for data := range dataChan {
				e = csrSave.Exec(tk.M{"data": data})
				if e != nil {
					ev.Log.AddLog(tk.Sprintf("Error on Save : %s", e.Error()), sError)
				}
			}
		}()
	}

	data := tk.M{}
	maxTimeStamp := time.Time{}

	for _, _data := range trendLineData {
		data = tk.M{}
		ids := _data.Get("_id", tk.M{}).(tk.M)
		timestamp := ids.Get("timestamp", time.Time{}).(time.Time).UTC()
		data.Set("projectname", projectname)
		data.Set("turbine", ids.GetString("turbine"))
		data.Set("timestamp", timestamp)
		_data.Set("_id", tk.Sprintf("%s_%s_%s", projectname, ids.GetString("turbine"), timestamp.Format("20060102")))
		for _dataKey, _dataVal := range _data {
			data.Set(_dataKey, _dataVal)
		}
		data.Set("type", "SCADA")

		if timestamp.After(maxTimeStamp) {
			maxTimeStamp = timestamp
		}

		dataChan <- data
	}

	close(dataChan)
	wg.Wait()

	ev.updateLastData(projectname, "SCADA", maxTimeStamp)
}

func (ev *TrendLinePlotSummary) projectWorkerMet(projectname string, lastupdate time.Time, wgProject *sync.WaitGroup) {
	defer wgProject.Done()

	pipe := []tk.M{
		tk.M{"$match": tk.M{
			"$and": []tk.M{
				tk.M{"dateinfo.dateid": tk.M{"$gte": lastupdate}},
				tk.M{"projectname": projectname},
			},
		}},
	}
	fieldName := "$trefhrefhumid855mavg"
	if lastupdate.Before(time.Date(2017, 1, 1, 0, 0, 0, 0, time.UTC)) {
		fieldName = "$trefhreftemp855mavg"
	}
	countCondition := tk.M{"$cond": tk.M{}.
		Set("if", tk.M{
			"$and": []tk.M{
				tk.M{"$ifNull": []interface{}{fieldName, false}},
				tk.M{"$lte": []interface{}{fieldName, 200}},
			},
		}).
		Set("then", 1).
		Set("else", 0)}
	pipe = append(pipe, tk.M{"$group": tk.M{
		"_id":              tk.M{"timestamp": "$dateinfo.dateid"},
		"tempoutdoortotal": tk.M{"$sum": fieldName},
		"tempoutdoorcount": tk.M{"$sum": countCondition},
	}})

	trendLineData := []tk.M{}
	csr, e := ev.Ctx.Connection.NewQuery().
		From("MetTower").
		Command("pipe", pipe).Cursor(nil)
	if e != nil {
		ev.Log.AddLog(tk.Sprintf("Error on cursor : %s", e.Error()), sError)
	}
	defer csr.Close()

	e = csr.Fetch(&trendLineData, 0, false)
	if e != nil {
		ev.Log.AddLog(tk.Sprintf("Error on Fetch : %s", e.Error()), sError)
	}

	var wg sync.WaitGroup
	totalData := len(trendLineData)
	totalWorker := 4
	dataChan := make(chan tk.M, totalData)

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
				From("rpt_trendlineplot").Save()
			defer csrSave.Close()
			for data := range dataChan {
				e = csrSave.Exec(tk.M{"data": data})
				if e != nil {
					ev.Log.AddLog(tk.Sprintf("Error on Save : %s", e.Error()), sError)
				}
			}
		}()
	}

	data := tk.M{}
	maxTimeStamp := time.Time{}

	for _, _data := range trendLineData {
		data = tk.M{}
		ids := _data.Get("_id", tk.M{}).(tk.M)
		timestamp := ids.Get("timestamp", time.Time{}).(time.Time).UTC()
		data.Set("projectname", projectname)
		data.Set("timestamp", timestamp)
		_data.Set("_id", tk.Sprintf("%s_%s", projectname, timestamp.Format("20060102")))
		for _dataKey, _dataVal := range _data {
			data.Set(_dataKey, _dataVal)
		}
		data.Set("type", "MET")

		if timestamp.After(maxTimeStamp) {
			maxTimeStamp = timestamp
		}

		dataChan <- data
	}

	close(dataChan)
	wg.Wait()

	ev.updateLastData(projectname, "MET", maxTimeStamp)
}

func (ev *TrendLinePlotSummary) processDataMet(wgMet *sync.WaitGroup) {
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
