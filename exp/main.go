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
}
