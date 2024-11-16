package utils

import (
	"encoding/csv"
	"fmt"
	"os"
	"retail_pulse_project/config"
	"retail_pulse_project/models"
)

// LoadStoresFromCSV loads store data from a CSV file into the Store table
func LoadStoresFromCSV(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open CSV file: %v", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return fmt.Errorf("failed to read CSV file: %v", err)
	}

	// Skip the header row (if present)
	for i, record := range records {
		if i == 0 {
			continue // Skip header row
		}

		// Parse fields
		areaCode := record[0]
		storeName := record[1]
		storeID := record[2]

		// Check if the store already exists
		var existingStore models.Store
		result := config.DB.First(&existingStore, "id = ?", storeID)
		if result.Error == nil {
			continue // Store already exists, skip
		}

		// Insert new store
		store := models.Store{
			ID:       storeID,
			Name:     storeName,
			AreaCode: areaCode,
		}
		if err := config.DB.Create(&store).Error; err != nil {
			return fmt.Errorf("failed to insert store: %v", err)
		}
	}

	fmt.Println("Store data successfully loaded from CSV.")
	return nil
}
