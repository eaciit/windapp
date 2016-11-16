package models

import (
	"github.com/eaciit/orm"
	"gopkg.in/mgo.v2/bson"
)

type UserModel struct {
	orm.ModelBase `bson:"-",json:"-"`
	Id            bson.ObjectId ` bson:"_id" , json:"_id" `
	UserId        string
	UserName      string
	Password      string
	Email         string
	DataAccess    map[string][]string
}

func NewUserModel() *UserModel {
	m := new(UserModel)
	m.Id = bson.NewObjectId()
	return m
}

func (e *UserModel) RecordID() interface{} {
	return e.Id
}

func (m *UserModel) TableName() string {
	return "SysUsers"
}
