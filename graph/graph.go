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
	adjacency  map[Vertex][]*Vertex
	addDiff    map[Vertex][]*Vertex
	removeDiff []*Vertex
}

func NewDirectedGraph() *DirectedGraph {
	return &DirectedGraph{
		adjacency:  make(map[Vertex][]*Vertex),
		addDiff:    make(map[Vertex][]*Vertex),
		removeDiff: make([]*Vertex, 0),
	}
}

func (g *DirectedGraph) AddEdge(src, dest *Vertex) {
	for _, destV := range g.adjacency[*src] {
		if destV.id == dest.id {
			return
		}
	}

	g.adjacency[*src] = append(g.adjacency[*src], dest)
	g.addDiff[*src] = append(g.addDiff[*src], dest)
}

func (g *DirectedGraph) RemoveEdge(src, dest *Vertex) {
	for idx, destVertex := range g.adjacency[*src] {
		if destVertex.id == dest.id {
			g.adjacency[*src] = append(g.adjacency[*src][:idx], g.adjacency[*src][idx+1:]...)
		}
	}

	g.addDiff[*src] = g.adjacency[*src]
	if len(g.adjacency[*src]) == 0 {
		delete(g.adjacency, *src)
		g.removeDiff = append(g.removeDiff, src)
	}
}

func (g *DirectedGraph) String() string {
	var builder strings.Builder
	for src, destList := range g.adjacency {
		s := fmt.Sprintf("Source vertex: %d\n", src.data)
		builder.WriteString(s)
		builder.WriteString("\tDestination verticies: ")
		for _, dest := range destList {
			s := fmt.Sprintf("%d    ", dest.data)
			builder.WriteString(s)
		}
		builder.WriteString("\n\n")
	}

	return builder.String()
}

func (g *DirectedGraph) MarshalBinary() ([]byte, error) {
	var b bytes.Buffer
	delim := " "
	log.Printf("Encoding %d source vertices to add", len(g.addDiff))
	for vertex, vertices := range g.addDiff {
		log.Printf("Encoding %d destination vertices", len(vertices))
		fmt.Fprint(&b, add, delim, vertex.data, delim, vertex.id.String(), delim, len(vertices), delim)
		for _, v := range vertices {
			fmt.Fprint(&b, v.data, delim, v.id.String(), delim)
		}

		fmt.Fprintln(&b)
	}

	log.Printf("Encoding %d source vertices to remove", len(g.removeDiff))
	for _, vertex := range g.removeDiff {
		log.Printf("Encoding %d destination vertices", len(g.removeDiff))
		fmt.Fprint(&b, remove, delim, vertex.data, delim, vertex.id.String(), delim, 0)
		fmt.Fprintln(&b)
	}

	g.addDiff = make(map[Vertex][]*Vertex)
	g.removeDiff = make([]*Vertex, 0)
	return b.Bytes(), nil
}

func (g *DirectedGraph) UnmarshalBinary(data []byte) error {
	defer func() {
		g.addDiff = make(map[Vertex][]*Vertex)
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
		var destVertices []*Vertex
		if op == remove {
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
			destVertices = append(destVertices, destV)
		}

		g.adjacency[*srcV] = destVertices
	}
}
