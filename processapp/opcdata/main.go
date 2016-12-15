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

	now := time.Now()
	indx:=6
	arrDataReader := []d.DataReader{}
	for ;indx>=0;indx--{
		prev:=now.Add(time.Duration(-indx)*time.Hour)
		year,month,day:=prev.Date()
		hour,_,_:=prev.Clock()
		modFileName := fmt.Sprintf("%d%02d%02d",year,month,day)
		modified = ReadModified(modFileName)
		processed = ReadProcessed(modFileName)
		targetFileName := fmt.Sprintf("DataFile%d%02d%02d-%02d.csv",year,month,day,hour)
			
		//watcher := w.NewFileWatcher(dirSources, dirProcess, wd,scpDir)
		if e=FileIsExist(dirSources,targetFileName);!e{
			continue
		}
		fmt.Println(dirSources+pathSep+targetFileName)
		info, _ := os.Stat(dirSources+pathSep+targetFileName)
		modifiedTimeFS := info.ModTime()
		modifiedTimeFS,_ = time.Parse("02-Jan-2006 15:04:05",modifiedTimeFS.Format("02-Jan-2006 15:04:05"))
		if _,ok:=modified[targetFileName];ok{
			lastModTimeLog,_ := time.Parse("02-Jan-2006 15:04:05",modified[targetFileName])
			fmt.Println(lastModTimeLog.Format("02-Jan-2006 15:04:05"),modifiedTimeFS.Format("02-Jan-2006 15:04:05"))
			if modifiedTimeFS.After(lastModTimeLog){
				fmt.Println("File ",targetFileName,"Modified")
				dr := d.NewDataReader(dirSources+"\\"+targetFileName, dirProcess, wd,scpDir,sshUser,sshServer)
				_,e,start,end,rows:=dr.Start(0)
				arrDataReader = append(arrDataReader,*dr)
				if e==nil{
					modified[targetFileName]=modifiedTimeFS.Format("02-Jan-2006 15:04:05")
					newPf:=processed[targetFileName]
					newPf.StartTime = start.Format("02-Jan-2006 15:04:05")
					newPf.EndTime = end.Format("02-Jan-2006 15:04:05")
					newPf.RowIndex = rows
					processed[targetFileName]=newPf
				}else{
					fmt.Println(e.Error())
				}
			}else{
				fmt.Println(lastModTimeLog.Nanosecond(),modifiedTimeFS.Nanosecond())
			}
		}else{
			dr := d.NewDataReader(dirSources+"\\"+targetFileName, dirProcess, wd,scpDir,sshUser,sshServer)
			_,e,start,end,rows:=dr.Start(0)
			arrDataReader = append(arrDataReader,*dr)
			if e==nil{
				modified[targetFileName]=modifiedTimeFS.Format("02-Jan-2006 15:04:05")
				newPf:=new(model.ProcessedLog)
				(*newPf).Filename = targetFileName
				(*newPf).StartTime = start.Format("02-Jan-2006 15:04:05")
				(*newPf).EndTime = end.Format("02-Jan-2006 15:04:05")
				(*newPf).RowIndex = rows
				processed[targetFileName]=*newPf
			}
			
		}
		
		WriteModified(modFileName)
		WriteProcessed(modFileName)
		
	}
	for oo:=len(arrDataReader)-1;oo>=0;oo--{
		arrDataReader[oo].SendFile(arrDataReader[oo].ZipName)
	}
	
//=======
	//watcher := w.NewFileWatcher(dirSources, dirProcess, wd,scpDir)
//	if e, FileName = FileIsExist(dirSources); !e {
//		os.Exit(1)
//	}
//	d.NewDataReader(dirSources+"\\"+FileName, dirProcess, wd, scpDir, sshUser, sshServer).Start()
//>>>>>>> e93aaba7699185484bd3b42e861e2779bb99342c
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
