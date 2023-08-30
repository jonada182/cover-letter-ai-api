package main

import (
	"log"

	"github.com/jonada182/cover-letter-ai-api/internal/handler"
	"github.com/jonada182/cover-letter-ai-api/internal/http"
	"github.com/jonada182/cover-letter-ai-api/internal/openai"
	"github.com/jonada182/cover-letter-ai-api/internal/store"
	"github.com/jonada182/cover-letter-ai-api/util"
)

func main() {
	err := util.LoadEnvFile(".env")
	if err != nil {
		log.Fatal("Error loading .env file:", err)
	}

	storeClient, err := store.NewStore()
	if err != nil {
		log.Fatal("Error initializing store:", err)
	}

	httpClient, err := http.NewHttpClient()
	if err != nil {
		log.Fatal("Error initializing http client:", err)
	}

	openAIClient, err := openai.NewOpenAIClient()
	if err != nil {
		log.Fatal("Error initializing OpenAI client:", err)
	}

	h := handler.NewHandler(storeClient, httpClient, openAIClient)

	r := http.SetupRouter(h.HttpClient, h.StoreClient, h.OpenAIClient)
	r.Run(":8080")
}
