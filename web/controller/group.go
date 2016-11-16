package controller

import (
	. "github.com/eaciit/windapp/library/models"
	"github.com/eaciit/windapp/web/helper"

	"github.com/eaciit/knot/knot.v1"
	"github.com/eaciit/toolkit"
)

type GroupController struct {
	App
}

func CreateGroupController() *GroupController {
	var controller = new(GroupController)
	return controller
}

func (a *GroupController) GetGroup(r *knot.WebContext) interface{} {
	r.Config.OutputType = knot.OutputJson
	data, err := GetGroup()
	if err != nil {
		return helper.CreateResult(false, nil, err.Error())
	}

	return helper.CreateResult(true, data, "success")

}
func (a *GroupController) EditGroup(r *knot.WebContext) interface{} {
	r.Config.OutputType = knot.OutputJson

	payload := toolkit.M{}

	if err := r.GetPayload(&payload); err != nil {
		return helper.CreateResult(false, nil, err.Error())
	}

	data, err := EditGroup(payload)
	if err != nil {
		return helper.CreateResult(false, nil, err.Error())
	}

	return helper.CreateResult(true, data["tGroup"], "success")

}

func (a *GroupController) Search(r *knot.WebContext) interface{} {
	r.Config.OutputType = knot.OutputJson

	payload := toolkit.M{}

	if err := r.GetPayload(&payload); err != nil {
		return helper.CreateResult(false, nil, err.Error())
	}

	data, err := SearchGroup(payload)
	if err != nil {
		return helper.CreateResult(false, nil, err.Error())
	}

	return helper.CreateResult(true, data, "success")
}

func (a *GroupController) GetAccessGroup(r *knot.WebContext) interface{} {
	r.Config.OutputType = knot.OutputJson

	payload := toolkit.M{}
	if err := r.GetPayload(&payload); err != nil {
		return helper.CreateResult(false, nil, err.Error())
	}

	data, err := GetAccessGroup(payload)
	if err != nil {
		return helper.CreateResult(false, nil, err.Error())
	}

	return helper.CreateResult(true, data, "success")
}

func (a *GroupController) DeleteGroup(r *knot.WebContext) interface{} {
	r.Config.OutputType = knot.OutputJson

	payload := toolkit.M{}
	if err := r.GetPayload(&payload); err != nil {
		return helper.CreateResult(false, nil, err.Error())
	}

	if err := DeleteGroup(payload); err != nil {
		return helper.CreateResult(false, nil, err.Error())
	}

	return helper.CreateResult(true, nil, "Delete Group Success")
}

func (a *GroupController) SaveGroup(r *knot.WebContext) interface{} {
	r.Config.OutputType = knot.OutputJson

	payload := toolkit.M{}
	if err := r.GetPayload(&payload); err != nil {
		return helper.CreateResult(false, nil, err.Error())
	}

	if err := SaveGroup(payload); err != nil {
		return helper.CreateResult(false, nil, err.Error())
	}

	return helper.CreateResult(true, nil, "Save Group Success")
}
