package store

import (
	"context"
	"cover-letter-ai-api/types"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Store struct {
	mongoURI string
	dbName   string
}

func NewStore() (*Store, error) {
	mongoURI := os.Getenv("MONGODB_URI")
	if mongoURI == "" {
		return nil, errors.New("no Mongo URI defined in env file")
	}

	return &Store{
		mongoURI: mongoURI,
		dbName:   "cover-letter-ai",
	}, nil
}

func (store *Store) Connect() (*mongo.Client, context.Context, error) {
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

func (store *Store) Disconnect(ctx context.Context, client *mongo.Client) {
	if err := client.Disconnect(ctx); err != nil {
		log.Printf("Failed to disconnect the database: %s", err.Error())
		return
	}
	fmt.Println("Disconnected from MongoDB")
}

func (store *Store) StoreCareerProfile(careerProfileRequest *types.CareerProfileRequest) (*types.CareerProfile, string, error) {
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
		fmt.Printf("%s:", result.ModifiedCount)
	}

	return careerProfileRow, responseMsg, nil
}

func (store *Store) GetCareerProfile(email string) (*types.CareerProfile, error) {
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
