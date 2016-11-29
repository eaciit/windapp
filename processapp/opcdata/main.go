package main

import (
	"bufio"
	//w "eaciit/wfdemo-git-dev/processapp/opcdata/filewatcher"
	d "eaciit/wfdemo-git-dev/processapp/opcdata/datareader"
	"github.com/eaciit/dbox"
	_ "github.com/eaciit/dbox/dbc/mongo"
	_ "github.com/eaciit/orm"
	tk "github.com/eaciit/toolkit"
	"os"
	"time"
	"fmt"
	"strings"
	"path/filepath"
	"eaciit/wfdemo-git/processapp/opcdata/model"
	"io/ioutil"
)
var pathSep = string(os.PathSeparator) 
var (
	wd = func() string {
		d, _ := filepath.Abs(filepath.Dir(os.Args[0]))
		return d + pathSep
	}()
)
var modified map[string]string
var	processed map[string]model.ProcessedLog

func main() {
	FileWatcher()
}

func FileIsExist(dirSources,targetFileName string)bool{
	
	//indx:=6
	//for indx>=0;indx--{
		
	//tk.Println(dirSources+"\\"+targetFileName)
	_,err:=os.Open(dirSources+"\\"+targetFileName)
	if err!=nil{
		if os.IsNotExist(err){
			tk.Println("File Not Exist")
			
		}else{
			tk.Println(err.Error())
		}
		return false
	}
	return true
	//}
	
}
func FileWatcher() {
	config := ReadConfig()
	
	dirSources := config["FileSources"]
	dirProcess := config["FileProcess"]
	scpDir := config["UploadDirectory"]
	sshUser := config["SSHUser"]
	sshServer := config["SSHServer"]
	//var FileName string
	var e bool
	now := time.Now()
	indx:=6
	
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
				_,e,start,end,rows:=d.NewDataReader(dirSources+"\\"+targetFileName, dirProcess, wd,scpDir,sshUser,sshServer).Start(0)
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
			
			_,e,start,end,rows:=d.NewDataReader(dirSources+"\\"+targetFileName, dirProcess, wd,scpDir,sshUser,sshServer).Start(0)
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
/*func WriteUnprocessedLog(){
	strBuff:=""
	for val,_:=range unprocessed{
		strBuff+=val+"\n"
	}
	buff:=[]byte{strBuff}
	err:=ioutil.WriteFile(wd + "conf/unprocessed.log",buff,0644)
	if err!=nil{
		panic(err)
	}
}func isInlist(list []string,item string) bool {
	for _,val:=range unprocessed{
		if val==item{
			return true
		}
	}
	return false
}
func ReadUnprocessedLog(){
	file, err := os.Open(wd + "conf/unprocessed.log")
	if err == nil {
		defer file.Close()
		reader := bufio.NewReader(file)
		for {
			line, _, e := reader.ReadLine()
			if e != nil {
				break
			}
			if !isInlist(unprocessed,line){
				unprocessed = append(unprocessed,line)
			}
			
			
			//sval := strings.Split(string(line), "=")
			//ret[sval[0]] = sval[1]
		}
	}
}*/

func WriteModified(filename string){
	strBuf:=""
	for val,key:=range modified{
		strBuf+=val+";"+key+"\n"
	}
	buff:=[]byte(strBuf)
	fmt.Println("Write Modified",wd+"log"+pathSep+"modified_"+filename+".csv")
	ioutil.WriteFile(wd+"log"+pathSep+"modified_"+filename+".csv",buff,0644)
}
func WriteProcessed(filename string){
	strBuf:=""
	for _,val:=range processed{
		strBuf+=val.ToString()
	}
	buff:=[]byte(strBuf)
	fmt.Println("Write Processed",wd+"log"+pathSep+"processed_"+filename+".csv")
	ioutil.WriteFile(wd+"log"+pathSep+"processed_"+filename+".csv",buff,0644)
}
func ReadProcessed(filename string) map[string]model.ProcessedLog {
	ret := make(map[string]model.ProcessedLog)
	file, err := os.Open(wd + "log"+pathSep+"processed_"+filename+".csv")
	if err == nil {
		defer file.Close()
		reader := bufio.NewReader(file)
		for {
			line, _, e := reader.ReadLine()
			if e != nil {
				break
			}
			newPF := model.FromString(string(line))
			ret[(*newPF).Filename] = *newPF
		}
	}
	return ret
}
func ReadModified(filename string) map[string]string {
	ret := make(map[string]string)
	file, err := os.Open(wd + "log"+pathSep+"modified_"+filename+".csv")
	if err == nil {
		defer file.Close()

		reader := bufio.NewReader(file)
		for {
			line, _, e := reader.ReadLine()
			if e != nil {
				break
			}

			sval := strings.Split(string(line), ";")
			ret[sval[0]] = sval[1]
		}
	} else {
		tk.Println(err.Error())
	}

	return ret
}
func ReadConfig() map[string]string {
	ret := make(map[string]string)
	file, err := os.Open(wd + "conf"+pathSep+"app.conf")
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
