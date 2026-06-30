package users_transport_http

import (
	"encoding/json"
	"fmt"
	"net/http"
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
	var request CreateUsersRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		fmt.Println("произошла ашибка")
	}
}