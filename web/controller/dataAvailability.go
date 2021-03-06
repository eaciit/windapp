package controller

import (
	. "eaciit/wfdemo-git/library/core"
	. "eaciit/wfdemo-git/library/helper"
	. "eaciit/wfdemo-git/library/models"
	"eaciit/wfdemo-git/web/helper"

	"time"

	"sort"

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

	if project != "" && project != "Fleet" {
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

	}

	rturbines := tk.ToFloat64(tEnd.UTC().Sub(tStart.UTC()).Hours()/24, 0, tk.RoundingUp) * 144
	iturbine, totalrows := float64(0), float64(0)
	for _, val := range resultScada {
		iturbine += 1
		totalrows += val.GetFloat64("totalrows")
	}

	return helper.CreateResult(true, tk.Div(totalrows, rturbines*iturbine), "success")
}

func setOpacity(scadaAvail float64) float64 {
	switch {
	case scadaAvail >= 0.0 && scadaAvail <= 0.1:
		return 1
	case scadaAvail > 0.1 && scadaAvail <= 0.2:
		return 0.9
	case scadaAvail > 0.2 && scadaAvail <= 0.3:
		return 0.8
	case scadaAvail > 0.3 && scadaAvail <= 0.4:
		return 0.7
	case scadaAvail > 0.4 && scadaAvail <= 0.5:
		return 0.6
	case scadaAvail > 0.5 && scadaAvail <= 0.6:
		return 0.6
	case scadaAvail > 0.6 && scadaAvail <= 0.7:
		return 0.7
	case scadaAvail > 0.7 && scadaAvail <= 0.8:
		return 0.8
	case scadaAvail > 0.8 && scadaAvail <= 0.9:
		return 0.9
	case scadaAvail > 0.9 && scadaAvail <= 1.0:
		return 1
	}

	return 1
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
		tk.Println("error parse time", e.Error())
		return tk.M{}
	}
	/* bulan di filter ditambah 1 bulan kemudian dikurangi 1 hari, didapatkan jumlah max hari pada bulan di filter */
	timeParse = timeParse.AddDate(0, 1, -1)

	pipes = append(pipes, tk.M{"$sort": tk.M{"_id": 1}})

	csr, e := DB().Connection.NewQuery().
		From(new(ScadaSummaryDaily).TableName()).
		Command("pipe", pipes).
		Cursor(nil)

	if e != nil {
		tk.Println("error cursor", e.Error())
		return tk.M{}
	}
	defer csr.Close()

	e = csr.Fetch(&dailyData, 0, false) /* data per hari selama 1 bulan untuk parent sesuai filter */
	if e != nil {
		tk.Println("error fetching data", e.Error())
		return tk.M{}
	}

	result := tk.M{}

	if len(dailyData) > 0 {
		totalDay := timeParse.Day() /* maksimum hari pada bulan filter */
		dataPerDay := tk.M{}
		timeConv := time.Time{}
		for _, val := range dailyData {
			timeConv = val.Get("_id", time.Time{}).(time.Time)
			dataPerDay.Set(tk.ToString(timeConv.Day()), val.GetFloat64("scadaavail"))
		}
		datas := []tk.M{}
		percentage := 0.0
		kelas := "progress-bar progress-bar-success"
		for idx := 1; idx <= totalDay; idx++ { /* untuk mendapatkan data tiap hari secara urut */
			percentage = 1.0 / tk.ToFloat64(totalDay, 6, tk.RoundingAuto) * 100 /* percentage 1 hari dibanding totalHari dalam 1 bulan*/
			if dataPerDay.Has(tk.ToString(idx)) {
				if dataPerDay.GetFloat64(tk.ToString(idx)) < 0.5 { /* jika kurang dari 0.5 availability nya, maka warnanya merah */
					kelas = "progress-bar progress-bar-red"
				} else {
					kelas = "progress-bar progress-bar-success"
				}

				percentageDay := tk.ToFloat64(dataPerDay.Get(tk.ToString(idx), 0.0).(float64), 4, tk.RoundingAuto) * 100

				datas = append(datas, tk.M{
					"tooltip":  "Day " + tk.ToString(idx),
					"class":    kelas,
					"value":    tk.ToString(percentage) + "%",
					"floatval": percentageDay,
					"opacity":  setOpacity(dataPerDay.GetFloat64(tk.ToString(idx))),
				})
			} else { /* default value jika tidak ada data availability pada hari tersebut */
				datas = append(datas, tk.M{
					"tooltip":  "Day " + tk.ToString(idx),
					"class":    "progress-bar progress-bar-red",
					"value":    tk.ToString(percentage) + "%",
					"floatval": 0,
					"opacity":  1,
				})
			}
		}
		result = tk.M{"Category": "Data Availability", "Turbine": []tk.M{}, "Data": datas}

		/* ========= Query Drill Down per Turbine ================ */
		dailyTurbine := []tk.M{}
		pipes = []tk.M{}
		pipes = append(pipes, tk.M{"$match": tk.M{"$and": query}})
		pipes = append(pipes, tk.M{"$group": tk.M{
			"_id":        tk.M{"tanggal": "$dateinfo.dateid", "turbine": "$turbine"},
			"scadaavail": tk.M{"$avg": "$scadaavail"},
		}})

		csrTurbine, e := DB().Connection.NewQuery().
			From(new(ScadaSummaryDaily).TableName()).
			Command("pipe", pipes).
			Cursor(nil)

		if e != nil {
			tk.Println("error cursor summary daily", e.Error())
			return tk.M{}
		}
		defer csrTurbine.Close()

		e = csrTurbine.Fetch(&dailyTurbine, 0, false) /* data per hari per turbine selama 1 bulan sesuai filter*/
		if e != nil {
			tk.Println("error fetching data daily turbine", e.Error())
			return tk.M{}
		}

		if len(dailyTurbine) > 0 {
			dataPerTurbine := tk.M{}
			timeConv = time.Time{}
			turbineList := []string{}
			turbineMap := map[string]bool{} /* untuk mengambil unique turbine */
			ids := tk.M{}
			for _, val := range dailyTurbine { /* pembentukan map data agar dapat diakses secara mudah sesuai turbine dan tanggal */
				ids = val.Get("_id", tk.M{}).(tk.M)
				turbineMap[ids.GetString("turbine")] = true
				timeConv = ids.Get("tanggal", time.Time{}).(time.Time)
				dataPerTurbine.Set(tk.ToString(timeConv.Day())+"_"+ids.GetString("turbine"), val.GetFloat64("scadaavail"))
			}
			for key := range turbineMap {
				turbineList = append(turbineList, key)
			}
			sort.Strings(turbineList)

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
			lastTurbine := ""

			for _, turbine := range turbineList { /* dengan metode ini sudah pasti urut turbinenya */
				_turbine = turbineName[turbine]
				if lastTurbine != turbine {
					lastTurbine = turbine
					turbineDetails = []tk.M{}
				}
				for idx := 1; idx <= totalDay; idx++ { /* dengan metode ini sudah pasti urut harinya */
					percentage = 1.0 / tk.ToFloat64(totalDay, 6, tk.RoundingAuto) * 100 /* percentage 1/totalHari selama 1 bulan */
					if dataPerTurbine.Has(tk.ToString(idx) + "_" + turbine) {
						if dataPerTurbine.GetFloat64(tk.ToString(idx)+"_"+turbine) < 0.5 {
							kelas = "progress-bar progress-bar-red"
						} else {
							kelas = "progress-bar progress-bar-success"
						}

						percentageTurbine := tk.ToFloat64(dataPerTurbine.Get(tk.ToString(idx)+"_"+turbine, 0.0).(float64), 4, tk.RoundingAuto) * 100

						turbineDetails = append(turbineDetails, tk.M{
							"tooltip":  "Day " + tk.ToString(idx),
							"class":    kelas,
							"value":    tk.ToString(percentage) + "%",
							"floatval": percentageTurbine,
							"opacity":  setOpacity(dataPerTurbine.GetFloat64(tk.ToString(idx) + "_" + turbine)),
						})
					} else { /* data default jika tidak ada data availability di hari tersebut */
						turbineDetails = append(turbineDetails, tk.M{
							"tooltip":  "Day " + tk.ToString(idx),
							"class":    "progress-bar progress-bar-red",
							"value":    tk.ToString(percentage) + "%",
							"floatval": 0,
							"opacity":  1,
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

func getParentData(project string, collType string) []tk.M {
	var (
		pipes []tk.M
		list  []tk.M
	)
	group := tk.M{
		"_id": tk.M{
			"name":    "$name",
			"project": "$details.projectname",
			"index":   "$details.id",
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
		"index":      "$_id.index",
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

	if match.Get("project") != nil {
		pipes = append(pipes, tk.M{"$match": match})
	}

	pipes = append(pipes, tk.M{"$sort": tk.M{"index": 1}})

	csr, e := DB().Connection.NewQuery().
		From(new(DataAvailability).TableName()).
		Command("pipe", pipes).
		Cursor(nil)

	if e != nil {
		tk.Println("error cursor", e.Error())
		return []tk.M{}
	}

	e = csr.Fetch(&list, 0, false)
	if e != nil {
		tk.Println("error fetching data", e.Error())
		return []tk.M{}
	}

	defer csr.Close()

	datas := []tk.M{}

	type DataStruct struct {
		Start    time.Time
		End      time.Time
		IsAvail  bool
		Duration float64
	}

	if len(list) > 0 {
		turbineDetailsData := []DataStruct{}
		currTotalDuration := 0.0
		for _, dt := range list {
			pTo := dt.Get("periodTo").(time.Time)
			pFrom := dt.Get("periodFrom").(time.Time)

			from = pFrom.UTC()
			to = pTo.UTC()
			totalDurationDays := to.UTC().Sub(from.UTC()).Hours() / 24 /* total hari parent yang ada di database */

			availList := dt.Get("list").([]interface{})

			for _, av := range availList {
				avail := av.(tk.M)
				start := avail.Get("start").(time.Time).UTC()
				end := avail.Get("end").(time.Time).UTC()
				duration := avail.GetFloat64("duration")
				isAvail := avail.Get("isavail").(bool)

				/* hitung dulu total durasi detail di database supaya bisa mendekati 100 %
				untuk jaga2 sapa tau hasil summary generator antara parent dan details berbeda */
				currTotalDuration += duration

				if (currTotalDuration - totalDurationDays) < 1 { /* jika total kelebihan waktu tidak lebih dari 1 hari */
					turbineDetailsData = append(turbineDetailsData, DataStruct{
						start,
						end,
						isAvail,
						duration,
					})
				} else { /* [SOLUSI TEMPORARY] jika lebih dari 1 hari maka mungkin kesalahan summary generator nya*/
					currTotalDuration -= duration
				}
				// datas = append(datas, setDataColumn(start, end, isAvail, duration))
			}
		}
		for _, dVal := range turbineDetailsData {
			datas = append(datas, setDataColumns(dVal.Start, dVal.End,
				dVal.IsAvail, dVal.Duration, currTotalDuration))
		}

	}
	return datas
}

func getAvailCollection(project string, turbines []interface{}, collType string) tk.M {
	var (
		pipes []tk.M
		list  []tk.M
	)
	group := tk.M{
		"_id": tk.M{
			"name":        "$name",
			"project":     "$details.projectname",
			"turbine":     "$details.turbine",
			"turbinename": "$details.turbinename",
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
		"name":        "$_id.name",
		"project":     "$_id.project",
		"turbine":     "$_id.turbine",
		"turbinename": "$_id.turbinename",
		"periodTo":    1,
		"periodFrom":  1,
		"list":        1,
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

	pipes = append(pipes, tk.M{"$sort": tk.M{"turbinename": 1}})

	csr, e := DB().Connection.NewQuery().
		From(new(DataAvailability).TableName()).
		Command("pipe", pipes).
		Cursor(nil)

	if e != nil {
		tk.Println("error cursor", e.Error())
		return tk.M{"Category": "", "Turbine": []tk.M{}, "Data": []tk.M{}}
	}

	e = csr.Fetch(&list, 0, false)
	if e != nil {
		tk.Println("error fetching data", e.Error())
		return tk.M{"Category": "", "Turbine": []tk.M{}, "Data": []tk.M{}}
	}

	defer csr.Close()

	res := []tk.M{}
	name := ""
	type DataStruct struct {
		Start    time.Time
		End      time.Time
		IsAvail  bool
		Duration float64
	}

	if len(list) > 0 {
		collTypeParent := collType + "_PROJECT"
		datas := getParentData(project, collTypeParent)

		for _, dt := range list {
			t := dt.GetString("turbinename")
			pTo := dt.Get("periodTo").(time.Time)
			pFrom := dt.Get("periodFrom").(time.Time)

			from = pFrom.UTC()
			to = pTo.UTC()
			totalDurationDays := to.UTC().Sub(from.UTC()).Hours() / 24 /* total hari parent yang ada di database */

			name = dt.GetString("name")
			availList := dt.Get("list").([]interface{})

			turbineDetails := []tk.M{}
			currTotalDuration := 0.0
			turbineDetailsData := []DataStruct{}

			// set availability data based on index ordering in collection
			for index := 1; index <= len(availList); index++ {
			breakAvail:
				for _, av := range availList {
					avail := av.(tk.M)
					if index == avail.GetInt("id") {
						start := avail.Get("start").(time.Time).UTC()
						end := avail.Get("end").(time.Time).UTC()
						duration := avail.GetFloat64("duration")
						isAvail := avail.Get("isavail").(bool)
						/* hitung dulu total durasi detail di database supaya bisa mendekati 100 %
						untuk jaga2 sapa tau hasil summary generator antara parent dan details berbeda */
						currTotalDuration += duration

						if (currTotalDuration - totalDurationDays) < 1 { /* jika total kelebihan waktu tidak lebih dari 1 hari */
							turbineDetailsData = append(turbineDetailsData, DataStruct{
								start,
								end,
								isAvail,
								duration,
							})
						} else { /* [SOLUSI TEMPORARY] jika lebih dari 1 hari maka mungkin kesalahan summary generator nya*/
							currTotalDuration -= duration
						}
						break breakAvail
					}
				}
			}
			for _, dVal := range turbineDetailsData {
				turbineDetails = append(turbineDetails, setDataColumns(dVal.Start, dVal.End,
					dVal.IsAvail, dVal.Duration, currTotalDuration))
			}

			turbine := tk.M{"TurbineName": t}
			turbine.Set("details", turbineDetails)

			res = append(res, turbine)
		}

		if collType == "MET_TOWER" && project != "Tejuva" { /* Khusus Met Tower selain Tejuva dimerahkan semuaa */
			datas = []tk.M{}
			datas = append(datas, tk.M{
				"tooltip":  from.Format("2 Jan 2006") + " until " + to.Format("2 Jan 2006"),
				"class":    "progress-bar progress-bar-red",
				"value":    "100%",
				"floatval": 100.0,
			})
		}

		if collType != "MET_TOWER" || (collType == "MET_TOWER" && project == "Tejuva") {
			return tk.M{"Category": name, "Turbine": res, "Data": datas}
		} else {
			return tk.M{"Category": name, "Turbine": []tk.M{}, "Data": datas}
		}

	}

	return tk.M{"Category": "", "Turbine": []tk.M{}, "Data": []tk.M{}}
}

func setDataColumns(start time.Time, end time.Time, avail bool, durationInDay, totalDurationDays float64) tk.M {
	res := tk.M{}
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
