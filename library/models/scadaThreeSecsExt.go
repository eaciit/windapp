package models

import (
	// . "eaciit/ostrowfm/library/helper"
	"time"

	"github.com/eaciit/orm"
)

type ScadaThreeSecsExt struct {
	orm.ModelBase `bson:"-",json:"-"`
	ID            string ` bson:"_id" , json:"_id" `
	/*TimeStamp1                          time.Time
	DateId1                             time.Time
	TimeStamp2                          time.Time
	DateId2                             time.Time
	DateId1Info                         DateInfo
	DateId2Info                         DateInfo*/
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

	Fast_CurrentL3CountSecs                      int
	Fast_ActivePower_kWCountSecs                 int
	Fast_CurrentL1CountSecs                      int
	Fast_ActivePowerSetpoint_kWCountSecs         int
	Fast_CurrentL2CountSecs                      int
	Fast_DrTrVibValueCountSecs                   int
	Fast_GenSpeed_RPMCountSecs                   int
	Fast_PitchAccuV1CountSecs                    int
	Fast_PitchAngleCountSecs                     int
	Fast_PitchAngle3CountSecs                    int
	Fast_PitchAngle2CountSecs                    int
	Fast_PitchConvCurrent1CountSecs              int
	Fast_PitchConvCurrent3CountSecs              int
	Fast_PitchConvCurrent2CountSecs              int
	Fast_PowerFactorCountSecs                    int
	Fast_ReactivePowerSetpointPPC_kVArCountSecs  int
	Fast_ReactivePower_kVArCountSecs             int
	Fast_RotorSpeed_RPMCountSecs                 int
	Fast_VoltageL1CountSecs                      int
	Fast_VoltageL2CountSecs                      int
	Fast_WindSpeed_msCountSecs                   int
	Slow_CapableCapacitiveReactPwr_kVArCountSecs int
	Slow_CapableInductiveReactPwr_kVArCountSecs  int
	Slow_DateTime_SecCountSecs                   int
	Slow_NacellePosCountSecs                     int
	Fast_PitchAngle1CountSecs                    int
	Fast_VoltageL3CountSecs                      int
	Slow_CapableCapacitivePwrFactorCountSecs     int
	Fast_Total_Production_kWhCountSecs           int
	Fast_Total_Prod_Day_kWhCountSecs             int
	Fast_Total_Prod_Month_kWhCountSecs           int
	Fast_ActivePowerOutPWCSell_kWCountSecs       int
	Fast_Frequency_HzCountSecs                   int
	Slow_TempG1L2CountSecs                       int
	Slow_TempG1L3CountSecs                       int
	Slow_TempGearBoxHSSDECountSecs               int
	Slow_TempGearBoxIMSNDECountSecs              int
	Slow_TempOutdoorCountSecs                    int
	Fast_PitchAccuV3CountSecs                    int
	Slow_TotalTurbineActiveHoursCountSecs        int
	Slow_TotalTurbineOKHoursCountSecs            int
	Slow_TotalTurbineTimeAllHoursCountSecs       int
	Slow_TempG1L1CountSecs                       int
	Slow_TempGearBoxOilSumpCountSecs             int
	Fast_PitchAccuV2CountSecs                    int
	Slow_TotalGridOkHoursCountSecs               int
	Slow_TotalActPowerOut_kWhCountSecs           int
	Fast_YawServiceCountSecs                     int
	Fast_YawAngleCountSecs                       int
	Slow_WindDirectionCountSecs                  int
	Slow_CapableInductivePwrFactorCountSecs      int
	Slow_TempGearBoxHSSNDECountSecs              int
	Slow_TempHubBearingCountSecs                 int
	Slow_TotalG1ActiveHoursCountSecs             int
	Slow_TotalActPowerOutG1_kWhCountSecs         int
	Slow_TotalReactPowerInG1_kVArhCountSecs      int
	Slow_NacelleDrillCountSecs                   int
	Slow_TempGearBoxIMSDECountSecs               int
	Fast_Total_Operating_hrsCountSecs            int
	Slow_TempNacelleCountSecs                    int
	Fast_Total_Grid_OK_hrsCountSecs              int
	Fast_Total_WTG_OK_hrsCountSecs               int
	Slow_TempCabinetTopBoxCountSecs              int
	Slow_TempGeneratorBearingNDECountSecs        int
	Fast_Total_Access_hrsCountSecs               int
	Slow_TempBottomPowerSectionCountSecs         int
	Slow_TempGeneratorBearingDECountSecs         int
	Slow_TotalReactPowerIn_kVArhCountSecs        int
	Slow_TempBottomControlSectionCountSecs       int
	Slow_TempConv1CountSecs                      int
	Fast_ActivePowerRated_kWCountSecs            int
	Fast_NodeIPCountSecs                         int
	Fast_PitchSpeed1CountSecs                    int
	Slow_CFCardSizeCountSecs                     int
	Slow_CPU_NumberCountSecs                     int
	Slow_CFCardSpaceLeftCountSecs                int
	Slow_TempBottomCapSectionCountSecs           int
	Slow_RatedPowerCountSecs                     int
	Slow_TempConv3CountSecs                      int
	Slow_TempConv2CountSecs                      int
	Slow_TotalActPowerIn_kWhCountSecs            int
	Slow_TotalActPowerInG1_kWhCountSecs          int
	Slow_TotalActPowerInG2_kWhCountSecs          int
	Slow_TotalActPowerOutG2_kWhCountSecs         int
	Slow_TotalG2ActiveHoursCountSecs             int
	Slow_TotalReactPowerInG2_kVArhCountSecs      int
	Slow_TotalReactPowerOut_kVArhCountSecs       int
	Slow_UTCoffset_intCountSecs                  int

	// Line int
	File string
}

func (m *ScadaThreeSecsExt) New() *ScadaThreeSecsExt {
	timeStampStr := m.TimeStampSecondGroup.Format("060102_150405")
	m.ID = timeStampStr + "#" + m.ProjectName + "#" + m.Turbine
	return m
}

func (m *ScadaThreeSecsExt) RecordID() interface{} {
	return m.ID
}

func (m *ScadaThreeSecsExt) TableName() string {
	return "ScadaThreeSecsExt"
}
