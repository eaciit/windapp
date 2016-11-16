package models

import (
	"github.com/eaciit/acl/v1.0"
	"github.com/eaciit/toolkit"
	"strings"
	"time"
)

func GetSession(payload toolkit.M) (toolkit.M, error) {
	tSession := new(acl.Session)
	take := toolkit.ToInt(payload["take"], toolkit.RoundingAuto)
	skip := toolkit.ToInt(payload["skip"], toolkit.RoundingAuto)

	c, err := acl.Find(tSession, nil, toolkit.M{})
	if err != nil {
		return nil, err
	}

	data := toolkit.M{}
	result := toolkit.Ms{}
	arrm := make([]toolkit.M, 0, 0)
	if err := c.Fetch(&arrm, 0, false); err != nil {
		return nil, err
	}
	find := strings.ToLower(payload.GetString("search"))
	for _, val := range arrm {
		val.Set("duration", time.Since(val["created"].(time.Time)).Hours())
		val.Set("status", "ACTIVE")
		if val["expired"].(time.Time).Before(time.Now()) {
			val.Set("duration", val["expired"].(time.Time).Sub(val["created"].(time.Time)).Hours())
			val.Set("status", "EXPIRED")
		}
		if find != "" {
			if strings.Contains(strings.ToLower(val.GetString("status")), find) {
				result = append(result, val)
			} else if strings.Contains(strings.ToLower(val.GetString("loginid")), find) {
				result = append(result, val)
			}
		} else {
			result = append(result, val)
		}
	}
	c.Close()
	results := toolkit.Ms{}
	maxIndex := toolkit.SliceLen(result)
	if skip+take < maxIndex {
		maxIndex = skip + take
	}
	for i := skip; i < maxIndex; i++ {
		results = append(results, result[i])
	}

	data.Set("Datas", results)
	data.Set("total", toolkit.SliceLen(result))

	return data, nil

}

func SetExpired(payload toolkit.M) error {
	tSession := new(acl.Session)
	if err := acl.FindByID(tSession, payload.GetString("_id")); err != nil {
		return err
	}

	tSession.Expired = time.Now().UTC()

	if err := acl.Save(tSession); err != nil {
		return err
	}

	return nil
}
