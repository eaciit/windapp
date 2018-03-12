package controller

import (
	. "eaciit/wfdemo-git/library/core"
	. "eaciit/wfdemo-git/library/models"
	"eaciit/wfdemo-git/web/helper"

	"sort"

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

	turbine := p.Turbine
	turbine = append(turbine, "") /* buat met tower */
	var (
		query      []tk.M
		pipes      []tk.M
		results    tk.M
		dataSeries []tk.M
	)

	colors := []string{"#ff9933", "#21c4af", "#ff7663", "#ffb74f", "#a2df53", "#1c9ec4", "#ff63a5", "#f44336", "#69d2e7", "#8877A9", "#9A12B3", "#26C281", "#E7505A", "#C49F47", "#ff5597", "#c3260c", "#d4735e", "#ff2ad7", "#34ac8b", "#11b2eb", "#004c79", "#ff0037", "#507ca3", "#ff6565", "#ffd664", "#72aaff", "#795548", "#383271", "#6a4795", "#bec554", "#ab5919", "#f5b1e1", "#7b3416", "#002fef", "#8d731b", "#1f8805", "#ff9900", "#9C27B0", "#6c7d8a", "#d73c1c", "#5be7a0", "#da02d4", "#afa56e", "#7e32cb", "#a2eaf7", "#9cb8f4", "#9E9E9E", "#065806", "#044082", "#18937d", "#2c787a", "#a57c0c", "#234341", "#1aae7a", "#7ac610", "#736f5f", "#4e741e", "#68349d", "#1df3b6", "#e02b09", "#d9cfab", "#6e4e52", "#f31880", "#7978ec", "#f5ace8", "#3db6ae", "#5e06b0", "#16d0b9", "#a25a5b", "#1e603a", "#4b0981", "#62975f", "#1c8f2f", "#b0c80c", "#642794", "#e2060d", "#2125f0"}

	if p.Project != "" {
		query = append(query, tk.M{"projectname": p.Project})
	}
	if len(turbine) > 0 {
		query = append(query, tk.M{"turbine": tk.M{"$in": turbine}})
	}

	query = append(query, tk.M{"timestamp": tk.M{"$gte": tStart}})
	query = append(query, tk.M{"timestamp": tk.M{"$lte": tEnd}})

	pipes = append(pipes, tk.M{"$match": tk.M{"$and": query}})
	pipes = append(pipes, tk.M{"$group": tk.M{
		"_id": tk.M{
			"turbine":      "$turbine",
			"windspeedbin": "$windspeedbin",
		},
		"windspeedtotal":    tk.M{"$sum": "$windspeedtotal"},
		"windspeedstdtotal": tk.M{"$sum": "$windspeedstdtotal"},
		"windspeedcount":    tk.M{"$sum": "$windspeedcount"},
		"windspeedstdcount": tk.M{"$sum": "$windspeedstdcount"},
	},
	})
	pipes = append(pipes, tk.M{"$sort": tk.M{"_id.windspeedbin": 1}})

	csr, e := DB().Connection.NewQuery().
		From("rpt_turbulenceintensity").
		Command("pipe", pipes).
		Cursor(nil)

	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	dataPerTurbine := map[string][][]float64{}

	for {
		item := tk.M{}
		e = csr.Fetch(&item, 1, false)
		if e != nil {
			break
		}
		ids := item.Get("_id", tk.M{}).(tk.M)
		_turbine := ids.GetString("turbine")
		if _turbine == "" {
			_turbine = "Met Tower"
		}
		windspeedbin := ids.GetFloat64("windspeedbin")
		sumws := item.GetFloat64("windspeedtotal")
		sumwsstddev := item.GetFloat64("windspeedstdtotal")
		countws := item.GetFloat64("windspeedcount")
		countwsstddev := item.GetFloat64("windspeedstdcount")
		avgws := tk.Div(sumws, countws)
		avgwsstddev := tk.Div(sumwsstddev, countwsstddev)
		tiValue := tk.Div(avgwsstddev, avgws)
		dataPerTurbine[_turbine] = append(dataPerTurbine[_turbine], []float64{windspeedbin, tiValue})
	}
	csr.Close()

	turbineName, e := helper.GetTurbineNameList(p.Project)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	sortTurbine := []string{}
	isMetExists := false
	for _turbine := range dataPerTurbine {
		if _turbine == "Met Tower" {
			isMetExists = true
			continue
		}
		sortTurbine = append(sortTurbine, _turbine)
	}
	sort.Strings(sortTurbine)
	if isMetExists {
		sortTurbine = append([]string{"Met Tower"}, sortTurbine...)
	}

	selArr := 0
	dataSeriesMet := []tk.M{}
	for _, turbineX := range sortTurbine {
		turbineData := tk.M{}
		nameTurbine := turbineName[turbineX]
		if nameTurbine == "" {
			nameTurbine = "Met Tower"
		}
		turbineData.Set("name", nameTurbine)
		turbineData.Set("turbineid", turbineX)
		turbineData.Set("type", "scatterLine")
		turbineData.Set("style", "smooth")
		turbineData.Set("dashType", "solid")
		turbineData.Set("markers", tk.M{"visible": false})
		turbineData.Set("width", 2)
		turbineData.Set("color", colors[selArr])
		turbineData.Set("idxseries", selArr)
		data := dataPerTurbine[turbineX]

		if len(data) > 0 {
			turbineData.Set("data", data)
		}

		if nameTurbine == "Met Tower" {
			dataSeriesMet = append(dataSeriesMet, turbineData)
		} else {
			dataSeries = append(dataSeries, turbineData)
		}
		selArr++
	}

	results = tk.M{
		"ChartSeries": dataSeries,
		"MetSeries":   dataSeriesMet,
		"ColorSeries": colors,
	}

	return results
}

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

	defer csrt.Close()

	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	e = csrt.Fetch(&metTowers, 0, false)

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
		filter = append(filter, dbox.Gte("windspeed_ms", -200))
		filter = append(filter, dbox.Gte("windspeed_ms_stddev", -200))

		filter = append(filter, dbox.Ne("windspeed_ms", nil))
		filter = append(filter, dbox.Ne("activepower_kw", nil))

		csr, e := DB().Connection.NewQuery().
			From("Scada10MinHFD").Where(dbox.And(filter...)).
			Cursor(nil)

		if e != nil {
			if csr != nil {
				csr.Close()
			}
			return helper.CreateResult(false, nil, e.Error())
		}

		e = csr.Fetch(&scadaHfds, 0, false)
		if e != nil {
			if csr != nil {
				csr.Close()
			}
			return helper.CreateResult(false, nil, e.Error())
		}
		csr.Close()

		tk.Println("Scada HFD Length", len(scadaHfds))

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
