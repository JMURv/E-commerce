package main

import (
	"context"
	"flag"
	"fmt"
	pb "github.com/JMURv/e-commerce/api/pb/item"
	controller "github.com/JMURv/e-commerce/items/internal/controller/item"
	usrgate "github.com/JMURv/e-commerce/items/internal/gateway/users"
	handler "github.com/JMURv/e-commerce/items/internal/handler/grpc"
	"github.com/JMURv/e-commerce/items/internal/repository/memory"
	"github.com/JMURv/e-commerce/pkg/discovery"
	"github.com/JMURv/e-commerce/pkg/discovery/consul"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const serviceName = "items"

func main() {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("Panic occurred: %v", err)
			os.Exit(1)
		}
	}()

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
	pb.RegisterCategoryServiceServer(srv, h)
	pb.RegisterTagServiceServer(srv, h)

	reflection.Register(srv)

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
		<-c
		log.Println("Shutting down gracefully...")
		registry.Deregister(ctx, instanceID, serviceName)
		os.Exit(0)
	}()

	log.Println("Item service is listening")
	srv.Serve(lis)
}
