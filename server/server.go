package server

import (
	"fmt"
	"jott55/go-shop/clog"
	"jott55/go-shop/database"
	"jott55/go-shop/product"
	"net/http"
	"os"

	"encoding/json"

	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/docgen"
)

type ProductRequest struct {
	Product *product.Product
}

var shopDB *database.DatabaseLink

func Run() {

	shopDB = database.Create()

	configure(shopDB)

	startDatabase(shopDB)

	doRouterShit()

	database.Close(shopDB)
}

func startDatabase(dl *database.DatabaseLink) {

	err := database.Init(dl)

	if err != nil {
		var str string

		clog.Log(clog.ERROR, "Database ERROR")
		fmt.Println("want to continue anyway? (y/n) change config (c)? rerun? (r)")
		fmt.Scanf("%s", &str)
		switch strings.ToLower(str) {
		case "n":
			os.Exit(1)
		case "c":
			createConfigFile("config.json")
			configure(dl)
			startDatabase(dl)
			return
		case "r":
			startDatabase(dl)
			return
		case "y":
			doRouterShit()
		default:
			startDatabase(dl)
		}
	} else {
		fmt.Println("database Initialized")
	}
}

func checkError(err error) {
	if err != nil {
		clog.Log(clog.ERROR, err)
	}
}

func configure(dl *database.DatabaseLink) {
	config_filename := "config.json"
	if checkFileExist(config_filename) {
		createConfigFile(config_filename)
	}
	dat, err := getConfigFileData()
	checkError(err)

	if len(dat) < 64 {
		createConfigFile(config_filename)
	}

	var db database.DatabaseInfo
	json.Unmarshal(dat, &db)

	database.Configure(dl, db)
}

func getConfigFileData() ([]byte, error) {
	return os.ReadFile("config.json")
}

func getDatabaseInfoFromUser() database.DatabaseInfo {
	var db database.DatabaseInfo

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

	if err != nil {
		clog.Log(clog.ERROR, err)
		return
	}

	file.Write(data)

	file.Close()
}

func getProduct(w http.ResponseWriter, r *http.Request) {
	product_id_param := chi.URLParam(r, "id")

	id, err := strconv.Atoi(product_id_param)

	if err != nil {
		clog.Log(clog.ERROR, err, "product of id: ", id)
	}

	if shopDB == nil {
		noDb()
		return
	}

	product, err := product.Get(shopDB, id)

	if err != nil {
		clog.Log(clog.ERROR, err)
		return
	}

	fmt.Println(product)

	w.Header().Set("Content-Type", "application/json")
	content, err := json.Marshal(product)

	if err != nil {
		clog.Log(clog.ERROR, err)
		return
	}
	clog.Log(clog.DEBUG, "product json")
	w.Write(content)
}

func getProductsSimplyfied(w http.ResponseWriter, r *http.Request) {

	if shopDB == nil {
		noDb()
		return
	}

	pd, err := product.GetAllSimplyfied(shopDB, 0, 100000)

	if err != nil {
		clog.Log(clog.ERROR, err)
	}

	w.Header().Set("Content-Type", "application/json")
	content, err := json.Marshal(pd)

	if err != nil {
		clog.Log(clog.ERROR, err)
	}
	clog.Log(clog.DEBUG, "Json was created ")
	w.Write(content)
}

func insertProduct(w http.ResponseWriter, r *http.Request) {

	var pr ProductRequest

	err := json.NewDecoder(r.Body).Decode(&pr)
	if err != nil {
		clog.Log(clog.ERROR, err)
	}

	if shopDB == nil {
		noDb()
		return
	}

	fmt.Println(pr)

	fmt.Println("inserting")

	product.Insert(shopDB, pr.Product)

}

func deleteProduct(w http.ResponseWriter, r *http.Request) {
	product_id_param := chi.URLParam(r, "id")
	id, err := strconv.Atoi(product_id_param)
	if err != nil {
		clog.Log(clog.ERROR, err)
	}

	if shopDB == nil {
		noDb()
		return
	}

	product.Delete(shopDB, id)
}

func mainPage(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello... what are you doing here???? anyway check my discord: @wasenokkami"))
}

func getImage(w http.ResponseWriter, r *http.Request) {
	product_name_param := chi.URLParam(r, "name")
	img, err := os.ReadFile(fmt.Sprintf("images/%s", product_name_param))
	if err != nil {
		clog.Log(clog.ERROR, err)
	}

	w.Write(img)
}

func admin(w http.ResponseWriter, r *http.Request) {
	file, err := os.ReadFile("admin.html")

	if err != nil {
		clog.Log(clog.ERROR, err)
	}
	clog.Log(clog.DEBUG, "admin page")
	w.Write(file)
}

func doRouterShit() {

	clog.Log(clog.INFO, "initializing router")

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

	router.Get("/generate", func(w http.ResponseWriter, r *http.Request) {
		generateMarkdown(router)
	})

	router.Get("/admin", admin)

	router.Get("/product", getProductsSimplyfied)

	router.Post("/post/product", insertProduct)

	router.Get("/images/{name}", getImage)

	router.Get("/product/{id}", getProduct)

	router.Get("/product/{id}/delete", deleteProduct)

	err := http.ListenAndServe(":8069", router)
	if err != nil {
		clog.Log(clog.ERROR, err)
		return
	}

}

func generateMarkdown(router *chi.Mux) {
	data := docgen.MarkdownRoutesDoc(router, docgen.MarkdownOpts{})

	if err := os.Remove("routes.md"); err != nil && os.IsExist(err) {
		clog.Log(clog.ERROR, err)
		return
	}

	file, err := os.Create("route.md")

	if err != nil {
		clog.Log(clog.ERROR, err)
		return
	}

	_, err = file.Write([]byte(data))
	if err != nil {
		clog.Log(clog.ERROR, err)
		return
	}

	defer file.Close()
}

func noDb() {
	clog.Log(clog.ERROR, "no db connection, returning")
}
