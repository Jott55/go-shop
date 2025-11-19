package routes

import (
	"jott55/go-shop/server/serverio"
	"jott55/go-shop/types"
	"jott55/go-shop/user/cart/item"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func Item(router *chi.Mux) {
	router.Post("/item/insert", func(w http.ResponseWriter, r *http.Request) {
		ir, _ := serverio.Insert[types.ItemRequest](r)
		item.Insert(dl, ir.Item)
	})

	router.Get("/item/{id}/delete", func(w http.ResponseWriter, r *http.Request) {
		id, _ := serverio.GetId(r)
		item.Delete(dl, id)
	})

	router.Get("/item", func(w http.ResponseWriter, r *http.Request) {
		is := item.GetAll(dl)
		serverio.SendJson(w, is)
	})

	router.Get("/item/{id}", func(w http.ResponseWriter, r *http.Request) {
		id, _ := serverio.GetId(r)
		i, _ := item.Get(dl, id)
		serverio.SendJson(w, i)
	})

}
