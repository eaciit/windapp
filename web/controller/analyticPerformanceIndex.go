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

	type Performance struct {
		Project                     string
		Turbine                     string
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
		ProductionIndex             float64
		ProductionIndexLast24Hours  float64
		ProductionIndexLastWeek     float64
		ProductionIndexMTD          float64
		ProductionIndexYTD          float64
		StartDate                   time.Time
		EndDate                     time.Time
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

	// results := make([]Performance, 0)
	results := map[string][]Performance{}

	aggrData := []tk.M{}

	query := []tk.M{}
	pipes := []tk.M{}
	group := tk.M{}
	var resulttemp Performance
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
			query = append(query, tk.M{"turbine": tk.M{"$in": turbine}})
		}

		query = append(query, tk.M{"available": 1})

		pipes = append(pipes, tk.M{"$match": tk.M{"$and": query}})
		group = tk.M{
			"_id": tk.M{
				"projectname": "$projectname",
				"turbine":     "$turbine",
			},
			"totaldenPower": tk.M{"$sum": "$denpower"},
			"totalPower":    tk.M{"$sum": "$power"},
			"totalOKTime":   tk.M{"$sum": "$oktime"},
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
				results["Summary"][len(results["Summary"])-1].PotentialPowerLast24Hours += val.GetFloat64("totaldenPower") / 1000
				results["Summary"][len(results["Summary"])-1].PowerLast24Hours += val.GetFloat64("totalPower") / 1000
				results["Summary"][len(results["Summary"])-1].PerformanceIndexLast24Hours = tk.Div(results["Summary"][len(results["Summary"])-1].PowerLast24Hours, results["Summary"][len(results["Summary"])-1].PotentialPowerLast24Hours) * 100
				results["Summary"][len(results["Summary"])-1].ProductionIndexLast24Hours = tk.Div(tk.Div(val.GetFloat64("totalPower"), 1000), tk.Div(val.GetFloat64("totalOKTime"), 3600))
				results["Data"][indexTurbine[turbinename]].PerformanceIndexLast24Hours = tk.Div(val.GetFloat64("totalPower")/1000, val.GetFloat64("totaldenPower")/1000) * 100
				results["Data"][indexTurbine[turbinename]].PotentialPowerLast24Hours = val.GetFloat64("totaldenPower") / 1000
				results["Data"][indexTurbine[turbinename]].PowerLast24Hours = val.GetFloat64("totalPower") / 1000
				results["Data"][indexTurbine[turbinename]].ProductionIndexLast24Hours = tk.Div(tk.Div(val.GetFloat64("totalPower"), 1000), tk.Div(val.GetFloat64("totalOKTime"), 3600))
				break
			case 2:
				// lastweek
				results["Summary"][len(results["Summary"])-1].PotentialPowerLastWeek += val.GetFloat64("totaldenPower") / 1000
				results["Summary"][len(results["Summary"])-1].PowerLastWeek += val.GetFloat64("totalPower") / 1000
				results["Summary"][len(results["Summary"])-1].PerformanceIndexLastWeek = tk.Div(results["Summary"][len(results["Summary"])-1].PowerLastWeek, results["Summary"][len(results["Summary"])-1].PotentialPowerLastWeek) * 100
				results["Summary"][len(results["Summary"])-1].ProductionIndexLastWeek = tk.Div(tk.Div(val.GetFloat64("totalPower"), 1000), tk.Div(val.GetFloat64("totalOKTime"), 3600))
				results["Data"][indexTurbine[turbinename]].PerformanceIndexLastWeek = tk.Div(val.GetFloat64("totalPower")/1000, val.GetFloat64("totaldenPower")/1000) * 100
				results["Data"][indexTurbine[turbinename]].PotentialPowerLastWeek = val.GetFloat64("totaldenPower") / 1000
				results["Data"][indexTurbine[turbinename]].PowerLastWeek = val.GetFloat64("totalPower") / 1000
				results["Data"][indexTurbine[turbinename]].ProductionIndexLastWeek = tk.Div(tk.Div(val.GetFloat64("totalPower"), 1000), tk.Div(val.GetFloat64("totalOKTime"), 3600))
				break
			case 1:
				// mtd
				results["Summary"][len(results["Summary"])-1].PotentialPowerMTD += val.GetFloat64("totaldenPower") / 1000
				results["Summary"][len(results["Summary"])-1].PowerMTD += val.GetFloat64("totalPower") / 1000
				results["Summary"][len(results["Summary"])-1].PerformanceIndexMTD = tk.Div(results["Summary"][len(results["Summary"])-1].PowerMTD, results["Summary"][len(results["Summary"])-1].PotentialPowerMTD) * 100
				results["Summary"][len(results["Summary"])-1].ProductionIndexMTD = tk.Div(tk.Div(val.GetFloat64("totalPower"), 1000), tk.Div(val.GetFloat64("totalOKTime"), 3600))
				results["Data"][indexTurbine[turbinename]].PerformanceIndexMTD = tk.Div(val.GetFloat64("totalPower")/1000, val.GetFloat64("totaldenPower")/1000) * 100
				results["Data"][indexTurbine[turbinename]].PotentialPowerMTD = val.GetFloat64("totaldenPower") / 1000
				results["Data"][indexTurbine[turbinename]].PowerMTD = val.GetFloat64("totalPower") / 1000
				results["Data"][indexTurbine[turbinename]].ProductionIndexMTD = tk.Div(tk.Div(val.GetFloat64("totalPower"), 1000), tk.Div(val.GetFloat64("totalOKTime"), 3600))
				break
			case 0:
				// ytd
				if lastProject != projectname {
					results["Summary"] = append(results["Summary"], resulttemp)
					results["Summary"][len(results["Summary"])-1].Project = projectname
					results["Summary"][len(results["Summary"])-1].StartDate = tStart
					results["Summary"][len(results["Summary"])-1].EndDate = tEnd
				}
				results["Summary"][len(results["Summary"])-1].PotentialPowerYTD += val.GetFloat64("totaldenPower") / 1000
				results["Summary"][len(results["Summary"])-1].PowerYTD += val.GetFloat64("totalPower") / 1000
				results["Summary"][len(results["Summary"])-1].PerformanceIndexYTD = tk.Div(results["Summary"][len(results["Summary"])-1].PowerYTD, results["Summary"][len(results["Summary"])-1].PotentialPowerYTD) * 100
				results["Summary"][len(results["Summary"])-1].ProductionIndexYTD = tk.Div(tk.Div(val.GetFloat64("totalPower"), 1000), tk.Div(val.GetFloat64("totalOKTime"), 3600))

				lastProject = projectname

				results["Data"] = append(results["Data"], resulttemp)
				indexTurbine[turbinename] = idxChild
				results["Data"][indexTurbine[turbinename]].Turbine = turbinename
				results["Data"][indexTurbine[turbinename]].Project = projectname
				results["Data"][indexTurbine[turbinename]].PerformanceIndexYTD = tk.Div(val.GetFloat64("totalPower")/1000, val.GetFloat64("totaldenPower")/1000) * 100
				results["Data"][indexTurbine[turbinename]].PotentialPowerYTD = val.GetFloat64("totaldenPower") / 1000
				results["Data"][indexTurbine[turbinename]].PowerYTD = val.GetFloat64("totalPower") / 1000
				results["Data"][indexTurbine[turbinename]].ProductionIndexYTD = tk.Div(tk.Div(val.GetFloat64("totalPower"), 1000), tk.Div(val.GetFloat64("totalOKTime"), 3600))
				results["Data"][indexTurbine[turbinename]].StartDate = tStart
				results["Data"][indexTurbine[turbinename]].EndDate = tEnd
				break
			default:
				// period
				results["Summary"][len(results["Summary"])-1].PotentialPower += val.GetFloat64("totaldenPower") / 1000
				results["Summary"][len(results["Summary"])-1].Power += val.GetFloat64("totalPower") / 1000
				results["Summary"][len(results["Summary"])-1].PerformanceIndex = tk.Div(results["Summary"][len(results["Summary"])-1].Power, results["Summary"][len(results["Summary"])-1].PotentialPower) * 100
				results["Summary"][len(results["Summary"])-1].ProductionIndex = tk.Div(tk.Div(val.GetFloat64("totalPower"), 1000), tk.Div(val.GetFloat64("totalOKTime"), 3600))
				results["Data"][indexTurbine[turbinename]].PerformanceIndex = tk.Div(val.GetFloat64("totalPower")/1000, val.GetFloat64("totaldenPower")/1000) * 100
				results["Data"][indexTurbine[turbinename]].PotentialPower = val.GetFloat64("totaldenPower") / 1000
				results["Data"][indexTurbine[turbinename]].Power = val.GetFloat64("totalPower") / 1000
				results["Data"][indexTurbine[turbinename]].ProductionIndex = tk.Div(tk.Div(val.GetFloat64("totalPower"), 1000), tk.Div(val.GetFloat64("totalOKTime"), 3600))
				break
			}
		}
	}

	return helper.CreateResult(true, results, "success")
}
