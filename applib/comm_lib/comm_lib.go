package comm_lib

// 从gin.Recover拷贝
// 过滤了堆栈开头的非自己的代码的信息
// 返回200且unknown error而不是500且白板

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"learn/server/applib/conf_lib"
	"learn/server/applib/log_lib"
	"learn/server/applib/tools_lib"
	"log"
	"net/http"
	"runtime"
	"sort"
	"strings"

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

var (
	dunno     = []byte("???")
	centerDot = []byte("·")
	dot       = []byte(".")
	slash     = []byte("/")
	reset     = string([]byte{27, 91, 48, 109})
	own_file  = string("learn/server/")
)

func MyRecover() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				stack := stack(3)
				//				httprequest, _ := httputil.DumpRequest(c.Request, false)
				//				logger.Printf("\n\n\x1b[31m[Recovery] panic recovered:\n%s\n%s\n%s%s", string(httprequest), err, stack, reset)
				httprequest := c.Request.URL.String()
				//				logger.Printf("\n\n\x1b[31m[Recovery] panic recovered:\n%s %s\n\nPANIC: %s\n%s%s", c.Request.Method, string(httprequest), err, stack, reset)
				logger.Printf("\x1b[31m[Recovery] panic recovered:\n%s\nPANIC: %s\n\n%s %s%s", stack, err, c.Request.Method, string(httprequest), reset)
				c.Abort()
				c.JSON(http.StatusOK, gin.H{"err_code": 10000, "err_msg": "unknown error"})
			}
		}()
		c.Next()
	}
}

// stack returns a nicely formated stack frame, skipping skip frames
func stack(skip int) []byte {
	buf := new(bytes.Buffer) // the returned data
	var my_lines sort.StringSlice
	// As we loop, we open files and read them. These variables record the currently
	// loaded file.
	var lines [][]byte
	var lastFile string
	line_skip := true
	for i := skip; ; i++ { // Skip the expected number of frames
		pc, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}

		if line_skip {
			if strings.Index(file, own_file) == -1 {
				continue
			} else {
				line_skip = false
			}
		}
		// Print this much at least.  If we can't find the source, it won't show.
		fmt.Fprintf(buf, "%s:%d (0x%x)\n", file, line, pc)
		if file != lastFile {
			data, err := ioutil.ReadFile(file)
			if err != nil {
				continue
			}
			lines = bytes.Split(data, []byte{'\n'})
			lastFile = file
		}
		fmt.Fprintf(buf, "\t%s: %s\n", function(pc), source(lines, line))
		my_lines = append(my_lines, buf.String())
		buf.Reset()
	}
	//	logger.Printf("len(my_lines) %v\n", len(my_lines))

	tools_lib.ReverseStringSlice(my_lines).ReverseOrder()
	x := strings.Join(my_lines, "")

	return []byte(x)
}

// source returns a space-trimmed slice of the n'th line.
func source(lines [][]byte, n int) []byte {
	n-- // in stack trace, lines are 1-indexed but our array is 0-indexed
	if n < 0 || n >= len(lines) {
		return dunno
	}
	return bytes.TrimSpace(lines[n])
}

// function returns, if possible, the name of the function containing the PC.
func function(pc uintptr) []byte {
	fn := runtime.FuncForPC(pc)
	if fn == nil {
		return dunno
	}
	name := []byte(fn.Name())
	// The name includes the path name to the package, which is unnecessary
	// since the file name is already included.  Plus, it has center dots.
	// That is, we see
	//	runtime/debug.*T·ptrmethod
	// and want
	//	*T.ptrmethod
	// Also the package path might contains dot (e.g. code.google.com/...),
	// so first eliminate the path prefix
	if lastslash := bytes.LastIndex(name, slash); lastslash >= 0 {
		name = name[lastslash+1:]
	}
	if period := bytes.Index(name, dot); period >= 0 {
		name = name[period+1:]
	}
	name = bytes.Replace(name, centerDot, dot, -1)
	return name
}
