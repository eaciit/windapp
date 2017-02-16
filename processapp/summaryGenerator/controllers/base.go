package controllers

import (
	"log"

	. "eaciit/wfdemo-git/library/helper"
	. "eaciit/wfdemo-git/library/models"

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

type LatestData struct {
	Alarm             []TurbineLatest
	EventDown         []TurbineLatest
	ScadaData         []TurbineLatest
	ScadaDataOEM      []TurbineLatest
	ScadaSummaryDaily []TurbineLatest

	MapAlarm             map[string]time.Time
	MapEventDown         map[string]time.Time
	MapScadaData         map[string]time.Time
	MapScadaDataOEM      map[string]time.Time
	MapScadaSummaryDaily map[string]time.Time
}

type TurbineLatest struct {
	project    string
	turbine    string
	latestTime time.Time
}

type BaseController struct {
	base        IBaseController
	Ctx         *orm.DataContext
	LatestData  LatestData
	RefTurbines tk.M
	RefAlarms   tk.M
}

func (b *BaseController) PrepareDataReff() {
	tk.Println("Getting data refference")
	logStart := time.Now()

	turbines := []TurbineMaster{}
	csrt, e := b.Ctx.Connection.NewQuery().From(new(TurbineMaster).TableName()).Cursor(nil)

	tk.Println("Get Turbines")

	e = csrt.Fetch(&turbines, 0, false)
	ErrorHandler(e, "get turbine master")
	csrt.Close()

	b.RefTurbines = tk.M{}
	for _, t := range turbines {
		b.RefTurbines.Set(t.TurbineId, tk.M{}.
			Set("turbinename", t.TurbineName).
			Set("turbineelevation", t.Elevation))
	}

	tk.Printf("Turbines: %v \n", len(turbines))

	tk.Println("Get EventDown")

	b.RefAlarms = tk.M{}

	for turbine, _ := range b.RefTurbines {
		filter := []*dbox.Filter{}
		filter = append(filter, dbox.Eq("projectname", "Tejuva"))
		filter = append(filter, dbox.Eq("turbine", turbine))
		filter = append(filter, dbox.Gt("timeend", b.LatestData.MapAlarm["Tejuva#"+turbine]))

		alarms := []EventDown{}
		csr2, e := b.Ctx.Connection.NewQuery().From(new(EventDown).TableName()).
			Where(filter...).Cursor(nil)

		e = csr2.Fetch(&alarms, 0, false)
		ErrorHandler(e, "get alarm data")
		csr2.Close()

		// tk.Printf("EventDown for: %v | %v \n", turbine, len(alarms))

		details := []EventDown{}
		for _, a := range alarms {
			if b.RefAlarms.Has(a.Turbine) {
				details = b.RefAlarms.Get(a.Turbine).([]EventDown)
			} else {
				details = []EventDown{}
			}

			details = append(details, a)
			b.RefAlarms.Set(a.Turbine, details)
		}
	}

	logDuration := time.Now().Sub(logStart).Seconds()
	tk.Printf("\nGetting refference data about %v secs\n", logDuration)
}

func (b *BaseController) SetCollectionLatestTime() {
	b.LatestData.Alarm, b.LatestData.MapAlarm = getLatestTime("farm", "turbine", "enddate", new(Alarm).TableName(), b.Ctx)
	b.LatestData.EventDown, b.LatestData.MapEventDown = getLatestTime("projectname", "turbine", "timeend", new(EventDown).TableName(), b.Ctx)
	b.LatestData.ScadaData, b.LatestData.MapScadaData = getLatestTime("projectname", "turbine", "timestamp", new(ScadaData).TableName(), b.Ctx)
	b.LatestData.ScadaDataOEM, b.LatestData.MapScadaDataOEM = getLatestTime("projectname", "turbine", "timestamp", new(ScadaDataOEM).TableName(), b.Ctx)
	b.LatestData.ScadaSummaryDaily, b.LatestData.MapScadaSummaryDaily = getLatestTime("projectname", "turbine", "dateinfo.dateid", new(ScadaSummaryDaily).TableName(), b.Ctx)
}

func getLatestTime(projectCol string, turbineCol string, timestampCol string, collection string, ctx *orm.DataContext) (res []TurbineLatest, resMap map[string]time.Time) {
	var (
		pipes []tk.M
	)

	group := tk.M{
		"_id": tk.M{
			"project": "$" + projectCol,
			"turbine": "$" + turbineCol,
		},
		"timestamp": tk.M{"$max": "$" + timestampCol},
	}

	pipes = append(pipes, tk.M{"$group": group})

	csr, e := ctx.Connection.NewQuery().
		From(collection).
		Command("pipe", pipes).
		Cursor(nil)

	defer csr.Close()

	if e != nil {
		log.Printf("Error getLatestTime", e.Error())
	}

	result := []tk.M{}
	e = csr.Fetch(&result, 0, false)
	if e != nil {
		log.Printf("Error getLatestTime", e.Error())
	}

	resMap = map[string]time.Time{}

	for _, r := range result {
		id := r.Get("_id").(tk.M)
		project := id.GetString("project")
		turbine := id.GetString("turbine")
		timestamp := r.Get("timestamp").(time.Time)

		res = append(res, TurbineLatest{
			project:    project,
			turbine:    turbine,
			latestTime: timestamp,
		})

		resMap[project+"#"+turbine] = timestamp
	}

	return
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
