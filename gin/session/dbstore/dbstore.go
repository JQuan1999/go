package dbstore

import (
	"database/sql"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var dsn = "test:123456@tcp(10.177.54.121:3311)/test_ms_proxy"

var dbstore *DBStore

type DBStore struct {
	db *sql.DB
}

func InitDBStore() {
	var err error
	dbstore, err = NewDBStore()
	if err != nil {
		panic(err)
	}
}

func GetDBConn() *sql.DB {
	return dbstore.getDBConn()
}

func NewDBStore() (*DBStore, error) {
	dbstore := DBStore{}
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	db.SetMaxIdleConns(128)
	db.SetConnMaxLifetime(120 * time.Second)
	dbstore.db = db
	return &dbstore, nil
}

func (dbstore *DBStore) getDBConn() *sql.DB {
	return dbstore.db
}
