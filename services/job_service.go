package services

import (
	"errors"
	"fmt"
	"image"
	_ "image/jpeg" // Support for JPEG images
	_ "image/png"  // Support for PNG images
	"math/rand"
	"net/http"
	"retail_pulse_project/config"
	"retail_pulse_project/models"
	"time"
)

// JobInput represents the structure of the incoming job submission request payload.
type JobInput struct {
	Count  int `json:"count"`
	Visits []struct {
		StoreID   string   `json:"store_id"`
		ImageURLs []string `json:"image_url"`
		VisitTime string   `json:"visit_time"`
	} `json:"visits"`
}

// SubmitJob handles the creation of jobs and processing of associated images.
func SubmitJob(input JobInput) (string, error) {
	// Validate that the count matches the number of visits
	if len(input.Visits) != input.Count {
		return "", errors.New("invalid count: does not match the number of visits")
	}

	var jobID string // To store the ID of the first job created

	// Iterate over visits to create jobs and process images
	for i, visit := range input.Visits {
		// Validate the StoreID by checking its existence in the database
		var store models.Store
		err := config.DB.First(&store, "id = ?", visit.StoreID).Error
		if err != nil {
			return "", fmt.Errorf("store with ID %s not found", visit.StoreID)
		}

		// Create a new job
		job := models.Job{
			ID:        generateID(),
			StoreID:   visit.StoreID,
			Status:    "ongoing", // Set initial status to 'ongoing'
			VisitTime: parseTime(visit.VisitTime),
		}

		// Save the job to the database
		config.DB.Create(&job)

		// Store the first job's ID
		if i == 0 {
			jobID = job.ID
		}

		// Process each image in the job
		for _, imageURL := range visit.ImageURLs {
			image := models.Image{
				ID:       generateID(),
				JobID:    job.ID,
				StoreID:  visit.StoreID,
				ImageURL: imageURL,
				Status:   "pending", // Initial status of the image
			}

			// Save the image record to the database
			config.DB.Create(&image)

			// Process the image asynchronously
			go processImage(image.ID, imageURL, job.ID)
		}
	}

	// Return the first job ID created
	return jobID, nil
}

// GetJobStatus retrieves the status of a job and returns the appropriate response.
func GetJobStatus(jobID string) (map[string]interface{}, error) {
	var job models.Job
	err := config.DB.First(&job, "id = ?", jobID).Error
	if err != nil {
		return nil, errors.New("job not found")
	}

	// If the job status is 'failed', collect errors for failed images
	if job.Status == "failed" {
		var failedImages []models.Image
		config.DB.Where("job_id = ? AND status = ?", jobID, "failed").Find(&failedImages)

		// Collect store-level errors for the failed images
		var storeErrors []map[string]interface{}
		for _, img := range failedImages {
			storeErrors = append(storeErrors, map[string]interface{}{
				"store_id": img.StoreID,
				"error":    "", // Error details can be added here if necessary
			})
		}

		// Return the failure response with errors
		return map[string]interface{}{
			"status": "failed",
			"job_id": "",
			"error":  storeErrors,
		}, nil
	}

	// Otherwise, return the job status as completed or ongoing
	return map[string]interface{}{
		"status": job.Status,
		"job_id": jobID,
	}, nil
}

// processImage fetches image dimensions, calculates the perimeter, and updates the database.
func processImage(imageID, imageURL, jobID string) {
	// Simulate GPU processing with a random delay
	time.Sleep(time.Duration(100+rand.Intn(300)) * time.Millisecond)

	// Get the dimensions of the image
	width, height, err := getImageDimensions(imageURL)
	if err != nil {
		// Update the image status to 'failed' in case of an error
		config.DB.Model(&models.Image{}).Where("id = ?", imageID).Updates(models.Image{
			Status: "failed",
			Error:  err.Error(),
		})

		// Update the job status to 'failed' if any image processing fails
		updateJobStatusToFailed(jobID)
		return
	}

	// Calculate the perimeter of the image
	perimeter := 2 * (float64(width) + float64(height))

	// Update the image processing result in the database
	config.DB.Model(&models.Image{}).Where("id = ?", imageID).Updates(models.Image{
		Perimeter: perimeter,
		Status:    "processed",
	})

	// Check if all images for the job are processed and update the job status
	updateJobStatus(jobID)
}

// updateJobStatus sets the job status to 'completed' if all associated images are processed.
func updateJobStatus(jobID string) {
	var images []models.Image

	// Retrieve all images associated with the job
	config.DB.Where("job_id = ?", jobID).Find(&images)

	// Check if all images are processed
	allProcessed := true
	for _, image := range images {
		if image.Status != "processed" {
			allProcessed = false
			break
		}
	}

	// If all images are processed, update the job status to 'completed'
	if allProcessed {
		config.DB.Model(&models.Job{}).Where("id = ?", jobID).Update("status", "completed")
	}
}

// updateJobStatusToFailed sets the job status to 'failed' if any image processing fails.
func updateJobStatusToFailed(jobID string) {
	config.DB.Model(&models.Job{}).Where("id = ?", jobID).Update("status", "failed")
}

// getImageDimensions retrieves the width and height of an image from its URL.
func getImageDimensions(imageURL string) (int, int, error) {
	resp, err := http.Get(imageURL)
	if err != nil {
		return 0, 0, err
	}
	defer resp.Body.Close()

	imgConfig, _, err := image.DecodeConfig(resp.Body)
	if err != nil {
		return 0, 0, err
	}

	return imgConfig.Width, imgConfig.Height, nil
}

// generateID generates a unique ID for jobs and images based on the current timestamp.
func generateID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}

// parseTime converts a string into a time.Time object.
func parseTime(t string) time.Time {
	parsed, _ := time.Parse(time.RFC3339, t)
	return parsed
}
