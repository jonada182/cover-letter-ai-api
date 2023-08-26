package main

import (
	"cover-letter-ai-api/clients"
	"cover-letter-ai-api/utils"
	"log"

	"github.com/gin-gonic/gin"
)

func setupRouter() *gin.Engine {
	router := gin.Default()
	httpClient := clients.NewHttpClient()
	router.GET("/", httpClient.HandleIndex)
	router.POST("/cover-letter", httpClient.HandleCoverLetter)
	router.POST("/career-profile", httpClient.HandleCreateCareerProfile)
	router.GET("/career-profile/:email", httpClient.HandleGetCareerProfile)
	return router
}

func main() {
	err := utils.LoadEnvFile(".env")
	if err != nil {
		log.Fatal("Error loading .env file:", err)
	}
	r := setupRouter()
	r.Run(":8080")
}
