package main

import (
	"bufio"
	"log"
	_ "math"
	"os"
	"strings"

	"github.com/eaciit/dbox"
	tk "github.com/eaciit/toolkit"

	// . "eaciit/wfdemo-git-dev/library/models"

	dc "eaciit/wfdemo-git-dev/processapp/threeextractor/dataconversion"

	_ "github.com/eaciit/dbox/dbc/mongo"
	"github.com/eaciit/orm"
)

var (
	wd = func() string {
		d, _ := os.Getwd()
		return d + separator
	}()
	separator = string(os.PathSeparator)
)

func main() {
	log.Println("Start Convert...")

	conn, err := PrepareConnection()
	if err != nil {
		tk.Println(err)
	}
	ctx := orm.New(conn)
	conv := dc.NewDataConversion(ctx)
	conv.Generate("DataFile20151223-00.csv")

	log.Println("End Convert.")

	/*conv := NewDataConversion(ctx)

	startTime, _ := time.Parse("20060102 15:04", "20151222 17:00")
	endTime, _ := time.Parse("20060102 15:04", "20160106 19:00")

	log.Println("Start Convert...")

	start := time.Now()
	for {
		if startTime.Format("2006-01-02 15:04") == endTime.Format("2006-01-02 15:04") {
			// log.Printf("idx: %v \n", idx)
			break
		}

		conv.Generate(startTime)
		startTime = hpp.GenNext10Minutes(startTime)
	}

	duration := time.Now().Sub(start).Seconds()
	log.Printf("End in: %v sec(s) \n", duration)*/
}

func PrepareConnection() (dbox.IConnection, error) {
	config := ReadConfig()

	// log.Printf("config: %#v \n", config)

	ci := &dbox.ConnectionInfo{config["host"], config["database"], config["username"], config["password"], tk.M{}.Set("timeout", 3000)}
	c, e := dbox.NewConnection("mongo", ci)

	if e != nil {
		return nil, e
	}

	e = c.Connect()
	if e != nil {
		return nil, e
	}

	return c, nil
}

func ReadConfig() map[string]string {
	ret := make(map[string]string)
	file, err := os.Open("../conf" + separator + "app.conf")
	if err == nil {
		defer file.Close()

		reader := bufio.NewReader(file)
		for {
			line, _, e := reader.ReadLine()
			if e != nil {
				break
			}

			sval := strings.Split(string(line), "=")
			ret[sval[0]] = sval[1]
		}
	} else {
		tk.Println(err.Error())
	}

	return ret
}

/*
type DataConversion struct {
	Ctx *orm.DataContext
}

var (
	emptyValueSmall = -0.000001
)

func NewDataConversion(ctx *orm.DataContext) *DataConversion {
	dc := new(DataConversion)
	dc.Ctx = ctx

	return dc
}

func (d *DataConversion) Generate(timestampconverted time.Time) (errorLine tk.M) {
	// funcName := "GenTenFromThreeSecond"
	// log.Printf("timeStamp: %v \n", timestampconverted.String())
	ctx := d.Ctx
	list := []tk.M{}
	pipes := []tk.M{}
	match := tk.M{"timestampconverted": timestampconverted}

	group := tk.M{
		"_id": tk.M{
			"timestamp":   "$timestampconverted",
			"projectname": "$projectname",
			"turbine":     "$turbine",
		},

		"count": tk.M{"$sum": 1},

		"fast_currentl3": tk.M{"$avg": "$fast_currentl3"},

		"fast_currentl1":                     tk.M{"$avg": "$fast_currentl1"},
		"fast_activepowersetpoint_kw":        tk.M{"$avg": "$fast_activepowersetpoint_kw"},
		"fast_currentl2":                     tk.M{"$avg": "$fast_currentl2"},
		"fast_drtrvibvalue":                  tk.M{"$avg": "$fast_drtrvibvalue"},
		"fast_genspeed_rpm":                  tk.M{"$avg": "$fast_genspeed_rpm"},
		"fast_pitchaccuv1":                   tk.M{"$avg": "$fast_pitchaccuv1"},
		"fast_pitchangle":                    tk.M{"$avg": "$fast_pitchangle"},
		"fast_pitchangle3":                   tk.M{"$avg": "$fast_pitchangle3"},
		"fast_pitchangle2":                   tk.M{"$avg": "$fast_pitchangle2"},
		"fast_pitchconvcurrent1":             tk.M{"$avg": "$fast_pitchconvcurrent1"},
		"fast_pitchconvcurrent3":             tk.M{"$avg": "$fast_pitchconvcurrent3"},
		"fast_pitchconvcurrent2":             tk.M{"$avg": "$fast_pitchconvcurrent2"},
		"fast_powerfactor":                   tk.M{"$avg": "$fast_powerfactor"},
		"fast_reactivepowersetpointppc_kvar": tk.M{"$avg": "$fast_reactivepowersetpointppc_kvar"},
		"fast_reactivepower_kvar":            tk.M{"$avg": "$fast_reactivepower_kvar"},
		"fast_rotorspeed_rpm":                tk.M{"$avg": "$fast_rotorspeed_rpm"},
		"fast_voltagel1":                     tk.M{"$avg": "$fast_voltagel1"},
		"fast_voltagel2":                     tk.M{"$avg": "$fast_voltagel2"},

		"slow_capablecapacitivereactpwr_kvar": tk.M{"$avg": "$slow_capablecapacitivereactpwr_kvar"},
		"slow_capableinductivereactpwr_kvar":  tk.M{"$avg": "$slow_capableinductivereactpwr_kvar"},
		"slow_datetime_sec":                   tk.M{"$avg": "$slow_datetime_sec"},

		"fast_pitchangle1":                tk.M{"$avg": "$fast_pitchangle1"},
		"fast_voltagel3":                  tk.M{"$avg": "$fast_voltagel3"},
		"slow_capablecapacitivepwrfactor": tk.M{"$avg": "$slow_capablecapacitivepwrfactor"},
		"fast_total_production_kwh":       tk.M{"$avg": "$fast_total_production_kwh"},
		"fast_total_prod_day_kwh":         tk.M{"$avg": "$fast_total_prod_day_kwh"},
		"fast_total_prod_month_kwh":       tk.M{"$avg": "$fast_total_prod_month_kwh"},
		"fast_activepoweroutpwcsell_kw":   tk.M{"$avg": "$fast_activepoweroutpwcsell_kw"},
		"fast_frequency_hz":               tk.M{"$avg": "$fast_frequency_hz"},
		"slow_tempg1l2":                   tk.M{"$avg": "$slow_tempg1l2"},
		"slow_tempg1l3":                   tk.M{"$avg": "$slow_tempg1l3"},
		"slow_tempgearboxhssde":           tk.M{"$avg": "$slow_tempgearboxhssde"},
		"slow_tempgearboximsnde":          tk.M{"$avg": "$slow_tempgearboximsnde"},
		"slow_tempoutdoor":                tk.M{"$avg": "$slow_tempoutdoor"},
		"fast_pitchaccuv3":                tk.M{"$avg": "$fast_pitchaccuv3"},
		"slow_totalturbineactivehours":    tk.M{"$avg": "$slow_totalturbineactivehours"},
		"slow_totalturbineokhours":        tk.M{"$avg": "$slow_totalturbineokhours"},
		"slow_totalturbinetimeallhours":   tk.M{"$avg": "$slow_totalturbinetimeallhours"},
		"slow_tempg1l1":                   tk.M{"$avg": "$slow_tempg1l1"},
		"slow_tempgearboxoilsump":         tk.M{"$avg": "$slow_tempgearboxoilsump"},
		"fast_pitchaccuv2":                tk.M{"$avg": "$fast_pitchaccuv2"},
		"slow_totalgridokhours":           tk.M{"$avg": "$slow_totalgridokhours"},
		"slow_totalactpowerout_kwh":       tk.M{"$avg": "$slow_totalactpowerout_kwh"},
		"fast_yawservice":                 tk.M{"$avg": "$fast_yawservice"},
		"fast_yawangle":                   tk.M{"$avg": "$fast_yawangle"},

		"slow_capableinductivepwrfactor": tk.M{"$avg": "$slow_capableinductivepwrfactor"},
		"slow_tempgearboxhssnde":         tk.M{"$avg": "$slow_tempgearboxhssnde"},
		"slow_temphubbearing":            tk.M{"$avg": "$slow_temphubbearing"},
		"slow_totalg1activehours":        tk.M{"$avg": "$slow_totalg1activehours"},
		"slow_totalactpoweroutg1_kwh":    tk.M{"$avg": "$slow_totalactpoweroutg1_kwh"},
		"slow_totalreactpowering1_kvarh": tk.M{"$avg": "$slow_totalreactpowering1_kvarh"},
		"slow_nacelledrill":              tk.M{"$avg": "$slow_nacelledrill"},
		"slow_tempgearboximsde":          tk.M{"$avg": "$slow_tempgearboximsde"},
		"fast_total_operating_hrs":       tk.M{"$avg": "$fast_total_operating_hrs"},
		"slow_tempnacelle":               tk.M{"$avg": "$slow_tempnacelle"},
		"fast_total_grid_ok_hrs":         tk.M{"$avg": "$fast_total_grid_ok_hrs"},
		"fast_total_wtg_ok_hrs":          tk.M{"$avg": "$fast_total_wtg_ok_hrs"},
		"slow_tempcabinettopbox":         tk.M{"$avg": "$slow_tempcabinettopbox"},
		"slow_tempgeneratorbearingnde":   tk.M{"$avg": "$slow_tempgeneratorbearingnde"},
		"fast_total_access_hrs":          tk.M{"$avg": "$fast_total_access_hrs"},
		"slow_tempbottompowersection":    tk.M{"$avg": "$slow_tempbottompowersection"},
		"slow_tempgeneratorbearingde":    tk.M{"$avg": "$slow_tempgeneratorbearingde"},
		"slow_totalreactpowerin_kvarh":   tk.M{"$avg": "$slow_totalreactpowerin_kvarh"},
		"slow_tempbottomcontrolsection":  tk.M{"$avg": "$slow_tempbottomcontrolsection"},
		"slow_tempconv1":                 tk.M{"$avg": "$slow_tempconv1"},
		"fast_activepowerrated_kw":       tk.M{"$avg": "$fast_activepowerrated_kw"},
		"fast_nodeip":                    tk.M{"$avg": "$fast_nodeip"},
		"fast_pitchspeed1":               tk.M{"$avg": "$fast_pitchspeed1"},
		"slow_cfcardsize":                tk.M{"$avg": "$slow_cfcardsize"},
		"slow_cpu_number":                tk.M{"$avg": "$slow_cpu_number"},
		"slow_cfcardspaceleft":           tk.M{"$avg": "$slow_cfcardspaceleft"},
		"slow_tempbottomcapsection":      tk.M{"$avg": "$slow_tempbottomcapsection"},
		"slow_ratedpower":                tk.M{"$avg": "$slow_ratedpower"},
		"slow_tempconv3":                 tk.M{"$avg": "$slow_tempconv3"},
		"slow_tempconv2":                 tk.M{"$avg": "$slow_tempconv2"},
		"slow_totalactpowerin_kwh":       tk.M{"$avg": "$slow_totalactpowerin_kwh"},
		"slow_totalactpowering1_kwh":     tk.M{"$avg": "$slow_totalactpowering1_kwh"},
		"slow_totalactpowering2_kwh":     tk.M{"$avg": "$slow_totalactpowering2_kwh"},
		"slow_totalactpoweroutg2_kwh":    tk.M{"$avg": "$slow_totalactpoweroutg2_kwh"},
		"slow_totalg2activehours":        tk.M{"$avg": "$slow_totalg2activehours"},
		"slow_totalreactpowering2_kvarh": tk.M{"$avg": "$slow_totalreactpowering2_kvarh"},
		"slow_totalreactpowerout_kvarh":  tk.M{"$avg": "$slow_totalreactpowerout_kvarh"},
		"slow_utcoffset_int":             tk.M{"$avg": "$slow_utcoffset_int"},
	}

	pipes = append(pipes, tk.M{"$match": match})
	pipes = append(pipes, tk.M{"$group": group})
	// pipes = append(pipes, tk.M{"$sort": tk.M{"_id": 1}})

	// tk.Printf("pipes: %#v \n", pipes)

	csr, e := ctx.Connection.NewQuery().
		From(new(ScadaThreeSecs).TableName()).
		Command("pipe", pipes).
		Cursor(nil)

	if e != nil {
		log.Printf("ERR: %#v \n", e.Error())
	} else {
		e = csr.Fetch(&list, 0, false)

		log.Printf("timeStamp: %v | %v \n", timestampconverted.String(), len(list))

		fastActivePowerKWList := d.getAvgMinMaxCount(ctx, timestampconverted, "fast_activepower_kw")
		fastActivePowerKWMap := d.getMap(fastActivePowerKWList, "fast_activepower_kw")

		fastWindSpeedMsList := d.getAvgMinMaxCount(ctx, timestampconverted, "fast_windspeed_ms")
		fastWindSpeedMsMap := d.getMap(fastWindSpeedMsList, "fast_windspeed_ms")

		slowNacellePosList := d.getAvgMinMaxCount(ctx, timestampconverted, "slow_nacellepos")
		slowNacellePosMap := d.getMap(slowNacellePosList, "slow_nacellepos")

		slowWindDirectionList := d.getAvgMinMaxCount(ctx, timestampconverted, "slow_winddirection")
		slowWindDirectionMap := d.getMap(slowWindDirectionList, "slow_winddirection")

		for idx, val := range list {
			// errorList := []error{}

			tenScada := new(ScadaConvTenMin)
			tenScada.No = idx + 1
			// tenScada.File = file

			id := val.Get("_id").(tk.M)
			timeStamp := id.Get("timestamp").(time.Time)
			projectName := id.GetString("projectname")
			turbine := id.GetString("turbine")

			timeStampStr := timeStamp.UTC().Format("060102_1504")
			key := timeStampStr + "_" + projectName + "_" + turbine

			tenScada.TimeStamp = timeStamp
			tenScada.DateInfo = helper.GetDateInfo(timeStamp)
			tenScada.ProjectName = projectName
			tenScada.Turbine = turbine

			tenScada = tenScada.New()

			tenScada.Fast_CurrentL3 = val.GetFloat64("fast_currentl3")

			fastActivePower := fastActivePowerKWMap[key]

			tenScada.Fast_ActivePower_kW = fastActivePower.GetFloat64("fast_activepower_kw")
			tenScada.Fast_ActivePower_kW_Min = fastActivePower.GetFloat64("fast_activepower_kw_min")
			tenScada.Fast_ActivePower_kW_Max = fastActivePower.GetFloat64("fast_activepower_kw_max")
			tenScada.Fast_ActivePower_kW_Count = fastActivePower.GetInt("fast_activepower_kw_count")

			tenScada.Fast_CurrentL1 = val.GetFloat64("fast_currentl1")
			tenScada.Fast_ActivePowerSetpoint_kW = val.GetFloat64("fast_activepowersetpoint_kw")
			tenScada.Fast_CurrentL2 = val.GetFloat64("fast_currentl2")
			tenScada.Fast_DrTrVibValue = val.GetFloat64("fast_drtrvibvalue")
			tenScada.Fast_GenSpeed_RPM = val.GetFloat64("fast_genspeed_rpm")
			tenScada.Fast_PitchAccuV1 = val.GetFloat64("fast_pitchaccuv1")
			tenScada.Fast_PitchAngle = val.GetFloat64("fast_pitchangle")
			tenScada.Fast_PitchAngle3 = val.GetFloat64("fast_pitchangle3")
			tenScada.Fast_PitchAngle2 = val.GetFloat64("fast_pitchangle2")
			tenScada.Fast_PitchConvCurrent1 = val.GetFloat64("fast_pitchconvcurrent1")
			tenScada.Fast_PitchConvCurrent3 = val.GetFloat64("fast_pitchconvcurrent3")
			tenScada.Fast_PitchConvCurrent2 = val.GetFloat64("fast_pitchconvcurrent2")
			tenScada.Fast_PowerFactor = val.GetFloat64("fast_powerfactor")
			tenScada.Fast_ReactivePowerSetpointPPC_kVAr = val.GetFloat64("fast_reactivepowersetpointppc_kvar")
			tenScada.Fast_ReactivePower_kVAr = val.GetFloat64("fast_reactivepower_kvar")
			tenScada.Fast_RotorSpeed_RPM = val.GetFloat64("fast_rotorspeed_rpm")
			tenScada.Fast_VoltageL1 = val.GetFloat64("fast_voltagel1")
			tenScada.Fast_VoltageL2 = val.GetFloat64("fast_voltagel2")

			fastWindSpeedMs := fastWindSpeedMsMap[key]

			tenScada.Fast_WindSpeed_ms = fastWindSpeedMs.GetFloat64("fast_windspeed_ms")
			tenScada.Fast_WindSpeed_ms_Min = fastWindSpeedMs.GetFloat64("fast_windspeed_ms_min")
			tenScada.Fast_WindSpeed_ms_Max = fastWindSpeedMs.GetFloat64("fast_windspeed_ms_max")
			tenScada.Fast_WindSpeed_ms_Count = fastWindSpeedMs.GetInt("fast_windspeed_ms_count")

			tenScada.Slow_CapableCapacitiveReactPwr_kVAr = val.GetFloat64("slow_capablecapacitivereactpwr_kvar")
			tenScada.Slow_CapableInductiveReactPwr_kVAr = val.GetFloat64("slow_capableinductivereactpwr_kvar")
			tenScada.Slow_DateTime_Sec = val.GetFloat64("slow_datetime_sec")

			slowNacellePos := slowNacellePosMap[key]

			tenScada.Slow_NacellePos = slowNacellePos.GetFloat64("slow_nacellepos")
			tenScada.Slow_NacellePos_Min = slowNacellePos.GetFloat64("slow_nacellepos_min")
			tenScada.Slow_NacellePos_Max = slowNacellePos.GetFloat64("slow_nacellepos_max")
			tenScada.Slow_NacellePos_Count = slowNacellePos.GetInt("slow_nacellepos_count")

			tenScada.Fast_PitchAngle1 = val.GetFloat64("fast_pitchangle1")
			tenScada.Fast_VoltageL3 = val.GetFloat64("fast_voltagel3")
			tenScada.Slow_CapableCapacitivePwrFactor = val.GetFloat64("slow_capablecapacitivepwrfactor")
			tenScada.Fast_Total_Production_kWh = val.GetFloat64("fast_total_production_kwh")
			tenScada.Fast_Total_Prod_Day_kWh = val.GetFloat64("fast_total_prod_day_kwh")
			tenScada.Fast_Total_Prod_Month_kWh = val.GetFloat64("fast_total_prod_month_kwh")
			tenScada.Fast_ActivePowerOutPWCSell_kW = val.GetFloat64("fast_activepoweroutpwcsell_kw")
			tenScada.Fast_Frequency_Hz = val.GetFloat64("fast_frequency_hz")
			tenScada.Slow_TempG1L2 = val.GetFloat64("slow_tempg1l2")
			tenScada.Slow_TempG1L3 = val.GetFloat64("slow_tempg1l3")
			tenScada.Slow_TempGearBoxHSSDE = val.GetFloat64("slow_tempgearboxhssde")
			tenScada.Slow_TempGearBoxIMSNDE = val.GetFloat64("slow_tempgearboximsnde")
			tenScada.Slow_TempOutdoor = val.GetFloat64("slow_tempoutdoor")
			tenScada.Fast_PitchAccuV3 = val.GetFloat64("fast_pitchaccuv3")
			tenScada.Slow_TotalTurbineActiveHours = val.GetFloat64("slow_totalturbineactivehours")
			tenScada.Slow_TotalTurbineOKHours = val.GetFloat64("slow_totalturbineokhours")
			tenScada.Slow_TotalTurbineTimeAllHours = val.GetFloat64("slow_totalturbinetimeallhours")
			tenScada.Slow_TempG1L1 = val.GetFloat64("slow_tempg1l1")
			tenScada.Slow_TempGearBoxOilSump = val.GetFloat64("slow_tempgearboxoilsump")
			tenScada.Fast_PitchAccuV2 = val.GetFloat64("fast_pitchaccuv2")
			tenScada.Slow_TotalGridOkHours = val.GetFloat64("slow_totalgridokhours")
			tenScada.Slow_TotalActPowerOut_kWh = val.GetFloat64("slow_totalactpowerout_kwh")
			tenScada.Fast_YawService = val.GetFloat64("fast_yawservice")
			tenScada.Fast_YawAngle = val.GetFloat64("fast_yawangle")

			slowWindDirection := slowWindDirectionMap[key]

			tenScada.Slow_WindDirection = slowWindDirection.GetFloat64("slow_winddirection")
			tenScada.Slow_WindDirection_Min = slowWindDirection.GetFloat64("slow_winddirection_min")
			tenScada.Slow_WindDirection_Max = slowWindDirection.GetFloat64("slow_winddirection_max")
			tenScada.Slow_WindDirection_Count = slowWindDirection.GetInt("slow_winddirection_count")

			tenScada.Slow_CapableInductivePwrFactor = val.GetFloat64("slow_capableinductivepwrfactor")
			tenScada.Slow_TempGearBoxHSSNDE = val.GetFloat64("slow_tempgearboxhssnde")
			tenScada.Slow_TempHubBearing = val.GetFloat64("slow_temphubbearing")
			tenScada.Slow_TotalG1ActiveHours = val.GetFloat64("slow_totalg1activehours")
			tenScada.Slow_TotalActPowerOutG1_kWh = val.GetFloat64("slow_totalactpoweroutg1_kwh")
			tenScada.Slow_TotalReactPowerInG1_kVArh = val.GetFloat64("slow_totalreactpowering1_kvarh")
			tenScada.Slow_NacelleDrill = val.GetFloat64("slow_nacelledrill")
			tenScada.Slow_TempGearBoxIMSDE = val.GetFloat64("slow_tempgearboximsde")
			tenScada.Fast_Total_Operating_hrs = val.GetFloat64("fast_total_operating_hrs")
			tenScada.Slow_TempNacelle = val.GetFloat64("slow_tempnacelle")
			tenScada.Fast_Total_Grid_OK_hrs = val.GetFloat64("fast_total_grid_ok_hrs")
			tenScada.Fast_Total_WTG_OK_hrs = val.GetFloat64("fast_total_wtg_ok_hrs")
			tenScada.Slow_TempCabinetTopBox = val.GetFloat64("slow_tempcabinettopbox")
			tenScada.Slow_TempGeneratorBearingNDE = val.GetFloat64("slow_tempgeneratorbearingnde")
			tenScada.Fast_Total_Access_hrs = val.GetFloat64("fast_total_access_hrs")
			tenScada.Slow_TempBottomPowerSection = val.GetFloat64("slow_tempbottompowersection")
			tenScada.Slow_TempGeneratorBearingDE = val.GetFloat64("slow_tempgeneratorbearingde")
			tenScada.Slow_TotalReactPowerIn_kVArh = val.GetFloat64("slow_totalreactpowerin_kvarh")
			tenScada.Slow_TempBottomControlSection = val.GetFloat64("slow_tempbottomcontrolsection")
			tenScada.Slow_TempConv1 = val.GetFloat64("slow_tempconv1")
			tenScada.Fast_ActivePowerRated_kW = val.GetFloat64("fast_activepowerrated_kw")
			tenScada.Fast_NodeIP = val.GetFloat64("fast_nodeip")
			tenScada.Fast_PitchSpeed1 = val.GetFloat64("fast_pitchspeed1")
			tenScada.Slow_CFCardSize = val.GetFloat64("slow_cfcardsize")
			tenScada.Slow_CPU_Number = val.GetFloat64("slow_cpu_number")
			tenScada.Slow_CFCardSpaceLeft = val.GetFloat64("slow_cfcardspaceleft")
			tenScada.Slow_TempBottomCapSection = val.GetFloat64("slow_tempbottomcapsection")
			tenScada.Slow_RatedPower = val.GetFloat64("slow_ratedpower")
			tenScada.Slow_TempConv3 = val.GetFloat64("slow_tempconv3")
			tenScada.Slow_TempConv2 = val.GetFloat64("slow_tempconv2")
			tenScada.Slow_TotalActPowerIn_kWh = val.GetFloat64("slow_totalactpowerin_kwh")
			tenScada.Slow_TotalActPowerInG1_kWh = val.GetFloat64("slow_totalactpowering1_kwh")
			tenScada.Slow_TotalActPowerInG2_kWh = val.GetFloat64("slow_totalactpowering2_kwh")
			tenScada.Slow_TotalActPowerOutG2_kWh = val.GetFloat64("slow_totalactpoweroutg2_kwh")
			tenScada.Slow_TotalG2ActiveHours = val.GetFloat64("slow_totalg2activehours")
			tenScada.Slow_TotalReactPowerInG2_kVArh = val.GetFloat64("slow_totalreactpowering2_kvarh")
			tenScada.Slow_TotalReactPowerOut_kVArh = val.GetFloat64("slow_totalreactpowerout_kvarh")
			tenScada.Slow_UTCoffset_int = val.GetFloat64("slow_utcoffset_int")

			tenScada.Count = val.GetInt("count")

			// log.Printf("%#v \n", val)

			// if len(errorList) == 0 {
			err := ctx.Insert(tenScada)

			if err != nil {
				log.Printf("err: %#v \n", err)
			}

			// ErrorLog(err, funcName, errorList)
			// ErrorHandler(err, "Saving")
			// }

			// if len(errorList) > 0 {
			// 	errorLine.Set(tk.ToString(tenScada.No), errorList)
			// }
		}

	}

	csr.Close()

	return
}

func (d *DataConversion) getAvgMinMaxCount(ctx *orm.DataContext, timestampconverted time.Time, field string) (result []tk.M) {
	pipes := []tk.M{}

	match := tk.M{
		"timestampconverted": timestampconverted,
		field:                tk.M{"$ne": emptyValueSmall},
	}

	group := tk.M{
		"_id": tk.M{
			"timestamp":   "$timestampconverted",
			"projectname": "$projectname",
			"turbine":     "$turbine",
		},
		field:            tk.M{"$avg": "$" + field},
		field + "_min":   tk.M{"$min": "$" + field},
		field + "_max":   tk.M{"$max": "$" + field},
		field + "_count": tk.M{"$sum": 1},
	}

	pipes = append(pipes, tk.M{"$match": match})
	pipes = append(pipes, tk.M{"$group": group})
	// pipes = append(pipes, tk.M{"$sort": tk.M{"_id": 1}})

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

func (d *DataConversion) getMap(list []tk.M, field string) (result map[string]tk.M) {
	result = map[string]tk.M{}

	for _, val := range list {
		id := val.Get("_id").(tk.M)
		timeStamp := id.Get("timestamp").(time.Time)
		projectName := id.GetString("projectname")
		turbine := id.GetString("turbine")

		timeStampStr := timeStamp.Format("060102_1504")
		key := timeStampStr + "#" + projectName + "#" + turbine

		value := tk.M{}
		value.Set(field, val.Get(field))
		value.Set(field+"_min", val.Get(field+"_min"))
		value.Set(field+"_max", val.Get(field+"_max"))
		value.Set(field+"_count", val.Get(field+"_count"))

		result[key] = value
	}

	return
}
*/
