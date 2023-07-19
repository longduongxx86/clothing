package incoming

type ProductCreate struct {
	CategoryId  int    ` json:"category_id" form:"category_id" param:"category_id" `
	Title       string ` json:"title" form:"title" param:"title" `
	Price       int    ` json:"price" form:"price" param:"price" `
	Discount    int    ` json:"discount" form:"discount" param:"discount" `
	Thumbnail   string ` json:"thumbnail" form:"thumbnail" param:"thumbnail" `
	Description string ` json:"description" form:"description" param:"description" `
}
type ProductUpdate struct {
	Title       string ` json:"title" form:"title" param:"title" `
	Price       int    ` json:"price" form:"price" param:"price" `
	Discount    int    ` json:"discount" form:"discount" param:"discount" `
	Thumbnail   string ` json:"thumbnail" form:"thumbnail" param:"thumbnail" `
	Description string ` json:"description" form:"description" param:"description" `
}
