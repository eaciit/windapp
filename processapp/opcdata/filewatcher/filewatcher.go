package filewatcher

import (
	d "eaciit/wfdemo-git/processapp/opcdata/datareader"
	tk "github.com/eaciit/toolkit"
	"github.com/howeyc/fsnotify"
	"strings"
	"fmt"
	"time"
)

type FileWatcher struct {
	PathSources string
	PathProcess string
	PathRoot    string
	PathUpload  string
}

func NewFileWatcher(pathSources string, pathProcess string, pathRoot string,uploadDir string) *FileWatcher {
	fw := new(FileWatcher)
	fw.PathSources = pathSources
	fw.PathProcess = pathProcess
	fw.PathRoot = pathRoot
	fw.PathUpload = uploadDir
	return fw
}

func (w *FileWatcher) StartWatcher() {
	tk.Println("Starting file watcher...")

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		tk.Println(err.Error())
	}

	done := make(chan bool)

	go func() {
		for {
			select {
			case ev := <-watcher.Event:
				w.doAction(ev)
			case err := <-watcher.Error:
				tk.Println("error:", err)
			}
		}
	}()

	err = watcher.Watch(w.PathSources)
	if err != nil {
		tk.Println(err.Error())
	}

	<-done

	watcher.Close()
}

func (w *FileWatcher) doAction(ev *fsnotify.FileEvent) {
	if ev.IsCreate() && strings.Contains(ev.Name, "DataFile") && !strings.Contains(ev.Name, "DataFile-T") {
		fileName := strings.Split(ev.Name,"\\")
		
		now := time.Now()
		year,month,day:=now.Date()
		hour,_,_:=now.Clock()
		
		targetFileName := fmt.Sprintf("DataFile%d%02d%02d-%02d.csv",year,month,day,hour)
		fmt.Println(targetFileName)
		tk.Println(fileName[len(fileName)-1],targetFileName)
		if fileName[len(fileName)-1]==targetFileName{
			d.NewDataReader(ev.Name, w.PathProcess, w.PathRoot,w.PathUpload).Start()
			//tk.Println(ev.Name, time.Now())
		}
		
	}
}
