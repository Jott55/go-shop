package routes

import "github.com/go-chi/chi/v5"

func Product(router *chi.Mux) {
	router.Get("/create/products", createProducts)

	router.Get("/product", getProductsSimplyfied)

	router.Get("/product/{id}", getProduct)

	router.Post("/product/insert", insertProduct)

	router.Get("/product/{id}/delete", deleteProduct)
}
