package controller

import (
	. "eaciit/wfdemo-git/library/core"
	hp "eaciit/wfdemo-git/library/helper"
	. "eaciit/wfdemo-git/library/models"
	"eaciit/wfdemo-git/web/helper"
	"time"

	"github.com/eaciit/dbox"

	"github.com/eaciit/knot/knot.v1"
	tk "github.com/eaciit/toolkit"
)

type AnalyticComparisonController struct {
	App
}

func CreateAnalyticComparisonController() *AnalyticComparisonController {
	var controller = new(AnalyticComparisonController)
	return controller
}

type PayloadComparison struct {
	Project   string
	Turbine   []interface{}
	Keys      []string
	Period    string
	DateStart time.Time
	DateEnd   time.Time
}

func (m *AnalyticComparisonController) GetData(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	var (
		pipes []tk.M
		list  []tk.M
	)

	p := new(PayloadComparison)
	e := k.GetPayload(&p)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	result := tk.M{}

	if len(p.Keys) > 0 {

		/*tStart, _ := time.Parse("2006-01-02", p.DateStart.UTC().Format("2006-01-02"))
		tEnd, _ := time.Parse("2006-01-02 15:04:05", p.DateEnd.UTC().Format("2006-01-02")+" 23:59:59")*/
		tStart, tEnd, e := helper.GetStartEndDate(k, p.Period, p.DateStart, p.DateEnd)

		// tk.Printfn("%v | %v", tStart.String(), tEnd.String())

		// log.Printf("EndDate: %v \n", tEnd)

		if e != nil {
			return helper.CreateResult(false, nil, e.Error())
		}

		for _, val := range p.Keys {

			if val == "DGR"{

				dateClause := tk.M{
					"$gte" : tStart,
					"$lte" : tEnd,
				}
				matchClause := tk.M{}
				matchClause.Set("dateinfo.dateid", dateClause)
				if len(p.Turbine) > 0{
					matchClause.Set("turbine", tk.M{"$in" : p.Turbine})	
				}
				matchClause.Set("chosensite", p.Project)
				groupClause := tk.M{}
				groupClause.Set("_id", "$chosensite")
				groupClause.Set("total", tk.M{"$sum": "$genkwhday"})
				dgrPipes := []tk.M{}
				dgrPipes = append(dgrPipes, tk.M{"$match" : matchClause})
				dgrPipes = append(dgrPipes, tk.M{"$group" : groupClause})
				csr, e := DB().Connection.NewQuery().
					From(new(DGRModel).TableName()).
					Command("pipe", dgrPipes).
					Cursor(nil)
				if e != nil {
					return helper.CreateResult(false, nil, e.Error())
				}
				tmp := []tk.M{}
				e = csr.Fetch(&tmp, 0, false)
				// add by ams, 2016-10-07
				csr.Close()
				if e != nil {
					return helper.CreateResult(false, nil, e.Error())
				}
				if len(tmp) > 0{
					result.Set(val, tk.Div(tmp[0].GetFloat64("total"), 1000))
				}

			}
		}


		match := tk.M{}
		match.Set("dateinfo.dateid", tk.M{"$gte": tStart, "$lte": tEnd})
		match.Set("power", tk.M{"$gte": -200})

		if len(p.Turbine) > 0 {
			match.Set("turbine", tk.M{"$in": p.Turbine})
		}

		group := tk.M{
			"power":           tk.M{"$sum": "$power"},
			"energy":          tk.M{"$sum": "$energy"},
			"machinedowntime": tk.M{"$sum": "$machinedowntime"},
			"griddowntime":    tk.M{"$sum": "$griddowntime"},
			"oktime":          tk.M{"$sum": "$oktime"},
			"powerlost":       tk.M{"$sum": "$powerlost"},
			"totaltimestamp":  tk.M{"$sum": 1},
			"available":       tk.M{"$sum": "$available"},
			"minutes":         tk.M{"$sum": "$minutes"},
			"maxdate":         tk.M{"$max": "$dateinfo.dateid"},
			"mindate":         tk.M{"$min": "$dateinfo.dateid"},
		}

		if p.Project != "" {
			match.Set("projectname", p.Project)
			group.Set("_id", "$projectname")
		} else {
			group.Set("_id", "all")
		}

		match.Set("available", 1)

		pipes = append(pipes, tk.M{"$match": match})
		pipes = append(pipes, tk.M{"$group": group})

		csr, e := DB().Connection.NewQuery().
			From(new(ScadaData).TableName()).
			Command("pipe", pipes).
			Cursor(nil)

		if e != nil {
			return helper.CreateResult(false, nil, e.Error())
		}

		e = csr.Fetch(&list, 0, false)
		csr.Close()

		if len(list) > 0 {
			val := list[0]

			var plf, trueAvail, machineAvail, gridAvail, dataAvail, prod, plfDivider float64
			var totalTurbine float64

			// totalTurbine = 1.0
			// hourValue := val.GetFloat64("minutes") / 60.0

			var turbineList []TurbineOut
			if p.Project != "" {
				turbineList, _ = helper.GetTurbineList([]interface{}{p.Project})
			} else {
				turbineList, _ = helper.GetTurbineList(nil)
			}

			if len(p.Turbine) == 0 {
				for _, v := range turbineList {
					plfDivider += v.Capacity
					totalTurbine += 1
				}

			} else {
				for _, vt := range p.Turbine {
					for _, v := range turbineList {
						if vt == v.Value {
							plfDivider += v.Capacity
							totalTurbine += 1
						}
					}
				}
			}

			// log.Printf(">> %v | %v \n", plfDivider, totalTurbine)

			minDate := val.Get("mindate").(time.Time)
			maxDate := val.Get("maxdate").(time.Time)

			hourValue := helper.GetHourValue(tStart.UTC(), tEnd.UTC(), minDate.UTC(), maxDate.UTC())

			// log.Printf(">> %v | %v - %v | %v >> %v \n", tStart.UTC(), tEnd.UTC(), minDate.UTC(), maxDate.UTC(), hourValue)

			okTime := val.GetFloat64("oktime")
			power := val.GetFloat64("power") / 1000.0
			_ = power
			energy := val.GetFloat64("energy") / 1000.0
			revenue := energy * 5.740 * 1000
			mDownTime := val.GetFloat64("machinedowntime") / 3600.0
			gDownTime := val.GetFloat64("griddowntime") / 3600.0
			sumTimeStamp := val.GetFloat64("totaltimestamp")
			minutes := val.GetFloat64("minutes") / 60

			machineAvail, gridAvail, dataAvail, trueAvail, plf = helper.GetAvailAndPLF(totalTurbine, okTime, energy, mDownTime, gDownTime, sumTimeStamp, hourValue, minutes, plfDivider)

			// log.Printf("%v | %v | %v | %v | %v | %v | %v | %v \n", totalTurbine, okTime, energy, mDownTime, gDownTime, sumTimeStamp, hourValue, minutes)

			prod = energy

			// log.Printf("%v | %v | %v | \n", trueAvail, machineAvail, hourValue)

			for _, val := range p.Keys {
				if val == "MachineAvailability" {
					result.Set(val, tk.ToFloat64(machineAvail, 2, tk.RoundingAuto))
				} else if val == "ActualProduction" {
					result.Set(val, tk.ToFloat64(prod, 2, tk.RoundingAuto))
				} else if val == "TotalAvailability" {
					result.Set(val, tk.ToFloat64(trueAvail, 2, tk.RoundingAuto))
				} else if val == "ActualPLF" {
					result.Set(val, tk.ToFloat64(plf, 2, tk.RoundingAuto))
				} else if val == "GridAvailability" {
					result.Set(val, tk.ToFloat64(gridAvail, 2, tk.RoundingAuto))
				} else if val == "DataAvailability" {
					result.Set(val, tk.ToFloat64(dataAvail, 2, tk.RoundingAuto))
				} else if val == "Revenue" {
					result.Set(val, tk.ToFloat64(revenue, 2, tk.RoundingAuto))
				}
			}

		}

		// pvalues -----------------------

		durationMonths := 0
		monthDay := tk.M{}
		var months []interface{}
		xDate := tStart
		year := xDate.Year()
		month := int(xDate.Month())
		day := 1

		daysInYear := hp.GetDayInYear(year)

		if (tk.ToString(xDate.Year()) + "" + tk.ToString(int(xDate.Month()))) != (tk.ToString(tEnd.Year()) + "" + tk.ToString(int(tEnd.Month()))) {
		out:
			for {
				xString := tk.ToString(xDate.Year()) + "" + tk.ToString(int(xDate.Month()))
				endString := tk.ToString(tEnd.Year()) + "" + tk.ToString(int(tEnd.Month()))

				if xString != endString {
					durationMonths++
					months = append(months, int(xDate.Month()))

					if (tk.ToString(xDate.Year()) + "" + tk.ToString(int(xDate.Month()))) == (tk.ToString(tStart.Year()) + "" + tk.ToString(int(tStart.Month()))) {
						monthDay.Set(tk.ToString(tStart.Year())+""+tk.ToString(int(tStart.Month())),
							tk.M{
								"days":         daysInYear.GetInt(tk.ToString(int(xDate.Month()))) - (int(tStart.Day()) - 1),
								"totalInMonth": daysInYear.GetInt(tk.ToString(int(xDate.Month()))),
							})
					} else {
						monthDay.Set(tk.ToString(xDate.Year())+""+tk.ToString(int(xDate.Month())),
							tk.M{
								"days":         daysInYear.GetInt(tk.ToString(int(xDate.Month()))),
								"totalInMonth": daysInYear.GetInt(tk.ToString(int(xDate.Month()))),
							})
					}

					month++

					if month > 12 {
						year = year + 1
						month = 1
						daysInYear = hp.GetDayInYear(year)
					}

					xDate, e = time.Parse("2006-1-2", tk.ToString(year)+"-"+tk.ToString(month)+"-"+tk.ToString(day))
				} else {
					durationMonths++
					months = append(months, int(tEnd.Month()))
					monthDay.Set(tk.ToString(tEnd.Year())+""+tk.ToString(int(tEnd.Month())), tk.M{
						"days":         int(tEnd.Day()),
						"totalInMonth": daysInYear.GetInt(tk.ToString(int(tEnd.Month()))),
					})
					break out
				}
			}
		}

		if durationMonths == 0 {
			months = append(months, int(tEnd.Month()))
			durationMonths = 1
			monthDay.Set(tk.ToString(tEnd.Year())+""+tk.ToString(int(tEnd.Month())), tk.M{
				"days":         int(tEnd.Day()) - (int(tStart.Day()) - 1),
				"totalInMonth": daysInYear.GetInt(tk.ToString(int(tEnd.Month()))),
			})
		}

		/*for x, y := range monthDay {
			tk.Printf("monthDay: %v | %#v \n", x, y)
		}

		for _, y := range months {
			tk.Printf("month: %v \n", y)
		}*/

		csr, e = DB().Connection.NewQuery().
			From(new(ExpPValueModel).TableName()).
			Where(dbox.And(dbox.In("monthno", months...), dbox.Eq("projectname", p.Project))).
			Cursor(nil)

		if e != nil {
			return helper.CreateResult(false, nil, e.Error())
		}

		pvalues := []ExpPValueModel{}

		e = csr.Fetch(&pvalues, 0, false)
		csr.Close()

		if e != nil {
			return helper.CreateResult(false, nil, e.Error())
		}

		var totP50PLF, totP75PLF, totP90PLF, totP50Generation, totP75Generation, totP90Generation float64
		tNow := getTimeNow()

		if len(pvalues) > 0 {
			for _, pval := range pvalues {
				for _, val := range p.Keys {
					if val == "P50PLF" {
						totP50PLF += pval.P50Plf
					} else if val == "P75PLF" {
						totP75PLF += pval.P75Plf
					} else if val == "P90PLF" {
						totP90PLF += pval.P90Plf
					} else if val == "P50Generation" || val == "P75Generation" || val == "P90Generation" {
					found:
						for yearDay, data := range monthDay {
							days := data.(tk.M).GetFloat64("days")
							totalInMonth := data.(tk.M).GetFloat64("totalInMonth")
							if tk.ToInt(yearDay[0:4], tk.RoundingAuto) == tNow.Year() &&
								tk.ToInt(yearDay[4:], tk.RoundingAuto) == int(tNow.Month()) &&
								days > float64(tNow.Day()) {
								days = float64(tNow.Day())
							}

							if tk.ToInt(yearDay[4:], tk.RoundingAuto) == pval.MonthNo {
								if val == "P50Generation" {
									totP50Generation += (pval.P50NetGenMWH / totalInMonth * days)
								} else if val == "P75Generation" {
									totP75Generation += (pval.P75NetGenMWH / totalInMonth * days)
								} else if val == "P90Generation" {
									totP90Generation += (pval.P90NetGenMWH / totalInMonth * days)
								}
								break found
							}
						}
					}
				}
			}

			for _, val := range p.Keys {
				if val == "P50PLF" {
					result.Set(val, tk.ToFloat64(totP50PLF/tk.ToFloat64(durationMonths, 0, tk.RoundingAuto), 2, tk.RoundingAuto)*100)
				} else if val == "P75PLF" {
					result.Set(val, tk.ToFloat64(totP75PLF/tk.ToFloat64(durationMonths, 0, tk.RoundingAuto), 2, tk.RoundingAuto)*100)
				} else if val == "P90PLF" {
					result.Set(val, tk.ToFloat64(totP90PLF/tk.ToFloat64(durationMonths, 0, tk.RoundingAuto), 2, tk.RoundingAuto)*100)
				} else if val == "P50Generation" {
					result.Set(val, tk.ToFloat64(totP50Generation, 2, tk.RoundingAuto))
				} else if val == "P75Generation" {
					result.Set(val, tk.ToFloat64(totP75Generation, 2, tk.RoundingAuto))
				} else if val == "P90Generation" {
					result.Set(val, tk.ToFloat64(totP90Generation, 2, tk.RoundingAuto))
				}
			}
		}
	}

	return helper.CreateResult(true, result, "success")
}
