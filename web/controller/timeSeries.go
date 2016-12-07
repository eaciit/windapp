package controller

import (
	// . "eaciit/wfdemo-git/library/core"
	// . "eaciit/wfdemo-git/library/models"
	"eaciit/wfdemo-git/web/helper"

	"github.com/eaciit/knot/knot.v1"
)

type TimeSeriesController struct {
	App
}

func CreateTimeSeriesController() *TimeSeriesController {
	var controller = new(TimeSeriesController)
	return controller
}

func (m *TimeSeriesController) GetData(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson
	return helper.CreateResult(true, nil, "success")
}
