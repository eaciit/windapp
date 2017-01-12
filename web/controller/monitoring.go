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

	p := new(PayloadAnalytic)
	e := k.GetPayload(&p)

	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	turbine := p.Turbine
	project := ""
	if p.Project != "" {
		anProject := strings.Split(p.Project, "(")
		project = strings.TrimRight(anProject[0], " ")
	}

	match := tk.M{}
	turbines := map[string]tk.M{}
	var projectList []interface{}
	if project != "" {
		match.Set("projectname", project)
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
		windspeed := v.GetFloat64("windspeed")
		production := v.GetFloat64("production")
		rotorspeedrpm := v.GetFloat64("rotorspeedrpm")
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

	res := []tk.M{}

	for proj, v := range projects {
		wsavg := 0.0
		turbineList := v.(tk.M).Get("turbines").([]tk.M)

		for _, i := range turbineList {
			wsavg += i.GetFloat64("windspeed")
		}

		wsavg = tk.Div(wsavg, tk.ToFloat64(len(turbineList), 0, tk.RoundingAuto))
		v.(tk.M).Set("totalwsavg", wsavg)
		v.(tk.M).Set("totalprod", v.(tk.M).GetFloat64("totalprod")/1000)

		v.(tk.M).Set("project", proj)
		res = append(res, v.(tk.M))
	}

	data := struct {
		Data []tk.M
	}{
		Data: res,
	}

	return helper.CreateResult(true, data, "success")
}
