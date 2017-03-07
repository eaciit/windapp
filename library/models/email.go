package models

import (
	// . "eaciit/wfdemo-git/library/helper"
	"time"

	"github.com/eaciit/orm"
)

type EmailManagement struct {
	orm.ModelBase `bson:"-",json:"-"`
	ID            string ` bson:"_id" , json:"_id" `
	Subject       string
	Category      string   // refer to ref_emailCategory
	Receivers     []string // list of userid
	AlarmCodes    []string // list of alarm code
	IntervalTime  int      // in minutes
	Template      string
	CreatedDate   time.Time
	LastUpdate    time.Time
	CreatedBy     string // userid
	UpdatedBy     string // userid
	Enable        bool
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
	Condition     string // "isAlarmCode,isInterval"
}

func (m *EmailCategory) RecordID() interface{} {
	return m.ID
}

func (m *EmailCategory) TableName() string {
	return "ref_emailCategory"
}
