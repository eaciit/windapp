package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"sync"

	"github.com/eaciit/database/base"
	"github.com/eaciit/orm"
	"github.com/fsnotify/fsnotify"
	// "github.com/metakeule/fmtdate"
	//dc "eaciit/wfdemo-git/processapp/threeextractor/dataconversion"
	"archive/tar"
	"archive/zip"
	"bufio"
	"bytes"
	"time"

	"eaciit/wfdemo-git/library/helper"
	. "eaciit/wfdemo-git/library/models"
	. "eaciit/wfdemo-git/processapp/watcher/controllers"

	econv "eaciit/wfdemo-git/processapp/eventHFDConverter/conversion"

	tk "github.com/eaciit/toolkit"
	"gopkg.in/mgo.v2/bson"

	"github.com/eaciit/dbox"
	_ "github.com/eaciit/dbox/dbc/mongo"
)

const (
	NOK = "NOK"
	OK  = "OK"
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
	conn base.IConnection
	ctx  *orm.DataContext
	conf Configuration
	mux  = &sync.Mutex{}

	pathSep = string(os.PathSeparator)

	wd = func() string {
		d, _ := os.Getwd()
		return d + "/"
	}()

	// masteralarmbrake = tk.M{}
)

func main() {

	_fconf := filepath.Join(wd, "..", "conf", "alert-config.json")

	fmt.Println(_fconf)
	file, _ := os.Open(_fconf)
	decoder := json.NewDecoder(file)
	err := decoder.Decode(&conf)

	if err != nil {
		panic(err)
	}

	log.Println("Starting the app..\n")

	log.Println()
	log.Printf("Draft: %v\n", conf.Draft)
	log.Printf("Process: %v\n", conf.Process)
	log.Printf("Fail: %v\n", conf.Fail)
	log.Printf("Success: %v\n", conf.Success)
	log.Printf("Errors: %v\n", conf.Errors)
	log.Println()

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
				if event.Op&fsnotify.Create == fsnotify.Create {
					go processfile(event.Name, conf.Commands)
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

	<-done
}

func processfile(filePath string, com []Command) {
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
		next := action.Action
	done:
		for {
			log.Printf("\n\nnext: %v \n", next)
			if next == "DONE" {
				break done
			} else {
				for _, act := range com {
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

	var errBuff bytes.Buffer
	cmd := exec.Command("sh", "-c", cmdStr)
	cmd.Stderr = &errBuff

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

	if action.Action == "COPY_TO_PROCESS" {
		cmdStr = fmt.Sprintf(action.Command, filepath.Join(conf.Draft, file), filepath.Join(conf.Process, file))
	} else if action.Action == "COPY_TO_SUCCESS" {
		runCommand = doprocess(filepath.Join(conf.Process, file))
		cmdStr = fmt.Sprintf(action.Command, filepath.Join(conf.Process, file), filepath.Join(conf.Success, file))
	} else if action.Action == "COPY_TO_FAIL" {
		cmdStr = fmt.Sprintf(action.Command, filepath.Join(conf.Process, file), filepath.Join(conf.Fail, file))
	}

	if runCommand {
		out, err := runCMD(cmdStr)

		if out != nil {
			log.Printf("%v \n", out)
		}
		if err != nil {
			log.Printf("result: %v %s\n%s", err.Error(), cmdStr, string(out))
		}
		next = action.Success
	} else {
		log.Println("DONE")
		next = action.Fail
	}

	return
}

func preparemasteralarmbrake() (_tkm tk.M) {
	_tkm = tk.M{}

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

	csr, err := workerconn.NewQuery().
		Select("brakeprogram", "alarmname", "alarmindex", "type").
		From("AlarmBrake").
		Cursor(nil)

	if err != nil {
		return
	}

	for {
		_atkm := tk.M{}
		err = csr.Fetch(&_atkm, 1, false)
		if err != nil {
			break
		}

		_tkm.Set(tk.Sprintf("%d", _atkm.GetInt("alarmindex")), _atkm)
	}

	return
}

func doprocess(file string) (success bool) {
	log.Printf("doProcess: %v \n", file)
	success = false
	t1 := time.Now()

	ilines, err := lineCounter(file)
	if err != nil {
		return
	}

	masteralarmbrake := preparemasteralarmbrake()

	sresult := make(chan int, ilines)
	sdata := make(chan string, ilines)
	for i := 0; i < 10; i++ {
		go workersave(i, sdata, sresult, &masteralarmbrake)
	}

	asend := 0
	_file, err := os.Open(file)
	if err != nil {
		log.Fatal(err)
	}
	defer _file.Close()
	scanner := bufio.NewScanner(_file)
	for scanner.Scan() {
		sdata <- tk.Sprintf("%s,%d", scanner.Text(), asend+1)
		asend++
	}
	close(sdata)

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	step := getstep(asend)
	for i := 0; i < asend; i++ {
		<-sresult
		if i%step == 0 {
			tk.Printfn("Done Saved Data %d to %d, in %s",
				i, asend, time.Since(t1).String())
		}
	}
	close(sresult)

	_t1_1 := time.Now()
	//Event Update
	var eventconn dbox.IConnection
	for {
		var err error
		eventconn, err = PrepareConnection()
		if err == nil {
			break
		} else {
			tk.Printfn("==#DB-ERRCONN==\n %s \n", err.Error())
			<-time.After(time.Second * 3)
		}
	}
	defer eventconn.Close()

	ctx := orm.New(eventconn)

	down := econv.NewHFDDownConversion(ctx)
	down.Run()

	tk.Println("Done update event in ", time.Since(_t1_1).String())
	//====================================================================

	_t1_1 = time.Now()
	//Update Monitoring
	UpdateLastMonitoring()
	tk.Println("Done update last monitor in ", time.Since(_t1_1).String())

	return
}

func workersave(wi int, jobs <-chan string, result chan<- int, msalarmbrake *tk.M) {
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

	dtablename := tk.Sprintf("%s", new(EventRawHFD).TableName())

	qSave := workerconn.NewQuery().
		From(dtablename).
		SetConfig("multiexec", true).
		Save()

	split := func(_astr string) (_erh EventRawHFD) {
		_erh = EventRawHFD{}

		_fdata := strings.Split(_astr, ",") //"stime", "project", "param", "id"
		if len(_fdata) < 4 {
			return
		}
		_ddata := strings.Split(_fdata[2], ".")

		// _erh.ProjectName = _fdata[1]
		_erh.ProjectName = "Tejuva"

		_erh.Turbine = _ddata[1]
		_erh.TimeStamp, _ = time.Parse("02-Jan-2006 15:04:05", _fdata[0])
		_erh.DateInfo = helper.GetDateInfo(_erh.TimeStamp)

		_erh.EventType = tk.Sprintf("%s.%s", _ddata[2], _ddata[3])

		ialarmid := "999"
		if _fdata[3] != "" {
			ialarmid = _fdata[3]
		}

		// if !msalarmbrake.Has(ialarmid) {
		// 	ialarmid = "999"
		// }else{

		// }

		_msabrake := msalarmbrake.Get(ialarmid, tk.M{}).(tk.M)

		_erh.BrakeProgram = _msabrake.GetInt("brakeprogram")
		if !msalarmbrake.Has(ialarmid) {
			_erh.BrakeProgram = 999
		}

		_erh.AlarmDescription = _msabrake.GetString("alarmname")
		_erh.AlarmId = tk.ToInt(ialarmid, tk.RoundingAuto)
		_erh.BrakeType = _msabrake.GetString("type")

		_erh.ID = tk.Sprintf("%s#%s#%s#%d#%s#%s", _erh.TimeStamp.Format("20060102_150405.000"),
			_erh.ProjectName, _erh.Turbine, _erh.AlarmId, _erh.EventType, _fdata[4])

		return
	}

	trx := string("")
	for trx = range jobs {
		_erh := split(trx)

		if _erh.ID != "" {
			err := qSave.Exec(tk.M{}.Set("data", _erh))
			if err != nil {
				tk.Println(err)
			}
		}

		result <- 1
	}

	return
}

func lineCounter(_fpath string) (int, error) {
	r, err := os.Open(_fpath)
	if err != nil {
		log.Fatal(err)
	}
	defer r.Close()

	buf := make([]byte, 32*1024)
	count := 0
	lineSep := []byte{'\n'}

	for {
		c, err := r.Read(buf)
		count += bytes.Count(buf[:c], lineSep)

		switch {
		case err == io.EOF:
			return count, nil

		case err != nil:
			return count, err
		}
	}
}

func getstep(count int) int {
	v := count / 5
	if v == 0 {
		return 1
	}
	return v
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

	tk.Println(">>> Delete monitoring before : ", speriode)

	err = workerconn.NewQuery().
		Delete().
		From(new(MonitoringEvent).TableName()).
		Where(dbox.Lte("grouptimestamp", speriode)).
		Exec(nil)

	if err != nil {
		tk.Println(">>> Error found on Delete : ", err.Error())
	}

	msmonitor, mskeys := PrepareMasterMonitoring()
	mseventraw := PrepareEventRawHFD(eperiode)
	tk.Println(">>> periode ", speriode, " ----- ", eperiode)
	//Change to event up down
	xcsr, err := workerconn.NewQuery().
		Select("grouptimestamp", "project", "turbine", "status", "type", "alarmdescription", "alarmid").
		From(new(MonitoringEvent).TableName()).
		Where(dbox.And(dbox.Lte("grouptimestamp", eperiode), dbox.Gt("grouptimestamp", speriode))).
		Order("timestamp").
		Cursor(nil)

	if err != nil {
		return
	}

	// _allkeys := tk.M{}
	for {
		_me := MonitoringEvent{}
		err = xcsr.Fetch(&_me, 1, false)
		if err != nil {
			break
		}

		_key := tk.Sprintf("%s#%s#%s",
			_me.Project,
			_me.Turbine,
			_me.GroupTimeStamp.Format("060102_150405"),
		)
		// tk.Println(">>> me key : ", _key)
		// _allkeys.Set(_key, 1)
		if _mo, _bo := msmonitor[_key]; _bo {
			_mo.Status = "brake"
			if _me.Status == "up" {
				_mo.Status = "ok"
			}

			_mo.Type = _me.Type
			_mo.StatusCode = _me.AlarmId
			_mo.StatusDesc = _me.AlarmDescription

			msmonitor[_key] = _mo
		}
	}
	xcsr.Close()

	sqsave := sworkerconn.NewQuery().
		From(new(Monitoring).TableName()).
		SetConfig("multiexec", true).
		Save()

	sort.Strings(mskeys)
	_lstatus := make(map[string]Monitoring, 0)

	// _ic := 0
	for _, _skey := range mskeys {
		_mo := msmonitor[_skey]

		if _mo.Status == "" || _mo.Status == "N/A" {
			_mo.Status = "N/A"
			_mo.Type = ""
			_mo.StatusCode = 0
			_mo.StatusDesc = ""
			if _erdata, _ercond := mseventraw[_skey]; _ercond { //&& _lsdata.Status == "brake"
				_lsdata := _lstatus[_mo.Turbine]
				// _ = _lsdata
				// === Look brake from previous status
				if _lsdata.Status != "" {
					_mo.Status = _lsdata.Status
					_mo.Type = _lsdata.Type
				}

				_mo.StatusCode = _erdata.AlarmId
				_mo.StatusDesc = _erdata.AlarmDescription
			}
		}

		_astatus := Monitoring{}
		_astatus.Status = _mo.Status
		_astatus.Type = _mo.Type
		_astatus.StatusCode = _mo.StatusCode
		_astatus.StatusDesc = _mo.StatusDesc

		_lstatus[_mo.Turbine] = _astatus
		// }

		_mo.LastUpdate = _nt0
		_mo.LastUpdateDateInfo = helper.GetDateInfo(_nt0)

		_ = sqsave.Exec(tk.M{}.Set("data", _mo))
		// _ic++
	}

	// for _, _mo := range msmonitor {
	// 	if _mo.Status == "" {
	// 		_mo.Status = "N/A"
	// 	}

	// 	_mo.LastUpdate = _nt0
	// 	_mo.LastUpdateDateInfo = helper.GetDateInfo(_nt0)

	// 	_ = sqsave.Exec(tk.M{}.Set("data", _mo))
	// }

	tk.Println(" >>> End Update Last Monitoring in ", time.Since(_nt0).String())
}

func PrepareMasterMonitoring() (_mnt map[string]Monitoring, _arkey []string) {
	_mnt, _arkey = make(map[string]Monitoring), make([]string, 0, 0)

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

		_arkey = append(_arkey, _amnt.ID)
	}

	return
}

func PrepareEventRawHFD(_ltime time.Time) (_mnt map[string]EventRawHFD) {
	_mnt = make(map[string]EventRawHFD)

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

	_stime := _ltime.AddDate(0, 0, -1)
	xcsr, err := workerconn.NewQuery().
		Select().
		From(new(EventRawHFD).TableName()).
		Where(dbox.And(dbox.Lte("timestamp", _ltime), dbox.Gt("timestamp", _stime))).
		Order("timestamp").
		Cursor(nil)

	if err != nil {
		return
	}

	defer xcsr.Close()

	for {
		_aerh := EventRawHFD{}
		err = xcsr.Fetch(&_aerh, 1, false)
		if err != nil {
			break
		}

		GroupTimeStamp := convertTo10min(_aerh.TimeStamp)
		_key := tk.Sprintf("%s#%s#%s",
			_aerh.ProjectName,
			_aerh.Turbine,
			GroupTimeStamp.Format("060102_150405"),
		)

		_mnt[_key] = _aerh
	}

	return
}

func convertTo10min(input time.Time) (output time.Time) {
	// THour := input.Hour()
	TMinute := input.Minute()
	TSecond := input.Second()
	TMinuteValue := float64(TMinute) + tk.Div(float64(TSecond), 60.0)
	TMinuteCategory := tk.ToInt(tk.RoundingUp64(tk.Div(TMinuteValue, 10), 0)*10, "0")

	tmpInput := input.Add(time.Duration(TMinuteCategory-TMinute) * time.Minute).Add(time.Duration(TSecond*-1) * time.Second).UTC()
	output, _ = time.Parse("20060102_150405", tmpInput.Format("20060102_150405"))
	return
}
