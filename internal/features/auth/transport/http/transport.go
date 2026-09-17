package auth_transport_http

import auth_grpc "github.com/KyoshiBlame/TodoKy/internal/features/auth/transport/grpc"

type AuthHTTPHandler struct {
	authClient *auth_grpc.AuthClient
}

func NewAuthHTTPHandler(authHandler *auth_grpc.AuthClient) *AuthHTTPHandler {
	return &AuthHTTPHandler{
		authClient: authHandler,
	}
}
