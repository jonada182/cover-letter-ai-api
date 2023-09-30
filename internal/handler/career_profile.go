package handler

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jonada182/cover-letter-ai-api/types"
)

// HandleCreateCareerProfile handles a POST method to create a career profile in MongoDB
func (h *Handler) HandleCreateCareerProfile(c *gin.Context) {
	// Receive CareerProfileRequest parameters from request payload
	var careerProfileRequest types.CareerProfile
	if err := c.ShouldBindJSON(&careerProfileRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("error retrieving JSON: %s", err.Error())})
		return
	}

	if careerProfileRequest.Headline == "" || careerProfileRequest.ExperienceYears == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "headline and experience are required"})
		return
	}

	// Call store method to upsert CareerProfile in MongoDB
	careerProfile, responseMsq, err := h.StoreClient.StoreCareerProfile(&careerProfileRequest)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": responseMsq, "data": careerProfile})
}

// HandleCreateCareerProfile handles a GET method to retrieve a career profile from MongoDB
func (h *Handler) HandleGetCareerProfile(c *gin.Context) {
	profileIdParam := c.Param("id")
	if profileIdParam == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no profile_id provided in the request"})
		return
	}

	profileId, err := uuid.Parse(profileIdParam)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Call store method to retrieve CareerProfile from MongoDB
	careerProfile, err := h.StoreClient.GetCareerProfileByID(profileId)
	if err != nil && strings.Contains(err.Error(), "no document") {
		c.JSON(http.StatusNotFound, gin.H{"error": "career profile not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": &careerProfile})
}
