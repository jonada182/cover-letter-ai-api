package store

//go:generate mockgen -destination=../../mocks/mock_store.go -package=mocks github.com/jonada182/cover-letter-ai-api/internal/store Store

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jonada182/cover-letter-ai-api/types"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
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

// StoreCareerProfile upserts a CareerProfile in MongoDB
func (store *StoreClient) StoreCareerProfile(careerProfile *types.CareerProfile) (*types.CareerProfile, string, error) {
	mongoClient, ctx, err := store.Connect()
	if err != nil {
		return nil, "", err
	}
	defer store.Disconnect(ctx, mongoClient)

	// Get the profiles collection from the database client
	collection := mongoClient.Database(store.dbName).Collection("profiles")
	careerProfileID := uuid.New()
	if careerProfile.ID != uuid.Nil {
		careerProfileID = careerProfile.ID
	}
	careerProfileRow := &types.CareerProfile{
		ID:              careerProfileID,
		FirstName:       careerProfile.FirstName,
		LastName:        careerProfile.LastName,
		Headline:        careerProfile.Headline,
		ExperienceYears: careerProfile.ExperienceYears,
		Summary:         careerProfile.Summary,
		Skills:          careerProfile.Skills,
		ContactInfo:     careerProfile.ContactInfo,
	}
	// Set up update options to ensure the values are overwritten in the database
	update := bson.M{"$set": careerProfileRow}
	updateOptions := options.Update().SetUpsert(true)
	result, err := collection.UpdateOne(
		ctx,
		bson.M{"contact_info.email": careerProfile.ContactInfo.Email},
		update,
		updateOptions,
	)
	if err != nil {
		log.Printf("Failed to update profile:%s", err.Error())
		return nil, "", err
	}

	// Check if upsert resulted in an insert (new document)
	var responseMsg string
	if result.UpsertedCount > 0 {
		responseMsg = "career profile has been inserted"
		fmt.Printf("%s:", result.UpsertedID)
	} else {
		responseMsg = "career profile has been updated"
		fmt.Printf("%d:", result.ModifiedCount)
	}

	return careerProfileRow, responseMsg, nil
}

// GetCareerProfile retrieves a CareerProfile from MongoDB
func (store *StoreClient) GetCareerProfileByEmail(email string) (*types.CareerProfile, error) {
	mongoClient, ctx, err := store.Connect()
	if err != nil {
		return nil, err
	}
	defer store.Disconnect(ctx, mongoClient)

	var careerProfile types.CareerProfile
	// Get the profiles collection from the database client
	collection := mongoClient.Database(store.dbName).Collection("profiles")
	// Find career profile using the contact_info.email and the given email address
	err = collection.FindOne(ctx, bson.M{"contact_info.email": email}).Decode(&careerProfile)
	if err != nil {
		log.Printf("Failed to find profile:%s", err.Error())
		return nil, err
	}

	return &careerProfile, nil
}

// GetCareerProfileByID retrieves a CareerProfile using the ID from MongoDB
func (store *StoreClient) GetCareerProfileByID(profileId uuid.UUID) (*types.CareerProfile, error) {
	mongoClient, ctx, err := store.Connect()
	if err != nil {
		return nil, err
	}
	defer store.Disconnect(ctx, mongoClient)

	var careerProfile types.CareerProfile
	// Get the profiles collection from the database client
	collection := mongoClient.Database(store.dbName).Collection("profiles")
	// Find career profile using the given profile ID
	err = collection.FindOne(ctx, bson.M{"id": profileId}).Decode(&careerProfile)
	if err != nil {
		log.Printf("Failed to find profile:%s", err.Error())
		return nil, err
	}

	return &careerProfile, nil
}

// GetJobApplications retrieves an array of job applications from MongoDB
func (store *StoreClient) GetJobApplications(profileId uuid.UUID) (*[]types.JobApplication, error) {
	mongoClient, ctx, err := store.Connect()
	if err != nil {
		return nil, err
	}
	defer store.Disconnect(ctx, mongoClient)

	var jobApplications []types.JobApplication
	// Get the job_applications collection from the database client
	collection := mongoClient.Database(store.dbName).Collection("job_applications")
	// Find job applications using the career profile id
	log.Printf("Find applications for %s", profileId.String())
	cur, err := collection.Find(ctx, bson.M{"profile_id": profileId})
	if err != nil {
		log.Printf("Failed to retrieve job applications:%s", err.Error())
		return nil, err
	}
	cur.All(ctx, &jobApplications)
	defer cur.Close(ctx)
	if cur.Err() != nil {
		log.Printf("Failed to retrieve job applications:%s", err.Error())
		return nil, err
	}

	if len(jobApplications) == 0 {
		return nil, errors.New("no job applications found")
	}

	return &jobApplications, nil
}

// StoreJobApplication upserts a JobApplication in MongoDB
func (store *StoreClient) StoreJobApplication(jobApplicationRequest *types.JobApplication) (*types.JobApplication, string, error) {
	mongoClient, ctx, err := store.Connect()
	if err != nil {
		return nil, "", err
	}
	defer store.Disconnect(ctx, mongoClient)
	isNew := true
	jobApplicationID := uuid.New()
	if jobApplicationRequest.ID != uuid.Nil {
		isNew = false
		jobApplicationID = jobApplicationRequest.ID
	}

	// Get the profiles collection from the database client
	collection := mongoClient.Database(store.dbName).Collection("job_applications")
	currentDateTime := time.Now().Format("2006-01-02 15:04:05")
	jobApplicationRow := &types.JobApplication{
		ID:          jobApplicationID,
		ProfileID:   jobApplicationRequest.ProfileID,
		CompanyName: jobApplicationRequest.CompanyName,
		JobRole:     jobApplicationRequest.JobRole,
		URL:         jobApplicationRequest.URL,
		Events:      jobApplicationRequest.Events,
		CreatedAt:   &currentDateTime,
		UpdatedAt:   &currentDateTime,
	}

	if isNew && jobApplicationRow.Events == nil {
		jobApplicationRow.Events = &[]types.JobApplicationEvent{}
		*jobApplicationRow.Events = append(*jobApplicationRow.Events, types.JobApplicationEvent{
			Type:        JobApplicationSubmission,
			Description: "Application Sent",
			Date:        currentDateTime,
		})
	}

	// Set up update options to ensure the values are overwritten in the database
	update := bson.M{"$set": jobApplicationRow}
	updateOptions := options.Update().SetUpsert(true)
	result, err := collection.UpdateOne(
		ctx,
		bson.M{"id": jobApplicationRow.ID},
		update,
		updateOptions,
	)
	if err != nil {
		log.Printf("Failed to update job application:%s", err.Error())
		return nil, "", err
	}

	// Check if upsert resulted in an insert (new document)
	var responseMsg string
	if result.UpsertedCount > 0 {
		responseMsg = "job application has been inserted"
		fmt.Printf("%s:", result.UpsertedID)
	} else {
		responseMsg = "job application has been updated"
		fmt.Printf("%d:", result.ModifiedCount)
	}

	return jobApplicationRow, responseMsg, nil
}
