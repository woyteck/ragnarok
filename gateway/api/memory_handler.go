package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"woyteck.pl/ragnarok/db"
	"woyteck.pl/ragnarok/di"
	"woyteck.pl/ragnarok/types"
)

type GetMemoriesResponse struct {
	Memories []*types.Memory `json:"memories"`
}

type MemoryHandler struct {
	store *db.Store
}

func NewMemoryHandler(container *di.Container) *MemoryHandler {
	store, ok := container.Get("store").(*db.Store)
	if !ok {
		panic("get store failed")
	}

	return &MemoryHandler{
		store: store,
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

	return c.JSON(response)
}
