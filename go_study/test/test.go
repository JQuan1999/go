package test

import (
	"database/sql"
	"fmt"
	"math/rand"

	_ "github.com/go-sql-driver/mysql"
)

var letters = []rune("abcdefghijklmnopqrstuvwxyz")

type User struct {
	Id   int
	Name string
	Age  int
}

func RandName() string {
	n := rand.Intn(10)
	name := make([]rune, n)
	for idx := range name {
		name[idx] = letters[rand.Intn(len(letters))]
	}
	return string(name)
}

func TestInsert(db *sql.DB, age int, name string) {
	sqlStr := fmt.Sprintf("insert into tbuser(age, name) values(%d,'%s')", age, name)
	_, err := db.Exec(sqlStr)
	if err != nil {
		fmt.Printf("exec [sql: %s] failed, err:%v\n", sqlStr, err)
		return
	}
}

func TestUpdate(db *sql.DB) {
	sqlStr := "update tbuser set age = 10 where id < 100"
	_, err := db.Exec(sqlStr)
	if err != nil {
		fmt.Printf("exec [sql: %s] failed, err:%v\n", sqlStr, err)
		return
	}
}

func TestSelectInnerJoin(db *sql.DB) {
	sqlStr := "select count(*) from tbuser t1 join (select * from tbuser where id < 100) as t2 on t1.age = t2.age where t1.id < 1000"
	_, err := db.Exec(sqlStr)
	if err != nil {
		fmt.Printf("exec [sql: %s] failed, err:%v\n", sqlStr, err)
		return
	}
}

func TestSelectSleep(db *sql.DB, sleepTime int) {
	sqlStr := fmt.Sprintf("select (%d) from tbuser where id = 1000", sleepTime)
	_, err := db.Exec(sqlStr)
	if err != nil {
		fmt.Printf("exec [sql: %s] failed, err:%v\n", sqlStr, err)
		return
	}
}

func TestSelectNormal(db *sql.DB) {
	sqlStr := "select * from tbuser where id < 10"
	rows, err := db.Query(sqlStr)
	if err != nil {
		fmt.Printf("exec [sql: %s] failed, err:%v\n", sqlStr, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var user User
		rows.Scan(&user.Id, &user.Name, &user.Age)
		fmt.Printf("use id: %d, name: %s, age: %d\n", user.Id, user.Name, user.Age)
	}
}

func TestTransaction(db *sql.DB) {
	tx, err := db.Begin()
	if err != nil {
		if tx != nil {
			_ = tx.Rollback()
		}
		fmt.Printf("begin trasaction failed, err:%v\n", err)
		return
	}

	sqlStr1 := "Update tbuser set age = 10 where id = 1"
	rs1, err := tx.Exec(sqlStr1)
	if err != nil {
		_ = tx.Rollback()
		fmt.Printf("exec failed, err:%v\n", err)
		return
	}
	row1, err := rs1.RowsAffected()
	if err != nil {
		_ = tx.Rollback()
		fmt.Printf("exec failed, err:%v\n", err)
		return
	}
	fmt.Println(row1)

	sqlStr2 := "Update tbuser set age = 10 where id = 100"
	rs2, err := tx.Exec(sqlStr2)
	if err != nil {
		_ = tx.Rollback()
		fmt.Printf("exec failed, err:%v\n", err)
		return
	}
	row2, err := rs2.RowsAffected()
	if err != nil {
		_ = tx.Rollback()
		fmt.Printf("exec failed, err:%v\n", err)
		return
	}
	fmt.Println(row2)
}
