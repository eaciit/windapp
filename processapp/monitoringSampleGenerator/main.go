package main

import (
	"bufio"
	"log"
	"math/rand"
	"os"
	"runtime"
	"strings"

	"time"

	. "eaciit/wfdemo-git/library/models"

	"eaciit/wfdemo-git/library/helper"

	"github.com/eaciit/dbox"
	_ "github.com/eaciit/dbox/dbc/mongo"
	"github.com/eaciit/orm"
	tk "github.com/eaciit/toolkit"
)

var (
	wd = func() string {
		d, _ := os.Getwd()
		return d + separator
	}()
	separator = string(os.PathSeparator)
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	log.Println("Generate Monitoring...")
	conn, err := PrepareConnection()
	if err != nil {
		log.Println(err)
	}
	ctx := orm.New(conn)
	// config := ReadConfig()

	project := map[string]int{}

	project["Tejuva"] = 24
	project["Beluguppa"] = 30
	project["Bercha"] = 40
	project["Bhesada"] = 20
	project["Dalot"] = 50

	tejuvaProject := []string{"HBR004", "HBR005", "HBR006", "HBR007", "SSE001", "SSE002", "SSE006", "SSE007", "SSE011", "SSE012", "SSE015", "SSE017", "SSE018", "SSE019", "SSE020", "TJ013", "TJ016", "TJ021", "TJ022", "TJ023", "TJ024", "TJ025", "HBR038", "TJW024"}

	timestamp := time.Now()

	for k, v := range project {
		for i := 0; i < v; i++ {
			m := new(Monitoring)
			m.TimeStamp = timestamp
			m.DateInfo = helper.GetDateInfo(timestamp)
			m.LastUpdate = m.TimeStamp
			m.LastUpdateDateInfo = m.DateInfo
			m.Project = k

			m.Production = tk.ToFloat64(rand.Intn(200), 0, tk.RoundingAuto)
			m.WindSpeed = tk.ToFloat64(rand.Intn(20), 0, tk.RoundingAuto)
			m.PerformanceIndex = tk.ToFloat64(rand.Intn(100), 0, tk.RoundingAuto)
			m.MachineAvail = tk.ToFloat64(rand.Intn(100), 0, tk.RoundingAuto)
			m.GridAvail = tk.ToFloat64(rand.Intn(100), 0, tk.RoundingAuto)

			alarm := rand.Intn(2)
			warning := rand.Intn(2)

			if alarm == 1 {
				m.IsAlarm = true
			}

			if warning == 1 {
				m.IsWarning = true
			}

			if k == "Tejuva" {
				m.Turbine = tejuvaProject[i]
			} else {
				m.Turbine = strings.ToUpper(k[:3] + tk.ToString(i))
			}

			ctx.Insert(m.New())
		}
	}

	log.Println("End generating Monitoring Data...")
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

	log.Println("DB Connect...")

	return c, nil
}

func ReadConfig() map[string]string {
	ret := make(map[string]string)
	file, err := os.Open(wd + "conf" + separator + "app.conf")
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
		log.Println(err.Error())
	}

	return ret
}
