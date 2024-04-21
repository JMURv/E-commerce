package main

import (
	"context"
	"fmt"
	kafka "github.com/JMURv/e-commerce/notifications/internal/broker/kafka"
	redis "github.com/JMURv/e-commerce/notifications/internal/cache/redis"
	ctrl "github.com/JMURv/e-commerce/notifications/internal/controller/notifications"
	hdlr "github.com/JMURv/e-commerce/notifications/internal/handler/grpc"
	tracing "github.com/JMURv/e-commerce/notifications/internal/metrics/jaeger"
	metrics "github.com/JMURv/e-commerce/notifications/internal/metrics/prometheus"
	cfg "github.com/JMURv/e-commerce/notifications/pkg/config"

	//db "github.com/JMURv/e-commerce/notifications/internal/repository/db"
	mem "github.com/JMURv/e-commerce/notifications/internal/repository/memory"
	"github.com/JMURv/e-commerce/pkg/discovery"
	"github.com/JMURv/e-commerce/pkg/discovery/consul"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"google.golang.org/grpc/reflection"

	pb "github.com/JMURv/e-commerce/api/pb/notification"
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
	cache := redis.New(&conf.Redis)
	repo := mem.New(conf)

	svc := ctrl.New(repo, cache)
	h := hdlr.New(svc)

	broker := kafka.New(&conf.Kafka, svc, h)
	go broker.Start()

	// TEST START
	//go func() {
	//	time.Sleep(time.Second * 10)
	//	n := &mdl.Notification{
	//		ID:         1,
	//		Type:       "TEST",
	//		UserID:     2066,
	//		User:       "",
	//		ReceiverID: 3066,
	//		Receiver:   "",
	//		Message:    "Test message",
	//		CreatedAt:  time.Now(),
	//	}
	//	bytes, err := json.Marshal(n)
	//	producer, err := sarama.NewAsyncProducer(conf.Kafka.Addrs, nil)
	//	if err != nil {
	//		log.Fatalf("Error creating Kafka producer: %v", err)
	//	}
	//	producer.Input() <- &sarama.ProducerMessage{
	//		Topic: conf.Kafka.NotificationTopic,
	//		Key:   sarama.StringEncoder(strconv.FormatUint(n.ID, 10)),
	//		Value: sarama.ByteEncoder(bytes),
	//	}
	//	log.Println("Message sent successfully!")
	//}()
	// TEST END

	srv := metric.ConfigureServerGRPC() // grpc.NewServer()
	pb.RegisterNotificationsServer(srv, h)
	pb.RegisterBroadcastServer(srv, h)
	reflection.Register(srv)

	// Start http server for prometheus
	go metric.Start(conf.Port + 1)

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
		srv.Stop()
		os.Exit(0)
	}()

	lis, err := net.Listen("tcp", fmt.Sprintf(":%v", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	log.Printf("%v service is listening", serviceName)
	srv.Serve(lis)
}
