package models

import "time"

type Job struct {
	ID int `json:"id"`
	Title string `json:"string"`
	Company string `json:"company"`
	URL string `json:"url"`
	Location string `json:"location"`
	Type string `json:"type"`
	CreatedAt time.Time `json:"created_at"`
}

type IngestRequest struct {
	Jobs []Job `json:"jobs"`
}

type Subscriber struct {
	Email string `json:"email"`
	Preference string `json:"preference"` // internship or new-grad roles
}