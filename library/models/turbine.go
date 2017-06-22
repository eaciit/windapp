package models

import (
	"github.com/eaciit/orm"
	"gopkg.in/mgo.v2/bson"
)

type TurbineModel struct {
	orm.ModelBase `bson:"-",json:"-"`
	Id            	bson.ObjectId ` bson:"_id" , json:"_id" `
	TurbineId       string
	TurbineName     string
	Feeder      	string
	Project         string
	Latitude		float64
	Longitude		float64
	Elevation		float64
	Capacitymw		float64
}

func (m *TurbineModel) New() *TurbineModel {
	m.Id = bson.NewObjectId()
	return m
}

func (m *TurbineModel) RecordID() interface{} {
	return m.Id
}

func (m *TurbineModel) TableName() string {
	return "ref_turbine"
}
