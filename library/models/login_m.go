package models

import (
	. "eaciit/wfdemo-git-dev/library/core"
	"errors"
	"fmt"
	"github.com/eaciit/acl/v1.0"
	"github.com/eaciit/dbox"
	"github.com/eaciit/orm"
	"github.com/eaciit/toolkit"
	"gopkg.in/gomail.v2"
)

type Login struct {
	orm.ModelBase
	ID       string `json:"_id",bson:"_id"`
	Password string
	Salt     string
}

func (l *Login) TableName() string {
	return "logins"
}

func (l *Login) RecordID() interface{} {
	return l.ID
}

func GetUserName(sessionId interface{}) (tUser acl.User, err error) {
	userid, err := acl.FindUserBySessionID(toolkit.ToString(sessionId))
	if err != nil {
		return
	}

	err = acl.FindByID(&tUser, userid)
	if err != nil {
		return
	}

	return
}

func GetListOfMenu(sessionid string) (result toolkit.Ms, err error) {
	if sessionid == "" {
		return nil, errors.New("Your session has expired, please log in again")
	}
	if !acl.IsSessionIDActive(sessionid) {
		return nil, errors.New("Your session has expired, please log in again")
	}
	result, err = acl.GetListMenuBySessionId(sessionid)
	if err != nil {
		return
	}

	return
}

func ResetPassword(payload toolkit.M) error {
	if !payload.Has("email") || !payload.Has("baseurl") {
		errors.New("Data is not complete")
	}

	uname, tokenid, err := acl.ResetPassword(toolkit.ToString(payload["email"]))

	if err != nil {
		err.Error()
	}

	linkstr := toolkit.Sprintf("<a href='%v/web/confirmreset?1=%v&2=%v'>Click</a>", toolkit.ToString(payload["baseurl"]), uname, tokenid)

	mailmsg := toolkit.Sprintf("Hi, <br/><br/> We received a request to reset your password, <br/><br/>")
	mailmsg = toolkit.Sprintf("%vFollow the link below to set a new password : <br/><br/> %v <br/><br/>", mailmsg, linkstr)
	mailmsg = toolkit.Sprintf("%vIf you don't want to change your password, you can ignore this email <br/><br/> Thanks,</body></html>", mailmsg)

	m := gomail.NewMessage()

	m.SetHeader("From", "admin.support@eaciit.com")
	m.SetHeader("To", toolkit.ToString(payload["email"]))

	m.SetHeader("Subject", "[no-reply] Self password reset")
	m.SetBody("text/html", mailmsg)

	d := gomail.NewPlainDialer("smtp.office365.com", 587, "admin.support@eaciit.com", "******")
	err = d.DialAndSend(m)

	if err != nil {
		return err
	}

	return nil
}

func SavePassword(payload toolkit.M) error {
	if !payload.Has("newpassword") || !payload.Has("userid") {
		return errors.New("Data is not complete")
	}

	switch {
	case payload.Has("tokenid"):
		acl.ChangePasswordToken(toolkit.ToString(payload["userid"]), toolkit.ToString(payload["newpassword"]), toolkit.ToString(payload["tokenid"]))
	default:
		// check sessionid first
		savedsessionid := "" //change with get session
		//=======================
		userid, err := acl.FindUserBySessionID(savedsessionid)
		if err == nil && userid == toolkit.ToString(payload["userid"]) {
			err = acl.ChangePassword(toolkit.ToString(payload["userid"]), toolkit.ToString(payload["newpassword"]))
		} else if err == nil {
			err = errors.New("Userid is not match")
		}
	}

	return nil
}

func AuthenticateProc(payload toolkit.M) (toolkit.M, error) {
	var iaccenum acl.AccessTypeEnum
	result := toolkit.M{}
	result.Set("hasaccess", false)

	switch toolkit.TypeName(payload["accesscheck"]) {
	case "[]interface {}":
		for _, val := range payload["accesscheck"].([]interface{}) {
			tacc := acl.GetAccessEnum(toolkit.ToString(val))
			if !acl.Matchaccess(int(tacc), int(iaccenum)) {
				iaccenum += tacc
			}
		}
	default:
		iaccenum = acl.GetAccessEnum(toolkit.ToString(payload["accesscheck"]))
	}

	found := acl.HasAccess(toolkit.ToString(payload["sessionid"]),
		acl.IDTypeSession,
		toolkit.ToString(payload["accessid"]),
		iaccenum)

	if found {
		result.Set("hasaccess", found)
	}

	return result, nil
}

func PrepareDefaultUser() (err error) {
	username := fmt.Sprintf("%v", GetConfig("default_username", ""))
	password := fmt.Sprintf("%v", GetConfig("default_password", ""))

	user := new(acl.User)
	filter := dbox.Contains("loginid", username)
	c, err := acl.Find(user, filter, nil)

	if err != nil {
		return
	}

	if c.Count() == 0 {
		user.ID = toolkit.RandomString(32)
		user.LoginID = username
		user.FullName = username
		user.Password = password
		user.Enable = true

		err = acl.Save(user)
		if err != nil {
			return
		}
		err = acl.ChangePassword(user.ID, password)
		if err != nil {
			return
		}

		toolkit.Printf(`Default user "%s" with standard password has been created%s`, username, "\n")
	}

	return
}

func LoginProcess(payload toolkit.M) ([]string, string, error) {
	switch {
	case !payload.Has("username") || !payload.Has("password"):
		return []string{}, "", errors.New("username or password not found")
	case payload.Has("username") && len(toolkit.ToString(payload["username"])) == 0:
		return []string{}, "", errors.New("username cannot be empty")
	case payload.Has("password") && len(toolkit.ToString(payload["password"])) == 0:
		return []string{}, "", errors.New("password cannot be empty")
	}

	sessid, err := acl.Login(toolkit.ToString(payload["username"]), toolkit.ToString(payload["password"]))
	if err != nil {
		return []string{}, "", err
	}
	result, err := getURLList(sessid)
	if err != nil {
		return []string{}, "", err
	}

	return result, sessid, nil
}

func getURLList(sessionid string) (urlList []string, err error) {
	isession := new(acl.Session)
	err = acl.FindByID(isession, sessionid)
	if err != nil {
		err = errors.New(toolkit.Sprintf("Get URL List : %s", err.Error()))
		return
	}

	iuser := new(acl.User)
	err = acl.FindUserByLoginID(iuser, isession.LoginID)
	if err != nil {
		err = errors.New(toolkit.Sprintf("Get URL List : %s", err.Error()))
		return
	}

	accessIDList := toolkit.M{}
	for _, val := range iuser.GetAccessList() {
		accessIDList.Set(val.AccessID, 1)
	}

	for _, val := range iuser.Groups {
		groupObj := new(acl.Group)
		_ = acl.FindByID(groupObj, val)
		for _, xval := range groupObj.GetAccessList() {
			accessIDList.Set(xval.AccessID, 1)
		}
	}

	listaccid := []interface{}{}
	for key := range accessIDList {
		listaccid = append(listaccid, key)
	}

	filter := dbox.And(dbox.In("_id", listaccid...), dbox.Eq("enable", true),
		dbox.Eq("category", acl.AccessMenu))

	csr, err := acl.Find(new(acl.Access), filter, nil)
	if err != nil {
		err = errors.New(toolkit.Sprintf("Get URL List : %s", err.Error()))
		return
	}

	defer csr.Close()

	accessList := make([]*acl.Access, 0)
	err = csr.Fetch(&accessList, 0, false)
	if err != nil {
		err = errors.New(toolkit.Sprintf("Get URL List : %s", err.Error()))
		return
	}

	for _, val := range accessList {
		if val.Url != "#" {
			urlList = append(urlList, val.Url)
		}
	}

	return
}
