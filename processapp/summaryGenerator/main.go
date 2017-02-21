package main

import (
	. "eaciit/wfdemo-git/processapp/summaryGenerator/controllers"
	. "eaciit/wfdemo-git/processapp/summaryGenerator/controllers/dataGenerator"

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

	db, e := PrepareConnection()
	if e != nil {
		tk.Println(e)
	} else {
		base := new(BaseController)
		base.Ctx = orm.New(db)
		defer base.Ctx.Close()

		// x, _ := GetPowerCurveCubicInterpolation(base.Ctx.Connection, "Tejuva", 3.3)
		// log.Println(x)

		base.SetCollectionLatestTime()
		base.PrepareDataReff()

		// new(UpdateScadaOemMinutes).GenerateDensity(base)    // step 0
		// new(UpdateOEMToScada).RunMapping(base)              // step 1
		// new(EventToAlarm).ConvertEventToAlarm(base)         // step 2
		// new(GenAlarmSummary).Generate(base)                 // step 3
		// new(GenDataPeriod).Generate(base)                   // step 4
		// new(GenScadaLast24).Generate(base)                  // step 5
		// new(GenScadaSummary).Generate(base)                 // step 6
		// new(GenScadaSummary).GenerateSummaryByFleet(base)   // step 7
		new(GenScadaSummary).GenerateSummaryByProject(base) // step 8
		// new(GenScadaSummary).GenerateSummaryDaily(base) // step 9
		// new(GenScadaSummary).GenWFAnalysisByProject(base)  // step 10
		// new(GenScadaSummary).GenWFAnalysisByTurbine1(base) // step 11
		// new(GenScadaSummary).GenWFAnalysisByTurbine2(base) // step 12
	}

	tk.Println("Application Close..")
}
