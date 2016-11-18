package main

import (
	"bufio"
	"log"
	_ "math"
	"os"
	"strconv"
	"strings"

	"github.com/eaciit/dbox"
	tk "github.com/eaciit/toolkit"

	. "github.com/eaciit/windapp/library/models"

	"time"

	"flag"
	_ "github.com/eaciit/dbox/dbc/mongo"
	"github.com/eaciit/orm"
	// dc "github.com/eaciit/windapp/processapp/threeextractor/dataconversion"
)

var (
	wd = func() string {
		d, _ := os.Getwd()
		return d + separator
	}()
	separator = string(os.PathSeparator)

	intstartdate = int(20160821)
	intenddate   = int(20160821)

	startdate, enddate time.Time
)

func main() {
	flag.IntVar(&intstartdate, "sdate", 20160821, "Start date for processing data")
	flag.IntVar(&intenddate, "edate", 20160821, "End date for processing data")
	flag.Parse()

	startdate = tk.String2Date(tk.Sprintf("%d", intstartdate), "YYYYMMdd").UTC()
	enddate = tk.String2Date(tk.Sprintf("%d", intenddate), "YYYYMMdd").UTC().AddDate(0, 0, 1)

	log.Println(tk.Sprintf("Convert Data from %v to %v", startdate, enddate))

	conn, err := PrepareConnection()
	if err != nil {
		tk.Println(err)
	}
	ctx := orm.New(conn)

	start := time.Now()
	log.Println("Start Convert...")

	UpdateThreeSecs("", conn, ctx)

	tk.Printfn("Update three second done in %s",
		time.Since(start).String())

	duration := time.Now().Sub(start).Seconds()
	log.Printf("End in: %v sec(s) \n", duration)
}

func UpdateThreeSecs(file string, conn dbox.IConnection, ctx *orm.DataContext) {

	t0 := time.Now()

	csr, e := ctx.Connection.NewQuery().
		From(new(ScadaThreeSecs).TableName()).
		Where(dbox.Gt("timestamp1", startdate.UTC()), dbox.Lte("timestamp1", enddate.UTC())).
		Cursor(nil)

	if e != nil {
		log.Printf("ERRROR: %v \n", e.Error())
		return
	}

	defer csr.Close()

	scount := csr.Count()
	tk.Printfn("Prepare Query in %s, count %d rows", time.Since(t0).String(), scount)

	iscount := 0
	step := getstep(scount)

	sresult := make(chan int, scount)
	sdata := make(chan ScadaThreeSecs, scount)
	for i := 0; i < 10; i++ {
		go workersave(i, sdata, sresult)
	}

	for {
		val := ScadaThreeSecs{}
		e = csr.Fetch(&val, 1, false)
		if e != nil {
			if strings.Contains(e.Error(), "Not found") {
				log.Printf("EOF")
			} else {
				log.Printf("ERRROR: %v \n", e.Error())
			}
			break
		}

		sdata <- val
		iscount++

		if iscount%step == 0 {
			tk.Printfn("Sending %d of %d (%d) in %s", iscount, scount, iscount*100/scount,
				time.Since(t0).String())
		}
	}

	close(sdata)

	for ri := 0; ri < scount; ri++ {
		<-sresult
		if ri%step == 0 {
			tk.Printfn("Updated %d of %d (%d pct) in %s",
				ri, scount, ri*100/scount, time.Since(t0).String())
		}
	}

	return
}

func getstep(count int) int {
	v := count / 10
	if v == 0 {
		return 1
	}
	return v
}

func workersave(wi int, jobs <-chan ScadaThreeSecs, result chan<- int) {
	workerconn, _ := PrepareConnection()
	defer workerconn.Close()

	dtablename := tk.Sprintf("%s", new(ScadaThreeSecs).TableName())

	qSave := workerconn.NewQuery().
		From(dtablename).
		SetConfig("multiexec", true).
		Save()

	trx := ScadaThreeSecs{}
	for trx = range jobs {

		// ==== From Prev Function
		timeStamp := trx.TimeStamp1.UTC()
		seconds := tk.Div(tk.ToFloat64(timeStamp.Nanosecond(), 1, tk.RoundingAuto), 1000000000)
		secondsInt := tk.ToInt(seconds, tk.RoundingAuto)
		newTimeTmp := timeStamp.Add(time.Duration(secondsInt) * time.Second)
		strTime := tk.ToString(newTimeTmp.Year()) + "-" + tk.ToString(int(newTimeTmp.Month())) + "-" + tk.ToString(newTimeTmp.Day()) + " " + tk.ToString(newTimeTmp.Hour()) + ":" + tk.ToString(newTimeTmp.Minute()) + ":" + tk.ToString(newTimeTmp.Second())

		// tk.Printfn("||| STR >> %#s", strTime)
		TimeStampSecondGroup, _ := time.Parse("2006-1-2 15:4:5", strTime)
		// tk.Println(">>>>>>>>>> ", TimeStampSecondGroup, "--", e)

		THour := TimeStampSecondGroup.Hour()
		TMinute := TimeStampSecondGroup.Minute()
		TSecond := TimeStampSecondGroup.Second()
		TMinuteValue := float64(TMinute) + tk.Div(float64(TSecond), 60.0)
		TMinuteCategory := tk.ToInt(tk.RoundingUp64(tk.Div(TMinuteValue, 10), 0)*10, "0")

		newTimeStamp := trx.DateId1.Add(time.Duration(THour) * time.Hour).Add(time.Duration(TMinuteCategory) * time.Minute)

		TimeStampConverted := newTimeStamp.UTC()
		TimeStampConvertedInt, _ := strconv.ParseInt(TimeStampConverted.Format("200601021504"), 10, 64)

		trx.THour = THour
		trx.TMinute = TMinute
		trx.TSecond = TSecond
		trx.TMinuteValue = TMinuteValue
		trx.TMinuteCategory = TMinuteCategory
		trx.TimeStampConverted = TimeStampConverted
		trx.TimeStampConvertedInt = TimeStampConvertedInt
		trx.TimeStampSecondGroup = TimeStampSecondGroup
		// ==== ==================

		err := qSave.Exec(tk.M{}.Set("data", trx))
		if err != nil {
			tk.Println(err)
		}

		result <- 1
	}
}

func PrepareConnection() (dbox.IConnection, error) {
	config := ReadConfig()

	ci := &dbox.ConnectionInfo{config["host"], config["database"], config["username"], config["password"], tk.M{}.Set("timeout", 3000)}
	c, e := dbox.NewConnection("mongo", ci)

	if e != nil {
		return nil, e
	}

	e = c.Connect()
	if e != nil {
		return nil, e
	}

	return c, nil
}

func ReadConfig() map[string]string {
	ret := make(map[string]string)
	file, err := os.Open("../conf" + separator + "app.conf")
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
