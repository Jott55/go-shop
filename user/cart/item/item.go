package item

import (
	"fmt"
	"jott55/go-shop/database"
	"jott55/go-shop/product"
	"jott55/go-shop/user/cart"
)

const Table = "cart_item"

type Item struct {
	Id       int64
	Quantity int
	Price    int
}

func CreateTable(dl *database.DatabaseLink) {
	field := fmt.Sprintf(`
		id bigint GENERATED ALWAYS AS IDENTITY,
		cart_id BIGINT NOT NULL REFERENCES %v(id) ON DELETE CASCADE,
		product_id BIGINT NOT NULL REFERENCES %v(id) ON DELETE RESTRICT,
		quantity INTEGER NOT NULL CHECK (quantity > 0),
		price NUMERIC(10, 2) NOT NULL,
		UNIQUE (cart_id, product_id)
	`, cart.Table, product.Table)

	database.CreateTable(dl, Table, field)
	database.CreateIndex(dl, Table, []string{"cart_id"})
}

func Drop(dl *database.DatabaseLink) {
	database.DropTable(dl, Table)
}
