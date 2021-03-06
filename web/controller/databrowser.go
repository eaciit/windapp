package controller

import (
	. "eaciit/wfdemo-git/library/core"
	// lh "eaciit/wfdemo-git/library/helper"
	. "eaciit/wfdemo-git/library/models"
	"eaciit/wfdemo-git/web/helper"
	"reflect"
	"strings"

	// "fmt"
	"os"
	"sort"
	// "strconv"
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

var (
	/*_amettower_label = []string{"V Hub WS 90m Avg", "V Hub WS 90m Std Dev", "V Ref WS 88m Avg", "V Ref WS 88m Std Dev",
		"V Tip WS 42m Avg", "V Tip WS 42m Std Dev", "D Hub WD 88m Avg", "D Hub WD 88m Std Dev", "D Ref WD 86m Avg",
		"D Ref WD 86m Std Dev", "T Hub & H Hub Humid 85m Avg", "T Hub & H Hub Humid 85m Std Dev", "T Ref & H Ref Humid 85.5m Avg", "T Ref & H Ref Humid 85.5m Std Dev",
		"T Hub & H Hub Temp 85.5m Avg", "T Hub & H Hub Temp 85.5m Std Dev", "T Ref & H Ref Temp 85.5 Avg", "T Ref & H Ref Temp 85.5 Std Dev", "Baro Air Pressure 85.5m Avg", "Baro Air Pressure 85.5m Std Dev",
	}

	_amettower_field = []string{"vhubws90mavg", "vhubws90mstddev", "vrefws88mavg", "vrefws88mstddev", "vtipws42mavg",
		"vtipws42mstddev", "dhubwd88mavg", "dhubwd88mstddev", "drefwd86mavg", "drefwd86mstddev",
		"thubhhubhumid855mavg", "thubhhubhumid855mstddev", "trefhrefhumid855mavg", "trefhrefhumid855mstddev", "thubhhubtemp855mavg",
		"thubhhubtemp855mstddev", "trefhreftemp855mavg", "trefhreftemp855mstddev", "baroairpress855mavg", "baroairpress855mstddev",
	}*/
	_amettower_label = []string{"D Hub WD 88m Avg", "D Ref WD 86m Avg", "T Hub & H Hub Temp 85.5m Avg", "T Ref & H Ref Temp 85.5 Avg",
		"V Hub WS 90m Avg", "V Hub WS 90m Std Dev", "V Ref WS 88m Avg", "V Ref WS 88m Std Dev", "V Tip WS 42m Avg", "V Tip WS 42m Std Dev"}

	_amettower_field = []string{"dhubwd88mavg", "drefwd86mavg", "thubhhubtemp855mavg", "trefhreftemp855mavg",
		"vhubws90mavg", "vhubws90mstddev", "vrefws88mavg", "vrefws88mstddev", "vtipws42mavg", "vtipws42mstddev"}
)

func GetCustomFieldList() []tk.M {
	atkm := []tk.M{}

	/*_ascadaoem_label := []string{"Ai Intern R Pid Angle Out", "Ai Intern I1", "Ai Intern I2",
		"Ai Dfig Torque Actual", "Ai Dr Tr Vib Value", "Ai Gear Oil Pressure", "Ai Hydr System Pressure", "Ai Intern Active Power",
		"Ai Intern Dfig Active Power Actual", "Ai Intern Nacelle Drill", "Ai Intern Nacelle Drill At North Pos Sensor", "Ai Intern Nacelle Pos",
		"Ai Intern Pitch Angle1", "Ai Intern Pitch Angle2", "Ai Intern Pitch Angle3", "Ai Intern Pitch Speed1", "Ai Intern Reactive Power",
		"Ai Intern Wind Direction", "Ai Intern Wind Speed", "Ai Intern Wind Speed Dif", "Ai Tower Vib Value Axial", "Ai Wind Speed1",
		"Ai Wind Speed2", "Ai Wind Vane1", "Ai Wind Vane2", "C Intern Speed Generator", "C Intern Speed Rotor", "Temp Bottom Control Section",
		"Temp Bottom Control Section Low", "Temp Bottom Power Section", "Temp Cabinet Top Box", "Temp Gearbox Hss De", "Temp Gear Box Hss Nde",
		"Temp Gear Box Ims De", "Temp Gear Box Ims Nde", "Temp Gear Oil Sump", "Temp Generator Bearing De", "Temp Generator Bearing Nde",
		"Temp Main Bearing", "Temp Nacelle", "Temp Outdoor", "Time Stamp", "Turbine",
	}*/

	/*_ascadaoem_field := []string{"ai_intern_r_pidangleout", "ai_intern_i1", "ai_intern_i2",
		"ai_dfig_torque_actual", "ai_drtrvibvalue", "ai_gearoilpressure", "ai_hydrsystempressure", "ai_intern_activpower",
		"ai_intern_dfig_active_power_actual", "ai_intern_nacelledrill", "ai_intern_nacelledrill_at_northpossensor", "ai_intern_nacellepos",
		"ai_intern_pitchangle1", "ai_intern_pitchangle2", "ai_intern_pitchangle3", "ai_intern_pitchspeed1", "ai_intern_reactivpower",
		"ai_intern_winddirection", "ai_intern_windspeed", "ai_intern_windspeeddif", "ai_towervibvalueaxial", "ai_windspeed1",
		"ai_windspeed2", "ai_windvane1", "ai_windvane2", "c_intern_speedgenerator", "c_intern_speedrotor", "temp_bottomcontrolsection",
		"temp_bottomcontrolsection_low", "temp_bottompowersection", "temp_cabinettopbox", "temp_gearbox_hss_de", "temp_gearbox_hss_nde",
		"temp_gearbox_ims_de", "temp_gearbox_ims_nde", "temp_gearoilsump", "temp_generatorbearing_de", "temp_generatorbearing_nde",
		"temp_mainbearing", "temp_nacelle", "temp_outdoor", "timestamp", "turbine",
	}*/
	_ascadaoem_label, _ascadaoem_field := GetScadaOEMHeader()

	for i, str := range _ascadaoem_field {
		tkm := tk.M{}.
			Set("_id", str).
			Set("label", _ascadaoem_label[i]).
			Set("source", "ScadaDataOEM")

		atkm = append(atkm, tkm)
	}

	/*for i, str := range _amettower_field {
		tkm := tk.M{}.
			Set("_id", str).
			Set("label", _amettower_label[i]).
			Set("source", "MetTower")

		atkm = append(atkm, tkm)
	}*/

	return atkm
}

//test commit
func GetHFDCustomFieldList() []tk.M {
	atkm := []tk.M{}

	// _ashfd_label := []string{"Fast ActivePower kW", "Fast ActivePower kW StdDev", "Fast ActivePower kW Min", "Fast ActivePower kW Max", "Fast ActivePower kW Count", "Fast WindSpeed ms", "Fast WindSpeed ms StdDev", "Fast WindSpeed ms Min", "Fast WindSpeed ms Max", "Fast WindSpeed ms Count", "Slow NacellePos", "Slow NacellePos StdDev", "Slow NacellePos Min", "Slow NacellePos Max", "Slow NacellePos Count", "Slow WindDirection", "Slow WindDirection StdDev", "Slow WindDirection Min", "Slow WindDirection Max", "Slow WindDirection Count", "Fast CurrentL3", "Fast CurrentL3 StdDev", "Fast CurrentL3 Min", "Fast CurrentL3 Max", "Fast CurrentL3 Count", "Fast CurrentL1", "Fast CurrentL1 StdDev", "Fast CurrentL1 Min", "Fast CurrentL1 Max", "Fast CurrentL1 Count", "Fast ActivePowerSetpoint kW", "Fast ActivePowerSetpoint kW StdDev", "Fast ActivePowerSetpoint kW Min", "Fast ActivePowerSetpoint kW Max", "Fast ActivePowerSetpoint kW Count", "Fast CurrentL2", "Fast CurrentL2 StdDev", "Fast CurrentL2 Min", "Fast CurrentL2 Max", "Fast CurrentL2 Count", "Fast DrTrVibValue", "Fast DrTrVibValue StdDev", "Fast DrTrVibValue Min", "Fast DrTrVibValue Max", "Fast DrTrVibValue Count", "Fast GenSpeed RPM", "Fast GenSpeed RPM StdDev", "Fast GenSpeed RPM Min", "Fast GenSpeed RPM Max", "Fast GenSpeed RPM Count", "Fast PitchAccuV1", "Fast PitchAccuV1 StdDev", "Fast PitchAccuV1 Min", "Fast PitchAccuV1 Max", "Fast PitchAccuV1 Count", "Fast PitchAngle", "Fast PitchAngle StdDev", "Fast PitchAngle Min", "Fast PitchAngle Max", "Fast PitchAngle Count", "Fast PitchAngle3", "Fast PitchAngle3 StdDev", "Fast PitchAngle3 Min", "Fast PitchAngle3 Max", "Fast PitchAngle3 Count", "Fast PitchAngle2", "Fast PitchAngle2 StdDev", "Fast PitchAngle2 Min", "Fast PitchAngle2 Max", "Fast PitchAngle2 Count", "Fast PitchConvCurrent1", "Fast PitchConvCurrent1 StdDev", "Fast PitchConvCurrent1 Min", "Fast PitchConvCurrent1 Max", "Fast PitchConvCurrent1 Count", "Fast PitchConvCurrent3", "Fast PitchConvCurrent3 StdDev", "Fast PitchConvCurrent3 Min", "Fast PitchConvCurrent3 Max", "Fast PitchConvCurrent3 Count", "Fast PitchConvCurrent2", "Fast PitchConvCurrent2 StdDev", "Fast PitchConvCurrent2 Min", "Fast PitchConvCurrent2 Max", "Fast PitchConvCurrent2 Count", "Fast PowerFactor", "Fast PowerFactor StdDev", "Fast PowerFactor Min", "Fast PowerFactor Max", "Fast PowerFactor Count", "Fast ReactivePowerSetpointPPC kVA", "Fast ReactivePowerSetpointPPC kVA StdDev", "Fast ReactivePowerSetpointPPC kVA Min", "Fast ReactivePowerSetpointPPC kVA Max", "Fast ReactivePowerSetpointPPC kVA Count", "Fast ReactivePower kVAr", "Fast ReactivePower kVAr StdDev", "Fast ReactivePower kVAr Min", "Fast ReactivePower kVAr Max", "Fast ReactivePower kVAr Count", "Fast RotorSpeed RPM", "Fast RotorSpeed RPM StdDev", "Fast RotorSpeed RPM Min", "Fast RotorSpeed RPM Max", "Fast RotorSpeed RPM Count", "Fast VoltageL1", "Fast VoltageL1 StdDev", "Fast VoltageL1 Min", "Fast VoltageL1 Max", "Fast VoltageL1 Count", "Fast VoltageL2", "Fast VoltageL2 StdDev", "Fast VoltageL2 Min", "Fast VoltageL2 Max", "Fast VoltageL2 Count", "Slow CapableCapacitiveReactPwr kVAr", "Slow CapableCapacitiveReactPwr kVAr StdDev", "Slow CapableCapacitiveReactPwr kVAr Min", "Slow CapableCapacitiveReactPwr kVAr Max", "Slow CapableCapacitiveReactPwr kVAr Count", "Slow CapableInductiveReactPwr kVAr", "Slow CapableInductiveReactPwr kVAr StdDev", "Slow CapableInductiveReactPwr kVAr Min", "Slow CapableInductiveReactPwr kVAr Max", "Slow CapableInductiveReactPwr kVAr Count", "Slow DateTime Sec", "Slow DateTime Sec StdDev", "Slow DateTime Sec Min", "Slow DateTime Sec Max", "Slow DateTime Sec Count", "Fast PitchAngle1", "Fast PitchAngle1 StdDev", "Fast PitchAngle1 Min", "Fast PitchAngle1 Max", "Fast PitchAngle1 Count", "Fast VoltageL3", "Fast VoltageL3 StdDev", "Fast VoltageL3 Min", "Fast VoltageL3 Max", "Fast VoltageL3 Count", "Slow CapableCapacitivePwrFactor", "Slow CapableCapacitivePwrFactor StdDev", "Slow CapableCapacitivePwrFactor Min", "Slow CapableCapacitivePwrFactor Max", "Slow CapableCapacitivePwrFactor Count", "Fast Total Production kWh", "Fast Total Production kWh StdDev", "Fast Total Production kWh Min", "Fast Total Production kWh Max", "Fast Total Production kWh Count", "Fast Total Prod Day kWh", "Fast Total Prod Day kWh StdDev", "Fast Total Prod Day kWh Min", "Fast Total Prod Day kWh Max", "Fast Total Prod Day kWh Count", "Fast Total Prod Month kWh", "Fast Total Prod Month kWh StdDev", "Fast Total Prod Month kWh Min", "Fast Total Prod Month kWh Max", "Fast Total Prod Month kWh Count", "Fast ActivePowerOutPWCSell kW", "Fast ActivePowerOutPWCSell kW StdDev", "Fast ActivePowerOutPWCSell kW Min", "Fast ActivePowerOutPWCSell kW Max", "Fast ActivePowerOutPWCSell kW Count", "Fast Frequency Hz", "Fast Frequency Hz StdDev", "Fast Frequency Hz Min", "Fast Frequency Hz Max", "Fast Frequency Hz Count", "Slow TempG1L2", "Slow TempG1L2 StdDev", "Slow TempG1L2 Min", "Slow TempG1L2 Max", "Slow TempG1L2 Count", "Slow TempG1L3", "Slow TempG1L3 StdDev", "Slow TempG1L3 Min", "Slow TempG1L3 Max", "Slow TempG1L3 Count", "Slow TempGearBoxHSSDE", "Slow TempGearBoxHSSDE StdDev", "Slow TempGearBoxHSSDE Min", "Slow TempGearBoxHSSDE Max", "Slow TempGearBoxHSSDE Count", "Slow TempGearBoxIMSNDE", "Slow TempGearBoxIMSNDE StdDev", "Slow TempGearBoxIMSNDE Min", "Slow TempGearBoxIMSNDE Max", "Slow TempGearBoxIMSNDE Count", "Slow TempOutdoor", "Slow TempOutdoor StdDev", "Slow TempOutdoor Min", "Slow TempOutdoor Max", "Slow TempOutdoor Count", "Fast PitchAccuV3", "Fast PitchAccuV3 StdDev", "Fast PitchAccuV3 Min", "Fast PitchAccuV3 Max", "Fast PitchAccuV3 Count", "Slow TotalTurbineActiveHours", "Slow TotalTurbineActiveHours StdDev", "Slow TotalTurbineActiveHours Min", "Slow TotalTurbineActiveHours Max", "Slow TotalTurbineActiveHours Count", "Slow TotalTurbineOKHours", "Slow TotalTurbineOKHours StdDev", "Slow TotalTurbineOKHours Min", "Slow TotalTurbineOKHours Max", "Slow TotalTurbineOKHours Count", "Slow TotalTurbineTimeAllHours", "Slow TotalTurbineTimeAllHours StdDev", "Slow TotalTurbineTimeAllHours Min", "Slow TotalTurbineTimeAllHours Max", "Slow TotalTurbineTimeAllHours Count", "Slow TempG1L1", "Slow TempG1L1 StdDev", "Slow TempG1L1 Min", "Slow TempG1L1 Max", "Slow TempG1L1 Count", "Slow TempGearBoxOilSump", "Slow TempGearBoxOilSump StdDev", "Slow TempGearBoxOilSump Min", "Slow TempGearBoxOilSump Max", "Slow TempGearBoxOilSump Count", "Fast PitchAccuV2", "Fast PitchAccuV2 StdDev", "Fast PitchAccuV2 Min", "Fast PitchAccuV2 Max", "Fast PitchAccuV2 Count", "Slow TotalGridOkHours", "Slow TotalGridOkHours StdDev", "Slow TotalGridOkHours Min", "Slow TotalGridOkHours Max", "Slow TotalGridOkHours Count", "Slow TotalActPowerOut kWh", "Slow TotalActPowerOut kWh StdDev", "Slow TotalActPowerOut kWh Min", "Slow TotalActPowerOut kWh Max", "Slow TotalActPowerOut kWh Count", "Fast YawService", "Fast YawService StdDev", "Fast YawService Min", "Fast YawService Max", "Fast YawService Count", "Fast YawAngle", "Fast YawAngle StdDev", "Fast YawAngle Min", "Fast YawAngle Max", "Fast YawAngle Count", "Slow CapableInductivePwrFactor", "Slow CapableInductivePwrFactor StdDev", "Slow CapableInductivePwrFactor Min", "Slow CapableInductivePwrFactor Max", "Slow CapableInductivePwrFactor Count", "Slow TempGearBoxHSSNDE", "Slow TempGearBoxHSSNDE StdDev", "Slow TempGearBoxHSSNDE Min", "Slow TempGearBoxHSSNDE Max", "Slow TempGearBoxHSSNDE Count", "Slow TempHubBearing", "Slow TempHubBearing StdDev", "Slow TempHubBearing Min", "Slow TempHubBearing Max", "Slow TempHubBearing Count", "Slow TotalG1ActiveHours", "Slow TotalG1ActiveHours StdDev", "Slow TotalG1ActiveHours Min", "Slow TotalG1ActiveHours Max", "Slow TotalG1ActiveHours Count", "Slow TotalActPowerOutG1 kWh", "Slow TotalActPowerOutG1 kWh StdDev", "Slow TotalActPowerOutG1 kWh Min", "Slow TotalActPowerOutG1 kWh Max", "Slow TotalActPowerOutG1 kWh Count", "Slow TotalReactPowerInG1 kVArh", "Slow TotalReactPowerInG1 kVArh StdDev", "Slow TotalReactPowerInG1 kVArh Min", "Slow TotalReactPowerInG1 kVArh Max", "Slow TotalReactPowerInG1 kVArh Count", "Slow NacelleDrill", "Slow NacelleDrill StdDev", "Slow NacelleDrill Min", "Slow NacelleDrill Max", "Slow NacelleDrill Count", "Slow TempGearBoxIMSDE", "Slow TempGearBoxIMSDE StdDev", "Slow TempGearBoxIMSDE Min", "Slow TempGearBoxIMSDE Max", "Slow TempGearBoxIMSDE Count", "Fast Total Operating hrs", "Fast Total Operating hrs StdDev", "Fast Total Operating hrs Min", "Fast Total Operating hrs Max", "Fast Total Operating hrs Count", "Slow TempNacelle", "Slow TempNacelle StdDev", "Slow TempNacelle Min", "Slow TempNacelle Max", "Slow TempNacelle Count", "Fast Total Grid OK hrs", "Fast Total Grid OK hrs StdDev", "Fast Total Grid OK hrs Min", "Fast Total Grid OK hrs Max", "Fast Total Grid OK hrs Count", "Fast Total WTG OK hrs", "Fast Total WTG OK hrs StdDev", "Fast Total WTG OK hrs Min", "Fast Total WTG OK hrs Max", "Fast Total WTG OK hrs Count", "Slow TempCabinetTopBox", "Slow TempCabinetTopBox StdDev", "Slow TempCabinetTopBox Min", "Slow TempCabinetTopBox Max", "Slow TempCabinetTopBox Count", "Slow TempGeneratorBearingNDE", "Slow TempGeneratorBearingNDE StdDev", "Slow TempGeneratorBearingNDE Min", "Slow TempGeneratorBearingNDE Max", "Slow TempGeneratorBearingNDE Count", "Fast Total Access hrs", "Fast Total Access hrs StdDev", "Fast Total Access hrs Min", "Fast Total Access hrs Max", "Fast Total Access hrs Count", "Slow TempBottomPowerSection", "Slow TempBottomPowerSection StdDev", "Slow TempBottomPowerSection Min", "Slow TempBottomPowerSection Max", "Slow TempBottomPowerSection Count", "Slow TempGeneratorBearingDE", "Slow TempGeneratorBearingDE StdDev", "Slow TempGeneratorBearingDE Min", "Slow TempGeneratorBearingDE Max", "Slow TempGeneratorBearingDE Count", "Slow TotalReactPowerIn kVArh", "Slow TotalReactPowerIn kVArh StdDev", "Slow TotalReactPowerIn kVArh Min", "Slow TotalReactPowerIn kVArh Max", "Slow TotalReactPowerIn kVArh Count", "Slow TempBottomControlSection", "Slow TempBottomControlSection StdDev", "Slow TempBottomControlSection Min", "Slow TempBottomControlSection Max", "Slow TempBottomControlSection Count", "Slow TempConv1", "Slow TempConv1 StdDev", "Slow TempConv1 Min", "Slow TempConv1 Max", "Slow TempConv1 Count", "Fast ActivePowerRated kW", "Fast ActivePowerRated kW StdDev", "Fast ActivePowerRated kW Min", "Fast ActivePowerRated kW Max", "Fast ActivePowerRated kW Count", "Fast NodeIP", "Fast NodeIP StdDev", "Fast NodeIP Min", "Fast NodeIP Max", "Fast NodeIP Count", "Fast PitchSpeed1", "Fast PitchSpeed1 StdDev", "Fast PitchSpeed1 Min", "Fast PitchSpeed1 Max", "Fast PitchSpeed1 Count", "Slow CFCardSize", "Slow CFCardSize StdDev", "Slow CFCardSize Min", "Slow CFCardSize Max", "Slow CFCardSize Count", "Slow CPU Number", "Slow CPU Number StdDev", "Slow CPU Number Min", "Slow CPU Number Max", "Slow CPU Number Count", "Slow CFCardSpaceLeft", "Slow CFCardSpaceLeft StdDev", "Slow CFCardSpaceLeft Min", "Slow CFCardSpaceLeft Max", "Slow CFCardSpaceLeft Count", "Slow TempBottomCapSection", "Slow TempBottomCapSection StdDev", "Slow TempBottomCapSection Min", "Slow TempBottomCapSection Max", "Slow TempBottomCapSection Count", "Slow RatedPower", "Slow RatedPower StdDev", "Slow RatedPower Min", "Slow RatedPower Max", "Slow RatedPower Count", "Slow TempConv3", "Slow TempConv3 StdDev", "Slow TempConv3 Min", "Slow TempConv3 Max", "Slow TempConv3 Count", "Slow TempConv2", "Slow TempConv2 StdDev", "Slow TempConv2 Min", "Slow TempConv2 Max", "Slow TempConv2 Count", "Slow TotalActPowerIn kWh", "Slow TotalActPowerIn kWh StdDev", "Slow TotalActPowerIn kWh Min", "Slow TotalActPowerIn kWh Max", "Slow TotalActPowerIn kWh Count", "Slow TotalActPowerInG1 kWh", "Slow TotalActPowerInG1 kWh StdDev", "Slow TotalActPowerInG1 kWh Min", "Slow TotalActPowerInG1 kWh Max", "Slow TotalActPowerInG1 kWh Count", "Slow TotalActPowerInG2 kWh", "Slow TotalActPowerInG2 kWh StdDev", "Slow TotalActPowerInG2 kWh Min", "Slow TotalActPowerInG2 kWh Max", "Slow TotalActPowerInG2 kWh Count", "Slow TotalActPowerOutG2 kWh", "Slow TotalActPowerOutG2 kWh StdDev", "Slow TotalActPowerOutG2 kWh Min", "Slow TotalActPowerOutG2 kWh Max", "Slow TotalActPowerOutG2 kWh Count", "Slow TotalG2ActiveHours", "Slow TotalG2ActiveHours StdDev", "Slow TotalG2ActiveHours Min", "Slow TotalG2ActiveHours Max", "Slow TotalG2ActiveHours Count", "Slow TotalReactPowerInG2 kVArh", "Slow TotalReactPowerInG2 kVArh StdDev", "Slow TotalReactPowerInG2 kVArh Min", "Slow TotalReactPowerInG2 kVArh Max", "Slow TotalReactPowerInG2 kVArh Count", "Slow TotalReactPowerOut kVArh", "Slow TotalReactPowerOut kVArh StdDev", "Slow TotalReactPowerOut kVArh Min", "Slow TotalReactPowerOut kVArh Max", "Slow TotalReactPowerOut kVArh Count", "Slow UTCoffset int", "Slow UTCoffset int StdDev", "Slow UTCoffset int Min", "Slow UTCoffset int Max", "Slow UTCoffset int Count", "Time Stamp", "Turbine"}

	// _ashfd_field := []string{"Fast_ActivePower_kW", "Fast_ActivePower_kW_StdDev", "Fast_ActivePower_kW_Min", "Fast_ActivePower_kW_Max", "Fast_ActivePower_kW_Count", "Fast_WindSpeed_ms", "Fast_WindSpeed_ms_StdDev", "Fast_WindSpeed_ms_Min", "Fast_WindSpeed_ms_Max", "Fast_WindSpeed_ms_Count", "Slow_NacellePos", "Slow_NacellePos_StdDev", "Slow_NacellePos_Min", "Slow_NacellePos_Max", "Slow_NacellePos_Count", "Slow_WindDirection", "Slow_WindDirection_StdDev", "Slow_WindDirection_Min", "Slow_WindDirection_Max", "Slow_WindDirection_Count", "Fast_CurrentL3", "Fast_CurrentL3_StdDev", "Fast_CurrentL3_Min", "Fast_CurrentL3_Max", "Fast_CurrentL3_Count", "Fast_CurrentL1", "Fast_CurrentL1_StdDev", "Fast_CurrentL1_Min", "Fast_CurrentL1_Max", "Fast_CurrentL1_Count", "Fast_ActivePowerSetpoint_kW", "Fast_ActivePowerSetpoint_kW_StdDev", "Fast_ActivePowerSetpoint_kW_Min", "Fast_ActivePowerSetpoint_kW_Max", "Fast_ActivePowerSetpoint_kW_Count", "Fast_CurrentL2", "Fast_CurrentL2_StdDev", "Fast_CurrentL2_Min", "Fast_CurrentL2_Max", "Fast_CurrentL2_Count", "Fast_DrTrVibValue", "Fast_DrTrVibValue_StdDev", "Fast_DrTrVibValue_Min", "Fast_DrTrVibValue_Max", "Fast_DrTrVibValue_Count", "Fast_GenSpeed_RPM", "Fast_GenSpeed_RPM_StdDev", "Fast_GenSpeed_RPM_Min", "Fast_GenSpeed_RPM_Max", "Fast_GenSpeed_RPM_Count", "Fast_PitchAccuV1", "Fast_PitchAccuV1_StdDev", "Fast_PitchAccuV1_Min", "Fast_PitchAccuV1_Max", "Fast_PitchAccuV1_Count", "Fast_PitchAngle", "Fast_PitchAngle_StdDev", "Fast_PitchAngle_Min", "Fast_PitchAngle_Max", "Fast_PitchAngle_Count", "Fast_PitchAngle3", "Fast_PitchAngle3_StdDev", "Fast_PitchAngle3_Min", "Fast_PitchAngle3_Max", "Fast_PitchAngle3_Count", "Fast_PitchAngle2", "Fast_PitchAngle2_StdDev", "Fast_PitchAngle2_Min", "Fast_PitchAngle2_Max", "Fast_PitchAngle2_Count", "Fast_PitchConvCurrent1", "Fast_PitchConvCurrent1_StdDev", "Fast_PitchConvCurrent1_Min", "Fast_PitchConvCurrent1_Max", "Fast_PitchConvCurrent1_Count", "Fast_PitchConvCurrent3", "Fast_PitchConvCurrent3_StdDev", "Fast_PitchConvCurrent3_Min", "Fast_PitchConvCurrent3_Max", "Fast_PitchConvCurrent3_Count", "Fast_PitchConvCurrent2", "Fast_PitchConvCurrent2_StdDev", "Fast_PitchConvCurrent2_Min", "Fast_PitchConvCurrent2_Max", "Fast_PitchConvCurrent2_Count", "Fast_PowerFactor", "Fast_PowerFactor_StdDev", "Fast_PowerFactor_Min", "Fast_PowerFactor_Max", "Fast_PowerFactor_Count", "Fast_ReactivePowerSetpointPPC_kVA", "Fast_ReactivePowerSetpointPPC_kVA_StdDev", "Fast_ReactivePowerSetpointPPC_kVA_Min", "Fast_ReactivePowerSetpointPPC_kVA_Max", "Fast_ReactivePowerSetpointPPC_kVA_Count", "Fast_ReactivePower_kVAr", "Fast_ReactivePower_kVAr_StdDev", "Fast_ReactivePower_kVAr_Min", "Fast_ReactivePower_kVAr_Max", "Fast_ReactivePower_kVAr_Count", "Fast_RotorSpeed_RPM", "Fast_RotorSpeed_RPM_StdDev", "Fast_RotorSpeed_RPM_Min", "Fast_RotorSpeed_RPM_Max", "Fast_RotorSpeed_RPM_Count", "Fast_VoltageL1", "Fast_VoltageL1_StdDev", "Fast_VoltageL1_Min", "Fast_VoltageL1_Max", "Fast_VoltageL1_Count", "Fast_VoltageL2", "Fast_VoltageL2_StdDev", "Fast_VoltageL2_Min", "Fast_VoltageL2_Max", "Fast_VoltageL2_Count", "Slow_CapableCapacitiveReactPwr_kVAr", "Slow_CapableCapacitiveReactPwr_kVAr_StdDev", "Slow_CapableCapacitiveReactPwr_kVAr_Min", "Slow_CapableCapacitiveReactPwr_kVAr_Max", "Slow_CapableCapacitiveReactPwr_kVAr_Count", "Slow_CapableInductiveReactPwr_kVAr", "Slow_CapableInductiveReactPwr_kVAr_StdDev", "Slow_CapableInductiveReactPwr_kVAr_Min", "Slow_CapableInductiveReactPwr_kVAr_Max", "Slow_CapableInductiveReactPwr_kVAr_Count", "Slow_DateTime_Sec", "Slow_DateTime_Sec_StdDev", "Slow_DateTime_Sec_Min", "Slow_DateTime_Sec_Max", "Slow_DateTime_Sec_Count", "Fast_PitchAngle1", "Fast_PitchAngle1_StdDev", "Fast_PitchAngle1_Min", "Fast_PitchAngle1_Max", "Fast_PitchAngle1_Count", "Fast_VoltageL3", "Fast_VoltageL3_StdDev", "Fast_VoltageL3_Min", "Fast_VoltageL3_Max", "Fast_VoltageL3_Count", "Slow_CapableCapacitivePwrFactor", "Slow_CapableCapacitivePwrFactor_StdDev", "Slow_CapableCapacitivePwrFactor_Min", "Slow_CapableCapacitivePwrFactor_Max", "Slow_CapableCapacitivePwrFactor_Count", "Fast_Total_Production_kWh", "Fast_Total_Production_kWh_StdDev", "Fast_Total_Production_kWh_Min", "Fast_Total_Production_kWh_Max", "Fast_Total_Production_kWh_Count", "Fast_Total_Prod_Day_kWh", "Fast_Total_Prod_Day_kWh_StdDev", "Fast_Total_Prod_Day_kWh_Min", "Fast_Total_Prod_Day_kWh_Max", "Fast_Total_Prod_Day_kWh_Count", "Fast_Total_Prod_Month_kWh", "Fast_Total_Prod_Month_kWh_StdDev", "Fast_Total_Prod_Month_kWh_Min", "Fast_Total_Prod_Month_kWh_Max", "Fast_Total_Prod_Month_kWh_Count", "Fast_ActivePowerOutPWCSell_kW", "Fast_ActivePowerOutPWCSell_kW_StdDev", "Fast_ActivePowerOutPWCSell_kW_Min", "Fast_ActivePowerOutPWCSell_kW_Max", "Fast_ActivePowerOutPWCSell_kW_Count", "Fast_Frequency_Hz", "Fast_Frequency_Hz_StdDev", "Fast_Frequency_Hz_Min", "Fast_Frequency_Hz_Max", "Fast_Frequency_Hz_Count", "Slow_TempG1L2", "Slow_TempG1L2_StdDev", "Slow_TempG1L2_Min", "Slow_TempG1L2_Max", "Slow_TempG1L2_Count", "Slow_TempG1L3", "Slow_TempG1L3_StdDev", "Slow_TempG1L3_Min", "Slow_TempG1L3_Max", "Slow_TempG1L3_Count", "Slow_TempGearBoxHSSDE", "Slow_TempGearBoxHSSDE_StdDev", "Slow_TempGearBoxHSSDE_Min", "Slow_TempGearBoxHSSDE_Max", "Slow_TempGearBoxHSSDE_Count", "Slow_TempGearBoxIMSNDE", "Slow_TempGearBoxIMSNDE_StdDev", "Slow_TempGearBoxIMSNDE_Min", "Slow_TempGearBoxIMSNDE_Max", "Slow_TempGearBoxIMSNDE_Count", "Slow_TempOutdoor", "Slow_TempOutdoor_StdDev", "Slow_TempOutdoor_Min", "Slow_TempOutdoor_Max", "Slow_TempOutdoor_Count", "Fast_PitchAccuV3", "Fast_PitchAccuV3_StdDev", "Fast_PitchAccuV3_Min", "Fast_PitchAccuV3_Max", "Fast_PitchAccuV3_Count", "Slow_TotalTurbineActiveHours", "Slow_TotalTurbineActiveHours_StdDev", "Slow_TotalTurbineActiveHours_Min", "Slow_TotalTurbineActiveHours_Max", "Slow_TotalTurbineActiveHours_Count", "Slow_TotalTurbineOKHours", "Slow_TotalTurbineOKHours_StdDev", "Slow_TotalTurbineOKHours_Min", "Slow_TotalTurbineOKHours_Max", "Slow_TotalTurbineOKHours_Count", "Slow_TotalTurbineTimeAllHours", "Slow_TotalTurbineTimeAllHours_StdDev", "Slow_TotalTurbineTimeAllHours_Min", "Slow_TotalTurbineTimeAllHours_Max", "Slow_TotalTurbineTimeAllHours_Count", "Slow_TempG1L1", "Slow_TempG1L1_StdDev", "Slow_TempG1L1_Min", "Slow_TempG1L1_Max", "Slow_TempG1L1_Count", "Slow_TempGearBoxOilSump", "Slow_TempGearBoxOilSump_StdDev", "Slow_TempGearBoxOilSump_Min", "Slow_TempGearBoxOilSump_Max", "Slow_TempGearBoxOilSump_Count", "Fast_PitchAccuV2", "Fast_PitchAccuV2_StdDev", "Fast_PitchAccuV2_Min", "Fast_PitchAccuV2_Max", "Fast_PitchAccuV2_Count", "Slow_TotalGridOkHours", "Slow_TotalGridOkHours_StdDev", "Slow_TotalGridOkHours_Min", "Slow_TotalGridOkHours_Max", "Slow_TotalGridOkHours_Count", "Slow_TotalActPowerOut_kWh", "Slow_TotalActPowerOut_kWh_StdDev", "Slow_TotalActPowerOut_kWh_Min", "Slow_TotalActPowerOut_kWh_Max", "Slow_TotalActPowerOut_kWh_Count", "Fast_YawService", "Fast_YawService_StdDev", "Fast_YawService_Min", "Fast_YawService_Max", "Fast_YawService_Count", "Fast_YawAngle", "Fast_YawAngle_StdDev", "Fast_YawAngle_Min", "Fast_YawAngle_Max", "Fast_YawAngle_Count", "Slow_CapableInductivePwrFactor", "Slow_CapableInductivePwrFactor_StdDev", "Slow_CapableInductivePwrFactor_Min", "Slow_CapableInductivePwrFactor_Max", "Slow_CapableInductivePwrFactor_Count", "Slow_TempGearBoxHSSNDE", "Slow_TempGearBoxHSSNDE_StdDev", "Slow_TempGearBoxHSSNDE_Min", "Slow_TempGearBoxHSSNDE_Max", "Slow_TempGearBoxHSSNDE_Count", "Slow_TempHubBearing", "Slow_TempHubBearing_StdDev", "Slow_TempHubBearing_Min", "Slow_TempHubBearing_Max", "Slow_TempHubBearing_Count", "Slow_TotalG1ActiveHours", "Slow_TotalG1ActiveHours_StdDev", "Slow_TotalG1ActiveHours_Min", "Slow_TotalG1ActiveHours_Max", "Slow_TotalG1ActiveHours_Count", "Slow_TotalActPowerOutG1_kWh", "Slow_TotalActPowerOutG1_kWh_StdDev", "Slow_TotalActPowerOutG1_kWh_Min", "Slow_TotalActPowerOutG1_kWh_Max", "Slow_TotalActPowerOutG1_kWh_Count", "Slow_TotalReactPowerInG1_kVArh", "Slow_TotalReactPowerInG1_kVArh_StdDev", "Slow_TotalReactPowerInG1_kVArh_Min", "Slow_TotalReactPowerInG1_kVArh_Max", "Slow_TotalReactPowerInG1_kVArh_Count", "Slow_NacelleDrill", "Slow_NacelleDrill_StdDev", "Slow_NacelleDrill_Min", "Slow_NacelleDrill_Max", "Slow_NacelleDrill_Count", "Slow_TempGearBoxIMSDE", "Slow_TempGearBoxIMSDE_StdDev", "Slow_TempGearBoxIMSDE_Min", "Slow_TempGearBoxIMSDE_Max", "Slow_TempGearBoxIMSDE_Count", "Fast_Total_Operating_hrs", "Fast_Total_Operating_hrs_StdDev", "Fast_Total_Operating_hrs_Min", "Fast_Total_Operating_hrs_Max", "Fast_Total_Operating_hrs_Count", "Slow_TempNacelle", "Slow_TempNacelle_StdDev", "Slow_TempNacelle_Min", "Slow_TempNacelle_Max", "Slow_TempNacelle_Count", "Fast_Total_Grid_OK_hrs", "Fast_Total_Grid_OK_hrs_StdDev", "Fast_Total_Grid_OK_hrs_Min", "Fast_Total_Grid_OK_hrs_Max", "Fast_Total_Grid_OK_hrs_Count", "Fast_Total_WTG_OK_hrs", "Fast_Total_WTG_OK_hrs_StdDev", "Fast_Total_WTG_OK_hrs_Min", "Fast_Total_WTG_OK_hrs_Max", "Fast_Total_WTG_OK_hrs_Count", "Slow_TempCabinetTopBox", "Slow_TempCabinetTopBox_StdDev", "Slow_TempCabinetTopBox_Min", "Slow_TempCabinetTopBox_Max", "Slow_TempCabinetTopBox_Count", "Slow_TempGeneratorBearingNDE", "Slow_TempGeneratorBearingNDE_StdDev", "Slow_TempGeneratorBearingNDE_Min", "Slow_TempGeneratorBearingNDE_Max", "Slow_TempGeneratorBearingNDE_Count", "Fast_Total_Access_hrs", "Fast_Total_Access_hrs_StdDev", "Fast_Total_Access_hrs_Min", "Fast_Total_Access_hrs_Max", "Fast_Total_Access_hrs_Count", "Slow_TempBottomPowerSection", "Slow_TempBottomPowerSection_StdDev", "Slow_TempBottomPowerSection_Min", "Slow_TempBottomPowerSection_Max", "Slow_TempBottomPowerSection_Count", "Slow_TempGeneratorBearingDE", "Slow_TempGeneratorBearingDE_StdDev", "Slow_TempGeneratorBearingDE_Min", "Slow_TempGeneratorBearingDE_Max", "Slow_TempGeneratorBearingDE_Count", "Slow_TotalReactPowerIn_kVArh", "Slow_TotalReactPowerIn_kVArh_StdDev", "Slow_TotalReactPowerIn_kVArh_Min", "Slow_TotalReactPowerIn_kVArh_Max", "Slow_TotalReactPowerIn_kVArh_Count", "Slow_TempBottomControlSection", "Slow_TempBottomControlSection_StdDev", "Slow_TempBottomControlSection_Min", "Slow_TempBottomControlSection_Max", "Slow_TempBottomControlSection_Count", "Slow_TempConv1", "Slow_TempConv1_StdDev", "Slow_TempConv1_Min", "Slow_TempConv1_Max", "Slow_TempConv1_Count", "Fast_ActivePowerRated_kW", "Fast_ActivePowerRated_kW_StdDev", "Fast_ActivePowerRated_kW_Min", "Fast_ActivePowerRated_kW_Max", "Fast_ActivePowerRated_kW_Count", "Fast_NodeIP", "Fast_NodeIP_StdDev", "Fast_NodeIP_Min", "Fast_NodeIP_Max", "Fast_NodeIP_Count", "Fast_PitchSpeed1", "Fast_PitchSpeed1_StdDev", "Fast_PitchSpeed1_Min", "Fast_PitchSpeed1_Max", "Fast_PitchSpeed1_Count", "Slow_CFCardSize", "Slow_CFCardSize_StdDev", "Slow_CFCardSize_Min", "Slow_CFCardSize_Max", "Slow_CFCardSize_Count", "Slow_CPU_Number", "Slow_CPU_Number_StdDev", "Slow_CPU_Number_Min", "Slow_CPU_Number_Max", "Slow_CPU_Number_Count", "Slow_CFCardSpaceLeft", "Slow_CFCardSpaceLeft_StdDev", "Slow_CFCardSpaceLeft_Min", "Slow_CFCardSpaceLeft_Max", "Slow_CFCardSpaceLeft_Count", "Slow_TempBottomCapSection", "Slow_TempBottomCapSection_StdDev", "Slow_TempBottomCapSection_Min", "Slow_TempBottomCapSection_Max", "Slow_TempBottomCapSection_Count", "Slow_RatedPower", "Slow_RatedPower_StdDev", "Slow_RatedPower_Min", "Slow_RatedPower_Max", "Slow_RatedPower_Count", "Slow_TempConv3", "Slow_TempConv3_StdDev", "Slow_TempConv3_Min", "Slow_TempConv3_Max", "Slow_TempConv3_Count", "Slow_TempConv2", "Slow_TempConv2_StdDev", "Slow_TempConv2_Min", "Slow_TempConv2_Max", "Slow_TempConv2_Count", "Slow_TotalActPowerIn_kWh", "Slow_TotalActPowerIn_kWh_StdDev", "Slow_TotalActPowerIn_kWh_Min", "Slow_TotalActPowerIn_kWh_Max", "Slow_TotalActPowerIn_kWh_Count", "Slow_TotalActPowerInG1_kWh", "Slow_TotalActPowerInG1_kWh_StdDev", "Slow_TotalActPowerInG1_kWh_Min", "Slow_TotalActPowerInG1_kWh_Max", "Slow_TotalActPowerInG1_kWh_Count", "Slow_TotalActPowerInG2_kWh", "Slow_TotalActPowerInG2_kWh_StdDev", "Slow_TotalActPowerInG2_kWh_Min", "Slow_TotalActPowerInG2_kWh_Max", "Slow_TotalActPowerInG2_kWh_Count", "Slow_TotalActPowerOutG2_kWh", "Slow_TotalActPowerOutG2_kWh_StdDev", "Slow_TotalActPowerOutG2_kWh_Min", "Slow_TotalActPowerOutG2_kWh_Max", "Slow_TotalActPowerOutG2_kWh_Count", "Slow_TotalG2ActiveHours", "Slow_TotalG2ActiveHours_StdDev", "Slow_TotalG2ActiveHours_Min", "Slow_TotalG2ActiveHours_Max", "Slow_TotalG2ActiveHours_Count", "Slow_TotalReactPowerInG2_kVArh", "Slow_TotalReactPowerInG2_kVArh_StdDev", "Slow_TotalReactPowerInG2_kVArh_Min", "Slow_TotalReactPowerInG2_kVArh_Max", "Slow_TotalReactPowerInG2_kVArh_Count", "Slow_TotalReactPowerOut_kVArh", "Slow_TotalReactPowerOut_kVArh_StdDev", "Slow_TotalReactPowerOut_kVArh_Min", "Slow_TotalReactPowerOut_kVArh_Max", "Slow_TotalReactPowerOut_kVArh_Count", "Slow_UTCoffset_int", "Slow_UTCoffset_int_StdDev", "Slow_UTCoffset_int_Min", "Slow_UTCoffset_int_Max", "Slow_UTCoffset_int_Count", "timestamp", "turbine"}

	// _ashfd_label, _ashfd_field := GetScadaHFDHeader()
	// for i, str := range _ashfd_field {
	// 	tkm := tk.M{}.
	// 		Set("_id", strings.ToLower(str)).
	// 		Set("label", _ashfd_label[i]).
	// 		Set("source", "ScadaDataHFD")

	// 	atkm = append(atkm, tkm)
	// }
	csr, e := DB().Connection.NewQuery().From("ref_databrowsertag").
		Order("order").Cursor(nil)
	defer csr.Close()
	if e != nil {
		tk.Println(e.Error())
	}

	minMaxTagList := map[string]bool{
		"windspeed_ms":       true,
		"activepower_kw":     true,
		"reactivepower_kvar": true,
		"rotorspeed_rpm":     true,
		"genspeed_rpm":       true,
		"pitchangle":         true,
		"pitchangle1":        true,
		"pitchangle2":        true,
		"pitchangle3":        true,
		"gridfrequencyhz":    true,
		"gridppvphaseab":     true,
		"gridppvphasebc":     true,
		"gridppvphaseca":     true,
		"gridcurrent":        true,
		"nacellepos":         true,
		"nacelledeviation":   true,
		"winddirection":      true,
	}
	minMaxList := []string{"min", "max", "stddev"}
	lastOrderPerProject := map[string]int{}

	_data := tk.M{}
	additionalData := []tk.M{}
	for {
		_data = tk.M{}
		e = csr.Fetch(&_data, 1, false)
		if e != nil {
			break
		}
		isTemp := strings.Contains(strings.ToLower(_data.GetString("realtimefield")), "temp")
		isEnable := _data.Get("enable", false).(bool)
		if isTemp || isEnable { /* jika tags temperature atau tags yang enable */
			idLower := strings.ToLower(_data.GetString("realtimefield"))
			atkm = append(atkm, tk.M{
				"_id":         idLower,
				"label":       _data.GetString("label"),
				"order":       _data.GetInt("order"),
				"projectname": _data.GetString("projectname"),
				"source":      _data.GetString("source"),
			})
			lastOrderPerProject[_data.GetString("projectname")] = _data.GetInt("order")
			/* kasih min, max dan stddev buat temperature tags */
			if isTemp || minMaxTagList[idLower] { /* jika tags temperature atau tags tertentu, tambahkan min, max, stddev */
				for _, minMaxVal := range minMaxList {
					additionalData = append(additionalData, tk.M{
						"_id":         idLower + "_" + minMaxVal,
						"label":       _data.GetString("label") + " " + strings.Title(minMaxVal),
						"order":       -999,
						"projectname": _data.GetString("projectname"),
						"source":      _data.GetString("source"),
					})
				}
			}
		}
	}

	for _, val := range additionalData { /* masukkan order yang benar (agak gak penting sih) */
		atkm = append(atkm, tk.M{
			"_id":         val.GetString("_id"),
			"label":       val.GetString("label"),
			"order":       lastOrderPerProject[val.GetString("projectname")] + 1,
			"projectname": val.GetString("projectname"),
			"source":      val.GetString("source"),
		})
		lastOrderPerProject[val.GetString("projectname")]++
	}

	/*startIndex := len(atkm)
	for i, str := range _amettower_field {
		startIndex++
		tkm := tk.M{}.
			Set("_id", str).
			Set("label", _amettower_label[i]).
			Set("source", "MetTower").
			Set("order", startIndex)

		atkm = append(atkm, tkm)
	}*/

	return atkm
}

func CheckData(tmpResult []tk.M, filter []*dbox.Filter, header map[string]string, tipe string) (result []tk.M, err error) {
	project := ""
	for _, val := range filter {
		if val.Field == "projectname" || val.Field == "project" {
			project = tk.ToString(val.Value)
		}
	}
	turbineName, err := helper.GetTurbineNameList(project)
	if err != nil {
		return
	}
	floatString := ""
	lowerField := ""
	for idx, val := range tmpResult {
		for field, dataType := range header {
			lowerField = strings.ToLower(field)
			if val.Has(lowerField) {
				switch lowerField {
				case "timestamp", "timestamputc", "timestart", "timeend":
					if tipe != "custom" {
						tmpResult[idx].Set(field, val.Get(lowerField, time.Time{}).(time.Time).UTC().Format("2006-01-02 15:04:05"))
					}
				case "turbine":
					tmpResult[idx].Set(field, turbineName[val.GetString(lowerField)])
				default:
					if dataType == "float64" {
						floatString = tk.Sprintf("%.2f", val.GetFloat64(lowerField))
						if val.GetFloat64(lowerField) == -999999 {
							floatString = "-"
						}
						if floatString != "-" && lowerField != "duration" {
							floatString = FormatThousandSeparator(floatString)
						}
						tmpResult[idx].Set(field, floatString)
					} else if dataType == "int" {
						floatString = tk.Sprintf("%d", val.GetInt(lowerField))
						if val.GetInt(lowerField) == -999999 {
							floatString = "-"
						}
						tmpResult[idx].Set(field, floatString)
					} else if dataType == "string" {
						tmpResult[idx].Set(field, val.GetString(lowerField))
					} else if dataType == "bool" {
						tmpResult[idx].Set(field, val[lowerField].(bool))
					}
				}
			}
			if tipe != "custom" {
				if tmpResult[idx].Has(lowerField) {
					tmpResult[idx].Unset(lowerField)
				}
			}
		}
	}
	result = tmpResult
	return
}

// GET DATA

func (m *DataBrowserController) GetDataBrowserList(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	p := new(helper.Payloads)
	e := k.GetPayload(&p)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}
	p.Misc.Set("knot_data", k)
	filter, _ := p.ParseFilter()
	tipe := p.Misc.GetString("tipe")
	needTotalTurbine := p.Misc["needtotalturbine"].(bool)
	tablename := new(ScadaDataOEM).TableName()
	var reflectVal reflect.Value
	var database dbox.IConnection

	switch tipe {
	case "scadaoem":
		obj := ScadaDataOEM{}
		tablename = obj.TableName()
		reflectVal = reflect.Indirect(reflect.ValueOf(obj))
		database = DB().Connection
	case "scadahfd":
		obj := Scada10Min{}
		// tablename = obj.TableName()
		tablename = "Scada10MinHFD"
		reflectVal = reflect.Indirect(reflect.ValueOf(obj))
		database = DB().Connection
	case "met":
		obj := MetTower{}
		tablename = obj.TableName()
		reflectVal = reflect.Indirect(reflect.ValueOf(obj))
		database = DB().Connection
	case "eventraw":
		obj := EventRaw{}
		tablename = obj.TableName()
		reflectVal = reflect.Indirect(reflect.ValueOf(obj))
		database = DB().Connection
	case "eventdown":
		obj := EventDown{}
		tablename = obj.TableName()
		reflectVal = reflect.Indirect(reflect.ValueOf(obj))
		database = DB().Connection
	case "eventdownhfd":
		obj := AlarmHFD{}
		tablename = obj.TableName()
		reflectVal = reflect.Indirect(reflect.ValueOf(obj))
		database = DBRealtime()
	}

	query := database.NewQuery().From(tablename).Skip(p.Skip).Take(p.Take)
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
	defer csr.Close()
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	tmpResult := make([]tk.M, 0)
	e = csr.Fetch(&tmpResult, 0, false)

	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	queryC := database.NewQuery().From(tablename).Where(dbox.And(filter...))
	ccount, e := queryC.Cursor(nil)
	defer ccount.Close()
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	totalPower := 0.0
	totalPowerLost := 0.0
	totalActivePower := 0.0
	avgWindSpeed := 0.0
	totalTurbine := 0
	totalEnergy := 0.0
	AvgWS := 0.0
	totalDuration := 0.0
	if needTotalTurbine {
		aggrData := []tk.M{}
		queryAggr := database.NewQuery().From(tablename)

		switch tipe {
		case "scadaoem":
			queryAggr = queryAggr.Aggr(dbox.AggrSum, "$power", "TotalPower").
				Aggr(dbox.AggrSum, "$powerlost", "TotalPowerLost").
				Aggr(dbox.AggrSum, "$ai_intern_activpower", "TotalActivePower").
				Aggr(dbox.AggrSum, "$ai_intern_windspeed", "AvgWindSpeed").
				Aggr(dbox.AggrSum, "$energy", "TotalEnergy").
				Group("turbine").Where(dbox.And(filter...))
		case "eventraw":
			queryAggr = queryAggr.Aggr(dbox.AggrSum, 1, "countData").
				Group("turbine").Where(dbox.And(filter...))
		case "eventdown", "eventdownhfd":
			queryAggr = queryAggr.Aggr(dbox.AggrSum, "$duration", "duration").
				Group("turbine").Where(dbox.And(filter...))
		case "scadahfd":
			queryAggr = queryAggr.Group("turbine").Where(dbox.And(filter...))
		}

		caggr, e := queryAggr.Cursor(nil)
		defer caggr.Close()
		if e != nil {
			return helper.CreateResult(false, nil, e.Error())
		}
		e = caggr.Fetch(&aggrData, 0, false)
		if e != nil {
			return helper.CreateResult(false, nil, e.Error())
		}
		totalTurbine = tk.SliceLen(aggrData)
		switch tipe {
		case "scadaoem":
			for _, val := range aggrData {
				totalPower += val.GetFloat64("TotalPower")
				totalPowerLost += val.GetFloat64("TotalPowerLost")
				totalActivePower += val.GetFloat64("TotalActivePower")
				avgWindSpeed += val.GetFloat64("AvgWindSpeed")
				totalEnergy += val.GetFloat64("TotalEnergy")
			}
			if ccount.Count() > 0.0 {
				AvgWS = avgWindSpeed / float64(ccount.Count())
			}
		case "eventdown", "eventdownhfd":
			for _, val := range aggrData {
				totalDuration += val.GetFloat64("duration")
			}
		case "scadahfd":
			/*totalActivePower = m.getSummaryColumn(filter, "fast_activepower_kw", "sum", tablename)
			AvgWS = m.getSummaryColumn(filter, "fast_windspeed_ms", "avg", tablename)*/
			totalActivePower = m.getSummaryColumn(filter, "activepower_kw", "sum", tablename)
			AvgWS = m.getSummaryColumn(filter, "windspeed_ms", "avg", tablename)
		}
	}

	header := map[string]string{}
	for i := 0; i < reflectVal.Type().NumField(); i++ {
		header[reflectVal.Type().Field(i).Name] = reflectVal.Field(i).Type().Name()
	}
	result, e := CheckData(tmpResult, filter, header, tipe)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
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
		TotalDuration    float64
		LastFilter       *helper.FilterJS
		LastSort         []helper.Sorting
	}{
		Data:             result,
		Total:            ccount.Count(),
		TotalPower:       totalPower,
		TotalPowerLost:   totalPowerLost,
		TotalActivePower: totalActivePower,
		AvgWindSpeed:     AvgWS, //avgWindSpeed / float64(ccount.Count()),
		TotalTurbine:     totalTurbine,
		TotalEnergy:      totalActivePower / 6,
		TotalDuration:    totalDuration,
		LastFilter:       p.Filter,
		LastSort:         p.Sort,
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
	p.Misc.Set("knot_data", k)
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
	defer csr.Close()
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

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
		Data       []JMR
		Total      int
		LastFilter *helper.FilterJS
		LastSort   []helper.Sorting
	}{
		Data:       tmpResult,
		Total:      ccount.Count(),
		LastFilter: p.Filter,
		LastSort:   p.Sort,
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
	p.Misc.Set("knot_data", k)
	filter, _ := p.ParseFilter()

	query := DB().Connection.NewQuery().From(new(JMR).TableName())
	query.Where(dbox.And(filter...))
	csr, e := query.Cursor(nil)
	defer csr.Close()
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	jmrResult := make([]JMR, 0)
	e = csr.Fetch(&jmrResult, 0, false)

	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

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
func (m *DataBrowserController) GetCustomList_DRAFT(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	p := new(helper.Payloads)
	e := k.GetPayload(&p)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}
	p.Misc.Set("knot_data", k)
	filter, _ := p.ParseFilter()
	tipe := p.Misc.GetString("tipe")

	tablename := "Scada10MinHFD"
	scadaFieldList := []string{"_id"}
	aggrFieldList := []string{"turbine", "power", "powerlost", "energy", "ai_intern_activpower", "ai_intern_windspeed"}
	source := "ScadaDataHFD"
	var val1 reflect.Value
	switch tipe {
	case "ScadaOEM":
		tablename = new(ScadaDataOEM).TableName()
		scadaFieldList = append(scadaFieldList, "timestamputc")
		source = "ScadaDataOEM"
		obj1 := ScadaDataOEM{}
		val1 = reflect.Indirect(reflect.ValueOf(obj1))
	case "ScadaHFD":
		filter = append(filter, dbox.Eq("isnull", false))
		aggrFieldList = []string{"turbine", "activepower_kw", "windspeed_ms"}
	}

	ids := ""
	projection := map[string]int{}
	if p.Custom.Has("ColumnList") {
		for _, _val := range p.Custom["ColumnList"].([]interface{}) {
			_tkm, _ := tk.ToM(_val)
			ids = strings.ToLower(_tkm.GetString("_id"))
			if _tkm.GetString("source") == source {
				projection[ids] = 1
				scadaFieldList = append(scadaFieldList, ids)
			}
		}
	}
	for _, val := range aggrFieldList {
		projection[val] = 1
	}
	sort.Float64s([]float64{})
	matches := []tk.M{}
	tStart := time.Time{}
	tEnd := time.Time{}
	for _, val := range filter {
		matches = append(matches, tk.M{
			val.Field: tk.M{val.Op: val.Value},
		})
		if val.Field == "timestamp" {
			if val.Op == "$gte" {
				tStart = val.Value.(time.Time).UTC()
			} else {
				tEnd = val.Value.(time.Time).UTC()
			}
		}
	}
	pipes := []tk.M{}
	pipes = append(pipes, tk.M{"$match": tk.M{"$and": matches}})
	pipes = append(pipes, tk.M{"$project": projection})

	csr, e := DB().Connection.NewQuery().
		From(tablename).Command("pipe", pipes).Cursor(nil)
	defer csr.Close()
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	countAll := 0.0
	totalPower := 0.0
	totalPowerLost := 0.0
	totalActivePower := 0.0
	avgWindSpeed := 0.0
	totalTurbine := 0
	totalEnergy := 0.0
	AvgWS := 0.0
	startIdx := tk.ToFloat64(p.Skip-1, 2, tk.RoundingAuto)
	endIdx := startIdx + tk.ToFloat64(p.Take, 2, tk.RoundingAuto)

	turbineList := map[string]int{}
	results := []tk.M{}
	_result := tk.M{}
	timestampCount := tStart
	timestampSorted := []time.Time{}
	dataPerTimestamp := map[time.Time][]tk.M{}
	for {
		timestampCount = timestampCount.Add(time.Duration(time.Minute * 10))
		if timestampCount.After(tEnd.UTC()) {
			break
		}
		timestampSorted = append(timestampSorted, timestampCount)
	}

	t0 := time.Now()
	for {
		_result = tk.M{}
		e = csr.Fetch(&_result, 1, false)
		if e != nil {
			break
		}
		switch tipe {
		case "ScadaOEM":
			totalPower += _result.GetFloat64("power")
			totalPowerLost += _result.GetFloat64("powerlost")
			totalActivePower += _result.GetFloat64("ai_intern_activpower")
			avgWindSpeed += _result.GetFloat64("ai_intern_windspeed")
			totalEnergy += _result.GetFloat64("energy")
		case "ScadaHFD":
			totalActivePower += _result.GetFloat64("activepower_kw")
			avgWindSpeed += _result.GetFloat64("windspeed_ms")
		}
		turbineList[_result.GetString("turbine")] = 1
		timestamp := _result.Get("timestamp", time.Time{}).(time.Time)
		dataPerTimestamp[timestamp] = append(dataPerTimestamp[timestamp], _result)

		countAll++
	}
	tk.Println("durasi", time.Since(t0).Seconds())
	counterIdx := 0.0
	counterTake := 0
	for _, _timestamp := range timestampSorted {
		_data, hasData := dataPerTimestamp[_timestamp]
		if hasData {
			for _, val := range _data {
				if counterIdx > startIdx && counterIdx <= endIdx {
					results = append(results, val)
					counterTake++
				}
				if counterTake == p.Take {
					break
				}
				counterIdx++
			}
		}
	}

	if countAll > 0.0 {
		AvgWS = avgWindSpeed / countAll
	}
	totalTurbine = len(turbineList)

	allFieldRequested := scadaFieldList
	allHeader := map[string]string{}
	header := map[string]string{}
	fieldName := ""
	switch tipe {
	case "ScadaOEM":
		for i := 0; i < val1.Type().NumField(); i++ {
			fieldName = strings.ToLower(val1.Type().Field(i).Name)
			allHeader[fieldName] = val1.Field(i).Type().Name()
		}
		for _, val := range allFieldRequested {
			header[val] = allHeader[val]
		}
	case "ScadaHFD":
		hfdexlist := []string{"timestamp", "projectname", "turbine", "turbinestate", "statedescription"}
		for _, val := range allFieldRequested {
			if tk.HasMember(hfdexlist, val) {
				if val == "timestamp" {
					header[val] = "time.Time"
				} else if val == "turbinestate" {
					header[val] = "int"
				} else {
					header[val] = "string"
				}
			} else {
				header[val] = "float64"
			}
		}
	}

	result, e := CheckData(results, filter, header, "custom")
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
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
		LastFilter       *helper.FilterJS
		LastSort         []helper.Sorting
	}{
		Data:             result,
		Total:            tk.ToInt(countAll, tk.RoundingAuto),
		TotalPower:       totalPower,
		TotalPowerLost:   totalPowerLost,
		TotalActivePower: totalActivePower,
		AvgWindSpeed:     AvgWS, //avgWindSpeed / float64(ccount.Count()),
		TotalTurbine:     totalTurbine,
		TotalEnergy:      totalActivePower / 6,
		LastFilter:       p.Filter,
		LastSort:         p.Sort,
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
	p.Misc.Set("knot_data", k)
	filter, _ := p.ParseFilter()
	tipe := p.Misc.GetString("tipe")

	// tablename := new(ScadaDataHFD).TableName()
	tablename := "Scada10MinHFD"
	arrscadaoem := []string{"_id"}
	source := "ScadaDataHFD"
	// timestamp := "timestamp"
	var val1 reflect.Value
	switch tipe {
	case "ScadaOEM":
		tablename = new(ScadaDataOEM).TableName()
		arrscadaoem = append(arrscadaoem, "timestamputc")
		source = "ScadaDataOEM"
		// timestamp = "timestamputc"
		obj1 := ScadaDataOEM{}
		val1 = reflect.Indirect(reflect.ValueOf(obj1))
	case "ScadaHFD":
		filter = append(filter, dbox.Eq("isnull", false))
	}

	// istimestamp := false
	// arrmettower := []string{}
	ids := ""
	projection := map[string]int{}
	if p.Custom.Has("ColumnList") {
		for _, _val := range p.Custom["ColumnList"].([]interface{}) {
			_tkm, _ := tk.ToM(_val)
			ids = strings.ToLower(_tkm.GetString("_id"))
			if _tkm.GetString("source") == source {
				projection[ids] = 1
				arrscadaoem = append(arrscadaoem, ids)
				/*if ids == "timestamp" {
					istimestamp = true
				}*/
			}
			/*else if _tkm.GetString("source") == "MetTower" {
				arrmettower = append(arrmettower, ids)
			}*/
		}
	}
	matches := []tk.M{}
	_tstart, _tend, _usethis := time.Time{}, time.Time{}, false
	for _, val := range p.Filter.Filters {
		for _, xval := range val.Filters {
			if xval.Field == "timestamp" {
				_xtime, _e := time.Parse("2006-01-02T15:04:05.000Z", tk.ToString(xval.Value))
				if _e != nil {
					_xtime, _ = time.Parse("2006-01-02 15:04:05", tk.ToString(xval.Value))
				}
				if xval.Op == "lte" {
					_tend = _xtime
				} else {
					_tstart = _xtime
				}
			}
		}
	}

	if _sub := _tend.UTC().Sub(_tstart.UTC()).Hours(); !_tstart.IsZero() && !_tend.IsZero() && _sub >= 0 && _sub < 24 {
		_usethis = true
	}

	for _, val := range filter {
		ttkm := tk.M{
			val.Field: tk.M{val.Op: val.Value},
		}

		if val.Field == "timestamp" {
			ttime := val.Value.(time.Time).UTC()
			if _usethis && val.Op == "$lte" {
				ttime = _tend
			} else if _usethis {
				ttime = _tstart
			}
			ttkm = tk.M{
				val.Field: tk.M{val.Op: ttime},
			}
			// tk.Println(val, ttkm)
		}

		matches = append(matches, ttkm)
	}
	if tipe == "ScadaHFD" {
		matches = append(matches, tk.M{
			"windspeed_ms": tk.M{"$gt": -200},
		})
		matches = append(matches, tk.M{
			"windspeed_ms": tk.M{"$exists": true},
		})
		matches = append(matches, tk.M{
			"activepower_kw": tk.M{"$exists": true},
		})
	}

	pipes := []tk.M{}
	// 	tk.M{"$match": tk.M{"$and": matches}},
	// 	tk.M{"$project": projection},
	// }
	sortList := map[string]int{}
	if len(p.Sort) > 0 {
		for _, val := range p.Sort {
			if val.Dir == "desc" {
				sortList[strings.ToLower(val.Field)] = -1
			} else {
				sortList[strings.ToLower(val.Field)] = 1
			}
		}
		pipes = append(pipes, tk.M{"$sort": sortList})
	}
	pipes = append(pipes, tk.M{"$match": tk.M{"$and": matches}})
	pipes = append(pipes, tk.M{"$project": projection})

	pipes = append(pipes, []tk.M{
		tk.M{"$skip": p.Skip},
		tk.M{"$limit": p.Take},
	}...)
	//tk.Printf("%#v\n", pipes)

	//timenow := time.Now()
	csr, e := DB().Connection.NewQuery().
		From(tablename).Command("pipe", pipes).Cursor(nil)
	defer csr.Close()
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}
	// defer csr.Close()
	// duration := time.Now().Sub(timenow).Seconds()
	// tk.Printf("Kondisi 1 = %v\n", duration)
	//tk.Printf("Total = %v\n", csr.Count())

	// timenow = time.Now()
	results := make([]tk.M, 0)
	// e = csr.Fetch(&results, 0, false)
	// item := tk.M{}
	for {
		item := tk.M{}
		e = csr.Fetch(&item, 1, false)
		if e != nil {
			e = nil
			break
		}
		results = append(results, item)
	}
	// tk.Printf("Total = %v\n", len(results))

	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}
	// csr.Close()
	// duration = time.Now().Sub(timenow).Seconds()
	// tk.Printf("Kondisi 2 = %v\n", duration)

	// arrmettowercond := []interface{}{}

	/*config := lh.ReadConfig()

	loc, err := time.LoadLocation(config["ReadTimeLoc"])
	if err != nil {
		tk.Printfn("Get time in %s found %s", config["ReadTimeLoc"], err.Error())
	}*/

	// timenow = time.Now()
	// for i, val := range results {
	// 	if val.Has("timestamputc") {
	// 		strangeTime := val.Get("timestamputc", time.Time{}).(time.Time).UTC().In(loc)
	// 		itime := time.Date(strangeTime.Year(), strangeTime.Month(), strangeTime.Day(),
	// 			strangeTime.Hour(), strangeTime.Minute(), strangeTime.Second(), strangeTime.Nanosecond(), time.UTC)
	// 		// arrmettowercond = append(arrmettowercond, itime)
	// 		val.Set("timestamputc", itime)
	// 		results[i] = val
	// 	}
	// 	if istimestamp {
	// 		itime := val.Get("timestamp", time.Time{}).(time.Time).UTC()
	// 		/*if tipe == "ScadaHFD" {
	// 			arrmettowercond = append(arrmettowercond, itime)
	// 		}*/
	// 		val.Set("timestamp", itime)
	// 		results[i] = val
	// 	}
	// }
	// duration = time.Now().Sub(timenow).Seconds()
	// tk.Printf("Kondisi 3 = %v\n", duration)

	/*tkmmet := tk.M{}
	if len(arrmettower) > 0 && len(arrmettowercond) > 0 {
		arrmettower = append(arrmettower, "timestamp")
		queryMet := DB().Connection.NewQuery().
			Select(arrmettower...).
			From("MetTower")

		filterMet := dbox.In("timestamp", arrmettowercond...)
		for _, val := range filter {
			if val.Field == "projectname" {
				filterMet = dbox.And(dbox.Eq(val.Field, val.Value), filterMet)
			}
		}
		_csr, _e := queryMet.Where(filterMet).Cursor(nil)

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
			itime := val.Get(timestamp, time.Time{}).(time.Time).UTC().String()
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
	}*/

	countAll := 0.0 //jangan dibaca dengan lantang, nanti pada gempar
	totalPower := 0.0
	totalPowerLost := 0.0
	totalActivePower := 0.0
	avgWindSpeed := 0.0
	totalTurbine := 0
	totalEnergy := 0.0
	AvgWS := 0.0

	aggrData := []tk.M{}
	groups, gprojects := tk.M{}, tk.M{}

	switch tipe {
	case "ScadaOEM":
		groups = tk.M{
			"_id":              "$turbine",
			"TotalPower":       tk.M{"$sum": "$power"},
			"TotalPowerLost":   tk.M{"$sum": "$powerlost"},
			"TotalActivePower": tk.M{"$sum": "$ai_intern_activpower"},
			"AvgWindSpeed":     tk.M{"$sum": "$ai_intern_windspeed"},
			"TotalEnergy":      tk.M{"$sum": "$energy"},
			"DataCount":        tk.M{"$sum": 1},
		}
		gprojects.Set("turbine", 1).Set("power", 1).Set("powerlost", 1).Set("ai_intern_activpower", 1).Set("ai_intern_windspeed", 1).Set("energy", 1)
	case "ScadaHFD":
		groups = tk.M{
			"_id":              "$turbine",
			"TotalActivePower": tk.M{"$sum": "$power"},
			"AvgWindSpeed":     tk.M{"$sum": "$avgwindspeed"},
			"DataCount":        tk.M{"$sum": 1},
		}
		gprojects.Set("turbine", 1).Set("power", 1).Set("avgwindspeed", 1)
		tablename = "ScadaData"
		for i, _match := range matches {
			if _match.Has("windspeed_ms") {
				_match.Set("avgwindspeed", _match.Get("windspeed_ms"))
				_match.Unset("windspeed_ms")
			}
			if _match.Has("activepower_kw") {
				_match.Set("power", _match.Get("activepower_kw"))
				_match.Unset("activepower_kw")
			}
			if _match.Has("isnull") {
				_match.Set("available", 1)
				_match.Unset("isnull")
			}
			matches[i] = _match
		}
	}

	pipes = []tk.M{
		tk.M{"$match": tk.M{"$and": matches}},
		tk.M{"$project": gprojects},
		tk.M{"$group": groups},
	}
	// tk.Printf("%#v\n", matches)
	// tk.Printf("%#v\n", groups)

	// tk.Printf("%#v\n", pipes)

	// timenow = time.Now()
	caggr, e := DB().Connection.NewQuery().
		From(tablename).Command("pipe", pipes).Cursor(nil)
	defer caggr.Close()

	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	// duration = time.Now().Sub(timenow).Seconds()
	// tk.Printf("Kondisi 4 = %v\n", duration)

	// timenow = time.Now()
	// e = caggr.Fetch(&aggrData, 0, false)

	for {
		item := tk.M{}
		e = caggr.Fetch(&item, 1, false)
		if e != nil {
			e = nil
			break
		}
		aggrData = append(aggrData, item)
	}

	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	// duration = time.Now().Sub(timenow).Seconds()
	// tk.Printf("Kondisi 5 = %v\n", duration)

	totalTurbine = tk.SliceLen(aggrData)
	switch tipe {
	case "ScadaOEM":
		for _, val := range aggrData {
			totalPower += val.GetFloat64("TotalPower")
			totalPowerLost += val.GetFloat64("TotalPowerLost")
			totalActivePower += val.GetFloat64("TotalActivePower")
			avgWindSpeed += val.GetFloat64("AvgWindSpeed")
			totalEnergy += val.GetFloat64("TotalEnergy")
			countAll += val.GetFloat64("DataCount")
		}
		if countAll > 0.0 {
			AvgWS = avgWindSpeed / countAll
		}
	case "ScadaHFD":
		/*totalActivePower = m.getSummaryColumn(filter, "activepower_kw", "sum", tablename)
		AvgWS = m.getSummaryColumn(filter, "windspeed_ms", "avg", tablename)*/
		for _, val := range aggrData {
			totalActivePower += val.GetFloat64("TotalActivePower")
			avgWindSpeed += val.GetFloat64("AvgWindSpeed")
			countAll += val.GetFloat64("DataCount")
		}
		if countAll > 0.0 {
			AvgWS = avgWindSpeed / countAll
		}
	}

	allFieldRequested := arrscadaoem
	// allFieldRequested = append(allFieldRequested, arrmettower...)
	allHeader := map[string]string{}
	header := map[string]string{}
	/*obj2 := MetTower{}
	val2 := reflect.Indirect(reflect.ValueOf(obj2))*/
	fieldName := ""
	switch tipe {
	case "ScadaOEM":
		for i := 0; i < val1.Type().NumField(); i++ {
			fieldName = strings.ToLower(val1.Type().Field(i).Name)
			allHeader[fieldName] = val1.Field(i).Type().Name()
		}
		/*for i := 0; i < val2.Type().NumField(); i++ {
			fieldName = strings.ToLower(val2.Type().Field(i).Name)
			allHeader[fieldName] = val2.Field(i).Type().Name()
		}*/
		for _, val := range allFieldRequested {
			header[val] = allHeader[val]
		}
	case "ScadaHFD":
		hfdexlist := []string{"timestamp", "projectname", "turbine", "turbinestate", "statedescription"}
		for _, val := range allFieldRequested {
			if tk.HasMember(hfdexlist, val) {
				if val == "timestamp" {
					header[val] = "time.Time"
				} else if val == "turbinestate" {
					header[val] = "int"
				} else {
					header[val] = "string"
				}
			} else {
				header[val] = "float64"
			}
		}
	}

	// timenow = time.Now()
	result, e := CheckData(results, filter, header, "custom")
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}
	// duration = time.Now().Sub(timenow).Seconds()
	// tk.Printf("Kondisi 6 = %v\n", duration)

	data := struct {
		Data             []tk.M
		Total            int
		TotalPower       float64
		TotalPowerLost   float64
		TotalActivePower float64
		AvgWindSpeed     float64
		TotalTurbine     int
		TotalEnergy      float64
		LastFilter       *helper.FilterJS
		LastSort         []helper.Sorting
	}{
		Data:             result,
		Total:            tk.ToInt(countAll, tk.RoundingAuto),
		TotalPower:       totalPower,
		TotalPowerLost:   totalPowerLost,
		TotalActivePower: totalActivePower,
		AvgWindSpeed:     AvgWS, //avgWindSpeed / float64(ccount.Count()),
		TotalTurbine:     totalTurbine,
		TotalEnergy:      totalActivePower / 6,
		LastFilter:       p.Filter,
		LastSort:         p.Sort,
	}

	return helper.CreateResult(true, data, "success")
}

func (m *DataBrowserController) GetCustomFarmWise(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	p := new(helper.Payloads)
	e := k.GetPayload(&p)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}
	p.Misc.Set("knot_data", k)
	filter, _ := p.ParseFilter()
	// tipe := p.Misc.GetString("tipe")

	tablename := "Scada10MinHFD"
	arrscadaoem := []string{"_id"}
	source := "ScadaDataHFD"

	filter = append(filter, dbox.Eq("isnull", false))

	ids := ""
	projection := map[string]int{}
	if p.Custom.Has("ColumnList") {
		for _, _val := range p.Custom["ColumnList"].([]interface{}) {
			_tkm, _ := tk.ToM(_val)
			ids = strings.ToLower(_tkm.GetString("_id"))
			if _tkm.GetString("source") == source {
				projection[ids] = 1
				arrscadaoem = append(arrscadaoem, ids)
			}
		}
	}

	matches := []tk.M{}
	_tstart, _tend, _usethis := time.Time{}, time.Time{}, false
	for _, val := range p.Filter.Filters {
		for _, xval := range val.Filters {
			if xval.Field == "timestamp" {
				_xtime, _e := time.Parse("2006-01-02T15:04:05.000Z", tk.ToString(xval.Value))
				if _e != nil {
					_xtime, _ = time.Parse("2006-01-02 15:04:05", tk.ToString(xval.Value))
				}
				if xval.Op == "lte" {
					_tend = _xtime
				} else {
					_tstart = _xtime
				}
			}
		}
	}

	if _sub := _tend.UTC().Sub(_tstart.UTC()).Hours(); !_tstart.IsZero() && !_tend.IsZero() && _sub >= 0 && _sub < 24 {
		_usethis = true
	}

	for _, val := range filter {

		if val.Field == "turbine" {
			continue
		}

		ttkm := tk.M{
			val.Field: tk.M{val.Op: val.Value},
		}

		if val.Field == "timestamp" {
			ttime := val.Value.(time.Time).UTC()
			if _usethis && val.Op == "$lte" {
				ttime = _tend
			} else if _usethis {
				ttime = _tstart
			}
			ttkm = tk.M{
				val.Field: tk.M{val.Op: ttime},
			}
		}

		matches = append(matches, ttkm)
	}

	matches = append(matches, tk.M{
		"windspeed_ms": tk.M{"$gt": -200},
	})
	matches = append(matches, tk.M{
		"windspeed_ms": tk.M{"$exists": true},
	})
	matches = append(matches, tk.M{
		"activepower_kw": tk.M{"$exists": true},
	})

	pipes := []tk.M{}
	agroups := tk.M{
		"_id": tk.M{"projectname": "$projectname", "timestamp": "$timestamp"},
	}

	fproject := map[string]int{"projectname": 1, "timestamp": 1}
	for field, _ := range projection {
		if field == "turbine" || field == "projectname" || field == "timestamp" {
			continue
		}

		// fproject[tk.Sprintf("%s_sum", field)] = 1
		// fproject[tk.Sprintf("%s_count", field)] = 1

		// agroups.Set(field+"_sum", tk.M{"$sum": tk.Sprintf("$%s_sum", field)})
		// agroups.Set(field+"_count", tk.M{"$sum": tk.Sprintf("$%s_count", field)})

		fproject[field] = 1
		agroups.Set(field, tk.M{"$avg": tk.Sprintf("$%s", field)})
	}

	sortList := map[string]int{}
	if len(p.Sort) > 0 {
		for _, val := range p.Sort {
			field := strings.ToLower(val.Field)
			if field == "turbine" || field == "projectname" || field == "timestamp" {
				continue
			}
			if val.Dir == "desc" {
				sortList[field] = -1
			} else {
				sortList[field] = 1
			}
		}
	}

	if len(sortList) == 0 {
		sortList["_id"] = 1
	}

	pipes = append(pipes, tk.M{"$match": tk.M{"$and": matches}})
	pipes = append(pipes, tk.M{"$project": fproject})
	pipes = append(pipes, tk.M{"$group": agroups})
	pipes = append(pipes, tk.M{"$sort": sortList})

	pipes = append(pipes, []tk.M{
		tk.M{"$skip": p.Skip},
		tk.M{"$limit": p.Take},
	}...)

	// tk.Println(" === ", pipes)
	// tk.Println(" === ", projection)

	// timenow := time.Now()
	csr, e := DB().Connection.NewQuery().
		From(tablename).Command("pipe", pipes).Cursor(nil)
	defer csr.Close()
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	// tk.Println("CSR : ", time.Since(timenow).String())
	// timenow = time.Now()

	results := make([]tk.M, 0)
	items := []tk.M{}
	e = csr.Fetch(&items, 0, false)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	// tk.Println("FETCH : ", time.Since(timenow).String())
	// timenow = time.Now()

	for _, item := range items {

		idtk, _ := tk.ToM(item.Get("_id"))

		// tk.Println(" -- ", idtk["timestamp"])

		item.Set("projectname", idtk["projectname"])
		item.Set("timestamp", idtk["timestamp"])
		// for field, _ := range projection {
		// 	if field == "turbine" || field == "projectname" || field == "timestamp" {
		// 		continue
		// 	}
		// 	resitem.Set(field, tk.Div(item.GetFloat64(field+"_sum"), item.GetFloat64(field+"_count")))
		// }

		results = append(results, item)
	}
	// tk.Println("PRE-PRO 01 : ", time.Since(timenow).String())
	// timenow = time.Now()

	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	countAll := 0.0 //jangan dibaca dengan lantang, nanti pada gempar
	totalPower := 0.0
	totalPowerLost := 0.0
	totalActivePower := 0.0
	avgWindSpeed := 0.0
	totalTurbine := 0
	AvgWS := 0.0

	groups, gprojects := tk.M{}, tk.M{}
	groups = tk.M{
		"_id":              tk.M{"projectname": "$projectname", "timestamp": "$timestamp"},
		"TotalActivePower": tk.M{"$sum": "$power"},
		"AvgWindSpeed":     tk.M{"$sum": "$avgwindspeed"},
	}
	gprojects.Set("projectname", 1).Set("power", 1).Set("avgwindspeed", 1).Set("timestamp", 1)
	tablename = "ScadaData"
	for i, _match := range matches {
		if _match.Has("windspeed_ms") {
			_match.Set("avgwindspeed", _match.Get("windspeed_ms"))
			_match.Unset("windspeed_ms")
		}
		if _match.Has("activepower_kw") {
			_match.Set("power", _match.Get("activepower_kw"))
			_match.Unset("activepower_kw")
		}
		if _match.Has("isnull") {
			_match.Set("available", 1)
			_match.Unset("isnull")
		}
		matches[i] = _match
	}

	pipes = []tk.M{
		tk.M{"$match": tk.M{"$and": matches}},
		tk.M{"$project": gprojects},
		tk.M{"$group": groups},
	}

	caggr, e := DB().Connection.NewQuery().
		From(tablename).Command("pipe", pipes).Cursor(nil)
	defer caggr.Close()

	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	tproject := tk.M{}
	for {
		item := tk.M{}
		e = caggr.Fetch(&item, 1, false)
		if e != nil {
			break
		}

		idtk := item.Get("_id", tk.M{}).(tk.M)
		tproject.Set(idtk.GetString("projectname"), 1)

		totalActivePower += item.GetFloat64("TotalActivePower")
		avgWindSpeed += item.GetFloat64("AvgWindSpeed")
		countAll += 1
	}

	if countAll > 0.0 {
		AvgWS = avgWindSpeed / countAll
	}

	allFieldRequested := arrscadaoem
	header := map[string]string{}

	hfdexlist := []string{"timestamp", "projectname", "turbine", "turbinestate", "statedescription"}
	for _, val := range allFieldRequested {
		if tk.HasMember(hfdexlist, val) {
			if val == "timestamp" {
				header[val] = "time.Time"
			} else if val == "turbinestate" {
				header[val] = "int"
			} else {
				header[val] = "string"
			}
		} else {
			header[val] = "float64"
		}
	}

	// tk.Println("AGGR 01 : ", time.Since(timenow).String())
	// timenow = time.Now()

	result, e := CheckData(results, filter, header, "custom")
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	// tk.Println("PRE-PRO 02 : ", time.Since(timenow).String())
	// timenow = time.Now()

	data := struct {
		Data             []tk.M
		Total            int
		TotalPower       float64
		TotalPowerLost   float64
		TotalActivePower float64
		AvgWindSpeed     float64
		TotalTurbine     int
		TotalProject     int
		TotalEnergy      float64
		LastFilter       *helper.FilterJS
		LastSort         []helper.Sorting
	}{
		Data:             result,
		Total:            tk.ToInt(countAll, tk.RoundingAuto),
		TotalPower:       totalPower,
		TotalPowerLost:   totalPowerLost,
		TotalActivePower: totalActivePower,
		AvgWindSpeed:     AvgWS, //avgWindSpeed / float64(ccount.Count()),
		TotalTurbine:     totalTurbine,
		TotalProject:     len(tproject),
		TotalEnergy:      totalActivePower / 6,
		LastFilter:       p.Filter,
		LastSort:         p.Sort,
	}

	return helper.CreateResult(true, data, "success")
}

func (m *DataBrowserController) getSummaryColumn(filter []*dbox.Filter, column, aggr, tablename string) float64 {
	xFilter := []*dbox.Filter{}
	queryAggr := DB().Connection.NewQuery().From(tablename)
	tkm := []tk.M{}

	switch column {
	case "windspeed_ms":
		xFilter = append(filter, dbox.Gte(column, 0))
		xFilter = append(xFilter, dbox.Lte(column, 25))
	case "activepower_kw":
		xFilter = append(filter, dbox.Gte(column, -200))
		xFilter = append(xFilter, dbox.Lte(column, 3000))
		xFilter = append(xFilter, dbox.Or(dbox.Eq("isnull", false), dbox.Eq("isnull", nil)))
	default:
		return 0
	}

	switch aggr {
	case "sum":
		queryAggr.Aggr(dbox.AggrSum, "$"+column, "xValue")
	case "avg":
		queryAggr.Aggr(dbox.AggrAvr, "$"+column, "xValue")
	default:
		return 0
	}

	caggr, e := queryAggr.
		Group("projectname").Where(dbox.And(xFilter...)).
		Cursor(nil)
	defer caggr.Close()
	if e != nil {
		return 0
	}
	e = caggr.Fetch(&tkm, 0, false)
	if e != nil {
		return 0
	}

	xVal := 0.0
	for _, val := range tkm {
		xVal = val.GetFloat64("xValue")
	}

	return xVal
}

func (m *DataBrowserController) GetLostEnergyDetail(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	p := new(helper.Payloads)
	e := k.GetPayload(&p)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}
	p.Misc.Set("knot_data", k)
	matches, pipes, project := []tk.M{}, []tk.M{}, ""

	tstart, tend := time.Time{}, time.Time{}
	for _, val := range p.Filter.Filters {
		if val.Field == "startdate" {
			if val.Op == "$gte" {
				tstart = tk.ToDate(tk.ToString(val.Value)[:10], "2006-01-02")
			} else {
				tend = tk.ToDate(tk.ToString(val.Value)[:10], "2006-01-02")
			}
			continue
		}

		if val.Field == "projectname" {
			project = tk.ToString(val.Value)
		}

		_match := tk.M{val.Field: tk.M{val.Op: val.Value}}
		if val.Op == "$regex" {
			_match = tk.M{val.Field: tk.M{val.Op: val.Value, "$options": "i"}}
		}

		matches = append(matches, _match)
	}

	tend = tend.AddDate(0, 0, 1)
	dates := []tk.M{
		tk.M{"$and": []tk.M{tk.M{"startdate": tk.M{"$gte": tstart}}, tk.M{"startdate": tk.M{"$lt": tend}}}},
		tk.M{"$and": []tk.M{tk.M{"enddate": tk.M{"$gte": tstart}}, tk.M{"enddate": tk.M{"$lt": tend}}}},
		tk.M{"$and": []tk.M{tk.M{"startdate": tk.M{"$gte": tstart}}, tk.M{"enddate": tk.M{"$lt": tend}}}},
	}

	matches = append(matches, tk.M{"$or": dates})

	matches = append(matches, tk.M{"reduceavailability": true})

	pipes = append(pipes, tk.M{"$match": tk.M{"$and": matches}})
	pipes = append(pipes, tk.M{"$project": tk.M{"turbine": 1, "startdate": 1, "enddate": 1, "reduceavailability": 1, "powerlost": 1, "alertdescription": 1, "detail": 1}})
	pipes = append(pipes, tk.M{"$unwind": tk.M{"path": "$detail"}})
	pipes = append(pipes, tk.M{"$project": tk.M{"turbine": 1, "startdate": 1, "enddate": 1, "reduceavailability": 1, "powerlost": 1, "alertdescription": 1,
		"detail.startdate": 1, "detail.enddate": 1, "detail.powerlost": 1, "detail.duration": 1, "detail.griddown": 1, "detail.machinedown": 1}})
	pipes = append(pipes, tk.M{"$match": tk.M{"detail.startdate": tk.M{"$gte": tstart, "$lt": tend}}})
	pipesAggr := []tk.M{}
	for _, val := range pipes {
		pipesAggr = append(pipesAggr, val)
	}
	pipesAggr = append(pipesAggr, tk.M{
		"$group": tk.M{
			"_id":            "$turbine",
			"totalpowerlost": tk.M{"$sum": "$detail.powerlost"},
			"totalduration":  tk.M{"$sum": "$detail.duration"},
			"totaldata":      tk.M{"$sum": 1},
		},
	})

	if len(p.Sort) > 0 {
		sortList := map[string]int{}
		for _, val := range p.Sort {
			if val.Dir == "desc" {
				sortList[strings.ToLower(val.Field)] = -1
			} else {
				sortList[strings.ToLower(val.Field)] = 1
			}
		}
		pipes = append(pipes, tk.M{"$sort": sortList})
	}

	pipes = append(pipes, []tk.M{
		tk.M{"$skip": p.Skip},
		tk.M{"$limit": p.Take},
	}...)

	csr, e := DB().Connection.NewQuery().
		From("Alarm").Command("pipe", pipes).Cursor(nil)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}
	defer csr.Close()

	Datas := []tk.M{}
	e = csr.Fetch(&Datas, 0, false)

	turbinename, e := helper.GetTurbineNameList(project)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	for i, data := range Datas {
		detail := data.Get("detail", tk.M{}).(tk.M)

		_tstart := detail.Get("startdate", time.Time{}).(time.Time)
		_tend := detail.Get("enddate", time.Time{}).(time.Time)

		if !_tstart.IsZero() && !_tend.IsZero() {
			detail.Set("duration", _tend.UTC().Sub(_tstart.UTC()).Seconds())
			data.Set("detail", detail)
		}

		idturbine := data.GetString("turbine")
		tname, cond := turbinename[idturbine]
		if !cond {
			tname = idturbine
		}
		data.Set("turbinename", tname)

		Datas[i] = data
	}

	csrAggr, e := DB().Connection.NewQuery().
		From("Alarm").Command("pipe", pipesAggr).Cursor(nil)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}
	defer csrAggr.Close()

	dataAggr := []tk.M{}
	e = csrAggr.Fetch(&dataAggr, 0, false)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	totalPower, totalDuration, totalData := 0.0, 0.0, 0
	for _, val := range dataAggr {
		totalPower += val.GetFloat64("totalpowerlost")
		totalDuration += val.GetFloat64("totalduration")
		totalData += val.GetInt("totaldata")
	}

	data := struct {
		Data           []tk.M
		Total          int
		TotalPowerLost float64
		TotalTurbine   int
		TotalDuration  float64
		LastFilter     *helper.FilterJS
		LastSort       []helper.Sorting
	}{
		Data:           Datas,
		Total:          totalData,
		TotalPowerLost: totalPower / 1000,
		TotalTurbine:   len(dataAggr),
		TotalDuration:  totalDuration,
		LastFilter:     p.Filter,
		LastSort:       p.Sort,
	}

	return helper.CreateResult(true, data, "success")
}

// Generate excel

func (m *DataBrowserController) GenExcelData(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	p := new(helper.Payloads)
	e := k.GetPayload(&p)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	p.Misc.Set("knot_data", k)
	filter, _ := p.ParseFilter()
	typeExcel := p.Misc.GetString("tipe")

	var pathDownload string
	header := []string{}
	tablename := ""
	fieldList := []string{}
	separator := "_"
	var database dbox.IConnection

	switch typeExcel {
	// case "ScadaOem":
	// 	header = []string{"TimeStamp", "Turbine", "AI intern R PidAngleOut", "AI intern ActivPower ", "AI intern I1 ", "AI intern I2", "AI intern I3", "AI intern NacelleDrill ", "AI intern NacellePos ", "AI intern PitchAkku V1 ", "AI intern PitchAkku V2 ", "AI intern PitchAkku V3 ", "AI intern PitchAngle1 ", "AI intern PitchAngle2 ", "AI intern PitchAngle3 ", "AI intern PitchConv Current1 ", "AI intern PitchConv Current2 ", "AI intern PitchConv Current3 ", "AI intern PitchAngleSP Diff1 ", "AI intern PitchAngleSP Diff2 ", "AI intern PitchAngleSP Diff3 ", "AI intern ReactivPower ", "AI intern RpmDiff ", "AI intern U1 ", "AI intern U2 ", "AI intern U3 ", "AI intern WindDirection ", "AI intern WindSpeed ", "AI Intern WindSpeedDif ", "AI speed RotFR ", "AI WindSpeed1 ", "AI WindSpeed2 ", "AI WindVane1 ", "AI WindVane2 ", "AI internCurrentAsym ", "Temp GearBox IMS NDE ", "AI intern WindVaneDiff ", "C intern SpeedGenerator ", "C intern SpeedRotor ", "AI intern Speed RPMDiff FR1 RotCNT ", "AI intern Frequency Grid ", "Temp GearBox HSS NDE ", "AI DrTrVibValue ", "AI intern InLastErrorConv1 ", "AI intern InLastErrorConv2 ", "AI intern InLastErrorConv3 ", "AI intern TempConv1 ", "AI intern TempConv2 ", "AI intern TempConv3 ", "AI intern PitchSpeed2", "Temp YawBrake 1 ", "Temp YawBrake 2 ", "Temp G1L1 ", "Temp G1L2 ", "Temp G1L3 ", "Temp YawBrake 4", "AI HydrSystemPressure ", "Temp BottomControlSection Low ", "Temp GearBox HSS DE ", "Temp GearOilSump ", "Temp GeneratorBearing DE ", "Temp GeneratorBearing NDE ", "Temp MainBearing ", "Temp GearBox IMS DE ", "Temp Nacelle ", "Temp Outdoor ", "AI TowerVibValueAxial ", "AI intern DiffGenSpeedSPToAct ", "Temp YawBrake 5", "AI intern SpeedGenerator Proximity ", "AI intern SpeedDiff Encoder Proximity ", "AI GearOilPressure ", "Temp CabinetTopBox Low ", "Temp CabinetTopBox ", "Temp BottomControlSection ", "Temp BottomPowerSection ", "Temp BottomPowerSection Low ", "AI intern Pitch1 Status High ", "AI intern Pitch2 Status High ", "AI intern Pitch3 Status High ", "AI intern InPosition1 ch3", "AI intern InPosition2 ch3", "AI intern InPosition3 ch3", "AI intern Temp Brake Blade1 ", "AI intern Temp Brake Blade2 ", "AI intern Temp Brake Blade3 ", "AI intern Temp PitchMotor Blade1 ", "AI intern Temp PitchMotor Blade2 ", "AI intern Temp PitchMotor Blade3 ", "AI intern Temp Hub Additional1 ", "AI intern Temp Hub Additional2 ", "AI intern Temp Hub Additional3 ", "AI intern Pitch1 Status Low ", "AI intern Pitch2 Status Low ", "AI intern Pitch3 Status Low ", "AI intern Battery VoltageBlade1 center ", "AI intern Battery VoltageBlade2 center ", "AI intern Battery VoltageBlade3 center ", "AI intern Battery ChargingCur Blade1 ", "AI intern Battery ChargingCur Blade2 ", "AI intern Battery ChargingCur Blade3 ", "AI intern Battery DischargingCur Blade1 ", "AI intern Battery DischargingCur Blade2 ", "AI intern Battery DischargingCur Blade3 ", "AI intern PitchMotor BrakeVoltage Blade1 ", "AI intern PitchMotor BrakeVoltage Blade2 ", "AI intern PitchMotor BrakeVoltage Blade3 ", "AI intern PitchMotor BrakeCurrent Blade1 ", "AI intern PitchMotor BrakeCurrent Blade2 ", "AI intern PitchMotor BrakeCurrent Blade3 ", "AI intern Temp HubBox Blade1 ", "AI intern Temp HubBox Blade2 ", "AI intern Temp HubBox Blade3 ", "AI intern Temp Pitch1 HeatSink ", "AI intern Temp Pitch2 HeatSink ", "AI intern Temp Pitch3 HeatSink ", "AI intern ErrorStackBlade1 ", "AI intern ErrorStackBlade2 ", "AI intern ErrorStackBlade3 ", "AI intern Temp BatteryBox Blade1 ", "AI intern Temp BatteryBox Blade2 ", "AI intern Temp BatteryBox Blade3 ", "AI intern DC LinkVoltage1 ", "AI intern DC LinkVoltage2 ", "AI intern DC LinkVoltage3 ", "Temp Yaw Motor1 ", "Temp Yaw Motor2 ", "Temp Yaw Motor3 ", "Temp Yaw Motor4 ", "AO DFIG Power Setpiont ", "AO DFIG Q Setpoint ", "AI DFIG Torque actual ", "AI DFIG SpeedGenerator Encoder ", "AI intern DFIG DC Link Voltage actual ", "AI intern DFIG MSC current ", "AI intern DFIG Main voltage ", "AI intern DFIG Main current ", "AI intern DFIG active power actual ", "AI intern DFIG reactive power actual ", "AI intern DFIG active power actual LSC ", "AI intern DFIG LSC current ", "AI intern DFIG Data log number ", "AI intern Damper OscMagnitude ", "AI intern Damper PassbandFullLoad ", "AI YawBrake TempRise1 ", "AI YawBrake TempRise2 ", "AI YawBrake TempRise3 ", "AI YawBrake TempRise4 ", "AI intern NacelleDrill at NorthPosSensor "}
	// 	tablename = new(ScadaDataOEM).TableName()
	// case "ScadaDataHFD":
	// 	header = []string{"TimeStamp", "ProjectName", "Turbine", "Fast_ActivePower_kW", "Fast_ActivePower_kW_StdDev", "Fast_ActivePower_kW_Min", "Fast_ActivePower_kW_Max", "Fast_ActivePower_kW_Count", "Fast_WindSpeed_ms", "Fast_WindSpeed_ms_StdDev", "Fast_WindSpeed_ms_Min", "Fast_WindSpeed_ms_Max", "Fast_WindSpeed_ms_Count", "Slow_NacellePos", "Slow_NacellePos_StdDev", "Slow_NacellePos_Min", "Slow_NacellePos_Max", "Slow_NacellePos_Count", "Slow_WindDirection", "Slow_WindDirection_StdDev", "Slow_WindDirection_Min", "Slow_WindDirection_Max", "Slow_WindDirection_Count", "Fast_CurrentL3", "Fast_CurrentL3_StdDev", "Fast_CurrentL3_Min", "Fast_CurrentL3_Max", "Fast_CurrentL3_Count", "Fast_CurrentL1", "Fast_CurrentL1_StdDev", "Fast_CurrentL1_Min", "Fast_CurrentL1_Max", "Fast_CurrentL1_Count", "Fast_ActivePowerSetpoint_kW", "Fast_ActivePowerSetpoint_kW_StdDev", "Fast_ActivePowerSetpoint_kW_Min", "Fast_ActivePowerSetpoint_kW_Max", "Fast_ActivePowerSetpoint_kW_Count", "Fast_CurrentL2", "Fast_CurrentL2_StdDev", "Fast_CurrentL2_Min", "Fast_CurrentL2_Max", "Fast_CurrentL2_Count", "Fast_DrTrVibValue", "Fast_DrTrVibValue_StdDev", "Fast_DrTrVibValue_Min", "Fast_DrTrVibValue_Max", "Fast_DrTrVibValue_Count", "Fast_GenSpeed_RPM", "Fast_GenSpeed_RPM_StdDev", "Fast_GenSpeed_RPM_Min", "Fast_GenSpeed_RPM_Max", "Fast_GenSpeed_RPM_Count", "Fast_PitchAccuV1", "Fast_PitchAccuV1_StdDev", "Fast_PitchAccuV1_Min", "Fast_PitchAccuV1_Max", "Fast_PitchAccuV1_Count", "Fast_PitchAngle", "Fast_PitchAngle_StdDev", "Fast_PitchAngle_Min", "Fast_PitchAngle_Max", "Fast_PitchAngle_Count", "Fast_PitchAngle3", "Fast_PitchAngle3_StdDev", "Fast_PitchAngle3_Min", "Fast_PitchAngle3_Max", "Fast_PitchAngle3_Count", "Fast_PitchAngle2", "Fast_PitchAngle2_StdDev", "Fast_PitchAngle2_Min", "Fast_PitchAngle2_Max", "Fast_PitchAngle2_Count", "Fast_PitchConvCurrent1", "Fast_PitchConvCurrent1_StdDev", "Fast_PitchConvCurrent1_Min", "Fast_PitchConvCurrent1_Max", "Fast_PitchConvCurrent1_Count", "Fast_PitchConvCurrent3", "Fast_PitchConvCurrent3_StdDev", "Fast_PitchConvCurrent3_Min", "Fast_PitchConvCurrent3_Max", "Fast_PitchConvCurrent3_Count", "Fast_PitchConvCurrent2", "Fast_PitchConvCurrent2_StdDev", "Fast_PitchConvCurrent2_Min", "Fast_PitchConvCurrent2_Max", "Fast_PitchConvCurrent2_Count", "Fast_PowerFactor", "Fast_PowerFactor_StdDev", "Fast_PowerFactor_Min", "Fast_PowerFactor_Max", "Fast_PowerFactor_Count", "Fast_ReactivePowerSetpointPPC_kVA", "Fast_ReactivePowerSetpointPPC_kVA_StdDev", "Fast_ReactivePowerSetpointPPC_kVA_Min", "Fast_ReactivePowerSetpointPPC_kVA_Max", "Fast_ReactivePowerSetpointPPC_kVA_Count", "Fast_ReactivePower_kVAr", "Fast_ReactivePower_kVAr_StdDev", "Fast_ReactivePower_kVAr_Min", "Fast_ReactivePower_kVAr_Max", "Fast_ReactivePower_kVAr_Count", "Fast_RotorSpeed_RPM", "Fast_RotorSpeed_RPM_StdDev", "Fast_RotorSpeed_RPM_Min", "Fast_RotorSpeed_RPM_Max", "Fast_RotorSpeed_RPM_Count", "Fast_VoltageL1", "Fast_VoltageL1_StdDev", "Fast_VoltageL1_Min", "Fast_VoltageL1_Max", "Fast_VoltageL1_Count", "Fast_VoltageL2", "Fast_VoltageL2_StdDev", "Fast_VoltageL2_Min", "Fast_VoltageL2_Max", "Fast_VoltageL2_Count", "Slow_CapableCapacitiveReactPwr_kVAr", "Slow_CapableCapacitiveReactPwr_kVAr_StdDev", "Slow_CapableCapacitiveReactPwr_kVAr_Min", "Slow_CapableCapacitiveReactPwr_kVAr_Max", "Slow_CapableCapacitiveReactPwr_kVAr_Count", "Slow_CapableInductiveReactPwr_kVAr", "Slow_CapableInductiveReactPwr_kVAr_StdDev", "Slow_CapableInductiveReactPwr_kVAr_Min", "Slow_CapableInductiveReactPwr_kVAr_Max", "Slow_CapableInductiveReactPwr_kVAr_Count", "Slow_DateTime_Sec", "Slow_DateTime_Sec_StdDev", "Slow_DateTime_Sec_Min", "Slow_DateTime_Sec_Max", "Slow_DateTime_Sec_Count", "Fast_PitchAngle1", "Fast_PitchAngle1_StdDev", "Fast_PitchAngle1_Min", "Fast_PitchAngle1_Max", "Fast_PitchAngle1_Count", "Fast_VoltageL3", "Fast_VoltageL3_StdDev", "Fast_VoltageL3_Min", "Fast_VoltageL3_Max", "Fast_VoltageL3_Count", "Slow_CapableCapacitivePwrFactor", "Slow_CapableCapacitivePwrFactor_StdDev", "Slow_CapableCapacitivePwrFactor_Min", "Slow_CapableCapacitivePwrFactor_Max", "Slow_CapableCapacitivePwrFactor_Count", "Fast_Total_Production_kWh", "Fast_Total_Production_kWh_StdDev", "Fast_Total_Production_kWh_Min", "Fast_Total_Production_kWh_Max", "Fast_Total_Production_kWh_Count", "Fast_Total_Prod_Day_kWh", "Fast_Total_Prod_Day_kWh_StdDev", "Fast_Total_Prod_Day_kWh_Min", "Fast_Total_Prod_Day_kWh_Max", "Fast_Total_Prod_Day_kWh_Count", "Fast_Total_Prod_Month_kWh", "Fast_Total_Prod_Month_kWh_StdDev", "Fast_Total_Prod_Month_kWh_Min", "Fast_Total_Prod_Month_kWh_Max", "Fast_Total_Prod_Month_kWh_Count", "Fast_ActivePowerOutPWCSell_kW", "Fast_ActivePowerOutPWCSell_kW_StdDev", "Fast_ActivePowerOutPWCSell_kW_Min", "Fast_ActivePowerOutPWCSell_kW_Max", "Fast_ActivePowerOutPWCSell_kW_Count", "Fast_Frequency_Hz", "Fast_Frequency_Hz_StdDev", "Fast_Frequency_Hz_Min", "Fast_Frequency_Hz_Max", "Fast_Frequency_Hz_Count", "Slow_TempG1L2", "Slow_TempG1L2_StdDev", "Slow_TempG1L2_Min", "Slow_TempG1L2_Max", "Slow_TempG1L2_Count", "Slow_TempG1L3", "Slow_TempG1L3_StdDev", "Slow_TempG1L3_Min", "Slow_TempG1L3_Max", "Slow_TempG1L3_Count", "Slow_TempGearBoxHSSDE", "Slow_TempGearBoxHSSDE_StdDev", "Slow_TempGearBoxHSSDE_Min", "Slow_TempGearBoxHSSDE_Max", "Slow_TempGearBoxHSSDE_Count", "Slow_TempGearBoxIMSNDE", "Slow_TempGearBoxIMSNDE_StdDev", "Slow_TempGearBoxIMSNDE_Min", "Slow_TempGearBoxIMSNDE_Max", "Slow_TempGearBoxIMSNDE_Count", "Slow_TempOutdoor", "Slow_TempOutdoor_StdDev", "Slow_TempOutdoor_Min", "Slow_TempOutdoor_Max", "Slow_TempOutdoor_Count", "Fast_PitchAccuV3", "Fast_PitchAccuV3_StdDev", "Fast_PitchAccuV3_Min", "Fast_PitchAccuV3_Max", "Fast_PitchAccuV3_Count", "Slow_TotalTurbineActiveHours", "Slow_TotalTurbineActiveHours_StdDev", "Slow_TotalTurbineActiveHours_Min", "Slow_TotalTurbineActiveHours_Max", "Slow_TotalTurbineActiveHours_Count", "Slow_TotalTurbineOKHours", "Slow_TotalTurbineOKHours_StdDev", "Slow_TotalTurbineOKHours_Min", "Slow_TotalTurbineOKHours_Max", "Slow_TotalTurbineOKHours_Count", "Slow_TotalTurbineTimeAllHours", "Slow_TotalTurbineTimeAllHours_StdDev", "Slow_TotalTurbineTimeAllHours_Min", "Slow_TotalTurbineTimeAllHours_Max", "Slow_TotalTurbineTimeAllHours_Count", "Slow_TempG1L1", "Slow_TempG1L1_StdDev", "Slow_TempG1L1_Min", "Slow_TempG1L1_Max", "Slow_TempG1L1_Count", "Slow_TempGearBoxOilSump", "Slow_TempGearBoxOilSump_StdDev", "Slow_TempGearBoxOilSump_Min", "Slow_TempGearBoxOilSump_Max", "Slow_TempGearBoxOilSump_Count", "Fast_PitchAccuV2", "Fast_PitchAccuV2_StdDev", "Fast_PitchAccuV2_Min", "Fast_PitchAccuV2_Max", "Fast_PitchAccuV2_Count", "Slow_TotalGridOkHours", "Slow_TotalGridOkHours_StdDev", "Slow_TotalGridOkHours_Min", "Slow_TotalGridOkHours_Max", "Slow_TotalGridOkHours_Count", "Slow_TotalActPowerOut_kWh", "Slow_TotalActPowerOut_kWh_StdDev", "Slow_TotalActPowerOut_kWh_Min", "Slow_TotalActPowerOut_kWh_Max", "Slow_TotalActPowerOut_kWh_Count", "Fast_YawService", "Fast_YawService_StdDev", "Fast_YawService_Min", "Fast_YawService_Max", "Fast_YawService_Count", "Fast_YawAngle", "Fast_YawAngle_StdDev", "Fast_YawAngle_Min", "Fast_YawAngle_Max", "Fast_YawAngle_Count", "Slow_CapableInductivePwrFactor", "Slow_CapableInductivePwrFactor_StdDev", "Slow_CapableInductivePwrFactor_Min", "Slow_CapableInductivePwrFactor_Max", "Slow_CapableInductivePwrFactor_Count", "Slow_TempGearBoxHSSNDE", "Slow_TempGearBoxHSSNDE_StdDev", "Slow_TempGearBoxHSSNDE_Min", "Slow_TempGearBoxHSSNDE_Max", "Slow_TempGearBoxHSSNDE_Count", "Slow_TempHubBearing", "Slow_TempHubBearing_StdDev", "Slow_TempHubBearing_Min", "Slow_TempHubBearing_Max", "Slow_TempHubBearing_Count", "Slow_TotalG1ActiveHours", "Slow_TotalG1ActiveHours_StdDev", "Slow_TotalG1ActiveHours_Min", "Slow_TotalG1ActiveHours_Max", "Slow_TotalG1ActiveHours_Count", "Slow_TotalActPowerOutG1_kWh", "Slow_TotalActPowerOutG1_kWh_StdDev", "Slow_TotalActPowerOutG1_kWh_Min", "Slow_TotalActPowerOutG1_kWh_Max", "Slow_TotalActPowerOutG1_kWh_Count", "Slow_TotalReactPowerInG1_kVArh", "Slow_TotalReactPowerInG1_kVArh_StdDev", "Slow_TotalReactPowerInG1_kVArh_Min", "Slow_TotalReactPowerInG1_kVArh_Max", "Slow_TotalReactPowerInG1_kVArh_Count", "Slow_NacelleDrill", "Slow_NacelleDrill_StdDev", "Slow_NacelleDrill_Min", "Slow_NacelleDrill_Max", "Slow_NacelleDrill_Count", "Slow_TempGearBoxIMSDE", "Slow_TempGearBoxIMSDE_StdDev", "Slow_TempGearBoxIMSDE_Min", "Slow_TempGearBoxIMSDE_Max", "Slow_TempGearBoxIMSDE_Count", "Fast_Total_Operating_hrs", "Fast_Total_Operating_hrs_StdDev", "Fast_Total_Operating_hrs_Min", "Fast_Total_Operating_hrs_Max", "Fast_Total_Operating_hrs_Count", "Slow_TempNacelle", "Slow_TempNacelle_StdDev", "Slow_TempNacelle_Min", "Slow_TempNacelle_Max", "Slow_TempNacelle_Count", "Fast_Total_Grid_OK_hrs", "Fast_Total_Grid_OK_hrs_StdDev", "Fast_Total_Grid_OK_hrs_Min", "Fast_Total_Grid_OK_hrs_Max", "Fast_Total_Grid_OK_hrs_Count", "Fast_Total_WTG_OK_hrs", "Fast_Total_WTG_OK_hrs_StdDev", "Fast_Total_WTG_OK_hrs_Min", "Fast_Total_WTG_OK_hrs_Max", "Fast_Total_WTG_OK_hrs_Count", "Slow_TempCabinetTopBox", "Slow_TempCabinetTopBox_StdDev", "Slow_TempCabinetTopBox_Min", "Slow_TempCabinetTopBox_Max", "Slow_TempCabinetTopBox_Count", "Slow_TempGeneratorBearingNDE", "Slow_TempGeneratorBearingNDE_StdDev", "Slow_TempGeneratorBearingNDE_Min", "Slow_TempGeneratorBearingNDE_Max", "Slow_TempGeneratorBearingNDE_Count", "Fast_Total_Access_hrs", "Fast_Total_Access_hrs_StdDev", "Fast_Total_Access_hrs_Min", "Fast_Total_Access_hrs_Max", "Fast_Total_Access_hrs_Count", "Slow_TempBottomPowerSection", "Slow_TempBottomPowerSection_StdDev", "Slow_TempBottomPowerSection_Min", "Slow_TempBottomPowerSection_Max", "Slow_TempBottomPowerSection_Count", "Slow_TempGeneratorBearingDE", "Slow_TempGeneratorBearingDE_StdDev", "Slow_TempGeneratorBearingDE_Min", "Slow_TempGeneratorBearingDE_Max", "Slow_TempGeneratorBearingDE_Count", "Slow_TotalReactPowerIn_kVArh", "Slow_TotalReactPowerIn_kVArh_StdDev", "Slow_TotalReactPowerIn_kVArh_Min", "Slow_TotalReactPowerIn_kVArh_Max", "Slow_TotalReactPowerIn_kVArh_Count", "Slow_TempBottomControlSection", "Slow_TempBottomControlSection_StdDev", "Slow_TempBottomControlSection_Min", "Slow_TempBottomControlSection_Max", "Slow_TempBottomControlSection_Count", "Slow_TempConv1", "Slow_TempConv1_StdDev", "Slow_TempConv1_Min", "Slow_TempConv1_Max", "Slow_TempConv1_Count", "Fast_ActivePowerRated_kW", "Fast_ActivePowerRated_kW_StdDev", "Fast_ActivePowerRated_kW_Min", "Fast_ActivePowerRated_kW_Max", "Fast_ActivePowerRated_kW_Count", "Fast_NodeIP", "Fast_NodeIP_StdDev", "Fast_NodeIP_Min", "Fast_NodeIP_Max", "Fast_NodeIP_Count", "Fast_PitchSpeed1", "Fast_PitchSpeed1_StdDev", "Fast_PitchSpeed1_Min", "Fast_PitchSpeed1_Max", "Fast_PitchSpeed1_Count", "Slow_CFCardSize", "Slow_CFCardSize_StdDev", "Slow_CFCardSize_Min", "Slow_CFCardSize_Max", "Slow_CFCardSize_Count", "Slow_CPU_Number", "Slow_CPU_Number_StdDev", "Slow_CPU_Number_Min", "Slow_CPU_Number_Max", "Slow_CPU_Number_Count", "Slow_CFCardSpaceLeft", "Slow_CFCardSpaceLeft_StdDev", "Slow_CFCardSpaceLeft_Min", "Slow_CFCardSpaceLeft_Max", "Slow_CFCardSpaceLeft_Count", "Slow_TempBottomCapSection", "Slow_TempBottomCapSection_StdDev", "Slow_TempBottomCapSection_Min", "Slow_TempBottomCapSection_Max", "Slow_TempBottomCapSection_Count", "Slow_RatedPower", "Slow_RatedPower_StdDev", "Slow_RatedPower_Min", "Slow_RatedPower_Max", "Slow_RatedPower_Count", "Slow_TempConv3", "Slow_TempConv3_StdDev", "Slow_TempConv3_Min", "Slow_TempConv3_Max", "Slow_TempConv3_Count", "Slow_TempConv2", "Slow_TempConv2_StdDev", "Slow_TempConv2_Min", "Slow_TempConv2_Max", "Slow_TempConv2_Count", "Slow_TotalActPowerIn_kWh", "Slow_TotalActPowerIn_kWh_StdDev", "Slow_TotalActPowerIn_kWh_Min", "Slow_TotalActPowerIn_kWh_Max", "Slow_TotalActPowerIn_kWh_Count", "Slow_TotalActPowerInG1_kWh", "Slow_TotalActPowerInG1_kWh_StdDev", "Slow_TotalActPowerInG1_kWh_Min", "Slow_TotalActPowerInG1_kWh_Max", "Slow_TotalActPowerInG1_kWh_Count", "Slow_TotalActPowerInG2_kWh", "Slow_TotalActPowerInG2_kWh_StdDev", "Slow_TotalActPowerInG2_kWh_Min", "Slow_TotalActPowerInG2_kWh_Max", "Slow_TotalActPowerInG2_kWh_Count", "Slow_TotalActPowerOutG2_kWh", "Slow_TotalActPowerOutG2_kWh_StdDev", "Slow_TotalActPowerOutG2_kWh_Min", "Slow_TotalActPowerOutG2_kWh_Max", "Slow_TotalActPowerOutG2_kWh_Count", "Slow_TotalG2ActiveHours", "Slow_TotalG2ActiveHours_StdDev", "Slow_TotalG2ActiveHours_Min", "Slow_TotalG2ActiveHours_Max", "Slow_TotalG2ActiveHours_Count", "Slow_TotalReactPowerInG2_kVArh", "Slow_TotalReactPowerInG2_kVArh_StdDev", "Slow_TotalReactPowerInG2_kVArh_Min", "Slow_TotalReactPowerInG2_kVArh_Max", "Slow_TotalReactPowerInG2_kVArh_Count", "Slow_TotalReactPowerOut_kVArh", "Slow_TotalReactPowerOut_kVArh_StdDev", "Slow_TotalReactPowerOut_kVArh_Min", "Slow_TotalReactPowerOut_kVArh_Max", "Slow_TotalReactPowerOut_kVArh_Count", "Slow_UTCoffset_int", "Slow_UTCoffset_int_StdDev", "Slow_UTCoffset_int_Min", "Slow_UTCoffset_int_Max", "Slow_UTCoffset_int_Count"}
	// 	tablename = new(ScadaDataHFD).TableName()
	case "DowntimeEvent":
		header = []string{"Turbine", "TimeStart", "TimeEnd", "Down Grid", "Down Environment", "Down Machine", "Alarm Description", "Duration", "Reduce Availability"}
		tablename = new(EventDown).TableName()
		separator = ""
		database = DB().Connection
	case "EventRaw":
		header = []string{"TimeStamp", "Project Name", "Turbine", "Event Type", "Alarm Description", "Turbine Status", "Brake Type", "Brake Program", "Alarm Id", "Alarm Toggle"}
		tablename = new(EventRaw).TableName()
		separator = ""
		database = DB().Connection
	case "MetTower":
		header = []string{"TimeStamp", "WindDirNo", "VHubWS90mAvg", "VHubWS90mMax", "VHubWS90mMin", "VHubWS90mStdDev", "VHubWS90mCount", "VRefWS88mAvg", "VRefWS88mMax", "VRefWS88mMin", "VRefWS88mStdDev", "VRefWS88mCount", "VTipWS42mAvg", "VTipWS42mMax", "VTipWS42mMin", "VTipWS42mStdDev", "VTipWS42mCount", "DHubWD88mAvg", "DHubWD88mMax", "DHubWD88mMin", "DHubWD88mStdDev", "DHubWD88mCount", "DRefWD86mAvg", "DRefWD86mMax", "DRefWD86mMin", "DRefWD86mStdDev", "DRefWD86mCount", "THubHHubHumid855mAvg", "THubHHubHumid855mMax", "THubHHubHumid855mMin", "THubHHubHumid855mStdDev", "THubHHubHumid855mCount", "TRefHRefHumid855mAvg", "TRefHRefHumid855mMax", "TRefHRefHumid855mMin", "TRefHRefHumid855mStdDev", "TRefHRefHumid855mCount", "THubHHubTemp855mAvg", "THubHHubTemp855mMax", "THubHHubTemp855mMin", "THubHHubTemp855mStdDev", "THubHHubTemp855mCount", "TRefHRefTemp855mAvg", "TRefHRefTemp855mMax", "TRefHRefTemp855mMin", "TRefHRefTemp855mStdDev", "TRefHRefTemp855mCount", "BaroAirPress855mAvg", "BaroAirPress855mMax", "BaroAirPress855mMin", "BaroAirPress855mStdDev", "BaroAirPress855mCount", "YawAngleVoltageAvg", "YawAngleVoltageMax", "YawAngleVoltageMin", "YawAngleVoltageStdDev", "YawAngleVoltageCount", "OtherSensorVoltageAI1Avg", "OtherSensorVoltageAI1Max", "OtherSensorVoltageAI1Min", "OtherSensorVoltageAI1StdDev", "OtherSensorVoltageAI1Count", "OtherSensorVoltageAI2Avg", "OtherSensorVoltageAI2Max", "OtherSensorVoltageAI2Min", "OtherSensorVoltageAI2StdDev", "OtherSensorVoltageAI2Count", "OtherSensorVoltageAI3Avg", "OtherSensorVoltageAI3Max", "OtherSensorVoltageAI3Min", "OtherSensorVoltageAI3StdDev", "OtherSensorVoltageAI3Count", "OtherSensorVoltageAI4Avg", "OtherSensorVoltageAI4Max", "OtherSensorVoltageAI4Min", "OtherSensorVoltageAI4StdDev", "OtherSensorVoltageAI4Count", "GenRPMCurrentAvg", "GenRPMCurrentMax", "GenRPMCurrentMin", "GenRPMCurrentStdDev", "GenRPMCurrentCount", "WS_SCSCurrentAvg", "WS_SCSCurrentMax", "WS_SCSCurrentMin", "WS_SCSCurrentStdDev", "WS_SCSCurrentCount", "RainStatusCount", "RainStatusSum", "OtherSensor2StatusIO1Avg", "OtherSensor2StatusIO1Max", "OtherSensor2StatusIO1Min", "OtherSensor2StatusIO1StdDev", "OtherSensor2StatusIO1Count", "OtherSensor2StatusIO2Avg", "OtherSensor2StatusIO2Max", "OtherSensor2StatusIO2Min", "OtherSensor2StatusIO2StdDev", "OtherSensor2StatusIO2Count", "OtherSensor2StatusIO3Avg", "OtherSensor2StatusIO3Max", "OtherSensor2StatusIO3Min", "OtherSensor2StatusIO3StdDev", "OtherSensor2StatusIO3Count", "OtherSensor2StatusIO4Avg", "OtherSensor2StatusIO4Max", "OtherSensor2StatusIO4Min", "OtherSensor2StatusIO4StdDev", "OtherSensor2StatusIO4Count", "OtherSensor2StatusIO5Avg", "OtherSensor2StatusIO5Max", "OtherSensor2StatusIO5Min", "OtherSensor2StatusIO5StdDev", "OtherSensor2StatusIO5Count", "A1Avg", "A1Max", "A1Min", "A1StdDev", "A1Count", "A2Avg", "A2Max", "A2Min", "A2StdDev", "A2Count", "A3Avg", "A3Max", "A3Min", "A3StdDev", "A3Count", "A4Avg", "A4Max", "A4Min", "A4StdDev", "A4Count", "A5Avg", "A5Max", "A5Min", "A5StdDev", "A5Count", "A6Avg", "A6Max", "A6Min", "A6StdDev", "A6Count", "A7Avg", "A7Max", "A7Min", "A7StdDev", "A7Count", "A8Avg", "A8Max", "A8Min", "A8StdDev", "A8Count", "A9Avg", "A9Max", "A9Min", "A9StdDev", "A9Count", "A10Avg", "A10Max", "A10Min", "A10StdDev", "A10Count", "AC1Avg", "AC1Max", "AC1Min", "AC1StdDev", "AC1Count", "AC2Avg", "AC2Max", "AC2Min", "AC2StdDev", "AC2Count", "C1Avg", "C1Max", "C1Min", "C1StdDev", "C1Count", "C2Avg", "C2Max", "C2Min", "C2StdDev", "C2Count", "C3Avg", "C3Max", "C3Min", "C3StdDev", "C3Count", "D1Avg", "D1Max", "D1Min", "D1StdDev", "M1_1Avg", "M1_1Max", "M1_1Min", "M1_1StdDev", "M1_1Count", "M1_2Avg", "M1_2Max", "M1_2Min", "M1_2StdDev", "M1_2Count", "M1_3Avg", "M1_3Max", "M1_3Min", "M1_3StdDev", "M1_3Count", "M1_4Avg", "M1_4Max", "M1_4Min", "M1_4StdDev", "M1_4Count", "M1_5Avg", "M1_5Max", "M1_5Min", "M1_5StdDev", "M1_5Count", "M2_1Avg", "M2_1Max", "M2_1Min", "M2_1StdDev", "M2_1Count", "M2_2Avg", "M2_2Max", "M2_2Min", "M2_2StdDev", "M2_2Count", "M2_3Avg", "M2_3Max", "M2_3Min", "M2_3StdDev", "M2_3Count", "M2_4Avg", "M2_4Max", "M2_4Min", "M2_4StdDev", "M2_4Count", "M2_5Avg", "M2_5Max", "M2_5Min", "M2_5StdDev", "M2_5Count", "M2_6Avg", "M2_6Max", "M2_6Min", "M2_6StdDev", "M2_6Count", "M2_7Avg", "M2_7Max", "M2_7Min", "M2_7StdDev", "M2_7Count", "M2_8Avg", "M2_8Max", "M2_8Min", "M2_8StdDev", "M2_8Count", "VAvg", "VMax", "VMin", "IAvg", "IMax", "IMin", "T", "Addr", "WindDirDesc", "WSCategoryNo", "WSCategoryDesc"}
		tablename = new(MetTower).TableName()
		separator = ""
		p.Project = ""
		database = DB().Connection
	case "DowntimeEventHFD":
		header = []string{"Turbine", "Time Start", "Time End", "Duration", "Break Down Group", "Turbine State", "Alarm Code", "Alarm Description"}
		tablename = new(AlarmHFD).TableName()
		separator = ""
		database = DBRealtime()
	}

	for _, val := range header {
		switch val {
		case "Break Down Group":
			val = "bdgroup"
		case "Alarm Description":
			val = "alarmdesc"
		}
		fieldList = append(fieldList, strings.ToLower(strings.Replace(strings.TrimSuffix(val, " "), " ", separator, -69)))
	}

	TimeCreate := time.Now().Format("2006-01-02_150405")
	CreateDateTime := typeExcel + TimeCreate

	if err := os.RemoveAll("web/assets/Excel/" + typeExcel + "/"); err != nil {
		tk.Println(err)
	}

	query := database.NewQuery().From(tablename).Where(dbox.And(filter...))
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
	defer csr.Close()
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	data := make([]tk.M, 0)
	e = csr.Fetch(&data, 0, false)

	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	//web/assets/Excel/

	if _, err := os.Stat("web/assets/Excel/" + typeExcel + "/"); os.IsNotExist(err) {
		os.MkdirAll("web/assets/Excel/"+typeExcel+"/", 0777)
	}
	turbinename, err := helper.GetTurbineNameList(p.Project)
	if err != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	DeserializeData(data, typeExcel, CreateDateTime, header, fieldList, turbinename)
	pathDownload = "res/Excel/" + typeExcel + "/" + CreateDateTime + ".xlsx"

	return helper.CreateResult(true, pathDownload, "success")
}

func (m *DataBrowserController) GenExcelCustom10Minutes(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	p := new(helper.Payloads)
	e := k.GetPayload(&p)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	p.Misc.Set("knot_data", k)
	filter, _ := p.ParseFilter()
	typeExcel := strings.Split(p.Misc.GetString("tipe"), "Custom")[0]

	arrscadaoem := []string{"_id"}
	// arrmettower := []string{}
	headerList := []string{}
	fieldList := []string{}
	source := "ScadaDataHFD" /* to filter field list from payload ColumnList */
	// timestamp := "timestamp"
	// tablename := new(ScadaDataHFD).TableName()
	tablename := "Scada10MinHFD"
	ids := ""
	switch typeExcel {
	case "ScadaOEM":
		tablename = new(ScadaDataOEM).TableName()
		source = "ScadaDataOEM"
		arrscadaoem = append(arrscadaoem, "timestamputc")
		// timestamp = "timestamputc"
	}

	// istimestamp := false
	if p.Custom.Has("ColumnList") {
		for _, _val := range p.Custom["ColumnList"].([]interface{}) {
			_tkm, _ := tk.ToM(_val)
			ids = strings.ToLower(_tkm.GetString("_id"))
			if _tkm.GetString("source") == source {
				arrscadaoem = append(arrscadaoem, ids)
				// if ids == "timestamp" {
				// 	istimestamp = true
				// }
			}
			/*else if _tkm.GetString("source") == "MetTower" {
				arrmettower = append(arrmettower, _tkm.GetString("_id"))
			}*/
			headerList = append(headerList, _tkm.GetString("label"))
			fieldList = append(fieldList, ids)
		}
	}

	/*query := DB().Connection.NewQuery().
	Select(arrscadaoem...).
	From(tablename).
	Where(dbox.And(filter...))*/

	projection := map[string]int{}
	for _, val := range arrscadaoem {
		projection[val] = 1
	}

	matches := []tk.M{}
	for _, f := range filter {
		value := f.Value
		if f.Field == "timestamp" {
			value = value.(time.Time).UTC()
		}
		matches = append(matches, tk.M{
			f.Field: tk.M{f.Op: value},
		})
	}

	pipes := []tk.M{}
	sortList := map[string]int{}
	if len(p.Sort) > 0 {
		for _, val := range p.Sort {
			if val.Dir == "desc" {
				sortList[strings.ToLower(val.Field)] = -1
			} else {
				sortList[strings.ToLower(val.Field)] = 1
			}
		}
		pipes = append(pipes, tk.M{"$sort": sortList})
	}
	pipes = append(pipes, tk.M{"$match": tk.M{"$and": matches}})
	pipes = append(pipes, tk.M{"$project": projection})
	/*if len(p.Sort) > 0 {
		var arrsort []string
		for _, val := range p.Sort {
			if val.Dir == "desc" {
				arrsort = append(arrsort, strings.ToLower("-"+val.Field))
			} else {
				arrsort = append(arrsort, strings.ToLower(val.Field))
			}
		}
		query = query.Order(arrsort...)
	}*/

	// csr, e := query.Cursor(nil)
	csr, e := DB().Connection.NewQuery().
		From(tablename).Command("pipe", pipes).Cursor(nil)
	defer csr.Close()
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	/*results := make([]tk.M, 0)
	e = csr.Fetch(&results, 0, false)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}*/
	results := make([]tk.M, 0)
	result := tk.M{}
	for {
		result = tk.M{}
		e = csr.Fetch(&result, 1, false)
		if e != nil {
			break
		}
		results = append(results, result)
	}

	// arrmettowercond := []interface{}{}

	// config := lh.ReadConfig()
	// loc, err := time.LoadLocation(config["ReadTimeLoc"])
	// if err != nil {
	// 	tk.Printfn("Get time in %s found %s", config["ReadTimeLoc"], err.Error())
	// }

	// for i, val := range results {
	// 	if val.Has("timestamputc") {
	// 		strangeTime := val.Get("timestamputc", time.Time{}).(time.Time).UTC().In(loc)
	// 		itime := time.Date(strangeTime.Year(), strangeTime.Month(), strangeTime.Day(),
	// 			strangeTime.Hour(), strangeTime.Minute(), strangeTime.Second(), strangeTime.Nanosecond(), time.UTC)
	// 		// arrmettowercond = append(arrmettowercond, itime)
	// 		val.Set("timestamputc", itime)
	// 		results[i] = val
	// 	}
	// 	if istimestamp {
	// 		itime := val.Get("timestamp", time.Time{}).(time.Time).UTC()
	// 		/*if typeExcel == "ScadaHFD" {
	// 			arrmettowercond = append(arrmettowercond, itime)
	// 		}*/
	// 		val.Set("timestamp", itime)
	// 		results[i] = val
	// 	}
	// }

	/*tkmmet := tk.M{}
	if len(arrmettower) > 0 && len(arrmettowercond) > 0 {
		arrmettower = append(arrmettower, "timestamp")
		queryMet := DB().Connection.NewQuery().Select(arrmettower...).From("MetTower")

		filterMet := dbox.In("timestamp", arrmettowercond...)
		for _, val := range filter {
			if val.Field == "projectname" {
				filterMet = dbox.And(dbox.Eq(val.Field, val.Value), filterMet)
			}
		}
		csrMet, _e := queryMet.Where(filterMet).Cursor(nil)
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
			itime := val.Get(timestamp, time.Time{}).(time.Time).UTC().String()
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
	}*/

	var pathDownload string
	TimeCreate := time.Now().Format("2006-01-02_150405")
	CreateDateTime := typeExcel + TimeCreate

	if err := os.RemoveAll("web/assets/Excel/" + typeExcel + "/"); err != nil {
		tk.Println(err)
	}

	if _, err := os.Stat("web/assets/Excel/" + typeExcel + "/"); os.IsNotExist(err) {
		os.MkdirAll("web/assets/Excel/"+typeExcel+"/", 0777)
	}

	turbineName, err := helper.GetTurbineNameList(p.Project)
	if err != nil {
		return helper.CreateResult(false, nil, e.Error())
	}
	DeserializeData(results, typeExcel, CreateDateTime, headerList, fieldList, turbineName)
	pathDownload = "res/Excel/" + typeExcel + "/" + CreateDateTime + ".xlsx"

	return helper.CreateResult(true, pathDownload, "success")
}

func (m *DataBrowserController) GenExcelCustom10FarmWise(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	p := new(helper.Payloads)
	e := k.GetPayload(&p)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	p.Misc.Set("knot_data", k)
	filter, _ := p.ParseFilter()
	typeExcel := strings.Split(p.Misc.GetString("tipe"), "Custom")[0]

	arrscadaoem := []string{}
	headerList := []string{}
	fieldList := []string{}
	source := "ScadaDataHFD" /* to filter field list from payload ColumnList */
	tablename := "Scada10MinHFD"
	ids := ""

	if p.Custom.Has("ColumnList") {
		for _, _val := range p.Custom["ColumnList"].([]interface{}) {
			_tkm, _ := tk.ToM(_val)
			ids = strings.ToLower(_tkm.GetString("_id"))
			if _tkm.GetString("source") == source {
				arrscadaoem = append(arrscadaoem, ids)
			}
			headerList = append(headerList, _tkm.GetString("label"))
			fieldList = append(fieldList, ids)
		}
	}

	projection := map[string]int{}
	fproject := map[string]int{"projectname": 1, "timestamp": 1}
	agroups := tk.M{
		"_id": tk.M{"projectname": "$projectname", "timestamp": "$timestamp"},
	}
	for _, val := range arrscadaoem {
		if val == "turbine" {
			continue
		}

		projection[val] = 1
		if val == "projectname" || val == "timestamp" {
			continue
		}

		// fproject[tk.Sprintf("%s_sum", val)] = 1
		// fproject[tk.Sprintf("%s_count", val)] = 1

		// agroups.Set(val+"_sum", tk.M{"$sum": tk.Sprintf("$%s_sum", val)})
		// agroups.Set(val+"_count", tk.M{"$sum": tk.Sprintf("$%s_count", val)})

		fproject[val] = 1
		agroups.Set(val, tk.M{"$avg": tk.Sprintf("$%s", val)})
	}

	matches := []tk.M{}
	for _, f := range filter {
		if f.Field == "turbine" {
			continue
		}

		value := f.Value
		if f.Field == "timestamp" {
			value = value.(time.Time).UTC()
		}
		matches = append(matches, tk.M{
			f.Field: tk.M{f.Op: value},
		})
	}

	pipes := []tk.M{}
	// sortList := map[string]int{}
	// if len(p.Sort) > 0 {
	// 	for _, val := range p.Sort {
	// 		if val.Dir == "desc" {
	// 			sortList[strings.ToLower(val.Field)] = -1
	// 		} else {
	// 			sortList[strings.ToLower(val.Field)] = 1
	// 		}
	// 	}
	// 	pipes = append(pipes, tk.M{"$sort": sortList})
	// }
	// pipes = append(pipes, tk.M{"$match": tk.M{"$and": matches}})
	// pipes = append(pipes, tk.M{"$project": projection})

	pipes = append(pipes, tk.M{"$match": tk.M{"$and": matches}})
	pipes = append(pipes, tk.M{"$project": fproject})
	pipes = append(pipes, tk.M{"$group": agroups})

	// csr, e := query.Cursor(nil)
	csr, e := DB().Connection.NewQuery().
		From(tablename).Command("pipe", pipes).Cursor(nil)
	defer csr.Close()
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	results := make([]tk.M, 0)
	result := tk.M{}
	for {
		result = tk.M{}
		e = csr.Fetch(&result, 1, false)
		if e != nil {
			break
		}

		// idtk, _ := tk.ToM(result.Get("_id"))
		idtk := result.Get("_id", tk.M{}).(tk.M)
		// resitem := tk.M{}
		// resitem.Set("projectname", idtk["projectname"])
		// resitem.Set("timestamp", idtk["timestamp"])
		// for field, _ := range projection {
		// 	if field == "turbine" || field == "projectname" || field == "timestamp" {
		// 		continue
		// 	}
		// 	resitem.Set(field, tk.Div(result.GetFloat64(field+"_sum"), result.GetFloat64(field+"_count")))
		// }
		// tk.Println(result)
		result.Set("projectname", idtk["projectname"])
		result.Set("timestamp", idtk.Get("timestamp", time.Time{}).(time.Time))

		results = append(results, result)
	}

	var pathDownload string
	TimeCreate := time.Now().Format("2006-01-02_150405")
	CreateDateTime := typeExcel + TimeCreate

	if err := os.RemoveAll("web/assets/Excel/" + typeExcel + "/"); err != nil {
		tk.Println(err)
	}

	if _, err := os.Stat("web/assets/Excel/" + typeExcel + "/"); os.IsNotExist(err) {
		os.MkdirAll("web/assets/Excel/"+typeExcel+"/", 0777)
	}

	turbineName, err := helper.GetTurbineNameList(p.Project)
	if err != nil {
		return helper.CreateResult(false, nil, e.Error())
	}
	// for _, res := range results {
	// 	tk.Println(" == ", res)
	// }
	DeserializeData(results, typeExcel, CreateDateTime, headerList, fieldList, turbineName)
	pathDownload = "res/Excel/" + typeExcel + "/" + CreateDateTime + ".xlsx"

	return helper.CreateResult(true, pathDownload, "success")
}

// Deserialize

func FormatThousandSeparator(floatString string) string {
	splitString := strings.Split(floatString, ".")
	beforeComma := splitString[0]
	afterComma := ""
	if len(splitString) > 1 {
		afterComma = splitString[1]
	} else {
		afterComma = "00"
	}
	resultNumber := ""
	idx := 0
	for i := len(beforeComma) - 1; i >= 0; i-- {
		if idx > 0 {
			if idx%3 == 0 {
				resultNumber = "," + resultNumber
			}
		}
		resultNumber = string(beforeComma[i]) + resultNumber
		idx++
	}
	resultNumber += "." + afterComma
	return resultNumber
}

func DeserializeData(data []tk.M, typeExcel, CreateDateTime string, header, fieldList []string, turbinename map[string]string) error {
	filename := ""
	filename = "web/assets/Excel/" + typeExcel + "/" + CreateDateTime + ".xlsx"

	file := x.NewFile()
	sheet, _ := file.AddSheet("Sheet1")
	// floatString := ""
	dataType := ""

	for i, each := range data {
		if i == 0 {
			rowHeader := sheet.AddRow()
			for _, hdr := range header {

				cell := rowHeader.AddCell()
				cell.Value = hdr
			}
		}

		// tk.Println(" == ", each)
		rowContent := sheet.AddRow()
		cell := rowContent.AddCell()
		for idx, field := range fieldList {
			if idx > 0 {
				cell = rowContent.AddCell()
			}
			if each.Has(field) && each[field] != nil {
				switch field {
				case "timestamp", "timestamputc", "timestart", "timeend":
					cell.Value = each[field].(time.Time).UTC().Format("2006-01-02 15:04:05")
				case "turbine":
					cell.Value = turbinename[each.GetString(field)]
				default:
					dataType = reflect.Indirect(reflect.ValueOf(each[field])).Type().Name()
					switch dataType {
					case "float64":
						/*floatString = tk.Sprintf("%.2f", each.GetFloat64(field))
						if each.GetFloat64(field) == -999999 {
							floatString = "-"
						}
						if floatString != "-" {
							floatString = FormatThousandSeparator(floatString)
						}
						cell.Value = floatString*/
						value := each.GetFloat64(field)
						if value != -999999 {
							cell.SetFloat(value)
						}
					case "int":
						/*floatString = tk.Sprintf("%d", each.GetInt(field))
						if each.GetInt(field) == -999999 {
							floatString = "-"
						}
						cell.Value = floatString*/
						value := each.GetInt(field)
						if value != -999999 {
							cell.SetInt(value)
						}
					case "string":
						cell.Value = each.GetString(field)
					case "bool":
						// cell.Value = strconv.FormatBool(each[field].(bool))
						cell.SetBool(each[field].(bool))
					}
				}
			}
		}
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
