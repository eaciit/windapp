package converterControllers

import (
	. "eaciit/wfdemo-git/library/helper"
	. "eaciit/wfdemo-git/library/models"
	. "eaciit/wfdemo-git/processapp/controllers"
	"os"
	"strings"
	"time"

	_ "github.com/eaciit/dbox/dbc/mongo"
	tk "github.com/eaciit/toolkit"
	"github.com/tealeg/xlsx"
)

// ConvPermanentMetTowerNewFormat
type ConvPermanentMetTowerNewFormat struct {
	*BaseController
}

// Generate
func (d *ConvPermanentMetTowerNewFormat) Generate(base *BaseController) {
	folderName := "met"
	funcName := "Converting Met Tower Data"
	count := 0
	total := 0
	errorLine := tk.M{}
	if base != nil {
		d.BaseController = base
		ctx := d.BaseController.Ctx
		_ = ctx
		dataSources, path := base.GetDataSource(folderName)
		tk.Println("Converting Met Tower Data from Excel File..")
		for _, source := range dataSources {
			if !strings.Contains(source.Name(), "~") {
				tk.Println(path + "/" + source.Name())
				file, e := xlsx.OpenFile(path + "/" + source.Name())
				if e != nil {
					ErrorHandler(e, funcName)
					os.Exit(0)
				}

				for _, sheet := range file.Sheet {
					for idx, row := range sheet.Rows {
						errorList := []error{}
						if idx > 0 {
							data := new(MetTower).New()
							data.Line = (idx + 1)

							tmp, e := row.Cells[0].Float()
							ErrorLog(e, funcName, errorList)
							row.Cells[0].SetDateTimeWithFormat(tmp, time.UnixDate)
							strDate, e := row.Cells[0].FormattedValue()
							ErrorLog(e, funcName, errorList)

							dtSplit := strings.Split(strings.Replace(strDate, "  ", " ", 1), " ")

							date, e := time.Parse("Jan 2 15:04:05 2006", dtSplit[1]+" "+dtSplit[2]+" "+dtSplit[3]+" "+dtSplit[5])

							// tk.Printf("%v | %v \n", strDate, date.String())

							data.TimeStamp = date
							data.DateInfo = GetDateInfo(date)

							data.VHubWS90mAvg, e = row.Cells[1].Float()
							ErrorLog(e, funcName, errorList)
							data.VHubWS90mStdDev, e = row.Cells[2].Float()
							ErrorLog(e, funcName, errorList)

							data.VRefWS88mAvg, e = row.Cells[3].Float()
							ErrorLog(e, funcName, errorList)
							data.VRefWS88mStdDev, e = row.Cells[4].Float()
							ErrorLog(e, funcName, errorList)

							data.VTipWS42mAvg, e = row.Cells[5].Float()
							ErrorLog(e, funcName, errorList)
							data.VTipWS42mStdDev, e = row.Cells[6].Float()
							ErrorLog(e, funcName, errorList)

							data.DHubWD88mAvg, e = row.Cells[7].Float()
							ErrorLog(e, funcName, errorList)

							data.DRefWD86mAvg, e = row.Cells[8].Float()
							ErrorLog(e, funcName, errorList)

							data.TRefHRefHumid855mAvg, e = row.Cells[9].Float()
							ErrorLog(e, funcName, errorList)
							data.TRefHRefHumid855mStdDev, e = row.Cells[10].Float()
							ErrorLog(e, funcName, errorList)
							data.TRefHRefHumid855mMax, e = row.Cells[11].Float()
							ErrorLog(e, funcName, errorList)
							data.TRefHRefHumid855mMin, e = row.Cells[12].Float()
							ErrorLog(e, funcName, errorList)

							data.BaroAirPress855mAvg, e = row.Cells[13].Float()
							ErrorLog(e, funcName, errorList)

							if e != nil {
								errorLine.Set(tk.ToString(idx+1), errorList)
							} else {
								e = ctx.Insert(data)
								ErrorLog(e, funcName, errorList)
							}
						} else {
							if idx != 0 {
								errorLine.Set(tk.ToString(idx+1), errorList)
							}
						}
					}
				}
				total += count
				tk.Printf("count: %v \n", total)
				tk.Printf("count line error: %v \n", len(errorLine))
				if len(errorLine) > 0 {
					WriteErrors(errorLine, source.Name())
				}
			}
		}
	}
}
