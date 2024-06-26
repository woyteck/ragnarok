package db

import (
	"context"

	"github.com/google/uuid"
	"woyteck.pl/ragnarok/types"
)

type MemoriesStore interface {
	Truncater
	GetMemoryByUUID(context.Context, uuid.UUID) (*types.Memory, error)
	GetMemoryBySource(context.Context, string) (isFound bool, memory *types.Memory, err error)
	GetMemories(ctx context.Context, allFields bool, limit int, offset int) ([]*types.Memory, error)
	InsertMemory(context.Context, *types.Memory) error
	UpdateMemory(context.Context, *types.Memory) error
}
