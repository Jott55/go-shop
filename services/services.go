package services

import "jott55/go-shop/database"

type IService interface {
	Init(dl *database.DatabaseLink, table_name string)
	Get(id int) (any, error)
	GetWhere(id_min int, id_max int) []any
	Insert(t *any) database.DatabaseResponse
	Drop()
	Create()
	Delete(id int) error
}

type Services struct {
	Cart      *CartService
	Cart_item *CartItemService
	Product   *ProductService
	User      *UserService
}
