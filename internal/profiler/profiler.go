package profiler

import (
	"fmt"
	"os"
	gruntime "runtime"
	"runtime/pprof"
	"time"

	"code.cloudfoundry.org/lager"
)

type ProfilerConfig struct {
	Path     string        `envconfig:"default=/tmp/profiler"`
	Sampling time.Duration `envconfig:"default=1s"`
	Memory   bool
}

func NewProfiler(config ProfilerConfig, logger lager.Logger) *Profiler {
	return &Profiler{
		config,
		logger,
	}
}

type Profiler struct {
	config ProfilerConfig
	logger lager.Logger
}

func (p *Profiler) PeriodicProfile() {
	if p.config.Memory == false {
		return
	}
	p.logger.Info(fmt.Sprintf("Starting periodic profiler %v", p.config))
	if err := os.MkdirAll(p.config.Path, os.ModePerm); err != nil {
		p.logger.Error(fmt.Sprintf("Failed to create dir %v for profile storage", p.config.Path), err)
	}
	for {
		profName := fmt.Sprintf("%v/mem-%v.pprof", p.config.Path, time.Now().Unix())
		p.logger.Info(fmt.Sprintf("Creating periodic memory profile %v", profName))
		profFile, err := os.Create(profName)
		if err != nil {
			p.logger.Error(fmt.Sprintf("Creating periodic memory profile %v failed", profName), err)
		}
		err = pprof.Lookup("allocs").WriteTo(profFile, 0)
		if err != nil {
			p.logger.Error(fmt.Sprintf("Error while looking and writing allocs with profile %v", profName), err)
		}

		gruntime.GC()
		time.Sleep(p.config.Sampling)
	}
}
