package controller

import (
	. "eaciit/wfdemo-git/library/core"
	. "eaciit/wfdemo-git/library/helper"
	. "eaciit/wfdemo-git/library/models"

	"eaciit/wfdemo-git/web/helper"
	// "fmt"
	// "strconv"
	"sort"
	"strings"
	// "time"
	// c "github.com/eaciit/crowd"
	// "github.com/eaciit/dbox"
	"github.com/eaciit/knot/knot.v1"
	tk "github.com/eaciit/toolkit"
)

type AnalyticMeteorologyController struct {
	App
}

func CreateAnalyticMeteorologyController() *AnalyticMeteorologyController {
	var controller = new(AnalyticMeteorologyController)
	return controller
}

func (m *AnalyticMeteorologyController) GetWindCorrelation(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	var dataSeries []tk.M
	p := new(PayloadAnalytic)
	e := k.GetPayload(&p)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	tStart, tEnd, e := helper.GetStartEndDate(k, p.Period, p.DateStart, p.DateEnd)
	query := []tk.M{}
	pipes := []tk.M{}
	query = append(query, tk.M{"_id": tk.M{"$ne": ""}})
	query = append(query, tk.M{"dateinfo.dateid": tk.M{"$gte": tStart}})
	query = append(query, tk.M{"dateinfo.dateid": tk.M{"$lte": tEnd}})
	// query = append(query, tk.M{"avgwindspeed": tk.M{"$gte": 0.5}})
	if p.Project != "" {
		anProject := strings.Split(p.Project, "(")
		query = append(query, tk.M{"projectname": strings.TrimRight(anProject[0], " ")})
	}

	pipes = append(pipes, tk.M{"$match": tk.M{"$and": query}})
	pipes = append(pipes, tk.M{"$project": tk.M{"turbine": 1, "avgwindspeed": 1, "timestamp": 1}})

	csr, err := DB().Connection.NewQuery().From(new(ScadaData).TableName()).
		Command("pipe", pipes).Cursor(nil)

	if err != nil {
		return helper.CreateResult(false, nil, err.Error())
	}

	allres := tk.M{}
	arrturbine := []string{}
	_tturbine := tk.M{}

	for {
		trx := new(ScadaData)
		e := csr.Fetch(trx, 1, false)
		if e != nil {
			break
		}

		dkey := trx.TimeStamp.Format("20060102030405")

		_tkm := allres.Get(trx.Turbine, tk.M{}).(tk.M)
		if trx.AvgWindSpeed != -99999.0 {
			_tkm.Set(dkey, tk.ToFloat64(trx.AvgWindSpeed, 6, tk.RoundingAuto))
		}

		allres.Set(trx.Turbine, _tkm)
		_tturbine.Set(trx.Turbine, 1)
	}

	for key, _ := range _tturbine {
		arrturbine = append(arrturbine, key)
	}

	sort.Strings(arrturbine)
	pturbine := append([]string{}, arrturbine...)
	arrturbine = append([]string{"Turbine"}, arrturbine...)

	if len(p.Turbine) > 0 {
		pturbine = []string{}
		for _, _v := range p.Turbine {
			pturbine = append(pturbine, tk.ToString(_v))
		}
	}

	for _, _turbine := range pturbine {
		_tkm := tk.M{}.Set("Turbine", _turbine)
		for i := 1; i < len(arrturbine); i++ {
			if arrturbine[i] != _turbine {
				_tkm.Set(arrturbine[i],
					GetCorrelation(allres.Get(_turbine, tk.M{}).(tk.M), allres.Get(arrturbine[i], tk.M{}).(tk.M)))
			} else {
				_tkm.Set(arrturbine[i], "")
			}
		}
		dataSeries = append(dataSeries, _tkm)
	}

	data := struct {
		Column []string
		Data   []tk.M
	}{
		Column: arrturbine,
		Data:   dataSeries,
	}

	return helper.CreateResult(true, data, "success")
}
