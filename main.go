package main

import (
	"fmt"
	"jott55/go-shop/database"
	"net/http"
	"os"

	"encoding/json"
	"log/slog"

	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/jackc/pgx/v5"
)

var shopDB *pgx.Conn

func main() {

	startDatabase()

	defer database.Close(shopDB)

	doRouterShit()

}

func startDatabase() {
	dbinfo := configure()

	var err error

	shopDB, err = database.Init(dbinfo)

	if err != nil {
		var str string
		slog.Error(err.Error())
		fmt.Println("want to continue anyway? (y/n) change config (c)? rerun? (r)")
		fmt.Scanf("%s", &str)
		if strings.ToLower(str) == "n" {
			os.Exit(1)
		}
		switch strings.ToLower(str) {
		case "n":
			os.Exit(1)
		case "c":
			createConfigFile("config.json")
			startDatabase()
			return
		case "r":
			startDatabase()
			return
		default:
			doRouterShit()
		}
	} else {
		fmt.Println("database Initialized")
	}
}

func checkError(err error) {
	if err != nil {
		slog.Error(err.Error())
	}
}

func configure() database.Database {
	config_filename := "config.json"
	if checkFileExist(config_filename) {
		createConfigFile(config_filename)
	}
	dat, err := getConfigFileData()
	checkError(err)

	if len(dat) < 64 {
		createConfigFile(config_filename)
	}

	var db database.Database
	json.Unmarshal(dat, &db)

	return db
}

func getConfigFileData() ([]byte, error) {
	return os.ReadFile("config.json")
}

func getDatabaseInfoFromUser() database.Database {
	var db database.Database

	fmt.Print("\nEnter user: ")
	fmt.Scanf("%s", &db.User)
	fmt.Print("\nEnter password: ")
	fmt.Scanf("%s", &db.Password)
	fmt.Print("\nEnter host: ")
	fmt.Scanf("%s", &db.Host)
	fmt.Print("\nEnter port: ")
	fmt.Scanf("%s", &db.Port)
	fmt.Print("\nEnter database name: ")
	fmt.Scanf("%s", &db.Database)

	return db
}

func checkFileExist(name string) bool {
	file, err := os.Open(name)
	file.Close()
	return os.IsNotExist(err)
}

func createConfigFile(name string) {
	file, err := os.Create(name)
	checkError(err)

	db := getDatabaseInfoFromUser()

	data, err := json.Marshal(db)
	checkError(err)

	file.Write(data)

	file.Close()
}

func getProduct(w http.ResponseWriter, r *http.Request) {
	product_id_param := chi.URLParam(r, "id")

	id, err := strconv.Atoi(product_id_param)

	if err != nil {
		slog.Error(err.Error(), "product of id: ", id)
	}

	if shopDB == nil {
		noDb()
		return
	}

	product := productGet(shopDB, id)
	w.Header().Set("Content-Type", "application/json")
	content, err := json.Marshal(product)
	handleErrorCustom(err, "product json")
	w.Write(content)
}

func getProductsSimplyfied(w http.ResponseWriter, r *http.Request) {

	if shopDB == nil {
		noDb()
		return
	}

	pd := productGetAllSimplyfied(shopDB, 0, 100000)
	w.Header().Set("Content-Type", "application/json")
	content, err := json.Marshal(pd)
	handleErrorCustom(err, "Json was created")
	w.Write(content)
}

func insertProduct(w http.ResponseWriter, r *http.Request) {

	var product ProductRequest

	err := json.NewDecoder(r.Body).Decode(&product)
	handleError(err)

	if shopDB == nil {
		noDb()
		return
	}

	fmt.Println(product)

	fmt.Println("inserting")

	productInsert(shopDB, product.Product)

}

func deleteProduct(w http.ResponseWriter, r *http.Request) {
	product_id_param := chi.URLParam(r, "id")
	id, err := strconv.Atoi(product_id_param)
	handleError(err)

	if shopDB == nil {
		noDb()
		return
	}

	productDelete(shopDB, id)
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

	slog.Info("initializing router")

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
	if err != nil {
		slog.Error(err.Error())
	}

}

func noDb() {
	slog.Error("no db connection, returning")
}
