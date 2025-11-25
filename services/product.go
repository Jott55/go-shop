package services

import (
	"fmt"
	"jott55/go-shop/database"
	"jott55/go-shop/types"
)

type ProductService struct {
	table string
	dl    *database.DatabaseLink
}

func (p *ProductService) Init(dl *database.DatabaseLink, table_name string) {
	p.dl = dl
	p.table = table_name
}

func (p *ProductService) Get(id int) (types.Product, error) {
	return database.GenericGet[types.Product](p.dl, p.table, id)
}

func (p *ProductService) GetWhere(id_min int, id_max int) []types.Product {
	return database.GenericGetWhere[types.Product](p.dl, p.table, fmt.Sprintf("id BETWEEN %v AND %v", id_min, id_max))
}

func (p *ProductService) Insert(product *types.ProductNoId) database.DatabaseResponse {
	return p.dl.Insert(p.table, product)
}
func (p *ProductService) Drop() {
	p.dl.DropTable(p.table)
}

func (p *ProductService) Create() {
	sql_table := `
		id bigint PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
		name VARCHAR(50),
		image_url VARCHAR(255),
		price NUMERIC(10,2),
		description VARCHAR(255)
	`

	p.dl.CreateTable(p.table, sql_table)
}

func (p *ProductService) Delete(id int) {
	p.dl.DeleteById(p.table, id)
}

func (p *ProductService) GetProductsFromItems(items []types.ItemNoIdCartId) []types.ProductItem {

	var productItems []types.ProductItem

	for _, it := range items {
		pl, _ := database.GenericGet[types.ProductLess](p.dl, p.table, it.Product_id)
		productItems = append(productItems, types.ProductItem{Id: it.Product_id, Name: pl.Name, Image_url: pl.Image_url, Price: it.Price, Quantity: it.Quantity})
	}

	return productItems
}
