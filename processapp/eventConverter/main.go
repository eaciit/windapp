package main

import (
	"bufio"
	. "eaciit/wfdemo-git-dev/processapp/eventConverter/conversion"
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

	tk.Println("Starting convert event data to down...")
	conn, err := PrepareConnection()
	if err != nil {
		tk.Println(err)
	}
	ctx := orm.New(conn)

	config := ReadConfig()
	dir := config["FileProcess"]

	event := NewEventRawConversion(ctx, dir)
	event.Run()

	down := NewDownConversion(ctx, dir)
	down.Run()

	tk.Println("End processing event data to down...")
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

	tk.Println("DB Connect...")

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
		tk.Println(err.Error())
	}

	return ret
}
