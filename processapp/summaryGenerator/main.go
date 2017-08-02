package main

import (
	. "eaciit/wfdemo-git/processapp/summaryGenerator/controllers"
	. "eaciit/wfdemo-git/processapp/summaryGenerator/controllers/dataGenerator"
	webhelper "eaciit/wfdemo-git/web/helper"
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

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	tk.Println("Starting the app..")

	start := time.Now().UTC()

	db, e := PrepareConnection()
	if e != nil {
		tk.Println(e)
	} else {
		base := new(BaseController)
		base.Ctx = orm.New(db)
		defer base.Ctx.Close()

		webhelper.HelperSetDb(db)

		base.GetTurbineScada()
		base.PrepareDataReff()
		base.SetCollectionLatestTime()

		// dependent Generate
		// new(UpdateScadaOemMinutes).GenerateDensity(base)    // step 0
		// new(UpdateOEMToScada).RunMapping(base)              // step 1
		// new(EventToAlarm).ConvertEventToAlarm(base)         // step 2

		tk.Println("step 3")
		new(GenAlarmSummary).Generate(base) // step 3
		tk.Println("step 4")
		new(GenDataPeriod).GenerateMinify(base) // step 4
		tk.Println("step 5")
		new(GenScadaLast24).Generate(base) // step 5
		// tk.Println("step 6")
		// new(GenScadaSummary).Generate(base) // step 6
		// tk.Println("step 8")
		// new(GenScadaSummary).GenerateSummaryByProject(base) // step 8
		tk.Println("step 9")
		new(GenScadaSummary).GenerateSummaryDaily(base) // step 9
		tk.Println(">> step 9.6")
		new(GenScadaSummary).GenerateSummaryByMonthUsingDaily(base)
		tk.Println(">> step 9.8")
		new(GenScadaSummary).GenerateSummaryByProjectUsingDaily(base)
		tk.Println("step 10")
		new(GenScadaSummary).GenWFAnalysisByProject(base) // step 10
		tk.Println("step 11")
		new(GenScadaSummary).GenWFAnalysisByTurbine1(base) // step 11
		tk.Println("step 12")
		new(GenScadaSummary).GenWFAnalysisByTurbine2(base) // step 12

		// additional step for optimization perpose
		tk.Println("step additional 01")
		new(GenDataWindDistribution).GenerateCurrentMonth(base) // step add.01

		// not dependent Generate
		new(DataAvailabilitySummary).ConvertDataAvailabilitySummary(base)
		// new(EventReduceAvailability).ConvertEventReduceAvailability(base)

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

	tk.Printf("DONE in %v Minutes \n", time.Now().UTC().Sub(start).Minutes())
}
