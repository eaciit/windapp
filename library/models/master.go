package models

import (
	"github.com/eaciit/orm"
	"github.com/eaciit/toolkit"
	"gopkg.in/mgo.v2/bson"
)

type TurbineMaster struct {
	orm.ModelBase  `bson:"-" json:"-"`
	ID             bson.ObjectId ` bson:"_id" json:"_id" `
	TurbineId      string
	TurbineName    string
	Feeder         string
	Project        string
	Latitude       float64
	Longitude      float64
	Elevation      float64
	CapacityMW     float64
	Routine        string
	TotalTurbine   int
	Active         bool
	TopCorrelation []string
	Engine         string
	Cluster        float64
	ProjectDgr     string
	TurbineDgr     string
}

func (m *TurbineMaster) New() *TurbineMaster {
	m.ID = bson.NewObjectId()
	return m
}

func (m *TurbineMaster) RecordID() interface{} {
	return m.ID
}

func (m *TurbineMaster) TableName() string {
	return "ref_turbine"
}

type StrangethresholdMaster struct {
	ID          string ` bson:"_id" json:"_id" `
	Max         float64
	Min         float64
	ProjectName []string
	Tags        string
	Type        string
}

func (m *StrangethresholdMaster) New() *StrangethresholdMaster {
	m.ID = m.Type + "_" + m.Tags
	return m
}

func (m *StrangethresholdMaster) RecordID() interface{} {
	return m.ID
}

func (m *StrangethresholdMaster) TableName() string {
	return "ref_strangethreshold"
}

type ProjectMaster struct {
	orm.ModelBase          `bson:"-" json:"-"`
	ID                     bson.ObjectId ` bson:"_id" json:"_id" `
	ProjectId              string
	ProjectName            string
	TotalPower             float64
	Latitude               float64
	Longitude              float64
	TotalTurbine           int
	RevenueMultiplier      float64
	State                  string
	City                   string
	SS_AirDensity          float64
	STD_AirDensity         float64
	Engine                 []string
	Forecast_Min_Cap       float64
	Forecast_Max_Cap       float64
	Forecast_Revision_Info []toolkit.M
}

func (m *ProjectMaster) New() *ProjectMaster {
	m.ID = bson.NewObjectId()
	return m
}

func (m *ProjectMaster) RecordID() interface{} {
	return m.ID
}

func (m *ProjectMaster) TableName() string {
	return "ref_project"
}

type TurbineOut struct {
	Project    string
	Turbine    string
	Value      string
	Capacity   float64
	Feeder     string
	Engine     string
	Coords     []float64
	Cluster    float64
	DgrProject string
	DgrTurbine string
}

type ProjectOut struct {
	ProjectId         string
	Name              string
	Value             string
	Coords            []float64
	RevenueMultiplier float64
	State             string
	City              string
	NoOfTurbine       int
	TotalMaxCapacity  float64
	SS_AirDensity     float64
	STD_AirDensity    float64
	Engine            []string
	ForecastMinCap    float64
	ForecastMaxCap    float64
	ForecastRevInfos  []toolkit.M
}

type StrangethresholdOld struct {
	StrangethresholdId string
	Max                float64
	Min                float64
	ProjectName        []string
	Tags               string
	Type               string
}
