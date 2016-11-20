package models

import (
	. "eaciit/ostrowfm/library/helper"
	"time"

	"github.com/eaciit/orm"
	tk "github.com/eaciit/toolkit"
	"gopkg.in/mgo.v2/bson"
)

type ScadaDataOEM struct {
	orm.ModelBase                            `bson:"-",json:"-"`
	ID                                       bson.ObjectId ` bson:"_id" , json:"_id" `
	TimeStamp                                time.Time
	TimeStampUTC                             time.Time
	DateInfo                                 DateInfo
	DateInfoUTC                              DateInfo
	ProjectName                              string
	Turbine                                  string
	Line                                     int
	AI_intern_R_PidAngleOut                  float64
	AI_intern_ActivPower                     float64
	AI_intern_I1                             float64
	AI_intern_I2                             float64
	AI_intern_I3                             float64
	AI_intern_NacelleDrill                   float64
	AI_intern_NacellePos                     float64
	AI_intern_PitchAkku_V1                   float64
	AI_intern_PitchAkku_V2                   float64
	AI_intern_PitchAkku_V3                   float64
	AI_intern_PitchAngle1                    float64
	AI_intern_PitchAngle2                    float64
	AI_intern_PitchAngle3                    float64
	AI_intern_PitchConv_Current1             float64
	AI_intern_PitchConv_Current2             float64
	AI_intern_PitchConv_Current3             float64
	AI_intern_PitchAngleSP_Diff1             float64
	AI_intern_PitchAngleSP_Diff2             float64
	AI_intern_PitchAngleSP_Diff3             float64
	AI_intern_ReactivPower                   float64
	AI_intern_RpmDiff                        float64
	AI_intern_U1                             float64
	AI_intern_U2                             float64
	AI_intern_U3                             float64
	AI_intern_WindDirection                  float64
	AI_intern_WindSpeed                      float64
	AI_Intern_WindSpeedDif                   float64
	AI_speed_RotFR                           float64
	AI_WindSpeed1                            float64
	AI_WindSpeed2                            float64
	AI_WindVane1                             float64
	AI_WindVane2                             float64
	AI_internCurrentAsym                     float64
	Temp_GearBox_IMS_NDE                     float64
	AI_intern_WindVaneDiff                   float64
	C_intern_SpeedGenerator                  float64
	C_intern_SpeedRotor                      float64
	AI_intern_Speed_RPMDiff_FR1_RotCNT       float64
	AI_intern_Frequency_Grid                 float64
	Temp_GearBox_HSS_NDE                     float64
	AI_DrTrVibValue                          float64
	AI_intern_InLastErrorConv1               float64
	AI_intern_InLastErrorConv2               float64
	AI_intern_InLastErrorConv3               float64
	AI_intern_TempConv1                      float64
	AI_intern_TempConv2                      float64
	AI_intern_TempConv3                      float64
	AI_intern_PitchSpeed1                    float64
	Temp_YawBrake_1                          float64
	Temp_YawBrake_2                          float64
	Temp_G1L1                                float64
	Temp_G1L2                                float64
	Temp_G1L3                                float64
	Temp_YawBrake_3                          float64
	AI_HydrSystemPressure                    float64
	Temp_BottomControlSection_Low            float64
	Temp_GearBox_HSS_DE                      float64
	Temp_GearOilSump                         float64
	Temp_GeneratorBearing_DE                 float64
	Temp_GeneratorBearing_NDE                float64
	Temp_MainBearing                         float64
	Temp_GearBox_IMS_DE                      float64
	Temp_Nacelle                             float64
	Temp_Outdoor                             float64
	AI_TowerVibValueAxial                    float64
	AI_intern_DiffGenSpeedSPToAct            float64
	Temp_YawBrake_4                          float64
	AI_intern_SpeedGenerator_Proximity       float64
	AI_intern_SpeedDiff_Encoder_Proximity    float64
	AI_GearOilPressure                       float64
	Temp_CabinetTopBox_Low                   float64
	Temp_CabinetTopBox                       float64
	Temp_BottomControlSection                float64
	Temp_BottomPowerSection                  float64
	Temp_BottomPowerSection_Low              float64
	AI_intern_Pitch1_Status_High             float64
	AI_intern_Pitch2_Status_High             float64
	AI_intern_Pitch3_Status_High             float64
	AI_intern_InPosition1_ch2                float64
	AI_intern_InPosition2_ch2                float64
	AI_intern_InPosition3_ch2                float64
	AI_intern_Temp_Brake_Blade1              float64
	AI_intern_Temp_Brake_Blade2              float64
	AI_intern_Temp_Brake_Blade3              float64
	AI_intern_Temp_PitchMotor_Blade1         float64
	AI_intern_Temp_PitchMotor_Blade2         float64
	AI_intern_Temp_PitchMotor_Blade3         float64
	AI_intern_Temp_Hub_Additional1           float64
	AI_intern_Temp_Hub_Additional2           float64
	AI_intern_Temp_Hub_Additional3           float64
	AI_intern_Pitch1_Status_Low              float64
	AI_intern_Pitch2_Status_Low              float64
	AI_intern_Pitch3_Status_Low              float64
	AI_intern_Battery_VoltageBlade1_center   float64
	AI_intern_Battery_VoltageBlade2_center   float64
	AI_intern_Battery_VoltageBlade3_center   float64
	AI_intern_Battery_ChargingCur_Blade1     float64
	AI_intern_Battery_ChargingCur_Blade2     float64
	AI_intern_Battery_ChargingCur_Blade3     float64
	AI_intern_Battery_DischargingCur_Blade1  float64
	AI_intern_Battery_DischargingCur_Blade2  float64
	AI_intern_Battery_DischargingCur_Blade3  float64
	AI_intern_PitchMotor_BrakeVoltage_Blade1 float64
	AI_intern_PitchMotor_BrakeVoltage_Blade2 float64
	AI_intern_PitchMotor_BrakeVoltage_Blade3 float64
	AI_intern_PitchMotor_BrakeCurrent_Blade1 float64
	AI_intern_PitchMotor_BrakeCurrent_Blade2 float64
	AI_intern_PitchMotor_BrakeCurrent_Blade3 float64
	AI_intern_Temp_HubBox_Blade1             float64
	AI_intern_Temp_HubBox_Blade2             float64
	AI_intern_Temp_HubBox_Blade3             float64
	AI_intern_Temp_Pitch1_HeatSink           float64
	AI_intern_Temp_Pitch2_HeatSink           float64
	AI_intern_Temp_Pitch3_HeatSink           float64
	AI_intern_ErrorStackBlade1               float64
	AI_intern_ErrorStackBlade2               float64
	AI_intern_ErrorStackBlade3               float64
	AI_intern_Temp_BatteryBox_Blade1         float64
	AI_intern_Temp_BatteryBox_Blade2         float64
	AI_intern_Temp_BatteryBox_Blade3         float64
	AI_intern_DC_LinkVoltage1                float64
	AI_intern_DC_LinkVoltage2                float64
	AI_intern_DC_LinkVoltage3                float64
	Temp_Yaw_Motor1                          float64
	Temp_Yaw_Motor2                          float64
	Temp_Yaw_Motor3                          float64
	Temp_Yaw_Motor4                          float64
	AO_DFIG_Power_Setpiont                   float64
	AO_DFIG_Q_Setpoint                       float64
	AI_DFIG_Torque_actual                    float64
	AI_DFIG_SpeedGenerator_Encoder           float64
	AI_intern_DFIG_DC_Link_Voltage_actual    float64
	AI_intern_DFIG_MSC_current               float64
	AI_intern_DFIG_Main_voltage              float64
	AI_intern_DFIG_Main_current              float64
	AI_intern_DFIG_active_power_actual       float64
	AI_intern_DFIG_reactive_power_actual     float64
	AI_intern_DFIG_active_power_actual_LSC   float64
	AI_intern_DFIG_LSC_current               float64
	AI_intern_DFIG_Data_log_number           float64
	AI_intern_Damper_OscMagnitude            float64
	AI_intern_Damper_PassbandFullLoad        float64
	AI_YawBrake_TempRise1                    float64
	AI_YawBrake_TempRise2                    float64
	AI_YawBrake_TempRise3                    float64
	AI_YawBrake_TempRise4                    float64
	AI_intern_NacelleDrill_at_NorthPosSensor float64
}

func (m *ScadaDataOEM) New() *ScadaDataOEM {
	m.ID = bson.NewObjectId()
	return m
}

func (m *ScadaDataOEM) RecordID() interface{} {
	return m.ID
}

func (m *ScadaDataOEM) TableName() string {
	return "ScadaDataOEM"
}

func (m *ScadaDataOEM) GetXlsColumns() []string {
	columns := []string{
		"TimeStamp",
		"AI_intern_R_PidAngleOut",
		"AI_intern_ActivPower",
		"AI_intern_I1",
		"AI_intern_I2",
		"AI_intern_I3",
		"AI_intern_NacelleDrill",
		"AI_intern_NacellePos",
		"AI_intern_PitchAkku_V1",
		"AI_intern_PitchAkku_V2",
		"AI_intern_PitchAkku_V3",
		"AI_intern_PitchAngle1",
		"AI_intern_PitchAngle2",
		"AI_intern_PitchAngle3",
		"AI_intern_PitchConv_Current1",
		"AI_intern_PitchConv_Current2",
		"AI_intern_PitchConv_Current3",
		"AI_intern_PitchAngleSP_Diff1",
		"AI_intern_PitchAngleSP_Diff2",
		"AI_intern_PitchAngleSP_Diff3",
		"AI_intern_ReactivPower",
		"AI_intern_RpmDiff",
		"AI_intern_U1",
		"AI_intern_U2",
		"AI_intern_U3",
		"AI_intern_WindDirection",
		"AI_intern_WindSpeed",
		"AI_Intern_WindSpeedDif",
		"AI_speed_RotFR",
		"AI_WindSpeed1",
		"AI_WindSpeed2",
		"AI_WindVane1",
		"AI_WindVane2",
		"AI_internCurrentAsym",
		"Temp_GearBox_IMS_NDE",
		"AI_intern_WindVaneDiff",
		"C_intern_SpeedGenerator",
		"C_intern_SpeedRotor",
		"AI_intern_Speed_RPMDiff_FR1_RotCNT",
		"AI_intern_Frequency_Grid",
		"Temp_GearBox_HSS_NDE",
		"AI_DrTrVibValue",
		"AI_intern_InLastErrorConv1",
		"AI_intern_InLastErrorConv2",
		"AI_intern_InLastErrorConv3",
		"AI_intern_TempConv1",
		"AI_intern_TempConv2",
		"AI_intern_TempConv3",
		"AI_intern_PitchSpeed1",
		"Temp_YawBrake_1",
		"Temp_YawBrake_2",
		"Temp_G1L1",
		"Temp_G1L2",
		"Temp_G1L3",
		"Temp_YawBrake_3",
		"AI_HydrSystemPressure",
		"Temp_BottomControlSection_Low",
		"Temp_GearBox_HSS_DE",
		"Temp_GearOilSump",
		"Temp_GeneratorBearing_DE",
		"Temp_GeneratorBearing_NDE",
		"Temp_MainBearing",
		"Temp_GearBox_IMS_DE",
		"Temp_Nacelle",
		"Temp_Outdoor",
		"AI_TowerVibValueAxial",
		"AI_intern_DiffGenSpeedSPToAct",
		"Temp_YawBrake_4",
		"AI_intern_SpeedGenerator_Proximity",
		"AI_intern_SpeedDiff_Encoder_Proximity",
		"AI_GearOilPressure",
		"Temp_CabinetTopBox_Low",
		"Temp_CabinetTopBox",
		"Temp_BottomControlSection",
		"Temp_BottomPowerSection",
		"Temp_BottomPowerSection_Low",
		"AI_intern_Pitch1_Status_High",
		"AI_intern_Pitch2_Status_High",
		"AI_intern_Pitch3_Status_High",
		"AI_intern_InPosition1_ch2",
		"AI_intern_InPosition2_ch2",
		"AI_intern_InPosition3_ch2",
		"AI_intern_Temp_Brake_Blade1",
		"AI_intern_Temp_Brake_Blade2",
		"AI_intern_Temp_Brake_Blade3",
		"AI_intern_Temp_PitchMotor_Blade1",
		"AI_intern_Temp_PitchMotor_Blade2",
		"AI_intern_Temp_PitchMotor_Blade3",
		"AI_intern_Temp_Hub_Additional1",
		"AI_intern_Temp_Hub_Additional2",
		"AI_intern_Temp_Hub_Additional3",
		"AI_intern_Pitch1_Status_Low",
		"AI_intern_Pitch2_Status_Low",
		"AI_intern_Pitch3_Status_Low",
		"AI_intern_Battery_VoltageBlade1_center",
		"AI_intern_Battery_VoltageBlade2_center",
		"AI_intern_Battery_VoltageBlade3_center",
		"AI_intern_Battery_ChargingCur_Blade1",
		"AI_intern_Battery_ChargingCur_Blade2",
		"AI_intern_Battery_ChargingCur_Blade3",
		"AI_intern_Battery_DischargingCur_Blade1",
		"AI_intern_Battery_DischargingCur_Blade2",
		"AI_intern_Battery_DischargingCur_Blade3",
		"AI_intern_PitchMotor_BrakeVoltage_Blade1",
		"AI_intern_PitchMotor_BrakeVoltage_Blade2",
		"AI_intern_PitchMotor_BrakeVoltage_Blade3",
		"AI_intern_PitchMotor_BrakeCurrent_Blade1",
		"AI_intern_PitchMotor_BrakeCurrent_Blade2",
		"AI_intern_PitchMotor_BrakeCurrent_Blade3",
		"AI_intern_Temp_HubBox_Blade1",
		"AI_intern_Temp_HubBox_Blade2",
		"AI_intern_Temp_HubBox_Blade3",
		"AI_intern_Temp_Pitch1_HeatSink",
		"AI_intern_Temp_Pitch2_HeatSink",
		"AI_intern_Temp_Pitch3_HeatSink",
		"AI_intern_ErrorStackBlade1",
		"AI_intern_ErrorStackBlade2",
		"AI_intern_ErrorStackBlade3",
		"AI_intern_Temp_BatteryBox_Blade1",
		"AI_intern_Temp_BatteryBox_Blade2",
		"AI_intern_Temp_BatteryBox_Blade3",
		"AI_intern_DC_LinkVoltage1",
		"AI_intern_DC_LinkVoltage2",
		"AI_intern_DC_LinkVoltage3",
		"Temp_Yaw_Motor1",
		"Temp_Yaw_Motor2",
		"Temp_Yaw_Motor3",
		"Temp_Yaw_Motor4",
		"AO_DFIG_Power_Setpiont",
		"AO_DFIG_Q_Setpoint",
		"AI_DFIG_Torque_actual",
		"AI_DFIG_SpeedGenerator_Encoder",
		"AI_intern_DFIG_DC_Link_Voltage_actual",
		"AI_intern_DFIG_MSC_current",
		"AI_intern_DFIG_Main_voltage",
		"AI_intern_DFIG_Main_current",
		"AI_intern_DFIG_active_power_actual",
		"AI_intern_DFIG_reactive_power_actual",
		"AI_intern_DFIG_active_power_actual_LSC",
		"AI_intern_DFIG_LSC_current",
		"AI_intern_DFIG_Data_log_number",
		"AI_intern_Damper_OscMagnitude",
		"AI_intern_Damper_PassbandFullLoad",
		"AI_YawBrake_TempRise1",
		"AI_YawBrake_TempRise2",
		"AI_YawBrake_TempRise3",
		"AI_YawBrake_TempRise4",
		"AI_intern_NacelleDrill_at_NorthPosSensor",
	}

	return columns
}

func (m *ScadaDataOEM) GetColumnInfos(colName string) tk.M {
	colInfos := tk.M{
		"TimeStamp":                                tk.M{"Column": "A", "Index": 0},
		"AI_intern_R_PidAngleOut":                  tk.M{"Column": "B", "Index": 1},
		"AI_intern_ActivPower":                     tk.M{"Column": "C", "Index": 2},
		"AI_intern_I1":                             tk.M{"Column": "D", "Index": 3},
		"AI_intern_I2":                             tk.M{"Column": "E", "Index": 4},
		"AI_intern_I3":                             tk.M{"Column": "F", "Index": 5},
		"AI_intern_NacelleDrill":                   tk.M{"Column": "G", "Index": 6},
		"AI_intern_NacellePos":                     tk.M{"Column": "H", "Index": 7},
		"AI_intern_PitchAkku_V1":                   tk.M{"Column": "I", "Index": 8},
		"AI_intern_PitchAkku_V2":                   tk.M{"Column": "J", "Index": 9},
		"AI_intern_PitchAkku_V3":                   tk.M{"Column": "K", "Index": 10},
		"AI_intern_PitchAngle1":                    tk.M{"Column": "L", "Index": 11},
		"AI_intern_PitchAngle2":                    tk.M{"Column": "M", "Index": 12},
		"AI_intern_PitchAngle3":                    tk.M{"Column": "N", "Index": 13},
		"AI_intern_PitchConv_Current1":             tk.M{"Column": "O", "Index": 14},
		"AI_intern_PitchConv_Current2":             tk.M{"Column": "P", "Index": 15},
		"AI_intern_PitchConv_Current3":             tk.M{"Column": "Q", "Index": 16},
		"AI_intern_PitchAngleSP_Diff1":             tk.M{"Column": "R", "Index": 17},
		"AI_intern_PitchAngleSP_Diff2":             tk.M{"Column": "S", "Index": 18},
		"AI_intern_PitchAngleSP_Diff3":             tk.M{"Column": "T", "Index": 19},
		"AI_intern_ReactivPower":                   tk.M{"Column": "U", "Index": 20},
		"AI_intern_RpmDiff":                        tk.M{"Column": "V", "Index": 21},
		"AI_intern_U1":                             tk.M{"Column": "W", "Index": 22},
		"AI_intern_U2":                             tk.M{"Column": "X", "Index": 23},
		"AI_intern_U3":                             tk.M{"Column": "Y", "Index": 24},
		"AI_intern_WindDirection":                  tk.M{"Column": "Z", "Index": 25},
		"AI_intern_WindSpeed":                      tk.M{"Column": "AA", "Index": 26},
		"AI_Intern_WindSpeedDif":                   tk.M{"Column": "AB", "Index": 27},
		"AI_speed_RotFR":                           tk.M{"Column": "AC", "Index": 28},
		"AI_WindSpeed1":                            tk.M{"Column": "AD", "Index": 29},
		"AI_WindSpeed2":                            tk.M{"Column": "AE", "Index": 30},
		"AI_WindVane1":                             tk.M{"Column": "AF", "Index": 31},
		"AI_WindVane2":                             tk.M{"Column": "AG", "Index": 32},
		"AI_internCurrentAsym":                     tk.M{"Column": "AH", "Index": 33},
		"Temp_GearBox_IMS_NDE":                     tk.M{"Column": "AI", "Index": 34},
		"AI_intern_WindVaneDiff":                   tk.M{"Column": "AJ", "Index": 35},
		"C_intern_SpeedGenerator":                  tk.M{"Column": "AK", "Index": 36},
		"C_intern_SpeedRotor":                      tk.M{"Column": "AL", "Index": 37},
		"AI_intern_Speed_RPMDiff_FR1_RotCNT":       tk.M{"Column": "AM", "Index": 38},
		"AI_intern_Frequency_Grid":                 tk.M{"Column": "AN", "Index": 39},
		"Temp_GearBox_HSS_NDE":                     tk.M{"Column": "AO", "Index": 40},
		"AI_DrTrVibValue":                          tk.M{"Column": "AP", "Index": 41},
		"AI_intern_InLastErrorConv1":               tk.M{"Column": "AQ", "Index": 42},
		"AI_intern_InLastErrorConv2":               tk.M{"Column": "AR", "Index": 43},
		"AI_intern_InLastErrorConv3":               tk.M{"Column": "AS", "Index": 44},
		"AI_intern_TempConv1":                      tk.M{"Column": "AT", "Index": 45},
		"AI_intern_TempConv2":                      tk.M{"Column": "AU", "Index": 46},
		"AI_intern_TempConv3":                      tk.M{"Column": "AV", "Index": 47},
		"AI_intern_PitchSpeed1":                    tk.M{"Column": "AW", "Index": 48},
		"Temp_YawBrake_1":                          tk.M{"Column": "AX", "Index": 49},
		"Temp_YawBrake_2":                          tk.M{"Column": "AY", "Index": 50},
		"Temp_G1L1":                                tk.M{"Column": "AZ", "Index": 51},
		"Temp_G1L2":                                tk.M{"Column": "BA", "Index": 52},
		"Temp_G1L3":                                tk.M{"Column": "BB", "Index": 53},
		"Temp_YawBrake_3":                          tk.M{"Column": "BC", "Index": 54},
		"AI_HydrSystemPressure":                    tk.M{"Column": "BD", "Index": 55},
		"Temp_BottomControlSection_Low":            tk.M{"Column": "BE", "Index": 56},
		"Temp_GearBox_HSS_DE":                      tk.M{"Column": "BF", "Index": 57},
		"Temp_GearOilSump":                         tk.M{"Column": "BG", "Index": 58},
		"Temp_GeneratorBearing_DE":                 tk.M{"Column": "BH", "Index": 59},
		"Temp_GeneratorBearing_NDE":                tk.M{"Column": "BI", "Index": 60},
		"Temp_MainBearing":                         tk.M{"Column": "BJ", "Index": 61},
		"Temp_GearBox_IMS_DE":                      tk.M{"Column": "BK", "Index": 62},
		"Temp_Nacelle":                             tk.M{"Column": "BL", "Index": 63},
		"Temp_Outdoor":                             tk.M{"Column": "BM", "Index": 64},
		"AI_TowerVibValueAxial":                    tk.M{"Column": "BN", "Index": 65},
		"AI_intern_DiffGenSpeedSPToAct":            tk.M{"Column": "BO", "Index": 66},
		"Temp_YawBrake_4":                          tk.M{"Column": "BP", "Index": 67},
		"AI_intern_SpeedGenerator_Proximity":       tk.M{"Column": "BQ", "Index": 68},
		"AI_intern_SpeedDiff_Encoder_Proximity":    tk.M{"Column": "BR", "Index": 69},
		"AI_GearOilPressure":                       tk.M{"Column": "BS", "Index": 70},
		"Temp_CabinetTopBox_Low":                   tk.M{"Column": "BT", "Index": 71},
		"Temp_CabinetTopBox":                       tk.M{"Column": "BU", "Index": 72},
		"Temp_BottomControlSection":                tk.M{"Column": "BV", "Index": 73},
		"Temp_BottomPowerSection":                  tk.M{"Column": "BW", "Index": 74},
		"Temp_BottomPowerSection_Low":              tk.M{"Column": "BX", "Index": 75},
		"AI_intern_Pitch1_Status_High":             tk.M{"Column": "BY", "Index": 76},
		"AI_intern_Pitch2_Status_High":             tk.M{"Column": "BZ", "Index": 77},
		"AI_intern_Pitch3_Status_High":             tk.M{"Column": "CA", "Index": 78},
		"AI_intern_InPosition1_ch2":                tk.M{"Column": "CB", "Index": 79},
		"AI_intern_InPosition2_ch2":                tk.M{"Column": "CC", "Index": 80},
		"AI_intern_InPosition3_ch2":                tk.M{"Column": "CD", "Index": 81},
		"AI_intern_Temp_Brake_Blade1":              tk.M{"Column": "CE", "Index": 82},
		"AI_intern_Temp_Brake_Blade2":              tk.M{"Column": "CF", "Index": 83},
		"AI_intern_Temp_Brake_Blade3":              tk.M{"Column": "CG", "Index": 84},
		"AI_intern_Temp_PitchMotor_Blade1":         tk.M{"Column": "CH", "Index": 85},
		"AI_intern_Temp_PitchMotor_Blade2":         tk.M{"Column": "CI", "Index": 86},
		"AI_intern_Temp_PitchMotor_Blade3":         tk.M{"Column": "CJ", "Index": 87},
		"AI_intern_Temp_Hub_Additional1":           tk.M{"Column": "CK", "Index": 88},
		"AI_intern_Temp_Hub_Additional2":           tk.M{"Column": "CL", "Index": 89},
		"AI_intern_Temp_Hub_Additional3":           tk.M{"Column": "CM", "Index": 90},
		"AI_intern_Pitch1_Status_Low":              tk.M{"Column": "CN", "Index": 91},
		"AI_intern_Pitch2_Status_Low":              tk.M{"Column": "CO", "Index": 92},
		"AI_intern_Pitch3_Status_Low":              tk.M{"Column": "CP", "Index": 93},
		"AI_intern_Battery_VoltageBlade1_center":   tk.M{"Column": "CQ", "Index": 94},
		"AI_intern_Battery_VoltageBlade2_center":   tk.M{"Column": "CR", "Index": 95},
		"AI_intern_Battery_VoltageBlade3_center":   tk.M{"Column": "CS", "Index": 96},
		"AI_intern_Battery_ChargingCur_Blade1":     tk.M{"Column": "CT", "Index": 97},
		"AI_intern_Battery_ChargingCur_Blade2":     tk.M{"Column": "CU", "Index": 98},
		"AI_intern_Battery_ChargingCur_Blade3":     tk.M{"Column": "CV", "Index": 99},
		"AI_intern_Battery_DischargingCur_Blade1":  tk.M{"Column": "CW", "Index": 100},
		"AI_intern_Battery_DischargingCur_Blade2":  tk.M{"Column": "CX", "Index": 101},
		"AI_intern_Battery_DischargingCur_Blade3":  tk.M{"Column": "CY", "Index": 102},
		"AI_intern_PitchMotor_BrakeVoltage_Blade1": tk.M{"Column": "CZ", "Index": 103},
		"AI_intern_PitchMotor_BrakeVoltage_Blade2": tk.M{"Column": "DA", "Index": 104},
		"AI_intern_PitchMotor_BrakeVoltage_Blade3": tk.M{"Column": "DB", "Index": 105},
		"AI_intern_PitchMotor_BrakeCurrent_Blade1": tk.M{"Column": "DC", "Index": 106},
		"AI_intern_PitchMotor_BrakeCurrent_Blade2": tk.M{"Column": "DD", "Index": 107},
		"AI_intern_PitchMotor_BrakeCurrent_Blade3": tk.M{"Column": "DE", "Index": 108},
		"AI_intern_Temp_HubBox_Blade1":             tk.M{"Column": "DF", "Index": 109},
		"AI_intern_Temp_HubBox_Blade2":             tk.M{"Column": "DG", "Index": 110},
		"AI_intern_Temp_HubBox_Blade3":             tk.M{"Column": "DH", "Index": 111},
		"AI_intern_Temp_Pitch1_HeatSink":           tk.M{"Column": "DI", "Index": 112},
		"AI_intern_Temp_Pitch2_HeatSink":           tk.M{"Column": "DJ", "Index": 113},
		"AI_intern_Temp_Pitch3_HeatSink":           tk.M{"Column": "DK", "Index": 114},
		"AI_intern_ErrorStackBlade1":               tk.M{"Column": "DL", "Index": 115},
		"AI_intern_ErrorStackBlade2":               tk.M{"Column": "DM", "Index": 116},
		"AI_intern_ErrorStackBlade3":               tk.M{"Column": "DN", "Index": 117},
		"AI_intern_Temp_BatteryBox_Blade1":         tk.M{"Column": "DO", "Index": 118},
		"AI_intern_Temp_BatteryBox_Blade2":         tk.M{"Column": "DP", "Index": 119},
		"AI_intern_Temp_BatteryBox_Blade3":         tk.M{"Column": "DQ", "Index": 120},
		"AI_intern_DC_LinkVoltage1":                tk.M{"Column": "DR", "Index": 121},
		"AI_intern_DC_LinkVoltage2":                tk.M{"Column": "DS", "Index": 122},
		"AI_intern_DC_LinkVoltage3":                tk.M{"Column": "DT", "Index": 123},
		"Temp_Yaw_Motor1":                          tk.M{"Column": "DU", "Index": 124},
		"Temp_Yaw_Motor2":                          tk.M{"Column": "DV", "Index": 125},
		"Temp_Yaw_Motor3":                          tk.M{"Column": "DW", "Index": 126},
		"Temp_Yaw_Motor4":                          tk.M{"Column": "DX", "Index": 127},
		"AO_DFIG_Power_Setpiont":                   tk.M{"Column": "DY", "Index": 128},
		"AO_DFIG_Q_Setpoint":                       tk.M{"Column": "DZ", "Index": 129},
		"AI_DFIG_Torque_actual":                    tk.M{"Column": "EA", "Index": 130},
		"AI_DFIG_SpeedGenerator_Encoder":           tk.M{"Column": "EB", "Index": 131},
		"AI_intern_DFIG_DC_Link_Voltage_actual":    tk.M{"Column": "EC", "Index": 132},
		"AI_intern_DFIG_MSC_current":               tk.M{"Column": "ED", "Index": 133},
		"AI_intern_DFIG_Main_voltage":              tk.M{"Column": "EE", "Index": 134},
		"AI_intern_DFIG_Main_current":              tk.M{"Column": "EF", "Index": 135},
		"AI_intern_DFIG_active_power_actual":       tk.M{"Column": "EG", "Index": 136},
		"AI_intern_DFIG_reactive_power_actual":     tk.M{"Column": "EH", "Index": 137},
		"AI_intern_DFIG_active_power_actual_LSC":   tk.M{"Column": "EI", "Index": 138},
		"AI_intern_DFIG_LSC_current":               tk.M{"Column": "EJ", "Index": 139},
		"AI_intern_DFIG_Data_log_number":           tk.M{"Column": "EK", "Index": 140},
		"AI_intern_Damper_OscMagnitude":            tk.M{"Column": "EL", "Index": 141},
		"AI_intern_Damper_PassbandFullLoad":        tk.M{"Column": "EM", "Index": 142},
		"AI_YawBrake_TempRise1":                    tk.M{"Column": "EN", "Index": 143},
		"AI_YawBrake_TempRise2":                    tk.M{"Column": "EO", "Index": 144},
		"AI_YawBrake_TempRise3":                    tk.M{"Column": "EP", "Index": 145},
		"AI_YawBrake_TempRise4":                    tk.M{"Column": "EQ", "Index": 146},
		"AI_intern_NacelleDrill_at_NorthPosSensor": tk.M{"Column": "ER", "Index": 147},
	}

	return colInfos[colName].(tk.M)
}
