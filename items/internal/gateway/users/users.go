package users

import (
	"context"
	"fmt"
	"github.com/JMURv/e-commerce/gateway/pkg/models"
	"github.com/JMURv/e-commerce/pkg/discovery"
)

type Gateway struct {
	registry discovery.Registry
}

func New(registry discovery.Registry) *Gateway {
	return &Gateway{registry}
}

func (g *Gateway) GetUserByID(ctx context.Context, userID string) (*models.User, error) {
	addrs, err := g.registry.ServiceAddresses(ctx, "users")
	fmt.Println(addrs)

	if err != nil {
		return nil, err
	}

	return &models.User{}, nil
}
