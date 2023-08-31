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
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type StoreClient struct {
	mongoURI string
	dbName   string
}

type Store interface {
	Connect() (*mongo.Client, context.Context, error)
	Disconnect(ctx context.Context, client *mongo.Client)
	GetCareerProfile(email string) (*types.CareerProfile, error)
	StoreCareerProfile(careerProfileRequest *types.CareerProfileRequest) (*types.CareerProfile, string, error)
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
func (store *StoreClient) StoreCareerProfile(careerProfileRequest *types.CareerProfileRequest) (*types.CareerProfile, string, error) {
	mongoClient, ctx, err := store.Connect()
	if err != nil {
		return nil, "", err
	}
	defer store.Disconnect(ctx, mongoClient)

	collection := mongoClient.Database(store.dbName).Collection("profiles")
	careerProfileRow := &types.CareerProfile{
		ID:              uuid.New(),
		FirstName:       careerProfileRequest.FirstName,
		LastName:        careerProfileRequest.LastName,
		Headline:        careerProfileRequest.Headline,
		ExperienceYears: careerProfileRequest.ExperienceYears,
		Summary:         careerProfileRequest.Summary,
		Skills:          careerProfileRequest.Skills,
		ContactInfo:     careerProfileRequest.ContactInfo,
	}
	update := bson.M{"$set": careerProfileRow}
	updateOptions := options.Update().SetUpsert(true)
	result, err := collection.UpdateOne(
		ctx,
		bson.M{"contact_info.email": careerProfileRequest.ContactInfo.Email},
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
func (store *StoreClient) GetCareerProfile(email string) (*types.CareerProfile, error) {
	mongoClient, ctx, err := store.Connect()
	if err != nil {
		return nil, err
	}
	defer store.Disconnect(ctx, mongoClient)

	var careerProfile types.CareerProfile
	collection := mongoClient.Database(store.dbName).Collection("profiles")
	err = collection.FindOne(ctx, bson.M{"contact_info.email": email}).Decode(&careerProfile)
	if err != nil {
		log.Printf("Failed to find profile:%s", err.Error())
		return nil, err
	}

	return &careerProfile, nil
}
