package controller

import (
	_ "eaciit/wfdemo-git-dev/library/core"
	_ "eaciit/wfdemo-git-dev/library/models"
	"eaciit/wfdemo-git-dev/web/helper"
	_ "fmt"
	_ "github.com/eaciit/crowd"
	_ "github.com/eaciit/dbox"
	"github.com/eaciit/knot/knot.v1"
	_ "github.com/eaciit/toolkit"
	_ "strconv"
	_ "strings"
	_ "time"
)

type AnalyticHistogramController struct {
	App
}

func CreateAnalyticHistogramController() *AnalyticHistogramController {
	var controller = new(AnalyticHistogramController)
	return controller
}

func (m *AnalyticHistogramController) GetHistogramData(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	// p := new(PayloadAnalytic)
	// e := k.GetPayload(&p)

	// if e != nil {
	// 	return helper.CreateResult(false, nil, e.Error())
	// }

	// tStart, _ := time.Parse("2006-01-02", p.DateStart.UTC().Format("2006-01-02"))
	// tEnd, _ := time.Parse("2006-01-02 15:04:05", p.DateEnd.UTC().Format("2006-01-02")+" 23:59:59")
	// turbine := p.Turbine
	// project := p.Project

	// match := tk.M{}
	// match.Set("dateinfo.dateid", tk.M{}.Set("$lte", tEnd).Set("$gte", tStart))
	// match.Set("projectname", project)
	// match.Set("avgwindspeed", tk.M{}.Set("$gte", 3).Set("$lte", 25))

	// if len(turbine) > 0 {
	// 	match.Set("turbine", tk.M{}.Set("$in", turbine))
	// }

	// group := tk.M{
	// 	"_id":   "$wsavgforpc",
	// 	"total": tk.M{}.Set("$sum", 1),
	// }

	// sort := tk.M{
	// 	"_id": 1,
	// }

	// var pipes []tk.M
	// pipes = append(pipes, tk.M{}.Set("$match", match))
	// pipes = append(pipes, tk.M{}.Set("$group", group))
	// pipes = append(pipes, tk.M{}.Set("$sort", sort))

	return helper.CreateResult(false, nil, "")

	// csr, e := DB().Connection.NewQuery().
	// 	From(new(ScadaData).TableName()).
	// 	Command("pipe", pipes).
	// 	Cursor(nil)

	// defer csr.Close()

	// if e != nil {
	// 	return helper.CreateResult(false, nil, "Error query : "+e.Error())
	// }

	// results := make([]tk.M, 0)
	// e = csr.Fetch(&results, 0, false)

	// if e != nil {
	// 	return helper.CreateResult(false, nil, "Error facing results : "+e.Error())
	// }

	// totalData := c.From(&results).Sum(func(x interface{}) interface{} {
	// 	dt := x.(tk.M)
	// 	return dt["total"].(int)
	// }).Exec().Result.Sum

	// valuewindspeed := tk.M{"3.0": 0}
	// valuewindspeed.Set("3.5", 0)

	// categorywindspeed := []string{}
	// categorywindspeed = append(categorywindspeed, "3 - 3.5")
	// categorywindspeed = append(categorywindspeed, "3.5 - 4")
	// for i := 4; i <= 24; i++ {
	// 	nextPhase := i + 1
	// 	categorywindspeed = append(categorywindspeed, strconv.Itoa(i)+" - "+strconv.Itoa(nextPhase))
	// 	valuewindspeed.Set(strconv.Itoa(i)+".0", 0)
	// }

	// for _, x := range results {
	// 	id := tk.RoundingAuto64(x["_id"].(float64), 1)
	// 	total := x["total"].(int)
	// 	value := tk.Div(float64(total), totalData)

	// 	sId := strconv.FormatFloat(id, 'f', 1, 64)

	// 	valuewindspeed.Set(sId, value)
	// }

	// retvaluews := []float64{}
	// retvaluews = append(retvaluews, valuewindspeed.GetFloat64("3.0"))
	// retvaluews = append(retvaluews, valuewindspeed.GetFloat64("3.5"))
	// for i := 4; i <= 24; i++ {
	// 	retvaluews = append(retvaluews, valuewindspeed.GetFloat64(strconv.Itoa(i)+".0")*100)
	// }

	// data := tk.M{
	// 	"categorywindspeed": categorywindspeed,
	// 	"valuewindspeed":    retvaluews,
	// 	"totaldata":         totalData,
	// }

	// return helper.CreateResult(true, data, "success")
}

// type PayloadHistogram struct {
// 	MaxValue float64
// 	MinValue float64
// 	BinValue int
// 	Filter   PayloadAnalytic
// }

// func (m *AnalyticHistogramController) GetProductionHistogramData(k *knot.WebContext) interface{} {
// 	k.Config.OutputType = knot.OutputJson

// 	p := new(PayloadHistogram)
// 	e := k.GetPayload(&p)

// 	if e != nil {
// 		return helper.CreateResult(false, nil, e.Error())
// 	}

// 	categoryproduction := []string{}
// 	valueproduction := []float64{}
// 	interval := (p.MaxValue - p.MinValue) / float64(p.BinValue)
// 	startcategory := p.MinValue
// 	totalData := 0.0

// 	for i := 0; i < (p.BinValue); i++ {
// 		categoryproduction = append(categoryproduction, fmt.Sprintf("%.2f", startcategory)+" : "+fmt.Sprintf("%.2f", (startcategory+interval)))

// 		match := tk.M{}
// 		match.Set("power", tk.M{}.Set("$lt", (startcategory+interval)).Set("$gte", startcategory))

// 		group := tk.M{
// 			"_id":   "",
// 			"total": tk.M{}.Set("$sum", 1),
// 		}

// 		var pipes []tk.M
// 		pipes = append(pipes, tk.M{}.Set("$match", match))
// 		pipes = append(pipes, tk.M{}.Set("$group", group))

// 		csr, e := DB().Connection.NewQuery().
// 			From(new(ScadaData).TableName()).
// 			Command("pipe", pipes).
// 			Cursor(nil)

// 		defer csr.Close()

// 		if e != nil {
// 			return helper.CreateResult(false, nil, "Error query : "+e.Error())
// 		}

// 		resultCategory := []tk.M{}
// 		e = csr.Fetch(&resultCategory, 0, false)

// 		valueproduction = append(valueproduction, float64(resultCategory[0]["total"].(int)))

// 		totalData = totalData + float64(resultCategory[0]["total"].(int))

// 		startcategory = startcategory + interval
// 	}

// 	for i := 0; i < len(valueproduction); i++ {
// 		valueproduction[i] = float64(int((valueproduction[i]/totalData*100)*100)) / 100
// 	}

// 	data := tk.M{
// 		"categoryproduction": categoryproduction,
// 		"valueproduction":    valueproduction,
// 		"totaldata":          totalData,
// 	}

// 	return helper.CreateResult(true, data, "success")
// }
