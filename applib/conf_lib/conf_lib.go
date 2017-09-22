package conf_lib

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"reflect"
	"strings"

	"gopkg.in/yaml.v2"
)

var Conf map[string]interface{}

func check(log *log.Logger, e error) {
	if e != nil {
		if log == nil {
			fmt.Printf("got error %s\n", e)
		} else {
			log.Fatalf("got error %s\n", e)
		}
	}
}

func init() {
	//	file_name := "./config/server.yaml"
	file_name := "/home/kevin/data_bk/go/src/learn/go_server/config/server.yaml"
	// path convert
	if !filepath.IsAbs(file_name) {
		if strings.HasPrefix(file_name, "~/") {
			usr, err := user.Current()
			check(nil, err)
			file_name = filepath.Join(usr.HomeDir, file_name[2:])
			fmt.Println("file path ->", file_name)
		} else {
			var err error
			file_name, err = filepath.Abs(file_name)
			fmt.Println("file path ->", file_name)
			check(nil, err)
		}
	}

	if _, err := os.Stat(file_name); os.IsNotExist(err) {
		fmt.Printf("conf file not found !!! %s\n", file_name)
	}

	data, err := ioutil.ReadFile(file_name)
	check(nil, err)
	tmp_conf := make(map[string]interface{})
	err = yaml.Unmarshal([]byte(data), &tmp_conf)
	if err != nil {
		fmt.Printf("error %v\n", err)
	} else {
		//		Conf = tmp_conf
		Conf = convert(tmp_conf).(map[string]interface{})
		fmt.Printf("Conf: %v\n", Conf)
		//		s_test, err := yaml.Marshal(tmp_conf)
		//		if err != nil {
		//			fmt.Printf("err: %s\n", err)
		//		} else {
		//			fmt.Printf("conf: %s\n", s_test)
		//		}
	}
}

func convert(m interface{}) interface{} {
	var r interface{}
	switch m.(type) {
	case map[interface{}]interface{}:
		//		fmt.Printf("!!!!!!!\n")
		r = make(map[string]interface{})
		x_m := m.(map[interface{}]interface{})
		for k, v := range x_m {
			str_k, ok := k.(string)
			if ok {
				r.(map[string]interface{})[str_k] = convert(v)
			} else {
				fmt.Printf("Unknown type key: %s\n", k)
			}

		}
	case map[string]interface{}:
		//		fmt.Printf("~~~~~~~~~~\n")
		r = make(map[string]interface{})
		for k, v := range m.(map[string]interface{}) {
			r.(map[string]interface{})[k] = convert(v)
		}
	case string:
		r = m.(string)
	case int:
		r = m.(int)
	case bool:
		r = m.(bool)
	case nil:
		r = nil
	case []interface{}:
		r = m
	default:
		fmt.Printf("type Unknown %v\n", reflect.TypeOf(m))
		r = 0
	}
	return r
}
