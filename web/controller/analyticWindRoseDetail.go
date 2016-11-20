package controller

import (
	. "eaciit/ostrowfm/library/core"
	// . "eaciit/ostrowfm/library/models"
	"eaciit/ostrowfm/web/helper"
	c "github.com/eaciit/crowd"
	"github.com/eaciit/dbox"
	"github.com/eaciit/knot/knot.v1"
	tk "github.com/eaciit/toolkit"
	// "time"
)

var arrDirection = [...]string{"N", "NE", "E", "SE", "S", "SW", "W", "NW"}
var arrCategory = [...]string{"0 to 4m/s", "4 to 7m/s", "7 to 9m/s", "9 to 14m/s", "14 and above"}

type AnalyticWindRoseDetailController struct {
	App
}

func CreateAnalyticWindRoseDetailController() *AnalyticWindRoseDetailController {
	var controller = new(AnalyticWindRoseDetailController)
	return controller
}

func (m *AnalyticWindRoseDetailController) GetWSData(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	var (
		// filter     		[]*dbox.Filter
		// pipes 			[]tk.M
		WindRoseResult []tk.M
	)

	p := new(PayloadAnalytic)
	e := k.GetPayload(&p)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	// tStart, _ := time.Parse("2006-01-02", p.DateStart.UTC().Format("2006-01-02"))
	// tEnd, _ := time.Parse("2006-01-02 15:04:05", p.DateEnd.UTC().Format("2006-01-02")+" 23:59:59")
	tStart, tEnd, e := helper.GetStartEndDate(k, p.Period, p.DateStart, p.DateEnd)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}
	turbine := p.Turbine
	project := p.Project

	coId := 0
	turbine = append(turbine, "MetTower")
	for _, turbineX := range turbine {
		tableName := "rpt_scadawindrosenew"
		groupdata := tk.M{}
		groupdata.Set("Index", coId)
		groupdata.Set("Name", turbineX.(string))
		coId++

		var filter []*dbox.Filter
		filter = append(filter, dbox.Ne("_id", ""))
		filter = append(filter, dbox.Gte("dateinfo.dateid", tStart))
		filter = append(filter, dbox.Lte("dateinfo.dateid", tEnd))

		if turbineX == "MetTower" {
			tableName = "rpt_scadawindrosemt"
			if project != "" {
				filter = append(filter, dbox.Eq("projectid", project))
			}
		} else {
			filter = append(filter, dbox.Eq("turbineid", turbineX))
		}

		fb := DB().Connection.Fb()
		fb.AddFilter(dbox.And(filter...))
		matches, e := fb.Build()
		if e != nil {
			helper.CreateResult(false, nil, e.Error())
		}
		groups := tk.M{}
		groupIds := tk.M{}
		group := []string{
			"directionno",
			"directiondesc",
			"wscategoryno",
			"wscategorydesc",
		}
		for _, val := range group {
			alias := val
			field := tk.Sprintf("$windroseitems.%s", val)
			groupIds[alias] = field
		}
		groups["_id"] = groupIds

		fields := []string{
			"hours",
			"contribute",
			"frequency",
		}

		for _, other := range fields {
			field := tk.Sprintf("$windroseitems.%s", other)
			op := ""
			if other == "contribute" {
				op = "$avg"
			} else {
				op = "$sum"
			}
			groups[other] = tk.M{op: field}
			//groups[other] = tk.M{"$sum": field}
		}

		pipes := []tk.M{{"$unwind": "$windroseitems"}, {"$match": matches}, {"$group": groups}}

		csr, e := DB().Connection.NewQuery().
			From(tableName). //From(new(WindRoseModel).TableName()).
			Command("pipe", pipes).
			Cursor(nil)

		if e != nil {
			helper.CreateResult(false, nil, e.Error())
		}
		defer csr.Close()

		resultItem := tk.Ms{}
		e = csr.Fetch(&resultItem, 0, false)
		if e != nil {
			return helper.CreateResult(false, nil, e.Error())
		}

		dummyRes := tk.Ms{}

		//Show All direction when data not avaible
		cWind := 0
		for _, ac := range arrCategory {
			cDir := 0
			for _, dir := range arrDirection {
				resdir := tk.M{}
				resdir.Set("directiondesc", dir)
				resdir.Set("directionno", cDir)
				resdir.Set("wscategorydesc", ac)
				resdir.Set("wscategoryno", cWind)

				exist := c.From(&resultItem).Where(func(x interface{}) interface{} {
					ids := x.(tk.M)["_id"].(tk.M)
					if ids["wscategoryno"] == cWind && ids["directionno"] == cDir {
						return true
					} else {
						return false
					}
				}).Exec().Result.Data().([]tk.M)

				// tk.Printf("exist %v\n", len(exist))
				if len(exist) <= 0 {
					resdir.Set("hours", 0.0)
					resdir.Set("contribute", 0.0)
					resdir.Set("frequency", 0.0)
					dummyRes = append(dummyRes, resdir)
				} else {
					resdir.Set("hours", exist[0]["hours"])
					resdir.Set("contribute", exist[0]["contribute"])
					resdir.Set("frequency", exist[0]["frequency"])
					dummyRes = append(dummyRes, resdir)
				}

				cDir++
			}
			cWind++
		}

		groupdata.Set("Data", dummyRes)
		WindRoseResult = append(WindRoseResult, groupdata)
	}
	//===================================================================================================================================================================

	data := struct {
		WindRose tk.Ms
	}{
		WindRose: WindRoseResult,
	}

	return helper.CreateResult(true, data, "success")
}
