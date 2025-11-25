package services

import (
	"fmt"
	"jott55/go-shop/database"
	"jott55/go-shop/types"
)

type CartItemService struct {
	dl    *database.DatabaseLink
	table string
}

func (ci *CartItemService) GetByCartId(cart_id int) []types.ItemNoIdCartId {
	return database.GenericGetWhere[types.ItemNoIdCartId](ci.dl, ci.table, fmt.Sprintf("%s=%d", "cart_id", cart_id))
}

func (ci *CartItemService) Init(dl *database.DatabaseLink, table_name string) {
	ci.dl = dl
	ci.table = table_name
}
func (ci *CartItemService) Get(id int) (types.Item, error) {
	return database.GenericGet[types.Item](ci.dl, ci.table, id)
}
func (ci *CartItemService) GetWhere(id_min int, id_max int) []types.Item {
	return database.GenericGetWhere[types.Item](ci.dl, ci.table, "true")
}
func (ci *CartItemService) Insert(item types.Item) database.DatabaseResponse {
	return ci.dl.Insert(ci.table, item)
}

func (ci *CartItemService) Drop() {
	ci.dl.DropTable(ci.table)
}

func (ci *CartItemService) Create() {

	const (
		const_id         = "id"
		const_cart_id    = "cart_id"
		const_product_id = "product_id"
		const_quantity   = "quantity"
		const_price      = "price"
	)
	field := fmt.Sprintf(`
		%s bigint GENERATED ALWAYS AS IDENTITY,
		%s BIGINT NOT NULL REFERENCES %v(id) ON DELETE CASCADE,
		%s BIGINT NOT NULL REFERENCES %v(id) ON DELETE RESTRICT,
		%s INTEGER NOT NULL CHECK (quantity > 0),
		%s NUMERIC(10, 2) NOT NULL,
		UNIQUE (cart_id, product_id)
	`, const_id, const_cart_id, "cart", const_product_id, "products", const_quantity, const_price)
	// TODO: Change constants and table names
	ci.dl.CreateTable(ci.table, field)
	ci.dl.CreateIndex(ci.table, []string{const_cart_id})
}

func (ci *CartItemService) Delete(id int) error {
	return ci.dl.DeleteById(ci.table, id)
}
