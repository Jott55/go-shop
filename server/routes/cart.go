package routes

import (
	"jott55/go-shop/server/serverio"
	"jott55/go-shop/types"
	"jott55/go-shop/user/cart"
	"jott55/go-shop/user/cart/item"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func Cart(router *chi.Mux) {
	router.Get("/cart/create", func(w http.ResponseWriter, r *http.Request) {
		cart.CreateTable(dl)
	})

	router.Get("/cart/item/create", func(w http.ResponseWriter, r *http.Request) {
		item.CreateTable(dl)
	})

	router.Post("/cart/insert", func(w http.ResponseWriter, r *http.Request) {
		cr, _ := serverio.GetStructFromRequestBody[types.CartRequest](r)
		cart.Insert(dl, cr.Cart)
	})

	router.Get("/cart", func(w http.ResponseWriter, r *http.Request) {
		cs := cart.GetAll(dl)
		serverio.SendJson(w, cs)
	})

	router.Get("/cart/{id}", func(w http.ResponseWriter, r *http.Request) {
		id, _ := serverio.GetId(r)
		c, _ := cart.Get(dl, id)
		serverio.SendJson(w, c)
	})

	router.Get("/cart/{id}/delete", func(w http.ResponseWriter, r *http.Request) {
		id, _ := serverio.GetId(r)
		cart.Delete(dl, id)
	})

}
