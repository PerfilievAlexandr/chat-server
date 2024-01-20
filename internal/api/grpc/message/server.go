package message

import (
	"chat-server/internal/api/grpc/message/mapper"
	"chat-server/internal/service"
	proto "chat-server/pkg/chat_v1"
	"context"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Server struct {
	proto.UnimplementedChatV1Server
	messageService service.MessageService
}

func NewServer(ctx context.Context, messageService service.MessageService) *Server {
	return &Server{messageService: messageService}
}

func (s *Server) Create(ctx context.Context, req *proto.CreteRequest) (*proto.CreateResponse, error) {
	id, err := s.messageService.Create(ctx, req.Usernames)
	if err != nil {
		return nil, err
	}

	return &proto.CreateResponse{
		Id: id,
	}, nil
}

func (s *Server) Delete(ctx context.Context, req *proto.DeleteRequest) (*emptypb.Empty, error) {
	return s.messageService.Delete(ctx, req.Id)
}

func (s *Server) SendMessage(ctx context.Context, req *proto.SendMessageRequest) (*emptypb.Empty, error) {
	return s.messageService.SendMessage(ctx, mapper.MapToSendMessageRequest(req))
}
