package db

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"woyteck.pl/ragnarok/types"
)

type MemoriesFragmentsStore interface {
	Truncater
	GetMemoryFragmentByUUID(context.Context, uuid.UUID) (*types.MemoryFragment, error)
	GetMemoryFragmentsByMemoryUUID(context.Context, uuid.UUID) ([]*types.MemoryFragment, error)
	InsertMemoryFragment(context.Context, *types.MemoryFragment) error
	UpdateMemoryFragment(context.Context, *types.MemoryFragment) error
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

func (s *PostgresMemoriesFragmentsStore) GetMemoryFragmentByUUID(ctx context.Context, id uuid.UUID) (*types.MemoryFragment, error) {
	var createdAt sql.NullString
	var contentOriginal string
	var contentRefined string
	var isRefined bool
	var isEmbedded bool
	var memoryID uuid.UUID

	query := fmt.Sprintf("SELECT created_at, content_original, content_refined, is_refined, is_embedded, memory_id FROM %s WHERE uuid = $1", s.table)
	row := s.db.QueryRow(query, id)

	switch err := row.Scan(&createdAt, &contentOriginal, &contentRefined, &isRefined, &isEmbedded, &memoryID); err {
	case sql.ErrNoRows:
		return nil, fmt.Errorf("memory not found")
	case nil:
		fragment := &types.MemoryFragment{
			ID:              id,
			ContentOriginal: contentOriginal,
			ContentRefined:  contentRefined,
			IsRefined:       isRefined,
			IsEmbedded:      isEmbedded,
			MemoryID:        memoryID,
		}

		createdAtTime, err := parseTimestamp(createdAt)
		if err == nil {
			fragment.CreatedAt = createdAtTime
		}

		return fragment, nil
	default:
		return nil, err
	}
}

func (s *PostgresMemoriesFragmentsStore) GetMemoryFragmentsByMemoryUUID(ctx context.Context, memoryID uuid.UUID) ([]*types.MemoryFragment, error) {
	query := fmt.Sprintf("SELECT uuid, created_at, content_original, content_refined, is_refined, is_embedded FROM %s WHERE memory_id = $1", s.table)
	rows, err := s.db.Query(query, memoryID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	fragments := []*types.MemoryFragment{}

	for rows.Next() {
		var id uuid.UUID
		var createdAt sql.NullString
		var contentOriginal string
		var contentRefined string
		var isRefined bool
		var isEmbedded bool

		err := rows.Scan(&id, &createdAt, &contentOriginal, &contentRefined, &isRefined, &isEmbedded)
		if err != nil {
			return nil, err
		}

		fragment := &types.MemoryFragment{
			ID:              id,
			ContentOriginal: contentOriginal,
			ContentRefined:  contentRefined,
			IsRefined:       isRefined,
			IsEmbedded:      isEmbedded,
			MemoryID:        memoryID,
		}

		createdAtTime, err := parseTimestamp(createdAt)
		if err == nil {
			fragment.CreatedAt = createdAtTime
		}

		fragments = append(fragments, fragment)
	}

	return fragments, nil
}

func (s *PostgresMemoriesFragmentsStore) InsertMemoryFragment(ctx context.Context, f *types.MemoryFragment) error {
	if f.ID == uuid.Nil {
		return fmt.Errorf("can't insert memory fragment with no ID")
	}

	if f.MemoryID == uuid.Nil {
		return fmt.Errorf("can't insert memory fragment with no MemoryId")
	}

	cols := []string{"uuid", "content_original", "content_refined", "is_refined", "is_embedded", "memory_id"}
	values := []any{f.ID, f.ContentOriginal, f.ContentRefined, f.IsRefined, f.IsEmbedded, f.MemoryID}

	if f.CreatedAt != nil && !f.CreatedAt.IsZero() {
		cols = append(cols, "created_at")
		values = append(values, f.CreatedAt)
	}

	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", s.table, strings.Join(cols, ","), makePlaceholders(len(values)))

	_, err := s.db.Exec(query, values...)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func (s *PostgresMemoriesFragmentsStore) UpdateMemoryFragment(ctx context.Context, f *types.MemoryFragment) error {
	if f.ID == uuid.Nil {
		return fmt.Errorf("can't update memory fragment with no ID")
	}

	if f.MemoryID == uuid.Nil {
		return fmt.Errorf("can't update memory fragment with no MemoryId")
	}

	query := fmt.Sprintf("UPDATE %s SET content_original=$1, content_refined=$2, is_refined=$3, is_embedded=$4, memory_id=$5 WHERE id=$6", s.table)
	_, err := s.db.Exec(query, f.ContentOriginal, f.ContentRefined, f.IsRefined, f.IsEmbedded, f.MemoryID, f.ID)

	return err
}
