package auth_transport_http

import (
	"net/http"

	core_logger "github.com/KyoshiBlame/TodoKy/internal/core/logger"
	core_http_request "github.com/KyoshiBlame/TodoKy/internal/core/transport/http/request"
	core_http_response "github.com/KyoshiBlame/TodoKy/internal/core/transport/http/response"
)

type LoginRequest struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type LoginResponse struct {
	UserID int64 `json:"user_id"`
}

func (h *AuthHTTPHandler) Login(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(
		log,
		rw,
	)

	var req LoginRequest

	if err := core_http_request.DecodeAndValidateRequest(r, &req); err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to decode login request",
		)
		return
	}

	token, userID, err := h.authClient.Login(
		ctx,
		req.Email,
		req.Password,
	)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to login",
		)
		return
	}

	http.SetCookie(rw, &http.Cookie{
		Name:     "access_token",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   86400,
	})

	responseHandler.JSONResponse(
		LoginResponse{UserID: userID},
		http.StatusOK,
	)
}
