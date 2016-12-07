package controller

import (
	// . "eaciit/wfdemo-git/library/core"
	// . "eaciit/wfdemo-git/library/models"
	"eaciit/wfdemo-git/web/helper"

	"github.com/eaciit/knot/knot.v1"
)

type SCMManagementController struct {
	App
}

func CreateSCMManagementController() *SCMManagementController {
	var controller = new(SCMManagementController)
	return controller
}

func (m *SCMManagementController) GetData(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson
	return helper.CreateResult(true, nil, "success")
}
