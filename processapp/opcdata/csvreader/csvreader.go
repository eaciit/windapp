package csvreader

import (
	"bufio"
	. "eaciit/wfdemo/library/helper"
	. "eaciit/wfdemo/library/models"
	"encoding/csv"
	"github.com/eaciit/orm"
	tk "github.com/eaciit/toolkit"
	"io/ioutil"
	"os"
	"reflect"
	"strings"
	"sync"
	"time"
)

type CsvReader struct {
	FileLocation string
	Ctx          *orm.DataContext
}

var (
	DataTranspose tk.M
	mutex         = &sync.Mutex{}
	mutexData     = &sync.Mutex{}
	idx           = 0
	FileCount     = 0
	DraftDir      = "Draft"
	ProcessDir    = "Process"
	SuccessDir    = "Success"
)

func NewCsvReader(fileLocation string, ctx *orm.DataContext) *CsvReader {
	csv := new(CsvReader)
	csv.FileLocation = fileLocation
	csv.Ctx = ctx

	return csv
}

func fileExists(fileLocation string) bool {
	if _, err := os.Stat(fileLocation); err == nil {
		return true
	}

	return false
}

func (c *CsvReader) Start() {
	if c.FileLocation != "" {
		if fileExists(c.FileLocation) {
			files, err := ioutil.ReadDir(c.FileLocation + "\\" + DraftDir)
			if err != nil {
				tk.Println(err)
			}

			for _, file := range files {
				start := time.Now()

				DataTranspose = tk.M{}
				FileCount++
				c.readFile(file.Name())

				duration := time.Now().Sub(start).Seconds()
				tk.Println(tk.Sprintf("Loading file %v data about %v sec(s)", file.Name(), duration))
			}
		}
	}
}

func (c *CsvReader) readFile(fileName string) {
	var wg sync.WaitGroup

	draftFile := c.FileLocation + "\\" + DraftDir + "\\" + fileName
	processFile := c.FileLocation + "\\" + ProcessDir + "\\" + fileName
	successFile := c.FileLocation + "\\" + SuccessDir + "\\" + fileName

	err := os.Rename(draftFile, processFile)
	if err != nil {
		tk.Println("Error Move Draft File : ", err.Error())
	}

	f, _ := os.Open(processFile)
	r, err := csv.NewReader(bufio.NewReader(f)).ReadAll()
	if err != nil {
		tk.Println("Error Read File : ", err.Error())
	}
	countData := len(r)

	countPerProcess := 5000
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

		data := r[startIndex:endIndex]

		wg.Add(1)
		go func(data [][]string) {
			for _, d := range data {
				mutex.Lock()

				contents := d
				parseContent(contents)
				mutex.Unlock()

				// runtime.Gosched()
			}
			wg.Done()
		}(data)

		counter++

		if endIndex >= countData {
			isFinish = true
		}
	}

	f.Close()

	wg.Wait()

	err = os.Rename(processFile, successFile)
	if err != nil {
		tk.Println("Error Move Process File : ", err.Error())
	}

	if len(DataTranspose) > 0 {
		c.createLog()
		c.insertData()
	}
}

func (c *CsvReader) createLog() string {
	scada := new(ScadaThreeSecs)
	ref := reflect.ValueOf(scada).Elem()
	typeOf := ref.Type()

	content := ""
	delim := ""
	for i := 0; i < ref.NumField(); i++ {
		if typeOf.Field(i).Name != "DateId1Info" && typeOf.Field(i).Name != "DateId2Info" && typeOf.Field(i).Name != "ModelBase" {
			content += delim + typeOf.Field(i).Name
			delim = ","
		}
	}

	f, _ := os.Create(c.FileLocation + "\\Results\\result_" + tk.ToString(FileCount) + ".csv")
	defer f.Close()

	f.WriteString(content + "\n")
	for _, value := range DataTranspose {
		content = ""
		delim = ""
		for i := 0; i < ref.NumField(); i++ {
			if typeOf.Field(i).Name != "DateId1Info" && typeOf.Field(i).Name != "DateId2Info" && typeOf.Field(i).Name != "ModelBase" {
				field := typeOf.Field(i).Name
				fieldType := ref.Field(i).Type()

				if field == "ID" {
					field = "Id"
				}

				var valData interface{}
				if value.(tk.M).Has(field) {
					if fieldType.String() == "time.Time" {
						valTime := value.(tk.M).Get(field).(time.Time)
						valData = valTime.Format("2006-01-02 15:04:05")
					} else {
						valData = value.(tk.M).Get(field)
					}
				} else {
					valData = ""
				}
				content += tk.Sprintf("%v%v", delim, valData)
				delim = ","
			}
		}
		f.WriteString(content + "\n")
	}

	return f.Name()
}

func (c *CsvReader) insertData() {
	tk.Println("Starting to insert data...")
	start := time.Now()

	var wg sync.WaitGroup

	datas := make([]tk.M, 0)
	for _, d := range DataTranspose {
		datas = append(datas, d.(tk.M))
	}
	countData := len(datas)

	// DataTranspose = tk.M{}

	countPerProcess := 5000
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

		data := datas[startIndex:endIndex]

		wg.Add(1)
		go func(data []tk.M) {
			for _, dt := range data {
				mutexData.Lock()

				mdl := new(ScadaThreeSecs).New()
				mdl.ID = dt.GetString("Id")
				mdl.ProjectName = dt.GetString("ProjectName")
				mdl.Turbine = dt.GetString("Turbine")
				mdl.TimeStamp1 = dt.Get("TimeStamp1").(time.Time)
				mdl.TimeStamp2 = dt.Get("TimeStamp2").(time.Time)
				mdl.DateId1 = dt.Get("DateId1").(time.Time)
				mdl.DateId2 = dt.Get("DateId2").(time.Time)
				mdl.THour = mdl.TimeStamp1.Hour()
				mdl.TMinute = mdl.TimeStamp1.Minute()
				mdl.TSecond = mdl.TimeStamp1.Second()
				mdl.TMinuteValue = float64(mdl.TMinute) + tk.Div(float64(mdl.TSecond), 60.0)
				mdl.TMinuteCategory = tk.ToInt(tk.RoundingUp64(tk.Div(mdl.TMinuteValue, 10), 0)*10, "0")
				newTimeStamp := mdl.DateId1.Add(time.Duration(mdl.THour) * time.Hour).Add(time.Duration(mdl.TMinuteCategory) * time.Minute)
				mdl.TimeStampConverted = newTimeStamp.UTC()
				mdl.DateId1Info = GetDateInfo(mdl.DateId1)
				mdl.DateId2Info = GetDateInfo(mdl.DateId2)

				ref := reflect.ValueOf(mdl).Elem()
				typeOf := ref.Type()
				for i := 0; i < ref.NumField(); i++ {
					if typeOf.Field(i).Name != "ID" && typeOf.Field(i).Name != "ModelBase" && typeOf.Field(i).Name != "ProjectName" && typeOf.Field(i).Name != "Turbine" && typeOf.Field(i).Name != "TimeStamp1" && typeOf.Field(i).Name != "TimeStamp2" && typeOf.Field(i).Name != "DateId1" && typeOf.Field(i).Name != "DateId2" && typeOf.Field(i).Name != "THour" && typeOf.Field(i).Name != "TMinute" && typeOf.Field(i).Name != "TSecond" && typeOf.Field(i).Name != "TMinuteValue" && typeOf.Field(i).Name != "TMinuteCategory" && typeOf.Field(i).Name != "TimeStampConverted" && typeOf.Field(i).Name != "DateId1Info" && typeOf.Field(i).Name != "DateId2Info" {
						fieldName := typeOf.Field(i).Name
						setVal := ref.Field(i)
						if setVal.IsValid() {
							if dt.Get(fieldName) != nil {
								setVal.SetFloat(dt.GetFloat64(fieldName))
							} else {
								setVal.SetFloat(-0.0000000001)
							}
						}
					}
				}

				err := c.Ctx.Save(mdl)
				if err != nil {
					tk.Println("Error saving data: " + err.Error())
				}

				mutexData.Unlock()
			}
			wg.Done()
		}(data)

		counter++

		if endIndex >= countData {
			isFinish = true
		}
	}

	wg.Wait()

	duration := time.Now().Sub(start).Seconds()
	tk.Println(tk.Sprintf("End insert data as long as %v secs", duration))
}

func parseContent(contents []string) {
	time1, _ := time.Parse("02-Jan-2006 15:04:05", contents[0])
	time2, _ := time.Parse("02-Jan-2006 15:04:05", contents[1])
	date1, _ := time.Parse("2006-01-02", time1.Format("2006-01-02"))
	date2, _ := time.Parse("2006-01-02", time2.Format("2006-01-02"))

	thour := time1.Hour()
	tminute := time1.Minute()
	tsecond := time1.Second()
	tminutevalue := float64(tminute) + tk.Div(float64(tsecond), 60.0)
	tminutecategory := tk.ToInt(tk.RoundingUp64(tk.Div(tminutevalue, 10), 0)*10, "0")
	if tminutecategory == 60 {
		tminutecategory = 0
		thour = thour + 1
	}
	newTimeStamp := date1.Add(time.Duration(thour) * time.Hour).Add(time.Duration(tminutecategory) * time.Minute)
	timestampconverted := newTimeStamp.UTC()

	infos := strings.Split(contents[2], ".")
	value := tk.ToFloat64(contents[3], 6, tk.RoundingAuto)

	project := "Tejuva"
	turbine := infos[2]
	column := infos[3] + "_" + infos[4]

	id := time1.Format("20060102_150405") + "_" + time2.Format("20060102_150405") + "_" + project + "_" + turbine

	if DataTranspose.Get(id) == nil {
		DataTranspose.Set(id, tk.M{}.Set("Id", id).Set("ProjectName", project).Set("Turbine", turbine).Set("TimeStamp1", time1).Set("TimeStamp2", time2).Set("DateId1", date1).Set("DateId2", date2).Set("THour", thour).Set("TMinute", tminute).Set("TSecond", tsecond).Set("TMinuteValue", tminutevalue).Set("TMinuteCategory", tminutecategory).Set("TimeStampConverted", timestampconverted).Set(column, value))
	} else {
		newData := DataTranspose.Get(id).(tk.M)
		DataTranspose.Set(id, newData.Set(column, value))
	}
}
