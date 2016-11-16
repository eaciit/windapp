package models

import (
	. "eaciit/wfdemo/library/helper"
	"time"

	"github.com/eaciit/orm"
	"gopkg.in/mgo.v2/bson"
)

type DowntimeEvent struct {
	orm.ModelBase    `bson:"-",json:"-"`
	ID               bson.ObjectId ` bson:"_id" , json:"_id" `
	ProjectName      string
	Turbine          string
	TimeStart        time.Time
	DateInfoStart    DateInfo
	TimeEnd          time.Time
	DateInfoEnd      DateInfo
	AlarmDescription string
	Duration         float64
	Detail           []DowntimeEventDetail
}

type DowntimeEventDetail struct {
	TimeStamp        time.Time
	DateInfo         DateInfo
	AlarmId          int
	AlarmDescription string
	AlarmToggle      bool
}

func (m *DowntimeEvent) New() *DowntimeEvent {
	m.ID = bson.NewObjectId()
	return m
}

func (m *DowntimeEvent) RecordID() interface{} {
	return m.ID
}

func (m *DowntimeEvent) TableName() string {
	return "DowntimeEvent"
}

type EventDown struct {
	orm.ModelBase    `bson:"-",json:"-"`
	ID               bson.ObjectId ` bson:"_id" , json:"_id" `
	ProjectName      string
	Turbine          string
	TimeStart        time.Time
	DateInfoStart    DateInfo
	TimeEnd          time.Time
	DateInfoEnd      DateInfo
	AlarmDescription string
	Duration         float64
	Detail           []EventDownDetail
}

type EventDownDetail struct {
	TimeStamp        time.Time
	DateInfo         DateInfo
	AlarmId          int
	AlarmDescription string
	AlarmToggle      bool
}

func (m *EventDown) New() *EventDown {
	m.ID = bson.NewObjectId()
	return m
}

func (m *EventDown) RecordID() interface{} {
	return m.ID
}

func (m *EventDown) TableName() string {
	return "EventDown"
}
