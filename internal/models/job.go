package models

import "time"

type Job struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Company     string    `json:"company"`
	Location    string    `json:"location"`
	Skills      []string  `json:"skills"`
	Salary      float64   `json:"salary"`
	CreatedAt   time.Time `json:"created_at"`
	Score       float64   `json:"score,omitempty"`
}
