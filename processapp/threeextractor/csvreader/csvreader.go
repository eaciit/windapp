package csvreader

import (
	"bufio"
	. "eaciit/wfdemo/library/helper"
	. "eaciit/wfdemo/library/models"
	dc "eaciit/wfdemo/processapp/threeextractor/dataconversion"
	"encoding/csv"
	"io/ioutil"
	"log"
	"os"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/eaciit/dbox"
	"github.com/eaciit/orm"
	tk "github.com/eaciit/toolkit"
)

type CsvReader struct {
	FileLocation string
	Ctx          *orm.DataContext
}

var (
	DataTranspose   tk.M
	mutex           = &sync.Mutex{}
	mutexData       = &sync.Mutex{}
	idx             = 0
	FileCount       = 0
	DraftDir        = "Draft"
	ProcessDir      = "Process"
	SuccessDir      = "Success"
	emptyValueSmall = -0.000001
	emptyValueBig   = -9999999.0
	separator       = string(os.PathSeparator)
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
	// tk.Printf("c.FileLocation %v \n", c.FileLocation)
	if c.FileLocation != "" {
		if fileExists(c.FileLocation) {
			files, err := ioutil.ReadDir(c.FileLocation + separator + DraftDir)
			if err != nil {
				tk.Println(err)
			}

			for _, file := range files {
				// locFile := c.FileLocation + separator + DraftDir + separator + file.Name()
				if strings.Contains(file.Name(), ".csv") {
					tk.Printf("Load file: %v \n", file.Name())
					start := time.Now()

					DataTranspose = tk.M{}
					FileCount++
					c.readFile(file.Name())

					duration := time.Now().Sub(start).Seconds()
					tk.Println(tk.Sprintf("Loading file %v data about %v sec(s)\n", file.Name(), duration))
				}
			}
		}
	}
}

func (c *CsvReader) readFile(fileName string) {
	var wg sync.WaitGroup

	draftFile := c.FileLocation + separator + DraftDir + separator + fileName
	processFile := c.FileLocation + separator + ProcessDir + separator + fileName
	successFile := c.FileLocation + separator + SuccessDir + separator + fileName

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

	countPerProcess := 1000
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
				parseContent(contents, fileName)
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

	if len(DataTranspose) > 0 {
		c.insertData()
		c.convert3Ext(fileName)
		c.convert10Min(fileName)
		// c.remove3Seconds(fileName)
	}

	err = os.Rename(processFile, successFile)
	if err != nil {
		tk.Println("Error Move Process File : ", err.Error())
	}
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

				timeStamp := val.TimeStamp1.UTC()
				seconds := tk.Div(tk.ToFloat64(timeStamp.Nanosecond(), 1, tk.RoundingAuto), 1000000000)
				secondsInt := tk.ToInt(seconds, tk.RoundingAuto)
				newTimeTmp := timeStamp.Add(time.Duration(secondsInt) * time.Second)
				strTime := tk.ToString(newTimeTmp.Year()) + tk.ToString(int(newTimeTmp.Month())) + tk.ToString(newTimeTmp.Day()) + " " + tk.ToString(newTimeTmp.Hour()) + ":" + tk.ToString(newTimeTmp.Minute()) + ":" + tk.ToString(newTimeTmp.Second())
				newTime, _ := time.Parse("200612 15:4:5", strTime)

				mdl.TimeStampSecondGroup = newTime

				mdl.THour = mdl.TimeStampSecondGroup.Hour()
				mdl.TMinute = mdl.TimeStampSecondGroup.Minute()
				mdl.TSecond = mdl.TimeStampSecondGroup.Second()
				mdl.TMinuteValue = float64(mdl.TMinute) + tk.Div(float64(mdl.TSecond), 60.0)
				mdl.TMinuteCategory = tk.ToInt(tk.RoundingUp64(tk.Div(mdl.TMinuteValue, 10), 0)*10, "0")
				newTimeStamp := mdl.DateId1.Add(time.Duration(mdl.THour) * time.Hour).Add(time.Duration(mdl.TMinuteCategory) * time.Minute)
				mdl.TimeStampConverted = newTimeStamp.UTC()
				mdl.TimeStampConvertedInt, _ = strconv.ParseInt(mdl.TimeStampConverted.Format("200601021504"), 10, 64)

				mdl.DateId1Info = GetDateInfo(mdl.DateId1)
				mdl.DateId2Info = GetDateInfo(mdl.DateId2)
				mdl.File = dt.GetString("file")

				ref := reflect.ValueOf(mdl).Elem()
				typeOf := ref.Type()
				for i := 0; i < ref.NumField(); i++ {
					if typeOf.Field(i).Name != "ID" && typeOf.Field(i).Name != "ModelBase" && typeOf.Field(i).Name != "ProjectName" && typeOf.Field(i).Name != "Turbine" && typeOf.Field(i).Name != "TimeStamp1" && typeOf.Field(i).Name != "TimeStamp2" && typeOf.Field(i).Name != "DateId1" && typeOf.Field(i).Name != "DateId2" && typeOf.Field(i).Name != "THour" && typeOf.Field(i).Name != "TMinute" && typeOf.Field(i).Name != "TSecond" && typeOf.Field(i).Name != "TMinuteValue" && typeOf.Field(i).Name != "TMinuteCategory" && typeOf.Field(i).Name != "TimeStampConverted" && typeOf.Field(i).Name != "DateId1Info" && typeOf.Field(i).Name != "DateId2Info" && typeOf.Field(i).Name != "Line" && typeOf.Field(i).Name != "File" && typeOf.Field(i).Name != "TimeStampConvertedInt" {
						fieldName := typeOf.Field(i).Name
						setVal := ref.Field(i)
						if setVal.IsValid() {
							if dt.Has(fieldName) {
								setVal.SetFloat(dt.GetFloat64(fieldName))
							} else {
								// setVal.SetFloat(emptyValueSmall)
								setVal.SetFloat(emptyValueBig)
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

func (c *CsvReader) convert3Ext(fileName string) {
	tk.Println("Start Convert data 3secs to 3ext")
	conv := dc.NewConvThreeExt(c.Ctx)
	conv.Generate(fileName)
	tk.Println("End Convert data 3secs to 3ext")

	//c.Ctx.Connection.NewQuery().From(new(ScadaThreeSecs).TableName()).Delete()
}

func (c *CsvReader) convert10Min(fileName string) {
	tk.Println("Start Convert data 3secs to 10mins")
	conv := dc.NewDataConversion(c.Ctx)
	conv.Generate(fileName)
	tk.Println("End Convert data 3secs to 10mins")

	//c.Ctx.Connection.NewQuery().From(new(ScadaThreeSecs).TableName()).Delete()
}

func (c *CsvReader) remove3Seconds(fileName string) {
	tk.Println("Start Remove 3secs data")
	e := c.Ctx.DeleteMany(new(ScadaThreeSecs), dbox.Eq("file", fileName))

	if e != nil {
		log.Printf("Error Remove 3secs data: %v \n", e.Error())
	} else {
		tk.Println("End Remove 3secs data")
	}

	//c.Ctx.Connection.NewQuery().From(new(ScadaThreeSecs).TableName()).Delete()
}

func parseContent(contents []string, fileName string) {
	time1, _ := time.Parse("02-Jan-2006 15:04:05", contents[0])
	time2, _ := time.Parse("02-Jan-2006 15:04:05", contents[1])
	date1, _ := time.Parse("2006-01-02", time1.Format("2006-01-02"))
	date2, _ := time.Parse("2006-01-02", time2.Format("2006-01-02"))
	infos := strings.Split(contents[2], ".")
	value := tk.ToFloat64(contents[3], 6, tk.RoundingAuto)

	project := "Tejuva"
	turbine := infos[2]
	column := infos[3] + "_" + infos[4]

	id := time1.Format("20060102_150405") + "_" + time2.Format("20060102_150405") + "_" + project + "_" + turbine

	if DataTranspose.Get(id) == nil {
		DataTranspose.Set(id, tk.M{}.Set("Id", id).Set("ProjectName", project).Set("Turbine", turbine).Set("TimeStamp1", time1).Set("TimeStamp2", time2).Set("DateId1", date1).Set("DateId2", date2).Set(column, value).Set("file", fileName))
	} else {
		newData := DataTranspose.Get(id).(tk.M)
		DataTranspose.Set(id, newData.Set(column, value))
	}
}
}
