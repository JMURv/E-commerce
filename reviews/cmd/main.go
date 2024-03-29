package main

import (
	"context"
	"fmt"
	"github.com/JMURv/e-commerce/pkg/discovery"
	"github.com/JMURv/e-commerce/pkg/discovery/consul"
	controller "github.com/JMURv/e-commerce/reviews/internal/controller/review"
	notifygate "github.com/JMURv/e-commerce/reviews/internal/gateway/notifications"
	handler "github.com/JMURv/e-commerce/reviews/internal/handler/grpc"
	//mem "github.com/JMURv/e-commerce/reviews/internal/repository/memory"
	db "github.com/JMURv/e-commerce/reviews/internal/repository/db"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"gopkg.in/yaml.v3"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	pb "github.com/JMURv/e-commerce/api/pb/review"
)

type Config struct {
	Port            int    `yaml:"port"`
	ServiceName     string `yaml:"serviceName"`
	RegistryAddress string `yaml:"registryAddress"`
}

func main() {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("Panic occurred: %v", err)
			os.Exit(1)
		}
	}()

	data, err := os.ReadFile("../dev.config.yaml")
	if err != nil {
		log.Fatalf("error reading configuration file: %v", err)
	}

	var conf Config
	if err = yaml.Unmarshal(data, &conf); err != nil {
		log.Fatalf("error parsing configuration data: %v", err)
	}

	port := conf.Port
	serviceName := conf.ServiceName
	registryAddress := conf.RegistryAddress

	// Setting up registry
	registry, err := consul.NewRegistry(registryAddress)
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
	repo := db.New()
	svc := controller.New(repo, notifygate.New(registry))
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
