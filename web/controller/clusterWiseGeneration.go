package controller

import (
	. "eaciit/wfdemo-git/library/core"
	. "eaciit/wfdemo-git/library/models"
	"eaciit/wfdemo-git/web/helper"
	"sort"
	"strings"
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

func (m *ClusterWiseGeneration) GetDataDGR(k *knot.WebContext) interface{} {
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

	_project, _arrturb := p.Project, []interface{}{}
	for _, aturb := range turbineList {
		cond := false

		for _, iturb := range p.Turbine {
			if tk.ToString(iturb) == aturb.Value {
				cond = true
				break
			}
		}

		if cond {
			_project = aturb.DgrProject
			_arrturb = append(_arrturb, aturb.DgrTurbine)
		}
	}

	if _project == "" {
		_project = p.Project
	}

	dataSeries := []tk.M{}

	reffdgrturb := getturbinedgr(p.Project)

	ids := tk.M{"project": "$chosensite", "turbine": "$turbine"}

	matches := tk.M{"dateinfo.dateid": tk.M{"$gte": tStart, "$lte": tEnd}}
	matches.Set("chosensite", _project)
	matches.Set("turbine", tk.M{}.Set("$in", _arrturb))

	pipe := []tk.M{{"$match": matches},
		{"$group": tk.M{
			"_id":       ids,
			"genkwhday": tk.M{"$sum": "$genkwhday"},
			"count":     tk.M{"$sum": 1},
		}},
		{"$sort": tk.M{"_id.project": 1}}}

	// tk.Println(pipe)

	csr, e := DB().Connection.NewQuery().
		From("rpt_generation").
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

		for _, v := range turbineList {
			compval := reffdgrturb.GetString(v.Value)
			if _id.GetString("turbine") == compval {
				ds.Set("turbine", v.Turbine)
				ds.Set("cluster", v.Cluster)
				ds.Set("capacity", v.Capacity)
			}
		}

		ds.Set("sumGeneration", tk.M{}.Set("value", tkm.GetFloat64("genkwhday")/1000).Set("type", "generation"))
		// ds.Set("averageMa", tk.M{}.Set("value", res.GetFloat64("machineavailability")/100).Set("type", "avail"))
		// ds.Set("averageGa", tk.M{}.Set("value", res.GetFloat64("gridavailability")/100).Set("type", "avail"))

		ds.Set("totalHours", tkm.GetFloat64("count")*24)

		// dataSeries = append(dataSeries, ds)
		alltkm.Set(ds.GetString("turbine"), ds)
		allturb = append(allturb, ds.GetString("turbine"))
	}

	ids = tk.M{"category": "$category", "turbine": "$turbine"}

	matches.Set("category", tk.M{"$ne": ""})
	matches.Set("subcategory", tk.M{"$nin": []string{"NOR", "Penalty Claim"}})

	pipe = []tk.M{{"$match": matches},
		{"$group": tk.M{
			"_id":            ids,
			"breakdownhours": tk.M{"$sum": "$breakdownhours"},
		}},
		{"$sort": tk.M{"_id.turbine": 1}}}

	csrdown, e := DB().Connection.NewQuery().
		From("rpt_downtime").
		Command("pipe", pipe).
		Cursor(nil)

	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}
	defer csrdown.Close()

	isRACategory := func(str string) bool {
		arr := []string{"S-M/C", "U-M/C", "ROW(M/C)-OEM", "S-IG", "U-IG", "S-PSS", "U-PSS", "ROW(IG)-OEM", "AOR"}
		for _, as := range arr {
			if as == str {
				return true
			}
		}
		return false
	}

	for {
		tkm, ds := tk.M{}, tk.M{}
		e = csrdown.Fetch(&tkm, 1, false)
		// tk.Println(tkm)

		if e != nil {
			break
		}

		_id := tkm.Get("_id", tk.M{}).(tk.M)
		for _, v := range turbineList {
			compval := reffdgrturb.GetString(v.Value)
			if _id.GetString("turbine") == compval {
				ds = alltkm.Get(v.Turbine, tk.M{}).(tk.M)
			}
		}

		sGridBreak, sMachineBreak := ds.GetFloat64("sGridBreak"), ds.GetFloat64("sMachineBreak")
		minA, sRABreak := ds.GetFloat64("minA"), ds.GetFloat64("sRABreak")

		switch strings.ToUpper(_id.GetString("category")) {
		case "S-M/C", "U-M/C", "ROW(M/C)-OEM", "AOR":
			sMachineBreak += tkm.GetFloat64("breakdownhours")
		case "S-EG", "U-EG":
			minA += tkm.GetFloat64("breakdownhours")
			sGridBreak += tkm.GetFloat64("breakdownhours")
		case "ROW(IG)-NONOEM", "ROW(M/C)-NONOEM", "FM-ENV", "FM-THEFT":
			minA += tkm.GetFloat64("breakdownhours")
		}

		if isRACategory(strings.ToUpper(_id.GetString("category"))) {
			sRABreak += tkm.GetFloat64("breakdownhours")
		}

		ds.Set("sGridBreak", sGridBreak).Set("sMachineBreak", sMachineBreak).Set("sRABreak", sRABreak).Set("minA", minA)

		alltkm.Set(ds.GetString("turbine"), ds)
	}

	sort.Strings(allturb)
	for _, turb := range allturb {
		ds := alltkm.Get(turb, tk.M{}).(tk.M)

		// Cacl MA, RA, GA
		totalHours := ds.GetFloat64("totalHours")
		sGridBreak, sMachineBreak := ds.GetFloat64("sGridBreak"), ds.GetFloat64("sMachineBreak")
		minA, sRABreak := ds.GetFloat64("minA"), ds.GetFloat64("sRABreak")

		A := totalHours - minA
		averageGa := tk.Div(totalHours-sGridBreak, totalHours) * 100
		averageMa := tk.Div(A-sMachineBreak, A) * 100
		averageRa := tk.Div(A-sRABreak, A) * 100

		// _ = sRABreak

		ds.Set("averageMa", tk.M{}.Set("value", averageMa/100).Set("type", "avail"))
		ds.Set("averageGa", tk.M{}.Set("value", averageGa/100).Set("type", "avail"))
		ds.Set("averageRa", tk.M{}.Set("value", averageRa/100).Set("type", "avail"))

		dataSeries = append(dataSeries, ds)
	}

	result := tk.M{}.Set("totalturbine", len(p.Turbine)).Set("data", dataSeries).Set("projectname", p.Project)
	return helper.CreateResult(true, result, "success")
}

func (m *ClusterWiseGeneration) GetDataDGRCompare(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	type DateDetail struct {
		Period    string
		DateStart time.Time
		DateEnd   time.Time
	}

	type PayloadDGRCompare struct {
		Period          string
		Project         string
		Turbine         []interface{}
		DateStart       time.Time
		DateEnd         time.Time
		ColName         string
		DeviationStatus bool
		Deviation       float64
		Details         []DateDetail
	}

	p := new(PayloadDGRCompare)
	e := k.GetPayload(&p)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	// tStart, tEnd, e := helper.GetStartEndDate(k, p.Period, p.DateStart, p.DateEnd)
	// if e != nil {
	// 	return helper.CreateResult(false, nil, e.Error())
	// }

	var turbineList []TurbineOut
	if p.Project != "" {
		turbineList, _ = helper.GetTurbineList([]interface{}{p.Project})
	} else {
		turbineList, _ = helper.GetTurbineList(nil)
	}

	_project, _arrturb := p.Project, []interface{}{}
	for _, aturb := range turbineList {
		cond := false

		for _, iturb := range p.Turbine {
			if tk.ToString(iturb) == aturb.Value {
				cond = true
				break
			}
		}

		if cond {
			_project = aturb.DgrProject
			_arrturb = append(_arrturb, aturb.DgrTurbine)
		}
	}

	if _project == "" {
		_project = p.Project
	}

	allSeries, allCompare := []interface{}{}, []string{}
	reffdgrturb := getturbinedgr(p.Project)

	ids := tk.M{"project": "$chosensite", "turbine": "$turbine"}

	for _, ddt := range p.Details {
		dataSeries := []tk.M{}
		// tk.Println(ddt)
		tStart, tEnd, e := helper.GetStartEndDate(k, ddt.Period, ddt.DateStart, ddt.DateEnd)
		if e != nil {
			continue
		}

		allCompare = append(allCompare, tk.Sprintf("%s to %s", tStart.Format("02-Jan-2006"), tEnd.Format("02-Jan-2006")))

		matches := tk.M{"dateinfo.dateid": tk.M{"$gte": tStart, "$lte": tEnd}}
		matches.Set("chosensite", _project)
		matches.Set("turbine", tk.M{}.Set("$in", _arrturb))

		pipe := []tk.M{{"$match": matches},
			{"$group": tk.M{
				"_id":       ids,
				"genkwhday": tk.M{"$sum": "$genkwhday"},
				"count":     tk.M{"$sum": 1},
			}},
			{"$sort": tk.M{"_id.project": 1}}}

		csr, e := DB().Connection.NewQuery().
			From("rpt_generation").
			Command("pipe", pipe).
			Cursor(nil)

		if e != nil {
			return helper.CreateResult(false, nil, e.Error())
		}

		alltkm, allturb := tk.M{}, []string{}
		for {
			tkm, ds := tk.M{}, tk.M{}
			e = csr.Fetch(&tkm, 1, false)

			if e != nil {
				break
			}

			_id := tkm.Get("_id", tk.M{}).(tk.M)

			for _, v := range turbineList {
				compval := reffdgrturb.GetString(v.Value)
				if _id.GetString("turbine") == compval {
					ds.Set("turbine", v.Turbine)
					ds.Set("cluster", v.Cluster)
					ds.Set("capacity", v.Capacity)
				}
			}

			ds.Set("sumGeneration", tk.M{}.Set("value", tkm.GetFloat64("genkwhday")/1000).Set("type", "generation"))
			// ds.Set("averageMa", tk.M{}.Set("value", res.GetFloat64("machineavailability")/100).Set("type", "avail"))
			// ds.Set("averageGa", tk.M{}.Set("value", res.GetFloat64("gridavailability")/100).Set("type", "avail"))

			ds.Set("totalHours", tkm.GetFloat64("count")*24)

			// dataSeries = append(dataSeries, ds)
			alltkm.Set(ds.GetString("turbine"), ds)
			allturb = append(allturb, ds.GetString("turbine"))
		}

		csr.Close()

		ids = tk.M{"category": "$category", "turbine": "$turbine"}

		matches.Set("category", tk.M{"$ne": ""})
		matches.Set("subcategory", tk.M{"$nin": []string{"NOR", "Penalty Claim"}})

		pipe = []tk.M{{"$match": matches},
			{"$group": tk.M{
				"_id":            ids,
				"breakdownhours": tk.M{"$sum": "$breakdownhours"},
			}},
			{"$sort": tk.M{"_id.turbine": 1}}}

		csrdown, e := DB().Connection.NewQuery().
			From("rpt_downtime").
			Command("pipe", pipe).
			Cursor(nil)

		if e != nil {
			return helper.CreateResult(false, nil, e.Error())
		}
		defer csrdown.Close()

		isRACategory := func(str string) bool {
			arr := []string{"S-M/C", "U-M/C", "ROW(M/C)-OEM", "S-IG", "U-IG", "S-PSS", "U-PSS", "ROW(IG)-OEM", "AOR"}
			for _, as := range arr {
				if as == str {
					return true
				}
			}
			return false
		}

		for {
			tkm, ds := tk.M{}, tk.M{}
			e = csrdown.Fetch(&tkm, 1, false)
			// tk.Println(tkm)

			if e != nil {
				break
			}

			_id := tkm.Get("_id", tk.M{}).(tk.M)
			for _, v := range turbineList {
				compval := reffdgrturb.GetString(v.Value)
				if _id.GetString("turbine") == compval {
					ds = alltkm.Get(v.Turbine, tk.M{}).(tk.M)
				}
			}

			sGridBreak, sMachineBreak := ds.GetFloat64("sGridBreak"), ds.GetFloat64("sMachineBreak")
			minA, sRABreak := ds.GetFloat64("minA"), ds.GetFloat64("sRABreak")

			switch strings.ToUpper(_id.GetString("category")) {
			case "S-M/C", "U-M/C", "ROW(M/C)-OEM", "AOR":
				sMachineBreak += tkm.GetFloat64("breakdownhours")
			case "S-EG", "U-EG":
				minA += tkm.GetFloat64("breakdownhours")
				sGridBreak += tkm.GetFloat64("breakdownhours")
			case "ROW(IG)-NONOEM", "ROW(M/C)-NONOEM", "FM-ENV", "FM-THEFT":
				minA += tkm.GetFloat64("breakdownhours")
			}

			if isRACategory(strings.ToUpper(_id.GetString("category"))) {
				sRABreak += tkm.GetFloat64("breakdownhours")
			}

			ds.Set("sGridBreak", sGridBreak).Set("sMachineBreak", sMachineBreak).Set("sRABreak", sRABreak).Set("minA", minA)

			alltkm.Set(ds.GetString("turbine"), ds)
		}

		sort.Strings(allturb)
		for _, turb := range allturb {
			ds := alltkm.Get(turb, tk.M{}).(tk.M)

			// Cacl MA, RA, GA
			totalHours := ds.GetFloat64("totalHours")
			sGridBreak, sMachineBreak := ds.GetFloat64("sGridBreak"), ds.GetFloat64("sMachineBreak")
			minA, sRABreak := ds.GetFloat64("minA"), ds.GetFloat64("sRABreak")

			A := totalHours - minA
			averageGa := tk.Div(totalHours-sGridBreak, totalHours) * 100
			averageMa := tk.Div(A-sMachineBreak, A) * 100
			averageRa := tk.Div(A-sRABreak, A) * 100

			// _ = sRABreak

			ds.Set("averageMa", tk.M{}.Set("value", averageMa/100).Set("type", "avail"))
			ds.Set("averageGa", tk.M{}.Set("value", averageGa/100).Set("type", "avail"))
			ds.Set("averageRa", tk.M{}.Set("value", averageRa/100).Set("type", "avail"))

			dataSeries = append(dataSeries, ds)
		}
		// tk.Println(len(dataSeries))
		allSeries = append(allSeries, dataSeries)

	}

	result := tk.M{}.Set("totalturbine", len(p.Turbine)).Set("data", allSeries).Set("projectname", p.Project).Set("compare", allCompare)
	return helper.CreateResult(true, result, "success")
}
