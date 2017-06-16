package gocore

import (
	"errors"
	"os"

	"github.com/eaciit/acl/v1.0"
	"github.com/eaciit/dbox"
	"github.com/eaciit/toolkit"
)

var ConfigPath string

func GetConnectionInfo(db_type string) (string, *dbox.ConnectionInfo) {
	config := new(Configuration)

	config.ID = db_type
	db, err := config.GetDB()
	if err != nil {
		toolkit.Printf("Error get DB config: %s \n", err.Error())
	}

	setting, _ := toolkit.ToM(db["setting"])

	ci := &dbox.ConnectionInfo{
		db.GetString("host"),
		db.GetString("db"),
		db.GetString("user"),
		db.GetString("pass"),
		setting}

	return db.GetString("driver"), ci
}

func PrepareConnection(db_type string) (conn dbox.IConnection, err error) {
	driver, ci := GetConnectionInfo(db_type)
	conn, err = dbox.NewConnection(driver, ci)
	if err != nil {
		return
	}
	err = conn.Connect()
	return
}

func InitialSetDatabase() error {
	conn_acl, err := PrepareConnection(CONF_DB_ACL)
	if err != nil {
		return err
	}

	if err = acl.SetDb(conn_acl); err != nil {
		return err
	}

	conn_ostro, err := PrepareConnection(CONF_DB_OSTRO)
	if err != nil {
		return err
	}

	if err := SetDb(conn_ostro); err != nil {
		toolkit.Printf("Error set ostro database: %s \n", err.Error())
		return err
	}

	conn_reatime, err := PrepareConnection(CONF_DB_REALTIME)
	if err != nil {
		return err
	}

	if err := SetDbRealTime(conn_reatime); err != nil {
		toolkit.Printf("Error set realtime database: %s \n", err.Error())
		return err
	}

	return nil
}

func validateConfig() error {
	if ConfigPath == "" {
		return errors.New("gocore.validateConfig: ConfigPath is empty")
	}
	_, e := os.Stat(ConfigPath)
	if e != nil {
		return errors.New("gocore.validateConfig: " + e.Error())
	}
	return nil
}

func getConnection() (dbox.IConnection, error) {
	if e := validateConfig(); e != nil {
		return nil, errors.New("gocore.GetConnection: " + e.Error())
	}
	c, e := dbox.NewConnection("jsons", &dbox.ConnectionInfo{ConfigPath, "", "", "", toolkit.M{}.Set("newfile", true)})
	if e != nil {
		return nil, errors.New("gocore.GetConnection: " + e.Error())
	}
	e = c.Connect()
	if e != nil {
		return nil, errors.New("gocore.GetConnection: Connect: " + e.Error())
	}
	return c, nil
}
