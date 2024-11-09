package server

import (
	"bytes"
	"net/http"
)

type CustomResponseWriter struct {
	Buffer   *bytes.Buffer
	Headers  http.Header
	HTTPCode int
}

func NewCustomResponseWriter() *CustomResponseWriter {
	return &CustomResponseWriter{
		Buffer:   new(bytes.Buffer),
		Headers:  make(http.Header),
		HTTPCode: http.StatusOK,
	}
}

func (rw *CustomResponseWriter) Header() http.Header {
	return rw.Headers
}

func (rw *CustomResponseWriter) Write(data []byte) (int, error) {
	return rw.Buffer.Write(data)
}

func (rw *CustomResponseWriter) WriteHeader(statusCode int) {
	rw.HTTPCode = statusCode
}
