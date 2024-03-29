// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v4.25.1
// source: review/review.proto

package review

import (
	context "context"
	common "github.com/JMURv/e-commerce/api/pb/common"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

const (
	ReviewService_GetReviewByID_FullMethodName           = "/review.ReviewService/GetReviewByID"
	ReviewService_GetReviewsByUserID_FullMethodName      = "/review.ReviewService/GetReviewsByUserID"
	ReviewService_AggregateUserRatingByID_FullMethodName = "/review.ReviewService/AggregateUserRatingByID"
	ReviewService_CreateReview_FullMethodName            = "/review.ReviewService/CreateReview"
	ReviewService_UpdateReview_FullMethodName            = "/review.ReviewService/UpdateReview"
	ReviewService_DeleteReview_FullMethodName            = "/review.ReviewService/DeleteReview"
)

// ReviewServiceClient is the client API for ReviewService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ReviewServiceClient interface {
	GetReviewByID(ctx context.Context, in *GetReviewByIDRequest, opts ...grpc.CallOption) (*common.Review, error)
	GetReviewsByUserID(ctx context.Context, in *ByUserIDRequest, opts ...grpc.CallOption) (*ListReviewResponse, error)
	AggregateUserRatingByID(ctx context.Context, in *ByUserIDRequest, opts ...grpc.CallOption) (*AggregateUserRatingByIDResponse, error)
	CreateReview(ctx context.Context, in *CreateReviewRequest, opts ...grpc.CallOption) (*common.Review, error)
	UpdateReview(ctx context.Context, in *UpdateReviewRequest, opts ...grpc.CallOption) (*common.Review, error)
	DeleteReview(ctx context.Context, in *DeleteReviewRequest, opts ...grpc.CallOption) (*EmptyResponse, error)
}

type reviewServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewReviewServiceClient(cc grpc.ClientConnInterface) ReviewServiceClient {
	return &reviewServiceClient{cc}
}

func (c *reviewServiceClient) GetReviewByID(ctx context.Context, in *GetReviewByIDRequest, opts ...grpc.CallOption) (*common.Review, error) {
	out := new(common.Review)
	err := c.cc.Invoke(ctx, ReviewService_GetReviewByID_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *reviewServiceClient) GetReviewsByUserID(ctx context.Context, in *ByUserIDRequest, opts ...grpc.CallOption) (*ListReviewResponse, error) {
	out := new(ListReviewResponse)
	err := c.cc.Invoke(ctx, ReviewService_GetReviewsByUserID_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *reviewServiceClient) AggregateUserRatingByID(ctx context.Context, in *ByUserIDRequest, opts ...grpc.CallOption) (*AggregateUserRatingByIDResponse, error) {
	out := new(AggregateUserRatingByIDResponse)
	err := c.cc.Invoke(ctx, ReviewService_AggregateUserRatingByID_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *reviewServiceClient) CreateReview(ctx context.Context, in *CreateReviewRequest, opts ...grpc.CallOption) (*common.Review, error) {
	out := new(common.Review)
	err := c.cc.Invoke(ctx, ReviewService_CreateReview_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *reviewServiceClient) UpdateReview(ctx context.Context, in *UpdateReviewRequest, opts ...grpc.CallOption) (*common.Review, error) {
	out := new(common.Review)
	err := c.cc.Invoke(ctx, ReviewService_UpdateReview_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *reviewServiceClient) DeleteReview(ctx context.Context, in *DeleteReviewRequest, opts ...grpc.CallOption) (*EmptyResponse, error) {
	out := new(EmptyResponse)
	err := c.cc.Invoke(ctx, ReviewService_DeleteReview_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ReviewServiceServer is the server API for ReviewService service.
// All implementations must embed UnimplementedReviewServiceServer
// for forward compatibility
type ReviewServiceServer interface {
	GetReviewByID(context.Context, *GetReviewByIDRequest) (*common.Review, error)
	GetReviewsByUserID(context.Context, *ByUserIDRequest) (*ListReviewResponse, error)
	AggregateUserRatingByID(context.Context, *ByUserIDRequest) (*AggregateUserRatingByIDResponse, error)
	CreateReview(context.Context, *CreateReviewRequest) (*common.Review, error)
	UpdateReview(context.Context, *UpdateReviewRequest) (*common.Review, error)
	DeleteReview(context.Context, *DeleteReviewRequest) (*EmptyResponse, error)
	mustEmbedUnimplementedReviewServiceServer()
}

// UnimplementedReviewServiceServer must be embedded to have forward compatible implementations.
type UnimplementedReviewServiceServer struct {
}

func (UnimplementedReviewServiceServer) GetReviewByID(context.Context, *GetReviewByIDRequest) (*common.Review, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetReviewByID not implemented")
}
func (UnimplementedReviewServiceServer) GetReviewsByUserID(context.Context, *ByUserIDRequest) (*ListReviewResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetReviewsByUserID not implemented")
}
func (UnimplementedReviewServiceServer) AggregateUserRatingByID(context.Context, *ByUserIDRequest) (*AggregateUserRatingByIDResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AggregateUserRatingByID not implemented")
}
func (UnimplementedReviewServiceServer) CreateReview(context.Context, *CreateReviewRequest) (*common.Review, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateReview not implemented")
}
func (UnimplementedReviewServiceServer) UpdateReview(context.Context, *UpdateReviewRequest) (*common.Review, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateReview not implemented")
}
func (UnimplementedReviewServiceServer) DeleteReview(context.Context, *DeleteReviewRequest) (*EmptyResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteReview not implemented")
}
func (UnimplementedReviewServiceServer) mustEmbedUnimplementedReviewServiceServer() {}

// UnsafeReviewServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ReviewServiceServer will
// result in compilation errors.
type UnsafeReviewServiceServer interface {
	mustEmbedUnimplementedReviewServiceServer()
}

func RegisterReviewServiceServer(s grpc.ServiceRegistrar, srv ReviewServiceServer) {
	s.RegisterService(&ReviewService_ServiceDesc, srv)
}

func _ReviewService_GetReviewByID_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetReviewByIDRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ReviewServiceServer).GetReviewByID(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ReviewService_GetReviewByID_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ReviewServiceServer).GetReviewByID(ctx, req.(*GetReviewByIDRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ReviewService_GetReviewsByUserID_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ByUserIDRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ReviewServiceServer).GetReviewsByUserID(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ReviewService_GetReviewsByUserID_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ReviewServiceServer).GetReviewsByUserID(ctx, req.(*ByUserIDRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ReviewService_AggregateUserRatingByID_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ByUserIDRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ReviewServiceServer).AggregateUserRatingByID(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ReviewService_AggregateUserRatingByID_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ReviewServiceServer).AggregateUserRatingByID(ctx, req.(*ByUserIDRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ReviewService_CreateReview_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateReviewRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ReviewServiceServer).CreateReview(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ReviewService_CreateReview_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ReviewServiceServer).CreateReview(ctx, req.(*CreateReviewRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ReviewService_UpdateReview_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateReviewRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ReviewServiceServer).UpdateReview(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ReviewService_UpdateReview_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ReviewServiceServer).UpdateReview(ctx, req.(*UpdateReviewRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ReviewService_DeleteReview_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteReviewRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ReviewServiceServer).DeleteReview(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ReviewService_DeleteReview_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ReviewServiceServer).DeleteReview(ctx, req.(*DeleteReviewRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// ReviewService_ServiceDesc is the grpc.ServiceDesc for ReviewService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ReviewService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "review.ReviewService",
	HandlerType: (*ReviewServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetReviewByID",
			Handler:    _ReviewService_GetReviewByID_Handler,
		},
		{
			MethodName: "GetReviewsByUserID",
			Handler:    _ReviewService_GetReviewsByUserID_Handler,
		},
		{
			MethodName: "AggregateUserRatingByID",
			Handler:    _ReviewService_AggregateUserRatingByID_Handler,
		},
		{
			MethodName: "CreateReview",
			Handler:    _ReviewService_CreateReview_Handler,
		},
		{
			MethodName: "UpdateReview",
			Handler:    _ReviewService_UpdateReview_Handler,
		},
		{
			MethodName: "DeleteReview",
			Handler:    _ReviewService_DeleteReview_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "review/review.proto",
}
