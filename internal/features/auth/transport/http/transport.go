package auth_transport_http

import (
	"net/http"

	core_http_server "github.com/KyoshiBlame/TodoKy/internal/core/transport/http/server"
	auth_grpc "github.com/KyoshiBlame/TodoKy/internal/features/auth/transport/grpc"
)

type AuthHTTPHandler struct {
	authClient *auth_grpc.AuthClient
}

func NewAuthHTTPHandler(authHandler *auth_grpc.AuthClient) *AuthHTTPHandler {
	return &AuthHTTPHandler{
		authClient: authHandler,
	}
}

func (h *AuthHTTPHandler) Routes() []core_http_server.Route {
	return []core_http_server.Route{
		{
			Method:  http.MethodPost,
			Path:    "auth/register",
			Handler: h.Register,
		},
		{
			Method:  http.MethodPost,
			Path:    "auth/Login",
			Handler: h.Login,
		},
	}
}
