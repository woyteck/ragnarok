package db

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"woyteck.pl/ragnarok/types"
)

type MemoriesStore interface {
	Truncater
	GetMemoryByUUID(context.Context, uuid.UUID) (*types.Memory, error)
	InsertMemory(context.Context, *types.Memory) (*types.Memory, error)
}

type PostgresMemoriesStore struct {
	db    *sql.DB
	table string
}

func NewPostgresMemoriesStore(db *sql.DB, table string) *PostgresMemoriesStore {
	return &PostgresMemoriesStore{
		db:    db,
		table: table,
	}
}

func (s *PostgresMemoriesStore) Truncate(ctx context.Context) error {
	query := fmt.Sprintf("DELETE FROM %s", s.table)
	fmt.Println(query)
	s.db.Exec(query)

	return nil
}

func (s *PostgresMemoriesStore) GetMemoryByUUID(ctx context.Context, id uuid.UUID) (*types.Memory, error) {
	var createdAt sql.NullString
	var updatedAt sql.NullString
	var deletedAt sql.NullString
	var memoryType string
	var source string
	var content string

	query := fmt.Sprintf("SELECT created_at, updated_at, deleted_at, memory_type, source, content FROM %s WHERE uuid = $1", s.table)
	row := s.db.QueryRow(query, id)

	switch err := row.Scan(&createdAt, &updatedAt, &deletedAt, &memoryType, &source, &content); err {
	case sql.ErrNoRows:
		return nil, fmt.Errorf("memory not found")
	case nil:
		memory := &types.Memory{
			ID:         id,
			MemoryType: memoryType,
			Source:     source,
			Content:    content,
		}

		createdAtTime, err := parseTimestamp(createdAt)
		if err == nil {
			memory.CreatedAt = createdAtTime
		}

		updatedAtTime, err := parseTimestamp(updatedAt)
		if err == nil {
			memory.UpdatedAt = updatedAtTime
		}

		deletedAtTime, err := parseTimestamp(deletedAt)
		if err == nil {
			memory.DeletedAt = deletedAtTime
		}

		return memory, nil
	default:
		return nil, err
	}
}

func (s *PostgresMemoriesStore) InsertMemory(ctx context.Context, m *types.Memory) (*types.Memory, error) {
	if m.ID == uuid.Nil {
		return nil, fmt.Errorf("can't insert memory with no ID")
	}

	cols := []string{"uuid", "memory_type", "source", "content"}
	values := []any{m.ID, m.MemoryType, m.Source, m.Content}

	if m.CreatedAt != nil && !m.CreatedAt.IsZero() {
		cols = append(cols, "created_at")
		values = append(values, m.CreatedAt)
	}

	if m.UpdatedAt != nil && !m.UpdatedAt.IsZero() {
		cols = append(cols, "updated_at")
		values = append(values, m.UpdatedAt)
	}

	if m.DeletedAt != nil && !m.DeletedAt.IsZero() {
		cols = append(cols, "deleted_at")
		values = append(values, m.DeletedAt)
	}

	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", strings.Join(cols, ","), s.table, makePlaceholders(len(values)))

	_, err := s.db.Exec(query, values...)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return m, nil
}
