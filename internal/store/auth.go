package store

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/jonada182/cover-letter-ai-api/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Set a token duration of 7 days
const TokenDuration = 7 * 24 * time.Hour

// StoreAccessToken stores an access_token for a given profile_id
func (store *StoreClient) StoreAccessToken(profileId uuid.UUID, accessToken string) (string, error) {
	mongoClient, ctx, err := store.Connect()
	if err != nil {
		return "", err
	}
	defer store.Disconnect(ctx, mongoClient)

	// Get the access_tokens collection from the database client
	collection := mongoClient.Database(store.dbName).Collection("access_tokens")
	expiresAt := time.Now().Add(TokenDuration).Format(DateTimeFormat)
	accessTokenRow := &types.AccessToken{
		ProfileID:   profileId,
		AccessToken: accessToken,
		ExpiresAt:   expiresAt,
	}
	// Set up update options to ensure the values are overwritten in the database
	update := bson.M{"$set": accessTokenRow}
	updateOptions := options.Update().SetUpsert(true)
	result, err := collection.UpdateOne(
		ctx,
		bson.M{"profile_id": profileId},
		update,
		updateOptions,
	)
	if err != nil {
		log.Printf("Failed to store access token:%s", err.Error())
		return "", err
	}

	// Check if upsert resulted in an insert (new document)
	var responseMsg string
	if result.UpsertedCount > 0 {
		responseMsg = "access token has been stored"
		fmt.Printf("%s:", result.UpsertedID)
	} else {
		responseMsg = "access token has been updated"
		fmt.Printf("%d:", result.ModifiedCount)
	}

	return responseMsg, nil
}

// ValidateAccessToken checks that a given access_token is valid for a profile_id
func (store *StoreClient) ValidateAccessToken(profileId uuid.UUID, accessToken string) (bool, error) {
	mongoClient, ctx, err := store.Connect()
	if err != nil {
		return false, err
	}
	defer store.Disconnect(ctx, mongoClient)

	var currentAccessToken types.AccessToken
	// Get the access_tokens collection from the database client
	collection := mongoClient.Database(store.dbName).Collection("access_tokens")
	// Find access token using the given profile ID and token
	err = collection.FindOne(ctx, bson.M{"profile_id": profileId}).Decode(&currentAccessToken)
	if err != nil {
		log.Printf("Failed to find access token:%s", err.Error())
		return false, err
	}

	if currentAccessToken.AccessToken != accessToken {
		err = errors.New("access_token provided is invalid")
		log.Println(err.Error())
		return false, err
	}

	expiresAt, err := time.Parse(DateTimeFormat, currentAccessToken.ExpiresAt)
	if err != nil {
		log.Printf("Failed to parse expires_at:%s", err.Error())
		return false, err
	}

	if time.Now().After(expiresAt) {
		err = errors.New("access_token provided has expired")
		log.Println(err.Error())
		return false, err
	}

	return true, nil
}
