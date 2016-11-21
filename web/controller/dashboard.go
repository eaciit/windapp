package controller

import (
	. "eaciit/wfdemo-git/library/core"
	. "eaciit/wfdemo-git/library/helper"
	. "eaciit/wfdemo-git/library/models"
	"eaciit/wfdemo-git/web/helper"
	"fmt"
	"strings"

	"github.com/eaciit/orm"

	"gopkg.in/mgo.v2/bson"

	"sort"
	// "strings"
	"math"
	"time"

	"github.com/eaciit/dbox"
	"github.com/eaciit/knot/knot.v1"
	tk "github.com/eaciit/toolkit"
)

type DashboardController struct {
	App
}

func CreateDashboardController() *DashboardController {
	var controller = new(DashboardController)
	return controller
}

var (
	turbineMW = 2.1
)

type PayloadDashboard struct {
	ProjectName string
	Turbine     string
	Type        string
	Date        time.Time
	DateStr     string
	Skip        int
	Take        int
	Sort        []Sorting
}

func (m *DashboardController) GetScadaSummary(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	csr, e := DB().Connection.NewQuery().From("rpt_scadasummary").Cursor(nil)

	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}
	defer csr.Close()

	data := []tk.M{}
	e = csr.Fetch(&data, 0, false)

	result := []string{}

	for _, val := range data {
		result = append(result, val.GetString("ID"))
	}
	sort.Strings(result)

	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	return helper.CreateResult(true, result, "success")
}

func (m *DashboardController) GetDashboardSummary(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	csr, e := DB().Connection.NewQuery().From("rpt_dashboardsummary").Cursor(nil)

	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}
	defer csr.Close()

	data := []tk.M{}
	e = csr.Fetch(&data, 0, false)

	result := []string{}

	for _, val := range data {
		result = append(result, val.GetString("ID"))
	}
	sort.Strings(result)

	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	return helper.CreateResult(true, result, "success")
}

func getMachineDownType() (map[string]string, error) {
	/*csr, e := DB().Connection.NewQuery().From("ref_machine_down").Cursor(nil)

	if e != nil {
		return nil, e
	}
	defer csr.Close()

	data := []tk.M{}
	e = csr.Fetch(&data, 0, false)

	if e != nil {
		return nil, e
	}*/

	result := map[string]string{
		"aebok":        "AEBOK",
		"externalstop": "External Stop",
		"griddown":     "Grid Down",
		"internalgrid": "Internal Grid",
		"machinedown":  "Machine Down",
		"unknown":      "Unknown",
		"weatherstop":  "Weather Stop",
	}

	return result, nil
}

func (m *DashboardController) GetMDTypeList(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	data, e := getMachineDownType()
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	result := []string{}
	for _, val := range data {
		result = append(result, val)
	}
	sort.Strings(result)

	return helper.CreateResult(true, result, "success")
}

func getProject() ([]string, error) {
	csr, e := DB().Connection.NewQuery().From("ref_project").Cursor(nil)

	if e != nil {
		return nil, e
	}
	defer csr.Close()

	data := []tk.M{}
	e = csr.Fetch(&data, 0, false)

	if e != nil {
		return nil, e
	}

	result := []string{}

	for _, val := range data {
		if val.GetString("projectid") == "Tejuva" {
			result = append(result, val.GetString("projectid"))
		}
	}
	sort.Strings(result)

	return result, nil
}

func (m *DashboardController) GetProjectList(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	result, e := getProject()
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	return helper.CreateResult(true, result, "success")
}

func (m *DashboardController) GetScadaLastUpdate(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	p := new(PayloadDashboard)
	e := k.GetPayload(&p)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	csr, e := DB().Connection.NewQuery().From("rpt_scadalastupdate").Where(dbox.And(dbox.Eq("projectname", p.ProjectName))).Cursor(nil)

	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}
	defer csr.Close()

	data := []ScadaLastUpdate{}
	e = csr.Fetch(&data, 0, false)

	// conver the time into utc
	result := []ScadaLastUpdate{}
	for _, val := range data {
		val.DateInfo = GetDateInfo(val.DateInfo.DateId.UTC())
		val.LastUpdate = val.LastUpdate.UTC()

		for idxProd, prod := range val.Productions {
			prod.TimeHour = prod.TimeHour.UTC()
			val.Productions[idxProd] = prod
		}

		for idxCumm, cumm := range val.CummulativeProductions {
			cumm.DateId = cumm.DateId.UTC()
			val.CummulativeProductions[idxCumm] = cumm
		}

		result = append(result, val)
	}

	return helper.CreateResult(true, result, "success")
}

type ScadaSummaryVariance struct {
	orm.ModelBase      `bson:"-",json:"-"`
	ID                 bson.ObjectId ` bson:"_id" , json:"_id" `
	DateInfo           DateInfo
	ProjectName        string
	Production         float64
	ProductionLastYear float64
	Revenue            float64
	RevenueInLacs      float64
	TrueAvail          float64
	ScadaAvail         float64
	MachineAvail       float64
	GridAvail          float64
	PLF                float64
	Budget             float64
	AvgWindSpeed       float64
	ExpWindSpeed       float64
	DowntimeHours      float64
	LostEnergy         float64
	RevenueLoss        float64
	Variance           float64
}

func (m *DashboardController) GetScadaSummaryByMonth(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson
	p := new(PayloadDashboard)
	e := k.GetPayload(&p)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	// Get lastmonth
	csrlastmonth, e := DB().Connection.NewQuery().From("rpt_scadalastupdate").Where(dbox.And(dbox.Eq("projectname", p.ProjectName))).Cursor(nil)

	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}
	defer csrlastmonth.Close()

	datalastmonth := make([]ScadaSummaryByMonth, 0)
	e = csrlastmonth.Fetch(&datalastmonth, 0, false)

	startmonth := 0
	endmonth := datalastmonth[0].DateInfo.MonthId
	month := endmonth - (int(endmonth/100) * 100)

	if month == 12 {
		startmonth = (int(endmonth/100) * 100) + 1
	} else {
		startmonth = (endmonth + 1) - 100
	}

	// result := make([]ScadaSummaryByMonth, 0)
	var result []interface{}
	dataVariance := new(ScadaSummaryVariance)

	for i := startmonth; i <= endmonth; i++ {
		//check if month more than 12
		if i-(int(i/100)*100) > 12 {
			i = (i - 12) + 100
		}

		yearloop := int(i / 100)
		monthloop := i - (int(i/100) * 100)

		csr, e := DB().Connection.NewQuery().From("rpt_scadasummarybymonth").Where(dbox.And(dbox.Eq("dateinfo.monthid", i), dbox.Eq("projectname", p.ProjectName))).Cursor(nil)

		if e != nil {
			return helper.CreateResult(false, nil, e.Error())
		}
		defer csr.Close()

		data := make([]ScadaSummaryByMonth, 0)
		e = csr.Fetch(&data, 0, false)

		if len(data) > 0 {
			dataVariance.ID = data[0].ID
			dataVariance.DateInfo = data[0].DateInfo
			dataVariance.ProjectName = data[0].ProjectName
			dataVariance.Production = data[0].Production / 1000
			dataVariance.ProductionLastYear = data[0].ProductionLastYear
			dataVariance.Revenue = data[0].Revenue
			dataVariance.RevenueInLacs = data[0].RevenueInLacs
			dataVariance.TrueAvail = data[0].TrueAvail
			dataVariance.ScadaAvail = data[0].ScadaAvail
			dataVariance.MachineAvail = data[0].MachineAvail
			dataVariance.GridAvail = data[0].GridAvail
			dataVariance.PLF = data[0].PLF
			dataVariance.Budget = data[0].Budget / 1000000
			dataVariance.AvgWindSpeed = data[0].AvgWindSpeed
			dataVariance.ExpWindSpeed = data[0].ExpWindSpeed
			dataVariance.DowntimeHours = data[0].DowntimeHours
			dataVariance.LostEnergy = data[0].LostEnergy
			dataVariance.RevenueLoss = data[0].RevenueLoss
			if data[0].ProductionLastYear == 0 {
				dataVariance.Variance = 100
			} else {
				dataVariance.Variance = math.Abs((data[0].Production - data[0].ProductionLastYear) / data[0].ProductionLastYear * 100)
			}

			result = append(result, *dataVariance)
		} else {
			// Temporary data to fill result if month doesn't exist
			datatemp := new(ScadaSummaryByMonth)

			datatemp.ID = ""
			datatemp.ProjectName = p.ProjectName
			dateInfo := GetDateInfo(time.Date(yearloop, time.Month(monthloop), 1, 17, 0, 0, 0, time.UTC))
			datatemp.DateInfo = dateInfo

			dataVariance.ID = datatemp.ID
			dataVariance.DateInfo = datatemp.DateInfo
			dataVariance.ProjectName = datatemp.ProjectName
			dataVariance.Production = datatemp.Production
			dataVariance.ProductionLastYear = datatemp.ProductionLastYear
			dataVariance.Revenue = datatemp.Revenue
			dataVariance.RevenueInLacs = datatemp.RevenueInLacs
			dataVariance.TrueAvail = datatemp.TrueAvail
			dataVariance.ScadaAvail = datatemp.ScadaAvail
			dataVariance.MachineAvail = datatemp.MachineAvail
			dataVariance.GridAvail = datatemp.GridAvail
			dataVariance.PLF = datatemp.PLF
			dataVariance.Budget = datatemp.Budget
			dataVariance.AvgWindSpeed = datatemp.AvgWindSpeed
			dataVariance.ExpWindSpeed = datatemp.ExpWindSpeed
			dataVariance.DowntimeHours = datatemp.DowntimeHours
			dataVariance.LostEnergy = datatemp.LostEnergy
			dataVariance.RevenueLoss = datatemp.RevenueLoss
			dataVariance.Variance = 0

			result = append(result, *dataVariance)
		}
	}

	return helper.CreateResult(true, result, "success")
}

func (m *DashboardController) GetDetailProd(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	p := tk.M{}
	e := k.GetPayload(&p)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	ids := tk.M{"project": "$projectname", "turbine": "$turbine"}
	matches := tk.M{"dateinfo.monthdesc": p.GetString("date")}
	if p.GetString("project") != "Fleet" {
		matches.Set("projectname", p.GetString("project"))
	}
	fields := tk.M{"$sum": "$power"}

	pipe := []tk.M{{"$match": matches}, {"$group": tk.M{"_id": ids, "production": fields}}}
	csrScada, e := DB().Connection.NewQuery().
		From(new(ScadaData).TableName()).
		Command("pipe", pipe).
		Cursor(nil)

	if e != nil {
		helper.CreateResult(false, nil, e.Error())
	}
	defer csrScada.Close()

	resultScada := []tk.M{}
	e = csrScada.Fetch(&resultScada, 0, false)

	matches.Unset("dateinfo.monthdesc")
	matches.Set("startdateinfo.monthdesc", p.GetString("date"))
	fields = tk.M{"$sum": "$powerlost"}

	pipe = []tk.M{{"$match": matches}, {"$group": tk.M{"_id": ids, "totalpowerlost": fields,
		"totalduration": tk.M{"$sum": "$duration"}}}}
	csrAlarm, e := DB().Connection.NewQuery().
		From(new(Alarm).TableName()).
		Command("pipe", pipe).
		Cursor(nil)

	if e != nil {
		helper.CreateResult(false, nil, e.Error())
	}
	defer csrAlarm.Close()

	resultAlarm := []tk.M{}
	e = csrAlarm.Fetch(&resultAlarm, 0, false)

	dataPowerLost := tk.M{}
	dataDuration := tk.M{}
	for _, val := range resultAlarm {
		data := val["_id"].(tk.M)
		dataPowerLost.Set(data.GetString("project")+"_"+data.GetString("turbine"), val.GetFloat64("totalpowerlost"))
		dataDuration.Set(data.GetString("project")+"_"+data.GetString("turbine"), val.GetFloat64("totalduration"))
	}
	totalPower := tk.M{}
	totalPowerLost := tk.M{}
	totalTurbines := tk.M{}
	detailData := tk.M{}
	detail := []tk.M{}
	for _, val := range resultScada {
		data := val["_id"].(tk.M)
		project := data.GetString("project")
		_id := project + "_" + data.GetString("turbine")
		if dataPowerLost.Has(_id) {
			val.Set("lostenergy", dataPowerLost.GetFloat64(_id))
			val.Set("downtime", dataDuration.GetFloat64(_id))
		}
		val.Unset("_id")
		val.Set("production", val.GetFloat64("production")/6)
		val.Set("turbine", data.GetString("turbine"))
		detail = append(detail, val)
		detailData.Set(project, detail)

		if totalPower.Has(project) {
			totalPower.Set(data.GetString("project"), totalPower.GetFloat64(project)+val.GetFloat64("production"))
		} else {
			totalPower.Set(data.GetString("project"), val.GetFloat64("production"))
		}
		if totalPowerLost.Has(project) {
			totalPowerLost.Set(data.GetString("project"), totalPowerLost.GetFloat64(project)+val.GetFloat64("lostenergy"))
		} else {
			totalPowerLost.Set(data.GetString("project"), val.GetFloat64("lostenergy"))
		}
		if totalTurbines.Has(project) {
			totalTurbines.Set(data.GetString("project"), totalTurbines.GetInt(project)+1)
		} else {
			totalTurbines.Set(data.GetString("project"), 1)
		}
	}

	dataItem := []tk.M{}
	for project, val := range totalPower {
		data := tk.M{
			"project":    project,
			"production": val.(float64),
			"lostenergy": totalPowerLost.GetFloat64(project),
			"wtg":        totalTurbines.GetInt(project),
			"detail":     detailData[project],
		}
		dataItem = append(dataItem, data)
	}
	dataItemTemp := dataItem
	dataItem = []tk.M{}
	for _, val := range dataItemTemp {
		newdata := helper.EnergyMeasurement(val, "production", "lostenergy")
		val = newdata[0]
		newdetail := helper.EnergyMeasurement(val["detail"].([]tk.M), "production", "lostenergy")
		val.Set("detail", newdetail)
		dataItem = append(dataItem, val)
	}

	return helper.CreateResult(true, dataItem, "success")
}

func (m *DashboardController) GetSummaryData(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	p := new(helper.Payloads)
	e := k.GetPayload(&p)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	filter, _ := p.ParseFilter()
	fb := DB().Connection.Fb()
	fb.AddFilter(dbox.And(filter...))
	matches, e := fb.Build()
	if e != nil {
		helper.CreateResult(false, nil, e.Error())
	}

	arrsort := tk.M{}
	if len(p.Sort) > 0 {
		for _, val := range p.Sort {
			if val.Dir == "desc" {
				arrsort.Set(strings.ToLower("dataitems."+val.Field), 0)
			} else {
				arrsort.Set(strings.ToLower("dataitems."+val.Field), 1)
			}
		}
	}

	pipe := []tk.M{{"$unwind": "$dataitems"}, {"$match": matches}, {"$sort": arrsort}, {"$skip": p.Skip}, {"$limit": p.Take}}
	csr, e := DB().Connection.NewQuery().
		From(new(ScadaSummaryByProject).TableName()).
		Command("pipe", pipe).
		Cursor(nil)

	if e != nil {
		helper.CreateResult(false, nil, e.Error())
	}
	defer csr.Close()

	result := []tk.M{}
	e = csr.Fetch(&result, 0, false)
	dataItem := []tk.M{}

	for _, val := range result {
		dataItem = append(dataItem, val["dataitems"].(tk.M))
	}

	csrCount, e := DB().Connection.NewQuery().From(new(ScadaSummaryByProject).TableName()).
		Where(dbox.And(filter...)).Cursor(nil)
	if e != nil {
		helper.CreateResult(false, nil, e.Error())
	}
	defer csrCount.Close()
	dataCount := new(ScadaSummaryByProject)
	e = csrCount.Fetch(&dataCount, 1, false)

	data := struct {
		Data  []tk.M
		Total int
	}{
		Data:  dataItem,
		Total: tk.SliceLen(dataCount.DataItems),
	}

	return helper.CreateResult(true, data, "success")
}

func (m *DashboardController) GetDownTime(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	result := tk.M{}

	p := new(PayloadDashboard)
	e := k.GetPayload(&p)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	downtimeDatas := getDownTimeLostEnergy("project", p)
	// tk.Printf("Payload => %#v\n", p)
	// tk.Printf("Data LostEnergy ==> %#v\n", downtimeDatas)
	result.Set("lostenergy", downtimeDatas)

	if p.Type == "" && p.ProjectName == "Fleet" {
		result.Set("lostenergybytype", getDownTimeLostEnergy("type", p))
	}
	result.Set("duration", getTurbineDownTimeTop("duration", p))
	result.Set("frequency", getTurbineDownTimeTop("frequency", p))
	result.Set("loss", getTurbineDownTimeTop("loss", p))

	result.Set("lossCatDuration", getLossCategoriesTop("duration", p))
	result.Set("lossCatFrequency", getLossCategoriesTop("frequency", p))
	result.Set("lossCatLoss", getLossCategoriesTop("loss", p))

	result.Set("machineAvailability", getAvailability("machine", p))
	result.Set("gridAvailability", getAvailability("grid", p))

	return helper.CreateResult(true, result, "success")
}

func (m *DashboardController) GetDownTimeFleetByDown(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	result := tk.M{}

	p := new(PayloadDashboard)
	e := k.GetPayload(&p)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	result.Set("lostenergy", getDownTimeLostEnergy("fleetdowntime", p))

	return helper.CreateResult(true, result, "success")
}

func getDownTimeLostEnergy(tipe string, p *PayloadDashboard) (result []tk.M) {
	var pipes []tk.M
	var pipesDown []tk.M
	var fromDate time.Time
	match := tk.M{}
	matchDown := tk.M{}

	if p.DateStr != "" {
		dateStr := strings.Split(p.DateStr, " ")
		date, e := time.Parse("Jan 2006 02 15:04:05", dateStr[0][0:3]+" "+dateStr[1]+" 01 00:00:00")
		if e != nil {
			return
		}

		dateInfo := GetDateInfo(date)

		if tipe == "fleetdowntime" {
			matchDown.Set("startdateinfo.monthid", dateInfo.MonthId)
		} else {
			match.Set("dateinfo.monthid", dateInfo.MonthId)
		}
	} else {
		fromDate = p.Date.AddDate(0, -12, 0)
		match.Set("dateinfo.dateid", tk.M{"$gte": fromDate.UTC(), "$lte": p.Date})
	}

	if p.ProjectName != "Fleet" {
		match.Set("projectname", p.ProjectName)
		matchDown.Set("projectname", p.ProjectName)
	}

	if p.Type != "" && tipe != "fleetdowntime" && p.Type != "All Types" {
		match.Set("type", p.Type)
	} else if p.Type != "" && tipe == "fleetdowntime" {
		matchDown.Set(strings.Replace(strings.ToLower(p.Type), " ", "", 1), true)
	}

	pipes = append(pipes, tk.M{"$match": match})

	if p.ProjectName != "Fleet" {
		// add a condition to check the type is project
		// regarding to next process can not catch the value for selecting downtime by project in dashboard
		// add by ams on 20161003
		if tipe == "project" {
			if p.Type == "All Types" {
				pipes = append(pipes,
					tk.M{
						"$group": tk.M{"_id": tk.M{"id1": "$dateinfo.monthid", "id2": "$dateinfo.monthdesc", "id3": "$type"},
							"result": tk.M{"$sum": "$lostenergy"},
						},
					},
				)
			} else {
				pipes = append(pipes,
					tk.M{
						/*"$group": tk.M{"_id": tk.M{"id1": "$dateinfo.monthid", "id2": "$dateinfo.monthdesc", "id3": "$projectname"},
						"result": tk.M{"$sum": "$lostenergy"},*/
						"$group": tk.M{"_id": tk.M{"id1": "$dateinfo.monthid", "id2": "$dateinfo.monthdesc", "id3": "$type"},
							"result": tk.M{"$sum": "$lostenergy"}, /*changed from by project to by MD type per 11 Oct 16 [RS]*/
						},
					},
				)
			}

		} else {
			pipes = append(pipes,
				tk.M{
					"$group": tk.M{"_id": tk.M{"id1": "$dateinfo.monthid", "id2": "$dateinfo.monthdesc", "id3": "$type"},
						"result": tk.M{"$sum": "$lostenergy"},
					},
				},
			)
		}
	} else {
		if tipe == "project" {
			pipes = append(pipes,
				tk.M{
					"$group": tk.M{"_id": tk.M{"id1": "$dateinfo.monthid", "id2": "$dateinfo.monthdesc", "id3": "$projectname"},
						"result": tk.M{"$sum": "$lostenergy"},
					},
				},
			)
		} else {
			pipes = append(pipes,
				tk.M{
					"$group": tk.M{"_id": tk.M{"id1": "$type", "id2": "$type", "id3": "$projectname"},
						"powerlost": tk.M{"$sum": "$lostenergy"},
					},
				},
			)
		}
	}

	machinedown, e := getMachineDownType()
	if e != nil {
		return nil
	}

	if p.DateStr == "" && tipe != "fleetdowntime" {
		pipes = append(pipes, tk.M{"$sort": tk.M{"_id.id3": 1}})

		csr, e := DB().Connection.NewQuery().
			From(new(AlarmSummaryByMonth).TableName()).
			Command("pipe", pipes).
			Cursor(nil)

		if e != nil {
			return
		}

		tmpResult := []tk.M{}
		e = csr.Fetch(&tmpResult, 0, false)

		// add by ams, 2016-10-07
		csr.Close()

		if e != nil {
			return
		}

		stack := map[string]string{}

		if p.ProjectName != "Fleet" {
			if p.Type == "All Types" {
				stack = machinedown
			} else {
				// stack[p.ProjectName] = p.ProjectName
				stack = machinedown
				/*changed from by project to by MD type per 11 Oct 16 [RS]*/
			}
		} else {
			project, e := getProject()
			if e != nil {
				return nil
			}
			for _, val := range project {
				stack[strings.ToLower(val)] = val
			}
		}

		dt, _ := time.Parse("2006-01-02 15:04:05", fromDate.UTC().Format("2006-01")+"-01 00:00:00")
		lineData := tk.M{}

		for _, title := range stack {
			if tipe != "type" {
				for i := 1; i < 13; i++ {
					currDate := dt.AddDate(0, i, 0)
					dateInfo := GetDateInfo(currDate)
					found := false

					for _, val := range tmpResult {
						id := val.Get("_id").(tk.M)
						id1 := id.GetInt("id1")
						id3 := id.GetString("id3")

						// tk.Printf("ID 1 => %#v\n", id1)
						// tk.Printf("MonthId => %#v\n", dateInfo.MonthId)
						// tk.Printf("ID 3 => %#v\n", id3)
						// tk.Printf("Title => %#v\n", tk.ToString(title))
						// tk.Printf("Value => %#v\n", val.GetFloat64("result"))

						if id1 == dateInfo.MonthId && id3 == tk.ToString(title) {
							found = true

							val.Set("result", val.GetFloat64("result")*0.001)
							result = append(result, val)
							break
						}
					}

					if !found {
						emptyRes := tk.M{}
						emptyRes.Set("_id", tk.M{"id1": dateInfo.MonthId, "id2": dateInfo.MonthDesc, "id3": tk.ToString(title)})
						emptyRes.Set("result", 0)

						result = append(result, emptyRes)
					}
				}
			} else if tipe == "type" {
				source := []tk.M{}
				groups := tk.M{
					"_id":       "$projectname",
					"duration":  tk.M{"$sum": "$duration"},
					"frequency": tk.M{"$sum": 1},
				}
				var bigPower, bigDuration float64
				var bigFreq int
				for field, mdName := range machinedown {
					matchX := tk.M{field: true}
					pipesX := []tk.M{{"$match": matchX}, {"$group": groups}}

					csr, e := DB().Connection.NewQuery().
						From(new(Alarm).TableName()).
						Command("pipe", pipesX).
						Cursor(nil)

					if e != nil {
						return
					}

					tmpRes := []tk.M{}
					e = csr.Fetch(&tmpRes, 0, false)
					// add by ams, 2016-10-07
					csr.Close()

					if e != nil {
						return
					}

					if len(tmpRes) > 0 {
						tmp := tmpRes[0]
						ids := tmp.GetString("_id")
						lineData.Set(mdName+"_"+ids, tk.M{"duration": tmp.GetFloat64("duration"),
							"frequency": tmp.GetInt("frequency")})
					}

					found := false
					for _, val := range tmpResult {
						id := val.Get("_id").(tk.M)
						id1 := id.GetString("id1")
						id3 := id.GetString("id3")
						if id1 == mdName && id3 == tk.ToString(title) {
							line := lineData[mdName+"_"+title].(tk.M)
							found = true
							powerlost := val.GetFloat64("powerlost") * 0.001
							duration := line.GetFloat64("duration")
							frequency := line.GetInt("frequency")

							if powerlost > bigPower {
								bigPower = powerlost
							}
							if duration > bigDuration {
								bigDuration = duration
							}
							if frequency > bigFreq {
								bigFreq = frequency
							}
							val.Set("powerlost", powerlost)
							val.Set("duration", duration)
							val.Set("frequency", frequency)
							source = append(source, val)
							break
						}
					}
					if !found {
						emptyRes := tk.M{}
						emptyRes.Set("_id", tk.M{"id1": mdName, "id2": mdName, "id3": tk.ToString(title)})
						emptyRes.Set("powerlost", 0)
						emptyRes.Set("duration", 0)
						emptyRes.Set("frequency", 0)

						source = append(source, emptyRes)
					}
				}
				data := tk.M{
					"maxPowerLost": bigPower * 3,
					"minPowerLost": 0,
					"maxDuration":  bigDuration * 3,
					"minDuration":  bigDuration - (bigDuration * 3),
					"maxFreq":      bigFreq * 2,
					"minFreq":      bigFreq - (bigFreq * 2),
					"source":       source}
				result = append(result, data)
			}
		}
	} else {
		if e != nil {
			return
		}

		if p.Type != "" && p.Type != "All Types" {
			pipesDown = append(pipesDown, tk.M{"$match": matchDown})
			pipesX := pipesDown
			pipesX = append(pipesX,
				tk.M{
					"$group": tk.M{"_id": tk.M{"id1": "$startdateinfo.monthid", "id2": "$startdateinfo.monthdesc", "id3": p.Type},
						"result": tk.M{"$sum": "$powerlost"},
					},
				},
			)

			// tk.Printf("%#v \n", pipesX)

			csr, e := DB().Connection.NewQuery().
				From(new(Alarm).TableName()).
				Command("pipe", pipesX).
				Cursor(nil)

			if e != nil {
				return
			}

			tmpRes := []tk.M{}
			e = csr.Fetch(&tmpRes, 0, false)
			// add by ams, 2016-10-07
			csr.Close()

			if e != nil {
				return
			}

			if len(tmpRes) > 0 {
				id := tmpRes[0].Get("_id").(tk.M)
				id3 := id.Get("id3")
				resVal := tmpRes[0].GetFloat64("result") / 1000.0

				for _, title := range machinedown {

					if id3 == title {
						result = append(result, tk.M{"type": title, "result": resVal})
					} else {
						result = append(result, tk.M{"type": title, "result": 0})
					}
				}
			}
		} else {

			doneField := []string{}

			for field, title := range machinedown {
				matchX := matchDown

				for _, done := range doneField {
					matchX.Unset(done)
				}

				doneField = append(doneField, field)
				matchX.Set(field, true)

				pipesDown = append(pipesDown, tk.M{"$match": matchX})
				pipesX := pipesDown

				pipesX = append(pipesX,
					tk.M{
						"$group": tk.M{"_id": tk.M{"id1": "$startdateinfo.monthid", "id2": "$startdateinfo.monthdesc", "id3": field},
							"result": tk.M{"$sum": "$powerlost"},
						},
					},
				)

				csr, e := DB().Connection.NewQuery().
					From(new(Alarm).TableName()).
					Command("pipe", pipesX).
					Cursor(nil)

				if e != nil {
					return
				}

				tmpRes := []tk.M{}
				e = csr.Fetch(&tmpRes, 0, false)
				// add by ams, 2016-10-07
				csr.Close()

				if e != nil {
					return
				}

				if len(tmpRes) > 0 {
					tmp := tmpRes[0]
					result = append(result, tk.M{"type": title, "result": tmp.GetFloat64("result") / 1000.0})
				} else {
					result = append(result, tk.M{"type": title, "result": 0})
				}
			}
		}

	}

	return
}

func getTurbineDownTimeTop(topType string, p *PayloadDashboard) (result []tk.M) {
	var pipes []tk.M
	var fromDate time.Time
	match := tk.M{}

	if p.DateStr == "" {
		fromDate = p.Date.AddDate(0, -12, 0)

		match.Set("startdate", tk.M{"$gte": fromDate.UTC(), "$lte": p.Date.UTC()})

		if p.ProjectName != "Fleet" {
			match.Set("projectname", p.ProjectName)
		}

		pipes = append(pipes, tk.M{"$match": match})

		if topType == "duration" {
			pipes = append(pipes, tk.M{"$group": tk.M{"_id": "$turbine", "result": tk.M{"$sum": "$duration"}}})
		} else if topType == "frequency" {
			pipes = append(pipes, tk.M{"$group": tk.M{"_id": "$turbine", "result": tk.M{"$sum": 1}}})
		} else if topType == "loss" {
			pipes = append(pipes, tk.M{"$group": tk.M{"_id": "$turbine", "result": tk.M{"$sum": "$powerlost"}}})
		}

		pipes = append(pipes, tk.M{"$sort": tk.M{"result": -1}})
		pipes = append(pipes, tk.M{"$limit": 10})

		// get the top 10

		csr, e := DB().Connection.NewQuery().
			Select("_id").
			From(new(Alarm).TableName()).
			Command("pipe", pipes).
			Cursor(nil)

		if e != nil {
			return
		}

		top10Turbines := []tk.M{}
		e = csr.Fetch(&top10Turbines, 0, false)
		// add by ams, 2016-10-07
		csr.Close()

		if e != nil {
			return
		}

		// get the downtime

		turbines := []string{}
		turbinesVal := tk.M{}

		for _, turbine := range top10Turbines {
			turbines = append(turbines, turbine.Get("_id").(string))
			turbinesVal.Set(turbine.Get("_id").(string), turbine.GetFloat64("result"))
		}

		// tk.Printf("topType: \n%#v \n", topType)
		// tk.Printf("turbines: %#v \n", turbines)

		match.Set("turbine", tk.M{"$in": turbines})

		downCause := tk.M{}
		downCause.Set("aebok", "AEBOK")
		downCause.Set("externalstop", "External Stop")
		downCause.Set("griddown", "Grid Down")
		downCause.Set("internalgrid", "Internal Grid")
		downCause.Set("machinedown", "Machine Down")
		downCause.Set("unknown", "Unknown")
		downCause.Set("weatherstop", "Weather Stop")

		tmpResult := []tk.M{}
		downDone := []string{}

		for f, t := range downCause {
			pipes = []tk.M{}
			loopMatch := match
			field := tk.ToString(f)
			title := tk.ToString(t)

			downDone = append(downDone, field)

			for _, done := range downDone {
				match.Unset(done)
			}

			loopMatch.Set(field, true)

			// tk.Printf("%#v \n", match)

			pipes = append(pipes, tk.M{"$match": loopMatch})
			if topType == "duration" {
				pipes = append(pipes,
					tk.M{
						"$group": tk.M{"_id": tk.M{"id1": "$dateinfo.monthid", "id2": "$dateinfo.monthdesc", "id3": "$turbine", "id4": title},
							"result": tk.M{"$sum": "$duration"},
						},
					},
				)
			} else if topType == "frequency" {
				pipes = append(pipes,
					tk.M{
						"$group": tk.M{"_id": tk.M{"id1": "$dateinfo.monthid", "id2": "$dateinfo.monthdesc", "id3": "$turbine", "id4": title},
							"result": tk.M{"$sum": 1},
						},
					},
				)
			} else if topType == "loss" {
				pipes = append(pipes,
					tk.M{
						"$group": tk.M{"_id": tk.M{"id1": "$dateinfo.monthid", "id2": "$dateinfo.monthdesc", "id3": "$turbine", "id4": title},
							"result": tk.M{"$sum": "$powerlost"},
						},
					},
				)
			}

			pipes = append(pipes, tk.M{"$sort": tk.M{"result": -1}})

			/*tk.Println()
			tk.Println(tk.ToString(title))
			for _, val := range pipes {
				tk.Printf("pipes: %v \n", val)
			}*/

			csr, e := DB().Connection.NewQuery().
				From(new(Alarm).TableName()).
				Command("pipe", pipes).
				Cursor(nil)

			if e != nil {
				return
			}

			resLoop := []tk.M{}
			e = csr.Fetch(&resLoop, 0, false)

			// add by ams, 2016-10-07
			csr.Close()

			// tk.Printf("resLoop: %v - %#v \n", tk.ToString(title), resLoop)

			for _, res := range resLoop {
				tmpResult = append(tmpResult, res)
			}
		}

		/*tk.Printf("len: %v \n", len(tmpResult))
		tk.Printf("%#v \n", tmpResult)*/

		/*for _, val := range tmpResult {
			tk.Printf("tmpResult: %v \n", val)
		}*/

		resY := []tk.M{}

		for _, t := range downCause {
			// field := tk.ToString(f)
			title := tk.ToString(t)

			for _, turbine := range turbines {
				resX := tk.M{}
				resX.Set("_id", tk.M{"id3": turbine, "id4": title})
				resX.Set("result", 0)

			out:
				for _, res := range tmpResult {
					id3 := res.Get("_id").(tk.M).GetString("id3")
					id4 := res.Get("_id").(tk.M).GetString("id4")

					if id3 == turbine && id4 == title {
						resX = res
						break out
					}
				}

				// if title == "External Stop" {
				resY = append(resY, resX)
				// }
			}
		}

		for _, turbine := range turbines {
			resVal := tk.M{}
			resVal.Set("_id", turbine)

			for _, val := range resY {
				valTurbine := val.Get("_id").(tk.M).GetString("id3")
				valResult := val.GetFloat64("result")
				valTitle := ""

				splitTitle := strings.Split(val.Get("_id").(tk.M).GetString("id4"), " ")

				if len(splitTitle) > 1 {
					valTitle = splitTitle[0] + "" + splitTitle[1]
				} else {
					valTitle = splitTitle[0]
				}

				if turbine == valTurbine && valResult != 0 {
					resVal.Set(valTitle, valResult)
				} else if resVal.Get(valTitle) == nil {
					resVal.Set(valTitle, 0)
				}
			}

			resVal.Set("Total", turbinesVal.GetFloat64(turbine))
			result = append(result, resVal)
		}
	}

	return
}

func getLossCategoriesTop(topType string, p *PayloadDashboard) (result []tk.M) {
	var pipes []tk.M
	var fromDate time.Time
	match := tk.M{}

	if p.DateStr == "" {
		fromDate = p.Date.AddDate(0, -12, 0)

		match.Set("startdate", tk.M{"$gte": fromDate.UTC(), "$lte": p.Date.UTC()})

		if p.ProjectName != "Fleet" {
			match.Set("projectname", p.ProjectName)
		}

		downCause := tk.M{}
		downCause.Set("aebok", "AEBOK")
		downCause.Set("externalstop", "External Stop")
		downCause.Set("griddown", "Grid Down")
		downCause.Set("internalgrid", "Internal Grid")
		downCause.Set("machinedown", "Machine Down")
		downCause.Set("unknown", "Unknown")
		downCause.Set("weatherstop", "Weather Stop")

		tmpResult := []tk.M{}
		downDone := []string{}

		for f, t := range downCause {
			pipes = []tk.M{}
			loopMatch := match
			field := tk.ToString(f)
			title := tk.ToString(t)

			downDone = append(downDone, field)

			for _, done := range downDone {
				match.Unset(done)
			}

			loopMatch.Set(field, true)

			pipes = append(pipes, tk.M{"$match": loopMatch})
			if topType == "duration" {
				pipes = append(pipes,
					tk.M{
						"$group": tk.M{"_id": tk.M{"id1": field, "id2": title}, "result": tk.M{"$sum": "$duration"}},
					},
				)
			} else if topType == "frequency" {
				pipes = append(pipes,
					tk.M{
						"$group": tk.M{"_id": tk.M{"id1": field, "id2": title}, "result": tk.M{"$sum": 1}},
					},
				)
			} else if topType == "loss" {
				pipes = append(pipes,
					tk.M{
						"$group": tk.M{"_id": tk.M{"id1": field, "id2": title}, "result": tk.M{"$sum": "$powerlost"}},
					},
				)
			}

			csr, e := DB().Connection.NewQuery().
				From(new(Alarm).TableName()).
				Command("pipe", pipes).
				Cursor(nil)

			if e != nil {
				return
			}

			resLoop := []tk.M{}
			e = csr.Fetch(&resLoop, 0, false)

			csr.Close()

			for _, res := range resLoop {
				tmpResult = append(tmpResult, res)
			}
		}

		size := len(tmpResult)

		if size > 1 {
			for i := 0; i < size; i++ {
				for j := size - 1; j >= i+1; j-- {
					a := tmpResult[j].GetFloat64("result")
					b := tmpResult[j-1].GetFloat64("result")

					if a > b {
						tmpResult[j], tmpResult[j-1] = tmpResult[j-1], tmpResult[j]
					}
				}
			}
		}

		result = tmpResult
	}

	return
}

func getAvailability(availType string, p *PayloadDashboard) (result []tk.M) {
	var fromDate time.Time
	match := tk.M{}

	if p.DateStr == "" {
		fromDate = p.Date.AddDate(0, -12, 0)

		match.Set("dateinfo.dateid", tk.M{"$gte": fromDate.UTC(), "$lte": p.Date.UTC()})

		group := tk.M{}

		if p.ProjectName != "Fleet" {
			match.Set("projectname", p.ProjectName)
		}

		group.Set("_id", tk.M{"id1": "$dateinfo.monthid", "id2": "$dateinfo.monthdesc", "id3": "$projectname"})

		if availType == "machine" {
			group.Set("result", tk.M{"$avg": "$machineavail"})
		} else if availType == "grid" {
			group.Set("result", tk.M{"$avg": "$gridavail"})
		}

		pipe := []tk.M{
			{"$match": match},
			{"$group": group},
			{"$sort": tk.M{"result": -1}},
			{"$limit": 10},
		}

		// tk.Printf("pipe: %#v \n", pipe)

		csr, e := DB().Connection.NewQuery().
			From(new(ScadaSummaryDaily).TableName()).
			Command("pipe", pipe).
			Cursor(nil)

		if e != nil {
			return
		}
		defer csr.Close()

		tmpResult := []tk.M{}

		e = csr.Fetch(&tmpResult, 0, false)
		if e != nil {
			return
		}

		// get project list, for now just using Tejuva
		// project list should come from ref_project collection

		projects := []string{}
		projects = append(projects, "Tejuva")

		// --------------

		// dayInYear := tk.M{}
		tmpFromDate := fromDate.AddDate(0, 1, 0)
		dateInfoTo := GetDateInfo(p.Date)

		for _, project := range projects {

		done:

			for {
				dateInfoFrom := GetDateInfo(tmpFromDate)
				// if dayInYear.Get(tk.ToString(dateInfoFrom.Year)) == nil {
				// 	dayInYear.Set(tk.ToString(dateInfoFrom.Year), GetDayInYear(dateInfoFrom.Year))
				// }
				// days := dayInYear.Get(tk.ToString(dateInfoFrom.Year)).(tk.M).GetInt(tk.ToString(int(tmpFromDate.Month())))

				var exist tk.M

			existData:

				for _, res := range tmpResult {
					id := res.Get("_id").(tk.M)
					if dateInfoFrom.MonthId == id.GetInt("id1") && project == id.GetString("id3") {
						exist = res
						break existData
					}
				}

				if exist != nil {
					// resVal := exist.GetFloat64("result") / tk.ToFloat64(days, 0, tk.RoundingAuto)
					// exist.Set("result", resVal)
					result = append(result, exist)
				} else {
					result = append(result, tk.M{
						"_id": tk.M{
							"id1": dateInfoFrom.MonthId,
							"id2": dateInfoFrom.MonthDesc,
							"id3": project,
						},
						"result": 0.00,
					})
				}

				if dateInfoFrom.MonthId == dateInfoTo.MonthId {
					break done
				}

				tmpFromDate = tmpFromDate.AddDate(0, 1, 0)
			}
		}

	}

	return
}

func (m *DashboardController) GetDownTimeTopDetail(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	p := new(PayloadDashboard)
	e := k.GetPayload(&p)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	var pipes []tk.M
	var fromDate time.Time

	fromDate = p.Date.AddDate(0, -12, 0)
	pipes = append(pipes, tk.M{"$match": tk.M{"turbine": p.Turbine, "startdate": tk.M{"$gte": fromDate.UTC(), "$lte": p.Date.UTC()}}})
	if p.Type == "Hours" {
		pipes = append(pipes,
			tk.M{
				"$group": tk.M{"_id": tk.M{"id1": "$startdateinfo.monthid", "id2": "$startdateinfo.monthdesc"},
					"result": tk.M{"$sum": "$duration"},
				},
			},
		)
	} else if p.Type == "Times" {
		pipes = append(pipes,
			tk.M{
				"$group": tk.M{"_id": tk.M{"id1": "$startdateinfo.monthid", "id2": "$startdateinfo.monthdesc"},
					"result": tk.M{"$sum": 1},
				},
			},
		)
	} else if p.Type == "MWh" {
		pipes = append(pipes,
			tk.M{
				"$group": tk.M{"_id": tk.M{"id1": "$startdateinfo.monthid", "id2": "$startdateinfo.monthdesc"},
					"result": tk.M{"$sum": "$powerlost"},
				},
			},
		)
	}

	pipes = append(pipes, tk.M{"$sort": tk.M{"_id.id1": 1}})

	csr, e := DB().Connection.NewQuery().
		From(new(Alarm).TableName()).
		Command("pipe", pipes).
		Cursor(nil)

	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	result := []tk.M{}
	tmpResult := []tk.M{}
	e = csr.Fetch(&tmpResult, 0, false)

	// add by ams, 2016-10-07
	csr.Close()

	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	dt, _ := time.Parse("2006-01-02 15:04:05", fromDate.UTC().Format("2006-01")+"-01 00:00:00")

	for i := 1; i < 13; i++ {
		currDate := dt.AddDate(0, i, 0)
		dateInfo := GetDateInfo(currDate)
		found := false

		for _, val := range tmpResult {
			id := val.Get("_id").(tk.M)
			id1 := id.GetInt("id1")

			if id1 == dateInfo.MonthId {
				found = true
				result = append(result, val)
			}
		}

		if !found {
			emptyRes := tk.M{}
			emptyRes.Set("_id", tk.M{"id1": dateInfo.MonthId, "id2": dateInfo.MonthDesc})
			emptyRes.Set("result", 0)
			result = append(result, emptyRes)
		}
	}

	return helper.CreateResult(true, result, "success")
}

func (m *DashboardController) GetDownTimeTurbines(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	p := new(PayloadDashboard)
	e := k.GetPayload(&p)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	var fromDate time.Time
	var pipes []tk.M

	fromDate = p.Date.AddDate(0, 0, -1)

	pipes = append(pipes, tk.M{"$match": tk.M{"startdate": tk.M{"$gte": fromDate.UTC(), "$lte": p.Date.UTC()}}})
	pipes = append(pipes, tk.M{"$group": tk.M{"_id": "$turbine", "result": tk.M{"$sum": "$duration"}}})
	pipes = append(pipes, tk.M{"$sort": tk.M{"_id": 1}})

	csr, e := DB().Connection.NewQuery().
		From(new(Alarm).TableName()).
		Command("pipe", pipes).
		Cursor(nil)

	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	tmpResult := []tk.M{}
	result := []tk.M{}

	e = csr.Fetch(&tmpResult, 0, false)
	// add by ams, 2016-10-07
	csr.Close()

	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	for _, val := range tmpResult {
		val.Set("isdown", false)
		if val.GetFloat64("result") > 24 {
			val.Set("isdown", true)
		}

		result = append(result, val)
	}

	return helper.CreateResult(true, result, "success")
}

func getMapCol(colname string) tk.Ms {
	csr, e := DB().Connection.NewQuery().From(colname).Cursor(nil)
	if e != nil {
		return nil
	}
	defer csr.Close()

	data := tk.Ms{}
	e = csr.Fetch(&data, 0, false)
	if e != nil {
		return nil
	}

	return data
}

func (m *DashboardController) GetMapData(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	payload := map[string]string{}
	e := k.GetPayload(&payload)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}
	colname := ""

	if payload["projectname"] == "Fleet" {
		colname = "ref_project"
	} else {
		colname = "ref_turbine"
	}
	data := getMapCol(colname)

	results := tk.Ms{}
	offset := []int{0, 2}
	coords := []float64{}
	for _, val := range data {
		result := tk.M{}
		coords = []float64{}
		coords = []float64{val.GetFloat64("latitude"), val.GetFloat64("longitude")}
		if payload["projectname"] == "Fleet" {
			result.Set("name", val.GetString("projectname"))
		} else {
			result.Set("name", val.GetString("turbinename"))
		}
		result.Set("coords", coords)
		result.Set("status", "closed")
		result.Set("offset", offset)
		results = append(results, result)
	}

	return helper.CreateResult(true, results, "success")
}

func (m *DashboardController) GetDownTimeLostEnergyDetail(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	p := new(PayloadDashboard)
	e := k.GetPayload(&p)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	var pipes []tk.M
	var date time.Time
	result := []tk.M{}

	dateStr := strings.Split(p.DateStr, " ")

	if dateStr[0] != "fleet" {
		date, e = time.Parse("Jan 2006 02 15:04:05", dateStr[0][0:3]+" "+dateStr[1]+" 01 00:00:00")

		if e != nil {
			return helper.CreateResult(false, nil, e.Error())
		}

		dateInfo := GetDateInfo(date)

		if p.Type != "" {
			pipes = append(pipes, tk.M{"$match": tk.M{"startdateinfo.monthid": dateInfo.MonthId, strings.ToLower(strings.Replace(p.Type, " ", "", 1)): true}})
		} else {
			pipes = append(pipes, tk.M{"$match": tk.M{"startdateinfo.monthid": dateInfo.MonthId}})
		}

	} else {
		dateStr = strings.Split("Jul 2015", " ")
		dateStr2 := strings.Split("Jun 2016", " ")
		date, e = time.Parse("Jan 2006 02 15:04:05", dateStr[0][0:3]+" "+dateStr[1]+" 01 00:00:00")
		date2, e := time.Parse("Jan 2006 02 15:04:05", dateStr2[0][0:3]+" "+dateStr2[1]+" 01 00:00:00")

		if e != nil {
			return helper.CreateResult(false, nil, e.Error())
		}

		if p.Type != "" {
			pipes = append(pipes, tk.M{"$match": tk.M{"startdateinfo.dateid": tk.M{"$gte": date, "$lte": date2}, strings.ToLower(strings.Replace(p.Type, " ", "", 1)): true}})
		} else {
			pipes = append(pipes, tk.M{"$match": tk.M{"startdateinfo.dateid": tk.M{"$gte": date, "$lte": date2}}})
		}
	}

	pipes = append(pipes,
		tk.M{
			"$group": tk.M{
				"_id":       "$turbine",
				"powerlost": tk.M{"$sum": "$powerlost"},
				"duration":  tk.M{"$sum": "$duration"},
			},
		},
	)

	pipes = append(pipes, tk.M{"$sort": tk.M{"powerlost": -1}})

	csr, e := DB().Connection.NewQuery().
		From(new(Alarm).TableName()).
		Command("pipe", pipes).
		Cursor(nil)

	// add by ams, 2016-10-07
	defer csr.Close()
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	e = csr.Fetch(&result, 0, false)

	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}
	return helper.CreateResult(true, result, "success")
}

func (m *DashboardController) GetDownTimeLostEnergyDetailTable(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	p := new(PayloadDashboard)
	e := k.GetPayload(&p)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	result := tk.M{}

	dateStr := []string{}
	date := time.Time{}
	date2 := time.Time{}
	if p.DateStr != "fleet date" {
		dateStr = strings.Split(p.DateStr, " ")
		date, e = time.Parse("Jan 2006 02 15:04:05", dateStr[0][0:3]+" "+dateStr[1]+" 01 00:00:00")
	} else {
		dateStr = strings.Split("Jul 2015", " ")
		dateStr2 := strings.Split("Jun 2016", " ")
		date, e = time.Parse("Jan 2006 02 15:04:05", dateStr[0][0:3]+" "+dateStr[1]+" 01 00:00:00")
		date2, e = time.Parse("Jan 2006 02 15:04:05", dateStr2[0][0:3]+" "+dateStr2[1]+" 01 00:00:00")
	}

	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	dateInfo := DateInfo{}
	dateInfo2 := DateInfo{}
	if p.DateStr != "fleet date" {
		dateInfo = GetDateInfo(date)
	} else {
		dateInfo = GetDateInfo(date)
		dateInfo2 = GetDateInfo(date2)
	}

	var filter []*dbox.Filter

	if p.DateStr != "fleet date" {
		filter = append(filter, dbox.Eq("startdateinfo.monthid", dateInfo.MonthId))
	} else {
		filter = append(filter, dbox.Gte("startdateinfo.monthid", dateInfo.MonthId))
		filter = append(filter, dbox.Lte("startdateinfo.monthid", dateInfo2.MonthId))
		// tk.Println(dateInfo.MonthId)
		// tk.Println(dateInfo2.MonthId)
	}
	if p.ProjectName != "Fleet" {
		if p.Type != "" && p.Type != "All Types" {
			filter = append(filter, dbox.Eq(strings.ToLower(strings.Replace(p.Type, " ", "", 1)), true))
		}
	}

	query := DB().Connection.NewQuery().
		From(new(Alarm).TableName()).Where(filter...).
		Skip(p.Skip).Take(p.Take)

	if len(p.Sort) > 0 {
		var arrsort []string
		for _, val := range p.Sort {
			if val.Dir == "desc" {
				arrsort = append(arrsort, strings.ToLower("-"+val.Field))
			} else {
				arrsort = append(arrsort, strings.ToLower(val.Field))
			}
		}
		query = query.Order(arrsort...)
	}
	csr, e := query.Cursor(nil)

	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	resTable := []Alarm{}
	e = csr.Fetch(&resTable, 0, false)

	// add by ams, 2016-10-07
	csr.Close()

	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	csr, e = DB().Connection.NewQuery().
		From(new(Alarm).TableName()).Where(filter...).
		Cursor(nil)

	// add by ams, 2016-10-07
	defer csr.Close()

	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	total := csr.Count()

	result.Set("data", resTable)
	result.Set("total", total)

	return helper.CreateResult(true, result, "success")
}

func (m *DashboardController) GetDownTimeTopDetailTable(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	p := new(PayloadDashboard)
	e := k.GetPayload(&p)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	var fromDate time.Time
	fromDate = p.Date.AddDate(0, -12, 0)

	var filter []*dbox.Filter

	filter = append(filter, dbox.Eq("turbine", p.Turbine))
	filter = append(filter, dbox.Gte("startdate", fromDate.UTC()))
	filter = append(filter, dbox.Lte("startdate", p.Date.UTC()))

	if p.ProjectName != "Fleet" {
		filter = append(filter, dbox.Lte("projectname", p.ProjectName))
	}

	query := DB().Connection.NewQuery().
		From(new(Alarm).TableName()).Where(filter...).
		Skip(p.Skip).Take(p.Take)

	if len(p.Sort) > 0 {
		var arrsort []string
		for _, val := range p.Sort {
			if val.Dir == "desc" {
				arrsort = append(arrsort, strings.ToLower("-"+val.Field))
			} else {
				arrsort = append(arrsort, strings.ToLower(val.Field))
			}
		}
		query = query.Order(arrsort...)
	}
	csr, e := query.Cursor(nil)

	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	result := tk.M{}
	resTable := []Alarm{}
	e = csr.Fetch(&resTable, 0, false)
	// add by ams, 2016-10-07
	csr.Close()

	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	csr, e = DB().Connection.NewQuery().
		From(new(Alarm).TableName()).Where(filter...).
		Cursor(nil)
	// add by ams, 2016-10-07
	defer csr.Close()

	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	total := csr.Count()

	result.Set("data", resTable)
	result.Set("total", total)

	return helper.CreateResult(true, result, "success")
}

func (m *DashboardController) GetSummaryDataDaily(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	p := new(helper.Payloads)
	e := k.GetPayload(&p)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	var from, to time.Time

	for _, filt := range p.Filter.Filters {
		if filt.Field == "dateinfo.dateid" && filt.Op == "gte" {
			b, err := time.Parse("2006-01-02T15:04:05.000Z", filt.Value.(string))
			t, _ := time.Parse("2006-01-02 15:04:05", b.UTC().Format("2006-01-02")+" 00:00:00")
			if err != nil {
				fmt.Println(err.Error())
			} else {
				from = t
			}
		} else if filt.Field == "dateinfo.dateid" && filt.Op == "lte" {
			b, err := time.Parse("2006-01-02T15:04:05.000Z", filt.Value.(string))
			t, _ := time.Parse("2006-01-02 15:04:05", b.UTC().Format("2006-01-02")+" 00:00:00")
			if err != nil {
				fmt.Println(err.Error())
			} else {
				to = t
			}
		}
	}

	totalHours := tk.ToFloat64(to.Sub(from).Hours(), 0, tk.RoundingUp)

	filter, _ := p.ParseFilter()
	fb := DB().Connection.Fb()
	fb.AddFilter(dbox.And(filter...))
	matches, e := fb.Build()
	if e != nil {
		helper.CreateResult(false, nil, e.Error())
	}

	var periodType string
	_ = periodType
	tmp := []tk.M{}
	projectName := ""

	for _, val := range matches.(tk.M).Get("$and").([]interface{}) {
		if val.(tk.M).Get("type") != nil {
			periodType = val.(tk.M).GetString("type")
		} else if val.(tk.M).Get("projectname") != nil {
			projectName = val.(tk.M).GetString("projectname")
			tmp = append(tmp, val.(tk.M))
		} else {
			tmp = append(tmp, val.(tk.M))
		}
	}

	matches = tk.M{"$and": tmp}
	group := tk.M{}
	if projectName != "" {
		group.Set("_id", tk.M{"id1": "$turbine", "id2": "$projectname"})
	} else {
		group.Set("_id", "$projectname")
	}

	group.Set("production", tk.M{"$sum": "$production"})
	group.Set("plf", tk.M{"$avg": "$plf"})
	group.Set("totalavail", tk.M{"$avg": "$totalavail"})
	group.Set("machineavail", tk.M{"$avg": "$machineavail"})
	group.Set("lowestmachineavail", tk.M{"$min": "$machineavail"})
	group.Set("lowestplf", tk.M{"$min": "$plf"})
	group.Set("maxlossenergy", tk.M{"$max": "$lostenergy"})
	group.Set("gridavail", tk.M{"$avg": "$gridavail"})
	group.Set("totalavail", tk.M{"$avg": "$totalavail"})

	pipe := []tk.M{
		{"$match": matches},
		{"$group": group},
		{"$sort": tk.M{"_id": 1}},
		{"$skip": p.Skip},
		{"$limit": p.Take},
	}

	// tk.Printf("%v\n", pipe)

	pipeCount := []tk.M{
		{"$match": matches},
		{"$group": group},
	}

	csr, e := DB().Connection.NewQuery().
		From(new(ScadaSummaryDaily).TableName()).
		Command("pipe", pipe).
		Cursor(nil)

	if e != nil {
		helper.CreateResult(false, nil, e.Error())
	}
	defer csr.Close()

	result := []tk.M{}
	e = csr.Fetch(&result, 0, false)

	totalTurbine := tk.M{}

	if projectName == "" {

		groupTotalTurbine := tk.M{
			"_id":   tk.M{"id1": "$turbine", "id2": "$projectname"},
			"count": tk.M{"$sum": 1},
		}

		pipeTotalTurbine := []tk.M{
			{"$match": matches},
			{"$group": groupTotalTurbine},
		}

		csrTotalTurbine, e := DB().Connection.NewQuery().
			From(new(ScadaSummaryDaily).TableName()).
			Command("pipe", pipeTotalTurbine).
			Cursor(nil)

		if e != nil {
			return helper.CreateResult(false, nil, e.Error())
		}

		defer csrTotalTurbine.Close()

		resultTotalTurbine := []tk.M{}
		e = csrTotalTurbine.Fetch(&resultTotalTurbine, 0, false)

		if e != nil {
			return helper.CreateResult(false, nil, e.Error())
		}

		for _, val := range resultTotalTurbine {
			var count int
			project := val.Get("_id").(tk.M).GetString("id2")

			if totalTurbine.Get(project) == nil {
				count = 1
			} else {
				count = totalTurbine.GetInt(project) + 1
			}

			totalTurbine.Set(project, count)
		}

	}

	for idx, dt := range result {
		lowestPlf := tk.ToFloat64(result[idx].GetFloat64("lowestplf"), 2, tk.RoundingAuto)
		lowestMachineAvail := tk.ToFloat64(result[idx].GetFloat64("lowestmachineavail"), 2, tk.RoundingAuto)
		maxLossEnergy := tk.ToFloat64(result[idx].GetFloat64("maxlossenergy"), 2, tk.RoundingAuto)

		if projectName != "" {
			result[idx].Set("name", dt.Get("_id").(tk.M).GetString("id1"))
			result[idx].Set("lowestplf", formatStringFloat(tk.ToString(lowestPlf), 2)+" %")
			result[idx].Set("lowestmachineavail", formatStringFloat(tk.ToString(lowestMachineAvail), 2)+" %")
			result[idx].Set("maxlossenergy", maxLossEnergy)

			maxCapacity := turbineMW * totalHours
			result[idx].Set("maxcapacity", maxCapacity)
		} else {
			result[idx].Set("name", dt.GetString("_id"))
			result[idx].Set("noofwtg", totalTurbine.GetInt(result[idx].GetString("name")))

			maxCapacity := turbineMW * totalTurbine.GetFloat64(result[idx].GetString("name")) * totalHours
			result[idx].Set("maxcapacity", maxCapacity)

			// ---- lowestplf

			pipeSub := []tk.M{
				{"$match": matches},
				{"$sort": tk.M{"plf": 1}},
				{"$limit": 1},
			}

			csr, e = DB().Connection.NewQuery().
				From(new(ScadaSummaryDaily).TableName()).
				Command("pipe", pipeSub).
				Cursor(nil)

			if e != nil {
				return helper.CreateResult(false, nil, e.Error())
			}

			defer csr.Close()

			lowest := []tk.M{}
			e = csr.Fetch(&lowest, 0, false)

			if e != nil {
				return helper.CreateResult(false, nil, e.Error())
			}

			if len(lowest) > 0 {
				result[idx].Set("lowestplf", formatStringFloat(tk.ToString(lowestPlf), 2)+"% ("+lowest[0].GetString("turbine")+")")
			}

			// ---- lowestmachineavail

			pipeSub = []tk.M{
				{"$match": matches},
				{"$sort": tk.M{"machineavail": 1}},
				{"$limit": 1},
			}

			csr, e = DB().Connection.NewQuery().
				From(new(ScadaSummaryDaily).TableName()).
				Command("pipe", pipeSub).
				Cursor(nil)

			if e != nil {
				return helper.CreateResult(false, nil, e.Error())
			}

			defer csr.Close()

			lowest = []tk.M{}
			e = csr.Fetch(&lowest, 0, false)

			if e != nil {
				return helper.CreateResult(false, nil, e.Error())
			}

			if len(lowest) > 0 {
				result[idx].Set("lowestmachineavail", formatStringFloat(tk.ToString(lowestMachineAvail), 2)+"% ("+lowest[0].GetString("turbine")+")")
			}

			// ---- maxlossenergy

			pipeSub = []tk.M{
				{"$match": matches},
				{"$sort": tk.M{"lostenergy": -1}},
				{"$limit": 1},
			}

			csr, e = DB().Connection.NewQuery().
				From(new(ScadaSummaryDaily).TableName()).
				Command("pipe", pipeSub).
				Cursor(nil)

			if e != nil {
				return helper.CreateResult(false, nil, e.Error())
			}

			defer csr.Close()

			highest := []tk.M{}
			e = csr.Fetch(&highest, 0, false)

			if e != nil {
				return helper.CreateResult(false, nil, e.Error())
			}

			if len(highest) > 0 {
				result[idx].Set("maxlossenergy", formatStringFloat(tk.ToString(maxLossEnergy), 2)+" ("+highest[0].GetString("turbine")+")")
			}
		}
	}

	csrCount, e := DB().Connection.NewQuery().
		From(new(ScadaSummaryDaily).TableName()).
		Command("pipe", pipeCount).
		Cursor(nil)
	defer csrCount.Close()

	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	rs := []tk.M{}
	e = csrCount.Fetch(&rs, 0, false)

	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	data := struct {
		Data  []tk.M
		Total int
	}{
		Data:  result,
		Total: tk.SliceLen(rs),
	}

	return helper.CreateResult(true, data, "success")
}

func formatStringFloat(str string, decimalPoint int) (result string) {
	anStr := strings.Split(str, ".")
	if len(anStr) > 0 {
		result = anStr[0] + "." + anStr[1][:decimalPoint]
	} else {
		result = str
	}
	return
}
