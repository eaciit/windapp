package controller

import (
	. "eaciit/wfdemo-git/library/models"
	"eaciit/wfdemo-git/web/helper"
	"errors"
	"strings"
	"time"

	"github.com/eaciit/acl/v1.0"
	"github.com/eaciit/dbox"
	"github.com/eaciit/knot/knot.v1"
	"github.com/eaciit/toolkit"

	. "eaciit/wfdemo-git/library/core"
)

type LoginController struct {
	App
}

func CreateLoginController() *LoginController {
	var controller = new(LoginController)
	return controller
}

func (l *LoginController) GetSession(r *knot.WebContext) interface{} {
	r.Config.OutputType = knot.OutputJson
	sessionId := r.Session("sessionid", "")
	return helper.CreateResult(true, sessionId, "")
}

func (l *LoginController) CheckCurrentSession(r *knot.WebContext) interface{} {
	r.Config.OutputType = knot.OutputJson
	sessionid := r.Session("sessionid", "")

	// toolkit.Printf("CheckCurrentSession: %#v \v", sessionid)

	if !acl.IsSessionIDActive(toolkit.ToString(sessionid)) {
		r.SetSession("sessionid", "")
		return helper.CreateResult(false, false, "inactive")
	}
	return helper.CreateResult(true, true, "active")
}

func (l *LoginController) GetMenuList(r *knot.WebContext) interface{} {
	r.Config.OutputType = knot.OutputJson
	menuList, err := GetListOfMenu(toolkit.ToString(r.Session("sessionid", "")))
	if err != nil {
		return helper.CreateResult(false, nil, err.Error())
	}

	// remarked by ams, 2016-10-15
	// issue ServerAddress is localhost blablabla
	// suggestion : if you want to check with the real address please put it (the address / base url) in the configuration file
	//              if not, please find by segment or split by /

	payload := toolkit.M{}
	if err := r.GetPayload(&payload); err != nil {
		return helper.CreateResult(false, "", err.Error())
	}
	maxURLLen := 4
	urlSplit := strings.SplitN(payload.GetString("url"), "/", maxURLLen)
	if len(urlSplit) == maxURLLen {
		url := "/" + urlSplit[maxURLLen-1]

		isFound := false
		if len(MenuList) > 0 {
			for _, val := range MenuList {
				if val == url {
					isFound = true
				}
			}
			if url == "/web/page/login" {
				isFound = true
			}
			if !isFound {
				return helper.CreateResult(false, "", "You don't have access to this page")
			}
		}
	}

	return helper.CreateResult(true, menuList, "")
}

func (l *LoginController) GetUserName(r *knot.WebContext) interface{} {
	r.Config.OutputType = knot.OutputJson
	sessionId := r.Session("sessionid", "")

	if toolkit.ToString(sessionId) == "" {
		err := error(errors.New("Sessionid is not found"))
		return helper.CreateResult(false, nil, err.Error())
	}
	tUser, err := GetUserName(sessionId)

	if err != nil {
		return helper.CreateResult(false, nil, "Get username failed")
	}

	return helper.CreateResult(true, tUser.LoginID, "")
}

func (l *LoginController) ProcessLogin(r *knot.WebContext) interface{} {
	r.Config.OutputType = knot.OutputJson

	lastDateData, _ := time.Parse("2006-01-02 15:04", "2016-10-31 23:59")

	payload := toolkit.M{}
	if err := r.GetPayload(&payload); err != nil {
		return helper.CreateResult(false, "", err.Error())
	}
	MenuList = []string{}
	menulis, sessid, err := LoginProcess(payload)
	if err != nil {
		return helper.CreateResult(false, "", err.Error())
	}
	WriteLog(sessid, "login", r.Request.URL.String())
	r.SetSession("sessionid", sessid)
	MenuList = menulis

	// temporary add last date hardcode, then will change to get it from database automatically
	// add by ams, 2016-10-04

	query := DB().Connection.NewQuery().From(new(ScadaData).TableName()).Order("-timestamp").Take(1)

	csr, e := query.Cursor(nil)
	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}
	defer csr.Close()

	Result := make([]ScadaData, 0)
	e = csr.Fetch(&Result, 0, false)

	if e != nil {
		return helper.CreateResult(false, nil, e.Error())
	}

	for _, val := range Result {
		// toolkit.Printf("Result : %s \n", val.TimeStamp.UTC())
		lastDateData = val.TimeStamp.UTC()
	}

	// toolkit.Printf("Result : %s \n", lastDateData)
	lastDateData = lastDateData.UTC()
	r.SetSession("lastdate_data", lastDateData)

	// Get Available Date All Collection
	Scadaresults := make([]time.Time, 0)
	Alarmresults := make([]time.Time, 0)
	JMRresults := make([]time.Time, 0)
	METresults := make([]time.Time, 0)
	Durationresults := make([]time.Time, 0)
	ScadaAnomalyresults := make([]time.Time, 0)
	AlarmOverlappingresults := make([]time.Time, 0)
	AlarmScadaAnomalyresults := make([]time.Time, 0)

	// Scada Data
	for i := 0; i < 2; i++ {
		var arrsort []string
		if i == 0 {
			arrsort = append(arrsort, "timestamp")
		} else {
			arrsort = append(arrsort, "-timestamp")
		}

		query := DB().Connection.NewQuery().From(new(ScadaData).TableName()).Skip(0).Take(1)
		query = query.Order(arrsort...)

		csr, e := query.Cursor(nil)
		if e != nil {
			return helper.CreateResult(false, nil, e.Error())
		}
		defer csr.Close()

		Result := make([]ScadaData, 0)
		e = csr.Fetch(&Result, 0, false)

		if e != nil {
			return helper.CreateResult(false, nil, e.Error())
		}

		for _, val := range Result {
			Scadaresults = append(Scadaresults, val.TimeStamp.UTC())
		}
	}

	// Alarm Data
	for i := 0; i < 2; i++ {
		var arrsort []string
		if i == 0 {
			arrsort = append(arrsort, "startdate")
		} else {
			arrsort = append(arrsort, "-startdate")
		}

		query := DB().Connection.NewQuery().From(new(Alarm).TableName()).Skip(0).Take(1)
		query = query.Order(arrsort...)

		csr, e := query.Cursor(nil)
		if e != nil {
			return helper.CreateResult(false, nil, e.Error())
		}
		defer csr.Close()

		Result := make([]Alarm, 0)
		e = csr.Fetch(&Result, 0, false)

		if e != nil {
			return helper.CreateResult(false, nil, e.Error())
		}

		for _, val := range Result {
			Alarmresults = append(Alarmresults, val.StartDate.UTC())
		}
	}

	// JMR Data
	for i := 0; i < 2; i++ {
		var arrsort []string
		if i == 0 {
			arrsort = append(arrsort, "dateinfo.dateid")
		} else {
			arrsort = append(arrsort, "-dateinfo.dateid")
		}

		query := DB().Connection.NewQuery().From(new(ScadaData).TableName()).Skip(0).Take(1)
		query = query.Order(arrsort...)

		csr, e := query.Cursor(nil)
		if e != nil {
			return helper.CreateResult(false, nil, e.Error())
		}
		defer csr.Close()

		Result := make([]ScadaData, 0)
		e = csr.Fetch(&Result, 0, false)

		if e != nil {
			return helper.CreateResult(false, nil, e.Error())
		}

		for _, val := range Result {
			JMRresults = append(JMRresults, val.DateInfo.DateId.UTC())
		}
	}

	// MET Data
	for i := 0; i < 2; i++ {
		var arrsort []string
		if i == 0 {
			arrsort = append(arrsort, "timestamp")
		} else {
			arrsort = append(arrsort, "-timestamp")
		}

		query := DB().Connection.NewQuery().From(new(MetTower).TableName()).Skip(0).Take(1)
		query = query.Order(arrsort...)

		csr, e := query.Cursor(nil)
		if e != nil {
			return helper.CreateResult(false, nil, e.Error())
		}
		defer csr.Close()

		Result := make([]MetTower, 0)
		e = csr.Fetch(&Result, 0, false)

		if e != nil {
			return helper.CreateResult(false, nil, e.Error())
		}

		for _, val := range Result {
			METresults = append(METresults, val.TimeStamp.UTC())
		}
	}

	// Duration Data
	for i := 0; i < 2; i++ {
		var arrsort []string
		if i == 0 {
			arrsort = append(arrsort, "timestamp")
		} else {
			arrsort = append(arrsort, "-timestamp")
		}

		query := DB().Connection.NewQuery().From(new(ScadaData).TableName()).Where(dbox.And(dbox.Eq("isvalidtimeduration", false))).Skip(0).Take(1)
		query = query.Order(arrsort...)

		csr, e := query.Cursor(nil)
		if e != nil {
			return helper.CreateResult(false, nil, e.Error())
		}
		defer csr.Close()

		Result := make([]ScadaData, 0)
		e = csr.Fetch(&Result, 0, false)

		if e != nil {
			return helper.CreateResult(false, nil, e.Error())
		}

		for _, val := range Result {
			Durationresults = append(Durationresults, val.TimeStamp.UTC())
		}
	}

	// Anomaly Data
	for i := 0; i < 2; i++ {
		var arrsort []string
		if i == 0 {
			arrsort = append(arrsort, "timestamp")
		} else {
			arrsort = append(arrsort, "-timestamp")
		}

		query := DB().Connection.NewQuery().From(new(ScadaData).TableName()).Where(dbox.And(dbox.Eq("isvalidtimeduration", true))).Skip(0).Take(1)
		query = query.Order(arrsort...)

		csr, e := query.Cursor(nil)
		if e != nil {
			return helper.CreateResult(false, nil, e.Error())
		}
		defer csr.Close()

		Result := make([]ScadaData, 0)
		e = csr.Fetch(&Result, 0, false)

		if e != nil {
			return helper.CreateResult(false, nil, e.Error())
		}

		for _, val := range Result {
			ScadaAnomalyresults = append(ScadaAnomalyresults, val.TimeStamp.UTC())
		}
	}

	// AlarmOverlapping Data
	for i := 0; i < 2; i++ {
		var arrsort []string
		if i == 0 {
			arrsort = append(arrsort, "startdate")
		} else {
			arrsort = append(arrsort, "-startdate")
		}

		query := DB().Connection.NewQuery().From(new(AlarmOverlapping).TableName()).Skip(0).Take(1)
		query = query.Order(arrsort...)

		csr, e := query.Cursor(nil)
		if e != nil {
			return helper.CreateResult(false, nil, e.Error())
		}
		defer csr.Close()

		Result := make([]AlarmOverlapping, 0)
		e = csr.Fetch(&Result, 0, false)

		if e != nil {
			return helper.CreateResult(false, nil, e.Error())
		}

		for _, val := range Result {
			AlarmOverlappingresults = append(AlarmOverlappingresults, val.StartDate.UTC())
		}
	}

	// AlarmScadaAnomaly Data
	for i := 0; i < 2; i++ {
		var arrsort []string
		if i == 0 {
			arrsort = append(arrsort, "startdate")
		} else {
			arrsort = append(arrsort, "-startdate")
		}

		query := DB().Connection.NewQuery().From(new(AlarmScadaAnomaly).TableName()).Skip(0).Take(1)
		query = query.Order(arrsort...)

		csr, e := query.Cursor(nil)
		if e != nil {
			return helper.CreateResult(false, nil, e.Error())
		}
		defer csr.Close()

		Result := make([]AlarmScadaAnomaly, 0)
		e = csr.Fetch(&Result, 0, false)

		if e != nil {
			return helper.CreateResult(false, nil, e.Error())
		}

		for _, val := range Result {
			AlarmScadaAnomalyresults = append(AlarmScadaAnomalyresults, val.StartDate.UTC())
		}
	}

	r.SetSession("scadaavaildate", Scadaresults)
	r.SetSession("alarmavaildate", Scadaresults)
	r.SetSession("jmravaildate", Scadaresults)
	r.SetSession("metavaildate", Scadaresults)
	r.SetSession("durationavaildate", Scadaresults)
	r.SetSession("scadaanomalyavaildate", Scadaresults)
	r.SetSession("alarmoverlappingavaildate", Scadaresults)
	r.SetSession("alarmscadaanomalyavaildate", Scadaresults)

	data := toolkit.M{
		"status":    true,
		"sessionid": sessid,
	}

	return helper.CreateResult(true, data, "Login Success")
}

func (l *LoginController) Logout(r *knot.WebContext) interface{} {
	r.Config.OutputType = knot.OutputJson

	err := SetExpired(toolkit.M{"_id": r.Session("sessionid", "")})
	if err != nil {
		return helper.CreateResult(false, nil, err.Error())
	}
	WriteLog(r.Session("sessionid", ""), "logout", r.Request.URL.String())
	r.SetSession("sessionid", "")

	return helper.CreateResult(true, nil, "Logout Success")
}

func (l *LoginController) ResetPassword(r *knot.WebContext) interface{} {
	r.Config.OutputType = knot.OutputJson

	payload := toolkit.M{}
	err := r.GetPayload(&payload)
	if err != nil {
		return helper.CreateResult(false, nil, err.Error())
	}
	if err = ResetPassword(payload); err != nil {
		return helper.CreateResult(false, nil, err.Error())
	}

	return helper.CreateResult(true, nil, "Reset Password Success")
}

func (l *LoginController) SavePassword(r *knot.WebContext) interface{} {
	r.Config.OutputType = knot.OutputJson

	payload := toolkit.M{}
	err := r.GetPayload(&payload)
	if err != nil {
		return helper.CreateResult(false, nil, err.Error())
	}

	if err = SavePassword(payload); err != nil {
		return helper.CreateResult(false, nil, err.Error())
	}

	return helper.CreateResult(true, nil, "Save Password Success")
}

func (l *LoginController) Authenticate(r *knot.WebContext) interface{} {
	r.Config.OutputType = knot.OutputJson

	payload := toolkit.M{}

	if err := r.GetPayload(&payload); err != nil {
		return helper.CreateResult(false, nil, err.Error())
	}

	result, err := AuthenticateProc(payload)
	if err != nil {
		return helper.CreateResult(false, nil, err.Error())
	}

	return helper.CreateResult(true, result, "Authenticate Success")
}
