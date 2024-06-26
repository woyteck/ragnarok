package db

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

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

func (s *PostgresMemoriesStore) GetMemoryBySource(ctx context.Context, source string) (isFound bool, memory *types.Memory, err error) {
	var id uuid.UUID
	var createdAt sql.NullString
	var updatedAt sql.NullString
	var deletedAt sql.NullString
	var memoryType string
	var content string

	query := fmt.Sprintf("SELECT uuid, created_at, updated_at, deleted_at, memory_type, source, content FROM %s WHERE source = $1", s.table)
	row := s.db.QueryRow(query, source)

	switch err := row.Scan(&id, &createdAt, &updatedAt, &deletedAt, &memoryType, &source, &content); err {
	case sql.ErrNoRows:
		return false, nil, nil
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

		return true, memory, nil
	default:
		return false, nil, err
	}
}

func (s *PostgresMemoriesStore) GetMemories(ctx context.Context, allFields bool, limit int, offset int) ([]*types.Memory, error) {
	//TODO: add pagination

	fields := []string{"uuid", "created_at", "updated_at", "deleted_at", "memory_type", "source"}
	if allFields {
		fields = append(fields, "content")
	}
	query := fmt.Sprintf("SELECT %s FROM %s ORDER BY created_at", strings.Join(fields, ","), s.table)
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	memories := []*types.Memory{}

	for rows.Next() {
		var memoryId uuid.UUID
		var createdAt sql.NullString
		var updatedAt sql.NullString
		var deletedAt sql.NullString
		var memoryType string
		var source string
		var content string

		var err error
		if allFields {
			err = rows.Scan(&memoryId, &createdAt, &updatedAt, &deletedAt, &memoryType, &source, &content)
		} else {
			err = rows.Scan(&memoryId, &createdAt, &updatedAt, &deletedAt, &memoryType, &source)
		}
		if err != nil {
			return nil, err
		}

		memory := &types.Memory{
			ID:         memoryId,
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

		memories = append(memories, memory)
	}

	return memories, nil
}

func (s *PostgresMemoriesStore) InsertMemory(ctx context.Context, m *types.Memory) error {
	if m.ID == uuid.Nil {
		return fmt.Errorf("can't insert memory with no ID")
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

	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", s.table, strings.Join(cols, ","), makePlaceholders(len(values)))

	_, err := s.db.Exec(query, values...)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func (s *PostgresMemoriesStore) UpdateMemory(ctx context.Context, m *types.Memory) error {
	if m.ID == uuid.Nil {
		return fmt.Errorf("can't update memory with no ID")
	}

	now := time.Now()
	m.UpdatedAt = &now
	query := fmt.Sprintf("UPDATE %s SET content=$1, memory_type=$2, source=$3, updated_at=$4 WHERE uuid=$5", s.table)
	_, err := s.db.Exec(query, m.Content, m.MemoryType, m.Source, m.UpdatedAt, m.ID)

	return err
}
