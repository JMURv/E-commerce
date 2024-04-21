package main

import (
	"context"
	"fmt"
	pb "github.com/JMURv/e-commerce/api/pb/item"
	kafka "github.com/JMURv/e-commerce/items/internal/broker/kafka"
	redis "github.com/JMURv/e-commerce/items/internal/cache/redis"
	ctrl "github.com/JMURv/e-commerce/items/internal/controller/item"
	usrgate "github.com/JMURv/e-commerce/items/internal/gateway/users"
	handler "github.com/JMURv/e-commerce/items/internal/handler/grpc"
	tracing "github.com/JMURv/e-commerce/items/internal/metrics/jaeger"
	metrics "github.com/JMURv/e-commerce/items/internal/metrics/prometheus"
	"github.com/JMURv/e-commerce/items/internal/repository/db"
	cfg "github.com/JMURv/e-commerce/items/pkg/config"
	"github.com/JMURv/e-commerce/pkg/discovery"
	"github.com/JMURv/e-commerce/pkg/discovery/consul"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const configName = "dev.config"

func main() {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("Panic occurred: %v", err)
			os.Exit(1)
		}
	}()

	// Load configuration
	conf, err := cfg.LoadConfig(configName)
	if err != nil {
		log.Fatalf("failed to parse config: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	port := conf.Port
	serviceName := conf.ServiceName

	// Setting up registry
	registry, err := consul.NewRegistry(conf.RegistryAddr)
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

	// Start metrics and tracing
	metric := metrics.New()
	trace := tracing.New(serviceName, &conf.Jaeger)

	// Setting up main app
	usrGateway := usrgate.New(registry)
	broker := kafka.New(&conf.Kafka)
	cache := redis.New(&conf.Redis)
	repo := db.New(conf)

	svc := ctrl.New(repo, cache, broker, usrGateway)
	h := handler.New(svc)

	srv := metric.ConfigureServerGRPC() // grpc.NewServer()
	pb.RegisterItemServiceServer(srv, h)
	pb.RegisterCategoryServiceServer(srv, h)
	pb.RegisterTagServiceServer(srv, h)
	reflection.Register(srv)

	// Graceful shutdown
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
		<-c
		log.Println("Shutting down gracefully...")

		cancel()
		broker.Close()
		cache.Close()
		metric.Close()
		if err = trace.Close(); err != nil {
			log.Printf("Error closing tracer: %v", err)
		}
		if err = registry.Deregister(ctx, instanceID, serviceName); err != nil {
			log.Printf("Error deregistering service: %v", err)
		}
		srv.GracefulStop()
		os.Exit(0)
	}()

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	log.Printf("%v service is listening", serviceName)
	srv.Serve(lis)
}
