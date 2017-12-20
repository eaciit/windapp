package controller

import (
	. "eaciit/wfdemo-git/library/core"
	. "eaciit/wfdemo-git/library/models"
	"eaciit/wfdemo-git/web/helper"
	"math"
	"sort"
	"time"

	"github.com/eaciit/crowd"
	"github.com/eaciit/dbox"
	"github.com/eaciit/knot/knot.v1"
	tk "github.com/eaciit/toolkit"
)

var (
	// colorFieldTLP = [...]string{"#87c5da","#cc2a35", "#d66b76", "#5d1b62", "#f1c175","#95204c","#8f4bc5","#7d287d","#00818e","#c8c8c8","#546698","#66c99a","#f3d752","#20adb8","#333d6b","#d077b1","#aab664","#01a278","#c1d41a","#807063","#ff5975","#01a3d4","#ca9d08","#026e51","#4c653f","#007ca7"}
	colorFieldTLP = [...]string{"#ff9933", "#21c4af", "#ff7663", "#ffb74f", "#a2df53", "#1c9ec4", "#ff63a5", "#f44336", "#69d2e7", "#8877A9", "#9A12B3", "#26C281", "#E7505A", "#C49F47", "#ff5597", "#c3260c", "#d4735e", "#ff2ad7", "#34ac8b", "#11b2eb", "#004c79", "#ff0037", "#507ca3", "#ff6565", "#ffd664", "#72aaff", "#795548", "#383271", "#6a4795", "#bec554", "#ab5919", "#f5b1e1", "#7b3416", "#002fef", "#8d731b", "#1f8805", "#ff9900", "#9C27B0", "#6c7d8a", "#d73c1c", "#5be7a0", "#da02d4", "#afa56e", "#7e32cb", "#a2eaf7", "#9cb8f4", "#9E9E9E", "#065806", "#044082", "#18937d", "#2c787a", "#a57c0c", "#234341", "#1aae7a", "#7ac610", "#736f5f", "#4e741e", "#68349d", "#1df3b6", "#e02b09", "#d9cfab", "#6e4e52", "#f31880", "#7978ec", "#f5ace8", "#3db6ae", "#5e06b0", "#16d0b9", "#a25a5b", "#1e603a", "#4b0981", "#62975f", "#1c8f2f", "#b0c80c", "#642794", "#e2060d", "#2125f0"}
)

type TrendLinePlotsController struct {
	App
}

func CreateTrendLinePlotsController() *TrendLinePlotsController {
	var controller = new(TrendLinePlotsController)
	return controller
}

func (m *TrendLinePlotsController) GetList(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	var (
		pipes        []tk.M
		filter       []*dbox.Filter
		list         []tk.M
		dataSeries   []tk.M
		sortTurbines []string
		categories   []string
	)

	p := new(PayloadAnalyticTLP)
	e := k.GetPayload(&p)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	tStart, tEnd, e := helper.GetStartEndDate(k, p.Period, p.DateStart, p.DateEnd)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}
	startdate := tStart.Day()
	enddate := tEnd.Day()
	turbine := p.Turbine
	colName := p.ColName
	deviationStatus := p.DeviationStatus
	deviation := p.Deviation
	project := p.Project

	// dateRange := 0

	minValue := 100.0
	maxValue := 0.0
	selArr := 1
	var listMonth []int
	catTitle := ""
	listOfYears := []int{}

	colId := "$dateinfo.dateid"

	/*==================== CREATING CATEGORY PART ====================*/
	for i := tStart.Year(); i <= tEnd.Year(); i++ {
		listOfYears = append(listOfYears, i)
	}

	_, months, monthDay := helper.GetDurationInMonth(tStart, tEnd)
	for _, val := range months {
		listMonth = append(listMonth, tk.ToInt(val, tk.RoundingAuto))
	}
	sort.Ints(listMonth)
	categoryChecker := []string{}
	lastMonth := 0
	idxYear := 0

	for lm, lMonth := range listMonth {
		if lm == 0 { /*bulan pertama*/
			catTitle = tStart.Month().String()
			if len(listMonth) == 1 { /*jika hanya 1 bulan*/
				for iDate := startdate; iDate <= enddate; iDate++ {
					categories = append(categories, tk.ToString(iDate))
					/*category checker akan berisi tanggal_bulan_tahun*/
					categoryChecker = append(categoryChecker, tk.ToString(iDate)+"_"+tk.ToString(lMonth)+"_"+tk.ToString(listOfYears[idxYear]))
				}
				catTitle += " " + tk.ToString(listOfYears[0]) /*Dec 2015*/
			} else {
				month := lMonth
				maxDays := monthDay.Get(tk.ToString(tStart.Year()) + tk.ToString(month)).(tk.M).GetInt("totalInMonth")
				for iDate := startdate; iDate <= maxDays; iDate++ {
					categories = append(categories, tk.ToString(iDate))
					categoryChecker = append(categoryChecker, tk.ToString(iDate)+"_"+tk.ToString(lMonth)+"_"+tk.ToString(listOfYears[idxYear]))
				}
				if len(listOfYears) > 1 { /*jika lebih dari 1 tahun, lanjut ke berikutnya*/
					catTitle += " " + tk.ToString(listOfYears[0]) /* Dec 2015*/
				}
				lastMonth = lMonth
			}
		} else { /*bulan selanjutnya*/
			if lastMonth > lMonth { /*jika bulan lalu lebih besar dari bulan saat ini maka ganti tahun*/
				idxYear++
			}
			if lm == len(listMonth)-1 { /*bulan terakhir*/
				catTitle += " - " + tEnd.Month().String()
				if len(listOfYears) == 1 {
					catTitle += " (" + tk.ToString(listOfYears[0]) + ")" /*Dec - Jan (2016)*/
				} else {
					catTitle += " " + tk.ToString(listOfYears[1]) /* - Jan 2016*/
				}
				for iDate := 1; iDate <= enddate; iDate++ {
					categories = append(categories, tk.ToString(iDate))
					categoryChecker = append(categoryChecker, tk.ToString(iDate)+"_"+tk.ToString(lMonth)+"_"+tk.ToString(listOfYears[idxYear]))
				}
			} else {
				month := lMonth
				maxDays := monthDay.Get(tk.ToString(tStart.Year()) + tk.ToString(month)).(tk.M).GetInt("totalInMonth")
				for iDate := 1; iDate <= maxDays; iDate++ {
					categories = append(categories, tk.ToString(iDate))
					categoryChecker = append(categoryChecker, tk.ToString(iDate)+"_"+tk.ToString(lMonth)+"_"+tk.ToString(listOfYears[idxYear]))
				}
				lastMonth = lMonth
			}
		}
	}
	/*==================== END OF CREATING CATEGORY PART ====================*/

	/*============================== AVG TLP PART ============================*/
	AvgTlp, TLPavgData, e := getTLPavgData(tStart, tEnd, colName, project, categoryChecker)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}
	dataSeries = append(dataSeries, TLPavgData)
	/*============================== END OF AVG TLP PART ============================*/

	/*================================= MET TOWER PART =================================*/
	metData := tk.M{}
	metData.Set("name", "Met Tower")
	metData.Set("type", "line")
	metData.Set("style", "smooth")
	metData.Set("dashType", "solid")
	metData.Set("markers", tk.M{"visible": false})
	metData.Set("width", 2)
	metData.Set("color", colorFieldTLP[selArr])
	metData.Set("idxseries", selArr)
	// if colName == "temp_yawbrake_1" {
	if colName == "temp_outdoor" {
		matches := []tk.M{
			tk.M{"_id": tk.M{"$ne": ""}},
			tk.M{"timestamp": tk.M{"$gte": tStart}},
			tk.M{"timestamp": tk.M{"$lte": tEnd}},
		}
		pipes = []tk.M{
			tk.M{"$match": tk.M{"$and": matches}},
		}
		if project != "" {
			filter = append(filter, dbox.Eq("projectname", project))
		}
		if tStart.Before(time.Date(2017, 1, 1, 0, 0, 0, 0, time.UTC)) {
			pipes = append(pipes, tk.M{"$group": tk.M{
				"_id":       tk.M{"colId": "$dateinfo.dateid"},
				"colresult": tk.M{"$avg": "$trefhreftemp855mavg"},
			}})
		} else {
			pipes = append(pipes, tk.M{"$group": tk.M{
				"_id":       tk.M{"colId": "$dateinfo.dateid"},
				"colresult": tk.M{"$avg": "$trefhrefhumid855mavg"},
			}})
		}

		pipes = append(pipes, tk.M{"$sort": tk.M{"_id": 1}})

		csrMet, e := DB().Connection.NewQuery().
			From(new(MetTower).TableName()).
			Command("pipe", pipes).
			Where(dbox.And(filter...)).
			Cursor(nil)

		if e != nil {
			return helper.CreateResult(false, nil, e.Error())
		}
		listMet := []tk.M{}
		e = csrMet.Fetch(&listMet, 0, false)
		defer csrMet.Close()

		var datas []float64

		idxAvgTlp := 0
		shownSeries := false

		dateFound := false
		for _, tanggal := range categoryChecker {
			dateFound = false
		metLoop:
			for _, val := range listMet {
				ids := val["_id"].(tk.M)
				tgl := ids.Get("colId").(time.Time)
				tglString := tk.ToString(tgl.Day()) + "_" + tk.ToString(int(tgl.Month())) + "_" + tk.ToString(tgl.Year())
				if tglString == tanggal {
					dateFound = true
					/*calculation process*/
					colresult := val.GetFloat64("colresult")
					if math.Abs(AvgTlp[idxAvgTlp]-colresult) > deviation {
						shownSeries = true
					}

					datas = append(datas, colresult)

					if colresult < minValue {
						minValue = colresult
					}
					if colresult > maxValue {
						maxValue = colresult
					}
					idxAvgTlp++
					break metLoop
				}
			}
			if !dateFound {
				datas = append(datas, 999999)
			}
		}
		if deviationStatus {
			if shownSeries {
				if len(datas) > 0 {
					metData.Set("data", datas)
				}
			}
		} else {
			if len(datas) > 0 {
				metData.Set("data", datas)
			}
		}
		selArr++
	} else {
		selArr++
	}
	dataSeries = append(dataSeries, metData)
	/*================================= END OF MET TOWER PART =================================*/

	/*==================== SCADA DATA OEM PART ====================*/
	matches := []tk.M{
		tk.M{"projectname": project},
		tk.M{"isnull": false},
		tk.M{"timestamp": tk.M{"$gte": tStart}},
		tk.M{"timestamp": tk.M{"$lte": tEnd}},
		// TEMPORARY SOLUTION,sampek disemprot mak e gara2 gak ngerti kudu lapho, kita selalu salah di mata wanita
		tk.M{colName: tk.M{"$lte": 200}},
	}

	pipes = []tk.M{
		tk.M{"$match": tk.M{"$and": matches}},
	}
	pipes = append(pipes, tk.M{"$group": tk.M{
		"_id":       tk.M{"colId": colId, "Turbine": "$turbine"},
		"colresult": tk.M{"$avg": "$" + colName}}})
	pipes = append(pipes, tk.M{"$sort": tk.M{"_id.Turbine": 1}})

	csr, e := DB().Connection.NewQuery().
		From("Scada10MinHFD").
		Command("pipe", pipes).
		Cursor(nil)

	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	_list := tk.M{}
	list = []tk.M{}
	for {
		_list = tk.M{}
		e = csr.Fetch(&_list, 1, false)
		if e != nil {
			break
		}
		list = append(list, _list)
	}
	defer csr.Close()

	if len(p.Turbine) == 0 {
		for _, listVal := range list {
			exist := false
			for _, val := range turbine {
				if listVal["_id"].(tk.M)["Turbine"] == val {
					exist = true
				}
			}
			if exist == false {
				turbine = append(turbine, listVal["_id"].(tk.M)["Turbine"])
			}
		}
	}

	for _, turX := range turbine {
		sortTurbines = append(sortTurbines, turX.(string))
	}
	sort.Strings(sortTurbines)
	turbineName, e := helper.GetTurbineNameList(p.Project)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}
	for _, turbineX := range sortTurbines {
		exist := crowd.From(&list).Where(func(x interface{}) interface{} {
			y := x.(tk.M)
			id := y.Get("_id").(tk.M)

			return id.GetString("Turbine") == turbineX
		}).Exec().Result.Data().([]tk.M)

		var datas []float64
		turbineData := tk.M{}
		turbineData.Set("name", turbineName[turbineX])
		turbineData.Set("type", "line")
		turbineData.Set("style", "smooth")
		turbineData.Set("dashType", "solid")
		turbineData.Set("markers", tk.M{"visible": false})
		turbineData.Set("width", 2)
		turbineData.Set("color", colorFieldTLP[selArr])
		turbineData.Set("idxseries", selArr)

		idxAvgTlp := 0
		shownSeries := false
		dateFound := false
		for _, tanggal := range categoryChecker {
			dateFound = false
		existLoop:
			for _, val := range exist {
				ids := val["_id"].(tk.M)
				tgl := ids.Get("colId").(time.Time)
				tglString := tk.ToString(tgl.Day()) + "_" + tk.ToString(int(tgl.Month())) + "_" + tk.ToString(tgl.Year())
				if tglString == tanggal { /*jika tanggal di dalam aggregate result ada di dalam category date*/
					dateFound = true
					/*calculation process*/
					colresult := val.GetFloat64("colresult")
					if math.Abs(AvgTlp[idxAvgTlp]-colresult) > deviation {
						shownSeries = true
					}

					datas = append(datas, colresult)

					if colresult < minValue {
						minValue = colresult
					}
					if colresult > maxValue {
						maxValue = colresult
					}
					idxAvgTlp++
					break existLoop
				}
			}
			if !dateFound { /*jika tanggal di dalam aggregate result tidak ditemukan di dalam category date*/
				datas = append(datas, 999999)
			}
		}

		if deviationStatus {
			if shownSeries {
				if len(datas) > 0 {
					turbineData.Set("data", datas)
				}
			}
		} else {
			if len(datas) > 0 {
				turbineData.Set("data", datas)
			}
		}

		dataSeries = append(dataSeries, turbineData)
		selArr++
	}
	/*==================== END OF SCADA DATA OEM PART ====================*/

	for _, val := range AvgTlp {

		if val < minValue && val < 999999 {
			minValue = val
		}
		if val > maxValue {
			maxValue = val
		}
	}

	data := struct {
		Data        []tk.M
		Categories  []string
		CatTitle    string
		Min         int
		Max         int
		TurbineName map[string]string
	}{
		Data:        dataSeries,
		Categories:  categories,
		CatTitle:    catTitle,
		Min:         tk.ToInt((minValue - 2), tk.RoundingAuto),
		Max:         tk.ToInt((maxValue + 2), tk.RoundingAuto),
		TurbineName: turbineName,
	}

	return helper.CreateResult(true, data, "success")
}

/**
 * @param  {[
 * Turbine    []interface{}
	DateStart  time.Time
	DateEnd    time.Time]}
 * @return {pcData}
*/

func getTLPavgData(DateStart time.Time, DateEnd time.Time, colName string, project string, categoryChecker []string) (datas []float64, pcData tk.M, e error) {

	var (
		pipes []tk.M
		list  []tk.M
	)

	matches := []tk.M{
		tk.M{"isnull": false},
		tk.M{"dateinfo.dateid": tk.M{"$gte": DateStart}},
		tk.M{"dateinfo.dateid": tk.M{"$lte": DateEnd}},
		// TEMPORARY SOLUTION
		tk.M{colName: tk.M{"$lte": 200}},
	}
	if project != "" {
		matches = append(matches, tk.M{"projectname": project})
	}
	pipes = []tk.M{
		tk.M{"$match": tk.M{"$and": matches}},
	}

	pipes = append(pipes, tk.M{"$group": tk.M{"_id": "$dateinfo.dateid", "colresult": tk.M{"$avg": "$" + colName}}})
	pipes = append(pipes, tk.M{"$sort": tk.M{"_id": 1}})

	csr, e := DB().Connection.NewQuery().
		From("Scada10MinHFD").
		Command("pipe", pipes).
		Cursor(nil)

	if e != nil {
		return
	}
	e = csr.Fetch(&list, 0, false)
	defer csr.Close()

	dateFound := false
	for _, tanggal := range categoryChecker {
		dateFound = false
	existLoop:
		for _, val := range list {
			tgl := val.Get("_id", time.Time{}).(time.Time)
			tglString := tk.ToString(tgl.Day()) + "_" + tk.ToString(int(tgl.Month())) + "_" + tk.ToString(tgl.Year())
			if tglString == tanggal { /*jika tanggal di dalam aggregate result ada di dalam category date*/
				dateFound = true
				colresult := val.GetFloat64("colresult")

				datas = append(datas, colresult)
				break existLoop
			}
		}
		if !dateFound { /*jika tanggal di dalam aggregate result tidak ditemukan di dalam category date*/
			datas = append(datas, 999999)
		}
	}

	pcData = tk.M{
		"name":      "Average",
		"idxseries": 0,
		"type":      "line",
		"dashType":  "longDash",
		"style":     "smooth",
		"color":     "#000000",
		"markers":   tk.M{"visible": false},
		"width":     3,
	}

	if len(datas) > 0 {
		pcData.Set("data", datas)
	}

	return
}
