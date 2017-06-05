package generatorControllers

import (
	. "eaciit/wfdemo-git/library/helper"
	. "eaciit/wfdemo-git/library/models"
	. "eaciit/wfdemo-git/processapp/summaryGenerator/controllers"
	"eaciit/wfdemo-git/web/helper"
	"log"
	"math"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/eaciit/dbox"
	_ "github.com/eaciit/dbox/dbc/mongo"
	tk "github.com/eaciit/toolkit"
)

type GenScadaSummary struct {
	*BaseController
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

		projectList := []ProjectOut{}
		projectList = append(projectList, ProjectOut{
			Name:   "",
			Value:  "Fleet",
			Coords: []float64{},
		})

		projects, _ := helper.GetProjectList()
		projectList = append(projectList, projects...)

		for _, v := range projectList {
			project := v.Value

			filter := []*dbox.Filter{}
			filter = append(filter, dbox.Gte("power", -200))
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

			for _, data := range datas {
				id := data["_id"].(tk.M)
				imonthid := id["dateinfo_monthid"].(int)
				monthid := strconv.Itoa(imonthid)
				year := monthid[0:4]
				month := monthid[4:6]
				day := "01"

				var turbineList []TurbineOut
				noOfTurbine := 0

				if project != "Fleet" {
					turbineList, _ = helper.GetTurbineList([]interface{}{project})
				} else {
					turbineList, _ = helper.GetTurbineList(nil)
				}
				noOfTurbine = len(turbineList)

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

				var plfDivider float64

				for _, v := range turbineList {
					plfDivider += v.Capacity
				}

				machineAvail, gridAvail, scadaAvail, trueAvail, plf := helper.GetAvailAndPLF(float64(noOfTurbine), oktime*3600, energy/1000, machinedowntime, griddowntime, float64(totaldata), hourValue, minutes, plfDivider)

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

		projectList := []ProjectOut{}
		projectList = append(projectList, ProjectOut{
			Name:   "",
			Value:  "Fleet",
			Coords: []float64{},
		})

		projects, _ := helper.GetProjectList()
		projectList = append(projectList, projects...)

		for _, v := range projectList {
			var turbineList []TurbineOut
			projectName := v.Value
			group := "projectname"

			filter := []*dbox.Filter{}
			filter = append(filter, dbox.Gte("power", -200))

			if projectName != "Fleet" {
				filter = append(filter, dbox.Eq("projectname", projectName))
				group = "turbine"
				turbineList, _ = helper.GetTurbineList([]interface{}{projectName})
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

			_, max, _ := GetDataDateAvailable(new(ScadaData).TableName(), "dateinfo.dateid", nil, d.Ctx.Connection)

			// daysInMonth := GetDayInYear(max.Year())
			// days := tk.ToString(daysInMonth.GetInt(tk.ToString(int(max.Month()))))
			// tmpdt, _ := time.Parse("060102_150405", max.UTC().Format("0601")+days+"_000000")
			tmpdt := max.UTC()
			endDate := tmpdt.UTC() //time.Parse("060102_150405", max.UTC().Format("0601")+"01_000000").UTC()
			startDate := GetNormalAddDateMonth(tmpdt.UTC(), -11)

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
				}

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

				maxDate := data.Get("max").(time.Time)
				minDate := data.Get("min").(time.Time)
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
			Where(dbox.Gte("power", -200)).
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

		mapRevenue := map[string]float64{}

		for _, v := range projectList {
			mapRevenue[v.Value] = v.RevenueMultiplier
		}

		var wg sync.WaitGroup
		counter := 0

		for turbine, v := range d.BaseController.RefTurbines {
			counter++
			wg.Add(1)

			go func(turbineX string, project string) {
				filter := tk.M{}
				filter = filter.Set("projectname", tk.M{}.Set("$eq", project))
				filter = filter.Set("turbine", tk.M{}.Set("$eq", turbineX))
				filter = filter.Set("power", tk.M{}.Set("$gte", -200))

				dt := d.BaseController.GetLatest("ScadaSummaryDaily", project, turbineX)

				if dt.Format("2006") != "0001" {
					filter = filter.Set("dateinfo.dateid", tk.M{}.Set("$gt", dt))
				}

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
					Set("minutes", tk.M{}.Set("$sum", "$minutes")).
					Set("totalts", tk.M{}.Set("$sum", 1)).
					Set("griddowntime", tk.M{}.Set("$sum", "$griddowntime")).
					Set("machinedowntime", tk.M{}.Set("$sum", "$machinedowntime")).
					Set("avgwindspeed", tk.M{}.Set("$avg", "$avgwindspeed"))))

				pipe = append(pipe, tk.M{"$sort": tk.M{"_id": 1}})

				csr, _ := ctx.NewQuery().
					Command("pipe", pipe).
					From(new(ScadaData).TableName()).
					Cursor(nil)

				scadaSums := []tk.M{}
				e = csr.Fetch(&scadaSums, 0, false)
				csr.Close()

				log.Printf("%v | %v | %v \n", project, turbineX, len(scadaSums))

				// if turbine == "HBR038" {
				// 	for _, t := range pipe {
				// 		log.Printf("%#v \n", t)
				// 	}
				// }

				revenueMultiplier := mapRevenue[project]
				revenueDividerInLacs := 100000.0
				count := 0
				total := 0

				for _, data := range scadaSums {
					id := data["_id"].(tk.M)
					project := id["projectname"].(string)
					// turbine := id["turbine"].(string)
					dtInfo := id["dateinfo"].(tk.M)
					dtId := dtInfo["dateid"].(time.Time)
					//totaltime := data["totaltime"].(float64)
					power := data["power"].(float64)
					energy := data["energy"].(float64)
					// pcvalue := data["pcvalue"].(float64)
					pcdeviation := data["pcdeviation"].(float64)
					oktime := data["oktime"].(float64)
					totalts := data["totalts"].(int)
					griddowntime := data["griddowntime"].(float64)
					machinedowntime := data["machinedowntime"].(float64)
					avgwindspeed := data["avgwindspeed"].(float64)

					dt := new(ScadaSummaryDaily).New()
					dt.DateInfo = GetDateInfo(dtId)
					dt.ProjectName = project
					dt.Turbine = turbineX
					dt.PowerKw = power
					dt.Production = energy
					dt.PCDeviation = pcdeviation
					dt.Revenue = power * revenueMultiplier
					dt.RevenueInLacs = tk.Div(dt.Revenue, revenueDividerInLacs)
					dt.OkTime = oktime
					dt.TrueAvail = tk.Div(oktime, 144*600)
					dt.ScadaAvail = tk.Div(float64(totalts), 144.0)
					dt.MachineAvail = tk.Div(((600.0 * 144.0) - machinedowntime), 144.0*600.0)
					dt.GridAvail = tk.Div(((600.0 * 144.0) - griddowntime), 144.0*600.0)
					dt.TotalAvail = dt.TrueAvail

					turbineList, _ := helper.GetTurbineList([]interface{}{projectName})
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
					monthId := dtInfo["monthid"].(int)
					sMonthNo := strconv.Itoa(monthId)[4:6]
					monthNo, _ = strconv.Atoi(sMonthNo)

					csrBudget, _ := ctx.NewQuery().From(new(ExpPValueModel).TableName()).
						Where(dbox.And(dbox.Eq("monthno", monthNo), dbox.Eq("projectname", project))).
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

					alarms := []tk.M{}
					_ = csrAlarm.Fetch(&alarms, 0, false)
					csrAlarm.Close()

					alarmDuration := 0.0
					alarmPowerLost := 0.0
					noOfFailures := 0

					if len(alarms) > 0 {
						alarmDuration = alarms[0]["duration"].(float64)
						alarmPowerLost = alarms[0]["powerlost"].(float64)
						noOfFailures = alarms[0].GetInt("count")
					}

					dt.DowntimeHours = alarmDuration
					dt.LostEnergy = alarmPowerLost
					dt.NoOfFailures = noOfFailures
					dt.RevenueLoss = (dt.LostEnergy * 6 * revenueMultiplier)

					pipeAlarm0 := []tk.M{
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
					csrAlarm0, _ := ctx.NewQuery().
						Command("pipe", pipeAlarm0).
						From(new(Alarm).TableName()).
						Cursor(nil)

					alarms0 := []tk.M{}
					_ = csrAlarm0.Fetch(&alarms0, 0, false)
					csrAlarm0.Close()

					alarmDuration0 := 0.0
					alarmPowerLost0 := 0.0
					if len(alarms0) > 0 {
						alarmDuration0 = alarms0[0]["duration"].(float64)
						alarmPowerLost0 = alarms0[0]["powerlost"].(float64)
					}

					pipeAlarm1 := []tk.M{
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
					csrAlarm1, _ := ctx.NewQuery().
						Command("pipe", pipeAlarm1).
						From(new(Alarm).TableName()).
						Cursor(nil)

					alarms1 := []tk.M{}
					_ = csrAlarm1.Fetch(&alarms1, 0, false)
					csrAlarm1.Close()

					alarmDuration1 := 0.0
					alarmPowerLost1 := 0.0
					if len(alarms1) > 0 {
						alarmDuration1 = alarms1[0]["duration"].(float64)
						alarmPowerLost1 = alarms1[0]["powerlost"].(float64)
					}

					pipeAlarm2 := []tk.M{
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
					csrAlarm2, _ := ctx.NewQuery().
						Command("pipe", pipeAlarm2).
						From(new(Alarm).TableName()).
						Cursor(nil)

					alarms2 := []tk.M{}
					_ = csrAlarm2.Fetch(&alarms2, 0, false)
					csrAlarm2.Close()

					alarmDuration2 := 0.0
					alarmPowerLost2 := 0.0
					if len(alarms2) > 0 {
						alarmDuration2 = alarms2[0]["duration"].(float64)
						alarmPowerLost2 = alarms2[0]["powerlost"].(float64)
					}

					dt.MachineDownHours = alarmDuration0
					dt.GridDownHours = alarmDuration1
					dt.OtherDowntimeHours = alarmDuration2
					dt.MachineDownLoss = alarmPowerLost0
					dt.GridDownLoss = alarmPowerLost1
					dt.OtherDownLoss = alarmPowerLost2

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
						boetotalloss = tk.Div(jmrs[0]["boetotalloss"].(float64), totalDayInMonth)
					}

					dt.ElectricalLosses = boetotalloss

					dt.ProductionRatio = 0.0

					d.BaseController.Ctx.Insert(dt)

					count++
					total++
					if count == 1000 {
						log.Printf("Total processed data %v\n", total)
						count = 0
					}

					// break
				}
				log.Printf("Total processed data %v | %v\n", turbineX, total)
				wg.Done()
			}(turbine, v.(tk.M).GetString("project"))

			if counter%10 == 0 || len(d.BaseController.RefTurbines) == counter {
				wg.Wait()
			}
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
			Set("totaltimestamp", tk.M{}.Set("$sum", 1)).
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
		log.Println(e.Error())
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
		//log.Println(_id)
		vid := "0"
		if groupBy != "dateinfo.dateid" {
			vid = strconv.Itoa(_id.GetInt("period"))
			vyearid, _ := strconv.Atoi(vid[0:4])
			vperiodid, _ := strconv.Atoi(vid[4:6])

			if groupBy == "dateinfo.monthid" {
				vdate, _ := time.Parse("2006-01-02", tk.Sprintf("%v-%v-%v", vyearid, vperiodid, 1))
				totalHour = float64(time.Date(vdate.Year(), vdate.Month(), 0, 0, 0, 0, 0, time.UTC).Day()) * 24.0
			}
			if groupBy == "dateinfo.qtrid" {
				totalHour = float64(GetDaysNoByQuarter(vyearid, vperiodid, endDate)) * 24.0
			}
		} else {
			dateId := _id.Get("period").(time.Time)
			vid = dateId.Format("20060102")
		}

		vgroup := _id.GetString("value")
		//log.Println(vgroup)
		vpower := d.GetFloat64("power")
		vws := d.GetFloat64("windspeed")
		vprod := d.GetFloat64("energy")
		oktime := d.GetFloat64("oktime")
		griddown := d.GetFloat64("griddowntime")
		machinedown := d.GetFloat64("machinedowntime")
		sumTimeStamp := d.GetFloat64("totaltimestamp")
		minutes := d.GetFloat64("minutes") / 60

		// vplf := tk.Div(vprod, (totalHour*float64(noOfTurbine)*2100.0)) * 100
		// vtotalavail := tk.Div(tk.Div(oktime, 3600.0), (totalHour*float64(noOfTurbine))) * 100
		// vgridavail := tk.Div(((totalHour*float64(noOfTurbine))-griddown), (totalHour*float64(noOfTurbine))) * 100
		// vmchavail := tk.Div(((totalHour*float64(noOfTurbine))-machinedown), (totalHour*float64(noOfTurbine))) * 100

		vmchavail, vgridavail, _, vtotalavail, vplf := helper.GetAvailAndPLF(float64(noOfTurbine), oktime, vprod/1000, machinedown, griddown, sumTimeStamp, totalHour, minutes, plfDivider)

		// if groupBy == "dateinfo.qtrid" {
		// 	log.Println(vid, "PLF = ", vplf, oktime, (totalHour * float64(noOfTurbine) * 3600.0))
		// 	log.Println(vid, "MD = ", vmchavail)
		// 	log.Println(vid, "GD = ", vgridavail)
		// 	log.Println(vid, "TV = ", vtotalavail)
		// }

		id = append(id, vid)
		group = append(group, vgroup)
		power = append(power, tk.Div(vpower, dividerPower))
		windspeed = append(windspeed, vws)
		production = append(production, tk.Div(vprod, dividerPower))
		plf = append(plf, vplf)
		totalavail = append(totalavail, vtotalavail)
		machineavail = append(machineavail, vmchavail)
		gridavail = append(gridavail, vgridavail)
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
