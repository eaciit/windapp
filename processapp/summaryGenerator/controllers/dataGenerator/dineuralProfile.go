package generatorControllers

import (
	. "eaciit/wfdemo-git/library/models"
	. "eaciit/wfdemo-git/processapp/summaryGenerator/controllers"
	_ "github.com/eaciit/dbox/dbc/mongo"
	tk "github.com/eaciit/toolkit"
	"time"
)

type DineuralProfileSummary struct {
	*BaseController
}

type DineuralProfile struct {
	ID               string ` bson:"_id" , json:"_id" `
	Projectname      string
	Turbine          string
	MonthDesc        string
	MonthID          int
	Hours            string
	WindSpeedTotal   float64
	TemperatureTotal float64
	PowerTotal       float64
	WindSpeedCount   float64
	TemperatureCount float64
	PowerCount       float64
	Type             string
}

func (m *DineuralProfile) TableName() string {
	return "rpt_dineuralprofile"
}

func (ev *DineuralProfileSummary) CreateDineuralProfileSummary(base *BaseController) {
	ev.BaseController = base

	ev.Log.AddLog("===================== Start processing Dineural Profile Summary...", sInfo)
	t0 := time.Now()
	ev.processDataScada()

	ev.Log.AddLog(tk.Sprintf("Duration processing scada data %f minutes", time.Since(t0).Minutes()), sInfo)

	t0 = time.Now()
	ev.processDataMet()

	ev.Log.AddLog(tk.Sprintf("Duration process met tower data %f minutes", time.Since(t0).Minutes()), sInfo)

	ev.Log.AddLog("===================== End processing Dineural Profile Summary...", sInfo)
}

func (ev *DineuralProfileSummary) processDataScada() {
	pipe := []tk.M{
		tk.M{"$match": tk.M{"available": 1}},
		tk.M{"$group": tk.M{
			"_id": tk.M{
				"projectname": "$projectname",
				"turbine":     "$turbine",
				"monthdesc":   "$dateinfo.monthdesc",
				"monthid":     "$dateinfo.monthid",
				"hours":       tk.M{"$dateToString": tk.M{"format": "%H:00", "date": "$timestamp"}},
			},
			"windspeedtotal":   tk.M{"$sum": "$avgwindspeed"},
			"temperaturetotal": tk.M{"$sum": "$nacelletemperature"},
			"powertotal":       tk.M{"$sum": "$power"},
			"windspeedcount":   tk.M{"$sum": 1},
			"temperaturecount": tk.M{"$sum": 1},
			"powercount":       tk.M{"$sum": 1},
		}},
	}

	dineuralData := []tk.M{}
	csr, e := ev.Ctx.Connection.NewQuery().
		From(new(ScadaData).TableName()).
		Command("pipe", pipe).Cursor(nil)
	if e != nil {
		ev.Log.AddLog(tk.Sprintf("Error on cursor : %s", e.Error()), sError)
	}
	defer csr.Close()

	e = csr.Fetch(&dineuralData, 0, false)
	if e != nil {
		ev.Log.AddLog(tk.Sprintf("Error on Fetch : %s", e.Error()), sError)
	}

	data := DineuralProfile{}

	csrSave := ev.Ctx.Connection.NewQuery().SetConfig("multiexec", true).
		From(new(DineuralProfile).TableName()).Save()
	defer csrSave.Close()

	for _, val := range dineuralData {
		data = DineuralProfile{}
		ids := val.Get("_id", tk.M{}).(tk.M)
		data.Projectname = ids.GetString("projectname")
		data.Turbine = ids.GetString("turbine")
		data.MonthDesc = ids.GetString("monthdesc")
		data.MonthID = ids.GetInt("monthid")
		data.Hours = ids.GetString("hours")
		data.ID = tk.Sprintf("%s_%s_%s_%s", data.Projectname, data.Turbine, tk.ToString(data.MonthID), data.Hours)

		data.WindSpeedTotal = val.GetFloat64("windspeedtotal")
		data.TemperatureTotal = val.GetFloat64("temperaturetotal")
		data.PowerTotal = val.GetFloat64("powertotal")
		data.WindSpeedCount = val.GetFloat64("windspeedcount")
		data.TemperatureCount = val.GetFloat64("temperaturecount")
		data.PowerCount = val.GetFloat64("powercount")
		data.Type = "SCADA"

		e = csrSave.Exec(tk.M{"data": data})
		if e != nil {
			ev.Log.AddLog(tk.Sprintf("Error on Save : %s", e.Error()), sError)
		}
	}
}

func (ev *DineuralProfileSummary) processDataMet() {
	pipe := []tk.M{
		tk.M{"$group": tk.M{
			"_id": tk.M{
				"projectname": "$projectname",
				"monthdesc":   "$dateinfo.monthdesc",
				"monthid":     "$dateinfo.monthid",
				"hours":       tk.M{"$dateToString": tk.M{"format": "%H:00", "date": "$timestamp"}},
			},
			"windspeedtotal":   tk.M{"$sum": "$vhubws90mavg"},
			"temperaturetotal": tk.M{"$sum": "$thubhhubtemp855mavg"},
			"windspeedcount":   tk.M{"$sum": "$vhubws90mcount"},
			"temperaturecount": tk.M{"$sum": "$thubhhubtemp855mcount"},
		}},
	}

	dineuralData := []tk.M{}
	csr, e := ev.Ctx.Connection.NewQuery().
		From(new(MetTower).TableName()).
		Command("pipe", pipe).Cursor(nil)
	if e != nil {
		ev.Log.AddLog(tk.Sprintf("Error on cursor : %s", e.Error()), sError)
	}
	defer csr.Close()

	e = csr.Fetch(&dineuralData, 0, false)
	if e != nil {
		ev.Log.AddLog(tk.Sprintf("Error on Fetch : %s", e.Error()), sError)
	}

	data := DineuralProfile{}

	csrSave := ev.Ctx.Connection.NewQuery().SetConfig("multiexec", true).
		From(new(DineuralProfile).TableName()).Save()
	defer csrSave.Close()

	for _, val := range dineuralData {
		data = DineuralProfile{}
		ids := val.Get("_id", tk.M{}).(tk.M)
		data.Projectname = ids.GetString("projectname")
		data.MonthDesc = ids.GetString("monthdesc")
		data.MonthID = ids.GetInt("monthid")
		data.Hours = ids.GetString("hours")
		data.ID = tk.Sprintf("%s_%s_%s", data.Projectname, tk.ToString(data.MonthID), data.Hours)

		data.WindSpeedTotal = val.GetFloat64("windspeedtotal")
		data.TemperatureTotal = val.GetFloat64("temperaturetotal")
		data.WindSpeedCount = val.GetFloat64("windspeedcount")
		data.TemperatureCount = val.GetFloat64("temperaturecount")
		data.Type = "MET"

		e = csrSave.Exec(tk.M{"data": data})
		if e != nil {
			ev.Log.AddLog(tk.Sprintf("Error on Save : %s", e.Error()), sError)
		}
	}
}
