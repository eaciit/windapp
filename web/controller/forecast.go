package controller

import (
	"bytes"
	. "eaciit/wfdemo-git/library/core"
	. "eaciit/wfdemo-git/library/helper"
	. "eaciit/wfdemo-git/library/models"
	"eaciit/wfdemo-git/web/helper"
	"io"
	"io/ioutil"
	"math"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"sync"
	"time"

	"github.com/eaciit/dbox"
	"github.com/eaciit/knot/knot.v1"
	tk "github.com/eaciit/toolkit"
	xls "github.com/tealeg/xlsx"
	gomail "gopkg.in/gomail.v2"
)

var (
	dataMtx sync.Mutex
)

type TimePeriodModel struct {
	TimePeriod time.Time
	DatePeriod time.Time
	TimeBlock  int
	TimeRange  string
}

type ForecastController struct {
	App
}

func CreateForecastController() *ForecastController {
	var controller = new(ForecastController)
	return controller
}

func getTimeBlock(currtime time.Time) int {
	currhour := currtime.Hour()
	currmint := currtime.Minute()
	mindiff := 4 - (60-currmint)/15 + 1
	timeblock := (currhour * 4) + 96 - (96 - mindiff)

	return timeblock
}

func get15MinPeriod(tstart time.Time, tend time.Time) []TimePeriodModel {
	timePeriods := []TimePeriodModel{}
	// tend = tend.Add(time.Duration(3) * time.Minute)
	if tend.Sub(tstart).Minutes() >= 0 {
		currTime := tstart
		dateid, _ := time.Parse("2006-01-02", currTime.Format("2006-01-02"))
		befTime := currTime.Add(time.Duration(-15) * time.Minute)
		timeBlock := getTimeBlock(befTime)
		item := TimePeriodModel{
			TimePeriod: currTime,
			DatePeriod: dateid,
			TimeBlock:  timeBlock,
			TimeRange:  tk.Sprintf("%s - %s", befTime.Format("15:04"), currTime.Format("15:04")),
		}
		timePeriods = append(timePeriods, item)

		for {
			currTime = currTime.Add(time.Duration(15) * time.Minute)
			dateid, _ = time.Parse("2006-01-02", currTime.Format("2006-01-02"))
			befTime = currTime.Add(time.Duration(-15) * time.Minute)
			timeBlock = getTimeBlock(befTime)
			if currTime.Sub(tend).Minutes() > 0 {
				break
			}
			item := TimePeriodModel{
				TimePeriod: currTime,
				DatePeriod: dateid,
				TimeBlock:  timeBlock,
				TimeRange:  tk.Sprintf("%s - %s", befTime.Format("15:04"), currTime.Format("15:04")),
			}
			timePeriods = append(timePeriods, item)
		}
	}

	return timePeriods
}

func getDirectoriesToRead(tstart time.Time, tend time.Time) (dirs []string) {
	dirs = []string{}

	if tend.Sub(tstart).Hours() >= 0 {
		tdir := tstart
		dir := tdir.Format("20060102")
		dirs = append(dirs, dir)
		for {
			tdir = tdir.AddDate(0, 0, 1)
			if tdir.Sub(tend).Hours() > 0 {
				break
			}
			dir = tdir.Format("20060102")
			dirs = append(dirs, dir)
		}
	}

	return
}

func (m *ForecastController) GetListTurbineDown(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	type Payload struct {
		Project string
	}
	p := new(Payload)
	e := k.GetPayload(&p)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	dataReturn := []tk.M{}
	today, _ := time.Parse("2006-01-02", time.Now().Format("2006-01-02"))
	dateend, _ := time.Parse("2006-01-02 15:04:05", "0001-01-01 00:00:00")
	matches := []tk.M{
		tk.M{"projectname": p.Project},
		//tk.M{"$or": []tk.M{
		tk.M{"dateinfostart.dateid": tk.M{"$eq": today}},
		tk.M{"dateinfoend.dateid": tk.M{"$eq": dateend}},
		tk.M{"isdeleted": false},
		//}},
	}
	pipes := []tk.M{
		tk.M{"$match": tk.M{"$and": matches}},
	}
	csrtd, e := DBRealtime().NewQuery().
		From("AlarmHFD").
		Command("pipe", pipes).
		Cursor(nil)
	defer csrtd.Close()

	e = csrtd.Fetch(&dataReturn, 0, false)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	return helper.CreateResult(true, dataReturn, "")
}

func (m *ForecastController) UpdateSldc(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	type Item struct {
		Id       string
		Value    float64
		ValueCap float64
	}
	type Payload struct {
		Project string
		Values  []Item
	}
	p := new(Payload)
	e := k.GetPayload(&p)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	if len(p.Values) > 0 {
		for _, v := range p.Values {
			err := DB().Connection.NewQuery().Update().From("ForecastData").Where(dbox.Eq("_id", v.Id)).Exec(tk.M{}.Set("data", tk.M{}.Set("schsdlc", v.Value).Set("avgcapacity", v.ValueCap).Set("isedited", 1)))
			if err != nil {
				tk.Printf("Update data failed : %s\n", err.Error())
			}
		}
	}

	dataReturn := tk.M{
		"success": true,
		"message": "",
	}

	return helper.CreateResult(true, dataReturn, "")
}

func (m *ForecastController) GetList(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	var wg sync.WaitGroup

	p := new(PayloadAnalytic)
	e := k.GetPayload(&p)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	tStart, tEnd, e := helper.GetStartEndDate(k, p.Period, p.DateStart, p.DateEnd)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}
	tStart = tStart.Add(time.Duration(15) * time.Minute)
	tEnd = tEnd.Add(time.Duration(15) * time.Minute)
	timeperiods := get15MinPeriod(tStart, tEnd)
	project := p.Project
	timeNowIst := time.Now().UTC().Add(time.Duration(330) * time.Minute)

	getscada15minpath := GetConfig("scada15min_path", "")
	scada15minpath := ""
	if getscada15minpath == "" || getscada15minpath == nil {
		scada15minpath = "/mnt/data/ostrorealtime/scada15minrev/data"
	} else {
		scada15minpath = tk.ToString(getscada15minpath)
	}

	// get pc reff
	csrpc, e := DB().Connection.NewQuery().From("ref_powercurve").
		Where(dbox.Eq("model", project)).
		Order("windspeed").Cursor(nil)
	defer csrpc.Close()
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	pcSrc := []tk.M{}
	e = csrpc.Fetch(&pcSrc, 0, false)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	latestSubject := []string{}
	// get latest subject
	subToday, _ := time.Parse("2006-01-02", time.Now().Format("2006-01-02"))
	subNextday := subToday.AddDate(0, 0, 1)
	matchessub := []tk.M{
		tk.M{"projectname": project},
		tk.M{"dateinfo.dateid": tk.M{"$gte": subToday}},
		tk.M{"dateinfo.dateid": tk.M{"$lte": subNextday}},
	}
	pipessub := []tk.M{
		tk.M{"$match": tk.M{"$and": matchessub}},
		tk.M{"$group": tk.M{
			"_id": tk.M{
				"projectname": "$projectname",
				"dateid":      "$dateinfo.dateid",
			},
			"max_subject": tk.M{"$max": "$mailsubject"},
		}},
		tk.M{"$sort": tk.M{"max_subject": 1}},
	}
	csrsub, e := DB().Connection.NewQuery().
		From(new(ForecastData).TableName()).
		Command("pipe", pipessub).
		Cursor(nil)
	defer csrsub.Close()

	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}
	forecastsubject := []tk.M{}
	e = csrsub.Fetch(&forecastsubject, 0, false)
	if len(forecastsubject) > 0 {
		prevSubject := ""
		for _, d := range forecastsubject {
			if prevSubject != d.GetString("max_subject") || prevSubject == "" {
				latestSubject = append(latestSubject, d.GetString("max_subject"))
			}
			prevSubject = d.GetString("max_subject")
		}
	}

	// get data forecast
	matches := []tk.M{
		tk.M{"projectname": project},
		tk.M{"timestamp": tk.M{"$gte": tStart}},
		tk.M{"timestamp": tk.M{"$lte": tEnd}},
	}
	pipes := []tk.M{
		tk.M{"$match": tk.M{"$and": matches}},
	}
	csr, e := DB().Connection.NewQuery().
		From(new(ForecastData).TableName()).
		Command("pipe", pipes).
		Cursor(nil)
	defer csr.Close()

	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}
	forecast := tk.M{}
	for {
		item := tk.M{}
		e = csr.Fetch(&item, 1, false)
		if e != nil {
			break
		}
		timestamp := item.Get("timestamp", time.Time{}).(time.Time).UTC() //.UTC().Add(time.Duration(330) * time.Minute)
		if !timestamp.IsZero() {
			forecast.Set(timestamp.Format("20060102_150405"), item)
		}
	}
	// tk.Printf("%#v\n", forecast)

	// get data scada 15 min
	scadaSrc := []tk.M{}
	scada := tk.M{}
	totalDirs := 0

	scadaChan := make(chan []tk.M)
	chanDone := make(chan bool)

	if scada15minpath != "" {
		if _, err := os.Stat(scada15minpath); !os.IsNotExist(err) {
			dirsToRead := getDirectoriesToRead(tStart, tEnd)
			if len(dirsToRead) > 0 {
				dirs := []string{}
				for _, dir := range dirsToRead {
					pathLoc := filepath.Join(scada15minpath, strings.ToLower(project), dir)
					if _, err := os.Stat(pathLoc); err == nil {
						totalDirs++
						dirs = append(dirs, pathLoc)
					}
				}

				scadaChan = make(chan []tk.M, totalDirs)
				go func() {
					for {
						scs, ok := <-scadaChan
						if ok {
							for _, s := range scs {
								scadaSrc = append(scadaSrc, s)
							}
						} else {
							chanDone <- true
							return
						}
					}
				}()

				if totalDirs > 0 {
					for _, dir := range dirs {
						wg.Add(1)
						go readFiles(&wg, scadaChan, dir, project, p.Turbine, pcSrc)
					}
				}

				if totalDirs > 0 {
					wg.Wait()
				}
				close(scadaChan)
				<-chanDone
			}
		}
	}

	if len(scadaSrc) > 0 {
		for _, src := range scadaSrc {
			for key, s := range src {
				scada.Set(key, s.(tk.M))
			}
		}
	}

	turbineDown := 0
	today, _ := time.Parse("2006-01-02", time.Now().Format("2006-01-02"))
	dateend, _ := time.Parse("2006-01-02 15:04:05", "0001-01-01 00:00:00")
	matches = []tk.M{
		tk.M{"projectname": project},
		//tk.M{"$or": []tk.M{
		tk.M{"dateinfostart.dateid": tk.M{"$eq": today}},
		tk.M{"dateinfoend.dateid": tk.M{"$eq": dateend}},
		tk.M{"isdeleted": false},
		//}},
	}
	pipes = []tk.M{
		tk.M{"$match": tk.M{"$and": matches}},
		tk.M{"$group": tk.M{
			"_id":   "$turbine",
			"total": tk.M{"$sum": 1},
		}},
	}
	csrtd, e := DBRealtime().NewQuery().
		From("AlarmHFD").
		Command("pipe", pipes).
		Cursor(nil)
	defer csrtd.Close()

	dtDowns := []tk.M{}
	e = csrtd.Fetch(&dtDowns, 0, false)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}
	if len(dtDowns) > 0 {
		turbineDown = len(dtDowns) //[0].GetInt("total")
	}

	// get total production from the realtime
	defaultValue := -999999.0
	powerRtd := defaultValue
	tsRtd := time.Time{}
	matches = []tk.M{
		tk.M{"projectname": project},
		tk.M{"timestamp": tk.M{"$gt": tStart}},
		tk.M{"tags": tk.M{"$eq": "ActivePower_kW"}},
	}
	pipes = []tk.M{
		tk.M{"$match": tk.M{"$and": matches}},
		tk.M{"$group": tk.M{
			"_id":       "$projectname",
			"total":     tk.M{"$sum": "$value"},
			"timestamp": tk.M{"$max": "$timestamp"},
		}},
	}
	csrpwr, e := DBRealtime().NewQuery().
		From("ScadaRealTimeNew").
		Command("pipe", pipes).
		Cursor(nil)
	defer csrpwr.Close()

	dtPwrRtd := []tk.M{}
	e = csrpwr.Fetch(&dtPwrRtd, 0, false)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}
	if len(dtPwrRtd) > 0 {
		powerRtd = dtPwrRtd[0].GetFloat64("total") / 1000.0
		tsRtd = dtPwrRtd[0].Get("timestamp", time.Time{}).(time.Time)
	}

	dataReturn := []tk.M{}
	for _, tp := range timeperiods {
		tpkey := tp.TimePeriod.Format("20060102_150405")
		dtForecast := tk.M{}
		dtScada := tk.M{}
		if forecast.Has(tpkey) {
			dtForecast = forecast.Get(tpkey).(tk.M)
		}
		if scada.Has(tpkey) {
			dtScada = scada.Get(tpkey).(tk.M)
		}
		avacap := defaultValue
		fcvalue := defaultValue
		schval := defaultValue
		expprod := defaultValue
		actual := defaultValue
		fcastws := defaultValue
		actualws := defaultValue
		devfcast := defaultValue
		devsch := defaultValue
		dsmpenalty := ""
		deviation := defaultValue
		//isschvalavg := true
		isedited := 0
		id := tk.Sprintf("%s_%v_%s", project, tp.TimeBlock, tpkey)

		if len(dtScada) > 0 {
			actual = dtScada.GetFloat64("power") / 1000
			actualws = dtScada.GetFloat64("windspeed")
			expprod = dtScada.GetFloat64("pcstd") / 1000
		}

		if len(dtForecast) > 0 {
			avacap = dtForecast.GetFloat64("avgcapacity")
			fcvalue = dtForecast.GetFloat64("schcapacity")
			schval = dtForecast.GetFloat64("schcapacity")
			// calculate sch value from actual power of 15 min data
			// if actual != defaultValue {
			// 	schval = (actual + fcvalue) / 2
			// }
			// calculate sch value from the realtime actual power
			if powerRtd != defaultValue {
				if !tsRtd.IsZero() {
					// check if the data coming within 10 mins in india time
					if timeNowIst.Sub(tsRtd).Minutes() <= 10 {
						// put the calculation using actual power from realtime for current cell edit allowed into next 6 time block
						if tp.TimePeriod.Sub(timeNowIst).Hours() >= 1 && tp.TimePeriod.Sub(timeNowIst).Hours() <= 2.5 {
							schval = (powerRtd + fcvalue) / 2
						}
					}
				}
			}
			if dtForecast.Has("schsdlc") {
				schval = dtForecast.GetFloat64("schsdlc")
				//isschvalavg = false
			}
			if dtForecast.Has("isedited") {
				isedited = dtForecast.GetInt("isedited")
			}
		}

		// cap value for sch
		// if isschvalavg && isedited < 1 {
		if isedited < 1 && schval != defaultValue {
			if schval < 8 {
				schval = 8
			}
			if schval > 52 {
				schval = 52
			}
		}

		actualsub := 0.0
		fcvaluesub := 0.0
		schvalsub := 0.0
		if fcvalue != defaultValue {
			fcvaluesub = fcvalue
		}
		if schval != defaultValue {
			schvalsub = schval
		}
		if actual >= 0 {
			actualsub = actual
		}

		// deviation calculation
		if actual != defaultValue {
			deviation = math.Abs(actualsub - schvalsub)

			if avacap != defaultValue {
				devfcast = (actualsub - fcvaluesub) / avacap
				devsch = (actualsub - schvalsub) / avacap
			}
		}

		dateToShow := tp.TimePeriod.Format("02/01/2006")
		if tp.TimeRange == "23:45 - 00:00" {
			dateToShow = tp.TimePeriod.AddDate(0, 0, -1).Format("02/01/2006")
		}
		item := tk.M{
			"ID":            id,
			"Date":          tp.TimePeriod.Format("02/01/2006"),
			"DateToShow":    dateToShow,
			"TimeBlock":     tp.TimeRange,
			"TimeStamp":     tp.TimePeriod.Format("2006-01-02 15:04"),
			"TimeBlockInt":  tp.TimeBlock,
			"AvaCap":        avacap,
			"Forecast":      fcvalue,
			"SchFcast":      schval,
			"ExpProd":       expprod,
			"Actual":        actual,
			"FcastWs":       fcastws,
			"ActualWs":      actualws,
			"DevFcast":      devfcast,
			"DevSchAct":     devsch,
			"DSMPenalty":    dsmpenalty,
			"Deviation":     deviation,
			"TurbineDown":   turbineDown,
			"LatestSubject": latestSubject,
		}
		if item.GetFloat64("AvaCap") == defaultValue {
			item.Set("AvaCap", nil)
		}
		if item.GetFloat64("Forecast") == defaultValue {
			item.Set("Forecast", nil)
		}
		if item.GetFloat64("SchFcast") == defaultValue {
			item.Set("SchFcast", nil)
		}
		if item.GetFloat64("ExpProd") == defaultValue {
			item.Set("ExpProd", nil)
		}
		if item.GetFloat64("Actual") == defaultValue {
			item.Set("Actual", nil)
		}
		if item.GetFloat64("FcastWs") == defaultValue {
			item.Set("FcastWs", nil)
		}
		if item.GetFloat64("ActualWs") == defaultValue {
			item.Set("ActualWs", nil)
		}
		if item.GetFloat64("DevFcast") == defaultValue {
			item.Set("DevFcast", nil)
		}
		if item.GetFloat64("DevSchAct") == defaultValue {
			item.Set("DevSchAct", nil)
		}
		if item.GetFloat64("Deviation") == defaultValue {
			item.Set("Deviation", nil)
		}
		dataReturn = append(dataReturn, item)
	}

	return helper.CreateResult(true, dataReturn, "")
}

func (m *ForecastController) SendMail(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	var wg sync.WaitGroup

	type Payload struct {
		Period  string
		Date    time.Time
		Turbine []interface{}
		Project string
		Tipe    string
		Subject string
	}

	p := new(Payload)
	e := k.GetPayload(&p)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	tStart, tEnd, e := helper.GetStartEndDate(k, p.Period, p.Date, p.Date)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}
	tStart = tStart.Add(time.Duration(15) * time.Minute)
	tEnd = tEnd.Add(time.Duration(15) * time.Minute)
	timeperiods := get15MinPeriod(tStart, tEnd)
	project := p.Project
	timeNowIst := time.Now().UTC().Add(time.Duration(330) * time.Minute)

	getscada15minpath := GetConfig("scada15min_path", "")
	scada15minpath := ""
	if getscada15minpath == "" || getscada15minpath == nil {
		scada15minpath = "/mnt/data/ostrorealtime/scada15minrev/data"
	} else {
		scada15minpath = tk.ToString(getscada15minpath)
	}

	// get pc reff
	csrpc, e := DB().Connection.NewQuery().From("ref_powercurve").
		Where(dbox.Eq("model", project)).
		Order("windspeed").Cursor(nil)
	defer csrpc.Close()
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	pcSrc := []tk.M{}
	e = csrpc.Fetch(&pcSrc, 0, false)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	latestSubject := []string{}
	// get latest subject
	subToday, _ := time.Parse("2006-01-02", time.Now().Format("2006-01-02"))
	subNextday := subToday.AddDate(0, 0, 1)
	matchessub := []tk.M{
		tk.M{"projectname": project},
		tk.M{"dateinfo.dateid": tk.M{"$gte": subToday}},
		tk.M{"dateinfo.dateid": tk.M{"$lte": subNextday}},
	}
	pipessub := []tk.M{
		tk.M{"$match": tk.M{"$and": matchessub}},
		tk.M{"$group": tk.M{
			"_id": tk.M{
				"projectname": "$projectname",
				"dateid":      "$dateinfo.dateid",
			},
			"max_subject": tk.M{"$max": "$mailsubject"},
		}},
		tk.M{"$sort": tk.M{"max_subject": 1}},
	}
	csrsub, e := DB().Connection.NewQuery().
		From(new(ForecastData).TableName()).
		Command("pipe", pipessub).
		Cursor(nil)
	defer csrsub.Close()

	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}
	forecastsubject := []tk.M{}
	e = csrsub.Fetch(&forecastsubject, 0, false)
	if len(forecastsubject) > 0 {
		for _, d := range forecastsubject {
			latestSubject = append(latestSubject, d.GetString("max_subject"))
		}
	}

	// get data forecast
	matches := []tk.M{
		tk.M{"projectname": project},
		tk.M{"timestamp": tk.M{"$gte": tStart}},
		tk.M{"timestamp": tk.M{"$lte": tEnd}},
	}
	pipes := []tk.M{
		tk.M{"$match": tk.M{"$and": matches}},
	}
	csr, e := DB().Connection.NewQuery().
		From(new(ForecastData).TableName()).
		Command("pipe", pipes).
		Cursor(nil)
	defer csr.Close()

	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}
	forecast := tk.M{}
	for {
		item := tk.M{}
		e = csr.Fetch(&item, 1, false)
		if e != nil {
			break
		}
		timestamp := item.Get("timestamp", time.Time{}).(time.Time).UTC() //.UTC().Add(time.Duration(330) * time.Minute)
		if !timestamp.IsZero() {
			forecast.Set(timestamp.Format("20060102_150405"), item)
		}
	}
	// tk.Printf("%#v\n", forecast)

	// get data scada 15 min
	scadaSrc := []tk.M{}
	scada := tk.M{}
	totalDirs := 0

	scadaChan := make(chan []tk.M)
	chanDone := make(chan bool)

	if scada15minpath != "" {
		if _, err := os.Stat(scada15minpath); err == nil {
			dirsToRead := getDirectoriesToRead(tStart, tEnd)
			if len(dirsToRead) > 0 {
				dirs := []string{}
				for _, dir := range dirsToRead {
					pathLoc := filepath.Join(scada15minpath, strings.ToLower(project), dir)
					if _, err := os.Stat(pathLoc); err == nil {
						totalDirs++
						dirs = append(dirs, pathLoc)
					}
				}

				scadaChan = make(chan []tk.M, totalDirs)
				go func() {
					for {
						scs, ok := <-scadaChan
						if ok {
							for _, s := range scs {
								scadaSrc = append(scadaSrc, s)
							}
						} else {
							chanDone <- true
							return
						}
					}
				}()

				if totalDirs > 0 {
					for _, dir := range dirs {
						wg.Add(1)
						go readFiles(&wg, scadaChan, dir, project, p.Turbine, pcSrc)
					}
				}
			}
		}

		if totalDirs > 0 {
			wg.Wait()
		}
		close(scadaChan)
		<-chanDone
	}

	if len(scadaSrc) > 0 {
		for _, src := range scadaSrc {
			for key, s := range src {
				scada.Set(key, s.(tk.M))
			}
		}
	}

	turbineDown := 0
	today, _ := time.Parse("2006-01-02", time.Now().Format("2006-01-02"))
	dateend, _ := time.Parse("2006-01-02 15:04:05", "0001-01-01 00:00:00")
	matches = []tk.M{
		tk.M{"projectname": project},
		//tk.M{"$or": []tk.M{
		tk.M{"dateinfostart.dateid": tk.M{"$eq": today}},
		tk.M{"dateinfoend.dateid": tk.M{"$eq": dateend}},
		tk.M{"isdeleted": false},
		//}},
	}
	pipes = []tk.M{
		tk.M{"$match": tk.M{"$and": matches}},
		tk.M{"$group": tk.M{
			"_id":   "$turbine",
			"total": tk.M{"$sum": 1},
		}},
	}
	csrtd, e := DBRealtime().NewQuery().
		From("AlarmHFD").
		Command("pipe", pipes).
		Cursor(nil)
	defer csrtd.Close()

	dtDowns := []tk.M{}
	e = csrtd.Fetch(&dtDowns, 0, false)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}
	if len(dtDowns) > 0 {
		turbineDown = len(dtDowns) //[0].GetInt("total")
	}

	// get total production from the realtime
	defaultValue := -999999.0
	powerRtd := defaultValue
	tsRtd := time.Time{}
	matches = []tk.M{
		tk.M{"projectname": project},
		tk.M{"timestamp": tk.M{"$gt": tStart}},
		tk.M{"tags": tk.M{"$eq": "ActivePower_kW"}},
	}
	pipes = []tk.M{
		tk.M{"$match": tk.M{"$and": matches}},
		tk.M{"$group": tk.M{
			"_id":       "$projectname",
			"total":     tk.M{"$sum": "$value"},
			"timestamp": tk.M{"$max": "$timestamp"},
		}},
	}
	csrpwr, e := DBRealtime().NewQuery().
		From("ScadaRealTimeNew").
		Command("pipe", pipes).
		Cursor(nil)
	defer csrpwr.Close()

	dtPwrRtd := []tk.M{}
	e = csrpwr.Fetch(&dtPwrRtd, 0, false)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}
	if len(dtPwrRtd) > 0 {
		powerRtd = dtPwrRtd[0].GetFloat64("total") / 1000.0
		tsRtd = dtPwrRtd[0].Get("timestamp", time.Time{}).(time.Time)
	}

	dataReturn := []tk.M{}
	for _, tp := range timeperiods {
		tpkey := tp.TimePeriod.Format("20060102_150405")
		dtForecast := tk.M{}
		dtScada := tk.M{}
		if forecast.Has(tpkey) {
			dtForecast = forecast.Get(tpkey).(tk.M)
		}
		if scada.Has(tpkey) {
			dtScada = scada.Get(tpkey).(tk.M)
		}
		avacap := defaultValue
		fcvalue := defaultValue
		schval := defaultValue
		expprod := defaultValue
		actual := defaultValue
		fcastws := defaultValue
		actualws := defaultValue
		devfcast := defaultValue
		devsch := defaultValue
		dsmpenalty := ""
		deviation := defaultValue
		// isschvalavg := true
		id := tk.Sprintf("%s_%v_%s", project, tp.TimeBlock, tpkey)

		if len(dtScada) > 0 {
			actual = dtScada.GetFloat64("power") / 1000
			actualws = dtScada.GetFloat64("windspeed")
			expprod = dtScada.GetFloat64("pcstd") / 1000
		}

		if len(dtForecast) > 0 {
			avacap = dtForecast.GetFloat64("avgcapacity")
			fcvalue = dtForecast.GetFloat64("schcapacity")
			schval = dtForecast.GetFloat64("schcapacity")
			// calculate sch value from actual power of 15 min data
			if actual != defaultValue {
				schval = (actual + fcvalue) / 2
			}
			// calculate sch value from the realtime actual power
			if powerRtd != defaultValue {
				if !tsRtd.IsZero() {
					// check if the data coming within 10 mins in india time
					if timeNowIst.Sub(tsRtd).Minutes() <= 10 {
						// put the calculation using actual power from realtime for current cell edit allowed into next 6 time block
						if tp.TimePeriod.Sub(timeNowIst).Hours() >= 1 && tp.TimePeriod.Sub(timeNowIst).Hours() <= 2.5 {
							schval = (powerRtd + fcvalue) / 2
						}
					}
				}
			}
			if dtForecast.Has("schsdlc") {
				schval = dtForecast.GetFloat64("schsdlc")
				// isschvalavg = false
			}
		}

		// cap value for sch
		// if isschvalavg {
		if schval < 8 {
			schval = 8
		}
		if schval > 52 {
			schval = 52
		}
		// }

		actualsub := 0.0
		fcvaluesub := 0.0
		schvalsub := 0.0
		if fcvalue != defaultValue {
			fcvaluesub = fcvalue
		}
		if schval != defaultValue {
			schvalsub = schval
		}
		if actual >= 0 {
			actualsub = actual
		}

		// deviation calculation
		if actual != defaultValue {
			deviation = math.Abs(actualsub - schvalsub)

			if avacap != defaultValue {
				devfcast = (actualsub - fcvaluesub) / avacap
				devsch = (actualsub - schvalsub) / avacap
			}
		}

		item := tk.M{
			"ID":            id,
			"Date":          tp.TimePeriod.Format("02/01/2006"),
			"TimeBlock":     tp.TimeRange,
			"TimeStamp":     tp.TimePeriod,
			"TimeBlockInt":  tp.TimeBlock,
			"AvaCap":        avacap,
			"Forecast":      fcvalue,
			"SchFcast":      schval,
			"ExpProd":       expprod,
			"Actual":        actual,
			"FcastWs":       fcastws,
			"ActualWs":      actualws,
			"DevFcast":      devfcast,
			"DevSchAct":     devsch,
			"DSMPenalty":    dsmpenalty,
			"Deviation":     deviation,
			"TurbineDown":   turbineDown,
			"LatestSubject": latestSubject,
		}
		if item.GetFloat64("AvaCap") == defaultValue {
			item.Set("AvaCap", nil)
		}
		if item.GetFloat64("Forecast") == defaultValue {
			item.Set("Forecast", nil)
		}
		if item.GetFloat64("SchFcast") == defaultValue {
			item.Set("SchFcast", nil)
		}
		if item.GetFloat64("ExpProd") == defaultValue {
			item.Set("ExpProd", nil)
		}
		if item.GetFloat64("Actual") == defaultValue {
			item.Set("Actual", nil)
		}
		if item.GetFloat64("FcastWs") == defaultValue {
			item.Set("FcastWs", nil)
		}
		if item.GetFloat64("ActualWs") == defaultValue {
			item.Set("ActualWs", nil)
		}
		if item.GetFloat64("DevFcast") == defaultValue {
			item.Set("DevFcast", nil)
		}
		if item.GetFloat64("DevSchAct") == defaultValue {
			item.Set("DevSchAct", nil)
		}
		if item.GetFloat64("Deviation") == defaultValue {
			item.Set("Deviation", nil)
		}
		dataReturn = append(dataReturn, item)
	}

	sendMail := true
	msg := ""
	if len(forecast) > 0 {
		from := tk.M{"email": "ostro.support@eaciit.com", "name": "Ostro Support"}
		tos := []tk.M{
			tk.M{"email": "oms@ostro.in", "name": "OMS"},
			tk.M{"email": "priyadarshan.b@ostro.in", "name": "Priyadarshan B"},
		}
		ccs := []tk.M{
			tk.M{"email": "sandhya@eaciit.com", "name": "Sandhya Jain"},
			tk.M{"email": "aris.meika@eaciit.com", "name": "Aris Meika"},
			tk.M{"email": "shreyas@eaciit.com", "name": "Shreyas Mithare"},
		}
		bccs := []tk.M{}
		createXlsAndSend(p.Project, p.Date, p.Subject, from, tos, ccs, bccs, dataReturn)
	} else {
		revNos := strings.Split(strings.ToLower(p.Subject), "rev")
		revNo := strings.TrimSpace(tk.Sprintf("%s ", revNos[1])[0:3])
		revNo2Digit := tk.Sprintf("%02v", revNo)
		sendMail = false
		msg = "No forecast data for " + p.Project + " for " + p.Date.Format("02/01/2006") + " rev no " + revNo2Digit
	}

	return helper.CreateResult(sendMail, []tk.M{}, msg)
}

func createXlsAndSend(project string, date time.Time, subject string, addressFrom tk.M, addressTo []tk.M, addressCc []tk.M, addressBcc []tk.M, data []tk.M) {
	revNos := strings.Split(strings.ToLower(subject), "rev")
	revNo := strings.TrimSpace(tk.Sprintf("%s ", revNos[1])[0:3])
	revNo2Digit := tk.Sprintf("%02v", revNo)
	pathToSave := tk.Sprintf("web/assets/forecast/%s", date.Format("20060102"))
	_, err := os.Stat(pathToSave)
	if os.IsNotExist(err) {
		err = os.MkdirAll(pathToSave, 0775)
	}

	if err != nil {
		tk.Println("Error creating directories:" + err.Error())
	}

	filename := tk.Sprintf("Forecast_%s_%s_Rev_%s.xlsx", project, date.Format("20060102"), revNo2Digit)
	filetosave := filepath.Join(pathToSave, filename)

	newXls := xls.NewFile()
	sheet, err := newXls.AddSheet(tk.Sprintf("Forecast Rev %s for %s", revNo, project))
	if err != nil {
		tk.Printf("Error adding new sheet: " + err.Error())
	}

	titleInfo := []tk.M{
		tk.M{"title": "Name of the pooling station", "value": "CHIKKOPPA"},
		tk.M{"title": "Date", "value": date.Format("02/01/2006")},
		tk.M{"title": "Installed Capacity", "value": "60 MW"},
		tk.M{"title": "Forecast provided by", "value": "Ostro Mahawind Power Pvt Ltd (Ostro Energy P LTD)"},
		tk.M{"title": "Developer of Pooling station", "value": "GAMESA (Ostro)"},
		tk.M{"title": "KPTCL/ESCOM Injecting station", "value": "KULGOD"},
		tk.M{"title": "Voltage level at injecting point", "value": "110/33/11 KV"},
		tk.M{"title": "Revision Number ", "value": revNo},
	}

	border := *xls.NewBorder("thin", "thin", "thin", "thin")
	font := *xls.NewFont(11, "Calibri")
	font.Bold = true

	headerStyle := xls.NewStyle()
	headerStyle.Alignment.Horizontal = "center"
	headerStyle.Alignment.Vertical = "center"
	headerStyle.ApplyAlignment = true
	headerStyle.Border = border
	headerStyle.ApplyBorder = true
	headerStyle.Font = font
	headerStyle.ApplyFont = true

	headerStyle2 := xls.NewStyle()
	headerStyle2.Alignment.Horizontal = "center"
	headerStyle2.Alignment.Vertical = "center"
	headerStyle2.ApplyAlignment = true
	headerStyle2.Border = border
	headerStyle2.ApplyBorder = true
	headerStyle2.Font = font
	headerStyle2.ApplyFont = true

	itemStyle := xls.NewStyle()
	itemStyle.Border = border
	itemStyle.ApplyBorder = true
	itemStyle.Alignment.Horizontal = "center"
	itemStyle.Alignment.Vertical = "center"
	itemStyle.ApplyAlignment = true
	itemStyle.Font = font
	itemStyle.Font.Bold = false
	itemStyle.ApplyFont = true

	itemStyle2 := xls.NewStyle()
	itemStyle2.Border = border
	itemStyle2.ApplyBorder = true
	itemStyle2.Alignment.Horizontal = "right"
	itemStyle2.Alignment.Vertical = "center"
	itemStyle2.ApplyAlignment = true
	itemStyle2.Font = font
	itemStyle2.Font.Bold = false
	itemStyle2.ApplyFont = true

	for _, t := range titleInfo {
		tRow := sheet.AddRow()
		tRow.AddCell()
		tRow.AddCell()
		tRow.AddCell()

		tRow.Cells[0].Merge(1, 0)
		tRow.Cells[0].Value = t.GetString("title")
		xStyle := xls.NewStyle()
		xStyle.Font = font
		xStyle.ApplyFont = true
		tRow.Cells[0].SetStyle(xStyle)

		tRow.Cells[2].Value = tk.Sprintf("%s", t.GetString("value"))
		yStyle := xls.NewStyle()
		yStyle.Font = font
		yStyle.Font.Bold = false
		yStyle.ApplyFont = true
		tRow.Cells[2].SetStyle(yStyle)
	}

	sheet.AddRow()
	sheet.AddRow()

	row0 := sheet.AddRow()
	row1 := sheet.AddRow()

	for i := 0; i < 5; i++ {
		row0.AddCell()
		row1.AddCell()
	}

	row0.Cells[0].Merge(0, 1)
	row0.Cells[0].Value = "Date"
	row0.Cells[0].SetStyle(headerStyle)
	row0.Cells[1].Merge(0, 1)
	row0.Cells[1].Value = "Time"
	row0.Cells[1].SetStyle(headerStyle)
	row0.Cells[2].Merge(0, 1)
	row0.Cells[2].Value = "Time Block"
	row0.Cells[2].SetStyle(headerStyle)
	row0.Cells[3].Merge(1, 0)
	row0.Cells[3].Value = "Rev"
	row0.Cells[3].SetStyle(headerStyle)

	row1.Cells[3].Value = "AvC"
	row1.Cells[3].SetStyle(headerStyle2)
	row1.Cells[4].Value = "SCH"
	row1.Cells[4].SetStyle(headerStyle2)

	if len(data) > 0 {
		for _, d := range data {
			currentRow := sheet.AddRow()
			for i := 0; i < 5; i++ {
				currentRow.AddCell()
			}

			sdate := d.GetString("Date")
			stimerange := d.GetString("TimeBlock")
			timeblock := d.GetInt("TimeBlockInt")
			avc := d.GetFloat64("AvaCap")
			sch := d.GetFloat64("SchFcast")

			currentRow.Cells[0].Value = sdate
			currentRow.Cells[0].SetStyle(itemStyle)
			currentRow.Cells[1].Value = stimerange
			currentRow.Cells[1].SetStyle(itemStyle)
			currentRow.Cells[2].SetInt(timeblock)
			currentRow.Cells[2].SetStyle(itemStyle)
			currentRow.Cells[3].SetFloatWithFormat(avc, "#,##0")
			currentRow.Cells[3].SetStyle(itemStyle)
			currentRow.Cells[4].SetFloatWithFormat(sch, "#,##0.00")
			currentRow.Cells[4].SetStyle(itemStyle2)
		}
	}

	sheet.Cols[0].Width = 16.67
	sheet.Cols[1].Width = 14.33
	sheet.Cols[2].Width = 16.83
	sheet.Cols[3].Width = 7.33
	sheet.Cols[4].Width = 7.33

	newXls.Save(filetosave)

	mailSubject := tk.Sprintf("%s Schedule Forecast For %s Rev %s", project, date.Format("02/01/2006"), revNo2Digit)
	mailContent := tk.Sprintf("<p>Dear All,</p><p>Please find the attachment for %s scheduler forecast for %s revision number %s.</p><p>&nbsp;<br /></p><p>Thank you.</p>", project, date.Format("02/01/2006"), revNo2Digit)
	err = sendEmail(mailSubject, addressFrom, addressTo, addressCc, addressBcc, mailContent, filetosave)
	if err != nil {
		tk.Printf("Error send email : %s \n", err.Error())
	}
}

func sendEmail(subject string, addressFrom tk.M, addressTo []tk.M, addressCc []tk.M, addressBcc []tk.M, content string, filename string) error {
	if subject != "" && len(addressTo) > 0 && len(addressFrom) > 0 {
		mailAddress := addressFrom.GetString("email")
		mailAddressPassword := "LfsAgXd546xmRGQJ"
		mailHost := "smtp.outlook.com"
		mailPort := 587

		mailTos := []string{}
		mailCcs := []string{}
		mailBccs := []string{}

		mail := gomail.NewMessage()
		mail.SetAddressHeader("From", mailAddress, addressFrom.GetString("name"))
		for _, to := range addressTo {
			mailTos = append(mailTos, to.GetString("email"))
		}
		if len(addressCc) > 0 {
			for _, cc := range addressCc {
				mailCcs = append(mailCcs, cc.GetString("email"))
			}
		}
		if len(addressBcc) > 0 {
			for _, bcc := range addressBcc {
				mailBccs = append(mailBccs, bcc.GetString("email"))
			}
		}
		mail.SetHeader("To", mailTos...)
		if len(mailCcs) > 0 {
			mail.SetHeader("Cc", mailCcs...)
		}
		if len(mailBccs) > 0 {
			mail.SetHeader("Cc", mailBccs...)
		}
		mail.SetHeader("Subject", subject)
		mail.SetBody("text/html", content)
		mail.Attach(filename)

		mailDial := gomail.NewPlainDialer(mailHost, mailPort, mailAddress, mailAddressPassword)
		err := mailDial.DialAndSend(mail)

		return err
	}

	return nil
}

func readFiles(wg *sync.WaitGroup, scadaChan chan []tk.M, pathLoc string, project string, turbines []interface{}, pcSrc []tk.M) {
	defer wg.Done()
	files := []string{}
	filepath.Walk(pathLoc, func(path string, info os.FileInfo, err error) error {
		files = append(files, path)
		return nil
	})
	if len(files) > 0 {
		data := []tk.M{}
		for _, file := range files {
			if strings.Contains(file, ".csv") {
				dataRows, _ := readLines(file)
				if len(dataRows) > 0 {
					items := tk.M{}
					for _, dt := range dataRows {
						dataMtx.Lock()
						lines := strings.Split(dt, ",")
						turbine := lines[0]
						ts, _ := time.Parse("2006-01-02 15:04:05", tk.Sprintf("%s %s", lines[1], lines[2]))
						power := tk.ToFloat64(lines[3], 8, tk.RoundingAuto)
						ws := tk.ToFloat64(lines[6], 8, tk.RoundingAuto)
						isexists, _ := inArray(turbine, turbines)
						stc, std := GetPowerCurveTkMSource(pcSrc, ws)
						key := ts.Format("20060102_150405")
						if isexists {
							if items.Has(key) {
								item := items.Get(key, tk.M{}).(tk.M)
								currpower := item.GetFloat64("power")
								currpower += power
								totalws := item.GetFloat64("totalws")
								totalws += ws
								currcount := item.GetInt("count")
								currcount++
								currstc := item.GetFloat64("pcspc")
								currstc += stc
								currstd := item.GetFloat64("pcstd")
								currstd += std

								currws := tk.Div(totalws, tk.ToFloat64(currcount, 4, tk.RoundingAuto))
								items.Unset(key)
								items.Set(key, tk.M{
									"timestamp": ts,
									"power":     currpower,
									"windspeed": currws,
									"totalws":   totalws,
									"pcspc":     currstc,
									"pcstd":     currstd,
									"count":     currcount,
								})
							} else {
								//stc, std := 0.0, 0.0 //GetPowerCurveByWs(DB().Connection, project, ws)
								items.Set(key, tk.M{
									"timestamp": ts,
									"power":     power,
									"windspeed": ws,
									"totalws":   ws,
									"pcspc":     stc,
									"pcstd":     std,
									"count":     1,
								})
							}
						}
						dataMtx.Unlock()
					}
					data = append(data, items)
				}
			}
		}
		scadaChan <- data
	}
}

func readLines(filename string) ([]string, error) {
	var lines []string
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		return lines, err
	}
	buf := bytes.NewBuffer(file)
	for {
		line, err := buf.ReadString('\n')
		if len(line) == 0 {
			if err != nil {
				if err == io.EOF {
					break
				}
				return lines, err
			}
		}
		lines = append(lines, line)
		if err != nil && err != io.EOF {
			return lines, err
		}
	}
	return lines, nil
}

func inArray(val interface{}, array interface{}) (exists bool, index int) {
	exists = false
	index = -1

	switch reflect.TypeOf(array).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(array)

		for i := 0; i < s.Len(); i++ {
			if reflect.DeepEqual(val, s.Index(i).Interface()) == true {
				index = i
				exists = true
				return
			}
		}
	}

	return
}
