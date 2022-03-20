package router

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"stablex/api"
)

func New(h api.OperatorHandler) *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/operators", h.GetOperators)
	r.Get("/operators/{id}", h.FindOperator)
	r.Post("/actions/{id}", h.InsertAction)

	return r
}
