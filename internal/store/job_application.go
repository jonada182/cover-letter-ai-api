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

// DeleteJobApplication retrieves an array of job applications from MongoDB
func (store *StoreClient) DeleteJobApplication(jobApplicationId uuid.UUID) error {
	mongoClient, ctx, err := store.Connect()
	if err != nil {
		return err
	}
	defer store.Disconnect(ctx, mongoClient)

	// Get the job_applications collection from the database client
	collection := mongoClient.Database(store.dbName).Collection("job_applications")
	// Delete job applications by id
	log.Printf("Deleting job application for %s", jobApplicationId.String())
	result, err := collection.DeleteOne(ctx, bson.M{"id": jobApplicationId})
	if err != nil {
		log.Printf("Failed to delete job application:%s", err.Error())
		return err
	}

	if result.DeletedCount == 0 {
		return errors.New("failed to delete job application")
	}

	log.Printf("Deleted job application for %s", jobApplicationId.String())
	return nil
}
