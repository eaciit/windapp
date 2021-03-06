package web

import (
	. "eaciit/wfdemo-git/library/core"
	. "eaciit/wfdemo-git/library/models"
	"eaciit/wfdemo-git/web/controller"
	"os"
	"path/filepath"
	"runtime"

	"github.com/eaciit/knot/knot.v1"
	"github.com/eaciit/toolkit"
)

var (
	AppName  string = "web"
	basePath string = (func(dir string, err error) string { return dir }(os.Getwd()))

// 	server *knot.Server
)

func init() {
	app := knot.NewApp(AppName)
	app.ViewsPath = filepath.Join(basePath, AppName, "view") + toolkit.PathSeparator

	runtime.GOMAXPROCS(4)
	ConfigPath = controller.CONFIG_PATH

	if err := setAclDatabase(); err != nil {
		toolkit.Printf("Error set acl database : %s \n", err.Error())
	}

	app.Static("res", filepath.Join(controller.AppBasePath, AppName, "assets"))
	app.Static("image", filepath.Join(controller.AppBasePath, AppName, "assets", "img"))

	app.Register(controller.CreatePageController(AppName))
	app.Register(controller.CreateLoginController())
	app.Register(controller.CreateDataBrowserController())
	app.Register(controller.CreateDataBrowserNewController())
	app.Register(controller.CreateAccessController())
	app.Register(controller.CreateSessionController())
	app.Register(controller.CreateUserController())
	app.Register(controller.CreateGroupController())
	app.Register(controller.CreateLogController())
	app.Register(controller.CreateAnalyticKpiController())
	app.Register(controller.CreateDashboardController())

	app.Register(controller.CreateAnalyticWindRoseController())
	app.Register(controller.CreateAnalyticWindRoseDetailController())
	app.Register(controller.CreateAnalyticWindRoseFlexiController())
	app.Register(controller.CreateAnalyticWindRoseFlexiDetailController())

	app.Register(controller.CreateAnalyticLossAnalysisController())
	app.Register(controller.CreateAnalyticPowerCurveController())
	app.Register(controller.CreateAnalyticWindDistributionController())
	app.Register(controller.CreateAnalyticDgrScadaController())
	app.Register(controller.CreateAnalyticAvailabilityController())
	app.Register(controller.CreateAnalyticWindAvailabilityController())
	app.Register(controller.CreateAnalyticKeyMetricsController())
	app.Register(controller.CreateAnalyticComparisonController())
	app.Register(controller.CreateTrendLinePlotsController())
	app.Register(controller.CreateAnalyticHistogramController())
	app.Register(controller.CreateHelperController())
	app.Register(controller.CreateUserPreferencesController())
	app.Register(controller.CreateClusterWiseGenerationController())

	app.Register(controller.CreateDataEntryPowerCurveController())
	app.Register(controller.CreateDataEntryTurbineController())
	app.Register(controller.CreateDataEntryThresholdController())
	app.Register(controller.CreateAnalyticPerformanceIndexController())
	app.Register(controller.CreateAnalyticMeteorologyController())

	app.Register(controller.CreateTurbineHealthController())
	app.Register(controller.CreateDataSensorGovernanceController())
	app.Register(controller.CreateTimeSeriesController())
	app.Register(controller.CreateDiyViewController())
	app.Register(controller.CreateXyAnalysisController())

	app.Register(controller.CreateReportingController())

	app.Register(controller.CreateWindFarmAnalysisController())
	app.Register(controller.CreateMonitoringController())
	app.Register(controller.CreateDataAvailabilityController())
	app.Register(controller.CreateEmailController())

	app.Register(controller.CreateMonitoringRealtimeController())
	app.Register(controller.CreateMonitoringCustomController())

	app.Register(controller.CreateTurbineCollaborationController())
	app.Register(controller.CreateForecastController())

	app.Register(controller.CreateAnalyticDgrReportController())

	app.Register(controller.CreateThreeDAnalyticController())

	// app.Route("/", func(r *knot.WebContext) interface{} {
	// 	regex := regexp.MustCompile("/web/report/[a-zA-Z0-9_]+(/.*)?$")
	// 	rURL := r.Request.URL.String()

	// 	if regex.MatchString(rURL) {
	// 		args := strings.Split(strings.Replace(rURL, "/web/report/", "", -1), "/")
	// 		return WebController.PageReport(r, args)
	// 	}

	// 	sessionid := r.Session("sessionid", "")
	// 	if sessionid == "" {
	// 		http.Redirect(r.Writer, r.Request, "/web/login", 301)
	// 	} else {
	// 		http.Redirect(r.Writer, r.Request, "/web/report/dashboard", 301)
	// 	}

	// 	return true
	// })

	// server.Listen()

	// app.LayoutTemplate = "_template.html"
	knot.RegisterApp(app)
}

func setAclDatabase() error {
	if err := InitialSetDatabase(); err != nil {
		return err
	}
	pipeProject := []toolkit.M{
		toolkit.M{"$match": toolkit.M{"active": true}},
		toolkit.M{"$sort": toolkit.M{"projectname": 1}},
	}
	csrProject, err := DB().Connection.NewQuery().
		From("ref_project").
		Command("pipe", pipeProject).
		Cursor(nil)
	defer csrProject.Close()

	if err != nil {
		return err
	}

	dataProjects := []toolkit.M{}
	err = csrProject.Fetch(&dataProjects, 0, false)
	if err != nil {
		return err
	}

	for _, val := range dataProjects {
		controller.NotAvailLimit[val.GetString("projectid")] = val.GetFloat64("unavailablelimit")
	}

	if err := PrepareDefaultUser(); err != nil {
		return err
	}
	return nil
}
