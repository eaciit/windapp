package generatorControllers

import (
	. "eaciit/wfdemo-git/library/helper"
	. "eaciit/wfdemo-git/library/models"
	. "eaciit/wfdemo-git/processapp/controllers"
	_ "fmt"
	"os"
	"strconv"
	_ "strings"
	"time"

	"github.com/eaciit/dbox"
	_ "github.com/eaciit/dbox/dbc/mongo"
	tk "github.com/eaciit/toolkit"
)

type GenScadaLast24 struct {
	*BaseController
}

func (d *GenScadaLast24) Generate(base *BaseController) {
	if base != nil {
		d.BaseController = base
		ctx, e := PrepareConnection()
		if e != nil {
			ErrorHandler(e, "Scada Summary")
			os.Exit(0)
		}

		/*csr, e := ctx.NewQuery().From(new(ScadaData).TableName()).
			Aggr(dbox.AggrMax, "$timestamp", "timestamp").
			Aggr(dbox.AggrMax, "$dateinfo.dateid", "dateid").
			Group("").
			Cursor(nil)
		defer csr.Close()*/

		d.BaseController.Ctx.DeleteMany(new(ScadaLastUpdate), dbox.And(dbox.Ne("_id", "")))

		csr, e := ctx.NewQuery().From(new(ScadaData).TableName()).
			//Where(dbox.Gte("power", 0)).
			Aggr(dbox.AggrMax, "$timestamp", "timestamp").
			Aggr(dbox.AggrMax, "$dateinfo.dateid", "dateid").
			Group("").
			Cursor(nil)
		defer csr.Close()

		datas := []tk.M{}
		e = csr.Fetch(&datas, 0, false)

		budgetMonths := []float64{
			5911.8744,
			6023.419200000001,
			7027.3224,
			8588.9496,
			14389.2792,
			16954.8096,
			15727.8168,
			12046.8384,
			9704.3976,
			5688.784799999999,
			3569.4336,
			5911.8744}

		//tk.Printf("#%v", datas)

		if datas != nil {
			dateId := datas[0]["dateid"].(time.Time).UTC()
			dtInfo := GetDateInfo(dateId)
			maxTimeStamp := datas[0]["timestamp"].(time.Time).UTC()
			//startTime := maxTimeStamp.Add(-24 * time.Hour)
			var budgetCurrMonthDaily float64

			if len(budgetMonths)-1 >= int(dateId.Month()) {
				budgetCurrMonths := budgetMonths[int(dateId.Month())] * 1000.0
				noOfDay := float64(daysIn(dateId.Month(), dateId.Year()))
				budgetCurrMonthDaily = tk.Div(budgetCurrMonths, noOfDay)
			}

			mdl := new(ScadaLastUpdate).New()
			mdl.ID = "SCADALASTUPDATE_FLEET"
			mdl.LastUpdate = maxTimeStamp
			mdl.DateInfo = dtInfo
			mdl.NoOfProjects = 1
			mdl.NoOfTurbines = 24
			mdl.TotalMaxCapacity = 24 * 2100
			mdl.CurrentDown = 0
			mdl.TwoDaysDown = 0
			mdl.ProjectName = "Fleet"

			items := make([]LastData24Hours, 0)
			for i := 0; i < 24; i++ {
				year := strconv.Itoa(dateId.Year())
				month := dateId.Month().String()
				day := strconv.Itoa(dateId.Day())
				strTime := year + "-" + month + "-" + day + " " + strconv.Itoa(i) + ":00:00"
				timeHr, _ := time.Parse("2006-January-2 15:04:05", strTime)

				timeHrStart := timeHr.Add(-1 * time.Hour)

				/*csr, e = ctx.NewQuery().From(new(ScadaData).TableName()).
					Where(dbox.And(dbox.Gt("timestamp", timeHrStart), dbox.Lte("timestamp", timeHr))).
					Aggr(dbox.AggrSum, "$power", "totalpower").
					Aggr(dbox.AggrSum, "$powerlost", "totalpowerlost").
					Aggr(dbox.AggrSum, "$oktime", "totaloktime").
					Aggr(dbox.AggrSum, "$griddowntime", "totalgriddowntime").
					Aggr(dbox.AggrAvr, "$windspeed", "avgwindspeed").
					Group("projectname").
					Cursor(nil)
				defer csr.Close()*/

				csr, e = ctx.NewQuery().From(new(ScadaData).TableName()).
					Where(dbox.And(dbox.Gt("timestamp", timeHrStart), dbox.Lte("timestamp", timeHr), dbox.Gte("power", 0))).
					Aggr(dbox.AggrSum, "$power", "totalpower").
					Aggr(dbox.AggrSum, "$powerlost", "totalpowerlost").
					Aggr(dbox.AggrSum, "$energylost", "energylost").
					Aggr(dbox.AggrSum, "$denpower", "denpower").
					Aggr(dbox.AggrSum, "$oktime", "totaloktime").
					Aggr(dbox.AggrSum, "$griddowntime", "totalgriddowntime").
					Aggr(dbox.AggrAvr, "$windspeed", "avgwindspeed").
					Group("projectname").
					Cursor(nil)
				defer csr.Close()

				scadas := []tk.M{}
				e = csr.Fetch(&scadas, 0, false)

				//tk.Printf("#%v\n", data)

				var last LastData24Hours
				if len(scadas) > 0 {
					data := scadas[0]
					trueAvail := 0.0
					gridAvail := 0.0

					// tk.Println(data)

					//id := data["_id"].(tk.M)

					ipower := data["totalpower"]
					power := 0.0
					if ipower != nil {
						power = ipower.(float64)
					}

					// ienergylost := data["energylost"]
					// energylost := 0.0
					// if ienergylost != nil {
					// 	energylost = ienergylost.(float64)
					// }

					ipotentialpower := data["denpower"]
					potentialpower := 0.0
					if ipotentialpower != nil {
						potentialpower = ipotentialpower.(float64)
					}

					// powerlost := 0.0
					// pipe := []tk.M{tk.M{}.Set("$match", tk.M{}.Set("$and", tk.M{}.Set("enddate", tk.M{}.Set("$lte", timeHr)).Set("startdate", tk.M{}.Set("$gt", timeHrStart)))), tk.M{}.Set("$group", tk.M{}.Set("_id", "$projectname").Set("duration", tk.M{}.Set("$sum", "$duration")).Set("powerlost", tk.M{}.Set("$sum", "$powerlost"))), tk.M{}.Set("$sort", tk.M{}.Set("_id", 1))}
					// csr1, _ := ctx.NewQuery().
					// 	Command("pipe", pipe).
					// 	From(new(Alarm).TableName()).
					// 	Cursor(nil)
					// defer csr1.Close()

					// alarms := []tk.M{}
					// e = csr1.Fetch(&alarms, 0, false)

					// if len(alarms) > 0 {
					// 	alarm := alarms[0]

					// 	ipowerlost := alarm["powerlost"]
					// 	if ipowerlost != nil {
					// 		powerlost = ipowerlost.(float64)
					// 	}
					// }

					iwindspeed := data["avgwindspeed"]
					windspeed := 0.0
					if iwindspeed != nil {
						windspeed = iwindspeed.(float64)
					}
					last.Hour = i
					last.TimeHour = timeHr
					last.AvgWindSpeed = windspeed
					last.PowerKw = power
					last.EnergyKwh = power / 6
					last.Potential = potentialpower
					last.PotentialKwh = potentialpower / 6
					last.TrueAvail = trueAvail
					last.GridAvail = gridAvail
				} else {
					last.Hour = i
					last.TimeHour = timeHr
					last.AvgWindSpeed = 0.0
					last.PowerKw = 0.0
					last.EnergyKwh = 0.0
					last.Potential = 0.0
					last.PotentialKwh = 0.0
					last.TrueAvail = 0.0
					last.GridAvail = 0.0
				}

				items = append(items, last)
			}

			// item30s := make([]Last30Days, 0)
			// for i := 0; i < 30; i++ {
			// 	year := strconv.Itoa(dateId.Year())
			// 	month := dateId.Month().String()
			// 	day := strconv.Itoa(i + 1)
			// 	strDate := year + "-" + month + "-" + day + " 00:00:00"
			// 	date, _ := time.Parse("2006-January-2 15:04:05", strDate)

			// 	var last30 Last30Days
			// 	last30.Day = (i + 1)
			// 	last30.DateId = date
			// 	last30.CurrProduction = 0
			// 	last30.CurrBudget = 0
			// 	last30.CumProduction = 0
			// 	last30.CumBudget = 0

			// 	item30s = append(item30s, last30)
			// }

			// csr, e = ctx.NewQuery().From(new(ScadaData).TableName()).
			// 	Where(dbox.Eq("dateinfo.monthid", dtInfo.MonthId)).
			// 	Order("dateinfo.dateid").
			// 	Aggr(dbox.AggrSum, "$power", "totalpower").
			// 	Aggr(dbox.AggrSum, "$powerlost", "totalpowerlost").
			// 	Aggr(dbox.AggrSum, "$oktime", "totaloktime").
			// 	Aggr(dbox.AggrSum, "$griddowntime", "totalgriddowntime").
			// 	Aggr(dbox.AggrAvr, "$windspeed", "avgwindspeed").
			// 	Group("dateinfo.dateid").
			// 	Order("id.dateinfo_dateid").
			// 	Cursor(nil)

			pipe := []tk.M{tk.M{}.Set("$match", tk.M{}.Set("dateinfo.monthid", tk.M{}.Set("$eq", dtInfo.MonthId)).Set("power", tk.M{}.Set("$gte", 0))), tk.M{}.Set("$group", tk.M{}.Set("_id", "$dateinfo.dateid").Set("totalpower", tk.M{}.Set("$sum", "$power"))), tk.M{}.Set("$sort", tk.M{}.Set("_id", 1))}
			/*csr, _ := ctx.NewQuery().
			Command("pipe", pipe).
			From(new(ScadaData).TableName()).
			Cursor(nil)*/
			csr, _ := ctx.NewQuery().
				Command("pipe", pipe).
				From(new(ScadaData).TableName()).
				Cursor(nil)
			defer csr.Close()

			scadas := []tk.M{}
			e = csr.Fetch(&scadas, 0, false)

			item30s := make([]Last30Days, 0)
			dateData := dateId
			cummProd := 0.0
			cummBudget := 0.0
			for _, data := range scadas {
				dateData = data["_id"].(time.Time)
				// dateData = id["dateinfo_dateid"].(time.Time)

				//tk.Printf("#%v\n", dateData)

				var last30 Last30Days
				last30.DateId = dateData
				last30.DayNo = dateData.Day()

				currProd := 0.0
				currBudget := budgetCurrMonthDaily // 565160.32
				if data != nil {
					ipower := data["totalpower"]
					power := 0.0
					if ipower != nil {
						power = ipower.(float64)
					}
					currProd = power / 6
				}
				cummProd = cummProd + currProd
				cummBudget = cummBudget + currBudget

				last30.CurrBudget = currBudget
				last30.CurrProduction = currProd
				last30.CumBudget = cummBudget / 1000000
				last30.CumProduction = cummProd / 1000000

				item30s = append(item30s, last30)

				dateData = dateId.Add(-1)
			}

			mdl.Productions = items
			mdl.CummulativeProductions = item30s

			d.BaseController.Ctx.Insert(mdl)

			mdl = new(ScadaLastUpdate).New()
			mdl.ID = "SCADALASTUPDATE_TEJUVA"
			mdl.LastUpdate = maxTimeStamp
			mdl.DateInfo = dtInfo
			mdl.NoOfProjects = 1
			mdl.NoOfTurbines = 24
			mdl.TotalMaxCapacity = 24 * 2100
			mdl.CurrentDown = 0
			mdl.TwoDaysDown = 0
			mdl.ProjectName = "Tejuva"

			// csr, e = ctx.NewQuery().From(new(ScadaData).TableName()).
			// 	Where(dbox.And(dbox.Gt("timestamp", startTime), dbox.Lte("timestamp", maxTimeStamp))).
			// 	Aggr(dbox.AggrSum, "$power", "totalpower").
			// 	Aggr(dbox.AggrSum, "$powerlost", "totalpowerlost").
			// 	Aggr(dbox.AggrSum, "$oktime", "totaloktime").
			// 	Aggr(dbox.AggrSum, "$griddowntime", "totalgriddowntime").
			// 	Aggr(dbox.AggrAvr, "$windspeed", "avgwindspeed").
			// 	Group("timestamp").
			// 	Cursor(nil)
			// defer csr.Close()

			// scadas := []tk.M{}
			// e = csr.Fetch(&scadas, 0, false)

			// items = make([]LastData24Hours, 0)
			// for i := 0; i < 24; i++ {
			// 	year := strconv.Itoa(dateId.Year())
			// 	month := dateId.Month().String()
			// 	day := strconv.Itoa(dateId.Day())
			// 	strTime := year + "-" + month + "-" + day + " " + strconv.Itoa(i) + ":00:00"
			// 	timeHr, _ := time.Parse("2006-January-2 15:04:05", strTime)

			// 	var last LastData24Hours
			// 	last.Hour = i
			// 	last.TimeHour = timeHr
			// 	last.PowerKw = 0.0
			// 	last.EnergyKwh = 0.0
			// 	last.Potential = 0.0
			// 	last.TrueAvail = 0.0
			// 	last.GridAvail = 0.0
			// 	last.AvgWindSpeed = 0.0

			// 	items = append(items, last)
			// }

			// item30s = make([]Last30Days, 0)
			// for i := 0; i < 30; i++ {
			// 	year := strconv.Itoa(dateId.Year())
			// 	month := dateId.Month().String()
			// 	day := strconv.Itoa(i + 1)
			// 	strDate := year + "-" + month + "-" + day + " 00:00:00"
			// 	date, _ := time.Parse("2006-January-2 15:04:05", strDate)

			// 	var last30 Last30Days
			// 	last30.Day = (i + 1)
			// 	last30.DateId = date
			// 	last30.CurrProduction = 0
			// 	last30.CurrBudget = 0
			// 	last30.CumProduction = 0
			// 	last30.CumBudget = 0

			// 	item30s = append(item30s, last30)
			// }

			mdl.Productions = items
			mdl.CummulativeProductions = item30s

			d.BaseController.Ctx.Insert(mdl)
		}
	}
}
