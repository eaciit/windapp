package main

import (
	. "eaciit/wfdemo-git/library/helper"
	"flag"
	"github.com/eaciit/dbox"
	_ "github.com/eaciit/dbox/dbc/mongo"
	tk "github.com/eaciit/toolkit"
	"runtime"
	"strings"
	"sync"
	"time"
)

var log *tk.LogEngine

const (
	sError   = "ERROR"
	sInfo    = "INFO"
	sWarning = "WARNING"
)

func main() {
	logpath := ""
	flag.StringVar(&logpath, "log", "", "Log folder place")
	flag.Parse()
	config := ReadConfig()
	if logpath == "" {
		logpath, _ = config["logpath"]
	}
	log, _ = tk.NewLog(false, true, logpath, "simpleHFDLog_%s", "20060102")
	ctx, e := PrepareConnection(config)
	if e != nil {
		log.AddLog(e.Error(), sError)
	}

	csrTag, e := ctx.NewQuery().From("ref_databrowsertag").
		Select("realtimefield").
		Where(dbox.And(
			dbox.Eq("source", "ScadaDataHFD"),
			dbox.Eq("enable", true)),
		).
		Cursor(nil)
	if e != nil {
		log.AddLog(e.Error(), sError)
	}
	defer csrTag.Close()
	tagList := []string{"_id", "timestamp", "dateinfo", "projectname", "turbine", "turbinestate"}
	tags := tk.M{}
	for {
		tags = tk.M{}
		e = csrTag.Fetch(&tags, 1, false)
		if e != nil {
			break
		}
		tagList = append(tagList, strings.ToLower(tags.GetString("realtimefield")))
	}

	csrLog, e := ctx.NewQuery().From("log_latestdaterun").
		Where(dbox.Eq("type", "databrowser")).Cursor(nil)
	if e != nil {
		log.AddLog(e.Error(), sError)
	}
	defer csrLog.Close()
	lastData := []struct {
		ProjectName string
		LastDate    time.Time
	}{}
	e = csrLog.Fetch(&lastData, 0, false)
	if e != nil {
		log.AddLog(e.Error(), sError)
	}
	lastDatePerProject := map[string]time.Time{}
	for _, val := range lastData {
		lastDatePerProject[val.ProjectName] = val.LastDate
	}

	csrProject, e := ctx.NewQuery().From("ref_project").
		Where(dbox.Eq("active", true)).Cursor(nil)
	if e != nil {
		log.AddLog(e.Error(), sError)
	}
	defer csrProject.Close()
	projectList := []struct {
		ProjectID   string
		ProjectName string
	}{}
	e = csrProject.Fetch(&projectList, 0, false)
	if e != nil {
		log.AddLog(e.Error(), sError)
	}
	var wgProject sync.WaitGroup
	wgProject.Add(len(projectList))

	for _, project := range projectList {
		go func(projectid string) {
			csrData, e := ctx.NewQuery().From("Scada10MinHFD").Select(tagList...).
				Where(dbox.And(
					dbox.Eq("projectname", projectid),
					dbox.Gte("timestamp", lastDatePerProject[projectid]),
					dbox.Eq("isnull", false))).
				Cursor(nil)
			if e != nil {
				log.AddLog(e.Error(), sError)
			}
			defer csrData.Close()

			maxTimeStamp := time.Time{}

			var wg sync.WaitGroup
			totalData := csrData.Count()
			totalWorker := runtime.NumCPU() * 2
			chanData := make(chan tk.M, totalData)
			step := getstep(totalData)
			tNow := time.Now()

			wg.Add(totalWorker)
			for i := 0; i < totalWorker; i++ {
				go func() {
					ctxWorker, e := PrepareConnection(config)
					if e != nil {
						log.AddLog(e.Error(), sError)
					}
					csrSave := ctxWorker.NewQuery().From("DatabrowserHFD").SetConfig("multiexec", true).Save()
					defer csrSave.Close()
					for data := range chanData {
						if data.GetInt("count")%step == 0 {
							percent := tk.ToInt(tk.Div(float64(data.GetInt("count"))*100.0, float64(totalData)), tk.RoundingUp)
							log.AddLog(tk.Sprintf("[%s] Saving %d of %d (%d percent) in %s\n",
								strings.ToUpper(projectid), data.GetInt("count"), totalData, percent,
								time.Since(tNow).String()), sInfo)
						}
						data.Unset("count")
						csrSave.Exec(tk.M{"data": data})
					}
					wg.Done()
				}()
			}

			log.AddLog(tk.Sprintf("Processing %d data [%s] with %d step using %d CPU since %s",
				totalData, strings.ToUpper(projectid), step, totalWorker, lastDatePerProject[projectid].Format("20060102_150405")), sInfo)

			count := 0
			_data := tk.M{}
			currTimeStamp := time.Time{}
			for {
				count++
				_data = tk.M{}
				e = csrData.Fetch(&_data, 1, false)
				if e != nil {
					if !strings.Contains(e.Error(), "Not found") {
						log.AddLog(e.Error(), sError)
					}
					break
				}
				currTimeStamp = _data.Get("timestamp", time.Time{}).(time.Time).UTC()
				if currTimeStamp.After(maxTimeStamp) {
					maxTimeStamp = currTimeStamp
				}

				_data.Set("count", count)
				chanData <- _data

				// if count%step == 0 {
				// 	log.AddLog(tk.Sprintf("Processing %d of %d in %s\n",
				// 		count, totalData,
				// 		time.Since(tNow).String()), sInfo)
				// }
			}
			close(chanData)
			wg.Wait()

			if maxTimeStamp.Year() > 1 {
				e = ctx.NewQuery().From("log_latestdaterun").Save().
					Exec(tk.M{"data": tk.M{
						"_id":         "databrowser_hfd_" + projectid,
						"lastdate":    maxTimeStamp,
						"projectname": projectid,
						"type":        "databrowser",
					}})
				if e != nil {
					log.AddLog(e.Error(), sError)
				}
			}
			wgProject.Done()
		}(project.ProjectID)
	}
	wgProject.Wait()
}

func getstep(count int) int {
	v := count / 5
	if v == 0 {
		return 1
	}
	return v
}

func PrepareConnection(config map[string]string) (dbox.IConnection, error) {
	ci := &dbox.ConnectionInfo{config["host"], config["database"], config["username"], config["password"], nil}
	c, e := dbox.NewConnection("mongo", ci)

	if e != nil {
		return nil, e
	}

	e = c.Connect()
	if e != nil {
		return nil, e
	}
	// log.AddLog(tk.Sprintf("DB Connect %s : %s", config["host"], config["database"]), sInfo)
	return c, nil
}
