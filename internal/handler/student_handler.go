package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/TakedaB/akademi-api/internal/model"
	"github.com/TakedaB/akademi-api/internal/service"
)

type StudentHandler struct {
	service *service.StudentService
}

func NewStudentHandler(service *service.StudentService) *StudentHandler {
	return &StudentHandler{service: service}
}

func (h *StudentHandler) Create(w http.ResponseWriter, r *http.Request) {
	var student model.Student
	if err := json.NewDecoder(r.Body).Decode(&student); err != nil {
		respondError(w, http.StatusBadRequest, "corpo da requisição inválido")
		return
	}

	if err := h.service.Create(&student); err != nil {
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
