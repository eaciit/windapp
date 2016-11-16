
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
	file, err := os.Open(wd + "conf/app.conf")
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
