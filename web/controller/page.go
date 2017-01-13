package controller

import (
	// . "eaciit/wfdemo-git/library/models"
	"eaciit/wfdemo-git/web/helper"
	"net/http"

	"github.com/eaciit/knot/knot.v1"
	"github.com/eaciit/toolkit"
)

type PageController struct {
	App
	Params toolkit.M
}

var (
	DefaultIncludes = []string{"_head.html", "_menu.html", "_loader.html", "_script_template.html"}
)

func CreatePageController(AppName string) *PageController {
	var controller = new(PageController)
	controller.Params = toolkit.M{"AppName": AppName}
	return controller
}

func (w *PageController) GetParams(r *knot.WebContext, isAnalyst bool) toolkit.M {
	w.Params.Set("AntiCache", toolkit.RandomString(20))
	w.Params.Set("CurrentDateData", helper.GetLastDateData(r))

	// w.Params.Set("Menus", r.Session("menus", []string{}))

	if isAnalyst {
		projectList, _ := helper.GetProjectList()
		turbineList, _ := helper.GetTurbineList("")

		w.Params.Set("ProjectList", projectList)
		w.Params.Set("TurbineList", turbineList)
	}

	r.Writer.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	r.Writer.Header().Set("Pragma", "no-cache")
	r.Writer.Header().Set("Expires", "0")
	// WriteLog(r.Session("sessionid", ""), "access", r.Request.URL.String())
	return w.Params
}

func (w *PageController) Index(r *knot.WebContext) interface{} {
	r.Config.OutputType = knot.OutputTemplate
	r.Config.LayoutTemplate = LayoutFile
	r.Config.ViewName = "page-index.html"
	return w.GetParams(r, false)
}

func (w *PageController) Login(r *knot.WebContext) interface{} {
	r.Config.OutputType = knot.OutputTemplate
	r.Config.ViewName = "page-login.html"

	if r.Session("sessionid", "") != "" {
		r.SetSession("sessionid", "")
	}

	return w.GetParams(r, false)
}

func (w *PageController) DataBrowser(r *knot.WebContext) interface{} {
	r.Config.OutputType = knot.OutputTemplate
	r.Config.LayoutTemplate = LayoutFile
	r.Config.ViewName = "page-databrowser.html"
	return w.GetParams(r, true).Set("ColumnList", GetCustomFieldList()).Set("HDFColList", GetHFDCustomFieldList())
}

func (w *PageController) DataBrowserNew(r *knot.WebContext) interface{} {
	r.Config.OutputType = knot.OutputTemplate
	r.Config.LayoutTemplate = LayoutFile
	r.Config.ViewName = "page-databrowser-new.html"

	return w.GetParams(r, false).Set("ColumnList", GetCustomFieldList()).Set("HDFColList", GetHFDCustomFieldList())
}

/*func (w *PageController) AnalyticWindDistribution(r *knot.WebContext) interface{} {
	r.Config.OutputType = knot.OutputTemplate
	r.Config.LayoutTemplate = LayoutFile
	r.Config.ViewName = "page-wind-distribution.html"
	return w.GetParams(r, true)
}*/

/*func (w *PageController) AnalyticWindAvailability(r *knot.WebContext) interface{} {
	r.Config.OutputType = knot.OutputTemplate
	r.Config.LayoutTemplate = LayoutFile
	r.Config.ViewName = "page-wind-availability-analysis.html"
	return w.GetParams(r, true)
}*/

/*func (w *PageController) AnalyticWindRose(r *knot.WebContext) interface{} {
	r.Config.OutputType = knot.OutputTemplate
	r.Config.LayoutTemplate = LayoutFile
	r.Config.ViewName = "page-analytic-wind-rose.html"
	return w.GetParams(r, true)
}*/

/*func (w *PageController) AnalyticWindRoseDetail(r *knot.WebContext) interface{} {
	r.Config.OutputType = knot.OutputTemplate
	r.Config.LayoutTemplate = LayoutFile
	r.Config.ViewName = "page-analytic-wind-rose-detail.html"
	return w.GetParams(r, true)
}*/

/*func (w *PageController) AnalyticWindRoseFlexi(r *knot.WebContext) interface{} {
	r.Config.OutputType = knot.OutputTemplate
	r.Config.LayoutTemplate = LayoutFile
	r.Config.ViewName = "page-analytic-wind-rose-flexi.html"
	return w.GetParams(r, true)
}*/

/*func (w *PageController) AnalyticWRFlexiDetail(r *knot.WebContext) interface{} {
	r.Config.OutputType = knot.OutputTemplate
	r.Config.LayoutTemplate = LayoutFile
	r.Config.ViewName = "page-analytic-wr-flexi-detail.html"
	return w.GetParams(r, true)
}*/

func (w *PageController) AnalyticPerformanceIndex(r *knot.WebContext) interface{} {
	r.Config.OutputType = knot.OutputTemplate
	r.Config.LayoutTemplate = LayoutFile
	r.Config.ViewName = "page-analytic-performance-index.html"

	return w.GetParams(r, true)
}

func (w *PageController) AnalyticPowerCurve(r *knot.WebContext) interface{} {
	r.Config.OutputType = knot.OutputTemplate
	r.Config.LayoutTemplate = LayoutFile
	r.Config.ViewName = "page-analytic-power-curve.html"

	return w.GetParams(r, true)
}

func (w *PageController) AnalyticPCMonthly(r *knot.WebContext) interface{} {
	r.Config.OutputType = knot.OutputTemplate
	r.Config.LayoutTemplate = LayoutFile
	r.Config.ViewName = "page-analytic-power-curve/individual-month.html"

	return w.GetParams(r, true)
}

func (w *PageController) AnalyticPCComparison(r *knot.WebContext) interface{} {
	r.Config.OutputType = knot.OutputTemplate
	r.Config.LayoutTemplate = LayoutFile
	r.Config.ViewName = "page-analytic-power-curve/comparison.html"

	return w.GetParams(r, true)
}

func (w *PageController) AnalyticPCScatter(r *knot.WebContext) interface{} {
	r.Config.OutputType = knot.OutputTemplate
	r.Config.LayoutTemplate = LayoutFile
	r.Config.ViewName = "page-analytic-power-curve/scatter.html"
	r.Config.IncludeFiles = append(DefaultIncludes, []string{"page-analytic-power-curve/page-filter-scatter.html"}...)

	return w.GetParams(r, true)
}

/*func (w *PageController) AnalyticDgrScada(r *knot.WebContext) interface{} {
	r.Config.OutputType = knot.OutputTemplate
	r.Config.LayoutTemplate = LayoutFile
	r.Config.ViewName = "page-analytic-dgr-scada.html"
	return w.GetParams(r, true)
}*/

func (w *PageController) AnalyticKeyMetrics(r *knot.WebContext) interface{} {
	r.Config.OutputType = knot.OutputTemplate
	r.Config.LayoutTemplate = LayoutFile
	r.Config.ViewName = "page-analytic-key-metrics.html"

	return w.GetParams(r, true)
}

func (w *PageController) AnalyticKpi(r *knot.WebContext) interface{} {
	r.Config.OutputType = knot.OutputTemplate
	r.Config.LayoutTemplate = LayoutFile
	r.Config.ViewName = "page-analytic-kpi.html"

	return w.GetParams(r, true)
}

/*func (w *PageController) AnalyticAvailability(r *knot.WebContext) interface{} {
	r.Config.OutputType = knot.OutputTemplate
	r.Config.LayoutTemplate = LayoutFile
	r.Config.ViewName = "page-availability-analysis.html"

	return w.GetParams(r, true)
}*/

func (w *PageController) AnalyticLoss(r *knot.WebContext) interface{} {
	r.Config.OutputType = knot.OutputTemplate
	r.Config.LayoutTemplate = LayoutFile
	r.Config.ViewName = "page-loss-analysis.html"

	return w.GetParams(r, true)
}

func (w *PageController) AnalyticDataConsistency(r *knot.WebContext) interface{} {
	r.Config.OutputType = knot.OutputTemplate
	r.Config.LayoutTemplate = LayoutFile
	r.Config.ViewName = "page-data-consistency.html"

	return w.GetParams(r, true)
}

func (w *PageController) AnalyticMeteorology(r *knot.WebContext) interface{} {
	r.Config.OutputType = knot.OutputTemplate
	r.Config.LayoutTemplate = LayoutFile
	r.Config.ViewName = "page-analytic-meteorology.html"

	return w.GetParams(r, true)
}

func (w *PageController) AnalyticComparison(r *knot.WebContext) interface{} {
	r.Config.OutputType = knot.OutputTemplate
	r.Config.LayoutTemplate = LayoutFile
	r.Config.ViewName = "page-analytic-comparison.html"

	return w.GetParams(r, true)
}
func (w *PageController) AnalyticDataHistogram(r *knot.WebContext) interface{} {
	r.Config.OutputType = knot.OutputTemplate
	r.Config.LayoutTemplate = LayoutFile
	r.Config.ViewName = "page-analytic-data-histogram.html"

	return w.GetParams(r, true)
}
func (w *PageController) DataEntryPowerCurve(r *knot.WebContext) interface{} {
	r.Config.OutputType = knot.OutputTemplate
	r.Config.LayoutTemplate = LayoutFile
	r.Config.ViewName = "page-dataentry-power-curve.html"

	return w.GetParams(r, false)
}
func (w *PageController) Access(r *knot.WebContext) interface{} {
	r.Config.OutputType = knot.OutputTemplate
	r.Config.LayoutTemplate = LayoutFile
	r.Config.ViewName = "page-access.html"

	return w.GetParams(r, false)
}

func (w *PageController) Group(r *knot.WebContext) interface{} {
	r.Config.OutputType = knot.OutputTemplate
	r.Config.LayoutTemplate = LayoutFile
	r.Config.ViewName = "page-group.html"

	return w.GetParams(r, false)
}

func (w *PageController) Session(r *knot.WebContext) interface{} {
	r.Config.OutputType = knot.OutputTemplate
	r.Config.LayoutTemplate = LayoutFile
	r.Config.ViewName = "page-session.html"

	return w.GetParams(r, false)
}

func (w *PageController) Log(r *knot.WebContext) interface{} {
	r.Config.OutputType = knot.OutputTemplate
	r.Config.LayoutTemplate = LayoutFile
	r.Config.ViewName = "page-log.html"

	return w.GetParams(r, false)
}

func (w *PageController) AdminTable(r *knot.WebContext) interface{} {
	r.Config.OutputType = knot.OutputTemplate
	r.Config.LayoutTemplate = LayoutFile
	r.Config.ViewName = "page-admintable.html"

	return w.GetParams(r, false)
}

func (w *PageController) User(r *knot.WebContext) interface{} {
	r.Config.OutputType = knot.OutputTemplate
	r.Config.LayoutTemplate = LayoutFile
	r.Config.ViewName = "page-user.html"

	return w.GetParams(r, false)
}

func (w *PageController) Monitoring(r *knot.WebContext) interface{} {
	r.Config.OutputType = knot.OutputTemplate
	r.Config.LayoutTemplate = LayoutFile
	r.Config.ViewName = "page-monitoring.html"
	r.Config.IncludeFiles = append(DefaultIncludes, []string{"_filter-monitoring.html"}...)
	return w.GetParams(r, true)
}

func (w *PageController) Dashboard(r *knot.WebContext) interface{} {
	r.Config.OutputType = knot.OutputTemplate
	r.Config.LayoutTemplate = LayoutFile
	r.Config.ViewName = "page-dashboard.html"
	r.Config.IncludeFiles = append(DefaultIncludes, []string{"page-dashboard-summary.html", "page-dashboard-production.html", "page-dashboard-availability.html"}...)
	return w.GetParams(r, false)
}

func (w *PageController) Home(r *knot.WebContext) interface{} {
	http.Redirect(r.Writer, r.Request, "dashboard", http.StatusTemporaryRedirect)
	return w.GetParams(r, false)
}

func (w *PageController) TurbineHealth(r *knot.WebContext) interface{} {
	r.Config.OutputType = knot.OutputTemplate
	r.Config.LayoutTemplate = LayoutFile
	r.Config.ViewName = "page-turbine-health.html"
	return w.GetParams(r, false)
}

func (w *PageController) DataSensorGovernance(r *knot.WebContext) interface{} {
	r.Config.OutputType = knot.OutputTemplate
	r.Config.LayoutTemplate = LayoutFile
	r.Config.ViewName = "page-data-sensor-governance.html"
	return w.GetParams(r, false)
}

func (w *PageController) TimeSeries(r *knot.WebContext) interface{} {
	r.Config.OutputType = knot.OutputTemplate
	r.Config.LayoutTemplate = LayoutFile
	r.Config.ViewName = "page-time-series.html"
	return w.GetParams(r, true)
}

func (w *PageController) DIYView(r *knot.WebContext) interface{} {
	r.Config.OutputType = knot.OutputTemplate
	r.Config.LayoutTemplate = LayoutFile
	r.Config.ViewName = "page-diy-view.html"
	return w.GetParams(r, false)
}

func (w *PageController) SCMManagement(r *knot.WebContext) interface{} {
	r.Config.OutputType = knot.OutputTemplate
	r.Config.LayoutTemplate = LayoutFile
	r.Config.ViewName = "page-scm.html"
	return w.GetParams(r, false)
}

func (w *PageController) IssueTracking(r *knot.WebContext) interface{} {
	r.Config.OutputType = knot.OutputTemplate
	r.Config.LayoutTemplate = LayoutFile
	r.Config.ViewName = "page-issue-tracking.html"
	return w.GetParams(r, false)
}

func (w *PageController) Reporting(r *knot.WebContext) interface{} {
	r.Config.OutputType = knot.OutputTemplate
	r.Config.LayoutTemplate = LayoutFile
	r.Config.ViewName = "page-reporting.html"
	return w.GetParams(r, false)
}

func (w *PageController) WindFarmAnalysis(r *knot.WebContext) interface{} {
	r.Config.OutputType = knot.OutputTemplate
	r.Config.LayoutTemplate = LayoutFile
	r.Config.ViewName = "page-windfarm-analysis.html"
	r.Config.IncludeFiles = append(DefaultIncludes, []string{"page-windfarm-analysis/project.html",
		"page-windfarm-analysis/turbine1.html",
		"page-windfarm-analysis/turbine2.html"}...)

	params := w.GetParams(r, true)
	params.Set("AvailableDate", r.Session("availdate"))

	return params
}
