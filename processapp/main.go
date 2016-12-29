package main

import (
	. "eaciit/wfdemo-git/processapp/controllers"
	. "eaciit/wfdemo-git/processapp/controllers/dataGenerator"
	// . "eaciit/wfdemo-git/processapp/controllers/excelConverter"

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

		// new(ConvScadaDataOEM).Generate(base)

		// =========================================================================================== //
		// step to prepare data for the application
		// =========================================================================================== //

		// new(UpdateScadaOemMinutes).GenerateDensity(base) // step 0

		// NewUpdateOEMToScada(base).RunMapping() // step 1

		NewEventToAlarm(base).ConvertEventToAlarm() // step 2

		new(GenAlarmSummary).Generate(base) // step 3

		NewGenDataPeriod(base).Generate() // step 4

		new(GenScadaLast24).Generate(base) // step 5

		new(GenScadaSummary).Generate(base)                 // step 6
		new(GenScadaSummary).GenerateSummaryByFleet(base)   // step 7
		new(GenScadaSummary).GenerateSummaryByProject(base) // step 8
		new(GenScadaSummary).GenerateSummaryDaily(base)     // step 9

		// =========================================================================================== //
		// step to prepare data for the application
		// =========================================================================================== //
	}

	tk.Println("Application Close..")
}
