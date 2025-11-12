package main

import (
	"fmt"
	"jott55/go-shop/database"

	"github.com/jackc/pgx/v5"
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

type ProductRequest struct {
	Product Product
}

func productInsert(conn *pgx.Conn, product Product) error {

	clog(DEBUG, "Inserting to db")

	clog(DEBUG, "Name: ", product.Name)
	clog(DEBUG, "Image: ", product.Image_url)
	clog(DEBUG, "Price: ", product.Price)
	clog(DEBUG, "Description: ", product.Description)

	var sql_insert string

	if product.Id < 0 {
		sql_insert = fmt.Sprintf("INSERT INTO products (name, image_url, price, description) VALUES ('%v', '%v', %v, '%v')", product.Name, product.Image_url, product.Price, product.Description)
	} else {
		sql_insert = fmt.Sprintf("INSERT INTO products (id, name, image_url, price, description) VALUES ('%v', '%v', '%v', %v, '%v')", product.Id, product.Name, product.Image_url, product.Price, product.Description)
	}

	clog(DEBUG, sql_insert)
	tag, err := database.Exec(conn, sql_insert)

	if err != nil {
		clog(ERROR, err)
		return err
	}
	clog(DEBUG, tag)
	return nil
}

func productDelete(conn *pgx.Conn, id int) error {
	clog(DEBUG, "Deleting product by id: ", id)
	sql_delete := fmt.Sprintf("DELETE FROM products WHERE id=%v", id)
	clog(DEBUG, sql_delete)

	tag, err := database.Exec(conn, sql_delete)

	if err != nil {
		clog(ERROR, err)
		return err
	}
	clog(DEBUG, tag)
	return nil
}

func productGet(conn *pgx.Conn, id int) (Product, error) {
	clog(DEBUG, "Getting product name")
	sql_select := fmt.Sprintf("SELECT name, image_url, price, description FROM products WHERE id=%v", id)

	product := Product{Id: id}
	err := database.QueryRow(conn, sql_select, &product.Name, &product.Image_url, &product.Price, &product.Description)

	if err != nil {
		return Product{}, err
	}

	return product, nil
}

func productGetAllSimplyfied(conn *pgx.Conn, id_min int, id_max int) ([]ProductView, error) {
	clog(DEBUG, "Getting products in range of: ", id_min, " ", id_max)

	sql_select := fmt.Sprintf("SELECT id, name, image_url, price FROM products WHERE id BETWEEN %v AND %v", id_min, id_max)

	rows, err := database.Query(conn, sql_select)

	if err != nil {
		clog(ERROR, err)
		return nil, err
	}

	clog(DEBUG, "Query was a success")

	products, err := pgx.CollectRows(rows, pgx.RowToStructByName[ProductView])

	if err != nil {
		clog(ERROR, err)
		return nil, err
	}

	clog(DEBUG, "Collect rows a success")

	return products, nil
}
