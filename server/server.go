package server

import (
	"fmt"
	"jott55/go-shop/clog"
	"jott55/go-shop/database"
	"jott55/go-shop/product"
	"jott55/go-shop/server/routes"
	"jott55/go-shop/user"
	"jott55/go-shop/user/cart"
	"jott55/go-shop/user/cart/item"
	"net/http"
	"os"

	"encoding/json"

	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/docgen"
)

var shopDB *database.DatabaseLink

func noDb(dl *database.DatabaseLink) bool {
	if dl == nil {
		clog.Log(clog.ERROR, "no db connection, returning")
		return true
	}
	return false
}

func checkFileExist(name string) bool {
	file, err := os.Open(name)
	file.Close()
	return os.IsNotExist(err)
}

func checkError(err error, msg ...any) bool {
	if err != nil {
		clog.Logger(clog.ERROR, 2, err, msg)
		return true
	}
	return false
}

func debug(msg ...any) {
	clog.Log(clog.DEBUG, msg...)
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

func createNewConfigFile(name string) {
	file, err := os.Create(name)
	if checkError(err) {
		return
	}

	db := getDatabaseInfoFromUser()

	data, err := json.Marshal(db)

	if checkError(err) {
		return
	}

	file.Write(data)

	file.Close()
}

func getConfigFileData() ([]byte, error) {
	return os.ReadFile("config.json")
}

func configure(dl *database.DatabaseLink) {
	config_filename := "config.json"
	if checkFileExist(config_filename) {
		createNewConfigFile(config_filename)
	}
	dat, err := getConfigFileData()

	if checkError(err) {
		return
	}

	if len(dat) < 64 {
		createNewConfigFile(config_filename)
		configure(dl)
		return
	}

	var db database.DatabaseInfo
	json.Unmarshal(dat, &db)

	database.Configure(dl, db)
}

func mainPage(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello... what are you doing here???? anyway check my discord: @wasenokkami"))
}

func admin(w http.ResponseWriter, r *http.Request) {
	file, err := os.ReadFile("admin.html")

	if checkError(err) {
		return
	}
	clog.Log(clog.DEBUG, "admin page")
	w.Write(file)
}

func generateMarkdown(router *chi.Mux) {
	data := docgen.MarkdownRoutesDoc(router, docgen.MarkdownOpts{})

	err := os.Remove("routes.md")

	if err != nil && os.IsExist(err) {
		clog.Log(clog.ERROR, err)
		return
	}

	file, err := os.Create("route.md")

	if checkError(err) {
		return
	}

	_, err = file.Write([]byte(data))

	if checkError(err) {
		return
	}

	defer file.Close()
}

func getImage(w http.ResponseWriter, r *http.Request) {
	product_name_param := chi.URLParam(r, "name")
	img, err := os.ReadFile(fmt.Sprintf("images/%s", product_name_param))

	if checkError(err) {
		return
	}

	w.Write(img)
}

func doRouterShit() {

	clog.Log(clog.INFO, "initializing router")

	router := chi.NewRouter()

	router.Use(middleware.Logger) // router logger
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

	router.Get("/admin", admin) // admin page

	router.Get("/generate", func(w http.ResponseWriter, r *http.Request) {
		generateMarkdown(router) // chi middleware for creating api doc
		data, _ := os.ReadFile("route.md")
		w.Write(data)
	})

	router.Get("/images/{name}", getImage)

	router.Get("/createAllTables", func(w http.ResponseWriter, r *http.Request) {
		user.CreateTable(shopDB)
		product.CreateTable(shopDB)
		cart.CreateTable(shopDB)
		item.CreateTable(shopDB)
	})

	router.Get("/deleteAllTables", func(w http.ResponseWriter, r *http.Request) {
		item.Drop(shopDB)
		cart.Drop(shopDB)
		product.Drop(shopDB)
		user.Drop(shopDB)
	})

	routes.SetDatabase(shopDB)
	routes.Start(router)

	err := http.ListenAndServe(":8069", router)
	if err != nil {
		clog.Log(clog.ERROR, err)
		return
	}

}

func startDatabase(dl *database.DatabaseLink) {

	err := database.Init(dl) // fill database link

	if err != nil {
		var str string

		clog.Log(clog.ERROR, "Database ERROR")
		fmt.Println("want to continue anyway? (y/n) change config (c)? rerun? (r)")
		fmt.Scanf("%s", &str)
		switch strings.ToLower(str) {
		case "n": // exit program
			os.Exit(1)
		case "c": // create config and retry
			createNewConfigFile("config.json")
			configure(dl) // apply config
			startDatabase(dl)
			return
		case "r": //  retry
			startDatabase(dl)
			return
		case "y": // exit function
			return
		default: // retry
			startDatabase(dl)
		}
	} else {
		fmt.Println("database Initialized")
	}
}

func Run() {

	shopDB = database.Create()

	configure(shopDB)

	startDatabase(shopDB)

	doRouterShit()

	database.Close(shopDB)
}
