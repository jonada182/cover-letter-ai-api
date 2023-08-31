package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/jonada182/cover-letter-ai-api/mocks"
	"github.com/jonada182/cover-letter-ai-api/types"
	"github.com/jonada182/cover-letter-ai-api/util"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestHandler(t *testing.T) {
	t.Run("HandleIndex", func(t *testing.T) {
		// Setup httptest and gin router
		router, recorder := util.SetupTestRouter()
		handler := &Handler{}
		router.GET("/", handler.HandleIndex)

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
		apiEndpoint := "/cover-letter"
		t.Run("invalid request", func(t *testing.T) {
			// Setup httptest, gin router and environment variables
			router, recorder := util.SetupTestRouter()
			util.SetupTestEnvironment(t)

			// Setup mocks and expectations
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockStore := mocks.NewMockStore(ctrl)
			mockOpenAI := mocks.NewMockOpenAI(ctrl)

			// Setup request handler
			handler := NewHandler(mockStore, mockOpenAI)
			router.POST(apiEndpoint, handler.HandleCoverLetter)
			// Create a new HTTP request with no payload
			req, err := http.NewRequest(http.MethodPost, apiEndpoint, nil)
			assert.NoError(t, err)

			// Serve the request
			router.ServeHTTP(recorder, req)
			// Check the response status code
			assert.Equal(t, http.StatusBadRequest, recorder.Code)
			// Check the response body
			expectedResponse := `{"error":"error retrieving JSON: invalid request"}`
			assert.Equal(t, expectedResponse, recorder.Body.String())
		})

		t.Run("valid request", func(t *testing.T) {
			// Setup httptest, gin router and environment variables
			router, recorder := util.SetupTestRouter()
			util.SetupTestEnvironment(t)

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

			_, message, err := util.SetupTestCareerProfile(email)
			assert.NoError(t, err)
			assert.Contains(t, message, "career profile")

			// Setup mocks and expectations
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockStore := mocks.NewMockStore(ctrl)
			mockOpenAI := mocks.NewMockOpenAI(ctrl)
			mockOpenAI.EXPECT().
				GenerateChatGPTCoverLetter(gomock.Any(), gomock.Eq(email), gomock.Eq(&requestData.JobPosting), gomock.Any()).
				Return("perfect cover letter", 200, nil).
				Times(1)

			// Setup mocks and expectations
			handler := NewHandler(mockStore, mockOpenAI)
			router.POST(apiEndpoint, handler.HandleCoverLetter)

			// Create a new HTTP request with valid payload
			requestBody, err := json.Marshal(requestData)
			assert.NoError(t, err)
			req, err := http.NewRequest(http.MethodPost, apiEndpoint, bytes.NewBuffer(requestBody))
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

	t.Run("HandleCreateCareerProfile", func(t *testing.T) {
		apiEndpoint := "/career-profile"
		t.Run("invalid request", func(t *testing.T) {
			// Setup httptest, gin router and environment variables
			router, recorder := util.SetupTestRouter()
			util.SetupTestEnvironment(t)

			// Setup mocks and expectations
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockStore := mocks.NewMockStore(ctrl)
			mockOpenAI := mocks.NewMockOpenAI(ctrl)

			// Setup mocks and expectations
			handler := NewHandler(mockStore, mockOpenAI)
			router.POST(apiEndpoint, handler.HandleCreateCareerProfile)
			// Create a new HTTP request with no payload
			req, err := http.NewRequest(http.MethodPost, apiEndpoint, nil)
			assert.NoError(t, err)

			// Serve the request
			router.ServeHTTP(recorder, req)
			// Check the response status code
			assert.Equal(t, http.StatusBadRequest, recorder.Code)
			// Check the response body
			expectedResponse := `{"error":"error retrieving JSON: invalid request"}`
			assert.Equal(t, expectedResponse, recorder.Body.String())
		})

		t.Run("valid request", func(t *testing.T) {
			// Setup httptest, gin router and environment variables
			router, recorder := util.SetupTestRouter()
			util.SetupTestEnvironment(t)

			email := "test@email"
			requestData := types.CareerProfileRequest{
				FirstName:       "John",
				LastName:        "Doe",
				Headline:        "Manager",
				ExperienceYears: 5,
				ContactInfo: &types.ContactInfo{
					Email: email,
				},
			}
			expectedResult := &types.CareerProfile{
				ID:              uuid.New(),
				FirstName:       "John",
				LastName:        "Doe",
				Headline:        "Manager",
				ExperienceYears: 5,
				ContactInfo: &types.ContactInfo{
					Email: email,
				},
			}
			expectedData, err := json.Marshal(expectedResult)
			assert.NoError(t, err)

			// Setup mocks and expectations
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockStore := mocks.NewMockStore(ctrl)
			mockOpenAI := mocks.NewMockOpenAI(ctrl)
			mockStore.
				EXPECT().
				StoreCareerProfile(gomock.Eq(&requestData)).
				Return(expectedResult, "success", nil).
				Times(1)

			// Setup mocks and expectations
			handler := NewHandler(mockStore, mockOpenAI)
			router.POST(apiEndpoint, handler.HandleCreateCareerProfile)

			// Create a new HTTP request with valid payload
			requestBody, err := json.Marshal(requestData)
			assert.NoError(t, err)
			req, err := http.NewRequest(http.MethodPost, apiEndpoint, bytes.NewBuffer(requestBody))
			assert.NoError(t, err)

			// Serve the request
			router.ServeHTTP(recorder, req)

			// Check the response status code
			assert.Equal(t, http.StatusOK, recorder.Code)

			// Check the response body
			expectedResponse := fmt.Sprintf("{\"data\":%s,\"message\":\"success\"}", string(expectedData))
			assert.Equal(t, expectedResponse, recorder.Body.String())
		})
	})

	t.Run("HandleGetCareerProfile", func(t *testing.T) {
		// Setup httptest, gin router and environment variables
		router, recorder := util.SetupTestRouter()
		util.SetupTestEnvironment(t)

		email := "test@email"
		expectedResult := &types.CareerProfile{
			ID:              uuid.New(),
			FirstName:       "John",
			LastName:        "Doe",
			Headline:        "Manager",
			ExperienceYears: 5,
			ContactInfo: &types.ContactInfo{
				Email: email,
			},
		}
		expectedData, err := json.Marshal(expectedResult)
		assert.NoError(t, err)

		// Setup mocks and expectations
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockStore := mocks.NewMockStore(ctrl)
		mockOpenAI := mocks.NewMockOpenAI(ctrl)
		mockStore.
			EXPECT().
			GetCareerProfile(gomock.Eq(email)).
			Return(expectedResult, nil).
			Times(1)

		// Setup mocks and expectations
		handler := NewHandler(mockStore, mockOpenAI)
		router.GET("/career-profile/:email", handler.HandleGetCareerProfile)

		req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/career-profile/%s", email), nil)
		assert.NoError(t, err)

		// Serve the request
		router.ServeHTTP(recorder, req)

		// Check the response status code
		assert.Equal(t, http.StatusOK, recorder.Code)

		// Check the response body
		expectedResponse := fmt.Sprintf("{\"data\":%s}", string(expectedData))
		assert.Equal(t, expectedResponse, recorder.Body.String())
	})
}
