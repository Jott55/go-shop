package routes

import (
	"jott55/go-shop/server/serverio"
	"jott55/go-shop/services/cart_item"
	"jott55/go-shop/types"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func Item(router *chi.Mux) {
	router.Post("/item/insert", func(w http.ResponseWriter, r *http.Request) {
		ir, _ := serverio.GetStructFromRequestBody[types.ItemRequest](r)
		cart_item.Insert(dl, ir.Item)
	})

	router.Get("/item/{id}/delete", func(w http.ResponseWriter, r *http.Request) {
		id, _ := serverio.GetId(r)
		cart_item.Delete(dl, id)
	})

	router.Get("/item", func(w http.ResponseWriter, r *http.Request) {
		is := cart_item.GetAll(dl)
		serverio.SendJson(w, is)
	})

	router.Get("/item/{id}", func(w http.ResponseWriter, r *http.Request) {
		id, _ := serverio.GetId(r)
		i, _ := cart_item.Get(dl, id)
		serverio.SendJson(w, i)
	})

	router.Get("/cart/item/create", func(w http.ResponseWriter, r *http.Request) {
		cart_item.CreateTable(dl)
	})

}
