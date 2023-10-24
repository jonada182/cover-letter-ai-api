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
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jonada182/cover-letter-ai-api/types"

	"github.com/gin-gonic/gin"
)

var OpenAICompletionsUrl = "https://api.openai.com/v1/chat/completions"
var (
	GPT35 = "gpt-3.5-turbo"
	GPT4 = "gpt-4"
)

type OpenAIClient struct {
	apiKey string
	model  string
}

type OpenAI interface {
	GenerateChatGPTCoverLetter(c *gin.Context, profileId uuid.UUID, jobPosting *types.JobPosting, s types.StoreClient) (string, int, error)
	GetCareerProfileInfoPrompt(profileId uuid.UUID, s types.StoreClient) (string, *types.CareerProfile, error)
	ParseCoverLetter(coverLetter *string, careerProfile *types.CareerProfile, jobPosting *types.JobPosting) (string, error)
}

// NewOpenAIClient initializes an OpenAI client with the API key from the .env file.
func NewOpenAIClient() (*OpenAIClient, error) {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		return nil, errors.New("no OpenAI API key present in env file")
	}
	return &OpenAIClient{
		apiKey: apiKey,
		model:  GPT35,
	}, nil
}

// GenerateChatGPTCoverLetter uses the OpenAI completions API to generate a cover letter using the given parameters
func (oa *OpenAIClient) GenerateChatGPTCoverLetter(c *gin.Context, profileId uuid.UUID, jobPosting *types.JobPosting, s types.StoreClient) (string, int, error) {
	promptMessages := []types.ChatGTPRequestMessage{
		{
			Role:    "system",
			Content: "You write cover letters when I give you job details. Limit: 3 paragraphs, 300 words. Start with: Dear [Employer's Name],",
		},
	}

	// Add career profile information to prompt
	careerProfilePrompt, careerProfile, err := oa.GetCareerProfileInfoPrompt(profileId, s)
	if err != nil {
		return "", http.StatusInternalServerError, err
	}

	// Create cover letter prompt using the jobPosting data
	emptyLinesPattern := `\s*\n`
	re := regexp.MustCompile(emptyLinesPattern)
	jobPostingDetails := re.ReplaceAllString(jobPosting.Details, "\n")
	promptFormat := "Company:%s\nJob Role:%s\nDetails:\n%s\nSkills:%s"
	prompt := fmt.Sprintf(promptFormat, jobPosting.CompanyName, jobPosting.JobRole, jobPostingDetails, jobPosting.Skills)
	coverLetterPrompt := fmt.Sprintf("Write a cover letter for this job:\n%s", prompt)
	if careerProfilePrompt != "" {
		coverLetterPrompt = fmt.Sprintf("%s\n\n%s", coverLetterPrompt, careerProfilePrompt)
	}

	promptMessages = append(promptMessages, types.ChatGTPRequestMessage{
		Role:    "user",
		Content: coverLetterPrompt,
	})
	fmt.Printf("OpenAI prompt:\n%s\n", coverLetterPrompt)

	requestBody := &types.ChatGPTRequest{
		Model:       oa.model,
		Messages:    promptMessages,
		Temperature: 0.5,
		MaxTokens:   512,
	}
	requestBodyBytes, err := json.Marshal(requestBody)
	if err != nil {
		return "", http.StatusInternalServerError, err
	}

	client := &http.Client{}
	// Make a request to the OpenAI completions API using the defined model and messages (prompts)
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
	coverLetter := responseData.Choices[0].Message.Content
	coverLetter, err = oa.ParseCoverLetter(&coverLetter, careerProfile, jobPosting)
	if err != nil {
		return "", http.StatusInternalServerError, err
	}

	// Return the content for the last received message from OpenAI if the response was successful
	return coverLetter, http.StatusOK, nil
}

// GetCareerProfileInfoPrompt returns a prompt string with the CareerProfile data retrieved using the given email
func (oa *OpenAIClient) GetCareerProfileInfoPrompt(profileId uuid.UUID, s types.StoreClient) (string, *types.CareerProfile, error) {
	info := ""

	careerProfile, err := s.GetCareerProfileByID(profileId)
	if err != nil {
		return "", &types.CareerProfile{}, err
	}

	// Concatenate each of the career profile fields into a single message
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

	return info, careerProfile, nil
}

// ParseCoverLetter adds contact information to the cover letter from the CareerProfile
func (oa *OpenAIClient) ParseCoverLetter(coverLetter *string, careerProfile *types.CareerProfile, jobPosting *types.JobPosting) (string, error) {
	if *coverLetter == "" {
		return "", errors.New("received an empty cover letter from OpenAI")
	}

	template := "[Your Name]\n[Your Address]\n[City, State, ZIP]\n[Email Address]\n[Phone Number]\n[Date]\n\n[Employer's Name]\n[Company Name]\n[Company Address]\n[City, State, ZIP]\n\n"
	parsedLetter := template + *coverLetter

	year, month, day := time.Now().Date()
	coverLetterKeywords := map[string]string{
		"[Your Name]":             fmt.Sprintf("%s %s", careerProfile.FirstName, careerProfile.LastName),
		"[Your Address]":          careerProfile.ContactInfo.Address,
		"[Email Address]":         careerProfile.ContactInfo.Email,
		"[Phone Number]":          careerProfile.ContactInfo.Phone,
		"[Date]":                  fmt.Sprintf("%s %d, %d", month.String(), day, year),
		"[Today's Date]":          fmt.Sprintf("%s %d, %d", month.String(), day, year),
		"[Employer's Name]":       "Hiring Manager",
		"[Recipient's Name]":      "Hiring Manager",
		"[Company Name]":          jobPosting.CompanyName,
		"[Company Address]":       "",
		"[City, State, ZIP]":      "",
		"[City, State, ZIP Code]": "",
	}

	for keyword, value := range coverLetterKeywords {
		parsedLetter = strings.Replace(parsedLetter, keyword, value, -1)
	}

	return parsedLetter, nil
}
