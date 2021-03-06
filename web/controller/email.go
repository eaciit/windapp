package controller

import (
	. "eaciit/wfdemo-git/library/core"
	. "eaciit/wfdemo-git/library/models"
	"eaciit/wfdemo-git/web/helper"
	"github.com/eaciit/acl/v1.0"
	"github.com/eaciit/dbox"
	"github.com/eaciit/knot/knot.v1"
	"github.com/eaciit/toolkit"
	"gopkg.in/gomail.v2"
	"sort"
	"strconv"
	"strings"
)

type EmailController struct {
	App
}

func CreateEmailController() *EmailController {
	var controller = new(EmailController)
	return controller
}

func GetCategoryMail() (result []toolkit.M, e error) {
	csr, e := DB().Connection.NewQuery().
		From("ref_emailCategory").
		Order("category").
		Cursor(nil)

	if e != nil {
		return
	}
	defer csr.Close()

	data := []toolkit.M{}
	e = csr.Fetch(&data, 0, false)

	keyVal := map[string]string{}
	condCat := map[string]string{}
	keyList := []string{}
	for _, val := range data {
		keyVal[val.GetString("category")] = val.GetString("_id")
		condCat[val.GetString("category")] = val.GetString("condition")
		keyList = append(keyList, val.GetString("category"))
	}
	sort.Strings(keyList)

	for _, val := range keyList {
		res := toolkit.M{
			"value":     keyVal[val],
			"text":      val,
			"condition": condCat[val],
		}
		result = append(result, res)
	}
	if e != nil {
		return
	}

	return
}

func GetUserMail() (result []toolkit.M, e error) {
	csr, e := DB().Connection.NewQuery().
		Select("_id", "fullname").
		From("acl_users").
		Order("fullname").
		Cursor(nil)

	if e != nil {
		return
	}
	defer csr.Close()

	data := []toolkit.M{}
	e = csr.Fetch(&data, 0, false)

	keyVal := map[string]string{}
	keyList := []string{}
	for _, val := range data {
		keyVal[val.GetString("fullname")] = val.GetString("_id")
		keyList = append(keyList, val.GetString("fullname"))
	}
	sort.Strings(keyList)

	for _, val := range keyList {
		res := toolkit.M{
			"value": keyVal[val],
			"text":  val,
		}
		result = append(result, res)
	}
	if e != nil {
		return
	}

	return
}

func GetAlarmCodesMail() (result []toolkit.M, e error) {
	csr, e := DB().Connection.NewQuery().
		Select("alarmname").
		From(new(AlarmBrake).TableName()).
		Order("alarmname").
		Cursor(nil)

	if e != nil {
		return
	}
	defer csr.Close()

	data := []toolkit.M{}
	e = csr.Fetch(&data, 0, false)

	for _, val := range data {
		res := toolkit.M{
			"value": val.GetString("alarmname"),
			"text":  val.GetString("alarmname"),
		}
		result = append(result, res)
	}
	if e != nil {
		return
	}

	return
}

func GetTemplateMail() (result toolkit.M, e error) {
	csr, e := DB().Connection.NewQuery().
		From(new(EmailManagement).TableName()).
		Where(dbox.In("_id", []interface{}{"alarmtemplate", "datatemplate"}...)).
		Cursor(nil)

	if e != nil {
		return
	}
	defer csr.Close()

	data := []toolkit.M{}
	e = csr.Fetch(&data, 0, false)
	if e != nil {
		return
	}

	res := toolkit.M{}
	for _, val := range data {
		if val.GetString("_id") == "alarmtemplate" {
			res.Set("alarmTemplate", val.GetString("template"))
		} else {
			res.Set("dataTemplate", val.GetString("template"))
		}
	}
	result = res

	return
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

func getMailListFromReceivers(receivers []string) (err error, mailList []string) {
	pipes := []toolkit.M{}
	pipes = append(pipes, toolkit.M{"$match": toolkit.M{"_id": toolkit.M{"$in": receivers}}})
	pipes = append(pipes, toolkit.M{"$group": toolkit.M{
		"_id":      "mailList",
		"mailList": toolkit.M{"$push": "$email"},
	}})
	csrUser, err := DB().Connection.NewQuery().
		From(new(acl.User).TableName()).
		Command("pipe", pipes).
		Cursor(nil)
	if err != nil {
		return
	}
	defer csrUser.Close()

	dataUser := map[string][]string{}
	err = csrUser.Fetch(&dataUser, 1, false)
	if err != nil {
		return
	}
	mailList = dataUser["mailList"]

	return
}

func SendEmail(templateID string) error {
	csr, err := DB().Connection.NewQuery().
		From(new(EmailManagement).TableName()).
		Where(dbox.Eq("_id", templateID)).
		Cursor(nil)
	if err != nil {
		return err
	}
	defer csr.Close()

	dataEmail := new(EmailManagement)
	err = csr.Fetch(&dataEmail, 1, false)
	if err != nil {
		return err
	}
	err, mailList := getMailListFromReceivers(dataEmail.Receivers)
	if err != nil {
		return err
	}

	m := gomail.NewMessage()

	m.SetHeader("From", "admin.support@eaciit.com")
	m.SetHeader("To", mailList...)
	m.SetHeader("Subject", dataEmail.Subject)
	m.SetBody("text/html", dataEmail.Template)

	d := gomail.NewPlainDialer("smtp.office365.com", 587, "admin.support@eaciit.com", "B920Support")
	err = d.DialAndSend(m)

	if err != nil {
		return err
	}

	return nil
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

	query := DB().Connection.NewQuery().
		From(new(EmailManagement).TableName()).
		Skip(payload.GetInt("skip")).
		Take(payload.GetInt("take"))

	if payload.GetString("search") != "" {
		query.Where(dbox.And(dbox.Nin("_id", []interface{}{"alarmtemplate", "datatemplate"}...), dbox.Or(filters...)))
	} else {
		query.Where(dbox.Nin("_id", []interface{}{"alarmtemplate", "datatemplate"}...))
	}
	csr, err := query.Cursor(nil)

	if err != nil {
		return helper.CreateResult(false, nil, err.Error())
	}
	defer csr.Close()

	data := toolkit.M{}
	result := []toolkit.M{}
	err = csr.Fetch(&result, 0, false)
	if err != nil {
		return helper.CreateResult(false, nil, err.Error())
	}

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
	userID, err := acl.FindUserBySessionID(toolkit.ToString(r.Session("sessionid", "")))
	if err != nil {
		return helper.CreateResult(false, nil, err.Error())
	}
	payload.UpdatedBy = userID
	if payload.CreatedBy == "" {
		payload.CreatedBy = userID
	}

	if err := DB().Save(payload); err != nil {
		return helper.CreateResult(false, nil, err.Error())
	}

	err, mailList := getMailListFromReceivers(payload.Receivers)
	if err != nil {
		return err
	}
	data := toolkit.M{}

	data.Set("_id", payload.ID)
	data.Set("subject", payload.Subject)
	data.Set("category", payload.Category)
	data.Set("receivers", strings.Join(mailList, ","))
	if len(payload.AlarmCodes) > 0 {
		data.Set("alarmcodes", strings.Join(payload.AlarmCodes, ","))
	} else {
		data.Set("alarmcodes", "")
	}
	data.Set("intervaltime", payload.IntervalTime)
	data.Set("template", payload.Template)
	data.Set("enable", payload.Enable)
	data.Set("createddate", payload.CreatedDate)
	data.Set("lastupdate", payload.LastUpdate)
	data.Set("createdby", payload.CreatedBy)
	data.Set("updatedby", payload.UpdatedBy)

	return helper.CreateResult(true, data, "Save Email Success")
}
