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
	for yearDay, data := range monthDay {
		days := data.(tk.M).GetFloat64("days")
		totalInMonth := data.(tk.M).GetFloat64("totalInMonth")
		if tk.ToInt(yearDay[4:], tk.RoundingAuto) == monthno {
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
	turbineCount := p.Misc.GetInt("totalturbine")
	if turbineCount == 0 {
		turbineCount = 24
	}
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
	listOfCategories := map[string][]string{}
	listOfCatTitles := map[string]string{}
	for i := tStart.Year(); i <= tEnd.Year(); i++ {
		listOfYears = append(listOfYears, i)
	}
	// totalTurbine := 1.0
	// if !strings.Contains(breakDown, "turbine") {
	totalTurbine := tk.ToFloat64(turbineCount, 0, tk.RoundingAuto)
	// }

	for i, key := range keys {
		list = []tk.M{}
		if tk.HasMember(keyList, key) {
			csrPValue, e := DB().Connection.NewQuery().
				From(new(ExpPValueModel).TableName()).
				Where(dbox.In("monthno", months...)).
				Cursor(nil)
			if e != nil {
				return helper.CreateResult(false, nil, e.Error())
			}
			e = csrPValue.Fetch(&list, 0, false)
			// add by ams, 2016-10-07
			csrPValue.Close()
			if e != nil {
				return helper.CreateResult(false, nil, e.Error())
			}
		} else {
			p.Misc.Set("knot_data", k)
			filter, e := p.ParseFilter()
			if e != nil {
				return helper.CreateResult(false, nil, e.Error())
			}
			fb := DB().Connection.Fb()
			fb.AddFilter(dbox.And(filter...))
			matches, e := fb.Build()

			group := tk.M{
				"power":           tk.M{"$sum": "$power"},
				"energy":          tk.M{"$sum": "$energy"},
				"machinedowntime": tk.M{"$sum": "$machinedowntime"},
				"griddowntime":    tk.M{"$sum": "$griddowntime"},
				"oktime":          tk.M{"$sum": "$oktime"},
				"totaltime":       tk.M{"$sum": "$totaltime"},
				"minutes":         tk.M{"$sum": "$minutes"},
				"countdata":       tk.M{"$sum": 1},
				"maxdate":         tk.M{"$max": "$dateinfo.dateid"},
				"mindate":         tk.M{"$min": "$dateinfo.dateid"},
			}

			group.Set("_id", tk.M{"id1": breakDown})
			if strings.Contains(breakDown, "month") {
				group.Set("_id", tk.M{"id1": breakDown, "id2": "$dateinfo.monthdesc"})
			}

			pipes := []tk.M{{"$match": matches}, {"$group": group}, {"$sort": tk.M{"_id.id1": 1}}}

			csr, e := DB().Connection.NewQuery().
				From(new(ScadaData).TableName()).
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
			// tk.Printf("breakDown : %s \n", breakDown)
		}
		series := tk.M{}
		measurement = "%"
		if i == 0 {
			series.Set("name", key)
			series.Set("type", "column")
			series.Set("axis", "Key1")
			series.Set("gap", 0.7)
			if key == "Actual Production" || strings.Contains(key, "Generation") {
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
			if key == "Actual Production" || strings.Contains(key, "Generation") {
				minKey2 = 99999999.99
				measurement = "MWh"
			}
			series.Set("satuan", measurement)
		}

		var datas []float64
		var values float64
		categories = []string{}
		for listCount, val := range list {
			var hourValue, minutes float64

			if strings.Contains(breakDown, "dateid") {
				id := val.Get("_id").(tk.M)
				id1 := id.Get("id1").(time.Time)
				// hourValue, minutes = getHourMinute(id1.UTC(), id1.UTC(), val.Get("mindate").(time.Time), val.Get("maxdate").(time.Time), val.GetFloat64("minutes"))

				hourValue = helper.GetHourValue(id1.UTC(), id1.UTC(), val.Get("mindate").(time.Time), val.Get("maxdate").(time.Time))

			} else {
				// hourValue, minutes = getHourMinute(tStart.UTC(), tEnd.UTC(), val.Get("mindate").(time.Time), val.Get("maxdate").(time.Time), val.GetFloat64("minutes"))

				hourValue = helper.GetHourValue(tStart.UTC(), tEnd.UTC(), val.Get("mindate").(time.Time), val.Get("maxdate").(time.Time))
			}

			oktime := val.GetFloat64("oktime")
			energy := val.GetFloat64("energy") / 1000
			mDownTime := (val.GetFloat64("machinedowntime") / 3600.0)
			gDownTime := (val.GetFloat64("griddowntime") / 3600.0)
			sumTimeStamp := val.GetFloat64("countdata")
			minutes = val.GetFloat64("minutes") / 60

			machineAvail, gridAvail, dataAvail, trueAvail, plf := helper.GetAvailAndPLF(totalTurbine, oktime, energy, mDownTime, gDownTime, sumTimeStamp, hourValue, minutes)

			// log.Printf("%v | %v | %v | %v | %v | %v | %v | %v \n", totalTurbine, oktime, energy, mDownTime, gDownTime, sumTimeStamp, hourValue, minutes)

			switch key {
			case "Machine Availability":
				values = machineAvail //tk.Div((minutes-(val.GetFloat64("machinedowntime")/3600.0)), (totalTurbine*hourValue)) * 100 /*percentage*/
			case "Grid Availability":
				values = gridAvail //tk.Div((minutes-(val.GetFloat64("griddowntime")/3600.0)), (totalTurbine*hourValue)) * 100 /*percentage*/
			case "Total Availability":
				values = trueAvail //tk.Div((val.GetFloat64("oktime")/3600), (totalTurbine*hourValue)) * 100
			case "Data Availability":
				values = dataAvail //tk.Div((tk.ToFloat64((val.GetInt("countdata")*10/60), 6, tk.RoundingAuto)), (hourValue*totalTurbine)) * 100
			case "Actual PLF":
				values = plf //tk.Div((val.GetFloat64("energy")/1000), (hourValue*2.1*totalTurbine)) * 100
			case "Actual Production":
				values = val.GetFloat64("energy") / 1000 /*MWh*/
			case "P50 Generation":
				values += checkPValue(monthDay, val.GetFloat64("p50netgenmwh"), val.GetInt("monthno"))
			case "P50 PLF":
				values += val.GetFloat64("p50plf")
			case "P75 Generation":
				values += checkPValue(monthDay, val.GetFloat64("p75netgenmwh"), val.GetInt("monthno"))
			case "P75 PLF":
				values += val.GetFloat64("p75plf")
			case "P90 Generation":
				values += checkPValue(monthDay, val.GetFloat64("p90netgenmwh"), val.GetInt("monthno"))
			case "P90 PLF":
				values += val.GetFloat64("p90plf")
			}

			// plf = energy / (totalTurbine * hourValue * 2.1) * 100
			// trueAvail = (okTime / 3600) / (totalTurbine * hourValue) * 100
			// dataAvail = (sumTimeStamp * 10 / 60) / (hourValue * totalTurbine) * 100

			/*p50netgen per day = (p50netgenmwh / jumlah hari dalam bulan tersebut) * jumlah hari periode
			plf e => p50netgen per hari ne / ( 2.1 x jumlah hari periode x 24 x 24 )*/

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

			if tk.HasMember(keyList, key) {
				if strings.Contains(breakDown, "dateid") {
					datas = []float64{}

					if listCount == 0 { /*bulan pertama*/
						catTitle = tStart.Month().String()
						if len(list) == 1 { /*jika hanya 1 bulan*/
							for iDate := startdate; iDate <= enddate; iDate++ {
								categories = append(categories, tk.ToString(iDate))
							}
							catTitle += " " + tk.ToString(listOfYears[0]) /*Dec 2015*/
						} else { /*jika lebih dari 1 bulan*/
							month := val.GetInt("monthno")
							maxDays := monthDay.Get(tk.ToString(tStart.Year()) + tk.ToString(month)).(tk.M).GetInt("totalInMonth")
							for iDate := startdate; iDate <= maxDays; iDate++ {
								categories = append(categories, tk.ToString(iDate))
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
						}

					}
					jumCat := tk.ToFloat64(len(categories), 6, tk.RoundingAuto)
					// tk.Printf("key : %s \n", key)
					for iCat := range categories {
						_ = iCat
						if strings.Contains(key, "PLF") {
							values = tk.Div(values, tk.ToFloat64(durationMonths, 0, tk.RoundingAuto))
							datas = append(datas, values)

							if i == 0 {
								maxKey1 = values
							} else {
								maxKey2 = values
								minKey2 = values
							}
						} else {
							datas = append(datas, tk.Div(values, jumCat))
							if i == 0 {
								maxKey1 = tk.Div(values, jumCat)
							} else {
								maxKey2 = tk.Div(values, jumCat)
								minKey2 = tk.Div(values, jumCat)
							}
						}
					}
				} else if strings.Contains(breakDown, "monthid") {
					categories = append(categories, time.Month(val.GetInt("monthno")).String())
					catTitle = "Month"
				} else if strings.Contains(breakDown, "year") {
					if listCount == 0 {
						for _, year := range listOfYears {
							categories = append(categories, tk.ToString(year))
						}
						catTitle = "Year"
					}
				} else if strings.Contains(breakDown, "project") {
					categories = append(categories, "Tejuva")
					catTitle = "Project"
				} else if strings.Contains(breakDown, "turbine") {
					temp := p.Filter.Filters[2].Value.([]interface{})
					for _, turbine := range temp {
						categories = append(categories, turbine.(string))
					}
					catTitle = "Turbine"
				}
				listOfCategories["pvalue"] = categories
				listOfCatTitles["pvalue"] = catTitle
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
					id2 := id.(tk.M).GetString("id2")
					if id2 != "" {
						categories = append(categories, id2)
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
				listOfCategories["biasa"] = categories
				listOfCatTitles["biasa"] = catTitle
			}
		}
		if i > 0 {
			if measurement == "MWh" {
				penambah := maxMinValue(maxKey2, 1.0)
				pengurang := maxMinValue(minKey2, 2.0)

				maxKey2 += penambah
				minKey2 -= pengurang
			} else {
				maxKey2 += 1
				minKey2 -= 5
			}
		}

		if len(datas) > 0 {
			series.Set("data", datas)
		}
		dataSeries = append(dataSeries, series)
	}
	categories = []string{}
	catTitle = ""
	for key, value := range listOfCategories {
		if key == "pvalue" {
			categories = value
			break
		} else {
			categories = value
		}
	}
	for key, value := range listOfCatTitles {
		if key == "pvalue" {
			catTitle = value
			break
		} else {
			catTitle = value
		}
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
		MaxKey1:    tk.ToInt((maxKey1*2 - (maxKey1 / 4)), tk.RoundingAuto),
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
