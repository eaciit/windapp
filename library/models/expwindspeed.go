package models

import (
	"github.com/eaciit/orm"
)

type ExpectedWindSpeedModel struct {
	orm.ModelBase `bson:"-",json:"-"`
	ID            string ` bson:"_id" , json:"_id" `
	MonthNo       int
	ProjectId     string
	EngineId      string
	AvgWindSpeed  float64
	DataItems     []ExpectedWindSpeedItem
}

type ExpectedWindSpeedItem struct {
	Hour      string
	WindSpeed float64
}

func (m *ExpectedWindSpeedModel) New() *ExpectedWindSpeedModel {
	return m
}

func (m *ExpectedWindSpeedModel) RecordID() interface{} {
	return m.ID
}

func (m *ExpectedWindSpeedModel) TableName() string {
	return "ref_expwindspeed"
}
