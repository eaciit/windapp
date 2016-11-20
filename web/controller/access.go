package controller

import (
	. "eaciit/ostrowfm/library/models"
	"eaciit/ostrowfm/web/helper"

	"github.com/eaciit/knot/knot.v1"
	"github.com/eaciit/toolkit"
)

type AccessController struct {
	App
}

func CreateAccessController() *AccessController {
	var controller = new(AccessController)
	return controller
}

func (a *AccessController) GetAccess(r *knot.WebContext) interface{} {
	r.Config.OutputType = knot.OutputJson

	payload := toolkit.M{}
	if err := r.GetPayload(&payload); err != nil {
		return helper.CreateResult(false, nil, err.Error())
	}

	data, err := GetAccessQuery(payload)
	if err != nil {
		return helper.CreateResult(false, nil, err.Error())
	}

	return helper.CreateResult(true, data, "success")

}

func (a *AccessController) GetParentID(r *knot.WebContext) interface{} {
	r.Config.OutputType = knot.OutputJson
	data, err := GetParentIDQuery()
	if err != nil {
		return helper.CreateResult(false, nil, err.Error())
	}

	return helper.CreateResult(true, data, "")
}

func (a *AccessController) Getaccessdropdown(r *knot.WebContext) interface{} {
	r.Config.OutputType = knot.OutputJson
	data, err := GetaccessDropDownQuery()
	if err != nil {
		return helper.CreateResult(false, nil, err.Error())
	}

	return helper.CreateResult(true, data, "")
}
func (a *AccessController) EditAccess(r *knot.WebContext) interface{} {
	r.Config.OutputType = knot.OutputJson

	payload := toolkit.M{}
	if err := r.GetPayload(&payload); err != nil {
		return helper.CreateResult(false, nil, err.Error())
	}
	data, err := FindAccess(payload)
	if err != nil {
		return helper.CreateResult(false, nil, err.Error())
	}

	return helper.CreateResult(true, data["tAccess"], "success")
}

func (a *AccessController) DeleteAccess(r *knot.WebContext) interface{} {
	r.Config.OutputType = knot.OutputJson

	payload := toolkit.M{}
	if err := r.GetPayload(&payload); err != nil {
		return helper.CreateResult(false, nil, err.Error())
	}

	err := DeleteAccessProc(payload)
	if err != nil {
		return helper.CreateResult(false, nil, err.Error())
	}

	return helper.CreateResult(true, nil, "delete access success")
}

func (a *AccessController) SaveAccess(r *knot.WebContext) interface{} {
	r.Config.OutputType = knot.OutputJson

	payload := toolkit.M{}
	if err := r.GetPayload(&payload); err != nil {
		return helper.CreateResult(false, nil, err.Error())
	}

	if err := SaveAccessProc(payload); err != nil {
		return helper.CreateResult(false, nil, err.Error())
	}

	return helper.CreateResult(true, nil, "save access success")
}
