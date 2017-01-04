package models

import (
	"github.com/eaciit/orm"
	"gopkg.in/mgo.v2/bson"
)

type WRA struct {
	orm.ModelBase `bson:"-",json:"-"`
	Id            bson.ObjectId ` bson:"_id" , json:"_id" `
	Time          string
	Jan           string
	Feb           string
	Mar           string
	Apr           string
	May           string
	Jun           string
	Jul           string
	Aug           string
	Sep           string
	Oct           string
	Nov           string
	Des           string
}

/*func NewWRA() *WRA {
	m := new(WRA)
	m.Id = bson.NewObjectId()
	return m
}*/

func (e *WRA) RecordID() interface{} {
	return e.Id
}

func (m *WRA) TableName() string {
	return "ref_wra"
}
