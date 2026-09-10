package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/TakedaB/akademi-api/internal/model"
	"github.com/TakedaB/akademi-api/internal/service"
)

type FinanceHandler struct {
	service *service.FinanceService
}

func NewFinanceHandler(service *service.FinanceService) *FinanceHandler {
	return &FinanceHandler{service: service}
}

func (h *FinanceHandler) Create(w http.ResponseWriter, r *http.Request) {
	var finance model.Finance
	if err := json.NewDecoder(r.Body).Decode(&finance); err != nil {
		respondError(w, http.StatusBadRequest, "corpo da requisição inválido")
		return
	}

	if err := h.service.Create(&finance); err != nil {
		respondFinanceError(w, err)
		return
	}

	respondJSON(w, http.StatusCreated, finance)
}

func (h *FinanceHandler) FindAll(w http.ResponseWriter, r *http.Request) {
	records, err := h.service.FindAll()
	if err != nil {
		respondError(w, http.StatusInternalServerError, "erro ao buscar cobranças")
		return
	}

	respondJSON(w, http.StatusOK, records)
}

func (h *FinanceHandler) FindByStudentID(w http.ResponseWriter, r *http.Request) {
	studentID := r.PathValue("studentId")

	records, err := h.service.FindByStudentID(studentID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "erro ao buscar cobranças")
		return
	}

	respondJSON(w, http.StatusOK, records)
}

type updateStatusRequest struct {
	Status string `json:"status"`
}

func (h *FinanceHandler) UpdateStatus(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	var req updateStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "corpo da requisição inválido")
		return
	}

	updated, err := h.service.UpdateStatus(id, req.Status)
	if err != nil {
		respondFinanceError(w, err)
		return
	}

	respondJSON(w, http.StatusOK, updated)
}

func (h *FinanceHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	if err := h.service.Delete(id); err != nil {
		respondError(w, http.StatusNotFound, "cobrança não encontrada")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func respondFinanceError(w http.ResponseWriter, err error) {
	validationErrors := []error{
		service.ErrFinanceStudentIDRequired,
		service.ErrFinanceDescriptionRequired,
		service.ErrFinanceAmountInvalid,
		service.ErrFinancePaymentMethodInvalid,
		service.ErrFinanceStatusInvalid,
		service.ErrFinanceDueDateRequired,
	}

	for _, ve := range validationErrors {
		if errors.Is(err, ve) {
			respondError(w, http.StatusBadRequest, err.Error())
			return
		}
	}

	respondError(w, http.StatusInternalServerError, "erro interno")
}
