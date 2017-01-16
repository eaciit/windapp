package controller

import (
	. "eaciit/wfdemo-git/library/core"
	. "eaciit/wfdemo-git/library/models"
	"strings"
	"time"

	"eaciit/wfdemo-git/web/helper"

	knot "github.com/eaciit/knot/knot.v1"
	tk "github.com/eaciit/toolkit"
)

type MonitoringController struct {
	App
}

func CreateMonitoringController() *MonitoringController {
	var controller = new(MonitoringController)
	return controller
}

func (m *MonitoringController) GetData(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	p := tk.M{}
	e := k.GetPayload(&p)

	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	// get the last data for monitoring

	turbine := p.Get("turbine").([]interface{})
	project := ""
	if p.GetString("project") != "" {
		anProject := strings.Split(p.GetString("project"), "(")
		project = strings.TrimRight(anProject[0], " ")
	}

	match := tk.M{}
	turbines := map[string]tk.M{}
	var projectList []interface{}
	if project != "" {
		match.Set("project", project)
		projectList = append(projectList, project)
	}

	turbines, _, _ = helper.GetProjectTurbineList(projectList)

	if len(turbine) > 0 {
		match.Set("turbine", tk.M{}.Set("$in", turbine))
	}

	group := tk.M{
		"_id": tk.M{
			"project": "$project",
			"turbine": "$turbine",
		},
		"timestamp":     tk.M{"$last": "$timestamp"},
		"windspeed":     tk.M{"$last": "$windspeed"},
		"production":    tk.M{"$last": "$production"},
		"rotorspeedrpm": tk.M{"$last": "$rotorspeedrpm"},
		"status":        tk.M{"$last": "$status"},
		"statuscode":    tk.M{"$last": "$statuscode"},
		"statusdesc":    tk.M{"$last": "$statusdesc"},
	}

	var pipes []tk.M
	pipes = append(pipes, tk.M{}.Set("$match", match))
	pipes = append(pipes, tk.M{}.Set("$group", group))
	pipes = append(pipes, tk.M{}.Set("$sort", tk.M{
		"_id.project": 1,
		"_id.turbine": 1,
	}))

	csr, e := DB().Connection.NewQuery().
		From(new(Monitoring).TableName()).
		Command("pipe", pipes).
		Cursor(nil)

	defer csr.Close()

	if e != nil {
		return helper.CreateResult(false, nil, "Error query : "+e.Error())
	}

	results := make([]tk.M, 0)
	e = csr.Fetch(&results, 0, false)

	projects := tk.M{}

	for _, v := range results {
		ID := v.Get("_id").(tk.M)
		project := ID.GetString("project")
		turbine := ID.GetString("turbine")
		timestamp := v.Get("timestamp").(time.Time)
		windspeed := checkNAValue(v.GetFloat64("windspeed"))
		production := checkNAValue(v.GetFloat64("production"))
		rotorspeedrpm := checkNAValue(v.GetFloat64("rotorspeedrpm"))
		status := v.GetString("status")
		statuscode := v.GetString("statuscode")
		statusdesc := v.GetString("statusdesc")

		list := []tk.M{}
		newRecord := tk.M{}

		updated := tk.M{}
		capacitymw := 0.0

		if projects.Get(project) != nil {
			if projects.Get(project).(tk.M).Get("turbines") != nil {
				list = projects.Get(project).(tk.M).Get("turbines").([]tk.M)
			}
			updated = projects.Get(project).(tk.M)
		}

		if turbines[project].Get(turbine) != nil {
			capacitymw = turbines[project].Get(turbine).(TurbineMaster).CapacityMW
		}

		newRecord.Set("turbine", turbine)
		newRecord.Set("timestamp", timestamp)
		newRecord.Set("windspeed", windspeed)
		newRecord.Set("production", production)
		newRecord.Set("rotorspeedrpm", rotorspeedrpm)
		newRecord.Set("status", status)
		newRecord.Set("statuscode", statuscode)
		newRecord.Set("statusdesc", statusdesc)

		list = append(list, newRecord)
		updated.Set("turbines", list)
		updated.Set("totalturbines", updated.GetInt("totalturbines")+1)
		updated.Set("totalcap", updated.GetFloat64("totalcap")+capacitymw)
		updated.Set("totalprod", updated.GetFloat64("totalprod")+production)

		projects.Set(project, updated)
	}

	// get today total production

	match = tk.M{}
	if project != "" {
		match.Set("project", project)
	}

	if len(turbine) > 0 {
		match.Set("turbine", tk.M{}.Set("$in", turbine))
	}

	dateid, _ := time.Parse("20060102_0504", time.Now().Format("20060102_")+"_0000")
	match.Set("dateinfo.dateid", dateid.UTC())

	group = tk.M{
		"_id":        "$project",
		"production": tk.M{"$sum": "$production"},
	}

	pipes = []tk.M{}
	pipes = append(pipes, tk.M{}.Set("$match", match))
	pipes = append(pipes, tk.M{}.Set("$group", group))

	csr, e = DB().Connection.NewQuery().
		From(new(Monitoring).TableName()).
		Command("pipe", pipes).
		Cursor(nil)

	defer csr.Close()

	if e != nil {
		return helper.CreateResult(false, nil, "Error query : "+e.Error())
	}

	results = make([]tk.M, 0)
	e = csr.Fetch(&results, 0, false)

	if e != nil {
		return helper.CreateResult(false, nil, "Error query : "+e.Error())
	}

	projectProdResult := map[string]float64{}

	for _, v := range results {
		id := v.GetString("_id")
		prod := v.GetFloat64("production")
		projectProdResult[id] = prod
	}

	// combine the data

	res := []tk.M{}

	for proj, v := range projects {
		wsavg := 0.0
		turbineList := v.(tk.M).Get("turbines").([]tk.M)

		for _, i := range turbineList {
			wsavg += i.GetFloat64("windspeed")
		}

		wsavg = tk.Div(wsavg, tk.ToFloat64(len(turbineList), 0, tk.RoundingAuto))
		v.(tk.M).Set("totalwsavg", wsavg)
		// v.(tk.M).Set("totalprod", v.(tk.M).GetFloat64("totalprod")/1000)
		v.(tk.M).Set("totalprod", projectProdResult[proj]/1000)

		v.(tk.M).Set("project", proj)
		res = append(res, v.(tk.M))
	}

	// get latest date update from ScadaDataHFD

	availDate := k.Session("availdate", "")
	date := availDate.(*Availdatedata).ScadaDataHFD[1].UTC()

	finalResult := tk.M{}
	finalResult.Set("data", res)
	finalResult.Set("timestamp", tk.M{"minute": date.Format("15:04"), "date": date.Format("02 Jan 2006")})

	data := struct {
		Data tk.M
	}{
		Data: finalResult,
	}

	return helper.CreateResult(true, data, "success")
}

func (m *MonitoringController) GetEvent(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	p := tk.M{}
	e := k.GetPayload(&p)

	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	turbine := p.Get("turbine").([]interface{})
	project := ""
	if p.GetString("project") != "" {
		anProject := strings.Split(p.GetString("project"), "(")
		project = strings.TrimRight(anProject[0], " ")
	}

	// log.Printf("%#v \n", project)

	match := tk.M{}
	var projectList []interface{}
	if project != "" {
		match.Set("project", project)
		projectList = append(projectList, project)
	}

	if len(turbine) > 0 {
		match.Set("turbine", tk.M{}.Set("$in", turbine))
	}

	var pipes []tk.M
	pipes = append(pipes, tk.M{}.Set("$match", match))
	pipes = append(pipes, tk.M{}.Set("$sort", tk.M{
		"timestamp": 1,
		"turbine":   1,
		"project":   1,
	}))

	csr, e := DB().Connection.NewQuery().
		From(new(Monitoring).TableName()).
		Command("pipe", pipes).
		Cursor(nil)

	defer csr.Close()

	if e != nil {
		return helper.CreateResult(false, nil, "Error query : "+e.Error())
	}

	results := make([]MonitoringEvent, 0)
	e = csr.Fetch(&results, 0, false)

	if e != nil {
		return helper.CreateResult(false, nil, "Error query : "+e.Error())
	}

	data := struct {
		Data []MonitoringEvent
	}{
		Data: results,
	}

	return helper.CreateResult(true, data, "success")
}

func checkNAValue(val float64) (result float64) {
	if val == -9999999.0 {
		result = 0
	} else {
		result = val
	}

	return
}
