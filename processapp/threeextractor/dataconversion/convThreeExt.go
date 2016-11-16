package dataconversion

import (
	. "eaciit/wfdemo/library/helper"
	. "eaciit/wfdemo/library/models"
	. "eaciit/wfdemo/processapp/watcher/controllers"
	"log"
	"strconv"
	"sync"
	"time"

	hpp "eaciit/wfdemo/processapp/helper"

	_ "github.com/eaciit/dbox/dbc/mongo"
	. "github.com/eaciit/orm"
	tk "github.com/eaciit/toolkit"
)

type ConvThreeExt struct {
	Ctx *DataContext
}

var (
	mutexX = &sync.Mutex{}
)

func NewConvThreeExt(ctx *DataContext) *ConvThreeExt {
	dc := new(ConvThreeExt)
	dc.Ctx = ctx

	return dc
}

func (d *ConvThreeExt) Generate(file string) (errorLine tk.M) {
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
		From(new(ScadaThreeSecs).TableName()).
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

				var startTime, endTime time.Time

				minTimeStamp := valList.Get("min_timestamp").(time.Time).UTC()
				maxTimeStamp := valList.Get("max_timestamp").(time.Time).UTC()

				// log.Printf("valList: %v | %v \n", minTimeStamp, maxTimeStamp)

				startTime, _ = time.Parse("20060102 15:04", minTimeStamp.Format("20060102 15:04"))
				endTime, _ = time.Parse("20060102 15:04", maxTimeStamp.Format("20060102 15:04"))

				for {
					if startTime.Format("2006-01-02 15:04") > endTime.Format("2006-01-02 15:04") {
						break
					}

					// log.Printf("startTime: %v \n", startTime)

					startTimeInt, _ := strconv.ParseInt(startTime.Format("200601021504"), 10, 64)

					Fast_CurrentL3List := d.getAvg(ctx, startTime, "fast_currentl3")
					Fast_CurrentL3Map := d.getMap(Fast_CurrentL3List, "fast_currentl3")

					Fast_ActivePower_kWList := d.getAvg(ctx, startTime, "fast_activepower_kw")
					Fast_ActivePower_kWMap := d.getMap(Fast_ActivePower_kWList, "fast_activepower_kw")

					Fast_CurrentL1List := d.getAvg(ctx, startTime, "fast_currentl1")
					Fast_CurrentL1Map := d.getMap(Fast_CurrentL1List, "fast_currentl1")

					Fast_ActivePowerSetpoint_kWList := d.getAvg(ctx, startTime, "fast_activepowersetpoint_kw")
					Fast_ActivePowerSetpoint_kWMap := d.getMap(Fast_ActivePowerSetpoint_kWList, "fast_activepowersetpoint_kw")

					Fast_CurrentL2List := d.getAvg(ctx, startTime, "fast_currentl2")
					Fast_CurrentL2Map := d.getMap(Fast_CurrentL2List, "fast_currentl2")

					Fast_DrTrVibValueList := d.getAvg(ctx, startTime, "fast_drtrvibvalue")
					Fast_DrTrVibValueMap := d.getMap(Fast_DrTrVibValueList, "fast_drtrvibvalue")

					Fast_GenSpeed_RPMList := d.getAvg(ctx, startTime, "fast_genspeed_rpm")
					Fast_GenSpeed_RPMMap := d.getMap(Fast_GenSpeed_RPMList, "fast_genspeed_rpm")

					Fast_PitchAccuV1List := d.getAvg(ctx, startTime, "fast_pitchaccuv1")
					Fast_PitchAccuV1Map := d.getMap(Fast_PitchAccuV1List, "fast_pitchaccuv1")

					Fast_PitchAngleList := d.getAvg(ctx, startTime, "fast_pitchangle")
					Fast_PitchAngleMap := d.getMap(Fast_PitchAngleList, "fast_pitchangle")

					Fast_PitchAngle3List := d.getAvg(ctx, startTime, "fast_pitchangle3")
					Fast_PitchAngle3Map := d.getMap(Fast_PitchAngle3List, "fast_pitchangle3")

					Fast_PitchAngle2List := d.getAvg(ctx, startTime, "fast_pitchangle2")
					Fast_PitchAngle2Map := d.getMap(Fast_PitchAngle2List, "fast_pitchangle2")

					Fast_PitchConvCurrent1List := d.getAvg(ctx, startTime, "fast_pitchconvcurrent1")
					Fast_PitchConvCurrent1Map := d.getMap(Fast_PitchConvCurrent1List, "fast_pitchconvcurrent1")

					Fast_PitchConvCurrent3List := d.getAvg(ctx, startTime, "fast_pitchconvcurrent3")
					Fast_PitchConvCurrent3Map := d.getMap(Fast_PitchConvCurrent3List, "fast_pitchconvcurrent3")

					Fast_PitchConvCurrent2List := d.getAvg(ctx, startTime, "fast_pitchconvcurrent2")
					Fast_PitchConvCurrent2Map := d.getMap(Fast_PitchConvCurrent2List, "fast_pitchconvcurrent2")

					Fast_PowerFactorList := d.getAvg(ctx, startTime, "fast_powerfactor")
					Fast_PowerFactorMap := d.getMap(Fast_PowerFactorList, "fast_powerfactor")

					Fast_ReactivePowerSetpointPPC_kVArList := d.getAvg(ctx, startTime, "fast_reactivepowersetpointppc_kvar")
					Fast_ReactivePowerSetpointPPC_kVArMap := d.getMap(Fast_ReactivePowerSetpointPPC_kVArList, "fast_reactivepowersetpointppc_kvar")

					Fast_ReactivePower_kVArList := d.getAvg(ctx, startTime, "fast_reactivepower_kvar")
					Fast_ReactivePower_kVArMap := d.getMap(Fast_ReactivePower_kVArList, "fast_reactivepower_kvar")

					Fast_RotorSpeed_RPMList := d.getAvg(ctx, startTime, "fast_rotorspeed_rpm")
					Fast_RotorSpeed_RPMMap := d.getMap(Fast_RotorSpeed_RPMList, "fast_rotorspeed_rpm")

					Fast_VoltageL1List := d.getAvg(ctx, startTime, "fast_voltagel1")
					Fast_VoltageL1Map := d.getMap(Fast_VoltageL1List, "fast_voltagel1")

					Fast_VoltageL2List := d.getAvg(ctx, startTime, "fast_voltagel2")
					Fast_VoltageL2Map := d.getMap(Fast_VoltageL2List, "fast_voltagel2")

					Fast_WindSpeed_msList := d.getAvg(ctx, startTime, "fast_windspeed_ms")
					Fast_WindSpeed_msMap := d.getMap(Fast_WindSpeed_msList, "fast_windspeed_ms")

					Slow_CapableCapacitiveReactPwr_kVArList := d.getAvg(ctx, startTime, "slow_capablecapacitivereactpwr_kvar")
					Slow_CapableCapacitiveReactPwr_kVArMap := d.getMap(Slow_CapableCapacitiveReactPwr_kVArList, "slow_capablecapacitivereactpwr_kvar")

					Slow_CapableInductiveReactPwr_kVArList := d.getAvg(ctx, startTime, "slow_capableinductivereactpwr_kvar")
					Slow_CapableInductiveReactPwr_kVArMap := d.getMap(Slow_CapableInductiveReactPwr_kVArList, "slow_capableinductivereactpwr_kvar")

					Slow_DateTime_SecList := d.getAvg(ctx, startTime, "slow_datetime_sec")
					Slow_DateTime_SecMap := d.getMap(Slow_DateTime_SecList, "slow_datetime_sec")

					Slow_NacellePosList := d.getAvg(ctx, startTime, "slow_nacellepos")
					Slow_NacellePosMap := d.getMap(Slow_NacellePosList, "slow_nacellepos")

					Fast_PitchAngle1List := d.getAvg(ctx, startTime, "fast_pitchangle1")
					Fast_PitchAngle1Map := d.getMap(Fast_PitchAngle1List, "fast_pitchangle1")

					Fast_VoltageL3List := d.getAvg(ctx, startTime, "fast_voltagel3")
					Fast_VoltageL3Map := d.getMap(Fast_VoltageL3List, "fast_voltagel3")

					Slow_CapableCapacitivePwrFactorList := d.getAvg(ctx, startTime, "slow_capablecapacitivepwrfactor")
					Slow_CapableCapacitivePwrFactorMap := d.getMap(Slow_CapableCapacitivePwrFactorList, "slow_capablecapacitivepwrfactor")

					Fast_Total_Production_kWhList := d.getAvg(ctx, startTime, "fast_total_production_kwh")
					Fast_Total_Production_kWhMap := d.getMap(Fast_Total_Production_kWhList, "fast_total_production_kwh")

					Fast_Total_Prod_Day_kWhList := d.getAvg(ctx, startTime, "fast_total_prod_day_kwh")
					Fast_Total_Prod_Day_kWhMap := d.getMap(Fast_Total_Prod_Day_kWhList, "fast_total_prod_day_kwh")

					Fast_Total_Prod_Month_kWhList := d.getAvg(ctx, startTime, "fast_total_prod_month_kwh")
					Fast_Total_Prod_Month_kWhMap := d.getMap(Fast_Total_Prod_Month_kWhList, "fast_total_prod_month_kwh")

					Fast_ActivePowerOutPWCSell_kWList := d.getAvg(ctx, startTime, "fast_activepoweroutpwcsell_kw")
					Fast_ActivePowerOutPWCSell_kWMap := d.getMap(Fast_ActivePowerOutPWCSell_kWList, "fast_activepoweroutpwcsell_kw")

					Fast_Frequency_HzList := d.getAvg(ctx, startTime, "fast_frequency_hz")
					Fast_Frequency_HzMap := d.getMap(Fast_Frequency_HzList, "fast_frequency_hz")

					Slow_TempG1L2List := d.getAvg(ctx, startTime, "slow_tempg1l2")
					Slow_TempG1L2Map := d.getMap(Slow_TempG1L2List, "slow_tempg1l2")

					Slow_TempG1L3List := d.getAvg(ctx, startTime, "slow_tempg1l3")
					Slow_TempG1L3Map := d.getMap(Slow_TempG1L3List, "slow_tempg1l3")

					Slow_TempGearBoxHSSDEList := d.getAvg(ctx, startTime, "slow_tempgearboxhssde")
					Slow_TempGearBoxHSSDEMap := d.getMap(Slow_TempGearBoxHSSDEList, "slow_tempgearboxhssde")

					Slow_TempGearBoxIMSNDEList := d.getAvg(ctx, startTime, "slow_tempgearboximsnde")
					Slow_TempGearBoxIMSNDEMap := d.getMap(Slow_TempGearBoxIMSNDEList, "slow_tempgearboximsnde")

					Slow_TempOutdoorList := d.getAvg(ctx, startTime, "slow_tempoutdoor")
					Slow_TempOutdoorMap := d.getMap(Slow_TempOutdoorList, "slow_tempoutdoor")

					Fast_PitchAccuV3List := d.getAvg(ctx, startTime, "fast_pitchaccuv3")
					Fast_PitchAccuV3Map := d.getMap(Fast_PitchAccuV3List, "fast_pitchaccuv3")

					Slow_TotalTurbineActiveHoursList := d.getAvg(ctx, startTime, "slow_totalturbineactivehours")
					Slow_TotalTurbineActiveHoursMap := d.getMap(Slow_TotalTurbineActiveHoursList, "slow_totalturbineactivehours")

					Slow_TotalTurbineOKHoursList := d.getAvg(ctx, startTime, "slow_totalturbineokhours")
					Slow_TotalTurbineOKHoursMap := d.getMap(Slow_TotalTurbineOKHoursList, "slow_totalturbineokhours")

					Slow_TotalTurbineTimeAllHoursList := d.getAvg(ctx, startTime, "slow_totalturbinetimeallhours")
					Slow_TotalTurbineTimeAllHoursMap := d.getMap(Slow_TotalTurbineTimeAllHoursList, "slow_totalturbinetimeallhours")

					Slow_TempG1L1List := d.getAvg(ctx, startTime, "slow_tempg1l1")
					Slow_TempG1L1Map := d.getMap(Slow_TempG1L1List, "slow_tempg1l1")

					Slow_TempGearBoxOilSumpList := d.getAvg(ctx, startTime, "slow_tempgearboxoilsump")
					Slow_TempGearBoxOilSumpMap := d.getMap(Slow_TempGearBoxOilSumpList, "slow_tempgearboxoilsump")

					Fast_PitchAccuV2List := d.getAvg(ctx, startTime, "fast_pitchaccuv2")
					Fast_PitchAccuV2Map := d.getMap(Fast_PitchAccuV2List, "fast_pitchaccuv2")

					Slow_TotalGridOkHoursList := d.getAvg(ctx, startTime, "slow_totalgridokhours")
					Slow_TotalGridOkHoursMap := d.getMap(Slow_TotalGridOkHoursList, "slow_totalgridokhours")

					Slow_TotalActPowerOut_kWhList := d.getAvg(ctx, startTime, "slow_totalactpowerout_kwh")
					Slow_TotalActPowerOut_kWhMap := d.getMap(Slow_TotalActPowerOut_kWhList, "slow_totalactpowerout_kwh")

					Fast_YawServiceList := d.getAvg(ctx, startTime, "fast_yawservice")
					Fast_YawServiceMap := d.getMap(Fast_YawServiceList, "fast_yawservice")

					Fast_YawAngleList := d.getAvg(ctx, startTime, "fast_yawangle")
					Fast_YawAngleMap := d.getMap(Fast_YawAngleList, "fast_yawangle")

					Slow_WindDirectionList := d.getAvg(ctx, startTime, "slow_winddirection")
					Slow_WindDirectionMap := d.getMap(Slow_WindDirectionList, "slow_winddirection")

					Slow_CapableInductivePwrFactorList := d.getAvg(ctx, startTime, "slow_capableinductivepwrfactor")
					Slow_CapableInductivePwrFactorMap := d.getMap(Slow_CapableInductivePwrFactorList, "slow_capableinductivepwrfactor")

					Slow_TempGearBoxHSSNDEList := d.getAvg(ctx, startTime, "slow_tempgearboxhssnde")
					Slow_TempGearBoxHSSNDEMap := d.getMap(Slow_TempGearBoxHSSNDEList, "slow_tempgearboxhssnde")

					Slow_TempHubBearingList := d.getAvg(ctx, startTime, "slow_temphubbearing")
					Slow_TempHubBearingMap := d.getMap(Slow_TempHubBearingList, "slow_temphubbearing")

					Slow_TotalG1ActiveHoursList := d.getAvg(ctx, startTime, "slow_totalg1activehours")
					Slow_TotalG1ActiveHoursMap := d.getMap(Slow_TotalG1ActiveHoursList, "slow_totalg1activehours")

					Slow_TotalActPowerOutG1_kWhList := d.getAvg(ctx, startTime, "slow_totalactpoweroutg1_kwh")
					Slow_TotalActPowerOutG1_kWhMap := d.getMap(Slow_TotalActPowerOutG1_kWhList, "slow_totalactpoweroutg1_kwh")

					Slow_TotalReactPowerInG1_kVArhList := d.getAvg(ctx, startTime, "slow_totalreactpowering1_kvarh")
					Slow_TotalReactPowerInG1_kVArhMap := d.getMap(Slow_TotalReactPowerInG1_kVArhList, "slow_totalreactpowering1_kvarh")

					Slow_NacelleDrillList := d.getAvg(ctx, startTime, "slow_nacelledrill")
					Slow_NacelleDrillMap := d.getMap(Slow_NacelleDrillList, "slow_nacelledrill")

					Slow_TempGearBoxIMSDEList := d.getAvg(ctx, startTime, "slow_tempgearboximsde")
					Slow_TempGearBoxIMSDEMap := d.getMap(Slow_TempGearBoxIMSDEList, "slow_tempgearboximsde")

					Fast_Total_Operating_hrsList := d.getAvg(ctx, startTime, "fast_total_operating_hrs")
					Fast_Total_Operating_hrsMap := d.getMap(Fast_Total_Operating_hrsList, "fast_total_operating_hrs")

					Slow_TempNacelleList := d.getAvg(ctx, startTime, "slow_tempnacelle")
					Slow_TempNacelleMap := d.getMap(Slow_TempNacelleList, "slow_tempnacelle")

					Fast_Total_Grid_OK_hrsList := d.getAvg(ctx, startTime, "fast_total_grid_ok_hrs")
					Fast_Total_Grid_OK_hrsMap := d.getMap(Fast_Total_Grid_OK_hrsList, "fast_total_grid_ok_hrs")

					Fast_Total_WTG_OK_hrsList := d.getAvg(ctx, startTime, "fast_total_wtg_ok_hrs")
					Fast_Total_WTG_OK_hrsMap := d.getMap(Fast_Total_WTG_OK_hrsList, "fast_total_wtg_ok_hrs")

					Slow_TempCabinetTopBoxList := d.getAvg(ctx, startTime, "slow_tempcabinettopbox")
					Slow_TempCabinetTopBoxMap := d.getMap(Slow_TempCabinetTopBoxList, "slow_tempcabinettopbox")

					Slow_TempGeneratorBearingNDEList := d.getAvg(ctx, startTime, "slow_tempgeneratorbearingnde")
					Slow_TempGeneratorBearingNDEMap := d.getMap(Slow_TempGeneratorBearingNDEList, "slow_tempgeneratorbearingnde")

					Fast_Total_Access_hrsList := d.getAvg(ctx, startTime, "fast_total_access_hrs")
					Fast_Total_Access_hrsMap := d.getMap(Fast_Total_Access_hrsList, "fast_total_access_hrs")

					Slow_TempBottomPowerSectionList := d.getAvg(ctx, startTime, "slow_tempbottompowersection")
					Slow_TempBottomPowerSectionMap := d.getMap(Slow_TempBottomPowerSectionList, "slow_tempbottompowersection")

					Slow_TempGeneratorBearingDEList := d.getAvg(ctx, startTime, "slow_tempgeneratorbearingde")
					Slow_TempGeneratorBearingDEMap := d.getMap(Slow_TempGeneratorBearingDEList, "slow_tempgeneratorbearingde")

					Slow_TotalReactPowerIn_kVArhList := d.getAvg(ctx, startTime, "slow_totalreactpowerin_kvarh")
					Slow_TotalReactPowerIn_kVArhMap := d.getMap(Slow_TotalReactPowerIn_kVArhList, "slow_totalreactpowerin_kvarh")

					Slow_TempBottomControlSectionList := d.getAvg(ctx, startTime, "slow_tempbottomcontrolsection")
					Slow_TempBottomControlSectionMap := d.getMap(Slow_TempBottomControlSectionList, "slow_tempbottomcontrolsection")

					Slow_TempConv1List := d.getAvg(ctx, startTime, "slow_tempconv1")
					Slow_TempConv1Map := d.getMap(Slow_TempConv1List, "slow_tempconv1")

					Fast_ActivePowerRated_kWList := d.getAvg(ctx, startTime, "fast_activepowerrated_kw")
					Fast_ActivePowerRated_kWMap := d.getMap(Fast_ActivePowerRated_kWList, "fast_activepowerrated_kw")

					Fast_NodeIPList := d.getAvg(ctx, startTime, "fast_nodeip")
					Fast_NodeIPMap := d.getMap(Fast_NodeIPList, "fast_nodeip")

					Fast_PitchSpeed1List := d.getAvg(ctx, startTime, "fast_pitchspeed1")
					Fast_PitchSpeed1Map := d.getMap(Fast_PitchSpeed1List, "fast_pitchspeed1")

					Slow_CFCardSizeList := d.getAvg(ctx, startTime, "slow_cfcardsize")
					Slow_CFCardSizeMap := d.getMap(Slow_CFCardSizeList, "slow_cfcardsize")

					Slow_CPU_NumberList := d.getAvg(ctx, startTime, "slow_cpu_number")
					Slow_CPU_NumberMap := d.getMap(Slow_CPU_NumberList, "slow_cpu_number")

					Slow_CFCardSpaceLeftList := d.getAvg(ctx, startTime, "slow_cfcardspaceleft")
					Slow_CFCardSpaceLeftMap := d.getMap(Slow_CFCardSpaceLeftList, "slow_cfcardspaceleft")

					Slow_TempBottomCapSectionList := d.getAvg(ctx, startTime, "slow_tempbottomcapsection")
					Slow_TempBottomCapSectionMap := d.getMap(Slow_TempBottomCapSectionList, "slow_tempbottomcapsection")

					Slow_RatedPowerList := d.getAvg(ctx, startTime, "slow_ratedpower")
					Slow_RatedPowerMap := d.getMap(Slow_RatedPowerList, "slow_ratedpower")

					Slow_TempConv3List := d.getAvg(ctx, startTime, "slow_tempconv3")
					Slow_TempConv3Map := d.getMap(Slow_TempConv3List, "slow_tempconv3")

					Slow_TempConv2List := d.getAvg(ctx, startTime, "slow_tempconv2")
					Slow_TempConv2Map := d.getMap(Slow_TempConv2List, "slow_tempconv2")

					Slow_TotalActPowerIn_kWhList := d.getAvg(ctx, startTime, "slow_totalactpowerin_kwh")
					Slow_TotalActPowerIn_kWhMap := d.getMap(Slow_TotalActPowerIn_kWhList, "slow_totalactpowerin_kwh")

					Slow_TotalActPowerInG1_kWhList := d.getAvg(ctx, startTime, "slow_totalactpowering1_kwh")
					Slow_TotalActPowerInG1_kWhMap := d.getMap(Slow_TotalActPowerInG1_kWhList, "slow_totalactpowering1_kwh")

					Slow_TotalActPowerInG2_kWhList := d.getAvg(ctx, startTime, "slow_totalactpowering2_kwh")
					Slow_TotalActPowerInG2_kWhMap := d.getMap(Slow_TotalActPowerInG2_kWhList, "slow_totalactpowering2_kwh")

					Slow_TotalActPowerOutG2_kWhList := d.getAvg(ctx, startTime, "slow_totalactpoweroutg2_kwh")
					Slow_TotalActPowerOutG2_kWhMap := d.getMap(Slow_TotalActPowerOutG2_kWhList, "slow_totalactpoweroutg2_kwh")

					Slow_TotalG2ActiveHoursList := d.getAvg(ctx, startTime, "slow_totalg2activehours")
					Slow_TotalG2ActiveHoursMap := d.getMap(Slow_TotalG2ActiveHoursList, "slow_totalg2activehours")

					Slow_TotalReactPowerInG2_kVArhList := d.getAvg(ctx, startTime, "slow_totalreactpowering2_kvarh")
					Slow_TotalReactPowerInG2_kVArhMap := d.getMap(Slow_TotalReactPowerInG2_kVArhList, "slow_totalreactpowering2_kvarh")

					Slow_TotalReactPowerOut_kVArhList := d.getAvg(ctx, startTime, "slow_totalreactpowerout_kvarh")
					Slow_TotalReactPowerOut_kVArhMap := d.getMap(Slow_TotalReactPowerOut_kVArhList, "slow_totalreactpowerout_kvarh")

					Slow_UTCoffset_intList := d.getAvg(ctx, startTime, "slow_utcoffset_int")
					Slow_UTCoffset_intMap := d.getMap(Slow_UTCoffset_intList, "slow_utcoffset_int")

					groupSub := tk.M{
						"_id": tk.M{
							"timestampsecondgroup": "$timestampsecondgroup",
							"projectname":          "$projectname",
							"turbine":              "$turbine",
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
						From(new(ScadaThreeSecs).TableName()).
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
								for _, val := range data {

									// log.Printf("val: %v \n", val)

									ext := new(ScadaThreeSecsExt)
									ext.File = file

									idSub := val.Get("_id").(tk.M)
									timeStampSub := idSub.Get("timestampsecondgroup").(time.Time)
									projectNameSub := idSub.GetString("projectname")
									turbineSub := idSub.GetString("turbine")

									timeStampStr := timeStampSub.UTC().Format("060102_150405")
									key := timeStampStr + "#" + projectNameSub + "#" + turbineSub

									ext.ProjectName = projectNameSub
									ext.Turbine = turbineSub

									tenMinuteInfo := GenTenMinuteInfo(timeStampSub)

									ext.THour = tenMinuteInfo.THour
									ext.TMinute = tenMinuteInfo.TMinute
									ext.TSecond = tenMinuteInfo.TSecond
									ext.TMinuteValue = tenMinuteInfo.TMinuteValue
									ext.TMinuteCategory = tenMinuteInfo.TMinuteCategory

									ext.TimeStampConverted = startTime
									ext.TimeStampConvertedInt, _ = strconv.ParseInt(ext.TimeStampConverted.Format("200601021504"), 10, 64)
									ext.TimeStampSecondGroup = timeStampSub

									ext = ext.New()

									Fast_CurrentL3 := Fast_CurrentL3Map[key]
									if Fast_CurrentL3 != nil {
										ext.Fast_CurrentL3 = Fast_CurrentL3.GetFloat64("fast_currentl3")
										ext.Fast_CurrentL3CountSecs = Fast_CurrentL3.GetInt("fast_currentl3countsecs")
									} else {
										ext.Fast_CurrentL3 = emptyValueBig
										ext.Fast_CurrentL3CountSecs = 0
									}

									Fast_ActivePower_kW := Fast_ActivePower_kWMap[key]
									if Fast_ActivePower_kW != nil {
										ext.Fast_ActivePower_kW = Fast_ActivePower_kW.GetFloat64("fast_activepower_kw")
										ext.Fast_ActivePower_kWCountSecs = Fast_ActivePower_kW.GetInt("fast_activepower_kwcountsecs")
									} else {
										ext.Fast_ActivePower_kW = emptyValueBig
										ext.Fast_ActivePower_kWCountSecs = 0
									}

									Fast_CurrentL1 := Fast_CurrentL1Map[key]
									if Fast_CurrentL1 != nil {
										ext.Fast_CurrentL1 = Fast_CurrentL1.GetFloat64("fast_currentl1")
										ext.Fast_CurrentL1CountSecs = Fast_CurrentL1.GetInt("fast_currentl1countsecs")
									} else {
										ext.Fast_CurrentL1 = emptyValueBig
										ext.Fast_CurrentL1CountSecs = 0
									}

									Fast_ActivePowerSetpoint_kW := Fast_ActivePowerSetpoint_kWMap[key]
									if Fast_ActivePowerSetpoint_kW != nil {
										ext.Fast_ActivePowerSetpoint_kW = Fast_ActivePowerSetpoint_kW.GetFloat64("fast_activepowersetpoint_kw")
										ext.Fast_ActivePowerSetpoint_kWCountSecs = Fast_ActivePowerSetpoint_kW.GetInt("fast_activepowersetpoint_kwcountsecs")
									} else {
										ext.Fast_ActivePowerSetpoint_kW = emptyValueBig
										ext.Fast_ActivePowerSetpoint_kWCountSecs = 0
									}

									Fast_CurrentL2 := Fast_CurrentL2Map[key]
									if Fast_CurrentL2 != nil {
										ext.Fast_CurrentL2 = Fast_CurrentL2.GetFloat64("fast_currentl2")
										ext.Fast_CurrentL2CountSecs = Fast_CurrentL2.GetInt("fast_currentl2countsecs")
									} else {
										ext.Fast_CurrentL2 = emptyValueBig
										ext.Fast_CurrentL2CountSecs = 0
									}

									Fast_DrTrVibValue := Fast_DrTrVibValueMap[key]
									if Fast_DrTrVibValue != nil {
										ext.Fast_DrTrVibValue = Fast_DrTrVibValue.GetFloat64("fast_drtrvibvalue")
										ext.Fast_DrTrVibValueCountSecs = Fast_DrTrVibValue.GetInt("fast_drtrvibvaluecountsecs")
									} else {
										ext.Fast_DrTrVibValue = emptyValueBig
										ext.Fast_DrTrVibValueCountSecs = 0
									}

									Fast_GenSpeed_RPM := Fast_GenSpeed_RPMMap[key]
									if Fast_GenSpeed_RPM != nil {
										ext.Fast_GenSpeed_RPM = Fast_GenSpeed_RPM.GetFloat64("fast_genspeed_rpm")
										ext.Fast_GenSpeed_RPMCountSecs = Fast_GenSpeed_RPM.GetInt("fast_genspeed_rpmcountsecs")
									} else {
										ext.Fast_GenSpeed_RPM = emptyValueBig
										ext.Fast_GenSpeed_RPMCountSecs = 0
									}

									Fast_PitchAccuV1 := Fast_PitchAccuV1Map[key]
									if Fast_PitchAccuV1 != nil {
										ext.Fast_PitchAccuV1 = Fast_PitchAccuV1.GetFloat64("fast_pitchaccuv1")
										ext.Fast_PitchAccuV1CountSecs = Fast_PitchAccuV1.GetInt("fast_pitchaccuv1countsecs")
									} else {
										ext.Fast_PitchAccuV1 = emptyValueBig
										ext.Fast_PitchAccuV1CountSecs = 0
									}

									Fast_PitchAngle := Fast_PitchAngleMap[key]
									if Fast_PitchAngle != nil {
										ext.Fast_PitchAngle = Fast_PitchAngle.GetFloat64("fast_pitchangle")
										ext.Fast_PitchAngleCountSecs = Fast_PitchAngle.GetInt("fast_pitchanglecountsecs")
									} else {
										ext.Fast_PitchAngle = emptyValueBig
										ext.Fast_PitchAngleCountSecs = 0
									}

									Fast_PitchAngle3 := Fast_PitchAngle3Map[key]
									if Fast_PitchAngle3 != nil {
										ext.Fast_PitchAngle3 = Fast_PitchAngle3.GetFloat64("fast_pitchangle3")
										ext.Fast_PitchAngle3CountSecs = Fast_PitchAngle3.GetInt("fast_pitchangle3countsecs")
									} else {
										ext.Fast_PitchAngle3 = emptyValueBig
										ext.Fast_PitchAngle3CountSecs = 0
									}

									Fast_PitchAngle2 := Fast_PitchAngle2Map[key]
									if Fast_PitchAngle2 != nil {
										ext.Fast_PitchAngle2 = Fast_PitchAngle2.GetFloat64("fast_pitchangle2")
										ext.Fast_PitchAngle2CountSecs = Fast_PitchAngle2.GetInt("fast_pitchangle2countsecs")
									} else {
										ext.Fast_PitchAngle2 = emptyValueBig
										ext.Fast_PitchAngle2CountSecs = 0
									}

									Fast_PitchConvCurrent1 := Fast_PitchConvCurrent1Map[key]
									if Fast_PitchConvCurrent1 != nil {
										ext.Fast_PitchConvCurrent1 = Fast_PitchConvCurrent1.GetFloat64("fast_pitchconvcurrent1")
										ext.Fast_PitchConvCurrent1CountSecs = Fast_PitchConvCurrent1.GetInt("fast_pitchconvcurrent1countsecs")
									} else {
										ext.Fast_PitchConvCurrent1 = emptyValueBig
										ext.Fast_PitchConvCurrent1CountSecs = 0
									}

									Fast_PitchConvCurrent3 := Fast_PitchConvCurrent3Map[key]
									if Fast_PitchConvCurrent3 != nil {
										ext.Fast_PitchConvCurrent3 = Fast_PitchConvCurrent3.GetFloat64("fast_pitchconvcurrent3")
										ext.Fast_PitchConvCurrent3CountSecs = Fast_PitchConvCurrent3.GetInt("fast_pitchconvcurrent3countsecs")
									} else {
										ext.Fast_PitchConvCurrent3 = emptyValueBig
										ext.Fast_PitchConvCurrent3CountSecs = 0
									}

									Fast_PitchConvCurrent2 := Fast_PitchConvCurrent2Map[key]
									if Fast_PitchConvCurrent2 != nil {
										ext.Fast_PitchConvCurrent2 = Fast_PitchConvCurrent2.GetFloat64("fast_pitchconvcurrent2")
										ext.Fast_PitchConvCurrent2CountSecs = Fast_PitchConvCurrent2.GetInt("fast_pitchconvcurrent2countsecs")
									} else {
										ext.Fast_PitchConvCurrent2 = emptyValueBig
										ext.Fast_PitchConvCurrent2CountSecs = 0
									}

									Fast_PowerFactor := Fast_PowerFactorMap[key]
									if Fast_PowerFactor != nil {
										ext.Fast_PowerFactor = Fast_PowerFactor.GetFloat64("fast_powerfactor")
										ext.Fast_PowerFactorCountSecs = Fast_PowerFactor.GetInt("fast_powerfactorcountsecs")
									} else {
										ext.Fast_PowerFactor = emptyValueBig
										ext.Fast_PowerFactorCountSecs = 0
									}

									Fast_ReactivePowerSetpointPPC_kVAr := Fast_ReactivePowerSetpointPPC_kVArMap[key]
									if Fast_ReactivePowerSetpointPPC_kVAr != nil {
										ext.Fast_ReactivePowerSetpointPPC_kVAr = Fast_ReactivePowerSetpointPPC_kVAr.GetFloat64("fast_reactivepowersetpointppc_kvar")
										ext.Fast_ReactivePowerSetpointPPC_kVArCountSecs = Fast_ReactivePowerSetpointPPC_kVAr.GetInt("fast_reactivepowersetpointppc_kvarcountsecs")
									} else {
										ext.Fast_ReactivePowerSetpointPPC_kVAr = emptyValueBig
										ext.Fast_ReactivePowerSetpointPPC_kVArCountSecs = 0
									}

									Fast_ReactivePower_kVAr := Fast_ReactivePower_kVArMap[key]
									if Fast_ReactivePower_kVAr != nil {
										ext.Fast_ReactivePower_kVAr = Fast_ReactivePower_kVAr.GetFloat64("fast_reactivepower_kvar")
										ext.Fast_ReactivePower_kVArCountSecs = Fast_ReactivePower_kVAr.GetInt("fast_reactivepower_kvarcountsecs")
									} else {
										ext.Fast_ReactivePower_kVAr = emptyValueBig
										ext.Fast_ReactivePower_kVArCountSecs = 0
									}

									Fast_RotorSpeed_RPM := Fast_RotorSpeed_RPMMap[key]
									if Fast_RotorSpeed_RPM != nil {
										ext.Fast_RotorSpeed_RPM = Fast_RotorSpeed_RPM.GetFloat64("fast_rotorspeed_rpm")
										ext.Fast_RotorSpeed_RPMCountSecs = Fast_RotorSpeed_RPM.GetInt("fast_rotorspeed_rpmcountsecs")
									} else {
										ext.Fast_RotorSpeed_RPM = emptyValueBig
										ext.Fast_RotorSpeed_RPMCountSecs = 0
									}

									Fast_VoltageL1 := Fast_VoltageL1Map[key]
									if Fast_VoltageL1 != nil {
										ext.Fast_VoltageL1 = Fast_VoltageL1.GetFloat64("fast_voltagel1")
										ext.Fast_VoltageL1CountSecs = Fast_VoltageL1.GetInt("fast_voltagel1countsecs")
									} else {
										ext.Fast_VoltageL1 = emptyValueBig
										ext.Fast_VoltageL1CountSecs = 0
									}

									Fast_VoltageL2 := Fast_VoltageL2Map[key]
									if Fast_VoltageL2 != nil {
										ext.Fast_VoltageL2 = Fast_VoltageL2.GetFloat64("fast_voltagel2")
										ext.Fast_VoltageL2CountSecs = Fast_VoltageL2.GetInt("fast_voltagel2countsecs")
									} else {
										ext.Fast_VoltageL2 = emptyValueBig
										ext.Fast_VoltageL2CountSecs = 0
									}

									Fast_WindSpeed_ms := Fast_WindSpeed_msMap[key]
									if Fast_WindSpeed_ms != nil {
										ext.Fast_WindSpeed_ms = Fast_WindSpeed_ms.GetFloat64("fast_windspeed_ms")
										ext.Fast_WindSpeed_msCountSecs = Fast_WindSpeed_ms.GetInt("fast_windspeed_mscountsecs")
									} else {
										ext.Fast_WindSpeed_ms = emptyValueBig
										ext.Fast_WindSpeed_msCountSecs = 0
									}

									Slow_CapableCapacitiveReactPwr_kVAr := Slow_CapableCapacitiveReactPwr_kVArMap[key]
									if Slow_CapableCapacitiveReactPwr_kVAr != nil {
										ext.Slow_CapableCapacitiveReactPwr_kVAr = Slow_CapableCapacitiveReactPwr_kVAr.GetFloat64("slow_capablecapacitivereactpwr_kvar")
										ext.Slow_CapableCapacitiveReactPwr_kVArCountSecs = Slow_CapableCapacitiveReactPwr_kVAr.GetInt("slow_capablecapacitivereactpwr_kvarcountsecs")
									} else {
										ext.Slow_CapableCapacitiveReactPwr_kVAr = emptyValueBig
										ext.Slow_CapableCapacitiveReactPwr_kVArCountSecs = 0
									}

									Slow_CapableInductiveReactPwr_kVAr := Slow_CapableInductiveReactPwr_kVArMap[key]
									if Slow_CapableInductiveReactPwr_kVAr != nil {
										ext.Slow_CapableInductiveReactPwr_kVAr = Slow_CapableInductiveReactPwr_kVAr.GetFloat64("slow_capableinductivereactpwr_kvar")
										ext.Slow_CapableInductiveReactPwr_kVArCountSecs = Slow_CapableInductiveReactPwr_kVAr.GetInt("slow_capableinductivereactpwr_kvarcountsecs")
									} else {
										ext.Slow_CapableInductiveReactPwr_kVAr = emptyValueBig
										ext.Slow_CapableInductiveReactPwr_kVArCountSecs = 0
									}

									Slow_DateTime_Sec := Slow_DateTime_SecMap[key]
									if Slow_DateTime_Sec != nil {
										ext.Slow_DateTime_Sec = Slow_DateTime_Sec.GetFloat64("slow_datetime_sec")
										ext.Slow_DateTime_SecCountSecs = Slow_DateTime_Sec.GetInt("slow_datetime_seccountsecs")
									} else {
										ext.Slow_DateTime_Sec = emptyValueBig
										ext.Slow_DateTime_SecCountSecs = 0
									}

									Slow_NacellePos := Slow_NacellePosMap[key]
									if Slow_NacellePos != nil {
										ext.Slow_NacellePos = Slow_NacellePos.GetFloat64("slow_nacellepos")
										ext.Slow_NacellePosCountSecs = Slow_NacellePos.GetInt("slow_nacelleposcountsecs")
									} else {
										ext.Slow_NacellePos = emptyValueBig
										ext.Slow_NacellePosCountSecs = 0
									}

									Fast_PitchAngle1 := Fast_PitchAngle1Map[key]
									if Fast_PitchAngle1 != nil {
										ext.Fast_PitchAngle1 = Fast_PitchAngle1.GetFloat64("fast_pitchangle1")
										ext.Fast_PitchAngle1CountSecs = Fast_PitchAngle1.GetInt("fast_pitchangle1countsecs")
									} else {
										ext.Fast_PitchAngle1 = emptyValueBig
										ext.Fast_PitchAngle1CountSecs = 0
									}

									Fast_VoltageL3 := Fast_VoltageL3Map[key]
									if Fast_VoltageL3 != nil {
										ext.Fast_VoltageL3 = Fast_VoltageL3.GetFloat64("fast_voltagel3")
										ext.Fast_VoltageL3CountSecs = Fast_VoltageL3.GetInt("fast_voltagel3countsecs")
									} else {
										ext.Fast_VoltageL3 = emptyValueBig
										ext.Fast_VoltageL3CountSecs = 0
									}

									Slow_CapableCapacitivePwrFactor := Slow_CapableCapacitivePwrFactorMap[key]
									if Slow_CapableCapacitivePwrFactor != nil {
										ext.Slow_CapableCapacitivePwrFactor = Slow_CapableCapacitivePwrFactor.GetFloat64("slow_capablecapacitivepwrfactor")
										ext.Slow_CapableCapacitivePwrFactorCountSecs = Slow_CapableCapacitivePwrFactor.GetInt("slow_capablecapacitivepwrfactorcountsecs")
									} else {
										ext.Slow_CapableCapacitivePwrFactor = emptyValueBig
										ext.Slow_CapableCapacitivePwrFactorCountSecs = 0
									}

									Fast_Total_Production_kWh := Fast_Total_Production_kWhMap[key]
									if Fast_Total_Production_kWh != nil {
										ext.Fast_Total_Production_kWh = Fast_Total_Production_kWh.GetFloat64("fast_total_production_kwh")
										ext.Fast_Total_Production_kWhCountSecs = Fast_Total_Production_kWh.GetInt("fast_total_production_kwhcountsecs")
									} else {
										ext.Fast_Total_Production_kWh = emptyValueBig
										ext.Fast_Total_Production_kWhCountSecs = 0
									}

									Fast_Total_Prod_Day_kWh := Fast_Total_Prod_Day_kWhMap[key]
									if Fast_Total_Prod_Day_kWh != nil {
										ext.Fast_Total_Prod_Day_kWh = Fast_Total_Prod_Day_kWh.GetFloat64("fast_total_prod_day_kwh")
										ext.Fast_Total_Prod_Day_kWhCountSecs = Fast_Total_Prod_Day_kWh.GetInt("fast_total_prod_day_kwhcountsecs")
									} else {
										ext.Fast_Total_Prod_Day_kWh = emptyValueBig
										ext.Fast_Total_Prod_Day_kWhCountSecs = 0
									}

									Fast_Total_Prod_Month_kWh := Fast_Total_Prod_Month_kWhMap[key]
									if Fast_Total_Prod_Month_kWh != nil {
										ext.Fast_Total_Prod_Month_kWh = Fast_Total_Prod_Month_kWh.GetFloat64("fast_total_prod_month_kwh")
										ext.Fast_Total_Prod_Month_kWhCountSecs = Fast_Total_Prod_Month_kWh.GetInt("fast_total_prod_month_kwhcountsecs")
									} else {
										ext.Fast_Total_Prod_Month_kWh = emptyValueBig
										ext.Fast_Total_Prod_Month_kWhCountSecs = 0
									}

									Fast_ActivePowerOutPWCSell_kW := Fast_ActivePowerOutPWCSell_kWMap[key]
									if Fast_ActivePowerOutPWCSell_kW != nil {
										ext.Fast_ActivePowerOutPWCSell_kW = Fast_ActivePowerOutPWCSell_kW.GetFloat64("fast_activepoweroutpwcsell_kw")
										ext.Fast_ActivePowerOutPWCSell_kWCountSecs = Fast_ActivePowerOutPWCSell_kW.GetInt("fast_activepoweroutpwcsell_kwcountsecs")
									} else {
										ext.Fast_ActivePowerOutPWCSell_kW = emptyValueBig
										ext.Fast_ActivePowerOutPWCSell_kWCountSecs = 0
									}

									Fast_Frequency_Hz := Fast_Frequency_HzMap[key]
									if Fast_Frequency_Hz != nil {
										ext.Fast_Frequency_Hz = Fast_Frequency_Hz.GetFloat64("fast_frequency_hz")
										ext.Fast_Frequency_HzCountSecs = Fast_Frequency_Hz.GetInt("fast_frequency_hzcountsecs")
									} else {
										ext.Fast_Frequency_Hz = emptyValueBig
										ext.Fast_Frequency_HzCountSecs = 0
									}

									Slow_TempG1L2 := Slow_TempG1L2Map[key]
									if Slow_TempG1L2 != nil {
										ext.Slow_TempG1L2 = Slow_TempG1L2.GetFloat64("slow_tempg1l2")
										ext.Slow_TempG1L2CountSecs = Slow_TempG1L2.GetInt("slow_tempg1l2countsecs")
									} else {
										ext.Slow_TempG1L2 = emptyValueBig
										ext.Slow_TempG1L2CountSecs = 0
									}

									Slow_TempG1L3 := Slow_TempG1L3Map[key]
									if Slow_TempG1L3 != nil {
										ext.Slow_TempG1L3 = Slow_TempG1L3.GetFloat64("slow_tempg1l3")
										ext.Slow_TempG1L3CountSecs = Slow_TempG1L3.GetInt("slow_tempg1l3countsecs")
									} else {
										ext.Slow_TempG1L3 = emptyValueBig
										ext.Slow_TempG1L3CountSecs = 0
									}

									Slow_TempGearBoxHSSDE := Slow_TempGearBoxHSSDEMap[key]
									if Slow_TempGearBoxHSSDE != nil {
										ext.Slow_TempGearBoxHSSDE = Slow_TempGearBoxHSSDE.GetFloat64("slow_tempgearboxhssde")
										ext.Slow_TempGearBoxHSSDECountSecs = Slow_TempGearBoxHSSDE.GetInt("slow_tempgearboxhssdecountsecs")
									} else {
										ext.Slow_TempGearBoxHSSDE = emptyValueBig
										ext.Slow_TempGearBoxHSSDECountSecs = 0
									}

									Slow_TempGearBoxIMSNDE := Slow_TempGearBoxIMSNDEMap[key]
									if Slow_TempGearBoxIMSNDE != nil {
										ext.Slow_TempGearBoxIMSNDE = Slow_TempGearBoxIMSNDE.GetFloat64("slow_tempgearboximsnde")
										ext.Slow_TempGearBoxIMSNDECountSecs = Slow_TempGearBoxIMSNDE.GetInt("slow_tempgearboximsndecountsecs")
									} else {
										ext.Slow_TempGearBoxIMSNDE = emptyValueBig
										ext.Slow_TempGearBoxIMSNDECountSecs = 0
									}

									Slow_TempOutdoor := Slow_TempOutdoorMap[key]
									if Slow_TempOutdoor != nil {
										ext.Slow_TempOutdoor = Slow_TempOutdoor.GetFloat64("slow_tempoutdoor")
										ext.Slow_TempOutdoorCountSecs = Slow_TempOutdoor.GetInt("slow_tempoutdoorcountsecs")
									} else {
										ext.Slow_TempOutdoor = emptyValueBig
										ext.Slow_TempOutdoorCountSecs = 0
									}

									Fast_PitchAccuV3 := Fast_PitchAccuV3Map[key]
									if Fast_PitchAccuV3 != nil {
										ext.Fast_PitchAccuV3 = Fast_PitchAccuV3.GetFloat64("fast_pitchaccuv3")
										ext.Fast_PitchAccuV3CountSecs = Fast_PitchAccuV3.GetInt("fast_pitchaccuv3countsecs")
									} else {
										ext.Fast_PitchAccuV3 = emptyValueBig
										ext.Fast_PitchAccuV3CountSecs = 0
									}

									Slow_TotalTurbineActiveHours := Slow_TotalTurbineActiveHoursMap[key]
									if Slow_TotalTurbineActiveHours != nil {
										ext.Slow_TotalTurbineActiveHours = Slow_TotalTurbineActiveHours.GetFloat64("slow_totalturbineactivehours")
										ext.Slow_TotalTurbineActiveHoursCountSecs = Slow_TotalTurbineActiveHours.GetInt("slow_totalturbineactivehourscountsecs")
									} else {
										ext.Slow_TotalTurbineActiveHours = emptyValueBig
										ext.Slow_TotalTurbineActiveHoursCountSecs = 0
									}

									Slow_TotalTurbineOKHours := Slow_TotalTurbineOKHoursMap[key]
									if Slow_TotalTurbineOKHours != nil {
										ext.Slow_TotalTurbineOKHours = Slow_TotalTurbineOKHours.GetFloat64("slow_totalturbineokhours")
										ext.Slow_TotalTurbineOKHoursCountSecs = Slow_TotalTurbineOKHours.GetInt("slow_totalturbineokhourscountsecs")
									} else {
										ext.Slow_TotalTurbineOKHours = emptyValueBig
										ext.Slow_TotalTurbineOKHoursCountSecs = 0
									}

									Slow_TotalTurbineTimeAllHours := Slow_TotalTurbineTimeAllHoursMap[key]
									if Slow_TotalTurbineTimeAllHours != nil {
										ext.Slow_TotalTurbineTimeAllHours = Slow_TotalTurbineTimeAllHours.GetFloat64("slow_totalturbinetimeallhours")
										ext.Slow_TotalTurbineTimeAllHoursCountSecs = Slow_TotalTurbineTimeAllHours.GetInt("slow_totalturbinetimeallhourscountsecs")
									} else {
										ext.Slow_TotalTurbineTimeAllHours = emptyValueBig
										ext.Slow_TotalTurbineTimeAllHoursCountSecs = 0
									}

									Slow_TempG1L1 := Slow_TempG1L1Map[key]
									if Slow_TempG1L1 != nil {
										ext.Slow_TempG1L1 = Slow_TempG1L1.GetFloat64("slow_tempg1l1")
										ext.Slow_TempG1L1CountSecs = Slow_TempG1L1.GetInt("slow_tempg1l1countsecs")
									} else {
										ext.Slow_TempG1L1 = emptyValueBig
										ext.Slow_TempG1L1CountSecs = 0
									}

									Slow_TempGearBoxOilSump := Slow_TempGearBoxOilSumpMap[key]
									if Slow_TempGearBoxOilSump != nil {
										ext.Slow_TempGearBoxOilSump = Slow_TempGearBoxOilSump.GetFloat64("slow_tempgearboxoilsump")
										ext.Slow_TempGearBoxOilSumpCountSecs = Slow_TempGearBoxOilSump.GetInt("slow_tempgearboxoilsumpcountsecs")
									} else {
										ext.Slow_TempGearBoxOilSump = emptyValueBig
										ext.Slow_TempGearBoxOilSumpCountSecs = 0
									}

									Fast_PitchAccuV2 := Fast_PitchAccuV2Map[key]
									if Fast_PitchAccuV2 != nil {
										ext.Fast_PitchAccuV2 = Fast_PitchAccuV2.GetFloat64("fast_pitchaccuv2")
										ext.Fast_PitchAccuV2CountSecs = Fast_PitchAccuV2.GetInt("fast_pitchaccuv2countsecs")
									} else {
										ext.Fast_PitchAccuV2 = emptyValueBig
										ext.Fast_PitchAccuV2CountSecs = 0
									}

									Slow_TotalGridOkHours := Slow_TotalGridOkHoursMap[key]
									if Slow_TotalGridOkHours != nil {
										ext.Slow_TotalGridOkHours = Slow_TotalGridOkHours.GetFloat64("slow_totalgridokhours")
										ext.Slow_TotalGridOkHoursCountSecs = Slow_TotalGridOkHours.GetInt("slow_totalgridokhourscountsecs")
									} else {
										ext.Slow_TotalGridOkHours = emptyValueBig
										ext.Slow_TotalGridOkHoursCountSecs = 0
									}

									Slow_TotalActPowerOut_kWh := Slow_TotalActPowerOut_kWhMap[key]
									if Slow_TotalActPowerOut_kWh != nil {
										ext.Slow_TotalActPowerOut_kWh = Slow_TotalActPowerOut_kWh.GetFloat64("slow_totalactpowerout_kwh")
										ext.Slow_TotalActPowerOut_kWhCountSecs = Slow_TotalActPowerOut_kWh.GetInt("slow_totalactpowerout_kwhcountsecs")
									} else {
										ext.Slow_TotalActPowerOut_kWh = emptyValueBig
										ext.Slow_TotalActPowerOut_kWhCountSecs = 0
									}

									Fast_YawService := Fast_YawServiceMap[key]
									if Fast_YawService != nil {
										ext.Fast_YawService = Fast_YawService.GetFloat64("fast_yawservice")
										ext.Fast_YawServiceCountSecs = Fast_YawService.GetInt("fast_yawservicecountsecs")
									} else {
										ext.Fast_YawService = emptyValueBig
										ext.Fast_YawServiceCountSecs = 0
									}

									Fast_YawAngle := Fast_YawAngleMap[key]
									if Fast_YawAngle != nil {
										ext.Fast_YawAngle = Fast_YawAngle.GetFloat64("fast_yawangle")
										ext.Fast_YawAngleCountSecs = Fast_YawAngle.GetInt("fast_yawanglecountsecs")
									} else {
										ext.Fast_YawAngle = emptyValueBig
										ext.Fast_YawAngleCountSecs = 0
									}

									Slow_WindDirection := Slow_WindDirectionMap[key]
									if Slow_WindDirection != nil {
										ext.Slow_WindDirection = Slow_WindDirection.GetFloat64("slow_winddirection")
										ext.Slow_WindDirectionCountSecs = Slow_WindDirection.GetInt("slow_winddirectioncountsecs")
									} else {
										ext.Slow_WindDirection = emptyValueBig
										ext.Slow_WindDirectionCountSecs = 0
									}

									Slow_CapableInductivePwrFactor := Slow_CapableInductivePwrFactorMap[key]
									if Slow_CapableInductivePwrFactor != nil {
										ext.Slow_CapableInductivePwrFactor = Slow_CapableInductivePwrFactor.GetFloat64("slow_capableinductivepwrfactor")
										ext.Slow_CapableInductivePwrFactorCountSecs = Slow_CapableInductivePwrFactor.GetInt("slow_capableinductivepwrfactorcountsecs")
									} else {
										ext.Slow_CapableInductivePwrFactor = emptyValueBig
										ext.Slow_CapableInductivePwrFactorCountSecs = 0
									}

									Slow_TempGearBoxHSSNDE := Slow_TempGearBoxHSSNDEMap[key]
									if Slow_TempGearBoxHSSNDE != nil {
										ext.Slow_TempGearBoxHSSNDE = Slow_TempGearBoxHSSNDE.GetFloat64("slow_tempgearboxhssnde")
										ext.Slow_TempGearBoxHSSNDECountSecs = Slow_TempGearBoxHSSNDE.GetInt("slow_tempgearboxhssndecountsecs")
									} else {
										ext.Slow_TempGearBoxHSSNDE = emptyValueBig
										ext.Slow_TempGearBoxHSSNDECountSecs = 0
									}

									Slow_TempHubBearing := Slow_TempHubBearingMap[key]
									if Slow_TempHubBearing != nil {
										ext.Slow_TempHubBearing = Slow_TempHubBearing.GetFloat64("slow_temphubbearing")
										ext.Slow_TempHubBearingCountSecs = Slow_TempHubBearing.GetInt("slow_temphubbearingcountsecs")
									} else {
										ext.Slow_TempHubBearing = emptyValueBig
										ext.Slow_TempHubBearingCountSecs = 0
									}

									Slow_TotalG1ActiveHours := Slow_TotalG1ActiveHoursMap[key]
									if Slow_TotalG1ActiveHours != nil {
										ext.Slow_TotalG1ActiveHours = Slow_TotalG1ActiveHours.GetFloat64("slow_totalg1activehours")
										ext.Slow_TotalG1ActiveHoursCountSecs = Slow_TotalG1ActiveHours.GetInt("slow_totalg1activehourscountsecs")
									} else {
										ext.Slow_TotalG1ActiveHours = emptyValueBig
										ext.Slow_TotalG1ActiveHoursCountSecs = 0
									}

									Slow_TotalActPowerOutG1_kWh := Slow_TotalActPowerOutG1_kWhMap[key]
									if Slow_TotalActPowerOutG1_kWh != nil {
										ext.Slow_TotalActPowerOutG1_kWh = Slow_TotalActPowerOutG1_kWh.GetFloat64("slow_totalactpoweroutg1_kwh")
										ext.Slow_TotalActPowerOutG1_kWhCountSecs = Slow_TotalActPowerOutG1_kWh.GetInt("slow_totalactpoweroutg1_kwhcountsecs")
									} else {
										ext.Slow_TotalActPowerOutG1_kWh = emptyValueBig
										ext.Slow_TotalActPowerOutG1_kWhCountSecs = 0
									}

									Slow_TotalReactPowerInG1_kVArh := Slow_TotalReactPowerInG1_kVArhMap[key]
									if Slow_TotalReactPowerInG1_kVArh != nil {
										ext.Slow_TotalReactPowerInG1_kVArh = Slow_TotalReactPowerInG1_kVArh.GetFloat64("slow_totalreactpowering1_kvarh")
										ext.Slow_TotalReactPowerInG1_kVArhCountSecs = Slow_TotalReactPowerInG1_kVArh.GetInt("slow_totalreactpowering1_kvarhcountsecs")
									} else {
										ext.Slow_TotalReactPowerInG1_kVArh = emptyValueBig
										ext.Slow_TotalReactPowerInG1_kVArhCountSecs = 0
									}

									Slow_NacelleDrill := Slow_NacelleDrillMap[key]
									if Slow_NacelleDrill != nil {
										ext.Slow_NacelleDrill = Slow_NacelleDrill.GetFloat64("slow_nacelledrill")
										ext.Slow_NacelleDrillCountSecs = Slow_NacelleDrill.GetInt("slow_nacelledrillcountsecs")
									} else {
										ext.Slow_NacelleDrill = emptyValueBig
										ext.Slow_NacelleDrillCountSecs = 0
									}

									Slow_TempGearBoxIMSDE := Slow_TempGearBoxIMSDEMap[key]
									if Slow_TempGearBoxIMSDE != nil {
										ext.Slow_TempGearBoxIMSDE = Slow_TempGearBoxIMSDE.GetFloat64("slow_tempgearboximsde")
										ext.Slow_TempGearBoxIMSDECountSecs = Slow_TempGearBoxIMSDE.GetInt("slow_tempgearboximsdecountsecs")
									} else {
										ext.Slow_TempGearBoxIMSDE = emptyValueBig
										ext.Slow_TempGearBoxIMSDECountSecs = 0
									}

									Fast_Total_Operating_hrs := Fast_Total_Operating_hrsMap[key]
									if Fast_Total_Operating_hrs != nil {
										ext.Fast_Total_Operating_hrs = Fast_Total_Operating_hrs.GetFloat64("fast_total_operating_hrs")
										ext.Fast_Total_Operating_hrsCountSecs = Fast_Total_Operating_hrs.GetInt("fast_total_operating_hrscountsecs")
									} else {
										ext.Fast_Total_Operating_hrs = emptyValueBig
										ext.Fast_Total_Operating_hrsCountSecs = 0
									}

									Slow_TempNacelle := Slow_TempNacelleMap[key]
									if Slow_TempNacelle != nil {
										ext.Slow_TempNacelle = Slow_TempNacelle.GetFloat64("slow_tempnacelle")
										ext.Slow_TempNacelleCountSecs = Slow_TempNacelle.GetInt("slow_tempnacellecountsecs")
									} else {
										ext.Slow_TempNacelle = emptyValueBig
										ext.Slow_TempNacelleCountSecs = 0
									}

									Fast_Total_Grid_OK_hrs := Fast_Total_Grid_OK_hrsMap[key]
									if Fast_Total_Grid_OK_hrs != nil {
										ext.Fast_Total_Grid_OK_hrs = Fast_Total_Grid_OK_hrs.GetFloat64("fast_total_grid_ok_hrs")
										ext.Fast_Total_Grid_OK_hrsCountSecs = Fast_Total_Grid_OK_hrs.GetInt("fast_total_grid_ok_hrscountsecs")
									} else {
										ext.Fast_Total_Grid_OK_hrs = emptyValueBig
										ext.Fast_Total_Grid_OK_hrsCountSecs = 0
									}

									Fast_Total_WTG_OK_hrs := Fast_Total_WTG_OK_hrsMap[key]
									if Fast_Total_WTG_OK_hrs != nil {
										ext.Fast_Total_WTG_OK_hrs = Fast_Total_WTG_OK_hrs.GetFloat64("fast_total_wtg_ok_hrs")
										ext.Fast_Total_WTG_OK_hrsCountSecs = Fast_Total_WTG_OK_hrs.GetInt("fast_total_wtg_ok_hrscountsecs")
									} else {
										ext.Fast_Total_WTG_OK_hrs = emptyValueBig
										ext.Fast_Total_WTG_OK_hrsCountSecs = 0
									}

									Slow_TempCabinetTopBox := Slow_TempCabinetTopBoxMap[key]
									if Slow_TempCabinetTopBox != nil {
										ext.Slow_TempCabinetTopBox = Slow_TempCabinetTopBox.GetFloat64("slow_tempcabinettopbox")
										ext.Slow_TempCabinetTopBoxCountSecs = Slow_TempCabinetTopBox.GetInt("slow_tempcabinettopboxcountsecs")
									} else {
										ext.Slow_TempCabinetTopBox = emptyValueBig
										ext.Slow_TempCabinetTopBoxCountSecs = 0
									}

									Slow_TempGeneratorBearingNDE := Slow_TempGeneratorBearingNDEMap[key]
									if Slow_TempGeneratorBearingNDE != nil {
										ext.Slow_TempGeneratorBearingNDE = Slow_TempGeneratorBearingNDE.GetFloat64("slow_tempgeneratorbearingnde")
										ext.Slow_TempGeneratorBearingNDECountSecs = Slow_TempGeneratorBearingNDE.GetInt("slow_tempgeneratorbearingndecountsecs")
									} else {
										ext.Slow_TempGeneratorBearingNDE = emptyValueBig
										ext.Slow_TempGeneratorBearingNDECountSecs = 0
									}

									Fast_Total_Access_hrs := Fast_Total_Access_hrsMap[key]
									if Fast_Total_Access_hrs != nil {
										ext.Fast_Total_Access_hrs = Fast_Total_Access_hrs.GetFloat64("fast_total_access_hrs")
										ext.Fast_Total_Access_hrsCountSecs = Fast_Total_Access_hrs.GetInt("fast_total_access_hrscountsecs")
									} else {
										ext.Fast_Total_Access_hrs = emptyValueBig
										ext.Fast_Total_Access_hrsCountSecs = 0
									}

									Slow_TempBottomPowerSection := Slow_TempBottomPowerSectionMap[key]
									if Slow_TempBottomPowerSection != nil {
										ext.Slow_TempBottomPowerSection = Slow_TempBottomPowerSection.GetFloat64("slow_tempbottompowersection")
										ext.Slow_TempBottomPowerSectionCountSecs = Slow_TempBottomPowerSection.GetInt("slow_tempbottompowersectioncountsecs")
									} else {
										ext.Slow_TempBottomPowerSection = emptyValueBig
										ext.Slow_TempBottomPowerSectionCountSecs = 0
									}

									Slow_TempGeneratorBearingDE := Slow_TempGeneratorBearingDEMap[key]
									if Slow_TempGeneratorBearingDE != nil {
										ext.Slow_TempGeneratorBearingDE = Slow_TempGeneratorBearingDE.GetFloat64("slow_tempgeneratorbearingde")
										ext.Slow_TempGeneratorBearingDECountSecs = Slow_TempGeneratorBearingDE.GetInt("slow_tempgeneratorbearingdecountsecs")
									} else {
										ext.Slow_TempGeneratorBearingDE = emptyValueBig
										ext.Slow_TempGeneratorBearingDECountSecs = 0
									}

									Slow_TotalReactPowerIn_kVArh := Slow_TotalReactPowerIn_kVArhMap[key]
									if Slow_TotalReactPowerIn_kVArh != nil {
										ext.Slow_TotalReactPowerIn_kVArh = Slow_TotalReactPowerIn_kVArh.GetFloat64("slow_totalreactpowerin_kvarh")
										ext.Slow_TotalReactPowerIn_kVArhCountSecs = Slow_TotalReactPowerIn_kVArh.GetInt("slow_totalreactpowerin_kvarhcountsecs")
									} else {
										ext.Slow_TotalReactPowerIn_kVArh = emptyValueBig
										ext.Slow_TotalReactPowerIn_kVArhCountSecs = 0
									}

									Slow_TempBottomControlSection := Slow_TempBottomControlSectionMap[key]
									if Slow_TempBottomControlSection != nil {
										ext.Slow_TempBottomControlSection = Slow_TempBottomControlSection.GetFloat64("slow_tempbottomcontrolsection")
										ext.Slow_TempBottomControlSectionCountSecs = Slow_TempBottomControlSection.GetInt("slow_tempbottomcontrolsectioncountsecs")
									} else {
										ext.Slow_TempBottomControlSection = emptyValueBig
										ext.Slow_TempBottomControlSectionCountSecs = 0
									}

									Slow_TempConv1 := Slow_TempConv1Map[key]
									if Slow_TempConv1 != nil {
										ext.Slow_TempConv1 = Slow_TempConv1.GetFloat64("slow_tempconv1")
										ext.Slow_TempConv1CountSecs = Slow_TempConv1.GetInt("slow_tempconv1countsecs")
									} else {
										ext.Slow_TempConv1 = emptyValueBig
										ext.Slow_TempConv1CountSecs = 0
									}

									Fast_ActivePowerRated_kW := Fast_ActivePowerRated_kWMap[key]
									if Fast_ActivePowerRated_kW != nil {
										ext.Fast_ActivePowerRated_kW = Fast_ActivePowerRated_kW.GetFloat64("fast_activepowerrated_kw")
										ext.Fast_ActivePowerRated_kWCountSecs = Fast_ActivePowerRated_kW.GetInt("fast_activepowerrated_kwcountsecs")
									} else {
										ext.Fast_ActivePowerRated_kW = emptyValueBig
										ext.Fast_ActivePowerRated_kWCountSecs = 0
									}

									Fast_NodeIP := Fast_NodeIPMap[key]
									if Fast_NodeIP != nil {
										ext.Fast_NodeIP = Fast_NodeIP.GetFloat64("fast_nodeip")
										ext.Fast_NodeIPCountSecs = Fast_NodeIP.GetInt("fast_nodeipcountsecs")
									} else {
										ext.Fast_NodeIP = emptyValueBig
										ext.Fast_NodeIPCountSecs = 0
									}

									Fast_PitchSpeed1 := Fast_PitchSpeed1Map[key]
									if Fast_PitchSpeed1 != nil {
										ext.Fast_PitchSpeed1 = Fast_PitchSpeed1.GetFloat64("fast_pitchspeed1")
										ext.Fast_PitchSpeed1CountSecs = Fast_PitchSpeed1.GetInt("fast_pitchspeed1countsecs")
									} else {
										ext.Fast_PitchSpeed1 = emptyValueBig
										ext.Fast_PitchSpeed1CountSecs = 0
									}

									Slow_CFCardSize := Slow_CFCardSizeMap[key]
									if Slow_CFCardSize != nil {
										ext.Slow_CFCardSize = Slow_CFCardSize.GetFloat64("slow_cfcardsize")
										ext.Slow_CFCardSizeCountSecs = Slow_CFCardSize.GetInt("slow_cfcardsizecountsecs")
									} else {
										ext.Slow_CFCardSize = emptyValueBig
										ext.Slow_CFCardSizeCountSecs = 0
									}

									Slow_CPU_Number := Slow_CPU_NumberMap[key]
									if Slow_CPU_Number != nil {
										ext.Slow_CPU_Number = Slow_CPU_Number.GetFloat64("slow_cpu_number")
										ext.Slow_CPU_NumberCountSecs = Slow_CPU_Number.GetInt("slow_cpu_numbercountsecs")
									} else {
										ext.Slow_CPU_Number = emptyValueBig
										ext.Slow_CPU_NumberCountSecs = 0
									}

									Slow_CFCardSpaceLeft := Slow_CFCardSpaceLeftMap[key]
									if Slow_CFCardSpaceLeft != nil {
										ext.Slow_CFCardSpaceLeft = Slow_CFCardSpaceLeft.GetFloat64("slow_cfcardspaceleft")
										ext.Slow_CFCardSpaceLeftCountSecs = Slow_CFCardSpaceLeft.GetInt("slow_cfcardspaceleftcountsecs")
									} else {
										ext.Slow_CFCardSpaceLeft = emptyValueBig
										ext.Slow_CFCardSpaceLeftCountSecs = 0
									}

									Slow_TempBottomCapSection := Slow_TempBottomCapSectionMap[key]
									if Slow_TempBottomCapSection != nil {
										ext.Slow_TempBottomCapSection = Slow_TempBottomCapSection.GetFloat64("slow_tempbottomcapsection")
										ext.Slow_TempBottomCapSectionCountSecs = Slow_TempBottomCapSection.GetInt("slow_tempbottomcapsectioncountsecs")
									} else {
										ext.Slow_TempBottomCapSection = emptyValueBig
										ext.Slow_TempBottomCapSectionCountSecs = 0
									}

									Slow_RatedPower := Slow_RatedPowerMap[key]
									if Slow_RatedPower != nil {
										ext.Slow_RatedPower = Slow_RatedPower.GetFloat64("slow_ratedpower")
										ext.Slow_RatedPowerCountSecs = Slow_RatedPower.GetInt("slow_ratedpowercountsecs")
									} else {
										ext.Slow_RatedPower = emptyValueBig
										ext.Slow_RatedPowerCountSecs = 0
									}

									Slow_TempConv3 := Slow_TempConv3Map[key]
									if Slow_TempConv3 != nil {
										ext.Slow_TempConv3 = Slow_TempConv3.GetFloat64("slow_tempconv3")
										ext.Slow_TempConv3CountSecs = Slow_TempConv3.GetInt("slow_tempconv3countsecs")
									} else {
										ext.Slow_TempConv3 = emptyValueBig
										ext.Slow_TempConv3CountSecs = 0
									}

									Slow_TempConv2 := Slow_TempConv2Map[key]
									if Slow_TempConv2 != nil {
										ext.Slow_TempConv2 = Slow_TempConv2.GetFloat64("slow_tempconv2")
										ext.Slow_TempConv2CountSecs = Slow_TempConv2.GetInt("slow_tempconv2countsecs")
									} else {
										ext.Slow_TempConv2 = emptyValueBig
										ext.Slow_TempConv2CountSecs = 0
									}

									Slow_TotalActPowerIn_kWh := Slow_TotalActPowerIn_kWhMap[key]
									if Slow_TotalActPowerIn_kWh != nil {
										ext.Slow_TotalActPowerIn_kWh = Slow_TotalActPowerIn_kWh.GetFloat64("slow_totalactpowerin_kwh")
										ext.Slow_TotalActPowerIn_kWhCountSecs = Slow_TotalActPowerIn_kWh.GetInt("slow_totalactpowerin_kwhcountsecs")
									} else {
										ext.Slow_TotalActPowerIn_kWh = emptyValueBig
										ext.Slow_TotalActPowerIn_kWhCountSecs = 0
									}

									Slow_TotalActPowerInG1_kWh := Slow_TotalActPowerInG1_kWhMap[key]
									if Slow_TotalActPowerInG1_kWh != nil {
										ext.Slow_TotalActPowerInG1_kWh = Slow_TotalActPowerInG1_kWh.GetFloat64("slow_totalactpowering1_kwh")
										ext.Slow_TotalActPowerInG1_kWhCountSecs = Slow_TotalActPowerInG1_kWh.GetInt("slow_totalactpowering1_kwhcountsecs")
									} else {
										ext.Slow_TotalActPowerInG1_kWh = emptyValueBig
										ext.Slow_TotalActPowerInG1_kWhCountSecs = 0
									}

									Slow_TotalActPowerInG2_kWh := Slow_TotalActPowerInG2_kWhMap[key]
									if Slow_TotalActPowerInG2_kWh != nil {
										ext.Slow_TotalActPowerInG2_kWh = Slow_TotalActPowerInG2_kWh.GetFloat64("slow_totalactpowering2_kwh")
										ext.Slow_TotalActPowerInG2_kWhCountSecs = Slow_TotalActPowerInG2_kWh.GetInt("slow_totalactpowering2_kwhcountsecs")
									} else {
										ext.Slow_TotalActPowerInG2_kWh = emptyValueBig
										ext.Slow_TotalActPowerInG2_kWhCountSecs = 0
									}

									Slow_TotalActPowerOutG2_kWh := Slow_TotalActPowerOutG2_kWhMap[key]
									if Slow_TotalActPowerOutG2_kWh != nil {
										ext.Slow_TotalActPowerOutG2_kWh = Slow_TotalActPowerOutG2_kWh.GetFloat64("slow_totalactpoweroutg2_kwh")
										ext.Slow_TotalActPowerOutG2_kWhCountSecs = Slow_TotalActPowerOutG2_kWh.GetInt("slow_totalactpoweroutg2_kwhcountsecs")
									} else {
										ext.Slow_TotalActPowerOutG2_kWh = emptyValueBig
										ext.Slow_TotalActPowerOutG2_kWhCountSecs = 0
									}

									Slow_TotalG2ActiveHours := Slow_TotalG2ActiveHoursMap[key]
									if Slow_TotalG2ActiveHours != nil {
										ext.Slow_TotalG2ActiveHours = Slow_TotalG2ActiveHours.GetFloat64("slow_totalg2activehours")
										ext.Slow_TotalG2ActiveHoursCountSecs = Slow_TotalG2ActiveHours.GetInt("slow_totalg2activehourscountsecs")
									} else {
										ext.Slow_TotalG2ActiveHours = emptyValueBig
										ext.Slow_TotalG2ActiveHoursCountSecs = 0
									}

									Slow_TotalReactPowerInG2_kVArh := Slow_TotalReactPowerInG2_kVArhMap[key]
									if Slow_TotalReactPowerInG2_kVArh != nil {
										ext.Slow_TotalReactPowerInG2_kVArh = Slow_TotalReactPowerInG2_kVArh.GetFloat64("slow_totalreactpowering2_kvarh")
										ext.Slow_TotalReactPowerInG2_kVArhCountSecs = Slow_TotalReactPowerInG2_kVArh.GetInt("slow_totalreactpowering2_kvarhcountsecs")
									} else {
										ext.Slow_TotalReactPowerInG2_kVArh = emptyValueBig
										ext.Slow_TotalReactPowerInG2_kVArhCountSecs = 0
									}

									Slow_TotalReactPowerOut_kVArh := Slow_TotalReactPowerOut_kVArhMap[key]
									if Slow_TotalReactPowerOut_kVArh != nil {
										ext.Slow_TotalReactPowerOut_kVArh = Slow_TotalReactPowerOut_kVArh.GetFloat64("slow_totalreactpowerout_kvarh")
										ext.Slow_TotalReactPowerOut_kVArhCountSecs = Slow_TotalReactPowerOut_kVArh.GetInt("slow_totalreactpowerout_kvarhcountsecs")
									} else {
										ext.Slow_TotalReactPowerOut_kVArh = emptyValueBig
										ext.Slow_TotalReactPowerOut_kVArhCountSecs = 0
									}

									Slow_UTCoffset_int := Slow_UTCoffset_intMap[key]
									if Slow_UTCoffset_int != nil {
										ext.Slow_UTCoffset_int = Slow_UTCoffset_int.GetFloat64("slow_utcoffset_int")
										ext.Slow_UTCoffset_intCountSecs = Slow_UTCoffset_int.GetInt("slow_utcoffset_intcountsecs")
									} else {
										ext.Slow_UTCoffset_int = emptyValueBig
										ext.Slow_UTCoffset_intCountSecs = 0
									}

									// log.Printf("%#v \n", tenScada)
									mutexX.Lock()

									/*if ext.Turbine == "HBR004" {
										log.Printf("tenScada: %v | %v | %v | %v \n", ext.ID, ext.TimeStamp.UTC().Format("20060102 15:04"), startTime.Format("20060102 15:04"), idSub.Get("timestampint").(int64))
									}*/

									err := ctx.Insert(ext)
									ErrorHandler(err, "Saving")
									mutexX.Unlock()
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

func (d *ConvThreeExt) getAvg(ctx *DataContext, timestampconverted time.Time, field string) (result []tk.M) {
	pipes := []tk.M{}

	match := tk.M{
		"timestampconverted": timestampconverted,
		field:                tk.M{"$gt": emptyValueBig},
	}

	group := tk.M{
		"_id": tk.M{
			"timestamp":   "$timestampsecondgroup",
			"projectname": "$projectname",
			"turbine":     "$turbine",
		},
		field: tk.M{"$avg": "$" + field},
	}

	pipes = append(pipes, tk.M{"$match": match})
	pipes = append(pipes, tk.M{"$group": group})

	csr, e := ctx.Connection.NewQuery().
		From(new(ScadaThreeSecs).TableName()).
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

func (d *ConvThreeExt) getMap(list []tk.M, field string) (result map[string]tk.M) {
	result = map[string]tk.M{}

	for _, val := range list {
		id := val.Get("_id").(tk.M)
		timeStamp := id.Get("timestamp").(time.Time)
		projectName := id.GetString("projectname")
		turbine := id.GetString("turbine")

		timeStampStr := timeStamp.UTC().Format("060102_150405")
		key := timeStampStr + "#" + projectName + "#" + turbine

		value := tk.M{}

		var avg float64
		var count int

		count = val.GetInt(field + "countsecs")

		// log.Printf("count: %v | %#v \n", val.GetInt(field+"_count"), key)

		if count == 0 {
			avg = emptyValueBig
			// log.Print("empty: %v \n", key)
		} else {
			avg = val.GetFloat64(field)
		}

		value.Set(field, avg)
		value.Set(field+"countsecs", count)
		result[key] = value
	}
	return
}
