package core_http_response

import "net/http"

//чтобы отличить что статс код точно был проставлен
var (
	StatusCodeUninitialized = -1
)

type ResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func NewResponseWriter(w http.ResponseWriter) *ResponseWriter{
	return &ResponseWriter {
		ResponseWriter: w,
	}
}

//записываем передаваемый статус код через WriteHeader и запоминаем его
func (rw *ResponseWriter) WriteHeader(statusCode int) {
	rw.ResponseWriter.WriteHeader(statusCode)
	rw.statusCode = StatusCodeUninitialized
}

func (rw *ResponseWriter) GetStatusCodeOrPanic() int{
	if rw.statusCode == StatusCodeUninitialized {
		panic("no status code set")
	}

	return rw.statusCode
}