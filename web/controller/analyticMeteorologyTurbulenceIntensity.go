package controller

import (
	. "eaciit/wfdemo-git/library/core"
	. "eaciit/wfdemo-git/library/models"
	"eaciit/wfdemo-git/web/helper"

	"github.com/eaciit/knot/knot.v1"
	tk "github.com/eaciit/toolkit"
	"github.com/eaciit/crowd"
	"sort"
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
		query    []tk.M
		pipes    []tk.M
		pipesmet []tk.M
		results  tk.M
		datas    []tk.M
		dataSeries []tk.M
		sortTurbines []string
	)

	scadaHfds := make([]tk.M, 0)
	metTowers := make([]tk.M, 0)

	colors := []string{
		"#ED1C24", "#A3238E", "#00A65D", "#F58220", "#0066B3", "#5C2D91", "#FFF200", "#579835", "#CF3834", "#00B274", "#74489D",
		"#C06616", "#5565AF", "#CCBE00", "#390A5D", "#006D6F", "#65C295", "#F04E4D", "#407927", "#00599D", "#A09600", "#0D1F63",
		"#C38312", "#003D73", "#454FA1", "#BC312E",
	}

	query = append(query, tk.M{"_id": tk.M{"$ne": ""}})
	query = append(query, tk.M{"timestamp": tk.M{"$gte": tStart}})
	query = append(query, tk.M{"timestamp": tk.M{"$lte": tEnd}})

	pipes = append(pipes, tk.M{"$match": tk.M{"$and": query}})
	pipes = append(pipes, tk.M{"$group": tk.M{"_id": tk.M{"turbine": "$turbine", "windspeedbin": "$fast_windspeed_bin"},
		"avgws": tk.M{"$avg": "$fast_windspeed_ms"}, "avgwsstddev": tk.M{"$avg": "$fast_activepower_kw_stddev"}},
	})
	pipes = append(pipes, tk.M{"$sort": tk.M{"_id": 1}})

	pipesmet = append(pipesmet, tk.M{"$match": tk.M{"$and": query}})
	pipesmet = append(pipesmet, tk.M{"$group": tk.M{"_id": tk.M{"turbine": "Met Tower", "windspeedbin": "$windspeedbin"},
		"avgws": tk.M{"$avg": "$vhubws90mavg"}, "avgwsstddev": tk.M{"$avg": "$vhubws90mstddev"}},
	})
	pipesmet = append(pipesmet, tk.M{"$sort": tk.M{"_id": 1}})

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

	tk.Printf("metTowers : %s \n", len(metTowers))

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
