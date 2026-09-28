package auth_grpc

import (
	"context"
	"fmt"

	pb "github.com/KyoshiBlame/TodoKy/proto"
)

func (c *AuthClient) Login(
	ctx context.Context,
	email, password string,
) (string, int64, error) {
	resp, err := c.client.Login(
		ctx,
		&pb.LoginRequest{
			Email:    email,
			Password: password,
		},
	)

	if err != nil {
		return "", 0, fmt.Errorf("failed to login: %v", err)
	}

	return resp.AccesToken, resp.UserId, nil
}
