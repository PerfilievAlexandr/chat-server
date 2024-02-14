package chat

import (
	"context"
	"github.com/PerfilievAlexandr/chat-server/internal/api/grpc/chat/mapper"
	"github.com/PerfilievAlexandr/chat-server/internal/service"
	proto "github.com/PerfilievAlexandr/chat-server/pkg/chat_v1"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Server struct {
	proto.UnimplementedChatV1Server
	chatService service.ChatService
}

func NewServer(_ context.Context, chatService service.ChatService) *Server {
	return &Server{chatService: chatService}
}

func (s *Server) SendMessage(ctx context.Context, req *proto.SendMessageRequest) (*emptypb.Empty, error) {
	err := s.chatService.SendMessage(ctx, mapper.MapToSendMessageRequest(req))
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (s *Server) CreateChat(ctx context.Context, req *proto.CreateChatRequest) (*proto.CreateChatResponse, error) {
	chatId, err := s.chatService.CreateChat(ctx, mapper.MapToCreateChatRequest(req))
	if err != nil {
		return nil, err
	}

	return &proto.CreateChatResponse{
		ChatId: chatId.String(),
	}, nil
}

func (s *Server) ConnectChat(req *proto.ConnectChatRequest, stream proto.ChatV1_ConnectChatServer) error {
	return s.chatService.ConnectChat(mapper.MapToConnectChatRequest(req), stream)
}
