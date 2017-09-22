package cache_lib

// https://github.com/go-redis/redis
// https://godoc.org/gopkg.in/redis.v5
// https://godoc.org/gopkg.in/redis.v5#pkg-examples

import (
	"fmt"
	"learn/go_server/applib/conf_lib"
	"learn/go_server/applib/log_lib"
	"learn/go_server/applib/tools_lib"
	"log"

	"gopkg.in/gin-gonic/gin.v1"

	"gopkg.in/redis.v5"
)

var logger *log.Logger
var conf map[string]interface{}
var m_cache map[string]*redis.Client
var check func(e error, args ...string)

func init() {
	logger = log_lib.Log
	conf = conf_lib.Conf
	m_cache = make(map[string]*redis.Client)
	check = tools_lib.Check
}

func _GetCacheOld(cache_name string) (conn *redis.Client, err error) {
	var ok bool
	if conn, ok = m_cache[cache_name]; !ok {
		cache_conf, _ := conf["cache"]
		if cache_conf, ok = cache_conf.(map[string]interface{})[cache_name]; ok {
			conn_conf := redis.Options{
				Addr:     fmt.Sprintf("%s:%v", cache_conf.(map[string]interface{})["host"], cache_conf.(map[string]interface{})["port"]),
				Password: cache_conf.(map[string]interface{})["password"].(string),
				DB:       cache_conf.(map[string]interface{})["db"].(int),
			}
			logger.Printf("conn_conf %v\n", conn_conf)
			conn = redis.NewClient(&conn_conf)
			m_cache[cache_name] = conn
		} else {
			logger.Printf("conf not found for %s\n", cache_name)
		}
	} else {
		logger.Printf("use per process cache conn %v\n", conn)
		err = conn.Ping().Err()
	}

	return conn, err
}

func GetCache(cache_name string) (conn *redis.Client, err error) {
	var ok bool
	cache_conf, _ := conf["cache"]
	if cache_conf, ok = cache_conf.(map[string]interface{})[cache_name]; ok {
		conn_conf := redis.Options{
			Addr:     fmt.Sprintf("%s:%v", cache_conf.(map[string]interface{})["host"], cache_conf.(map[string]interface{})["port"]),
			Password: cache_conf.(map[string]interface{})["password"].(string),
			DB:       cache_conf.(map[string]interface{})["db"].(int),
		}
		logger.Printf("conn_conf %+v\n", conn_conf)
		conn = redis.NewClient(&conn_conf)
		m_cache[cache_name] = conn
	} else {
		err = logger.Output(1, fmt.Sprintf("conf not found for %s\n", cache_name))
	}

	return conn, err
}

type MyContext gin.Context

func (c *MyContext) GetCache(cache_name string) (conn *redis.Client, err error) {
	c_tmp := (*gin.Context)(c)
	var (
		cache_map map[string]*redis.Client
		exist     bool
	)
	if v, exist := c_tmp.Get("cache"); !exist {
		cache_map = make(map[string]*redis.Client)
		c_tmp.Set("cache", cache_map)
	} else {
		cache_map = v.(map[string]*redis.Client)
	}

	if conn, exist = cache_map[cache_name]; !exist {
		if conn, err = GetCache(cache_name); err == nil {
			cache_map[cache_name] = conn
		}
	} else {
		logger.Printf("use per request existing cache conn %v\n", conn)
	}

	return conn, err
}
