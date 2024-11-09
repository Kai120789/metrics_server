package server

import (
	"bytes"
	"net/http"
)

// CustomResponseWriter реализует интерфейс http.ResponseWriter для захвата данных вывода
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

// Header возвращает заголовки ответа
func (rw *CustomResponseWriter) Header() http.Header {
	return rw.Headers
}

// Write записывает данные в буфер
func (rw *CustomResponseWriter) Write(data []byte) (int, error) {
	return rw.Buffer.Write(data)
}

// WriteHeader устанавливает код состояния ответа
func (rw *CustomResponseWriter) WriteHeader(statusCode int) {
	rw.HTTPCode = statusCode
}
