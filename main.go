package main

import (
	"fmt"
	"jott55/go-shop/database"
	"net/http"
	"os"

	"encoding/json"
	"log"

	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/jackc/pgx/v5"
)

var shopConn *pgx.Conn

func main() {
	log.SetPrefix("log: ")
	log.SetFlags(0)

	shopConn = database.Init(database.Database{})

	fmt.Println("database Initialized")

	defer database.Close(shopConn)

	doRouterShit()

}

func getProduct(w http.ResponseWriter, r *http.Request) {
	product_id_param := chi.URLParam(r, "id")

	id, err := strconv.Atoi(product_id_param)

	handleErrorCustom(err, fmt.Sprintf("product get by id: %v", id))

	product := productGet(shopConn, id)
	w.Header().Set("Content-Type", "application/json")
	content, err := json.Marshal(product)
	handleErrorCustom(err, "product json")
	w.Write(content)
}

func getProductsSimplyfied(w http.ResponseWriter, r *http.Request) {
	pd := productGetAllSimplyfied(shopConn, 0, 100000)
	w.Header().Set("Content-Type", "application/json")
	content, err := json.Marshal(pd)
	handleErrorCustom(err, "Json was created")
	w.Write(content)
}

func insertProduct(w http.ResponseWriter, r *http.Request) {

	var product Product

	err := json.NewDecoder(r.Body).Decode(&product)
	handleError(err)

	productInsert(shopConn, product)

}

func deleteProduct(w http.ResponseWriter, r *http.Request) {
	product_id_param := chi.URLParam(r, "id")
	id, err := strconv.Atoi(product_id_param)
	handleError(err)
	productDelete(shopConn, id)
}

func mainPage(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello... what are you doing here???? anyway check my discord: @wasenokkami"))
}

func getImage(w http.ResponseWriter, r *http.Request) {
	product_name_param := chi.URLParam(r, "name")
	img, err := os.ReadFile(fmt.Sprintf("images/%s", product_name_param))
	handleError(err)

	w.Write(img)
}

func admin(w http.ResponseWriter, r *http.Request) {
	file, err := os.ReadFile("admin.html")

	handleErrorCustom(err, "Admin page")
	w.Write(file)
}

func doRouterShit() {
	router := chi.NewRouter()

	router.Use(middleware.Logger)
	router.Use(cors.Handler(
		cors.Options{
			AllowedOrigins:   []string{"https://*", "http://*"},
			AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowedHeaders:   []string{"Accept", "Authorization", "Content-type", "X-CSRF-TOKEN"},
			ExposedHeaders:   []string{"link"},
			AllowCredentials: false,
			MaxAge:           300,
		},
	))

	router.Get("/", mainPage)

	router.Get("/admin", admin)

	router.Get("/product", getProductsSimplyfied)

	router.Post("/post/product", insertProduct)

	router.Get("/images/{name}", getImage)

	router.Get("/product/{id}", getProduct)

	router.Get("/product/{id}/delete", deleteProduct)

	err := http.ListenAndServe(":8069", router)
	handleError(err)
}
