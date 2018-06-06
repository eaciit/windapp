package main

import (
	. "eaciit/wfdemo-git/processapp/summaryGenerator/controllers"
	. "eaciit/wfdemo-git/processapp/summaryGenerator/controllers/dataGenerator"
	webhelper "eaciit/wfdemo-git/web/helper"
	"strings"
	"time"

	"os"
	"runtime"

	"github.com/eaciit/orm"
	tk "github.com/eaciit/toolkit"
)

var (
	wd = func() string {
		d, _ := os.Getwd()
		return d + "/"
	}()
)

const (
	sError   = "ERROR"
	sInfo    = "INFO"
	sWarning = "WARNING"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	tk.Println("Starting the app..")

	if runtime.GOOS != "windows" {
		res, err := tk.RunCommand("pidof", "ostroAllSummary.v02")
		ares := strings.Split(strings.TrimSpace(res), " ")
		if err != nil || len(ares) > 1 {
			tk.Printfn("EXIT FOUND : %s", res)
			os.Exit(0)
		}
	}

	start := time.Now().UTC()

	base := new(BaseController)
	base.Log, _ = tk.NewLog(true, false, "", "", "")

	db, e := PrepareConnection()
	if e != nil {
		tk.Println(e)
	} else {
		base.Ctx = orm.New(db)
		defer base.Ctx.Close()

		webhelper.HelperSetDb(db)

		base.GetTurbineScada()
		base.PrepareDataReff()
		// base.SetCollectionLatestTime()

		// dependent Generate
		// new(UpdateScadaOemMinutes).GenerateDensity(base)    // step 0
		// new(UpdateOEMToScada).RunMapping(base)              // step 1
		// new(EventToAlarm).ConvertEventToAlarm(base)         // step 2
		base.Log.AddLog("step 3", sInfo)
		new(GenAlarmSummary).Generate(base) // step 3
		base.Log.AddLog("step 4", sInfo)
		new(GenDataPeriod).GenerateMinify(base) // step 4
		base.Log.AddLog("step 5", sInfo)
		new(GenScadaLast24).Generate(base) // step 5
		// tk.Println("step 6")
		// new(GenScadaSummary).Generate(base) // step 6
		// tk.Println("step 8")
		// new(GenScadaSummary).GenerateSummaryByProject(base) // step 8
		base.Log.AddLog("step 9", sInfo)
		new(GenScadaSummary).GenerateSummaryDaily(base) // step 9
		base.Log.AddLog(">> step 9.6", sInfo)
		new(GenScadaSummary).GenerateSummaryByMonthUsingDaily(base)
		base.Log.AddLog(">> step 9.8", sInfo)
		new(GenScadaSummary).GenerateSummaryByProjectUsingDaily(base)
		base.Log.AddLog("step 10", sInfo)
		new(GenScadaSummary).GenWFAnalysisByProject(base) // step 10
		base.Log.AddLog("step 11", sInfo)
		new(GenScadaSummary).GenWFAnalysisByTurbine1(base) // step 11
		base.Log.AddLog("step 12", sInfo)
		new(GenScadaSummary).GenWFAnalysisByTurbine2(base) // step 12

		// additional step for optimization perpose
		base.Log.AddLog("step additional 01", sInfo)
		new(GenDataWindDistribution).GenerateCurrentMonth(base) // step add.01

		// not dependent Generate
		// new(DataAvailabilitySummary).ConvertDataAvailabilitySummary(base)
		// new(EventReduceAvailability).ConvertEventReduceAvailability(base)
		// new(DineuralProfileSummary).CreateDineuralProfileSummary(base)
		// new(TurbulenceIntensitySummary).CreateTurbulenceIntensitySummary(base)
		// new(TrendLinePlotSummary).CreateTrendLinePlotSummary(base)
		// new(TurbulenceIntensityGenerator).CreateTurbulenceIntensity10Min(base)

		// // custom function temporary running
		// new(UpdateScadaOemMinutes).UpdateDeviation(base)

		/* data that need to copy:

		Alarm
		EventDown
		ScadaData
		ScadaDataOEM
		EventRaw
		GWFAnalysisBy***
		LatestDataPeriod -> just copy the data that changed
		rpt_***
		DataAvailability
		*/
	}

	base.Log.AddLog(tk.Sprintf("DONE in %v Minutes \n", time.Now().UTC().Sub(start).Minutes()), sInfo)
}
