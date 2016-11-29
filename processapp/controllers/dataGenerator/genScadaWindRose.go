package generatorControllers

import (
	. "eaciit/wfdemo-git-dev/library/helper"
	. "eaciit/wfdemo-git-dev/library/models"
	. "eaciit/wfdemo-git-dev/processapp/controllers"
	"fmt"
	c "github.com/eaciit/crowd"
	"github.com/eaciit/dbox"
	_ "github.com/eaciit/dbox/dbc/mongo"
	tk "github.com/eaciit/toolkit"
	"math"
	"os"
	_ "strconv"
	_ "strings"
	"time"
)

type GenScadaWindRose struct {
	*BaseController
}

func (d *GenScadaWindRose) Generate(base *BaseController) {
	if base != nil {
		d.BaseController = base

		timeNow := time.Now()

		ctx, e := PrepareConnection()
		if e != nil {
			ErrorHandler(e, "Scada Summary")
			os.Exit(0)
		}

		tk.Println("Starting getting scada data")

		dateStart, _ := time.Parse("2006-01-02", "2016-01-01")
		dateEnd, _ := time.Parse("2006-01-02", "2016-06-30")

		csr, e := ctx.NewQuery().From(new(ScadaData).TableName()).
			//Where(dbox.And(dbox.Gte("dateinfo.dateid", dateStart), dbox.Lte("dateinfo.dateid", dateEnd), dbox.Eq("turbine", "SSE002"))).
			Where(dbox.And(dbox.Gte("dateinfo.dateid", dateStart), dbox.Lte("dateinfo.dateid", dateEnd), dbox.Gte("power", 0))).
			Cursor(nil)
		defer csr.Close()

		scadas := make([]ScadaData, 0)
		_ = csr.Fetch(&scadas, 0, false)

		tk.Println("Starting to processing data for wind rose")

		totalDuration := float64((len(scadas) * 10.0)) / 60.0

		type DataItem struct {
			ProjectName    string
			Turbine        string
			DateId         time.Time
			DirectionNo    int
			DirectionDesc  string
			WsCategoryNo   int
			WsCategoryDesc string
			Hours          float64
			Frequency      int
		}

		type DataItemResult struct {
			ProjectName    string
			Turbine        string
			DateId         time.Time
			DirectionNo    int
			DirectionDesc  string
			WsCategoryNo   int
			WsCategoryDesc string
			Hours          float64
			Contribution   float64
			Frequency      int
		}

		type DataItemGroup struct {
			ProjectName    string
			Turbine        string
			DateId         time.Time
			DirectionNo    int
			DirectionDesc  string
			WsCategoryNo   int
			WsCategoryDesc string
		}

		type DataGroupResult struct {
			ProjectName string
			Turbine     string
			DateId      time.Time
		}

		type ContributeGroupResult struct {
			WsCategoryNo   int
			WsCategoryDesc string
		}

		type ContributeItemResult struct {
			WsCategoryNo   int
			WsCategoryDesc string
			Hours          float64
			Contribution   float64
			Frequency      int
		}

		datas := c.From(&scadas).Apply(func(x interface{}) interface{} {
			dt := x.(ScadaData)
			var di DataItem

			dateId := dt.DateInfo.DateId
			project := "Tejuva"
			turbine := dt.Turbine
			dirNo, dirDesc := GetDirection(dt.WindDirection, dt.NacelDirection)
			wsNo, wsDesc := GetWsCategory(dt.AvgWindSpeed)

			di.DateId = dateId
			di.ProjectName = project
			di.Turbine = turbine
			di.DirectionNo = dirNo
			di.DirectionDesc = dirDesc
			di.WsCategoryNo = wsNo
			di.WsCategoryDesc = wsDesc
			di.Hours = 10.0 / 60.0
			di.Frequency = 1

			return di
		}).Exec().Group(func(x interface{}) interface{} {
			dt := x.(DataItem)

			var dig DataItemGroup
			dig.DateId = dt.DateId
			dig.ProjectName = dt.ProjectName
			dig.Turbine = dt.Turbine
			dig.DirectionNo = dt.DirectionNo
			dig.DirectionDesc = dt.DirectionDesc
			dig.WsCategoryNo = dt.WsCategoryNo
			dig.WsCategoryDesc = dt.WsCategoryDesc

			return dig
		}, nil).Exec()

		// fmt.Println(datas.Result.Data())

		dts := datas.Apply(func(x interface{}) interface{} {
			kv := x.(c.KV)
			vv := kv.Key.(DataItemGroup)
			vs := kv.Value.([]DataItem)

			sumDuration := c.From(&vs).Sum(func(x interface{}) interface{} {
				dt := x.(DataItem)
				return dt.Hours
			}).Exec().Result.Sum

			var di DataItemResult
			di.DateId = vv.DateId
			di.ProjectName = vv.ProjectName
			di.Turbine = vv.Turbine
			di.DirectionNo = vv.DirectionNo
			di.DirectionDesc = vv.DirectionDesc
			di.WsCategoryNo = vv.WsCategoryNo
			di.WsCategoryDesc = vv.WsCategoryDesc
			di.Hours = sumDuration
			contribute := 0.0
			contribute = sumDuration / totalDuration
			di.Contribution = RoundUp(contribute, .5, 2)

			sumFreq := c.From(&vs).Sum(func(x interface{}) interface{} {
				dt := x.(DataItem)
				return dt.Frequency
			}).Exec().Result.Sum
			di.Frequency = int(sumFreq)

			return di
		}).Exec().Group(func(x interface{}) interface{} {
			dt := x.(DataItemResult)

			var dg DataGroupResult
			dg.DateId = dt.DateId
			dg.ProjectName = dt.ProjectName
			dg.Turbine = dt.Turbine

			return dg
		}, nil).Exec()

		results := dts.Result.Data().([]c.KV)
		for _, res := range results {
			//tk.Println("Processing index-" + tk.ToString(index))

			keys := res.Key.(DataGroupResult)
			values := res.Value.([]DataItemResult)

			totalDurationPerKey := c.From(&values).Sum(func(x interface{}) interface{} {
				dt := x.(DataItemResult)
				return dt.Hours
			}).Exec().Result.Sum

			data := new(WindRoseModel).New()
			data.DateInfo = GetDateInfo(keys.DateId)
			data.ProjectId = keys.ProjectName
			data.TurbineId = keys.Turbine

			wrItems := make([]WindRoseItem, 0)
			for _, x := range values {
				var wri WindRoseItem
				wri.DirectionNo = x.DirectionNo
				wri.DirectionDesc = x.DirectionDesc
				wri.WSCategoryNo = x.WsCategoryNo
				wri.WSCategoryDesc = x.WsCategoryDesc
				wri.Hours = x.Hours
				contribute := x.Hours / totalDurationPerKey
				wri.Contribute = contribute
				wri.Frequency = x.Frequency

				wrItems = append(wrItems, wri)
			}

			data.WindRoseItems = wrItems

			groups := c.From(&values).Group(func(x interface{}) interface{} {
				dt := x.(DataItemResult)

				var cgr ContributeGroupResult
				cgr.WsCategoryDesc = dt.WsCategoryDesc
				cgr.WsCategoryNo = dt.WsCategoryNo

				return cgr
			}, nil).Exec().Apply(func(x interface{}) interface{} {
				kv := x.(c.KV)
				vv := kv.Key.(ContributeGroupResult)
				vs := kv.Value.([]DataItemResult)

				sumHours := c.From(&vs).Sum(func(x interface{}) interface{} {
					dt := x.(DataItemResult)
					return dt.Hours
				}).Exec().Result.Sum

				sumFreq := c.From(&vs).Sum(func(x interface{}) interface{} {
					dt := x.(DataItemResult)
					return dt.Frequency
				}).Exec().Result.Sum

				var cir ContributeItemResult
				cir.WsCategoryNo = vv.WsCategoryNo
				cir.WsCategoryDesc = vv.WsCategoryDesc
				cir.Hours = sumHours
				cir.Contribution = sumHours / totalDurationPerKey
				cir.Frequency = int(sumFreq)

				return cir
			}).Exec().Result.Data()

			contributes := make([]WindRoseContribute, 0)
			for _, dt := range groups.([]ContributeItemResult) {
				var wrc WindRoseContribute
				wrc.WSCategoryNo = dt.WsCategoryNo
				wrc.WSCategoryDesc = dt.WsCategoryDesc
				wrc.Hours = dt.Hours
				wrc.Contribute = dt.Contribution
				wrc.Frequency = dt.Frequency

				contributes = append(contributes, wrc)
			}

			data.TotalContributes = contributes

			d.BaseController.Ctx.Insert(data)
		}

		duration := timeNow.Sub(time.Now())
		fmt.Println(duration.Seconds())
	}
}

func (d *GenScadaWindRose) GenerateFromScadaNew(base *BaseController) {
	if base != nil {
		d.BaseController = base

		timeNow := time.Now()

		ctx, e := PrepareConnection()
		if e != nil {
			ErrorHandler(e, "Scada Summary")
			os.Exit(0)
		}

		tk.Println("Starting getting scada data")

		dateStart, _ := time.Parse("2006-01-02", "2016-07-01")
		dateEnd, _ := time.Parse("2006-01-02", "2016-08-27")

		csr, e := ctx.NewQuery().From(new(ScadaDataNew).TableName()).
			//Where(dbox.And(dbox.Gte("dateinfo.dateid", dateStart), dbox.Lte("dateinfo.dateid", dateEnd), dbox.Eq("turbine", "SSE002"))).
			Where(dbox.And(dbox.Gte("dateinfo.dateid", dateStart), dbox.Lte("dateinfo.dateid", dateEnd), dbox.Gte("power", 0))).
			Cursor(nil)
		defer csr.Close()

		scadas := make([]ScadaDataNew, 0)
		_ = csr.Fetch(&scadas, 0, false)

		tk.Printf("%v\n", scadas)

		tk.Println("Starting to processing data for wind rose")

		totalDuration := float64((len(scadas) * 10.0)) / 60.0

		type DataItem struct {
			ProjectName    string
			Turbine        string
			DateId         time.Time
			DirectionNo    int
			DirectionDesc  string
			WsCategoryNo   int
			WsCategoryDesc string
			Hours          float64
			Frequency      int
		}

		type DataItemResult struct {
			ProjectName    string
			Turbine        string
			DateId         time.Time
			DirectionNo    int
			DirectionDesc  string
			WsCategoryNo   int
			WsCategoryDesc string
			Hours          float64
			Contribution   float64
			Frequency      int
		}

		type DataItemGroup struct {
			ProjectName    string
			Turbine        string
			DateId         time.Time
			DirectionNo    int
			DirectionDesc  string
			WsCategoryNo   int
			WsCategoryDesc string
		}

		type DataGroupResult struct {
			ProjectName string
			Turbine     string
			DateId      time.Time
		}

		type ContributeGroupResult struct {
			WsCategoryNo   int
			WsCategoryDesc string
		}

		type ContributeItemResult struct {
			WsCategoryNo   int
			WsCategoryDesc string
			Hours          float64
			Contribution   float64
			Frequency      int
		}

		datas := c.From(&scadas).Apply(func(x interface{}) interface{} {
			dt := x.(ScadaDataNew)
			var di DataItem

			dateId := dt.DateInfo.DateId
			project := "Tejuva"
			turbine := dt.Turbine
			dirNo, dirDesc := GetDirectionOnlyFromWD(dt.NacelDirection)
			wsNo, wsDesc := GetWsCategory(dt.AvgWindSpeed)

			di.DateId = dateId
			di.ProjectName = project
			di.Turbine = turbine
			di.DirectionNo = dirNo
			di.DirectionDesc = dirDesc
			di.WsCategoryNo = wsNo
			di.WsCategoryDesc = wsDesc
			di.Hours = 10.0 / 60.0
			di.Frequency = 1

			return di
		}).Exec().Group(func(x interface{}) interface{} {
			dt := x.(DataItem)

			var dig DataItemGroup
			dig.DateId = dt.DateId
			dig.ProjectName = dt.ProjectName
			dig.Turbine = dt.Turbine
			dig.DirectionNo = dt.DirectionNo
			dig.DirectionDesc = dt.DirectionDesc
			dig.WsCategoryNo = dt.WsCategoryNo
			dig.WsCategoryDesc = dt.WsCategoryDesc

			return dig
		}, nil).Exec()

		// fmt.Println(datas.Result.Data())

		dts := datas.Apply(func(x interface{}) interface{} {
			kv := x.(c.KV)
			vv := kv.Key.(DataItemGroup)
			vs := kv.Value.([]DataItem)

			sumDuration := c.From(&vs).Sum(func(x interface{}) interface{} {
				dt := x.(DataItem)
				return dt.Hours
			}).Exec().Result.Sum

			var di DataItemResult
			di.DateId = vv.DateId
			di.ProjectName = vv.ProjectName
			di.Turbine = vv.Turbine
			di.DirectionNo = vv.DirectionNo
			di.DirectionDesc = vv.DirectionDesc
			di.WsCategoryNo = vv.WsCategoryNo
			di.WsCategoryDesc = vv.WsCategoryDesc
			di.Hours = sumDuration
			contribute := 0.0
			contribute = sumDuration / totalDuration
			di.Contribution = RoundUp(contribute, .5, 2)

			sumFreq := c.From(&vs).Sum(func(x interface{}) interface{} {
				dt := x.(DataItem)
				return dt.Frequency
			}).Exec().Result.Sum
			di.Frequency = int(sumFreq)

			return di
		}).Exec().Group(func(x interface{}) interface{} {
			dt := x.(DataItemResult)

			var dg DataGroupResult
			dg.DateId = dt.DateId
			dg.ProjectName = dt.ProjectName
			dg.Turbine = dt.Turbine

			return dg
		}, nil).Exec()

		results := dts.Result.Data().([]c.KV)
		for _, res := range results {
			//tk.Println("Processing index-" + tk.ToString(index))

			keys := res.Key.(DataGroupResult)
			values := res.Value.([]DataItemResult)

			totalDurationPerKey := c.From(&values).Sum(func(x interface{}) interface{} {
				dt := x.(DataItemResult)
				return dt.Hours
			}).Exec().Result.Sum

			data := new(WindRoseNewModel).New()
			data.DateInfo = GetDateInfo(keys.DateId)
			data.ProjectId = keys.ProjectName
			data.TurbineId = keys.Turbine

			wrItems := make([]WindRoseItemNew, 0)
			for _, x := range values {
				var wri WindRoseItemNew
				wri.DirectionNo = x.DirectionNo
				wri.DirectionDesc = x.DirectionDesc
				wri.WSCategoryNo = x.WsCategoryNo
				wri.WSCategoryDesc = x.WsCategoryDesc
				wri.Hours = x.Hours
				contribute := x.Hours / totalDurationPerKey
				wri.Contribute = contribute
				wri.Frequency = x.Frequency

				wrItems = append(wrItems, wri)
			}

			data.WindRoseItems = wrItems

			groups := c.From(&values).Group(func(x interface{}) interface{} {
				dt := x.(DataItemResult)

				var cgr ContributeGroupResult
				cgr.WsCategoryDesc = dt.WsCategoryDesc
				cgr.WsCategoryNo = dt.WsCategoryNo

				return cgr
			}, nil).Exec().Apply(func(x interface{}) interface{} {
				kv := x.(c.KV)
				vv := kv.Key.(ContributeGroupResult)
				vs := kv.Value.([]DataItemResult)

				sumHours := c.From(&vs).Sum(func(x interface{}) interface{} {
					dt := x.(DataItemResult)
					return dt.Hours
				}).Exec().Result.Sum

				sumFreq := c.From(&vs).Sum(func(x interface{}) interface{} {
					dt := x.(DataItemResult)
					return dt.Frequency
				}).Exec().Result.Sum

				var cir ContributeItemResult
				cir.WsCategoryNo = vv.WsCategoryNo
				cir.WsCategoryDesc = vv.WsCategoryDesc
				cir.Hours = sumHours
				cir.Contribution = sumHours / totalDurationPerKey
				cir.Frequency = int(sumFreq)

				return cir
			}).Exec().Result.Data()

			contributes := make([]WindRoseContributeNew, 0)
			for _, dt := range groups.([]ContributeItemResult) {
				var wrc WindRoseContributeNew
				wrc.WSCategoryNo = dt.WsCategoryNo
				wrc.WSCategoryDesc = dt.WsCategoryDesc
				wrc.Hours = dt.Hours
				wrc.Contribute = dt.Contribution
				wrc.Frequency = dt.Frequency

				contributes = append(contributes, wrc)
			}

			data.TotalContributes = contributes

			d.BaseController.Ctx.Insert(data)
		}

		duration := timeNow.Sub(time.Now())
		fmt.Println(duration.Seconds())
	}
}

func GetWsCategory(ws float64) (int, string) {
	catNo := 0
	catDesc := "0 to 4m/s"
	if ws >= 14 {
		catNo = 4
		catDesc = "14 and above"
	} else if ws >= 9 {
		catNo = 3
		catDesc = "9 to 14m/s"
	} else if ws >= 7 {
		catNo = 2
		catDesc = "7 to 9m/s"
	} else if ws >= 4 {
		catNo = 1
		catDesc = "4 to 7m/s"
	}

	return catNo, catDesc
}

func GetDirection(windDir float64, nacelPos float64) (int, string) {
	dirNo := 0
	dirDescs := []string{"N", "NE", "E", "SE", "S", "SW", "W", "NW"}
	if windDir < 0 {
		windDir = 360.0 + windDir
	}
	if nacelPos < 0 {
		nacelPos = 360.0 + nacelPos
	}
	if nacelPos < 0 {
		nacelPos = 360.0 + nacelPos
	}
	dirCalc := math.Mod((nacelPos + windDir), 360.0)
	dirNo = int(tk.RoundingAuto64(dirCalc/45.0, 0))

	if dirNo > 7 {
		dirNo = 0
	}

	return dirNo, dirDescs[dirNo]
}

func GetDirectionOnlyFromWD(windDir float64) (int, string) {
	dirNo := 0
	dirDescs := []string{"N", "NE", "E", "SE", "S", "SW", "W", "NW"}
	if windDir < 0 {
		windDir = 360.0 + windDir
	}
	dirNo = int(tk.RoundingAuto64(windDir/45.0, 0))

	if dirNo > 7 {
		dirNo = 0
	}

	return dirNo, dirDescs[dirNo]
}
