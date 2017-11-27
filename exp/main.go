package main

import (
	"fmt"

	"../models"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "photoz"
	password = "photoz"
	dbname   = "photoz_test"
)

func main() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	us, err := models.NewUserService(psqlInfo)
	if err != nil {
		panic(err)
	}
	defer us.Close()
	//us.DestructiveReset()
	//us.AutoMigrate()

	user := models.User{
		Name:     "test100",
		Email:    "test100@test.com",
		Password: "test100",
		Remember: "abc123",
	}
	err = us.Create(&user)
	if err != nil {
		panic(err)
	}

	user2, err := us.ByRemember("abc123")
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", user2)
}
