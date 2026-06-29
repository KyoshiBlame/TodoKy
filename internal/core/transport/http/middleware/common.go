package core_http_middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

//middleware для идентификации запросов для удобства логов
func RequestID() gin.HandlerFunc {
	const requestIDHeader = "X-Request-ID"

	return func (c *gin.Context) {
		requestID := c.GetHeader(requestIDHeader)
		if requestID == "" {
			requestID = uuid.NewString()
		}

		c.Request.Header.Set(requestIDHeader, requestID)
		c.Header(requestIDHeader, requestID)
		c.Set("request_id", requestID)
		c.Next()
	}
}