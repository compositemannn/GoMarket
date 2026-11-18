package core_http_response

import "net/http"

var (
	StatusCodeUninitialized = -1
)

type ResponseWriter struct { // структура нужна, чтобы узнать какой статус-код был отправлен после http-ответа
	http.ResponseWriter // здесь интерфейс не имеет метода GetStatusCode, лежит указатель на структуру
	statusCode          int
}

func NewResponseWriter(w http.ResponseWriter) *ResponseWriter {
	return &ResponseWriter{
		ResponseWriter: w,
		statusCode:     StatusCodeUninitialized,
	}
}

func (rw *ResponseWriter) WriteHeader(statusCode int) {
	if rw.statusCode != StatusCodeUninitialized {
		return
	}

	rw.statusCode = statusCode
	rw.ResponseWriter.WriteHeader(statusCode)
}

func (rw *ResponseWriter) Write(p []byte) (int, error) {
	if rw.statusCode == StatusCodeUninitialized {
		rw.statusCode = http.StatusOK
	}

	return rw.ResponseWriter.Write(p)
}

func (rw *ResponseWriter) GetStatusCode() int {
	if rw.statusCode == StatusCodeUninitialized {
		return http.StatusOK
	}

	return rw.statusCode
}
