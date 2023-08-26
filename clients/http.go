package clients

import (
	"fmt"
	"net/http"
	"resu-mate-api/types"

	"github.com/gin-gonic/gin"
)

type HttpClient struct{}

func NewHttpClient() (httpClient *HttpClient) {
	return &HttpClient{}
}

func (client *HttpClient) GenerateCoverLetter(c *gin.Context) {
	var jobPosting types.JobPosting
	if err := c.ShouldBindJSON(&jobPosting); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

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
	openAIClient.HandleChatGPT(c, prompt)
}
