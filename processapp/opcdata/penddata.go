package main

import (
	"bufio"
	//w "eaciit/wfdemo-git/processapp/opcdata/filewatcher"
	d "eaciit/wfdemo-git/processapp/opcdata/datareader"
	"eaciit/wfdemo-git/processapp/opcdata/model"
	"flag"
	"fmt"
	"github.com/eaciit/dbox"
	_ "github.com/eaciit/dbox/dbc/mongo"
	_ "github.com/eaciit/orm"
	. "github.com/eaciit/sshclient"
	tk "github.com/eaciit/toolkit"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var (
	wd = func() string {
		d, _ := os.Getwd()
		return d + "/"
	}()
)
var pathSep = string(os.PathSeparator)
var modified = map[string]string{}
var processed = map[string]model.ProcessedLog{}

var sdate, edate time.Time
var sshUser, sshServer, scpDir string

func main() {
	intstartdate, intenddate := int(20161208), int(20161208)
	flag.IntVar(&intstartdate, "sdate", 20161208, "Start date for processing data")
	flag.IntVar(&intenddate, "edate", 20161208, "End date for processing data")
	flag.Parse()

	sdate = tk.String2Date(tk.ToString(intstartdate), "YYYYMMdd").UTC()
	edate = tk.String2Date(tk.ToString(intenddate), "YYYYMMdd").UTC()

	tk.Printfn(">>> %v to %v ", sdate.Format("2006-01-02"), edate.Format("2006-01-02"))

	FileWatcher()
}

func FileIsExist(dirSources string, targetFileName string) bool {
	_, err := os.Open(dirSources + "\\" + targetFileName)
	if err != nil {
		if os.IsNotExist(err) {
			tk.Println(">>> File Not Exist : ", targetFileName)

		} else {
			tk.Println(err.Error())
		}
		return false
	}
	return true
}

func FileWatcher() {
	config := ReadConfig()
	dirSources := config["FileSources"]
	dirProcess := config["FileProcess"]
	scpDir = config["UploadDirectory"]
	sshUser = config["SSHUser"]
	sshServer = config["SSHServer"]

	// tk.Printfn(">>> Start Program : %v ==================\n\n", t0)

	now := sdate
	for !now.After(edate) {
		t0 := time.Now()
		ifile := 0
		// arrfile := []string{}
		tk.Printfn(">>> Start Date : %v ==================\n\n", now.Format("2006-01-02"))

		// _ci := make(chan int, 24)
		_si := make(chan int, 24)
		_fname := make(chan string, 24)

		go sendFile(_fname, _si)

		year, month, day := now.Date()
		modFileName := fmt.Sprintf("%d%02d%02d", year, month, day)

		for hour := 0; hour < 24; hour++ {
			// go func(_hour int) {
			t1 := time.Now()
			targetFileName := fmt.Sprintf("DataFile%d%02d%02d-%02d.csv", year, month, day, hour)

			tk.Printfn(">>> Start Process : %s ==================", targetFileName)

			if cond := FileIsExist(dirSources, targetFileName); !cond {
				// _ci <- 1
				tk.Printf(">>> : ================== \n\n")
				continue
			}

			fmt.Println(">>>", dirSources+pathSep+targetFileName)
			info, _ := os.Stat(dirSources + pathSep + targetFileName)

			modifiedTimeFS := info.ModTime()
			modifiedTimeFS, _ = time.Parse("02-Jan-2006 15:04:05", modifiedTimeFS.Format("02-Jan-2006 15:04:05"))

			dr := d.NewDataReader(dirSources+"\\"+targetFileName, dirProcess, wd, scpDir, sshUser, sshServer)
			_, e, start, end, rows := dr.Start(0)
			if e == nil {
				modified[targetFileName] = modifiedTimeFS.Format("02-Jan-2006 15:04:05")
				newPf := new(model.ProcessedLog)
				(*newPf).Filename = targetFileName
				(*newPf).StartTime = start.Format("02-Jan-2006 15:04:05")
				(*newPf).EndTime = end.Format("02-Jan-2006 15:04:05")
				(*newPf).RowIndex = rows
				processed[targetFileName] = *newPf
			}

			_fname <- dr.ZipName
			// arrfile = append(arrfile, dr.ZipName)
			// dr.SendFile(dr.ZipName)

			ifile++
			tk.Printfn(">>> End Process : %s in %s ==================", targetFileName, time.Since(t1).String())
			tk.Printf("\n\n")
			// 	_ci <- 1
			// }(hour)

			// if hour != 0 {
			// 	<-_ci
			// }
		}

		close(_fname)

		WriteModified(modFileName)
		WriteProcessed(modFileName)
		_ptstring := time.Since(t0).String()
		tk.Printfn(">>> %v Processed in %s ==================\n\n", now.Format("2006-01-02"), _ptstring)

		for _i := 0; _i < ifile; _i++ {
			<-_si
		}

		close(_si)
		// close(_ci)
		// sendFile(arrfile)

		WriteDateProcess(now.Format("2006-01-02"), _ptstring, time.Since(t0).String(), ifile)
		tk.Printfn(">>> %v Done in %s ==================\n\n", now.Format("2006-01-02"), time.Since(t0).String())

		now = now.AddDate(0, 0, 1)
	}
}

func sendFile(cfilename chan string, result chan<- int) {
	locationTarget := scpDir
	_listfail := tk.M{}
	tk.Println("Sending file")
	filename := ""
	for filename = range cfilename {
		_t0 := time.Now()
		ssh := new(SshSetting)

		ssh.SSHAuthType = SSHAuthType_Certificate
		ssh.SSHHost = sshServer
		ssh.SSHUser = sshUser
		ssh.SSHKeyLocation = wd + "\\conf\\key\\developer.pem"

		// _, err := ssh.Connect()
		// if err != nil {
		// 	tk.Println("Error connecting to server: " + err.Error())
		// } else {
		// 	tk.Println("Connected to server!")
		// }

		file, err := os.Open(filename)
		if err != nil {
			tk.Println("Error opening file: " + err.Error())
			os.Exit(1)
		}

		fileStat, err := file.Stat()
		if err != nil {
			tk.Println("Error opening file: " + err.Error())
			os.Exit(1)
		}

		tk.Println(">>> SENDING ", filename)
		for true {
			ilf := _listfail.GetInt(filename)
			err = ssh.SshCopyByFile(file, fileStat.Size(), fileStat.Mode().Perm(), filepath.Base(fileStat.Name()), locationTarget)
			if err != nil {
				tk.Println("Error : ", err.Error())
				// _, err = ssh.Connect()
				// if err != nil {
				// 	tk.Println("Error connecting to server: " + err.Error())
				// } else {
				// 	tk.Println("Connected to server!")
				// }
				if ilf < 5 {
					ilf += 1
					_listfail.Set(filename, ilf)

					tk.Println("Try to resend : ", filename, " - ", ilf)
					cfilename <- filename
				} else {
					WriteFileFail(filename)
					result <- 1
				}
				break
			} else {
				arrname := strings.Split(file.Name(), "\\")
				_lnote := tk.Sprintf("%v;%s;%v;%s;%d \r\n", time.Now().UTC(), arrname[len(arrname)-1], fileStat.Size(), time.Since(_t0).String(), ilf)
				WriteFileSend(_lnote)
				tk.Println(_lnote)

				//Remove file
				//===========
				file.Close()
				_ = os.Remove(filename)
				//===========
				result <- 1
				break
			}
		}

	}
	tk.Println("Sending file successfully")
}

func ReadModified(filename string) map[string]string {
	ret := make(map[string]string)
	file, err := os.Open(wd + "log" + pathSep + "modified_" + filename + ".csv")
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

func WriteModified(filename string) {
	strBuf := ""
	for val, key := range modified {
		strBuf += val + ";" + key + "\n"
	}
	buff := []byte(strBuf)
	fmt.Println("Write Modified", wd+"log"+pathSep+"modified_"+filename+".csv")
	ioutil.WriteFile(wd+"log"+pathSep+"modified_"+filename+".csv", buff, 0644)
}

func WriteDateProcess(date, pdurr, durr string, ifile int) {
	f, _ := os.OpenFile(wd+pathSep+"datelog.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	defer f.Close()
	_, _ = f.WriteString(tk.Sprintf("%s;%s;%s;%v files\r\n", date, pdurr, durr, ifile))
}

func WriteFileSend(_txt string) {
	f, _ := os.OpenFile(wd+pathSep+"filesend.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	defer f.Close()
	_, _ = f.WriteString(_txt)
}

func WriteFileFail(_txt string) {
	f, _ := os.OpenFile(wd+pathSep+"failsend.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	defer f.Close()
	_, _ = f.WriteString(_txt)
}

func WriteProcessed(filename string) {
	strBuf := ""
	for _, val := range processed {
		strBuf += val.ToString()
	}
	buff := []byte(strBuf)
	fmt.Println("Write Processed", wd+"log"+pathSep+"processed_"+filename+".csv")
	ioutil.WriteFile(wd+"log"+pathSep+"processed_"+filename+".csv", buff, 0644)
}

func ReadProcessed(filename string) map[string]model.ProcessedLog {
	ret := make(map[string]model.ProcessedLog)
	file, err := os.Open(wd + "log" + pathSep + "processed_" + filename + ".csv")
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
