package models

import (
	// . "eaciit/wfdemo-git/library/helper"
	"time"

	"github.com/eaciit/orm"
	"gopkg.in/mgo.v2/bson"
)

type EmailManagement struct {
	orm.ModelBase `bson:"-",json:"-"`
	ID            bson.ObjectId ` bson:"_id" , json:"_id" `
	Subject       string
	Receivers     []string // list of userid
	AlarmCodes    []string // list of alarm code
	IntervalTime  int      // in minutes
	Template      string
	CreatedDate   time.Time
	LastUpdate    time.Time
	CreateBy      string // userid
	UpdateBy      string // userid
}

func (m *EmailManagement) New() *EmailManagement {
	m.ID = bson.NewObjectId()
	return m
}

func (m *EmailManagement) RecordID() interface{} {
	return m.ID
}

func (m *EmailManagement) TableName() string {
	return "EmailManagement"
}

type EmailCategory struct {
	orm.ModelBase `bson:"-",json:"-"`
	ID            string ` bson:"_id" , json:"_id" `
	Category      string
	Condition     string // isAlarmCode, isInterval
}

func (m *EmailCategory) RecordID() interface{} {
	return m.ID
}

func (m *EmailCategory) TableName() string {
	return "ref_emailCategory"
}
