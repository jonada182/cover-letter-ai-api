package clients

import (
	"context"
	"errors"
	"log"
	"os"
	"resu-mate-api/types"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Store struct {
	dbName      string
	mongoClient *mongo.Client
	mongoCtx    context.Context
}

func NewStore() (*Store, error) {
	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		return nil, errors.New("no Mongo URI defined in env file")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)

	return &Store{
		dbName:      "career",
		mongoClient: client,
		mongoCtx:    ctx,
	}, nil
}

func (store *Store) StoreCareerProfile(careerProfile *types.CareerProfile) (*mongo.InsertOneResult, error) {
	collection := store.mongoClient.Database(store.dbName).Collection("profiles")
	result, err := collection.InsertOne(store.mongoCtx, careerProfile)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return result, nil
}

func (store *Store) GetCareerProfile(careerProfileID int) (*types.CareerProfile, error) {
	var careerProfile types.CareerProfile
	collection := store.mongoClient.Database(store.dbName).Collection("profiles")
	err := collection.FindOne(store.mongoCtx, bson.M{"id": careerProfileID}).Decode(&careerProfile)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return &careerProfile, nil
}
