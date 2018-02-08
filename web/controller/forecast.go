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
		timestamp := item.Get("timestamp", time.Time{}).(time.Time)
		if !timestamp.IsZero() {
			forecast.Set(timestamp.Format("20060102_150405"), item)
		}
	}

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

	defaultValue := -999999.0
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
		devfcast := 0.0
		devsch := 0.0
		dsmpenalty := ""
		deviation := defaultValue

		if len(dtForecast) > 0 {
			avacap = dtForecast.GetFloat64("avgcapacity")
			fcvalue = dtForecast.GetFloat64("schcapacity")
			schval = dtForecast.GetFloat64("schcapacity")
			if dtForecast.Has("schsdlc") {
				schval = dtForecast.GetFloat64("schsdlc")
			}
		}

		if len(dtScada) > 0 {
			actual = dtScada.GetFloat64("power") / 1000
			actualws = dtScada.GetFloat64("windspeed")
			expprod = dtScada.GetFloat64("pcstd") / 1000
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
		deviation = math.Abs(actualsub - schvalsub)

		if avacap != defaultValue {
			devfcast = (actualsub - fcvaluesub) / avacap
			devsch = (actualsub - schvalsub) / avacap
		}

		item := tk.M{
			"Date":         tp.TimePeriod.Format("02/01/2006"),
			"TimeBlock":    tp.TimeRange,
			"TimeStamp":    tp.TimePeriod,
			"TimeBlockInt": tp.TimeBlock,
			"AvaCap":       avacap,
			"Forecast":     fcvalue,
			"SchFcast":     schval,
			"ExpProd":      expprod,
			"Actual":       actual,
			"FcastWs":      fcastws,
			"ActualWs":     actualws,
			"DevFcast":     devfcast,
			"DevSchAct":    devsch,
			"DSMPenalty":   dsmpenalty,
			"Deviation":    deviation,
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
		if item.GetFloat64("Deviation") == defaultValue {
			item.Set("Deviation", nil)
		}
		dataReturn = append(dataReturn, item)
	}

	return helper.CreateResult(true, dataReturn, "")
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
