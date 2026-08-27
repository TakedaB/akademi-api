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

func (r *StudentRepository) Create(s *model.Student) error {
	query := `
		INSERT INTO students (enrollment_number, name, birth_date, parent_name, city, phone, email, grade )
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, created_at, updated_at`

	return r.db.QueryRow(
		query,
		s.EnrollmentNumber, s.Name, s.BirthDate, s.ParentName, s.City, s.Phone, s.Email, s.Grade,
	).Scan(&s.ID, &s.CreatedAt, &s.UpdatedAt)
}

func (r *StudentRepository) FindAll() ([]model.Student, error) {
	query := `SELECT id, enrollment_number, name, birth_date, parent_name, city, phone, email, grade, created_at, updated_at FROM students ORDER BY name `

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var students []model.Student
	for rows.Next() {
		var s model.Student
		if err := rows.Scan(&s.ID, &s.EnrollmentNumber, &s.Name, &s.BirthDate, &s.ParentName, &s.City, &s.Phone, &s.Email, &s.Grade, &s.CreatedAt, &s.UpdatedAt); err != nil {
			return nil, err
		}
		students = append(students, s)
	}

	return students, rows.Err()
}

func (r *StudentRepository) FindByID(id string) (*model.Student, error) {
	query := `SELECT id, enrollment_number, name, birth_date, parent_name, city, phone, email, grade, created_at, updated_at FROM students WHERE id = $1`

	var s model.Student
	err := r.db.QueryRow(query, id).Scan(&s.ID, &s.EnrollmentNumber, &s.Name, &s.BirthDate, &s.ParentName, &s.City, &s.Phone, &s.Email, &s.Grade, &s.CreatedAt, &s.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *StudentRepository) Update(s *model.Student) error {
	query := `
		UPDATE students
		SET enrollment_number = $1, name = $2, birth_date = $3, parent_name = $4, city = $5, phone = $6, email = $7, grade = $8, updated_at = now()
		WHERE id = $9
		RETURNING updated_at`

	return r.db.QueryRow(
		query,
		s.EnrollmentNumber, s.Name, s.BirthDate, s.ParentName, s.City, s.Phone, s.Email, s.Grade, s.ID,
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
