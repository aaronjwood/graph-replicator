package main

import (
	"fmt"
	"graph-replicator/client"
	"graph-replicator/graph"
	"log"
	"math/rand"
	"time"
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
	initData()
	g := graph.NewDirectedGraph()
	client := client.New(g)
	err := client.Connect()
	defer client.Disconnect()
	if err != nil {
		log.Fatalf("Failed to connect: %s", err.Error())
	}

	remoteGraph, err := client.RemoteGraph()
	if err != nil {
		log.Fatalf("Failed to get remote graph: %s", err.Error())
	}

	fmt.Println("Local graph:")
	fmt.Println(client.LocalGraph())

	fmt.Println("Remote graph:")
	fmt.Println(remoteGraph)
	syncs := 20
	start := time.Now()
	for i := 0; i < syncs; i++ {
		fillGraph(g)
		fmt.Println("Syncing local graph to remote")
		start := time.Now()
		err := client.SyncGraph()
		if err != nil {
			log.Fatalf("Failed to sync graph: %s", err.Error())
		}

		fmt.Printf("Took %s to sync changes\n", time.Since(start))
	}

	fmt.Println()
	fmt.Printf("Took %s to sync all changes\n", time.Since(start))
	remoteGraph, err = client.RemoteGraph()
	if err != nil {
		log.Fatalf("Failed to get remote graph: %s", err.Error())
	}

	fmt.Println("Local graph:")
	fmt.Println(client.LocalGraph())

	fmt.Println("Remote graph:")
	fmt.Println(remoteGraph)
}
