package models

import (
	. "eaciit/wfdemo-git/library/helper"
	"time"

	"math/rand"

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
	orm.ModelBase `bson:"-",json:"-"`
	ID            string ` bson:"_id" , json:"_id" `
	ProjectName   string
	Turbine       string
	TimeStart     time.Time
	TimeStartInt  int64
	// TimeStartUTC     time.Time
	DateInfoStart DateInfo
	// DateInfoStartUTC DateInfo
	TimeEnd time.Time
	// TimeEndUTC       time.Time
	TimeEndInt  int64
	DateInfoEnd DateInfo
	// DateInfoEndUTC   DateInfo
	AlarmDescription   string
	BrakeType          string // add by ams, regarding to add new req | 20170130
	Duration           float64
	Detail             []EventDownDetail
	DownGrid           bool
	DownEnvironment    bool
	DownMachine        bool
	ReduceAvailability bool
}

type EventDownDetail struct {
	TimeStamp    time.Time
	TimeStampInt int64
	// TimeStampUTC     time.Time
	DateInfo DateInfo
	// DateInfoUTC      DateInfo
	AlarmId          int
	AlarmDescription string
	BrakeType        string // add by ams, regarding to add new req | 20170130
	AlarmToggle      bool
}

func (m *EventDown) New() *EventDown {
	milistr := tk.ToString(m.TimeStart.Nanosecond() / 1000000)
	timeStampStr := m.TimeStart.Format("060102_150405") + "_" + milistr
	m.ID = timeStampStr + "#" + m.ProjectName + "#" + m.Turbine + "#" + time.Now().Format("060102150405_000000000") + "_" + tk.ToString(rand.Intn(999999))
	return m
}

func (m *EventDown) RecordID() interface{} {
	return m.ID
}

func (m *EventDown) TableName() string {
	return "EventDown"
}

type EventAlarm struct {
	orm.ModelBase `bson:"-",json:"-"`
	ID            string ` bson:"_id" , json:"_id" `
	ProjectName   string
	Turbine       string
	TimeStart     time.Time
	TimeStartInt  int64
	// TimeStartUTC     time.Time
	DateInfoStart DateInfo
	// DateInfoStartUTC DateInfo
	TimeEnd time.Time
	// TimeEndUTC       time.Time
	TimeEndInt  int64
	DateInfoEnd DateInfo
	// DateInfoEndUTC   DateInfo
	AlarmDescription string
	Duration         float64
	Detail           []EventAlarmDetail
	/*DownGrid         bool
	DownEnvironment  bool
	DownMachine      bool*/
}

type EventAlarmDetail struct {
	TimeStamp    time.Time
	TimeStampInt int64
	// TimeStampUTC     time.Time
	DateInfo DateInfo
	// DateInfoUTC      DateInfo
	AlarmId          int
	AlarmDescription string
	AlarmToggle      bool
}

func (m *EventAlarm) New() *EventAlarm {
	milistr := tk.ToString(m.TimeStart.Nanosecond() / 1000000)
	timeStampStr := m.TimeStart.Format("060102_150405") + "_" + milistr
	m.ID = timeStampStr + "#" + m.ProjectName + "#" + m.Turbine + "#" + time.Now().Format("060102150405_000000000") + "_" + tk.ToString(rand.Intn(999999))
	return m
}

func (m *EventAlarm) RecordID() interface{} {
	return m.ID
}

func (m *EventAlarm) TableName() string {
	return "EventAlarm"
}
