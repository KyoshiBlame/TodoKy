package web_transport_http

import (
	"net/http"

	core_logger "github.com/KyoshiBlame/TodoKy/internal/core/logger"
	core_http_response "github.com/KyoshiBlame/TodoKy/internal/core/transport/http/response"
)

func (h *WebHTTPHandler) GetAuthPage(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	html, err := h.webService.GetAuthPage()
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get auth.html")
		return
	}

	responseHandler.HTMLResponse(html)
}
