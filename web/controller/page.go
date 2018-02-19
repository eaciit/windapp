package controller

import (
	. "eaciit/wfdemo-git/library/models"
	"eaciit/wfdemo-git/web/helper"
	"net/http"

	"sync"

	"github.com/eaciit/knot/knot.v1"
	"github.com/eaciit/toolkit"
)

type PageController struct {
	App
	Params toolkit.M
	mux    sync.Mutex
}

var (
	DefaultIncludes = []string{"_head.html", "_menu.html", "_loader.html", "_script_template.html", "_dateinfo.html"}
)

func CreatePageController(AppName string) *PageController {
	var controller = new(PageController)
	controller.Params = toolkit.M{"AppName": AppName}
	return controller
}

func (w *PageController) GetParams(r *knot.WebContext, isAnalyst bool) toolkit.M {
	w.Params.Set("AntiCache", toolkit.RandomString(20))
	w.Params.Set("CurrentDateData", GetLastDateData())
	menus, _ := GetListOfMenu(toolkit.ToString(r.Session("sessionid", "")))
	w.Params.Set("Menus", menus)

	if isAnalyst {
		projectList, _ := helper.GetProjectList()
		turbineList, _ := helper.GetTurbineList(nil)
		temperatureList, _ := helper.GetTemperatureList()
		alarmTagList, _ := helper.GetAlarmTagsList()

		w.Params.Set("ProjectList", projectList)
		w.Params.Set("TurbineList", turbineList)
		w.Params.Set("TemperatureList", temperatureList)
		w.Params.Set("AlarmTagList", alarmTagList)
	}

	r.Writer.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	r.Writer.Header().Set("Pragma", "no-cache")
	r.Writer.Header().Set("Expires", "0")

	// check includeFiles
	var defIncludes []string
	if len(DefaultIncludes) > 0 {
		for _, d := range DefaultIncludes {
			if !toolkit.HasMember(r.Config.IncludeFiles, d) {
				defIncludes = append(defIncludes, d)
			}
		}
	}

	r.Config.IncludeFiles = append(r.Config.IncludeFiles, defIncludes...)
	// r.Config.IncludeFiles = append(r.Config.IncludeFiles, DefaultIncludes...)
	// WriteLog(r.Session("sessionid", ""), "access", r.Request.URL.String())
	return w.Params
}

/*func (w *PageController) Index(r *knot.WebContext) interface{} {
	r.Config.OutputType = knot.OutputTemplate
	r.Config.LayoutTemplate = LayoutFile
	r.Config.ViewName = "page-index.html"
	return w.GetParams(r, false)
}*/

func (w *PageController) Login(r *knot.WebContext) interface{} {
	r.Config.OutputType = knot.OutputTemplate
	r.Config.ViewName = "page-login.html"

	if r.Session("sessionid", "") != "" {
		w.mux.Lock()
		r.SetSession("sessionid", "")
		w.mux.Unlock()
	}

	return w.GetParams(r, false)
}

func (w *PageController) DataBrowser(r *knot.WebContext) interface{} {
	r.Config.OutputType = knot.OutputTemplate
	r.Config.LayoutTemplate = LayoutFile
	r.Config.ViewName = "page-databrowser.html"
	return w.GetParams(r, true).Set("ColumnList", GetCustomFieldList()).Set("HFDColList", GetHFDCustomFieldList())
}

/*func (w *PageController) DataBrowserNew(r *knot.WebContext) interface{} {
	r.Config.OutputType = knot.OutputTemplate
	r.Config.LayoutTemplate = LayoutFile
	r.Config.ViewName = "page-databrowser-new.html"

	return w.GetParams(r, false).Set("ColumnList", GetCustomFieldList()).Set("HDFColList", GetHFDCustomFieldList())
}*/

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
	r.Config.IncludeFiles = []string{"_filter-powercurve.html"}
	r.Config.ViewName = "page-analytic-power-curve.html"

	return w.GetParams(r, true)
}

func (w *PageController) AnalyticPCMonthly(r *knot.WebContext) interface{} {
	r.Config.OutputType = knot.OutputTemplate
	r.Config.LayoutTemplate = LayoutFile
	r.Config.ViewName = "page-analytic-power-curve/individual-month.html"

	return w.GetParams(r, true)
}

func (w *PageController) AnalyticPCMonthlyScatter(r *knot.WebContext) interface{} {
	r.Config.OutputType = knot.OutputTemplate
	r.Config.LayoutTemplate = LayoutFile
	r.Config.ViewName = "page-analytic-power-curve/individual-month-scatter.html"

	return w.GetParams(r, true)
}

func (w *PageController) AnalyticPCComparison(r *knot.WebContext) interface{} {
	r.Config.OutputType = knot.OutputTemplate
	r.Config.LayoutTemplate = LayoutFile
	r.Config.ViewName = "page-analytic-power-curve/comparison.html"

	return w.GetParams(r, true)
}

func (w *PageController) AnalyticTrendLinePlots(r *knot.WebContext) interface{} {
	r.Config.OutputType = knot.OutputTemplate
	r.Config.LayoutTemplate = LayoutFile
	r.Config.ViewName = "trend-line-plots/trendlineplots.html"

	return w.GetParams(r, true)
}

func (w *PageController) AnalyticPCScatter(r *knot.WebContext) interface{} {
	r.Config.OutputType = knot.OutputTemplate
	r.Config.LayoutTemplate = LayoutFile
	r.Config.ViewName = "page-analytic-power-curve/scatter.html"
	r.Config.IncludeFiles = []string{"_filter-scatter.html"}

	return w.GetParams(r, true)
}

func (w *PageController) AnalyticPCScatterAnalysis(r *knot.WebContext) interface{} {
	r.Config.OutputType = knot.OutputTemplate
	r.Config.LayoutTemplate = LayoutFile
	r.Config.ViewName = "page-analytic-power-curve/scatter-analysis.html"
	r.Config.IncludeFiles = []string{"_filter-scatter.html"}

	return w.GetParams(r, true)
}

func (w *PageController) AnalyticPCScatterOperational(r *knot.WebContext) interface{} {
	r.Config.OutputType = knot.OutputTemplate
	r.Config.LayoutTemplate = LayoutFile
	r.Config.ViewName = "page-analytic-power-curve/scatter-operational.html"
	r.Config.IncludeFiles = []string{"page-analytic-power-curve/page-filter-scatter.html"}

	return w.GetParams(r, true)
}

/*func (w *PageController) AnalyticDgrScada(r *knot.WebContext) interface{} {
	r.Config.OutputType = knot.OutputTemplate
	r.Config.LayoutTemplate = LayoutFile
	r.Config.ViewName = "page-analytic-dgr-scada.html"
	return w.GetParams(r, true)
}*/

func (w *PageController) ClusterWiseGeneration(r *knot.WebContext) interface{} {
	r.Config.OutputType = knot.OutputTemplate
	r.Config.LayoutTemplate = LayoutFile
	r.Config.ViewName = "page-cluster-wise-generation.html"

	return w.GetParams(r, true)
}

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
	r.Config.IncludeFiles = []string{"_filter-analytic.html",
		"page-loss-analysis/static-view.html",
		"page-loss-analysis/downtime-top10.html",
		"page-loss-analysis/availability.html",
		"page-loss-analysis/lost-energy.html",
		"page-loss-analysis/windspeed-availability.html",
		"page-loss-analysis/warning-frequency.html",
		"page-loss-analysis/component-alarm.html",
		"page-loss-analysis/mtbf.html",
	}
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
	r.Config.IncludeFiles = []string{"_filter-analytic.html",
		"page-analytic-meteorology/turbulence-intensity.html",
		"page-analytic-meteorology/table1224.html",
		"page-analytic-meteorology/windrose.html",
		"page-analytic-meteorology/windrose-comparison.html",
		"page-analytic-meteorology/wind-distribution.html",
		"page-analytic-meteorology/nacelle-distribution.html",
		"page-analytic-meteorology/average-windspeed.html",
		"page-analytic-meteorology/turbine-correlation.html",
		"page-analytic-meteorology/energy-correlation.html",
	}
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
func (w *PageController) DataEntryTurbine(r *knot.WebContext) interface{} {
	r.Config.OutputType = knot.OutputTemplate
	r.Config.LayoutTemplate = LayoutFile
	r.Config.ViewName = "page-dataentry-turbine.html"

	return w.GetParams(r, true)
}
func (w *PageController) DataEntryThreshold(r *knot.WebContext) interface{} {
	r.Config.OutputType = knot.OutputTemplate
	r.Config.LayoutTemplate = LayoutFile
	r.Config.ViewName = "page-strangethreshold.html"

	return w.GetParams(r, true)
}
func (w *PageController) DgrReport(r *knot.WebContext) interface{} {
	r.Config.OutputType = knot.OutputTemplate
	r.Config.LayoutTemplate = LayoutFile
	r.Config.ViewName = "page-dgr-report.html"

	return w.GetParams(r, true)
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

func (w *PageController) EmailManagement(r *knot.WebContext) interface{} {
	r.Config.OutputType = knot.OutputTemplate
	r.Config.LayoutTemplate = LayoutFile
	r.Config.ViewName = "page-email-management.html"

	categorylist, _ := GetCategoryMail()
	userList, _ := GetUserMail()
	alarmCodes, _ := GetAlarmCodesMail()
	template, _ := GetTemplateMail()

	w.Params.Set("CategoryMailList", categorylist)
	w.Params.Set("UserMailList", userList)
	w.Params.Set("AlarmCodesMailList", alarmCodes)
	w.Params.Set("TemplateMailList", template)

	return w.GetParams(r, false)
}

func (w *PageController) Monitoring(r *knot.WebContext) interface{} {
	r.Config.OutputType = knot.OutputTemplate
	r.Config.LayoutTemplate = LayoutFile
	r.Config.ViewName = "page-monitoring.html"
	r.Config.IncludeFiles = []string{"_filter-monitoring.html"}

	if r.Session("sessionid", "") == "" {
		w.mux.Lock()
		r.SetSession("sessionid", "monitoring-page")
		w.mux.Unlock()
	}

	return w.GetParams(r, true)
}

func (w *PageController) MonitoringByProject(r *knot.WebContext) interface{} {
	r.Config.OutputType = knot.OutputTemplate
	r.Config.LayoutTemplate = LayoutFile
	r.Config.ViewName = "page-monitoring/by-project.html"
	r.Config.IncludeFiles = []string{"page-monitoring/turbine-collaboration.html"}

	if r.Session("sessionid", "") == "" {
		w.mux.Lock()
		r.SetSession("sessionid", "monitoring-page")
		w.mux.Unlock()
	}

	return w.GetParams(r, true)
}
func (w *PageController) MonitoringAllProject(r *knot.WebContext) interface{} {
	r.Config.OutputType = knot.OutputTemplate
	r.Config.LayoutTemplate = LayoutFile
	r.Config.ViewName = "page-monitoring/all-project.html"

	if r.Session("sessionid", "") == "" {
		w.mux.Lock()
		r.SetSession("sessionid", "monitoring-page")
		w.mux.Unlock()
	}

	return w.GetParams(r, true)
}
func (w *PageController) MonitoringSummary(r *knot.WebContext) interface{} {
	r.Config.OutputType = knot.OutputTemplate
	r.Config.LayoutTemplate = LayoutFile
	r.Config.ViewName = "page-monitoring/summary.html"
	r.Config.IncludeFiles = []string{"page-monitoring/all-project.html", "page-monitoring/custom.html", "page-monitoring/turbine-collaboration.html"}

	if r.Session("sessionid", "") == "" {
		w.mux.Lock()
		r.SetSession("sessionid", "monitoring-page")
		w.mux.Unlock()
	}

	return w.GetParams(r, true)
}
func (w *PageController) MonitoringByProjectCustom(r *knot.WebContext) interface{} {
	r.Config.OutputType = knot.OutputTemplate
	r.Config.LayoutTemplate = LayoutFile
	r.Config.ViewName = "page-monitoring/custom.html"
	r.Config.IncludeFiles = []string{"page-monitoring/turbine-collaboration.html"}

	if r.Session("sessionid", "") == "" {
		w.mux.Lock()
		r.SetSession("sessionid", "monitoring-page")
		w.mux.Unlock()
	}

	return w.GetParams(r, true)
}
func (w *PageController) MonitoringByTurbine(r *knot.WebContext) interface{} {
	r.Config.OutputType = knot.OutputTemplate
	r.Config.LayoutTemplate = LayoutFile
	r.Config.IncludeFiles = []string{"page-monitoring/turbine-pic.html", "page-monitoring/turbine-collaboration.html"}
	r.Config.ViewName = "page-monitoring/individual-turbine.html"
	// allTurbineList, _ := helper.GetAllTurbineList()
	// w.Params.Set("AllTurbineList", allTurbineList)

	if r.Session("sessionid", "") == "" {
		w.mux.Lock()
		r.SetSession("sessionid", "monitoring-page")
		w.mux.Unlock()
	}

	return w.GetParams(r, true)
}

func (w *PageController) MonitoringAlarm(r *knot.WebContext) interface{} {
	r.Config.OutputType = knot.OutputTemplate
	r.Config.LayoutTemplate = LayoutFile
	r.Config.IncludeFiles = []string{"_filter-analytic.html"}
	r.Config.ViewName = "page-monitoring/monitoring-alarm.html"

	if r.Session("sessionid", "") == "" {
		w.mux.Lock()
		r.SetSession("sessionid", "monitoring-page")
		w.mux.Unlock()
	}

	return w.GetParams(r, true)
}

func (w *PageController) MonitoringTemperature(r *knot.WebContext) interface{} {
	r.Config.OutputType = knot.OutputTemplate
	r.Config.LayoutTemplate = LayoutFile
	r.Config.ViewName = "page-monitoring/temperature.html"

	if r.Session("sessionid", "") == "" {
		w.mux.Lock()
		r.SetSession("sessionid", "monitoring-page")
		w.mux.Unlock()
	}

	return w.GetParams(r, true)
}

func (w *PageController) MonitoringNotification(r *knot.WebContext) interface{} {
	r.Config.OutputType = knot.OutputTemplate
	r.Config.LayoutTemplate = LayoutFile
	r.Config.IncludeFiles = []string{"_filter-analytic.html"}
	r.Config.ViewName = "page-monitoring/notification.html"

	if r.Session("sessionid", "") == "" {
		w.mux.Lock()
		r.SetSession("sessionid", "monitoring-page")
		w.mux.Unlock()
	}

	return w.GetParams(r, true)
}

func (w *PageController) MonitoringWeather(r *knot.WebContext) interface{} {
	r.Config.OutputType = knot.OutputTemplate
	r.Config.LayoutTemplate = LayoutFile
	r.Config.ViewName = "page-monitoring/weather-forecast.html"

	if r.Session("sessionid", "") == "" {
		w.mux.Lock()
		r.SetSession("sessionid", "monitoring-page")
		w.mux.Unlock()
	}

	return w.GetParams(r, true)
}

func (w *PageController) Dashboard(r *knot.WebContext) interface{} {
	r.Config.OutputType = knot.OutputTemplate
	r.Config.LayoutTemplate = LayoutFile
	r.Config.ViewName = "page-dashboard.html"
	r.Config.IncludeFiles = []string{"page-dashboard-summary.html", "page-dashboard-production.html", "page-dashboard-availability.html"}

	projectList, _ := helper.GetProjectList()
	w.Params.Set("ProjectList", projectList)

	return w.GetParams(r, false)
}
func (w *PageController) Starter(r *knot.WebContext) interface{} {
	r.Config.OutputType = knot.OutputTemplate
	r.Config.LayoutTemplate = LayoutFile
	r.Config.ViewName = "page-starter.html"
	// r.Config.IncludeFiles = []string{"page-dashboard-summary.html", "page-dashboard-production.html", "page-dashboard-availability.html"}
	projectList, _ := helper.GetProjectList()
	w.Params.Set("ProjectList", projectList)

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

func (w *PageController) Forecasting(r *knot.WebContext) interface{} {
	r.Config.OutputType = knot.OutputTemplate
	r.Config.LayoutTemplate = LayoutFile
	r.Config.ViewName = "page-forecast.html"
	return w.GetParams(r, true)
}

func (w *PageController) ThreeDAnalytic(r *knot.WebContext) interface{} {
	r.Config.OutputType = knot.OutputTemplate
	r.Config.LayoutTemplate = LayoutFile
	r.Config.ViewName = "page-threedanalytic.html"
	return w.GetParams(r, true)
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
	r.Config.IncludeFiles = []string{"page-analytic-power-curve/page-filter-scatter.html"}
	return w.GetParams(r, true).Set("PageType", "OEM")
}

func (w *PageController) TimeSeriesHFD(r *knot.WebContext) interface{} {
	r.Config.OutputType = knot.OutputTemplate
	r.Config.LayoutTemplate = LayoutFile
	r.Config.ViewName = "page-time-series.html"
	r.Config.IncludeFiles = []string{"page-analytic-power-curve/page-filter-scatter.html"}
	return w.GetParams(r, true).Set("PageType", "HFD")
}

func (w *PageController) XYAnalysis(r *knot.WebContext) interface{} {
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
	r.Config.IncludeFiles = []string{"page-windfarm-analysis/project.html",
		"page-windfarm-analysis/turbine1.html",
		"page-windfarm-analysis/turbine2.html"}

	params := w.GetParams(r, true)

	return params
}

func (w *PageController) DataAvailability(r *knot.WebContext) interface{} {
	r.Config.OutputType = knot.OutputTemplate
	r.Config.LayoutTemplate = LayoutFile
	r.Config.ViewName = "page-data-availability.html"

	params := w.GetParams(r, true)

	return params
}
