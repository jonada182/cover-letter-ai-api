package main

import (
	"log"
	"resu-mate-api/clients"
	"resu-mate-api/utils"

	"github.com/gin-gonic/gin"
)

func setupRouter() *gin.Engine {
	router := gin.Default()
	httpClient := clients.NewHttpClient()
	router.POST("/cover-letter", httpClient.GenerateCoverLetter)
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
