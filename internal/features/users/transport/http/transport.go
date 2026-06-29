package users_transport_http

type UsersHTTPHandler struct {
	UsersService UserService
}

type UserService interface {

}

func NewUsersHTTPHandler(
	UserService UserService,
) *UsersHTTPHandler {
	return &UsersHTTPHandler{
		UsersService: UserService,
	}
}	