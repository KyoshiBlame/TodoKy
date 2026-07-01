package users_transport_http

import (
	"context"
	"net/http"

	"github.com/KyoshiBlame/TodoKy/internal/core/domain"
	core_http_server "github.com/KyoshiBlame/TodoKy/internal/core/transport/http/server"
)

type UsersHTTPHandler struct {
	UsersService UserService
}

type UserService interface {
	CreateUser (
		ctx context.Context,
		user domain.User,
	) (domain.User, error)
}

func NewUsersHTTPHandler(
	UserService UserService,
) *UsersHTTPHandler {
	return &UsersHTTPHandler{
		UsersService: UserService,
	}
}

func (h *UsersHTTPHandler) Routes() []core_http_server.Route {
	return []core_http_server.Route {
		{
			Method: http.MethodPost,
			Path: "/users",
			Handler: h.CreateUser,
		},
	}
}