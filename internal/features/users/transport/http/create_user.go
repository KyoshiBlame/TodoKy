package users_transport_http

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/KyoshiBlame/TodoKy/internal/core/domain"
	core_logger "github.com/KyoshiBlame/TodoKy/internal/core/logger"
	core_http_request "github.com/KyoshiBlame/TodoKy/internal/core/transport/http/request"
	core_http_response "github.com/KyoshiBlame/TodoKy/internal/core/transport/http/response"
)

type CreateUsersRequest struct {
	FullName    string  `json:"full_name" validate:"required,min=3,max=100"`
	PhoneNumber *string `json:"phone_number" validate:"omiempty,min=10,max=15,startswith=+"`
}

//omiempty - если не передали в dto то и правила валидации применять не нужно | required - обязательное поле

type CreateUserResponse struct {
	ID          int     `json:"id"`
	Version     int     `json:"version"`
	FullName    string  `json:"full_name"`
	PhoneNumber *string `json:"phone_number"`
}

func (h *UsersHTTPHandler) CreateUser(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHander := core_http_response.NewHTTPResponseHandler(log, w)

	var request CreateUsersRequest

	if err := core_http_request.DecodeAndValidateRequest(r, &request); err != nil {
		responseHander.ErrorResponse(err, "failed to decode and validate request")

		return
	}

	h.UsersService.CreateUser(ctx, domainFromDTO())
}

func domainFromDTO(dto CreateUsersRequest) domain.User {
	return domain.NewUser(dto.FullName, dto.PhoneNumber)
}