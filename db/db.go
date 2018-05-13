package db

import (
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/last911/tools"
	"reflect"
	"strconv"
	"strings"
)

// Row interfacer{}
type Row interface{}

// MapRow map[string]string
type MapRow map[string]string

// Rows rows type
type Rows []Row

// DB struct
type DB struct {
	*sqlx.DB
}

// NewDB return DB struct
// dsn [username[:password]@][protocol[(address)]]/dbname[?param1=value1&...&paramN=valueN]
// Examples:
// root:scnjl@tcp(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=true
func NewDB(driver, dsn string) (*DB, error) {
	sqldb, err := sql.Open(driver, dsn)
	if err != nil {
		return nil, err
	}

	return &DB{sqlx.NewDb(sqldb, driver)}, nil
}

// Execute execute sql
// returns last insert id、affect row, error
func (db *DB) Execute(sql string, v ...interface{}) (int64, int64, error) {
	result, err := db.Exec(sql, v...)
	if err != nil {
		return 0, 0, err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return 0, 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, 0, err
	}

	return id, rows, nil
}

// query return num rows
func (db *DB) query(sql string, num int, v ...interface{}) (Rows, error) {
	var st interface{}
	if len(v) > 0 {
		stValue := reflect.ValueOf(v[0])
		if reflect.TypeOf(stValue).Kind() == reflect.Struct {
			st = v[0]
			v = v[1:]
		}
	}

	rows, err := db.Queryx(sql, v...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result Rows
	if st != nil {
		for rows.Next() {
			err := rows.StructScan(st)
			if err != nil {
				return nil, err
			}
			result = append(result, reflect.ValueOf(st).Elem().Interface())
			if num == 1 {
				break
			}
		}
	} else {
		// get column's name
		columns, err := rows.Columns()
		if err != nil {
			return nil, err
		}
		var cLen = len(columns)
		valuePtrs := make([]interface{}, cLen)
		values := make([][]byte, cLen)
		for i := 0; i < cLen; i++ {
			valuePtrs[i] = &values[i]
		}
		for rows.Next() {
			err := rows.Scan(valuePtrs...)
			if err != nil {
				return nil, err
			}
			row := make(MapRow, cLen)
			for i, v := range columns {
				row[v] = string(values[i])
			}

			result = append(result, row)
			if num == 1 {
				break
			}
		}
	}

	return result, nil
}

// FetchAll fetch all rows
func (db *DB) FetchAll(sql string, v ...interface{}) (Rows, error) {
	return db.query(sql, 0, v...)
}

// FetchOne fetch one row
func (db *DB) FetchOne(sql string, v ...interface{}) (Row, error) {
	rows, err := db.query(sql, 1, v...)
	if err != nil {
		return nil, err
	}

	if len(rows) > 0 {
		return rows[0], nil
	}
	return nil, nil
}

// Count return count
func (db *DB) Count(sql string, v ...interface{}) (int64, error) {
	row, err := db.FetchOne(sql, v...)
	if err != nil {
		return 0, err
	}

	for _, v := range row.(MapRow) {
		i, err := strconv.ParseInt(v, 10, 0)
		if err != nil {
			return 0, err
		}
		return i, nil
	}

	return 0, fmt.Errorf("Count error: no fields")
}

// doExec handle insert/replace action
func (db *DB) doExec(handle, table string, data ...map[string]interface{}) (int64, error) {
	dataLen := len(data)
	if dataLen > 0 {
		fields := tools.Keys(data[0])
		fieldLen := len(fields)
		sql := handle + " INTO " + table + "(" + strings.Join(fields, ", ") + ") VALUES " +
			strings.Repeat(", ("+strings.Repeat(", ?", fieldLen)[2:]+")", dataLen)[2:]
		values := make([]interface{}, dataLen*fieldLen)
		var i int
		for _, field := range fields {
			for _, row := range data {
				values[i] = row[field]
				i++
			}
		}

		id, affectRows, err := db.Execute(sql, values...)
		if err != nil {
			return 0, err
		}

		if handle == "INSERT" {
			return id, nil
		}
		return affectRows, nil
	}

	return 0, fmt.Errorf("Execute SQL[" + handle + "] not enough parameters")
}

// Insert insert to table
func (db *DB) Insert(table string, data ...map[string]interface{}) (int64, error) {
	return db.doExec("INSERT", table, data...)
}

// Replace replace to table
func (db *DB) Replace(table string, data ...map[string]interface{}) (int64, error) {
	return db.doExec("REPLACE", table, data...)
}
