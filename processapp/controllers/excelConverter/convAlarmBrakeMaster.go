package converterControllers

import (
	. "github.com/eaciit/windapp/library/helper"
	. "github.com/eaciit/windapp/library/models"
	. "github.com/eaciit/windapp/processapp/controllers"
	_ "github.com/eaciit/dbox/dbc/mongo"
	tk "github.com/eaciit/toolkit"
	"github.com/tealeg/xlsx"
	"os"
	"strings"
)

type ConvAlarmBrakeMaster struct {
	*BaseController
}

func (d *ConvAlarmBrakeMaster) Generate(base *BaseController) {
	folderName := "alarmbrakemaster"
	funcName := "Converting Alarm Brake Data"
	count := 0
	total := 0
	if base != nil {
		d.BaseController = base

		ctx := d.BaseController.Ctx
		dataSources, path := base.GetDataSource(folderName)
		tk.Println("Converting Alarm Brake Data from Excel File..")
		for _, source := range dataSources {
			if strings.Contains(source.Name(), "AlarmBrake") {
				tk.Println(path + "\\" + source.Name())
				file, e := xlsx.OpenFile(path + "\\" + source.Name())
				if e != nil {
					ErrorHandler(e, funcName)
					os.Exit(0)
				}

				for _, sheet := range file.Sheet {
					errorLine := tk.M{}
					for idx, row := range sheet.Rows {
						if idx > 0 {
							data := new(AlarmBrake).New()

							data.TypeCode, _ = row.Cells[0].Int()
							data.AlarmIndex, e = row.Cells[1].Int()
							if e != nil {
								data.AlarmIndex = 0
							}
							data.AlarmName, _ = row.Cells[2].String()
							data.AlarmTypeId, _ = row.Cells[3].String()
							data.TypeId, _ = row.Cells[4].Int()
							data.Type, _ = row.Cells[5].String()
							data.Set = row.Cells[6].Bool()
							data.Disabled = row.Cells[7].Bool()
							data.DefaultDisabled = row.Cells[8].Bool()
							data.BrakeProgram, _ = row.Cells[9].Int()
							data.DefaultBrakeProgram, _ = row.Cells[10].Int()
							data.YawProgram, _ = row.Cells[11].Int()
							data.DefaultYawProgram, _ = row.Cells[12].Int()
							data.AlarmPaging = row.Cells[13].Bool()
							data.DefaultAlarmPaging = row.Cells[14].Bool()
							data.AlarmDelay, _ = row.Cells[15].Int()
							data.DefaultAlarmDelay, _ = row.Cells[16].Int()
							data.AlarmDelayUnit, _ = row.Cells[17].String()
							data.ReducesAvailability = row.Cells[18].Bool()
							data.DefaultReducesAvailability = row.Cells[19].Bool()
							data.OnTimeCounter, _ = row.Cells[20].Int()
							data.AlarmCounter, _ = row.Cells[21].Int()
							data.RepeatAlarmCode, _ = row.Cells[22].Int()
							data.RepeatAlarmName, _ = row.Cells[23].String()
							data.RepeatAlarmNumber, _ = row.Cells[24].Int()
							data.DefaultRepeatAlarmCounter, _ = row.Cells[25].Int()
							data.RepeatAlarmTime, _ = row.Cells[26].Int()
							data.DefaultRepeatAlarmTime, _ = row.Cells[27].Int()
							data.LevelDisableAlarm, _ = row.Cells[28].Int()
							data.LevelResetAlarm, _ = row.Cells[29].Int()

							// data.GridDown = row.Cells[6].Bool()
							// data.InternalGrid = row.Cells[7].Bool()
							// data.MachineDown = row.Cells[8].Bool()
							// data.AEbOK = row.Cells[9].Bool()
							// data.Unknown = row.Cells[10].Bool()
							// data.WeatherStop = row.Cells[11].Bool()
							// data.EndDate, e = GetDateCellAuto(row.Cells[12], row.Cells[13])
							// ErrorLog(e, funcName, errorList)
							// if row.Cells[12].GetNumberFormat() != "general" {
							// 	data.EndDate, e = ReverseMonthDate(data.EndDate)
							// }

							e = ctx.Insert(data)

							count++
							if count == 1000 {
								total += count
								tk.Printf("count: %v \n", total)
								count = 0
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
