package main

import (
	. "eaciit/wfdemo-git/library/models"
	"github.com/eaciit/dbox"
	_ "github.com/eaciit/dbox/dbc/jsons"
	_ "github.com/eaciit/dbox/dbc/mongo"
	"github.com/eaciit/toolkit"
	"gopkg.in/mgo.v2/bson"
	"os"
)

func main() {
	conn, _ := connectiontofile()
	query, _ := conn.NewQuery().From("project").Cursor(nil)
	dataJSON := make([]ProjectMaster, 0)
	e := query.Fetch(&dataJSON, 0, false)
	if e != nil {
		toolkit.Println(e.Error())
	}
	connM, _ := connectiontomongo()

	for _, data := range dataJSON {
		q := connM.NewQuery().From(new(ProjectMaster).TableName()).Save()
		data.ID = bson.ObjectId(toolkit.RandomString(12))
		e = q.Exec(toolkit.M{"data": data})
		if e != nil {
			toolkit.Println(e.Error())
		}
	}

}

func connectiontofile() (dbox.IConnection, error) {
	wd, _ := os.Getwd()
	ci := &dbox.ConnectionInfo{wd, "", "", "", nil}
	conn, err := dbox.NewConnection("jsons", ci)
	if err != nil {
		return nil, err
	}
	err = conn.Connect()
	return conn, nil
}

func connectiontomongo() (dbox.IConnection, error) {
	ci := &dbox.ConnectionInfo{"192.168.0.200:27123", "wfdemo", "", "", nil}
	conn, err := dbox.NewConnection("mongo", ci)
	if err != nil {
		return nil, err
	}
	err = conn.Connect()
	return conn, nil
}
