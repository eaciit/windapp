package controllers

import (
	"log"

	. "eaciit/wfdemo-git-dev/library/helper"
	. "eaciit/wfdemo-git-dev/library/models"

	_ "github.com/eaciit/dbox/dbc/mongo"
	"github.com/eaciit/orm"

	"github.com/eaciit/dbox"

	tk "github.com/eaciit/toolkit"
	// "math"
	"io/ioutil"
	"os"
	"reflect"
	"strings"
	"sync"
	"time"
)

var (
	// mu                 = &sync.Mutex{}
	retry              = 10
	worker             = 100
	maxDataEachProcess = 100000
	idx                = 0
	mu                 = &sync.Mutex{}
	muinsert           = &sync.Mutex{}
)

type IBaseController interface {
	// not implemented anything yet
}

type BaseController struct {
	base IBaseController
	Ctx  *orm.DataContext
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
	config := ReadConfig()
	ci := &dbox.ConnectionInfo{config["host"], config["database"], config["username"], config["password"], nil}
	c, e := dbox.NewConnection("mongo", ci)

	if e != nil {
		return nil, e
	}

	e = c.Connect()
	if e != nil {
		return nil, e
	}
	tk.Println("DB Connect...\n")
	return c, nil
}

func (b *BaseController) GetDataSource(dataSourceFolder string) ([]os.FileInfo, string) {
	config := ReadConfig()
	source := config["datasource"]
	files, e := ioutil.ReadDir(source + string(os.PathSeparator) + dataSourceFolder)
	if e != nil {
		tk.Println(e)
		os.Exit(0)
	}
	return files, source + string(os.PathSeparator) + dataSourceFolder
}

func (b *BaseController) GetDataSourceDirect(dataSourceFolder string) ([]os.FileInfo, string) {
	files, e := ioutil.ReadDir(dataSourceFolder)
	if e != nil {
		tk.Println(e)
		os.Exit(0)
	}
	return files, dataSourceFolder
}
