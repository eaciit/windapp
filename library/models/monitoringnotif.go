package models

import (
	"fmt"
	"time"

	"github.com/eaciit/orm"
)

type MonitoringNotification struct {
	orm.ModelBase `bson:"-",json:"-"`
	Id            string
	ProjectName   string
	Turbine       string
	GTags         string
	Tags          string
	TimeStart     time.Time
	TimeEnd       time.Time
	Duration      float64
	Status        bool
	IsError       bool
	CompareVal    float64
	Value         float64
	LastValue     float64
	NoteStart     string
	NoteEnd       string
}

func (m *MonitoringNotification) New() *MonitoringNotification {
	m.Id = fmt.Sprintf("%s_%s_%s_%s", m.ProjectName, m.Turbine, m.GTags, m.TimeStart.Format("20060102150405"))
	return m
}

func (m *MonitoringNotification) RecordID() interface{} {
	return fmt.Sprintf("%s_%s_%s_%s", m.ProjectName, m.Turbine, m.GTags, m.TimeStart.Format("20060102150405"))
}

func (m *MonitoringNotification) TableName() string {
	return "RealMonitoringNotification"
}
