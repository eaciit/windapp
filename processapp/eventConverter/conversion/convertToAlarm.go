package conversion

import (
	"log"
	"sync"
	"time"

	// . "eaciit/wfdemo-git/library/helper"
	. "eaciit/wfdemo-git/library/models"

	_ "github.com/eaciit/dbox/dbc/mongo"
	"github.com/eaciit/orm"

	"math"

	tk "github.com/eaciit/toolkit"
)

type AlarmConversion struct {
	Ctx      *orm.DataContext
	FilePath string
}

func NewAlarmConversion(ctx *orm.DataContext, filePath string) *AlarmConversion {
	ev := new(AlarmConversion)
	ev.Ctx = ctx
	ev.FilePath = filePath

	return ev
}

func (ev *AlarmConversion) Run() {
	// _ = ev.getLatest()
	var wg sync.WaitGroup
	loops := ev.getLatest()
	maxRun := 2

	for idx, loop := range loops {
		// if loop.Turbine == "SSE017" {
		// log.Printf("loop: %v | %v \n", loop.Turbine, loop.LatestProcessTime)
		wg.Add(1)
		go ev.processTurbine(loop, &wg)
		// }

		div := math.Mod(tk.ToFloat64(idx, 0, tk.RoundingAuto), tk.ToFloat64(maxRun, 0, tk.RoundingAuto))

		if div == 0 || idx == len(loops)-1 {
			wg.Wait()
		}
	}
}

func (ev *AlarmConversion) processTurbine(loop GroupResult, wg *sync.WaitGroup) {
	// mutex.Lock()

	now := time.Now()
	log.Printf("Starting process %v | %v | %v \n", loop.Project, loop.Turbine, loop.LatestProcessTime.String())

	pipes := []tk.M{}

	match := tk.M{
		"projectname": loop.Project,
		"turbine":     loop.Turbine,
		"eventtype":   "alarmchanged",
		// "eventtype":    tk.M{"$ne": "000"},
		"brakeprogram": 0,
	}

	if loop.LatestFrom == "Raw" {
		match.Set("timestamp", tk.M{"$gte": loop.LatestProcessTime})
	} else {
		match.Set("timestamp", tk.M{"$gt": loop.LatestProcessTime})
	}

	pipes = append(pipes, tk.M{"$match": match})
	pipes = append(pipes, tk.M{"$sort": tk.M{"timestamp": 1}})

	csr, err := ev.Ctx.Connection.NewQuery().From(new(EventRaw).TableName()).Command("pipe", pipes).Cursor(nil)
	defer csr.Close()

	eventRaws := []EventRaw{}

	if err != nil {
		tk.Println("Error: " + err.Error())
	} else {
		err = csr.Fetch(&eventRaws, 0, false)
		if err != nil {
			tk.Println("Error: " + err.Error())
		} else {

			loopData := eventRaws

		mainLoop:
			for {
				// log.Printf("loopData: %v \n", len(loopData))
				trueFound := map[int]EventAlarmDetail{}
				details := []EventAlarmDetail{}
				startIdx := -1
				endIdx := -1

				var start, end EventRaw
				// log.Printf("loopData: %#v \n", loopData)
				// log.Printf("trueFound: %#v \n", trueFound[1].DateInfo.MonthId)

				/*for idx, data := range loopData {
					log.Printf("loopData: %v | %#v \n", idx, data)
				}*/

			reloop:
				for idx, data := range loopData {
					// log.Printf("data: %v | %v | %v \n", data.TimeStamp.UTC(), data.AlarmToggle, data.AlarmId)
					// log.Printf("loopData: %v \n", len(loopData))
					// log.Printf("trueFound: %v | %#v \n", idx, len(trueFound))
					if data.DateInfo.MonthId != 0 {

						tmp := EventAlarmDetail{}
						tmp.AlarmId = data.AlarmId
						tmp.AlarmToggle = data.AlarmToggle
						tmp.TimeStamp = data.TimeStamp.UTC()
						tmp.TimeStampInt = data.TimeStampInt
						// tmp.TimeStampUTC = data.TimeStampUTC
						tmp.DateInfo = data.DateInfo
						// tmp.DateInfoUTC = data.DateInfoUTC
						tmp.AlarmDescription = data.AlarmDescription

						// log.Printf("trueFound: %v | %#v \n", idx, len(trueFound))

						if data.AlarmToggle {
							//log.Printf("x: %v \n", data.AlarmId)

							if startIdx == -1 {
								startIdx = idx
								start = data
							}

							trueFound[data.AlarmId] = tmp
							// log.Printf("n: %v \n", trueFound[data.AlarmId].DateInfo.MonthId)
						} else if !data.AlarmToggle && trueFound[data.AlarmId].DateInfo.MonthId != 0 {
							// log.Printf("y: %v \n", data.AlarmId)
							// log.Printf("y: %v \n", trueFound[data.AlarmId])

							details = append(details, trueFound[data.AlarmId])
							details = append(details, tmp)

							tmpFound := map[int]EventAlarmDetail{}

							for id, found := range trueFound {
								if id != tmp.AlarmId {
									tmpFound[id] = found
								}
							}
							trueFound = tmpFound
							// log.Printf("trueFoundXXXXXX: %v | %#v \n", idx, len(trueFound))
							if len(trueFound) == 0 || trueFound == nil {
								end = data

								alarm := new(EventAlarm).New()

								alarm.ProjectName = loop.Project
								alarm.Turbine = loop.Turbine

								alarm.TimeStart = start.TimeStamp.UTC()
								alarm.TimeStartInt = start.TimeStampInt
								alarm.DateInfoStart = start.DateInfo

								alarm.TimeEnd = end.TimeStamp.UTC()
								alarm.TimeEndInt = end.TimeStampInt
								alarm.DateInfoEnd = end.DateInfo

								alarm.AlarmDescription = start.AlarmDescription
								alarm.Duration = end.TimeStamp.UTC().Sub(start.TimeStamp.UTC()).Seconds()

								alarm.Detail = details

								if alarm.DateInfoStart.MonthId != 0 && alarm.TimeStart.UTC().Year() != 1 {
									mutex.Lock()
									/*brakeType := data.BrakeType
									if strings.Contains(strings.ToLower(brakeType), "grid") {
										alarm.DownGrid = true
									}
									if strings.Contains(strings.ToLower(brakeType), "environment") {
										alarm.DownEnvironment = true
									}
									if !strings.Contains(strings.ToLower(brakeType), "grid") && !strings.Contains(strings.ToLower(brakeType), "environment") {
										alarm.DownMachine = true
									}*/

									alarm := alarm.New()
									count := 0
									for {
										e := ev.Ctx.Insert(alarm)
										if e != nil {
											log.Printf("error: %v \n", e.Error())
											alarm = alarm.New()
										} else {
											break
										}

										if count == 2 {
											break
										}
										count++
									}

									mutex.Unlock()
									// log.Print("Insert Event Down")
								}

								details = []EventAlarmDetail{}
								endIdx = idx
								break reloop
							}
						}
					}
				}

				// log.Printf("loopData: %v \n", len(loopData))

				tmpLoopData := []EventRaw{}

				if endIdx > 0 {
					tmpLoopData = append(tmpLoopData, loopData[endIdx+1:]...)
				}

				loopData = tmpLoopData

				if len(loopData) == 0 {
					break mainLoop
				}
			}
		}
	}

	/*for idx, data := range eventRaws {
		log.Printf("idx: %v | %v | %v | %v \n", idx, data.AlarmId, data.TimeStamp.UTC(), data.AlarmToggle)
	}*/

	duration := time.Now().Sub(now)
	log.Printf("Process %v | %v about %v sec(s) \n", loop.Project, loop.Turbine, duration.Seconds())
	// mutex.Unlock()
	csr.Close()
	wg.Done()
}

func (ev *AlarmConversion) getLatest() []GroupResult {
	pipes := []tk.M{}
	result := []GroupResult{}
	// get max from down
	// loop check max, if not exist then check min from eventraw

	group := tk.M{
		"_id": tk.M{
			"project": "$projectname",
			"turbine": "$turbine",
		},
		"timestamp": tk.M{"$max": "$timeend"},
	}

	pipes = append(pipes, tk.M{"$group": group})

	csr, err := ev.Ctx.Connection.NewQuery().
		From(new(EventAlarm).TableName()).
		Command("pipe", pipes).
		Cursor(nil)

	defer csr.Close()

	eventAlarms := []tk.M{}

	if err != nil {
		tk.Println("Error: " + err.Error())
		return nil
	}
	err = csr.Fetch(&eventAlarms, 0, false)

	if err != nil {
		tk.Println("Error: " + err.Error())
		return nil
	}

	/*for _, val := range eventAlarms {
		id := val.Get("_id").(tk.M)
		log.Printf("%#v \n", id)
	}*/

	if len(eventAlarms) > 0 {
		// log.Printf("len(eventAlarms): %v \n", len(eventAlarms))
		for _, val := range eventAlarms {
			id := val.Get("_id").(tk.M)

			tmp := GroupResult{}
			tmp.Project = id.GetString("project")
			tmp.Turbine = id.GetString("turbine")
			tmp.LatestProcessTime = val.Get("timestamp").(time.Time).UTC()
			tmp.LatestFrom = "Alarm"
			result = append(result, tmp)
		}
	}

	turbines := []string{}

	for _, val := range result {
		turbines = append(turbines, val.Turbine)
	}

	// check min from eventraw

	match := tk.M{}
	pipes = []tk.M{}
	match = tk.M{"eventtype": "alarmchanged", "brakeprogram": 0}

	if len(turbines) > 0 {
		// checking new turbine that not in eventdown yet
		match.Set("turbine", tk.M{"$nin": turbines})
	}

	pipes = append(pipes, tk.M{"$match": match})

	group = tk.M{
		"_id": tk.M{
			"project": "$projectname",
			"turbine": "$turbine",
		},
		"timestamp": tk.M{"$min": "$timestamp"},
	}

	pipes = append(pipes, tk.M{"$group": group})

	csr, err = ev.Ctx.Connection.NewQuery().
		From(new(EventRaw).TableName()).
		Command("pipe", pipes).
		Cursor(nil)

	defer csr.Close()

	eventRaws := []tk.M{}

	if err != nil {
		tk.Println("Error: " + err.Error())
		return nil
	}
	err = csr.Fetch(&eventRaws, 0, false)

	if err != nil {
		tk.Println("Error: " + err.Error())
		return nil
	}

	/*for _, res := range result {
		log.Printf("res: %v | %v \n", res.Turbine, res.LatestProcessTime)
	}

	tk.Println()*/

	// log.Printf("len(eventRaws): %v \n", len(eventRaws))
	for _, val := range eventRaws {
		id := val.Get("_id").(tk.M)
		tmp := GroupResult{}
		tmp.Project = id.GetString("project")
		tmp.Turbine = id.GetString("turbine")
		tmp.LatestProcessTime = val.Get("timestamp").(time.Time).UTC()
		tmp.LatestFrom = "Raw"
		result = append(result, tmp)
	}
	csr.Close()

	/*for _, res := range result {
		log.Printf("res: %v | %v \n", res.Turbine, res.LatestProcessTime)
	}*/

	return result
}
