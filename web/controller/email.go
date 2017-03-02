package controller

import (
	. "eaciit/wfdemo-git/library/core"
	. "eaciit/wfdemo-git/library/models"
	"eaciit/wfdemo-git/web/helper"

	// "github.com/eaciit/acl/v1.0"
	"github.com/eaciit/dbox"
	"github.com/eaciit/knot/knot.v1"
	"github.com/eaciit/toolkit"
	"strconv"
)

type EmailController struct {
	App
}

func CreateEmailController() *EmailController {
	var controller = new(EmailController)
	return controller
}

func (a *EmailController) EditEmail(r *knot.WebContext) interface{} {
	r.Config.OutputType = knot.OutputJson

	payload := toolkit.M{}

	if err := r.GetPayload(&payload); err != nil {
		return helper.CreateResult(false, nil, err.Error())
	}

	csr, err := DB().Connection.NewQuery().
		From(new(EmailManagement).TableName()).
		Where(dbox.Eq("_id", payload.GetString("_id"))).
		Cursor(nil)
	if err != nil {
		return helper.CreateResult(false, nil, err.Error())
	}
	defer csr.Close()

	data := toolkit.M{}
	err = csr.Fetch(&data, 1, false)

	if err != nil {
		return helper.CreateResult(false, nil, err.Error())
	}

	return helper.CreateResult(true, data, "success")

}

func (a *EmailController) Search(r *knot.WebContext) interface{} {
	r.Config.OutputType = knot.OutputJson

	payload := toolkit.M{}

	if err := r.GetPayload(&payload); err != nil {
		return helper.CreateResult(false, nil, err.Error())
	}

	var filters []*dbox.Filter
	boolList := []string{"false", "true"}

	if find := payload.GetString("search"); find != "" {
		bfind, err := strconv.ParseBool(find)
		if err == nil && toolkit.HasMember(boolList, find) {
			filters = append(filters, dbox.Eq("enable", bfind))
		} else {
			_filters := []*dbox.Filter{
				dbox.Contains("id", find),
				dbox.Contains("subject", find),
				dbox.Contains("category", find),
				dbox.Contains("template", find),
			}
			filters = append(filters, _filters...)
		}
	}

	csr, err := DB().Connection.NewQuery().
		From(new(EmailManagement).TableName()).
		Where(dbox.Or(filters...)).
		Skip(payload.GetInt("take")).
		Take(payload.GetInt("skip")).
		Cursor(nil)

	if err != nil {
		return helper.CreateResult(false, nil, err.Error())
	}
	defer csr.Close()

	data := toolkit.M{}
	result := make([]toolkit.M, 0, 0)
	err = csr.Fetch(&result, 0, false)

	data.Set("Data", result)
	data.Set("total", csr.Count())
	if err != nil {
		return helper.CreateResult(false, nil, err.Error())
	}

	return helper.CreateResult(true, data, "success")
}

func (a *EmailController) DeleteEmail(r *knot.WebContext) interface{} {
	r.Config.OutputType = knot.OutputJson

	payload := toolkit.M{}
	if err := r.GetPayload(&payload); err != nil {
		return helper.CreateResult(false, nil, err.Error())
	}

	idArray := payload.Get("_id").([]interface{})

	for _, id := range idArray {
		o := new(EmailManagement)
		o.ID = toolkit.ToString(id)
		if err := DB().Delete(o); err != nil {
			return helper.CreateResult(false, nil, err.Error())
		}
	}

	return helper.CreateResult(true, nil, "Delete Email Success")
}

func (a *EmailController) SaveEmail(r *knot.WebContext) interface{} {
	r.Config.OutputType = knot.OutputJson

	payload := new(EmailManagement)
	if err := r.GetPayload(&payload); err != nil {
		return helper.CreateResult(false, nil, err.Error())
	}

	if err := DB().Save(payload); err != nil {
		return helper.CreateResult(false, nil, err.Error())
	}

	return helper.CreateResult(true, nil, "Save Email Success")
}
