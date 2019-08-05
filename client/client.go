package client

import (
	"bytes"
	"context"
	"encoding/gob"
	"graph-replicator/api"
	"graph-replicator/graph"
	"log"

	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc"
)

type client struct {
	graph  graph.Graph
	conn   *grpc.ClientConn
	client api.GraphReplicatorClient
}

func New(g graph.Graph) *client {
	return &client{
		graph: g,
	}
}

func (c *client) LocalGraph() string {
	log.Println("Getting local client graph")
	return c.graph.String()
}

func (c *client) RemoteGraph() (string, error) {
	log.Println("Getting graph from server")
	res, err := c.client.ShowGraph(context.Background(), &empty.Empty{})
	if err != nil {
		log.Fatalf("Failed to show remote graph: %s", err.Error())
	}

	return res.Graph, err
}

func (c *client) SyncGraph() error {
	var b bytes.Buffer
	enc := gob.NewEncoder(&b)
	err := enc.Encode(c.graph)
	if err != nil {
		log.Fatalf("Failed to encode: %s", err.Error())
	}

	_, err = c.client.SyncGraph(context.Background(), &api.SyncRequest{
		Graph: b.Bytes(),
	})
	return err
}

func (c *client) Connect() error {
	log.Println("Connecting to remote graph")
	conn, err := grpc.Dial("127.0.0.1:3000", grpc.WithInsecure())
	if err != nil {
		return err
	}

	c.conn = conn
	c.client = api.NewGraphReplicatorClient(conn)
	return nil
}

func (c *client) Disconnect() {
	c.conn.Close()
}
