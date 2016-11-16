package main

import (
	. "github.com/eaciit/windapp/library/models"
	"github.com/eaciit/dbox"
	_ "github.com/eaciit/dbox/dbc/jsons"
	_ "github.com/eaciit/dbox/dbc/mongo"
	"github.com/eaciit/toolkit"
	"gopkg.in/mgo.v2/bson"
	"os"
	"time"
)

func main() {
	conn, _ := connectiontofile()
	query, _ := conn.NewQuery().From("windrose-data").Cursor(nil)
	dataJSON := make([]WindRoseItem, 0)
	turbineList := getTurbine()
	e := query.Fetch(&dataJSON, 0, false)
	if e != nil {
		toolkit.Println(e.Error())
	}
	connM, _ := connectiontomongo()

	for i, item := range dataJSON {
		q := connM.NewQuery().From(new(WindRoseModel).TableName()).Save()
		data := new(WindRoseModel)
		waktu := time.Now()
		data.ID = bson.ObjectId(toolkit.RandomString(12))
		data.DateInfo.DateId = waktu.AddDate(0, 0, -i)
		data.ProjectId = "Tejuva"
		data.TurbineId = turbineList[i%len(turbineList)]
		item.Hours = item.Contribute * 2.4
		data.WindRoseItems = append(data.WindRoseItems, item)

		e = q.Exec(toolkit.M{"data": data})
		if e != nil {
			toolkit.Println(e.Error())
		}
	}

}

func getTurbine() []string {
	connM, _ := connectiontomongo()
	csr, e := connM.NewQuery().From("ref_turbine").Cursor(nil)

	if e != nil {
		toolkit.Println(e.Error())
	}
	defer csr.Close()

	data := []toolkit.M{}
	e = csr.Fetch(&data, 0, false)

	result := []string{}

	for _, val := range data {
		result = append(result, val.GetString("turbineid"))
	}

	if e != nil {
		toolkit.Println(e.Error())
	}

	return result
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
	ci := &dbox.ConnectionInfo{"localhost:27017", "wfdemo", "", "", nil}
	conn, err := dbox.NewConnection("mongo", ci)
	if err != nil {
		return nil, err
	}
	err = conn.Connect()
	return conn, nil
}
