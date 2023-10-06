package handler

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jonada182/cover-letter-ai-api/types"
)

// HandleCreateJobApplication handles a POST method to create a job application in MongoDB
func (h *Handler) HandleCreateJobApplication(c *gin.Context) {
	// Receive JobApplication parameters from request payload
	var jobApplicationRequest types.JobApplication
	if err := c.ShouldBindJSON(&jobApplicationRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("error retrieving JSON: %s", err.Error())})
		return
	}

	if jobApplicationRequest.ProfileID.String() == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "profile id is required"})
		return
	}

	if jobApplicationRequest.CompanyName == "" || jobApplicationRequest.JobRole == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "company name and job role are required"})
		return
	}

	// Call store method to upsert Job Application in MongoDB
	jobApplication, responseMsq, err := h.StoreClient.StoreJobApplication(&jobApplicationRequest)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": responseMsq, "data": jobApplication})
}

// HandleGetJobApplications handles a GET method to retrieve job applications from MongoDB
func (h *Handler) HandleGetJobApplications(c *gin.Context) {
	profileIdParam, exists := c.Get("ProfileID")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no profile_id provided in the request"})
		return
	}

	profileId, ok := profileIdParam.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid profile id"})
		return
	}

	// Call store method to retrieve []JobApplication from MongoDB
	jobApplications, err := h.StoreClient.GetJobApplications(profileId)
	if err != nil && strings.Contains(err.Error(), "no job applications found") {
		c.JSON(http.StatusNotFound, gin.H{"error": "no job applications found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": &jobApplications})
}

// HandleGetJobApplicationByID handles a GET method to retrieve a job application from MongoDB
func (h *Handler) HandleGetJobApplicationByID(c *gin.Context) {
	jobApplicationIdParam := c.Param("id")
	if jobApplicationIdParam == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no job application id provided in the request"})
		return
	}

	jobApplicationId, err := uuid.Parse(jobApplicationIdParam)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Call store method to retrieve JobApplication from MongoDB
	jobApplication, err := h.StoreClient.GetJobApplicationByID(jobApplicationId)
	if err != nil && strings.Contains(err.Error(), "no document") {
		c.JSON(http.StatusNotFound, gin.H{"error": "job application not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": &jobApplication})
}

// HandleDeleteJobApplication handles a DELETE request to delete a job application by ID
func (h *Handler) HandleDeleteJobApplication(c *gin.Context) {
	jobApplicationID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid job application id"})
	}

	err = h.StoreClient.DeleteJobApplication(jobApplicationID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "job application deleted successfully"})
}
