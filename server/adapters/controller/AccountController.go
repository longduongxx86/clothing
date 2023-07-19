package controller

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"main/adapters/middleware"
	"main/adapters/ports/incoming"
	"main/adapters/ports/outgoing"
	"main/driver"
	"main/model"
	repoimplement "main/repository/repo_implement"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

type controllerAuth struct{}

func NewControllerAuth() controllerAuth {
	return controllerAuth{}
}

func (controllerAuth) Login(db *driver.PostgresDB) echo.HandlerFunc {
	return func(c echo.Context) error {
		accountLogin := new(incoming.AccountLogin)

		if errWhenBinding := c.Bind(accountLogin); errWhenBinding != nil {
			return errWhenBinding
		}

		account, errAccount := repoimplement.NewAccountRepo(db.SQL).Login(accountLogin.Email, accountLogin.Password)

		if errAccount != nil {
			return c.JSON(http.StatusUnauthorized, model.Message{
				Text:   "Tên đăng nhập hoặc mật khẩu không đúng",
				Code:   http.StatusUnauthorized,
				Output: errAccount.Error(),
			})
		} else {

			token, errorWhenCreateToken := middleware.CreateToken(&account)
			if errorWhenCreateToken != nil {
				return c.JSON(http.StatusBadRequest, model.Message{
					Text:   "Có vấn đề xảy ra khi tạo token!",
					Code:   http.StatusBadRequest,
					Output: errorWhenCreateToken.Error(),
				})
			}

			return c.JSON(http.StatusOK, model.Message{
				Text: "Đăng nhập thành công",
				Code: http.StatusOK,
				Output: outgoing.AccountLogin{
					FullName:    account.Fullname,
					Email:       account.Email,
					PhoneNumber: account.PhoneNumber,
					Address:     account.Address,
					Token:       token,
				},
			})
		}

	}
}

func (controllerAuth) ResetPassword(db *driver.PostgresDB) echo.HandlerFunc {
	return func(c echo.Context) error {
		var data map[string]string

		if err := c.Bind(&data); err != nil {
			return c.JSON(http.StatusBadRequest, model.Message{
				Text: err.Error(),
				Code: http.StatusBadRequest,
			})
		}
		newPassword := data["new-password"]

		repoimplement.NewAccountRepo(db.SQL).ResetPassword(fmt.Sprintf("%v", c.Get("email")), newPassword)

		return c.JSON(http.StatusUnauthorized, model.Message{
			Text: "Đổi mật khẩu thành công",
			Code: 200,
		})
	}
}

func (controllerAuth) SignIn(db *driver.PostgresDB) echo.HandlerFunc {
	return func(c echo.Context) error {
		var incomingNewAccount incoming.AccountSignIn

		if errorWhenBinding := c.Bind(&incomingNewAccount); errorWhenBinding != nil {
			return c.JSON(http.StatusBadRequest, model.Message{
				Text:   "Thông tin sai",
				Code:   http.StatusBadRequest,
				Output: incomingNewAccount,
			})
		}
		newAccount := model.Account{
			Fullname:    incomingNewAccount.Fullname,
			Email:       incomingNewAccount.Email,
			PhoneNumber: incomingNewAccount.PhoneNumber,
			Address:     incomingNewAccount.Address,
			Password:    incomingNewAccount.Password,
			RoleId:      0,
		}

		if error := repoimplement.NewAccountRepo(db.SQL).SignIn(&newAccount, generateSalt()); error != nil {
			return c.JSON(http.StatusBadRequest, model.Message{
				Text:   "Lỗi khi đăng ký tài khoản",
				Code:   http.StatusBadRequest,
				Output: error.Error(),
			})
		}

		account, err := repoimplement.NewAccountRepo(db.SQL).Login(newAccount.Email, newAccount.Password)
		if err != nil {
			return c.JSON(http.StatusBadRequest, model.Message{
				Text:   "something went wrong",
				Code:   http.StatusBadRequest,
				Output: err.Error(),
			})
		}

		return c.JSON(http.StatusOK, model.Message{
			Text:   "Đăng ký thành công",
			Code:   http.StatusOK,
			Output: account,
		})

	}
}

func generateSalt() string {
	// Lấy thời gian hiện tại
	timestamp := time.Now().UnixNano()

	// Tạo chuỗi salt từ thời gian hiện tại và một chuỗi ngẫu nhiên
	randomBytes := make([]byte, 32)
	rand.Read(randomBytes)
	salt := base64.StdEncoding.EncodeToString([]byte(fmt.Sprint(timestamp) + string(randomBytes)))

	return salt
}
