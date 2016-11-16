package controller

import (
	. "eaciit/wfdemo/library/models"
	"eaciit/wfdemo/web/helper"
	"github.com/eaciit/knot/knot.v1"
	"github.com/eaciit/toolkit"
)

type UserController struct {
	App
}

func CreateUserController() *UserController {
	var controller = new(UserController)
	return controller
}

func (a *UserController) GetUser(r *knot.WebContext) interface{} {
	r.Config.OutputType = knot.OutputJson

	payload := toolkit.M{}
	if err := r.GetPayload(&payload); err != nil {
		return helper.CreateResult(false, nil, err.Error())
	}

	data, err := GetUser(payload)
	if err != nil {
		return helper.CreateResult(false, nil, err.Error())
	}

	return helper.CreateResult(true, data, "Success")
}
func (a *UserController) EditUser(r *knot.WebContext) interface{} {
	r.Config.OutputType = knot.OutputJson

	payload := toolkit.M{}
	if err := r.GetPayload(&payload); err != nil {
		return helper.CreateResult(false, nil, err.Error())
	}
	data, err := EditUser(payload)
	if err != nil {
		return helper.CreateResult(false, nil, err.Error())
	}

	return helper.CreateResult(true, data["tUser"], "Success")

}

func (a *UserController) DeleteUser(r *knot.WebContext) interface{} {
	r.Config.OutputType = knot.OutputJson

	payload := toolkit.M{}
	if err := r.GetPayload(&payload); err != nil {
		return helper.CreateResult(false, nil, err.Error())
	}

	err := DeleteUser(payload)
	if err != nil {
		return helper.CreateResult(false, nil, err.Error())
	}

	return helper.CreateResult(true, nil, "Delete User Success")
}
func (a *UserController) GetAccessUser(r *knot.WebContext) interface{} {
	r.Config.OutputType = knot.OutputJson

	payload := toolkit.M{}
	if err := r.GetPayload(&payload); err != nil {
		return helper.CreateResult(true, nil, err.Error())
	}
	AccessGrants, err := GetAccessUser(payload)
	if err != nil {
		return helper.CreateResult(true, nil, err.Error())
	}

	return helper.CreateResult(true, AccessGrants, "")
}

func (a *UserController) SaveUser(r *knot.WebContext) interface{} {
	r.Config.OutputType = knot.OutputJson

	payload := toolkit.M{}
	if err := r.GetPayload(&payload); err != nil {
		return helper.CreateResult(false, nil, err.Error())
	}

	if err := SaveUser(payload); err != nil {
		return helper.CreateResult(false, nil, err.Error())
	}

	return helper.CreateResult(true, nil, "save user success")
}

func (a *UserController) ChangePass(r *knot.WebContext) interface{} {
	r.Config.OutputType = knot.OutputJson

	payload := toolkit.M{}
	if err := r.GetPayload(&payload); err != nil {
		return helper.CreateResult(false, nil, err.Error())
	}

	if err := ChangePass(payload); err != nil {
		return helper.CreateResult(true, nil, err.Error())
	}
	return helper.CreateResult(true, nil, "Change Password Success")
}
