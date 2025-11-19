package product

import (
	"fmt"
	"jott55/go-shop/clog"
	"jott55/go-shop/database"
	"jott55/go-shop/types"
)

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

func Insert(dl *database.DatabaseLink, product *types.Product) error {

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

func Get(dl *database.DatabaseLink, id int) (types.Product, error) {
	debug("Getting product name")
	return database.GenericGet[types.Product](dl, Table, id)
}

func GetAllSimplyfied(dl *database.DatabaseLink, id_min int, id_max int) ([]types.ProductView, error) {
	debug("Getting products in range of: ", id_min, " ", id_max)

	sql_select := fmt.Sprintf("SELECT id, name, image_url, price FROM products WHERE id BETWEEN %v AND %v", id_min, id_max)

	rows, err := database.Query(dl, sql_select)

	if checkError(err) {
		return nil, err
	}

	debug("Query was a success")

	products, err := database.CollectRows[types.ProductView](rows)

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

func GetProductsFromItems(dl *database.DatabaseLink, items []types.ItemNoIdCartId) []types.ProductItem {

	var productItems []types.ProductItem

	for _, it := range items {
		pl, _ := database.GenericGet[types.ProductLess](dl, Table, it.Product_id)
		productItems = append(productItems, types.ProductItem{Id: it.Product_id, Name: pl.Name, Image_url: pl.Image_url, Price: it.Price, Quantity: it.Quantity})
	}

	return productItems
}
