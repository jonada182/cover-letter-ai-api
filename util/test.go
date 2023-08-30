package util

import (
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/jonada182/cover-letter-ai-api/internal/store"
	"github.com/jonada182/cover-letter-ai-api/types"
)

func SetupTestRouter() (router *gin.Engine, recorder *httptest.ResponseRecorder) {
	// Initialize a new Gin router
	router = gin.Default()
	// Create a response recorder to capture the response
	recorder = httptest.NewRecorder()
	return router, recorder
}

func SetupTestEnvironment(t *testing.T) {
	t.Setenv("OPENAI_API_KEY", "some_key")
	t.Setenv("MONGODB_URI", "mongodb://localhost:27018")
}

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
