.PHONY: protos build

dep:
	@go get -u github.com/google/uuid
	@go get -u github.com/golang/protobuf/protoc-gen-go
	@go get -u google.golang.org/grpc

protos:
	@protoc -I protos --go_out=plugins=grpc:api/ protos/server.proto
	@protoc -I protos --go_out=plugins=grpc:api/ protos/client.proto

build: dep
	@go build -o bin/client cmd/client/main.go
	@go build -o bin/server cmd/server/main.go
