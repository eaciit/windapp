package models

import (
	"fmt"

	"github.com/eaciit/orm"
)

type TurbineModel struct {
	orm.ModelBase  `bson:"-",json:"-"`
	Id             string ` bson:"_id" , json:"_id" `
	TurbineId      string
	TurbineName    string
	Feeder         string
	Project        string
	Latitude       float64
	Longitude      float64
	Elevation      float64
	Capacitymw     float64
	TopCorrelation []string
	Routine        string
}

func (m *TurbineModel) New() *TurbineModel {
	m.Id = fmt.Sprintf("%s_%s", m.TurbineId, m.Project)
	return m
}

func (m *TurbineModel) RecordID() interface{} {
	if m.Id == "" {
		m.Id = fmt.Sprintf("%s_%s", m.TurbineId, m.Project)
	}
	return m.Id
}

func (m *TurbineModel) TableName() string {
	return "ref_turbine"
}
