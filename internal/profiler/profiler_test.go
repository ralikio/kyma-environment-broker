package profiler

import (
	"context"
	"sync"
	"testing"

	"net/http"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)


func TestProfiler_StartIfSwitched(t *testing.T) {
	t.Run("should create profiler files when profiler feature flag is turned on", func(t *testing.T) {
		// given
		profilerConfig := ProfilerConfig{
			Path:                  "dummy",
			Sampling:              1,
			DebugEndpointsEnabled: true,
			Memory:                false,
		}
		profiler := NewProfiler(profilerConfig, nil)
		router := mux.NewRouter()
		
		// attach routes
		profiler.AttachRoutesIfSwitched(router)

		// start server
		srv := &http.Server{Addr: ":8080"}
		var wg sync.WaitGroup
		wg.Add(1)
		go func() {
			srv.ListenAndServe()
			wg.Done()
		}()

		// when
		resp, err := http.Get("http://localhost:8080/debug/pprof/")
		assert.NoError(t, err)
		defer resp.Body.Close()

		// then
		assert.Equal(t, http.StatusOK, resp.StatusCode)	
		assert.NoError(t, err)
		
		// cleanup
		srv.Shutdown(context.Background())
		wg.Wait()
	})
}