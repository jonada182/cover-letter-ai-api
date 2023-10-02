package types

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/mongo"
)

type CoverLetterRequest struct {
	ProfileID  uuid.UUID  `json:"profile_id"`
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
	Model       string                  `json:"model"`
	Messages    []ChatGTPRequestMessage `json:"messages"`
	Temperature float32                 `json:"temperature"`
	MaxTokens   int                     `json:"max_tokens"`
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

type LinkedInTokenResponse struct {
	AccessToken           string  `json:"access_token"`
	ExpiresIn             float64 `json:"expires_in"`
	RefreshToken          string  `json:"refresh_token"`
	RefreshTokenExpiresIn float64 `json:"refresh_token_expires_in"`
	Scope                 string  `json:"scope"`
}

func MapToLinkedInTokenResponse(data map[string]interface{}) LinkedInTokenResponse {
	var mappedData LinkedInTokenResponse
	if accessToken, ok := data["access_token"].(string); ok {
		mappedData.AccessToken = accessToken
	}
	if expiresIn, ok := data["expires_in"].(float64); ok {
		mappedData.ExpiresIn = expiresIn
	}
	if refreshToken, ok := data["refresh_token"].(string); ok {
		mappedData.RefreshToken = refreshToken
	}
	if refreshTokenExpiresIn, ok := data["refresh_token_expires_in"].(float64); ok {
		mappedData.RefreshTokenExpiresIn = refreshTokenExpiresIn
	}
	if scope, ok := data["scope"].(string); ok {
		mappedData.Scope = scope
	}
	return mappedData
}

type LinkedInUserData struct {
	Sub        string `json:"sub"`
	Name       string `json:"name"`
	GivenName  string `json:"given_name"`
	FamilyName string `json:"family_name"`
	Picture    string `json:"picture"`
	Email      string `json:"email"`
}

type AccessToken struct {
	ProfileID   uuid.UUID `bson:"profile_id" json:"profile_id"`
	AccessToken string    `bson:"access_token" json:"access_token"`
	ExpiresAt   string    `bson:"expires_at" json:"expires_at"`
}

func MapToLinkedInUserData(data map[string]interface{}) LinkedInUserData {
	var mappedData LinkedInUserData
	if sub, ok := data["sub"].(string); ok {
		mappedData.Sub = sub
	}
	if name, ok := data["name"].(string); ok {
		mappedData.Name = name
	}
	if givenName, ok := data["given_name"].(string); ok {
		mappedData.GivenName = givenName
	}
	if familyName, ok := data["family_name"].(string); ok {
		mappedData.FamilyName = familyName
	}
	if picture, ok := data["picture"].(string); ok {
		mappedData.Picture = picture
	}
	if email, ok := data["email"].(string); ok {
		mappedData.Email = email
	}
	return mappedData
}

type Handler interface {
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

type StoreClient interface {
	Connect() (*mongo.Client, context.Context, error)
	Disconnect(ctx context.Context, client *mongo.Client)
	GetCareerProfileByEmail(email string) (*CareerProfile, error)
	GetCareerProfileByID(profileId uuid.UUID) (*CareerProfile, error)
	StoreCareerProfile(careerProfileRequest *CareerProfile) (*CareerProfile, string, error)
	GetJobApplications(profileId uuid.UUID) (*[]JobApplication, error)
	StoreJobApplication(jobApplicationRequest *JobApplication) (*JobApplication, string, error)
	DeleteJobApplication(jobApplicationId uuid.UUID) error
	StoreAccessToken(profileId uuid.UUID, accessToken string) (string, error)
	ValidateAccessToken(profileId uuid.UUID, accessToken string) (bool, error)
}

type OpenAIClient interface {
	GenerateChatGPTCoverLetter(c *gin.Context, profileId uuid.UUID, jobPosting *JobPosting, s StoreClient) (string, int, error)
	GetCareerProfileInfoPrompt(profileId uuid.UUID, s StoreClient) (string, *CareerProfile, error)
	ParseCoverLetter(coverLetter *string, careerProfile *CareerProfile, jobPosting *JobPosting) (string, error)
}
