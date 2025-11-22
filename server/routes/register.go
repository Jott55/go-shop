package routes

import (
	"jott55/go-shop/server/serverio"
	"jott55/go-shop/services/user"
	"jott55/go-shop/types"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func Register(router *chi.Mux) {
	router.Post("/register/user", func(w http.ResponseWriter, r *http.Request) {
		request, err := serverio.GetStructFromRequestBody[types.UserRequest](r)

		if checkError(err) ||
			len(request.User.Email) < 10 ||
			len(request.User.Name) < 4 ||
			len(request.User.Password) < 8 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		user.Insert(dl, request.User)

		token := createTokenString(request.User.Name)

		w.Write([]byte(token))
	})

	router.Get("/test/{key}", func(w http.ResponseWriter, r *http.Request) {

	})
}
