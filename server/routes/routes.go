package routes

import (
	"jott55/go-shop/clog"
	"jott55/go-shop/database"

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
	Cart(router)
	Item(router)
	Login(router)
	Product(router)
	Register(router)
	User(router)
}

func SetDatabase(databaselink *database.DatabaseLink) {
	dl = databaselink
}
