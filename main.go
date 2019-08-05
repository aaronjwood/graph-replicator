package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"graph-replicator/graph"
	"graph-replicator/server"
	"log"
)

func main() {
	s := server.New()
	fmt.Println("Creating new graph")
	g := graph.NewDirectedGraph()
	fmt.Println("Adding edge to graph")
	v1 := graph.NewVertex(1)
	v2 := graph.NewVertex(2)
	v3 := graph.NewVertex(3)
	g.AddEdge(v1, v2)
	g.AddEdge(v1, v3)
	fmt.Println("Adding vertex to graph")
	v4 := graph.NewVertex(4)
	g.AddVertex(v4)
	g.RemoveEdge(v1, v2)
	var b bytes.Buffer
	enc := gob.NewEncoder(&b)
	err := enc.Encode(g)
	if err != nil {
		log.Fatalf("Failed to encode: %s", err.Error())
	}

	dec := gob.NewDecoder(&b)
	newG := graph.NewDirectedGraph()
	err = dec.Decode(&newG)
	if err != nil {
		log.Fatalf("Failed to decode: %s", err.Error())
	}

	fmt.Println("Displaying graph\n")
	fmt.Println(g.String())

	fmt.Println("Displaying new graph\n")
	fmt.Println(newG.String())
}
