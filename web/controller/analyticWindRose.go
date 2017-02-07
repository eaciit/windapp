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
	}

	type MiniScada struct {
		NacelDirection float64
		AvgWindSpeed   float64
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

	degree = toolkit.ToInt(p.BreakDown, toolkit.RoundingAuto)
	section = degree
	getFullWSCategory()

	coId := 0

	query := []toolkit.M{}
	pipes := []toolkit.M{}
	query = append(query, toolkit.M{"_id": toolkit.M{"$ne": nil}})
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
			pipes = append(pipes, toolkit.M{"$project": toolkit.M{"naceldirection": 1, "avgwindspeed": 1}})
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

					dirNo, dirDesc := getDirection(dt.NacelDirection, section)
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
					var di DataItems

					dirNo, dirDesc := getDirection(dt.DHubWD88mAvg, section)
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
