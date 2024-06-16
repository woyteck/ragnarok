package types

import "github.com/google/uuid"

type ScrapTaskEvent struct {
	Url         string `json:"url"`
	CssSelector string `json:"cssSelector"`
}

type IndexMemoryFragmentEvent struct {
	MemoryFragmentID uuid.UUID
}

type EmbedMemoryFragmentEvent struct {
	MemoryFragmentID uuid.UUID
}
