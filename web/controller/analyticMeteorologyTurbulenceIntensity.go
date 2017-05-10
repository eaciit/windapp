package controller

import (
	. "eaciit/wfdemo-git/library/core"
	. "eaciit/wfdemo-git/library/models"
	"eaciit/wfdemo-git/web/helper"

	"sort"

	"github.com/eaciit/crowd"
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

	colors := []string{"#87c5da", "#cc2a35", "#d66b76", "#5d1b62", "#f1c175", "#95204c", "#8f4bc5", "#7d287d", "#00818e", "#c8c8c8", "#546698", "#66c99a", "#f3d752", "#20adb8", "#333d6b", "#d077b1", "#aab664", "#01a278", "#c1d41a", "#807063", "#ff5975", "#01a3d4", "#ca9d08", "#026e51", "#4c653f", "#007ca7"}

	if p.Project != "" {
		query = append(query, tk.M{"projectname": tk.M{"$eq": p.Project}})
	}

	query = append(query, tk.M{"_id": tk.M{"$ne": ""}})
	query = append(query, tk.M{"timestamp": tk.M{"$gte": tStart}})
	query = append(query, tk.M{"timestamp": tk.M{"$lte": tEnd}})
	query = append(query, tk.M{"fast_windspeed_bin": tk.M{"$gte": 0}})
	query = append(query, tk.M{"fast_windspeed_bin": tk.M{"$lte": 25}})
	query = append(query, tk.M{"fast_windspeed_ms": tk.M{"$gte": -200}})
	query = append(query, tk.M{"fast_windspeed_ms_stddev": tk.M{"$gte": -200}})

	pipes = append(pipes, tk.M{"$match": tk.M{"$and": query}})
	pipes = append(pipes, tk.M{"$group": tk.M{"_id": tk.M{
		"turbine":      "$turbine",
		"windspeedbin": "$fast_windspeed_bin"},
		"avgws":       tk.M{"$avg": "$fast_windspeed_ms"},
		"avgwsstddev": tk.M{"$avg": "$fast_windspeed_ms_stddev"},
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
		From(new(ScadaDataHFD).TableName()).
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
		turbineData.Set("name", turbineX)
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
