package models

import (
	. "eaciit/wfdemo/library/helper"
	"github.com/eaciit/orm"
	"gopkg.in/mgo.v2/bson"
	"time"
)

type DGRDowntimeModel struct {
	orm.ModelBase        `bson:"-",json:"-"`
	Id                   bson.ObjectId ` bson:"_id" , json:"_id" `
	DateInfo             DateInfo
	CustomerName         string
	State                string
	Site                 string
	Section              string
	Turbine              string
	MaxCapacity          float64
	BreakdownRemark      string
	FormulaParameter     string
	BreakdownHours       float64
	StartTime            time.Time
	EndTime              time.Time
	WindSpeed            float64
	GenerationPowerCurve float64
	LostEnergy           float64
	RevenueLoss          float64
	GBILoss              float64
	BreakdownNote        string
}

func NewDGRDowntimeModel() *DGRDowntimeModel {
	m := new(DGRDowntimeModel)
	m.Id = bson.NewObjectId()
	return m
}

func (e *DGRDowntimeModel) RecordID() interface{} {
	return e.Id
}

func (m *DGRDowntimeModel) TableName() string {
	return "rpt_downtime"
}
