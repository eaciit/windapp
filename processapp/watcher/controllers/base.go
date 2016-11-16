package controllers

import (
	"bufio"
	"fmt"
	"log"

	// . "eaciit/wfdemo/library/helper"
	. "eaciit/wfdemo/library/models"

	"github.com/eaciit/dbox"
	"github.com/eaciit/orm"
	tk "github.com/eaciit/toolkit"
	// "math"
	"io/ioutil"
	"os"
	"reflect"
	"strings"
	"sync"
	"time"

	_ "github.com/eaciit/dbox/dbc/mongo"
)

var (
	wd = func() string {
		d, _ := os.Getwd()
		return d + separator
	}()

	// mu                 = &sync.Mutex{}
	retry              = 10
	worker             = 100
	maxDataEachProcess = 100000
	idx                = 0
	mu                 = &sync.Mutex{}
	muinsert           = &sync.Mutex{}
	emptyValueSmall    = -0.000001
	emptyValueBig      = -9999999.0
	separator          = string(os.PathSeparator)
)

type IBaseController interface {
	// not implemented anything yet
}

type BaseController struct {
	base IBaseController
	Ctx  *orm.DataContext
}

type TenMinuteInfo struct {
	THour              int
	TMinute            int
	TSecond            int
	TMinuteValue       float64
	TMinuteCategory    int
	TimeStampConverted time.Time
}

func (b *BaseController) InsertBulk(result []tk.M, m orm.IModel, wg *sync.WaitGroup) {
	var datas []orm.IModel
	for _, i := range result {
		valueType := reflect.TypeOf(m).Elem()
		for f := 0; f < valueType.NumField(); f++ {
			field := valueType.Field(f)
			bsonField := field.Tag.Get("bson")
			jsonField := field.Tag.Get("json")

			if jsonField != bsonField && field.Name != "RWMutex" && field.Name != "ModelBase" {
				i.Set(field.Name, GetMgoValue(i, bsonField))
			}
			switch field.Type.Name() {
			case "string":
				if GetMgoValue(i, bsonField) == nil {
					i.Set(field.Name, "")
				}
				break
			case "Time":
				if GetMgoValue(i, bsonField) == nil {
					i.Set(field.Name, time.Time{})
				} else {
					i.Set(field.Name, GetMgoValue(i, bsonField).(time.Time).UTC())
				}
				break
			default:
				break
			}

		}

		newPointer := getNewPointer(m)
		e := tk.Serde(i, newPointer, "json")
		datas = append(datas, newPointer)

		if e != nil {
			tk.Printf("\n----------- ERROR -------------- \n %v \n\n %#v \n\n %#v \n-------------------------  \n", e.Error(), i, newPointer)
			wg.Done()
		}

	}

	if nil != datas {
		muinsert.Lock()
		for {
			e := b.Ctx.InsertBulk(datas)
			if e == nil {
				ctn := len(result)
				idx += ctn
				tk.Printf("saved: %v data(s)\n", idx)
				break
			} else {
				b.Ctx.Connection.Connect()
			}
		}
		muinsert.Unlock()
	}

	wg.Done()
}

func (b *BaseController) Insert(result []tk.M, m orm.IModel, wg *sync.WaitGroup) {
	// muinsert := &sync.Mutex{}
	for _, i := range result {
		valueType := reflect.TypeOf(m).Elem()
		for f := 0; f < valueType.NumField(); f++ {
			field := valueType.Field(f)
			bsonField := field.Tag.Get("bson")
			jsonField := field.Tag.Get("json")

			if jsonField != bsonField && field.Name != "RWMutex" && field.Name != "ModelBase" {
				i.Set(field.Name, GetMgoValue(i, bsonField))
			}
			switch field.Type.Name() {
			case "string":
				if GetMgoValue(i, bsonField) == nil {
					i.Set(field.Name, "")
				}
				break
			case "Time":
				if GetMgoValue(i, bsonField) == nil {
					i.Set(field.Name, time.Time{})
				} else {
					i.Set(field.Name, GetMgoValue(i, bsonField).(time.Time).UTC())
				}
				break
			default:
				break
			}

		}

		newPointer := getNewPointer(m)
		e := tk.Serde(i, newPointer, "json")
		var newId int64
		for index := 0; index < retry; index++ {
			muinsert.Lock()
			newId, e = b.Ctx.InsertOut(newPointer)
			_ = newId
			muinsert.Unlock()
			if e == nil {
				wg.Done()
				break
			} else {
				b.Ctx.Connection.Connect()
			}
		}

		if e != nil {
			tk.Printf("\n----------- ERROR -------------- \n %v \n\n %#v \n\n %#v \n-------------------------  \n", e.Error(), i, newPointer)
			wg.Done()
		}

	}
	wg.Done()
}
func GetMgoValue(d tk.M, fieldName string) interface{} {
	index := strings.Index(fieldName, ".")
	if index < 0 {
		return d.Get(fieldName)
	} else {
		data := d.Get(fieldName[0:index])
		if data != nil {
			return GetMgoValue(data.(tk.M), fieldName[(index+1):len(fieldName)])
		} else {
			return nil
		}
	}
}

func (b *BaseController) GetById(m orm.IModel, id interface{}, column_name ...string) error {
	var e error
	c := b.Ctx.Connection
	column_id := "Id"
	if column_name != nil && len(column_name) > 0 {
		column_id = column_name[0]
	}
	csr, e := c.NewQuery().From(m.(orm.IModel).TableName()).Where(dbox.Eq(column_id, id)).Cursor(nil)
	defer csr.Close()
	if e != nil {
		return e
	}
	e = csr.Fetch(m, 1, false)
	if e != nil {
		return e
	}

	return nil
}

func getNewPointer(m orm.IModel) orm.IModel {
	switch m.TableName() {
	case "ScadaData":
		return new(ScadaData)
	case "Alarm":
		return new(Alarm)
	default:
		return m
	}

}

func (b *BaseController) Delete(m orm.IModel, id interface{}, column_name ...string) error {
	column_id := "Id"
	if column_name != nil && len(column_name) > 0 {
		column_id = column_name[0]
	}
	e := b.Ctx.Connection.NewQuery().From(m.(orm.IModel).TableName()).Where(dbox.Eq(column_id, id)).Delete().Exec(nil)
	if e != nil {
		return e
	}
	return nil
}

func (b *BaseController) Update(m orm.IModel, id interface{}, column_name ...string) error {
	column_id := "Id"
	if column_name != nil && len(column_name) > 0 {
		column_id = column_name[0]
	}
	e := b.Ctx.Connection.NewQuery().From(m.(orm.IModel).TableName()).Where(dbox.Eq(column_id, id)).Update().Exec(tk.M{"data": m})
	if e != nil {
		return e
	}
	return nil
}

func (b *BaseController) Truncate(m orm.IModel) error {
	c := b.Ctx.Connection
	e := c.NewQuery().From(m.(orm.IModel).TableName()).Delete().Exec(nil)
	if e != nil {
		return e
	}

	return nil
}
func (b *BaseController) CloseDb() {
	if b.Ctx != nil {
		b.Ctx.Close()
	}
}

func (b *BaseController) WriteLog(msg interface{}) {
	log.Printf("%#v\n\r", msg)
	return
}

func PrepareConnection() (dbox.IConnection, error) {
	config := Config()
	ci := &dbox.ConnectionInfo{config["host"], config["database"], config["username"], config["password"], nil}
	c, e := dbox.NewConnection("mongo", ci)

	if e != nil {
		return nil, e
	}

	e = c.Connect()

	tk.Printf("%#v \n", e)

	if e != nil {
		return nil, e
	}

	return c, nil
}

func (b *BaseController) GetDataSource(dataSourceFolder string) ([]os.FileInfo, string) {
	config := Config()
	source := config["datasource"]
	files, e := ioutil.ReadDir(source + "\\" + dataSourceFolder)
	if e != nil {
		tk.Println(e)
		os.Exit(0)
	}
	return files, source + "\\" + dataSourceFolder
}

func (b *BaseController) GetDataSourceDirect(dataSourceFolder string) ([]os.FileInfo, string) {
	files, e := ioutil.ReadDir(dataSourceFolder)
	if e != nil {
		tk.Println(e)
		os.Exit(0)
	}
	return files, dataSourceFolder
}

func Config() map[string]string {
	ret := make(map[string]string)
	file, err := os.Open(wd + ".." + separator + "conf" + separator + "app.conf")
	if err == nil {
		defer file.Close()

		reader := bufio.NewReader(file)
		for {
			line, _, e := reader.ReadLine()
			if e != nil {
				break
			}

			sval := strings.Split(string(line), "=")
			ret[sval[0]] = sval[1]
		}
	} else {
		tk.Println(err.Error())
	}

	return ret
}

func WriteWatcherErrors(errorLine tk.M, fileName string, folder string) (e error) {
	fileName = fileName + "_" + tk.GenerateRandomString("", 5) + ".txt"
	errors := ""

	for x, err := range errorLine {
		errors = errors + "" + fmt.Sprintf("#%v: %#v \n", x, err)
	}

	if errors != "" {
		e = ioutil.WriteFile(folder+separator+fileName, []byte(errors), 0644)
		log.Printf("Saving Errors... %v\n", fileName)
	}

	return
}

func GenTenMinuteInfo(timestamp time.Time) (result TenMinuteInfo) {
	result.THour = timestamp.Hour()
	result.TMinute = timestamp.Minute()
	result.TSecond = timestamp.Second()
	result.TMinuteValue = tk.ToFloat64(result.TMinute, 2, tk.RoundingAuto) + (tk.ToFloat64(result.TSecond, 2, tk.RoundingAuto) / 60)

	switch {
	case result.TMinuteValue <= 10:
		result.TMinuteCategory = 10
	case result.TMinuteValue <= 20:
		result.TMinuteCategory = 20
	case result.TMinuteValue <= 30:
		result.TMinuteCategory = 30
	case result.TMinuteValue <= 40:
		result.TMinuteCategory = 40
	case result.TMinuteValue <= 50:
		result.TMinuteCategory = 50
	case result.TMinuteValue <= 60:
		result.TMinuteCategory = 0

		// log.Printf("timestamp_before: %v | ", timestamp)

		if result.THour+1 >= 24 {
			timestamp = timestamp.AddDate(0, 0, 1)
		}

		result.TimeStampConverted, _ = time.Parse("2006-1-2 15:4:05",
			tk.ToString(timestamp.Year())+"-"+tk.ToString(int(timestamp.Month()))+"-"+tk.ToString(timestamp.Day())+" 00:0:00")

		// log.Printf(" timestamp_after: %v \n", result.TimeStampConverted.String())

		return
	}

	result.TimeStampConverted, _ = time.Parse("2006-1-2 15:4:05",
		tk.ToString(timestamp.Year())+"-"+tk.ToString(int(timestamp.Month()))+"-"+tk.ToString(timestamp.Day())+" "+tk.ToString(timestamp.Hour())+":"+tk.ToString(result.TMinuteCategory)+":00")

	return result
}

func SetFloatValue(strVal string) (result float64) {
	if strVal == "" {
		// result = emptyValue
		result = emptyValueBig
	} else {
		result, _ = tk.StringToFloat(strVal)
	}

	return
}
