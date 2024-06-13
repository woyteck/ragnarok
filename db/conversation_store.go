package db

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"woyteck.pl/ragnarok/types"
)

type ConversationStore interface {
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
	query := fmt.Sprintf("DELETE FROM %s", s.table)
	fmt.Println(query)
	s.db.Exec(query)
	return nil
}

func (s *PostgresConversationStore) GetConversationByUUID(ctx context.Context, uuid uuid.UUID) (*types.Conversation, error) {
	var createdAt, updatedAt sql.NullString

	query := fmt.Sprintf("SELECT created_at, updated_at FROM %s WHERE uuid = $1", s.table)
	row := s.db.QueryRow(query, uuid)

	switch err := row.Scan(&createdAt, &updatedAt); err {
	case sql.ErrNoRows:
		return nil, fmt.Errorf("conversation not found")
	case nil:
		conv := &types.Conversation{
			ID: uuid,
		}

		createdAtTime, err := parseTimestamp(createdAt)
		if err == nil {
			conv.CreatedAt = createdAtTime
		}

		updatedAtTime, err := parseTimestamp(updatedAt)
		if err == nil {
			conv.UpdatedAt = updatedAtTime
		}

		return conv, nil
	default:
		return nil, err
	}
}

func (s *PostgresConversationStore) InsertConversation(ctx context.Context, c *types.Conversation) (*types.Conversation, error) {
	if c.ID == uuid.Nil {
		return nil, fmt.Errorf("can't insert conversation with no ID")
	}

	cols := []string{"uuid"}
	values := []any{c.ID}

	if c.CreatedAt != nil && !c.CreatedAt.IsZero() {
		cols = append(cols, "created_at")
		values = append(values, c.CreatedAt)
	}

	if c.UpdatedAt != nil && !c.UpdatedAt.IsZero() {
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
