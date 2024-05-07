package logger

import (
	"strings"

	"code.cloudfoundry.org/lager"
)
type FilterSink struct {
	stdSink lager.Sink
	errSink lager.Sink
}

func NewFilterSink(stdSink lager.Sink, errSink lager.Sink) *FilterSink {
	return &FilterSink{
		stdSink: stdSink,
		errSink: errSink,
	}
}

func (sink *FilterSink) Log(log lager.LogFormat) {
    if strings.HasPrefix(log.Message, "kyma-env-broker.getInstance.provisioning of instanceID") {
        log.LogLevel = lager.INFO
		sink.stdSink.Log(log)
    } else {
		sink.errSink.Log(log)
	}

}
