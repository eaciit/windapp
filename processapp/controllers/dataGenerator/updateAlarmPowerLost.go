package generatorControllers

import (
	. "eaciit/ostrowfm/library/helper"
	. "eaciit/ostrowfm/library/models"
	. "eaciit/ostrowfm/processapp/controllers"
	"os"
	"sync"
	"time"

	"github.com/eaciit/crowd"
	"github.com/eaciit/dbox"
	_ "github.com/eaciit/dbox/dbc/mongo"
	tk "github.com/eaciit/toolkit"
)

var (
	maxDataProcess = 50
)

// UpdateAlarmPowerLost
type UpdateAlarmPowerLost struct {
	*BaseController
}

// Generate
func (d *UpdateAlarmPowerLost) Generate(base *BaseController) {
	funcName := "UpdateAlarmPowerLost Data"
	startProcessTime := time.Now()

	/*fromDate, _ := time.Parse("2006-01-02 15:04 MST", "2016-06-30 00:00 UTC")
	toDate, _ := time.Parse("2006-01-02 15:04 MST", "2016-06-30 23:59 UTC")

	var filter []*dbox.Filter
	filter = append(filter, dbox.Gte("startdate", fromDate))
	filter = append(filter, dbox.Lte("startdate", toDate))
	filter = append(filter, dbox.Eq("turbine", "HBR006"))

	var filterX []*dbox.Filter
	filterX = append(filterX, dbox.Gte("timestamp", fromDate))
	filterX = append(filterX, dbox.Lte("timestamp", toDate))
	filterX = append(filterX, dbox.Eq("turbine", "HBR006"))*/

	if base != nil {
		d.BaseController = base

		ctx, e := PrepareConnection()
		if e != nil {
			ErrorHandler(e, funcName)
			os.Exit(0)
		}

		alarms := []AlarmClean{}
		// csr, e := ctx.NewQuery().From(new(Alarm).TableName()).Where(filter...).Cursor(nil)
		csr, e := ctx.NewQuery().From(new(AlarmClean).TableName()).Cursor(nil)
		e = csr.Fetch(&alarms, 0, false)
		ErrorHandler(e, funcName)
		csr.Close()

		tk.Println("UpdateAlarmPowerLost Data")

		totalData := len(alarms)
		worker := tk.ToInt(tk.ToFloat64(totalData/maxDataProcess, 0, tk.RoundingUp), tk.RoundingUp)

		wg := &sync.WaitGroup{}
		wg.Add(worker + 1)

		tk.Printf("totalData: %v worker: %v \n", totalData, worker+1)

		scadaDatas := []ScadaClean{}

		// csr, e = ctx.NewQuery().From(new(ScadaData).TableName()).Where(filterX...).Order("timestamp").Cursor(nil)

		var pipes []tk.M

		pipes = append(pipes, tk.M{"$sort": tk.M{"timestamp": 1}})

		// csr, e = ctx.NewQuery().From(new(ScadaData).TableName()).Order("timestamp").Cursor(nil)

		csr, e = ctx.NewQuery().
			From(new(ScadaClean).TableName()).
			Command("pipe", pipes).
			Cursor(nil)

		e = csr.Fetch(&scadaDatas, 0, false)
		ErrorHandler(e, funcName)
		csr.Close()

		powerCurve := []PowerCurveModel{}
		csr, e = ctx.NewQuery().From(new(PowerCurveModel).TableName()).Cursor(nil)
		e = csr.Fetch(&powerCurve, 0, false)
		ErrorHandler(e, funcName)
		csr.Close()

		for work := 0; work <= worker; work++ {
			from := work * maxDataProcess
			until := maxDataProcess + from

			if work == worker {
				until = totalData
			}

			// tk.Printf("%v | %v | %v \n", work, from, until)
			go d.doGen(alarms[from:until], scadaDatas, powerCurve, wg)
		}

		wg.Wait()

		// tk.Printf("totaldata: %v \n", total)
		// tk.Printf("totaldata: %v \n", len(result))
	}
	duration := time.Since(startProcessTime)
	tk.Printf("finish in: %v \n", duration.String())
}

func (d *UpdateAlarmPowerLost) doGen(alarms []AlarmClean, scadaDatas []ScadaClean, powerCurve []PowerCurveModel, wg *sync.WaitGroup) {
	// ctx := d.BaseController.Ctx.Connection
	count := 0
	total := 0
	// var result []orm.IModel

	for _, data := range alarms {
		startDate := GetDateRange(data.StartDate.UTC(), false)
		endDate := GetDateRange(data.EndDate.UTC(), false)

		exist := crowd.From(&scadaDatas).Where(func(x interface{}) interface{} {
			y := x.(ScadaClean)
			isBefore := y.TimeStamp.UTC().Before(endDate.UTC()) || y.TimeStamp.UTC().Equal(endDate.UTC())
			isAfter := y.TimeStamp.UTC().After(startDate.UTC()) || y.TimeStamp.UTC().Equal(startDate.UTC())
			return isBefore && isAfter && (y.Turbine == data.Turbine)
		}).Exec().Result.Data().([]ScadaClean)

		var totalPowerLost float64

		// tk.Printf("len(exist): %v \n\n", len(exist))

		for idx, val := range exist {
			// tk.Printf("timestamp: %v | ", val.TimeStamp.UTC().String())

			startDuration := val.TimeStamp.UTC().Sub(data.StartDate.UTC()).Hours()
			endDuration := (10 - val.TimeStamp.UTC().Sub(data.EndDate.UTC()).Minutes()) / 60
			windSpeed := tk.ToFloat64(val.AvgWindSpeed, 0, tk.RoundingDown)

			powers := crowd.From(&powerCurve).Where(func(x interface{}) interface{} {
				y := x.(PowerCurveModel)
				return y.WindSpeed == windSpeed
			}).Exec().Result.Data().([]PowerCurveModel)

			var powerLost float64

			if len(powers) > 0 {
				power := powers[0].Power1

				if idx == 0 {
					// tk.Printf("#%v (%v * %v)", windSpeed, startDuration, (power + 0.0))
					powerLost = startDuration * (power + 0.0)
				} else if idx == len(exist)-1 {
					// tk.Printf("#%v (%v * %v)", windSpeed, endDuration, (power + 0.0))
					powerLost = endDuration * (power + 0.0)
				} else {
					// tk.Printf("#%v (%v * %v)", windSpeed, (10.0 / 60.0), (power + 0.0))
					powerLost = (10 / 60) * (power + 0.0)
				}

				// tk.Printf(" | %v \n", powerLost)
			}

			totalPowerLost += powerLost
			// tk.Printf(" | %v \n", totalPowerLost)
		}

		// muinsert.Lock()
		e := d.BaseController.Ctx.Connection.NewQuery().Update().From(new(AlarmClean).TableName()).Where(dbox.Eq("_id", data.ID)).Exec(tk.M{}.Set("data", tk.M{}.Set("powerlost", totalPowerLost)))
		ErrorHandler(e, funcName)
		// muinsert.Unlock()

		// tk.Printf("totalPowerLost: %v \n", totalPowerLost)

		count++
		total++
	}
	tk.Printf("count: %v \n", total)
	wg.Done()
}
