package controller

import (
	. "github.com/eaciit/windapp/library/core"
	. "github.com/eaciit/windapp/library/models"
	"github.com/eaciit/windapp/web/helper"
	"time"

	"github.com/eaciit/knot/knot.v1"
	tk "github.com/eaciit/toolkit"
)

type AnalyticAvailabilityController struct {
	App
}

func CreateAnalyticAvailabilityController() *AnalyticAvailabilityController {
	var controller = new(AnalyticAvailabilityController)
	return controller
}

func (m *AnalyticAvailabilityController) GetData(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	var (
		pipes           []tk.M
		list            []tk.M
		dataSeriesAvail []tk.M
		dataSeriesProd  []tk.M
	)

	p := new(PayloadAnalytic)
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
	turbine := p.Turbine
	project := p.Project
	breakDown := p.BreakDown
	duration := tEnd.Sub(tStart).Hours() / 24

	match := tk.M{}
	match.Set("dateinfo.dateid", tk.M{"$gte": tStart, "$lte": tEnd})

	if len(turbine) > 0 {
		match.Set("turbine", tk.M{"$in": turbine})
	}

	group := tk.M{
		"power":           tk.M{"$sum": "$power"},
		"machinedowntime": tk.M{"$sum": "$machinedowntime"},
		"griddowntime":    tk.M{"$sum": "$griddowntime"},
		"oktime":          tk.M{"$sum": "$oktime"},
		"powerlost":       tk.M{"$sum": "$powerlost"},
		"totaltimestamp":  tk.M{"$sum": 1},
		"available":       tk.M{"$sum": "$available"},
		"minutes":         tk.M{"$sum": "$minutes"},
	}

	if project != "" {
		match.Set("projectname", project)
	}

	if breakDown == "Date" {
		group.Set("_id", tk.M{"id1": "$dateinfo.dateid"})
	} else if breakDown == "Month" {
		group.Set("_id", tk.M{"id1": "$dateinfo.monthid", "id2": "$dateinfo.monthdesc"})
	} else if breakDown == "Year" {
		group.Set("_id", tk.M{"id1": "$dateinfo.year"})
	} else if breakDown == "Project" {
		group.Set("_id", tk.M{"id1": "$projectname"})
	} else if breakDown == "Turbine" {
		group.Set("_id", tk.M{"id1": "$turbine"})
	}

	pipes = append(pipes, tk.M{"$match": match})
	pipes = append(pipes, tk.M{"$group": group})
	pipes = append(pipes, tk.M{"$sort": tk.M{"_id.id1": 1}})

	csr, e := DB().Connection.NewQuery().
		From(new(ScadaData).TableName()).
		Command("pipe", pipes).
		Cursor(nil)

	// add by ams, 2016-10-07
	defer csr.Close()

	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	e = csr.Fetch(&list, 0, false)

	keys := []string{"Production", "Machine Availability", "Grid Availability", "Total Availability"} //, "Data Availability", "PLF"

	categories := []string{}

	var totalEnergy float64

	for _, val := range list {
		power := val.GetFloat64("power") / 1000.0
		totalEnergy += power / 6
	}

	max := 0
	min := 0

	for _, key := range keys {
		series := tk.M{}
		//seriesProd := tk.M{}

		series.Set("name", key)
		series.Set("type", "column")
		series.Set("axis", "availpercentage")

		if key == "Production" {
			//series.Set("type", "line")
			//series.Set("dashType", "solid")
			//series.Set("markers", tk.M{"visible": false})
			//series.Set("width", 2)
			series.Set("axis", "availline")
		}

		var datas []float64
		for _, val := range list {
			var plf, trueAvail, machineAvail, gridAvail, dataAvail, prod float64
			var totalTurbine float64

			minutesInHour := val.GetFloat64("minutes") / 60.0
			okTime := val.GetFloat64("oktime")
			power := val.GetFloat64("power") / 1000.0
			energy := power / 6

			mDownTime := val.GetFloat64("machinedowntime") / 3600.0
			gDownTime := val.GetFloat64("griddowntime") / 3600.0
			sumTimeStamp := val.GetFloat64("totaltimestamp")

			/*if breakDown == "Date" {}
			if breakDown == "Month" {}
			if breakDown == "Year" {}
			if breakDown == "Project" {}*/
			if breakDown == "Turbine" || breakDown == "Date" || breakDown == "Month" || breakDown == "Year" {
				totalTurbine = 1.0
			} else {
				totalTurbine = tk.ToFloat64(len(turbine), 0, tk.RoundingAuto)
			}

			/*plf = energy / (totalTurbine * duration * 24.0 * 2100) * 100 * 1000
			trueAvail = (okTime / 3600) / (duration * totalTurbine * 24.0) * 100
			machineAvail = (minutesInHour - mDownTime) / (totalTurbine * 24.0 * duration) * 100
			gridAvail = (minutesInHour - gDownTime) / (totalTurbine * 24.0 * duration) * 100
			dataAvail = sumTimeStamp / (144 * duration * totalTurbine) * 100
			// prod = energy / totalEnergy * 100
			prod = energy*/

			plf = energy / (totalTurbine * minutesInHour * 2100) * 100 * 1000
			trueAvail = (okTime / 3600) / (totalTurbine * minutesInHour) * 100
			machineAvail = (minutesInHour - mDownTime) / (totalTurbine * minutesInHour) * 100
			gridAvail = (minutesInHour - gDownTime) / (totalTurbine * minutesInHour) * 100
			dataAvail = (sumTimeStamp * 10 / 60) / (minutesInHour * totalTurbine) * 100
			prod = energy

			_ = duration

			if key == "Machine Availability" {
				datas = append(datas, tk.ToFloat64(machineAvail, 2, tk.RoundingAuto))
				val := tk.ToInt(tk.ToFloat64(machineAvail, 2, tk.RoundingAuto), tk.RoundingUp)
				if val > max {
					max = val
				}
			} else if key == "Grid Availability" {
				datas = append(datas, tk.ToFloat64(gridAvail, 2, tk.RoundingAuto))
				val := tk.ToInt(tk.ToFloat64(gridAvail, 2, tk.RoundingAuto), tk.RoundingUp)
				if val > max {
					max = val
				}
			} else if key == "Total Availability" {
				datas = append(datas, tk.ToFloat64(trueAvail, 2, tk.RoundingAuto))
				val := tk.ToInt(tk.ToFloat64(trueAvail, 2, tk.RoundingAuto), tk.RoundingUp)
				if val > max {
					max = val
				}
			} else if key == "Data Availability" {
				datas = append(datas, tk.ToFloat64(dataAvail, 2, tk.RoundingAuto))
				val := tk.ToInt(tk.ToFloat64(dataAvail, 2, tk.RoundingAuto), tk.RoundingUp)
				if val > max {
					max = val
				}
			} else if key == "Production" {
				datas = append(datas, tk.ToFloat64(prod, 2, tk.RoundingAuto))
				val := tk.ToInt(tk.ToFloat64(prod, 2, tk.RoundingAuto), tk.RoundingUp)
				if val > min {
					min = val
				}
			} else if key == "PLF" {
				datas = append(datas, tk.ToFloat64(plf, 2, tk.RoundingAuto))
				val := tk.ToInt(tk.ToFloat64(plf, 2, tk.RoundingAuto), tk.RoundingUp)
				if val > max {
					max = val
				}
			}
		}

		if len(datas) > 0 {
			series.Set("data", datas)
		}

		if key == "Production" {
			dataSeriesProd = append(dataSeriesProd, series)
		} else {
			dataSeriesAvail = append(dataSeriesAvail, series)
		}
	}

	for _, val := range list {
		id := val.Get("_id")
		id1 := id.(tk.M).Get("id1")

		if breakDown == "Date" {
			dt := id1.(time.Time)
			categories = append(categories, tk.ToString(dt.Day())+"/"+dt.Month().String()[:3])
		} else if breakDown == "Month" {
			id2 := id.(tk.M).GetString("id2")
			if id2 != "" {
				categories = append(categories, id2)
			}
		} else if breakDown == "Year" {
			categories = append(categories, tk.ToString(id1))
		} else if breakDown == "Project" {
			categories = append(categories, tk.ToString(id1))
		} else if breakDown == "Turbine" {
			categories = append(categories, tk.ToString(id1))
		}
	}

	result := struct {
		SeriesAvail []tk.M
		SeriesProd  []tk.M
		Categories  []string
		Max         int
		Min         int
	}{
		SeriesAvail: dataSeriesAvail,
		SeriesProd:  dataSeriesProd,
		Categories:  categories,
		Max:         max * 2,
		Min:         min * -3,
	}

	return helper.CreateResult(true, result, "success")
}
