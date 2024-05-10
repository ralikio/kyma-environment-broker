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
	if strings.HasSuffix(log.Message, "does not exist") ||
		 strings.HasSuffix(log.Message, "failed") ||
		 strings.HasSuffix(log.Message, "in progress") {

		log.LogLevel = lager.INFO
		sink.stdSink.Log(log)
	} else {
		sink.errSink.Log(log)
	}

}
