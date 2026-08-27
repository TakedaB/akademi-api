package service

import (
	"errors"
	"time"

	"github.com/TakedaB/akademi-api/internal/model"
	"github.com/TakedaB/akademi-api/internal/repository"
)

var (
	ErrTeacherNameRequired     = errors.New("name é obrigatório")
	ErrTeacherEmailRequired    = errors.New("email é obrigatório")
	ErrTeacherPasswordRequired = errors.New("password é obrigatório")
	ErrTeacherSubjectRequired  = errors.New("subject é obrigatório")
	ErrTeacherPhoneRequired    = errors.New("phone é obrigatório")
	ErrTeacherClassRequired    = errors.New("class_assigned é obrigatório")
	ErrTeacherHireDateRequired = errors.New("hire_date é obrigatório")
)

type CreateTeacherInput struct {
	Name          string
	Email         string
	Password      string
	Subject       string
	Phone         string
	HireDate      time.Time
	ClassAssigned string
	WorkloadHours int
}

type TeacherService struct {
	teacherRepo *repository.TeacherRepository
	userRepo    *repository.UserRepository
}

func NewTeacherService(teacherRepo *repository.TeacherRepository, userRepo *repository.UserRepository) *TeacherService {
	return &TeacherService{teacherRepo: teacherRepo, userRepo: userRepo}
}

func (s *TeacherService) Create(input CreateTeacherInput) (*model.Teacher, error) {
	if input.Name == "" {
		return nil, ErrTeacherNameRequired
	}
	if input.Email == "" {
		return nil, ErrTeacherEmailRequired
	}
	if input.Password == "" {
		return nil, ErrTeacherPasswordRequired
	}
	if input.Subject == "" {
		return nil, ErrTeacherSubjectRequired
	}
	if input.Phone == "" {
		return nil, ErrTeacherPhoneRequired
	}
	if input.ClassAssigned == "" {
		return nil, ErrTeacherClassRequired
	}
	if input.HireDate.IsZero() {
		return nil, ErrTeacherHireDateRequired
	}

	passwordHash, err := HashPassword(input.Password)
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
		Role:         model.RoleProfessor,
	}
	if err := s.userRepo.Create(tx, user); err != nil {
		return nil, err
	}

	teacher := &model.Teacher{
		UserID:        user.ID,
		Name:          user.Name,
		Email:         user.Email,
		Subject:       input.Subject,
		Phone:         input.Phone,
		HireDate:      input.HireDate,
		ClassAssigned: input.ClassAssigned,
		WorkloadHours: input.WorkloadHours,
	}
	if err := s.teacherRepo.Create(tx, teacher); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return teacher, nil
}

func (s *TeacherService) FindAll() ([]model.Teacher, error) {
	return s.teacherRepo.FindAll()
}

func (s *TeacherService) FindByID(id string) (*model.Teacher, error) {
	return s.teacherRepo.FindByID(id)
}

func (s *TeacherService) Delete(id string) error {
	return s.teacherRepo.Delete(id)
}
