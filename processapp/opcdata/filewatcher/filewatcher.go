package filewatcher

import (
	d "eaciit/wfdemo/processapp/opcdata/datareader"
	tk "github.com/eaciit/toolkit"
	"github.com/howeyc/fsnotify"
	"strings"
	_ "time"
)

type FileWatcher struct {
	PathSources string
	PathProcess string
	PathRoot    string
}

func NewFileWatcher(pathSources string, pathProcess string, pathRoot string) *FileWatcher {
	fw := new(FileWatcher)
	fw.PathSources = pathSources
	fw.PathProcess = pathProcess
	fw.PathRoot = pathRoot

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
	if (ev.IsCreate() || ev.IsModify()) && strings.Contains(ev.Name, "DataFile") && !strings.Contains(ev.Name, "DataFile-T") {
		d.NewDataReader(ev.Name, w.PathProcess, w.PathRoot).Start()
		// tk.Println(ev.Name, time.Now())
	}
}
