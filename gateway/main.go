package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"woyteck.pl/ragnarok/di"
	"woyteck.pl/ragnarok/gateway/api"
)

var config = fiber.Config{
	ErrorHandler: api.ErrorHandler,
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	container := di.NewContainer(di.Services)

	conversationHandler := api.NewConversationHandler(container)

	app := fiber.New(config)
	apiRoot := app.Group("/api")
	v1 := apiRoot.Group("v1")
	v1.Get("/conversation/:uuid?", conversationHandler.HandleGetConversation)
	v1.Post("/conversation/:uuid", conversationHandler.HandlePostConversation)

	listernAddr := os.Getenv("API_HTTP_LISTEN_ADDRESS")
	log.Fatal(app.Listen(listernAddr))
}
