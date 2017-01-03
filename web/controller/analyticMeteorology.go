package controller

import (
	. "eaciit/wfdemo-git/library/core"
	. "eaciit/wfdemo-git/library/models"
	"eaciit/wfdemo-git/web/helper"
	"strings"

	"time"

	"github.com/eaciit/knot/knot.v1"
	tk "github.com/eaciit/toolkit"
)

type AnalyticMeteorologyController struct {
	App
}

func CreateAnalyticMeteorologyController() *AnalyticMeteorologyController {
	var controller = new(AnalyticMeteorologyController)
	return controller
}

func (c *AnalyticMeteorologyController) AverageWindSpeed(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	type PayloadAvgWindSpeed struct {
		Period          string
		Project         string
		Turbine         []interface{}
		DateStart       time.Time
		DateEnd         time.Time
		SeriesBreakDown string
		TimeBreakDown   string
	}

	var (
		pipes    []tk.M
		metTower []tk.M
		turbines []tk.M
		list     []tk.M
	)

	p := new(PayloadAvgWindSpeed)
	e := k.GetPayload(&p)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	tStart, tEnd, e := helper.GetStartEndDate(k, p.Period, p.DateStart, p.DateEnd)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	match := tk.M{}

	match.Set("dateinfo.dateid", tk.M{"$gte": tStart, "$lte": tEnd})

	if p.Project != "" {
		anProject := strings.Split(p.Project, "(")
		match.Set("projectname", strings.TrimRight(anProject[0], " "))
	}

	if len(p.Turbine) > 0 {
		match.Set("turbine", tk.M{"$in": p.Turbine})
	}

	group := tk.M{
		"windspeed": tk.M{"$avg": "$avgwindspeed"},
	}

	groupID := tk.M{}

	if p.TimeBreakDown == "daily" {
		groupID.Set("dateid", "$dateinfo.dateid")
	} else if p.TimeBreakDown == "monthly" {
		groupID.Set("monthdesc", "$dateinfo.monthdesc")
	}

	if p.SeriesBreakDown == "byturbine" {
		groupID.Set("turbine", "$turbine")
	}

	group.Set("_id", groupID)

	pipes = append(pipes, tk.M{"$match": match})
	pipes = append(pipes, tk.M{"$group": group})
	pipes = append(pipes, tk.M{"$sort": tk.M{"_id": 1}})

	csr, e := DB().Connection.NewQuery().
		From(new(ScadaData).TableName()).
		Command("pipe", pipes).
		Cursor(nil)

	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	e = csr.Fetch(&list, 0, false)

	csr.Close()

	for _, val := range list {
		id := val.Get("_id").(tk.M)
		turVal := tk.M{}

		if id.Get("dateid") == nil {
			turVal.Set("time", id.Get("dateid").(time.Time))
		} else {
			turVal.Set("time", id.GetString("monthdesc"))
		}

		wind := val.GetFloat64("windspeed")
		turVal.Set("value", wind)

		turbines = append(turbines, turVal)
	}

	list = []tk.M{}

	// met tower

	match = tk.M{}

	match.Set("dateinfo.dateid", tk.M{"$gte": tStart, "$lte": tEnd})

	group = tk.M{
		"windspeed": tk.M{"$avg": "$vhubws90mavg"},
	}

	groupID = tk.M{}

	if p.TimeBreakDown == "daily" {
		groupID.Set("dateid", "$dateinfo.dateid")
	} else if p.TimeBreakDown == "monthly" {
		groupID.Set("monthdesc", "$dateinfo.monthdesc")
	}

	group.Set("_id", groupID)

	pipes = []tk.M{}
	pipes = append(pipes, tk.M{"$match": match})
	pipes = append(pipes, tk.M{"$group": group})
	pipes = append(pipes, tk.M{"$sort": tk.M{"_id": 1}})

	csr, e = DB().Connection.NewQuery().
		From(new(MetTower).TableName()).
		Command("pipe", pipes).
		Cursor(nil)

	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	e = csr.Fetch(&list, 0, false)

	csr.Close()

	for _, val := range list {
		id := val.Get("_id").(tk.M)
		turVal := tk.M{}

		if id.Get("dateid") == nil {
			turVal.Set("time", id.Get("dateid").(time.Time))
		} else {
			turVal.Set("time", id.GetString("monthdesc"))
		}

		wind := val.GetFloat64("windspeed")
		turVal.Set("value", wind)

		metTower = append(turbines, turVal)
	}

	data := struct {
		Data tk.M
	}{
		Data: tk.M{"mettower": metTower, "turbine": turbines},
	}

	return helper.CreateResult(true, data, "success")
}
