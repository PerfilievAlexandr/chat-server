package interceptor

import (
	"context"
	"errors"
	"fmt"
	authProto "github.com/PerfilievAlexandr/auth/pkg/access_v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"log"
)

const servicePort = 8080

func AccessInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, errors.New("metadata is not provided")
	}
	outgoingContext := metadata.NewOutgoingContext(ctx, md)

	conn, err := grpc.Dial(
		fmt.Sprintf(":%d", servicePort),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("failed to dial GRPC client: %v", err)
	}

	cl := authProto.NewAccessV1Client(conn)

	_, err = cl.Check(outgoingContext, &authProto.CheckRequest{
		EndpointAddress: "",
	})
	if err != nil {
		log.Fatal(err.Error())
	}
	fmt.Printf("%s", outgoingContext)

	return handler(ctx, req)
}
