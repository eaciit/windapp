package main

import (
	. "eaciit/wfdemo-git/library/helper"
	"flag"
	"github.com/eaciit/dbox"
	_ "github.com/eaciit/dbox/dbc/mongo"
	tk "github.com/eaciit/toolkit"
	"runtime"
	"strings"
	"sync"
	"time"
)

var (
	log      *tk.LogEngine
	mapField = map[string]MappingColumn{
		"windspeed":                   MappingColumn{"Wind Speed", "WindSpeed_ms", "m/s", 0.0, 50.0, "$avg"},
		"power":                       MappingColumn{"Power", "ActivePower_kW", "kW", -200, 2100.0 + (2100.0 * 0.10), "$sum"},
		"production":                  MappingColumn{"Production", "Production", "kWh", -200, 2100.0, "$sum"},
		"winddirection":               MappingColumn{"Wind Direction", "WindDirection", "Degree", 0.0, 360.0, "$avg"},
		"nacellepos":                  MappingColumn{"Nacelle Direction", "NacellePos", "Degree", 0.0, 360.0, "$avg"},
		"rotorrpm":                    MappingColumn{"Rotor RPM", "RotorSpeed_RPM", "RPM", 0.0, 30.0, "$avg"},
		"genrpm":                      MappingColumn{"Generator RPM", "GenSpeed_RPM", "RPM", 0.0, 30.0, "$avg"},
		"pitchangle":                  MappingColumn{"Pitch Angle", "PitchAngle", "Degree", -10.0, 120.0, "$avg"},
		"PitchCabinetTempBlade1":      MappingColumn{"Pitch Cabinet Temp Blade 1", "PitchCabinetTempBlade1", "Degree", -10.0, 120.0, "$avg"},
		"PitchCabinetTempBlade2":      MappingColumn{"Pitch Cabinet Temp Blade 2", "PitchCabinetTempBlade2", "Degree", -10.0, 120.0, "$avg"},
		"PitchCabinetTempBlade3":      MappingColumn{"Pitch Cabinet Temp Blade 3", "PitchCabinetTempBlade3", "Degree", -10.0, 120.0, "$avg"},
		"PitchConvInternalTempBlade1": MappingColumn{"Pitch Conv Internal Temp Blade 1", "PitchConvInternalTempBlade1", "Degree", -10.0, 120.0, "$avg"},
		"PitchConvInternalTempBlade2": MappingColumn{"Pitch Conv Internal Temp Blade 2", "PitchConvInternalTempBlade2", "Degree", -10.0, 120.0, "$avg"},
		"PitchConvInternalTempBlade3": MappingColumn{"Pitch Conv Internal Temp Blade 3", "PitchConvInternalTempBlade3", "Degree", -10.0, 120.0, "$avg"},
		"TempG1L1":                    MappingColumn{"TempG1 L1", "TempG1L1", "Degree", -10.0, 200.0, "$avg"},
		"TempG1L2":                    MappingColumn{"TempG1 L2", "TempG1L2", "Degree", -10.0, 200.0, "$avg"},
		"TempG1L3":                    MappingColumn{"TempG1 L3", "TempG1L3", "Degree", -10.0, 200.0, "$avg"},
		"TempGeneratorBearingDE":      MappingColumn{"Temp Generator Bearing DE", "TempGeneratorBearingDE", "Degree", -10.0, 200.0, "$avg"},
		"TempGeneratorBearingNDE":     MappingColumn{"Temp Generator Bearing NDE", "TempGeneratorBearingNDE", "Degree", -10.0, 200.0, "$avg"},
		"TempGearBoxOilSump":          MappingColumn{"Temp Gear Box Oil Sump", "TempGearBoxOilSump", "Degree", -10.0, 200.0, "$avg"},
		"TempHubBearing":              MappingColumn{"Temp Hub Bearing", "TempHubBearing", "Degree", -10.0, 200.0, "$avg"},
		"TempGeneratorChoke":          MappingColumn{"Temp Generator Choke", "TempGeneratorChoke", "Degree", -10.0, 200.0, "$avg"},
		"TempGridChoke":               MappingColumn{"Temp Grid Choke", "TempGridChoke", "Degree", -10.0, 200.0, "$avg"},
		"TempGeneratorCoolingUnit":    MappingColumn{"Temp Generator Cooling Unit", "TempGeneratorCoolingUnit", "Degree", -10.0, 200.0, "$avg"},
		"TempConvCabinet2":            MappingColumn{"Temp Conv Cabinet 2", "TempConvCabinet2", "Degree", -10.0, 200.0, "$avg"},
		"TempOutdoor":                 MappingColumn{"Temp Outdoor", "TempOutdoor", "Degree", -10.0, 200.0, "$avg"},
		"TempSlipRing":                MappingColumn{"Temp Slip Ring", "TempSlipRing", "Degree", -10.0, 200.0, "$avg"},
		"TransformerWindingTemp1":     MappingColumn{"Transformer Winding Temp 1", "TransformerWindingTemp1", "Degree", -10.0, 200.0, "$avg"},
		"TransformerWindingTemp2":     MappingColumn{"Transformer Winding Temp 2", "TransformerWindingTemp2", "Degree", -10.0, 200.0, "$avg"},
		"TransformerWindingTemp3":     MappingColumn{"Transformer Winding Temp 3", "TransformerWindingTemp3", "Degree", -10.0, 200.0, "$avg"},
		"TempShaftBearing1":           MappingColumn{"Temp Shaft Bearing 1", "TempShaftBearing1", "Degree", -10.0, 200.0, "$avg"},
		"TempShaftBearing2":           MappingColumn{"Temp Shaft Bearing 2", "TempShaftBearing2", "Degree", -10.0, 200.0, "$avg"},
		"TempShaftBearing3":           MappingColumn{"Temp Shaft Bearing 3", "TempShaftBearing3", "Degree", -10.0, 200.0, "$avg"},
		"TempGearBoxIMSDE":            MappingColumn{"Temp Gear Box IMS DE", "TempGearBoxIMSDE", "Degree", -10.0, 200.0, "$avg"},
		"TempBottomControlSection":    MappingColumn{"Temp Bottom Control Section", "TempBottomControlSection", "Degree", -10.0, 200.0, "$avg"},
		"TempBottomPowerSection":      MappingColumn{"Temp Bottom Power Section", "TempBottomPowerSection", "Degree", -10.0, 200.0, "$avg"},
		"TempCabinetTopBox":           MappingColumn{"Temp Cabinet Top Box", "TempCabinetTopBox", "Degree", -10.0, 200.0, "$avg"},
		"TempNacelle":                 MappingColumn{"Temp Nacelle", "TempNacelle", "Degree", -10.0, 200.0, "$avg"},
		"VoltageL1":                   MappingColumn{"Voltage L 1", "GridPPVPhaseAB", "Volt", -10.0, 1000.0, "$avg"},
		"VoltageL2":                   MappingColumn{"Voltage L 2", "GridPPVPhaseBC", "Volt", -10.0, 1000.0, "$avg"},
		"VoltageL3":                   MappingColumn{"Voltage L 3", "GridPPVPhaseCA", "Volt", -10.0, 1000.0, "$avg"},
		"PitchAccuV1":                 MappingColumn{"Blade Voltage V1", "PitchAccuV1", "Volt", -10.0, 1000.0, "$avg"},
		"PitchAccuV2":                 MappingColumn{"Blade Voltage V2", "PitchAccuV2", "Volt", -10.0, 1000.0, "$avg"},
		"PitchAccuV3":                 MappingColumn{"Blade Voltage V3", "PitchAccuV3", "Volt", -10.0, 1000.0, "$avg"},
	}
)

type MappingColumn struct {
	Name      string
	SecField  string
	Unit      string
	MinValue  float64
	MaxValue  float64
	Aggregate string
}

const (
	sError   = "ERROR"
	sInfo    = "INFO"
	sWarning = "WARNING"
)

func main() {
	logpath := ""
	flag.StringVar(&logpath, "log", "", "Log folder place")
	flag.Parse()
	config := ReadConfig()
	if logpath == "" {
		logpath, _ = config["logpathtimeseries"]
	}
	log, _ = tk.NewLog(false, true, logpath, "simpleHFDLog_%s", "20060102")
	ctx, e := PrepareConnection(config)
	if e != nil {
		log.AddLog(e.Error(), sError)
	}
	defer ctx.Close()

	tagList := []string{"_id", "timestamp", "dateinfo", "projectname", "turbine"}
	for _, val := range mapField {
		tagList = append(tagList, strings.ToLower(val.SecField))
	}

	csrLog, e := ctx.NewQuery().From("log_latestdaterun").
		Where(dbox.Eq("type", "timeseries")).Cursor(nil)
	if e != nil {
		log.AddLog(e.Error(), sError)
	}
	defer csrLog.Close()
	lastData := []struct {
		ProjectName string
		LastDate    time.Time
	}{}
	e = csrLog.Fetch(&lastData, 0, false)
	if e != nil {
		log.AddLog(e.Error(), sError)
	}
	lastDatePerProject := map[string]time.Time{}
	for _, val := range lastData {
		lastDatePerProject[val.ProjectName] = val.LastDate
	}

	csrProject, e := ctx.NewQuery().From("ref_project").
		Where(dbox.Eq("active", true)).Cursor(nil)
	if e != nil {
		log.AddLog(e.Error(), sError)
	}
	defer csrProject.Close()
	projectList := []struct {
		ProjectID   string
		ProjectName string
	}{}
	e = csrProject.Fetch(&projectList, 0, false)
	if e != nil {
		log.AddLog(e.Error(), sError)
	}
	var wgProject sync.WaitGroup
	wgProject.Add(len(projectList))

	for _, project := range projectList {
		go func(projectid string) {
			csrData, e := ctx.NewQuery().From("Scada10MinHFD").Select(tagList...).
				Where(dbox.And(
					dbox.Eq("projectname", projectid),
					dbox.Gte("timestamp", lastDatePerProject[projectid]),
					dbox.Eq("isnull", false))).
				Cursor(nil)
			if e != nil {
				log.AddLog(e.Error(), sError)
			}
			defer csrData.Close()

			maxTimeStamp := time.Time{}

			var wg sync.WaitGroup
			totalData := csrData.Count()
			totalWorker := runtime.NumCPU() * 2
			chanData := make(chan tk.M, totalData)
			step := getstep(totalData)
			tNow := time.Now()

			wg.Add(totalWorker)
			for i := 0; i < totalWorker; i++ {
				go func() {
					ctxWorker, e := PrepareConnection(config)
					if e != nil {
						log.AddLog(e.Error(), sError)
					}
					defer ctxWorker.Close()
					csrSave := ctxWorker.NewQuery().From("TimeSeriesHFD").SetConfig("multiexec", true).Save()
					defer csrSave.Close()
					for data := range chanData {
						if data.GetInt("count")%step == 0 {
							percent := tk.ToInt(tk.Div(float64(data.GetInt("count"))*100.0, float64(totalData)), tk.RoundingUp)
							log.AddLog(tk.Sprintf("[%s] Saving %d of %d (%d percent) in %s\n",
								strings.ToUpper(projectid), data.GetInt("count"), totalData, percent,
								time.Since(tNow).String()), sInfo)
						}
						data.Unset("count")
						csrSave.Exec(tk.M{"data": data})
					}
					wg.Done()
				}()
			}

			log.AddLog(tk.Sprintf("Processing %d data [%s] with %d step using %d CPU since %s",
				totalData, strings.ToUpper(projectid), step, totalWorker, lastDatePerProject[projectid].Format("20060102_150405")), sInfo)

			count := 0
			_data := tk.M{}
			currTimeStamp := time.Time{}
			for {
				count++
				_data = tk.M{}
				e = csrData.Fetch(&_data, 1, false)
				if e != nil {
					if !strings.Contains(e.Error(), "Not found") {
						log.AddLog(e.Error(), sError)
					}
					break
				}
				currTimeStamp = _data.Get("timestamp", time.Time{}).(time.Time).UTC()
				if currTimeStamp.After(maxTimeStamp) {
					maxTimeStamp = currTimeStamp
				}

				_data.Set("count", count)
				chanData <- _data

				// if count%step == 0 {
				// 	log.AddLog(tk.Sprintf("Processing %d of %d in %s\n",
				// 		count, totalData,
				// 		time.Since(tNow).String()), sInfo)
				// }
			}
			close(chanData)
			wg.Wait()

			if maxTimeStamp.Year() > 1 {
				e = ctx.NewQuery().From("log_latestdaterun").Save().
					Exec(tk.M{"data": tk.M{
						"_id":         "timeseries_hfd_" + projectid,
						"lastdate":    maxTimeStamp,
						"projectname": projectid,
						"type":        "timeseries",
					}})
				if e != nil {
					log.AddLog(e.Error(), sError)
				}
			}
			wgProject.Done()
		}(project.ProjectID)
	}
	wgProject.Wait()
}

func getstep(count int) int {
	v := count / 5
	if v == 0 {
		return 1
	}
	return v
}

func PrepareConnection(config map[string]string) (dbox.IConnection, error) {
	ci := &dbox.ConnectionInfo{config["host"], config["database"], config["username"], config["password"], nil}
	c, e := dbox.NewConnection("mongo", ci)

	if e != nil {
		return nil, e
	}

	e = c.Connect()
	if e != nil {
		return nil, e
	}
	// log.AddLog(tk.Sprintf("DB Connect %s : %s", config["host"], config["database"]), sInfo)
	return c, nil
}
