package handler

import (
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (h *Handler) middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With, UserID")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		if c.Request.URL.Path != "/linkedin/callback" {
			authorizationHeader := c.GetHeader("Authorization")
			tokenParts := strings.Split(authorizationHeader, " ")
			if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized request!"})
				return
			}
			accessToken := tokenParts[1]

			if c.Request.URL.Path != "/user" {
				UserID := c.GetHeader("UserID")
				profileId, err := uuid.Parse(UserID)
				if err != nil {
					log.Printf("error when parsing UserID: %s", err.Error())
					c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized request"})
					return
				}
				validToken, err := h.StoreClient.ValidateAccessToken(profileId, accessToken)
				if !validToken || err != nil {
					log.Printf("error when validating access token: %s", err.Error())
					c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized request"})
					return
				}
				c.Set("ProfileID", profileId)
			}
			c.Set("AccessToken", accessToken)
		}

		c.Next()
	}
}
