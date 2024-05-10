package logger

import (
	"bytes"
	"fmt"
	"testing"

	"code.cloudfoundry.org/lager"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestFilterSink(t *testing.T) {

	t.Run("should filter out preconfigured message", func(t *testing.T) {
		logger := lager.NewLogger("kyma-env-broker")
		stdBuf := bytes.NewBufferString("")
		errBuf := bytes.NewBufferString("")
		logger.RegisterSink(NewFilterSink(lager.NewWriterSink(stdBuf, lager.INFO), lager.NewWriterSink(errBuf, lager.ERROR)))

		// when
		inputMessage := "any message"
		logger.Error(inputMessage, fmt.Errorf("sample error"))

		assert.Contains(t, errBuf.String(), fmt.Sprintf("\"log_level\":%d", lager.ERROR))
		assert.Empty(t, stdBuf.String())

		errBuf.Reset()
		stdBuf.Reset()
		assert.Empty(t, errBuf.String())
		assert.Empty(t, stdBuf.String())

		// when
		inputMessage = "getInstance.provisioning of instanceID " + uuid.New().String()
		logger.Error(inputMessage, fmt.Errorf("sample error"))

		outputMessage := stdBuf.String()
		assert.Empty(t, errBuf.String())
		assert.NotEmpty(t, outputMessage)

		assert.Contains(t, outputMessage, fmt.Sprintf("\"log_level\":%d", lager.INFO))
	})
}
