package profiler

import (
	"fmt"
	"net/http"
	"os"
	gruntime "runtime"
	"runtime/pprof"
	"runtime/trace"
	"time"

	httpPprof "net/http/pprof"

	"code.cloudfoundry.org/lager"
	"github.com/gorilla/mux"
)

type ProfilerConfig struct {
	Path                  string        `envconfig:"default=/tmp/profiler"`
	Sampling              time.Duration `envconfig:"default=1s"`
	DebugEndpointsEnabled bool          `envconfig:"default=true"`
	Memory                bool
}

func NewProfiler(config ProfilerConfig, logger lager.Logger) *Profiler {
    timestamp := time.Now().Unix()

	return &Profiler{
		config,
		logger,
		nil,
		nil,
        fmt.Sprintf("%v/profile-%v.prof", config.Path, timestamp),
	    fmt.Sprintf("%v/trace-%v.prof", config.Path, timestamp),
	}
}

type Profiler struct {
	config      ProfilerConfig
	logger      lager.Logger
	profileFile *os.File
	traceFile   *os.File
	profilePath   string
	tracePath   string
}

func (p *Profiler) PeriodicProfile() {
	// if p.config.Memory == false {
	// 	return
	// }
	// p.logger.Info(fmt.Sprintf("Starting periodic profiler %v", p.config))
	// if err := os.MkdirAll(p.config.Path, os.ModePerm); err != nil {
	// 	p.logger.Error(fmt.Sprintf("Failed to create dir %v for profile storage", p.config.Path), err)
	// }
	// for {
	// 	profName := fmt.Sprintf("%v/mem-%v.pprof", p.config.Path, time.Now().Unix())
	// 	p.logger.Info(fmt.Sprintf("Creating periodic memory profile %v", profName))
	// 	profFile, err := os.Create(profName)
	// 	if err != nil {
	// 		p.logger.Error(fmt.Sprintf("Creating periodic memory profile %v failed", profName), err)
	// 	}
	// 	err = pprof.Lookup("allocs").WriteTo(profFile, 0)
	// 	if err != nil {
	// 		p.logger.Error(fmt.Sprintf("Error while looking and writing allocs with profile %v", profName), err)
	// 	}

	// 	gruntime.GC()
	// 	time.Sleep(p.config.Sampling)
	// }
}

func (p *Profiler) StartIfSwitched() {
	if !p.config.DebugEndpointsEnabled {
		return
	}

	// var err error
	// p.profileFile, err = os.Create(p.profilePath)
	// if err != nil {
	// 	p.logger.Error("Failed to create profile file", err)
	// }

	if err := pprof.StartCPUProfile(p.profileFile); err != nil {
		p.logger.Error("Failed to start CPU profile", err)
	}

	// p.traceFile, err = os.Create(fmt.Sprintf("%v/trace-%v.prof", p.config.Path, time.Now().Unix()))
	// if err != nil {
	// 	p.logger.Error("Failed to create trace file", err)
	// }

	// if err := trace.Start(p.traceFile); err != nil {
	// 	p.logger.Error("Failed to start trace", err)
	// }
}

func (p *Profiler) AttachRoutesIfSwitched(router *mux.Router) {
	if !p.config.DebugEndpointsEnabled {
		return
	}

	// Attach routes
	router.HandleFunc("/debug/pprof/", httpPprof.Index).Methods(http.MethodGet)
	router.HandleFunc("/debug/pprof/cmdline", httpPprof.Cmdline).Methods(http.MethodGet)
	router.HandleFunc("/debug/pprof/profile", httpPprof.Profile).Methods(http.MethodGet)
	router.HandleFunc("/debug/pprof/symbol", httpPprof.Symbol).Methods(http.MethodGet)
	router.HandleFunc("/debug/pprof/trace", httpPprof.Trace).Methods(http.MethodGet)
}

func (p *Profiler) StopIfSwitched() {
	if !p.config.DebugEndpointsEnabled {
		return
	}

	p.profileFile.Close()
	pprof.StopCPUProfile()
	p.traceFile.Close()
	trace.Stop()
}
