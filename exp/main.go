package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "photoz"
	password = "photoz"
	dbname   = "photoz"
)

func main() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	/*
		_, err = db.Exec(`
			INSERT INTO users(name, email)
			VALUES($1, $2)`, "Reid Collins", "test0@test.com")
		if err != nil {
			panic(err)
		}
	*/

	//you can get the newly created id using QueryRow & scan
	/*
		var id int
		err = db.QueryRow(`
			INSERT INTO users(name, email)
			VALUES($1, $2)
			RETURNING id`,
			"Reid2 Collins", "test2@test.com").Scan(&id)
		if err != nil {
			panic(err)
		}
		fmt.Print("the newly created id is: ", id)
	*/

	/*querying the db for single row
	var id int
	var name, email string
	row := db.QueryRow(`
		SELECT id, name, email
		FROM users
		where id=$1`, 1)
	err = row.Scan(&id, &name, &email)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("no records found")
		} else {
			panic(err)
		}
	}
	*/

	//querying db for multiple rows
	type User struct {
		ID    int
		Name  string
		Email string
	}

	//var id int
	//var name, email string
	var users []User
	rows, err := db.Query(`
		SELECT id, name, email
		FROM users`)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	for rows.Next() {
		var user User
		err = rows.Scan(&user.ID, &user.Name, &user.Email)
		if err != nil {
			panic(err)
		}
		users = append(users, user)
	}
	if rows.Err() != nil {
		fmt.Println(rows.Err())
	}
	fmt.Println(users)

	/*
		creating a relationship between users and orders
		for i := 1; i <= 6; i++ {
			userID := 1
			if i > 3 {
				userID = 5
			}
			amount := i * 100
			description := fmt.Sprintf("usb adapter x%d", i)

			_, err = db.Query(`
				INSERT INTO orders(user_id, amount, description)
				VALUES($1, $2, $3)`, userID, amount, description)
			if err != nil {
				panic(err)
			}
		}
	*/

	//querying relational data records
	rows, err = db.Query(`
		SELECT *
		FROM users
		INNER JOIN orders ON users.id=orders.id`)
	if err != nil {
		panic(err)
	}
	for rows.Next() {
		var userID, orderID, amount int
		var email, name, description string
		if err := rows.Scan(&userID, &name, &email, &orderID, &userID, &amount, &description); err != nil {
			panic(err)
		}
		fmt.Println("userID: ", userID, "name: ", name, "email: ", email,
			"orderID: ", orderID, "amount: ", amount, "description: ", description)
	}
	if rows.Err() != nil {
		panic(rows.Err())
	}
}
