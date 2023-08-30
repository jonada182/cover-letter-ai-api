package types

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/mongo"
)

type CoverLetterRequest struct {
	Email      string     `json:"email"`
	JobPosting JobPosting `json:"job_posting"`
}

type JobPosting struct {
	CompanyName string `json:"company_name" bind:"required"`
	JobRole     string `json:"job_role"`
	Details     string `json:"job_details"`
	Skills      string `json:"skills"`
}

type ChatGTPRequestMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatGPTRequest struct {
	Model    string                  `json:"model"`
	Messages []ChatGTPRequestMessage `json:"messages"`
}

type ChatGPTResponseChoice struct {
	Index        int            `json:"index"`
	Message      ChatGPTMessage `json:"message"`
	FinishReason string         `json:"finish_reason"`
}

type ChatGPTMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatGPTResponseData struct {
	Choices []ChatGPTResponseChoice `json:"choices"`
}

type CareerProfileRequest struct {
	FirstName       string       `json:"first_name"`
	LastName        string       `json:"last_name"`
	Headline        string       `json:"headline"`
	ExperienceYears int          `json:"experience_years"`
	Summary         *string      `json:"summary"`
	Skills          *[]string    `json:"skills"`
	ContactInfo     *ContactInfo `json:"contact_info"`
}

type CareerProfile struct {
	ID              uuid.UUID    `bson:"id"`
	FirstName       string       `bson:"first_name"`
	LastName        string       `bson:"last_name"`
	Headline        string       `bson:"headline"`
	ExperienceYears int          `bson:"experience_years"`
	Summary         *string      `bson:"summary"`
	Skills          *[]string    `bson:"skills"`
	ContactInfo     *ContactInfo `bson:"contact_info"`
}

type ContactInfo struct {
	Email   string `bson:"email"`
	Address string `bson:"address"`
	Phone   string `bson:"phone"`
	Website string `bson:"website"`
}

type StoreClient interface {
	Connect() (*mongo.Client, context.Context, error)
	Disconnect(ctx context.Context, client *mongo.Client)
	StoreCareerProfile(careerProfileRequest *CareerProfileRequest) (*CareerProfile, string, error)
	GetCareerProfile(email string) (*CareerProfile, error)
}

type OpenAIClient interface {
	GenerateChatGPTCoverLetter(c *gin.Context, email string, prompt string, s StoreClient) (string, int, error)
	GetCareerProfileInfoPrompt(email string, s StoreClient) (string, error)
}
