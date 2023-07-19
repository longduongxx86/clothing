package incoming

type order struct {
	ProductId int `json:"productId" form:"productId"`
	Quantity  int `json:"quantity" form:"quantity"`
}

type Bill struct {
	Orders  []order
	Address string
}
