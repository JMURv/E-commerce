package memory

import (
	"context"
	"fmt"
	pb "github.com/JMURv/e-commerce/api/pb/item"
	"github.com/JMURv/e-commerce/items/pkg/model"
	"github.com/JMURv/e-commerce/pkg/discovery"
	"github.com/JMURv/e-commerce/users/internal/grpcutil"
)

type Gateway struct {
	registry discovery.Registry
}

func New(registry discovery.Registry) *Gateway {
	return &Gateway{registry}
}

func (g *Gateway) ListUserItemsByID(ctx context.Context, userID uint64) ([]*model.Item, error) {
	conn, err := grpcutil.ServiceConnection(ctx, "items", g.registry)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := pb.NewItemServiceClient(conn)
	resp, err := client.ListUserItemsByID(ctx, &pb.ListUserItemsByIDRequest{UserId: userID})
	if err != nil {
		return nil, err
	}

	fmt.Println(resp)
	return []*model.Item{}, nil
}
