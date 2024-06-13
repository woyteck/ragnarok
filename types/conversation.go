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

type Coords struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}
