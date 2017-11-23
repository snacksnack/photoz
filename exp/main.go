package main

import (
	"fmt"

	"../models"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "photoz"
	password = "photoz"
	dbname   = "photoz"
)

type User struct {
	gorm.Model //not inheritance - embedding
	Name       string
	Email      string `gorm:"not null; unique_index"`
	Color      string
	Orders     []Order //orders will not be pull by default - must tell gorm to preload
}

type Order struct {
	gorm.Model
	UserId      uint
	Amount      int
	Description string
}

func main() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	us, err := models.NewUserService(psqlInfo)
	if err != nil {
		panic(err)
	}
	defer us.Close()
	us.DestructiveReset()

	user := models.User{
		Name:  "test1",
		Email: "test1@test.com",
	}
	if err := us.Create(&user); err != nil {
		panic(err)
	}

	//if err := us.Delete(user.ID); err != nil {
	//	panic(err)
	//}

	user.Email = "test1_changed@test.com"
	if err := us.Update(&user); err != nil {
		panic(err)
	}

	userByID, err := us.ById(user.ID)
	if err != nil {
		panic(err)
	}
	fmt.Println(userByID)

	userByEmail, err := us.ByEmail(user.Email)
	if err != nil {
		panic(err)
	}
	fmt.Println(userByEmail)
}
