package db

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"woyteck.pl/ragnarok/types"
)

type MessagesStore interface {
	Truncater
	GetMessageByUUID(context.Context, uuid.UUID) (*types.Message, error)
	GetMessagesByConversationUUID(context.Context, uuid.UUID) ([]*types.Message, error)
	InsertMessage(context.Context, *types.Message) (*types.Message, error)
}

type PostgresMessagesStore struct {
	db    *sql.DB
	table string
}

func NewPostgresMessagesStore(db *sql.DB, table string) *PostgresMessagesStore {
	return &PostgresMessagesStore{
		db:    db,
		table: table,
	}
}

func (s *PostgresMessagesStore) Truncate(ctx context.Context) error {
	query := fmt.Sprintf("DELETE FROM %s", s.table)
	fmt.Println(query)
	s.db.Exec(query)
	return nil
}

func (s *PostgresMessagesStore) GetMessageByUUID(ctx context.Context, id uuid.UUID) (*types.Message, error) {
	var conversationId uuid.UUID
	var role string
	var content string
	var createdAt sql.NullString

	query := fmt.Sprintf("SELECT conversation_id, role, content, created_at FROM %s WHERE uuid = $1", s.table)
	row := s.db.QueryRow(query, id)

	switch err := row.Scan(&conversationId, &role, &content, &createdAt); err {
	case sql.ErrNoRows:
		return nil, fmt.Errorf("conversation not found")
	case nil:
		message := &types.Message{
			ID:             id,
			ConversationId: conversationId,
			Role:           role,
			Content:        content,
		}

		createdAtTime, err := parseTimestamp(createdAt)
		if err == nil {
			message.CreatedAt = createdAtTime
		}

		return message, nil
	default:
		return nil, err
	}
}

func (s *PostgresMessagesStore) GetMessagesByConversationUUID(ctx context.Context, id uuid.UUID) ([]*types.Message, error) {
	query := fmt.Sprintf("SELECT conversation_id, role, content, created_at FROM %s WHERE conversation_id = $1", s.table)
	rows, err := s.db.Query(query, id)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	messages := []*types.Message{}

	for rows.Next() {
		var conversationId uuid.UUID
		var role string
		var content string
		var createdAt sql.NullString

		err := rows.Scan(&conversationId, &role, &content, &createdAt)
		if err != nil {
			return nil, err
		}

		message := &types.Message{
			ID:             id,
			ConversationId: conversationId,
			Role:           role,
			Content:        content,
		}

		createdAtTime, err := parseTimestamp(createdAt)
		if err == nil {
			message.CreatedAt = createdAtTime
		}

		messages = append(messages, message)
	}

	return messages, nil
}

func (s *PostgresMessagesStore) InsertMessage(ctx context.Context, m *types.Message) (*types.Message, error) {
	if m.ID == uuid.Nil {
		return nil, fmt.Errorf("can't insert message with no ID")
	}

	cols := []string{"uuid", "conversation_id", "role", "content"}
	values := []any{m.ID, m.ConversationId, m.Role, m.Content}

	if !m.CreatedAt.IsZero() {
		cols = append(cols, "created_at")
		values = append(values, m.CreatedAt)
	}

	query := fmt.Sprintf("INSERT INTO Messages (%s) VALUES (%s)", strings.Join(cols, ","), makePlaceholders(len(values)))

	_, err := s.db.Exec(query, values...)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return m, nil
}
