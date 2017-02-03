package checker

import (
	"log"
	"sync"
	"time"

	// . "eaciit/wfdemo-git/library/helper"
	. "eaciit/wfdemo-git/library/models"

	"github.com/eaciit/dbox"
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
	var wg sync.WaitGroup
	turbines := ev.getMaxMin()

	wg.Add(len(turbines))
	for _, val := range turbines {
		go func(val tk.M) {
			turbine := val.GetString("_id")
			// if turbine == "HBR007" {
			max := val.Get("max").(time.Time).UTC()
			min := val.Get("min").(time.Time).UTC()
			current, _ := time.Parse("20060102_150405", min.Format("20060102_")+"000000")

			log.Printf("%v | %v | %v \n", turbine, min.String(), max.String())

			// maxDays := max.Sub(min).Hours() / 24
			counter := 0

			for {
				// log.Printf("%v-%v => %v | %v of %v \n", turbine, current.Format("20060102_150405"), float64(counter)/maxDays*100, counter, maxDays)

				list := ev.getDatas(current, turbine)
				var ids []interface{}
				for idx, val := range list {
					if idx != 0 {
						before := list[idx-1]
						now := val

						if before.TimeStamp.Format("20060102_150405") == now.TimeStamp.Format("20060102_150405") {
							ids = append(ids, before.ID)
						}
					}
				}

				if len(ids) > 0 {
					ev.Ctx.DeleteMany(new(ScadaDataOEM), dbox.And(dbox.In("_id", ids...)))
					log.Printf(">>>>>>>> %v - %v | %v \n", turbine, current.Format("2006-01-02"), len(ids))
				}

				// log.Printf("%v | %v \n", current.Format("2006-01-02"), len(ids))
				if current.Format("20060102") == max.Format("20060102") {
					break
				} else {
					current = current.AddDate(0, 0, 1)
					counter++
				}
			}
			// }

			wg.Done()
		}(val)
	}

	wg.Wait()
}

func (ev *OEMChecker) getDatas(dateid time.Time, turbine string) (result []ScadaDataOEM) {
	pipes := make([]tk.M, 0)

	pipes = append(pipes, tk.M{
		"$match": tk.M{}.Set("dateinfo.dateid", dateid).Set("turbine", turbine),
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

func (ev *OEMChecker) getMaxMin() (res []tk.M) {
	pipes := make([]tk.M, 0)
	pipes = append(pipes, tk.M{
		"$group": tk.M{}.Set("_id", "$turbine").
			Set("min", tk.M{}.Set("$min", "$timestamp")).
			Set("max", tk.M{}.Set("$max", "$timestamp")),
	})

	csr, _ := ev.Ctx.Connection.NewQuery().
		Command("pipe", pipes).
		From(new(ScadaDataOEM).TableName()).
		Cursor(nil)

	csr.Fetch(&res, 0, false)

	csr.Close()

	return
}
