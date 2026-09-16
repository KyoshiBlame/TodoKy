package auth_grpc

import (
	"fmt"

	pb "github.com/KyoshiBlame/TodoKy/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type authClient struct {
	client pb.AuthServiceClient
	conn   *grpc.ClientConn
}

func NewAuthClient(authSeviceAddr string) (*authClient, error) {
	conn, err := grpc.NewClient(
		authSeviceAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)

	if err != nil {
		return nil, fmt.Errorf("failed to connect auth-service: %v", err)
	}

	client := pb.NewAuthServiceClient(conn)
	return &authClient{client: client, conn: conn}, nil
}

func (c *authClient) Close() error {
	return c.conn.Close()
}
