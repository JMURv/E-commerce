package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/JMURv/e-commerce/pkg/discovery"
	"github.com/JMURv/e-commerce/pkg/discovery/consul"
	controller "github.com/JMURv/e-commerce/reviews/internal/controller/review"
	handler "github.com/JMURv/e-commerce/reviews/internal/handler/grpc"
	"github.com/JMURv/e-commerce/reviews/internal/repository/memory"
	"log"
	"net"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	pb "github.com/JMURv/e-commerce/api/pb/review"
)

const serviceName = "reviews"

func main() {
	var port int
	flag.IntVar(&port, "port", 50085, "gRPC handler port")
	flag.Parse()

	// Setting up registry
	registry, err := consul.NewRegistry("localhost:8500")
	if err != nil {
		panic(err)
	}

	// Register service
	ctx := context.Background()
	instanceID := discovery.GenerateInstanceID(serviceName)
	if err = registry.Register(ctx, instanceID, serviceName, fmt.Sprintf("localhost:%d", port)); err != nil {
		panic(err)
	}
	go func() {
		for {
			if err := registry.ReportHealthyState(instanceID, serviceName); err != nil {
				log.Println("Failed to report healthy state: " + err.Error())
			}
			time.Sleep(1 * time.Second)
		}
	}()
	defer registry.Deregister(ctx, instanceID, serviceName)

	// Setting up main app
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
