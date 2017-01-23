package controller

import (
	. "eaciit/wfdemo-git/library/core"
	. "eaciit/wfdemo-git/library/models"
	"eaciit/wfdemo-git/web/helper"
	"time"
	"sort"
	"github.com/eaciit/crowd"
	"github.com/eaciit/dbox"
	"github.com/eaciit/knot/knot.v1"
	tk "github.com/eaciit/toolkit"
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
		categories		[]string
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
	// project := ""
	// if p.Project != "" {
	// 	anProject := strings.Split(p.Project, "(")
	// 	project = strings.TrimRight(anProject[0], " ")
	// }
	colName := p.ColName
	// breakDown := p.BreakDown

	colId := "$dateinfoutc.dateid"

	pipes = append(pipes, tk.M{"$group": tk.M{"_id": tk.M{"colId": colId, "Turbine": "$turbine"}, "colresult": tk.M{"$sum": "$"+colName}, "totaldata": tk.M{"$sum": 1}}})
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
		turbineData.Set("color", colorField[selArr])
		turbineData.Set("idxseries", selArr)

		for _, val := range exist {
			datas = append(datas, val.GetFloat64("colresult"))
		}

		if len(datas) > 0 {
			turbineData.Set("data", datas)
		}

		dataSeries = append(dataSeries, turbineData)
		selArr++
	}
	catTitle := ""

	MStart := tStart.Month()
	MEnd := tEnd.Month()

	iStart := int(MStart)
	iEnd := int(MEnd)

	jumMonth := iEnd - iStart

	listCount := 0
	monthNum := 0
	var listMonth []int


	if(jumMonth == 0){
		listMonth = append(listMonth, iStart)	
	}else{		
		for listCount <= jumMonth {
			monthNum = iStart + listCount	
			listMonth = append(listMonth, monthNum)		
			listCount = listCount + 1
		}
	}


	listOfDays := []int{0, 31, 28, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31}
	monthString := []string{"", "January", "February", "March", "April", "May", "June", "July",
		"August", "September", "October", "November", "December"}

	listOfYears := []int{}
	for i := tStart.Year(); i <= tEnd.Year(); i++ {
		listOfYears = append(listOfYears, i)
	}

		// tk.Printf("tEnd : %s \n", tEnd)

	for lm, lMonth := range listMonth {
		month := lMonth
		jumHari := listOfDays[month]
		catTitle = monthString[iStart]

		if lm == 0 { /*bulan pertama*/
			catTitle = monthString[iStart]
			if len(listMonth) == 1 {
				for iDate := startdate; iDate <= enddate; iDate++ {
					categories = append(categories, tk.ToString(iDate))
				}
				catTitle += " " + tk.ToString(listOfYears[0]) /*Dec 2015*/
			} else {
				for iDate := startdate; iDate <= jumHari; iDate++ {
					categories = append(categories, tk.ToString(iDate))
				}
				if len(listOfYears) > 1 { /*jika cuma 1 tahun, lanjut ke berikutnya*/
					catTitle += " " + tk.ToString(listOfYears[0]) /* Dec 2015*/
				}
			}
		} else { /*bulan kedua*/
			catTitle += " - " + monthString[iEnd]
			if len(listOfYears) == 1 {
				catTitle += " (" + tk.ToString(listOfYears[0]) + ")" /*Dec - Jan (2016)*/
			} else {
				catTitle += " " + tk.ToString(listOfYears[1]) /* - Jan 2016*/
			}
			for iDate := 1; iDate <= enddate; iDate++ {
				categories = append(categories, tk.ToString(iDate))
			}

		}

	}

	data := struct {
		Data []tk.M
		Categories []string
		CatTitle   string
	}{
		Data: dataSeries,
		Categories: categories,
		CatTitle: catTitle,
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