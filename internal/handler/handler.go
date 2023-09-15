package handler

//go:generate mockgen -destination=../../mocks/mock_handler.go -package=mocks github.com/jonada182/cover-letter-ai-api/internal/handler HandlerInterface

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jonada182/cover-letter-ai-api/types"
)

type HandlerInterface interface {
	HandleIndex(c *gin.Context)
	HandleCoverLetter(c *gin.Context, s types.StoreClient, oa types.OpenAIClient)
	HandleCreateCareerProfile(c *gin.Context)
	HandleGetCareerProfile(c *gin.Context)
	HandleCreateJobApplication(c *gin.Context)
	HandleGetJobApplications(c *gin.Context)
	HandleLinkedInCallback(c *gin.Context)
	HandleGetUser(c *gin.Context)
}

type Handler struct {
	StoreClient  types.StoreClient
	OpenAIClient types.OpenAIClient
}

// NewHandler Initializes application handler allowing the injection of clients
func NewHandler(s types.StoreClient, o types.OpenAIClient) *Handler {
	return &Handler{
		StoreClient:  s,
		OpenAIClient: o,
	}
}

// SetupRouter sets all the API endpoints and returns a gin router
func (h *Handler) SetupRouter() *gin.Engine {
	router := gin.Default()
	router.Use(Middleware())
	router.GET("/", h.HandleIndex)
	router.POST("/cover-letter", h.HandleCoverLetter)
	router.POST("/career-profile", h.HandleCreateCareerProfile)
	router.GET("/career-profile/:profile_id", h.HandleGetCareerProfile)
	router.POST("/job-applications", h.HandleCreateJobApplication)
	router.GET("/job-applications/:profile_id", h.HandleGetJobApplications)
	router.GET("/linkedin/callback", h.HandleLinkedInCallback)
	router.GET("/user", h.HandleGetUser)
	return router
}

// HandleIndex returns a welcome message when "/" is accessed
func (h *Handler) HandleIndex(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Welcome to the CoverLetterAI API"})
}

// HandleCoverLetter handles a POST method that returns a cover letter from OpenAI
func (h *Handler) HandleCoverLetter(c *gin.Context) {
	// Receive CoverLetterRequest parameters from request payload
	var coverLetterRequest types.CoverLetterRequest
	if err := c.ShouldBindJSON(&coverLetterRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("error retrieving JSON: %s", err.Error())})
		return
	}

	jobPosting := coverLetterRequest.JobPosting
	if jobPosting.CompanyName == "" || jobPosting.JobRole == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "company name and job role are required"})
		return
	}

	// Call OpenAI to generate a cover letter with the given parameters
	coverLetter, statusCode, err := h.OpenAIClient.GenerateChatGPTCoverLetter(c, coverLetterRequest.ProfileID, &jobPosting, h.StoreClient)
	if err != nil {
		c.JSON(statusCode, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": coverLetter})
}

// HandleCreateCareerProfile handles a POST method to create a career profile in MongoDB
func (h *Handler) HandleCreateCareerProfile(c *gin.Context) {
	// Receive CareerProfileRequest parameters from request payload
	var careerProfileRequest types.CareerProfile
	if err := c.ShouldBindJSON(&careerProfileRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("error retrieving JSON: %s", err.Error())})
		return
	}

	if careerProfileRequest.Headline == "" || careerProfileRequest.ExperienceYears == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "headline and experience are required"})
		return
	}

	// Call store method to upsert CareerProfile in MongoDB
	careerProfile, responseMsq, err := h.StoreClient.StoreCareerProfile(&careerProfileRequest)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": responseMsq, "data": careerProfile})
}

// HandleCreateCareerProfile handles a GET method to retrieve a career profile from MongoDB
func (h *Handler) HandleGetCareerProfile(c *gin.Context) {
	profileIdParam := c.Param("profile_id")
	if profileIdParam == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no profile_id provided in the request"})
		return
	}

	profileId, err := uuid.Parse(profileIdParam)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Call store method to retrieve CareerProfile from MongoDB
	careerProfile, err := h.StoreClient.GetCareerProfileByID(profileId)
	if err != nil && strings.Contains(err.Error(), "no document") {
		c.JSON(http.StatusNotFound, gin.H{"error": "career profile not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": &careerProfile})
}

// HandleCreateJobApplication handles a POST method to create a job application in MongoDB
func (h *Handler) HandleCreateJobApplication(c *gin.Context) {
	// Receive JobApplication parameters from request payload
	var jobApplicationRequest types.JobApplication
	if err := c.ShouldBindJSON(&jobApplicationRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("error retrieving JSON: %s", err.Error())})
		return
	}

	if jobApplicationRequest.ProfileID.String() == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "profile id is required"})
		return
	}

	if jobApplicationRequest.CompanyName == "" || jobApplicationRequest.JobRole == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "company name and job role are required"})
		return
	}

	// Call store method to upsert Job Application in MongoDB
	jobApplication, responseMsq, err := h.StoreClient.StoreJobApplication(&jobApplicationRequest)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": responseMsq, "data": jobApplication})
}

// HandleGetJobApplications handles a GET method to retrieve job applications from MongoDB
func (h *Handler) HandleGetJobApplications(c *gin.Context) {
	profileIdParam := c.Param("profile_id")
	if profileIdParam == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no profile_id provided in the request"})
		return
	}

	profileId, err := uuid.Parse(profileIdParam)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Call store method to retrieve []JobApplication from MongoDB
	jobApplications, err := h.StoreClient.GetJobApplications(profileId)
	if err != nil && strings.Contains(err.Error(), "no document") {
		c.JSON(http.StatusNotFound, gin.H{"error": "no job applications found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": &jobApplications})
}

// HandleLinkedInCallback handles a oAuth callback from LinkedIn
func (h *Handler) HandleLinkedInCallback(c *gin.Context) {
	linkedInClientID := os.Getenv("LINKEDIN_CLIENT_ID")
	if linkedInClientID == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "no LinkedIn Client ID env variable"})
		return
	}

	linkedInClientSecret := os.Getenv("LINKEDIN_CLIENT_SECRET")
	if linkedInClientSecret == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "no LinkedIn Client Secret env variable"})
		return
	}

	baseUrl := os.Getenv("BASE_API_URL")
	if baseUrl == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "no base api url env variable"})
		return
	}

	state := c.Query("state")
	if state == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no state provided in the request"})
		return
	}

	code := c.Query("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no code provided in the request"})
		return
	}

	client := &http.Client{}

	// Set parameters for LinkedIn access token request
	data := url.Values{}
	data.Set("grant_type", "authorization_code")
	data.Set("code", code)
	data.Set("client_id", linkedInClientID)
	data.Set("client_secret", linkedInClientSecret)
	data.Set("redirect_uri", fmt.Sprintf("%s/linkedin/callback", baseUrl))

	// Create LinkedIn access token request
	tokenRequest, err := http.NewRequest("POST", "https://www.linkedin.com/oauth/v2/accessToken", strings.NewReader(data.Encode()))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	tokenRequest.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	tokenResponse, err := client.Do(tokenRequest)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer tokenResponse.Body.Close()
	// Get data from LinkedIn access token response
	tokenResponseBody, err := readResponse(tokenResponse)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	tokenResponseData := types.MapToLinkedInTokenResponse(tokenResponseBody)

	c.Redirect(http.StatusPermanentRedirect, fmt.Sprintf("http://localhost:3000/?access_token=%s", tokenResponseData.AccessToken))
}

func (h *Handler) HandleGetUser(c *gin.Context) {
	accessToken, exists := c.Get("AccessToken")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "no authorization token provided"})
	}

	client := &http.Client{}
	// Create LinkedIn user data request
	userDataRequest, err := http.NewRequest("GET", "https://api.linkedin.com/v2/userinfo", nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	userDataRequest.Header.Add("Authorization", fmt.Sprintf("Bearer %s", accessToken))
	userDataResponse, err := client.Do(userDataRequest)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer userDataResponse.Body.Close()
	// Get data from LinkedIn user data response
	userDataResponseBody, err := readResponse(userDataResponse)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	linkedInUserData := types.MapToLinkedInUserData(userDataResponseBody)

	var profileID uuid.UUID
	isNewProfile := false
	existingProfile, err := h.StoreClient.GetCareerProfileByEmail(linkedInUserData.Email)
	if err != nil && strings.Contains(err.Error(), "no document") {
		fmt.Println("creating new career profile from LinkedIn user data")
		newCareerProfile, _, err := h.StoreClient.StoreCareerProfile(&types.CareerProfile{
			FirstName: linkedInUserData.GivenName,
			LastName:  linkedInUserData.FamilyName,
			ContactInfo: &types.ContactInfo{
				Email: linkedInUserData.Email,
			},
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		isNewProfile = true
		profileID = newCareerProfile.ID
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if !isNewProfile {
		profileID = existingProfile.ID
	}

	c.JSON(http.StatusOK, gin.H{"profile_id": profileID})
}

func readResponse(response *http.Response) (map[string]interface{}, error) {
	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	if response.StatusCode != http.StatusOK {
		fmt.Println("error response body: ", string(responseBody), response.Request.URL)
		return nil, fmt.Errorf("HTTP request failed with status code:%d", response.StatusCode)
	}
	var data map[string]interface{}
	if err := json.Unmarshal(responseBody, &data); err != nil {
		return nil, err
	}
	return data, nil
}
