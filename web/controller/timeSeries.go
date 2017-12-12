package controller

import (
	"bufio"
	. "eaciit/wfdemo-git/library/core"
	. "eaciit/wfdemo-git/library/models"
	"eaciit/wfdemo-git/web/helper"
	"encoding/csv"
	"fmt"
	"io"
	// "math"
	"os"

	"time"

	"sort"

	"path/filepath"
	"strings"

	"github.com/eaciit/dbox"
	"github.com/eaciit/knot/knot.v1"
	tk "github.com/eaciit/toolkit"
)

var (
	// notAvailValue    = -9999999.0
	// notAvailValueOEM = -99999.0
	mapField = map[string]MappingColumn{
		"windspeed":                   MappingColumn{"Wind Speed", "WindSpeed_ms", "m/s", 0.0, 50.0, "$avg"},
		"power":                       MappingColumn{"Power", "ActivePower_kW", "kW", -200, 2100.0 + (2100.0 * 0.10), "$sum"},
		"production":                  MappingColumn{"Production", "Production", "kWh", -200, 2100.0, "$sum"},
		"winddirection":               MappingColumn{"Wind Direction", "WindDirection", "Degree", 0.0, 360.0, "$avg"},
		"nacellepos":                  MappingColumn{"Nacelle Direction", "NacellePos", "Degree", 0.0, 360.0, "$avg"},
		"rotorrpm":                    MappingColumn{"Rotor RPM", "RotorSpeed_RPM", "RPM", 0.0, 30.0, "$avg"},
		"genrpm":                      MappingColumn{"Generator RPM", "GenSpeed_RPM", "RPM", 0.0, 30.0, "$avg"},
		"pitchangle":                  MappingColumn{"Pitch Angle", "PitchAngle", "Degree", -10.0, 120.0, "$avg"},
		"PitchCabinetTempBlade1":      MappingColumn{"Pitch Cabinet Temp Blade 1", "PitchCabinetTempBlade1", "Degree", -10.0, 120.0, "$avg"},
		"PitchCabinetTempBlade2":      MappingColumn{"Pitch Cabinet Temp Blade 2", "PitchCabinetTempBlade2", "Degree", -10.0, 120.0, "$avg"},
		"PitchCabinetTempBlade3":      MappingColumn{"Pitch Cabinet Temp Blade 3", "PitchCabinetTempBlade3", "Degree", -10.0, 120.0, "$avg"},
		"PitchConvInternalTempBlade1": MappingColumn{"Pitch Conv Internal Temp Blade 1", "PitchConvInternalTempBlade1", "Degree", -10.0, 120.0, "$avg"},
		"PitchConvInternalTempBlade2": MappingColumn{"Pitch Conv Internal Temp Blade 2", "PitchConvInternalTempBlade2", "Degree", -10.0, 120.0, "$avg"},
		"PitchConvInternalTempBlade3": MappingColumn{"Pitch Conv Internal Temp Blade 3", "PitchConvInternalTempBlade3", "Degree", -10.0, 120.0, "$avg"},
		"TempG1L1":                    MappingColumn{"TempG1 L1", "TempG1L1", "Degree", -10.0, 200.0, "$avg"},
		"TempG1L2":                    MappingColumn{"TempG1 L2", "TempG1L2", "Degree", -10.0, 200.0, "$avg"},
		"TempG1L3":                    MappingColumn{"TempG1 L3", "TempG1L3", "Degree", -10.0, 200.0, "$avg"},
		"TempGeneratorBearingDE":      MappingColumn{"Temp Generator Bearing DE", "TempGeneratorBearingDE", "Degree", -10.0, 200.0, "$avg"},
		"TempGeneratorBearingNDE":     MappingColumn{"Temp Generator Bearing NDE", "TempGeneratorBearingNDE", "Degree", -10.0, 200.0, "$avg"},
		"TempGearBoxOilSump":          MappingColumn{"Temp Gear Box Oil Sump", "TempGearBoxOilSump", "Degree", -10.0, 200.0, "$avg"},
		"TempHubBearing":              MappingColumn{"Temp Hub Bearing", "TempHubBearing", "Degree", -10.0, 200.0, "$avg"},
		"TempGeneratorChoke":          MappingColumn{"Temp Generator Choke", "TempGeneratorChoke", "Degree", -10.0, 200.0, "$avg"},
		"TempGridChoke":               MappingColumn{"Temp Grid Choke", "TempGridChoke", "Degree", -10.0, 200.0, "$avg"},
		"TempGeneratorCoolingUnit":    MappingColumn{"Temp Generator Cooling Unit", "TempGeneratorCoolingUnit", "Degree", -10.0, 200.0, "$avg"},
		"TempConvCabinet2":            MappingColumn{"Temp Conv Cabinet 2", "TempConvCabinet2", "Degree", -10.0, 200.0, "$avg"},
		"TempOutdoor":                 MappingColumn{"Temp Outdoor", "TempOutdoor", "Degree", -10.0, 200.0, "$avg"},
		"TempSlipRing":                MappingColumn{"Temp Slip Ring", "TempSlipRing", "Degree", -10.0, 200.0, "$avg"},
		"TransformerWindingTemp1":     MappingColumn{"Transformer Winding Temp 1", "TransformerWindingTemp1", "Degree", -10.0, 200.0, "$avg"},
		"TransformerWindingTemp2":     MappingColumn{"Transformer Winding Temp 2", "TransformerWindingTemp2", "Degree", -10.0, 200.0, "$avg"},
		"TransformerWindingTemp3":     MappingColumn{"Transformer Winding Temp 3", "TransformerWindingTemp3", "Degree", -10.0, 200.0, "$avg"},
		"TempShaftBearing1":           MappingColumn{"Temp Shaft Bearing 1", "TempShaftBearing1", "Degree", -10.0, 200.0, "$avg"},
		"TempShaftBearing2":           MappingColumn{"Temp Shaft Bearing 2", "TempShaftBearing2", "Degree", -10.0, 200.0, "$avg"},
		"TempShaftBearing3":           MappingColumn{"Temp Shaft Bearing 3", "TempShaftBearing3", "Degree", -10.0, 200.0, "$avg"},
		"TempGearBoxIMSDE":            MappingColumn{"Temp Gear Box IMS DE", "TempGearBoxIMSDE", "Degree", -10.0, 200.0, "$avg"},
		"TempBottomControlSection":    MappingColumn{"Temp Bottom Control Section", "TempBottomControlSection", "Degree", -10.0, 200.0, "$avg"},
		"TempBottomPowerSection":      MappingColumn{"Temp Bottom Power Section", "TempBottomPowerSection", "Degree", -10.0, 200.0, "$avg"},
		"TempCabinetTopBox":           MappingColumn{"Temp Cabinet Top Box", "TempCabinetTopBox", "Degree", -10.0, 200.0, "$avg"},
		"TempNacelle":                 MappingColumn{"Temp Nacelle", "TempNacelle", "Degree", -10.0, 200.0, "$avg"},
		"VoltageL1":                   MappingColumn{"Voltage L 1", "GridPPVPhaseAB", "Volt", -10.0, 1000.0, "$avg"},
		"VoltageL2":                   MappingColumn{"Voltage L 2", "GridPPVPhaseBC", "Volt", -10.0, 1000.0, "$avg"},
		"VoltageL3":                   MappingColumn{"Voltage L 3", "GridPPVPhaseCA", "Volt", -10.0, 1000.0, "$avg"},
		"PitchAccuV1":                 MappingColumn{"Blade Voltage V1", "PitchAccuV1", "Volt", -10.0, 1000.0, "$avg"},
		"PitchAccuV2":                 MappingColumn{"Blade Voltage V2", "PitchAccuV2", "Volt", -10.0, 1000.0, "$avg"},
		"PitchAccuV3":                 MappingColumn{"Blade Voltage V3", "PitchAccuV3", "Volt", -10.0, 1000.0, "$avg"},
	}
)

type TimeSeriesController struct {
	App
}

type RHFDModel struct {
	HFDModel
	SValue float64
	CValue float64
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
	Outliers   [][]interface{}
}

type MappingColumn struct {
	Name      string
	SecField  string
	Unit      string
	MinValue  float64
	MaxValue  float64
	Aggregate string
}

func CreateTimeSeriesController() *TimeSeriesController {
	var controller = new(TimeSeriesController)
	return controller
}

func (m *TimeSeriesController) GetAvailDate(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	return helper.CreateResult(true, k.Session("availdate", ""), "success")
}

func (m *TimeSeriesController) GetDataHFD(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	breaks := []tk.M{}
	resultChart := []tk.M{}
	periodList := []tk.M{}
	outliers := map[int]bool{}

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

	projectName := p.Project

	turbine := p.Turbine
	dataType := p.DataType
	pageType := p.PageType

	// log.Printf(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> %v | %v \n", dataType, pageType)

	// default tags
	tags := p.TagList

	// if len(p.TagList) == 0 {
	// 	tags = []string{"windspeed", "power"}
	// }

	if pageType == "HFD" && dataType == "SEC" {
		secTags := []string{}
		for _, tg := range tags {
			secTags = append(secTags, mapField[tg].SecField)
		}

		// tags = secTags

		for {
			if tStart.Sub(tEnd).Seconds() >= 0 {
				break
			}

			before := tStart
			// tStart = tStart.UTC().Add(time.Duration(24) * time.Hour)
			// log.Printf(">>>>>> %v | %v | %v \n", tEnd.UTC().Sub(tStart.UTC()).Seconds(), tStart.UTC(), tEnd.UTC())
			tStart = tStart.Add(time.Duration(tEnd.Sub(tStart).Seconds()) * time.Second)

			beforeInt := tk.ToInt(tk.ToString(before.Unix())+"000", tk.RoundingAuto)
			afterInt := tk.ToInt(tk.ToString(tStart.Unix())+"000", tk.RoundingAuto)

			periodList = append(periodList, tk.M{"starttime": before, "endtime": tStart, "starttimeint": beforeInt, "endtimeint": afterInt})
		}

		if len(periodList) > 0 || p.IsHour {
			for idx, pl := range periodList {
				current := pl
				currStar := current.Get("starttime").(time.Time)
				currEnd := current.Get("endtime").(time.Time)

				hfds, empty, e := GetHFDDataRev(projectName, turbine, currStar, currEnd, tags, secTags)

				breaks = append(breaks, empty...)

				if e != nil {
					return helper.CreateResult(false, nil, e.Error())
				}

				if len(hfds) > 0 || p.IsHour {
					for _, tag := range tags {
						var dts [][]interface{}
						var dterr [][]interface{}
						columnTag := mapField[tag]

						// log.Printf(">> %v \n", len(hfds))

						/*if len(hfds) > 0 {
							log.Printf(">> tag: %v \n", tag)
							dts, dterr, outliers = constructData(hfds, tStart, tEnd, 5, pageType, tag, columnTag)
						}*/

						for _, val := range hfds {
							timestamp := tk.ToInt(tk.ToString(val.Get("timestamp").(time.Time).Unix())+"000", tk.RoundingAuto)
							var dt []interface{}

							if val.Get(columnTag.SecField) != nil {
								tagVal := val.GetFloat64(columnTag.SecField)
								dt = []interface{}{timestamp, tagVal}

								if tagVal == float64(-99999.00) || tagVal == float64(-999999.00) || tagVal == float64(-9999999.00) {
									dt = []interface{}{timestamp, nil}
								} else if tagVal < columnTag.MinValue || tagVal > columnTag.MaxValue {
									dterr = append(dterr, []interface{}{timestamp, 100.0})
									outliers[timestamp] = true
								}
							} else {
								dt = []interface{}{timestamp, nil}
							}

							dts = append(dts, dt)
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

		if pageType == "HFD" || pageType == "MIX" {
			collName = "Scada10MinHFD"
			match.Set("dateinfo.dateid", tk.M{"$gte": tStart, "$lte": tEnd})
			match.Set("isnull", false)
			match.Set("turbine", turbine)
			group.Set("_id", "$timestamp")
			for _, tag := range tags {
				group.Set(tag, tk.M{mapField[tag].Aggregate: "$" + strings.ToLower(mapField[tag].SecField)})
			}
		} else if pageType == "OEM" {
			collName = new(ScadaDataOEM).TableName()
			match.Set("dateinfo.dateid", tk.M{"$gte": tStart, "$lte": tEnd})
			match.Set("turbine", turbine)

			group = tk.M{
				"_id":        "$timestamp",
				"power":      tk.M{"$sum": "$denpower"},
				"windspeed":  tk.M{"$avg": "$denwindspeed"},
				"production": tk.M{"$avg": "$energy"},
			}
		}

		list := []tk.M{}

		if pageType != "LIVE" {
			pipes = append(pipes, tk.M{"$match": match})
			pipes = append(pipes, tk.M{"$group": group})

			if tStart.Year() != 1 && tEnd.Year() != 1 {
				pipes = append(pipes, tk.M{"$sort": tk.M{"_id": 1}})
			} else {
				pipes = append(pipes, tk.M{"$sort": tk.M{"_id": -1}})
				pipes = append(pipes, tk.M{"$limit": 5})
			}

			csr, e := DB().Connection.NewQuery().
				From(collName).
				Command("pipe", pipes).
				Cursor(nil)
			defer csr.Close()

			if e != nil {
				return helper.CreateResult(false, nil, e.Error())
			}

			e = csr.Fetch(&list, 0, false)

		} else {
			list = getDataLiveNew(projectName, turbine, tStart, p.TagList)
		}

		if pageType == "MIX" {
			pipes = []tk.M{}
			match.Unset("isnull")
			group = tk.M{
				"_id":        "$timestamp",
				"power":      tk.M{"$sum": "$denpower"},
				"windspeed":  tk.M{"$avg": "$denwindspeed"},
				"production": tk.M{"$avg": "$energy"},
			}

			pipes = append(pipes, tk.M{"$match": match})
			pipes = append(pipes, tk.M{"$group": group})

			if tStart.Year() != 1 && tEnd.Year() != 1 {
				pipes = append(pipes, tk.M{"$sort": tk.M{"_id": 1}})
			} else {
				pipes = append(pipes, tk.M{"$sort": tk.M{"_id": -1}})
				pipes = append(pipes, tk.M{"$limit": 5})
			}

			list = MixDataIfAny(list, pipes, new(ScadaDataOEM).TableName())
		}

		for _, tag := range tags {
			var dts [][]interface{}
			var dterr [][]interface{}
			columnTag := mapField[tag]
			// log.Printf("> %v | %v | %v \n", tag, tmpStart)

			if len(list) > 0 {
				dts, dterr, outliers = constructData(list, tStart, tEnd, 600, pageType, tag, columnTag)
			}
			resultChart = append(resultChart, tk.M{"name": columnTag.Name, "data": dts, "dataerr": dterr, "unit": columnTag.Unit, "minval": columnTag.MinValue, "maxval": columnTag.MaxValue})
		}

	}

	revOutliers := [][]interface{}{}
	tmpOutliers := []int{}

	for it := range outliers {
		found := false
		for _, t := range tmpOutliers {
			if t == it {
				found = true
				break
			}
		}

		if !found {
			tmpOutliers = append(tmpOutliers, it)
		}
	}

	sort.Ints(tmpOutliers)

	for _, timestamp := range tmpOutliers {
		revOutliers = append(revOutliers, []interface{}{timestamp, 100.0})
	}

	data := struct {
		Data ResDataAvail
	}{
		Data: ResDataAvail{
			Chart:      resultChart,
			PeriodList: periodList,
			Breaks:     breaks,
			Outliers:   revOutliers,
		},
	}

	return helper.CreateResult(true, data, "success")
}

func constructData(list []tk.M, tStart time.Time, tEnd time.Time, seconds float64, pageType string, tag string, columnTag MappingColumn) (dts [][]interface{}, dterr [][]interface{}, outliers map[int]bool) {
	outliers = map[int]bool{}

	tmpStart := tStart.UTC()

	for _, val := range list {
		timestamp := val.Get("_id").(time.Time).UTC()
		timestampInt := tk.ToInt(tk.ToString(timestamp.Unix())+"000", tk.RoundingAuto)
		var tagVal float64

		if tmpStart.Before(timestamp) && timestamp.Sub(tmpStart).Seconds() > seconds && pageType != "LIVE" {
			// get the time is not exist in collection
			notExist := constructDataNotExist(tmpStart, timestamp, seconds)
			if len(notExist) > 0 {
				dts = append(dts, notExist...)
			}
			// log.Printf("> %v > %v | %v - %v > %v | %v > %v > %v \n", tag, tmpStart, timestamp, tmpStart.Before(timestamp), len(notExist), len(dts), dts[0], notExist[0])
		}

		// if tag == "production" {
		// 	tagVal = val.GetFloat64(tag)
		// } else if tag == "windspeed" {
		// 	tagVal = val.GetFloat64(tag)
		// if tagVal < 0 {
		// 	tagVal = 0.0
		// }
		// } else {

		var dt []interface{}
		isNill := true

		if seconds >= 600 {
			if val.Get(tag) != nil {
				isNill = false
				tagVal = val.GetFloat64(tag)
			}
		} else {
			if val.Get(columnTag.SecField) != nil {
				isNill = false
				tagVal = val.GetFloat64(columnTag.SecField)
			}
		}

		// }

		// log.Printf(">> %v | %v | %v | %v \n", isNill, tagVal, timestamp, timestampInt)

		if !isNill {
			dt = []interface{}{timestampInt, tagVal}
			if tagVal == float64(-99999.00) || tagVal == float64(-999999.00) || tagVal == float64(-9999999.00) {
				dt = []interface{}{timestampInt, nil}
			} else if tagVal < columnTag.MinValue || tagVal > columnTag.MaxValue {
				dterr = append(dterr, []interface{}{timestampInt, 100.0})
				outliers[timestampInt] = true
			}
		} else {
			dt = []interface{}{timestampInt, nil}
		}

		dts = append(dts, dt)
		tmpStart = timestamp
	}

	// log.Printf("> %v > %v | %v - %v \n", tag, tmpStart, tmpEnd, tmpStart.Before(tEnd))

	if tmpStart.Before(tEnd) && tEnd.Sub(tmpStart).Seconds() > seconds && pageType != "LIVE" {
		// get the time is not exist in collection
		notExist := constructDataNotExist(tmpStart, tEnd, seconds)
		if len(notExist) > 0 {
			dts = append(dts, notExist...)
		}
		// log.Printf("> %v > %v | %v - %v > %v | %v \n", tag, tmpStart, tEnd, tmpStart.Before(tEnd), len(notExist), len(dts))
	}

	return
}

func constructDataNotExist(start time.Time, end time.Time, seconds float64) (result [][]interface{}) {
	for {
		if end.UTC().Sub(start.UTC()).Seconds() <= seconds {
			// log.Println(">>>> BREAK")
			break
		}

		// before := end.UTC()
		timestamp := tk.ToInt(tk.ToString(start.UTC().Unix())+"000", tk.RoundingAuto)
		result = append(result, []interface{}{timestamp, nil})

		start = start.UTC().Add(time.Duration(seconds) * time.Second)
	}

	return
}

func getDataLiveNew(project string, turbine string, tStart time.Time, tags []string) (result []tk.M) {
	filter := []*dbox.Filter{}
	tmpRes := map[string]interface{}{}

	filter = append(filter, dbox.Eq("projectname", project))
	filter = append(filter, dbox.Eq("turbine", turbine))
	// filter = append(filter, dbox.In("tags", []interface{}{"ActivePower_kW", "WindSpeed_ms"}))

	if tStart.Year() != 1 {
		filter = append(filter, dbox.Gt("timestamp", tStart.UTC()))
	}
	rconn := DBRealtime()

	csr, err := rconn.NewQuery().From(new(ScadaRealTimeNew).TableName()).
		Where(dbox.And(filter...)).
		Order("-timestamp").
		Cursor(nil)

	defer csr.Close()

	if err != nil {
		tk.Println(err.Error())
	}

	listtag := tk.M{}.Set("power", "ActivePower_kW").Set("windspeed", "WindSpeed_ms").Set("rotorrpm", "RotorSpeed_RPM").Set("pitchangle", "PitchAngle")
	tstamp := time.Time{}
	if csr.Count() > 0 {
		for {
			data := tk.M{}
			err = csr.Fetch(&data, 1, false)
			if err != nil {
				break
			}
			itime := data.Get("timestamp", time.Time{}).(time.Time)
			if tstamp.IsZero() || tstamp.Before(itime) {
				tstamp = itime
			}

			for _, xTag := range tags {
				if listtag.GetString(xTag) == data.GetString("tags") {
					tmpRes[xTag] = data.GetFloat64("value")
				}
			}
		}
		csr.Close()
		result = append(result, tk.M{"_id": tstamp})

		for tag, mp := range tmpRes {
			result[0].Set(tag, mp)
		}

	}
	return
}

func GetHFDDataRev(project string, turbine string, tStart time.Time, tEnd time.Time, tags []string, secTags []string) (result []tk.M, empty []tk.M, e error) {

	dhfd, draw, allkey := make(map[time.Time]tk.M, 0), make(map[time.Time]tk.M, 0), make(map[time.Time]int, 0)

	match := tk.M{}
	pipes := []tk.M{}
	projection := map[string]int{"timestamp": 1}
	for _, tag := range secTags {
		projection[strings.ToLower(tag)] = 1
	}
	match.Set("timestamp", tk.M{"$gte": tStart.UTC(), "$lt": tEnd.UTC()})
	match.Set("projectname", project)
	match.Set("turbine", turbine)
	match.Set("isnull", false)

	pipes = append(pipes, tk.M{"$match": match})
	pipes = append(pipes, tk.M{"$project": projection})
	pipes = append(pipes, tk.M{"$sort": tk.M{"timestamp": 1}})
	csr, e := DB().Connection.NewQuery().
		// Select("timestamp", "fast_windspeed_ms", "fast_activepower_kw", "slow_winddirection",
		// 	"slow_nacellepos", "fast_rotorspeed_rpm", "fast_genspeed_rpm", "fast_pitchangle",
		// 	"pitchcabinettempblade1", "pitchcabinettempblade2", "pitchcabinettempblade3",
		// 	"pitchconvinternaltempblade1", "pitchconvinternaltempblade2", "pitchconvinternaltempblade3").
		From("Scada10MinHFD").
		Command("pipe", pipes).
		Cursor(nil)
	defer csr.Close()

	if e == nil {
		// fname := map[string]string{"windspeed": "fast_windspeed_ms", "power": "fast_activepower_kw", "winddirection": "slow_winddirection",
		// 	"nacellepos": "slow_nacellepos", "rotorrpm": "fast_rotorspeed_rpm", "genrpm": "fast_genspeed_rpm", "pitchangle": "fast_pitchangle",
		// 	"PitchCabinetTempBlade1": "pitchcabinettempblade1", "PitchCabinetTempBlade2": "pitchcabinettempblade2", "PitchCabinetTempBlade3": "pitchcabinettempblade3",
		// 	"PitchConvInternalTempBlade1": "pitchconvinternaltempblade1", "PitchConvInternalTempBlade2": "pitchconvinternaltempblade2", "PitchConvInternalTempBlade3": "pitchconvinternaltempblade3"}
		for {
			dt, dts := tk.M{}, tk.M{}
			err := csr.Fetch(&dt, 1, false)
			if err != nil {
				break
			}

			dtime := dt.Get("timestamp", time.Time{}).(time.Time).UTC()
			dts.Set("timestamp", dtime)

			for _, tag := range tags {
				mc := mapField[tag]
				val := dt.GetFloat64(strings.ToLower(mc.SecField))

				if val == float64(-99999.00) || val == float64(-999999.00) || val == float64(-9999999.00) {
					dts.Set(mc.SecField, nil)
				} else {
					dts.Set(mc.SecField, val)
				}
			}

			dhfd[dtime] = dts
			allkey[dtime] = 1
		}
	}

	tStart, tEnd = tStart.UTC(), tEnd.UTC()
	tstart10m, tend10m := getNext10Min(tStart), getNext10Min(tEnd)

	for {
		f1 := tstart10m.Format("20060102")
		f2 := tstart10m.Format("15")
		f3 := tstart10m.Format("1504")

		_fpath := filepath.Join(helper.GetHFDFolder(), strings.ToLower(project), f1, f2, f3, tk.Sprintf("%s.csv", turbine))
		// tk.Println(">>> ", _fpath)
		tmpresult, _ := ReadHFDFile(_fpath, secTags)

		for _, res := range tmpresult {
			dts := tk.M{}
			if _val, _cond := draw[res.Timestamp]; _cond {
				dts = _val
			}
			dts.Set("timestamp", res.Timestamp)
			dts.Set(res.Tag, res.Value)

			draw[res.Timestamp] = dts
			allkey[res.Timestamp] = 1
		}

		if tstart10m.After(tend10m) {
			break
		}

		tstart10m = tstart10m.Add(10 * time.Minute)
	}

	arrkey := []time.Time{}
	for ktime, _ := range allkey {
		arrkey = append(arrkey, ktime)
	}

	sort.Sort(ByTime(arrkey))

	for _, key := range arrkey {
		if rval, rcond := draw[key]; rcond {
			result = append(result, rval)
		} else if hval, hcond := dhfd[key]; hcond {
			result = append(result, hval)
		}
	}

	return
}

func ReadHFDFile(path string, tags []string) (result []HFDModel, e error) {
	fr, e := os.Open(path)
	defer fr.Close()
	if e != nil {
		return
	}

	rawres := make(map[string]RHFDModel, 0)
	read := csv.NewReader(bufio.NewReader(fr))
	for {
		record, err := read.Read()
		if err == io.EOF {
			fr.Close()
			break
		}

		if len(record) > 4 && strings.ToLower(record[4]) != "good" && strings.ToLower(record[4]) != "" {
			continue
		}

		timestamp, _ := time.Parse("2006-01-02 15:04:05", record[0])
		second := tk.ToInt(timestamp.Format("5"), tk.RoundingAuto)
		if val := second % 5; val != 0 {
			timestamp = timestamp.Add(time.Second * time.Duration(5-val))
		}

		turbine := record[1]
		tag := record[2]
		for _, tg := range tags {
			if tg == tag {
				skey := fmt.Sprintf("%s_%s_%s", tag, turbine, timestamp.Format("20060102150405"))
				value, _ := tk.StringToFloat(record[3])

				_rhfd := rawres[skey]
				_rhfd.Timestamp = timestamp
				_rhfd.Turbine = turbine
				_rhfd.Tag = tag
				_rhfd.SValue += value
				_rhfd.CValue += 1

				rawres[skey] = _rhfd

			}
		}
	}

	result = []HFDModel{}
	for _, val := range rawres {
		result = append(result, HFDModel{
			Timestamp: val.Timestamp,
			Turbine:   val.Turbine,
			Tag:       val.Tag,
			Value:     tk.Div(val.SValue, val.CValue),
		})
	}

	return
}

type ByTime []time.Time

func (b ByTime) Len() int {
	return len(b)
}

func (b ByTime) Swap(i, j int) {
	b[i], b[j] = b[j], b[i]
}

func (b ByTime) Less(i, j int) bool {
	return b[i].Before(b[j])
}

func MixDataIfAny(ilist []tk.M, pipes []tk.M, tname string) (list []tk.M) {
	list = ilist

	csr, e := DB().Connection.NewQuery().
		From(tname).
		Command("pipe", pipes).
		Cursor(nil)
	defer csr.Close()

	// tk.Println(pipes, ">>>", csr.Count(), tname, e)

	if e != nil {
		return
	}

	itkm, skey := map[string]tk.M{}, []string{}
	for _, val := range ilist {
		itime := val.Get("_id", time.Time{}).(time.Time).UTC()
		if itime.IsZero() {
			continue
		}

		key := itime.Format("20060102150405")
		itkm[key] = val
	}

	icount := 0
	for {
		tkm := tk.M{}
		e = csr.Fetch(&tkm, 1, false)
		if e != nil {
			break
		}

		itime := tkm.Get("_id", time.Time{}).(time.Time).UTC()
		if itime.IsZero() {
			continue
		}

		key := itime.Format("20060102150405")
		itkm[key] = tkm
		icount++
	}

	if icount == 0 {
		list = ilist
		return
	}

	for key, _ := range itkm {
		skey = append(skey, key)
	}

	sort.Strings(skey)

	list = []tk.M{}
	for _, val := range skey {
		list = append(list, itkm[val])
	}

	return
}

// func GetHFDData(project string, turbine string, tStart time.Time, tEnd time.Time, tags []string, secTags []string) (result []tk.M, empty []tk.M, e error) {
// 	// log.Printf(">>> %v - %v | %v - %v \n", tStart.String(), tStart.UTC().String(), tEnd.String(), tEnd.UTC().String())
// 	prefix := "data_"
// 	emptyLen := 0
// 	emptyStartStr := ""

// 	var emptyStart time.Time
// 	var emptySeconds float64

// 	for {
// 		startStr := tStart.UTC().Format("20060102150405")
// 		// endStr := tEnd.Format("20060102150405")

// 		// log.Printf("%v | %v \n", tStart.String(), tStart.UTC().String())

// 		if emptyStartStr != "" {
// 			// fill in empty seconds HFD data with minutes HFD data
// 			emptyStart, _ = time.Parse("20060102150405", emptyStartStr)
// 			emptyStart = emptyStart
// 			emptySeconds = tStart.Sub(emptyStart).Seconds()

// 			// log.Printf(">>>> %v | %v ===> %v \n", emptyStart.String(), tStart.UTC().String(), emptySeconds)

// 			if emptySeconds >= float64(600) || tStart.UTC().Sub(tEnd.UTC()).Seconds() >= 0 {
// 				match := tk.M{}
// 				group := tk.M{}
// 				pipes := []tk.M{}

// 				match.Set("dateinfo.dateid", tk.M{"$gte": emptyStart, "$lt": tStart.UTC()})
// 				// match.Set("fast_windspeed_ms_stddev", tk.M{"$lte": 25})
// 				match.Set("projectname", project)
// 				match.Set("turbine", turbine)
// 				// match.Set("available", 1)

// 				group = tk.M{
// 					"_id": "$timestamp",
// 					// "energy":    tk.M{"$sum": "$energy"},
// 					"windspeed":     tk.M{"$avg": "$fast_windspeed_ms"},
// 					"power":         tk.M{"$sum": "$fast_activepower_kw"},
// 					"winddirection": tk.M{"$avg": "$slow_winddirection"},
// 					"nacellepos":    tk.M{"$avg": "$slow_nacellepos"},
// 					"rotorrpm":      tk.M{"$avg": "$fast_rotorspeed_rpm"},
// 					"genrpm":        tk.M{"$avg": "$fast_genspeed_rpm"},
// 					"pitchangle":    tk.M{"$avg": "$fast_pitchangle"},
// 				}

// 				pipes = append(pipes, tk.M{"$match": match})
// 				pipes = append(pipes, tk.M{"$group": group})
// 				pipes = append(pipes, tk.M{"$sort": tk.M{"_id": 1}})

// 				csr, e := DB().Connection.NewQuery().
// 					From(new(ScadaDataHFD).TableName()).
// 					Command("pipe", pipes).
// 					Cursor(nil)
// 				defer csr.Close()

// 				if e == nil {
// 					list := []tk.M{}
// 					e = csr.Fetch(&list, 0, false)

// 					// log.Printf(">>>>> %v | %v => %v \n", emptyStart.String(), tStart.UTC().String(), len(list))

// 					if len(list) > 0 {
// 						for _, val := range list {
// 							dts := tk.M{}
// 							timestamp := val.Get("_id").(time.Time)

// 							dts.Set("timestamp", timestamp)

// 							for _, tag := range tags {
// 								tagVal := val.GetFloat64(tag)
// 								mc := mapField[tag]

// 								if tagVal == float64(-99999.00) || tagVal == float64(-999999.00) || tagVal == float64(-9999999.00) {
// 									dts.Set(mc.SecField, nil)
// 								} else {
// 									dts.Set(mc.SecField, tagVal)
// 								}
// 							}

// 							result = append(result, dts)
// 						}
// 					}

// 					emptyStartStr = ""
// 					emptyLen = 0
// 				}
// 			}
// 		}

// 		minute := tk.ToFloat64(tk.ToInt(tStart.Format("4"), tk.RoundingAuto)*60, 0, tk.RoundingAuto)
// 		second := tk.ToFloat64(tStart.Format("5"), 0, tk.RoundingAuto)

// 		totalSeconds := minute + second
// 		minuteDiv := math.Mod(totalSeconds, float64(600))

// 		newTime := tStart.Add(time.Duration(600-minuteDiv) * time.Second)

// 		f1 := newTime.Format("20060102")
// 		f2 := newTime.Format("15")
// 		f3 := newTime.Format("1504")

// 		separator := string(os.PathSeparator)

// 		folder := strings.ToLower(project) + separator + f1 + separator + f2 + separator + f3
// 		file1 := prefix + startStr + ".csv"
// 		file2 := turbine + "_" + startStr + ".csv"

// 		path := helper.GetHFDFolder() + folder + separator + file1
// 		if _, _err := os.Stat(path); os.IsNotExist(_err) {
// 			path = helper.GetHFDFolder() + folder + separator + file2
// 		}

// 		// log.Printf("%v | %v | %v \n", tStart.UTC().String(), newTime.String(), path)

// 		tmpResult, err := ReadHFDFile(path, secTags)

// 		if err != nil {
// 			// log.Printf("Err: %v \n", err.Error())
// 		}

// 		if len(tmpResult) > 0 {
// 			// log.Printf("len(tmpResult) > 0 || %v \n", emptyStartStr)
// 			mapTag := map[string][]float64{}
// 			res := tk.M{}

// 			for _, r := range tmpResult {
// 				for _, tag := range secTags {
// 					if tag == r.Tag && turbine == r.Turbine {
// 						mapTag[tag] = append(mapTag[tag], r.Value)
// 					}
// 				}
// 			}

// 			res.Set("timestamp", tStart)
// 			for n, mp := range mapTag {
// 				var value float64
// 				if len(mp) > 0 {
// 					for _, v := range mp {
// 						value += v
// 					}

// 					if value == float64(-99999.00) && value == float64(-999999.00) && value == float64(-9999999.00) {
// 						res.Set(n, nil)
// 					} else {
// 						value = value / float64(len(mp))
// 						res.Set(n, value)
// 					}

// 					// if startStr == "20170418225105" {
// 					// 	log.Printf(">>> %v > %#v | %v | %v \n", n, mp, value, res.Get(n))
// 					// }

// 				} else {
// 					res.Set(n, nil)
// 				}
// 			}

// 			result = append(result, res)
// 		} else {
// 			// log.Printf("else || %v \n", emptyStartStr)
// 			// res := tk.M{}
// 			// res.Set("from", tStart)
// 			// res.Set("to", tStart.Add(5*time.Second))

// 			// empty = append(empty, res)

// 			if emptyStartStr == "" {
// 				emptyLen++
// 				emptyStartStr = startStr
// 				// log.Printf("else if || %v | %v \n", emptySeconds, float64(emptyLen))
// 			} else {
// 				// log.Printf("else else || %v | %v \n", emptySeconds, float64(emptyLen))
// 				if emptySeconds/5 != float64(emptyLen) {
// 					emptyStartStr = ""
// 					emptyLen = 0
// 				} else {
// 					emptyLen++
// 				}
// 			}
// 		}

// 		if tStart.Sub(tEnd).Seconds() >= 0 {
// 			break
// 		}

// 		// log.Printf(">>> %v \n", startStr)

// 		modSecond := math.Mod(second, float64(5))
// 		if modSecond == 0 {
// 			tStart = tStart.Add(5 * time.Second)
// 		} else {
// 			tStart = tStart.Add(time.Duration(5-modSecond) * time.Second)
// 		}

// 		// tStart = tStart.Add(1 * time.Second)
// 	}

// 	return
// }

// func getDataLive(project string, turbine string, tStart time.Time, tags []string) (result []tk.M) {
// 	filter := []*dbox.Filter{}
// 	tmpRes := map[string]interface{}{}

// 	filter = append(filter, dbox.Eq("projectname", project))
// 	filter = append(filter, dbox.Eq("turbine", turbine))

// 	if tStart.Year() != 1 {
// 		filter = append(filter, dbox.Gt("timestamp", tStart.UTC()))
// 	}
// 	// rconn := lh.GetConnRealtime()
// 	// defer rconn.Close()
// 	rconn := DBRealtime()

// 	csr, err := rconn.NewQuery().From(new(ScadaRealTime).TableName()).
// 		Where(dbox.And(filter...)).
// 		Order("-timestamp").
// 		Cursor(nil)

// 	defer csr.Close()

// 	if err != nil {
// 		tk.Println(err.Error())
// 	}

// 	tstamp := time.Time{}

// 	mapFound := map[string]bool{}

// 	// log.Printf(">>> %v | %v \n", tStart.String(), csr.Count())

// 	if csr.Count() > 0 {
// 		for {
// 			data := tk.M{}
// 			err = csr.Fetch(&data, 1, false)
// 			if err != nil {
// 				break
// 			}

// 			tstamp = data.Get("timestamp", time.Time{}).(time.Time)

// 			for _, xTag := range tags {
// 				tag := xTag
// 				if xTag == "power" {
// 					tag = "activepower"
// 				}

// 				if !mapFound[tag] {
// 					tagVal := data.GetFloat64(tag)

// 					if tagVal == float64(-99999.00) || tagVal == float64(-999999.00) || tagVal == float64(-9999999.00) {
// 						tmpRes[xTag] = nil
// 					} else {
// 						tmpRes[xTag] = tagVal
// 					}

// 					mapFound[tag] = true
// 				}
// 			}

// 			count := 0
// 			for _, mp := range mapFound {
// 				if mp {
// 					count++
// 				}
// 			}

// 			if count == len(tags) {
// 				break
// 			}
// 		}
// 		csr.Close()

// 		result = append(result, tk.M{"_id": tstamp})

// 		for tag, mp := range tmpRes {
// 			result[0].Set(tag, mp)
// 		}

// 	}
// 	return
// }
