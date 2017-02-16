package generatorControllers

import (
	. "eaciit/wfdemo-git/library/helper"
	. "eaciit/wfdemo-git/library/models"
	. "eaciit/wfdemo-git/processapp/controllers"
	"github.com/eaciit/crowd"
	"github.com/eaciit/dbox"
	_ "github.com/eaciit/dbox/dbc/mongo"
	tk "github.com/eaciit/toolkit"
	"os"
	"time"
)

type UpdateMetTower struct {
	*BaseController
}

func (d *UpdateMetTower) Generate(base *BaseController) {
	funcName := "UpdateMetTower Data"
	if base != nil {
		d.BaseController = base

		ctx, e := PrepareConnection()
		if e != nil {
			ErrorHandler(e, funcName)
			os.Exit(0)
		}

		metTowers := []MetTower{}

		csr, e := ctx.NewQuery().From(new(MetTower).TableName()).Cursor(nil)

		e = csr.Fetch(&metTowers, 0, false)
		ErrorHandler(e, funcName)
		csr.Close()
		tk.Println(funcName)
		count := 0
		total := 0
		for _, data := range metTowers {
			windSpeed := data.VHubWS90mAvg
			windDir := data.DHubWD88mAvg

			wsCatNo, wsCatDesc := GetWsCategory(windSpeed)
			wdNo, wdDesc := GetDirectionOnlyFromWD(windDir)
			windSpeedBin := tk.RoundingUp64(windSpeed, 0)

			e = ctx.NewQuery().Update().From(new(MetTower).TableName()).
				Where(dbox.Eq("_id", data.ID)).
				Exec(tk.M{}.Set("data", tk.M{}.Set("winddirno", wdNo).Set("winddirdesc", wdDesc).Set("wscategoryno", wsCatNo).Set("wscategorydesc", wsCatDesc).Set("windspeedbin", windSpeedBin)))
			if e != nil {
				tk.Printf("Update fail: %s", e.Error())
			}

			count++
			total++
			if count == 1000 {
				tk.Printf("Total processed data %v\n", total)
				count = 0
			}
		}
		tk.Printf("Total processed data %v\n", total)
	}
}

func (d *UpdateMetTower) GenerateWindRose(base *BaseController) {
	funcName := "GenerateWindRose from Met Tower Data"
	if base != nil {
		d.BaseController = base

		conn := d.BaseController.Ctx.Connection
		defer conn.Close()

		csrt, e := conn.NewQuery().From(new(MetTower).TableName()).
			Group("dateinfo.dateid").
			Cursor(nil)

		if e != nil {
			ErrorHandler(e, funcName)
			os.Exit(0)
		}

		parents := []tk.M{}
		e = csrt.Fetch(&parents, 0, false)
		csrt.Close()

		if len(parents) > 0 {
			for _, p := range parents {
				tId := p["_id"].(tk.M)
				tDateId := tId["dateinfo_dateid"].(time.Time)

				mdl := new(WindRoseMTModel).New()
				mdl.ProjectId = "Tejuva"
				mdl.DateInfo = GetDateInfo(tDateId.UTC())

				pipe := []tk.M{tk.M{}.Set("$match", tk.M{}.Set("dateinfo.dateid", tDateId.UTC())), tk.M{}.Set("$group", tk.M{}.Set("_id", tk.M{}.Set("dateinfo", "$dateinfo").Set("winddirno", "$winddirno").Set("winddirdesc", "$winddirdesc").Set("wscategoryno", "$wscategoryno").Set("wscategorydesc", "$wscategorydesc")).
					Set("duration", tk.M{}.Set("$sum", "$vhubws90mcount")).Set("totalcount", tk.M{}.Set("$sum", 1)))}
				csr, _ := conn.NewQuery().
					Command("pipe", pipe).
					From(new(MetTower).TableName()).
					Cursor(nil)

				datas := []tk.M{}
				_ = csr.Fetch(&datas, 0, false)
				csr.Close()

				dataTotal := crowd.From(&datas).Sum(func(x interface{}) interface{} {
					dt := x.(tk.M)
					return dt["duration"].(float64)
				}).Exec().Result.Sum

				writems := make([]WindRoseItemMT, 0)
				if len(datas) > 0 {
					for _, d := range datas {
						ids := d["_id"].(tk.M)
						dirNo := int(ids["winddirno"].(float64))
						dirDesc := ids["winddirdesc"].(string)
						wsNo := int(ids["wscategoryno"].(float64))
						wsDesc := ids["wscategorydesc"].(string)
						duration := d["duration"].(float64)
						contribution := tk.Div(duration, dataTotal)
						frequency := d["totalcount"].(int)

						var item WindRoseItemMT
						item.DirectionNo = dirNo
						item.DirectionDesc = dirDesc
						item.WSCategoryNo = wsNo
						item.WSCategoryDesc = wsDesc
						item.Hours = tk.Div(duration, 3600.0)
						item.Frequency = frequency
						item.Contribute = contribution

						writems = append(writems, item)
					}
				}
				mdl.WindRoseItems = writems

				pipe1 := []tk.M{tk.M{}.Set("$match", tk.M{}.Set("dateinfo.dateid", tDateId.UTC())), tk.M{}.Set("$group", tk.M{}.Set("_id", tk.M{}.Set("dateinfo", "$dateinfo").Set("wscategoryno", "$wscategoryno").Set("wscategorydesc", "$wscategorydesc")).
					Set("duration", tk.M{}.Set("$sum", "$vhubws90mcount")).Set("totalcount", tk.M{}.Set("$sum", 1)))}
				csr1, _ := conn.NewQuery().
					Command("pipe", pipe1).
					From(new(MetTower).TableName()).
					Cursor(nil)

				datas1 := []tk.M{}
				_ = csr1.Fetch(&datas1, 0, false)
				csr1.Close()

				dataTotalContribute := crowd.From(&datas1).Sum(func(x interface{}) interface{} {
					dt := x.(tk.M)
					return dt["duration"].(float64)
				}).Exec().Result.Sum

				totalconts := make([]WindRoseContributeMT, 0)
				if len(datas1) > 0 {
					for _, d := range datas1 {
						ids := d["_id"].(tk.M)
						wsNo := int(ids["wscategoryno"].(float64))
						wsDesc := ids["wscategorydesc"].(string)
						duration := d["duration"].(float64)
						contribution := tk.Div(duration, dataTotalContribute)
						frequency := d["totalcount"].(int)

						var item WindRoseContributeMT
						item.WSCategoryNo = wsNo
						item.WSCategoryDesc = wsDesc
						item.Hours = tk.Div(duration, 3600.0)
						item.Contribute = contribution
						item.Frequency = frequency

						totalconts = append(totalconts, item)
					}
				}
				mdl.TotalContributes = totalconts

				e := d.BaseController.Ctx.Insert(mdl)
				if e != nil {
					ErrorHandler(e, funcName)
					os.Exit(0)
				}
			}
		}
	}
}
