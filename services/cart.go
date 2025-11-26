package services

import (
	"fmt"
	"jott55/go-shop/database"
	"jott55/go-shop/types"
)

type CartService struct {
	dl    *database.DatabaseLink
	table string
}

func (c *CartService) GetIdByUserId(user_id int) (int, error) {
	ar := database.GenericGetWhere[types.CartId](c.dl, c.table, fmt.Sprintf("%s=%d", "user_id", user_id))
	if len(ar) >= 1 {
		return ar[0].Id, nil
	}
	return 0, CreateError(NOT_FOUND, "no carts for user of id: %d", user_id)
}

func (c *CartService) Init(dl *database.DatabaseLink, table_name string) {
	c.dl = dl
	c.table = table_name
}

func (c *CartService) Get(id int) (types.Cart, error) {
	return database.GenericGet[types.Cart](c.dl, c.table, id)
}

func (c *CartService) GetWhere(id_min int, id_max int) []types.Cart {
	return database.GenericGetWhere[types.Cart](c.dl, c.table, "true")
}

func (c *CartService) Insert(cart *types.CartNoId) database.DatabaseResponse {
	return c.dl.Insert(c.table, cart)
}

func (c *CartService) Drop() {
	c.dl.DropTable(c.table)
}

func (c *CartService) Create() {

	const (
		const_id      = "id"
		const_user_id = "user_id"
	)

	field := fmt.Sprintf(`
		%s BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
		%s BIGINT REFERENCES %s(id) ON DELETE CASCADE
	`, const_id, const_user_id, "users")

	c.dl.CreateTable(c.table, field)
	c.dl.CreateIndex(c.table, []string{"user_id"})
}

func (c *CartService) Delete(id int) error {
	return c.dl.DeleteById(c.table, id)
}
