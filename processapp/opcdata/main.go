package main

import (
	"bufio"
	//w "eaciit/wfdemo-git/processapp/opcdata/filewatcher"
	d "eaciit/wfdemo-git/processapp/opcdata/datareader"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/eaciit/dbox"
	_ "github.com/eaciit/dbox/dbc/mongo"
	_ "github.com/eaciit/orm"
	tk "github.com/eaciit/toolkit"
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
func FileIsExist(dirSources string) (bool, string) {
	now := time.Now()
	year, month, day := now.Date()
	hour, _, _ := now.Clock()

	targetFileName := fmt.Sprintf("DataFile%d%02d%02d-%02d.csv", year, month, day, hour)
	tk.Println(dirSources + "\\" + targetFileName)
	_, err := os.Open(dirSources + "\\" + targetFileName)
	if err != nil {
		if os.IsNotExist(err) {
			tk.Println("File Not Exist")

		} else {
			tk.Println(err.Error())
		}
		return false, ""
	}
	return true, targetFileName
}
func FileWatcher() {
	config := ReadConfig()
	dirSources := config["FileSources"]
	dirProcess := config["FileProcess"]
	scpDir := config["UploadDirectory"]
	sshUser := config["SSHUser"]
	sshServer := config["SSHServer"]
	var FileName string
	var e bool
	//watcher := w.NewFileWatcher(dirSources, dirProcess, wd,scpDir)
	if e, FileName = FileIsExist(dirSources); !e {
		os.Exit(1)
	}
	d.NewDataReader(dirSources+"\\"+FileName, dirProcess, wd, scpDir, sshUser, sshServer).Start()
	//watcher.StartWatcher()
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
