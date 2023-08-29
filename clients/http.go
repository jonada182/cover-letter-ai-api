package clients

import (
	"cover-letter-ai-api/types"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type HttpClient struct{}

func NewHttpClient() (httpClient *HttpClient) {
	return &HttpClient{}
}

func (client *HttpClient) HandleIndex(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Welcome to the CoverLetterAI API"})
}

func (client *HttpClient) HandleCoverLetter(c *gin.Context) {
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

	openAIClient, err := NewOpenAIClient()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	openAIClient.GenerateChatGPTCoverLetter(c, coverLetterRequest.Email, prompt)
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

	s, err := NewStore()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	_, err = s.StoreCareerProfile(&careerProfileRequest)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "career profile has been created"})
}

func (client *HttpClient) HandleGetCareerProfile(c *gin.Context) {
	email := c.Param("email")
	if email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no email provided in the request"})
		return
	}

	s, err := NewStore()
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
