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

/*func (d *ConvPermanentMetTowerNewFormat) Validate(base *BaseController) {
	list := []tk.M{}
	d.BaseController = base
	ctx := d.BaseController.Ctx
	_ = ctx

	csr, e := ctx.Connection.NewQuery().
		From(new(MetTower).TableName()).
		Order("line").
		Cursor(nil)

	if e != nil {
		log.Printf("err: %v \n", e.Error())
	}

	e = csr.Fetch(&list, 0, false)
	csr.Close()

	mp := map[int]tk.M{}

	for _, val := range list {
		mp[val.GetInt("line")] = val
	}

	for i := 1; i < 4322; i++ {
		x := mp[i]

		if x == nil {
			log.Printf("%#v \n", i)
		} else {
		}
	}

}*/

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
							if len(row.Cells) > 0 {
								data := new(MetTower).New()
								data.Line = (idx + 1)

								tmp, e := row.Cells[0].Float()
								ErrorLog(e, funcName, errorList)
								row.Cells[0].SetDateTimeWithFormat(tmp, time.UnixDate)
								strDate, e := row.Cells[0].FormattedValue()
								ErrorLog(e, funcName, errorList)

								dtSplit := strings.Split(strings.Replace(strDate, "  ", " ", 1), " ")

								date, e := time.Parse("Jan 2 15:04:05 2006", dtSplit[1]+" "+dtSplit[2]+" "+dtSplit[3]+" "+dtSplit[5])

								// tk.Printf("%v | %v | %v \n", strDate, date.String(), e.Error())
								if e != nil {
									tk.Printf("%v | %v | %v \n", idx+1, e.Error(), strDate)
								}

								data.TimeStamp = date.UTC()
								data.DateInfo = GetDateInfo(date.UTC())

								/*for _, valx := range row.Cells {
									y, _ := valx.Float()
									log.Printf("cell: %#v \n", y)
								}*/

								if len(row.Cells) > 1 {
									VHubWS90mAvg, err := row.Cells[1].Float()
									if err == nil {
										data.VHubWS90mAvg = VHubWS90mAvg
									}

									ErrorLog(e, funcName, errorList)

									VHubWS90mStdDev, err := row.Cells[2].Float()
									if err == nil {
										data.VHubWS90mStdDev = VHubWS90mStdDev
									}
									ErrorLog(e, funcName, errorList)

									VRefWS88mAvg, err := row.Cells[3].Float()
									if err == nil {
										data.VRefWS88mAvg = VRefWS88mAvg
									}
									ErrorLog(e, funcName, errorList)

									VRefWS88mStdDev, err := row.Cells[4].Float()
									if err == nil {
										data.VRefWS88mStdDev = VRefWS88mStdDev
									}
									ErrorLog(e, funcName, errorList)

									VTipWS42mAvg, err := row.Cells[5].Float()
									if err == nil {
										data.VTipWS42mAvg = VTipWS42mAvg
									}
									ErrorLog(e, funcName, errorList)

									VTipWS42mStdDev, err := row.Cells[6].Float()
									if err == nil {
										data.VTipWS42mStdDev = VTipWS42mStdDev
									}
									ErrorLog(e, funcName, errorList)

									DHubWD88mAvg, err := row.Cells[7].Float()
									if err == nil {
										data.DHubWD88mAvg = DHubWD88mAvg
									}
									ErrorLog(e, funcName, errorList)

									DRefWD86mAvg, err := row.Cells[8].Float()
									if err == nil {
										data.DRefWD86mAvg = DRefWD86mAvg
									}
									ErrorLog(e, funcName, errorList)

									TRefHRefHumid855mAvg, err := row.Cells[9].Float()
									if err == nil {
										data.TRefHRefHumid855mAvg = TRefHRefHumid855mAvg
									}
									ErrorLog(e, funcName, errorList)

									TRefHRefHumid855mStdDev, err := row.Cells[10].Float()
									if err == nil {
										data.TRefHRefHumid855mStdDev = TRefHRefHumid855mStdDev
									}
									ErrorLog(e, funcName, errorList)

									TRefHRefHumid855mMax, err := row.Cells[11].Float()
									if err == nil {
										data.TRefHRefHumid855mMax = TRefHRefHumid855mMax
									}
									ErrorLog(e, funcName, errorList)

									TRefHRefHumid855mMin, err := row.Cells[12].Float()
									if err == nil {
										data.TRefHRefHumid855mMin = TRefHRefHumid855mMin
									}
									ErrorLog(e, funcName, errorList)

									BaroAirPress855mAvg, err := row.Cells[13].Float()
									if err == nil {
										data.BaroAirPress855mAvg = BaroAirPress855mAvg
									}
									ErrorLog(e, funcName, errorList)
								}

								if e != nil {
									errorLine.Set(tk.ToString(idx+1), errorList)
								} else {
									e = ctx.Insert(data)
									ErrorLog(e, funcName, errorList)
								}
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
