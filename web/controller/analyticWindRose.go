package controller

import (
	. "eaciit/wfdemo-git/library/core"
	. "eaciit/wfdemo-git/library/models"
	"eaciit/wfdemo-git/web/helper"
	"sort"
	"strings"
	"time"

	"github.com/eaciit/dbox"
	"github.com/eaciit/knot/knot.v1"
	"github.com/eaciit/toolkit"
)

type AnalyticWindRoseController struct {
	App
}

func CreateAnalyticWindRoseController() *AnalyticWindRoseController {
	var controller = new(AnalyticWindRoseController)
	return controller
}

/*color palette below already remove some colors that not sharp enough, beware out of index*/
// var colorWindrose = []string{"#87c5da","#cc2a35", "#d66b76", "#5d1b62", "#f1c175","#95204c","#8f4bc5","#7d287d","#00818e","#c8c8c8","#546698","#66c99a","#f3d752","#20adb8","#333d6b","#d077b1","#aab664","#01a278","#c1d41a","#807063","#ff5975","#01a3d4","#ca9d08","#026e51","#4c653f","#007ca7"}
var colorWindrose = []string{"#ff9933", "#21c4af", "#ff7663", "#ffb74f", "#a2df53", "#1c9ec4", "#ff63a5", "#f44336", "#69d2e7", "#8877A9", "#9A12B3", "#26C281", "#E7505A", "#C49F47", "#ff5597", "#c3260c", "#d4735e", "#ff2ad7", "#34ac8b", "#11b2eb", "#004c79", "#ff0037", "#507ca3", "#ff6565", "#ffd664", "#72aaff", "#795548", "#383271", "#6a4795", "#bec554", "#ab5919", "#f5b1e1", "#7b3416", "#002fef", "#8d731b", "#1f8805", "#ff9900", "#9C27B0", "#6c7d8a", "#d73c1c", "#5be7a0", "#da02d4", "#afa56e", "#7e32cb", "#a2eaf7", "#9cb8f4", "#9E9E9E", "#065806", "#044082", "#18937d", "#2c787a", "#a57c0c", "#234341", "#1aae7a", "#7ac610", "#736f5f", "#4e741e", "#68349d", "#1df3b6", "#e02b09", "#d9cfab", "#6e4e52", "#f31880", "#7978ec", "#f5ace8", "#3db6ae", "#5e06b0", "#16d0b9", "#a25a5b", "#1e603a", "#4b0981", "#62975f", "#1c8f2f", "#b0c80c", "#642794", "#e2060d", "#2125f0"}
var calibrateTime, _ = time.Parse("01022006_150405", "12012016_000000")

// var colorWindrose = []string{
// 	"#B71C1C", "#F44336", "#D81B60", "#F06292", "#880E4F",
// 	"#4A148C", "#7B1FA2", "#9C27B0", "#BA68C8", "#1A237E",
// 	"#5C6BC0", "#1E88E5", "#0277BD", "#0097A7", "#26A69A",
// 	"#81C784", "#8BC34A", "#24752A", "#827717", "#004D40",
// 	"#C0CA33", "#FF6F00", "#D6C847", "#FFB300", "#BA8914",
// 	"#9999FF",
// }

/*var colorWindrose = []string{
	"#ff880e", "#21c4af", "#ff7663", "#ffb74f", "#a2df53",
	"#1c9ec4", "#ff63a5", "#f44336", "#D91E18", "#8877A9",
	"#9A12B3", "#26C281", "#E7505A", "#C49F47", "#ff5597",
	"#c3260c", "#d4735e", "#ff2ad7", "#34ac8b", "#11b2eb",
	"#f35838", "#ff0037", "#507ca3", "#ff6565", "#ffd664",
	"#72aaff", "#795548",
}*/

type DataItemsResultComp struct {
	DirectionNo   int
	DirectionDesc int
	Hours         float64
	Contribution  float64
}

func (m *AnalyticWindRoseController) GetWSData(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	p := new(helper.Payloads)
	e := k.GetPayload(&p)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	filter, _ := p.ParseFilter()
	fb := DB().Connection.Fb()
	fb.AddFilter(dbox.And(filter...))
	matches, e := fb.Build()
	if e != nil {
		helper.CreateResult(false, nil, e.Error())
	}

	groups := toolkit.M{}
	groupIds := toolkit.M{}
	group := []string{
		"directionno",
		"directiondesc",
		"wscategoryno",
		"wscategorydesc",
	}
	for _, val := range group {
		alias := val
		field := toolkit.Sprintf("$windroseitems.%s", val)
		groupIds[alias] = field
	}
	groups["_id"] = groupIds

	fields := []string{
		"hours",
		"contribute",
		"frequency",
	}

	for _, other := range fields {
		field := toolkit.Sprintf("$windroseitems.%s", other)
		op := ""
		if other == "contribute" {
			op = "$avg"
		} else {
			op = "$sum"
		}
		groups[other] = toolkit.M{op: field}
	}

	pipe := []toolkit.M{{"$unwind": "$windroseitems"}, {"$match": matches}, {"$group": groups}}

	csr, e := DB().Connection.NewQuery().
		From(new(WindRoseModel).TableName()).
		Command("pipe", pipe).
		Cursor(nil)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}
	defer csr.Close()

	resultItem := toolkit.Ms{}
	e = csr.Fetch(&resultItem, 0, false)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	results := toolkit.Ms{}
	for _, val := range resultItem {
		result, _ := toolkit.ToM(val["_id"])
		result.Set("hours", val["hours"])
		result.Set("contribute", val["contribute"])
		result.Set("frequency", val["frequency"])
		results = append(results, result)
	}

	data := struct {
		WindRose toolkit.Ms
	}{
		WindRose: results,
	}

	return helper.CreateResult(true, data, "success")
}

func (m *AnalyticWindRoseController) GetWSCategory(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	p := new(helper.Payloads)
	e := k.GetPayload(&p)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	filter, _ := p.ParseFilter()
	fb := DB().Connection.Fb()
	fb.AddFilter(dbox.And(filter...))
	matches, e := fb.Build()
	if e != nil {
		helper.CreateResult(false, nil, e.Error())
	}

	groupsWS := toolkit.M{}
	groupIdWS := toolkit.M{}
	group := []string{
		"wscategoryno",
		"wscategorydesc",
	}
	for _, val := range group {
		alias := val
		fieldWS := toolkit.Sprintf("$totalcontributes.%s", val)
		groupIdWS[alias] = fieldWS
	}
	groupsWS["_id"] = groupIdWS

	fields := []string{
		"hours",
		"contribute",
		"frequency",
	}

	for _, other := range fields {
		fieldWS := toolkit.Sprintf("$totalcontributes.%s", other)
		op := ""
		if other == "contribute" {
			op = "$avg"
		} else {
			op = "$sum"
		}
		groupsWS[other] = toolkit.M{op: fieldWS}
	}

	pipeWS := []toolkit.M{{"$unwind": "$totalcontributes"}, {"$match": matches}, {"$group": groupsWS}}

	csrWS, e := DB().Connection.NewQuery().
		From(new(WindRoseModel).TableName()).
		Command("pipe", pipeWS).
		Cursor(nil)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}
	defer csrWS.Close()

	WSData := toolkit.Ms{}
	e = csrWS.Fetch(&WSData, 0, false)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	resultsWS := toolkit.Ms{}
	for _, val := range WSData {
		result, _ := toolkit.ToM(val["_id"])
		result.Set("hours", val["hours"])
		result.Set("contribute", val["contribute"])
		result.Set("frequency", val["frequency"])
		resultsWS = append(resultsWS, result)
	}

	data := struct {
		WSCategory toolkit.Ms
		Total      int
	}{
		WSCategory: resultsWS,
		Total:      csrWS.Count(),
	}

	return helper.CreateResult(true, data, "success")
}

func setDataWS(dataDirNoDesc map[string]float64, tkMaxVal *toolkit.M,
	dataCount float64, divider int) (results []DataItemsResult) {
	maxValKey := ""
	splitDirNoDesc := []string{}
	diRes := DataItemsResult{}
	results = []DataItemsResult{}
	wsCategoryList := []string{}
	for dirNoDesc, sumFreq := range dataDirNoDesc {
		splitDirNoDesc = strings.Split(dirNoDesc, "_")
		diRes = DataItemsResult{}
		diRes.DirectionNo = toolkit.ToInt(splitDirNoDesc[1], toolkit.RoundingAuto)
		diRes.DirectionDesc = toolkit.ToInt(splitDirNoDesc[2], toolkit.RoundingAuto)
		diRes.WsCategoryNo = toolkit.ToInt(splitDirNoDesc[3], toolkit.RoundingAuto)
		diRes.WsCategoryDesc = splitDirNoDesc[4]
		diRes.Hours = toolkit.Div(sumFreq, 6.0)
		diRes.Contribution = toolkit.RoundingAuto64(toolkit.Div(sumFreq, dataCount)*100.0, 2)
		diRes.Frequency = toolkit.ToInt(sumFreq, toolkit.RoundingAuto)
		results = append(results, diRes)

		maxValKey = splitDirNoDesc[0] + "_" + toolkit.ToString(diRes.DirectionNo)
		if !tkMaxVal.Has(maxValKey) {
			tkMaxVal.Set(maxValKey, diRes.Contribution)
		} else {
			tkMaxVal.Set(maxValKey, tkMaxVal.GetFloat64(maxValKey)+diRes.Contribution)
		}

		wsCategoryList = append(wsCategoryList, splitDirNoDesc[1]+
			"_"+splitDirNoDesc[3]+"_"+splitDirNoDesc[4])
	}
	splitCatList := []string{}
	emptyRes := DataItemsResult{}
	for _, wsCat := range fullWSCatList {
		if !toolkit.HasMember(wsCategoryList, wsCat) {
			splitCatList = strings.Split(wsCat, "_")
			emptyRes = DataItemsResult{}
			emptyRes.DirectionNo = toolkit.ToInt(splitCatList[0], toolkit.RoundingAuto)

			emptyRes.DirectionDesc = divider * emptyRes.DirectionNo
			emptyRes.WsCategoryNo = toolkit.ToInt(splitCatList[1], toolkit.RoundingAuto)
			emptyRes.WsCategoryDesc = splitCatList[2]
			results = append(results, emptyRes)
		}
	}
	return
}

func (m *AnalyticWindRoseController) GetFlexiDataEachTurbine(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson
	WindRoseResult = []toolkit.M{}
	maxValue := 0.0
	tkMaxVal := toolkit.M{}

	type PayloadWindRose struct {
		Period       string
		Project      string
		IsMonitoring bool
		Turbine      []interface{}
		DateStart    time.Time
		DateEnd      time.Time
		BreakDown    string
	}
	type MiniMetTower struct {
		DHubWD88mAvg float64
		VHubWS90mAvg float64
		TimeStamp    time.Time
	}

	type MiniScada struct {
		NacelDirection float64
		AvgWindSpeed   float64
		WindDirection  float64
		Turbine        string
	}
	p := new(PayloadWindRose)
	e := k.GetPayload(&p)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	var tStart, tEnd time.Time
	if p.IsMonitoring {
		now := time.Now()
		last := now.AddDate(0, 0, -24)

		tStart, _ = time.Parse("20060102", last.Format("200601")+"01")
		tEnd, _ = time.Parse("20060102", now.Format("200601")+"01")
	} else {
		tStart, tEnd, e = helper.GetStartEndDate(k, p.Period, p.DateStart, p.DateEnd)
		if e != nil {
			return helper.CreateResult(false, nil, e.Error())
		}
	}

	section = toolkit.ToInt(p.BreakDown, toolkit.RoundingAuto)
	divider := 360 / section
	getFullWSCategory()

	coId := 0

	query := []toolkit.M{}
	pipes := []toolkit.M{}
	query = append(query, toolkit.M{"dateinfo.dateid": toolkit.M{"$gte": tStart}})
	if p.IsMonitoring {
		query = append(query, toolkit.M{"dateinfo.dateid": toolkit.M{"$lt": tEnd}})
	} else {
		query = append(query, toolkit.M{"dateinfo.dateid": toolkit.M{"$lte": tEnd}})
	}
	if p.Project != "" {
		query = append(query, toolkit.M{"projectname": p.Project})
	}
	query = append(query, toolkit.M{"available": 1})

	turbine := []string{}
	turbineInt := []interface{}{}
	if len(p.Turbine) == 0 {
		pipes = append(pipes, toolkit.M{"$match": toolkit.M{"$and": query}})
		pipes = append(pipes, toolkit.M{"$group": toolkit.M{"_id": "$turbine"}})

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
	for _, val := range turbine {
		turbineInt = append(turbineInt, val)
	}
	if !p.IsMonitoring && p.Project == "Tejuva" {
		turbine = append([]string{"MetTower"}, turbine...)
	}
	query = append(query, toolkit.M{"turbine": toolkit.M{"$in": turbine}})

	_data := toolkit.M{}
	dataDirNoDesc := map[string]float64{}
	dataPerTurbine := toolkit.M{}
	calibratedWindDir := 0.0
	_turbine := ""
	groupKey := ""
	lastTurbine := ""
	dataCount := 0.0

	turbineName, e := helper.GetTurbineNameList(p.Project)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	pipes = []toolkit.M{}
	pipes = append(pipes, toolkit.M{"$match": toolkit.M{"$and": query}})
	pipes = append(pipes, toolkit.M{"$project": toolkit.M{"naceldirection": 1, "winddirection": 1, "avgwindspeed": 1, "turbine": 1}})
	pipes = append(pipes, toolkit.M{"$sort": toolkit.M{"turbine": 1}})
	csrData, e := DB().Connection.NewQuery().From(new(ScadaData).TableName()).
		Command("pipe", pipes).Cursor(nil)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	for {
		_data = toolkit.M{}
		e = csrData.Fetch(&_data, 1, false)
		if e != nil {
			break
		}
		_turbine = _data.GetString("turbine")
		if lastTurbine != _turbine {
			dataPerTurbine.Set(lastTurbine, setDataWS(dataDirNoDesc, &tkMaxVal, dataCount, divider))
			dataDirNoDesc = map[string]float64{}
			lastTurbine = _turbine
			dataCount = 0.0
		}
		if _data.Has("naceldirection") {
			dataCount++ /* total frequency per turbine */
			// calibratedWindDir = _data.GetFloat64("winddirection")
			calibratedWindDir = 0
			dirNo, dirDesc := getDirection(_data.GetFloat64("naceldirection")+calibratedWindDir, section)
			wsNo, wsDesc := getWsCategory(_data.GetFloat64("avgwindspeed"))
			groupKey = _turbine + "_" + toolkit.ToString(dirNo) + "_" + toolkit.ToString(dirDesc) +
				"_" + toolkit.ToString(wsNo) + "_" + wsDesc
			dataDirNoDesc[groupKey] = dataDirNoDesc[groupKey] + 1
		}
	}
	if lastTurbine != "" {
		dataPerTurbine.Set(lastTurbine, setDataWS(dataDirNoDesc, &tkMaxVal, dataCount, divider))
	}
	csrData.Close()

	for _, turbineVal := range turbine {
		coId++
		groupdata := toolkit.M{}
		groupdata.Set("Index", coId)

		if turbineVal != "MetTower" {
			groupdata.Set("Name", turbineName[turbineVal])

			if dataPerTurbine.Has(turbineVal) {
				results := dataPerTurbine[turbineVal].([]DataItemsResult)
				groupdata.Set("Data", results)

				WindRoseResult = append(WindRoseResult, groupdata)
			} else {
				splitCatList := []string{}
				results := []DataItemsResult{}
				for _, wsCat := range fullWSCatList {
					splitCatList = strings.Split(wsCat, "_")
					emptyRes := DataItemsResult{}
					emptyRes.DirectionNo = toolkit.ToInt(splitCatList[0], toolkit.RoundingAuto)
					divider := section

					emptyRes.DirectionDesc = (360 / divider) * emptyRes.DirectionNo
					emptyRes.WsCategoryNo = toolkit.ToInt(splitCatList[1], toolkit.RoundingAuto)
					emptyRes.WsCategoryDesc = splitCatList[2]
					results = append(results, emptyRes)
				}
				groupdata.Set("Data", results)
				WindRoseResult = append(WindRoseResult, groupdata)
			}
		} else {
			groupdata.Set("Name", turbineVal)

			queryMet := query[0:3]
			_dataMetTower := toolkit.M{}

			pipes = []toolkit.M{}
			pipes = append(pipes, toolkit.M{"$match": toolkit.M{"$and": queryMet}})
			pipes = append(pipes, toolkit.M{"$project": toolkit.M{"dhubwd88mavg": 1, "vhubws90mavg": 1, "timestamp": 1}})
			csrMet, e := DB().Connection.NewQuery().From(new(MetTower).TableName()).
				Command("pipe", pipes).Cursor(nil)
			if e != nil {
				return helper.CreateResult(false, nil, e.Error())
			}
			dataCount = 0.0
			dataDirNoDesc = map[string]float64{}
			dataMetTower := toolkit.M{}

			for {
				_dataMetTower = toolkit.M{}
				e = csrMet.Fetch(&_dataMetTower, 1, false)
				if e != nil {
					break
				}
				if _dataMetTower.Has("dhubwd88mavg") {
					dataCount++
					if _dataMetTower.Get("timestamp", time.Time{}).(time.Time).UTC().Before(calibrateTime.UTC()) {
						calibratedWindDir = _dataMetTower.GetFloat64("dhubwd88mavg") + 300
					} else {
						calibratedWindDir = _dataMetTower.GetFloat64("dhubwd88mavg")
					}
					dirNo, dirDesc := getDirection(calibratedWindDir, section)
					wsNo, wsDesc := getWsCategory(_dataMetTower.GetFloat64("vhubws90mavg"))
					groupKey = turbineVal + "_" + toolkit.ToString(dirNo) + "_" + toolkit.ToString(dirDesc) +
						"_" + toolkit.ToString(wsNo) + "_" + wsDesc
					dataDirNoDesc[groupKey] = dataDirNoDesc[groupKey] + 1
				}
			}
			dataMetTower.Set(turbineVal, setDataWS(dataDirNoDesc, &tkMaxVal, dataCount, divider))
			csrMet.Close()

			if dataCount > 0 {
				results := dataMetTower[turbineVal].([]DataItemsResult)
				groupdata.Set("Data", results)

				WindRoseResult = append(WindRoseResult, groupdata)
			} else {
				splitCatList := []string{}
				results := []DataItemsResult{}
				for _, wsCat := range fullWSCatList {
					splitCatList = strings.Split(wsCat, "_")
					emptyRes := DataItemsResult{}
					emptyRes.DirectionNo = toolkit.ToInt(splitCatList[0], toolkit.RoundingAuto)
					divider := section

					emptyRes.DirectionDesc = (360 / divider) * emptyRes.DirectionNo
					emptyRes.WsCategoryNo = toolkit.ToInt(splitCatList[1], toolkit.RoundingAuto)
					emptyRes.WsCategoryDesc = splitCatList[2]
					results = append(results, emptyRes)
				}
				groupdata.Set("Data", results)
				WindRoseResult = append(WindRoseResult, groupdata)
			}
		}
	}

	for _, val := range tkMaxVal {
		if val.(float64) > maxValue {
			maxValue = val.(float64)
		}
	}

	switch {
	case maxValue >= 90 && maxValue <= 100:
		maxValue = 100
	case maxValue >= 80 && maxValue < 90:
		maxValue = 90
	case maxValue >= 70 && maxValue < 80:
		maxValue = 80
	case maxValue >= 60 && maxValue < 70:
		maxValue = 70
	case maxValue >= 50 && maxValue < 60:
		maxValue = 60
	case maxValue >= 40 && maxValue < 50:
		maxValue = 50
	case maxValue >= 30 && maxValue < 40:
		maxValue = 40
	case maxValue >= 20 && maxValue < 30:
		maxValue = 30
	case maxValue >= 10 && maxValue < 20:
		maxValue = 20
	case maxValue >= 0 && maxValue < 10:
		maxValue = 10
	}

	datas := struct {
		WindRose toolkit.Ms
		MaxValue float64
	}{
		WindRose: WindRoseResult,
		MaxValue: maxValue,
	}

	return helper.CreateResult(true, datas, "success")

}

func setDataPerTurbine(dataDirNoDesc map[string]float64, tkMaxVal *toolkit.M,
	dataCount float64, categories []string, divider int) (results []DataItemsResultComp) {
	maxValKey := ""
	splitDirNoDesc := []string{}
	diRes := DataItemsResultComp{}
	dirCatList := []string{}
	dirContribute := map[string]float64{}
	dirHours := map[string]float64{}
	results = []DataItemsResultComp{}
	for dirNoDesc, sumFreq := range dataDirNoDesc {
		splitDirNoDesc = strings.Split(dirNoDesc, "_")
		diRes = DataItemsResultComp{}
		diRes.DirectionNo = toolkit.ToInt(splitDirNoDesc[1], toolkit.RoundingAuto)
		diRes.DirectionDesc = toolkit.ToInt(splitDirNoDesc[2], toolkit.RoundingAuto)
		diRes.Hours = toolkit.Div(sumFreq, 6.0)
		diRes.Contribution = toolkit.RoundingAuto64(toolkit.Div(sumFreq, dataCount)*100.0, 2)

		maxValKey = splitDirNoDesc[0] + "_" + toolkit.ToString(diRes.DirectionNo)
		if !tkMaxVal.Has(maxValKey) {
			tkMaxVal.Set(maxValKey, diRes.Contribution)
		} else {
			tkMaxVal.Set(maxValKey, tkMaxVal.GetFloat64(maxValKey)+diRes.Contribution)
		}

		dirCatList = append(dirCatList, toolkit.ToString(diRes.DirectionDesc))
		dirContribute[toolkit.ToString(diRes.DirectionDesc)] = diRes.Contribution
		dirHours[toolkit.ToString(diRes.DirectionDesc)] = diRes.Hours
	}
	firstData := DataItemsResultComp{}
	for idx, dirCat := range categories {
		dataRes := DataItemsResultComp{}
		if !toolkit.HasMember(dirCatList, dirCat) { /*if empty*/
			dataRes.DirectionDesc = toolkit.ToInt(dirCat, toolkit.RoundingAuto)
			dataRes.DirectionNo = dataRes.DirectionDesc / divider
			dataRes.Contribution = 0.0
			dataRes.Hours = 0.0
			results = append(results, dataRes)
		} else {
			dataRes.DirectionDesc = toolkit.ToInt(dirCat, toolkit.RoundingAuto)
			dataRes.DirectionNo = dataRes.DirectionDesc / divider
			dataRes.Contribution = dirContribute[dirCat]
			dataRes.Hours = dirHours[dirCat]
			results = append(results, dataRes)
		}
		if idx == 0 {
			firstData = dataRes
		}
	}
	results = append(results, firstData)
	return
}

func (m *AnalyticWindRoseController) GetWindRoseData(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson
	WindRoseResult = []toolkit.M{}
	maxValue := 0.0
	tkMaxVal := toolkit.M{}

	type PayloadWindRose struct {
		Period    string
		Project   string
		Turbine   []interface{}
		DateStart time.Time
		DateEnd   time.Time
		BreakDown string
	}

	p := new(PayloadWindRose)
	e := k.GetPayload(&p)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	var tStart, tEnd time.Time
	tStart, tEnd, e = helper.GetStartEndDate(k, p.Period, p.DateStart, p.DateEnd)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	section = toolkit.ToInt(p.BreakDown, toolkit.RoundingAuto)
	categories := []string{}
	direction := 0
	divider := 360 / section
	for i := 0; i < section; i++ {
		categories = append(categories, toolkit.ToString(direction))
		direction += divider
	}

	query := []toolkit.M{}
	pipes := []toolkit.M{}
	query = append(query, toolkit.M{"_id": toolkit.M{"$ne": nil}})
	query = append(query, toolkit.M{"dateinfo.dateid": toolkit.M{"$gte": tStart}})
	query = append(query, toolkit.M{"dateinfo.dateid": toolkit.M{"$lte": tEnd}})
	query = append(query, toolkit.M{"available": 1})

	if p.Project != "" {
		query = append(query, toolkit.M{"projectname": p.Project})
	}

	turbine := []string{}
	turbineInt := []interface{}{}
	if len(p.Turbine) == 0 {
		pipes = append(pipes, toolkit.M{"$match": toolkit.M{"$and": query}})
		pipes = append(pipes, toolkit.M{"$group": toolkit.M{"_id": "$turbine"}})

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
	for _, val := range turbine {
		turbineInt = append(turbineInt, val)
	}
	queryT := []*dbox.Filter{}
	queryT = append(queryT, dbox.Gte("dateinfo.dateid", tStart))
	queryT = append(queryT, dbox.Lte("dateinfo.dateid", tEnd))
	queryT = append(queryT, dbox.Eq("available", 1))
	if p.Project != "" {
		queryT = append(queryT, dbox.Eq("projectname", p.Project))
	}
	queryT = append(queryT, dbox.In("turbine", turbineInt...))
	if p.Project == "Tejuva" {
		turbine = append([]string{"Met Tower"}, turbine...)
	}

	_data := toolkit.M{}
	dataDirNoDesc := map[string]float64{}
	selArr := 0
	calibratedWindDir := 0.0
	groupKey := ""
	dataCount := 0.0

	turbineName, e := helper.GetTurbineNameList(p.Project)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	query = append(query, toolkit.M{"turbine": ""})
	lastQueryIdx := len(query) - 1

	for _, turbineVal := range turbine {
		turbineData := toolkit.M{}
		namaTurbine := turbineVal
		if turbineVal != "Met Tower" {
			namaTurbine = turbineName[turbineVal]
		}
		turbineData.Set("name", namaTurbine)
		turbineData.Set("type", "polarLine")
		turbineData.Set("color", colorWindrose[selArr])
		turbineData.Set("idxseries", selArr)
		turbineData.Set("xField", "DirectionDesc")
		turbineData.Set("yField", "Contribution")
		selArr++

		if turbineVal != "Met Tower" {
			query[lastQueryIdx] = toolkit.M{"turbine": turbineVal}
			pipes = []toolkit.M{
				toolkit.M{"$match": toolkit.M{"$and": query}},
				toolkit.M{"$project": toolkit.M{"naceldirection": 1, "winddirection": 1, "turbine": 1}},
			}
			csrData, e := DB().Connection.NewQuery().From(new(ScadaData).TableName()).
				Command("pipe", pipes).Cursor(nil)
			if e != nil {
				return helper.CreateResult(false, nil, e.Error())
			}
			dataCount = 0
			dataDirNoDesc = map[string]float64{}
		loopFetchData:
			for {
				e = csrData.Fetch(&_data, 1, false)
				if e != nil {
					break loopFetchData
				}
				if _data.Has("naceldirection") {
					dataCount++
					// calibratedWindDir = _data.GetFloat64("winddirection")
					calibratedWindDir = 0
					dirNo, dirDesc := getDirection(_data.GetFloat64("naceldirection")+calibratedWindDir, section)
					groupKey = turbineVal + "_" + toolkit.ToString(dirNo) + "_" + toolkit.ToString(dirDesc)
					dataDirNoDesc[groupKey] = dataDirNoDesc[groupKey] + 1
				}
			}
			hasil := setDataPerTurbine(dataDirNoDesc, &tkMaxVal, dataCount, categories, divider)
			csrData.Close()

			if dataCount > 0 {
				turbineData.Set("data", hasil)
				WindRoseResult = append(WindRoseResult, turbineData)
			} else {
				results := []DataItemsResultComp{}
				firstData := DataItemsResultComp{}
				for idx, dirCat := range categories {
					emptyRes := DataItemsResultComp{}
					emptyRes.DirectionDesc = toolkit.ToInt(dirCat, toolkit.RoundingAuto)
					emptyRes.DirectionNo = emptyRes.DirectionDesc / divider
					emptyRes.Contribution = 0.0
					emptyRes.Hours = 0.0
					results = append(results, emptyRes)
					if idx == 0 {
						firstData = emptyRes
					}
				}
				results = append(results, firstData)
				turbineData.Set("data", results)
				WindRoseResult = append(WindRoseResult, turbineData)
			}
		} else {
			queryMet := []*dbox.Filter{}
			_dataMetTower := toolkit.M{}

			for _, filter := range queryT {
				if filter.Field != "available" && filter.Field != "turbine" {
					queryMet = append(queryMet, filter)
				}
			}

			csrMet, _ := DB().Connection.NewQuery().From(new(MetTower).TableName()).
				Select("dhubwd88mavg").
				Where(dbox.And(queryMet...)).Cursor(nil)
			dataCount = 0.0
			dataDirNoDesc = map[string]float64{}
			for {
				e = csrMet.Fetch(&_dataMetTower, 1, false)
				if e != nil {
					break
				}
				if _dataMetTower.Has("dhubwd88mavg") {
					dataCount++
					calibratedWindDir = _dataMetTower.GetFloat64("dhubwd88mavg") + 300
					dirNo, dirDesc := getDirection(calibratedWindDir, section)
					groupKey = turbineVal + "_" + toolkit.ToString(dirNo) + "_" + toolkit.ToString(dirDesc)
					dataDirNoDesc[groupKey] = dataDirNoDesc[groupKey] + 1
				}
			}
			hasil := setDataPerTurbine(dataDirNoDesc, &tkMaxVal, dataCount, categories, divider)
			csrMet.Close()

			if csrMet.Count() > 0 {
				turbineData.Set("data", hasil)
				WindRoseResult = append(WindRoseResult, turbineData)
			} else {
				results := []DataItemsResultComp{}
				firstData := DataItemsResultComp{}
				for idx, dirCat := range categories {
					emptyRes := DataItemsResultComp{}
					emptyRes.DirectionDesc = toolkit.ToInt(dirCat, toolkit.RoundingAuto)
					emptyRes.DirectionNo = emptyRes.DirectionDesc / divider
					emptyRes.Contribution = 0.0
					emptyRes.Hours = 0.0
					results = append(results, emptyRes)
					if idx == 0 {
						firstData = emptyRes
					}
				}
				results = append(results, firstData)
				turbineData.Set("data", results)
				WindRoseResult = append(WindRoseResult, turbineData)
			}
		}
	}

	for _, val := range tkMaxVal {
		if val.(float64) > maxValue {
			maxValue = val.(float64)
		}
	}

	switch {
	case maxValue >= 90 && maxValue <= 100:
		maxValue = 100
	case maxValue >= 80 && maxValue < 90:
		maxValue = 90
	case maxValue >= 70 && maxValue < 80:
		maxValue = 80
	case maxValue >= 60 && maxValue < 70:
		maxValue = 70
	case maxValue >= 50 && maxValue < 60:
		maxValue = 60
	case maxValue >= 40 && maxValue < 50:
		maxValue = 50
	case maxValue >= 35 && maxValue < 40:
		maxValue = 40
	case maxValue >= 30 && maxValue < 35:
		maxValue = 35
	case maxValue >= 25 && maxValue < 30:
		maxValue = 30
	case maxValue >= 20 && maxValue < 25:
		maxValue = 25
	case maxValue >= 15 && maxValue < 20:
		maxValue = 20
	case maxValue >= 10 && maxValue < 15:
		maxValue = 15
	case maxValue >= 5 && maxValue < 10:
		maxValue = 10
	case maxValue >= 0 && maxValue < 5:
		maxValue = 5
	}

	datas := struct {
		Data toolkit.Ms
		// Categories []string
		MaxValue float64
	}{
		Data: WindRoseResult,
		// Categories: categories,
		MaxValue: maxValue,
	}

	return helper.CreateResult(true, datas, "success")

}
