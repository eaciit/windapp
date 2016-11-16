package converterControllers

import (
	"bufio"
	. "eaciit/wfdemo/library/helper"
	. "eaciit/wfdemo/library/models"
	. "eaciit/wfdemo/processapp/controllers"
	"encoding/csv"
	"io"
	"os"
	"strings"
	"time"

	"errors"

	_ "github.com/eaciit/dbox/dbc/mongo"
	tk "github.com/eaciit/toolkit"
	"github.com/tealeg/xlsx"
)

// ConvScadaData
type ConvScadaData struct {
	*BaseController
}

var (
	funcName = "Converting Scada Data"
)

// Generate
func (d *ConvScadaData) Generate(base *BaseController) {
	scadaConf := []ScadaConf{}
	ReadJson("conf/genScadaConf.json", &scadaConf)
	project := "Tejuva"
	if len(scadaConf) > 0 {
		// please select the index of the scada kind of
		conf := scadaConf[1]
		if base != nil {
			d.BaseController = base

			ctx := d.BaseController.Ctx
			_ = ctx
			dataSources, path := base.GetDataSource(conf.Folder)
			tk.Println("Converting Scada Data from Excel File..")
			for _, source := range dataSources {
				count := 0
				total := 0
				errorLine := tk.M{}
				if conf.DocType == "excel" {
					if strings.Contains(source.Name(), "Scada") {
						tk.Println(path + "\\" + source.Name())
						file, e := xlsx.OpenFile(path + "\\" + source.Name())
						if e != nil {
							ErrorHandler(e, funcName)
							os.Exit(0)
						}

						for _, sheet := range file.Sheet {
							errorLine = tk.M{}
							for idx, row := range sheet.Rows {
								errorList := []error{}
								if idx > 0 { //&& len(row.Cells) == 35 {
									data, errorList := ConstructScadaDataExcel(conf, row)
									data.Line = idx + 1
									data.ProjectName = project
									/*totalTimeDuration := data.AlarmUnknownTime + data.AlarmWeatherStop + data.ExternalStopTime + data.GridDownTime + data.InternalLineDown + data.MachineDownTime + data.OkTime
									if totalTimeDuration > 600.0 || totalTimeDuration < 600.0 {
										data.IsValidTimeDuration = false
									}*/

									/*tk.Printf("%#v \n", data)
									tk.Println()*/

									if len(errorList) > 0 {
										errorLine.Set(tk.ToString(idx+1), errorList)
									} else {
										e = ctx.Insert(data)
										ErrorHandler(e, "Saving")
										count++
										if count == 1000 {
											total += count
											tk.Printf("count: %v \n", total)
											count = 0
										}
									}
								} else {
									if idx != 0 {
										errorLine.Set(tk.ToString(idx+1), errorList)
									}
								}
							}
						}
					}
				} else if conf.DocType == "csv" {
					if strings.Contains(source.Name(), "DataFile-T") {
						tk.Println(path + "\\" + source.Name())
						fr, _ := os.Open(path + "\\" + source.Name())
						read := csv.NewReader(bufio.NewReader(fr))
						idx := 0
						for {
							record, err := read.Read()
							if err == io.EOF {
								break
							}
							if idx > 0 {
								errorLine = tk.M{}

								data, errorList := ConstructScadaDataCSV(conf, record)
								data.Line = idx + 1
								data.ProjectName = project

								if len(errorList) > 0 {
									errorLine.Set(tk.ToString(idx+1), errorList)
								} else {
									e := ctx.Insert(data)
									ErrorHandler(e, "Saving")
									count++
									if count == 1000 {
										total += count
										tk.Printf("count: %v \n", total)
										count = 0
									}

								}
							}
							idx++
						}
					}
				}

				total += count
				tk.Printf("count: %v \n", total)
				tk.Printf("count line error: %v \n", len(errorLine))
				if len(errorLine) > 0 {
					WriteErrors(errorLine, source.Name())
				}

				// tk.Printf("\n --------- \nTotal Data: %v for: %v \n---------\n", total+count, path+"\\"+source.Name())
			}
		}
	}
}

type ScadaConf struct {
	KindOfDoc string                         `json:"kindOfDoc"`
	Folder    string                         `json:"folder"`
	DocType   string                         `json:"docType"`
	Columns   []map[string][]ScadaConfColumn `json:"columns"`
}

type ScadaConfColumn struct {
	Index int    `json:"index"`
	Type  string `json:"type"`
}

func ConstructScadaDataExcel(conf ScadaConf, row *xlsx.Row) (res *ScadaData, errorList []error) {
	var e error
	data := new(ScadaData).New()

	// tk.Printf("%#v \n", conf.Columns)

	for _, val := range conf.Columns {
		if val["TimeStamp"] != nil {
			col := val["TimeStamp"]
			if len(col) > 1 {
				data.TimeStamp, e = GetDateCellAuto(row.Cells[col[0].Index], row.Cells[col[1].Index])
				ErrorLog(e, funcName, errorList)
			} else {
				// not implemented yet

				/*data.TimeStamp, e = GetDateCellAuto(row.Cells[col[0].Index])
				ErrorLog(e, funcName, errorList)*/
			}

			data.DateInfo = GetDateInfo(data.TimeStamp)
		} else if val["Turbine"] != nil {
			col := val["Turbine"]
			data.Turbine, e = row.Cells[col[0].Index].String()
			ErrorLog(e, funcName, errorList)
		} else if val["GridFrequency"] != nil {
			col := val["GridFrequency"]
			data.GridFrequency, e = GetFloatCell(row.Cells[col[0].Index])
			ErrorLog(e, funcName, errorList)
		} else if val["ReactivePower"] != nil {
			col := val["ReactivePower"]
			data.ReactivePower, e = GetFloatCell(row.Cells[col[0].Index])
			ErrorLog(e, funcName, errorList)
		} else if val["AlarmExtStopTime"] != nil {
			col := val["AlarmExtStopTime"]
			data.AlarmExtStopTime, e = GetFloatCell(row.Cells[col[0].Index])
			ErrorLog(e, funcName, errorList)
		} else if val["AlarmGridDownTime"] != nil {
			col := val["AlarmGridDownTime"]
			data.AlarmGridDownTime, e = GetFloatCell(row.Cells[col[0].Index])
			ErrorLog(e, funcName, errorList)
		} else if val["AlarmInterLineDown"] != nil {
			col := val["AlarmInterLineDown"]
			data.AlarmInterLineDown, e = GetFloatCell(row.Cells[col[0].Index])
			ErrorLog(e, funcName, errorList)
		} else if val["AlarmMachDownTime"] != nil {
			col := val["AlarmMachDownTime"]
			data.AlarmMachDownTime, e = GetFloatCell(row.Cells[col[0].Index])
			ErrorLog(e, funcName, errorList)
		} else if val["AlarmOkTime"] != nil {
			col := val["AlarmOkTime"]
			data.AlarmOkTime, e = GetFloatCell(row.Cells[col[0].Index])
			ErrorLog(e, funcName, errorList)
		} else if val["AlarmUnknownTime"] != nil {
			col := val["AlarmUnknownTime"]
			data.AlarmUnknownTime, e = GetFloatCell(row.Cells[col[0].Index])
			ErrorLog(e, funcName, errorList)
		} else if val["AlarmWeatherStop"] != nil {
			col := val["AlarmWeatherStop"]
			data.AlarmWeatherStop, e = GetFloatCell(row.Cells[col[0].Index])
			ErrorLog(e, funcName, errorList)
		} else if val["ExternalStopTime"] != nil {
			col := val["ExternalStopTime"]
			data.ExternalStopTime, e = GetFloatCell(row.Cells[col[0].Index])
			ErrorLog(e, funcName, errorList)
		} else if val["GridDownTime"] != nil {
			col := val["GridDownTime"]
			data.GridDownTime, e = GetFloatCell(row.Cells[col[0].Index])
			ErrorLog(e, funcName, errorList)
		} else if val["GridOkSecs"] != nil {
			col := val["GridOkSecs"]
			data.GridOkSecs, e = GetFloatCell(row.Cells[col[0].Index])
			ErrorLog(e, funcName, errorList)
		} else if val["InternalLineDown"] != nil {
			col := val["InternalLineDown"]
			data.InternalLineDown, e = GetFloatCell(row.Cells[col[0].Index])
			ErrorLog(e, funcName, errorList)
		} else if val["MachineDownTime"] != nil {
			col := val["MachineDownTime"]
			data.MachineDownTime, e = GetFloatCell(row.Cells[col[0].Index])
			ErrorLog(e, funcName, errorList)
		} else if val["OkSecs"] != nil {
			col := val["OkSecs"]
			data.OkSecs, e = GetFloatCell(row.Cells[col[0].Index])
			ErrorLog(e, funcName, errorList)
		} else if val["OkTime"] != nil {
			col := val["OkTime"]
			data.OkTime, e = GetFloatCell(row.Cells[col[0].Index])
			ErrorLog(e, funcName, errorList)
		} else if val["UnknownTime"] != nil {
			col := val["UnknownTime"]
			data.UnknownTime, e = GetFloatCell(row.Cells[col[0].Index])
			ErrorLog(e, funcName, errorList)
		} else if val["WeatherStopTime"] != nil {
			col := val["WeatherStopTime"]
			data.WeatherStopTime, e = GetFloatCell(row.Cells[col[0].Index])
			ErrorLog(e, funcName, errorList)
		} else if val["GeneratorRPM"] != nil {
			col := val["GeneratorRPM"]
			data.GeneratorRPM, e = GetFloatCell(row.Cells[col[0].Index])
			ErrorLog(e, funcName, errorList)
		} else if val["NacelleYawPositionUntwist"] != nil {
			col := val["NacelleYawPositionUntwist"]
			data.NacelleYawPositionUntwist, e = GetFloatCell(row.Cells[col[0].Index])
			ErrorLog(e, funcName, errorList)
		} else if val["NacelleTemperature"] != nil {
			col := val["NacelleTemperature"]
			data.NacelleTemperature, e = GetFloatCell(row.Cells[col[0].Index])
			ErrorLog(e, funcName, errorList)
		} else if val["AdjWindSpeed"] != nil {
			col := val["AdjWindSpeed"]
			data.AdjWindSpeed, e = GetFloatCell(row.Cells[col[0].Index])
			ErrorLog(e, funcName, errorList)
		} else if val["AmbientTemperature"] != nil {
			col := val["AmbientTemperature"]
			data.AmbientTemperature, e = GetFloatCell(row.Cells[col[0].Index])
			ErrorLog(e, funcName, errorList)
		} else if val["AvgBladeAngle"] != nil {
			col := val["AvgBladeAngle"]
			data.AvgBladeAngle, e = GetFloatCell(row.Cells[col[0].Index])
			ErrorLog(e, funcName, errorList)
		} else if val["AvgWindSpeed"] != nil {
			col := val["AvgWindSpeed"]
			data.AvgWindSpeed, e = GetFloatCell(row.Cells[col[0].Index])
			ErrorLog(e, funcName, errorList)
		} else if val["UnitsGenerated"] != nil {
			col := val["UnitsGenerated"]
			data.UnitsGenerated, e = GetFloatCell(row.Cells[col[0].Index])
			ErrorLog(e, funcName, errorList)
		} else if val["EstimatedPower"] != nil {
			col := val["EstimatedPower"]
			data.EstimatedPower, e = GetFloatCell(row.Cells[col[0].Index])
			ErrorLog(e, funcName, errorList)
		} else if val["NacelDirection"] != nil {
			col := val["NacelDirection"]
			data.NacelDirection, e = GetFloatCell(row.Cells[col[0].Index])
			ErrorLog(e, funcName, errorList)
		} else if val["Power"] != nil {
			col := val["Power"]
			data.Power, e = GetFloatCell(row.Cells[col[0].Index])
			ErrorLog(e, funcName, errorList)
		} else if val["PowerLost"] != nil {
			col := val["PowerLost"]
			data.PowerLost, e = GetFloatCell(row.Cells[col[0].Index])
			ErrorLog(e, funcName, errorList)
		} else if val["RotorRPM"] != nil {
			if len(row.Cells) > 33 {
				col := val["RotorRPM"]
				data.RotorRPM, e = GetFloatCell(row.Cells[col[0].Index])
				ErrorLog(e, funcName, errorList)
			}
		} else if val["WindDirection"] != nil {
			if len(row.Cells) > 34 {
				col := val["WindDirection"]
				data.WindDirection, e = GetFloatCell(row.Cells[col[0].Index])
				ErrorLog(e, funcName, errorList)
			}
		}

		data.TotalTime = data.ExternalStopTime + data.GridDownTime + data.InternalLineDown + data.MachineDownTime + data.OkTime
		data.Minutes = 10
		data.IsValidTimeDuration = true

		if data.AvgWindSpeed < 4 || (data.AvgWindSpeed >= 4 && data.Power > 0) {
			data.Available = 1
		} else if data.AvgWindSpeed >= 4 && data.Power <= 0 {
			data.Available = 0
		}
	}

	res = data
	return
}

func ConstructScadaDataCSV(conf ScadaConf, row []string) (res *ScadaDataNew, errorList []error) {
	var err error
	var stat bool

	res = new(ScadaDataNew).New()
out:
	for _, val := range conf.Columns {
		if val["TimeStamp"] != nil {
			col := val["TimeStamp"]
			res.TimeStamp, err = time.Parse("02-Jan-2006 15:04", tk.ToString(row[col[0].Index]))
			ErrorLog(err, funcName, errorList)
			if err != nil {
				break out
			}
			res.DateInfo = GetDateInfo(res.TimeStamp)
		} else if val["Turbine"] != nil {
			col := val["Turbine"]
			turbines := strings.Split(row[col[0].Index], ".")
			res.Turbine = turbines[2]
		} else if val["Power"] != nil {
			col := val["Power"]
			res.Power, stat = tk.StringToFloat(row[col[0].Index])
			if !stat {
				ErrorLog(errors.New("data not found for: Power"), funcName, errorList)
			} else {
				res.Energy = tk.Div(res.Power, 6.0)
			}
		} else if val["AvgWindSpeed"] != nil {
			col := val["AvgWindSpeed"]
			res.AvgWindSpeed, stat = tk.StringToFloat(row[col[0].Index])
			if !stat {
				ErrorLog(errors.New("data not found for: AvgWindSpeed"), funcName, errorList)
			} else {
				res.AdjWindSpeed = tk.RoundingAuto64(res.AvgWindSpeed, 1)
			}
		} else if val["WindDirection"] != nil {
			col := val["WindDirection"]
			res.WindDirection, stat = tk.StringToFloat(row[col[0].Index])
			if !stat {
				ErrorLog(errors.New("data not found for: WindDirection"), funcName, errorList)
			}
		} else if val["NacelDirection"] != nil {
			col := val["NacelDirection"]
			res.NacelDirection, stat = tk.StringToFloat(row[col[0].Index])
			if !stat {
				res.NacelDirection = 0.0
				ErrorLog(errors.New("data not found for: NacelDirection"), funcName, errorList)
			}
		} else if val["RotorRPM"] != nil {
			col := val["RotorRPM"]
			res.RotorRPM, stat = tk.StringToFloat(row[col[0].Index])
			if !stat {
				ErrorLog(errors.New("data not found for: RotorRPM"), funcName, errorList)
			}
		} else if val["ReactivePower"] != nil {
			col := val["ReactivePower"]
			res.ReactivePower, stat = tk.StringToFloat(row[col[0].Index])
			if !stat {
				ErrorLog(errors.New("data not found for: ReactivePower"), funcName, errorList)
			}
		}
	}

	if err == nil {
		res.EnergyLost = 0.0
		res.GridFrequency = 0.0
		res.AlarmExtStopTime = 0.0
		res.AlarmGridDownTime = 0.0
		res.AlarmInterLineDown = 0.0
		res.AlarmMachDownTime = 0.0
		res.AlarmOkTime = 0.0
		res.AlarmUnknownTime = 0.0
		res.AlarmWeatherStop = 0.0
		res.ExternalStopTime = 0.0
		res.GridDownTime = 0.0
		res.GridOkSecs = 0.0
		res.InternalLineDown = 0.0
		res.MachineDownTime = 0.0
		res.OkSecs = 0.0
		res.OkTime = 0.0
		res.UnknownTime = 0.0
		res.WeatherStopTime = 0.0
		res.GeneratorRPM = 0.0
		res.NacelleYawPositionUntwist = 0.0
		res.NacelleTemperature = 0.0
		res.AmbientTemperature = 0.0
		res.AvgBladeAngle = 0.0
		res.UnitsGenerated = 0.0
		res.EstimatedPower = 0.0
		res.EstimatedEnergy = 0.0
		res.TotalTime = 600.0
		res.Minutes = 10
		res.Available = 1
		res.DenWindSpeed = 0.0
		res.DenAdjWindSpeed = 0.0
		res.DenPower = 0.0
		res.DenEnergy = 0.0
		res.IsValidTimeDuration = false

		if res.AvgWindSpeed >= 3.5 && res.Power <= 0 {
			res.Available = 0
		} else {
			res.OkSecs = 600.0
			res.IsValidTimeDuration = true
			res.TotalTime = 600.0
		}
	}

	// tk.Printf("\n%#v \n", res)

	return
}
