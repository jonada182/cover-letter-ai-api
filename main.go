package main

import (
	"log"
	"resu-mate-api/handlers"
	"resu-mate-api/utils"

	"github.com/gin-gonic/gin"
)

func setupRouter() *gin.Engine {
	router := gin.Default()
	router.POST("/cover-letter", handlers.GenerateCoverLetter)
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
