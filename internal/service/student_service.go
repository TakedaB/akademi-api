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
	ErrEmailRequired     = errors.New("email é obrigatório")
	ErrPasswordRequired  = errors.New("password é obrigatório")
	ErrPhoneRequired     = errors.New("phone é obrigatório")
	ErrParentRequired    = errors.New("parent_name é obrigatório")
	ErrBirthDateRequired = errors.New("birth_date é obrigatório")
)

type CreateStudentInput struct {
	Name       string
	Email      string
	Password   string
	BirthDate  time.Time
	ParentName string
	City       string
	Phone      string
	Grade      string
}

type StudentService struct {
	studentRepo *repository.StudentRepository
	userRepo    *repository.UserRepository
}

func NewStudentService(studentRepo *repository.StudentRepository, userRepo *repository.UserRepository) *StudentService {
	return &StudentService{studentRepo: studentRepo, userRepo: userRepo}
}

func (s *StudentService) Create(input CreateStudentInput) (*model.Student, error) {
	if input.Name == "" {
		return nil, ErrNameRequired
	}
	if input.Email == "" {
		return nil, ErrEmailRequired
	}
	if input.Password == "" {
		return nil, ErrPasswordRequired
	}
	if input.ParentName == "" {
		return nil, ErrParentRequired
	}
	if input.Phone == "" {
		return nil, ErrPhoneRequired
	}
	if input.BirthDate.IsZero() {
		return nil, ErrBirthDateRequired
	}

	passwordHash, err := HashPassword(input.Password)
	if err != nil {
		return nil, err
	}

	enrollmentNumber, err := s.generateEnrollmentNumber()
	if err != nil {
		return nil, err
	}

	tx, err := s.userRepo.BeginTx()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	user := &model.User{
		Name:         input.Name,
		Email:        input.Email,
		PasswordHash: passwordHash,
		Role:         model.RoleAluno,
	}
	if err := s.userRepo.Create(tx, user); err != nil {
		return nil, err
	}

	student := &model.Student{
		UserID:           user.ID,
		Name:             user.Name,
		Email:            user.Email,
		EnrollmentNumber: enrollmentNumber,
		BirthDate:        input.BirthDate,
		ParentName:       input.ParentName,
		City:             input.City,
		Phone:            input.Phone,
		Grade:            input.Grade,
	}
	if err := s.studentRepo.Create(tx, student); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return student, nil
}

func (s *StudentService) FindAll() ([]model.Student, error) {
	return s.studentRepo.FindAll()
}

func (s *StudentService) FindByID(id string) (*model.Student, error) {
	return s.studentRepo.FindByID(id)
}

func (s *StudentService) Update(student *model.Student) error {
	if student.ParentName == "" {
		return ErrParentRequired
	}
	if student.Phone == "" {
		return ErrPhoneRequired
	}
	if student.BirthDate.IsZero() {
		return ErrBirthDateRequired
	}
	return s.studentRepo.Update(student)
}

func (s *StudentService) Delete(id string) error {
	return s.studentRepo.Delete(id)
}

func (s *StudentService) generateEnrollmentNumber() (string, error) {
	year := time.Now().Year()
	sequential, err := s.studentRepo.NextSequential(year)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%d%03d", year, sequential), nil
}
