package store

//go:generate mockgen -destination=../../mocks/mock_store.go -package=mocks github.com/jonada182/cover-letter-ai-api/internal/store Store

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/jonada182/cover-letter-ai-api/types"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	JobApplicationSubmission = 0
	JobApplicationInterview  = 1
	JobApplicationAssessment = 2
	JobApplicationOffer      = 3
	JobApplicationCompletion = 4
	JobApplicationRejection  = 5
)

const DateTimeFormat = "2006-01-02 15:04:05"

type StoreClient struct {
	mongoURI string
	dbName   string
}

type Store interface {
	Connect() (*mongo.Client, context.Context, error)
	Disconnect(ctx context.Context, client *mongo.Client)
	GetCareerProfileByEmail(email string) (*types.CareerProfile, error)
	GetCareerProfileByID(profileId uuid.UUID) (*types.CareerProfile, error)
	StoreCareerProfile(careerProfileRequest *types.CareerProfile) (*types.CareerProfile, string, error)
	GetJobApplications(profileId uuid.UUID) (*[]types.JobApplication, error)
	StoreJobApplication(jobApplicationRequest *types.JobApplication) (*types.JobApplication, string, error)
	DeleteJobApplication(jobApplicationId uuid.UUID) error
	StoreAccessToken(profileId uuid.UUID, accessToken string) (string, error)
	ValidateAccessToken(profileId uuid.UUID, accessToken string) (bool, error)
}

// NewStore returns a store client, which has methods to interact with MongoDB
func NewStore() (*StoreClient, error) {
	mongoURI := os.Getenv("MONGODB_URI")
	if mongoURI == "" {
		return nil, errors.New("no Mongo URI defined in env file")
	}

	return &StoreClient{
		mongoURI: mongoURI,
		dbName:   "cover-letter-ai",
	}, nil
}

// Connect establishes a Mongo connection, and returns a MongoDB client
func (store *StoreClient) Connect() (*mongo.Client, context.Context, error) {
	ctx := context.Background()
	// Connect to the database with the given mongoURI
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(store.mongoURI))
	if err != nil {
		log.Printf("Failed to connect to the database: %s", err.Error())
		return nil, nil, err
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Printf("Failed to ping the database:%s", err.Error())
		return nil, nil, err
	}

	fmt.Println("Connected to MongoDB")
	return client, ctx, nil
}

// Disconnect closes the MongoDB connection for the given client
func (store *StoreClient) Disconnect(ctx context.Context, client *mongo.Client) {
	if err := client.Disconnect(ctx); err != nil {
		log.Printf("Failed to disconnect the database: %s", err.Error())
		return
	}
	fmt.Println("Disconnected from MongoDB")
}
