package controller

import (
	. "eaciit/wfdemo-git/library/core"
	. "eaciit/wfdemo-git/library/models"
	"eaciit/wfdemo-git/web/helper"
	"fmt"
	"math"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	c "github.com/eaciit/crowd"
	"github.com/eaciit/dbox"
	"github.com/eaciit/knot/knot.v1"
	tk "github.com/eaciit/toolkit"
)

type AnalyticLossAnalysisController struct {
	App
	mux sync.Mutex
}

var colorFields = [...]string{"#ff880e", "#21c4af", "#ff7663", "#ffb74f", "#a2df53", "#1c9ec4", "#ff63a5", "#f44336", "#D91E18", "#8877A9", "#9A12B3", "#26C281", "#E7505A", "#C49F47", "#ff5597", "#c3260c", "#d4735e", "#ff2ad7", "#34ac8b", "#11b2eb", "#f35838", "#ff0037", "#507ca3", "#ff6565", "#ffd664", "#72aaff", "#795548"}

func CreateAnalyticLossAnalysisController() *AnalyticLossAnalysisController {
	var controller = new(AnalyticLossAnalysisController)
	return controller
}

func (m *AnalyticLossAnalysisController) GetScadaSummaryList(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	var (
		pipes []tk.M
	)

	p := new(PayloadAnalytic)
	e := k.GetPayload(&p)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	// tStart, _ := time.Parse("2006-01-02", p.DateStart.UTC().Format("2006-01-02"))
	// tEnd, _ := time.Parse("2006-01-02 15:04:05", p.DateEnd.UTC().Format("2006-01-02")+" 23:59:59")
	tStart, tEnd, e := helper.GetStartEndDate(k, p.Period, p.DateStart, p.DateEnd)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}
	turbine := p.Turbine
	project := p.Project
	matches := []tk.M{}
	matches = append(matches, tk.M{"dateinfo.dateid": tk.M{"$gte": tStart}})
	matches = append(matches, tk.M{"dateinfo.dateid": tk.M{"$lte": tEnd}})

	if project != "" {
		matches = append(matches, tk.M{"projectname": project})
	}

	if len(turbine) != 0 {
		matches = append(matches, tk.M{"turbine": tk.M{"$in": turbine}})
	}

	breakdown := "Turbine"
	ids := "$turbine"
	if project == "" {
		ids = "$projectname"
		breakdown = "Project"
	}
	pipes = append(pipes, tk.M{"$match": tk.M{"$and": matches}})
	// Aggr(dbox.AggrMax, "$dateinfo.dateid", "max").
	// Aggr(dbox.AggrMin, "$dateinfo.dateid", "min").
	pipes = append(pipes, tk.M{"$group": tk.M{"_id": ids,
		"Production":       tk.M{"$sum": "$production"},
		"MachineDownLoss":  tk.M{"$sum": "$machinedownloss"},
		"GridDownLoss":     tk.M{"$sum": "$griddownloss"},
		"PCDeviation":      tk.M{"$sum": "$pcdeviation"},
		"ElectricLoss":     tk.M{"$sum": "$electricallosses"},
		"OtherDownLoss":    tk.M{"$sum": "$otherdownloss"},
		"DownTimeDuration": tk.M{"$sum": "$downtimehours"},
		"MachineDownHours": tk.M{"$sum": "$machinedownhours"},
		"GridDownHours":    tk.M{"$sum": "$griddownhours"},
		"OtherDownHours":   tk.M{"$sum": "$otherdowntimehours"},
		"maxdate":          tk.M{"$max": "$dateinfo.dateid"},
		"mindate":          tk.M{"$min": "$dateinfo.dateid"},
		"scadaavail":       tk.M{"$avg": "$scadaavail"},
		"totalrows":        tk.M{"$sum": "$totalrows"},
		"LossEnergy":       tk.M{"$sum": "$lostenergy"}}})

	csr, e := DB().Connection.NewQuery().
		From(new(ScadaSummaryDaily).TableName()).
		Command("pipe", pipes).
		Cursor(nil)

	if e != nil {
		helper.CreateResult(false, nil, e.Error())
	}
	defer csr.Close()

	resultScada := []tk.M{}
	e = csr.Fetch(&resultScada, 0, false)
	if e != nil {
		helper.CreateResult(false, nil, e.Error())
	}

	// maxDate := data.Get("maxdate", time.Time{}).(time.Time)
	// minDate := data.Get("mindate", time.Time{}).(time.Time).AddDate(0, 0, 1)
	availability := getAvailabilityValue(tStart, tEnd, project, turbine, breakdown)

	// for _, avail := range availability {
	// 	log.Printf(">>> %#v \n", avail)
	// }

	/*======== JMR PART ==================*/
	_, _, monthDay := helper.GetDurationInMonth(tStart, tEnd)
	newKey := ""
	monthList := []int{}
	for key, val := range monthDay { /*ubah jika ada key yang hanya 5 huruf ==> 20161 menjadi 201601*/
		newKey = key
		if len(key) < 6 {
			newKey = key[0:4] + "0" + key[4:]
			monthDay.Set(newKey, val)
			monthDay.Unset(key)
		}
		monthList = append(monthList, tk.ToInt(newKey, tk.RoundingAuto))
	}
	pipesJMR := []tk.M{}
	match := []tk.M{}
	match = append(match, tk.M{"dateinfo.monthid": tk.M{"$in": monthList}})
	if len(turbine) != 0 {
		match = append(match, tk.M{"sections.turbine": tk.M{"$in": turbine}})
	}
	projection := tk.M{
		"dateinfo.monthid":  1,
		"sections.turbine":  1,
		"sections.contrgen": 1,
		"sections.boenet":   1,
	}
	pipesJMR = append(pipesJMR, tk.M{"$unwind": "$sections"})
	pipesJMR = append(pipesJMR, tk.M{"$match": tk.M{"$and": match}})
	pipesJMR = append(pipesJMR, tk.M{"$project": projection})

	dataJMR := []tk.M{}
	csrJMR, e := DB().Connection.NewQuery().
		From(new(JMR).TableName()).
		Command("pipe", pipesJMR).
		Cursor(nil)

	if e != nil {
		helper.CreateResult(false, nil, e.Error())
	}
	defer csrJMR.Close()
	e = csrJMR.Fetch(&dataJMR, 0, false)
	if e != nil {
		helper.CreateResult(false, nil, e.Error())
	}
	resultJMR := map[string]float64{}
	for _, val := range dataJMR { /*buat data jmr agar sesuai dengan data scada (resultScada)*/
		month := val.Get("dateinfo").(tk.M).GetInt("monthid") /*isinya 201601, 201602, dst....*/
		months := monthDay.Get(tk.ToString(month)).(tk.M)     /*isinya days(jumlah hari sesuai filter) dan totalInMonth (total hari dalam 1 bulan)*/
		sections := val.Get("sections").(tk.M)                /*isinya turbine, contrgen dan boenet*/
		contrgen := sections.GetFloat64("contrgen") / months.GetFloat64("totalInMonth") * months.GetFloat64("days")
		boenet := sections.GetFloat64("boenet") / months.GetFloat64("totalInMonth") * months.GetFloat64("days")
		resultJMR[sections.GetString("turbine")] += (contrgen - boenet) / 1000.0
	}
	/*======== END OF JMR PART ==================*/

	/*======== Denominator for DATA AVAIL */
	listavaildate := getAvailDateByCondition(project, "ScadaData")
	_availdate := listavaildate.Get(project, tk.M{}).(tk.M).Get("ScadaData", []time.Time{}).([]time.Time)
	if len(_availdate) > 0 {
		if _availdate[0].UTC().After(tStart.UTC()) {
			tStart = _availdate[0]
		}

		if _availdate[1].UTC().Before(tEnd.UTC()) {
			tEnd = _availdate[1]
		}
	}

	totalDataShouldBe := tk.ToFloat64(tEnd.UTC().Sub(tStart.UTC()).Hours()/24, 0, tk.RoundingUp) * 144

	tk.Println(tEnd, tStart, totalDataShouldBe)

	/*======== END OF Denominator for DATA AVAIL */

	LossAnalysisResult := []tk.M{}
	turbineName, e := helper.GetTurbineNameList(p.Project)
	if e != nil {
		helper.CreateResult(false, nil, e.Error())
	}

	scadaavail := ""

	for _, val := range resultScada {
		id := val.GetString("_id")
		var oktime, totalavail float64
		for _, avail := range availability {
			av := avail.Get(id)
			if av != nil {
				oktime = av.(tk.M).GetFloat64("oktime")
				totalavail = av.(tk.M).GetFloat64("totalavail")
				break
			}
		}
		scadaavail = tk.ToString(tk.Div(val.GetFloat64("totalrows"), totalDataShouldBe) * 100)
		scadaavail = " (" + strings.Split(scadaavail, ".")[0] + "." + strings.Split(scadaavail, ".")[1][0:2] + "%)"

		LossAnalysisResult = append(LossAnalysisResult, tk.M{
			"Id": turbineName[val.GetString("_id")] +
				scadaavail,
			"Production":       val.GetFloat64("Production") / 1000,
			"LossEnergy":       val.GetFloat64("LossEnergy") / 1000,
			"MachineDownHours": val.GetFloat64("MachineDownHours"),
			"GridDownHours":    val.GetFloat64("GridDownHours"),
			"OtherDownHours":   val.GetFloat64("OtherDownHours"),
			"EnergyyMD":        val.GetFloat64("MachineDownLoss") / 1000,
			"EnergyyGD":        val.GetFloat64("GridDownLoss") / 1000,
			"EnergyyOD":        val.GetFloat64("OtherDownLoss") / 1000,
			// "ElectricLoss":     val.GetFloat64("ElectricalLosses") / 1000,
			"ElectricLoss":     resultJMR[val.GetString("_id")],
			"PCDeviation":      val.GetFloat64("PCDeviation") / 1000,
			"Others":           val.GetFloat64("OtherDownLoss") / 1000,
			"DownTimeDuration": val.GetFloat64("DownTimeDuration"),
			"OKTime":           oktime / 3600,
			"TotalAvail":       totalavail,
		})
	}

	data := struct {
		Data []tk.M
	}{
		Data: LossAnalysisResult,
	}

	return helper.CreateResult(true, data, "success")
}

func (m *AnalyticLossAnalysisController) GetScadaSummaryChart(k *knot.WebContext) interface{} {
	keys := []string{"Energy Lost Due to Machine Down", "Energy Lost Due to Grid Down", "Others", "Electrical Losses", "PC Deviation"}
	categories := []string{}

	k.Config.OutputType = knot.OutputJson

	var (
		filter     []*dbox.Filter
		pipes      []tk.M
		dataSeries []tk.M
	)

	p := new(PayloadAnalytic)
	e := k.GetPayload(&p)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	// tStart, _ := time.Parse("2006-01-02", p.DateStart.UTC().Format("2006-01-02"))
	// tEnd, _ := time.Parse("2006-01-02 15:04:05", p.DateEnd.UTC().Format("2006-01-02")+" 23:59:59")
	tStart, tEnd, e := helper.GetStartEndDate(k, p.Period, p.DateStart, p.DateEnd)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}
	turbine := p.Turbine
	project := p.Project
	breakDown := p.BreakDown

	filter = append(filter, dbox.Ne("_id", ""))
	filter = append(filter, dbox.Gte("dateinfo.dateid", tStart))
	filter = append(filter, dbox.Lte("dateinfo.dateid", tEnd))

	if project != "" {
		filter = append(filter, dbox.Eq("projectname", project))
	}
	if len(turbine) != 0 {
		filter = append(filter, dbox.In("turbine", turbine...))
	}

	ids := "$" + breakDown

	if breakDown == "dateinfo.monthdesc" {
		pipes = append(pipes, tk.M{"$group": tk.M{"_id": tk.M{}.Set("monthid", "$dateinfo.monthid").Set("monthdesc", "$dateinfo.monthdesc"),
			"LostEnergy":       tk.M{"$sum": "$lostenergy"},
			"MachineDownLoss":  tk.M{"$sum": "$machinedownloss"},
			"GridDownLoss":     tk.M{"$sum": "$griddownloss"},
			"PCDeviation":      tk.M{"$sum": "$pcdeviation"},
			"ElectricalLosses": tk.M{"$sum": "$electricallosses"},
			"OtherDownLoss":    tk.M{"$sum": "$otherdownloss"},
			"Production":       tk.M{"$sum": "$production"}}})
		pipes = append(pipes, tk.M{"$sort": tk.M{"_id.monthid": 1}})
	} else {
		pipes = append(pipes, tk.M{"$group": tk.M{"_id": ids,
			"LostEnergy":       tk.M{"$sum": "$lostenergy"},
			"MachineDownLoss":  tk.M{"$sum": "$machinedownloss"},
			"GridDownLoss":     tk.M{"$sum": "$griddownloss"},
			"PCDeviation":      tk.M{"$sum": "$pcdeviation"},
			"ElectricalLosses": tk.M{"$sum": "$electricallosses"},
			"OtherDownLoss":    tk.M{"$sum": "$otherdownloss"},
			"Production":       tk.M{"$sum": "$production"}}})
		pipes = append(pipes, tk.M{"$sort": tk.M{"_id": 1}})
	}
	csr, e := DB().Connection.NewQuery().
		From(new(ScadaSummaryDaily).TableName()).
		Command("pipe", pipes).
		Where(dbox.And(filter...)).
		Cursor(nil)

	if e != nil {
		helper.CreateResult(false, nil, e.Error())
	}
	defer csr.Close()

	resultScada := []tk.M{}
	e = csr.Fetch(&resultScada, 0, false)
	// add by ams, 2016-10-07
	csr.Close()

	selArr := 1
	for _, key := range keys {
		// series := tk.M{}
		// series.Set("name", key)
		// series.Set("type", "bar")
		// series.Set("axis", "lossBar")
		// series.Set("opacity", 0.6)

		// if key == "Energy Lost Due to Machine Down" ||  key == "Energy Lost Due to Grid Down" {
		// 	series.Set("type", "line")
		// 	series.Set("style", "smooth")
		// 	series.Set("axis", "lossLine")
		// 	series.Set("dashType", "solid")
		// 	series.Set("markers", tk.M{"visible": false})
		// 	series.Set("width", 3)
		// }

		unitMeter := " (MWh)"
		if key == "Energy Lost Due to Machine Down" || key == "Energy Lost Due to Grid Down" {
			unitMeter = " (%)"
		} else if key == "PC Deviation" {
			unitMeter = " (MW)"
		}

		series := tk.M{}
		series.Set("name", key+unitMeter)
		series.Set("type", "line")
		series.Set("style", "smooth")
		series.Set("axis", "lossLine")
		series.Set("dashType", "solid")
		series.Set("markers", tk.M{"visible": false})
		series.Set("width", 3)
		series.Set("color", colorFields[selArr])

		if key == "PC Deviation" || key == "Electrical Losses" {
			series.Set("type", "bar")
			series.Set("axis", "lossBar")
			series.Set("opacity", 0.6)
		}

		var datas []float64
		for _, val := range resultScada {
			LostEnergy := val.GetFloat64("LostEnergy")
			if LostEnergy == 0 {
				LostEnergy = 1
			}
			// tk.Printf("dt %v\n", LostEnergy)

			if key == "Energy Lost Due to Machine Down" {
				datas = append(datas, tk.ToFloat64((val.GetFloat64("MachineDownLoss")/LostEnergy)*100, 2, tk.RoundingAuto))
			} else if key == "Energy Lost Due to Grid Down" {
				datas = append(datas, tk.ToFloat64((val.GetFloat64("GridDownLoss")/LostEnergy)*100, 2, tk.RoundingAuto))
			} else if key == "PC Deviation" {
				datas = append(datas, tk.ToFloat64((val.GetFloat64("PCDeviation")/1000), 2, tk.RoundingAuto))
			} else if key == "Electrical Losses" {
				datas = append(datas, tk.ToFloat64((val.GetFloat64("ElectricalLosses")/1000), 2, tk.RoundingAuto))
			} else if key == "Others" {
				datas = append(datas, tk.ToFloat64((val.GetFloat64("OtherDownLoss")/LostEnergy)*100, 2, tk.RoundingAuto))
			}
		}

		if len(datas) > 0 {
			series.Set("data", datas)
		}
		dataSeries = append(dataSeries, series)
		selArr++
	}

	for _, val := range resultScada {
		id := val.Get("_id")
		monthdescKey := ""
		if breakDown == "dateinfo.monthdesc" {
			_id := val.Get("_id").(tk.M)
			// tk.Printf("_id => %#v\n", _id)
			monthdescKey = tk.ToString(_id["monthdesc"])
		}

		if breakDown == "dateinfo.dateid" {
			dt := id.(time.Time)
			// if index == 0 || dt.Day() == 1 {
			// 	categories = append(categories, tk.ToString(dt.Day())+" "+dt.Month().String()[:3])
			// } else {
			categories = append(categories, tk.ToString(dt.Day()))
			// }
			// categories = append(categories, tk.ToString(dt.Day())+"/"+dt.Month().String()[:3])
		} else if breakDown == "dateinfo.monthdesc" {
			if id != "" {
				categories = append(categories, monthdescKey)
			}
		} else if breakDown == "dateinfo.year" {
			categories = append(categories, tk.ToString(id))
		} else if breakDown == "projectname" {
			categories = append(categories, tk.ToString(id))
		} else if breakDown == "turbine" {
			categories = append(categories, tk.ToString(id))
		}
	}

	result := struct {
		Series     []tk.M
		Categories []string
	}{
		Series:     dataSeries,
		Categories: categories,
	}

	return helper.CreateResult(true, result, "success")
}

func (m *AnalyticLossAnalysisController) GetDowntimeTabDuration(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	p := new(PayloadAnalytic)
	e := k.GetPayload(&p)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	result := tk.M{}

	duration, e := getDownTimeTopLossDuration("duration", p, k)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}
	result.Set("duration", duration)

	return helper.CreateResult(true, result, "success")
}

func (m *AnalyticLossAnalysisController) GetDowntimeTabFreq(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	p := new(PayloadAnalytic)
	e := k.GetPayload(&p)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	result := tk.M{}

	frequency, e := getDownTimeTopFrequency(p, k)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}
	result.Set("frequency", frequency)

	return helper.CreateResult(true, result, "success")
}

func (m *AnalyticLossAnalysisController) GetDowntimeTabLoss(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	p := new(PayloadAnalytic)
	e := k.GetPayload(&p)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	result := tk.M{}

	loss, e := getDownTimeTopLossDuration("loss", p, k)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}
	result.Set("loss", loss)

	return helper.CreateResult(true, result, "success")
}

func (m *AnalyticLossAnalysisController) GetComponentAlarmTab(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	p := new(PayloadAnalytic)
	e := k.GetPayload(&p)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	result := tk.M{}
	// =============== Component Alarm =============
	componentduration, e := getTopComponentAlarm("braketype", "duration", p, k)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}
	result.Set("componentduration", componentduration)

	componentfrequency, e := getTopComponentAlarm("braketype", "frequency", p, k)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}
	result.Set("componentfrequency", componentfrequency)

	componentloss, e := getTopComponentAlarm("braketype", "loss", p, k)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}
	result.Set("componentloss", componentloss)

	// ======= Alarm
	alarmduration, e := getTopComponentAlarm("alertdescription", "duration", p, k)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}
	result.Set("alarmduration", alarmduration)

	alarmfrequency, e := getTopComponentAlarm("alertdescription", "frequency", p, k)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}
	result.Set("alarmfrequency", alarmfrequency)

	alarmloss, e := getTopComponentAlarm("alertdescription", "loss", p, k)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}
	result.Set("alarmloss", alarmloss)

	return helper.CreateResult(true, result, "success")
}

func (m *AnalyticLossAnalysisController) GetFrequencyTab(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	p := new(PayloadAnalytic)
	e := k.GetPayload(&p)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	result := tk.M{}

	catlossfreq, e := getCatLossTopFrequency(p, k)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}
	result.Set("catlossfreq", catlossfreq)

	return helper.CreateResult(true, result, "success")
}

func (m *AnalyticLossAnalysisController) GetLostEnergyTab(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	p := new(PayloadAnalytic)
	e := k.GetPayload(&p)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	result := tk.M{}

	catloss, e := getCatLossTopLossDuration("loss", p, k)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}
	result.Set("catloss", catloss)

	catlossduration, e := getCatLossTopLossDuration("duration", p, k)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}
	result.Set("catlossduration", catlossduration)

	return helper.CreateResult(true, result, "success")
}

func getCatLossTopLossDuration(topType string, p *PayloadAnalytic, k *knot.WebContext) ([]tk.M, error) {
	result := []tk.M{}
	var e error

	if p != nil {
		data := tk.M{}
		lossDurData, e := getLossDuration(topType, p, k)
		if e != nil {
			return result, e
		}
		if len(lossDurData) > 0 {
			data = lossDurData[0]
		}
		machineDown, _ := getMachineDownType()
		tmpResult := []tk.M{}
		for field, title := range machineDown {
			tmpResult = append(tmpResult, tk.M{
				"title":  title,
				"result": data.GetFloat64(field),
			})
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

	return result, e
}

func getCatLossTopFrequency(p *PayloadAnalytic, k *knot.WebContext) ([]tk.M, error) {
	result := []tk.M{}
	var e error

	if p != nil {
		freqData, e := getLossFrequency(p, k)
		if e != nil {
			return result, e
		}
		machineDown, _ := getMachineDownType()
		tmpResult := []tk.M{}
		for field, title := range machineDown {
			tmpResult = append(tmpResult, tk.M{
				"title":  title,
				"result": freqData.GetFloat64(field),
			})
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

	return result, e
}

func getDownTimeTopLossDuration(topType string, p *PayloadAnalytic, k *knot.WebContext) ([]tk.M, error) {
	var result []tk.M
	var e error
	var pipes []tk.M
	match := tk.M{}
	breakdown := "$turbine" /* default grouping */

	if p != nil {
		tStart, tEnd, e := helper.GetStartEndDate(k, p.Period, p.DateStart, p.DateEnd)
		if e != nil {
			return result, e
		}

		if p.Project != "" {
			match.Set("projectname", p.Project)
		}

		if len(p.Turbine) != 0 {
			match.Set("turbine", tk.M{"$in": p.Turbine})
		}
		if p.BreakDown == "$projectname" { /* untuk drilldown loss energy saat series di klik*/
			breakdown = p.BreakDown
			match.Set("dateinfo.monthdesc", p.DateStr) /* drilldown nya per bulan makanya filternya dibedakan */
		} else {
			match.Set("dateinfo.dateid", tk.M{"$gte": tStart, "$lte": tEnd})
		}
		pipes = append(pipes, tk.M{"$match": match})
		if topType == "duration" {
			pipes = append(pipes, tk.M{"$group": tk.M{
				"_id":         breakdown,
				"result":      tk.M{"$sum": "$downtimehours"},
				"machinedown": tk.M{"$sum": "$machinedownhours"},
				"griddown":    tk.M{"$sum": "$griddownhours"},
				"unknown":     tk.M{"$sum": "$otherdowntimehours"},
			}})
		} else if topType == "loss" {
			pipes = append(pipes, tk.M{"$group": tk.M{
				"_id":         breakdown,
				"result":      tk.M{"$sum": "$lostenergy"},
				"machinedown": tk.M{"$sum": "$machinedownloss"},
				"griddown":    tk.M{"$sum": "$griddownloss"},
				"unknown":     tk.M{"$sum": "$otherdownloss"},
			}})
		}

		pipes = append(pipes, tk.M{"$sort": tk.M{"result": -1}})
		pipes = append(pipes, tk.M{"$limit": 10})

		// get the top 10
		csr, e := DB().Connection.NewQuery().
			From(new(ScadaSummaryDaily).TableName()).
			Command("pipe", pipes).
			Cursor(nil)

		if e != nil {
			return result, e
		}

		top10Turbines := []tk.M{}
		e = csr.Fetch(&top10Turbines, 0, false)
		csr.Close()

		if e != nil {
			return result, e
		}

		turbineList := []string{}
		for _, turbine := range top10Turbines {
			turbineList = append(turbineList, turbine.Get("_id").(string)) /* if breakdown by project so it become project */
		}

		downCause, _ := getMachineDownType()

		turbineName, e := helper.GetTurbineNameList(p.Project)
		if e != nil {
			return result, e
		}

		for _, turbine := range turbineList {
			resVal := tk.M{}
			resVal.Set("_id", turbineName[turbine])
			if p.BreakDown == "$projectname" {
				resVal.Set("_id", turbine)
			}
			lossPerTurbine := 0.0
			for _, val := range top10Turbines {
				valTurbine := val.GetString("_id")
				if turbine == valTurbine {
					lossPerTurbine = val.GetFloat64("result")
					if topType == "loss" {
						lossPerTurbine /= 1000 /* jadikan MWh */
					}
					for field, down := range downCause {
						valResultType := val.GetFloat64(field)
						if topType == "loss" {
							valResultType /= 1000 /* jadikan MWh */
						}
						valTitle := strings.Replace(down, " ", "", -69)
						if valResultType >= 0 {
							resVal.Set(valTitle, valResultType) /* MachineDown : 7.6666 */
						}
					}
				}
			}

			resVal.Set("Total", lossPerTurbine)
			result = append(result, resVal)
		}
	}

	return result, e
}

func getDownTimeTopFrequency(p *PayloadAnalytic, k *knot.WebContext) ([]tk.M, error) {
	var result []tk.M
	var e error
	var pipes []tk.M
	match := tk.M{}

	if p != nil {
		tStart, tEnd, e := helper.GetStartEndDate(k, p.Period, p.DateStart, p.DateEnd)
		if e != nil {
			return result, e
		}
		match.Set("_id", tk.M{"$ne": ""})
		match.Set("reduceavailability", true)

		if p.Project != "" {
			match.Set("projectname", p.Project)
		}

		if len(p.Turbine) != 0 {
			match.Set("turbine", tk.M{"$in": p.Turbine})
		}
		match.Set("startdate", tk.M{"$gte": tStart, "$lte": tEnd})
		pipes = append(pipes, tk.M{"$match": match})
		pipes = append(pipes, tk.M{"$group": tk.M{"_id": "$turbine", "result": tk.M{"$sum": 1}}})
		pipes = append(pipes, tk.M{"$sort": tk.M{"result": -1}})
		pipes = append(pipes, tk.M{"$limit": 10})

		// get the top 10
		csr, e := DB().Connection.NewQuery().
			From(new(Alarm).TableName()).
			Command("pipe", pipes).
			Cursor(nil)

		if e != nil {
			return result, e
		}

		top10Turbines := []tk.M{}
		e = csr.Fetch(&top10Turbines, 0, false)

		// add by ams, 2016-10-07
		csr.Close()

		if e != nil {
			return result, e
		}

		// get the downtime
		turbines := []string{}
		turbinesVal := tk.M{}

		for _, turbine := range top10Turbines {
			turbines = append(turbines, turbine.Get("_id").(string))
			turbinesVal.Set(turbine.Get("_id").(string), turbine.GetFloat64("result"))
		}

		match.Set("turbine", tk.M{"$in": turbines})

		downCause := tk.M{}
		downCause.Set("griddown", "Grid Down")
		downCause.Set("machinedown", "Machine Down")
		downCause.Set("unknown", "Unknown")

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
			pipes = append(pipes,
				tk.M{
					"$group": tk.M{"_id": tk.M{"id3": "$turbine", "id4": title},
						"result": tk.M{"$sum": 1},
					},
				},
			)

			pipes = append(pipes, tk.M{"$sort": tk.M{"result": -1}})

			csr, e := DB().Connection.NewQuery().
				From(new(Alarm).TableName()).
				Command("pipe", pipes).
				Cursor(nil)

			if e != nil {
				return result, e
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

		turbineName, e := helper.GetTurbineNameList(p.Project)
		if e != nil {
			return result, e
		}
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

				resY = append(resY, resX)
			}
		}

		for _, turbine := range turbines {
			resVal := tk.M{}
			resVal.Set("_id", turbineName[turbine])

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

	return result, e
}

func getTopComponentAlarm(Id string, topType string, p *PayloadAnalytic, k *knot.WebContext) ([]tk.M, error) {
	var result []tk.M
	var e error
	var pipes []tk.M
	match := tk.M{}
	var dataSeries []tk.M

	if p != nil {
		tStart, tEnd, e := helper.GetStartEndDate(k, p.Period, p.DateStart, p.DateEnd)
		if e != nil {
			return result, e
		}
		match.Set("_id", tk.M{"$ne": ""})
		match.Set("detail.startdate", tk.M{"$gte": tStart, "$lte": tEnd})

		if p.Project != "" {
			match.Set("projectname", p.Project)
		}

		if len(p.Turbine) != 0 {
			match.Set("turbine", tk.M{"$in": p.Turbine})
		}

		if topType == "frequency" {
			match.Set("startdate", tk.M{"$gte": tStart, "$lte": tEnd})
		} else {
			pipes = append(pipes, tk.M{"$unwind": "$detail"})
		}

		pipes = append(pipes, tk.M{"$match": match})
		if topType == "duration" {
			pipes = append(pipes, tk.M{"$group": tk.M{"_id": "$" + Id, "result": tk.M{"$sum": "$detail.duration"}}})
		} else if topType == "frequency" {
			pipes = append(pipes, tk.M{"$group": tk.M{"_id": "$" + Id, "result": tk.M{"$sum": 1}}})
		} else if topType == "loss" {
			pipes = append(pipes, tk.M{"$group": tk.M{"_id": "$" + Id, "result": tk.M{"$sum": "$detail.powerlost"}}})
		}

		pipes = append(pipes, tk.M{"$sort": tk.M{"result": -1}})
		pipes = append(pipes, tk.M{"$limit": 10})

		// get the top 10
		csr, e := DB().Connection.NewQuery().
			From(new(Alarm).TableName()).
			Command("pipe", pipes).
			Cursor(nil)

		if e != nil {
			return result, e
		}

		e = csr.Fetch(&result, 0, false)

		csr.Close()

		for _, val := range result {

			series := tk.M{}
			valueResult := val.GetFloat64("result")
			id := strings.Title(strings.Replace(val.GetString("_id"), "_", " ", -69))

			if topType == "loss" {
				valueResult = tk.Div(valueResult, 1000)
			}

			series.Set("_id", id)
			series.Set("result", valueResult)

			dataSeries = append(dataSeries, series)
		}

		if e != nil {
			return dataSeries, e
		}
	}

	return dataSeries, e
}

func getAvailabilityValue(tStart time.Time, tEnd time.Time, project string, turbine []interface{}, breakDown string) (result []tk.M) {
	pipes := []tk.M{}
	list := []tk.M{}
	match := tk.M{}

	match.Set("dateinfo.dateid", tk.M{"$gte": tStart, "$lte": tEnd})

	if len(turbine) > 0 {
		match.Set("turbine", tk.M{"$in": turbine})
	}

	group := tk.M{
		"power":           tk.M{"$sum": "$powerkw"},
		"machinedowntime": tk.M{"$sum": "$machinedownhours"},
		"griddowntime":    tk.M{"$sum": "$griddownhours"},
		"unknowndowntime": tk.M{"$sum": "$griddownhours"},
		"oktime":          tk.M{"$sum": "$oktime"},
		"powerlost":       tk.M{"$sum": "$lostenergy"},
		"totaltimestamp":  tk.M{"$sum": "$totalrows"},
		// "available":       tk.M{"$sum": "$available"},
		// "minutes":         tk.M{"$sum": "$minutes"},
		"maxdate": tk.M{"$max": "$dateinfo.dateid"},
		"mindate": tk.M{"$min": "$dateinfo.dateid"},
	}

	if project != "" {
		match.Set("projectname", project)
	}

	if breakDown == "Date" {
		group.Set("_id", tk.M{"id1": "$dateinfo.dateid"})
	} else if breakDown == "Month" {
		group.Set("_id", tk.M{"id1": "$dateinfo.monthid", "id2": "$dateinfo.monthdesc"})
	} else if breakDown == "Year" {
		group.Set("_id", tk.M{"id1": "$dateinfo.year"})
	} else if breakDown == "Project" {
		group.Set("_id", tk.M{"id1": "$projectname"})
	} else if breakDown == "Turbine" {
		group.Set("_id", tk.M{"id1": "$turbine"})
	}

	pipes = append(pipes, tk.M{"$match": match})
	pipes = append(pipes, tk.M{"$group": group})
	pipes = append(pipes, tk.M{"$sort": tk.M{"_id.id1": 1}})

	csr, e := DB().Connection.NewQuery().
		From(new(ScadaSummaryDaily).TableName()).
		Command("pipe", pipes).
		Cursor(nil)

	defer csr.Close()

	if e != nil {
		return
	}

	e = csr.Fetch(&list, 0, false)

	var turbineList []TurbineOut
	if project != "" {
		turbineList, _ = helper.GetTurbineList([]interface{}{project})
	} else {
		turbineList, _ = helper.GetTurbineList(nil)
	}

	for _, val := range list {
		var totalTurbine, hourValue, plfDivider float64

		id := val.Get("_id").(tk.M)
		key := id.GetString("id1")

		if breakDown == "Turbine" {
			totalTurbine = 1.0

			for _, v := range turbineList {
				if key == v.Value {
					plfDivider += v.Capacity
				}
			}
		} else if len(turbine) == 0 {
			totalTurbine = float64(len(turbineList))
			for _, v := range turbineList {
				if key == v.Project {
					plfDivider += v.Capacity
				}
			}
		} else {
			totalTurbine = tk.ToFloat64(len(turbine), 1, tk.RoundingAuto)
			for _, vt := range turbine {
				for _, v := range turbineList {
					if vt.(string) == v.Value && key == v.Project {
						plfDivider += v.Capacity
					}
				}
			}
		}

		minDate := val.Get("mindate", time.Time{}).(time.Time)
		maxDate := val.Get("maxdate", time.Time{}).(time.Time)

		// if breakDown == "Date" {
		// 	id1 := id.Get("id1").(time.Time)
		// 	key = id1.Format("20060102_1504050000")
		// 	hourValue = helper.GetHourValue(id1.UTC(), id1.UTC(), minDate.UTC(), maxDate.UTC())
		// } else {

		// hourValue = helper.GetHourValue(tStart.UTC(), tEnd.UTC(), minDate.UTC(), maxDate.UTC())
		// }

		hourValue = maxDate.AddDate(0, 0, 1).UTC().Sub(minDate.UTC()).Hours()

		okTime := val.GetFloat64("oktime") / 3600
		power := val.GetFloat64("power")
		energy := power / 6
		mDownTime := val.GetFloat64("machinedowntime")
		gDownTime := val.GetFloat64("griddowntime")
		uDownTime := val.GetFloat64("unknowndowntime")
		sumTimeStamp := val.GetFloat64("totaltimestamp")
		minutes := val.GetFloat64("minutes") / 60

		// machineAvail, gridAvail, dataAvail, trueAvail, plf := helper.GetAvailAndPLF(totalTurbine, okTime, energy, mDownTime, gDownTime, sumTimeStamp, hourValue, minutes, plfDivider)

		in := tk.M{}.Set("noofturbine", totalTurbine).Set("oktime", okTime).Set("energy", energy/1000).
			Set("totalhour", hourValue).Set("totalcapacity", plfDivider).Set("counttimestamp", sumTimeStamp).
			Set("machinedowntime", mDownTime).Set("griddowntime", gDownTime).Set("otherdowntime", uDownTime)

		calc := helper.CalcAvailabilityAndPLF(in)

		res := tk.M{
			key: tk.M{
				"oktime":          okTime * 3600,
				"power":           power,
				"energy":          energy,
				"machinedowntime": mDownTime,
				"griddowntime":    gDownTime,
				"count":           sumTimeStamp,
				"minutes":         minutes,
				"machineavail":    calc.GetFloat64("machineavailability"),
				"gridavail":       calc.GetFloat64("gridavailability"),
				"dataavail":       calc.GetFloat64("dataavailability"),
				"totalavail":      calc.GetFloat64("totalavailability"),
				"plf":             calc.GetFloat64("plf"),
			},
		}
		result = append(result, res)
	}

	return
}

func (m *AnalyticLossAnalysisController) GetHistogramProduction(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	p := new(PayloadAnalytic)
	e := k.GetPayload(&p)

	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	tStart, _ := time.Parse("2006-01-02", p.DateStart.UTC().Format("2006-01-02"))
	tEnd, _ := time.Parse("2006-01-02 15:04:05", p.DateEnd.UTC().Format("2006-01-02")+" 23:59:59")
	turbine := p.Turbine
	project := p.Project

	match := tk.M{}
	match.Set("dateinfo.dateid", tk.M{}.Set("$lte", tEnd).Set("$gte", tStart))
	match.Set("projectname", project)
	match.Set("avgwindspeed", tk.M{}.Set("$gte", 3).Set("$lt", 25))

	if len(turbine) > 0 {
		match.Set("turbine", tk.M{}.Set("$in", turbine))
	}

	group := tk.M{
		"_id":   "$wsavgforpc",
		"total": tk.M{}.Set("$sum", 1),
	}

	sort := tk.M{
		"_id": 1,
	}

	match.Set("available", 1)

	var pipes []tk.M
	pipes = append(pipes, tk.M{}.Set("$match", match))
	pipes = append(pipes, tk.M{}.Set("$group", group))
	pipes = append(pipes, tk.M{}.Set("$sort", sort))

	csr, e := DB().Connection.NewQuery().
		From(new(ScadaData).TableName()).
		Command("pipe", pipes).
		Cursor(nil)

	defer csr.Close()

	if e != nil {
		return helper.CreateResult(false, nil, "Error query : "+e.Error())
	}

	results := make([]tk.M, 0)
	e = csr.Fetch(&results, 0, false)

	if e != nil {
		return helper.CreateResult(false, nil, "Error facing results : "+e.Error())
	}

	totalData := c.From(&results).Sum(func(x interface{}) interface{} {
		dt := x.(tk.M)
		return dt["total"].(int)
	}).Exec().Result.Sum

	valuewindspeed := tk.M{"3.0": 0}
	valuewindspeed.Set("3.5", 0)

	categorywindspeed := []string{}
	categorywindspeed = append(categorywindspeed, "3 - 3.5")
	categorywindspeed = append(categorywindspeed, "3.5 - 4")
	for i := 4; i <= 24; i++ {
		nextPhase := i + 1
		categorywindspeed = append(categorywindspeed, strconv.Itoa(i)+" - "+strconv.Itoa(nextPhase))
		valuewindspeed.Set(strconv.Itoa(i)+".0", 0)
	}

	for _, x := range results {
		id := tk.RoundingAuto64(x["_id"].(float64), 1)
		total := x["total"].(int)
		value := tk.Div(float64(total), totalData)

		sId := strconv.FormatFloat(id, 'f', 1, 64)

		valuewindspeed.Set(sId, value)
	}

	retvaluews := []float64{}
	retvaluews = append(retvaluews, valuewindspeed.GetFloat64("3.0"))
	retvaluews = append(retvaluews, valuewindspeed.GetFloat64("3.5"))
	for i := 4; i <= 24; i++ {
		retvaluews = append(retvaluews, valuewindspeed.GetFloat64(strconv.Itoa(i)+".0"))
	}

	data := tk.M{
		"categorywindspeed": categorywindspeed,
		"valuewindspeed":    retvaluews,
		"totaldata":         totalData,
	}

	return helper.CreateResult(true, data, "success")
}

type PayloadHistogram struct {
	MaxValue float64
	MinValue float64
	BinValue int
	Filter   PayloadAnalytic
}

func (m *AnalyticLossAnalysisController) GetHistogramData(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	p := new(PayloadHistogram)
	e := k.GetPayload(&p)

	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	tStart, tEnd, e := helper.GetStartEndDate(k, p.Filter.Period, p.Filter.DateStart, p.Filter.DateEnd)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}
	turbine := p.Filter.Turbine
	project := p.Filter.Project

	categorywindspeed := []string{}
	valuewindspeed := []float64{}
	interval := (p.MaxValue - p.MinValue) / float64(p.BinValue)
	interval = tk.ToFloat64(interval, 0, tk.RoundingUp)
	startcategory := p.MinValue
	totalData := 0.0

	match := tk.M{}
	match.Set("dateinfo.dateid", tk.M{}.Set("$lte", tEnd).Set("$gte", tStart))
	if project != "" {
		match.Set("projectname", project)
	}
	if len(turbine) > 0 {
		match.Set("turbine", tk.M{}.Set("$in", turbine))
	}

	group := tk.M{
		"_id":   "",
		"total": tk.M{}.Set("$sum", 1),
	}

	for i := 0; i < (p.BinValue); i++ {
		catformat := "%.1f"
		if startcategory == math.Trunc(startcategory) {
			catformat = "%.0f"
		}
		// categorywindspeed = append(categorywindspeed, fmt.Sprintf("%.0f", startcategory)+" ~ "+fmt.Sprintf("%.0f", (startcategory+interval)))
		categorywindspeed = append(categorywindspeed, fmt.Sprintf(catformat, startcategory))
		match.Set("avgwindspeed", tk.M{}.Set("$lt", (startcategory+(interval*0.5))).Set("$gte", startcategory-(0.5*interval)))
		match.Set("available", 1)

		var pipes []tk.M
		pipes = append(pipes, tk.M{}.Set("$match", match))
		pipes = append(pipes, tk.M{}.Set("$group", group))

		csr, e := DB().Connection.NewQuery().
			From(new(ScadaData).TableName()).
			Command("pipe", pipes).
			Cursor(nil)

		if e != nil {
			return helper.CreateResult(false, nil, "Error query : "+e.Error())
		}

		//timenow := time.Now()
		resultCategory := []tk.M{}
		e = csr.Fetch(&resultCategory, 0, false)
		//duration := time.Now().Sub(timenow).Seconds()
		//tk.Printf("Kondisi utk bin %v = %v\n", i, duration)

		if len(resultCategory) > 0 {
			valuewindspeed = append(valuewindspeed, float64(resultCategory[0]["total"].(int)))
			totalData = totalData + float64(resultCategory[0]["total"].(int))
		} else {
			valuewindspeed = append(valuewindspeed, 0.00)
		}

		csr.Close()

		startcategory = startcategory + interval
	}

	for i := 0; i < len(valuewindspeed); i++ {
		value := float64(int((valuewindspeed[i]/totalData*100)*100)) / 100
		if value < 0 {
			valuewindspeed[i] = 0
		} else {
			valuewindspeed[i] = value
		}
	}
	turbineName, e := helper.GetTurbineNameList(project)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	data := tk.M{
		"categorywindspeed": categorywindspeed,
		"valuewindspeed":    valuewindspeed,
		"totaldata":         totalData,
		"turbinename":       turbineName,
	}

	return helper.CreateResult(true, data, "success")
}

func (m *AnalyticLossAnalysisController) GetProductionHistogramData(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	p := new(PayloadHistogram)
	e := k.GetPayload(&p)

	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	tStart, tEnd, e := helper.GetStartEndDate(k, p.Filter.Period, p.Filter.DateStart, p.Filter.DateEnd)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}
	turbine := p.Filter.Turbine
	project := p.Filter.Project

	categoryproduction := []string{}
	valueproduction := []float64{}
	interval := (p.MaxValue - p.MinValue) / float64(p.BinValue)
	interval = tk.ToFloat64(interval, 0, tk.RoundingUp)
	startcategory := p.MinValue
	totalData := 0.0

	match := tk.M{}
	match.Set("dateinfo.dateid", tk.M{}.Set("$lte", tEnd).Set("$gte", tStart))
	if project != "" {
		match.Set("projectname", project)
	}
	if len(turbine) > 0 {
		match.Set("turbine", tk.M{}.Set("$in", turbine))
	}
	group := tk.M{
		"_id":   "",
		"total": tk.M{}.Set("$sum", 1),
	}

	for i := 0; i < (p.BinValue); i++ {
		// categoryproduction = append(categoryproduction, fmt.Sprintf("%.0f", startcategory)+" ~ "+fmt.Sprintf("%.0f", (startcategory+interval)))
		catformat := "%.1f"
		if startcategory == math.Trunc(startcategory) {
			catformat = "%.0f"
		}
		categoryproduction = append(categoryproduction, fmt.Sprintf(catformat, startcategory))
		match.Set("power", tk.M{}.Set("$lt", (startcategory+(interval*0.5))).Set("$gte", startcategory-(0.5*interval)))
		match.Set("available", 1)

		var pipes []tk.M
		pipes = append(pipes, tk.M{}.Set("$match", match))
		pipes = append(pipes, tk.M{}.Set("$group", group))

		csr, e := DB().Connection.NewQuery().
			From(new(ScadaData).TableName()).
			Command("pipe", pipes).
			Cursor(nil)

		if e != nil {
			csr.Close()
			return helper.CreateResult(false, nil, "Error query : "+e.Error())
		}

		resultCategory := []tk.M{}
		e = csr.Fetch(&resultCategory, 0, false)

		if len(resultCategory) > 0 {
			valueproduction = append(valueproduction, float64(resultCategory[0]["total"].(int)))
			totalData = totalData + float64(resultCategory[0]["total"].(int))
		} else {
			valueproduction = append(valueproduction, 0.00)
		}
		csr.Close()

		startcategory = startcategory + interval
	}

	for i := 0; i < len(valueproduction); i++ {
		value := float64(int((valueproduction[i]/totalData*100)*100)) / 100
		if value < 0 {
			valueproduction[i] = 0
		} else {
			valueproduction[i] = value
		}
	}
	turbineName, e := helper.GetTurbineNameList(project)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	data := tk.M{
		"categoryproduction": categoryproduction,
		"valueproduction":    valueproduction,
		"totaldata":          totalData,
		"turbinename":        turbineName,
	}

	return helper.CreateResult(true, data, "success")
}

func (m *AnalyticLossAnalysisController) GetTempHistogramData(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	type TempHistoPayload struct {
		MaxValue  float64
		MinValue  float64
		BinValue  int
		FieldName string
		Filter    PayloadAnalytic
	}

	p := new(TempHistoPayload)
	e := k.GetPayload(&p)

	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	category := []string{}
	value := []float64{}
	interval := (p.MaxValue - p.MinValue) / float64(p.BinValue)
	interval = tk.ToFloat64(interval, 0, tk.RoundingUp)
	startcategory := p.MinValue
	totalData := 0.0
	tStart, tEnd, e := helper.GetStartEndDate(k, p.Filter.Period, p.Filter.DateStart, p.Filter.DateEnd)

	var wg sync.WaitGroup
	wg.Add(p.BinValue)
	valueMap := tk.M{}

	for i := 0; i < (p.BinValue); i++ {
		catformat := "%.1f"
		if startcategory == math.Trunc(startcategory) {
			catformat = "%.0f"
		}

		category = append(category, fmt.Sprintf(catformat, startcategory))
		if e != nil {
			tk.Println("error on get start end date GetTempHistogramData()", e.Error())
			return helper.CreateResult(false, nil, e.Error())
		}

		// match.Set("isnull", false)
		go func(p *TempHistoPayload, valueMap tk.M, totalData *float64, startcategory, interval float64,
			wg *sync.WaitGroup, tStart, tEnd time.Time) {
			defer wg.Done()
			turbine := p.Filter.Turbine
			project := p.Filter.Project
			match := tk.M{}
			match.Set("dateinfo.dateid", tk.M{}.Set("$lte", tEnd).Set("$gte", tStart))
			if project != "" {
				match.Set("projectname", project)
			}
			if len(turbine) > 0 {
				match.Set("turbine", tk.M{}.Set("$in", turbine))
			}
			match.Set(p.FieldName, tk.M{}.Set("$lt", (startcategory+(interval*0.5))).Set("$gte", startcategory-(0.5*interval)))
			group := tk.M{
				"_id":   "",
				"total": tk.M{}.Set("$sum", 1),
			}
			var pipes []tk.M
			pipes = append(pipes, tk.M{}.Set("$match", match))
			pipes = append(pipes, tk.M{}.Set("$group", group))

			csr, e := DB().Connection.NewQuery().
				From("Scada10MinHFD").
				Command("pipe", pipes).
				Cursor(nil)

			defer csr.Close()
			if e != nil {
				csr.Close()
				tk.Println("error on cursor go func get temperature histogram", e.Error())
				return
			}

			// timenow := time.Now()
			resultCategory := []tk.M{}
			e = csr.Fetch(&resultCategory, 0, false)
			if e != nil {
				csr.Close()
				tk.Println("error on fetch go func get temperature histogram", e.Error())
				return
			}
			// duration := time.Now().Sub(timenow).Seconds()
			// tk.Printf("KOndisi temp %v = %v\n", i, duration)

			catformat := "%.1f"
			if startcategory == math.Trunc(startcategory) {
				catformat = "%.0f"
			}

			m.mux.Lock()
			if len(resultCategory) > 0 {
				valueMap.Set(fmt.Sprintf(catformat, startcategory), float64(resultCategory[0]["total"].(int)))
				*totalData = *totalData + float64(resultCategory[0]["total"].(int))
			} else {
				valueMap.Set(fmt.Sprintf(catformat, startcategory), 0.0)
			}
			m.mux.Unlock()
		}(p, valueMap, &totalData, startcategory, interval, &wg, tStart, tEnd)

		startcategory = startcategory + interval
	}
	wg.Wait()

	for _, val := range category {
		value = append(value, valueMap.GetFloat64(val))
	}

	for i := 0; i < len(value); i++ {
		valuex := float64(int((value[i]/totalData*100)*100)) / 100
		if valuex < 0 {
			value[i] = 0
		} else {
			value[i] = valuex
		}
	}
	turbineName, e := helper.GetTurbineNameList(p.Filter.Project)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	data := tk.M{
		"category":    category,
		"value":       value,
		"totaldata":   totalData,
		"turbinename": turbineName,
	}

	return helper.CreateResult(true, data, "success")
}

func (m *AnalyticLossAnalysisController) GetWarning(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	p := new(PayloadAnalytic)
	e := k.GetPayload(&p)

	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	tStart, tEnd, e := helper.GetStartEndDate(k, p.Period, p.DateStart, p.DateEnd)
	/*tStart, _ := time.Parse("2006-01-02", p.DateStart.UTC().Format("2006-01-02"))
	tEnd, _ := time.Parse("2006-01-02 15:04:05", p.DateEnd.UTC().Format("2006-01-02")+" 23:59:59")*/
	turbine := p.Turbine
	project := p.Project

	match := tk.M{}
	match.Set("dateinfostart.dateid", tk.M{}.Set("$lte", tEnd).Set("$gte", tStart))
	match.Set("projectname", project)
	if len(turbine) > 0 {
		match.Set("turbine", tk.M{}.Set("$in", turbine))
	}

	group := tk.M{
		"_id":   tk.M{"desc": "$alarmdescription", "turbine": "$turbine"},
		"count": tk.M{"$sum": 1},
	}

	var pipes []tk.M
	pipes = append(pipes, tk.M{}.Set("$match", match))
	pipes = append(pipes, tk.M{}.Set("$group", group))
	pipes = append(pipes, tk.M{}.Set("$sort", tk.M{
		"_id": 1,
	}))

	csr, e := DB().Connection.NewQuery().
		From(new(EventAlarm).TableName()).
		Command("pipe", pipes).
		Cursor(nil)

	defer csr.Close()

	if e != nil {
		return helper.CreateResult(false, nil, "Error query : "+e.Error())
	}

	results := make([]tk.M, 0)
	e = csr.Fetch(&results, 0, false)

	// log.Printf("results: %v \n", len(results))

	/*for _, v := range results {
		log.Printf("results: %#v \n", v)
	}*/

	if e != nil {
		return helper.CreateResult(false, nil, "Error facing results : "+e.Error())
	}

	turbines := []TurbineOut{}
	if len(turbine) == 0 {
		turbines, _ = helper.GetTurbineList([]interface{}{project})
	} else {
		for _, v := range turbine {
			turbines = append(turbines, TurbineOut{
				Project: "",
				Turbine: v.(string),
			})
		}
	}

	// sort.Strings(turbines)

	descs := []string{}
	mapRes := map[string][]tk.M{}
	for _, v := range results {
		id := v.Get("_id").(tk.M)
		desc := id.GetString("desc")
		turbine := id.GetString("turbine")
		count := v.GetInt("count")

		if len(mapRes[desc]) == 0 {
			defHeader := []tk.M{}
			for _, v := range turbines {
				defHeader = append(defHeader, tk.M{"turbine": v, "count": 0})
			}
			mapRes[desc] = defHeader
		}

		var tmp []tk.M
		tmp = mapRes[desc]

		for _, t := range tmp {
			if t.GetString("turbine") == turbine {
				t.Set("count", t.GetInt("count")+count)
				break
			}
		}
		mapRes[desc] = tmp

		found := false

		for _, loop := range descs {
			if loop == desc {
				found = true
				break
			}
		}

		if !found {
			descs = append(descs, desc)
		}

	}

	sort.Strings(descs)

	res := []tk.M{}
	for _, v := range descs {
		total := 0
		for _, x := range mapRes[v] {
			total += x.GetInt("count")
		}

		res = append(res, tk.M{"desc": v, "turbines": mapRes[v], "total": total})
	}

	data := struct {
		Data []tk.M
	}{
		Data: res,
	}

	return helper.CreateResult(true, data, "success")
}

func (m *AnalyticLossAnalysisController) GetAvailDateAll(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	return helper.CreateResult(true, getLastAvailDateAll(), "success")
}

func (m *AnalyticLossAnalysisController) GetAvailDate_DRAFT(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson
	type AvailDatePayload struct {
		Project string
	}

	p := AvailDatePayload{}
	e := k.GetPayload(&p)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}
	datePeriod := k.Session("availdate", "").(map[string]*Availdatedata)
	if p.Project == "" {
		p.Project = "All"
	}
	result := tk.M{}
	result["availabledate"] = datePeriod[p.Project]
	lastDateData, _ := time.Parse("2006-01-02 15:04", "2016-10-31 23:59")
	lastDateData = datePeriod[p.Project].ScadaData[1].UTC()
	result["lastdate"] = lastDateData

	return helper.CreateResult(true, result, "success")
}

func (m *AnalyticLossAnalysisController) GetTempTags(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	tempTags, _ := helper.GetTemperatureList()

	return helper.CreateResult(true, tempTags, "success")
}

func (m *AnalyticLossAnalysisController) GetMaxValTempTags(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	type TempMaxValTempPayload struct {
		FieldList []string
		Filter    PayloadAnalytic
	}

	p := new(TempMaxValTempPayload)
	e := k.GetPayload(&p)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	tStart, tEnd, e := helper.GetStartEndDate(k, p.Filter.Period, p.Filter.DateStart, p.Filter.DateEnd)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}
	turbine := p.Filter.Turbine
	project := p.Filter.Project
	match := tk.M{}
	match.Set("dateinfo.dateid", tk.M{}.Set("$lte", tEnd).Set("$gte", tStart))
	if project != "" {
		match.Set("projectname", project)
	}
	if len(turbine) > 0 {
		match.Set("turbine", tk.M{}.Set("$in", turbine))
	}

	group := tk.M{
		"_id":   "",
		"total": tk.M{}.Set("$sum", 1),
	}

	for _, field := range p.FieldList {
		group.Set(field, tk.M{}.Set("$max", "$"+field))
	}

	var pipes []tk.M
	pipes = append(pipes, tk.M{}.Set("$match", match))
	pipes = append(pipes, tk.M{}.Set("$group", group))

	//timenow := time.Now()
	csr, e := DB().Connection.NewQuery().
		From("Scada10MinHFD").
		Command("pipe", pipes).
		Cursor(nil)

	datamax := tk.M{}
	_ = csr.Fetch(&datamax, 1, false)
	//duration := time.Now().Sub(timenow).Seconds()
	//tk.Printf("Kondisi 1 = %v\n", duration)

	//timenow = time.Now()
	for key, _ := range datamax {
		if key == "_id" || key == "total" {
			continue
		}
		datamax.Set(key, datamax.GetFloat64(key)+10)
	}
	//duration = time.Now().Sub(timenow).Seconds()
	//tk.Printf("Kondisi 2 = %v\n", duration)
	csr.Close()

	return helper.CreateResult(true, datamax, "success")
}

type PayloadEventAnalysis struct {
	Period           string
	Project          string
	Turbine          []interface{}
	DateStart        time.Time
	DateEnd          time.Time
	BreakDown        string
	AdditionalFilter tk.M
	RealDesc         tk.M
}

func (m *AnalyticLossAnalysisController) GetEventAnalysisTab(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	p := new(PayloadEventAnalysis)
	e := k.GetPayload(&p)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	result, resultpct, realDesc, e := getEventAnalysis(p.BreakDown, p.AdditionalFilter, p.RealDesc, p, k)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}
	data := tk.M{
		"data":           result,
		"datapercentage": resultpct,
		"realdesc":       realDesc,
	}

	return helper.CreateResult(true, data, "success")
}

func getEventAnalysis(breakDown string, addFilter, realDesc tk.M, p *PayloadEventAnalysis, k *knot.WebContext) ([]tk.M, []tk.M, tk.M, error) {
	var e error
	dfilter := []*dbox.Filter{}
	var dataSeries []tk.M
	var dataSeriesPct []tk.M

	if p != nil {
		tStart, tEnd, e := helper.GetStartEndDate(k, p.Period, p.DateStart, p.DateEnd)
		if e != nil {
			return dataSeries, dataSeriesPct, realDesc, e
		}
		if p.Project != "" {
			dfilter = append(dfilter, dbox.Eq("projectname", p.Project))
		}
		orFilter := dbox.Or(dbox.And(dbox.Gte("timestart", tStart), dbox.Lte("timestart", tEnd)),
			dbox.And(dbox.Gte("timeend", tStart), dbox.Lte("timeend", tEnd)),
			dbox.And(dbox.Lte("timestart", tStart), dbox.Gte("timeend", tEnd)))
		if len(p.Turbine) > 0 {
			dfilter = append(dfilter, dbox.In("turbine", p.Turbine...))
		}
		dfilter = append(dfilter, dbox.Eq("isdeleted", false))
		dfilter = append(dfilter, dbox.Eq("reduceavailability", true))
		dfilter = append(dfilter, orFilter)
		/*
			filter tambahan untuk tiap breakdown
			jika breakdown = alarmdesc maka dapat filter tambahan detailgroup
			jika breakdown = turbine maka dapat filter tambahan alarmdesc dan detailgroup
		*/
		if len(addFilter) > 0 {
			for key, val := range addFilter {
				if breakDown == "turbine" && key == "alarmdesc" { /* jika breakdown turbine menggunakan alarmdesc sesuai DB */
					dfilter = append(dfilter, dbox.Eq(key, realDesc.GetString(tk.ToString(val))))
				} else {
					dfilter = append(dfilter, dbox.Eq(key, val))
				}
			}
		}

		csr, e := DBRealtime().NewQuery().
			Select("turbine", "timestart", "timeend", "duration", "alarmdesc", "detailgroup").
			From("AlarmHFD").
			Where(dbox.And(dfilter...)).
			Cursor(nil)

		if e != nil {
			return dataSeries, dataSeriesPct, realDesc, e
		}

		_data := tk.M{}
		dataPerGroup := map[string]float64{} /* data per breakdown */
		groupName := ""
		timestart := time.Time{}
		timeend := time.Time{}
		duration := 0.0
		timeIndia := getTimeNow()
		totalHours := 0.0
		for {
			_data = tk.M{}
			e = csr.Fetch(&_data, 1, false)
			if e != nil {
				e = nil
				break
			}
			timestart = _data.Get("timestart", time.Time{}).(time.Time).UTC()
			timeend = _data.Get("timeend", time.Time{}).(time.Time).UTC()
			duration = timeend.Sub(timestart).Hours()
			groupName = _data.GetString(breakDown)

			if !timestart.Before(tStart) && timeend.After(tEnd) { /* jika timeend melebihi tEnd filter */
				duration = tEnd.Sub(timestart).Hours() /* tEnd filter - timestart db */
			} else if !timeend.After(tEnd) && timestart.Before(tStart) { /* jika timestart sebelum tStart */
				duration = timeend.Sub(tStart).Hours() /* timeend db - tStart filter */
			} else if !timestart.Before(tStart) && timeend.IsZero() { /* jika alarm belum selesai */
				if timeIndia.After(tEnd) { /* jika time now India lebih besar dari tEnd filter maka gunakan tEnd filter */
					duration = tEnd.Sub(timestart).Hours()
				} else {
					duration = timeIndia.Sub(timestart).Hours()
				}
			}
			dataPerGroup[groupName] += duration
			totalHours += duration
		}

		csr.Close()

		if breakDown == "alarmdesc" {
			realDesc = tk.M{} /* untuk alarm desc sesuai value di DB */
		}
		turbineName, e := helper.GetTurbineNameList(p.Project)
		if e != nil {
			return dataSeries, dataSeriesPct, realDesc, e
		}
		for key, val := range dataPerGroup {
			series := tk.M{}
			seriesPct := tk.M{}
			id := strings.Title(strings.Replace(key, "_", " ", -69)) /* underscore diganti spasi dan huruf awal besar semua */
			if breakDown == "turbine" {
				id = turbineName[key]
			}
			if breakDown == "alarmdesc" {
				realDesc.Set(id, key)
			}

			series.Set("_id", id)
			series.Set("result", val)
			seriesPct.Set("_id", id)
			seriesPct.Set("result", tk.Div(val, totalHours)*100)

			dataSeries = append(dataSeries, series)
			dataSeriesPct = append(dataSeriesPct, seriesPct)
		}
	}

	return dataSeries, dataSeriesPct, realDesc, e
}
