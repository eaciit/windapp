package controller

import (
	"strings"
)

func GetScadaOEMHeader() (headerResult []string, fieldResult []string) {
	oldHeader := `AI intern ActivPower
AI intern WindSpeed
AI intern NacellePos
AI intern WindDirection
AI intern PitchAngle1
AI intern PitchAngle2
AI intern PitchAngle3
C intern SpeedGenerator
C intern SpeedRotor
AI intern ReactivPower
AI intern Frequency Grid
AI GearOilPressure
Temp Outdoor
Temp Nacelle
Temp GearBox HSS NDE
Temp GearBox HSS DE
Temp GearBox IMS DE
Temp GearOilSump
Temp GearBox IMS NDE
Temp GeneratorBearing DE
Temp GeneratorBearing NDE
Temp MainBearing
Temp YawBrake 1
Temp YawBrake 2
Temp G1L1
Temp G1L2
Temp G1L3
Temp YawBrake 4
AI HydrSystemPressure
Temp BottomControlSection Low
AI intern TempConv1
AI intern TempConv2
AI intern TempConv3
AI intern R PidAngleOut
AI intern I1
AI intern I2
AI intern I3
AI intern NacelleDrill
AI intern PitchAkku V1
AI intern PitchAkku V2
AI intern PitchAkku V3
AI intern PitchConv Current1
AI intern PitchConv Current2
AI intern PitchConv Current3
AI intern PitchAngleSP Diff1
AI intern PitchAngleSP Diff2
AI intern PitchAngleSP Diff3
AI intern RpmDiff
AI intern U1
AI intern U2
AI intern U3
AI Intern WindSpeedDif
AI speed RotFR
AI WindSpeed1
AI WindSpeed2
AI WindVane1
AI WindVane2
AI internCurrentAsym
AI intern WindVaneDiff
AI intern PitchSpeed2
AI intern Speed RPMDiff FR1 RotCNT
AI DrTrVibValue
AI intern InLastErrorConv1
AI intern InLastErrorConv2
AI intern InLastErrorConv3
AI TowerVibValueAxial
AI intern DiffGenSpeedSPToAct
Temp YawBrake 5
AI intern SpeedGenerator Proximity
AI intern SpeedDiff Encoder Proximity
Temp CabinetTopBox Low
Temp CabinetTopBox
Temp BottomControlSection
Temp BottomPowerSection
Temp BottomPowerSection Low
AI intern Pitch1 Status High
AI intern Pitch2 Status High
AI intern Pitch3 Status High
AI intern InPosition1 ch3
AI intern InPosition2 ch3
AI intern InPosition3 ch3
AI intern Temp Brake Blade1
AI intern Temp Brake Blade2
AI intern Temp Brake Blade3
AI intern Temp PitchMotor Blade1
AI intern Temp PitchMotor Blade2
AI intern Temp PitchMotor Blade3
AI intern Temp Hub Additional1
AI intern Temp Hub Additional2
AI intern Temp Hub Additional3
AI intern Pitch1 Status Low
AI intern Pitch2 Status Low
AI intern Pitch3 Status Low
AI intern Battery VoltageBlade1 center
AI intern Battery VoltageBlade2 center
AI intern Battery VoltageBlade3 center
AI intern Battery ChargingCur Blade1
AI intern Battery ChargingCur Blade2
AI intern Battery ChargingCur Blade3
AI intern Battery DischargingCur Blade1
AI intern Battery DischargingCur Blade2
AI intern Battery DischargingCur Blade3
AI intern PitchMotor BrakeVoltage Blade1
AI intern PitchMotor BrakeVoltage Blade2
AI intern PitchMotor BrakeVoltage Blade3
AI intern PitchMotor BrakeCurrent Blade1
AI intern PitchMotor BrakeCurrent Blade2
AI intern PitchMotor BrakeCurrent Blade3
AI intern Temp HubBox Blade1
AI intern Temp HubBox Blade2
AI intern Temp HubBox Blade3
AI intern Temp Pitch1 HeatSink
AI intern Temp Pitch2 HeatSink
AI intern Temp Pitch3 HeatSink
AI intern ErrorStackBlade1
AI intern ErrorStackBlade2
AI intern ErrorStackBlade3
AI intern Temp BatteryBox Blade1
AI intern Temp BatteryBox Blade2
AI intern Temp BatteryBox Blade3
AI intern DC LinkVoltage1
AI intern DC LinkVoltage2
AI intern DC LinkVoltage3
Temp Yaw Motor1
Temp Yaw Motor2
Temp Yaw Motor3
Temp Yaw Motor4
AO DFIG Power Setpiont
AO DFIG Q Setpoint
AI DFIG Torque actual
AI DFIG SpeedGenerator Encoder
AI intern DFIG DC Link Voltage actual
AI intern DFIG MSC current
AI intern DFIG Main voltage
AI intern DFIG Main current
AI intern DFIG active power actual
AI intern DFIG reactive power actual
AI intern DFIG active power actual LSC
AI intern DFIG LSC current
AI intern DFIG Data log number
AI intern Damper OscMagnitude
AI intern Damper PassbandFullLoad
AI YawBrake TempRise1
AI YawBrake TempRise2
AI YawBrake TempRise3
AI YawBrake TempRise4
AI intern NacelleDrill at NorthPosSensor`

	newHeader := `Active Power
Wind Speed
Nacelle Pos
Wind Direction
Pitch Angle1
Pitch Angle2
Pitch Angle3
Generator Speed
Rotor Speed
Reactive Power
Frequency Grid
Gear Oil Pressure
Ambient Temp
Temp Nacelle
Temp GearBox HSS NDE
Temp GearBox HSS DE
Temp GearBox IMS DE
Temp Gear Oil Sump
Temp GearBox IMS NDE
Temp Generator Bearing DE
Temp Generator Bearing NDE
Temp Main Bearing
Temp Yaw Brake 1
Temp Yaw Brake 2
Temp G1L1
Temp G1L2
Temp G1L3
Temp Yaw Brake 4
Hydr System Pressure
Temp Bottom Control Section Low
Temp Conv1
Temp Conv2
Temp Conv3
R Pid Angle Out
I1
I2
I3
Nacelle Drill
Pitch Akku V1
Pitch Akku V2
Pitch Akku V3
Pitch Conv Current1
Pitch Conv Current2
Pitch Conv Current3
Pitch Angle SP Diff1
Pitch Angle SP Diff2
Pitch Angle SP Diff3
Rpm Diff
U1
U2
U3
Wind Speed Dif
Speed Rot FR
Wind Speed1
Wind Speed2
Wind Vane1
Wind Vane2
Current Asym
Wind Vane Diff
Pitch Speed2
Speed RPM Diff FR1 Rot CNT
Dr Tr Vib Value
In Last Error Conv1
In Last Error Conv2
In Last Error Conv3
AI Tower Vib Value Axial
Diff Gen Speed SP To Act
Temp YawBrake 5
Speed Generator Proximity
Speed Diff Encoder Proximity
Temp CabinetTopBox Low
Temp CabinetTopBox
Temp Bottom Control Section
Temp Bottom Power Section
Temp Bottom Power Section Low
Pitch1 Status High
Pitch2 Status High
Pitch3 Status High
In Position1 ch3
In Position2 ch3
In Position3 ch3
Temp Brake Blade1
Temp Brake Blade2
Temp Brake Blade3
Temp Pitch Motor Blade1
Temp Pitch Motor Blade2
Temp Pitch Motor Blade3
Temp Hub Additional1
Temp Hub Additional2
Temp Hub Additional3
Pitch1 Status Low
Pitch2 Status Low
Pitch3 Status Low
Battery Voltage Blade1 center
Battery Voltage Blade2 center
Battery Voltage Blade3 center
Battery Charging Cur Blade1
Battery Charging Cur Blade2
Battery Charging Cur Blade3
Battery Discharging Cur Blade1
Battery Discharging Cur Blade2
Battery Discharging Cur Blade3
Pitch Motor Brake Voltage Blade1
Pitch Motor Brake Voltage Blade2
Pitch Motor Brake Voltage Blade3
Pitch Motor Brake Current Blade1
Pitch Motor Brake Current Blade2
Pitch Motor Brake Current Blade3
Temp HubBox Blade1
Temp HubBox Blade2
Temp HubBox Blade3
Temp Pitch1 HeatSink
Temp Pitch2 HeatSink
Temp Pitch3 HeatSink
Error Stack Blade1
Error Stack Blade2
Error Stack Blade3
Temp Battery Box Blade1
Temp Battery Box Blade2
Temp Battery Box Blade3
DC Link Voltage1
DC Link Voltage2
DC Link Voltage3
Temp Yaw Motor1
Temp Yaw Motor2
Temp Yaw Motor3
Temp Yaw Motor4
AO DFIG Power Setpiont
AO DFIG Q Setpoint
AI DFIG Torque actual
AI DFIG Speed Generator Encoder
DFIG DC Link Voltage actual
DFIG MSC current
DFIG Main voltage
DFIG Main current
DFIG Active power actual
DFIG reactive power actual
DFIG active power actual LSC
DFIG LSC current
DFIG Data log number
Damper Osc Magnitude
Damper Pass band FullLoad
AI Yaw Brake Temp Rise1
AI Yaw Brake Temp Rise2
AI Yaw Brake Temp Rise3
AI Yaw Brake Temp Rise4
Nacelle Drill at North Pos Sensor`

	oldHeaderList := strings.Split(oldHeader, "\n")
	newHeaderList := strings.Split(newHeader, "\n")
	for idx, val := range oldHeaderList {
		fieldResult = append(fieldResult, strings.ToLower(strings.Replace(strings.TrimSuffix(val, " "), " ", "_", -69)))
		headerResult = append(headerResult, newHeaderList[idx])
	}

	return
}

func GetScadaHFDHeader() (headerResult []string, fieldResult []string) {
	// vArrRealtime := []string{"TimeStamp", "Turbine", "ActivePower_kW", "WindSpeed_ms", "NacellePos", "WindDirection",
	// 	"PitchAngle", "PitchAngle1", "PitchAngle2", "PitchAngle3", "GenSpeed_RPM", "RotorSpeed_RPM",
	// 	"ReactivePower_kVAr", "Frequency_Hz", "TempOutdoor", "TempNacelle", "TempGearBoxHSSNDE", "TempGearBoxHSSDE",
	// 	"TempGearBoxIMSDE", "TempGearBoxOilSump", "TempGearBoxIMSNDE", "TempGeneratorBearingDE", "TempGeneratorBearingNDE",
	// 	"TempHubBearing", "TempG1L1", "TempG1L2", "TempG1L3", "TempBottomControlSection", "TempConv1",
	// 	"TempConv2", "TempConv3", "PitchAccuV1", "PitchAccuV2", "PitchAccuV3", "PowerFactor",
	// 	"Total_Prod_Day_kWh"}

	fieldResult = []string{"Fast_ActivePower_kW", "Fast_WindSpeed_ms", "Slow_NacellePos", "Slow_WindDirection",
		"Fast_PitchAngle", "Fast_PitchAngle1", "Fast_PitchAngle2", "Fast_PitchAngle3", "Fast_GenSpeed_RPM", "Fast_RotorSpeed_RPM",
		"Fast_ReactivePower_kVAr", "Fast_Frequency_Hz", "Slow_TempOutdoor", "Slow_TempNacelle", "Slow_TempGearBoxHSSNDE", "Slow_TempGearBoxHSSDE",
		"Slow_TempGearBoxIMSDE", "Slow_TempGearBoxOilSump", "Slow_TempGearBoxIMSNDE", "Slow_TempGeneratorBearingDE", "Slow_TempGeneratorBearingNDE",
		"Slow_TempHubBearing", "Slow_TempG1L1", "Slow_TempG1L2", "Slow_TempG1L3", "Slow_TempBottomControlSection", "Slow_TempConv1",
		"Slow_TempConv2", "Slow_TempConv3", "Fast_PitchAccuV1", "Fast_PitchAccuV2", "Fast_PitchAccuV3", "Fast_PowerFactor",
		"Fast_Total_Prod_Day_kWh",
		"Fast_ActivePower_kW_Count", "Fast_WindSpeed_ms_Count", "Slow_NacellePos_Count", "Slow_WindDirection_Count",
		"Fast_PitchAngle_Count", "Fast_PitchAngle1_Count", "Fast_PitchAngle2_Count", "Fast_PitchAngle3_Count", "Fast_GenSpeed_RPM_Count", "Fast_RotorSpeed_RPM_Count",
		"Fast_ReactivePower_kVAr_Count", "Fast_Frequency_Hz_Count", "Slow_TempOutdoor_Count", "Slow_TempNacelle_Count", "Slow_TempGearBoxHSSNDE_Count", "Slow_TempGearBoxHSSDE_Count",
		"Slow_TempGearBoxIMSDE_Count", "Slow_TempGearBoxOilSump_Count", "Slow_TempGearBoxIMSNDE_Count", "Slow_TempGeneratorBearingDE_Count", "Slow_TempGeneratorBearingNDE_Count",
		"Slow_TempHubBearing_Count", "Slow_TempG1L1_Count", "Slow_TempG1L2_Count", "Slow_TempG1L3_Count", "Slow_TempBottomControlSection_Count", "Slow_TempConv1_Count",
		"Slow_TempConv2_Count", "Slow_TempConv3_Count", "Fast_PitchAccuV1_Count", "Fast_PitchAccuV2_Count", "Fast_PitchAccuV3_Count", "Fast_PowerFactor_Count",
		"Fast_Total_Prod_Day_kWh_Count"}

	headerResult = []string{"Active Power", "Wind Speed", "Nacelle Pos", "Wind Direction",
		"Pitch Angle", "Pitch Angle1", "Pitch Angle2", "Pitch Angle3", "Generator Speed", "Rotor Speed",
		"Reactive Power", "Frequency Grid", "Ambient Temp", "Temp Nacelle", "Temp GearBox HSS NDE", "Temp GearBox HSS DE",
		"Temp GearBox IMS DE", "Temp GearBox Oil Sump", "Temp GearBox IMS NDE", "Temp Generator Bearing DE", "Temp Generator Bearing NDE",
		"Temp Main Bearing", "Temp G1L1", "Temp G1L2", "Temp G1L3", "Temp Bottom Control Section", "Temp Conv1",
		"Temp Conv2", "Temp Conv3", "Pitch Accu V1", "Pitch Accu V2", "Pitch Accu V3", "Power Factor",
		"Total Production Day",
		"Active Power Count", "Wind Speed Count", "Nacelle Pos Count", "Wind Direction Count",
		"Pitch Angle Count", "Pitch Angle1 Count", "Pitch Angle2 Count", "Pitch Angle3 Count", "Generator Speed Count", "Rotor Speed Count",
		"Reactive Power Count", "Frequency Grid Count", "Ambient Temp Count", "Temp Nacelle Count", "Temp GearBox HSS NDE Count", "Temp GearBox HSS DE Count",
		"Temp GearBox IMS DE Count", "Temp GearBox Oil Sump Count", "Temp GearBox IMS NDE Count", "Temp Generator Bearing DE Count", "Temp Generator Bearing NDE Count",
		"Temp Main Bearing Count", "Temp G1L1 Count", "Temp G1L2 Count", "Temp G1L3 Count", "Temp Bottom Control Section Count", "Temp Conv1 Count",
		"Temp Conv2 Count", "Temp Conv3 Count", "Pitch Accu V1 Count", "Pitch Accu V2 Count", "Pitch Accu V3 Count", "Power Factor Count",
		"Total Production Day Count"}

	return
}
