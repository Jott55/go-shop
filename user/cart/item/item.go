package item

import (
	"fmt"
	"jott55/go-shop/database"
	"jott55/go-shop/product"
	"jott55/go-shop/types"
	"jott55/go-shop/user/cart"
)

const Table = "cart_item"

const (
	const_id         = "id"
	const_cart_id    = "cart_id"
	const_product_id = "product_id"
	const_quantity   = "quantity"
	const_price      = "price"
)

func CreateTable(dl *database.DatabaseLink) {
	field := fmt.Sprintf(`
		%s bigint GENERATED ALWAYS AS IDENTITY,
		%s BIGINT NOT NULL REFERENCES %v(id) ON DELETE CASCADE,
		%s BIGINT NOT NULL REFERENCES %v(id) ON DELETE RESTRICT,
		%s INTEGER NOT NULL CHECK (quantity > 0),
		%s NUMERIC(10, 2) NOT NULL,
		UNIQUE (cart_id, product_id)
	`, const_id, const_cart_id, cart.Table, const_product_id, product.Table, const_quantity, const_price)

	database.CreateTable(dl, Table, field)
	database.CreateIndex(dl, Table, []string{const_cart_id})
}

func Drop(dl *database.DatabaseLink) {
	database.DropTable(dl, Table)
}

func Insert(dl *database.DatabaseLink, item *types.ItemNoId) {
	database.GenericInsert(dl, Table, item)
}

func Delete(dl *database.DatabaseLink, id int) {
	database.DeleteById(dl, Table, id)
}

func Get(dl *database.DatabaseLink, id int) (types.Item, error) {
	return database.GenericGet[types.Item](dl, Table, id)
}

func GetAll(dl *database.DatabaseLink) []types.Item {
	return database.GenericGetWhere[types.Item](dl, Table, "true")
}

func GetByCartId(dl *database.DatabaseLink, cart_id int) []types.ItemNoIdCartId {
	return database.GenericGetWhere[types.ItemNoIdCartId](dl, Table, fmt.Sprintf("%s=%d", const_cart_id, cart_id))
}
