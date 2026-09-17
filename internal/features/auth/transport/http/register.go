package auth_transport_http

import (
	"net/http"

	core_logger "github.com/KyoshiBlame/TodoKy/internal/core/logger"
	core_http_request "github.com/KyoshiBlame/TodoKy/internal/core/transport/http/request"
	core_http_response "github.com/KyoshiBlame/TodoKy/internal/core/transport/http/response"
)

type RegisterRequest struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
	Fullname string `json:"full_name" validate:"required"`
}

type RegisterResponse struct {
	UserID int64 `json:"user_id"`
}

func (h *AuthHTTPHandler) Register(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(
		log,
		rw,
	)

	var req RegisterRequest

	if err := core_http_request.DecodeAndValidateRequest(r, &req); err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to decode register request",
		)
		return
	}

	userID, err := h.authClient.Register(
		ctx,
		req.Email,
		req.Password,
		req.Fullname,
	)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to register",
		)
		return
	}

	responseHandler.JSONResponse(
		RegisterResponse{UserID: userID},
		http.StatusCreated,
	)
}
