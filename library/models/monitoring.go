package models

import (
	. "eaciit/wfdemo-git/library/helper"
	"time"

	"github.com/eaciit/orm"
)

type Monitoring struct {
	orm.ModelBase      `bson:"-",json:"-"`
	ID                 string ` bson:"_id" , json:"_id" `
	TimeStamp          time.Time
	DateInfo           DateInfo
	LastUpdate         time.Time
	LastUpdateDateInfo DateInfo
	Project            string
	Turbine            string

	Production       float64 // MWh - Energy
	WindSpeed        float64
	PerformanceIndex float64 // skip
	MachineAvail     float64 // skip
	GridAvail        float64 // skip

	RotorSpeedRPM float64

	IsAlarm   bool
	IsWarning bool

	Status     string
	StatusCode string
	StatusDesc string
}

func (m *Monitoring) New() *Monitoring {
	timeStampStr := m.TimeStamp.Format("060102_150405")
	m.ID = m.Project + "#" + m.Turbine + "#" + timeStampStr
	return m
}

func (m *Monitoring) RecordID() interface{} {
	return m.ID
}

func (m *Monitoring) TableName() string {
	return "Monitoring"
}

type MonitoringEvent struct {
	orm.ModelBase    `bson:"-",json:"-"`
	ID               string ` bson:"_id" , json:"_id" `
	Project          string
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
	Type             string // Alarm, Brake, Warning
}

func (m *MonitoringEvent) New() *MonitoringEvent {
	timeStartStr := m.TimeStart.Format("060102_150405")
	m.ID = m.Project + "#" + m.Turbine + "#" + timeStartStr
	return m
}

func (m *MonitoringEvent) RecordID() interface{} {
	return m.ID
}

func (m *MonitoringEvent) TableName() string {
	return "MonitoringEvent"
}
