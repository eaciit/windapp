package converterControllers

import (
	. "eaciit/wfdemo-git-dev/library/helper"
	. "eaciit/wfdemo-git-dev/library/models"
	. "eaciit/wfdemo-git-dev/processapp/controllers"
	"os"
	"time"

	"strings"

	tk "github.com/eaciit/toolkit"
	"github.com/tealeg/xlsx"
)

// ConvScadaDataOEM
type ConvScadaDataOEM struct {
	*BaseController
}

// Generate
func (d *ConvScadaDataOEM) Generate(base *BaseController) {
	funcName := "Converting Scada Data OEM"
	project := "Tejuva"
	folder := "scadaOEM/11_NOV/16.11.16to24.11.16/TML"
	if base != nil {
		d.BaseController = base

		ctx := d.BaseController.Ctx
		_ = ctx
		dataSources, path := base.GetDataSource(folder)
		tk.Println("Converting Scada Data OEM from Excel File..")
		for _, source := range dataSources {
			count := 0
			total := 0
			errorLine := tk.M{}
			// EXCEL Type
			// if strings.Contains(source.Name(), "Scada") {
			anName := strings.Split(source.Name(), ".")
			turbine := anName[0]
			tk.Printf("turbine: %v \n", turbine)
			tk.Println(path + "/" + source.Name())
			file, e := xlsx.OpenFile(path + "/" + source.Name())
			if e != nil {
				ErrorHandler(e, funcName)
				os.Exit(0)
			}

			for _, sheet := range file.Sheet {
				errorLine = tk.M{}
				for idx, row := range sheet.Rows {
					errorList := []error{}
					if idx > 0 { //&& idx < 5 { //&& len(row.Cells) == 35 {
						data, errorList := ConstructScadaDataOEMExcel(row)
						data.Line = idx + 1
						data.ProjectName = project
						data.Turbine = turbine

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
			// }

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

func ConstructScadaDataOEMExcel(row *xlsx.Row) (res *ScadaDataOEM, errorList []error) {
	var e error
	data := new(ScadaDataOEM).New()
	columns := data.GetXlsColumns()
	// tk.Printf("%#v \n", conf.Columns)

	for _, col := range columns {
		var colInfo tk.M
		colInfo = data.GetColumnInfos(col)
		idx := colInfo.GetInt("Index")

		if col == "TimeStamp" {
			str, e := row.Cells[idx].String()
			dtStr := strings.Split(str, "T")
			dtStrTime := strings.Split(dtStr[1], ".000")
			locale := strings.Replace(dtStrTime[1], ":", "", 1)
			data.TimeStamp, e = time.Parse("2006-01-02 15:04:05", dtStr[0]+" "+dtStrTime[0])
			ErrorLog(e, funcName, errorList)
			if e == nil {
				data.DateInfo = GetDateInfo(data.TimeStamp)

				x, e := time.Parse("2006-01-02 15:04:05 -0700", dtStr[0]+" "+dtStrTime[0]+" "+locale)
				data.TimeStampUTC = x.UTC()

				ErrorLog(e, funcName, errorList)
				if e == nil {
					data.DateInfoUTC = GetDateInfo(data.TimeStampUTC)
				}
			}
		} else if col == "AI_intern_R_PidAngleOut" {
			data.AI_intern_R_PidAngleOut, e = GetFloatCell(row.Cells[idx])
			ErrorLog(e, funcName, errorList)

		} else if col == "AI_intern_ActivPower" {
			data.AI_intern_ActivPower, e = GetFloatCell(row.Cells[idx])
			ErrorLog(e, funcName, errorList)

		} else if col == "AI_intern_I1" {
			data.AI_intern_I1, e = GetFloatCell(row.Cells[idx])
			ErrorLog(e, funcName, errorList)

		} else if col == "AI_intern_I2" {
			data.AI_intern_I2, e = GetFloatCell(row.Cells[idx])
			ErrorLog(e, funcName, errorList)

		} else if col == "AI_intern_I3" {
			data.AI_intern_I3, e = GetFloatCell(row.Cells[idx])
			ErrorLog(e, funcName, errorList)

		} else if col == "AI_intern_NacelleDrill" {
			data.AI_intern_NacelleDrill, e = GetFloatCell(row.Cells[idx])
			ErrorLog(e, funcName, errorList)

		} else if col == "AI_intern_NacellePos" {
			data.AI_intern_NacellePos, e = GetFloatCell(row.Cells[idx])
			ErrorLog(e, funcName, errorList)

		} else if col == "AI_intern_PitchAkku_V1" {
			data.AI_intern_PitchAkku_V1, e = GetFloatCell(row.Cells[idx])
			ErrorLog(e, funcName, errorList)

		} else if col == "AI_intern_PitchAkku_V2" {
			data.AI_intern_PitchAkku_V2, e = GetFloatCell(row.Cells[idx])
			ErrorLog(e, funcName, errorList)

		} else if col == "AI_intern_PitchAkku_V3" {
			data.AI_intern_PitchAkku_V3, e = GetFloatCell(row.Cells[idx])
			ErrorLog(e, funcName, errorList)

		} else if col == "AI_intern_PitchAngle1" {
			data.AI_intern_PitchAngle1, e = GetFloatCell(row.Cells[idx])
			ErrorLog(e, funcName, errorList)

		} else if col == "AI_intern_PitchAngle2" {
			data.AI_intern_PitchAngle2, e = GetFloatCell(row.Cells[idx])
			ErrorLog(e, funcName, errorList)

		} else if col == "AI_intern_PitchAngle3" {
			data.AI_intern_PitchAngle3, e = GetFloatCell(row.Cells[idx])
			ErrorLog(e, funcName, errorList)

		} else if col == "AI_intern_PitchConv_Current1" {
			data.AI_intern_PitchConv_Current1, e = GetFloatCell(row.Cells[idx])
			ErrorLog(e, funcName, errorList)

		} else if col == "AI_intern_PitchConv_Current2" {
			data.AI_intern_PitchConv_Current2, e = GetFloatCell(row.Cells[idx])
			ErrorLog(e, funcName, errorList)

		} else if col == "AI_intern_PitchConv_Current3" {
			data.AI_intern_PitchConv_Current3, e = GetFloatCell(row.Cells[idx])
			ErrorLog(e, funcName, errorList)

		} else if col == "AI_intern_PitchAngleSP_Diff1" {
			data.AI_intern_PitchAngleSP_Diff1, e = GetFloatCell(row.Cells[idx])
			ErrorLog(e, funcName, errorList)

		} else if col == "AI_intern_PitchAngleSP_Diff2" {
			data.AI_intern_PitchAngleSP_Diff2, e = GetFloatCell(row.Cells[idx])
			ErrorLog(e, funcName, errorList)

		} else if col == "AI_intern_PitchAngleSP_Diff3" {
			data.AI_intern_PitchAngleSP_Diff3, e = GetFloatCell(row.Cells[idx])
			ErrorLog(e, funcName, errorList)

		} else if col == "AI_intern_ReactivPower" {
			data.AI_intern_ReactivPower, e = GetFloatCell(row.Cells[idx])
			ErrorLog(e, funcName, errorList)

		} else if col == "AI_intern_RpmDiff" {
			data.AI_intern_RpmDiff, e = GetFloatCell(row.Cells[idx])
			ErrorLog(e, funcName, errorList)

		} else if col == "AI_intern_U1" {
			data.AI_intern_U1, e = GetFloatCell(row.Cells[idx])
			ErrorLog(e, funcName, errorList)

		} else if col == "AI_intern_U2" {
			data.AI_intern_U2, e = GetFloatCell(row.Cells[idx])
			ErrorLog(e, funcName, errorList)

		} else if col == "AI_intern_U3" {
			data.AI_intern_U3, e = GetFloatCell(row.Cells[idx])
			ErrorLog(e, funcName, errorList)

		} else if col == "AI_intern_WindDirection" {
			data.AI_intern_WindDirection, e = GetFloatCell(row.Cells[idx])
			ErrorLog(e, funcName, errorList)

		} else if col == "AI_intern_WindSpeed" {
			data.AI_intern_WindSpeed, e = GetFloatCell(row.Cells[idx])
			ErrorLog(e, funcName, errorList)

		} else if col == "AI_Intern_WindSpeedDif" {
			data.AI_Intern_WindSpeedDif, e = GetFloatCell(row.Cells[idx])
			ErrorLog(e, funcName, errorList)

		} else if col == "AI_speed_RotFR" {
			data.AI_speed_RotFR, e = GetFloatCell(row.Cells[idx])
			ErrorLog(e, funcName, errorList)

		} else if col == "AI_WindSpeed1" {
			data.AI_WindSpeed1, e = GetFloatCell(row.Cells[idx])
			ErrorLog(e, funcName, errorList)

		} else if col == "AI_WindSpeed2" {
			data.AI_WindSpeed2, e = GetFloatCell(row.Cells[idx])
			ErrorLog(e, funcName, errorList)

		} else if col == "AI_WindVane1" {
			data.AI_WindVane1, e = GetFloatCell(row.Cells[idx])
			ErrorLog(e, funcName, errorList)

		} else if col == "AI_WindVane2" {
			data.AI_WindVane2, e = GetFloatCell(row.Cells[idx])
			ErrorLog(e, funcName, errorList)

		} else if col == "AI_internCurrentAsym" {
			data.AI_internCurrentAsym, e = GetFloatCell(row.Cells[idx])
			ErrorLog(e, funcName, errorList)

		} else if col == "Temp_GearBox_IMS_NDE" {
			data.Temp_GearBox_IMS_NDE, e = GetFloatCell(row.Cells[idx])
			ErrorLog(e, funcName, errorList)

		} else if col == "AI_intern_WindVaneDiff" {
			data.AI_intern_WindVaneDiff, e = GetFloatCell(row.Cells[idx])
			ErrorLog(e, funcName, errorList)

		} else if col == "C_intern_SpeedGenerator" {
			data.C_intern_SpeedGenerator, e = GetFloatCell(row.Cells[idx])
			ErrorLog(e, funcName, errorList)

		} else if col == "C_intern_SpeedRotor" {
			data.C_intern_SpeedRotor, e = GetFloatCell(row.Cells[idx])
			ErrorLog(e, funcName, errorList)

		} else if col == "AI_intern_Speed_RPMDiff_FR1_RotCNT" {
			data.AI_intern_Speed_RPMDiff_FR1_RotCNT, e = GetFloatCell(row.Cells[idx])
			ErrorLog(e, funcName, errorList)

		} else if col == "AI_intern_Frequency_Grid" {
			data.AI_intern_Frequency_Grid, e = GetFloatCell(row.Cells[idx])
			ErrorLog(e, funcName, errorList)

		} else if col == "Temp_GearBox_HSS_NDE" {
			data.Temp_GearBox_HSS_NDE, e = GetFloatCell(row.Cells[idx])
			ErrorLog(e, funcName, errorList)

		} else if col == "AI_DrTrVibValue" {
			data.AI_DrTrVibValue, e = GetFloatCell(row.Cells[idx])
			ErrorLog(e, funcName, errorList)

		} else if col == "AI_intern_InLastErrorConv1" {
			data.AI_intern_InLastErrorConv1, e = GetFloatCell(row.Cells[idx])
			ErrorLog(e, funcName, errorList)

		} else if col == "AI_intern_InLastErrorConv2" {
			data.AI_intern_InLastErrorConv2, e = GetFloatCell(row.Cells[idx])
			ErrorLog(e, funcName, errorList)

		} else if col == "AI_intern_InLastErrorConv3" {
			data.AI_intern_InLastErrorConv3, e = GetFloatCell(row.Cells[idx])
			ErrorLog(e, funcName, errorList)

		} else if col == "AI_intern_TempConv1" {
			data.AI_intern_TempConv1, e = GetFloatCell(row.Cells[idx])
			ErrorLog(e, funcName, errorList)

		} else if col == "AI_intern_TempConv2" {
			data.AI_intern_TempConv2, e = GetFloatCell(row.Cells[idx])
			ErrorLog(e, funcName, errorList)

		} else if col == "AI_intern_TempConv3" {
			data.AI_intern_TempConv3, e = GetFloatCell(row.Cells[idx])
			ErrorLog(e, funcName, errorList)

		} else if col == "AI_intern_PitchSpeed1" {
			data.AI_intern_PitchSpeed1, e = GetFloatCell(row.Cells[idx])
			ErrorLog(e, funcName, errorList)

		} else if col == "Temp_YawBrake_1" {
			data.Temp_YawBrake_1, e = GetFloatCell(row.Cells[idx])
			ErrorLog(e, funcName, errorList)

		} else if col == "Temp_YawBrake_2" {
			data.Temp_YawBrake_2, e = GetFloatCell(row.Cells[idx])
			ErrorLog(e, funcName, errorList)

		} else if col == "Temp_G1L1" {
			data.Temp_G1L1, e = GetFloatCell(row.Cells[idx])
			ErrorLog(e, funcName, errorList)

		} else if col == "Temp_G1L2" {
			data.Temp_G1L2, e = GetFloatCell(row.Cells[idx])
			ErrorLog(e, funcName, errorList)

		} else if col == "Temp_G1L3" {
			data.Temp_G1L3, e = GetFloatCell(row.Cells[idx])
			ErrorLog(e, funcName, errorList)

		} else if col == "Temp_YawBrake_3" {
			data.Temp_YawBrake_3, e = GetFloatCell(row.Cells[idx])
			ErrorLog(e, funcName, errorList)

		} else if col == "AI_HydrSystemPressure" {
			data.AI_HydrSystemPressure, e = GetFloatCell(row.Cells[idx])
			ErrorLog(e, funcName, errorList)

		} else if col == "Temp_BottomControlSection_Low" {
			data.Temp_BottomControlSection_Low, e = GetFloatCell(row.Cells[idx])
			ErrorLog(e, funcName, errorList)

		} else if col == "Temp_GearBox_HSS_DE" {
			data.Temp_GearBox_HSS_DE, e = GetFloatCell(row.Cells[idx])
			ErrorLog(e, funcName, errorList)

		} else if col == "Temp_GearOilSump" {
			data.Temp_GearOilSump, e = GetFloatCell(row.Cells[idx])
			ErrorLog(e, funcName, errorList)

		} else if col == "Temp_GeneratorBearing_DE" {
			data.Temp_GeneratorBearing_DE, e = GetFloatCell(row.Cells[idx])
			ErrorLog(e, funcName, errorList)

		} else if col == "Temp_GeneratorBearing_NDE" {
			data.Temp_GeneratorBearing_NDE, e = GetFloatCell(row.Cells[idx])
			ErrorLog(e, funcName, errorList)

		} else if col == "Temp_MainBearing" {
			data.Temp_MainBearing, e = GetFloatCell(row.Cells[idx])
			ErrorLog(e, funcName, errorList)

		} else if col == "Temp_GearBox_IMS_DE" {
			data.Temp_GearBox_IMS_DE, e = GetFloatCell(row.Cells[idx])
			ErrorLog(e, funcName, errorList)

		} else if col == "Temp_Nacelle" {
			data.Temp_Nacelle, e = GetFloatCell(row.Cells[idx])
			ErrorLog(e, funcName, errorList)

		} else if col == "Temp_Outdoor" {
			data.Temp_Outdoor, e = GetFloatCell(row.Cells[idx])
			ErrorLog(e, funcName, errorList)

		} else if col == "AI_TowerVibValueAxial" {
			data.AI_TowerVibValueAxial, e = GetFloatCell(row.Cells[idx])
			ErrorLog(e, funcName, errorList)

		} else if col == "AI_intern_DiffGenSpeedSPToAct" {
			data.AI_intern_DiffGenSpeedSPToAct, e = GetFloatCell(row.Cells[idx])
			ErrorLog(e, funcName, errorList)

		} else if col == "Temp_YawBrake_4" {
			data.Temp_YawBrake_4, e = GetFloatCell(row.Cells[idx])
			ErrorLog(e, funcName, errorList)

		} else if col == "AI_intern_SpeedGenerator_Proximity" {
			data.AI_intern_SpeedGenerator_Proximity, e = GetFloatCell(row.Cells[idx])
			ErrorLog(e, funcName, errorList)

		} else if col == "AI_intern_SpeedDiff_Encoder_Proximity" {
			data.AI_intern_SpeedDiff_Encoder_Proximity, e = GetFloatCell(row.Cells[idx])
			ErrorLog(e, funcName, errorList)

		} else if col == "AI_GearOilPressure" {
			data.AI_GearOilPressure, e = GetFloatCell(row.Cells[idx])
			ErrorLog(e, funcName, errorList)

		} else if col == "Temp_CabinetTopBox_Low" {
			data.Temp_CabinetTopBox_Low, e = GetFloatCell(row.Cells[idx])
			ErrorLog(e, funcName, errorList)

		} else if col == "Temp_CabinetTopBox" {
			data.Temp_CabinetTopBox, e = GetFloatCell(row.Cells[idx])
			ErrorLog(e, funcName, errorList)

		} else if col == "Temp_BottomControlSection" {
			data.Temp_BottomControlSection, e = GetFloatCell(row.Cells[idx])
			ErrorLog(e, funcName, errorList)

		} else if col == "Temp_BottomPowerSection" {
			data.Temp_BottomPowerSection, e = GetFloatCell(row.Cells[idx])
			ErrorLog(e, funcName, errorList)

		} else if col == "Temp_BottomPowerSection_Low" {
			data.Temp_BottomPowerSection_Low, e = GetFloatCell(row.Cells[idx])
			ErrorLog(e, funcName, errorList)

		} else if col == "AI_intern_Pitch1_Status_High" {
			data.AI_intern_Pitch1_Status_High, e = GetFloatCell(row.Cells[idx])
			ErrorLog(e, funcName, errorList)

		} else if col == "AI_intern_Pitch2_Status_High" {
			data.AI_intern_Pitch2_Status_High, e = GetFloatCell(row.Cells[idx])
			ErrorLog(e, funcName, errorList)

		} else if col == "AI_intern_Pitch3_Status_High" {
			data.AI_intern_Pitch3_Status_High, e = GetFloatCell(row.Cells[idx])
			ErrorLog(e, funcName, errorList)

		} else if col == "AI_intern_InPosition1_ch2" {
			data.AI_intern_InPosition1_ch2, e = GetFloatCell(row.Cells[idx])
			ErrorLog(e, funcName, errorList)

		} else if col == "AI_intern_InPosition2_ch2" {
			data.AI_intern_InPosition2_ch2, e = GetFloatCell(row.Cells[idx])
			ErrorLog(e, funcName, errorList)

		} else if col == "AI_intern_InPosition3_ch2" {
			data.AI_intern_InPosition3_ch2, e = GetFloatCell(row.Cells[idx])
			ErrorLog(e, funcName, errorList)

		} else if col == "AI_intern_Temp_Brake_Blade1" {
			data.AI_intern_Temp_Brake_Blade1, e = GetFloatCell(row.Cells[idx])
			ErrorLog(e, funcName, errorList)

		} else if col == "AI_intern_Temp_Brake_Blade2" {
			data.AI_intern_Temp_Brake_Blade2, e = GetFloatCell(row.Cells[idx])
			ErrorLog(e, funcName, errorList)

		} else if col == "AI_intern_Temp_Brake_Blade3" {
			data.AI_intern_Temp_Brake_Blade3, e = GetFloatCell(row.Cells[idx])
			ErrorLog(e, funcName, errorList)

		} else if col == "AI_intern_Temp_PitchMotor_Blade1" {
			data.AI_intern_Temp_PitchMotor_Blade1, e = GetFloatCell(row.Cells[idx])
			ErrorLog(e, funcName, errorList)

		} else if col == "AI_intern_Temp_PitchMotor_Blade2" {
			data.AI_intern_Temp_PitchMotor_Blade2, e = GetFloatCell(row.Cells[idx])
			ErrorLog(e, funcName, errorList)

		} else if col == "AI_intern_Temp_PitchMotor_Blade3" {
			data.AI_intern_Temp_PitchMotor_Blade3, e = GetFloatCell(row.Cells[idx])
			ErrorLog(e, funcName, errorList)

		} else if col == "AI_intern_Temp_Hub_Additional1" {
			data.AI_intern_Temp_Hub_Additional1, e = GetFloatCell(row.Cells[idx])
			ErrorLog(e, funcName, errorList)

		} else if col == "AI_intern_Temp_Hub_Additional2" {
			data.AI_intern_Temp_Hub_Additional2, e = GetFloatCell(row.Cells[idx])
			ErrorLog(e, funcName, errorList)

		} else if col == "AI_intern_Temp_Hub_Additional3" {
			data.AI_intern_Temp_Hub_Additional3, e = GetFloatCell(row.Cells[idx])
			ErrorLog(e, funcName, errorList)

		} else if col == "AI_intern_Pitch1_Status_Low" {
			data.AI_intern_Pitch1_Status_Low, e = GetFloatCell(row.Cells[idx])
			ErrorLog(e, funcName, errorList)

		} else if col == "AI_intern_Pitch2_Status_Low" {
			data.AI_intern_Pitch2_Status_Low, e = GetFloatCell(row.Cells[idx])
			ErrorLog(e, funcName, errorList)

		} else if col == "AI_intern_Pitch3_Status_Low" {
			data.AI_intern_Pitch3_Status_Low, e = GetFloatCell(row.Cells[idx])
			ErrorLog(e, funcName, errorList)

		} else if col == "AI_intern_Battery_VoltageBlade1_center" {
			data.AI_intern_Battery_VoltageBlade1_center, e = GetFloatCell(row.Cells[idx])
			ErrorLog(e, funcName, errorList)

		} else if col == "AI_intern_Battery_VoltageBlade2_center" {
			data.AI_intern_Battery_VoltageBlade2_center, e = GetFloatCell(row.Cells[idx])
			ErrorLog(e, funcName, errorList)

		} else if col == "AI_intern_Battery_VoltageBlade3_center" {
			data.AI_intern_Battery_VoltageBlade3_center, e = GetFloatCell(row.Cells[idx])
			ErrorLog(e, funcName, errorList)

		} else if col == "AI_intern_Battery_ChargingCur_Blade1" {
			data.AI_intern_Battery_ChargingCur_Blade1, e = GetFloatCell(row.Cells[idx])
			ErrorLog(e, funcName, errorList)

		} else if col == "AI_intern_Battery_ChargingCur_Blade2" {
			data.AI_intern_Battery_ChargingCur_Blade2, e = GetFloatCell(row.Cells[idx])
			ErrorLog(e, funcName, errorList)

		} else if col == "AI_intern_Battery_ChargingCur_Blade3" {
			data.AI_intern_Battery_ChargingCur_Blade3, e = GetFloatCell(row.Cells[idx])
			ErrorLog(e, funcName, errorList)

		} else if col == "AI_intern_Battery_DischargingCur_Blade1" {
			data.AI_intern_Battery_DischargingCur_Blade1, e = GetFloatCell(row.Cells[idx])
			ErrorLog(e, funcName, errorList)

		} else if col == "AI_intern_Battery_DischargingCur_Blade2" {
			data.AI_intern_Battery_DischargingCur_Blade2, e = GetFloatCell(row.Cells[idx])
			ErrorLog(e, funcName, errorList)

		} else if col == "AI_intern_Battery_DischargingCur_Blade3" {
			data.AI_intern_Battery_DischargingCur_Blade3, e = GetFloatCell(row.Cells[idx])
			ErrorLog(e, funcName, errorList)

		} else if col == "AI_intern_PitchMotor_BrakeVoltage_Blade1" {
			data.AI_intern_PitchMotor_BrakeVoltage_Blade1, e = GetFloatCell(row.Cells[idx])
			ErrorLog(e, funcName, errorList)

		} else if col == "AI_intern_PitchMotor_BrakeVoltage_Blade2" {
			data.AI_intern_PitchMotor_BrakeVoltage_Blade2, e = GetFloatCell(row.Cells[idx])
			ErrorLog(e, funcName, errorList)

		} else if col == "AI_intern_PitchMotor_BrakeVoltage_Blade3" {
			data.AI_intern_PitchMotor_BrakeVoltage_Blade3, e = GetFloatCell(row.Cells[idx])
			ErrorLog(e, funcName, errorList)

		} else if col == "AI_intern_PitchMotor_BrakeCurrent_Blade1" {
			data.AI_intern_PitchMotor_BrakeCurrent_Blade1, e = GetFloatCell(row.Cells[idx])
			ErrorLog(e, funcName, errorList)

		} else if col == "AI_intern_PitchMotor_BrakeCurrent_Blade2" {
			data.AI_intern_PitchMotor_BrakeCurrent_Blade2, e = GetFloatCell(row.Cells[idx])
			ErrorLog(e, funcName, errorList)

		} else if col == "AI_intern_PitchMotor_BrakeCurrent_Blade3" {
			data.AI_intern_PitchMotor_BrakeCurrent_Blade3, e = GetFloatCell(row.Cells[idx])
			ErrorLog(e, funcName, errorList)

		} else if col == "AI_intern_Temp_HubBox_Blade1" {
			data.AI_intern_Temp_HubBox_Blade1, e = GetFloatCell(row.Cells[idx])
			ErrorLog(e, funcName, errorList)

		} else if col == "AI_intern_Temp_HubBox_Blade2" {
			data.AI_intern_Temp_HubBox_Blade2, e = GetFloatCell(row.Cells[idx])
			ErrorLog(e, funcName, errorList)

		} else if col == "AI_intern_Temp_HubBox_Blade3" {
			data.AI_intern_Temp_HubBox_Blade3, e = GetFloatCell(row.Cells[idx])
			ErrorLog(e, funcName, errorList)

		} else if col == "AI_intern_Temp_Pitch1_HeatSink" {
			data.AI_intern_Temp_Pitch1_HeatSink, e = GetFloatCell(row.Cells[idx])
			ErrorLog(e, funcName, errorList)

		} else if col == "AI_intern_Temp_Pitch2_HeatSink" {
			data.AI_intern_Temp_Pitch2_HeatSink, e = GetFloatCell(row.Cells[idx])
			ErrorLog(e, funcName, errorList)

		} else if col == "AI_intern_Temp_Pitch3_HeatSink" {
			data.AI_intern_Temp_Pitch3_HeatSink, e = GetFloatCell(row.Cells[idx])
			ErrorLog(e, funcName, errorList)

		} else if col == "AI_intern_ErrorStackBlade1" {
			data.AI_intern_ErrorStackBlade1, e = GetFloatCell(row.Cells[idx])
			ErrorLog(e, funcName, errorList)

		} else if col == "AI_intern_ErrorStackBlade2" {
			data.AI_intern_ErrorStackBlade2, e = GetFloatCell(row.Cells[idx])
			ErrorLog(e, funcName, errorList)

		} else if col == "AI_intern_ErrorStackBlade3" {
			data.AI_intern_ErrorStackBlade3, e = GetFloatCell(row.Cells[idx])
			ErrorLog(e, funcName, errorList)

		} else if col == "AI_intern_Temp_BatteryBox_Blade1" {
			data.AI_intern_Temp_BatteryBox_Blade1, e = GetFloatCell(row.Cells[idx])
			ErrorLog(e, funcName, errorList)

		} else if col == "AI_intern_Temp_BatteryBox_Blade2" {
			data.AI_intern_Temp_BatteryBox_Blade2, e = GetFloatCell(row.Cells[idx])
			ErrorLog(e, funcName, errorList)

		} else if col == "AI_intern_Temp_BatteryBox_Blade3" {
			data.AI_intern_Temp_BatteryBox_Blade3, e = GetFloatCell(row.Cells[idx])
			ErrorLog(e, funcName, errorList)

		} else if col == "AI_intern_DC_LinkVoltage1" {
			data.AI_intern_DC_LinkVoltage1, e = GetFloatCell(row.Cells[idx])
			ErrorLog(e, funcName, errorList)

		} else if col == "AI_intern_DC_LinkVoltage2" {
			data.AI_intern_DC_LinkVoltage2, e = GetFloatCell(row.Cells[idx])
			ErrorLog(e, funcName, errorList)

		} else if col == "AI_intern_DC_LinkVoltage3" {
			data.AI_intern_DC_LinkVoltage3, e = GetFloatCell(row.Cells[idx])
			ErrorLog(e, funcName, errorList)

		} else if col == "Temp_Yaw_Motor1" {
			data.Temp_Yaw_Motor1, e = GetFloatCell(row.Cells[idx])
			ErrorLog(e, funcName, errorList)

		} else if col == "Temp_Yaw_Motor2" {
			data.Temp_Yaw_Motor2, e = GetFloatCell(row.Cells[idx])
			ErrorLog(e, funcName, errorList)

		} else if col == "Temp_Yaw_Motor3" {
			data.Temp_Yaw_Motor3, e = GetFloatCell(row.Cells[idx])
			ErrorLog(e, funcName, errorList)

		} else if col == "Temp_Yaw_Motor4" {
			data.Temp_Yaw_Motor4, e = GetFloatCell(row.Cells[idx])
			ErrorLog(e, funcName, errorList)

		} else if col == "AO_DFIG_Power_Setpiont" {
			data.AO_DFIG_Power_Setpiont, e = GetFloatCell(row.Cells[idx])
			ErrorLog(e, funcName, errorList)

		} else if col == "AO_DFIG_Q_Setpoint" {
			data.AO_DFIG_Q_Setpoint, e = GetFloatCell(row.Cells[idx])
			ErrorLog(e, funcName, errorList)

		} else if col == "AI_DFIG_Torque_actual" {
			data.AI_DFIG_Torque_actual, e = GetFloatCell(row.Cells[idx])
			ErrorLog(e, funcName, errorList)

		} else if col == "AI_DFIG_SpeedGenerator_Encoder" {
			data.AI_DFIG_SpeedGenerator_Encoder, e = GetFloatCell(row.Cells[idx])
			ErrorLog(e, funcName, errorList)

		} else if col == "AI_intern_DFIG_DC_Link_Voltage_actual" {
			data.AI_intern_DFIG_DC_Link_Voltage_actual, e = GetFloatCell(row.Cells[idx])
			ErrorLog(e, funcName, errorList)

		} else if col == "AI_intern_DFIG_MSC_current" {
			data.AI_intern_DFIG_MSC_current, e = GetFloatCell(row.Cells[idx])
			ErrorLog(e, funcName, errorList)

		} else if col == "AI_intern_DFIG_Main_voltage" {
			data.AI_intern_DFIG_Main_voltage, e = GetFloatCell(row.Cells[idx])
			ErrorLog(e, funcName, errorList)

		} else if col == "AI_intern_DFIG_Main_current" {
			data.AI_intern_DFIG_Main_current, e = GetFloatCell(row.Cells[idx])
			ErrorLog(e, funcName, errorList)

		} else if col == "AI_intern_DFIG_active_power_actual" {
			data.AI_intern_DFIG_active_power_actual, e = GetFloatCell(row.Cells[idx])
			ErrorLog(e, funcName, errorList)

		} else if col == "AI_intern_DFIG_reactive_power_actual" {
			data.AI_intern_DFIG_reactive_power_actual, e = GetFloatCell(row.Cells[idx])
			ErrorLog(e, funcName, errorList)

		} else if col == "AI_intern_DFIG_active_power_actual_LSC" {
			data.AI_intern_DFIG_active_power_actual_LSC, e = GetFloatCell(row.Cells[idx])
			ErrorLog(e, funcName, errorList)

		} else if col == "AI_intern_DFIG_LSC_current" {
			data.AI_intern_DFIG_LSC_current, e = GetFloatCell(row.Cells[idx])
			ErrorLog(e, funcName, errorList)

		} else if col == "AI_intern_DFIG_Data_log_number" {
			data.AI_intern_DFIG_Data_log_number, e = GetFloatCell(row.Cells[idx])
			ErrorLog(e, funcName, errorList)

		} else if col == "AI_intern_Damper_OscMagnitude" {
			data.AI_intern_Damper_OscMagnitude, e = GetFloatCell(row.Cells[idx])
			ErrorLog(e, funcName, errorList)

		} else if col == "AI_intern_Damper_PassbandFullLoad" {
			data.AI_intern_Damper_PassbandFullLoad, e = GetFloatCell(row.Cells[idx])
			ErrorLog(e, funcName, errorList)

		} else if col == "AI_YawBrake_TempRise1" {
			data.AI_YawBrake_TempRise1, e = GetFloatCell(row.Cells[idx])
			ErrorLog(e, funcName, errorList)

		} else if col == "AI_YawBrake_TempRise2" {
			data.AI_YawBrake_TempRise2, e = GetFloatCell(row.Cells[idx])
			ErrorLog(e, funcName, errorList)

		} else if col == "AI_YawBrake_TempRise3" {
			data.AI_YawBrake_TempRise3, e = GetFloatCell(row.Cells[idx])
			ErrorLog(e, funcName, errorList)

		} else if col == "AI_YawBrake_TempRise4" {
			data.AI_YawBrake_TempRise4, e = GetFloatCell(row.Cells[idx])
			ErrorLog(e, funcName, errorList)

		} else if col == "AI_intern_NacelleDrill_at_NorthPosSensor" {
			data.AI_intern_NacelleDrill_at_NorthPosSensor, e = GetFloatCell(row.Cells[idx])
			ErrorLog(e, funcName, errorList)

		}
	}

	res = data
	return
}
