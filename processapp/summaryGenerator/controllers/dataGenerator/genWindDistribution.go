package generatorControllers

import (
	. "eaciit/wfdemo-git/library/models"
	. "eaciit/wfdemo-git/processapp/summaryGenerator/controllers"

	"eaciit/wfdemo-git/web/helper"

	"os"
	"time"

	"github.com/eaciit/crowd"
	"github.com/eaciit/dbox"
	_ "github.com/eaciit/dbox/dbc/mongo"
	tk "github.com/eaciit/toolkit"
)

type GenDataWindDistribution struct {
	*BaseController
}

var windCats = [...]float64{1, 1.5, 2, 2.5, 3, 3.5, 4, 4.5, 5, 5.5, 6, 6.5, 7, 7.5, 8, 8.5, 9, 9.5, 10, 10.5, 11, 11.5, 12, 12.5, 13, 13.5, 14, 14.5, 15}

func (d *GenDataWindDistribution) GenerateCurrentMonth(base *BaseController) {
	d.BaseController = base

	type ScadaAnalyticsWDDataX struct {
		Project  string
		Category float64
		Minutes  float64
	}

	type ScadaAnalyticsWDDataGroup struct {
		Project  string
		Category float64
	}

	type MiniScada struct {
		AvgWindSpeed float64
		Project      string
		Count        int
	}

	conn, e := PrepareConnection()
	if e != nil {
		d.Log.AddLog(tk.Sprintf("Wind Distribution : %s"+e.Error()), sWarning)
		os.Exit(0)
	}
	defer conn.Close()

	projects, _ := helper.GetProjectList()

	mdl := new(LatestDataPeriod)
	csr, e := conn.NewQuery().
		Select().
		From(mdl.TableName()).
		Where(dbox.Eq("type", "ScadaData")).
		Cursor(nil)

	if e != nil {
		d.Log.AddLog(tk.Sprintf("Wind Distribution : %s"+e.Error()), sWarning)
	}
	defer csr.Close()

	latesttime := time.Time{}
	for {
		mdl = new(LatestDataPeriod)
		e = csr.Fetch(mdl, 1, false)
		if e != nil {
			break
		}

		if latesttime.IsZero() || latesttime.UTC().Before(mdl.Data[1].UTC()) {
			latesttime = mdl.Data[1].UTC()
		}
	}

	stime := time.Date(latesttime.Year(), latesttime.Month(), 1, 0, 0, 0, 0, latesttime.Location())

	query, pipes := []tk.M{}, []tk.M{}
	query = append(query, tk.M{"_id": tk.M{"$ne": ""}})
	query = append(query, tk.M{"dateinfo.dateid": tk.M{"$gte": stime}})
	query = append(query, tk.M{"dateinfo.dateid": tk.M{"$lte": latesttime}})
	query = append(query, tk.M{"avgwindspeed": tk.M{"$gte": 0.5}})
	query = append(query, tk.M{"available": tk.M{"$eq": 1}})

	iquery := []tk.M{}
	qSave := conn.NewQuery().
		From("rpt_winddistributioncurrentmonth").
		SetConfig("multiexec", true).
		Save()

	defer qSave.Close()

	for _, oproject := range projects {
		proj := oproject.Value

		// ==========================================
		_ = conn.NewQuery().
			Delete().
			From("rpt_winddistributioncurrentmonth").
			Where(dbox.Eq("Project", proj)).
			SetConfig("multiexec", true).
			Exec(nil)
		// ==========================================

		_data := []tk.M{}
		pipes = []tk.M{}

		tmpResult := []MiniScada{}

		iquery = query
		iquery = append(iquery, tk.M{"projectname": proj})
		pipes = append(pipes, tk.M{"$match": tk.M{"$and": iquery}})
		pipes = append(pipes, tk.M{"$group": tk.M{"_id": tk.M{"projectname": "$projectname", "avgwindspeed": "$avgwindspeed"}, "count": tk.M{"$sum": 1}}})
		pipes = append(pipes, tk.M{"$project": tk.M{"_id.projectname": 1, "_id.avgwindspeed": 1, "count": 1}})

		csrx, _ := conn.NewQuery().
			From(new(ScadaData).TableName()).
			Command("pipe", pipes).Cursor(nil)

		e = csrx.Fetch(&_data, 0, false)
		if e != nil {
			break
		}
		csrx.Close()

		for _, v := range _data {
			id := v.Get("_id").(tk.M)
			tmpResult = append(tmpResult, MiniScada{
				AvgWindSpeed: id.GetFloat64("avgwindspeed"),
				Project:      id.GetString("projectname"),
				Count:        v.GetInt("count"),
			})
		}

		if len(tmpResult) > 0 {
			totalCount := 0
			datas := crowd.From(&tmpResult).Apply(func(x interface{}) interface{} {
				dt := x.(MiniScada)

				var di ScadaAnalyticsWDDataX
				di.Project = dt.Project
				di.Category = getWindDistrCategory(dt.AvgWindSpeed)
				di.Minutes = float64(10 * dt.Count)
				totalCount += dt.Count

				return di
			}).Exec().Group(func(x interface{}) interface{} {
				dt := x.(ScadaAnalyticsWDDataX)

				var dig ScadaAnalyticsWDDataGroup
				dig.Project = dt.Project
				dig.Category = dt.Category

				return dig
			}, nil).Exec()

			dts := datas.Apply(func(x interface{}) interface{} {
				kv := x.(crowd.KV)
				keys := kv.Key.(ScadaAnalyticsWDDataGroup)
				vs := kv.Value.([]ScadaAnalyticsWDDataX)
				total := 0.0

				for _, v := range vs {
					total += v.Minutes
				}

				var di ScadaAnalyticsWDDataX
				di.Project = keys.Project
				di.Category = keys.Category
				di.Minutes = total

				return di
			}).Exec().Result.Data().([]ScadaAnalyticsWDDataX)

			totalMinutes := float64(totalCount * 10)

			for _, wc := range windCats {
				exist := crowd.From(&dts).Where(func(x interface{}) interface{} {
					y := x.(ScadaAnalyticsWDDataX)
					Project := y.Project == proj
					Category := y.Category == wc
					return Project && Category
				}).Exec().Result.Data().([]ScadaAnalyticsWDDataX)

				distHelper := tk.M{}

				if len(exist) > 0 {
					distHelper.Set("Project", proj)
					distHelper.Set("Category", wc)

					Minute := crowd.From(&exist).Sum(func(x interface{}) interface{} {
						dt := x.(ScadaAnalyticsWDDataX)
						return dt.Minutes
					}).Exec().Result.Sum

					distHelper.Set("Contribute", Minute/totalMinutes)
				} else {
					distHelper.Set("Project", proj)
					distHelper.Set("Category", wc)
					distHelper.Set("Contribute", -0.0)
				}

				distHelper.Set("_id", tk.Sprintf("%s_%v", proj, wc))
				_ = qSave.Exec(tk.M{}.Set("data", distHelper))
			}
		}
	}

}

func getWindDistrCategory(windValue float64) float64 {
	var datas float64

	for _, val := range windCats {
		if val >= windValue {
			datas = val
			return datas
		}
	}

	return datas
}