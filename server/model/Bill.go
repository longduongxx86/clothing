package model

type OrderClient struct {
	Title    string
	Quantity int
	Status   string
}

type Bill struct {
	Id      int
	Orders  []*OrderClient
	Total   int
	Address string
}
