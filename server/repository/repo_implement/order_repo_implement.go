package repoimplement

import (
	"database/sql"
	"fmt"
	"main/model"
	"main/repository"
	"sync"
)

type orderRepoImplement struct {
	Db *sql.DB
}

func NewOrderRepo(db *sql.DB) repository.OrderRepo {
	return &orderRepoImplement{
		Db: db,
	}
}

func (orderRepo orderRepoImplement) Create(orders []*model.Order) error {
	var billId int
	moneyOrder := make(chan int, len(orders))
	errChannel := make(chan error, len(orders)) //tạo channel kiểm soát lỗi

	var wg sync.WaitGroup
	for _, order := range orders {
		billId = order.BillId

		wg.Add(1)
		go func(order model.Order, errChannel chan<- error) {

			tx, err := orderRepo.Db.Begin()
			if err != nil {
				errChannel <- err
				return
			}

			var productPrice int

			errWhenGetPrice := tx.QueryRow(
				"select price from products where id = $1",
				order.ProductId,
			).Scan(&productPrice)
			if errWhenGetPrice != nil {
				errChannel <- err
				return
			}

			moneyOrder <- order.ProductId * productPrice

			_, errExec := tx.Exec(
				"insert into orders(bill_id, product_id, quantity, status_order, product_price) values($1, $2, $3, $4, $5)",
				order.BillId, order.ProductId, order.Quantity, order.Status, productPrice,
			)
			if errExec != nil {
				errChannel <- err
				return
			}

			errChannel <- tx.Commit()

			defer wg.Done()
		}(*order, errChannel)
	}

	for i := 0; i < len(orders); i++ {
		if err := <-errChannel; err != nil {
			return err // Trả về lỗi đầu tiên phát hiện được
		}
	}

	wg.Wait()
	close(moneyOrder)

	var billTotal int

	for val := range moneyOrder {
		billTotal += val
	}

	fmt.Println(billId, billTotal)

	NewBillRepo(orderRepo.Db).InsertTotal(billId, billTotal)

	return nil
}
