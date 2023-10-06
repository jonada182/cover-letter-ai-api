package main

import (
	"log"

	"github.com/jonada182/cover-letter-ai-api/internal/handler"
	"github.com/jonada182/cover-letter-ai-api/internal/openai"
	"github.com/jonada182/cover-letter-ai-api/internal/store"
	"github.com/jonada182/cover-letter-ai-api/util"
)

func main() {
	err := util.LoadEnvFile(".env")
	if err != nil {
		log.Println("Error loading .env file:", err)
	}

	storeClient, err := store.NewStore()
	if err != nil {
		log.Fatal("Error initializing store:", err)
	}

	openAIClient, err := openai.NewOpenAIClient()
	if err != nil {
		log.Fatal("Error initializing OpenAI client:", err)
	}

	h := handler.NewHandler(storeClient, openAIClient)

	r := h.SetupRouter()
	r.Run(":8080")
}
