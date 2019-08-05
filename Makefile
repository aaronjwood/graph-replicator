.PHONY: protos run build

protos:
	@protoc -I protos --go_out=plugins=grpc:api/ protos/server.proto
	@protoc -I protos --go_out=plugins=grpc:api/ protos/client.proto

run:
	@go run .

build:
	@go build .
