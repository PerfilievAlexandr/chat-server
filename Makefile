include .env
LOCAL_BIN:=$(CURDIR)/bin
PG_DSN:="host=$(DB_HOST) port=$(DB_PORT) user=$(DB_USER) password=$(DB_PASSWORD) dbname=$(DB_NAME) sslmode=$(DB_SSL_MODE)"

install-deps:
	GOBIN=$(LOCAL_BIN) go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28.1
	GOBIN=$(LOCAL_BIN) go install -mod=mod google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
	GOBIN=$(LOCAL_BIN) go install github.com/pressly/goose/v3/cmd/goose@v3.14.0

get-deps:
	go get -u google.golang.org/protobuf/cmd/protoc-gen-go
	go get -u google.golang.org/grpc/cmd/protoc-gen-go-grpc

generate:
	make generate-chat-api

generate-chat-api:
	mkdir -p pkg/chat_v1
	protoc --proto_path api/chat_v1 \
	--go_out=pkg/chat_v1 --go_opt=paths=source_relative \
	--plugin=protoc-gen-go=bin/protoc-gen-go \
	--go-grpc_out=pkg/chat_v1 --go-grpc_opt=paths=source_relative \
	--plugin=protoc-gen-go-grpc=bin/protoc-gen-go-grpc \
	api/chat_v1/chat.proto

local-migration-status:
	$(LOCAL_BIN)/goose -dir "./migrations" postgres $(PG_DSN) status -v

local-migration-up:
	$(LOCAL_BIN)/goose -dir "./migrations" postgres $(PG_DSN) up -v

local-migration-down:
	$(LOCAL_BIN)/goose -dir "./migrations" postgres $(PG_DSN) down -v
