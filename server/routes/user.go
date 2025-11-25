package routes

import (
	"encoding/json"
	"jott55/go-shop/server/serverio"
	"jott55/go-shop/services"
	"jott55/go-shop/types"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func getUsers(w http.ResponseWriter, r *http.Request) {
	users := ser.User.GetWhere(0, 10)

	serverio.SendJson(w, users)
}

func insertUser(w http.ResponseWriter, r *http.Request) {
	var us types.UserRequest

	err := json.NewDecoder(r.Body).Decode(&us)

	if checkError(err) {
		return
	}

	ser.User.Insert(us.User)
}

func dropUsers(w http.ResponseWriter, r *http.Request) {
	ser.User.Drop()
}

func deleteUser(w http.ResponseWriter, r *http.Request) {
	id, err := serverio.GetId(r)
	checkError(err)
	debug("deleting user: ", id)

	ser.User.Delete(id)

}

func createUsers(w http.ResponseWriter, r *http.Request) {
	ser.User.Create()
}

func getUser(w http.ResponseWriter, r *http.Request) {

	id, err := serverio.GetId(r)

	if checkError(err) {
		return
	}

	user, err := ser.User.Get(id)

	if checkError(err) {
		return
	}

	serverio.SendJson(w, user)
}

func User(router chi.Router) {

	router.Get("/user/create", createUsers)

	router.Get("/user/drop", dropUsers)

	router.Get("/user/{id}", getUser)

	router.Get("/user", getUsers)

	router.Post("/user/insert", insertUser)

	router.Get("/user/{id}/delete", deleteUser)

	router.Get("/user/cart", func(w http.ResponseWriter, r *http.Request) {

		username := getKey[string](r, username_key)

		debug(username)
		user_id, err := ser.User.GetIdByName(username)

		checkError(err)

		if err.(services.ServiceError).Code() == services.NOT_FOUND {
			ser.Cart.Insert(&types.CartNoId{User_id: user_id})
		}

		cart_id, err := ser.Cart.GetIdByUserId(user_id)

		checkError(err)

		items := ser.Cart_item.GetByCartId(cart_id)

		productsItems := ser.Product.GetProductsFromItems(items)
		serverio.SendJson(w, productsItems)
	})

	router.Get("/user/item/add", func(w http.ResponseWriter, r *http.Request) {
		// Get product id and cart id

		// Check existent item by product id and cart id

		// if exists sum
	})
}
