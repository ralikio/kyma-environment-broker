package logger

import (
	"bytes"
	"fmt"
	"testing"

	"code.cloudfoundry.org/lager"

	"github.com/stretchr/testify/assert"
)

func TestFilterSink(t *testing.T) {

	t.Run("should filter out preconfigured message", func(t *testing.T) {
		// given
		logger := lager.NewLogger("kyma-env-broker")
		stdBuf := bytes.NewBufferString("")
		errBuf := bytes.NewBufferString("")
		logger.RegisterSink(NewFilterSink(lager.NewWriterSink(stdBuf, lager.INFO), lager.NewWriterSink(errBuf, lager.ERROR)))

		// when
		inputMessage := "any message"
		logger.Error(inputMessage, fmt.Errorf("sample error"))

		// then
		assert.Contains(t, errBuf.String(), fmt.Sprintf("\"log_level\":%d", lager.ERROR))
		assert.Empty(t, stdBuf.String())

		// when
		errBuf.Reset()
		stdBuf.Reset()
		assert.Empty(t, errBuf.String())
		assert.Empty(t, stdBuf.String())

		inputMessage = "kyma-env-broker.getInstance.provisioning of instanceID 6CEA47A3-157F-4978-B364-AFAA9775DF9A failed"
		logger.Error(inputMessage, fmt.Errorf("sample error"))
		outputMessage := stdBuf.String()

		// then
		assert.Empty(t, errBuf.String())
		assert.NotEmpty(t, outputMessage)
		assert.Contains(t, outputMessage, fmt.Sprintf("\"log_level\":%d", lager.INFO))

		// when
		errBuf.Reset()
		stdBuf.Reset()
		assert.Empty(t, errBuf.String())
		assert.Empty(t, stdBuf.String())

		inputMessage = "kyma-env-broker.getInstance.provisioning of instanceID 6CEA47A3-157F-4978-B364-AFAA9775DF9A does not exist"
		logger.Error(inputMessage, fmt.Errorf("sample error"))
		outputMessage = stdBuf.String()

		// then
		assert.Empty(t, errBuf.String())
		assert.NotEmpty(t, outputMessage)
		assert.Contains(t, outputMessage, fmt.Sprintf("\"log_level\":%d", lager.INFO))

		// when
		errBuf.Reset()
		stdBuf.Reset()
		assert.Empty(t, errBuf.String())
		assert.Empty(t, stdBuf.String())

		inputMessage = "kyma-env-broker.getInstance.provisioning of instanceID 6CEA47A3-157F-4978-B364-AFAA9775DF9A in progress"
		logger.Error(inputMessage, fmt.Errorf("sample error"))
		outputMessage = stdBuf.String()

		// then
		assert.Empty(t, errBuf.String())
		assert.NotEmpty(t, outputMessage)
		assert.Contains(t, outputMessage, fmt.Sprintf("\"log_level\":%d", lager.INFO))
		
	})
}
