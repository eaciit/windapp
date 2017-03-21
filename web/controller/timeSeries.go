package controller

import (
	"bufio"
	. "eaciit/wfdemo-git/library/core"
	. "eaciit/wfdemo-git/library/models"
	"eaciit/wfdemo-git/web/helper"
	"encoding/csv"
	"io"
	"log"
	"math"
	"os"
	"strings"

	"time"

	"github.com/eaciit/knot/knot.v1"
	tk "github.com/eaciit/toolkit"
)

type TimeSeriesController struct {
	App
}

type HFDModel struct {
	Timestamp time.Time
	Turbine   string
	Tag       string
	Value     float64
}

type ResDataAvail struct {
	Chart      []tk.M
	PeriodList []tk.M
	Breaks     []tk.M
}

type MappingColumn struct {
	Name     string
	SecField string
	Unit     string
	MinValue float64
	MaxValue float64
}

func CreateTimeSeriesController() *TimeSeriesController {
	var controller = new(TimeSeriesController)
	return controller
}

// func (m *TimeSeriesController) GetData(k *knot.WebContext) interface{} {
// 	k.Config.OutputType = knot.OutputJson

// 	var (
// 		pipes       []tk.M
// 		list        []tk.M
// 		resultChart []tk.M
// 		periodList  []tk.M
// 	)

// 	p := new(PayloadTimeSeries)
// 	e := k.GetPayload(&p)
// 	if e != nil {
// 		return helper.CreateResult(false, nil, e.Error())
// 	}

// 	tStart, tEnd, e := helper.GetStartEndDate(k, p.Period, p.DateStart, p.DateEnd)
// 	if e != nil {
// 		return helper.CreateResult(false, nil, e.Error())
// 	}

// 	match := tk.M{}

// 	match.Set("dateinfo.dateid", tk.M{"$gte": tStart, "$lte": tEnd})
// 	match.Set("avgwindspeed", tk.M{"$lte": 25})
// 	// match.Set("turbine", p.Turbine)

// 	if p.Project != "" {
// 		anProject := strings.Split(p.Project, "(")
// 		match.Set("projectname", strings.TrimRight(anProject[0], " "))
// 	}

// 	group := tk.M{
// 		"energy":    tk.M{"$sum": "$energy"},
// 		"windspeed": tk.M{"$avg": "$avgwindspeed"},
// 	}

// 	group.Set("_id", "$timestamp")

// 	pipes = append(pipes, tk.M{"$match": match})
// 	pipes = append(pipes, tk.M{"$group": group})
// 	pipes = append(pipes, tk.M{"$sort": tk.M{"_id": 1}})

// 	csr, e := DB().Connection.NewQuery().
// 		From(new(ScadaData).TableName()).
// 		Command("pipe", pipes).
// 		Cursor(nil)

// 	if e != nil {
// 		return helper.CreateResult(false, nil, e.Error())
// 	}

// 	e = csr.Fetch(&list, 0, false)

// 	csr.Close()

// 	var dtProd, dtWS [][]interface{}

// 	for _, val := range list {
// 		// intTimestamp :=
// 		intTimestamp := tk.ToInt(tk.ToString(val.Get("_id").(time.Time).UTC().Unix())+"000", tk.RoundingAuto)

// 		energy := val.GetFloat64("energy") / 1000
// 		wind := val.GetFloat64("windspeed")

// 		dtProd = append(dtProd, []interface{}{intTimestamp, energy})
// 		dtWS = append(dtWS, []interface{}{intTimestamp, wind})
// 	}

// 	resultChart = append(resultChart, tk.M{"name": "Production", "data": dtProd, "unit": "MWh"})
// 	resultChart = append(resultChart, tk.M{"name": "Windspeed", "data": dtWS, "unit": "m/s"})

// 	data := struct {
// 		Data ResDataAvail
// 	}{
// 		Data: ResDataAvail{
// 			Chart:      resultChart,
// 			PeriodList: periodList,
// 		},
// 	}

// 	return helper.CreateResult(true, data, "success")
// }

func (m *TimeSeriesController) GetAvailDate(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	return helper.CreateResult(true, k.Session("availdate", ""), "success")
}

func (m *TimeSeriesController) GetDataHFD(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	breaks := []tk.M{}
	resultChart := []tk.M{}
	periodList := []tk.M{}

	p := new(PayloadTimeSeries)
	e := k.GetPayload(&p)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	var tStart, tEnd time.Time

	if p.IsHour {
		tStart = p.DateStart.UTC()
		tEnd = p.DateEnd.UTC()
	} else {
		tStart, tEnd, e = helper.GetStartEndDate(k, p.Period, p.DateStart, p.DateEnd)
	}

	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	projectName := ""
	_ = projectName

	if p.Project != "" {
		anProject := strings.Split(p.Project, "(")
		projectName = strings.TrimRight(anProject[0], " ")
	}

	turbine := ""

	if len(p.Turbine) == 1 {
		turbine = p.Turbine[0].(string)
	}

	dataType := p.DataType
	pageType := p.PageType

	// log.Printf(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> %v | %v \n", dataType, pageType)

	// default tags
	tags := []string{}
	tags = []string{"windspeed", "power"}

	mapField := map[string]MappingColumn{}
	mapField["windspeed"] = MappingColumn{"Wind Speed", "WindSpeed_ms", "m/s", 0.0, 25.0}
	mapField["power"] = MappingColumn{"Power", "ActivePower_kW", "kW", -100, 2121.0}
	mapField["production"] = MappingColumn{"Production", "", "kWh", -100, 2121.0}
	mapField["winddirection"] = MappingColumn{"Wind Direction", "WindDirection", "Degree", -10.0, 120.0}
	mapField["nacellepos"] = MappingColumn{"Nacelle Direction", "NacellePos", "Degree", -10.0, 120.0}
	mapField["rotorrpm"] = MappingColumn{"Rotor RPM", "RotorSpeed_RPM", "RPM", 0.0, 30.0}
	mapField["genrpm"] = MappingColumn{"Generator RPM", "WindSpeed_ms", "RPM", 0.0, 30.0}
	mapField["pitchangle"] = MappingColumn{"Pitch Angle", "PitchAngle1", "Degree", -10.0, 120.0}

	if len(p.TagList) > 0 {
		tags = p.TagList
	}

	if pageType == "HFD" && dataType == "SEC" {
		tmpTag := []string{}
		for _, tg := range tags {
			tmpTag = append(tmpTag, mapField[tg].SecField)
		}

		tags = tmpTag

		for {
			if tStart.UTC().Sub(tEnd.UTC()).Seconds() >= 0 {
				break
			}

			before := tStart.UTC()
			tStart = tStart.UTC().Add(time.Duration(24) * time.Hour)

			beforeInt := tk.ToInt(tk.ToString(before.UTC().Unix())+"000", tk.RoundingAuto)
			afterInt := tk.ToInt(tk.ToString(tStart.UTC().Unix())+"000", tk.RoundingAuto)

			periodList = append(periodList, tk.M{"starttime": before.UTC(), "endtime": tStart.UTC(), "starttimeint": beforeInt, "endtimeint": afterInt})
		}

		if len(periodList) > 0 || p.IsHour {
			for idx, pl := range periodList {
				current := pl
				currStar := current.Get("starttime").(time.Time)
				currEnd := current.Get("endtime").(time.Time)

				hfds, empty, e := GetHFDData(turbine, currStar, currEnd, tags)

				breaks = append(breaks, empty...)

				if e != nil {
					return helper.CreateResult(false, nil, e.Error())
				}

				if len(hfds) > 0 || p.IsHour {
					for _, tag := range p.TagList {
						var dts [][]interface{}
						var dterr [][]interface{}
						columnTag := mapField[tag]
						for _, val := range hfds {
							timestamp := tk.ToInt(tk.ToString(val.Get("timestamp").(time.Time).Unix())+"000", tk.RoundingAuto)
							tagVal := val.GetFloat64(columnTag.SecField)
							dt := []interface{}{timestamp, tagVal}
							dts = append(dts, dt)

							if tagVal < columnTag.MinValue || tagVal > columnTag.MaxValue {
								// dte := []interface{}{timestamp, tagVal}
								// dterr = append(dterr, tk.M{"x": timestamp})
								dterr = append(dterr, []interface{}{timestamp, 100.0})
							}
						}

						resultChart = append(resultChart, tk.M{"name": mapField[tag].Name, "data": dts, "dataerr": dterr, "unit": mapField[tag].Unit, "minval": mapField[tag].MinValue, "maxval": mapField[tag].MaxValue})
					}

					periodList = periodList[idx:]
					break
				}

			}
		}
	} else {
		match := tk.M{}
		group := tk.M{}
		pipes := []tk.M{}
		var collName string

		if projectName != "" {
			match.Set("projectname", projectName)
		}

		if pageType == "HFD" {
			collName = new(ScadaDataHFD).TableName()
			match.Set("dateinfo.dateid", tk.M{"$gte": tStart, "$lte": tEnd})
			// match.Set("fast_windspeed_ms_stddev", tk.M{"$lte": 25})
			match.Set("turbine", turbine)

			group = tk.M{
				"_id": "$timestamp",
				// "energy":    tk.M{"$sum": "$energy"},
				"windspeed":     tk.M{"$avg": "$fast_windspeed_ms"},
				"power":         tk.M{"$sum": "$fast_activepower_kw"},
				"winddirection": tk.M{"$avg": "$slow_winddirection"},
				"nacellepos":    tk.M{"$avg": "$slow_nacellepos"},
				"rotorrpm":      tk.M{"$avg": "$fast_rotorspeed_rpm"},
				"genrpm":        tk.M{"$avg": "$fast_genspeed_rpm"},
				"pitchangle":    tk.M{"$avg": "$fast_pitchangle"},
			}
		} else if pageType == "LIVE" {
			collName = new(ScadaRealTime).TableName()
			if tStart.Year() != 1 && tEnd.Year() != 1 {
				match.Set("dateinfo.dateid", tk.M{"$gte": tStart, "$lte": tEnd})
			}

			// match.Set("windspeed", tk.M{"$lte": 25})
			match.Set("turbine", turbine)

			group = tk.M{
				"_id": "$timestamp",
				// "energy":    tk.M{"$sum": "$energy"},
				"windspeed":     tk.M{"$avg": "$windspeed"},
				"power":         tk.M{"$sum": "$activepower"},
				"winddirection": tk.M{"$avg": "$winddirection"},
				"nacellepos":    tk.M{"$avg": "$nacelleposition"},
				"rotorrpm":      tk.M{"$avg": "$rotorrpm"},
				"pitchangle":    tk.M{"$avg": "$pitchangle"},
			}
		} else {
			collName = new(ScadaDataOEM).TableName()
			match.Set("dateinfo.dateid", tk.M{"$gte": tStart, "$lte": tEnd})
			// match.Set("denwindspeed", tk.M{"$lte": 25})
			match.Set("turbine", turbine)

			group = tk.M{
				"_id":        "$timestamp",
				"power":      tk.M{"$sum": "$denpower"},
				"windspeed":  tk.M{"$avg": "$denwindspeed"},
				"production": tk.M{"$avg": "$energy"},
			}
		}

		pipes = append(pipes, tk.M{"$match": match})
		pipes = append(pipes, tk.M{"$group": group})

		if tStart.Year() != 1 && tEnd.Year() != 1 {
			pipes = append(pipes, tk.M{"$sort": tk.M{"_id": 1}})
		} else {
			pipes = append(pipes, tk.M{"$sort": tk.M{"_id": -1}})
			pipes = append(pipes, tk.M{"$limit": 5})
		}

		// log.Printf("%v \n", collName)
		// log.Printf("%v | %v \n", tStart.String(), tEnd.String())

		// for _, p := range pipes {
		// 	log.Printf(">>>>>> %#v \n", p)
		// }

		csr, e := DB().Connection.NewQuery().
			From(collName).
			Command("pipe", pipes).
			Cursor(nil)

		if e != nil {
			return helper.CreateResult(false, nil, e.Error())
		}

		list := []tk.M{}
		e = csr.Fetch(&list, 0, false)

		defer csr.Close()

		// log.Printf(">>> %v \n", len(list))

		for _, tag := range tags {
			var dts [][]interface{}
			var dterr [][]interface{}
			columnTag := mapField[tag]
			if len(list) > 0 {
				// get the time is not exist in collection
				if pageType != "LIVE" {
					_first := list[0]
					_firstTimestamp := _first.Get("_id").(time.Time).UTC()
					first := getTime10MinutesNotExist(tStart, _firstTimestamp)

					if len(first) > 0 {
						dts = append(dts, first...)
					}
				}

				for _, val := range list {
					timestamp := tk.ToInt(tk.ToString(val.Get("_id").(time.Time).Unix())+"000", tk.RoundingAuto)
					var tagVal float64

					// if tag == "production" {
					// 	tagVal = val.GetFloat64(tag)
					// } else if tag == "windspeed" {
					// 	tagVal = val.GetFloat64(tag)
					// if tagVal < 0 {
					// 	tagVal = 0.0
					// }
					// } else {
					tagVal = val.GetFloat64(tag)
					// }

					dt := []interface{}{timestamp, tagVal}
					if tagVal <= -99999.0 {
						// res := tk.M{}
						// res.Set("from", val.Get("_id").(time.Time).UTC())
						// res.Set("to", val.Get("_id").(time.Time).UTC().Add(10*time.Minute))
						// breaks = append(breaks, res)

						dt = []interface{}{timestamp, nil}
					}

					dts = append(dts, dt)

					if (tagVal < columnTag.MinValue || tagVal > columnTag.MaxValue) && tagVal > -99999.0 {
						dterr = append(dterr, []interface{}{timestamp, 100.0})
					}
				}

				// get the time is not exist in collection
				if pageType != "LIVE" {
					_last := list[len(list)-1]
					_lastTimestamp := _last.Get("_id").(time.Time).UTC()

					last := getTime10MinutesNotExist(_lastTimestamp, tEnd)

					if len(last) > 0 {
						dts = append(dts, last...)
					}
				}

				resultChart = append(resultChart, tk.M{"name": columnTag.Name, "data": dts, "dataerr": dterr, "unit": columnTag.Unit, "minval": columnTag.MinValue, "maxval": columnTag.MaxValue})
			}
		}

		csr.Close()
	}

	data := struct {
		Data ResDataAvail
	}{
		Data: ResDataAvail{
			Chart:      resultChart,
			PeriodList: periodList,
			Breaks:     breaks,
		},
	}

	return helper.CreateResult(true, data, "success")
}

func getTime10MinutesNotExist(start time.Time, end time.Time) (result [][]interface{}) {
	for {
		// log.Printf(">> %v | %v \n", start.String(), end.String())

		if end.UTC().Sub(start.UTC()).Minutes() <= 10 {
			// log.Println(">>>> BREAK")
			break
		}

		before := end.UTC()
		timestamp := tk.ToInt(tk.ToString(before.UTC().Unix())+"000", tk.RoundingAuto)
		result = append(result, []interface{}{timestamp, nil})

		start = start.UTC().Add(time.Duration(10) * time.Minute)
	}

	return
}

func (m *TimeSeriesController) GetDataHFDX(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	breaks := []tk.M{}
	resultChart := []tk.M{}
	periodList := []tk.M{}

	p := new(PayloadTimeSeries)
	e := k.GetPayload(&p)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	var tStart, tEnd time.Time

	if p.IsHour {
		tStart = p.DateStart.UTC()
		tEnd = p.DateEnd.UTC()
	} else {
		tStart, tEnd, e = helper.GetStartEndDate(k, p.Period, p.DateStart, p.DateEnd)
	}

	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	projectName := ""
	_ = projectName

	if p.Project != "" {
		anProject := strings.Split(p.Project, "(")
		projectName = strings.TrimRight(anProject[0], " ")
	}

	turbine := ""

	if len(p.Turbine) == 1 {
		turbine = p.Turbine[0].(string)
	}

	dataType := p.DataType
	pageType := p.PageType

	// log.Printf(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> %v | %v \n", dataType, pageType)

	// default tags
	tags := []string{}
	tags = []string{"windspeed", "power"}

	mapField := map[string]MappingColumn{}
	mapField["windspeed"] = MappingColumn{"Wind Speed", "WindSpeed_ms", "m/s", 0.0, 25.0}
	mapField["power"] = MappingColumn{"Power", "ActivePower_kW", "kW", 0, 500.0}
	mapField["production"] = MappingColumn{"Production", "", "kWh", 0, 1000.0}
	mapField["winddirection"] = MappingColumn{"Wind Direction", "WindDirection", "Degree", -10.0, 120.0}
	mapField["rotorrpm"] = MappingColumn{"Rotor RPM", "RotorSpeed_RPM", "RPM", 0.0, 30.0}
	mapField["genrpm"] = MappingColumn{"Generator RPM", "WindSpeed_ms", "RPM", 0.0, 30.0}
	mapField["pitchangle"] = MappingColumn{"Pitch Angle", "PitchAngle1", "Degree", -10.0, 120.0}

	if len(p.TagList) > 0 {
		tags = p.TagList
	}

	if pageType == "HFD" && dataType == "SEC" {
		tmpTag := []string{}
		for _, tg := range tags {
			tmpTag = append(tmpTag, mapField[tg].SecField)
		}

		tags = tmpTag

		for {
			if tStart.UTC().Sub(tEnd.UTC()).Seconds() >= 0 {
				break
			}

			before := tStart.UTC()
			tStart = tStart.UTC().Add(time.Duration(24) * time.Hour)

			beforeInt := tk.ToInt(tk.ToString(before.UTC().Unix())+"000", tk.RoundingAuto)
			afterInt := tk.ToInt(tk.ToString(tStart.UTC().Unix())+"000", tk.RoundingAuto)

			periodList = append(periodList, tk.M{"starttime": before.UTC(), "endtime": tStart.UTC(), "starttimeint": beforeInt, "endtimeint": afterInt})
		}

		if len(periodList) > 0 || p.IsHour {
			for idx, pl := range periodList {
				current := pl
				currStar := current.Get("starttime").(time.Time)
				currEnd := current.Get("endtime").(time.Time)

				hfds, empty, e := GetHFDData(turbine, currStar, currEnd, tags)

				breaks = append(breaks, empty...)

				if e != nil {
					return helper.CreateResult(false, nil, e.Error())
				}

				if len(hfds) > 0 || p.IsHour {
					for _, tag := range p.TagList {
						var dts [][]interface{}
						var dterr [][]interface{}
						columnTag := mapField[tag]
						for _, val := range hfds {
							timestamp := tk.ToInt(tk.ToString(val.Get("timestamp").(time.Time).Unix())+"000", tk.RoundingAuto)
							tagVal := val.GetFloat64(columnTag.SecField)
							dt := []interface{}{timestamp, tagVal}
							dts = append(dts, dt)

							if tagVal < columnTag.MinValue || tagVal > columnTag.MaxValue {
								// dte := []interface{}{timestamp, tagVal}
								// dterr = append(dterr, tk.M{"x": timestamp})
								dterr = append(dterr, []interface{}{timestamp, columnTag.MaxValue})
							}
						}

						resultChart = append(resultChart, tk.M{"name": mapField[tag].Name, "data": dts, "dataerr": dterr, "unit": mapField[tag].Unit, "minval": mapField[tag].MinValue, "maxval": mapField[tag].MaxValue})
					}

					periodList = periodList[idx:]
					break
				}

			}
		}
	} else {
		match := tk.M{}
		group := tk.M{}
		pipes := []tk.M{}
		var collName string

		if projectName != "" {
			match.Set("projectname", projectName)
		}

		if pageType == "HFD" {
			collName = new(ScadaDataHFD).TableName()
			match.Set("dateinfo.dateid", tk.M{"$gte": tStart, "$lte": tEnd})
			// match.Set("fast_windspeed_ms_stddev", tk.M{"$lte": 25})
			match.Set("turbine", turbine)

			group = tk.M{
				"_id": "$timestamp",
				// "energy":    tk.M{"$sum": "$energy"},
				"windspeed":     tk.M{"$avg": "$fast_windspeed_ms_stddev"},
				"power":         tk.M{"$sum": "$fast_activepower_kw_stddevv"},
				"winddirection": tk.M{"$avg": "$slow_winddirection_stddev"},
				"nacellepos":    tk.M{"$avg": "$slow_nacellepos_stddev"},
				"rotorrpm":      tk.M{"$avg": "$fast_rotorspeed_rpm_stddev"},
				"genrpm":        tk.M{"$avg": "$fast_genspeed_rpm_stddev"},
				"pitchangle":    tk.M{"$avg": "$fast_pitchangle_stddev"},
			}
		} else if pageType == "LIVE" {
			collName = new(ScadaRealTime).TableName()
			if tStart.Year() != 1 && tEnd.Year() != 1 {
				match.Set("dateinfo.dateid", tk.M{"$gte": tStart, "$lte": tEnd})
			}

			// match.Set("windspeed", tk.M{"$lte": 25})
			match.Set("turbine", turbine)

			group = tk.M{
				"_id": "$timestamp",
				// "energy":    tk.M{"$sum": "$energy"},
				"windspeed":     tk.M{"$avg": "$windspeed"},
				"power":         tk.M{"$sum": "$activepower"},
				"winddirection": tk.M{"$avg": "$winddirection"},
				"nacellepos":    tk.M{"$avg": "$nacelleposition"},
				"rotorrpm":      tk.M{"$avg": "$rotorrpm"},
				"pitchangle":    tk.M{"$avg": "$pitchangle"},
			}
		} else {
			collName = new(ScadaDataOEM).TableName()
			match.Set("dateinfo.dateid", tk.M{"$gte": tStart, "$lte": tEnd})
			// match.Set("denwindspeed", tk.M{"$lte": 25})
			match.Set("turbine", turbine)

			group = tk.M{
				"_id":        "$timestamp",
				"power":      tk.M{"$sum": "$denpower"},
				"windspeed":  tk.M{"$avg": "$denwindspeed"},
				"production": tk.M{"$avg": "$energy"},
			}
		}

		pipes = append(pipes, tk.M{"$match": match})
		pipes = append(pipes, tk.M{"$group": group})

		if tStart.Year() != 1 && tEnd.Year() != 1 {
			pipes = append(pipes, tk.M{"$sort": tk.M{"_id": 1}})
		} else {
			pipes = append(pipes, tk.M{"$sort": tk.M{"_id": -1}})
			pipes = append(pipes, tk.M{"$limit": 10})
		}

		// log.Printf("%v \n", collName)
		// log.Printf("%v | %v \n", tStart.String(), tEnd.String())

		// for _, p := range pipes {
		// 	log.Printf(">>>>>> %#v \n", p)
		// }

		csr, e := DB().Connection.NewQuery().
			From(collName).
			Command("pipe", pipes).
			Cursor(nil)

		if e != nil {
			return helper.CreateResult(false, nil, e.Error())
		}

		list := []tk.M{}
		e = csr.Fetch(&list, 0, false)

		defer csr.Close()

		// log.Printf(">>> %v \n", len(list))

		for _, tag := range tags {
			var dts [][]interface{}
			var dterr [][]interface{}
			columnTag := mapField[tag]
			if len(list) > 0 {
				for _, val := range list {
					timestamp := tk.ToInt(tk.ToString(val.Get("_id").(time.Time).Unix())+"000", tk.RoundingAuto)
					var tagVal float64

					// if tag == "production" {
					// 	tagVal = val.GetFloat64(tag)
					// } else if tag == "windspeed" {
					// 	tagVal = val.GetFloat64(tag)
					// if tagVal < 0 {
					// 	tagVal = 0.0
					// }
					// } else {
					tagVal = val.GetFloat64(tag)
					// }

					dt := []interface{}{timestamp, tagVal}
					dts = append(dts, dt)

					if tagVal < columnTag.MinValue || tagVal > columnTag.MaxValue {
						dterr = append(dterr, []interface{}{timestamp, columnTag.MaxValue})
					}
				}

				resultChart = append(resultChart, tk.M{"name": columnTag.Name, "data": dts, "dataerr": dterr, "unit": columnTag.Unit, "minval": columnTag.MinValue, "maxval": columnTag.MaxValue})
			}
		}

		csr.Close()
	}

	data := struct {
		Data ResDataAvail
	}{
		Data: ResDataAvail{
			Chart:      resultChart,
			PeriodList: periodList,
			Breaks:     breaks,
		},
	}

	return helper.CreateResult(true, data, "success")
}

func GetHFDData(turbine string, tStart time.Time, tEnd time.Time, tags []string) (result []tk.M, empty []tk.M, e error) {
	log.Printf(">>> %v - %v | %v - %v \n", tStart.String(), tStart.UTC().String(), tEnd.String(), tEnd.UTC().String())
	prefix := "data_"

	for {
		startStr := tStart.Format("20060102150405")
		// endStr := tEnd.Format("20060102150405")

		minute := tk.ToFloat64(tk.ToInt(tStart.UTC().Format("4"), tk.RoundingAuto)*60, 0, tk.RoundingAuto)
		second := tk.ToFloat64(tStart.UTC().Format("5"), 0, tk.RoundingAuto)

		totalSeconds := minute + second
		minuteDiv := math.Mod(totalSeconds, float64(600))

		newTime := tStart.UTC().Add(time.Duration(600-minuteDiv) * time.Second).UTC()

		f1 := newTime.Format("20060102")
		f2 := newTime.Format("15")
		f3 := newTime.Format("1504")

		separator := string(os.PathSeparator)

		folder := f1 + separator + f2 + separator + f3
		file := prefix + startStr + ".csv"

		path := helper.GetHFDFolder() + folder + separator + file
		tmpResult, err := ReadHFDFile(path, tags)

		if err != nil {
			// log.Printf("Err: %v \n", err.Error())
		}

		// if len(tmpResult) > 0 {
		mapTag := map[string][]float64{}
		res := tk.M{}

		for _, r := range tmpResult {
			for _, tag := range tags {
				if tag == r.Tag {
					mapTag[tag] = append(mapTag[tag], r.Value)
				}
			}
		}

		res.Set("timestamp", tStart)
		for n, mp := range mapTag {
			var value float64
			if len(mp) > 0 {
				for _, v := range mp {
					value += v
				}

				value = value / float64(len(mp))
				res.Set(n, value)
			} else {
				res.Set(n, nil)
			}
		}

		result = append(result, res)
		// }

		// else {
		// 	res := tk.M{}
		// 	res.Set("from", tStart)
		// 	res.Set("to", tStart.Add(5*time.Second))

		// 	empty = append(empty, res)
		// }

		if tStart.UTC().Sub(tEnd.UTC()).Seconds() >= 0 {
			break
		}

		// log.Printf(">>> %v \n", startStr)

		modSecond := math.Mod(second, float64(5))
		if modSecond == 0 {
			tStart = tStart.Add(5 * time.Second)
		} else {
			tStart = tStart.UTC().Add(time.Duration(5-modSecond) * time.Second).UTC()
		}

		// tStart = tStart.Add(1 * time.Second)
	}

	return
}

func ReadHFDFile(path string, tags []string) (result []HFDModel, e error) {
	fr, e := os.Open(path)
	defer fr.Close()
	if e != nil {
		fr.Close()
		return
	}

	read := csv.NewReader(bufio.NewReader(fr))
	for {
		record, err := read.Read()
		if err == io.EOF {
			fr.Close()
			break
		}

		timestamp, _ := time.Parse("2006-01-02 15:04:05", record[0])
		turbine := record[1]
		tag := record[2]
		value, _ := tk.StringToFloat(record[3])

		result = append(result, HFDModel{
			Timestamp: timestamp,
			Turbine:   turbine,
			Tag:       tag,
			Value:     value,
		})
	}

	return
}
