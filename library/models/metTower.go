package models

import (
	. "eaciit/ostrowfm/library/helper"
	"time"

	"github.com/eaciit/orm"
	"gopkg.in/mgo.v2/bson"
)

type MetTower struct {
	orm.ModelBase `bson:"-",json:"-"`
	ID            bson.ObjectId ` bson:"_id" , json:"_id" `
	Line          int
	TimeStamp     time.Time
	DateInfo      DateInfo

	VHubWS90mAvg    float64
	VHubWS90mMax    float64
	VHubWS90mMin    float64
	VHubWS90mStdDev float64
	VHubWS90mCount  float64

	VRefWS88mAvg    float64
	VRefWS88mMax    float64
	VRefWS88mMin    float64
	VRefWS88mStdDev float64
	VRefWS88mCount  float64

	VTipWS42mAvg    float64
	VTipWS42mMax    float64
	VTipWS42mMin    float64
	VTipWS42mStdDev float64
	VTipWS42mCount  float64

	DHubWD88mAvg    float64
	DHubWD88mMax    float64
	DHubWD88mMin    float64
	DHubWD88mStdDev float64
	DHubWD88mCount  float64

	DRefWD86mAvg    float64
	DRefWD86mMax    float64
	DRefWD86mMin    float64
	DRefWD86mStdDev float64
	DRefWD86mCount  float64

	THubHHubHumid855mAvg    float64
	THubHHubHumid855mMax    float64
	THubHHubHumid855mMin    float64
	THubHHubHumid855mStdDev float64
	THubHHubHumid855mCount  float64

	TRefHRefHumid855mAvg    float64
	TRefHRefHumid855mMax    float64
	TRefHRefHumid855mMin    float64
	TRefHRefHumid855mStdDev float64
	TRefHRefHumid855mCount  float64

	THubHHubTemp855mAvg    float64
	THubHHubTemp855mMax    float64
	THubHHubTemp855mMin    float64
	THubHHubTemp855mStdDev float64
	THubHHubTemp855mCount  float64

	TRefHRefTemp855mAvg    float64
	TRefHRefTemp855mMax    float64
	TRefHRefTemp855mMin    float64
	TRefHRefTemp855mStdDev float64
	TRefHRefTemp855mCount  float64

	BaroAirPress855mAvg    float64
	BaroAirPress855mMax    float64
	BaroAirPress855mMin    float64
	BaroAirPress855mStdDev float64
	BaroAirPress855mCount  float64

	WindDirNo      int    // added by ams, Sep 19, 2016
	WindDirDesc    string // added by ams, Sep 19, 2016
	WSCategoryNo   int    // added by ams, Sep 19, 2016
	WSCategoryDesc string // added by ams, Sep 19, 2016

	YawAngleVoltageAvg          float64
	YawAngleVoltageMax          float64
	YawAngleVoltageMin          float64
	YawAngleVoltageStdDev       float64
	YawAngleVoltageCount        float64
	OtherSensorVoltageAI1Avg    float64
	OtherSensorVoltageAI1Max    float64
	OtherSensorVoltageAI1Min    float64
	OtherSensorVoltageAI1StdDev float64
	OtherSensorVoltageAI1Count  float64
	OtherSensorVoltageAI2Avg    float64
	OtherSensorVoltageAI2Max    float64
	OtherSensorVoltageAI2Min    float64
	OtherSensorVoltageAI2StdDev float64
	OtherSensorVoltageAI2Count  float64
	OtherSensorVoltageAI3Avg    float64
	OtherSensorVoltageAI3Max    float64
	OtherSensorVoltageAI3Min    float64
	OtherSensorVoltageAI3StdDev float64
	OtherSensorVoltageAI3Count  float64
	OtherSensorVoltageAI4Avg    float64
	OtherSensorVoltageAI4Max    float64
	OtherSensorVoltageAI4Min    float64
	OtherSensorVoltageAI4StdDev float64
	OtherSensorVoltageAI4Count  float64
	GenRPMCurrentAvg            float64
	GenRPMCurrentMax            float64
	GenRPMCurrentMin            float64
	GenRPMCurrentStdDev         float64
	GenRPMCurrentCount          float64
	WS_SCSCurrentAvg            float64
	WS_SCSCurrentMax            float64
	WS_SCSCurrentMin            float64
	WS_SCSCurrentStdDev         float64
	WS_SCSCurrentCount          float64
	RainStatusCount             float64
	RainStatusSum               float64
	OtherSensor2StatusIO1Avg    float64
	OtherSensor2StatusIO1Max    float64
	OtherSensor2StatusIO1Min    float64
	OtherSensor2StatusIO1StdDev float64
	OtherSensor2StatusIO1Count  float64
	OtherSensor2StatusIO2Avg    float64
	OtherSensor2StatusIO2Max    float64
	OtherSensor2StatusIO2Min    float64
	OtherSensor2StatusIO2StdDev float64
	OtherSensor2StatusIO2Count  float64
	OtherSensor2StatusIO3Avg    float64
	OtherSensor2StatusIO3Max    float64
	OtherSensor2StatusIO3Min    float64
	OtherSensor2StatusIO3StdDev float64
	OtherSensor2StatusIO3Count  float64
	OtherSensor2StatusIO4Avg    float64
	OtherSensor2StatusIO4Max    float64
	OtherSensor2StatusIO4Min    float64
	OtherSensor2StatusIO4StdDev float64
	OtherSensor2StatusIO4Count  float64
	OtherSensor2StatusIO5Avg    float64
	OtherSensor2StatusIO5Max    float64
	OtherSensor2StatusIO5Min    float64
	OtherSensor2StatusIO5StdDev float64
	OtherSensor2StatusIO5Count  float64
	A1Avg                       float64
	A1Max                       float64
	A1Min                       float64
	A1StdDev                    float64
	A1Count                     float64
	A2Avg                       float64
	A2Max                       float64
	A2Min                       float64
	A2StdDev                    float64
	A2Count                     float64
	A3Avg                       float64
	A3Max                       float64
	A3Min                       float64
	A3StdDev                    float64
	A3Count                     float64
	A4Avg                       float64
	A4Max                       float64
	A4Min                       float64
	A4StdDev                    float64
	A4Count                     float64
	A5Avg                       float64
	A5Max                       float64
	A5Min                       float64
	A5StdDev                    float64
	A5Count                     float64
	A6Avg                       float64
	A6Max                       float64
	A6Min                       float64
	A6StdDev                    float64
	A6Count                     float64
	A7Avg                       float64
	A7Max                       float64
	A7Min                       float64
	A7StdDev                    float64
	A7Count                     float64
	A8Avg                       float64
	A8Max                       float64
	A8Min                       float64
	A8StdDev                    float64
	A8Count                     float64
	A9Avg                       float64
	A9Max                       float64
	A9Min                       float64
	A9StdDev                    float64
	A9Count                     float64
	A10Avg                      float64
	A10Max                      float64
	A10Min                      float64
	A10StdDev                   float64
	A10Count                    float64
	AC1Avg                      float64
	AC1Max                      float64
	AC1Min                      float64
	AC1StdDev                   float64
	AC1Count                    float64
	AC2Avg                      float64
	AC2Max                      float64
	AC2Min                      float64
	AC2StdDev                   float64
	AC2Count                    float64
	C1Avg                       float64
	C1Max                       float64
	C1Min                       float64
	C1StdDev                    float64
	C1Count                     float64
	C2Avg                       float64
	C2Max                       float64
	C2Min                       float64
	C2StdDev                    float64
	C2Count                     float64
	C3Avg                       float64
	C3Max                       float64
	C3Min                       float64
	C3StdDev                    float64
	C3Count                     float64
	D1Avg                       float64
	D1Max                       float64
	D1Min                       float64
	D1StdDev                    float64
	M1_1Avg                     float64
	M1_1Max                     float64
	M1_1Min                     float64
	M1_1StdDev                  float64
	M1_1Count                   float64
	M1_2Avg                     float64
	M1_2Max                     float64
	M1_2Min                     float64
	M1_2StdDev                  float64
	M1_2Count                   float64
	M1_3Avg                     float64
	M1_3Max                     float64
	M1_3Min                     float64
	M1_3StdDev                  float64
	M1_3Count                   float64
	M1_4Avg                     float64
	M1_4Max                     float64
	M1_4Min                     float64
	M1_4StdDev                  float64
	M1_4Count                   float64
	M1_5Avg                     float64
	M1_5Max                     float64
	M1_5Min                     float64
	M1_5StdDev                  float64
	M1_5Count                   float64
	M2_1Avg                     float64
	M2_1Max                     float64
	M2_1Min                     float64
	M2_1StdDev                  float64
	M2_1Count                   float64
	M2_2Avg                     float64
	M2_2Max                     float64
	M2_2Min                     float64
	M2_2StdDev                  float64
	M2_2Count                   float64
	M2_3Avg                     float64
	M2_3Max                     float64
	M2_3Min                     float64
	M2_3StdDev                  float64
	M2_3Count                   float64
	M2_4Avg                     float64
	M2_4Max                     float64
	M2_4Min                     float64
	M2_4StdDev                  float64
	M2_4Count                   float64
	M2_5Avg                     float64
	M2_5Max                     float64
	M2_5Min                     float64
	M2_5StdDev                  float64
	M2_5Count                   float64
	M2_6Avg                     float64
	M2_6Max                     float64
	M2_6Min                     float64
	M2_6StdDev                  float64
	M2_6Count                   float64
	M2_7Avg                     float64
	M2_7Max                     float64
	M2_7Min                     float64
	M2_7StdDev                  float64
	M2_7Count                   float64
	M2_8Avg                     float64
	M2_8Max                     float64
	M2_8Min                     float64
	M2_8StdDev                  float64
	M2_8Count                   float64
	VAvg                        float64
	VMax                        float64
	VMin                        float64
	IAvg                        float64
	IMax                        float64
	IMin                        float64
	T                           float64
	Addr                        float64
}

func (m *MetTower) New() *MetTower {
	m.ID = bson.NewObjectId()
	return m
}

func (m *MetTower) RecordID() interface{} {
	return m.ID
}

func (m *MetTower) TableName() string {
	return "MetTower"
}
