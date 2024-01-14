package model

import (
	"context"
	"errors"
	"gin/session/dbstore"
)

type Account struct {
	WorkId   string `json:"work_id" binding:"required" example:"123"`
	PassWord string `json:"password" binding:"required" example:"passWord"`
}

func CreateAccount(ctx context.Context, account *Account) error {
	db := dbstore.GetDBConn()
	_, err := db.ExecContext(ctx, "insert into session(workid, password) values(?, ?)", account.WorkId, account.PassWord)
	if err != nil {
		return err
	}
	return nil
}

func LoginAccount(ctx context.Context, account *Account) error {
	db := dbstore.GetDBConn()
	row := db.QueryRowContext(ctx, "select password from session where workid = ?", account.WorkId)
	var password string
	if err := row.Scan(&password); err != nil {
		return err
	}
	if account.PassWord != password {
		return errors.New("password error")
	}
	return nil
}
