package main

import (
	"cover-letter-ai-api/internal/http"
	"cover-letter-ai-api/util"
	"log"

	"github.com/gin-gonic/gin"
)

func setupRouter() *gin.Engine {
	router := gin.Default()
	httpClient := http.NewHttpClient()
	router.GET("/", httpClient.HandleIndex)
	router.POST("/cover-letter", httpClient.HandleCoverLetter)
	router.POST("/career-profile", httpClient.HandleCreateCareerProfile)
	router.GET("/career-profile/:email", httpClient.HandleGetCareerProfile)
	return router
}

func main() {
	err := util.LoadEnvFile(".env")
	if err != nil {
		log.Fatal("Error loading .env file:", err)
	}
	r := setupRouter()
	r.Run(":8080")
}
