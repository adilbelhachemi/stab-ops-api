package api

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"stablex/auth"
	"stablex/domain"
	"stablex/domain/helper"
	"time"
)

type OperatorHandler interface {
	FindOperator(http.ResponseWriter, *http.Request)
	GetOperators(http.ResponseWriter, *http.Request)
	InsertAction(http.ResponseWriter, *http.Request)
	Signin(http.ResponseWriter, *http.Request)
	Signup(http.ResponseWriter, *http.Request)
}

type Handler struct {
	operatorService domain.OperatorService
}

func NewHandler(operatorService domain.OperatorService) OperatorHandler {
	return &Handler{operatorService: operatorService}
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
	helper.SetupResponse(w, operator, 200)
}

func (h *Handler) GetOperators(w http.ResponseWriter, r *http.Request) {
	var filter domain.OperatorFilter

	err := json.NewDecoder(r.Body).Decode(&filter)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	res, _ := h.operatorService.GetOperators(filter)
	helper.SetupResponse(w, res, 200)
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

	helper.SetupResponse(w, action, 200)
}

func getOperatorRequest(w http.ResponseWriter, r *http.Request) domain.OperatorRequest {
	var opReq domain.OperatorRequest

	err := json.NewDecoder(r.Body).Decode(&opReq)
	if err != nil {
		helper.SetupResponse(w, err.Error(), http.StatusInternalServerError)
	}

	if opReq.ID == "" || opReq.Password == "" {
		helper.SetupResponse(w, "Username or password are missing!", http.StatusBadRequest)
	}

	return opReq
}

func (h *Handler) Signin(w http.ResponseWriter, r *http.Request) {
	opReq := getOperatorRequest(w, r)
	operator, err := h.operatorService.FindOperator(opReq.ID, domain.OperatorFilter{})
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	userPass := []byte(opReq.Password)
	dbPass := []byte(operator.Password)

	passErr := bcrypt.CompareHashAndPassword(dbPass, userPass)

	if passErr != nil {
		log.Println(passErr)
		http.Error(w, "Wrong username or password", http.StatusBadRequest)
		return
	}

	token, err := auth.CreateToken(opReq.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}
	helper.SetupResponse(w, token, 200)
}

func (h *Handler) Signup(w http.ResponseWriter, r *http.Request) {
	opReq := getOperatorRequest(w, r)
	password := auth.GetHash([]byte(opReq.Password))

	if err := h.operatorService.UpdateOperator(opReq.ID, domain.OperatorFilter{Password: password}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}
