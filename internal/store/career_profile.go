package store

import (
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/jonada182/cover-letter-ai-api/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

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
