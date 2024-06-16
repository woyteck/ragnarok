package api

import (
	"database/sql"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/google/uuid"
	"woyteck.pl/ragnarok/db"
	"woyteck.pl/ragnarok/di"
	"woyteck.pl/ragnarok/openai"
	"woyteck.pl/ragnarok/prompter"
	"woyteck.pl/ragnarok/rag"
	"woyteck.pl/ragnarok/text_to_speech"
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
	llm   *openai.Client
	rag   *rag.Rag
	tts   *text_to_speech.ElevenLabsTTS
}

func NewConversationHandler(container *di.Container) *ConversationHandler {
	dbConn, ok := container.Get("db").(*sql.DB)
	if !ok {
		panic("get db failed")
	}
	store, ok := container.Get("store").(*db.Store)
	if !ok {
		panic("get store failed")
	}
	llm, ok := container.Get("llm").(*openai.Client)
	if !ok {
		panic("get llm failed")
	}
	rag, ok := container.Get("rag").(*rag.Rag)
	if !ok {
		panic("get rag failed")
	}
	tts, ok := container.Get("tts").(*text_to_speech.ElevenLabsTTS)
	if !ok {
		panic("get tts failed")
	}

	return &ConversationHandler{
		db:    dbConn,
		store: store,
		llm:   llm,
		rag:   rag,
		tts:   tts,
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
		conv = &types.Conversation{
			ID: uuid.New(),
		}
		err := h.store.Conversation.InsertConversation(c.Context(), conv)
		if err != nil {
			log.Error(err)
			return ErrInternalError("something went wrong")
		}
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

	if len(conv.Messages) == 0 {
		pr := prompter.New()
		convContext, err := pr.Get("main_context")
		if err != nil {
			return err
		}

		contextMessage := types.NewMessage(conv.ID, "system", convContext)
		conv.Messages = append(conv.Messages, contextMessage)
		err = h.store.Message.InsertMessage(c.Context(), contextMessage)
		if err != nil {
			log.Error(err)
			return ErrBadRequest()
		}
	}

	message := types.NewMessage(conv.ID, "user", request.Text)
	conv.Messages = append(conv.Messages, message)
	err = h.store.Message.InsertMessage(c.Context(), message)
	if err != nil {
		log.Error(err)
		return ErrBadRequest()
	}

	err = h.rag.Ask(conv)
	if err != nil {
		log.Error(err)
		return ErrInternalError("something went wrong")
	}

	lastMessage := conv.Messages[len(conv.Messages)-1]
	b, err := h.tts.TextToSpeech(lastMessage.Content)
	if err != nil {
		log.Error(err)
		return ErrInternalError("something went wrong")
	}

	c.Set("Content-Type", "audio/mpeg")
	return c.Send(b)
}
