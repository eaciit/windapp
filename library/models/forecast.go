package models

import (
	. "eaciit/wfdemo-git/library/helper"
	"time"

	tk "github.com/eaciit/toolkit"
)

// ForecastData for data modeling
type ForecastData struct {
	Id           string `bson:"_id" json:"_id"`
	DateReceived time.Time
	DateUpdated  time.Time
	Sender       string
	MailSubject  string
	ProjectName  string
	TimeStamp    time.Time
	DateInfo     DateInfo
	TimeRange    string
	TimeBlock    int
	AvgCapacity  float64
	SchCapacity  float64
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
