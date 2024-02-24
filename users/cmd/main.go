package main

import (
	"context"
	"flag"
	"fmt"
	pb "github.com/JMURv/e-commerce/api/pb/user"
	"github.com/JMURv/e-commerce/pkg/discovery"
	"github.com/JMURv/e-commerce/pkg/discovery/consul"
	controller "github.com/JMURv/e-commerce/users/internal/controller/user"
	itmgate "github.com/JMURv/e-commerce/users/internal/gateway/items"
	handler "github.com/JMURv/e-commerce/users/internal/handler/grpc"
	"github.com/JMURv/e-commerce/users/internal/repository/memory"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const serviceName = "users"

func main() {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("Panic occurred: %v", err)
			os.Exit(1)
		}
	}()

	var port int
	flag.IntVar(&port, "port", 50075, "gRPC handler port")
	flag.Parse()

	// Setting up registry
	registry, err := consul.NewRegistry("localhost:8500")
	if err != nil {
		panic(err)
	}

	// Register service
	ctx, cancel := context.WithCancel(context.Background())
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

	// Setting up main app
	repo := memory.New()
	svc := controller.New(repo, *itmgate.New(registry))
	h := handler.New(svc)

	lis, err := net.Listen("tcp", ":50075")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	srv := grpc.NewServer()
	pb.RegisterUserServiceServer(srv, h)
	reflection.Register(srv)

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
		<-c
		log.Println("Shutting down gracefully...")
		cancel()
		registry.Deregister(ctx, instanceID, serviceName)
		srv.GracefulStop()
	}()

	log.Printf("%v service is listening", serviceName)
	srv.Serve(lis)
}
