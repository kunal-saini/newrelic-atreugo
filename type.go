package newrelic_atreugo

import "github.com/newrelic/go-agent/v3/newrelic"

type NewRelicAtreugo struct {
	newrelicApp *newrelic.Application
	instrumentResponseBody bool
}

type Options struct {
	InstrumentResponseBody bool
}

func NewRelicAtreugoWrapper(nra *newrelic.Application) *NewRelicAtreugo {
	return &NewRelicAtreugo{newrelicApp: nra, instrumentResponseBody: true}
}

func NewRelicAtreugoWrapperWithOptions(nra *newrelic.Application, ops *Options) *NewRelicAtreugo {
	na := &NewRelicAtreugo{newrelicApp: nra}
	if ops != nil {
		na.instrumentResponseBody = ops.InstrumentResponseBody
	}
	return na
}
