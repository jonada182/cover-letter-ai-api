package http

//go:generate mockgen -destination=../../mocks/mock_http.go -package=mocks github.com/jonada182/cover-letter-ai-api/internal/http Http

import (
	"fmt"
	"net/http"

	"github.com/jonada182/cover-letter-ai-api/internal/store"
	"github.com/jonada182/cover-letter-ai-api/types"

	"github.com/gin-gonic/gin"
)

type HttpClient struct{}

type Http interface {
	HandleIndex(c *gin.Context)
	HandleCoverLetter(c *gin.Context)
	HandleCreateCareerProfile(c *gin.Context)
	HandleGetCareerProfile(c *gin.Context)
}

func NewHttpClient() (*HttpClient, error) {
	return &HttpClient{}, nil
}

func (client *HttpClient) HandleIndex(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Welcome to the CoverLetterAI API"})
}

func (client *HttpClient) HandleCoverLetter(c *gin.Context, s types.StoreClient, oa types.OpenAIClient) {
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

	promptFormat := "Company:%s\nJob Role:%s\nDetails:%s\nSkills:%s"
	prompt := fmt.Sprintf(promptFormat, jobPosting.CompanyName, jobPosting.JobRole, jobPosting.Details, jobPosting.Skills)

	coverLetter, statusCode, err := oa.GenerateChatGPTCoverLetter(c, coverLetterRequest.Email, prompt, s)
	if err != nil {
		c.JSON(statusCode, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": coverLetter})
}

func (client *HttpClient) HandleCreateCareerProfile(c *gin.Context) {
	var careerProfileRequest types.CareerProfileRequest
	if err := c.ShouldBindJSON(&careerProfileRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if careerProfileRequest.Headline == "" || careerProfileRequest.ExperienceYears == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "headline and experience are required"})
		return
	}

	s, err := store.NewStore()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	careerProfile, responseMsq, err := s.StoreCareerProfile(&careerProfileRequest)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": responseMsq, "data": careerProfile})
}

func (client *HttpClient) HandleGetCareerProfile(c *gin.Context) {
	email := c.Param("email")
	if email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no email provided in the request"})
		return
	}

	s, err := store.NewStore()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	careerProfile, err := s.GetCareerProfile(email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": &careerProfile})
}
