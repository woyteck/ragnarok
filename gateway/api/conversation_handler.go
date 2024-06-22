package api

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"woyteck.pl/ragnarok/db"
	"woyteck.pl/ragnarok/di"
	"woyteck.pl/ragnarok/openai"
	"woyteck.pl/ragnarok/rag"
	"woyteck.pl/ragnarok/tts"
	"woyteck.pl/ragnarok/types"
)

type ConversationResponse struct {
	Conversation types.Conversation `json:"conversation"`
}

type TalkRequest struct {
	ConversationId uuid.UUID    `json:"conversationId,omitempty"`
	Text           string       `json:"text"`
	Coords         types.Coords `json:"coords,omitempty"`
}

type TalkResponse struct {
	ConversationId uuid.UUID `json:"conversationId,omitempty"`
	Text           string    `json:"text"`
	Audio          string    `json:"audio"`
}

type ConversationHandler struct {
	db    *sql.DB
	store *db.Store
	llm   *openai.Client
	rag   *rag.Rag
	tts   *tts.ElevenLabsTTS
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
	tts, ok := container.Get("tts").(*tts.ElevenLabsTTS)
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
	id, err := uuid.Parse(c.Params("uuid"))
	if err != nil {
		return ErrBadRequest()
	}

	request := TalkRequest{}
	if err := c.BodyParser(&request); err != nil {
		return ErrBadRequest()
	}

	conv, err := h.rag.Ask(c.Context(), id, request.Text)
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

func (h *ConversationHandler) HandleWsConversation(conn *websocket.Conn, req *TalkRequest) error {
	conv, err := h.rag.Ask(context.Background(), req.ConversationId, req.Text)
	if err != nil {
		return err
	}

	lastMessage := conv.Messages[len(conv.Messages)-1]

	handler := func(msg tts.TextToSpeechStreamingResponse) {
		if len(msg.Audio) == 0 {
			return
		}

		chunk := strings.Join(msg.Alignment.Chars, "")
		fmt.Println(chunk)
		conn.WriteJSON(TalkResponse{
			ConversationId: req.ConversationId,
			Text:           chunk,
			Audio:          msg.Audio,
		})
	}
	wsClient, err := h.tts.NewWsClient(handler)
	if err != nil {
		return err
	}

	wsClient.Send(lastMessage.Content, true)

	return nil
}
