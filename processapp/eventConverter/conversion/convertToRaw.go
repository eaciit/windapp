package conversion

import (
	"io/ioutil"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	. "eaciit/wfdemo-git-dev/library/helper"
	. "eaciit/wfdemo-git-dev/library/models"

	"strconv"

	_ "github.com/eaciit/dbox/dbc/mongo"
	"github.com/eaciit/orm"
	tk "github.com/eaciit/toolkit"
	"github.com/tealeg/xlsx"
)

var (
	separatorRaw = string(os.PathSeparator)
	// mutex           = &sync.Mutex{}
	countPerProcessRaw = 2
)

type EventRawConversion struct {
	Ctx      *orm.DataContext
	FilePath string
}

func NewEventRawConversion(ctx *orm.DataContext, filePath string) *EventRawConversion {
	ev := new(EventRawConversion)
	ev.Ctx = ctx
	ev.FilePath = filePath
	return ev
}

func (ev *EventRawConversion) Run() {
	var wg sync.WaitGroup

	if ev.FilePath != "" {
		if fileExists(ev.FilePath) {

			brakes := GetAlarmBrake(ev.Ctx)

			files, err := ioutil.ReadDir(ev.FilePath)
			if err != nil {
				log.Println(err)
			}

			counter := 0
			countData := len(files)
			isDone := false
			startIndex := 0
			endIndex := 0

			for !isDone {
				startIndex = counter * countPerProcessRaw
				endIndex = (counter + 1) * countPerProcessRaw

				if endIndex > countData {
					endIndex = countData
				}

				fileprocess := files[startIndex:endIndex]
				for _, f := range fileprocess {
					if strings.Contains(f.Name(), ".xlsx") && !strings.Contains(f.Name(), "~") {
						wg.Add(1)
						go ev.processFile(f.Name(), &wg, brakes)
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

func (ev *EventRawConversion) processFile(filename string, wg *sync.WaitGroup, brakes map[int]AlarmBrake) {
	// mutex.Lock()

	now := time.Now()
	log.Println("Starting process file ", filename)

	total := 0

	fLoc := ev.FilePath + separatorRaw + filename
	fi, err := os.Stat(fLoc)
	turbine := strings.Replace(fi.Name(), ".xlsx", "", 1)
	project := "Tejuva"

	xls, err := xlsx.OpenFile(fLoc)
	if err != nil {
		log.Println("Error open excel file : ", err.Error())
	}

	// var result []orm.IModel

	for _, sheet := range xls.Sheet {
		for idx, row := range sheet.Rows {
			if idx > 0 {
				affectedItem, _ := row.Cells[3].String()
				// eventType, _ := row.Cells[1].String()
				if strings.TrimSpace(affectedItem) != "" { //&& strings.ToUpper(eventType) == "ALARMCHANGED" {
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

					var brakeProgram int
					var brakeType string
					brakeProgram = brakes[alarmIdx].BrakeProgram
					brakeType = brakes[alarmIdx].Type

					sTurbineStatus, _ := row.Cells[7].String()
					turbineStatus := ""
					if sTurbineStatus != "" {
						arrTurbineStatus := strings.Split(sTurbineStatus, " ")
						if len(arrTurbineStatus) > 1 {
							turbineStatus = arrTurbineStatus[1]
						}
					}

					rawdata := new(EventRaw)
					rawdata.ProjectName = project
					rawdata.Turbine = turbine
					sTimeStamp, _ := row.Cells[0].String()
					// rawdata.TimeStamp, _ = time.Parse("2006-01-02 15:04:05", strings.Replace(strings.Replace(sTimeStamp, "T", " ", 1), "+05:30", "", 1))
					//rawdata.TimeStamp, _ = time.Parse("2006-01-02 15:04:05-07:00", strings.Replace(sTimeStamp, "T", " ", 1))
					rawdata.TimeStamp, _ = time.Parse("2006-01-02 15:04:05.000", strings.Replace(strings.Replace(sTimeStamp, "T", " ", 1), "+05:30", "", 1))

					// rawdata.TimeStampUTC = rawdata.TimeStamp.UTC()

					milistr := tk.ToString(rawdata.TimeStamp.Nanosecond() / 1000000)
					timeStampStr := rawdata.TimeStamp.Format("060102150405") + milistr

					rawdata.TimeStampInt, _ = strconv.ParseInt(timeStampStr, 10, 64)

					// log.Printf("%v | %v || %v \n", strings.Replace(strings.Replace(sTimeStamp, "T", " ", 1), "+05:30", " ", 1), sTimeStamp, rawdata.TimeStamp.String())
					// log.Printf("%v | %v || %v \n", rawdata.TimeStamp.String(), rawdata.TimeStampUTC.String(), strings.Replace(sTimeStamp, "T", " ", 1))

					sEventType, _ := row.Cells[1].String()
					rawdata.EventType = strings.TrimSpace(sEventType)

					rawdata.BrakeProgram = brakeProgram
					rawdata.DateInfo = GetDateInfo(rawdata.TimeStamp)
					// rawdata.DateInfoUTC = GetDateInfo(rawdata.TimeStampUTC)
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

					// var tmpResult orm.IModel
					// tmpResult = rawdata

					// result = append(result, tmpResult)
					count := 0
					for {
						rawdata = rawdata.New()
						e := ev.Ctx.Insert(rawdata)
						if e != nil {
							log.Printf("error: %v \n", e.Error())
							if count == 0 {
								total++
							}
						} else {
							break
						}

						if count == 2 {
							break
						}
						count++
					}
				} else {
					total++
				}
			}
		}
	}

	// ctx := orm.New(ev.Conn)

	/*resLen := len(result)
	maxInsert := 2

	if resLen > 0 {
		count := 0
		tmpResult := []orm.IModel{}

	doneLoop:
		for {
			start := count * maxInsert
			end := count*maxInsert + maxInsert

			if len(result[start:end]) == maxInsert {
				tmpResult = result[start:end]
			} else {
				tmpResult = result[start:]
			}

		doneInsert:
			for {
				e := ev.Ctx.InsertBulk(tmpResult)
				if e != nil {
					log.Printf("error: %v \n", e.Error())
					e = ev.Ctx.Connection.Connect()

					end = resLen + 1
					break doneInsert

					if e != nil {
						log.Printf("error: %v \n", e.Error())
					}
				} else {
					break doneInsert
				}
			}

			if end >= resLen {
				break doneLoop
			}
		}
	}*/

	/*for {
		e := ev.Ctx.InsertBulk(result)
		if e != nil {
			log.Printf("error: %v \n", e.Error())
			ev.Ctx.Connection.Close()
			e = ev.Ctx.Connection.Connect()
			if e != nil {
				log.Printf("error: %v \n", e.Error())
			}
		} else {
			break
		}
	}*/

	duration := time.Now().Sub(now)
	log.Println(tk.Sprintf("Process file %v about %v sec(s) | total error: %v", filename, duration.Seconds(), total))

	// mutex.Unlock()

	wg.Done()
}

func GetAlarmBrake(ctx *orm.DataContext) map[int]AlarmBrake {
	alarmbrakes := make([]AlarmBrake, 0)

	csr, err := ctx.Connection.NewQuery().
		From(new(AlarmBrake).TableName()).
		Cursor(nil)

	// defer csr.Close()

	if err != nil {
		log.Printf("ERROR: %v \n", err.Error())
		return nil
	}

	err = csr.Fetch(&alarmbrakes, 0, false)

	if err != nil {
		log.Printf("ERROR: %v \n", err.Error())
		return nil
	}

	// log.Printf("GetAlarmBrake: %v \n", len(alarmbrakes))

	result := map[int]AlarmBrake{}

	for _, val := range alarmbrakes {
		// if tk.IsNilOrEmpty(result[val.AlarmIndex]) {
		result[val.AlarmIndex] = val
		// }
	}

	/*for x, val := range result {
		fmt.Printf("%#v | %#v \n", x, val)
	}*/
	csr.Close()
	return result
}

func fileExists(fileLocation string) bool {
	if _, err := os.Stat(fileLocation); err == nil {
		return true
	}

	return false
}
