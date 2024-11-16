package main

import (
	"fmt"
	"retail_pulse_project/config"
	"retail_pulse_project/routes"
	"retail_pulse_project/utils"

	"github.com/gin-gonic/gin"
)

func main() {
	config.ConnectDB()

	err := utils.LoadStoresFromCSV("stores.csv") // Update the path to your CSV file
	if err != nil {
		fmt.Println("Error loading store data:", err)
		return
	}	

	router := gin.Default()
	routes.SetupRoutes(router)

	router.Run(":8080")
}
