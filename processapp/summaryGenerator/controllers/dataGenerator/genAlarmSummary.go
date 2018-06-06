package generatorControllers

import (
	. "eaciit/wfdemo-git/library/helper"
	. "eaciit/wfdemo-git/library/models"
	. "eaciit/wfdemo-git/processapp/summaryGenerator/controllers"
	"eaciit/wfdemo-git/web/helper"
	"strconv"
	"time"

	"github.com/eaciit/dbox"
	_ "github.com/eaciit/dbox/dbc/mongo"
	tk "github.com/eaciit/toolkit"
)

type GenAlarmSummary struct {
	*BaseController
}

func (d *GenAlarmSummary) Generate(base *BaseController) {
	if base != nil {
		d.BaseController = base

		// ctx, e := PrepareConnection()
		// if e != nil {
		// 	ErrorHandler(e, "Alarm Summary")
		// 	os.Exit(0)
		// }
		ctx := d.BaseController.Ctx.Connection

		// #faisal
		// remove delete function
		base.Ctx.DeleteMany(new(AlarmSummaryByMonth), dbox.Ne("projectname", ""))

		downCause := tk.M{}
		downCause.Set("externalstop", "External Stop")
		downCause.Set("griddown", "Grid Down")
		downCause.Set("internalgrid", "Internal Grid")
		downCause.Set("machinedown", "Machine Down")
		downCause.Set("aebok", "AEBOK")
		downCause.Set("unknown", "Unknown")
		downCause.Set("weatherstop", "Weather Stop")

		tk.Println("Generate Alarm Summary By Month")

		projects, _ := helper.GetProjectList()

		for _, proj := range projects {
			projectName := proj.Value
			d.Log.AddLog(tk.Sprintf("> %v \n", projectName), sInfo)
			// turbineList, _ := helper.GetTurbineList([]interface{}{projectName})
			// totalTurbine := len(turbineList)

			for field, title := range downCause {

				var pipes []tk.M
				match := tk.M{}

				match.Set(field, true)
				match.Set("farm", projectName)

				pipes = append(pipes,
					tk.M{"$unwind": "$detail"},
				)
				pipes = append(pipes, tk.M{"$match": match})
				pipes = append(pipes,
					tk.M{
						"$group": tk.M{"_id": tk.M{"id1": "$detail.detaildateinfo.monthid", "id2": "$detail.detaildateinfo.monthdesc", "id3": title.(string)},
							"result": tk.M{"$sum": "$powerlost"},
						},
					},
				)

				// tk.Println("====================pipes===========================")
				// tk.Println(pipes)

				/*csr, e := ctx.NewQuery().
				From(new(Alarm).TableName()).
				Command("pipe", pipes).
				Cursor(nil)*/

				// #faisal
				// add condition to check the latest data, and start the generator from that latest data

				csr, e := ctx.NewQuery().
					From(new(Alarm).TableName()).
					Command("pipe", pipes).
					Cursor(nil)

				defer csr.Close()

				ErrorHandler(e, "Generate Alarm Summary")

				result := []tk.M{}
				e = csr.Fetch(&result, 0, false)

				ErrorHandler(e, "Generate Alarm Summary")

				// tk.Println("====================result===========================")
				// tk.Println(result)

				for _, val := range result {

					id := val.Get("_id").(tk.M)

					// tk.Printf("%v \n", id)

					monthid := tk.ToString(id.GetInt("id1"))
					if monthid != "101" {
						year := monthid[0:4]
						month := monthid[4:6]
						day := "01"

						iMonth, _ := strconv.Atoi(string(month))
						iMonth = iMonth - 1

						dtStr := year + "-" + month + "-" + day
						dtId, _ := time.Parse("2006-01-02", dtStr)
						dtinfo := GetDateInfo(dtId)

						mdl := new(AlarmSummaryByMonth)
						mdl.ProjectName = projectName
						mdl.DateInfo = dtinfo
						// tk.Println(val.GetFloat64("result"))
						mdl.LostEnergy = val.GetFloat64("result")
						mdl.Type = title.(string)
						mdl.ID = mdl.ProjectName + "-" + tk.ToString(dtinfo.MonthId) + "-" + field

						if mdl != nil {
							d.BaseController.Ctx.Insert(mdl)
							// log.Printf(">>> %#v \n", mdl)
						}
					}

				}
			}

		}
	}
}
