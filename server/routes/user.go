package routes

import (
	"encoding/json"
	"fmt"
	"jott55/go-shop/server/serverio"
	"jott55/go-shop/services"
	"jott55/go-shop/types"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func getUsers(w http.ResponseWriter, r *http.Request) {
	// Get username from token
	username := getKey[string](r, username_key)
	debug("user name: ", username)

	user, err := ser.User.GetProfileByName(username)
	checkError(err)

	serverio.SendJson(w, user)
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

		if checkError(err) {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		cart_id, err := ser.Cart.GetIdByUserId(user_id)

		if checkError(err) {
			if err.(services.ServiceError).Code() == services.NOT_FOUND {
				w.WriteHeader(http.StatusNotFound)
				return
			}
		}

		items := ser.Cart_item.GetByCartId(cart_id)

		productsItems := ser.Product.GetProductsFromItems(items)
		serverio.SendJson(w, productsItems)
	})

	router.Post("/user/item/add", func(w http.ResponseWriter, r *http.Request) {
		var (
			err        error
			pr         types.ProductIdRequest
			product_id int
			product    types.Product
			cart_id    int
			user_id    int
			username   string
			item       types.ItemNoId
			itemless   *types.ItemNoIdCartIdProductId
		)
		// Get product id from request body
		pr, err = serverio.GetStructFromRequestBody[types.ProductIdRequest](r)
		checkError(err)
		product_id = pr.Product_id

		// Get user id
		username = getKey[string](r, username_key)
		debug("user", username)
		user_id, err = ser.User.GetIdByName(username)
		checkError(err)

		// Get cart id
		cart_id, err = ser.Cart.GetIdByUserId(user_id)
		checkError(err)

		// Check existent item by product id and cart id
		itemless, err = ser.Cart_item.GetByCartIdProductId(cart_id, product_id)
		checkError(err)

		// Get product price
		product, err = ser.Product.Get(product_id)
		checkError(err)

		// if exists sum quantity
		if itemless != nil && itemless.Price == product.Price {
			// update item
			checkError(fmt.Errorf("not implemented yet"))
			w.WriteHeader(http.StatusNotImplemented)
			return
		}
		// else create new item

		// insert product id in new item
		item.Product_id = product_id
		// insert quantity in new item
		item.Quantity = 1
		// insert cart_id in new item
		item.Cart_id = cart_id
		// insert price from product in new item
		item.Price = product.Price
		// insert item on item table
		ser.Cart_item.Insert(&item)
		w.WriteHeader(http.StatusCreated)
	})
}
