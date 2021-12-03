# newrelic-atreugo
NewRelic instrumentation middleware for Atreugo 

## example
```go
import (
"github.com/newrelic/go-agent/v3/newrelic"
newrelicAtreugo "github.com/kunal-saini/newrelic-atreugo"
)

newrelicApp, err := newrelic.NewApplication(
    newrelic.ConfigAppName("app name"),
    newrelic.ConfigLicense("license"),
    newrelic.ConfigDistributedTracerEnabled(true),
)

newrelicAtreugoMiddleware := newrelicAtreugo.NewRelicAtreugoWrapper(newrelicApp)
server := atreugo.New(atreugo.Config{
	Addr: "0.0.0.0:3000", 
	PanicView: newrelicAtreugoMiddleware.PanicView()
})

server.UseBefore(newrelicAtreugoMiddleware.BeginInstrumentationMiddleware())
server.UseAfter(newrelicAtreugoMiddleware.EndInstrumentationMiddleware())
```