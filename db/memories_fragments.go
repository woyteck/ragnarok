package db

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"woyteck.pl/ragnarok/types"
)

type MemoriesFragmentsStore interface {
	Truncater
	GetMemoryFragmentByUUID(context.Context, uuid.UUID) (*types.MemoryFragment, error)
	GetMemoryFragmentByMemoryUUID(context.Context, uuid.UUID) (*types.MemoryFragment, error)
	InsertMemoryFragment(context.Context, *types.MemoryFragment) (*types.MemoryFragment, error)
}

type PostgresMemoriesFragmentsStore struct {
	db    *sql.DB
	table string
}

func NewPostgresMemoriesFragmentsStore(db *sql.DB, table string) *PostgresMemoriesFragmentsStore {
	return &PostgresMemoriesFragmentsStore{
		db:    db,
		table: table,
	}
}

//TODO
