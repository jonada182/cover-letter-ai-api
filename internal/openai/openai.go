package openai

import (
	"bytes"
	"cover-letter-ai-api/internal/store"
	"cover-letter-ai-api/types"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

type OpenAIClient struct {
	apiKey string
	model  string
}

func NewOpenAIClient() (*OpenAIClient, error) {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		return nil, errors.New("no OpenAI API key present in env file")
	}
	return &OpenAIClient{
		apiKey: apiKey,
		model:  "gpt-3.5-turbo",
	}, nil
}

func (oa *OpenAIClient) GenerateChatGPTCoverLetter(c *gin.Context, email string, prompt string) {
	apiKey := oa.apiKey
	apiUrl := "https://api.openai.com/v1/chat/completions"
	promptMessages := []types.ChatGTPRequestMessage{
		{
			Role:    "system",
			Content: "You are professional career advisor.",
		},
		{
			Role:    "user",
			Content: "I need a cover letter for job. 3 paragraphs, 300 words. ONLY letter body. Details below:",
		},
	}

	// Add career profile information to prompt
	careerProfileInfo, err := generateCareerProfileInfoPrompt(email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if careerProfileInfo != "" {
		promptMessages = append(promptMessages, types.ChatGTPRequestMessage{
			Role:    "user",
			Content: careerProfileInfo,
		})
	}

	// Add cover letter details to prompt
	promptMessages = append(promptMessages, types.ChatGTPRequestMessage{
		Role:    "user",
		Content: prompt,
	})

	requestBody := &types.ChatGPTRequest{
		Model:    oa.model,
		Messages: promptMessages,
	}
	requestBodyBytes, err := json.Marshal(requestBody)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	client := &http.Client{}
	req, err := http.NewRequest("POST", apiUrl, bytes.NewReader(requestBodyBytes))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	// Send the request and handle the response
	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if resp.StatusCode != http.StatusOK {
		c.JSON(http.StatusInternalServerError, gin.H{"error": resp.Status})
		return
	}

	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}

	var responseData types.ChatGPTResponseData
	if err := json.Unmarshal(responseBody, &responseData); err != nil {
		fmt.Println("Error unmarshaling JSON:", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": responseData.Choices[0].Message.Content})
}

func generateCareerProfileInfoPrompt(email string) (string, error) {
	info := ""

	s, err := store.NewStore()
	if err != nil {
		return "", err
	}

	careerProfile, err := s.GetCareerProfile(email)
	if err != nil {
		return "", err
	}

	var builder strings.Builder
	builder.WriteString("Here is my career information:")
	if careerProfile.Headline != "" {
		builder.WriteString(fmt.Sprintf("\nHeadline:%s,", careerProfile.Headline))
	}
	if careerProfile.ExperienceYears > 0 {
		builder.WriteString(fmt.Sprintf("\nExperience:%d years,", careerProfile.ExperienceYears))
	}
	if careerProfile.Skills != nil && len(*careerProfile.Skills) > 0 {
		builder.WriteString(fmt.Sprintf("\nSkills:%s,", strings.Join(*careerProfile.Skills, ",")))
	}
	if careerProfile.Summary != nil && *careerProfile.Summary != "" {
		builder.WriteString(fmt.Sprintf("\nSummary:%s,", *careerProfile.Summary))
	}
	info = builder.String()

	return info, nil
}
