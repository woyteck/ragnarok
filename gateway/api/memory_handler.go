package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/google/uuid"
	"woyteck.pl/ragnarok/db"
	"woyteck.pl/ragnarok/di"
	"woyteck.pl/ragnarok/kafka"
	"woyteck.pl/ragnarok/types"
)

type GetMemoriesResponse struct {
	Memories []*types.Memory `json:"memories"`
}

type PostMemoryRequest struct {
	Source      string `json:"source"`
	MemoryType  string `json:"type"`
	CssSelector string `json:"cssSelector,omitempty"`
}

type MemoryHandler struct {
	store *db.Store
	kafka *kafka.Kafka
}

func NewMemoryHandler(container *di.Container) *MemoryHandler {
	store, ok := container.Get("store").(*db.Store)
	if !ok {
		panic("get store failed")
	}

	kafka, ok := container.Get("kafka").(*kafka.Kafka)
	if !ok {
		panic("get kafka failed")
	}

	return &MemoryHandler{
		store: store,
		kafka: kafka,
	}
}

func (h *MemoryHandler) HandleGetMemories(c *fiber.Ctx) error {
	memories, err := h.store.Memory.GetMemories(c.Context(), false, 0, 0)
	if err != nil {
		log.Error(err)
		return ErrInternalError("something went wrong")
	}

	response := GetMemoriesResponse{
		Memories: memories,
	}

	c.Response().Header.Add("Access-Control-Allow-Origin", "*") //FIXME: create CORS middleware and move this there

	return c.JSON(response)
}

func (h *MemoryHandler) HandleGetMemory(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("uuid"))
	if err != nil {
		return ErrBadRequest()
	}

	memory, err := h.store.Memory.GetMemoryByUUID(c.Context(), id)
	if err != nil {
		log.Error(err)
		return ErrInternalError("something went wrong")
	}

	fragments, err := h.store.MemoryFragment.GetMemoryFragmentsByMemoryUUID(c.Context(), id)
	if err != nil {
		log.Error(err)
		return ErrInternalError("something went wrong")
	}
	memory.Fragments = fragments

	c.Response().Header.Add("Access-Control-Allow-Origin", "*") //FIXME: create CORS middleware and move this there

	return c.JSON(memory)
}

func (h *MemoryHandler) HandlePostMemory(c *fiber.Ctx) error {
	var req PostMemoryRequest

	err := json.Unmarshal(c.Request().Body(), &req)
	if err != nil {
		return ErrBadRequest()
	}

	if req.MemoryType != types.MemoryTypeWebArticle {
		return ErrBadRequest()
	}

	if req.Source == "" {
		return ErrBadRequest()
	}

	if req.MemoryType == types.MemoryTypeWebArticle && req.CssSelector == "" {
		return ErrBadRequest()
	}

	isFound, _, err := h.store.Memory.GetMemoryBySource(c.Context(), req.Source)
	if err != nil {
		log.Error(err)
		return ErrInternalError("something went wrong")
	}

	if isFound {
		return ErrResourceAlreadyExists(fmt.Sprintf("memory: %s", req.Source))
	}

	memory := types.NewMemory(req.MemoryType, req.Source, "")
	h.store.Memory.InsertMemory(c.Context(), memory)

	//produce kafka event
	task := types.ScrapTaskEvent{
		Url:         req.Source,
		CssSelector: req.CssSelector,
	}
	message, err := json.Marshal(task)
	if err != nil {
		log.Error(err)
		return ErrInternalError("something went wrong")
	}
	h.kafka.Produce("scrap_jobs", message)

	c.Response().Header.Add("Access-Control-Allow-Origin", "*") //FIXME: create CORS middleware and move this there

	return c.Status(http.StatusCreated).JSON(memory)
}
