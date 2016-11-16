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

	. "eaciit/wfdemo/library/models"

	"time"

	_ "github.com/eaciit/dbox/dbc/mongo"
	"github.com/eaciit/orm"

	dc "eaciit/wfdemo/processapp/threeextractor/dataconversion"
)

var (
	wd = func() string {
		d, _ := os.Getwd()
		return d + separator
	}()
	separator = string(os.PathSeparator)
)

func main() {
	conn, err := PrepareConnection()
	if err != nil {
		tk.Println(err)
	}
	ctx := orm.New(conn)

	start := time.Now()
	log.Println("Start Convert...")

	UpdateThreeSecs("", conn, ctx)

	convExt := dc.NewConvThreeExt(ctx)
	convExt.Generate("")

	conv10 := dc.NewDataConversion(ctx)
	conv10.Generate("")

	duration := time.Now().Sub(start).Seconds()
	log.Printf("End in: %v sec(s) \n", duration)
}

func UpdateThreeSecs(file string, conn dbox.IConnection, ctx *orm.DataContext) {
	/*file := ""

	pipes := []tk.M{}

	match := tk.M{}
	if file != "" {
		match = tk.M{"file": file}
		pipes = append(pipes, tk.M{"$match": match})
	}*/

	csr, e := ctx.Connection.NewQuery().
		From(new(ScadaThreeSecs).TableName()).
		// Command("pipe", pipes).
		Cursor(nil)

	if e != nil {
		log.Printf("ERRROR: %v \n", e.Error())
	}

	defer csr.Close()

	list := []ScadaThreeSecs{}
	e = csr.Fetch(&list, 0, false)

	if e != nil {
		log.Printf("ERRROR: %v \n", e.Error())
	}

	updateList := []ScadaThreeSecs{}

	for _, val := range list {
		timeStamp := val.TimeStamp1.UTC()
		seconds := tk.Div(tk.ToFloat64(timeStamp.Nanosecond(), 1, tk.RoundingAuto), 1000000000)
		secondsInt := tk.ToInt(seconds, tk.RoundingAuto)
		newTimeTmp := timeStamp.Add(time.Duration(secondsInt) * time.Second)
		strTime := tk.ToString(newTimeTmp.Year()) + tk.ToString(int(newTimeTmp.Month())) + tk.ToString(newTimeTmp.Day()) + " " + tk.ToString(newTimeTmp.Hour()) + ":" + tk.ToString(newTimeTmp.Minute()) + ":" + tk.ToString(newTimeTmp.Second())

		TimeStampSecondGroup, _ := time.Parse("200612 15:4:5", strTime)

		THour := TimeStampSecondGroup.Hour()
		TMinute := TimeStampSecondGroup.Minute()
		TSecond := TimeStampSecondGroup.Second()
		TMinuteValue := float64(TMinute) + tk.Div(float64(TSecond), 60.0)
		TMinuteCategory := tk.ToInt(tk.RoundingUp64(tk.Div(TMinuteValue, 10), 0)*10, "0")

		newTimeStamp := val.DateId1.Add(time.Duration(THour) * time.Hour).Add(time.Duration(TMinuteCategory) * time.Minute)

		TimeStampConverted := newTimeStamp.UTC()
		TimeStampConvertedInt, _ := strconv.ParseInt(TimeStampConverted.Format("200601021504"), 10, 64)

		/*updateData := tk.M{}
		updateData.Set("timestampsecondgroup", TimeStampSecondGroup)
		updateData.Set("timestampconverted", TimeStampConverted)
		updateData.Set("timestampconvertedint", TimeStampConvertedInt)

		updateData.Set("thour", THour)
		updateData.Set("tminute", TMinute)
		updateData.Set("tsecond", TSecond)
		updateData.Set("tminutevalue", TMinuteValue)
		updateData.Set("tminutecategory", TMinuteCategory)

		updateList[val.ID] = updateData*/

		val.THour = THour
		val.TMinute = TMinute
		val.TSecond = TSecond
		val.TMinuteValue = TMinuteValue
		val.TMinuteCategory = TMinuteCategory
		val.TimeStampConverted = TimeStampConverted
		val.TimeStampConvertedInt = TimeStampConvertedInt
		val.TimeStampSecondGroup = TimeStampSecondGroup

		updateList = append(updateList, val)

		e := ctx.Save(&val)
		if e != nil {
			log.Printf("ERRROR: %v \n", e.Error())
		}
	}

	for _, val := range updateList {
		e := ctx.Save(&val)
		if e != nil {
			log.Printf("ERRROR: %v \n", e.Error())
		}
	}

	/*if len(updateList) > 0 {
		for id, updateData := range updateList {
			e := conn.NewQuery().Update().From(new(ScadaThreeSecs).TableName()).Where(dbox.Eq("_id", id)).Exec(tk.M{}.Set("data", updateData))

			if e != nil {
				log.Printf("ERRROR: %v \n", e.Error())
			}
		}
	}*/
}

func PrepareConnection() (dbox.IConnection, error) {
	config := ReadConfig()

	// log.Printf("config: %#v \n", config)

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
