package log

import (
	"fmt"
	"os"
	"time"
)

var FPROXY *os.File
var FPROXY_LOG *os.File

func init() {
	fproxy, err := os.Create("log.txt")
	if err == nil {
		FPROXY = fproxy
	}
	fproxy_log, err := os.Create("proxy.txt")
	if err == nil {
		FPROXY_LOG = fproxy_log
	}
}
func Pr(typex string, ip string, text string, a ...interface{}) {
	fmt.Fprintln(FPROXY, "["+typex+"]^"+ip+"^["+time.Now().Format("2006-01-02-15:04:05")+"]^"+text+" ", a)
}
func Proxy_log(typex string, ip string, text string, a ...interface{}) {
	fmt.Fprintln(FPROXY_LOG, "["+typex+"]^"+ip+"^["+time.Now().Format("2006-01-02-15:04:05")+"]^"+text+" ", a)
}
