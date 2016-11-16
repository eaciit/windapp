package models

import (
	. "eaciit/wfdemo/library/helper"
	"github.com/eaciit/orm"
	"gopkg.in/mgo.v2/bson"
)

type WindRoseModel struct {
	orm.ModelBase    `bson:"-",json:"-"`
	ID               bson.ObjectId ` bson:"_id" , json:"_id" `
	DateInfo         DateInfo
	ProjectId        string
	TurbineId        string
	WindRoseItems    []WindRoseItem
	TotalContributes []WindRoseContribute
}

type WindRoseItem struct {
	DirectionNo    int
	DirectionDesc  string
	WSCategoryNo   int
	WSCategoryDesc string
	Contribute     float64
	Frequency      int
	Hours          float64
}

type WindRoseContribute struct {
	WSCategoryNo   int
	WSCategoryDesc string
	Contribute     float64
	Frequency      int
	Hours          float64
}

func (m *WindRoseModel) New() *WindRoseModel {
	m.ID = bson.NewObjectId()
	return m
}

func (m *WindRoseModel) RecordID() interface{} {
	return m.ID
}

func (m *WindRoseModel) TableName() string {
	return "rpt_scadawindrose"
}

type WindRoseMTModel struct {
	orm.ModelBase    `bson:"-",json:"-"`
	ID               bson.ObjectId ` bson:"_id" , json:"_id" `
	DateInfo         DateInfo
	ProjectId        string
	WindRoseItems    []WindRoseItemMT
	TotalContributes []WindRoseContributeMT
}

type WindRoseItemMT struct {
	DirectionNo    int
	DirectionDesc  string
	WSCategoryNo   int
	WSCategoryDesc string
	Contribute     float64
	Frequency      int
	Hours          float64
}

type WindRoseContributeMT struct {
	WSCategoryNo   int
	WSCategoryDesc string
	Contribute     float64
	Frequency      int
	Hours          float64
}

func (m *WindRoseMTModel) New() *WindRoseMTModel {
	m.ID = bson.NewObjectId()
	return m
}

func (m *WindRoseMTModel) RecordID() interface{} {
	return m.ID
}

func (m *WindRoseMTModel) TableName() string {
	return "rpt_scadawindrosemt"
}

type WindRoseNewModel struct {
	orm.ModelBase    `bson:"-",json:"-"`
	ID               bson.ObjectId ` bson:"_id" , json:"_id" `
	DateInfo         DateInfo
	ProjectId        string
	TurbineId        string
	WindRoseItems    []WindRoseItemNew
	TotalContributes []WindRoseContributeNew
}

type WindRoseItemNew struct {
	DirectionNo    int
	DirectionDesc  string
	WSCategoryNo   int
	WSCategoryDesc string
	Contribute     float64
	Frequency      int
	Hours          float64
}

type WindRoseContributeNew struct {
	WSCategoryNo   int
	WSCategoryDesc string
	Contribute     float64
	Frequency      int
	Hours          float64
}

func (m *WindRoseNewModel) New() *WindRoseNewModel {
	m.ID = bson.NewObjectId()
	return m
}

func (m *WindRoseNewModel) RecordID() interface{} {
	return m.ID
}

func (m *WindRoseNewModel) TableName() string {
	return "rpt_scadawindrosenew"
}
