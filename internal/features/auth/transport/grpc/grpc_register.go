package auth_grpc

import (
	"context"
	"fmt"

	pb "github.com/KyoshiBlame/TodoKy/proto"
)

func (c *authClient) Register(
	ctx context.Context,
	email, password, fullName string,
) (int64, error) {
	resp, err := c.client.Register(ctx, &pb.RegisterRequest{
		Email:    email,
		Password: password,
		FullName: fullName,
	})

	if err != nil {
		return 0, fmt.Errorf("failed to register: %v", err)
	}

	return resp.UserId, nil
}
