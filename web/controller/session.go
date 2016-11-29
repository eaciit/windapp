package controller

import (
	. "eaciit/wfdemo-git-dev/library/models"
	"eaciit/wfdemo-git-dev/web/helper"
	"github.com/eaciit/knot/knot.v1"
	"github.com/eaciit/toolkit"
)

type SessionController struct {
	App
}

func CreateSessionController() *SessionController {
	var controller = new(SessionController)
	return controller
}

func (a *SessionController) GetSession(r *knot.WebContext) interface{} {
	r.Config.OutputType = knot.OutputJson

	payload := toolkit.M{}
	if err := r.GetPayload(&payload); err != nil {
		return helper.CreateResult(false, nil, err.Error())
	}

	data, err := GetSession(payload)
	if err != nil {
		return helper.CreateResult(false, nil, err.Error())
	}

	return helper.CreateResult(true, data, "success")

}

func (a *SessionController) SetExpired(r *knot.WebContext) interface{} {
	r.Config.OutputType = knot.OutputJson

	payload := toolkit.M{}
	if err := r.GetPayload(&payload); err != nil {
		return helper.CreateResult(false, nil, err.Error())
	}

	err := SetExpired(payload)
	if err != nil {
		return helper.CreateResult(false, nil, err.Error())
	}

	return helper.CreateResult(true, nil, "Set Expired Success")
}
