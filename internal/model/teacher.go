package model

import "time"

type Teacher struct {
	ID            string    `json:"id"`
	UserID        string    `json:"user_id"`
	Name          string    `json:"name"`
	Email         string    `json:"email"`
	Subject       string    `json:"subject"`
	Phone         string    `json:"phone"`
	HireDate      time.Time `json:"hire_date"`
	ClassAssigned string    `json:"class_assigned"`
	WorkloadHours int       `json:"workload_hours"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
