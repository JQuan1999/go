package history

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/juju/errors"
)

var dsn = "DBProxyServiceData:71ljq9WfljT0u05H4w3I@tcp(dcdb.qa.17usoft.com:20034)/DBProxyServiceData"

var dbstore *DBStore

type DBStore struct {
	db *sql.DB
}

func init() {
	var err error
	dbstore, err = NewDBStore()
	if err != nil {
		panic(err)
	}
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

// ProxyInstance 定义Proxy运行实例
type ProxyInstance struct {
	// Host Proxy IP 地址
	Host string `json:"host" example:"10.100.217.17"`
	// Port Proxy 监听端口
	Port int32 `json:"port" example:"3306"`
	// Idc Proxy 部署机房
	Idc string `json:"idc" example:"xhy"`
	// Offlined Proxy下线标志，0：online；1：Offline
	Offlined   int    `json:"-"`
	GroupName  string `json:"group_name"`
	CollectSql int    `json:"collect_sql"`
}

func (s *DBStore) GetProxyInst() ([]string, error) {
	db := s.db
	if db == nil {
		return nil, errors.New("db is nil")
	}
	rows, err := db.QueryContext(context.Background(), `select address from ProxyInstance where isDeleted = 0`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var proxyAddress []string
	var address string
	for rows.Next() {
		if err := rows.Scan(&address); err != nil {
			fmt.Println("scan failed")
			return nil, err
		}
		proxyAddress = append(proxyAddress, address)
	}
	return proxyAddress, nil
}
