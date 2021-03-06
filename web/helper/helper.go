package helper

import (
	. "eaciit/wfdemo-git/library/core"
	hp "eaciit/wfdemo-git/library/helper"
	md "eaciit/wfdemo-git/library/models"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"gopkg.in/mgo.v2/bson"

	"github.com/eaciit/dbox"
	"github.com/eaciit/knot/knot.v1"
	"github.com/eaciit/toolkit"
)

var (
	DebugMode bool
	DownTypes = []toolkit.M{
		{"down": "aebok", "label": "AEBok"},
		{"down": "externalstop", "label": "External Stop"},
		{"down": "griddown", "label": "Grid Down"},
		{"down": "internalgrid", "label": "Internal Grid"},
		{"down": "machinedown", "label": "Machine Down"},
		{"down": "unknown", "label": "Unknown"},
		{"down": "weatherstop", "label": "Weather Stop"},
	}
	WC *knot.WebContext
)

type PayloadsDB struct {
	Project   string
	Turbine   []interface{}
	DateStart time.Time
	DateEnd   time.Time
	Skip      int
	Take      int
	Sort      []Sorting
	Filter    *FilterJS `json:"filter"`
	Misc      toolkit.M `json:"misc"`
	Custom    toolkit.M `json:"custom"`
}

type Payloads struct {
	Project string
	Skip    int
	Take    int
	Sort    []Sorting
	Filter  *FilterJS `json:"filter"`
	Misc    toolkit.M `json:"misc"`
	Custom  toolkit.M `json:"custom"`
}

type Sorting struct {
	Field string
	Dir   string
}

type FilterJS struct {
	Filters []*Filter `json:"filters"`
	Logic   string
}

type Filter struct {
	Field   string      `json:"field"`
	Op      string      `json:"operator"`
	Value   interface{} `json:"value"`
	Filters []Filter    `json:"filters"`
}

func (s *Payloads) ParseFilter() (filters []*dbox.Filter, err error) {
	if s != nil {
		for _, each := range s.Filter.Filters {
			filtersTmp := []*dbox.Filter{}
			filtersTmp, err = doParseFilter(each, s)
			for _, eachTmp := range filtersTmp {
				filters = append(filters, eachTmp)
			}
		}
	}

	/*for _, val := range filters {
		log.Printf("filter: %#v \n", val)
		if val.Field == "timestamp" {
			log.Printf("timestamp: %#v \n", val.Value.(time.Time).String())
		}
	}*/

	return
}

func doParseFilter(each *Filter, s *Payloads) (filters []*dbox.Filter, err error) {
	datelist := []string{
		"timestamp",
		"dateinfo.dateid",
		"startdate",
		"timestart",
	}

	field := strings.ToLower(each.Field)

	if each.Filters != nil || len(each.Filters) > 0 {
		for _, eachF := range each.Filters {
			filtersTmp := []*dbox.Filter{}
			filtersTmp, err = doParseFilter(&eachF, s)
			for _, eachTmp := range filtersTmp {
				filters = append(filters, eachTmp)
			}
		}
	} else {
		switch each.Op {
		case "gte":
			var value interface{} = each.Value
			if toolkit.TypeName(value) == "string" {
				if value.(string) != "" {
					if toolkit.HasMember(datelist, field) {
						var t time.Time
						b, err := time.Parse("2006-01-02T15:04:05.000Z", value.(string))
						if err != nil {
							b, err = time.Parse("2006-01-02 15:04:05", value.(string))
							if err != nil {
								toolkit.Println(err.Error())
							}
						}
						if s.Misc.Has("period") {
							t, _, err = GetStartEndDate(s.Misc["knot_data"].(*knot.WebContext), s.Misc.GetString("period"), b, b)
						} else {
							t, _ = time.Parse("2006-01-02 15:04:05", b.UTC().Format("2006-01-02")+" 00:00:00")
						}
						value = t
					}
					filters = append(filters, dbox.Gte(field, value))
				}
			} else {
				filters = append(filters, dbox.Gte(field, value))
			}
		case "gt":
			var value interface{} = each.Value
			if toolkit.TypeName(value) == "string" {
				if value.(string) != "" {
					if toolkit.HasMember(datelist, field) {
						var t time.Time
						b, err := time.Parse("2006-01-02T15:04:05.000Z", value.(string))
						if err != nil {
							b, err = time.Parse("2006-01-02 15:04:05", value.(string))
							if err != nil {
								toolkit.Println(err.Error())
							}
						}
						if s.Misc.Has("period") {
							t, _, err = GetStartEndDate(s.Misc["knot_data"].(*knot.WebContext), s.Misc.GetString("period"), b, b)
						} else {
							t, _ = time.Parse("2006-01-02 15:04:05", b.UTC().Format("2006-01-02")+" 00:00:00")
						}
						value = t
					}
					filters = append(filters, dbox.Gt(field, value))
				}
			} else {
				filters = append(filters, dbox.Gt(field, value))
			}
		case "lte":
			var value interface{} = each.Value

			if toolkit.TypeName(value) == "string" {
				if value.(string) != "" {
					if toolkit.HasMember(datelist, field) {
						var t time.Time
						b, err := time.Parse("2006-01-02T15:04:05.000Z", value.(string))
						if err != nil {
							b, err = time.Parse("2006-01-02 15:04:05", value.(string))
							if err != nil {
								toolkit.Println(err.Error())
							}
						}
						if s.Misc.Has("period") {
							_, t, err = GetStartEndDate(s.Misc["knot_data"].(*knot.WebContext), s.Misc.GetString("period"), b, b)
						} else {
							t, _ = time.Parse("2006-01-02 15:04:05", b.UTC().Format("2006-01-02")+" 23:59:59")
						}
						value = t
					}
					filters = append(filters, dbox.Lte(field, value))
				}
			} else {
				filters = append(filters, dbox.Lte(field, value))
			}
		case "lt":
			var value interface{} = each.Value

			if toolkit.TypeName(value) == "string" {
				if value.(string) != "" {
					if toolkit.HasMember(datelist, field) {
						var t time.Time
						b, err := time.Parse("2006-01-02T15:04:05.000Z", value.(string))
						if err != nil {
							b, err = time.Parse("2006-01-02 15:04:05", value.(string))
							if err != nil {
								toolkit.Println(err.Error())
							}
						}
						if s.Misc.Has("period") {
							_, t, err = GetStartEndDate(s.Misc["knot_data"].(*knot.WebContext), s.Misc.GetString("period"), b, b)
						} else {
							t, _ = time.Parse("2006-01-02 15:04:05", b.UTC().Format("2006-01-02")+" 23:59:59")
						}
						value = t
					}
					filters = append(filters, dbox.Lt(field, value))
				}
			} else {
				filters = append(filters, dbox.Lt(field, value))
			}
		case "eq":
			value := each.Value

			if field == "turbine" && value.(string) == "" {
				return
			} else if field == "isvalidtimeduration" && value.(bool) == true {
				return
			} else if field == "projectid" && value.(string) == "" {
				return
			}

			if field == "projectname" && value.(string) != "" {
				anProject := strings.Split(value.(string), "(")
				project := strings.TrimRight(anProject[0], " ")
				filters = append(filters, dbox.Eq(field, project))
			} else if field == "_id" && bson.IsObjectIdHex(toolkit.ToString(value)) {
				filters = append(filters, dbox.Eq(field, bson.ObjectIdHex(toolkit.ToString(value))))
			} else {
				filters = append(filters, dbox.Eq(field, value))
			}
		case "neq":
			value := each.Value
			filters = append(filters, dbox.Ne(field, value))
		case "in":
			value := each.Value
			if (field == "turbineid" && toolkit.SliceLen(value) == 0) ||
				field == "turbine" && toolkit.SliceLen(value) == 0 {
				return
			}
			filters = append(filters, dbox.In(field, value.([]interface{})...))
		case "contains":
			value := each.Value
			filters = append(filters, dbox.Contains(strings.ToLower(field), toolkit.ToString(value)))
		}
	}

	return
}

func HandleError(err error, optionalArgs ...interface{}) bool {
	if err != nil {
		toolkit.Printf("error occured: %s", err.Error())

		if len(optionalArgs) > 0 {
			optionalArgs[0].(func(bool))(false)
		}

		return false
	}

	if len(optionalArgs) > 0 {
		optionalArgs[0].(func(bool))(true)
	}

	return true
}

func CheckEnergyComparison(newdata toolkit.Ms, key1 string, key2 string) toolkit.Ms {
	// countData1 := 0
	// countData2 := 0
	result := toolkit.Ms{}
	measurement := "MWh"
	// for _, data := range newdata {
	// 	if data.GetFloat64(key1) < data.GetFloat64(key2) {
	// 		countData1++
	// 	} else {
	// 		countData2++
	// 	}
	// }

	// kunciData := ""
	// if countData1 > countData2 {
	// 	kunciData = key1
	// } else {
	// 	kunciData = key2
	// }

	// countSatuan := toolkit.M{}

	// for _, data := range newdata {
	// 	cekVal := data.GetFloat64(kunciData) / 1000000
	// 	energyType := "MWh"
	// 	if cekVal < 1 {
	// 		cekVal = data.GetFloat64(kunciData) / 1000
	// 		energyType = "MWh"
	// 		if cekVal < 1 {
	// 			cekVal = data.GetFloat64(kunciData)
	// 			energyType = "kWh"
	// 		}
	// 	}
	// 	if countSatuan.Has(energyType) {
	// 		countSatuan.Set(energyType, countSatuan.GetInt(energyType)+1)
	// 	} else {
	// 		countSatuan.Set(energyType, 1)
	// 	}
	// }

	pembagi := 1000.0
	// if (countSatuan.GetInt("GWh") > countSatuan.GetInt("MWh")) && (countSatuan.GetInt("GWh") > countSatuan.GetInt("kWh")) {
	// 	pembagi = 1000000
	// 	measurement = "GWh"
	// } else if (countSatuan.GetInt("MWh") > countSatuan.GetInt("GWh")) && (countSatuan.GetInt("MWh") > countSatuan.GetInt("kWh")) {
	// 	pembagi = 1000
	// 	measurement = "MWh"
	// } else {
	// 	pembagi = 1
	// 	measurement = "kWh"
	// }

	for _, data := range newdata {
		data.Set(key1, data.GetFloat64(key1)/pembagi)
		data.Set(key2, data.GetFloat64(key2)/pembagi)
		data.Set("measurement", measurement)
		result = append(result, data)
	}
	return result
}
func EnergyMeasurement(data interface{}, key1 string, key2 string) toolkit.Ms {
	result := toolkit.Ms{}
	newdata := toolkit.Ms{}
	if strings.Contains(toolkit.TypeName(data), "[]toolkit") {
		newdata = data.([]toolkit.M)
		result = CheckEnergyComparison(newdata, key1, key2)
	} else {
		_data := data.(toolkit.M)
		newdata = append(newdata, _data)
		result = CheckEnergyComparison(newdata, key1, key2)
	}

	return result
}

func CreateResultWithoutSession(success bool, data interface{}, message string) map[string]interface{} {
	if !success {
		toolkit.Println("ERROR! ", message)
		if DebugMode {
			panic(message)
		}
	}

	return map[string]interface{}{
		"data":    data,
		"success": success,
		"message": message,
	}
}

func CreateResult(success bool, data interface{}, message string) map[string]interface{} {
	if !success {
		toolkit.Println("ERROR! ", message)
		if DebugMode {
			panic(message)
		}
	}

	sessionid := "baypass session"
	// sessionid := WC.Session("sessionid", "")
	// if toolkit.ToString(sessionid) == "" {
	// 	sessionid = "baypass session"
	// }

	// log.Printf(">> %v \n", sessionid)

	if toolkit.ToString(sessionid) == "" {
		// if !success && data == nil && !strings.Contains(WC.Request.URL.String(), "login/processlogin") {
		if !strings.Contains(WC.Request.URL.String(), "login/processlogin") {
			dataX := struct {
				Data []toolkit.M
			}{
				Data: []toolkit.M{},
			}

			data = dataX
			success = false
			message = "Your session has expired, please login"
		}
	} else {
		if !success && data == nil {
			dataX := struct {
				Data []toolkit.M
			}{
				Data: []toolkit.M{},
			}

			data = dataX
			success = false
			message = "data is empty"
		}
	}

	return map[string]interface{}{
		"data":    data,
		"success": success,
		"message": message,
	}
}

func CreateResultX(success bool, data interface{}, message string, r *knot.WebContext) map[string]interface{} {
	if !success {
		toolkit.Println("ERROR! ", message)
		if DebugMode {
			panic(message)
		}
	}
	sessionid := r.Session("sessionid", "")

	// log.Printf(">> %v \n", sessionid)

	if toolkit.ToString(sessionid) == "" {
		// if !success && data == nil && !strings.Contains(WC.Request.URL.String(), "login/processlogin") {
		if !strings.Contains(WC.Request.URL.String(), "login/processlogin") {
			// dataX := struct {
			// 	Data []toolkit.M
			// }{
			// 	Data: []toolkit.M{},
			// }

			// data = dataX
			// success = false
			// message = "Your session has expired, please login"
		}
	} else {
		if !success && data == nil {
			dataX := struct {
				Data []toolkit.M
			}{
				Data: []toolkit.M{},
			}

			data = dataX
			success = false
			message = "data is empty"
		}
	}

	return map[string]interface{}{
		"data":    data,
		"success": success,
		"message": message,
	}
}

func ImageUploadHandler(r *knot.WebContext, filename, dstpath string) (error, string) {
	file, handler, err := r.Request.FormFile(filename)
	if err != nil {
		return err, ""
	}
	defer file.Close()

	newImageName := toolkit.RandomString(32) + filepath.Ext(handler.Filename)
	dstSource := dstpath + toolkit.PathSeparator + newImageName
	f, err := os.OpenFile(dstSource, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return err, ""
	}
	defer f.Close()
	io.Copy(f, file)

	return nil, newImageName
}

func UploadFileHandler(r *knot.WebContext, tempfile, dstpath, filename string) (error, string, string, string) {
	file, handler, err := r.Request.FormFile(tempfile)
	if err != nil {
		return err, "", "", ""
	}
	defer file.Close()

	ext := filepath.Ext(handler.Filename)
	newFileName := filename + ext
	dstSource := dstpath + toolkit.PathSeparator + newFileName
	f, err := os.OpenFile(dstSource, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return err, "", "", ""
	}
	defer f.Close()
	io.Copy(f, file)

	return nil, handler.Filename, newFileName, strings.Split(ext, ".")[1]
}

func UploadHandler(r *knot.WebContext, filename, dstpath string) (error, string) {
	file, handler, err := r.Request.FormFile(filename)
	if err != nil {
		return err, ""
	}
	defer file.Close()

	dstSource := dstpath + toolkit.PathSeparator + handler.Filename
	f, err := os.OpenFile(dstSource, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return err, ""
	}
	defer f.Close()
	_, err = io.Copy(f, file)
	if err != nil {
		return err, ""
	}
	toolkit.Println("Write file: " + dstSource)
	return nil, handler.Filename
}

func GetDayInYear(year int) toolkit.M {
	result := toolkit.M{}
	for m := time.January; m <= time.December; m++ {
		t := time.Date(year, m+1, 1, 0, 0, 0, 0, time.UTC)
		result.Set(toolkit.ToString(int(m)), t.Add(-24*time.Hour).Day())
	}
	return result
}

func GetDurationInMonth(tStart time.Time, tEnd time.Time) (int, []interface{}, toolkit.M) {
	durationMonths := 0
	monthDay := toolkit.M{}
	var months []interface{}
	xDate := tStart
	year := xDate.Year()
	month := int(xDate.Month())
	day := 1

	daysInYear := GetDayInYear(year)

	if (toolkit.ToString(xDate.Year()) + "" + toolkit.ToString(int(xDate.Month()))) != (toolkit.ToString(tEnd.Year()) + "" + toolkit.ToString(int(tEnd.Month()))) {
	out:
		for {
			xString := toolkit.ToString(xDate.Year()) + "" + toolkit.ToString(int(xDate.Month()))
			endString := toolkit.ToString(tEnd.Year()) + "" + toolkit.ToString(int(tEnd.Month()))

			if xString != endString {
				durationMonths++
				months = append(months, int(xDate.Month()))

				if (toolkit.ToString(xDate.Year()) + "" + toolkit.ToString(int(xDate.Month()))) == (toolkit.ToString(tStart.Year()) + "" + toolkit.ToString(int(tStart.Month()))) {
					monthDay.Set(toolkit.ToString(tStart.Year())+""+toolkit.ToString(int(tStart.Month())),
						toolkit.M{
							"days":         daysInYear.GetInt(toolkit.ToString(int(xDate.Month()))) - (int(tStart.Day()) - 1),
							"totalInMonth": daysInYear.GetInt(toolkit.ToString(int(xDate.Month()))),
						})
				} else {
					monthDay.Set(toolkit.ToString(xDate.Year())+""+toolkit.ToString(int(xDate.Month())),
						toolkit.M{
							"days":         daysInYear.GetInt(toolkit.ToString(int(xDate.Month()))),
							"totalInMonth": daysInYear.GetInt(toolkit.ToString(int(xDate.Month()))),
						})
				}

				month++
				if month > 12 {
					year = year + 1
					month = 1
					daysInYear = GetDayInYear(year)
				}

				xDate, _ = time.Parse("2006-1-2", toolkit.ToString(year)+"-"+toolkit.ToString(month)+"-"+toolkit.ToString(day))
			} else {
				durationMonths++
				months = append(months, int(tEnd.Month()))
				monthDay.Set(toolkit.ToString(tEnd.Year())+""+toolkit.ToString(int(tEnd.Month())), toolkit.M{
					"days":         int(tEnd.Day()),
					"totalInMonth": daysInYear.GetInt(toolkit.ToString(int(tEnd.Month()))),
				})
				break out
			}
		}
	}

	if durationMonths == 0 {
		months = append(months, int(tEnd.Month()))
		durationMonths = 1
		monthDay.Set(toolkit.ToString(tEnd.Year())+""+toolkit.ToString(int(tEnd.Month())), toolkit.M{
			"days":         int(tEnd.Day()) - (int(tStart.Day()) - 1),
			"totalInMonth": daysInYear.GetInt(toolkit.ToString(int(tEnd.Month()))),
		})
	}

	return durationMonths, months, monthDay
}

// add by RS, 2016-10-26 to assign start date & end date based on period type
func GetStartEndDate(r *knot.WebContext, period string, tStart, tEnd time.Time) (startDate, endDate time.Time, err error) {
	currentDate := time.Now().UTC()
	if period == "custom" {
		if tStart.Year() > 2012 || tEnd.Year() > 2012 {
			startDate, _ = time.Parse("2006-01-02", tStart.UTC().Format("2006-01-02"))
			/*if tEnd.Truncate(24 * time.Hour).Equal(currentDate.Truncate(24 * time.Hour)) {
				endDate = currentDate
			} else {
				endDate, _ = time.Parse("2006-01-02 15:04:05", tEnd.UTC().Format("2006-01-02")+" 23:59:59")
			}*/
			endDate, _ = time.Parse("2006-01-02 15:04:05.00", tEnd.UTC().Format("2006-01-02")+" 23:59:59.99")
		} else {
			err = errors.New("Date Cannot be Less Than 2013")
		}
	} else {
		endDate = currentDate

		switch period {
		case "last24hours":
			startDate = endDate.Add(-24 * time.Hour)
		case "last7days":
			startDate = time.Date(endDate.Year(), endDate.Month(), endDate.Day()-7, 0, 0, 0, 0, time.UTC)
		case "monthly":
			if tStart.Year() > 2012 || tEnd.Year() > 2012 {
				/*start date sudah tanggal 1 dari frontend*/
				startDate, _ = time.Parse("2006-01-02", tStart.UTC().Format("2006-01-02"))
				/*if (tEnd.Year() == currentDate.Year()) && (tEnd.Month() == currentDate.Month()) {
					endDate = currentDate
				} else {
					t := time.Date(tEnd.Year(), tEnd.Month()+1, 1, 0, 0, 0, 0, time.UTC)
					endDate = time.Date(tEnd.Year(), tEnd.Month(), t.Add(-24*time.Hour).Day(), 23, 59, 59, 0, time.UTC)
				}*/
				/*dari end date frontend ditambah 1 bulan trus dikurangi 1 hari untuk dapet max day di bulan tsb*/
				t := time.Date(tEnd.Year(), tEnd.Month()+1, 1, 0, 0, 0, 0, time.UTC)
				endDate = time.Date(tEnd.Year(), tEnd.Month(), t.Add(-24*time.Hour).Day(), 23, 59, 59, 999999999, time.UTC)
			} else {
				err = errors.New("Date Cannot be Less Than 2013")
			}
		case "annual":
			if tStart.Year() > 2012 || tEnd.Year() > 2012 {
				if tEnd.Year() != endDate.Year() {
					endDate = time.Date(tEnd.Year(), 12, 31, 23, 59, 59, 999999999, time.UTC)
				}
				startDate = time.Date(tStart.Year(), 1, 1, 0, 0, 0, 0, time.UTC)
			} else {
				err = errors.New("Date Cannot be Less Than 2013")
			}
		}
	}
	return
}

/*func GetAllTurbineList() (result []toolkit.M, e error) {
	var projects []interface{}
	resProj, e := GetProjectList()

	for _, v := range resProj {
		projects = append(projects, v.Value)
	}

	csr, e := DB().Connection.
		NewQuery().
		From(new(md.TurbineMaster).TableName()).
		Where(dbox.In("Project", projects...)).
		Order("project, turbineid").
		Cursor(nil)

	if e != nil {
		return
	}
	defer csr.Close()
	e = csr.Fetch(&result, 0, false)

	return
}*/

func HelperSetDb(conn dbox.IConnection) {
	_ = SetDb(conn)
}

// { "value": 1, "text": "Ambient Temp", "colname": "tempoutdoor" }

func GetTemperatureList() (result toolkit.M, e error) {
	csr, e := DB().Connection.NewQuery().
		From("ref_databrowsertag").
		Order("projectname", "label").
		Cursor(nil)
	defer csr.Close()

	if e != nil {
		return
	}

	_data := toolkit.M{}
	lastProject := ""
	currProject := ""
	indexCount := 1
	tempList := []toolkit.M{}
	result = toolkit.M{}
	for {
		_data = toolkit.M{}
		e = csr.Fetch(&_data, 1, false)
		if e != nil {
			break
		}
		currProject = _data.GetString("projectname")
		if lastProject != currProject {
			if lastProject != "" {
				result.Set(lastProject, tempList)
				indexCount = 1
				tempList = []toolkit.M{}
			}
			lastProject = currProject
		}
		if strings.Contains(strings.ToLower(_data.GetString("realtimefield")), "temp") {
			tempList = append(tempList, toolkit.M{
				"value":   indexCount,
				"text":    _data.GetString("label"),
				"colname": strings.ToLower(_data.GetString("realtimefield")),
			})
			indexCount++
		}
	}
	if lastProject != "" {
		result.Set(lastProject, tempList)
	}

	return
}

func GetAlarmTagsList() (result toolkit.M, e error) {
	csr, e := DBRealtime().NewQuery().
		From("ref_alarmtaglist").
		Where(dbox.Eq("enable", true)).
		Order("projectname", "tagsdesc").
		Cursor(nil)
	defer csr.Close()

	if e != nil {
		return
	}

	_data := toolkit.M{}
	lastProject := ""
	currProject := ""
	indexCount := 1
	tagList := []toolkit.M{}
	result = toolkit.M{}
	allType := []toolkit.M{
		toolkit.M{
			"value":   0,
			"text":    "All Types",
			"colname": "alltypes",
		},
	}
	for {
		_data = toolkit.M{}
		e = csr.Fetch(&_data, 1, false)
		if e != nil {
			break
		}
		currProject = _data.GetString("projectname")
		if lastProject != currProject {
			if lastProject != "" {
				tagList = append(allType, tagList...)
				result.Set(lastProject, tagList)
				indexCount = 1
				tagList = []toolkit.M{}
			}
			lastProject = currProject
		}
		tagList = append(tagList, toolkit.M{
			"value":   indexCount,
			"text":    _data.GetString("tagsdesc"),
			"colname": _data.GetString("tags"),
		})

		indexCount++
	}
	if lastProject != "" {
		tagList = append(allType, tagList...)
		result.Set(lastProject, tagList)
	}

	return
}

func GetStateList(projlist []md.ProjectOut) (result []string, e error) {
	result = []string{}
	_tkm := toolkit.M{}
	for _, proj := range projlist {
		_tkm.Set(proj.State, 1)
	}

	for state, _ := range _tkm {
		result = append(result, state)
	}

	sort.Strings(result)

	return
}

func GetProjectList() (result []md.ProjectOut, e error) {
	pipes := []toolkit.M{
		toolkit.M{"$match": toolkit.M{"active": true}},
	}
	csr, e := DB().Connection.NewQuery().
		From(new(md.ProjectMaster).TableName()).
		Command("pipe", pipes).
		Cursor(nil)
	defer csr.Close()

	if e != nil {
		return
	}

	data := []md.ProjectMaster{}
	e = csr.Fetch(&data, 0, false)

	for _, val := range data {
		result = append(result, md.ProjectOut{
			ProjectId:         val.ProjectId,
			NoOfTurbine:       val.TotalTurbine,
			TotalMaxCapacity:  val.TotalPower,
			Name:              fmt.Sprintf("%v (%v | %v MW)", val.ProjectId, val.TotalTurbine, val.TotalPower),
			Value:             val.ProjectId,
			Coords:            []float64{val.Latitude, val.Longitude},
			RevenueMultiplier: val.RevenueMultiplier,
			City:              val.City,
			SS_AirDensity:     val.SS_AirDensity,
			STD_AirDensity:    val.STD_AirDensity,
			Engine:            val.Engine,
			State:             val.State,
			ForecastMinCap:    val.Forecast_Min_Cap,
			ForecastMaxCap:    val.Forecast_Max_Cap,
			ForecastRevInfos:  val.Forecast_Revision_Info,
		})
	}

	// sort.Strings(result)
	return
}

func GetTurbineList(projects []interface{}) (result []md.TurbineOut, e error) {
	query := DB().Connection.NewQuery().From("ref_turbine")
	pipes := []toolkit.M{}
	if len(projects) > 0 {
		pipes = []toolkit.M{
			toolkit.M{"$match": toolkit.M{"project": toolkit.M{"$in": projects}}},
		}
	}
	pipes = append(pipes, toolkit.M{"$sort": toolkit.M{"turbinename": 1}})

	query = query.Command("pipe", pipes)

	csr, e := query.
		Cursor(nil)
	defer csr.Close()
	if e != nil {
		return
	}

	data := []md.TurbineMaster{}
	e = csr.Fetch(&data, 0, false)

	for _, val := range data {
		result = append(result, md.TurbineOut{
			Project:    val.Project,
			Turbine:    val.TurbineName,
			Value:      val.TurbineId,
			Capacity:   val.CapacityMW,
			Feeder:     val.Feeder,
			Engine:     val.Engine,
			Coords:     []float64{val.Latitude, val.Longitude},
			Cluster:    val.Cluster,
			DgrProject: val.ProjectDgr,
			DgrTurbine: val.TurbineDgr,
		})
	}

	return
}

func GetTurbineNameList(project string) (turbineName map[string]string, err error) {
	query := DBRealtime().NewQuery().From("ref_turbine")
	if project != "" && project != "Fleet" {
		pipes := []toolkit.M{
			toolkit.M{"$match": toolkit.M{"project": project}},
		}
		query = query.Command("pipe", pipes)
	}
	csrTurbine, err := query.Cursor(nil)
	defer csrTurbine.Close()
	if err != nil {
		return
	}
	turbineList := []toolkit.M{}
	err = csrTurbine.Fetch(&turbineList, 0, false)
	if err != nil {
		return
	}
	turbineName = map[string]string{}
	for _, val := range turbineList {
		// if project != "" {
		// 	turbineName[val.GetString("turbineid")] = val.GetString("turbinename")
		// } else {
		// 	turbineName[toolkit.Sprintf("%s_%s", val.GetString("project"), val.GetString("turbineid"))] = val.GetString("turbinename")
		// }
		turbineName[val.GetString("turbineid")] = val.GetString("turbinename")
	}
	return
}

func GetProjectTurbineList(projects []interface{}) (result map[string]toolkit.M, sortedKey []string, e error) {
	var filter []*dbox.Filter
	result = map[string]toolkit.M{}
	sortedKey = []string{}

	if len(projects) > 0 {
		filter = append(filter, dbox.In("project", projects...))
	}

	csr, e := DB().Connection.
		NewQuery().
		From(new(md.TurbineMaster).TableName()).
		Where(filter...).
		Cursor(nil)
	defer csr.Close()

	if e != nil {
		return
	}

	data := []md.TurbineMaster{}
	e = csr.Fetch(&data, 0, false)

	keys := []string{}

	for _, val := range data {
		list := toolkit.M{}

		if result[val.Project] != nil {
			list = result[val.Project]
		} else {
			keys = append(keys, val.Project)
		}
		list.Set(val.TurbineId, val)
		result[val.Project] = list
	}

	sort.Strings(keys)

	return
}

func GetHourValue(tStart time.Time, tEnd time.Time, minDate time.Time, maxDate time.Time) (hourValue float64) {
	startStr := tStart.Format("0601")
	endStr := tEnd.Format("0601")

	minDateStr := minDate.Format("0601")
	maxDateStr := maxDate.Format("0601")

	if startStr == minDateStr {
		minDate = tStart
	} else {
		minDate, _ = time.Parse("060102", minDateStr+"01")
	}

	if endStr != maxDateStr {
		daysInMonth := GetDayInYear(maxDate.Year())
		maxDate, _ = time.Parse("060102", maxDateStr+toolkit.ToString(daysInMonth.GetInt(toolkit.ToString(int(maxDate.Month())))))
	}

	start, _ := time.Parse("060102150405", minDate.Format("060102")+"000000")
	end, _ := time.Parse("060102150405", maxDate.Format("060102")+"235959")

	// log.Printf("hours: %v | %v | %v  \n", end.Sub(start).Hours(), start.String(), end.String())

	hourValue = toolkit.ToFloat64(end.Sub(start).Hours(), 0, toolkit.RoundingUp)

	return
}

// totalTurbine in float64
// okTime sum ok time
// energy should be div by 1000
// machineDownTime, gridDownTime already in hour value
// minutes should be div by 60
func GetAvailAndPLF(totalTurbine float64, okTime float64, energy float64, machineDownTime float64, gridDownTime float64, countTimeStamp float64, hourValue float64, totalMinutes float64, plfDivider float64) (machineAvail float64, gridAvail float64, dataAvail float64, totalAvail float64, plf float64) {
	divider := (totalTurbine * hourValue)
	plf = energy / (plfDivider * hourValue) * 100
	// log.Printf(">>> %v >>> %v | %v | %v \n", plf, energy, plfDivider, hourValue)
	totalAvail = (okTime / 3600) / divider * 100
	machineAvail = (totalMinutes - machineDownTime) / divider * 100
	gridAvail = (totalMinutes - gridDownTime) / divider * 100
	dataAvail = (countTimeStamp * 10 / 60) / divider * 100
	return
}

//=============================
//Revision Of GetAvailAndPLF @asp:20170725
//=============================
//	Input : noofturbine, oktime, energy, counttimestamp, totalhour, totalcapacity,
//			machinedowntime, griddowntime, otherdowntime
//
//	counttimestamp -> count total data that available (*every data in 10 mins conv)
//	oktime, totalhour -> in hour
//	machinedowntime, griddowntime, otherdowntime -> in hour
//	totalcapacity -> in MWatt
//  energy -> in MWh
//
//	Output : totalavailability, plf, machineavailability, gridavailability, dataavailability
//=============================
func CalcAvailabilityAndPLF(in toolkit.M) (res toolkit.M) {
	res = toolkit.M{}
	totalhour := in.GetFloat64("totalhour")
	divider := in.GetFloat64("noofturbine") * totalhour

	plf := toolkit.Div(in.GetFloat64("energy"), (totalhour*in.GetFloat64("totalcapacity"))) * 100
	if plf <= 0 {
		plf = 0
	}
	res.Set("plf", plf)

	totalavailability := toolkit.Div(in.GetFloat64("oktime"), divider) * 100
	res.Set("totalavailability", totalavailability)

	gdown, mdown, odown := in.GetFloat64("griddowntime"), in.GetFloat64("machinedowntime"), in.GetFloat64("otherdowntime")
	mdowndivider := divider - gdown - odown

	machineavailability := toolkit.Div(mdowndivider-mdown, mdowndivider) * 100
	res.Set("machineavailability", machineavailability)

	gridavailability := toolkit.Div(divider-gdown, divider) * 100
	res.Set("gridavailability", gridavailability)

	dataavailability := toolkit.Div(in.GetFloat64("counttimestamp")/6, divider) * 100
	res.Set("dataavailability", dataavailability)

	return
}

func GetDataDateAvailable(collectionName string, timestampColumn string, where *dbox.Filter) (min time.Time, max time.Time, err error) {
	min, max, err = hp.GetDataDateAvailable(collectionName, timestampColumn, where, DB().Connection)
	return
}

func GetHFDFolder() string {
	config := hp.ReadConfig()
	source := config["hfdfolder"]
	return source + string(os.PathSeparator)
}
