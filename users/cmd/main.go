package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/JMURv/e-commerce/pkg/discovery"
	"github.com/JMURv/e-commerce/pkg/discovery/consul"
	controller "github.com/JMURv/e-commerce/users/internal/controller/user"
	itmgate "github.com/JMURv/e-commerce/users/internal/gateway/items"
	handler "github.com/JMURv/e-commerce/users/internal/handler/grpc"
	"github.com/JMURv/e-commerce/users/internal/repository/memory"
	pb "github.com/JMURv/protos/ecom/user"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"time"
)

const serviceName = "users"

func main() {
	var port int
	flag.IntVar(&port, "port", 50075, "gRPC handler port")
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

	// Setting up other services
	itemGateway := itmgate.New(registry)

	repo := memory.New()
	svc := controller.New(repo, *itemGateway)
	h := handler.New(svc)

	lis, err := net.Listen("tcp", ":50075")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	srv := grpc.NewServer()
	pb.RegisterUserServiceServer(srv, h)

	reflection.Register(srv)

	log.Println("User service is listening")
	srv.Serve(lis)
}
