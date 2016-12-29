package models

import (
	"github.com/eaciit/acl/v1.0"
	"github.com/eaciit/dbox"
	"github.com/eaciit/toolkit"
	"sort"
	"strconv"
)

func GetAccessQuery(payload toolkit.M) (toolkit.M, error) {
	var filters []*dbox.Filter
	var filter *dbox.Filter
	tAccess := new(acl.Access)
	boolList := []string{"false", "true"}

	if find := payload.GetString("search"); find != "" {
		bfind, err := strconv.ParseBool(find)

		if err == nil && toolkit.HasMember(boolList, find) {
			filters = append(filters, dbox.Eq("enable", bfind))
		}
		ifind, err := strconv.Atoi(find)
		if err == nil {
			filters = append(filters, dbox.Eq("index", ifind))
		} else {
			_filters := []*dbox.Filter{
				dbox.Contains("id", find),
				dbox.Contains("title", find),
				dbox.Contains("icon", find),
				dbox.Contains("url", find),
				dbox.Contains("parentid", find),
			}
			filters = append(filters, _filters...)
		}
	}

	if len(filters) > 0 {
		filter = dbox.Or(filters...)
	}
	data := toolkit.M{}
	arrm := make([]toolkit.M, 0, 0)

	take := toolkit.ToInt(payload["take"], toolkit.RoundingAuto)
	skip := toolkit.ToInt(payload["skip"], toolkit.RoundingAuto)

	c, err := acl.Find(tAccess, filter, toolkit.M{}.Set("take", take).Set("skip", skip))
	if err != nil {
		return nil, err
	}

	if err := c.Fetch(&arrm, 0, false); err != nil {
		return nil, err
	}

	c.Close()

	c, err = acl.Find(tAccess, filter, nil) /*for counting all the data*/
	if err != nil {
		return nil, err
	}

	data.Set("Data", arrm)
	data.Set("total", c.Count())

	return data, nil

}

func GetaccessDropDownQuery() ([]toolkit.M, error) {
	tAccess := new(acl.Access)

	arrm := make([]toolkit.M, 0, 0)
	c, err := acl.Find(tAccess, nil, nil)
	if err != nil {
		return nil, err
	}
	if err := c.Fetch(&arrm, 0, false); err != nil {
		return nil, err
	}
	orderID := []string{}
	for _, val := range arrm {
		orderID = append(orderID, val.GetString("_id"))
	}
	sort.Strings(orderID)
	result := toolkit.Ms{}
	for _, id := range orderID {
		for _, val := range arrm {
			if id == val.GetString("_id") {
				result = append(result, val)
			}
		}
	}

	return result, nil
}

func GetParentIDQuery() ([]string, error) {
	tAccess := new(acl.Access)

	arrm := make([]toolkit.M, 0, 0)
	c, err := acl.Find(tAccess, nil, nil)
	if err != nil {
		return nil, err
	}
	if err := c.Fetch(&arrm, 0, false); err != nil {
		return nil, err
	}
	result := []string{}
	for _, val := range arrm {
		result = append(result, val.GetString("id"))
	}
	sort.Strings(result)

	return result, nil
}

func FindAccess(payload toolkit.M) (toolkit.M, error) {
	tAccess := new(acl.Access)
	result := toolkit.M{}

	if err := acl.FindByID(tAccess, payload.GetString("_id")); err != nil {
		return nil, err
	}
	result.Set("tAccess", tAccess)

	return result, nil
}

func DeleteAccessProc(payload toolkit.M) error {
	idArray := payload.Get("_id").([]interface{})

	for _, id := range idArray {
		o := new(acl.Access)
		o.ID = toolkit.ToString(id)
		if err := acl.Delete(o); err != nil {
			return err
		}
	}
	return nil
}

func SaveAccessProc(payload toolkit.M) error {

	initAccess := new(acl.Access)
	initAccess.ID = payload.GetString("_id")
	initAccess.Title = payload.GetString("Title")
	initAccess.Icon = payload.GetString("Icon")
	initAccess.Index = payload.GetInt("Index")
	initAccess.ParentId = payload.GetString("ParentId")
	initAccess.Url = payload.GetString("Url")
	initAccess.Enable = payload["Enable"].(bool)
	initAccess.Category = 2

	if err := acl.Save(initAccess); err != nil {
		return err
	}

	return nil
}
