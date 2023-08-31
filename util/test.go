package util

import (
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/jonada182/cover-letter-ai-api/internal/store"
	"github.com/jonada182/cover-letter-ai-api/types"
)

// SetupTestRouter initializes a gin router and a httptest recorder to mock http calls
func SetupTestRouter() (router *gin.Engine, recorder *httptest.ResponseRecorder) {
	// Initialize a new Gin router
	router = gin.Default()
	// Create a response recorder to capture the response
	recorder = httptest.NewRecorder()
	return router, recorder
}

// SetupTestEnvironment sets fake environment variables used for testing.
// These should match the variables from .env
func SetupTestEnvironment(t *testing.T) {
	t.Setenv("OPENAI_API_KEY", "some_key")
	t.Setenv("MONGODB_URI", "mongodb://localhost:27018")
}

// SetupTestCareerProfile inserts a fake CareerProfile in MongoDB that can be used for testing
func SetupTestCareerProfile(email string) (careerProfile *types.CareerProfile, message string, err error) {
	s, err := store.NewStore()
	if err != nil {
		return &types.CareerProfile{}, "", err
	}
	careerProfile, message, err = s.StoreCareerProfile(&types.CareerProfileRequest{
		FirstName:       "John",
		LastName:        "Doe",
		Headline:        "Manager",
		ExperienceYears: 5,
		ContactInfo: &types.ContactInfo{
			Email: email,
		},
	})
	if err != nil {
		return &types.CareerProfile{}, "", err
	}
	return careerProfile, message, nil
}
