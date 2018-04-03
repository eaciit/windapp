package models

import (
	. "eaciit/wfdemo-git/library/helper"

	"github.com/eaciit/orm"
	"gopkg.in/mgo.v2/bson"
)

// type DGRScadaModel struct {
// 	orm.ModelBase       `bson:"-",json:"-"`
// 	Id                  bson.ObjectId ` bson:"_id" , json:"_id" `
// 	DateInfo            DateInfo
// 	ProjectName        	string
// 	Turbine             string
// 	PowerKW             float64
// 	Production          float64
// 	PCDeviation         float64
// 	Revenue         	float64
// 	RevenueInLacs       float64
// 	OkTime           	float64
// 	TrueAvail           float64
// 	ScadaAvail          float64
// 	MachineAvail        float64
// 	GridAvail           float64
// 	TotalAvail 			float64
// 	PLF        			float64
// 	Budget              float64
// 	AvgWindSpeed        float64
// 	ExpWindSpeed        float64
// 	DowntimeHours       float64
// 	LostEnergy    		float64
// 	RevenueLoss    		float64
// 	MachineDownHours    float64
// 	GridDownHours       float64
// 	otherdowntimehours  float64
// 	machinedownloss     float64
// 	griddownloss        float64
// 	otherdownloss       float64
// 	electricallosses    float64
// 	productionratio     float64
// 	nooffailures  		float64
// 	totalminutes		float64
// 	DetWindSpeed      []DetWindSpeed
// 	totalrows			float64
// }

type DGRScadaModel struct {
	orm.ModelBase      `bson:"-",json:"-"`
	Id                 bson.ObjectId ` bson:"_id" , json:"_id" `
	DateInfo           DateInfo
	ProjectName        string
	Turbine            string
	TurbineName        string
	PowerKW            float64
	Production         float64
	OkTime             float64
	TrueAvail          float64
	ScadaAvail         float64
	PLF                float64
	DowntimeHours      float64
	LostEnergy         float64
	MachineDownHours   float64
	GridDownHours      float64
	Otherdowntimehours float64
	LoWindTime         float64
}

func NewDGRScadaModel() *DGRScadaModel {
	m := new(DGRScadaModel)
	m.Id = bson.NewObjectId()
	return m
}

func (e *DGRScadaModel) RecordID() interface{} {
	return e.Id
}

func (m *DGRScadaModel) TableName() string {
	return "rpt_scadasummarydaily"
}
