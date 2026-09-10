package repository

import (
	"database/sql"

	"github.com/TakedaB/akademi-api/internal/model"
)

type StudentRepository struct {
	db *sql.DB
}

func NewStudentRepository(db *sql.DB) *StudentRepository {
	return &StudentRepository{db: db}
}

func (r *StudentRepository) Create(tx *sql.Tx, s *model.Student) error {
	query := `
		INSERT INTO students (user_id, enrollment_number, birth_date, parent_name, city, phone, grade)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, created_at, updated_at`

	return tx.QueryRow(
		query,
		s.UserID, s.EnrollmentNumber, s.BirthDate, s.ParentName, s.City, s.Phone, s.Grade,
	).Scan(&s.ID, &s.CreatedAt, &s.UpdatedAt)
}

func (r *StudentRepository) FindAll() ([]model.Student, error) {
	query := `
		SELECT s.id, s.user_id, u.name, u.email, s.enrollment_number, s.birth_date, s.parent_name, s.city, s.phone, s.grade, s.created_at, s.updated_at
		FROM students s
		JOIN users u ON u.id = s.user_id
		ORDER BY u.name`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var students []model.Student
	for rows.Next() {
		var s model.Student
		if err := rows.Scan(&s.ID, &s.UserID, &s.Name, &s.Email, &s.EnrollmentNumber, &s.BirthDate, &s.ParentName, &s.City, &s.Phone, &s.Grade, &s.CreatedAt, &s.UpdatedAt); err != nil {
			return nil, err
		}
		students = append(students, s)
	}

	return students, rows.Err()
}

func (r *StudentRepository) FindByID(id string) (*model.Student, error) {
	query := `
		SELECT s.id, s.user_id, u.name, u.email, s.enrollment_number, s.birth_date, s.parent_name, s.city, s.phone, s.grade, s.created_at, s.updated_at
		FROM students s
		JOIN users u ON u.id = s.user_id
		WHERE s.id = $1`

	var s model.Student
	err := r.db.QueryRow(query, id).Scan(&s.ID, &s.UserID, &s.Name, &s.Email, &s.EnrollmentNumber, &s.BirthDate, &s.ParentName, &s.City, &s.Phone, &s.Grade, &s.CreatedAt, &s.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *StudentRepository) Update(s *model.Student) error {
	query := `
		UPDATE students
		SET birth_date = $1, parent_name = $2, city = $3, phone = $4, grade = $5, updated_at = now()
		WHERE id = $6
		RETURNING updated_at`

	return r.db.QueryRow(
		query,
		s.BirthDate, s.ParentName, s.City, s.Phone, s.Grade, s.ID,
	).Scan(&s.UpdatedAt)
}

func (r *StudentRepository) Delete(id string) error {
	query := `DELETE FROM students WHERE id = $1`

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

func (r *StudentRepository) NextSequential(year int) (int, error) {
	query := `
		INSERT INTO enrollment_counters (year, count)
		VALUES ($1, 1)
		ON CONFLICT (year) DO UPDATE SET count = enrollment_counters.count + 1
		RETURNING count`

	var count int
	err := r.db.QueryRow(query, year).Scan(&count)
	return count, err
}
