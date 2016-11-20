package web

import (
	. "eaciit/ostrowfm/library/core"
	. "eaciit/ostrowfm/library/models"
	"eaciit/ostrowfm/web/controller"
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
	// port := new(Ports)
	// port.ID = "port"
	// if err := port.GetPort(); err != nil {
	// 	toolkit.Printf("Error get port: %s \n", err.Error())
	// }
	// if port.Port == 0 {
	// 	if err := setup.Port(port); err != nil {
	// 		toolkit.Printf("Error set port: %s \n", err.Error())
	// 	}
	// }

	if err := setAclDatabase(); err != nil {
		toolkit.Printf("Error set acl database : %s \n", err.Error())
	}

	// server.Address = toolkit.Sprintf("localhost:%v", toolkit.ToString(port.Port))
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
	app.Register(controller.CreateAnalyticHistogramController())
	app.Register(controller.CreateHelperController())
	app.Register(controller.CreateUserPreferencesController())

	app.Register(controller.CreateDataEntryPowerCurveController())

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

	if err := PrepareDefaultUser(); err != nil {
		return err
	}
	return nil
}
