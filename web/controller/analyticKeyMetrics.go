package controller

import (
	. "eaciit/wfdemo-git/library/core"
	. "eaciit/wfdemo-git/library/models"
	"eaciit/wfdemo-git/web/helper"
	"strings"
	"time"

	"github.com/eaciit/dbox"
	"github.com/eaciit/knot/knot.v1"
	tk "github.com/eaciit/toolkit"
)

type AnalyticKeyMetrics struct {
	App
}

func CreateAnalyticKeyMetricsController() *AnalyticKeyMetrics {
	var controller = new(AnalyticKeyMetrics)
	return controller
}

func checkPValue(monthDay tk.M, value float64, monthno int) float64 {
	tNow := getTimeNow()
	for yearDay, data := range monthDay {
		days := data.(tk.M).GetFloat64("days")
		if tk.ToInt(yearDay[0:4], tk.RoundingAuto) == tNow.Year() &&
			tk.ToInt(yearDay[4:], tk.RoundingAuto) == int(tNow.Month()) &&
			days > float64(tNow.Day()) {
			days = float64(tNow.Day())
		}
		totalInMonth := data.(tk.M).GetFloat64("totalInMonth")
		if tk.ToInt(yearDay[4:], tk.RoundingAuto) == monthno { /* hanya bulan yang dimaksud yang di return */
			return value / totalInMonth * days
		}
	}
	return 0.0
}

func getHourMinute(tStart, tEnd, minDate, maxDate time.Time, minute float64) (hourValue, minutes float64) {
	hourValue = helper.GetHourValue(tStart.UTC(), tEnd.UTC(), minDate.UTC(), maxDate.UTC())
	minutes = minute / 60
	return
}

func (m *AnalyticKeyMetrics) GetKeyMetrics(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	p := new(helper.Payloads)
	e := k.GetPayload(&p)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	var list, dataSeries []tk.M

	keyList := []string{"P50 Generation", "P50 PLF", "P75 Generation", "P75 PLF", "P90 Generation", "P90 PLF"}
	keys := []string{p.Misc.GetString("key1"), p.Misc.GetString("key2")}
	breakDown := p.Misc.GetString("breakdown")
	// duration := p.Misc.GetInt("duration")
	// turbineCount := p.Misc.GetInt("totalturbine")
	projectName := ""
	turbines := []string{}

	if len(p.Filter.Filters) > 2 {
		for _, v := range p.Filter.Filters[2].Value.([]interface{}) {
			turbines = append(turbines, v.(string))
		}
	}

	if len(p.Filter.Filters) > 3 {
		projectName = p.Filter.Filters[3].Value.(string)
	}

	var turbineList []TurbineOut
	if projectName != "" {
		turbineList, _ = helper.GetTurbineList([]interface{}{projectName})
	} else {
		turbineList, _ = helper.GetTurbineList(nil)
	}

	turbineCount := len(turbines)

	var plfDivider float64

	if len(turbines) == 0 {
		turbineCount = len(turbineList)
		for _, v := range turbineList {
			plfDivider += v.Capacity
		}
	} else {
		for _, vt := range turbines {
			for _, v := range turbineList {
				if vt == v.Value {
					plfDivider += v.Capacity
				}
			}
		}
	}

	plfDivider = tk.ToFloat64(plfDivider, 2, tk.RoundingAuto)

	categories := []string{}

	var maxKey1, maxKey2, minKey2 float64
	catTitle := ""
	start, _ := time.Parse("2006-01-02T15:04:05.000Z", p.Filter.Filters[0].Value.(string))
	end, _ := time.Parse("2006-01-02T15:04:05.000Z", p.Filter.Filters[1].Value.(string))
	tStart, tEnd, e := helper.GetStartEndDate(k, p.Misc.GetString("period"), start.UTC(), end.UTC())

	// log.Printf("%v | %v \n", start.String(), end.String())

	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}
	startdate := tStart.Day()
	enddate := tEnd.Day()
	durationMonths, months, monthDay := helper.GetDurationInMonth(tStart, tEnd)
	monthList := tk.M{}
	measurement := ""
	// totalData := 0
	listOfYears := []int{}
	listOfMonthPVal := []int{}              /* isinya list bulan dengan format [1 2 3 4 5, dst]*/
	listOfMonths := []int{}                 /* isinya list bulan dengan format [201801 201802 201803 dst] */
	listMonthPerYearPVal := map[int][]int{} /* isinya list bulan per tahun map[2017][11 12] map[2018][1 2 3] */
	for i := tStart.Year(); i <= tEnd.Year(); i++ {
		listOfYears = append(listOfYears, i)

		if strings.Contains(breakDown, "monthid") || strings.Contains(breakDown, "year") {
			if i == tEnd.Year() {
				if tStart.Year()-tEnd.Year() == 0 {
					for j := int(tStart.Month()); j <= int(tEnd.Month()); j++ {
						listOfMonthPVal = append(listOfMonthPVal, j)
						listOfMonths = append(listOfMonths, (i*100)+j)
						listMonthPerYearPVal[i] = append(listMonthPerYearPVal[i], j)
						categories = append(categories, tk.Sprintf("%s %s", time.Month(j).String()[0:3], tk.ToString(i)[2:]))
					}
				} else {
					for j := 1; j <= int(tEnd.Month()); j++ {
						listOfMonthPVal = append(listOfMonthPVal, j)
						listOfMonths = append(listOfMonths, (i*100)+j)
						listMonthPerYearPVal[i] = append(listMonthPerYearPVal[i], j)
						categories = append(categories, tk.Sprintf("%s %s", time.Month(j).String()[0:3], tk.ToString(i)[2:]))
					}
				}
			} else {
				if i == tStart.Year() {
					for j := int(tStart.Month()); j <= 12; j++ {
						listOfMonthPVal = append(listOfMonthPVal, j)
						listOfMonths = append(listOfMonths, (i*100)+j)
						listMonthPerYearPVal[i] = append(listMonthPerYearPVal[i], j)
						categories = append(categories, tk.Sprintf("%s %s", time.Month(j).String()[0:3], tk.ToString(i)[2:]))
					}
				} else {
					for j := 1; j <= 12; j++ {
						listOfMonthPVal = append(listOfMonthPVal, j)
						listOfMonths = append(listOfMonths, (i*100)+j)
						listMonthPerYearPVal[i] = append(listMonthPerYearPVal[i], j)
						categories = append(categories, tk.Sprintf("%s %s", time.Month(j).String()[0:3], tk.ToString(i)[2:]))
					}
				}
			}
		}
	}
	// totalTurbine := 1.0
	// if !strings.Contains(breakDown, "turbine") {
	totalTurbine := tk.ToFloat64(turbineCount, 0, tk.RoundingAuto)
	totalTurbinePValue := tk.ToFloat64(len(turbineList), 0, tk.RoundingAuto)
	// }
	isExpPValue := false
	pValueMonth := map[int]float64{}

	for i, key := range keys {
		pValueMonth = map[int]float64{}
		list = []tk.M{}
		series := tk.M{}

		if !strings.Contains(breakDown, "monthid") {
			categories = []string{} /* reset categories content for second key */
		}
		if tk.HasMember(keyList, key) {
			csrPValue, e := DB().Connection.NewQuery().
				From(new(ExpPValueModel).TableName()).
				Where(dbox.And(dbox.In("monthno", months...), dbox.Eq("projectname", projectName))).
				Cursor(nil)

			defer csrPValue.Close()
			if e != nil {
				return helper.CreateResult(false, nil, e.Error())
			}

			e = csrPValue.Fetch(&list, 0, false)
			// add by ams, 2016-10-07
			if e != nil {
				return helper.CreateResult(false, nil, e.Error())
			}
			isExpPValue = true
			/* rawan duplicate makanya taruh sini karena langsung looping all turbine bukan satu per satu */
			if strings.Contains(breakDown, "turbine") {
				temp := p.Filter.Filters[2].Value.([]interface{})
				for _, turbine := range temp {
					categories = append(categories, turbine.(string))
				}
			}
		} else if key == "DGR" {
			dateClause := tk.M{
				"$gte": tStart,
				"$lte": tEnd,
			}
			matchClause := tk.M{}
			matchClause.Set("dateinfo.dateid", dateClause)
			matchClause.Set("turbine", tk.M{"$in": turbines})
			matchClause.Set("chosensite", projectName)
			groupClause := tk.M{}
			groupClause.Set("_id", p.Misc.GetString("breakdown"))
			groupClause.Set("total", tk.M{"$sum": "$genkwhday"})
			sortClause := tk.M{
				"_id": 1,
			}
			pipes := []tk.M{}
			pipes = append(pipes, tk.M{"$match": matchClause})
			pipes = append(pipes, tk.M{"$group": groupClause})
			pipes = append(pipes, tk.M{"$sort": sortClause})

			csr, e := DB().Connection.NewQuery().
				From(new(DGRModel).TableName()).
				Command("pipe", pipes).
				Cursor(nil)
			if e != nil {
				return helper.CreateResult(false, nil, e.Error())
			}
			// results := []tk.M{}
			e = csr.Fetch(&list, 0, false)
			// add by ams, 2016-10-07
			csr.Close()
			if e != nil {
				return helper.CreateResult(false, nil, e.Error())
			}
			// var tmp []float64
			// for _, res := range results {
			// 	tmp = append(tmp, tk.Div(res.GetFloat64("total"), 1000))
			// }
			// series.Set("data", tmp)
			isExpPValue = true
		} else {
			p.Misc.Set("knot_data", k)
			filter, e := p.ParseFilter()
			if e != nil {
				return helper.CreateResult(false, nil, e.Error())
			}
			fb := DB().Connection.Fb()

			// filter = append(filter, dbox.Eq("available", 1))

			fb.AddFilter(dbox.And(filter...))
			matches, e := fb.Build()

			group := tk.M{
				"powerkw":         tk.M{"$sum": "$powerkw"},
				"machinedowntime": tk.M{"$sum": "$machinedownhours"},
				"griddowntime":    tk.M{"$sum": "$griddownhours"},
				"otherdowntime":   tk.M{"$sum": "$otherdowntimehours"},
				"oktime":          tk.M{"$sum": "$oktime"},
				"counttimestamp":  tk.M{"$sum": "$totalrows"},
				"maxdate":         tk.M{"$max": "$dateinfo.dateid"},
				"mindate":         tk.M{"$min": "$dateinfo.dateid"},
			}

			group.Set("_id", tk.M{"id1": breakDown})
			if strings.Contains(breakDown, "month") {
				group.Set("_id", tk.M{"id1": breakDown, "id2": "$dateinfo.monthdesc"})
			}

			pipes := []tk.M{{"$match": matches}, {"$group": group}, {"$sort": tk.M{"_id.id1": 1}}}

			csr, e := DB().Connection.NewQuery().
				From(new(ScadaSummaryDaily).TableName()).
				Command("pipe", pipes).
				Cursor(nil)

			// for _, v := range pipes {
			// 	log.Printf("%#v \n", v)
			// }

			if e != nil {
				return helper.CreateResult(false, nil, e.Error())
			}

			e = csr.Fetch(&list, 0, false)

			// add by ams, 2016-10-07
			csr.Close()

			if e != nil {
				return helper.CreateResult(false, nil, e.Error())
			}
			isExpPValue = false
			// tk.Printf("breakDown : %s \n", breakDown)
		}

		measurement = "%"
		if i == 0 {
			series.Set("name", key)
			series.Set("type", "column")
			series.Set("axis", "Key1")
			series.Set("gap", 0.7)
			if key == "Actual Production" || strings.Contains(key, "Generation") || key == "DGR" {
				measurement = "MWh"
			}
			series.Set("satuan", measurement)
		} else {
			series.Set("name", key)
			series.Set("type", "line")
			series.Set("dashType", "solid")
			series.Set("markers", tk.M{"visible": false})
			series.Set("width", 2)
			series.Set("axis", "Key2")
			minKey2 = 100.00
			if key == "Actual Production" || strings.Contains(key, "Generation") || key == "DGR" {
				minKey2 = 99999999.99
				measurement = "MWh"
			}
			series.Set("satuan", measurement)
		}

		var datas []float64
		var values float64
		for listCount, val := range list {
			var hourValue float64
			if !isExpPValue {
				id := val.Get("_id").(tk.M)
				if strings.Contains(breakDown, "dateid") {

					id1 := id.Get("id1").(time.Time)
					// hourValue, minutes = getHourMinute(id1.UTC(), id1.UTC(), val.Get("mindate").(time.Time), val.Get("maxdate").(time.Time), val.GetFloat64("minutes"))

					hourValue = helper.GetHourValue(id1.UTC(), id1.UTC(), val.Get("mindate").(time.Time), val.Get("maxdate").(time.Time))

				} else {
					// hourValue, minutes = getHourMinute(tStart.UTC(), tEnd.UTC(), val.Get("mindate").(time.Time), val.Get("maxdate").(time.Time), val.GetFloat64("minutes"))

					hourValue = helper.GetHourValue(tStart.UTC(), tEnd.UTC(), val.Get("mindate").(time.Time), val.Get("maxdate").(time.Time))
				}
			}

			val.Set("totalhour", hourValue)
			val.Set("noofturbine", totalTurbine)
			val.Set("energy", val.GetFloat64("powerkw")/6/1000)
			val.Set("totalcapacity", plfDivider)

			// machineAvail, gridAvail, dataAvail, trueAvail, plf := helper.GetAvailAndPLF(totalTurbine, oktime, energy, mDownTime, gDownTime, sumTimeStamp, hourValue, minutes, plfDivider)
			availPLF := helper.CalcAvailabilityAndPLF(val)

			// log.Printf("%v | %v | %v | %v | %v | %v | %v | %v \n", totalTurbine, oktime, energy, mDownTime, gDownTime, sumTimeStamp, hourValue, minutes)

			switch key {
			case "Machine Availability":
				values = availPLF.GetFloat64("machineavailability") //tk.Div((minutes-(val.GetFloat64("machinedowntime")/3600.0)), (totalTurbine*hourValue)) * 100 /*percentage*/
				if strings.Contains(breakDown, "monthid") {
					ids := val.Get("_id", tk.M{}).(tk.M)
					pValueMonth[ids.GetInt("id1")] = values
				}
			case "Grid Availability":
				values = availPLF.GetFloat64("gridavailability") //tk.Div((minutes-(val.GetFloat64("griddowntime")/3600.0)), (totalTurbine*hourValue)) * 100 /*percentage*/
				if strings.Contains(breakDown, "monthid") {
					ids := val.Get("_id", tk.M{}).(tk.M)
					pValueMonth[ids.GetInt("id1")] = values
				}
			case "Total Availability":
				values = availPLF.GetFloat64("totalavailability") //tk.Div((val.GetFloat64("oktime")/3600), (totalTurbine*hourValue)) * 100
				if strings.Contains(breakDown, "monthid") {
					ids := val.Get("_id", tk.M{}).(tk.M)
					pValueMonth[ids.GetInt("id1")] = values
				}
			case "Data Availability":
				values = availPLF.GetFloat64("dataavailability") //tk.Div((tk.ToFloat64((val.GetInt("countdata")*10/60), 6, tk.RoundingAuto)), (hourValue*totalTurbine)) * 100
				if strings.Contains(breakDown, "monthid") {
					ids := val.Get("_id", tk.M{}).(tk.M)
					pValueMonth[ids.GetInt("id1")] = values
				}
			case "Actual PLF":
				values = availPLF.GetFloat64("plf") //tk.Div((val.GetFloat64("energy")/1000), (hourValue*2.1*totalTurbine)) * 100
				if strings.Contains(breakDown, "monthid") {
					ids := val.Get("_id", tk.M{}).(tk.M)
					pValueMonth[ids.GetInt("id1")] = values
				}
			case "Actual Production":
				values = val.GetFloat64("powerkw") / 6 / 1000 /*MWh*/
				if strings.Contains(breakDown, "monthid") {
					ids := val.Get("_id", tk.M{}).(tk.M)
					pValueMonth[ids.GetInt("id1")] = values
				}
			case "P50 Generation":
				value := checkPValue(monthDay, val.GetFloat64("p50netgenmwh"), val.GetInt("monthno"))
				if strings.Contains(breakDown, "dateid") || strings.Contains(breakDown, "monthid") { /* jika per hari maka akan berbeda nilai jika lintas bulan */
					values = value
				} else { /* jika selain per hari nilainya bisa di rata2 jika lintas bulan */
					values += value
				}
				if strings.Contains(breakDown, "monthid") || strings.Contains(breakDown, "year") {
					pValueMonth[val.GetInt("monthno")] = value
				}
			case "P50 PLF":
				// values += val.GetFloat64("p50plf") * 100
				value := val.GetFloat64("p50plf") * 100
				if strings.Contains(breakDown, "dateid") || strings.Contains(breakDown, "monthid") { /* jika per hari maka akan berbeda nilai jika lintas bulan */
					values = value
				} else { /* jika selain per hari nilainya bisa di rata2 jika lintas bulan */
					values += value
				}
				if strings.Contains(breakDown, "monthid") || strings.Contains(breakDown, "year") {
					pValueMonth[val.GetInt("monthno")] = value
				}
			case "P75 Generation":
				// values += checkPValue(monthDay, val.GetFloat64("p75netgenmwh"), val.GetInt("monthno"))
				value := checkPValue(monthDay, val.GetFloat64("p75netgenmwh"), val.GetInt("monthno"))
				if strings.Contains(breakDown, "dateid") || strings.Contains(breakDown, "monthid") { /* jika per hari maka akan berbeda nilai jika lintas bulan */
					values = value
				} else { /* jika selain per hari nilainya bisa di rata2 jika lintas bulan */
					values += value
				}
				if strings.Contains(breakDown, "monthid") || strings.Contains(breakDown, "year") {
					pValueMonth[val.GetInt("monthno")] = value
				}
			case "P75 PLF":
				// values += val.GetFloat64("p75plf") * 100
				value := val.GetFloat64("p75plf") * 100
				if strings.Contains(breakDown, "dateid") || strings.Contains(breakDown, "monthid") { /* jika per hari maka akan berbeda nilai jika lintas bulan */
					values = value
				} else { /* jika selain per hari nilainya bisa di rata2 jika lintas bulan */
					values += value
				}
				if strings.Contains(breakDown, "monthid") || strings.Contains(breakDown, "year") {
					pValueMonth[val.GetInt("monthno")] = value
				}
			case "P90 Generation":
				// values += checkPValue(monthDay, val.GetFloat64("p90netgenmwh"), val.GetInt("monthno"))
				value := checkPValue(monthDay, val.GetFloat64("p90netgenmwh"), val.GetInt("monthno"))
				if strings.Contains(breakDown, "dateid") || strings.Contains(breakDown, "monthid") { /* jika per hari maka akan berbeda nilai jika lintas bulan */
					values = value
				} else { /* jika selain per hari nilainya bisa di rata2 jika lintas bulan */
					values += value
				}
				if strings.Contains(breakDown, "monthid") || strings.Contains(breakDown, "year") {
					pValueMonth[val.GetInt("monthno")] = value
				}
			case "P90 PLF":
				// values += val.GetFloat64("p90plf") * 100
				value := val.GetFloat64("p90plf") * 100
				if strings.Contains(breakDown, "dateid") || strings.Contains(breakDown, "monthid") { /* jika per hari maka akan berbeda nilai jika lintas bulan */
					values = value
				} else { /* jika selain per hari nilainya bisa di rata2 jika lintas bulan */
					values += value
				}
				if strings.Contains(breakDown, "monthid") || strings.Contains(breakDown, "year") {
					pValueMonth[val.GetInt("monthno")] = value
				}
			case "DGR":
				values = tk.Div(val.GetFloat64("total"), 1000)
			}

			// plf = energy / (totalTurbine * hourValue * 2.1) * 100
			// trueAvail = (okTime / 3600) / (totalTurbine * hourValue) * 100
			// dataAvail = (sumTimeStamp * 10 / 60) / (hourValue * totalTurbine) * 100

			/*p50netgen per day = (p50netgenmwh / jumlah hari dalam bulan tersebut) * jumlah hari periode
			plf e => p50netgen per hari ne / ( 2.1 x jumlah hari periode x 24 x 24 )*/

			if !isExpPValue {
				datas = append(datas, tk.ToFloat64(values, 2, tk.RoundingAuto))
				if i == 0 {
					if values > maxKey1 {
						maxKey1 = values
					}
				} else {
					if values > maxKey2 {
						maxKey2 = values
					}
					if values < minKey2 {
						minKey2 = values
					}
				}
			}
			if key == "DGR" {
				continue
			}

			if isExpPValue {
				if strings.Contains(breakDown, "dateid") {
					jumCat := 0.0
					if listCount == 0 { /*bulan pertama*/
						catTitle = tStart.Month().String()
						if len(list) == 1 { /*jika hanya 1 bulan*/
							for iDate := startdate; iDate <= enddate; iDate++ {
								categories = append(categories, tk.ToString(iDate))
								jumCat++
							}
							catTitle += " " + tk.ToString(listOfYears[0]) /*Dec 2015*/
						} else { /*jika lebih dari 1 bulan*/
							month := val.GetInt("monthno")
							maxDays := monthDay.Get(tk.ToString(tStart.Year())+tk.ToString(month), tk.M{}).(tk.M).GetInt("totalInMonth")
							for iDate := startdate; iDate <= maxDays; iDate++ {
								categories = append(categories, tk.ToString(iDate))
								jumCat++
							}
							if len(listOfYears) > 1 { /*jika cuma 1 tahun, lanjut ke berikutnya*/
								catTitle += " " + tk.ToString(listOfYears[0]) /* Dec 2015*/
							}
						}
					} else { /*bulan kedua*/
						catTitle += " - " + tEnd.Month().String()
						if len(listOfYears) == 1 {
							catTitle += " (" + tk.ToString(listOfYears[0]) + ")" /*Dec - Jan (2016)*/
						} else {
							catTitle += " " + tk.ToString(listOfYears[1]) /* - Jan 2016*/
						}
						for iDate := 1; iDate <= enddate; iDate++ {
							categories = append(categories, tk.ToString(iDate))
							jumCat++
						}
					}
					// tk.Printf("key : %s \n", key)
					for jc := 1.0; jc <= jumCat; jc++ {
						if strings.Contains(key, "PLF") {
							datas = append(datas, values)

							if i == 0 {
								if values > maxKey1 {
									maxKey1 = values
								}
							} else {
								if values > maxKey2 {
									maxKey2 = values
								}
								if values < minKey2 {
									minKey2 = values
								}
							}
						} else {
							newData := tk.Div(values, jumCat)
							datas = append(datas, newData)
							if i == 0 {
								if values > maxKey1 {
									maxKey1 = values
								}
							} else {
								if values > maxKey2 {
									maxKey2 = values
								}
								if values < minKey2 {
									minKey2 = values
								}
							}
						}
					}
				} else if strings.Contains(breakDown, "monthid") {
					if i == 0 { /* jika menjadi Key 1 saja */
						categories = append(categories, time.Month(val.GetInt("monthno")).String())
					}
					catTitle = "Month"

					if listCount == len(list)-1 {
						for _, monthSec := range listOfMonthPVal {
							values = pValueMonth[monthSec]
							if strings.Contains(key, "PLF") {
								datas = append(datas, values)
								if i == 0 {
									if values > maxKey1 {
										maxKey1 = values
									}
								} else {
									if values > maxKey2 {
										maxKey2 = values
									}
									if values < minKey2 {
										minKey2 = values
									}
								}
							} else {
								datas = append(datas, values)
								if i == 0 {
									if values > maxKey1 {
										maxKey1 = values
									}
								} else {
									if values > maxKey2 {
										maxKey2 = values
									}
									if values < minKey2 {
										minKey2 = values
									}
								}
							}
						}
					}
				} else if strings.Contains(breakDown, "year") {
					if listCount == 0 {
						for _, year := range listOfYears {
							categories = append(categories, tk.ToString(year))
						}
						catTitle = "Year"
					}
					if listCount == len(list)-1 { /* append hanya jika akumulasi terakhir saja */
						for _, year := range listOfYears {
							monthCumm := 0.0
							listBulan := listMonthPerYearPVal[year] /* bulan apa aja untuk tahun ini */
							for _, bulan := range listBulan {
								monthCumm += pValueMonth[bulan] /* total data per bulan selama 1 tahun */
							}
							if strings.Contains(key, "PLF") {
								monthCumm = tk.Div(monthCumm, tk.ToFloat64(len(listBulan), 0, tk.RoundingAuto))
							}

							datas = append(datas, monthCumm)

							if i == 0 {
								if monthCumm > maxKey1 {
									maxKey1 = monthCumm
								}
							} else {
								if monthCumm > maxKey2 {
									maxKey2 = monthCumm
								}
								if monthCumm < minKey2 {
									minKey2 = monthCumm
								}
							}
						}
					}
				} else if strings.Contains(breakDown, "project") {
					categories = append(categories, projectName)
					catTitle = "Project"
				} else if strings.Contains(breakDown, "turbine") {
					catTitle = "Turbine"

					jumCat := tk.ToFloat64(turbineCount, 0, tk.RoundingAuto)
					turbineDiv := tk.ToFloat64(totalTurbinePValue, 0, tk.RoundingAuto)
					if listCount == len(list)-1 {
						for jc := 1.0; jc <= jumCat; jc++ {
							if strings.Contains(key, "PLF") {
								newValues := tk.Div(values, tk.ToFloat64(durationMonths, 0, tk.RoundingAuto))
								datas = append(datas, newValues)
								if i == 0 {
									if newValues > maxKey1 {
										maxKey1 = newValues
									}
								} else {
									if newValues > maxKey2 {
										maxKey2 = newValues
									}
									if newValues < minKey2 {
										minKey2 = newValues
									}
								}
							} else {
								/* menggunakan values hasil akumulasi PER HARI */
								newData := tk.Div(values, turbineDiv)
								datas = append(datas, newData)
								if i == 0 {
									if newData > maxKey1 {
										maxKey1 = newData
									}
								} else {
									if newData > maxKey2 {
										maxKey2 = newData
									}
									if newData < minKey2 {
										minKey2 = newData
									}
								}
							}
						}
					}
				}
			} else {
				id := val.Get("_id")
				id1 := id.(tk.M).Get("id1")
				// tk.Printf("id1 : %s \n", id1)
				if strings.Contains(breakDown, "dateid") {
					dt := id1.(time.Time)
					monthList.Set(dt.Month().String(), 1)
					categories = append(categories, tk.ToString(dt.Day()))
					count := 0
					for field := range monthList {
						if count == 0 {
							catTitle = field
						} else {
							catTitle += " - " + field
						}
						count++
					}
					if len(monthList) == 1 {
						catTitle += " " + tk.ToString(dt.Year())
					} else {
						catTitle += " (" + tk.ToString(dt.Year()) + ")"
					}
				} else if strings.Contains(breakDown, "monthid") {
					if listCount == len(list)-1 {
						datas = []float64{}
						for _, monthSec := range listOfMonths {
							values = pValueMonth[monthSec]
							if strings.Contains(key, "PLF") {
								datas = append(datas, values)
								if i == 0 {
									if values > maxKey1 {
										maxKey1 = values
									}
								} else {
									if values > maxKey2 {
										maxKey2 = values
									}
									if values < minKey2 {
										minKey2 = values
									}
								}
							} else {
								datas = append(datas, values)
								if i == 0 {
									if values > maxKey1 {
										maxKey1 = values
									}
								} else {
									if values > maxKey2 {
										maxKey2 = values
									}
									if values < minKey2 {
										minKey2 = values
									}
								}
							}
						}
					}
					catTitle = "Month"
				} else if strings.Contains(breakDown, "year") {
					categories = append(categories, tk.ToString(id1))
					catTitle = "Year"
				} else if strings.Contains(breakDown, "project") {
					categories = append(categories, tk.ToString(id1))
					catTitle = "Project"
				} else if strings.Contains(breakDown, "turbine") {
					categories = append(categories, tk.ToString(id1))
					catTitle = "Turbine"
				}
			}

		}

		if len(datas) > 0 {
			series.Set("data", datas)
		}
		dataSeries = append(dataSeries, series)
	}
	minKey2 = 0.0 /* pake 0 aja deh minimum nya karena hampir gak mungkin generation atau plf di bawah 0 */
	if maxKey1 > maxKey2 {
		if measurement == "MWh" { /* generation */
			penambah := maxMinValue(maxKey1, 1.0)
			maxKey1 += penambah /* biar gak terlalu mentok ujung chart plotting nya */
		} else { /* plf which is selalu percentage */
			maxKey1 += 5 /* hanya ditambah 1 persen biar gak mentok chart */
		}
		maxKey2 = maxKey1
	} else {
		if measurement == "MWh" { /* generation */
			penambah := maxMinValue(maxKey2, 1.0)
			maxKey2 += penambah /* biar gak terlalu mentok ujung chart plotting nya */
		} else { /* plf which is selalu percentage */
			maxKey2 += 5 /* hanya ditambah 1 persen biar gak mentok chart */
		}
		maxKey1 = maxKey2
	}

	result := struct {
		Series     []tk.M
		Categories []string
		MinKey1    int
		MaxKey1    int
		MinKey2    int
		MaxKey2    int
		CatTitle   string
	}{
		Series:     dataSeries,
		Categories: categories,
		MinKey1:    0,
		MaxKey1:    tk.ToInt(maxKey1, tk.RoundingAuto), //tk.ToInt((maxKey1*2 - (maxKey1 / 4)), tk.RoundingAuto),
		MinKey2:    tk.ToInt(minKey2, tk.RoundingAuto),
		MaxKey2:    tk.ToInt(maxKey2, tk.RoundingAuto),
		CatTitle:   catTitle,
	}

	return helper.CreateResult(true, result, "success")
}

func maxMinValue(value float64, pengali float64) float64 {
	result := 0.0
	if value < 10 {
		result = 1.0 * pengali
	} else if value < 100 {
		result = 10.0 * pengali
	} else if value < 1000 {
		result = 100.0 * pengali
	} else if value < 10000 {
		result = 1000.0 * pengali
	} else if value < 100000 {
		result = 10000.0 * pengali
	} else if value < 1000000 {
		result = 100000.0 * pengali
	}

	return result
}
