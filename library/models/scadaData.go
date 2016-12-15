package models

import (
	. "eaciit/wfdemo-git/library/helper"
	"time"

	"github.com/eaciit/orm"
	"gopkg.in/mgo.v2/bson"
)

type ScadaData struct {
	orm.ModelBase             `bson:"-",json:"-"`
	ID                        bson.ObjectId ` bson:"_id" , json:"_id" `
	DateInfo                  DateInfo
	TimeStamp                 time.Time
	Turbine                   string
	GridFrequency             float64
	ReactivePower             float64
	AlarmExtStopTime          float64
	AlarmGridDownTime         float64
	AlarmInterLineDown        float64
	AlarmMachDownTime         float64
	AlarmOkTime               float64
	AlarmUnknownTime          float64
	AlarmWeatherStop          float64
	ExternalStopTime          float64
	GridDownTime              float64
	GridOkSecs                float64
	InternalLineDown          float64
	MachineDownTime           float64
	OkSecs                    float64
	OkTime                    float64
	UnknownTime               float64
	WeatherStopTime           float64
	GeneratorRPM              float64
	NacelleYawPositionUntwist float64
	NacelleTemperature        float64
	AdjWindSpeed              float64
	AmbientTemperature        float64
	AvgBladeAngle             float64
	AvgWindSpeed              float64
	UnitsGenerated            float64
	EstimatedPower            float64
	EstimatedEnergy           float64 // new added on Sep 14, 2016 by ams
	NacelDirection            float64
	Power                     float64
	PowerLost                 float64
	Energy                    float64 // new added on Sep 14, 2016 by ams
	EnergyLost                float64 // new added on Sep 14, 2016 by ams
	RotorRPM                  float64
	WindDirection             float64
	Line                      int
	IsValidTimeDuration       bool
	TotalTime                 float64
	Minutes                   int
	ProjectName               string
	Available                 int
	DenValue                  float64 // new added on Sep 14, 2016 by ams
	DenPh                     float64 // new added on Sep 14, 2016 by ams
	DenWindSpeed              float64 // new added on Sep 14, 2016 by ams
	DenAdjWindSpeed           float64 // new added on Sep 14, 2016 by ams
	DenPower                  float64 // new added on Sep 14, 2016 by ams
	DenEnergy                 float64 // new added on Sep 14, 2016 by ams
	PCValue                   float64 // new added on Sep 15, 2016 by ams
	PCValueAdj                float64 // new added on Sep 15, 2016 by ams
	PCDeviation               float64 // new added on Sep 15, 2016 by ams
	WSAdjForPC                float64 // new added on Sep 16, 2016 by ams
	WSAvgForPC                float64 // new added on Sep 16, 2016 by ams
	TotalAvail                float64 // new added on Sep 27, 2016 by ams
	MachineAvail              float64 // new added on Sep 27, 2016 by ams
	GridAvail                 float64 // new added on Sep 27, 2016 by ams
	DenPcDeviation            float64 // new added on Sep 27, 2016 by ams
	DenDeviationPct           float64 // new added on Sep 27, 2016 by ams
	DenPcValue                float64 // new added on Sep 27, 2016 by ams
	DeviationPct              float64 // new added on Sep 27, 2016 by ams
	MTTR                      float64
	MTTF                      float64
	PerformanceIndex          float64
}

func (m *ScadaData) New() *ScadaData {
	m.ID = bson.NewObjectId()
	return m
}

func (m *ScadaData) RecordID() interface{} {
	return m.ID
}

func (m *ScadaData) TableName() string {
	return "ScadaData"
}

type ScadaAlarmAnomaly struct {
	orm.ModelBase             `bson:"-",json:"-"`
	ID                        bson.ObjectId ` bson:"_id" , json:"_id" `
	DateInfo                  DateInfo
	TimeStamp                 time.Time
	Turbine                   string
	GridFrequency             float64
	ReactivePower             float64
	AlarmExtStopTime          float64
	AlarmGridDownTime         float64
	AlarmInterLineDown        float64
	AlarmMachDownTime         float64
	AlarmOkTime               float64
	AlarmUnknownTime          float64
	AlarmWeatherStop          float64
	ExternalStopTime          float64
	GridDownTime              float64
	GridOkSecs                float64
	InternalLineDown          float64
	MachineDownTime           float64
	OkSecs                    float64
	OkTime                    float64
	UnknownTime               float64
	WeatherStopTime           float64
	GeneratorRPM              float64
	NacelleYawPositionUntwist float64
	NacelleTemperature        float64
	AdjWindSpeed              float64
	AmbientTemperature        float64
	AvgBladeAngle             float64
	AvgWindSpeed              float64
	UnitsGenerated            float64
	EstimatedPower            float64
	NacelDirection            float64
	Power                     float64
	PowerLost                 float64
	RotorRPM                  float64
	WindDirection             float64
	Line                      int
	IsValidTimeDuration       bool
	TotalTime                 float64
	Minutes                   int
	ProjectName               string
	Available                 int
}

func (m *ScadaAlarmAnomaly) New() *ScadaAlarmAnomaly {
	m.ID = bson.NewObjectId()
	return m
}

func (m *ScadaAlarmAnomaly) RecordID() interface{} {
	return m.ID
}

func (m *ScadaAlarmAnomaly) TableName() string {
	return "ScadaAlarmAnomaly"
}

type ScadaClean struct {
	orm.ModelBase             `bson:"-",json:"-"`
	ID                        bson.ObjectId ` bson:"_id" , json:"_id" `
	DateInfo                  DateInfo
	TimeStamp                 time.Time
	Turbine                   string
	GridFrequency             float64
	ReactivePower             float64
	AlarmExtStopTime          float64
	AlarmGridDownTime         float64
	AlarmInterLineDown        float64
	AlarmMachDownTime         float64
	AlarmOkTime               float64
	AlarmUnknownTime          float64
	AlarmWeatherStop          float64
	ExternalStopTime          float64
	GridDownTime              float64
	GridOkSecs                float64
	InternalLineDown          float64
	MachineDownTime           float64
	OkSecs                    float64
	OkTime                    float64
	UnknownTime               float64
	WeatherStopTime           float64
	GeneratorRPM              float64
	NacelleYawPositionUntwist float64
	NacelleTemperature        float64
	AdjWindSpeed              float64
	AmbientTemperature        float64
	AvgBladeAngle             float64
	AvgWindSpeed              float64
	UnitsGenerated            float64
	EstimatedPower            float64
	NacelDirection            float64
	Power                     float64
	PowerLost                 float64
	RotorRPM                  float64
	WindDirection             float64
	Line                      int
	IsValidTimeDuration       bool
	TotalTime                 float64
	Minutes                   int
	ProjectName               string
	Available                 int
}

func (m *ScadaClean) New() *ScadaClean {
	m.ID = bson.NewObjectId()
	return m
}

func (m *ScadaClean) RecordID() interface{} {
	return m.ID
}

func (m *ScadaClean) TableName() string {
	return "ScadaClean"
}

type ScadaDataNew struct {
	orm.ModelBase             `bson:"-",json:"-"`
	ID                        bson.ObjectId ` bson:"_id" , json:"_id" `
	DateInfo                  DateInfo
	TimeStamp                 time.Time
	Turbine                   string
	GridFrequency             float64
	ReactivePower             float64
	AlarmExtStopTime          float64
	AlarmGridDownTime         float64
	AlarmInterLineDown        float64
	AlarmMachDownTime         float64
	AlarmOkTime               float64
	AlarmUnknownTime          float64
	AlarmWeatherStop          float64
	ExternalStopTime          float64
	GridDownTime              float64
	GridOkSecs                float64
	InternalLineDown          float64
	MachineDownTime           float64
	OkSecs                    float64
	OkTime                    float64
	UnknownTime               float64
	WeatherStopTime           float64
	GeneratorRPM              float64
	NacelleYawPositionUntwist float64
	NacelleTemperature        float64
	AdjWindSpeed              float64
	AmbientTemperature        float64
	AvgBladeAngle             float64
	AvgWindSpeed              float64
	UnitsGenerated            float64
	EstimatedPower            float64
	EstimatedEnergy           float64 // new added on Sep 14, 2016 by ams
	NacelDirection            float64
	Power                     float64
	PowerLost                 float64
	Energy                    float64 // new added on Sep 14, 2016 by ams
	EnergyLost                float64 // new added on Sep 14, 2016 by ams
	RotorRPM                  float64
	WindDirection             float64
	Line                      int
	IsValidTimeDuration       bool
	TotalTime                 float64
	Minutes                   int
	ProjectName               string
	Available                 int
	DenValue                  float64 // new added on Sep 14, 2016 by ams
	DenPh                     float64 // new added on Sep 14, 2016 by ams
	DenWindSpeed              float64 // new added on Sep 14, 2016 by ams
	DenAdjWindSpeed           float64 // new added on Sep 14, 2016 by ams
	DenPower                  float64 // new added on Sep 14, 2016 by ams
	DenEnergy                 float64 // new added on Sep 14, 2016 by ams
	PCValue                   float64 // new added on Sep 15, 2016 by ams
	PCDeviation               float64 // new added on Sep 15, 2016 by ams
	WSAdjForPC                float64 // new added on Sep 16, 2016 by ams
	WSAvgForPC                float64 // new added on Sep 16, 2016 by ams
}

func (m *ScadaDataNew) New() *ScadaDataNew {
	m.ID = bson.NewObjectId()
	return m
}

func (m *ScadaDataNew) RecordID() interface{} {
	return m.ID
}

func (m *ScadaDataNew) TableName() string {
	return "ScadaDataNew"
}
