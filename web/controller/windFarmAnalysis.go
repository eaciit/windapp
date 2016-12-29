package controller

import (
	. "eaciit/wfdemo-git/library/core"
	. "eaciit/wfdemo-git/library/models"
	"eaciit/wfdemo-git/web/helper"
	// "sort"
	// "strings"
	// "time"

	"github.com/eaciit/dbox"
	"github.com/eaciit/knot/knot.v1"
	tk "github.com/eaciit/toolkit"
)

type WindFarmAnalysisController struct {
	App
}

func CreateWindFarmAnalysisController() *WindFarmAnalysisController {
	c := new(WindFarmAnalysisController)
	return c
}

func (c *WindFarmAnalysisController) GetDataByProject(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	p := struct {
		Project string
	}{}

	e := k.GetPayload(&p)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	csr, e := DB().Connection.NewQuery().From(new(GWFAnalysisByProject).TableName()).
		Where(dbox.And(dbox.Eq("projectname", p.Project))).Order("OrderNo").Cursor(nil)

	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	datas := make([]GWFAnalysisByProject, 0)
	e = csr.Fetch(&datas, 0, false)
	csr.Close()

	return helper.CreateResult(true, datas, "success")
}

func (c *WindFarmAnalysisController) GetDataByTurbine1(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	p := struct {
		Project  string
		Turbines []interface{}
	}{}

	e := k.GetPayload(&p)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	tk.Println("Turbines: ", p.Turbines)

	tQry := make([]*dbox.Filter, 0)
	tQry = append(tQry, dbox.Eq("projectname", p.Project))
	if len(p.Turbines) > 0 {
		tQry = append(tQry, dbox.In("turbine", p.Turbines...))
	}

	tk.Println(tQry)

	csr, e := DB().Connection.NewQuery().From(new(GWFAnalysisByTurbine1).TableName()).
		Where(dbox.And(tQry...)).Order([]string{"turbine", "orderno"}...).Cursor(nil)

	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	datas := make([]GWFAnalysisByTurbine1, 0)
	e = csr.Fetch(&datas, 0, false)
	csr.Close()

	return helper.CreateResult(true, datas, "success")
}

func (c *WindFarmAnalysisController) GetDataByTurbine2(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	p := struct {
		Project  string
		Turbines []string
	}{}

	e := k.GetPayload(&p)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	tQry := make([]*dbox.Filter, 0)
	tQry = append(tQry, dbox.Eq("projectname", p.Project))
	if len(p.Turbines) > 0 {
		tQry = append(tQry, dbox.In("turbine", p.Turbines))
	}

	csr, e := DB().Connection.NewQuery().From(new(GWFAnalysisByTurbine2).TableName()).
		Where(dbox.And(tQry...)).Order("orderno").Cursor(nil)

	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	datas := make([]GWFAnalysisByTurbine2, 0)
	e = csr.Fetch(&datas, 0, false)
	csr.Close()

	return helper.CreateResult(true, datas, "success")
}
