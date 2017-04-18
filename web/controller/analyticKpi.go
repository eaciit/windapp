package controller

import (
	. "eaciit/wfdemo-git/library/core"
	hp "eaciit/wfdemo-git/library/helper"
	. "eaciit/wfdemo-git/library/models"
	"eaciit/wfdemo-git/web/helper"
	"strconv"
	"strings"
	"time"

	"github.com/eaciit/dbox"
	"github.com/eaciit/knot/knot.v1"
	tk "github.com/eaciit/toolkit"
)

type AnalyticKpiController struct {
	App
}

func CreateAnalyticKpiController() *AnalyticKpiController {
	var controller = new(AnalyticKpiController)
	return controller
}

func FloatToString(input_num float64) string {
	// to convert a float number to a string
	return strconv.FormatFloat(input_num, 'f', 2, 64)
}

func (m *AnalyticKpiController) GetScadaSummaryList(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	var (
		pipes             []tk.M
		kpiAnalysisResult []tk.M
		list              []tk.M
	)

	p := new(PayloadKPI)
	e := k.GetPayload(&p)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	/*tStart, _ := time.Parse("2006-01-02", p.DateStart.UTC().Format("2006-01-02"))
	tEnd, _ := time.Parse("2006-01-02 15:04:05", p.DateEnd.UTC().Format("2006-01-02")+" 23:59:59")*/
	tStart, tEnd, e := helper.GetStartEndDate(k, p.Period, p.DateStart, p.DateEnd)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}
	colBreakdown := p.ColumnBreakDown
	rowsBreakdown := p.RowBreakDown
	keys := []string{p.KeyA, p.KeyB, p.KeyC}

	// create period divider

	periodDivider := tk.M{}

	tStartTmp := tStart
	tEndTmp := tEnd

	if colBreakdown == "Month" {
	outMonth:
		for {
			daysInYear := hp.GetDayInYear(tStartTmp.Year())
			title := ""
			if tStartTmp.Format("200601") == tEnd.Format("200601") {
				title = strings.Join([]string{tStartTmp.Format("2 Jan 2006"), tEnd.Format("2 Jan 2006")}, " to ")
				periodDivider.Set(tStartTmp.Format("200601"), title)
				break outMonth
			} else {
				tEndTmp = tStartTmp.AddDate(0, 0, daysInYear.GetInt(tk.ToString(int(tStartTmp.Month())))-tStartTmp.Day())
				title = strings.Join([]string{tStartTmp.Format("2 Jan 2006"), tEndTmp.Format("2 Jan 2006")}, " to ")
				periodDivider.Set(tStartTmp.Format("200601"), title)
				tStartTmp = tEndTmp.AddDate(0, 0, 1)
			}
		}
	} else if colBreakdown == "Year" {
	outYear:
		for {
			daysInYear := hp.GetDayInYear(tStartTmp.Year())
			title := ""
			if tStartTmp.Format("2006") == tEnd.Format("2006") {
				title = strings.Join([]string{tStartTmp.Format("2 Jan 2006"), tEnd.Format("2 Jan 2006")}, " to ")
				periodDivider.Set(tStartTmp.Format("2006"), title)
				break outYear
			} else {
				tEndTmp, e = time.Parse("2006-01-02", strings.Join([]string{
					tk.ToString(tStartTmp.Year()),
					"12",
					tk.ToString(daysInYear.GetInt("12"))},
					"-"))

				if e == nil {
					title = strings.Join([]string{tStartTmp.Format("2 Jan 2006"), tEndTmp.Format("2 Jan 2006")}, " to ")
					periodDivider.Set(tStartTmp.Format("2006"), title)
					tStartTmp = tEndTmp.AddDate(0, 0, 1)
				} else {
					tk.Println(e.Error())
				}

			}
		}
	}

	// tk.Printf("%#v \n", periodDivider)

	// ------

	match := tk.M{}
	groupId := tk.M{}

	match.Set("dateinfo.dateid", tk.M{"$gte": tStart, "$lte": tEnd})

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

	if rowsBreakdown == "Project" {
		if p.Project != "" {
			anProject := strings.Split(p.Project, "(")
			match.Set("projectname", strings.TrimRight(anProject[0], " "))
		}
		groupId.Set("id1", "$projectname")
	} else if rowsBreakdown == "Turbine" {
		if len(p.Turbine) > 0 {
			match.Set("turbine", tk.M{"$in": p.Turbine})
		}
		groupId.Set("id1", "$turbine")
	}

	if colBreakdown == "Date" {
		groupId.Set("id2", "$dateinfo.dateid")
	} else if colBreakdown == "Month" {
		groupId.Set("id2", "$dateinfo.monthid")
		groupId.Set("id3", "$dateinfo.monthdesc")
	} else if colBreakdown == "Year" {
		groupId.Set("id2", "$dateinfo.year")
	}

	group.Set("_id", groupId)

	match.Set("available", 1)

	pipes = append(pipes, tk.M{"$match": match})
	pipes = append(pipes, tk.M{"$group": group})
	pipes = append(pipes, tk.M{"$sort": tk.M{"_id": 1}})

	csr, e := DB().Connection.NewQuery().
		From(new(ScadaData).TableName()).
		Command("pipe", pipes).
		Cursor(nil)

	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	e = csr.Fetch(&list, 0, false)

	// add by ams, 2016-10-07
	csr.Close()

	result := make(map[string]interface{})

	for _, val := range list {
		var plf, trueAvail, machineAvail, gridAvail, dataAvail, prod, revenue, totalTurbine, hourValue float64

		if rowsBreakdown == "Turbine" {
			totalTurbine = 1.0
		} else if len(p.Turbine) == 0 {
			totalTurbine = 24.0
		} else {
			totalTurbine = tk.ToFloat64(len(p.Turbine), 1, tk.RoundingAuto)
		}

		minDate := val.Get("mindate").(time.Time)
		maxDate := val.Get("maxdate").(time.Time)

		if colBreakdown == "Date" {
			id := val.Get("_id").(tk.M)
			id1 := id.Get("id2").(time.Time)
			hourValue = helper.GetHourValue(id1.UTC(), id1.UTC(), minDate.UTC(), maxDate.UTC())
		} else {
			hourValue = helper.GetHourValue(tStart.UTC(), tEnd.UTC(), minDate.UTC(), maxDate.UTC())
		}

		okTime := val.GetFloat64("oktime")
		power := val.GetFloat64("power") / 1000.0
		energy := val.GetFloat64("energy") / 1000 //power / 6
		mDownTime := val.GetFloat64("machinedowntime") / 3600.0
		gDownTime := val.GetFloat64("griddowntime") / 3600.0
		sumTimeStamp := val.GetFloat64("totaltimestamp")
		minutes := val.GetFloat64("minutes") / 60

		machineAvail, gridAvail, dataAvail, trueAvail, plf = helper.GetAvailAndPLF(totalTurbine, okTime, energy, mDownTime, gDownTime, sumTimeStamp, hourValue, minutes)

		prod = energy
		revenue = power * 5.740 * 1000

		resVal := tk.M{}
		resVal.Set("MachineAvailability", tk.ToFloat64((machineAvail), 2, tk.RoundingAuto))
		resVal.Set("Production", tk.ToFloat64(prod, 2, tk.RoundingAuto))
		resVal.Set("TotalAvailability", tk.ToFloat64((trueAvail), 2, tk.RoundingAuto))
		resVal.Set("PLF", tk.ToFloat64(plf, 2, tk.RoundingAuto))
		resVal.Set("GridAvailability", tk.ToFloat64((gridAvail), 2, tk.RoundingAuto))
		resVal.Set("DataAvailability", tk.ToFloat64((dataAvail), 2, tk.RoundingAuto))

		resVal.Set("MachineAvailabilityUnit", "%")
		resVal.Set("ProductionUnit", "MWh")
		resVal.Set("TotalAvailabilityUnit", "%")
		resVal.Set("PLFUnit", "%")
		resVal.Set("GridAvailabilityUnit", "%")
		resVal.Set("DataAvailabilityUnit", "%")

		/*if revenue/100000 < 0 {
			resVal.Set("Revenue", tk.ToFloat64(revenue, 2, tk.RoundingAuto))
			resVal.Set("RevenueUnit", "Rupee")
		} else {*/
		resVal.Set("Revenue", tk.ToFloat64(revenue/100000.0, 2, tk.RoundingAuto))
		resVal.Set("RevenueUnit", "Lacs")
		// }

		tmpRes := tk.M{}
		for idx, key := range keys {
			if resVal.Get(key) != nil {
				unit := resVal.GetString(key + "Unit")
				res := resVal.GetFloat64(key)
				if idx == 0 {
					tmpRes.Set("KeyA", res)
					tmpRes.Set("TitleKeyA", unit)
				} else if idx == 1 && strings.ToLower(p.KeyB) != "none" {
					tmpRes.Set("KeyB", res)
					tmpRes.Set("TitleKeyB", unit)
				} else if idx == 2 && strings.ToLower(p.KeyC) != "none" {
					tmpRes.Set("KeyC", res)
					tmpRes.Set("TitleKeyC", unit)
				}
			}
		}

		id := val.Get("_id").(tk.M)

		if colBreakdown == "Date" {
			dt := id.Get("id2").(time.Time).UTC()
			tmpRes.Set("Name", dt.Format("02 Jan 2006"))
			tmpRes.Set("YearMonth", dt.Format("200601"))
		} else if colBreakdown == "Month" {
			tmpRes.Set("Name", id.GetString("id3")+" <br/> "+periodDivider.GetString(id.GetString("id2")))
			tmpRes.Set("YearMonth", id.GetString("id2"))
		} else if colBreakdown == "Year" {
			tmpRes.Set("Name", id.GetString("id2")+" <br/> "+periodDivider.GetString(id.GetString("id2")))
			tmpRes.Set("YearMonth", id.GetString("id2")+"00")
		}

		id1 := id.GetString("id1")

		if result[id1] != nil {
			tmp := result[id1].([]tk.M)
			tmp = append(tmp, tmpRes)
			result[id1] = tmp
		} else {
			tmp := []tk.M{}
			tmp = append(tmp, tmpRes)
			result[id1] = tmp
		}
	}

	// pvalues -----------------------

	isExp := false
	expList := []string{"P50Generation", "P50PLF", "P75Generation", "P75PLF", "P90Generation", "P90PLF"}
	expFields := tk.M{}
	monthInYear := tk.M{}

	for idx, key := range keys {
		for _, exp := range expList {
			if key == exp {
				isExp = true
				expFields.Set(key, idx)
			}
		}
	}

	pValues := tk.M{}
	monthDay := tk.M{}
	var months []interface{}

	if isExp {
		durationMonths := 0

		xDate := tStart
		year := xDate.Year()
		month := int(xDate.Month())
		day := 1

		daysInYear := hp.GetDayInYear(year)

		if (tk.ToString(xDate.Year()) + "" + tk.ToString(int(xDate.Month()))) != (tk.ToString(tEnd.Year()) + "" + tk.ToString(int(tEnd.Month()))) {
		out:
			for {
				xString := xDate.Format("200601")
				endString := tEnd.Format("200601")

				if xString != endString {
					durationMonths++
					months = append(months, int(xDate.Month()))

					if xDate.Format("200601") == tStart.Format("200601") {
						monthDay.Set(tStart.Format("200601"),
							tk.M{
								"days":         daysInYear.GetInt(tk.ToString(int(xDate.Month()))) - (int(tStart.Day()) - 1),
								"totalInMonth": daysInYear.GetInt(tk.ToString(int(xDate.Month()))),
							})

					} else {
						monthDay.Set(xDate.Format("200601"),
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
					monthStr := tk.ToString(int(tEnd.Month()))
					if len(monthStr) == 1 {
						monthStr = "0" + monthStr
					}
					monthDay.Set(tk.ToString(tEnd.Year())+""+monthStr,
						tk.M{
							"days":         int(tEnd.Day()),
							"totalInMonth": daysInYear.GetInt(tk.ToString(int(tEnd.Month()))),
						})
					break out
				}
			}
		}

		if durationMonths == 0 {
			months = append(months, int(tEnd.Month()))
			monthDay.Set(tEnd.Format("200601"),
				tk.M{
					"days":         int(tEnd.Day()) - (int(tStart.Day()) - 1),
					"totalInMonth": daysInYear.GetInt(tk.ToString(int(tEnd.Month()))),
				})
		}

		csr, e = DB().Connection.NewQuery().
			From(new(ExpPValueModel).TableName()).
			Where(dbox.In("monthno", months...)).
			Cursor(nil)

		if e != nil {
			return helper.CreateResult(false, nil, e.Error())
		}

		expS := []ExpPValueModel{}

		e = csr.Fetch(&expS, 0, false)
		// add by ams, 2016-10-07
		csr.Close()

		if e != nil {
			return helper.CreateResult(false, nil, e.Error())
		}

		if len(expS) > 0 {
			for keyYearDay, data := range monthDay {

				year := keyYearDay[:4]
				month := keyYearDay[4:]

				days := data.(tk.M).GetFloat64("days")
				if colBreakdown == "Date" {
					days = 1
				}
				totalInMonth := data.(tk.M).GetFloat64("totalInMonth")

				if colBreakdown == "Year" {
					keyYearDay = year + "00"
				}

				for _, pval := range expS {
					if tk.ToInt(month, tk.RoundingAuto) == pval.MonthNo {

						if colBreakdown == "Year" {
							if monthInYear.Get(year) != nil {
								monthInYear.Set(year, monthInYear.GetFloat64(year)+1)
							} else {
								monthInYear.Set(year, 1.0)
							}
						}

						tmp := tk.M{}

						for key := range expFields {
							var tmpVal float64
							if key == "P50PLF" {
								tmpVal = tk.ToFloat64(pval.P50Plf*100, 2, tk.RoundingAuto)
							} else if key == "P75PLF" {
								tmpVal = tk.ToFloat64(pval.P75Plf*100, 2, tk.RoundingAuto)
							} else if key == "P90PLF" {
								tmpVal = tk.ToFloat64(pval.P90Plf*100, 2, tk.RoundingAuto)
							} else if key == "P50Generation" || key == "P75Generation" || key == "P90Generation" {
								if key == "P50Generation" {
									tmpVal = (pval.P50NetGenMWH / totalInMonth * days)
								} else if key == "P75Generation" {
									tmpVal = (pval.P75NetGenMWH / totalInMonth * days)
								} else if key == "P90Generation" {
									tmpVal = (pval.P90NetGenMWH / totalInMonth * days)
								}
							}
							if pValues != nil {
								if pValues.Get(keyYearDay) != nil {
									tmp = pValues.Get(keyYearDay).(tk.M)
								}
							}

							if (colBreakdown == "Year" || colBreakdown == "Month") && tmp.Get(key) != nil {
								tmpVal = tmp.GetFloat64(key) + tmpVal
								tmp.Set(key, tmpVal)
							}

							if tmp.Get(key) == nil {
								tmp.Set(key, tmpVal)
							}
						}
						pValues.Set(keyYearDay, tmp)
					}
				}
			}
		}
	}
	resVal := tk.M{}
	resVal.Set("MachineAvailabilityUnit", "%")
	resVal.Set("ProductionUnit", "MWh")
	resVal.Set("TotalAvailabilityUnit", "%")
	resVal.Set("PLFUnit", "%")
	resVal.Set("GridAvailabilityUnit", "%")
	resVal.Set("DataAvailabilityUnit", "%")
	resVal.Set("RevenueUnit", "Lacs")

	resVal.Set("P50GenerationUnit", "MWh")
	resVal.Set("P50PLFUnit", "%")

	resVal.Set("P75GenerationUnit", "MWh")
	resVal.Set("P75PLFUnit", "%")

	resVal.Set("P90GenerationUnit", "MWh")
	resVal.Set("P90PLFUnit", "%")

	//Based Date=============
	_basedcolsdata := []tk.M{}
	if colBreakdown == "Date" {
		countDays := tk.ToInt(tk.Div(tEnd.Sub(tStart).Hours(), 24.0), tk.RoundingAuto)

		for i := 0; i < countDays; i++ {
			tmpCol := tk.M{}
			tmpCol.Set("KeyA", 0.0)
			if strings.ToLower(p.KeyB) != "none" {
				tmpCol.Set("KeyB", 0.0)
			}
			if strings.ToLower(p.KeyC) != "none" {
				tmpCol.Set("KeyC", 0.0)
			}

			loopmonth := tStart.AddDate(0, 0, i).UTC()
			// log.Println(tStart.AddDate(0, 0, i).String())
			tmpCol.Set("Name", tk.ToString(loopmonth.Format("02 Jan 2006")))
			tmpCol.Set("TitleKeyA", resVal.GetString(p.KeyA+"Unit"))
			if strings.ToLower(p.KeyB) != "none" {
				tmpCol.Set("TitleKeyB", resVal.GetString(p.KeyB+"Unit"))
			}
			if strings.ToLower(p.KeyC) != "none" {
				tmpCol.Set("TitleKeyC", resVal.GetString(p.KeyC+"Unit"))
			}
			tmpCol.Set("YearMonth", tk.ToString(loopmonth.Format("20060102")))

			_basedcolsdata = append(_basedcolsdata, tmpCol)
		}
	}
	//=============

	for row, column := range result {
		tmpRes := tk.M{}
		tmpRes.Set("Row", row)

		tmpCol := []tk.M{}

		for _, col := range column.([]tk.M) {
			yearMonth := col.GetString("YearMonth")
			for x, idx := range expFields {
				idExp := idx.(int)
				if pValues.Get(yearMonth) != nil {
					yearMonthVal := pValues.Get(yearMonth).(tk.M)
					val := yearMonthVal.GetFloat64(x)
					unit := "%"
					if x[3:] == "Generation" {
						unit = "MWh"
					}

					var res float64

					if (colBreakdown == "Year") && (x == "P50PLF" || x == "P75PLF" || x == "P90PLF") {
						res = tk.ToFloat64(val/monthInYear.GetFloat64(yearMonth[:4]), 2, tk.RoundingAuto)
					} else if (colBreakdown == "Month") && (x == "P50PLF" || x == "P75PLF" || x == "P90PLF") {
						res = tk.ToFloat64(val/tk.ToFloat64(len(months), 2, tk.RoundingAuto), 2, tk.RoundingAuto)
					} else {
						res = tk.ToFloat64(val, 2, tk.RoundingAuto)
					}

					if idExp == 0 {
						col.Set("KeyA", res)
						col.Set("TitleKeyA", unit)
					} else if idExp == 1 {
						col.Set("KeyB", res)
						col.Set("TitleKeyB", unit)
					} else if idExp == 2 {
						col.Set("KeyC", res)
						col.Set("TitleKeyC", unit)
					}
				}
			}
			tmpCol = append(tmpCol, col)
		}

		if colBreakdown == "Date" && len(tmpCol) != len(_basedcolsdata) {
			_tmpCol := []tk.M{}
			ix := 0
			for _, _val := range _basedcolsdata {
				// tk.Printfn(">>> %s - %s | %d >>>", _val.GetString("Name"), tmpCol[ix].GetString("Name"), len(tmpCol))
				if _val.GetString("Name") == tmpCol[ix].GetString("Name") {
					_tmpCol = append(_tmpCol, tmpCol[ix])
					ix++
					if ix >= len(tmpCol)-1 {
						ix = len(tmpCol) - 1
					}
				} else {
					_tmpCol = append(_tmpCol, _val)
				}
			}
			tmpRes.Set("Column", _tmpCol)
		} else {
			tmpRes.Set("Column", tmpCol)
		}

		kpiAnalysisResult = append(kpiAnalysisResult, tmpRes)
	}

	for _, dt := range kpiAnalysisResult {
		col := dt.Get("Column").([]tk.M)[0]
		units := make([]string, 3)
		if col != nil {
			for key, val := range col {
				if key == "TitleKeyA" {
					units[0] = val.(string)
				} else if key == "TitleKeyB" && strings.ToLower(p.KeyB) != "none" {
					units[1] = val.(string)
				} else if key == "TitleKeyC" && strings.ToLower(p.KeyC) != "none" {
					units[2] = val.(string)
				}
			}
		}
		dt.Set("Unit", units)
	}

	if len(kpiAnalysisResult) == 0 {
		var (
			turbinelist []tk.M
		)

		var unit []string

		unit = append(unit, resVal.GetString(p.KeyA+"Unit"))
		if strings.ToLower(p.KeyB) != "none" {
			unit = append(unit, resVal.GetString(p.KeyB+"Unit"))
		}
		if strings.ToLower(p.KeyC) != "none" {
			unit = append(unit, resVal.GetString(p.KeyC+"Unit"))
		}

		tmpColArr := []tk.M{}
		if colBreakdown == "Date" {
			countDays := tk.ToInt(tk.Div(tEnd.Sub(tStart).Hours(), 24.0), tk.RoundingAuto)
			for i := 0; i < countDays; i++ {
				tmpCol := tk.M{}
				tmpCol.Set("KeyA", 0.0)
				if strings.ToLower(p.KeyB) != "none" {
					tmpCol.Set("KeyB", 0.0)
				}
				if strings.ToLower(p.KeyC) != "none" {
					tmpCol.Set("KeyC", 0.0)
				}

				loopmonth := tStart.AddDate(0, 0, i).UTC()

				tmpCol.Set("Name", tk.ToString(loopmonth.Format("02 January 2006")))
				tmpCol.Set("TitleKeyA", resVal.GetString(p.KeyA+"Unit"))
				if strings.ToLower(p.KeyB) != "none" {
					tmpCol.Set("TitleKeyB", resVal.GetString(p.KeyB+"Unit"))
				}
				if strings.ToLower(p.KeyC) != "none" {
					tmpCol.Set("TitleKeyC", resVal.GetString(p.KeyC+"Unit"))
				}
				tmpCol.Set("YearMonth", tk.ToString(loopmonth.Format("20060102")))

				tmpColArr = append(tmpColArr, tmpCol)
			}
		} else if colBreakdown == "Month" {
			count := 0
			done := false
			for {
				if done {
					break
				}

				loopmonth := tStart.AddDate(0, count, 0)
				if loopmonth.Format("200601") == tEnd.Format("200601") {
					done = true
				}

				tmpCol := tk.M{}
				tmpCol.Set("KeyA", 0.0)
				if strings.ToLower(p.KeyB) != "none" {
					tmpCol.Set("KeyB", 0.0)
				}
				if strings.ToLower(p.KeyC) != "none" {
					tmpCol.Set("KeyC", 0.0)
				}

				tmpCol.Set("Name", tk.ToString(loopmonth.Format("January 2006"))+" <br/> "+periodDivider.GetString(tk.ToString(loopmonth.Format("200601"))))

				// log.Println(loopmonth.String())

				tmpCol.Set("TitleKeyA", resVal.GetString(p.KeyA+"Unit"))
				if strings.ToLower(p.KeyB) != "none" {
					tmpCol.Set("TitleKeyB", resVal.GetString(p.KeyB+"Unit"))
				}
				if strings.ToLower(p.KeyC) != "none" {
					tmpCol.Set("TitleKeyC", resVal.GetString(p.KeyC+"Unit"))
				}
				tmpCol.Set("YearMonth", tk.ToString(loopmonth.Format("200601")))

				tmpColArr = append(tmpColArr, tmpCol)

				count++
			}
		} else if colBreakdown == "Year" {
			for i := 0; i < (tk.ToInt(tEnd.Format("2006"), "") - tk.ToInt(tStart.Format("2006"), "") + 1); i++ {
				tmpCol := tk.M{}
				tmpCol.Set("KeyA", 0.0)
				if strings.ToLower(p.KeyB) != "none" {
					tmpCol.Set("KeyB", 0.0)
				}
				if strings.ToLower(p.KeyC) != "none" {
					tmpCol.Set("KeyC", 0.0)
				}

				loopmonth := tk.String2Date(tk.ToString(tk.ToInt(tStart.Format("2006"), "")+i), "YYYY").UTC()

				tmpCol.Set("Name", tk.ToString(loopmonth.Format("2006"))+" <br/> "+periodDivider.GetString(tk.ToString(loopmonth.Format("2006"))))
				tmpCol.Set("TitleKeyA", resVal.GetString(p.KeyA+"Unit"))
				if strings.ToLower(p.KeyB) != "none" {
					tmpCol.Set("TitleKeyB", resVal.GetString(p.KeyB+"Unit"))
				}
				if strings.ToLower(p.KeyC) != "none" {
					tmpCol.Set("TitleKeyC", resVal.GetString(p.KeyC+"Unit"))
				}
				tmpCol.Set("YearMonth", tk.ToString(loopmonth.Format("2006"))+"00")

				tmpColArr = append(tmpColArr, tmpCol)
			}
		}

		if len(p.Turbine) > 0 {
			for i := 0; i < len(p.Turbine); i++ {
				tmpRes := tk.M{}

				tmpRes.Set("Column", tmpColArr)
				tmpRes.Set("Row", p.Turbine[i])
				tmpRes.Set("Unit", unit)
				kpiAnalysisResult = append(kpiAnalysisResult, tmpRes)
			}
		} else {
			queryAggr := DB().Connection.NewQuery().From(new(TurbineMaster).TableName()).
				Group("turbineid")

			caggr, e := queryAggr.Cursor(nil)
			if e != nil {
				return helper.CreateResult(false, nil, e.Error())
			}
			defer caggr.Close()
			e = caggr.Fetch(&turbinelist, 0, false)
			if e != nil {
				return helper.CreateResult(false, nil, e.Error())
			}

			for _, turbine := range turbinelist {
				tmpRes := tk.M{}

				tmpRes.Set("Column", tmpColArr)
				tmpRes.Set("Row", turbine["_id"].(tk.M)["turbineid"].(string))
				tmpRes.Set("Unit", unit)
				kpiAnalysisResult = append(kpiAnalysisResult, tmpRes)
			}
		}
	}

	data := struct {
		Data []tk.M
	}{
		Data: kpiAnalysisResult,
	}

	return helper.CreateResult(true, data, "success")
}
