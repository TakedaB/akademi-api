package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/TakedaB/akademi-api/internal/model"
	"github.com/TakedaB/akademi-api/internal/service"
)

type StudentHandler struct {
	service *service.StudentService
}

func NewStudentHandler(service *service.StudentService) *StudentHandler {
	return &StudentHandler{service: service}
}

type createStudentRequest struct {
	Name       string `json:"name"`
	Email      string `json:"email"`
	Password   string `json:"password"`
	BirthDate  string `json:"birth_date"`
	ParentName string `json:"parent_name"`
	City       string `json:"city"`
	Phone      string `json:"phone"`
	Grade      string `json:"grade"`
}

func (h *StudentHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req createStudentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "corpo da requisição inválido")
		return
	}

	birthDate, err := time.Parse(time.RFC3339, req.BirthDate)
	if err != nil {
		respondError(w, http.StatusBadRequest, "birth_date inválida (use formato RFC3339)")
		return
	}

	input := service.CreateStudentInput{
		Name:       req.Name,
		Email:      req.Email,
		Password:   req.Password,
		BirthDate:  birthDate,
		ParentName: req.ParentName,
		City:       req.City,
		Phone:      req.Phone,
		Grade:      req.Grade,
	}

	student, err := h.service.Create(input)
	if err != nil {
		respondValidationOrServerError(w, err)
		return
	}

	respondJSON(w, http.StatusCreated, student)
}

func (h *StudentHandler) FindAll(w http.ResponseWriter, r *http.Request) {
	students, err := h.service.FindAll()
	if err != nil {
		respondError(w, http.StatusInternalServerError, "erro ao buscar alunos")
		return
	}

	respondJSON(w, http.StatusOK, students)
}

func (h *StudentHandler) FindByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	student, err := h.service.FindByID(id)
	if err != nil {
		respondError(w, http.StatusNotFound, "aluno não encontrado")
		return
	}

	respondJSON(w, http.StatusOK, student)
}

func (h *StudentHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	var student model.Student
	if err := json.NewDecoder(r.Body).Decode(&student); err != nil {
		respondError(w, http.StatusBadRequest, "corpo da requisição inválido")
		return
	}

	student.ID = id

	if err := h.service.Update(&student); err != nil {
		respondValidationOrServerError(w, err)
		return
	}

	respondJSON(w, http.StatusOK, student)
}

func (h *StudentHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	if err := h.service.Delete(id); err != nil {
		respondError(w, http.StatusNotFound, "aluno não encontrado")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func respondValidationOrServerError(w http.ResponseWriter, err error) {
	validationErrors := []error{
		service.ErrNameRequired,
		service.ErrEmailRequired,
		service.ErrPasswordRequired,
		service.ErrPhoneRequired,
		service.ErrParentRequired,
		service.ErrBirthDateRequired,
	}

	for _, ve := range validationErrors {
		if errors.Is(err, ve) {
			respondError(w, http.StatusBadRequest, err.Error())
			return
		}
	}

	respondError(w, http.StatusInternalServerError, "erro interno")
}
func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func respondError(w http.ResponseWriter, status int, message string) {
	respondJSON(w, status, map[string]string{"error": message})
}
