package generatorControllers

import (
	. "eaciit/wfdemo/library/helper"
	. "eaciit/wfdemo/library/models"
	. "eaciit/wfdemo/processapp/controllers"
	"os"
	"sync"
	"time"

	"github.com/eaciit/crowd"
	"github.com/eaciit/dbox"
	_ "github.com/eaciit/dbox/dbc/mongo"
	tk "github.com/eaciit/toolkit"
)

var (
	maxDataEachProcess = 50
	muinsert           = &sync.Mutex{}
	funcName           = "GenAlarmScadaAnomaly Data"
)

// GenAlarmScadaAnomaly
type GenAlarmScadaAnomaly struct {
	*BaseController
}

// Generate
func (d *GenAlarmScadaAnomaly) Generate(base *BaseController) {
	startProcessTime := time.Now()

	if base != nil {
		d.BaseController = base

		ctx, e := PrepareConnection()
		if e != nil {
			ErrorHandler(e, funcName)
			os.Exit(0)
		}

		alarms := []Alarm{}
		csr, e := ctx.NewQuery().From(new(Alarm).TableName()).Cursor(nil)
		e = csr.Fetch(&alarms, 0, false)
		ErrorHandler(e, funcName)
		csr.Close()

		tk.Println("GenAlarmScadaAnomaly Data")

		totalData := len(alarms)
		worker := tk.ToInt(tk.ToFloat64(totalData/maxDataEachProcess, 0, tk.RoundingUp), tk.RoundingUp)

		wg := &sync.WaitGroup{}
		wg.Add(worker + 1)

		tk.Printf("totalData: %v worker: %v \n", totalData, worker+1)

		scadaDatas := []ScadaData{}

		csr, e = ctx.NewQuery().From(new(ScadaData).TableName()).Where(dbox.Lt("oktime", 600)).Cursor(nil)
		e = csr.Fetch(&scadaDatas, 0, false)
		ErrorHandler(e, funcName)
		csr.Close()

		scadaDatasAlarm := []ScadaData{}

		csr, e = ctx.NewQuery().From(new(ScadaData).TableName()).Where(dbox.Lt("alarmoktime", 600)).Cursor(nil)
		e = csr.Fetch(&scadaDatasAlarm, 0, false)
		ErrorHandler(e, funcName)
		csr.Close()

		for work := 0; work <= worker; work++ {
			from := work * maxDataEachProcess
			until := maxDataEachProcess + from

			if work == worker {
				until = totalData
			}

			// tk.Printf("%v | %v | %v \n", work, from, until)
			go d.doGen(alarms[from:until], scadaDatas, scadaDatasAlarm, wg)
		}

		wg.Wait()

		// tk.Printf("totaldata: %v \n", total)
		// tk.Printf("totaldata: %v \n", len(result))
	}
	duration := time.Since(startProcessTime)
	tk.Printf("finish in: %v \n", duration.String())
}

func (d *GenAlarmScadaAnomaly) doGen(alarms []Alarm, scadaDatas []ScadaData, scadaDatasAlarm []ScadaData, wg *sync.WaitGroup) {
	// ctx := d.BaseController.Ctx.Connection
	count := 0
	total := 0
	countClean := 0
	totalClean := 0
	// var result []orm.IModel

	for _, data := range alarms {
		startDate := GetDateRange(data.StartDate.UTC(), true)
		endDate := GetDateRange(data.EndDate.UTC(), false)

		exist := len(crowd.From(&scadaDatas).Where(func(x interface{}) interface{} {
			y := x.(ScadaData)
			isBefore := y.TimeStamp.UTC().Before(endDate.UTC()) || y.TimeStamp.UTC().Equal(endDate.UTC())
			isAfter := y.TimeStamp.UTC().After(startDate.UTC()) || y.TimeStamp.UTC().Equal(startDate.UTC())
			return isBefore && isAfter && y.Turbine == data.Turbine
		}).Exec().Result.Data().([]ScadaData))

		if exist == 0 {
			anomaly := new(AlarmScadaAnomaly).New()
			anomaly.Farm = data.Farm
			anomaly.StartDate = data.StartDate
			anomaly.StartDateInfo = data.StartDateInfo
			anomaly.EndDate = data.EndDate
			anomaly.Turbine = data.Turbine
			anomaly.AlertDescription = data.AlertDescription
			anomaly.ExternalStop = data.ExternalStop
			anomaly.GridDown = data.GridDown
			anomaly.InternalGrid = data.InternalGrid
			anomaly.MachineDown = data.MachineDown
			anomaly.AEbOK = data.AEbOK
			anomaly.Unknown = data.Unknown
			anomaly.WeatherStop = data.WeatherStop
			anomaly.Line = data.Line
			anomaly.Duration = data.Duration

			existOnAlarm := len(crowd.From(&scadaDatasAlarm).Where(func(x interface{}) interface{} {
				y := x.(ScadaData)
				isBefore := y.TimeStamp.UTC().Before(endDate.UTC()) || y.TimeStamp.UTC().Equal(endDate.UTC())
				isAfter := y.TimeStamp.UTC().After(startDate.UTC()) || y.TimeStamp.UTC().Equal(startDate.UTC())
				return isBefore && isAfter && y.Turbine == data.Turbine
			}).Exec().Result.Data().([]ScadaData))

			anomaly.IsAlarmOk = true
			if existOnAlarm == 0 {
				anomaly.IsAlarmOk = false
			}

			// result = append(result, anomaly)

			// muinsert.Lock()
			e := d.BaseController.Ctx.Insert(anomaly)
			ErrorHandler(e, funcName)
			// muinsert.Unlock()

			count++
			total++
		} else {
			clean := new(AlarmClean).New()
			clean.Farm = data.Farm
			clean.StartDate = data.StartDate
			clean.StartDateInfo = data.StartDateInfo
			clean.EndDate = data.EndDate
			clean.Turbine = data.Turbine
			clean.AlertDescription = data.AlertDescription
			clean.ExternalStop = data.ExternalStop
			clean.GridDown = data.GridDown
			clean.InternalGrid = data.InternalGrid
			clean.MachineDown = data.MachineDown
			clean.AEbOK = data.AEbOK
			clean.Unknown = data.Unknown
			clean.WeatherStop = data.WeatherStop
			clean.Line = data.Line
			clean.Duration = data.Duration

			/*existOnAlarm := len(crowd.From(&scadaDatasAlarm).Where(func(x interface{}) interface{} {
				y := x.(ScadaData)
				isBefore := y.TimeStamp.UTC().Before(endDate.UTC()) || y.TimeStamp.UTC().Equal(endDate.UTC())
				isAfter := y.TimeStamp.UTC().After(startDate.UTC()) || y.TimeStamp.UTC().Equal(startDate.UTC())
				return isBefore && isAfter && y.Turbine == data.Turbine
			}).Exec().Result.Data().([]ScadaData))

			clean.IsAlarmOk = true
			if existOnAlarm == 0 {
				clean.IsAlarmOk = false
			}*/

			// result = append(result, clean)

			// muinsert.Lock()
			e := d.BaseController.Ctx.Insert(clean)
			ErrorHandler(e, funcName)
			// muinsert.Unlock()

			countClean++
			totalClean++
		}

		// if count == 10 {
		// 	/*muinsert.Lock()
		// 	e := d.BaseController.Ctx.InsertBulk(result)
		// 	ErrorHandler(e, funcName)
		// 	muinsert.Unlock()*/

		// 	// result = []orm.IModel{}

		// 	tk.Printf("count: %v \n", total)
		// 	count = 0
		// }
	}

	/*if len(result) > 0 {
		muinsert.Lock()
		e := d.BaseController.Ctx.InsertBulk(result)
		ErrorHandler(e, funcName)
		muinsert.Unlock()

		tk.Printf("count: %v \n", total)
		count = 0
	}*/
	tk.Printf("count: %v \n", total)
	tk.Printf("count clean: %v \n", totalClean)
	wg.Done()
}
