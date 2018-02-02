package controller

import (
	. "eaciit/wfdemo-git/library/core"
	. "eaciit/wfdemo-git/library/helper"
	. "eaciit/wfdemo-git/library/models"
	"eaciit/wfdemo-git/web/helper"
	"log"
	"strings"

	"sync"

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
	projects, e := helper.GetProjectList()
	result := []string{}

	if e != nil {
		return result, e
	}

	for _, val := range projects {
		result = append(result, val.Value)
	}
	sort.Strings(result)

	return result, e
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

		// turbineDownOneDays := getTotalDownTurbine(val.ProjectName, val.LastUpdate, 0)
		// turbineDownTwoDays := getTotalDownTurbine(val.ProjectName, val.LastUpdate, 2)
		turbineDownOneDays := len(getDownTurbineStatus(val.ProjectName, val.LastUpdate, 0))
		turbineDownTwoDays := 0

		val.CurrentDown = turbineDownOneDays
		val.TwoDaysDown = turbineDownTwoDays

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
	TrueAvail          interface{}
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

func GetDataAvailability(datalastmonth []ScadaSummaryByMonth, projectname string) []interface{} {
	var result []interface{}

	startmonth := 0
	endmonth := datalastmonth[0].DateInfo.MonthId
	month := endmonth - (int(endmonth/100) * 100)

	if month == 12 {
		startmonth = (int(endmonth/100) * 100) + 1
	} else {
		startmonth = (endmonth + 1) - 100
	}

	dataVariance := new(ScadaSummaryVariance)

	csr, e := DB().Connection.NewQuery().
		From("rpt_scadasummarybymonth").
		Where(dbox.Eq("projectname", projectname)).
		Cursor(nil)

	if e != nil {
		return nil
	}
	defer csr.Close()

	data := make([]ScadaSummaryByMonth, 0)
	e = csr.Fetch(&data, 0, false)
	dataPerMonth := map[string]ScadaSummaryByMonth{}

	for _, val := range data {
		dataPerMonth[val.ProjectName+"_"+tk.ToString(val.DateInfo.MonthId)] = val
	}

	hasValue := false
	keys := ""

	if len(data) > 0 {
		for i := startmonth; i <= endmonth; i++ {
			//check if month more than 12
			if i-(int(i/100)*100) > 12 {
				i = (i - 12) + 100
			}

			yearloop := int(i / 100)
			monthloop := i - (int(i/100) * 100)
			keys = projectname + "_" + tk.ToString(i)
			_, hasValue = dataPerMonth[keys]

			if hasValue {
				dataVariance.ID = dataPerMonth[keys].ID
				dataVariance.DateInfo = dataPerMonth[keys].DateInfo
				dataVariance.ProjectName = dataPerMonth[keys].ProjectName
				dataVariance.Production = dataPerMonth[keys].Production / 1000
				dataVariance.ProductionLastYear = dataPerMonth[keys].ProductionLastYear
				dataVariance.Revenue = dataPerMonth[keys].Revenue
				dataVariance.RevenueInLacs = dataPerMonth[keys].RevenueInLacs
				dataVariance.TrueAvail = dataPerMonth[keys].TrueAvail
				if dataPerMonth[keys].TrueAvail == 0 && dataPerMonth[keys].ScadaAvail == 0 {
					dataVariance.TrueAvail = nil
				}
				dataVariance.ScadaAvail = dataPerMonth[keys].ScadaAvail
				dataVariance.MachineAvail = dataPerMonth[keys].MachineAvail
				dataVariance.GridAvail = dataPerMonth[keys].GridAvail
				dataVariance.PLF = dataPerMonth[keys].PLF
				dataVariance.Budget = dataPerMonth[keys].Budget / 1000000
				dataVariance.AvgWindSpeed = dataPerMonth[keys].AvgWindSpeed
				dataVariance.ExpWindSpeed = dataPerMonth[keys].ExpWindSpeed
				dataVariance.DowntimeHours = dataPerMonth[keys].DowntimeHours
				dataVariance.LostEnergy = dataPerMonth[keys].LostEnergy
				dataVariance.RevenueLoss = dataPerMonth[keys].RevenueLoss
				if dataPerMonth[keys].ProductionLastYear == 0 {
					dataVariance.Variance = 100
				} else {
					dataVariance.Variance = math.Abs((dataPerMonth[keys].Production - dataPerMonth[keys].ProductionLastYear) / dataPerMonth[keys].ProductionLastYear * 100)
				}

				result = append(result, *dataVariance)
			} else {
				// Temporary data to fill result if month doesn't exist
				datatemp := new(ScadaSummaryByMonth)

				datatemp.ID = ""
				datatemp.ProjectName = projectname
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
				if datatemp.TrueAvail == 0 && datatemp.ScadaAvail == 0 {
					dataVariance.TrueAvail = nil
				}
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

	return result
}

func (m *DashboardController) GetScadaSummaryByMonth(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson
	type PayloadSummaryByMonth struct {
		ProjectName string
		Date        time.Time
		ProjectList []string
	}
	p := new(PayloadSummaryByMonth)
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

	var finalResult tk.M
	finalResult = tk.M{}

	if len(datalastmonth) > 0 {
		result := GetDataAvailability(datalastmonth, p.ProjectName)
		finalResult.Set("Data", result)
		if p.ProjectName == "Fleet" {
			availData := tk.M{}
			for _, projectName := range p.ProjectList {
				avail := GetDataAvailability(datalastmonth, projectName)
				availData.Set(projectName, avail)
			}
			finalResult.Set("Availability", availData)
		}
	}

	return helper.CreateResult(true, finalResult, "success")
}

func (m *DashboardController) GetDetailProd(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	p := tk.M{}
	e := k.GetPayload(&p)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}
	var turbineList []TurbineOut
	ids := tk.M{"project": "$projectname", "turbine": "$turbine"}
	matches := tk.M{"dateinfo.monthdesc": p.GetString("date")}
	if p.GetString("project") != "Fleet" {
		matches.Set("projectname", p.GetString("project"))
		turbineList, _ = helper.GetTurbineList([]interface{}{p.GetString("project")})
	} else {
		turbineList, _ = helper.GetTurbineList(nil)
	}
	turbineCapacity := map[string]float64{}
	for _, v := range turbineList {
		turbineCapacity[tk.Sprintf("%s_%s", v.Project, v.Value)] = v.Capacity
	}

	pipe := []tk.M{{"$match": matches},
		{"$group": tk.M{
			"_id":            ids,
			"production":     tk.M{"$sum": "$production"},
			"lostenergy":     tk.M{"$sum": "$lostenergy"},
			"dateid":         tk.M{"$max": "$dateinfo.dateid"},
			"mdownhours":     tk.M{"$sum": "$machinedownhours"},
			"gdownhours":     tk.M{"$sum": "$griddownhours"},
			"odownhours":     tk.M{"$sum": "$otherdowntimehours"},
			"oktime":         tk.M{"$sum": "$oktime"},
			"totaltimestamp": tk.M{"$sum": 1},
			"power":          tk.M{"$sum": "$powerkw"},
			"maxdate":        tk.M{"$max": "$dateinfo.dateid"},
			"mindate":        tk.M{"$min": "$dateinfo.dateid"},
		}},
		{"$sort": tk.M{"_id.project": 1}}}

	csrScada, e := DB().Connection.NewQuery().
		From(new(ScadaSummaryDaily).TableName()).
		Command("pipe", pipe).
		Cursor(nil)

	if e != nil {
		helper.CreateResult(false, nil, e.Error())
	}
	defer csrScada.Close()

	resultScada := []tk.M{}
	e = csrScada.Fetch(&resultScada, 0, false)
	if e != nil {
		helper.CreateResult(false, nil, e.Error())
	}

	totalPower := tk.M{}
	totalPowerLost := tk.M{}
	totalTurbines := tk.M{}
	detailData := tk.M{}
	detail := []tk.M{}

	listturbine := tk.M{}
	tproject := ""
	maxdate := time.Time{}
	for _, val := range resultScada {
		data := val["_id"].(tk.M)
		project := data.GetString("project")
		if tproject != project {
			tproject = project
			detail = []tk.M{}
			listturbine = PopulateTurbines(DB().Connection, tproject)
		}
		val.Unset("_id")
		hourValue := val.Get("maxdate").(time.Time).AddDate(0, 0, 1).UTC().Sub(val.Get("mindate").(time.Time).UTC()).Hours()

		okTime := val.GetFloat64("oktime") / 3600 /*jadikan hour*/
		power := val.GetFloat64("power") / 1000.0 /*megaWatt*/
		energy := power / 6                       /*MWh karena kapasitas turbine maksimal hanya MWh*/
		mDownTime := val.GetFloat64("mdownhours")
		gDownTime := val.GetFloat64("gdownhours")
		uDownTime := val.GetFloat64("odownhours")
		sumTimeStamp := val.GetFloat64("totaltimestamp")

		in := tk.M{}.Set("noofturbine", 1).Set("oktime", okTime).Set("energy", energy).
			Set("totalhour", hourValue).Set("totalcapacity", turbineCapacity[tk.Sprintf("%s_%s", project, data.GetString("turbine"))]).
			Set("counttimestamp", sumTimeStamp).Set("machinedowntime", mDownTime).Set("griddowntime", gDownTime).
			Set("otherdowntime", uDownTime)
		res := helper.CalcAvailabilityAndPLF(in)
		val.Set("plf", res.GetFloat64("plf"))
		val.Set("trueavail", res.GetFloat64("totalavailability"))

		val.Set("turbine", data.GetString("turbine"))
		if listturbine.Has(data.GetString("turbine")) {
			val.Set("turbine", listturbine.GetString(data.GetString("turbine")))
		}

		downtimehours := val.GetFloat64("mdownhours") + val.GetFloat64("gdownhours") + val.GetFloat64("odownhours")
		val.Set("downtimehours", downtimehours)
		val.Unset("maxdate")
		val.Unset("mindate")

		detail = append(detail, val)
		detailData.Set(project, detail)

		if totalTurbines.Has(project) {
			totalTurbines.Set(project, totalTurbines.GetInt(project)+1)
		} else {
			totalTurbines.Set(project, 1)
		}
		if totalPower.Has(project) {
			totalPower.Set(project, totalPower.GetFloat64(project)+val.GetFloat64("production"))
		} else {
			totalPower.Set(project, val.GetFloat64("production"))
		}

		if totalPowerLost.Has(project) {
			totalPowerLost.Set(project, totalPowerLost.GetFloat64(project)+val.GetFloat64("lostenergy"))
		} else {
			totalPowerLost.Set(project, val.GetFloat64("lostenergy"))
		}

		mdateid := val.Get("dateid", time.Time{}).(time.Time)
		if maxdate.UTC().Before(mdateid.UTC()) {
			maxdate = mdateid
		}
	}

	csrMonthly, e := DB().Connection.NewQuery().
		From(new(ScadaSummaryByMonth).TableName()).
		Where(dbox.And(
			dbox.Eq("projectname", p.GetString("project")),
			dbox.Eq("dateinfo.monthdesc", p.GetString("date")),
		)).
		Cursor(nil)

	if e != nil {
		helper.CreateResult(false, nil, e.Error())
	}
	defer csrMonthly.Close()

	resultMonthly := []ScadaSummaryByMonth{}
	e = csrMonthly.Fetch(&resultMonthly, 0, false)
	if e != nil {
		helper.CreateResult(false, nil, e.Error())
	}

	xbudget := float64(1)
	dataItemTemp := []tk.M{}

	if len(resultMonthly) > 0 {
		bulan := resultMonthly[0].DateInfo.MonthId - (resultMonthly[0].DateInfo.Year * 100)
		filterBudget := []*dbox.Filter{}
		query := DB().Connection.NewQuery().From(new(ExpPValueModel).TableName())
		filterBudget = append(filterBudget, dbox.Eq("monthno", bulan))
		if p.GetString("project") != "Fleet" {
			filterBudget = append(filterBudget, dbox.Eq("projectname", p.GetString("project")))
		}
		csrBudget, e := query.Where(dbox.And(filterBudget...)).Cursor(nil)

		if tnow := getTimeNow(); int(tnow.Month()) == bulan {
			maxdate = maxdate.AddDate(0, 0, 1)
			tdays := maxdate.UTC().Sub(resultMonthly[0].DateInfo.DateId.UTC()).Hours() / 24
			mdays := tk.ToFloat64(time.Date(maxdate.Year(), maxdate.Month(), 0, 0, 0, 0, 0, time.UTC).Day(), 0, tk.RoundingAuto)
			xbudget = tdays / mdays
		}

		if e != nil {
			helper.CreateResult(false, nil, e.Error())
		}
		defer csrBudget.Close()

		resultBudget := []ExpPValueModel{}
		e = csrBudget.Fetch(&resultBudget, 0, false)
		if e != nil {
			helper.CreateResult(false, nil, e.Error())
		}
		budgetPerProject := map[string]ExpPValueModel{}
		for _, val := range resultBudget {
			budgetPerProject[val.ProjectName] = val
		}

		projectList, _ := helper.GetProjectList()
		dataItem := []tk.M{}
		for project, val := range totalPower {

			labelproject := ""
			for _, _info := range projectList {
				if strings.ToLower(_info.Value) == strings.ToLower(project) {
					labelproject = _info.Name
					break
				}
			}
			budget := budgetPerProject[project]

			data := tk.M{
				"project":       project,
				"production":    val.(float64),
				"lostenergy":    totalPowerLost.GetFloat64(project),
				"wtg":           totalTurbines.GetInt(project),
				"labelproject":  labelproject,
				"detail":        detailData[project],
				"avgwindspeed":  resultMonthly[0].AvgWindSpeed,
				"downtimehours": resultMonthly[0].DowntimeHours,
				"plf":           resultMonthly[0].PLF,
				"trueavail":     resultMonthly[0].TrueAvail,
				"budget_p50":    budget.P50NetGenMWH * xbudget,
				"budget_p75":    budget.P75NetGenMWH * xbudget,
				"budget_p90":    budget.P90NetGenMWH * xbudget,
			}
			dataItem = append(dataItem, data)
		}
		dataItemTemp = dataItem
	}

	dataOutput := []tk.M{}
	for _, val := range dataItemTemp {
		newdata := helper.EnergyMeasurement(val, "production", "lostenergy")
		val = newdata[0]
		newdetail := helper.EnergyMeasurement(val["detail"].([]tk.M), "production", "lostenergy")
		val.Set("detail", newdetail)
		dataOutput = append(dataOutput, val)
	}

	return helper.CreateResult(true, dataOutput, "success")
}

func (m *DashboardController) GetDetailProdLevel1(k *knot.WebContext) interface{} {
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
	pipe := []tk.M{{"$match": matches},
		{"$group": tk.M{
			"_id":        ids,
			"production": tk.M{"$sum": "$production"},
			"lostenergy": tk.M{"$sum": "$lostenergy"},
			"mdownhours": tk.M{"$sum": "$machinedownhours"},
			"gdownhours": tk.M{"$sum": "$griddownhours"},
			"odownhours": tk.M{"$sum": "$otherdowntimehours"},
		}},
		{"$sort": tk.M{"projectname": 1}}}

	csrScada, e := DB().Connection.NewQuery().
		From(new(ScadaSummaryDaily).TableName()).
		Command("pipe", pipe).
		Cursor(nil)

	if e != nil {
		helper.CreateResult(false, nil, e.Error())
	}
	defer csrScada.Close()

	resultScada := []tk.M{}
	e = csrScada.Fetch(&resultScada, 0, false)
	if e != nil {
		helper.CreateResult(false, nil, e.Error())
	}

	totalEnergy := map[string]float64{}
	totalEnergyLost := map[string]float64{}
	totalTurbines := map[string]int{}
	detailData := tk.M{}
	detail := []tk.M{}
	downtimehours := 0.0

	listturbine := tk.M{}
	tproject := ""
	for _, val := range resultScada {
		data := val["_id"].(tk.M)
		project := data.GetString("project")
		if tproject != project {
			tproject = project
			listturbine = PopulateTurbines(DB().Connection, tproject)
		}
		val.Unset("_id")
		val.Set("turbine", data.GetString("turbine"))
		if listturbine.Has(data.GetString("turbine")) {
			val.Set("turbine", listturbine.GetString(data.GetString("turbine")))
		}
		downtimehours = val.GetFloat64("mdownhours") + val.GetFloat64("gdownhours") + val.GetFloat64("odownhours")
		val.Set("downtimehours", downtimehours)

		detail = append(detail, val)
		detailData.Set(project, detail)

		totalEnergy[project] += val.GetFloat64("production")
		totalEnergyLost[project] += val.GetFloat64("lostenergy")
		totalTurbines[project]++
	}

	dataItem := []tk.M{}
	for project, val := range totalEnergy {
		data := tk.M{
			"project":    project,
			"production": val,
			"lostenergy": totalEnergyLost[project],
			"wtg":        totalTurbines[project],
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

	listname := tk.M{}
	if p.GetString("project") != "Fleet" {
		listname = PopulateTurbines(DB().Connection, p.GetString("project"))
	}

	for _, val := range result {
		vtkm := val["dataitems"].(tk.M)
		if listname.Has(vtkm.GetString("name")) {
			vtkm.Set("name", listname.GetString(vtkm.GetString("name")))
		}
		dataItem = append(dataItem, vtkm)
	}

	data := struct {
		Data  []tk.M
		Total int
	}{
		Data:  dataItem,
		Total: tk.SliceLen(dataItem),
	}

	// log.Printf("> %#v \n", data)

	return helper.CreateResult(true, data, "success")
}

func getLossDuration(topType string, p *PayloadAnalytic, k *knot.WebContext) ([]tk.M, error) {
	result := []tk.M{}
	breakdown := "$projectname"
	if p.BreakDown == "$turbine" {
		breakdown = p.BreakDown
	}
	var e error
	var pipes []tk.M
	match := tk.M{}

	if p != nil {
		tStart, tEnd, e := helper.GetStartEndDate(k, p.Period, p.DateStart, p.DateEnd)
		if e != nil {
			return result, e
		}
		match.Set("dateinfo.dateid", tk.M{"$gte": tStart, "$lte": tEnd})

		if p.Project != "" {
			match.Set("projectname", p.Project)
		}

		if len(p.Turbine) != 0 {
			match.Set("turbine", tk.M{"$in": p.Turbine})
		}

		pipes = append(pipes, tk.M{"$match": match})
		if topType == "duration" {
			pipes = append(pipes, tk.M{"$group": tk.M{
				"_id":         breakdown,
				"machinedown": tk.M{"$sum": "$machinedownhours"},
				"griddown":    tk.M{"$sum": "$griddownhours"},
				"unknown":     tk.M{"$sum": "$otherdowntimehours"},
			}})
		} else if topType == "loss" {
			pipes = append(pipes, tk.M{"$group": tk.M{
				"_id":         breakdown,
				"machinedown": tk.M{"$sum": "$machinedownloss"},
				"griddown":    tk.M{"$sum": "$griddownloss"},
				"unknown":     tk.M{"$sum": "$otherdownloss"},
			}})
		} else if topType == "both" {
			pipes = append(pipes, tk.M{"$group": tk.M{
				"_id": breakdown,
				"machinedownduration": tk.M{"$sum": "$machinedownhours"},
				"griddownduration":    tk.M{"$sum": "$griddownhours"},
				"unknownduration":     tk.M{"$sum": "$otherdowntimehours"},
				"machinedownloss":     tk.M{"$sum": "$machinedownloss"},
				"griddownloss":        tk.M{"$sum": "$griddownloss"},
				"unknownloss":         tk.M{"$sum": "$otherdownloss"},
			}})
		}

		csr, e := DB().Connection.NewQuery().
			From(new(ScadaSummaryDaily).TableName()).
			Command("pipe", pipes).
			Cursor(nil)

		if e != nil {
			return result, e
		}

		e = csr.Fetch(&result, 0, false)
		if e != nil {
			return result, e
		}
		defer csr.Close()
	}

	return result, e
}

func getLossFrequency(p *PayloadAnalytic, k *knot.WebContext) (tk.M, error) {
	result := tk.M{}
	var e error
	var pipes []tk.M
	matches := []tk.M{}

	if p != nil {
		tStart, tEnd, e := helper.GetStartEndDate(k, p.Period, p.DateStart, p.DateEnd)
		if e != nil {
			return result, e
		}
		matches = append(matches, tk.M{"startdate": tk.M{"$gte": tStart}})
		matches = append(matches, tk.M{"startdate": tk.M{"$lte": tEnd}})
		matches = append(matches, tk.M{"reduceavailability": true})
		matches = append(matches, tk.M{"duration": tk.M{"$gt": 0}})

		if p.Project != "" {
			matches = append(matches, tk.M{"projectname": p.Project})
		}

		if len(p.Turbine) != 0 {
			matches = append(matches, tk.M{"turbine": tk.M{"$in": p.Turbine}})
		}

		downCause, _ := getMachineDownType()

		for field := range downCause {
			pipes = []tk.M{}
			loopMatch := matches

			loopMatch = append(loopMatch, tk.M{field: true})

			pipes = append(pipes, tk.M{"$match": tk.M{"$and": loopMatch}})
			pipes = append(pipes,
				tk.M{
					"$group": tk.M{"_id": "", "result": tk.M{"$sum": 1}},
				},
			)

			csr, e := DB().Connection.NewQuery().
				From(new(Alarm).TableName()).
				Command("pipe", pipes).
				Cursor(nil)

			if e != nil {
				return result, e
			}

			resLoop := []tk.M{}
			e = csr.Fetch(&resLoop, 0, false)

			csr.Close()

			for _, res := range resLoop {
				res.Unset("_id")
				result.Set(field, res.GetFloat64("result"))
			}
		}
	}

	return result, e
}

func (m *DashboardController) GetDownTimeLoss(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	p := new(PayloadAnalytic)
	e := k.GetPayload(&p)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	machinedown, _ := getMachineDownType()

	result := []tk.M{}
	lossDurationData := tk.M{}
	_lossDurData, e := getLossDuration("both", p, k)
	if len(_lossDurData) > 0 {
		lossDurationData = _lossDurData[0]
	}
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}
	freqData, e := getLossFrequency(p, k)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	for field, mdName := range machinedown {
		res := tk.M{}
		res.Set("field", mdName)
		res.Set("title", mdName)
		res.Set("projectname", p.Project)
		res.Set("powerlost", lossDurationData.GetFloat64(field+"loss")/1000)
		res.Set("duration", lossDurationData.GetFloat64(field+"duration"))
		res.Set("frequency", freqData.GetFloat64(field))
		result = append(result, res)
	}
	return helper.CreateResult(true, result, "success")
}

func (m *DashboardController) GetLostEnergy(k *knot.WebContext) interface{} { /* hanya dipakai di dashboard availability */
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

func (m *DashboardController) GetDetailLossLevel1(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	results := tk.M{}

	p := new(PayloadDashboard)
	e := k.GetPayload(&p)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	var pipes []tk.M
	var fromDate time.Time
	match := tk.M{}

	fromDate = p.Date.AddDate(0, -11, 0)
	match.Set("dateinfo.dateid", tk.M{"$gte": fromDate.UTC(), "$lte": p.Date.UTC()})

	if p.ProjectName != "Fleet" {
		match.Set("projectname", p.ProjectName)
	}

	pipes = append(pipes, tk.M{"$match": match})
	pipes = append(pipes, tk.M{"$group": tk.M{
		"_id": tk.M{
			"projectname": "$projectname",
			"bulan":       "$dateinfo.monthid",
		},
		"result":      tk.M{"$sum": "$lostenergy"},
		"machinedown": tk.M{"$sum": "$machinedownloss"},
		"griddown":    tk.M{"$sum": "$griddownloss"},
		"unknown":     tk.M{"$sum": "$otherdownloss"},
	}})
	pipes = append(pipes, tk.M{"$sort": tk.M{"_id.bulan": 1}})

	// get the top 10 of turbine dan mengambil total

	csr, e := DB().Connection.NewQuery().
		From(new(ScadaSummaryDaily).TableName()).
		Command("pipe", pipes).
		Cursor(nil)

	if e != nil {
		return helper.CreateResult(false, results, e.Error())
	}

	allLossData := []tk.M{}
	e = csr.Fetch(&allLossData, 0, false)
	csr.Close()

	if e != nil {
		return helper.CreateResult(false, results, e.Error())
	}

	monthList := []int{}
	monthDescList := map[int]string{}
	monthCount := (fromDate.Year() * 100) + int(fromDate.Month()) /*201611*/
	maxMonth := (fromDate.Year() * 100) + 12                      /*201612*/
	monthInt := int(fromDate.Month())
	yearInt := fromDate.Year()
	for i := 1; i <= 12; i++ {
		if monthCount > maxMonth {
			monthCount = monthCount - maxMonth + (p.Date.Year() * 100) /*(201613 - 201612) + 201700*/
			maxMonth = (p.Date.Year() * 100) + 12
			yearInt = p.Date.Year()
			monthInt = 1
		}
		monthList = append(monthList, monthCount)
		monthDescList[monthCount] = tk.Sprintf("%s %d", time.Month(monthInt).String(), yearInt)
		monthCount++
		monthInt++
	}
	projectList := []string{}
	projects, _ := helper.GetProjectList()
	for _, val := range projects {
		projectList = append(projectList, val.ProjectId)
	}

	downCause := map[string]string{}
	downCause["griddown"] = "Grid Down"
	downCause["machinedown"] = "Machine Down"
	downCause["unknown"] = "Unknown"

	totalLossPerYear := tk.M{}
	/* expected output :
	{
		Amba: [
		{
			"name": "MachineDown",
			"value": xxxxxx,
		},
		{
			"name": "GridDown",
			"value": xxxxxx,
		}]
	}
	*/
	totalLoss := map[string]float64{}
	totalLossPerType := []tk.M{}

	result := []tk.M{}
	resultPerProject := tk.M{}

	for _, project := range projectList {
		if p.ProjectName != "Fleet" && project != p.ProjectName {
			continue
		}
		result = []tk.M{}
		totalLossPerType = []tk.M{}
		totalLoss = map[string]float64{}
		for _, month := range monthList {
			resVal := tk.M{}
			resVal.Set("_id", monthDescList[month]) /* _id: "August 2017" */
			for _, down := range downCause {
				valTitle := strings.Replace(down, " ", "", -69)
				resVal.Set(valTitle, 0.0) /* MachineDown : 0.0 ==> default value */
			}
			lossPerMonth := 0.0
			for _, val := range allLossData {
				valProject := val.Get("_id").(tk.M).GetString("projectname")
				valMonth := val.Get("_id").(tk.M).GetInt("bulan")
				if month == valMonth && project == valProject {
					lossPerMonth = val.GetFloat64("result")
					for field, down := range downCause {
						valResultType := val.GetFloat64(field) / 1000 /* jadikan MWh */
						valTitle := strings.Replace(down, " ", "", -69)
						if valResultType >= 0 {
							resVal.Set(valTitle, valResultType) /* MachineDown : 7.6666 */
							totalLoss[valTitle] += valResultType
						}
					}
				}
			}
			resVal.Set("Total", lossPerMonth/1000) /* jadikan MWh */
			result = append(result, resVal)
		}
		for keyTotal, total := range totalLoss {
			totalLossPerType = append(totalLossPerType, tk.M{
				"name":  keyTotal,
				"value": total / 1000,
			})
		}
		resultPerProject.Set(project, result)
		totalLossPerYear[project] = totalLossPerType
	}
	results.Set("datachart", resultPerProject)
	results.Set("datapie", totalLossPerYear)
	return helper.CreateResult(true, results, "success")
}

func (m *DashboardController) GetDetailLossLevel2(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	result := []tk.M{}

	p := tk.M{}
	e := k.GetPayload(&p)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	var pipes []tk.M
	match := tk.M{}
	match.Set("dateinfo.monthdesc", p.GetString("date"))

	if p.GetString("project") != "Fleet" {
		match.Set("projectname", p.GetString("project"))
	}

	pipes = append(pipes, tk.M{"$match": match})
	pipes = append(pipes, tk.M{"$group": tk.M{
		"_id":         "$turbine",
		"result":      tk.M{"$sum": "$lostenergy"},
		"machinedown": tk.M{"$sum": "$machinedownloss"},
		"griddown":    tk.M{"$sum": "$griddownloss"},
		"unknown":     tk.M{"$sum": "$otherdownloss"},
	}})
	pipes = append(pipes, tk.M{"$sort": tk.M{"_id": 1}})

	csr, e := DB().Connection.NewQuery().
		Select("_id").
		From(new(ScadaSummaryDaily).TableName()).
		Command("pipe", pipes).
		Cursor(nil)

	if e != nil {
		return helper.CreateResult(false, result, e.Error())
	}

	allLossData := []tk.M{}
	e = csr.Fetch(&allLossData, 0, false)
	csr.Close()

	if e != nil {
		return helper.CreateResult(false, result, e.Error())
	}

	turbineList := []string{}
	turbines, _ := helper.GetTurbineList([]interface{}{p.GetString("project")})

	for _, turbine := range turbines {
		turbineList = append(turbineList, turbine.Value)
	}
	sort.Strings(turbineList)

	downCause, _ := getMachineDownType()

	project := p.GetString("project")
	if p.GetString("project") == "Fleet" {
		project = ""
	}
	turbineName, _ := helper.GetTurbineNameList(project)
	for _, turbine := range turbineList {
		resVal := tk.M{}
		resVal.Set("_id", turbineName[turbine])
		for _, down := range downCause {
			valTitle := strings.Replace(down, " ", "", -69)
			resVal.Set(valTitle, 0.0) /* MachineDown : 0.0 ==> default value */
		}
		lossPerTurbine := 0.0
		for _, val := range allLossData {
			valTurbine := val.GetString("_id")
			if turbine == valTurbine {
				lossPerTurbine = val.GetFloat64("result")
				for field, down := range downCause {
					valResultType := val.GetFloat64(field) / 1000 /* jadikan MWh */
					valTitle := strings.Replace(down, " ", "", -69)
					if valResultType >= 0 {
						resVal.Set(valTitle, valResultType) /* MachineDown : 7.6666 */
					}
				}
			}
		}

		resVal.Set("Total", lossPerTurbine/1000) /* ubah jadi MWh */
		result = append(result, resVal)
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
		var wg sync.WaitGroup
		var mux sync.Mutex
		if p.ProjectName != "Fleet" {
			wg.Add(2)
			//tidak bisa dicombine karena tiap top 10 kategori beda urutan top 10 nya
			go getTurbineDownTimeTop(result, "duration", p, k, &wg, &mux)
			go getTurbineDownTimeTop(result, "frequency", p, k, &wg, &mux)
		}
		if p.Type == "project" {
			wg.Add(1)
			go getTurbineDownTimeTop(result, p.Type, p, k, &wg, &mux)
		}
		wg.Add(1)
		go getTurbineDownTimeTop(result, "loss", p, k, &wg, &mux)
		wg.Wait()
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
		var lossD, lossF, loss, dataSeries []tk.M
		if p.ProjectName != "Fleet" {
			lossD, lossF, loss = getLossCategoriesTopDFP(p, k)
		} else {
			lossD, lossF, loss, dataSeries = getLossCategoriesTopStack(p, k)
		}
		result.Set("lossCatDuration", lossD)
		result.Set("lossCatFrequency", lossF)
		result.Set("lossCatLoss", loss)
		result.Set("dataseries", dataSeries)
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

	result.Set("lostenergy", getDownTimeProjectMonthly("fleetdowntime", p, k))

	return helper.CreateResult(true, result, "success")
}

func getDownTimeProjectMonthly(tipe string, p *PayloadDashboard, k *knot.WebContext) (result []tk.M) {
	splitted := strings.Split(p.DateStr, " ")
	bulanStr := splitted[0]
	tahunStr := splitted[1]
	awal := time.Time{}
	for i := 1; i <= 12; i++ {
		if time.Month(i).String() == bulanStr {
			awal = time.Date(tk.ToInt(tahunStr, tk.RoundingAuto), time.Month(i), 01, 0, 0, 0, 0, time.UTC)
			break
		}
	}
	akhir := awal.AddDate(0, 1, -1)
	pAnalytic := new(PayloadAnalytic)
	pAnalytic.Project = p.ProjectName
	pAnalytic.DateStart = awal
	pAnalytic.DateEnd = akhir
	pAnalytic.Period = "custom"

	lossData, _ := getLossDuration("loss", pAnalytic, k)
	downCause, _ := getMachineDownType()
	sortedDown := []string{}
	for key := range downCause {
		sortedDown = append(sortedDown, key)
	}
	sort.Strings(sortedDown)
	result = []tk.M{}
	if len(lossData) > 0 {
		for _, field := range sortedDown {
			result = append(result, tk.M{
				"result": lossData[0].GetFloat64(field) / 1000,
				"type":   downCause[field],
			})
		}
	}

	return
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

		for field, title := range stack {
			if tipe != "type" {

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

func getTurbineDownTimeTop(result tk.M, topType string, p *PayloadDashboard, k *knot.WebContext, wg *sync.WaitGroup, mux *sync.Mutex) {
	defer wg.Done()
	if p.DateStr == "" || (p.DateStr != "" && topType == "project") {
		var e error
		var fromDate time.Time
		fromDate = p.Date.AddDate(0, -12, 0)
		newPayload := new(PayloadAnalytic)
		if topType != "project" {
			newPayload.DateStart = fromDate.UTC()
			newPayload.DateEnd = p.Date.UTC()
		} else {
			newPayload.DateStr = p.DateStr
		}

		if p.ProjectName != "Fleet" && p.ProjectName != "" {
			newPayload.Project = p.ProjectName
		}
		dataResult := []tk.M{}
		if topType != "frequency" {
			tipe := topType
			if topType == "project" { /* untuk tipe ini hanya get loss saja */
				tipe = "loss"
				newPayload.BreakDown = "$projectname"
			}
			dataResult, e = getDownTimeTopLossDuration(tipe, newPayload, k)
		} else {
			dataResult, e = getDownTimeTopFrequency(newPayload, k)
		}
		if e != nil {
			return
		}
		mux.Lock()
		result.Set(topType, dataResult)
		mux.Unlock()
	}
}

func getLossCategoriesFreq(matchSource tk.M, downCause map[string]string, val string) (resLoop []tk.M) {
	pipes := []tk.M{}
	match := tk.M{}
	for key, valMatch := range matchSource {
		if key == "detail.startdate" {
			key = "startdate"
		}
		match.Set(key, valMatch)
	}

	loopMatch := match
	field := val
	title := downCause[val]

	loopMatch.Set(field, true)

	pipes = append(pipes, tk.M{"$match": loopMatch})
	pipes = append(pipes,
		tk.M{
			"$group": tk.M{
				"_id":  tk.M{"id1": field, "id2": title, "project": "$projectname"},
				"freq": tk.M{"$sum": 1}},
		},
	)

	csr, e := DB().Connection.NewQuery().
		From(new(Alarm).TableName()).
		Command("pipe", pipes).
		Cursor(nil)

	if e != nil {
		return
	}

	e = csr.Fetch(&resLoop, 0, false)
	if e != nil {
		return
	}

	csr.Close()

	return

}

func workerLossDur(lossDurData *[]tk.M, newP *PayloadAnalytic, k *knot.WebContext, wg *sync.WaitGroup, mux *sync.Mutex) {
	defer wg.Done()
	data, e := getLossDuration("both", newP, k)
	if e != nil {
		return
	}
	mux.Lock()
	for _, val := range data {
		*lossDurData = append(*lossDurData, val)
	}
	mux.Unlock()
}

func workerFreq(field string, downCause map[string]string, resFreqAsync tk.M, p *PayloadDashboard, wg *sync.WaitGroup, mux *sync.Mutex) {
	defer wg.Done()
	var fromDate time.Time
	match := tk.M{}
	fromDate = p.Date.AddDate(0, -12, 0)
	match.Set("startdate", tk.M{"$gte": fromDate.UTC(), "$lte": p.Date.UTC()})
	match.Set("reduceavailability", true)

	if p.ProjectName != "Fleet" {
		match.Set("projectname", p.ProjectName)
	}
	match.Set(field, true)
	resLoopFreq := getLossCategoriesFreq(match, downCause, field)

	tmpData := tk.M{}
	for _, res := range resLoopFreq {
		resID, _ := tk.ToM(res["_id"])
		tmpData.Set(resID.GetString("project"), res.GetFloat64("freq"))
	}
	mux.Lock()
	resFreqAsync.Set(field, tmpData)
	mux.Unlock()
}

func getLossCategoriesTopStack(p *PayloadDashboard, k *knot.WebContext) (resultDuration, resultFreq, resultPowerLost, dataSeries []tk.M) {
	if p.DateStr == "" {
		downCause, _ := getMachineDownType()
		sortedDown := []string{}
		for key := range downCause {
			sortedDown = append(sortedDown, key)
		}
		sort.Strings(sortedDown)

		var fromDate time.Time
		fromDate = p.Date.AddDate(0, -12, 0)

		newP := new(PayloadAnalytic)
		newP.DateStart = fromDate.UTC()
		newP.DateEnd = p.Date.UTC()
		newP.Period = "custom"
		if p.ProjectName != "Fleet" {
			newP.Project = p.ProjectName
		}
		lossDurData := []tk.M{}
		resFreqAsync := tk.M{}
		var wg sync.WaitGroup
		var mux sync.Mutex

		wg.Add(len(sortedDown) + 1)
		go workerLossDur(&lossDurData, newP, k, &wg, &mux)

		for _, field := range sortedDown {
			go workerFreq(field, downCause, resFreqAsync, p, &wg, &mux)
		}
		wg.Wait()

		tmpResultFreq := tk.M{}
		tmpResultPowerLost := tk.M{}
		tmpResultDuration := tk.M{}
		projectList := map[string]int{}

		for _, val := range sortedDown {
			field := val
			title := downCause[val]

			tmpResultFreq = tk.M{
				"_id": tk.M{"id2": title},
			}
			tmpResultPowerLost = tk.M{
				"_id": tk.M{"id2": title},
			}
			tmpResultDuration = tk.M{
				"_id": tk.M{"id2": title},
			}
			for _, val := range lossDurData {
				projectList[val.GetString("_id")] = 1
				tmpResultPowerLost.Set(val.GetString("_id"), val.GetFloat64(field+"loss"))
				tmpResultDuration.Set(val.GetString("_id"), val.GetFloat64(field+"duration"))
			}

			for field, res := range resFreqAsync.Get(field, tk.M{}).(tk.M) {
				tmpResultFreq.Set(field, res)
			}
			resultFreq = append(resultFreq, tmpResultFreq)
			resultDuration = append(resultDuration, tmpResultDuration)
			resultPowerLost = append(resultPowerLost, tmpResultPowerLost)
		}
		projectSorted := []string{}
		for key := range projectList {
			projectSorted = append(projectSorted, key)
		}
		sort.Strings(projectSorted)

		for _, key := range projectSorted {
			dataSeries = append(dataSeries, tk.M{
				"field": key,
				"name":  key,
			})
		}
	}

	return
}

func getLossCategoriesTopDFP(p *PayloadDashboard, k *knot.WebContext) (resultDuration, resultFreq, resultPowerLost []tk.M) {
	if p.DateStr == "" {
		var fromDate time.Time
		fromDate = p.Date.AddDate(0, -12, 0)

		downCause, _ := getMachineDownType()
		sortedDown := []string{}
		for key := range downCause {
			sortedDown = append(sortedDown, key)
		}
		sort.Strings(sortedDown)

		newP := new(PayloadAnalytic)
		newP.DateStart = fromDate.UTC()
		newP.DateEnd = p.Date.UTC()
		newP.Period = "custom"
		if p.ProjectName != "Fleet" {
			newP.Project = p.ProjectName
		}
		lossDurData := []tk.M{}
		resFreqAsync := tk.M{}
		var wg sync.WaitGroup
		var mux sync.Mutex

		wg.Add(len(sortedDown) + 1)
		go workerLossDur(&lossDurData, newP, k, &wg, &mux)

		for _, field := range sortedDown {
			go workerFreq(field, downCause, resFreqAsync, p, &wg, &mux)
		}
		wg.Wait()

		tmpResultPowerLost := tk.M{}
		tmpResultDuration := tk.M{}
		tmpResultFreq := tk.M{}

		for _, val := range sortedDown {
			field := val
			title := downCause[val]

			tmpResultFreq = tk.M{
				"_id": tk.M{"id2": title},
			}
			tmpResultPowerLost = tk.M{
				"_id": tk.M{"id2": title},
			}
			tmpResultDuration = tk.M{
				"_id": tk.M{"id2": title},
			}
			for _, val := range lossDurData {
				tmpResultPowerLost.Set("result", val.GetFloat64(field+"loss"))
				tmpResultDuration.Set("result", val.GetFloat64(field+"duration"))
			}
			for _, res := range resFreqAsync.Get(field, tk.M{}).(tk.M) {
				tmpResultFreq.Set("result", res)
			}
			resultFreq = append(resultFreq, tmpResultFreq)
			resultDuration = append(resultDuration, tmpResultDuration)
			resultPowerLost = append(resultPowerLost, tmpResultPowerLost)
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

	// log.Printf(">>> %v \n", totalTurbine)

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
	// p.Date, _ = time.Parse("2006-01-02 15:04:05", p.Date.UTC().Format("2006-01")+"-01"+" 00:00:00")

	projects := []string{}

	if p.ProjectName != "Fleet" {
		projects = append(projects, p.ProjectName)
	} else {
		projectList, _ := helper.GetProjectList()
		for _, v := range projectList {
			projects = append(projects, v.Value)
		}
	}
	sort.Strings(projects)

	// rprojects := tk.M{}

	if p.DateStr == "" {
		fromDate = p.Date.AddDate(0, -12, 0)
		if fromDate.Format("20060102")[6:] != "01" {
			fromDate, _ = time.Parse("20060102_150405", fromDate.Format("200601")+"01"+"_000000")
		}

		match.Set("dateinfo.dateid", tk.M{"$gte": fromDate.UTC(), "$lte": p.Date.UTC()})
		// match.Set("available", 1)

		if p.ProjectName != "Fleet" {
			match.Set("projectname", p.ProjectName)
		}

		group := tk.M{
			"_id":           tk.M{"id1": "$dateinfo.monthid", "id2": "$dateinfo.monthdesc", "id3": "$projectname"},
			"count":         tk.M{"$sum": 1},
			"mindate":       tk.M{"$min": "$dateinfo.dateid"},
			"maxdate":       tk.M{"$max": "$dateinfo.dateid"},
			"machineResult": tk.M{"$sum": "$machinedownhours"},
			"gridResult":    tk.M{"$sum": "$griddownhours"},
			"unknownResult": tk.M{"$sum": "$otherdowntimehours"},
		}

		pipe := []tk.M{
			{"$match": match},
			{"$group": group},
			{"$sort": tk.M{"_id.id1": -1}},
			// {"$limit": 12},
		}

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

		// --------------
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
				id := scada.Get("_id").(tk.M)
				project := id.GetString("id3")
				m := scada.GetFloat64("machineResult")
				g := scada.GetFloat64("gridResult")
				u := scada.GetFloat64("unknownResult")

				minDate := scada.Get("mindate").(time.Time).UTC()
				maxDate := scada.Get("maxdate").(time.Time).UTC()
				// minutes := scada.GetFloat64("minutes") / 60

				// fromDateSub, _ := time.Parse("060102_150405", minDate.Format("0601")+"01_000000")
				// tmpDt, _ := time.Parse("060102_150405", minDate.AddDate(0, 1, 0).Format("0601")+"01_000000")
				// toDateSub := tmpDt.AddDate(0, 0, -1)

				turbineList, _ := helper.GetTurbineList([]interface{}{project})
				totalTurbine := float64(len(turbineList))

				// hourValue := helper.GetHourValue(fromDateSub.UTC(), toDateSub.UTC(), minDate.UTC(), maxDate.UTC())
				hourValue := maxDate.AddDate(0, 0, 1).UTC().Sub(minDate.UTC()).Hours()
				// mAvail, gAvail, _, _, _ := helper.GetAvailAndPLF(totalTurbine, float64(0), float64(0), m, g, float64(0), hourValue, minutes, float64(0))

				in := tk.M{}.Set("noofturbine", totalTurbine).Set("oktime", 0).Set("energy", 0).
					Set("totalhour", hourValue).Set("totalcapacity", 0).
					Set("machinedowntime", m).Set("griddowntime", g).Set("otherdowntime", u)

				res := helper.CalcAvailabilityAndPLF(in)

				scada.Set("machineResult", res.GetFloat64("machineavailability")/100)
				scada.Set("gridResult", res.GetFloat64("gridavailability")/100)
				// log.Printf(">>> %#v \n", scada)

				// log.Printf("SCADA: %v | %v | %v | %v = %v | %v - %v - %v - %v \n", minutes, res/3600.0, totalTurbine, hourValue, tk.ToFloat64(avail, 2, tk.RoundingAuto), fromDate.UTC().String(), p.Date.UTC().String(), minDate.UTC().String(), maxDate.UTC().String())
			}

			scada.Unset("maxdate")
			scada.Unset("mindate")
			scada.Unset("minutes")
		}
	}
	if p.ProjectName != "Fleet" {
		for _, res := range result {
			machineResult = append(machineResult, tk.M{"_id": res.Get("_id"), "result": res.GetFloat64("machineResult")})
			gridResult = append(gridResult, tk.M{"_id": res.Get("_id"), "result": res.GetFloat64("gridResult")})
		}
	} else {
		orderNo := 0
		ids := tk.M{}
		for _, res := range result {
			orderNo++
			ids, _ = tk.ToM(res.Get("_id"))
			var ima, iga interface{}
			ima, iga = res.GetFloat64("machineResult"), res.GetFloat64("gridResult")
			if res.GetFloat64("count") == 0 {
				ima, iga = nil, nil
			}
			machineResult = append(machineResult, tk.M{
				"DataId":  tk.ToString(ids.GetInt("id1")),
				"Title":   ids.GetString("id2"),
				"OrderNo": orderNo,
				"Value":   ima,
				"Project": ids.GetString("id3"),
			})
			gridResult = append(gridResult, tk.M{
				"DataId":  tk.ToString(ids.GetInt("id1")),
				"Title":   ids.GetString("id2"),
				"OrderNo": orderNo,
				"Value":   iga,
				"Project": ids.GetString("id3"),
			})
		}
	}

	mrTmp := []tk.M{}
	grTmp := []tk.M{}
	// tk.Println(len(machineResult), " > ", (len(projects) * 12))
	if len(machineResult) >= (len(projects) * 12) {
		length := len(machineResult)
		div := (length / len(projects))

		for i := 1; i < len(projects)+1; i++ {
			offerX := (div * i) - 12
			// log.Printf("> %v | %v | %v | %v \n", div, div*i, offerX, offerX+12)
			// tk.Println(offerX+12, " > ", len(gridResult), " || ", len(machineResult))
			// if offerX+12 > len(gridResult) || offerX+12 > len(machineResult) {
			// 	continue
			// }

			if p.ProjectName != "Fleet" {
				mrTmp = append(mrTmp, machineResult[offerX:offerX+12]...)
				grTmp = append(grTmp, gridResult[offerX:offerX+12]...)
			} else {
				mrTmp = append(mrTmp, tk.M{
					"Project": machineResult[offerX].GetString("Project"),
					"Details": machineResult[offerX : offerX+12],
				})
				grTmp = append(grTmp, tk.M{
					"Project": gridResult[offerX].GetString("Project"),
					"Details": gridResult[offerX : offerX+12],
				})
			}
		}
		machineResult = mrTmp
		gridResult = grTmp
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

	turbineName, e := helper.GetTurbineNameList(p.ProjectName)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	for key, val := range turbineName {
		if p.Turbine == val {
			p.Turbine = key
		}
	}
	tipe := strings.Split(p.Type, "_")
	if len(tipe) < 2 {
		return helper.CreateResult(false, nil, e.Error())
	}

	fromDate = p.Date.AddDate(0, -12, 0)
	tipeDown := "$" + strings.ToLower(tipe[0])
	tableName := new(ScadaSummaryDaily).TableName()
	if tipe[1] == "Hours" {
		if tipeDown == "$unknown" {
			tipeDown = "$otherdowntimehours"
		} else {
			tipeDown += "hours"
		}
		pipes = []tk.M{
			tk.M{
				"$match": tk.M{
					"turbine":         p.Turbine,
					"dateinfo.dateid": tk.M{"$gte": fromDate.UTC(), "$lte": p.Date.UTC()},
				},
			},
			tk.M{
				"$group": tk.M{"_id": tk.M{"id1": "$dateinfo.monthid", "id2": "$dateinfo.monthdesc"},
					"result": tk.M{"$sum": tipeDown},
				},
			},
		}
	} else if tipe[1] == "Times" {
		pipes = []tk.M{
			tk.M{
				"$match": tk.M{
					"turbine":                p.Turbine,
					"startdate":              tk.M{"$gte": fromDate.UTC(), "$lte": p.Date.UTC()},
					strings.ToLower(tipe[0]): true,
				},
			},
			tk.M{
				"$group": tk.M{"_id": tk.M{"id1": "$startdateinfo.monthid", "id2": "$startdateinfo.monthdesc"},
					"result": tk.M{"$sum": 1},
				},
			},
		}
		tableName = new(Alarm).TableName()
	} else if tipe[1] == "MWh" {
		if tipeDown == "$unknown" {
			tipeDown = "$otherdownloss"
		} else {
			tipeDown += "loss"
		}
		pipes = []tk.M{
			tk.M{
				"$match": tk.M{
					"turbine":         p.Turbine,
					"dateinfo.dateid": tk.M{"$gte": fromDate.UTC(), "$lte": p.Date.UTC()},
				},
			},
			tk.M{
				"$group": tk.M{"_id": tk.M{"id1": "$dateinfo.monthid", "id2": "$dateinfo.monthdesc"},
					"result": tk.M{"$sum": tipeDown},
				},
			},
		}
	}

	pipes = append(pipes, tk.M{"$sort": tk.M{"_id.id1": 1}})

	csr, e := DB().Connection.NewQuery().
		From(tableName).
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

func (m *DashboardController) GetWindDistributionRev(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	var dataSeries []tk.M

	type PayloadWindDist struct {
		ProjectName string
		Date        time.Time
		PeriodList  []string
	}

	p := new(PayloadWindDist)
	e := k.GetPayload(&p)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	csr, _ := DB().
		Connection.
		NewQuery().
		Select("Contribute", "Project", "Category").
		From("rpt_winddistributioncurrentmonth").
		Cursor(nil)

	defer csr.Close()

	for {
		tkm := tk.M{}
		e = csr.Fetch(&tkm, 1, false)
		if e != nil {
			break
		}

		dataSeries = append(dataSeries, tkm)
	}

	result := tk.M{}
	result["currentmonth"] = dataSeries
	data := struct {
		Data tk.M
	}{
		Data: result,
	}

	return helper.CreateResult(true, data, "success")

}

func (m *DashboardController) GetDownTimeTurbines(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	p := new(PayloadDashboard)
	e := k.GetPayload(&p)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	// result := getDownTurbine(p.ProjectName, p.Date, 1)
	result := getDownTurbineStatus(p.ProjectName, p.Date, 1)

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

func getTotalDownTurbine(project string, currentDate time.Time, dayDuration int) (result int) {
	// var fromDate time.Time
	var pipes []tk.M
	match := tk.M{}
	// currentDate = getTimeNow()

	// fromDate = currentDate.UTC().AddDate(0, 0, dayDuration*-1)

	// match.Set("datestart", tk.M{"$gte": fromDate.UTC()})
	match.Set("status", 0)

	if project != "Fleet" {
		match.Set("projectname", project)
	}

	pipes = append(pipes, tk.M{"$match": match})
	pipes = append(pipes, tk.M{"$sort": tk.M{"_id": 1}})

	rconn := DBRealtime()

	hasil := []tk.M{}
	csr, e := rconn.NewQuery().
		From(new(TurbineStatus).TableName()).
		Command("pipe", pipes).
		Cursor(nil)

	if e != nil {
		return
	}
	e = csr.Fetch(&hasil, 0, false)
	if e != nil {
		return
	}
	defer csr.Close()
	result = len(hasil)

	return
}

func getDownTurbineStatus(project string, currentDate time.Time, dayDuration int) (result []tk.M) {
	var fromDate time.Time
	var pipes []tk.M
	match := tk.M{}

	currentDate = getTimeNow()
	fromDate = currentDate.UTC().AddDate(0, 0, dayDuration*-1)

	if dayDuration > 1 {
		match.Set("datestart", tk.M{"$gte": fromDate.UTC(), "$lte": currentDate.UTC()})
	}
	match.Set("status", tk.M{"$eq": 0})

	if project != "Fleet" {
		match.Set("projectname", project)
	}

	pipes = append(pipes, tk.M{"$match": match})
	pipes = append(pipes, tk.M{"$sort": tk.M{"_id": 1}})

	rconn := DBRealtime()

	csr, e := rconn.NewQuery().
		From(new(TurbineStatus).TableName()).
		Command("pipe", pipes).
		Cursor(nil)

	defer csr.Close()

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

	lastProject := ""
	turbineName := map[string]string{}
	for _, val := range tmpResult {
		if lastProject != val.GetString("projectname") {
			lastProject = val.GetString("projectname")
			turbineName, _ = helper.GetTurbineNameList(lastProject)
		}
		if val.Get("datestart") != nil {
			start := val.Get("datestart").(time.Time)
			downHours := currentDate.UTC().Sub(start.UTC()).Hours()
			if dayDuration > 1 {
				val.Set("_id", turbineName[val.GetString("turbine")])
				if downHours >= float64(24*dayDuration) {
					val.Set("result", downHours)
					val.Set("isdown", true)
					result = append(result, val)
				}
			} else {
				val.Set("_id", turbineName[val.GetString("turbine")])
				val.Set("result", downHours)
				val.Set("isdown", true)
				result = append(result, val)
			}
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

func setMapData() (result tk.M) {
	// initiate all variables
	result = tk.M{}
	// set database to realtime data db
	rconn := DBRealtime()
	t0, servt0 := getTimeNow(), time.Now().UTC()

	pipes := []tk.M{
		tk.M{"$match": tk.M{"projectname": tk.M{"$ne": ""}}}}
	pipes = append(pipes, tk.M{"$group": tk.M{
		"_id":            tk.M{"projectname": "$projectname", "turbine": "$turbine"},
		"lastupdated":    tk.M{"$max": "$timestamp"},
		"lasttimeserver": tk.M{"$max": "$servertimestamp"},
	}})
	pipes = append(pipes, tk.M{
		"$sort": tk.M{
			"_id.projectname": 1,
		},
	})

	csrNa, err := rconn.NewQuery().From(new(ScadaRealTimeNew).TableName()).Command("pipe", pipes).Cursor(nil)
	if err != nil {
		tk.Println(err.Error())
	}
	defer csrNa.Close()
	lastUpdateRealtime := []tk.M{}
	err = csrNa.Fetch(&lastUpdateRealtime, 0, false)
	if err != nil {
		tk.Println(err.Error())
	}
	arrturbinestatus := GetTurbineStatus("", "")
	// get no of turbine waiting for wind status
	waitingForWs := getDataPerTurbine("_waitingforwindspeed", tk.M{
		"$and": []tk.M{
			tk.M{"status": true},
		}}, false)

	waitingForWsProject := 0
	dataNa := 0
	dataDowns, greyDowns := 0, 0
	_tTurbine := ""
	_tProject := ""
	isDataComing := false
	var tstamp, servtstamp time.Time
	keys := ""
	lastProject := ""
	turbineStatus := map[string]string{}
	turbineDownList, turbineDownListByProject := []tk.M{}, []tk.M{}
	currturbine, turbineName := tk.M{}, map[string]string{}
	currentDate := getTimeNow()
	downPerProject := map[string]int{}

	for _, dt := range lastUpdateRealtime {
		ids, _ := tk.ToM(dt.Get("_id"))
		tstamp = dt.Get("lastupdated", time.Time{}).(time.Time).UTC()
		servtstamp = dt.Get("lasttimeserver", time.Time{}).(time.Time).UTC()
		_tTurbine = ids.GetString("turbine")
		_tProject = ids.GetString("projectname")
		if lastProject != _tProject {
			if lastProject != "" {
				if len(currturbine) == dataNa {
					dataDowns = 0
					turbineDownListByProject = []tk.M{}
					downPerProject[lastProject] = 0
				} else {
					dataNa = dataNa - greyDowns
				}

				if len(turbineDownListByProject) > 0 {
					turbineDownList = append(turbineDownList, turbineDownListByProject...)
				}

				result.Set(lastProject, tk.M{
					"allna":       len(currturbine) == dataNa,
					"grey":        dataNa,
					"orange":      waitingForWsProject,
					"red":         dataDowns,
					"turbineList": turbineStatus,
				})
			}

			currturbine = tk.M{}
			downPerProject[_tProject] = 0
			turbineName, _ = helper.GetTurbineNameList(_tProject)
			lastProject = _tProject
			turbineStatus = map[string]string{}
			waitingForWsProject = 0
			dataNa, dataDowns, greyDowns = 0, 0, 0
			turbineDownListByProject = []tk.M{}
		}
		currturbine.Set(_tTurbine, 1)
		turbineStatus[_tTurbine] = "green"
		limitVal, hasLimit := NotAvailLimit[_tProject]
		if hasLimit && (t0.Sub(tstamp.UTC()).Minutes() <= limitVal || servt0.Sub(servtstamp.UTC()).Minutes() <= limitVal) {
			isDataComing = true
		} else {
			turbineStatus[_tTurbine] = "grey"
			isDataComing = false
			dataNa++
		}
		keys = _tProject + "_" + _tTurbine

		if _idt, _cond := arrturbinestatus[_tTurbine]; _cond {
			if _idt.Status == 0 {
				downHours := currentDate.UTC().Sub(_idt.DateStart.UTC()).Hours()
				dtDown := tk.M{
					"_id":    turbineName[_idt.Turbine],
					"result": downHours,
					"isdown": true,
					"color":  "red",
				}
				turbineStatus[_tTurbine] = "red"

				if !isDataComing {
					turbineStatus[_tTurbine] = "grey"
					greyDowns++
					dtDown.Set("color", "grey")
				}

				downPerProject[_tProject]++
				turbineDownListByProject = append(turbineDownListByProject, dtDown)
				dataDowns++
			} else if waitingForWs.Has(keys) && isDataComing {
				turbineStatus[_tTurbine] = "orange"
				waitingForWsProject++
			}
		}
	}

	if lastProject != "" {
		if len(currturbine) == dataNa {
			dataDowns = 0
			turbineDownListByProject = []tk.M{}
			downPerProject[lastProject] = 0
		} else {
			dataNa = dataNa - greyDowns
		}

		if len(turbineDownListByProject) > 0 {
			turbineDownList = append(turbineDownList, turbineDownListByProject...)
		}

		result.Set(lastProject, tk.M{
			"allna":       len(currturbine) == dataNa,
			"grey":        dataNa,
			"orange":      waitingForWsProject,
			"red":         dataDowns,
			"turbineList": turbineStatus,
		})
	}

	result.Set("turbineDownList", turbineDownList)
	result.Set("downAll", len(turbineDownList))
	result.Set("downPerProject", downPerProject)

	return
}

func (m *DashboardController) GetMapData(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	payload := map[string]string{}
	e := k.GetPayload(&payload)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}
	projectName := payload["projectname"]
	projectList, _ := helper.GetProjectList()
	projects := []ProjectOut{}

	if projectName != "Fleet" {
		for _, v := range projectList {
			if v.Value == projectName {
				projects = append(projects, v)
				break
			}
		}
	} else {
		projects = projectList
	}

	resultMap := []tk.M{}
	projectTurbineStatus := setMapData()
	projectVal := tk.M{}
	turbineCount, turbinena := 0, 0
	turbineStatus := map[string]string{}

	for _, project := range projects {
		res := []tk.M{}
		stsProj := ""
		turbineList, _ := helper.GetTurbineList([]interface{}{project.Value})
		turbineCount = len(turbineList)
		if projectTurbineStatus.Has(project.Value) {
			projectVal = projectTurbineStatus.Get(project.Value, tk.M{}).(tk.M)
		}
		if projectName != "Fleet" {
			turbineStatus = projectVal.Get("turbineList", map[string]string{}).(map[string]string)
			for _, turbine := range turbineList {
				res = append(res, tk.M{
					"name":   turbine.Turbine,
					"value":  turbine.Value,
					"coords": turbine.Coords,
					"status": turbineStatus[turbine.Value],
				})
			}
			resultMap = res
		} else {
			if projectVal.Get("allna", false).(bool) {
				stsProj = "grey"
			} else if projectVal.GetInt("red") == turbineCount {
				stsProj = "red"
			} else if projectVal.GetInt("orange") == turbineCount {
				stsProj = "orange"
			} else {
				stsProj = "green"
			}
			resultMap = append(resultMap, tk.M{
				"name":   project.Value,
				"value":  project.Value,
				"coords": project.Coords,
				"status": stsProj,
			})
		}

		if projectVal.GetInt("grey") != turbineCount {
			turbinena += projectVal.GetInt("grey")
		}
	}

	results := tk.M{}
	results.Set("resultMap", resultMap)
	results.Set("turbineDownList", projectTurbineStatus.Get("turbineDownList"))
	results.Set("totalDownFleet", projectTurbineStatus.GetInt("downAll"))
	results.Set("totalNAFleet", turbinena)
	results.Set("downPerProject", projectTurbineStatus.Get("downPerProject"))

	// probably its temporary solution to handle fatal error: concurrent map writes
	//return helper.CreateResult(true, results, "success")
	return helper.CreateResultWithoutSession(true, results, "success")
}

func (m *DashboardController) GetMapData_old(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	payload := map[string]string{}
	e := k.GetPayload(&payload)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}
	projectName := payload["projectname"]

	mapTurbines := map[string]string{}

	if projectName != "Fleet" {
		result := getDownTurbineStatus(projectName, time.Now(), 0)
		for _, v := range result {
			turbine := v.GetString("_id")
			mapTurbines[turbine] = turbine
		}
	}

	data := getMapCol(projectName)

	results := tk.Ms{}
	offset := []int{0, 2}
	coords := []float64{}
	for _, val := range data {
		result := tk.M{}
		coords = []float64{}
		coords = []float64{val.GetFloat64("latitude"), val.GetFloat64("longitude")}
		status := true
		if projectName == "Fleet" {
			result.Set("name", val.GetString("projectname"))
			result.Set("status", status)
		} else {
			result.Set("name", val.GetString("turbineid"))
			if mapTurbines[val.GetString("turbineid")] != "" {
				status = false
			}
			result.Set("status", status)
		}
		result.Set("coords", coords)
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
	addLoss := []string{}
	addDuration := []string{}
	matches := []tk.M{}
	if p.ProjectName != "" && strings.ToLower(p.ProjectName) != "fleet" {
		matches = append(matches, tk.M{
			"projectname": p.ProjectName,
		})
	}
	downTypeList, _ := getMachineDownType()

	if p.Type != "" {
		switch strings.ToLower(strings.Replace(p.Type, " ", "", 1)) {
		case "machinedown":
			addLoss = append(addLoss, "$machinedownloss")
			addDuration = append(addDuration, "$machinedownduration")
		case "griddown":
			addLoss = append(addLoss, "$griddownloss")
			addDuration = append(addDuration, "$griddownduration")
		case "unknown":
			addLoss = append(addLoss, "$unknownloss")
			addDuration = append(addDuration, "$unknownduration")
		}
	} else {
		for key := range downTypeList {
			addLoss = append(addLoss, "$"+key+"loss")
			addDuration = append(addDuration, "$"+key+"duration")
		}
	}

	if dateStr[0] != "fleet" {
		date, e = time.Parse("Jan 2006 02 15:04:05", dateStr[0][0:3]+" "+dateStr[1]+" 01 00:00:00")

		if e != nil {
			return helper.CreateResult(false, nil, e.Error())
		}

		dateInfo := GetDateInfo(date)

		matches = append(matches, tk.M{"dateinfo.monthid": dateInfo.MonthId})
		pipes = append(pipes, tk.M{"$match": tk.M{"$and": matches}})

	} else {
		dateStr = strings.Split("Jul 2015", " ")
		dateStr2 := strings.Split("Jun 2016", " ")
		date, e = time.Parse("Jan 2006 02 15:04:05", dateStr[0][0:3]+" "+dateStr[1]+" 01 00:00:00")
		date2, e := time.Parse("Jan 2006 02 15:04:05", dateStr2[0][0:3]+" "+dateStr2[1]+" 01 00:00:00")

		if e != nil {
			return helper.CreateResult(false, nil, e.Error())
		}

		matches = append(matches, tk.M{"dateinfo.dateid": tk.M{"$gte": date}})
		matches = append(matches, tk.M{"dateinfo.dateid": tk.M{"$lte": date2}})
		pipes = append(pipes, tk.M{"$match": tk.M{"$and": matches}})
	}

	pipes = append(pipes,
		tk.M{
			"$group": tk.M{
				"_id": "$turbine",
				"machinedownduration": tk.M{"$sum": "$machinedownhours"},
				"griddownduration":    tk.M{"$sum": "$griddownhours"},
				"unknownduration":     tk.M{"$sum": "$otherdowntimehours"},
				"machinedownloss":     tk.M{"$sum": "$machinedownloss"},
				"griddownloss":        tk.M{"$sum": "$griddownloss"},
				"unknownloss":         tk.M{"$sum": "$otherdownloss"},
			},
		},
	)
	pipes = append(pipes,
		tk.M{
			"$project": tk.M{
				"powerlost": tk.M{"$add": addLoss},
				"duration":  tk.M{"$add": addDuration},
			},
		},
	)

	pipes = append(pipes, tk.M{"$sort": tk.M{"powerlost": -1}})

	csr, e := DB().Connection.NewQuery().
		From(new(ScadaSummaryDaily).TableName()).
		Command("pipe", pipes).
		Cursor(nil)

	// add by ams, 2016-10-07
	defer csr.Close()
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	e = csr.Fetch(&result, 0, false)

	if strings.ToLower(p.ProjectName) != "fleet" {
		listturbines := PopulateTurbines(DB().Connection, p.ProjectName)
		for i, itkm := range result {
			if listturbines.Has(itkm.GetString("_id")) {
				itkm.Set("_id", listturbines.GetString(itkm.GetString("_id")))
				result[i] = itkm
			}
		}
	}

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

	listturbine := tk.M{}
	tprojectname := ""
	for i, _alarm := range resTable {
		if _alarm.Farm != tprojectname {
			tprojectname = _alarm.Farm
			listturbine = PopulateTurbines(DB().Connection, tprojectname)
		}

		if listturbine.Has(_alarm.Turbine) {
			_alarm.Turbine = listturbine.GetString(_alarm.Turbine)
			resTable[i] = _alarm
		}
	}

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

	turbineName, e := helper.GetTurbineNameList(p.ProjectName)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	for key, val := range turbineName {
		if p.Turbine == val {
			p.Turbine = key
		}
	}
	tipe := strings.Split(p.Type, "_")

	filter = append(filter, dbox.Eq("turbine", p.Turbine))
	filter = append(filter, dbox.Gte("startdate", fromDate.UTC()))
	filter = append(filter, dbox.Lte("startdate", p.Date.UTC()))
	filter = append(filter, dbox.Eq(strings.ToLower(tipe[0]), true))

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

	pipe := []tk.M{{"$match": matches}}

	group := tk.M{}
	group.Set("_id", tk.M{"id1": "$turbine", "id2": "$projectname"})
	group.Set("machinedownhours", tk.M{"$sum": "$machinedownhours"})
	group.Set("lossenergy", tk.M{"$sum": "$lostenergy"})
	group.Set("production", tk.M{"$sum": "$production"})
	group.Set("oktime", tk.M{"$sum": "$oktime"})
	group.Set("totalrows", tk.M{"$sum": "$totalrows"})
	group.Set("maxDate", tk.M{"$max": "$dateinfo.dateid"})
	group.Set("minDate", tk.M{"$min": "$dateinfo.dateid"})
	pipe = append(pipe, tk.M{"$group": group})

	pipe = append(pipe, tk.M{"$sort": tk.M{"_id": 1}})

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

	listturbine, listcapacity := PopulateTurbines(DB().Connection, projectName), tk.M{}

	listcapacity = PopulateTurbinesCapacity(DB().Connection, projectName)
	databyproject, arrproj := map[string][]tk.M{}, []string{}

	for _, dres := range result {
		_proj := dres.Get("_id").(tk.M).GetString("id2")
		_name := dres.Get("_id").(tk.M).GetString("id1")

		if _, cond := databyproject[_proj]; !cond {
			databyproject[_proj] = []tk.M{}
		}

		if !tk.HasMember(arrproj, _proj) {
			arrproj = append(arrproj, _proj)
		}

		dres.Set("name", _name)
		if listturbine.Has(_name) {
			dres.Set("name", listturbine.GetString(_name))
		}

		machinedownhours := tk.ToFloat64(dres.GetFloat64("machinedownhours"), 2, tk.RoundingAuto)

		lowestMachineAvail := float64(0)
		maxLossEnergy := tk.ToFloat64(dres.GetFloat64("lossenergy"), 2, tk.RoundingAuto)

		minDate := dres.Get("minDate", time.Time{}).(time.Time)
		maxDate := dres.Get("maxDate", time.Time{}).(time.Time)
		totalHours = maxDate.AddDate(0, 0, 1).UTC().Sub(minDate.UTC()).Hours()

		turbineMW = listcapacity.GetFloat64(_name)
		maxCapacity := turbineMW * totalHours

		lowestMachineAvail = tk.Div(totalHours-machinedownhours, totalHours)

		// dres.Set("lowestplf", formatStringFloat(tk.ToString(lowestPlf), 2)+" %")
		dres.Set("lowestmachineavail", formatStringFloat(tk.ToString(lowestMachineAvail), 2)+" %")

		// dres.Set("plffloat", lowestPlf)
		dres.Set("machineavailfloat", lowestMachineAvail)

		dres.Set("maxlossenergy", maxLossEnergy)
		dres.Set("maxcapacity", maxCapacity)
		dres.Set("plf", (dres.GetFloat64("production")/1000000)/(maxCapacity/1000))
		dres.Set("totalavail", tk.Div(dres.GetFloat64("oktime")/3600, totalHours))
		dres.Set("dataavail", tk.Div(dres.GetFloat64("totalrows")/6, totalHours))

		databyproject[_proj] = append(databyproject[_proj], dres)
	}

	lasresult := []tk.M{}
	if projectName != "" {
		lasresult = databyproject[projectName]
	} else {
		for _, proj := range arrproj {
			ltkm := tk.M{}
			ltkm.Set("name", proj)
			noofwtg := float64(len(databyproject[proj]))
			ltkm.Set("noofwtg", len(databyproject[proj]))

			minDate, maxDate := time.Time{}, time.Time{}

			lplf, lmachineavail, llostenergy := float64(0), float64(0), float64(0)
			for _, tkm := range databyproject[proj] {
				ltkm.Set("production", ltkm.GetFloat64("production")+tkm.GetFloat64("production"))
				ltkm.Set("oktime", ltkm.GetFloat64("oktime")+tkm.GetFloat64("oktime"))
				ltkm.Set("totalrows", ltkm.GetFloat64("totalrows")+tkm.GetFloat64("totalrows"))

				iMinDate := tkm.Get("minDate", time.Time{}).(time.Time)
				iMaxDate := tkm.Get("maxDate", time.Time{}).(time.Time)

				if (minDate.IsZero() || minDate.UTC().After(iMinDate.UTC())) && !iMinDate.IsZero() {
					minDate = iMinDate
				}

				if (maxDate.IsZero() || maxDate.UTC().Before(iMaxDate.UTC())) && !iMaxDate.IsZero() {
					maxDate = iMaxDate
				}

				if lmachineavail == 0 || lmachineavail > tkm.GetFloat64("machineavailfloat") {
					lmachineavail = tkm.GetFloat64("machineavailfloat")
					ltkm.Set("lowestmachineavail", formatStringFloat(tk.ToString(lmachineavail*100), 2)+"% ("+tkm.GetString("name")+")")
				}

				if lmachineavail == 1 {
					ltkm.Set("lowestmachineavail", "-")
				}

				if lplf == 0 || lplf > tkm.GetFloat64("plf") {
					lplf = tkm.GetFloat64("plf")
					ltkm.Set("lowestplf", formatStringFloat(tk.ToString(lplf*100), 2)+"% ("+tkm.GetString("name")+")")
				}

				if llostenergy == 0 || llostenergy < tkm.GetFloat64("maxlossenergy") {
					llostenergy = tkm.GetFloat64("maxlossenergy")
					ltkm.Set("maxlossenergy", formatStringFloat(tk.ToString(llostenergy), 2)+" ("+tkm.GetString("name")+")")
				}
			}

			totalHours = maxDate.AddDate(0, 0, 1).UTC().Sub(minDate.UTC()).Hours()
			turbineMW = listcapacity.GetFloat64(proj)
			maxCapacity := turbineMW * totalHours

			ltkm.Set("maxcapacity", maxCapacity)
			ltkm.Set("plf", (ltkm.GetFloat64("production")/1000000)/(maxCapacity/1000))
			ltkm.Set("totalavail", tk.Div(ltkm.GetFloat64("oktime")/3600, totalHours*noofwtg))
			ltkm.Set("dataavail", tk.Div(ltkm.GetFloat64("totalrows")/6, totalHours*noofwtg))

			lasresult = append(lasresult, ltkm)
		}
	}

	data := struct {
		Data  []tk.M
		Total int
	}{
		Data:  lasresult,
		Total: len(lasresult),
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

func (m *DashboardController) GetMonthlyProject(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	// initiate data return
	data := map[string][]tk.M{}

	// get payload
	type PayloadMonthlyProject struct {
		Projects []string
	}

	p := new(PayloadMonthlyProject)
	e := k.GetPayload(&p)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	// initiate last 12 months
	startYear := time.Now().Year()
	startMonth := int(time.Now().Month()) + 1
	if startMonth == 13 {
		startMonth = 1
	} else {
		startYear--
	}

	// getting data last 12 months for all projects
	monthIdFilter := tk.ToInt((tk.ToString(startYear) + LeftPad2Len(tk.ToString(startMonth), "0", 2)), "0")
	csrScada, e := DB().Connection.NewQuery().
		From(new(ScadaSummaryByMonth).TableName()).
		Where(dbox.And(dbox.Ne("projectname", "Fleet"), dbox.Gte("dateinfo.monthid", monthIdFilter))).
		Order("projectname", "dateinfo.monthid").Cursor(nil)
	defer csrScada.Close()
	if e != nil {
		helper.CreateResult(false, nil, e.Error())
	}

	results := []tk.M{}
	e = csrScada.Fetch(&results, 0, false)

	// plot data
	dataPlot := tk.M{}
	if len(results) > 0 {
		for _, res := range results {
			project := res.GetString("projectname")
			dateinfo := res.Get("dateinfo").(tk.M)
			monthid := dateinfo.GetInt("monthid")
			production := tk.Div(res.GetFloat64("production"), 1000.0)
			lostenergy := res.GetFloat64("lostenergy")
			dataPlot.Set(project+"|"+tk.ToString(monthid), tk.M{
				"production": production,
				"lostenergy": lostenergy,
				"plf":        res.GetFloat64("plf"),
				"trueavail":  res.GetFloat64("trueavail"),
			})
		}
	}

	// define data last 12 months for each projects
	if len(p.Projects) > 0 {
		for _, project := range p.Projects {
			for i := 0; i < 12; i++ {
				monthId := tk.ToInt((tk.ToString(startYear) + LeftPad2Len(tk.ToString(startMonth), "0", 2)), "0")

				production := 0.0
				lostenergy := 0.0
				plf := 0.0
				trueavail := 0.0
				if dataPlot.Has(project + "|" + tk.ToString(monthId)) {
					dtp := dataPlot[project+"|"+tk.ToString(monthId)]
					production = dtp.(tk.M).GetFloat64("production")
					lostenergy = dtp.(tk.M).GetFloat64("lostenergy")
					plf = dtp.(tk.M).GetFloat64("plf")
					trueavail = dtp.(tk.M).GetFloat64("trueavail")
				}

				dateInfo := MonthIDToDateInfo(monthId)
				data[project] = append(data[project], tk.M{
					"monthid":      dateInfo.MonthId,
					"monthdesc":    dateInfo.MonthDesc,
					"production":   production,
					"lostenergy":   lostenergy,
					"plf":          plf,
					"availability": trueavail,
				})

				startMonth++
				if startMonth == 13 {
					startMonth = 1
					startYear++
				}
			}

			startYear = time.Now().Year()
			startMonth = int(time.Now().Month()) + 1
			if startMonth == 13 {
				startMonth = 1
			} else {
				startYear--
			}
		}
	}

	return helper.CreateResult(true, data, "success")
}
