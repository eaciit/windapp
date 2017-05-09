package controller

import (
	. "eaciit/wfdemo-git/library/core"
	. "eaciit/wfdemo-git/library/models"
	"eaciit/wfdemo-git/web/helper"
	"sort"
	"strings"
	"time"

	c "github.com/eaciit/crowd"
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
var colorWindrose = []string{"#ff880e", "#21c4af", "#ff7663", "#ffb74f", "#a2df53", "#1c9ec4", "#ff63a5", "#f44336", "#69d2e7", "#8877A9", "#9A12B3", "#26C281", "#E7505A", "#C49F47", "#ff5597", "#c3260c", "#d4735e", "#ff2ad7", "#34ac8b", "#11b2eb", "#004c79", "#ff0037", "#507ca3", "#ff6565", "#ffd664", "#72aaff", "#795548"}
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

	/*lastDateData, _ := time.Parse("2006-01-02 15:04", "2016-09-30 23:59")
	k.SetSession("custom_lastdate", lastDateData.UTC())*/
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
	getFullWSCategory()

	coId := 0

	query := []toolkit.M{}
	pipes := []toolkit.M{}
	query = append(query, toolkit.M{"_id": toolkit.M{"$ne": nil}})
	query = append(query, toolkit.M{"available": 1})
	query = append(query, toolkit.M{"dateinfo.dateid": toolkit.M{"$gte": tStart}})
	if p.IsMonitoring {
		query = append(query, toolkit.M{"dateinfo.dateid": toolkit.M{"$lt": tEnd}})
	} else {
		query = append(query, toolkit.M{"dateinfo.dateid": toolkit.M{"$lte": tEnd}})
	}
	if p.Project != "" {
		anProject := strings.Split(p.Project, "(")
		query = append(query, toolkit.M{"projectname": strings.TrimRight(anProject[0], " ")})
	}

	turbine := []string{}
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
	if !p.IsMonitoring {
		turbine = append([]string{"MetTower"}, turbine...)
	}

	data := []MiniScada{}
	_data := MiniScada{}
	calibratedWindDir := 0.0

	for _, turbineVal := range turbine {
		coId++
		groupdata := toolkit.M{}
		groupdata.Set("Index", coId)
		groupdata.Set("Name", turbineVal)

		if turbineVal != "MetTower" {
			pipes = []toolkit.M{}
			data = []MiniScada{}
			queryT := query
			queryT = append(queryT, toolkit.M{"turbine": turbineVal})
			pipes = append(pipes, toolkit.M{"$match": toolkit.M{"$and": queryT}})
			pipes = append(pipes, toolkit.M{"$project": toolkit.M{"naceldirection": 1, "avgwindspeed": 1, "winddirection": 1}})
			csr, _ := DB().Connection.NewQuery().From(new(ScadaData).TableName()).
				Command("pipe", pipes).Cursor(nil)

			for {
				e = csr.Fetch(&_data, 1, false)
				if e != nil {
					break
				}
				data = append(data, _data)
			}
			csr.Close()

			if toolkit.SliceLen(data) > 0 {
				totalDuration := float64(len(data)) /* Tot data * 2 for get total minutes*/
				datas := c.From(&data).Apply(func(x interface{}) interface{} {
					dt := x.(MiniScada)
					var di DataItems
					calibratedWindDir = dt.WindDirection //+ 300
					dirNo, dirDesc := getDirection(dt.NacelDirection+calibratedWindDir, section)
					wsNo, wsDesc := getWsCategory(dt.AvgWindSpeed)

					di.DirectionNo = dirNo
					di.DirectionDesc = dirDesc
					di.WsCategoryNo = wsNo
					di.WsCategoryDesc = wsDesc
					di.Frequency = 1

					return di
				}).Exec().Group(func(x interface{}) interface{} {
					dt := x.(DataItems)

					var dig DataItemsGroup
					dig.DirectionNo = dt.DirectionNo
					dig.DirectionDesc = dt.DirectionDesc
					dig.WsCategoryNo = dt.WsCategoryNo
					dig.WsCategoryDesc = dt.WsCategoryDesc

					return dig
				}, nil).Exec()

				dts := datas.Apply(func(x interface{}) interface{} {
					kv := x.(c.KV)
					vv := kv.Key.(DataItemsGroup)
					vs := kv.Value.([]DataItems)

					sumFreq := c.From(&vs).Sum(func(x interface{}) interface{} {
						dt := x.(DataItems)
						return dt.Frequency
					}).Exec().Result.Sum

					var di DataItemsResult

					di.DirectionNo = vv.DirectionNo
					di.DirectionDesc = vv.DirectionDesc
					di.WsCategoryNo = vv.WsCategoryNo
					di.WsCategoryDesc = vv.WsCategoryDesc
					di.Hours = toolkit.Div(sumFreq, 6.0)
					di.Contribution = toolkit.RoundingAuto64(toolkit.Div(sumFreq, totalDuration)*100.0, 2)

					key := turbineVal + "_" + toolkit.ToString(di.DirectionNo)

					if !tkMaxVal.Has(key) {
						tkMaxVal.Set(key, di.Contribution)
					} else {
						tkMaxVal.Set(key, tkMaxVal.GetFloat64(key)+di.Contribution)
					}

					di.Frequency = int(sumFreq)

					return di
				}).Exec()

				results := dts.Result.Data().([]DataItemsResult)
				wsCategoryList := []string{}
				for _, dataRes := range results {
					wsCategoryList = append(wsCategoryList, toolkit.ToString(dataRes.DirectionNo)+
						"_"+toolkit.ToString(dataRes.WsCategoryNo)+"_"+dataRes.WsCategoryDesc)
				}
				splitCatList := []string{}
				for _, wsCat := range fullWSCatList {
					if !toolkit.HasMember(wsCategoryList, wsCat) {
						splitCatList = strings.Split(wsCat, "_")
						emptyRes := DataItemsResult{}
						emptyRes.DirectionNo = toolkit.ToInt(splitCatList[0], toolkit.RoundingAuto)
						divider := section

						emptyRes.DirectionDesc = (360 / divider) * emptyRes.DirectionNo
						emptyRes.WsCategoryNo = toolkit.ToInt(splitCatList[1], toolkit.RoundingAuto)
						emptyRes.WsCategoryDesc = splitCatList[2]
						results = append(results, emptyRes)
					}
				}
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
			pipes = []toolkit.M{}
			queryT := []toolkit.M{}
			dataMetTower := []MiniMetTower{}
			_dataMetTower := MiniMetTower{}
			queryT = append(queryT, toolkit.M{"_id": toolkit.M{"$ne": nil}})
			queryT = append(queryT, toolkit.M{"dateinfo.dateid": toolkit.M{"$gte": tStart}})
			queryT = append(queryT, toolkit.M{"dateinfo.dateid": toolkit.M{"$lte": tEnd}})

			pipes = append(pipes, toolkit.M{"$match": toolkit.M{"$and": queryT}})
			pipes = append(pipes, toolkit.M{"$project": toolkit.M{"vhubws90mavg": 1, "dhubwd88mavg": 1, "timestamp": 1}})
			csrMet, _ := DB().Connection.NewQuery().From(new(MetTower).TableName()).
				Command("pipe", pipes).Cursor(nil)

			for {
				e = csrMet.Fetch(&_dataMetTower, 1, false)
				if e != nil {
					break
				}
				dataMetTower = append(dataMetTower, _dataMetTower)
			}
			csrMet.Close()

			if toolkit.SliceLen(dataMetTower) > 0 {
				totalDuration := float64(len(dataMetTower)) // * 10.0 /* Tot data * 2 for get total minutes*/
				datas := c.From(&dataMetTower).Apply(func(x interface{}) interface{} {
					dt := x.(MiniMetTower)
					var di DataItems

					if dt.TimeStamp.UTC().Before(calibrateTime.UTC()) {
						calibratedWindDir = dt.DHubWD88mAvg + 300
					} else {
						calibratedWindDir = dt.DHubWD88mAvg
					}

					dirNo, dirDesc := getDirection(calibratedWindDir, section)
					wsNo, wsDesc := getWsCategory(dt.VHubWS90mAvg)

					di.DirectionNo = dirNo
					di.DirectionDesc = dirDesc
					di.WsCategoryNo = wsNo
					di.WsCategoryDesc = wsDesc
					di.Frequency = 1

					return di
				}).Exec().Group(func(x interface{}) interface{} {
					dt := x.(DataItems)

					var dig DataItemsGroup
					dig.DirectionNo = dt.DirectionNo
					dig.DirectionDesc = dt.DirectionDesc
					dig.WsCategoryNo = dt.WsCategoryNo
					dig.WsCategoryDesc = dt.WsCategoryDesc

					return dig
				}, nil).Exec()

				dts := datas.Apply(func(x interface{}) interface{} {
					kv := x.(c.KV)
					vv := kv.Key.(DataItemsGroup)
					vs := kv.Value.([]DataItems)

					sumFreq := c.From(&vs).Sum(func(x interface{}) interface{} {
						dt := x.(DataItems)
						return dt.Frequency
					}).Exec().Result.Sum

					var di DataItemsResult

					di.DirectionNo = vv.DirectionNo
					di.DirectionDesc = vv.DirectionDesc
					di.WsCategoryNo = vv.WsCategoryNo
					di.WsCategoryDesc = vv.WsCategoryDesc
					di.Hours = toolkit.Div(sumFreq, 6.0)
					di.Contribution = toolkit.RoundingAuto64(toolkit.Div(sumFreq, totalDuration)*100.0, 2)

					key := turbineVal + "_" + toolkit.ToString(di.DirectionNo)
					if !tkMaxVal.Has(key) {
						tkMaxVal.Set(key, di.Contribution)
					} else {
						tkMaxVal.Set(key, tkMaxVal.GetFloat64(key)+di.Contribution)
					}

					di.Frequency = int(sumFreq)

					return di
				}).Exec()

				results := dts.Result.Data().([]DataItemsResult)
				wsCategoryList := []string{}
				for _, dataRes := range results {
					wsCategoryList = append(wsCategoryList, toolkit.ToString(dataRes.DirectionNo)+
						"_"+toolkit.ToString(dataRes.WsCategoryNo)+"_"+dataRes.WsCategoryDesc)
				}
				splitCatList := []string{}
				for _, wsCat := range fullWSCatList {
					if !toolkit.HasMember(wsCategoryList, wsCat) {
						splitCatList = strings.Split(wsCat, "_")
						emptyRes := DataItemsResult{}
						emptyRes.DirectionNo = toolkit.ToInt(splitCatList[0], toolkit.RoundingAuto)
						divider := section

						emptyRes.DirectionDesc = (360 / divider) * emptyRes.DirectionNo
						emptyRes.WsCategoryNo = toolkit.ToInt(splitCatList[1], toolkit.RoundingAuto)
						emptyRes.WsCategoryDesc = splitCatList[2]
						results = append(results, emptyRes)
					}
				}
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

	type DataItemsComp struct {
		DirectionNo   int
		DirectionDesc int
		Frequency     int
	}

	type DataItemsResultComp struct {
		DirectionNo   int
		DirectionDesc int
		Hours         float64
		Contribution  float64
	}

	type DataItemsGroupComp struct {
		DirectionNo   int
		DirectionDesc int
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

	coId := 0

	query := []toolkit.M{}
	pipes := []toolkit.M{}
	query = append(query, toolkit.M{"_id": toolkit.M{"$ne": nil}})
	query = append(query, toolkit.M{"dateinfo.dateid": toolkit.M{"$gte": tStart}})
	query = append(query, toolkit.M{"dateinfo.dateid": toolkit.M{"$lte": tEnd}})
	query = append(query, toolkit.M{"available": 1})

	if p.Project != "" {
		anProject := strings.Split(p.Project, "(")
		query = append(query, toolkit.M{"projectname": strings.TrimRight(anProject[0], " ")})
	}

	turbine := []string{}
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
	turbine = append([]string{"Met Tower"}, turbine...)

	data := []MiniScada{}
	_data := MiniScada{}
	selArr := 0
	calibratedWindDir := 0.0

	for _, turbineVal := range turbine {
		coId++
		turbineData := toolkit.M{}
		turbineData.Set("name", turbineVal)
		turbineData.Set("type", "polarLine")
		turbineData.Set("color", colorWindrose[selArr])
		turbineData.Set("idxseries", selArr)
		turbineData.Set("xField", "DirectionDesc")
		turbineData.Set("yField", "Contribution")
		selArr++
		// dataDir := []float64{}

		if turbineVal != "Met Tower" {
			pipes = []toolkit.M{}
			data = []MiniScada{}
			queryT := query
			queryT = append(queryT, toolkit.M{"turbine": turbineVal})
			pipes = append(pipes, toolkit.M{"$match": toolkit.M{"$and": queryT}})
			pipes = append(pipes, toolkit.M{"$project": toolkit.M{"naceldirection": 1, "avgwindspeed": 1, "winddirection": 1}})
			csr, _ := DB().Connection.NewQuery().From(new(ScadaData).TableName()).
				Command("pipe", pipes).Cursor(nil)

			for {
				e = csr.Fetch(&_data, 1, false)
				if e != nil {
					break
				}
				data = append(data, _data)
			}
			csr.Close()

			if toolkit.SliceLen(data) > 0 {
				totalDuration := float64(len(data)) /* Tot data * 2 for get total minutes*/
				datas := c.From(&data).Apply(func(x interface{}) interface{} {
					dt := x.(MiniScada)
					var di DataItemsComp

					calibratedWindDir = dt.WindDirection //+ 300
					dirNo, dirDesc := getDirection(dt.NacelDirection+calibratedWindDir, section)

					di.DirectionNo = dirNo
					di.DirectionDesc = dirDesc
					di.Frequency = 1

					return di
				}).Exec().Group(func(x interface{}) interface{} {
					dt := x.(DataItemsComp)

					var dig DataItemsGroupComp
					dig.DirectionNo = dt.DirectionNo
					dig.DirectionDesc = dt.DirectionDesc

					return dig
				}, nil).Exec()

				dts := datas.Apply(func(x interface{}) interface{} {
					kv := x.(c.KV)
					vv := kv.Key.(DataItemsGroupComp)
					vs := kv.Value.([]DataItemsComp)

					sumFreq := c.From(&vs).Sum(func(x interface{}) interface{} {
						dt := x.(DataItemsComp)
						return dt.Frequency
					}).Exec().Result.Sum

					var di DataItemsResultComp

					di.DirectionNo = vv.DirectionNo
					di.DirectionDesc = vv.DirectionDesc
					di.Hours = toolkit.Div(sumFreq, 6.0)
					di.Contribution = toolkit.RoundingAuto64(toolkit.Div(sumFreq, totalDuration)*100.0, 2)

					key := turbineVal + "_" + toolkit.ToString(di.DirectionNo)

					if !tkMaxVal.Has(key) {
						tkMaxVal.Set(key, di.Contribution)
					} else {
						tkMaxVal.Set(key, tkMaxVal.GetFloat64(key)+di.Contribution)
					}

					return di
				}).Exec()

				result := dts.Result.Data().([]DataItemsResultComp)
				dirCatList := []string{}
				dirContribute := map[string]float64{}
				dirHours := map[string]float64{}
				for _, dataRes := range result {
					dirCatList = append(dirCatList, toolkit.ToString(dataRes.DirectionDesc))
					dirContribute[toolkit.ToString(dataRes.DirectionDesc)] = dataRes.Contribution
					dirHours[toolkit.ToString(dataRes.DirectionDesc)] = dataRes.Hours
				}
				results := []DataItemsResultComp{}
				firstData := DataItemsResultComp{}
				for idx, dirCat := range categories {
					dataRes := DataItemsResultComp{}
					if !toolkit.HasMember(dirCatList, dirCat) { /*if empty*/
						dataRes.DirectionDesc = toolkit.ToInt(dirCat, toolkit.RoundingAuto)
						dataRes.DirectionNo = dataRes.DirectionDesc / divider
						dataRes.Contribution = 0.0
						dataRes.Hours = 0.0
						results = append(results, dataRes)
						// dataDir = append(dataDir, 0.0)
					} else {
						dataRes.DirectionDesc = toolkit.ToInt(dirCat, toolkit.RoundingAuto)
						dataRes.DirectionNo = dataRes.DirectionDesc / divider
						dataRes.Contribution = dirContribute[dirCat]
						dataRes.Hours = dirHours[dirCat]
						results = append(results, dataRes)
						// dataDir = append(dataDir, dirContribute[dirCat])
					}
					if idx == 0 {
						firstData = dataRes
					}
				}
				// turbineData.Set("data", dataDir)
				results = append(results, firstData)
				turbineData.Set("data", results)
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
					// dataDir = append(dataDir, 0.0)
				}
				// turbineData.Set("data", dataDir)
				results = append(results, firstData)
				turbineData.Set("data", results)
				WindRoseResult = append(WindRoseResult, turbineData)
			}
		} else {
			pipes = []toolkit.M{}
			queryT := []toolkit.M{}
			dataMetTower := []MiniMetTower{}
			_dataMetTower := MiniMetTower{}
			queryT = append(queryT, toolkit.M{"_id": toolkit.M{"$ne": nil}})
			queryT = append(queryT, toolkit.M{"dateinfo.dateid": toolkit.M{"$gte": tStart}})
			queryT = append(queryT, toolkit.M{"dateinfo.dateid": toolkit.M{"$lte": tEnd}})

			pipes = append(pipes, toolkit.M{"$match": toolkit.M{"$and": queryT}})
			pipes = append(pipes, toolkit.M{"$project": toolkit.M{"vhubws90mavg": 1, "dhubwd88mavg": 1}})
			csrMet, _ := DB().Connection.NewQuery().From(new(MetTower).TableName()).
				Command("pipe", pipes).Cursor(nil)

			for {
				e = csrMet.Fetch(&_dataMetTower, 1, false)
				if e != nil {
					break
				}
				dataMetTower = append(dataMetTower, _dataMetTower)
			}
			csrMet.Close()

			if toolkit.SliceLen(dataMetTower) > 0 {
				totalDuration := float64(len(dataMetTower)) // * 10.0 /* Tot data * 2 for get total minutes*/
				datas := c.From(&dataMetTower).Apply(func(x interface{}) interface{} {
					dt := x.(MiniMetTower)
					var di DataItemsComp
					calibratedWindDir = dt.DHubWD88mAvg + 300
					dirNo, dirDesc := getDirection(calibratedWindDir, section)

					di.DirectionNo = dirNo
					di.DirectionDesc = dirDesc
					di.Frequency = 1

					return di
				}).Exec().Group(func(x interface{}) interface{} {
					dt := x.(DataItemsComp)

					var dig DataItemsGroupComp
					dig.DirectionNo = dt.DirectionNo
					dig.DirectionDesc = dt.DirectionDesc

					return dig
				}, nil).Exec()

				dts := datas.Apply(func(x interface{}) interface{} {
					kv := x.(c.KV)
					vv := kv.Key.(DataItemsGroupComp)
					vs := kv.Value.([]DataItemsComp)

					sumFreq := c.From(&vs).Sum(func(x interface{}) interface{} {
						dt := x.(DataItemsComp)
						return dt.Frequency
					}).Exec().Result.Sum

					var di DataItemsResultComp

					di.DirectionNo = vv.DirectionNo
					di.DirectionDesc = vv.DirectionDesc
					di.Hours = toolkit.Div(sumFreq, 6.0)
					di.Contribution = toolkit.RoundingAuto64(toolkit.Div(sumFreq, totalDuration)*100.0, 2)

					key := turbineVal + "_" + toolkit.ToString(di.DirectionNo)
					if !tkMaxVal.Has(key) {
						tkMaxVal.Set(key, di.Contribution)
					} else {
						tkMaxVal.Set(key, tkMaxVal.GetFloat64(key)+di.Contribution)
					}

					return di
				}).Exec()

				result := dts.Result.Data().([]DataItemsResultComp)
				dirCatList := []string{}
				dirContribute := map[string]float64{}
				dirHours := map[string]float64{}
				for _, dataRes := range result {
					dirCatList = append(dirCatList, toolkit.ToString(dataRes.DirectionDesc))
					dirContribute[toolkit.ToString(dataRes.DirectionDesc)] = dataRes.Contribution
					dirHours[toolkit.ToString(dataRes.DirectionDesc)] = dataRes.Hours
				}
				results := []DataItemsResultComp{}
				firstData := DataItemsResultComp{}
				for idx, dirCat := range categories {
					dataRes := DataItemsResultComp{}
					if !toolkit.HasMember(dirCatList, dirCat) {
						dataRes.DirectionDesc = toolkit.ToInt(dirCat, toolkit.RoundingAuto)
						dataRes.DirectionNo = dataRes.DirectionDesc / divider
						dataRes.Contribution = 0.0
						dataRes.Hours = 0.0
						results = append(results, dataRes)
						// dataDir = append(dataDir, 0.0)
					} else {
						dataRes.DirectionDesc = toolkit.ToInt(dirCat, toolkit.RoundingAuto)
						dataRes.DirectionNo = dataRes.DirectionDesc / divider
						dataRes.Contribution = dirContribute[dirCat]
						dataRes.Hours = dirHours[dirCat]
						results = append(results, dataRes)
						// dataDir = append(dataDir, dirContribute[dirCat])
					}
					if idx == 0 {
						firstData = dataRes
					}
				}
				// turbineData.Set("data", dataDir)
				results = append(results, firstData)
				turbineData.Set("data", results)
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
					// dataDir = append(dataDir, 0.0)
				}
				// turbineData.Set("data", dataDir)
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
