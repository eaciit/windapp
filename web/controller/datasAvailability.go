package controller

import (
	// . "eaciit/wfdemo-git/library/core"
	// . "eaciit/wfdemo-git/library/models"
	"eaciit/wfdemo-git/web/helper"

	"github.com/eaciit/knot/knot.v1"
	tk "github.com/eaciit/toolkit"
)

type DatasAvailabilityController struct {
	App
}

// DatasAvailabilityController
func CreateDatasAvailabilityController() *DatasAvailabilityController {
	var controller = new(DatasAvailabilityController)
	return controller
}

// GetDataAvailability
func (m *DatasAvailabilityController) GetDataAvailability(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	/*var (
		pipes             []tk.M
		kpiAnalysisResult []tk.M
		list              []tk.M
	)

	p := new(PayloadKPI)
	e := k.GetPayload(&p)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	tStart, tEnd, e := helper.GetStartEndDate(k, p.Period, p.DateStart, p.DateEnd)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	match := tk.M{}
	groupId := tk.M{}

	match.Set("dateinfo.dateid", tk.M{"$gte": tStart, "$lte": tEnd})

	group := tk.M{
		"power":           tk.M{"$sum": "$power"},
		"energy":          tk.M{"$sum": "$energy"},
		"machinedowntime": tk.M{"$sum": "$machinedowntime"},
		"griddowntime":    tk.M{"$sum": "$griddowntime"},
		"oktime":          tk.M{"$sum": "$oktime"},
		"powerlost":       tk.M{"$sum": "$powerlost"},
		"totaltimestamp":  tk.M{"$sum": 1},
		"available":       tk.M{"$sum": "$available"},
		"minutes":         tk.M{"$sum": "$minutes"},
		"maxdate":         tk.M{"$max": "$dateinfo.dateid"},
		"mindate":         tk.M{"$min": "$dateinfo.dateid"},
	}

	group.Set("_id", groupId)

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

	// add by ams, 2016-10-07
	defer csr.Close()*/

	data := struct {
		Data []tk.M
	}{
		Data: nil,
	}

	return helper.CreateResult(true, data, "success")
}
