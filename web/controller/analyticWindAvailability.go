package controller

import (
	. "eaciit/ostrowfm/library/core"
	. "eaciit/ostrowfm/library/models"
	"eaciit/ostrowfm/web/helper"
	"github.com/eaciit/crowd"
	"github.com/eaciit/knot/knot.v1"
	tk "github.com/eaciit/toolkit"
	// "time"
)

type AnalyticWindAvailabilityController struct {
	App
}

func CreateAnalyticWindAvailabilityController() *AnalyticWindAvailabilityController {
	var controller = new(AnalyticWindAvailabilityController)
	return controller
}

func (m *AnalyticWindAvailabilityController) GetData(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	p := new(PayloadAnalytic)
	e := k.GetPayload(&p)

	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	// tStart, _ := time.Parse("2006-01-02", p.DateStart.UTC().Format("2006-01-02"))
	// tEnd, _ := time.Parse("2006-01-02 15:04:05", p.DateEnd.UTC().Format("2006-01-02")+" 23:59:59")
	tStart, tEnd, e := helper.GetStartEndDate(k, p.Period, p.DateStart, p.DateEnd)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}
	turbine := p.Turbine
	project := p.Project

	match := tk.M{}
	match.Set("dateinfo.dateid", tk.M{}.Set("$lte", tEnd).Set("$gte", tStart))
	match.Set("projectname", project)
	match.Set("avgwindspeed", tk.M{}.Set("$gte", 3))

	if len(turbine) > 0 {
		match.Set("turbine", tk.M{}.Set("$in", turbine))
	}

	group := tk.M{
		"_id":        "$wsavgforpc",
		"energy":     tk.M{}.Set("$sum", "$energy"),
		"totalavail": tk.M{}.Set("$avg", "$totalavail"),
		"totalcount": tk.M{}.Set("$sum", 1),
	}

	sort := tk.M{
		"_id": 1,
	}

	var pipes []tk.M
	pipes = append(pipes, tk.M{}.Set("$match", match))
	pipes = append(pipes, tk.M{}.Set("$group", group))
	pipes = append(pipes, tk.M{}.Set("$sort", sort))

	csr, e := DB().Connection.NewQuery().
		From(new(ScadaData).TableName()).
		Command("pipe", pipes).
		Cursor(nil)

	defer csr.Close()

	if e != nil {
		return helper.CreateResult(false, nil, "Error query : "+e.Error())
	}

	results := make([]tk.M, 0)
	e = csr.Fetch(&results, 0, false)

	if e != nil {
		return helper.CreateResult(false, nil, "Error facing results : "+e.Error())
	}

	totalEnergy := crowd.From(&results).Sum(func(x interface{}) interface{} {
		dt := x.(tk.M)
		return dt["energy"].(float64)
	}).Exec().Result.Sum

	totalTime := crowd.From(&results).Sum(func(x interface{}) interface{} {
		dt := x.(tk.M)
		return dt["totalcount"].(int)
	}).Exec().Result.Sum

	type DataReturn struct {
		WindSpeed  float64
		Energy     float64
		Time       float64
		TotalAvail float64
	}
	dataReturn := make([]DataReturn, 0)
	// duration := tk.RoundingAuto64(tEnd.Sub(tStart).Hours()*3600, 0)
	totalTurbine := 0.0
	_ = totalTurbine
	if len(turbine) > 0 {
		totalTurbine = float64(len(turbine))
	}
	energyAcumulative := 0.0
	timeAcumulative := 0
	for _, d := range results {
		windSpeed := d["_id"].(float64)
		energy := d["energy"].(float64)
		totalavail := d["totalavail"].(float64)
		time := d["totalcount"].(int)

		energyAcumulative = energyAcumulative + energy
		timeAcumulative = timeAcumulative + time
		totalAvail := totalavail * 100
		energyPros := tk.Div(energyAcumulative, totalEnergy) * 100
		timePros := tk.Div(float64(timeAcumulative), float64(totalTime)) * 100

		var dr DataReturn
		dr.WindSpeed = tk.RoundingAuto64(windSpeed, 1)
		dr.Energy = tk.RoundingAuto64(energyPros, 1)
		dr.Time = tk.RoundingAuto64(timePros, 1)
		dr.TotalAvail = tk.RoundingAuto64(totalAvail, 1)

		dataReturn = append(dataReturn, dr)
	}

	return helper.CreateResult(true, dataReturn, "success")
}
