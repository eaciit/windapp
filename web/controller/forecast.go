package controller

import (
	. "eaciit/wfdemo-git/library/core"
	. "eaciit/wfdemo-git/library/models"
	"eaciit/wfdemo-git/web/helper"
	"os"
	"time"

	"github.com/eaciit/knot/knot.v1"
	tk "github.com/eaciit/toolkit"
)

type ForecastController struct {
	App
}

func CreateForecastController() *ForecastController {
	var controller = new(ForecastController)
	return controller
}

func get15MinPeriod(tstart time.Time, tend time.Time) []time.Time {
	timePeriods := []time.Time{}

	if tend.Sub(tstart).Minutes() >= 0 {
		currTime := tstart
		timePeriods = append(timePeriods, currTime)
		for {
			currTime = currTime.Add(time.Duration(15) * time.Minute)
			if currTime.Sub(tend).Minutes() > 0 {
				break
			}
			timePeriods = append(timePeriods, currTime)
		}
	}

	return timePeriods
}

func (m *ForecastController) GetList(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	p := new(PayloadAnalytic)
	e := k.GetPayload(&p)
	tk.Printf("%#v\n", p)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	tStart, tEnd, e := helper.GetStartEndDate(k, p.Period, p.DateStart, p.DateEnd)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	timeperiods := get15MinPeriod(tStart, tEnd)
	_ = timeperiods

	getscada15minpath := GetConfig("scada15min_path", "")
	scada15minpath := ""
	if getscada15minpath == "" || getscada15minpath == nil {
		scada15minpath = "/Users/masmeka/Works/Windfarm/Ostro/Scada15Min/data/"
	} else {
		scada15minpath = tk.ToString(getscada15minpath)
	}

	// get data forecast
	project := p.Project
	matches := []tk.M{
		tk.M{"projectname": project},
		tk.M{"timestamp": tk.M{"$gte": tStart}},
		tk.M{"timestamp": tk.M{"$lte": tEnd}},
	}
	pipes := []tk.M{
		tk.M{"$match": tk.M{"$and": matches}},
	}
	csr, e := DB().Connection.NewQuery().
		From(new(ForecastData).TableName()).
		Command("pipe", pipes).
		Cursor(nil)
	defer csr.Close()

	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}
	list := []tk.M{}
	for {
		item := tk.M{}
		e = csr.Fetch(&item, 1, false)
		if e != nil {
			break
		}
		timestamp := item.Get("timestamp", time.Time{}).(time.Time)
		if !timestamp.IsZero() {
			itemtoappend := tk.M{
				timestamp.Format("20060102_150405"): item,
			}
			list = append(list, itemtoappend)
		}
	}

	// get data scada 15 min
	scada := []tk.M{}
	_ = scada
	if _, err := os.Stat(scada15minpath); err == nil {

	}

	dataReturn := []tk.M{}
	for _, tp := range timeperiods {
		tpkey := tp.Format("20060102_150405")
		timeBefore := tp.Add(time.Duration(-15) * time.Minute)
		_ = tpkey
		item := tk.M{
			"Date":       tp.Format("02/01/2006"),
			"TimeBlock":  tk.Sprintf("%s - %s", tp.Format("15:04"), timeBefore.Format("15:04")),
			"AvaCap":     0.0,
			"Forecast":   0.0,
			"SchFcast":   0.0,
			"ExpProd":    0.0,
			"Actual":     0.0,
			"FcastWs":    0.0,
			"ActualWs":   0.0,
			"DevFcast":   0.0,
			"DevSchAct":  0.0,
			"DSMPenalty": "",
		}
		dataReturn = append(dataReturn, item)
	}

	return helper.CreateResult(true, dataReturn, "")
}
