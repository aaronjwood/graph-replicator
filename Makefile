.PHONY: protos build

protos:
	@protoc -I protos --go_out=plugins=grpc:api/ protos/server.proto
	@protoc -I protos --go_out=plugins=grpc:api/ protos/client.proto

build:
	@go build -o bin/client cmd/client/main.go
	@go build -o bin/server cmd/server/main.go
