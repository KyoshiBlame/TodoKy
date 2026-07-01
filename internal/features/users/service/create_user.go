package users_service

import (
	"context"

	"github.com/KyoshiBlame/TodoKy/internal/core/domain"
)


func (s *UsersService) CreateUser(
	ctx context.Context,
	user domain.User, 
) (domain.User, error) {
	
}