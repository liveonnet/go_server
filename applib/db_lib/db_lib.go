package db_lib

// https://godoc.org/github.com/go-sql-driver/mysql
// https://github.com/go-sql-driver/mysql/wiki/Examples
// https://github.com/go-sql-driver/mysql

import (
	"database/sql"
	"fmt"
	"log"

	"learn/go_server/applib/conf_lib"
	"learn/go_server/applib/log_lib"
	"learn/go_server/applib/tools_lib"

	_ "github.com/go-sql-driver/mysql"
	"gopkg.in/gin-gonic/gin.v1"
)

var logger *log.Logger
var conf map[string]interface{}
var m_db map[string]*sql.DB
var check func(e error, args ...string)

func init() {
	logger = log_lib.Log
	conf = conf_lib.Conf
	m_db = make(map[string]*sql.DB)
	check = tools_lib.Check
}

func QueryRows(db *sql.DB, close bool, sql_str string, args ...interface{}) (rslt [][]interface{}, col_name []string, err error) {
	if close {
		defer db.Close()
	}
	//	err = db.Ping()
	//	check(err, "Ping()", "")
	stmt, err := db.Prepare(sql_str)
	check(err, "Prepare", sql_str)
	defer stmt.Close()
	rcds, err := stmt.Query(args...)
	check(err, "query")
	defer rcds.Close()
	col_name, err = rcds.Columns()
	nr_cols := len(col_name)
	check(err, "columns")
	scanArgs := make([]interface{}, nr_cols)

	rslt = make([][]interface{}, 0, 10)

	i_row := 0
	for rcds.Next() {
		i_row += 1
		vals := make([]interface{}, nr_cols)
		for i := range vals {
			scanArgs[i] = &vals[i]
		}
		err := rcds.Scan(scanArgs...)
		check(err, "scan", string(i_row))
		rslt = append(rslt, vals)
	}

	return rslt, col_name, err
}

func QueryRowsRaw(db *sql.DB, close bool, sql_str string, args ...interface{}) (rslt [][]sql.RawBytes, col_name []string, err error) {
	if close {
		defer db.Close()
	}
	//	err = db.Ping()
	//	check(err, "Ping()", "")
	stmt, err := db.Prepare(sql_str)
	check(err, "Prepare", sql_str)
	defer stmt.Close()
	rcds, err := stmt.Query(args...)
	check(err, "query")
	defer rcds.Close()
	col_name, err = rcds.Columns()
	nr_cols := len(col_name)
	check(err, "columns")
	scanArgs := make([]interface{}, nr_cols)

	rslt = make([][]sql.RawBytes, 0, 10)

	i_row := 0
	for rcds.Next() {
		i_row += 1
		vals := make([]sql.RawBytes, nr_cols)
		for i := range vals {
			scanArgs[i] = &vals[i]
		}
		err := rcds.Scan(scanArgs...)
		check(err, "scan", string(i_row))
		rslt = append(rslt, vals)
	}

	return rslt, col_name, err
}

/*
func _GetDBOld(db_name string) (db *sql.DB, err error) {
	var ok bool
	if db, ok = m_db[db_name]; !ok {
		db_conf, _ := conf["database"]
		if db_conf, ok = db_conf.(map[string]interface{})[db_name]; ok {
			conn_str := fmt.Sprintf("%s:%s@tcp(%s:%v)/%s", db_conf.(map[string]interface{})["user"], db_conf.(map[string]interface{})["password"], db_conf.(map[string]interface{})["host"], db_conf.(map[string]interface{})["port"], db_conf.(map[string]interface{})["db_name"])
			logger.Printf("conn_str %s\n", conn_str)
			db, err = sql.Open("mysql", conn_str)
			db.SetMaxOpenConns(500)
			m_db[db_name] = db
		} else {
			logger.Printf("conf not found for %s\n", db_name)
		}
	} else {
		logger.Printf("use per process db conn %v\n", db)
		err = db.Ping()
	}

	return db, err
}
*/

func GetDB(db_name string) (db *sql.DB, err error) {
	var ok bool
	db_conf, _ := conf["database"]
	if db_conf, ok = db_conf.(map[string]interface{})[db_name]; ok {
		conn_str := fmt.Sprintf("%s:%s@tcp(%s:%v)/%s", db_conf.(map[string]interface{})["user"], db_conf.(map[string]interface{})["password"], db_conf.(map[string]interface{})["host"], db_conf.(map[string]interface{})["port"], db_conf.(map[string]interface{})["db_name"])
		logger.Printf("conn_str %s\n", conn_str)
		db, err = sql.Open("mysql", conn_str)
		db.SetMaxOpenConns(500)
		m_db[db_name] = db
	} else {
		err = logger.Output(1, fmt.Sprintf("conf not found for %s\n", db_name))
	}

	return db, err
}

type MyContext gin.Context

func (c *MyContext) GetDB(db_name string) (db *sql.DB, err error) {
	c_tmp := (*gin.Context)(c)
	var (
		db_map map[string]*sql.DB
		exist  bool
	)
	if v, exist := c_tmp.Get("db"); !exist {
		db_map = make(map[string]*sql.DB)
		c_tmp.Set("db", db_map)
	} else {
		db_map = v.(map[string]*sql.DB)
	}

	if db, exist = db_map[db_name]; !exist {
		if db, err = GetDB(db_name); err == nil {
			db_map[db_name] = db
		}
	} else {
		logger.Printf("use per request existing db conn %v\n", db)
	}
	return db, err
}
