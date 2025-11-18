package product

import (
	"fmt"
	"jott55/go-shop/clog"
	"jott55/go-shop/database"
)

type ProductView struct {
	Id        int
	Name      string
	Image_url string
	Price     int
}

type Product struct {
	Id          int
	Name        string
	Image_url   string
	Price       int
	Description string
}

const Table = "products"

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

func Insert(dl *database.DatabaseLink, product *Product) error {

	debug("Inserting to db")
	debug("Name: ", product.Name)
	debug("Image: ", product.Image_url)
	debug("Price: ", product.Price)
	debug("Description: ", product.Description)

	sql_insert := fmt.Sprintf("INSERT INTO products (name, image_url, price, description) VALUES ('%v', '%v', %v, '%v')", product.Name, product.Image_url, product.Price, product.Description)

	debug(sql_insert)
	tag, err := database.Exec(dl, sql_insert)

	if checkError(err) {
		return err
	}
	debug(tag)
	return nil
}

func Delete(dl *database.DatabaseLink, id int) error {
	debug("Deleting product by id: ", id)
	sql_delete := fmt.Sprintf("DELETE FROM products WHERE id=%v", id)
	debug(sql_delete)

	tag, err := database.Exec(dl, sql_delete)

	if checkError(err) {
		return err
	}
	debug(tag)
	return nil
}

func Get(dl *database.DatabaseLink, id int) (Product, error) {
	debug("Getting product name")
	return database.GenericGet[Product](dl, Table, id)
}

func GetAllSimplyfied(dl *database.DatabaseLink, id_min int, id_max int) ([]ProductView, error) {
	debug("Getting products in range of: ", id_min, " ", id_max)

	sql_select := fmt.Sprintf("SELECT id, name, image_url, price FROM products WHERE id BETWEEN %v AND %v", id_min, id_max)

	rows, err := database.Query(dl, sql_select)

	if checkError(err) {
		return nil, err
	}

	debug("Query was a success")

	products, err := database.CollectRows[ProductView](rows)

	if checkError(err) {
		return nil, err
	}

	debug("Collect rows a success")

	return products, nil
}

func CreateTable(dl *database.DatabaseLink) {
	sql_table := `
		id bigint PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
		name VARCHAR(50),
		image_url VARCHAR(255),
		price NUMERIC(10,2),
		description VARCHAR(255)
	`

	database.CreateTable(dl, Table, sql_table)

}
func Drop(dl *database.DatabaseLink) {
	database.DropTable(dl, Table)
}
