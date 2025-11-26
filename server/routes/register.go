package routes

import (
	"jott55/go-shop/server/serverio"
	"jott55/go-shop/types"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func Register(router *chi.Mux) {
	router.Post("/register/user", func(w http.ResponseWriter, r *http.Request) {
		request, err := serverio.GetStructFromRequestBody[types.UserRequest](r)

		// check if request is valid
		if checkError(err) ||
			len(request.User.Email) < 10 ||
			len(request.User.Name) < 4 ||
			len(request.User.Password) < 8 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		// Create user
		ser.User.Insert(request.User)
		// Get user id
		user_id, err := ser.User.GetIdByName(request.User.Name)
		checkError(err)

		// Create user cart
		ser.Cart.Insert(&types.CartNoId{User_id: user_id})

		// Get signed token adding username to it
		token := createTokenString(request.User.Name)

		// send token
		w.Write([]byte(token))
	})

	// test your key here
	router.Get("/test/{key}", func(w http.ResponseWriter, r *http.Request) {
		key := chi.URLParam(r, "key")
		claism := decryptTokenString(key)
		debug(claism)
	})
}
