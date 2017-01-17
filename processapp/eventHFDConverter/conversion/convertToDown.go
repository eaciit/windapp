package conversion

import (
	"log"
	"os"
	"strings"
	"sync"
	"time"

	// . "eaciit/wfdemo-git/library/helper"
	. "eaciit/wfdemo-git/library/models"

	"github.com/eaciit/dbox"
	_ "github.com/eaciit/dbox/dbc/mongo"
	"github.com/eaciit/orm"

	tk "github.com/eaciit/toolkit"
)

var (
	separator       = string(os.PathSeparator)
	mutex           = &sync.Mutex{}
	countPerProcess = 1
)

type GroupResult struct {
	Project           string
	Turbine           string
	LatestProcessTime time.Time
	LatestFrom        string
}

type HFDDownConversion struct {
	Ctx *orm.DataContext
}

func NewHFDDownConversion(ctx *orm.DataContext) *HFDDownConversion {
	ev := new(HFDDownConversion)
	ev.Ctx = ctx
	return ev
}

func (ev *HFDDownConversion) Run() {
	var wg sync.WaitGroup
	loops := ev.getLatest()

	for _, loop := range loops {
		// if loop.Turbine == "SSE017" {
		// log.Printf("loop: %v | %v \n", loop.Turbine, loop.LatestProcessTime)
		wg.Add(1)
		go ev.processTurbine(loop, &wg)
		// }
	}

	wg.Wait()
}

func (ev *HFDDownConversion) processTurbine(loop GroupResult, wg *sync.WaitGroup) {
	// mutex.Lock()

	now := time.Now()
	log.Printf("Starting process %v | %v | %v \n", loop.Project, loop.Turbine, loop.LatestProcessTime.String())

	pipes := []tk.M{}

	match := tk.M{
		"projectname": loop.Project,
		"turbine":     loop.Turbine,
		// "eventtype":    "alarmchanged",
		// "brakeprogram": tk.M{"$gt": 0},
	}

	if loop.LatestFrom == "Raw" {
		match.Set("timestamp", tk.M{"$gte": loop.LatestProcessTime})
	} else {
		match.Set("timestamp", tk.M{"$gt": loop.LatestProcessTime})
	}

	pipes = append(pipes, tk.M{"$match": match})
	pipes = append(pipes, tk.M{"$sort": tk.M{"timestamp": 1}})

	csr, err := ev.Ctx.Connection.NewQuery().From(new(EventRawHFD).TableName()).Command("pipe", pipes).Cursor(nil)
	defer csr.Close()

	eventRaws := []EventRawHFD{}

	if err != nil {
		tk.Println("Error: " + err.Error())
	} else {
		err = csr.Fetch(&eventRaws, 0, false)
		if err != nil {
			tk.Println("Error: " + err.Error())
		} else {

			loopData := eventRaws

			AlarmCode := ".AlarmCode"
			ErrorState := ".ErrorState"
			TurbineState := ".TurbineState"

		mainLoop:
			for {
				startIdx := -1
				endIdx := -1

				var start, end EventRawHFD
				lastAlarmCode := -1
				lastAlarmDesc := ""
				// lastBrakeProgram := 0

			reloop:
				for idx, data := range loopData {
					if idx > 0 {
						lastAlarmCode = loopData[idx-1].AlarmId
						lastAlarmDesc = loopData[idx-1].AlarmDescription
						// lastBrakeProgram = loopData[idx-1].BrakeProgram
					}

					if data.DateInfo.MonthId != 0 {
						/*else if lastBrakeProgram != 0 && lastBrakeProgram != 999 && lastAlarmCode != 0 && strings.Contains(data.EventType, AlarmCode) && startIdx == -1 {
							startIdx = idx
							start = data
							ev.InsertToMonitoringEvent(start, "down")
						} */

						if data.BrakeProgram != 0 && data.BrakeProgram != 999 && data.AlarmId != 0 && strings.Contains(data.EventType, AlarmCode) && startIdx == -1 {
							startIdx = idx
							start = data
							ev.InsertToMonitoringEvent(start, "down")
						} else if ((strings.Contains(data.EventType, ErrorState) && data.AlarmId == 0) || (strings.Contains(data.EventType, AlarmCode) && data.AlarmId == 0) || (strings.Contains(data.EventType, TurbineState) && (data.AlarmId >= 0 && data.AlarmId <= 11) || data.AlarmId == 0)) && startIdx != -1 {
							end = data
						} else if data.AlarmId == 999 && startIdx != -1 {
							if (strings.Contains(data.EventType, ErrorState) && lastAlarmCode == 0) || (strings.Contains(data.EventType, AlarmCode) && lastAlarmCode == 0) || (strings.Contains(data.EventType, TurbineState) && (lastAlarmCode >= 0 && lastAlarmCode <= 11) || lastAlarmCode == 0) {
								end = data
							}
						}

						if end.TimeStamp.Year() != 1 && startIdx != -1 {
							tmp := EventRawHFD{}
							tmp = end

							if tmp.AlarmId == 999 {
								tmp.AlarmId = lastAlarmCode
								tmp.AlarmDescription = lastAlarmDesc
							}

							ev.InsertToMonitoringEvent(tmp, "up")

							down := new(EventDownHFD).New()
							down.AlarmID = start.AlarmId
							down.ProjectName = loop.Project
							down.Turbine = loop.Turbine

							down.TimeStart = start.TimeStamp.UTC()
							down.DateInfoStart = start.DateInfo

							down.TimeEnd = end.TimeStamp.UTC()
							down.DateInfoEnd = end.DateInfo

							down.AlarmDescription = start.AlarmDescription
							down.Duration = end.TimeStamp.UTC().Sub(start.TimeStamp.UTC()).Seconds()

							down.GroupTimeStart = convertTo10min(down.TimeStart)
							down.GroupTimeEnd = convertTo10min(down.TimeEnd)

							if down.DateInfoStart.MonthId != 0 && down.TimeStart.UTC().Year() != 1 {
								mutex.Lock()
								brakeType := data.BrakeType
								if strings.Contains(strings.ToLower(brakeType), "grid") {
									down.DownGrid = true
								}
								if strings.Contains(strings.ToLower(brakeType), "environment") {
									down.DownEnvironment = true
								}
								if !strings.Contains(strings.ToLower(brakeType), "grid") && !strings.Contains(strings.ToLower(brakeType), "environment") {
									down.DownMachine = true
								}

								down = down.New()
								count := 0
								for {
									e := ev.Ctx.Insert(down)
									if e != nil {
										log.Printf("error: %v \n", e.Error())
										down = down.New()
									} else {
										break
									}

									if count == 2 {
										break
									}
									count++
								}

								mutex.Unlock()
							}

							endIdx = idx
							break reloop
						}

					}
				}

				// log.Printf("loopData: %v \n", len(loopData))

				tmpLoopData := []EventRawHFD{}

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

	duration := time.Now().Sub(now)
	log.Printf("Process %v | %v about %v sec(s) \n", loop.Project, loop.Turbine, duration.Seconds())
	csr.Close()
	wg.Done()
}

func (ev *HFDDownConversion) getLatest() []GroupResult {
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
		From(new(EventDownHFD).TableName()).
		Command("pipe", pipes).
		Cursor(nil)

	defer csr.Close()

	eventDowns := []tk.M{}

	if err != nil {
		tk.Println("Error: " + err.Error())
		return nil
	}
	err = csr.Fetch(&eventDowns, 0, false)

	if err != nil {
		tk.Println("Error: " + err.Error())
		return nil
	}

	/*for _, val := range eventDowns {
		id := val.Get("_id").(tk.M)
		log.Printf("%#v \n", id)
	}*/

	if len(eventDowns) > 0 {
		// log.Printf("len(eventDowns): %v \n", len(eventDowns))
		for _, val := range eventDowns {
			id := val.Get("_id").(tk.M)

			tmp := GroupResult{}
			tmp.Project = id.GetString("project")
			tmp.Turbine = id.GetString("turbine")
			tmp.LatestProcessTime = val.Get("timestamp").(time.Time).UTC()
			tmp.LatestFrom = "Down"
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
	match = tk.M{"brakeprogram": tk.M{"$gt": 0}}

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
		From(new(EventRawHFD).TableName()).
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

func convertTo10min(input time.Time) (output time.Time) {
	// THour := input.Hour()
	TMinute := input.Minute()
	TSecond := input.Second()
	TMinuteValue := float64(TMinute) + tk.Div(float64(TSecond), 60.0)
	TMinuteCategory := tk.ToInt(tk.RoundingUp64(tk.Div(TMinuteValue, 10), 0)*10, "0")

	tmpInput := input.Add(time.Duration(TMinuteCategory-TMinute) * time.Minute).Add(time.Duration(TSecond*-1) * time.Second).UTC()
	output, _ = time.Parse("20060102_150405", tmpInput.Format("20060102_150405"))
	return
}

func (ev *HFDDownConversion) InsertToMonitoringEvent(data EventRawHFD, status string) error {
	mutex.Lock()

	mEvent := new(MonitoringEvent).New()
	mEvent.Project = data.ProjectName
	mEvent.Turbine = data.Turbine
	mEvent.TimeStamp = data.TimeStamp
	mEvent.DateInfo = data.DateInfo
	mEvent.GroupTimeStamp = convertTo10min(mEvent.TimeStamp)
	mEvent.AlarmId = data.AlarmId
	mEvent.AlarmDescription = data.AlarmDescription
	mEvent.Type = "brake"
	mEvent.Status = status

	if status == "up" {
		filter := []*dbox.Filter{}
		filter = append(filter, dbox.Eq("project", data.ProjectName))
		filter = append(filter, dbox.Eq("turbine", data.Turbine))

		csr, e := ev.Ctx.Connection.NewQuery().
			From(new(MonitoringEvent).TableName()).
			Where(dbox.And(filter...)).
			Order("-timestamp").
			Take(1).
			Cursor(nil)

		defer csr.Close()

		events := []MonitoringEvent{}

		if e != nil {
			log.Printf("Error: %v \n" + e.Error())
			return e
		}
		e = csr.Fetch(&events, 0, false)

		if e != nil {
			log.Printf("Error: %v \n" + e.Error())
			return e
		}

		if len(events) == 1 {
			start := events[0]
			log.Printf("%v | %v \n", mEvent.TimeStamp.UTC().String(), start.TimeStamp.UTC().String())
			mEvent.Duration = mEvent.TimeStamp.UTC().Sub(start.TimeStamp.UTC()).Seconds()
		}
	}

	mEvent = mEvent.New()

	e := ev.Ctx.Save(mEvent)

	mutex.Unlock()

	return e
}
