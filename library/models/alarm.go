package models

import (
	. "eaciit/wfdemo-git/library/helper"
	"time"

	"github.com/eaciit/orm"
	"gopkg.in/mgo.v2/bson"
)

/*type AlarmRAW struct {
	orm.ModelBase    `bson:"-",json:"-"`
	ID               bson.ObjectId ` bson:"_id" , json:"_id" `
	Farm             string
	StartDate        time.Time
	StartDateInfo    DateInfo
	EndDate          time.Time
	Duration         float64
	Turbine          string
	AlertDescription string
	ExternalStop     bool
	GridDown         bool
	InternalGrid     bool
	MachineDown      bool
	AEbOK            bool
	Unknown          bool
	WeatherStop      bool
	Line             int
	ProjectName                   string
}

func (m *AlarmRAW) New() *AlarmRAW {
	m.ID = bson.NewObjectId()
	return m
}

func (m *AlarmRAW) RecordID() interface{} {
	return m.ID
}

func (m *AlarmRAW) TableName() string {
	return "AlarmRAW"
}*/

type AlarmOverlapping struct {
	orm.ModelBase    `bson:"-",json:"-"`
	ID               bson.ObjectId ` bson:"_id" , json:"_id" `
	Farm             string
	StartDate        time.Time
	StartDateInfo    DateInfo
	EndDate          time.Time
	Duration         float64
	Turbine          string
	AlertDescription string
	ExternalStop     bool
	GridDown         bool
	InternalGrid     bool
	MachineDown      bool
	AEbOK            bool
	Unknown          bool
	WeatherStop      bool
	Alarms           []Alarm
	ProjectName      string
}

func (m *AlarmOverlapping) New() *AlarmOverlapping {
	m.ID = bson.NewObjectId()
	return m
}

func (m *AlarmOverlapping) RecordID() interface{} {
	return m.ID
}

func (m *AlarmOverlapping) TableName() string {
	return "AlarmOverlapping"
}

type AlarmIncorrect struct {
	orm.ModelBase    `bson:"-",json:"-"`
	ID               bson.ObjectId ` bson:"_id" , json:"_id" `
	Farm             string
	StartDate        time.Time
	StartDateInfo    DateInfo
	EndDate          time.Time
	Turbine          string
	AlertDescription string
	ExternalStop     bool
	GridDown         bool
	InternalGrid     bool
	MachineDown      bool
	AEbOK            bool
	Unknown          bool
	WeatherStop      bool
	Line             int
	ProjectName      string
}

func (m *AlarmIncorrect) New() *AlarmIncorrect {
	m.ID = bson.NewObjectId()
	return m
}

func (m *AlarmIncorrect) RecordID() interface{} {
	return m.ID
}

func (m *AlarmIncorrect) TableName() string {
	return "AlarmIncorrect"
}

type Alarm struct {
	orm.ModelBase    `bson:"-",json:"-"`
	ID               bson.ObjectId ` bson:"_id" , json:"_id" `
	Farm             string
	StartDate        time.Time
	StartDateInfo    DateInfo
	EndDate          time.Time
	Duration         float64 // duration in hours
	Turbine          string
	AlertDescription string
	ExternalStop     bool
	GridDown         bool
	InternalGrid     bool
	MachineDown      bool
	AEbOK            bool
	Unknown          bool
	WeatherStop      bool
	Line             int
	ProjectName      string
	PowerLost        float64
}

func (m *Alarm) New() *Alarm {
	m.ID = bson.NewObjectId()
	return m
}

func (m *Alarm) RecordID() interface{} {
	return m.ID
}

func (m *Alarm) TableName() string {
	return "Alarm"
}

type AlarmScadaAnomaly struct {
	orm.ModelBase    `bson:"-",json:"-"`
	ID               bson.ObjectId ` bson:"_id" , json:"_id" `
	Farm             string
	StartDate        time.Time
	StartDateInfo    DateInfo
	EndDate          time.Time
	Duration         float64
	Turbine          string
	AlertDescription string
	ExternalStop     bool
	GridDown         bool
	InternalGrid     bool
	MachineDown      bool
	AEbOK            bool
	Unknown          bool
	WeatherStop      bool
	Line             int
	IsAlarmOk        bool
	ProjectName      string
	PowerLost        float64
}

func (m *AlarmScadaAnomaly) New() *AlarmScadaAnomaly {
	m.ID = bson.NewObjectId()
	return m
}

func (m *AlarmScadaAnomaly) RecordID() interface{} {
	return m.ID
}

func (m *AlarmScadaAnomaly) TableName() string {
	return "AlarmScadaAnomaly"
}

type AlarmClean struct {
	orm.ModelBase    `bson:"-",json:"-"`
	ID               bson.ObjectId ` bson:"_id" , json:"_id" `
	Farm             string
	StartDate        time.Time
	StartDateInfo    DateInfo
	EndDate          time.Time
	Duration         float64 // duration in hours
	Turbine          string
	AlertDescription string
	ExternalStop     bool
	GridDown         bool
	InternalGrid     bool
	MachineDown      bool
	AEbOK            bool
	Unknown          bool
	WeatherStop      bool
	Line             int
	ProjectName      string
	PowerLost        float64
}

func (m *AlarmClean) New() *AlarmClean {
	m.ID = bson.NewObjectId()
	return m
}

func (m *AlarmClean) RecordID() interface{} {
	return m.ID
}

func (m *AlarmClean) TableName() string {
	return "AlarmClean"
}
