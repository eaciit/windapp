package controller

import (
	. "eaciit/wfdemo-git/library/core"
	"eaciit/wfdemo-git/web/helper"

	"github.com/eaciit/knot/knot.v1"
	"github.com/eaciit/toolkit"

	. "eaciit/wfdemo-git/library/models"

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
	var projects []interface{}
	resProj, e := helper.GetProjectList()

	for _, v := range resProj {
		projects = append(projects, v.Value)
	}

	result, e := helper.GetTurbineList(projects)

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

	filter = append(filter, dbox.Eq("active", true))

	if p.Project != "" {
		filter = append(filter, dbox.Eq("project", p.Project))
	}

	if len(p.Turbines) > 0 {
		filter = append(filter, dbox.In("turbineid", p.Turbines...))
	}

	csr, e := DB().Connection.NewQuery().From(new(TurbineMaster).TableName()).
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

	result, e := helper.GetProjectList()

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

func getAvailDateByCondition(project, dtype string) toolkit.M {
	latestDataPeriods := make([]LatestDataPeriod, 0)
	iQuery := DB().Connection.NewQuery().From(new(LatestDataPeriod).TableName())
	_acond := make([]*dbox.Filter, 0)

	if project != "" {
		_acond = append(_acond, dbox.Eq("projectname", project))
	}

	if dtype != "" {
		_acond = append(_acond, dbox.Eq("type", dtype))
	}

	if len(_acond) == 1 {
		iQuery.Where(_acond[0])
	} else if len(_acond) > 1 {
		iQuery.Where(dbox.And(_acond...))
	}

	csr, e := iQuery.Cursor(nil)
	if e != nil {
		return nil
	}

	e = csr.Fetch(&latestDataPeriods, 0, false)
	csr.Close()

	result := toolkit.M{}
	for _, val := range latestDataPeriods {
		for i, tval := range val.Data {
			val.Data[i] = tval.UTC()
		}
		if result.Has(val.ProjectName) {
			currData, _ := toolkit.ToM(result[val.ProjectName])
			currData.Set(val.Type, val.Data)
			result.Set(val.ProjectName, currData)
		} else {
			result.Set(val.ProjectName, toolkit.M{val.Type: val.Data})
		}
	}

	return result
}
