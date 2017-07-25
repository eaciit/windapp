package controller

import (
	. "eaciit/wfdemo-git/library/core"
	. "eaciit/wfdemo-git/library/helper"
	. "eaciit/wfdemo-git/library/models"
	"eaciit/wfdemo-git/web/helper"

	"time"

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

	result = append(result, getAvailCollection(project, turbine, "SCADA_DATA_OEM"))
	result = append(result, getAvailCollection(project, turbine, "SCADA_DATA_HFD"))
	result = append(result, getAvailCollection(project, turbine, "MET_TOWER"))

	for {
		months = append(months, from.Format("Jan"))
		if from.Format("0601") == to.Format("0601") {
			break
		}
		from = GetNormalAddDateMonth(from.UTC(), 1)
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
		return helper.CreateResult(false, nil, e.Error())
	}

	e = csr.Fetch(&list, 0, false)

	defer csr.Close()

	res := []tk.M{}
	name := ""

	if len(list) > 0 {
		datas := []tk.M{}

		for _, dt := range list {
			t := dt.GetString("turbine")
			p := dt.GetString("project")
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

		if collType != "MET_TOWER" {
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
		"tooltip": start.Format("2 Jan 2006") + " until " + end.Format("2 Jan 2006"),
		"class":   class,
		"value":   tk.ToString(percentage) + "%",
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
