package user_lib

// https://github.com/dgrijalva/jwt-go
// https://godoc.org/github.com/dgrijalva/jwt-go

import (
	"learn/server/applib/conf_lib"
	"learn/server/applib/log_lib"
	"learn/server/applib/tools_lib"
	"log"
	"net/http"
	"strconv"
	"time"

	jwt "github.com/dgrijalva/jwt-go"

	"gopkg.in/gin-gonic/gin.v1"
)

var logger *log.Logger
var conf map[string]interface{}
var check func(e error, args ...string)

func init() {
	logger = log_lib.Log
	conf = conf_lib.Conf
	check = tools_lib.Check
}

type MyCustomClaims struct {
	Uid string `json:"uid"`
	jwt.StandardClaims
}

func CreateTicket(uid int) (ticket string, err error) {
	sign_key := conf["user_login_key"].(string)
	claims := MyCustomClaims{
		strconv.Itoa(uid),
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(24 * 365 * time.Duration(time.Hour)).Unix(),
			Issuer:    "xxhb",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ticket, err = token.SignedString([]byte(sign_key))
	logger.Printf("%v, %v\n", ticket, err)
	return ticket, err
}

func ValidateTicket(ticket string) (uid int, err error) {
	defer func() {
		if x := recover(); x != nil {
			logger.Printf("recovered from %s\n", x)
		}
	}()
	token, err := jwt.ParseWithClaims(ticket, &MyCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(conf["user_login_key"].(string)), nil
	})
	if err == nil {
		if claims, ok := token.Claims.(*MyCustomClaims); ok && token.Valid {
			uid, err = strconv.Atoi(claims.Uid)
			if err == nil {
				logger.Printf("ticket ok for uid: %v, expire time %v(%s) issuer: %v\n", uid, claims.StandardClaims.ExpiresAt, time.Unix(claims.ExpiresAt, 0), claims.Issuer)
			} else {
				logger.Printf("uid string -> int error: %s\n", err)
			}
		} else {
			logger.Println(err)
		}
	} else {
		logger.Printf("error: %s\n", err)
	}

	return uid, err
}

func UserAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		//		defer func() {
		//			if x := recover(); x != nil {
		//				logger.Println("recover from ", x)
		//				c.Abort()
		//				c.JSON(http.StatusUnauthorized, gin.H{"err_code": 10000, "err_msg": "param 'ticket' is missing"})
		//			}
		//		}()
		ticket := c.DefaultQuery("ticket", "")
		if ticket != "" {
			uid, err := ValidateTicket(ticket)
			if err == nil {
				logger.Printf("auth ok uid %v\n", uid)
				c.Set("uid", uid)
				c.Next()
				c.Set("uid", nil) // not necessary
			} else {
				c.Abort()
				c.JSON(http.StatusUnauthorized, gin.H{"err_code": 10000, "err_msg": "'ticket' invalid/expired"})
			}
		} else {
			c.Abort()
			c.JSON(http.StatusUnauthorized, gin.H{"err_code": 10000, "err_msg": "param 'ticket' is missing/empty"})
		}
	}
}
