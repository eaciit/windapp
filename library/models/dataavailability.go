package models

import (
	. "eaciit/wfdemo-git/library/helper"

	"time"

	"github.com/eaciit/orm"
)

// should be applied for last 1 year datas
type DataAvailability struct {
	orm.ModelBase `bson:"-",json:"-"`
	ID            string ` bson:"_id" , json:"_id" `
	Type          string
	Name          string
	PeriodTo      time.Time
	PeriodFrom    time.Time
	Details       []DataAvailabilityDetail
}

type DataAvailabilityDetail struct {
	ProjectName string
	Turbine     string
	Start       time.Time
	StartInfo   DateInfo
	End         time.Time
	EndInfo     DateInfo
	Duration    float64
	IsAvail     bool
}

func (m *DataAvailability) TableName() string {
	return "DataAvailability"
}
