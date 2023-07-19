package repository

import "main/model"

type ProductRepoWithoutCache interface {
	Create(*model.Product) error
	Update(*model.Product) error
	DeleteProduct(id string) error
}

type ProductRepoWithCache interface {
	GetProductDetail(id string) (*model.Product, error)
	GetListProductWithCondition(*model.ConditionSelected) ([]*model.Product, error)
}
