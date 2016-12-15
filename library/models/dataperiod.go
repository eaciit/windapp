package models

import (
	"github.com/eaciit/orm"
	"gopkg.in/mgo.v2/bson"
	"time"
)

type LatestDataPeriod struct {
	orm.ModelBase `bson:"-",json:"-"`
	Id            bson.ObjectId ` bson:"_id" , json:"_id" `
	ProjectName   string
	Type          string
	Data          []time.Time
}

func NewLatestDataPeriod() *LatestDataPeriod {
	m := new(LatestDataPeriod)
	m.Id = bson.NewObjectId()

	return m
}

func (m *LatestDataPeriod) TableName() string {
	return "LatestDataPeriod"
}
