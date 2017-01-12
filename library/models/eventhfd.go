package models

import (
	. "eaciit/wfdemo-git/library/helper"

	"math/rand"

	"time"

	"github.com/eaciit/orm"
	tk "github.com/eaciit/toolkit"
)

type EventRawHFD struct {
	orm.ModelBase `bson:"-",json:"-"`
	ID            string ` bson:"_id" , json:"_id" `
	ProjectName   string
	Turbine       string
	TimeStamp     time.Time
	DateInfo      DateInfo

	EventType        string
	BrakeProgram     int    // alarmbrake > brakeprogram
	AlarmDescription string // alarmbrake > alarmame
	AlarmId          int    // alarmcode

	// TurbineStatus string
	// AlarmToggle   bool
	BrakeType string // AlarmBrake > type
}

func (m *EventRawHFD) New() *EventRawHFD {
	timeStampStr := m.TimeStamp.Format("060102_150405")
	m.ID = timeStampStr + "#" + m.ProjectName + "#" + m.Turbine + "#" + tk.ToString(m.AlarmId) + "#" + m.EventType + "#" + time.Now().Format("060102150405_000000000") + "_" + tk.ToString(rand.Intn(999999))
	return m
}

func (m *EventRawHFD) TableName() string {
	return "EventRawHFD"
}

type EventDownHFD struct {
	orm.ModelBase    `bson:"-",json:"-"`
	ID               string ` bson:"_id" , json:"_id" `
	ProjectName      string
	Turbine          string
	TimeStart        time.Time
	DateInfoStart    DateInfo
	TimeEnd          time.Time
	TimeEndInt       int64
	DateInfoEnd      DateInfo
	AlarmDescription string
	Duration         float64
	DownGrid         bool
	DownEnvironment  bool
	DownMachine      bool
}

func (m *EventDownHFD) New() *EventDownHFD {
	timeStampStr := m.TimeStart.Format("060102_150405")
	m.ID = timeStampStr + "#" + m.ProjectName + "#" + m.Turbine + "#" + time.Now().Format("060102150405_000000000") + "_" + tk.ToString(rand.Intn(999999))
	return m
}

func (m *EventDownHFD) RecordID() interface{} {
	return m.ID
}

func (m *EventDownHFD) TableName() string {
	return "EventDownHFD"
}
