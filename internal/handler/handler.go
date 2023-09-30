package handler

//go:generate mockgen -destination=../../mocks/mock_handler.go -package=mocks github.com/jonada182/cover-letter-ai-api/internal/handler HandlerInterface

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jonada182/cover-letter-ai-api/types"
)

type HandlerInterface interface {
	HandleIndex(c *gin.Context)
	HandleCoverLetter(c *gin.Context)
	HandleCreateCareerProfile(c *gin.Context)
	HandleGetCareerProfile(c *gin.Context)
	HandleCreateJobApplication(c *gin.Context)
	HandleGetJobApplications(c *gin.Context)
	HandleDeleteJobApplication(c *gin.Context)
	HandleLinkedInCallback(c *gin.Context)
	HandleGetUser(c *gin.Context)
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
	router.Use(Middleware())
	router.GET("/", h.HandleIndex)
	router.POST("/cover-letter", h.HandleCoverLetter)
	router.POST("/career-profile", h.HandleCreateCareerProfile)
	router.GET("/career-profile/:id", h.HandleGetCareerProfile)
	router.POST("/job-applications", h.HandleCreateJobApplication)
	router.GET("/job-applications/:profile_id", h.HandleGetJobApplications)
	router.DELETE("/job-applications/:id", h.HandleDeleteJobApplication)
	router.GET("/linkedin/callback", h.HandleLinkedInCallback)
	router.GET("/user", h.HandleGetUser)
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
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("error retrieving JSON: %s", err.Error())})
		return
	}

	jobPosting := coverLetterRequest.JobPosting
	if jobPosting.CompanyName == "" || jobPosting.JobRole == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "company name and job role are required"})
		return
	}

	// Call OpenAI to generate a cover letter with the given parameters
	coverLetter, statusCode, err := h.OpenAIClient.GenerateChatGPTCoverLetter(c, coverLetterRequest.ProfileID, &jobPosting, h.StoreClient)
	if err != nil {
		c.JSON(statusCode, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": coverLetter})
}

// readResponse returns a string map from a response body
func readResponse(response *http.Response) (map[string]interface{}, error) {
	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	if response.StatusCode != http.StatusOK {
		fmt.Println("error response body: ", string(responseBody), response.Request.URL)
		return nil, fmt.Errorf("HTTP request failed with status code:%d", response.StatusCode)
	}
	var data map[string]interface{}
	if err := json.Unmarshal(responseBody, &data); err != nil {
		return nil, err
	}
	return data, nil
}
