package models

import (
	. "eaciit/wfdemo-git/library/helper"
	"time"

	"gopkg.in/mgo.v2/bson"
)

type Monitoring struct {
	ID                 bson.ObjectId ` bson:"_id" , json:"_id" `
	Timestamp          time.Time
	DateInfo           DateInfo
	LastUpdate         time.Time
	LastUpdateDateInfo DateInfo
	Project            string
	Turbine            string

	Production       float64
	WindSpeed        float64
	PerformanceIndex float64
	MachineAvail     float64
	GridAvail        float64

	IsAlarm   bool
	IsWarning bool
}

type MonitoringEvent struct {
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
