package types

import (
	"time"

	"github.com/google/uuid"
)

type Conversation struct {
	ID        uuid.UUID  `json:"id" db:"uuid"`
	CreatedAt *time.Time `json:"created_at,omitempty" db:"created_at"`
	UpdatedAt *time.Time `json:"updated_at,omitempty" db:"updated_at"`
	Messages  []*Message `json:"messages,omitempty"`
}

type Message struct {
	ID             uuid.UUID  `json:"id"`
	ConversationId uuid.UUID  `json:"conversation_id"`
	Role           string     `json:"role"`
	Content        string     `json:"content"`
	CreatedAt      *time.Time `json:"created_at,omitempty"`
}

func NewMessage(conversationID uuid.UUID, role string, content string) *Message {
	now := time.Now()
	message := &Message{
		ID:        uuid.New(),
		CreatedAt: &now,
	}

	message.ConversationId = conversationID
	message.Role = role
	message.Content = content

	return message
}

type Coords struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type Memory struct {
	ID         uuid.UUID  `json:"id"`
	CreatedAt  *time.Time `json:"createdAt"`
	UpdatedAt  *time.Time `json:"updatedAt"`
	DeletedAt  *time.Time `json:"deletedAt"`
	MemoryType string     `json:"memoryType"`
	Source     string     `json:"source"`
	Content    string     `json:"content"`
}

type MemoryFragment struct {
	ID              uuid.UUID  `json:"id"`
	CreatedAt       *time.Time `json:"createdAt"`
	ContentOriginal string     `json:"contentOriginal"`
	ContentRefined  string     `json:"contentRefined"`
	IsRefined       bool       `json:"isRefined"`
	IsEmbedded      bool       `json:"isEmbedded"`
	MemoryId        uuid.UUID  `json:"memoryId"`
}
