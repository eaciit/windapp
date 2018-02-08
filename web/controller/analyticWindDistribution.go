package controller

import (
	. "eaciit/wfdemo-git/library/core"
	. "eaciit/wfdemo-git/library/models"
	"eaciit/wfdemo-git/web/helper"
	"math"
	"time"
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
		// tk.Println(val, "=>", dataCatCount[tk.ToString(val)], "/", countPerWSCat)
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
	fieldName := "vhubws90mavg"

	csrData, _ := DB().Connection.NewQuery().
		Select(fieldName).
		From(new(MetTower).TableName()).
		Where(dbox.And(queryT...)).
		Order("turbine").
		Cursor(nil)

	groupKey := ""
	countPerWSCat := 0.0
	dataCatCount := map[string]float64{}
	category := 0.0
	modus := 0.0
	step := stepWS
	bin := step / 2
	maxStep := maxWS
	_data := tk.M{}
	for {
		_data = tk.M{}
		e = csrData.Fetch(&_data, 1, false)
		if e != nil {
			break
		}

		if _data.Has(fieldName) {
			countPerWSCat++
			value := _data.GetFloat64(fieldName)
			if value > maxStep { /* jika datanya melebihi max step, maka ubah menjadi max step*/
				_data.Set(fieldName, maxStep)
			}

			valueBin := value + bin

			if value < 0 {
				valueBin = value - bin
			}

			modus = math.Mod(valueBin, step)
			if modus == 0 {
				category = valueBin - step
				if value < 0 {
					category = valueBin + step
				}
			} else {
				category = valueBin - modus
			}
			groupKey = tk.ToString(category)
			dataCatCount[groupKey] = dataCatCount[groupKey] + 1
		}
	}
	csrData.Close()
	dataSeriesVal := setContribution("Met Tower", "avgwindspeed", dataCatCount, countPerWSCat)
	dataSeries = append(dataSeries, setSeries("Met Tower", colorList[0], 0, dataSeriesVal))

	return dataSeries
}

func GetScadaDataLahori(turbineName map[string]string, turbineNameSorted, turbineIDSorted []string, queryT []*dbox.Filter,
	tipe string, tStart, tEnd time.Time) ([]tk.M, []string) {
	var e error
	maxStep := maxNacelle
	step := stepNacelle
	fieldName := []string{"turbine", "winddirection", "naceldirection"}

	csrData, e := DB().Connection.NewQuery().
		Select("turbine", "winddirection", "naceldirection", "timestamp").
		From(new(ScadaData).TableName()).
		Where(dbox.And(queryT...)).
		Cursor(nil)
	defer csrData.Close()
	if e != nil {
		tk.Println("error", e.Error())
		return []tk.M{}, []string{}
	}

	sortedTimeStamp := []time.Time{}
	timeCount := tStart.UTC()
	for {
		timeCount = timeCount.Add(time.Duration(time.Minute) * 10).UTC()
		if timeCount.After(tEnd.UTC()) {
			break
		}
		sortedTimeStamp = append(sortedTimeStamp, timeCount)
	}

	_turbine := ""
	groupKey := ""
	countPerWSCat := 0.0
	dataCatCount := map[string]float64{}
	category := 0.0
	modus := 0.0
	dataSeries := []tk.M{}
	dataSeriesPerTurbine := map[string][]float64{}
	bin := step / 2
	_data := tk.M{}
	lastNacelTime := map[string]time.Time{}
	lastNacelValue := map[time.Time]float64{}
	currTimeStamp := time.Time{}
	hasNacel := true
	datas := []tk.M{}
	e = csrData.Fetch(&datas, 0, false)
	if e != nil {
		tk.Println("error", e.Error())
		return []tk.M{}, []string{}
	}
	sortedData := []tk.M{}
	dataPerTurbine := map[string][]tk.M{}
	dataPerTimeStamp := map[time.Time][]tk.M{}
	sortedDataByTimeStamp := []tk.M{}

	for _, _data = range datas { /* grouping data per timestamp */
		waktu := _data.Get("timestamp", time.Time{}).(time.Time).UTC()
		dataPerTimeStamp[waktu] = append(dataPerTimeStamp[waktu], _data)
	}
	for _, _waktu := range sortedTimeStamp { /* masukkan data yang sudah urut timestamp ke variable */
		data, hasData := dataPerTimeStamp[_waktu]
		if hasData {
			sortedDataByTimeStamp = append(sortedDataByTimeStamp, data...)
		}
	}
	for _, _data = range sortedDataByTimeStamp { /* grouping sortedDataByTimeStamp per turbine */
		dataPerTurbine[_data.GetString("turbine")] = append(dataPerTurbine[_data.GetString("turbine")], _data)
	}
	for _, _turbine = range turbineIDSorted {
		sortedData = dataPerTurbine[_turbine]
		for _, _data = range sortedData {
			_turbine = turbineName[_data.GetString("turbine")]
			currTimeStamp = _data.Get("timestamp", time.Time{}).(time.Time).UTC()
			if _data.Has(fieldName[2]) {
				lastNacelTime[_turbine] = currTimeStamp                            /* store latest timestamp per turbine for nacelle */
				lastNacelValue[currTimeStamp] = _data.GetFloat64("naceldirection") /* store latest value per timestamp */
				hasNacel = true
			} else {
				hasNacel = false
			}
			if _data.Has(fieldName[1]) {
				if !hasNacel {
					latestNacel, hasLatest := lastNacelTime[_turbine]
					if hasLatest {
						if currTimeStamp.Sub(latestNacel).Hours() < 1 {
							_data.Set(fieldName[2], lastNacelValue[lastNacelTime[_turbine]])
						} else {
							continue
						}
					} else {
						continue
					}
				}
				countPerWSCat++
				value := math.Abs(_data.GetFloat64(fieldName[1]) - _data.GetFloat64(fieldName[2]))
				if value > maxStep { /* jika datanya melebihi max step, maka ubah menjadi max step*/
					value = maxStep
				}

				valueBin := value + bin

				if value < 0 {
					valueBin = value - bin
				}

				modus = math.Mod(valueBin, step)
				if modus == 0 {
					category = valueBin - step
					if value < 0 {
						category = valueBin + step
					}
				} else {
					category = valueBin - modus
				}

				groupKey = tk.ToString(category)
				dataCatCount[groupKey] = dataCatCount[groupKey] + 1
			}
		}
		dataSeriesPerTurbine[_turbine] = setContribution(_turbine, tipe, dataCatCount, countPerWSCat)
		dataCatCount = map[string]float64{}
		countPerWSCat = 0.0
	}

	turbineAvail := []string{}
	index := 0
	for _, _turbinename := range turbineNameSorted {
		if dataSeriesVal, hasKey := dataSeriesPerTurbine[_turbinename]; hasKey {
			dataSeries = append(dataSeries, setSeries(_turbinename, colorList[index], index, dataSeriesVal))
			turbineAvail = append(turbineAvail, _turbinename)
			index++
		}
	}

	return dataSeries, turbineAvail
}

func GetScadaData(turbineName map[string]string, turbineNameSorted, turbineIDSorted []string, queryT []*dbox.Filter,
	tipe string, withMet bool) ([]tk.M, []string) {
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

	csrData, e := DB().Connection.NewQuery().
		Select("turbine", fieldName).
		From(new(ScadaData).TableName()).
		Where(dbox.And(queryT...)).
		// Order("turbine").
		Cursor(nil)
	defer csrData.Close()
	if e != nil {
		tk.Println("error on fetch", e.Error())
		return []tk.M{}, []string{}
	}

	_turbine := ""
	groupKey := ""
	countPerWSCat := 0.0
	dataCatCount := map[string]float64{}
	category := 0.0
	modus := 0.0
	dataSeries := []tk.M{}
	dataSeriesPerTurbine := map[string][]float64{}
	bin := step / 2
	_data := tk.M{}
	datas := []tk.M{}
	e = csrData.Fetch(&datas, 0, false)
	if e != nil {
		tk.Println("error on fetch", e.Error())
		return []tk.M{}, []string{}
	}
	dataPerTurbine := map[string][]tk.M{}
	for _, _data = range datas { /* grouping data per turbine */
		dataPerTurbine[_data.GetString("turbine")] = append(dataPerTurbine[_data.GetString("turbine")], _data)
	}
	for _, _turbine = range turbineIDSorted {
		datas = dataPerTurbine[_turbine]
		for _, _data = range datas {
			if _data.Has(fieldName) {
				countPerWSCat++
				value := _data.GetFloat64(fieldName)
				if value > maxStep { /* jika datanya melebihi max step, maka ubah menjadi max step*/
					value = maxStep
				}

				valueBin := value + bin

				if value < 0 {
					valueBin = value - bin
				}

				modus = math.Mod(valueBin, step)
				if modus == 0 {
					category = valueBin - step
					if value < 0 {
						category = valueBin + step
					}
				} else {
					category = valueBin - modus
				}

				groupKey = tk.ToString(category)
				dataCatCount[groupKey] = dataCatCount[groupKey] + 1

				// uncomment following codes to debug if there is wrong plotted value
				/*if category == 6 {
					tk.Printf("%f, ", value)
					if value <= 5.75 && value > 6.25 {
						tk.Printf("salah plot => %f", value)
					}
				} else if category == -15 {
					if value >= -7.5 && value < -22.5 {
						tk.Printf("salah plot => %f", value)
					}
				} else if category == 15 {
					if value <= 7.5 && value > 22.5 {
						tk.Printf("salah plot => %f", value)
					}
				}*/
			}
		}
		_turbine = turbineName[_data.GetString("turbine")]

		dataSeriesPerTurbine[_turbine] = setContribution(_turbine, tipe, dataCatCount, countPerWSCat)
		dataCatCount = map[string]float64{}
		countPerWSCat = 0.0
	}

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
			dataScada, turbineAvailTemp = GetScadaData(turbineName, turbineNameSorted, turbine, queryT, p.BreakDown, true)
			wg.Done()
		}()
		wg.Wait()
		dataSeries = append(dataMetTower, dataScada...)
		turbineAvail = append(turbineAvail, turbineAvailTemp...)
	} else {
		if p.Project == "Lahori" && p.BreakDown == "nacelledeviation" {
			dataSeries, turbineAvail = GetScadaDataLahori(turbineName, turbineNameSorted, turbine, queryT, p.BreakDown, tStart, tEnd)
		} else {
			dataSeries, turbineAvail = GetScadaData(turbineName, turbineNameSorted, turbine, queryT, p.BreakDown, false)
		}
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
