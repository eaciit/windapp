package controller

import (
	. "eaciit/wfdemo/library/core"
	. "eaciit/wfdemo/library/models"
	"eaciit/wfdemo/web/helper"
	"errors"
	"github.com/eaciit/acl/v1.0"
	"github.com/eaciit/dbox"
	"github.com/eaciit/knot/knot.v1"
	"github.com/eaciit/toolkit"
	// "time"
)

type UserPreferencesController struct {
	App
}

var (
	maxViews = 3
)

func CreateUserPreferencesController() *UserPreferencesController {
	var controller = new(UserPreferencesController)
	return controller
}

func (m *UserPreferencesController) GetSavedKPIAnalysis(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	tUser, e := GetUserLoginDetails(k)

	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	csr, e := DB().Connection.NewQuery().
		From(new(UserPreferences).TableName()).
		Where(dbox.Eq("loginid", tUser.LoginID)).
		Cursor(nil)

	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}
	defer csr.Close()

	data := []UserPreferences{}
	e = csr.Fetch(&data, 0, false)

	result := []KPIAnalysis{}
	if len(data) == 1 {
		result = data[0].KPIAnalysis
	}

	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	return helper.CreateResult(true, result, "success")
}

func (m *UserPreferencesController) SaveKPI(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	tUser, e := GetUserLoginDetails(k)

	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	p := struct {
		OldName         string
		Name            string
		KeyA            string
		KeyB            string
		KeyC            string
		ColumnBreakdown string
		RowBreakdown    string
	}{}
	e = k.GetPayload(&p)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	csr, e := DB().Connection.NewQuery().
		From(new(UserPreferences).TableName()).
		Where(dbox.Eq("loginid", tUser.LoginID)).
		Cursor(nil)

	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}
	defer csr.Close()

	data := []UserPreferences{}
	e = csr.Fetch(&data, 0, false)

	mdl := UserPreferences{}
	kpiAnalysis := []KPIAnalysis{}

	x := KPIAnalysis{
		Name:            p.Name,
		KeyA:            p.KeyA,
		KeyB:            p.KeyB,
		KeyC:            p.KeyC,
		ColumnBreakdown: p.ColumnBreakdown,
		RowBreakdown:    p.RowBreakdown,
	}

	// toolkit.Printf("%#v \n", x)

	if len(data) == 1 {
		mdl = data[0]
		updateExisting := false
		for _, val := range data[0].KPIAnalysis {
			if val.Name == p.OldName {
				kpiAnalysis = append(kpiAnalysis, x)
				updateExisting = true
			} else {
				kpiAnalysis = append(kpiAnalysis, val)
			}
		}

		if !updateExisting {
			kpiAnalysis = append(kpiAnalysis, x)
		}

		mdl.KPIAnalysis = kpiAnalysis
	} else {
		mdl.Id = tUser.ID
		mdl.LoginID = tUser.LoginID
		mdl.KPIAnalysis = []KPIAnalysis{}
		mdl.KPIAnalysis = append(mdl.KPIAnalysis, x)
	}

	// toolkit.Printf("%#v \n", mdl)

	if len(mdl.KPIAnalysis) > maxViews {
		return helper.CreateResult(false, nil, "Maximum 3 Views are allowed")
	}

	e = DB().Save(&mdl)

	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	return helper.CreateResult(true, mdl.KPIAnalysis, "success")
}

func (m *UserPreferencesController) GetAnalysisStudioViews(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	tUser, e := GetUserLoginDetails(k)

	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	csr, e := DB().Connection.NewQuery().
		From(new(UserPreferences).TableName()).
		Where(dbox.Eq("loginid", tUser.LoginID)).
		Cursor(nil)

	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}
	defer csr.Close()

	data := []UserPreferences{}
	e = csr.Fetch(&data, 0, false)

	result := []AnalysisStudio{}
	if len(data) == 1 {
		result = data[0].AnalysisStudio
	}

	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	return helper.CreateResult(true, result, "success")
}

func (m *UserPreferencesController) SaveAnalysisStudioViews(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	tUser, e := GetUserLoginDetails(k)

	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	p := struct {
		OldName string
		Name    string
		Keys    []string
		Filters []Filter
	}{}
	e = k.GetPayload(&p)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	csr, e := DB().Connection.NewQuery().
		From(new(UserPreferences).TableName()).
		Where(dbox.Eq("loginid", tUser.LoginID)).
		Cursor(nil)

	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}
	defer csr.Close()

	data := []UserPreferences{}
	e = csr.Fetch(&data, 0, false)

	mdl := UserPreferences{}
	analysisStudio := []AnalysisStudio{}

	x := AnalysisStudio{
		Name:    p.Name,
		Keys:    p.Keys,
		Filters: p.Filters,
	}

	// toolkit.Printf("%#v \n", x)

	if len(data) == 1 {
		mdl = data[0]
		updateExisting := false
		for _, val := range data[0].AnalysisStudio {
			if val.Name == p.OldName {
				analysisStudio = append(analysisStudio, x)
				updateExisting = true
			} else {
				analysisStudio = append(analysisStudio, val)
			}
		}

		if !updateExisting {
			analysisStudio = append(analysisStudio, x)
		}

		mdl.AnalysisStudio = analysisStudio
	} else {
		mdl.Id = tUser.ID
		mdl.LoginID = tUser.LoginID
		mdl.AnalysisStudio = []AnalysisStudio{}
		mdl.AnalysisStudio = append(mdl.AnalysisStudio, x)
	}

	// toolkit.Printf("%#v \n", mdl)

	if len(mdl.AnalysisStudio) > maxViews {
		return helper.CreateResult(false, nil, "Maximum 3 Views are allowed")
	}

	e = DB().Save(&mdl)

	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	return helper.CreateResult(true, mdl.AnalysisStudio, "success")
}

func GetUserLoginDetails(k *knot.WebContext) (tUser acl.User, err error) {
	sessionId := k.Session("sessionid", "")

	if toolkit.ToString(sessionId) == "" {
		err = error(errors.New("Sessionid is not found"))
		return
	}
	tUser, err = GetUserName(sessionId)
	return
}
