package api

import (
	"database/sql"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/google/uuid"
	"woyteck.pl/ragnarok/conversation"
)

type ConversationHandler struct {
	db *sql.DB
}

func NewConversationHandler(db *sql.DB) *ConversationHandler {
	return &ConversationHandler{
		db: db,
	}
}

func (h *ConversationHandler) HandleGetConversation(c *fiber.Ctx) error {

	isNew := true
	var validId uuid.UUID
	var conv *conversation.Conversation

	uuidParam := c.Params("uuid")
	err := uuid.Validate(uuidParam)
	if err == nil {
		validId = uuid.MustParse(uuidParam)
		isNew = false
	}

	if isNew {
		conv = conversation.New(h.db)
		_, err = conv.Create("test context")
		if err != nil {
			log.Error(err)
			return ErrInternalError("something went wrong")
		}
	} else {
		conv, err = conversation.Get(h.db, validId)
		if err != nil {
			log.Error(err)
			return ErrResourceNotFound("conversation")
		}
	}

	return c.JSON(conv)
}
