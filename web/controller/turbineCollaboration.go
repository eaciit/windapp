package controller

import (
	"time"

	. "eaciit/wfdemo-git/library/core"
	. "eaciit/wfdemo-git/library/helper"
	. "eaciit/wfdemo-git/library/models"
	"eaciit/wfdemo-git/web/helper"

	"github.com/eaciit/dbox"
	_ "github.com/eaciit/dbox/dbc/mongo"
	"github.com/eaciit/knot/knot.v1"
	// tk "github.com/eaciit/toolkit"
	_ "gopkg.in/mgo.v2/bson"
)

type TurbineCollaborationController struct {
	App
}

func CreateTurbineCollaborationController() *TurbineCollaborationController {
	var controller = new(TurbineCollaborationController)
	return controller
}

func (m *TurbineCollaborationController) GetLatest(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	payload := struct {
		Take    int
		Turbine string
		Project string
	}{}

	e := k.GetPayload(&payload)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	turbine := payload.Turbine
	project := payload.Project
	take := payload.Take
	tNow := time.Now()
	timestamp := time.Date(tNow.Year(), tNow.Month(), tNow.Day(), 0, 0, 0, 0, time.UTC)

	csr, e := DB().Connection.NewQuery().
		From(new(TurbineCollaborationModel).TableName()).
		Where(
			dbox.And(
				dbox.Eq("turbineid", turbine),
				dbox.Eq("projectid", project),
				dbox.Gte("date", timestamp),
			),
		).
		Order("-date").Take(take).
		Cursor(nil)
	defer csr.Close()

	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	turbColls := TurbineCollaborationModel{}
	e = csr.Fetch(&turbColls, 1, false)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	return helper.CreateResult(true, turbColls, "success")
}

func (m *TurbineCollaborationController) Save(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	p := struct {
		Id          string
		ResponseFor string
		TurbineId   string
		TurbineName string
		Feeder      string
		Project     string
		Date        string
		Status      string
		Remark      string
	}{}
	e := k.GetPayload(&p)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	date, _ := time.Parse("2006-01-02T15:04:05Z", p.Date)

	createdBy := ""     //tUser.LoginID
	createdByName := "" //tUser.FullName
	createdIp := ""     //tUser.Email
	createdLoc := ""

	mdl := new(TurbineCollaborationModel)

	mdl.ResponseFor = p.ResponseFor
	mdl.ProjectId = p.Project
	mdl.TurbineId = p.TurbineId
	mdl.TurbineName = p.TurbineName
	mdl.Feeder = p.Feeder
	mdl.Date = date
	mdl.DateInfo = GetDateInfo(date)
	mdl.Status = p.Status
	mdl.Remark = p.Remark
	mdl.CreatedBy = createdBy
	mdl.CreatedByName = createdByName
	mdl.CreatedOn = time.Now()
	mdl.CreatedIp = createdIp
	mdl.CreatedLoc = createdLoc

	if p.Id != "" {
		mdl.Id = p.Id
	} else {
		mdl = mdl.New()
	}

	e = DB().Save(mdl)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	return helper.CreateResult(true, nil, "success")
}
