package generatorControllers

import (
	. "eaciit/wfdemo-git/library/helper"
	. "eaciit/wfdemo-git/library/models"
	. "eaciit/wfdemo-git/processapp/summaryGenerator/controllers"
	"eaciit/wfdemo-git/web/helper"
	_ "fmt"
	"log"
	_ "strings"
	"time"

	"strings"

	"github.com/eaciit/dbox"
	_ "github.com/eaciit/dbox/dbc/mongo"
	tk "github.com/eaciit/toolkit"
)

type GenScadaLast24 struct {
	*BaseController
}

func (d *GenScadaLast24) Generate(base *BaseController) {
	if base != nil {
		// d.BaseController = base
		// ctx, e := PrepareConnection()
		// if e != nil {
		// 	ErrorHandler(e, "Scada Summary")
		// 	os.Exit(0)
		// }

		t0 := time.Now()
		tk.Println("Start generating data last 24")

		ctx := d.BaseController.Ctx.Connection
		// d.BaseController.Ctx.DeleteMany(new(ScadaLastUpdate), dbox.And(dbox.Ne("_id", "")))

		projectList, _ := helper.GetProjectList()

		inprojectactive := func(str string) bool {
			for _, v := range projectList {
				if v.Value == str {
					return true
				}
			}
			return false
		}

		mapbudget := map[string]float64{}
		csrBudget, _ := ctx.NewQuery().From(new(ExpPValueModel).TableName()).
			Cursor(nil)

		budgets := make([]ExpPValueModel, 0)
		_ = csrBudget.Fetch(&budgets, 0, false)
		csrBudget.Close()

		for _, budget := range budgets {
			mapbudget[tk.Sprintf("%s_%d_75", budget.ProjectName, budget.MonthNo)] = budget.P75NetGenMWH
			mapbudget[tk.Sprintf("%s_%d_50", budget.ProjectName, budget.MonthNo)] = budget.P50NetGenMWH
			mapbudget[tk.Sprintf("%s_%d_90", budget.ProjectName, budget.MonthNo)] = budget.P90NetGenMWH
			if inprojectactive(budget.ProjectName) {
				mapbudget[tk.Sprintf("fleet_%d_75", budget.MonthNo)] = budget.P75NetGenMWH
				mapbudget[tk.Sprintf("fleet_%d_50", budget.MonthNo)] = budget.P50NetGenMWH
				mapbudget[tk.Sprintf("fleet_%d_90", budget.MonthNo)] = budget.P90NetGenMWH
			}
		}

		for _, proj := range d.BaseController.ProjectList {

			tk.Println("Start : ", proj.Name)

			projectName := proj.Value
			turbineList := []TurbineOut{}

			if projectName != "Fleet" {
				turbineList, _ = helper.GetTurbineList([]interface{}{projectName})
			} else {
				turbineList, _ = helper.GetTurbineList(nil)
			}

			totalTurbine := len(turbineList)

			filter := dbox.Eq("available", 1)
			if projectName != "Fleet" {
				filter = dbox.And(dbox.Eq("projectname", projectName), filter)
			}

			/*for _, v := range filter {
				log.Printf(">> %#v \n", v)
			}*/

			csr, e := ctx.NewQuery().
				From(new(ScadaData).TableName()).
				Where(filter).
				Aggr(dbox.AggrMax, "$timestamp", "timestamp").
				Aggr(dbox.AggrMax, "$dateinfo.dateid", "dateid").
				Group("").
				Cursor(nil)

			if e != nil {
				log.Printf("Error: %v \n", e.Error())
			} else {
				datas := []tk.M{}
				e = csr.Fetch(&datas, 0, false)
				csr.Close()

				tk.Printf(">> %#v \n", datas)

				if len(datas) > 0 {
					dateId := datas[0].Get("dateid", time.Time{}).(time.Time).UTC()
					dtInfo := GetDateInfo(dateId)
					maxTimeStamp := datas[0].Get("timestamp", time.Time{}).(time.Time).UTC()

					var budgetCurrMonthDaily float64
					var budgetCurrMonthDaily50 float64
					var budgetCurrMonthDaily90 float64

					_id := tk.Sprintf("%s_%d", projectName, dateId.Month())
					if val, cond := mapbudget[_id+"_75"]; cond {
						budgetCurrMonths := val * 1000.0
						noOfDay := float64(daysIn(dateId.Month(), dateId.Year()))
						budgetCurrMonthDaily = tk.Div(budgetCurrMonths, noOfDay)
					}
					if val, cond := mapbudget[_id+"_50"]; cond {
						budgetCurrMonths := val * 1000.0
						noOfDay := float64(daysIn(dateId.Month(), dateId.Year()))
						budgetCurrMonthDaily50 = tk.Div(budgetCurrMonths, noOfDay)
					}
					if val, cond := mapbudget[_id+"_90"]; cond {
						budgetCurrMonths := val * 1000.0
						noOfDay := float64(daysIn(dateId.Month(), dateId.Year()))
						budgetCurrMonthDaily90 = tk.Div(budgetCurrMonths, noOfDay)
					}

					mdl := new(ScadaLastUpdate).New()

					if projectName != "Fleet" {
						mdl.ID = "SCADALASTUPDATE_" + strings.ToUpper(projectName)
						mdl.ProjectName = projectName
						mdl.NoOfProjects = 1
					} else {
						mdl.ID = "SCADALASTUPDATE_FLEET"
						mdl.ProjectName = "Fleet"
						mdl.NoOfProjects = len(d.BaseController.ProjectList) - 1
					}

					for _, t := range turbineList {
						mdl.TotalMaxCapacity += t.Capacity
					}

					mdl.TotalMaxCapacity = tk.ToFloat64(mdl.TotalMaxCapacity*1000.0, 2, tk.RoundingAuto)
					mdl.LastUpdate = maxTimeStamp
					mdl.DateInfo = dtInfo
					mdl.NoOfTurbines = totalTurbine

					items := make([]LastData24Hours, 0)
					cdatehour := dateId.UTC().Add(-1 * time.Hour)
					for i := 0; i < 24; i++ {
						cdatehour = cdatehour.Add(1 * time.Hour)

						// year := strconv.Itoa(dateId.Year())
						// month := dateId.Month().String()
						// day := strconv.Itoa(dateId.Day())
						// strTime := year + "-" + month + "-" + day + " " + strconv.Itoa(i) + ":00:00"
						// timeHr, _ := time.Parse("2006-January-2 15:04:05", strTime)

						// timeHrStart := timeHr.Add(-1 * time.Hour)

						filterSub := []*dbox.Filter{}
						filterSub = append(filterSub, dbox.Gt("timestamp", cdatehour.Add(time.Hour*-1)))
						filterSub = append(filterSub, dbox.Lte("timestamp", cdatehour))
						filterSub = append(filterSub, dbox.Eq("available", 1))

						if projectName != "Fleet" {
							filterSub = append(filterSub, dbox.Eq("projectname", projectName))
						}

						csr, e = ctx.NewQuery().From(new(ScadaData).TableName()).
							Where(dbox.And(filterSub...)).
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

						var last LastData24Hours
						if len(scadas) > 0 {
							data := scadas[0]
							trueAvail := 0.0
							gridAvail := 0.0

							ipower := data["totalpower"]
							power := 0.0
							if ipower != nil {
								power = tk.ToFloat64(ipower, 6, tk.RoundingAuto)
							}

							ipotentialpower := data["denpower"]
							potentialpower := 0.0
							if ipotentialpower != nil {
								potentialpower = tk.ToFloat64(ipotentialpower, 6, tk.RoundingAuto)
							}

							iwindspeed := data["avgwindspeed"]
							windspeed := 0.0
							if iwindspeed != nil {
								windspeed = tk.ToFloat64(iwindspeed, 6, tk.RoundingAuto)
							}
							last.Hour = i
							last.TimeHour = cdatehour
							last.AvgWindSpeed = windspeed
							last.PowerKw = power
							last.EnergyKwh = power / 6
							last.Potential = potentialpower
							last.PotentialKwh = potentialpower / 6
							last.TrueAvail = trueAvail
							last.GridAvail = gridAvail
						} else {
							last.Hour = i
							last.TimeHour = cdatehour
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

					match := tk.M{}

					match.Set("dateinfo.monthid", tk.M{}.Set("$eq", dtInfo.MonthId)).Set("available", tk.M{}.Set("$eq", 1))

					if projectName != "Fleet" {
						match.Set("projectname", projectName)
					}

					pipe := []tk.M{tk.M{}.Set("$match", match), tk.M{}.Set("$group", tk.M{}.Set("_id", "$dateinfo.dateid").Set("totalpower", tk.M{}.Set("$sum", "$power"))), tk.M{}.Set("$sort", tk.M{}.Set("_id", 1))}

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
					cummBudget50 := 0.0
					cummBudget90 := 0.0
					for _, data := range scadas {
						dateData = data["_id"].(time.Time)
						var last30 Last30Days
						last30.DateId = dateData
						last30.DayNo = dateData.Day()

						currProd := 0.0
						currBudget := budgetCurrMonthDaily // 565160.32
						currBudget50 := budgetCurrMonthDaily50
						currBudget90 := budgetCurrMonthDaily90
						if data != nil {
							ipower := data["totalpower"]
							power := 0.0
							if ipower != nil {
								power = data.GetFloat64("totalpower")
							}
							currProd = power / 6
						}
						cummProd = cummProd + currProd
						cummBudget = cummBudget + currBudget
						cummBudget50 += currBudget50
						cummBudget90 += currBudget90

						last30.CurrBudget = currBudget
						last30.CurrBudget50 = currBudget50
						last30.CurrBudget90 = currBudget90
						last30.CurrProduction = currProd
						last30.CumBudget = cummBudget / 1000000
						last30.CumBudget50 = cummBudget50 / 1000000
						last30.CumBudget90 = cummBudget90 / 1000000
						last30.CumProduction = cummProd / 1000000

						item30s = append(item30s, last30)

						dateData = dateId.Add(-1)
					}

					mdl.Productions = items
					mdl.CummulativeProductions = item30s

					d.BaseController.Ctx.Save(mdl)
				}
			}
		}

		tk.Println("End generating data last 24 in ", time.Since(t0).String())
	}
}
