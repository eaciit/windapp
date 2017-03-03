package controller

import (
	. "eaciit/wfdemo-git/library/core"
	. "eaciit/wfdemo-git/library/helper"
	. "eaciit/wfdemo-git/library/models"
	"eaciit/wfdemo-git/web/helper"

	"time"

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
		pipes  []tk.M
		oems   []tk.M
		result []tk.M
	)

	p := new(tk.M)
	e := k.GetPayload(&p)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	// ScadaDataOEM

	group := tk.M{
		"_id": tk.M{
			"name":    "$name",
			"project": "$details.projectname",
			"turbine": "$details.turbine",
		},
		"periodTo":   tk.M{"$max": "$periodto"},
		"periodFrom": tk.M{"$min": "$periodfrom"},
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
		"name":       "$_id.name",
		"project":    "$_id.project",
		"turbine":    "$_id.turbine",
		"periodTo":   1,
		"periodFrom": 1,
		"list":       1,
	}

	// pipes = append(pipes, tk.M{"$match": match})
	pipes = append(pipes, tk.M{"$unwind": "$details"})
	pipes = append(pipes, tk.M{"$group": group})
	pipes = append(pipes, tk.M{"$sort": tk.M{"_id.turbine": 1}})
	pipes = append(pipes, tk.M{"$project": projection})

	// pipes = append(pipes, tk.M{"$sort": tk.M{"_id": 1}})

	csr, e := DB().Connection.NewQuery().
		From(new(DataAvailability).TableName()).
		Command("pipe", pipes).
		Cursor(nil)

	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	e = csr.Fetch(&oems, 0, false)

	defer csr.Close()

	resOEM := []tk.M{}
	name := ""

	from := time.Now()
	to := time.Now()

	for _, oem := range oems {
		t := oem.GetString("turbine")
		p := oem.GetString("project")
		_ = p
		pTo := oem.Get("periodTo").(time.Time)
		pFrom := oem.Get("periodFrom").(time.Time)

		durationDays := pFrom.UTC().Sub(pTo.UTC()).Hours() / 24

		name = oem.GetString("name")
		availList := oem.Get("list").([]tk.M)

		turbineDetails := []tk.M{}

		for _, avail := range availList {
			start := avail.Get("start").(time.Time).UTC()
			end := avail.Get("end").(time.Time).UTC()
			duration := avail.GetFloat64("duration")
			isAvail := avail.Get("isavail").(bool)
			class := "progress-bar progress-bar-success"

			if !isAvail {
				class = "progress-bar progress-bar-red"
			}

			percentage := duration / durationDays * 100

			turbineDetails = append(turbineDetails, tk.M{
				"tooltip": start.Format("2 Jan 2006") + " until " + end.Format("2 Jan 2006"),
				"class":   class,
				"value":   tk.ToString(percentage) + "%",
			})
		}

		turbine := tk.M{"TurbineName": t}
		turbine.Set("details", turbineDetails)
	}

	result = append(result, tk.M{"Category": name, "Turbine": resOEM})

	months := []string{}

	for {
		months = append(months, from.Format("Jan"))
		from = GetNormalAddDateMonth(to.UTC(), 1)

		if from.Format("060102") == to.Format("060102") {
			break
		}
	}

	data := struct {
		Data  []tk.M
		Month []string
	}{
		Data:  result,
		Month: months,
	}

	return helper.CreateResult(true, data, "success")
}
