package controller

import (
	. "eaciit/wfdemo-git/library/core"
	. "eaciit/wfdemo-git/library/models"
	"eaciit/wfdemo-git/web/helper"
	// "time"

	"github.com/eaciit/knot/knot.v1"
)

type MonitoringController struct {
	App
}

func CreateMonitoringController() *MonitoringController {
	var controller = new(MonitoringController)
	return controller
}

func (m *MonitoringController) GetData(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	p := new(helper.Payloads)
	e := k.GetPayload(&p)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	filter := p.ParseFilter()

	if e != nil {
		helper.CreateResult(false, nil, e.Error())
	}

	return helper.CreateResult(true, data, "success")
}
