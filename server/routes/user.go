package routes

import (
	"encoding/json"
	"jott55/go-shop/server/serverio"
	"jott55/go-shop/services/cart"
	"jott55/go-shop/services/cart_item"
	"jott55/go-shop/services/product"
	"jott55/go-shop/services/user"
	"jott55/go-shop/types"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func getUsers(w http.ResponseWriter, r *http.Request) {
	users := user.GetWhere(dl, 0, 10)

	serverio.SendJson(w, users)
}

func insertUser(w http.ResponseWriter, r *http.Request) {
	var us types.UserRequest

	err := json.NewDecoder(r.Body).Decode(&us)

	if checkError(err) {
		return
	}

	if noDb(dl) {
		return
	}

	user.Insert(dl, us.User)
}

func dropUsers(w http.ResponseWriter, r *http.Request) {
	user.Drop(dl)
}

func deleteUser(w http.ResponseWriter, r *http.Request) {
	id, err := serverio.GetId(r)
	checkError(err)
	debug("deleting user: ", id)

	user.Delete(dl, id)

}

func createUsers(w http.ResponseWriter, r *http.Request) {
	user.CreateTable(dl)
}

func getUser(w http.ResponseWriter, r *http.Request) {

	id, err := serverio.GetId(r)

	if checkError(err) {
		return
	}

	user, err := user.Get(dl, id)

	if checkError(err) {
		return
	}

	serverio.SendJson(w, user)
}

func User(router chi.Router) {
	noDb(dl)

	router.Get("/user/create", createUsers)

	router.Get("/user/drop", dropUsers)

	router.Get("/user/{id}", getUser)

	router.Get("/user", getUsers)

	router.Post("/user/insert", insertUser)

	router.Get("/user/{id}/delete", deleteUser)

	router.Get("/user/cart", func(w http.ResponseWriter, r *http.Request) {

		username := r.Context().Value(username_key).(string)

		debug(username)

		cart_id := cart.GetIdByUserId(dl, 1)
		items := cart_item.GetByCartId(dl, cart_id)

		productsItems := product.GetProductsFromItems(dl, items)
		serverio.SendJson(w, productsItems)
	})
}
