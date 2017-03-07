package db_lib

import (
	"fmt"
	"reflect"
	"testing"
)

func Test_Query(t *testing.T) {
	if db, err := GetDB("uc_read"); err == nil {
		sql_str := "select uid, device_id, ctime, os_type, pnum from z_user where left(uid, 1) = '4' limit 10"
		rslt, col_name, err := QueryRows(db, false, sql_str)
		check(err, "query", sql_str)
		// find max column name length
		max_len := 0
		for _, _name := range col_name {
			if len(_name) > max_len {
				max_len = len(_name)
			}
		}
		s := fmt.Sprintf("%%+%ds: %%v %%v\n", max_len)

		for i, cols := range rslt {
			t.Logf("-------------- %v --------------------\n", i+1)
			for j, col := range cols {
				t.Logf(s, col_name[j], reflect.TypeOf(col), col)
			}
		}
	} else {
		fmt.Printf("err %s\n", err)
	}

}

func Test_QueryRaw(t *testing.T) {
	if db, err := GetDB("uc_read"); err == nil {
		sql_str := "select uid, device_id, ctime, os_type, pnum from z_user where left(uid, 1) = '4' limit 10"
		rslt, col_name, err := QueryRowsRaw(db, false, sql_str)
		check(err, "query", sql_str)
		// find max column name length
		max_len := 0
		for _, _name := range col_name {
			if len(_name) > max_len {
				max_len = len(_name)
			}
		}
		s := fmt.Sprintf("%%+%ds: %%s\n", max_len)

		var val string
		for i, cols := range rslt {
			t.Logf("-------------- %v --------------------\n", i+1)
			for j, col := range cols {
				if col == nil {
					val = "NULL"
				} else {
					val = string(col)
				}
				t.Logf(s, col_name[j], val)
			}
		}
	} else {
		fmt.Printf("err %s\n", err)
	}

}
