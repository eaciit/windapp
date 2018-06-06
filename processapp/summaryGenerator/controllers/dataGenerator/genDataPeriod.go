package generatorControllers

import (
	. "eaciit/wfdemo-git/library/models"
	. "eaciit/wfdemo-git/processapp/summaryGenerator/controllers"
	"eaciit/wfdemo-git/web/helper"
	"time"

	"github.com/eaciit/dbox"
	_ "github.com/eaciit/dbox/dbc/mongo"
	"github.com/eaciit/toolkit"
)

type GenDataPeriod struct {
	*BaseController
}

func (d *GenDataPeriod) Generate(base *BaseController) {
	d.BaseController = base
	// conn, e := PrepareConnection()
	// if e != nil {
	// 	toolkit.Println("Scada Summary : " + e.Error())
	// 	os.Exit(0)
	// }
	// defer conn.Close()
	var e error
	conn := d.BaseController.Ctx.Connection

	// reset data
	// #faisal
	// remove this function
	// d.BaseController.Ctx.DeleteMany(NewLatestDataPeriod(), dbox.Ne("projectname", ""))

	projects, _ := helper.GetProjectList()

	for _, proj := range projects {
		projectName := proj.Value
		scadaResults := make([]time.Time, 2)
		dgrResults := make([]time.Time, 2)
		alarmResults := make([]time.Time, 2)
		jmrResults := make([]time.Time, 2)
		metResults := make([]time.Time, 2)
		durationResults := make([]time.Time, 2)
		scadaHFDResult := make([]time.Time, 2)
		warningResult := make([]time.Time, 2)
		eventDownResult := make([]time.Time, 2)
		eventRawResult := make([]time.Time, 2)
		eventDownHFDResult := make([]time.Time, 2)
		scadaOEMResult := make([]time.Time, 2)
		// scadaAnomalyresults := make([]time.Time, 2)
		// alarmOverlappingresults := make([]time.Time, 2)
		// alarmScadaAnomalyresults := make([]time.Time, 2)

		scadaResults[0], scadaResults[1], e = getDataDateAvailable(conn, new(ScadaData).TableName(), "timestamp", dbox.Eq("projectname", projectName))
		dgrResults[0], dgrResults[1], e = getDataDateAvailable(conn, new(DGRModel).TableName(), "dateinfo.dateid", dbox.Eq("site", projectName))
		alarmResults[0], alarmResults[1], e = getDataDateAvailable(conn, new(Alarm).TableName(), "startdate", dbox.Eq("farm", projectName))
		jmrResults[0], jmrResults[1], e = getDataDateAvailable(conn, new(JMR).TableName(), "dateinfo.dateid", nil)
		metResults[0], metResults[1], e = getDataDateAvailable(conn, new(MetTower).TableName(), "timestamp", nil)
		durationResults[0], durationResults[1], e = getDataDateAvailable(conn, new(ScadaData).TableName(), "timestamp", dbox.And(dbox.Eq("isvalidtimeduration", false), dbox.Eq("projectname", projectName)))
		scadaHFDResult[0], scadaHFDResult[1], e = getDataDateAvailable(conn, new(ScadaDataHFD).TableName(), "timestamp", dbox.Eq("projectname", projectName))
		warningResult[0], warningResult[1], e = getDataDateAvailable(conn, new(EventAlarm).TableName(), "timestart", dbox.Eq("projectname", projectName))
		eventDownResult[0], eventDownResult[1], e = getDataDateAvailable(conn, new(EventDown).TableName(), "timestart", dbox.Eq("projectname", projectName))
		eventRawResult[0], eventRawResult[1], e = getDataDateAvailable(conn, new(EventRaw).TableName(), "timestamp", dbox.Eq("projectname", projectName))
		eventDownHFDResult[0], eventDownHFDResult[1], e = getDataDateAvailable(conn, new(EventDownHFD).TableName(), "timestart", dbox.Eq("projectname", projectName))
		scadaOEMResult[0], scadaOEMResult[1], e = getDataDateAvailable(conn, new(ScadaDataOEM).TableName(), "timestamp", dbox.Eq("projectname", projectName))
		// alarmOverlappingresults[0], alarmOverlappingresults[1], e = getDataDateAvailable(conn, new(AlarmOverlapping).TableName(), "startdate", dbox.Eq("farm", projectName))
		// alarmScadaAnomalyresults[0], alarmScadaAnomalyresults[1], e = getDataDateAvailable(conn, new(AlarmScadaAnomaly).TableName(), "startdate", dbox.Eq("farm", projectName))
		// scadaAnomalyresults[0], scadaAnomalyresults[1], e = getDataDateAvailable(conn, new(ScadaData).TableName(), "timestamp", dbox.And(dbox.Eq("isvalidtimeduration", true), dbox.Eq("projectname", projectName)))
		_ = e
		availdatedata := struct {
			ScadaData    []time.Time
			DGRData      []time.Time
			Alarm        []time.Time
			JMR          []time.Time
			MET          []time.Time
			Duration     []time.Time
			ScadaDataHFD []time.Time
			Warning      []time.Time
			ScadaDataOEM []time.Time
			EventDown    []time.Time
			EventRaw     []time.Time
			EventDownHFD []time.Time
			// ScadaAnomaly      []time.Time
			// AlarmOverlapping  []time.Time
			// AlarmScadaAnomaly []time.Time
		}{
			ScadaData:    scadaResults,
			DGRData:      dgrResults,
			Alarm:        alarmResults,
			JMR:          jmrResults,
			MET:          metResults,
			Duration:     durationResults,
			ScadaDataHFD: scadaHFDResult,
			Warning:      warningResult,
			ScadaDataOEM: scadaOEMResult,
			EventDown:    eventDownResult,
			EventRaw:     eventRawResult,
			EventDownHFD: eventDownHFDResult,
			// ScadaAnomaly:      scadaAnomalyresults,
			// AlarmOverlapping:  alarmOverlappingresults,
			// AlarmScadaAnomaly: alarmScadaAnomalyresults,
		}

		mdl := new(LatestDataPeriod)

		mdl.Type = "ScadaData"
		mdl.ProjectName = projectName
		mdl.Data = availdatedata.ScadaData
		mdl = mdl.New()

		d.BaseController.Ctx.Save(mdl)

		mdl = new(LatestDataPeriod)
		mdl.Type = "DGRData"
		mdl.ProjectName = projectName
		mdl.Data = availdatedata.DGRData
		mdl = mdl.New()

		d.BaseController.Ctx.Save(mdl)

		mdl = new(LatestDataPeriod)
		mdl.Type = "Alarm"
		mdl.ProjectName = projectName
		mdl.Data = availdatedata.Alarm
		mdl = mdl.New()

		d.BaseController.Ctx.Save(mdl)

		mdl = new(LatestDataPeriod)
		mdl.Type = "JMR"
		mdl.ProjectName = projectName
		mdl.Data = availdatedata.JMR
		mdl = mdl.New()

		d.BaseController.Ctx.Save(mdl)

		mdl = new(LatestDataPeriod)
		mdl.Type = "MET"
		mdl.ProjectName = projectName
		mdl.Data = availdatedata.MET
		mdl = mdl.New()

		d.BaseController.Ctx.Save(mdl)

		mdl = new(LatestDataPeriod)
		mdl.Type = "Duration"
		mdl.ProjectName = projectName
		mdl.Data = availdatedata.Duration
		mdl = mdl.New()

		d.BaseController.Ctx.Save(mdl)

		mdl = new(LatestDataPeriod)
		mdl.Type = "ScadaDataHFD"
		mdl.ProjectName = projectName
		mdl.Data = availdatedata.ScadaDataHFD
		mdl = mdl.New()

		d.BaseController.Ctx.Save(mdl)

		mdl = new(LatestDataPeriod)
		mdl.Type = "Warning"
		mdl.ProjectName = projectName
		mdl.Data = availdatedata.Warning
		mdl = mdl.New()

		d.BaseController.Ctx.Save(mdl)

		mdl = new(LatestDataPeriod)
		mdl.Type = "ScadaDataOEM"
		mdl.ProjectName = projectName
		mdl.Data = availdatedata.ScadaDataOEM
		mdl = mdl.New()

		d.BaseController.Ctx.Save(mdl)

		mdl = new(LatestDataPeriod)
		mdl.Type = "EventDown"
		mdl.ProjectName = projectName
		mdl.Data = availdatedata.EventDown
		mdl = mdl.New()

		d.BaseController.Ctx.Save(mdl)

		mdl = new(LatestDataPeriod)
		mdl.Type = "EventRaw"
		mdl.ProjectName = projectName
		mdl.Data = availdatedata.EventRaw
		mdl = mdl.New()

		d.BaseController.Ctx.Save(mdl)

		mdl = new(LatestDataPeriod)
		mdl.Type = "EventDownHFD"
		mdl.ProjectName = projectName
		mdl.Data = availdatedata.EventDownHFD
		mdl = mdl.New()

		d.BaseController.Ctx.Save(mdl)

		// mdl = new(LatestDataPeriod)
		// mdl.Type = "ScadaAnomaly"
		// mdl.ProjectName = projectName
		// mdl.Data = availdatedata.ScadaAnomaly
		// mdl = mdl.New()

		// d.BaseController.Ctx.Save(mdl)

		// mdl = new(LatestDataPeriod)
		// mdl.Type = "AlarmOverlapping"
		// mdl.ProjectName = projectName
		// mdl.Data = availdatedata.AlarmOverlapping
		// mdl = mdl.New()

		// d.BaseController.Ctx.Save(mdl)

		// mdl = new(LatestDataPeriod)
		// mdl.Type = "AlarmScadaAnomaly"
		// mdl.ProjectName = projectName
		// mdl.Data = availdatedata.AlarmScadaAnomaly
		// mdl = mdl.New()

		// d.BaseController.Ctx.Save(mdl)
	}

}

func (d *GenDataPeriod) GenerateMinify(base *BaseController) {
	d.BaseController = base
	// conn, e := PrepareConnection()
	// if e != nil {
	// 	toolkit.Println("Scada Summary : " + e.Error())
	// 	os.Exit(0)
	// }
	t0 := time.Now()
	toolkit.Println("Start generating data available date : ", t0)

	var e error
	// conn := &d.BaseController.Ctx.Connection
	projects, _ := helper.GetProjectList()

	for _, proj := range projects {
		toolkit.Println("Start : ", proj.Name, " - ", t0)
		projectName := proj.Value
		scadaResults := make([]time.Time, 2)
		alarmResults := make([]time.Time, 2)
		durationResults := make([]time.Time, 2)

		scadaResults[0], scadaResults[1], e = getDataDateAvailable(d.BaseController.Ctx.Connection, new(ScadaData).TableName(), "timestamp", dbox.Eq("projectname", projectName))
		alarmResults[0], alarmResults[1], e = getDataDateAvailable(d.BaseController.Ctx.Connection, new(Alarm).TableName(), "startdate", dbox.Eq("farm", projectName))
		durationResults[0], durationResults[1], e = getDataDateAvailable(d.BaseController.Ctx.Connection, new(ScadaData).TableName(), "timestamp", dbox.And(dbox.Eq("isvalidtimeduration", false), dbox.Eq("projectname", projectName)))
		_ = e
		availdatedata := struct {
			ScadaData []time.Time
			Alarm     []time.Time
			Duration  []time.Time
		}{
			ScadaData: scadaResults,
			Alarm:     alarmResults,
			Duration:  durationResults,
		}

		mdl := new(LatestDataPeriod)

		mdl.Type = "ScadaData"
		mdl.ProjectName = projectName
		mdl.Data = availdatedata.ScadaData
		mdl = mdl.New()

		d.BaseController.Ctx.Save(mdl)

		mdl = new(LatestDataPeriod)
		mdl.Type = "Alarm"
		mdl.ProjectName = projectName
		mdl.Data = availdatedata.Alarm
		mdl = mdl.New()

		d.BaseController.Ctx.Save(mdl)

		mdl = new(LatestDataPeriod)
		mdl.Type = "Duration"
		mdl.ProjectName = projectName
		mdl.Data = availdatedata.Duration
		mdl = mdl.New()

		d.BaseController.Ctx.Save(mdl)

	}
	toolkit.Println("End generating data available date in ", time.Since(t0).String())

}

func getDataDateAvailable(conn dbox.IConnection, collectionName string, timestampColumn string, where *dbox.Filter) (min time.Time, max time.Time, err error) {
	q := conn.NewQuery().From(collectionName)

	if where != nil {
		q.Where(where)
	}

	csr, err := q.
		Aggr(dbox.AggrMin, "$"+timestampColumn, "min").
		Aggr(dbox.AggrMax, "$"+timestampColumn, "max").
		Group("enable").
		Cursor(nil)

	defer csr.Close()

	if err != nil {
		return
	}

	data := []toolkit.M{}
	err = csr.Fetch(&data, 0, false)

	if err != nil || len(data) == 0 {
		return
	}

	min = data[0].Get("min", time.Time{}).(time.Time).UTC()
	max = data[0].Get("max", time.Time{}).(time.Time).UTC()

	return
}
