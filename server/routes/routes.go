package routes

import (
	"context"
	"fmt"
	"jott55/go-shop/clog"
	"jott55/go-shop/database"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type key int

const (
	username_key key = iota
)

var dl *database.DatabaseLink

func checkError(err error, msg ...any) bool {
	if err != nil {
		clog.Logger(clog.ERROR, 2, err, msg)
		return true
	}
	return false
}

func noDb(dl *database.DatabaseLink) bool {
	if dl == nil {
		clog.Log(clog.ERROR, "no db connection, returning")
		return true
	}
	return false
}
func debug(msg ...any) {
	clog.Log(clog.DEBUG, msg...)
}

func LoginMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		authHeader := r.Header.Get("Authorization")

		token, err := getTokenFromHeader(authHeader)

		checkError(err)

		debug(token)
		claims := decryptTokenString(token)

		if claims == nil {
			w.WriteHeader(http.StatusForbidden)
			return
		}

		user_id := claims["sub"].(string)

		ctx := context.WithValue(r.Context(), username_key, user_id)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func Start(mux *chi.Mux) {

	mux.Group(func(router chi.Router) {
		router.Use(LoginMiddleware)

		router.Get("/viponly", func(w http.ResponseWriter, r *http.Request) {

			s := fmt.Sprintf("You made it! %v", r.Context().Value(username_key))

			w.Write([]byte(s))
		})
		User(router)
	})

	Cart(mux)
	Item(mux)
	Login(mux)
	Product(mux)
	Register(mux)
}

func SetDatabase(databaselink *database.DatabaseLink) {
	dl = databaselink
}
