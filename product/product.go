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

func Insert(dl *database.DatabaseLink, product *Product) error {

	clog.Log(clog.DEBUG, "Inserting to db")

	clog.Log(clog.DEBUG, "Name: ", product.Name)
	clog.Log(clog.DEBUG, "Image: ", product.Image_url)
	clog.Log(clog.DEBUG, "Price: ", product.Price)
	clog.Log(clog.DEBUG, "Description: ", product.Description)

	var sql_insert string

	if product.Id < 0 {
		sql_insert = fmt.Sprintf("INSERT INTO products (name, image_url, price, description) VALUES ('%v', '%v', %v, '%v')", product.Name, product.Image_url, product.Price, product.Description)
	} else {
		sql_insert = fmt.Sprintf("INSERT INTO products (id, name, image_url, price, description) VALUES ('%v', '%v', '%v', %v, '%v')", product.Id, product.Name, product.Image_url, product.Price, product.Description)
	}

	clog.Log(clog.DEBUG, sql_insert)
	tag, err := database.Exec(dl, sql_insert)

	if err != nil {
		clog.Log(clog.ERROR, err)
		return err
	}
	clog.Log(clog.DEBUG, tag)
	return nil
}

func Delete(dl *database.DatabaseLink, id int) error {
	clog.Log(clog.DEBUG, "Deleting product by id: ", id)
	sql_delete := fmt.Sprintf("DELETE FROM products WHERE id=%v", id)
	clog.Log(clog.DEBUG, sql_delete)

	tag, err := database.Exec(dl, sql_delete)

	if err != nil {
		clog.Log(clog.ERROR, err)
		return err
	}
	clog.Log(clog.DEBUG, tag)
	return nil
}

func Get(dl *database.DatabaseLink, id int) (Product, error) {
	clog.Log(clog.DEBUG, "Getting product name")
	sql_select := fmt.Sprintf("SELECT name, image_url, price, description FROM products WHERE id=%v", id)

	product := Product{Id: id}
	err := database.QueryRow(dl, sql_select, &product.Name, &product.Image_url, &product.Price, &product.Description)

	if err != nil {
		return Product{}, err
	}

	return product, nil
}

func GetAllSimplyfied(dl *database.DatabaseLink, id_min int, id_max int) ([]ProductView, error) {
	clog.Log(clog.DEBUG, "Getting products in range of: ", id_min, " ", id_max)

	sql_select := fmt.Sprintf("SELECT id, name, image_url, price FROM products WHERE id BETWEEN %v AND %v", id_min, id_max)

	rows, err := database.Query(dl, sql_select)

	if err != nil {
		clog.Log(clog.ERROR, err)
		return nil, err
	}

	clog.Log(clog.DEBUG, "Query was a success")

	products, err := database.CollectRows[ProductView](rows)

	if err != nil {
		clog.Log(clog.ERROR, err)
		return nil, err
	}

	clog.Log(clog.DEBUG, "Collect rows a success")

	return products, nil
}

func CreateTable(dl *database.DatabaseLink) {
	sql_table := `CREATE TABLE products (
		id bigint GENERATED ALWAYS AS IDENTITY,
		name VARCHAR(50),
		image_url VARCHAR(255),
		price NUMERIC(10,2),
		description VARCHAR(255)
	)`

	table, err := database.Exec(dl, sql_table)

	if err != nil {
		clog.Log(clog.ERROR, err)
	}
	clog.Log(clog.DEBUG, table)
}
