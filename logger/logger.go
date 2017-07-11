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

func Info(log string) {
	slog.Info(log)
}

func Warning(waring string) {
	buf := make([]byte, 4096)
	l := runtime.Stack(buf, false)
	slog.Warnf("%s:%s", waring, buf[:l])
}

func Error(errorStr string) {
	buf := make([]byte, 4096)
	l := runtime.Stack(buf, false)
	slog.Errorf("%s:%s", errorStr, buf[:l])
}
