package gocore

import (
	"errors"

	"github.com/eaciit/dbox"
	_ "github.com/eaciit/dbox/dbc/jsons"
	_ "github.com/eaciit/dbox/dbc/mongo"
	"github.com/eaciit/orm"
	"github.com/eaciit/toolkit"
	"gopkg.in/mgo.v2"
	"time"
)

var _ctx *orm.DataContext
var _ctxErr error
var _db *orm.DataContext
var _conn dbox.IConnection
var _connrealtime dbox.IConnection
var _dbErr error
var _session *mgo.Session

func ctx() *orm.DataContext {
	var _conn dbox.IConnection
	var econn error
	if _ctx == nil {
		if _conn == nil {
			_conn, econn = getConnection()
			if econn != nil {
				_ctxErr = errors.New("Connection error: " + econn.Error())
				return nil
			}
		}
		_ctx = orm.New(_conn)
	}
	return _ctx
}

func Delete(o orm.IModel) error {
	e := DB().Delete(o)
	if e != nil {
		return errors.New("Delete: " + e.Error())
	}
	return e

}

func Save(o orm.IModel) error {
	e := DB().Save(o)
	if e != nil {
		return errors.New("Save: " + e.Error())
	}
	return e
}

func Get(o orm.IModel, id interface{}) error {
	filter := dbox.Eq("_id", id)
	e := DB().Get(o, toolkit.M{}.Set(orm.ConfigWhere, filter))
	if e != nil {
		return errors.New("Get: " + e.Error())
	}
	return e
}

func Find(o orm.IModel, filter *dbox.Filter, config toolkit.M) (dbox.ICursor, error) {
	var filters []*dbox.Filter
	if filter != nil {
		filters = append(filters, filter)
	}

	dconf := toolkit.M{}.Set("where", filters)
	if config != nil {
		if config.Has("take") {
			dconf.Set("limit", config["take"])
		}
		if config.Has("skip") {
			dconf.Set("skip", config["skip"])
		}
		if config.Has("order") && toolkit.TypeName(config["order"]) == "[]string" {
			dconf.Set("order", config["order"])
		}
	}

	c, e := DB().Find(o, dconf)
	if e != nil {
		return nil, errors.New("Find: " + e.Error())
	}
	return c, nil
}

func GetData(o orm.IModel, id interface{}) error {
	filter := dbox.Eq("_id", id)
	e := ctx().Get(o, toolkit.M{}.Set(orm.ConfigWhere, filter))
	if e != nil {
		return errors.New("Core.Get: " + e.Error())
	}
	return e
}

func SetDb(conn dbox.IConnection) error {
	CloseDb()

	e := conn.Connect()
	if e != nil {
		_dbErr = errors.New("ostro.SetDB: Test Connect: " + e.Error())
		return _dbErr
	}

	_db = orm.New(conn)
	return nil
}

func CloseDb() {
	if _db != nil {
		_db.Close()
	}
}

func DB() *orm.DataContext {
	return _db
}

func SetDbRealTime(conn dbox.IConnection) error {
	e := conn.Connect()
	if e != nil {
		_dbErr = errors.New("ostro.SetDB: Test Connect: " + e.Error())
		return _dbErr
	}

	_connrealtime = conn
	return nil
}

func DBRealtime() dbox.IConnection {
	return _connrealtime
}

func SetSession(conn dbox.IConnection) (e error) {
	CloseSession()

	dboxInfo := conn.Info()
	mgoInfo := new(mgo.DialInfo)
	mgoInfo.Addrs = []string{dboxInfo.Host}
	mgoInfo.Database = dboxInfo.Database
	mgoInfo.Username = dboxInfo.UserName
	mgoInfo.Password = dboxInfo.Password
	if dboxInfo.Settings == nil {
		dboxInfo.Settings = toolkit.M{}
	}
	poollimit := dboxInfo.Settings.GetInt("poollimit")
	if poollimit > 0 {
		mgoInfo.PoolLimit = poollimit
	}
	timeout := dboxInfo.Settings.GetInt("timeout")
	if timeout > 0 {
		mgoInfo.Timeout = time.Duration(timeout) * time.Second
	}
	_session, e = mgo.DialWithInfo(mgoInfo)
	if e != nil {
		return
	}
	_session.SetMode(mgo.Monotonic, true)

	return
}

func DBSession() *mgo.Session {
	return _session
}

func CloseSession() {
	if _session != nil {
		_session.Close()
	}
}
