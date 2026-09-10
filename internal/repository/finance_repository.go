package repository

import (
	"database/sql"

	"github.com/TakedaB/akademi-api/internal/model"
)

type FinanceRepository struct {
	db *sql.DB
}

func NewFinanceRepository(db *sql.DB) *FinanceRepository {
	return &FinanceRepository{db: db}
}

func (r *FinanceRepository) Create(f *model.Finance) error {
	query := `
		INSERT INTO finance (student_id, description, amount, payment_method, status,due_date)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at, updated_at`

	return r.db.QueryRow(
		query,
		f.StudentID, f.Description, f.Amount, f.PaymentMethod, f.Status, f.DueDate,
	).Scan(&f.ID, &f.CreatedAt, &f.UpdatedAt)
}

func (r *FinanceRepository) FindAll() ([]model.Finance, error) {
	query := `SELECT id, student_id, description, amount, payment_method, status, due_date, created_at, updated_at FROM finance ORDER BY due_date`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var records []model.Finance
	for rows.Next() {
		var f model.Finance
		if err := rows.Scan(&f.ID, &f.StudentID, &f.Description, &f.Amount, &f.PaymentMethod, &f.Status, &f.DueDate, &f.CreatedAt, &f.UpdatedAt); err != nil {
			return nil, err
		}
		records = append(records, f)
	}
	return records, rows.Err()
}

func (r *FinanceRepository) FindByStudentID(studentID string) ([]model.Finance, error) {
	query := `SELECT id, student_id, description, amount, payment_method, status, due_date, created_at, updated_at FROM finance WHERE student_id = $1 ORDER BY due_date`

	rows, err := r.db.Query(query, studentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var records []model.Finance
	for rows.Next() {
		var f model.Finance
		if err := rows.Scan(&f.ID, &f.StudentID, &f.Description, &f.Amount, &f.PaymentMethod, &f.Status, &f.DueDate, &f.CreatedAt, &f.UpdatedAt); err != nil {
			return nil, err
		}
		records = append(records, f)
	}
	return records, rows.Err()
}

func (r *FinanceRepository) UpdateStatus(id, status string) (*model.Finance, error) {
	query := `
		UPDATE finance
		SET status = $1, updated_at = now()
		WHERE id = $2
		RETURNING id, student_id, description, amount, payment_method, status, due_date, created_at, updated_at`

	var f model.Finance
	err := r.db.QueryRow(query, status, id).Scan(
		&f.ID, &f.StudentID, &f.Description, &f.Amount, &f.PaymentMethod, &f.Status, &f.DueDate, &f.CreatedAt, &f.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &f, nil
}

func (r *FinanceRepository) Delete(id string) error {
	query := `DELETE FROM finance WHERE id = $1`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}
	return nil
}
