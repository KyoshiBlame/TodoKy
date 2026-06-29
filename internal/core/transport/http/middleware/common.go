package core_http_middleware

import (
	"net/http"

	"github.com/google/uuid"
)

//middleware для идентификации запросов для удобства логов
func RequestID() Middleware {

	const RequestIDHeader = "X-Request-ID"
	
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