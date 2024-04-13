// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v4.25.1
// source: chat/chat.proto

package chat

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

const (
	Broadcast_CreateStream_FullMethodName  = "/chat.Broadcast/CreateStream"
	Broadcast_CreateMessage_FullMethodName = "/chat.Broadcast/CreateMessage"
	Broadcast_UpdateMessage_FullMethodName = "/chat.Broadcast/UpdateMessage"
	Broadcast_DeleteMessage_FullMethodName = "/chat.Broadcast/DeleteMessage"
)

// BroadcastClient is the client API for Broadcast service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type BroadcastClient interface {
	CreateStream(ctx context.Context, in *Connect, opts ...grpc.CallOption) (Broadcast_CreateStreamClient, error)
	CreateMessage(ctx context.Context, in *CreateMessageRequest, opts ...grpc.CallOption) (*EmptyResponse, error)
	UpdateMessage(ctx context.Context, in *UpdateMessageRequest, opts ...grpc.CallOption) (*EmptyResponse, error)
	DeleteMessage(ctx context.Context, in *DeleteMessageRequest, opts ...grpc.CallOption) (*EmptyResponse, error)
}

type broadcastClient struct {
	cc grpc.ClientConnInterface
}

func NewBroadcastClient(cc grpc.ClientConnInterface) BroadcastClient {
	return &broadcastClient{cc}
}

func (c *broadcastClient) CreateStream(ctx context.Context, in *Connect, opts ...grpc.CallOption) (Broadcast_CreateStreamClient, error) {
	stream, err := c.cc.NewStream(ctx, &Broadcast_ServiceDesc.Streams[0], Broadcast_CreateStream_FullMethodName, opts...)
	if err != nil {
		return nil, err
	}
	x := &broadcastCreateStreamClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type Broadcast_CreateStreamClient interface {
	Recv() (*Message, error)
	grpc.ClientStream
}

type broadcastCreateStreamClient struct {
	grpc.ClientStream
}

func (x *broadcastCreateStreamClient) Recv() (*Message, error) {
	m := new(Message)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *broadcastClient) CreateMessage(ctx context.Context, in *CreateMessageRequest, opts ...grpc.CallOption) (*EmptyResponse, error) {
	out := new(EmptyResponse)
	err := c.cc.Invoke(ctx, Broadcast_CreateMessage_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *broadcastClient) UpdateMessage(ctx context.Context, in *UpdateMessageRequest, opts ...grpc.CallOption) (*EmptyResponse, error) {
	out := new(EmptyResponse)
	err := c.cc.Invoke(ctx, Broadcast_UpdateMessage_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *broadcastClient) DeleteMessage(ctx context.Context, in *DeleteMessageRequest, opts ...grpc.CallOption) (*EmptyResponse, error) {
	out := new(EmptyResponse)
	err := c.cc.Invoke(ctx, Broadcast_DeleteMessage_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// BroadcastServer is the server API for Broadcast service.
// All implementations must embed UnimplementedBroadcastServer
// for forward compatibility
type BroadcastServer interface {
	CreateStream(*Connect, Broadcast_CreateStreamServer) error
	CreateMessage(context.Context, *CreateMessageRequest) (*EmptyResponse, error)
	UpdateMessage(context.Context, *UpdateMessageRequest) (*EmptyResponse, error)
	DeleteMessage(context.Context, *DeleteMessageRequest) (*EmptyResponse, error)
	mustEmbedUnimplementedBroadcastServer()
}

// UnimplementedBroadcastServer must be embedded to have forward compatible implementations.
type UnimplementedBroadcastServer struct {
}

func (UnimplementedBroadcastServer) CreateStream(*Connect, Broadcast_CreateStreamServer) error {
	return status.Errorf(codes.Unimplemented, "method CreateStream not implemented")
}
func (UnimplementedBroadcastServer) CreateMessage(context.Context, *CreateMessageRequest) (*EmptyResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateMessage not implemented")
}
func (UnimplementedBroadcastServer) UpdateMessage(context.Context, *UpdateMessageRequest) (*EmptyResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateMessage not implemented")
}
func (UnimplementedBroadcastServer) DeleteMessage(context.Context, *DeleteMessageRequest) (*EmptyResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteMessage not implemented")
}
func (UnimplementedBroadcastServer) mustEmbedUnimplementedBroadcastServer() {}

// UnsafeBroadcastServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to BroadcastServer will
// result in compilation errors.
type UnsafeBroadcastServer interface {
	mustEmbedUnimplementedBroadcastServer()
}

func RegisterBroadcastServer(s grpc.ServiceRegistrar, srv BroadcastServer) {
	s.RegisterService(&Broadcast_ServiceDesc, srv)
}

func _Broadcast_CreateStream_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(Connect)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(BroadcastServer).CreateStream(m, &broadcastCreateStreamServer{stream})
}

type Broadcast_CreateStreamServer interface {
	Send(*Message) error
	grpc.ServerStream
}

type broadcastCreateStreamServer struct {
	grpc.ServerStream
}

func (x *broadcastCreateStreamServer) Send(m *Message) error {
	return x.ServerStream.SendMsg(m)
}

func _Broadcast_CreateMessage_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateMessageRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BroadcastServer).CreateMessage(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Broadcast_CreateMessage_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BroadcastServer).CreateMessage(ctx, req.(*CreateMessageRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Broadcast_UpdateMessage_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateMessageRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BroadcastServer).UpdateMessage(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Broadcast_UpdateMessage_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BroadcastServer).UpdateMessage(ctx, req.(*UpdateMessageRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Broadcast_DeleteMessage_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteMessageRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BroadcastServer).DeleteMessage(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Broadcast_DeleteMessage_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BroadcastServer).DeleteMessage(ctx, req.(*DeleteMessageRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Broadcast_ServiceDesc is the grpc.ServiceDesc for Broadcast service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Broadcast_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "chat.Broadcast",
	HandlerType: (*BroadcastServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateMessage",
			Handler:    _Broadcast_CreateMessage_Handler,
		},
		{
			MethodName: "UpdateMessage",
			Handler:    _Broadcast_UpdateMessage_Handler,
		},
		{
			MethodName: "DeleteMessage",
			Handler:    _Broadcast_DeleteMessage_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "CreateStream",
			Handler:       _Broadcast_CreateStream_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "chat/chat.proto",
}

const (
	Rooms_CreateRoom_FullMethodName   = "/chat.Rooms/CreateRoom"
	Rooms_GetUserRooms_FullMethodName = "/chat.Rooms/GetUserRooms"
	Rooms_DeleteRoom_FullMethodName   = "/chat.Rooms/DeleteRoom"
)

// RoomsClient is the client API for Rooms service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type RoomsClient interface {
	CreateRoom(ctx context.Context, in *CreateRoomRequest, opts ...grpc.CallOption) (*Room, error)
	GetUserRooms(ctx context.Context, in *ListRoomRequest, opts ...grpc.CallOption) (*ListRoomResponse, error)
	DeleteRoom(ctx context.Context, in *DeleteRoomRequest, opts ...grpc.CallOption) (*EmptyResponse, error)
}

type roomsClient struct {
	cc grpc.ClientConnInterface
}

func NewRoomsClient(cc grpc.ClientConnInterface) RoomsClient {
	return &roomsClient{cc}
}

func (c *roomsClient) CreateRoom(ctx context.Context, in *CreateRoomRequest, opts ...grpc.CallOption) (*Room, error) {
	out := new(Room)
	err := c.cc.Invoke(ctx, Rooms_CreateRoom_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *roomsClient) GetUserRooms(ctx context.Context, in *ListRoomRequest, opts ...grpc.CallOption) (*ListRoomResponse, error) {
	out := new(ListRoomResponse)
	err := c.cc.Invoke(ctx, Rooms_GetUserRooms_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *roomsClient) DeleteRoom(ctx context.Context, in *DeleteRoomRequest, opts ...grpc.CallOption) (*EmptyResponse, error) {
	out := new(EmptyResponse)
	err := c.cc.Invoke(ctx, Rooms_DeleteRoom_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// RoomsServer is the server API for Rooms service.
// All implementations must embed UnimplementedRoomsServer
// for forward compatibility
type RoomsServer interface {
	CreateRoom(context.Context, *CreateRoomRequest) (*Room, error)
	GetUserRooms(context.Context, *ListRoomRequest) (*ListRoomResponse, error)
	DeleteRoom(context.Context, *DeleteRoomRequest) (*EmptyResponse, error)
	mustEmbedUnimplementedRoomsServer()
}

// UnimplementedRoomsServer must be embedded to have forward compatible implementations.
type UnimplementedRoomsServer struct {
}

func (UnimplementedRoomsServer) CreateRoom(context.Context, *CreateRoomRequest) (*Room, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateRoom not implemented")
}
func (UnimplementedRoomsServer) GetUserRooms(context.Context, *ListRoomRequest) (*ListRoomResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetUserRooms not implemented")
}
func (UnimplementedRoomsServer) DeleteRoom(context.Context, *DeleteRoomRequest) (*EmptyResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteRoom not implemented")
}
func (UnimplementedRoomsServer) mustEmbedUnimplementedRoomsServer() {}

// UnsafeRoomsServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to RoomsServer will
// result in compilation errors.
type UnsafeRoomsServer interface {
	mustEmbedUnimplementedRoomsServer()
}

func RegisterRoomsServer(s grpc.ServiceRegistrar, srv RoomsServer) {
	s.RegisterService(&Rooms_ServiceDesc, srv)
}

func _Rooms_CreateRoom_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateRoomRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RoomsServer).CreateRoom(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Rooms_CreateRoom_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RoomsServer).CreateRoom(ctx, req.(*CreateRoomRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Rooms_GetUserRooms_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListRoomRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RoomsServer).GetUserRooms(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Rooms_GetUserRooms_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RoomsServer).GetUserRooms(ctx, req.(*ListRoomRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Rooms_DeleteRoom_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteRoomRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RoomsServer).DeleteRoom(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Rooms_DeleteRoom_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RoomsServer).DeleteRoom(ctx, req.(*DeleteRoomRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Rooms_ServiceDesc is the grpc.ServiceDesc for Rooms service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Rooms_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "chat.Rooms",
	HandlerType: (*RoomsServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateRoom",
			Handler:    _Rooms_CreateRoom_Handler,
		},
		{
			MethodName: "GetUserRooms",
			Handler:    _Rooms_GetUserRooms_Handler,
		},
		{
			MethodName: "DeleteRoom",
			Handler:    _Rooms_DeleteRoom_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "chat/chat.proto",
}

const (
	Messages_GetMessageByID_FullMethodName = "/chat.Messages/GetMessageByID"
)

// MessagesClient is the client API for Messages service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type MessagesClient interface {
	GetMessageByID(ctx context.Context, in *GetMessageByIDRequest, opts ...grpc.CallOption) (*Message, error)
}

type messagesClient struct {
	cc grpc.ClientConnInterface
}

func NewMessagesClient(cc grpc.ClientConnInterface) MessagesClient {
	return &messagesClient{cc}
}

func (c *messagesClient) GetMessageByID(ctx context.Context, in *GetMessageByIDRequest, opts ...grpc.CallOption) (*Message, error) {
	out := new(Message)
	err := c.cc.Invoke(ctx, Messages_GetMessageByID_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// MessagesServer is the server API for Messages service.
// All implementations must embed UnimplementedMessagesServer
// for forward compatibility
type MessagesServer interface {
	GetMessageByID(context.Context, *GetMessageByIDRequest) (*Message, error)
	mustEmbedUnimplementedMessagesServer()
}

// UnimplementedMessagesServer must be embedded to have forward compatible implementations.
type UnimplementedMessagesServer struct {
}

func (UnimplementedMessagesServer) GetMessageByID(context.Context, *GetMessageByIDRequest) (*Message, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetMessageByID not implemented")
}
func (UnimplementedMessagesServer) mustEmbedUnimplementedMessagesServer() {}

// UnsafeMessagesServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to MessagesServer will
// result in compilation errors.
type UnsafeMessagesServer interface {
	mustEmbedUnimplementedMessagesServer()
}

func RegisterMessagesServer(s grpc.ServiceRegistrar, srv MessagesServer) {
	s.RegisterService(&Messages_ServiceDesc, srv)
}

func _Messages_GetMessageByID_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetMessageByIDRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MessagesServer).GetMessageByID(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Messages_GetMessageByID_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MessagesServer).GetMessageByID(ctx, req.(*GetMessageByIDRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Messages_ServiceDesc is the grpc.ServiceDesc for Messages service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Messages_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "chat.Messages",
	HandlerType: (*MessagesServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetMessageByID",
			Handler:    _Messages_GetMessageByID_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "chat/chat.proto",
}
