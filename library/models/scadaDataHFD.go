package models

import (
	. "eaciit/wfdemo-git/library/helper"
	"time"

	"github.com/eaciit/orm"
)

type ScadaDataHFD struct {
	orm.ModelBase `bson:"-",json:"-"`
	ID            string ` bson:"_id" json:"_id" `
	TimeStamp     time.Time
	TimeStampInt  int64
	DateInfo      DateInfo
	ProjectName   string
	Turbine       string
	IsNull        bool
	TurbineState  float64

	Fast_ActivePower_kW        float64
	Fast_ActivePower_kW_StdDev float64
	Fast_ActivePower_kW_Min    float64
	Fast_ActivePower_kW_Max    float64
	Fast_ActivePower_kW_Count  int

	Fast_WindSpeed_ms        float64
	Fast_WindSpeed_ms_StdDev float64
	Fast_WindSpeed_ms_Min    float64
	Fast_WindSpeed_ms_Max    float64
	Fast_WindSpeed_ms_Count  int
	Fast_WindSpeed_Bin       float64

	Slow_NacellePos        float64
	Slow_NacellePos_StdDev float64
	Slow_NacellePos_Min    float64
	Slow_NacellePos_Max    float64
	Slow_NacellePos_Count  int

	Slow_WindDirection        float64
	Slow_WindDirection_StdDev float64
	Slow_WindDirection_Min    float64
	Slow_WindDirection_Max    float64
	Slow_WindDirection_Count  int

	Fast_CurrentL3        float64
	Fast_CurrentL3_StdDev float64
	Fast_CurrentL3_Min    float64
	Fast_CurrentL3_Max    float64
	Fast_CurrentL3_Count  int

	Fast_CurrentL1        float64
	Fast_CurrentL1_StdDev float64
	Fast_CurrentL1_Min    float64
	Fast_CurrentL1_Max    float64
	Fast_CurrentL1_Count  int

	Fast_ActivePowerSetpoint_kW        float64
	Fast_ActivePowerSetpoint_kW_StdDev float64
	Fast_ActivePowerSetpoint_kW_Min    float64
	Fast_ActivePowerSetpoint_kW_Max    float64
	Fast_ActivePowerSetpoint_kW_Count  int

	Fast_CurrentL2        float64
	Fast_CurrentL2_StdDev float64
	Fast_CurrentL2_Min    float64
	Fast_CurrentL2_Max    float64
	Fast_CurrentL2_Count  int

	Fast_DrTrVibValue        float64
	Fast_DrTrVibValue_StdDev float64
	Fast_DrTrVibValue_Min    float64
	Fast_DrTrVibValue_Max    float64
	Fast_DrTrVibValue_Count  int

	Fast_GenSpeed_RPM        float64
	Fast_GenSpeed_RPM_StdDev float64
	Fast_GenSpeed_RPM_Min    float64
	Fast_GenSpeed_RPM_Max    float64
	Fast_GenSpeed_RPM_Count  int

	Fast_PitchAccuV1        float64
	Fast_PitchAccuV1_StdDev float64
	Fast_PitchAccuV1_Min    float64
	Fast_PitchAccuV1_Max    float64
	Fast_PitchAccuV1_Count  int

	Fast_PitchAngle        float64
	Fast_PitchAngle_StdDev float64
	Fast_PitchAngle_Min    float64
	Fast_PitchAngle_Max    float64
	Fast_PitchAngle_Count  int

	Fast_PitchAngle3        float64
	Fast_PitchAngle3_StdDev float64
	Fast_PitchAngle3_Min    float64
	Fast_PitchAngle3_Max    float64
	Fast_PitchAngle3_Count  int

	Fast_PitchAngle2        float64
	Fast_PitchAngle2_StdDev float64
	Fast_PitchAngle2_Min    float64
	Fast_PitchAngle2_Max    float64
	Fast_PitchAngle2_Count  int

	Fast_PitchConvCurrent1        float64
	Fast_PitchConvCurrent1_StdDev float64
	Fast_PitchConvCurrent1_Min    float64
	Fast_PitchConvCurrent1_Max    float64
	Fast_PitchConvCurrent1_Count  int

	Fast_PitchConvCurrent3        float64
	Fast_PitchConvCurrent3_StdDev float64
	Fast_PitchConvCurrent3_Min    float64
	Fast_PitchConvCurrent3_Max    float64
	Fast_PitchConvCurrent3_Count  int

	Fast_PitchConvCurrent2        float64
	Fast_PitchConvCurrent2_StdDev float64
	Fast_PitchConvCurrent2_Min    float64
	Fast_PitchConvCurrent2_Max    float64
	Fast_PitchConvCurrent2_Count  int

	Fast_PowerFactor        float64
	Fast_PowerFactor_StdDev float64
	Fast_PowerFactor_Min    float64
	Fast_PowerFactor_Max    float64
	Fast_PowerFactor_Count  int

	Fast_ReactivePowerSetpointPPC_kVA        float64
	Fast_ReactivePowerSetpointPPC_kVA_StdDev float64
	Fast_ReactivePowerSetpointPPC_kVA_Min    float64
	Fast_ReactivePowerSetpointPPC_kVA_Max    float64
	Fast_ReactivePowerSetpointPPC_kVA_Count  int

	Fast_ReactivePower_kVAr        float64
	Fast_ReactivePower_kVAr_StdDev float64
	Fast_ReactivePower_kVAr_Min    float64
	Fast_ReactivePower_kVAr_Max    float64
	Fast_ReactivePower_kVAr_Count  int

	Fast_RotorSpeed_RPM        float64
	Fast_RotorSpeed_RPM_StdDev float64
	Fast_RotorSpeed_RPM_Min    float64
	Fast_RotorSpeed_RPM_Max    float64
	Fast_RotorSpeed_RPM_Count  int

	Fast_VoltageL1        float64
	Fast_VoltageL1_StdDev float64
	Fast_VoltageL1_Min    float64
	Fast_VoltageL1_Max    float64
	Fast_VoltageL1_Count  int

	Fast_VoltageL2        float64
	Fast_VoltageL2_StdDev float64
	Fast_VoltageL2_Min    float64
	Fast_VoltageL2_Max    float64
	Fast_VoltageL2_Count  int

	Slow_CapableCapacitiveReactPwr_kVAr        float64
	Slow_CapableCapacitiveReactPwr_kVAr_StdDev float64
	Slow_CapableCapacitiveReactPwr_kVAr_Min    float64
	Slow_CapableCapacitiveReactPwr_kVAr_Max    float64
	Slow_CapableCapacitiveReactPwr_kVAr_Count  int

	Slow_CapableInductiveReactPwr_kVAr        float64
	Slow_CapableInductiveReactPwr_kVAr_StdDev float64
	Slow_CapableInductiveReactPwr_kVAr_Min    float64
	Slow_CapableInductiveReactPwr_kVAr_Max    float64
	Slow_CapableInductiveReactPwr_kVAr_Count  int

	Slow_DateTime_Sec        float64
	Slow_DateTime_Sec_StdDev float64
	Slow_DateTime_Sec_Min    float64
	Slow_DateTime_Sec_Max    float64
	Slow_DateTime_Sec_Count  int

	Fast_PitchAngle1        float64
	Fast_PitchAngle1_StdDev float64
	Fast_PitchAngle1_Min    float64
	Fast_PitchAngle1_Max    float64
	Fast_PitchAngle1_Count  int

	Fast_VoltageL3        float64
	Fast_VoltageL3_StdDev float64
	Fast_VoltageL3_Min    float64
	Fast_VoltageL3_Max    float64
	Fast_VoltageL3_Count  int

	Slow_CapableCapacitivePwrFactor        float64
	Slow_CapableCapacitivePwrFactor_StdDev float64
	Slow_CapableCapacitivePwrFactor_Min    float64
	Slow_CapableCapacitivePwrFactor_Max    float64
	Slow_CapableCapacitivePwrFactor_Count  int

	Fast_Total_Production_kWh        float64
	Fast_Total_Production_kWh_StdDev float64
	Fast_Total_Production_kWh_Min    float64
	Fast_Total_Production_kWh_Max    float64
	Fast_Total_Production_kWh_Count  int

	Fast_Total_Prod_Day_kWh        float64
	Fast_Total_Prod_Day_kWh_StdDev float64
	Fast_Total_Prod_Day_kWh_Min    float64
	Fast_Total_Prod_Day_kWh_Max    float64
	Fast_Total_Prod_Day_kWh_Count  int

	Fast_Total_Prod_Month_kWh        float64
	Fast_Total_Prod_Month_kWh_StdDev float64
	Fast_Total_Prod_Month_kWh_Min    float64
	Fast_Total_Prod_Month_kWh_Max    float64
	Fast_Total_Prod_Month_kWh_Count  int

	Fast_ActivePowerOutPWCSell_kW        float64
	Fast_ActivePowerOutPWCSell_kW_StdDev float64
	Fast_ActivePowerOutPWCSell_kW_Min    float64
	Fast_ActivePowerOutPWCSell_kW_Max    float64
	Fast_ActivePowerOutPWCSell_kW_Count  int

	Fast_Frequency_Hz        float64
	Fast_Frequency_Hz_StdDev float64
	Fast_Frequency_Hz_Min    float64
	Fast_Frequency_Hz_Max    float64
	Fast_Frequency_Hz_Count  int

	Slow_TempG1L2        float64
	Slow_TempG1L2_StdDev float64
	Slow_TempG1L2_Min    float64
	Slow_TempG1L2_Max    float64
	Slow_TempG1L2_Count  int

	Slow_TempG1L3        float64
	Slow_TempG1L3_StdDev float64
	Slow_TempG1L3_Min    float64
	Slow_TempG1L3_Max    float64
	Slow_TempG1L3_Count  int

	Slow_TempGearBoxHSSDE        float64
	Slow_TempGearBoxHSSDE_StdDev float64
	Slow_TempGearBoxHSSDE_Min    float64
	Slow_TempGearBoxHSSDE_Max    float64
	Slow_TempGearBoxHSSDE_Count  int

	Slow_TempGearBoxIMSNDE        float64
	Slow_TempGearBoxIMSNDE_StdDev float64
	Slow_TempGearBoxIMSNDE_Min    float64
	Slow_TempGearBoxIMSNDE_Max    float64
	Slow_TempGearBoxIMSNDE_Count  int

	Slow_TempOutdoor        float64
	Slow_TempOutdoor_StdDev float64
	Slow_TempOutdoor_Min    float64
	Slow_TempOutdoor_Max    float64
	Slow_TempOutdoor_Count  int

	Fast_PitchAccuV3        float64
	Fast_PitchAccuV3_StdDev float64
	Fast_PitchAccuV3_Min    float64
	Fast_PitchAccuV3_Max    float64
	Fast_PitchAccuV3_Count  int

	Slow_TotalTurbineActiveHours        float64
	Slow_TotalTurbineActiveHours_StdDev float64
	Slow_TotalTurbineActiveHours_Min    float64
	Slow_TotalTurbineActiveHours_Max    float64
	Slow_TotalTurbineActiveHours_Count  int

	Slow_TotalTurbineOKHours        float64
	Slow_TotalTurbineOKHours_StdDev float64
	Slow_TotalTurbineOKHours_Min    float64
	Slow_TotalTurbineOKHours_Max    float64
	Slow_TotalTurbineOKHours_Count  int

	Slow_TotalTurbineTimeAllHours        float64
	Slow_TotalTurbineTimeAllHours_StdDev float64
	Slow_TotalTurbineTimeAllHours_Min    float64
	Slow_TotalTurbineTimeAllHours_Max    float64
	Slow_TotalTurbineTimeAllHours_Count  int

	Slow_TempG1L1        float64
	Slow_TempG1L1_StdDev float64
	Slow_TempG1L1_Min    float64
	Slow_TempG1L1_Max    float64
	Slow_TempG1L1_Count  int

	Slow_TempGearBoxOilSump        float64
	Slow_TempGearBoxOilSump_StdDev float64
	Slow_TempGearBoxOilSump_Min    float64
	Slow_TempGearBoxOilSump_Max    float64
	Slow_TempGearBoxOilSump_Count  int

	Fast_PitchAccuV2        float64
	Fast_PitchAccuV2_StdDev float64
	Fast_PitchAccuV2_Min    float64
	Fast_PitchAccuV2_Max    float64
	Fast_PitchAccuV2_Count  int

	Slow_TotalGridOkHours        float64
	Slow_TotalGridOkHours_StdDev float64
	Slow_TotalGridOkHours_Min    float64
	Slow_TotalGridOkHours_Max    float64
	Slow_TotalGridOkHours_Count  int

	Slow_TotalActPowerOut_kWh        float64
	Slow_TotalActPowerOut_kWh_StdDev float64
	Slow_TotalActPowerOut_kWh_Min    float64
	Slow_TotalActPowerOut_kWh_Max    float64
	Slow_TotalActPowerOut_kWh_Count  int

	Fast_YawService        float64
	Fast_YawService_StdDev float64
	Fast_YawService_Min    float64
	Fast_YawService_Max    float64
	Fast_YawService_Count  int

	Fast_YawAngle        float64
	Fast_YawAngle_StdDev float64
	Fast_YawAngle_Min    float64
	Fast_YawAngle_Max    float64
	Fast_YawAngle_Count  int

	Slow_CapableInductivePwrFactor        float64
	Slow_CapableInductivePwrFactor_StdDev float64
	Slow_CapableInductivePwrFactor_Min    float64
	Slow_CapableInductivePwrFactor_Max    float64
	Slow_CapableInductivePwrFactor_Count  int

	Slow_TempGearBoxHSSNDE        float64
	Slow_TempGearBoxHSSNDE_StdDev float64
	Slow_TempGearBoxHSSNDE_Min    float64
	Slow_TempGearBoxHSSNDE_Max    float64
	Slow_TempGearBoxHSSNDE_Count  int

	Slow_TempHubBearing        float64
	Slow_TempHubBearing_StdDev float64
	Slow_TempHubBearing_Min    float64
	Slow_TempHubBearing_Max    float64
	Slow_TempHubBearing_Count  int

	Slow_TotalG1ActiveHours        float64
	Slow_TotalG1ActiveHours_StdDev float64
	Slow_TotalG1ActiveHours_Min    float64
	Slow_TotalG1ActiveHours_Max    float64
	Slow_TotalG1ActiveHours_Count  int

	Slow_TotalActPowerOutG1_kWh        float64
	Slow_TotalActPowerOutG1_kWh_StdDev float64
	Slow_TotalActPowerOutG1_kWh_Min    float64
	Slow_TotalActPowerOutG1_kWh_Max    float64
	Slow_TotalActPowerOutG1_kWh_Count  int

	Slow_TotalReactPowerInG1_kVArh        float64
	Slow_TotalReactPowerInG1_kVArh_StdDev float64
	Slow_TotalReactPowerInG1_kVArh_Min    float64
	Slow_TotalReactPowerInG1_kVArh_Max    float64
	Slow_TotalReactPowerInG1_kVArh_Count  int

	Slow_NacelleDrill        float64
	Slow_NacelleDrill_StdDev float64
	Slow_NacelleDrill_Min    float64
	Slow_NacelleDrill_Max    float64
	Slow_NacelleDrill_Count  int

	Slow_TempGearBoxIMSDE        float64
	Slow_TempGearBoxIMSDE_StdDev float64
	Slow_TempGearBoxIMSDE_Min    float64
	Slow_TempGearBoxIMSDE_Max    float64
	Slow_TempGearBoxIMSDE_Count  int

	Fast_Total_Operating_hrs        float64
	Fast_Total_Operating_hrs_StdDev float64
	Fast_Total_Operating_hrs_Min    float64
	Fast_Total_Operating_hrs_Max    float64
	Fast_Total_Operating_hrs_Count  int

	Slow_TempNacelle        float64
	Slow_TempNacelle_StdDev float64
	Slow_TempNacelle_Min    float64
	Slow_TempNacelle_Max    float64
	Slow_TempNacelle_Count  int

	Fast_Total_Grid_OK_hrs        float64
	Fast_Total_Grid_OK_hrs_StdDev float64
	Fast_Total_Grid_OK_hrs_Min    float64
	Fast_Total_Grid_OK_hrs_Max    float64
	Fast_Total_Grid_OK_hrs_Count  int

	Fast_Total_WTG_OK_hrs        float64
	Fast_Total_WTG_OK_hrs_StdDev float64
	Fast_Total_WTG_OK_hrs_Min    float64
	Fast_Total_WTG_OK_hrs_Max    float64
	Fast_Total_WTG_OK_hrs_Count  int

	Slow_TempCabinetTopBox        float64
	Slow_TempCabinetTopBox_StdDev float64
	Slow_TempCabinetTopBox_Min    float64
	Slow_TempCabinetTopBox_Max    float64
	Slow_TempCabinetTopBox_Count  int

	Slow_TempGeneratorBearingNDE        float64
	Slow_TempGeneratorBearingNDE_StdDev float64
	Slow_TempGeneratorBearingNDE_Min    float64
	Slow_TempGeneratorBearingNDE_Max    float64
	Slow_TempGeneratorBearingNDE_Count  int

	Fast_Total_Access_hrs        float64
	Fast_Total_Access_hrs_StdDev float64
	Fast_Total_Access_hrs_Min    float64
	Fast_Total_Access_hrs_Max    float64
	Fast_Total_Access_hrs_Count  int

	Slow_TempBottomPowerSection        float64
	Slow_TempBottomPowerSection_StdDev float64
	Slow_TempBottomPowerSection_Min    float64
	Slow_TempBottomPowerSection_Max    float64
	Slow_TempBottomPowerSection_Count  int

	Slow_TempGeneratorBearingDE        float64
	Slow_TempGeneratorBearingDE_StdDev float64
	Slow_TempGeneratorBearingDE_Min    float64
	Slow_TempGeneratorBearingDE_Max    float64
	Slow_TempGeneratorBearingDE_Count  int

	Slow_TotalReactPowerIn_kVArh        float64
	Slow_TotalReactPowerIn_kVArh_StdDev float64
	Slow_TotalReactPowerIn_kVArh_Min    float64
	Slow_TotalReactPowerIn_kVArh_Max    float64
	Slow_TotalReactPowerIn_kVArh_Count  int

	Slow_TempBottomControlSection        float64
	Slow_TempBottomControlSection_StdDev float64
	Slow_TempBottomControlSection_Min    float64
	Slow_TempBottomControlSection_Max    float64
	Slow_TempBottomControlSection_Count  int

	Slow_TempConv1        float64
	Slow_TempConv1_StdDev float64
	Slow_TempConv1_Min    float64
	Slow_TempConv1_Max    float64
	Slow_TempConv1_Count  int

	Fast_ActivePowerRated_kW        float64
	Fast_ActivePowerRated_kW_StdDev float64
	Fast_ActivePowerRated_kW_Min    float64
	Fast_ActivePowerRated_kW_Max    float64
	Fast_ActivePowerRated_kW_Count  int

	Fast_NodeIP        float64
	Fast_NodeIP_StdDev float64
	Fast_NodeIP_Min    float64
	Fast_NodeIP_Max    float64
	Fast_NodeIP_Count  int

	Fast_PitchSpeed1        float64
	Fast_PitchSpeed1_StdDev float64
	Fast_PitchSpeed1_Min    float64
	Fast_PitchSpeed1_Max    float64
	Fast_PitchSpeed1_Count  int

	Slow_CFCardSize        float64
	Slow_CFCardSize_StdDev float64
	Slow_CFCardSize_Min    float64
	Slow_CFCardSize_Max    float64
	Slow_CFCardSize_Count  int

	Slow_CPU_Number        float64
	Slow_CPU_Number_StdDev float64
	Slow_CPU_Number_Min    float64
	Slow_CPU_Number_Max    float64
	Slow_CPU_Number_Count  int

	Slow_CFCardSpaceLeft        float64
	Slow_CFCardSpaceLeft_StdDev float64
	Slow_CFCardSpaceLeft_Min    float64
	Slow_CFCardSpaceLeft_Max    float64
	Slow_CFCardSpaceLeft_Count  int

	Slow_TempBottomCapSection        float64
	Slow_TempBottomCapSection_StdDev float64
	Slow_TempBottomCapSection_Min    float64
	Slow_TempBottomCapSection_Max    float64
	Slow_TempBottomCapSection_Count  int

	Slow_RatedPower        float64
	Slow_RatedPower_StdDev float64
	Slow_RatedPower_Min    float64
	Slow_RatedPower_Max    float64
	Slow_RatedPower_Count  int

	Slow_TempConv3        float64
	Slow_TempConv3_StdDev float64
	Slow_TempConv3_Min    float64
	Slow_TempConv3_Max    float64
	Slow_TempConv3_Count  int

	Slow_TempConv2        float64
	Slow_TempConv2_StdDev float64
	Slow_TempConv2_Min    float64
	Slow_TempConv2_Max    float64
	Slow_TempConv2_Count  int

	Slow_TotalActPowerIn_kWh        float64
	Slow_TotalActPowerIn_kWh_StdDev float64
	Slow_TotalActPowerIn_kWh_Min    float64
	Slow_TotalActPowerIn_kWh_Max    float64
	Slow_TotalActPowerIn_kWh_Count  int

	Slow_TotalActPowerInG1_kWh        float64
	Slow_TotalActPowerInG1_kWh_StdDev float64
	Slow_TotalActPowerInG1_kWh_Min    float64
	Slow_TotalActPowerInG1_kWh_Max    float64
	Slow_TotalActPowerInG1_kWh_Count  int

	Slow_TotalActPowerInG2_kWh        float64
	Slow_TotalActPowerInG2_kWh_StdDev float64
	Slow_TotalActPowerInG2_kWh_Min    float64
	Slow_TotalActPowerInG2_kWh_Max    float64
	Slow_TotalActPowerInG2_kWh_Count  int

	Slow_TotalActPowerOutG2_kWh        float64
	Slow_TotalActPowerOutG2_kWh_StdDev float64
	Slow_TotalActPowerOutG2_kWh_Min    float64
	Slow_TotalActPowerOutG2_kWh_Max    float64
	Slow_TotalActPowerOutG2_kWh_Count  int

	Slow_TotalG2ActiveHours        float64
	Slow_TotalG2ActiveHours_StdDev float64
	Slow_TotalG2ActiveHours_Min    float64
	Slow_TotalG2ActiveHours_Max    float64
	Slow_TotalG2ActiveHours_Count  int

	Slow_TotalReactPowerInG2_kVArh        float64
	Slow_TotalReactPowerInG2_kVArh_StdDev float64
	Slow_TotalReactPowerInG2_kVArh_Min    float64
	Slow_TotalReactPowerInG2_kVArh_Max    float64
	Slow_TotalReactPowerInG2_kVArh_Count  int

	Slow_TotalReactPowerOut_kVArh        float64
	Slow_TotalReactPowerOut_kVArh_StdDev float64
	Slow_TotalReactPowerOut_kVArh_Min    float64
	Slow_TotalReactPowerOut_kVArh_Max    float64
	Slow_TotalReactPowerOut_kVArh_Count  int

	Slow_UTCoffset_int        float64
	Slow_UTCoffset_int_StdDev float64
	Slow_UTCoffset_int_Min    float64
	Slow_UTCoffset_int_Max    float64
	Slow_UTCoffset_int_Count  int

	File  string
	No    int
	Count int
}

func (m *ScadaDataHFD) New() *ScadaDataHFD {
	timeStampStr := m.TimeStamp.UTC().Format("060102_1504")
	m.ID = timeStampStr + "#" + m.ProjectName + "#" + m.Turbine // + "#" + tk.ToString(m.No)
	return m
}

func (m *ScadaDataHFD) RecordID() interface{} {
	return m.ID
}

func (m *ScadaDataHFD) TableName() string {
	return "ScadaDataHFD"
}
