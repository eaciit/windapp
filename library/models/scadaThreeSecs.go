package models

import (
	. "eaciit/wfdemo/library/helper"
	"time"

	"github.com/eaciit/orm"
	tk "github.com/eaciit/toolkit"
)

type ScadaThreeSecs struct {
	orm.ModelBase                       `bson:"-",json:"-"`
	ID                                  string ` bson:"_id" , json:"_id" `
	TimeStamp1                          time.Time
	DateId1                             time.Time
	TimeStamp2                          time.Time
	DateId2                             time.Time
	DateId1Info                         DateInfo
	DateId2Info                         DateInfo
	ProjectName                         string
	Turbine                             string
	THour                               int
	TMinute                             int
	TSecond                             int
	TMinuteValue                        float64
	TMinuteCategory                     int
	TimeStampConverted                  time.Time
	TimeStampConvertedInt               int64
	TimeStampSecondGroup                time.Time
	Fast_CurrentL3                      float64
	Fast_ActivePower_kW                 float64
	Fast_CurrentL1                      float64
	Fast_ActivePowerSetpoint_kW         float64
	Fast_CurrentL2                      float64
	Fast_DrTrVibValue                   float64
	Fast_GenSpeed_RPM                   float64
	Fast_PitchAccuV1                    float64
	Fast_PitchAngle                     float64
	Fast_PitchAngle3                    float64
	Fast_PitchAngle2                    float64
	Fast_PitchConvCurrent1              float64
	Fast_PitchConvCurrent3              float64
	Fast_PitchConvCurrent2              float64
	Fast_PowerFactor                    float64
	Fast_ReactivePowerSetpointPPC_kVAr  float64
	Fast_ReactivePower_kVAr             float64
	Fast_RotorSpeed_RPM                 float64
	Fast_VoltageL1                      float64
	Fast_VoltageL2                      float64
	Fast_WindSpeed_ms                   float64
	Slow_CapableCapacitiveReactPwr_kVAr float64
	Slow_CapableInductiveReactPwr_kVAr  float64
	Slow_DateTime_Sec                   float64
	Slow_NacellePos                     float64
	Fast_PitchAngle1                    float64
	Fast_VoltageL3                      float64
	Slow_CapableCapacitivePwrFactor     float64
	Fast_Total_Production_kWh           float64
	Fast_Total_Prod_Day_kWh             float64
	Fast_Total_Prod_Month_kWh           float64
	Fast_ActivePowerOutPWCSell_kW       float64
	Fast_Frequency_Hz                   float64
	Slow_TempG1L2                       float64
	Slow_TempG1L3                       float64
	Slow_TempGearBoxHSSDE               float64
	Slow_TempGearBoxIMSNDE              float64
	Slow_TempOutdoor                    float64
	Fast_PitchAccuV3                    float64
	Slow_TotalTurbineActiveHours        float64
	Slow_TotalTurbineOKHours            float64
	Slow_TotalTurbineTimeAllHours       float64
	Slow_TempG1L1                       float64
	Slow_TempGearBoxOilSump             float64
	Fast_PitchAccuV2                    float64
	Slow_TotalGridOkHours               float64
	Slow_TotalActPowerOut_kWh           float64
	Fast_YawService                     float64
	Fast_YawAngle                       float64
	Slow_WindDirection                  float64
	Slow_CapableInductivePwrFactor      float64
	Slow_TempGearBoxHSSNDE              float64
	Slow_TempHubBearing                 float64
	Slow_TotalG1ActiveHours             float64
	Slow_TotalActPowerOutG1_kWh         float64
	Slow_TotalReactPowerInG1_kVArh      float64
	Slow_NacelleDrill                   float64
	Slow_TempGearBoxIMSDE               float64
	Fast_Total_Operating_hrs            float64
	Slow_TempNacelle                    float64
	Fast_Total_Grid_OK_hrs              float64
	Fast_Total_WTG_OK_hrs               float64
	Slow_TempCabinetTopBox              float64
	Slow_TempGeneratorBearingNDE        float64
	Fast_Total_Access_hrs               float64
	Slow_TempBottomPowerSection         float64
	Slow_TempGeneratorBearingDE         float64
	Slow_TotalReactPowerIn_kVArh        float64
	Slow_TempBottomControlSection       float64
	Slow_TempConv1                      float64
	Fast_ActivePowerRated_kW            float64
	Fast_NodeIP                         float64
	Fast_PitchSpeed1                    float64
	Slow_CFCardSize                     float64
	Slow_CPU_Number                     float64
	Slow_CFCardSpaceLeft                float64
	Slow_TempBottomCapSection           float64
	Slow_RatedPower                     float64
	Slow_TempConv3                      float64
	Slow_TempConv2                      float64
	Slow_TotalActPowerIn_kWh            float64
	Slow_TotalActPowerInG1_kWh          float64
	Slow_TotalActPowerInG2_kWh          float64
	Slow_TotalActPowerOutG2_kWh         float64
	Slow_TotalG2ActiveHours             float64
	Slow_TotalReactPowerInG2_kVArh      float64
	Slow_TotalReactPowerOut_kVArh       float64
	Slow_UTCoffset_int                  float64
	Line                                int
	File                                string
}

func (m *ScadaThreeSecs) New() *ScadaThreeSecs {
	timeStampStr := m.TimeStamp1.Format("060102_150405")
	m.ID = timeStampStr + "#" + m.ProjectName + "#" + m.Turbine + "#" + tk.ToString(m.Line)
	return m
}

func (m *ScadaThreeSecs) RecordID() interface{} {
	return m.ID
}

func (m *ScadaThreeSecs) TableName() string {
	return "ScadaThreeSecs"
}
