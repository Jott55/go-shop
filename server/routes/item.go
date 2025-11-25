package routes

import (
	"jott55/go-shop/server/serverio"
	"jott55/go-shop/types"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func Item(router *chi.Mux) {
	router.Post("/item/insert", func(w http.ResponseWriter, r *http.Request) {
		ir, _ := serverio.GetStructFromRequestBody[types.ItemRequest](r)
		ser.Cart_item.Insert(ir.Item)
	})

	router.Get("/item/{id}/delete", func(w http.ResponseWriter, r *http.Request) {
		id, _ := serverio.GetId(r)
		ser.Cart_item.Delete(id)
	})

	router.Get("/item", func(w http.ResponseWriter, r *http.Request) {
		is := ser.Cart_item.GetWhere(0, 100)
		serverio.SendJson(w, is)
	})

	router.Get("/item/{id}", func(w http.ResponseWriter, r *http.Request) {
		id, _ := serverio.GetId(r)
		i, _ := ser.Cart_item.Get(id)
		serverio.SendJson(w, i)
	})

	router.Get("/cart/item/create", func(w http.ResponseWriter, r *http.Request) {
		ser.Cart_item.Create()
	})

}
