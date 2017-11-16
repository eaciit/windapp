package controller

import (
	. "eaciit/wfdemo-git/library/core"
	. "eaciit/wfdemo-git/library/models"
	"eaciit/wfdemo-git/web/helper"
	"math"
	// "time"
	// "fmt"
	"github.com/eaciit/dbox"
	// "runtime"
	"sort"
	"sync"

	"github.com/eaciit/knot/knot.v1"
	tk "github.com/eaciit/toolkit"
)

type AnalyticWindDistributionController struct {
	App
}

const (
	minWS       = 1.0
	maxWS       = 15.0
	stepWS      = 0.5
	minNacelle  = -180.0
	maxNacelle  = 180.0
	stepNacelle = 15.0
)

func CreateAnalyticWindDistributionController() *AnalyticWindDistributionController {
	var controller = new(AnalyticWindDistributionController)
	return controller
}

var (
	windCats    []float64 /* step nya muncul 1, 1.5, 2, 2.5, 3, ... dst */
	nacelleCats []float64 /* step nya muncul -180, -165, -150, -135, -120, ... dst */
	colorList   []string  /* color list dari UI supaya sinkron */
)

func setContribution(turbine, tipe string, dataCatCount map[string]float64, countPerWSCat float64) (results []float64) {
	results = []float64{}
	category := []float64{}
	switch tipe { /* pemilihan category sesuai tipe nya */
	case "nacelledeviation":
		category = nacelleCats
	case "avgwindspeed":
		category = windCats
	}

	for _, val := range category {
		results = append(results, tk.Div(dataCatCount[tk.ToString(val)], countPerWSCat))
	}
	return
}

func setSeries(turbinename, color string, index int, data []float64) tk.M {
	return tk.M{
		"name":  turbinename,
		"color": color,
		"style": "smooth",
		"width": 2,
		"type":  "line",
		"index": index,
		"data":  data,
		"markers": tk.M{
			"visible": false,
			"size":    3,
		},
	}
}

func GetMetTowerData(p *PayloadAnalytic, k *knot.WebContext) []tk.M {
	dataSeries := []tk.M{}

	tStart, tEnd, e := helper.GetStartEndDate(k, p.Period, p.DateStart, p.DateEnd)

	queryT := []*dbox.Filter{}
	queryT = append(queryT, dbox.Gte("dateinfo.dateid", tStart))
	queryT = append(queryT, dbox.Lte("dateinfo.dateid", tEnd))
	queryT = append(queryT, dbox.Gt("vhubws90mavg", 0.5))

	if p.Project != "" {
		queryT = append(queryT, dbox.Eq("projectname", p.Project))
	}

	csrData, _ := DB().Connection.NewQuery().
		Select("vhubws90mavg").
		From(new(MetTower).TableName()).
		Where(dbox.And(queryT...)).
		Order("turbine").
		Cursor(nil)

	groupKey := ""
	countPerWSCat := 0.0
	dataCatCount := map[string]float64{}
	category := 0.0
	modus := 0.0
	_data := tk.M{}
	for {
		_data = tk.M{}
		e = csrData.Fetch(&_data, 1, false)
		if e != nil {
			break
		}
		countPerWSCat++
		if _data.GetFloat64("vhubws90mavg") > maxWS {
			_data.Set("vhubws90mavg", maxWS)
		}
		modus = math.Mod(_data.GetFloat64("vhubws90mavg"), stepWS)
		if modus == 0 {
			category = _data.GetFloat64("vhubws90mavg")
		} else {
			category = _data.GetFloat64("vhubws90mavg") - modus + stepWS
		}
		groupKey = tk.ToString(category)
		dataCatCount[groupKey] = dataCatCount[groupKey] + 1
	}
	csrData.Close()
	dataSeriesVal := setContribution("Met Tower", "avgwindspeed", dataCatCount, countPerWSCat)
	dataSeries = append(dataSeries, setSeries("Met Tower", colorList[0], 0, dataSeriesVal))

	return dataSeries
}

func GetScadaData(turbineName map[string]string, turbineNameSorted []string, queryT []*dbox.Filter, tipe string, withMet bool) ([]tk.M, []string) {
	var e error
	fieldName := ""
	maxStep := 0.0
	step := 0.0
	switch tipe {
	case "nacelledeviation":
		fieldName = "winddirection"
		maxStep = maxNacelle
		step = stepNacelle
	case "avgwindspeed":
		fieldName = "avgwindspeed"
		maxStep = maxWS
		step = stepWS
	}

	csrData, _ := DB().Connection.NewQuery().
		Select("turbine", fieldName).
		From(new(ScadaData).TableName()).
		Where(dbox.And(queryT...)).
		Order("turbine").
		Cursor(nil)

	lastTurbine := ""
	_turbine := ""
	groupKey := ""
	countPerWSCat := 0.0
	dataCatCount := map[string]float64{}
	category := 0.0
	modus := 0.0
	dataSeries := []tk.M{}
	dataSeriesPerTurbine := map[string][]float64{}
	_data := tk.M{}
	for {
		_data = tk.M{}
		e = csrData.Fetch(&_data, 1, false)
		if e != nil {
			break
		}
		_turbine = turbineName[_data.GetString("turbine")]
		if lastTurbine != _turbine {
			if lastTurbine != "" {
				dataSeriesPerTurbine[lastTurbine] = setContribution(lastTurbine, tipe, dataCatCount, countPerWSCat)
			}
			dataCatCount = map[string]float64{}
			lastTurbine = _turbine
			countPerWSCat = 0.0
		}
		if _data.Has(fieldName) {
			countPerWSCat++
			if _data.GetFloat64(fieldName) > maxStep { /* jika datanya melebihi max step, maka ubah menjadi max step*/
				_data.Set(fieldName, maxStep)
			}
			modus = math.Mod(_data.GetFloat64(fieldName), step)
			if modus == 0 { /* jika habis dibagi step maka value itu sendiri yang di assign*/
				category = _data.GetFloat64(fieldName)
			} else { /* jika tidak habis dibagi step maka diikutkan value + step setelahnya */
				category = _data.GetFloat64(fieldName) - modus + step
			}
			groupKey = tk.ToString(category)
			dataCatCount[groupKey] = dataCatCount[groupKey] + 1
		}
	}
	if lastTurbine != "" {
		dataSeriesPerTurbine[lastTurbine] = setContribution(lastTurbine, tipe, dataCatCount, countPerWSCat)
	}
	csrData.Close()
	turbineAvail := []string{}
	index := 0
	if withMet {
		index = 1
	}
	for _, _turbinename := range turbineNameSorted {
		if dataSeriesVal, hasKey := dataSeriesPerTurbine[_turbinename]; hasKey {
			dataSeries = append(dataSeries, setSeries(_turbinename, colorList[index], index, dataSeriesVal))
			turbineAvail = append(turbineAvail, _turbinename)
			index++
		}
	}

	return dataSeries, turbineAvail
}

func (m *AnalyticWindDistributionController) GetList(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	dataSeries := []tk.M{}

	p := new(PayloadAnalytic)
	e := k.GetPayload(&p)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}
	for _, color := range p.Color {
		colorList = append(colorList, tk.ToString(color))
	}

	tStart, tEnd, e := helper.GetStartEndDate(k, p.Period, p.DateStart, p.DateEnd)
	query := []tk.M{}
	pipes := []tk.M{}
	query = append(query, tk.M{"_id": tk.M{"$ne": ""}})
	query = append(query, tk.M{"dateinfo.dateid": tk.M{"$gte": tStart}})
	query = append(query, tk.M{"dateinfo.dateid": tk.M{"$lte": tEnd}})

	category := []float64{}
	switch p.BreakDown {
	case "nacelledeviation":
		nacelleCats = []float64{}
		start := minNacelle
		for {
			if start > maxNacelle {
				break
			}
			nacelleCats = append(nacelleCats, start)
			start += stepNacelle
		}
		category = nacelleCats
		query = append(query, tk.M{"winddirection": tk.M{"$gte": -180}})
	case "avgwindspeed":
		windCats = []float64{}
		start := minWS
		for {
			if start > maxWS {
				break
			}
			windCats = append(windCats, start)
			start += stepWS
		}
		category = windCats
		query = append(query, tk.M{"avgwindspeed": tk.M{"$gt": 0.5}})
	}
	query = append(query, tk.M{"available": 1})
	if p.Project != "" {
		query = append(query, tk.M{"projectname": p.Project})
	}

	turbine := []string{}
	if len(p.Turbine) == 0 {
		pipes = append(pipes, tk.M{"$match": tk.M{"$and": query}})
		pipes = append(pipes, tk.M{"$group": tk.M{"_id": "$turbine"}})

		csr, _ := DB().Connection.NewQuery().From(new(ScadaData).TableName()).
			Command("pipe", pipes).Cursor(nil)
		_turbine := map[string]string{}
		for {
			e = csr.Fetch(&_turbine, 1, false)
			if e != nil {
				break
			}
			turbine = append(turbine, _turbine["_id"])
		}
		csr.Close()
	} else {
		bufferTurbine := []string{}
		for _, val := range p.Turbine {
			bufferTurbine = append(bufferTurbine, val.(string))
		}
		turbine = append(turbine, bufferTurbine...)
	}
	sort.Strings(turbine)
	turbineInt := []interface{}{}
	for _, val := range turbine {
		turbineInt = append(turbineInt, val)
	}
	turbineName, e := helper.GetTurbineNameList(p.Project)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	turbineNameSorted := []string{}
	for _, _turbinename := range turbineName {
		turbineNameSorted = append(turbineNameSorted, _turbinename)
	}
	sort.Strings(turbineNameSorted)
	turbineNameSorted = append([]string{"Met Tower"}, turbineNameSorted...)

	queryT := []*dbox.Filter{}
	queryT = append(queryT, dbox.Gte("dateinfo.dateid", tStart))
	queryT = append(queryT, dbox.Lte("dateinfo.dateid", tEnd))
	switch p.BreakDown {
	case "nacelledeviation":
		queryT = append(queryT, dbox.Gte("winddirection", -180))
	case "avgwindspeed":
		queryT = append(queryT, dbox.Gt("avgwindspeed", 0.5))
	}
	queryT = append(queryT, dbox.Eq("available", 1))
	if p.Project != "" {
		queryT = append(queryT, dbox.Eq("projectname", p.Project))
	}
	queryT = append(queryT, dbox.In("turbine", turbineInt...))
	turbineAvail := []string{"Met Tower"}
	turbineAvailTemp := []string{}

	if p.Project == "Tejuva" && p.BreakDown == "avgwindspeed" {
		dataMetTower := []tk.M{}
		dataScada := []tk.M{}
		var wg sync.WaitGroup
		wg.Add(2)
		go func() {
			dataMetTower = GetMetTowerData(p, k)
			wg.Done()
		}()
		go func() {
			dataScada, turbineAvailTemp = GetScadaData(turbineName, turbineNameSorted, queryT, p.BreakDown, true)
			wg.Done()
		}()
		wg.Wait()
		dataSeries = append(dataMetTower, dataScada...)
		turbineAvail = append(turbineAvail, turbineAvailTemp...)
	} else {
		dataSeries, turbineAvail = GetScadaData(turbineName, turbineNameSorted, queryT, p.BreakDown, false)
	}

	data := struct {
		Data        []tk.M
		TurbineList []string
		Categories  []float64
	}{
		Data:        dataSeries,
		TurbineList: turbineAvail,
		Categories:  category,
	}

	return helper.CreateResult(true, data, "success")
}
