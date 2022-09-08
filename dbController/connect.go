package dbController

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	// エディタで開くときたまに赤線が入るが、動く
	_ "github.com/go-sql-driver/mysql"
)

func Connect() *sql.DB {
	err := godotenv.Load()

	if err != nil {
		fmt.Println(err.Error())
	}
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	host := os.Getenv("BD_HOST")
	port := os.Getenv("DB_PORT")
	database_name := os.Getenv("DB_DATABASE_NAME")
	// [ユーザ名]:[パスワード]@tcp([ホスト名]:[ポート番号])/[データベース名]?charset=[文字コード]
	dbconf := user + ":" + password + "@tcp(" + host + ":" + port + ")/" + database_name + "?charset=utf8mb4"
	db, err := sql.Open("mysql", dbconf)
	if err != nil {
		fmt.Println(err.Error())
	}
	// 重要なセッティングだそうです。.
	db.SetConnMaxLifetime(time.Minute * 3) //接続が再使用される最大時間
	db.SetMaxOpenConns(5)                  // オープン接続の最大数
	db.SetMaxIdleConns(5)                  // アイドル接続プールの最大接続数(オープン接続と同じが推奨)
	return db
}

type User struct {
	id       int
	name     string
	password string
	email    string
}

func GetRows(db *sql.DB) {

	rows, err := db.Query("SELECT * FROM user")
	if err != nil {
		fmt.Printf("getRows db.Query error:%v", err)
	}
	defer rows.Close()

	for rows.Next() {
		u := &User{}
		if err := rows.Scan(&u.id, &u.name, &u.password, &u.email); err != nil {
			fmt.Printf("getRows rows.Scan error:%v", err)
		}
		fmt.Println(u.id, u.name, u.email, u.password)
	}
	err = rows.Err()
	if err != nil {
		fmt.Printf("getRows rows.Err error:%v", err)
	}
}

func GetSingleRow(db *sql.DB, userID int) {
	u := &User{}
	err := db.QueryRow("SELECT * FROM user WHERE id = ?", userID).
		Scan(&u.id, &u.name, &u.email, &u.password)
	if errors.Is(err, sql.ErrNoRows) {
		fmt.Println("getSingleRow no records.")
		return
	}
	if err != nil {
		log.Fatalf("getSingleRow db.QueryRow error err:%v", err)
	}
	fmt.Println(u)
}

func InsertUser(db *sql.DB, name, password, email string) int64 {
	res, err := db.Exec(
		"INSERT INTO user (name, password, email) VALUES (?, ?, ?)",
		name,
		password,
		email,
	)
	if err != nil {
		log.Fatalf("insertUser db.Exec error err:%v", err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		log.Fatalf("insertUser res.LastInsertId error err:%v", err)
	}
	return id
}
