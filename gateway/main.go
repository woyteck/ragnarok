package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"woyteck.pl/ragnarok/di"
	"woyteck.pl/ragnarok/gateway/api"

	_ "github.com/lib/pq"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	go startApi()

	listernAddr := os.Getenv("WEBSOCKET_LISTEN_ADDRESS")
	fmt.Println("Starting ws server on", listernAddr)

	container := di.NewContainer(di.Services)
	conversationHandler := api.NewConversationHandler(container)
	server := NewServer(listernAddr, conversationHandler.HandleWsConversation)

	http.HandleFunc("/ws", server.handleWS)
	log.Fatal(http.ListenAndServe(listernAddr, nil))
}

func corsMiddleware(c *fiber.Ctx) error {
	c.Response().Header.Add("Access-Control-Allow-Origin", "*")
	return c.Next()
}

func startApi() {
	var config = fiber.Config{
		ErrorHandler: api.ErrorHandler,
	}

	var (
		listernAddr         = os.Getenv("API_HTTP_LISTEN_ADDRESS")
		container           = di.NewContainer(di.Services)
		conversationHandler = api.NewConversationHandler(container)
		memoryHandler       = api.NewMemoryHandler(container)
		app                 = fiber.New(config)
		apiRoot             = app.Group("/api")
		v1                  = apiRoot.Group("v1")
	)

	v1.Get("/conversation/:uuid?", corsMiddleware, conversationHandler.HandleGetConversation)
	v1.Post("/conversation/:uuid", corsMiddleware, conversationHandler.HandlePostConversation)

	v1.Get("/memories", corsMiddleware, memoryHandler.HandleGetMemories)
	v1.Get("/memories/:uuid", corsMiddleware, memoryHandler.HandleGetMemory)
	v1.Post("/memories", corsMiddleware, memoryHandler.HandlePostMemory)

	log.Fatal(app.Listen(listernAddr))
}
