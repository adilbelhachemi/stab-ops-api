package router

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"stablex/api"
)

func New(h api.OperatorHandler) *chi.Mux {

	// login := auth.JWT{}.New()

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	// r.Use(login.Authenticator())

	r.Group(func(mux chi.Router) {
		// mux.Use(login.Authenticator())

		r.Route("/operators", func(r chi.Router) {
			r.Get("/", h.GetOperators)
			r.Get("/{id}", h.FindOperator)
			r.Put("/{id}", h.UpdateOperator)
			r.Post("/signin", h.Signin)
			r.Post("/signup", h.Signup)
		})
	})

	r.Post("/actions/{id}", h.InsertAction)

	return r
}
