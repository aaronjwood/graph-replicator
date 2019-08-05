package graph

import "github.com/google/uuid"

type Vertex struct {
	data int
	id   uuid.UUID
}

func NewVertex(data int) *Vertex {
	return &Vertex{
		data: data,
		id:   uuid.New(),
	}
}
