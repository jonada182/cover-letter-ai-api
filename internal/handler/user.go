package handler

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jonada182/cover-letter-ai-api/types"
)

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

// HandleGetUser returns a career profile with an email from LinkedIn using a token
func (h *Handler) HandleGetUser(c *gin.Context) {
	accessTokenParam, exists := c.Get("AccessToken")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "no authorization token provided"})
		return
	}
	var accessToken string
	if accessTokenParam != nil {
		if str, ok := accessTokenParam.(string); ok {
			accessToken = str
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "no authorization token provided"})
			return
		}
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

	// Store access token in DB
	_, err = h.StoreClient.StoreAccessToken(profileID, accessToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"profile_id": profileID})
}
