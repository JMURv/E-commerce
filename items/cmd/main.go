package main

import (
	"context"
	"flag"
	"fmt"
	controller "github.com/JMURv/e-commerce/items/internal/controller/item"
	usrgate "github.com/JMURv/e-commerce/items/internal/gateway/users"
	handler "github.com/JMURv/e-commerce/items/internal/handler/grpc"
	"github.com/JMURv/e-commerce/items/internal/repository/memory"
	"github.com/JMURv/e-commerce/pkg/discovery"
	"github.com/JMURv/e-commerce/pkg/discovery/consul"
	pb "github.com/JMURv/protos/ecom/item"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"time"
)

const serviceName = "items"

func main() {
	var port int
	flag.IntVar(&port, "port", 50080, "gRPC handler port")
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
	usrGateway := usrgate.New(registry)

	// Setting up main app
	repo := memory.New()
	svc := controller.New(repo, *usrGateway)
	h := handler.New(svc)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	srv := grpc.NewServer()
	pb.RegisterItemServiceServer(srv, h)

	reflection.Register(srv)

	log.Println("Item service is listening")
	srv.Serve(lis)
}
