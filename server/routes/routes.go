package routes

import (
	"context"
	"jott55/go-shop/clog"
	"jott55/go-shop/database"
	"net/http"

	"github.com/go-chi/chi/v5"
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

func Start(router *chi.Mux) {

	router.Group(func(r chi.Router) {
	})

	router.Use(someMiddleWare)

	User(router)
	Cart(router)
	Item(router)
	Login(router)
	Product(router)
	Register(router)
}

func SetDatabase(databaselink *database.DatabaseLink) {
	dl = databaselink
}

func someMiddleWare(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), "user", "jotter")

		next.ServeHTTP(w, r.WithContext(ctx))

	})
}

func someHandler(w http.Response, r *http.Request) {
	user := r.Context().Value("user").(string)

	debug(user)
}
