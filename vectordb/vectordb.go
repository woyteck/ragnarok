package vectordb

import (
	"github.com/google/uuid"
)

type Collection struct {
	Name        string
	PointsCount int
}

type Record struct {
	ID      uuid.UUID
	Score   float64
	Payload map[string]any
	Version int
}

type VectorClientCollections interface {
	GetCollection(collectionName string) (*Collection, error)
	GetCollections() ([]*Collection, error)
	CreateCollection(collectionName string) error
	DeleteCollection(collectionName string) error
}

type VectorClientPoints interface {
	UpsertPoints(collectionName string, vector []float64, id uuid.UUID, payload map[string]any) error
	Search(collectionName string, vector []float64, resultsCount int) ([]*Record, error)
}
