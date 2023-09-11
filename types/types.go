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
	JobRole     string `json:"job_role" bind:"required"`
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
	ExperienceYears uint         `json:"experience_years"`
	Summary         *string      `json:"summary"`
	Skills          *[]string    `json:"skills"`
	ContactInfo     *ContactInfo `json:"contact_info"`
}

type CareerProfile struct {
	ID              uuid.UUID    `bson:"id" json:"id"`
	FirstName       string       `bson:"first_name" json:"first_name"`
	LastName        string       `bson:"last_name" json:"last_name"`
	Headline        string       `bson:"headline" json:"headline"`
	ExperienceYears uint         `bson:"experience_years" json:"experience_years"`
	Summary         *string      `bson:"summary" json:"summary"`
	Skills          *[]string    `bson:"skills" json:"skills"`
	ContactInfo     *ContactInfo `bson:"contact_info" json:"contact_info"`
}

type ContactInfo struct {
	Email   string `bson:"email" json:"email"`
	Address string `bson:"address" json:"address"`
	Phone   string `bson:"phone" json:"phone"`
	Website string `bson:"website" json:"website"`
}

type JobApplication struct {
	ID          uuid.UUID              `bson:"id" json:"id"`
	ProfileID   uuid.UUID              `bson:"profile_id" json:"profile_id"`
	CompanyName string                 `bson:"company_name" json:"company_name" bind:"required"`
	JobRole     string                 `bson:"job_role" json:"job_role" bind:"required"`
	URL         *string                `bson:"url" json:"url"`
	Events      *[]JobApplicationEvent `bson:"events" json:"events"`
	CreatedAt   *string                `bson:"created_at" json:"created_at"`
	UpdatedAt   *string                `bson:"updated_at" json:"updated_at"`
}

type JobApplicationEvent struct {
	Type            uint    `bson:"type" json:"type"`
	Description     string  `bson:"description" json:"description"`
	Date            string  `bson:"date" json:"date"`
	AdditionalNotes *string `bson:"additional_notes" json:"additional_notes"`
}

type Handler interface {
	HandleIndex(c *gin.Context)
	HandleCoverLetter(c *gin.Context)
	HandleCreateCareerProfile(c *gin.Context)
	HandleGetCareerProfile(c *gin.Context)
	HandleCreateJobApplication(c *gin.Context)
	HandleGetJobApplications(c *gin.Context)
}

type StoreClient interface {
	Connect() (*mongo.Client, context.Context, error)
	Disconnect(ctx context.Context, client *mongo.Client)
	StoreCareerProfile(careerProfileRequest *CareerProfileRequest) (*CareerProfile, string, error)
	GetCareerProfile(email string) (*CareerProfile, error)
	GetJobApplications(profileId uuid.UUID) (*[]JobApplication, error)
	StoreJobApplication(jobApplicationRequest *JobApplication) (*JobApplication, string, error)
}

type OpenAIClient interface {
	GenerateChatGPTCoverLetter(c *gin.Context, email string, jobPosting *JobPosting, s StoreClient) (string, int, error)
	GetCareerProfileInfoPrompt(email string, s StoreClient) (string, *CareerProfile, error)
	ParseCoverLetter(coverLetter *string, careerProfile *CareerProfile, jobPosting *JobPosting) (string, error)
}
