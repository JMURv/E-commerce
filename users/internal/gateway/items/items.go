package items

import (
	"context"
	"fmt"
	"github.com/JMURv/e-commerce/gateway/pkg/models"
	"github.com/JMURv/e-commerce/pkg/discovery"
	"github.com/JMURv/e-commerce/users/internal/grpcutil"
	pb "github.com/JMURv/protos/ecom/item"
)

type Gateway struct {
	registry discovery.Registry
}

func New(registry discovery.Registry) *Gateway {
	return &Gateway{registry}
}

func (g *Gateway) ListUserItemsByID(ctx context.Context, userID uint64) ([]*models.Item, error) {
	conn, err := grpcutil.ServiceConnection(ctx, "items", g.registry)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := pb.NewItemServiceClient(conn)
	resp, err := client.ListUserItemsByID(ctx, &pb.GetUserByIDRequest{UserId: userID})
	if err != nil {
		return nil, err
	}

	fmt.Println(resp)
	return []*models.Item{}, nil
}
