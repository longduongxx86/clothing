package repoimplement

import (
	"database/sql"
	"errors"
	"main/model"
	"main/repository"
	"strings"
)

type billRepoImplement struct {
	Db *sql.DB
}

func NewBillRepo(db *sql.DB) repository.BillRepo {
	return &billRepoImplement{
		Db: db,
	}
}

func (billRepo billRepoImplement) Create(email, address string) (int, error) {
	accountId, errAccountId := billRepo.getAccountId(email)
	if errAccountId != nil {
		return 0, nil
	}

	var billId int
	if err := billRepo.Db.QueryRow("insert into bills (account_id, address) values ($1, $2) returning id", accountId, address).Scan(&billId); err != nil {

		return 0, err
	}

	return billId, nil
}

func (billRepo billRepoImplement) GetAllBills(email string) ([]*model.Bill, error) {

	// lấy account_id
	accountId, errorWhenGetAccountIt := billRepo.getAccountId(email)
	if errorWhenGetAccountIt != nil {
		return nil, errorWhenGetAccountIt
	}

	// từ account_id, lấy danh sách (billId, billTotal, billAddres)
	billsRows, errWhenGetBills := billRepo.Db.Query(
		"select id, total, address from bills where account_id = $1",
		accountId,
	)
	if errWhenGetBills != nil {
		return nil, errWhenGetBills
	}

	defer billsRows.Close()

	var bills []*model.Bill

	for billsRows.Next() {
		var billId, billTotal int
		var address string

		errScanBill := billsRows.Scan(&billId, &billTotal, &address)

		if errScanBill != nil {
			return nil, errScanBill
		}

		bill := new(model.Bill)
		bill.Address = strings.TrimSpace(address)
		bill.Total = billTotal
		bill.Id = billId

		orders, err := billRepo.getOrdersWithBillId(billId)
		if err != nil {
			return nil, err
		}

		bill.Orders = orders

		bills = append(bills, bill)
	}

	return bills, nil
}

func (billRepo billRepoImplement) GetBillDetail(email string, billId int) (*model.Bill, error) {
	// lấy account_id
	accountId, errorWhenGetAccountIt := billRepo.getAccountId(email)
	if errorWhenGetAccountIt != nil {
		return nil, errorWhenGetAccountIt
	}

	billRows, errBill := billRepo.Db.Query(
		"select total, address from bills where account_id = $1 and id = $2",
		accountId, billId,
	)

	if errBill != nil {
		return nil, errBill
	}

	if !billRows.Next() {
		return nil, errors.New("không tồn tại")
	}

	var billTotal int
	var billAddress string
	errScanBill := billRows.Scan(&billTotal, &billAddress)
	if errScanBill != nil {
		return nil, errScanBill
	}

	bill := new(model.Bill)
	bill.Id = billId
	bill.Total = billTotal
	bill.Address = billAddress

	var err error
	bill.Orders, err = billRepo.getOrdersWithBillId(billId)
	if err != nil {
		return nil, err
	}

	return bill, nil
}

func (billRepo billRepoImplement) getAccountId(email string) (int, error) {
	var accountId int

	err := billRepo.Db.QueryRow("select id from accounts where email = $1", email).Scan(&accountId)
	if err != nil {
		return 0, err
	}
	return accountId, nil
}

func (billRepo billRepoImplement) getOrdersWithBillId(billId int) ([]*model.OrderClient, error) {
	orderRows, errWhenGetOrders := billRepo.Db.Query(
		"select products.title, orders.quantity, status.name from orders join products on orders.product_id = products.id join status on orders.status_order = status.id where bill_id= $1;",
		billId,
	)
	if errWhenGetOrders != nil {
		return nil, errWhenGetOrders
	}

	var ordersClient []*model.OrderClient

	for orderRows.Next() {
		order := new(model.OrderClient)

		err := orderRows.Scan(&order.Title, &order.Quantity, &order.Status)
		if err != nil {
			return nil, err
		}
		ordersClient = append(ordersClient, order)
	}

	return ordersClient, nil
}

func (billRepo billRepoImplement) InsertTotal(id, money int) {
	billRepo.Db.Exec("update bills set total = $1 where id = $2", money, id)
}
