package generatorControllers

import (
	. "eaciit/wfdemo-git/library/helper"
	. "eaciit/wfdemo-git/library/models"
	. "eaciit/wfdemo-git/processapp/summaryGenerator/controllers"
	"eaciit/wfdemo-git/web/helper"
	"math"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/eaciit/dbox"
	_ "github.com/eaciit/dbox/dbc/mongo"
	tk "github.com/eaciit/toolkit"
)

const (
	sError   = "ERROR"
	sInfo    = "INFO"
	sWarning = "WARNING"
)

type GenScadaSummary struct {
	*BaseController
	sync.RWMutex
}

func daysIn(m time.Month, year int) int {
	return time.Date(year, m+1, 0, 0, 0, 0, 0, time.UTC).Day()
}

func (d *GenScadaSummary) Generate(base *BaseController) {
	if base != nil {
		d.BaseController = base

		ctx, e := PrepareConnection()
		if e != nil {
			ErrorHandler(e, "Scada Summary")
			os.Exit(0)
		}

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

		d.BaseController.Ctx.DeleteMany(new(ScadaSummaryByMonth), dbox.Ne("projectname", ""))

		for _, v := range d.BaseController.ProjectList {
			project := v.Value

			filter := []*dbox.Filter{}
			//
			// filter = append(filter, dbox.Gte("power", -200))
			// using same filter for all generation function
			// @asp 21-07-2017
			filter = append(filter, dbox.Eq("available", 1))

			group := []string{}

			if project != "Fleet" {
				filter = append(filter, dbox.Eq("projectname", project))
				group = []string{"projectname", "dateinfo.monthid"}
			} else {
				group = []string{"dateinfo.monthid"}
			}

			csr, e := ctx.NewQuery().From(new(ScadaData).TableName()).
				Where(dbox.And(filter...)).
				Aggr(dbox.AggrSum, "$power", "totalpower").
				Aggr(dbox.AggrSum, "$energy", "energy").
				Aggr(dbox.AggrSum, "$energylost", "totalenergylost").
				Aggr(dbox.AggrSum, "$oktime", "totaloktime").
				Aggr(dbox.AggrSum, "$minutes", "totalminutes").
				Aggr(dbox.AggrSum, "$griddowntime", "totalgriddowntime").
				Aggr(dbox.AggrSum, "$unknowntime", "totalunknowntime").
				Aggr(dbox.AggrSum, "$machinedowntime", "totalmachinedowntime").
				Aggr(dbox.AggrAvr, "$avgwindspeed", "avgwindspeed").
				Aggr(dbox.AggrSum, 1, "totaltimestamp").
				Aggr(dbox.AggrMax, "$dateinfo.dateid", "maxdateid").
				Aggr(dbox.AggrMin, "$dateinfo.dateid", "mindateid").
				Group(group...).
				Cursor(nil)
			defer csr.Close()

			if e != nil {
				ErrorHandler(e, "Scada Summary")
				os.Exit(0)
			}

			datas := []tk.M{}
			e = csr.Fetch(&datas, 0, false)

			divider := 1000.0
			noOfTurbine := 0.0
			plfDivider := 0.0

			for _, data := range datas {
				id := data["_id"].(tk.M)
				imonthid := id["dateinfo_monthid"].(int)
				monthid := strconv.Itoa(imonthid)
				year := monthid[0:4]
				month := monthid[4:6]
				day := "01"

				iMonth, _ := strconv.Atoi(string(month))
				iMonth = iMonth - 1

				dtStr := year + "-" + month + "-" + day
				dtId, _ := time.Parse("2006-01-02", dtStr)
				dtinfo := GetDateInfo(dtId)
				ioktime := data["totaloktime"]
				oktime := 0.0
				if ioktime != nil {
					oktime = (ioktime.(float64)) / 3600 // divide by 3600 secs, result in hours
				}
				noOfTurbine = d.BaseController.TotalTurbinePerMonth[project+"_"+monthid]
				plfDivider = d.BaseController.CapacityPerMonth[project+"_"+monthid]

				revenueTimes := 5.74

				duration := 0.0
				lostEnergy := 0.0

				ipower := data["totalpower"]
				power := 0.0
				if ipower != nil {
					power = ipower.(float64)
				}

				ienergy := data["energy"]
				energy := 0.0
				if ienergy != nil {
					energy = ienergy.(float64)
				}

				iminutes := data["totalminutes"]
				minutes := 0.0
				if iminutes != nil {
					minutes = float64(iminutes.(int)) / 60
				}

				imaxdate := data["maxdateid"]
				imindate := data["mindateid"]
				maxdate := time.Now()
				mindate := time.Now()
				if imaxdate != nil {
					maxdate = imaxdate.(time.Time)
				}

				if imindate != nil {
					mindate = imindate.(time.Time)
				}
				//log.Printf("#%v\n", (power / 6000000))

				pipe := []tk.M{tk.M{}.Set("$unwind", "$detail"),
					tk.M{}.Set("$match", tk.M{}.Set("detail.detaildateinfo.monthid", imonthid)), tk.M{}.Set("$group", tk.M{}.Set("_id", "$projectname").Set("duration", tk.M{}.Set("$sum", "$detail.duration")).Set("powerlost", tk.M{}.Set("$sum", "$detail.powerlost")))}
				// log.Printf("#%v\n", pipe)
				csr1, _ := ctx.NewQuery().
					Command("pipe", pipe).
					From(new(Alarm).TableName()).
					Cursor(nil)
				defer csr1.Close()

				alarms := []tk.M{}
				e = csr1.Fetch(&alarms, 0, false)

				// log.Printf("#%v\n", alarms)

				if len(alarms) > 0 {
					alarm := alarms[0]
					ipowerlost := alarm["powerlost"]
					if ipowerlost != nil {
						lostEnergy = ipowerlost.(float64)
					}
					iduration := alarm["duration"]
					if iduration != nil {
						duration = iduration.(float64)
					}
				}
				//log.Println("Lost:", lostEnergy)
				//log.Println("Duration:", duration)

				powerlastyear := 0.0
				// powerlost := data["totalpowerlost"].(float64)
				// griddowntime := data["totalgriddowntime"].(float64)
				iwindspeed := data["avgwindspeed"]
				windspeed := 0.0
				if iwindspeed != nil {
					windspeed = iwindspeed.(float64)
				}

				itotaldata := data["totaltimestamp"]
				totaldata := 0
				if itotaldata != nil {
					totaldata = itotaldata.(int)
				}

				igriddowntime := data["totalgriddowntime"]
				griddowntime := 0.0
				if igriddowntime != nil {
					griddowntime = igriddowntime.(float64) / 3600
				}

				imachinedowntime := data["totalmachinedowntime"]
				machinedowntime := 0.0
				if imachinedowntime != nil {
					machinedowntime = imachinedowntime.(float64) / 3600
				}

				iunknowntime := data["totalunknowntime"]
				unknowntime := 0.0
				if iunknowntime != nil {
					unknowntime = iunknowntime.(float64) / 3600
				}
				_ = machinedowntime
				_ = griddowntime
				_ = unknowntime

				expwstimes := 0.133
				randno := tk.RandInt(5)
				//log.Printf("#%v\n", randno)
				if randno > 3 {
					expwstimes = -0.125
				}

				expWindSpeed := (windspeed + (windspeed * expwstimes))
				revenue := power * revenueTimes
				revenueInLacs := revenue / 100000

				tStart, _ := time.Parse("020601_150405", "01"+maxdate.UTC().Format("0601")+"_000000")
				daysInYear := GetDayInYear(tStart.Year())
				days := daysInYear.GetInt(tk.ToString(int(tStart.Month())))

				tEnd, _ := time.Parse("020601_150405", tk.ToString(days)+maxdate.UTC().Format("0601")+"_235959")

				hourValue := helper.GetHourValue(tStart.UTC(), tEnd.UTC(), mindate.UTC(), maxdate.UTC())

				machineAvail, gridAvail, scadaAvail, trueAvail, plf := helper.GetAvailAndPLF(noOfTurbine, oktime*3600, energy/1000, machinedowntime, griddowntime, float64(totaldata), hourValue, minutes, plfDivider)

				if plf > 100 {
					plf = 100
				}
				if gridAvail > 100 {
					gridAvail = 100
				}
				if scadaAvail > 100 {
					scadaAvail = 100
				}
				if machineAvail > 100 {
					machineAvail = 100
				}

				budget := budgetMonths[iMonth]

				mdl := new(ScadaSummaryByMonth).New()
				mdl.ProjectName = project
				mdl.DateInfo = dtinfo
				mdl.Production = ((power / 6) / divider)
				mdl.ProductionLastYear = (powerlastyear / divider)
				mdl.Revenue = revenue
				mdl.RevenueInLacs = revenueInLacs
				mdl.TrueAvail = trueAvail
				mdl.ScadaAvail = scadaAvail
				mdl.MachineAvail = machineAvail
				mdl.GridAvail = gridAvail
				mdl.PLF = plf
				mdl.Budget = budget * 1000
				mdl.AvgWindSpeed = windspeed
				mdl.ExpWindSpeed = expWindSpeed
				mdl.DowntimeHours = duration
				mdl.LostEnergy = lostEnergy / (divider * divider) // convert to giga
				mdl.RevenueLoss = (lostEnergy * revenueTimes)

				if mdl != nil {
					d.BaseController.Ctx.Insert(mdl)
				}

			}

		}

	}
}

func (d *GenScadaSummary) GenerateSummaryByProject(base *BaseController) {
	if base != nil {

		d.BaseController = base

		ctx, e := PrepareConnection()
		if e != nil {
			ErrorHandler(e, "Scada Summary")
			os.Exit(0)
		}

		d.BaseController.Ctx.DeleteMany(new(ScadaSummaryByProject), dbox.Ne("_id", ""))

		for _, v := range d.BaseController.ProjectList {
			var turbineList []TurbineOut
			projectName := v.Value
			group := "projectname"

			filter := []*dbox.Filter{}
			// filter = append(filter, dbox.Gte("power", -200))
			// using same filter for all generation function
			// @asp 21-07-2017
			filter = append(filter, dbox.Eq("available", 1))

			var max time.Time

			if projectName != "Fleet" {
				filter = append(filter, dbox.Eq("projectname", projectName))
				group = "turbine"
				turbineList, _ = helper.GetTurbineList([]interface{}{projectName})
				_, max, _ = GetDataDateAvailable(new(ScadaData).TableName(), "dateinfo.dateid", dbox.Eq("projectname", projectName), d.Ctx.Connection)
			}

			csr, e := ctx.NewQuery().From(new(ScadaData).TableName()).
				Where(dbox.And(filter...)).
				Aggr(dbox.AggrSum, "$power", "totalpower").
				Aggr(dbox.AggrSum, "$energy", "energy").
				Aggr(dbox.AggrSum, "$energylost", "totalenergylost").
				Aggr(dbox.AggrSum, "$oktime", "totaloktime").
				Aggr(dbox.AggrSum, "$minutes", "totalminutes").
				Aggr(dbox.AggrSum, "$griddowntime", "totalgriddowntime").
				Aggr(dbox.AggrSum, "$unknowntime", "totalunknowntime").
				Aggr(dbox.AggrSum, "$machinedowntime", "totalmachinedowntime").
				Aggr(dbox.AggrAvr, "$avgwindspeed", "avgwindspeed").
				Aggr(dbox.AggrMax, "$dateinfo.dateid", "max").
				Aggr(dbox.AggrMin, "$dateinfo.dateid", "min").
				Group(group).
				Cursor(nil)
			defer csr.Close()

			_ = e

			datas := []tk.M{}
			e = csr.Fetch(&datas, 0, false)

			// daysInMonth := GetDayInYear(max.Year())
			// days := tk.ToString(daysInMonth.GetInt(tk.ToString(int(max.Month()))))
			// tmpdt, _ := time.Parse("060102_150405", max.UTC().Format("0601")+days+"_000000")

			mdl := new(ScadaSummaryByProject).New()
			mdl.ID = projectName

			items := make([]ScadaSummaryByProjectItem, 0)
			for _, data := range datas {
				id := data["_id"].(tk.M)
				turbine := id[group].(string)
				idfield := "turbine"
				totalWtg := 0

				if projectName == "Fleet" {
					turbineList, _ = helper.GetTurbineList([]interface{}{turbine})
					idfield = "farm"
					_, max, _ = GetDataDateAvailable(new(ScadaData).TableName(), "dateinfo.dateid", dbox.Eq("projectname", turbine), d.Ctx.Connection)
				}
				tmpdt := max.UTC()
				endDate := tmpdt.UTC() //time.Parse("060102_150405", max.UTC().Format("0601")+"01_000000").UTC()
				startDate := GetNormalAddDateMonth(tmpdt.UTC(), -11)

				match := tk.M{}.
					Set(idfield, tk.M{}.Set("$eq", turbine)).
					Set("detail.detaildateinfo.dateid", tk.M{}.Set("$gte", startDate).Set("$lte", endDate))

				ioktime := data["totaloktime"]
				oktime := 0.0
				if ioktime != nil {
					oktime = (ioktime.(float64)) /// 3600 // divide by 3600 secs, result in hours
				}

				iminutes := data["totalminutes"]
				minutes := 0.0
				if iminutes != nil {
					valminutes, ok := iminutes.(float64) // divide by 60 secs, result in hours
					if !ok {
						ivalminutes := iminutes.(int)
						minutes = float64(ivalminutes) / 60.0
					} else {
						minutes = valminutes / 60.0
					}
				}
				_ = minutes

				downtimeHours := 0.0
				lostEnergy := 0.0

				// ilostEnergy := data["totalenergylost"]
				// if ilostEnergy != nil {
				// 	lostEnergy = ilostEnergy.(float64)
				// }

				pipe := []tk.M{
					tk.M{}.Set("$unwind", "$detail"),
					tk.M{}.Set("$match", match),
					tk.M{}.Set("$group", tk.M{}.Set("_id", "$"+idfield).
						Set("duration", tk.M{}.Set("$sum", "$detail.duration")).
						Set("powerlost", tk.M{}.Set("$sum", "$detail.powerlost"))),
					tk.M{}.Set("$sort", tk.M{}.Set("_id", 1)),
				}
				// log.Printf("#%v\n", pipe)
				csr1, _ := ctx.NewQuery().
					Command("pipe", pipe).
					From(new(Alarm).TableName()).
					Cursor(nil)
				defer csr1.Close()

				alarms := []tk.M{}
				e = csr1.Fetch(&alarms, 0, false)

				// log.Printf("#%v\n", alarms)

				if len(alarms) > 0 {
					alarm := alarms[0]
					ipowerlost := alarm["powerlost"]
					if ipowerlost != nil {
						lostEnergy = ipowerlost.(float64)
					}
					iduration := alarm["duration"]
					if iduration != nil {
						downtimeHours = iduration.(float64)
					}
				}

				ipower := data["totalpower"]
				power := 0.0
				if ipower != nil {
					power = ipower.(float64)
				}

				imachinedowntime := data["totalmachinedowntime"]
				machinedowntime := 0.0
				if imachinedowntime != nil {
					machinedowntime = imachinedowntime.(float64) / 3600
				}

				igriddowntime := data["totalgriddowntime"]
				griddowntime := 0.0
				if igriddowntime != nil {
					griddowntime = igriddowntime.(float64) / 3600
				}

				iunknowntime := data["totalunknowntime"]
				unknowntime := 0.0
				if iunknowntime != nil {
					unknowntime = iunknowntime.(float64) / 3600
				}

				_ = machinedowntime
				_ = griddowntime
				_ = unknowntime

				//downtimeHours = machinedowntime + griddowntime + unknowntime

				// machineAvail := (minutes - machinedowntime) / (durationInMonth * 24)

				maxDate := data.Get("max", time.Time{}).(time.Time)
				minDate := data.Get("min", time.Time{}).(time.Time)
				energy := data.GetFloat64("energy") / 1000

				hourValue := helper.GetHourValue(startDate.UTC(), endDate.UTC(), minDate.UTC(), maxDate.UTC())

				// log.Printf("%v | %v | %v \n", hourValue, minDate.String(), maxDate.String())
				var plfDivider float64
				for _, v := range turbineList {
					if projectName != "Fleet" {
						if v.Value == turbine {
							plfDivider += v.Capacity
							totalWtg += 1
						}
					} else {
						if v.Project == turbine {
							plfDivider += v.Capacity
							totalWtg += 1
						}
					}
				}

				machineAvail, _, _, trueAvail, plf := helper.GetAvailAndPLF(float64(totalWtg), oktime, energy, machinedowntime, griddowntime, float64(0), hourValue, minutes, plfDivider)

				var item ScadaSummaryByProjectItem

				item.Name = turbine
				item.NoOfWtg = totalWtg
				item.Production = power / 6
				item.PLF = plf / 100 //(power / 6) / (durationInMonth * 24 * 2100)
				item.MachineAvail = machineAvail / 100
				item.TrueAvail = trueAvail / 100 //oktime / (durationInMonth * 24)
				item.LostEnergy = lostEnergy
				item.DowntimeHours = downtimeHours

				items = append(items, item)
			}

			mdl.DataItems = items

			d.BaseController.Ctx.Insert(mdl)
		}
	}
}

func (d *GenScadaSummary) GenerateSummaryByFleet(base *BaseController) {
	if base != nil {
		d.BaseController = base

		ctx, e := PrepareConnection()
		if e != nil {
			ErrorHandler(e, "Scada Summary")
			os.Exit(0)
		}

		d.BaseController.Ctx.DeleteMany(new(ScadaSummaryByProject), dbox.Ne("_id", "Fleet"))

		csr, e := ctx.NewQuery().From(new(ScadaData).TableName()).
			// using same filter for all generation function
			// @asp 21-07-2017
			Where(dbox.Eq("available", 1)).
			// Where(dbox.Gte("power", -200)).
			Aggr(dbox.AggrSum, "$power", "totalpower").
			Aggr(dbox.AggrSum, "$energy", "energy").
			Aggr(dbox.AggrSum, "$energylost", "totalenergylost").
			Aggr(dbox.AggrSum, "$oktime", "totaloktime").
			Aggr(dbox.AggrSum, "$minutes", "totalminutes").
			Aggr(dbox.AggrSum, "$griddowntime", "totalgriddowntime").
			Aggr(dbox.AggrSum, "$unknowntime", "totalunknowntime").
			Aggr(dbox.AggrSum, "$machinedowntime", "totalmachinedowntime").
			Aggr(dbox.AggrAvr, "$avgwindspeed", "avgwindspeed").
			Aggr(dbox.AggrMax, "$dateinfo.dateid", "max").
			Aggr(dbox.AggrMin, "$dateinfo.dateid", "min").
			Group("projectname").
			Cursor(nil)
		defer csr.Close()

		datas := []tk.M{}
		e = csr.Fetch(&datas, 0, false)

		_, max, _ := GetDataDateAvailable(new(ScadaData).TableName(), "dateinfo.dateid", nil, d.Ctx.Connection)

		daysInMonth := GetDayInYear(max.Year())
		days := tk.ToString(daysInMonth.GetInt(tk.ToString(int(max.Month()))))
		tmpdt, _ := time.Parse("060102_150405", max.UTC().Format("0601")+days+"_235959")
		endDate := tmpdt.UTC() //time.Parse("060102_150405", max.UTC().Format("0601")+"01_000000").UTC()

		startDate := GetNormalAddDateMonth(tmpdt.UTC(), -11)

		//log.Println(durationInMonth)

		noOfTurbine := 0
		turbineList, _ := helper.GetTurbineList(nil)
		noOfTurbine = len(turbineList)

		mdl := new(ScadaSummaryByProject).New()
		mdl.ID = "Fleet"

		items := make([]ScadaSummaryByProjectItem, 0)
		for _, data := range datas {
			id := data["_id"].(tk.M)
			name := id["projectname"].(string)

			ioktime := data["totaloktime"]
			oktime := 0.0
			if ioktime != nil {
				oktime = (ioktime.(float64)) / 3600 // divide by 3600 secs, result in hours
			}

			iminutes := data["totalminutes"]
			minutes := 0.0
			if iminutes != nil {
				valminutes, ok := iminutes.(float64) // divide by 60 secs, result in hours
				if !ok {
					ivalminutes := iminutes.(int)
					minutes = float64(ivalminutes) / 60.0
				} else {
					minutes = valminutes / 60.0
				}
			}
			_ = minutes

			downtimeHours := 0.0
			lostEnergy := 0.0

			ilostEnergy := data["totalenergylost"]
			if ilostEnergy != nil {
				lostEnergy = ilostEnergy.(float64)
			}

			// pipe := []tk.M{tk.M{}.Set("$match", tk.M{}.Set("projectname", tk.M{}.Set("$eq", name)).Set("startdateinfo.dateid", tk.M{}.Set("$lte", endDate))), tk.M{}.Set("$group", tk.M{}.Set("_id", "$projectname").Set("duration", tk.M{}.Set("$sum", "$duration")).Set("powerlost", tk.M{}.Set("$sum", "$powerlost"))), tk.M{}.Set("$sort", tk.M{}.Set("_id", 1))}
			// csr1, _ := ctx.NewQuery().
			// 	Command("pipe", pipe).
			// 	From(new(Alarm).TableName()).
			// 	Cursor(nil)
			// defer csr1.Close()

			// alarms := []tk.M{}
			// e = csr1.Fetch(&alarms, 0, false)

			// if len(alarms) > 0 {
			// 	alarm := alarms[0]
			// 	iduration := alarm["duration"]
			// 	if iduration != nil {
			// 		downtimeHours = iduration.(float64)
			// 	}
			// 	ipowerlost := alarm["powerlost"]
			// 	if ipowerlost != nil {
			// 		lostEnergy = ipowerlost.(float64)
			// 	}
			// }

			pipe := []tk.M{tk.M{}.Set("$unwind", "$detail"),
				tk.M{}.Set("$match", tk.M{}.Set("detail.detaildateinfo.dateid", tk.M{}.Set("$lte", endDate))), tk.M{}.Set("$group", tk.M{}.Set("_id", "$projectname").Set("duration", tk.M{}.Set("$sum", "$detail.duration")).Set("powerlost", tk.M{}.Set("$sum", "$detail.powerlost"))), tk.M{}.Set("$sort", tk.M{}.Set("_id", 1))}
			// log.Printf("#%v\n", pipe)
			csr1, _ := ctx.NewQuery().
				Command("pipe", pipe).
				From(new(Alarm).TableName()).
				Cursor(nil)
			defer csr1.Close()

			alarms := []tk.M{}
			e = csr1.Fetch(&alarms, 0, false)

			// log.Printf("#%v\n", alarms)

			if len(alarms) > 0 {
				alarm := alarms[0]
				ipowerlost := alarm["powerlost"]
				if ipowerlost != nil {
					lostEnergy = ipowerlost.(float64)
				}
				iduration := alarm["duration"]
				if iduration != nil {
					downtimeHours = iduration.(float64)
				}
			}

			ipower := data["totalpower"]
			power := 0.0
			if ipower != nil {
				power = ipower.(float64)
			}

			imachinedowntime := data["totalmachinedowntime"]
			machinedowntime := 0.0
			if imachinedowntime != nil {
				machinedowntime = imachinedowntime.(float64) / 3600
			}

			igriddowntime := data["totalgriddowntime"]
			griddowntime := 0.0
			if igriddowntime != nil {
				griddowntime = igriddowntime.(float64) / 3600
			}

			iunknowntime := data["totalunknowntime"]
			unknowntime := 0.0
			if iunknowntime != nil {
				unknowntime = iunknowntime.(float64) / 3600
			}

			_ = machinedowntime + griddowntime + unknowntime

			maxDate := data.Get("max").(time.Time)
			minDate := data.Get("min").(time.Time)
			energy := data.GetFloat64("energy") / 1000

			hourValue := helper.GetHourValue(startDate.UTC(), endDate.UTC(), minDate.UTC(), maxDate.UTC())

			// log.Printf("%v | %v | %v | %v : %v  \n", startDate.UTC().String(), endDate.UTC().String(), minDate.UTC().String(), maxDate.UTC().String(), hourValue)

			var plfDivider float64
			for _, v := range turbineList {
				plfDivider += v.Capacity
			}

			machineAvail, _, _, trueAvail, plf := helper.GetAvailAndPLF(float64(noOfTurbine), oktime*3600, energy, machinedowntime, griddowntime, float64(0), hourValue, minutes, plfDivider)

			// log.Printf("%v | %v | %v | %v \n", machineAvail, trueAvail, plf, energy)

			var item ScadaSummaryByProjectItem

			item.Name = name
			item.NoOfWtg = int(noOfTurbine)
			item.Production = power / 6
			item.PLF = plf / 100 //(power / 6) / (durationInMonth * 24 * 2100 * noOfTurbine)
			item.MachineAvail = machineAvail / 100
			item.TrueAvail = trueAvail / 100 //oktime / (durationInMonth * 24 * noOfTurbine)
			item.LostEnergy = lostEnergy
			item.DowntimeHours = downtimeHours

			items = append(items, item)
		}

		mdl.DataItems = items

		d.BaseController.Ctx.Insert(mdl)
	}
}

func (d *GenScadaSummary) GenerateSummaryDaily(base *BaseController) {
	if base != nil {
		d.BaseController = base

		ctx, e := PrepareConnection()
		if e != nil {
			ErrorHandler(e, "Scada Summary Daily")
			os.Exit(0)
		}

		projectList, _ := helper.GetProjectList()
		mapRevenue, arrproject := map[string]float64{}, []string{}
		for _, v := range projectList {
			mapRevenue[v.Value] = v.RevenueMultiplier
			arrproject = append(arrproject, v.ProjectId)
		}

		_t0 := time.Now().UTC()
		maxproc, cmaxproc := d.getLastExec(ctx, "summary", "daily", time.Date(_t0.Year(), _t0.Month(), _t0.Day(), 0, 0, 0, 0, _t0.Location()).AddDate(0, 0, -5), arrproject), map[string]time.Time{}

		var wg sync.WaitGroup
		counter := 0

		for turbine, v := range d.BaseController.RefTurbines {
			counter++
			wg.Add(1)

			go func(turbineX string, project string) {
				filter := tk.M{}
				filter = filter.Set("projectname", tk.M{}.Set("$eq", project))
				filter = filter.Set("turbine", tk.M{}.Set("$eq", turbineX))
				filter = filter.Set("available", 1)

				// dt := d.BaseController.GetLatest("ScadaSummaryDaily", project, turbineX)
				lep := maxproc.Get(project, new(LastExecProcess)).(*LastExecProcess)
				dt := lep.LastDate

				if dt.Format("2006") != "0001" {
					dt = dt.AddDate(0, 0, -1)
					filter = filter.Set("dateinfo.dateid", tk.M{}.Set("$gte", dt))
				}

				countws := tk.M{"$cond": tk.M{}.
					Set("if", tk.M{"$ifNull": []interface{}{"$avgwindspeed", false}}).
					Set("then", 1).
					Set("else", 0)}

				pipe := []tk.M{}
				pipe = append(pipe, tk.M{}.Set("$match", filter))
				pipe = append(pipe, tk.M{}.Set("$group", tk.M{}.
					Set("_id", tk.M{}.
						Set("projectname", "$projectname").
						Set("dateinfo", "$dateinfo")).
					Set("totaltime", tk.M{}.Set("$sum", "$totaltime")).
					Set("power", tk.M{}.Set("$sum", "$power")).
					Set("energy", tk.M{}.Set("$sum", "$energy")).
					Set("pcvalue", tk.M{}.Set("$sum", "$pcvalue")).
					Set("pcdeviation", tk.M{}.Set("$sum", "$pcdeviation")).
					Set("oktime", tk.M{}.Set("$sum", "$oktime")).
					Set("oksecs", tk.M{}.Set("$sum", "$oksecs")).
					Set("minutes", tk.M{}.Set("$sum", "$minutes")).
					Set("totalts", tk.M{}.Set("$sum", 1)).
					Set("griddowntime", tk.M{}.Set("$sum", "$griddowntime")).
					Set("machinedowntime", tk.M{}.Set("$sum", "$machinedowntime")).
					Set("avgwindspeed", tk.M{}.Set("$avg", "$avgwindspeed")).
					Set("sumwindspeed", tk.M{}.Set("$sum", "$avgwindspeed")).
					Set("countwindspeed", tk.M{}.Set("$sum", countws)).
					Set("sumlowindtime", tk.M{}.Set("$sum", "$lowindtime")).
					Set("totalrows", tk.M{}.Set("$sum", 1))))

				pipe = append(pipe, tk.M{"$sort": tk.M{"_id": 1}})

				csr, _ := ctx.NewQuery().
					Command("pipe", pipe).
					From(new(ScadaData).TableName()).
					Cursor(nil)

				scadaSums := []tk.M{}
				e = csr.Fetch(&scadaSums, 0, false)
				csr.Close()

				d.Log.AddLog(tk.Sprintf("%v | %v | %v \n", project, turbineX, len(scadaSums)), sInfo)

				revenueMultiplier := mapRevenue[project]
				revenueDividerInLacs := 100000.0
				count := 0
				total := 0

				maxdate := time.Time{}

				_ValidPCDev := PopulateValidPCDev(ctx, filter)

				for _, data := range scadaSums {
					id := data["_id"].(tk.M)
					project := id.GetString("projectname")
					// turbine := id["turbine"].(string)
					dtInfo := id.Get("dateinfo", tk.M{}).(tk.M)
					dtId := dtInfo.Get("dateid", time.Time{}).(time.Time)
					//totaltime := data["totaltime"].(float64)
					power := data.GetFloat64("power")
					energy := data.GetFloat64("energy")
					// pcvalue := data["pcvalue"].(float64)
					pcdeviation := data.GetFloat64("pcdeviation")
					oktime := data.GetFloat64("oktime")
					// oksecs := data.GetFloat64("oksecs")
					totalts := data.GetInt("totalts")
					griddowntime := data.GetFloat64("griddowntime")
					machinedowntime := data.GetFloat64("machinedowntime")
					avgwindspeed := data.GetFloat64("avgwindspeed")

					dt := new(ScadaSummaryDaily)
					dt.DateInfo = GetDateInfo(dtId)
					dt.ProjectName = project
					dt.Turbine = turbineX
					dt.PowerKw = power
					dt.Production = energy
					dt.PCDeviation = pcdeviation
					dt.Revenue = power * revenueMultiplier
					dt.RevenueInLacs = tk.Div(dt.Revenue, revenueDividerInLacs)
					dt = dt.New()

					//LackOfWind
					dt.LoWindTime = data.GetFloat64("sumlowindtime")

					//only using valid data
					dt.PCDeviation = _ValidPCDev.GetFloat64(dt.ID)

					dt.TotalRows = data.GetFloat64("totalrows")

					dt.DetWindSpeed = DetailWindSpeed{SumWindSpeed: data.GetFloat64("sumwindspeed"),
						CountWindSpeed: data.GetFloat64("countwindspeed")}

					dt.OkTime = oktime
					dt.TrueAvail = tk.Div(oktime, 144*600)
					dt.ScadaAvail = tk.Div(float64(totalts), 144.0)
					dt.TotalAvail = dt.TrueAvail

					// obsolete please see bellow the calculation using data from alarm
					// @asp 20-07-2017
					dt.MachineAvail = tk.Div(((600.0 * 144.0) - machinedowntime), 144.0*600.0)
					dt.GridAvail = tk.Div(((600.0 * 144.0) - griddowntime), 144.0*600.0)
					// ===================================================================

					turbineList, _ := helper.GetTurbineList([]interface{}{dt.ProjectName})
					capacity := 0.0

					for _, v := range turbineList {
						if turbineX == v.Value {
							capacity = v.Capacity
							break
						}
					}

					dt.PLF = tk.Div(energy, capacity*1000*24.0)
					dt.TotalMinutes = data.GetInt("minutes")

					monthNo := 0
					monthId := dtInfo.GetInt("monthid")
					sMonthNo := strconv.Itoa(monthId)[4:6]
					monthNo, _ = strconv.Atoi(sMonthNo)

					pipebudget := []tk.M{
						tk.M{}.Set("$match", tk.M{}.
							Set("projectname", project).
							Set("monthno", monthNo)),
					}

					csrBudget, _ := ctx.NewQuery().From(new(ExpPValueModel).TableName()).
						Command("pipe", pipebudget).
						// Where(dbox.And(dbox.Eq("monthno", monthNo), dbox.Eq("projectname", project))).
						Cursor(nil)

					budgets := make([]ExpPValueModel, 0)
					_ = csrBudget.Fetch(&budgets, 0, false)
					csrBudget.Close()

					budget := 0.0
					totalDayInMonth := float64(daysIn(dtId.Month(), dtId.Year()))
					if len(budgets) > 0 {
						budget = tk.Div(budgets[0].P75NetGenMWH, totalDayInMonth)
					}
					dt.Budget = budget

					dt.AvgWindSpeed = avgwindspeed

					expws := 0.0
					dt.ExpWindSpeed = expws

					pipeAlarmParent := []tk.M{
						tk.M{}.Set("$match", tk.M{}.
							Set("projectname", project).
							Set("turbine", turbineX).
							Set("reduceavailability", true).
							Set("startdateinfo.dateid", dtId)),
						tk.M{}.Set("$group", tk.M{}.
							Set("_id", "").
							Set("count", tk.M{}.Set("$sum", 1))),
					}
					csrAlarmParent, _ := ctx.NewQuery().
						Command("pipe", pipeAlarmParent).
						From(new(Alarm).TableName()).
						Cursor(nil)

					alarmsParent := []tk.M{}
					_ = csrAlarmParent.Fetch(&alarmsParent, 0, false)
					csrAlarmParent.Close()
					noOfFailures := 0

					if len(alarmsParent) > 0 {
						noOfFailures = alarmsParent[0].GetInt("count")
					}

					pipeAlarm := []tk.M{
						tk.M{}.Set("$match", tk.M{}.
							Set("projectname", project).
							Set("turbine", turbineX)),
						tk.M{}.Set("$unwind", "$detail"),
						tk.M{}.Set("$match", tk.M{}.
							Set("detail.detaildateinfo.dateid", dtId)),
						tk.M{}.Set("$group", tk.M{}.
							Set("_id", "").
							Set("duration", tk.M{}.Set("$sum", "$detail.duration")).
							Set("powerlost", tk.M{}.Set("$sum", "$detail.powerlost")).
							Set("count", tk.M{}.Set("$sum", 1))),
					}
					csrAlarm, _ := ctx.NewQuery().
						Command("pipe", pipeAlarm).
						From(new(Alarm).TableName()).
						Cursor(nil)

					// tk.Printfn("DEBUG-01 >> %v || \n %v || \n %v", pipeAlarm, erx, new(Alarm).TableName())
					alarms := []tk.M{}
					_ = csrAlarm.Fetch(&alarms, 0, false)
					csrAlarm.Close()

					alarmDuration := 0.0
					alarmPowerLost := 0.0

					if len(alarms) > 0 {
						alarmDuration = alarms[0].GetFloat64("duration")
						alarmPowerLost = alarms[0].GetFloat64("powerlost")
					}

					dt.DowntimeHours = alarmDuration
					dt.LostEnergy = alarmPowerLost
					dt.NoOfFailures = noOfFailures
					dt.RevenueLoss = (dt.LostEnergy * 6 * revenueMultiplier)

					pipeAlarmMachine := []tk.M{
						tk.M{}.Set("$match", tk.M{}.
							Set("machinedown", true).
							Set("projectname", project).
							Set("turbine", turbineX)),
						tk.M{}.Set("$unwind", "$detail"),
						tk.M{}.Set("$match", tk.M{}.
							Set("detail.detaildateinfo.dateid", dtId)),
						tk.M{}.
							Set("$group", tk.M{}.
								Set("_id", "").
								Set("duration", tk.M{}.Set("$sum", "$detail.duration")).
								Set("powerlost", tk.M{}.Set("$sum", "$detail.powerlost")))}
					csrAlarmMachine, _ := ctx.NewQuery().
						Command("pipe", pipeAlarmMachine).
						From(new(Alarm).TableName()).
						Cursor(nil)

					alarmsMachine := []tk.M{}
					_ = csrAlarmMachine.Fetch(&alarmsMachine, 0, false)
					csrAlarmMachine.Close()

					alarmDurationMachine := 0.0
					alarmPowerLostMachine := 0.0
					if len(alarmsMachine) > 0 {
						alarmDurationMachine = alarmsMachine[0].GetFloat64("duration")
						alarmPowerLostMachine = alarmsMachine[0].GetFloat64("powerlost")
					}

					pipeAlarmGrid := []tk.M{
						tk.M{}.Set("$match", tk.M{}.
							Set("griddown", true).
							Set("projectname", project).
							Set("turbine", turbineX)),
						tk.M{}.Set("$unwind", "$detail"),
						tk.M{}.Set("$match", tk.M{}.
							Set("detail.detaildateinfo.dateid", dtId)),
						tk.M{}.
							Set("$group", tk.M{}.
								Set("_id", "").
								Set("duration", tk.M{}.Set("$sum", "$detail.duration")).
								Set("powerlost", tk.M{}.Set("$sum", "$detail.powerlost")))}

					// []tk.M{tk.M{}.Set("$match", tk.M{}.Set("griddown", true).Set("projectname", project).Set("turbine", turbine).Set("startdateinfo.dateid", dtId).Set("reduceavailability", true)), tk.M{}.Set("$group", tk.M{}.Set("_id", "").Set("duration", tk.M{}.Set("$sum", "$duration")).Set("powerlost", tk.M{}.Set("$sum", "$powerlost")))}
					csrAlarmGrid, _ := ctx.NewQuery().
						Command("pipe", pipeAlarmGrid).
						From(new(Alarm).TableName()).
						Cursor(nil)

					alarmsGrid := []tk.M{}
					_ = csrAlarmGrid.Fetch(&alarmsGrid, 0, false)
					csrAlarmGrid.Close()

					alarmDurationGrid := 0.0
					alarmPowerLostGrid := 0.0
					if len(alarmsGrid) > 0 {
						alarmDurationGrid = alarmsGrid[0].GetFloat64("duration")
						alarmPowerLostGrid = alarmsGrid[0].GetFloat64("powerlost")
					}

					pipeAlarmOther := []tk.M{
						tk.M{}.Set("$match", tk.M{}.
							Set("machinedown", false).
							Set("griddown", false).
							Set("projectname", project).
							Set("turbine", turbineX)),
						tk.M{}.Set("$unwind", "$detail"),
						tk.M{}.Set("$match", tk.M{}.
							Set("detail.detaildateinfo.dateid", dtId)),
						tk.M{}.
							Set("$group", tk.M{}.
								Set("_id", "").
								Set("duration", tk.M{}.Set("$sum", "$detail.duration")).
								Set("powerlost", tk.M{}.Set("$sum", "$detail.powerlost")))}
					// []tk.M{tk.M{}.Set("$match", tk.M{}.Set("machinedown", false).Set("griddown", false).Set("projectname", project).Set("turbine", turbine).Set("startdateinfo.dateid", dtId).Set("reduceavailability", true)), tk.M{}.Set("$group", tk.M{}.Set("_id", "").Set("duration", tk.M{}.Set("$sum", "$duration")).Set("powerlost", tk.M{}.Set("$sum", "$powerlost")))}
					csrAlarmOther, _ := ctx.NewQuery().
						Command("pipe", pipeAlarmOther).
						From(new(Alarm).TableName()).
						Cursor(nil)

					alarmsOther := []tk.M{}
					_ = csrAlarmOther.Fetch(&alarmsOther, 0, false)
					csrAlarmOther.Close()

					alarmDurationOther := 0.0
					alarmPowerLostOther := 0.0
					if len(alarmsOther) > 0 {
						alarmDurationOther = alarmsOther[0].GetFloat64("duration")
						alarmPowerLostOther = alarmsOther[0].GetFloat64("powerlost")
					}

					dt.MachineDownHours = alarmDurationMachine
					dt.GridDownHours = alarmDurationGrid
					dt.OtherDowntimeHours = alarmDurationOther
					dt.MachineDownLoss = alarmPowerLostMachine
					dt.GridDownLoss = alarmPowerLostGrid
					dt.OtherDownLoss = alarmPowerLostOther

					// the calculation using data from alarm
					dt.MachineAvail = tk.Div(((600.0 * 144.0) - (dt.MachineDownHours * 3600)), 144.0*600.0)
					dt.GridAvail = tk.Div(((600.0 * 144.0) - (dt.GridDownHours * 3600)), 144.0*600.0)
					// ===================================================================

					pipeJmr := []tk.M{tk.M{}.Set("$unwind", "$sections"), tk.M{}.Set("$match", tk.M{}.Set("sections.turbine", turbineX).Set("dateinfo.monthid", monthId)), tk.M{}.Set("$group", tk.M{}.Set("_id", "$sections.turbine").Set("boetotalloss", tk.M{}.Set("$sum", "$sections.boetotalloss")))}
					csrJmr, _ := ctx.NewQuery().
						Command("pipe", pipeJmr).
						From(new(JMR).TableName()).
						Cursor(nil)

					// log.Printf("%v\n", pipeJmr)

					jmrs := []tk.M{}
					_ = csrJmr.Fetch(&jmrs, 0, false)
					csrJmr.Close()

					// log.Printf("%#v\n", jmrs)

					boetotalloss := 0.0
					if len(jmrs) > 0 {
						boetotalloss = tk.Div(jmrs[0].GetFloat64("boetotalloss"), totalDayInMonth)
					}

					dt.ElectricalLosses = boetotalloss

					dt.ProductionRatio = 0.0

					d.BaseController.Ctx.Save(dt)

					if maxdate.IsZero() || maxdate.UTC().Before(dt.DateInfo.DateId) {
						maxdate = dt.DateInfo.DateId
					}

					count++
					total++
					if count == 1000 {
						d.Log.AddLog(tk.Sprintf("Total processed data %v\n", total), sInfo)
						count = 0
					}

					// break
				}

				d.Lock()
				if _mp, _cmp := cmaxproc[project]; !_cmp || maxdate.UTC().After(_mp.UTC()) {
					cmaxproc[project] = maxdate
				}
				d.Unlock()

				d.Log.AddLog(tk.Sprintf("Total processed data %v | %v\n", turbineX, total), sInfo)
				wg.Done()
			}(turbine, v.(tk.M).GetString("project"))

			if counter%10 == 0 || len(d.BaseController.RefTurbines) == counter {
				wg.Wait()
			}
		}

		for project, _ := range maxproc {
			lep := maxproc.Get(project, new(LastExecProcess)).(*LastExecProcess)
			lep.LastDate = cmaxproc[project]
			d.saveLastExecProject(ctx, lep)
		}
	}
}

func (d *GenScadaSummary) GenerateSummaryByProjectUsingDaily(base *BaseController) {
	if base != nil {

		d.BaseController = base

		ctx, e := PrepareConnection()
		if e != nil {
			ErrorHandler(e, "Scada Summary")
			os.Exit(0)
		}

		// d.BaseController.Ctx.DeleteMany(new(ScadaSummaryByProject), dbox.Ne("_id", ""))
		arrScadaSummaryByProject := make([]*ScadaSummaryByProject, 0)

		for _, v := range d.BaseController.ProjectList {
			var turbineList []TurbineOut
			projectName := v.Value
			group := "projectname"

			if projectName != "Fleet" {
				group = "turbine"
				turbineList, _ = helper.GetTurbineList([]interface{}{projectName})
			}

			iQuery := ctx.NewQuery().From(new(ScadaSummaryDaily).TableName())

			if projectName != "Fleet" {
				iQuery = iQuery.Where(dbox.Eq("projectname", projectName))
			}

			csr, e := iQuery.
				Aggr(dbox.AggrSum, "$powerkw", "totalpower").
				Aggr(dbox.AggrSum, "$production", "energy").
				Aggr(dbox.AggrSum, "$lostenergy", "totalenergylost").
				Aggr(dbox.AggrSum, "$oktime", "totaloktime").
				// Aggr(dbox.AggrSum, "$minutes", "totalminutes").
				Aggr(dbox.AggrSum, "$griddownhours", "totalgriddowntime").
				Aggr(dbox.AggrSum, "$otherdowntimehours", "totalunknowntime").
				Aggr(dbox.AggrSum, "$machinedownhours", "totalmachinedowntime").
				// Aggr(dbox.AggrAvr, "$avgwindspeed", "avgwindspeed").
				Aggr(dbox.AggrMax, "$dateinfo.dateid", "max").
				Aggr(dbox.AggrMin, "$dateinfo.dateid", "min").
				Aggr(dbox.AggrSum, "$totalrows", "rows10min").
				Group(group).
				Cursor(nil)
			defer csr.Close()

			_ = e

			datas := []tk.M{}
			e = csr.Fetch(&datas, 0, false)

			mdl := new(ScadaSummaryByProject).New()
			mdl.ID = projectName

			items := make([]ScadaSummaryByProjectItem, 0)
			for _, data := range datas {
				id := data["_id"].(tk.M)
				turbine := id[group].(string)

				if projectName == "Fleet" {
					turbineList, _ = helper.GetTurbineList([]interface{}{turbine})
				}

				oktime := data.GetFloat64("totaloktime") / 3600 // in hour
				power := data.GetFloat64("totalpower")          //KW ke MW

				imachinedowntime := data.GetFloat64("totalmachinedowntime")
				igriddowntime := data.GetFloat64("totalgriddowntime")
				iunknowntime := data.GetFloat64("totalunknowntime")

				icount10min := data.GetFloat64("rows10min")

				maxDate := data.Get("max", time.Time{}).(time.Time)
				minDate := data.Get("min", time.Time{}).(time.Time)
				totalhour := maxDate.AddDate(0, 0, 1).UTC().Sub(minDate.UTC()).Hours()

				energy := data.GetFloat64("energy") / 1000 //KWh ke MWh

				noofturbine, capacity := int(0), float64(0)
				for _, v := range turbineList {
					if projectName != "Fleet" {
						if v.Value == turbine {
							capacity += v.Capacity
							noofturbine += 1
						}
					} else {
						if v.Project == turbine {
							capacity += v.Capacity
							noofturbine += 1
						}
					}
				}

				// tk.Println(">>", noofturbine, oktime, energy, totalhour, capacity)

				in := tk.M{}.Set("noofturbine", noofturbine).Set("oktime", oktime).Set("energy", energy).
					Set("totalhour", totalhour).Set("totalcapacity", capacity).
					Set("machinedowntime", imachinedowntime).Set("griddowntime", igriddowntime).Set("otherdowntime", iunknowntime).
					Set("counttimestamp", icount10min)

				res := helper.CalcAvailabilityAndPLF(in)

				var item ScadaSummaryByProjectItem

				item.Name = turbine
				item.NoOfWtg = noofturbine
				item.Production = power / 6
				item.PLF = res.GetFloat64("plf") / 100
				item.MachineAvail = res.GetFloat64("machineavailability") / 100
				item.TrueAvail = res.GetFloat64("totalavailability") / 100
				item.LostEnergy = data.GetFloat64("totalenergylost")
				item.DowntimeHours = imachinedowntime + igriddowntime + iunknowntime
				item.DataAvail = res.GetFloat64("dataavailability")

				items = append(items, item)
			}

			mdl.DataItems = items

			// d.BaseController.Ctx.Insert(mdl)
			arrScadaSummaryByProject = append(arrScadaSummaryByProject, mdl)
		}

		d.BaseController.Ctx.DeleteMany(new(ScadaSummaryByProject), dbox.Ne("_id", ""))
		for _, mdl := range arrScadaSummaryByProject {
			d.BaseController.Ctx.Insert(mdl)
		}
	}
}

func (d *GenScadaSummary) GenerateSummaryByMonthUsingDaily(base *BaseController) {
	if base != nil {
		d.BaseController = base

		ctx, e := PrepareConnection()
		if e != nil {
			ErrorHandler(e, "Scada Summary")
			os.Exit(0)
		}

		projectList, _ := helper.GetProjectList()
		mapRevenue := map[string]float64{}
		for _, v := range projectList {
			mapRevenue[v.Value] = v.RevenueMultiplier
		}

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
			mapbudget[tk.Sprintf("%s_%d", budget.ProjectName, budget.MonthNo)] = budget.P75NetGenMWH
			if inprojectactive(budget.ProjectName) {
				mapbudget[tk.Sprintf("fleet_%d", budget.MonthNo)] = budget.P75NetGenMWH
			}
		}

		reffexpectedws := tk.M{}
		csrexws, _ := ctx.NewQuery().From("ref_expectedwindspeed").
			Cursor(nil)

		arrtkm := []tk.M{}
		_ = csrexws.Fetch(&arrtkm, 0, false)
		csrexws.Close()
		for _, _tkm := range arrtkm {
			reffexpectedws.Set(_tkm.GetString("_id"), _tkm.GetFloat64("value"))
		}

		// d.BaseController.Ctx.DeleteMany(new(ScadaSummaryByMonth), dbox.Ne("projectname", ""))
		// ScadaSummaryByMonth
		arrScadaSummaryByMonth := make([]*ScadaSummaryByMonth, 0)
		for _, v := range d.BaseController.ProjectList {
			project := v.Value

			group := []string{}
			if project != "Fleet" {
				group = []string{"projectname", "dateinfo.monthid"}
			} else {
				group = []string{"dateinfo.monthid"}
			}

			iQuery := ctx.NewQuery().From(new(ScadaSummaryDaily).TableName())
			if project != "Fleet" {
				iQuery = iQuery.Where(dbox.Eq("projectname", project))
			}
			csr, e := iQuery.
				Aggr(dbox.AggrSum, "$powerkw", "totalpower").
				Aggr(dbox.AggrSum, "$production", "energy").
				Aggr(dbox.AggrSum, "$lostenergy", "totalenergylost").
				Aggr(dbox.AggrSum, "$oktime", "totaloktime").
				Aggr(dbox.AggrSum, "$griddownhours", "totalgriddowntime").
				Aggr(dbox.AggrSum, "$otherdowntimehours", "totalunknowntime").
				Aggr(dbox.AggrSum, "$machinedownhours", "totalmachinedowntime").
				Aggr(dbox.AggrAvr, "$avgwindspeed", "avgwindspeed"). //check if this can happened or not
				Aggr(dbox.AggrSum, "$detwindspeed.sumwindspeed", "sumwindspeed").
				Aggr(dbox.AggrSum, "$detwindspeed.countwindspeed", "countwindspeed").
				Aggr(dbox.AggrSum, "$totalrows", "totalrows").
				Aggr(dbox.AggrMax, "$dateinfo.dateid", "max").
				Aggr(dbox.AggrMin, "$dateinfo.dateid", "min").
				Group(group...).
				Cursor(nil)
			defer csr.Close()

			if e != nil {
				ErrorHandler(e, "Scada Summary")
				os.Exit(0)
			}

			datas := []tk.M{}
			e = csr.Fetch(&datas, 0, false)

			for _, data := range datas {
				id := data["_id"].(tk.M)
				imonthid := id["dateinfo_monthid"].(int)
				monthid := strconv.Itoa(imonthid)
				year := monthid[0:4]
				month := monthid[4:6]
				day := "01"

				iMonth, _ := strconv.Atoi(string(month))
				// iMonth = iMonth - 1

				dtStr := year + "-" + month + "-" + day
				dtId, _ := time.Parse("2006-01-02", dtStr)
				dtinfo := GetDateInfo(dtId)

				noofturbine := d.BaseController.TotalTurbinePerMonth[project+"_"+monthid]
				totalcapacity := d.BaseController.CapacityPerMonth[project+"_"+monthid]

				oktime := data.GetFloat64("totaloktime") / 3600 // in hour

				power := data.GetFloat64("totalpower") / 1000 //kW to MW
				energy := data.GetFloat64("energy") / 1000    // kWh to MWh

				// @CLARIFYTHIS
				// revenueTimes := 5.74                                 // check this hardcoded
				revenueTimes := mapRevenue[project]
				if project == "fleet" {
					revenueTimes = 5.74 // check this hardcoded
				}

				revenue := revenueTimes * power * 1000 //MW to kWatt
				revenueInLacs := revenue / 100000

				maxdate := data.Get("max", time.Time{}).(time.Time)
				mindate := data.Get("min", time.Time{}).(time.Time)

				totalhour := maxdate.AddDate(0, 0, 1).UTC().Sub(mindate.UTC()).Hours()
				totalrows := data.GetFloat64("totalrows")

				imachinedowntime := data.GetFloat64("totalmachinedowntime")
				igriddowntime := data.GetFloat64("totalgriddowntime")
				iunknowntime := data.GetFloat64("totalunknowntime")

				in := tk.M{}.Set("noofturbine", noofturbine).Set("oktime", oktime).Set("energy", energy).
					Set("totalhour", totalhour).Set("totalcapacity", totalcapacity).Set("counttimestamp", totalrows).
					Set("machinedowntime", imachinedowntime).Set("griddowntime", igriddowntime).Set("otherdowntime", iunknowntime)

				res := helper.CalcAvailabilityAndPLF(in)

				budget := mapbudget[tk.Sprintf("%s_%d", project, iMonth)]

				mdl := new(ScadaSummaryByMonth).New()
				mdl.ProjectName = project
				mdl.DateInfo = dtinfo
				mdl.Production = power / 6
				// mdl.ProductionLastYear = (powerlastyear / divider)
				mdl.Revenue = revenue
				mdl.RevenueInLacs = revenueInLacs
				mdl.TrueAvail = res.GetFloat64("totalavailability")
				mdl.ScadaAvail = res.GetFloat64("dataavailability")
				mdl.MachineAvail = res.GetFloat64("machineavailability")
				mdl.GridAvail = res.GetFloat64("gridavailability")
				mdl.PLF = res.GetFloat64("plf")
				mdl.Budget = budget * 1000
				mdl.AvgWindSpeed = tk.Div(data.GetFloat64("sumwindspeed"), data.GetFloat64("countwindspeed"))

				// @CLARIFYTHIS
				expwstimes := 0.133
				randno := tk.RandInt(5)
				if randno > 3 {
					expwstimes = -0.125
				}

				mdl.ExpWindSpeed = reffexpectedws.GetFloat64(tk.Sprintf("%s_%d", project, iMonth))
				if mdl.ExpWindSpeed == 0 {
					mdl.ExpWindSpeed = mdl.AvgWindSpeed + (mdl.AvgWindSpeed * expwstimes)
				}

				mdl.DowntimeHours = imachinedowntime + iunknowntime + igriddowntime
				mdl.LostEnergy = data.GetFloat64("totalenergylost") / 1000000 // Watt to GWatt
				mdl.RevenueLoss = (data.GetFloat64("totalenergylost") * revenueTimes)

				if mdl != nil {
					// d.BaseController.Ctx.Insert(mdl)
					arrScadaSummaryByMonth = append(arrScadaSummaryByMonth, mdl)
				}

			}

		}

		d.BaseController.Ctx.DeleteMany(new(ScadaSummaryByMonth), dbox.Ne("projectname", ""))
		for _, mdl := range arrScadaSummaryByMonth {
			d.BaseController.Ctx.Insert(mdl)
		}

	}
}

func (d *GenScadaSummary) getWFAnalysisData(ctx dbox.IConnection, projectName string,
	startDate time.Time, endDate time.Time, groupBy string, totalHour float64, noOfTurbine int, dividerPower float64,
	pipeMatch tk.M, pipeGroup tk.M, plfDivider float64) tk.M {

	switch groupBy {
	case "dateinfo.dateid":
		totalHour = 24.0
	case "dateinfo.weekid":
		totalHour = 7 * 24.0
	case "dateinfo.monthid":
	case "dateinfo.qtrid":
		totalHour = 0
	}

	pipes := make([]tk.M, 0)
	// pipes = append(pipes, tk.M{
	// 	"$match": tk.M{}.Set("projectname", tk.M{}.Set("$eq", projectName)).
	// 		Set("dateinfo.dateid", tk.M{}.Set("$gte", startDate).Set("$lte", endDate)),
	// })
	pipes = append(pipes, pipeMatch)
	pipes = append(pipes, tk.M{
		"$group": tk.M{}.Set("_id", pipeGroup).
			Set("power", tk.M{}.Set("$sum", "$powerkw")).
			Set("energy", tk.M{}.Set("$sum", "$production")).
			Set("windspeed", tk.M{}.Set("$avg", "$avgwindspeed")).
			Set("oktime", tk.M{}.Set("$sum", "$oktime")).
			Set("griddowntime", tk.M{}.Set("$sum", "$griddownhours")).
			Set("machinedowntime", tk.M{}.Set("$sum", "$machinedownhours")).
			Set("unknowndowntime", tk.M{}.Set("$sum", "$otherdowntimehours")).
			Set("totaltimestamp", tk.M{}.Set("$sum", 1)).
			Set("maxdate", tk.M{}.Set("$max", "$dateinfo.dateid")).
			Set("mindate", tk.M{}.Set("$min", "$dateinfo.dateid")).
			Set("minutes", tk.M{}.Set("$sum", "$totalminutes")),
	})
	pipes = append(pipes, tk.M{"$sort": tk.M{"_id": 1}})

	csr, _ := ctx.NewQuery().
		Command("pipe", pipes).
		From(new(ScadaSummaryDaily).TableName()).
		Cursor(nil)

	scadaSums := []tk.M{}
	e := csr.Fetch(&scadaSums, 0, false)
	if e != nil {
		d.Log.AddLog(e.Error(), sError)
	}
	csr.Close()

	// log.Println(scadaSums)

	id := make([]string, 0)
	group := make([]string, 0)
	power := make([]float64, 0)
	windspeed := make([]float64, 0)
	production := make([]float64, 0)
	plf := make([]float64, 0)
	totalavail := make([]float64, 0)
	machineavail := make([]float64, 0)
	gridavail := make([]float64, 0)

	// turbineList, _ := helper.GetTurbineList([]interface{}{projectName})

	// var plfDivider float64
	// for _, v := range turbineList {
	// 	plfDivider += v.Capacity
	// }

	for _, d := range scadaSums {
		_id := d.Get("_id").(tk.M)

		vid := "0"
		if groupBy != "dateinfo.dateid" {
			vid = strconv.Itoa(_id.GetInt("period"))
			vyearid, _ := strconv.Atoi(vid[0:4])
			vperiodid, _ := strconv.Atoi(vid[4:6])

			if groupBy == "dateinfo.monthid" {
				vdate, _ := time.Parse("2006-1-2", tk.Sprintf("%d-%d-%v", vyearid, vperiodid, 1))
				totalHour = float64(time.Date(vdate.Year(), vdate.Month(), 0, 0, 0, 0, 0, time.UTC).Day()) * 24.0
				if next := endDate.UTC().AddDate(0, 0, 1); next.Before(vdate.UTC().AddDate(0, 1, 0)) {
					totalHour = next.Sub(vdate.UTC()).Hours()
				}
			}
			if groupBy == "dateinfo.qtrid" {
				totalHour = float64(GetDaysNoByQuarter(vyearid, vperiodid, endDate)) * 24.0
			}
		} else {
			dateId := _id.Get("period").(time.Time)
			vid = dateId.Format("20060102")
		}

		vgroup := _id.GetString("value")
		vpower := d.GetFloat64("power") / 1000 // kW to MW
		vws := d.GetFloat64("windspeed")
		vprod := d.GetFloat64("energy") / 1000  // kWh to MWh
		oktime := d.GetFloat64("oktime") / 3600 // in hour
		griddown := d.GetFloat64("griddowntime")
		machinedown := d.GetFloat64("machinedowntime")
		unknowndown := d.GetFloat64("unknowndowntime")
		// sumTimeStamp := d.GetFloat64("totaltimestamp")
		// minutes := d.GetFloat64("minutes") / 60

		// vplf := tk.Div(vprod, (totalHour*float64(noOfTurbine)*2100.0)) * 100
		// vtotalavail := tk.Div(tk.Div(oktime, 3600.0), (totalHour*float64(noOfTurbine))) * 100
		// vgridavail := tk.Div(((totalHour*float64(noOfTurbine))-griddown), (totalHour*float64(noOfTurbine))) * 100
		// vmchavail := tk.Div(((totalHour*float64(noOfTurbine))-machinedown), (totalHour*float64(noOfTurbine))) * 100

		maxDate := d.Get("maxdate", time.Time{}).(time.Time)
		minDate := d.Get("mindate", time.Time{}).(time.Time)
		totalHour := maxDate.AddDate(0, 0, 1).UTC().Sub(minDate.UTC()).Hours()

		//vmchavail, vgridavail, _, vtotalavail, vplf := helper.GetAvailAndPLF(float64(noOfTurbine), oktime, vprod/1000, machinedown, griddown, sumTimeStamp, totalHour, minutes, plfDivider)

		// if groupBy == "dateinfo.monthid" || groupBy == "dateinfo.qtrid" {
		// 	log.Println(vid, "data = ", _id, oktime, totalHour, noOfTurbine, plfDivider, groupBy)
		// 	log.Println(vid, "MD = ", vmchavail)
		// 	log.Println(vid, "GD = ", vgridavail)
		// 	log.Println(vid, "TV = ", vtotalavail)
		// }

		in := tk.M{}.Set("noofturbine", noOfTurbine).Set("oktime", oktime).Set("energy", vprod).
			Set("totalhour", totalHour).Set("totalcapacity", plfDivider).
			Set("machinedowntime", machinedown).Set("griddowntime", griddown).Set("otherdowntime", unknowndown)

		res := helper.CalcAvailabilityAndPLF(in)

		id = append(id, vid)
		group = append(group, vgroup)
		power = append(power, tk.Div(vpower*1000, dividerPower))
		windspeed = append(windspeed, vws)
		production = append(production, tk.Div(vprod*1000, dividerPower))
		plf = append(plf, res.GetFloat64("plf"))
		totalavail = append(totalavail, res.GetFloat64("totalavailability"))
		machineavail = append(machineavail, res.GetFloat64("machineavailability"))
		gridavail = append(gridavail, res.GetFloat64("gridavailability"))
	}

	ret := tk.M{
		"Id":           id,
		"Group":        group,
		"Power":        power,
		"WindSpeed":    windspeed,
		"Production":   production,
		"PLF":          plf,
		"TotalAvail":   totalavail,
		"MachineAvail": machineavail,
		"GridAvail":    gridavail,
	}

	//log.Println(ret)

	return ret
}

func (d *GenScadaSummary) GenWFAnalysisByProject(base *BaseController) {
	if base != nil {
		d.BaseController = base

		d.BaseController.Ctx.DeleteMany(new(GWFAnalysisByProject), dbox.Ne("projectname", ""))

		ctx, e := PrepareConnection()
		if e != nil {
			ErrorHandler(e, "Scada Summary WF Analysis")
			os.Exit(0)
		}

		keys := []string{
			"Power",
			"WindSpeed",
			"Production",
			"PLF",
			"TotalAvail",
			"MachineAvail",
			"GridAvail",
		}

		projectList, _ := helper.GetProjectList()

		for _, v := range projectList {
			projectName := v.Value

			var turbineList []TurbineOut
			var plfDivider float64
			if projectName != "" {
				turbineList, _ = helper.GetTurbineList([]interface{}{projectName})
			} else {
				turbineList, _ = helper.GetTurbineList(nil)
			}
			noOfTurbine := len(turbineList)

			for _, v := range turbineList {
				plfDivider += v.Capacity
			}

			byproject := make([]*GWFAnalysisByProject, 0)
			for idx, k := range keys {
				m := new(GWFAnalysisByProject).New()
				m.Key = k
				m.ProjectName = projectName
				m.OrderNo = idx

				byproject = append(byproject, m)
			}

			_, max, _ := GetDataDateAvailable(new(ScadaSummaryDaily).TableName(), "dateinfo.dateid", nil, d.Ctx.Connection)

			lastDate := max.UTC() //time.Parse("2006-01-02", "2016-12-21")
			strEnd := lastDate.Format("02-Jan-2006")

			// getting daily data
			last12Day := lastDate.AddDate(0, 0, -11)
			strStart := last12Day.Format("02-Jan-2006")
			dateText := strStart + " to " + strEnd
			totalHourPerDay := 24.0

			pipeMatch := tk.M{
				"$match": tk.M{}.Set("projectname", tk.M{}.Set("$eq", projectName)).
					Set("dateinfo.dateid", tk.M{}.Set("$gte", last12Day).Set("$lte", lastDate)),
			}
			pipeGroupId := tk.M{
				"value":  "$projectname",
				"period": "$dateinfo.dateid",
			}
			dailyData := d.getWFAnalysisData(ctx, projectName, last12Day, lastDate, "dateinfo.dateid", totalHourPerDay, noOfTurbine, 1000000.0, pipeMatch, pipeGroupId, plfDivider)

			for _, p := range byproject {
				var val GWFAnalysisValue

				val.DateText = dateText
				items := make([]GWFAnalysisItem, 0)
				for i := 11; i >= 0; i-- {
					dDayId := lastDate.AddDate(0, 0, -1*i)
					sId := dDayId.Format("20060102")
					Ids := dailyData.Get("Id").([]string)

					var item GWFAnalysisItem

					item.OrderNo = 12 - i
					item.Title = dDayId.Format("02-01-2006")
					item.DataId = sId

					isFound := false
					for idx, id := range Ids {
						if id == sId {
							isFound = true
							dataItems := dailyData.Get(p.Key).([]float64)
							item.Value = dataItems[idx]
							break
						}
					}
					if !isFound {
						item.Value = 0
					}

					items = append(items, item)
				}

				val.ValueAvg = items
				val.ValueMin = make([]GWFAnalysisItem, 0)
				val.ValueMax = make([]GWFAnalysisItem, 0)

				p.Roll12Days = val
			}

			// getting weekly data
			lastYear, lastWeek := lastDate.ISOWeek()
			last12Week := GetPeriodBackByDate("WEEK", lastDate, 12) // lastDate.Add(-83 * 24 * time.Hour)
			strStart = last12Week.Format("02-Jan-2006")
			dateText = strStart + " to " + strEnd
			totalHourPerWeek := 24.0 * 7.0
			pipeMatch = tk.M{
				"$match": tk.M{}.Set("projectname", tk.M{}.Set("$eq", projectName)).
					Set("dateinfo.dateid", tk.M{}.Set("$gte", last12Week).Set("$lte", lastDate)),
			}
			pipeGroupId = tk.M{
				"value":  "$projectname",
				"period": "$dateinfo.weekid",
			}
			weeklyData := d.getWFAnalysisData(ctx, projectName, last12Week, lastDate, "dateinfo.weekid", totalHourPerWeek, noOfTurbine, 1000000.0, pipeMatch, pipeGroupId, plfDivider)

			for _, p := range byproject {
				var val GWFAnalysisValue

				val.DateText = dateText

				startYear := lastYear
				startWeek := lastWeek
				items := make([]GWFAnalysisItem, 0)
				for i := 0; i < 12; i++ {

					sId := tk.Sprintf("%v%v", startYear, LeftPad2Len(strconv.Itoa(startWeek), "0", 2))
					Ids := weeklyData.Get("Id").([]string)

					var item GWFAnalysisItem

					item.OrderNo = 12 - i
					item.Title = "W " + LeftPad2Len(strconv.Itoa(startWeek), "0", 2) + " " + strconv.Itoa(startYear)
					item.DataId = sId

					isFound := false
					for idx, id := range Ids {
						if id == sId {
							isFound = true
							dataItems := weeklyData.Get(p.Key).([]float64)
							item.Value = dataItems[idx]
							break
						}
					}
					if !isFound {
						item.Value = 0
					}

					items = append(items, item)

					startWeek--
					if startWeek == 0 {
						startWeek = 52
						startYear--
					}
				}

				val.ValueAvg = items
				val.ValueMin = make([]GWFAnalysisItem, 0)
				val.ValueMax = make([]GWFAnalysisItem, 0)

				p.Roll12Weeks = val
			}

			// getting monthly data
			lastMonth := int(lastDate.Month())
			last12Month := GetPeriodBackByDate("MONTH", lastDate, 12) // lastDate.Add(-364 * 24 * time.Hour)
			strStart = last12Month.Format("02-Jan-2006")
			dateText = strStart + " to " + strEnd
			totalHourPerMonth := 24.0 * 30.0
			pipeMatch = tk.M{
				"$match": tk.M{}.Set("projectname", tk.M{}.Set("$eq", projectName)).
					Set("dateinfo.dateid", tk.M{}.Set("$gte", last12Month).Set("$lte", lastDate)),
			}
			pipeGroupId = tk.M{
				"value":  "$projectname",
				"period": "$dateinfo.monthid",
			}
			monthlyData := d.getWFAnalysisData(ctx, projectName, last12Month, lastDate, "dateinfo.monthid", totalHourPerMonth, noOfTurbine, 1000000.0, pipeMatch, pipeGroupId, plfDivider)

			for _, p := range byproject {
				var val GWFAnalysisValue

				val.DateText = dateText

				startYear := lastYear
				startMonth := lastMonth
				items := make([]GWFAnalysisItem, 0)
				for i := 0; i < 12; i++ {

					sId := tk.Sprintf("%v%v", startYear, LeftPad2Len(strconv.Itoa(startMonth), "0", 2))
					Ids := monthlyData.Get("Id").([]string)

					var item GWFAnalysisItem

					item.OrderNo = 12 - i
					item.Title = LeftPad2Len(strconv.Itoa(startMonth), "0", 2) + "-" + strconv.Itoa(startYear)
					item.DataId = sId

					isFound := false
					for idx, id := range Ids {
						if id == sId {
							isFound = true
							dataItems := monthlyData.Get(p.Key).([]float64)
							item.Value = dataItems[idx]
							break
						}
					}
					if !isFound {
						item.Value = 0
					}

					items = append(items, item)

					startMonth--
					if startMonth == 0 {
						startMonth = 12
						startYear--
					}
				}

				val.ValueAvg = items
				val.ValueMin = make([]GWFAnalysisItem, 0)
				val.ValueMax = make([]GWFAnalysisItem, 0)

				p.Roll12Months = val
			}

			// getting monthly data
			last12Qtr := GetPeriodBackByDate("QTR", lastDate, 12) // lastDate.Add(3 * -365 * 24 * time.Hour)
			strStart = last12Qtr.Format("02-Jan-2006")
			dateText = strStart + " to " + strEnd
			totalHourPerQtr := 24.0 * 90.0
			pipeMatch = tk.M{
				"$match": tk.M{}.Set("projectname", tk.M{}.Set("$eq", projectName)).
					Set("dateinfo.dateid", tk.M{}.Set("$gte", last12Qtr).Set("$lte", lastDate)),
			}
			pipeGroupId = tk.M{
				"value":  "$projectname",
				"period": "$dateinfo.qtrid",
			}
			qtrData := d.getWFAnalysisData(ctx, projectName, last12Qtr, lastDate, "dateinfo.qtrid", totalHourPerQtr, noOfTurbine, 1000000.0, pipeMatch, pipeGroupId, plfDivider)

			for _, p := range byproject {
				var val GWFAnalysisValue

				val.DateText = dateText

				startYear := lastYear

				qtr := 0
				if lastMonth%3 > 0 {
					qtr = int(math.Ceil(float64(lastMonth / 3)))
					qtr = qtr + 1
				} else {
					qtr = lastMonth / 3
				}

				startQtr := qtr
				items := make([]GWFAnalysisItem, 0)
				for i := 0; i < 12; i++ {

					sId := tk.Sprintf("%v%v", startYear, LeftPad2Len(strconv.Itoa(startQtr), "0", 2))
					Ids := qtrData.Get("Id").([]string)

					var item GWFAnalysisItem

					item.OrderNo = 12 - i
					item.Title = "Q" + LeftPad2Len(strconv.Itoa(startQtr), "0", 2) + "-" + strconv.Itoa(startYear)
					item.DataId = sId

					isFound := false
					for idx, id := range Ids {
						if id == sId {
							isFound = true
							dataItems := qtrData.Get(p.Key).([]float64)
							// log.Println(p.Key, id, sId, dataItems[idx], dataItems)
							item.Value = dataItems[idx]
							break
						}
					}
					if !isFound {
						item.Value = 0
					}

					items = append(items, item)

					startQtr--
					if startQtr == 0 {
						startQtr = 4
						startYear--
					}
				}

				val.ValueAvg = items
				val.ValueMin = make([]GWFAnalysisItem, 0)
				val.ValueMax = make([]GWFAnalysisItem, 0)

				p.Roll12Quarters = val
			}

			for _, p := range byproject {
				d.BaseController.Ctx.Save(p)
			}
		}

		// log.Println(dailyData)
		// log.Println(weeklyData)
		// log.Println(monthlyData)
		// log.Println(qtrData)
	}
}

func (d *GenScadaSummary) GenWFAnalysisByTurbine1(base *BaseController) {
	if base != nil {
		d.BaseController = base

		d.BaseController.Ctx.DeleteMany(new(GWFAnalysisByTurbine1), dbox.And(dbox.Ne("projectname", ""), dbox.Ne("turbine", "")))

		ctx, e := PrepareConnection()
		if e != nil {
			ErrorHandler(e, "Scada Summary WF Analysis")
			os.Exit(0)
		}

		keys := []string{
			"Power",
			"WindSpeed",
			"Production",
			"PLF",
			"TotalAvail",
			"MachineAvail",
			"GridAvail",
		}

		projectList, _ := helper.GetProjectList()
		for _, v := range projectList {
			projectName := v.Value

			// turbines := make([]TurbineMaster, 0)
			turbines, _ := helper.GetTurbineList([]interface{}{projectName})

			// csr, e := ctx.NewQuery().From(new(TurbineMaster).TableName()).
			// 	Where(dbox.And(dbox.Eq("project", projectName))).Order("turbineid").Cursor(nil)

			// if e != nil {
			// 	log.Println(e.Error())
			// }
			// e = csr.Fetch(&turbines, 0, false)
			// csr.Close()

			_, max, _ := GetDataDateAvailable(new(ScadaSummaryDaily).TableName(), "dateinfo.dateid", nil, d.Ctx.Connection)

			lastDate := max.UTC() //time.Parse("2006-01-02", "2016-12-21")
			strEnd := lastDate.Format("02-Jan-2006")

			for _, t := range turbines {
				var plfDivider float64
				for _, v := range turbines {
					if t.Value == v.Value {
						plfDivider = v.Capacity
					}
				}

				datas := make([]*GWFAnalysisByTurbine1, 0)
				for idx, k := range keys {
					m := new(GWFAnalysisByTurbine1).New()
					m.Key = k
					m.ProjectName = projectName
					m.Turbine = t.Value
					m.OrderNo = idx

					datas = append(datas, m)
				}

				// getting daily data
				last12Day := lastDate.Add(-11 * 24 * time.Hour)
				strStart := last12Day.Format("02-Jan-2006")
				dateText := strStart + " to " + strEnd
				totalHourPerDay := 24.0
				pipeMatch := tk.M{
					"$match": tk.M{}.Set("projectname", tk.M{}.Set("$eq", projectName)).
						Set("turbine", tk.M{}.Set("$eq", t.Value)).
						Set("dateinfo.dateid", tk.M{}.Set("$gte", last12Day).Set("$lte", lastDate)),
				}
				pipeGroupId := tk.M{
					"value":  "$turbine",
					"period": "$dateinfo.dateid",
				}
				dailyData := d.getWFAnalysisData(ctx, projectName, last12Day, lastDate, "dateinfo.dateid", totalHourPerDay, 1, 1000.0, pipeMatch, pipeGroupId, plfDivider)

				for _, p := range datas {
					var val GWFAnalysisValue

					val.DateText = dateText
					items := make([]GWFAnalysisItem, 0)
					for i := 11; i >= 0; i-- {
						dDayId := lastDate.AddDate(0, 0, -1*i)
						sId := dDayId.Format("20060102")
						Ids := dailyData.Get("Id").([]string)

						var item GWFAnalysisItem

						item.OrderNo = 12 - i
						item.Title = dDayId.Format("02-01-2006")
						item.DataId = sId

						isFound := false
						for idx, id := range Ids {
							if id == sId {
								isFound = true
								dataItems := dailyData.Get(p.Key).([]float64)
								item.Value = dataItems[idx]
								break
							}
						}
						if !isFound {
							item.Value = 0
						}

						items = append(items, item)
					}

					val.ValueAvg = items
					val.ValueMin = make([]GWFAnalysisItem, 0)
					val.ValueMax = make([]GWFAnalysisItem, 0)

					p.Roll12Days = val
				}

				// getting weekly data
				lastYear, lastWeek := lastDate.ISOWeek()
				last12Week := GetPeriodBackByDate("WEEK", lastDate, 12) // lastDate.Add(-83 * 24 * time.Hour)
				strStart = last12Week.Format("02-Jan-2006")
				dateText = strStart + " to " + strEnd
				totalHourPerWeek := 24.0 * 7.0
				pipeMatch = tk.M{
					"$match": tk.M{}.Set("projectname", tk.M{}.Set("$eq", projectName)).
						Set("turbine", tk.M{}.Set("$eq", t.Value)).
						Set("dateinfo.dateid", tk.M{}.Set("$gte", last12Week).Set("$lte", lastDate)),
				}
				pipeGroupId = tk.M{
					"value":  "$turbine",
					"period": "$dateinfo.weekid",
				}
				weeklyData := d.getWFAnalysisData(ctx, projectName, last12Week, lastDate, "dateinfo.weekid", totalHourPerWeek, 1, 1000.0, pipeMatch, pipeGroupId, plfDivider)

				for _, p := range datas {
					var val GWFAnalysisValue

					val.DateText = dateText

					startYear := lastYear
					startWeek := lastWeek
					items := make([]GWFAnalysisItem, 0)
					for i := 0; i < 12; i++ {

						sId := tk.Sprintf("%v%v", startYear, LeftPad2Len(strconv.Itoa(startWeek), "0", 2))
						Ids := weeklyData.Get("Id").([]string)

						var item GWFAnalysisItem

						item.OrderNo = 12 - i
						item.Title = "W " + LeftPad2Len(strconv.Itoa(startWeek), "0", 2) + " " + strconv.Itoa(startYear)
						item.DataId = sId

						isFound := false
						for idx, id := range Ids {
							if id == sId {
								isFound = true
								dataItems := weeklyData.Get(p.Key).([]float64)
								item.Value = dataItems[idx]
								break
							}
						}
						if !isFound {
							item.Value = 0
						}

						items = append(items, item)

						startWeek--
						if startWeek == 0 {
							startWeek = 52
							startYear--
						}
					}

					val.ValueAvg = items
					val.ValueMin = make([]GWFAnalysisItem, 0)
					val.ValueMax = make([]GWFAnalysisItem, 0)

					p.Roll12Weeks = val
				}

				// getting monthly data
				lastMonth := int(lastDate.Month())
				last12Month := GetPeriodBackByDate("MONTH", lastDate, 12) // lastDate.Add(-364 * 24 * time.Hour)
				strStart = last12Month.Format("02-Jan-2006")
				dateText = strStart + " to " + strEnd
				totalHourPerMonth := 24.0 * 30.0
				pipeMatch = tk.M{
					"$match": tk.M{}.Set("projectname", tk.M{}.Set("$eq", projectName)).
						Set("turbine", tk.M{}.Set("$eq", t.Value)).
						Set("dateinfo.dateid", tk.M{}.Set("$gte", last12Month).Set("$lte", lastDate)),
				}
				pipeGroupId = tk.M{
					"value":  "$turbine",
					"period": "$dateinfo.monthid",
				}
				monthlyData := d.getWFAnalysisData(ctx, projectName, last12Month, lastDate, "dateinfo.monthid", totalHourPerMonth, 1, 1000.0, pipeMatch, pipeGroupId, plfDivider)

				for _, p := range datas {
					var val GWFAnalysisValue

					val.DateText = dateText

					startYear := lastYear
					startMonth := lastMonth
					items := make([]GWFAnalysisItem, 0)
					for i := 0; i < 12; i++ {

						sId := tk.Sprintf("%v%v", startYear, LeftPad2Len(strconv.Itoa(startMonth), "0", 2))
						Ids := monthlyData.Get("Id").([]string)

						var item GWFAnalysisItem

						item.OrderNo = 12 - i
						item.Title = LeftPad2Len(strconv.Itoa(startMonth), "0", 2) + "-" + strconv.Itoa(startYear)
						item.DataId = sId

						isFound := false
						for idx, id := range Ids {
							if id == sId {
								isFound = true
								dataItems := monthlyData.Get(p.Key).([]float64)
								item.Value = dataItems[idx]
								break
							}
						}
						if !isFound {
							item.Value = 0
						}

						items = append(items, item)

						startMonth--
						if startMonth == 0 {
							startMonth = 12
							startYear--
						}
					}

					val.ValueAvg = items
					val.ValueMin = make([]GWFAnalysisItem, 0)
					val.ValueMax = make([]GWFAnalysisItem, 0)

					p.Roll12Months = val
				}

				// getting monthly data
				last12Qtr := GetPeriodBackByDate("QTR", lastDate, 12) // lastDate.Add(3 * -365 * 24 * time.Hour)
				strStart = last12Qtr.Format("02-Jan-2006")
				dateText = strStart + " to " + strEnd
				totalHourPerQtr := 24.0 * 90.0
				pipeMatch = tk.M{
					"$match": tk.M{}.Set("projectname", tk.M{}.Set("$eq", projectName)).
						Set("turbine", tk.M{}.Set("$eq", t.Value)).
						Set("dateinfo.dateid", tk.M{}.Set("$gte", last12Qtr).Set("$lte", lastDate)),
				}
				pipeGroupId = tk.M{
					"value":  "$turbine",
					"period": "$dateinfo.qtrid",
				}
				qtrData := d.getWFAnalysisData(ctx, projectName, last12Qtr, lastDate, "dateinfo.qtrid", totalHourPerQtr, 1, 1000.0, pipeMatch, pipeGroupId, plfDivider)

				for _, p := range datas {
					var val GWFAnalysisValue

					val.DateText = dateText

					startYear := lastYear

					qtr := 0
					if lastMonth%3 > 0 {
						qtr = int(math.Ceil(float64(lastMonth / 3)))
						qtr = qtr + 1
					} else {
						qtr = lastMonth / 3
					}

					startQtr := qtr
					items := make([]GWFAnalysisItem, 0)
					for i := 0; i < 12; i++ {

						sId := tk.Sprintf("%v%v", startYear, LeftPad2Len(strconv.Itoa(startQtr), "0", 2))
						Ids := qtrData.Get("Id").([]string)

						var item GWFAnalysisItem

						item.OrderNo = 12 - i
						item.Title = "Q" + LeftPad2Len(strconv.Itoa(startQtr), "0", 2) + "-" + strconv.Itoa(startYear)
						item.DataId = sId

						isFound := false
						for idx, id := range Ids {
							if id == sId {
								isFound = true
								dataItems := qtrData.Get(p.Key).([]float64)
								// log.Println(p.Key, id, sId, dataItems[idx], dataItems)
								item.Value = dataItems[idx]
								break
							}
						}
						if !isFound {
							item.Value = 0
						}

						items = append(items, item)

						startQtr--
						if startQtr == 0 {
							startQtr = 4
							startYear--
						}
					}

					val.ValueAvg = items
					val.ValueMin = make([]GWFAnalysisItem, 0)
					val.ValueMax = make([]GWFAnalysisItem, 0)

					p.Roll12Quarters = val
				}

				for _, p := range datas {
					d.BaseController.Ctx.Save(p)
				}
			}
		}
	}
}

func (d *GenScadaSummary) GenWFAnalysisByTurbine2(base *BaseController) {
	if base != nil {
		d.BaseController = base

		d.BaseController.Ctx.DeleteMany(new(GWFAnalysisByTurbine2), dbox.And(dbox.Ne("projectname", "")))

		ctx, e := PrepareConnection()
		if e != nil {
			ErrorHandler(e, "Scada Summary WF Analysis Turbine 2")
			os.Exit(0)
		}

		keys := []string{
			"Power",
			"WindSpeed",
			"Production",
			"PLF",
			"TotalAvail",
			"MachineAvail",
			"GridAvail",
		}

		projectList, _ := helper.GetProjectList()
		for _, v := range projectList {
			projectName := v.Value
			turbines, _ := helper.GetTurbineList([]interface{}{projectName})
			noOfTurbines := len(turbines)

			var plfDivider float64
			for _, v := range turbines {
				plfDivider += v.Capacity
			}

			_, max, _ := GetDataDateAvailable(new(ScadaSummaryDaily).TableName(), "dateinfo.dateid", nil, d.Ctx.Connection)

			lastDate := max.UTC() //time.Parse("2006-01-02", "2016-12-21")
			strEnd := lastDate.Format("02-Jan-2006")

			datas := make([]*GWFAnalysisByTurbine2, 0)
			for idx, k := range keys {
				m := new(GWFAnalysisByTurbine2).New()
				m.Key = k
				m.ProjectName = projectName
				m.OrderNo = idx

				datas = append(datas, m)
			}

			// getting daily data
			last12Day := lastDate.Add(-11 * 24 * time.Hour)
			strStart := last12Day.Format("02-Jan-2006")
			dateText := strStart + " to " + strEnd
			_ = dateText
			totalHourPerDay := 24.0
			pipeMatch := tk.M{
				"$match": tk.M{}.Set("projectname", tk.M{}.Set("$eq", projectName)).
					Set("dateinfo.dateid", tk.M{}.Set("$gte", last12Day).Set("$lte", lastDate)),
			}
			pipeGroupId := tk.M{
				"value":  "$turbine",
				"period": "$dateinfo.dateid",
			}
			dailyData := d.getWFAnalysisData(ctx, projectName, last12Day, lastDate, "dateinfo.dateid", totalHourPerDay, 1, 1000.0, pipeMatch, pipeGroupId, plfDivider)

			for _, p := range datas {
				items := make([]GWFAnalysisItem2, 0)
				totalTurbine := len(turbines)
				totalValues := make([]float64, 12)
				for _, tb := range turbines {
					for i := 11; i >= 0; i-- {
						dDayId := lastDate.AddDate(0, 0, -1*i)
						sId := dDayId.Format("20060102")
						Ids := dailyData.Get("Id").([]string)
						sValue := dailyData.Get("Group").([]string)

						var item GWFAnalysisItem2

						item.OrderNo = 12 - i
						item.Title = dDayId.Format("02-01-2006")
						item.DataId = sId
						item.Turbine = tb.Value

						isFound := false
						for idx, id := range Ids {
							if id == sId && tb.Value == sValue[idx] {
								isFound = true
								dataItems := dailyData.Get(p.Key).([]float64)
								item.Value = dataItems[idx]
								// item.Turbine = dailyData.Get("Group").([]string)[idx]
								break
							}
						}
						if !isFound {
							item.Value = 0
						}

						totalValues[12-(i+1)] += item.Value

						items = append(items, item)
					}
				}

				for idx, v := range totalValues {
					dDayId := lastDate.AddDate(0, 0, -1*(11-idx))
					sId := dDayId.Format("20060102")

					var item GWFAnalysisItem2
					item.OrderNo = (idx + 1)
					item.Title = dDayId.Format("02-01-2006")
					item.DataId = sId
					item.Turbine = "Average"
					item.Value = tk.Div(v, float64(totalTurbine))

					items = append(items, item)
				}

				p.Roll12Days = items
			}

			// getting weekly data
			lastYear, lastWeek := lastDate.ISOWeek()
			last12Week := GetPeriodBackByDate("WEEK", lastDate, 12) // lastDate.Add(-83 * 24 * time.Hour)
			strStart = last12Week.Format("02-Jan-2006")
			dateText = strStart + " to " + strEnd
			totalHourPerWeek := 24.0 * 7.0
			pipeMatch = tk.M{
				"$match": tk.M{}.Set("projectname", tk.M{}.Set("$eq", projectName)).
					//Set("turbine", tk.M{}.Set("$eq", t.TurbineId)).
					Set("dateinfo.dateid", tk.M{}.Set("$gte", last12Week).Set("$lte", lastDate)),
			}
			pipeGroupId = tk.M{
				"value":  "$turbine",
				"period": "$dateinfo.weekid",
			}
			weeklyData := d.getWFAnalysisData(ctx, projectName, last12Week, lastDate, "dateinfo.weekid", totalHourPerWeek, 1, 1000.0, pipeMatch, pipeGroupId, plfDivider)

			for _, p := range datas {
				items := make([]GWFAnalysisItem2, 0)
				totalTurbine := len(turbines)
				totalValues := make([]float64, 12)
				startYear := lastYear
				startWeek := lastWeek
				for _, tb := range turbines {
					for i := 11; i >= 0; i-- {
						sId := tk.Sprintf("%v%v", startYear, LeftPad2Len(strconv.Itoa(startWeek), "0", 2))
						Ids := weeklyData.Get("Id").([]string)

						var item GWFAnalysisItem2

						item.OrderNo = 12 - i
						item.Title = "W " + LeftPad2Len(strconv.Itoa(startWeek), "0", 2) + " " + strconv.Itoa(startYear)
						item.DataId = sId
						item.Turbine = tb.Value

						isFound := false
						for idx, id := range Ids {
							if id == sId {
								isFound = true
								dataItems := weeklyData.Get(p.Key).([]float64)
								item.Value = dataItems[idx]
								// item.Turbine = dailyData.Get("Group").([]string)[idx]
								break
							}
						}
						if !isFound {
							item.Value = 0
						}

						totalValues[11-i] += item.Value

						startWeek--
						if startWeek == 0 {
							startWeek = 52
							startYear--
						}

						items = append(items, item)
					}
				}

				for idx, v := range totalValues {
					dDayId := lastDate.AddDate(0, 0, -1*(11-idx))
					sId := dDayId.Format("20060102")

					var item GWFAnalysisItem2
					item.OrderNo = (idx + 1)
					item.Title = dDayId.Format("02-01-2006")
					item.DataId = sId
					item.Turbine = "Average"
					item.Value = tk.Div(v, float64(totalTurbine))

					items = append(items, item)
				}

				p.Roll12Weeks = items
			}

			// getting monthly data
			lastMonth := int(lastDate.Month())
			last12Month := GetPeriodBackByDate("MONTH", lastDate, 12) // lastDate.Add(-364 * 24 * time.Hour)
			strStart = last12Month.Format("02-Jan-2006")
			dateText = strStart + " to " + strEnd
			totalHourPerMonth := 24.0 * 30.0
			pipeMatch = tk.M{
				"$match": tk.M{}.Set("projectname", tk.M{}.Set("$eq", projectName)).
					Set("dateinfo.dateid", tk.M{}.Set("$gte", last12Month).Set("$lte", lastDate)),
			}
			pipeGroupId = tk.M{
				"value":  "$projectname",
				"period": "$dateinfo.monthid",
			}
			monthlyData := d.getWFAnalysisData(ctx, projectName, last12Month, lastDate, "dateinfo.monthid", totalHourPerMonth, noOfTurbines, 1000000.0, pipeMatch, pipeGroupId, plfDivider)

			for _, p := range datas {
				items := make([]GWFAnalysisItem2, 0)
				totalTurbine := len(turbines)
				totalValues := make([]float64, 12)
				startYear := lastYear
				startMonth := lastMonth
				for _, tb := range turbines {
					for i := 11; i >= 0; i-- {
						sId := tk.Sprintf("%v%v", startYear, LeftPad2Len(strconv.Itoa(startMonth), "0", 2))
						Ids := monthlyData.Get("Id").([]string)

						var item GWFAnalysisItem2

						item.OrderNo = 12 - i
						item.Title = LeftPad2Len(strconv.Itoa(startMonth), "0", 2) + "-" + strconv.Itoa(startYear)
						item.DataId = sId
						item.Turbine = tb.Value

						isFound := false
						for idx, id := range Ids {
							if id == sId {
								isFound = true
								dataItems := monthlyData.Get(p.Key).([]float64)
								item.Value = dataItems[idx]
								// item.Turbine = dailyData.Get("Group").([]string)[idx]
								break
							}
						}
						if !isFound {
							item.Value = 0
						}

						totalValues[11-i] += item.Value

						startMonth--
						if startMonth == 0 {
							startMonth = 12
							startYear--
						}

						items = append(items, item)
					}
				}

				startYear = lastYear
				startMonth = lastMonth
				for idx, v := range totalValues {
					sId := tk.Sprintf("%v%v", startYear, LeftPad2Len(strconv.Itoa(startMonth), "0", 2))

					var item GWFAnalysisItem2
					item.OrderNo = (idx + 1)
					item.Title = LeftPad2Len(strconv.Itoa(startMonth), "0", 2) + "-" + strconv.Itoa(startYear)
					item.DataId = sId
					item.Turbine = "Average"
					item.Value = tk.Div(v, float64(totalTurbine))

					startMonth--
					if startMonth == 0 {
						startMonth = 12
						startYear--
					}

					items = append(items, item)
				}

				p.Roll12Months = items
			}

			// getting monthly data
			last12Qtr := GetPeriodBackByDate("QTR", lastDate, 12) // lastDate.Add(3 * -365 * 24 * time.Hour)
			strStart = last12Qtr.Format("02-Jan-2006")
			dateText = strStart + " to " + strEnd
			totalHourPerQtr := 24.0 * 90.0
			pipeMatch = tk.M{
				"$match": tk.M{}.Set("projectname", tk.M{}.Set("$eq", projectName)).
					Set("dateinfo.dateid", tk.M{}.Set("$gte", last12Qtr).Set("$lte", lastDate)),
			}
			pipeGroupId = tk.M{
				"value":  "$projectname",
				"period": "$dateinfo.qtrid",
			}
			qtrData := d.getWFAnalysisData(ctx, projectName, last12Qtr, lastDate, "dateinfo.qtrid", totalHourPerQtr, noOfTurbines, 1000000.0, pipeMatch, pipeGroupId, plfDivider)

			for _, p := range datas {
				items := make([]GWFAnalysisItem2, 0)
				totalTurbine := len(turbines)
				totalValues := make([]float64, 12)
				startYear := lastYear
				qtr := 0
				if lastMonth%3 > 0 {
					qtr = int(math.Ceil(float64(lastMonth / 3)))
					qtr = qtr + 1
				} else {
					qtr = lastMonth / 3
				}
				startQtr := qtr
				for _, tb := range turbines {
					for i := 11; i >= 0; i-- {
						sId := tk.Sprintf("%v%v", startYear, LeftPad2Len(strconv.Itoa(startQtr), "0", 2))
						Ids := qtrData.Get("Id").([]string)

						var item GWFAnalysisItem2

						item.OrderNo = 12 - i
						item.Title = "Q" + LeftPad2Len(strconv.Itoa(startQtr), "0", 2) + "-" + strconv.Itoa(startYear)
						item.DataId = sId
						item.Turbine = tb.Value

						isFound := false
						for idx, id := range Ids {
							if id == sId {
								isFound = true
								dataItems := qtrData.Get(p.Key).([]float64)
								// log.Println(p.Key, id, sId, dataItems[idx], dataItems)
								item.Value = dataItems[idx]
								break
							}
						}
						if !isFound {
							item.Value = 0
						}

						items = append(items, item)

						startQtr--
						if startQtr == 0 {
							startQtr = 4
							startYear--
						}
					}
				}

				startYear = lastYear
				startQtr = qtr
				for idx, v := range totalValues {
					sId := tk.Sprintf("%v%v", startYear, LeftPad2Len(strconv.Itoa(startQtr), "0", 2))

					var item GWFAnalysisItem2
					item.OrderNo = (idx + 1)
					item.Title = "Q" + LeftPad2Len(strconv.Itoa(startQtr), "0", 2) + "-" + strconv.Itoa(startYear)
					item.DataId = sId
					item.Turbine = "Average"
					item.Value = tk.Div(v, float64(totalTurbine))

					startQtr--
					if startQtr == 0 {
						startQtr = 4
						startYear--
					}

					items = append(items, item)
				}

				p.Roll12Quarters = items
			}

			for _, p := range datas {
				d.BaseController.Ctx.Save(p)
			}
		}
	}
}

func (d *GenScadaSummary) getLastExec(rconn dbox.IConnection, _process, _type string, _dtime time.Time, aproject []string) tk.M {

	csr, err := rconn.NewQuery().
		From(new(LastExecProcess).TableName()).
		Where(dbox.And(dbox.Eq("process", _process), dbox.Eq("type", _type))).
		Cursor(nil)

	if err != nil {
		d.Log.AddLog(tk.Sprintf("Get last exec project found %s", err.Error()), sInfo)
		return tk.M{}
	}

	defer csr.Close()

	_alep, tkm := []*LastExecProcess{}, tk.M{}

	_ = csr.Fetch(&_alep, 0, false)
	for _, _ilep := range _alep {
		tkm.Set(_ilep.ProjectName, _ilep)
	}

	for _, str := range aproject {
		if !tkm.Has(str) {
			_ilp := new(LastExecProcess)
			_ilp.Process = _process
			_ilp.Type = _type
			_ilp.ProjectName = str
			_ilp.LastDate = _dtime.AddDate(0, 0, -1)

			tkm.Set(str, _ilp.New())
			d.saveLastExecProject(rconn, _ilp.New())
		}
	}

	return tkm
}

func (d *GenScadaSummary) saveLastExecProject(rconn dbox.IConnection, lepData *LastExecProcess) {
	qSave := rconn.NewQuery().
		From(new(LastExecProcess).TableName()).
		SetConfig("multiexec", true).
		Save()
	defer qSave.Close()

	_ = qSave.Exec(tk.M{}.Set("data", lepData))

	return
}

func PopulateValidPCDev(iconn dbox.IConnection, filter tk.M) (res tk.M) {
	res = tk.M{}

	filter.Set("isvalidstate", true)

	pipe := []tk.M{}
	pipe = append(pipe, tk.M{}.Set("$match", filter))
	pipe = append(pipe, tk.M{}.Set("$group", tk.M{}.
		Set("_id", tk.M{}.Set("projectname", "$projectname").Set("turbine", "$turbine").Set("dateid", "$dateinfo.dateid")).
		Set("pcdeviation", tk.M{}.Set("$sum", "$pcdeviation"))))

	csr, err := iconn.NewQuery().
		Command("pipe", pipe).
		From(new(ScadaData).TableName()).
		Cursor(nil)

	if err != nil {
		return
	}
	defer csr.Close()

	rawres := []tk.M{}
	err = csr.Fetch(&rawres, 0, false)
	if err != nil {
		return
	}

	for _, raw := range rawres {
		id := raw["_id"].(tk.M)
		dateId := id.Get("dateid", time.Time{}).(time.Time)
		key := tk.Sprintf("%s_%s_%s", id.GetString("projectname"), id.GetString("turbine"), dateId.UTC().Format("20060102"))

		res.Set(key, raw.GetFloat64("pcdeviation"))
	}

	return
}
