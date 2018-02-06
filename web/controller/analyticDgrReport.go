package controller

import (
	. "eaciit/wfdemo-git/library/core"
	. "eaciit/wfdemo-git/library/models"
	"eaciit/wfdemo-git/web/helper"

	"github.com/eaciit/knot/knot.v1"
	// "github.com/eaciit/toolkit"
	"strings"
	"time"

	"github.com/eaciit/dbox"
	_ "github.com/eaciit/dbox/dbc/mongo"
	// "gopkg.in/mgo.v2/bson"
)

type AnalyticDgrReportController struct {
	App
}

func CreateAnalyticDgrReportController() *AnalyticDgrReportController {
	var controller = new(AnalyticDgrReportController)
	return controller
}

func (m *AnalyticDgrReportController) GetList(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	p := struct {
		DateStart   time.Time
		DateEnd     time.Time
		Period      string
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

	turbine := p.Turbine
	project := p.Project

	query := DB().Connection.NewQuery().From(new(DGRScadaModel).TableName()).Skip(p.Skip).Take(p.Take)

	tStart, tEnd, e := helper.GetStartEndDate(k, p.Period, p.DateStart, p.DateEnd)
	query = query.Where(dbox.Gte("dateinfo.dateid", tStart))
	query = query.Where(dbox.Lte("dateinfo.dateid", tEnd))

	if project != "" {
		query = query.Where(dbox.Eq("projectname", project))
	}

	if len(turbine) > 0 {
		query = query.Where(dbox.In("turbine", turbine...))
	}


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

	results := make([]DGRScadaModel, 0)
	e = csr.Fetch(&results, 0, false)

	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}
	results2 := results

	data := struct {
		Data  []DGRScadaModel
		Total int
	}{
		Data:  results2,
		Total: csr.Count(),
	}

	return helper.CreateResult(true, data, "success")
}

