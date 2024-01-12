package main

import (
	controller "github.com/JMURv/e-commerce/items/internal/controller/item"
	handler "github.com/JMURv/e-commerce/items/internal/handler/grpc"
	"github.com/JMURv/e-commerce/items/internal/repository/memory"
	pb "github.com/JMURv/protos/ecom/item"
	//"google.golang.org/protobuf/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
)

func main() {
	repo := memory.New()
	svc := controller.New(repo)
	h := handler.New(svc)

	lis, err := net.Listen("tcp", ":50080")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	srv := grpc.NewServer()
	pb.RegisterItemServiceServer(srv, h)

	reflection.Register(srv)

	log.Println("Item service is listening")
	srv.Serve(lis)
}
