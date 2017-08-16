package controller

import (
	. "eaciit/wfdemo-git/library/core"
	. "eaciit/wfdemo-git/library/helper"
	. "eaciit/wfdemo-git/library/models"
	"eaciit/wfdemo-git/web/helper"

	"time"

	"github.com/eaciit/dbox"
	"github.com/eaciit/knot/knot.v1"
	tk "github.com/eaciit/toolkit"
)

var (
	from time.Time
	to   time.Time
)

type DataAvailabilityController struct {
	App
}

type FalseContainer struct {
	Order int
	Start time.Time
	End   time.Time
}

// DataAvailabilityController
func CreateDataAvailabilityController() *DataAvailabilityController {
	var controller = new(DataAvailabilityController)
	return controller
}

// GetDataAvailability
func (m *DataAvailabilityController) GetDataAvailability(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	var (
		result []tk.M
		months []string
	)

	from = time.Now()
	to = time.Now()

	p := new(PayloadAnalytic)
	e := k.GetPayload(&p)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	turbine := p.Turbine
	project := p.Project

	if p.BreakDown == "daily" {
		result = append(result, getAvailDaily(project, turbine, p.Period))

		timeParse, e := time.Parse("January 2006", p.Period)
		if e != nil {
			return nil
		}
		timeParse = timeParse.AddDate(0, 1, -1)
		for idx := 1; idx <= timeParse.Day(); idx++ {
			months = append(months, tk.ToString(idx))
		}

	} else {
		result = append(result, getAvailCollection(project, turbine, "SCADA_DATA_OEM"))
		result = append(result, getAvailCollection(project, turbine, "SCADA_DATA_HFD"))
		result = append(result, getAvailCollection(project, turbine, "MET_TOWER"))

		for {
			months = append(months, from.Format("January 2006"))
			if from.Format("0601") == to.Format("0601") {
				break
			}
			from = GetNormalAddDateMonth(from.UTC(), 1)
		}
	}

	data := struct {
		Data  []tk.M
		Month []string
	}{
		Data:  result,
		Month: months,
	}

	return helper.CreateResult(true, data, "success")
}

func (m *DataAvailabilityController) GetCurrentDataAvailability(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	var (
		filter []*dbox.Filter
		pipes  []tk.M
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

	filter = append(filter, dbox.Ne("_id", ""))
	filter = append(filter, dbox.Gte("dateinfo.dateid", tStart))
	filter = append(filter, dbox.Lte("dateinfo.dateid", tEnd))

	if project != "" {
		filter = append(filter, dbox.Eq("projectname", project))
	}

	if len(turbine) != 0 {
		filter = append(filter, dbox.In("turbine", turbine...))
	}

	pipes = append(pipes, tk.M{"$group": tk.M{"_id": "$turbine",
		"totalrows": tk.M{"$sum": "$totalrows"}}})

	csr, e := DB().Connection.NewQuery().
		From(new(ScadaSummaryDaily).TableName()).
		Command("pipe", pipes).
		Where(dbox.And(filter...)).
		Cursor(nil)

	if e != nil {
		helper.CreateResult(false, 0, e.Error())
	}
	defer csr.Close()

	resultScada := []tk.M{}
	e = csr.Fetch(&resultScada, 0, false)
	if e != nil {
		helper.CreateResult(false, 0, e.Error())
	}

	rturbines := tEnd.UTC().Sub(tStart.UTC()).Hours() * 6
	iturbine, totalrows := float64(0), float64(0)
	for _, val := range resultScada {
		iturbine += 1
		totalrows += val.GetFloat64("totalrows")
	}

	return helper.CreateResult(true, tk.Div(totalrows, rturbines*iturbine), "success")
}

func setOpacity(scadaAvail float64) string {
	switch {
	case scadaAvail >= 0.0 && scadaAvail <= 0.1:
		return "100%"
	case scadaAvail > 0.1 && scadaAvail <= 0.2:
		return "90%"
	case scadaAvail > 0.2 && scadaAvail <= 0.3:
		return "80%"
	case scadaAvail > 0.3 && scadaAvail <= 0.4:
		return "70%"
	case scadaAvail > 0.4 && scadaAvail <= 0.5:
		return "60%"
	case scadaAvail > 0.5 && scadaAvail <= 0.6:
		return "60%"
	case scadaAvail > 0.6 && scadaAvail <= 0.7:
		return "70%"
	case scadaAvail > 0.7 && scadaAvail <= 0.8:
		return "80%"
	case scadaAvail > 0.8 && scadaAvail <= 0.9:
		return "90%"
	case scadaAvail > 0.9 && scadaAvail <= 1.0:
		return "100%"
	}

	return "100%"
}

func getAvailDaily(project string, turbines []interface{}, monthdesc string) tk.M {
	pipes := []tk.M{}
	query := []tk.M{}
	dailyData := []tk.M{}
	if project != "" {
		query = append(query, tk.M{"projectname": project})
	}

	if len(turbines) > 0 {
		query = append(query, tk.M{"turbine": tk.M{"$in": turbines}})
	}
	query = append(query, tk.M{"dateinfo.monthdesc": monthdesc})

	pipes = append(pipes, tk.M{"$match": tk.M{"$and": query}})
	pipes = append(pipes, tk.M{"$group": tk.M{
		"_id":        "$dateinfo.dateid",
		"scadaavail": tk.M{"$avg": "$scadaavail"},
	}})

	timeParse, e := time.Parse("January 2006", monthdesc)
	if e != nil {
		return nil
	}
	timeParse = timeParse.AddDate(0, 1, -1)

	pipes = append(pipes, tk.M{"$sort": tk.M{"_id": 1}})

	csr, e := DB().Connection.NewQuery().
		From(new(ScadaSummaryDaily).TableName()).
		Command("pipe", pipes).
		Cursor(nil)

	if e != nil {
		return nil
	}
	defer csr.Close()

	e = csr.Fetch(&dailyData, 0, false)
	if e != nil {
		return nil
	}

	result := tk.M{}

	if len(dailyData) > 0 {
		totalDay := timeParse.Day()
		dataPerDay := tk.M{}
		timeConv := time.Time{}
		for _, val := range dailyData {
			timeConv = val.Get("_id", time.Time{}).(time.Time)
			dataPerDay.Set(tk.ToString(timeConv.Day()), val.GetFloat64("scadaavail"))
		}
		datas := []tk.M{}
		percentage := 0.0
		kelas := "progress-bar progress-bar-success"
		for idx := 1; idx <= totalDay; idx++ {
			percentage = 1.0 / tk.ToFloat64(totalDay, 6, tk.RoundingAuto)
			if dataPerDay.Has(tk.ToString(idx)) {
				if dataPerDay.GetFloat64(tk.ToString(idx)) < 0.5 {
					kelas = "progress-bar progress-bar-red"
				} else {
					kelas = "progress-bar progress-bar-success"
				}
				datas = append(datas, tk.M{
					"tooltip":  "Day " + tk.ToString(idx),
					"class":    kelas,
					"value":    tk.ToString(percentage) + "%",
					"floatval": percentage,
					"opacity":  setOpacity(dataPerDay.GetFloat64(tk.ToString(idx))),
				})
			} else {
				datas = append(datas, tk.M{
					"tooltip":  "Day " + tk.ToString(idx),
					"class":    "progress-bar progress-bar-red",
					"value":    tk.ToString(percentage) + "%",
					"floatval": percentage,
					"opacity":  "100%",
				})
			}
		}
		result = tk.M{"Category": "Data Availability", "Turbine": []tk.M{}, "Data": datas}

		dailyTurbine := []tk.M{}
		pipes = []tk.M{}
		pipes = append(pipes, tk.M{"$match": tk.M{"$and": query}})
		pipes = append(pipes, tk.M{"$group": tk.M{
			"_id":        tk.M{"tanggal": "$dateinfo.dateid", "turbine": "$turbine"},
			"scadaavail": tk.M{"$avg": "$scadaavail"},
		}})

		pipes = append(pipes, tk.M{"$sort": tk.M{"_id.turbine": 1, "_id.tanggal": 1}})

		csrTurbine, e := DB().Connection.NewQuery().
			From(new(ScadaSummaryDaily).TableName()).
			Command("pipe", pipes).
			Cursor(nil)

		if e != nil {
			return nil
		}
		defer csrTurbine.Close()

		e = csrTurbine.Fetch(&dailyTurbine, 0, false)
		if e != nil {
			return nil
		}

		if len(dailyTurbine) > 0 {
			dataPerTurbine := tk.M{}
			timeConv = time.Time{}
			turbineList := []string{}
			lastTurbine := ""
			ids := tk.M{}
			for _, val := range dailyTurbine {
				ids, _ = tk.ToM(val.Get("_id"))
				if lastTurbine != ids.GetString("turbine") {
					lastTurbine = ids.GetString("turbine")
					turbineList = append(turbineList, lastTurbine)
				}
				if tk.TypeName(ids.Get("tanggal")) == "string" {
					timeConv, e = time.Parse("2006-01-02T15:04:05Z07:00", tk.ToString(ids.GetString("tanggal")))
					if e != nil {
						tk.Println(e.Error())
					}
				} else {
					timeConv = ids.Get("tanggal", time.Time{}).(time.Time)
				}
				dataPerTurbine.Set(tk.ToString(timeConv.Day())+"_"+ids.GetString("turbine"), val.GetFloat64("scadaavail"))
			}
			turbineDetails := []tk.M{}
			turbineDatas := []tk.M{}
			turbineItem := tk.M{}
			percentage := 0.0
			kelas := "progress-bar progress-bar-success"
			turbineName, e := helper.GetTurbineNameList(project)
			if e != nil {
				tk.Println(e.Error())
			}
			_turbine := ""
			lastTurbine = ""

			for _, turbine := range turbineList {
				_turbine = turbineName[turbine]
				if lastTurbine != turbine {
					lastTurbine = turbine
					turbineDetails = []tk.M{}
				}
				for idx := 1; idx <= totalDay; idx++ {
					percentage = 1.0 / tk.ToFloat64(totalDay, 6, tk.RoundingAuto)
					if dataPerTurbine.Has(tk.ToString(idx) + "_" + turbine) {
						if dataPerTurbine.GetFloat64(tk.ToString(idx)+"_"+turbine) < 0.5 {
							kelas = "progress-bar progress-bar-red"
						} else {
							kelas = "progress-bar progress-bar-success"
						}
						turbineDetails = append(turbineDetails, tk.M{
							"tooltip":  "Day " + tk.ToString(idx),
							"class":    kelas,
							"value":    tk.ToString(percentage) + "%",
							"floatval": percentage,
							"opacity":  setOpacity(dataPerTurbine.GetFloat64(tk.ToString(idx) + "_" + turbine)),
						})
					} else {
						turbineDetails = append(turbineDetails, tk.M{
							"tooltip":  "Day " + tk.ToString(idx),
							"class":    "progress-bar progress-bar-red",
							"value":    tk.ToString(percentage) + "%",
							"floatval": percentage,
							"opacity":  "100%",
						})
					}
				}
				turbineItem = tk.M{
					"TurbineName": _turbine,
					"details":     turbineDetails,
				}
				turbineDatas = append(turbineDatas, turbineItem)
			}
			result.Set("Turbine", turbineDatas)
		}
	}

	return result
}

func getAvailCollection(project string, turbines []interface{}, collType string) tk.M {
	var (
		pipes          []tk.M
		list           []tk.M
		falseContainer []FalseContainer
	)
	group := tk.M{
		"_id": tk.M{
			"name":    "$name",
			"project": "$details.projectname",
			"turbine": "$details.turbine",
		},
		"periodTo":   tk.M{"$max": "$periodto"},
		"periodFrom": tk.M{"$min": "$periodfrom"},
		"list": tk.M{
			"$push": tk.M{
				"id":       "$details.id",
				"start":    "$details.start",
				"end":      "$details.end",
				"duration": "$details.duration",
				"isavail":  "$details.isavail",
			},
		},
	}

	projection := tk.M{
		"name":       "$_id.name",
		"project":    "$_id.project",
		"turbine":    "$_id.turbine",
		"periodTo":   1,
		"periodFrom": 1,
		"list":       1,
	}

	pipes = append(pipes, tk.M{"$match": tk.M{"type": tk.M{"$eq": collType}}})
	pipes = append(pipes, tk.M{"$unwind": "$details"})
	pipes = append(pipes, tk.M{"$group": group})
	pipes = append(pipes, tk.M{"$project": projection})

	match := tk.M{}

	if project != "" {
		match.Set("project", project)
	}

	if len(turbines) > 0 {
		match.Set("turbine", tk.M{"$in": turbines})
	}

	if match.Get("turbine") != nil || match.Get("project") != nil {
		pipes = append(pipes, tk.M{"$match": match})
	}

	pipes = append(pipes, tk.M{"$sort": tk.M{"turbine": 1}})

	csr, e := DB().Connection.NewQuery().
		From(new(DataAvailability).TableName()).
		Command("pipe", pipes).
		Cursor(nil)

	if e != nil {
		return nil
	}

	e = csr.Fetch(&list, 0, false)

	defer csr.Close()

	res := []tk.M{}
	name := ""

	if len(list) > 0 {
		totalPercent := 0.0
		diffPercent := 0.0
		datas := []tk.M{}
		turbineName := map[string]string{}
		latestProject := ""

		for _, dt := range list {
			p := dt.GetString("project")
			if latestProject != p {
				turbineName, e = helper.GetTurbineNameList(p)
				if e != nil {
					return nil
				}
			}
			t := turbineName[dt.GetString("turbine")]
			_ = p
			pTo := dt.Get("periodTo").(time.Time)
			pFrom := dt.Get("periodFrom").(time.Time)

			from = pFrom.UTC()
			to = pTo.UTC()

			name = dt.GetString("name")
			availList := dt.Get("list").([]interface{})

			turbineDetails := []tk.M{}

			// log.Printf(">> %v | %v | %v | %v | %v | %v \n", p, t, pFrom.String(), pTo.String(), totalDurationDays, name)

			// set availabiility data based on index ordering in collection
			// log.Printf(">>>> turbine: %v \n", t)
			for index := 1; index <= len(availList); index++ {
			breakAvail:
				for _, av := range availList {
					avail := av.(tk.M)
					if index == avail.GetInt("id") {
						start := avail.Get("start").(time.Time).UTC()
						end := avail.Get("end").(time.Time).UTC()
						duration := avail.GetFloat64("duration")
						isAvail := avail.Get("isavail").(bool)

						if !isAvail {
							falseContainer = setFalseContainer(start, end, falseContainer)
							// log.Printf(">> !avail: %v | %v | %v \n", start.String(), end.String(), duration)
							// for _, fc := range falseContainer {
							// 	log.Printf(">> falsecontainer: %v | %v | %v \n", fc.Order, fc.Start.String(), fc.End.String())
							// }
						}

						turbineDetails = append(turbineDetails, setDataColumn(start, end, isAvail, duration))

						// log.Printf(">>>> %v | %v | %v \n", start.Format("2 Jan 2006")+" until "+end.Format("2 Jan 2006"), class, tk.ToString(percentage)+"%")

						break breakAvail
					}
				}
			}
			totalPercent = 0.0
			for idx, val := range turbineDetails {
				totalPercent += val.GetFloat64("floatval")
				if idx == len(turbineDetails)-1 {
					diffPercent = totalPercent - 100.0
					turbineDetails[idx].Set("value", tk.ToString(val.GetFloat64("floatval")-diffPercent)+"%")
				}
			}

			turbine := tk.M{"TurbineName": t}
			turbine.Set("details", turbineDetails)

			res = append(res, turbine)
		}

		// for _, fc := range falseContainer {
		// 	log.Printf("%v | %v | %v \n", fc.Order, fc.Start.String(), fc.End.String())
		// }

		// set turbine parent availabililty
		var before time.Time
		for idx, fc := range falseContainer {
			if idx == 0 {
				if fc.Start.Sub(from.UTC()).Seconds() > 0 {
					datas = append(datas, setDataColumn(from, fc.Start, true, fc.Start.Sub(from.UTC()).Hours()/24))
				}
				datas = append(datas, setDataColumn(fc.Start, fc.End, false, fc.End.Sub(fc.Start.UTC()).Hours()/24))
				before = fc.End
			} else {
				if fc.Start.Sub(before.UTC()).Seconds() > 0 {
					datas = append(datas, setDataColumn(before, fc.Start, true, fc.Start.Sub(before.UTC()).Hours()/24))
				}
				datas = append(datas, setDataColumn(fc.Start, fc.End, false, fc.End.Sub(fc.Start.UTC()).Hours()/24))
				before = fc.End
			}
		}

		if collType == "MET_TOWER" && project != "Tejuva" {
			datas = []tk.M{}
			datas = append(datas, tk.M{
				"tooltip":  from.Format("2 Jan 2006") + " until " + to.Format("2 Jan 2006"),
				"class":    "progress-bar progress-bar-red",
				"value":    "100%",
				"floatval": 100.0,
			})
		}

		if collType != "MET_TOWER" {
			totalPercent = 0.0
			for idx, val := range datas {
				totalPercent += val.GetFloat64("floatval")
				if idx == len(datas)-1 {
					diffPercent = totalPercent - 100.0
					datas[idx].Set("value", tk.ToString(val.GetFloat64("floatval")-diffPercent)+"%")
				}
			}
			return tk.M{"Category": name, "Turbine": res, "Data": datas}
		} else {
			return tk.M{"Category": name, "Turbine": []tk.M{}, "Data": datas}
		}

	}

	return nil
}

func setDataColumn(start time.Time, end time.Time, avail bool, durationInDay float64) tk.M {
	res := tk.M{}
	totalDurationDays := to.UTC().Sub(from.UTC()).Hours() / 24
	class := "progress-bar progress-bar-success"

	if !avail {
		class = "progress-bar progress-bar-red"
	}

	percentage := durationInDay / totalDurationDays * 100

	res = tk.M{
		"tooltip":  start.Format("2 Jan 2006") + " until " + end.Format("2 Jan 2006"),
		"class":    class,
		"value":    tk.ToString(percentage) + "%",
		"floatval": tk.ToFloat64(percentage, 6, tk.RoundingAuto),
	}

	return res
}

func setFalseContainer(start time.Time, end time.Time, falseContainer []FalseContainer) []FalseContainer {
	newFalseContainer := []FalseContainer{}
	if len(falseContainer) == 0 {
		newFalseContainer = append(newFalseContainer, FalseContainer{1, start.UTC(), end.UTC()})
		// log.Printf("new: %v \n", newFalseContainer[0])
	} else {
		// current := FalseContainer{}

		startInt := tk.ToInt(start.Format("20060102150504"), tk.RoundingAuto)
		endInt := tk.ToInt(end.Format("20060102150504"), tk.RoundingAuto)

		// found := false

		idx := 0
		insertedMap := map[string]bool{}
		for _, ct := range falseContainer {
			var ctStartInt, ctEndInt int
			idx++

			ctStartInt = tk.ToInt(ct.Start.Format("20060102150504"), tk.RoundingAuto)
			ctEndInt = tk.ToInt(ct.End.Format("20060102150504"), tk.RoundingAuto)
			// if !found {

			// log.Printf("%v - %v | %v - %v \n", startInt, endInt, ctStartInt, ctEndInt)

			if startInt >= ctStartInt && endInt <= ctEndInt { // inside all
				//log.Println(1)
				newFalseContainer = append(newFalseContainer, FalseContainer{idx, ct.Start, ct.End})
				insertedMap[tk.ToString(ctStartInt)+"-"+tk.ToString(ctEndInt)] = true

				xCount := idx
				for _, x := range falseContainer[idx:] {
					if !insertedMap[x.Start.Format("20060102150504")+"-"+x.End.Format("20060102150504")] {
						xCount++
						newFalseContainer = append(newFalseContainer, FalseContainer{xCount, x.Start, x.End})
					}
				}
				break

			} else if startInt < ctStartInt && endInt > ctEndInt { // start outside, end outside
				//log.Println(2)
				newFalseContainer = append(newFalseContainer, FalseContainer{idx, start, end})
				insertedMap[tk.ToString(startInt)+"-"+tk.ToString(endInt)] = true

				xCount := idx
				for _, x := range falseContainer[idx:] {
					if !insertedMap[x.Start.Format("20060102150504")+"-"+x.End.Format("20060102150504")] {
						xCount++
						newFalseContainer = append(newFalseContainer, FalseContainer{xCount, x.Start, x.End})
					}
				}
				break

			} else if (startInt >= ctStartInt && startInt <= ctEndInt) && endInt > ctEndInt { // start inside, end outside
				//log.Println(3)
				newFalseContainer = append(newFalseContainer, FalseContainer{idx, ct.Start, end})
				insertedMap[tk.ToString(ctStartInt)+"-"+tk.ToString(endInt)] = true
			} else if startInt < ctStartInt && (endInt >= ctStartInt && endInt <= ctEndInt) { // end inside, start outside
				//log.Println(4)
				newFalseContainer = append(newFalseContainer, FalseContainer{idx, start, ct.End})
				insertedMap[tk.ToString(startInt)+"-"+tk.ToString(ctEndInt)] = true

				xCount := idx
				for _, x := range falseContainer[idx:] {
					if !insertedMap[x.Start.Format("20060102150504")+"-"+x.End.Format("20060102150504")] {
						xCount++
						newFalseContainer = append(newFalseContainer, FalseContainer{xCount, x.Start, x.End})
					}
				}
				break
			} else if startInt < ctStartInt && endInt < ctStartInt { // outside all, before
				//log.Println(5)
				newFalseContainer = append(newFalseContainer, FalseContainer{idx, start, end})
				idx++
				newFalseContainer = append(newFalseContainer, FalseContainer{idx, ct.Start, ct.End})

				insertedMap[tk.ToString(startInt)+"-"+tk.ToString(endInt)] = true
				insertedMap[tk.ToString(ctStartInt)+"-"+tk.ToString(ctEndInt)] = true

				xCount := idx
				for _, x := range falseContainer[idx-1:] {
					if !insertedMap[x.Start.Format("20060102150504")+"-"+x.End.Format("20060102150504")] {
						xCount++
						newFalseContainer = append(newFalseContainer, FalseContainer{xCount, x.Start, x.End})
					}
				}
				break
			} else {
				//log.Println(6)
				newFalseContainer = append(newFalseContainer, FalseContainer{idx, ct.Start, ct.End})
				insertedMap[tk.ToString(ctStartInt)+"-"+tk.ToString(ctEndInt)] = true

				if idx == len(falseContainer) {
					//log.Println(7)
					idx++
					newFalseContainer = append(newFalseContainer, FalseContainer{idx, start, end})
					insertedMap[tk.ToString(startInt)+"-"+tk.ToString(endInt)] = true
				}
			}
		}
	}

	return newFalseContainer
}
