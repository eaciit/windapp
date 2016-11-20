package converterControllers

import (
	. "eaciit/ostrowfm/library/helper"
	. "eaciit/ostrowfm/library/models"
	. "eaciit/ostrowfm/processapp/controllers"
	"os"
	"strings"
	"time"

	_ "github.com/eaciit/dbox/dbc/mongo"
	tk "github.com/eaciit/toolkit"
	"github.com/tealeg/xlsx"
)

// ConvJMRBreakup
type ConvJMRBreakup struct {
	*BaseController
}

// Generate
func (d *ConvJMRBreakup) Generate(base *BaseController) {
	folderName := "jmr"
	funcName := "Converting JMR & Breakup"
	count := 0
	total := 0
	errorLine := tk.M{}
	if base != nil {
		d.BaseController = base

		ctx := d.BaseController.Ctx
		_ = ctx
		dataSources, path := base.GetDataSource(folderName)
		tk.Println("Converting JMR & Breakup from Excel File..")
		for _, source := range dataSources {
			if !strings.Contains(source.Name(), "~") {
				tk.Println(path + "\\" + source.Name())
				file, e := xlsx.OpenFile(path + "\\" + source.Name())
				if e != nil {
					ErrorHandler(e, funcName)
					os.Exit(0)
				}

				for _, sheet := range file.Sheet {
					var e error
					var description string
					var dateStr string
					data := new(JMR).New()

					for idx, row := range sheet.Rows {
						errorList := []error{}
						if idx == 0 {
							description, e = row.Cells[1].String()
							ErrorHandler(e, funcName)

							data.Description = description

							if e != nil {
								errorLine.Set(tk.ToString(idx+1), errorList)
							}
						} else if idx == 1 {
							dateStr, e = row.Cells[1].String()
							ErrorHandler(e, funcName)

							split := strings.Split(dateStr, "-")
							if len(split) == 2 {
								monthStr := split[0]
								yearStr := "20" + split[1]

								date, e := time.Parse("2006 January 02", yearStr+" "+monthStr+" 01")
								ErrorHandler(e, funcName)

								if e != nil {
									errorLine.Set(tk.ToString(idx+1), errorList)
								} else {
									dateInfo := GetDateInfo(date)
									data.DateInfo = dateInfo
								}
							}

							if e != nil {
								errorLine.Set(tk.ToString(idx+1), errorList)
							}
						} else if idx > 1 {
							cell0, e := row.Cells[0].String()
							cell2, e := row.Cells[2].String()

							if cell0 != "" || cell2 != "Total" {
								var section JMRSection

								section.Description, e = row.Cells[0].String()
								ErrorLog(e, funcName, errorList)

								section.Turbine, e = row.Cells[1].String()
								ErrorLog(e, funcName, errorList)

								section.Company, e = row.Cells[2].String()
								ErrorLog(e, funcName, errorList)

								section.ContrGen, e = row.Cells[3].Float()
								ErrorLog(e, funcName, errorList)

								section.BoEExport, e = row.Cells[4].Float()
								ErrorLog(e, funcName, errorList)
								section.BoEImport, e = row.Cells[5].Float()
								ErrorLog(e, funcName, errorList)
								section.BoENet, e = row.Cells[6].Float()
								ErrorLog(e, funcName, errorList)

								section.BoETotalLoss = section.ContrGen - section.BoEExport

								section.BoLExport, e = row.Cells[11].Float()
								ErrorLog(e, funcName, errorList)
								section.BoLImport, e = row.Cells[12].Float()
								ErrorLog(e, funcName, errorList)
								section.BoLNet, e = row.Cells[13].Float()
								ErrorLog(e, funcName, errorList)

								section.BoE2Export, e = row.Cells[15].Float()
								ErrorLog(e, funcName, errorList)
								section.BoE2Import, e = row.Cells[16].Float()
								ErrorLog(e, funcName, errorList)
								section.BoE2Net, e = row.Cells[17].Float()
								ErrorLog(e, funcName, errorList)

								if e != nil {
									errorLine.Set(tk.ToString(idx+1), errorList)
								} else {
									data.Sections = append(data.Sections, section)
								}
							}

						} else {
							if idx != 0 {
								errorLine.Set(tk.ToString(idx+1), errorList)
							}
						}
					}

					data.SetTotalDetails()
					e = ctx.Insert(data)
					ErrorHandler(e, funcName)
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
