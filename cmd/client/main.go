package main

import (
	"bytes"
	"context"
	"encoding/gob"
	"fmt"
	"graph-replicator/api"
	"graph-replicator/graph"
	"log"
	"math/rand"
	"time"

	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc"
)

const (
	fill = 10
)

var data = make([]*graph.Vertex, fill)

// Setup a bunch of vertices that we can reuse to show off the differential syncing.
func initData() {
	for i := 0; i < fill; i++ {
		data[i] = graph.NewVertex(rand.Intn(100))
	}
}

// Add edges to the graph from the pool of vertices.
// Randomly reuse vertices so that we can exercise the graph's ability to sync only diffs and not the entire graph.
func fillGraph(g graph.Graph) {
	for i := 0; i < fill; i++ {
		g.AddEdge(data[rand.Intn(fill)], data[rand.Intn(fill)])
	}
}

func main() {
	fmt.Println("Creating new graph")
	g := graph.NewDirectedGraph()
	initData()
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

	syncs := 20
	start := time.Now()
	for i := 0; i < syncs; i++ {
		fillGraph(g)
		fmt.Println("Syncing local graph to remote")
		var b bytes.Buffer
		enc := gob.NewEncoder(&b)
		err = enc.Encode(g)
		if err != nil {
			log.Fatalf("Failed to encode: %s", err.Error())
		}

		start := time.Now()
		_, err = client.SyncGraph(context.Background(), &api.SyncRequest{
			Graph: b.Bytes(),
		})
		if err != nil {
			log.Fatalf("Failed to sync graph: %s", err.Error())
		}

		fmt.Printf("Took %s to sync changes\n", time.Since(start))
	}

	fmt.Printf("Took %s to sync all changes\n", time.Since(start))
	res, err = client.ShowGraph(context.Background(), &empty.Empty{})
	if err != nil {
		log.Fatalf("Failed to show remote graph: %s", err.Error())
	}

	fmt.Println("Remote graph:")
	fmt.Println(res.Graph)
}
