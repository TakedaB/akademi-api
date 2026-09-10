package repository

import (
	"database/sql"

	"github.com/TakedaB/akademi-api/internal/model"
)

type TeacherRepository struct {
	db *sql.DB
}

func NewTeacherRepository(db *sql.DB) *TeacherRepository {
	return &TeacherRepository{db: db}
}

func (r *TeacherRepository) Create(tx *sql.Tx, t *model.Teacher) error {
	query := `
		INSERT INTO teachers (user_id, subject, phone, hire_date, class_assigned, workload_hours)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at, updated_at`

	return tx.QueryRow(
		query,
		t.UserID, t.Subject, t.Phone, t.HireDate, t.ClassAssigned, t.WorkloadHours,
	).Scan(&t.ID, &t.CreatedAt, &t.UpdatedAt)
}

func (r *TeacherRepository) FindAll() ([]model.Teacher, error) {
	query := `
		SELECT t.id, t.user_id, u.name, u.email, t.subject, t.phone, t.hire_date, t.class_assigned, t.workload_hours, t.created_at, t.updated_at
		FROM teachers t
		JOIN users u ON u.id = t.user_id
		ORDER BY u.name`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	teachers := []model.Teacher{}
	for rows.Next() {
		var t model.Teacher
		if err := rows.Scan(&t.ID, &t.UserID, &t.Name, &t.Email, &t.Subject, &t.Phone, &t.HireDate, &t.ClassAssigned, &t.WorkloadHours, &t.CreatedAt, &t.UpdatedAt); err != nil {
			return nil, err
		}
		teachers = append(teachers, t)
	}

	return teachers, rows.Err()
}

func (r *TeacherRepository) FindByID(id string) (*model.Teacher, error) {
	query := `
		SELECT t.id, t.user_id, u.name, u.email, t.subject, t.phone, t.hire_date, t.class_assigned, t.workload_hours, t.created_at, t.updated_at
		FROM teachers t
		JOIN users u ON u.id = t.user_id
		WHERE t.id = $1`

	var t model.Teacher
	err := r.db.QueryRow(query, id).Scan(&t.ID, &t.UserID, &t.Name, &t.Email, &t.Subject, &t.Phone, &t.HireDate, &t.ClassAssigned, &t.WorkloadHours, &t.CreatedAt, &t.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &t, nil
}

func (r *TeacherRepository) Delete(id string) error {
	query := `DELETE FROM teachers WHERE id = $1`

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
