package api

import (
	"database/sql"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/google/uuid"
	"woyteck.pl/ragnarok/db"
	"woyteck.pl/ragnarok/openai"
	"woyteck.pl/ragnarok/rag"
	"woyteck.pl/ragnarok/types"
)

type ConversationResponse struct {
	Conversation types.Conversation `json:"conversation"`
}

type TalkRequest struct {
	Text   string       `json:"text"`
	Coords types.Coords `json:"coords,omitempty"`
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
			return ErrResourceNotFound("conversation")
		}

		messages, err := h.store.Message.GetMessagesByConversationUUID(c.Context(), validId)
		if err != nil {
			log.Error(err)
			return ErrInternalError("something went wrong")
		}
		conv.Messages = messages
	}

	return c.JSON(conv)
}

func (h *ConversationHandler) HandlePostConversation(c *fiber.Ctx) error {
	uuidParam := c.Params("uuid")
	err := uuid.Validate(uuidParam)
	if err != nil {
		return ErrBadRequest()
	}

	validId := uuid.MustParse(uuidParam)
	conv, err := h.store.Conversation.GetConversationByUUID(c.Context(), validId)
	if err != nil {
		return ErrResourceNotFound("conversation")
	}

	messages, err := h.store.Message.GetMessagesByConversationUUID(c.Context(), validId)
	if err != nil {
		log.Error(err)
		return ErrInternalError("something went wrong")
	}
	conv.Messages = messages

	request := TalkRequest{}
	if err := c.BodyParser(&request); err != nil {
		return ErrBadRequest()
	}

	now := time.Now()
	message := &types.Message{
		ID:             uuid.New(),
		ConversationId: conv.ID,
		Role:           "user",
		Content:        request.Text,
		CreatedAt:      &now,
	}
	conv.Messages = append(conv.Messages, message)
	_, err = h.store.Message.InsertMessage(c.Context(), message)
	if err != nil {
		log.Error(err)
		return ErrBadRequest()
	}

	config := openai.Config{ //TODO: dependency injection
		ApiKey: os.Getenv("OPENAI_API_KEY"),
	}
	llm := openai.NewClient(config)
	rag := rag.New(&llm, h.store.Message)
	err = rag.Ask(conv)
	if err != nil {
		log.Error(err)
		return ErrInternalError("something went wrong")
	}

	return c.JSON(conv)
}
