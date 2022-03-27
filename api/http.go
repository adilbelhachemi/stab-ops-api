package api

import (
	"encoding/json"
	"fmt"
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
	Signin(http.ResponseWriter, *http.Request)
	Signup(http.ResponseWriter, *http.Request)
	FindOperator(http.ResponseWriter, *http.Request)
	GetOperators(http.ResponseWriter, *http.Request)
	UpdateOperator(http.ResponseWriter, *http.Request)
	InsertAction(http.ResponseWriter, *http.Request)
}

type Handler struct {
	operatorService domain.OperatorService
}

func NewHandler(operatorService domain.OperatorService) OperatorHandler {
	return &Handler{operatorService: operatorService}
}

func (h *Handler) Signin(w http.ResponseWriter, r *http.Request) {
	opReq := getOperatorSigninRequest(w, r)
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
	helper.SetCookieHandler(w, r, "jwt-token", token)
}

func (h *Handler) Signup(w http.ResponseWriter, r *http.Request) {
	opReq := getOperatorSigninRequest(w, r)
	password := auth.GetHash([]byte(opReq.Password))

	operator, err := h.operatorService.FindOperator(opReq.ID, domain.OperatorFilter{})
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if operator.Password != "" {
		helper.Response(w, "You're already a member", http.StatusConflict)
	}

	opr, err := h.operatorService.UpdateOperator(opReq.ID, domain.OperatorFilter{Password: password})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	helper.Response(w, fmt.Sprintf("Operator %d updated successfully", opr.ID), http.StatusOK)
}

func getOperatorSigninRequest(w http.ResponseWriter, r *http.Request) domain.OperatorSigninRequest {
	var opReq domain.OperatorSigninRequest

	err := json.NewDecoder(r.Body).Decode(&opReq)
	if err != nil {
		helper.Response(w, err.Error(), http.StatusInternalServerError)
	}

	if opReq.ID == "" || opReq.Password == "" {
		helper.Response(w, "Username or password are missing!", http.StatusBadRequest)
	}

	return opReq
}

func (h *Handler) GetOperators(w http.ResponseWriter, r *http.Request) {
	var filter domain.OperatorFilter

	err := json.NewDecoder(r.Body).Decode(&filter)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	res, _ := h.operatorService.GetOperators(filter)
	helper.Response(w, res, 200)
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
	helper.Response(w, operator, 200)
}

func (h *Handler) UpdateOperator(w http.ResponseWriter, r *http.Request) {
	var opFilter domain.OperatorFilter

	id := chi.URLParam(r, "id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err := json.NewDecoder(r.Body).Decode(&opFilter)
	if err != nil {
		helper.Response(w, err.Error(), http.StatusInternalServerError)
	}

	res, err := h.operatorService.UpdateOperator(id, opFilter)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	helper.Response(w, res, http.StatusOK)
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

	helper.Response(w, action, 200)
}
