package controllers

import (
	"bufio"
	. "eaciit/wfdemo-git/library/helper"
	. "eaciit/wfdemo-git/library/models"
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
		if string(os.PathSeparator)=="\\"{
			anFile = strings.Split(file, "\\")
		}
		fileName := anFile[len(anFile)-1]

		// tk.Printf("file: %v \n", file)

		read := csv.NewReader(bufio.NewReader(fr))
		count := 0
		fieldToIndex := map[string]int{}
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

				scada.TimeStamp1, err = time.Parse("2006-01-02 15:04:05", string(record[fieldToIndex["TimeStamp1"]]))
				ErrorLog(err, funcName, errorList)
				scada.DateId1, err = time.Parse("2006-01-02 15:04:05", string(record[fieldToIndex["DateId1"]]))
				ErrorLog(err, funcName, errorList)
				scada.TimeStamp2, err = time.Parse("2006-01-02 15:04:05", string(record[fieldToIndex["TimeStamp2"]]))
				ErrorLog(err, funcName, errorList)
				scada.DateId2, err = time.Parse("2006-01-02 15:04:05", string(record[fieldToIndex["DateId2"]]))
				ErrorLog(err, funcName, errorList)
				scada.TimeStampConverted,err = time.Parse("2006-01-02 15:04:05", string(record[fieldToIndex["TimeStampConverted"]]))
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

				//tenMinuteInfo := GenTenMinuteInfo(newTime)

				scada.THour,_ = strconv.Atoi(record[fieldToIndex["THour"]])
				scada.TMinute,_ = strconv.Atoi(record[fieldToIndex["TMinute"]])
				scada.TSecond,_ = strconv.Atoi(record[fieldToIndex["TSecond"]])
				scada.TMinuteValue = SetFloatValue(record[fieldToIndex["TMinuteValue"]])
				scada.TMinuteCategory,_ = strconv.Atoi(record[fieldToIndex["TMinuteCategory"]])
				scada.TimeStampConvertedInt, _ = strconv.ParseInt(scada.TimeStampConverted.Format("200601021504"), 10, 64)

				scada.ProjectName = string(record[fieldToIndex["ProjectName"]])
				scada.Turbine = string(record[fieldToIndex["Turbine"]])

				scada = scada.New()

				scada.Fast_CurrentL3 = SetFloatValue(record[fieldToIndex["Fast_CurrentL3"]])
				scada.Fast_ActivePower_kW = SetFloatValue(record[fieldToIndex["Fast_ActivePower_kW"]])
				scada.Fast_CurrentL1 = SetFloatValue(record[fieldToIndex["Fast_CurrentL1"]])
				scada.Fast_ActivePowerSetpoint_kW = SetFloatValue(record[fieldToIndex["Fast_ActivePowerSetpoint_kW"]])
				scada.Fast_CurrentL2 = SetFloatValue(record[fieldToIndex["Fast_CurrentL2"]])
				scada.Fast_DrTrVibValue = SetFloatValue(record[fieldToIndex["Fast_DrTrVibValue"]])
				scada.Fast_GenSpeed_RPM = SetFloatValue(record[fieldToIndex["Fast_GenSpeed_RPM"]])
				scada.Fast_PitchAccuV1 = SetFloatValue(record[fieldToIndex["Fast_PitchAccuV1"]])
				scada.Fast_PitchAngle = SetFloatValue(record[fieldToIndex["Fast_PitchAngle"]])
				scada.Fast_PitchAngle3 = SetFloatValue(record[fieldToIndex["Fast_PitchAngle3"]])
				scada.Fast_PitchAngle2 = SetFloatValue(record[fieldToIndex["Fast_PitchAngle2"]])
				scada.Fast_PitchConvCurrent1 = SetFloatValue(record[fieldToIndex["Fast_PitchConvCurrent1"]])
				scada.Fast_PitchConvCurrent3 = SetFloatValue(record[fieldToIndex["Fast_PitchConvCurrent3"]])
				scada.Fast_PitchConvCurrent2 = SetFloatValue(record[fieldToIndex["Fast_PitchConvCurrent2"]])
				scada.Fast_PowerFactor = SetFloatValue(record[fieldToIndex["Fast_PowerFactor"]])
				scada.Fast_ReactivePowerSetpointPPC_kVAr = SetFloatValue(record[fieldToIndex["Fast_ReactivePowerSetpointPPC_kVAr"]])
				scada.Fast_ReactivePower_kVAr = SetFloatValue(record[fieldToIndex["Fast_ReactivePower_kVAr"]])
				scada.Fast_RotorSpeed_RPM = SetFloatValue(record[fieldToIndex["Fast_RotorSpeed_RPM"]])
				scada.Fast_VoltageL1 = SetFloatValue(record[fieldToIndex["Fast_VoltageL1"]])
				scada.Fast_VoltageL2 = SetFloatValue(record[fieldToIndex["Fast_VoltageL2"]])
				scada.Fast_WindSpeed_ms = SetFloatValue(record[fieldToIndex["Fast_WindSpeed_ms"]])
				scada.Slow_CapableCapacitiveReactPwr_kVAr = SetFloatValue(record[fieldToIndex["Slow_CapableCapacitiveReactPwr_kVAr"]])
				scada.Slow_CapableInductiveReactPwr_kVAr = SetFloatValue(record[fieldToIndex["Slow_CapableInductiveReactPwr_kVAr"]])
				scada.Slow_DateTime_Sec = SetFloatValue(record[fieldToIndex["Slow_DateTime_Sec"]])
				scada.Slow_NacellePos = SetFloatValue(record[fieldToIndex["Slow_NacellePos"]])
				scada.Fast_PitchAngle1 = SetFloatValue(record[fieldToIndex["Fast_PitchAngle1"]])
				scada.Fast_VoltageL3 = SetFloatValue(record[fieldToIndex["Fast_VoltageL3"]])
				scada.Slow_CapableCapacitivePwrFactor = SetFloatValue(record[fieldToIndex["Slow_CapableCapacitivePwrFactor"]])
				scada.Fast_Total_Production_kWh = SetFloatValue(record[fieldToIndex["Fast_Total_Production_kWh"]])
				scada.Fast_Total_Prod_Day_kWh = SetFloatValue(record[fieldToIndex["Fast_Total_Prod_Day_kWh"]])
				scada.Fast_Total_Prod_Month_kWh = SetFloatValue(record[fieldToIndex["Fast_Total_Prod_Month_kWh"]])
				scada.Fast_ActivePowerOutPWCSell_kW = SetFloatValue(record[fieldToIndex["Fast_ActivePowerOutPWCSell_kW"]])
				scada.Fast_Frequency_Hz = SetFloatValue(record[fieldToIndex["Fast_Frequency_Hz"]])
				scada.Slow_TempG1L2 = SetFloatValue(record[fieldToIndex["Slow_TempG1L2"]])
				scada.Slow_TempG1L3 = SetFloatValue(record[fieldToIndex["Slow_TempG1L3"]])
				scada.Slow_TempGearBoxHSSDE = SetFloatValue(record[fieldToIndex["Slow_TempGearBoxHSSDE"]])
				scada.Slow_TempGearBoxIMSNDE = SetFloatValue(record[fieldToIndex["Slow_TempGearBoxIMSNDE"]])
				scada.Slow_TempOutdoor = SetFloatValue(record[fieldToIndex["Slow_TempOutdoor"]])
				scada.Fast_PitchAccuV3 = SetFloatValue(record[fieldToIndex["Fast_PitchAccuV3"]])
				scada.Slow_TotalTurbineActiveHours = SetFloatValue(record[fieldToIndex["Slow_TotalTurbineActiveHours"]])
				scada.Slow_TotalTurbineOKHours = SetFloatValue(record[fieldToIndex["Slow_TotalTurbineOKHours"]])
				scada.Slow_TotalTurbineTimeAllHours = SetFloatValue(record[fieldToIndex["Slow_TotalTurbineTimeAllHours"]])
				scada.Slow_TempG1L1 = SetFloatValue(record[fieldToIndex["Slow_TempG1L1"]])
				scada.Slow_TempGearBoxOilSump = SetFloatValue(record[fieldToIndex["Slow_TempGearBoxOilSump"]])
				scada.Fast_PitchAccuV2 = SetFloatValue(record[fieldToIndex["Fast_PitchAccuV2"]])
				scada.Slow_TotalGridOkHours = SetFloatValue(record[fieldToIndex["Slow_TotalGridOkHours"]])
				scada.Slow_TotalActPowerOut_kWh = SetFloatValue(record[fieldToIndex["Slow_TotalActPowerOut_kWh"]])
				scada.Fast_YawService = SetFloatValue(record[fieldToIndex["Fast_YawService"]])
				scada.Fast_YawAngle = SetFloatValue(record[fieldToIndex["Fast_YawAngle"]])
				scada.Slow_WindDirection = SetFloatValue(record[fieldToIndex["Slow_WindDirection"]])
				scada.Slow_CapableInductivePwrFactor = SetFloatValue(record[fieldToIndex["Slow_CapableInductivePwrFactor"]])
				scada.Slow_TempGearBoxHSSNDE = SetFloatValue(record[fieldToIndex["Slow_TempGearBoxHSSNDE"]])
				scada.Slow_TempHubBearing = SetFloatValue(record[fieldToIndex["Slow_TempHubBearing"]])
				scada.Slow_TotalG1ActiveHours = SetFloatValue(record[fieldToIndex["Slow_TotalG1ActiveHours"]])
				scada.Slow_TotalActPowerOutG1_kWh = SetFloatValue(record[fieldToIndex["Slow_TotalActPowerOutG1_kWh"]])
				scada.Slow_TotalReactPowerInG1_kVArh = SetFloatValue(record[fieldToIndex["Slow_TotalReactPowerInG1_kVArh"]])
				scada.Slow_NacelleDrill = SetFloatValue(record[fieldToIndex["Slow_NacelleDrill"]])
				scada.Slow_TempGearBoxIMSDE = SetFloatValue(record[fieldToIndex["Slow_TempGearBoxIMSDE"]])
				scada.Fast_Total_Operating_hrs = SetFloatValue(record[fieldToIndex["Fast_Total_Operating_hrs"]])
				scada.Slow_TempNacelle = SetFloatValue(record[fieldToIndex["Slow_TempNacelle"]])
				scada.Fast_Total_Grid_OK_hrs = SetFloatValue(record[fieldToIndex["Fast_Total_Grid_OK_hrs"]])
				scada.Fast_Total_WTG_OK_hrs = SetFloatValue(record[fieldToIndex["Fast_Total_WTG_OK_hrs"]])
				scada.Slow_TempCabinetTopBox = SetFloatValue(record[fieldToIndex["Slow_TempCabinetTopBox"]])
				scada.Slow_TempGeneratorBearingNDE = SetFloatValue(record[fieldToIndex["Slow_TempGeneratorBearingNDE"]])
				scada.Fast_Total_Access_hrs = SetFloatValue(record[fieldToIndex["Fast_Total_Access_hrs"]])
				scada.Slow_TempBottomPowerSection = SetFloatValue(record[fieldToIndex["Slow_TempBottomPowerSection"]])
				scada.Slow_TempGeneratorBearingDE = SetFloatValue(record[fieldToIndex["Slow_TempGeneratorBearingDE"]])
				scada.Slow_TotalReactPowerIn_kVArh = SetFloatValue(record[fieldToIndex["Slow_TotalReactPowerIn_kVArh"]])
				scada.Slow_TempBottomControlSection = SetFloatValue(record[fieldToIndex["Slow_TempBottomControlSection"]])
				scada.Slow_TempConv1 = SetFloatValue(record[fieldToIndex["Slow_TempConv1"]])
				scada.Fast_ActivePowerRated_kW = SetFloatValue(record[fieldToIndex["Fast_ActivePowerRated_kW"]])
				scada.Fast_NodeIP = SetFloatValue(record[fieldToIndex["Fast_NodeIP"]])
				scada.Fast_PitchSpeed1 = SetFloatValue(record[fieldToIndex["Fast_PitchSpeed1"]])
				scada.Slow_CFCardSize = SetFloatValue(record[fieldToIndex["Slow_CFCardSize"]])
				scada.Slow_CPU_Number = SetFloatValue(record[fieldToIndex["Slow_CPU_Number"]])
				scada.Slow_CFCardSpaceLeft = SetFloatValue(record[fieldToIndex["Slow_CFCardSpaceLeft"]])
				scada.Slow_TempBottomCapSection = SetFloatValue(record[fieldToIndex["Slow_TempBottomCapSection"]])
				scada.Slow_RatedPower = SetFloatValue(record[fieldToIndex["Slow_RatedPower"]])
				scada.Slow_TempConv3 = SetFloatValue(record[fieldToIndex["Slow_TempConv3"]])
				scada.Slow_TempConv2 = SetFloatValue(record[fieldToIndex["Slow_TempConv2"]])
				scada.Slow_TotalActPowerIn_kWh = SetFloatValue(record[fieldToIndex["Slow_TotalActPowerIn_kWh"]])
				scada.Slow_TotalActPowerInG1_kWh = SetFloatValue(record[fieldToIndex["Slow_TotalActPowerInG1_kWh"]])
				scada.Slow_TotalActPowerInG2_kWh = SetFloatValue(record[fieldToIndex["Slow_TotalActPowerInG2_kWh"]])
				scada.Slow_TotalActPowerOutG2_kWh = SetFloatValue(record[fieldToIndex["Slow_TotalActPowerOutG2_kWh"]])
				scada.Slow_TotalG2ActiveHours = SetFloatValue(record[fieldToIndex["Slow_TotalG2ActiveHours"]])
				scada.Slow_TotalReactPowerInG2_kVArh = SetFloatValue(record[fieldToIndex["Slow_TotalReactPowerInG2_kVArh"]])
				scada.Slow_TotalReactPowerOut_kVArh = SetFloatValue(record[fieldToIndex["Slow_TotalReactPowerOut_kVArh"]])
				scada.Slow_UTCoffset_int = SetFloatValue(record[fieldToIndex["Slow_UTCoffset_int"]])

				// log.Printf("scada: %#v \n", scada)

				if len(errorList) == 0 {
					//tk.Println("Adding")
					err = ctx.Insert(scada)
					ErrorLog(err, funcName, errorList)
					ErrorHandler(err, "Saving")
				}

				if len(errorList) > 0 {
					errorLine.Set(tk.ToString(scada.Line), errorList)
				}
			}else{
				for idx,val:=range record{
					
					fieldToIndex[val]=idx
				}				
			}

			count++
		}
	}

	return
}
