package core_http_middleware

import (
	"context"
	"net/http"
	"time"

	core_logger "github.com/compositemannn/GoMarket/internal/core/logger"
	core_http_response "github.com/compositemannn/GoMarket/internal/core/transport/http/response"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

const requestIDHeader = "X-Request-ID"

func RequestID() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestID := r.Header.Get(requestIDHeader)
			if requestID == "" {
				requestID = uuid.NewString()
			}
			r.Header.Set(requestIDHeader, requestID)   // проставляем id в хэддерах http-запроса
			w.Header().Set(requestIDHeader, requestID) // проставляем id в хэддерах http-ответа

			next.ServeHTTP(w, r) // вызов обработчика
		})
	}
}

// связывает все логи с http-запросом через производный логгер
func Logger(log *core_logger.Logger) Middleware { // возвращает Middleware, что имеет доступ к переданному логгеру
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestID := r.Header.Get(requestIDHeader)
			l := log.With( // указатель на новый логгер с дополнительными полями, который после передадим в контекст через ключ
				zap.String("request_id", requestID),
				zap.String("url", r.URL.String()),
			)

			ctx := context.WithValue(r.Context(), "log", l) // полученный контекст сохраняет связь с исходным, включая отмену и дедлайн

			next.ServeHTTP(w, r.WithContext(ctx)) // вызов обработчика, у которого Request с новым ctx
		})
	}
}

func Panic() Middleware { // отлавливает паники внутри обработчиков и записывает логи
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log := core_logger.FromContext(r.Context())
			responseHandler := core_http_response.NewHTTPResponseHandler(log, w)
			defer func() {
				if p := recover(); p != nil { // recover перехватывает панику в текущей горутине
					// лог с ID запроса, информацией о панике, http.500, тело ответа
					responseHandler.PanicResponse(
						p,
						"during handle HTTP request got unexpected panic",
					)
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}

func Trace() Middleware { // записывает необходимую инфу в логи при выполнении обработчика
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log := core_logger.FromContext(r.Context())
			rw := core_http_response.NewResponseWriter(w)
			before := time.Now()

			log.Debug(
				">>> incoming HTTP request",
				zap.Time("time", before.UTC()),
			)

			next.ServeHTTP(rw, r)

			log.Debug(
				"<<< done HTTP request",
				zap.Duration("latency", time.Since(before)),
				zap.Int("status_code", rw.GetStatusCode()),
			)
		})
	}
}
