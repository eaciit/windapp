package models

import (
	. "eaciit/wfdemo-git/library/helper"

	"github.com/eaciit/orm"
	"gopkg.in/mgo.v2/bson"
)

type PowerCurveModel struct {
	orm.ModelBase `bson:"-",json:"-"`
	ID            bson.ObjectId ` bson:"_id" , json:"_id" `
	Model         string
	WindSpeed     float64
	Power1        float64
	Standard      float64
	Engine        string
}

func (m *PowerCurveModel) New() *PowerCurveModel {
	m.ID = bson.NewObjectId()
	return m
}

func (m *PowerCurveModel) RecordID() interface{} {
	return m.ID
}

func (m *PowerCurveModel) TableName() string {
	return "ref_powercurve"
}

var SPCTableName string

type ScadaPowerCurveModel struct {
	orm.ModelBase `bson:"-",json:"-"`
	ID            bson.ObjectId ` bson:"_id" , json:"_id" `
	DateInfo      DateInfo
	ProjectName   string
	TurbineId     string
	DataItems     []ScadaPowerCurveItem
}

type ScadaPowerCurveItem struct {
	WSClass    float64
	Production float64
	TotalData  int
}

func (m *ScadaPowerCurveModel) New() *ScadaPowerCurveModel {
	m.ID = bson.NewObjectId()
	return m
}

func (m *ScadaPowerCurveModel) RecordID() interface{} {
	return m.ID
}

func (m *ScadaPowerCurveModel) SetTableName(tblname string) {
	SPCTableName = tblname
}

func (m *ScadaPowerCurveModel) TableName() string {
	if SPCTableName == "" {
		SPCTableName = "rpt_scadapowercurve"
	}
	return SPCTableName
}
