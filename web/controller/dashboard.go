package controller

import (
	. "eaciit/wfdemo-git/library/core"
	. "eaciit/wfdemo-git/library/helper"
	. "eaciit/wfdemo-git/library/models"
	"eaciit/wfdemo-git/web/helper"
	"log"
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
	IsDetail    bool
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
		// "aebok":        "AEBOK",
		// "externalstop": "External Stop",
		"griddown": "Grid Down",
		// "internalgrid": "Internal Grid",
		"machinedown": "Machine Down",
		"unknown":     "Unknown",
		// "weatherstop":  "Weather Stop",
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

	csr, e := DB().Connection.NewQuery().From(new(ScadaLastUpdate).TableName()).Where(dbox.And(dbox.Eq("projectname", p.ProjectName))).Cursor(nil)

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

		turbineDownOneDays := getDownTurbine(val.ProjectName, val.LastUpdate, 1)
		turbineDownTwoDays := getDownTurbine(val.ProjectName, val.LastUpdate, 2)

		val.CurrentDown = len(turbineDownOneDays)
		val.TwoDaysDown = len(turbineDownTwoDays)

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
	csrlastmonth, e := DB().Connection.NewQuery().
		From("rpt_scadalastupdate").
		Where(dbox.And(dbox.Eq("projectname", p.ProjectName))).
		Cursor(nil)

	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}
	defer csrlastmonth.Close()

	datalastmonth := make([]ScadaSummaryByMonth, 0)
	e = csrlastmonth.Fetch(&datalastmonth, 0, false)

	var result []interface{}

	if len(datalastmonth) > 0 {

		startmonth := 0
		endmonth := datalastmonth[0].DateInfo.MonthId
		month := endmonth - (int(endmonth/100) * 100)

		if month == 12 {
			startmonth = (int(endmonth/100) * 100) + 1
		} else {
			startmonth = (endmonth + 1) - 100
		}

		// result := make([]ScadaSummaryByMonth, 0)

		dataVariance := new(ScadaSummaryVariance)

		for i := startmonth; i <= endmonth; i++ {
			//check if month more than 12
			if i-(int(i/100)*100) > 12 {
				i = (i - 12) + 100
			}

			yearloop := int(i / 100)
			monthloop := i - (int(i/100) * 100)

			csr, e := DB().Connection.NewQuery().
				From("rpt_scadasummarybymonth").
				Where(dbox.And(dbox.Eq("dateinfo.monthid", i), dbox.Eq("projectname", p.ProjectName))).
				Cursor(nil)

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
	matches := tk.M{"dateinfo.monthdesc": p.GetString("date"), "available": 1}
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

	p := tk.M{}
	e := k.GetPayload(&p)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	pipe := []tk.M{}
	pipe = append(pipe, tk.M{"$unwind": "$dataitems"})
	// if p.GetString("project") == "Fleet" {
	// 	pipe = append(pipe, tk.M{"$group": tk.M{
	// 		"_id":           "$_id",
	// 		"noofwtg":       tk.M{"$sum": "$dataitems.noofwtg"},
	// 		"production":    tk.M{"$sum": "$dataitems.production"},
	// 		"plf":           tk.M{"$avg": "$dataitems.plf"},
	// 		"lostenergy":    tk.M{"$sum": "$dataitems.lostenergy"},
	// 		"downtimehours": tk.M{"$sum": "$dataitems.downtimehours"},
	// 		"machineavail":  tk.M{"$avg": "$dataitems.machineavail"},
	// 		"trueavail":     tk.M{"$avg": "$dataitems.trueavail"},
	// 	}})
	// 	pipe = append(pipe, tk.M{"$sort": tk.M{"_id": 1}})
	// } else {
	// 	pipe = append(pipe, tk.M{"$match": tk.M{"_id": p.GetString("project")}})
	// 	pipe = append(pipe, tk.M{"$sort": tk.M{"dataitems.name": 1}})
	// }
	pipe = append(pipe, tk.M{"$match": tk.M{"_id": p.GetString("project")}})
	pipe = append(pipe, tk.M{"$sort": tk.M{"dataitems.name": 1}})
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

	// if p.GetString("project") == "Fleet" {
	// 	for _, val := range result {
	// 		val.Set("name", val.GetString("_id"))
	// 		dataItem = append(dataItem, val)
	// 	}
	// } else {
	// 	for _, val := range result {
	// 		dataItem = append(dataItem, val["dataitems"].(tk.M))
	// 	}
	// }
	for _, val := range result {
		dataItem = append(dataItem, val["dataitems"].(tk.M))
	}

	data := struct {
		Data  []tk.M
		Total int
	}{
		Data:  dataItem,
		Total: tk.SliceLen(dataItem),
	}

	return helper.CreateResult(true, data, "success")
}

func (m *DashboardController) GetDownTimeLoss(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	p := new(PayloadAnalytic)
	e := k.GetPayload(&p)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	machinedown, _ := getMachineDownType()
	projectList := []string{}
	if p.Project == "" {
		projectList, e = getProject()
		if e != nil {
			return helper.CreateResult(false, nil, e.Error())
		}
	} else {
		projectList = []string{p.Project}
	}

	tStart, tEnd, e := helper.GetStartEndDate(k, p.Period, p.DateStart, p.DateEnd)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}
	var pipes []tk.M
	match := tk.M{}

	match.Set("detail.startdate", tk.M{"$gte": tStart, "$lte": tEnd})

	if len(p.Turbine) != 0 {
		match.Set("turbine", tk.M{"$in": p.Turbine})
	}
	result := []tk.M{}

	for _, project := range projectList {
		for field, mdName := range machinedown {
			pipes = []tk.M{}
			match.Set("projectname", project)
			match.Set("detail."+field, true)
			pipes = append(pipes, tk.M{"$unwind": "$detail"})
			pipes = append(pipes, tk.M{"$match": match})
			groups := tk.M{
				"_id": tk.M{
					"id1": mdName,
					"id2": mdName,
					"id3": project,
				},
				"powerlost": tk.M{"$sum": "$detail.powerlost"},
				"duration":  tk.M{"$sum": "$detail.duration"},
				"frequency": tk.M{"$sum": 1},
			}
			pipes = append(pipes, tk.M{"$group": groups})

			csr, e := DB().Connection.NewQuery().
				From(new(Alarm).TableName()).
				Command("pipe", pipes).
				Cursor(nil)

			if e != nil {
				return helper.CreateResult(false, nil, e.Error())
			}

			tmpRes := []tk.M{}
			e = csr.Fetch(&tmpRes, 0, false)
			if e != nil {
				return helper.CreateResult(false, nil, e.Error())
			}
			csr.Close()

			found := false
			if tk.SliceLen(tmpRes) > 0 {
				found = true
				tmpRes[0]["powerlost"] = tmpRes[0].GetFloat64("powerlost") / 1000
				result = append(result, tmpRes[0])
			}

			if !found {
				emptyRes := tk.M{}
				emptyRes.Set("_id", tk.M{"id1": mdName, "id2": mdName, "id3": tk.ToString(project)})
				emptyRes.Set("powerlost", 0)
				emptyRes.Set("duration", 0)
				emptyRes.Set("frequency", 0)

				result = append(result, emptyRes)
			}
			match.Unset("detail." + field)
		}
	}
	return helper.CreateResult(true, result, "success")
}

func (m *DashboardController) GetLostEnergy(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	result := tk.M{}

	p := new(PayloadDashboard)
	e := k.GetPayload(&p)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}
	downtimeDatas := getDownTimeLostEnergy("project", p)
	result.Set("lostenergy", downtimeDatas)

	if !p.IsDetail {
		if p.Type == "" && p.ProjectName == "Fleet" {
			result.Set("lostenergybytype", getDownTimeLostEnergy("type", p))
		}
	}
	return helper.CreateResult(true, result, "success")
}

func (m *DashboardController) GetDowntimeTop(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	result := tk.M{}

	p := new(PayloadDashboard)
	e := k.GetPayload(&p)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	if !p.IsDetail {
		if p.ProjectName != "Fleet" {
			//tidak bisa dicombine karena tiap top 10 kategori beda urutan top 10 nya
			result.Set("duration", getTurbineDownTimeTop("duration", p))
			result.Set("frequency", getTurbineDownTimeTop("frequency", p))
		}
		result.Set("loss", getTurbineDownTimeTop("loss", p))
	}
	return helper.CreateResult(true, result, "success")
}

func (m *DashboardController) GetLossCategories(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	result := tk.M{}

	p := new(PayloadDashboard)
	e := k.GetPayload(&p)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	if !p.IsDetail {
		lossD, lossF, loss := getLossCategoriesTopDFP(p)
		result.Set("lossCatDuration", lossD)
		result.Set("lossCatFrequency", lossF)
		result.Set("lossCatLoss", loss)
	}
	return helper.CreateResult(true, result, "success")
}

func (m *DashboardController) GetMachGridAvailability(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	result := tk.M{}

	p := new(PayloadDashboard)
	e := k.GetPayload(&p)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	if !p.IsDetail {
		machAvail, gridAvail := getMGAvailability(p)

		result.Set("machineAvailability", machAvail)
		result.Set("gridAvailability", gridAvail)
	}
	return helper.CreateResult(true, result, "success")
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
	result.Set("lostenergy", downtimeDatas)

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
	machinedown, e := getMachineDownType()
	var typeX string

	for i, v := range machinedown {
		if v == p.Type {
			typeX = i
			break
		}
	}

	if p.DateStr != "" {
		dateStr := strings.Split(p.DateStr, " ")
		date, e := time.Parse("Jan 2006 02 15:04:05", dateStr[0][0:3]+" "+dateStr[1]+" 01 00:00:00")
		if e != nil {
			return
		}

		dateInfo := GetDateInfo(date)

		if tipe == "fleetdowntime" {
			matchDown.Set("detail.detaildateinfo.monthid", dateInfo.MonthId)
		} else {
			match.Set("dateinfo.monthid", dateInfo.MonthId)
		}
	} else {
		fromDate = p.Date.AddDate(0, -12, 0)
		match.Set("detail.detaildateinfo.dateid", tk.M{"$gte": fromDate.UTC(), "$lte": p.Date})
		/*tk.Println("From Date: ", fromDate)
		tk.Println("PayLoad Date: ", p.Date)*/
	}

	if p.ProjectName != "Fleet" {
		match.Set("projectname", p.ProjectName)
		matchDown.Set("projectname", p.ProjectName)
	}

	if p.Type != "" && tipe != "fleetdowntime" && p.Type != "All Types" {
		match.Set(typeX, true)
	} else if p.Type != "" && tipe == "fleetdowntime" {
		matchDown.Set(strings.Replace(strings.ToLower(p.Type), " ", "", 1), true)
	}

	// pipes = append(pipes, tk.M{"$match": match})

	if p.ProjectName != "Fleet" {
		// add a condition to check the type is project
		// regarding to next process can not catch the value for selecting downtime by project in dashboard
		// add by ams on 20161003
		if tipe == "project" {
			if p.Type == "All Types" {
				pipeIds := tk.M{
					"id1": "$detail.detaildateinfo.monthid",
					"id2": "$detail.detaildateinfo.monthdesc",
					"id3": p.Type,
				}

				for mcd := range machinedown {
					pipeIds.Set(mcd, "$"+mcd)
				}

				pipes = append(pipes, tk.M{"$unwind": "$detail"})
				pipes = append(pipes, tk.M{"$match": match})
				pipes = append(pipes,
					tk.M{
						"$group": tk.M{
							"_id":    pipeIds,
							"result": tk.M{"$sum": "$detail.powerlost"},
						},
					},
				)
			} else {
				// pipes = append(pipes,
				// 	tk.M{
				// 		/*"$group": tk.M{"_id": tk.M{"id1": "$dateinfo.monthid", "id2": "$dateinfo.monthdesc", "id3": "$projectname"},
				// 		"result": tk.M{"$sum": "$lostenergy"},*/
				// 		"$group": tk.M{"_id": tk.M{"id1": "$dateinfo.monthid", "id2": "$dateinfo.monthdesc", "id3": "$type"},
				// 			"result": tk.M{"$sum": "$lostenergy"}, /*changed from by project to by MD type per 11 Oct 16 [RS]*/
				// 		},
				// 	},
				// )

				pipes = append(pipes, tk.M{"$unwind": "$detail"})
				pipes = append(pipes, tk.M{"$match": match})
				pipes = append(pipes,
					tk.M{
						"$group": tk.M{"_id": tk.M{
							"id1": "$detail.detaildateinfo.monthid",
							"id2": "$detail.detaildateinfo.monthdesc",
							"id3": typeX},
							"result": tk.M{"$sum": "$detail.powerlost"},
						},
					},
				)
			}

		} else {
			pipes = append(pipes, tk.M{"$match": match})
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
			pipes = append(pipes, tk.M{"$unwind": "$detail"})
			pipes = append(pipes, tk.M{"$match": match})
			pipes = append(pipes,
				tk.M{
					"$group": tk.M{"_id": tk.M{"id1": "$detail.detaildateinfo.monthid", "id2": "$detail.detaildateinfo.monthdesc", "id3": "$projectname"},
						"result": tk.M{"$sum": "$detail.powerlost"},
					},
				},
			)
		} else {
			/*pipes = append(pipes,
				tk.M{
					"$group": tk.M{"_id": tk.M{"id1": "$type", "id2": "$type", "id3": "$projectname"},
						"powerlost": tk.M{"$sum": "$powerlost"},
					},
				},
			)*/

			pipeIds := tk.M{
				"id1": "tipe",
				"id2": "tipe",
				"id3": "$projectname",
			}

			for mcd := range machinedown {
				pipeIds.Set(mcd, "$"+mcd)
			}

			pipes = append(pipes, tk.M{"$unwind": "$detail"})
			pipes = append(pipes, tk.M{"$match": match})
			pipes = append(pipes,
				tk.M{
					"$group": tk.M{
						"_id":       pipeIds,
						"powerlost": tk.M{"$sum": "$detail.powerlost"},
						"duration":  tk.M{"$sum": "$detail.duration"},
						"frequency": tk.M{"$sum": 1},
					},
				},
			)
		}
	}

	if e != nil {
		return nil
	}

	if p.DateStr == "" && tipe != "fleetdowntime" {
		pipes = append(pipes, tk.M{"$sort": tk.M{"_id.id3": 1}})

		/*for _, pip := range pipes {
			log.Printf("%#v \n", pip)
		}*/

		csr, e := DB().Connection.NewQuery().
			From(new(Alarm).TableName()).
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
		// lineData := tk.M{}

		for field, title := range stack {
			if tipe != "type" {
				/*for _, val := range tmpResult {
					log.Printf("val: %#v \n", val)
				}*/

				for i := 1; i < 13; i++ {
					currDate := dt.AddDate(0, i, 0)
					dateInfo := GetDateInfo(currDate)
					found := false

					for _, val := range tmpResult {
						id := val.Get("_id").(tk.M)
						id1 := id.GetInt("id1")
						id3 := id.GetString("id3")
						if p.Type == "All Types" && id.Get(field) != nil {
							// log.Printf("id: %#v || %v \n", id, field)
							idDown := id.Get(field).(bool)
							if id1 == dateInfo.MonthId && idDown {
								found = true
								val.Set("_id", tk.M{
									"id1": id1,
									"id2": dateInfo.MonthDesc,
									"id3": tk.ToString(title),
								})
								val.Set("result", val.GetFloat64("result")*0.001)
								result = append(result, val)
								break
							}

						} else {
							if id1 == dateInfo.MonthId && (id3 == tk.ToString(title) || id3 == tk.ToString(field)) {
								found = true
								val.Set("_id", tk.M{
									"id1": id1,
									"id2": dateInfo.MonthDesc,
									"id3": tk.ToString(title),
								})
								val.Set("result", val.GetFloat64("result")*0.001)
								result = append(result, val)
								break
							}
						}
						// tk.Printf("ID 1 => %#v\n", id1)
						// tk.Printf("MonthId => %#v\n", dateInfo.MonthId)
						// tk.Printf("ID 3 => %#v\n", id3)
						// tk.Printf("Title => %#v\n", tk.ToString(title))
						// tk.Printf("Value => %#v\n", val.GetFloat64("result"))

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
				var bigPower, bigDuration float64
				var bigFreq int

				for field, mdName := range machinedown {
					found := false
					for _, val := range tmpResult {
						id := val.Get("_id").(tk.M)
						id3 := id.GetString("id3")

						found = id.Get(field).(bool)
						if found && id3 == tk.ToString(title) {
							powerlost := val.GetFloat64("powerlost") * 0.001
							duration := val.GetFloat64("duration")
							frequency := val.GetInt("frequency")

							if powerlost > bigPower {
								bigPower = powerlost
							}
							if duration > bigDuration {
								bigDuration = duration
							}
							if frequency > bigFreq {
								bigFreq = frequency
							}

							foundRes := tk.M{}
							foundRes.Set("_id", tk.M{"id1": mdName, "id2": mdName, "id3": tk.ToString(title)})
							foundRes.Set("powerlost", powerlost)
							foundRes.Set("duration", duration)
							foundRes.Set("frequency", frequency)

							source = append(source, foundRes)
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
			pipesDown = []tk.M{}
			pipesDown = append(pipesDown, tk.M{"$unwind": "$detail"})
			pipesDown = append(pipesDown, tk.M{"$match": matchDown})
			pipesX := pipesDown
			pipesX = append(pipesX,
				tk.M{
					"$group": tk.M{"_id": tk.M{"id1": "$startdateinfo.monthid", "id2": "$startdateinfo.monthdesc", "id3": p.Type},
						"result": tk.M{"$sum": "$detail.powerlost"},
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

				pipesDown = []tk.M{}
				pipesDown = append(pipesDown, tk.M{"$unwind": "$detail"})
				pipesDown = append(pipesDown, tk.M{"$match": matchX})
				pipesX := pipesDown
				pipesX = append(pipesX,
					tk.M{
						"$group": tk.M{"_id": tk.M{"id1": "$startdateinfo.monthid", "id2": "$startdateinfo.monthdesc", "id3": field},
							"result": tk.M{"$sum": "$detail.powerlost"},
						},
					},
				)

				// tk.Printf("\n%#v \n\n", pipesX)

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

	if p.Type != "All Types" && p.ProjectName != "Fleet" && tipe != "fleetdowntime" {
		hasil := []tk.M{}
		ids := tk.M{}
		for _, val := range result {
			ids, _ = tk.ToM(val["_id"])
			if ids.GetString("id3") == p.Type {
				hasil = append(hasil, val)
			}
		}
		result = hasil
	} else if p.Type != "All Types" && p.Type != "" && tipe == "fleetdowntime" {
		hasil := []tk.M{}
		for _, val := range result {
			if val.GetString("type") == p.Type {
				hasil = append(hasil, val)
			}
		}
		result = hasil
	}

	return
}

func getTurbineDownTimeTop(topType string, p *PayloadDashboard) (result []tk.M) {
	var pipes []tk.M
	var fromDate time.Time
	match := tk.M{}

	if p.DateStr == "" {
		fromDate = p.Date.AddDate(0, -12, 0)

		match.Set("detail.startdate", tk.M{"$gte": fromDate.UTC(), "$lte": p.Date.UTC()})

		if p.ProjectName != "Fleet" {
			match.Set("projectname", p.ProjectName)
		}

		pipes = append(pipes, tk.M{"$unwind": "$detail"})
		pipes = append(pipes, tk.M{"$match": match})

		if topType == "duration" {
			pipes = append(pipes, tk.M{"$group": tk.M{"_id": "$turbine", "result": tk.M{"$sum": "$detail.duration"}}})
		} else if topType == "frequency" {
			pipes = append(pipes, tk.M{"$group": tk.M{"_id": "$turbine", "result": tk.M{"$sum": 1}}})
		} else if topType == "loss" {
			pipes = append(pipes, tk.M{"$group": tk.M{"_id": "$turbine", "result": tk.M{"$sum": "$detail.powerlost"}}})
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
			turbines = append(turbines, turbine.Get("_id").(string))                   /*untuk turbine list*/
			turbinesVal.Set(turbine.Get("_id").(string), turbine.GetFloat64("result")) /*untuk total list tiap turbine*/
		}

		match.Set("turbine", tk.M{"$in": turbines})

		downCause := tk.M{}
		// downCause.Set("aebok", "AEBOK")
		// downCause.Set("externalstop", "External Stop")
		downCause.Set("griddown", "Grid Down")
		// downCause.Set("internalgrid", "Internal Grid")
		downCause.Set("machinedown", "Machine Down")
		downCause.Set("unknown", "Unknown")
		// downCause.Set("weatherstop", "Weather Stop")

		tmpResult := []tk.M{}
		downDone := []string{}

		for f, t := range downCause {
			pipes = []tk.M{}
			loopMatch := match
			field := tk.ToString(f)
			title := tk.ToString(t)

			downDone = append(downDone, field)

			for _, done := range downDone {
				match.Unset("detail." + done)
			}

			loopMatch.Set("detail."+field, true)

			pipes = append(pipes, tk.M{"$unwind": "$detail"})
			pipes = append(pipes, tk.M{"$match": loopMatch})
			if topType == "duration" {
				pipes = append(pipes,
					tk.M{
						"$group": tk.M{"_id": tk.M{"id3": "$turbine", "id4": title},
							"result": tk.M{"$sum": "$detail.duration"},
						},
					},
				)
			} else if topType == "frequency" {
				pipes = append(pipes,
					tk.M{
						"$group": tk.M{"_id": tk.M{"id3": "$turbine", "id4": title},
							"result": tk.M{"$sum": 1},
						},
					},
				)
			} else if topType == "loss" {
				pipes = append(pipes,
					tk.M{
						"$group": tk.M{"_id": tk.M{"id3": "$turbine", "id4": title},
							"result": tk.M{"$sum": "$detail.powerlost"},
						},
					},
				)
			}

			pipes = append(pipes, tk.M{"$sort": tk.M{"result": -1}})

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

			for _, res := range resLoop {
				tmpResult = append(tmpResult, res)
			}
		}

		resY := []tk.M{}

		for _, t := range downCause {
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

// func getTurbineDownTimeTopAll(p *PayloadDashboard) (duration []tk.M, frequency []tk.M, loss []tk.M) {
// 	var pipes []tk.M
// 	var fromDate time.Time
// 	match := tk.M{}

// 	if p.DateStr == "" {
// 		fromDate = p.Date.AddDate(0, -12, 0)

// 		match.Set("startdate", tk.M{"$gte": fromDate.UTC(), "$lte": p.Date.UTC()})

// 		if p.ProjectName != "Fleet" {
// 			match.Set("projectname", p.ProjectName)
// 		}

// 		pipes = append(pipes, tk.M{"$match": match})

// 		if topType == "duration" {
// 			pipes = append(pipes, tk.M{"$group": tk.M{"_id": "$turbine", "result": tk.M{"$sum": "$duration"}}})
// 		} else if topType == "frequency" {
// 			pipes = append(pipes, tk.M{"$group": tk.M{"_id": "$turbine", "result": tk.M{"$sum": 1}}})
// 		} else if topType == "loss" {
// 			pipes = append(pipes, tk.M{"$group": tk.M{"_id": "$turbine", "result": tk.M{"$sum": "$powerlost"}}})
// 		}

// 		pipes = append(pipes, tk.M{"$group": tk.M{"_id": "$turbine", "result": tk.M{"$sum": "$duration"}}})

// 		pipes = append(pipes, tk.M{"$sort": tk.M{"result": -1}})
// 		pipes = append(pipes, tk.M{"$limit": 10})

// 		// get the top 10

// 		csr, e := DB().Connection.NewQuery().
// 			Select("_id").
// 			From(new(Alarm).TableName()).
// 			Command("pipe", pipes).
// 			Cursor(nil)

// 		if e != nil {
// 			return
// 		}

// 		top10Turbines := []tk.M{}
// 		e = csr.Fetch(&top10Turbines, 0, false)
// 		// add by ams, 2016-10-07
// 		csr.Close()

// 		if e != nil {
// 			return
// 		}

// 		// get the downtime

// 		turbines := []string{}
// 		turbinesVal := tk.M{}

// 		for _, turbine := range top10Turbines {
// 			turbines = append(turbines, turbine.Get("_id").(string))
// 			turbinesVal.Set(turbine.Get("_id").(string), turbine.GetFloat64("result"))
// 		}

// 		// tk.Printf("topType: \n%#v \n", topType)
// 		// tk.Printf("turbines: %#v \n", turbines)

// 		match.Set("turbine", tk.M{"$in": turbines})

// 		down, e := getMachineDownType()
// 		downCause := tk.M{}

// 		for f, v := range down {
// 			downCause.Set(f, v)
// 		}

// 		/*downCause.Set("aebok", "AEBOK")
// 		downCause.Set("externalstop", "External Stop")
// 		downCause.Set("griddown", "Grid Down")
// 		downCause.Set("internalgrid", "Internal Grid")
// 		downCause.Set("machinedown", "Machine Down")
// 		downCause.Set("unknown", "Unknown")
// 		downCause.Set("weatherstop", "Weather Stop")*/

// 		tmpResult := []tk.M{}
// 		downDone := []string{}

// 		for f, t := range downCause {
// 			pipes = []tk.M{}
// 			loopMatch := match
// 			field := tk.ToString(f)
// 			title := tk.ToString(t)

// 			downDone = append(downDone, field)

// 			for _, done := range downDone {
// 				match.Unset(done)
// 			}

// 			loopMatch.Set(field, true)

// 			// tk.Printf("%#v \n", match)

// 			pipes = append(pipes, tk.M{"$match": loopMatch})
// 			if topType == "duration" {
// 				pipes = append(pipes,
// 					tk.M{
// 						"$group": tk.M{"_id": tk.M{"id1": "$dateinfo.monthid", "id2": "$dateinfo.monthdesc", "id3": "$turbine", "id4": title},
// 							"result": tk.M{"$sum": "$duration"},
// 						},
// 					},
// 				)
// 			} else if topType == "frequency" {
// 				pipes = append(pipes,
// 					tk.M{
// 						"$group": tk.M{"_id": tk.M{"id1": "$dateinfo.monthid", "id2": "$dateinfo.monthdesc", "id3": "$turbine", "id4": title},
// 							"result": tk.M{"$sum": 1},
// 						},
// 					},
// 				)
// 			} else if topType == "loss" {
// 				pipes = append(pipes,
// 					tk.M{
// 						"$group": tk.M{"_id": tk.M{"id1": "$dateinfo.monthid", "id2": "$dateinfo.monthdesc", "id3": "$turbine", "id4": title},
// 							"result": tk.M{"$sum": "$powerlost"},
// 						},
// 					},
// 				)
// 			}

// 			pipes = append(pipes, tk.M{"$sort": tk.M{"result": -1}})

// 			/*tk.Println()
// 			tk.Println(tk.ToString(title))
// 			for _, val := range pipes {
// 				tk.Printf("pipes: %v \n", val)
// 			}*/

// 			csr, e := DB().Connection.NewQuery().
// 				From(new(Alarm).TableName()).
// 				Command("pipe", pipes).
// 				Cursor(nil)

// 			if e != nil {
// 				return
// 			}

// 			resLoop := []tk.M{}
// 			e = csr.Fetch(&resLoop, 0, false)

// 			// add by ams, 2016-10-07
// 			csr.Close()

// 			// tk.Printf("resLoop: %v - %#v \n", tk.ToString(title), resLoop)

// 			for _, res := range resLoop {
// 				tmpResult = append(tmpResult, res)
// 			}
// 		}

// 		/*tk.Printf("len: %v \n", len(tmpResult))
// 		tk.Printf("%#v \n", tmpResult)*/

// 		/*for _, val := range tmpResult {
// 			tk.Printf("tmpResult: %v \n", val)
// 		}*/

// 		resY := []tk.M{}

// 		for _, t := range downCause {
// 			// field := tk.ToString(f)
// 			title := tk.ToString(t)

// 			for _, turbine := range turbines {
// 				resX := tk.M{}
// 				resX.Set("_id", tk.M{"id3": turbine, "id4": title})
// 				resX.Set("result", 0)

// 			out:
// 				for _, res := range tmpResult {
// 					id3 := res.Get("_id").(tk.M).GetString("id3")
// 					id4 := res.Get("_id").(tk.M).GetString("id4")

// 					if id3 == turbine && id4 == title {
// 						resX = res
// 						break out
// 					}
// 				}

// 				// if title == "External Stop" {
// 				resY = append(resY, resX)
// 				// }
// 			}
// 		}

// 		for _, turbine := range turbines {
// 			resVal := tk.M{}
// 			resVal.Set("_id", turbine)

// 			for _, val := range resY {
// 				valTurbine := val.Get("_id").(tk.M).GetString("id3")
// 				valResult := val.GetFloat64("result")
// 				valTitle := ""

// 				splitTitle := strings.Split(val.Get("_id").(tk.M).GetString("id4"), " ")

// 				if len(splitTitle) > 1 {
// 					valTitle = splitTitle[0] + "" + splitTitle[1]
// 				} else {
// 					valTitle = splitTitle[0]
// 				}

// 				if turbine == valTurbine && valResult != 0 {
// 					resVal.Set(valTitle, valResult)
// 				} else if resVal.Get(valTitle) == nil {
// 					resVal.Set(valTitle, 0)
// 				}
// 			}

// 			resVal.Set("Total", turbinesVal.GetFloat64(turbine))
// 			result = append(result, resVal)
// 		}
// 	}

// 	return
// }

func getLossCategoriesTopDFP(p *PayloadDashboard) (resultDuration, resultFreq, resultPowerLost []tk.M) {
	var pipes []tk.M
	var fromDate time.Time
	match := tk.M{}

	if p.DateStr == "" {
		fromDate = p.Date.AddDate(0, -12, 0)

		match.Set("detail.startdate", tk.M{"$gte": fromDate.UTC(), "$lte": p.Date.UTC()})

		if p.ProjectName != "Fleet" {
			match.Set("projectname", p.ProjectName)
		}

		downCause := tk.M{}
		// downCause.Set("aebok", "AEBOK")
		// downCause.Set("externalstop", "External Stop")
		downCause.Set("griddown", "Grid Down")
		// downCause.Set("internalgrid", "Internal Grid")
		downCause.Set("machinedown", "Machine Down")
		downCause.Set("unknown", "Unknown")
		// downCause.Set("weatherstop", "Weather Stop")

		tmpResult := []tk.M{}
		tmpResultFreq := []tk.M{}
		tmpResultPower := []tk.M{}
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

			pipes = append(pipes, tk.M{"$unwind": "$detail"})
			pipes = append(pipes, tk.M{"$match": loopMatch})
			pipes = append(pipes,
				tk.M{
					"$group": tk.M{
						"_id":       tk.M{"id1": "detail." + field, "id2": title},
						"duration":  tk.M{"$sum": "$detail.duration"},
						"freq":      tk.M{"$sum": 1},
						"powerlost": tk.M{"$sum": "$detail.powerlost"}},
				},
			)

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
				tmpResultFreq = append(tmpResultFreq, res)
				tmpResultPower = append(tmpResultPower, res)
			}
		}

		size := len(tmpResult)
		sizeF := len(tmpResultFreq)
		sizeP := len(tmpResultPower)

		if size > 1 {
			for i := 0; i < size; i++ {
				for j := size - 1; j >= i+1; j-- {
					a := tmpResult[j].GetFloat64("duration")
					b := tmpResult[j-1].GetFloat64("duration")

					if a > b {
						tmpResult[j], tmpResult[j-1] = tmpResult[j-1], tmpResult[j]
					}
				}
			}
		}

		if sizeF > 1 {
			for i := 0; i < sizeF; i++ {
				for j := sizeF - 1; j >= i+1; j-- {
					a := tmpResultFreq[j].GetInt("freq")
					b := tmpResultFreq[j-1].GetInt("freq")

					if a > b {
						tmpResultFreq[j], tmpResultFreq[j-1] = tmpResultFreq[j-1], tmpResultFreq[j]
					}
				}
			}
		}

		if sizeP > 1 {
			for i := 0; i < size; i++ {
				for j := size - 1; j >= i+1; j-- {
					a := tmpResultPower[j].GetFloat64("powerlost")
					b := tmpResultPower[j-1].GetFloat64("powerlost")

					if a > b {
						tmpResultPower[j], tmpResultPower[j-1] = tmpResultPower[j-1], tmpResultPower[j]
					}
				}
			}
		}

		for _, res := range tmpResult {
			resultDuration = append(resultDuration, tk.M{"_id": res["_id"], "result": res.GetFloat64("duration")})
		}
		for _, res := range tmpResultFreq {
			resultFreq = append(resultFreq, tk.M{"_id": res["_id"], "result": res.GetInt("freq")})
		}
		for _, res := range tmpResultPower {
			resultPowerLost = append(resultPowerLost, tk.M{"_id": res["_id"], "result": res.GetFloat64("powerlost")})
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
	var turbineList []TurbineOut
	if p.ProjectName != "Fleet" {
		turbineList, _ = helper.GetTurbineList([]interface{}{p.ProjectName})
	} else {
		turbineList, _ = helper.GetTurbineList(nil)
	}
	totalTurbine := float64(len(turbineList))

	log.Printf(">>> %v \n", totalTurbine)

	if p.DateStr == "" {
		fromDate = p.Date.AddDate(0, -12, 0)

		match.Set("dateinfo.dateid", tk.M{"$gte": fromDate.UTC(), "$lte": p.Date.UTC()})
		match.Set("available", 1)

		if p.ProjectName != "Fleet" {
			match.Set("projectname", p.ProjectName)
		}

		group := tk.M{
			"_id":     tk.M{"id1": "$dateinfo.monthid", "id2": "$dateinfo.monthdesc", "id3": "$projectname"},
			"minutes": tk.M{"$sum": "$minutes"},
			"maxdate": tk.M{"$max": "$dateinfo.dateid"},
			"mindate": tk.M{"$min": "$dateinfo.dateid"},
		}

		if availType == "machine" {
			group.Set("result", tk.M{"$sum": "$machinedowntime"})
			// group.Set("result", tk.M{"$avg": "$machineavail"})
		} else if availType == "grid" {
			group.Set("result", tk.M{"$sum": "$griddowntime"})
			// group.Set("result", tk.M{"$avg": "$gridavail"})
		}

		pipe := []tk.M{
			{"$match": match},
			{"$group": group},
			{"$sort": tk.M{"_id.id1": -1}},
			{"$limit": 12},
		}

		// tk.Printf("pipe: %#v \n", pipe)

		/*csr, e := DB().Connection.NewQuery().
		From(new(ScadaSummaryDaily).TableName()).
		Command("pipe", pipe).
		Cursor(nil)*/

		csr, e := DB().Connection.NewQuery().
			From(new(ScadaData).TableName()).
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
		// tk.Println(availType)
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
					// log.Printf("LOOP: %#v | %#v %v \n", id, dateInfoFrom, project)
					if dateInfoFrom.MonthId == id.GetInt("id1") && project == id.GetString("id3") {
						exist = res
						break existData
					}
				}
				// tk.Println()
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
		for _, scada := range result {
			res := scada.GetFloat64("result")
			if scada.Get("mindate") != nil {
				minDate := scada.Get("mindate").(time.Time)
				maxDate := scada.Get("maxdate").(time.Time)
				minutes := scada.GetFloat64("minutes") / 60

				hourValue := helper.GetHourValue(fromDate.UTC(), p.Date.UTC(), minDate.UTC(), maxDate.UTC())
				avail := (minutes - (res / 3600.0)) / (totalTurbine * hourValue)
				scada.Set("result", avail)

				// log.Printf("SCADA: %v | %v | %v | %v = %v | %v - %v - %v - %v \n", minutes, res/3600.0, totalTurbine, hourValue, tk.ToFloat64(avail, 2, tk.RoundingAuto), fromDate.UTC().String(), p.Date.UTC().String(), minDate.UTC().String(), maxDate.UTC().String())
			} else {
				// log.Printf("SCADA_X: %v | %#v \n", res, scada.Get("_id"))
			}
			scada.Unset("maxdate")
			scada.Unset("mindate")
			scada.Unset("minutes")
		}
	}
	// tk.Println()
	return
}

func getMGAvailability(p *PayloadDashboard) (machineResult []tk.M, gridResult []tk.M) {
	result := []tk.M{}
	var fromDate time.Time
	match := tk.M{}
	// total turbine should follow projects, for now it's hardcoded
	var turbineList []TurbineOut
	if p.ProjectName != "Fleet" {
		turbineList, _ = helper.GetTurbineList([]interface{}{p.ProjectName})
	} else {
		turbineList, _ = helper.GetTurbineList(nil)
	}
	totalTurbine := float64(len(turbineList))

	p.Date, _ = time.Parse("2006-01-02 15:04:05", p.Date.UTC().Format("2006-01")+"-01"+" 00:00:00")

	if p.DateStr == "" {
		fromDate = p.Date.AddDate(0, -12, 0)

		match.Set("dateinfo.dateid", tk.M{"$gte": fromDate.UTC(), "$lte": p.Date.UTC()})
		match.Set("available", 1)

		if p.ProjectName != "Fleet" {
			match.Set("projectname", p.ProjectName)
		}

		group := tk.M{
			"_id":           tk.M{"id1": "$dateinfo.monthid", "id2": "$dateinfo.monthdesc", "id3": "$projectname"},
			"minutes":       tk.M{"$sum": "$minutes"},
			"mindate":       tk.M{"$min": "$dateinfo.dateid"},
			"maxdate":       tk.M{"$max": "$dateinfo.dateid"},
			"machineResult": tk.M{"$sum": "$machinedowntime"},
			"gridResult":    tk.M{"$sum": "$griddowntime"},
		}

		pipe := []tk.M{
			{"$match": match},
			{"$group": group},
			{"$sort": tk.M{"_id.id1": -1}},
			{"$limit": 12},
		}

		csr, e := DB().Connection.NewQuery().
			From(new(ScadaData).TableName()).
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

		projects := []string{"Tejuva"}

		// --------------

		// dayInYear := tk.M{}
		tmpFromDate := fromDate.AddDate(0, 1, 0)
		dateInfoTo := GetDateInfo(p.Date)
		for _, project := range projects {

		done:

			for {
				dateInfoFrom := GetDateInfo(tmpFromDate)
				// log.Printf("%v \n", dateInfoFrom.MonthDesc)
				// if dayInYear.Get(tk.ToString(dateInfoFrom.Year)) == nil {
				// 	dayInYear.Set(tk.ToString(dateInfoFrom.Year), GetDayInYear(dateInfoFrom.Year))
				// }
				// days := dayInYear.Get(tk.ToString(dateInfoFrom.Year)).(tk.M).GetInt(tk.ToString(int(tmpFromDate.Month())))

				var exist tk.M

			existData:

				for _, res := range tmpResult {
					id := res.Get("_id").(tk.M)
					// log.Printf("LOOP: %#v | %#v %v \n", id, dateInfoFrom, project)
					if dateInfoFrom.MonthId == id.GetInt("id1") && project == id.GetString("id3") {
						exist = res
						break existData
					}
				}

				if exist != nil {
					result = append(result, exist)
				} else {
					result = append(result, tk.M{
						"_id": tk.M{
							"id1": dateInfoFrom.MonthId,
							"id2": dateInfoFrom.MonthDesc,
							"id3": project,
						},
						"machineResult": 0.00,
						"gridResult":    0.00,
					})
				}

				if dateInfoFrom.MonthId == dateInfoTo.MonthId {
					tmpFromDate = fromDate.AddDate(0, 1, 0)
					break done
				}

				tmpFromDate = tmpFromDate.AddDate(0, 1, 0)
			}
		}

		for _, scada := range result {
			if scada.Get("mindate") != nil {
				m := scada.GetFloat64("machineResult") / 3600.0
				g := scada.GetFloat64("gridResult") / 3600.0
				minDate := scada.Get("mindate").(time.Time)
				maxDate := scada.Get("maxdate").(time.Time)
				minutes := scada.GetFloat64("minutes") / 60

				fromDateSub, _ := time.Parse("060102_150405", minDate.Format("0601")+"01_000000")
				tmpDt, _ := time.Parse("060102_150405", minDate.AddDate(0, 1, 0).Format("0601")+"01_000000")
				toDateSub := tmpDt.AddDate(0, 0, -1)

				hourValue := helper.GetHourValue(fromDateSub.UTC(), toDateSub.UTC(), minDate.UTC(), maxDate.UTC())
				mAvail, gAvail, _, _, _ := helper.GetAvailAndPLF(totalTurbine, float64(0), float64(0), m, g, float64(0), hourValue, minutes)

				// log.Printf("%v | %v \n", mAvail, gAvail)

				scada.Set("machineResult", tk.ToFloat64((mAvail), 2, tk.RoundingAuto)/100)
				scada.Set("gridResult", tk.ToFloat64((gAvail), 2, tk.RoundingAuto)/100)

				// log.Printf(">>> %#v \n", scada)

				// log.Printf("SCADA: %v | %v | %v | %v = %v | %v - %v - %v - %v \n", minutes, res/3600.0, totalTurbine, hourValue, tk.ToFloat64(avail, 2, tk.RoundingAuto), fromDate.UTC().String(), p.Date.UTC().String(), minDate.UTC().String(), maxDate.UTC().String())
			}

			scada.Unset("maxdate")
			scada.Unset("mindate")
			scada.Unset("minutes")
		}
	}

	for _, res := range result {
		id := res.Get("_id")

		machineResult = append(machineResult, tk.M{"_id": id, "result": res.GetFloat64("machineResult")})
		gridResult = append(gridResult, tk.M{"_id": id, "result": res.GetFloat64("gridResult")})
	}

	// tk.Println()
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

	result := getDownTurbine(p.ProjectName, p.Date, 1)

	return helper.CreateResult(true, result, "success")
}

func getDownTurbine(project string, currentDate time.Time, dayDuration int) (result []tk.M) {
	var fromDate time.Time
	var pipes []tk.M

	fromDate = currentDate.UTC().AddDate(0, 0, dayDuration*-1)

	match := tk.M{"startdate": tk.M{"$gte": fromDate.UTC(), "$lte": currentDate.UTC()}}

	if project != "Fleet" {
		match.Set("farm", project)
	}

	pipes = append(pipes, tk.M{"$match": match})
	pipes = append(pipes, tk.M{"$group": tk.M{"_id": "$turbine", "result": tk.M{"$sum": "$duration"}}})
	pipes = append(pipes, tk.M{"$sort": tk.M{"_id": 1}})

	csr, e := DB().Connection.NewQuery().
		From(new(Alarm).TableName()).
		Command("pipe", pipes).
		Cursor(nil)

	if e != nil {
		return
	}

	tmpResult := []tk.M{}

	e = csr.Fetch(&tmpResult, 0, false)
	defer csr.Close()
	csr.Close()

	if e != nil {
		return
	}

	for _, val := range tmpResult {
		if val.GetFloat64("result") >= float64(24*dayDuration) {
			val.Set("isdown", true)
			result = append(result, val)
		}
	}

	return
}

func getMapCol(project string) tk.Ms {
	filter := []*dbox.Filter{}
	colname := new(ProjectMaster).TableName()

	if project != "Fleet" {
		colname = new(TurbineMaster).TableName()
		filter = append(filter, dbox.Eq("project", project))
	} else {
		filter = append(filter, dbox.Eq("active", true))
	}

	csr, e := DB().Connection.NewQuery().
		From(colname).
		Where(filter...).
		Cursor(nil)
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
	projectName := payload["projectname"]

	data := getMapCol(projectName)

	results := tk.Ms{}
	offset := []int{0, 2}
	coords := []float64{}
	for _, val := range data {
		result := tk.M{}
		coords = []float64{}
		coords = []float64{val.GetFloat64("latitude"), val.GetFloat64("longitude")}
		if projectName == "Fleet" {
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
	pipes := []tk.M{}
	match := []tk.M{}

	p := new(PayloadDashboard)
	e := k.GetPayload(&p)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	machinedown, e := getMachineDownType()
	var typeX string

	for i, v := range machinedown {
		if v == p.Type {
			typeX = i
			break
		}
	}

	result := tk.M{}

	dateStr := []string{}
	date := time.Time{}
	date2 := time.Time{}
	if p.DateStr != "fleet date" {
		dateStr = strings.Split(p.DateStr, " ")
		date, e = time.Parse("Jan 2006 02 15:04:05", dateStr[0][0:3]+" "+dateStr[1]+" 01 00:00:00")
	} else {
		date2 = helper.GetLastDateData(k)
		date = date2.AddDate(0, -12, 0)

		/*dateStr = strings.Split("Jul 2015", " ")
		dateStr2 := strings.Split("Jun 2016", " ")
		date, e = time.Parse("Jan 2006 02 15:04:05", dateStr[0][0:3]+" "+dateStr[1]+" 01 00:00:00")
		date2, e = time.Parse("Jan 2006 02 15:04:05", dateStr2[0][0:3]+" "+dateStr2[1]+" 01 00:00:00")*/
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
		match = append(match, tk.M{"detail.detaildateinfo.monthid": dateInfo.MonthId})
		filter = append(filter, dbox.Eq("detail.detaildateinfo.monthid", dateInfo.MonthId))
	} else {
		filter = append(filter, dbox.Gte("startdateinfo.monthid", dateInfo.MonthId))
		filter = append(filter, dbox.Lte("startdateinfo.monthid", dateInfo2.MonthId))

		match = append(match, tk.M{"detail.detaildateinfo.monthid": tk.M{"$gt": dateInfo.MonthId, "$lte": dateInfo2.MonthId}})
		// tk.Println(dateInfo.MonthId)
		// tk.Println(dateInfo2.MonthId)
	}
	if p.ProjectName != "Fleet" {
		if p.Type != "" && p.Type != "All Types" {
			filter = append(filter, dbox.Eq("detail."+typeX, true))
			match = append(match, tk.M{"detail." + typeX: true})
		}

		filter = append(filter, dbox.Eq("projectname", p.ProjectName))
		match = append(match, tk.M{"projectname": p.ProjectName})
	} else {
		if p.Type != "" && p.Type != "All Types" {
			filter = append(filter, dbox.Eq("detail."+typeX, true))
			match = append(match, tk.M{"detail." + typeX: true})
		}
	}

	pipes = append(pipes, tk.M{"$unwind": "$detail"})
	pipes = append(pipes, tk.M{"$match": match})

	query := DB().Connection.NewQuery().
		From(new(Alarm).TableName()).Where(filter...).
		Skip(p.Skip).Take(p.Take)

	/*query := DB().Connection.NewQuery().
	From(new(Alarm).TableName()).
	Command("pipes", pipes).
	Skip(p.Skip).Take(p.Take)*/

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

	/*csr, e = DB().Connection.NewQuery().
	From(new(Alarm).TableName()).
	Command("pipes", pipes).
	Cursor(nil)*/

	// add by ams, 2016-10-07
	defer csr.Close()

	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	total := csr.Count()

	/*for _, v := range pipes {
		log.Printf("pipes: %#v \n", v)
	}

	if len(resTable) > 0 {
		log.Printf("resTable: %#v \n", resTable[0])
	}*/

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
				log.Println(err.Error())
			} else {
				from = t
			}
		} else if filt.Field == "dateinfo.dateid" && filt.Op == "lte" {
			b, err := time.Parse("2006-01-02T15:04:05.000Z", filt.Value.(string))
			t, _ := time.Parse("2006-01-02 15:04:05.999999999", b.UTC().Format("2006-01-02")+" 23:59:59.999999999")
			if err != nil {
				log.Println(err.Error())
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

	// log.Printf("tmp: \n%#v \n", tmp)

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
	group.Set("maxDate", tk.M{"$max": "$dateinfo.dateid"})

	pipe := []tk.M{
		{"$match": matches},
		{"$group": group},
		{"$sort": tk.M{"_id": 1}},
		// {"$skip": p.Skip},
		// {"$limit": p.Take},
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
			result[idx].Set("plf", (result[idx].GetFloat64("production")/1000000)/(maxCapacity/1000))
		} else {
			result[idx].Set("name", dt.GetString("_id"))
			result[idx].Set("noofwtg", totalTurbine.GetInt(result[idx].GetString("name")))

			maxCapacity := turbineMW * totalTurbine.GetFloat64(result[idx].GetString("name")) * totalHours
			result[idx].Set("maxcapacity", maxCapacity)
			result[idx].Set("plf", (result[idx].GetFloat64("production")/1000000)/(maxCapacity/1000))
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
