package controller

import (
	. "eaciit/wfdemo-git/library/core"
	lh "eaciit/wfdemo-git/library/helper"
	. "eaciit/wfdemo-git/library/models"
	"eaciit/wfdemo-git/web/helper"
	"reflect"
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

var (
	_amettower_label = []string{"V Hub WS 90m Avg", "V Hub WS 90m Std Dev", "V Ref WS 88m Avg", "V Ref WS 88m Std Dev",
		"V Tip WS 42m Avg", "V Tip WS 42m Std Dev", "D Hub WD 88m Avg", "D Hub WD 88m Std Dev", "D Ref WD 86m Avg",
		"D Ref WD 86m Std Dev", "T Hub & H Hub Humid 85m Avg", "T Hub & H Hub Humid 85m Std Dev", "T Ref & H Ref Humid 85.5m Avg", "T Ref & H Ref Humid 85.5m Std Dev",
		"T Hub & H Hub Temp 85.5m Avg", "T Hub & H Hub Temp 85.5m Std Dev", "T Ref & H Ref Temp 85.5 Avg", "T Ref & H Ref Temp 85.5 Std Dev", "Baro Air Pressure 85.5m Avg", "Baro Air Pressure 85.5m Std Dev",
	}

	_amettower_field = []string{"vhubws90mavg", "vhubws90mstddev", "vrefws88mavg", "vrefws88mstddev", "vtipws42mavg",
		"vtipws42mstddev", "dhubwd88mavg", "dhubwd88mstddev", "drefwd86mavg", "drefwd86mstddev",
		"thubhhubhumid855mavg", "thubhhubhumid855mstddev", "trefhrefhumid855mavg", "trefhrefhumid855mstddev", "thubhhubtemp855mavg",
		"thubhhubtemp855mstddev", "trefhreftemp855mavg", "trefhreftemp855mstddev", "baroairpress855mavg", "baroairpress855mstddev",
	}
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

	// _ashfd_label := []string{"Fast ActivePower kW", "Fast ActivePower kW StdDev", "Fast ActivePower kW Min", "Fast ActivePower kW Max", "Fast ActivePower kW Count", "Fast WindSpeed ms", "Fast WindSpeed ms StdDev", "Fast WindSpeed ms Min", "Fast WindSpeed ms Max", "Fast WindSpeed ms Count", "Slow NacellePos", "Slow NacellePos StdDev", "Slow NacellePos Min", "Slow NacellePos Max", "Slow NacellePos Count", "Slow WindDirection", "Slow WindDirection StdDev", "Slow WindDirection Min", "Slow WindDirection Max", "Slow WindDirection Count", "Fast CurrentL3", "Fast CurrentL3 StdDev", "Fast CurrentL3 Min", "Fast CurrentL3 Max", "Fast CurrentL3 Count", "Fast CurrentL1", "Fast CurrentL1 StdDev", "Fast CurrentL1 Min", "Fast CurrentL1 Max", "Fast CurrentL1 Count", "Fast ActivePowerSetpoint kW", "Fast ActivePowerSetpoint kW StdDev", "Fast ActivePowerSetpoint kW Min", "Fast ActivePowerSetpoint kW Max", "Fast ActivePowerSetpoint kW Count", "Fast CurrentL2", "Fast CurrentL2 StdDev", "Fast CurrentL2 Min", "Fast CurrentL2 Max", "Fast CurrentL2 Count", "Fast DrTrVibValue", "Fast DrTrVibValue StdDev", "Fast DrTrVibValue Min", "Fast DrTrVibValue Max", "Fast DrTrVibValue Count", "Fast GenSpeed RPM", "Fast GenSpeed RPM StdDev", "Fast GenSpeed RPM Min", "Fast GenSpeed RPM Max", "Fast GenSpeed RPM Count", "Fast PitchAccuV1", "Fast PitchAccuV1 StdDev", "Fast PitchAccuV1 Min", "Fast PitchAccuV1 Max", "Fast PitchAccuV1 Count", "Fast PitchAngle", "Fast PitchAngle StdDev", "Fast PitchAngle Min", "Fast PitchAngle Max", "Fast PitchAngle Count", "Fast PitchAngle3", "Fast PitchAngle3 StdDev", "Fast PitchAngle3 Min", "Fast PitchAngle3 Max", "Fast PitchAngle3 Count", "Fast PitchAngle2", "Fast PitchAngle2 StdDev", "Fast PitchAngle2 Min", "Fast PitchAngle2 Max", "Fast PitchAngle2 Count", "Fast PitchConvCurrent1", "Fast PitchConvCurrent1 StdDev", "Fast PitchConvCurrent1 Min", "Fast PitchConvCurrent1 Max", "Fast PitchConvCurrent1 Count", "Fast PitchConvCurrent3", "Fast PitchConvCurrent3 StdDev", "Fast PitchConvCurrent3 Min", "Fast PitchConvCurrent3 Max", "Fast PitchConvCurrent3 Count", "Fast PitchConvCurrent2", "Fast PitchConvCurrent2 StdDev", "Fast PitchConvCurrent2 Min", "Fast PitchConvCurrent2 Max", "Fast PitchConvCurrent2 Count", "Fast PowerFactor", "Fast PowerFactor StdDev", "Fast PowerFactor Min", "Fast PowerFactor Max", "Fast PowerFactor Count", "Fast ReactivePowerSetpointPPC kVA", "Fast ReactivePowerSetpointPPC kVA StdDev", "Fast ReactivePowerSetpointPPC kVA Min", "Fast ReactivePowerSetpointPPC kVA Max", "Fast ReactivePowerSetpointPPC kVA Count", "Fast ReactivePower kVAr", "Fast ReactivePower kVAr StdDev", "Fast ReactivePower kVAr Min", "Fast ReactivePower kVAr Max", "Fast ReactivePower kVAr Count", "Fast RotorSpeed RPM", "Fast RotorSpeed RPM StdDev", "Fast RotorSpeed RPM Min", "Fast RotorSpeed RPM Max", "Fast RotorSpeed RPM Count", "Fast VoltageL1", "Fast VoltageL1 StdDev", "Fast VoltageL1 Min", "Fast VoltageL1 Max", "Fast VoltageL1 Count", "Fast VoltageL2", "Fast VoltageL2 StdDev", "Fast VoltageL2 Min", "Fast VoltageL2 Max", "Fast VoltageL2 Count", "Slow CapableCapacitiveReactPwr kVAr", "Slow CapableCapacitiveReactPwr kVAr StdDev", "Slow CapableCapacitiveReactPwr kVAr Min", "Slow CapableCapacitiveReactPwr kVAr Max", "Slow CapableCapacitiveReactPwr kVAr Count", "Slow CapableInductiveReactPwr kVAr", "Slow CapableInductiveReactPwr kVAr StdDev", "Slow CapableInductiveReactPwr kVAr Min", "Slow CapableInductiveReactPwr kVAr Max", "Slow CapableInductiveReactPwr kVAr Count", "Slow DateTime Sec", "Slow DateTime Sec StdDev", "Slow DateTime Sec Min", "Slow DateTime Sec Max", "Slow DateTime Sec Count", "Fast PitchAngle1", "Fast PitchAngle1 StdDev", "Fast PitchAngle1 Min", "Fast PitchAngle1 Max", "Fast PitchAngle1 Count", "Fast VoltageL3", "Fast VoltageL3 StdDev", "Fast VoltageL3 Min", "Fast VoltageL3 Max", "Fast VoltageL3 Count", "Slow CapableCapacitivePwrFactor", "Slow CapableCapacitivePwrFactor StdDev", "Slow CapableCapacitivePwrFactor Min", "Slow CapableCapacitivePwrFactor Max", "Slow CapableCapacitivePwrFactor Count", "Fast Total Production kWh", "Fast Total Production kWh StdDev", "Fast Total Production kWh Min", "Fast Total Production kWh Max", "Fast Total Production kWh Count", "Fast Total Prod Day kWh", "Fast Total Prod Day kWh StdDev", "Fast Total Prod Day kWh Min", "Fast Total Prod Day kWh Max", "Fast Total Prod Day kWh Count", "Fast Total Prod Month kWh", "Fast Total Prod Month kWh StdDev", "Fast Total Prod Month kWh Min", "Fast Total Prod Month kWh Max", "Fast Total Prod Month kWh Count", "Fast ActivePowerOutPWCSell kW", "Fast ActivePowerOutPWCSell kW StdDev", "Fast ActivePowerOutPWCSell kW Min", "Fast ActivePowerOutPWCSell kW Max", "Fast ActivePowerOutPWCSell kW Count", "Fast Frequency Hz", "Fast Frequency Hz StdDev", "Fast Frequency Hz Min", "Fast Frequency Hz Max", "Fast Frequency Hz Count", "Slow TempG1L2", "Slow TempG1L2 StdDev", "Slow TempG1L2 Min", "Slow TempG1L2 Max", "Slow TempG1L2 Count", "Slow TempG1L3", "Slow TempG1L3 StdDev", "Slow TempG1L3 Min", "Slow TempG1L3 Max", "Slow TempG1L3 Count", "Slow TempGearBoxHSSDE", "Slow TempGearBoxHSSDE StdDev", "Slow TempGearBoxHSSDE Min", "Slow TempGearBoxHSSDE Max", "Slow TempGearBoxHSSDE Count", "Slow TempGearBoxIMSNDE", "Slow TempGearBoxIMSNDE StdDev", "Slow TempGearBoxIMSNDE Min", "Slow TempGearBoxIMSNDE Max", "Slow TempGearBoxIMSNDE Count", "Slow TempOutdoor", "Slow TempOutdoor StdDev", "Slow TempOutdoor Min", "Slow TempOutdoor Max", "Slow TempOutdoor Count", "Fast PitchAccuV3", "Fast PitchAccuV3 StdDev", "Fast PitchAccuV3 Min", "Fast PitchAccuV3 Max", "Fast PitchAccuV3 Count", "Slow TotalTurbineActiveHours", "Slow TotalTurbineActiveHours StdDev", "Slow TotalTurbineActiveHours Min", "Slow TotalTurbineActiveHours Max", "Slow TotalTurbineActiveHours Count", "Slow TotalTurbineOKHours", "Slow TotalTurbineOKHours StdDev", "Slow TotalTurbineOKHours Min", "Slow TotalTurbineOKHours Max", "Slow TotalTurbineOKHours Count", "Slow TotalTurbineTimeAllHours", "Slow TotalTurbineTimeAllHours StdDev", "Slow TotalTurbineTimeAllHours Min", "Slow TotalTurbineTimeAllHours Max", "Slow TotalTurbineTimeAllHours Count", "Slow TempG1L1", "Slow TempG1L1 StdDev", "Slow TempG1L1 Min", "Slow TempG1L1 Max", "Slow TempG1L1 Count", "Slow TempGearBoxOilSump", "Slow TempGearBoxOilSump StdDev", "Slow TempGearBoxOilSump Min", "Slow TempGearBoxOilSump Max", "Slow TempGearBoxOilSump Count", "Fast PitchAccuV2", "Fast PitchAccuV2 StdDev", "Fast PitchAccuV2 Min", "Fast PitchAccuV2 Max", "Fast PitchAccuV2 Count", "Slow TotalGridOkHours", "Slow TotalGridOkHours StdDev", "Slow TotalGridOkHours Min", "Slow TotalGridOkHours Max", "Slow TotalGridOkHours Count", "Slow TotalActPowerOut kWh", "Slow TotalActPowerOut kWh StdDev", "Slow TotalActPowerOut kWh Min", "Slow TotalActPowerOut kWh Max", "Slow TotalActPowerOut kWh Count", "Fast YawService", "Fast YawService StdDev", "Fast YawService Min", "Fast YawService Max", "Fast YawService Count", "Fast YawAngle", "Fast YawAngle StdDev", "Fast YawAngle Min", "Fast YawAngle Max", "Fast YawAngle Count", "Slow CapableInductivePwrFactor", "Slow CapableInductivePwrFactor StdDev", "Slow CapableInductivePwrFactor Min", "Slow CapableInductivePwrFactor Max", "Slow CapableInductivePwrFactor Count", "Slow TempGearBoxHSSNDE", "Slow TempGearBoxHSSNDE StdDev", "Slow TempGearBoxHSSNDE Min", "Slow TempGearBoxHSSNDE Max", "Slow TempGearBoxHSSNDE Count", "Slow TempHubBearing", "Slow TempHubBearing StdDev", "Slow TempHubBearing Min", "Slow TempHubBearing Max", "Slow TempHubBearing Count", "Slow TotalG1ActiveHours", "Slow TotalG1ActiveHours StdDev", "Slow TotalG1ActiveHours Min", "Slow TotalG1ActiveHours Max", "Slow TotalG1ActiveHours Count", "Slow TotalActPowerOutG1 kWh", "Slow TotalActPowerOutG1 kWh StdDev", "Slow TotalActPowerOutG1 kWh Min", "Slow TotalActPowerOutG1 kWh Max", "Slow TotalActPowerOutG1 kWh Count", "Slow TotalReactPowerInG1 kVArh", "Slow TotalReactPowerInG1 kVArh StdDev", "Slow TotalReactPowerInG1 kVArh Min", "Slow TotalReactPowerInG1 kVArh Max", "Slow TotalReactPowerInG1 kVArh Count", "Slow NacelleDrill", "Slow NacelleDrill StdDev", "Slow NacelleDrill Min", "Slow NacelleDrill Max", "Slow NacelleDrill Count", "Slow TempGearBoxIMSDE", "Slow TempGearBoxIMSDE StdDev", "Slow TempGearBoxIMSDE Min", "Slow TempGearBoxIMSDE Max", "Slow TempGearBoxIMSDE Count", "Fast Total Operating hrs", "Fast Total Operating hrs StdDev", "Fast Total Operating hrs Min", "Fast Total Operating hrs Max", "Fast Total Operating hrs Count", "Slow TempNacelle", "Slow TempNacelle StdDev", "Slow TempNacelle Min", "Slow TempNacelle Max", "Slow TempNacelle Count", "Fast Total Grid OK hrs", "Fast Total Grid OK hrs StdDev", "Fast Total Grid OK hrs Min", "Fast Total Grid OK hrs Max", "Fast Total Grid OK hrs Count", "Fast Total WTG OK hrs", "Fast Total WTG OK hrs StdDev", "Fast Total WTG OK hrs Min", "Fast Total WTG OK hrs Max", "Fast Total WTG OK hrs Count", "Slow TempCabinetTopBox", "Slow TempCabinetTopBox StdDev", "Slow TempCabinetTopBox Min", "Slow TempCabinetTopBox Max", "Slow TempCabinetTopBox Count", "Slow TempGeneratorBearingNDE", "Slow TempGeneratorBearingNDE StdDev", "Slow TempGeneratorBearingNDE Min", "Slow TempGeneratorBearingNDE Max", "Slow TempGeneratorBearingNDE Count", "Fast Total Access hrs", "Fast Total Access hrs StdDev", "Fast Total Access hrs Min", "Fast Total Access hrs Max", "Fast Total Access hrs Count", "Slow TempBottomPowerSection", "Slow TempBottomPowerSection StdDev", "Slow TempBottomPowerSection Min", "Slow TempBottomPowerSection Max", "Slow TempBottomPowerSection Count", "Slow TempGeneratorBearingDE", "Slow TempGeneratorBearingDE StdDev", "Slow TempGeneratorBearingDE Min", "Slow TempGeneratorBearingDE Max", "Slow TempGeneratorBearingDE Count", "Slow TotalReactPowerIn kVArh", "Slow TotalReactPowerIn kVArh StdDev", "Slow TotalReactPowerIn kVArh Min", "Slow TotalReactPowerIn kVArh Max", "Slow TotalReactPowerIn kVArh Count", "Slow TempBottomControlSection", "Slow TempBottomControlSection StdDev", "Slow TempBottomControlSection Min", "Slow TempBottomControlSection Max", "Slow TempBottomControlSection Count", "Slow TempConv1", "Slow TempConv1 StdDev", "Slow TempConv1 Min", "Slow TempConv1 Max", "Slow TempConv1 Count", "Fast ActivePowerRated kW", "Fast ActivePowerRated kW StdDev", "Fast ActivePowerRated kW Min", "Fast ActivePowerRated kW Max", "Fast ActivePowerRated kW Count", "Fast NodeIP", "Fast NodeIP StdDev", "Fast NodeIP Min", "Fast NodeIP Max", "Fast NodeIP Count", "Fast PitchSpeed1", "Fast PitchSpeed1 StdDev", "Fast PitchSpeed1 Min", "Fast PitchSpeed1 Max", "Fast PitchSpeed1 Count", "Slow CFCardSize", "Slow CFCardSize StdDev", "Slow CFCardSize Min", "Slow CFCardSize Max", "Slow CFCardSize Count", "Slow CPU Number", "Slow CPU Number StdDev", "Slow CPU Number Min", "Slow CPU Number Max", "Slow CPU Number Count", "Slow CFCardSpaceLeft", "Slow CFCardSpaceLeft StdDev", "Slow CFCardSpaceLeft Min", "Slow CFCardSpaceLeft Max", "Slow CFCardSpaceLeft Count", "Slow TempBottomCapSection", "Slow TempBottomCapSection StdDev", "Slow TempBottomCapSection Min", "Slow TempBottomCapSection Max", "Slow TempBottomCapSection Count", "Slow RatedPower", "Slow RatedPower StdDev", "Slow RatedPower Min", "Slow RatedPower Max", "Slow RatedPower Count", "Slow TempConv3", "Slow TempConv3 StdDev", "Slow TempConv3 Min", "Slow TempConv3 Max", "Slow TempConv3 Count", "Slow TempConv2", "Slow TempConv2 StdDev", "Slow TempConv2 Min", "Slow TempConv2 Max", "Slow TempConv2 Count", "Slow TotalActPowerIn kWh", "Slow TotalActPowerIn kWh StdDev", "Slow TotalActPowerIn kWh Min", "Slow TotalActPowerIn kWh Max", "Slow TotalActPowerIn kWh Count", "Slow TotalActPowerInG1 kWh", "Slow TotalActPowerInG1 kWh StdDev", "Slow TotalActPowerInG1 kWh Min", "Slow TotalActPowerInG1 kWh Max", "Slow TotalActPowerInG1 kWh Count", "Slow TotalActPowerInG2 kWh", "Slow TotalActPowerInG2 kWh StdDev", "Slow TotalActPowerInG2 kWh Min", "Slow TotalActPowerInG2 kWh Max", "Slow TotalActPowerInG2 kWh Count", "Slow TotalActPowerOutG2 kWh", "Slow TotalActPowerOutG2 kWh StdDev", "Slow TotalActPowerOutG2 kWh Min", "Slow TotalActPowerOutG2 kWh Max", "Slow TotalActPowerOutG2 kWh Count", "Slow TotalG2ActiveHours", "Slow TotalG2ActiveHours StdDev", "Slow TotalG2ActiveHours Min", "Slow TotalG2ActiveHours Max", "Slow TotalG2ActiveHours Count", "Slow TotalReactPowerInG2 kVArh", "Slow TotalReactPowerInG2 kVArh StdDev", "Slow TotalReactPowerInG2 kVArh Min", "Slow TotalReactPowerInG2 kVArh Max", "Slow TotalReactPowerInG2 kVArh Count", "Slow TotalReactPowerOut kVArh", "Slow TotalReactPowerOut kVArh StdDev", "Slow TotalReactPowerOut kVArh Min", "Slow TotalReactPowerOut kVArh Max", "Slow TotalReactPowerOut kVArh Count", "Slow UTCoffset int", "Slow UTCoffset int StdDev", "Slow UTCoffset int Min", "Slow UTCoffset int Max", "Slow UTCoffset int Count", "Time Stamp", "Turbine"}

	// _ashfd_field := []string{"Fast_ActivePower_kW", "Fast_ActivePower_kW_StdDev", "Fast_ActivePower_kW_Min", "Fast_ActivePower_kW_Max", "Fast_ActivePower_kW_Count", "Fast_WindSpeed_ms", "Fast_WindSpeed_ms_StdDev", "Fast_WindSpeed_ms_Min", "Fast_WindSpeed_ms_Max", "Fast_WindSpeed_ms_Count", "Slow_NacellePos", "Slow_NacellePos_StdDev", "Slow_NacellePos_Min", "Slow_NacellePos_Max", "Slow_NacellePos_Count", "Slow_WindDirection", "Slow_WindDirection_StdDev", "Slow_WindDirection_Min", "Slow_WindDirection_Max", "Slow_WindDirection_Count", "Fast_CurrentL3", "Fast_CurrentL3_StdDev", "Fast_CurrentL3_Min", "Fast_CurrentL3_Max", "Fast_CurrentL3_Count", "Fast_CurrentL1", "Fast_CurrentL1_StdDev", "Fast_CurrentL1_Min", "Fast_CurrentL1_Max", "Fast_CurrentL1_Count", "Fast_ActivePowerSetpoint_kW", "Fast_ActivePowerSetpoint_kW_StdDev", "Fast_ActivePowerSetpoint_kW_Min", "Fast_ActivePowerSetpoint_kW_Max", "Fast_ActivePowerSetpoint_kW_Count", "Fast_CurrentL2", "Fast_CurrentL2_StdDev", "Fast_CurrentL2_Min", "Fast_CurrentL2_Max", "Fast_CurrentL2_Count", "Fast_DrTrVibValue", "Fast_DrTrVibValue_StdDev", "Fast_DrTrVibValue_Min", "Fast_DrTrVibValue_Max", "Fast_DrTrVibValue_Count", "Fast_GenSpeed_RPM", "Fast_GenSpeed_RPM_StdDev", "Fast_GenSpeed_RPM_Min", "Fast_GenSpeed_RPM_Max", "Fast_GenSpeed_RPM_Count", "Fast_PitchAccuV1", "Fast_PitchAccuV1_StdDev", "Fast_PitchAccuV1_Min", "Fast_PitchAccuV1_Max", "Fast_PitchAccuV1_Count", "Fast_PitchAngle", "Fast_PitchAngle_StdDev", "Fast_PitchAngle_Min", "Fast_PitchAngle_Max", "Fast_PitchAngle_Count", "Fast_PitchAngle3", "Fast_PitchAngle3_StdDev", "Fast_PitchAngle3_Min", "Fast_PitchAngle3_Max", "Fast_PitchAngle3_Count", "Fast_PitchAngle2", "Fast_PitchAngle2_StdDev", "Fast_PitchAngle2_Min", "Fast_PitchAngle2_Max", "Fast_PitchAngle2_Count", "Fast_PitchConvCurrent1", "Fast_PitchConvCurrent1_StdDev", "Fast_PitchConvCurrent1_Min", "Fast_PitchConvCurrent1_Max", "Fast_PitchConvCurrent1_Count", "Fast_PitchConvCurrent3", "Fast_PitchConvCurrent3_StdDev", "Fast_PitchConvCurrent3_Min", "Fast_PitchConvCurrent3_Max", "Fast_PitchConvCurrent3_Count", "Fast_PitchConvCurrent2", "Fast_PitchConvCurrent2_StdDev", "Fast_PitchConvCurrent2_Min", "Fast_PitchConvCurrent2_Max", "Fast_PitchConvCurrent2_Count", "Fast_PowerFactor", "Fast_PowerFactor_StdDev", "Fast_PowerFactor_Min", "Fast_PowerFactor_Max", "Fast_PowerFactor_Count", "Fast_ReactivePowerSetpointPPC_kVA", "Fast_ReactivePowerSetpointPPC_kVA_StdDev", "Fast_ReactivePowerSetpointPPC_kVA_Min", "Fast_ReactivePowerSetpointPPC_kVA_Max", "Fast_ReactivePowerSetpointPPC_kVA_Count", "Fast_ReactivePower_kVAr", "Fast_ReactivePower_kVAr_StdDev", "Fast_ReactivePower_kVAr_Min", "Fast_ReactivePower_kVAr_Max", "Fast_ReactivePower_kVAr_Count", "Fast_RotorSpeed_RPM", "Fast_RotorSpeed_RPM_StdDev", "Fast_RotorSpeed_RPM_Min", "Fast_RotorSpeed_RPM_Max", "Fast_RotorSpeed_RPM_Count", "Fast_VoltageL1", "Fast_VoltageL1_StdDev", "Fast_VoltageL1_Min", "Fast_VoltageL1_Max", "Fast_VoltageL1_Count", "Fast_VoltageL2", "Fast_VoltageL2_StdDev", "Fast_VoltageL2_Min", "Fast_VoltageL2_Max", "Fast_VoltageL2_Count", "Slow_CapableCapacitiveReactPwr_kVAr", "Slow_CapableCapacitiveReactPwr_kVAr_StdDev", "Slow_CapableCapacitiveReactPwr_kVAr_Min", "Slow_CapableCapacitiveReactPwr_kVAr_Max", "Slow_CapableCapacitiveReactPwr_kVAr_Count", "Slow_CapableInductiveReactPwr_kVAr", "Slow_CapableInductiveReactPwr_kVAr_StdDev", "Slow_CapableInductiveReactPwr_kVAr_Min", "Slow_CapableInductiveReactPwr_kVAr_Max", "Slow_CapableInductiveReactPwr_kVAr_Count", "Slow_DateTime_Sec", "Slow_DateTime_Sec_StdDev", "Slow_DateTime_Sec_Min", "Slow_DateTime_Sec_Max", "Slow_DateTime_Sec_Count", "Fast_PitchAngle1", "Fast_PitchAngle1_StdDev", "Fast_PitchAngle1_Min", "Fast_PitchAngle1_Max", "Fast_PitchAngle1_Count", "Fast_VoltageL3", "Fast_VoltageL3_StdDev", "Fast_VoltageL3_Min", "Fast_VoltageL3_Max", "Fast_VoltageL3_Count", "Slow_CapableCapacitivePwrFactor", "Slow_CapableCapacitivePwrFactor_StdDev", "Slow_CapableCapacitivePwrFactor_Min", "Slow_CapableCapacitivePwrFactor_Max", "Slow_CapableCapacitivePwrFactor_Count", "Fast_Total_Production_kWh", "Fast_Total_Production_kWh_StdDev", "Fast_Total_Production_kWh_Min", "Fast_Total_Production_kWh_Max", "Fast_Total_Production_kWh_Count", "Fast_Total_Prod_Day_kWh", "Fast_Total_Prod_Day_kWh_StdDev", "Fast_Total_Prod_Day_kWh_Min", "Fast_Total_Prod_Day_kWh_Max", "Fast_Total_Prod_Day_kWh_Count", "Fast_Total_Prod_Month_kWh", "Fast_Total_Prod_Month_kWh_StdDev", "Fast_Total_Prod_Month_kWh_Min", "Fast_Total_Prod_Month_kWh_Max", "Fast_Total_Prod_Month_kWh_Count", "Fast_ActivePowerOutPWCSell_kW", "Fast_ActivePowerOutPWCSell_kW_StdDev", "Fast_ActivePowerOutPWCSell_kW_Min", "Fast_ActivePowerOutPWCSell_kW_Max", "Fast_ActivePowerOutPWCSell_kW_Count", "Fast_Frequency_Hz", "Fast_Frequency_Hz_StdDev", "Fast_Frequency_Hz_Min", "Fast_Frequency_Hz_Max", "Fast_Frequency_Hz_Count", "Slow_TempG1L2", "Slow_TempG1L2_StdDev", "Slow_TempG1L2_Min", "Slow_TempG1L2_Max", "Slow_TempG1L2_Count", "Slow_TempG1L3", "Slow_TempG1L3_StdDev", "Slow_TempG1L3_Min", "Slow_TempG1L3_Max", "Slow_TempG1L3_Count", "Slow_TempGearBoxHSSDE", "Slow_TempGearBoxHSSDE_StdDev", "Slow_TempGearBoxHSSDE_Min", "Slow_TempGearBoxHSSDE_Max", "Slow_TempGearBoxHSSDE_Count", "Slow_TempGearBoxIMSNDE", "Slow_TempGearBoxIMSNDE_StdDev", "Slow_TempGearBoxIMSNDE_Min", "Slow_TempGearBoxIMSNDE_Max", "Slow_TempGearBoxIMSNDE_Count", "Slow_TempOutdoor", "Slow_TempOutdoor_StdDev", "Slow_TempOutdoor_Min", "Slow_TempOutdoor_Max", "Slow_TempOutdoor_Count", "Fast_PitchAccuV3", "Fast_PitchAccuV3_StdDev", "Fast_PitchAccuV3_Min", "Fast_PitchAccuV3_Max", "Fast_PitchAccuV3_Count", "Slow_TotalTurbineActiveHours", "Slow_TotalTurbineActiveHours_StdDev", "Slow_TotalTurbineActiveHours_Min", "Slow_TotalTurbineActiveHours_Max", "Slow_TotalTurbineActiveHours_Count", "Slow_TotalTurbineOKHours", "Slow_TotalTurbineOKHours_StdDev", "Slow_TotalTurbineOKHours_Min", "Slow_TotalTurbineOKHours_Max", "Slow_TotalTurbineOKHours_Count", "Slow_TotalTurbineTimeAllHours", "Slow_TotalTurbineTimeAllHours_StdDev", "Slow_TotalTurbineTimeAllHours_Min", "Slow_TotalTurbineTimeAllHours_Max", "Slow_TotalTurbineTimeAllHours_Count", "Slow_TempG1L1", "Slow_TempG1L1_StdDev", "Slow_TempG1L1_Min", "Slow_TempG1L1_Max", "Slow_TempG1L1_Count", "Slow_TempGearBoxOilSump", "Slow_TempGearBoxOilSump_StdDev", "Slow_TempGearBoxOilSump_Min", "Slow_TempGearBoxOilSump_Max", "Slow_TempGearBoxOilSump_Count", "Fast_PitchAccuV2", "Fast_PitchAccuV2_StdDev", "Fast_PitchAccuV2_Min", "Fast_PitchAccuV2_Max", "Fast_PitchAccuV2_Count", "Slow_TotalGridOkHours", "Slow_TotalGridOkHours_StdDev", "Slow_TotalGridOkHours_Min", "Slow_TotalGridOkHours_Max", "Slow_TotalGridOkHours_Count", "Slow_TotalActPowerOut_kWh", "Slow_TotalActPowerOut_kWh_StdDev", "Slow_TotalActPowerOut_kWh_Min", "Slow_TotalActPowerOut_kWh_Max", "Slow_TotalActPowerOut_kWh_Count", "Fast_YawService", "Fast_YawService_StdDev", "Fast_YawService_Min", "Fast_YawService_Max", "Fast_YawService_Count", "Fast_YawAngle", "Fast_YawAngle_StdDev", "Fast_YawAngle_Min", "Fast_YawAngle_Max", "Fast_YawAngle_Count", "Slow_CapableInductivePwrFactor", "Slow_CapableInductivePwrFactor_StdDev", "Slow_CapableInductivePwrFactor_Min", "Slow_CapableInductivePwrFactor_Max", "Slow_CapableInductivePwrFactor_Count", "Slow_TempGearBoxHSSNDE", "Slow_TempGearBoxHSSNDE_StdDev", "Slow_TempGearBoxHSSNDE_Min", "Slow_TempGearBoxHSSNDE_Max", "Slow_TempGearBoxHSSNDE_Count", "Slow_TempHubBearing", "Slow_TempHubBearing_StdDev", "Slow_TempHubBearing_Min", "Slow_TempHubBearing_Max", "Slow_TempHubBearing_Count", "Slow_TotalG1ActiveHours", "Slow_TotalG1ActiveHours_StdDev", "Slow_TotalG1ActiveHours_Min", "Slow_TotalG1ActiveHours_Max", "Slow_TotalG1ActiveHours_Count", "Slow_TotalActPowerOutG1_kWh", "Slow_TotalActPowerOutG1_kWh_StdDev", "Slow_TotalActPowerOutG1_kWh_Min", "Slow_TotalActPowerOutG1_kWh_Max", "Slow_TotalActPowerOutG1_kWh_Count", "Slow_TotalReactPowerInG1_kVArh", "Slow_TotalReactPowerInG1_kVArh_StdDev", "Slow_TotalReactPowerInG1_kVArh_Min", "Slow_TotalReactPowerInG1_kVArh_Max", "Slow_TotalReactPowerInG1_kVArh_Count", "Slow_NacelleDrill", "Slow_NacelleDrill_StdDev", "Slow_NacelleDrill_Min", "Slow_NacelleDrill_Max", "Slow_NacelleDrill_Count", "Slow_TempGearBoxIMSDE", "Slow_TempGearBoxIMSDE_StdDev", "Slow_TempGearBoxIMSDE_Min", "Slow_TempGearBoxIMSDE_Max", "Slow_TempGearBoxIMSDE_Count", "Fast_Total_Operating_hrs", "Fast_Total_Operating_hrs_StdDev", "Fast_Total_Operating_hrs_Min", "Fast_Total_Operating_hrs_Max", "Fast_Total_Operating_hrs_Count", "Slow_TempNacelle", "Slow_TempNacelle_StdDev", "Slow_TempNacelle_Min", "Slow_TempNacelle_Max", "Slow_TempNacelle_Count", "Fast_Total_Grid_OK_hrs", "Fast_Total_Grid_OK_hrs_StdDev", "Fast_Total_Grid_OK_hrs_Min", "Fast_Total_Grid_OK_hrs_Max", "Fast_Total_Grid_OK_hrs_Count", "Fast_Total_WTG_OK_hrs", "Fast_Total_WTG_OK_hrs_StdDev", "Fast_Total_WTG_OK_hrs_Min", "Fast_Total_WTG_OK_hrs_Max", "Fast_Total_WTG_OK_hrs_Count", "Slow_TempCabinetTopBox", "Slow_TempCabinetTopBox_StdDev", "Slow_TempCabinetTopBox_Min", "Slow_TempCabinetTopBox_Max", "Slow_TempCabinetTopBox_Count", "Slow_TempGeneratorBearingNDE", "Slow_TempGeneratorBearingNDE_StdDev", "Slow_TempGeneratorBearingNDE_Min", "Slow_TempGeneratorBearingNDE_Max", "Slow_TempGeneratorBearingNDE_Count", "Fast_Total_Access_hrs", "Fast_Total_Access_hrs_StdDev", "Fast_Total_Access_hrs_Min", "Fast_Total_Access_hrs_Max", "Fast_Total_Access_hrs_Count", "Slow_TempBottomPowerSection", "Slow_TempBottomPowerSection_StdDev", "Slow_TempBottomPowerSection_Min", "Slow_TempBottomPowerSection_Max", "Slow_TempBottomPowerSection_Count", "Slow_TempGeneratorBearingDE", "Slow_TempGeneratorBearingDE_StdDev", "Slow_TempGeneratorBearingDE_Min", "Slow_TempGeneratorBearingDE_Max", "Slow_TempGeneratorBearingDE_Count", "Slow_TotalReactPowerIn_kVArh", "Slow_TotalReactPowerIn_kVArh_StdDev", "Slow_TotalReactPowerIn_kVArh_Min", "Slow_TotalReactPowerIn_kVArh_Max", "Slow_TotalReactPowerIn_kVArh_Count", "Slow_TempBottomControlSection", "Slow_TempBottomControlSection_StdDev", "Slow_TempBottomControlSection_Min", "Slow_TempBottomControlSection_Max", "Slow_TempBottomControlSection_Count", "Slow_TempConv1", "Slow_TempConv1_StdDev", "Slow_TempConv1_Min", "Slow_TempConv1_Max", "Slow_TempConv1_Count", "Fast_ActivePowerRated_kW", "Fast_ActivePowerRated_kW_StdDev", "Fast_ActivePowerRated_kW_Min", "Fast_ActivePowerRated_kW_Max", "Fast_ActivePowerRated_kW_Count", "Fast_NodeIP", "Fast_NodeIP_StdDev", "Fast_NodeIP_Min", "Fast_NodeIP_Max", "Fast_NodeIP_Count", "Fast_PitchSpeed1", "Fast_PitchSpeed1_StdDev", "Fast_PitchSpeed1_Min", "Fast_PitchSpeed1_Max", "Fast_PitchSpeed1_Count", "Slow_CFCardSize", "Slow_CFCardSize_StdDev", "Slow_CFCardSize_Min", "Slow_CFCardSize_Max", "Slow_CFCardSize_Count", "Slow_CPU_Number", "Slow_CPU_Number_StdDev", "Slow_CPU_Number_Min", "Slow_CPU_Number_Max", "Slow_CPU_Number_Count", "Slow_CFCardSpaceLeft", "Slow_CFCardSpaceLeft_StdDev", "Slow_CFCardSpaceLeft_Min", "Slow_CFCardSpaceLeft_Max", "Slow_CFCardSpaceLeft_Count", "Slow_TempBottomCapSection", "Slow_TempBottomCapSection_StdDev", "Slow_TempBottomCapSection_Min", "Slow_TempBottomCapSection_Max", "Slow_TempBottomCapSection_Count", "Slow_RatedPower", "Slow_RatedPower_StdDev", "Slow_RatedPower_Min", "Slow_RatedPower_Max", "Slow_RatedPower_Count", "Slow_TempConv3", "Slow_TempConv3_StdDev", "Slow_TempConv3_Min", "Slow_TempConv3_Max", "Slow_TempConv3_Count", "Slow_TempConv2", "Slow_TempConv2_StdDev", "Slow_TempConv2_Min", "Slow_TempConv2_Max", "Slow_TempConv2_Count", "Slow_TotalActPowerIn_kWh", "Slow_TotalActPowerIn_kWh_StdDev", "Slow_TotalActPowerIn_kWh_Min", "Slow_TotalActPowerIn_kWh_Max", "Slow_TotalActPowerIn_kWh_Count", "Slow_TotalActPowerInG1_kWh", "Slow_TotalActPowerInG1_kWh_StdDev", "Slow_TotalActPowerInG1_kWh_Min", "Slow_TotalActPowerInG1_kWh_Max", "Slow_TotalActPowerInG1_kWh_Count", "Slow_TotalActPowerInG2_kWh", "Slow_TotalActPowerInG2_kWh_StdDev", "Slow_TotalActPowerInG2_kWh_Min", "Slow_TotalActPowerInG2_kWh_Max", "Slow_TotalActPowerInG2_kWh_Count", "Slow_TotalActPowerOutG2_kWh", "Slow_TotalActPowerOutG2_kWh_StdDev", "Slow_TotalActPowerOutG2_kWh_Min", "Slow_TotalActPowerOutG2_kWh_Max", "Slow_TotalActPowerOutG2_kWh_Count", "Slow_TotalG2ActiveHours", "Slow_TotalG2ActiveHours_StdDev", "Slow_TotalG2ActiveHours_Min", "Slow_TotalG2ActiveHours_Max", "Slow_TotalG2ActiveHours_Count", "Slow_TotalReactPowerInG2_kVArh", "Slow_TotalReactPowerInG2_kVArh_StdDev", "Slow_TotalReactPowerInG2_kVArh_Min", "Slow_TotalReactPowerInG2_kVArh_Max", "Slow_TotalReactPowerInG2_kVArh_Count", "Slow_TotalReactPowerOut_kVArh", "Slow_TotalReactPowerOut_kVArh_StdDev", "Slow_TotalReactPowerOut_kVArh_Min", "Slow_TotalReactPowerOut_kVArh_Max", "Slow_TotalReactPowerOut_kVArh_Count", "Slow_UTCoffset_int", "Slow_UTCoffset_int_StdDev", "Slow_UTCoffset_int_Min", "Slow_UTCoffset_int_Max", "Slow_UTCoffset_int_Count", "timestamp", "turbine"}

	_ashfd_label, _ashfd_field := GetScadaHFDHeader()
	for i, str := range _ashfd_field {
		tkm := tk.M{}.
			Set("_id", strings.ToLower(str)).
			Set("label", _ashfd_label[i]).
			Set("source", "ScadaDataHFD")

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
			if val[lowerField] != nil {
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

	switch tipe {
	case "scadaoem":
		obj := ScadaDataOEM{}
		tablename = obj.TableName()
		reflectVal = reflect.Indirect(reflect.ValueOf(obj))
	case "scadahfd":
		obj := ScadaDataHFD{}
		tablename = obj.TableName()
		reflectVal = reflect.Indirect(reflect.ValueOf(obj))
	case "met":
		obj := MetTower{}
		tablename = obj.TableName()
		reflectVal = reflect.Indirect(reflect.ValueOf(obj))
	case "eventraw":
		obj := EventRaw{}
		tablename = obj.TableName()
		reflectVal = reflect.Indirect(reflect.ValueOf(obj))
	case "eventdown":
		obj := EventDown{}
		tablename = obj.TableName()
		reflectVal = reflect.Indirect(reflect.ValueOf(obj))
	case "eventdownhfd":
		obj := EventDownHFD{}
		tablename = obj.TableName()
		reflectVal = reflect.Indirect(reflect.ValueOf(obj))
	}

	query := DB().Connection.NewQuery().From(tablename).Skip(p.Skip).Take(p.Take)
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

	tmpResult := make([]tk.M, 0)
	e = csr.Fetch(&tmpResult, 0, false)

	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	queryC := DB().Connection.NewQuery().From(tablename).Where(dbox.And(filter...))
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
	totalDuration := 0.0
	if needTotalTurbine {
		aggrData := []tk.M{}
		queryAggr := DB().Connection.NewQuery().From(tablename)
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
		if e != nil {
			return helper.CreateResult(false, nil, e.Error())
		}
		defer caggr.Close()
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
			totalActivePower = m.getSummaryColumn(filter, "fast_activepower_kw", "sum", tablename)
			AvgWS = m.getSummaryColumn(filter, "fast_windspeed_ms", "avg", tablename)
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

	tablename := new(ScadaDataHFD).TableName()
	arrscadaoem := []string{"_id"}
	source := "ScadaDataHFD"
	timestamp := "timestamp"
	var val1 reflect.Value
	switch tipe {
	case "ScadaOEM":
		tablename = new(ScadaDataOEM).TableName()
		arrscadaoem = append(arrscadaoem, "timestamputc")
		source = "ScadaDataOEM"
		timestamp = "timestamputc"
		obj1 := ScadaDataOEM{}
		val1 = reflect.Indirect(reflect.ValueOf(obj1))
	case "ScadaHFD":
		obj1 := ScadaDataHFD{}
		val1 = reflect.Indirect(reflect.ValueOf(obj1))
	}

	istimestamp := false
	arrmettower := []string{}
	ids := ""
	if p.Custom.Has("ColumnList") {
		for _, _val := range p.Custom["ColumnList"].([]interface{}) {
			_tkm, _ := tk.ToM(_val)
			ids = strings.ToLower(_tkm.GetString("_id"))
			if _tkm.GetString("source") == source {
				arrscadaoem = append(arrscadaoem, ids)
				if ids == "timestamp" {
					istimestamp = true
				}
			} else if _tkm.GetString("source") == "MetTower" {
				arrmettower = append(arrmettower, ids)
			}
		}
	}

	query := DB().Connection.NewQuery().
		Select(arrscadaoem...).
		From(tablename).
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

	config := lh.ReadConfig()

	loc, err := time.LoadLocation(config["ReadTimeLoc"])
	if err != nil {
		tk.Printfn("Get time in %s found %s", config["ReadTimeLoc"], err.Error())
	}

	for i, val := range results {
		if val.Has("timestamputc") {
			strangeTime := val.Get("timestamputc", time.Time{}).(time.Time).UTC().In(loc)
			itime := time.Date(strangeTime.Year(), strangeTime.Month(), strangeTime.Day(),
				strangeTime.Hour(), strangeTime.Minute(), strangeTime.Second(), strangeTime.Nanosecond(), time.UTC)
			arrmettowercond = append(arrmettowercond, itime)
			val.Set("timestamputc", itime)
			results[i] = val
		}
		if istimestamp {
			itime := val.Get("timestamp", time.Time{}).(time.Time).UTC()
			if tipe == "ScadaHFD" {
				arrmettowercond = append(arrmettowercond, itime)
			}
			val.Set("timestamp", itime)
			results[i] = val
		}
	}

	tkmmet := tk.M{}
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
	}

	queryC := DB().Connection.NewQuery().From(tablename).Where(dbox.And(filter...))
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

	queryAggr := DB().Connection.NewQuery().From(tablename)
	switch tipe {
	case "ScadaOEM":
		queryAggr = queryAggr.Aggr(dbox.AggrSum, "$power", "TotalPower").
			Aggr(dbox.AggrSum, "$powerlost", "TotalPowerLost").
			Aggr(dbox.AggrSum, "$ai_intern_activpower", "TotalActivePower").
			Aggr(dbox.AggrSum, "$ai_intern_windspeed", "AvgWindSpeed").
			Aggr(dbox.AggrSum, "$energy", "TotalEnergy").
			Group("turbine").Where(dbox.And(filter...))
	case "ScadaHFD":
		queryAggr = queryAggr.Group("turbine").Where(dbox.And(filter...))
	}

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
	switch tipe {
	case "ScadaOEM":
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
	case "ScadaHFD":
		totalActivePower = m.getSummaryColumn(filter, "fast_activepower_kw", "sum", tablename)
		AvgWS = m.getSummaryColumn(filter, "fast_windspeed_ms", "avg", tablename)
	}

	allFieldRequested := arrscadaoem
	allFieldRequested = append(allFieldRequested, arrmettower...)
	allHeader := map[string]string{}
	header := map[string]string{}
	obj2 := MetTower{}
	val2 := reflect.Indirect(reflect.ValueOf(obj2))
	fieldName := ""
	for i := 0; i < val1.Type().NumField(); i++ {
		fieldName = strings.ToLower(val1.Type().Field(i).Name)
		allHeader[fieldName] = val1.Field(i).Type().Name()
	}
	for i := 0; i < val2.Type().NumField(); i++ {
		fieldName = strings.ToLower(val2.Type().Field(i).Name)
		allHeader[fieldName] = val2.Field(i).Type().Name()
	}

	for _, val := range allFieldRequested {
		header[val] = allHeader[val]
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
		Total:            ccount.Count(),
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

func (m *DataBrowserController) getSummaryColumn(filter []*dbox.Filter, column, aggr, tablename string) float64 {
	xFilter := []*dbox.Filter{}
	queryAggr := DB().Connection.NewQuery().From(tablename)
	tkm := []tk.M{}

	switch column {
	case "fast_windspeed_ms":
		xFilter = append(filter, dbox.Gte("fast_windspeed_ms", 0))
		xFilter = append(xFilter, dbox.Lte("fast_windspeed_ms", 25))
	case "fast_activepower_kw":
		xFilter = append(filter, dbox.Ne("fast_activepower_kw", -999999.0))
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
	if e != nil {
		return 0
	}
	defer caggr.Close()
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

	switch typeExcel {
	case "ScadaOem":
		header = []string{"TimeStamp", "Turbine", "AI intern R PidAngleOut", "AI intern ActivPower ", "AI intern I1 ", "AI intern I2", "AI intern I3", "AI intern NacelleDrill ", "AI intern NacellePos ", "AI intern PitchAkku V1 ", "AI intern PitchAkku V2 ", "AI intern PitchAkku V3 ", "AI intern PitchAngle1 ", "AI intern PitchAngle2 ", "AI intern PitchAngle3 ", "AI intern PitchConv Current1 ", "AI intern PitchConv Current2 ", "AI intern PitchConv Current3 ", "AI intern PitchAngleSP Diff1 ", "AI intern PitchAngleSP Diff2 ", "AI intern PitchAngleSP Diff3 ", "AI intern ReactivPower ", "AI intern RpmDiff ", "AI intern U1 ", "AI intern U2 ", "AI intern U3 ", "AI intern WindDirection ", "AI intern WindSpeed ", "AI Intern WindSpeedDif ", "AI speed RotFR ", "AI WindSpeed1 ", "AI WindSpeed2 ", "AI WindVane1 ", "AI WindVane2 ", "AI internCurrentAsym ", "Temp GearBox IMS NDE ", "AI intern WindVaneDiff ", "C intern SpeedGenerator ", "C intern SpeedRotor ", "AI intern Speed RPMDiff FR1 RotCNT ", "AI intern Frequency Grid ", "Temp GearBox HSS NDE ", "AI DrTrVibValue ", "AI intern InLastErrorConv1 ", "AI intern InLastErrorConv2 ", "AI intern InLastErrorConv3 ", "AI intern TempConv1 ", "AI intern TempConv2 ", "AI intern TempConv3 ", "AI intern PitchSpeed2", "Temp YawBrake 1 ", "Temp YawBrake 2 ", "Temp G1L1 ", "Temp G1L2 ", "Temp G1L3 ", "Temp YawBrake 4", "AI HydrSystemPressure ", "Temp BottomControlSection Low ", "Temp GearBox HSS DE ", "Temp GearOilSump ", "Temp GeneratorBearing DE ", "Temp GeneratorBearing NDE ", "Temp MainBearing ", "Temp GearBox IMS DE ", "Temp Nacelle ", "Temp Outdoor ", "AI TowerVibValueAxial ", "AI intern DiffGenSpeedSPToAct ", "Temp YawBrake 5", "AI intern SpeedGenerator Proximity ", "AI intern SpeedDiff Encoder Proximity ", "AI GearOilPressure ", "Temp CabinetTopBox Low ", "Temp CabinetTopBox ", "Temp BottomControlSection ", "Temp BottomPowerSection ", "Temp BottomPowerSection Low ", "AI intern Pitch1 Status High ", "AI intern Pitch2 Status High ", "AI intern Pitch3 Status High ", "AI intern InPosition1 ch3", "AI intern InPosition2 ch3", "AI intern InPosition3 ch3", "AI intern Temp Brake Blade1 ", "AI intern Temp Brake Blade2 ", "AI intern Temp Brake Blade3 ", "AI intern Temp PitchMotor Blade1 ", "AI intern Temp PitchMotor Blade2 ", "AI intern Temp PitchMotor Blade3 ", "AI intern Temp Hub Additional1 ", "AI intern Temp Hub Additional2 ", "AI intern Temp Hub Additional3 ", "AI intern Pitch1 Status Low ", "AI intern Pitch2 Status Low ", "AI intern Pitch3 Status Low ", "AI intern Battery VoltageBlade1 center ", "AI intern Battery VoltageBlade2 center ", "AI intern Battery VoltageBlade3 center ", "AI intern Battery ChargingCur Blade1 ", "AI intern Battery ChargingCur Blade2 ", "AI intern Battery ChargingCur Blade3 ", "AI intern Battery DischargingCur Blade1 ", "AI intern Battery DischargingCur Blade2 ", "AI intern Battery DischargingCur Blade3 ", "AI intern PitchMotor BrakeVoltage Blade1 ", "AI intern PitchMotor BrakeVoltage Blade2 ", "AI intern PitchMotor BrakeVoltage Blade3 ", "AI intern PitchMotor BrakeCurrent Blade1 ", "AI intern PitchMotor BrakeCurrent Blade2 ", "AI intern PitchMotor BrakeCurrent Blade3 ", "AI intern Temp HubBox Blade1 ", "AI intern Temp HubBox Blade2 ", "AI intern Temp HubBox Blade3 ", "AI intern Temp Pitch1 HeatSink ", "AI intern Temp Pitch2 HeatSink ", "AI intern Temp Pitch3 HeatSink ", "AI intern ErrorStackBlade1 ", "AI intern ErrorStackBlade2 ", "AI intern ErrorStackBlade3 ", "AI intern Temp BatteryBox Blade1 ", "AI intern Temp BatteryBox Blade2 ", "AI intern Temp BatteryBox Blade3 ", "AI intern DC LinkVoltage1 ", "AI intern DC LinkVoltage2 ", "AI intern DC LinkVoltage3 ", "Temp Yaw Motor1 ", "Temp Yaw Motor2 ", "Temp Yaw Motor3 ", "Temp Yaw Motor4 ", "AO DFIG Power Setpiont ", "AO DFIG Q Setpoint ", "AI DFIG Torque actual ", "AI DFIG SpeedGenerator Encoder ", "AI intern DFIG DC Link Voltage actual ", "AI intern DFIG MSC current ", "AI intern DFIG Main voltage ", "AI intern DFIG Main current ", "AI intern DFIG active power actual ", "AI intern DFIG reactive power actual ", "AI intern DFIG active power actual LSC ", "AI intern DFIG LSC current ", "AI intern DFIG Data log number ", "AI intern Damper OscMagnitude ", "AI intern Damper PassbandFullLoad ", "AI YawBrake TempRise1 ", "AI YawBrake TempRise2 ", "AI YawBrake TempRise3 ", "AI YawBrake TempRise4 ", "AI intern NacelleDrill at NorthPosSensor "}
		tablename = new(ScadaDataOEM).TableName()
	case "ScadaDataHFD":
		header = []string{"TimeStamp", "ProjectName", "Turbine", "Fast_ActivePower_kW", "Fast_ActivePower_kW_StdDev", "Fast_ActivePower_kW_Min", "Fast_ActivePower_kW_Max", "Fast_ActivePower_kW_Count", "Fast_WindSpeed_ms", "Fast_WindSpeed_ms_StdDev", "Fast_WindSpeed_ms_Min", "Fast_WindSpeed_ms_Max", "Fast_WindSpeed_ms_Count", "Slow_NacellePos", "Slow_NacellePos_StdDev", "Slow_NacellePos_Min", "Slow_NacellePos_Max", "Slow_NacellePos_Count", "Slow_WindDirection", "Slow_WindDirection_StdDev", "Slow_WindDirection_Min", "Slow_WindDirection_Max", "Slow_WindDirection_Count", "Fast_CurrentL3", "Fast_CurrentL3_StdDev", "Fast_CurrentL3_Min", "Fast_CurrentL3_Max", "Fast_CurrentL3_Count", "Fast_CurrentL1", "Fast_CurrentL1_StdDev", "Fast_CurrentL1_Min", "Fast_CurrentL1_Max", "Fast_CurrentL1_Count", "Fast_ActivePowerSetpoint_kW", "Fast_ActivePowerSetpoint_kW_StdDev", "Fast_ActivePowerSetpoint_kW_Min", "Fast_ActivePowerSetpoint_kW_Max", "Fast_ActivePowerSetpoint_kW_Count", "Fast_CurrentL2", "Fast_CurrentL2_StdDev", "Fast_CurrentL2_Min", "Fast_CurrentL2_Max", "Fast_CurrentL2_Count", "Fast_DrTrVibValue", "Fast_DrTrVibValue_StdDev", "Fast_DrTrVibValue_Min", "Fast_DrTrVibValue_Max", "Fast_DrTrVibValue_Count", "Fast_GenSpeed_RPM", "Fast_GenSpeed_RPM_StdDev", "Fast_GenSpeed_RPM_Min", "Fast_GenSpeed_RPM_Max", "Fast_GenSpeed_RPM_Count", "Fast_PitchAccuV1", "Fast_PitchAccuV1_StdDev", "Fast_PitchAccuV1_Min", "Fast_PitchAccuV1_Max", "Fast_PitchAccuV1_Count", "Fast_PitchAngle", "Fast_PitchAngle_StdDev", "Fast_PitchAngle_Min", "Fast_PitchAngle_Max", "Fast_PitchAngle_Count", "Fast_PitchAngle3", "Fast_PitchAngle3_StdDev", "Fast_PitchAngle3_Min", "Fast_PitchAngle3_Max", "Fast_PitchAngle3_Count", "Fast_PitchAngle2", "Fast_PitchAngle2_StdDev", "Fast_PitchAngle2_Min", "Fast_PitchAngle2_Max", "Fast_PitchAngle2_Count", "Fast_PitchConvCurrent1", "Fast_PitchConvCurrent1_StdDev", "Fast_PitchConvCurrent1_Min", "Fast_PitchConvCurrent1_Max", "Fast_PitchConvCurrent1_Count", "Fast_PitchConvCurrent3", "Fast_PitchConvCurrent3_StdDev", "Fast_PitchConvCurrent3_Min", "Fast_PitchConvCurrent3_Max", "Fast_PitchConvCurrent3_Count", "Fast_PitchConvCurrent2", "Fast_PitchConvCurrent2_StdDev", "Fast_PitchConvCurrent2_Min", "Fast_PitchConvCurrent2_Max", "Fast_PitchConvCurrent2_Count", "Fast_PowerFactor", "Fast_PowerFactor_StdDev", "Fast_PowerFactor_Min", "Fast_PowerFactor_Max", "Fast_PowerFactor_Count", "Fast_ReactivePowerSetpointPPC_kVA", "Fast_ReactivePowerSetpointPPC_kVA_StdDev", "Fast_ReactivePowerSetpointPPC_kVA_Min", "Fast_ReactivePowerSetpointPPC_kVA_Max", "Fast_ReactivePowerSetpointPPC_kVA_Count", "Fast_ReactivePower_kVAr", "Fast_ReactivePower_kVAr_StdDev", "Fast_ReactivePower_kVAr_Min", "Fast_ReactivePower_kVAr_Max", "Fast_ReactivePower_kVAr_Count", "Fast_RotorSpeed_RPM", "Fast_RotorSpeed_RPM_StdDev", "Fast_RotorSpeed_RPM_Min", "Fast_RotorSpeed_RPM_Max", "Fast_RotorSpeed_RPM_Count", "Fast_VoltageL1", "Fast_VoltageL1_StdDev", "Fast_VoltageL1_Min", "Fast_VoltageL1_Max", "Fast_VoltageL1_Count", "Fast_VoltageL2", "Fast_VoltageL2_StdDev", "Fast_VoltageL2_Min", "Fast_VoltageL2_Max", "Fast_VoltageL2_Count", "Slow_CapableCapacitiveReactPwr_kVAr", "Slow_CapableCapacitiveReactPwr_kVAr_StdDev", "Slow_CapableCapacitiveReactPwr_kVAr_Min", "Slow_CapableCapacitiveReactPwr_kVAr_Max", "Slow_CapableCapacitiveReactPwr_kVAr_Count", "Slow_CapableInductiveReactPwr_kVAr", "Slow_CapableInductiveReactPwr_kVAr_StdDev", "Slow_CapableInductiveReactPwr_kVAr_Min", "Slow_CapableInductiveReactPwr_kVAr_Max", "Slow_CapableInductiveReactPwr_kVAr_Count", "Slow_DateTime_Sec", "Slow_DateTime_Sec_StdDev", "Slow_DateTime_Sec_Min", "Slow_DateTime_Sec_Max", "Slow_DateTime_Sec_Count", "Fast_PitchAngle1", "Fast_PitchAngle1_StdDev", "Fast_PitchAngle1_Min", "Fast_PitchAngle1_Max", "Fast_PitchAngle1_Count", "Fast_VoltageL3", "Fast_VoltageL3_StdDev", "Fast_VoltageL3_Min", "Fast_VoltageL3_Max", "Fast_VoltageL3_Count", "Slow_CapableCapacitivePwrFactor", "Slow_CapableCapacitivePwrFactor_StdDev", "Slow_CapableCapacitivePwrFactor_Min", "Slow_CapableCapacitivePwrFactor_Max", "Slow_CapableCapacitivePwrFactor_Count", "Fast_Total_Production_kWh", "Fast_Total_Production_kWh_StdDev", "Fast_Total_Production_kWh_Min", "Fast_Total_Production_kWh_Max", "Fast_Total_Production_kWh_Count", "Fast_Total_Prod_Day_kWh", "Fast_Total_Prod_Day_kWh_StdDev", "Fast_Total_Prod_Day_kWh_Min", "Fast_Total_Prod_Day_kWh_Max", "Fast_Total_Prod_Day_kWh_Count", "Fast_Total_Prod_Month_kWh", "Fast_Total_Prod_Month_kWh_StdDev", "Fast_Total_Prod_Month_kWh_Min", "Fast_Total_Prod_Month_kWh_Max", "Fast_Total_Prod_Month_kWh_Count", "Fast_ActivePowerOutPWCSell_kW", "Fast_ActivePowerOutPWCSell_kW_StdDev", "Fast_ActivePowerOutPWCSell_kW_Min", "Fast_ActivePowerOutPWCSell_kW_Max", "Fast_ActivePowerOutPWCSell_kW_Count", "Fast_Frequency_Hz", "Fast_Frequency_Hz_StdDev", "Fast_Frequency_Hz_Min", "Fast_Frequency_Hz_Max", "Fast_Frequency_Hz_Count", "Slow_TempG1L2", "Slow_TempG1L2_StdDev", "Slow_TempG1L2_Min", "Slow_TempG1L2_Max", "Slow_TempG1L2_Count", "Slow_TempG1L3", "Slow_TempG1L3_StdDev", "Slow_TempG1L3_Min", "Slow_TempG1L3_Max", "Slow_TempG1L3_Count", "Slow_TempGearBoxHSSDE", "Slow_TempGearBoxHSSDE_StdDev", "Slow_TempGearBoxHSSDE_Min", "Slow_TempGearBoxHSSDE_Max", "Slow_TempGearBoxHSSDE_Count", "Slow_TempGearBoxIMSNDE", "Slow_TempGearBoxIMSNDE_StdDev", "Slow_TempGearBoxIMSNDE_Min", "Slow_TempGearBoxIMSNDE_Max", "Slow_TempGearBoxIMSNDE_Count", "Slow_TempOutdoor", "Slow_TempOutdoor_StdDev", "Slow_TempOutdoor_Min", "Slow_TempOutdoor_Max", "Slow_TempOutdoor_Count", "Fast_PitchAccuV3", "Fast_PitchAccuV3_StdDev", "Fast_PitchAccuV3_Min", "Fast_PitchAccuV3_Max", "Fast_PitchAccuV3_Count", "Slow_TotalTurbineActiveHours", "Slow_TotalTurbineActiveHours_StdDev", "Slow_TotalTurbineActiveHours_Min", "Slow_TotalTurbineActiveHours_Max", "Slow_TotalTurbineActiveHours_Count", "Slow_TotalTurbineOKHours", "Slow_TotalTurbineOKHours_StdDev", "Slow_TotalTurbineOKHours_Min", "Slow_TotalTurbineOKHours_Max", "Slow_TotalTurbineOKHours_Count", "Slow_TotalTurbineTimeAllHours", "Slow_TotalTurbineTimeAllHours_StdDev", "Slow_TotalTurbineTimeAllHours_Min", "Slow_TotalTurbineTimeAllHours_Max", "Slow_TotalTurbineTimeAllHours_Count", "Slow_TempG1L1", "Slow_TempG1L1_StdDev", "Slow_TempG1L1_Min", "Slow_TempG1L1_Max", "Slow_TempG1L1_Count", "Slow_TempGearBoxOilSump", "Slow_TempGearBoxOilSump_StdDev", "Slow_TempGearBoxOilSump_Min", "Slow_TempGearBoxOilSump_Max", "Slow_TempGearBoxOilSump_Count", "Fast_PitchAccuV2", "Fast_PitchAccuV2_StdDev", "Fast_PitchAccuV2_Min", "Fast_PitchAccuV2_Max", "Fast_PitchAccuV2_Count", "Slow_TotalGridOkHours", "Slow_TotalGridOkHours_StdDev", "Slow_TotalGridOkHours_Min", "Slow_TotalGridOkHours_Max", "Slow_TotalGridOkHours_Count", "Slow_TotalActPowerOut_kWh", "Slow_TotalActPowerOut_kWh_StdDev", "Slow_TotalActPowerOut_kWh_Min", "Slow_TotalActPowerOut_kWh_Max", "Slow_TotalActPowerOut_kWh_Count", "Fast_YawService", "Fast_YawService_StdDev", "Fast_YawService_Min", "Fast_YawService_Max", "Fast_YawService_Count", "Fast_YawAngle", "Fast_YawAngle_StdDev", "Fast_YawAngle_Min", "Fast_YawAngle_Max", "Fast_YawAngle_Count", "Slow_CapableInductivePwrFactor", "Slow_CapableInductivePwrFactor_StdDev", "Slow_CapableInductivePwrFactor_Min", "Slow_CapableInductivePwrFactor_Max", "Slow_CapableInductivePwrFactor_Count", "Slow_TempGearBoxHSSNDE", "Slow_TempGearBoxHSSNDE_StdDev", "Slow_TempGearBoxHSSNDE_Min", "Slow_TempGearBoxHSSNDE_Max", "Slow_TempGearBoxHSSNDE_Count", "Slow_TempHubBearing", "Slow_TempHubBearing_StdDev", "Slow_TempHubBearing_Min", "Slow_TempHubBearing_Max", "Slow_TempHubBearing_Count", "Slow_TotalG1ActiveHours", "Slow_TotalG1ActiveHours_StdDev", "Slow_TotalG1ActiveHours_Min", "Slow_TotalG1ActiveHours_Max", "Slow_TotalG1ActiveHours_Count", "Slow_TotalActPowerOutG1_kWh", "Slow_TotalActPowerOutG1_kWh_StdDev", "Slow_TotalActPowerOutG1_kWh_Min", "Slow_TotalActPowerOutG1_kWh_Max", "Slow_TotalActPowerOutG1_kWh_Count", "Slow_TotalReactPowerInG1_kVArh", "Slow_TotalReactPowerInG1_kVArh_StdDev", "Slow_TotalReactPowerInG1_kVArh_Min", "Slow_TotalReactPowerInG1_kVArh_Max", "Slow_TotalReactPowerInG1_kVArh_Count", "Slow_NacelleDrill", "Slow_NacelleDrill_StdDev", "Slow_NacelleDrill_Min", "Slow_NacelleDrill_Max", "Slow_NacelleDrill_Count", "Slow_TempGearBoxIMSDE", "Slow_TempGearBoxIMSDE_StdDev", "Slow_TempGearBoxIMSDE_Min", "Slow_TempGearBoxIMSDE_Max", "Slow_TempGearBoxIMSDE_Count", "Fast_Total_Operating_hrs", "Fast_Total_Operating_hrs_StdDev", "Fast_Total_Operating_hrs_Min", "Fast_Total_Operating_hrs_Max", "Fast_Total_Operating_hrs_Count", "Slow_TempNacelle", "Slow_TempNacelle_StdDev", "Slow_TempNacelle_Min", "Slow_TempNacelle_Max", "Slow_TempNacelle_Count", "Fast_Total_Grid_OK_hrs", "Fast_Total_Grid_OK_hrs_StdDev", "Fast_Total_Grid_OK_hrs_Min", "Fast_Total_Grid_OK_hrs_Max", "Fast_Total_Grid_OK_hrs_Count", "Fast_Total_WTG_OK_hrs", "Fast_Total_WTG_OK_hrs_StdDev", "Fast_Total_WTG_OK_hrs_Min", "Fast_Total_WTG_OK_hrs_Max", "Fast_Total_WTG_OK_hrs_Count", "Slow_TempCabinetTopBox", "Slow_TempCabinetTopBox_StdDev", "Slow_TempCabinetTopBox_Min", "Slow_TempCabinetTopBox_Max", "Slow_TempCabinetTopBox_Count", "Slow_TempGeneratorBearingNDE", "Slow_TempGeneratorBearingNDE_StdDev", "Slow_TempGeneratorBearingNDE_Min", "Slow_TempGeneratorBearingNDE_Max", "Slow_TempGeneratorBearingNDE_Count", "Fast_Total_Access_hrs", "Fast_Total_Access_hrs_StdDev", "Fast_Total_Access_hrs_Min", "Fast_Total_Access_hrs_Max", "Fast_Total_Access_hrs_Count", "Slow_TempBottomPowerSection", "Slow_TempBottomPowerSection_StdDev", "Slow_TempBottomPowerSection_Min", "Slow_TempBottomPowerSection_Max", "Slow_TempBottomPowerSection_Count", "Slow_TempGeneratorBearingDE", "Slow_TempGeneratorBearingDE_StdDev", "Slow_TempGeneratorBearingDE_Min", "Slow_TempGeneratorBearingDE_Max", "Slow_TempGeneratorBearingDE_Count", "Slow_TotalReactPowerIn_kVArh", "Slow_TotalReactPowerIn_kVArh_StdDev", "Slow_TotalReactPowerIn_kVArh_Min", "Slow_TotalReactPowerIn_kVArh_Max", "Slow_TotalReactPowerIn_kVArh_Count", "Slow_TempBottomControlSection", "Slow_TempBottomControlSection_StdDev", "Slow_TempBottomControlSection_Min", "Slow_TempBottomControlSection_Max", "Slow_TempBottomControlSection_Count", "Slow_TempConv1", "Slow_TempConv1_StdDev", "Slow_TempConv1_Min", "Slow_TempConv1_Max", "Slow_TempConv1_Count", "Fast_ActivePowerRated_kW", "Fast_ActivePowerRated_kW_StdDev", "Fast_ActivePowerRated_kW_Min", "Fast_ActivePowerRated_kW_Max", "Fast_ActivePowerRated_kW_Count", "Fast_NodeIP", "Fast_NodeIP_StdDev", "Fast_NodeIP_Min", "Fast_NodeIP_Max", "Fast_NodeIP_Count", "Fast_PitchSpeed1", "Fast_PitchSpeed1_StdDev", "Fast_PitchSpeed1_Min", "Fast_PitchSpeed1_Max", "Fast_PitchSpeed1_Count", "Slow_CFCardSize", "Slow_CFCardSize_StdDev", "Slow_CFCardSize_Min", "Slow_CFCardSize_Max", "Slow_CFCardSize_Count", "Slow_CPU_Number", "Slow_CPU_Number_StdDev", "Slow_CPU_Number_Min", "Slow_CPU_Number_Max", "Slow_CPU_Number_Count", "Slow_CFCardSpaceLeft", "Slow_CFCardSpaceLeft_StdDev", "Slow_CFCardSpaceLeft_Min", "Slow_CFCardSpaceLeft_Max", "Slow_CFCardSpaceLeft_Count", "Slow_TempBottomCapSection", "Slow_TempBottomCapSection_StdDev", "Slow_TempBottomCapSection_Min", "Slow_TempBottomCapSection_Max", "Slow_TempBottomCapSection_Count", "Slow_RatedPower", "Slow_RatedPower_StdDev", "Slow_RatedPower_Min", "Slow_RatedPower_Max", "Slow_RatedPower_Count", "Slow_TempConv3", "Slow_TempConv3_StdDev", "Slow_TempConv3_Min", "Slow_TempConv3_Max", "Slow_TempConv3_Count", "Slow_TempConv2", "Slow_TempConv2_StdDev", "Slow_TempConv2_Min", "Slow_TempConv2_Max", "Slow_TempConv2_Count", "Slow_TotalActPowerIn_kWh", "Slow_TotalActPowerIn_kWh_StdDev", "Slow_TotalActPowerIn_kWh_Min", "Slow_TotalActPowerIn_kWh_Max", "Slow_TotalActPowerIn_kWh_Count", "Slow_TotalActPowerInG1_kWh", "Slow_TotalActPowerInG1_kWh_StdDev", "Slow_TotalActPowerInG1_kWh_Min", "Slow_TotalActPowerInG1_kWh_Max", "Slow_TotalActPowerInG1_kWh_Count", "Slow_TotalActPowerInG2_kWh", "Slow_TotalActPowerInG2_kWh_StdDev", "Slow_TotalActPowerInG2_kWh_Min", "Slow_TotalActPowerInG2_kWh_Max", "Slow_TotalActPowerInG2_kWh_Count", "Slow_TotalActPowerOutG2_kWh", "Slow_TotalActPowerOutG2_kWh_StdDev", "Slow_TotalActPowerOutG2_kWh_Min", "Slow_TotalActPowerOutG2_kWh_Max", "Slow_TotalActPowerOutG2_kWh_Count", "Slow_TotalG2ActiveHours", "Slow_TotalG2ActiveHours_StdDev", "Slow_TotalG2ActiveHours_Min", "Slow_TotalG2ActiveHours_Max", "Slow_TotalG2ActiveHours_Count", "Slow_TotalReactPowerInG2_kVArh", "Slow_TotalReactPowerInG2_kVArh_StdDev", "Slow_TotalReactPowerInG2_kVArh_Min", "Slow_TotalReactPowerInG2_kVArh_Max", "Slow_TotalReactPowerInG2_kVArh_Count", "Slow_TotalReactPowerOut_kVArh", "Slow_TotalReactPowerOut_kVArh_StdDev", "Slow_TotalReactPowerOut_kVArh_Min", "Slow_TotalReactPowerOut_kVArh_Max", "Slow_TotalReactPowerOut_kVArh_Count", "Slow_UTCoffset_int", "Slow_UTCoffset_int_StdDev", "Slow_UTCoffset_int_Min", "Slow_UTCoffset_int_Max", "Slow_UTCoffset_int_Count"}
		tablename = new(ScadaDataHFD).TableName()
	case "DowntimeEvent":
		header = []string{"Turbine", "TimeStart", "TimeEnd", "Down Grid", "Down Environment", "Down Machine", "Alarm Description", "Duration", "Reduce Availability"}
		tablename = new(EventDown).TableName()
		separator = ""
	case "EventRaw":
		header = []string{"TimeStamp", "Project Name", "Turbine", "Event Type", "Alarm Description", "Turbine Status", "Brake Type", "Brake Program", "Alarm Id", "Alarm Toggle"}
		tablename = new(EventRaw).TableName()
		separator = ""
	case "MetTower":
		header = []string{"TimeStamp", "WindDirNo", "VHubWS90mAvg", "VHubWS90mMax", "VHubWS90mMin", "VHubWS90mStdDev", "VHubWS90mCount", "VRefWS88mAvg", "VRefWS88mMax", "VRefWS88mMin", "VRefWS88mStdDev", "VRefWS88mCount", "VTipWS42mAvg", "VTipWS42mMax", "VTipWS42mMin", "VTipWS42mStdDev", "VTipWS42mCount", "DHubWD88mAvg", "DHubWD88mMax", "DHubWD88mMin", "DHubWD88mStdDev", "DHubWD88mCount", "DRefWD86mAvg", "DRefWD86mMax", "DRefWD86mMin", "DRefWD86mStdDev", "DRefWD86mCount", "THubHHubHumid855mAvg", "THubHHubHumid855mMax", "THubHHubHumid855mMin", "THubHHubHumid855mStdDev", "THubHHubHumid855mCount", "TRefHRefHumid855mAvg", "TRefHRefHumid855mMax", "TRefHRefHumid855mMin", "TRefHRefHumid855mStdDev", "TRefHRefHumid855mCount", "THubHHubTemp855mAvg", "THubHHubTemp855mMax", "THubHHubTemp855mMin", "THubHHubTemp855mStdDev", "THubHHubTemp855mCount", "TRefHRefTemp855mAvg", "TRefHRefTemp855mMax", "TRefHRefTemp855mMin", "TRefHRefTemp855mStdDev", "TRefHRefTemp855mCount", "BaroAirPress855mAvg", "BaroAirPress855mMax", "BaroAirPress855mMin", "BaroAirPress855mStdDev", "BaroAirPress855mCount", "YawAngleVoltageAvg", "YawAngleVoltageMax", "YawAngleVoltageMin", "YawAngleVoltageStdDev", "YawAngleVoltageCount", "OtherSensorVoltageAI1Avg", "OtherSensorVoltageAI1Max", "OtherSensorVoltageAI1Min", "OtherSensorVoltageAI1StdDev", "OtherSensorVoltageAI1Count", "OtherSensorVoltageAI2Avg", "OtherSensorVoltageAI2Max", "OtherSensorVoltageAI2Min", "OtherSensorVoltageAI2StdDev", "OtherSensorVoltageAI2Count", "OtherSensorVoltageAI3Avg", "OtherSensorVoltageAI3Max", "OtherSensorVoltageAI3Min", "OtherSensorVoltageAI3StdDev", "OtherSensorVoltageAI3Count", "OtherSensorVoltageAI4Avg", "OtherSensorVoltageAI4Max", "OtherSensorVoltageAI4Min", "OtherSensorVoltageAI4StdDev", "OtherSensorVoltageAI4Count", "GenRPMCurrentAvg", "GenRPMCurrentMax", "GenRPMCurrentMin", "GenRPMCurrentStdDev", "GenRPMCurrentCount", "WS_SCSCurrentAvg", "WS_SCSCurrentMax", "WS_SCSCurrentMin", "WS_SCSCurrentStdDev", "WS_SCSCurrentCount", "RainStatusCount", "RainStatusSum", "OtherSensor2StatusIO1Avg", "OtherSensor2StatusIO1Max", "OtherSensor2StatusIO1Min", "OtherSensor2StatusIO1StdDev", "OtherSensor2StatusIO1Count", "OtherSensor2StatusIO2Avg", "OtherSensor2StatusIO2Max", "OtherSensor2StatusIO2Min", "OtherSensor2StatusIO2StdDev", "OtherSensor2StatusIO2Count", "OtherSensor2StatusIO3Avg", "OtherSensor2StatusIO3Max", "OtherSensor2StatusIO3Min", "OtherSensor2StatusIO3StdDev", "OtherSensor2StatusIO3Count", "OtherSensor2StatusIO4Avg", "OtherSensor2StatusIO4Max", "OtherSensor2StatusIO4Min", "OtherSensor2StatusIO4StdDev", "OtherSensor2StatusIO4Count", "OtherSensor2StatusIO5Avg", "OtherSensor2StatusIO5Max", "OtherSensor2StatusIO5Min", "OtherSensor2StatusIO5StdDev", "OtherSensor2StatusIO5Count", "A1Avg", "A1Max", "A1Min", "A1StdDev", "A1Count", "A2Avg", "A2Max", "A2Min", "A2StdDev", "A2Count", "A3Avg", "A3Max", "A3Min", "A3StdDev", "A3Count", "A4Avg", "A4Max", "A4Min", "A4StdDev", "A4Count", "A5Avg", "A5Max", "A5Min", "A5StdDev", "A5Count", "A6Avg", "A6Max", "A6Min", "A6StdDev", "A6Count", "A7Avg", "A7Max", "A7Min", "A7StdDev", "A7Count", "A8Avg", "A8Max", "A8Min", "A8StdDev", "A8Count", "A9Avg", "A9Max", "A9Min", "A9StdDev", "A9Count", "A10Avg", "A10Max", "A10Min", "A10StdDev", "A10Count", "AC1Avg", "AC1Max", "AC1Min", "AC1StdDev", "AC1Count", "AC2Avg", "AC2Max", "AC2Min", "AC2StdDev", "AC2Count", "C1Avg", "C1Max", "C1Min", "C1StdDev", "C1Count", "C2Avg", "C2Max", "C2Min", "C2StdDev", "C2Count", "C3Avg", "C3Max", "C3Min", "C3StdDev", "C3Count", "D1Avg", "D1Max", "D1Min", "D1StdDev", "M1_1Avg", "M1_1Max", "M1_1Min", "M1_1StdDev", "M1_1Count", "M1_2Avg", "M1_2Max", "M1_2Min", "M1_2StdDev", "M1_2Count", "M1_3Avg", "M1_3Max", "M1_3Min", "M1_3StdDev", "M1_3Count", "M1_4Avg", "M1_4Max", "M1_4Min", "M1_4StdDev", "M1_4Count", "M1_5Avg", "M1_5Max", "M1_5Min", "M1_5StdDev", "M1_5Count", "M2_1Avg", "M2_1Max", "M2_1Min", "M2_1StdDev", "M2_1Count", "M2_2Avg", "M2_2Max", "M2_2Min", "M2_2StdDev", "M2_2Count", "M2_3Avg", "M2_3Max", "M2_3Min", "M2_3StdDev", "M2_3Count", "M2_4Avg", "M2_4Max", "M2_4Min", "M2_4StdDev", "M2_4Count", "M2_5Avg", "M2_5Max", "M2_5Min", "M2_5StdDev", "M2_5Count", "M2_6Avg", "M2_6Max", "M2_6Min", "M2_6StdDev", "M2_6Count", "M2_7Avg", "M2_7Max", "M2_7Min", "M2_7StdDev", "M2_7Count", "M2_8Avg", "M2_8Max", "M2_8Min", "M2_8StdDev", "M2_8Count", "VAvg", "VMax", "VMin", "IAvg", "IMax", "IMin", "T", "Addr", "WindDirDesc", "WSCategoryNo", "WSCategoryDesc"}
		tablename = new(MetTower).TableName()
		separator = ""
		p.Project = ""
	case "DowntimeEventHFD":
		header = []string{"TimeStart", "TimeEnd", "Down Grid", "Down Environment", "Down Machine", "Turbine", "Alarm Description", "Duration"}
		tablename = new(EventDownHFD).TableName()
		separator = ""
	}

	for _, val := range header {
		fieldList = append(fieldList, strings.ToLower(strings.Replace(strings.TrimSuffix(val, " "), " ", separator, -69)))
	}

	TimeCreate := time.Now().Format("2006-01-02_150405")
	CreateDateTime := typeExcel + TimeCreate

	if err := os.RemoveAll("web/assets/Excel/" + typeExcel + "/"); err != nil {
		tk.Println(err)
	}

	query := DB().Connection.NewQuery().From(tablename).Where(dbox.And(filter...))
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
	tablename := new(ScadaDataOEM).TableName()
	switch typeExcel {
	case "ScadaHFD":
		tablename = new(ScadaDataHFD).TableName()
	}

	query := DB().Connection.NewQuery().
		Select(arrscadaoem...).
		From(tablename).
		Where(dbox.And(filter...))

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

	config := lh.ReadConfig()
	loc, err := time.LoadLocation(config["ReadTimeLoc"])
	if err != nil {
		tk.Printfn("Get time in %s found %s", config["ReadTimeLoc"], err.Error())
	}

	for i, val := range results {
		if val.Has("timestamputc") {
			strangeTime := val.Get("timestamputc", time.Time{}).(time.Time).UTC().In(loc)
			itime := time.Date(strangeTime.Year(), strangeTime.Month(), strangeTime.Day(),
				strangeTime.Hour(), strangeTime.Minute(), strangeTime.Second(), strangeTime.Nanosecond(), time.UTC)
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
	floatString := ""
	dataType := ""

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
			if each[field] != nil {
				switch field {
				case "timestamp", "timestamputc", "timestart", "timeend":
					cell.Value = each[field].(time.Time).UTC().Format("2006-01-02 15:04:05")
				case "turbine":
					cell.Value = turbinename[each.GetString(field)]
				default:
					dataType = reflect.Indirect(reflect.ValueOf(each[field])).Type().Name()
					switch dataType {
					case "float64":
						floatString = tk.Sprintf("%.2f", each.GetFloat64(field))
						if each.GetFloat64(field) == -999999 {
							floatString = "-"
						}
						if floatString != "-" {
							floatString = FormatThousandSeparator(floatString)
						}
						cell.Value = floatString
					case "int":
						floatString = tk.Sprintf("%d", each.GetInt(field))
						if each.GetInt(field) == -999999 {
							floatString = "-"
						}
						cell.Value = floatString
					case "string":
						cell.Value = each.GetString(field)
					case "bool":
						cell.Value = strconv.FormatBool(each[field].(bool))
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
