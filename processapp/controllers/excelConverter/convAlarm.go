package converterControllers

import (
	. "eaciit/wfdemo-git-dev/library/helper"
	. "eaciit/wfdemo-git-dev/library/models"
	. "eaciit/wfdemo-git-dev/processapp/controllers"
	"os"
	"strings"

	_ "github.com/eaciit/dbox/dbc/mongo"
	tk "github.com/eaciit/toolkit"
	"github.com/tealeg/xlsx"
)

// ConvAlarm
type ConvAlarm struct {
	*BaseController
}

// Generate
func (d *ConvAlarm) Generate(base *BaseController) {
	folderName := "alarm"
	funcName := "Converting Alarm Data"
	count := 0
	total := 0
	if base != nil {
		d.BaseController = base

		ctx := d.BaseController.Ctx
		_ = ctx
		dataSources, path := base.GetDataSource(folderName)
		tk.Println("Converting Alarm Data from Excel File..")
		for _, source := range dataSources {
			if strings.Contains(source.Name(), "Alarm") {
				tk.Println(path + "\\" + source.Name())
				file, e := xlsx.OpenFile(path + "\\" + source.Name())
				if e != nil {
					ErrorHandler(e, funcName)
					os.Exit(0)
				}

				for _, sheet := range file.Sheet {
					errorLine := tk.M{}
					tmpData := new(Alarm).New()
					for idx, row := range sheet.Rows {
						/*farm, e := row.Cells[0].String()
						_ = e*/
						errorList := []error{}
						if idx > 0 { //&& farm != "" && len(row.Cells) == 14 {
							// if idx > 0 && len(row.Cells) == 14 {
							/*tk.Printf("%v ", idx+1)
							for _, val := range row.Cells {
								x, _ := val.String()
								tk.Printf("%v ", x)
							}
							tk.Println()*/
							x := ""
							data := new(Alarm).New()

							if tmpData.Farm == "" {
								data.Farm, e = row.Cells[0].String()
								ErrorLog(e, funcName, errorList)

								data.StartDate, e = GetDateCellAuto(row.Cells[1], row.Cells[2])
								ErrorLog(e, funcName, errorList)

								if row.Cells[1].GetNumberFormat() != "general" {
									data.StartDate, e = ReverseMonthDate(data.StartDate)
								}

								data.StartDateInfo = GetDateInfo(data.StartDate)

								data.Turbine, e = row.Cells[3].String()
								ErrorLog(e, funcName, errorList)

								data.AlertDescription, e = row.Cells[4].String()
								ErrorLog(e, funcName, errorList)

								data.Line = idx + 1

								x, e = row.Cells[5].String()
								ErrorLog(e, funcName, errorList)
							}

							if x != "" {
								data.ExternalStop = row.Cells[5].Bool()
								data.GridDown = row.Cells[6].Bool()
								data.InternalGrid = row.Cells[7].Bool()
								data.MachineDown = row.Cells[8].Bool()
								data.AEbOK = row.Cells[9].Bool()
								data.Unknown = row.Cells[10].Bool()
								data.WeatherStop = row.Cells[11].Bool()
								data.EndDate, e = GetDateCellAuto(row.Cells[12], row.Cells[13])
								ErrorLog(e, funcName, errorList)
								if row.Cells[12].GetNumberFormat() != "general" {
									data.EndDate, e = ReverseMonthDate(data.EndDate)
								}
								tmpData = new(Alarm).New()
							} else if tmpData.Farm != "" {
								data = tmpData
								data.ExternalStop = row.Cells[1].Bool()
								data.GridDown = row.Cells[2].Bool()
								data.InternalGrid = row.Cells[3].Bool()
								data.MachineDown = row.Cells[4].Bool()
								data.AEbOK = row.Cells[5].Bool()
								data.Unknown = row.Cells[6].Bool()
								data.WeatherStop = row.Cells[7].Bool()
								data.EndDate, e = GetDateCellAuto(row.Cells[8], row.Cells[9])
								ErrorLog(e, funcName, errorList)
								if row.Cells[8].GetNumberFormat() != "general" {
									data.EndDate, e = ReverseMonthDate(data.EndDate)
								}
								errorLine.Set(tk.ToString(idx+1), errorList)
								// tk.Printf("up up: %v: %v - %v \n", idx+1, data.Farm, data.ExternalStop)
								tmpData = new(Alarm).New()
							} else {
								// tk.Printf("before: %v: %v - %v \n", idx+1, tmpData.Farm, data.ExternalStop)
								tmpData = data
								data = new(Alarm).New()
								errorLine.Set(tk.ToString(idx+1), errorList)
								// tk.Printf("after: %v: %v - %v \n", idx+1, tmpData.Farm, data.ExternalStop)
							}

							if data.Farm != "" {
								if e != nil || (data.ExternalStop == false && data.GridDown == false && data.InternalGrid == false && data.MachineDown == false && data.AEbOK == false && data.Unknown == false && data.WeatherStop == false) {
									errorLine.Set(tk.ToString(idx+1), errorList)
								} else {
									// adding duration for Alarm Data
									duration := data.EndDate.Sub(data.StartDate)
									data.Duration = duration.Hours()

									e = ctx.Insert(data)
									count++
									if count == 1000 {
										total += count
										tk.Printf("count: %v \n", total)
										count = 0
									}
								}
							}
						} else {
							if idx != 0 {
								errorLine.Set(tk.ToString(idx+1), errorList)
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
}
