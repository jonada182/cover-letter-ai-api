package clients

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"resu-mate-api/types"

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

func (oa *OpenAIClient) HandleChatGPT(c *gin.Context, prompt string) {
	apiKey := oa.apiKey
	apiUrl := "https://api.openai.com/v1/chat/completions"
	requestBody := &types.ChatGPTRequest{
		Model: oa.model,
		Messages: []types.ChatGTPRequestMessage{
			{
				Role:    "system",
				Content: "You are professional career advisor.",
			},
			{
				Role:    "user",
				Content: "I need a cover letter for job. 3 paragraphs, 300 words. ONLY letter body. Details below:",
			},
			{
				Role:    "user",
				Content: prompt,
			},
		},
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
