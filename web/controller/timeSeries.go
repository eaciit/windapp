package controller

import (
	"bufio"
	. "eaciit/wfdemo-git/library/core"
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

func CreateTimeSeriesController() *TimeSeriesController {
	var controller = new(TimeSeriesController)
	return controller
}

func (m *TimeSeriesController) GetData(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	var (
		pipes []tk.M
		list  []tk.M
	)

	p := new(PayloadAnalytic)
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
		intTimestamp := val.Get("_id").(time.Time).UTC().Unix()

		energy := val.GetFloat64("energy") / 1000
		wind := val.GetFloat64("windspeed")

		dtProd = append(dtProd, []interface{}{intTimestamp, energy})
		dtWS = append(dtWS, []interface{}{intTimestamp, wind})
	}

	result := []tk.M{}

	result = append(result, tk.M{"name": "Production", "data": dtProd, "unit": "MWh"})
	result = append(result, tk.M{"name": "Windspeed", "data": dtWS, "unit": "m/s"})

	data := struct {
		Data []tk.M
	}{
		Data: result,
	}

	return helper.CreateResult(true, data, "success")
}

func (m *TimeSeriesController) GetAvailDate(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	return helper.CreateResult(true, k.Session("availdate", ""), "success")
}

func (m *TimeSeriesController) GetDataHFD(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	var (
		pipes       []tk.M
		list        []tk.M
		resultChart []tk.M
		periodList  []tk.M
	)

	p := new(PayloadAnalytic)
	e := k.GetPayload(&p)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	tStart, tEnd, e := helper.GetStartEndDate(k, p.Period, p.DateStart, p.DateEnd)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	dataType := "seconds"

	if dataType == "seconds" {
		mapUnit := map[string]string{}
		mapUnit["WindSpeed_ms"] = "m/s"
		mapUnit["ActivePower_kW"] = "kW"

		for {
			// log.Printf(">> %v \n", tEnd.UTC().Sub(tStart.UTC()).Seconds())
			// log.Printf(">> %v \n", tStart.UTC().Sub(tEnd.UTC()).Seconds())

			if tStart.UTC().Sub(tEnd.UTC()).Seconds() >= 0 {
				break
			}

			// log.Printf(">>> %v \n", tStart.String())

			before := tStart.UTC()
			tStart = tStart.UTC().Add(time.Duration(3) * time.Hour)
			periodList = append(periodList, tk.M{"starttime": before, "endtime": tStart.UTC()})
		}

		// log.Printf("> %#v \n", periodList)

		if len(periodList) > 0 {
			current := periodList[0]
			currStar := current.Get("starttime").(time.Time)
			currEnd := current.Get("endtime").(time.Time)

			tags := []string{"WindSpeed_ms", "ActivePower_kW"}
			hdfs, e := GetHFDData("HBR004", currStar, currEnd, tags)

			if e != nil {
				return helper.CreateResult(false, nil, e.Error())
			}

			projectName := ""
			_ = projectName
			if p.Project != "" {
				anProject := strings.Split(p.Project, "(")
				projectName = strings.TrimRight(anProject[0], " ")
			}

			for _, val := range hdfs {
				timestamp := val.Get("timestamp").(time.Time).Unix()
				for _, tag := range tags {
					tagVal := val.GetFloat64(tag)
					dt := []interface{}{timestamp, tagVal}
					resultChart = append(resultChart, tk.M{"name": tag, "data": dt, "unit": mapUnit[tag]})
				}
			}

			for _, tag := range tags {
				var dts [][]interface{}
				for _, val := range hdfs {
					timestamp := val.Get("timestamp").(time.Time).Unix()
					tagVal := val.GetFloat64(tag)
					dt := []interface{}{timestamp, tagVal}
					dts = append(dts, dt)
				}

				resultChart = append(resultChart, tk.M{"name": tag, "data": dts, "unit": mapUnit[tag]})
			}
		}
	} else {
		match := tk.M{}

		match.Set("dateinfo.dateid", tk.M{"$gte": tStart, "$lte": tEnd})
		// match.Set("avgwindspeed", tk.M{"$lte": 25})

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
			From(new(ScadaDataHFD).TableName()).
			Command("pipe", pipes).
			Cursor(nil)

		if e != nil {
			return helper.CreateResult(false, nil, e.Error())
		}

		e = csr.Fetch(&list, 0, false)

		csr.Close()

		var dtProd, dtWS [][]interface{}

		for _, val := range list {
			intTimestamp := val.Get("_id").(time.Time).UTC().Unix()

			energy := val.GetFloat64("energy") / 1000
			wind := val.GetFloat64("windspeed")

			dtProd = append(dtProd, []interface{}{intTimestamp, energy})
			dtWS = append(dtWS, []interface{}{intTimestamp, wind})
		}

		resultChart = append(resultChart, tk.M{"name": "Production", "data": dtProd, "unit": "MWh"})
		resultChart = append(resultChart, tk.M{"name": "Windspeed", "data": dtWS, "unit": "m/s"})
	}

	type ResDataAvail struct {
		Chart      []tk.M
		PeriodList []tk.M
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

		// groupTime := hpp.GenNext10Minutes(newTime.UTC()).UTC()
		// log.Printf(">> %v | %v \n", newTime, tStart.UTC())
		folder := newTime.Format("20060102/15/") + newTime.Format("1504")
		file := prefix + startStr + ".csv"
		path := "/Volumes/Development/ostrorealtime/" + folder + "/" + file
		tmpResult, err := ReadHFDFile(path, tags)

		if err != nil {
			// log.Printf("Err: %v \n", err.Error())
		}

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
