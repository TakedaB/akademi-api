package model

import "time"

type Student struct {
	ID               string    `json:"id"`
	EnrollmentNumber string    `json:"enrollment_number"`
	Name             string    `json:"name"`
	BirthDate        string    `json:"birth_date"`
	ParentName       string    `json:"parent_name"`
	City             string    `json:"city,omitempty"`
	Phone            string    `json:"phone"`
	Email            string    `json:"email,omitempty"`
	Grade            string    `json:"grade,omitempty"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}
