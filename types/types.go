package types

type ProductRequest struct {
	Product *ProductNoId
}

type ProductIdRequest struct {
	Product_id int
}

type UserRequest struct {
	User *UserNoId
}

type CartRequest struct {
	Cart *CartNoId
}

type ItemRequest struct {
	Item *ItemNoId
}

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

type ProductNoId struct {
	Name        string
	Image_url   string
	Price       int
	Description string
}

type ProductLess struct {
	Name      string
	Image_url string
}

type ProductItem struct {
	Id        int
	Name      string
	Image_url string
	Price     int
	Quantity  int
}

type Item struct {
	Id         int64
	Cart_id    int
	Product_id int
	Quantity   int
	Price      int
}

type ItemNoId struct {
	Cart_id    int
	Product_id int
	Quantity   int
	Price      int
}

type ItemNoIdCartId struct {
	Product_id int
	Quantity   int
	Price      int
}
type CartId struct {
	Id int
}

type Cart struct {
	Id      int
	User_id int
}

type CartNoId struct {
	User_id int
}

type User struct {
	Id        int
	Name      string
	Email     string
	Password  string // plan on using argon2id
	Photo_url string
}

type UserNoId struct {
	Name      string
	Email     string
	Password  string
	Photo_url string
}

type LoginUser struct {
	Email    string
	Password string
}
