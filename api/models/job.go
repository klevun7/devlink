package models

import (
	"time"
	"gorm.io/gorm"
)

type Job struct {
	gorm.Model
	Title    string    `json:"title"`
	Company  string    `json:"company"`
	Location string 	`json:"location"`
	URL      string    `json:"url" gorm:"uniqueIndex"` 
	PostedAt time.Time `json:"posted_at"`
	
}