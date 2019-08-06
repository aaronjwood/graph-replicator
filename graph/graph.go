package graph

import (
	"bytes"
	"encoding"
	"fmt"
	"io"
	"log"
	"strings"

	"github.com/google/uuid"
)

const (
	remove = iota
	add
)

type Graph interface {
	AddEdge(src, dest *Vertex)
	RemoveEdge(src, dest *Vertex)
	fmt.Stringer
	encoding.BinaryMarshaler
	encoding.BinaryUnmarshaler
}

type DirectedGraph struct {
	adjacency  map[Vertex]vertexSet
	addDiff    map[Vertex]vertexSet
	removeDiff []*Vertex
}

func NewDirectedGraph() *DirectedGraph {
	return &DirectedGraph{
		adjacency:  make(map[Vertex]vertexSet),
		addDiff:    make(map[Vertex]vertexSet),
		removeDiff: make([]*Vertex, 0),
	}
}

func (g *DirectedGraph) AddEdge(src, dest *Vertex) {
	if g.adjacency[*src] == nil {
		g.adjacency[*src] = make(vertexSet)
	}

	g.adjacency[*src].Add(*dest)
	g.addDiff[*src] = g.adjacency[*src]
}

func (g *DirectedGraph) RemoveEdge(src, dest *Vertex) {
	g.adjacency[*src].Remove(*dest)
	g.addDiff[*src] = g.adjacency[*src]
	if len(g.adjacency[*src]) == 0 {
		delete(g.adjacency, *src)
		g.removeDiff = append(g.removeDiff, src)
	}
}

func (g *DirectedGraph) String() string {
	var builder strings.Builder
	for src, dest := range g.adjacency {
		s := fmt.Sprintf("Source vertex: %d\n", src.data)
		builder.WriteString(s)
		builder.WriteString("\tDestination verticies: ")
		for vertex := range dest {
			s := fmt.Sprintf("%d    ", vertex.data)
			builder.WriteString(s)
		}
		builder.WriteString("\n\n")
	}

	return builder.String()
}

func (g *DirectedGraph) MarshalBinary() ([]byte, error) {
	var b bytes.Buffer
	delim := " "
	for vertex, vertices := range g.addDiff {
		fmt.Fprint(&b, add, delim, vertex.data, delim, vertex.id.String(), delim, len(vertices), delim)
		for v := range vertices {
			fmt.Fprint(&b, v.data, delim, v.id.String(), delim)
		}

		fmt.Fprintln(&b)
	}

	for _, vertex := range g.removeDiff {
		fmt.Fprint(&b, remove, delim, vertex.data, delim, vertex.id.String(), delim, 0)
		fmt.Fprintln(&b)
	}

	g.addDiff = make(map[Vertex]vertexSet)
	g.removeDiff = make([]*Vertex, 0)
	return b.Bytes(), nil
}

func (g *DirectedGraph) UnmarshalBinary(data []byte) error {
	defer func() {
		g.addDiff = make(map[Vertex]vertexSet)
		g.removeDiff = make([]*Vertex, 0)
	}()

	b := bytes.NewBuffer(data)
	var op int
	var srcVData int
	var srcVID string
	var destLen int
	for {
		_, err := fmt.Fscan(b, &op, &srcVData, &srcVID, &destLen)
		if err == io.EOF {
			return nil
		}

		if err != nil {
			return err
		}

		log.Printf("Decoding vertex %s", srcVID)
		srcV := NewVertex(srcVData)
		srcV.id = uuid.MustParse(srcVID)

		var destVData int
		var destVID string
		destVertices := make(vertexSet)
		if op == remove {
			log.Printf("Deleting vertex %s", srcVID)
			delete(g.adjacency, *srcV)
			continue
		}

		log.Printf("Decoding %d destination vertices", destLen)
		for i := 0; i < destLen; i++ {
			_, err = fmt.Fscan(b, &destVData, &destVID)
			if err == io.EOF {
				return nil
			}

			if err != nil {
				return err
			}

			destV := NewVertex(destVData)
			destV.id = uuid.MustParse(destVID)
			destVertices.Add(*destV)
		}

		g.adjacency[*srcV] = destVertices
	}
}
