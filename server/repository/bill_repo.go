package repository

import "main/model"

type BillRepo interface {
	Create(email, address string) (int, error)
	GetAllBills(email string) ([]*model.Bill, error)
	GetBillDetail(email string, id int) (*model.Bill, error)

	InsertTotal(id, money int)
}
