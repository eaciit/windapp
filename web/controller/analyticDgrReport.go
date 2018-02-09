package controller

import (
	. "eaciit/wfdemo-git/library/core"
	. "eaciit/wfdemo-git/library/models"
	"eaciit/wfdemo-git/web/helper"
	"strings"

	"github.com/eaciit/knot/knot.v1"

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
		DateStart time.Time
		DateEnd   time.Time
		Period    string
		Project   string
		Turbine   []interface{}
		Skip      int
		Take      int
		Sort      []Sorting
	}{}

	e := k.GetPayload(&p)

	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	turbine := p.Turbine
	project := p.Project

	turbineName, e := helper.GetTurbineNameList(project)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	query := DB().Connection.NewQuery().From("rpt_scadasummarydaily")
	tStart, tEnd, e := helper.GetStartEndDate(k, p.Period, p.DateStart, p.DateEnd)

	_dfilter := make([]*dbox.Filter, 0)
	_dfilter = append(_dfilter, dbox.Gte("dateinfo.dateid", tStart))
	_dfilter = append(_dfilter, dbox.Lte("dateinfo.dateid", tEnd))

	if project != "" {
		_dfilter = append(_dfilter, dbox.Eq("projectname", project))
	}

	if len(turbine) > 0 {
		_dfilter = append(_dfilter, dbox.In("turbine", turbine...))
	}

	query.Where(dbox.And(_dfilter...))

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
	for {
		result := DGRScadaModel{}
		e = csr.Fetch(&result, 1, false)
		if e != nil {
			break
		}
		result.DowntimeHours = (result.GridDownHours + result.MachineDownHours + result.Otherdowntimehours) * 3600
		result.TurbineName = turbineName[result.Turbine]

		results = append(results, result)
	}

	data := struct {
		Data  []DGRScadaModel
		Total int
	}{
		Data:  results,
		Total: csr.Count(),
	}

	return helper.CreateResult(true, data, "success")
}
