package controller

import (
	. "eaciit/wfdemo/library/core"
	. "eaciit/wfdemo/library/models"
	"eaciit/wfdemo/web/helper"
	c "github.com/eaciit/crowd"
	"github.com/eaciit/dbox"
	"github.com/eaciit/knot/knot.v1"
	"github.com/eaciit/toolkit"
	"sort"
	"strings"
	"time"
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

	lastDateData, _ := time.Parse("2006-01-02 15:04", "2016-07-31 23:59")
	k.SetSession("custom_lastdate", lastDateData.UTC())
	tStart, tEnd, e := helper.GetStartEndDate(k, p.Period, p.DateStart, p.DateEnd)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	bufferTurbine := []string{}
	for _, val := range p.Turbine {
		bufferTurbine = append(bufferTurbine, val.(string))
	}
	degree = toolkit.ToInt(p.BreakDown, toolkit.RoundingAuto)
	section = degree
	getFullWSCategory()

	coId := 0
	sort.Strings(bufferTurbine)
	turbine := []string{"MetTower"}
	turbine = append(turbine, bufferTurbine...)
	var filter []*dbox.Filter

	filter = append(filter, dbox.Gte("dateinfo.dateid", tStart))
	filter = append(filter, dbox.Lte("dateinfo.dateid", tEnd))
	if len(p.Turbine) > 0 {
		filter = append(filter, dbox.In("turbine", p.Turbine...))
	}
	filter = append(filter, dbox.Ne("_id", ""))

	scadas := []MiniScada{}

	csr, _ := DB().Connection.NewQuery().From(new(ScadaDataNew).TableName()).
		Where(dbox.And(filter...)).Cursor(nil) //.Order("turbine")

	e = csr.Fetch(&scadas, 0, false)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}
	defer csr.Close()

	if len(p.Turbine) == 0 {
		for _, scadaVal := range scadas {
			exist := false
			for _, val := range turbine {
				if scadaVal.Turbine == val {
					exist = true
				}
			}
			if exist == false {
				turbine = append(turbine, scadaVal.Turbine)
			}
		}
	}

	for _, turbineVal := range turbine {
		coId++
		groupdata := toolkit.M{}
		groupdata.Set("Index", coId)
		groupdata.Set("Name", turbineVal)

		if turbineVal != "MetTower" {

			data := c.From(&scadas).Where(func(x interface{}) interface{} {
				y := x.(MiniScada)

				return y.Turbine == turbineVal
			}).Exec().Result.Data().([]MiniScada)

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
				for _, data := range results {
					wsCategoryList = append(wsCategoryList, toolkit.ToString(data.DirectionNo)+
						"_"+toolkit.ToString(data.WsCategoryNo)+"_"+data.WsCategoryDesc)
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
			data := []MiniMetTower{}

			var metTowerFilter []*dbox.Filter
			metTowerFilter = append(metTowerFilter, dbox.Gte("dateinfo.dateid", tStart))
			metTowerFilter = append(metTowerFilter, dbox.Lte("dateinfo.dateid", tEnd))

			csrMet, _ := DB().Connection.NewQuery().From(new(MetTower).TableName()).
				Where(dbox.And(metTowerFilter...)).Cursor(nil)

			e = csrMet.Fetch(&data, 0, false)
			if e != nil {
				return helper.CreateResult(false, nil, e.Error())
			}
			defer csrMet.Close()

			if toolkit.SliceLen(data) > 0 {
				totalDuration := float64(len(data)) // * 10.0 /* Tot data * 2 for get total minutes*/
				datas := c.From(&data).Apply(func(x interface{}) interface{} {
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
				for _, data := range results {
					wsCategoryList = append(wsCategoryList, toolkit.ToString(data.DirectionNo)+
						"_"+toolkit.ToString(data.WsCategoryNo)+"_"+data.WsCategoryDesc)
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

	data := struct {
		WindRose toolkit.Ms
		MaxValue float64
	}{
		WindRose: WindRoseResult,
		MaxValue: maxValue,
	}

	return helper.CreateResult(true, data, "success")

}
