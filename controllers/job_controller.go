package controllers

import (
	"net/http"
	"retail_pulse_project/services"

	"github.com/gin-gonic/gin"
)

// SubmitJob handles the HTTP POST request to submit a new job for processing images.
// It accepts a JSON payload containing the job details and returns a job ID upon success.
//
// Endpoint: POST /api/submit/
// Request Body:
// {
//    "count": 2,
//    "visits": [
//       {
//          "store_id": "S00339218",
//          "image_url": [
//             "https://www.gstatic.com/webp/gallery/2.jpg",
//             "https://www.gstatic.com/webp/gallery/3.jpg"
//          ],
//          "visit_time": "2024-11-16T10:00:00Z"
//       }
//    ]
// }
//
// Responses:
// 201 Created: {"job_id": "123"}
// 400 Bad Request: {"error": "Invalid input format"}
//
func SubmitJob(c *gin.Context) {
	var jobInput services.JobInput

	// Bind and validate the JSON request body
	if err := c.ShouldBindJSON(&jobInput); err != nil {
		// Return a 400 Bad Request with an appropriate error message
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input format"})
		return
	}

	// Call the service layer to handle the job submission
	jobID, err := services.SubmitJob(jobInput)
	if err != nil {
		// Return a 400 Bad Request with details of the error
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Respond with the job ID and a 201 Created status
	c.JSON(http.StatusCreated, gin.H{"job_id": jobID})
}

// GetJobStatus handles the HTTP GET request to retrieve the status of a specific job.
// It takes a query parameter `jobid` and returns the job's current status.
//
// Endpoint: GET /api/status?jobid=123
// Query Parameters:
// - jobid: The unique identifier of the job.
//
// Responses:
// 200 OK: {"status": "completed", "job_id": "123"}
// 200 OK (Failed Job): {
//    "status": "failed",
//    "job_id": "",
//    "error": [
//       {"store_id": "S00339218", "error": ""}
//    ]
// }
// 400 Bad Request: {"error": "jobid is required"}
// 404 Not Found: {}
// 500 Internal Server Error: {"error": "unexpected server error"}
//
func GetJobStatus(c *gin.Context) {
	// Retrieve the job ID from query parameters
	jobID := c.Query("jobid")
	if jobID == "" {
		// Return a 400 Bad Request if jobID is not provided
		c.JSON(http.StatusBadRequest, gin.H{"error": "jobid is required"})
		return
	}

	// Call the service layer to fetch the job status
	status, err := services.GetJobStatus(jobID)
	if err != nil {
		// If the job is not found, return a 404 Not Found
		if err.Error() == "job not found" {
			c.JSON(http.StatusNotFound, gin.H{})
		} else {
			// Return a 500 Internal Server Error for unexpected errors
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	// Respond with the job status
	c.JSON(http.StatusOK, status)
}
