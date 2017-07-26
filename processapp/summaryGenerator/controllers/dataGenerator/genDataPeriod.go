package generatorControllers

import (
	. "eaciit/wfdemo-git/library/models"
	. "eaciit/wfdemo-git/processapp/summaryGenerator/controllers"
	"eaciit/wfdemo-git/web/helper"
	"os"
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
	conn, e := PrepareConnection()
	if e != nil {
		toolkit.Println("Scada Summary : " + e.Error())
		os.Exit(0)
	}

	// reset data
	// #faisal
	// remove this function
	d.BaseController.Ctx.DeleteMany(NewLatestDataPeriod(), dbox.Ne("projectname", ""))

	projects, _ := helper.GetProjectList()

	for _, proj := range projects {
		projectName := proj.Value
		scadaResults := make([]time.Time, 2)
		dgrResults := make([]time.Time, 2)
		alarmResults := make([]time.Time, 2)
		jmrResults := make([]time.Time, 2)
		metResults := make([]time.Time, 2)
		durationResults := make([]time.Time, 2)
		scadaAnomalyresults := make([]time.Time, 2)
		alarmOverlappingresults := make([]time.Time, 2)
		alarmScadaAnomalyresults := make([]time.Time, 2)
		scadaHFDResult := make([]time.Time, 2)
		warningResult := make([]time.Time, 2)
		scadaOEMResult := make([]time.Time, 2)

		scadaResults[0], scadaResults[1], e = getDataDateAvailable(conn, new(ScadaData).TableName(), "timestamp", dbox.Eq("projectname", projectName))
		dgrResults[0], dgrResults[1], e = getDataDateAvailable(conn, new(DGRModel).TableName(), "dateinfo.dateid", dbox.Eq("site", projectName))
		alarmResults[0], alarmResults[1], e = getDataDateAvailable(conn, new(Alarm).TableName(), "startdate", dbox.Eq("farm", projectName))
		alarmOverlappingresults[0], alarmOverlappingresults[1], e = getDataDateAvailable(conn, new(AlarmOverlapping).TableName(), "startdate", dbox.Eq("farm", projectName))
		alarmScadaAnomalyresults[0], alarmScadaAnomalyresults[1], e = getDataDateAvailable(conn, new(AlarmScadaAnomaly).TableName(), "startdate", dbox.Eq("farm", projectName))
		jmrResults[0], jmrResults[1], e = getDataDateAvailable(conn, new(ScadaData).TableName(), "dateinfo.dateid", dbox.Eq("projectname", projectName))
		metResults[0], metResults[1], e = getDataDateAvailable(conn, new(MetTower).TableName(), "timestamp", nil)
		durationResults[0], durationResults[1], e = getDataDateAvailable(conn, new(ScadaData).TableName(), "timestamp", dbox.And(dbox.Eq("isvalidtimeduration", false), dbox.Eq("projectname", projectName)))
		scadaAnomalyresults[0], scadaAnomalyresults[1], e = getDataDateAvailable(conn, new(ScadaData).TableName(), "timestamp", dbox.And(dbox.Eq("isvalidtimeduration", true), dbox.Eq("projectname", projectName)))
		scadaHFDResult[0], scadaHFDResult[1], e = getDataDateAvailable(conn, new(ScadaDataHFD).TableName(), "timestamp", dbox.Eq("projectname", projectName))
		warningResult[0], warningResult[1], e = getDataDateAvailable(conn, new(EventAlarm).TableName(), "timestart", dbox.Eq("projectname", projectName))
		scadaOEMResult[0], scadaOEMResult[1], e = getDataDateAvailable(conn, new(ScadaDataOEM).TableName(), "timestamp", dbox.Eq("projectname", projectName))

		availdatedata := struct {
			ScadaData         []time.Time
			DGRData           []time.Time
			Alarm             []time.Time
			JMR               []time.Time
			MET               []time.Time
			Duration          []time.Time
			ScadaAnomaly      []time.Time
			AlarmOverlapping  []time.Time
			AlarmScadaAnomaly []time.Time
			ScadaDataHFD      []time.Time
			Warning           []time.Time
			ScadaDataOEM      []time.Time
		}{
			ScadaData:         scadaResults,
			DGRData:           dgrResults,
			Alarm:             alarmResults,
			JMR:               jmrResults,
			MET:               metResults,
			Duration:          durationResults,
			ScadaAnomaly:      scadaAnomalyresults,
			AlarmOverlapping:  alarmOverlappingresults,
			AlarmScadaAnomaly: alarmScadaAnomalyresults,
			ScadaDataHFD:      scadaHFDResult,
			Warning:           warningResult,
			ScadaDataOEM:      scadaOEMResult,
		}

		mdl := NewLatestDataPeriod()
		mdl.Type = "ScadaData"
		mdl.ProjectName = projectName
		mdl.Data = availdatedata.ScadaData

		d.BaseController.Ctx.Insert(mdl)

		mdl = NewLatestDataPeriod()
		mdl.Type = "DGRData"
		mdl.ProjectName = projectName
		mdl.Data = availdatedata.DGRData

		d.BaseController.Ctx.Insert(mdl)

		mdl = NewLatestDataPeriod()
		mdl.Type = "Alarm"
		mdl.ProjectName = projectName
		mdl.Data = availdatedata.Alarm

		d.BaseController.Ctx.Insert(mdl)

		mdl = NewLatestDataPeriod()
		mdl.Type = "JMR"
		mdl.ProjectName = projectName
		mdl.Data = availdatedata.JMR

		d.BaseController.Ctx.Insert(mdl)

		mdl = NewLatestDataPeriod()
		mdl.Type = "MET"
		mdl.ProjectName = projectName
		mdl.Data = availdatedata.MET

		d.BaseController.Ctx.Insert(mdl)

		mdl = NewLatestDataPeriod()
		mdl.Type = "Duration"
		mdl.ProjectName = projectName
		mdl.Data = availdatedata.Duration

		d.BaseController.Ctx.Insert(mdl)

		mdl = NewLatestDataPeriod()
		mdl.Type = "ScadaAnomaly"
		mdl.ProjectName = projectName
		mdl.Data = availdatedata.ScadaAnomaly

		d.BaseController.Ctx.Insert(mdl)

		mdl = NewLatestDataPeriod()
		mdl.Type = "AlarmOverlapping"
		mdl.ProjectName = projectName
		mdl.Data = availdatedata.AlarmOverlapping

		d.BaseController.Ctx.Insert(mdl)

		mdl = NewLatestDataPeriod()
		mdl.Type = "AlarmScadaAnomaly"
		mdl.ProjectName = projectName
		mdl.Data = availdatedata.AlarmScadaAnomaly

		d.BaseController.Ctx.Insert(mdl)

		mdl = NewLatestDataPeriod()
		mdl.Type = "ScadaDataHFD"
		mdl.ProjectName = projectName
		mdl.Data = availdatedata.ScadaDataHFD

		d.BaseController.Ctx.Insert(mdl)

		mdl = NewLatestDataPeriod()
		mdl.Type = "Warning"
		mdl.ProjectName = projectName
		mdl.Data = availdatedata.Warning

		d.BaseController.Ctx.Insert(mdl)

		mdl = NewLatestDataPeriod()
		mdl.Type = "ScadaDataOEM"
		mdl.ProjectName = projectName
		mdl.Data = availdatedata.ScadaDataOEM

		d.BaseController.Ctx.Insert(mdl)
	}

}

func getDataDateAvailable(conn dbox.IConnection, collectionName string, timestampColumn string, where *dbox.Filter) (min time.Time, max time.Time, err error) {
	q := conn.
		NewQuery().
		From(collectionName)

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
		csr.Close()
		return
	}

	data := []toolkit.M{}
	err = csr.Fetch(&data, 0, false)

	if err != nil || len(data) == 0 {
		csr.Close()
		return
	}

	min = data[0].Get("min", time.Time{}).(time.Time).UTC()
	max = data[0].Get("max", time.Time{}).(time.Time).UTC()

	csr.Close()
	return
}
