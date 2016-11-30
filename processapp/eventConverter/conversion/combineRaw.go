package conversion

import (
	"time"

	// . "eaciit/wfdemo-git/library/helper"
	. "eaciit/wfdemo-git/library/models"

	_ "github.com/eaciit/dbox/dbc/mongo"
	"github.com/eaciit/orm"

	"github.com/eaciit/dbox"
	tk "github.com/eaciit/toolkit"
)

type CombineRaw struct {
	Ctx *orm.DataContext
}

func NewCombineRaw(ctx *orm.DataContext) *CombineRaw {
	ev := new(CombineRaw)
	ev.Ctx = ctx

	return ev
}

func (ev *CombineRaw) Run() {
	/*existing, compare := ev.getDatas("EventRaw_", "EventRaw", time.Now(), time.Now())

	for _, comp := range compare {
		for _, exist := range existing {

		}
	}*/
}

func (ev *CombineRaw) getRaw(coll string, start time.Time, end time.Time) (result map[string]EventRaw) {
	filter := []*dbox.Filter{}
	filter = append(filter, dbox.Gte("timestamp", start))
	filter = append(filter, dbox.Lte("timestamp", end))

	csr, err := ev.Ctx.Connection.NewQuery().
		From(coll).
		Where(filter...).
		Cursor(nil)

	defer csr.Close()

	eventRaws := []EventRaw{}

	if err != nil {
		tk.Println("Error: " + err.Error())
		return nil
	}
	err = csr.Fetch(&eventRaws, 0, false)

	if err != nil {
		tk.Println("Error: " + err.Error())
		return nil
	}

	for _, val := range eventRaws {
		key := ""

		result[key] = val
	}
	csr.Close()
	return
}

func (ev *CombineRaw) getDatas(existingColl string, compareColl string, start time.Time, end time.Time) (existing map[string]EventRaw, compare map[string]EventRaw) {
	existing = ev.getRaw(existingColl, start, end)
	compare = ev.getRaw(compareColl, start, end)
	return
}
