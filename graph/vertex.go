package graph

import "github.com/google/uuid"

type Vertex struct {
	data int
	id   uuid.UUID
}

type vertexSet map[Vertex]bool

func (s vertexSet) Add(v Vertex) {
	s[v] = true
}

func (s vertexSet) Remove(v Vertex) {
	delete(s, v)
}

func (s vertexSet) Contains(v Vertex) bool {
	_, ok := s[v]
	return ok
}

func NewVertex(data int) *Vertex {
	return &Vertex{
		data: data,
		id:   uuid.New(),
	}
}
