package repository

import "main/model"

type OrderRepo interface {
	Create([]*model.Order) error
}
