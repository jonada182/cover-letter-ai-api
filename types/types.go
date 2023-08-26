package types

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

type CareerProfile struct {
	ID              int32             `bson:"id"`
	Headline        string            `bson:"headline"`
	ExperienceYears int               `bson:"experience_years"`
	Summary         string            `bson:"summary"`
	Skills          []CareerSkills    `bson:"skills"`
	ContactInfo     CareerContactInfo `bson:"contact_info"`
}

type CareerSkills struct {
	Skill string `bson:"skill"`
}

type CareerContactInfo struct {
	Email   string `bson:"email"`
	Address string `bson:"address"`
	Phone   string `bson:"phone"`
	Website string `bson:"website"`
}
