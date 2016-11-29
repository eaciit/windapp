package dataconversion

import (
	. "eaciit/wfdemo-git-dev/library/helper"
	. "eaciit/wfdemo-git-dev/library/models"
	"log"
	"strconv"
	"sync"
	"time"

	hpp "eaciit/wfdemo-git-dev/processapp/helper"

	_ "github.com/eaciit/dbox/dbc/mongo"
	. "github.com/eaciit/orm"
	tk "github.com/eaciit/toolkit"
)

type DataConversion struct {
	Ctx *DataContext
}

var (
	emptyValueSmall = -0.000001
	emptyValueBig   = -9999999.0
	mutex           = &sync.Mutex{}
)

func NewDataConversion(ctx *DataContext) *DataConversion {
	dc := new(DataConversion)
	dc.Ctx = ctx

	return dc
}

func (d *DataConversion) Generate(file string) (errorLine tk.M) {
	log.Println("Start Conversion...")
	// funcName := "GenTenFromThreeSecond"
	var wg sync.WaitGroup
	_ = wg
	ctx := d.Ctx
	list := []tk.M{}
	pipes := []tk.M{}

	match := tk.M{}

	if file != "" {
		match = tk.M{"file": file}
		pipes = append(pipes, tk.M{"$match": match})
	}

	group := tk.M{
		"_id":           "$file",
		"min_timestamp": tk.M{"$min": "$timestampconverted"},
		"max_timestamp": tk.M{"$max": "$timestampconverted"},
	}

	pipes = append(pipes, tk.M{"$group": group})
	// pipes = append(pipes, tk.M{"$sort": tk.M{"_id": 1}})

	// tk.Printf("pipes: %#v \n", pipes)

	csr, e := ctx.Connection.NewQuery().
		From(new(ScadaThreeSecsExt).TableName()).
		Command("pipe", pipes).
		Cursor(nil)

	defer csr.Close()

	if e != nil {
		log.Printf("ERR: %#v \n", e.Error())
	} else {
		e = csr.Fetch(&list, 0, false)
		// log.Printf("list: %#v \n", list)
		if len(list) > 0 {
			for _, valList := range list {
				// log.Printf("valList: %#v \n", valList)
				startTime, _ := time.Parse("20060102 15:04", valList.Get("min_timestamp").(time.Time).UTC().Format("20060102 15:04"))
				endTime, _ := time.Parse("20060102 15:04", valList.Get("max_timestamp").(time.Time).UTC().Format("20060102 15:04"))

				for {
					if startTime.Format("2006-01-02 15:04") > endTime.Format("2006-01-02 15:04") {
						break
					}

					startTimeInt, _ := strconv.ParseInt(startTime.Format("200601021504"), 10, 64)

					fastActivePowerKWList := d.getStdDevAvgMinMaxCount(ctx, startTime, "fast_activepower_kw")
					fastActivePowerKWMap := d.getMap(fastActivePowerKWList, "fast_activepower_kw")

					// log.Printf("fastActivePowerKWMap: %#v \n\n", fastActivePowerKWMap)

					fastWindSpeedMsList := d.getStdDevAvgMinMaxCount(ctx, startTime, "fast_windspeed_ms")
					fastWindSpeedMsMap := d.getMap(fastWindSpeedMsList, "fast_windspeed_ms")

					slowNacellePosList := d.getStdDevAvgMinMaxCount(ctx, startTime, "slow_nacellepos")
					slowNacellePosMap := d.getMap(slowNacellePosList, "slow_nacellepos")

					slowWindDirectionList := d.getStdDevAvgMinMaxCount(ctx, startTime, "slow_winddirection")
					slowWindDirectionMap := d.getMap(slowWindDirectionList, "slow_winddirection")

					Fast_CurrentL3List := d.getStdDevAvgMinMaxCount(ctx, startTime, "fast_currentl3")
					Fast_CurrentL3Map := d.getMap(Fast_CurrentL3List, "fast_currentl3")

					Fast_CurrentL1List := d.getStdDevAvgMinMaxCount(ctx, startTime, "fast_currentl1")
					Fast_CurrentL1Map := d.getMap(Fast_CurrentL1List, "fast_currentl1")

					Fast_ActivePowerSetpoint_kWList := d.getStdDevAvgMinMaxCount(ctx, startTime, "fast_activepowersetpoint_kw")
					Fast_ActivePowerSetpoint_kWMap := d.getMap(Fast_ActivePowerSetpoint_kWList, "fast_activepowersetpoint_kw")

					Fast_CurrentL2List := d.getStdDevAvgMinMaxCount(ctx, startTime, "fast_currentl2")
					Fast_CurrentL2Map := d.getMap(Fast_CurrentL2List, "fast_currentl2")

					Fast_DrTrVibValueList := d.getStdDevAvgMinMaxCount(ctx, startTime, "fast_drtrvibvalue")
					Fast_DrTrVibValueMap := d.getMap(Fast_DrTrVibValueList, "fast_drtrvibvalue")

					Fast_GenSpeed_RPMList := d.getStdDevAvgMinMaxCount(ctx, startTime, "fast_genspeed_rpm")
					Fast_GenSpeed_RPMMap := d.getMap(Fast_GenSpeed_RPMList, "fast_genspeed_rpm")

					Fast_PitchAccuV1List := d.getStdDevAvgMinMaxCount(ctx, startTime, "fast_pitchaccuv1")
					Fast_PitchAccuV1Map := d.getMap(Fast_PitchAccuV1List, "fast_pitchaccuv1")

					Fast_PitchAngleList := d.getStdDevAvgMinMaxCount(ctx, startTime, "fast_pitchangle")
					Fast_PitchAngleMap := d.getMap(Fast_PitchAngleList, "fast_pitchangle")

					Fast_PitchAngle3List := d.getStdDevAvgMinMaxCount(ctx, startTime, "fast_pitchangle3")
					Fast_PitchAngle3Map := d.getMap(Fast_PitchAngle3List, "fast_pitchangle3")

					Fast_PitchAngle2List := d.getStdDevAvgMinMaxCount(ctx, startTime, "fast_pitchangle2")
					Fast_PitchAngle2Map := d.getMap(Fast_PitchAngle2List, "fast_pitchangle2")

					Fast_PitchConvCurrent1List := d.getStdDevAvgMinMaxCount(ctx, startTime, "fast_pitchconvcurrent1")
					Fast_PitchConvCurrent1Map := d.getMap(Fast_PitchConvCurrent1List, "fast_pitchconvcurrent1")

					Fast_PitchConvCurrent3List := d.getStdDevAvgMinMaxCount(ctx, startTime, "fast_pitchconvcurrent3")
					Fast_PitchConvCurrent3Map := d.getMap(Fast_PitchConvCurrent3List, "fast_pitchconvcurrent3")

					Fast_PitchConvCurrent2List := d.getStdDevAvgMinMaxCount(ctx, startTime, "fast_pitchconvcurrent2")
					Fast_PitchConvCurrent2Map := d.getMap(Fast_PitchConvCurrent2List, "fast_pitchconvcurrent2")

					Fast_PowerFactorList := d.getStdDevAvgMinMaxCount(ctx, startTime, "fast_powerfactor")
					Fast_PowerFactorMap := d.getMap(Fast_PowerFactorList, "fast_powerfactor")

					Fast_ReactivePowerSetpointPPC_kVAList := d.getStdDevAvgMinMaxCount(ctx, startTime, "fast_reactivepowersetpointppc_kva")
					Fast_ReactivePowerSetpointPPC_kVAMap := d.getMap(Fast_ReactivePowerSetpointPPC_kVAList, "fast_reactivepowersetpointppc_kva")

					Fast_ReactivePower_kVArList := d.getStdDevAvgMinMaxCount(ctx, startTime, "fast_reactivepower_kvar")
					Fast_ReactivePower_kVArMap := d.getMap(Fast_ReactivePower_kVArList, "fast_reactivepower_kvar")

					Fast_RotorSpeed_RPMList := d.getStdDevAvgMinMaxCount(ctx, startTime, "fast_rotorspeed_rpm")
					Fast_RotorSpeed_RPMMap := d.getMap(Fast_RotorSpeed_RPMList, "fast_rotorspeed_rpm")

					Fast_VoltageL1List := d.getStdDevAvgMinMaxCount(ctx, startTime, "fast_voltagel1")
					Fast_VoltageL1Map := d.getMap(Fast_VoltageL1List, "fast_voltagel1")

					Fast_VoltageL2List := d.getStdDevAvgMinMaxCount(ctx, startTime, "fast_voltagel2")
					Fast_VoltageL2Map := d.getMap(Fast_VoltageL2List, "fast_voltagel2")

					Slow_CapableCapacitiveReactPwr_kVArList := d.getStdDevAvgMinMaxCount(ctx, startTime, "slow_capablecapacitivereactpwr_kvar")
					Slow_CapableCapacitiveReactPwr_kVArMap := d.getMap(Slow_CapableCapacitiveReactPwr_kVArList, "slow_capablecapacitivereactpwr_kvar")

					Slow_CapableInductiveReactPwr_kVArList := d.getStdDevAvgMinMaxCount(ctx, startTime, "slow_capableinductivereactpwr_kvar")
					Slow_CapableInductiveReactPwr_kVArMap := d.getMap(Slow_CapableInductiveReactPwr_kVArList, "slow_capableinductivereactpwr_kvar")

					Slow_DateTime_SecList := d.getStdDevAvgMinMaxCount(ctx, startTime, "slow_datetime_sec")
					Slow_DateTime_SecMap := d.getMap(Slow_DateTime_SecList, "slow_datetime_sec")

					Fast_PitchAngle1List := d.getStdDevAvgMinMaxCount(ctx, startTime, "fast_pitchangle1")
					Fast_PitchAngle1Map := d.getMap(Fast_PitchAngle1List, "fast_pitchangle1")

					Fast_VoltageL3List := d.getStdDevAvgMinMaxCount(ctx, startTime, "fast_voltagel3")
					Fast_VoltageL3Map := d.getMap(Fast_VoltageL3List, "fast_voltagel3")

					Slow_CapableCapacitivePwrFactorList := d.getStdDevAvgMinMaxCount(ctx, startTime, "slow_capablecapacitivepwrfactor")
					Slow_CapableCapacitivePwrFactorMap := d.getMap(Slow_CapableCapacitivePwrFactorList, "slow_capablecapacitivepwrfactor")

					Fast_Total_Production_kWhList := d.getStdDevAvgMinMaxCount(ctx, startTime, "fast_total_production_kwh")
					Fast_Total_Production_kWhMap := d.getMap(Fast_Total_Production_kWhList, "fast_total_production_kwh")

					Fast_Total_Prod_Day_kWhList := d.getStdDevAvgMinMaxCount(ctx, startTime, "fast_total_prod_day_kwh")
					Fast_Total_Prod_Day_kWhMap := d.getMap(Fast_Total_Prod_Day_kWhList, "fast_total_prod_day_kwh")

					Fast_Total_Prod_Month_kWhList := d.getStdDevAvgMinMaxCount(ctx, startTime, "fast_total_prod_month_kwh")
					Fast_Total_Prod_Month_kWhMap := d.getMap(Fast_Total_Prod_Month_kWhList, "fast_total_prod_month_kwh")

					Fast_ActivePowerOutPWCSell_kWList := d.getStdDevAvgMinMaxCount(ctx, startTime, "fast_activepoweroutpwcsell_kw")
					Fast_ActivePowerOutPWCSell_kWMap := d.getMap(Fast_ActivePowerOutPWCSell_kWList, "fast_activepoweroutpwcsell_kw")

					Fast_Frequency_HzList := d.getStdDevAvgMinMaxCount(ctx, startTime, "fast_frequency_hz")
					Fast_Frequency_HzMap := d.getMap(Fast_Frequency_HzList, "fast_frequency_hz")

					Slow_TempG1L2List := d.getStdDevAvgMinMaxCount(ctx, startTime, "slow_tempg1l2")
					Slow_TempG1L2Map := d.getMap(Slow_TempG1L2List, "slow_tempg1l2")

					Slow_TempG1L3List := d.getStdDevAvgMinMaxCount(ctx, startTime, "slow_tempg1l3")
					Slow_TempG1L3Map := d.getMap(Slow_TempG1L3List, "slow_tempg1l3")

					Slow_TempGearBoxHSSDEList := d.getStdDevAvgMinMaxCount(ctx, startTime, "slow_tempgearboxhssde")
					Slow_TempGearBoxHSSDEMap := d.getMap(Slow_TempGearBoxHSSDEList, "slow_tempgearboxhssde")

					Slow_TempGearBoxIMSNDEList := d.getStdDevAvgMinMaxCount(ctx, startTime, "slow_tempgearboximsnde")
					Slow_TempGearBoxIMSNDEMap := d.getMap(Slow_TempGearBoxIMSNDEList, "slow_tempgearboximsnde")

					Slow_TempOutdoorList := d.getStdDevAvgMinMaxCount(ctx, startTime, "slow_tempoutdoor")
					Slow_TempOutdoorMap := d.getMap(Slow_TempOutdoorList, "slow_tempoutdoor")

					Fast_PitchAccuV3List := d.getStdDevAvgMinMaxCount(ctx, startTime, "fast_pitchaccuv3")
					Fast_PitchAccuV3Map := d.getMap(Fast_PitchAccuV3List, "fast_pitchaccuv3")

					Slow_TotalTurbineActiveHoursList := d.getStdDevAvgMinMaxCount(ctx, startTime, "slow_totalturbineactivehours")
					Slow_TotalTurbineActiveHoursMap := d.getMap(Slow_TotalTurbineActiveHoursList, "slow_totalturbineactivehours")

					Slow_TotalTurbineOKHoursList := d.getStdDevAvgMinMaxCount(ctx, startTime, "slow_totalturbineokhours")
					Slow_TotalTurbineOKHoursMap := d.getMap(Slow_TotalTurbineOKHoursList, "slow_totalturbineokhours")

					Slow_TotalTurbineTimeAllHoursList := d.getStdDevAvgMinMaxCount(ctx, startTime, "slow_totalturbinetimeallhours")
					Slow_TotalTurbineTimeAllHoursMap := d.getMap(Slow_TotalTurbineTimeAllHoursList, "slow_totalturbinetimeallhours")

					Slow_TempG1L1List := d.getStdDevAvgMinMaxCount(ctx, startTime, "slow_tempg1l1")
					Slow_TempG1L1Map := d.getMap(Slow_TempG1L1List, "slow_tempg1l1")

					Slow_TempGearBoxOilSumpList := d.getStdDevAvgMinMaxCount(ctx, startTime, "slow_tempgearboxoilsump")
					Slow_TempGearBoxOilSumpMap := d.getMap(Slow_TempGearBoxOilSumpList, "slow_tempgearboxoilsump")

					Fast_PitchAccuV2List := d.getStdDevAvgMinMaxCount(ctx, startTime, "fast_pitchaccuv2")
					Fast_PitchAccuV2Map := d.getMap(Fast_PitchAccuV2List, "fast_pitchaccuv2")

					Slow_TotalGridOkHoursList := d.getStdDevAvgMinMaxCount(ctx, startTime, "slow_totalgridokhours")
					Slow_TotalGridOkHoursMap := d.getMap(Slow_TotalGridOkHoursList, "slow_totalgridokhours")

					Slow_TotalActPowerOut_kWhList := d.getStdDevAvgMinMaxCount(ctx, startTime, "slow_totalactpowerout_kwh")
					Slow_TotalActPowerOut_kWhMap := d.getMap(Slow_TotalActPowerOut_kWhList, "slow_totalactpowerout_kwh")

					Fast_YawServiceList := d.getStdDevAvgMinMaxCount(ctx, startTime, "fast_yawservice")
					Fast_YawServiceMap := d.getMap(Fast_YawServiceList, "fast_yawservice")

					Fast_YawAngleList := d.getStdDevAvgMinMaxCount(ctx, startTime, "fast_yawangle")
					Fast_YawAngleMap := d.getMap(Fast_YawAngleList, "fast_yawangle")

					Slow_CapableInductivePwrFactorList := d.getStdDevAvgMinMaxCount(ctx, startTime, "slow_capableinductivepwrfactor")
					Slow_CapableInductivePwrFactorMap := d.getMap(Slow_CapableInductivePwrFactorList, "slow_capableinductivepwrfactor")

					Slow_TempGearBoxHSSNDEList := d.getStdDevAvgMinMaxCount(ctx, startTime, "slow_tempgearboxhssnde")
					Slow_TempGearBoxHSSNDEMap := d.getMap(Slow_TempGearBoxHSSNDEList, "slow_tempgearboxhssnde")

					Slow_TempHubBearingList := d.getStdDevAvgMinMaxCount(ctx, startTime, "slow_temphubbearing")
					Slow_TempHubBearingMap := d.getMap(Slow_TempHubBearingList, "slow_temphubbearing")

					Slow_TotalG1ActiveHoursList := d.getStdDevAvgMinMaxCount(ctx, startTime, "slow_totalg1activehours")
					Slow_TotalG1ActiveHoursMap := d.getMap(Slow_TotalG1ActiveHoursList, "slow_totalg1activehours")

					Slow_TotalActPowerOutG1_kWhList := d.getStdDevAvgMinMaxCount(ctx, startTime, "slow_totalactpoweroutg1_kwh")
					Slow_TotalActPowerOutG1_kWhMap := d.getMap(Slow_TotalActPowerOutG1_kWhList, "slow_totalactpoweroutg1_kwh")

					Slow_TotalReactPowerInG1_kVArhList := d.getStdDevAvgMinMaxCount(ctx, startTime, "slow_totalreactpowering1_kvarh")
					Slow_TotalReactPowerInG1_kVArhMap := d.getMap(Slow_TotalReactPowerInG1_kVArhList, "slow_totalreactpowering1_kvarh")

					Slow_NacelleDrillList := d.getStdDevAvgMinMaxCount(ctx, startTime, "slow_nacelledrill")
					Slow_NacelleDrillMap := d.getMap(Slow_NacelleDrillList, "slow_nacelledrill")

					Slow_TempGearBoxIMSDEList := d.getStdDevAvgMinMaxCount(ctx, startTime, "slow_tempgearboximsde")
					Slow_TempGearBoxIMSDEMap := d.getMap(Slow_TempGearBoxIMSDEList, "slow_tempgearboximsde")

					Fast_Total_Operating_hrsList := d.getStdDevAvgMinMaxCount(ctx, startTime, "fast_total_operating_hrs")
					Fast_Total_Operating_hrsMap := d.getMap(Fast_Total_Operating_hrsList, "fast_total_operating_hrs")

					Slow_TempNacelleList := d.getStdDevAvgMinMaxCount(ctx, startTime, "slow_tempnacelle")
					Slow_TempNacelleMap := d.getMap(Slow_TempNacelleList, "slow_tempnacelle")

					Fast_Total_Grid_OK_hrsList := d.getStdDevAvgMinMaxCount(ctx, startTime, "fast_total_grid_ok_hrs")
					Fast_Total_Grid_OK_hrsMap := d.getMap(Fast_Total_Grid_OK_hrsList, "fast_total_grid_ok_hrs")

					Fast_Total_WTG_OK_hrsList := d.getStdDevAvgMinMaxCount(ctx, startTime, "fast_total_wtg_ok_hrs")
					Fast_Total_WTG_OK_hrsMap := d.getMap(Fast_Total_WTG_OK_hrsList, "fast_total_wtg_ok_hrs")

					Slow_TempCabinetTopBoxList := d.getStdDevAvgMinMaxCount(ctx, startTime, "slow_tempcabinettopbox")
					Slow_TempCabinetTopBoxMap := d.getMap(Slow_TempCabinetTopBoxList, "slow_tempcabinettopbox")

					Slow_TempGeneratorBearingNDEList := d.getStdDevAvgMinMaxCount(ctx, startTime, "slow_tempgeneratorbearingnde")
					Slow_TempGeneratorBearingNDEMap := d.getMap(Slow_TempGeneratorBearingNDEList, "slow_tempgeneratorbearingnde")

					Fast_Total_Access_hrsList := d.getStdDevAvgMinMaxCount(ctx, startTime, "fast_total_access_hrs")
					Fast_Total_Access_hrsMap := d.getMap(Fast_Total_Access_hrsList, "fast_total_access_hrs")

					Slow_TempBottomPowerSectionList := d.getStdDevAvgMinMaxCount(ctx, startTime, "slow_tempbottompowersection")
					Slow_TempBottomPowerSectionMap := d.getMap(Slow_TempBottomPowerSectionList, "slow_tempbottompowersection")

					Slow_TempGeneratorBearingDEList := d.getStdDevAvgMinMaxCount(ctx, startTime, "slow_tempgeneratorbearingde")
					Slow_TempGeneratorBearingDEMap := d.getMap(Slow_TempGeneratorBearingDEList, "slow_tempgeneratorbearingde")

					Slow_TotalReactPowerIn_kVArhList := d.getStdDevAvgMinMaxCount(ctx, startTime, "slow_totalreactpowerin_kvarh")
					Slow_TotalReactPowerIn_kVArhMap := d.getMap(Slow_TotalReactPowerIn_kVArhList, "slow_totalreactpowerin_kvarh")

					Slow_TempBottomControlSectionList := d.getStdDevAvgMinMaxCount(ctx, startTime, "slow_tempbottomcontrolsection")
					Slow_TempBottomControlSectionMap := d.getMap(Slow_TempBottomControlSectionList, "slow_tempbottomcontrolsection")

					Slow_TempConv1List := d.getStdDevAvgMinMaxCount(ctx, startTime, "slow_tempconv1")
					Slow_TempConv1Map := d.getMap(Slow_TempConv1List, "slow_tempconv1")

					Fast_ActivePowerRated_kWList := d.getStdDevAvgMinMaxCount(ctx, startTime, "fast_activepowerrated_kw")
					Fast_ActivePowerRated_kWMap := d.getMap(Fast_ActivePowerRated_kWList, "fast_activepowerrated_kw")

					Fast_NodeIPList := d.getStdDevAvgMinMaxCount(ctx, startTime, "fast_nodeip")
					Fast_NodeIPMap := d.getMap(Fast_NodeIPList, "fast_nodeip")

					Fast_PitchSpeed1List := d.getStdDevAvgMinMaxCount(ctx, startTime, "fast_pitchspeed1")
					Fast_PitchSpeed1Map := d.getMap(Fast_PitchSpeed1List, "fast_pitchspeed1")

					Slow_CFCardSizeList := d.getStdDevAvgMinMaxCount(ctx, startTime, "slow_cfcardsize")
					Slow_CFCardSizeMap := d.getMap(Slow_CFCardSizeList, "slow_cfcardsize")

					Slow_CPU_NumberList := d.getStdDevAvgMinMaxCount(ctx, startTime, "slow_cpu_number")
					Slow_CPU_NumberMap := d.getMap(Slow_CPU_NumberList, "slow_cpu_number")

					Slow_CFCardSpaceLeftList := d.getStdDevAvgMinMaxCount(ctx, startTime, "slow_cfcardspaceleft")
					Slow_CFCardSpaceLeftMap := d.getMap(Slow_CFCardSpaceLeftList, "slow_cfcardspaceleft")

					Slow_TempBottomCapSectionList := d.getStdDevAvgMinMaxCount(ctx, startTime, "slow_tempbottomcapsection")
					Slow_TempBottomCapSectionMap := d.getMap(Slow_TempBottomCapSectionList, "slow_tempbottomcapsection")

					Slow_RatedPowerList := d.getStdDevAvgMinMaxCount(ctx, startTime, "slow_ratedpower")
					Slow_RatedPowerMap := d.getMap(Slow_RatedPowerList, "slow_ratedpower")

					Slow_TempConv3List := d.getStdDevAvgMinMaxCount(ctx, startTime, "slow_tempconv3")
					Slow_TempConv3Map := d.getMap(Slow_TempConv3List, "slow_tempconv3")

					Slow_TempConv2List := d.getStdDevAvgMinMaxCount(ctx, startTime, "slow_tempconv2")
					Slow_TempConv2Map := d.getMap(Slow_TempConv2List, "slow_tempconv2")

					Slow_TotalActPowerIn_kWhList := d.getStdDevAvgMinMaxCount(ctx, startTime, "slow_totalactpowerin_kwh")
					Slow_TotalActPowerIn_kWhMap := d.getMap(Slow_TotalActPowerIn_kWhList, "slow_totalactpowerin_kwh")

					Slow_TotalActPowerInG1_kWhList := d.getStdDevAvgMinMaxCount(ctx, startTime, "slow_totalactpowering1_kwh")
					Slow_TotalActPowerInG1_kWhMap := d.getMap(Slow_TotalActPowerInG1_kWhList, "slow_totalactpowering1_kwh")

					Slow_TotalActPowerInG2_kWhList := d.getStdDevAvgMinMaxCount(ctx, startTime, "slow_totalactpowering2_kwh")
					Slow_TotalActPowerInG2_kWhMap := d.getMap(Slow_TotalActPowerInG2_kWhList, "slow_totalactpowering2_kwh")

					Slow_TotalActPowerOutG2_kWhList := d.getStdDevAvgMinMaxCount(ctx, startTime, "slow_totalactpoweroutg2_kwh")
					Slow_TotalActPowerOutG2_kWhMap := d.getMap(Slow_TotalActPowerOutG2_kWhList, "slow_totalactpoweroutg2_kwh")

					Slow_TotalG2ActiveHoursList := d.getStdDevAvgMinMaxCount(ctx, startTime, "slow_totalg2activehours")
					Slow_TotalG2ActiveHoursMap := d.getMap(Slow_TotalG2ActiveHoursList, "slow_totalg2activehours")

					Slow_TotalReactPowerInG2_kVArhList := d.getStdDevAvgMinMaxCount(ctx, startTime, "slow_totalreactpowering2_kvarh")
					Slow_TotalReactPowerInG2_kVArhMap := d.getMap(Slow_TotalReactPowerInG2_kVArhList, "slow_totalreactpowering2_kvarh")

					Slow_TotalReactPowerOut_kVArhList := d.getStdDevAvgMinMaxCount(ctx, startTime, "slow_totalreactpowerout_kvarh")
					Slow_TotalReactPowerOut_kVArhMap := d.getMap(Slow_TotalReactPowerOut_kVArhList, "slow_totalreactpowerout_kvarh")

					Slow_UTCoffset_intList := d.getStdDevAvgMinMaxCount(ctx, startTime, "slow_utcoffset_int")
					Slow_UTCoffset_intMap := d.getMap(Slow_UTCoffset_intList, "slow_utcoffset_int")

					groupSub := tk.M{
						"_id": tk.M{
							// "timestamp":   "$timestampconverted",
							"timestampint": "$timestampconvertedint",
							"projectname":  "$projectname",
							"turbine":      "$turbine",
						},
						"count": tk.M{"$sum": 1},
					}

					pipesSub := []tk.M{}
					matchSub := tk.M{}

					if file != "" {
						matchSub.Set("file", file)
					}
					matchSub.Set("timestampconvertedint", startTimeInt)

					pipesSub = append(pipesSub, tk.M{"$match": matchSub})
					pipesSub = append(pipesSub, tk.M{"$group": groupSub})
					// pipesSub = append(pipesSub, tk.M{"$sort": tk.M{"_id": 1}})

					csr, e := ctx.Connection.NewQuery().
						From(new(ScadaThreeSecsExt).TableName()).
						Command("pipe", pipesSub).
						Cursor(nil)

					if e != nil {
						log.Printf("Error: %v \n", e.Error())
					}

					listSub := []tk.M{}
					e = csr.Fetch(&listSub, 0, false)

					countData := len(listSub)

					if countData > 0 {
						// log.Printf("timestamp: %v | %v | %v \n", startTime.Format("20060102 15:04"), countData, listSub[0].Get("_id").(tk.M).Get("timestampint"))

						countPerProcess := 1000
						counter := 0
						startIndex := counter * countPerProcess
						endIndex := (counter+1)*countPerProcess - 1
						isFinish := false

						for !isFinish {
							startIndex = counter * countPerProcess
							endIndex = (counter+1)*countPerProcess - 1

							if endIndex > countData {
								endIndex = countData
							}

							data := listSub[startIndex:endIndex]

							wg.Add(1)
							go func(data []tk.M) {
								for idx, val := range data {

									tenScada := new(ScadaConvTenMin)
									tenScada.No = idx + 1
									tenScada.File = file

									idSub := val.Get("_id").(tk.M)
									timeStampSubInt := idSub.Get("timestampint").(int64)
									timeStampSub, _ := time.Parse("200601021504", tk.ToString(timeStampSubInt))
									// timeStamp := startTime
									projectNameSub := idSub.GetString("projectname")
									turbineSub := idSub.GetString("turbine")

									timeStampStr := timeStampSub.UTC().Format("060102_1504")
									key := timeStampStr + "#" + projectNameSub + "#" + turbineSub

									tenScada.TimeStamp = timeStampSub
									tenScada.TimeStampInt, _ = strconv.ParseInt(tenScada.TimeStamp.Format("200601021504"), 10, 64)
									tenScada.DateInfo = GetDateInfo(timeStampSub)
									tenScada.ProjectName = projectNameSub
									tenScada.Turbine = turbineSub
									tenScada.Count = val.GetInt("count")

									tenScada = tenScada.New()

									fastActivePower := fastActivePowerKWMap[key]
									if fastActivePower != nil {
										tenScada.Fast_ActivePower_kW = fastActivePower.GetFloat64("fast_activepower_kw")
										tenScada.Fast_ActivePower_kW_StdDev = fastActivePower.GetFloat64("fast_activepower_kw_stddev")
										tenScada.Fast_ActivePower_kW_Min = fastActivePower.GetFloat64("fast_activepower_kw_min")
										tenScada.Fast_ActivePower_kW_Max = fastActivePower.GetFloat64("fast_activepower_kw_max")
										tenScada.Fast_ActivePower_kW_Count = fastActivePower.GetInt("fast_activepower_kw_count")
									} else {
										tenScada.Fast_ActivePower_kW = emptyValueBig
										tenScada.Fast_ActivePower_kW_StdDev = emptyValueBig
										tenScada.Fast_ActivePower_kW_Min = emptyValueBig
										tenScada.Fast_ActivePower_kW_Max = emptyValueBig

									}

									fastWindSpeedMs := fastWindSpeedMsMap[key]
									if fastWindSpeedMs != nil {
										tenScada.Fast_WindSpeed_ms = fastWindSpeedMs.GetFloat64("fast_windspeed_ms")
										tenScada.Fast_WindSpeed_ms_StdDev = fastWindSpeedMs.GetFloat64("fast_windspeed_ms_stddev")
										tenScada.Fast_WindSpeed_ms_Min = fastWindSpeedMs.GetFloat64("fast_windspeed_ms_min")
										tenScada.Fast_WindSpeed_ms_Max = fastWindSpeedMs.GetFloat64("fast_windspeed_ms_max")
										tenScada.Fast_WindSpeed_ms_Count = fastWindSpeedMs.GetInt("fast_windspeed_ms_count")
									} else {
										tenScada.Fast_WindSpeed_ms = emptyValueBig
										tenScada.Fast_WindSpeed_ms_StdDev = emptyValueBig
										tenScada.Fast_WindSpeed_ms_Min = emptyValueBig
										tenScada.Fast_WindSpeed_ms_Max = emptyValueBig

									}

									slowNacellePos := slowNacellePosMap[key]
									if slowNacellePos != nil {
										tenScada.Slow_NacellePos = slowNacellePos.GetFloat64("slow_nacellepos")
										tenScada.Slow_NacellePos_StdDev = slowNacellePos.GetFloat64("slow_nacellepos_stddev")
										tenScada.Slow_NacellePos_Min = slowNacellePos.GetFloat64("slow_nacellepos_min")
										tenScada.Slow_NacellePos_Max = slowNacellePos.GetFloat64("slow_nacellepos_max")
										tenScada.Slow_NacellePos_Count = slowNacellePos.GetInt("slow_nacellepos_count")
									} else {
										tenScada.Slow_NacellePos = emptyValueBig
										tenScada.Slow_NacellePos_StdDev = emptyValueBig
										tenScada.Slow_NacellePos_Min = emptyValueBig
										tenScada.Slow_NacellePos_Max = emptyValueBig

									}

									slowWindDirection := slowWindDirectionMap[key]
									if slowWindDirection != nil {
										tenScada.Slow_WindDirection = slowWindDirection.GetFloat64("slow_winddirection")
										tenScada.Slow_WindDirection_StdDev = slowWindDirection.GetFloat64("slow_winddirection_stddev")
										tenScada.Slow_WindDirection_Min = slowWindDirection.GetFloat64("slow_winddirection_min")
										tenScada.Slow_WindDirection_Max = slowWindDirection.GetFloat64("slow_winddirection_max")
										tenScada.Slow_WindDirection_Count = slowWindDirection.GetInt("slow_winddirection_count")
									} else {
										tenScada.Slow_WindDirection = emptyValueBig
										tenScada.Slow_WindDirection_StdDev = emptyValueBig
										tenScada.Slow_WindDirection_Min = emptyValueBig
										tenScada.Slow_WindDirection_Max = emptyValueBig

									}

									Fast_CurrentL3 := Fast_CurrentL3Map[key]
									if Fast_CurrentL3 != nil {
										tenScada.Fast_CurrentL3 = Fast_CurrentL3.GetFloat64("fast_currentl3")
										tenScada.Fast_CurrentL3_StdDev = Fast_CurrentL3.GetFloat64("fast_currentl3_stddev")
										tenScada.Fast_CurrentL3_Min = Fast_CurrentL3.GetFloat64("fast_currentl3_min")
										tenScada.Fast_CurrentL3_Max = Fast_CurrentL3.GetFloat64("fast_currentl3_max")
										tenScada.Fast_CurrentL3_Count = Fast_CurrentL3.GetInt("fast_currentl3_count")
									} else {
										tenScada.Fast_CurrentL3 = emptyValueBig
										tenScada.Fast_CurrentL3_StdDev = emptyValueBig
										tenScada.Fast_CurrentL3_Min = emptyValueBig
										tenScada.Fast_CurrentL3_Max = emptyValueBig

									}

									Fast_CurrentL1 := Fast_CurrentL1Map[key]
									if Fast_CurrentL1 != nil {
										tenScada.Fast_CurrentL1 = Fast_CurrentL1.GetFloat64("fast_currentl1")
										tenScada.Fast_CurrentL1_StdDev = Fast_CurrentL1.GetFloat64("fast_currentl1_stddev")
										tenScada.Fast_CurrentL1_Min = Fast_CurrentL1.GetFloat64("fast_currentl1_min")
										tenScada.Fast_CurrentL1_Max = Fast_CurrentL1.GetFloat64("fast_currentl1_max")
										tenScada.Fast_CurrentL1_Count = Fast_CurrentL1.GetInt("fast_currentl1_count")
									} else {
										tenScada.Fast_CurrentL1 = emptyValueBig
										tenScada.Fast_CurrentL1_StdDev = emptyValueBig
										tenScada.Fast_CurrentL1_Min = emptyValueBig
										tenScada.Fast_CurrentL1_Max = emptyValueBig

									}

									Fast_ActivePowerSetpoint_kW := Fast_ActivePowerSetpoint_kWMap[key]
									if Fast_ActivePowerSetpoint_kW != nil {
										tenScada.Fast_ActivePowerSetpoint_kW = Fast_ActivePowerSetpoint_kW.GetFloat64("fast_activepowersetpoint_kw")
										tenScada.Fast_ActivePowerSetpoint_kW_StdDev = Fast_ActivePowerSetpoint_kW.GetFloat64("fast_activepowersetpoint_kw_stddev")
										tenScada.Fast_ActivePowerSetpoint_kW_Min = Fast_ActivePowerSetpoint_kW.GetFloat64("fast_activepowersetpoint_kw_min")
										tenScada.Fast_ActivePowerSetpoint_kW_Max = Fast_ActivePowerSetpoint_kW.GetFloat64("fast_activepowersetpoint_kw_max")
										tenScada.Fast_ActivePowerSetpoint_kW_Count = Fast_ActivePowerSetpoint_kW.GetInt("fast_activepowersetpoint_kw_count")
									} else {
										tenScada.Fast_ActivePowerSetpoint_kW = emptyValueBig
										tenScada.Fast_ActivePowerSetpoint_kW_StdDev = emptyValueBig
										tenScada.Fast_ActivePowerSetpoint_kW_Min = emptyValueBig
										tenScada.Fast_ActivePowerSetpoint_kW_Max = emptyValueBig

									}

									Fast_CurrentL2 := Fast_CurrentL2Map[key]
									if Fast_CurrentL2 != nil {
										tenScada.Fast_CurrentL2 = Fast_CurrentL2.GetFloat64("fast_currentl2")
										tenScada.Fast_CurrentL2_StdDev = Fast_CurrentL2.GetFloat64("fast_currentl2_stddev")
										tenScada.Fast_CurrentL2_Min = Fast_CurrentL2.GetFloat64("fast_currentl2_min")
										tenScada.Fast_CurrentL2_Max = Fast_CurrentL2.GetFloat64("fast_currentl2_max")
										tenScada.Fast_CurrentL2_Count = Fast_CurrentL2.GetInt("fast_currentl2_count")
									} else {
										tenScada.Fast_CurrentL2 = emptyValueBig
										tenScada.Fast_CurrentL2_StdDev = emptyValueBig
										tenScada.Fast_CurrentL2_Min = emptyValueBig
										tenScada.Fast_CurrentL2_Max = emptyValueBig

									}

									Fast_DrTrVibValue := Fast_DrTrVibValueMap[key]
									if Fast_DrTrVibValue != nil {
										tenScada.Fast_DrTrVibValue = Fast_DrTrVibValue.GetFloat64("fast_drtrvibvalue")
										tenScada.Fast_DrTrVibValue_StdDev = Fast_DrTrVibValue.GetFloat64("fast_drtrvibvalue_stddev")
										tenScada.Fast_DrTrVibValue_Min = Fast_DrTrVibValue.GetFloat64("fast_drtrvibvalue_min")
										tenScada.Fast_DrTrVibValue_Max = Fast_DrTrVibValue.GetFloat64("fast_drtrvibvalue_max")
										tenScada.Fast_DrTrVibValue_Count = Fast_DrTrVibValue.GetInt("fast_drtrvibvalue_count")
									} else {
										tenScada.Fast_DrTrVibValue = emptyValueBig
										tenScada.Fast_DrTrVibValue_StdDev = emptyValueBig
										tenScada.Fast_DrTrVibValue_Min = emptyValueBig
										tenScada.Fast_DrTrVibValue_Max = emptyValueBig

									}

									Fast_GenSpeed_RPM := Fast_GenSpeed_RPMMap[key]
									if Fast_GenSpeed_RPM != nil {
										tenScada.Fast_GenSpeed_RPM = Fast_GenSpeed_RPM.GetFloat64("fast_genspeed_rpm")
										tenScada.Fast_GenSpeed_RPM_StdDev = Fast_GenSpeed_RPM.GetFloat64("fast_genspeed_rpm_stddev")
										tenScada.Fast_GenSpeed_RPM_Min = Fast_GenSpeed_RPM.GetFloat64("fast_genspeed_rpm_min")
										tenScada.Fast_GenSpeed_RPM_Max = Fast_GenSpeed_RPM.GetFloat64("fast_genspeed_rpm_max")
										tenScada.Fast_GenSpeed_RPM_Count = Fast_GenSpeed_RPM.GetInt("fast_genspeed_rpm_count")
									} else {
										tenScada.Fast_GenSpeed_RPM = emptyValueBig
										tenScada.Fast_GenSpeed_RPM_StdDev = emptyValueBig
										tenScada.Fast_GenSpeed_RPM_Min = emptyValueBig
										tenScada.Fast_GenSpeed_RPM_Max = emptyValueBig

									}

									Fast_PitchAccuV1 := Fast_PitchAccuV1Map[key]
									if Fast_PitchAccuV1 != nil {
										tenScada.Fast_PitchAccuV1 = Fast_PitchAccuV1.GetFloat64("fast_pitchaccuv1")
										tenScada.Fast_PitchAccuV1_StdDev = Fast_PitchAccuV1.GetFloat64("fast_pitchaccuv1_stddev")
										tenScada.Fast_PitchAccuV1_Min = Fast_PitchAccuV1.GetFloat64("fast_pitchaccuv1_min")
										tenScada.Fast_PitchAccuV1_Max = Fast_PitchAccuV1.GetFloat64("fast_pitchaccuv1_max")
										tenScada.Fast_PitchAccuV1_Count = Fast_PitchAccuV1.GetInt("fast_pitchaccuv1_count")
									} else {
										tenScada.Fast_PitchAccuV1 = emptyValueBig
										tenScada.Fast_PitchAccuV1_StdDev = emptyValueBig
										tenScada.Fast_PitchAccuV1_Min = emptyValueBig
										tenScada.Fast_PitchAccuV1_Max = emptyValueBig

									}

									Fast_PitchAngle := Fast_PitchAngleMap[key]
									if Fast_PitchAngle != nil {
										tenScada.Fast_PitchAngle = Fast_PitchAngle.GetFloat64("fast_pitchangle")
										tenScada.Fast_PitchAngle_StdDev = Fast_PitchAngle.GetFloat64("fast_pitchangle_stddev")
										tenScada.Fast_PitchAngle_Min = Fast_PitchAngle.GetFloat64("fast_pitchangle_min")
										tenScada.Fast_PitchAngle_Max = Fast_PitchAngle.GetFloat64("fast_pitchangle_max")
										tenScada.Fast_PitchAngle_Count = Fast_PitchAngle.GetInt("fast_pitchangle_count")
									} else {
										tenScada.Fast_PitchAngle = emptyValueBig
										tenScada.Fast_PitchAngle_StdDev = emptyValueBig
										tenScada.Fast_PitchAngle_Min = emptyValueBig
										tenScada.Fast_PitchAngle_Max = emptyValueBig

									}

									Fast_PitchAngle3 := Fast_PitchAngle3Map[key]
									if Fast_PitchAngle3 != nil {
										tenScada.Fast_PitchAngle3 = Fast_PitchAngle3.GetFloat64("fast_pitchangle3")
										tenScada.Fast_PitchAngle3_StdDev = Fast_PitchAngle3.GetFloat64("fast_pitchangle3_stddev")
										tenScada.Fast_PitchAngle3_Min = Fast_PitchAngle3.GetFloat64("fast_pitchangle3_min")
										tenScada.Fast_PitchAngle3_Max = Fast_PitchAngle3.GetFloat64("fast_pitchangle3_max")
										tenScada.Fast_PitchAngle3_Count = Fast_PitchAngle3.GetInt("fast_pitchangle3_count")
									} else {
										tenScada.Fast_PitchAngle3 = emptyValueBig
										tenScada.Fast_PitchAngle3_StdDev = emptyValueBig
										tenScada.Fast_PitchAngle3_Min = emptyValueBig
										tenScada.Fast_PitchAngle3_Max = emptyValueBig

									}

									Fast_PitchAngle2 := Fast_PitchAngle2Map[key]
									if Fast_PitchAngle2 != nil {
										tenScada.Fast_PitchAngle2 = Fast_PitchAngle2.GetFloat64("fast_pitchangle2")
										tenScada.Fast_PitchAngle2_StdDev = Fast_PitchAngle2.GetFloat64("fast_pitchangle2_stddev")
										tenScada.Fast_PitchAngle2_Min = Fast_PitchAngle2.GetFloat64("fast_pitchangle2_min")
										tenScada.Fast_PitchAngle2_Max = Fast_PitchAngle2.GetFloat64("fast_pitchangle2_max")
										tenScada.Fast_PitchAngle2_Count = Fast_PitchAngle2.GetInt("fast_pitchangle2_count")
									} else {
										tenScada.Fast_PitchAngle2 = emptyValueBig
										tenScada.Fast_PitchAngle2_StdDev = emptyValueBig
										tenScada.Fast_PitchAngle2_Min = emptyValueBig
										tenScada.Fast_PitchAngle2_Max = emptyValueBig

									}

									Fast_PitchConvCurrent1 := Fast_PitchConvCurrent1Map[key]
									if Fast_PitchConvCurrent1 != nil {
										tenScada.Fast_PitchConvCurrent1 = Fast_PitchConvCurrent1.GetFloat64("fast_pitchconvcurrent1")
										tenScada.Fast_PitchConvCurrent1_StdDev = Fast_PitchConvCurrent1.GetFloat64("fast_pitchconvcurrent1_stddev")
										tenScada.Fast_PitchConvCurrent1_Min = Fast_PitchConvCurrent1.GetFloat64("fast_pitchconvcurrent1_min")
										tenScada.Fast_PitchConvCurrent1_Max = Fast_PitchConvCurrent1.GetFloat64("fast_pitchconvcurrent1_max")
										tenScada.Fast_PitchConvCurrent1_Count = Fast_PitchConvCurrent1.GetInt("fast_pitchconvcurrent1_count")
									} else {
										tenScada.Fast_PitchConvCurrent1 = emptyValueBig
										tenScada.Fast_PitchConvCurrent1_StdDev = emptyValueBig
										tenScada.Fast_PitchConvCurrent1_Min = emptyValueBig
										tenScada.Fast_PitchConvCurrent1_Max = emptyValueBig

									}

									Fast_PitchConvCurrent3 := Fast_PitchConvCurrent3Map[key]
									if Fast_PitchConvCurrent3 != nil {
										tenScada.Fast_PitchConvCurrent3 = Fast_PitchConvCurrent3.GetFloat64("fast_pitchconvcurrent3")
										tenScada.Fast_PitchConvCurrent3_StdDev = Fast_PitchConvCurrent3.GetFloat64("fast_pitchconvcurrent3_stddev")
										tenScada.Fast_PitchConvCurrent3_Min = Fast_PitchConvCurrent3.GetFloat64("fast_pitchconvcurrent3_min")
										tenScada.Fast_PitchConvCurrent3_Max = Fast_PitchConvCurrent3.GetFloat64("fast_pitchconvcurrent3_max")
										tenScada.Fast_PitchConvCurrent3_Count = Fast_PitchConvCurrent3.GetInt("fast_pitchconvcurrent3_count")
									} else {
										tenScada.Fast_PitchConvCurrent3 = emptyValueBig
										tenScada.Fast_PitchConvCurrent3_StdDev = emptyValueBig
										tenScada.Fast_PitchConvCurrent3_Min = emptyValueBig
										tenScada.Fast_PitchConvCurrent3_Max = emptyValueBig

									}

									Fast_PitchConvCurrent2 := Fast_PitchConvCurrent2Map[key]
									if Fast_PitchConvCurrent2 != nil {
										tenScada.Fast_PitchConvCurrent2 = Fast_PitchConvCurrent2.GetFloat64("fast_pitchconvcurrent2")
										tenScada.Fast_PitchConvCurrent2_StdDev = Fast_PitchConvCurrent2.GetFloat64("fast_pitchconvcurrent2_stddev")
										tenScada.Fast_PitchConvCurrent2_Min = Fast_PitchConvCurrent2.GetFloat64("fast_pitchconvcurrent2_min")
										tenScada.Fast_PitchConvCurrent2_Max = Fast_PitchConvCurrent2.GetFloat64("fast_pitchconvcurrent2_max")
										tenScada.Fast_PitchConvCurrent2_Count = Fast_PitchConvCurrent2.GetInt("fast_pitchconvcurrent2_count")
									} else {
										tenScada.Fast_PitchConvCurrent2 = emptyValueBig
										tenScada.Fast_PitchConvCurrent2_StdDev = emptyValueBig
										tenScada.Fast_PitchConvCurrent2_Min = emptyValueBig
										tenScada.Fast_PitchConvCurrent2_Max = emptyValueBig

									}

									Fast_PowerFactor := Fast_PowerFactorMap[key]
									if Fast_PowerFactor != nil {
										tenScada.Fast_PowerFactor = Fast_PowerFactor.GetFloat64("fast_powerfactor")
										tenScada.Fast_PowerFactor_StdDev = Fast_PowerFactor.GetFloat64("fast_powerfactor_stddev")
										tenScada.Fast_PowerFactor_Min = Fast_PowerFactor.GetFloat64("fast_powerfactor_min")
										tenScada.Fast_PowerFactor_Max = Fast_PowerFactor.GetFloat64("fast_powerfactor_max")
										tenScada.Fast_PowerFactor_Count = Fast_PowerFactor.GetInt("fast_powerfactor_count")
									} else {
										tenScada.Fast_PowerFactor = emptyValueBig
										tenScada.Fast_PowerFactor_StdDev = emptyValueBig
										tenScada.Fast_PowerFactor_Min = emptyValueBig
										tenScada.Fast_PowerFactor_Max = emptyValueBig

									}

									Fast_ReactivePowerSetpointPPC_kVA := Fast_ReactivePowerSetpointPPC_kVAMap[key]
									if Fast_ReactivePowerSetpointPPC_kVA != nil {
										tenScada.Fast_ReactivePowerSetpointPPC_kVA = Fast_ReactivePowerSetpointPPC_kVA.GetFloat64("fast_reactivepowersetpointppc_kva")
										tenScada.Fast_ReactivePowerSetpointPPC_kVA_StdDev = Fast_ReactivePowerSetpointPPC_kVA.GetFloat64("fast_reactivepowersetpointppc_kva_stddev")
										tenScada.Fast_ReactivePowerSetpointPPC_kVA_Min = Fast_ReactivePowerSetpointPPC_kVA.GetFloat64("fast_reactivepowersetpointppc_kva_min")
										tenScada.Fast_ReactivePowerSetpointPPC_kVA_Max = Fast_ReactivePowerSetpointPPC_kVA.GetFloat64("fast_reactivepowersetpointppc_kva_max")
										tenScada.Fast_ReactivePowerSetpointPPC_kVA_Count = Fast_ReactivePowerSetpointPPC_kVA.GetInt("fast_reactivepowersetpointppc_kva_count")
									} else {
										tenScada.Fast_ReactivePowerSetpointPPC_kVA = emptyValueBig
										tenScada.Fast_ReactivePowerSetpointPPC_kVA_StdDev = emptyValueBig
										tenScada.Fast_ReactivePowerSetpointPPC_kVA_Min = emptyValueBig
										tenScada.Fast_ReactivePowerSetpointPPC_kVA_Max = emptyValueBig

									}

									Fast_ReactivePower_kVAr := Fast_ReactivePower_kVArMap[key]
									if Fast_ReactivePower_kVAr != nil {
										tenScada.Fast_ReactivePower_kVAr = Fast_ReactivePower_kVAr.GetFloat64("fast_reactivepower_kvar")
										tenScada.Fast_ReactivePower_kVAr_StdDev = Fast_ReactivePower_kVAr.GetFloat64("fast_reactivepower_kvar_stddev")
										tenScada.Fast_ReactivePower_kVAr_Min = Fast_ReactivePower_kVAr.GetFloat64("fast_reactivepower_kvar_min")
										tenScada.Fast_ReactivePower_kVAr_Max = Fast_ReactivePower_kVAr.GetFloat64("fast_reactivepower_kvar_max")
										tenScada.Fast_ReactivePower_kVAr_Count = Fast_ReactivePower_kVAr.GetInt("fast_reactivepower_kvar_count")
									} else {
										tenScada.Fast_ReactivePower_kVAr = emptyValueBig
										tenScada.Fast_ReactivePower_kVAr_StdDev = emptyValueBig
										tenScada.Fast_ReactivePower_kVAr_Min = emptyValueBig
										tenScada.Fast_ReactivePower_kVAr_Max = emptyValueBig

									}

									Fast_RotorSpeed_RPM := Fast_RotorSpeed_RPMMap[key]
									if Fast_RotorSpeed_RPM != nil {
										tenScada.Fast_RotorSpeed_RPM = Fast_RotorSpeed_RPM.GetFloat64("fast_rotorspeed_rpm")
										tenScada.Fast_RotorSpeed_RPM_StdDev = Fast_RotorSpeed_RPM.GetFloat64("fast_rotorspeed_rpm_stddev")
										tenScada.Fast_RotorSpeed_RPM_Min = Fast_RotorSpeed_RPM.GetFloat64("fast_rotorspeed_rpm_min")
										tenScada.Fast_RotorSpeed_RPM_Max = Fast_RotorSpeed_RPM.GetFloat64("fast_rotorspeed_rpm_max")
										tenScada.Fast_RotorSpeed_RPM_Count = Fast_RotorSpeed_RPM.GetInt("fast_rotorspeed_rpm_count")
									} else {
										tenScada.Fast_RotorSpeed_RPM = emptyValueBig
										tenScada.Fast_RotorSpeed_RPM_StdDev = emptyValueBig
										tenScada.Fast_RotorSpeed_RPM_Min = emptyValueBig
										tenScada.Fast_RotorSpeed_RPM_Max = emptyValueBig

									}

									Fast_VoltageL1 := Fast_VoltageL1Map[key]
									if Fast_VoltageL1 != nil {
										tenScada.Fast_VoltageL1 = Fast_VoltageL1.GetFloat64("fast_voltagel1")
										tenScada.Fast_VoltageL1_StdDev = Fast_VoltageL1.GetFloat64("fast_voltagel1_stddev")
										tenScada.Fast_VoltageL1_Min = Fast_VoltageL1.GetFloat64("fast_voltagel1_min")
										tenScada.Fast_VoltageL1_Max = Fast_VoltageL1.GetFloat64("fast_voltagel1_max")
										tenScada.Fast_VoltageL1_Count = Fast_VoltageL1.GetInt("fast_voltagel1_count")
									} else {
										tenScada.Fast_VoltageL1 = emptyValueBig
										tenScada.Fast_VoltageL1_StdDev = emptyValueBig
										tenScada.Fast_VoltageL1_Min = emptyValueBig
										tenScada.Fast_VoltageL1_Max = emptyValueBig

									}

									Fast_VoltageL2 := Fast_VoltageL2Map[key]
									if Fast_VoltageL2 != nil {
										tenScada.Fast_VoltageL2 = Fast_VoltageL2.GetFloat64("fast_voltagel2")
										tenScada.Fast_VoltageL2_StdDev = Fast_VoltageL2.GetFloat64("fast_voltagel2_stddev")
										tenScada.Fast_VoltageL2_Min = Fast_VoltageL2.GetFloat64("fast_voltagel2_min")
										tenScada.Fast_VoltageL2_Max = Fast_VoltageL2.GetFloat64("fast_voltagel2_max")
										tenScada.Fast_VoltageL2_Count = Fast_VoltageL2.GetInt("fast_voltagel2_count")
									} else {
										tenScada.Fast_VoltageL2 = emptyValueBig
										tenScada.Fast_VoltageL2_StdDev = emptyValueBig
										tenScada.Fast_VoltageL2_Min = emptyValueBig
										tenScada.Fast_VoltageL2_Max = emptyValueBig

									}

									Slow_CapableCapacitiveReactPwr_kVAr := Slow_CapableCapacitiveReactPwr_kVArMap[key]
									if Slow_CapableCapacitiveReactPwr_kVAr != nil {
										tenScada.Slow_CapableCapacitiveReactPwr_kVAr = Slow_CapableCapacitiveReactPwr_kVAr.GetFloat64("slow_capablecapacitivereactpwr_kvar")
										tenScada.Slow_CapableCapacitiveReactPwr_kVAr_StdDev = Slow_CapableCapacitiveReactPwr_kVAr.GetFloat64("slow_capablecapacitivereactpwr_kvar_stddev")
										tenScada.Slow_CapableCapacitiveReactPwr_kVAr_Min = Slow_CapableCapacitiveReactPwr_kVAr.GetFloat64("slow_capablecapacitivereactpwr_kvar_min")
										tenScada.Slow_CapableCapacitiveReactPwr_kVAr_Max = Slow_CapableCapacitiveReactPwr_kVAr.GetFloat64("slow_capablecapacitivereactpwr_kvar_max")
										tenScada.Slow_CapableCapacitiveReactPwr_kVAr_Count = Slow_CapableCapacitiveReactPwr_kVAr.GetInt("slow_capablecapacitivereactpwr_kvar_count")
									} else {
										tenScada.Slow_CapableCapacitiveReactPwr_kVAr = emptyValueBig
										tenScada.Slow_CapableCapacitiveReactPwr_kVAr_StdDev = emptyValueBig
										tenScada.Slow_CapableCapacitiveReactPwr_kVAr_Min = emptyValueBig
										tenScada.Slow_CapableCapacitiveReactPwr_kVAr_Max = emptyValueBig

									}

									Slow_CapableInductiveReactPwr_kVAr := Slow_CapableInductiveReactPwr_kVArMap[key]
									if Slow_CapableInductiveReactPwr_kVAr != nil {
										tenScada.Slow_CapableInductiveReactPwr_kVAr = Slow_CapableInductiveReactPwr_kVAr.GetFloat64("slow_capableinductivereactpwr_kvar")
										tenScada.Slow_CapableInductiveReactPwr_kVAr_StdDev = Slow_CapableInductiveReactPwr_kVAr.GetFloat64("slow_capableinductivereactpwr_kvar_stddev")
										tenScada.Slow_CapableInductiveReactPwr_kVAr_Min = Slow_CapableInductiveReactPwr_kVAr.GetFloat64("slow_capableinductivereactpwr_kvar_min")
										tenScada.Slow_CapableInductiveReactPwr_kVAr_Max = Slow_CapableInductiveReactPwr_kVAr.GetFloat64("slow_capableinductivereactpwr_kvar_max")
										tenScada.Slow_CapableInductiveReactPwr_kVAr_Count = Slow_CapableInductiveReactPwr_kVAr.GetInt("slow_capableinductivereactpwr_kvar_count")
									} else {
										tenScada.Slow_CapableInductiveReactPwr_kVAr = emptyValueBig
										tenScada.Slow_CapableInductiveReactPwr_kVAr_StdDev = emptyValueBig
										tenScada.Slow_CapableInductiveReactPwr_kVAr_Min = emptyValueBig
										tenScada.Slow_CapableInductiveReactPwr_kVAr_Max = emptyValueBig

									}

									Slow_DateTime_Sec := Slow_DateTime_SecMap[key]
									if Slow_DateTime_Sec != nil {
										tenScada.Slow_DateTime_Sec = Slow_DateTime_Sec.GetFloat64("slow_datetime_sec")
										tenScada.Slow_DateTime_Sec_StdDev = Slow_DateTime_Sec.GetFloat64("slow_datetime_sec_stddev")
										tenScada.Slow_DateTime_Sec_Min = Slow_DateTime_Sec.GetFloat64("slow_datetime_sec_min")
										tenScada.Slow_DateTime_Sec_Max = Slow_DateTime_Sec.GetFloat64("slow_datetime_sec_max")
										tenScada.Slow_DateTime_Sec_Count = Slow_DateTime_Sec.GetInt("slow_datetime_sec_count")
									} else {
										tenScada.Slow_DateTime_Sec = emptyValueBig
										tenScada.Slow_DateTime_Sec_StdDev = emptyValueBig
										tenScada.Slow_DateTime_Sec_Min = emptyValueBig
										tenScada.Slow_DateTime_Sec_Max = emptyValueBig

									}

									Fast_PitchAngle1 := Fast_PitchAngle1Map[key]
									if Fast_PitchAngle1 != nil {
										tenScada.Fast_PitchAngle1 = Fast_PitchAngle1.GetFloat64("fast_pitchangle1")
										tenScada.Fast_PitchAngle1_StdDev = Fast_PitchAngle1.GetFloat64("fast_pitchangle1_stddev")
										tenScada.Fast_PitchAngle1_Min = Fast_PitchAngle1.GetFloat64("fast_pitchangle1_min")
										tenScada.Fast_PitchAngle1_Max = Fast_PitchAngle1.GetFloat64("fast_pitchangle1_max")
										tenScada.Fast_PitchAngle1_Count = Fast_PitchAngle1.GetInt("fast_pitchangle1_count")
									} else {
										tenScada.Fast_PitchAngle1 = emptyValueBig
										tenScada.Fast_PitchAngle1_StdDev = emptyValueBig
										tenScada.Fast_PitchAngle1_Min = emptyValueBig
										tenScada.Fast_PitchAngle1_Max = emptyValueBig

									}

									Fast_VoltageL3 := Fast_VoltageL3Map[key]
									if Fast_VoltageL3 != nil {
										tenScada.Fast_VoltageL3 = Fast_VoltageL3.GetFloat64("fast_voltagel3")
										tenScada.Fast_VoltageL3_StdDev = Fast_VoltageL3.GetFloat64("fast_voltagel3_stddev")
										tenScada.Fast_VoltageL3_Min = Fast_VoltageL3.GetFloat64("fast_voltagel3_min")
										tenScada.Fast_VoltageL3_Max = Fast_VoltageL3.GetFloat64("fast_voltagel3_max")
										tenScada.Fast_VoltageL3_Count = Fast_VoltageL3.GetInt("fast_voltagel3_count")
									} else {
										tenScada.Fast_VoltageL3 = emptyValueBig
										tenScada.Fast_VoltageL3_StdDev = emptyValueBig
										tenScada.Fast_VoltageL3_Min = emptyValueBig
										tenScada.Fast_VoltageL3_Max = emptyValueBig

									}

									Slow_CapableCapacitivePwrFactor := Slow_CapableCapacitivePwrFactorMap[key]
									if Slow_CapableCapacitivePwrFactor != nil {
										tenScada.Slow_CapableCapacitivePwrFactor = Slow_CapableCapacitivePwrFactor.GetFloat64("slow_capablecapacitivepwrfactor")
										tenScada.Slow_CapableCapacitivePwrFactor_StdDev = Slow_CapableCapacitivePwrFactor.GetFloat64("slow_capablecapacitivepwrfactor_stddev")
										tenScada.Slow_CapableCapacitivePwrFactor_Min = Slow_CapableCapacitivePwrFactor.GetFloat64("slow_capablecapacitivepwrfactor_min")
										tenScada.Slow_CapableCapacitivePwrFactor_Max = Slow_CapableCapacitivePwrFactor.GetFloat64("slow_capablecapacitivepwrfactor_max")
										tenScada.Slow_CapableCapacitivePwrFactor_Count = Slow_CapableCapacitivePwrFactor.GetInt("slow_capablecapacitivepwrfactor_count")
									} else {
										tenScada.Slow_CapableCapacitivePwrFactor = emptyValueBig
										tenScada.Slow_CapableCapacitivePwrFactor_StdDev = emptyValueBig
										tenScada.Slow_CapableCapacitivePwrFactor_Min = emptyValueBig
										tenScada.Slow_CapableCapacitivePwrFactor_Max = emptyValueBig

									}

									Fast_Total_Production_kWh := Fast_Total_Production_kWhMap[key]
									if Fast_Total_Production_kWh != nil {
										tenScada.Fast_Total_Production_kWh = Fast_Total_Production_kWh.GetFloat64("fast_total_production_kwh")
										tenScada.Fast_Total_Production_kWh_StdDev = Fast_Total_Production_kWh.GetFloat64("fast_total_production_kwh_stddev")
										tenScada.Fast_Total_Production_kWh_Min = Fast_Total_Production_kWh.GetFloat64("fast_total_production_kwh_min")
										tenScada.Fast_Total_Production_kWh_Max = Fast_Total_Production_kWh.GetFloat64("fast_total_production_kwh_max")
										tenScada.Fast_Total_Production_kWh_Count = Fast_Total_Production_kWh.GetInt("fast_total_production_kwh_count")
									} else {
										tenScada.Fast_Total_Production_kWh = emptyValueBig
										tenScada.Fast_Total_Production_kWh_StdDev = emptyValueBig
										tenScada.Fast_Total_Production_kWh_Min = emptyValueBig
										tenScada.Fast_Total_Production_kWh_Max = emptyValueBig

									}

									Fast_Total_Prod_Day_kWh := Fast_Total_Prod_Day_kWhMap[key]
									if Fast_Total_Prod_Day_kWh != nil {
										tenScada.Fast_Total_Prod_Day_kWh = Fast_Total_Prod_Day_kWh.GetFloat64("fast_total_prod_day_kwh")
										tenScada.Fast_Total_Prod_Day_kWh_StdDev = Fast_Total_Prod_Day_kWh.GetFloat64("fast_total_prod_day_kwh_stddev")
										tenScada.Fast_Total_Prod_Day_kWh_Min = Fast_Total_Prod_Day_kWh.GetFloat64("fast_total_prod_day_kwh_min")
										tenScada.Fast_Total_Prod_Day_kWh_Max = Fast_Total_Prod_Day_kWh.GetFloat64("fast_total_prod_day_kwh_max")
										tenScada.Fast_Total_Prod_Day_kWh_Count = Fast_Total_Prod_Day_kWh.GetInt("fast_total_prod_day_kwh_count")
									} else {
										tenScada.Fast_Total_Prod_Day_kWh = emptyValueBig
										tenScada.Fast_Total_Prod_Day_kWh_StdDev = emptyValueBig
										tenScada.Fast_Total_Prod_Day_kWh_Min = emptyValueBig
										tenScada.Fast_Total_Prod_Day_kWh_Max = emptyValueBig

									}

									Fast_Total_Prod_Month_kWh := Fast_Total_Prod_Month_kWhMap[key]
									if Fast_Total_Prod_Month_kWh != nil {
										tenScada.Fast_Total_Prod_Month_kWh = Fast_Total_Prod_Month_kWh.GetFloat64("fast_total_prod_month_kwh")
										tenScada.Fast_Total_Prod_Month_kWh_StdDev = Fast_Total_Prod_Month_kWh.GetFloat64("fast_total_prod_month_kwh_stddev")
										tenScada.Fast_Total_Prod_Month_kWh_Min = Fast_Total_Prod_Month_kWh.GetFloat64("fast_total_prod_month_kwh_min")
										tenScada.Fast_Total_Prod_Month_kWh_Max = Fast_Total_Prod_Month_kWh.GetFloat64("fast_total_prod_month_kwh_max")
										tenScada.Fast_Total_Prod_Month_kWh_Count = Fast_Total_Prod_Month_kWh.GetInt("fast_total_prod_month_kwh_count")
									} else {
										tenScada.Fast_Total_Prod_Month_kWh = emptyValueBig
										tenScada.Fast_Total_Prod_Month_kWh_StdDev = emptyValueBig
										tenScada.Fast_Total_Prod_Month_kWh_Min = emptyValueBig
										tenScada.Fast_Total_Prod_Month_kWh_Max = emptyValueBig

									}

									Fast_ActivePowerOutPWCSell_kW := Fast_ActivePowerOutPWCSell_kWMap[key]
									if Fast_ActivePowerOutPWCSell_kW != nil {
										tenScada.Fast_ActivePowerOutPWCSell_kW = Fast_ActivePowerOutPWCSell_kW.GetFloat64("fast_activepoweroutpwcsell_kw")
										tenScada.Fast_ActivePowerOutPWCSell_kW_StdDev = Fast_ActivePowerOutPWCSell_kW.GetFloat64("fast_activepoweroutpwcsell_kw_stddev")
										tenScada.Fast_ActivePowerOutPWCSell_kW_Min = Fast_ActivePowerOutPWCSell_kW.GetFloat64("fast_activepoweroutpwcsell_kw_min")
										tenScada.Fast_ActivePowerOutPWCSell_kW_Max = Fast_ActivePowerOutPWCSell_kW.GetFloat64("fast_activepoweroutpwcsell_kw_max")
										tenScada.Fast_ActivePowerOutPWCSell_kW_Count = Fast_ActivePowerOutPWCSell_kW.GetInt("fast_activepoweroutpwcsell_kw_count")
									} else {
										tenScada.Fast_ActivePowerOutPWCSell_kW = emptyValueBig
										tenScada.Fast_ActivePowerOutPWCSell_kW_StdDev = emptyValueBig
										tenScada.Fast_ActivePowerOutPWCSell_kW_Min = emptyValueBig
										tenScada.Fast_ActivePowerOutPWCSell_kW_Max = emptyValueBig

									}

									Fast_Frequency_Hz := Fast_Frequency_HzMap[key]
									if Fast_Frequency_Hz != nil {
										tenScada.Fast_Frequency_Hz = Fast_Frequency_Hz.GetFloat64("fast_frequency_hz")
										tenScada.Fast_Frequency_Hz_StdDev = Fast_Frequency_Hz.GetFloat64("fast_frequency_hz_stddev")
										tenScada.Fast_Frequency_Hz_Min = Fast_Frequency_Hz.GetFloat64("fast_frequency_hz_min")
										tenScada.Fast_Frequency_Hz_Max = Fast_Frequency_Hz.GetFloat64("fast_frequency_hz_max")
										tenScada.Fast_Frequency_Hz_Count = Fast_Frequency_Hz.GetInt("fast_frequency_hz_count")
									} else {
										tenScada.Fast_Frequency_Hz = emptyValueBig
										tenScada.Fast_Frequency_Hz_StdDev = emptyValueBig
										tenScada.Fast_Frequency_Hz_Min = emptyValueBig
										tenScada.Fast_Frequency_Hz_Max = emptyValueBig

									}

									Slow_TempG1L2 := Slow_TempG1L2Map[key]
									if Slow_TempG1L2 != nil {
										tenScada.Slow_TempG1L2 = Slow_TempG1L2.GetFloat64("slow_tempg1l2")
										tenScada.Slow_TempG1L2_StdDev = Slow_TempG1L2.GetFloat64("slow_tempg1l2_stddev")
										tenScada.Slow_TempG1L2_Min = Slow_TempG1L2.GetFloat64("slow_tempg1l2_min")
										tenScada.Slow_TempG1L2_Max = Slow_TempG1L2.GetFloat64("slow_tempg1l2_max")
										tenScada.Slow_TempG1L2_Count = Slow_TempG1L2.GetInt("slow_tempg1l2_count")
									} else {
										tenScada.Slow_TempG1L2 = emptyValueBig
										tenScada.Slow_TempG1L2_StdDev = emptyValueBig
										tenScada.Slow_TempG1L2_Min = emptyValueBig
										tenScada.Slow_TempG1L2_Max = emptyValueBig

									}

									Slow_TempG1L3 := Slow_TempG1L3Map[key]
									if Slow_TempG1L3 != nil {
										tenScada.Slow_TempG1L3 = Slow_TempG1L3.GetFloat64("slow_tempg1l3")
										tenScada.Slow_TempG1L3_StdDev = Slow_TempG1L3.GetFloat64("slow_tempg1l3_stddev")
										tenScada.Slow_TempG1L3_Min = Slow_TempG1L3.GetFloat64("slow_tempg1l3_min")
										tenScada.Slow_TempG1L3_Max = Slow_TempG1L3.GetFloat64("slow_tempg1l3_max")
										tenScada.Slow_TempG1L3_Count = Slow_TempG1L3.GetInt("slow_tempg1l3_count")
									} else {
										tenScada.Slow_TempG1L3 = emptyValueBig
										tenScada.Slow_TempG1L3_StdDev = emptyValueBig
										tenScada.Slow_TempG1L3_Min = emptyValueBig
										tenScada.Slow_TempG1L3_Max = emptyValueBig

									}

									Slow_TempGearBoxHSSDE := Slow_TempGearBoxHSSDEMap[key]
									if Slow_TempGearBoxHSSDE != nil {
										tenScada.Slow_TempGearBoxHSSDE = Slow_TempGearBoxHSSDE.GetFloat64("slow_tempgearboxhssde")
										tenScada.Slow_TempGearBoxHSSDE_StdDev = Slow_TempGearBoxHSSDE.GetFloat64("slow_tempgearboxhssde_stddev")
										tenScada.Slow_TempGearBoxHSSDE_Min = Slow_TempGearBoxHSSDE.GetFloat64("slow_tempgearboxhssde_min")
										tenScada.Slow_TempGearBoxHSSDE_Max = Slow_TempGearBoxHSSDE.GetFloat64("slow_tempgearboxhssde_max")
										tenScada.Slow_TempGearBoxHSSDE_Count = Slow_TempGearBoxHSSDE.GetInt("slow_tempgearboxhssde_count")
									} else {
										tenScada.Slow_TempGearBoxHSSDE = emptyValueBig
										tenScada.Slow_TempGearBoxHSSDE_StdDev = emptyValueBig
										tenScada.Slow_TempGearBoxHSSDE_Min = emptyValueBig
										tenScada.Slow_TempGearBoxHSSDE_Max = emptyValueBig

									}

									Slow_TempGearBoxIMSNDE := Slow_TempGearBoxIMSNDEMap[key]
									if Slow_TempGearBoxIMSNDE != nil {
										tenScada.Slow_TempGearBoxIMSNDE = Slow_TempGearBoxIMSNDE.GetFloat64("slow_tempgearboximsnde")
										tenScada.Slow_TempGearBoxIMSNDE_StdDev = Slow_TempGearBoxIMSNDE.GetFloat64("slow_tempgearboximsnde_stddev")
										tenScada.Slow_TempGearBoxIMSNDE_Min = Slow_TempGearBoxIMSNDE.GetFloat64("slow_tempgearboximsnde_min")
										tenScada.Slow_TempGearBoxIMSNDE_Max = Slow_TempGearBoxIMSNDE.GetFloat64("slow_tempgearboximsnde_max")
										tenScada.Slow_TempGearBoxIMSNDE_Count = Slow_TempGearBoxIMSNDE.GetInt("slow_tempgearboximsnde_count")
									} else {
										tenScada.Slow_TempGearBoxIMSNDE = emptyValueBig
										tenScada.Slow_TempGearBoxIMSNDE_StdDev = emptyValueBig
										tenScada.Slow_TempGearBoxIMSNDE_Min = emptyValueBig
										tenScada.Slow_TempGearBoxIMSNDE_Max = emptyValueBig

									}

									Slow_TempOutdoor := Slow_TempOutdoorMap[key]
									if Slow_TempOutdoor != nil {
										tenScada.Slow_TempOutdoor = Slow_TempOutdoor.GetFloat64("slow_tempoutdoor")
										tenScada.Slow_TempOutdoor_StdDev = Slow_TempOutdoor.GetFloat64("slow_tempoutdoor_stddev")
										tenScada.Slow_TempOutdoor_Min = Slow_TempOutdoor.GetFloat64("slow_tempoutdoor_min")
										tenScada.Slow_TempOutdoor_Max = Slow_TempOutdoor.GetFloat64("slow_tempoutdoor_max")
										tenScada.Slow_TempOutdoor_Count = Slow_TempOutdoor.GetInt("slow_tempoutdoor_count")
									} else {
										tenScada.Slow_TempOutdoor = emptyValueBig
										tenScada.Slow_TempOutdoor_StdDev = emptyValueBig
										tenScada.Slow_TempOutdoor_Min = emptyValueBig
										tenScada.Slow_TempOutdoor_Max = emptyValueBig

									}

									Fast_PitchAccuV3 := Fast_PitchAccuV3Map[key]
									if Fast_PitchAccuV3 != nil {
										tenScada.Fast_PitchAccuV3 = Fast_PitchAccuV3.GetFloat64("fast_pitchaccuv3")
										tenScada.Fast_PitchAccuV3_StdDev = Fast_PitchAccuV3.GetFloat64("fast_pitchaccuv3_stddev")
										tenScada.Fast_PitchAccuV3_Min = Fast_PitchAccuV3.GetFloat64("fast_pitchaccuv3_min")
										tenScada.Fast_PitchAccuV3_Max = Fast_PitchAccuV3.GetFloat64("fast_pitchaccuv3_max")
										tenScada.Fast_PitchAccuV3_Count = Fast_PitchAccuV3.GetInt("fast_pitchaccuv3_count")
									} else {
										tenScada.Fast_PitchAccuV3 = emptyValueBig
										tenScada.Fast_PitchAccuV3_StdDev = emptyValueBig
										tenScada.Fast_PitchAccuV3_Min = emptyValueBig
										tenScada.Fast_PitchAccuV3_Max = emptyValueBig

									}

									Slow_TotalTurbineActiveHours := Slow_TotalTurbineActiveHoursMap[key]
									if Slow_TotalTurbineActiveHours != nil {
										tenScada.Slow_TotalTurbineActiveHours = Slow_TotalTurbineActiveHours.GetFloat64("slow_totalturbineactivehours")
										tenScada.Slow_TotalTurbineActiveHours_StdDev = Slow_TotalTurbineActiveHours.GetFloat64("slow_totalturbineactivehours_stddev")
										tenScada.Slow_TotalTurbineActiveHours_Min = Slow_TotalTurbineActiveHours.GetFloat64("slow_totalturbineactivehours_min")
										tenScada.Slow_TotalTurbineActiveHours_Max = Slow_TotalTurbineActiveHours.GetFloat64("slow_totalturbineactivehours_max")
										tenScada.Slow_TotalTurbineActiveHours_Count = Slow_TotalTurbineActiveHours.GetInt("slow_totalturbineactivehours_count")
									} else {
										tenScada.Slow_TotalTurbineActiveHours = emptyValueBig
										tenScada.Slow_TotalTurbineActiveHours_StdDev = emptyValueBig
										tenScada.Slow_TotalTurbineActiveHours_Min = emptyValueBig
										tenScada.Slow_TotalTurbineActiveHours_Max = emptyValueBig

									}

									Slow_TotalTurbineOKHours := Slow_TotalTurbineOKHoursMap[key]
									if Slow_TotalTurbineOKHours != nil {
										tenScada.Slow_TotalTurbineOKHours = Slow_TotalTurbineOKHours.GetFloat64("slow_totalturbineokhours")
										tenScada.Slow_TotalTurbineOKHours_StdDev = Slow_TotalTurbineOKHours.GetFloat64("slow_totalturbineokhours_stddev")
										tenScada.Slow_TotalTurbineOKHours_Min = Slow_TotalTurbineOKHours.GetFloat64("slow_totalturbineokhours_min")
										tenScada.Slow_TotalTurbineOKHours_Max = Slow_TotalTurbineOKHours.GetFloat64("slow_totalturbineokhours_max")
										tenScada.Slow_TotalTurbineOKHours_Count = Slow_TotalTurbineOKHours.GetInt("slow_totalturbineokhours_count")
									} else {
										tenScada.Slow_TotalTurbineOKHours = emptyValueBig
										tenScada.Slow_TotalTurbineOKHours_StdDev = emptyValueBig
										tenScada.Slow_TotalTurbineOKHours_Min = emptyValueBig
										tenScada.Slow_TotalTurbineOKHours_Max = emptyValueBig

									}

									Slow_TotalTurbineTimeAllHours := Slow_TotalTurbineTimeAllHoursMap[key]
									if Slow_TotalTurbineTimeAllHours != nil {
										tenScada.Slow_TotalTurbineTimeAllHours = Slow_TotalTurbineTimeAllHours.GetFloat64("slow_totalturbinetimeallhours")
										tenScada.Slow_TotalTurbineTimeAllHours_StdDev = Slow_TotalTurbineTimeAllHours.GetFloat64("slow_totalturbinetimeallhours_stddev")
										tenScada.Slow_TotalTurbineTimeAllHours_Min = Slow_TotalTurbineTimeAllHours.GetFloat64("slow_totalturbinetimeallhours_min")
										tenScada.Slow_TotalTurbineTimeAllHours_Max = Slow_TotalTurbineTimeAllHours.GetFloat64("slow_totalturbinetimeallhours_max")
										tenScada.Slow_TotalTurbineTimeAllHours_Count = Slow_TotalTurbineTimeAllHours.GetInt("slow_totalturbinetimeallhours_count")
									} else {
										tenScada.Slow_TotalTurbineTimeAllHours = emptyValueBig
										tenScada.Slow_TotalTurbineTimeAllHours_StdDev = emptyValueBig
										tenScada.Slow_TotalTurbineTimeAllHours_Min = emptyValueBig
										tenScada.Slow_TotalTurbineTimeAllHours_Max = emptyValueBig

									}

									Slow_TempG1L1 := Slow_TempG1L1Map[key]
									if Slow_TempG1L1 != nil {
										tenScada.Slow_TempG1L1 = Slow_TempG1L1.GetFloat64("slow_tempg1l1")
										tenScada.Slow_TempG1L1_StdDev = Slow_TempG1L1.GetFloat64("slow_tempg1l1_stddev")
										tenScada.Slow_TempG1L1_Min = Slow_TempG1L1.GetFloat64("slow_tempg1l1_min")
										tenScada.Slow_TempG1L1_Max = Slow_TempG1L1.GetFloat64("slow_tempg1l1_max")
										tenScada.Slow_TempG1L1_Count = Slow_TempG1L1.GetInt("slow_tempg1l1_count")
									} else {
										tenScada.Slow_TempG1L1 = emptyValueBig
										tenScada.Slow_TempG1L1_StdDev = emptyValueBig
										tenScada.Slow_TempG1L1_Min = emptyValueBig
										tenScada.Slow_TempG1L1_Max = emptyValueBig

									}

									Slow_TempGearBoxOilSump := Slow_TempGearBoxOilSumpMap[key]
									if Slow_TempGearBoxOilSump != nil {
										tenScada.Slow_TempGearBoxOilSump = Slow_TempGearBoxOilSump.GetFloat64("slow_tempgearboxoilsump")
										tenScada.Slow_TempGearBoxOilSump_StdDev = Slow_TempGearBoxOilSump.GetFloat64("slow_tempgearboxoilsump_stddev")
										tenScada.Slow_TempGearBoxOilSump_Min = Slow_TempGearBoxOilSump.GetFloat64("slow_tempgearboxoilsump_min")
										tenScada.Slow_TempGearBoxOilSump_Max = Slow_TempGearBoxOilSump.GetFloat64("slow_tempgearboxoilsump_max")
										tenScada.Slow_TempGearBoxOilSump_Count = Slow_TempGearBoxOilSump.GetInt("slow_tempgearboxoilsump_count")
									} else {
										tenScada.Slow_TempGearBoxOilSump = emptyValueBig
										tenScada.Slow_TempGearBoxOilSump_StdDev = emptyValueBig
										tenScada.Slow_TempGearBoxOilSump_Min = emptyValueBig
										tenScada.Slow_TempGearBoxOilSump_Max = emptyValueBig

									}

									Fast_PitchAccuV2 := Fast_PitchAccuV2Map[key]
									if Fast_PitchAccuV2 != nil {
										tenScada.Fast_PitchAccuV2 = Fast_PitchAccuV2.GetFloat64("fast_pitchaccuv2")
										tenScada.Fast_PitchAccuV2_StdDev = Fast_PitchAccuV2.GetFloat64("fast_pitchaccuv2_stddev")
										tenScada.Fast_PitchAccuV2_Min = Fast_PitchAccuV2.GetFloat64("fast_pitchaccuv2_min")
										tenScada.Fast_PitchAccuV2_Max = Fast_PitchAccuV2.GetFloat64("fast_pitchaccuv2_max")
										tenScada.Fast_PitchAccuV2_Count = Fast_PitchAccuV2.GetInt("fast_pitchaccuv2_count")
									} else {
										tenScada.Fast_PitchAccuV2 = emptyValueBig
										tenScada.Fast_PitchAccuV2_StdDev = emptyValueBig
										tenScada.Fast_PitchAccuV2_Min = emptyValueBig
										tenScada.Fast_PitchAccuV2_Max = emptyValueBig

									}

									Slow_TotalGridOkHours := Slow_TotalGridOkHoursMap[key]
									if Slow_TotalGridOkHours != nil {
										tenScada.Slow_TotalGridOkHours = Slow_TotalGridOkHours.GetFloat64("slow_totalgridokhours")
										tenScada.Slow_TotalGridOkHours_StdDev = Slow_TotalGridOkHours.GetFloat64("slow_totalgridokhours_stddev")
										tenScada.Slow_TotalGridOkHours_Min = Slow_TotalGridOkHours.GetFloat64("slow_totalgridokhours_min")
										tenScada.Slow_TotalGridOkHours_Max = Slow_TotalGridOkHours.GetFloat64("slow_totalgridokhours_max")
										tenScada.Slow_TotalGridOkHours_Count = Slow_TotalGridOkHours.GetInt("slow_totalgridokhours_count")
									} else {
										tenScada.Slow_TotalGridOkHours = emptyValueBig
										tenScada.Slow_TotalGridOkHours_StdDev = emptyValueBig
										tenScada.Slow_TotalGridOkHours_Min = emptyValueBig
										tenScada.Slow_TotalGridOkHours_Max = emptyValueBig

									}

									Slow_TotalActPowerOut_kWh := Slow_TotalActPowerOut_kWhMap[key]
									if Slow_TotalActPowerOut_kWh != nil {
										tenScada.Slow_TotalActPowerOut_kWh = Slow_TotalActPowerOut_kWh.GetFloat64("slow_totalactpowerout_kwh")
										tenScada.Slow_TotalActPowerOut_kWh_StdDev = Slow_TotalActPowerOut_kWh.GetFloat64("slow_totalactpowerout_kwh_stddev")
										tenScada.Slow_TotalActPowerOut_kWh_Min = Slow_TotalActPowerOut_kWh.GetFloat64("slow_totalactpowerout_kwh_min")
										tenScada.Slow_TotalActPowerOut_kWh_Max = Slow_TotalActPowerOut_kWh.GetFloat64("slow_totalactpowerout_kwh_max")
										tenScada.Slow_TotalActPowerOut_kWh_Count = Slow_TotalActPowerOut_kWh.GetInt("slow_totalactpowerout_kwh_count")
									} else {
										tenScada.Slow_TotalActPowerOut_kWh = emptyValueBig
										tenScada.Slow_TotalActPowerOut_kWh_StdDev = emptyValueBig
										tenScada.Slow_TotalActPowerOut_kWh_Min = emptyValueBig
										tenScada.Slow_TotalActPowerOut_kWh_Max = emptyValueBig

									}

									Fast_YawService := Fast_YawServiceMap[key]
									if Fast_YawService != nil {
										tenScada.Fast_YawService = Fast_YawService.GetFloat64("fast_yawservice")
										tenScada.Fast_YawService_StdDev = Fast_YawService.GetFloat64("fast_yawservice_stddev")
										tenScada.Fast_YawService_Min = Fast_YawService.GetFloat64("fast_yawservice_min")
										tenScada.Fast_YawService_Max = Fast_YawService.GetFloat64("fast_yawservice_max")
										tenScada.Fast_YawService_Count = Fast_YawService.GetInt("fast_yawservice_count")
									} else {
										tenScada.Fast_YawService = emptyValueBig
										tenScada.Fast_YawService_StdDev = emptyValueBig
										tenScada.Fast_YawService_Min = emptyValueBig
										tenScada.Fast_YawService_Max = emptyValueBig

									}

									Fast_YawAngle := Fast_YawAngleMap[key]
									if Fast_YawAngle != nil {
										tenScada.Fast_YawAngle = Fast_YawAngle.GetFloat64("fast_yawangle")
										tenScada.Fast_YawAngle_StdDev = Fast_YawAngle.GetFloat64("fast_yawangle_stddev")
										tenScada.Fast_YawAngle_Min = Fast_YawAngle.GetFloat64("fast_yawangle_min")
										tenScada.Fast_YawAngle_Max = Fast_YawAngle.GetFloat64("fast_yawangle_max")
										tenScada.Fast_YawAngle_Count = Fast_YawAngle.GetInt("fast_yawangle_count")
									} else {
										tenScada.Fast_YawAngle = emptyValueBig
										tenScada.Fast_YawAngle_StdDev = emptyValueBig
										tenScada.Fast_YawAngle_Min = emptyValueBig
										tenScada.Fast_YawAngle_Max = emptyValueBig

									}

									Slow_CapableInductivePwrFactor := Slow_CapableInductivePwrFactorMap[key]
									if Slow_CapableInductivePwrFactor != nil {
										tenScada.Slow_CapableInductivePwrFactor = Slow_CapableInductivePwrFactor.GetFloat64("slow_capableinductivepwrfactor")
										tenScada.Slow_CapableInductivePwrFactor_StdDev = Slow_CapableInductivePwrFactor.GetFloat64("slow_capableinductivepwrfactor_stddev")
										tenScada.Slow_CapableInductivePwrFactor_Min = Slow_CapableInductivePwrFactor.GetFloat64("slow_capableinductivepwrfactor_min")
										tenScada.Slow_CapableInductivePwrFactor_Max = Slow_CapableInductivePwrFactor.GetFloat64("slow_capableinductivepwrfactor_max")
										tenScada.Slow_CapableInductivePwrFactor_Count = Slow_CapableInductivePwrFactor.GetInt("slow_capableinductivepwrfactor_count")
									} else {
										tenScada.Slow_CapableInductivePwrFactor = emptyValueBig
										tenScada.Slow_CapableInductivePwrFactor_StdDev = emptyValueBig
										tenScada.Slow_CapableInductivePwrFactor_Min = emptyValueBig
										tenScada.Slow_CapableInductivePwrFactor_Max = emptyValueBig

									}

									Slow_TempGearBoxHSSNDE := Slow_TempGearBoxHSSNDEMap[key]
									if Slow_TempGearBoxHSSNDE != nil {
										tenScada.Slow_TempGearBoxHSSNDE = Slow_TempGearBoxHSSNDE.GetFloat64("slow_tempgearboxhssnde")
										tenScada.Slow_TempGearBoxHSSNDE_StdDev = Slow_TempGearBoxHSSNDE.GetFloat64("slow_tempgearboxhssnde_stddev")
										tenScada.Slow_TempGearBoxHSSNDE_Min = Slow_TempGearBoxHSSNDE.GetFloat64("slow_tempgearboxhssnde_min")
										tenScada.Slow_TempGearBoxHSSNDE_Max = Slow_TempGearBoxHSSNDE.GetFloat64("slow_tempgearboxhssnde_max")
										tenScada.Slow_TempGearBoxHSSNDE_Count = Slow_TempGearBoxHSSNDE.GetInt("slow_tempgearboxhssnde_count")
									} else {
										tenScada.Slow_TempGearBoxHSSNDE = emptyValueBig
										tenScada.Slow_TempGearBoxHSSNDE_StdDev = emptyValueBig
										tenScada.Slow_TempGearBoxHSSNDE_Min = emptyValueBig
										tenScada.Slow_TempGearBoxHSSNDE_Max = emptyValueBig

									}

									Slow_TempHubBearing := Slow_TempHubBearingMap[key]
									if Slow_TempHubBearing != nil {
										tenScada.Slow_TempHubBearing = Slow_TempHubBearing.GetFloat64("slow_temphubbearing")
										tenScada.Slow_TempHubBearing_StdDev = Slow_TempHubBearing.GetFloat64("slow_temphubbearing_stddev")
										tenScada.Slow_TempHubBearing_Min = Slow_TempHubBearing.GetFloat64("slow_temphubbearing_min")
										tenScada.Slow_TempHubBearing_Max = Slow_TempHubBearing.GetFloat64("slow_temphubbearing_max")
										tenScada.Slow_TempHubBearing_Count = Slow_TempHubBearing.GetInt("slow_temphubbearing_count")
									} else {
										tenScada.Slow_TempHubBearing = emptyValueBig
										tenScada.Slow_TempHubBearing_StdDev = emptyValueBig
										tenScada.Slow_TempHubBearing_Min = emptyValueBig
										tenScada.Slow_TempHubBearing_Max = emptyValueBig

									}

									Slow_TotalG1ActiveHours := Slow_TotalG1ActiveHoursMap[key]
									if Slow_TotalG1ActiveHours != nil {
										tenScada.Slow_TotalG1ActiveHours = Slow_TotalG1ActiveHours.GetFloat64("slow_totalg1activehours")
										tenScada.Slow_TotalG1ActiveHours_StdDev = Slow_TotalG1ActiveHours.GetFloat64("slow_totalg1activehours_stddev")
										tenScada.Slow_TotalG1ActiveHours_Min = Slow_TotalG1ActiveHours.GetFloat64("slow_totalg1activehours_min")
										tenScada.Slow_TotalG1ActiveHours_Max = Slow_TotalG1ActiveHours.GetFloat64("slow_totalg1activehours_max")
										tenScada.Slow_TotalG1ActiveHours_Count = Slow_TotalG1ActiveHours.GetInt("slow_totalg1activehours_count")
									} else {
										tenScada.Slow_TotalG1ActiveHours = emptyValueBig
										tenScada.Slow_TotalG1ActiveHours_StdDev = emptyValueBig
										tenScada.Slow_TotalG1ActiveHours_Min = emptyValueBig
										tenScada.Slow_TotalG1ActiveHours_Max = emptyValueBig

									}

									Slow_TotalActPowerOutG1_kWh := Slow_TotalActPowerOutG1_kWhMap[key]
									if Slow_TotalActPowerOutG1_kWh != nil {
										tenScada.Slow_TotalActPowerOutG1_kWh = Slow_TotalActPowerOutG1_kWh.GetFloat64("slow_totalactpoweroutg1_kwh")
										tenScada.Slow_TotalActPowerOutG1_kWh_StdDev = Slow_TotalActPowerOutG1_kWh.GetFloat64("slow_totalactpoweroutg1_kwh_stddev")
										tenScada.Slow_TotalActPowerOutG1_kWh_Min = Slow_TotalActPowerOutG1_kWh.GetFloat64("slow_totalactpoweroutg1_kwh_min")
										tenScada.Slow_TotalActPowerOutG1_kWh_Max = Slow_TotalActPowerOutG1_kWh.GetFloat64("slow_totalactpoweroutg1_kwh_max")
										tenScada.Slow_TotalActPowerOutG1_kWh_Count = Slow_TotalActPowerOutG1_kWh.GetInt("slow_totalactpoweroutg1_kwh_count")
									} else {
										tenScada.Slow_TotalActPowerOutG1_kWh = emptyValueBig
										tenScada.Slow_TotalActPowerOutG1_kWh_StdDev = emptyValueBig
										tenScada.Slow_TotalActPowerOutG1_kWh_Min = emptyValueBig
										tenScada.Slow_TotalActPowerOutG1_kWh_Max = emptyValueBig

									}

									Slow_TotalReactPowerInG1_kVArh := Slow_TotalReactPowerInG1_kVArhMap[key]
									if Slow_TotalReactPowerInG1_kVArh != nil {
										tenScada.Slow_TotalReactPowerInG1_kVArh = Slow_TotalReactPowerInG1_kVArh.GetFloat64("slow_totalreactpowering1_kvarh")
										tenScada.Slow_TotalReactPowerInG1_kVArh_StdDev = Slow_TotalReactPowerInG1_kVArh.GetFloat64("slow_totalreactpowering1_kvarh_stddev")
										tenScada.Slow_TotalReactPowerInG1_kVArh_Min = Slow_TotalReactPowerInG1_kVArh.GetFloat64("slow_totalreactpowering1_kvarh_min")
										tenScada.Slow_TotalReactPowerInG1_kVArh_Max = Slow_TotalReactPowerInG1_kVArh.GetFloat64("slow_totalreactpowering1_kvarh_max")
										tenScada.Slow_TotalReactPowerInG1_kVArh_Count = Slow_TotalReactPowerInG1_kVArh.GetInt("slow_totalreactpowering1_kvarh_count")
									} else {
										tenScada.Slow_TotalReactPowerInG1_kVArh = emptyValueBig
										tenScada.Slow_TotalReactPowerInG1_kVArh_StdDev = emptyValueBig
										tenScada.Slow_TotalReactPowerInG1_kVArh_Min = emptyValueBig
										tenScada.Slow_TotalReactPowerInG1_kVArh_Max = emptyValueBig

									}

									Slow_NacelleDrill := Slow_NacelleDrillMap[key]
									if Slow_NacelleDrill != nil {
										tenScada.Slow_NacelleDrill = Slow_NacelleDrill.GetFloat64("slow_nacelledrill")
										tenScada.Slow_NacelleDrill_StdDev = Slow_NacelleDrill.GetFloat64("slow_nacelledrill_stddev")
										tenScada.Slow_NacelleDrill_Min = Slow_NacelleDrill.GetFloat64("slow_nacelledrill_min")
										tenScada.Slow_NacelleDrill_Max = Slow_NacelleDrill.GetFloat64("slow_nacelledrill_max")
										tenScada.Slow_NacelleDrill_Count = Slow_NacelleDrill.GetInt("slow_nacelledrill_count")
									} else {
										tenScada.Slow_NacelleDrill = emptyValueBig
										tenScada.Slow_NacelleDrill_StdDev = emptyValueBig
										tenScada.Slow_NacelleDrill_Min = emptyValueBig
										tenScada.Slow_NacelleDrill_Max = emptyValueBig

									}

									Slow_TempGearBoxIMSDE := Slow_TempGearBoxIMSDEMap[key]
									if Slow_TempGearBoxIMSDE != nil {
										tenScada.Slow_TempGearBoxIMSDE = Slow_TempGearBoxIMSDE.GetFloat64("slow_tempgearboximsde")
										tenScada.Slow_TempGearBoxIMSDE_StdDev = Slow_TempGearBoxIMSDE.GetFloat64("slow_tempgearboximsde_stddev")
										tenScada.Slow_TempGearBoxIMSDE_Min = Slow_TempGearBoxIMSDE.GetFloat64("slow_tempgearboximsde_min")
										tenScada.Slow_TempGearBoxIMSDE_Max = Slow_TempGearBoxIMSDE.GetFloat64("slow_tempgearboximsde_max")
										tenScada.Slow_TempGearBoxIMSDE_Count = Slow_TempGearBoxIMSDE.GetInt("slow_tempgearboximsde_count")
									} else {
										tenScada.Slow_TempGearBoxIMSDE = emptyValueBig
										tenScada.Slow_TempGearBoxIMSDE_StdDev = emptyValueBig
										tenScada.Slow_TempGearBoxIMSDE_Min = emptyValueBig
										tenScada.Slow_TempGearBoxIMSDE_Max = emptyValueBig

									}

									Fast_Total_Operating_hrs := Fast_Total_Operating_hrsMap[key]
									if Fast_Total_Operating_hrs != nil {
										tenScada.Fast_Total_Operating_hrs = Fast_Total_Operating_hrs.GetFloat64("fast_total_operating_hrs")
										tenScada.Fast_Total_Operating_hrs_StdDev = Fast_Total_Operating_hrs.GetFloat64("fast_total_operating_hrs_stddev")
										tenScada.Fast_Total_Operating_hrs_Min = Fast_Total_Operating_hrs.GetFloat64("fast_total_operating_hrs_min")
										tenScada.Fast_Total_Operating_hrs_Max = Fast_Total_Operating_hrs.GetFloat64("fast_total_operating_hrs_max")
										tenScada.Fast_Total_Operating_hrs_Count = Fast_Total_Operating_hrs.GetInt("fast_total_operating_hrs_count")
									} else {
										tenScada.Fast_Total_Operating_hrs = emptyValueBig
										tenScada.Fast_Total_Operating_hrs_StdDev = emptyValueBig
										tenScada.Fast_Total_Operating_hrs_Min = emptyValueBig
										tenScada.Fast_Total_Operating_hrs_Max = emptyValueBig

									}

									Slow_TempNacelle := Slow_TempNacelleMap[key]
									if Slow_TempNacelle != nil {
										tenScada.Slow_TempNacelle = Slow_TempNacelle.GetFloat64("slow_tempnacelle")
										tenScada.Slow_TempNacelle_StdDev = Slow_TempNacelle.GetFloat64("slow_tempnacelle_stddev")
										tenScada.Slow_TempNacelle_Min = Slow_TempNacelle.GetFloat64("slow_tempnacelle_min")
										tenScada.Slow_TempNacelle_Max = Slow_TempNacelle.GetFloat64("slow_tempnacelle_max")
										tenScada.Slow_TempNacelle_Count = Slow_TempNacelle.GetInt("slow_tempnacelle_count")
									} else {
										tenScada.Slow_TempNacelle = emptyValueBig
										tenScada.Slow_TempNacelle_StdDev = emptyValueBig
										tenScada.Slow_TempNacelle_Min = emptyValueBig
										tenScada.Slow_TempNacelle_Max = emptyValueBig

									}

									Fast_Total_Grid_OK_hrs := Fast_Total_Grid_OK_hrsMap[key]
									if Fast_Total_Grid_OK_hrs != nil {
										tenScada.Fast_Total_Grid_OK_hrs = Fast_Total_Grid_OK_hrs.GetFloat64("fast_total_grid_ok_hrs")
										tenScada.Fast_Total_Grid_OK_hrs_StdDev = Fast_Total_Grid_OK_hrs.GetFloat64("fast_total_grid_ok_hrs_stddev")
										tenScada.Fast_Total_Grid_OK_hrs_Min = Fast_Total_Grid_OK_hrs.GetFloat64("fast_total_grid_ok_hrs_min")
										tenScada.Fast_Total_Grid_OK_hrs_Max = Fast_Total_Grid_OK_hrs.GetFloat64("fast_total_grid_ok_hrs_max")
										tenScada.Fast_Total_Grid_OK_hrs_Count = Fast_Total_Grid_OK_hrs.GetInt("fast_total_grid_ok_hrs_count")
									} else {
										tenScada.Fast_Total_Grid_OK_hrs = emptyValueBig
										tenScada.Fast_Total_Grid_OK_hrs_StdDev = emptyValueBig
										tenScada.Fast_Total_Grid_OK_hrs_Min = emptyValueBig
										tenScada.Fast_Total_Grid_OK_hrs_Max = emptyValueBig

									}

									Fast_Total_WTG_OK_hrs := Fast_Total_WTG_OK_hrsMap[key]
									if Fast_Total_WTG_OK_hrs != nil {
										tenScada.Fast_Total_WTG_OK_hrs = Fast_Total_WTG_OK_hrs.GetFloat64("fast_total_wtg_ok_hrs")
										tenScada.Fast_Total_WTG_OK_hrs_StdDev = Fast_Total_WTG_OK_hrs.GetFloat64("fast_total_wtg_ok_hrs_stddev")
										tenScada.Fast_Total_WTG_OK_hrs_Min = Fast_Total_WTG_OK_hrs.GetFloat64("fast_total_wtg_ok_hrs_min")
										tenScada.Fast_Total_WTG_OK_hrs_Max = Fast_Total_WTG_OK_hrs.GetFloat64("fast_total_wtg_ok_hrs_max")
										tenScada.Fast_Total_WTG_OK_hrs_Count = Fast_Total_WTG_OK_hrs.GetInt("fast_total_wtg_ok_hrs_count")
									} else {
										tenScada.Fast_Total_WTG_OK_hrs = emptyValueBig
										tenScada.Fast_Total_WTG_OK_hrs_StdDev = emptyValueBig
										tenScada.Fast_Total_WTG_OK_hrs_Min = emptyValueBig
										tenScada.Fast_Total_WTG_OK_hrs_Max = emptyValueBig

									}

									Slow_TempCabinetTopBox := Slow_TempCabinetTopBoxMap[key]
									if Slow_TempCabinetTopBox != nil {
										tenScada.Slow_TempCabinetTopBox = Slow_TempCabinetTopBox.GetFloat64("slow_tempcabinettopbox")
										tenScada.Slow_TempCabinetTopBox_StdDev = Slow_TempCabinetTopBox.GetFloat64("slow_tempcabinettopbox_stddev")
										tenScada.Slow_TempCabinetTopBox_Min = Slow_TempCabinetTopBox.GetFloat64("slow_tempcabinettopbox_min")
										tenScada.Slow_TempCabinetTopBox_Max = Slow_TempCabinetTopBox.GetFloat64("slow_tempcabinettopbox_max")
										tenScada.Slow_TempCabinetTopBox_Count = Slow_TempCabinetTopBox.GetInt("slow_tempcabinettopbox_count")
									} else {
										tenScada.Slow_TempCabinetTopBox = emptyValueBig
										tenScada.Slow_TempCabinetTopBox_StdDev = emptyValueBig
										tenScada.Slow_TempCabinetTopBox_Min = emptyValueBig
										tenScada.Slow_TempCabinetTopBox_Max = emptyValueBig

									}

									Slow_TempGeneratorBearingNDE := Slow_TempGeneratorBearingNDEMap[key]
									if Slow_TempGeneratorBearingNDE != nil {
										tenScada.Slow_TempGeneratorBearingNDE = Slow_TempGeneratorBearingNDE.GetFloat64("slow_tempgeneratorbearingnde")
										tenScada.Slow_TempGeneratorBearingNDE_StdDev = Slow_TempGeneratorBearingNDE.GetFloat64("slow_tempgeneratorbearingnde_stddev")
										tenScada.Slow_TempGeneratorBearingNDE_Min = Slow_TempGeneratorBearingNDE.GetFloat64("slow_tempgeneratorbearingnde_min")
										tenScada.Slow_TempGeneratorBearingNDE_Max = Slow_TempGeneratorBearingNDE.GetFloat64("slow_tempgeneratorbearingnde_max")
										tenScada.Slow_TempGeneratorBearingNDE_Count = Slow_TempGeneratorBearingNDE.GetInt("slow_tempgeneratorbearingnde_count")
									} else {
										tenScada.Slow_TempGeneratorBearingNDE = emptyValueBig
										tenScada.Slow_TempGeneratorBearingNDE_StdDev = emptyValueBig
										tenScada.Slow_TempGeneratorBearingNDE_Min = emptyValueBig
										tenScada.Slow_TempGeneratorBearingNDE_Max = emptyValueBig

									}

									Fast_Total_Access_hrs := Fast_Total_Access_hrsMap[key]
									if Fast_Total_Access_hrs != nil {
										tenScada.Fast_Total_Access_hrs = Fast_Total_Access_hrs.GetFloat64("fast_total_access_hrs")
										tenScada.Fast_Total_Access_hrs_StdDev = Fast_Total_Access_hrs.GetFloat64("fast_total_access_hrs_stddev")
										tenScada.Fast_Total_Access_hrs_Min = Fast_Total_Access_hrs.GetFloat64("fast_total_access_hrs_min")
										tenScada.Fast_Total_Access_hrs_Max = Fast_Total_Access_hrs.GetFloat64("fast_total_access_hrs_max")
										tenScada.Fast_Total_Access_hrs_Count = Fast_Total_Access_hrs.GetInt("fast_total_access_hrs_count")
									} else {
										tenScada.Fast_Total_Access_hrs = emptyValueBig
										tenScada.Fast_Total_Access_hrs_StdDev = emptyValueBig
										tenScada.Fast_Total_Access_hrs_Min = emptyValueBig
										tenScada.Fast_Total_Access_hrs_Max = emptyValueBig

									}

									Slow_TempBottomPowerSection := Slow_TempBottomPowerSectionMap[key]
									if Slow_TempBottomPowerSection != nil {
										tenScada.Slow_TempBottomPowerSection = Slow_TempBottomPowerSection.GetFloat64("slow_tempbottompowersection")
										tenScada.Slow_TempBottomPowerSection_StdDev = Slow_TempBottomPowerSection.GetFloat64("slow_tempbottompowersection_stddev")
										tenScada.Slow_TempBottomPowerSection_Min = Slow_TempBottomPowerSection.GetFloat64("slow_tempbottompowersection_min")
										tenScada.Slow_TempBottomPowerSection_Max = Slow_TempBottomPowerSection.GetFloat64("slow_tempbottompowersection_max")
										tenScada.Slow_TempBottomPowerSection_Count = Slow_TempBottomPowerSection.GetInt("slow_tempbottompowersection_count")
									} else {
										tenScada.Slow_TempBottomPowerSection = emptyValueBig
										tenScada.Slow_TempBottomPowerSection_StdDev = emptyValueBig
										tenScada.Slow_TempBottomPowerSection_Min = emptyValueBig
										tenScada.Slow_TempBottomPowerSection_Max = emptyValueBig

									}

									Slow_TempGeneratorBearingDE := Slow_TempGeneratorBearingDEMap[key]
									if Slow_TempGeneratorBearingDE != nil {
										tenScada.Slow_TempGeneratorBearingDE = Slow_TempGeneratorBearingDE.GetFloat64("slow_tempgeneratorbearingde")
										tenScada.Slow_TempGeneratorBearingDE_StdDev = Slow_TempGeneratorBearingDE.GetFloat64("slow_tempgeneratorbearingde_stddev")
										tenScada.Slow_TempGeneratorBearingDE_Min = Slow_TempGeneratorBearingDE.GetFloat64("slow_tempgeneratorbearingde_min")
										tenScada.Slow_TempGeneratorBearingDE_Max = Slow_TempGeneratorBearingDE.GetFloat64("slow_tempgeneratorbearingde_max")
										tenScada.Slow_TempGeneratorBearingDE_Count = Slow_TempGeneratorBearingDE.GetInt("slow_tempgeneratorbearingde_count")
									} else {
										tenScada.Slow_TempGeneratorBearingDE = emptyValueBig
										tenScada.Slow_TempGeneratorBearingDE_StdDev = emptyValueBig
										tenScada.Slow_TempGeneratorBearingDE_Min = emptyValueBig
										tenScada.Slow_TempGeneratorBearingDE_Max = emptyValueBig

									}

									Slow_TotalReactPowerIn_kVArh := Slow_TotalReactPowerIn_kVArhMap[key]
									if Slow_TotalReactPowerIn_kVArh != nil {
										tenScada.Slow_TotalReactPowerIn_kVArh = Slow_TotalReactPowerIn_kVArh.GetFloat64("slow_totalreactpowerin_kvarh")
										tenScada.Slow_TotalReactPowerIn_kVArh_StdDev = Slow_TotalReactPowerIn_kVArh.GetFloat64("slow_totalreactpowerin_kvarh_stddev")
										tenScada.Slow_TotalReactPowerIn_kVArh_Min = Slow_TotalReactPowerIn_kVArh.GetFloat64("slow_totalreactpowerin_kvarh_min")
										tenScada.Slow_TotalReactPowerIn_kVArh_Max = Slow_TotalReactPowerIn_kVArh.GetFloat64("slow_totalreactpowerin_kvarh_max")
										tenScada.Slow_TotalReactPowerIn_kVArh_Count = Slow_TotalReactPowerIn_kVArh.GetInt("slow_totalreactpowerin_kvarh_count")
									} else {
										tenScada.Slow_TotalReactPowerIn_kVArh = emptyValueBig
										tenScada.Slow_TotalReactPowerIn_kVArh_StdDev = emptyValueBig
										tenScada.Slow_TotalReactPowerIn_kVArh_Min = emptyValueBig
										tenScada.Slow_TotalReactPowerIn_kVArh_Max = emptyValueBig

									}

									Slow_TempBottomControlSection := Slow_TempBottomControlSectionMap[key]
									if Slow_TempBottomControlSection != nil {
										tenScada.Slow_TempBottomControlSection = Slow_TempBottomControlSection.GetFloat64("slow_tempbottomcontrolsection")
										tenScada.Slow_TempBottomControlSection_StdDev = Slow_TempBottomControlSection.GetFloat64("slow_tempbottomcontrolsection_stddev")
										tenScada.Slow_TempBottomControlSection_Min = Slow_TempBottomControlSection.GetFloat64("slow_tempbottomcontrolsection_min")
										tenScada.Slow_TempBottomControlSection_Max = Slow_TempBottomControlSection.GetFloat64("slow_tempbottomcontrolsection_max")
										tenScada.Slow_TempBottomControlSection_Count = Slow_TempBottomControlSection.GetInt("slow_tempbottomcontrolsection_count")
									} else {
										tenScada.Slow_TempBottomControlSection = emptyValueBig
										tenScada.Slow_TempBottomControlSection_StdDev = emptyValueBig
										tenScada.Slow_TempBottomControlSection_Min = emptyValueBig
										tenScada.Slow_TempBottomControlSection_Max = emptyValueBig

									}

									Slow_TempConv1 := Slow_TempConv1Map[key]
									if Slow_TempConv1 != nil {
										tenScada.Slow_TempConv1 = Slow_TempConv1.GetFloat64("slow_tempconv1")
										tenScada.Slow_TempConv1_StdDev = Slow_TempConv1.GetFloat64("slow_tempconv1_stddev")
										tenScada.Slow_TempConv1_Min = Slow_TempConv1.GetFloat64("slow_tempconv1_min")
										tenScada.Slow_TempConv1_Max = Slow_TempConv1.GetFloat64("slow_tempconv1_max")
										tenScada.Slow_TempConv1_Count = Slow_TempConv1.GetInt("slow_tempconv1_count")
									} else {
										tenScada.Slow_TempConv1 = emptyValueBig
										tenScada.Slow_TempConv1_StdDev = emptyValueBig
										tenScada.Slow_TempConv1_Min = emptyValueBig
										tenScada.Slow_TempConv1_Max = emptyValueBig

									}

									Fast_ActivePowerRated_kW := Fast_ActivePowerRated_kWMap[key]
									if Fast_ActivePowerRated_kW != nil {
										tenScada.Fast_ActivePowerRated_kW = Fast_ActivePowerRated_kW.GetFloat64("fast_activepowerrated_kw")
										tenScada.Fast_ActivePowerRated_kW_StdDev = Fast_ActivePowerRated_kW.GetFloat64("fast_activepowerrated_kw_stddev")
										tenScada.Fast_ActivePowerRated_kW_Min = Fast_ActivePowerRated_kW.GetFloat64("fast_activepowerrated_kw_min")
										tenScada.Fast_ActivePowerRated_kW_Max = Fast_ActivePowerRated_kW.GetFloat64("fast_activepowerrated_kw_max")
										tenScada.Fast_ActivePowerRated_kW_Count = Fast_ActivePowerRated_kW.GetInt("fast_activepowerrated_kw_count")
									} else {
										tenScada.Fast_ActivePowerRated_kW = emptyValueBig
										tenScada.Fast_ActivePowerRated_kW_StdDev = emptyValueBig
										tenScada.Fast_ActivePowerRated_kW_Min = emptyValueBig
										tenScada.Fast_ActivePowerRated_kW_Max = emptyValueBig

									}

									Fast_NodeIP := Fast_NodeIPMap[key]
									if Fast_NodeIP != nil {
										tenScada.Fast_NodeIP = Fast_NodeIP.GetFloat64("fast_nodeip")
										tenScada.Fast_NodeIP_StdDev = Fast_NodeIP.GetFloat64("fast_nodeip_stddev")
										tenScada.Fast_NodeIP_Min = Fast_NodeIP.GetFloat64("fast_nodeip_min")
										tenScada.Fast_NodeIP_Max = Fast_NodeIP.GetFloat64("fast_nodeip_max")
										tenScada.Fast_NodeIP_Count = Fast_NodeIP.GetInt("fast_nodeip_count")
									} else {
										tenScada.Fast_NodeIP = emptyValueBig
										tenScada.Fast_NodeIP_StdDev = emptyValueBig
										tenScada.Fast_NodeIP_Min = emptyValueBig
										tenScada.Fast_NodeIP_Max = emptyValueBig

									}

									Fast_PitchSpeed1 := Fast_PitchSpeed1Map[key]
									if Fast_PitchSpeed1 != nil {
										tenScada.Fast_PitchSpeed1 = Fast_PitchSpeed1.GetFloat64("fast_pitchspeed1")
										tenScada.Fast_PitchSpeed1_StdDev = Fast_PitchSpeed1.GetFloat64("fast_pitchspeed1_stddev")
										tenScada.Fast_PitchSpeed1_Min = Fast_PitchSpeed1.GetFloat64("fast_pitchspeed1_min")
										tenScada.Fast_PitchSpeed1_Max = Fast_PitchSpeed1.GetFloat64("fast_pitchspeed1_max")
										tenScada.Fast_PitchSpeed1_Count = Fast_PitchSpeed1.GetInt("fast_pitchspeed1_count")
									} else {
										tenScada.Fast_PitchSpeed1 = emptyValueBig
										tenScada.Fast_PitchSpeed1_StdDev = emptyValueBig
										tenScada.Fast_PitchSpeed1_Min = emptyValueBig
										tenScada.Fast_PitchSpeed1_Max = emptyValueBig

									}

									Slow_CFCardSize := Slow_CFCardSizeMap[key]
									if Slow_CFCardSize != nil {
										tenScada.Slow_CFCardSize = Slow_CFCardSize.GetFloat64("slow_cfcardsize")
										tenScada.Slow_CFCardSize_StdDev = Slow_CFCardSize.GetFloat64("slow_cfcardsize_stddev")
										tenScada.Slow_CFCardSize_Min = Slow_CFCardSize.GetFloat64("slow_cfcardsize_min")
										tenScada.Slow_CFCardSize_Max = Slow_CFCardSize.GetFloat64("slow_cfcardsize_max")
										tenScada.Slow_CFCardSize_Count = Slow_CFCardSize.GetInt("slow_cfcardsize_count")
									} else {
										tenScada.Slow_CFCardSize = emptyValueBig
										tenScada.Slow_CFCardSize_StdDev = emptyValueBig
										tenScada.Slow_CFCardSize_Min = emptyValueBig
										tenScada.Slow_CFCardSize_Max = emptyValueBig

									}

									Slow_CPU_Number := Slow_CPU_NumberMap[key]
									if Slow_CPU_Number != nil {
										tenScada.Slow_CPU_Number = Slow_CPU_Number.GetFloat64("slow_cpu_number")
										tenScada.Slow_CPU_Number_StdDev = Slow_CPU_Number.GetFloat64("slow_cpu_number_stddev")
										tenScada.Slow_CPU_Number_Min = Slow_CPU_Number.GetFloat64("slow_cpu_number_min")
										tenScada.Slow_CPU_Number_Max = Slow_CPU_Number.GetFloat64("slow_cpu_number_max")
										tenScada.Slow_CPU_Number_Count = Slow_CPU_Number.GetInt("slow_cpu_number_count")
									} else {
										tenScada.Slow_CPU_Number = emptyValueBig
										tenScada.Slow_CPU_Number_StdDev = emptyValueBig
										tenScada.Slow_CPU_Number_Min = emptyValueBig
										tenScada.Slow_CPU_Number_Max = emptyValueBig

									}

									Slow_CFCardSpaceLeft := Slow_CFCardSpaceLeftMap[key]
									if Slow_CFCardSpaceLeft != nil {
										tenScada.Slow_CFCardSpaceLeft = Slow_CFCardSpaceLeft.GetFloat64("slow_cfcardspaceleft")
										tenScada.Slow_CFCardSpaceLeft_StdDev = Slow_CFCardSpaceLeft.GetFloat64("slow_cfcardspaceleft_stddev")
										tenScada.Slow_CFCardSpaceLeft_Min = Slow_CFCardSpaceLeft.GetFloat64("slow_cfcardspaceleft_min")
										tenScada.Slow_CFCardSpaceLeft_Max = Slow_CFCardSpaceLeft.GetFloat64("slow_cfcardspaceleft_max")
										tenScada.Slow_CFCardSpaceLeft_Count = Slow_CFCardSpaceLeft.GetInt("slow_cfcardspaceleft_count")
									} else {
										tenScada.Slow_CFCardSpaceLeft = emptyValueBig
										tenScada.Slow_CFCardSpaceLeft_StdDev = emptyValueBig
										tenScada.Slow_CFCardSpaceLeft_Min = emptyValueBig
										tenScada.Slow_CFCardSpaceLeft_Max = emptyValueBig

									}

									Slow_TempBottomCapSection := Slow_TempBottomCapSectionMap[key]
									if Slow_TempBottomCapSection != nil {
										tenScada.Slow_TempBottomCapSection = Slow_TempBottomCapSection.GetFloat64("slow_tempbottomcapsection")
										tenScada.Slow_TempBottomCapSection_StdDev = Slow_TempBottomCapSection.GetFloat64("slow_tempbottomcapsection_stddev")
										tenScada.Slow_TempBottomCapSection_Min = Slow_TempBottomCapSection.GetFloat64("slow_tempbottomcapsection_min")
										tenScada.Slow_TempBottomCapSection_Max = Slow_TempBottomCapSection.GetFloat64("slow_tempbottomcapsection_max")
										tenScada.Slow_TempBottomCapSection_Count = Slow_TempBottomCapSection.GetInt("slow_tempbottomcapsection_count")
									} else {
										tenScada.Slow_TempBottomCapSection = emptyValueBig
										tenScada.Slow_TempBottomCapSection_StdDev = emptyValueBig
										tenScada.Slow_TempBottomCapSection_Min = emptyValueBig
										tenScada.Slow_TempBottomCapSection_Max = emptyValueBig

									}

									Slow_RatedPower := Slow_RatedPowerMap[key]
									if Slow_RatedPower != nil {
										tenScada.Slow_RatedPower = Slow_RatedPower.GetFloat64("slow_ratedpower")
										tenScada.Slow_RatedPower_StdDev = Slow_RatedPower.GetFloat64("slow_ratedpower_stddev")
										tenScada.Slow_RatedPower_Min = Slow_RatedPower.GetFloat64("slow_ratedpower_min")
										tenScada.Slow_RatedPower_Max = Slow_RatedPower.GetFloat64("slow_ratedpower_max")
										tenScada.Slow_RatedPower_Count = Slow_RatedPower.GetInt("slow_ratedpower_count")
									} else {
										tenScada.Slow_RatedPower = emptyValueBig
										tenScada.Slow_RatedPower_StdDev = emptyValueBig
										tenScada.Slow_RatedPower_Min = emptyValueBig
										tenScada.Slow_RatedPower_Max = emptyValueBig

									}

									SlowTempConv3 := Slow_TempConv3Map[key]
									if SlowTempConv3 != nil {
										tenScada.Slow_TempConv3 = SlowTempConv3.GetFloat64("slow_tempconv3")
										tenScada.Slow_TempConv3_StdDev = SlowTempConv3.GetFloat64("slow_tempconv3_stddev")
										tenScada.Slow_TempConv3_Min = SlowTempConv3.GetFloat64("slow_tempconv3_min")
										tenScada.Slow_TempConv3_Max = SlowTempConv3.GetFloat64("slow_tempconv3_max")
										tenScada.Slow_TempConv3_Count = SlowTempConv3.GetInt("slow_tempconv3_count")
									} else {
										tenScada.Slow_TempConv3 = emptyValueBig
										tenScada.Slow_TempConv3_StdDev = emptyValueBig
										tenScada.Slow_TempConv3_Min = emptyValueBig
										tenScada.Slow_TempConv3_Max = emptyValueBig

									}

									SlowTempConv2 := Slow_TempConv2Map[key]
									if SlowTempConv2 != nil {
										tenScada.Slow_TempConv2 = SlowTempConv2.GetFloat64("slow_tempconv2")
										tenScada.Slow_TempConv2_StdDev = SlowTempConv2.GetFloat64("slow_tempconv2_stddev")
										tenScada.Slow_TempConv2_Min = SlowTempConv2.GetFloat64("slow_tempconv2_min")
										tenScada.Slow_TempConv2_Max = SlowTempConv2.GetFloat64("slow_tempconv2_max")
										tenScada.Slow_TempConv2_Count = SlowTempConv2.GetInt("slow_tempconv2_count")
									} else {
										tenScada.Slow_TempConv2 = emptyValueBig
										tenScada.Slow_TempConv2_StdDev = emptyValueBig
										tenScada.Slow_TempConv2_Min = emptyValueBig
										tenScada.Slow_TempConv2_Max = emptyValueBig

									}

									Slow_TotalActPowerIn_kWh := Slow_TotalActPowerIn_kWhMap[key]
									if Slow_TotalActPowerIn_kWh != nil {
										tenScada.Slow_TotalActPowerIn_kWh = Slow_TotalActPowerIn_kWh.GetFloat64("slow_totalactpowerin_kwh")
										tenScada.Slow_TotalActPowerIn_kWh_StdDev = Slow_TotalActPowerIn_kWh.GetFloat64("slow_totalactpowerin_kwh_stddev")
										tenScada.Slow_TotalActPowerIn_kWh_Min = Slow_TotalActPowerIn_kWh.GetFloat64("slow_totalactpowerin_kwh_min")
										tenScada.Slow_TotalActPowerIn_kWh_Max = Slow_TotalActPowerIn_kWh.GetFloat64("slow_totalactpowerin_kwh_max")
										tenScada.Slow_TotalActPowerIn_kWh_Count = Slow_TotalActPowerIn_kWh.GetInt("slow_totalactpowerin_kwh_count")
									} else {
										tenScada.Slow_TotalActPowerIn_kWh = emptyValueBig
										tenScada.Slow_TotalActPowerIn_kWh_StdDev = emptyValueBig
										tenScada.Slow_TotalActPowerIn_kWh_Min = emptyValueBig
										tenScada.Slow_TotalActPowerIn_kWh_Max = emptyValueBig

									}

									Slow_TotalActPowerInG1_kWh := Slow_TotalActPowerInG1_kWhMap[key]
									if Slow_TotalActPowerInG1_kWh != nil {
										tenScada.Slow_TotalActPowerInG1_kWh = Slow_TotalActPowerInG1_kWh.GetFloat64("slow_totalactpowering1_kwh")
										tenScada.Slow_TotalActPowerInG1_kWh_StdDev = Slow_TotalActPowerInG1_kWh.GetFloat64("slow_totalactpowering1_kwh_stddev")
										tenScada.Slow_TotalActPowerInG1_kWh_Min = Slow_TotalActPowerInG1_kWh.GetFloat64("slow_totalactpowering1_kwh_min")
										tenScada.Slow_TotalActPowerInG1_kWh_Max = Slow_TotalActPowerInG1_kWh.GetFloat64("slow_totalactpowering1_kwh_max")
										tenScada.Slow_TotalActPowerInG1_kWh_Count = Slow_TotalActPowerInG1_kWh.GetInt("slow_totalactpowering1_kwh_count")
									} else {
										tenScada.Slow_TotalActPowerInG1_kWh = emptyValueBig
										tenScada.Slow_TotalActPowerInG1_kWh_StdDev = emptyValueBig
										tenScada.Slow_TotalActPowerInG1_kWh_Min = emptyValueBig
										tenScada.Slow_TotalActPowerInG1_kWh_Max = emptyValueBig

									}

									Slow_TotalActPowerInG2_kWh := Slow_TotalActPowerInG2_kWhMap[key]
									if Slow_TotalActPowerInG2_kWh != nil {
										tenScada.Slow_TotalActPowerInG2_kWh = Slow_TotalActPowerInG2_kWh.GetFloat64("slow_totalactpowering2_kwh")
										tenScada.Slow_TotalActPowerInG2_kWh_StdDev = Slow_TotalActPowerInG2_kWh.GetFloat64("slow_totalactpowering2_kwh_stddev")
										tenScada.Slow_TotalActPowerInG2_kWh_Min = Slow_TotalActPowerInG2_kWh.GetFloat64("slow_totalactpowering2_kwh_min")
										tenScada.Slow_TotalActPowerInG2_kWh_Max = Slow_TotalActPowerInG2_kWh.GetFloat64("slow_totalactpowering2_kwh_max")
										tenScada.Slow_TotalActPowerInG2_kWh_Count = Slow_TotalActPowerInG2_kWh.GetInt("slow_totalactpowering2_kwh_count")
									} else {
										tenScada.Slow_TotalActPowerInG2_kWh = emptyValueBig
										tenScada.Slow_TotalActPowerInG2_kWh_StdDev = emptyValueBig
										tenScada.Slow_TotalActPowerInG2_kWh_Min = emptyValueBig
										tenScada.Slow_TotalActPowerInG2_kWh_Max = emptyValueBig

									}

									Slow_TotalActPowerOutG2_kWh := Slow_TotalActPowerOutG2_kWhMap[key]
									if Slow_TotalActPowerOutG2_kWh != nil {
										tenScada.Slow_TotalActPowerOutG2_kWh = Slow_TotalActPowerOutG2_kWh.GetFloat64("slow_totalactpoweroutg2_kwh")
										tenScada.Slow_TotalActPowerOutG2_kWh_StdDev = Slow_TotalActPowerOutG2_kWh.GetFloat64("slow_totalactpoweroutg2_kwh_stddev")
										tenScada.Slow_TotalActPowerOutG2_kWh_Min = Slow_TotalActPowerOutG2_kWh.GetFloat64("slow_totalactpoweroutg2_kwh_min")
										tenScada.Slow_TotalActPowerOutG2_kWh_Max = Slow_TotalActPowerOutG2_kWh.GetFloat64("slow_totalactpoweroutg2_kwh_max")
										tenScada.Slow_TotalActPowerOutG2_kWh_Count = Slow_TotalActPowerOutG2_kWh.GetInt("slow_totalactpoweroutg2_kwh_count")
									} else {
										tenScada.Slow_TotalActPowerOutG2_kWh = emptyValueBig
										tenScada.Slow_TotalActPowerOutG2_kWh_StdDev = emptyValueBig
										tenScada.Slow_TotalActPowerOutG2_kWh_Min = emptyValueBig
										tenScada.Slow_TotalActPowerOutG2_kWh_Max = emptyValueBig

									}

									Slow_TotalG2ActiveHours := Slow_TotalG2ActiveHoursMap[key]
									if Slow_TotalG2ActiveHours != nil {
										tenScada.Slow_TotalG2ActiveHours = Slow_TotalG2ActiveHours.GetFloat64("slow_totalg2activehours")
										tenScada.Slow_TotalG2ActiveHours_StdDev = Slow_TotalG2ActiveHours.GetFloat64("slow_totalg2activehours_stddev")
										tenScada.Slow_TotalG2ActiveHours_Min = Slow_TotalG2ActiveHours.GetFloat64("slow_totalg2activehours_min")
										tenScada.Slow_TotalG2ActiveHours_Max = Slow_TotalG2ActiveHours.GetFloat64("slow_totalg2activehours_max")
										tenScada.Slow_TotalG2ActiveHours_Count = Slow_TotalG2ActiveHours.GetInt("slow_totalg2activehours_count")
									} else {
										tenScada.Slow_TotalG2ActiveHours = emptyValueBig
										tenScada.Slow_TotalG2ActiveHours_StdDev = emptyValueBig
										tenScada.Slow_TotalG2ActiveHours_Min = emptyValueBig
										tenScada.Slow_TotalG2ActiveHours_Max = emptyValueBig

									}

									Slow_TotalReactPowerInG2_kVArh := Slow_TotalReactPowerInG2_kVArhMap[key]
									if Slow_TotalReactPowerInG2_kVArh != nil {
										tenScada.Slow_TotalReactPowerInG2_kVArh = Slow_TotalReactPowerInG2_kVArh.GetFloat64("slow_totalreactpowering2_kvarh")
										tenScada.Slow_TotalReactPowerInG2_kVArh_StdDev = Slow_TotalReactPowerInG2_kVArh.GetFloat64("slow_totalreactpowering2_kvarh_stddev")
										tenScada.Slow_TotalReactPowerInG2_kVArh_Min = Slow_TotalReactPowerInG2_kVArh.GetFloat64("slow_totalreactpowering2_kvarh_min")
										tenScada.Slow_TotalReactPowerInG2_kVArh_Max = Slow_TotalReactPowerInG2_kVArh.GetFloat64("slow_totalreactpowering2_kvarh_max")
										tenScada.Slow_TotalReactPowerInG2_kVArh_Count = Slow_TotalReactPowerInG2_kVArh.GetInt("slow_totalreactpowering2_kvarh_count")
									} else {
										tenScada.Slow_TotalReactPowerInG2_kVArh = emptyValueBig
										tenScada.Slow_TotalReactPowerInG2_kVArh_StdDev = emptyValueBig
										tenScada.Slow_TotalReactPowerInG2_kVArh_Min = emptyValueBig
										tenScada.Slow_TotalReactPowerInG2_kVArh_Max = emptyValueBig

									}

									Slow_TotalReactPowerOut_kVArh := Slow_TotalReactPowerOut_kVArhMap[key]
									if Slow_TotalReactPowerOut_kVArh != nil {
										tenScada.Slow_TotalReactPowerOut_kVArh = Slow_TotalReactPowerOut_kVArh.GetFloat64("slow_totalreactpowerout_kvarh")
										tenScada.Slow_TotalReactPowerOut_kVArh_StdDev = Slow_TotalReactPowerOut_kVArh.GetFloat64("slow_totalreactpowerout_kvarh_stddev")
										tenScada.Slow_TotalReactPowerOut_kVArh_Min = Slow_TotalReactPowerOut_kVArh.GetFloat64("slow_totalreactpowerout_kvarh_min")
										tenScada.Slow_TotalReactPowerOut_kVArh_Max = Slow_TotalReactPowerOut_kVArh.GetFloat64("slow_totalreactpowerout_kvarh_max")
										tenScada.Slow_TotalReactPowerOut_kVArh_Count = Slow_TotalReactPowerOut_kVArh.GetInt("slow_totalreactpowerout_kvarh_count")
									} else {
										tenScada.Slow_TotalReactPowerOut_kVArh = emptyValueBig
										tenScada.Slow_TotalReactPowerOut_kVArh_StdDev = emptyValueBig
										tenScada.Slow_TotalReactPowerOut_kVArh_Min = emptyValueBig
										tenScada.Slow_TotalReactPowerOut_kVArh_Max = emptyValueBig

									}

									Slow_UTCoffset_int := Slow_UTCoffset_intMap[key]
									if Slow_UTCoffset_int != nil {
										tenScada.Slow_UTCoffset_int = Slow_UTCoffset_int.GetFloat64("slow_utcoffset_int")
										tenScada.Slow_UTCoffset_int_StdDev = Slow_UTCoffset_int.GetFloat64("slow_utcoffset_int_stddev")
										tenScada.Slow_UTCoffset_int_Min = Slow_UTCoffset_int.GetFloat64("slow_utcoffset_int_min")
										tenScada.Slow_UTCoffset_int_Max = Slow_UTCoffset_int.GetFloat64("slow_utcoffset_int_max")
										tenScada.Slow_UTCoffset_int_Count = Slow_UTCoffset_int.GetInt("slow_utcoffset_int_count")
									} else {
										tenScada.Slow_UTCoffset_int = emptyValueBig
										tenScada.Slow_UTCoffset_int_StdDev = emptyValueBig
										tenScada.Slow_UTCoffset_int_Min = emptyValueBig
										tenScada.Slow_UTCoffset_int_Max = emptyValueBig

									}

									// log.Printf("%#v \n", tenScada)
									mutex.Lock()

									/*if tenScada.Turbine == "HBR004" {
										log.Printf("tenScada: %v | %v | %v | %v \n", tenScada.ID, tenScada.TimeStamp.UTC().Format("20060102 15:04"), startTime.Format("20060102 15:04"), idSub.Get("timestampint").(int64))
									}*/

									err := ctx.Insert(tenScada)
									ErrorHandler(err, "Saving")
									mutex.Unlock()
								}

								wg.Done()
							}(data)

							counter++

							if endIndex >= countData {
								isFinish = true
							}
						}

						wg.Wait()
					}

					startTime = hpp.GenNext10Minutes(startTime)
				}

			}
		}
	}

	csr.Close()
	log.Println("End Conversion.")
	return
}

func (d *DataConversion) getStdDevAvgMinMaxCount(ctx *DataContext, timestampconverted time.Time, field string) (result []tk.M) {
	pipes := []tk.M{}

	match := tk.M{
		"timestampconverted": timestampconverted,
		field:                tk.M{"$gt": emptyValueBig},
	}

	group := tk.M{
		"_id": tk.M{
			"timestamp":   "$timestampconverted",
			"projectname": "$projectname",
			"turbine":     "$turbine",
		},
		field:             tk.M{"$avg": "$" + field},
		field + "_stddev": tk.M{"$stdDevPop": "$" + field},
		field + "_min":    tk.M{"$min": "$" + field},
		field + "_max":    tk.M{"$max": "$" + field},
		field + "_count":  tk.M{"$sum": 1},
	}

	pipes = append(pipes, tk.M{"$match": match})
	pipes = append(pipes, tk.M{"$group": group})
	pipes = append(pipes, tk.M{"$sort": tk.M{"_id": 1}})

	csr, e := ctx.Connection.NewQuery().
		From(new(ScadaThreeSecsExt).TableName()).
		Command("pipe", pipes).
		Cursor(nil)

	if e != nil {
		log.Printf("ERR: %#v \n", e.Error())
	} else {
		e = csr.Fetch(&result, 0, false)
	}

	csr.Close()

	return
}

func (d *DataConversion) getMap(list []tk.M, field string) (result map[string]tk.M) {
	result = map[string]tk.M{}

	for _, val := range list {
		id := val.Get("_id").(tk.M)
		timeStamp := id.Get("timestamp").(time.Time)
		projectName := id.GetString("projectname")
		turbine := id.GetString("turbine")

		timeStampStr := timeStamp.UTC().Format("060102_1504")
		key := timeStampStr + "#" + projectName + "#" + turbine

		value := tk.M{}

		var avg, stddev, min, max float64
		var count int

		count = val.GetInt(field + "_count")

		// log.Printf("count: %v | %#v \n", val.GetInt(field+"_count"), key)

		if count == 0 {
			avg, stddev, min, max = emptyValueBig, emptyValueBig, emptyValueBig, emptyValueBig
			log.Print("empty: %v \n", key)
		} else {
			avg, stddev, min, max = val.GetFloat64(field), val.GetFloat64(field+"_stddev"), val.GetFloat64(field+"_min"), val.GetFloat64(field+"_max")
		}

		value.Set(field, avg)
		value.Set(field+"_stddev", stddev)
		value.Set(field+"_min", min)
		value.Set(field+"_max", max)
		value.Set(field+"_count", count)

		result[key] = value
	}
	return
}
