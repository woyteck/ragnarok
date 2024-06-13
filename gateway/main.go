package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"woyteck.pl/ragnarok/db"
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

	user := os.Getenv("DB_USER")
	pass := os.Getenv("DB_PASSWORD")
	name := os.Getenv("DB_NAME")
	host := os.Getenv("DB_HOST")
	db := db.Connect(host, user, pass, name)

	listernAddr := os.Getenv("API_HTTP_LISTEN_ADDRESS")

	conversationHandler := api.NewConversationHandler(db)

	app := fiber.New(config)
	apiRoot := app.Group("/api")
	v1 := apiRoot.Group("v1")
	v1.Get("/conversation/:uuid?", conversationHandler.HandleGetConversation)

	log.Fatal(app.Listen(listernAddr))
}
