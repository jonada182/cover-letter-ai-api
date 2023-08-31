package handler

//go:generate mockgen -destination=../../mocks/mock_handler.go -package=mocks github.com/jonada182/cover-letter-ai-api/internal/handler HandlerInterface

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jonada182/cover-letter-ai-api/types"
)

type HandlerInterface interface {
	HandleIndex(c *gin.Context)
	HandleCoverLetter(c *gin.Context, s types.StoreClient, oa types.OpenAIClient)
	HandleCreateCareerProfile(c *gin.Context)
	HandleGetCareerProfile(c *gin.Context)
}

type Handler struct {
	StoreClient  types.StoreClient
	OpenAIClient types.OpenAIClient
}

// NewHandler Initializes application handler allowing the injection of clients
func NewHandler(s types.StoreClient, o types.OpenAIClient) *Handler {
	return &Handler{
		StoreClient:  s,
		OpenAIClient: o,
	}
}

// SetupRouter sets all the API endpoints and returns a gin router
func (h *Handler) SetupRouter() *gin.Engine {
	router := gin.Default()
	router.GET("/", h.HandleIndex)
	router.POST("/cover-letter", h.HandleCoverLetter)
	router.POST("/career-profile", h.HandleCreateCareerProfile)
	router.GET("/career-profile/:email", h.HandleGetCareerProfile)
	return router
}

// HandleIndex returns a welcome message when "/" is accessed
func (h *Handler) HandleIndex(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Welcome to the CoverLetterAI API"})
}

// HandleCoverLetter handles a POST method that returns a cover letter from OpenAI
func (h *Handler) HandleCoverLetter(c *gin.Context) {
	// Receive CoverLetterRequest parameters from request payload
	var coverLetterRequest types.CoverLetterRequest
	if err := c.ShouldBindJSON(&coverLetterRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if coverLetterRequest.Email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "email is required"})
		return
	}

	jobPosting := coverLetterRequest.JobPosting
	if jobPosting.CompanyName == "" || jobPosting.JobRole == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "company name and job role are required"})
		return
	}

	// Create cover letter prompt using the jobPosting data
	promptFormat := "Company:%s\nJob Role:%s\nDetails:%s\nSkills:%s"
	prompt := fmt.Sprintf(promptFormat, jobPosting.CompanyName, jobPosting.JobRole, jobPosting.Details, jobPosting.Skills)

	// Call OpenAI to generate a cover letter with the given parameters
	coverLetter, statusCode, err := h.OpenAIClient.GenerateChatGPTCoverLetter(c, coverLetterRequest.Email, prompt, h.StoreClient)
	if err != nil {
		c.JSON(statusCode, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": coverLetter})
}

// HandleCreateCareerProfile handles a POST method to create a career profile in MongoDB
func (h *Handler) HandleCreateCareerProfile(c *gin.Context) {
	// Receive CareerProfileRequest parameters from request payload
	var careerProfileRequest types.CareerProfileRequest
	if err := c.ShouldBindJSON(&careerProfileRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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
	email := c.Param("email")
	if email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no email provided in the request"})
		return
	}

	// Call store method to retrieve CareerProfile from MongoDB
	careerProfile, err := h.StoreClient.GetCareerProfile(email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": &careerProfile})
}
