package main

import (
	"fmt"

	// module名/パッケージ名
	"goweb/dbController"
)

type User struct {
	id       int
	name     string
	password string
	email    string
}

func main() {
	db := dbController.Connect()

	// 接続が終了したらクローズする
	defer db.Close()

	err := db.Ping()

	if err != nil {
		fmt.Println("データベース接続失敗")
		return
	} else {
		fmt.Println("データベース接続成功")
	}

	//dbController.GetRows(db)
	//dbController.GetSingleRow(db, 7)
	//dbController.InsertUser(db, "fname3", "fpass3", "fmail3")
	//dbController.GetRows(db)

}
