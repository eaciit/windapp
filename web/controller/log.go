package controller

import (
	. "eaciit/wfdemo-git-dev/library/models"
	"eaciit/wfdemo-git-dev/web/helper"

	"github.com/eaciit/knot/knot.v1"
	"github.com/eaciit/toolkit"
)

type LogController struct {
	App
}

func CreateLogController() *LogController {
	var controller = new(LogController)
	return controller
}

func (a *LogController) GetLog(r *knot.WebContext) interface{} {
	r.Config.OutputType = knot.OutputJson

	payload := toolkit.M{}
	if err := r.GetForms(&payload); err != nil {
		return helper.CreateResult(false, nil, err.Error())
	}

	data, err := GetLog(payload)
	if err != nil {
		return helper.CreateResult(false, nil, err.Error())
	}

	return helper.CreateResult(true, data, "success")

}
