package controller

import (
	. "eaciit/wfdemo-git/library/core"
	. "eaciit/wfdemo-git/library/models"
	"github.com/eaciit/knot/knot.v1"
	"github.com/eaciit/toolkit"
	"gopkg.in/mgo.v2/bson"
	"strings"
)

type DataEntryPowerCurveController struct {
	App
}

func CreateDataEntryPowerCurveController() *DataEntryPowerCurveController {
	var controller = new(DataEntryPowerCurveController)
	return controller
}

func (m *DataEntryPowerCurveController) GetList(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	p := struct {
		Skip int
		Take int
		Sort []Sorting
	}{}
	e := k.GetPayload(&p)
	if e != nil {
		toolkit.Println(e.Error())
	}

	csr, e := DB().Find(new(PowerCurveModel), toolkit.M{}.Set("skip", p.Skip).Set("limit", p.Take))
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
		csr, e = DB().Find(new(PowerCurveModel), toolkit.M{}.Set("order", arrsort).Set("skip", p.Skip).Set("limit", p.Take))
		defer csr.Close()
	}

	results := make([]PowerCurveModel, 0)
	e = csr.Fetch(&results, 0, false)
	if e != nil {
		return e.Error()
	}
	results2 := results

	csr, e = DB().Find(new(PowerCurveModel), toolkit.M{})
	defer csr.Close()

	results = make([]PowerCurveModel, 0)
	e = csr.Fetch(&results, 0, false)
	if e != nil {
		return e.Error()
	}

	data := struct {
		Data  []PowerCurveModel
		Total int
	}{
		Data:  results2,
		Total: len(results),
	}

	return data
}

func (m *DataEntryPowerCurveController) GetData(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	p := struct {
		Id bson.ObjectId
	}{}
	e := k.GetPayload(&p)
	if e != nil {
		toolkit.Println(e.Error())
	}

	result := new(PowerCurveModel)
	e = DB().GetById(result, p.Id)
	if e != nil {
		toolkit.Println(e.Error())
	}

	return result
}

func (m *DataEntryPowerCurveController) Delete(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	p := struct {
		Id bson.ObjectId
	}{}
	e := k.GetPayload(&p)
	if e != nil {
		toolkit.Println(e.Error())
	}

	result := new(PowerCurveModel)
	e = DB().GetById(result, p.Id)
	if e != nil {
		toolkit.Println(e.Error())
	}

	e = DB().Delete(result)

	return e
}

func (m *DataEntryPowerCurveController) Save(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	p := struct {
		Id        bson.ObjectId
		Model     string
		WindSpeed float64
		Power1    float64
	}{}
	e := k.GetPayload(&p)
	if e != nil {
		toolkit.Println(e.Error())
	}

	mdl := new(PowerCurveModel)
	if p.Id == "" {
		mdl.ID = bson.NewObjectId()
	} else {
		mdl.ID = p.Id
	}
	mdl.Model = p.Model
	mdl.WindSpeed = p.WindSpeed
	mdl.Power1 = p.Power1

	e = DB().Save(mdl)

	return ""
}
