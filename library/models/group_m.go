package models

import (
	"github.com/eaciit/acl/v1.0"
	"github.com/eaciit/dbox"
	"github.com/eaciit/toolkit"
	"strconv"
)

func GetGroup() ([]toolkit.M, error) {
	tGroup := new(acl.Group)

	arrm := make([]toolkit.M, 0, 0)
	c, err := acl.Find(tGroup, nil, nil)
	if err != nil {
		return nil, err
	}

	if err := c.Fetch(&arrm, 0, false); err != nil {
		return nil, err
	}

	return arrm, nil

}
func EditGroup(payload toolkit.M) (toolkit.M, error) {
	tGroup := new(acl.Group)
	result := toolkit.M{}

	if err := acl.FindByID(tGroup, payload.GetString("_id")); err != nil {
		return nil, err
	}
	result.Set("tGroup", tGroup)

	return result, nil

}

func SearchGroup(payload toolkit.M) (toolkit.M, error) {
	var filters []*dbox.Filter
	var filter *dbox.Filter
	tGroup := new(acl.Group)
	boolList := []string{"false", "true"}

	if find := payload.GetString("search"); find != "" {
		bfind, err := strconv.ParseBool(find)
		if err == nil && toolkit.HasMember(boolList, find) {
			filters = append(filters, dbox.Eq("enable", bfind))
		} else {
			_filters := []*dbox.Filter{
				dbox.Contains("id", find),
				dbox.Contains("title", find),
				dbox.Contains("owner", find),
			}
			filters = append(filters, _filters...)
		}
	}

	if len(filters) > 0 {
		filter = dbox.Or(filters...)
	}
	take := toolkit.ToInt(payload["take"], toolkit.RoundingAuto)
	skip := toolkit.ToInt(payload["skip"], toolkit.RoundingAuto)

	c, err := acl.Find(tGroup, filter, toolkit.M{}.Set("take", take).Set("skip", skip))
	if err != nil {
		return nil, err
	}

	data := toolkit.M{}
	arrm := make([]toolkit.M, 0, 0)
	if err := c.Fetch(&arrm, 0, false); err != nil {
		return nil, err
	}

	data.Set("Data", arrm)
	data.Set("total", c.Count())

	return data, nil

}

func GetAccessGroup(payload toolkit.M) ([]interface{}, error) {
	tGroup := new(acl.Group)

	if err := acl.FindByID(tGroup, payload.GetString("_id")); err != nil {
		return nil, err
	}
	var AccessGrants = []interface{}{}
	for _, v := range tGroup.Grants {
		var access = toolkit.M{}
		access.Set("AccessID", v.AccessID)
		access.Set("AccessValue", acl.Splitinttogrant(int(v.AccessValue)))
		AccessGrants = append(AccessGrants, access)
	}

	return AccessGrants, nil
}

func DeleteGroup(payload toolkit.M) error {
	idArray := payload.Get("_id").([]interface{})

	for _, id := range idArray {
		o := new(acl.Group)
		o.ID = toolkit.ToString(id)
		if err := acl.Delete(o); err != nil {
			return err
		}
	}
	return nil
}

func SaveGroup(payload toolkit.M) error {
	g := payload["group"].(map[string]interface{})

	initGroup := new(acl.Group)
	initGroup.ID = g["_id"].(string)
	initGroup.Title = g["Title"].(string)
	initGroup.Owner = g["Owner"].(string)
	initGroup.Enable = g["Enable"].(bool)

	if g["GroupType"].(float64) == 1 {
		initGroup.GroupType = acl.GroupTypeLdap
	} else if g["GroupType"].(float64) == 0 {
		initGroup.GroupType = acl.GroupTypeBasic
	}

	if err := acl.Save(initGroup); err != nil {
		return err
	}

	var grant map[string]interface{}
	for _, p := range payload["grants"].([]interface{}) {
		grant = p.(map[string]interface{})
		AccessID := grant["AccessID"].(string)
		Accessvalue := grant["AccessValue"]
		for _, v := range Accessvalue.([]interface{}) {
			switch v {
			case "AccessCreate":
				initGroup.Grant(AccessID, acl.AccessCreate)
			case "AccessRead":
				initGroup.Grant(AccessID, acl.AccessRead)
			case "AccessUpdate":
				initGroup.Grant(AccessID, acl.AccessUpdate)
			case "AccessDelete":
				initGroup.Grant(AccessID, acl.AccessDelete)
			case "AccessSpecial1":
				initGroup.Grant(AccessID, acl.AccessSpecial1)
			case "AccessSpecial2":
				initGroup.Grant(AccessID, acl.AccessSpecial2)
			case "AccessSpecial3":
				initGroup.Grant(AccessID, acl.AccessSpecial3)
			case "AccessSpecial4":
				initGroup.Grant(AccessID, acl.AccessSpecial4)
			}
		}
	}

	if err := acl.Save(initGroup); err != nil {
		return err
	}

	return nil
}
