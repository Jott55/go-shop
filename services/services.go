package services

import (
	"fmt"
	"jott55/go-shop/database"
)

const (
	NOT_FOUND int = iota
)

type IServiceError interface {
	Error()
	Code() int
}

type ServiceError struct {
	message string
	code    int
}

func (e ServiceError) Error() string {
	return e.message
}
func (e ServiceError) Code() int {
	return e.code
}

func CreateError(code int, format string, f ...any) error {
	return ServiceError{message: fmt.Sprintf(format, f...), code: code}
}

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

func CreateServices(dal *database.DatabaseLink, cart string, cart_item string, product string, user string) Services {
	var serve Services

	serve.Cart = &CartService{table: cart, dl: dal}
	serve.Cart_item = &CartItemService{table: cart_item, dl: dal}
	serve.Product = &ProductService{table: product, dl: dal}
	serve.User = &UserService{table: user, dl: dal}
	return serve
}
