package models

import (
	. "eaciit/wfdemo-git-dev/library/helper"

	"github.com/eaciit/orm"
)

type AlarmSummaryByMonth struct {
	orm.ModelBase `bson:"-",json:"-"`
	ID            string ` bson:"_id" , json:"_id" `
	DateInfo      DateInfo
	ProjectName   string
	Type          string
	LostEnergy    float64
}

func (m *AlarmSummaryByMonth) TableName() string {
	return "rpt_AlarmSummaryByMonth"
}
