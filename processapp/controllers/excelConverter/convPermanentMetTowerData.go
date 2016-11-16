package converterControllers

import (
	. "github.com/eaciit/windapp/library/helper"
	. "github.com/eaciit/windapp/library/models"
	. "github.com/eaciit/windapp/processapp/controllers"
	"os"
	"strings"
	"time"

	_ "github.com/eaciit/dbox/dbc/mongo"
	tk "github.com/eaciit/toolkit"
	"github.com/tealeg/xlsx"
)

// ConvPermanentMetTower
type ConvPermanentMetTower struct {
	*BaseController
}

// Generate
func (d *ConvPermanentMetTower) Generate(base *BaseController) {
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
				tk.Println(path + "\\" + source.Name())
				file, e := xlsx.OpenFile(path + "\\" + source.Name())
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
							data.VHubWS90mMax, e = row.Cells[2].Float()
							ErrorLog(e, funcName, errorList)
							data.VHubWS90mMin, e = row.Cells[3].Float()
							ErrorLog(e, funcName, errorList)
							data.VHubWS90mStdDev, e = row.Cells[4].Float()
							ErrorLog(e, funcName, errorList)
							data.VHubWS90mCount, e = row.Cells[5].Float()
							ErrorLog(e, funcName, errorList)

							data.VRefWS88mAvg, e = row.Cells[6].Float()
							ErrorLog(e, funcName, errorList)
							data.VRefWS88mMax, e = row.Cells[7].Float()
							ErrorLog(e, funcName, errorList)
							data.VRefWS88mMin, e = row.Cells[8].Float()
							ErrorLog(e, funcName, errorList)
							data.VRefWS88mStdDev, e = row.Cells[9].Float()
							ErrorLog(e, funcName, errorList)
							data.VRefWS88mCount, e = row.Cells[10].Float()
							ErrorLog(e, funcName, errorList)

							data.VTipWS42mAvg, e = row.Cells[11].Float()
							ErrorLog(e, funcName, errorList)
							data.VTipWS42mMax, e = row.Cells[12].Float()
							ErrorLog(e, funcName, errorList)
							data.VTipWS42mMin, e = row.Cells[13].Float()
							ErrorLog(e, funcName, errorList)
							data.VTipWS42mStdDev, e = row.Cells[14].Float()
							ErrorLog(e, funcName, errorList)
							data.VTipWS42mCount, e = row.Cells[15].Float()
							ErrorLog(e, funcName, errorList)

							data.DHubWD88mAvg, e = row.Cells[16].Float()
							ErrorLog(e, funcName, errorList)
							data.DHubWD88mMax, e = row.Cells[17].Float()
							ErrorLog(e, funcName, errorList)
							data.DHubWD88mMin, e = row.Cells[18].Float()
							ErrorLog(e, funcName, errorList)
							data.DHubWD88mStdDev, e = row.Cells[19].Float()
							ErrorLog(e, funcName, errorList)
							data.DHubWD88mCount, e = row.Cells[20].Float()
							ErrorLog(e, funcName, errorList)

							data.DRefWD86mAvg, e = row.Cells[21].Float()
							ErrorLog(e, funcName, errorList)
							data.DRefWD86mMax, e = row.Cells[22].Float()
							ErrorLog(e, funcName, errorList)
							data.DRefWD86mMin, e = row.Cells[23].Float()
							ErrorLog(e, funcName, errorList)
							data.DRefWD86mStdDev, e = row.Cells[24].Float()
							ErrorLog(e, funcName, errorList)
							data.DRefWD86mCount, e = row.Cells[25].Float()
							ErrorLog(e, funcName, errorList)

							data.THubHHubHumid855mAvg, e = row.Cells[26].Float()
							ErrorLog(e, funcName, errorList)
							data.THubHHubHumid855mMax, e = row.Cells[27].Float()
							ErrorLog(e, funcName, errorList)
							data.THubHHubHumid855mMin, e = row.Cells[28].Float()
							ErrorLog(e, funcName, errorList)
							data.THubHHubHumid855mStdDev, e = row.Cells[29].Float()
							ErrorLog(e, funcName, errorList)
							data.THubHHubHumid855mCount, e = row.Cells[30].Float()
							ErrorLog(e, funcName, errorList)

							data.TRefHRefHumid855mAvg, e = row.Cells[31].Float()
							ErrorLog(e, funcName, errorList)
							data.TRefHRefHumid855mMax, e = row.Cells[32].Float()
							ErrorLog(e, funcName, errorList)
							data.TRefHRefHumid855mMin, e = row.Cells[33].Float()
							ErrorLog(e, funcName, errorList)
							data.TRefHRefHumid855mStdDev, e = row.Cells[34].Float()
							ErrorLog(e, funcName, errorList)
							data.TRefHRefHumid855mCount, e = row.Cells[35].Float()
							ErrorLog(e, funcName, errorList)

							data.THubHHubTemp855mAvg, e = row.Cells[36].Float()
							ErrorLog(e, funcName, errorList)
							data.THubHHubTemp855mMax, e = row.Cells[37].Float()
							ErrorLog(e, funcName, errorList)
							data.THubHHubTemp855mMin, e = row.Cells[38].Float()
							ErrorLog(e, funcName, errorList)
							data.THubHHubTemp855mStdDev, e = row.Cells[39].Float()
							ErrorLog(e, funcName, errorList)
							data.THubHHubTemp855mCount, e = row.Cells[40].Float()
							ErrorLog(e, funcName, errorList)

							data.TRefHRefTemp855mAvg, e = row.Cells[41].Float()
							ErrorLog(e, funcName, errorList)
							data.TRefHRefTemp855mMax, e = row.Cells[42].Float()
							ErrorLog(e, funcName, errorList)
							data.TRefHRefTemp855mMin, e = row.Cells[43].Float()
							ErrorLog(e, funcName, errorList)
							data.TRefHRefTemp855mStdDev, e = row.Cells[44].Float()
							ErrorLog(e, funcName, errorList)
							data.TRefHRefTemp855mCount, e = row.Cells[45].Float()
							ErrorLog(e, funcName, errorList)

							data.BaroAirPress855mAvg, e = row.Cells[46].Float()
							ErrorLog(e, funcName, errorList)
							data.BaroAirPress855mMax, e = row.Cells[47].Float()
							ErrorLog(e, funcName, errorList)
							data.BaroAirPress855mMin, e = row.Cells[48].Float()
							ErrorLog(e, funcName, errorList)
							data.BaroAirPress855mStdDev, e = row.Cells[49].Float()
							ErrorLog(e, funcName, errorList)
							data.BaroAirPress855mCount, e = row.Cells[50].Float()
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
