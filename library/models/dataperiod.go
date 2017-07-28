package models

import (
	"fmt"
	"time"

	"github.com/eaciit/orm"
)

type LatestDataPeriod struct {
	orm.ModelBase `bson:"-",json:"-"`
	Id            string ` bson:"_id" , json:"_id" `
	ProjectName   string
	Type          string
	Data          []time.Time
}

func (m *LatestDataPeriod) New() *LatestDataPeriod {
	m.Id = fmt.Sprintf("%s_%s", m.ProjectName, m.Type)
	return m
}

func (m *LatestDataPeriod) RecordID() interface{} {
	return fmt.Sprintf("%s_%s", m.ProjectName, m.Type)
}

func (m *LatestDataPeriod) TableName() string {
	return "LatestDataPeriod"
}
