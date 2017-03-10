package models

import (
	. "eaciit/wfdemo-git/library/helper"
	"fmt"
	"github.com/eaciit/orm"
	"gopkg.in/mgo.v2/bson"
	"time"
)

type ScadaHFD struct {
	orm.ModelBase   `bson:"-",json:"-"`
	Id              bson.ObjectId ` bson:"_id" , json:"_id" `
	TimeStamp       time.Time
	DateInfo        DateInfo
	ProjectName     string
	Turbine         string
	ActivePower     float64
	Production      float64
	OprHours        float64
	WtgOkHours      float64
	WindSpeed       float64
	WindDirection   float64
	NacellePosition float64
	Temperature     float64
	PitchAngle      float64
	RotorRPM        float64
}

func (m *ScadaHFD) New() *ScadaHFD {
	m.Id = bson.NewObjectId()
	return m
}

func (m *ScadaHFD) RecordID() interface{} {
	return m.Id
}

func (m *ScadaHFD) TableName() string {
	return "ScadaHFD"
}

type ScadaMonitoring struct {
	orm.ModelBase        `bson:"-",json:"-"`
	Id                   string ` bson:"_id" json:"_id" `
	TimeStamp            time.Time
	DateInfo             DateInfo
	ProjectName          string
	ActivePower          float64
	Production           float64
	OprHours             float64
	WtgOkHours           float64
	WindSpeed            float64
	WindSpeedCount       int
	WindDirection        float64
	WindDirectionCount   int
	NacellePosition      float64
	NacellePositionCount int
	Temperature          float64
	TemperatureCount     int
	PitchAngle           float64
	PitchAngleCount      int
	RotorRPM             float64
	RotorRPMCount        int
	Detail               []ScadaMonitoringItem
}

type ScadaMonitoringItem struct {
	Turbine         string
	ActivePower     float64
	WindSpeed       float64
	WindDirection   float64
	NacellePosition float64
	Temperature     float64
	PitchAngle      float64
	RotorRPM        float64
	TimeUpdate      time.Time
	DataComing      int
}

func (m *ScadaMonitoring) New() *ScadaMonitoring {
	m.Id = "Bhesada_Update"
	return m
}

func (m *ScadaMonitoring) RecordID() interface{} {
	return m.Id
}

func (m *ScadaMonitoring) TableName() string {
	return "ScadaMonitoring"
}

type ScadaRealTime struct {
	orm.ModelBase   `bson:"-",json:"-"`
	Id              bson.ObjectId ` bson:"_id" json:"_id" `
	TimeStamp       time.Time
	DateInfo        DateInfo
	ProjectName     string
	Turbine         string
	ActivePower     float64
	Production      float64
	OprHours        float64
	WtgOkHours      float64
	WindSpeed       float64
	WindDirection   float64
	NacellePosition float64
	Temperature     float64
	PitchAngle      float64
	RotorRPM        float64
	LastUpdate      time.Time

	RotorSpeed_RPM float64
	GenSpeed_RPM   float64

	PitchAngle1 float64
	PitchAngle2 float64
	PitchAngle3 float64

	PitchAccuV1 float64
	PitchAccuV2 float64
	PitchAccuV3 float64

	PitchConvCurrent1 float64
	PitchConvCurrent2 float64
	PitchConvCurrent3 float64

	TempConv1 float64
	TempConv2 float64
	TempConv3 float64

	VoltageL1 float64
	VoltageL2 float64
	VoltageL3 float64

	CurrentL1 float64
	CurrentL2 float64
	CurrentL3 float64

	ReactivePower_kVAr float64
	Frequency_Hz       float64

	Total_Prod_Day_kWh float64
	PowerFactor        float64

	TempG1L1                float64
	TempG1L2                float64
	TempG1L3                float64
	TempGeneratorBearingDE  float64
	TempGeneratorBearingNDE float64

	TempGearBoxHSSDE   float64
	TempGearBoxHSSNDE  float64
	TempGearBoxIMSDE   float64
	TempGearBoxIMSNDE  float64
	TempGearBoxOilSump float64

	TempNacelle    float64
	TempOutdoor    float64
	TempHubBearing float64

	DrTrVibValue float64
}

func (m *ScadaRealTime) New() *ScadaRealTime {
	m.Id = bson.NewObjectId()
	return m
}

func (m *ScadaRealTime) RecordID() interface{} {
	return m.Id
}

func (m *ScadaRealTime) TableName() string {
	return "ScadaRealTime"
}

type Scada10Min struct {
	orm.ModelBase `bson:"-",json:"-"`
	Id            string
	TimeStamp     time.Time
	DateInfo      DateInfo
	ProjectName   string
	Turbine       string
	IsNull        bool

	ErrorState float64
	NodeIP     float64

	ActivePower_kW        float64
	ActivePower_kW_sum    float64
	ActivePower_kW_min    float64
	ActivePower_kW_max    float64
	ActivePower_kW_stddev float64
	ActivePower_kW_count  float64

	ActivePowerCurtailmentSource        float64
	ActivePowerCurtailmentSource_sum    float64
	ActivePowerCurtailmentSource_min    float64
	ActivePowerCurtailmentSource_max    float64
	ActivePowerCurtailmentSource_stddev float64
	ActivePowerCurtailmentSource_count  float64

	ActivePowerOutPWC_kW        float64
	ActivePowerOutPWC_kW_sum    float64
	ActivePowerOutPWC_kW_min    float64
	ActivePowerOutPWC_kW_max    float64
	ActivePowerOutPWC_kW_stddev float64
	ActivePowerOutPWC_kW_count  float64

	ActivePowerOutPWCSell_kW        float64
	ActivePowerOutPWCSell_kW_sum    float64
	ActivePowerOutPWCSell_kW_min    float64
	ActivePowerOutPWCSell_kW_max    float64
	ActivePowerOutPWCSell_kW_stddev float64
	ActivePowerOutPWCSell_kW_count  float64

	ActivePowerRated_kW        float64
	ActivePowerRated_kW_sum    float64
	ActivePowerRated_kW_min    float64
	ActivePowerRated_kW_max    float64
	ActivePowerRated_kW_stddev float64
	ActivePowerRated_kW_count  float64

	ActivePowerSetpoint_kW        float64
	ActivePowerSetpoint_kW_sum    float64
	ActivePowerSetpoint_kW_min    float64
	ActivePowerSetpoint_kW_max    float64
	ActivePowerSetpoint_kW_stddev float64
	ActivePowerSetpoint_kW_count  float64

	ActivePowerSetpointPPC_kW        float64
	ActivePowerSetpointPPC_kW_sum    float64
	ActivePowerSetpointPPC_kW_min    float64
	ActivePowerSetpointPPC_kW_max    float64
	ActivePowerSetpointPPC_kW_stddev float64
	ActivePowerSetpointPPC_kW_count  float64

	AlarmCode        float64
	AlarmCode_sum    float64
	AlarmCode_min    float64
	AlarmCode_max    float64
	AlarmCode_stddev float64
	AlarmCode_count  float64

	AlarmCode_DetectTime        float64
	AlarmCode_DetectTime_sum    float64
	AlarmCode_DetectTime_min    float64
	AlarmCode_DetectTime_max    float64
	AlarmCode_DetectTime_stddev float64
	AlarmCode_DetectTime_count  float64

	CapableCapacitivePwrFactor        float64
	CapableCapacitivePwrFactor_sum    float64
	CapableCapacitivePwrFactor_min    float64
	CapableCapacitivePwrFactor_max    float64
	CapableCapacitivePwrFactor_stddev float64
	CapableCapacitivePwrFactor_count  float64

	CapableCapacitiveReactPwr_kVAr        float64
	CapableCapacitiveReactPwr_kVAr_sum    float64
	CapableCapacitiveReactPwr_kVAr_min    float64
	CapableCapacitiveReactPwr_kVAr_max    float64
	CapableCapacitiveReactPwr_kVAr_stddev float64
	CapableCapacitiveReactPwr_kVAr_count  float64

	CapableInductivePwrFactor        float64
	CapableInductivePwrFactor_sum    float64
	CapableInductivePwrFactor_min    float64
	CapableInductivePwrFactor_max    float64
	CapableInductivePwrFactor_stddev float64
	CapableInductivePwrFactor_count  float64

	CapableInductiveReactPwr_kVAr        float64
	CapableInductiveReactPwr_kVAr_sum    float64
	CapableInductiveReactPwr_kVAr_min    float64
	CapableInductiveReactPwr_kVAr_max    float64
	CapableInductiveReactPwr_kVAr_stddev float64
	CapableInductiveReactPwr_kVAr_count  float64

	CFCardSize        float64
	CFCardSize_sum    float64
	CFCardSize_min    float64
	CFCardSize_max    float64
	CFCardSize_stddev float64
	CFCardSize_count  float64

	CFCardSpaceLeft        float64
	CFCardSpaceLeft_sum    float64
	CFCardSpaceLeft_min    float64
	CFCardSpaceLeft_max    float64
	CFCardSpaceLeft_stddev float64
	CFCardSpaceLeft_count  float64

	CPU_Number        float64
	CPU_Number_sum    float64
	CPU_Number_min    float64
	CPU_Number_max    float64
	CPU_Number_stddev float64
	CPU_Number_count  float64

	CurrentL1        float64
	CurrentL1_sum    float64
	CurrentL1_min    float64
	CurrentL1_max    float64
	CurrentL1_stddev float64
	CurrentL1_count  float64

	CurrentL2        float64
	CurrentL2_sum    float64
	CurrentL2_min    float64
	CurrentL2_max    float64
	CurrentL2_stddev float64
	CurrentL2_count  float64

	CurrentL3        float64
	CurrentL3_sum    float64
	CurrentL3_min    float64
	CurrentL3_max    float64
	CurrentL3_stddev float64
	CurrentL3_count  float64

	DateTime        float64
	DateTime_sum    float64
	DateTime_min    float64
	DateTime_max    float64
	DateTime_stddev float64
	DateTime_count  float64

	DateTime_Sec        float64
	DateTime_Sec_sum    float64
	DateTime_Sec_min    float64
	DateTime_Sec_max    float64
	DateTime_Sec_stddev float64
	DateTime_Sec_count  float64

	DrTrVibValue        float64
	DrTrVibValue_sum    float64
	DrTrVibValue_min    float64
	DrTrVibValue_max    float64
	DrTrVibValue_stddev float64
	DrTrVibValue_count  float64

	Frequency_Hz        float64
	Frequency_Hz_sum    float64
	Frequency_Hz_min    float64
	Frequency_Hz_max    float64
	Frequency_Hz_stddev float64
	Frequency_Hz_count  float64

	GenSpeed_RPM        float64
	GenSpeed_RPM_sum    float64
	GenSpeed_RPM_min    float64
	GenSpeed_RPM_max    float64
	GenSpeed_RPM_stddev float64
	GenSpeed_RPM_count  float64

	NacelleDrill        float64
	NacelleDrill_sum    float64
	NacelleDrill_min    float64
	NacelleDrill_max    float64
	NacelleDrill_stddev float64
	NacelleDrill_count  float64

	NacellePos        float64
	NacellePos_sum    float64
	NacellePos_min    float64
	NacellePos_max    float64
	NacellePos_stddev float64
	NacellePos_count  float64

	PitchAccuV1        float64
	PitchAccuV1_sum    float64
	PitchAccuV1_min    float64
	PitchAccuV1_max    float64
	PitchAccuV1_stddev float64
	PitchAccuV1_count  float64

	PitchAccuV2        float64
	PitchAccuV2_sum    float64
	PitchAccuV2_min    float64
	PitchAccuV2_max    float64
	PitchAccuV2_stddev float64
	PitchAccuV2_count  float64

	PitchAccuV3        float64
	PitchAccuV3_sum    float64
	PitchAccuV3_min    float64
	PitchAccuV3_max    float64
	PitchAccuV3_stddev float64
	PitchAccuV3_count  float64

	PitchAngle        float64
	PitchAngle_sum    float64
	PitchAngle_min    float64
	PitchAngle_max    float64
	PitchAngle_stddev float64
	PitchAngle_count  float64

	PitchAngle1        float64
	PitchAngle1_sum    float64
	PitchAngle1_min    float64
	PitchAngle1_max    float64
	PitchAngle1_stddev float64
	PitchAngle1_count  float64

	PitchAngle2        float64
	PitchAngle2_sum    float64
	PitchAngle2_min    float64
	PitchAngle2_max    float64
	PitchAngle2_stddev float64
	PitchAngle2_count  float64

	PitchAngle3        float64
	PitchAngle3_sum    float64
	PitchAngle3_min    float64
	PitchAngle3_max    float64
	PitchAngle3_stddev float64
	PitchAngle3_count  float64

	PitchConvCurrent1        float64
	PitchConvCurrent1_sum    float64
	PitchConvCurrent1_min    float64
	PitchConvCurrent1_max    float64
	PitchConvCurrent1_stddev float64
	PitchConvCurrent1_count  float64

	PitchConvCurrent2        float64
	PitchConvCurrent2_sum    float64
	PitchConvCurrent2_min    float64
	PitchConvCurrent2_max    float64
	PitchConvCurrent2_stddev float64
	PitchConvCurrent2_count  float64

	PitchConvCurrent3        float64
	PitchConvCurrent3_sum    float64
	PitchConvCurrent3_min    float64
	PitchConvCurrent3_max    float64
	PitchConvCurrent3_stddev float64
	PitchConvCurrent3_count  float64

	PitchSpeed1        float64
	PitchSpeed1_sum    float64
	PitchSpeed1_min    float64
	PitchSpeed1_max    float64
	PitchSpeed1_stddev float64
	PitchSpeed1_count  float64

	PowerFactor        float64
	PowerFactor_sum    float64
	PowerFactor_min    float64
	PowerFactor_max    float64
	PowerFactor_stddev float64
	PowerFactor_count  float64

	RatedPower        float64
	RatedPower_sum    float64
	RatedPower_min    float64
	RatedPower_max    float64
	RatedPower_stddev float64
	RatedPower_count  float64

	ReactivePower_kVAr        float64
	ReactivePower_kVAr_sum    float64
	ReactivePower_kVAr_min    float64
	ReactivePower_kVAr_max    float64
	ReactivePower_kVAr_stddev float64
	ReactivePower_kVAr_count  float64

	ReactivePowerSetpointPPC_kVAr        float64
	ReactivePowerSetpointPPC_kVAr_sum    float64
	ReactivePowerSetpointPPC_kVAr_min    float64
	ReactivePowerSetpointPPC_kVAr_max    float64
	ReactivePowerSetpointPPC_kVAr_stddev float64
	ReactivePowerSetpointPPC_kVAr_count  float64

	ReturnHeartbeat        float64
	ReturnHeartbeat_sum    float64
	ReturnHeartbeat_min    float64
	ReturnHeartbeat_max    float64
	ReturnHeartbeat_stddev float64
	ReturnHeartbeat_count  float64

	RotorSpeed_RPM        float64
	RotorSpeed_RPM_sum    float64
	RotorSpeed_RPM_min    float64
	RotorSpeed_RPM_max    float64
	RotorSpeed_RPM_stddev float64
	RotorSpeed_RPM_count  float64

	SoftwareRelease        float64
	SoftwareRelease_sum    float64
	SoftwareRelease_min    float64
	SoftwareRelease_max    float64
	SoftwareRelease_stddev float64
	SoftwareRelease_count  float64

	TempBottomCapSection        float64
	TempBottomCapSection_sum    float64
	TempBottomCapSection_min    float64
	TempBottomCapSection_max    float64
	TempBottomCapSection_stddev float64
	TempBottomCapSection_count  float64

	TempBottomControlSection        float64
	TempBottomControlSection_sum    float64
	TempBottomControlSection_min    float64
	TempBottomControlSection_max    float64
	TempBottomControlSection_stddev float64
	TempBottomControlSection_count  float64

	TempBottomPowerSection        float64
	TempBottomPowerSection_sum    float64
	TempBottomPowerSection_min    float64
	TempBottomPowerSection_max    float64
	TempBottomPowerSection_stddev float64
	TempBottomPowerSection_count  float64

	TempCabinetTopBox        float64
	TempCabinetTopBox_sum    float64
	TempCabinetTopBox_min    float64
	TempCabinetTopBox_max    float64
	TempCabinetTopBox_stddev float64
	TempCabinetTopBox_count  float64

	TempConv1        float64
	TempConv1_sum    float64
	TempConv1_min    float64
	TempConv1_max    float64
	TempConv1_stddev float64
	TempConv1_count  float64

	TempConv2        float64
	TempConv2_sum    float64
	TempConv2_min    float64
	TempConv2_max    float64
	TempConv2_stddev float64
	TempConv2_count  float64

	TempConv3        float64
	TempConv3_sum    float64
	TempConv3_min    float64
	TempConv3_max    float64
	TempConv3_stddev float64
	TempConv3_count  float64

	TempG1L1        float64
	TempG1L1_sum    float64
	TempG1L1_min    float64
	TempG1L1_max    float64
	TempG1L1_stddev float64
	TempG1L1_count  float64

	TempG1L2        float64
	TempG1L2_sum    float64
	TempG1L2_min    float64
	TempG1L2_max    float64
	TempG1L2_stddev float64
	TempG1L2_count  float64

	TempG1L3        float64
	TempG1L3_sum    float64
	TempG1L3_min    float64
	TempG1L3_max    float64
	TempG1L3_stddev float64
	TempG1L3_count  float64

	TempGearBoxHSSDE        float64
	TempGearBoxHSSDE_sum    float64
	TempGearBoxHSSDE_min    float64
	TempGearBoxHSSDE_max    float64
	TempGearBoxHSSDE_stddev float64
	TempGearBoxHSSDE_count  float64

	TempGearBoxHSSNDE        float64
	TempGearBoxHSSNDE_sum    float64
	TempGearBoxHSSNDE_min    float64
	TempGearBoxHSSNDE_max    float64
	TempGearBoxHSSNDE_stddev float64
	TempGearBoxHSSNDE_count  float64

	TempGearBoxIMSDE        float64
	TempGearBoxIMSDE_sum    float64
	TempGearBoxIMSDE_min    float64
	TempGearBoxIMSDE_max    float64
	TempGearBoxIMSDE_stddev float64
	TempGearBoxIMSDE_count  float64

	TempGearBoxIMSNDE        float64
	TempGearBoxIMSNDE_sum    float64
	TempGearBoxIMSNDE_min    float64
	TempGearBoxIMSNDE_max    float64
	TempGearBoxIMSNDE_stddev float64
	TempGearBoxIMSNDE_count  float64

	TempGearBoxOilSump        float64
	TempGearBoxOilSump_sum    float64
	TempGearBoxOilSump_min    float64
	TempGearBoxOilSump_max    float64
	TempGearBoxOilSump_stddev float64
	TempGearBoxOilSump_count  float64

	TempGeneratorBearingDE        float64
	TempGeneratorBearingDE_sum    float64
	TempGeneratorBearingDE_min    float64
	TempGeneratorBearingDE_max    float64
	TempGeneratorBearingDE_stddev float64
	TempGeneratorBearingDE_count  float64

	TempGeneratorBearingNDE        float64
	TempGeneratorBearingNDE_sum    float64
	TempGeneratorBearingNDE_min    float64
	TempGeneratorBearingNDE_max    float64
	TempGeneratorBearingNDE_stddev float64
	TempGeneratorBearingNDE_count  float64

	TempHubBearing        float64
	TempHubBearing_sum    float64
	TempHubBearing_min    float64
	TempHubBearing_max    float64
	TempHubBearing_stddev float64
	TempHubBearing_count  float64

	TempNacelle        float64
	TempNacelle_sum    float64
	TempNacelle_min    float64
	TempNacelle_max    float64
	TempNacelle_stddev float64
	TempNacelle_count  float64

	TempOutdoor        float64
	TempOutdoor_sum    float64
	TempOutdoor_min    float64
	TempOutdoor_max    float64
	TempOutdoor_stddev float64
	TempOutdoor_count  float64

	Total_Access_hrs        float64
	Total_Access_hrs_sum    float64
	Total_Access_hrs_min    float64
	Total_Access_hrs_max    float64
	Total_Access_hrs_stddev float64
	Total_Access_hrs_count  float64

	Total_Grid_OK_hrs        float64
	Total_Grid_OK_hrs_sum    float64
	Total_Grid_OK_hrs_min    float64
	Total_Grid_OK_hrs_max    float64
	Total_Grid_OK_hrs_stddev float64
	Total_Grid_OK_hrs_count  float64

	Total_Operating_hrs        float64
	Total_Operating_hrs_sum    float64
	Total_Operating_hrs_min    float64
	Total_Operating_hrs_max    float64
	Total_Operating_hrs_stddev float64
	Total_Operating_hrs_count  float64

	Total_Prod_Day_kWh        float64
	Total_Prod_Day_kWh_sum    float64
	Total_Prod_Day_kWh_min    float64
	Total_Prod_Day_kWh_max    float64
	Total_Prod_Day_kWh_stddev float64
	Total_Prod_Day_kWh_count  float64

	Total_Prod_Month_kWh        float64
	Total_Prod_Month_kWh_sum    float64
	Total_Prod_Month_kWh_min    float64
	Total_Prod_Month_kWh_max    float64
	Total_Prod_Month_kWh_stddev float64
	Total_Prod_Month_kWh_count  float64

	Total_Production_kWh        float64
	Total_Production_kWh_sum    float64
	Total_Production_kWh_min    float64
	Total_Production_kWh_max    float64
	Total_Production_kWh_stddev float64
	Total_Production_kWh_count  float64

	Total_WTG_OK_hrs        float64
	Total_WTG_OK_hrs_sum    float64
	Total_WTG_OK_hrs_min    float64
	Total_WTG_OK_hrs_max    float64
	Total_WTG_OK_hrs_stddev float64
	Total_WTG_OK_hrs_count  float64

	TotalActPowerIn_kWh        float64
	TotalActPowerIn_kWh_sum    float64
	TotalActPowerIn_kWh_min    float64
	TotalActPowerIn_kWh_max    float64
	TotalActPowerIn_kWh_stddev float64
	TotalActPowerIn_kWh_count  float64

	TotalActPowerInG1_kWh        float64
	TotalActPowerInG1_kWh_sum    float64
	TotalActPowerInG1_kWh_min    float64
	TotalActPowerInG1_kWh_max    float64
	TotalActPowerInG1_kWh_stddev float64
	TotalActPowerInG1_kWh_count  float64

	TotalActPowerInG2_kWh        float64
	TotalActPowerInG2_kWh_sum    float64
	TotalActPowerInG2_kWh_min    float64
	TotalActPowerInG2_kWh_max    float64
	TotalActPowerInG2_kWh_stddev float64
	TotalActPowerInG2_kWh_count  float64

	TotalActPowerOut_kWh        float64
	TotalActPowerOut_kWh_sum    float64
	TotalActPowerOut_kWh_min    float64
	TotalActPowerOut_kWh_max    float64
	TotalActPowerOut_kWh_stddev float64
	TotalActPowerOut_kWh_count  float64

	TotalActPowerOutG1_kWh        float64
	TotalActPowerOutG1_kWh_sum    float64
	TotalActPowerOutG1_kWh_min    float64
	TotalActPowerOutG1_kWh_max    float64
	TotalActPowerOutG1_kWh_stddev float64
	TotalActPowerOutG1_kWh_count  float64

	TotalActPowerOutG2_kWh        float64
	TotalActPowerOutG2_kWh_sum    float64
	TotalActPowerOutG2_kWh_min    float64
	TotalActPowerOutG2_kWh_max    float64
	TotalActPowerOutG2_kWh_stddev float64
	TotalActPowerOutG2_kWh_count  float64

	TotalG1ActiveHours        float64
	TotalG1ActiveHours_sum    float64
	TotalG1ActiveHours_min    float64
	TotalG1ActiveHours_max    float64
	TotalG1ActiveHours_stddev float64
	TotalG1ActiveHours_count  float64

	TotalG2ActiveHours        float64
	TotalG2ActiveHours_sum    float64
	TotalG2ActiveHours_min    float64
	TotalG2ActiveHours_max    float64
	TotalG2ActiveHours_stddev float64
	TotalG2ActiveHours_count  float64

	TotalGridOkHours        float64
	TotalGridOkHours_sum    float64
	TotalGridOkHours_min    float64
	TotalGridOkHours_max    float64
	TotalGridOkHours_stddev float64
	TotalGridOkHours_count  float64

	TotalReactPowerIn_kVArh        float64
	TotalReactPowerIn_kVArh_sum    float64
	TotalReactPowerIn_kVArh_min    float64
	TotalReactPowerIn_kVArh_max    float64
	TotalReactPowerIn_kVArh_stddev float64
	TotalReactPowerIn_kVArh_count  float64

	TotalReactPowerInG1_kVArh        float64
	TotalReactPowerInG1_kVArh_sum    float64
	TotalReactPowerInG1_kVArh_min    float64
	TotalReactPowerInG1_kVArh_max    float64
	TotalReactPowerInG1_kVArh_stddev float64
	TotalReactPowerInG1_kVArh_count  float64

	TotalReactPowerInG2_kVArh        float64
	TotalReactPowerInG2_kVArh_sum    float64
	TotalReactPowerInG2_kVArh_min    float64
	TotalReactPowerInG2_kVArh_max    float64
	TotalReactPowerInG2_kVArh_stddev float64
	TotalReactPowerInG2_kVArh_count  float64

	TotalReactPowerOut_kVArh        float64
	TotalReactPowerOut_kVArh_sum    float64
	TotalReactPowerOut_kVArh_min    float64
	TotalReactPowerOut_kVArh_max    float64
	TotalReactPowerOut_kVArh_stddev float64
	TotalReactPowerOut_kVArh_count  float64

	TotalReactPowerOutG1_kVArh        float64
	TotalReactPowerOutG1_kVArh_sum    float64
	TotalReactPowerOutG1_kVArh_min    float64
	TotalReactPowerOutG1_kVArh_max    float64
	TotalReactPowerOutG1_kVArh_stddev float64
	TotalReactPowerOutG1_kVArh_count  float64

	TotalReactPowerOutG2_kVArh        float64
	TotalReactPowerOutG2_kVArh_sum    float64
	TotalReactPowerOutG2_kVArh_min    float64
	TotalReactPowerOutG2_kVArh_max    float64
	TotalReactPowerOutG2_kVArh_stddev float64
	TotalReactPowerOutG2_kVArh_count  float64

	TotalTurbineActiveHours        float64
	TotalTurbineActiveHours_sum    float64
	TotalTurbineActiveHours_min    float64
	TotalTurbineActiveHours_max    float64
	TotalTurbineActiveHours_stddev float64
	TotalTurbineActiveHours_count  float64

	TotalTurbineOKHours        float64
	TotalTurbineOKHours_sum    float64
	TotalTurbineOKHours_min    float64
	TotalTurbineOKHours_max    float64
	TotalTurbineOKHours_stddev float64
	TotalTurbineOKHours_count  float64

	TotalTurbineTimeAllHours        float64
	TotalTurbineTimeAllHours_sum    float64
	TotalTurbineTimeAllHours_min    float64
	TotalTurbineTimeAllHours_max    float64
	TotalTurbineTimeAllHours_stddev float64
	TotalTurbineTimeAllHours_count  float64

	TurbineState        float64
	TurbineState_sum    float64
	TurbineState_min    float64
	TurbineState_max    float64
	TurbineState_stddev float64
	TurbineState_count  float64

	UTCoffset        float64
	UTCoffset_sum    float64
	UTCoffset_min    float64
	UTCoffset_max    float64
	UTCoffset_stddev float64
	UTCoffset_count  float64

	UTCoffset_int        float64
	UTCoffset_int_sum    float64
	UTCoffset_int_min    float64
	UTCoffset_int_max    float64
	UTCoffset_int_stddev float64
	UTCoffset_int_count  float64

	VoltageL1        float64
	VoltageL1_sum    float64
	VoltageL1_min    float64
	VoltageL1_max    float64
	VoltageL1_stddev float64
	VoltageL1_count  float64

	VoltageL2        float64
	VoltageL2_sum    float64
	VoltageL2_min    float64
	VoltageL2_max    float64
	VoltageL2_stddev float64
	VoltageL2_count  float64

	VoltageL3        float64
	VoltageL3_sum    float64
	VoltageL3_min    float64
	VoltageL3_max    float64
	VoltageL3_stddev float64
	VoltageL3_count  float64

	WindDirection        float64
	WindDirection_sum    float64
	WindDirection_min    float64
	WindDirection_max    float64
	WindDirection_stddev float64
	WindDirection_count  float64

	WindSpeed_ms        float64
	WindSpeed_ms_sum    float64
	WindSpeed_ms_min    float64
	WindSpeed_ms_max    float64
	WindSpeed_ms_stddev float64
	WindSpeed_ms_count  float64

	YawAngle        float64
	YawAngle_sum    float64
	YawAngle_min    float64
	YawAngle_max    float64
	YawAngle_stddev float64
	YawAngle_count  float64

	YawService        float64
	YawService_sum    float64
	YawService_min    float64
	YawService_max    float64
	YawService_stddev float64
	YawService_count  float64
}

func (m *Scada10Min) New() *Scada10Min {
	m.Id = fmt.Sprintf("%s_%s_%s", m.ProjectName, m.Turbine, m.TimeStamp.Format("20060102150405"))
	return m
}

func (m *Scada10Min) RecordID() interface{} {
	if m.Id == "" {
		m.Id = fmt.Sprintf("%s_%s_%s", m.ProjectName, m.Turbine, m.TimeStamp.Format("20060102150405"))
	}
	return m.Id
}

func (m *Scada10Min) TableName() string {
	return "Scada10Min"
}

/*
type ScadaDataHFD struct {
	orm.ModelBase `bson:"-",json:"-"`
	ID            string ` bson:"_id" json:"_id" `
	TimeStamp     time.Time
	TimeStampInt  int64
	DateInfo      DateInfo
	ProjectName   string
	Turbine       string
	IsNull        bool

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
	m.ID = timeStampStr + "#" + m.ProjectName + "#" + m.Turbine
	return m
}

func (m *ScadaDataHFD) RecordID() interface{} {
	return m.ID
}

func (m *ScadaDataHFD) TableName() string {
	return "ScadaDataHFD"
}
*/

type TurbineStatus struct {
	orm.ModelBase `bson:"-" json:"-"`
	ID            string ` bson:"_id" , json:"_id" `
	ProjectName   string
	TimeUpdate    time.Time
	Status        int // 0 : down, 1 : up
	AlarmCode     int
	AlarmDesc     string
}

func (m *TurbineStatus) New() *TurbineStatus {
	return m
}

func (m *TurbineStatus) RecordID() interface{} {
	return m.ID
}

func (m *TurbineStatus) TableName() string {
	return "TurbineStatus"
}

type AlarmRawHFD struct {
	orm.ModelBase `bson:"-" json:"-"`
	ID            bson.ObjectId ` bson:"_id" json:"_id" `
	ProjectName   string
	Turbine       string
	Time          time.Time
	DateInfo      DateInfo
	AlarmCode     int
	AlarmDesc     string
	BrakeProgram  int
	BrakeType     string
}

func (m *AlarmRawHFD) New() *AlarmRawHFD {
	m.ID = bson.NewObjectId()
	return m
}

func (m *AlarmRawHFD) RecordID() interface{} {
	return m.ID
}

func (m *AlarmRawHFD) TableName() string {
	return "AlarmRawHFD"
}

type AlarmHFD struct {
	orm.ModelBase `bson:"-" json:"-"`
	ID            bson.ObjectId ` bson:"_id" json:"_id" `
	ProjectName   string
	Turbine       string
	TimeStart     time.Time
	TimeEnd       time.Time
	DateInfoStart DateInfo
	DateInfoEnd   DateInfo
	Duration      float64
	AlarmCode     int
	AlarmDesc     string
	BrakeProgram  int
	BrakeType     string
	Finish        int
}

func (m *AlarmHFD) New() *AlarmHFD {
	m.ID = bson.NewObjectId()
	return m
}

func (m *AlarmHFD) RecordID() interface{} {
	return m.ID
}

func (m *AlarmHFD) TableName() string {
	return "AlarmHFD"
}
