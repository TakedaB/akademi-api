package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/TakedaB/akademi-api/internal/service"
)

type TeacherHandler struct {
	service *service.TeacherService
}

func NewTeacherHandler(service *service.TeacherService) *TeacherHandler {
	return &TeacherHandler{service: service}
}

type createTeacherRequest struct {
	Name          string `json:"name"`
	Email         string `json:"email"`
	Password      string `json:"password"`
	Subject       string `json:"subject"`
	Phone         string `json:"phone"`
	HireDate      string `json:"hire_date"`
	ClassAssigned string `json:"class_assigned"`
	WorkloadHours int    `json:"workload_hours"`
}

func (h *TeacherHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req createTeacherRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "corpo da requisição inválido")
		return
	}

	hireDate, err := time.Parse(time.RFC3339, req.HireDate)
	if err != nil {
		respondError(w, http.StatusBadRequest, "hire_date inválido, use o formato RFC3339")
		return
	}

	input := service.CreateTeacherInput{
		Name:          req.Name,
		Email:         req.Email,
		Password:      req.Password,
		Subject:       req.Subject,
		Phone:         req.Phone,
		HireDate:      hireDate,
		ClassAssigned: req.ClassAssigned,
		WorkloadHours: req.WorkloadHours,
	}

	teacher, err := h.service.Create(input)
	if err != nil {
		respondTeacherError(w, err)
		return
	}

	respondJSON(w, http.StatusCreated, teacher)
}

func (h *TeacherHandler) FindAll(w http.ResponseWriter, r *http.Request) {
	teachers, err := h.service.FindAll()
	if err != nil {
		respondError(w, http.StatusInternalServerError, "erro ao buscar professores")
		return
	}

	respondJSON(w, http.StatusOK, teachers)
}

func (h *TeacherHandler) FindByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	teacher, err := h.service.FindByID(id)
	if err != nil {
		respondError(w, http.StatusNotFound, "professor não encontrado")
		return
	}

	respondJSON(w, http.StatusOK, teacher)
}

func (h *TeacherHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	if err := h.service.Delete(id); err != nil {
		respondError(w, http.StatusNotFound, "professor não encontrado")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func respondTeacherError(w http.ResponseWriter, err error) {
	validationErrors := []error{
		service.ErrTeacherNameRequired,
		service.ErrTeacherEmailRequired,
		service.ErrTeacherPasswordRequired,
		service.ErrTeacherSubjectRequired,
		service.ErrTeacherPhoneRequired,
		service.ErrTeacherClassRequired,
		service.ErrTeacherHireDateRequired,
	}

	for _, ve := range validationErrors {
		if errors.Is(err, ve) {
			respondError(w, http.StatusBadRequest, err.Error())
			return
		}
	}

	respondError(w, http.StatusInternalServerError, "erro interno")
}
