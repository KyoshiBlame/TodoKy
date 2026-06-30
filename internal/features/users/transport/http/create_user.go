package users_transport_http

import (
	"encoding/json"
	"fmt"
	"net/http"

	core_logger "github.com/KyoshiBlame/TodoKy/internal/core/logger"
)

type CreateUsersRequest struct {
	FullName 	string 	`json:"full_name"`
	PhoneNumber *string	`json:"phone_number"`
}

type CreateUserResponse struct {
	ID 			int 	`json:"id"`
	Version 	int		`json:"version"`
	FullName 	string	`json:"full_name"`
	PhoneNumber *string `json:"phone_number"`
}

func (h *UsersHTTPHandler) CreateUser(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	log := core_logger.FromContext(ctx)

	log.Debug("invoke CreateUser handler")

	var request CreateUsersRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		fmt.Println("произошла ашибка")
	}
}