package controller

import (
	. "eaciit/wfdemo-git/library/core"
	. "eaciit/wfdemo-git/library/helper"
	. "eaciit/wfdemo-git/library/models"
	"eaciit/wfdemo-git/web/helper"
	"strings"

	"time"

	"github.com/eaciit/knot/knot.v1"
	tk "github.com/eaciit/toolkit"
)

var (
	from time.Time
	to   time.Time
)

type DataAvailabilityController struct {
	App
}

// DataAvailabilityController
func CreateDataAvailabilityController() *DataAvailabilityController {
	var controller = new(DataAvailabilityController)
	return controller
}

// GetDataAvailability
func (m *DataAvailabilityController) GetDataAvailability(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	var (
		result []tk.M
		months []string
	)

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

	result = append(result, getAvailCollection(project, turbine, "SCADA_DATA_OEM"))

	for {
		months = append(months, from.Format("Jan"))
		if from.Format("0601") == to.Format("0601") {
			break
		}
		from = GetNormalAddDateMonth(from.UTC(), 1)
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

func getAvailCollection(project string, turbines []interface{}, collType string) tk.M {
	var (
		pipes []tk.M
		list  []tk.M
	)
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
				"id":       "$details.id",
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

	pipes = append(pipes, tk.M{"$match": tk.M{"type": tk.M{"$eq": collType}}})
	pipes = append(pipes, tk.M{"$unwind": "$details"})
	pipes = append(pipes, tk.M{"$group": group})
	pipes = append(pipes, tk.M{"$project": projection})

	match := tk.M{}

	if project != "" {
		match.Set("project", project)
	}

	if len(turbines) > 0 {
		match.Set("turbine", tk.M{"$in": turbines})
	}

	if match.Get("turbine") != nil || match.Get("project") != nil {
		pipes = append(pipes, tk.M{"$match": match})
	}

	pipes = append(pipes, tk.M{"$sort": tk.M{"turbine": 1}})

	csr, e := DB().Connection.NewQuery().
		From(new(DataAvailability).TableName()).
		Command("pipe", pipes).
		Cursor(nil)

	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	e = csr.Fetch(&list, 0, false)

	defer csr.Close()

	res := []tk.M{}
	name := ""

	from = time.Now()
	to = time.Now()

	if len(list) > 0 {
		datas := []tk.M{}

		for _, oem := range list {
			t := oem.GetString("turbine")
			p := oem.GetString("project")
			_ = p
			pTo := oem.Get("periodTo").(time.Time)
			pFrom := oem.Get("periodFrom").(time.Time)

			from = pFrom.UTC()
			to = pTo.UTC()

			durationDays := pTo.UTC().Sub(pFrom.UTC()).Hours() / 24

			name = oem.GetString("name")
			availList := oem.Get("list").([]interface{})

			turbineDetails := []tk.M{}

			// log.Printf(">> %v | %v | %v | %v | %v | %v \n", p, t, pFrom.String(), pTo.String(), durationDays, name)

			for index := 1; index <= len(availList); index++ {
			breakAvail:
				for _, av := range availList {
					avail := av.(tk.M)
					if index == avail.GetInt("id") {

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

						// log.Printf(">>>> %v | %v | %v \n", start.Format("2 Jan 2006")+" until "+end.Format("2 Jan 2006"), class, tk.ToString(percentage)+"%")

						break breakAvail
					}
				}
			}

			turbine := tk.M{"TurbineName": t}
			turbine.Set("details", turbineDetails)

			res = append(res, turbine)
		}

		// dummy
		datas = append(datas, tk.M{
			"tooltip": "xxx until xxx",
			"class":   "progress-bar progress-bar-success",
			"value":   "100%",
		})

		return tk.M{"Category": name, "Turbine": res, "Data": datas}
	}

	return nil
}
