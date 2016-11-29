package generatorControllers

import (
	. "eaciit/wfdemo-git/library/helper"
	. "eaciit/wfdemo-git/library/models"
	. "eaciit/wfdemo-git/processapp/controllers"
	"os"
	"strconv"
	"time"

	_ "github.com/eaciit/dbox/dbc/mongo"
	tk "github.com/eaciit/toolkit"
)

type GenAlarmSummary struct {
	*BaseController
}

func (d *GenAlarmSummary) Generate(base *BaseController) {
	if base != nil {
		d.BaseController = base

		ctx, e := PrepareConnection()
		if e != nil {
			ErrorHandler(e, "Alarm Summary")
			os.Exit(0)
		}

		downCause := tk.M{}
		downCause.Set("externalstop", "External Stop")
		downCause.Set("griddown", "Grid Down")
		downCause.Set("internalgrid", "Internal Grid")
		downCause.Set("machinedown", "Machine Down")
		downCause.Set("aebok", "AEBOK")
		downCause.Set("unknown", "Unknown")
		downCause.Set("weatherstop", "Weather Stop")

		tk.Println("Generate Alarm Summary By Month")

		for field, title := range downCause {

			var pipes []tk.M

			pipes = append(pipes, tk.M{"$match": tk.M{field: true}})
			pipes = append(pipes,
				tk.M{
					"$group": tk.M{"_id": tk.M{"id1": "$startdateinfo.monthid", "id2": "$startdateinfo.monthdesc", "id3": title.(string)},
						"result": tk.M{"$sum": "$powerlost"},
					},
				},
			)

			/*csr, e := ctx.NewQuery().
			From(new(Alarm).TableName()).
			Command("pipe", pipes).
			Cursor(nil)*/

			csr, e := ctx.NewQuery().
				From(new(AlarmClean).TableName()).
				Command("pipe", pipes).
				Cursor(nil)

			ErrorHandler(e, "Generate Alarm Summary")

			result := []tk.M{}
			e = csr.Fetch(&result, 0, false)

			ErrorHandler(e, "Generate Alarm Summary")

			for _, val := range result {

				id := val.Get("_id").(tk.M)

				// tk.Printf("%v \n", id)

				monthid := strconv.Itoa(id.GetInt("id1"))
				year := monthid[0:4]
				month := monthid[4:6]
				day := "01"

				iMonth, _ := strconv.Atoi(string(month))
				iMonth = iMonth - 1

				dtStr := year + "-" + month + "-" + day
				dtId, _ := time.Parse("2006-01-02", dtStr)
				dtinfo := GetDateInfo(dtId)

				mdl := new(AlarmSummaryByMonth)
				mdl.ProjectName = "Tejuva"
				mdl.DateInfo = dtinfo
				mdl.LostEnergy = val.GetFloat64("result")
				mdl.Type = title.(string)
				mdl.ID = mdl.ProjectName + "-" + tk.ToString(dtinfo.MonthId) + "-" + field

				if mdl != nil {
					d.BaseController.Ctx.Insert(mdl)
				}

				/*mdl = new(AlarmSummaryByMonth)
				mdl.ProjectName = "Fleet"
				mdl.DateInfo = dtinfo
				mdl.LostEnergy = val.GetFloat64("result")
				mdl.Type = title.(string)
				mdl.ID = mdl.ProjectName + "-" + tk.ToString(dtinfo.MonthId) + "-" + field

				if mdl != nil {
					d.BaseController.Ctx.Insert(mdl)
				}*/

			}
		}

	}
}
