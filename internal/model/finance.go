package model

import "time"

type Finance struct {
	ID            string    `json:"id"`
	StudentID     string    `json:"student_id"`
	Description   string    `json:"description"`
	Amount        float64   `json:"amount"`
	PaymentMethod string    `json:"payment_method"`
	Status        string    `json:"status"`
	DueDate       time.Time `json:"due_date"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
