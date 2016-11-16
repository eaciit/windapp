package controllers

import (
	"bufio"
	. "github.com/eaciit/windapp/library/helper"
	. "github.com/eaciit/windapp/library/models"
	"encoding/csv"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	_ "github.com/eaciit/dbox/dbc/mongo"
	tk "github.com/eaciit/toolkit"
)

type ConvScadaTreeSecs struct {
	*BaseController
}

func (d *ConvScadaTreeSecs) Generate(base *BaseController, file string) (errorLine tk.M) {
	funcName := "ConvScadaTreeSecs"
	if base != nil {
		d.BaseController = base
		ctx := d.BaseController.Ctx
		fr, _ := os.Open(file)

		anFile := strings.Split(file, "/")
		fileName := anFile[len(anFile)-1]

		// tk.Printf("file: %v \n", file)

		read := csv.NewReader(bufio.NewReader(fr))
		count := 0

		for {
			record, err := read.Read()
			// tk.Printf("record: %#v \n", record)
			if count != 0 {
				errorList := []error{}

				if err == io.EOF {
					// tk.Printf("ERR: %#v \n", err.Error())
					break
				}

				scada := new(ScadaThreeSecs)
				scada.File = fileName
				scada.Line = count

				scada.TimeStamp1, err = time.Parse("2006-01-02 15:04:05", string(record[0]))
				ErrorLog(err, funcName, errorList)
				scada.DateId1, err = time.Parse("2006-01-02 15:04:05", string(record[1]))
				ErrorLog(err, funcName, errorList)
				scada.TimeStamp2, err = time.Parse("2006-01-02 15:04:05", string(record[2]))
				ErrorLog(err, funcName, errorList)
				scada.DateId2, err = time.Parse("2006-01-02 15:04:05", string(record[3]))
				ErrorLog(err, funcName, errorList)

				scada.DateId1Info = GetDateInfo(scada.TimeStamp1)
				scada.DateId2Info = GetDateInfo(scada.TimeStamp2)

				timeStamp := scada.TimeStamp1.UTC()
				seconds := tk.Div(tk.ToFloat64(timeStamp.Nanosecond(), 1, tk.RoundingAuto), 1000000000)
				secondsInt := tk.ToInt(seconds, tk.RoundingAuto)
				newTimeTmp := timeStamp.Add(time.Duration(secondsInt) * time.Second)
				strTime := tk.ToString(newTimeTmp.Year()) + tk.ToString(int(newTimeTmp.Month())) + tk.ToString(newTimeTmp.Day()) + " " + tk.ToString(newTimeTmp.Hour()) + ":" + tk.ToString(newTimeTmp.Minute()) + ":" + tk.ToString(newTimeTmp.Second())
				newTime, _ := time.Parse("200612 15:4:5", strTime)

				scada.TimeStampSecondGroup = newTime

				tenMinuteInfo := GenTenMinuteInfo(newTime)

				scada.THour = tenMinuteInfo.THour
				scada.TMinute = tenMinuteInfo.TMinute
				scada.TSecond = tenMinuteInfo.TSecond
				scada.TMinuteValue = tenMinuteInfo.TMinuteValue
				scada.TMinuteCategory = tenMinuteInfo.TMinuteCategory
				scada.TimeStampConverted = tenMinuteInfo.TimeStampConverted
				scada.TimeStampConvertedInt, _ = strconv.ParseInt(scada.TimeStampConverted.Format("200601021504"), 10, 64)

				scada.ProjectName = string(record[4])
				scada.Turbine = string(record[5])

				scada = scada.New()

				scada.Fast_CurrentL3 = SetFloatValue(record[6])
				scada.Fast_ActivePower_kW = SetFloatValue(record[7])
				scada.Fast_CurrentL1 = SetFloatValue(record[8])
				scada.Fast_ActivePowerSetpoint_kW = SetFloatValue(record[9])
				scada.Fast_CurrentL2 = SetFloatValue(record[10])
				scada.Fast_DrTrVibValue = SetFloatValue(record[11])
				scada.Fast_GenSpeed_RPM = SetFloatValue(record[12])
				scada.Fast_PitchAccuV1 = SetFloatValue(record[13])
				scada.Fast_PitchAngle = SetFloatValue(record[14])
				scada.Fast_PitchAngle3 = SetFloatValue(record[15])
				scada.Fast_PitchAngle2 = SetFloatValue(record[16])
				scada.Fast_PitchConvCurrent1 = SetFloatValue(record[17])
				scada.Fast_PitchConvCurrent3 = SetFloatValue(record[18])
				scada.Fast_PitchConvCurrent2 = SetFloatValue(record[19])
				scada.Fast_PowerFactor = SetFloatValue(record[20])
				scada.Fast_ReactivePowerSetpointPPC_kVAr = SetFloatValue(record[21])
				scada.Fast_ReactivePower_kVAr = SetFloatValue(record[22])
				scada.Fast_RotorSpeed_RPM = SetFloatValue(record[23])
				scada.Fast_VoltageL1 = SetFloatValue(record[24])
				scada.Fast_VoltageL2 = SetFloatValue(record[25])
				scada.Fast_WindSpeed_ms = SetFloatValue(record[26])
				scada.Slow_CapableCapacitiveReactPwr_kVAr = SetFloatValue(record[27])
				scada.Slow_CapableInductiveReactPwr_kVAr = SetFloatValue(record[28])
				scada.Slow_DateTime_Sec = SetFloatValue(record[29])
				scada.Slow_NacellePos = SetFloatValue(record[30])
				scada.Fast_PitchAngle1 = SetFloatValue(record[31])
				scada.Fast_VoltageL3 = SetFloatValue(record[32])
				scada.Slow_CapableCapacitivePwrFactor = SetFloatValue(record[33])
				scada.Fast_Total_Production_kWh = SetFloatValue(record[34])
				scada.Fast_Total_Prod_Day_kWh = SetFloatValue(record[35])
				scada.Fast_Total_Prod_Month_kWh = SetFloatValue(record[36])
				scada.Fast_ActivePowerOutPWCSell_kW = SetFloatValue(record[37])
				scada.Fast_Frequency_Hz = SetFloatValue(record[38])
				scada.Slow_TempG1L2 = SetFloatValue(record[39])
				scada.Slow_TempG1L3 = SetFloatValue(record[40])
				scada.Slow_TempGearBoxHSSDE = SetFloatValue(record[41])
				scada.Slow_TempGearBoxIMSNDE = SetFloatValue(record[42])
				scada.Slow_TempOutdoor = SetFloatValue(record[43])
				scada.Fast_PitchAccuV3 = SetFloatValue(record[44])
				scada.Slow_TotalTurbineActiveHours = SetFloatValue(record[45])
				scada.Slow_TotalTurbineOKHours = SetFloatValue(record[46])
				scada.Slow_TotalTurbineTimeAllHours = SetFloatValue(record[47])
				scada.Slow_TempG1L1 = SetFloatValue(record[48])
				scada.Slow_TempGearBoxOilSump = SetFloatValue(record[49])
				scada.Fast_PitchAccuV2 = SetFloatValue(record[50])
				scada.Slow_TotalGridOkHours = SetFloatValue(record[51])
				scada.Slow_TotalActPowerOut_kWh = SetFloatValue(record[52])
				scada.Fast_YawService = SetFloatValue(record[53])
				scada.Fast_YawAngle = SetFloatValue(record[54])
				scada.Slow_WindDirection = SetFloatValue(record[55])
				scada.Slow_CapableInductivePwrFactor = SetFloatValue(record[56])
				scada.Slow_TempGearBoxHSSNDE = SetFloatValue(record[57])
				scada.Slow_TempHubBearing = SetFloatValue(record[58])
				scada.Slow_TotalG1ActiveHours = SetFloatValue(record[59])
				scada.Slow_TotalActPowerOutG1_kWh = SetFloatValue(record[60])
				scada.Slow_TotalReactPowerInG1_kVArh = SetFloatValue(record[61])
				scada.Slow_NacelleDrill = SetFloatValue(record[62])
				scada.Slow_TempGearBoxIMSDE = SetFloatValue(record[63])
				scada.Fast_Total_Operating_hrs = SetFloatValue(record[64])
				scada.Slow_TempNacelle = SetFloatValue(record[65])
				scada.Fast_Total_Grid_OK_hrs = SetFloatValue(record[66])
				scada.Fast_Total_WTG_OK_hrs = SetFloatValue(record[67])
				scada.Slow_TempCabinetTopBox = SetFloatValue(record[68])
				scada.Slow_TempGeneratorBearingNDE = SetFloatValue(record[69])
				scada.Fast_Total_Access_hrs = SetFloatValue(record[70])
				scada.Slow_TempBottomPowerSection = SetFloatValue(record[71])
				scada.Slow_TempGeneratorBearingDE = SetFloatValue(record[72])
				scada.Slow_TotalReactPowerIn_kVArh = SetFloatValue(record[73])
				scada.Slow_TempBottomControlSection = SetFloatValue(record[74])
				scada.Slow_TempConv1 = SetFloatValue(record[75])
				scada.Fast_ActivePowerRated_kW = SetFloatValue(record[76])
				scada.Fast_NodeIP = SetFloatValue(record[77])
				scada.Fast_PitchSpeed1 = SetFloatValue(record[78])
				scada.Slow_CFCardSize = SetFloatValue(record[79])
				scada.Slow_CPU_Number = SetFloatValue(record[80])
				scada.Slow_CFCardSpaceLeft = SetFloatValue(record[81])
				scada.Slow_TempBottomCapSection = SetFloatValue(record[82])
				scada.Slow_RatedPower = SetFloatValue(record[83])
				scada.Slow_TempConv3 = SetFloatValue(record[84])
				scada.Slow_TempConv2 = SetFloatValue(record[85])
				scada.Slow_TotalActPowerIn_kWh = SetFloatValue(record[86])
				scada.Slow_TotalActPowerInG1_kWh = SetFloatValue(record[87])
				scada.Slow_TotalActPowerInG2_kWh = SetFloatValue(record[88])
				scada.Slow_TotalActPowerOutG2_kWh = SetFloatValue(record[89])
				scada.Slow_TotalG2ActiveHours = SetFloatValue(record[90])
				scada.Slow_TotalReactPowerInG2_kVArh = SetFloatValue(record[91])
				scada.Slow_TotalReactPowerOut_kVArh = SetFloatValue(record[92])
				scada.Slow_UTCoffset_int = SetFloatValue(record[93])

				// log.Printf("scada: %#v \n", scada)

				if len(errorList) == 0 {
					err = ctx.Insert(scada)
					ErrorLog(err, funcName, errorList)
					ErrorHandler(err, "Saving")
				}

				if len(errorList) > 0 {
					errorLine.Set(tk.ToString(scada.Line), errorList)
				}
			}

			count++
		}
	}

	return
}
