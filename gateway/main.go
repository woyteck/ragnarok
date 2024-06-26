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
	http.ListenAndServe(listernAddr, nil)
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

	v1.Get("/conversation/:uuid?", conversationHandler.HandleGetConversation)
	v1.Post("/conversation/:uuid", conversationHandler.HandlePostConversation)

	v1.Get("/memories", memoryHandler.HandleGetMemories)
	v1.Get("/memories/:uuid", memoryHandler.HandleGetMemory)
	v1.Post("/memories", memoryHandler.HandlePostMemory)

	log.Fatal(app.Listen(listernAddr))
}
