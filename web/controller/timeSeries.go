package controller

import (
	"bufio"
	. "eaciit/wfdemo-git/library/core"
	hp "eaciit/wfdemo-git/library/helper"
	. "eaciit/wfdemo-git/library/models"
	"eaciit/wfdemo-git/web/helper"
	"encoding/csv"
	"io"
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
}

func CreateTimeSeriesController() *TimeSeriesController {
	var controller = new(TimeSeriesController)
	return controller
}

func (m *TimeSeriesController) GetData(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	var (
		pipes       []tk.M
		list        []tk.M
		resultChart []tk.M
		periodList  []tk.M
	)

	p := new(PayloadTimeSeries)
	e := k.GetPayload(&p)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	tStart, tEnd, e := helper.GetStartEndDate(k, p.Period, p.DateStart, p.DateEnd)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	match := tk.M{}

	match.Set("dateinfo.dateid", tk.M{"$gte": tStart, "$lte": tEnd})
	match.Set("avgwindspeed", tk.M{"$lte": 25})
	// match.Set("turbine", p.Turbine)

	if p.Project != "" {
		anProject := strings.Split(p.Project, "(")
		match.Set("projectname", strings.TrimRight(anProject[0], " "))
	}

	group := tk.M{
		"energy":    tk.M{"$sum": "$energy"},
		"windspeed": tk.M{"$avg": "$avgwindspeed"},
	}

	group.Set("_id", "$timestamp")

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

	csr.Close()

	var dtProd, dtWS [][]interface{}

	for _, val := range list {
		// intTimestamp :=
		intTimestamp := tk.ToInt(tk.ToString(val.Get("_id").(time.Time).UTC().Unix())+"000", tk.RoundingAuto)

		energy := val.GetFloat64("energy") / 1000
		wind := val.GetFloat64("windspeed")

		dtProd = append(dtProd, []interface{}{intTimestamp, energy})
		dtWS = append(dtWS, []interface{}{intTimestamp, wind})
	}

	resultChart = append(resultChart, tk.M{"name": "Production", "data": dtProd, "unit": "MWh"})
	resultChart = append(resultChart, tk.M{"name": "Windspeed", "data": dtWS, "unit": "m/s"})

	data := struct {
		Data ResDataAvail
	}{
		Data: ResDataAvail{
			Chart:      resultChart,
			PeriodList: periodList,
		},
	}

	return helper.CreateResult(true, data, "success")
}

func (m *TimeSeriesController) GetAvailDate(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	return helper.CreateResult(true, k.Session("availdate", ""), "success")
}

func (m *TimeSeriesController) GetDataHFD(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

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

	dataType := p.DataType
	pageType := p.PageType

	// log.Printf(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> %v | %v \n", dataType, pageType)

	tags := []string{}
	mapUnit := map[string]string{}

	tags = []string{"windspeed", "power", "production"}
	mapUnit["windspeed"] = "m/s"
	mapUnit["power"] = "kW"
	mapUnit["production"] = "MWh"

	// if pageType == "HFD" {
	// set default value for HFD
	// tags = []string{"windspeed", "power"}
	// mapUnit["windspeed"] = "m/s"
	// mapUnit["power"] = "kW"

	if len(p.TagList) > 0 {
		tags = p.TagList
	}

	// } else if pageType == "OEM" {
	// 	// set default value for OEM

	// }

	if pageType == "HFD" && dataType == "SEC" {
		for {
			if tStart.UTC().Sub(tEnd.UTC()).Seconds() >= 0 {
				break
			}

			before := tStart.UTC()
			tStart = tStart.UTC().Add(time.Duration(3) * time.Hour)
			periodList = append(periodList, tk.M{"starttime": before, "endtime": tStart.UTC()})
		}

		if len(periodList) > 0 {
			for idx, pl := range periodList {
				current := pl
				currStar := current.Get("starttime").(time.Time)
				currEnd := current.Get("endtime").(time.Time)

				hfds, e := GetHFDData("HBR004", currStar, currEnd, tags)

				if e != nil {
					return helper.CreateResult(false, nil, e.Error())
				}

				if len(hfds) > 0 {
					for _, tag := range tags {
						var dts [][]interface{}
						for _, val := range hfds {
							timestamp := val.Get("timestamp").(time.Time).Unix()
							tagVal := val.GetFloat64(tag)
							dt := []interface{}{timestamp, tagVal}
							dts = append(dts, dt)
						}

						resultChart = append(resultChart, tk.M{"name": hp.UpperFirstLetter(tag), "data": dts, "unit": mapUnit[tag]})
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
			// match.Set("avgwindspeed", tk.M{"$lte": 25})
			// match.Set("turbine", p.Turbine)

			group = tk.M{
				"_id": "$timestamp",
				// "energy":    tk.M{"$sum": "$energy"},
				"power":     tk.M{"$sum": "$fast_activepower_kw_stddevv"},
				"windspeed": tk.M{"$avg": "$fast_windspeed_ms_stddev"},
			}
		} else {
			collName = new(ScadaDataOEM).TableName()
			match.Set("dateinfo.dateid", tk.M{"$gte": tStart, "$lte": tEnd})
			// match.Set("avgwindspeed", tk.M{"$lte": 25})
			// match.Set("turbine", p.Turbine)

			group = tk.M{
				"_id":       "$timestamp",
				"power":     tk.M{"$sum": "$denpower"},
				"windspeed": tk.M{"$avg": "$denwindspeed"},
			}
		}

		pipes = append(pipes, tk.M{"$match": match})
		pipes = append(pipes, tk.M{"$group": group})
		pipes = append(pipes, tk.M{"$sort": tk.M{"_id": 1}})

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
			if len(list) > 0 {
				for _, val := range list {
					timestamp := tk.ToInt(tk.ToString(val.Get("_id").(time.Time).Unix())+"000", tk.RoundingAuto)
					var tagVal float64

					if tag == "production" {
						tagVal = val.GetFloat64(tag) / 1000
					} else if tag == "windspeed" {
						tagVal = val.GetFloat64(tag)
						if tagVal < 0 {
							tagVal = 0.0
						}
					} else {
						tagVal = val.GetFloat64(tag)
					}

					dt := []interface{}{timestamp, tagVal}
					dts = append(dts, dt)
				}

				resultChart = append(resultChart, tk.M{"name": hp.UpperFirstLetter(tag), "data": dts, "unit": mapUnit[tag]})
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
		},
	}

	return helper.CreateResult(true, data, "success")
}

func GetHFDData(turbine string, tStart time.Time, tEnd time.Time, tags []string) (result []tk.M, e error) {
	prefix := "data_"

	for {
		startStr := tStart.Format("20060102150405")
		endStr := tEnd.Format("20060102150405")

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

		if len(tmpResult) > 0 {

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
				for _, v := range mp {
					value += v
				}

				value = value / float64(len(mp))
				res.Set(n, value)
			}

			result = append(result, res)
		}
		if startStr == endStr {
			break
		}

		tStart = tStart.Add(1 * time.Second)
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
