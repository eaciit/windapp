package models

import (
	. "eaciit/wfdemo-git/library/helper"
	"fmt"
	"time"

	"github.com/eaciit/orm"
	"gopkg.in/mgo.v2/bson"
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
	CurrBudget50   float64
	CurrBudget90   float64
	CumProduction  float64
	CumBudget      float64
	CumBudget50    float64
	CumBudget90    float64
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

type DetailWindSpeed struct {
	SumWindSpeed   float64
	CountWindSpeed float64
}

type ScadaSummaryDaily struct {
	orm.ModelBase      `bson:"-",json:"-"`
	ID                 string ` bson:"_id" , json:"_id" `
	DateInfo           DateInfo
	ProjectName        string
	Turbine            string
	PowerKw            float64
	Production         float64 // it also called energy, measurement in kwh
	PCDeviation        float64
	Revenue            float64
	RevenueInLacs      float64
	OkTime             float64
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
	NoOfFailures       int
	TotalMinutes       int
	DetWindSpeed       DetailWindSpeed
	TotalRows          float64
	LoWindTime         float64
}

func (m *ScadaSummaryDaily) New() *ScadaSummaryDaily {
	m.ID = fmt.Sprintf("%s_%s_%s", m.ProjectName, m.Turbine, m.DateInfo.DateId.UTC().Format("20060102"))
	return m
}

func (m *ScadaSummaryDaily) RecordID() interface{} {
	if m.ID == "" {
		m.ID = fmt.Sprintf("%s_%s_%s", m.ProjectName, m.Turbine, m.DateInfo.DateId.UTC().Format("20060102"))
	}
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
	DataAvail     float64
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

type GWFAnalysisByProject struct {
	orm.ModelBase  `bson:"-",json:"-"`
	ID             bson.ObjectId ` bson:"_id" json:"_id" `
	ProjectName    string
	Key            string
	OrderNo        int
	Roll12Days     GWFAnalysisValue
	Roll12Weeks    GWFAnalysisValue
	Roll12Months   GWFAnalysisValue
	Roll12Quarters GWFAnalysisValue
}

func (m *GWFAnalysisByProject) New() *GWFAnalysisByProject {
	m.ID = bson.NewObjectId()
	return m
}

func (m *GWFAnalysisByProject) RecordID() interface{} {
	return m.ID
}

func (m *GWFAnalysisByProject) TableName() string {
	return "GWFAnalaysisByProject"
}

type GWFAnalysisByTurbine1 struct {
	orm.ModelBase  `bson:"-",json:"-"`
	ID             bson.ObjectId ` bson:"_id" json:"_id" `
	ProjectName    string
	Turbine        string
	Key            string
	OrderNo        int
	Roll12Days     GWFAnalysisValue
	Roll12Weeks    GWFAnalysisValue
	Roll12Months   GWFAnalysisValue
	Roll12Quarters GWFAnalysisValue
}

func (m *GWFAnalysisByTurbine1) New() *GWFAnalysisByTurbine1 {
	m.ID = bson.NewObjectId()
	return m
}

func (m *GWFAnalysisByTurbine1) RecordID() interface{} {
	return m.ID
}

func (m *GWFAnalysisByTurbine1) TableName() string {
	return "GWFAnalaysisByTurbine1"
}

type GWFAnalysisByTurbine2 struct {
	orm.ModelBase  `bson:"-",json:"-"`
	ID             bson.ObjectId ` bson:"_id" json:"_id" `
	ProjectName    string
	Key            string
	OrderNo        int
	Roll12Days     []GWFAnalysisItem2
	Roll12Weeks    []GWFAnalysisItem2
	Roll12Months   []GWFAnalysisItem2
	Roll12Quarters []GWFAnalysisItem2
}

func (m *GWFAnalysisByTurbine2) New() *GWFAnalysisByTurbine2 {
	m.ID = bson.NewObjectId()
	return m
}

func (m *GWFAnalysisByTurbine2) RecordID() interface{} {
	return m.ID
}

func (m *GWFAnalysisByTurbine2) TableName() string {
	return "GWFAnalaysisByTurbine2"
}

type GWFAnalysisValue struct {
	DateText string
	ValueAvg []GWFAnalysisItem
	ValueMin []GWFAnalysisItem // used only for turbine 2 analysis
	ValueMax []GWFAnalysisItem // used only for turbine 2 analysis
}

type GWFAnalysisItem struct {
	DataId  string
	Title   string
	OrderNo int
	Value   float64
}

type GWFAnalysisItem2 struct {
	Turbine string
	DataId  string
	Title   string
	OrderNo int
	Value   float64
}
