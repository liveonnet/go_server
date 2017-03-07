package log_lib

import (
	"log"
	"os"

	"learn/server/applib/conf_lib"
)

var Log *log.Logger

func init() {
	Log = log.New(os.Stdout, "root:", log.Lshortfile|log.Ldate|log.Lmicroseconds)
	Log.Println("logger init done")
	if conf_lib.Conf != nil {
		Log.Printf("debug mode: %v\n", conf_lib.Conf["debug"])
	}
}
