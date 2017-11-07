package controller

import (
	. "eaciit/wfdemo-git/library/core"
	. "eaciit/wfdemo-git/library/models"
	"eaciit/wfdemo-git/web/helper"

	"sort"

	"github.com/eaciit/crowd"
	"github.com/eaciit/dbox"
	"github.com/eaciit/knot/knot.v1"
	tk "github.com/eaciit/toolkit"
	// "time"
)

func (m *AnalyticMeteorologyController) GetTurbulenceIntensity(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	p := new(PayloadAnalytic)
	e := k.GetPayload(&p)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}
	tStart, tEnd, e := helper.GetStartEndDate(k, p.Period, p.DateStart, p.DateEnd)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	// tStart, _ = time.Parse("2006-01-02 15:04:05", "2016-08-21 00:00:00")
	// tEnd, _ = time.Parse("2006-01-02 15:04:05", "2016-08-23 00:00:00")

	turbine := p.Turbine
	var (
		query        []tk.M
		querymet     []tk.M
		pipes        []tk.M
		pipesmet     []tk.M
		results      tk.M
		datas        []tk.M
		dataSeries   []tk.M
		sortTurbines []string
	)

	scadaHfds := make([]tk.M, 0)
	metTowers := make([]tk.M, 0)

	colors := []string{"#ff9933", "#21c4af", "#ff7663", "#ffb74f", "#a2df53", "#1c9ec4", "#ff63a5", "#f44336", "#69d2e7", "#8877A9", "#9A12B3", "#26C281", "#E7505A", "#C49F47", "#ff5597", "#c3260c", "#d4735e", "#ff2ad7", "#34ac8b", "#11b2eb", "#004c79", "#ff0037", "#507ca3", "#ff6565", "#ffd664", "#72aaff", "#795548", "#383271", "#6a4795", "#bec554", "#ab5919", "#f5b1e1", "#7b3416", "#002fef", "#8d731b", "#1f8805", "#ff9900", "#9C27B0", "#6c7d8a", "#d73c1c", "#5be7a0", "#da02d4", "#afa56e", "#7e32cb", "#a2eaf7", "#9cb8f4", "#9E9E9E", "#065806", "#044082", "#18937d", "#2c787a", "#a57c0c", "#234341", "#1aae7a", "#7ac610", "#736f5f", "#4e741e", "#68349d", "#1df3b6", "#e02b09", "#d9cfab", "#6e4e52", "#f31880", "#7978ec", "#f5ace8", "#3db6ae", "#5e06b0", "#16d0b9", "#a25a5b", "#1e603a", "#4b0981", "#62975f", "#1c8f2f", "#b0c80c", "#642794", "#e2060d", "#2125f0"}
	// colors := []string{"#87c5da", "#cc2a35", "#d66b76", "#5d1b62", "#f1c175", "#95204c", "#8f4bc5", "#7d287d", "#00818e", "#c8c8c8", "#546698", "#66c99a", "#f3d752", "#20adb8", "#333d6b", "#d077b1", "#aab664", "#01a278", "#c1d41a", "#807063", "#ff5975", "#01a3d4", "#ca9d08", "#026e51", "#4c653f", "#007ca7"}

	if p.Project != "" {
		query = append(query, tk.M{"projectname": tk.M{"$eq": p.Project}})
		querymet = append(querymet, tk.M{"projectname": p.Project})
	}

	query = append(query, tk.M{"_id": tk.M{"$ne": ""}})
	query = append(query, tk.M{"timestamp": tk.M{"$gte": tStart}})
	query = append(query, tk.M{"timestamp": tk.M{"$lte": tEnd}})
	query = append(query, tk.M{"windspeed_ms_bin": tk.M{"$gte": 0}})
	query = append(query, tk.M{"windspeed_ms_bin": tk.M{"$lte": 25}})
	query = append(query, tk.M{"windspeed_ms": tk.M{"$gte": -200}})
	query = append(query, tk.M{"windspeed_ms_stddev": tk.M{"$gte": -200}})

	pipes = append(pipes, tk.M{"$match": tk.M{"$and": query}})
	pipes = append(pipes, tk.M{"$group": tk.M{"_id": tk.M{
		"turbine":      "$turbine",
		"windspeedbin": "$windspeed_ms_bin"},
		"avgws":       tk.M{"$avg": "$windspeed_ms"},
		"avgwsstddev": tk.M{"$avg": "$windspeed_ms_stddev"},
	},
	})
	pipes = append(pipes, tk.M{"$sort": tk.M{"_id.windspeedbin": 1}})

	querymet = append(querymet, tk.M{"_id": tk.M{"$ne": ""}})
	querymet = append(querymet, tk.M{"timestamp": tk.M{"$gte": tStart}})
	querymet = append(querymet, tk.M{"timestamp": tk.M{"$lte": tEnd}})
	querymet = append(querymet, tk.M{"windspeedbin": tk.M{"$gte": 0}})
	querymet = append(querymet, tk.M{"windspeedbin": tk.M{"$lte": 25}})
	querymet = append(querymet, tk.M{"vhubws90mavg": tk.M{"$gte": -200}})
	querymet = append(querymet, tk.M{"vhubws90mstddev": tk.M{"$gte": -200}})

	pipesmet = append(pipesmet, tk.M{"$match": tk.M{"$and": querymet}})
	pipesmet = append(pipesmet, tk.M{"$group": tk.M{"_id": tk.M{"turbine": "Met Tower", "windspeedbin": "$windspeedbin"},
		"avgws": tk.M{"$avg": "$vhubws90mavg"}, "avgwsstddev": tk.M{"$avg": "$vhubws90mstddev"}},
	})
	pipesmet = append(pipesmet, tk.M{"$sort": tk.M{"_id.windspeedbin": 1}})

	csr, e := DB().Connection.NewQuery().
		From("Scada10MinHFD").
		Command("pipe", pipes).
		Cursor(nil)

	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	e = csr.Fetch(&scadaHfds, 0, false)

	csr.Close()

	csrt, e := DB().Connection.NewQuery().
		From(new(MetTower).TableName()).
		Command("pipe", pipesmet).
		Cursor(nil)

	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	e = csrt.Fetch(&metTowers, 0, false)

	csrt.Close()

	// tk.Printf("metTowers : %s \n", len(metTowers))
	turbineName, e := helper.GetTurbineNameList(p.Project)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	for _, m := range metTowers {
		iDs := m.Get("_id").(tk.M)
		turbine := iDs.GetString("turbine")
		windspeedbin := iDs.GetInt("windspeedbin")
		avgws := m.GetFloat64("avgws")
		avgwsstddev := m.GetFloat64("avgwsstddev")
		datas = append(datas, tk.M{
			"turbine":     turbine,
			"binws":       windspeedbin,
			"avgws":       avgws,
			"avgwsstddev": avgwsstddev,
			"tivalue":     tk.Div(avgwsstddev, avgws),
		})
	}

	for _, m := range scadaHfds {
		iDs := m.Get("_id").(tk.M)
		turbine := iDs.GetString("turbine")
		windspeedbin := iDs.GetFloat64("windspeedbin")
		avgws := m.GetFloat64("avgws")
		avgwsstddev := m.GetFloat64("avgwsstddev")
		datas = append(datas, tk.M{
			"turbine":     turbine,
			"binws":       windspeedbin,
			"avgws":       avgws,
			"avgwsstddev": avgwsstddev,
			"tivalue":     tk.Div(avgwsstddev, avgws),
		})
	}

	if len(p.Turbine) == 0 {
		for _, listVal := range datas {
			exist := false
			for _, val := range turbine {
				if listVal["turbine"] == val {
					exist = true
				}
			}
			if exist == false {
				turbine = append(turbine, listVal["turbine"])
			}
		}
	}

	for _, turX := range turbine {
		sortTurbines = append(sortTurbines, turX.(string))
	}
	sort.Strings(sortTurbines)

	selArr := 0
	for _, turbineX := range sortTurbines {

		exist := crowd.From(&datas).Where(func(x interface{}) interface{} {

			return x.(tk.M).GetString("turbine") == turbineX
		}).Exec().Result.Data().([]tk.M)

		var dts [][]float64
		turbineData := tk.M{}
		turbineData.Set("name", turbineName[turbineX])
		turbineData.Set("type", "scatterLine")
		turbineData.Set("style", "smooth")
		turbineData.Set("dashType", "solid")
		turbineData.Set("markers", tk.M{"visible": false})
		turbineData.Set("width", 2)
		turbineData.Set("color", colors[selArr])
		turbineData.Set("idxseries", selArr)

		for _, val := range exist {
			dts = append(dts, []float64{val.GetFloat64("binws"), val.GetFloat64("tivalue")})
		}

		if len(dts) > 0 {
			turbineData.Set("data", dts)
		}

		dataSeries = append(dataSeries, turbineData)
		selArr++
	}

	if datas == nil {
		datas = make([]tk.M, 0)
	}

	results = tk.M{
		"Data":        datas,
		"ChartSeries": dataSeries,
	}

	return results
}

/// getting scatter for TI
func (m *AnalyticMeteorologyController) GetTurbulenceIntensityScatter(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	p := new(PayloadAnalytic)
	e := k.GetPayload(&p)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}
	tStart, tEnd, e := helper.GetStartEndDate(k, p.Period, p.DateStart, p.DateEnd)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	var (
		filtermet []*dbox.Filter
		results   tk.M
		datas     []tk.M
	)

	metTowers := make([]tk.M, 0)

	colors := p.Color
	turbines := p.Turbine
	project := p.Project

	if project != "" {
		filtermet = append(filtermet, dbox.Eq("projectname", project))
	}

	filtermet = append(filtermet, dbox.Gte("timestamp", tStart))
	filtermet = append(filtermet, dbox.Lte("timestamp", tEnd))
	filtermet = append(filtermet, dbox.Gte("windspeedbin", 0))
	filtermet = append(filtermet, dbox.Gte("windspeedbin", 25))

	csrt, e := DB().Connection.NewQuery().
		From(new(MetTower).TableName()).Where(dbox.And(filtermet...)).
		Cursor(nil)

	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	e = csrt.Fetch(&metTowers, 0, false)

	csrt.Close()

	for _, m := range metTowers {
		iDs := m.Get("_id").(tk.M)
		turbine := "Met Tower" //iDs.GetString("turbine")
		windspeedbin := iDs.GetInt("windspeedbin")
		avgws := m.GetFloat64("avgws")
		avgwsstddev := m.GetFloat64("avgwsstddev")
		datas = append(datas, tk.M{
			"turbine":     turbine,
			"binws":       windspeedbin,
			"avgws":       avgws,
			"avgwsstddev": avgwsstddev,
			"tivalue":     tk.Div(avgwsstddev, avgws),
		})
	}

	selArr := 0
	for _, t := range turbines {
		currColor := colors[selArr]

		var filter []*dbox.Filter
		var scadaHfds []tk.M

		tb := t.(string)

		filter = append(filter, dbox.Eq("projectname", project))
		filter = append(filter, dbox.Eq("turbine", tb))
		filter = append(filter, dbox.Gte("timestamp", tStart))
		filter = append(filter, dbox.Lte("timestamp", tEnd))
		filter = append(filter, dbox.Gte("windspeed_ms_bin", 0))
		filter = append(filter, dbox.Lte("windspeed_ms_bin", 25))

		csr, e := DB().Connection.NewQuery().
			From("Scada10MinHFD").Where(dbox.And(filter...)).
			Cursor(nil)

		if e != nil {
			return helper.CreateResult(false, nil, e.Error())
		}

		e = csr.Fetch(&scadaHfds, 0, false)
		if e != nil {
			return helper.CreateResult(false, nil, e.Error())
		}
		csr.Close()

		item := tk.M{}
		item.Set("colorField", "valueColor")
		item.Set("markers", tk.M{}.Set("size", 2))
		item.Set("name", "Scatter-"+tb)
		item.Set("type", "scatter")
		item.Set("xField", "WindSpeed")
		item.Set("yField", "TurbulenceIntensity")

		dataItem := []tk.M{}
		for _, m := range scadaHfds {
			ws := m.GetFloat64("windspeed_ms")
			wsBin := m.GetFloat64("windspeed_ms_bin")
			wsStdDev := m.GetFloat64("windspeed_ms_stddev")
			tiVal := tk.Div(wsStdDev, ws)

			dt := tk.M{}
			dt.Set("valueColor", currColor)
			dt.Set("WindSpeed", ws)
			dt.Set("WindSpeedBin", wsBin)
			dt.Set("WindSpeedStdDev", wsStdDev)
			dt.Set("TurbulenceIntensity", tiVal)

			dataItem = append(dataItem, dt)
		}
		item.Set("data", dataItem)

		datas = append(datas, item)

		selArr++
	}

	results = tk.M{
		"data":    tk.M{"Data": datas},
		"message": "success",
		"success": true,
	}

	return results
}
