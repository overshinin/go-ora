package main

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "github.com/overshinin/go-ora/v2"
)

func main() {
	conn, err := sql.Open("oracle", os.Getenv("DSN"))
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func() {
		err = conn.Close()
		if err != nil {
			fmt.Println("can't close connection: ", err)
		}
	}()
	var timestamp time.Time
	err = conn.QueryRow(`select systimestamp AT TIME ZONE 'Asia/Kolkata' from dual`).Scan(&timestamp)
	if err != nil {
		fmt.Println("can't query row: ", err)
		return
	}
	fmt.Println(timestamp)
}
