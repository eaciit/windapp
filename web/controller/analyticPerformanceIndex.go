package controller

import (
	. "eaciit/wfdemo-git/library/core"
	. "eaciit/wfdemo-git/library/models"
	"eaciit/wfdemo-git/web/helper"

	"time"

	"github.com/eaciit/dbox"
	"github.com/eaciit/knot/knot.v1"
	tk "github.com/eaciit/toolkit"
)

type AnalyticPerformanceIndexController struct {
	App
}

func CreateAnalyticPerformanceIndexController() *AnalyticPerformanceIndexController {
	var controller = new(AnalyticPerformanceIndexController)
	return controller
}

func (m *AnalyticPerformanceIndexController) GetPerformanceIndex(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	type PerformanceDetail struct {
		Turbine                     string
		PerformanceIndex            float64
		PerformanceIndexLast24Hours float64
		PerformanceIndexLastWeek    float64
		PerformanceIndexMTD         float64
		PerformanceIndexYTD         float64
		Production                  float64
		ProductionLast24Hours       float64
		ProductionLastWeek          float64
		ProductionMTD               float64
		ProductionYTD               float64
		Power                       float64
		PowerLast24Hours            float64
		PowerLastWeek               float64
		PowerMTD                    float64
		PowerYTD                    float64
		StartDate                   time.Time
		EndDate                     time.Time
	}

	type Performance struct {
		Project                     string
		PerformanceIndex            float64
		PerformanceIndexLast24Hours float64
		PerformanceIndexLastWeek    float64
		PerformanceIndexMTD         float64
		PerformanceIndexYTD         float64
		Production                  float64
		ProductionLast24Hours       float64
		ProductionLastWeek          float64
		ProductionMTD               float64
		ProductionYTD               float64
		Power                       float64
		PowerLast24Hours            float64
		PowerLastWeek               float64
		PowerMTD                    float64
		PowerYTD                    float64
		StartDate                   time.Time
		EndDate                     time.Time
		Details                     []PerformanceDetail
	}

	p := new(PayloadAnalytic)
	e := k.GetPayload(&p)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}
	// filter, _ := p.ParseFilter()

	tStart, tEnd, e := helper.GetStartEndDate(k, p.Period, p.DateStart, p.DateEnd)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	turbine := p.Turbine
	project := p.Project

	results := make([]Performance, 0)

	aggrData := []tk.M{}

	for i := 0; i < 5; i++ {
		var filter []*dbox.Filter

		switch i {
		case 3:
			// last24hours
			filter = append(filter, dbox.Gte("dateinfo.dateid", tEnd.Add(time.Hour*24*(-1))))
			break
		case 2:
			// lastweek
			filter = append(filter, dbox.Gte("dateinfo.dateid", tEnd.Add(time.Hour*24*(-7))))
			break
		case 1:
			// mtd
			filter = append(filter, dbox.Gte("dateinfo.dateid", time.Date(tEnd.Year(), tEnd.Month(), 1, 0, 0, 0, 0, tEnd.Location())))
			break
		case 0:
			// ytd
			filter = append(filter, dbox.Gte("dateinfo.dateid", time.Date(tEnd.Year(), 1, 1, 0, 0, 0, 0, tEnd.Location())))
			break
		default:
			// period
			filter = append(filter, dbox.Gte("dateinfo.dateid", tStart))
			break
		}

		filter = append(filter, dbox.Lte("dateinfo.dateid", tEnd))

		if project != "" {
			filter = append(filter, dbox.Eq("projectname", project))
		}

		if len(turbine) != 0 {
			filter = append(filter, dbox.In("turbine", turbine...))
		}

		queryAggr := DB().Connection.NewQuery().From(new(ScadaData).TableName()).
			Aggr(dbox.AggrAvr, "$energy", "totalProduction").
			Aggr(dbox.AggrAvr, "$power", "totalPower").
			Aggr(dbox.AggrAvr, "$denpower", "totaldenPower").
			Group("projectname").Where(dbox.And(filter...))

		caggr, e := queryAggr.Cursor(nil)
		if e != nil {
			return helper.CreateResult(false, nil, e.Error())
		}
		defer caggr.Close()
		e = caggr.Fetch(&aggrData, 0, false)
		if e != nil {
			return helper.CreateResult(false, nil, e.Error())
		}

		var resulttemp Performance

		for _, val := range aggrData {
			switch i {
			case 3:
				// last24hours
				results[len(results)-1].PerformanceIndexLast24Hours = val.GetFloat64("totalPower") / val.GetFloat64("totaldenPower") * 100
				results[len(results)-1].ProductionLast24Hours = val.GetFloat64("totalProduction")
				results[len(results)-1].PowerLast24Hours = val.GetFloat64("totalPower")
				break
			case 2:
				// lastweek
				results[len(results)-1].PerformanceIndexLastWeek = val.GetFloat64("totalPower") / val.GetFloat64("totaldenPower") * 100
				results[len(results)-1].ProductionLastWeek = val.GetFloat64("totalProduction")
				results[len(results)-1].PowerLastWeek = val.GetFloat64("totalPower")
				break
			case 1:
				// mtd
				results[len(results)-1].PerformanceIndexMTD = val.GetFloat64("totalPower") / val.GetFloat64("totaldenPower") * 100
				results[len(results)-1].ProductionMTD = val.GetFloat64("totalProduction")
				results[len(results)-1].PowerMTD = val.GetFloat64("totalPower")
				break
			case 0:
				// ytd
				results = append(results, resulttemp)
				results[len(results)-1].Project = val["_id"].(tk.M)["projectname"].(string)
				results[len(results)-1].PerformanceIndexYTD = val.GetFloat64("totalPower") / val.GetFloat64("totaldenPower") * 100
				results[len(results)-1].ProductionYTD = val.GetFloat64("totalProduction")
				results[len(results)-1].PowerYTD = val.GetFloat64("totalPower")
				results[len(results)-1].StartDate = tStart
				results[len(results)-1].EndDate = tEnd
				break
			default:
				// period
				results[len(results)-1].PerformanceIndex = val.GetFloat64("totalPower") / val.GetFloat64("totaldenPower") * 100
				results[len(results)-1].Production = val.GetFloat64("totalProduction")
				results[len(results)-1].Power = val.GetFloat64("totalPower")
				break
			}
		}
	}

	if len(turbine) == 0 {
		var filter []*dbox.Filter
		filter = append(filter, dbox.Gte("dateinfo.dateid", tStart))
		filter = append(filter, dbox.Lte("dateinfo.dateid", tEnd))

		queryAggr := DB().Connection.NewQuery().From(new(ScadaData).TableName()).
			// Aggr(dbox.AggrMax, "$energy", "energy").
			Group("turbine").Where(dbox.And(filter...))

		caggr, e := queryAggr.Cursor(nil)
		if e != nil {
			return helper.CreateResult(false, nil, e.Error())
		}
		defer caggr.Close()
		e = caggr.Fetch(&aggrData, 0, false)
		if e != nil {
			return helper.CreateResult(false, nil, e.Error())
		}

		for _, val := range aggrData {
			turbine = append(turbine, val["_id"].(tk.M)["turbine"].(string))
		}
	}

	// detail turbines
	for i := 0; i < len(results); i++ {
		for _, valturbine := range turbine {
			for j := 0; j < 5; j++ {
				var filter []*dbox.Filter

				switch j {
				case 3:
					// last24hours
					filter = append(filter, dbox.Gte("dateinfo.dateid", tEnd.Add(time.Hour*24*(-1))))
					break
				case 2:
					// lastweek
					filter = append(filter, dbox.Gte("dateinfo.dateid", tEnd.Add(time.Hour*24*(-7))))
					break
				case 1:
					// mtd
					filter = append(filter, dbox.Gte("dateinfo.dateid", time.Date(tEnd.Year(), tEnd.Month(), 1, 0, 0, 0, 0, tEnd.Location())))
					break
				case 0:
					// ytd
					filter = append(filter, dbox.Gte("dateinfo.dateid", time.Date(tEnd.Year(), 1, 1, 0, 0, 0, 0, tEnd.Location())))
					break
				default:
					// period
					filter = append(filter, dbox.Gte("dateinfo.dateid", tStart))
					break
				}

				filter = append(filter, dbox.Lte("dateinfo.dateid", tEnd))

				// if project != "" {
				filter = append(filter, dbox.Eq("projectname", results[i].Project))
				// }

				// if len(turbine) != 0 {
				filter = append(filter, dbox.In("turbine", valturbine))
				// }

				queryAggr := DB().Connection.NewQuery().From(new(ScadaData).TableName()).
					Aggr(dbox.AggrAvr, "$energy", "totalProduction").
					Aggr(dbox.AggrAvr, "$power", "totalPower").
					Aggr(dbox.AggrAvr, "$denpower", "totaldenPower").
					Group("turbine").Where(dbox.And(filter...))

				caggr, e := queryAggr.Cursor(nil)
				if e != nil {
					return helper.CreateResult(false, nil, e.Error())
				}
				defer caggr.Close()
				e = caggr.Fetch(&aggrData, 0, false)
				if e != nil {
					return helper.CreateResult(false, nil, e.Error())
				}

				// fmt.Println(aggrData)

				var resultdetailtemp PerformanceDetail

				for _, val := range aggrData {
					switch j {
					case 3:
						// last24hours
						results[i].Details[len(results[i].Details)-1].PerformanceIndexLast24Hours = val.GetFloat64("totalPower") / val.GetFloat64("totaldenPower") * 100
						results[i].Details[len(results[i].Details)-1].ProductionLast24Hours = val.GetFloat64("totalProduction")
						results[i].Details[len(results[i].Details)-1].PowerLast24Hours = val.GetFloat64("totalPower")
						break
					case 2:
						// lastweek
						results[i].Details[len(results[i].Details)-1].PerformanceIndexLastWeek = val.GetFloat64("totalPower") / val.GetFloat64("totaldenPower") * 100
						results[i].Details[len(results[i].Details)-1].ProductionLastWeek = val.GetFloat64("totalProduction")
						results[i].Details[len(results[i].Details)-1].PowerLastWeek = val.GetFloat64("totalPower")
						break
					case 1:
						// mtd
						results[i].Details[len(results[i].Details)-1].PerformanceIndexMTD = val.GetFloat64("totalPower") / val.GetFloat64("totaldenPower") * 100
						results[i].Details[len(results[i].Details)-1].ProductionMTD = val.GetFloat64("totalProduction")
						results[i].Details[len(results[i].Details)-1].PowerMTD = val.GetFloat64("totalPower")
						break
					case 0:
						// ytd
						results[i].Details = append(results[i].Details, resultdetailtemp)
						results[i].Details[len(results[i].Details)-1].Turbine = val["_id"].(tk.M)["turbine"].(string)
						results[i].Details[len(results[i].Details)-1].PerformanceIndexYTD = val.GetFloat64("totalPower") / val.GetFloat64("totaldenPower") * 100
						results[i].Details[len(results[i].Details)-1].ProductionYTD = val.GetFloat64("totalProduction")
						results[i].Details[len(results[i].Details)-1].PowerYTD = val.GetFloat64("totalPower")
						results[i].Details[len(results[i].Details)-1].StartDate = tStart
						results[i].Details[len(results[i].Details)-1].EndDate = tEnd
						break
					default:
						// period
						results[i].Details[len(results[i].Details)-1].PerformanceIndex = val.GetFloat64("totalPower") / val.GetFloat64("totaldenPower") * 100
						results[i].Details[len(results[i].Details)-1].Production = val.GetFloat64("totalProduction")
						results[i].Details[len(results[i].Details)-1].Power = val.GetFloat64("totalPower")
						break
					}
				}
			}
		}
	}

	data := struct {
		Data []Performance
	}{
		Data: results,
	}

	return helper.CreateResult(true, data, "success")
}
