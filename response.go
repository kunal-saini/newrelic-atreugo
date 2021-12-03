package newrelic_atreugo

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
)

type Response struct {
	io.Writer
	h int
}

func newResponse(statusCode int) http.ResponseWriter {
	return &Response{h : statusCode}
}

func newResponseWithBody(data []byte, statusCode int) http.ResponseWriter {
	return &Response{Writer: bytes.NewBuffer(data), h : statusCode}
}

func (w *Response) Header() http.Header { return make(http.Header) }
func (w *Response) WriteHeader(h int)   { w.h = h }
func (w *Response) String() string      { return fmt.Sprintf("[%v] %q", w.h, w.Writer) }