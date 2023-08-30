package http

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/jonada182/cover-letter-ai-api/internal/store"
	"github.com/jonada182/cover-letter-ai-api/mocks"
	"github.com/jonada182/cover-letter-ai-api/types"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func SetupTest() (router *gin.Engine, recorder *httptest.ResponseRecorder) {
	// Initialize a new Gin router
	router = gin.Default()
	// Create a response recorder to capture the response
	recorder = httptest.NewRecorder()
	return router, recorder
}

func SetEnvironment(t *testing.T) {
	t.Setenv("OPENAI_API_KEY", "some_key")
	t.Setenv("MONGODB_URI", "mongodb://localhost:27018")
}

func TestHttpClient(t *testing.T) {
	t.Run("HandleIndex", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		client := &HttpClient{}
		router, recorder := SetupTest()
		router.GET("/", client.HandleIndex)

		// Create a new HTTP request to test the handler
		req, err := http.NewRequest(http.MethodGet, "/", nil)
		assert.NoError(t, err)

		// Serve the request
		router.ServeHTTP(recorder, req)

		// Check the response status code
		assert.Equal(t, http.StatusOK, recorder.Code)

		// Check the response body
		expectedResponse := `{"message":"Welcome to the CoverLetterAI API"}`
		assert.Equal(t, expectedResponse, recorder.Body.String())
	})

	t.Run("HandleCoverLetter", func(t *testing.T) {
		t.Run("invalid request", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			SetEnvironment(t)
			client := &HttpClient{}
			mockStore := mocks.NewMockStore(ctrl)
			mockOpenAI := mocks.NewMockOpenAI(ctrl)
			router, recorder := SetupTest()
			router.POST("/cover-letter", func(ctx *gin.Context) {
				client.HandleCoverLetter(ctx, mockStore, mockOpenAI)
			})
			// Create a new HTTP request with no payload
			req, err := http.NewRequest(http.MethodPost, "/cover-letter", nil)
			assert.NoError(t, err)

			// Serve the request
			router.ServeHTTP(recorder, req)
			// Check the response status code
			assert.Equal(t, http.StatusBadRequest, recorder.Code)
			// Check the response body
			expectedResponse := `{"error":"invalid request"}`
			assert.Equal(t, expectedResponse, recorder.Body.String())
		})

		t.Run("valid request", func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			SetEnvironment(t)
			client := &HttpClient{}
			email := "test@email"
			requestData := types.CoverLetterRequest{
				Email: email,
				JobPosting: types.JobPosting{
					CompanyName: "Acme",
					JobRole:     "Manager",
					Details:     "Great worker",
					Skills:      "Management",
				},
			}
			s, err := store.NewStore()
			assert.NoError(t, err)
			_, message, err := s.StoreCareerProfile(&types.CareerProfileRequest{
				FirstName:       "John",
				LastName:        "Doe",
				Headline:        "Manager",
				ExperienceYears: 5,
				ContactInfo: &types.ContactInfo{
					Email: email,
				},
			})
			assert.NoError(t, err)
			assert.Contains(t, message, "career profile")

			mockStore := mocks.NewMockStore(ctrl)
			mockOpenAI := mocks.NewMockOpenAI(ctrl)
			mockOpenAI.EXPECT().
				GenerateChatGPTCoverLetter(gomock.Any(), gomock.Eq(email), gomock.Any(), gomock.Any()).
				Return("perfect cover letter", 200, nil).
				Times(1)

			router, recorder := SetupTest()
			router.POST("/cover-letter", func(ctx *gin.Context) {
				client.HandleCoverLetter(ctx, mockStore, mockOpenAI)
			})

			// Create a new HTTP request with valid payload
			requestBody, err := json.Marshal(requestData)
			assert.NoError(t, err)
			req, err := http.NewRequest(http.MethodPost, "/cover-letter", bytes.NewBuffer(requestBody))
			assert.NoError(t, err)

			// Serve the request
			router.ServeHTTP(recorder, req)

			// Check the response status code
			assert.Equal(t, http.StatusOK, recorder.Code)

			// Check the response body
			expectedResponse := `{"data":"perfect cover letter"}`
			assert.Equal(t, expectedResponse, recorder.Body.String())
		})
	})
}
