package models

import (
	. "eaciit/wfdemo-git/library/helper"
	"time"

	"github.com/eaciit/orm"
	tk "github.com/eaciit/toolkit"
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
	/*===================================*/
	IsAlarm   bool
	IsWarning bool

	Type       string // Alarm, Brake, Warning
	Status     string // ok, brake, N/A
	StatusCode int    // brake : AlarmID
	StatusDesc string // brake : AlarmDescription
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
	orm.ModelBase     `bson:"-",json:"-"`
	ID                string ` bson:"_id" , json:"_id" `
	Project           string
	Turbine           string
	TimeStamp         time.Time
	TimeStampStr      string
	DateInfo          DateInfo
	GroupTimeStamp    time.Time
	GroupTimeStampStr string
	AlarmId           int
	AlarmDescription  string
	Type              string // Alarm, Brake, Warning
	Status            string /// down, up
	Duration          float64
	PitchAngle        float64
	WindDirection     float64
}

func (m *MonitoringEvent) New() *MonitoringEvent {
	timestampstr := m.TimeStamp.Format("060102_150405")
	// nowstr := time.Now().Format("060102_150405")
	m.ID = m.Project + "#" + m.Turbine + "#" + timestampstr + "#" + tk.ToString(m.AlarmId) + "#" + m.Status + "#" + m.Type //+ "_" + nowstr
	return m
}

func (m *MonitoringEvent) RecordID() interface{} {
	return m.ID
}

func (m *MonitoringEvent) TableName() string {
	return "MonitoringEvent"
}
