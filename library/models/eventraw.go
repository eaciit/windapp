package models

import (
	. "eaciit/wfdemo-git/library/helper"

	"math/rand"

	"time"

	"github.com/eaciit/orm"
	tk "github.com/eaciit/toolkit"
	"gopkg.in/mgo.v2/bson"
)

type DowntimeEventRaw struct {
	orm.ModelBase    `bson:"-",json:"-"`
	ID               bson.ObjectId ` bson:"_id" , json:"_id" `
	ProjectName      string
	Turbine          string
	TimeStamp        time.Time
	DateInfo         DateInfo
	EventType        string
	BrakeProgram     int
	AlarmDescription string
	AlarmId          int
	TurbineStatus    string
	AlarmToggle      bool
	BrakeType        string
}

func (m *DowntimeEventRaw) New() *DowntimeEventRaw {
	m.ID = bson.NewObjectId()
	return m
}

func (m *DowntimeEventRaw) RecordID() interface{} {
	return m.ID
}

func (m *DowntimeEventRaw) TableName() string {
	return "DowntimeEventRaw"
}

type EventRaw struct {
	orm.ModelBase `bson:"-",json:"-"`
	ID            string ` bson:"_id" , json:"_id" `
	ProjectName   string
	Turbine       string
	TimeStamp     time.Time
	TimeStampInt  int64
	// TimeStampUTC     time.Time
	DateInfo DateInfo
	// DateInfoUTC      DateInfo
	EventType        string
	BrakeProgram     int
	AlarmDescription string
	AlarmId          int
	TurbineStatus    string
	AlarmToggle      bool
	BrakeType        string
}

func (m *EventRaw) New() *EventRaw {
	milistr := tk.ToString(m.TimeStamp.Nanosecond() / 1000000)
	timeStampStr := m.TimeStamp.Format("060102_150405") + "_" + milistr
	m.ID = timeStampStr + "#" + m.ProjectName + "#" + m.Turbine + "#" + tk.ToString(m.AlarmId) + "#" + m.EventType + "#" + time.Now().Format("060102150405_000000000") + "_" + tk.ToString(rand.Intn(999999))
	return m
}

func (m *EventRaw) RecordID() interface{} {
	return m.ID
}

func (m *EventRaw) TableName() string {
	return "EventRaw"
}
