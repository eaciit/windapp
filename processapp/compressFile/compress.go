package main

import (
	lh "eaciit/wfdemo-git/library/helper"
	"flag"
	tk "github.com/eaciit/toolkit"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

const (
	sError   = "ERROR"
	sInfo    = "INFO"
	sWarning = "WARNING"
)

var (
	wd = func() string {
		d, _ := os.Getwd()
		return d
	}()

	Log *tk.LogEngine
)

func main() {
	var tipe string
	flag.StringVar(&tipe, "tipe", "data", "to determine whether compressing data folder or rawdata folder")
	flag.Parse()

	Log, _ := tk.NewLog(false, true, wd, "compress_%s", "20060102")
	t0 := time.Now()
	config := lh.ReadConfig()
	dataPath := ""
	destDataPath := ""
	if tipe == "data" {
		dataPath = config["datapath"]
		destDataPath = config["destdatapath"]
		Log.AddLog(tk.Sprintf("starting compressing file on %s folder\n", "data"), sInfo)
	} else {
		dataPath = config["rawdatapath"]
		destDataPath = config["destrawdatapath"]
		Log.AddLog(tk.Sprintf("starting compressing file on %s folder\n", "rawdata"), sInfo)
	}
	if _, err := os.Stat(destDataPath); os.IsNotExist(err) {
		os.Mkdir(destDataPath, 0777) /*jika folder destinasi belum ada maka dibuat*/
	}
	platform := runtime.GOOS
	sourcePathLevel1 := ""
	sourcePathLevel2 := ""
	destPathLevel1 := ""
	destPathLevel2 := ""
	dirLevel1, e := ioutil.ReadDir(dataPath) /*list nama folder per project*/
	deletedList := []string{}
	compressedList := []string{}
	extension := ""
	if e != nil {
		Log.AddLog(e.Error(), sError)
	}
	for _, f := range dirLevel1 {
		sourcePathLevel1 = filepath.Join(dataPath, f.Name())
		destPathLevel1 = filepath.Join(destDataPath, f.Name())
		if _, err := os.Stat(destPathLevel1); os.IsNotExist(err) {
			os.Mkdir(destPathLevel1, 0777) /*jika destinasi belum ada folder project maka dibuat*/
		}
		dirLevel2, e := ioutil.ReadDir(sourcePathLevel1) /*list nama folder per hari*/
		if e != nil {
			Log.AddLog(e.Error(), sError)
		}
		for _, f2 := range dirLevel2 {
			sourcePathLevel2 = filepath.Join(sourcePathLevel1, f2.Name())
			destPathLevel2 = filepath.Join(destPathLevel1, f2.Name())
			if platform == "windows" {
				extension = ".zip"
				destPathLevel2 += extension
				e = tk.ZipCompress(sourcePathLevel2, destPathLevel2)
			} else {
				extension = ".tar.gz"
				destPathLevel2 += extension
				e = tk.TarCompress(sourcePathLevel2, destPathLevel2)
			}
			if e != nil {
				Log.AddLog(e.Error(), sError)
			}
			Log.AddLog(tk.Sprintf("%s%s created", filepath.Join(f.Name(), f2.Name()), extension), sInfo)
			deletedList = append(deletedList, sourcePathLevel2)
			compressedList = append(compressedList, destPathLevel2)
		}
	}
	for idx, delPath := range deletedList {
		_, err := os.Stat(compressedList[idx])
		if err == nil {
			/*jika file compress sudah ada maka delete folder source*/
			e = os.RemoveAll(delPath)
			if e != nil {
				Log.AddLog(e.Error(), sError)
			}
			if idx == len(deletedList)-1 {
				Log.AddLog(tk.Sprintf("folder %s deleted\n", delPath), sInfo)
			} else {
				Log.AddLog(tk.Sprintf("folder %s deleted", delPath), sInfo)
			}
		}
	}
	Log.AddLog(tk.Sprintf("finish compressing file in %s minutes\n=======================================================================",
		tk.ToString(time.Since(t0).Minutes())), sInfo)
}
