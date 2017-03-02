package controller

import (
	. "eaciit/wfdemo-git/library/core"
	. "eaciit/wfdemo-git/library/models"
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

	var (
		pipes []tk.M
		list  []tk.M
	)

	p := new(tk.M)
	e := k.GetPayload(&p)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	group := tk.M{
		"_id": tk.M{
			"name":    "$name",
			"project": "$details.projectname",
			"turbine": "$details.turbine",
		},
		"timestamp": tk.M{"$max": "$timestamp"},
		"list": tk.M{
			"$push": tk.M{
				"start":    "$details.start",
				"end":      "$details.end",
				"duration": "$details.duration",
				"isavail":  "$details.isavail",
			},
		},
	}

	projection := tk.M{
		"name":      "$_id.name",
		"project":   "$_id.projectname",
		"turbine":   "$_id.turbine",
		"timestamp": 1,
		"list":      1,
	}

	// pipes = append(pipes, tk.M{"$match": match})
	pipes = append(pipes, tk.M{"$unwind": "$details"})
	pipes = append(pipes, tk.M{"$group": group})
	pipes = append(pipes, tk.M{"$project": projection})

	// pipes = append(pipes, tk.M{"$sort": tk.M{"_id": 1}})

	csr, e := DB().Connection.NewQuery().
		From(new(DataAvailability).TableName()).
		Command("pipe", pipes).
		Cursor(nil)

	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	e = csr.Fetch(&list, 0, false)

	defer csr.Close()

	data := struct {
		Data  []tk.M
		Month []string
	}{
		Data:  nil,
		Month: nil,
	}

	return helper.CreateResult(true, data, "success")
}
