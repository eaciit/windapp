package controller

import (
	. "eaciit/wfdemo-git/library/core"
	. "eaciit/wfdemo-git/library/helper"
	. "eaciit/wfdemo-git/library/models"
	"eaciit/wfdemo-git/web/helper"
	c "github.com/eaciit/crowd"
	"github.com/eaciit/dbox"
	"github.com/eaciit/knot/knot.v1"
	"github.com/eaciit/toolkit"

	_ "github.com/eaciit/dbox/dbc/mongo"

	"math"
	_ "strconv"
	_ "strings"
	"time"
)

type AnalyticWindRoseFlexiController struct {
	App
}

func CreateAnalyticWindRoseFlexiController() *AnalyticWindRoseFlexiController {
	var controller = new(AnalyticWindRoseFlexiController)
	return controller
}

type DataItem struct {
	DirectionNo    int
	DirectionDesc  int
	WsCategoryNo   int
	WsCategoryDesc string
	Hours          float64
	Frequency      int
}

type DataItemResult struct {
	DirectionNo    int
	DirectionDesc  int
	WsCategoryNo   int
	WsCategoryDesc string
	Hours          float64
	Contribution   float64
	Frequency      int
}

type DataItemResultGrid struct {
	WsCategoryNo   int
	WsCategoryDesc string
	Hours          float64
	Frequency      float64
}

type DataItemGroup struct {
	DirectionNo    int
	DirectionDesc  int
	WsCategoryNo   int
	WsCategoryDesc string
}

type DataItemGroupGrid struct {
	WsCategoryNo   int
	WsCategoryDesc string
}

type DataGroupResult struct {
	ProjectName string
	Turbine     string
	DateId      time.Time
}

type ContributeGroupResult struct {
	WsCategoryNo   int
	WsCategoryDesc string
}

type ContributeItemResult struct {
	WsCategoryNo   int
	WsCategoryDesc string
	Hours          float64
	Contribution   float64
	Frequency      int
}

func GetWsCategory(ws float64) (int, string) {
	catNo := 0
	catDesc := "0 to 4m/s"
	if ws >= 14 {
		catNo = 4
		catDesc = "14 and above"
	} else if ws >= 9 {
		catNo = 3
		catDesc = "9 to 14m/s"
	} else if ws >= 7 {
		catNo = 2
		catDesc = "7 to 9m/s"
	} else if ws >= 4 {
		catNo = 1
		catDesc = "4 to 7m/s"
	}

	return catNo, catDesc
}

func GetDirection(windDir float64, nacelPos float64, breakDown int) (int, int) {
	dirNo := 0
	devide := 360.0 / toolkit.ToFloat64(breakDown, 0, toolkit.RoundingUp)

	// dirDescs := make([]int)
	dirDescs := make([]int, breakDown)
	direction := 0
	plusValue := 360 / breakDown
	for j := 0; j < breakDown; j++ {
		dirDescs[j] = direction
		direction = direction + plusValue
		// dirDescs := append(dirDescs, direction)
	}
	// dirDescs := []string{"N", "NE", "E", "SE", "S", "SW", "W", "NW"}
	if windDir < 0 {
		windDir = 360.0 + windDir
	}
	if nacelPos < 0 {
		nacelPos = 360.0 + nacelPos
	}
	if nacelPos < 0 {
		nacelPos = 360.0 + nacelPos
	}
	dirCalc := math.Mod((nacelPos + windDir), 360.0)
	dirNo = int(toolkit.RoundingAuto64(dirCalc/devide, 0))

	if dirNo > (breakDown - 1) {
		dirNo = 0
	}

	return dirNo, dirDescs[dirNo]
}

func (m *AnalyticWindRoseFlexiController) GetFlexiData(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	var (
		filter []*dbox.Filter
		// pipes 			[]toolkit.M
		scadas []ScadaData
	)

	p := new(PayloadAnalytic)
	e := k.GetPayload(&p)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	tStart, _ := time.Parse("2006-01-02", p.DateStart.UTC().Format("2006-01-02"))
	tEnd, _ := time.Parse("2006-01-02 15:04:05", p.DateEnd.UTC().Format("2006-01-02")+" 23:59:59")
	turbine := p.Turbine
	project := p.Project
	breakDown := p.BreakDown // Section

	filter = append(filter, dbox.Ne("_id", ""))
	filter = append(filter, dbox.Gte("dateinfo.dateid", tStart))
	filter = append(filter, dbox.Lte("dateinfo.dateid", tEnd))

	if project != "" {
		filter = append(filter, dbox.Eq("projectname", project))
	}
	if len(turbine) != 0 {
		filter = append(filter, dbox.In("turbine", turbine...))
	}

	csr, e := DB().Connection.NewQuery().From(new(ScadaData).TableName()).Where(dbox.And(filter...)).Cursor(nil)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	e = csr.Fetch(&scadas, 0, false)

	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	defer csr.Close()

	totalDuration := float64((len(scadas) * 10.0)) / 60.0

	datas := c.From(&scadas).Apply(func(x interface{}) interface{} {
		dt := x.(ScadaData)
		var di DataItem

		dirNo, dirDesc := GetDirection(dt.WindDirection, dt.NacelDirection, toolkit.ToInt(breakDown, toolkit.RoundingAuto))
		wsNo, wsDesc := GetWsCategory(dt.AvgWindSpeed)

		di.DirectionNo = dirNo
		di.DirectionDesc = dirDesc
		di.WsCategoryNo = wsNo
		di.WsCategoryDesc = wsDesc
		di.Hours = 10.0 / 60.0
		di.Frequency = 1

		return di
	}).Exec().Group(func(x interface{}) interface{} {
		dt := x.(DataItem)

		var dig DataItemGroup
		dig.DirectionNo = dt.DirectionNo
		dig.DirectionDesc = dt.DirectionDesc
		dig.WsCategoryNo = dt.WsCategoryNo
		dig.WsCategoryDesc = dt.WsCategoryDesc

		return dig
	}, nil).Exec()

	dts := datas.Apply(func(x interface{}) interface{} {
		kv := x.(c.KV)
		vv := kv.Key.(DataItemGroup)
		vs := kv.Value.([]DataItem)

		sumDuration := c.From(&vs).Sum(func(x interface{}) interface{} {
			dt := x.(DataItem)
			return dt.Hours
		}).Exec().Result.Sum

		var di DataItemResult
		di.DirectionNo = vv.DirectionNo
		di.DirectionDesc = vv.DirectionDesc
		di.WsCategoryNo = vv.WsCategoryNo
		di.WsCategoryDesc = vv.WsCategoryDesc
		di.Hours = sumDuration
		contribute := 0.0
		contribute = sumDuration / totalDuration
		di.Contribution = RoundUp(contribute, .5, 2)

		sumFreq := c.From(&vs).Sum(func(x interface{}) interface{} {
			dt := x.(DataItem)
			return dt.Frequency
		}).Exec().Result.Sum
		di.Frequency = int(sumFreq)

		return di
	}).Exec()

	results := dts.Result.Data().([]DataItemResult)

	gridDt := c.From(&results).Apply(func(x interface{}) interface{} {
		dt := x.(DataItemResult)

		return dt
	}).Exec().Group(func(x interface{}) interface{} {
		dt := x.(DataItemResult)

		var dig DataItemGroupGrid
		dig.WsCategoryNo = dt.WsCategoryNo
		dig.WsCategoryDesc = dt.WsCategoryDesc

		return dig
	}, nil).Exec()

	resGridDt := gridDt.Apply(func(x interface{}) interface{} {
		kv := x.(c.KV)
		keys := kv.Key.(DataItemGroupGrid)
		vs := kv.Value.([]DataItemResult)

		Hours := c.From(&vs).Sum(func(x interface{}) interface{} {
			dt := x.(DataItemResult)
			return dt.Hours
		}).Exec().Result.Sum

		Frequency := c.From(&vs).Sum(func(x interface{}) interface{} {
			dt := x.(DataItemResult)
			return dt.Frequency
		}).Exec().Result.Sum

		var di DataItemResultGrid
		di.WsCategoryNo = keys.WsCategoryNo
		di.WsCategoryDesc = keys.WsCategoryDesc
		di.Hours = Hours
		di.Frequency = Frequency

		return di
	}).Exec().Result.Data().([]DataItemResultGrid)

	data := struct {
		WindRose     []DataItemResult
		GridWindrose []DataItemResultGrid
	}{
		WindRose:     results,
		GridWindrose: resGridDt,
	}

	return helper.CreateResult(true, data, "success")

}

func (m *AnalyticWindRoseFlexiController) GetFlexiDataEachTurbine(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	var (
		scadas         []ScadaData
		WindRoseResult []toolkit.M
	)

	p := new(PayloadAnalytic)
	e := k.GetPayload(&p)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	tStart, _ := time.Parse("2006-01-02", p.DateStart.UTC().Format("2006-01-02"))
	tEnd, _ := time.Parse("2006-01-02 15:04:05", p.DateEnd.UTC().Format("2006-01-02")+" 23:59:59")
	turbine := p.Turbine
	breakDown := p.BreakDown // Section

	coId := 0
	turbine = append(turbine, "All")
	for _, turbineX := range turbine {
		groupdata := toolkit.M{}
		groupdata.Set("Index", coId)
		groupdata.Set("Name", turbineX.(string))
		coId++

		var filter []*dbox.Filter
		filter = append(filter, dbox.Ne("_id", ""))
		filter = append(filter, dbox.Gte("dateinfo.dateid", tStart))
		filter = append(filter, dbox.Lte("dateinfo.dateid", tEnd))
		if turbineX != "All" {
			filter = append(filter, dbox.Eq("turbine", turbineX))
		}

		csr, e := DB().Connection.NewQuery().From(new(ScadaData).TableName()).Where(dbox.And(filter...)).Cursor(nil)
		if e != nil {
			return helper.CreateResult(false, nil, e.Error())
		}

		e = csr.Fetch(&scadas, 0, false)

		if e != nil {
			return helper.CreateResult(false, nil, e.Error())
		}

		defer csr.Close()

		totalDuration := float64((len(scadas) * 10.0)) / 60.0
		toolkit.Printf("Total Data : %s \n", totalDuration)

		datas := c.From(&scadas).Apply(func(x interface{}) interface{} {
			dt := x.(ScadaData)
			var di DataItem

			dirNo, dirDesc := GetDirection(dt.WindDirection, dt.NacelDirection, toolkit.ToInt(breakDown, toolkit.RoundingAuto))
			wsNo, wsDesc := GetWsCategory(dt.AvgWindSpeed)

			di.DirectionNo = dirNo
			di.DirectionDesc = dirDesc
			di.WsCategoryNo = wsNo
			di.WsCategoryDesc = wsDesc
			di.Hours = 10.0 / 60.0
			di.Frequency = 1

			return di
		}).Exec().Group(func(x interface{}) interface{} {
			dt := x.(DataItem)

			var dig DataItemGroup
			dig.DirectionNo = dt.DirectionNo
			dig.DirectionDesc = dt.DirectionDesc
			dig.WsCategoryNo = dt.WsCategoryNo
			dig.WsCategoryDesc = dt.WsCategoryDesc

			return dig
		}, nil).Exec()

		dts := datas.Apply(func(x interface{}) interface{} {
			kv := x.(c.KV)
			vv := kv.Key.(DataItemGroup)
			vs := kv.Value.([]DataItem)

			sumDuration := c.From(&vs).Sum(func(x interface{}) interface{} {
				dt := x.(DataItem)
				return dt.Hours
			}).Exec().Result.Sum

			var di DataItemResult
			di.DirectionNo = vv.DirectionNo
			di.DirectionDesc = vv.DirectionDesc
			di.WsCategoryNo = vv.WsCategoryNo
			di.WsCategoryDesc = vv.WsCategoryDesc
			di.Hours = sumDuration
			contribute := 0.0
			contribute = sumDuration / totalDuration
			di.Contribution = RoundUp(contribute, .5, 2)

			sumFreq := c.From(&vs).Sum(func(x interface{}) interface{} {
				dt := x.(DataItem)
				return dt.Frequency
			}).Exec().Result.Sum
			di.Frequency = int(sumFreq)

			return di
		}).Exec()

		results := dts.Result.Data().([]DataItemResult)

		groupdata.Set("Data", results)
		WindRoseResult = append(WindRoseResult, groupdata)

	}

	data := struct {
		WindRose []toolkit.M
	}{
		WindRose: WindRoseResult,
	}

	return helper.CreateResult(true, data, "success")

}
