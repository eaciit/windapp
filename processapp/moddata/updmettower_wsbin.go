package main

import (
	tk "github.com/eaciit/toolkit"

	"bufio"
	"os"
	"strings"
	"time"

	"github.com/eaciit/dbox"
	_ "github.com/eaciit/dbox/dbc/mongo"

	. "eaciit/wfdemo-git/library/models"
)

var (
	stablename = "MetTower"
	dtablename = "MetTower"
)

func main() {
	t0 := time.Now()
	// from 1 nov - 30 nov
	_conn, _ := PrepareConnection()
	defer _conn.Close()

	stablename = new(MetTower).TableName()
	dtablename = tk.Sprintf("%s-wsbin", new(MetTower).TableName())
	// stime := time.Date(2016, 11, 1, 0, 0, 0, 0, time.UTC)

	csr, _ := _conn.NewQuery().Select().From(stablename).
		// Where(dbox.And(dbox.Gte("timestamp", stime), dbox.Lt("timestamp", stime.AddDate(0, 1, 0)))).
		// Order("timestamp").
		Cursor(nil)
	defer csr.Close()

	count := csr.Count()
	sresult := make(chan int, count)
	sdata := make(chan *MetTower, count)
	for i := 0; i < 5; i++ {
		go workersave(i, sdata, sresult)
	}

	// tk.Println(stime)
	step := getstep(count)
	_i := 0
	for {
		_sd := new(MetTower)
		err := csr.Fetch(_sd, 1, false)
		if err != nil {
			break
		}

		// i++
		// tk.Println(">>> ", _sd.TimeStamp.UTC())
		if _i%step == 0 {
			tk.Printfn("Process Data %d to %d, in %s",
				_i, count, time.Since(t0).String())
		}
		sdata <- _sd
	}

	close(sdata)

	for i := 0; i < count; i++ {
		<-sresult
		if i%step == 0 {
			tk.Printfn("Done Saved Data %d to %d, in %s",
				i, count, time.Since(t0).String())
		}
	}
	close(sresult)

	// tk.Printfn("Done All Process Data in %s", time.Since(t0).String())
}

func workersave(wi int, jobs <-chan *MetTower, result chan<- int) {
	workerconn, _ := PrepareConnection()
	defer workerconn.Close()

	qSave := workerconn.NewQuery().
		From(dtablename).
		SetConfig("multiexec", true).
		Save()

	trx := new(MetTower)
	for trx = range jobs {

		trx.WindSpeedBin = tk.RoundingAuto64(trx.VHubWS90mAvg, 0)

		err := qSave.Exec(tk.M{}.Set("data", trx))
		if err != nil {
			tk.Println(err)
		}

		result <- 1
	}

	return
}

func PrepareConnection() (dbox.IConnection, error) {
	// config := ReadConfig()

	// ci := &dbox.ConnectionInfo{config["host"], config["database"], config["username"], config["password"], tk.M{}.Set("timeout", 3000)}
	ci := &dbox.ConnectionInfo{"localhost:27017", "wfdemo", "admin", "qwerty", tk.M{}.Set("timeout", 3000)}
	tk.Println("Connect : ", ci)
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
	file, err := os.Open("../conf" + "/" + "app.conf")
	if err == nil {
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

	file.Close()
	return ret
}

func getstep(count int) int {
	v := count / 5
	if v == 0 {
		return 1
	}
	return v
}
