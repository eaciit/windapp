package converterControllers

import (
	. "eaciit/wfdemo-git/library/helper"
	. "eaciit/wfdemo-git/library/models"
	. "eaciit/wfdemo-git/processapp/controllers"
	"fmt"
	"os"
	"strconv"
	"strings"

	_ "github.com/eaciit/dbox/dbc/mongo"
	"github.com/tealeg/xlsx"
)

type ConvTurbine struct {
	*BaseController
}

func (d *ConvTurbine) Generate(base *BaseController) {
	folderName := "turbine"
	// funcName := "Converting Turbine Data"
	if base != nil {
		d.BaseController = base
		ctx := d.BaseController.Ctx
		dataSource, path := base.GetDataSource(folderName) //base.GetDataSourceDirect("D:\\Works\\2016\\Ostro Wind Farm\\Documents\\Ostro Files\\turbine")
		for _, src := range dataSource {
			if strings.Contains(src.Name(), ".xlsx") {
				xlFile, e := xlsx.OpenFile(path + "\\" + src.Name())

				if e != nil {
					ErrorHandler(e, "Importing Data Turbine")
					os.Exit(0)
				}
				fmt.Println(path + "\\" + src.Name())
				for _, sheet := range xlFile.Sheet {
					for idx, row := range sheet.Rows {
						if idx > 1 {
							mdl := new(TurbineMaster).New()

							turbineValue, _ := row.Cells[2].String()
							mdl.TurbineId = strings.Replace(turbineValue, "-", "", 1)
							mdl.TurbineName = turbineValue
							mdl.Project = "Tejuva"
							strLat := row.Cells[6].Value
							floLat, _ := strconv.ParseFloat(strLat, 64)
							strLon := row.Cells[7].Value
							floLon, _ := strconv.ParseFloat(strLon, 64)
							mdl.Latitude = RoundUp(floLat, .5, 16)
							mdl.Longitude = RoundUp(floLon, .5, 16)

							if mdl != nil {
								ctx.Insert(mdl)
							}
						}
					}
				}
			}
		}
	}
}
