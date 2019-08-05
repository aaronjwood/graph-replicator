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

type Graph interface {
	AddEdge(src, dest *Vertex)
	RemoveEdge(src, dest *Vertex)
	fmt.Stringer
	encoding.BinaryMarshaler
	encoding.BinaryUnmarshaler
}

type DirectedGraph struct {
	adjacency map[*Vertex][]*Vertex
	diff      map[*Vertex][]*Vertex // Trade off space for caching only what's changed for when we encode.
}

func NewDirectedGraph() *DirectedGraph {
	return &DirectedGraph{
		adjacency: make(map[*Vertex][]*Vertex),
		diff:      make(map[*Vertex][]*Vertex),
	}
}

func (g *DirectedGraph) AddEdge(src, dest *Vertex) {
	for _, destV := range g.adjacency[src] {
		if destV.id == dest.id {
			return
		}
	}

	g.adjacency[src] = append(g.adjacency[src], dest)
	g.diff[src] = append(g.diff[src], dest)
}

func (g *DirectedGraph) RemoveEdge(src, dest *Vertex) {
	for idx, destVertex := range g.adjacency[src] {
		if destVertex.id == dest.id {
			g.adjacency[src] = append(g.adjacency[src][:idx], g.adjacency[src][idx+1:]...)
			g.diff[src] = append(g.diff[src][:idx], g.diff[src][idx+1:]...)
		}
	}

	if len(g.adjacency[src]) == 0 {
		delete(g.adjacency, src)
		delete(g.diff, src)
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
	log.Printf("Encoding %d source vertices", len(g.diff))
	for vertex, vertices := range g.diff {
		log.Printf("Encoding %d destination vertices", len(vertices))
		fmt.Fprint(&b, vertex.data, delim, vertex.id.String(), delim, len(vertices), delim)
		for _, v := range vertices {
			fmt.Fprint(&b, v.data, delim, v.id.String(), delim)
		}

		fmt.Fprintln(&b)
	}

	g.diff = make(map[*Vertex][]*Vertex)
	return b.Bytes(), nil
}

func (g *DirectedGraph) UnmarshalBinary(data []byte) error {
	defer func() {
		g.diff = make(map[*Vertex][]*Vertex)
	}()

	b := bytes.NewBuffer(data)
	var srcVData int
	var srcVID string
	var destLen int
	for {
		_, err := fmt.Fscan(b, &srcVData, &srcVID, &destLen)
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

		g.adjacency[srcV] = destVertices
	}
}
