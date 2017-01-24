package main

import (
	c3to10 "eaciit/wfdemo-git/processapp/threeextractor/convert3secto10min/lib"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/eaciit/database/base"
	"github.com/eaciit/orm"
	"github.com/fsnotify/fsnotify"
	// "github.com/metakeule/fmtdate"
	//dc "eaciit/wfdemo-git/processapp/threeextractor/dataconversion"
	"archive/tar"
	"archive/zip"
	"bytes"

	"eaciit/wfdemo-git/library/helper"
	. "eaciit/wfdemo-git/library/models"
	. "eaciit/wfdemo-git/processapp/watcher/controllers"

	tk "github.com/eaciit/toolkit"
	"gopkg.in/mgo.v2/bson"

	"github.com/eaciit/dbox"
	_ "github.com/eaciit/dbox/dbc/mongo"
)

const (
	NOK = "NOK"
	OK  = "OK"
	// TIME_FORMAT_STR = "DDMMYYhhmmss"
)

type Command struct {
	Action  string
	Command string
	Success string
	Fail    string
}
type Configuration struct {
	Draft    string
	Process  string
	Fail     string
	Success  string
	Errors   string
	Archive  string
	Commands []Command
}

var (
	conn    base.IConnection
	ctx     *orm.DataContext
	conf    Configuration
	mux     = &sync.Mutex{}
	pathSep = string(os.PathSeparator)
)

func main() {

	fmt.Println(".." + pathSep + "conf" + pathSep + "receiver-config.json")

	file, _ := os.Open(".." + pathSep + "conf" + pathSep + "receiver-config.json")
	decoder := json.NewDecoder(file)
	err := decoder.Decode(&conf)

	if err != nil {
		panic(err)
	}

	// now, _ := time.Parse("2006-1-2 15:4:05", "2016-10-22 23:57:30")
	// tenMinInfo := GenTenMinuteInfo(now)

	// log.Printf("%#v \n", tenMinInfo)

	runtime.GOMAXPROCS(runtime.NumCPU())
	log.Println("Starting the app..\n")

	log.Println()
	log.Printf("Draft: %v\n", conf.Draft)
	log.Printf("Process: %v\n", conf.Process)
	log.Printf("Fail: %v\n", conf.Fail)
	log.Printf("Success: %v\n", conf.Success)
	log.Printf("Errors: %v\n", conf.Errors)
	log.Println()

	/*db, _ := PrepareConnection()
	base := new(BaseController)
	base.Ctx = orm.New(db)
	defer base.Ctx.Close()

	log.Println("xxx")

	new(GenTenFromThreeSecond).Generate(base, "sample.csv")*/

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	done := make(chan bool)
	go func() {
		for {
			select {
			case event := <-watcher.Events:
				// log.Println("event:", event)
				if event.Op&fsnotify.Create == fsnotify.Create {
					go processFile(event.Name, conf.Commands)
				}
			case err := <-watcher.Errors:
				log.Println("error:", err)
			}
		}
	}()

	// watch draft
	err = watcher.Add(conf.Draft)
	if err != nil {
		log.Fatal(err)
	}

	// watch process
	err = watcher.Add(conf.Process)
	if err != nil {
		log.Fatal(err)
	}

	/*// watch fail
	err = watcher.Add("/Users/frezadev/Documents/watch/fail")
	if err != nil {
		log.Fatal(err)
	}

	// watch success
	err = watcher.Add("/Users/frezadev/Documents/watch/success")
	if err != nil {
		log.Fatal(err)
	}*/
	<-done
}

func processFile(filePath string, com []Command) {
	// dt := time.Now()
	// dtStr := fmtdate.Format(TIME_FORMAT_STR, dt)
	fmt.Println(" >>> Process : ", filePath)
	for true {
		byteOut, err := runCMD("lsof " + filePath)
		if err != nil {
			log.Print("Gagal")
		}
		if len(byteOut) == 0 {
			break
		} else {
			// fmt.Println(string(byteOut))
			time.Sleep(5 * time.Second)
		}

	}
	file := strings.Split(filePath, pathSep)
	fileName := file[len(file)-1]

	if strings.Contains(fileName, ".csv") {
		time.Sleep(100 * time.Millisecond)
		log.Printf("Proccess file: %v \n", filePath)
		var action Command

		if strings.Contains(filePath, conf.Draft) {
			action = com[0]
		} else if strings.Contains(filePath, conf.Process) {
			action = com[1]
		}

		// log.Printf("action: %v \n", action)

		next := action.Action
	done:

		for {
			log.Printf("\n\nnext: %v \n", next)

			if next == "DONE" {
				break done
			} else {
				for _, act := range com {
					//fmt.Println("PPPP")
					if act.Action == next {
						next = run(act, fileName)

						break
					}
				}
			}
		}

		log.Printf("DONE for file: %v\n", filePath)
	} else if strings.Contains(fileName, ".tar") {
		time.Sleep(100 * time.Millisecond)
		untar(filePath, conf.Draft)
		_, _ = runCMD(fmt.Sprintf("mv %v %v", filePath, conf.Archive))
	} else if strings.Contains(fileName, ".zip") {
		time.Sleep(100 * time.Millisecond)
		e := unzip(filePath, conf.Draft)
		if e != nil {
			log.Printf("%s - %s", filePath, e.Error())
			_, _ = runCMD(fmt.Sprintf("mv %v %v", filePath, conf.Errors))
		} else {
			fmt.Println("Unzip Done")
			_, _ = runCMD(fmt.Sprintf("mv %v %v", filePath, conf.Archive))
		}
	}
}

func untar(tarball, target string) error {
	reader, err := os.Open(tarball)
	if err != nil {
		return err
	}
	defer reader.Close()
	tarReader := tar.NewReader(reader)

	for {
		header, err := tarReader.Next()

		if err == io.EOF {
			break
		} else if err != nil {
			log.Printf("ERR: %#v \n", err.Error())
			return err
		}

		path := filepath.Join(target, header.Name)
		info := header.FileInfo()
		// log.Printf("path: %#v \n", path)
		// log.Printf("info: %#v \n", info)
		if strings.Contains(info.Name(), ".csv") {
			if info.IsDir() {
				if err = os.MkdirAll(path, info.Mode()); err != nil {
					return err
				}
				continue
			}

			file, err := os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, info.Mode())
			if err != nil {
				return err
			}
			defer file.Close()
			_, err = io.Copy(file, tarReader)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func unzip(archive, target string) error {
	fmt.Println("Unzip", archive, target)
	reader, err := zip.OpenReader(archive)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(target, 0755); err != nil {
		return err
	}

	for _, file := range reader.File {
		if strings.Contains(file.Name, ".csv") {

			path := filepath.Join(target, file.Name)
			fileReader, err := file.Open()

			if err != nil {
				log.Printf("ERR: %#v \n", err.Error())
				return err
			}
			defer fileReader.Close()

			targetFile, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
			if err != nil {
				return err
			}
			defer targetFile.Close()

			if _, err := io.Copy(targetFile, fileReader); err != nil {
				return err
			}
		}
	}

	reader.Close()
	return nil
}

func runCMD(cmdStr string) (out []byte, err error) {
	if pathSep == "\\" {
		cmdStr = strings.Replace(cmdStr, "\\", "/", -1)

	}

	if !strings.Contains(cmdStr, "lsof") {
		fmt.Println("sh", "-c", cmdStr)
	}

	var errBuff bytes.Buffer
	cmd := exec.Command("sh", "-c", cmdStr)
	cmd.Stderr = &errBuff
	//cmd.Run()
	//cmd.Path = os.Getenv("Path")
	//fmt.Println(cmd.Path)
	out, err = cmd.Output()
	if err != nil {
		fmt.Println(errBuff.String())
	}
	return
}

func run(action Command, file string) (next string) {
	cmdStr := ""
	runCommand := true
	log.Printf("run: %v | %v \n", action, file)
	// mux.Lock()

	if action.Action == "COPY_TO_PROCESS" {
		cmdStr = fmt.Sprintf(action.Command, conf.Draft+pathSep+file, conf.Process+pathSep+file)
	} else if action.Action == "COPY_TO_SUCCESS" {
		runCommand = doProcess(conf.Process + pathSep + file)

		cmdStr = fmt.Sprintf(action.Command, conf.Process+pathSep+file, conf.Success+pathSep+file)
	} else if action.Action == "COPY_TO_FAIL" {
		cmdStr = fmt.Sprintf(action.Command, conf.Process+pathSep+file, conf.Fail+pathSep+file)
	}
	//log.Printf("cmdstr: %v \n", cmdStr)
	// mux.Unlock()

	_ismvsucces := false
	if runCommand {
		out, err := runCMD(cmdStr)

		if out != nil {
			log.Printf("%v \n", out)
		}

		if err != nil {
			log.Printf("result: %v %s\n%s", err.Error(), cmdStr, string(out))
		} else {
			_ismvsucces = true
		}

		next = action.Success
	} else {
		log.Println("DONE")
		next = action.Fail
	}

	if action.Action == "COPY_TO_SUCCESS" && runCommand && _ismvsucces {
		_cmd := fmt.Sprintf("rm %v", filepath.Join(conf.Success, file))
		_, _ = runCMD(_cmd)
	}

	// log.Printf("next: %v \n", next)

	return
}

func doProcess(file string) (success bool) {
	log.Printf("doProcess: %v \n", file)
	db, e := PrepareConnection()
	if e != nil {
		log.Printf("ERROR on Process: %v\n", e.Error())
		success = false
	} else {
		//start := time.Now()
		base := new(BaseController)
		base.Ctx = orm.New(db)
		defer base.Ctx.Close()

		anFile := strings.Split(file, pathSep)
		fileName := anFile[len(anFile)-1]

		muxDo := &sync.Mutex{}

		muxDo.Lock()
		errorLine := new(ConvScadaTreeSecs).Generate(base, file)
		WriteWatcherErrors(errorLine, fileName, conf.Errors)
		log.Println("ConvScadaTreeSecs: DONE")
		muxDo.Unlock()
		log.Println("Begin Converting ", file)
		err := c3to10.Generate(tk.M{}.Set("selector", "file").Set("file", fileName))
		if err != nil {
			tk.Println(err)
		} else {
			UpdateLastHFDAvail()
			UpdateLastMonitoring()
			tk.Println(">> DONE <<")
		}
	}

	return
}

func UpdateLastHFDAvail() {

	_nt0 := time.Now()
	tk.Println("Start Update Last HDF Available ...")

	var workerconn dbox.IConnection
	for {
		var err error
		workerconn, err = PrepareConnection()
		if err == nil {
			break
		} else {
			tk.Printfn("==#DB-ERRCONN==\n %s \n", err.Error())
			<-time.After(time.Second * 3)
		}
	}
	defer workerconn.Close()

	type latestdataperiod struct {
		ID          bson.ObjectId ` bson:"_id" , json:"_id" `
		Projectname string
		Type        string
		Data        []time.Time
	}

	csr, err := workerconn.NewQuery().
		Select().
		From("LatestDataPeriod").
		Where(dbox.Eq("type", "ScadaDataHFD")).
		Cursor(nil)

	if err != nil || csr.Count() == 0 {
		return
	}

	_dt := new(latestdataperiod)

	_ = csr.Fetch(_dt, 1, false)
	csr.Close()

	pipes := []tk.M{tk.M{"$group": tk.M{"_id": "$projectname",
		"mintimestamp": tk.M{"$min": "$timestamp"},
		"maxtimestamp": tk.M{"$max": "$timestamp"},
	}}}

	xcsr, err := workerconn.NewQuery().
		From(new(ScadaConvTenMin).TableName()).
		Command("pipe", pipes).
		Cursor(nil)
	if err != nil {
		return
	}

	_tkm := tk.M{}
	_ = xcsr.Fetch(&_tkm, 1, false)
	xcsr.Close()

	_min := _tkm.Get("mintimestamp", time.Time{}).(time.Time)
	_max := _tkm.Get("maxtimestamp", time.Time{}).(time.Time)

	_dt.Data[0] = _min
	_dt.Data[1] = _max

	_ = workerconn.NewQuery().
		From("LatestDataPeriod").
		SetConfig("multiexec", true).
		Save().Exec(tk.M{}.Set("data", _dt))

	tk.Println(" >>> End Update Last HDF Available in ", time.Since(_nt0).String())
}

func UpdateLastMonitoring() {
	_nt0 := time.Now()
	tk.Println(" >>> Start Update Last Monitoring ...")

	var workerconn dbox.IConnection
	for {
		var err error
		workerconn, err = PrepareConnection()
		if err == nil {
			break
		} else {
			tk.Printfn("==#DB-ERRCONN==\n %s \n", err.Error())
			<-time.After(time.Second * 3)
		}
	}
	defer workerconn.Close()

	var sworkerconn dbox.IConnection
	for {
		var err error
		sworkerconn, err = PrepareConnection()
		if err == nil {
			break
		} else {
			tk.Printfn("==#DB-ERRCONN==\n %s \n", err.Error())
			<-time.After(time.Second * 3)
		}
	}
	defer sworkerconn.Close()

	type latestdataperiod struct {
		ID          bson.ObjectId ` bson:"_id" , json:"_id" `
		Projectname string
		Type        string
		Data        []time.Time
	}

	csr, err := workerconn.NewQuery().
		Select().
		From("LatestDataPeriod").
		Where(dbox.Eq("type", "ScadaDataHFD")).
		Cursor(nil)

	if err != nil || csr.Count() == 0 {
		return
	}

	_dt := new(latestdataperiod)

	_ = csr.Fetch(_dt, 1, false)
	csr.Close()

	speriode := _dt.Data[1].AddDate(0, 0, -1)
	eperiode := _dt.Data[1]

	err = workerconn.NewQuery().
		Delete().
		From(new(Monitoring).TableName()).
		Where(dbox.Lte("timestamp", speriode)).
		Exec(nil)

	if err != nil {
		tk.Println(">>> Error found on Delete : ", err.Error())
	}

	msmonitor := PrepareMasterMonitoring()
	tk.Println(">>> periode ", speriode, " ----- ", eperiode)
	xcsr, err := workerconn.NewQuery().
		Select("timestamp", "projectname", "turbine", "fast_activepower_kw", "fast_windspeed_ms", "fast_rotorspeed_rpm",
			"slow_tempnacelle", "fast_pitchangle").
		From(new(ScadaConvTenMin).TableName()).
		Where(dbox.And(dbox.Lte("timestamp", eperiode), dbox.Gt("timestamp", speriode))).
		Order("timestamp").
		Cursor(nil)

	if err != nil {
		return
	}

	sqsave := sworkerconn.NewQuery().
		From(new(Monitoring).TableName()).
		SetConfig("multiexec", true).
		Save()

	// _lstatus := make(map[string]Monitoring, 0)
	for {
		_tkm := tk.M{}
		err = xcsr.Fetch(&_tkm, 1, false)
		if err != nil {
			break
		}

		_timestamp := _tkm.Get("timestamp", time.Time{}).(time.Time)
		_key := tk.Sprintf("%s#%s#%s",
			_tkm.GetString("projectname"),
			_tkm.GetString("turbine"),
			_timestamp.Format("060102_150405"),
		)
		_monitor := Monitoring{}

		if _mo, _bo := msmonitor[_key]; _bo {
			_monitor = _mo

			// 	if _mo.Status != "" {
			// 		_astatus := Monitoring{}

			// 		_astatus.Status = _mo.Status
			// 		_astatus.Type = _mo.Type
			// 		_astatus.StatusCode = _mo.StatusCode
			// 		_astatus.StatusDesc = _mo.StatusDesc

			// 		_lstatus[_mo.Turbine] = _astatus
			// 	}
			// } else if _lsdata, _lscond := _lstatus[_tkm.GetString("turbine")]; _lscond {
			// 	_monitor.Status = _lsdata.Status
			// 	_monitor.Type = _lsdata.Type
			// 	_monitor.StatusCode = _lsdata.StatusCode
			// 	_monitor.StatusDesc = _lsdata.StatusDesc
		}

		if _monitor.Status == "" {
			_monitor.Status = "N/A"
		}

		_monitor.ID = _key
		_monitor.TimeStamp = _timestamp
		_monitor.DateInfo = helper.GetDateInfo(_timestamp)
		_monitor.LastUpdate = _nt0
		_monitor.LastUpdateDateInfo = helper.GetDateInfo(_nt0)
		_monitor.Project = _tkm.GetString("projectname")
		_monitor.Turbine = _tkm.GetString("turbine")

		if _val := _tkm.GetFloat64("fast_activepower_kw"); _val != -9999999 {
			_monitor.Production = (_val / 1000) / 6
		} else {
			_monitor.Production = _val
		}

		// if _val := _tkm.GetFloat64("fast_windspeed_ms"); _val != -9999999 {
		_monitor.WindSpeed = _tkm.GetFloat64("fast_windspeed_ms")
		// }

		// if _val := _tkm.GetFloat64("fast_rotorspeed_rpm"); _val != -9999999 {
		_monitor.RotorSpeedRPM = _tkm.GetFloat64("fast_rotorspeed_rpm")
		// }
		_monitor.PitchAngle = _tkm.GetFloat64("fast_pitchangle")
		_monitor.WindDirection = _tkm.GetFloat64("slow_tempnacelle")
		// WindDirection
		//"slow_tempnacelle","fast_pitchangle"

		_ = sqsave.Exec(tk.M{}.Set("data", _monitor))

	}
	xcsr.Close()

	tk.Println(" >>> End Update Last Monitoring in ", time.Since(_nt0).String())
}

func PrepareMasterMonitoring() (_mnt map[string]Monitoring) {
	_mnt = make(map[string]Monitoring)

	var workerconn dbox.IConnection
	for {
		var err error
		workerconn, err = PrepareConnection()
		if err == nil {
			break
		} else {
			tk.Printfn("==#DB-ERRCONN==\n %s \n", err.Error())
			<-time.After(time.Second * 3)
		}
	}
	defer workerconn.Close()

	xcsr, err := workerconn.NewQuery().
		Select().
		From(new(Monitoring).TableName()).
		Cursor(nil)

	if err != nil {
		return
	}

	defer xcsr.Close()

	for {
		_amnt := Monitoring{}
		err = xcsr.Fetch(&_amnt, 1, false)
		if err != nil {
			break
		}

		_mnt[_amnt.ID] = _amnt
	}

	return
}
