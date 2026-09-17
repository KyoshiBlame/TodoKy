package auth_grpc

import (
	"context"
	"fmt"

	pb "github.com/KyoshiBlame/TodoKy/proto"
)

func (c *AuthClient) Validate(
	ctx context.Context,
	token string,
) (int64, error) {
	resp, err := c.client.Validate(
		ctx,
		&pb.ValidateRequestToken{
			Token: token,
		},
	)

	if err != nil {
		return 0, fmt.Errorf("validate toeken: %v", err)
	}

	if !resp.Valid {
		return 0, fmt.Errorf("token is invalid")
	}

	return resp.UserId, nil
}
