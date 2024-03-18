package profiler

import (
	"fmt"
	"net/http"
	"os"
	gruntime "runtime"
	"runtime/pprof"
	"time"

	. "net/http/pprof"

	"code.cloudfoundry.org/lager"
	"github.com/gorilla/mux"
)

type ProfilerConfig struct {
	Path                  string        `envconfig:"default=/tmp/profiler"`
	Sampling              time.Duration `envconfig:"default=1s"`
	DebugEndpointsEnabled bool          `envconfig:"default=false"`
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
	profilePath string
	tracePath   string
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

func (p *Profiler) AttachRoutesIfSwitched(router *mux.Router) {
	if !p.config.DebugEndpointsEnabled {
		router.HandleFunc(`/{test:debug\/.*}`, http.NotFound)
		router.HandleFunc(`/debug\/{test:.*}`, http.NotFound)
		router.HandleFunc(`/debug`, http.NotFound)
		router.HandleFunc(`/debug/`, http.NotFound)
		return
	}

	// Attach routes
	router.HandleFunc("/debug/pprof", Index).Methods(http.MethodGet)
	router.HandleFunc("/debug/pprof/", Index).Methods(http.MethodGet)
	router.HandleFunc("/debug/pprof/heap", Index).Methods(http.MethodGet)
	router.HandleFunc("/debug/pprof/cmdline", Cmdline).Methods(http.MethodGet)
	router.HandleFunc("/debug/pprof/profile", Profile).Methods(http.MethodGet)
	router.HandleFunc("/debug/pprof/symbol", Symbol).Methods(http.MethodGet)
	router.HandleFunc("/debug/pprof/trace", Trace).Methods(http.MethodGet)
}
