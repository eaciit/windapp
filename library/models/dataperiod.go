package models

import (
	"github.com/eaciit/orm"
	"time"
)

type LatestDataPeriod struct {
	orm.ModelBase `bson:"-",json:"-"`
	Id            string ` bson:"_id" , json:"_id" `
	ProjectName   string
	Type          string
	Data          []time.Time
}

func (m *LatestDataPeriod) NewLatestDataPeriod() *LatestDataPeriod {
	m.Id = m.ProjectName + "_" + m.Type

	return m
}

func (m *LatestDataPeriod) TableName() string {
	return "LatestDataPeriod"
}
