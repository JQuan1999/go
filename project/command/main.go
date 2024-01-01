package main

import (
	"context"
	"database/sql"
	"fmt"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

const (
	MAXID        = 200
	SELECTOFFSET = 100
	UPDATEOFFSET = 100
	USERNAME     = "test"
	PASSWORD     = "123456"
	NETWORK      = "tcp"
	IP           = "10.177.54.121"
	PORT         = 3388
	TESTTABLE    = "test"
	DATABASE     = "test_ms_proxy"
	TBNAME       = TESTTABLE
	TIMEFORMAT   = "2006-01-02 15:04:05"
)

func main() {
	dsn := fmt.Sprintf("%s:%s@%s(%s:%d)/%s", USERNAME, PASSWORD, NETWORK, IP, PORT, DATABASE)
	db, err := sql.Open("mysql", dsn)

	if err != nil {
		fmt.Printf("dsn:%s invalid, error:%v\n", dsn, err)
		return
	}

	err = db.Ping()
	if err != nil {
		fmt.Printf("open %s failed, err:%v\n", dsn, err)
		return
	}
	db.SetConnMaxLifetime(60 * time.Second)
	TestSlowSql(db)
	// TestTransaction(db)
	// TestConcurrency(db)
	// TestSqlError(db)
}

func TestSlowSql(db *sql.DB) {
	conn, err := db.Conn(context.Background())
	if err != nil {
		panic(err)
	}
	rows, err := conn.QueryContext(context.Background(), "select sleep(100) from test")
	if err != nil {
		fmt.Printf("exec failed, err:%v\n", err)
		return
	}
	rows.Close()
	if err := conn.PingContext(context.Background()); err != nil {
		fmt.Printf("ping failed, err:%v\n", err)
		return
	}
}

func TestTransaction(db *sql.DB) {
	tx, err := db.Begin()
	if err != nil {
		panic(err)
	}
	defer tx.Commit()

	rows, err := tx.Query("select sleep(100) from test")
	if err != nil {
		fmt.Printf("exec failed, err:%v\n", err)
		return
	}
	rows.Close()
}

func TestConcurrency(db *sql.DB) {
	k := 1000
	var wg sync.WaitGroup

	for i := 1; i <= k; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			conn, err := db.Conn(context.Background())
			if err != nil {
				fmt.Printf("get one connection failed, err:%v\n", err)
				return
			}
			defer conn.Close()
			rows, err := conn.QueryContext(context.Background(),
				"select 1 from test")
			if err != nil {
				fmt.Printf("query failed, err:%v\n", err)
				return
			}
			rows.Close()
		}(i)
	}
	wg.Wait()
}

func TestSqlError(db *sql.DB) {
	_, err := db.Query("update 1 from test.id")
	if err != nil {
		fmt.Printf("query failed, err:%v\n", err)
		return
	}
}
