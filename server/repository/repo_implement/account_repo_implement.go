package repoimplement

import (
	"crypto/sha256"
	"database/sql"
	"errors"
	"fmt"
	"main/model"
	"main/repository"
)

type AccountRepoImplement struct {
	Db *sql.DB
}

func NewAccountRepo(db *sql.DB) repository.AccountRepo {
	return &AccountRepoImplement{
		Db: db,
	}
}

func (accountRepo AccountRepoImplement) Login(email, password string) (model.Account, error) {
	var err error
	var salt string
	var account model.Account

	err = accountRepo.Db.QueryRow(
		"SELECT salts.salt FROM accounts JOIN salts ON accounts.id = salts.account_id WHERE accounts.email = $1",
		email,
	).Scan(&salt)

	if err != nil {
		return account, err
	}

	result := accountRepo.Db.QueryRow(
		`select id, fullname, phone_number, address, role_id from accounts where email = $1 and password = $2`,
		email, hash(password, salt),
	)

	if result.Err() != nil {
		return account, result.Err()
	}

	err = result.Scan(&account.Id, &account.Fullname, &account.PhoneNumber, &account.Address, &account.RoleId)

	if err != nil {
		return account, err
	}
	account.Email = email

	return account, nil
}

func (accountRepo AccountRepoImplement) ResetPassword(email, newPassword string) error {
	var err error
	var salt string

	err = accountRepo.Db.QueryRow(
		"SELECT salts.salt FROM accounts JOIN salts ON accounts.id = salts.account_id WHERE accounts.email = $1",
		email,
	).Scan(&salt)

	if err != nil {
		return err
	}

	_, err = accountRepo.Db.Exec(`update accounts set password = $1 where email = $2`,
		hash(newPassword, salt), email,
	)

	if err != nil {
		return err
	}

	return nil
}

func (accountRepo AccountRepoImplement) SignIn(account *model.Account, salt string) error {
	var err error
	var count int
	err = accountRepo.Db.QueryRow("SELECT COUNT(*) FROM accounts WHERE email = $1", account.Email).Scan(&count)
	if err != nil {
		return err
	}

	// Kiểm tra xem địa chỉ email đã tồn tại trong bảng hay chưa
	if count > 0 {
		return errors.New("email has already exists")
	}

	var accountId int

	err = accountRepo.Db.QueryRow(`insert into accounts("fullname", "email", "phone_number", "address", "password") values ($1, $2, $3, $4, $5) returning id`,
		account.Fullname, account.Email, account.PhoneNumber, account.Address, hash(account.Password, salt),
	).Scan(&accountId)

	if err != nil {
		return err
	}

	_, err = accountRepo.Db.Exec("INSERT INTO salts (account_id, salt) VALUES ($1, $2)", accountId, salt)
	if err != nil {
		return err
	}

	return nil
}

func hash(password, salt string) string {
	return fmt.Sprintf("%x", sha256.Sum256([]byte(password+salt)))
}
