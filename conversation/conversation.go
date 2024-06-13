package conversation

import (
	"database/sql"
	"fmt"

	"github.com/google/uuid"
)

type Conversationer interface {
	Create(context string) (uuid.UUID, error)
	AddMessage(role string, content string) error
	GetMessages() ([]*Message, error)
}

type Message struct {
	Role    string
	Content string
}

type Conversation struct {
	db       *sql.DB
	ID       uuid.UUID  `json:"id" db:"uuid"`
	Messages []*Message `json:"messages"`
}

func New(db *sql.DB) *Conversation {
	return &Conversation{
		db:       db,
		Messages: []*Message{},
	}
}

func Get(db *sql.DB, id uuid.UUID) (*Conversation, error) {
	var created_at, updated_at string

	query := "SELECT created_at, updated_at FROM conversations WHERE uuid = $1 AND deleted_at IS NULL"
	row := db.QueryRow(query, id)
	switch err := row.Scan(&created_at, &updated_at); err {
	case sql.ErrNoRows:
		return nil, fmt.Errorf("conversation not found")
	case nil:
		conv := &Conversation{
			db:       db,
			ID:       id,
			Messages: []*Message{},
		}

		messages, err := conv.getMessages()
		if err != nil {
			return nil, err
		}
		conv.Messages = messages

		return conv, nil
	default:
		return nil, err
	}
}

func (c *Conversation) Create(context string) (uuid.UUID, error) {
	c.ID = uuid.New()

	row, err := c.db.Exec("INSERT INTO conversations (uuid) VALUES ($1)", c.ID)
	if err != nil {
		return uuid.Nil, err
	}
	fmt.Println(row)

	err = c.AddMessaage("system", context)
	if err != nil {
		return uuid.Nil, err
	}

	return c.ID, nil
}

func (c *Conversation) AddMessaage(role string, content string) error {
	c.Messages = append(c.Messages, &Message{
		Role:    role,
		Content: content,
	})
	_, err := c.db.Exec("INSERT INTO messages (role, content, conversation_id) VALUES ($1, $2, $3)", role, content, c.ID)

	return err
}

func (c *Conversation) getMessages() ([]*Message, error) {
	query := "SELECT role, content FROM messages WHERE conversation_id = $1"
	rows, err := c.db.Query(query, c.ID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	messages := []*Message{}
	for rows.Next() {
		message := Message{}
		rows.Scan(&message.Role, &message.Content)
		messages = append(messages, &message)
	}

	return messages, nil
}
