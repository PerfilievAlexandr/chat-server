package main

import (
	desc "chat-server/pkg/chat_v1"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/emptypb"
	"log"
	"net"
)

const grpcPort = 50051

type server struct {
	desc.UnimplementedChatV1Server
}

func (s *server) Create(ctx context.Context, req *desc.CreteRequest) (*desc.CreateResponse, error) {
	log.Printf("Request: %d", req)

	return &desc.CreateResponse{}, nil
}

func (s *server) Delete(ctx context.Context, req *desc.DeleteRequest) (*emptypb.Empty, error) {
	log.Printf("Request: %d", req)

	return &emptypb.Empty{}, nil
}

func (s *server) SendMessage(ctx context.Context, req *desc.SendMessageRequest) (*emptypb.Empty, error) {
	log.Printf("Request: %d", req)

	return &emptypb.Empty{}, nil
}

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	reflection.Register(s)
	desc.RegisterChatV1Server(s, &server{})

	log.Printf("server listening at %v", lis.Addr())

	if err = s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}