package helper

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"strconv"
	"strings"
	"time"

	tk "github.com/eaciit/toolkit"
	"github.com/tealeg/xlsx"
	// _ "github.com/tealeg/xlsx"

	"github.com/eaciit/dbox"
	_ "github.com/eaciit/dbox/dbc/mongo"
)

var (
	wd = func() string {
		d, _ := os.Getwd()
		return d + "/"
	}()

	DateFormat1 = "02-01-2006 15:04:05"
	DateFormat2 = "02-01-2006 04:05."
	DateFormat3 = "02-01-06 15:04:05"
	DateFormat4 = "01-02-06 15:04:05"
	DateFormat5 = "2-1-2006 15:4:5"
)

type DateInfo struct {
	DateId    time.Time
	MonthId   int
	MonthDesc string
	QtrId     int
	QtrDesc   string
	WeekId    int    `bson:"weekid,omitempty" json:"weekid,omitempty"`
	WeekDesc  string `bson:"weekdesc,omitempty" json:"weekdesc,omitempty"`
	Year      int
}

type SortDirection struct {
	Field string
	Dir   string
}

func GetDateInfo(t time.Time) DateInfo {
	di := DateInfo{}

	year := t.Year()
	month := int(t.Month())

	monthid := strconv.Itoa(year) + LeftPad2Len(strconv.Itoa(month), "0", 2)
	monthdesc := t.Month().String() + " " + strconv.Itoa(year)

	qtr := 0
	if month%3 > 0 {
		qtr = int(math.Ceil(float64(month / 3)))
		qtr = qtr + 1
	} else {
		qtr = month / 3
	}

	qtrid := strconv.Itoa(year) + LeftPad2Len(strconv.Itoa(qtr), "0", 2)
	qtrdesc := "Q" + strconv.Itoa(qtr) + " " + strconv.Itoa(year)

	di.DateId, _ = time.Parse("2006-01-02 15:04:05", t.UTC().Format("2006-01-02")+" 00:00:00")
	di.Year = year
	di.MonthDesc = monthdesc
	di.MonthId, _ = strconv.Atoi(monthid)
	di.QtrDesc = qtrdesc
	di.QtrId, _ = strconv.Atoi(qtrid)

	_, week := di.DateId.ISOWeek()
	weekid := strconv.Itoa(year) + LeftPad2Len(strconv.Itoa(week), "0", 2)
	di.WeekId, _ = strconv.Atoi(weekid)
	di.WeekDesc = tk.Sprintf("W %v", weekid)

	return di
}

func MonthIDToDateInfo(mid int) (dateInfo DateInfo) {
	monthid := strconv.Itoa(mid)
	year := monthid[0:4]
	month := monthid[4:6]
	day := "01"

	iMonth, _ := strconv.Atoi(string(month))
	iMonth = iMonth - 1

	dtStr := year + "-" + month + "-" + day
	date, _ := time.Parse("2006-01-02", dtStr)

	dateInfo = GetDateInfo(date)

	return
}

func LeftPad2Len(s string, padStr string, overallLen int) string {
	var padCountInt int
	padCountInt = 1 + ((overallLen - len(padStr)) / len(padStr))
	var retStr = strings.Repeat(padStr, padCountInt) + s
	return retStr[(len(retStr) - overallLen):]
}

func ErrorHandler(e error, position string) {
	if e != nil {
		tk.Printf("ERROR on %v: %v \n", position, e.Error())
	}
}

func ErrorLog(e error, position string, errorList []error) []error {
	if e != nil {
		errorList = append(errorList, e)
		// tk.Printf("ERROR on %v: %v \n", position, e.Error())
	}
	return errorList
}

func GetFloatCell(cell *xlsx.Cell) (result float64, e error) {
	str, e := cell.String()
	result = 0

	if str != "" {
		result, e = cell.Float()
	}

	return
}

func GetDateCell(strDate string) (result time.Time, e error) {
	result, e = time.Parse(DateFormat1, strDate)
	if e != nil {
		// tk.Printf("DateFormat1 ERROR: %v \n", strDate)
		result, e = time.Parse(DateFormat2, strDate)
		if e != nil {
			// tk.Printf("DateFormat2 ERROR: %v \n", strDate)
			result, e = time.Parse(DateFormat3, strDate)
			if e != nil {
				// tk.Printf("DateFormat3 ERROR: %v \n", strDate)
				result, e = time.Parse(DateFormat4, strDate)
				if e != nil {
					tk.Printf("GetDateCell ERROR: %v \n", strDate)
				}
			}
		}
	}

	return
}

func GetDateCellAuto(cellDate *xlsx.Cell, cellTime *xlsx.Cell) (result time.Time, e error) {
	strDate := ""
	strTime := ""

	if cellDate != nil && cellTime != nil {
		var tmp float64

		tmp, e = cellTime.Float()
		cellTime.SetDateTimeWithFormat(tmp, "15:04:05")
		strTime, _ = cellTime.FormattedValue()

		if strTime == "" {
			e = errors.New("Date or Time is not Valid")
			return
		}

		tmp, e = cellDate.Float()
		if e != nil {
			tmpStr := ""
			tmpStr, e = cellDate.String()
			if tmpStr == "" {
				e = errors.New("Date or Time is not Valid")
				return
			}

			result, e = GetDateCell(tmpStr + " " + strTime)
			return
		}

		cellDate.SetDateTimeWithFormat(tmp, time.UnixDate)
		strDate, e = cellDate.FormattedValue()
		if e != nil {
			return
		}

		result, e = time.Parse(time.UnixDate, strings.Replace(strDate, "00:00:00", strTime, 1))
	} else {
		e = errors.New("Please Input Date and Time")
	}
	return
}

func ReverseMonthDate(date time.Time) (result time.Time, e error) {
	year, month, day := date.Date()
	hour, minute, second := date.Clock()
	dtStr := tk.ToString(month) + "-" + tk.ToString(day) + "-" + tk.ToString(year) + " " + tk.ToString(hour) + ":" + tk.ToString(minute) + ":" + tk.ToString(second)
	if day <= 12 {
		result, e = time.Parse(DateFormat5, dtStr)
	} else {
		e = errors.New("Date is not valid")
		result = date
	}

	return
}

func WriteErrors(errorList tk.M, fileName string) (e error) {
	config := ReadConfig()
	source := config["datasource"]
	dataSourceFolder := "errors"
	fileName = fileName + "_" + tk.GenerateRandomString("", 5) + ".txt"
	tk.Printf("Saving Errors... %v\n", fileName)

	errors := ""

	for x, err := range errorList {
		errors = errors + "" + fmt.Sprintf("#%v: %#v \n", x, err)
	}

	e = ioutil.WriteFile(source+"\\"+dataSourceFolder+"\\"+fileName, []byte(errors), 0644)
	return
}

func ReadConfig() map[string]string {
	ret := make(map[string]string)
	file, err := os.Open(wd + "config/app.conf")
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

func GetConnRealtime() dbox.IConnection {
	config := ReadConfig()

	ci := &dbox.ConnectionInfo{config["host"], config["dbrealtime"], config["username"], config["password"], tk.M{}.Set("timeout", 3000)}
	c, _ := dbox.NewConnection("mongo", ci)

	for {
		e := c.Connect()
		if e != nil {
			tk.Println("Realtime DB Connection Found ", e.Error())
		} else {
			break
		}
	}

	return c
}

func GetDateRange(dt time.Time, isBefore bool) (result time.Time) {
	_, minute, _ := dt.Clock()

	if isBefore {
		if minute > 10 && minute < 20 {
			minute = minute - 10
		} else if minute > 20 && minute < 30 {
			minute = minute - 20
		} else if minute > 30 && minute < 40 {
			minute = minute - 30
		} else if minute > 40 && minute < 50 {
			minute = minute - 40
		} else if minute > 50 && minute < 60 {
			minute = minute - 50
		}
		switch minute {
		case 1:
			result = dt.Add(-1 * time.Minute)
			break
		case 2:
			result = dt.Add(-2 * time.Minute)
			break
		case 3:
			result = dt.Add(-3 * time.Minute)
			break
		case 4:
			result = dt.Add(-4 * time.Minute)
			break
		case 5:
			result = dt.Add(-5 * time.Minute)
			break
		case 6:
			result = dt.Add(-6 * time.Minute)
			break
		case 7:
			result = dt.Add(-7 * time.Minute)
			break
		case 8:
			result = dt.Add(-8 * time.Minute)
			break
		case 9:
			result = dt.Add(-9 * time.Minute)
			break
		default:
			result = dt.Add(-10 * time.Minute)
			break
		}
	} else {
		if minute > 50 && minute < 60 {
			minute = 60 - minute
		}
		switch minute {
		case 1:
			result = dt.Add(1 * time.Minute)
			break
		case 2:
			result = dt.Add(2 * time.Minute)
			break
		case 3:
			result = dt.Add(3 * time.Minute)
			break
		case 4:
			result = dt.Add(4 * time.Minute)
			break
		case 5:
			result = dt.Add(5 * time.Minute)
			break
		case 6:
			result = dt.Add(6 * time.Minute)
			break
		case 7:
			result = dt.Add(7 * time.Minute)
			break
		case 8:
			result = dt.Add(8 * time.Minute)
			break
		case 9:
			result = dt.Add(9 * time.Minute)
			break
		default:
			result = dt.Add(10 * time.Minute)
			break
		}
	}

	resultStr := result.Format("2006-01-02 15:04")
	result, e := time.Parse("2006-01-02 15:04:05", resultStr[0:len(resultStr)-1]+"0"+":00")
	// tk.Printf("%v | %v | %v \n", dt.Format("2006-01-02 15:04"), resultStr, result.Format("2006-01-02 15:04:05"))
	ErrorHandler(e, "GetDateRange")
	return
}

func Round(f float64) float64 {
	return math.Floor(f + .5)
}

func RoundPlus(f float64, places int) float64 {
	shift := math.Pow(10, float64(places))
	return Round(f*shift) / shift
}

func RoundUp(val float64, roundOn float64, places int) (newVal float64) {
	var round float64
	pow := math.Pow(10, float64(places))
	digit := pow * val
	_, div := math.Modf(digit)
	if div >= roundOn {
		round = math.Ceil(digit)
	} else {
		round = math.Floor(digit)
	}
	newVal = round / pow
	return
}

func GetDayInYear(year int) tk.M {
	result := tk.M{}
	for m := time.January; m <= time.December; m++ {
		t := time.Date(year, m+1, 1, 0, 0, 0, 0, time.UTC)
		result.Set(tk.ToString(int(m)), t.Add(-24*time.Hour).Day())
	}
	return result
}

func ReadJson(source string, result interface{}) {
	file, err := os.Open(wd + source)
	if err == nil {
		defer file.Close()

		jsonParser := json.NewDecoder(file)
		err = jsonParser.Decode(&result)

		if err != nil {
			tk.Println(err.Error())
		}
	} else {
		tk.Println(err.Error())
	}
}

func GetDaysNoByQuarter(year int, qtr int, lastDate time.Time) int {
	totalDays := 0
	lastMonth := qtr * 3
	for i := 1; i <= lastMonth; i++ {
		date, _ := time.Parse("2006-01-02", tk.Sprintf("%v-%v-%v", year, i, 1))
		dateMonth := time.Date(date.Year(), date.Month(), 0, 0, 0, 0, 0, time.UTC)
		if dateMonth.Month() != lastDate.Month() && dateMonth.Year() != lastDate.Year() {
			totalDays += dateMonth.Day()
		} else {
			totalDays += lastDate.Day()
		}
	}
	// tk.Println(totalDays, lastMonth, year, qtr)

	return totalDays
}

// periodType:
//		YEAR => in Years
//		QTR => in Quarters
//		MONTH => in Months
//		WEEK => in Weeks
//		DAY	=> in Days
func GetPeriodBackByDate(periodType string, lastDate time.Time, noPeriodBack int) time.Time {
	var ret time.Time

	exactNoPeriodBack := noPeriodBack - 1
	dateLayout := "2006-01-02"
	lastMonth := int(lastDate.Month())
	lastYear := lastDate.Year()
	switch periodType {
	case "YEAR":
		startYear := lastYear - exactNoPeriodBack
		ret, _ = time.Parse(dateLayout, tk.Sprintf("%v-%v-%v", startYear, 1, 1))
	case "QTR":
		lastQtr := 0
		if lastMonth%3 > 0 {
			lastQtr = int(math.Ceil(float64(lastMonth / 3)))
			lastQtr = lastQtr + 1
		} else {
			lastQtr = lastMonth / 3
		}
		startQtr := lastQtr
		startYear := lastYear
		noQtr := exactNoPeriodBack
		for noQtr > 0 {
			startQtr--
			if startQtr == 0 {
				startQtr = 4
				startYear--
			}
			noQtr--
		}

		startMonthOfStartQtr := (startQtr * 3) - 2
		ret, _ = time.Parse(dateLayout, tk.Sprintf("%v-%v-%v", startYear, LeftPad2Len(tk.ToString(startMonthOfStartQtr), "0", 2), "01"))
	case "MONTH":
		startMonth := lastMonth
		startYear := lastYear
		noMonths := exactNoPeriodBack
		for noMonths > 0 {
			startMonth--
			if startMonth == 0 {
				startMonth = 12
				startYear--
			}
			noMonths--
		}
		ret, _ = time.Parse(dateLayout, tk.Sprintf("%v-%v-%v", startYear, LeftPad2Len(tk.ToString(startMonth), "0", 2), "01"))
	case "WEEK":
		lastYear, lastWeek := lastDate.ISOWeek()
		startWeek := lastWeek
		startYear := lastYear
		for i := lastWeek; i > 0; i-- {
			startWeek--
			if startWeek == 0 {
				startWeek = 52
				startYear--
			}
		}
		ret = FirstDayOfISOWeek(startYear, startWeek, time.UTC)
	case "DAY":
		ret = lastDate.AddDate(0, 0, -1*exactNoPeriodBack)
	}

	return ret
}

func FirstDayOfISOWeek(year int, week int, timezone *time.Location) time.Time {
	date := time.Date(year, 0, 0, 0, 0, 0, 0, timezone)
	isoYear, isoWeek := date.ISOWeek()
	for date.Weekday() != time.Monday { // iterate back to Monday
		date = date.AddDate(0, 0, -1)
		isoYear, isoWeek = date.ISOWeek()
	}
	for isoYear < year { // iterate forward to the first day of the first week
		date = date.AddDate(0, 0, 1)
		isoYear, isoWeek = date.ISOWeek()
	}
	for isoWeek < week { // iterate forward to the first day of the given week
		date = date.AddDate(0, 0, 1)
		isoYear, isoWeek = date.ISOWeek()
	}
	return date
}

func GetDataDateAvailable(collectionName string, timestampColumn string, where *dbox.Filter, ctx dbox.IConnection) (min time.Time, max time.Time, err error) {
	q := ctx.
		NewQuery().
		From(collectionName)

	if where != nil {
		q.Where(where)
	}

	csr, err := q.
		Aggr(dbox.AggrMin, "$"+timestampColumn, "min").
		Aggr(dbox.AggrMax, "$"+timestampColumn, "max").
		Group("enable").
		Cursor(nil)

	defer csr.Close()

	if err != nil {
		csr.Close()
		return
	}

	data := []tk.M{}
	err = csr.Fetch(&data, 0, false)

	if err != nil || len(data) == 0 {
		csr.Close()
		return
	}

	min = data[0].Get("min").(time.Time)
	max = data[0].Get("max").(time.Time)

	csr.Close()
	return
}

func GetNormalAddDateMonth(dt time.Time, month int) (res time.Time) {
	tmp, _ := time.Parse("060102_150405", dt.Format("0601")+"01_"+dt.Format("150405"))
	res = tmp.AddDate(0, month, 0)

	return
}

func UpperFirstLetter(str string) string {
	if len(str) > 0 {
		str = strings.ToUpper(str[:1]) + str[1:]
	}

	return str
}
