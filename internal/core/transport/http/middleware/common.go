package core_http_middleware

import (
	"context"
	"net/http"

	core_logger "github.com/KyoshiBlame/TodoKy/internal/core/logger"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

const RequestIDHeader = "X-Request-ID"

//middleware для идентификации запросов для удобства логов
func RequestID() Middleware {
	
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
			requestID := r.Header.Get(RequestIDHeader)
			if requestID == "" {
				requestID = uuid.NewString()
			}

			r.Header.Set(RequestIDHeader, requestID)
			w.Header().Set(RequestIDHeader, requestID)
			
			next.ServeHTTP(w,r)
		})
	}
}

func Logger(log *core_logger.Logger) Middleware {
	return func (next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
			RequestID := r.Header.Get(RequestIDHeader)

			l := log.With(
				zap.String("request_id", RequestID),
				zap.String("url", r.URL.String()),
			)
			
			ctx := context.WithValue(r.Context(), "log", l)//по ключу log передаём наш логгер

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
