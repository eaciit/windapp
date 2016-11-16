package models

import (
	. "eaciit/wfdemo/library/helper"
	"github.com/eaciit/orm"
	"gopkg.in/mgo.v2/bson"
	"time"
)

type ScadaLastUpdate struct {
	orm.ModelBase          `bson:"-",json:"-"`
	ID                     string ` bson:"_id" , json:"_id" `
	LastUpdate             time.Time
	DateInfo               DateInfo
	ProjectName            string
	NoOfProjects           int
	NoOfTurbines           int
	TotalMaxCapacity       float64
	CurrentDown            int
	TwoDaysDown            int
	Productions            []LastData24Hours
	CummulativeProductions []Last30Days
}

type LastData24Hours struct {
	Hour         int
	TimeHour     time.Time
	PowerKw      float64
	EnergyKwh    float64
	Potential    float64
	PotentialKwh float64
	TrueAvail    float64
	GridAvail    float64
	AvgWindSpeed float64
}

type Last30Days struct {
	DayNo          int
	DateId         time.Time
	CurrProduction float64
	CurrBudget     float64
	CumProduction  float64
	CumBudget      float64
}

func (m *ScadaLastUpdate) New() *ScadaLastUpdate {
	m.ID = "SCADALASTUPDATE"
	return m
}

func (m *ScadaLastUpdate) RecordID() interface{} {
	return m.ID
}

func (m *ScadaLastUpdate) TableName() string {
	return "rpt_scadalastupdate"
}

type ScadaSummaryByMonth struct {
	orm.ModelBase      `bson:"-",json:"-"`
	ID                 bson.ObjectId ` bson:"_id" , json:"_id" `
	DateInfo           DateInfo
	ProjectName        string
	Production         float64
	ProductionLastYear float64
	Revenue            float64
	RevenueInLacs      float64
	TrueAvail          float64
	ScadaAvail         float64
	MachineAvail       float64
	GridAvail          float64
	PLF                float64
	Budget             float64
	AvgWindSpeed       float64
	ExpWindSpeed       float64
	DowntimeHours      float64
	LostEnergy         float64
	RevenueLoss        float64
}

func (m *ScadaSummaryByMonth) New() *ScadaSummaryByMonth {
	m.ID = bson.NewObjectId()
	return m
}

func (m *ScadaSummaryByMonth) RecordID() interface{} {
	return m.ID
}

func (m *ScadaSummaryByMonth) TableName() string {
	return "rpt_scadasummarybymonth"
}

type ScadaSummaryDaily struct {
	orm.ModelBase      `bson:"-",json:"-"`
	ID                 bson.ObjectId ` bson:"_id" , json:"_id" `
	DateInfo           DateInfo
	ProjectName        string
	Turbine            string
	PowerKw            float64
	Production         float64 // it also called energy, measurement in kwh
	PCDeviation        float64
	Revenue            float64
	RevenueInLacs      float64
	TrueAvail          float64
	ScadaAvail         float64
	MachineAvail       float64
	GridAvail          float64
	TotalAvail         float64
	PLF                float64
	Budget             float64
	AvgWindSpeed       float64
	ExpWindSpeed       float64
	DowntimeHours      float64
	LostEnergy         float64
	RevenueLoss        float64
	MachineDownHours   float64
	GridDownHours      float64
	OtherDowntimeHours float64
	MachineDownLoss    float64
	GridDownLoss       float64
	OtherDownLoss      float64
	ElectricalLosses   float64
	ProductionRatio    float64
}

func (m *ScadaSummaryDaily) New() *ScadaSummaryDaily {
	m.ID = bson.NewObjectId()
	return m
}

func (m *ScadaSummaryDaily) RecordID() interface{} {
	return m.ID
}

func (m *ScadaSummaryDaily) TableName() string {
	return "rpt_scadasummarydaily"
}

type ScadaSummaryByProject struct {
	orm.ModelBase `bson:"-",json:"-"`
	ID            string ` bson:"_id" , json:"_id" `
	DataItems     []ScadaSummaryByProjectItem
}

type ScadaSummaryByProjectItem struct {
	Name          string
	NoOfWtg       int
	Production    float64
	PLF           float64
	LostEnergy    float64
	DowntimeHours float64
	MachineAvail  float64
	TrueAvail     float64
}

func (m *ScadaSummaryByProject) New() *ScadaSummaryByProject {
	return m
}

func (m *ScadaSummaryByProject) RecordID() interface{} {
	return m.ID
}

func (m *ScadaSummaryByProject) TableName() string {
	return "rpt_scadasummarybyproject"
}
