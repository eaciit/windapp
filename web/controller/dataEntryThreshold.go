package controller

import (
	. "eaciit/wfdemo-git/library/core"
	. "eaciit/wfdemo-git/library/models"
	"eaciit/wfdemo-git/web/helper"

	"github.com/eaciit/knot/knot.v1"
	// "github.com/eaciit/toolkit"
	"strings"

	_ "github.com/eaciit/dbox/dbc/mongo"
	"gopkg.in/mgo.v2/bson"
)

type DataEntryThresholdController struct {
	App
}

func CreateDataEntryThresholdController() *DataEntryThresholdController {
	var controller = new(DataEntryThresholdController)
	return controller
}

func (m *DataEntryThresholdController) GetList(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	p := struct {
		Project string
		Turbine []interface{}
		Skip    int
		Take    int
		Sort    []Sorting
	}{}

	e := k.GetPayload(&p)

	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	// turbine := p.Turbine
	// project := p.Project

	query := DBRealtime().NewQuery().From(new(StrangethresholdMaster).TableName()).Skip(p.Skip).Take(p.Take)

	// if project != "" {
	// 	query = query.Where(dbox.Eq("project", project))
	// }

	// if len(turbine) > 0 {
	// 	query = query.Where(dbox.In("turbineid", turbine...))
	// }

	if len(p.Sort) > 0 {
		var arrsort []string
		for _, val := range p.Sort {
			if val.Dir == "desc" {
				arrsort = append(arrsort, strings.ToLower("-"+val.Field))
			} else {
				arrsort = append(arrsort, strings.ToLower(val.Field))
			}
		}
		query = query.Order(arrsort...)
	}

	csr, e := query.Cursor(nil)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}
	defer csr.Close()

	results := make([]StrangethresholdMaster, 0)
	e = csr.Fetch(&results, 0, false)

	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}
	results2 := results

	data := struct {
		Data  []StrangethresholdMaster
		Total int
	}{
		Data:  results2,
		Total: csr.Count(),
	}

	return helper.CreateResult(true, data, "success")
}

func (m *DataEntryThresholdController) GetData(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	p := struct {
		Id bson.ObjectId
	}{}
	e := k.GetPayload(&p)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	result := new(TurbineModel)
	e = DB().GetById(result, p.Id)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	return helper.CreateResult(true, result, "success")
}

func (m *DataEntryThresholdController) Delete(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	p := struct {
		Id bson.ObjectId
	}{}
	e := k.GetPayload(&p)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	result := new(TurbineModel)
	e = DB().GetById(result, p.Id)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	e = DB().Delete(result)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	return helper.CreateResult(true, nil, "success")
}

func (m *DataEntryThresholdController) Save(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	p := struct {
		Id          string
		TurbineId   string
		TurbineName string
		Feeder      string
		Project     string
		Latitude    float64
		Longitude   float64
		Elevation   float64
		Capacitymw  float64
	}{}
	e := k.GetPayload(&p)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	mdl := new(TurbineModel)

	mdl.TurbineId = p.TurbineId
	mdl.TurbineName = p.TurbineName
	mdl.Feeder = p.Feeder
	mdl.Project = p.Project
	mdl.Latitude = p.Latitude
	mdl.Longitude = p.Longitude
	mdl.Elevation = p.Elevation
	mdl.Capacitymw = p.Capacitymw

	if p.Id != "" {
		mdl.Id = p.Id
	} else {
		mdl = mdl.New()
	}

	e = DB().Save(mdl)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	return helper.CreateResult(true, nil, "success")
}
