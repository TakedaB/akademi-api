package service

import (
	"errors"

	"github.com/TakedaB/akademi-api/internal/model"
	"github.com/TakedaB/akademi-api/internal/repository"
)

var (
	ErrFinanceStudentIDRequired    = errors.New("student_id é obrigatório")
	ErrFinanceDescriptionRequired  = errors.New("description é obrigatório")
	ErrFinanceAmountInvalid        = errors.New("amount deve ser maior que zero")
	ErrFinancePaymentMethodInvalid = errors.New("payment_method inválido")
	ErrFinanceStatusInvalid        = errors.New("status inválido")
	ErrFinanceDueDateRequired      = errors.New("due_date é obrigatório")
)

var validPaymentMethods = map[string]bool{
	"cash":          true,
	"pix":           true,
	"cartao":        true,
	"boleto":        true,
	"transferencia": true,
}

var validStatuses = map[string]bool{
	"pendente": true,
	"pago":     true,
	"atrasado": true,
	"isento":   true,
}

type FinanceService struct {
	repo *repository.FinanceRepository
}

func NewFinanceService(repo *repository.FinanceRepository) *FinanceService {
	return &FinanceService{repo: repo}
}

func (s *FinanceService) Create(f *model.Finance) error {
	if f.StudentID == "" {
		return ErrFinanceStudentIDRequired
	}
	if f.Description == "" {
		return ErrFinanceDescriptionRequired
	}
	if f.Amount <= 0 {
		return ErrFinanceAmountInvalid
	}
	if !validPaymentMethods[f.PaymentMethod] {
		return ErrFinancePaymentMethodInvalid
	}
	if f.DueDate.IsZero() {
		return ErrFinanceDueDateRequired
	}

	f.Status = "pendente" // Default status

	return s.repo.Create(f)
}

func (s *FinanceService) FindAll() ([]model.Finance, error) {
	return s.repo.FindAll()
}

func (s *FinanceService) FindByStudentID(studentID string) ([]model.Finance, error) {
	return s.repo.FindByStudentID(studentID)
}

func (s *FinanceService) UpdateStatus(id, status string) (*model.Finance, error) {
	if !validStatuses[status] {
		return nil, ErrFinanceStatusInvalid
	}
	return s.repo.UpdateStatus(id, status)
}

func (s *FinanceService) Delete(id string) error {
	return s.repo.Delete(id)
}
