package main

// https://github.com/gin-gonic/gin

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	gin "gopkg.in/gin-gonic/gin.v1"

	redis "gopkg.in/redis.v5"

	"learn/go_server/applib/cache_lib"
	"learn/go_server/applib/comm_lib"
	"learn/go_server/applib/conf_lib"
	"learn/go_server/applib/db_lib"
	"learn/go_server/applib/log_lib"
	"learn/go_server/applib/tools_lib"
	"learn/go_server/applib/user_lib"
)

//func check(log *log.Logger, e error, prefix string, postfix string) {
//	if e != nil {
//		log.Panic("got error ", e)
//	}
//}

var Build string
var logger *log.Logger
var conf map[string]interface{}
var check func(e error, args ...string)

func init() {
	logger = log_lib.Log
	conf = conf_lib.Conf
	check = tools_lib.Check
}

func LoginHandler(c *gin.Context) {
	//	defer func() {
	//		if x := recover(); x != nil {
	//			//			func_name := runtime.FuncForPC(reflect.ValueOf(x).Pointer()).Name()
	//			pcs := make([]uintptr, 100)
	//			n := runtime.Callers(2, pcs)
	//			if n != 0 {
	//				func_info := runtime.FuncForPC(pcs[0])
	//				func_name := func_info.Name()
	//				file_name, line_no := func_info.FileLine(pcs[0])
	//				logger.Printf("recover from %v|%s %s:%v| %s\n", n, file_name, func_name, line_no, x)
	//				frms := runtime.CallersFrames(pcs)
	//				for {
	//					frm, more := frms.Next()
	//					logger.Printf("----| %s.%s:%v\n", frm.File, frm.Function, frm.Line)
	//					if !more {
	//						break
	//					}
	//					n--
	//					if n == 0 {
	//						break
	//					}
	//				}
	//			}
	//
	//			//			logger.Printf("+++++++++++++++++++++++\n%s\n----------------------------------\n", gin.stack(2))
	//			c.JSON(http.StatusBadRequest, gin.H{"err_code": 10000, "err_msg": "no op"})
	//		}
	//	}()
	op := c.DefaultQuery("op", "")

	var data interface{}
	c_tmp := (*cache_lib.MyContext)(c)
	cache, err := c_tmp.GetCache("default")
	check(err, "get cache conn")

	c_key := "xxxx"
	data, err = cache.Get(c_key).Result()
	if err == redis.Nil {
		c_tmp := (*db_lib.MyContext)(c)
		db, err := c_tmp.GetDB("uc_read")
		check(err, "get db")

		sql_str := "select uid, device_id, ctime, os_type, pnum from z_user where left(uid, 1) = '4' limit ?"
		rslt, col_name, err := db_lib.QueryRowsRaw(db, false, sql_str, 20)
		check(err, "query", sql_str)

		logger.Printf("op %s\n", op)
		data = gin.H{}
		for i, cols := range rslt {
			tmp := gin.H{}
			for j, v := range cols {
				if v == nil {
					v = sql.RawBytes("")
				}
				tmp[col_name[j]] = string(v)
			}
			data.(gin.H)[strconv.Itoa(i)] = tmp
		}

		j_data, err := json.Marshal(data) // map[string]interface{} -> []byte
		check(err, "json.Marshal")
		cache.Set(c_key, string(j_data), 1*time.Minute)
	} else {
		if err != nil {
			check(err, "cache", c_key)
		} else {
			// string -> map[string]interface{}
			data_tmp := gin.H{}
			json.Unmarshal([]byte(data.(string)), &data_tmp)
			//			json.Unmarshal(data.([]byte), &data_tmp) // trigger panic
			data = data_tmp

			// test get dup cache conn
			logger.Printf("cache hit %s\n", c_key)
			cache, err := (*cache_lib.MyContext)(c).GetCache("default")
			check(err, "cache default")
			logger.Printf("get default conn again %v\n", cache)
			// test get other cache conn
			cache, err = (*cache_lib.MyContext)(c).GetCache("ad")
			check(err, "cache ad")
			logger.Printf("get ad conn %v\n", cache)

			//			c.String(http.StatusOK, data.(string))
			//			return
		}

	}

	c.JSON(http.StatusOK, data)
}

func ScoreListHandler(c *gin.Context) {
	logger.Printf("get score_list\n")
}

func main() {
	logger.Println("start logging ...")
	logger.Printf("Using build: %s\n", Build)

	r := gin.New()
	if conf["debug"].(bool) {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	//	r.Use(gin.Logger(), gin.Recovery(), tools_lib.EnvMgr())
	r.Use(gin.Logger(), comm_lib.MyRecover(), tools_lib.EnvMgr())
	need_auth := r.Group("/user")
	need_auth.Use(user_lib.UserAuth())
	{
		need_auth.GET("/login", LoginHandler)
		need_auth.GET("/score_list", ScoreListHandler)
	}

	logger.Printf("listening on %v\n", conf["port"])
	r.Run(fmt.Sprintf(":%v", conf["port"]))
}
