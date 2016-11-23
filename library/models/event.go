package models

import (
	. "eaciit/wfdemo-git/library/helper"
	"time"

	"github.com/eaciit/orm"
	tk "github.com/eaciit/toolkit"
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
	DownGrid         bool
	DownEnvironment  bool
	DownMachine      bool
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
	ID               string ` bson:"_id" , json:"_id" `
	ProjectName      string
	Turbine          string
	TimeStart        time.Time
	DateInfoStart    DateInfo
	TimeEnd          time.Time
	DateInfoEnd      DateInfo
	AlarmDescription string
	Duration         float64
	Detail           []EventDownDetail
	DownGrid         bool
	DownEnvironment  bool
	DownMachine      bool
}

type EventDownDetail struct {
	TimeStamp        time.Time
	DateInfo         DateInfo
	AlarmId          int
	AlarmDescription string
	AlarmToggle      bool
}

func (m *EventDown) New() *EventDown {
	milistr := tk.ToString(m.TimeStart.Nanosecond() / 1000000)
	timeStampStr := m.TimeStart.Format("060102_150405") + "_" + milistr
	m.ID = timeStampStr + "#" + m.ProjectName + "#" + m.Turbine
	return m
}

func (m *EventDown) RecordID() interface{} {
	return m.ID
}

func (m *EventDown) TableName() string {
	return "EventDown"
}
