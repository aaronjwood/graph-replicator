package main

import (
	"bytes"
	"context"
	"encoding/gob"
	"fmt"
	"graph-replicator/api"
	"graph-replicator/graph"
	"log"

	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc"
)

func main() {
	fmt.Println("Creating new graph")
	g := graph.NewDirectedGraph()

	fmt.Println("Adding some edges to the graph")
	v1 := graph.NewVertex(1)
	v2 := graph.NewVertex(2)
	v3 := graph.NewVertex(3)
	v4 := graph.NewVertex(4)
	v5 := graph.NewVertex(5)
	g.AddEdge(v1, v2)
	g.AddEdge(v1, v3)
	g.AddEdge(v4, v5)

	fmt.Println("Local graph:")
	fmt.Println(g.String())

	fmt.Println("Adding vertex with no connections to the graph")
	v6 := graph.NewVertex(6)
	g.AddVertex(v6)

	fmt.Println("Local graph:")
	fmt.Println(g.String())

	fmt.Println("Removing an edge from the graph")
	g.RemoveEdge(v1, v2)

	fmt.Println("Local graph:")
	fmt.Println(g.String())

	fmt.Println("Connecting to remote graph")
	conn, err := grpc.Dial("127.0.0.1:3000", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect: %s", err.Error())
	}

	defer conn.Close()
	client := api.NewGraphReplicatorClient(conn)
	res, err := client.ShowGraph(context.Background(), &empty.Empty{})
	if err != nil {
		log.Fatalf("Failed to show remote graph: %s", err.Error())
	}

	fmt.Println("Remote graph:")
	fmt.Println(res.Graph)

	var b bytes.Buffer
	enc := gob.NewEncoder(&b)
	err = enc.Encode(g)
	if err != nil {
		log.Fatalf("Failed to encode: %s", err.Error())
	}

	fmt.Println("Syncing local graph to remote")
	_, err = client.SyncGraph(context.Background(), &api.SyncRequest{
		Graph: b.Bytes(),
	})
	if err != nil {
		log.Fatalf("Failed to sync graph: %s", err.Error())
	}

	res, err = client.ShowGraph(context.Background(), &empty.Empty{})
	if err != nil {
		log.Fatalf("Failed to show remote graph: %s", err.Error())
	}

	fmt.Println("Remote graph:")
	fmt.Println(res.Graph)
}
