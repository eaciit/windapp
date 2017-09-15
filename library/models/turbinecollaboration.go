package models

import (
	"fmt"
	"time"

	. "eaciit/wfdemo-git/library/helper"

	"github.com/eaciit/orm"
	"gopkg.in/mgo.v2/bson"
)

type TurbineCollaborationModel struct {
	orm.ModelBase `bson:"-" json:"-"`
	Id            string ` bson:"_id" json:"_id" `
	ResponseFor   string
	ProjectId     string
	TurbineId     string
	TurbineName   string
	Feeder        string
	Date          time.Time
	DateInfo      DateInfo
	Status        string
	Remark        string
	CreatedBy     string
	CreatedByName string
	CreatedOn     time.Time
	CreatedIp     string
	CreatedLoc    string
	IsDeleted     bool
}

func (m *TurbineCollaborationModel) New() *TurbineCollaborationModel {
	if m.TurbineId != "" && m.CreatedBy != "" && !m.CreatedOn.IsZero() {
		sTime := m.CreatedOn.Format("20060102_150405")
		m.Id = fmt.Sprintf("%s_%s_%s", m.TurbineId, m.ProjectId, sTime)
	} else {
		m.Id = bson.NewObjectId().String()
	}
	return m
}

func (m *TurbineCollaborationModel) RecordID() interface{} {
	if m.Id == "" {
		sTime := m.CreatedOn.Format("20060102_150405")
		m.Id = fmt.Sprintf("%s_%s_%s", m.TurbineId, m.CreatedBy, sTime)
	}
	return m.Id
}

func (m *TurbineCollaborationModel) TableName() string {
	return "TurbineCollaboration"
}
