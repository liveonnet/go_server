package tools_lib

import (
	"fmt"
	"log"
	"sort"

	"learn/server/applib/log_lib"

	"gopkg.in/gin-gonic/gin.v1"
)

var logger *log.Logger

func init() {
	logger = log_lib.Log
}

func Check(e error, args ...string) {
	if e != nil {
		var prefix, postfix string
		switch len(args) {
		case 0:
		case 1:
			prefix = args[0]
		default:
			prefix = args[0]
			postfix = args[1]
		}
		log.Fatalf("%s got error: %s\nextra msg: %s\n", prefix, e, postfix)
		panic(fmt.Errorf("%s got error: %s", prefix, e))
	}
}

func EnvMgr() gin.HandlerFunc {
	return func(c *gin.Context) {
		if tmp, ok := c.Get("db"); ok {
			logger.Printf("db_map %s\n", tmp)
		}
		if tmp, ok := c.Get("cache"); ok {
			logger.Printf("cache_map %s\n", tmp)
		}
		//		c.Set("db", make(map[string]*sql.DB))
		//		c.Set("cache", make(map[string]*redis.Client))
		c.Next()
		//		if db_map, exist := c.Get("db"); exist {
		//			for k, v := range db_map.(map[string]*sql.DB) {
		//				logger.Printf("release %v %v\n", k, v)
		//				v.Close()
		//			}
		//		}
		//		if cache_map, exist := c.Get("cache"); exist {
		//			for k, v := range cache_map.(map[string]*redis.Client) {
		//				logger.Printf("release %v %v\n", k, v)
		//				v.Close()
		//			}
		//		}
	}
}

type ReverseStringSlice sort.StringSlice

func (p ReverseStringSlice) ReverseOrder() {
	l := len(p)
	logger.Printf("len(p) %v\n", l)
	if l > 1 {
		a, b := 0, l-1
		for {
			if a == b {
				break
			}
			p[a], p[b] = p[b], p[a]
			a++
			b--
		}

	}
}
