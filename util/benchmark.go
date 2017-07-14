package util

import (
	"os"
	"runtime"
	"runtime/pprof"

	"github.com/gfandada/gserver/logger"
)

var prof struct {
	cpu *os.File
	mem *os.File
}

func StartProfile(cpuprofile, memprofile string) {
	if cpuprofile != "" {
		f, err := os.Create(cpuprofile)
		if err != nil {
			logger.Error("cpuprofile: %v", err)
			return
		}
		logger.Info("writing CPU profile to: %s\n", cpuprofile)
		prof.cpu = f
		pprof.StartCPUProfile(prof.cpu)
	}
	if memprofile != "" {
		f, err := os.Create(memprofile)
		if err != nil {
			logger.Error("memprofile: %v", err)
		}
		logger.Info("writing mem profile to: %s\n", memprofile)
		prof.mem = f
		runtime.MemProfileRate = 4096
	}
}

func StopProfile() {
	if prof.cpu != nil {
		pprof.StopCPUProfile()
		prof.cpu.Close()
		logger.Info("CPU profile stopped")
	}
	if prof.mem != nil {
		pprof.Lookup("heap").WriteTo(prof.mem, 0)
		prof.mem.Close()
		logger.Info("mem profile stopped")
	}
}
