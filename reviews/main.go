package main

import (
	"e-commerce/reviews/model"
	"golang.org/x/net/context"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	common "github.com/JMURv/protos/ecom/common"
	pb "github.com/JMURv/protos/ecom/review"
)

type reviewServer struct {
	pb.ReviewServiceServer
}

func (s *reviewServer) GetReviewByID(ctx context.Context, req *pb.GetReviewByIDRequest) (*common.Review, error) {
	reviewID := req.GetReviewId()

	review, err := model.GetReviewByID(reviewID)
	if err != nil {
		log.Printf("Error getting review by ID: %v", err)
		return &common.Review{}, err
	}

	response := &common.Review{
		ReviewId:       review.ID,
		UserId:         review.UserID,
		ItemId:         review.ItemID,
		ReviewedUserId: review.ReviewedUserID,
		Advantages:     review.Advantages,
		Disadvantages:  review.Disadvantages,
		ReviewText:     review.ReviewText,
		Rating:         review.Rating,
	}

	return response, nil
}

func (s *reviewServer) CreateReview(ctx context.Context, req *pb.CreateReviewRequest) (*common.Review, error) {
	newReview := &model.Review{
		UserID:         req.GetUserId(),
		ItemID:         req.GetItemId(),
		ReviewedUserID: req.GetReviewedUserId(),
		Advantages:     req.GetAdvantages(),
		Disadvantages:  req.GetDisadvantages(),
		ReviewText:     req.GetReviewText(),
		Rating:         req.GetRating(),
	}

	// TODO: Fetch all related fields: user, items etc.
	review, err := newReview.CreateReview()
	if err != nil {
		return nil, err
	}

	response := &common.Review{
		ReviewId:       review.ID,
		UserId:         review.UserID,
		ItemId:         review.ItemID,
		ReviewedUserId: review.ReviewedUserID,
		Advantages:     review.Advantages,
		Disadvantages:  review.Disadvantages,
		ReviewText:     review.ReviewText,
		Rating:         review.Rating,
	}

	return response, nil
}

func (s *reviewServer) UpdateReview(ctx context.Context, req *pb.UpdateReviewRequest) (*common.Review, error) {
	reviewID := req.GetReviewId()

	newData := &model.Review{
		UserID:         req.GetUserId(),
		ItemID:         req.GetItemId(),
		ReviewedUserID: req.GetReviewedUserId(),
		Advantages:     req.GetAdvantages(),
		Disadvantages:  req.GetDisadvantages(),
		ReviewText:     req.GetReviewText(),
		Rating:         req.GetRating(),
	}

	updatedReview, err := model.UpdateReview(reviewID, newData)
	if err != nil {
		return nil, err
	}

	response := &common.Review{
		ReviewId:       updatedReview.ID,
		UserId:         updatedReview.UserID,
		ItemId:         updatedReview.ItemID,
		ReviewedUserId: updatedReview.ReviewedUserID,
		Advantages:     updatedReview.Advantages,
		Disadvantages:  updatedReview.Disadvantages,
		ReviewText:     updatedReview.ReviewText,
		Rating:         updatedReview.Rating,
	}

	return response, nil
}

func (s *reviewServer) DeleteReview(ctx context.Context, req *pb.DeleteReviewRequest) (*pb.EmptyResponse, error) {
	reviewID := req.GetReviewId()

	err := model.DeleteReview(reviewID)
	if err != nil {
		return nil, err
	}

	return &pb.EmptyResponse{}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50070")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	server := grpc.NewServer()
	pb.RegisterReviewServiceServer(server, &reviewServer{})

	reflection.Register(server)

	log.Println("Review service is listening")
	if err := server.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
