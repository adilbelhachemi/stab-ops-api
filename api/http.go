package api

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
	"stablex/domain"
	"time"
)

type OperatorHandler interface {
	FindOperator(http.ResponseWriter, *http.Request)
	GetOperators(http.ResponseWriter, *http.Request)
	InsertAction(http.ResponseWriter, *http.Request)
}

type Handler struct {
	operatorService domain.OperatorService
}

func NewHandler(operatorService domain.OperatorService) OperatorHandler {
	return &Handler{operatorService: operatorService}
}

func setupResponse(w http.ResponseWriter, input interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	jsonResp, err := json.Marshal(input)
	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
	}
	w.Write(jsonResp)
	return
}

func (h *Handler) FindOperator(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	operator, err := h.operatorService.FindOperator(id, domain.OperatorFilter{})
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	setupResponse(w, operator, 200)
}

func (h *Handler) GetOperators(w http.ResponseWriter, r *http.Request) {
	var filter domain.OperatorFilter

	err := json.NewDecoder(r.Body).Decode(&filter)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	res, _ := h.operatorService.GetOperators(filter)
	setupResponse(w, res, 200)
}

func (h *Handler) InsertAction(w http.ResponseWriter, r *http.Request) {
	var action domain.Action

	id := chi.URLParam(r, "id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err := json.NewDecoder(r.Body).Decode(&action)
	if err != nil || action.Type == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	action.CreatedAt = time.Now().UTC()
	_ = h.operatorService.InsertAction(id, action)

	setupResponse(w, action, 200)
}
