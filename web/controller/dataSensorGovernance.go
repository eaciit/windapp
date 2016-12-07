package controller

import (
	// . "eaciit/wfdemo-git/library/core"
	// . "eaciit/wfdemo-git/library/models"
	"eaciit/wfdemo-git/web/helper"

	"github.com/eaciit/knot/knot.v1"
)

type DataSensorGovernanceController struct {
	App
}

func CreateDataSensorGovernanceController() *DataSensorGovernanceController {
	var controller = new(DataSensorGovernanceController)
	return controller
}

func (m *DataSensorGovernanceController) GetData(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson
	return helper.CreateResult(true, nil, "success")
}
