package conversion

import (
	"log"
	"os"
	"strings"
	"sync"
	"time"

	// . "eaciit/wfdemo-git/library/helper"
	. "eaciit/wfdemo-git/library/models"

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

type DownConversion struct {
	Ctx      *orm.DataContext
	FilePath string
}

func NewDownConversion(ctx *orm.DataContext, filePath string) *DownConversion {
	ev := new(DownConversion)
	ev.Ctx = ctx
	ev.FilePath = filePath

	return ev
}

func (ev *DownConversion) Run() {
	// _ = ev.getLatest()
	var wg sync.WaitGroup
	loops := ev.getLatest()
	counter := 0
	for _, loop := range loops {
		counter++
		// if loop.Turbine == "SSE017" {
		// log.Printf("loop: %v | %v \n", loop.Turbine, loop.LatestProcessTime)
		wg.Add(1)
		go ev.processTurbine(loop, &wg)
		// }

		if counter%5 == 0 || len(loops) == counter {
			wg.Wait()
		}
	}
}

func (ev *DownConversion) processTurbine(loop GroupResult, wg *sync.WaitGroup) {
	// mutex.Lock()

	now := time.Now()
	log.Printf("Starting process %v | %v | %v \n", loop.Project, loop.Turbine, loop.LatestProcessTime.String())

	pipes := []tk.M{}

	match := tk.M{
		"projectname":  loop.Project,
		"turbine":      loop.Turbine,
		"eventtype":    "alarmchanged",
		"brakeprogram": tk.M{"$gt": 0},
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
		log.Println("Error: " + err.Error())
	} else {
		err = csr.Fetch(&eventRaws, 0, false)
		if err != nil {
			log.Println("Error: " + err.Error())
		} else {

			loopData := eventRaws

			dataInserted := 0

		mainLoop:
			for {
				// log.Printf("loopData: %v \n", len(loopData))
				trueFound := map[int]EventDownDetail{}
				details := []EventDownDetail{}
				_ = details
				startIdx := -1
				endIdx := -1
				foundInProduction := false

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
					// log.Printf("\n\nData: %v | %v | %v | %v | %v \n", data.AlarmToggle, data.AlarmDescription, data.TurbineStatus, data.TimeStamp.String(), data.AlarmId)
					if data.DateInfo.MonthId != 0 {
						tmp := EventDownDetail{}
						tmp.AlarmId = data.AlarmId
						tmp.AlarmToggle = data.AlarmToggle
						tmp.TimeStamp = data.TimeStamp.UTC()
						tmp.TimeStampInt = data.TimeStampInt
						// tmp.TimeStampUTC = data.TimeStampUTC
						tmp.DateInfo = data.DateInfo
						// tmp.DateInfoUTC = data.DateInfoUTC
						tmp.AlarmDescription = data.AlarmDescription

						// add by ams, regarding to add new req | 20170130
						tmp.BrakeType = data.BrakeType

						// log.Printf("trueFound: %v | %#v \n", idx, len(trueFound))

						var isTrue bool

						if len(trueFound) > 0 {
							isStartDone := false
							if trueFound[start.AlarmId].TimeStampInt == 0 {
								isStartDone = true
							}

							if start.AlarmDescription != data.AlarmDescription && isStartDone {
								isTrue = true
							}

							// log.Printf("condition: %v | %v || %#v \n", tk.IsNilOrEmpty(trueFound[start.AlarmId]), start.AlarmDescription != data.AlarmDescription, trueFound[start.AlarmId].TimeStampInt)
						}

						if ev.isProduction(data.TurbineStatus) && isTrue {
							// log.Printf("Production: %#v \n", data)
							trueFound = map[int]EventDownDetail{}
							// log.Printf("trueFoundXXXXXX: %v | %#v \n", idx, len(trueFound))
							end = loopData[idx-1]

							productionCondition := ev.getTurbineProduction(loop, start.TimeStamp.UTC())
							if end.TimeStamp.UTC().Sub(productionCondition.TimeStamp.UTC()).Seconds() > 0.0 && productionCondition.TimeStamp.UTC().Sub(start.TimeStamp.UTC()).Seconds() != 0.0 {
								end = productionCondition
								foundInProduction = true
							}

							ev.insertEventDown(loop, start, end)
							dataInserted++

							details = []EventDownDetail{}
							endIdx = idx - 1
							break reloop

						} else if data.AlarmToggle {
							// log.Printf("x: %v \n", data.AlarmId)

							if startIdx == -1 {
								startIdx = idx
								start = data
							}

							trueFound[data.AlarmId] = tmp
							// log.Printf("n: %v \n", trueFound[data.AlarmId].DateInfo.MonthId)
						} else if !data.AlarmToggle && trueFound[data.AlarmId].DateInfo.MonthId != 0 {
							// log.Printf("y: %v \n", data.AlarmId)
							// log.Printf("y: %v \n", trueFound[data.AlarmId])

							// details = append(details, trueFound[data.AlarmId])
							// details = append(details, tmp)

							tmpFound := map[int]EventDownDetail{}

							for id, found := range trueFound {
								if id != tmp.AlarmId {
									tmpFound[id] = found
								}
							}
							trueFound = tmpFound
							// log.Printf("trueFoundXXXXXX: %v | %#v \n", idx, len(trueFound))
							if len(trueFound) == 0 || trueFound == nil {
								end = data

								productionCondition := ev.getTurbineProduction(loop, start.TimeStamp.UTC())
								if end.TimeStamp.UTC().Sub(productionCondition.TimeStamp.UTC()).Seconds() > 0.0 && productionCondition.TimeStamp.UTC().Sub(start.TimeStamp.UTC()).Seconds() != 0.0 {
									end = productionCondition
									foundInProduction = true
								}

								ev.insertEventDown(loop, start, end)
								dataInserted++

								details = []EventDownDetail{}
								endIdx = idx
								break reloop
							}
						}
					}
				}

				// log.Printf("loopData: %v \n", len(loopData))

				if !foundInProduction {

					tmpLoopData := []EventRaw{}

					if endIdx > 0 {
						tmpLoopData = append(tmpLoopData, loopData[endIdx+1:]...)
					}

					loopData = tmpLoopData
				} else {
					loopData = ev.getLoopData(loop, end)
				}

				if len(loopData) == 0 && dataInserted > 0 {
					break mainLoop
				} else if len(loopData) == 0 && dataInserted == 0 {
					end = ev.getTurbineProduction(loop, start.TimeStamp.UTC())

					ev.insertEventDown(loop, start, end)
					loopData = ev.getLoopData(loop, end)

					details = []EventDownDetail{}
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
	// csr.Close()
	wg.Done()
}

func (ev *DownConversion) getLatest() []GroupResult {
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
		From(new(EventDown).TableName()).
		Command("pipe", pipes).
		Cursor(nil)

	defer csr.Close()

	eventDowns := []tk.M{}

	if err != nil {
		log.Println("Error: " + err.Error())
		return nil
	}
	err = csr.Fetch(&eventDowns, 0, false)

	if err != nil {
		log.Println("Error: " + err.Error())
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
	match = tk.M{"eventtype": "alarmchanged", "brakeprogram": tk.M{"$gt": 0}}

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
		log.Println("Error: " + err.Error())
		return nil
	}
	err = csr.Fetch(&eventRaws, 0, false)

	if err != nil {
		log.Println("Error: " + err.Error())
		return nil
	}

	/*for _, res := range result {
		log.Printf("res: %v | %v \n", res.Turbine, res.LatestProcessTime)
	}

	log.Println()*/

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

	/*for _, res := range result {
		log.Printf("res: %v | %v \n", res.Turbine, res.LatestProcessTime)
	}*/

	return result
}

func (ev *DownConversion) isProduction(check string) (status bool) {
	strList := []string{"Production", "Boot", "Start", "Waiting", "LimSw", "Pitch", "Anemometer", "Accu", "Slow", "Syncron.", "Fast", "Turb."}

	for _, b := range strList {
		if b == strings.ToLower(check) {
			return true
		}
	}
	return false
}

func (ev *DownConversion) getTurbineProduction(loop GroupResult, startTime time.Time) (result EventRaw) {
	match := tk.M{
		"projectname":   loop.Project,
		"turbine":       loop.Turbine,
		"eventtype":     "turbinestatechanged",
		"brakeprogram":  tk.M{"$eq": 0},
		"turbinestatus": tk.M{"$in": []string{"Production", "Boot", "Start", "Waiting", "LimSw", "Pitch", "Anemometer", "Accu", "Slow", "Syncron.", "Fast", "Turb."}},
	}

	// if loop.LatestFrom == "Raw" {
	match.Set("timestamp", tk.M{"$gte": startTime})
	// } else {
	// 	match.Set("timestamp", tk.M{"$gt": loop.LatestProcessTime})
	// }

	pipes := []tk.M{}
	pipes = append(pipes, tk.M{"$match": match})
	pipes = append(pipes, tk.M{"$sort": tk.M{"timestamp": 1}})
	pipes = append(pipes, tk.M{"$limit": 1})

	csr, err := ev.Ctx.Connection.NewQuery().From(new(EventRaw).TableName()).Command("pipe", pipes).Cursor(nil)
	defer csr.Close()

	eventRaws := []EventRaw{}

	if err != nil {
		log.Println("Error: " + err.Error())
	} else {
		err = csr.Fetch(&eventRaws, 0, false)
		if len(eventRaws) > 0 {
			result = eventRaws[0]
		}
	}

	return
}

func (ev *DownConversion) getLoopData(loop GroupResult, end EventRaw) (eventRaws []EventRaw) {
	pipes := []tk.M{}
	match := tk.M{
		"projectname":  loop.Project,
		"turbine":      loop.Turbine,
		"eventtype":    "alarmchanged",
		"brakeprogram": tk.M{"$gt": 0},
	}

	match.Set("timestamp", tk.M{"$gte": end.TimeStamp.UTC()})

	pipes = append(pipes, tk.M{"$match": match})
	pipes = append(pipes, tk.M{"$sort": tk.M{"timestamp": 1}})

	csrx, err := ev.Ctx.Connection.NewQuery().From(new(EventRaw).TableName()).Command("pipe", pipes).Cursor(nil)
	defer csrx.Close()

	if err != nil {
		log.Println("Error: " + err.Error())
	} else {
		err = csrx.Fetch(&eventRaws, 0, false)
		if err != nil {
			log.Println("Error: " + err.Error())
		}
	}

	return
}

func (ev *DownConversion) insertEventDown(loop GroupResult, start EventRaw, end EventRaw) {
	down := new(EventDown).New()

	down.ProjectName = loop.Project
	down.Turbine = loop.Turbine

	down.TimeStart = start.TimeStamp.UTC()
	down.TimeStartInt = start.TimeStampInt
	down.DateInfoStart = start.DateInfo

	down.TimeEnd = end.TimeStamp.UTC()
	down.TimeEndInt = end.TimeStampInt
	down.DateInfoEnd = end.DateInfo

	down.AlarmDescription = start.AlarmDescription

	// add by ams, regarding to add new req | 20170130
	down.BrakeType = start.BrakeType

	down.Duration = end.TimeStamp.UTC().Sub(start.TimeStamp.UTC()).Seconds()

	// down.Detail = details

	if down.DateInfoStart.MonthId != 0 && down.TimeStart.UTC().Year() != 1 {
		mutex.Lock()
		brakeType := start.BrakeType
		if strings.Contains(strings.ToLower(brakeType), "grid") {
			down.DownGrid = true
		} else if strings.Contains(strings.ToLower(brakeType), "environment") {
			down.DownEnvironment = true
		} else if !strings.Contains(strings.ToLower(brakeType), "grid") && !strings.Contains(strings.ToLower(brakeType), "environment") {
			down.DownMachine = true
		}

		down := down.New()
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
		// log.Print("Insert Event Down")
	}
}
