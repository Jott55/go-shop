package routes

import (
	"jott55/go-shop/server/serverio"
	"jott55/go-shop/types"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func Login(router *chi.Mux) {
	router.Post("/login/user", func(w http.ResponseWriter, r *http.Request) {
		luser, err := serverio.GetStructFromRequestBody[types.LoginUser](r)
		if checkError(err) {
			w.WriteHeader(400)
			return
		}

		debug(luser)
	})
}
