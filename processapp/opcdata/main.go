package main

import (
	"bufio"
	w "eaciit/ostrowfm/processapp/opcdata/filewatcher"
	"github.com/eaciit/dbox"
	_ "github.com/eaciit/dbox/dbc/mongo"
	_ "github.com/eaciit/orm"
	tk "github.com/eaciit/toolkit"
	"os"
	"strings"
)

var (
	wd = func() string {
		d, _ := os.Getwd()
		return d + "/"
	}()
)

func main() {
	FileWatcher()
}

func FileWatcher() {
	config := ReadConfig()
	dirSources := config["FileSources"]
	dirProcess := config["FileProcess"]

	watcher := w.NewFileWatcher(dirSources, dirProcess, wd)
	watcher.StartWatcher()
}

func CsvExtractor() {
	//conn, err := PrepareConnection()
	// if err != nil {
	// 	tk.Println("Error connection: ", err.Error())
	// }
	//ctx := orm.New(conn)

	// config := ReadConfig()
	// dirSources := config["FileSources"]
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
