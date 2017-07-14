package controller

import (
	// . "eaciit/wfdemo-git/library/core"
	// . "eaciit/wfdemo-git/library/models"
	"eaciit/wfdemo-git/web/helper"

	"github.com/eaciit/knot/knot.v1"
)

type MonitoringCustomController struct {
	App
}

func CreateMonitoringCustomController() *MonitoringCustomController {
	var controller = new(MonitoringCustomController)
	return controller
}

func (m *MonitoringCustomController) GetData(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson
	return helper.CreateResult(true, nil, "success")
}
