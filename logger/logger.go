package logger

import (
	"flag"
	"runtime"

	"github.com/golang/glog"
)

const (
	InfoLevel    = 0
	WarningLevel = 1
	ErrorLevel   = 2
)

// TODO 暂时不提供配置
var level = InfoLevel

func Start() {
	flag.Parse()
	defer glog.Flush()
}

func Info(info string) {
	if level <= InfoLevel {
		glog.Infof("%s", info)
	}
}

func Warning(waring string) {
	if level <= WarningLevel {
		buf := make([]byte, 4096)
		l := runtime.Stack(buf, false)
		glog.Warningf("%s:%s", waring, buf[:l])
	}
}

func Error(errorStr string) {
	if level <= ErrorLevel {
		buf := make([]byte, 4096)
		l := runtime.Stack(buf, false)
		glog.Errorf("%s:%s", errorStr, buf[:l])
	}
}
