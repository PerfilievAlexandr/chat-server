package chat

import (
	"context"
	"github.com/PerfilievAlexandr/chat-server/internal/api/grpc/chat/dto"
	"github.com/PerfilievAlexandr/chat-server/internal/api/grpc/chat/mapper"
	"github.com/PerfilievAlexandr/chat-server/internal/domain"
	"github.com/PerfilievAlexandr/chat-server/internal/repository"
	"github.com/PerfilievAlexandr/chat-server/internal/service"
	proto "github.com/PerfilievAlexandr/chat-server/pkg/chat_v1"
	"github.com/PerfilievAlexandr/platform_common/pkg/db"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"sync"
)

type chat struct {
	streams  map[string]proto.ChatV1_ConnectChatServer
	mxStream sync.RWMutex
}

type chatService struct {
	chats  map[string]*chat
	mxChat sync.RWMutex

	channels  map[string]chan domain.Message
	mxChannel sync.RWMutex

	messageRepository repository.MessageRepository
	chatRepository    repository.ChatRepository
	tx                db.TxManager
}

func NewChatService(
	_ context.Context,
	messageRepo repository.MessageRepository,
	chatRepo repository.ChatRepository,
	tx db.TxManager,
) service.ChatService {
	return &chatService{
		chats:             make(map[string]*chat),
		channels:          make(map[string]chan domain.Message),
		messageRepository: messageRepo,
		chatRepository:    chatRepo,
		tx:                tx,
	}
}

func (c *chatService) CreateChat(ctx context.Context, req dto.CreateChatRequest) (uuid.UUID, error) {
	chatId, err := c.chatRepository.SaveChat(ctx, req)
	if err != nil {
		return uuid.UUID{}, err
	}

	c.channels[chatId.String()] = make(chan domain.Message, 100)

	return chatId, nil
}

func (c *chatService) ConnectChat(req dto.ConnectChatRequest, stream proto.ChatV1_ConnectChatServer) error {
	chanel, err := c.createChannelIfNotExist(stream.Context(), req.ChatId)
	if err != nil {
		return err
	}
	err = c.fillChanelFromDb(stream.Context(), req.ChatId, chanel)
	if err != nil {
		return err
	}
	c.setStreamToChatByUsername(stream.Context(), req, stream)

	for {
		select {
		case msg, okCh := <-chanel:
			if !okCh {
				return nil
			}
			err = c.sendToChatMembers(stream.Context(), req.ChatId, msg)
			if err != nil {
				return err
			}
		case <-stream.Context().Done():
			c.chats[req.ChatId].mxStream.Lock()
			delete(c.chats[req.ChatId].streams, req.Username)
			c.chats[req.ChatId].mxStream.Unlock()
			return nil
		}
	}
}

func (c *chatService) createChannelIfNotExist(ctx context.Context, chatId string) (chan domain.Message, error) {
	isChatExists, err := c.chatRepository.IsExists(ctx, uuid.MustParse(chatId))
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "chat not created")
	}
	c.mxChannel.RLock()
	chanel, ok := c.channels[chatId]
	if !ok && isChatExists {
		c.channels[chatId] = make(chan domain.Message, 100)
		chanel = c.channels[chatId]
	} else if !ok && !isChatExists {
		return nil, status.Errorf(codes.NotFound, "chat not found")
	}
	c.mxChannel.RUnlock()

	return chanel, nil
}

func (c *chatService) fillChanelFromDb(ctx context.Context, chatId string, chanel chan domain.Message) error {
	messages, err := c.messageRepository.GetMessagesByChatId(ctx, uuid.MustParse(chatId))
	if err != nil {
		return status.Errorf(codes.Internal, "load messages error")
	}
	for _, message := range messages {
		chanel <- message
	}

	return nil
}

func (c *chatService) setStreamToChatByUsername(_ context.Context, req dto.ConnectChatRequest, stream proto.ChatV1_ConnectChatServer) {
	c.mxChat.Lock()
	_, chatOk := c.chats[req.ChatId]
	if !chatOk {
		c.chats[req.ChatId] = &chat{
			streams: make(map[string]proto.ChatV1_ConnectChatServer),
		}
	}
	c.mxChat.Unlock()

	c.chats[req.ChatId].mxStream.Lock()
	c.chats[req.ChatId].streams[req.Username] = stream
	c.chats[req.ChatId].mxStream.Unlock()
}

func (c *chatService) sendToChatMembers(_ context.Context, chatID string, message domain.Message) error {
	for key, st := range c.chats[chatID].streams {
		if key != message.From {
			if err := st.Send(mapper.MapToProtoMessage(message)); err != nil {
				return err
			}
		}
	}

	return nil
}

func (c *chatService) SendMessage(ctx context.Context, message dto.SendMessageRequest) error {
	c.mxChannel.Lock()
	chanel, ok := c.channels[message.ChatId]
	c.mxChannel.Unlock()
	if !ok {
		return status.Errorf(codes.NotFound, "chat not found")
	}

	saveMessage, err := c.messageRepository.SaveMessage(ctx, message)
	if err != nil {
		return status.Errorf(codes.Internal, "error save message")
	}

	chanel <- saveMessage
	return nil
}
