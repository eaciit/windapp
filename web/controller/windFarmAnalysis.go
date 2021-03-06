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

	project := p.Project

	csr, e := DB().Connection.NewQuery().From(new(GWFAnalysisByProject).TableName()).
		Where(dbox.And(dbox.Eq("projectname", project))).Order("OrderNo").Cursor(nil)

	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	datas := make([]GWFAnalysisByProject, 0)
	e = csr.Fetch(&datas, 0, false)
	csr.Close()

	title := make([]string, 4)
	if len(datas) > 0 {
		d := datas[0]
		title[0] = "Rolling 12 Days<br /><span class='k-info'>" + d.Roll12Days.DateText + "</span>"
		title[1] = "Rolling 12 Weeks<br /><span class='k-info'>" + d.Roll12Weeks.DateText + "</span>"
		title[2] = "Rolling 12 Months<br /><span class='k-info'>" + d.Roll12Months.DateText + "</span>"
		title[3] = "Rolling 12 Quarters<br /><span class='k-info'>" + d.Roll12Quarters.DateText + "</span>"
	}

	datareturn := tk.M{}.Set("data", datas).Set("header", title)

	return helper.CreateResult(true, datareturn, "success")
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

	project := p.Project

	tQry := make([]*dbox.Filter, 0)
	tQry = append(tQry, dbox.Eq("projectname", project))
	if len(p.Turbines) > 0 {
		tQry = append(tQry, dbox.In("turbine", p.Turbines...))
	}

	csr, e := DB().Connection.NewQuery().From(new(GWFAnalysisByTurbine1).TableName()).
		Where(dbox.And(tQry...)).Order([]string{"turbine", "orderno"}...).Cursor(nil)

	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	datas := make([]GWFAnalysisByTurbine1, 0)
	e = csr.Fetch(&datas, 0, false)
	csr.Close()

	title := make([]string, 4)
	turbineName, e := helper.GetTurbineNameList(project)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}
	if len(datas) > 0 {
		d := datas[0]
		title[0] = "Rolling 12 Days<br /><span class='k-info'>" + d.Roll12Days.DateText + "</span>"
		title[1] = "Rolling 12 Weeks<br /><span class='k-info'>" + d.Roll12Weeks.DateText + "</span>"
		title[2] = "Rolling 12 Months<br /><span class='k-info'>" + d.Roll12Months.DateText + "</span>"
		title[3] = "Rolling 12 Quarters<br /><span class='k-info'>" + d.Roll12Quarters.DateText + "</span>"
		for idx, val := range datas {
			datas[idx].Turbine = turbineName[val.Turbine]
		}
	}

	datareturn := tk.M{}.Set("data", datas).Set("header", title)

	return helper.CreateResult(true, datareturn, "success")
}

func (c *WindFarmAnalysisController) GetDataByTurbine2(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	p := struct {
		Project  string
		Turbines []string
	}{}

	colors := []string{
		"#ED1C24", "#A3238E", "#00A65D", "#F58220", "#0066B3", "#5C2D91", "#FFF200", "#579835", "#CF3834", "#00B274", "#74489D",
		"#C06616", "#5565AF", "#CCBE00", "#390A5D", "#006D6F", "#65C295", "#F04E4D", "#407927", "#00599D", "#A09600", "#0D1F63",
		"#C38312", "#003D73", "#454FA1", "#BC312E",
	}

	e := k.GetPayload(&p)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	project := p.Project

	tQry := make([]*dbox.Filter, 0)
	tQry = append(tQry, dbox.Eq("projectname", project))
	if len(p.Turbines) > 0 {
		tQry = append(tQry, dbox.In("turbine", p.Turbines))
	} else {
		turbines := make([]TurbineMaster, 0)
		csrt, et := DB().Connection.NewQuery().From(new(TurbineMaster).TableName()).
			Where(dbox.And(dbox.Eq("project", project))).Order("turbineid").Cursor(nil)

		if et != nil {
			tk.Println(et.Error())
		}
		et = csrt.Fetch(&turbines, 0, false)
		csrt.Close()

		for _, t := range turbines {
			p.Turbines = append(p.Turbines, t.TurbineId)
		}
	}

	csr, e := DB().Connection.NewQuery().From(new(GWFAnalysisByTurbine2).TableName()).
		Where(dbox.And(tQry...)).Order("orderno").Cursor(nil)

	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	chartSeries := []tk.M{}
	chartSeries = append(chartSeries, tk.M{
		"field": "Average",
		"name":  "Average",
		"color": colors[0],
	})
	for idx, t := range p.Turbines {
		chartSeries = append(chartSeries, tk.M{
			"field": t,
			"name":  t,
			"color": colors[(idx + 1)],
		})
	}

	datas := make([]GWFAnalysisByTurbine2, 0)
	e = csr.Fetch(&datas, 0, false)
	csr.Close()

	retVal := tk.M{}.Set("ChartSeries", chartSeries).Set("ChartData", datas)

	return helper.CreateResult(true, retVal, "success")
}
