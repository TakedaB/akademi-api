package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/TakedaB/akademi-api/internal/model"
	"github.com/TakedaB/akademi-api/internal/repository"
)

var (
	ErrNameRequired      = errors.New("name é obrigatório")
	ErrPhoneRequired     = errors.New("phone é obrigatório")
	ErrParentRequired    = errors.New("parent_name é obrigatório")
	ErrBirthDateRequired = errors.New("birth_name é obrigatório")
)

type StudentService struct {
	repo *repository.StudentRepository
}

func NewStudentService(repo *repository.StudentRepository) *StudentService {
	return &StudentService{repo: repo}
}

func (s *StudentService) Create(student *model.Student) error {
	if err := validateStudent(student); err != nil {
		return err
	}

	enrollmentNumber, err := s.generateEnrollmentNumber()
	if err != nil {
		return err
	}
	student.EnrollmentNumber = enrollmentNumber

	return s.repo.Create(student)
}

func (s *StudentService) FindAll() ([]model.Student, error) {
	return s.repo.FindAll()
}

func (s *StudentService) FindByID(id string) (*model.Student, error) {
	return s.repo.FindByID(id)
}

func (s *StudentService) Update(student *model.Student) error {
	if err := validateStudent(student); err != nil {
		return err
	}
	return s.repo.Update(student)
}

func (s *StudentService) Delete(id string) error {
	return s.repo.Delete(id)
}

func (s *StudentService) generateEnrollmentNumber() (string, error) {
	year := time.Now().Year()

	sequential, err := s.repo.NextSequential(year)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%d%03d", year, sequential), nil
}

func validateStudent(student *model.Student) error {
	if student.Name == "" {
		return ErrNameRequired
	}
	if student.ParentName == "" {
		return ErrParentRequired
	}
	if student.Phone == "" {
		return ErrPhoneRequired
	}
	if student.BirthDate.IsZero() {
		return ErrBirthDateRequired
	}
	return nil
}
