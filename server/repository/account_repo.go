package repository

import "main/model"

type AccountRepo interface {
	Login(string, string) (model.Account, error)
	ResetPassword(string, string) error
	SignIn(*model.Account, string) error
}
