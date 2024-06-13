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

// type Dropper interface {
// 	Drop(context.Context) error
// }

type Truncater interface {
	Truncate(context.Context) error
}

type ConversationStore interface {
	// Dropper
	Truncater
	GetConversationByUUID(context.Context, uuid.UUID) (*types.Conversation, error)
	InsertConversation(context.Context, *types.Conversation) (*types.Conversation, error)
}

type PostgresConversationStore struct {
	db    *sql.DB
	table string
}

func NewPostgresConversationStore(db *sql.DB, table string) *PostgresConversationStore {
	return &PostgresConversationStore{
		db:    db,
		table: table,
	}
}

func (s *PostgresConversationStore) Truncate(ctx context.Context) error {
	s.db.Exec(fmt.Sprintf("TRUNCATE %s", s.table))
	return nil
}

func (s *PostgresConversationStore) GetConversationByUUID(ctx context.Context, uuid uuid.UUID) (*types.Conversation, error) {
	var createdAt, updatedAt time.Time
	query := fmt.Sprintf("SELECT created_at, updated_at FROM %s WHERE uuid = $1 AND deleted_at IS NULL", s.table)
	row := s.db.QueryRow(query, uuid)
	switch err := row.Scan(&createdAt, &updatedAt); err {
	case sql.ErrNoRows:
		return nil, fmt.Errorf("conversation not found")
	case nil:
		conv := &types.Conversation{
			ID:        uuid,
			CreatedAt: createdAt,
			UpdatedAt: updatedAt,
		}

		return conv, nil
	default:
		return nil, err
	}
}

func (s *PostgresConversationStore) InsertConversation(ctx context.Context, c *types.Conversation) (*types.Conversation, error) {

	cols := []string{"uuid"}
	values := []any{c.ID}

	if !c.CreatedAt.IsZero() {
		cols = append(cols, "created_at")
		values = append(values, c.CreatedAt)
	}

	if !c.UpdatedAt.IsZero() {
		cols = append(cols, "updated_at")
		values = append(values, c.UpdatedAt)
	}

	query := fmt.Sprintf("INSERT INTO conversations (%s) VALUES (%s)", strings.Join(cols, ","), makePlaceholders(len(values)))

	_, err := s.db.Exec(query, values...)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return c, nil
}

func makePlaceholders(count int) string {
	list := []string{}
	for i := 0; i < count; i++ {
		list = append(list, fmt.Sprintf("$%d", i+1))
	}

	return strings.Join(list, ",")
}
