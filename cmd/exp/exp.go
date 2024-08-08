package main

import (
	"database/sql"
	"fmt"

	_ "github.com/jackc/pgx/v4/stdlib"
)

type Order struct {
	ID          int
	UserID      int
	Amount      int
	Description string
}

func main() {
	var orders []Order
	db, err := sql.Open("pgx", "host=localhost port=5432 user=rwauser password=rwapassword dbname=rwa sslmode=disable")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	fmt.Println("db connected")
	userID := 1
	rows, err := db.Query(`
		SELECT id, amount, description
		FROM orders
		WHERE user_id=$1
	`, userID)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	for rows.Next() {
		var order Order
		order.UserID = userID
		err := rows.Scan(&order.ID, &order.Amount, &order.Description)
		if err != nil {
			panic(err)
		}
		orders = append(orders, order)
	}
	err = rows.Err()
	if err != nil {
		panic(err)
	}
	fmt.Println(orders)
}
