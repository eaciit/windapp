package controller

import (
	. "eaciit/wfdemo-git/library/core"
	. "eaciit/wfdemo-git/library/models"
	"eaciit/wfdemo-git/web/helper"
	"sort"
	"time"

	"github.com/eaciit/knot/knot.v1"
	tk "github.com/eaciit/toolkit"
)

type ClusterWiseGeneration struct {
	App
}

func CreateClusterWiseGenerationController() *ClusterWiseGeneration {
	var controller = new(ClusterWiseGeneration)
	return controller
}

func (m *ClusterWiseGeneration) GetData(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	p := new(PayloadAnalyticTLP)
	e := k.GetPayload(&p)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	tStart, tEnd, e := helper.GetStartEndDate(k, p.Period, p.DateStart, p.DateEnd)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	var turbineList []TurbineOut
	if p.Project != "" {
		turbineList, _ = helper.GetTurbineList([]interface{}{p.Project})
	} else {
		turbineList, _ = helper.GetTurbineList(nil)
	}

	dataSeries := []tk.M{}

	ids := tk.M{"project": "$projectname", "turbine": "$turbine"}

	matches := tk.M{"dateinfo.dateid": tk.M{"$gte": tStart, "$lte": tEnd}}
	matches.Set("projectname", p.Project)
	matches.Set("turbine", tk.M{}.Set("$in", p.Turbine))

	pipe := []tk.M{{"$match": matches},
		{"$group": tk.M{
			"_id":            ids,
			"production":     tk.M{"$sum": "$production"},
			"lostenergy":     tk.M{"$sum": "$lostenergy"},
			"mdownhours":     tk.M{"$sum": "$machinedownhours"},
			"gdownhours":     tk.M{"$sum": "$griddownhours"},
			"odownhours":     tk.M{"$sum": "$otherdowntimehours"},
			"oktime":         tk.M{"$sum": "$oktime"},
			"totaltimestamp": tk.M{"$sum": 1},
			"power":          tk.M{"$sum": "$powerkw"},
			"maxdate":        tk.M{"$max": "$dateinfo.dateid"},
			"mindate":        tk.M{"$min": "$dateinfo.dateid"},
		}},
		{"$sort": tk.M{"_id.project": 1}}}

	csr, e := DB().Connection.NewQuery().
		From(new(ScadaSummaryDaily).TableName()).
		Command("pipe", pipe).
		Cursor(nil)

	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}
	defer csr.Close()

	alltkm, allturb := tk.M{}, []string{}
	for {
		tkm, ds := tk.M{}, tk.M{}
		e = csr.Fetch(&tkm, 1, false)

		if e != nil {
			break
		}

		_id := tkm.Get("_id", tk.M{}).(tk.M)

		minDate := tkm.Get("mindate").(time.Time).UTC()
		maxDate := tkm.Get("maxdate").(time.Time).UTC()

		hourValue := maxDate.UTC().Sub(minDate.UTC()).Hours()

		capacity := float64(0)
		for _, v := range turbineList {
			if _id.GetString("turbine") == v.Value {
				capacity = v.Capacity
				ds.Set("turbine", v.Turbine)
				ds.Set("cluster", v.Cluster)
			}
		}

		in := tk.M{}.Set("noofturbine", 1).Set("oktime", tkm.GetFloat64("oktime")/3600).Set("energy", tkm.GetFloat64("production")).
			Set("totalhour", hourValue).Set("totalcapacity", capacity).
			Set("machinedowntime", tkm.GetFloat64("mdownhours")).Set("griddowntime", tkm.GetFloat64("gdownhours")).
			Set("otherdowntime", tkm.GetFloat64("odownhours"))

		res := helper.CalcAvailabilityAndPLF(in)

		ds.Set("sumGeneration", tk.M{}.Set("value", tkm.GetFloat64("power")/1000).Set("type", "generation"))
		ds.Set("averageMa", tk.M{}.Set("value", res.GetFloat64("machineavailability")/100).Set("type", "avail"))
		ds.Set("averageGa", tk.M{}.Set("value", res.GetFloat64("gridavailability")/100).Set("type", "avail"))

		// dataSeries = append(dataSeries, ds)
		alltkm.Set(ds.GetString("turbine"), ds)
		allturb = append(allturb, ds.GetString("turbine"))
	}

	sort.Strings(allturb)
	for _, turb := range allturb {
		dataSeries = append(dataSeries, alltkm.Get(turb, tk.M{}).(tk.M))
	}

	result := tk.M{}.Set("totalturbine", len(p.Turbine)).Set("data", dataSeries).Set("projectname", p.Project)
	return helper.CreateResult(true, result, "success")
}
