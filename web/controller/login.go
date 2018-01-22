package controller

import (
	. "eaciit/wfdemo-git/library/models"
	"eaciit/wfdemo-git/web/helper"
	"errors"
	"reflect"
	"strings"
	"time"

	"github.com/eaciit/acl/v1.0"
	// "github.com/eaciit/dbox"
	"github.com/eaciit/knot/knot.v1"
	"github.com/eaciit/toolkit"

	. "eaciit/wfdemo-git/library/core"
	"sync"
)

type LoginController struct {
	App
	mux sync.RWMutex
}

type Availdatedata struct {
	ScadaData         []time.Time
	DGRData           []time.Time
	Alarm             []time.Time
	JMR               []time.Time
	MET               []time.Time
	Duration          []time.Time
	ScadaAnomaly      []time.Time
	AlarmOverlapping  []time.Time
	AlarmScadaAnomaly []time.Time
	ScadaDataOEM      []time.Time
	ScadaDataHFD      []time.Time
	Warning           []time.Time
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

	if !acl.IsSessionIDActive(toolkit.ToString(sessionid)) {
		toolkit.Printf(">> CheckCurrentSession - notactive: %#v \v", sessionid)
		r.SetSession("sessionid", "")
		return helper.CreateResult(false, false, "inactive")
	}
	return helper.CreateResult(true, true, "active")
}

func (l *LoginController) GetMenuList(r *knot.WebContext) interface{} {
	r.Config.OutputType = knot.OutputJson
	menuList, err := GetListOfMenu(toolkit.ToString(r.Session("sessionid", "")))
	if err != nil {
		return helper.CreateResult(false, "", err.Error())
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

	return helper.CreateResult(true, menuList, "success")
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

func (l *LoginController) LoginRealtime(r *knot.WebContext) interface{} {
	r.Config.OutputType = knot.OutputJson

	lastDateData, _ := time.Parse("2006-01-02 15:04", "2016-10-31 23:59")

	if r.Session("keyRealtime", "") == "58e4931965d1041094641f0f" {
		MenuList = []string{}
		credentials := toolkit.M{"username": "realtime", "password": "realtime@1234"}
		menus, sessid, err := LoginProcess(credentials)
		if err != nil {
			return helper.CreateResult(false, "", err.Error())
		}

		WriteLog(sessid, "realtime login", r.Request.URL.String())
		r.SetSession("sessionid", sessid)
		r.SetSession("menus", menus)
		l.mux.Lock()
		helper.WC = r
		l.mux.Unlock()
		MenuList = menus

		datePeriod := getLastAvailDate()
		r.SetSession("availdate", datePeriod)

		lastDateData = datePeriod.ScadaData[1].UTC()
		r.SetSession("lastdate_data", lastDateData)

		data := toolkit.M{
			"status":    true,
			"sessionid": sessid,
		}

		return helper.CreateResult(true, data, "Login Success")
	}

	return helper.CreateResult(false, "", "You have no right to access this page")
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

	// log.Printf("sessid: %v \n", sessid)

	WriteLog(sessid, "login", r.Request.URL.String())
	r.SetSession("sessionid", sessid)
	r.SetSession("menus", menus)
	l.mux.Lock()
	helper.WC = r
	l.mux.Unlock()
	MenuList = menus

	// Get Available Date All Collection
	datePeriod := getLastAvailDate()
	r.SetSession("availdate", datePeriod)
	// r.SetSession("availdateall", getLastAvailDateAll())

	// log.Printf("availdate: %v \n", r.Session("availdate", ""))

	lastDateData = datePeriod.ScadaData[1].UTC()
	r.SetSession("lastdate_data", lastDateData)

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

func getLastAvailDate() *Availdatedata {
	latestDataPeriods := make([]LatestDataPeriod, 0)
	csr, e := DB().Connection.NewQuery().From(new(LatestDataPeriod).TableName()).Cursor(nil)
	if e != nil {
		return nil
	}

	e = csr.Fetch(&latestDataPeriods, 0, false)
	csr.Close()

	// toolkit.Println(latestDataPeriods)

	// mapCheck := map[string]time.Time{}

	datePeriod := new(Availdatedata)
	xdp := reflect.ValueOf(datePeriod).Elem()
	for _, d := range latestDataPeriods {

		for i, tval := range d.Data {
			d.Data[i] = tval.UTC()
		}

		f := xdp.FieldByName(d.Type)
		if f.IsValid() {
			if f.CanSet() {
				if f.Len() > 0 {
					if f.Len() == 2 {
						if (d.Data[0].Sub(f.Index(0).Interface().(time.Time)) < 0 && d.Data[0].Year() > 1) ||
							f.Index(0).Interface().(time.Time).Year() == 1 {
							f.Index(0).Set(reflect.ValueOf(d.Data[0]))
						}
						if d.Data[1].Sub(f.Index(1).Interface().(time.Time)) > 0 {
							f.Index(1).Set(reflect.ValueOf(d.Data[1]))
						}
					}
				} else {
					f.Set(reflect.ValueOf(d.Data))
				}
			}
		}
	}

	return datePeriod
}

func getLastAvailDateAll() toolkit.M {
	latestDataPeriods := make([]LatestDataPeriod, 0)
	csr, e := DB().Connection.NewQuery().From(new(LatestDataPeriod).TableName()).Cursor(nil)
	if e != nil {
		return nil
	}

	e = csr.Fetch(&latestDataPeriods, 0, false)
	csr.Close()

	result := toolkit.M{}
	for _, val := range latestDataPeriods {
		for i, tval := range val.Data {
			val.Data[i] = tval.UTC()
		}
		if result.Has(val.ProjectName) {
			currData, _ := toolkit.ToM(result[val.ProjectName])
			currData.Set(val.Type, val.Data)
			result.Set(val.ProjectName, currData)
		} else {
			result.Set(val.ProjectName, toolkit.M{val.Type: val.Data})
		}
	}

	return result
}

func getLastAvailDate_DRAFT() map[string]*Availdatedata {
	//contoh akses data
	// lastDateData = datePeriod["All"].ScadaData[1].UTC()
	allProject := []string{"All"}
	projectList, e := getProject()
	if e != nil {
		return nil
	}
	allProject = append(allProject, projectList...)

	result := map[string]*Availdatedata{}

	latestDataPeriodsList := make([]LatestDataPeriod, 0)
	query := DB().Connection.NewQuery().From(new(LatestDataPeriod).TableName())
	csr, e := query.Cursor(nil)
	if e != nil {
		return nil
	}

	e = csr.Fetch(&latestDataPeriodsList, 0, false)
	csr.Close()

	latestDataPeriods := make([]LatestDataPeriod, 0)
	for _, project := range allProject {
		latestDataPeriods = make([]LatestDataPeriod, 0)
		for _, latestDate := range latestDataPeriodsList {
			if project != "All" {
				if latestDate.ProjectName == project {
					latestDataPeriods = append(latestDataPeriods, latestDate)
				}
			} else {
				latestDataPeriods = append(latestDataPeriods, latestDate)
			}
		}
		datePeriod := new(Availdatedata)
		xdp := reflect.ValueOf(datePeriod).Elem()
		for _, d := range latestDataPeriods {
			f := xdp.FieldByName(d.Type)
			if f.IsValid() {
				if f.CanSet() {
					if f.Len() > 0 {
						if f.Len() == 2 {
							if (d.Data[0].Sub(f.Index(0).Interface().(time.Time)) < 0 && d.Data[0].Year() > 1) ||
								f.Index(0).Interface().(time.Time).Year() == 1 {
								f.Index(0).Set(reflect.ValueOf(d.Data[0]))
							}
							if d.Data[1].Sub(f.Index(1).Interface().(time.Time)) > 0 {
								f.Index(1).Set(reflect.ValueOf(d.Data[1]))
							}
						}
					} else {
						f.Set(reflect.ValueOf(d.Data))
					}
				}
			}
		}
		result[project] = datePeriod
	}

	return result
}
