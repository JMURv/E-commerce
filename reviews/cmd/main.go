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
	"os"
	"os/signal"
	"syscall"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	pb "github.com/JMURv/e-commerce/api/pb/review"
)

const serviceName = "reviews"
const RegistryAddress = "localhost:8500"

func main() {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("Panic occurred: %v", err)
			os.Exit(1)
		}
	}()

	var port int
	flag.IntVar(&port, "port", 50085, "gRPC handler port")
	flag.Parse()

	// Setting up registry
	registry, err := consul.NewRegistry(RegistryAddress)
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

	// Setting up main app
	repo := memory.New()
	svc := controller.New(repo)
	h := handler.New(svc)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%v", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	srv := grpc.NewServer()
	pb.RegisterReviewServiceServer(srv, h)

	reflection.Register(srv)

	// Setting up signal handling for graceful shutdown
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
		<-c
		log.Println("Shutting down gracefully...")
		registry.Deregister(ctx, instanceID, serviceName)
		os.Exit(0)
	}()

	log.Printf("%v service is listening", serviceName)
	srv.Serve(lis)

}
