package models

import (
	"github.com/eaciit/orm"
)

type ExpPValueModel struct {
	orm.ModelBase      `bson:"-",json:"-"`
	ID                 string ` bson:"_id" , json:"_id" `
	MonthNo            int
	EnergyDistribution float64
	P50NetGenMWH       float64
	P50Plf             float64
	P75NetGenMWH       float64
	P75Plf             float64
	P90NetGenMWH       float64
	P90Plf             float64
}

func (m *ExpPValueModel) New() *ExpPValueModel {
	return m
}

func (m *ExpPValueModel) RecordID() interface{} {
	return m.ID
}

func (m *ExpPValueModel) TableName() string {
	return "ref_exp_p_value"
}
