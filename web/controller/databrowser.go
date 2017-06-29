package controller

import (
	. "eaciit/wfdemo-git/library/core"
	. "eaciit/wfdemo-git/library/models"
	"eaciit/wfdemo-git/web/helper"
	"strings"

	// "fmt"
	"os"
	"strconv"
	"time"

	"github.com/eaciit/dbox"
	"github.com/eaciit/knot/knot.v1"
	tk "github.com/eaciit/toolkit"
	x "github.com/tealeg/xlsx"
	// f "path/filepath"
)

type DataBrowserController struct {
	App
}

func CreateDataBrowserController() *DataBrowserController {
	var controller = new(DataBrowserController)
	return controller
}

func GetCustomFieldList() []tk.M {
	atkm := []tk.M{}

	_ascadaoem_label := []string{"Ai Intern R Pid Angle Out", "Ai Intern I1", "Ai Intern I2",
		"Ai Dfig Torque Actual", "Ai Dr Tr Vib Value", "Ai Gear Oil Pressure", "Ai Hydr System Pressure", "Ai Intern Active Power",
		"Ai Intern Dfig Active Power Actual", "Ai Intern Nacelle Drill", "Ai Intern Nacelle Drill At North Pos Sensor", "Ai Intern Nacelle Pos", "Ai Intern Pitch Angle1",
		"Ai Intern Pitch Angle2", "Ai Intern Pitch Angle3", "Ai Intern Pitch Speed1", "Ai Intern Reactive Power", "Ai Intern Wind Direction",
		"Ai Intern Wind Speed", "Ai Intern Wind Speed Dif", "Ai Tower Vib Value Axial", "Ai Wind Speed1", "Ai Wind Speed2",
		"Ai Wind Vane1", "Ai Wind Vane2", "C Intern Speed Generator", "C Intern Speed Rotor", "Temp Bottom Control Section",
		"Temp Bottom Control Section Low", "Temp Bottom Power Section", "Temp Cabinet Top Box", "Temp Gearbox Hss De", "Temp Gear Box Hss Nde",
		"Temp Gear Box Ims De", "Temp Gear Box Ims Nde", "Temp Gear Oil Sump", "Temp Generator Bearing De", "Temp Generator Bearing Nde",
		"Temp Main Bearing", "Temp Nacelle", "Temp Outdoor", "Time Stamp", "Turbine",
	}

	_ascadaoem_field := []string{"ai_intern_r_pidangleout", "ai_intern_i1", "ai_intern_i2",
		"ai_dfig_torque_actual", "ai_drtrvibvalue", "ai_gearoilpressure", "ai_hydrsystempressure", "ai_intern_activpower",
		"ai_intern_dfig_active_power_actual", "ai_intern_nacelledrill", "ai_intern_nacelledrill_at_northpossensor", "ai_intern_nacellepos", "ai_intern_pitchangle1",
		"ai_intern_pitchangle2", "ai_intern_pitchangle3", "ai_intern_pitchspeed1", "ai_intern_reactivpower", "ai_intern_winddirection",
		"ai_intern_windspeed", "ai_intern_windspeeddif", "ai_towervibvalueaxial", "ai_windspeed1", "ai_windspeed2",
		"ai_windvane1", "ai_windvane2", "c_intern_speedgenerator", "c_intern_speedrotor", "temp_bottomcontrolsection",
		"temp_bottomcontrolsection_low", "temp_bottompowersection", "temp_cabinettopbox", "temp_gearbox_hss_de", "temp_gearbox_hss_nde",
		"temp_gearbox_ims_de", "temp_gearbox_ims_nde", "temp_gearoilsump", "temp_generatorbearing_de", "temp_generatorbearing_nde",
		"temp_mainbearing", "temp_nacelle", "temp_outdoor", "timestamp", "turbine",
	}

	_amettower_label := []string{"V Hub WS 90m Avg", "V Hub WS 90m Std Dev", "V Ref WS 88m Avg", "V Ref WS 88m Std Dev",
		"V Tip WS 42m Avg", "V Tip WS 42m Std Dev", "D Hub WD 88m Avg", "D Hub WD 88m Std Dev", "D Ref WD 86m Avg",
		"D Ref WD 86m Std Dev", "T Hub & H Hub Humid 85m Avg", "T Hub & H Hub Humid 85m Std Dev", "T Ref & H Ref Humid 85.5m Avg", "T Ref & H Ref Humid 85.5m Std Dev",
		"T Hub & H Hub Temp 85.5m Avg", "T Hub & H Hub Temp 85.5m Std Dev", "T Ref & H Ref Temp 85.5 Avg", "T Ref & H Ref Temp 85.5 Std Dev", "Baro Air Pressure 85.5m Avg", "Baro Air Pressure 85.5m Std Dev",
	}

	_amettower_field := []string{"vhubws90mavg", "vhubws90mstddev", "vrefws88mavg", "vrefws88mstddev", "vtipws42mavg",
		"vtipws42mstddev", "dhubwd88mavg", "dhubwd88mstddev", "drefwd86mavg", "drefwd86mstddev",
		"thubhhubhumid855mavg", "thubhhubhumid855mstddev", "trefhrefhumid855mavg", "trefhrefhumid855mstddev", "thubhhubtemp855mavg",
		"thubhhubtemp855mstddev", "trefhreftemp855mavg", "trefhreftemp855mstddev", "baroairpress855mavg", "baroairpress855mstddev",
	}

	for i, str := range _ascadaoem_field {
		tkm := tk.M{}.
			Set("_id", str).
			Set("label", _ascadaoem_label[i]).
			Set("source", "ScadaDataOEM")

		atkm = append(atkm, tkm)
	}

	for i, str := range _amettower_field {
		tkm := tk.M{}.
			Set("_id", str).
			Set("label", _amettower_label[i]).
			Set("source", "MetTower")

		atkm = append(atkm, tkm)
	}

	return atkm
}

func GetHFDCustomFieldList() []tk.M {
	atkm := []tk.M{}

	_ashfd_label := []string{"Fast ActivePower kW", "Fast ActivePower kW StdDev", "Fast ActivePower kW Min", "Fast ActivePower kW Max", "Fast ActivePower kW Count", "Fast WindSpeed ms", "Fast WindSpeed ms StdDev", "Fast WindSpeed ms Min", "Fast WindSpeed ms Max", "Fast WindSpeed ms Count", "Slow NacellePos", "Slow NacellePos StdDev", "Slow NacellePos Min", "Slow NacellePos Max", "Slow NacellePos Count", "Slow WindDirection", "Slow WindDirection StdDev", "Slow WindDirection Min", "Slow WindDirection Max", "Slow WindDirection Count", "Fast CurrentL3", "Fast CurrentL3 StdDev", "Fast CurrentL3 Min", "Fast CurrentL3 Max", "Fast CurrentL3 Count", "Fast CurrentL1", "Fast CurrentL1 StdDev", "Fast CurrentL1 Min", "Fast CurrentL1 Max", "Fast CurrentL1 Count", "Fast ActivePowerSetpoint kW", "Fast ActivePowerSetpoint kW StdDev", "Fast ActivePowerSetpoint kW Min", "Fast ActivePowerSetpoint kW Max", "Fast ActivePowerSetpoint kW Count", "Fast CurrentL2", "Fast CurrentL2 StdDev", "Fast CurrentL2 Min", "Fast CurrentL2 Max", "Fast CurrentL2 Count", "Fast DrTrVibValue", "Fast DrTrVibValue StdDev", "Fast DrTrVibValue Min", "Fast DrTrVibValue Max", "Fast DrTrVibValue Count", "Fast GenSpeed RPM", "Fast GenSpeed RPM StdDev", "Fast GenSpeed RPM Min", "Fast GenSpeed RPM Max", "Fast GenSpeed RPM Count", "Fast PitchAccuV1", "Fast PitchAccuV1 StdDev", "Fast PitchAccuV1 Min", "Fast PitchAccuV1 Max", "Fast PitchAccuV1 Count", "Fast PitchAngle", "Fast PitchAngle StdDev", "Fast PitchAngle Min", "Fast PitchAngle Max", "Fast PitchAngle Count", "Fast PitchAngle3", "Fast PitchAngle3 StdDev", "Fast PitchAngle3 Min", "Fast PitchAngle3 Max", "Fast PitchAngle3 Count", "Fast PitchAngle2", "Fast PitchAngle2 StdDev", "Fast PitchAngle2 Min", "Fast PitchAngle2 Max", "Fast PitchAngle2 Count", "Fast PitchConvCurrent1", "Fast PitchConvCurrent1 StdDev", "Fast PitchConvCurrent1 Min", "Fast PitchConvCurrent1 Max", "Fast PitchConvCurrent1 Count", "Fast PitchConvCurrent3", "Fast PitchConvCurrent3 StdDev", "Fast PitchConvCurrent3 Min", "Fast PitchConvCurrent3 Max", "Fast PitchConvCurrent3 Count", "Fast PitchConvCurrent2", "Fast PitchConvCurrent2 StdDev", "Fast PitchConvCurrent2 Min", "Fast PitchConvCurrent2 Max", "Fast PitchConvCurrent2 Count", "Fast PowerFactor", "Fast PowerFactor StdDev", "Fast PowerFactor Min", "Fast PowerFactor Max", "Fast PowerFactor Count", "Fast ReactivePowerSetpointPPC kVA", "Fast ReactivePowerSetpointPPC kVA StdDev", "Fast ReactivePowerSetpointPPC kVA Min", "Fast ReactivePowerSetpointPPC kVA Max", "Fast ReactivePowerSetpointPPC kVA Count", "Fast ReactivePower kVAr", "Fast ReactivePower kVAr StdDev", "Fast ReactivePower kVAr Min", "Fast ReactivePower kVAr Max", "Fast ReactivePower kVAr Count", "Fast RotorSpeed RPM", "Fast RotorSpeed RPM StdDev", "Fast RotorSpeed RPM Min", "Fast RotorSpeed RPM Max", "Fast RotorSpeed RPM Count", "Fast VoltageL1", "Fast VoltageL1 StdDev", "Fast VoltageL1 Min", "Fast VoltageL1 Max", "Fast VoltageL1 Count", "Fast VoltageL2", "Fast VoltageL2 StdDev", "Fast VoltageL2 Min", "Fast VoltageL2 Max", "Fast VoltageL2 Count", "Slow CapableCapacitiveReactPwr kVAr", "Slow CapableCapacitiveReactPwr kVAr StdDev", "Slow CapableCapacitiveReactPwr kVAr Min", "Slow CapableCapacitiveReactPwr kVAr Max", "Slow CapableCapacitiveReactPwr kVAr Count", "Slow CapableInductiveReactPwr kVAr", "Slow CapableInductiveReactPwr kVAr StdDev", "Slow CapableInductiveReactPwr kVAr Min", "Slow CapableInductiveReactPwr kVAr Max", "Slow CapableInductiveReactPwr kVAr Count", "Slow DateTime Sec", "Slow DateTime Sec StdDev", "Slow DateTime Sec Min", "Slow DateTime Sec Max", "Slow DateTime Sec Count", "Fast PitchAngle1", "Fast PitchAngle1 StdDev", "Fast PitchAngle1 Min", "Fast PitchAngle1 Max", "Fast PitchAngle1 Count", "Fast VoltageL3", "Fast VoltageL3 StdDev", "Fast VoltageL3 Min", "Fast VoltageL3 Max", "Fast VoltageL3 Count", "Slow CapableCapacitivePwrFactor", "Slow CapableCapacitivePwrFactor StdDev", "Slow CapableCapacitivePwrFactor Min", "Slow CapableCapacitivePwrFactor Max", "Slow CapableCapacitivePwrFactor Count", "Fast Total Production kWh", "Fast Total Production kWh StdDev", "Fast Total Production kWh Min", "Fast Total Production kWh Max", "Fast Total Production kWh Count", "Fast Total Prod Day kWh", "Fast Total Prod Day kWh StdDev", "Fast Total Prod Day kWh Min", "Fast Total Prod Day kWh Max", "Fast Total Prod Day kWh Count", "Fast Total Prod Month kWh", "Fast Total Prod Month kWh StdDev", "Fast Total Prod Month kWh Min", "Fast Total Prod Month kWh Max", "Fast Total Prod Month kWh Count", "Fast ActivePowerOutPWCSell kW", "Fast ActivePowerOutPWCSell kW StdDev", "Fast ActivePowerOutPWCSell kW Min", "Fast ActivePowerOutPWCSell kW Max", "Fast ActivePowerOutPWCSell kW Count", "Fast Frequency Hz", "Fast Frequency Hz StdDev", "Fast Frequency Hz Min", "Fast Frequency Hz Max", "Fast Frequency Hz Count", "Slow TempG1L2", "Slow TempG1L2 StdDev", "Slow TempG1L2 Min", "Slow TempG1L2 Max", "Slow TempG1L2 Count", "Slow TempG1L3", "Slow TempG1L3 StdDev", "Slow TempG1L3 Min", "Slow TempG1L3 Max", "Slow TempG1L3 Count", "Slow TempGearBoxHSSDE", "Slow TempGearBoxHSSDE StdDev", "Slow TempGearBoxHSSDE Min", "Slow TempGearBoxHSSDE Max", "Slow TempGearBoxHSSDE Count", "Slow TempGearBoxIMSNDE", "Slow TempGearBoxIMSNDE StdDev", "Slow TempGearBoxIMSNDE Min", "Slow TempGearBoxIMSNDE Max", "Slow TempGearBoxIMSNDE Count", "Slow TempOutdoor", "Slow TempOutdoor StdDev", "Slow TempOutdoor Min", "Slow TempOutdoor Max", "Slow TempOutdoor Count", "Fast PitchAccuV3", "Fast PitchAccuV3 StdDev", "Fast PitchAccuV3 Min", "Fast PitchAccuV3 Max", "Fast PitchAccuV3 Count", "Slow TotalTurbineActiveHours", "Slow TotalTurbineActiveHours StdDev", "Slow TotalTurbineActiveHours Min", "Slow TotalTurbineActiveHours Max", "Slow TotalTurbineActiveHours Count", "Slow TotalTurbineOKHours", "Slow TotalTurbineOKHours StdDev", "Slow TotalTurbineOKHours Min", "Slow TotalTurbineOKHours Max", "Slow TotalTurbineOKHours Count", "Slow TotalTurbineTimeAllHours", "Slow TotalTurbineTimeAllHours StdDev", "Slow TotalTurbineTimeAllHours Min", "Slow TotalTurbineTimeAllHours Max", "Slow TotalTurbineTimeAllHours Count", "Slow TempG1L1", "Slow TempG1L1 StdDev", "Slow TempG1L1 Min", "Slow TempG1L1 Max", "Slow TempG1L1 Count", "Slow TempGearBoxOilSump", "Slow TempGearBoxOilSump StdDev", "Slow TempGearBoxOilSump Min", "Slow TempGearBoxOilSump Max", "Slow TempGearBoxOilSump Count", "Fast PitchAccuV2", "Fast PitchAccuV2 StdDev", "Fast PitchAccuV2 Min", "Fast PitchAccuV2 Max", "Fast PitchAccuV2 Count", "Slow TotalGridOkHours", "Slow TotalGridOkHours StdDev", "Slow TotalGridOkHours Min", "Slow TotalGridOkHours Max", "Slow TotalGridOkHours Count", "Slow TotalActPowerOut kWh", "Slow TotalActPowerOut kWh StdDev", "Slow TotalActPowerOut kWh Min", "Slow TotalActPowerOut kWh Max", "Slow TotalActPowerOut kWh Count", "Fast YawService", "Fast YawService StdDev", "Fast YawService Min", "Fast YawService Max", "Fast YawService Count", "Fast YawAngle", "Fast YawAngle StdDev", "Fast YawAngle Min", "Fast YawAngle Max", "Fast YawAngle Count", "Slow CapableInductivePwrFactor", "Slow CapableInductivePwrFactor StdDev", "Slow CapableInductivePwrFactor Min", "Slow CapableInductivePwrFactor Max", "Slow CapableInductivePwrFactor Count", "Slow TempGearBoxHSSNDE", "Slow TempGearBoxHSSNDE StdDev", "Slow TempGearBoxHSSNDE Min", "Slow TempGearBoxHSSNDE Max", "Slow TempGearBoxHSSNDE Count", "Slow TempHubBearing", "Slow TempHubBearing StdDev", "Slow TempHubBearing Min", "Slow TempHubBearing Max", "Slow TempHubBearing Count", "Slow TotalG1ActiveHours", "Slow TotalG1ActiveHours StdDev", "Slow TotalG1ActiveHours Min", "Slow TotalG1ActiveHours Max", "Slow TotalG1ActiveHours Count", "Slow TotalActPowerOutG1 kWh", "Slow TotalActPowerOutG1 kWh StdDev", "Slow TotalActPowerOutG1 kWh Min", "Slow TotalActPowerOutG1 kWh Max", "Slow TotalActPowerOutG1 kWh Count", "Slow TotalReactPowerInG1 kVArh", "Slow TotalReactPowerInG1 kVArh StdDev", "Slow TotalReactPowerInG1 kVArh Min", "Slow TotalReactPowerInG1 kVArh Max", "Slow TotalReactPowerInG1 kVArh Count", "Slow NacelleDrill", "Slow NacelleDrill StdDev", "Slow NacelleDrill Min", "Slow NacelleDrill Max", "Slow NacelleDrill Count", "Slow TempGearBoxIMSDE", "Slow TempGearBoxIMSDE StdDev", "Slow TempGearBoxIMSDE Min", "Slow TempGearBoxIMSDE Max", "Slow TempGearBoxIMSDE Count", "Fast Total Operating hrs", "Fast Total Operating hrs StdDev", "Fast Total Operating hrs Min", "Fast Total Operating hrs Max", "Fast Total Operating hrs Count", "Slow TempNacelle", "Slow TempNacelle StdDev", "Slow TempNacelle Min", "Slow TempNacelle Max", "Slow TempNacelle Count", "Fast Total Grid OK hrs", "Fast Total Grid OK hrs StdDev", "Fast Total Grid OK hrs Min", "Fast Total Grid OK hrs Max", "Fast Total Grid OK hrs Count", "Fast Total WTG OK hrs", "Fast Total WTG OK hrs StdDev", "Fast Total WTG OK hrs Min", "Fast Total WTG OK hrs Max", "Fast Total WTG OK hrs Count", "Slow TempCabinetTopBox", "Slow TempCabinetTopBox StdDev", "Slow TempCabinetTopBox Min", "Slow TempCabinetTopBox Max", "Slow TempCabinetTopBox Count", "Slow TempGeneratorBearingNDE", "Slow TempGeneratorBearingNDE StdDev", "Slow TempGeneratorBearingNDE Min", "Slow TempGeneratorBearingNDE Max", "Slow TempGeneratorBearingNDE Count", "Fast Total Access hrs", "Fast Total Access hrs StdDev", "Fast Total Access hrs Min", "Fast Total Access hrs Max", "Fast Total Access hrs Count", "Slow TempBottomPowerSection", "Slow TempBottomPowerSection StdDev", "Slow TempBottomPowerSection Min", "Slow TempBottomPowerSection Max", "Slow TempBottomPowerSection Count", "Slow TempGeneratorBearingDE", "Slow TempGeneratorBearingDE StdDev", "Slow TempGeneratorBearingDE Min", "Slow TempGeneratorBearingDE Max", "Slow TempGeneratorBearingDE Count", "Slow TotalReactPowerIn kVArh", "Slow TotalReactPowerIn kVArh StdDev", "Slow TotalReactPowerIn kVArh Min", "Slow TotalReactPowerIn kVArh Max", "Slow TotalReactPowerIn kVArh Count", "Slow TempBottomControlSection", "Slow TempBottomControlSection StdDev", "Slow TempBottomControlSection Min", "Slow TempBottomControlSection Max", "Slow TempBottomControlSection Count", "Slow TempConv1", "Slow TempConv1 StdDev", "Slow TempConv1 Min", "Slow TempConv1 Max", "Slow TempConv1 Count", "Fast ActivePowerRated kW", "Fast ActivePowerRated kW StdDev", "Fast ActivePowerRated kW Min", "Fast ActivePowerRated kW Max", "Fast ActivePowerRated kW Count", "Fast NodeIP", "Fast NodeIP StdDev", "Fast NodeIP Min", "Fast NodeIP Max", "Fast NodeIP Count", "Fast PitchSpeed1", "Fast PitchSpeed1 StdDev", "Fast PitchSpeed1 Min", "Fast PitchSpeed1 Max", "Fast PitchSpeed1 Count", "Slow CFCardSize", "Slow CFCardSize StdDev", "Slow CFCardSize Min", "Slow CFCardSize Max", "Slow CFCardSize Count", "Slow CPU Number", "Slow CPU Number StdDev", "Slow CPU Number Min", "Slow CPU Number Max", "Slow CPU Number Count", "Slow CFCardSpaceLeft", "Slow CFCardSpaceLeft StdDev", "Slow CFCardSpaceLeft Min", "Slow CFCardSpaceLeft Max", "Slow CFCardSpaceLeft Count", "Slow TempBottomCapSection", "Slow TempBottomCapSection StdDev", "Slow TempBottomCapSection Min", "Slow TempBottomCapSection Max", "Slow TempBottomCapSection Count", "Slow RatedPower", "Slow RatedPower StdDev", "Slow RatedPower Min", "Slow RatedPower Max", "Slow RatedPower Count", "Slow TempConv3", "Slow TempConv3 StdDev", "Slow TempConv3 Min", "Slow TempConv3 Max", "Slow TempConv3 Count", "Slow TempConv2", "Slow TempConv2 StdDev", "Slow TempConv2 Min", "Slow TempConv2 Max", "Slow TempConv2 Count", "Slow TotalActPowerIn kWh", "Slow TotalActPowerIn kWh StdDev", "Slow TotalActPowerIn kWh Min", "Slow TotalActPowerIn kWh Max", "Slow TotalActPowerIn kWh Count", "Slow TotalActPowerInG1 kWh", "Slow TotalActPowerInG1 kWh StdDev", "Slow TotalActPowerInG1 kWh Min", "Slow TotalActPowerInG1 kWh Max", "Slow TotalActPowerInG1 kWh Count", "Slow TotalActPowerInG2 kWh", "Slow TotalActPowerInG2 kWh StdDev", "Slow TotalActPowerInG2 kWh Min", "Slow TotalActPowerInG2 kWh Max", "Slow TotalActPowerInG2 kWh Count", "Slow TotalActPowerOutG2 kWh", "Slow TotalActPowerOutG2 kWh StdDev", "Slow TotalActPowerOutG2 kWh Min", "Slow TotalActPowerOutG2 kWh Max", "Slow TotalActPowerOutG2 kWh Count", "Slow TotalG2ActiveHours", "Slow TotalG2ActiveHours StdDev", "Slow TotalG2ActiveHours Min", "Slow TotalG2ActiveHours Max", "Slow TotalG2ActiveHours Count", "Slow TotalReactPowerInG2 kVArh", "Slow TotalReactPowerInG2 kVArh StdDev", "Slow TotalReactPowerInG2 kVArh Min", "Slow TotalReactPowerInG2 kVArh Max", "Slow TotalReactPowerInG2 kVArh Count", "Slow TotalReactPowerOut kVArh", "Slow TotalReactPowerOut kVArh StdDev", "Slow TotalReactPowerOut kVArh Min", "Slow TotalReactPowerOut kVArh Max", "Slow TotalReactPowerOut kVArh Count", "Slow UTCoffset int", "Slow UTCoffset int StdDev", "Slow UTCoffset int Min", "Slow UTCoffset int Max", "Slow UTCoffset int Count"}

	_ashfd_field := []string{"Fast_ActivePower_kW", "Fast_ActivePower_kW_StdDev", "Fast_ActivePower_kW_Min", "Fast_ActivePower_kW_Max", "Fast_ActivePower_kW_Count", "Fast_WindSpeed_ms", "Fast_WindSpeed_ms_StdDev", "Fast_WindSpeed_ms_Min", "Fast_WindSpeed_ms_Max", "Fast_WindSpeed_ms_Count", "Slow_NacellePos", "Slow_NacellePos_StdDev", "Slow_NacellePos_Min", "Slow_NacellePos_Max", "Slow_NacellePos_Count", "Slow_WindDirection", "Slow_WindDirection_StdDev", "Slow_WindDirection_Min", "Slow_WindDirection_Max", "Slow_WindDirection_Count", "Fast_CurrentL3", "Fast_CurrentL3_StdDev", "Fast_CurrentL3_Min", "Fast_CurrentL3_Max", "Fast_CurrentL3_Count", "Fast_CurrentL1", "Fast_CurrentL1_StdDev", "Fast_CurrentL1_Min", "Fast_CurrentL1_Max", "Fast_CurrentL1_Count", "Fast_ActivePowerSetpoint_kW", "Fast_ActivePowerSetpoint_kW_StdDev", "Fast_ActivePowerSetpoint_kW_Min", "Fast_ActivePowerSetpoint_kW_Max", "Fast_ActivePowerSetpoint_kW_Count", "Fast_CurrentL2", "Fast_CurrentL2_StdDev", "Fast_CurrentL2_Min", "Fast_CurrentL2_Max", "Fast_CurrentL2_Count", "Fast_DrTrVibValue", "Fast_DrTrVibValue_StdDev", "Fast_DrTrVibValue_Min", "Fast_DrTrVibValue_Max", "Fast_DrTrVibValue_Count", "Fast_GenSpeed_RPM", "Fast_GenSpeed_RPM_StdDev", "Fast_GenSpeed_RPM_Min", "Fast_GenSpeed_RPM_Max", "Fast_GenSpeed_RPM_Count", "Fast_PitchAccuV1", "Fast_PitchAccuV1_StdDev", "Fast_PitchAccuV1_Min", "Fast_PitchAccuV1_Max", "Fast_PitchAccuV1_Count", "Fast_PitchAngle", "Fast_PitchAngle_StdDev", "Fast_PitchAngle_Min", "Fast_PitchAngle_Max", "Fast_PitchAngle_Count", "Fast_PitchAngle3", "Fast_PitchAngle3_StdDev", "Fast_PitchAngle3_Min", "Fast_PitchAngle3_Max", "Fast_PitchAngle3_Count", "Fast_PitchAngle2", "Fast_PitchAngle2_StdDev", "Fast_PitchAngle2_Min", "Fast_PitchAngle2_Max", "Fast_PitchAngle2_Count", "Fast_PitchConvCurrent1", "Fast_PitchConvCurrent1_StdDev", "Fast_PitchConvCurrent1_Min", "Fast_PitchConvCurrent1_Max", "Fast_PitchConvCurrent1_Count", "Fast_PitchConvCurrent3", "Fast_PitchConvCurrent3_StdDev", "Fast_PitchConvCurrent3_Min", "Fast_PitchConvCurrent3_Max", "Fast_PitchConvCurrent3_Count", "Fast_PitchConvCurrent2", "Fast_PitchConvCurrent2_StdDev", "Fast_PitchConvCurrent2_Min", "Fast_PitchConvCurrent2_Max", "Fast_PitchConvCurrent2_Count", "Fast_PowerFactor", "Fast_PowerFactor_StdDev", "Fast_PowerFactor_Min", "Fast_PowerFactor_Max", "Fast_PowerFactor_Count", "Fast_ReactivePowerSetpointPPC_kVA", "Fast_ReactivePowerSetpointPPC_kVA_StdDev", "Fast_ReactivePowerSetpointPPC_kVA_Min", "Fast_ReactivePowerSetpointPPC_kVA_Max", "Fast_ReactivePowerSetpointPPC_kVA_Count", "Fast_ReactivePower_kVAr", "Fast_ReactivePower_kVAr_StdDev", "Fast_ReactivePower_kVAr_Min", "Fast_ReactivePower_kVAr_Max", "Fast_ReactivePower_kVAr_Count", "Fast_RotorSpeed_RPM", "Fast_RotorSpeed_RPM_StdDev", "Fast_RotorSpeed_RPM_Min", "Fast_RotorSpeed_RPM_Max", "Fast_RotorSpeed_RPM_Count", "Fast_VoltageL1", "Fast_VoltageL1_StdDev", "Fast_VoltageL1_Min", "Fast_VoltageL1_Max", "Fast_VoltageL1_Count", "Fast_VoltageL2", "Fast_VoltageL2_StdDev", "Fast_VoltageL2_Min", "Fast_VoltageL2_Max", "Fast_VoltageL2_Count", "Slow_CapableCapacitiveReactPwr_kVAr", "Slow_CapableCapacitiveReactPwr_kVAr_StdDev", "Slow_CapableCapacitiveReactPwr_kVAr_Min", "Slow_CapableCapacitiveReactPwr_kVAr_Max", "Slow_CapableCapacitiveReactPwr_kVAr_Count", "Slow_CapableInductiveReactPwr_kVAr", "Slow_CapableInductiveReactPwr_kVAr_StdDev", "Slow_CapableInductiveReactPwr_kVAr_Min", "Slow_CapableInductiveReactPwr_kVAr_Max", "Slow_CapableInductiveReactPwr_kVAr_Count", "Slow_DateTime_Sec", "Slow_DateTime_Sec_StdDev", "Slow_DateTime_Sec_Min", "Slow_DateTime_Sec_Max", "Slow_DateTime_Sec_Count", "Fast_PitchAngle1", "Fast_PitchAngle1_StdDev", "Fast_PitchAngle1_Min", "Fast_PitchAngle1_Max", "Fast_PitchAngle1_Count", "Fast_VoltageL3", "Fast_VoltageL3_StdDev", "Fast_VoltageL3_Min", "Fast_VoltageL3_Max", "Fast_VoltageL3_Count", "Slow_CapableCapacitivePwrFactor", "Slow_CapableCapacitivePwrFactor_StdDev", "Slow_CapableCapacitivePwrFactor_Min", "Slow_CapableCapacitivePwrFactor_Max", "Slow_CapableCapacitivePwrFactor_Count", "Fast_Total_Production_kWh", "Fast_Total_Production_kWh_StdDev", "Fast_Total_Production_kWh_Min", "Fast_Total_Production_kWh_Max", "Fast_Total_Production_kWh_Count", "Fast_Total_Prod_Day_kWh", "Fast_Total_Prod_Day_kWh_StdDev", "Fast_Total_Prod_Day_kWh_Min", "Fast_Total_Prod_Day_kWh_Max", "Fast_Total_Prod_Day_kWh_Count", "Fast_Total_Prod_Month_kWh", "Fast_Total_Prod_Month_kWh_StdDev", "Fast_Total_Prod_Month_kWh_Min", "Fast_Total_Prod_Month_kWh_Max", "Fast_Total_Prod_Month_kWh_Count", "Fast_ActivePowerOutPWCSell_kW", "Fast_ActivePowerOutPWCSell_kW_StdDev", "Fast_ActivePowerOutPWCSell_kW_Min", "Fast_ActivePowerOutPWCSell_kW_Max", "Fast_ActivePowerOutPWCSell_kW_Count", "Fast_Frequency_Hz", "Fast_Frequency_Hz_StdDev", "Fast_Frequency_Hz_Min", "Fast_Frequency_Hz_Max", "Fast_Frequency_Hz_Count", "Slow_TempG1L2", "Slow_TempG1L2_StdDev", "Slow_TempG1L2_Min", "Slow_TempG1L2_Max", "Slow_TempG1L2_Count", "Slow_TempG1L3", "Slow_TempG1L3_StdDev", "Slow_TempG1L3_Min", "Slow_TempG1L3_Max", "Slow_TempG1L3_Count", "Slow_TempGearBoxHSSDE", "Slow_TempGearBoxHSSDE_StdDev", "Slow_TempGearBoxHSSDE_Min", "Slow_TempGearBoxHSSDE_Max", "Slow_TempGearBoxHSSDE_Count", "Slow_TempGearBoxIMSNDE", "Slow_TempGearBoxIMSNDE_StdDev", "Slow_TempGearBoxIMSNDE_Min", "Slow_TempGearBoxIMSNDE_Max", "Slow_TempGearBoxIMSNDE_Count", "Slow_TempOutdoor", "Slow_TempOutdoor_StdDev", "Slow_TempOutdoor_Min", "Slow_TempOutdoor_Max", "Slow_TempOutdoor_Count", "Fast_PitchAccuV3", "Fast_PitchAccuV3_StdDev", "Fast_PitchAccuV3_Min", "Fast_PitchAccuV3_Max", "Fast_PitchAccuV3_Count", "Slow_TotalTurbineActiveHours", "Slow_TotalTurbineActiveHours_StdDev", "Slow_TotalTurbineActiveHours_Min", "Slow_TotalTurbineActiveHours_Max", "Slow_TotalTurbineActiveHours_Count", "Slow_TotalTurbineOKHours", "Slow_TotalTurbineOKHours_StdDev", "Slow_TotalTurbineOKHours_Min", "Slow_TotalTurbineOKHours_Max", "Slow_TotalTurbineOKHours_Count", "Slow_TotalTurbineTimeAllHours", "Slow_TotalTurbineTimeAllHours_StdDev", "Slow_TotalTurbineTimeAllHours_Min", "Slow_TotalTurbineTimeAllHours_Max", "Slow_TotalTurbineTimeAllHours_Count", "Slow_TempG1L1", "Slow_TempG1L1_StdDev", "Slow_TempG1L1_Min", "Slow_TempG1L1_Max", "Slow_TempG1L1_Count", "Slow_TempGearBoxOilSump", "Slow_TempGearBoxOilSump_StdDev", "Slow_TempGearBoxOilSump_Min", "Slow_TempGearBoxOilSump_Max", "Slow_TempGearBoxOilSump_Count", "Fast_PitchAccuV2", "Fast_PitchAccuV2_StdDev", "Fast_PitchAccuV2_Min", "Fast_PitchAccuV2_Max", "Fast_PitchAccuV2_Count", "Slow_TotalGridOkHours", "Slow_TotalGridOkHours_StdDev", "Slow_TotalGridOkHours_Min", "Slow_TotalGridOkHours_Max", "Slow_TotalGridOkHours_Count", "Slow_TotalActPowerOut_kWh", "Slow_TotalActPowerOut_kWh_StdDev", "Slow_TotalActPowerOut_kWh_Min", "Slow_TotalActPowerOut_kWh_Max", "Slow_TotalActPowerOut_kWh_Count", "Fast_YawService", "Fast_YawService_StdDev", "Fast_YawService_Min", "Fast_YawService_Max", "Fast_YawService_Count", "Fast_YawAngle", "Fast_YawAngle_StdDev", "Fast_YawAngle_Min", "Fast_YawAngle_Max", "Fast_YawAngle_Count", "Slow_CapableInductivePwrFactor", "Slow_CapableInductivePwrFactor_StdDev", "Slow_CapableInductivePwrFactor_Min", "Slow_CapableInductivePwrFactor_Max", "Slow_CapableInductivePwrFactor_Count", "Slow_TempGearBoxHSSNDE", "Slow_TempGearBoxHSSNDE_StdDev", "Slow_TempGearBoxHSSNDE_Min", "Slow_TempGearBoxHSSNDE_Max", "Slow_TempGearBoxHSSNDE_Count", "Slow_TempHubBearing", "Slow_TempHubBearing_StdDev", "Slow_TempHubBearing_Min", "Slow_TempHubBearing_Max", "Slow_TempHubBearing_Count", "Slow_TotalG1ActiveHours", "Slow_TotalG1ActiveHours_StdDev", "Slow_TotalG1ActiveHours_Min", "Slow_TotalG1ActiveHours_Max", "Slow_TotalG1ActiveHours_Count", "Slow_TotalActPowerOutG1_kWh", "Slow_TotalActPowerOutG1_kWh_StdDev", "Slow_TotalActPowerOutG1_kWh_Min", "Slow_TotalActPowerOutG1_kWh_Max", "Slow_TotalActPowerOutG1_kWh_Count", "Slow_TotalReactPowerInG1_kVArh", "Slow_TotalReactPowerInG1_kVArh_StdDev", "Slow_TotalReactPowerInG1_kVArh_Min", "Slow_TotalReactPowerInG1_kVArh_Max", "Slow_TotalReactPowerInG1_kVArh_Count", "Slow_NacelleDrill", "Slow_NacelleDrill_StdDev", "Slow_NacelleDrill_Min", "Slow_NacelleDrill_Max", "Slow_NacelleDrill_Count", "Slow_TempGearBoxIMSDE", "Slow_TempGearBoxIMSDE_StdDev", "Slow_TempGearBoxIMSDE_Min", "Slow_TempGearBoxIMSDE_Max", "Slow_TempGearBoxIMSDE_Count", "Fast_Total_Operating_hrs", "Fast_Total_Operating_hrs_StdDev", "Fast_Total_Operating_hrs_Min", "Fast_Total_Operating_hrs_Max", "Fast_Total_Operating_hrs_Count", "Slow_TempNacelle", "Slow_TempNacelle_StdDev", "Slow_TempNacelle_Min", "Slow_TempNacelle_Max", "Slow_TempNacelle_Count", "Fast_Total_Grid_OK_hrs", "Fast_Total_Grid_OK_hrs_StdDev", "Fast_Total_Grid_OK_hrs_Min", "Fast_Total_Grid_OK_hrs_Max", "Fast_Total_Grid_OK_hrs_Count", "Fast_Total_WTG_OK_hrs", "Fast_Total_WTG_OK_hrs_StdDev", "Fast_Total_WTG_OK_hrs_Min", "Fast_Total_WTG_OK_hrs_Max", "Fast_Total_WTG_OK_hrs_Count", "Slow_TempCabinetTopBox", "Slow_TempCabinetTopBox_StdDev", "Slow_TempCabinetTopBox_Min", "Slow_TempCabinetTopBox_Max", "Slow_TempCabinetTopBox_Count", "Slow_TempGeneratorBearingNDE", "Slow_TempGeneratorBearingNDE_StdDev", "Slow_TempGeneratorBearingNDE_Min", "Slow_TempGeneratorBearingNDE_Max", "Slow_TempGeneratorBearingNDE_Count", "Fast_Total_Access_hrs", "Fast_Total_Access_hrs_StdDev", "Fast_Total_Access_hrs_Min", "Fast_Total_Access_hrs_Max", "Fast_Total_Access_hrs_Count", "Slow_TempBottomPowerSection", "Slow_TempBottomPowerSection_StdDev", "Slow_TempBottomPowerSection_Min", "Slow_TempBottomPowerSection_Max", "Slow_TempBottomPowerSection_Count", "Slow_TempGeneratorBearingDE", "Slow_TempGeneratorBearingDE_StdDev", "Slow_TempGeneratorBearingDE_Min", "Slow_TempGeneratorBearingDE_Max", "Slow_TempGeneratorBearingDE_Count", "Slow_TotalReactPowerIn_kVArh", "Slow_TotalReactPowerIn_kVArh_StdDev", "Slow_TotalReactPowerIn_kVArh_Min", "Slow_TotalReactPowerIn_kVArh_Max", "Slow_TotalReactPowerIn_kVArh_Count", "Slow_TempBottomControlSection", "Slow_TempBottomControlSection_StdDev", "Slow_TempBottomControlSection_Min", "Slow_TempBottomControlSection_Max", "Slow_TempBottomControlSection_Count", "Slow_TempConv1", "Slow_TempConv1_StdDev", "Slow_TempConv1_Min", "Slow_TempConv1_Max", "Slow_TempConv1_Count", "Fast_ActivePowerRated_kW", "Fast_ActivePowerRated_kW_StdDev", "Fast_ActivePowerRated_kW_Min", "Fast_ActivePowerRated_kW_Max", "Fast_ActivePowerRated_kW_Count", "Fast_NodeIP", "Fast_NodeIP_StdDev", "Fast_NodeIP_Min", "Fast_NodeIP_Max", "Fast_NodeIP_Count", "Fast_PitchSpeed1", "Fast_PitchSpeed1_StdDev", "Fast_PitchSpeed1_Min", "Fast_PitchSpeed1_Max", "Fast_PitchSpeed1_Count", "Slow_CFCardSize", "Slow_CFCardSize_StdDev", "Slow_CFCardSize_Min", "Slow_CFCardSize_Max", "Slow_CFCardSize_Count", "Slow_CPU_Number", "Slow_CPU_Number_StdDev", "Slow_CPU_Number_Min", "Slow_CPU_Number_Max", "Slow_CPU_Number_Count", "Slow_CFCardSpaceLeft", "Slow_CFCardSpaceLeft_StdDev", "Slow_CFCardSpaceLeft_Min", "Slow_CFCardSpaceLeft_Max", "Slow_CFCardSpaceLeft_Count", "Slow_TempBottomCapSection", "Slow_TempBottomCapSection_StdDev", "Slow_TempBottomCapSection_Min", "Slow_TempBottomCapSection_Max", "Slow_TempBottomCapSection_Count", "Slow_RatedPower", "Slow_RatedPower_StdDev", "Slow_RatedPower_Min", "Slow_RatedPower_Max", "Slow_RatedPower_Count", "Slow_TempConv3", "Slow_TempConv3_StdDev", "Slow_TempConv3_Min", "Slow_TempConv3_Max", "Slow_TempConv3_Count", "Slow_TempConv2", "Slow_TempConv2_StdDev", "Slow_TempConv2_Min", "Slow_TempConv2_Max", "Slow_TempConv2_Count", "Slow_TotalActPowerIn_kWh", "Slow_TotalActPowerIn_kWh_StdDev", "Slow_TotalActPowerIn_kWh_Min", "Slow_TotalActPowerIn_kWh_Max", "Slow_TotalActPowerIn_kWh_Count", "Slow_TotalActPowerInG1_kWh", "Slow_TotalActPowerInG1_kWh_StdDev", "Slow_TotalActPowerInG1_kWh_Min", "Slow_TotalActPowerInG1_kWh_Max", "Slow_TotalActPowerInG1_kWh_Count", "Slow_TotalActPowerInG2_kWh", "Slow_TotalActPowerInG2_kWh_StdDev", "Slow_TotalActPowerInG2_kWh_Min", "Slow_TotalActPowerInG2_kWh_Max", "Slow_TotalActPowerInG2_kWh_Count", "Slow_TotalActPowerOutG2_kWh", "Slow_TotalActPowerOutG2_kWh_StdDev", "Slow_TotalActPowerOutG2_kWh_Min", "Slow_TotalActPowerOutG2_kWh_Max", "Slow_TotalActPowerOutG2_kWh_Count", "Slow_TotalG2ActiveHours", "Slow_TotalG2ActiveHours_StdDev", "Slow_TotalG2ActiveHours_Min", "Slow_TotalG2ActiveHours_Max", "Slow_TotalG2ActiveHours_Count", "Slow_TotalReactPowerInG2_kVArh", "Slow_TotalReactPowerInG2_kVArh_StdDev", "Slow_TotalReactPowerInG2_kVArh_Min", "Slow_TotalReactPowerInG2_kVArh_Max", "Slow_TotalReactPowerInG2_kVArh_Count", "Slow_TotalReactPowerOut_kVArh", "Slow_TotalReactPowerOut_kVArh_StdDev", "Slow_TotalReactPowerOut_kVArh_Min", "Slow_TotalReactPowerOut_kVArh_Max", "Slow_TotalReactPowerOut_kVArh_Count", "Slow_UTCoffset_int", "Slow_UTCoffset_int_StdDev", "Slow_UTCoffset_int_Min", "Slow_UTCoffset_int_Max", "Slow_UTCoffset_int_Count"}

	for i, str := range _ashfd_field {
		tkm := tk.M{}.
			Set("_id", str).
			Set("label", _ashfd_label[i]).
			Set("source", "ScadaDataHFD")

		atkm = append(atkm, tkm)
	}

	return atkm
}

// GET DATA

func (m *DataBrowserController) GetScadaList(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson
	var filter []*dbox.Filter

	p := new(helper.Payloads)
	e := k.GetPayload(&p)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	filter, _ = p.ParseFilter()

	query := DB().Connection.NewQuery().From(new(ScadaData).TableName()).Skip(p.Skip).Take(p.Take)
	query.Where(dbox.And(filter...))

	if len(p.Sort) > 0 {
		var arrsort []string
		for _, val := range p.Sort {
			if val.Dir == "desc" {
				arrsort = append(arrsort, strings.ToLower("-"+val.Field))
			} else {
				arrsort = append(arrsort, strings.ToLower(val.Field))
			}
		}
		query = query.Order(arrsort...)
	}
	csr, e := query.Cursor(nil)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}
	defer csr.Close()

	tmpResult := make([]ScadaData, 0)
	e = csr.Fetch(&tmpResult, 0, false)

	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	totalPower := 0.0
	totalPowerLost := 0.0
	totalTurbine := 0
	totalProduction := 0.0
	sumWindSpeed := 0.0
	countData := 0.0
	AvgWindSpeed := 0.0

	aggrData := []tk.M{}

	queryAggr := DB().Connection.NewQuery().From(new(ScadaData).TableName()).
		Aggr(dbox.AggrSum, "$power", "TotalPower").
		Aggr(dbox.AggrSum, "$powerlost", "TotalPowerLost").
		Aggr(dbox.AggrSum, "$power", "totalProduction").
		Aggr(dbox.AggrSum, "$avgwindspeed", "sumWindSpeed").
		Aggr(dbox.AggrSum, 1, "countData").
		Group("turbine").Where(dbox.And(filter...))

	caggr, e := queryAggr.Cursor(nil)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}
	defer caggr.Close()
	e = caggr.Fetch(&aggrData, 0, false)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	for _, val := range aggrData {
		totalPower += val.GetFloat64("TotalPower")
		totalPowerLost += val.GetFloat64("TotalPowerLost")
		totalProduction += val.GetFloat64("totalProduction")
		sumWindSpeed += val.GetFloat64("sumWindSpeed")
		countData += val.GetFloat64("countData")
	}
	totalTurbine = tk.SliceLen(aggrData)

	if countData > 0.0 {
		AvgWindSpeed = sumWindSpeed / countData
	}

	data := struct {
		Data            []ScadaData
		Total           float64
		TotalPower      float64
		TotalPowerLost  float64
		TotalProduction float64
		AvgWindSpeed    float64
		TotalTurbine    int
	}{
		Data:            tmpResult,
		Total:           countData,
		TotalPower:      totalPower,
		TotalPowerLost:  totalPowerLost,
		TotalProduction: totalProduction, // / 6,
		AvgWindSpeed:    AvgWindSpeed,    //sumWindSpeed / countData,
		TotalTurbine:    totalTurbine,
	}

	return helper.CreateResult(true, data, "success")
}

func (m *DataBrowserController) GetScadaAnomalyList(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	p := new(helper.Payloads)
	e := k.GetPayload(&p)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	filter, _ := p.ParseFilter()

	query := DB().Connection.NewQuery().From(new(ScadaAlarmAnomaly).TableName()).Skip(p.Skip).Take(p.Take)
	query.Where(dbox.And(filter...))

	if len(p.Sort) > 0 {
		var arrsort []string
		for _, val := range p.Sort {
			if val.Dir == "desc" {
				arrsort = append(arrsort, strings.ToLower("-"+val.Field))
			} else {
				arrsort = append(arrsort, strings.ToLower(val.Field))
			}
		}
		query = query.Order(arrsort...)
	}
	csr, e := query.Cursor(nil)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}
	defer csr.Close()

	tmpResult := make([]ScadaAlarmAnomaly, 0)
	e = csr.Fetch(&tmpResult, 0, false)

	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	queryC := DB().Connection.NewQuery().From(new(ScadaAlarmAnomaly).TableName()).Where(dbox.And(filter...))
	ccount, e := queryC.Cursor(nil)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}
	defer ccount.Close()

	totalPower := 0.0
	totalPowerLost := 0.0
	sumWindSpeed := 0.0
	totalTurbine := 0
	AvgWS := 0.0

	aggrData := []tk.M{}

	queryAggr := DB().Connection.NewQuery().From(new(ScadaAlarmAnomaly).TableName()).
		Aggr(dbox.AggrSum, "$power", "TotalPower").
		Aggr(dbox.AggrSum, "$powerlost", "TotalPowerLost").
		Aggr(dbox.AggrSum, "$avgwindspeed", "sumWindSpeed").
		Group("turbine").Where(dbox.And(filter...))

	caggr, e := queryAggr.Cursor(nil)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}
	defer caggr.Close()
	e = caggr.Fetch(&aggrData, 0, false)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	for _, val := range aggrData {
		totalPower += val.GetFloat64("TotalPower")
		totalPowerLost += val.GetFloat64("TotalPowerLost")
		sumWindSpeed += val.GetFloat64("sumWindSpeed")
	}
	totalTurbine = tk.SliceLen(aggrData)

	if ccount.Count() > 0.0 {
		AvgWS = sumWindSpeed / float64(ccount.Count())
	}

	data := struct {
		Data            []ScadaAlarmAnomaly
		Total           int
		TotalPower      float64
		TotalPowerLost  float64
		TotalProduction float64
		AvgWindSpeed    float64
		TotalTurbine    int
	}{
		Data:            tmpResult,
		Total:           ccount.Count(),
		TotalPower:      totalPower,
		TotalPowerLost:  totalPowerLost,
		TotalProduction: totalPower / 6,
		AvgWindSpeed:    AvgWS, //sumWindSpeed / float64(ccount.Count()),
		TotalTurbine:    totalTurbine,
	}

	return helper.CreateResult(true, data, "success")
}

func (m *DataBrowserController) GetAlarmList(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	p := new(helper.Payloads)

	e := k.GetPayload(&p)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	filter, _ := p.ParseFilter()
	query := DB().Connection.NewQuery().From(new(Alarm).TableName()).
		Skip(p.Skip).Take(p.Take).Where(dbox.And(filter...))

	if len(p.Sort) > 0 {
		var arrsort []string
		for _, val := range p.Sort {
			if val.Dir == "desc" {
				arrsort = append(arrsort, strings.ToLower("-"+val.Field))
			} else {
				arrsort = append(arrsort, strings.ToLower(val.Field))
			}
		}
		query = query.Order(arrsort...)
	}

	csr, e := query.Cursor(nil)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}
	defer csr.Close()

	tmpResult := make([]Alarm, 0)
	e = csr.Fetch(&tmpResult, 0, false)

	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	queryC := DB().Connection.NewQuery().From(new(Alarm).TableName()).Where(dbox.And(filter...))

	ccount, e := queryC.Cursor(nil)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}
	defer ccount.Close()

	totalTurbine := 0
	aggrData := []tk.M{}

	queryAggr := DB().Connection.NewQuery().From(new(Alarm).TableName()).
		Group("turbine").Where(dbox.And(filter...))

	caggr, e := queryAggr.Cursor(nil)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}
	defer caggr.Close()
	e = caggr.Fetch(&aggrData, 0, false)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	totalTurbine = tk.SliceLen(aggrData)

	data := struct {
		Data         []Alarm
		Total        int
		TotalTurbine int
	}{
		Data:         tmpResult,
		Total:        ccount.Count(),
		TotalTurbine: totalTurbine,
	}

	return helper.CreateResult(true, data, "success")
}

func (m *DataBrowserController) GetAlarmScadaAnomalyList(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	p := new(helper.Payloads)

	e := k.GetPayload(&p)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	filter, _ := p.ParseFilter()
	query := DB().Connection.NewQuery().From(new(AlarmScadaAnomaly).TableName()).
		Skip(p.Skip).Take(p.Take).Where(dbox.And(filter...))

	if len(p.Sort) > 0 {
		var arrsort []string
		for _, val := range p.Sort {
			if val.Dir == "desc" {
				arrsort = append(arrsort, strings.ToLower("-"+val.Field))
			} else {
				arrsort = append(arrsort, strings.ToLower(val.Field))
			}
		}
		query = query.Order(arrsort...)
	}

	csr, e := query.Cursor(nil)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}
	defer csr.Close()

	tmpResult := make([]AlarmScadaAnomaly, 0)
	e = csr.Fetch(&tmpResult, 0, false)

	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	queryC := DB().Connection.NewQuery().From(new(AlarmScadaAnomaly).TableName()).
		Where(dbox.And(filter...))

	ccount, e := queryC.Cursor(nil)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}
	defer ccount.Close()

	totalTurbine := 0
	aggrData := []tk.M{}

	queryAggr := DB().Connection.NewQuery().From(new(AlarmScadaAnomaly).TableName()).
		Group("turbine").Where(dbox.And(filter...))

	caggr, e := queryAggr.Cursor(nil)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}
	defer caggr.Close()
	e = caggr.Fetch(&aggrData, 0, false)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	totalTurbine = tk.SliceLen(aggrData)

	data := struct {
		Data         []AlarmScadaAnomaly
		Total        int
		TotalTurbine int
	}{
		Data:         tmpResult,
		Total:        ccount.Count(),
		TotalTurbine: totalTurbine,
	}

	return helper.CreateResult(true, data, "success")
}

func (m *DataBrowserController) GetAlarmOverlappingList(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	p := new(helper.Payloads)

	e := k.GetPayload(&p)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	filter, _ := p.ParseFilter()
	query := DB().Connection.NewQuery().From(new(AlarmOverlapping).TableName()).
		Skip(p.Skip).Take(p.Take).Where(dbox.And(filter...))

	if len(p.Sort) > 0 {
		var arrsort []string
		for _, val := range p.Sort {
			if val.Dir == "desc" {
				arrsort = append(arrsort, strings.ToLower("-"+val.Field))
			} else {
				arrsort = append(arrsort, strings.ToLower(val.Field))
			}
		}
		query = query.Order(arrsort...)
	}

	csr, e := query.Cursor(nil)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}
	defer csr.Close()

	tmpResult := make([]AlarmOverlapping, 0)
	e = csr.Fetch(&tmpResult, 0, false)

	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	queryC := DB().Connection.NewQuery().From(new(AlarmOverlapping).TableName()).
		Where(dbox.And(filter...))

	ccount, e := queryC.Cursor(nil)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}
	defer ccount.Close()

	totalTurbine := 0
	aggrData := []tk.M{}

	queryAggr := DB().Connection.NewQuery().From(new(AlarmOverlapping).TableName()).
		Group("turbine").Where(dbox.And(filter...))

	caggr, e := queryAggr.Cursor(nil)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}
	defer caggr.Close()
	e = caggr.Fetch(&aggrData, 0, false)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	totalTurbine = tk.SliceLen(aggrData)

	data := struct {
		Data         []AlarmOverlapping
		Total        int
		TotalTurbine int
	}{
		Data:         tmpResult,
		Total:        ccount.Count(),
		TotalTurbine: totalTurbine,
	}

	return helper.CreateResult(true, data, "success")
}

func (m *DataBrowserController) GetAlarmOverlappingDetails(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	p := new(helper.Payloads)

	e := k.GetPayload(&p)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	filter, _ := p.ParseFilter()

	var (
		pipes []tk.M
	)

	pipes = append(pipes, tk.M{"$unwind": "$alarms"})
	pipes = append(pipes, tk.M{"$project": tk.M{
		"farm":             "$alarms.farm",
		"startdate":        "$alarms.startdate",
		"enddate":          "$alarms.enddate",
		"turbine":          "$alarms.turbine",
		"alertdescription": "$alarms.alertdescription",
		"externalstop":     "$alarms.externalstop",
		"griddown":         "$alarms.griddown",
		"internalgrid":     "$alarms.internalgrid",
		"machinedown":      "$alarms.machinedown",
		"aebok":            "$alarms.aebok",
		"unknown":          "$alarms.unknown",
		"weatherstop":      "$alarms.weatherstop",
		"line":             "$alarms.line",
	}})

	query := DB().Connection.NewQuery().From(new(AlarmOverlapping).TableName()).
		Command("pipe", pipes).Where(dbox.And(filter...))

	if len(p.Sort) > 0 {
		var arrsort []string
		for _, val := range p.Sort {
			if val.Dir == "desc" {
				arrsort = append(arrsort, strings.ToLower("-"+val.Field))
			} else {
				arrsort = append(arrsort, strings.ToLower(val.Field))
			}
		}
		query = query.Order(arrsort...)
	}
	queryCount := query

	query.Skip(p.Skip).Take(p.Take)
	csr, e := query.Cursor(nil)
	if e != nil {
		return e.Error()
	}
	defer csr.Close()

	tmpResult := make([]Alarm, 0)
	e = csr.Fetch(&tmpResult, 0, false)

	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	ccount, e := queryCount.Cursor(nil)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}
	defer ccount.Close()

	data := struct {
		Data  []Alarm
		Total int
	}{
		Data:  tmpResult,
		Total: ccount.Count(),
	}

	return helper.CreateResult(true, data, "success")
}

func (m *DataBrowserController) GetJMRList(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	p := new(helper.Payloads)
	e := k.GetPayload(&p)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}
	filter, _ := p.ParseFilter()

	query := DB().Connection.NewQuery().From(new(JMR).TableName()).Skip(p.Skip).Take(p.Take)
	query.Where(dbox.And(filter...))

	if len(p.Sort) > 0 {
		var arrsort []string
		for _, val := range p.Sort {
			if val.Dir == "desc" {
				arrsort = append(arrsort, strings.ToLower("-"+val.Field))
			} else {
				arrsort = append(arrsort, strings.ToLower(val.Field))
			}
		}
		query = query.Order(arrsort...)
	}
	csr, e := query.Cursor(nil)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}
	defer csr.Close()

	tmpResult := make([]JMR, 0)
	e = csr.Fetch(&tmpResult, 0, false)

	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	queryC := DB().Connection.NewQuery().From(new(JMR).TableName()).Where(dbox.And(filter...))
	ccount, e := queryC.Cursor(nil)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}
	defer ccount.Close()

	data := struct {
		Data  []JMR
		Total int
	}{
		Data:  tmpResult,
		Total: ccount.Count(),
	}

	return helper.CreateResult(true, data, "success")
}

func (m *DataBrowserController) GetJMRDetails(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	p := new(helper.Payloads)

	e := k.GetPayload(&p)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	filter, _ := p.ParseFilter()

	query := DB().Connection.NewQuery().From(new(JMR).TableName())
	query.Where(dbox.And(filter...))
	csr, e := query.Cursor(nil)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	jmrResult := make([]JMR, 0)
	e = csr.Fetch(&jmrResult, 0, false)

	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}
	defer csr.Close()

	result := make([]JMRSection, 0)

	turbines := make(map[string]bool, 0)

	for _, fil := range filter {
		if fil.Field == "sections.turbine" {
			for _, str := range fil.Value.([]interface{}) {
				turbines[tk.ToString(str)] = true
			}
		}
	}

	clean := tk.M{}

	for _, jmr := range jmrResult {
		sectionsClean := []JMRSection{}
		for _, section := range jmr.Sections {
			if turbines[section.Turbine] {
				sectionsClean = append(sectionsClean, section)
			}
		}

		if len(sectionsClean) > 0 {
			clean.Set(jmr.DateInfo.MonthDesc, sectionsClean)
		}
	}

	for _, jmr := range jmrResult {
		for _, total := range jmr.TotalDetails {

			var contrGenTotal float64
			var boEExportTotal float64
			var boEImportTotal float64
			var boENetTotal float64

			sectionsClean := clean.Get(jmr.DateInfo.MonthDesc).([]JMRSection)

			for _, section := range sectionsClean {
				if total.Section == section.Description {
					result = append(result, section)
					contrGenTotal += section.ContrGen
					boEExportTotal += section.BoEExport
					boEImportTotal += section.BoEImport
					boENetTotal += section.BoENet
				}
			}

			if contrGenTotal != 0 {
				tmpSection := JMRSection{}
				tmpSection.Company = "Total"
				tmpSection.ContrGen = contrGenTotal
				tmpSection.BoEExport = boEExportTotal
				tmpSection.BoEImport = boEImportTotal
				tmpSection.BoENet = boENetTotal

				result = append(result, tmpSection)
			}

		}
	}

	return helper.CreateResult(true, result, "success")
}

func (m *DataBrowserController) GetMETList(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	p := new(helper.Payloads)
	e := k.GetPayload(&p)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}
	filter, _ := p.ParseFilter()

	query := DB().Connection.NewQuery().From(new(MetTower).TableName()).Skip(p.Skip).Take(p.Take)
	query.Where(dbox.And(filter...))

	if len(p.Sort) > 0 {
		var arrsort []string
		for _, val := range p.Sort {
			if val.Dir == "desc" {
				arrsort = append(arrsort, strings.ToLower("-"+val.Field))
			} else {
				arrsort = append(arrsort, strings.ToLower(val.Field))
			}
		}
		query = query.Order(arrsort...)
	}
	csr, e := query.Cursor(nil)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}
	defer csr.Close()

	tmpResult := make([]MetTower, 0)
	e = csr.Fetch(&tmpResult, 0, false)

	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	queryC := DB().Connection.NewQuery().From(new(MetTower).TableName()).Where(dbox.And(filter...))
	ccount, e := queryC.Cursor(nil)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}
	defer ccount.Close()

	data := struct {
		Data  []MetTower
		Total int
	}{
		Data:  tmpResult,
		Total: ccount.Count(),
	}

	return helper.CreateResult(true, data, "success")
}

func (m *DataBrowserController) GetEventList(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson
	var filter []*dbox.Filter

	p := new(helper.Payloads)
	e := k.GetPayload(&p)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	filter, _ = p.ParseFilter()

	query := DB().Connection.NewQuery().From(new(EventRaw).TableName()).Skip(p.Skip).Take(p.Take)
	query.Where(dbox.And(filter...))

	if len(p.Sort) > 0 {
		var arrsort []string
		for _, val := range p.Sort {
			if val.Dir == "desc" {
				arrsort = append(arrsort, strings.ToLower("-"+val.Field))
			} else {
				arrsort = append(arrsort, strings.ToLower(val.Field))
			}
		}
		query = query.Order(arrsort...)
	}
	csr, e := query.Cursor(nil)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}
	defer csr.Close()

	results := make([]EventRaw, 0)
	e = csr.Fetch(&results, 0, false)

	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	totalTurbine := 0
	countData := 0.0

	aggrData := []tk.M{}

	queryAggr := DB().Connection.NewQuery().From(new(EventRaw).TableName()).
		Aggr(dbox.AggrSum, 1, "countData").
		Group("turbine").Where(dbox.And(filter...))

	caggr, e := queryAggr.Cursor(nil)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}
	defer caggr.Close()
	e = caggr.Fetch(&aggrData, 0, false)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	for _, val := range aggrData {
		countData += val.GetFloat64("countData")
	}
	totalTurbine = tk.SliceLen(aggrData)

	data := struct {
		Data         []EventRaw
		Total        float64
		TotalTurbine int
	}{
		Data:         results,
		Total:        countData,
		TotalTurbine: totalTurbine,
	}

	return helper.CreateResult(true, data, "success")
}

func (m *DataBrowserController) GetScadaOemList(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	p := new(helper.Payloads)
	e := k.GetPayload(&p)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}
	filter, _ := p.ParseFilter()

	query := DB().Connection.NewQuery().From(new(ScadaDataOEM).TableName()).Skip(p.Skip).Take(p.Take)
	query.Where(dbox.And(filter...))

	if len(p.Sort) > 0 {
		var arrsort []string
		for _, val := range p.Sort {
			if val.Dir == "desc" {
				arrsort = append(arrsort, strings.ToLower("-"+val.Field))
			} else {
				arrsort = append(arrsort, strings.ToLower(val.Field))
			}
		}
		query = query.Order(arrsort...)
	}
	csr, e := query.Cursor(nil)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}
	defer csr.Close()

	tmpResult := make([]ScadaDataOEM, 0)
	e = csr.Fetch(&tmpResult, 0, false)

	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	queryC := DB().Connection.NewQuery().From(new(ScadaDataOEM).TableName()).Where(dbox.And(filter...))
	ccount, e := queryC.Cursor(nil)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}
	defer ccount.Close()

	totalPower := 0.0
	totalPowerLost := 0.0
	totalActivePower := 0.0
	avgWindSpeed := 0.0
	totalTurbine := 0
	totalEnergy := 0.0
	AvgWS := 0.0

	aggrData := []tk.M{}

	queryAggr := DB().Connection.NewQuery().From(new(ScadaDataOEM).TableName()).
		Aggr(dbox.AggrSum, "$power", "TotalPower").
		Aggr(dbox.AggrSum, "$powerlost", "TotalPowerLost").
		Aggr(dbox.AggrSum, "$ai_intern_activpower", "TotalActivePower").
		Aggr(dbox.AggrSum, "$ai_intern_windspeed", "AvgWindSpeed").
		Aggr(dbox.AggrSum, "$energy", "TotalEnergy").
		Group("turbine").Where(dbox.And(filter...))

	caggr, e := queryAggr.Cursor(nil)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}
	defer caggr.Close()
	e = caggr.Fetch(&aggrData, 0, false)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	for _, val := range aggrData {
		totalPower += val.GetFloat64("TotalPower")
		totalPowerLost += val.GetFloat64("TotalPowerLost")
		totalActivePower += val.GetFloat64("TotalActivePower")
		avgWindSpeed += val.GetFloat64("AvgWindSpeed")
		totalEnergy += val.GetFloat64("TotalEnergy")
	}
	totalTurbine = tk.SliceLen(aggrData)

	if ccount.Count() > 0.0 {
		AvgWS = avgWindSpeed / float64(ccount.Count())
	}

	data := struct {
		Data             []ScadaDataOEM
		Total            int
		TotalPower       float64
		TotalPowerLost   float64
		TotalActivePower float64
		AvgWindSpeed     float64
		TotalTurbine     int
		TotalEnergy      float64
	}{
		Data:             tmpResult,
		Total:            ccount.Count(),
		TotalPower:       totalPower,
		TotalPowerLost:   totalPowerLost,
		TotalActivePower: totalActivePower,
		AvgWindSpeed:     AvgWS, //avgWindSpeed / float64(ccount.Count()),
		TotalTurbine:     totalTurbine,
		TotalEnergy:      totalActivePower / 6,
	}

	return helper.CreateResult(true, data, "success")
}

func (m *DataBrowserController) GetCustomList(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	p := new(helper.Payloads)
	e := k.GetPayload(&p)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}
	filter, _ := p.ParseFilter()

	istimestamp := false
	arrscadaoem := []string{"_id", "timestamputc"}
	arrmettower := []string{}
	if p.Custom.Has("ColumnList") {
		for _, _val := range p.Custom["ColumnList"].([]interface{}) {
			_tkm, _ := tk.ToM(_val)
			if _tkm.GetString("source") == "ScadaDataOEM" {
				arrscadaoem = append(arrscadaoem, _tkm.GetString("_id"))
				if _tkm.GetString("_id") == "timestamp" {
					istimestamp = true
				}
			} else if _tkm.GetString("source") == "MetTower" {
				arrmettower = append(arrmettower, _tkm.GetString("_id"))
			}
		}
	}

	query := DB().Connection.NewQuery().
		Select(arrscadaoem...).
		From(new(ScadaDataOEM).TableName()).
		Skip(p.Skip).
		Take(p.Take)
	query.Where(dbox.And(filter...))

	if len(p.Sort) > 0 {
		var arrsort []string
		for _, val := range p.Sort {
			if val.Dir == "desc" {
				arrsort = append(arrsort, strings.ToLower("-"+val.Field))
			} else {
				arrsort = append(arrsort, strings.ToLower(val.Field))
			}
		}
		query = query.Order(arrsort...)
	}
	csr, e := query.Cursor(nil)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}
	defer csr.Close()

	results := make([]tk.M, 0)
	e = csr.Fetch(&results, 0, false)

	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	arrmettowercond := []interface{}{}

	for i, val := range results {
		if val.Has("timestamputc") {
			itime := val.Get("timestamputc", time.Time{}).(time.Time).UTC()
			arrmettowercond = append(arrmettowercond, itime)
			val.Set("timestamputc", itime)
			results[i] = val
		}
		if istimestamp {
			itime := val.Get("timestamp", time.Time{}).(time.Time)
			val.Set("timestamp", itime.UTC())
			results[i] = val
		}
	}

	tkmmet := tk.M{}
	if len(arrmettower) > 0 && len(arrmettowercond) > 0 {
		arrmettower = append(arrmettower, "timestamp")
		_csr, _e := DB().Connection.NewQuery().
			Select(arrmettower...).
			From("MetTower").
			Where(dbox.In("timestamp", arrmettowercond...)).Cursor(nil)
		if _e != nil {
			return helper.CreateResult(false, nil, _e.Error())
		}
		defer _csr.Close()

		_resmet := make([]tk.M, 0)
		_e = _csr.Fetch(&_resmet, 0, false)

		if _e != nil {
			return helper.CreateResult(false, nil, _e.Error())
		}

		for _, val := range _resmet {
			itime := val.Get("timestamp", time.Time{}).(time.Time).UTC().String()
			tkmmet.Set(itime, val)
		}
	}

	if len(tkmmet) > 0 {
		for i, val := range results {
			itime := val.Get("timestamputc", time.Time{}).(time.Time).UTC().String()
			if tkmmet.Has(itime) {
				for _key, _val := range tkmmet[itime].(tk.M) {
					if _key != "timestamp" {
						val.Set(_key, _val)
					}
				}
			}
			val.Unset("timestamputc")
			results[i] = val
		}
	}

	queryC := DB().Connection.NewQuery().From(new(ScadaDataOEM).TableName()).Where(dbox.And(filter...))
	ccount, e := queryC.Cursor(nil)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}
	defer ccount.Close()

	totalPower := 0.0
	totalPowerLost := 0.0
	totalActivePower := 0.0
	avgWindSpeed := 0.0
	totalTurbine := 0
	totalEnergy := 0.0
	AvgWS := 0.0

	aggrData := []tk.M{}

	queryAggr := DB().Connection.NewQuery().From(new(ScadaDataOEM).TableName()).
		Aggr(dbox.AggrSum, "$power", "TotalPower").
		Aggr(dbox.AggrSum, "$powerlost", "TotalPowerLost").
		Aggr(dbox.AggrSum, "$ai_intern_activpower", "TotalActivePower").
		Aggr(dbox.AggrSum, "$ai_intern_windspeed", "AvgWindSpeed").
		Aggr(dbox.AggrSum, "$energy", "TotalEnergy").
		Group("turbine").Where(dbox.And(filter...))

	caggr, e := queryAggr.Cursor(nil)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}
	defer caggr.Close()
	e = caggr.Fetch(&aggrData, 0, false)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	for _, val := range aggrData {
		totalPower += val.GetFloat64("TotalPower")
		totalPowerLost += val.GetFloat64("TotalPowerLost")
		totalActivePower += val.GetFloat64("TotalActivePower")
		avgWindSpeed += val.GetFloat64("AvgWindSpeed")
		totalEnergy += val.GetFloat64("TotalEnergy")
	}
	totalTurbine = tk.SliceLen(aggrData)

	if ccount.Count() > 0.0 {
		AvgWS = avgWindSpeed / float64(ccount.Count())
	}

	data := struct {
		Data             []tk.M
		Total            int
		TotalPower       float64
		TotalPowerLost   float64
		TotalActivePower float64
		AvgWindSpeed     float64
		TotalTurbine     int
		TotalEnergy      float64
	}{
		Data:             results,
		Total:            ccount.Count(),
		TotalPower:       totalPower,
		TotalPowerLost:   totalPowerLost,
		TotalActivePower: totalActivePower,
		AvgWindSpeed:     AvgWS, //avgWindSpeed / float64(ccount.Count()),
		TotalTurbine:     totalTurbine,
		TotalEnergy:      totalEnergy,
	}

	return helper.CreateResult(true, data, "success")
}

func (m *DataBrowserController) GetDowntimeEventList(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	var filter []*dbox.Filter

	p := new(helper.Payloads)
	e := k.GetPayload(&p)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	filter, _ = p.ParseFilter()

	query := DB().Connection.NewQuery().From(new(EventDown).TableName()).Skip(p.Skip).Take(p.Take)
	query.Where(dbox.And(filter...))

	if len(p.Sort) > 0 {
		var arrsort []string
		for _, val := range p.Sort {
			if val.Dir == "desc" {
				arrsort = append(arrsort, strings.ToLower("-"+val.Field))
			} else {
				arrsort = append(arrsort, strings.ToLower(val.Field))
			}
		}
		query = query.Order(arrsort...)
	}
	csr, e := query.Cursor(nil)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}
	defer csr.Close()

	tmpResult := make([]EventDown, 0)
	e = csr.Fetch(&tmpResult, 0, false)

	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	queryC := DB().Connection.NewQuery().From(new(EventDown).TableName()).Where(dbox.And(filter...))
	ccount, e := queryC.Cursor(nil)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}
	defer ccount.Close()

	totalDuration := 0.0
	totalTurbine := 0

	aggrData := []tk.M{}

	queryAggr := DB().Connection.NewQuery().From(new(EventDown).TableName()).
		Aggr(dbox.AggrSum, "$duration", "duration").
		Group("turbine").Where(dbox.And(filter...))

	caggr, e := queryAggr.Cursor(nil)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}
	defer caggr.Close()
	e = caggr.Fetch(&aggrData, 0, false)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	for _, val := range aggrData {
		totalDuration += val.GetFloat64("duration")
	}
	totalTurbine = tk.SliceLen(aggrData)

	data := struct {
		Data          []EventDown
		Total         int
		TotalTurbine  int
		TotalDuration float64
	}{
		Data:          tmpResult,
		Total:         ccount.Count(),
		TotalTurbine:  totalTurbine,
		TotalDuration: totalDuration,
	}

	return helper.CreateResult(true, data, "success")
}

func (m *DataBrowserController) GetScadaHFDList(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	p := new(helper.Payloads)
	e := k.GetPayload(&p)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}
	filter, _ := p.ParseFilter()

	query := DB().Connection.NewQuery().From(new(ScadaDataHFD).TableName()).Skip(p.Skip).Take(p.Take)
	query.Where(dbox.And(filter...))

	if len(p.Sort) > 0 {
		var arrsort []string
		for _, val := range p.Sort {
			if val.Dir == "desc" {
				arrsort = append(arrsort, strings.ToLower("-"+val.Field))
			} else {
				arrsort = append(arrsort, strings.ToLower(val.Field))
			}
		}
		query = query.Order(arrsort...)
	}
	csr, e := query.Cursor(nil)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}
	defer csr.Close()

	tmpResult := make([]ScadaDataHFD, 0)
	e = csr.Fetch(&tmpResult, 0, false)

	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	queryC := DB().Connection.NewQuery().From(new(ScadaDataHFD).TableName()).Where(dbox.And(filter...))
	ccount, e := queryC.Cursor(nil)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}
	defer ccount.Close()

	totalTurbine := 0

	aggrData := []tk.M{}

	queryAggr := DB().Connection.NewQuery().From(new(ScadaDataHFD).TableName()).
		Group("turbine").Where(dbox.And(filter...))

	caggr, e := queryAggr.Cursor(nil)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}
	defer caggr.Close()
	e = caggr.Fetch(&aggrData, 0, false)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	totalTurbine = tk.SliceLen(aggrData)

	project := ""
	for _, val := range filter {
		if val.Field == "projectname" {
			project = tk.ToString(val.Value)
		}
	}
	turbineName, err := helper.GetTurbineNameList(project)
	if err != nil {
		return helper.CreateResult(false, nil, e.Error())
	}
	for idx, val := range tmpResult {
		tmpResult[idx].Turbine = turbineName[val.Turbine]
	}

	data := struct {
		Data         []ScadaDataHFD
		Total        int
		TotalTurbine int
	}{
		Data:         tmpResult,
		Total:        ccount.Count(),
		TotalTurbine: totalTurbine,
	}

	return helper.CreateResult(true, data, "success")
}

// Get date info each tab

func (m *DataBrowserController) GetAvailDate(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	Scadaresults := make([]time.Time, 0)
	Alarmresults := make([]time.Time, 0)
	JMRresults := make([]time.Time, 0)
	METresults := make([]time.Time, 0)
	Durationresults := make([]time.Time, 0)
	ScadaAnomalyresults := make([]time.Time, 0)
	AlarmOverlappingresults := make([]time.Time, 0)
	AlarmScadaAnomalyresults := make([]time.Time, 0)

	// Scada Data
	for i := 0; i < 2; i++ {
		var arrsort []string
		if i == 0 {
			arrsort = append(arrsort, "timestamp")
		} else {
			arrsort = append(arrsort, "-timestamp")
		}

		query := DB().Connection.NewQuery().From(new(ScadaData).TableName()).Skip(0).Take(1)
		query = query.Order(arrsort...)

		csr, e := query.Cursor(nil)
		if e != nil {
			return helper.CreateResult(false, nil, e.Error())
		}
		defer csr.Close()

		Result := make([]ScadaData, 0)
		e = csr.Fetch(&Result, 0, false)

		if e != nil {
			return helper.CreateResult(false, nil, e.Error())
		}

		for _, val := range Result {
			Scadaresults = append(Scadaresults, val.TimeStamp.UTC())
		}
	}

	// Alarm Data
	for i := 0; i < 2; i++ {
		var arrsort []string
		if i == 0 {
			arrsort = append(arrsort, "startdate")
		} else {
			arrsort = append(arrsort, "-startdate")
		}

		query := DB().Connection.NewQuery().From(new(Alarm).TableName()).Skip(0).Take(1)
		query = query.Order(arrsort...)

		csr, e := query.Cursor(nil)
		if e != nil {
			return helper.CreateResult(false, nil, e.Error())
		}
		defer csr.Close()

		Result := make([]Alarm, 0)
		e = csr.Fetch(&Result, 0, false)

		if e != nil {
			return helper.CreateResult(false, nil, e.Error())
		}

		for _, val := range Result {
			Alarmresults = append(Alarmresults, val.StartDate.UTC())
		}
	}

	// JMR Data
	for i := 0; i < 2; i++ {
		var arrsort []string
		if i == 0 {
			arrsort = append(arrsort, "dateinfo.dateid")
		} else {
			arrsort = append(arrsort, "-dateinfo.dateid")
		}

		query := DB().Connection.NewQuery().From(new(ScadaData).TableName()).Skip(0).Take(1)
		query = query.Order(arrsort...)

		csr, e := query.Cursor(nil)
		if e != nil {
			return helper.CreateResult(false, nil, e.Error())
		}
		defer csr.Close()

		Result := make([]ScadaData, 0)
		e = csr.Fetch(&Result, 0, false)

		if e != nil {
			return helper.CreateResult(false, nil, e.Error())
		}

		for _, val := range Result {
			JMRresults = append(JMRresults, val.DateInfo.DateId.UTC())
		}
	}

	// MET Data
	for i := 0; i < 2; i++ {
		var arrsort []string
		if i == 0 {
			arrsort = append(arrsort, "timestamp")
		} else {
			arrsort = append(arrsort, "-timestamp")
		}

		query := DB().Connection.NewQuery().From(new(MetTower).TableName()).Skip(0).Take(1)
		query = query.Order(arrsort...)

		csr, e := query.Cursor(nil)
		if e != nil {
			return helper.CreateResult(false, nil, e.Error())
		}
		defer csr.Close()

		Result := make([]MetTower, 0)
		e = csr.Fetch(&Result, 0, false)

		if e != nil {
			return helper.CreateResult(false, nil, e.Error())
		}

		for _, val := range Result {
			METresults = append(METresults, val.TimeStamp.UTC())
		}
	}

	// Duration Data
	for i := 0; i < 2; i++ {
		var arrsort []string
		if i == 0 {
			arrsort = append(arrsort, "timestamp")
		} else {
			arrsort = append(arrsort, "-timestamp")
		}

		query := DB().Connection.NewQuery().From(new(ScadaData).TableName()).Where(dbox.And(dbox.Eq("isvalidtimeduration", false))).Skip(0).Take(1)
		query = query.Order(arrsort...)

		csr, e := query.Cursor(nil)
		if e != nil {
			return helper.CreateResult(false, nil, e.Error())
		}
		defer csr.Close()

		Result := make([]ScadaData, 0)
		e = csr.Fetch(&Result, 0, false)

		if e != nil {
			return helper.CreateResult(false, nil, e.Error())
		}

		for _, val := range Result {
			Durationresults = append(Durationresults, val.TimeStamp.UTC())
		}
	}

	// Anomaly Data
	for i := 0; i < 2; i++ {
		var arrsort []string
		if i == 0 {
			arrsort = append(arrsort, "timestamp")
		} else {
			arrsort = append(arrsort, "-timestamp")
		}

		query := DB().Connection.NewQuery().From(new(ScadaData).TableName()).Where(dbox.And(dbox.Eq("isvalidtimeduration", true))).Skip(0).Take(1)
		query = query.Order(arrsort...)

		csr, e := query.Cursor(nil)
		if e != nil {
			return helper.CreateResult(false, nil, e.Error())
		}
		defer csr.Close()

		Result := make([]ScadaData, 0)
		e = csr.Fetch(&Result, 0, false)

		if e != nil {
			return helper.CreateResult(false, nil, e.Error())
		}

		for _, val := range Result {
			ScadaAnomalyresults = append(ScadaAnomalyresults, val.TimeStamp.UTC())
		}
	}

	// AlarmOverlapping Data
	for i := 0; i < 2; i++ {
		var arrsort []string
		if i == 0 {
			arrsort = append(arrsort, "startdate")
		} else {
			arrsort = append(arrsort, "-startdate")
		}

		query := DB().Connection.NewQuery().From(new(AlarmOverlapping).TableName()).Skip(0).Take(1)
		query = query.Order(arrsort...)

		csr, e := query.Cursor(nil)
		if e != nil {
			return helper.CreateResult(false, nil, e.Error())
		}
		defer csr.Close()

		Result := make([]AlarmOverlapping, 0)
		e = csr.Fetch(&Result, 0, false)

		if e != nil {
			return helper.CreateResult(false, nil, e.Error())
		}

		for _, val := range Result {
			AlarmOverlappingresults = append(AlarmOverlappingresults, val.StartDate.UTC())
		}
	}

	// AlarmScadaAnomaly Data
	for i := 0; i < 2; i++ {
		var arrsort []string
		if i == 0 {
			arrsort = append(arrsort, "startdate")
		} else {
			arrsort = append(arrsort, "-startdate")
		}

		query := DB().Connection.NewQuery().From(new(AlarmScadaAnomaly).TableName()).Skip(0).Take(1)
		query = query.Order(arrsort...)

		csr, e := query.Cursor(nil)
		if e != nil {
			return helper.CreateResult(false, nil, e.Error())
		}
		defer csr.Close()

		Result := make([]AlarmScadaAnomaly, 0)
		e = csr.Fetch(&Result, 0, false)

		if e != nil {
			return helper.CreateResult(false, nil, e.Error())
		}

		for _, val := range Result {
			AlarmScadaAnomalyresults = append(AlarmScadaAnomalyresults, val.StartDate.UTC())
		}
	}

	data := struct {
		ScadaData         []time.Time
		Alarm             []time.Time
		JMR               []time.Time
		MET               []time.Time
		Duration          []time.Time
		ScadaAnomaly      []time.Time
		AlarmOverlapping  []time.Time
		AlarmScadaAnomaly []time.Time
	}{
		ScadaData:         Scadaresults,
		Alarm:             Alarmresults,
		JMR:               JMRresults,
		MET:               METresults,
		Duration:          Durationresults,
		ScadaAnomaly:      ScadaAnomalyresults,
		AlarmOverlapping:  AlarmOverlappingresults,
		AlarmScadaAnomaly: AlarmScadaAnomalyresults,
	}

	return helper.CreateResult(true, data, "success")
}

func (m *DataBrowserController) GetScadaAvailDate(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	Scadaresults := make([]time.Time, 0)

	// Scada Data
	for i := 0; i < 2; i++ {
		var arrsort []string
		if i == 0 {
			arrsort = append(arrsort, "timestamp")
		} else {
			arrsort = append(arrsort, "-timestamp")
		}

		query := DB().Connection.NewQuery().From(new(ScadaData).TableName()).Skip(0).Take(1)
		query = query.Order(arrsort...)

		csr, e := query.Cursor(nil)
		if e != nil {
			return helper.CreateResult(false, nil, e.Error())
		}
		defer csr.Close()

		Result := make([]ScadaData, 0)
		e = csr.Fetch(&Result, 0, false)

		if e != nil {
			return helper.CreateResult(false, nil, e.Error())
		}

		for _, val := range Result {
			Scadaresults = append(Scadaresults, val.TimeStamp.UTC())
		}
	}

	data := struct {
		ScadaData []time.Time
	}{
		ScadaData: Scadaresults,
	}

	return helper.CreateResult(true, data, "success")
}

func (m *DataBrowserController) GetAlarmAvailDate(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	Alarmresults := make([]time.Time, 0)

	// Alarm Data
	for i := 0; i < 2; i++ {
		var arrsort []string
		if i == 0 {
			arrsort = append(arrsort, "startdate")
		} else {
			arrsort = append(arrsort, "-startdate")
		}

		query := DB().Connection.NewQuery().From(new(Alarm).TableName()).Skip(0).Take(1)
		query = query.Order(arrsort...)

		csr, e := query.Cursor(nil)
		if e != nil {
			return helper.CreateResult(false, nil, e.Error())
		}
		defer csr.Close()

		Result := make([]Alarm, 0)
		e = csr.Fetch(&Result, 0, false)

		if e != nil {
			return helper.CreateResult(false, nil, e.Error())
		}

		for _, val := range Result {
			Alarmresults = append(Alarmresults, val.StartDate.UTC())
		}
	}

	data := struct {
		Alarm []time.Time
	}{
		Alarm: Alarmresults,
	}

	return helper.CreateResult(true, data, "success")
}

func (m *DataBrowserController) GetJMRAvailDate(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	JMRresults := make([]time.Time, 0)

	// JMR Data
	for i := 0; i < 2; i++ {
		var arrsort []string
		if i == 0 {
			arrsort = append(arrsort, "dateinfo.dateid")
		} else {
			arrsort = append(arrsort, "-dateinfo.dateid")
		}

		query := DB().Connection.NewQuery().From(new(JMR).TableName()).Skip(0).Take(1)
		query = query.Order(arrsort...)

		csr, e := query.Cursor(nil)
		if e != nil {
			return helper.CreateResult(false, nil, e.Error())
		}
		defer csr.Close()

		Result := make([]ScadaData, 0)
		e = csr.Fetch(&Result, 0, false)

		if e != nil {
			return helper.CreateResult(false, nil, e.Error())
		}

		for _, val := range Result {
			JMRresults = append(JMRresults, val.DateInfo.DateId.UTC())
		}
	}

	data := struct {
		JMR []time.Time
	}{
		JMR: JMRresults,
	}

	return helper.CreateResult(true, data, "success")
}

func (m *DataBrowserController) GetMETAvailDate(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	METresults := make([]time.Time, 0)

	// MET Data
	for i := 0; i < 2; i++ {
		var arrsort []string
		if i == 0 {
			arrsort = append(arrsort, "timestamp")
		} else {
			arrsort = append(arrsort, "-timestamp")
		}

		query := DB().Connection.NewQuery().From(new(MetTower).TableName()).Skip(0).Take(1)
		query = query.Order(arrsort...)

		csr, e := query.Cursor(nil)
		if e != nil {
			return helper.CreateResult(false, nil, e.Error())
		}
		defer csr.Close()

		Result := make([]MetTower, 0)
		e = csr.Fetch(&Result, 0, false)

		if e != nil {
			return helper.CreateResult(false, nil, e.Error())
		}

		for _, val := range Result {
			METresults = append(METresults, val.TimeStamp.UTC())
		}
	}

	data := struct {
		MET []time.Time
	}{
		MET: METresults,
	}

	return helper.CreateResult(true, data, "success")
}

func (m *DataBrowserController) GetDurationAvailDate(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	Durationresults := make([]time.Time, 0)

	// Duration Data
	for i := 0; i < 2; i++ {
		var arrsort []string
		if i == 0 {
			arrsort = append(arrsort, "timestamp")
		} else {
			arrsort = append(arrsort, "-timestamp")
		}

		// query := DB().Connection.NewQuery().From(new(ScadaData).TableName()).Where(dbox.And(dbox.Eq("isvalidtimeduration", false))).Skip(0).Take(1)
		query := DB().Connection.NewQuery().From(new(ScadaData).TableName()).Skip(0).Take(1)
		query = query.Order(arrsort...)

		csr, e := query.Cursor(nil)
		if e != nil {
			return helper.CreateResult(false, nil, e.Error())
		}
		defer csr.Close()

		Result := make([]ScadaData, 0)
		e = csr.Fetch(&Result, 0, false)

		if e != nil {
			return helper.CreateResult(false, nil, e.Error())
		}

		for _, val := range Result {
			Durationresults = append(Durationresults, val.TimeStamp.UTC())
		}
	}

	data := struct {
		Duration []time.Time
	}{
		Duration: Durationresults,
	}

	return helper.CreateResult(true, data, "success")
}

func (m *DataBrowserController) GetScadaAnomalyAvailDate(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	ScadaAnomalyresults := make([]time.Time, 0)

	// Anomaly Data
	for i := 0; i < 2; i++ {
		var arrsort []string
		if i == 0 {
			arrsort = append(arrsort, "timestamp")
		} else {
			arrsort = append(arrsort, "-timestamp")
		}

		query := DB().Connection.NewQuery().From(new(ScadaData).TableName()).Where(dbox.And(dbox.Eq("isvalidtimeduration", true))).Skip(0).Take(1)
		query = query.Order(arrsort...)

		csr, e := query.Cursor(nil)
		if e != nil {
			return helper.CreateResult(false, nil, e.Error())
		}
		defer csr.Close()

		Result := make([]ScadaData, 0)
		e = csr.Fetch(&Result, 0, false)

		if e != nil {
			return helper.CreateResult(false, nil, e.Error())
		}

		for _, val := range Result {
			ScadaAnomalyresults = append(ScadaAnomalyresults, val.TimeStamp.UTC())
		}
	}

	data := struct {
		ScadaAnomaly []time.Time
	}{
		ScadaAnomaly: ScadaAnomalyresults,
	}

	return helper.CreateResult(true, data, "success")
}

func (m *DataBrowserController) GetAlarmOverlappingAvailDate(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	AlarmOverlappingresults := make([]time.Time, 0)

	// AlarmOverlapping Data
	for i := 0; i < 2; i++ {
		var arrsort []string
		if i == 0 {
			arrsort = append(arrsort, "startdate")
		} else {
			arrsort = append(arrsort, "-startdate")
		}

		query := DB().Connection.NewQuery().From(new(AlarmOverlapping).TableName()).Skip(0).Take(1)
		query = query.Order(arrsort...)

		csr, e := query.Cursor(nil)
		if e != nil {
			return helper.CreateResult(false, nil, e.Error())
		}
		defer csr.Close()

		Result := make([]AlarmOverlapping, 0)
		e = csr.Fetch(&Result, 0, false)

		if e != nil {
			return helper.CreateResult(false, nil, e.Error())
		}

		for _, val := range Result {
			AlarmOverlappingresults = append(AlarmOverlappingresults, val.StartDate.UTC())
		}
	}

	data := struct {
		AlarmOverlapping []time.Time
	}{
		AlarmOverlapping: AlarmOverlappingresults,
	}

	return helper.CreateResult(true, data, "success")
}

func (m *DataBrowserController) GetAlarmScadaAnomalyAvailDate(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	AlarmScadaAnomalyresults := make([]time.Time, 0)

	// AlarmScadaAnomaly Data
	for i := 0; i < 2; i++ {
		var arrsort []string
		if i == 0 {
			arrsort = append(arrsort, "startdate")
		} else {
			arrsort = append(arrsort, "-startdate")
		}

		query := DB().Connection.NewQuery().From(new(AlarmScadaAnomaly).TableName()).Skip(0).Take(1)
		query = query.Order(arrsort...)

		csr, e := query.Cursor(nil)
		if e != nil {
			return helper.CreateResult(false, nil, e.Error())
		}
		defer csr.Close()

		Result := make([]AlarmScadaAnomaly, 0)
		e = csr.Fetch(&Result, 0, false)

		if e != nil {
			return helper.CreateResult(false, nil, e.Error())
		}

		for _, val := range Result {
			AlarmScadaAnomalyresults = append(AlarmScadaAnomalyresults, val.StartDate.UTC())
		}
	}

	data := struct {
		AlarmScadaAnomaly []time.Time
	}{
		AlarmScadaAnomaly: AlarmScadaAnomalyresults,
	}

	return helper.CreateResult(true, data, "success")
}

func (m *DataBrowserController) GetEventAvailDate(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	EventDateresults := make([]time.Time, 0)

	// AlarmScadaAnomaly Data
	for i := 0; i < 2; i++ {
		var arrsort []string
		if i == 0 {
			arrsort = append(arrsort, "timestamp")
		} else {
			arrsort = append(arrsort, "-timestamp")
		}

		query := DB().Connection.NewQuery().From(new(EventRaw).TableName()).Skip(0).Take(1)
		query = query.Order(arrsort...)

		csr, e := query.Cursor(nil)
		if e != nil {
			return helper.CreateResult(false, nil, e.Error())
		}
		defer csr.Close()

		Result := make([]EventRaw, 0)
		e = csr.Fetch(&Result, 0, false)

		if e != nil {
			return helper.CreateResult(false, nil, e.Error())
		}

		for _, val := range Result {
			EventDateresults = append(EventDateresults, val.TimeStamp.UTC())
		}
	}

	data := struct {
		EventDate []time.Time
	}{
		EventDate: EventDateresults,
	}

	return helper.CreateResult(true, data, "success")
}

func (m *DataBrowserController) GetCustomAvailDate(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	Dateresults := make([]time.Time, 0)

	// ScadaDataOEM
	for i := 0; i < 2; i++ {
		var arrsort []string
		if i == 0 {
			arrsort = append(arrsort, "timestamp")
		} else {
			arrsort = append(arrsort, "-timestamp")
		}

		query := DB().Connection.NewQuery().From(new(ScadaDataOEM).TableName()).Skip(0).Take(1)
		query = query.Order(arrsort...)

		queryMetTower := DB().Connection.NewQuery().From(new(MetTower).TableName()).Skip(0).Take(1)
		queryMetTower = queryMetTower.Order(arrsort...)

		csr, e := query.Cursor(nil)
		csrM, eM := queryMetTower.Cursor(nil)
		if e != nil {
			return helper.CreateResult(false, nil, e.Error())
		}
		defer csr.Close()

		if eM != nil {
			return helper.CreateResult(false, nil, e.Error())
		}
		defer csrM.Close()

		Result := make([]ScadaDataOEM, 0)
		e = csr.Fetch(&Result, 0, false)

		ResultMetTower := make([]MetTower, 0)
		eM = csrM.Fetch(&ResultMetTower, 0, false)

		for _, val := range Result {
			Dateresults = append(Dateresults, val.TimeStamp.UTC())
		}
		for _, val := range ResultMetTower {
			Dateresults = append(Dateresults, val.TimeStamp.UTC())
		}
	}

	data := struct {
		CustomDate []time.Time
	}{
		CustomDate: Dateresults,
	}

	return helper.CreateResult(true, data, "success")
}

func (m *DataBrowserController) GetScadaDataOEMAvailDate(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	ScadaDataOEMresults := make([]time.Time, 0)

	// ScadaDataOEM Data
	for i := 0; i < 2; i++ {
		var arrsort []string
		if i == 0 {
			arrsort = append(arrsort, "timestamp")
		} else {
			arrsort = append(arrsort, "-timestamp")
		}

		query := DB().Connection.NewQuery().From(new(ScadaDataOEM).TableName()).Skip(0).Take(1)
		query = query.Order(arrsort...)

		csr, e := query.Cursor(nil)
		if e != nil {
			return helper.CreateResult(false, nil, e.Error())
		}
		defer csr.Close()

		Result := make([]ScadaDataOEM, 0)
		e = csr.Fetch(&Result, 0, false)

		if e != nil {
			return helper.CreateResult(false, nil, e.Error())
		}

		for _, val := range Result {
			ScadaDataOEMresults = append(ScadaDataOEMresults, val.TimeStamp.UTC())
		}
	}

	data := struct {
		ScadaDataOEM []time.Time
	}{
		ScadaDataOEM: ScadaDataOEMresults,
	}

	return helper.CreateResult(true, data, "success")
}

func (m *DataBrowserController) GetDowntimeEventvailDate(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	DowntimeEventresults := make([]time.Time, 0)

	// Downtime Event Data
	for i := 0; i < 2; i++ {
		var arrsort []string
		if i == 0 {
			arrsort = append(arrsort, "timestart")
		} else {
			arrsort = append(arrsort, "-timestart")
		}

		query := DB().Connection.NewQuery().From(new(EventDown).TableName()).Skip(0).Take(1)
		query = query.Order(arrsort...)

		csr, e := query.Cursor(nil)
		if e != nil {
			return helper.CreateResult(false, nil, e.Error())
		}
		defer csr.Close()

		Result := make([]EventDown, 0)
		e = csr.Fetch(&Result, 0, false)

		if e != nil {
			return helper.CreateResult(false, nil, e.Error())
		}

		for _, val := range Result {
			DowntimeEventresults = append(DowntimeEventresults, val.TimeStart.UTC())
		}
	}

	data := struct {
		DowntimeEvent []time.Time
	}{
		DowntimeEvent: DowntimeEventresults,
	}

	return helper.CreateResult(true, data, "success")
}

func (m *DataBrowserController) GetScadaDataHFDAvailDate(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	ScadaDataHFDresults := make([]time.Time, 0)

	// ScadaDataHFD Data
	for i := 0; i < 2; i++ {
		var arrsort []string
		if i == 0 {
			arrsort = append(arrsort, "timestamp")
		} else {
			arrsort = append(arrsort, "-timestamp")
		}

		query := DB().Connection.NewQuery().From(new(ScadaDataHFD).TableName()).Skip(0).Take(1)
		query = query.Order(arrsort...)

		csr, e := query.Cursor(nil)
		if e != nil {
			return helper.CreateResult(false, nil, e.Error())
		}
		defer csr.Close()

		Result := make([]ScadaDataHFD, 0)
		e = csr.Fetch(&Result, 0, false)

		if e != nil {
			return helper.CreateResult(false, nil, e.Error())
		}

		for _, val := range Result {
			ScadaDataHFDresults = append(ScadaDataHFDresults, val.TimeStamp.UTC())
		}
	}

	data := struct {
		ScadaDataHFD []time.Time
	}{
		ScadaDataHFD: ScadaDataHFDresults,
	}

	return helper.CreateResult(true, data, "success")
}

// Generate excel

func (m *DataBrowserController) GenExcelCustom10Minutes(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	p := new(helper.PayloadsDB)
	e := k.GetPayload(&p)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	tStart, _ := time.Parse("2006-01-02", p.DateStart.UTC().Format("2006-01-02"))
	tEnd, _ := time.Parse("2006-01-02 15:04:05", p.DateEnd.UTC().Format("2006-01-02")+" 23:59:59")
	turbine := p.Turbine

	istimestamp := false
	arrscadaoem := []string{"_id", "timestamputc"}
	arrmettower := []string{}
	headerList := []string{}
	fieldList := []string{}
	if p.Custom.Has("ColumnList") {
		for _, _val := range p.Custom["ColumnList"].([]interface{}) {
			_tkm, _ := tk.ToM(_val)
			if _tkm.GetString("source") == "ScadaDataOEM" {
				arrscadaoem = append(arrscadaoem, _tkm.GetString("_id"))
				if _tkm.GetString("_id") == "timestamp" {
					istimestamp = true
				}
			} else if _tkm.GetString("source") == "MetTower" {
				arrmettower = append(arrmettower, _tkm.GetString("_id"))
			}
			headerList = append(headerList, _tkm.GetString("label"))
			fieldList = append(fieldList, _tkm.GetString("_id"))
		}
	}
	var filter []*dbox.Filter
	filter = append(filter, dbox.Gte("timestamp", tStart))
	filter = append(filter, dbox.Lte("timestamp", tEnd))
	if len(turbine) > 0 {
		filter = append(filter, dbox.In("turbine", turbine...))
	}
	if p.Project != "" {
		filter = append(filter, dbox.Eq("projectname", p.Project))
	}

	csr, e := DB().Connection.NewQuery().Select(arrscadaoem...).
		From(new(ScadaDataOEM).TableName()).Where(dbox.And(filter...)).Cursor(nil)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}
	defer csr.Close()

	results := make([]tk.M, 0)
	e = csr.Fetch(&results, 0, false)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	arrmettowercond := []interface{}{}

	for i, val := range results {
		if val.Has("timestamputc") {
			itime := val.Get("timestamputc", time.Time{}).(time.Time).UTC()
			arrmettowercond = append(arrmettowercond, itime)
			val.Set("timestamputc", itime)
			results[i] = val
		}
		if istimestamp {
			itime := val.Get("timestamp", time.Time{}).(time.Time)
			val.Set("timestamp", itime.UTC())
			results[i] = val
		}
	}

	tkmmet := tk.M{}
	if len(arrmettower) > 0 && len(arrmettowercond) > 0 {
		arrmettower = append(arrmettower, "timestamp")
		csrMet, _e := DB().Connection.NewQuery().Select(arrmettower...).From("MetTower").
			Where(dbox.In("timestamp", arrmettowercond...)).Cursor(nil)
		if _e != nil {
			return helper.CreateResult(false, nil, _e.Error())
		}
		defer csrMet.Close()

		_resmet := make([]tk.M, 0)
		_e = csrMet.Fetch(&_resmet, 0, false)

		if _e != nil {
			return helper.CreateResult(false, nil, _e.Error())
		}

		for _, val := range _resmet {
			itime := val.Get("timestamp", time.Time{}).(time.Time).UTC().String()
			tkmmet.Set(itime, val)
		}
	}

	if len(tkmmet) > 0 {
		for i, val := range results {
			itime := val.Get("timestamputc", time.Time{}).(time.Time).UTC().String()
			if tkmmet.Has(itime) {
				for _key, _val := range tkmmet[itime].(tk.M) {
					if _key != "timestamp" {
						val.Set(_key, _val)
					}
				}
			}
			val.Unset("timestamputc")
			results[i] = val
		}
	}

	var pathDownload string
	typeExcel := "Custom10Minutes"
	TimeCreate := time.Now().Format("2006-01-02_150405")
	CreateDateTime := typeExcel + TimeCreate

	if err := os.RemoveAll("web/assets/Excel/" + typeExcel + "/"); err != nil {
		tk.Println(err)
	}

	if _, err := os.Stat("web/assets/Excel/" + typeExcel + "/"); os.IsNotExist(err) {
		os.MkdirAll("web/assets/Excel/"+typeExcel+"/", 0777)
	}

	DeserializeCustom10Minutes(results, typeExcel, CreateDateTime, headerList, fieldList)
	pathDownload = "res/Excel/" + typeExcel + "/" + CreateDateTime + ".xlsx"

	return helper.CreateResult(true, pathDownload, "success")
}

func (m *DataBrowserController) GenExcelScadaOem(k *knot.WebContext) interface{} {

	k.Config.OutputType = knot.OutputJson

	var filter []*dbox.Filter

	p := new(helper.PayloadsDB)
	e := k.GetPayload(&p)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	tStart, _ := time.Parse("2006-01-02", p.DateStart.UTC().Format("2006-01-02"))
	tEnd, _ := time.Parse("2006-01-02 15:04:05", p.DateEnd.UTC().Format("2006-01-02")+" 23:59:59")
	turbine := p.Turbine

	var pathDownload string
	typeExcel := "ScadaOem"
	TimeCreate := time.Now().Format("2006-01-02_150405")
	CreateDateTime := typeExcel + TimeCreate

	if err := os.RemoveAll("web/assets/Excel/" + typeExcel + "/"); err != nil {
		tk.Println(err)
	}

	filter = append(filter, dbox.Ne("_id", ""))
	// filter = append(filter, dbox.Ne("powerlost", ""))
	// filter = append(filter, dbox.Ne("ai_intern_activpower", ""))
	// filter = append(filter, dbox.Ne("ai_intern_windspeed", ""))
	filter = append(filter, dbox.Gte("timestamp", tStart))
	filter = append(filter, dbox.Lte("timestamp", tEnd))
	if len(turbine) != 0 {
		filter = append(filter, dbox.In("turbine", turbine...))
	}
	if p.Project != "" {
		filter = append(filter, dbox.Eq("projectname", p.Project))
	}

	query := DB().Connection.NewQuery().From(new(ScadaDataOEM).TableName()).Where(dbox.And(filter...))

	csr, e := query.Cursor(nil)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}
	defer csr.Close()

	tmpResult := make([]ScadaDataOEM, 0)
	results := make([]ScadaDataOEM, 0)
	e = csr.Fetch(&tmpResult, 0, false)

	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	for _, val := range tmpResult {
		val.TimeStamp = val.TimeStamp.UTC()
		results = append(results, val)
	}
	//web/assets/Excel/

	if _, err := os.Stat("web/assets/Excel/" + typeExcel + "/"); os.IsNotExist(err) {
		os.MkdirAll("web/assets/Excel/"+typeExcel+"/", 0777)
	}
	turbineName, err := helper.GetTurbineNameList(p.Project)
	if err != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	DeserializeScadaOem(results, 0, typeExcel, CreateDateTime, turbineName)
	pathDownload = "res/Excel/" + typeExcel + "/" + CreateDateTime + ".xlsx"
	// tk.Println(pathDownload)

	return helper.CreateResult(true, pathDownload, "success")
}

func (m *DataBrowserController) GenExcelDowntimeEvent(k *knot.WebContext) interface{} {

	k.Config.OutputType = knot.OutputJson

	var filter []*dbox.Filter

	p := new(helper.PayloadsDB)
	e := k.GetPayload(&p)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	tStart, _ := time.Parse("2006-01-02", p.DateStart.UTC().Format("2006-01-02"))
	tEnd, _ := time.Parse("2006-01-02 15:04:05", p.DateEnd.UTC().Format("2006-01-02")+" 23:59:59")
	turbine := p.Turbine

	var pathDownload string
	typeExcel := "DowntimeEvent"
	TimeCreate := time.Now().Format("2006-01-02_150405")
	CreateDateTime := typeExcel + TimeCreate

	if err := os.RemoveAll("web/assets/Excel/" + typeExcel + "/"); err != nil {
		tk.Println(err)
	}

	filter = append(filter, dbox.Ne("_id", ""))
	filter = append(filter, dbox.Gte("timestart", tStart))
	filter = append(filter, dbox.Lte("timestart", tEnd))
	if len(turbine) != 0 {
		filter = append(filter, dbox.In("turbine", turbine...))
	}
	if p.Project != "" {
		filter = append(filter, dbox.Eq("projectname", p.Project))
	}

	query := DB().Connection.NewQuery().From(new(EventDown).TableName()).Where(dbox.And(filter...))

	csr, e := query.Cursor(nil)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}
	defer csr.Close()

	tmpResult := make([]EventDown, 0)
	// results := make([]EventDown, 0)
	e = csr.Fetch(&tmpResult, 0, false)

	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	// for _, val := range tmpResult {
	// 	val.TimeStart = val.TimeStart.UTC()
	// 	results = append(results, val)
	// }
	//web/assets/Excel/

	if _, err := os.Stat("web/assets/Excel/" + typeExcel + "/"); os.IsNotExist(err) {
		os.MkdirAll("web/assets/Excel/"+typeExcel+"/", 0777)
	}

	turbineName, err := helper.GetTurbineNameList(p.Project)
	if err != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	DeserializeEventDown(tmpResult, 0, typeExcel, CreateDateTime, turbineName)
	pathDownload = "res/Excel/" + typeExcel + "/" + CreateDateTime + ".xlsx"
	// tk.Println(pathDownload)

	return helper.CreateResult(true, pathDownload, "success")
}

func (m *DataBrowserController) GenExcelEventRaw(k *knot.WebContext) interface{} {

	k.Config.OutputType = knot.OutputJson

	var filter []*dbox.Filter

	p := new(helper.PayloadsDB)
	e := k.GetPayload(&p)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	tStart, _ := time.Parse("2006-01-02", p.DateStart.UTC().Format("2006-01-02"))
	tEnd, _ := time.Parse("2006-01-02 15:04:05", p.DateEnd.UTC().Format("2006-01-02")+" 23:59:59")
	turbine := p.Turbine

	var pathDownload string
	typeExcel := "EventRaw"
	TimeCreate := time.Now().Format("2006-01-02_150405")
	CreateDateTime := typeExcel + TimeCreate

	if err := os.RemoveAll("web/assets/Excel/" + typeExcel + "/"); err != nil {
		tk.Println(err)
	}

	filter = append(filter, dbox.Ne("_id", ""))
	filter = append(filter, dbox.Gte("timestamp", tStart))
	filter = append(filter, dbox.Lte("timestamp", tEnd))
	if len(turbine) != 0 {
		filter = append(filter, dbox.In("turbine", turbine...))
	}
	if p.Project != "" {
		filter = append(filter, dbox.Eq("projectname", p.Project))
	}

	query := DB().Connection.NewQuery().From(new(EventRaw).TableName()).Where(dbox.And(filter...))

	csr, e := query.Cursor(nil)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}
	defer csr.Close()

	tmpResult := make([]EventRaw, 0)
	e = csr.Fetch(&tmpResult, 0, false)

	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	if _, err := os.Stat("web/assets/Excel/" + typeExcel + "/"); os.IsNotExist(err) {
		os.MkdirAll("web/assets/Excel/"+typeExcel+"/", 0777)
	}

	turbineName, err := helper.GetTurbineNameList(p.Project)
	if err != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	DeserializeEventRaw(tmpResult, 0, typeExcel, CreateDateTime, turbineName)
	pathDownload = "res/Excel/" + typeExcel + "/" + CreateDateTime + ".xlsx"

	return helper.CreateResult(true, pathDownload, "success")
}

func (m *DataBrowserController) GenExcelMet(k *knot.WebContext) interface{} {

	k.Config.OutputType = knot.OutputJson

	var filter []*dbox.Filter

	p := new(helper.PayloadsDB)
	e := k.GetPayload(&p)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	tStart, _ := time.Parse("2006-01-02", p.DateStart.UTC().Format("2006-01-02"))
	tEnd, _ := time.Parse("2006-01-02 15:04:05", p.DateEnd.UTC().Format("2006-01-02")+" 23:59:59")
	// turbine := p.Turbine

	var pathDownload string
	typeExcel := "MetTower"
	TimeCreate := time.Now().Format("2006-01-02_150405")
	CreateDateTime := typeExcel + TimeCreate

	if err := os.RemoveAll("web/assets/Excel/" + typeExcel + "/"); err != nil {
		tk.Println(err)
	}

	filter = append(filter, dbox.Ne("_id", ""))
	filter = append(filter, dbox.Gte("timestamp", tStart))
	filter = append(filter, dbox.Lte("timestamp", tEnd))
	if p.Project != "" {
		filter = append(filter, dbox.Eq("project", p.Project))
	}

	query := DB().Connection.NewQuery().From(new(MetTower).TableName()).Where(dbox.And(filter...))

	csr, e := query.Cursor(nil)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}
	defer csr.Close()

	tmpResult := make([]MetTower, 0)
	e = csr.Fetch(&tmpResult, 0, false)

	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	if _, err := os.Stat("web/assets/Excel/" + typeExcel + "/"); os.IsNotExist(err) {
		os.MkdirAll("web/assets/Excel/"+typeExcel+"/", 0777)
	}

	DeserializeMetTower(tmpResult, 0, typeExcel, CreateDateTime)
	pathDownload = "res/Excel/" + typeExcel + "/" + CreateDateTime + ".xlsx"

	return helper.CreateResult(true, pathDownload, "success")
}

func (m *DataBrowserController) GenExcelScada(k *knot.WebContext) interface{} {

	k.Config.OutputType = knot.OutputJson

	var filter []*dbox.Filter

	p := new(helper.PayloadsDB)
	e := k.GetPayload(&p)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	tStart, _ := time.Parse("2006-01-02", p.DateStart.UTC().Format("2006-01-02"))
	tEnd, _ := time.Parse("2006-01-02 15:04:05", p.DateEnd.UTC().Format("2006-01-02")+" 23:59:59")
	turbine := p.Turbine

	var pathDownload string
	typeExcel := "ScadaData"
	TimeCreate := time.Now().Format("2006-01-02_150405")
	CreateDateTime := typeExcel + TimeCreate

	if err := os.RemoveAll("web/assets/Excel/" + typeExcel + "/"); err != nil {
		tk.Println(err)
	}

	filter = append(filter, dbox.Ne("_id", ""))
	filter = append(filter, dbox.Gte("timestamp", tStart))
	filter = append(filter, dbox.Lte("timestamp", tEnd))
	if len(turbine) != 0 {
		filter = append(filter, dbox.In("turbine", turbine...))
	}
	if p.Project != "" {
		filter = append(filter, dbox.Eq("projectname", p.Project))
	}

	query := DB().Connection.NewQuery().From(new(ScadaData).TableName()).Where(dbox.And(filter...))

	csr, e := query.Cursor(nil)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}
	defer csr.Close()

	tmpResult := make([]ScadaData, 0)
	e = csr.Fetch(&tmpResult, 0, false)

	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	if _, err := os.Stat("web/assets/Excel/" + typeExcel + "/"); os.IsNotExist(err) {
		os.MkdirAll("web/assets/Excel/"+typeExcel+"/", 0777)
	}
	turbineName, err := helper.GetTurbineNameList(p.Project)
	if err != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	DeserializeScadaData(tmpResult, 0, typeExcel, CreateDateTime, turbineName)
	pathDownload = "res/Excel/" + typeExcel + "/" + CreateDateTime + ".xlsx"

	return helper.CreateResult(true, pathDownload, "success")
}

func (m *DataBrowserController) GenExcelScadaHFD(k *knot.WebContext) interface{} {

	k.Config.OutputType = knot.OutputJson

	var filter []*dbox.Filter

	p := new(helper.PayloadsDB)
	e := k.GetPayload(&p)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	tStart, _ := time.Parse("2006-01-02", p.DateStart.UTC().Format("2006-01-02"))
	tEnd, _ := time.Parse("2006-01-02 15:04:05", p.DateEnd.UTC().Format("2006-01-02")+" 23:59:59")
	turbine := p.Turbine
	project := p.Project

	var pathDownload string
	typeExcel := "ScadaDataHFD"
	TimeCreate := time.Now().Format("2006-01-02_150405")
	CreateDateTime := typeExcel + TimeCreate

	if err := os.RemoveAll("web/assets/Excel/" + typeExcel + "/"); err != nil {
		tk.Println(err)
	}

	filter = append(filter, dbox.Ne("_id", ""))
	filter = append(filter, dbox.Gte("timestamp", tStart))
	filter = append(filter, dbox.Lte("timestamp", tEnd))
	if len(turbine) != 0 {
		filter = append(filter, dbox.In("turbine", turbine...))
	}
	if project != "" {
		filter = append(filter, dbox.Eq("projectname", project))
	}

	query := DB().Connection.NewQuery().From(new(ScadaDataHFD).TableName()).Where(dbox.And(filter...))

	csr, e := query.Cursor(nil)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}
	defer csr.Close()

	tmpResult := make([]ScadaDataHFD, 0)
	e = csr.Fetch(&tmpResult, 0, false)

	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	if _, err := os.Stat("web/assets/Excel/" + typeExcel + "/"); os.IsNotExist(err) {
		os.MkdirAll("web/assets/Excel/"+typeExcel+"/", 0777)
	}

	turbineName, err := helper.GetTurbineNameList(project)
	if err != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	DeserializeScadaDataHFD(tmpResult, 0, typeExcel, CreateDateTime, turbineName)
	pathDownload = "res/Excel/" + typeExcel + "/" + CreateDateTime + ".xlsx"

	return helper.CreateResult(true, pathDownload, "success")
}

// Deserialize

func DeserializeCustom10Minutes(data []tk.M, typeExcel string, CreateDateTime string, header []string, fieldList []string) error {
	filename := ""
	filename = "web/assets/Excel/" + typeExcel + "/" + CreateDateTime + ".xlsx"

	file := x.NewFile()
	sheet, _ := file.AddSheet("Sheet1")

	for i, each := range data {
		if i == 0 {
			rowHeader := sheet.AddRow()
			for _, hdr := range header {

				cell := rowHeader.AddCell()
				cell.Value = hdr
			}
		}

		rowContent := sheet.AddRow()
		cell := rowContent.AddCell()
		for idx, field := range fieldList {
			if idx > 0 {
				cell = rowContent.AddCell()
			}
			switch field {
			case "timestamp", "timestamputc":
				cell.Value = each[field].(time.Time).Format("2006-01-02 15:04:05")
			case "turbine":
				cell.Value = each.GetString(field)
			default:
				cell.Value = tk.ToString(each.GetFloat64(field))
			}
		}
	}

	tk.Println(filename)

	err := file.Save(filename)
	if err != nil {
		return err
	}

	return nil
}

func DeserializeScadaOem(data []ScadaDataOEM, j int, typeExcel string, CreateDateTime string, turbinename map[string]string) error {
	//savecipo += 1
	filename := ""
	filename = "web/assets/Excel/" + typeExcel + "/" + CreateDateTime + ".xlsx"

	file := x.NewFile()
	sheet, _ := file.AddSheet("Sheet1")
	header := []string{"TimeStamp", "Turbine", "AI intern R PidAngleOut", "AI intern ActivPower ", "AI intern I1 ", "AI intern I2", "AI intern I3", "AI intern NacelleDrill ", "AI intern NacellePos ", "AI intern PitchAkku V1 ", "AI intern PitchAkku V2 ", "AI intern PitchAkku V3 ", "AI intern PitchAngle1 ", "AI intern PitchAngle2 ", "AI intern PitchAngle3 ", "AI intern PitchConv Current1 ", "AI intern PitchConv Current2 ", "AI intern PitchConv Current3 ", "AI intern PitchAngleSP Diff1 ", "AI intern PitchAngleSP Diff2 ", "AI intern PitchAngleSP Diff3 ", "AI intern ReactivPower ", "AI intern RpmDiff ", "AI intern U1 ", "AI intern U2 ", "AI intern U3 ", "AI intern WindDirection ", "AI intern WindSpeed ", "AI Intern WindSpeedDif ", "AI speed RotFR ", "AI WindSpeed1 ", "AI WindSpeed2 ", "AI WindVane1 ", "AI WindVane2 ", "AI internCurrentAsym ", "Temp GearBox IMS NDE ", "AI intern WindVaneDiff ", "C intern SpeedGenerator ", "C intern SpeedRotor ", "AI intern Speed RPMDiff FR1 RotCNT ", "AI intern Frequency Grid ", "Temp GearBox HSS NDE ", "AI DrTrVibValue ", "AI intern InLastErrorConv1 ", "AI intern InLastErrorConv2 ", "AI intern InLastErrorConv3 ", "AI intern TempConv1 ", "AI intern TempConv2 ", "AI intern TempConv3 ", "AI intern PitchSpeed2", "Temp YawBrake 1 ", "Temp YawBrake 2 ", "Temp G1L1 ", "Temp G1L2 ", "Temp G1L3 ", "Temp YawBrake 4", "AI HydrSystemPressure ", "Temp BottomControlSection Low ", "Temp GearBox HSS DE ", "Temp GearOilSump ", "Temp GeneratorBearing DE ", "Temp GeneratorBearing NDE ", "Temp MainBearing ", "Temp GearBox IMS DE ", "Temp Nacelle ", "Temp Outdoor ", "AI TowerVibValueAxial ", "AI intern DiffGenSpeedSPToAct ", "Temp YawBrake 5", "AI intern SpeedGenerator Proximity ", "AI intern SpeedDiff Encoder Proximity ", "AI GearOilPressure ", "Temp CabinetTopBox Low ", "Temp CabinetTopBox ", "Temp BottomControlSection ", "Temp BottomPowerSection ", "Temp BottomPowerSection Low ", "AI intern Pitch1 Status High ", "AI intern Pitch2 Status High ", "AI intern Pitch3 Status High ", "AI intern InPosition1 ch3", "AI intern InPosition2 ch3", "AI intern InPosition3 ch3", "AI intern Temp Brake Blade1 ", "AI intern Temp Brake Blade2 ", "AI intern Temp Brake Blade3 ", "AI intern Temp PitchMotor Blade1 ", "AI intern Temp PitchMotor Blade2 ", "AI intern Temp PitchMotor Blade3 ", "AI intern Temp Hub Additional1 ", "AI intern Temp Hub Additional2 ", "AI intern Temp Hub Additional3 ", "AI intern Pitch1 Status Low ", "AI intern Pitch2 Status Low ", "AI intern Pitch3 Status Low ", "AI intern Battery VoltageBlade1 center ", "AI intern Battery VoltageBlade2 center ", "AI intern Battery VoltageBlade3 center ", "AI intern Battery ChargingCur Blade1 ", "AI intern Battery ChargingCur Blade2 ", "AI intern Battery ChargingCur Blade3 ", "AI intern Battery DischargingCur Blade1 ", "AI intern Battery DischargingCur Blade2 ", "AI intern Battery DischargingCur Blade3 ", "AI intern PitchMotor BrakeVoltage Blade1 ", "AI intern PitchMotor BrakeVoltage Blade2 ", "AI intern PitchMotor BrakeVoltage Blade3 ", "AI intern PitchMotor BrakeCurrent Blade1 ", "AI intern PitchMotor BrakeCurrent Blade2 ", "AI intern PitchMotor BrakeCurrent Blade3 ", "AI intern Temp HubBox Blade1 ", "AI intern Temp HubBox Blade2 ", "AI intern Temp HubBox Blade3 ", "AI intern Temp Pitch1 HeatSink ", "AI intern Temp Pitch2 HeatSink ", "AI intern Temp Pitch3 HeatSink ", "AI intern ErrorStackBlade1 ", "AI intern ErrorStackBlade2 ", "AI intern ErrorStackBlade3 ", "AI intern Temp BatteryBox Blade1 ", "AI intern Temp BatteryBox Blade2 ", "AI intern Temp BatteryBox Blade3 ", "AI intern DC LinkVoltage1 ", "AI intern DC LinkVoltage2 ", "AI intern DC LinkVoltage3 ", "Temp Yaw Motor1 ", "Temp Yaw Motor2 ", "Temp Yaw Motor3 ", "Temp Yaw Motor4 ", "AO DFIG Power Setpiont ", "AO DFIG Q Setpoint ", "AI DFIG Torque actual ", "AI DFIG SpeedGenerator Encoder ", "AI intern DFIG DC Link Voltage actual ", "AI intern DFIG MSC current ", "AI intern DFIG Main voltage ", "AI intern DFIG Main current ", "AI intern DFIG active power actual ", "AI intern DFIG reactive power actual ", "AI intern DFIG active power actual LSC ", "AI intern DFIG LSC current ", "AI intern DFIG Data log number ", "AI intern Damper OscMagnitude ", "AI intern Damper PassbandFullLoad ", "AI YawBrake TempRise1 ", "AI YawBrake TempRise2 ", "AI YawBrake TempRise3 ", "AI YawBrake TempRise4 ", "AI intern NacelleDrill at NorthPosSensor "}

	for i, each := range data {
		if i == 0 {
			rowHeader := sheet.AddRow()
			for _, hdr := range header {

				cell := rowHeader.AddCell()
				cell.Value = hdr
			}
		}

		rowContent := sheet.AddRow()

		cell := rowContent.AddCell()
		cell.Value = each.TimeStamp.Format("2006-01-02 15:04:05")

		cell = rowContent.AddCell()
		cell.Value = turbinename[each.Turbine]

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.AI_intern_R_PidAngleOut, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.AI_intern_ActivPower, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.AI_intern_I1, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.AI_intern_I2, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.AI_intern_I3, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.AI_intern_NacelleDrill, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.AI_intern_NacellePos, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.AI_intern_PitchAkku_V1, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.AI_intern_PitchAkku_V2, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.AI_intern_PitchAkku_V3, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.AI_intern_PitchAngle1, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.AI_intern_PitchAngle2, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.AI_intern_PitchAngle3, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.AI_intern_PitchConv_Current1, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.AI_intern_PitchConv_Current2, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.AI_intern_PitchConv_Current3, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.AI_intern_PitchAngleSP_Diff1, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.AI_intern_PitchAngleSP_Diff2, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.AI_intern_PitchAngleSP_Diff3, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.AI_intern_ReactivPower, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.AI_intern_RpmDiff, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.AI_intern_U1, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.AI_intern_U2, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.AI_intern_U3, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.AI_intern_WindDirection, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.AI_intern_WindSpeed, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.AI_Intern_WindSpeedDif, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.AI_speed_RotFR, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.AI_WindSpeed1, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.AI_WindSpeed2, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.AI_WindVane1, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.AI_WindVane2, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.AI_internCurrentAsym, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Temp_GearBox_IMS_NDE, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.AI_intern_WindVaneDiff, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.C_intern_SpeedGenerator, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.C_intern_SpeedRotor, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.AI_intern_Speed_RPMDiff_FR1_RotCNT, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.AI_intern_Frequency_Grid, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Temp_GearBox_HSS_NDE, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.AI_DrTrVibValue, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.AI_intern_InLastErrorConv1, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.AI_intern_InLastErrorConv2, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.AI_intern_InLastErrorConv3, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.AI_intern_TempConv1, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.AI_intern_TempConv2, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.AI_intern_TempConv3, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.AI_intern_PitchSpeed1, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Temp_YawBrake_1, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Temp_YawBrake_2, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Temp_G1L1, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Temp_G1L2, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Temp_G1L3, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Temp_YawBrake_3, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.AI_HydrSystemPressure, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Temp_BottomControlSection_Low, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Temp_GearBox_HSS_DE, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Temp_GearOilSump, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Temp_GeneratorBearing_DE, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Temp_GeneratorBearing_NDE, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Temp_MainBearing, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Temp_GearBox_IMS_DE, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Temp_Nacelle, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Temp_Outdoor, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.AI_TowerVibValueAxial, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.AI_intern_DiffGenSpeedSPToAct, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Temp_YawBrake_4, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.AI_intern_SpeedGenerator_Proximity, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.AI_intern_SpeedDiff_Encoder_Proximity, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.AI_GearOilPressure, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Temp_CabinetTopBox_Low, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Temp_CabinetTopBox, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Temp_BottomControlSection, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Temp_BottomPowerSection, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Temp_BottomPowerSection_Low, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.AI_intern_Pitch1_Status_High, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.AI_intern_Pitch2_Status_High, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.AI_intern_Pitch3_Status_High, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.AI_intern_InPosition1_ch2, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.AI_intern_InPosition2_ch2, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.AI_intern_InPosition3_ch2, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.AI_intern_Temp_Brake_Blade1, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.AI_intern_Temp_Brake_Blade2, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.AI_intern_Temp_Brake_Blade3, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.AI_intern_Temp_PitchMotor_Blade1, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.AI_intern_Temp_PitchMotor_Blade2, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.AI_intern_Temp_PitchMotor_Blade3, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.AI_intern_Temp_Hub_Additional1, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.AI_intern_Temp_Hub_Additional2, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.AI_intern_Temp_Hub_Additional3, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.AI_intern_Pitch1_Status_Low, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.AI_intern_Pitch2_Status_Low, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.AI_intern_Pitch3_Status_Low, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.AI_intern_Battery_VoltageBlade1_center, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.AI_intern_Battery_VoltageBlade2_center, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.AI_intern_Battery_VoltageBlade3_center, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.AI_intern_Battery_ChargingCur_Blade1, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.AI_intern_Battery_ChargingCur_Blade2, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.AI_intern_Battery_ChargingCur_Blade3, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.AI_intern_Battery_DischargingCur_Blade1, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.AI_intern_Battery_DischargingCur_Blade2, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.AI_intern_Battery_DischargingCur_Blade3, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.AI_intern_PitchMotor_BrakeVoltage_Blade1, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.AI_intern_PitchMotor_BrakeVoltage_Blade2, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.AI_intern_PitchMotor_BrakeVoltage_Blade3, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.AI_intern_PitchMotor_BrakeCurrent_Blade1, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.AI_intern_PitchMotor_BrakeCurrent_Blade2, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.AI_intern_PitchMotor_BrakeCurrent_Blade3, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.AI_intern_Temp_HubBox_Blade1, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.AI_intern_Temp_HubBox_Blade2, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.AI_intern_Temp_HubBox_Blade3, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.AI_intern_Temp_Pitch1_HeatSink, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.AI_intern_Temp_Pitch2_HeatSink, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.AI_intern_Temp_Pitch3_HeatSink, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.AI_intern_ErrorStackBlade1, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.AI_intern_ErrorStackBlade2, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.AI_intern_ErrorStackBlade3, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.AI_intern_Temp_BatteryBox_Blade1, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.AI_intern_Temp_BatteryBox_Blade2, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.AI_intern_Temp_BatteryBox_Blade3, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.AI_intern_DC_LinkVoltage1, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.AI_intern_DC_LinkVoltage2, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.AI_intern_DC_LinkVoltage3, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Temp_Yaw_Motor1, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Temp_Yaw_Motor2, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Temp_Yaw_Motor3, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Temp_Yaw_Motor4, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.AO_DFIG_Power_Setpiont, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.AO_DFIG_Q_Setpoint, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.AI_DFIG_Torque_actual, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.AI_DFIG_SpeedGenerator_Encoder, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.AI_intern_DFIG_DC_Link_Voltage_actual, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.AI_intern_DFIG_MSC_current, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.AI_intern_DFIG_Main_voltage, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.AI_intern_DFIG_Main_current, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.AI_intern_DFIG_active_power_actual, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.AI_intern_DFIG_reactive_power_actual, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.AI_intern_DFIG_active_power_actual_LSC, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.AI_intern_DFIG_LSC_current, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.AI_intern_DFIG_Data_log_number, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.AI_intern_Damper_OscMagnitude, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.AI_intern_Damper_PassbandFullLoad, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.AI_YawBrake_TempRise1, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.AI_YawBrake_TempRise2, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.AI_YawBrake_TempRise3, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.AI_YawBrake_TempRise4, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.AI_intern_NacelleDrill_at_NorthPosSensor, 'f', -1, 64)

	}

	tk.Println(filename)

	err := file.Save(filename)
	if err != nil {
		return err
	}

	return nil
}

func DeserializeEventDown(data []EventDown, j int, typeExcel string, CreateDateTime string, turbinename map[string]string) error {
	//savecipo += 1
	filename := ""
	filename = "web/assets/Excel/" + typeExcel + "/" + CreateDateTime + ".xlsx"

	file := x.NewFile()
	sheet, _ := file.AddSheet("Sheet1")
	header := []string{"Turbine", "TimeStart", "TimeEnd", "Down Grid", "Down Environment", "Down Machine", "Alarm Description", "Duration"}

	for i, each := range data {
		if i == 0 {
			rowHeader := sheet.AddRow()
			for _, hdr := range header {

				cell := rowHeader.AddCell()
				cell.Value = hdr
			}
		}

		rowContent := sheet.AddRow()

		cell := rowContent.AddCell()
		cell.Value = turbinename[each.Turbine]

		cell = rowContent.AddCell()
		cell.Value = each.TimeStart.Format("2006-01-02 15:04:05")

		cell = rowContent.AddCell()
		cell.Value = each.TimeEnd.Format("2006-01-02 15:04:05")

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatBool(each.DownGrid)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatBool(each.DownEnvironment)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatBool(each.DownMachine)

		cell = rowContent.AddCell()
		cell.Value = each.AlarmDescription //strconv.FormatFloat(each.AI_intern_R_PidAngleOut , 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = SecondsToHms(each.Duration) //strconv.FormatFloat(each.Duration, 'f', -1, 64)

	}

	err := file.Save(filename)
	if err != nil {
		return err
	}

	return nil
}

func DeserializeEventRaw(data []EventRaw, j int, typeExcel string, CreateDateTime string, turbinename map[string]string) error {
	//savecipo += 1
	filename := ""
	filename = "web/assets/Excel/" + typeExcel + "/" + CreateDateTime + ".xlsx"

	file := x.NewFile()
	sheet, _ := file.AddSheet("Sheet1")
	header := []string{"TimeStamp", "Project Name", "Turbine", "Event Type", "Alarm Description", "Turbine Status", "Brake Type", "Brake Program", "Alarm Id", "Alarm Toggle"}

	for i, each := range data {
		if i == 0 {
			rowHeader := sheet.AddRow()
			for _, hdr := range header {

				cell := rowHeader.AddCell()
				cell.Value = hdr
			}
		}

		rowContent := sheet.AddRow()

		cell := rowContent.AddCell()
		cell.Value = each.TimeStamp.Format("2006-01-02 15:04:05")

		cell = rowContent.AddCell()
		cell.Value = each.ProjectName

		cell = rowContent.AddCell()
		cell.Value = turbinename[each.Turbine]

		cell = rowContent.AddCell()
		cell.Value = each.EventType

		cell = rowContent.AddCell()
		cell.Value = each.AlarmDescription

		cell = rowContent.AddCell()
		cell.Value = each.TurbineStatus

		cell = rowContent.AddCell()
		cell.Value = each.BrakeType

		cell = rowContent.AddCell()
		cell.Value = strconv.Itoa(each.BrakeProgram) //strconv.Formatint(each.BrakeProgram , 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.Itoa(each.AlarmId) //strconv.Formatint(each.AlarmId , 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatBool(each.AlarmToggle)

	}

	err := file.Save(filename)
	if err != nil {
		return err
	}

	return nil
}

func DeserializeMetTower(data []MetTower, j int, typeExcel string, CreateDateTime string) error {
	//savecipo += 1
	filename := ""
	filename = "web/assets/Excel/" + typeExcel + "/" + CreateDateTime + ".xlsx"

	file := x.NewFile()
	sheet, _ := file.AddSheet("Sheet1")
	header := []string{"TimeStamp", "WindDirNo", "WindDirDesc", "WSCategoryNo", "WSCategoryDesc", "VHubWS90mAvg", "VHubWS90mMax", "VHubWS90mMin", "VHubWS90mStdDev", "VHubWS90mCount", "VRefWS88mAvg", "VRefWS88mMax", "VRefWS88mMin", "VRefWS88mStdDev", "VRefWS88mCount", "VTipWS42mAvg", "VTipWS42mMax", "VTipWS42mMin", "VTipWS42mStdDev", "VTipWS42mCount", "DHubWD88mAvg", "DHubWD88mMax", "DHubWD88mMin", "DHubWD88mStdDev", "DHubWD88mCount", "DRefWD86mAvg", "DRefWD86mMax", "DRefWD86mMin", "DRefWD86mStdDev", "DRefWD86mCount", "THubHHubHumid855mAvg", "THubHHubHumid855mMax", "THubHHubHumid855mMin", "THubHHubHumid855mStdDev", "THubHHubHumid855mCount", "TRefHRefHumid855mAvg", "TRefHRefHumid855mMax", "TRefHRefHumid855mMin", "TRefHRefHumid855mStdDev", "TRefHRefHumid855mCount", "THubHHubTemp855mAvg", "THubHHubTemp855mMax", "THubHHubTemp855mMin", "THubHHubTemp855mStdDev", "THubHHubTemp855mCount", "TRefHRefTemp855mAvg", "TRefHRefTemp855mMax", "TRefHRefTemp855mMin", "TRefHRefTemp855mStdDev", "TRefHRefTemp855mCount", "BaroAirPress855mAvg", "BaroAirPress855mMax", "BaroAirPress855mMin", "BaroAirPress855mStdDev", "BaroAirPress855mCount", "YawAngleVoltageAvg", "YawAngleVoltageMax", "YawAngleVoltageMin", "YawAngleVoltageStdDev", "YawAngleVoltageCount", "OtherSensorVoltageAI1Avg", "OtherSensorVoltageAI1Max", "OtherSensorVoltageAI1Min", "OtherSensorVoltageAI1StdDev", "OtherSensorVoltageAI1Count", "OtherSensorVoltageAI2Avg", "OtherSensorVoltageAI2Max", "OtherSensorVoltageAI2Min", "OtherSensorVoltageAI2StdDev", "OtherSensorVoltageAI2Count", "OtherSensorVoltageAI3Avg", "OtherSensorVoltageAI3Max", "OtherSensorVoltageAI3Min", "OtherSensorVoltageAI3StdDev", "OtherSensorVoltageAI3Count", "OtherSensorVoltageAI4Avg", "OtherSensorVoltageAI4Max", "OtherSensorVoltageAI4Min", "OtherSensorVoltageAI4StdDev", "OtherSensorVoltageAI4Count", "GenRPMCurrentAvg", "GenRPMCurrentMax", "GenRPMCurrentMin", "GenRPMCurrentStdDev", "GenRPMCurrentCount", "WS_SCSCurrentAvg", "WS_SCSCurrentMax", "WS_SCSCurrentMin", "WS_SCSCurrentStdDev", "WS_SCSCurrentCount", "RainStatusCount", "RainStatusSum", "OtherSensor2StatusIO1Avg", "OtherSensor2StatusIO1Max", "OtherSensor2StatusIO1Min", "OtherSensor2StatusIO1StdDev", "OtherSensor2StatusIO1Count", "OtherSensor2StatusIO2Avg", "OtherSensor2StatusIO2Max", "OtherSensor2StatusIO2Min", "OtherSensor2StatusIO2StdDev", "OtherSensor2StatusIO2Count", "OtherSensor2StatusIO3Avg", "OtherSensor2StatusIO3Max", "OtherSensor2StatusIO3Min", "OtherSensor2StatusIO3StdDev", "OtherSensor2StatusIO3Count", "OtherSensor2StatusIO4Avg", "OtherSensor2StatusIO4Max", "OtherSensor2StatusIO4Min", "OtherSensor2StatusIO4StdDev", "OtherSensor2StatusIO4Count", "OtherSensor2StatusIO5Avg", "OtherSensor2StatusIO5Max", "OtherSensor2StatusIO5Min", "OtherSensor2StatusIO5StdDev", "OtherSensor2StatusIO5Count", "A1Avg", "A1Max", "A1Min", "A1StdDev", "A1Count", "A2Avg", "A2Max", "A2Min", "A2StdDev", "A2Count", "A3Avg", "A3Max", "A3Min", "A3StdDev", "A3Count", "A4Avg", "A4Max", "A4Min", "A4StdDev", "A4Count", "A5Avg", "A5Max", "A5Min", "A5StdDev", "A5Count", "A6Avg", "A6Max", "A6Min", "A6StdDev", "A6Count", "A7Avg", "A7Max", "A7Min", "A7StdDev", "A7Count", "A8Avg", "A8Max", "A8Min", "A8StdDev", "A8Count", "A9Avg", "A9Max", "A9Min", "A9StdDev", "A9Count", "A10Avg", "A10Max", "A10Min", "A10StdDev", "A10Count", "AC1Avg", "AC1Max", "AC1Min", "AC1StdDev", "AC1Count", "AC2Avg", "AC2Max", "AC2Min", "AC2StdDev", "AC2Count", "C1Avg", "C1Max", "C1Min", "C1StdDev", "C1Count", "C2Avg", "C2Max", "C2Min", "C2StdDev", "C2Count", "C3Avg", "C3Max", "C3Min", "C3StdDev", "C3Count", "D1Avg", "D1Max", "D1Min", "D1StdDev", "M1_1Avg", "M1_1Max", "M1_1Min", "M1_1StdDev", "M1_1Count", "M1_2Avg", "M1_2Max", "M1_2Min", "M1_2StdDev", "M1_2Count", "M1_3Avg", "M1_3Max", "M1_3Min", "M1_3StdDev", "M1_3Count", "M1_4Avg", "M1_4Max", "M1_4Min", "M1_4StdDev", "M1_4Count", "M1_5Avg", "M1_5Max", "M1_5Min", "M1_5StdDev", "M1_5Count", "M2_1Avg", "M2_1Max", "M2_1Min", "M2_1StdDev", "M2_1Count", "M2_2Avg", "M2_2Max", "M2_2Min", "M2_2StdDev", "M2_2Count", "M2_3Avg", "M2_3Max", "M2_3Min", "M2_3StdDev", "M2_3Count", "M2_4Avg", "M2_4Max", "M2_4Min", "M2_4StdDev", "M2_4Count", "M2_5Avg", "M2_5Max", "M2_5Min", "M2_5StdDev", "M2_5Count", "M2_6Avg", "M2_6Max", "M2_6Min", "M2_6StdDev", "M2_6Count", "M2_7Avg", "M2_7Max", "M2_7Min", "M2_7StdDev", "M2_7Count", "M2_8Avg", "M2_8Max", "M2_8Min", "M2_8StdDev", "M2_8Count", "VAvg", "VMax", "VMin", "IAvg", "IMax", "IMin", "T", "Addr"}

	for i, each := range data {
		if i == 0 {
			rowHeader := sheet.AddRow()
			for _, hdr := range header {

				cell := rowHeader.AddCell()
				cell.Value = hdr
			}
		}

		rowContent := sheet.AddRow()

		cell := rowContent.AddCell()
		cell.Value = each.TimeStamp.Format("2006-01-02 15:04:05")

		cell = rowContent.AddCell()
		cell.Value = each.WindDirDesc

		cell = rowContent.AddCell()
		cell.Value = each.WSCategoryDesc

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.VHubWS90mAvg, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.VHubWS90mMax, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.VHubWS90mMin, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.VHubWS90mStdDev, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.VHubWS90mCount, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.VRefWS88mAvg, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.VRefWS88mMax, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.VRefWS88mMin, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.VRefWS88mStdDev, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.VRefWS88mCount, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.VTipWS42mAvg, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.VTipWS42mMax, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.VTipWS42mMin, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.VTipWS42mStdDev, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.VTipWS42mCount, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.DHubWD88mAvg, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.DHubWD88mMax, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.DHubWD88mMin, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.DHubWD88mStdDev, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.DHubWD88mCount, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.DRefWD86mAvg, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.DRefWD86mMax, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.DRefWD86mMin, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.DRefWD86mStdDev, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.DRefWD86mCount, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.THubHHubHumid855mAvg, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.THubHHubHumid855mMax, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.THubHHubHumid855mMin, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.THubHHubHumid855mStdDev, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.THubHHubHumid855mCount, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.TRefHRefHumid855mAvg, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.TRefHRefHumid855mMax, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.TRefHRefHumid855mMin, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.TRefHRefHumid855mStdDev, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.TRefHRefHumid855mCount, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.THubHHubTemp855mAvg, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.THubHHubTemp855mMax, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.THubHHubTemp855mMin, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.THubHHubTemp855mStdDev, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.THubHHubTemp855mCount, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.TRefHRefTemp855mAvg, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.TRefHRefTemp855mMax, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.TRefHRefTemp855mMin, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.TRefHRefTemp855mStdDev, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.TRefHRefTemp855mCount, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.BaroAirPress855mAvg, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.BaroAirPress855mMax, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.BaroAirPress855mMin, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.BaroAirPress855mStdDev, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.BaroAirPress855mCount, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.YawAngleVoltageAvg, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.YawAngleVoltageMax, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.YawAngleVoltageMin, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.YawAngleVoltageStdDev, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.YawAngleVoltageCount, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.OtherSensorVoltageAI1Avg, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.OtherSensorVoltageAI1Max, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.OtherSensorVoltageAI1Min, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.OtherSensorVoltageAI1StdDev, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.OtherSensorVoltageAI1Count, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.OtherSensorVoltageAI2Avg, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.OtherSensorVoltageAI2Max, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.OtherSensorVoltageAI2Min, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.OtherSensorVoltageAI2StdDev, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.OtherSensorVoltageAI2Count, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.OtherSensorVoltageAI3Avg, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.OtherSensorVoltageAI3Max, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.OtherSensorVoltageAI3Min, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.OtherSensorVoltageAI3StdDev, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.OtherSensorVoltageAI3Count, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.OtherSensorVoltageAI4Avg, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.OtherSensorVoltageAI4Max, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.OtherSensorVoltageAI4Min, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.OtherSensorVoltageAI4StdDev, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.OtherSensorVoltageAI4Count, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.GenRPMCurrentAvg, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.GenRPMCurrentMax, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.GenRPMCurrentMin, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.GenRPMCurrentStdDev, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.GenRPMCurrentCount, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.WS_SCSCurrentAvg, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.WS_SCSCurrentMax, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.WS_SCSCurrentMin, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.WS_SCSCurrentStdDev, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.WS_SCSCurrentCount, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.RainStatusCount, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.RainStatusSum, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.OtherSensor2StatusIO1Avg, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.OtherSensor2StatusIO1Max, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.OtherSensor2StatusIO1Min, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.OtherSensor2StatusIO1StdDev, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.OtherSensor2StatusIO1Count, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.OtherSensor2StatusIO2Avg, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.OtherSensor2StatusIO2Max, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.OtherSensor2StatusIO2Min, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.OtherSensor2StatusIO2StdDev, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.OtherSensor2StatusIO2Count, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.OtherSensor2StatusIO3Avg, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.OtherSensor2StatusIO3Max, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.OtherSensor2StatusIO3Min, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.OtherSensor2StatusIO3StdDev, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.OtherSensor2StatusIO3Count, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.OtherSensor2StatusIO4Avg, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.OtherSensor2StatusIO4Max, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.OtherSensor2StatusIO4Min, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.OtherSensor2StatusIO4StdDev, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.OtherSensor2StatusIO4Count, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.OtherSensor2StatusIO5Avg, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.OtherSensor2StatusIO5Max, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.OtherSensor2StatusIO5Min, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.OtherSensor2StatusIO5StdDev, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.OtherSensor2StatusIO5Count, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.A1Avg, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.A1Max, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.A1Min, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.A1StdDev, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.A1Count, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.A2Avg, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.A2Max, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.A2Min, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.A2StdDev, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.A2Count, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.A3Avg, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.A3Max, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.A3Min, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.A3StdDev, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.A3Count, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.A4Avg, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.A4Max, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.A4Min, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.A4StdDev, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.A4Count, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.A5Avg, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.A5Max, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.A5Min, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.A5StdDev, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.A5Count, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.A6Avg, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.A6Max, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.A6Min, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.A6StdDev, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.A6Count, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.A7Avg, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.A7Max, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.A7Min, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.A7StdDev, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.A7Count, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.A8Avg, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.A8Max, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.A8Min, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.A8StdDev, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.A8Count, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.A9Avg, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.A9Max, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.A9Min, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.A9StdDev, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.A9Count, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.A10Avg, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.A10Max, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.A10Min, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.A10StdDev, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.A10Count, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.AC1Avg, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.AC1Max, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.AC1Min, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.AC1StdDev, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.AC1Count, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.AC2Avg, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.AC2Max, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.AC2Min, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.AC2StdDev, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.AC2Count, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.C1Avg, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.C1Max, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.C1Min, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.C1StdDev, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.C1Count, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.C2Avg, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.C2Max, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.C2Min, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.C2StdDev, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.C2Count, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.C3Avg, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.C3Max, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.C3Min, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.C3StdDev, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.C3Count, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.D1Avg, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.D1Max, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.D1Min, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.D1StdDev, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.M1_1Avg, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.M1_1Max, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.M1_1Min, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.M1_1StdDev, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.M1_1Count, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.M1_2Avg, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.M1_2Max, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.M1_2Min, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.M1_2StdDev, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.M1_2Count, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.M1_3Avg, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.M1_3Max, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.M1_3Min, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.M1_3StdDev, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.M1_3Count, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.M1_4Avg, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.M1_4Max, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.M1_4Min, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.M1_4StdDev, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.M1_4Count, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.M1_5Avg, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.M1_5Max, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.M1_5Min, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.M1_5StdDev, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.M1_5Count, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.M2_1Avg, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.M2_1Max, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.M2_1Min, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.M2_1StdDev, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.M2_1Count, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.M2_2Avg, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.M2_2Max, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.M2_2Min, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.M2_2StdDev, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.M2_2Count, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.M2_3Avg, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.M2_3Max, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.M2_3Min, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.M2_3StdDev, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.M2_3Count, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.M2_4Avg, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.M2_4Max, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.M2_4Min, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.M2_4StdDev, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.M2_4Count, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.M2_5Avg, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.M2_5Max, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.M2_5Min, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.M2_5StdDev, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.M2_5Count, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.M2_6Avg, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.M2_6Max, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.M2_6Min, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.M2_6StdDev, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.M2_6Count, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.M2_7Avg, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.M2_7Max, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.M2_7Min, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.M2_7StdDev, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.M2_7Count, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.M2_8Avg, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.M2_8Max, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.M2_8Min, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.M2_8StdDev, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.M2_8Count, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.VAvg, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.VMax, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.VMin, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.IAvg, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.IMax, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.IMin, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.T, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Addr, 'f', -1, 64)

	}

	err := file.Save(filename)
	if err != nil {
		return err
	}

	return nil
}

func DeserializeScadaData(data []ScadaData, j int, typeExcel string, CreateDateTime string, turbinename map[string]string) error {
	//savecipo += 1
	filename := ""
	filename = "web/assets/Excel/" + typeExcel + "/" + CreateDateTime + ".xlsx"

	file := x.NewFile()
	sheet, _ := file.AddSheet("Sheet1")
	header := []string{"TimeStamp", "ProjectName", "Turbine", "Minutes", "TotalTime  ", "GridFrequency  ", "ReactivePower  ", "AlarmExtStopTime ", "AlarmGridDownTime  ", "AlarmInterLineDown ", "AlarmMachDownTime  ", "AlarmOkTime  ", "AlarmUnknownTime ", "AlarmWeatherStop ", "ExternalStopTime ", "GridDownTime ", "GridOkSecs ", "InternalLineDown ", "MachineDownTime  ", "OkSecs ", "OkTime ", "UnknownTime  ", "WeatherStopTime  ", "GeneratorRPM ", "NacelleYawPositionUntwist  ", "NacelleTemperature ", "AdjWindSpeed ", "AmbientTemperature ", "AvgBladeAngle  ", "AvgWindSpeed ", "UnitsGenerated ", "EstimatedPower ", "EstimatedEnergy", "NacelDirection ", "Power  ", "PowerLost  ", "Energy  ", "EnergyLost  ", "RotorRPM ", "WindDirection  ", "DenValue  ", "DenPh", "DenWindSpeed  ", "DenAdjWindSpeed", "DenPower  ", "DenEnergy", "PCValue", "PCValueAdj  ", "PCDeviation", "WSAdjForPC  ", "WSAvgForPC  ", "TotalAvail  ", "MachineAvail  ", "GridAvail", "DenPcDeviation  ", "DenDeviationPct", "DenPcValue  ", "DeviationPct  ", "MTTR ", "MTTF ", "PerformanceIndex"}

	for i, each := range data {
		if i == 0 {
			rowHeader := sheet.AddRow()
			for _, hdr := range header {

				cell := rowHeader.AddCell()
				cell.Value = hdr
			}
		}

		rowContent := sheet.AddRow()

		cell := rowContent.AddCell()
		cell.Value = each.TimeStamp.Format("2006-01-02 15:04:05")

		cell = rowContent.AddCell()
		cell.Value = each.ProjectName

		cell = rowContent.AddCell()
		cell.Value = turbinename[each.Turbine]

		cell = rowContent.AddCell()
		cell.Value = strconv.Itoa(each.Minutes)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.TotalTime, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.GridFrequency, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.ReactivePower, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.AlarmExtStopTime, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.AlarmGridDownTime, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.AlarmInterLineDown, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.AlarmMachDownTime, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.AlarmOkTime, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.AlarmUnknownTime, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.AlarmWeatherStop, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.ExternalStopTime, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.GridDownTime, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.GridOkSecs, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.InternalLineDown, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.MachineDownTime, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.OkSecs, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.OkTime, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.UnknownTime, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.WeatherStopTime, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.GeneratorRPM, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.NacelleYawPositionUntwist, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.NacelleTemperature, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.AdjWindSpeed, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.AmbientTemperature, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.AvgBladeAngle, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.AvgWindSpeed, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.UnitsGenerated, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.EstimatedPower, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.EstimatedEnergy, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.NacelDirection, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Power, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.PowerLost, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Energy, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.EnergyLost, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.RotorRPM, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.WindDirection, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.DenValue, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.DenPh, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.DenWindSpeed, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.DenAdjWindSpeed, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.DenPower, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.DenEnergy, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.PCValue, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.PCValueAdj, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.PCDeviation, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.WSAdjForPC, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.WSAvgForPC, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.TotalAvail, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.MachineAvail, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.GridAvail, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.DenPcDeviation, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.DenDeviationPct, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.DenPcValue, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.DeviationPct, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.MTTR, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.MTTF, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.PerformanceIndex, 'f', -1, 64)

	}

	err := file.Save(filename)
	if err != nil {
		return err
	}

	return nil
}

func DeserializeScadaDataHFD(data []ScadaDataHFD, j int, typeExcel, CreateDateTime string, turbinename map[string]string) error {
	//savecipo += 1
	filename := ""
	filename = "web/assets/Excel/" + typeExcel + "/" + CreateDateTime + ".xlsx"

	file := x.NewFile()
	sheet, _ := file.AddSheet("Sheet1")
	header := []string{"TimeStamp", "ProjectName", "Turbine", "Fast_ActivePower_kW", "Fast_ActivePower_kW_StdDev", "Fast_ActivePower_kW_Min", "Fast_ActivePower_kW_Max", "Fast_ActivePower_kW_Count", "Fast_WindSpeed_ms", "Fast_WindSpeed_ms_StdDev", "Fast_WindSpeed_ms_Min", "Fast_WindSpeed_ms_Max", "Fast_WindSpeed_ms_Count", "Slow_NacellePos", "Slow_NacellePos_StdDev", "Slow_NacellePos_Min", "Slow_NacellePos_Max", "Slow_NacellePos_Count", "Slow_WindDirection", "Slow_WindDirection_StdDev", "Slow_WindDirection_Min", "Slow_WindDirection_Max", "Slow_WindDirection_Count", "Fast_CurrentL3", "Fast_CurrentL3_StdDev", "Fast_CurrentL3_Min", "Fast_CurrentL3_Max", "Fast_CurrentL3_Count", "Fast_CurrentL1", "Fast_CurrentL1_StdDev", "Fast_CurrentL1_Min", "Fast_CurrentL1_Max", "Fast_CurrentL1_Count", "Fast_ActivePowerSetpoint_kW", "Fast_ActivePowerSetpoint_kW_StdDev", "Fast_ActivePowerSetpoint_kW_Min", "Fast_ActivePowerSetpoint_kW_Max", "Fast_ActivePowerSetpoint_kW_Count", "Fast_CurrentL2", "Fast_CurrentL2_StdDev", "Fast_CurrentL2_Min", "Fast_CurrentL2_Max", "Fast_CurrentL2_Count", "Fast_DrTrVibValue", "Fast_DrTrVibValue_StdDev", "Fast_DrTrVibValue_Min", "Fast_DrTrVibValue_Max", "Fast_DrTrVibValue_Count", "Fast_GenSpeed_RPM", "Fast_GenSpeed_RPM_StdDev", "Fast_GenSpeed_RPM_Min", "Fast_GenSpeed_RPM_Max", "Fast_GenSpeed_RPM_Count", "Fast_PitchAccuV1", "Fast_PitchAccuV1_StdDev", "Fast_PitchAccuV1_Min", "Fast_PitchAccuV1_Max", "Fast_PitchAccuV1_Count", "Fast_PitchAngle", "Fast_PitchAngle_StdDev", "Fast_PitchAngle_Min", "Fast_PitchAngle_Max", "Fast_PitchAngle_Count", "Fast_PitchAngle3", "Fast_PitchAngle3_StdDev", "Fast_PitchAngle3_Min", "Fast_PitchAngle3_Max", "Fast_PitchAngle3_Count", "Fast_PitchAngle2", "Fast_PitchAngle2_StdDev", "Fast_PitchAngle2_Min", "Fast_PitchAngle2_Max", "Fast_PitchAngle2_Count", "Fast_PitchConvCurrent1", "Fast_PitchConvCurrent1_StdDev", "Fast_PitchConvCurrent1_Min", "Fast_PitchConvCurrent1_Max", "Fast_PitchConvCurrent1_Count", "Fast_PitchConvCurrent3", "Fast_PitchConvCurrent3_StdDev", "Fast_PitchConvCurrent3_Min", "Fast_PitchConvCurrent3_Max", "Fast_PitchConvCurrent3_Count", "Fast_PitchConvCurrent2", "Fast_PitchConvCurrent2_StdDev", "Fast_PitchConvCurrent2_Min", "Fast_PitchConvCurrent2_Max", "Fast_PitchConvCurrent2_Count", "Fast_PowerFactor", "Fast_PowerFactor_StdDev", "Fast_PowerFactor_Min", "Fast_PowerFactor_Max", "Fast_PowerFactor_Count", "Fast_ReactivePowerSetpointPPC_kVA", "Fast_ReactivePowerSetpointPPC_kVA_StdDev", "Fast_ReactivePowerSetpointPPC_kVA_Min", "Fast_ReactivePowerSetpointPPC_kVA_Max", "Fast_ReactivePowerSetpointPPC_kVA_Count", "Fast_ReactivePower_kVAr", "Fast_ReactivePower_kVAr_StdDev", "Fast_ReactivePower_kVAr_Min", "Fast_ReactivePower_kVAr_Max", "Fast_ReactivePower_kVAr_Count", "Fast_RotorSpeed_RPM", "Fast_RotorSpeed_RPM_StdDev", "Fast_RotorSpeed_RPM_Min", "Fast_RotorSpeed_RPM_Max", "Fast_RotorSpeed_RPM_Count", "Fast_VoltageL1", "Fast_VoltageL1_StdDev", "Fast_VoltageL1_Min", "Fast_VoltageL1_Max", "Fast_VoltageL1_Count", "Fast_VoltageL2", "Fast_VoltageL2_StdDev", "Fast_VoltageL2_Min", "Fast_VoltageL2_Max", "Fast_VoltageL2_Count", "Slow_CapableCapacitiveReactPwr_kVAr", "Slow_CapableCapacitiveReactPwr_kVAr_StdDev", "Slow_CapableCapacitiveReactPwr_kVAr_Min", "Slow_CapableCapacitiveReactPwr_kVAr_Max", "Slow_CapableCapacitiveReactPwr_kVAr_Count", "Slow_CapableInductiveReactPwr_kVAr", "Slow_CapableInductiveReactPwr_kVAr_StdDev", "Slow_CapableInductiveReactPwr_kVAr_Min", "Slow_CapableInductiveReactPwr_kVAr_Max", "Slow_CapableInductiveReactPwr_kVAr_Count", "Slow_DateTime_Sec", "Slow_DateTime_Sec_StdDev", "Slow_DateTime_Sec_Min", "Slow_DateTime_Sec_Max", "Slow_DateTime_Sec_Count", "Fast_PitchAngle1", "Fast_PitchAngle1_StdDev", "Fast_PitchAngle1_Min", "Fast_PitchAngle1_Max", "Fast_PitchAngle1_Count", "Fast_VoltageL3", "Fast_VoltageL3_StdDev", "Fast_VoltageL3_Min", "Fast_VoltageL3_Max", "Fast_VoltageL3_Count", "Slow_CapableCapacitivePwrFactor", "Slow_CapableCapacitivePwrFactor_StdDev", "Slow_CapableCapacitivePwrFactor_Min", "Slow_CapableCapacitivePwrFactor_Max", "Slow_CapableCapacitivePwrFactor_Count", "Fast_Total_Production_kWh", "Fast_Total_Production_kWh_StdDev", "Fast_Total_Production_kWh_Min", "Fast_Total_Production_kWh_Max", "Fast_Total_Production_kWh_Count", "Fast_Total_Prod_Day_kWh", "Fast_Total_Prod_Day_kWh_StdDev", "Fast_Total_Prod_Day_kWh_Min", "Fast_Total_Prod_Day_kWh_Max", "Fast_Total_Prod_Day_kWh_Count", "Fast_Total_Prod_Month_kWh", "Fast_Total_Prod_Month_kWh_StdDev", "Fast_Total_Prod_Month_kWh_Min", "Fast_Total_Prod_Month_kWh_Max", "Fast_Total_Prod_Month_kWh_Count", "Fast_ActivePowerOutPWCSell_kW", "Fast_ActivePowerOutPWCSell_kW_StdDev", "Fast_ActivePowerOutPWCSell_kW_Min", "Fast_ActivePowerOutPWCSell_kW_Max", "Fast_ActivePowerOutPWCSell_kW_Count", "Fast_Frequency_Hz", "Fast_Frequency_Hz_StdDev", "Fast_Frequency_Hz_Min", "Fast_Frequency_Hz_Max", "Fast_Frequency_Hz_Count", "Slow_TempG1L2", "Slow_TempG1L2_StdDev", "Slow_TempG1L2_Min", "Slow_TempG1L2_Max", "Slow_TempG1L2_Count", "Slow_TempG1L3", "Slow_TempG1L3_StdDev", "Slow_TempG1L3_Min", "Slow_TempG1L3_Max", "Slow_TempG1L3_Count", "Slow_TempGearBoxHSSDE", "Slow_TempGearBoxHSSDE_StdDev", "Slow_TempGearBoxHSSDE_Min", "Slow_TempGearBoxHSSDE_Max", "Slow_TempGearBoxHSSDE_Count", "Slow_TempGearBoxIMSNDE", "Slow_TempGearBoxIMSNDE_StdDev", "Slow_TempGearBoxIMSNDE_Min", "Slow_TempGearBoxIMSNDE_Max", "Slow_TempGearBoxIMSNDE_Count", "Slow_TempOutdoor", "Slow_TempOutdoor_StdDev", "Slow_TempOutdoor_Min", "Slow_TempOutdoor_Max", "Slow_TempOutdoor_Count", "Fast_PitchAccuV3", "Fast_PitchAccuV3_StdDev", "Fast_PitchAccuV3_Min", "Fast_PitchAccuV3_Max", "Fast_PitchAccuV3_Count", "Slow_TotalTurbineActiveHours", "Slow_TotalTurbineActiveHours_StdDev", "Slow_TotalTurbineActiveHours_Min", "Slow_TotalTurbineActiveHours_Max", "Slow_TotalTurbineActiveHours_Count", "Slow_TotalTurbineOKHours", "Slow_TotalTurbineOKHours_StdDev", "Slow_TotalTurbineOKHours_Min", "Slow_TotalTurbineOKHours_Max", "Slow_TotalTurbineOKHours_Count", "Slow_TotalTurbineTimeAllHours", "Slow_TotalTurbineTimeAllHours_StdDev", "Slow_TotalTurbineTimeAllHours_Min", "Slow_TotalTurbineTimeAllHours_Max", "Slow_TotalTurbineTimeAllHours_Count", "Slow_TempG1L1", "Slow_TempG1L1_StdDev", "Slow_TempG1L1_Min", "Slow_TempG1L1_Max", "Slow_TempG1L1_Count", "Slow_TempGearBoxOilSump", "Slow_TempGearBoxOilSump_StdDev", "Slow_TempGearBoxOilSump_Min", "Slow_TempGearBoxOilSump_Max", "Slow_TempGearBoxOilSump_Count", "Fast_PitchAccuV2", "Fast_PitchAccuV2_StdDev", "Fast_PitchAccuV2_Min", "Fast_PitchAccuV2_Max", "Fast_PitchAccuV2_Count", "Slow_TotalGridOkHours", "Slow_TotalGridOkHours_StdDev", "Slow_TotalGridOkHours_Min", "Slow_TotalGridOkHours_Max", "Slow_TotalGridOkHours_Count", "Slow_TotalActPowerOut_kWh", "Slow_TotalActPowerOut_kWh_StdDev", "Slow_TotalActPowerOut_kWh_Min", "Slow_TotalActPowerOut_kWh_Max", "Slow_TotalActPowerOut_kWh_Count", "Fast_YawService", "Fast_YawService_StdDev", "Fast_YawService_Min", "Fast_YawService_Max", "Fast_YawService_Count", "Fast_YawAngle", "Fast_YawAngle_StdDev", "Fast_YawAngle_Min", "Fast_YawAngle_Max", "Fast_YawAngle_Count", "Slow_CapableInductivePwrFactor", "Slow_CapableInductivePwrFactor_StdDev", "Slow_CapableInductivePwrFactor_Min", "Slow_CapableInductivePwrFactor_Max", "Slow_CapableInductivePwrFactor_Count", "Slow_TempGearBoxHSSNDE", "Slow_TempGearBoxHSSNDE_StdDev", "Slow_TempGearBoxHSSNDE_Min", "Slow_TempGearBoxHSSNDE_Max", "Slow_TempGearBoxHSSNDE_Count", "Slow_TempHubBearing", "Slow_TempHubBearing_StdDev", "Slow_TempHubBearing_Min", "Slow_TempHubBearing_Max", "Slow_TempHubBearing_Count", "Slow_TotalG1ActiveHours", "Slow_TotalG1ActiveHours_StdDev", "Slow_TotalG1ActiveHours_Min", "Slow_TotalG1ActiveHours_Max", "Slow_TotalG1ActiveHours_Count", "Slow_TotalActPowerOutG1_kWh", "Slow_TotalActPowerOutG1_kWh_StdDev", "Slow_TotalActPowerOutG1_kWh_Min", "Slow_TotalActPowerOutG1_kWh_Max", "Slow_TotalActPowerOutG1_kWh_Count", "Slow_TotalReactPowerInG1_kVArh", "Slow_TotalReactPowerInG1_kVArh_StdDev", "Slow_TotalReactPowerInG1_kVArh_Min", "Slow_TotalReactPowerInG1_kVArh_Max", "Slow_TotalReactPowerInG1_kVArh_Count", "Slow_NacelleDrill", "Slow_NacelleDrill_StdDev", "Slow_NacelleDrill_Min", "Slow_NacelleDrill_Max", "Slow_NacelleDrill_Count", "Slow_TempGearBoxIMSDE", "Slow_TempGearBoxIMSDE_StdDev", "Slow_TempGearBoxIMSDE_Min", "Slow_TempGearBoxIMSDE_Max", "Slow_TempGearBoxIMSDE_Count", "Fast_Total_Operating_hrs", "Fast_Total_Operating_hrs_StdDev", "Fast_Total_Operating_hrs_Min", "Fast_Total_Operating_hrs_Max", "Fast_Total_Operating_hrs_Count", "Slow_TempNacelle", "Slow_TempNacelle_StdDev", "Slow_TempNacelle_Min", "Slow_TempNacelle_Max", "Slow_TempNacelle_Count", "Fast_Total_Grid_OK_hrs", "Fast_Total_Grid_OK_hrs_StdDev", "Fast_Total_Grid_OK_hrs_Min", "Fast_Total_Grid_OK_hrs_Max", "Fast_Total_Grid_OK_hrs_Count", "Fast_Total_WTG_OK_hrs", "Fast_Total_WTG_OK_hrs_StdDev", "Fast_Total_WTG_OK_hrs_Min", "Fast_Total_WTG_OK_hrs_Max", "Fast_Total_WTG_OK_hrs_Count", "Slow_TempCabinetTopBox", "Slow_TempCabinetTopBox_StdDev", "Slow_TempCabinetTopBox_Min", "Slow_TempCabinetTopBox_Max", "Slow_TempCabinetTopBox_Count", "Slow_TempGeneratorBearingNDE", "Slow_TempGeneratorBearingNDE_StdDev", "Slow_TempGeneratorBearingNDE_Min", "Slow_TempGeneratorBearingNDE_Max", "Slow_TempGeneratorBearingNDE_Count", "Fast_Total_Access_hrs", "Fast_Total_Access_hrs_StdDev", "Fast_Total_Access_hrs_Min", "Fast_Total_Access_hrs_Max", "Fast_Total_Access_hrs_Count", "Slow_TempBottomPowerSection", "Slow_TempBottomPowerSection_StdDev", "Slow_TempBottomPowerSection_Min", "Slow_TempBottomPowerSection_Max", "Slow_TempBottomPowerSection_Count", "Slow_TempGeneratorBearingDE", "Slow_TempGeneratorBearingDE_StdDev", "Slow_TempGeneratorBearingDE_Min", "Slow_TempGeneratorBearingDE_Max", "Slow_TempGeneratorBearingDE_Count", "Slow_TotalReactPowerIn_kVArh", "Slow_TotalReactPowerIn_kVArh_StdDev", "Slow_TotalReactPowerIn_kVArh_Min", "Slow_TotalReactPowerIn_kVArh_Max", "Slow_TotalReactPowerIn_kVArh_Count", "Slow_TempBottomControlSection", "Slow_TempBottomControlSection_StdDev", "Slow_TempBottomControlSection_Min", "Slow_TempBottomControlSection_Max", "Slow_TempBottomControlSection_Count", "Slow_TempConv1", "Slow_TempConv1_StdDev", "Slow_TempConv1_Min", "Slow_TempConv1_Max", "Slow_TempConv1_Count", "Fast_ActivePowerRated_kW", "Fast_ActivePowerRated_kW_StdDev", "Fast_ActivePowerRated_kW_Min", "Fast_ActivePowerRated_kW_Max", "Fast_ActivePowerRated_kW_Count", "Fast_NodeIP", "Fast_NodeIP_StdDev", "Fast_NodeIP_Min", "Fast_NodeIP_Max", "Fast_NodeIP_Count", "Fast_PitchSpeed1", "Fast_PitchSpeed1_StdDev", "Fast_PitchSpeed1_Min", "Fast_PitchSpeed1_Max", "Fast_PitchSpeed1_Count", "Slow_CFCardSize", "Slow_CFCardSize_StdDev", "Slow_CFCardSize_Min", "Slow_CFCardSize_Max", "Slow_CFCardSize_Count", "Slow_CPU_Number", "Slow_CPU_Number_StdDev", "Slow_CPU_Number_Min", "Slow_CPU_Number_Max", "Slow_CPU_Number_Count", "Slow_CFCardSpaceLeft", "Slow_CFCardSpaceLeft_StdDev", "Slow_CFCardSpaceLeft_Min", "Slow_CFCardSpaceLeft_Max", "Slow_CFCardSpaceLeft_Count", "Slow_TempBottomCapSection", "Slow_TempBottomCapSection_StdDev", "Slow_TempBottomCapSection_Min", "Slow_TempBottomCapSection_Max", "Slow_TempBottomCapSection_Count", "Slow_RatedPower", "Slow_RatedPower_StdDev", "Slow_RatedPower_Min", "Slow_RatedPower_Max", "Slow_RatedPower_Count", "Slow_TempConv3", "Slow_TempConv3_StdDev", "Slow_TempConv3_Min", "Slow_TempConv3_Max", "Slow_TempConv3_Count", "Slow_TempConv2", "Slow_TempConv2_StdDev", "Slow_TempConv2_Min", "Slow_TempConv2_Max", "Slow_TempConv2_Count", "Slow_TotalActPowerIn_kWh", "Slow_TotalActPowerIn_kWh_StdDev", "Slow_TotalActPowerIn_kWh_Min", "Slow_TotalActPowerIn_kWh_Max", "Slow_TotalActPowerIn_kWh_Count", "Slow_TotalActPowerInG1_kWh", "Slow_TotalActPowerInG1_kWh_StdDev", "Slow_TotalActPowerInG1_kWh_Min", "Slow_TotalActPowerInG1_kWh_Max", "Slow_TotalActPowerInG1_kWh_Count", "Slow_TotalActPowerInG2_kWh", "Slow_TotalActPowerInG2_kWh_StdDev", "Slow_TotalActPowerInG2_kWh_Min", "Slow_TotalActPowerInG2_kWh_Max", "Slow_TotalActPowerInG2_kWh_Count", "Slow_TotalActPowerOutG2_kWh", "Slow_TotalActPowerOutG2_kWh_StdDev", "Slow_TotalActPowerOutG2_kWh_Min", "Slow_TotalActPowerOutG2_kWh_Max", "Slow_TotalActPowerOutG2_kWh_Count", "Slow_TotalG2ActiveHours", "Slow_TotalG2ActiveHours_StdDev", "Slow_TotalG2ActiveHours_Min", "Slow_TotalG2ActiveHours_Max", "Slow_TotalG2ActiveHours_Count", "Slow_TotalReactPowerInG2_kVArh", "Slow_TotalReactPowerInG2_kVArh_StdDev", "Slow_TotalReactPowerInG2_kVArh_Min", "Slow_TotalReactPowerInG2_kVArh_Max", "Slow_TotalReactPowerInG2_kVArh_Count", "Slow_TotalReactPowerOut_kVArh", "Slow_TotalReactPowerOut_kVArh_StdDev", "Slow_TotalReactPowerOut_kVArh_Min", "Slow_TotalReactPowerOut_kVArh_Max", "Slow_TotalReactPowerOut_kVArh_Count", "Slow_UTCoffset_int", "Slow_UTCoffset_int_StdDev", "Slow_UTCoffset_int_Min", "Slow_UTCoffset_int_Max", "Slow_UTCoffset_int_Count"}

	for i, each := range data {
		if i == 0 {
			rowHeader := sheet.AddRow()
			for _, hdr := range header {

				cell := rowHeader.AddCell()
				cell.Value = hdr
			}
		}

		rowContent := sheet.AddRow()

		cell := rowContent.AddCell()
		cell.Value = each.TimeStamp.Format("2006-01-02 15:04:05")

		cell = rowContent.AddCell()
		cell.Value = each.ProjectName

		cell = rowContent.AddCell()

		cell.Value = turbinename[each.Turbine]

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Fast_ActivePower_kW, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Fast_ActivePower_kW_StdDev, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Fast_ActivePower_kW_Min, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Fast_ActivePower_kW_Max, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.Itoa(each.Fast_ActivePower_kW_Count)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Fast_WindSpeed_ms, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Fast_WindSpeed_ms_StdDev, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Fast_WindSpeed_ms_Min, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Fast_WindSpeed_ms_Max, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.Itoa(each.Fast_WindSpeed_ms_Count)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Slow_NacellePos, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Slow_NacellePos_StdDev, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Slow_NacellePos_Min, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Slow_NacellePos_Max, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.Itoa(each.Slow_NacellePos_Count)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Slow_WindDirection, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Slow_WindDirection_StdDev, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Slow_WindDirection_Min, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Slow_WindDirection_Max, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.Itoa(each.Slow_WindDirection_Count)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Fast_CurrentL3, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Fast_CurrentL3_StdDev, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Fast_CurrentL3_Min, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Fast_CurrentL3_Max, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.Itoa(each.Fast_CurrentL3_Count)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Fast_CurrentL1, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Fast_CurrentL1_StdDev, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Fast_CurrentL1_Min, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Fast_CurrentL1_Max, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.Itoa(each.Fast_CurrentL1_Count)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Fast_ActivePowerSetpoint_kW, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Fast_ActivePowerSetpoint_kW_StdDev, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Fast_ActivePowerSetpoint_kW_Min, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Fast_ActivePowerSetpoint_kW_Max, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.Itoa(each.Fast_ActivePowerSetpoint_kW_Count)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Fast_CurrentL2, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Fast_CurrentL2_StdDev, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Fast_CurrentL2_Min, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Fast_CurrentL2_Max, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.Itoa(each.Fast_CurrentL2_Count)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Fast_DrTrVibValue, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Fast_DrTrVibValue_StdDev, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Fast_DrTrVibValue_Min, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Fast_DrTrVibValue_Max, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.Itoa(each.Fast_DrTrVibValue_Count)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Fast_GenSpeed_RPM, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Fast_GenSpeed_RPM_StdDev, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Fast_GenSpeed_RPM_Min, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Fast_GenSpeed_RPM_Max, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.Itoa(each.Fast_GenSpeed_RPM_Count)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Fast_PitchAccuV1, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Fast_PitchAccuV1_StdDev, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Fast_PitchAccuV1_Min, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Fast_PitchAccuV1_Max, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.Itoa(each.Fast_PitchAccuV1_Count)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Fast_PitchAngle, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Fast_PitchAngle_StdDev, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Fast_PitchAngle_Min, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Fast_PitchAngle_Max, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.Itoa(each.Fast_PitchAngle_Count)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Fast_PitchAngle3, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Fast_PitchAngle3_StdDev, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Fast_PitchAngle3_Min, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Fast_PitchAngle3_Max, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.Itoa(each.Fast_PitchAngle3_Count)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Fast_PitchAngle2, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Fast_PitchAngle2_StdDev, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Fast_PitchAngle2_Min, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Fast_PitchAngle2_Max, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.Itoa(each.Fast_PitchAngle2_Count)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Fast_PitchConvCurrent1, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Fast_PitchConvCurrent1_StdDev, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Fast_PitchConvCurrent1_Min, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Fast_PitchConvCurrent1_Max, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.Itoa(each.Fast_PitchConvCurrent1_Count)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Fast_PitchConvCurrent3, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Fast_PitchConvCurrent3_StdDev, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Fast_PitchConvCurrent3_Min, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Fast_PitchConvCurrent3_Max, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.Itoa(each.Fast_PitchConvCurrent3_Count)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Fast_PitchConvCurrent2, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Fast_PitchConvCurrent2_StdDev, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Fast_PitchConvCurrent2_Min, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Fast_PitchConvCurrent2_Max, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.Itoa(each.Fast_PitchConvCurrent2_Count)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Fast_PowerFactor, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Fast_PowerFactor_StdDev, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Fast_PowerFactor_Min, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Fast_PowerFactor_Max, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.Itoa(each.Fast_PowerFactor_Count)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Fast_ReactivePowerSetpointPPC_kVA, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Fast_ReactivePowerSetpointPPC_kVA_StdDev, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Fast_ReactivePowerSetpointPPC_kVA_Min, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Fast_ReactivePowerSetpointPPC_kVA_Max, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.Itoa(each.Fast_ReactivePowerSetpointPPC_kVA_Count)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Fast_ReactivePower_kVAr, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Fast_ReactivePower_kVAr_StdDev, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Fast_ReactivePower_kVAr_Min, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Fast_ReactivePower_kVAr_Max, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.Itoa(each.Fast_ReactivePower_kVAr_Count)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Fast_RotorSpeed_RPM, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Fast_RotorSpeed_RPM_StdDev, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Fast_RotorSpeed_RPM_Min, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Fast_RotorSpeed_RPM_Max, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.Itoa(each.Fast_RotorSpeed_RPM_Count)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Fast_VoltageL1, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Fast_VoltageL1_StdDev, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Fast_VoltageL1_Min, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Fast_VoltageL1_Max, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.Itoa(each.Fast_VoltageL1_Count)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Fast_VoltageL2, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Fast_VoltageL2_StdDev, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Fast_VoltageL2_Min, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Fast_VoltageL2_Max, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.Itoa(each.Fast_VoltageL2_Count)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Slow_CapableCapacitiveReactPwr_kVAr, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Slow_CapableCapacitiveReactPwr_kVAr_StdDev, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Slow_CapableCapacitiveReactPwr_kVAr_Min, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Slow_CapableCapacitiveReactPwr_kVAr_Max, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.Itoa(each.Slow_CapableCapacitiveReactPwr_kVAr_Count)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Slow_CapableInductiveReactPwr_kVAr, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Slow_CapableInductiveReactPwr_kVAr_StdDev, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Slow_CapableInductiveReactPwr_kVAr_Min, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Slow_CapableInductiveReactPwr_kVAr_Max, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.Itoa(each.Slow_CapableInductiveReactPwr_kVAr_Count)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Slow_DateTime_Sec, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Slow_DateTime_Sec_StdDev, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Slow_DateTime_Sec_Min, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Slow_DateTime_Sec_Max, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.Itoa(each.Slow_DateTime_Sec_Count)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Fast_PitchAngle1, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Fast_PitchAngle1_StdDev, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Fast_PitchAngle1_Min, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Fast_PitchAngle1_Max, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.Itoa(each.Fast_PitchAngle1_Count)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Fast_VoltageL3, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Fast_VoltageL3_StdDev, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Fast_VoltageL3_Min, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Fast_VoltageL3_Max, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.Itoa(each.Fast_VoltageL3_Count)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Slow_CapableCapacitivePwrFactor, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Slow_CapableCapacitivePwrFactor_StdDev, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Slow_CapableCapacitivePwrFactor_Min, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Slow_CapableCapacitivePwrFactor_Max, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.Itoa(each.Slow_CapableCapacitivePwrFactor_Count)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Fast_Total_Production_kWh, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Fast_Total_Production_kWh_StdDev, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Fast_Total_Production_kWh_Min, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Fast_Total_Production_kWh_Max, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.Itoa(each.Fast_Total_Production_kWh_Count)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Fast_Total_Prod_Day_kWh, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Fast_Total_Prod_Day_kWh_StdDev, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Fast_Total_Prod_Day_kWh_Min, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Fast_Total_Prod_Day_kWh_Max, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.Itoa(each.Fast_Total_Prod_Day_kWh_Count)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Fast_Total_Prod_Month_kWh, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Fast_Total_Prod_Month_kWh_StdDev, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Fast_Total_Prod_Month_kWh_Min, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Fast_Total_Prod_Month_kWh_Max, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.Itoa(each.Fast_Total_Prod_Month_kWh_Count)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Fast_ActivePowerOutPWCSell_kW, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Fast_ActivePowerOutPWCSell_kW_StdDev, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Fast_ActivePowerOutPWCSell_kW_Min, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Fast_ActivePowerOutPWCSell_kW_Max, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.Itoa(each.Fast_ActivePowerOutPWCSell_kW_Count)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Fast_Frequency_Hz, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Fast_Frequency_Hz_StdDev, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Fast_Frequency_Hz_Min, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Fast_Frequency_Hz_Max, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.Itoa(each.Fast_Frequency_Hz_Count)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Slow_TempG1L2, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Slow_TempG1L2_StdDev, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Slow_TempG1L2_Min, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Slow_TempG1L2_Max, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.Itoa(each.Slow_TempG1L2_Count)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Slow_TempG1L3, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Slow_TempG1L3_StdDev, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Slow_TempG1L3_Min, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Slow_TempG1L3_Max, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.Itoa(each.Slow_TempG1L3_Count)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Slow_TempGearBoxHSSDE, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Slow_TempGearBoxHSSDE_StdDev, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Slow_TempGearBoxHSSDE_Min, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Slow_TempGearBoxHSSDE_Max, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.Itoa(each.Slow_TempGearBoxHSSDE_Count)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Slow_TempGearBoxIMSNDE, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Slow_TempGearBoxIMSNDE_StdDev, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Slow_TempGearBoxIMSNDE_Min, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Slow_TempGearBoxIMSNDE_Max, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.Itoa(each.Slow_TempGearBoxIMSNDE_Count)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Slow_TempOutdoor, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Slow_TempOutdoor_StdDev, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Slow_TempOutdoor_Min, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Slow_TempOutdoor_Max, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.Itoa(each.Slow_TempOutdoor_Count)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Fast_PitchAccuV3, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Fast_PitchAccuV3_StdDev, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Fast_PitchAccuV3_Min, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Fast_PitchAccuV3_Max, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.Itoa(each.Fast_PitchAccuV3_Count)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Slow_TotalTurbineActiveHours, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Slow_TotalTurbineActiveHours_StdDev, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Slow_TotalTurbineActiveHours_Min, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Slow_TotalTurbineActiveHours_Max, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.Itoa(each.Slow_TotalTurbineActiveHours_Count)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Slow_TotalTurbineOKHours, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Slow_TotalTurbineOKHours_StdDev, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Slow_TotalTurbineOKHours_Min, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Slow_TotalTurbineOKHours_Max, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.Itoa(each.Slow_TotalTurbineOKHours_Count)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Slow_TotalTurbineTimeAllHours, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Slow_TotalTurbineTimeAllHours_StdDev, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Slow_TotalTurbineTimeAllHours_Min, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Slow_TotalTurbineTimeAllHours_Max, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.Itoa(each.Slow_TotalTurbineTimeAllHours_Count)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Slow_TempG1L1, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Slow_TempG1L1_StdDev, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Slow_TempG1L1_Min, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Slow_TempG1L1_Max, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.Itoa(each.Slow_TempG1L1_Count)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Slow_TempGearBoxOilSump, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Slow_TempGearBoxOilSump_StdDev, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Slow_TempGearBoxOilSump_Min, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Slow_TempGearBoxOilSump_Max, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.Itoa(each.Slow_TempGearBoxOilSump_Count)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Fast_PitchAccuV2, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Fast_PitchAccuV2_StdDev, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Fast_PitchAccuV2_Min, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Fast_PitchAccuV2_Max, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.Itoa(each.Fast_PitchAccuV2_Count)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Slow_TotalGridOkHours, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Slow_TotalGridOkHours_StdDev, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Slow_TotalGridOkHours_Min, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Slow_TotalGridOkHours_Max, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.Itoa(each.Slow_TotalGridOkHours_Count)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Slow_TotalActPowerOut_kWh, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Slow_TotalActPowerOut_kWh_StdDev, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Slow_TotalActPowerOut_kWh_Min, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Slow_TotalActPowerOut_kWh_Max, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.Itoa(each.Slow_TotalActPowerOut_kWh_Count)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Fast_YawService, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Fast_YawService_StdDev, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Fast_YawService_Min, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Fast_YawService_Max, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.Itoa(each.Fast_YawService_Count)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Fast_YawAngle, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Fast_YawAngle_StdDev, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Fast_YawAngle_Min, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Fast_YawAngle_Max, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.Itoa(each.Fast_YawAngle_Count)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Slow_CapableInductivePwrFactor, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Slow_CapableInductivePwrFactor_StdDev, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Slow_CapableInductivePwrFactor_Min, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Slow_CapableInductivePwrFactor_Max, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.Itoa(each.Slow_CapableInductivePwrFactor_Count)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Slow_TempGearBoxHSSNDE, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Slow_TempGearBoxHSSNDE_StdDev, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Slow_TempGearBoxHSSNDE_Min, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Slow_TempGearBoxHSSNDE_Max, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.Itoa(each.Slow_TempGearBoxHSSNDE_Count)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Slow_TempHubBearing, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Slow_TempHubBearing_StdDev, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Slow_TempHubBearing_Min, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Slow_TempHubBearing_Max, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.Itoa(each.Slow_TempHubBearing_Count)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Slow_TotalG1ActiveHours, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Slow_TotalG1ActiveHours_StdDev, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Slow_TotalG1ActiveHours_Min, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Slow_TotalG1ActiveHours_Max, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.Itoa(each.Slow_TotalG1ActiveHours_Count)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Slow_TotalActPowerOutG1_kWh, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Slow_TotalActPowerOutG1_kWh_StdDev, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Slow_TotalActPowerOutG1_kWh_Min, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Slow_TotalActPowerOutG1_kWh_Max, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.Itoa(each.Slow_TotalActPowerOutG1_kWh_Count)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Slow_TotalReactPowerInG1_kVArh, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Slow_TotalReactPowerInG1_kVArh_StdDev, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Slow_TotalReactPowerInG1_kVArh_Min, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Slow_TotalReactPowerInG1_kVArh_Max, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.Itoa(each.Slow_TotalReactPowerInG1_kVArh_Count)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Slow_NacelleDrill, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Slow_NacelleDrill_StdDev, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Slow_NacelleDrill_Min, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Slow_NacelleDrill_Max, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.Itoa(each.Slow_NacelleDrill_Count)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Slow_TempGearBoxIMSDE, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Slow_TempGearBoxIMSDE_StdDev, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Slow_TempGearBoxIMSDE_Min, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Slow_TempGearBoxIMSDE_Max, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.Itoa(each.Slow_TempGearBoxIMSDE_Count)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Fast_Total_Operating_hrs, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Fast_Total_Operating_hrs_StdDev, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Fast_Total_Operating_hrs_Min, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Fast_Total_Operating_hrs_Max, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.Itoa(each.Fast_Total_Operating_hrs_Count)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Slow_TempNacelle, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Slow_TempNacelle_StdDev, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Slow_TempNacelle_Min, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Slow_TempNacelle_Max, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.Itoa(each.Slow_TempNacelle_Count)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Fast_Total_Grid_OK_hrs, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Fast_Total_Grid_OK_hrs_StdDev, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Fast_Total_Grid_OK_hrs_Min, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Fast_Total_Grid_OK_hrs_Max, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.Itoa(each.Fast_Total_Grid_OK_hrs_Count)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Fast_Total_WTG_OK_hrs, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Fast_Total_WTG_OK_hrs_StdDev, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Fast_Total_WTG_OK_hrs_Min, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Fast_Total_WTG_OK_hrs_Max, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.Itoa(each.Fast_Total_WTG_OK_hrs_Count)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Slow_TempCabinetTopBox, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Slow_TempCabinetTopBox_StdDev, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Slow_TempCabinetTopBox_Min, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Slow_TempCabinetTopBox_Max, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.Itoa(each.Slow_TempCabinetTopBox_Count)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Slow_TempGeneratorBearingNDE, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Slow_TempGeneratorBearingNDE_StdDev, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Slow_TempGeneratorBearingNDE_Min, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Slow_TempGeneratorBearingNDE_Max, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.Itoa(each.Slow_TempGeneratorBearingNDE_Count)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Fast_Total_Access_hrs, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Fast_Total_Access_hrs_StdDev, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Fast_Total_Access_hrs_Min, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Fast_Total_Access_hrs_Max, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.Itoa(each.Fast_Total_Access_hrs_Count)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Slow_TempBottomPowerSection, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Slow_TempBottomPowerSection_StdDev, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Slow_TempBottomPowerSection_Min, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Slow_TempBottomPowerSection_Max, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.Itoa(each.Slow_TempBottomPowerSection_Count)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Slow_TempGeneratorBearingDE, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Slow_TempGeneratorBearingDE_StdDev, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Slow_TempGeneratorBearingDE_Min, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Slow_TempGeneratorBearingDE_Max, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.Itoa(each.Slow_TempGeneratorBearingDE_Count)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Slow_TotalReactPowerIn_kVArh, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Slow_TotalReactPowerIn_kVArh_StdDev, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Slow_TotalReactPowerIn_kVArh_Min, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Slow_TotalReactPowerIn_kVArh_Max, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.Itoa(each.Slow_TotalReactPowerIn_kVArh_Count)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Slow_TempBottomControlSection, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Slow_TempBottomControlSection_StdDev, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Slow_TempBottomControlSection_Min, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Slow_TempBottomControlSection_Max, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.Itoa(each.Slow_TempBottomControlSection_Count)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Slow_TempConv1, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Slow_TempConv1_StdDev, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Slow_TempConv1_Min, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Slow_TempConv1_Max, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.Itoa(each.Slow_TempConv1_Count)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Fast_ActivePowerRated_kW, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Fast_ActivePowerRated_kW_StdDev, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Fast_ActivePowerRated_kW_Min, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Fast_ActivePowerRated_kW_Max, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.Itoa(each.Fast_ActivePowerRated_kW_Count)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Fast_NodeIP, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Fast_NodeIP_StdDev, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Fast_NodeIP_Min, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Fast_NodeIP_Max, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.Itoa(each.Fast_NodeIP_Count)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Fast_PitchSpeed1, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Fast_PitchSpeed1_StdDev, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Fast_PitchSpeed1_Min, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Fast_PitchSpeed1_Max, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.Itoa(each.Fast_PitchSpeed1_Count)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Slow_CFCardSize, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Slow_CFCardSize_StdDev, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Slow_CFCardSize_Min, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Slow_CFCardSize_Max, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.Itoa(each.Slow_CFCardSize_Count)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Slow_CPU_Number, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Slow_CPU_Number_StdDev, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Slow_CPU_Number_Min, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Slow_CPU_Number_Max, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.Itoa(each.Slow_CPU_Number_Count)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Slow_CFCardSpaceLeft, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Slow_CFCardSpaceLeft_StdDev, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Slow_CFCardSpaceLeft_Min, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Slow_CFCardSpaceLeft_Max, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.Itoa(each.Slow_CFCardSpaceLeft_Count)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Slow_TempBottomCapSection, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Slow_TempBottomCapSection_StdDev, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Slow_TempBottomCapSection_Min, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Slow_TempBottomCapSection_Max, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.Itoa(each.Slow_TempBottomCapSection_Count)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Slow_RatedPower, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Slow_RatedPower_StdDev, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Slow_RatedPower_Min, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Slow_RatedPower_Max, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.Itoa(each.Slow_RatedPower_Count)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Slow_TempConv3, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Slow_TempConv3_StdDev, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Slow_TempConv3_Min, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Slow_TempConv3_Max, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.Itoa(each.Slow_TempConv3_Count)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Slow_TempConv2, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Slow_TempConv2_StdDev, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Slow_TempConv2_Min, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Slow_TempConv2_Max, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.Itoa(each.Slow_TempConv2_Count)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Slow_TotalActPowerIn_kWh, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Slow_TotalActPowerIn_kWh_StdDev, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Slow_TotalActPowerIn_kWh_Min, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Slow_TotalActPowerIn_kWh_Max, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.Itoa(each.Slow_TotalActPowerIn_kWh_Count)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Slow_TotalActPowerInG1_kWh, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Slow_TotalActPowerInG1_kWh_StdDev, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Slow_TotalActPowerInG1_kWh_Min, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Slow_TotalActPowerInG1_kWh_Max, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.Itoa(each.Slow_TotalActPowerInG1_kWh_Count)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Slow_TotalActPowerInG2_kWh, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Slow_TotalActPowerInG2_kWh_StdDev, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Slow_TotalActPowerInG2_kWh_Min, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Slow_TotalActPowerInG2_kWh_Max, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.Itoa(each.Slow_TotalActPowerInG2_kWh_Count)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Slow_TotalActPowerOutG2_kWh, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Slow_TotalActPowerOutG2_kWh_StdDev, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Slow_TotalActPowerOutG2_kWh_Min, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Slow_TotalActPowerOutG2_kWh_Max, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.Itoa(each.Slow_TotalActPowerOutG2_kWh_Count)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Slow_TotalG2ActiveHours, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Slow_TotalG2ActiveHours_StdDev, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Slow_TotalG2ActiveHours_Min, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Slow_TotalG2ActiveHours_Max, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.Itoa(each.Slow_TotalG2ActiveHours_Count)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Slow_TotalReactPowerInG2_kVArh, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Slow_TotalReactPowerInG2_kVArh_StdDev, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Slow_TotalReactPowerInG2_kVArh_Min, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Slow_TotalReactPowerInG2_kVArh_Max, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.Itoa(each.Slow_TotalReactPowerInG2_kVArh_Count)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Slow_TotalReactPowerOut_kVArh, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Slow_TotalReactPowerOut_kVArh_StdDev, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Slow_TotalReactPowerOut_kVArh_Min, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Slow_TotalReactPowerOut_kVArh_Max, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.Itoa(each.Slow_TotalReactPowerOut_kVArh_Count)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Slow_UTCoffset_int, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Slow_UTCoffset_int_StdDev, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Slow_UTCoffset_int_Min, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.FormatFloat(each.Slow_UTCoffset_int_Max, 'f', -1, 64)

		cell = rowContent.AddCell()
		cell.Value = strconv.Itoa(each.Slow_UTCoffset_int_Count)

	}

	err := file.Save(filename)
	if err != nil {
		return err
	}

	return nil
}

// FUNCTION

func SecondsToHms(d float64) string {

	duration := tk.ToInt(d, tk.RoundingUp)

	var h = duration / 3600
	var m = duration % 3600 / 60
	var s = duration % 3600 % 60
	res := ""
	hstring := ""
	mstring := ""
	sstring := ""

	if h > 0 {
		if h < 10 {
			hstring = "0" + tk.ToString(h)
		} else {
			hstring = tk.ToString(h)
		}
	} else {
		hstring = "00"
	}

	if m > 0 {
		if m < 10 {
			mstring = "0" + tk.ToString(m)
		} else {
			mstring = tk.ToString(m)
		}
	} else {
		mstring = "00"
	}

	if s > 0 {
		if s < 10 {
			sstring = "0" + tk.ToString(s)
		} else {
			sstring = tk.ToString(s)
		}
	} else {
		sstring = "00"
	}

	res = hstring + ":" + mstring + ":" + sstring

	return res
}

func (m *DataBrowserController) GetDowntimeEventListHFD(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	var filter []*dbox.Filter

	p := new(helper.Payloads)
	e := k.GetPayload(&p)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	filter, _ = p.ParseFilter()

	query := DB().Connection.NewQuery().From(new(EventDownHFD).TableName()).Skip(p.Skip).Take(p.Take)
	query.Where(dbox.And(filter...))

	if len(p.Sort) > 0 {
		var arrsort []string
		for _, val := range p.Sort {
			if val.Dir == "desc" {
				arrsort = append(arrsort, strings.ToLower("-"+val.Field))
			} else {
				arrsort = append(arrsort, strings.ToLower(val.Field))
			}
		}
		query = query.Order(arrsort...)
	}
	csr, e := query.Cursor(nil)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}
	defer csr.Close()

	tmpResult := make([]EventDown, 0)
	e = csr.Fetch(&tmpResult, 0, false)

	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	queryC := DB().Connection.NewQuery().From(new(EventDownHFD).TableName()).Where(dbox.And(filter...))
	ccount, e := queryC.Cursor(nil)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}
	defer ccount.Close()

	totalDuration := 0.0
	totalTurbine := 0

	aggrData := []tk.M{}

	queryAggr := DB().Connection.NewQuery().From(new(EventDownHFD).TableName()).
		Aggr(dbox.AggrSum, "$duration", "duration").
		Group("turbine").Where(dbox.And(filter...))

	caggr, e := queryAggr.Cursor(nil)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}
	defer caggr.Close()
	e = caggr.Fetch(&aggrData, 0, false)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	for _, val := range aggrData {
		totalDuration += val.GetFloat64("duration")
	}
	totalTurbine = tk.SliceLen(aggrData)

	data := struct {
		Data          []EventDown
		Total         int
		TotalTurbine  int
		TotalDuration float64
	}{
		Data:          tmpResult,
		Total:         ccount.Count(),
		TotalTurbine:  totalTurbine,
		TotalDuration: totalDuration,
	}

	return helper.CreateResult(true, data, "success")
}

func (m *DataBrowserController) GetDowntimeEventvailDateHFD(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	DowntimeEventresults := make([]time.Time, 0)

	// Downtime Event Data
	for i := 0; i < 2; i++ {
		var arrsort []string
		if i == 0 {
			arrsort = append(arrsort, "timestart")
		} else {
			arrsort = append(arrsort, "-timestart")
		}

		query := DB().Connection.NewQuery().From(new(EventDownHFD).TableName()).Skip(0).Take(1)
		query = query.Order(arrsort...)

		csr, e := query.Cursor(nil)
		if e != nil {
			return helper.CreateResult(false, nil, e.Error())
		}
		defer csr.Close()

		Result := make([]EventDown, 0)
		e = csr.Fetch(&Result, 0, false)

		if e != nil {
			return helper.CreateResult(false, nil, e.Error())
		}

		for _, val := range Result {
			DowntimeEventresults = append(DowntimeEventresults, val.TimeStart.UTC())
		}
	}

	data := struct {
		DowntimeEvent []time.Time
	}{
		DowntimeEvent: DowntimeEventresults,
	}

	return helper.CreateResult(true, data, "success")
}

func (m *DataBrowserController) GenExcelDowntimeEventHFD(k *knot.WebContext) interface{} {

	k.Config.OutputType = knot.OutputJson

	var filter []*dbox.Filter

	p := new(helper.PayloadsDB)
	e := k.GetPayload(&p)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	tStart, _ := time.Parse("2006-01-02", p.DateStart.UTC().Format("2006-01-02"))
	tEnd, _ := time.Parse("2006-01-02 15:04:05", p.DateEnd.UTC().Format("2006-01-02")+" 23:59:59")
	turbine := p.Turbine

	var pathDownload string
	typeExcel := "DowntimeEventHFD"
	TimeCreate := time.Now().Format("2006-01-02_150405")
	CreateDateTime := typeExcel + TimeCreate

	if err := os.RemoveAll("web/assets/Excel/" + typeExcel + "/"); err != nil {
		tk.Println(err)
	}

	filter = append(filter, dbox.Ne("_id", ""))
	filter = append(filter, dbox.Gte("timestart", tStart))
	filter = append(filter, dbox.Lte("timestart", tEnd))
	if len(turbine) != 0 {
		filter = append(filter, dbox.In("turbine", turbine...))
	}

	query := DB().Connection.NewQuery().From(new(EventDownHFD).TableName()).Where(dbox.And(filter...))

	csr, e := query.Cursor(nil)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}
	defer csr.Close()

	tmpResult := make([]EventDown, 0)
	// results := make([]EventDown, 0)
	e = csr.Fetch(&tmpResult, 0, false)

	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	// for _, val := range tmpResult {
	// 	val.TimeStart = val.TimeStart.UTC()
	// 	results = append(results, val)
	// }
	//web/assets/Excel/

	if _, err := os.Stat("web/assets/Excel/" + typeExcel + "/"); os.IsNotExist(err) {
		os.MkdirAll("web/assets/Excel/"+typeExcel+"/", 0777)
	}
	turbineName, err := helper.GetTurbineNameList(p.Project)
	if err != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	DeserializeEventDown(tmpResult, 0, typeExcel, CreateDateTime, turbineName)
	pathDownload = "res/Excel/" + typeExcel + "/" + CreateDateTime + ".xlsx"
	// tk.Println(pathDownload)

	return helper.CreateResult(true, pathDownload, "success")
}
