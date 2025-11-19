package routes

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func Register(router *chi.Mux) {
	router.Post("/register/user", func(w http.ResponseWriter, r *http.Request) {
		// ruser, err := insert[types.User](r)
	})
}
