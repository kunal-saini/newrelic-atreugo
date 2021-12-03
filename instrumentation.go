package newrelic_atreugo

import (
	"fmt"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/savsgio/atreugo/v11"
	"github.com/valyala/fasthttp"
	"net/http"
	"net/url"
)

const NewRelicTransaction = "__newrelic_transaction__"

// BeginInstrumentationMiddleware returns a middleware.
func (n *NewRelicAtreugo) BeginInstrumentationMiddleware() atreugo.Middleware {
	return func(ctx *atreugo.RequestCtx) error {
		txn := n.newrelicApp.StartTransaction(string(ctx.Method()) + " " + string(ctx.Path()))
		defer logError(txn)

		txn.SetWebRequest(convertToNewRelicWebRequest(ctx, txn))
		ctx.SetUserValue(NewRelicTransaction, txn)
		return ctx.Next()
	}
}

// EndInstrumentationMiddleware returns a middleware.
func (n *NewRelicAtreugo) EndInstrumentationMiddleware() atreugo.Middleware {
	return func(ctx *atreugo.RequestCtx) error {
		val := ctx.UserValue(NewRelicTransaction)
		if txn, ok := val.(*newrelic.Transaction); ok {
			defer logError(txn)
			statusCode := ctx.Response.StatusCode()
			response := newResponse(statusCode)

			if n.instrumentResponseBody {
				respCopy := &fasthttp.Response{}
				ctx.Response.CopyTo(respCopy)
				respBody := respCopy.Body()
				if len(respBody) > 0 {
					response = newResponseWithBody(respBody, statusCode)
				}
			}

			txn.SetWebResponse(response)
			txn.End()
		}
		return ctx.Next()
	}
}

// PanicView returns a panic view.
func (n *NewRelicAtreugo) PanicView() atreugo.PanicView {
	return func(ctx *atreugo.RequestCtx, err interface{}) {
		val := ctx.UserValue(NewRelicTransaction)
		if txn, ok := val.(*newrelic.Transaction); ok {
			logError(txn)
		}
	}
}

func logError(txn *newrelic.Transaction)  {
	if err := recover(); err != nil {
		switch err := err.(type) {
		case error:
			txn.NoticeError(err)
		default:
			txn.NoticeError(errWrapper{err})
		}
	}
}

func convertToNewRelicWebRequest(ctx *atreugo.RequestCtx, txn *newrelic.Transaction) newrelic.WebRequest {
	r := new(newrelic.WebRequest)

	r.Method = string(ctx.Method())
	uri := ctx.URI()
	// Ignore error.
	r.URL, _ = url.Parse(fmt.Sprintf("%s://%s%s", uri.Scheme(), uri.Host(), uri.Path()))

	// Headers
	r.Header = make(http.Header)
	r.Header.Add("Host", string(ctx.Host()))
	ctx.Request.Header.VisitAll(func(key, value []byte) {
		r.Header.Add(string(key), string(value))
	})
	r.Host = string(ctx.Host())

	// QueryString
	r.URL.RawQuery = string(ctx.URI().QueryString())

	return *r
}
