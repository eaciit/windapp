package gocore

import (
	"errors"

	"github.com/eaciit/dbox"
	_ "github.com/eaciit/dbox/dbc/csv"
	_ "github.com/eaciit/dbox/dbc/csvs"
	_ "github.com/eaciit/dbox/dbc/json"
	_ "github.com/eaciit/dbox/dbc/jsons"
	_ "github.com/eaciit/dbox/dbc/mongo"
	"github.com/eaciit/toolkit"
)

type queryWrapper struct {
	ci         *dbox.ConnectionInfo
	connection dbox.IConnection
	err        error
}

type MetaSave struct {
	keyword string
	data    string
}

func Query(driver string, host string, other ...interface{}) *queryWrapper {

	wrapper := queryWrapper{}
	wrapper.ci = &dbox.ConnectionInfo{host, "", "", "", nil}

	if len(other) > 0 {
		wrapper.ci.Database = other[0].(string)
	}
	if len(other) > 1 {
		wrapper.ci.UserName = other[1].(string)
	}
	if len(other) > 2 {
		wrapper.ci.Password = other[2].(string)
	}
	if len(other) > 3 {
		wrapper.ci.Settings = other[3].(toolkit.M)
	}

	wrapper.connection, wrapper.err = dbox.NewConnection(driver, wrapper.ci)
	if wrapper.err != nil {
		return &wrapper
	}

	wrapper.err = wrapper.connection.Connect()
	if wrapper.err != nil {
		return &wrapper
	}

	return &wrapper
}

func (c *queryWrapper) CheckIfConnected() error {
	return c.err
}

func (c *queryWrapper) Connect() (dbox.IConnection, error) {
	if c.err != nil {
		return nil, c.err
	}

	return c.connection, nil
}

func (c *queryWrapper) SelectOne(tableName string, clause ...*dbox.Filter) (toolkit.M, error) {
	if c.err != nil {
		return nil, c.err
	}

	connection := c.connection
	defer connection.Close()

	query := connection.NewQuery().Select().Take(1)
	if tableName != "" {
		query = query.From(tableName)
	}
	if len(clause) > 0 {
		query = query.Where(clause[0])
	}

	cursor, err := query.Cursor(nil)
	if err != nil {
		return nil, err
	}
	defer cursor.Close()

	data := make([]toolkit.M, 0)
	err = cursor.Fetch(&data, 0, false)
	if err != nil {
		return nil, err
	}

	if len(data) == 0 {
		return nil, errors.New("No data found")
	}

	return data[0], nil
}

func (c *queryWrapper) Delete(tableName string, clause *dbox.Filter) error {
	if c.err != nil {
		return c.err
	}

	connection := c.connection
	defer connection.Close()

	query := connection.NewQuery().Delete()
	if tableName != "" {
		query = query.From(tableName)
	}

	err := query.Where(clause).Exec(nil)
	if err != nil {
		return err
	}

	return nil
}

func (c *queryWrapper) SelectAll(tableName string, clause ...*dbox.Filter) ([]toolkit.M, error) {
	if c.err != nil {
		return nil, c.err
	}

	connection := c.connection
	defer connection.Close()

	query := connection.NewQuery().Select()
	if tableName != "" {
		query = query.From(tableName)
	}
	if len(clause) > 0 {
		query = query.Where(clause[0])
	}

	cursor, err := query.Cursor(nil)
	if err != nil {
		return nil, err
	}
	defer cursor.Close()

	data := make([]toolkit.M, 0)
	err = cursor.Fetch(&data, 0, false)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (c *queryWrapper) Save(tableName string, payload map[string]interface{}, clause ...*dbox.Filter) error {
	if c.err != nil {
		return c.err
	}

	connection := c.connection
	defer connection.Close()

	query := connection.NewQuery()
	if tableName != "" {
		query = query.From(tableName)
	}

	if len(clause) == 0 {
		err := query.Insert().Exec(toolkit.M{"data": payload})
		if err != nil {
			return err
		}

		return nil
	} else {
		err := query.Update().Where(clause[0]).Exec(toolkit.M{"data": payload})
		if err != nil {
			return err
		}

		return nil
	}

	return errors.New("nothing changes")
}

/*func FilterParse(where toolkit.M) *dbox.Filter {
	field := where.Get("field", "").(string)
	value := toolkit.Sprintf("%v", where["value"])

	if key := where.Get("key", "").(string); key == "Eq" {
		valueInt, errv := strconv.Atoi(toolkit.Sprintf("%v", where["value"]))
		if errv == nil {
			return dbox.Eq(field, valueInt)
		} else {
			return dbox.Eq(field, value)
		}
	} else if key == "Ne" {
		valueInt, errv := strconv.Atoi(toolkit.Sprintf("%v", where["value"]))
		if errv == nil {
			return dbox.Ne(field, valueInt)
		} else {
			return dbox.Ne(field, value)
		}
	} else if key == "Lt" {
		valueInt, errv := strconv.Atoi(toolkit.Sprintf("%v", where["value"]))
		if errv == nil {
			return dbox.Lt(field, valueInt)
		} else {
			return dbox.Lt(field, value)
		}
	} else if key == "Lte" {
		valueInt, errv := strconv.Atoi(toolkit.Sprintf("%v", where["value"]))
		if errv == nil {
			return dbox.Lte(field, valueInt)
		} else {
			return dbox.Lte(field, value)
		}
	} else if key == "Gt" {
		valueInt, errv := strconv.Atoi(toolkit.Sprintf("%v", where["value"]))
		if errv == nil {
			return dbox.Gt(field, valueInt)
		} else {
			return dbox.Gt(field, value)
		}
	} else if key == "Gte" {
		valueInt, errv := strconv.Atoi(toolkit.Sprintf("%v", where["value"]))
		if errv == nil {
			return dbox.Gte(field, valueInt)
		} else {
			return dbox.Gte(field, value)
		}
	} else if key == "In" {
		valueArray := []interface{}{}
		for _, e := range strings.Split(value, ",") {
			valueArray = append(valueArray, strings.Trim(e, ""))
		}
		return dbox.In(field, valueArray...)
	} else if key == "Nin" {
		valueArray := []interface{}{}
		for _, e := range strings.Split(value, ",") {
			valueArray = append(valueArray, strings.Trim(e, ""))
		}
		return dbox.Nin(field, valueArray...)
	} else if key == "Contains" {
		return dbox.Contains(field, value)
	} else if key == "Or" {
		subs := where.Get("value", []interface{}{}).([]interface{})
		filtersToMerge := []*dbox.Filter{}
		for _, eachSub := range subs {
			eachWhere, _ := toolkit.ToM(eachSub)
			filtersToMerge = append(filtersToMerge, FilterParse(eachWhere))
		}
		return dbox.Or(filtersToMerge...)
	} else if key == "And" {
		subs := where.Get("value", []interface{}{}).([]interface{})
		filtersToMerge := []*dbox.Filter{}
		for _, eachSub := range subs {
			eachWhere, _ := toolkit.ToM(eachSub)
			filtersToMerge = append(filtersToMerge, FilterParse(eachWhere))
		}
		return dbox.And(filtersToMerge...)
	}

	return nil
}*/
