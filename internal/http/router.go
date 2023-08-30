package http

import (
	"github.com/gin-gonic/gin"
	"github.com/jonada182/cover-letter-ai-api/types"
)

func SetupRouter(httpClient *HttpClient, storeClient types.StoreClient, openaiClient types.OpenAIClient) *gin.Engine {
	router := gin.Default()
	router.GET("/", httpClient.HandleIndex)
	router.POST("/cover-letter", func(ctx *gin.Context) {
		httpClient.HandleCoverLetter(ctx, storeClient, openaiClient)
	})
	router.POST("/career-profile", httpClient.HandleCreateCareerProfile)
	router.GET("/career-profile/:email", httpClient.HandleGetCareerProfile)
	return router
}
