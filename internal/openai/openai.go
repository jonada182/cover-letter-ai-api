package openai

//go:generate mockgen -destination=../../mocks/mock_openai.go -package=mocks github.com/jonada182/cover-letter-ai-api/internal/openai OpenAI

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/jonada182/cover-letter-ai-api/types"

	"github.com/gin-gonic/gin"
)

var OpenAICompletionsUrl = "https://api.openai.com/v1/chat/completions"

type OpenAIClient struct {
	apiKey string
	model  string
}

type OpenAI interface {
	GenerateChatGPTCoverLetter(c *gin.Context, email string, prompt string, s types.StoreClient) (string, int, error)
	GetCareerProfileInfoPrompt(email string, s types.StoreClient) (string, error)
}

// NewOpenAIClient initializes an OpenAI client with the API key from the .env file.
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

// GenerateChatGPTCoverLetter uses the OpenAI completions API to generate a cover letter using the given parameters
func (oa *OpenAIClient) GenerateChatGPTCoverLetter(c *gin.Context, email string, prompt string, s types.StoreClient) (string, int, error) {
	promptMessages := []types.ChatGTPRequestMessage{
		{
			Role:    "system",
			Content: "You are professional career advisor.",
		},
	}

	// Add career profile information to prompt
	careerProfileInfo, err := oa.GetCareerProfileInfoPrompt(email, s)
	if err != nil {
		return "", http.StatusInternalServerError, err
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
		Content: fmt.Sprintf("I need the message body of a cover letter for job. 3 paragraphs, 300 words. ONLY letter body. Details below:%s", prompt),
	})

	requestBody := &types.ChatGPTRequest{
		Model:    oa.model,
		Messages: promptMessages,
	}
	requestBodyBytes, err := json.Marshal(requestBody)
	if err != nil {
		return "", http.StatusInternalServerError, err
	}

	client := &http.Client{}
	req, err := http.NewRequest("POST", OpenAICompletionsUrl, bytes.NewReader(requestBodyBytes))
	if err != nil {
		return "", http.StatusInternalServerError, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+oa.apiKey)

	// Send the request and handle the response
	resp, err := client.Do(req)
	if err != nil {
		return "", http.StatusInternalServerError, err
	}

	if resp.StatusCode != http.StatusOK {
		return "", http.StatusInternalServerError, err
	}

	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", http.StatusInternalServerError, err
	}

	var responseData types.ChatGPTResponseData
	if err := json.Unmarshal(responseBody, &responseData); err != nil {
		return "", http.StatusInternalServerError, err
	}

	return responseData.Choices[0].Message.Content, http.StatusOK, nil
}

// GetCareerProfileInfoPrompt returns a prompt string with the CareerProfile data retrieved using the given email
func (oa *OpenAIClient) GetCareerProfileInfoPrompt(email string, s types.StoreClient) (string, error) {
	info := ""

	careerProfile, err := s.GetCareerProfile(email)
	if err != nil {
		return "", err
	}

	var builder strings.Builder
	builder.WriteString("Here is my career profile information:")
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
