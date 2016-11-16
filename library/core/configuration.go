package gocore

import (
	"github.com/eaciit/orm"
	"github.com/eaciit/toolkit"
)

const (
	CONF_DB_ACL   string = "db_acl"
	CONF_DB_OSTRO string = "db_ostro"
)

type Configuration struct {
	orm.ModelBase
	ID   string `json:"_id",bson:"_id"`
	Data interface{}
}

func (a *Configuration) TableName() string {
	return "configurations"
}

func (a *Configuration) RecordID() interface{} {
	return a.ID
}

func GetConfig(key string, args ...string) interface{} {
	data := new(Configuration)
	if err := GetData(data, key); err != nil {
		return err
	}
	return data.Data
}

func (p *Configuration) GetPort() (int, error) {
	if err := GetData(p, p.ID); err != nil {
		return 0, err
	}
	toolkit.Println("port", p.ID, p.Data)

	return toolkit.ToInt(p.Data, toolkit.RoundingAuto), nil
}

func (p *Configuration) GetDB() (toolkit.M, error) {
	if err := GetData(p, p.ID); err != nil {
		return nil, err
	}
	data, err := toolkit.ToM(p.Data)
	if err != nil {
		return nil, err
	}

	return data, nil
}
