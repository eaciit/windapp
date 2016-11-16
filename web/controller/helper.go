package controller

import (
	. "eaciit/wfdemo/library/core"
	"eaciit/wfdemo/web/helper"
	"sort"

	"github.com/eaciit/knot/knot.v1"
	"github.com/eaciit/toolkit"

	"github.com/eaciit/dbox"
)

type HelperController struct {
	App
}

type Sorting struct {
	Field string
	Dir   string
}

func CreateHelperController() *HelperController {
	var controller = new(HelperController)
	return controller
}

func (m *HelperController) GetTurbineList(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	csr, e := DB().Connection.NewQuery().From("ref_turbine").Cursor(nil)

	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}
	defer csr.Close()

	data := []toolkit.M{}
	e = csr.Fetch(&data, 0, false)

	result := []string{}

	for _, val := range data {
		result = append(result, val.GetString("turbineid"))
	}
	sort.Strings(result)

	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	return helper.CreateResult(true, result, "success")
}

func (m *HelperController) GetProjectInfo(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	var filter []*dbox.Filter

	type ProjectInfoFilter struct {
		Project  string
		Turbines []interface{}
	}

	p := new(ProjectInfoFilter)
	e := k.GetPayload(&p)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	if p.Project != "" {
		filter = append(filter, dbox.Eq("project", p.Project))
	}

	if len(p.Turbines) > 0 {
		filter = append(filter, dbox.In("turbineid", p.Turbines...))
	}

	csr, e := DB().Connection.NewQuery().From("ref_turbine").
		Where(filter...).
		Order("turbineid").
		Cursor(nil)

	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}
	defer csr.Close()

	data := []toolkit.M{}
	e = csr.Fetch(&data, 0, false)

	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	var totalCapacity float64

	for _, turbine := range data {
		totalCapacity += turbine.GetFloat64("capacitymw")
	}

	result := struct {
		Turbines      []toolkit.M
		TotalTurbine  int
		TotalCapacity float64
	}{
		Turbines:      data,
		TotalTurbine:  len(data),
		TotalCapacity: toolkit.ToFloat64(totalCapacity, 2, toolkit.RoundingAuto),
	}

	return helper.CreateResult(true, result, "success")
}

func (m *HelperController) GetProjectList(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	csr, e := DB().Connection.NewQuery().From("ref_project").Cursor(nil)

	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}
	defer csr.Close()

	data := []toolkit.M{}
	e = csr.Fetch(&data, 0, false)

	result := []string{}

	for _, val := range data {
		if val.GetString("projectid") == "Tejuva" {
			result = append(result, val.GetString("projectid"))
		}
	}
	sort.Strings(result)

	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	return helper.CreateResult(true, result, "success")
}

func (m *HelperController) GetModelList(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson
	result := []string{}

	return helper.CreateResult(true, result, "success")
}
