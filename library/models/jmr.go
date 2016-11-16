package models

import (
	. "eaciit/wfdemo/library/helper"

	"github.com/eaciit/orm"
	"gopkg.in/mgo.v2/bson"
)

type JMR struct {
	orm.ModelBase `bson:"-",json:"-"`
	ID            bson.ObjectId ` bson:"_id" , json:"_id" `
	DateInfo      DateInfo
	Description   string
	Sections      []JMRSection
	TotalDetails  []JMRTotalDetails
}

type JMRSection struct {
	Turbine     string
	Description string
	Company     string
	ContrGen    float64

	BoEExport    float64
	BoEImport    float64
	BoENet       float64
	BoETotalLoss float64

	BoLExport float64
	BoLImport float64
	BoLNet    float64

	BoE2Export float64
	BoE2Import float64
	BoE2Net    float64
}

type JMRTotalDetails struct {
	Section        string
	ContrGenTotal  float64
	BoEExportTotal float64
	BoEImportTotal float64
	BoENetTotal    float64
	TotalLoss      float64

	/*BoLExportTotal float64
	BoLImportTotal float64
	BoLNetTotal    float64

	BoE2ExportTotal float64
	BoE2ImportTotal float64
	BoE2NetTotal    float64*/
}

func (m *JMR) New() *JMR {
	m.ID = bson.NewObjectId()
	return m
}

func (m *JMR) RecordID() interface{} {
	return m.ID
}

func (m *JMR) TableName() string {
	return "JMR"
}

func (m *JMR) SetTotalDetails() {
	result := []JMRTotalDetails{}

	for _, val := range m.Sections {
		found := false
	out:
		for idx, resVal := range result {
			if resVal.Section == val.Description {
				result[idx].ContrGenTotal = (resVal.ContrGenTotal + val.ContrGen)
				result[idx].BoEExportTotal = (resVal.BoEExportTotal + val.BoEExport)
				result[idx].BoEImportTotal = (resVal.BoEImportTotal + val.BoEImport)
				result[idx].BoENetTotal = (resVal.BoENetTotal + val.BoENet)
				result[idx].TotalLoss = (resVal.TotalLoss + val.BoETotalLoss)
				found = true
				break out
			}
		}

		if !found {
			tmpResult := JMRTotalDetails{}
			tmpResult.Section = val.Description
			tmpResult.ContrGenTotal = val.ContrGen
			tmpResult.BoEExportTotal = val.BoEExport
			tmpResult.BoEImportTotal = val.BoEImport
			tmpResult.BoENetTotal = val.BoENet
			tmpResult.TotalLoss = val.BoETotalLoss
			result = append(result, tmpResult)
		}
	}

	m.TotalDetails = result
}
