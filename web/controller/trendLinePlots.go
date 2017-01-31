package controller

import (
	. "eaciit/wfdemo-git/library/core"
	. "eaciit/wfdemo-git/library/models"
	"eaciit/wfdemo-git/web/helper"
	"github.com/eaciit/crowd"
	"github.com/eaciit/dbox"
	"github.com/eaciit/knot/knot.v1"
	tk "github.com/eaciit/toolkit"
	"sort"
	"time"
)

var (
	colorFieldTLP = [...]string{"#B71C1C", "#E57373", "#F44336", "#D81B60", "#F06292", "#880E4F",
		"#4A148C", "#7B1FA2", "#9C27B0", "#BA68C8", "#1A237E", "#5C6BC0",
		"#1E88E5", "#0277BD", "#0097A7", "#26A69A", "#4DD0E1", "#81C784",
		"#8BC34A", "#1B5E20", "#827717", "#C0CA33", "#DCE775", "#FF6F00", "#A1887F",
		"#FFEE58", "#004D40", "#212121", "#607D8B", "#BDBDBD", "#FF00CC", "#9999FF"}
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

	// dateRange := 0

	minValue := 100.0
	maxValue := 0.0
	var listMonth []int
	catTitle := ""
	listOfYears := []int{}

	colId := "$dateinfoutc.dateid"

	// if(dateRange == 0){
	// 	colId = "$timestamputc"
	// }

	AvgTlp, TLPavgData, e := getTLPavgData(tStart, tEnd, colName)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	dataSeries = append(dataSeries, TLPavgData)

	pipes = append(pipes, tk.M{"$group": tk.M{"_id": tk.M{"colId": colId, "Turbine": "$turbine"}, "colresult": tk.M{"$avg": "$" + colName}, "totaldata": tk.M{"$sum": 1}}})
	pipes = append(pipes, tk.M{"$sort": tk.M{"_id": 1}})

	selArr := 1

	filter = nil
	filter = append(filter, dbox.Ne("_id", ""))
	filter = append(filter, dbox.Gte("dateinfoutc.dateid", tStart))
	filter = append(filter, dbox.Lte("dateinfoutc.dateid", tEnd))
	filter = append(filter, dbox.Ne("turbine", ""))
	filter = append(filter, dbox.Ne("timestamp", ""))
	filter = append(filter, dbox.Ne("powerlost", ""))
	filter = append(filter, dbox.Ne("ai_intern_activpower", ""))
	filter = append(filter, dbox.Ne("ai_intern_windspeed", ""))

	csr, e := DB().Connection.NewQuery().
		From(new(ScadaDataOEM).TableName()).
		Command("pipe", pipes).
		Where(dbox.And(filter...)).
		Cursor(nil)

	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}
	e = csr.Fetch(&list, 0, false)
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

	for _, turbineX := range sortTurbines {

		exist := crowd.From(&list).Where(func(x interface{}) interface{} {
			y := x.(tk.M)
			id := y.Get("_id").(tk.M)

			return id.GetString("Turbine") == turbineX
		}).Exec().Result.Data().([]tk.M)

		var datas []float64
		turbineData := tk.M{}
		turbineData.Set("name", turbineX)
		turbineData.Set("type", "line")
		turbineData.Set("style", "smooth")
		turbineData.Set("dashType", "solid")
		turbineData.Set("markers", tk.M{"visible": false})
		turbineData.Set("width", 2)
		turbineData.Set("color", colorFieldTLP[selArr])
		turbineData.Set("idxseries", selArr)

		idxAvgTlp := 0
		shownSeries := false
		//colresult + deviation[idxAvgTlp]
		for _, val := range exist {

			calcColResult := 0.0
			colresult := val.GetFloat64("colresult")
			colresultMinus := colresult - deviation
			colresultPlus := colresult + deviation

			if colresult > AvgTlp[idxAvgTlp] {
				calcColResult = colresultMinus - AvgTlp[idxAvgTlp]
			} else {
				calcColResult = AvgTlp[idxAvgTlp] - colresultPlus
			}

			if calcColResult > 0.0 {
				shownSeries = true

				// tk.Printf("calcColResult : %s \n", calcColResult)
			}

			// if !shownSeries {
			// 	if calcColResult < 0 {
			// 		shownSeries = true
			// 	}
			// }

			datas = append(datas, colresult)

			if val.GetFloat64("colresult") < minValue {
				minValue = colresult
			}
			if val.GetFloat64("colresult") > maxValue {
				maxValue = colresult
			}
			idxAvgTlp = idxAvgTlp + 1

		}

		// tk.Printf("shownSeries : %s \n", shownSeries)
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
	metData := tk.M{}
	metData.Set("name", "Met Tower")
	metData.Set("type", "line")
	metData.Set("style", "smooth")
	metData.Set("dashType", "solid")
	metData.Set("markers", tk.M{"visible": false})
	metData.Set("width", 2)
	metData.Set("color", colorFieldTLP[selArr])
	metData.Set("idxseries", selArr)
	/*================================= MET TOWER PART =================================*/
	if colName == "temp_outdoor" {
		tk.Println(AvgTlp)
		// if colName == "kikuk" {
		pipes = []tk.M{}
		pipes = append(pipes, tk.M{"$group": tk.M{
			"_id":       tk.M{"colId": "$dateinfo.dateid"},
			"colresult": tk.M{"$avg": "$trefhreftemp855mavg"},
			"totaldata": tk.M{"$sum": 1},
		}})
		pipes = append(pipes, tk.M{"$sort": tk.M{"_id": 1}})

		filter = nil
		filter = append(filter, dbox.Ne("_id", ""))
		filter = append(filter, dbox.Gte("dateinfo.dateid", tStart))
		filter = append(filter, dbox.Lte("dateinfo.dateid", tEnd))

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
		for _, val := range listMet {

			calcColResult := 0.0
			colresult := val.GetFloat64("colresult")
			colresultMinus := colresult - deviation
			colresultPlus := colresult + deviation

			if colresult > AvgTlp[idxAvgTlp] {
				calcColResult = colresultMinus - AvgTlp[idxAvgTlp]
			} else {
				calcColResult = AvgTlp[idxAvgTlp] - colresultPlus
			}

			if calcColResult > 0.0 {
				shownSeries = true
			}

			datas = append(datas, colresult)

			if val.GetFloat64("colresult") < minValue {
				minValue = colresult
			}
			if val.GetFloat64("colresult") > maxValue {
				maxValue = colresult
			}
			idxAvgTlp = idxAvgTlp + 1

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
	}
	/*================================= END OF MET TOWER PART =================================*/
	dataSeries = append(dataSeries, metData)

	for _, val := range AvgTlp {

		if val < minValue {
			minValue = val
		}
		if val > maxValue {
			maxValue = val
		}
	}

	for i := tStart.Year(); i <= tEnd.Year(); i++ {
		listOfYears = append(listOfYears, i)
	}

	_, months, monthDay := helper.GetDurationInMonth(tStart, tEnd)
	for _, val := range months {
		listMonth = append(listMonth, tk.ToInt(val, tk.RoundingAuto))
	}
	sort.Ints(listMonth)

	for lm, lMonth := range listMonth {
		if lm == 0 { /*bulan pertama*/
			catTitle = tStart.Month().String()
			if len(listMonth) == 1 {
				for iDate := startdate; iDate <= enddate; iDate++ {
					categories = append(categories, tk.ToString(iDate))
				}
				catTitle += " " + tk.ToString(listOfYears[0]) /*Dec 2015*/
			} else {
				month := lMonth
				maxDays := monthDay.Get(tk.ToString(tStart.Year()) + tk.ToString(month)).(tk.M).GetInt("totalInMonth")
				for iDate := startdate; iDate <= maxDays; iDate++ {
					categories = append(categories, tk.ToString(iDate))
				}
				if len(listOfYears) > 1 { /*jika cuma 1 tahun, lanjut ke berikutnya*/
					catTitle += " " + tk.ToString(listOfYears[0]) /* Dec 2015*/
				}
			}
		} else { /*bulan selanjutnya*/
			if lm == len(listMonth)-1 { /*bulan terakhir*/
				catTitle += " - " + tEnd.Month().String()
				if len(listOfYears) == 1 {
					catTitle += " (" + tk.ToString(listOfYears[0]) + ")" /*Dec - Jan (2016)*/
				} else {
					catTitle += " " + tk.ToString(listOfYears[1]) /* - Jan 2016*/
				}
				for iDate := 1; iDate <= enddate; iDate++ {
					categories = append(categories, tk.ToString(iDate))
				}
			} else {
				month := lMonth
				maxDays := monthDay.Get(tk.ToString(tStart.Year()) + tk.ToString(month)).(tk.M).GetInt("totalInMonth")
				for iDate := 1; iDate <= maxDays; iDate++ {
					categories = append(categories, tk.ToString(iDate))
				}
			}
		}

	}

	data := struct {
		Data       []tk.M
		Categories []string
		CatTitle   string
		Min        int
		Max        int
	}{
		Data:       dataSeries,
		Categories: categories,
		CatTitle:   catTitle,
		Min:        tk.ToInt((minValue - 2), tk.RoundingAuto),
		Max:        tk.ToInt((maxValue + 2), tk.RoundingAuto),
	}

	return helper.CreateResult(true, data, "success")
}

func (m *TrendLinePlotsController) GetScadaOemAvailDate(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	Scadaresults := make([]time.Time, 0)

	// Scada Data
	for i := 0; i < 2; i++ {
		var arrsort []string
		if i == 0 {
			arrsort = append(arrsort, "timestamp")
		} else {
			arrsort = append(arrsort, "-timestamp")
		}

		query := DB().Connection.NewQuery().From(new(ScadaDataOEM).TableName()).Skip(0).Take(1)
		query = query.Order(arrsort...)

		csr, e := query.Cursor(nil)
		if e != nil {
			return helper.CreateResult(false, nil, e.Error())
		}
		defer csr.Close()

		Result := make([]ScadaDataOEM, 0)
		e = csr.Fetch(&Result, 0, false)

		if e != nil {
			return helper.CreateResult(false, nil, e.Error())
		}

		for _, val := range Result {
			Scadaresults = append(Scadaresults, val.TimeStampUTC.UTC())
		}
	}

	data := struct {
		ScadaOemAvailDate []time.Time
	}{
		ScadaOemAvailDate: Scadaresults,
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

func getTLPavgData(DateStart time.Time, DateEnd time.Time, colName string) (datas []float64, pcData tk.M, e error) {

	var (
		pipes  []tk.M
		filter []*dbox.Filter
		list   []tk.M
	)

	pipes = append(pipes, tk.M{"$group": tk.M{"_id": "$dateinfoutc.dateid", "colresult": tk.M{"$avg": "$" + colName}, "totaldata": tk.M{"$sum": 1}}})
	pipes = append(pipes, tk.M{"$sort": tk.M{"_id": 1}})

	filter = nil
	filter = append(filter, dbox.Ne("_id", ""))
	filter = append(filter, dbox.Gte("dateinfoutc.dateid", DateStart))
	filter = append(filter, dbox.Lte("dateinfoutc.dateid", DateEnd))
	// if(len(Turbine) > 0){
	// 	filter = append(filter, dbox.In("turbine", Turbine...))
	// }
	filter = append(filter, dbox.Ne("timestamp", ""))
	filter = append(filter, dbox.Ne("powerlost", ""))
	filter = append(filter, dbox.Ne("ai_intern_activpower", ""))
	filter = append(filter, dbox.Ne("ai_intern_windspeed", ""))

	csr, e := DB().Connection.NewQuery().
		From(new(ScadaDataOEM).TableName()).
		Command("pipe", pipes).
		Where(dbox.And(filter...)).
		Cursor(nil)

	if e != nil {
		return
	}
	e = csr.Fetch(&list, 0, false)
	defer csr.Close()

	// var datas []float64

	for _, val := range list {
		datas = append(datas, val.GetFloat64("colresult"))
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
