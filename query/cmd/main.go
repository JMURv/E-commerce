package main

import (
	"context"
	"fmt"
	"github.com/JMURv/e-commerce/pkg/discovery"
	"github.com/JMURv/e-commerce/pkg/discovery/consul"
	"github.com/JMURv/e-commerce/query/internal/broker/kafka"
	redis "github.com/JMURv/e-commerce/query/internal/cache"
	ctrl "github.com/JMURv/e-commerce/query/internal/controller"
	"github.com/JMURv/e-commerce/query/internal/repository/db"
	config "github.com/JMURv/e-commerce/query/pkg/config"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("Panic occurred: %v", err)
			os.Exit(1)
		}
	}()

	conf, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("failed to parse config: %v", err)
	}

	ctx := context.Background()
	port := conf.Port
	serviceName := conf.ServiceName
	registryAddress := conf.RegistryAddr

	// Setting up registry
	registry, err := consul.NewRegistry(registryAddress)
	if err != nil {
		panic(err)
	}

	// Register service
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
	broker := kafka.New(conf.KafkaAddr)
	cache := redis.New(conf.RedisAddr, conf.RedisPass)
	repo := db.New(conf.MongoAddr)
	svc := ctrl.New(repo, cache)
	//h := handler.New(svc)
	//
	//lis, err := net.Listen("tcp", fmt.Sprintf(":%v", port))
	//if err != nil {
	//	log.Fatalf("failed to listen: %v", err)
	//}

	//srv := grpc.NewServer()
	//pb.RegisterReviewServiceServer(srv, h)

	//reflection.Register(srv)

	// Setting up signal handling for graceful shutdown
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
		<-c
		log.Println("Shutting down gracefully...")

		repo.Close()
		cache.Close()
		broker.Close()
		if err = registry.Deregister(ctx, instanceID, serviceName); err != nil {
			log.Printf("Error deregistering service: %v", err)
		}

		os.Exit(0)
	}()

	//log.Printf("%v service is listening", serviceName)
	//srv.Serve(lis)
}
