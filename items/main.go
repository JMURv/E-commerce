package main

import (
	"e-commerce/items/model"
	common "github.com/JMURv/protos/ecom/common"
	pb "github.com/JMURv/protos/ecom/item"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
)

type itemServer struct {
	pb.ItemServiceServer
}

func (s *itemServer) ListItem(ctx context.Context, req *pb.EmptyRequest) (*pb.ListItemResponse, error) {
	return &pb.ListItemResponse{}, nil
}

func (s *itemServer) GetItemByID(ctx context.Context, req *pb.GetItemByIDRequest) (*common.Item, error) {
	itemID := req.GetItemId()

	i, err := model.GetItemByID(itemID)
	if err != nil {
		log.Printf("Error getting item by ID: %v", err)
		return &common.Item{}, err
	}

	return responseItem(i)
}

func (s *itemServer) CreateItem(ctx context.Context, req *pb.CreateItemRequest) (*common.Item, error) {
	newItem := &model.Item{
		UserID:      req.GetUserId(),
		CategoryID:  req.GetCategoryId(),
		Name:        req.GetName(),
		Description: req.GetDescription(),
		Price:       req.GetPrice(),
		Quantity:    req.GetQuantity(),
		Tags:        tagsToModel(req.GetTags()),
	}

	i, err := newItem.CreateItem()
	if err != nil {
		return nil, err
	}

	return responseItem(i)
}

func (s *itemServer) UpdateItem(ctx context.Context, req *pb.UpdateItemRequest) (*common.Item, error) {
	itemID := req.GetItemId()

	newData := &model.Item{
		Name:        req.GetName(),
		Description: req.GetDescription(),
		Price:       req.GetPrice(),
		CategoryID:  req.GetCategoryId(),
		UserID:      req.GetUserId(),
		Tags:        tagsToModel(req.GetTags()),
		Quantity:    req.GetQuantity(),
	}

	i, err := model.UpdateItem(itemID, newData)
	if err != nil {
		return nil, err
	}

	return responseItem(i)
}

func (s *itemServer) DeleteItem(ctx context.Context, req *pb.DeleteItemRequest) (*pb.EmptyResponse, error) {
	itemID := req.GetItemId()

	err := model.DeleteItem(itemID)
	if err != nil {
		return nil, err
	}

	return &pb.EmptyResponse{}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50080")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	server := grpc.NewServer()
	pb.RegisterItemServiceServer(server, &itemServer{})

	reflection.Register(server)

	log.Println("Item service is listening")
	if err := server.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
