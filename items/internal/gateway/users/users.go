package users

import (
	"context"
	"fmt"
	cm "github.com/JMURv/e-commerce/api/pb/common"
	pb "github.com/JMURv/e-commerce/api/pb/user"
	"github.com/JMURv/e-commerce/items/internal/grpcutil"
	"github.com/JMURv/e-commerce/pkg/discovery"
)

type Gateway struct {
	registry discovery.Registry
}

func New(registry discovery.Registry) *Gateway {
	return &Gateway{registry}
}

func (g *Gateway) GetUserByID(ctx context.Context, userID uint64) (*cm.User, error) {
	conn, err := grpcutil.ServiceConnection(ctx, "users", g.registry)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := pb.NewUserServiceClient(conn)
	resp, err := client.GetUserByID(ctx, &pb.GetUserByIDRequest{UserId: userID})
	if err != nil {
		return nil, err
	}

	fmt.Println(resp)
	return &cm.User{}, nil
}
