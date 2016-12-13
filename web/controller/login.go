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

func getMenus(r *knot.WebContext) (interface{}, error) {
	menuList, err := GetListOfMenu(toolkit.ToString(r.Session("sessionid", "")))
	if err != nil {
		return nil, err
	}

	payload := toolkit.M{}
	if err := r.GetPayload(&payload); err != nil {
		return nil, err
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
				return nil, err
			}
		}
	}

	return menuList, err
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
	menus, sessid, err := LoginProcess(payload)
	if err != nil {
		return helper.CreateResult(false, "", err.Error())
	}
	WriteLog(sessid, "login", r.Request.URL.String())
	r.SetSession("sessionid", sessid)
	r.SetSession("menus", menus)
	MenuList = menus

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
	scadaResults := make([]time.Time, 2)
	dgrResults := make([]time.Time, 2)
	alarmResults := make([]time.Time, 2)
	jmrResults := make([]time.Time, 2)
	metResults := make([]time.Time, 2)
	durationResults := make([]time.Time, 2)
	scadaAnomalyresults := make([]time.Time, 2)
	alarmOverlappingresults := make([]time.Time, 2)
	alarmScadaAnomalyresults := make([]time.Time, 2)

	scadaResults[0], scadaResults[1], e = helper.GetDataDateAvailable(new(ScadaData).TableName(), "timestamp", nil)
	dgrResults[0], dgrResults[1], e = helper.GetDataDateAvailable(new(DGRModel).TableName(), "dateinfo.dateid", nil)
	alarmResults[0], alarmResults[1], e = helper.GetDataDateAvailable(new(Alarm).TableName(), "startdate", nil)
	alarmOverlappingresults[0], alarmOverlappingresults[1], e = helper.GetDataDateAvailable(new(AlarmOverlapping).TableName(), "startdate", nil)
	alarmScadaAnomalyresults[0], alarmScadaAnomalyresults[1], e = helper.GetDataDateAvailable(new(AlarmScadaAnomaly).TableName(), "startdate", nil)
	jmrResults[0], jmrResults[1], e = helper.GetDataDateAvailable(new(ScadaData).TableName(), "dateinfo.dateid", nil)
	metResults[0], metResults[1], e = helper.GetDataDateAvailable(new(MetTower).TableName(), "timestamp", nil)
	durationResults[0], durationResults[1], e = helper.GetDataDateAvailable(new(ScadaData).TableName(), "timestamp", dbox.And(dbox.Eq("isvalidtimeduration", false)))
	scadaAnomalyresults[0], scadaAnomalyresults[1], e = helper.GetDataDateAvailable(new(ScadaData).TableName(), "timestamp", dbox.And(dbox.Eq("isvalidtimeduration", true)))

	availdatedata := struct {
		ScadaData         []time.Time
		DGRData           []time.Time
		Alarm             []time.Time
		JMR               []time.Time
		MET               []time.Time
		Duration          []time.Time
		ScadaAnomaly      []time.Time
		AlarmOverlapping  []time.Time
		AlarmScadaAnomaly []time.Time
	}{
		ScadaData:         scadaResults,
		DGRData:           dgrResults,
		Alarm:             alarmResults,
		JMR:               jmrResults,
		MET:               metResults,
		Duration:          durationResults,
		ScadaAnomaly:      scadaAnomalyresults,
		AlarmOverlapping:  alarmOverlappingresults,
		AlarmScadaAnomaly: alarmScadaAnomalyresults,
	}

	r.SetSession("availdate", availdatedata)

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
