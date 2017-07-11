package logger

import (
	"runtime"

	"github.com/cihub/seelog"
)

var slog seelog.LoggerInterface

func Start(path string) {
	slog, _ = seelog.LoggerFromConfigAsFile(path)
	defer slog.Flush()
}

func Debug(debug string, params ...interface{}) {
	slog.Debugf(debug, params...)
}

func Info(info string, params ...interface{}) {
	slog.Infof(info, params...)
}

func Warning(waring string, params ...interface{}) {
	buf := make([]byte, 4096)
	l := runtime.Stack(buf, false)
	slog.Warnf(waring, params...)
	slog.Warnf("%s:%s", waring, buf[:l])
}

func Error(errorStr string, params ...interface{}) {
	buf := make([]byte, 4096)
	l := runtime.Stack(buf, false)
	slog.Warnf(errorStr, params...)
	slog.Warnf("%s:%s", errorStr, buf[:l])
}
