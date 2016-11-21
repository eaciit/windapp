package generatorControllers

import (
	. "eaciit/wfdemo-git/library/helper"
	. "eaciit/wfdemo-git/library/models"
	. "eaciit/wfdemo-git/processapp/controllers"
	_ "fmt"
	_ "github.com/eaciit/dbox/dbc/mongo"
	tk "github.com/eaciit/toolkit"
	"os"
	_ "strconv"
	_ "strings"
	"time"
)

type GenScadaPowerCurve struct {
	*BaseController
}

func (d *GenScadaPowerCurve) GenerateAdj(base *BaseController) {
	d.Generate(base, "ADJ", "rpt_scadapowercurve")
}

func (d *GenScadaPowerCurve) GenerateAvg(base *BaseController) {
	d.Generate(base, "AVG", "rpt_scadapowercurve_avg")
}

func (d *GenScadaPowerCurve) Generate(base *BaseController, dataType string, tblName string) {
	if base != nil {
		d.BaseController = base

		conn := d.BaseController.Ctx.Connection
		defer conn.Close()

		csrt, e := conn.NewQuery().From(new(ScadaData).TableName()).
			Group("dateinfo.dateid", "turbine").
			Cursor(nil)
		defer csrt.Close()

		if e != nil {
			ErrorHandler(e, "Scada Power Curve")
			os.Exit(0)
		}

		dataTurbines := []tk.M{}
		e = csrt.Fetch(&dataTurbines, 0, false)

		total := 0
		count := 0
		for _, dataTurbine := range dataTurbines {
			tId := dataTurbine["_id"].(tk.M)
			tDateId := tId["dateinfo_dateid"].(time.Time)
			tTurbine := tId["turbine"].(string)
			tProject := "Tejuva"

			mdl := new(ScadaPowerCurveModel).New()
			mdl.DateInfo = GetDateInfo(tDateId)
			mdl.ProjectName = tProject
			mdl.TurbineId = tTurbine

			if e != nil {
				ErrorHandler(e, "Scada Power Curve Items")
				os.Exit(0)
			}

			fieldWs := "$wsavgforpc"
			if dataType == "ADJ" {
				fieldWs = "$wsadjforpc"
			}

			pipe := []tk.M{tk.M{}.Set("$match", tk.M{}.Set("dateinfo.dateid", tDateId).Set("turbine", tTurbine)), tk.M{}.Set("$group", tk.M{}.Set("_id", fieldWs).Set("power", tk.M{}.Set("$sum", "$power")).Set("total", tk.M{}.Set("$sum", 1)))}
			//tk.Printf("#%v\n", pipe)
			csr1, _ := conn.NewQuery().
				Command("pipe", pipe).
				From(new(ScadaData).TableName()).
				Cursor(nil)

			dtvalues := []tk.M{}
			_ = csr1.Fetch(&dtvalues, 0, false)

			// tk.Printf("#%v\n", dtvalues)

			csr1.Close()

			items := make([]ScadaPowerCurveItem, 0)
			if len(dtvalues) > 0 {
				for _, d := range dtvalues {
					ids := d["_id"]
					vws := 0.0
					if ids != nil {
						vws = ids.(float64)
					}
					vpower := 0.0
					ipower := d["power"]
					if ipower != nil {
						vpower = ipower.(float64)
					}
					total := d["total"].(int)

					var item ScadaPowerCurveItem
					item.WSClass = vws
					item.Production = vpower
					item.TotalData = total

					items = append(items, item)
				}
			}

			mdl.DataItems = items

			mdl.SetTableName(tblName)

			d.BaseController.Ctx.Insert(mdl)

			count++
			total++
			if count == 500 {
				tk.Printf("Total data processed %v \n", total)
				count = 0
			}
		}
		tk.Printf("Total data processed %v \n", total)
	}
}
