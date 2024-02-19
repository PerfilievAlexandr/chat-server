package chat

import (
	"context"
	"github.com/PerfilievAlexandr/chat-server/internal/api/grpc/chat/mapper"
	authClient "github.com/PerfilievAlexandr/chat-server/internal/integration/auth"
	"github.com/PerfilievAlexandr/chat-server/internal/service"
	proto "github.com/PerfilievAlexandr/chat-server/pkg/chat_v1"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Server struct {
	proto.UnimplementedChatV1Server
	chatService      service.ChatService
	checkRoleService service.CheckRoleService
	authClient       authClient.AuthServiceClient
}

func NewServer(
	_ context.Context,
	chatService service.ChatService,
	checkRoleService service.CheckRoleService,
	authClient authClient.AuthServiceClient,
) *Server {
	return &Server{chatService: chatService, checkRoleService: checkRoleService, authClient: authClient}
}

func (s *Server) SendMessage(ctx context.Context, req *proto.SendMessageRequest) (*emptypb.Empty, error) {
	err := s.chatService.SendMessage(ctx, mapper.MapToSendMessageRequest(req))
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (s *Server) CreateChat(ctx context.Context, req *proto.CreateChatRequest) (*proto.CreateChatResponse, error) {
	userInfo, err := s.authClient.Check(ctx)
	if err != nil {
		return nil, err
	}
	err = s.checkRoleService.CheckAdmin(ctx, userInfo.Role)
	if err != nil {
		return nil, err
	}

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
