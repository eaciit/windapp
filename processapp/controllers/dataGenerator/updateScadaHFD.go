package generatorControllers

import (
	. "eaciit/wfdemo-git/library/helper"
	. "eaciit/wfdemo-git/library/models"
	. "eaciit/wfdemo-git/processapp/controllers"
	"os"
	"sync"
	"time"

	"github.com/eaciit/dbox"
	_ "github.com/eaciit/dbox/dbc/mongo"
	tk "github.com/eaciit/toolkit"
)

type UpdateScadaHFD struct {
	*BaseController
}

var (
	hmtx = &sync.Mutex{}
)

func (c *UpdateScadaHFD) DoUpdateWsBin(base *BaseController) {
	funcName := "Update ws bins for ScadaHFD"
	c.BaseController = base

	var wg sync.WaitGroup

	if base != nil {
		ctx, e := PrepareConnection()
		if e != nil {
			ErrorHandler(e, "Scada Summary")
			os.Exit(0)
		}

		csr, e := ctx.NewQuery().From(new(ScadaDataHFD).TableName()).Cursor(nil)

		defer csr.Close()

		counter := 0
		countData := csr.Count()
		isDone := false
		countPerProcess := 1000

		for !isDone && countData > 0 {
			scadas := []*ScadaDataHFD{}

			e = csr.Fetch(&scadas, countPerProcess, false)
			ErrorHandler(e, funcName)

			if len(scadas) < countPerProcess {
				isDone = true
			}

			wg.Add(1)
			go func(datas []*ScadaDataHFD, counter int) {
				tk.Println("start process ", countPerProcess*(counter+1))
				for _, d := range datas {
					hmtx.Lock()

					dId := d.ID
					wsBin := tk.RoundingAuto64(d.Fast_WindSpeed_ms, 0)
					tk.Println("Updating data for ID = ", dId, wsBin)
					e = ctx.NewQuery().Update().From(new(ScadaDataHFD).TableName()).
						Where(dbox.Eq("_id", dId)).
						Exec(tk.M{}.Set("data", tk.M{}.Set("fast_windspeed_bin", wsBin)))
					ErrorHandler(e, funcName)

					hmtx.Unlock()
				}
				tk.Println("end process ", countPerProcess*(counter+1))
				wg.Done()
			}(scadas, counter)

			counter++
			if counter%10 == 0 || isDone {
				wg.Wait()
			}
		}
	}

	tk.Println("End process updating wind speed bin for ScadaData HFD...")
}

func (c *UpdateScadaHFD) DoUpdateTurbineStateForScadaData(base *BaseController) {
	funcName := "Update turbine state for Scada"
	c.BaseController = base

	startTime := time.Now()

	var wg sync.WaitGroup

	if base != nil {
		ctx, e := PrepareConnection()
		if e != nil {
			ErrorHandler(e, "Update Turbine State Scada Data")
			os.Exit(0)
		}

		validState := tk.M{
			"Tejuva": tk.M{
				"10": true,
				"11": true,
			},
			"Lahori": tk.M{
				"7": true,
			},
			"Amba": tk.M{
				"100": true,
			},
		}

		dtAlarmBrakes := []tk.M{}
		csrAb, e := ctx.NewQuery().From("ref_turbinestate").Order("projectname").Cursor(nil)
		e = csrAb.Fetch(&dtAlarmBrakes, 0, false)

		defer csrAb.Close()
		alarmBrakes := make(map[string]map[int]string)
		if len(dtAlarmBrakes) > 0 {
			for _, dt := range dtAlarmBrakes {
				project := dt.GetString("projectname")
				turbineState := dt.GetInt("turbinestate")
				turbineStateDesc := dt.GetString("description")
				if _, ok := alarmBrakes[project]; !ok {
					alarmBrakes[project] = make(map[int]string)
				}
				alarmBrakes[project][turbineState] = turbineStateDesc
			}
		}

		// get min max ts per project
		var pipes []tk.M
		// timeMaxFilter, _ := time.Parse("2006-01-02 15:04:05", "2017-11-01 00:00:00")
		pipes = append(pipes, tk.M{
			// "$match": tk.M{
			// 	"timestamp": tk.M{"$lte": timeMaxFilter},
			// },
			"$group": tk.M{
				"_id": "$projectname",
				"minTs": tk.M{
					"$min": "$timestamp",
				},
				"maxTs": tk.M{
					"$max": "$timestamp",
				},
			},
		})
		csrts, e := ctx.NewQuery().From("ScadaData").Command("pipe", pipes).Cursor(nil)
		ErrorHandler(e, "Getting min max TS")

		dataTs := []tk.M{}
		e = csrts.Fetch(&dataTs, 0, false)
		defer csrts.Close()

		tk.Printf("%#v\n", dataTs)

		ErrorHandler(e, "Fetching data min max TS")

		dayToProcess := 7
		if len(dataTs) > 0 {
			for _, dt := range dataTs {
				projectname := dt.GetString("_id")
				tmin := dt.Get("minTs").(time.Time)
				tmax := dt.Get("maxTs").(time.Time)
				tstart := tmin
				isFinish := false

				for {
					tend := tstart.Add(time.Duration(dayToProcess) * 24 * time.Hour)
					// tk.Printf("tend = %#v\n", tend.Format("2006-01-02 15:04:05"))
					if tend.Sub(tmax) > 0 {
						isFinish = true
					}

					// process updating data
					tk.Printf("Start processing %v from %v to %v\n", projectname, tstart.Format("2006-01-02 15:04:05"), tend.Format("2006-01-02 15:04:05"))
					csr, e := ctx.NewQuery().From("ScadaData").Where(dbox.And(dbox.Eq("projectname", projectname), dbox.Gte("timestamp", tstart), dbox.Lte("timestamp", tend))).Cursor(nil)

					defer csr.Close()

					counter := 0
					countData := csr.Count()
					tk.Println("Total data will processed :", countData)
					isDone := false
					countPerProcess := 1000

					for !isDone && countData > 0 {
						scadas := []*tk.M{}

						startTimeFetch := time.Now()

						e = csr.Fetch(&scadas, countPerProcess, false)
						ErrorHandler(e, funcName)

						if len(scadas) < countPerProcess {
							isDone = true
						}

						durationFetch := time.Now().Sub(startTimeFetch).Seconds()
						tk.Printf("Fetch %v data about %v sec(s)\n", len(scadas), durationFetch)

						wg.Add(1)
						go func(datas []*tk.M, counter int) {
							startTimeGo := time.Now()
							tk.Println("start process ", countPerProcess*(counter+1))
							countupdate := 0
							for _, d := range datas {
								hmtx.Lock()

								dId := d.GetString("_id")
								timestamp, _ := time.Parse("2006-01-02T15:04:05Z07:00", d.GetString("timestamp"))
								timestamp = timestamp.UTC().Add(7 * time.Hour) // this add by 7 because run in local
								turbine := d.GetString("turbine")

								turbineState := -999997
								turbineStateDesc := ""
								isvalidstate := false

								dtsource := []tk.M{}
								csrTs, e := ctx.NewQuery().From("Scada10MinHFD").
									Where(dbox.And(dbox.Eq("turbine", turbine), dbox.Eq("timestamp", timestamp))).
									Cursor(nil)
								e = csrTs.Fetch(&dtsource, 1, false)
								csrTs.Close()

								if len(dtsource) > 0 {
									turbineState = tk.ToInt(dtsource[0].GetFloat64("turbinestate"), tk.RoundingAuto)
								}

								// getting turbine state from alarm
								if validState.Has(projectname) {
									dtVs := validState.Get(projectname).(tk.M)
									if dtVs.Has(tk.ToString(turbineState)) {
										isvalidstate = true
									}
								}

								if desc, ok := alarmBrakes[projectname][turbineState]; ok {
									turbineStateDesc = desc
								}

								e = ctx.NewQuery().Update().From("ScadaData").
									Where(dbox.Eq("_id", dId)).
									Exec(tk.M{}.Set("data", tk.M{}.Set("turbinestate", turbineState).Set("statedescription", turbineStateDesc).Set("isvalidstate", isvalidstate)))
								ErrorHandler(e, "Update scada data")

								countupdate++

								hmtx.Unlock()
							}
							durationGo := time.Now().Sub(startTimeGo).Seconds()
							tk.Printf("\nEnd of process %v data about %v sec(s) \n", countPerProcess*(counter+1), durationGo)
							wg.Done()
						}(scadas, counter)

						counter++
						if counter%10 == 0 || isDone {
							wg.Wait()
						}
					}

					if isFinish {
						break
					}
					tstart = tend
					// tk.Printf("tend = %#v\n", tstart.Format("2006-01-02 15:04:05"))
				}
			}
		}
	}

	duration := time.Now().Sub(startTime).Seconds()

	tk.Printf("End process updating wind speed bin for ScadaData HFD about %v sec(s)\n", duration)
}
