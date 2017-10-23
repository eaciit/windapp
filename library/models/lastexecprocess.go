package models

import (
	"fmt"
	"time"

	"github.com/eaciit/orm"
)

type LastExecProcess struct {
	orm.ModelBase `bson:"-" json:"-"`
	ID            string `json:"_id" bson:"_id"`
	Process       string //conv10min
	Type          string //null, alarm, realtime
	ProjectName   string
	LastDate      time.Time
}

func (m *LastExecProcess) New() *LastExecProcess {
	m.ID = fmt.Sprintf("%s_%s_%s", m.Process, m.Type, m.ProjectName)
	return m
}

func (m *LastExecProcess) RecordID() interface{} {
	return fmt.Sprintf("%s_%s_%s", m.Process, m.Type, m.ProjectName)
}

func (m *LastExecProcess) TableName() string {
	return "log_lastexecprocess"
}
