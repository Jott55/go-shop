package main

import (
	"fmt"
	"jott55/go-shop/database"
	"log/slog"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
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

func productInsert(conn *pgx.Conn, product Product) {

	fmt.Println("Inserting to db")

	fmt.Println("Name: ", product.Name)
	fmt.Println("Image: ", product.Image_url)
	fmt.Println("Price: ", product.Price)
	fmt.Println("Description: ", product.Description)

	var sql_insert string

	if product.Id < 0 {
		sql_insert = fmt.Sprintf("INSERT INTO products (name, image_url, price, description) VALUES ('%v', '%v', %v, '%v')", product.Name, product.Image_url, product.Price, product.Description)
	} else {
		sql_insert = fmt.Sprintf("INSERT INTO products (id, name, image_url, price, description) VALUES ('%v', '%v', '%v', %v, '%v')", product.Id, product.Name, product.Image_url, product.Price, product.Description)
	}

	handleSql(sql_insert)

	handleResponse(database.Exec(conn, sql_insert))
}

func productDelete(conn *pgx.Conn, id int) {
	fmt.Println("Deleting product by id: ", id)
	sql_delete := fmt.Sprintf("DELETE FROM products WHERE id=%v", id)
	handleSql(sql_delete)

	handleResponse(database.Exec(conn, sql_delete))
}

func productGet(conn *pgx.Conn, id int) Product {
	fmt.Println("Getting product name")
	sql_select := fmt.Sprintf("SELECT name, image_url, price, description FROM products WHERE id=%v", id)

	product := Product{}
	handleError(database.QueryRow(conn, sql_select, &product.Name, &product.Image_url, &product.Price, &product.Description))

	return product
}

func productGetAllSimplyfied(conn *pgx.Conn, id_min int, id_max int) []ProductView {
	fmt.Println("Getting products in range of: ", id_min, " ", id_max)

	sql_select := fmt.Sprintf("SELECT id, name, image_url, price FROM products WHERE id BETWEEN %v AND %v", id_min, id_max)

	rows, err := database.Query(conn, sql_select)

	handleErrorCustom(err, "Querry was a success")

	products, err := pgx.CollectRows(rows, pgx.RowToStructByName[ProductView])
	handleErrorCustom(err, "Collect rows was a success")

	return products
}

func handleResponse(res pgconn.CommandTag, err error) {
	handleError(err)
	fmt.Println(res)
}

func handleError(err error) {
	if err != nil {
		slog.Error(err.Error())
	}
}

func handleErrorCustom(err error, message string) {
	handleError(err)
	fmt.Println(message)
}

func handleSql(sql string) {
	fmt.Println("Sql syntax: ", sql)
}
