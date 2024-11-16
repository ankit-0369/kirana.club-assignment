package routes

import (
	"retail_pulse_project/controllers"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {
	router.POST("/api/submit", controllers.SubmitJob)
	router.GET("/api/status", controllers.GetJobStatus)
}
