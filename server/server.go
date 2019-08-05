package server

import (
	"bytes"
	"context"
	"encoding/gob"
	"fmt"
	"graph-replicator/api"
	"graph-replicator/graph"
	"log"
	"net"

	"google.golang.org/grpc"

	"github.com/golang/protobuf/ptypes/empty"
)

type server struct {
	graph graph.Graph
}

func New(g graph.Graph) *server {
	return &server{
		graph: g,
	}
}

func (s *server) ShowGraph(ctx context.Context, in *empty.Empty) (*api.Response, error) {
	return &api.Response{
		Graph: s.graph.String(),
	}, nil
}

func (s *server) SyncGraph(ctx context.Context, in *api.SyncRequest) (*api.Response, error) {
	b := bytes.NewBuffer(in.Graph)
	dec := gob.NewDecoder(&b)
	newG := graph.NewDirectedGraph()
	err = dec.Decode(&newG)
	if err != nil {
		log.Fatalf("Failed to decode: %s", err.Error())
	}

	return &api.Response{
		Graph: s.graph.String(),
	}, nil
}

func (s *server) Start(port int) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("Failed to start server: %s", err.Error())
	}

	grpcServer := grpc.NewServer()
	api.RegisterGraphReplicatorServer(grpcServer, s)
	grpcServer.Serve(lis)
}
