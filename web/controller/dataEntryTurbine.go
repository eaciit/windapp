package controller

import (
	. "eaciit/wfdemo-git/library/core"
	. "eaciit/wfdemo-git/library/models"
	"eaciit/wfdemo-git/web/helper"
	"github.com/eaciit/knot/knot.v1"
	"github.com/eaciit/toolkit"
	"gopkg.in/mgo.v2/bson"
	"strings"
)

type DataEntryTurbineController struct {
	App
}

func CreateDataEntryTurbineController() *DataEntryTurbineController {
	var controller = new(DataEntryTurbineController)
	return controller
}

func (m *DataEntryTurbineController) GetList(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	p := struct {
		Skip int
		Take int
		Sort []Sorting
	}{}
	e := k.GetPayload(&p)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	csr, e := DB().Find(new(TurbineModel), toolkit.M{}.Set("skip", p.Skip).Set("limit", p.Take))
	defer csr.Close()

	if len(p.Sort) > 0 {
		var arrsort []string
		for _, val := range p.Sort {
			if val.Dir == "desc" {
				arrsort = append(arrsort, strings.ToLower("-"+p.Sort[0].Field))
			} else {
				arrsort = append(arrsort, strings.ToLower(p.Sort[0].Field))
			}
		}
		csr, e = DB().Find(new(TurbineModel), toolkit.M{}.Set("order", arrsort).Set("skip", p.Skip).Set("limit", p.Take))
		if e != nil {
			return helper.CreateResult(false, nil, e.Error())
		}
		defer csr.Close()
	}

	results := make([]TurbineModel, 0)
	e = csr.Fetch(&results, 0, false)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}
	results2 := results

	csr, e = DB().Find(new(TurbineModel), toolkit.M{})
	defer csr.Close()

	results = make([]TurbineModel, 0)
	e = csr.Fetch(&results, 0, false)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	data := struct {
		Data  []TurbineModel
		Total int
	}{
		Data:  results2,
		Total: len(results),
	}

	return helper.CreateResult(true, data, "success")
}

func (m *DataEntryTurbineController) GetData(k *knot.WebContext) interface{} {
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

func (m *DataEntryTurbineController) Delete(k *knot.WebContext) interface{} {
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

func (m *DataEntryTurbineController) Save(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	p := struct {
		Id        		bson.ObjectId
		TurbineId       string
		TurbineName     string
		Feeder      	string
		Project         string
		Latitude		float64
		Longitude		float64
		Elevation		float64
		Capacitymw		float64
	}{}
	e := k.GetPayload(&p)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	mdl := new(TurbineModel)
	if p.Id == "" {
		mdl.Id = bson.NewObjectId()
	} else {
		mdl.Id = p.Id
	}

	mdl.TurbineId 		= p.TurbineId
	mdl.TurbineName		= p.TurbineName
	mdl.Feeder 			= p.Feeder
	mdl.Project			= p.Project
	mdl.Latitude		= p.Latitude
	mdl.Longitude		= p.Longitude
	mdl.Elevation		= p.Elevation
	mdl.Capacitymw		= p.Capacitymw

	e = DB().Save(mdl)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	return helper.CreateResult(true, nil, "success")
}
