package api

import (
	"database/sql"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/google/uuid"
	"woyteck.pl/ragnarok/db"
	"woyteck.pl/ragnarok/types"
)

type ConversationResponse struct {
	Conversation types.Conversation `json:"conversation"`
}

type ConversationHandler struct {
	db    *sql.DB
	store *db.Store
}

func NewConversationHandler(db *sql.DB, store *db.Store) *ConversationHandler {
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

	var conv *types.Conversation
	if isNew {
		conv = &types.Conversation{}
	} else {
		conv, err = h.store.Conversation.GetConversationByUUID(c.Context(), validId)
		if err != nil {
			log.Error(err)
			return ErrResourceNotFound("conversation")
		}

		messages, err := h.store.Message.GetMessagesByConversationUUID(c.Context(), validId)
		if err != nil {
			log.Error(err)
		}

		conv.Messages = messages
	}

	return c.JSON(conv)
}
