package models

import (
	. "eaciit/wfdemo-git/library/helper"
	"time"

	"gopkg.in/mgo.v2/bson"

	"github.com/eaciit/orm"
	tk "github.com/eaciit/toolkit"
)

// ForecastData for data modeling
type ForecastData struct {
	orm.ModelBase   `bson:"-" json:"-"`
	Id              string `bson:"_id" json:"_id"`
	DateReceived    time.Time
	DateUpdated     time.Time
	Sender          string
	MailSubject     string
	ProjectName     string
	TimeStamp       time.Time
	DateInfo        DateInfo
	TimeRange       string
	TimeBlock       int
	AvgCapacity     float64
	SchCapacity     float64
	SchSdlc         float64
	WindSpeed       float64
	PowerRtd        float64
	TsRtd           time.Time
	MinCap          float64
	MaxCap          float64
	AvgCapacityMail float64
	SchCapacityMail float64
	IsEdited        int
}

//New instance for ForecastData
func (c *ForecastData) New() *ForecastData {
	c.Id = tk.Sprintf("%s_%v_%s", c.ProjectName, c.TimeBlock, c.TimeStamp.Format("20060102_150405"))
	return c
}

//TableName of ForecastData
func (c *ForecastData) TableName() string {
	return "ForecastData"
}

// ForecastConfig for data modeling
type ForecastConfig struct {
	orm.ModelBase `bson:"-" json:"-"`
	Id            string `bson:"_id" json:"_id"` // fill Id as projectname
	IsAutoSend    int
	AllowedUsers  []string
	LastSetAuto   time.Time
	LastSetBy     string
}

//New instance for ForecastConfig
func (c *ForecastConfig) New() *ForecastConfig {
	c.IsAutoSend = 0
	return c
}

//TableName of ForecastConfig
func (c *ForecastConfig) TableName() string {
	return "ForecastConfig"
}

// ForecastRecipients for data modeling
type ForecastRecipients struct {
	orm.ModelBase `bson:"-" json:"-"`
	Id            bson.ObjectId `bson:"_id" json:"_id"`
	ProjectName   string
	Email         string
	Name          string
	RecipientType string
}

//New instance for ForecastRecipients
func (c *ForecastRecipients) New() *ForecastRecipients {
	c.Id = bson.NewObjectId()
	return c
}

//RecordID to get id
func (c *ForecastRecipients) RecordID() interface{} {
	return c.Id
}

//TableName of ForecastRecipients
func (c *ForecastRecipients) TableName() string {
	return "ForecastRecipients"
}
