package controller

import (
	. "eaciit/wfdemo-git/library/core"
	. "eaciit/wfdemo-git/library/models"
	"eaciit/wfdemo-git/web/helper"

	"time"

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
		PerformanceIndex            float64 //power / den power
		PerformanceIndexLast24Hours float64
		PerformanceIndexLastWeek    float64
		PerformanceIndexMTD         float64
		PerformanceIndexYTD         float64
		PotentialPower              float64 //den power
		PotentialPowerLast24Hours   float64
		PotentialPowerLastWeek      float64
		PotentialPowerMTD           float64
		PotentialPowerYTD           float64
		Power                       float64 // power
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
		PotentialPower              float64
		PotentialPowerLast24Hours   float64
		PotentialPowerLastWeek      float64
		PotentialPowerMTD           float64
		PotentialPowerYTD           float64
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

	tStart, tEnd, e := helper.GetStartEndDate(k, p.Period, p.DateStart, p.DateEnd)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	turbine := p.Turbine
	project := p.Project

	results := make([]Performance, 0)

	aggrData := []tk.M{}

	query := []tk.M{}
	pipes := []tk.M{}
	group := tk.M{}
	var resulttemp Performance
	var resultdetailtemp PerformanceDetail
	lastProject := ""
	indexTurbine := map[string]int{}
	turbinename := ""
	projectname := ""

	for i := 0; i < 5; i++ {
		query = []tk.M{}
		pipes = []tk.M{}

		switch i {
		case 3:
			// last24hours
			query = append(query, tk.M{"dateinfo.dateid": tk.M{"$gte": tEnd.Add(time.Hour * 24 * (-1))}})
			break
		case 2:
			// lastweek
			query = append(query, tk.M{"dateinfo.dateid": tk.M{"$gte": time.Date(tEnd.Year(), tEnd.Month(), tEnd.Day()-7, 0, 0, 0, 0, time.UTC)}})
			break
		case 1:
			// mtd
			query = append(query, tk.M{"dateinfo.dateid": tk.M{"$gte": time.Date(tEnd.Year(), tEnd.Month(), 1, 0, 0, 0, 0, time.UTC)}})
			break
		case 0:
			// ytd
			query = append(query, tk.M{"dateinfo.dateid": tk.M{"$gte": time.Date(tEnd.Year(), 1, 1, 0, 0, 0, 0, time.UTC)}})
			break
		default:
			// period
			query = append(query, tk.M{"dateinfo.dateid": tk.M{"$gte": tStart}})
			break
		}
		query = append(query, tk.M{"dateinfo.dateid": tk.M{"$lte": tEnd}})

		if project != "" {
			query = append(query, tk.M{"projectname": project})
		}

		if len(turbine) != 0 {
			query = append(query, tk.M{"projectname": tk.M{"$in": turbine}})
		}

		pipes = append(pipes, tk.M{"$match": tk.M{"$and": query}})
		group = tk.M{
			"_id": tk.M{
				"projectname": "$projectname",
				"turbine":     "$turbine",
			},
			"totaldenPower": tk.M{"$sum": "$denpower"},
			"totalPower":    tk.M{"$sum": "$power"},
		}
		pipes = append(pipes, tk.M{"$group": group})
		pipes = append(pipes, tk.M{"$sort": tk.M{"_id": 1}})

		caggr, e := DB().Connection.NewQuery().From(new(ScadaData).TableName()).
			Command("pipe", pipes).Cursor(nil)
		if e != nil {
			return helper.CreateResult(false, nil, e.Error())
		}
		e = caggr.Fetch(&aggrData, 0, false)
		if e != nil {
			return helper.CreateResult(false, nil, e.Error())
		}
		caggr.Close()

		for idxChild, val := range aggrData {
			projectname = val["_id"].(tk.M)["projectname"].(string)
			turbinename = val["_id"].(tk.M)["turbine"].(string)
			switch i {
			case 3:
				// last24hours
				results[len(results)-1].PotentialPowerLast24Hours += val.GetFloat64("totaldenPower") / 1000
				results[len(results)-1].PowerLast24Hours += val.GetFloat64("totalPower") / 1000
				results[len(results)-1].PerformanceIndexLast24Hours = tk.Div(results[len(results)-1].PowerLast24Hours, results[len(results)-1].PotentialPowerLast24Hours) * 100
				results[len(results)-1].Details[indexTurbine[turbinename]].PerformanceIndexLast24Hours = tk.Div(val.GetFloat64("totalPower")/1000, val.GetFloat64("totaldenPower")/1000) * 100
				results[len(results)-1].Details[indexTurbine[turbinename]].PotentialPowerLast24Hours = val.GetFloat64("totaldenPower") / 1000
				results[len(results)-1].Details[indexTurbine[turbinename]].PowerLast24Hours = val.GetFloat64("totalPower") / 1000
				break
			case 2:
				// lastweek
				results[len(results)-1].PotentialPowerLastWeek += val.GetFloat64("totaldenPower") / 1000
				results[len(results)-1].PowerLastWeek += val.GetFloat64("totalPower") / 1000
				results[len(results)-1].PerformanceIndexLastWeek = tk.Div(results[len(results)-1].PowerLastWeek, results[len(results)-1].PotentialPowerLastWeek) * 100
				results[len(results)-1].Details[indexTurbine[turbinename]].PerformanceIndexLastWeek = tk.Div(val.GetFloat64("totalPower")/1000, val.GetFloat64("totaldenPower")/1000) * 100
				results[len(results)-1].Details[indexTurbine[turbinename]].PotentialPowerLastWeek = val.GetFloat64("totaldenPower") / 1000
				results[len(results)-1].Details[indexTurbine[turbinename]].PowerLastWeek = val.GetFloat64("totalPower") / 1000
				break
			case 1:
				// mtd
				results[len(results)-1].PotentialPowerMTD += val.GetFloat64("totaldenPower") / 1000
				results[len(results)-1].PowerMTD += val.GetFloat64("totalPower") / 1000
				results[len(results)-1].PerformanceIndexMTD = tk.Div(results[len(results)-1].PowerMTD, results[len(results)-1].PotentialPowerMTD) * 100
				results[len(results)-1].Details[indexTurbine[turbinename]].PerformanceIndexMTD = tk.Div(val.GetFloat64("totalPower")/1000, val.GetFloat64("totaldenPower")/1000) * 100
				results[len(results)-1].Details[indexTurbine[turbinename]].PotentialPowerMTD = val.GetFloat64("totaldenPower") / 1000
				results[len(results)-1].Details[indexTurbine[turbinename]].PowerMTD = val.GetFloat64("totalPower") / 1000
				break
			case 0:
				// ytd
				if lastProject != projectname {
					results = append(results, resulttemp)
					results[len(results)-1].Project = projectname
				}
				lastProject = projectname
				results[len(results)-1].PotentialPowerYTD += val.GetFloat64("totaldenPower") / 1000
				results[len(results)-1].PowerYTD += val.GetFloat64("totalPower") / 1000
				results[len(results)-1].PerformanceIndexYTD = tk.Div(results[len(results)-1].PowerYTD, results[len(results)-1].PotentialPowerYTD) * 100
				results[len(results)-1].StartDate = tStart
				results[len(results)-1].EndDate = tEnd
				results[len(results)-1].Details = append(results[len(results)-1].Details, resultdetailtemp)
				results[len(results)-1].Details[idxChild].Turbine = turbinename
				indexTurbine[turbinename] = idxChild
				results[len(results)-1].Details[indexTurbine[turbinename]].PerformanceIndexYTD = tk.Div(val.GetFloat64("totalPower")/1000, val.GetFloat64("totaldenPower")/1000) * 100
				results[len(results)-1].Details[indexTurbine[turbinename]].PotentialPowerYTD = val.GetFloat64("totaldenPower") / 1000
				results[len(results)-1].Details[indexTurbine[turbinename]].PowerYTD = val.GetFloat64("totalPower") / 1000
				results[len(results)-1].Details[indexTurbine[turbinename]].StartDate = tStart
				results[len(results)-1].Details[indexTurbine[turbinename]].EndDate = tEnd
				break
			default:
				// period
				results[len(results)-1].PotentialPower += val.GetFloat64("totaldenPower") / 1000
				results[len(results)-1].Power += val.GetFloat64("totalPower") / 1000
				results[len(results)-1].PerformanceIndex = tk.Div(results[len(results)-1].Power, results[len(results)-1].PotentialPower) * 100
				results[len(results)-1].Details[indexTurbine[turbinename]].PerformanceIndex = tk.Div(val.GetFloat64("totalPower")/1000, val.GetFloat64("totaldenPower")/1000) * 100
				results[len(results)-1].Details[indexTurbine[turbinename]].PotentialPower = val.GetFloat64("totaldenPower") / 1000
				results[len(results)-1].Details[indexTurbine[turbinename]].Power = val.GetFloat64("totalPower") / 1000
				break
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
