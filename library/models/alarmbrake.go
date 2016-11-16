package models

import (
	"github.com/eaciit/orm"
	"gopkg.in/mgo.v2/bson"
)

type AlarmBrake struct {
	orm.ModelBase              `bson:"-",json:"-"`
	ID                         bson.ObjectId ` bson:"_id" , json:"_id" `
	TypeCode                   int
	AlarmIndex                 int
	AlarmName                  string
	AlarmTypeId                string
	TypeId                     int
	Type                       string
	Set                        bool
	Disabled                   bool
	DefaultDisabled            bool
	BrakeProgram               int
	DefaultBrakeProgram        int
	YawProgram                 int
	DefaultYawProgram          int
	AlarmPaging                bool
	DefaultAlarmPaging         bool
	AlarmDelay                 int
	DefaultAlarmDelay          int
	AlarmDelayUnit             string
	ReducesAvailability        bool
	DefaultReducesAvailability bool
	OnTimeCounter              int
	AlarmCounter               int
	RepeatAlarmCode            int
	RepeatAlarmName            string
	RepeatAlarmNumber          int
	DefaultRepeatAlarmCounter  int
	RepeatAlarmTime            int
	DefaultRepeatAlarmTime     int
	LevelDisableAlarm          int
	LevelResetAlarm            int
}

func (m *AlarmBrake) New() *AlarmBrake {
	m.ID = bson.NewObjectId()
	return m
}

func (m *AlarmBrake) RecordID() interface{} {
	return m.ID
}

func (m *AlarmBrake) TableName() string {
	return "AlarmBrake"
}
