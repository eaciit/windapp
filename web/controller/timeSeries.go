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

type TimeSeriesController struct {
	App
}

func CreateTimeSeriesController() *TimeSeriesController {
	var controller = new(TimeSeriesController)
	return controller
}

func (m *TimeSeriesController) GetData(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	var (
		pipes []tk.M
		list  []tk.M
	)

	p := new(PayloadAnalytic)
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
	match.Set("avgwindspeed", tk.M{"$lte": 25})

	if p.Project != "" {
		anProject := strings.Split(p.Project, "(")
		match.Set("projectname", strings.TrimRight(anProject[0], " "))
	}

	group := tk.M{
		"energy":    tk.M{"$sum": "$energy"},
		"windspeed": tk.M{"$avg": "$avgwindspeed"},
	}

	group.Set("_id", "$timestamp")

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

	var dtProd, dtWS [][]interface{}

	for _, val := range list {
		intTimestamp := val.Get("_id").(time.Time).UTC().Unix()

		energy := val.GetFloat64("energy") / 1000
		wind := val.GetFloat64("windspeed")

		dtProd = append(dtProd, []interface{}{intTimestamp, energy})
		dtWS = append(dtWS, []interface{}{intTimestamp, wind})
	}

	result := []tk.M{}

	result = append(result, tk.M{"name": "Production", "data": dtProd, "unit": "MWh"})
	result = append(result, tk.M{"name": "Windspeed", "data": dtWS, "unit": "m/s"})

	data := struct {
		Data []tk.M
	}{
		Data: result,
	}

	return helper.CreateResult(true, data, "success")
}

func (m *TimeSeriesController) GetAvailDate(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	return helper.CreateResult(true, k.Session("availdate", ""), "success")
}
