package main

import (
	"graph-replicator/graph"
	"graph-replicator/server"
	"log"
)

func main() {
	log.Println("Creating new graph")
	g := graph.NewDirectedGraph()
	s := server.New(g)
	log.Println("Starting server on port 3000")
	s.Start(3000)
}
