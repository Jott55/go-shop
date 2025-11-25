package routes

import (
	"encoding/json"
	"fmt"
	"jott55/go-shop/clog"

	"jott55/go-shop/types"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func getProduct(w http.ResponseWriter, r *http.Request) {
	product_id_param := chi.URLParam(r, "id")

	id, err := strconv.Atoi(product_id_param)

	if err != nil {
		clog.Log(clog.ERROR, err, "product of id: ", id)
	}

	if checkError(err, "product of id: ", id) {
		return
	}

	product, err := ser.Product.Get(id)

	if err != nil {
		clog.Log(clog.ERROR, err)
		return
	}

	fmt.Println(product)

	w.Header().Set("Content-Type", "application/json")
	content, err := json.Marshal(product)

	if checkError(err) {
		return
	}

	clog.Log(clog.DEBUG, "product json")
	w.Write(content)
}

func insertProduct(w http.ResponseWriter, r *http.Request) {

	var pr types.ProductRequest

	err := json.NewDecoder(r.Body).Decode(&pr)

	if checkError(err) {
		return
	}

	fmt.Println(pr)

	fmt.Println("inserting")

	ser.Product.Insert(pr.Product)

}

func createProducts(w http.ResponseWriter, r *http.Request) {
	ser.Product.Create()
}

func getProductsSimplyfied(w http.ResponseWriter, r *http.Request) {

	pd := ser.Product.GetWhere(0, 100)

	w.Header().Set("Content-Type", "application/json")
	content, err := json.Marshal(pd)

	if checkError(err) {
		return
	}

	clog.Log(clog.DEBUG, "Json was created ")
	w.Write(content)
}

func deleteProduct(w http.ResponseWriter, r *http.Request) {
	product_id_param := chi.URLParam(r, "id")
	id, err := strconv.Atoi(product_id_param)
	if checkError(err) {
		return
	}

	ser.Product.Delete(id)
}

func Product(router *chi.Mux) {
	router.Get("/create/products", createProducts)

	router.Get("/product", getProductsSimplyfied)

	router.Get("/product/{id}", getProduct)

	router.Post("/product/insert", insertProduct)

	router.Get("/product/{id}/delete", deleteProduct)
}
