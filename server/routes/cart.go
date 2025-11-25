package routes

import (
	"jott55/go-shop/server/serverio"
	"jott55/go-shop/services"
	"jott55/go-shop/types"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func Cart(router *chi.Mux, cart *services.CartService) {
	router.Get("/cart/create", func(w http.ResponseWriter, r *http.Request) {
		cart.Create()
	})

	router.Post("/cart/insert", func(w http.ResponseWriter, r *http.Request) {
		cr, _ := serverio.GetStructFromRequestBody[types.CartRequest](r)
		cart.Insert(cr.Cart)
	})

	router.Get("/cart", func(w http.ResponseWriter, r *http.Request) {
		cs := cart.GetWhere(0, 100)
		serverio.SendJson(w, cs)
	})

	router.Get("/cart/{id}", func(w http.ResponseWriter, r *http.Request) {
		id, _ := serverio.GetId(r)
		c, _ := cart.Get(id)
		serverio.SendJson(w, c)
	})

	router.Get("/cart/{id}/delete", func(w http.ResponseWriter, r *http.Request) {
		id, _ := serverio.GetId(r)
		cart.Delete(id)
	})

}
