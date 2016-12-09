package controller

import (
	// . "eaciit/wfdemo-git/library/core"
	// . "eaciit/wfdemo-git/library/models"
	"eaciit/wfdemo-git/web/helper"

	"github.com/eaciit/knot/knot.v1"
)

type IssueTrackingController struct {
	App
}

func CreateIssueTrackingController() *IssueTrackingController {
	var controller = new(IssueTrackingController)
	return controller
}

func (m *IssueTrackingController) GetData(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson
	return helper.CreateResult(true, nil, "success")
}
