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

type ScadaAnalyticsWDDataGroup struct {
	Turbine  string
	Category float64
}

const (
	minWS       = 0.5
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
	windCats    []float64
	nacelleCats []float64
)

func getWindDistrCategory(windValue float64) float64 {
	var datas float64

	for _, val := range windCats {
		if val >= windValue {
			datas = val
			return datas
		}
	}

	return datas
}

type ScadaAnalyticsWDData struct {
	Turbine    string
	Category   float64
	Contribute float64
}

func setContribution(turbine, tipe string, dataCatCount map[string]float64, countPerWSCat float64) (results []ScadaAnalyticsWDData) {
	results = []ScadaAnalyticsWDData{}
	category := []float64{}
	switch tipe {
	case "nacelledeviation":
		category = nacelleCats
	case "avgwindspeed":
		category = windCats
	}

	for _, val := range category {
		results = append(results, ScadaAnalyticsWDData{
			Turbine:    turbine,
			Category:   val,
			Contribute: tk.Div(dataCatCount[tk.ToString(val)], countPerWSCat),
		})
	}
	return
}

func GetMetTowerData(p *PayloadAnalytic, k *knot.WebContext) []ScadaAnalyticsWDData {
	dataSeries := []ScadaAnalyticsWDData{}

	tStart, tEnd, e := helper.GetStartEndDate(k, p.Period, p.DateStart, p.DateEnd)

	type MiniMetTower struct {
		VHubWS90mAvg float64
	}

	queryT := []*dbox.Filter{}
	queryT = append(queryT, dbox.Gte("dateinfo.dateid", tStart))
	queryT = append(queryT, dbox.Lte("dateinfo.dateid", tEnd))
	queryT = append(queryT, dbox.Gte("vhubws90mavg", 0.5))

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
	_data := MiniMetTower{}
	for {
		_data = MiniMetTower{}
		e = csrData.Fetch(&_data, 1, false)
		if e != nil {
			break
		}
		countPerWSCat++
		if _data.VHubWS90mAvg > maxWS {
			_data.VHubWS90mAvg = maxWS
		}
		modus = math.Mod(_data.VHubWS90mAvg, stepWS)
		if modus == 0 {
			category = _data.VHubWS90mAvg
		} else {
			category = _data.VHubWS90mAvg - modus + stepWS
		}
		groupKey = tk.ToString(category)
		dataCatCount[groupKey] = dataCatCount[groupKey] + 1
	}
	csrData.Close()
	dataSeries = append(dataSeries, setContribution("Met Tower", "avgwindspeed", dataCatCount, countPerWSCat)...)

	return dataSeries
}

func GetScadaData(turbineName map[string]string, turbineNameSorted []string, queryT []*dbox.Filter, tipe string) ([]ScadaAnalyticsWDData, []string) {
	var e error
	fieldName := ""
	maxStep := 0.0
	step := 0.0
	switch tipe {
	case "nacelledeviation":
		fieldName = "nacelledeviation"
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
	dataSeries := []ScadaAnalyticsWDData{}
	dataSeriesTempMap := map[string][]ScadaAnalyticsWDData{}
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
				dataSeriesTempMap[lastTurbine] = setContribution(lastTurbine, tipe, dataCatCount, countPerWSCat)
			}
			dataCatCount = map[string]float64{}
			lastTurbine = _turbine
			countPerWSCat = 0.0
		}
		countPerWSCat++
		if _data.GetFloat64(fieldName) > maxStep {
			_data.Set(fieldName, maxStep)
		}
		modus = math.Mod(_data.GetFloat64(fieldName), step)
		if modus == 0 {
			category = _data.GetFloat64(fieldName)
		} else {
			category = _data.GetFloat64(fieldName) - modus + step
		}
		groupKey = tk.ToString(category)
		dataCatCount[groupKey] = dataCatCount[groupKey] + 1
	}
	if lastTurbine != "" {
		dataSeriesTempMap[lastTurbine] = setContribution(lastTurbine, tipe, dataCatCount, countPerWSCat)
	}
	csrData.Close()
	turbineAvail := []string{}
	for _, _turbinename := range turbineNameSorted {
		if dataSeriesVal, hasKey := dataSeriesTempMap[_turbinename]; hasKey {
			dataSeries = append(dataSeries, dataSeriesVal...)
			turbineAvail = append(turbineAvail, _turbinename)
		}
	}

	return dataSeries, turbineAvail
}

func (m *AnalyticWindDistributionController) GetList(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	dataSeries := []ScadaAnalyticsWDData{}

	p := new(PayloadAnalytic)
	e := k.GetPayload(&p)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	tStart, tEnd, e := helper.GetStartEndDate(k, p.Period, p.DateStart, p.DateEnd)
	query := []tk.M{}
	pipes := []tk.M{}
	query = append(query, tk.M{"_id": tk.M{"$ne": ""}})
	query = append(query, tk.M{"dateinfo.dateid": tk.M{"$gte": tStart}})
	query = append(query, tk.M{"dateinfo.dateid": tk.M{"$lte": tEnd}})
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
		query = append(query, tk.M{"nacelledeviation": tk.M{"$gte": -180}})
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
		query = append(query, tk.M{"avgwindspeed": tk.M{"$gte": 0.5}})
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
		queryT = append(queryT, dbox.Gte("nacelledeviation", -180))
	case "avgwindspeed":
		queryT = append(queryT, dbox.Gte("avgwindspeed", 0.5))
	}
	queryT = append(queryT, dbox.Eq("available", 1))
	if p.Project != "" {
		queryT = append(queryT, dbox.Eq("projectname", p.Project))
	}
	queryT = append(queryT, dbox.In("turbine", turbineInt...))
	turbineAvail := []string{"Met Tower"}
	turbineAvailTemp := []string{}

	if p.Project == "Tejuva" && p.BreakDown == "avgwindspeed" {
		dataMetTower := []ScadaAnalyticsWDData{}
		dataScada := []ScadaAnalyticsWDData{}
		var wg sync.WaitGroup
		wg.Add(2)
		go func() {
			dataMetTower = GetMetTowerData(p, k)
			wg.Done()
		}()
		go func() {
			dataScada, turbineAvailTemp = GetScadaData(turbineName, turbineNameSorted, queryT, p.BreakDown)
			wg.Done()
		}()
		wg.Wait()
		dataSeries = append(dataMetTower, dataScada...)
	} else {
		dataSeries, turbineAvail = GetScadaData(turbineName, turbineNameSorted, queryT, p.BreakDown)
	}
	turbineAvail = append(turbineAvail, turbineAvailTemp...)

	data := struct {
		Data        []ScadaAnalyticsWDData
		TurbineList []string
	}{
		Data:        dataSeries,
		TurbineList: turbineAvail,
	}

	return helper.CreateResult(true, data, "success")
}
