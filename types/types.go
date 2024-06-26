package types

import (
	"time"

	"github.com/google/uuid"
)

const MemoryTypeWebArticle = "web_article"
const MemoryTypeTextFile = "text_file"

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
	UpdatedAt  *time.Time `json:"updatedAt,omitempty"`
	DeletedAt  *time.Time `json:"deletedAt,omitempty"`
	MemoryType string     `json:"memoryType"`
	Source     string     `json:"source,omitempty"`
	Content    string     `json:"content,omitempty"`
}

func NewMemory(memoryType string, source string, content string) *Memory {
	now := time.Now()
	return &Memory{
		ID:         uuid.New(),
		CreatedAt:  &now,
		MemoryType: memoryType,
		Source:     source,
		Content:    content,
	}
}

type MemoryFragment struct {
	ID              uuid.UUID  `json:"id"`
	CreatedAt       *time.Time `json:"createdAt"`
	ContentOriginal string     `json:"contentOriginal"`
	ContentRefined  string     `json:"contentRefined"`
	IsRefined       bool       `json:"isRefined"`
	IsEmbedded      bool       `json:"isEmbedded"`
	MemoryID        uuid.UUID  `json:"memoryId"`
}

func NewMemoryFragment(contentOriginal string, contentRefined string, isRefined bool, isEmbedded bool, memoryID uuid.UUID) *MemoryFragment {
	now := time.Now()
	return &MemoryFragment{
		ID:              uuid.New(),
		CreatedAt:       &now,
		ContentOriginal: contentOriginal,
		ContentRefined:  contentRefined,
		IsRefined:       isRefined,
		IsEmbedded:      isEmbedded,
		MemoryID:        memoryID,
	}
}
