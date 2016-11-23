package conversion

import (
	"log"
	"os"
	"strings"
	"sync"
	"time"

	. "eaciit/wfdemo-git/library/helper"
	. "eaciit/wfdemo-git/library/models"

	_ "github.com/eaciit/dbox/dbc/mongo"
	"github.com/eaciit/orm"

	tk "github.com/eaciit/toolkit"
)

var (
	separator       = string(os.PathSeparator)
	mutex           = &sync.Mutex{}
	countPerProcess = 3
)

type GroupResult struct {
	Project string
	Turbine string
	Min     time.Time
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
	var wg sync.WaitGroup
	loops := ev.getLatest()

	for _, loop := range loops {
		wg.Add(1)
		go ev.processTurbine(loop, &wg)
	}

	wg.Wait()
}

func (ev *DownConversion) processTurbine(loop GroupResult, wg *sync.WaitGroup) {
	// mutex.Lock()

	now := time.Now()
	log.Printf("Starting process %v | %v | %v \n", loop.Project, loop.Turbine, loop.Min.UTC().String())

	pipes := []tk.M{}

	match := tk.M{
		"projectname":  loop.Project,
		"turbine":      loop.Turbine,
		"eventtype":    "alarmchanged",
		"timestamp":    tk.M{"$gte": loop.Min},
		"brakeprogram": tk.M{"$gt": 0},
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

			// log.Printf("loopData: %v \n", len(loopData))

		mainLoop:
			for {
				trueFound := map[int]EventDownDetail{}
				details := []EventDownDetail{}
				startIdx := -1
				endIdx := -1

				var start, end EventRaw

			reloop:
				for idx, data := range loopData {

					if data.DateInfo.MonthId != 0 {

						tmp := EventDownDetail{}
						tmp.AlarmId = data.AlarmId
						tmp.AlarmToggle = data.AlarmToggle
						tmp.TimeStamp = data.TimeStamp
						tmp.DateInfo = data.DateInfo
						tmp.AlarmDescription = data.AlarmDescription

						// log.Printf("trueFound: %v | %#v \n", idx, len(trueFound))

						if data.AlarmToggle {
							//log.Printf("x: %v \n", data.AlarmId)

							if startIdx == -1 {
								startIdx = idx
								start = data
							}

							trueFound[data.AlarmId] = tmp
						} else if !data.AlarmToggle && !tk.IsNilOrEmpty(trueFound[data.AlarmId]) {
							// log.Printf("y: %v \n", data.AlarmId)

							details = append(details, trueFound[data.AlarmId])
							details = append(details, tmp)

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

								down := new(EventDown).New()

								down.ProjectName = loop.Project
								down.Turbine = loop.Turbine

								down.TimeStart = start.TimeStamp
								down.DateInfoStart = GetDateInfo(start.TimeStamp)

								down.TimeEnd = end.TimeStamp
								down.DateInfoEnd = GetDateInfo(end.TimeStamp)

								down.AlarmDescription = start.AlarmDescription
								down.Duration = end.TimeStamp.Sub(start.TimeStamp).Seconds()

								down.Detail = details

								if down.DateInfoStart.MonthId != 0 && down.TimeStart.Year() != 1 {
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

									down := down.New()

									ev.Ctx.Insert(down)
									mutex.Unlock()
									// log.Print("Insert Event Down")
								}

								details = []EventDownDetail{}
								endIdx = idx
								break reloop
							}
						}
					}
				}

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
		"min_timestamp": tk.M{"$min": "$timeend"},
	}

	pipes = append(pipes, tk.M{"$group": group})

	csr, err := ev.Ctx.Connection.NewQuery().
		From(new(EventDown).TableName()).
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

	for _, val := range eventDowns {
		id := val.Get("_id").(tk.M)
		log.Printf("%#v \n", id)
	}

	if len(eventDowns) > 0 {
		for _, val := range eventDowns {
			id := val.Get("_id").(tk.M)

			tmp := GroupResult{}
			tmp.Project = id.GetString("project")
			tmp.Turbine = id.GetString("turbine")
			tmp.Min = val.Get("min_timestamp").(time.Time)

			result = append(result, tmp)
		}
	} else {

		// check min from eventraw
		match := tk.M{}
		pipes = []tk.M{}
		match = tk.M{"eventtype": "alarmchanged", "brakeprogram": tk.M{"$gt": 0}}
		pipes = append(pipes, tk.M{"$match": match})

		group = tk.M{
			"_id": tk.M{
				"project": "$projectname",
				"turbine": "$turbine",
			},
			"min_timestamp": tk.M{"$min": "$timestamp"},
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

		for _, val := range eventRaws {
			id := val.Get("_id").(tk.M)

			tmp := GroupResult{}
			tmp.Project = id.GetString("project")
			tmp.Turbine = id.GetString("turbine")
			tmp.Min = val.Get("min_timestamp").(time.Time)

			result = append(result, tmp)
		}
	}
	return result
}
