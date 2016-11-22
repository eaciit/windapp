package main

import (
	"os"
	"time"

	// "strings"

	"github.com/eaciit/dbox"
	_ "github.com/eaciit/dbox/dbc/mongo"

	tk "github.com/eaciit/toolkit"
	"io/ioutil"
	// "path/filepath"
	"regexp"
)

var (
	checker    tk.M
	listfile   tk.M
	dir        = string("E:\\AlgoEngines")
	tablename  = string("ScadaThreeSecs")
	dtablename = string("filecheck")
)

func PrepareDatabaseConnection() (dbox.IConnection, error) {
	var config = tk.M{}.Set("timeout", 10)

	ci := &dbox.ConnectionInfo{"192.168.0.220:27017", "ecwfdemo", "", "", config}
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

func getstep(count int) int {
	v := count / 100
	if v == 0 {
		return 1
	}
	return v
}

func PrepareDataChecker() {
	checker = tk.M{}

	dbconn, _ := PrepareDatabaseConnection()
	defer dbconn.Close()
	csr, _ := dbconn.NewQuery().Select().From("check_count").Cursor(nil)
	defer csr.Close()

	for {
		_tkm := tk.M{}
		err := csr.Fetch(&_tkm, 1, false)
		if err != nil {
			break
		}

		checker.Set(_tkm.GetString("_id"), _tkm.GetInt("count"))
	}

	return
}

func PrepareFileChecker() {
	listfile = tk.M{}

	files, e := ioutil.ReadDir(dir)
	if e != nil {
		tk.Println(e)
		os.Exit(0)
	}

	// scount := len(files)
	icount := 0
	// step := getstep(scount)

	for _, file := range files {
		icount++
		tkm := tk.M{}
		filename := file.Name()
		if cond, _ := regexp.MatchString("^(DataFile20.*)(\\.[Cc][Ss][Vv])$", filename); cond {
			tkm.Set("FileName", filename).Set("_id", filename).Set("CountData", 0)
			if !checker.Has(filename) {
				listfile.Set(filename, file.Size())
			}
		}
	}
}

func main() {

	t0 := time.Now()

	// stime := time.Date(2015, 11, 1, 0, 0, 0, 0, time.UTC)
	// etime := time.Date(2016, 1, 1, 0, 0, 0, 0, time.UTC)

	tk.Println("Prepare file for checker,")
	PrepareFileChecker()

	scount := len(listfile)

	sresult := make(chan tk.M, scount)
	sdata := make(chan string, scount)
	for i := 0; i < 10; i++ {
		go getaggrdata(i, sdata, sresult)
	}

	for key, _ := range listfile {
		// tk.Printfn("%s : %#v", key, val)
		sdata <- key
	}

	close(sdata)

	sconn, _ := PrepareDatabaseConnection()
	defer sconn.Close()

	qSave := sconn.NewQuery().
		From(dtablename).
		SetConfig("multiexec", true).
		Save()

	_astr := tk.M{}
	for i := 0; i < scount; i++ {
		tkm := <-sresult
		_id := tkm.GetString("_id")
		_astr.Set(_id, 1)

		if listfile.Has(_id) {
			tkm.Set("filesize", listfile[_id])
		}

		err := qSave.Exec(tk.M{}.Set("data", tkm))
		if err != nil {
			tk.Println(err)
		}
	}

	for key, val := range listfile {
		if !_astr.Has(key) {

			tkm := tk.M{}.Set("_id", key).
				Set("filesize", val).
				Set("count", 0)

			err := qSave.Exec(tk.M{}.Set("data", tkm))
			if err != nil {
				tk.Println(err)
			}
		}
	}

	tk.Printf("All data check done in %s \n",
		time.Since(t0).String())
}

func getaggrdata(wi int, jobs <-chan string, result chan<- tk.M) {

	conn, err := PrepareDatabaseConnection()
	if err != nil {
		tk.Println(err)
		return
	}
	defer conn.Close()

	for _file := range jobs {

		pipes := []tk.M{}
		group := tk.M{}.Set("_id", "$file").
			Set("count", tk.M{}.Set("$sum", 1)).
			Set("min", tk.M{}.Set("$min", "$timestampconverted")).
			Set("max", tk.M{}.Set("$max", "$timestampconverted"))

		// sort := tk.M{}.Set("_id", 1)

		pipes = append(pipes, tk.M{"$match": tk.M{}.Set("file", _file)})
		pipes = append(pipes, tk.M{"$group": group})

		csr, e := conn.NewQuery().
			From(tablename).
			Command("pipe", pipes).
			Cursor(nil)

		if e != nil {
			tk.Printf("ERRROR: %v \n", e.Error())
			os.Exit(1)
		}

		tkm := tk.M{}
		_ = csr.Fetch(&tkm, 1, false)
		csr.Close()

		result <- tkm
	}

	return
}
