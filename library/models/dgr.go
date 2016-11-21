package models

import (
	. "eaciit/wfdemo-git/library/helper"
	"github.com/eaciit/orm"
	"gopkg.in/mgo.v2/bson"
	"time"
)

type DGRModel struct {
	orm.ModelBase       `bson:"-",json:"-"`
	Id                  bson.ObjectId ` bson:"_id" , json:"_id" `
	DateInfo            DateInfo
	CustomerName        string
	State               string
	Site                string
	Section             string
	Turbine             string
	MaxCapacity         float64
	GenKwhDay           float64
	GenKwhMtd           float64
	GenKwhYtd           float64
	PLFDay              float64
	PLFMtd              float64
	PLFYtd              float64
	MachineAvailability float64
	ForceMajeure        float64
	Schedule            float64
	Unschedule          float64
	NonOperational      float64
	GenerationHours     float64
	OperationalHours    float64
	GridAvailability    float64
	GFGF                float64
	GFFM                float64
	GFS                 float64
	GFU                 float64
	DowntimeHours       float64
	LostEnergy          float64
	RevenueLoss         float64
	GBILoss             float64
	DowntimeDetail      []DowntimeDetail
}

type DowntimeDetail struct {
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

func NewDGRModel() *DGRModel {
	m := new(DGRModel)
	m.Id = bson.NewObjectId()
	return m
}

func (e *DGRModel) RecordID() interface{} {
	return e.Id
}

func (m *DGRModel) TableName() string {
	return "rpt_generation"
}
