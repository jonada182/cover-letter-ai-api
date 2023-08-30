package handler

import (
	"github.com/jonada182/cover-letter-ai-api/internal/http"
	"github.com/jonada182/cover-letter-ai-api/types"
)

type Handler struct {
	StoreClient  types.StoreClient
	HttpClient   *http.HttpClient
	OpenAIClient types.OpenAIClient
}

func NewHandler(s types.StoreClient, h *http.HttpClient, o types.OpenAIClient) *Handler {
	return &Handler{
		StoreClient:  s,
		HttpClient:   h,
		OpenAIClient: o,
	}
}
