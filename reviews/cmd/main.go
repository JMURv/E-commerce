package main

import (
	controller "github.com/JMURv/e-commerce/reviews/internal/controller/review"
	handler "github.com/JMURv/e-commerce/reviews/internal/handler/grpc"
	"github.com/JMURv/e-commerce/reviews/internal/repository/memory"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	pb "github.com/JMURv/e-commerce/api/pb/review"
)

func main() {
	repo := memory.New()
	svc := controller.New(repo)
	h := handler.New(svc)

	lis, err := net.Listen("tcp", ":50085")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	srv := grpc.NewServer()
	pb.RegisterReviewServiceServer(srv, h)

	reflection.Register(srv)

	log.Println("Review service is listening")
	srv.Serve(lis)
}
