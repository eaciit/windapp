package main

import (
<<<<<<< HEAD
	. "eaciit/wfdemo-git/processapp/controllers"
	// . "eaciit/wfdemo-git/processapp/controllers/dataGenerator"
	. "eaciit/wfdemo-git/processapp/controllers/excelConverter"
=======

	. "eaciit/wfdemo-git-dev/processapp/controllers"
	. "eaciit/wfdemo-git-dev/processapp/controllers/dataGenerator"
	. "eaciit/wfdemo-git-dev/processapp/controllers/excelConverter"
>>>>>>> a0bda0a9905baf4d0e67cbdc86fddbd6e8910e93
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

		// new(ConvTurbine).Generate(base)

		// new(ConvScadaData).Generate(base)
		// new(ConvScadaDataCsv).Generate(base)
		// new(GenScadaDataExceptionDurationTime).Generate(base)

		// new(ConvAlarm).Generate(base)
		// new(ConvAlarmBrakeMaster).Generate(base)

		// new(GenAlarmOverlapping).Generate(base)
		// new(ConvTurbine).Generate(base)

		// new(GenAlarmScadaAnomaly).Generate(base)
		// new(GenScadaAlarmAnomaly).Generate(base)

		// new(GenScadaSummary).Generate(base)
		// new(GenScadaSummary).GenerateSummaryByFleet(base)
		// new(GenScadaSummary).GenerateSummaryByProject(base)
		// new(GenScadaLast24).Generate(base)
		// new(GenScadaSummary).GenerateSummaryDaily(base)

		// new(UpdateAlarmPowerLost).Generate(base)
		// new(GenAlarmSummary).Generate(base)

		// new(GenScadaWindRose).Generate(base)
		// new(GenScadaWindRose).GenerateFromScadaNew(base)
		// new(GenScadaPowerCurve).GenerateAvg(base)
		// new(GenScadaPowerCurve).GenerateAdj(base)
		// new(GenScadaPowerCurvePlus).GeneratePlusAdj(base)
		// new(GenScadaPowerCurvePlus).GeneratePlusAvg(base)
		// new(UpdateScadaMinutes).Generate(base)
		// new(UpdateScadaMinutes).GenerateDensity(base)
		// new(UpdateProjectScadaAndAlarm).Generate(base)

		// new(ConvJMRBreakup).Generate(base)

		// new(ConvPermanentMetTower).Generate(base)
		// new(ConvPermanentMetTowerCSV).Generate(base)
		/*met := new(UpdateMetTower)
		met.Generate(base)
		met.GenerateWindRose(base)*/

		new(ConvScadaDataOEM).Generate(base)
	}

	tk.Println("Application Close..")
}
