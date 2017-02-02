package checker

import (
	"log"
	"time"

	"gopkg.in/mgo.v2/bson"

	// . "eaciit/wfdemo-git/library/helper"
	. "eaciit/wfdemo-git/library/models"

	_ "github.com/eaciit/dbox/dbc/mongo"
	"github.com/eaciit/orm"

	tk "github.com/eaciit/toolkit"
)

type OEMChecker struct {
	Ctx *orm.DataContext
}

func NewOEMChecker(ctx *orm.DataContext) *OEMChecker {
	ev := new(OEMChecker)
	ev.Ctx = ctx

	return ev
}

func (ev *OEMChecker) Run() {
	max, min := ev.getMaxMin()
	current := min

	log.Printf("%v | %v \n", min.String(), max.String())

	for {
		list := ev.getDatas(current)
		var ids []bson.ObjectId

		log.Println(len(list))

		for idx, val := range list {
			log.Println(idx)
			if idx != 0 {
				before := list[idx-1]
				now := val

				if before.TimeStamp.Format("20060102_150405") == now.TimeStamp.Format("20060102_150405") {
					ids = append(ids, before.ID)
				}
			}
		}

		// ev.Ctx.DeleteMany(new(ScadaDataOEM), dbox.And(dbox.In("_id", ids)))
		log.Printf("%v | %v \n", current.Format("2006-01-02"), len(ids))
		if current.Format("20060102") == max.Format("20060102") {
			break
		} else {
			current = current.AddDate(0, 0, 1)
		}
	}
}

func (ev *OEMChecker) getDatas(dateid time.Time) (result []ScadaDataOEM) {
	pipes := make([]tk.M, 0)

	pipes = append(pipes, tk.M{
		"$match": tk.M{}.Set("dateinfo.dateid", dateid),
	})
	pipes = append(pipes, tk.M{
		"$sort": tk.M{}.Set("timestamp", 1),
	})

	csr, _ := ev.Ctx.Connection.NewQuery().
		Command("pipe", pipes).
		From(new(ScadaDataOEM).TableName()).
		Cursor(nil)

	csr.Fetch(&result, 0, false)
	csr.Close()

	return
}

func (ev *OEMChecker) getMaxMin() (max time.Time, min time.Time) {
	pipes := make([]tk.M, 0)
	pipes = append(pipes, tk.M{
		"$group": tk.M{}.Set("_id", "oem").
			Set("min", tk.M{}.Set("$min", "$timestamp")).
			Set("max", tk.M{}.Set("$max", "$timestamp")),
	})

	csr, _ := ev.Ctx.Connection.NewQuery().
		Command("pipe", pipes).
		From(new(ScadaDataOEM).TableName()).
		Cursor(nil)

	res := []tk.M{}
	csr.Fetch(&res, 0, false)

	csr.Close()

	if len(res) > 0 {
		min = res[0].Get("min").(time.Time).UTC()
		max = res[0].Get("max").(time.Time).UTC()
	}

	return
}
