package eventconversion

import (
	"io/ioutil"
	"os"
	"strings"
	"sync"
	"time"

	. "github.com/eaciit/windapp/library/helper"
	. "github.com/eaciit/windapp/library/models"

	_ "github.com/eaciit/dbox/dbc/mongo"
	"github.com/eaciit/orm"
	tk "github.com/eaciit/toolkit"
	"github.com/tealeg/xlsx"
)

var (
	separator       = string(os.PathSeparator)
	mutex           = &sync.Mutex{}
	countPerProcess = 3
)

type EventConversion struct {
	Ctx      *orm.DataContext
	FilePath string
}

func NewEventConversion(ctx *orm.DataContext, filePath string) *EventConversion {
	ev := new(EventConversion)
	ev.Ctx = ctx
	ev.FilePath = filePath

	return ev
}

func fileExists(fileLocation string) bool {
	if _, err := os.Stat(fileLocation); err == nil {
		return true
	}

	return false
}

func (ev *EventConversion) Run() {
	var wg sync.WaitGroup

	if ev.FilePath != "" {
		if fileExists(ev.FilePath) {
			files, err := ioutil.ReadDir(ev.FilePath)
			if err != nil {
				tk.Println(err)
			}

			counter := 0
			countData := len(files)
			isDone := false
			startIndex := 0
			endIndex := 0

			for !isDone {
				startIndex = counter * countPerProcess
				endIndex = (counter + 1) * countPerProcess

				if endIndex > countData {
					endIndex = countData
				}

				fileprocess := files[startIndex:endIndex]
				for _, f := range fileprocess {
					if strings.Contains(f.Name(), ".xlsx") && !strings.Contains(f.Name(), "~") {
						wg.Add(1)
						go ev.processFile(f.Name(), &wg, ev.Ctx)
					}
				}
				wg.Wait()

				counter++

				if endIndex >= countData {
					isDone = true
				}
			}
		}
	}
}

func (ev *EventConversion) processFile(filename string, wg *sync.WaitGroup, ctx *orm.DataContext) {
	// mutex.Lock()

	now := time.Now()
	tk.Println("Starting process file ", filename)

	fLoc := ev.FilePath + separator + filename
	fi, err := os.Stat(fLoc)
	turbine := strings.Replace(fi.Name(), ".xlsx", "", 1)
	project := "Tejuva"

	xls, err := xlsx.OpenFile(fLoc)
	if err != nil {
		tk.Println("Error open excel file : ", err.Error())
	}

	for _, sheet := range xls.Sheet {
		dataAlarm := tk.M{}
		id := ""
		cekAlarms := tk.M{}
		counterDowntime := 0
		dataDetail := make([]tk.M, 0)
		counterDetail := 0
		for idx, row := range sheet.Rows {
			if idx > 0 {
				affectedItem, _ := row.Cells[3].String()
				if strings.TrimSpace(affectedItem) != "" {
					alarmIdx := 0
					alarmDesc := affectedItem
					if strings.Contains(affectedItem, ")") {
						alarms := strings.Split(affectedItem, ")")
						if len(alarms) > 0 {
							alarmIdx = tk.ToInt(strings.TrimSpace(strings.Replace(alarms[0], "(", "", 1)), "0")
							if len(alarms) > 1 {
								alarmDesc = strings.TrimSpace(alarms[1])
							}
						}
					}
					// tk.Printf("alarmIdx: %v \n", alarmIdx)
					alarmbrakes := make([]AlarmBrake, 0)
					/*csr, err := ctx.Connection.NewQuery().From(new(AlarmBrake).TableName()).
					Where(dbox.Eq("alarmindex", alarmIdx)).Cursor(nil)*/

					match := tk.M{}
					match = tk.M{"alarmindex": alarmIdx}

					pipes := []tk.M{}
					pipes = append(pipes, tk.M{"$match": match})

					csr, err := ctx.Connection.NewQuery().From(new(AlarmBrake).TableName()).
						Command("pipe", pipes).Cursor(nil)

					if err != nil {
						tk.Println("Error getting alarm brake: " + err.Error())
					} else {
						err = csr.Fetch(&alarmbrakes, 0, false)
						if err != nil {
							tk.Println("Error fetch data: " + err.Error())
						}
						csr.Close()

						brakeProgram := 0
						brakeType := ""
						if len(alarmbrakes) > 0 {
							brakeProgram = alarmbrakes[0].BrakeProgram
							brakeType = alarmbrakes[0].Type
						}

						sTurbineStatus, _ := row.Cells[7].String()
						turbineStatus := ""
						if sTurbineStatus != "" {
							arrTurbineStatus := strings.Split(sTurbineStatus, " ")
							if len(arrTurbineStatus) > 1 {
								turbineStatus = arrTurbineStatus[1]
							}
						}

						rawdata := new(DowntimeEventRaw).New()
						rawdata.ProjectName = project
						rawdata.Turbine = turbine

						sTimeStamp, _ := row.Cells[0].String()
						rawdata.TimeStamp, _ = time.Parse("2006-01-02 15:04:05", strings.Replace(strings.Replace(sTimeStamp, "T", " ", 1), "+05:30", "", 1))

						sEventType, _ := row.Cells[1].String()
						rawdata.EventType = strings.TrimSpace(sEventType)

						rawdata.BrakeProgram = brakeProgram
						rawdata.DateInfo = GetDateInfo(rawdata.TimeStamp)
						rawdata.AlarmDescription = alarmDesc
						rawdata.AlarmId = alarmIdx
						rawdata.TurbineStatus = strings.TrimSpace(turbineStatus)
						rawdata.BrakeType = brakeType
						sAlarmToggle, _ := row.Cells[2].String()

						if strings.TrimSpace(strings.ToUpper(sAlarmToggle)) == "TRUE" || strings.TrimSpace(strings.ToUpper(sAlarmToggle)) == "1" {
							rawdata.AlarmToggle = true
						} else {
							rawdata.AlarmToggle = false
						}

						ctx.Insert(rawdata)

						////////////////////////////////// PROCESSING DOWNTIME EVENT /////////////////////////////////////////
						if strings.TrimSpace(sEventType) == "alarmchanged" && brakeProgram > 0 {
							if len(cekAlarms.Keys()) <= 0 && rawdata.AlarmToggle {
								counterDowntime++
								id = tk.ToString(counterDowntime)
								dataAlarm.Set(id, tk.M{}.
									Set("timestart", rawdata.TimeStamp).
									Set("project", project).
									Set("turbine", turbine).
									Set("alarmdescription", rawdata.AlarmDescription).
									Set("braketype", rawdata.BrakeType))
							}

							cekAlarms.Set(tk.ToString(rawdata.AlarmId), rawdata.AlarmToggle)

							// set detail downtime event
							counterDetail++
							detail := tk.M{}.
								Set("timestamp", rawdata.TimeStamp).
								Set("alarmid", rawdata.AlarmId).
								Set("alarmdescription", rawdata.AlarmDescription).
								Set("alarmtoggle", rawdata.AlarmToggle)
							dataDetail = append(dataDetail, detail)

							if isAllClosed(cekAlarms) && len(dataDetail) > 1 {
								// insert downtime event into database here
								if dataAlarm.Has(id) {
									getdata := dataAlarm.Get(id).(tk.M)
									timestart := getdata.Get("timestart").(time.Time)
									duration := rawdata.TimeStamp.Sub(timestart)
									dataAlarm.Set(id, getdata.
										Set("timeend", rawdata.TimeStamp).
										Set("duration", duration.Seconds()).
										Set("detail", dataDetail))

									latestData := dataAlarm.Get(id).(tk.M)

									downEvent := new(DowntimeEvent).New()
									downEvent.ProjectName = latestData.GetString("project")
									downEvent.Turbine = latestData.GetString("turbine")
									downEvent.AlarmDescription = latestData.GetString("alarmdescription")
									downEvent.Duration = latestData.GetFloat64("duration")
									downEvent.TimeStart = latestData.Get("timestart").(time.Time)
									downEvent.DateInfoStart = GetDateInfo(downEvent.TimeStart)
									downEvent.TimeEnd = latestData.Get("timeend").(time.Time)
									downEvent.DateInfoEnd = GetDateInfo(downEvent.TimeEnd)

									brakeType := latestData.GetString("braketype")

									if strings.Contains(strings.ToLower(brakeType), "grid") {
										downEvent.DownGrid = true
									}

									if strings.Contains(strings.ToLower(brakeType), "environment") {
										downEvent.DownEnvironment = true
									}

									if !strings.Contains(strings.ToLower(brakeType), "grid") && !strings.Contains(strings.ToLower(brakeType), "environment") {
										downEvent.DownMachine = true
									}

									downDetail := make([]DowntimeEventDetail, 0)
									details := latestData.Get("detail").([]tk.M)
									for _, d := range details {
										var dd DowntimeEventDetail
										dd.AlarmDescription = d.GetString("alarmdescription")
										dd.AlarmId = d.GetInt("alarmid")
										dd.TimeStamp = d.Get("timestamp").(time.Time)
										dd.AlarmToggle = d.Get("alarmtoggle").(bool)
										dd.DateInfo = GetDateInfo(dd.TimeStamp)

										downDetail = append(downDetail, dd)
									}

									downEvent.Detail = downDetail

									// insert data
									tk.Println("Inserting data downtime event")
									ctx.Insert(downEvent)

									// reset cek alarms
									cekAlarms = tk.M{}
									counterDetail = 0
									id = ""
									dataDetail = make([]tk.M, 0)
								}
							} else {
								if len(dataDetail) <= 1 && !rawdata.AlarmToggle {
									dataDetail = make([]tk.M, 0)
								}
							}
						}
					}

					////////////////////////////////// PROCESSING DOWNTIME EVENT /////////////////////////////////////////
				}
			}
		}
	}

	duration := time.Now().Sub(now)
	tk.Println(tk.Sprintf("Process file %v about %v sec(s)", filename, duration.Seconds()))

	// mutex.Unlock()

	wg.Done()
}

func isAllClosed(alarms tk.M) bool {
	isclosed := true
	if len(alarms.Keys()) > 0 {
		for _, a := range alarms {
			if a.(bool) {
				isclosed = false
				break
			}
		}
	} else {
		isclosed = false
	}

	return isclosed
}
