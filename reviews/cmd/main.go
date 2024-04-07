package main

import (
	"context"
	"fmt"
	"github.com/JMURv/e-commerce/pkg/discovery"
	"github.com/JMURv/e-commerce/pkg/discovery/consul"
	cache "github.com/JMURv/e-commerce/reviews/internal/cache"
	controller "github.com/JMURv/e-commerce/reviews/internal/controller/review"
	notifygate "github.com/JMURv/e-commerce/reviews/internal/gateway/notifications"
	handler "github.com/JMURv/e-commerce/reviews/internal/handler/grpc"
	"github.com/go-redis/redis/v8"

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
	Port         int    `yaml:"port"`
	ServiceName  string `yaml:"serviceName"`
	RegistryAddr string `yaml:"registryAddr"`
	RedisAddr    string `yaml:"redisAddr"`
	RedisPass    string `yaml:"redisPass"`
}

func loadConfig() (*Config, error) {
	var conf Config

	data, err := os.ReadFile("../dev.config.yaml")
	if err != nil {
		return nil, err
	}

	if err = yaml.Unmarshal(data, &conf); err != nil {
		return nil, err
	}
	return &conf, nil
}

func main() {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("Panic occurred: %v", err)
			os.Exit(1)
		}
	}()

	conf, err := loadConfig()
	if err != nil {
		log.Fatalf("failed to parse config: %v", err)
	}

	port := conf.Port
	serviceName := conf.ServiceName
	registryAddress := conf.RegistryAddr

	// Setting up redis TODO: Убрать подклбчение в реализацию
	redisCli := redis.NewClient(&redis.Options{
		Addr:     conf.RedisAddr,
		Password: conf.RedisPass,
		DB:       0,
	})
	pong, err := redisCli.Ping(context.Background()).Result()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	fmt.Println("Connected to Redis:", pong)

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
	svc := controller.New(repo, cache.New(redisCli), notifygate.New(registry))
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
		redisCli.Close()
		os.Exit(0)
	}()

	log.Printf("%v service is listening", serviceName)
	srv.Serve(lis)
}
