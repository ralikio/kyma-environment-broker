package profiler

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)


func TestProfiler_StartIfSwitched(t *testing.T) {
	t.Run("should create profiler files when profiler feature flag is turned on", func(t *testing.T) {
		// given
		profilerConfig := ProfilerConfig{
			Path:                  "/tmp/profiler",
			Sampling:              1,
			DebugEndpointsEnabled: true,
			Memory:                false,
		}
		profiler := NewProfiler(profilerConfig, nil)
		defer profiler.StopIfSwitched()

		// when
		profiler.StartIfSwitched()

		// then
		profileFile, err := os.Stat(profiler.profilePath)
		assert.NoError(t, err)
		assert.False(t, profileFile.IsDir())
		assert.True(t, profileFile.Size() > 0)
	})
}