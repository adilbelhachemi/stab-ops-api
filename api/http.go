package api

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
	"stablex/domain"
	"time"
)

type OperatorHandler interface {
	FindOperator(http.ResponseWriter, *http.Request)
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

func (h *Handler) FindOperator(w http.ResponseWriter, r *http.Request) {}

func (h *Handler) InsertAction(w http.ResponseWriter, r *http.Request) {
	// Declare a new Person struct.
	var action domain.Action

	id := chi.URLParam(r, "id")
	fmt.Println("--- id: ", id)

	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Try to decode the request body into the struct. If there is an error,
	// respond to the client with the error message and a 400 status code.
	err := json.NewDecoder(r.Body).Decode(&action)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	location, err := time.LoadLocation("America/New_York")
	if err != nil {
		fmt.Println(err)
	}
	action.CreatedAt = time.Now().In(location) // time.Now().UTC().Unix()

	h.operatorService.InsertAction(id, action)

	setupResponse(w, action, 200)
}
