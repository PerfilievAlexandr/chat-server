package main

import (
	"context"
	"fmt"
	proto "github.com/PerfilievAlexandr/chat-server/pkg/chat_v1"
	"io"
	"log"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const (
	address = "localhost:8008"
)

// TODO времянка
func main() {
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect to server: %v", err)
	}
	defer conn.Close()

	ctx := context.Background()
	client := proto.NewChatV1Client(conn)
	user := "oleg"

	// Создаем новый чат на сервере
	//chatID, err := createChat(ctx, client, user)
	//if err != nil {
	//	log.Fatalf("failed to create chat: %v", err)
	//}
	//
	//log.Printf(fmt.Sprintf("%s: %s\n", "Chat created", chatID))

	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		defer wg.Done()

		err = connectChat(ctx, client, "9785b424-d14d-11ee-92ea-de0c48c12eae", user, 5*time.Second)
		if err != nil {
			log.Fatalf("failed to connect chat: %v", err)
		}
	}()

	wg.Wait()
}

func connectChat(ctx context.Context, client proto.ChatV1Client, chatID string, username string, period time.Duration) error {
	stream, err := client.ConnectChat(ctx, &proto.ConnectChatRequest{
		ChatId:   chatID,
		Username: username,
	})
	if err != nil {
		return err
	}

	go func() {
		for {
			message, errRecv := stream.Recv()
			if errRecv == io.EOF {
				return
			}
			if errRecv != nil {
				log.Println("failed to receive message from stream: ", errRecv)
				return
			}

			log.Printf("[%v] - [from: %s] - [to: %s]: %s\n",
				message.Timestamp.String(),
				message.GetFrom(),
				username,
				message.GetText(),
			)
		}
	}()

	for {
		time.Sleep(period)
		text := fmt.Sprintf("Message from %s, hello!", username)

		_, err = client.SendMessage(ctx, &proto.SendMessageRequest{
			ChatId: chatID,
			Message: &proto.Message{
				From:      username,
				Text:      text,
				Timestamp: timestamppb.Now(),
			},
		})
		if err != nil {
			log.Println("failed to send message: ", err)
			return err
		}
	}
}

func createChat(ctx context.Context, client proto.ChatV1Client, user string) (string, error) {
	res, err := client.CreateChat(ctx, &proto.CreateChatRequest{
		Username: user,
	})
	if err != nil {
		return "", err
	}

	return res.GetChatId(), nil
}
