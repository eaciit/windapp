package models

import (
	"time"

	"github.com/eaciit/orm"
)

type UserPreferences struct {
	orm.ModelBase  `bson:"-",json:"-"`
	Id             string ` bson:"_id" , json:"_id" ` // the id will be the user id
	LoginID        string
	KPIAnalysis    []KPIAnalysis
	AnalysisStudio []AnalysisStudio
}

func (e *UserPreferences) RecordID() interface{} {
	return e.Id
}

func (m *UserPreferences) TableName() string {
	return "UserPreferences"
}

type KPIAnalysis struct {
	Name            string
	KeyA            string
	KeyB            string
	KeyC            string
	ColumnBreakdown string
	RowBreakdown    string
}

type AnalysisStudio struct {
	Name    string
	Keys    []string
	Filters []Filter
}

type Filter struct {
	Project   string
	Turbine   []string
	Period    string
	DateStart time.Time
	DateEnd   time.Time
}
