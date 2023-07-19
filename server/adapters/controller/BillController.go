package controller

import (
	"fmt"
	"main/adapters/ports/incoming"
	"main/driver"
	"main/model"
	repoimplement "main/repository/repo_implement"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type controllerBill struct{}

func NewControllerBill() controllerBill {
	return controllerBill{}
}

func (controllerBill) Create(db *driver.PostgresDB) echo.HandlerFunc {
	return func(c echo.Context) error {
		var incomingBill incoming.Bill
		if errWhenBinding := c.Bind(&incomingBill); errWhenBinding != nil {
			return c.JSON(http.StatusBadRequest, model.Message{
				Code:   http.StatusBadRequest,
				Text:   "error when binding",
				Output: errWhenBinding.Error(),
			})
		}

		billId, errWhenCreateBill := repoimplement.NewBillRepo(db.SQL).Create(fmt.Sprintf("%v", c.Get("email")), incomingBill.Address)
		if errWhenCreateBill != nil {
			return c.JSON(http.StatusBadRequest, model.Message{
				Code:   http.StatusBadRequest,
				Text:   "error when create bill",
				Output: errWhenCreateBill.Error(),
			})
		}

		orders := make([]*model.Order, len(incomingBill.Orders))

		for i, incommingOrder := range incomingBill.Orders {
			orders[i] = &model.Order{
				BillId:    billId,
				ProductId: incommingOrder.ProductId,
				Quantity:  incommingOrder.Quantity,
				Status:    7,
			}
		}

		errWhenCreateOrdersToBill := repoimplement.NewOrderRepo(db.SQL).Create(orders)
		if errWhenCreateOrdersToBill != nil {
			return c.JSON(http.StatusBadRequest, model.Message{
				Code:   http.StatusBadRequest,
				Text:   "error when insert to db",
				Output: errWhenCreateOrdersToBill.Error(),
			})
		}

		return c.JSON(http.StatusOK, model.Message{
			Code:   http.StatusOK,
			Text:   "oke",
			Output: orders,
		})
	}
}

func (controllerBill) GetAllBills(db *driver.PostgresDB) echo.HandlerFunc {
	return func(c echo.Context) error {
		bills, err := repoimplement.NewBillRepo(db.SQL).GetAllBills(fmt.Sprintf("%v", c.Get("email")))
		if err != nil {
			return c.JSON(http.StatusBadRequest, model.Message{
				Code:   http.StatusBadRequest,
				Text:   "fails",
				Output: err.Error(),
			})
		}

		return c.JSON(http.StatusOK, model.Message{
			Code:   http.StatusOK,
			Text:   "oke",
			Output: bills,
		})
	}
}

func (controllerBill) GetBillDetail(db *driver.PostgresDB) echo.HandlerFunc {
	return func(c echo.Context) error {
		billId, err := strconv.Atoi(fmt.Sprintf("%v", c.Param("bill-id")))
		if err != nil {
			return err
		}
		bills, err := repoimplement.NewBillRepo(db.SQL).GetBillDetail(fmt.Sprintf("%v", c.Get("email")), billId)
		if err != nil {
			return c.JSON(http.StatusBadRequest, model.Message{
				Code:   http.StatusBadRequest,
				Text:   "fails",
				Output: err.Error(),
			})
		}

		return c.JSON(http.StatusOK, model.Message{
			Code:   http.StatusOK,
			Text:   "oke",
			Output: bills,
		})
	}
}
