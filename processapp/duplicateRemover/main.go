package main

import (
	"bufio"
	. "eaciit/wfdemo-git/processapp/duplicateRemover/checker"
	"log"
	"os"
	"runtime"
	"strings"

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

	log.Println("Starting checking duplicate data...")
	conn, err := PrepareConnection()
	if err != nil {
		log.Println(err)
	}
	ctx := orm.New(conn)

	/*scadaOEM := NewOEMChecker(ctx)
	scadaOEM.Run()*/

	scadaData := NewScadaChecker(ctx)
	scadaData.Run()

	log.Println("End remove duplicate data...")
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
