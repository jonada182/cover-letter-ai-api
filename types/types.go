package types

import "github.com/google/uuid"

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
	FirstName       string             `json:"first_name"`
	LastName        string             `json:"last_name"`
	Headline        string             `json:"headline"`
	ExperienceYears int                `json:"experience_years"`
	Summary         *string            `json:"summary"`
	Skills          *[]string          `json:"skills"`
	ContactInfo     *CareerContactInfo `json:"contact_info"`
}

type CareerProfile struct {
	ID              uuid.UUID          `bson:"id"`
	FirstName       string             `bson:"first_name"`
	LastName        string             `bson:"last_name"`
	Headline        string             `bson:"headline"`
	ExperienceYears int                `bson:"experience_years"`
	Summary         *string            `bson:"summary"`
	Skills          *[]string          `bson:"skills"`
	ContactInfo     *CareerContactInfo `bson:"contact_info"`
}

type CareerContactInfo struct {
	Email   string `bson:"email"`
	Address string `bson:"address"`
	Phone   string `bson:"phone"`
	Website string `bson:"website"`
}
