package di

import (
	"database/sql"
	"os"

	"github.com/gocolly/colly"
	"woyteck.pl/ragnarok/db"
	"woyteck.pl/ragnarok/openai"
	"woyteck.pl/ragnarok/prompter"
	"woyteck.pl/ragnarok/rag"
	"woyteck.pl/ragnarok/scraper"
	"woyteck.pl/ragnarok/text_to_speech"
	"woyteck.pl/ragnarok/vectordb"
)

var Services = map[string]ServiceFactoryFn{
	"llm": func(c *Container) any {
		config := openai.Config{
			ApiKey: os.Getenv("OPENAI_API_KEY"),
		}

		return openai.NewClient(config)
	},
	"prompter": func(c *Container) any {
		return prompter.New()
	},
	"db": func(c *Container) any {
		user := os.Getenv("DB_USER")
		pass := os.Getenv("DB_PASSWORD")
		name := os.Getenv("DB_NAME")
		host := os.Getenv("DB_HOST")
		dbConn := db.Connect(host, user, pass, name)

		return dbConn
	},
	"store": func(c *Container) any {
		dbConn, ok := c.Get("db").(*sql.DB)
		if !ok {
			panic("dbConn factory failed")
		}

		return &db.Store{
			Conversation:   db.NewPostgresConversationStore(dbConn, "conversations"),
			Message:        db.NewPostgresMessagesStore(dbConn, "messages"),
			Cache:          db.NewPostgresCacheStore(dbConn, "cache"),
			Memory:         db.NewPostgresMemoriesStore(dbConn, "memories"),
			MemoryFragment: db.NewPostgresMemoriesFragmentsStore(dbConn, "memory_fragments"),
		}
	},
	"rag": func(c *Container) any {
		llm, ok := c.Get("llm").(*openai.Client)
		if !ok {
			panic("openai factory failed")
		}
		store, ok := c.Get("store").(*db.Store)
		if !ok {
			panic("store factory failed")
		}
		prompter, ok := c.Get("prompter").(*prompter.Prompter)
		if !ok {
			panic("store factory failed")
		}
		vectordb, ok := c.Get("vectordb").(*vectordb.QdrantClient)
		if !ok {
			panic("vectordb factory failed")
		}

		return rag.New(llm, store, prompter, vectordb)
	},
	"vectordb": func(c *Container) any {
		user := os.Getenv("QDRANT_BASEURL")
		apiKey := os.Getenv("QDRANT_API_KEY")

		return vectordb.NewQdrantClient(user, apiKey)
	},
	"tts": func(c *Container) any {
		apiKey := os.Getenv("ELEVENLABS_API_KEY")
		config := text_to_speech.NewElevenLabsConfig().WithApiKey(apiKey)

		return text_to_speech.NewElevenLabsTTS(config)
	},
	"scraper": func(c *Container) any {
		return scraper.NewCollyScraper(colly.NewCollector())
	},
}
