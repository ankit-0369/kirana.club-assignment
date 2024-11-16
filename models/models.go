package models

import "time"

// Job represents a processing job for a store
type Job struct {
	ID        string    `gorm:"primaryKey" json:"id"`
	StoreID   string    `gorm:"index;not null;column:store_id" json:"store_id"` // Maps to `store_id` in the database
	Status    string    `json:"status"`                                       // ongoing, completed, or failed
	Error     string    `json:"error"`                                        // Optional field for job-level errors
	VisitTime time.Time `json:"visit_time"`                                   // Time of store visit
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	Images []Image `gorm:"foreignKey:JobID;constraint:OnDelete:CASCADE;" json:"images"` // One-to-Many relationship with Images
	Store  Store   `gorm:"foreignKey:StoreID;constraint:OnDelete:CASCADE;" json:"store"` // Links to the Store table
}

// Image represents an image being processed as part of a job
type Image struct {
	ID        string    `gorm:"primaryKey" json:"id"`
	JobID     string    `gorm:"index;not null;column:job_id;constraint:OnDelete:CASCADE;" json:"job_id"` // Maps to `job_id` in the database
	StoreID   string    `gorm:"column:store_id" json:"store_id"`                                         // Maps to `store_id` in the database
	ImageURL  string    `json:"image_url"`
	Perimeter float64   `json:"perimeter"` // Calculated as 2 * (Height + Width)
	Status    string    `json:"status"`    // processed or failed
	Error     string    `json:"error"`     // Optional field for image-level errors
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

// Store represents a store where jobs are created
type Store struct {
	ID       string `gorm:"primaryKey" json:"id"`
	Name     string `json:"name"`
	AreaCode string `json:"area_code"`

	Jobs []Job `gorm:"foreignKey:StoreID;constraint:OnDelete:CASCADE;" json:"jobs"` // One-to-Many relationship with Jobs
}
