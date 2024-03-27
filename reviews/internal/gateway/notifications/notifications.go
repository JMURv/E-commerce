package notifications

import (
	"context"
	pb "github.com/JMURv/e-commerce/api/pb/notification"
	"github.com/JMURv/e-commerce/pkg/discovery"
	"github.com/JMURv/e-commerce/reviews/internal/grpcutil"
	"github.com/JMURv/e-commerce/reviews/pkg/model"
)

type Gateway struct {
	registry discovery.Registry
}

func New(registry discovery.Registry) *Gateway {
	return &Gateway{registry}
}

func (g *Gateway) CreateReviewNotification(ctx context.Context, rev *model.Review) error {
	conn, err := grpcutil.ServiceConnection(ctx, "notifications", g.registry)
	if err != nil {
		return err
	}
	defer conn.Close()

	client := pb.NewNotificationsClient(conn)
	_, err = client.CreateNotification(ctx, &pb.Notification{
		Type:       "new_review",
		UserId:     rev.UserID,
		ReceiverId: rev.ReviewedUserID,
		Message:    "New review recieved!",
	})
	if err != nil {
		return err
	}
	return nil
}
