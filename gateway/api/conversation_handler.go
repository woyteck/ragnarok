package api

import (
	"context"
	"database/sql"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"woyteck.pl/ragnarok/db"
	"woyteck.pl/ragnarok/types"
)

type ConversationHandler struct {
	db    *sql.DB
	store db.ConversationStore
}

func NewConversationHandler(db *sql.DB, store db.ConversationStore) *ConversationHandler {
	return &ConversationHandler{
		db:    db,
		store: store,
	}
}

func (h *ConversationHandler) HandleGetConversation(c *fiber.Ctx) error {
	isNew := true
	var validId uuid.UUID

	uuidParam := c.Params("uuid")
	err := uuid.Validate(uuidParam)
	if err == nil {
		validId = uuid.MustParse(uuidParam)
		isNew = false
	}

	if isNew {
		conv := types.Conversation{}
		return c.JSON(conv)
	}

	conv, err := h.store.GetConversationByUUID(context.Background(), validId)
	if err != nil {
		return ErrResourceNotFound("conversation")
	}

	return c.JSON(conv)

}
