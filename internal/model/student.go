package model

import "time"

type Student struct {
	ID               string    `json:"id"`
	UserID           string    `json:"user_id"`
	EnrollmentNumber string    `json:"enrollment_number"`
	Name             string    `json:"name"`
	BirthDate        time.Time `json:"birth_date"`
	ParentName       string    `json:"parent_name"`
	City             string    `json:"city,omitempty"`
	Phone            string    `json:"phone"`
	Email            string    `json:"email,omitempty"`
	Grade            string    `json:"grade,omitempty"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}
